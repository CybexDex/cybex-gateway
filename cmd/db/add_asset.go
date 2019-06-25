package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shopspring/decimal"

	"cybex-gateway/config"
	"cybex-gateway/controller/jp"
	"cybex-gateway/model"

	apim "github.com/CybexDex/cybex-go/api"
	cybTypes "github.com/CybexDex/cybex-go/types"
	"github.com/spf13/viper"
)

var coinviper *viper.Viper
var api apim.CybexAPI

func main() {
	env := os.Getenv("env")
	if len(env) == 0 {
		env = "dev"
	}
	config.LoadConfig(env)
	model.INITFromViper()
	coinviper = viper.New()
	coin := os.Getenv("coin")
	fmt.Println(coin)
	coinviper.SetConfigName(coin)
	coinviper.AddConfigPath("./cmd/db/coins")
	err := coinviper.ReadInConfig() // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	node := viper.GetString("cybserver.node")
	api = apim.New(node, "")
	if err := api.Connect(); err != nil {
		panic(err)
	}
	asset()
}
func asset() {
	var db = model.GetDB()
	name := coinviper.GetString("name")
	fmt.Println(name)
	md, _ := decimal.NewFromString(coinviper.GetString("minDeposit"))
	mw, _ := decimal.NewFromString(coinviper.GetString("minWithdraw"))
	df, _ := decimal.NewFromString(coinviper.GetString("depositFee"))
	wf, _ := decimal.NewFromString(coinviper.GetString("withdrawFee"))
	// 瑶池是否支持
	coinName := coinviper.GetString("name")
	_, err := jp.VerifyAddress(coinName, "fortestaddress")
	if strings.Contains(err.Error(), "不支持该币种类型") {
		fmt.Println("瑶池不支持币种", err)
		return
	}
	// 这个cybex 网关账号能否转账1CYB成功
	gatewayAccount := coinviper.GetString("gatewayAccount")
	gatewayPass := coinviper.GetString("gatewayPass")
	tosends := []cybTypes.SimpleSend{}
	tosend := cybTypes.SimpleSend{
		From:     gatewayAccount,
		To:       "a-new-world",
		Amount:   "1",
		Asset:    "CYB",
		Password: gatewayPass,
		Memo:     "address:1",
	}
	tosends = append(tosends, tosend)
	tx, err := api.PreSends(tosends)
	err = api.ValidateTransaction(tx)
	if err != nil {
		fmt.Println("cybex测试失败", err)
		return
	}
	b := model.Asset{
		Name:           coinName,
		Blockchain:     coinviper.GetString("blockchain"),
		CYBName:        coinviper.GetString("cybname"),
		CYBID:          coinviper.GetString("cybid"),
		SmartContract:  coinviper.GetString("smartContract"),
		GatewayAccount: gatewayAccount,
		GatewayPass:    gatewayPass,
		WithdrawPrefix: coinviper.GetString("withdrawPrefix"),

		DepositSwitch:  coinviper.GetBool("depositSwitch"),
		WithdrawSwitch: coinviper.GetBool("withdrawSwitch"),

		MinDeposit:  md,
		MinWithdraw: mw,
		WithdrawFee: wf,
		DepositFee:  df,

		Precision: coinviper.GetString("precision"),
		ImgURL:    coinviper.GetString("imgURL"),
		HashLink:  coinviper.GetString("hashLink"),
	}
	_, err = model.AssetsFind(name)
	if err != nil {
		err = db.Create(&b).Error
		fmt.Println("新记录", err)
	} else {
		fmt.Println(err)
		err = db.Debug().Model(model.Asset{}).Where(&model.Asset{
			Name: name,
		}).UpdateColumn(&b).Error
		fmt.Println("更新记录", err)
	}
}
