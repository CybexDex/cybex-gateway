package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/blockchain"
	"github.com/gorilla/mux"

	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

//GetAllBlockchain ...
func GetAllBlockchain(w http.ResponseWriter, r *http.Request) {
	blockchainRepo := blockchain.NewRepo(model.GetDB())
	blockchains, err := blockchainRepo.FetchAll()
	if err != nil {
		utils.Errorf("Update error: %v", err)
		w.WriteHeader(400)
		return
	}

	resp := utils.Message(true, "success")
	resp["data"] = blockchains
	utils.Respond(w, resp)
}

//CreateBlockchain ...
func CreateBlockchain(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Errorf("ReadAll error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
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
	err = blockchainRepo.Create(&blockchainEntity)
	if err != nil {
		utils.Errorf("Update error: %v", err)
		w.WriteHeader(400)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

//GetBlockchain ...
func GetBlockchain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Errorf("Atoi error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	blockchainRepo := blockchain.NewRepo(model.GetDB())
	blockchain, err := blockchainRepo.GetByID(uint(id))
	if err != nil {
		utils.Errorf("Update error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	resp["data"] = blockchain
	utils.Respond(w, resp)
}

//UpdateBlockchain ...
func UpdateBlockchain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Errorf("Atoi error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Errorf("ReadAll error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	utils.Debugf("request: %s", requestBody)

	blockchainEntity := model.Blockchain{}
	err = json.Unmarshal(requestBody, &blockchainEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	blockchainRepo := blockchain.NewRepo(model.GetDB())
	err = blockchainRepo.Update(uint(id), &blockchainEntity)
	if err != nil {
		utils.Errorf("Update error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}
