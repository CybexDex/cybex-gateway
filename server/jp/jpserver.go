package jp

import (
	"bytes"
	"encoding/json"

	"bitbucket.org/woyoutlz/bbb-gateway/controller/jp"
	"bitbucket.org/woyoutlz/bbb-gateway/types"
	"bitbucket.org/woyoutlz/bbb-gateway/utils"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/eventlog"
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
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		str := buf.String()
		eventlog.Log("jpnoti", str)
		// 充提记录的进一步处理
		reqBody := new(types.JPEvent)
		err := json.Unmarshal([]byte(str), reqBody)
		if err != nil {
			log.Errorln("Error", err)
			c.JSON(400, gin.H{
				"message": "Unmarshal error",
			})
		}
		result := types.JPOrderResult{}
		err = utils.ResultToStruct(reqBody.Result, &result)
		if err != nil {
			log.Errorln("Error", err)
			c.JSON(400, gin.H{
				"message": "Unmarshal error",
			})
		}
		if result.BizType == "WITHDRAW" {
			// 提现订单
			err = jp.HandleWithdraw(result)
			if err != nil {
				log.Errorln("Error", err)
				c.JSON(400, gin.H{
					"message": "HandleWithdraw Error",
				})
			}
		} else if result.BizType == "DEPOSIT" {
			// 充值订单
			err = jp.HandleDeposit(result)
			if err != nil {
				log.Errorln("Error", err)
				c.JSON(400, gin.H{
					"message": "HandleDeposit Error",
				})
			}
		} else {
			c.JSON(400, gin.H{
				"message": "BizType Error",
			})
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
