package model

import (
	"fmt"
	"github.com/cockroachdb/apd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateBalance(t *testing.T) {
	db := GetDB()

	bal := new(Balance)
	err := db.First(bal).Error
	assert.Nil(t, err)
	assert.Equal(t, false, db.NewRecord(bal))

	fmt.Println(bal)

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
	data["OutLockedFee"].Value = five

	err = ComputeBalance(db, bal, &data)
	assert.Nil(t, err)
}
