package jp

import (
	"fmt"
	"time"

	jpc "cybex-gateway/controller/sass"
	// jpc "cybex-gateway/controller/sass"
	"cybex-gateway/model"
	"cybex-gateway/utils/log"
)

// HandleWorker ...
func HandleWorker(seconds int) {
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
	result, err := jpc.Withdraw(order.Asset, order.OutAddr, order.Amount.String(), order.ID*100+order.BNRetry)
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
