package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

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
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", blockchains)
	utils.Respond(w, resp)
}

//CreateBlockchain ...
func CreateBlockchain(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Errorf("ReadAll error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	utils.Debugf("request:\n %s", requestBody)

	blockchainEntity := model.Blockchain{}
	if len(blockchainEntity.Name) == 0 {
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(requestBody, &blockchainEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	blockchainRepo := blockchain.NewRepo(model.GetDB())
	blockchain, err := blockchainRepo.FetchWith(&model.Blockchain{Name: blockchainEntity.Name})
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	if len(blockchain) != 0 {
		utils.Respond(w, utils.Message(false, "find one"), http.StatusBadRequest)
		return
	}
	err = blockchainEntity.Create()
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
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
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	blockchainRepo := blockchain.NewRepo(model.GetDB())
	blockchain, err := blockchainRepo.GetByID(uint(id))
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		utils.Errorf("GetByID error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", blockchain)
	utils.Respond(w, resp)
}

//UpdateBlockchain ...
func UpdateBlockchain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Errorf("Atoi error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Errorf("ReadAll error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	utils.Debugf("request: %s", requestBody)

	blockchainEntity := model.Blockchain{}
	err = json.Unmarshal(requestBody, &blockchainEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	//blockchainRepo := blockchain.NewRepo(model.GetDB())
	blockchainEntity.ID = uint(id)
	err = blockchainEntity.UpdateColumns(&blockchainEntity)
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

//DeleteBlockchain ...
func DeleteBlockchain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Errorf("Atoi error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	blockchainRepo := blockchain.NewRepo(model.GetDB())
	err = blockchainRepo.DeleteByID(uint(id))
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}
