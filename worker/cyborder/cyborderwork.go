package cyborder

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	"coding.net/bobxuyang/cy-gateway-BN/utils"
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

func handleOrders(order *model.JPOrder) (err error) {
	// 是否可处理的asset
	assets := viper.GetStringMap("assets")
	orderAsset := strings.ToLower(order.Asset)
	if assets[orderAsset] == nil {
		err = fmt.Errorf("asset_cannot_find %s", order.Asset)
		log.Errorln(err)
		return err
	}
	actionPath := fmt.Sprintf("assets.%s.handle_action", orderAsset)
	action := viper.GetString(actionPath)
	accountPath := fmt.Sprintf("assets.%s.deposit.gateway", orderAsset)
	gatewayAccount := viper.GetString(accountPath)
	accountPass := fmt.Sprintf("assets.%s.deposit.gatewaypass", orderAsset)
	gatewayPassword := viper.GetString(accountPass)
	sendTO := viper.GetStringSlice(fmt.Sprintf("assets.%s.deposit.sendto", orderAsset))
	if action == "BBB" {
		log.Infoln(gatewayAccount, gatewayPassword, sendTO)
		//
		tosends := []cybTypes.SimpleSend{}
		for _, dasset := range sendTO {
			ds := strings.Split(dasset, ":")
			utils.Infoln("ds", ds)
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
		fmt.Println(tosends)
		stx, err := mySend(tosends)
		fmt.Println("xxxx", err)
		if err != nil {
			order.SetCurrent("cyborder", model.JPOrderStatusFailed, "send error")
			return err
		}
		utils.Infoln("sendorder tx is ", *stx)
	} else {
		log.Infoln("cannot handle this action,order", order.ID)
		order.SetCurrent("cyborder", model.JPOrderStatusFailed, "cannot handle this action")
		return nil
	}
	order.SetCurrent("done", model.JPOrderStatusDone, "")
	order.SetStatus(model.JPOrderStatusDone)
	return nil
}

func mySend(tosends []cybTypes.SimpleSend) (tx *cybTypes.SignedTransaction, err error) {
	defer func() {
		if r := recover(); r != nil {
			// utils.Errorf("%v, stack: %s", r, debug.Stack())
			utils.Errorf("%v, stack: %s", r)
			err = fmt.Errorf("send Error")
		}
	}()
	tx, err = api.Sends(tosends)
	return tx, err
}
