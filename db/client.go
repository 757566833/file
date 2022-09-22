package db

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"strconv"
)

var MinIoClient *minio.Client

func InitDB() *minio.Client {
	MinIoUrl := os.Getenv("MINIO_URL")
	MinioRootUser := os.Getenv("MINIO_ROOT_USER")
	MinioRootPassword := os.Getenv("MINIO_ROOT_PASSWORD")
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
		Creds:  credentials.NewStaticV4(MinioRootUser, MinioRootPassword, ""),
		Secure: ssl,
	})

	if err != nil {
		log.Fatalln(err)
	}
	MinIoClient = minioClient

	return minioClient
}
