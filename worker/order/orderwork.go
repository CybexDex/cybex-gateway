package order

import (
	"fmt"
	"time"

	"cybex-gateway/model"
	"cybex-gateway/utils/log"
)

// HoldOne ...
func HoldOne() (*model.JPOrder, error) {
	order, err := model.HoldOrderOne()
	return order, err
}
func updateAllUnDone() {
	res, err := model.JPOrderCurrentNotDone("order", "1m", 0, 10)
	if err != nil {
		log.Errorln("updateAllUnDone", err)
		return
	}
	for _, order := range res {
		switch order.CurrentState {
		case model.JPOrderStatusFailed:
			order.SetCurrent("order", model.JPOrderStatusInit, "fail to init")
		case model.JPOrderStatusProcessing:
			order.SetCurrent("order", model.JPOrderStatusInit, "processing to init")
		}
		err = order.Save()
		if err != nil {
			log.Errorln("updateAllUnDone", err)
		}
	}
}

// HandleWorker ...
func HandleWorker(seconds int) {
	for {
		updateAllUnDone()
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
	if order1.Type == model.JPOrderTypeDeposit {
		handleOrders(order1)
	} else if order1.Type == model.JPOrderTypeWithdraw {
		handleWithdrawOrders(order1)
	}
	order1.Save()
	return 0
	// check order to process it
}
func handleWithdrawOrders(order *model.JPOrder) (err error) {
	// 是否可处理的asset
	asset, err := model.AssetsFind(order.Asset)
	if err != nil {
		err = fmt.Errorf("asset_cannot_find %v", err)
		log.Errorln(err)
		return err
	}
	// 是否可以提现
	if !asset.WithdrawSwitch {
		err = fmt.Errorf("withdrawSwitch false ")
		order.CurrentState = model.JPOrderStatusFailed
		order.CurrentReason = "withdrawSwitch"
		log.Warningln(err)
		return err
	}
	// 提现金额是否过小
	if order.TotalAmount.LessThan(asset.MinWithdraw) {
		err = fmt.Errorf("MinWithdraw false ")
		order.CurrentState = model.JPOrderStatusTerminate
		order.CurrentReason = "MinWithdraw"
		log.Errorln(err)
		return err
	}
	// 是否黑名单
	isblack, bs, err := IsBlack(order)
	if err != nil {
		err = fmt.Errorf("IsBlack err,%v", err)
		order.CurrentState = model.JPOrderStatusFailed
		order.CurrentReason = err.Error()
		log.Errorln(err)
		return err
	}
	if isblack {
		err = fmt.Errorf("Order IsBlack,%v %v", bs[0].Address, bs[0].Blockchain)
		order.CurrentState = model.JPOrderStatusTerminate
		order.CurrentReason = err.Error()
		log.Errorln(err)
		return err
	}
	// 计算费率
	fee := asset.WithdrawFee
	order.Fee = fee
	order.Amount = order.TotalAmount.Sub(order.Fee)
	// 是否大额
	// order 通过，进入下一阶段
	order.SetCurrent("jp", model.JPOrderStatusInit, "")
	return nil
}

func handleOrders(order *model.JPOrder) (err error) {
	// 是否可处理的asset
	asset, err := model.AssetsFind(order.Asset)
	if err != nil {
		err = fmt.Errorf("asset_cannot_find %s", order.Asset)
		log.Errorln(err)
		return err
	}
	// 是否可以充值
	if !asset.DepositSwitch {
		err = fmt.Errorf("depositSwitch false ")
		order.CurrentState = model.JPOrderStatusFailed
		order.CurrentReason = "depositSwitch"
		log.Warningln(err)
		return err
	}
	// 充值金额是否过小
	if order.TotalAmount.LessThan(asset.MinDeposit) {
		err = fmt.Errorf("MinDeposit false ")
		order.CurrentState = model.JPOrderStatusTerminate
		order.CurrentReason = "MinDeposit"
		log.Errorln(err)
		return err
	}
	// 是否黑名单
	isblack, bs, err := IsBlack(order)
	if err != nil {
		err = fmt.Errorf("IsBlack err,%v", err)
		order.CurrentState = model.JPOrderStatusFailed
		order.CurrentReason = err.Error()
		log.Errorln(err)
		return err
	}
	if isblack {
		err = fmt.Errorf("Order IsBlack,%v %v", bs[0].Address, bs[0].Blockchain)
		order.CurrentState = model.JPOrderStatusTerminate
		order.CurrentReason = err.Error()
		log.Errorln(err)
		return err
	}
	// 计算费率
	order.Fee = asset.DepositFee
	order.Amount = order.TotalAmount.Sub(order.Fee)
	// 是否大额
	// order 通过，进入下一阶段
	order.SetCurrent("cyborder", model.JPOrderStatusInit, "")
	return nil
}
