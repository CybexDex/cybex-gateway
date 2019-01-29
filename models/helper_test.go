package model

import (
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	"testing"

	"github.com/cockroachdb/apd"
)

func TestUpdateBalance(t *testing.T) {
	db := GetDB()

	bal := Balance{}
	db.Where("ID=?", 1).First(&bal)
	//fmt.Println(bal)

	data := GetBalanceInitData()
	data["Balance"].Oper = "ADD"
	one := apd.New(123, 0)
	data["Balance"].Value = one

	err := ComputeBalance(db, &bal, &data)
	if err != nil {
		utils.Errorln(err.Error())
	}
}
