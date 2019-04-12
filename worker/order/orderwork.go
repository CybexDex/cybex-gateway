package order

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"

	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
)

// HoldOne ...
func HoldOne() (*model.JPOrder, error) {
	order, err := model.HoldOrderOne()
	return order, err
}

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
	handleOrders(order1)
	order1.Save()
	return 0
	// check order to process it
}

func handleOrders(order *model.JPOrder) (err error) {
	// 是否可处理的asset
	assets := viper.GetStringMap("assets")
	orderAsset := strings.ToLower(order.Asset)
	if assets[orderAsset] == nil {
		err = fmt.Errorf("asset_cannot_find %s", order.Asset)
		log.Errorln(err)
		return err
	}
	switchPath := fmt.Sprintf("assets.%s.deposit.switch", orderAsset)
	depositSwitch := viper.GetBool(switchPath)
	// 是否可以充值
	if !depositSwitch {
		err = fmt.Errorf("depositSwitch false ")
		order.CurrentState = model.JPOrderStatusFailed
		order.CurrentReason = "depositSwitch"
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
	fee, err := decimal.NewFromString("0.0")
	if err != nil {
		return err
	}
	order.Fee = fee
	order.Amount = order.TotalAmount.Sub(order.Fee)
	// 是否大额
	// order 通过，进入下一阶段
	order.SetCurrent("cyborder", model.JPOrderStatusInit, "")
	return nil
}
