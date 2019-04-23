package user

import (
	"bitbucket.org/woyoutlz/bbb-gateway/utils/ecc"

	userc "bitbucket.org/woyoutlz/bbb-gateway/controller/user"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// StartServer ...
func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/t", func(c *gin.Context) {
		ecc.TestECCSign()
		c.JSON(200, gin.H{})
	})
	r.GET("/v1/users/:user/assets/:asset/address", func(c *gin.Context) {
		user := c.Param("user")
		asset := c.Param("asset")
		log.Infoln("GetAddress", user, asset)
		address, err := userc.GetAddress(user, asset)
		if err != nil {
			log.Errorln("user address", err)
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, address)
	})
	r.GET("/v1/bbb", func(c *gin.Context) {
		address, err := userc.GetBBBAssets()
		if err != nil {
			log.Errorln("user address", err)
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, address)
	})
	port := viper.GetString("userserver.port")
	log.Infoln("userserver start at", port)
	r.Run(port) // listen and serve on 0.0.0.0:8080
}
