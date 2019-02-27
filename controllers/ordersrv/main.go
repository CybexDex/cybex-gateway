package ordersrv

import (
	"runtime/debug"
	"time"

	apim "coding.net/yundkyy/cybexgolib/api"
	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

var api apim.BitsharesAPI
var gatewayPassword string

func init() {

}
func findOrders() (*m.Order, error) {
	return rep.Order.HoldingOne(), nil
}
func handleOrders(order1 *m.Order) {
	_, err := rep.Asset.GetByID(order1.AssetID)
	if err != nil {
		utils.Errorln("handleOrders asset", err)
		order1.UpdateColumns(&m.Order{
			Status: m.OrderStatusFailed,
		})
		return
	}
	tx := m.GetDB().Begin()
	defer func() {
		tx.Save(order1)
		tx.Commit()
		if r := recover(); r != nil {
			utils.Errorf("%v, stack: %s", r, debug.Stack())
			tx.Rollback()
		}
	}()
	isopen, err := IsOpen(order1)
	if err != nil {
		utils.Errorln("handleOrders IsOpen", err)
		order1.Status = m.OrderStatusFailed
		return
	}
	if !isopen {
		order1.Status = m.OrderStatusTerminated
		return
	}
	isblack, err := IsBlack(order1)
	if err != nil {
		utils.Errorln("handleOrders isblack", err)
		order1.Status = m.OrderStatusFailed
		return
	}
	if isblack {
		utils.Warningln("handleOrders isblack", order1.ID)
		order1.Status = m.OrderStatusTerminated
		return
	}
	isbig, err := IsBig(order1)
	if err != nil {
		utils.Errorln("handleOrders isbig", err)
		order1.Status = m.OrderStatusFailed
		return
	}
	if isbig {
		utils.Warningln("handleOrders isbig", order1.ID)
		order1.Status = m.OrderStatusWaiting
		return
	}
	order1.Status = m.OrderStatusDone
	order1.CreateNext(tx)
}

// HandleWorker ...
func HandleWorker() {
	for {
		utils.Infoln("start...")
		for {
			ret := HandleOneTime()
			if ret != 0 {
				break
			}
		}
		time.Sleep(time.Second * 10)
	}
}

// HandleOneTime ...
func HandleOneTime() int {
	order1, _ := findOrders()
	if order1.ID == 0 {
		return 1
	}
	handleOrders(order1)
	// check order to process it
	return 0
}
