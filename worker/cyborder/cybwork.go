package cyborder

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/utils"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	apim "github.com/CybexDex/cybex-go/api"
	cybTypes "github.com/CybexDex/cybex-go/types"
	"github.com/spf13/viper"
)

var api apim.CybexAPI

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
func updateAllUnDone(current string) {
	res, err := model.JPOrderCurrentNotDone(current, "1m", 0, 10)
	if err != nil {
		log.Errorln("updateAllUnDone", err)
		return
	}
	for _, order := range res {
		switch order.CurrentState {
		case model.JPOrderStatusFailed:
			order.SetCurrent(current, model.JPOrderStatusInit, "fail to init")
		case model.JPOrderStatusProcessing:
			// 如果sig不存在
			if order.Sig == "" {
				order.SetCurrent(current, model.JPOrderStatusInit, "processing to init")
			}
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
		updateAllUnBalance()
		updateAllUnDone("cyborder")
		updateAllUnDone("cybinner")
		for {
			ret := HandleDepositOneTime()
			if ret != 0 {
				break
			}
		}
		// for {
		// 	ret := HandleInnerOneTime()
		// 	if ret != 0 {
		// 		break
		// 	}
		// }
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
		stx, err := mySend(tosends, order)
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
		// order.Sig = stx.Signatures[0].String()
		order.SetCurrent("cybinner", model.JPOrderStatusPending, "")
	} else {
		log.Infoln("cannot handle this action,order", order.ID)
		order.SetCurrent("cybinner", model.JPOrderStatusTerminate, "cannot handle this action")
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
	assetC, err := model.AssetsFind(order.Asset)
	if err != nil {
		return fmt.Errorf("AssetsFind %v", err)
	}
	gatewayAccount := assetC.GatewayAccount
	gatewayPassword := utils.SeedString(assetC.GatewayPass)
	//
	tosends := []cybTypes.SimpleSend{}
	tosend := cybTypes.SimpleSend{
		From:     gatewayAccount,
		To:       order.CybUser,
		Amount:   order.Amount.String(),
		Asset:    assetC.CYBName,
		Password: gatewayPassword,
		Memo:     "address:" + order.To,
	}
	tosends = append(tosends, tosend)
	order.Memo = tosend.Memo
	// log.Infoln(tosends)
	stx, err := mySend(tosends, order)
	if err != nil {
		log.Errorln("xxxx", err)
		order.SetCurrent("cyborder", model.JPOrderStatusFailed, "send error")
		return err
	}
	log.Infof("order:%d,%s:%+v\n", order.ID, "sendorder", *stx)
	order.SetCurrent("cyborder", model.JPOrderStatusPending, "")
	return nil

	// order.SetStatus(model.JPOrderStatusDone)
}

func mySend(tosends []cybTypes.SimpleSend, order *model.JPOrder) (tx *cybTypes.SignedTransaction, err error) {
	defer func() {
		if r := recover(); r != nil {
			// log.Errorf("%v, stack: %s", r, debug.Stack())
			log.Errorf("%v, stack: %s", r)
			err = fmt.Errorf("send Error")
		}
	}()
	log.Infoln(tosends)
	tx, err = api.PreSends(tosends)
	if err != nil {
		log.Errorln("updateAllUnDone", err)
		return tx, err
	}
	order.Sig = tx.Signatures[0].String()
	err = order.Save()
	if err != nil {
		log.Errorln("updateAllUnDone", err)
		return tx, err
	}
	if err := api.BroadcastTransaction(tx); err != nil {
		//log.Fatal(errors.Annotate(err, "BroadcastTransaction"))
		return nil, err
	}
	return tx, nil
}
