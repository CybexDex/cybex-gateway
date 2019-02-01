package cybsrv

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	apim "git.coding.net/yundkyy/cybexgolib/api"
	"github.com/joho/godotenv"
	"github.com/juju/errors"
)

var api apim.BitsharesAPI
var gatewayPassword string

func init() {
	api = apim.New("wss://shanghai.51nebula.com/", "")
	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "OnConnect"))
	}
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	gatewayPassword = os.Getenv("gatewayPassword")
}
func findOrders() (*m.CybOrder, error) {
	db := m.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// do some database operations in the transaction (use 'tx' from this point, not 'db')
	var order1 m.CybOrder
	tx.Raw("select * from cyb_orders  where status='INIT' limit 1  for update").Scan(&order1)

	// time.Sleep(time.Second * 1)
	// fmt.Println("ID", order1.ID)
	tx.Exec("update cyb_orders set status='HOLDING' where id=?", order1.ID).Scan(&order1)
	// ...
	// Or commit the transaction
	tx.Commit()
	return &order1, nil
}
func handleOrders(order1 *m.CybOrder) {
	// fmt.Println("order1", order1)
	fmt.Println("send cyborder id", order1.ID, order1.From, order1.To, order1.Amount, order1.AssetID)
	db := m.GetDB()
	asset := &m.Asset{}
	db.First(asset, order1.AssetID)
	amount, _ := order1.Amount.Float64()
	re, err := api.Send(order1.From, order1.To, amount, asset.CybID, "", gatewayPassword)
	if err != nil {
		fmt.Println(1, err)
		if strings.Contains(err.Error(), "skip_transaction_dupe_check") {
			order1.UpdateColumns(&m.CybOrder{
				Status: "FAIL",
			})
		}
	} else {
		fmt.Println(re)
		order1.UpdateColumns(&m.CybOrder{
			Status: "PENDING",
		})
	}
}

// HandleWorker ...
func HandleWorker() {
	for {
		fmt.Println("start handle...")
		for {
			ret := HandleOneTime()
			if ret != 0 {
				break
			}
		}
		db := m.GetDB()
		rownum := db.Exec("update cyb_orders set status='INIT' where status='FAIL'").RowsAffected
		fmt.Println("fails=>init", rownum, "waiting next...", 10)
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
func WorkerStart() {

}
func WorkerStop() {

}
