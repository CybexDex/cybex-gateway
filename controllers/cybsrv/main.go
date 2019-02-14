package cybsrv

import (
	"log"
	"strings"
	"time"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	apim "git.coding.net/yundkyy/cybexgolib/api"
	"github.com/juju/errors"
	"github.com/spf13/viper"
)

var api apim.BitsharesAPI
var gatewayPassword string

func init() {
	utils.InitConfig()
	node := viper.GetString("cybsrv.node")
	api = apim.New(node, "")
	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "OnConnect"))
	}
	gatewayPassword = viper.GetString("cybsrv.gatewayPassword")
}
func findOrders() (*m.CybOrder, error) {
	db := m.GetDB()
	// do some database operations in the transaction (use 'tx' from this point, not 'db')
	var order1 m.CybOrder

	// time.Sleep(time.Second * 2)
	// fmt.Println("ID", order1.ID)
	s := `update cyb_orders 
	set status = 'HOLDING' 
	where id = (
				select id 
				from cyb_orders 
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
func handleOrders(order1 *m.CybOrder) {
	// fmt.Println("order1", order1)
	// utils.Infof("order noti request:\n %s", requestBody)
	// fmt.Println("send cyborder id", order1.ID, order1.From, order1.To, order1.Amount, order1.AssetID)
	db := m.GetDB()
	asset := &m.Asset{}
	db.First(asset, order1.AssetID)
	amount, _ := order1.Amount.Float64()
	_, err := api.Send(order1.From, order1.To, amount, asset.CybID, "", gatewayPassword)
	if err != nil {
		if strings.Contains(err.Error(), "skip_transaction_dupe_check") {
			order1.UpdateColumns(&m.CybOrder{
				Status: "FAIL",
			})
		}
	} else {
		order1.UpdateColumns(&m.CybOrder{
			Status: "PENDING",
		})
	}
}

// HandleWorker ...
func HandleWorker() {
	for {
		utils.Infof("start...")
		for {
			ret := HandleOneTime()
			if ret != 0 {
				break
			}
		}
		db := m.GetDB()
		rownum := db.Exec("update cyb_orders set status='INIT' where status='FAIL'").RowsAffected
		// fmt.Println("fails=>init", rownum, "waiting next...", 10)
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
