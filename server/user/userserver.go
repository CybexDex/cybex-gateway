package user

import (
	userc "bitbucket.org/woyoutlz/bbb-gateway/controller/user"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// StartServer ...
func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/users/:user/assets/:asset/address", func(c *gin.Context) {
		user := c.Param("user")
		asset := c.Param("asset")
		address, err := userc.GetAddress(user, asset)
		if err != nil {
			log.Errorln("Error", err)
			c.JSON(400, gin.H{
				"message": err,
			})
		}
		c.JSON(200, address)
	})
	port := viper.GetString("userserver.port")
	log.Infoln("userserver start at", port)
	r.Run(port) // listen and serve on 0.0.0.0:8080
}
