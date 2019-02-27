package adminsrv

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

	"git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/repository/jadepool"
	"github.com/gorilla/mux"

	utils "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

//GetAllJadepool ...
func GetAllJadepool(w http.ResponseWriter, r *http.Request) {
	if !checkAccount(r) {
		utils.Respond(w, utils.Message(false, "Unauthorized"), http.StatusUnauthorized)
		return
	}

	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepools, err := jadepoolRepo.FetchAll()
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", jadepools)
	utils.Respond(w, resp)
}

//CreateJadepool ...
func CreateJadepool(w http.ResponseWriter, r *http.Request) {
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

	jadepoolEntity := model.Jadepool{}
	err = json.Unmarshal(requestBody, &jadepoolEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	if len(jadepoolEntity.Name) == 0 {
		utils.Errorf("error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepool, err := jadepoolRepo.FetchWith(&model.Jadepool{Name: jadepoolEntity.Name})
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	if len(jadepool) != 0 {
		utils.Respond(w, utils.Message(false, "find one"), http.StatusBadRequest)
		return
	}
	err = jadepoolEntity.Create()
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

//GetJadepool ...
func GetJadepool(w http.ResponseWriter, r *http.Request) {
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

	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepool, err := jadepoolRepo.GetByID(uint(id))
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		utils.Errorf("GetByID error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success", jadepool)
	utils.Respond(w, resp)
}

//UpdateJadepool ...
func UpdateJadepool(w http.ResponseWriter, r *http.Request) {
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

	updateEntity := model.Jadepool{}
	err = json.Unmarshal(requestBody, &updateEntity)
	if err != nil {
		utils.Errorf("json.Unmarshal error: %v", err)
		utils.Respond(w, utils.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	jadepoolEntity, err := jadepoolRepo.GetByID(uint(id))
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

	err = jadepoolEntity.UpdateColumns(&updateEntity)
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

//DeleteJadepool ...
func DeleteJadepool(w http.ResponseWriter, r *http.Request) {
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

	jadepoolRepo := jadepool.NewRepo(model.GetDB())
	err = jadepoolRepo.DeleteByID(uint(id))
	if err != nil {
		utils.Errorf("Update error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}
