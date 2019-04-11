package jp

import (
	"bitbucket.org/woyoutlz/bbb-gateway/controller/jp"
	"bitbucket.org/woyoutlz/bbb-gateway/types"
	"bitbucket.org/woyoutlz/bbb-gateway/utils"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// StartServer ...
func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/api/order/noti", func(c *gin.Context) {
		// 记录日志
		reqBody := &types.JPEvent{}
		err := c.Bind(&reqBody)
		if err != nil {
			log.Errorln("Error", err)
			c.JSON(400, gin.H{
				"message": "Unmarshal error",
			})
			return
		}
		// ok, err := ecc.VerifyECCSign(reqBody.Result, reqBody.Sig, "04ace32532c90652e1bae916248e427a7ab10aeeea1067949669a3f4da10965ef90d7297f538f23006a31f94fdcfaed9e8dd38c85ba7e285f727430332925aefe5")
		err = jp.CheckComing(reqBody)
		if err != nil {
			log.Errorln("Error", err)
			c.JSON(400, gin.H{
				"message": err,
			})
			return
		}
		result := types.JPOrderResult{}
		err = utils.ResultToStruct(reqBody.Result, &result)
		if err != nil {
			log.Errorln("Error", err)
			c.JSON(400, gin.H{
				"message": "Unmarshal error",
			})
			return
		}
		if result.BizType == "WITHDRAW" {
			// 提现订单
			err = jp.HandleWithdraw(result)
			if err != nil {
				log.Errorln("Error", err)
				c.JSON(400, gin.H{
					"message": "HandleWithdraw Error",
				})
				return
			}
		} else if result.BizType == "DEPOSIT" {
			// 充值订单
			err = jp.HandleDeposit(result)
			if err != nil {
				log.Errorln("Error", err)
				c.JSON(400, gin.H{
					"message": "HandleDeposit Error",
				})
				return
			}
		} else {
			c.JSON(400, gin.H{
				"message": "BizType Error",
			})
			return
		}
		// 返回
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	port := viper.GetString("jpserver.port")
	log.Infoln("jpserver start at", port)
	err := r.Run(port) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Errorln(err)
	}
}
