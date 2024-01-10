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
		return nil, nil, err
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

var MinioPreviewClient *minio.Client
var MinioPreviewCore *minio.Core

func InitMinioPreviewClient() {
	MinIoUrl := os.Getenv("MINIO_URL")
	MinioSsl := os.Getenv("MINIO_SSL")
	MINIO_PREVIEW_USER := os.Getenv("MINIO_PREVIEW_USER")
	MINIO_PREVIEW_PASSWORD := os.Getenv("MINIO_PREVIEW_PASSWORD")
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
		Creds:  credentials.NewStaticV4(MINIO_PREVIEW_USER, MINIO_PREVIEW_PASSWORD, ""),
		Secure: ssl,
	})

	if err != nil {
		os.Exit(1)
	}

	minIoCore, err := minio.NewCore(MinIoUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(MINIO_PREVIEW_USER, MINIO_PREVIEW_PASSWORD, ""),
		Secure: ssl,
	})

	if err != nil {
		os.Exit(1)
	}
	MinioPreviewClient = minioClient
	MinioPreviewCore = minIoCore
	// return minioClient, minIoCore, nil
}

func Close(file multipart.File) {
	err := file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
}
