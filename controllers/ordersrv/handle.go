package ordersrv

import (
	"errors"

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
	// order.asset.blockchain.Name  order.jporder.to order.jporder.from
	// order.app.cybname

	db := m.GetDB()
	blockchain := &m.Blockchain{}
	asset := &m.Asset{}
	jporder := &m.JPOrder{}
	app := &m.App{}
	db.First(asset, order1.AssetID)
	db.First(blockchain, asset.BlockchainID)
	db.First(jporder, order1.JPOrderID)
	db.First(app, order1.AppID)
	// fmt.Println(blockchain.Name, jporder.To, jporder.From, app.CybAccount)
	black := &m.Black{}
	db.Where(&m.Black{
		Blockchain: blockchain.Name,
		Address:    jporder.To,
	}).Or(&m.Black{
		Blockchain: blockchain.Name,
		Address:    jporder.From,
	}).First(black)
	if black.ID > 0 {
		return true, errors.New("black:" + black.Address)
	}
	db.Where(&m.Black{
		Blockchain: "CYB",
		Address:    app.CybAccount,
	}).First(black)
	if black.ID > 0 {
		return true, errors.New("black:" + black.Address)
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
