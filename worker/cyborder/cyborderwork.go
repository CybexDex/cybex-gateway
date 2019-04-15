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
		log.Infoln(gatewayAccount, gatewayPassword, sendTO)
		//
		tosends := []cybTypes.SimpleSend{}
		for _, dasset := range sendTO {
			ds := strings.Split(dasset, ":")
			log.Infoln("ds", ds)
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
		log.Errorln("xxxx", err)
		if err != nil {
			order.SetCurrent("cyborder", model.JPOrderStatusFailed, "send error")
			return err
		}
		log.Infoln("sendorder tx is ", *stx)
	} else {
		log.Infoln("cannot handle this action,order", order.ID)
		order.SetCurrent("cyborder", model.JPOrderStatusFailed, "cannot handle this action")
		return nil
	}
	order.SetCurrent("cyborder", model.JPOrderStatusPending, "")
	// order.SetStatus(model.JPOrderStatusDone)
	return nil
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
