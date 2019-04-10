package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// StartServer ...
func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	port := viper.GetString("userserver.port")
	fmt.Println("userserver start at", port)
	r.Run(port) // listen and serve on 0.0.0.0:8080
}
