package usersrvc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	apim "coding.net/yundkyy/cybexgolib/api"
	"coding.net/yundkyy/cybexgolib/types"
	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	u "git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/gorilla/mux"
	"github.com/juju/errors"
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
	Expiration  int    `json:"expiration"`
}

var api apim.BitsharesAPI

func init() {
	u.InitConfig()
	node := viper.GetString("cybsrv.node")
	api = apim.New(node, "")
	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "OnConnect"))
	}
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
func checkIsUser(sig string, op Op) (isok bool, expiration int, err error) {
	ss := types.Fund{
		AccountName: []byte(op.AccountName),
		Expiration:  op.Expiration,
	}
	re, err := api.LoginVerify(ss, sig)
	if err != nil {
		return false, 0, err
	}
	if op.Expiration > 1048239662000 {
		expiration = op.Expiration / 1000
	} else {
		expiration = op.Expiration
	}
	return re, expiration, nil
}
func saveTokenExp(user string, token string, expiration uint) error {
	//TODO: redis may be the better solution
	cybtoken := &m.CybToken{
		CybAccount: user,
		Signer:     token,
		Expiration: expiration,
	}
	err := cybtoken.SaveUniqueBy(m.CybToken{CybAccount: user})
	return err
}

// IsTokenOK ...
func IsTokenOK(token string) (bool, error) {

	t, err := rep.CybToken.FetchWith(&m.CybToken{
		Signer: token,
	})
	if err != nil {
		return false, err
	}

	if len(t) > 0 {
		timenow := time.Now().Unix()
		u.Infoln(t, token, t[0].Expiration, uint(timenow))
		if t[0].Expiration < uint(timenow) {
			u.Infoln("expire")
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	opMsg := &OpMsg{}
	err := json.NewDecoder(r.Body).Decode(opMsg)
	if err != nil {
		u.Errorf("err: %v", err)
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	isok, expiration, err := checkIsUser(opMsg.Signer, opMsg.Op)
	if err != nil {
		u.Errorf("err: %v", err)
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	if !isok {
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	err = saveTokenExp(opMsg.Op.AccountName, opMsg.Signer, uint(expiration))
	if err != nil {
		u.Errorf("err: %v", err)
		u.Respond(w, u.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	msg := makeRes(opMsg.Op.AccountName, opMsg.Signer)
	u.Respond(w, msg, 200)

}
func findAsset() ([]*m.Asset, error) {
	assets, err := rep.Asset.FetchAll()
	if err != nil {
		return nil, err
	}
	return assets, err
}

// AllAsset ...
func AllAsset(w http.ResponseWriter, r *http.Request) {
	assets, err := findAsset()
	if err != nil {
		u.Errorf("err: %v", err)
		u.Respond(w, u.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	msg := map[string]interface{}{
		"code": 200, // 200:ok  400:fail
		"data": assets,
	}
	u.Respond(w, msg, 200)
}
func createCybexUserApp(user string) (*m.App, error) {
	newapp := &m.App{
		CybAccount: user,
	}
	err := newapp.Save()
	return newapp, err
}

type JPData struct {
	Status bool `json:"status"`
	Data   struct {
		Address string `json:"address"`
		Type    string `json:"type"`
	} `json:"data"`
}

func createCybexUserAddress(addrQ *m.Address) (*m.Address, error) {
	asset, err := rep.Asset.GetByID(addrQ.AssetID)
	if err != nil {
		return nil, err
	}
	url := viper.GetString("usersrv.jpsrv_url") + "/api/address/new?type=" + asset.Name
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	_bodyBytes, err := ioutil.ReadAll(resp.Body)
	resObj := &JPData{}
	err = json.Unmarshal(_bodyBytes, resObj)
	if err != nil {
		u.Errorf("error %v", err)
		return nil, err
	}
	addrQ.Address = resObj.Data.Address
	err = addrQ.Save()
	if err != nil {
		u.Errorf("error %v", err)
		return nil, err
	}
	return addrQ, nil
}

// func findAppOrCreate(user string) (app1 *m.App, err error) {
// 	appQ := &m.App{
// 		CybAccount: user,
// 	}
// 	apps, err := rep.App.FetchWith(appQ)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(apps) == 0 {
// 		app1, err = createCybexUserApp(user)
// 		if err != nil {
// 			return app1, err
// 		}
// 	} else {
// 		app1 = apps[0]
// 	}
// 	return app1, nil
// }
func findAssetByName(name string) (*m.Asset, error) {
	return rep.Asset.GetByName(name)
}
func findAddrOrCreate(app *m.App, asset *m.Asset) (*m.Address, error) {
	addrQ := &m.Address{
		AppID:   app.ID,
		AssetID: asset.ID,
	}
	addrs, err := rep.Address.FetchWith(addrQ)
	if err != nil {
		u.Errorln(err)
		return nil, err
	}
	var addr1 *m.Address
	if len(addrs) == 0 {
		addr1, err = createCybexUserAddress(addrQ)
		if err != nil {
			return nil, err
		}
	} else {
		addr1 = addrs[0]
	}
	return addr1, nil
}

// DepositAddress ...
func DepositAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	asset := vars["asset"]
	app1, err := rep.App.FindAppOrCreate(user)
	if err != nil {
		u.Errorf("error %v", err)
		u.Respond(w, u.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	asset1, err := findAssetByName(asset)
	if err != nil {
		u.Errorf("error %v", err)
		u.Respond(w, u.Message(false, "asset not support!"), http.StatusBadRequest)
		return
	}
	addr, err := findAddrOrCreate(app1, asset1)
	if err != nil {
		u.Errorf("error %v", err)
		u.Respond(w, u.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
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
