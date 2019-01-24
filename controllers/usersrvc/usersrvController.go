package usersrvc

import (
	"encoding/json"
	"fmt"
	"net/http"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	u "git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/gorilla/mux"
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
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	opMsg := &OpMsg{}
	err := json.NewDecoder(r.Body).Decode(opMsg)
	if err != nil {

	}
	user, token, timestamp := checkIsUser(opMsg.Signer, opMsg.Op)
	fmt.Println(user, token, timestamp)
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
func findUserAddr(user string, asset string) string {
	return "nil"
}

// DepositAddress ...
func DepositAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	asset := vars["asset"]
	addr := findUserAddr(user, asset)
	// if addr == nil {
	// 	addr = createUserAddr(user, asset)
	// }
	msg := map[string]interface{}{
		"code": 200, // 200:ok  400:fail
		"data": map[string]interface{}{
			"address": addr,
			"asset":   "BTC",
		},
	}
	u.Respond(w, msg, 200)
}
