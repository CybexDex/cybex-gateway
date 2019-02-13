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
	if !checkAccount(r) {
		utils.Respond(w, utils.Message(false, "Unauthorized"), http.StatusUnauthorized)
		return
	}

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
	if !checkAccount(r) {
		utils.Respond(w, utils.Message(false, "Unauthorized"), http.StatusUnauthorized)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Errorf("ReadAll error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	utils.Debugf("request:\n %s", requestBody)

	blockchainEntity := model.Blockchain{}
	err = json.Unmarshal(requestBody, &blockchainEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	if len(blockchainEntity.Name) == 0 {
		utils.Errorf("error: %v", err)
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
		utils.Errorf("Create error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

//GetBlockchain ...
func GetBlockchain(w http.ResponseWriter, r *http.Request) {
	if !checkAccount(r) {
		utils.Respond(w, utils.Message(false, "Unauthorized"), http.StatusUnauthorized)
		return
	}

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
	if !checkAccount(r) {
		utils.Respond(w, utils.Message(false, "Unauthorized"), http.StatusUnauthorized)
		return
	}

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

	updateEntity := model.Blockchain{}
	err = json.Unmarshal(requestBody, &updateEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	blockchainRepo := blockchain.NewRepo(model.GetDB())
	blockchainEntity, err := blockchainRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Errorf("error: %v", err)
			utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadGateway)
			return
		}
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusInternalServerError)
		return
	}

	err = blockchainEntity.UpdateColumns(&updateEntity)
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
	if !checkAccount(r) {
		utils.Respond(w, utils.Message(false, "Unauthorized"), http.StatusUnauthorized)
		return
	}

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
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

func checkAccount(r *http.Request) bool {
	/*id := r.Context().Value("UserID")
	if id == nil {
		return false
	}
	accountRepo := account.NewRepo(model.GetDB())
	_, err := accountRepo.GetByID(id.(uint))
	if err != nil {
		utils.Errorf("Update error: %v", err)
		return false
	}*/

	// todo: check account role

	return true
}
