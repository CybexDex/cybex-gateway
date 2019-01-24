package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	app "git.coding.net/bobxuyang/cy-gateway-BN/app"
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
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	// generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.Errorf("Operation failed when compute hash of password, %s", err.Error())
		u.Respond(w, u.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	// save data to DB
	account.PasswordHash = string(hash)
	//repo := accRepo.NewRepo(m.GetDB())
	err = account.Save()
	if err != nil {
		u.Errorf("Create account error, \"%s\"", err.Error())
		if !strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			u.Respond(w, u.Message(false, "Create account error"), http.StatusInternalServerError)
		} else {
			u.Respond(w, u.Message(false, "Username or email is duplicated"), http.StatusBadRequest)
		}
		return
	}

	// return data to client
	account.PasswordHash = ""
	u.Respond(w, u.Message(true, "OK", account))
}

//Authenticate ...
var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	// parse data into object
	account := &m.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Errorf("Invalid request: %v, %s", account, err.Error())
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	password := account.Password
	name := account.Name

	// get account data from DB
	repo := accRepo.NewRepo(m.GetDB())
	account, err = repo.GetByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Errorf("Username or passord incorrect, %s", err.Error())
			u.Respond(w, u.Message(false, "Username or passord incorrect"), http.StatusUnauthorized)
			return
		}

		u.Errorf("Fetch account data error, %s", err.Error())
		u.Respond(w, u.Message(false, "Fetch account data error"), http.StatusInternalServerError)
		return
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		u.Errorf("Username or passord incorrect, %s", err.Error())
		u.Respond(w, u.Message(false, "Username or passord incorrect"), http.StatusUnauthorized)
		return
	} else if err != nil {
		u.Errorf("Internal server error, %s", err.Error())
		u.Respond(w, u.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	//Create JWT token
	tk := &app.Token{
		UserID: account.ID,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// TODO: need to remove hard code
	// TODO: need to remove hard code
	tokenString, err := token.SignedString([]byte("token_password"))
	if err != nil {
		u.Errorf("Internal server error, %s", err.Error())
		u.Respond(w, u.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	// return data to client
	u.Respond(w, u.Message(true, "Login OK", map[string]string{"token": tokenString}))
}

//DebugInfo ...
var DebugInfo = func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("UserID")
	u.Infof("Login UserID:", id)
	u.Respond(w, u.Message(true, "Login OK", map[string]interface{}{"UserID": id}))
}
