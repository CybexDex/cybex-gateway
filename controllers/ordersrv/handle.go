package ordersrv

import (
	"errors"

	rep "coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "coding.net/bobxuyang/cy-gateway-BN/models"
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
	var addressBlockChain string
	if order1.Type == m.OrderTypeDeposit {
		jporder, err := rep.JPOrder.GetByID(order1.JPOrderID)
		if err != nil {
			return false, err
		}
		addressBlockChain = jporder.From
	} else if order1.Type == m.OrderTypeWithdraw {
		cyborder, err := rep.CybOrder.GetByID(order1.CybOrderID)
		if err != nil {
			return false, err
		}
		addressBlockChain = cyborder.WithdrawAddr
	}
	app, err := rep.App.GetByID(order1.AppID)
	if err != nil {
		return false, err
	}
	blacks, err := rep.Black.FetchWith(&m.Black{
		Blockchain: blockchain.Name,
		Address:    addressBlockChain,
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
	bigs, err := rep.BigAsset.FindBig(&m.BigAsset{
		AssetID: order1.AssetID,
		Type:    order1.Type,
	}, order1.TotalAmount)
	if err != nil {
		return false, err
	}
	if len(bigs) > 0 {
		return true, nil
	}
	return false, nil
}
