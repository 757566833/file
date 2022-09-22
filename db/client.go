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
	USER := os.Getenv("USER")
	PASSWORD := os.Getenv("PASSWORD")
	SSL := os.Getenv("SSL")
	ssl, _ := strconv.ParseBool(SSL)
	//MinIoUrl := "127.0.0.1:9000"
	//USER := "root"
	//PASSWORD := "123456789"
	//ssl := false
	// Initialize minio client object.
	minioClient, err := minio.New(MinIoUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(USER, PASSWORD, ""),
		Secure: ssl,
	})

	if err != nil {
		log.Fatalln(err)
	}
	MinIoClient = minioClient

	return minioClient
}
