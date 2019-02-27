package adminsrv

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

	"git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/asset"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/blockchain"
	"github.com/gorilla/mux"

	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

//GetAllAsset ...
func GetAllAsset(w http.ResponseWriter, r *http.Request) {
	if !checkAccount(r) {
		utils.Respond(w, utils.Message(false, "Unauthorized"), http.StatusUnauthorized)
		return
	}

	assetRepo := asset.NewRepo(model.GetDB())
	assets, err := assetRepo.FetchAll()
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", assets)
	utils.Respond(w, resp)
}

//CreateAsset ...
func CreateAsset(w http.ResponseWriter, r *http.Request) {
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

	assetEntity := model.Asset{}
	err = json.Unmarshal(requestBody, &assetEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	if len(assetEntity.Name) == 0 {
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	assetRepo := blockchain.NewRepo(model.GetDB())
	asset, err := assetRepo.FetchWith(&model.Blockchain{Name: assetEntity.Name})
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	if len(asset) != 0 {
		utils.Respond(w, utils.Message(false, "find one"), http.StatusBadRequest)
		return
	}
	err = assetEntity.Create()
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

//GetAsset ...
func GetAsset(w http.ResponseWriter, r *http.Request) {
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

	assetRepo := asset.NewRepo(model.GetDB())
	asset, err := assetRepo.GetByID(uint(id))
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		utils.Errorf("GetByID error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", asset)
	utils.Respond(w, resp)
}

//UpdateAsset ...
func UpdateAsset(w http.ResponseWriter, r *http.Request) {
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

	updateEntity := model.Asset{}
	err = json.Unmarshal(requestBody, &updateEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	assetRepo := asset.NewRepo(model.GetDB())
	assetEntity, err := assetRepo.GetByID(uint(id))
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
	err = assetEntity.UpdateColumns(&updateEntity)
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

//DeleteAsset ...
func DeleteAsset(w http.ResponseWriter, r *http.Request) {
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

	assetRepo := asset.NewRepo(model.GetDB())
	err = assetRepo.DeleteByID(uint(id))
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}
