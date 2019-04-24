package jp

import (
	"time"

	jpc "bitbucket.org/woyoutlz/bbb-gateway/controller/jp"
	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
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
	result, err := jpc.Withdraw(order.Asset, order.OutAddr, order.Amount.String(), order.ID)
	if err != nil {
		log.Errorln("jpc.Withdraw", err)
		return err
	}
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
