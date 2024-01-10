package main

import (
	"file/db"
	"file/log"
	"file/route"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	log.InitLogger()
	db.InitMinioPreviewClient()

	if err != nil {
		os.Exit(1)
	}
	ExplorerServerPort := os.Getenv("EXPLORER_SERVER_PORT")
	//ExplorerServerPort := "8090"
	router := route.InitRouter()
	err = router.Run("0.0.0.0:" + ExplorerServerPort)
	if err != nil {
		log.Logger.Error(err.Error())
	}
}
