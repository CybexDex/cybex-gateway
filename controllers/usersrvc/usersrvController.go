package usersrvc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	u "git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// OpMsg ...
type OpMsg struct {
	Op     Op     `json:"op"`
	Signer string `json:"signer"`
}

// Op ...
type Op struct {
	AccountName string `json:"accountName"`
	Expiration  uint   `json:"expiration"`
}

func makeRes(user string, sig string) map[string]interface{} {
	return map[string]interface{}{
		"code": 200, // 200:ok  400:fail
		"data": map[string]string{
			"accountName": user,
			"signer":      sig,
		},
	}
}
func checkIsUser(sig string, op Op) (user string, token string, expiration uint) {
	if op.Expiration > 1048239662000 {
		expiration = op.Expiration / 1000
	} else {
		expiration = op.Expiration
	}
	return op.AccountName, sig, expiration
}
func saveTokenExp(user string, token string, expiration uint) error {
	cybtoken := &m.CybToken{
		CybAccount: user,
		Signer:     token,
		Expiration: expiration,
	}
	err := cybtoken.SaveUniqueBy(m.CybToken{CybAccount: user})
	return err
}

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	opMsg := &OpMsg{}
	err := json.NewDecoder(r.Body).Decode(opMsg)
	if err != nil {

	}
	user, token, timestamp := checkIsUser(opMsg.Signer, opMsg.Op)
	saveTokenExp(user, token, timestamp)
	msg := makeRes(user, opMsg.Signer)
	u.Respond(w, msg, 200)

}
func findAsset() map[string]interface{} {
	return map[string]interface{}{
		"code": 200, // 200:ok  400:fail
		"data": []string{"BTC", "ETH"},
	}
}

// Asset ...
func Asset(w http.ResponseWriter, r *http.Request) {
	msg := findAsset()
	u.Respond(w, msg, 200)
}
func createCybexUserApp(user string) *m.App {
	newapp := &m.App{
		CybAccount: user,
	}
	newapp.Save()
	return newapp
}

type JPData struct {
	Status string `json:"status"`
	Data   struct {
		Address string `json:"address"`
		Type    string `json:"type"`
	} `json:"data"`
}

func createCybexUserAddress(addrQ *m.Address) *m.Address {
	asset := &m.Asset{}
	db := m.GetDB()
	db.First(asset, addrQ.AssetID)
	url := viper.GetString("usersrv.jpsrv_url") + "/api/address/new?type=" + asset.Name
	resp, err := http.Get(url)
	if err != nil {
	}
	_bodyBytes, err := ioutil.ReadAll(resp.Body)
	resObj := &JPData{}
	err = json.Unmarshal(_bodyBytes, resObj)
	addrQ.Address = resObj.Data.Address
	addrQ.Save()
	return addrQ
}
func findAppOrCreate(user string) *m.App {
	appQ := &m.App{
		CybAccount: user,
	}
	apps, err := rep.App.FetchWith(appQ)
	if err != nil {

	}
	var app1 *m.App
	if len(apps) == 0 {
		app1 = createCybexUserApp(user)
	} else {
		app1 = apps[0]
	}
	return app1
}
func findAssetByName(name string) (*m.Asset, error) {
	return rep.Asset.GetByName(name)
}
func findAddrOrCreate(app *m.App, asset *m.Asset) *m.Address {
	addrQ := &m.Address{
		AppID:   app.ID,
		AssetID: asset.ID,
	}
	addrs, _ := rep.Address.FetchWith(addrQ)
	var addr1 *m.Address
	if len(addrs) == 0 {
		addr1 = createCybexUserAddress(addrQ)
	} else {
		addr1 = addrs[0]
	}
	return addr1
}

// DepositAddress ...
func DepositAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	asset := vars["asset"]
	app1 := findAppOrCreate(user)
	asset1, err := findAssetByName(asset)
	if err != nil {
		u.Respond(w, u.Message(false, "asset not support!"), http.StatusBadRequest)
		return
	}
	addr := findAddrOrCreate(app1, asset1)
	msg := map[string]interface{}{
		"code": 200, // 200:ok  400:fail
		"data": map[string]interface{}{
			"address":  addr.Address,
			"asset":    asset1.Name, //TODO: cyb链的资产名
			"type":     asset1.Name,
			"createAt": addr.CreatedAt,
		},
	}
	u.Respond(w, msg, 200)
}
