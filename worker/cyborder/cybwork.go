package cyborder

import (
	"fmt"
	"strings"
	"time"

	"cybex-gateway/model"
	"cybex-gateway/utils"
	"cybex-gateway/utils/log"

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

// UpdateExpire ...
func UpdateExpire() {
	log.Infoln("UpdateExpire...")
	current, err := model.EasyFrist("cybLastBlockNum")
	if err != nil {
		log.Errorln("updateExpire", err)
		return
	}
	res, err := model.CYBOrderExpire(current.RecordTime, "1m", 0, 10)
	if err != nil {
		log.Errorln("updateExpire", err)
		return
	}
	for _, order := range res {
		switch order.CurrentState {
		case model.JPOrderStatusPending:
			order.SetCurrent("cyborder", model.JPOrderStatusFailed, "expire Pending to fail")
		case model.JPOrderStatusProcessing:
			// 如果sig不存在,有sig了才会发
			order.SetCurrent("cyborder", model.JPOrderStatusFailed, "expire processing to fail")
		}
		err = order.Save()
		if err != nil {
			log.Errorln("UpdateExpire save", err)
		}
	}
}

// HoldInnerOne ...
// func HoldInnerOne() (*model.JPOrder, error) {
// 	order, err := model.HoldCYBInnerOrderOne()
// 	return order, err
// }
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
			// 如果sig不存在,有sig了才会发
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
		// updateAllUnBalance()
		updateAllUnDone("cyborder")
		isauthFail := viper.GetBool("cybserver.expireAuthFail")
		if isauthFail {
			UpdateExpire()
		}
		// updateAllUnDone("cybinner")
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

// func updateAllUnBalance() {
// 	model.JPOrderUnBalanceInit()
// }

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

// func setOrderThen(order1 *model.JPOrder, order2 *model.JPOrder, timeout int) {
// 	time.Sleep(time.Second * time.Duration(timeout))
// 	order1.Update(order2)
// }
func handleOrders(order *model.JPOrder) (err error) {
	// 是否可处理的asset
	log.Infoln("start handle order", order.ID)
	assetC, err := model.AssetsFind(order.Asset)
	if err != nil {
		log.Errorln("AssetsFind", err)
		return fmt.Errorf("AssetsFind %v", err)
	}
	gatewayAccount := assetC.GatewayAccount
	// gatewayPassword := utils.SeedString(assetC.GatewayPass)
	keybag := utils.KeyBagByUserSeedPass(gatewayAccount, assetC.GatewayPass)
	prikeys := keybag.Privates()
	priWifs := []string{}
	for _, prikey := range prikeys {
		priWifs = append(priWifs, prikey.ToWIF())
	}
	prikeysStr := strings.Join(priWifs, ",")
	//
	tosends := []cybTypes.SimpleSend{}
	if gatewayAccount == order.CybUser {
		order.SetCurrent("cyborder", model.JPOrderStatusTerminate, "网关账号不能发给自己")
		return nil
	}
	tosend := cybTypes.SimpleSend{
		From:     gatewayAccount,
		To:       order.CybUser,
		Amount:   order.Amount.String(),
		Asset:    assetC.CYBName,
		Password: "," + prikeysStr,
		Memo:     "address:" + order.To,
	}
	tosends = append(tosends, tosend)
	order.Memo = tosend.Memo
	stx, err := mySend(tosends, order)
	if err != nil {
		log.Errorln("xxxx", err)
		order.SetCurrent("cyborder", model.JPOrderStatusFailed, "send error")
		errmsg := fmt.Sprintf("id:%d\nerr:%v", order.ID, err)
		model.WxSendTaskCreate("cybex充值失败", errmsg)
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
	// log.Infoln(tosends)
	tx, err = api.PreSends(tosends)
	if err != nil {
		log.Errorln("updateAllUnDone", err)
		return tx, err
	}
	order.ExpireTime = tx.Transaction.Expiration.Time
	order.Sig = tx.Signatures[0].String()
	err = order.Save()
	if err != nil {
		log.Errorln("order send before", err)
		return tx, err
	}
	if err := api.BroadcastTransaction(tx); err != nil {
		//log.Fatal(errors.Annotate(err, "BroadcastTransaction"))
		return nil, err
	}
	return tx, nil
}
