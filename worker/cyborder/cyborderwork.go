package cyborder

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	apim "coding.net/yundkyy/cybexgolib/api"
	cybTypes "coding.net/yundkyy/cybexgolib/types"
	"github.com/spf13/viper"
)

var api apim.BitsharesAPI

// InitNode ...
func InitNode() {
	node := viper.GetString("cybserver.node")
	api = apim.New(node, "")
	if err := api.Connect(); err != nil {
		panic(err)
	}
}

// HoldOne ...
func HoldOne() (*model.JPOrder, error) {
	order, err := model.HoldCYBOrderOne()
	return order, err
}

// HoldInnerOne ...
func HoldInnerOne() (*model.JPOrder, error) {
	order, err := model.HoldCYBInnerOrderOne()
	return order, err
}

// HandleWorker ...
func HandleWorker(seconds int) {
	for {
		updateAllUnBalance()
		for {
			ret := HandleDepositOneTime()
			if ret != 0 {
				break
			}
		}
		for {
			ret := HandleInnerOneTime()
			if ret != 0 {
				break
			}
		}
		time.Sleep(time.Second * time.Duration(seconds))
	}
}
func updateAllUnBalance() {
	model.JPOrderUnBalanceInit()
}

// HandleInnerOneTime ...
func HandleInnerOneTime() int {
	order1, _ := HoldInnerOne()
	if order1.ID == 0 {
		return 1
	}
	if order1.Current == "cybinner" {
		handleInnerOrders(order1)
		order1.Save()
	}
	return 0
	// check order to process it
}

// HandleDepositOneTime ...
func HandleDepositOneTime() int {
	order1, _ := HoldOne()
	if order1.ID == 0 {
		return 1
	}
	if order1.Current == "cyborder" {
		handleOrders(order1)
	}
	order1.Save()
	return 0
	// check order to process it
}
func findAsset(name string) (out interface{}, err error) {
	// 是否可处理的asset
	assets := viper.GetStringMap("assets")
	orderAsset := strings.ToLower(name)
	if assets[orderAsset] == nil {
		err = fmt.Errorf("asset_cannot_find %s", name)
		return nil, err
	}
	out = assets[orderAsset]
	return out, nil
}
func handleInnerOrders(order *model.JPOrder) (err error) {
	assetC := allAssets[order.Asset]
	if assetC == nil {
		err = fmt.Errorf("asset_cannot_find %s", order.Asset)
		return err
	}
	action := assetC.HandleAction
	gatewayAccount := assetC.Withdraw.Gateway
	gatewayPassword := assetC.Withdraw.Gatewaypass
	waitCoin := assetC.Withdraw.Wait
	SendTo := assetC.Withdraw.Send
	if action == "BBB" {
		// log.Infoln(gatewayAccount, gatewayPassword, waitCoin, SendTo)
		// 构造两个send to send
		tosends := []cybTypes.SimpleSend{}
		tosend1 := cybTypes.SimpleSend{
			From:     gatewayAccount,
			To:       SendTo,
			Amount:   order.TotalAmount.String(),
			Asset:    assetC.Withdraw.Coin,
			Password: gatewayPassword,
		}
		tosend2 := cybTypes.SimpleSend{
			From:     gatewayAccount,
			To:       SendTo,
			Amount:   order.TotalAmount.String(),
			Asset:    waitCoin,
			Password: gatewayPassword,
		}
		tosends = append(tosends, tosend1, tosend2)
		// log.Infoln(tosends)
		stx, err := mySend(tosends)
		if err != nil {
			if strings.Contains(err.Error(), "insufficient_balance") {
				// 金额不够,等待去
				log.Errorln("insufficient_balance", err)
				order.SetCurrent("cybinner", model.JPOrderStatusUnbalance, err.Error())
				return nil
			}
			log.Errorln("xxxx", err)
			order.SetCurrent("cybinner", model.JPOrderStatusFailed, err.Error())
			return err
		}
		// log.Infoln("sendorder tx is ", *stx)
		log.Infof("order:%d,%s:%+v\n", order.ID, "sendInnerOrder", *stx)
		order.Sig = stx.Signatures[0].String()
		order.SetCurrent("cybinner", model.JPOrderStatusPending, "")
	} else {
		log.Infoln("cannot handle this action,order", order.ID)
		order.SetCurrent("cybinner", model.JPOrderStatusFailed, "cannot handle this action")
		return nil
	}
	return nil
}

// func setOrderThen(order1 *model.JPOrder, order2 *model.JPOrder, timeout int) {
// 	time.Sleep(time.Second * time.Duration(timeout))
// 	order1.Update(order2)
// }
func handleOrders(order *model.JPOrder) (err error) {
	// 是否可处理的asset
	assetC := allAssets[order.Asset]
	if assetC == nil {
		err = fmt.Errorf("asset_cannot_find %s", order.Asset)
		return err
	}
	action := assetC.HandleAction
	gatewayAccount := assetC.Deposit.Gateway
	gatewayPassword := assetC.Deposit.Gatewaypass
	sendTO := assetC.Deposit.Sendto
	if action == "BBB" {
		// log.Infoln(gatewayAccount, sendTO)
		//
		tosends := []cybTypes.SimpleSend{}
		for _, dasset := range sendTO {
			ds := strings.Split(dasset, ":")
			// log.Infoln("ds", ds)
			assetname := ds[0]
			to := ""
			if len(ds) > 1 {
				to = ds[1]
			}
			value := ""
			if len(ds) > 2 {
				value = ds[2]
			}
			tosend := cybTypes.SimpleSend{
				From:     gatewayAccount,
				To:       order.CybUser,
				Amount:   order.Amount.String(),
				Asset:    assetname,
				Password: gatewayPassword,
			}
			if to != "" {
				tosend.To = to
			}
			if value != "" {
				tosend.Amount = value
			}
			tosends = append(tosends, tosend)
		}
		// log.Infoln(tosends)
		stx, err := mySend(tosends)
		if err != nil {
			log.Errorln("xxxx", err)
			order.SetCurrent("cyborder", model.JPOrderStatusFailed, "send error")
			return err
		}
		log.Infof("order:%d,%s:%+v\n", order.ID, "sendorder", *stx)
		// log.Infoln("sendorder tx is ", *stx)
		order.Sig = stx.Signatures[0].String()
		order.SetCurrent("cyborder", model.JPOrderStatusPending, "")
		return nil
	}
	log.Infoln("cannot handle this action,order", order.ID)
	order.SetCurrent("cyborder", model.JPOrderStatusFailed, "cannot handle this action")
	return nil

	// order.SetStatus(model.JPOrderStatusDone)
}

func mySend(tosends []cybTypes.SimpleSend) (tx *cybTypes.SignedTransaction, err error) {
	defer func() {
		if r := recover(); r != nil {
			// log.Errorf("%v, stack: %s", r, debug.Stack())
			log.Errorf("%v, stack: %s", r)
			err = fmt.Errorf("send Error")
		}
	}()
	tx, err = api.Sends(tosends)
	return tx, err
}
