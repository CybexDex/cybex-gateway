package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

//OrderNotiRequest ...
type OrderNotiRequest struct {
	ID            string                 `json:"id"`
	State         string                 `json:"state"`
	BizType       string                 `json:"bizType"`
	CoinType      string                 `json:"coinType"`
	From          string                 `json:"from"`
	To            string                 `json:"to"`
	Value         string                 `json:"value"`
	Confirmations int                    `json:"confirmations"`
	CreateAt      int64                  `json:"create_at"`
	UpdateAt      int64                  `json:"update_at"`
	Fee           string                 `json:"fee"`
	Hash          string                 `json:"hash"`
	Data          map[string]interface{} `json:"data"`
	ExtraData     string                 `json:"extraData"`
	Memo          string                 `json:"memo"`
}

//OrderNoti ...
func OrderNoti(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.GetLogger().Errorf("ReadAll error: %v", err)
		w.WriteHeader(400)
		return
	}
	request := OrderNotiRequest{}
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		utils.GetLogger().Errorf("Unmarshal error: %v, request body: %s", err, requestBody)
		w.WriteHeader(400)
		return
	}
	utils.GetLogger().Infof("oder noti request: %s", requestBody)
	//todo: save log

	//todo: find exorder and save/update exorder

	//todo: response
	resp := map[string]interface{}{}
	utils.Respond(w, resp)
}
