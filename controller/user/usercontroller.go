package user

import (
	"fmt"
	"strconv"

	jp "cybex-gateway/controller/jpselect"
	model "cybex-gateway/model"
	"cybex-gateway/types"
	"cybex-gateway/utils"
	"cybex-gateway/utils/log"

	apim "github.com/CybexDex/cybex-go/api"
	"github.com/spf13/viper"
)

var api apim.CybexAPI

// InitNode ...
func InitNode() {
	node := viper.GetString("cybserver.node")
	api = apim.New(node, "")
	if err := api.Connect(); err != nil {
		panic(err)
	}
}

// GetRecordAsset ...
func GetRecordAsset(user string) ([]*model.RecordAsset, error) {
	res, err := model.JPOrderRecordAsset(user)
	return res, err
}

//GetRecord ...
func GetRecord(query *types.RecordsQuery) ([]*model.JPOrder, int, error) {
	res, count, err := model.JPOrderRecord(query.User, query.Asset, query.FundType, query.Size, query.LastID, query.Offset)
	return res, count, err
}

// RecordNotDone ...
func RecordNotDone(fromUpdate string, offset int, limit int) (res []*model.JPOrder, err error) {
	res, err = model.JPOrderNotDone(fromUpdate, offset, limit)
	return res, err
}

// CheckUser ...
func CheckUser(expiration string, user string, sig string) (isok bool, ex int, err error) {
	toSign := expiration + user
	log.Infoln(user, toSign, sig)
	re, err := api.VerifySign(user, toSign, sig)
	if err != nil {
		return false, 0, err
	}
	i, err := strconv.Atoi(expiration)
	if err != nil {
		return false, 0, err
	}
	if i > 1048239662000 {
		ex = i / 1000
	} else {
		ex = i
	}
	return re, ex, nil
}

// GetAssetsOne ...
func GetAssetsOne(asset string) (out *model.Asset, err error) {
	out, err = model.AssetsFind(asset)
	//
	// out = toResult(assets)
	return out, err
}

// GetAssets ...
func GetAssets() (out []*model.Asset, err error) {
	assets, err := model.AssetsAll()
	//
	// out = toResult(assets)
	return assets, err
}

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

// NewAddress ...
func NewAddress(user string, asset string) (address *types.UserResultAddress, err error) {
	address = &types.UserResultAddress{}
	assetobj, err := model.AssetsFind(asset)
	if err != nil {
		errmsg := fmt.Sprintf("%s %s", asset, "不是合法币种")
		log.Infoln(errmsg)
		return address, err
	}
	assetName := assetobj.Name
	if assetobj.JadeName != "" {
		assetName = assetobj.JadeName
	}
	//找user,asset的address
	//没有找到，获取，创建，返回
	var newaddr *types.JPAddressResult

	newaddr, err = jp.DepositAddress(assetName)

	if err != nil {
		return address, err
	}
	address1 := &model.Address{
		Address:    newaddr.Address,
		User:       user,
		Asset:      asset,
		BlockChain: "",
	}
	err = model.AddrssCreate(address1)
	if err != nil {
		return address, err
	}
	address.Address = address1.Address
	address.Asset = address1.Asset
	address.CreateAt = address1.CreatedAt
	address.CybName = assetobj.CYBName
	return address, nil
}

// VerifyAddress ...
func VerifyAddress(asset string, address string) (verifyRes *types.VerifyRes, err error) {
	var res *types.VerifyRes
	assetobj, err := model.AssetsFind(asset)
	if err != nil {
		errmsg := fmt.Sprintf("%s %s", asset, "不是合法币种")
		log.Infoln(errmsg)
		return res, err
	}
	assetName := assetobj.Name
	if assetobj.JadeName != "" {
		assetName = assetobj.JadeName
	}
	res, err = jp.VerifyAddress(assetName, address)
	return res, err
}

//GetAddress ...
func GetAddress(user string, asset string) (address *types.UserResultAddress, err error) {
	address = &types.UserResultAddress{}
	assetobj, err := model.AssetsFind(asset)
	if err != nil {
		errmsg := fmt.Sprintf("%s %s", asset, "不是合法币种")
		log.Infoln(errmsg)
		return address, err
	}
	assetName := assetobj.Name
	if assetobj.JadeName != "" {
		assetName = assetobj.JadeName
	}
	//找user,asset的address
	address1, err := model.AddressLast(user, asset)
	if err != nil {
		if err.Error() != "record not found" {
			return address, err
		}
		//没有找到，获取，创建，返回
		var newaddr *types.JPAddressResult
		newaddr, err = jp.DepositAddress(assetName)
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
	address.CybName = assetobj.CYBName
	return address, nil
}
