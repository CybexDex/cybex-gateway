package user

import (
	"bitbucket.org/woyoutlz/bbb-gateway/controller/jp"
	model "bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/types"
)

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
