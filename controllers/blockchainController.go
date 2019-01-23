package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/blockchain"

	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

//CreateBlockchain ...
func CreateBlockchain(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Errorf("ReadAll error: %v", err)
		w.WriteHeader(400)
		return
	}
	utils.Debugf("request:\n %s", requestBody)

	blockchainEntity := model.Blockchain{}
	err = json.Unmarshal(requestBody, &blockchainEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		w.WriteHeader(400)
		return
	}

	blockchainRepo := blockchain.NewRepo(model.GetDB())
	err = blockchainRepo.Update(&blockchainEntity)
	if err != nil {
		utils.Errorf("Update error: %v", err)
		w.WriteHeader(400)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}
