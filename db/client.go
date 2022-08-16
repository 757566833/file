package db

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"strconv"
)

var MinIoClient *minio.Client

var BuketName = os.Getenv("BUKET_NAME")
var Region = os.Getenv("Region")

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
	bucketName := BuketName

	exists, err := minioClient.BucketExists(context.Background(), bucketName)

	if err == nil && exists {
		log.Printf("We already own %s\n", bucketName)
	} else {
		op := new(minio.MakeBucketOptions)
		op.Region = Region

		err = minioClient.MakeBucket(context.Background(), bucketName, *op)
	}

	return minioClient
}
