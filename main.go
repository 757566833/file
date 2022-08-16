package main

import (
	"file/db"
	"file/route"
	"os"
)

var BuketName = os.Getenv("BUKET_NAME")
var Region = os.Getenv("BUKET_NAME")

func main() {
	db.InitDB()
	ExplorerServerPort := os.Getenv("EXPLORER_SERVER_PORT")
	//ExplorerServerPort := "8090"
	router := route.InitRouter()
	router.Run("0.0.0.0:" + ExplorerServerPort)
}
