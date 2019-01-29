package model

import (
	"testing"

	"git.coding.net/bobxuyang/cy-gateway-BN/utils"

	"github.com/cockroachdb/apd"
)

func TestUpdateBalance(t *testing.T) {
	db := GetDB()

	bal := Balance{}
	db.Where("ID=?", 1).First(&bal)
	//fmt.Println(bal)

	data := GetBalanceInitData()

	data["Balance"].Oper = "ADD"
	one := apd.New(1, 0)
	data["Balance"].Value = one

	data["InLocked"].Oper = "ADD"
	two := apd.New(2, 0)
	data["InLocked"].Value = two

	data["OutLocked"].Oper = "SUB"
	three := apd.New(3, 0)
	data["OutLocked"].Value = three

	data["InLockedFee"].Oper = "ADD"
	four := apd.New(123, 0)
	data["InLockedFee"].Value = four

	data["OutLockedFee"].Oper = "SUB"
	five := apd.New(123, 0)
	data["OutLocked"].Value = five

	err := ComputeBalance(db, &bal, &data)
	if err != nil {
		utils.Errorln(err.Error())
	}
}
