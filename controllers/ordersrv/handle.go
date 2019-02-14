package ordersrv

import (
	"errors"

	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

// IsOpen ...
func IsOpen(order1 *m.Order) (bool, error) {
	asset, err := rep.Asset.GetByID(order1.AssetID)
	if err != nil {
		return false, err
	}
	if order1.Type == m.OrderTypeDeposit {
		return asset.DepositSwitch, nil
	} else if order1.Type == m.OrderTypeWithdraw {
		return asset.WithdrawSwitch, nil
	}
	return false, errors.New("unknown order type")
}

// IsBlack ...
func IsBlack(order1 *m.Order) (bool, error) {
	asset, err := rep.Asset.GetByID(order1.AssetID)
	if err != nil {
		return false, err
	}
	blockchain, err := rep.Blockchain.GetByID(asset.BlockchainID)
	if err != nil {
		return false, err
	}
	jporder, err := rep.JPOrder.GetByID(order1.JPOrderID)
	if err != nil {
		return false, err
	}
	app, err := rep.App.GetByID(order1.AppID)
	if err != nil {
		return false, err
	}
	blacks, err := rep.Black.FetchWithOr(&m.Black{
		Blockchain: blockchain.Name,
		Address:    jporder.To,
	}, &m.Black{
		Blockchain: blockchain.Name,
		Address:    jporder.From,
	})
	if err != nil {
		return false, err
	}
	if len(blacks) > 0 {
		return true, nil
	}
	blacks, err = rep.Black.FetchWith(&m.Black{
		Blockchain: "CYB",
		Address:    app.CybAccount,
	})
	if err != nil {
		return false, err
	}
	if len(blacks) > 0 {
		return true, nil
	}
	return false, nil
}

// IsBig ...
func IsBig(order1 *m.Order) (bool, error) {
	db := m.GetDB()
	big := &m.BigAsset{}
	db.Raw("SELECT * FROM big_assets WHERE asset_id = ? and type =? and big_amount < ?", order1.AssetID, order1.Type, order1.TotalAmount).Scan(&big)
	// fmt.Println("big", big)
	if big.AssetID > 0 {
		return true, nil
	}
	return false, nil
}
