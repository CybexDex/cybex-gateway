package user

import (
	"bitbucket.org/woyoutlz/bbb-gateway/controller/jp"
	model "bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/types"
	"bitbucket.org/woyoutlz/bbb-gateway/utils"
	"github.com/spf13/viper"
)

// GetBBBAssets ...
func GetBBBAssets() (out []*types.UserResultBBB, err error) {
	assetsConf := viper.GetStringMap("assets")
	for _, conf := range assetsConf {
		assetC := types.AssetConfig{}
		err := utils.V2S(conf, &assetC)
		if err != nil {
			return nil, err
		}
		assetO := &types.UserResultBBB{
			Name:            assetC.Name,
			Blockchain:      assetC.BlockChain,
			DepositAs:       assetC.Deposit.JustAsset,
			WithdrawAsset:   assetC.Withdraw.Coin,
			WithdrawGateway: assetC.Withdraw.Gateway,
			WithdrawPrefix:  assetC.Withdraw.Memopre,

			DepositSwitch:  assetC.Deposit.Switch,
			WithdrawSwitch: assetC.Withdraw.Switch,

			// MinDeposit:  assetC.Withdraw.Memopre,
			// MinWithdraw: assetC.Withdraw.Memopre,
			WithdrawFee: assetC.Withdraw.Fee,
			// DepositFee:  assetC.Withdraw.Memopre,
		}
		out = append(out, assetO)
	}
	return out, nil
}

//GetAddress ...
func GetAddress(user string, asset string) (address *types.UserResultAddress, err error) {
	address = &types.UserResultAddress{}
	//找user,asset的address
	address1, err := model.AddressLast(user, asset)
	if err != nil {
		if err.Error() != "record not found" {
			return address, err
		}
		//没有找到，获取，创建，返回
		newaddr, err := jp.DepositAddress(asset)
		if err != nil {
			return address, err
		}
		address1 = &model.Address{
			Address:    newaddr.Address,
			User:       user,
			Asset:      asset,
			BlockChain: "",
		}
		err = model.AddrssCreate(address1)
		if err != nil {
			return address, err
		}
	}
	//返回
	address.Address = address1.Address
	address.Asset = address1.Asset
	address.CreateAt = address1.CreatedAt
	return address, nil
}
