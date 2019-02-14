package ordersrv

import (
	"runtime/debug"
	"time"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	apim "git.coding.net/yundkyy/cybexgolib/api"
)

var api apim.BitsharesAPI
var gatewayPassword string

func init() {

}
func findOrders() (*m.Order, error) {
	db := m.GetDB()
	// do some database operations in the transaction (use 'tx' from this point, not 'db')
	var order1 m.Order

	// time.Sleep(time.Second * 2)
	// fmt.Println("ID", order1.ID)
	s := `update orders 
	set status = 'PROCESSING' 
	where id = (
				select id 
				from orders 
				where status = 'INIT' 
				order by id
				limit 1
			)
	returning *`
	db.Raw(s).Scan(&order1)
	// ...
	// Or commit the transaction
	return &order1, nil
}
func handleOrders(order1 *m.Order) {
	// fmt.Println("order1", order1)
	// fmt.Println("send order id", order1.ID, order1.From, order1.To, order1.Amount, order1.AssetID)
	db := m.GetDB()
	asset := &m.Asset{}
	db.First(asset, order1.AssetID)
	tx := m.GetDB().Begin()
	defer func() {
		tx.Save(order1)
		tx.Commit()
		if r := recover(); r != nil {
			utils.Errorf("%v, stack: %s", r, debug.Stack())
			tx.Rollback()
		}
	}()
	isopen, _ := IsOpen(order1)

	if !isopen {
		order1.Status = "TERMINATED"
		return
	}
	isblack, _ := IsBlack(order1)
	if isblack {
		order1.Status = "TERMINATED"
		return
	}
	isbig, _ := IsBig(order1)
	if isbig {
		order1.Status = "WAITING"
		return
	}
	order1.Status = "DONE"
	order1.CreateNext(tx)
}

// HandleWorker ...
func HandleWorker() {
	for {
		for {
			ret := HandleOneTime()
			if ret != 0 {
				break
			}
		}
		db := m.GetDB()
		rownum := db.Exec("update orders set status='INIT' where status='FAIL'").RowsAffected
		utils.Infof("fails=>init... %d", rownum)
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
