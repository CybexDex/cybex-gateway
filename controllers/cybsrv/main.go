package cybsrv

import (
	"log"
	"strings"
	"time"

	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
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
	return rep.CybOrder.HoldingOne(), nil
}
func handleOrders(order1 *m.CybOrder) {
	var err error
	asset, err := rep.Asset.GetByID(order1.AssetID)
	if err != nil {
		order1.UpdateColumns(&m.CybOrder{
			Status: m.CybOrderStatusFailed,
		})
		return
	}
	amount, _ := order1.Amount.Float64()
	_, err = api.Send(order1.From, order1.To, amount, asset.CybID, "", gatewayPassword)
	if err != nil {
		if strings.Contains(err.Error(), "skip_transaction_dupe_check") {
			order1.UpdateColumns(&m.CybOrder{
				Status: m.CybOrderStatusFailed,
			})
		}
	} else {
		order1.UpdateColumns(&m.CybOrder{
			Status: m.CybOrderStatusPending,
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
		re := rep.CybOrder.UpdateAll(&m.CybOrder{Status: m.CybOrderStatusFailed}, &m.CybOrder{
			Status: m.CybOrderStatusInit,
		})
		rownum := re.RowsAffected
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
