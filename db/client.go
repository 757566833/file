package db

import (
	"fmt"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinioClient(user string, password string) (*minio.Client, *minio.Core, error) {
	MinIoUrl := os.Getenv("MINIO_URL")
	MinioSsl := os.Getenv("MINIO_SSL")
	ssl, err := strconv.ParseBool(MinioSsl)
	if err != nil {
		os.Exit(1)
	}
	//MinIoUrl := "127.0.0.1:9000"
	//USER := "root"
	//PASSWORD := "123456789"
	//ssl := false
	// Initialize minio client object.
	minioClient, err := minio.New(MinIoUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(user, password, ""),
		Secure: ssl,
	})

	if err != nil {
		return nil, nil, err
	}

	minIoCore, err := minio.NewCore(MinIoUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(user, password, ""),
		Secure: ssl,
	})

	if err != nil {
		return nil, nil, err
	}

	return minioClient, minIoCore, nil
}

func Close(file multipart.File) {
	err := file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
}
