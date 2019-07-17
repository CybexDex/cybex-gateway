package jp

import (
	"cybex-gateway/controller/jp"
	"cybex-gateway/server/middleware"
	"cybex-gateway/types"
	"cybex-gateway/utils"
	"cybex-gateway/utils/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func jpNotify(c *gin.Context) {
	// 记录日志
	reqBody := &types.JPEvent{}
	err := c.Bind(&reqBody)
	if err != nil {
		log.Errorln("Error", err)
		errorRes(c, 400, gin.H{
			"message": "Unmarshal error",
		})
		return
	}
	err = jp.CheckComing(reqBody)
	if err != nil {
		log.Errorln("Error CheckComing", err)
		errorRes(c, 400, gin.H{
			"message": err,
		})
		return
	}
	result := types.JPOrderResult{}
	err = utils.ResultToStruct(reqBody.Result, &result)
	if err != nil {
		log.Errorln("Error", err)
		errorRes(c, 400, gin.H{
			"message": "Unmarshal Error",
		})
		return
	}
	if result.BizType == "WITHDRAW" {
		// 提现订单
		err = jp.HandleWithdraw(result)
		if err != nil {
			log.Errorln("Error HandleWithdraw", result.ID, err)
			errorRes(c, 400, gin.H{
				"message": "HandleWithdraw Error",
			})
			return
		}
	} else if result.BizType == "DEPOSIT" {
		// 充值订单
		err = jp.HandleDeposit(result)
		if err != nil {
			log.Errorln("Error HandleDeposit", result.ID, err)
			errorRes(c, 400, gin.H{
				"message": "HandleDeposit Error",
			})
			return
		}
	} else {
		// TODO 可能一直发无法处理的类型
		errorRes(c, 400, gin.H{
			"message": "BizType Error",
		})
		return
	}
	// 返回
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// StartServer ...
func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.RequestLogger())
	r.Use(middleware.GinBodyLogMiddleware)
	r.POST("/api/order/noti", jpNotify)
	port := viper.GetString("jpserver.port")
	log.Infoln("jpserver start at", port)
	err := r.Run(port) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Errorln(err)
	}
}

func errorRes(c *gin.Context, code int, obj interface{}) {
	c.JSON(code, obj)
}
