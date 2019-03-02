package cybsrv

import (
	"fmt"
	"log"
	"strings"
	"time"

	rep "coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "coding.net/bobxuyang/cy-gateway-BN/models"
	"coding.net/bobxuyang/cy-gateway-BN/utils"
	apim "coding.net/yundkyy/cybexgolib/api"
	"coding.net/yundkyy/cybexgolib/crypto"
	"coding.net/yundkyy/cybexgolib/types"
	"github.com/joho/godotenv"
	"github.com/juju/errors"
	"github.com/spf13/viper"
)

var api apim.BitsharesAPI
var gatewayPassword string
var gatewayAccount *types.Account
var coldAccount *types.Account
var gatewaykeyBag *crypto.KeyBag
var gatewayMemoPri types.PrivateKeys
var gatewayPrefix string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	utils.InitConfig()
	// init db
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	m.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	rep.Init()

	fmt.Println("after m.InitDB")

	node := viper.GetString("cybsrv.node")
	api = apim.New(node, "")
	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "OnConnect"))
	}
	gatewayPassword = viper.GetString("cybsrv.gatewayPassword")
	gatewayAccountStr := viper.GetString("cybsrv.gatewayAccount")
	gatewayAccount, err = api.GetAccountByName(gatewayAccountStr)
	if err != nil {
		panic(err)
	}
	gatewayPrefix = viper.GetString("cybsrv.gatewayPrefix")
	coldStr := viper.GetString("cybsrv.coldAccount")
	coldAccount, err = api.GetAccountByName(coldStr)
	if err != nil {
		panic(err)
	}

	gatewaykeyBag = apim.KeyBagByUserPass(gatewayAccountStr, gatewayPassword)
	memokey := gatewayAccount.Options.MemoKey
	pubkeys := types.PublicKeys{memokey}
	gatewayMemoPri = gatewaykeyBag.PrivatesByPublics(pubkeys)
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
	amount := order1.Amount.Text('f')
	if order1.From == "" {
		order1.From = gatewayAccount.Name
	}
	tx, err := api.Send(order1.From, order1.To, amount, asset.CybID, "", gatewayPassword)
	if err != nil {
		if strings.Contains(err.Error(), "skip_transaction_dupe_check") {
			order1.UpdateColumns(&m.CybOrder{
				Status: m.CybOrderStatusFailed,
			})
		} else {
			utils.Errorln("api.Send", err)
		}
	} else {
		utils.Infoln("sendorder tx is ", *tx)
		signed := tx.Signatures[0].String()
		order1.UpdateColumns(&m.CybOrder{
			Status: m.CybOrderStatusPending,
			Sig:    signed,
		})
	}
}

// HandleWorker ...
func HandleWorker() {
	for {
		// utils.Debugln("start...")
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
		utils.Debugf("fails=>init... %d", rownum)
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
