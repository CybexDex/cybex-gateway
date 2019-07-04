package wx

import (
	"cybex-gateway/model"
	"cybex-gateway/utils/log"
	wxu "cybex-gateway/utils/wx"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// HoldOne ...
func HoldOne() (*model.Task, error) {
	order, err := model.HoldWxOne()
	return order, err
}

// HandleOneTime ...
func HandleOneTime() int {
	order1, _ := HoldOne()
	if order1.ID == 0 {
		return 1
	}
	// fmt.Println(order1)
	prefix := viper.GetString("wx.prefix")
	msg := order1.Value
	if prefix != "" {
		msg = fmt.Sprintf("%s_%s", prefix, order1.Value)
	}
	err := wxu.WeixinSend(msg)
	if err != nil {
		order1.Status = model.JPOrderStatusFailed
		order1.Adds1 = fmt.Sprintf("fail:%v", err)
		log.Errorln("wx send Fail", order1.ID)
	} else {
		order1.Status = model.JPOrderStatusDone
		log.Infoln("wx send ok", order1.ID)
	}
	order1.Save()
	return 0
	// check order to process it
}

// InitWx ...
func InitWx() {
	corpid := viper.GetString("wx.corpid")
	corpsecret := viper.GetString("wx.corpsecret")
	agentid := viper.GetInt("wx.agentid")
	users := viper.GetString("wx.users")
	wxu.InitWeixin(corpid, corpsecret, agentid, users)
	wxu.GetToken()
	// log.Infoln("token:", token, "err", err)
}

// HandleWorker ...
func HandleWorker(seconds int) {
	log.Infoln("wx start...")
	InitWx()
	for {
		for {
			ret := HandleOneTime()
			if ret != 0 {
				break
			}
		}
		time.Sleep(time.Second * time.Duration(seconds))
	}
}
