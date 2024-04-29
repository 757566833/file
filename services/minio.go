package services

import (
	"bytes"
	"context"
	"encoding/json"
	"file/db"
	"file/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

const Region = "local"

const contentJSONType = "application/json"

func checkBucket(user string, password string, bucket string) error {
	client, _, err := db.InitMinioClient(user, password)
	if err != nil {
		return err
	}
	exists, err := client.BucketExists(context.Background(), bucket)
	if err == nil && exists {
		// 已存在
	} else {
		op := new(minio.MakeBucketOptions)
		op.Region = Region

		err = client.MakeBucket(context.Background(), bucket, *op)
		if err != nil {
			return err
		}
	}
	return nil
}

type UploadResponse struct {
	Preview string `json:"preview"`
}

func Upload(c *gin.Context) {
	bucket := c.Param("bucket")
	usernameRaw, ok := c.Get("username")
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	passwordRaw, ok := c.Get("password")
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	username, ok := usernameRaw.(string)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	password, ok := passwordRaw.(string)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	err := checkBucket(username, password, bucket)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	f, err := c.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	fileIo, err := f.Open()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	_type, err := utils.GetFileContentType(fileIo)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	defer db.Close(fileIo)

	name := utils.RandStr(32)
	split := strings.Split(f.Filename, ".")

	if len(split) > 1 {
		name = name + "." + split[len(split)-1]
	}
	client, _, err := db.InitMinioClient(username, password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	_, err = client.PutObject(context.Background(), bucket, name, fileIo, f.Size, minio.PutObjectOptions{ContentType: _type})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, UploadResponse{Preview: "/preview/" + bucket + "/" + name})
}

func Preview(c *gin.Context) {
	bucket := c.Param("bucket")
	file := c.Param("file")
	if file == "" {
		c.IndentedJSON(http.StatusBadRequest, "")
		return
	}

	object, err := db.MinioPreviewClient.GetObject(context.Background(), bucket, file, minio.GetObjectOptions{})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	objectInfo, err := object.Stat()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close(object)
	extraHeaders := map[string]string{
		"Content-Disposition": "inline",
	}
	for key, value := range objectInfo.Metadata {
		extraHeaders[key] = strings.Join(value, ";")
	}
	c.DataFromReader(http.StatusOK, objectInfo.Size, objectInfo.Metadata.Get("Content-Type"), object, extraHeaders)

}

type DeleteFIleStruct struct {
	Filename string `json:"filename"`
}

func Delete(c *gin.Context) {
	// c.IndentedJSON(http.StatusOK, "success")
	bucket := c.Param("bucket")
	usernameRaw, ok := c.Get("username")
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	passwordRaw, ok := c.Get("password")
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	username, ok := usernameRaw.(string)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	password, ok := passwordRaw.(string)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	err := checkBucket(username, password, bucket)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	var requestBody DeleteFIleStruct
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	client, _, err := db.InitMinioClient(username, password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = client.RemoveObject(context.Background(), bucket, requestBody.Filename, minio.RemoveObjectOptions{})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, "success")
}

func Download(c *gin.Context) {
	bucket := c.Param("bucket")
	file := c.Param("file")
	if file == "" {
		c.IndentedJSON(http.StatusBadRequest, "")
		return
	}

	object, err := db.MinioPreviewClient.GetObject(context.Background(), bucket, file, minio.GetObjectOptions{})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	objectInfo, err := object.Stat()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	defer db.Close(object)

	extraHeaders := map[string]string{
		"Content-Disposition": "attachment; filename=" + file,
	}
	for key, value := range objectInfo.Metadata {
		extraHeaders[key] = strings.Join(value, ";")
	}
	c.DataFromReader(http.StatusOK, objectInfo.Size, objectInfo.Metadata.Get("Content-Type"), object, extraHeaders)
}

func All(c *gin.Context) {
	bucket := c.Param("bucket")
	object, err := db.MinioPreviewCore.ListObjects(bucket, "", "", "", 1000)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, object.Contents)
}

type CreateJsonStruct struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

func CreateJson(c *gin.Context) {
	bucket := c.Param("bucket")

	usernameRaw, ok := c.Get("username")
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	passwordRaw, ok := c.Get("password")
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	username, ok := usernameRaw.(string)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	password, ok := passwordRaw.(string)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	err := checkBucket(username, password, bucket)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	var requestBody CreateJsonStruct
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	name := requestBody.Name
	client, _, err := db.InitMinioClient(username, password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	_object, err := client.GetObject(context.Background(), bucket, name, minio.GetObjectOptions{})
	defer db.Close(_object)
	if err == nil && _object != nil {
		c.IndentedJSON(http.StatusBadRequest, name+" existed")
		return
	}
	data := requestBody.Data
	str, err := json.Marshal(data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	_data := bytes.NewReader(str)

	object, err := client.PutObject(context.Background(), bucket, name, _data, int64(len(str)), minio.PutObjectOptions{ContentType: contentJSONType})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, object)
}
func ForceJson(c *gin.Context) {
	bucket := c.Param("bucket")

	usernameRaw, ok := c.Get("username")
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	passwordRaw, ok := c.Get("password")
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	username, ok := usernameRaw.(string)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	password, ok := passwordRaw.(string)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, "")
		return
	}
	err := checkBucket(username, password, bucket)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	var requestBody CreateJsonStruct
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	name := requestBody.Name
	data := requestBody.Data
	str, err := json.Marshal(data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	_data := bytes.NewReader(str)
	client, _, err := db.InitMinioClient(username, password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	object, err := client.PutObject(context.Background(), bucket, name, _data, int64(len(str)), minio.PutObjectOptions{ContentType: contentJSONType})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, object)
}
