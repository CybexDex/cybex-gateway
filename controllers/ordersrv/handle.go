package ordersrv

import (
	"errors"
	"fmt"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

// IsOpen ...
func IsOpen(order1 *m.Order) (bool, error) {
	db := m.GetDB()
	asset := &m.Asset{}
	db.First(asset, order1.AssetID)
	if order1.Type == m.OrderTypeDeposit {
		return asset.DepositSwitch, nil
	} else if order1.Type == m.OrderTypeWithdraw {
		return asset.WithdrawSwitch, nil
	}
	return false, errors.New("unknown order type")
}

// IsBlack ...
func IsBlack(order1 *m.Order) (bool, error) {
	db := m.GetDB()
	app := &m.App{}
	db.First(app, order1.AppID)
	fmt.Println(app)
	if app.Status == "NORMAL" {
		return false, nil
	}
	errstr := fmt.Sprintf("error app status:%s", app.Status)
	return true, errors.New(errstr)
}

// IsBig ...
func IsBig(order1 *m.Order) (bool, error) {
	db := m.GetDB()
	big := &m.BigAsset{}
	db.Raw("SELECT * FROM big_assets WHERE asset_id = ? and type =? and big_amount < ?", order1.AssetID, order1.Type, order1.TotalAmount).Scan(&big)
	fmt.Println("big", big)
	if big.AssetID > 0 {
		return true, nil
	}
	return false, nil
}
