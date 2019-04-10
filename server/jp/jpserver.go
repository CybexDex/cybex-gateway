package jp

import (
	"bytes"
	"fmt"

	"bitbucket.org/woyoutlz/bbb-gateway/utils/eventlog"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// StartServer ...
func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/api/order/noti", func(c *gin.Context) {
		// 记录日志
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		str := buf.String()
		eventlog.Log("jpnoti", str)
		// 充提记录的进一步处理
		// 返回
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	port := viper.GetString("jpserver.port")
	fmt.Println("jpserver start at", port)
	err := r.Run(port) // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println(err)
	}
}
