package services

import (
	"bytes"
	"context"
	"encoding/json"
	"file/db"
	"file/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"net/http"
)

const contentJSONType = "application/json"

func Upload(c *gin.Context) {
	buket := c.Param("buket")
	f, err := c.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	io, err := f.Open()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	defer io.Close()
	_type, err := utils.GetFileContentType(io)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	info, err := db.MinIoClient.PutObject(context.Background(), buket, f.Filename, io, f.Size, minio.PutObjectOptions{ContentType: _type})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	c.IndentedJSON(http.StatusOK, info)
}

func Preview(c *gin.Context) {
	buket := c.Param("buket")
	file := c.Param("file")
	if file == "" {
		c.IndentedJSON(http.StatusBadRequest, "")
	}

	object, err := db.MinIoClient.GetObject(context.Background(), buket, file, minio.GetObjectOptions{})
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
	buket := c.Param("buket")
	file := c.Param("file")
	if file == "" {
		c.IndentedJSON(http.StatusBadRequest, "")
	}
	object, err := db.MinIoClient.GetObject(context.Background(), buket, file, minio.GetObjectOptions{})
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
	buket := c.Param("buket")
	var requestBody CreateJsonStruct
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	name := requestBody.Name
	_object, err := db.MinIoClient.GetObject(context.Background(), buket, name, minio.GetObjectOptions{})
	defer _object.Close()
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

	object, err := db.MinIoClient.PutObject(context.Background(), buket, name, _data, int64(len(str)), minio.PutObjectOptions{ContentType: contentJSONType})

	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	c.IndentedJSON(http.StatusOK, object)
}
func ForceJson(c *gin.Context) {
	buket := c.Param("buket")
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

	object, err := db.MinIoClient.PutObject(context.Background(), buket, name, _data, int64(len(str)), minio.PutObjectOptions{ContentType: contentJSONType})
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	c.IndentedJSON(http.StatusOK, object)
}
