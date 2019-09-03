package sass

import (
	"fmt"
	"time"

	jpc "cybex-gateway/controller/sass"
	// jpc "cybex-gateway/controller/sass"
	"cybex-gateway/model"
	"cybex-gateway/utils/log"

	"github.com/spf13/viper"
)

func updateAllUnDone() {
	// log.Infoln("jp fail => init")
	res, err := model.JPWithdrawFailed("1m", 0, 10)
	if err != nil {
		log.Errorln("updateAllUnDone", err)
		return
	}
	for _, order := range res {
		switch order.CurrentState {
		case model.JPOrderStatusFailed:
			order.SetCurrent("jp", model.JPOrderStatusInit, "fail to init")
		case model.JPOrderStatusProcessing:
			order.SetCurrent("jp", model.JPOrderStatusInit, "processing to init")
		}
		err = order.Save()
		if err != nil {
			log.Errorln("updateAllUnDone", err)
		}
	}
}

// HandleWorker ...
func HandleWorker(seconds int) {
	log.Infoln("sass worker start")
	for {
		isfail2init := viper.GetBool("sassserver.isfail2init")
		if isfail2init {
			updateAllUnDone()
		}
		for {
			ret := HandleOneTime()
			if ret != 0 {
				break
			}
		}
		time.Sleep(time.Second * time.Duration(seconds))
	}
}

// HandleOneTime ...
func HandleOneTime() int {
	order1, _ := HoldOne()
	if order1.ID == 0 {
		return 1
	}
	if order1.Type == model.JPOrderTypeWithdraw {
		handleOrders(order1)
	}
	order1.Save()
	return 0
	// check order to process it
}

// HoldOne ...
func HoldOne() (*model.JPOrder, error) {
	order, err := model.HoldJPWithdrawOne()
	return order, err
}
func handleOrders(order *model.JPOrder) error {
	log.Infof("order:%d,%s:%+v\n", order.ID, "jpwork_handle", *order)
	// 订单序列号设置
	evt := fmt.Sprintf("sequence:%d,%+v", order.ID*100+order.BNRetry, *order)
	order.Log("before_BNWithdraw", evt)
	assetobj, err := model.AssetsFind(order.Asset)
	if err != nil {
		errmsg := fmt.Sprintf("%s %s", order.Asset, "不是合法币种")
		log.Infoln(errmsg)
		return err
	}
	assetName := assetobj.Name
	if assetobj.JadeName != "" {
		assetName = assetobj.JadeName
	}
	result, err := jpc.Withdraw(assetName, order.OutAddr, order.Amount.String(), order.ID*100+order.BNRetry)
	if err != nil {
		errstr := fmt.Sprintf("jpc.Withdraw:%v", err)
		log.Errorf("order:%d,%s:%+v\n", order.ID, "jpc.Withdraw", err)
		order.SetCurrent(order.Current, model.JPOrderStatusFailed, errstr)
		errmsg := fmt.Sprintf("id:%d\nerr:%s", order.ID, errstr)
		model.WxSendTaskCreate("网关提现失败", errmsg)
		return err
	}
	evt2 := fmt.Sprintf("sequence:%d,%+v", order.ID*100+order.BNRetry, *result)
	order.Log("after_BNWithdraw", evt2)
	order.BNOrderID = &result.ID
	order.Current = "jpsended"
	order.CurrentState = result.State
	order.Confirmations = result.Confirmations
	err = order.Save()
	if err != nil {
		log.Errorln("order.Save", err)
		return err
	}
	return nil
}
