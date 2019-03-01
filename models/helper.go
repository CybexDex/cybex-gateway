package model

import (
	"reflect"

	u "coding.net/bobxuyang/cy-gateway-BN/utils"
	"github.com/cockroachdb/apd"
	"github.com/jinzhu/gorm"
)

//BalanceData ...
type BalanceData struct {
	Value *apd.Decimal
	Oper  string // ADD or SUB
}

//GetBalanceInitData ...
func GetBalanceInitData() map[string]*BalanceData {

	result := make(map[string]*BalanceData)
	result["Balance"] = &BalanceData{Value: nil, Oper: ""}
	result["InLocked"] = &BalanceData{Value: nil, Oper: ""}
	result["OutLocked"] = &BalanceData{Value: nil, Oper: ""}
	result["InLockedFee"] = &BalanceData{Value: nil, Oper: ""}
	result["OutLockedFee"] = &BalanceData{Value: nil, Oper: ""}

	return result
}

//ComputeBalance ...
func ComputeBalance(tx *gorm.DB, bal *Balance, d *map[string]*BalanceData) error {
	var err error

	for name, data := range *d {
		if data.Oper == "" || data.Value == nil {
			continue
		}

		vv := reflect.ValueOf(bal)
		vv = vv.Elem()
		vv = vv.FieldByName(name)
		ptr := vv.Interface().(*apd.Decimal)

		if data.Oper == "ADD" {
			condition, err := apd.BaseContext.Add(ptr, ptr, data.Value)
			if err != nil {
				u.Errorf("apd decimal ADD error,", err, bal.ID)
				return err
			}
			if condition.Any() {
				u.Errorf("apd decimal ADD error,", condition.String(), bal.ID)
				return err
			}
		} else {
			condition, err := apd.BaseContext.Sub(ptr, ptr, data.Value)
			if err != nil {
				u.Errorf("apd decimal SUB error,", err, bal.ID)
				return err
			}
			if condition.Any() {
				u.Errorf("apd decimal SUB error,", condition.String(), bal.ID)
				return err
			}
		}
	}

	// save balance
	err = tx.Save(bal).Error
	if err != nil {
		u.Errorf("save balance error,", err, bal.ID)
		return err
	}

	return err
}
