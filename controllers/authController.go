package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	accRepo "git.coding.net/bobxuyang/cy-gateway-BN/repository/account"
	u "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

//CreateAccount ...
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	// parse data into object
	account := &m.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	password := account.Password
	account.Password = ""
	if err != nil {
		u.Errorf("Invalid request: %v, %s", account, err.Error())
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	// generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.Errorf("Operation failed when compute hash of password, %s", err.Error())
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	// save data to DB
	account.PasswordHash = string(hash)
	repo := accRepo.NewRepo(m.GetDB())
	err = repo.Create(account)
	if err != nil {
		u.Errorf("Create account error, %s", err.Error())
		u.Respond(w, u.Message(false, "Create account error"))
		return
	}

	// return data to client
	account.PasswordHash = ""
	u.RespondObj(w, account)
}

//Authenticate ...
var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	// parse data into object
	account := &m.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Errorf("Invalid request: %v, %s", account, err.Error())
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	password := account.Password
	name := account.Name

	// get account data from DB
	repo := accRepo.NewRepo(m.GetDB())
	account, err = repo.GetByName(name)
	if err != nil {
		u.Errorf("Fetch account data error, %s", err.Error())
		u.Respond(w, u.Message(false, "Fetch account data error"))
		return
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		u.Errorf("Passord incorrect, %s", err.Error())
		u.Respond(w, u.Message(false, "Passord incorrect"))
		return
	}

	// return data to client
	u.Respond(w, u.Message(true, "password correct"))
}
