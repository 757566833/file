package route

import (
	"encoding/base64"
	"file/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			c.Next()
			return
		} else {
			// 获取 Authorization 头的值
			authHeader := c.GetHeader("Authorization")

			// 检查 Authorization 头是否存在
			if authHeader == "" {
				c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
				return
			}

			// 检查 Authorization 头的格式是否为 "Basic <base64-encoded-credentials>"
			authParts := strings.Fields(authHeader)
			if len(authParts) != 2 || authParts[0] != "Basic" {
				c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Authorization Header"})
				return
			}

			// 解码 base64 编码的凭据
			credentials, err := base64.StdEncoding.DecodeString(authParts[1])
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Base64 Encoding"})
				return
			}

			// 将凭据拆分为用户名和密码
			credParts := strings.SplitN(string(credentials), ":", 2)
			if len(credParts) != 2 {
				c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Credentials Format"})
				return
			}

			username := credParts[0]
			password := credParts[1]

			// 在实际应用中，你可以在这里验证用户名和密码
			// 这里只是一个简单的示例，实际应用中应使用更安全的认证方式

			// fmt.Println("Username:", username)
			// fmt.Println("Password:", password)

			// 将用户名和密码存储到 Context 中，以便在后续的处理函数中使用
			c.Set("username", username)
			c.Set("password", password)

			// 认证通过，继续执行后续中间件和处理程序
			c.Next()
		}

	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
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
	router.Use(authMiddleware())
	router.GET("/all/:bucket", services.All)
	router.POST("/upload/:bucket", services.Upload)
	router.GET("/preview/:bucket/:file", services.Preview)
	router.GET("/download/:bucket/:file", services.Download)
	router.POST("/create/json/:bucket", services.CreateJson)
	router.POST("/force/json/:bucket", services.ForceJson)
	return router
}
