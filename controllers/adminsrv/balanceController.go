package adminsrv

import (
	"net/http"

	"coding.net/bobxuyang/cy-gateway-BN/models"
	"coding.net/bobxuyang/cy-gateway-BN/repository/app"
	"coding.net/bobxuyang/cy-gateway-BN/repository/balance"
	"github.com/gorilla/mux"

	utils "coding.net/bobxuyang/cy-gateway-BN/utils"
)

//GetBalance ...
func GetBalance(w http.ResponseWriter, r *http.Request) {
	if !checkAccount(r) {
		utils.Respond(w, utils.Message(false, "Unauthorized"), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	account := vars["account"]

	appRepo := app.NewRepo(model.GetDB())
	appEntitys, err := appRepo.FetchWith(&model.App{CybAccount: account})
	if err != nil {
		utils.Errorf("FetchWith error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	if len(appEntitys) == 0 {
		utils.Errorf("FetchWith empty")
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}
	balanceRepo := balance.NewRepo(model.GetDB())
	balances, err := balanceRepo.FetchWith(&model.Balance{AppID: appEntitys[0].ID})
	if err != nil {
		utils.Errorf("FetchWith error: %v", err)
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	if len(balances) == 0 {
		utils.Errorf("FetchWith empty")
		utils.Respond(w, utils.Message(false, "Internal server error"), http.StatusInternalServerError)
		return
	}

	m := make(map[string]interface{})
	m["balance"] = balances[0].Balance.Text('f')
	m["inLocked"] = balances[0].InLocked.Text('f')
	m["outLocked"] = balances[0].OutLocked.Text('f')
	m["inLockedFee"] = balances[0].InLockedFee.Text('f')
	m["outLockedFee"] = balances[0].OutLockedFee.Text('f')

	resp := utils.Message(true, "success", m)
	utils.Respond(w, resp)
}
