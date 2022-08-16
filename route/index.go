package route

import (
	"file/services"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())
	//router.POST("/upload", controller.Upload)
	router.GET("/preview/:file", services.Preview)
	router.GET("/download", services.Download)
	router.POST("/create/json", services.CreateJson)
	router.POST("/force/json", services.ForceJson)
	return router
}
