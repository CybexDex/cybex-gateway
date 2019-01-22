package controllers

import (
	"encoding/json"
	"net/http"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	accRepo "git.coding.net/bobxuyang/cy-gateway-BN/repository/account"
	u "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

//CreateAccount ...
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &m.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	repo := accRepo.NewRepo(m.GetDB())
	err = repo.Create(account)
	if err != nil {
		u.Respond(w, u.Message(false, "create account error"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

////Authenticate ...
//var Authenticate = func(w http.ResponseWriter, r *http.Request) {
//	account := &models.Account{}
//	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
//	if err != nil {
//		u.Respond(w, u.Message(false, "Invalid request"))
//		return
//	}
//	resp := models.Login(account.Email, account.Password)
//	u.Respond(w, resp)
//}
