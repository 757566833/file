package services

import (
	"bytes"
	"context"
	"encoding/json"
	"file/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"net/http"
	"os"
)

const ContentJpgType = "application/jpg"
const contentJSONType = "application/json"

var BuketName = os.Getenv("BUKET_NAME")

func Upload(c *gin.Context) {
	// todo
	//objectName := "222.jpg"
	//filePath := "/go/test/222.jpg"
	//
	//// Upload the zip file with FPutObject
	//
	//n, err := db.MinIoClient.FPutObject(context.Background(), BuketName, objectName, filePath, minio.PutObjectOptions{ContentType: ContentJpgType})
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//c.IndentedJSON(http.StatusOK, n)
}
func Preview(c *gin.Context) {
	file := c.Param("file")
	if file == "" {
		c.IndentedJSON(http.StatusBadRequest, "")
	}

	object, err := db.MinIoClient.GetObject(context.Background(), BuketName, file, minio.GetObjectOptions{})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	objectInfo, err := object.Stat()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	defer object.Close()
	extraHeaders := map[string]string{
		"Content-Disposition": "inline",
	}
	c.DataFromReader(http.StatusOK, objectInfo.Size, contentJSONType, object, extraHeaders)

}
func Download(c *gin.Context) {
	file := c.Param("file")
	if file == "" {
		c.IndentedJSON(http.StatusBadRequest, "")
	}
	object, err := db.MinIoClient.GetObject(context.Background(), BuketName, file, minio.GetObjectOptions{})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	objectInfo, err := object.Stat()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	defer object.Close()
	extraHeaders := map[string]string{
		"Content-Disposition": "attachment; filename=" + file,
	}
	c.DataFromReader(http.StatusOK, objectInfo.Size, contentJSONType, object, extraHeaders)
}

type CreateJsonStruct struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

func CreateJson(c *gin.Context) {

	var requestBody CreateJsonStruct
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	name := requestBody.Name
	_object, err := db.MinIoClient.GetObject(context.Background(), BuketName, name, minio.GetObjectOptions{})
	defer _object.Close()
	fmt.Println(err)
	fmt.Println(_object)
	if err == nil && _object != nil {
		c.IndentedJSON(http.StatusBadRequest, name+" existed")
		return
	}
	data := requestBody.Data
	str, err := json.Marshal(data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	_data := bytes.NewReader(str)

	object, err := db.MinIoClient.PutObject(context.Background(), BuketName, name, _data, int64(len(str)), minio.PutObjectOptions{ContentType: contentJSONType})

	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	c.IndentedJSON(http.StatusOK, object)
}
func ForceJson(c *gin.Context) {

	var requestBody CreateJsonStruct
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	name := requestBody.Name
	data := requestBody.Data
	str, err := json.Marshal(data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	_data := bytes.NewReader(str)

	object, err := db.MinIoClient.PutObject(context.Background(), BuketName, name, _data, int64(len(str)), minio.PutObjectOptions{ContentType: contentJSONType})
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	c.IndentedJSON(http.StatusOK, object)
}
