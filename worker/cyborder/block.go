package cyborder

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/types"
	"bitbucket.org/woyoutlz/bbb-gateway/utils"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	"coding.net/yundkyy/cybexgolib/operations"
	cybTypes "coding.net/yundkyy/cybexgolib/types"
	"github.com/spf13/viper"
)

// BBBHandler ...
type BBBHandler struct {
}

// HandleTR ...
func (a *BBBHandler) HandleTR(op *operations.TransferOperation, tx *cybTypes.SignedTransaction) {
	log.Infoln("bbb", op.From, tx.Signatures)
	//to gatewayin 的话就是充值,排除from gatewayout, from 特殊账户的
}

var allgateways []types.GatewayAccount
var allAssets map[string]*types.AssetConfig

// InitAsset ...初始化asset gateway 账户
func InitAsset() {
	allAssets = make(map[string]*types.AssetConfig)
	assets := viper.GetStringMap("assets")
	for keyname, asset := range assets {
		keyName := strings.ToUpper(keyname)
		assetC := types.AssetConfig{}
		err := utils.V2S(asset, &assetC)
		if err != nil {
			panic(err)
		}
		log.Infoln(assetC.Name)
		allAssets[assetC.Name] = &assetC
		account1, err := api.GetAccountByName(assetC.Deposit.Gateway)
		if account1 == nil {
			log.Errorln("gateway account 不存在", assetC.Deposit.Gateway)
			panic("")
		}
		account2, err := api.GetAccountByName(assetC.Withdraw.Gateway)
		if account2 == nil {
			log.Errorln("gateway account 不存在", assetC.Withdraw.Gateway)
			panic("")
		}
		g1 := types.GatewayAccount{
			Account: account1,
			Type:    "DEPOSIT",
			Asset:   keyName,
		}
		g2 := types.GatewayAccount{
			Account: account2,
			Type:    "WITHDRAW",
			Asset:   keyName,
		}
		allgateways = append(allgateways, g1, g2)
	}
	log.Infoln(allgateways)
}

// Test ...
func Test() {
	handleBlockNum(6966405)
}
func readBlock(cnum int, handler types.HandleInterface) ([]string, error) {
	blockInfo, err := api.GetBlock(uint64(cnum))
	if err != nil {
		return nil, err
	}
	bs, _ := json.Marshal(blockInfo)
	var block cybTypes.Block
	err = json.Unmarshal(bs, &block)
	if err != nil {
		return nil, err
	}
	for _, tx := range block.Transactions {
		for _, op := range tx.Operations {
			opbyte, _ := json.Marshal(op)
			var opt operations.TransferOperation
			err = json.Unmarshal(opbyte, &opt)
			if err != nil {
				return nil, err
			}
			if opt.From.String() != "" {
				handler.HandleTR(&opt, &tx)
			}
		}
	}
	return nil, nil
}
func handleBlockNum(cnum int) {
	handler := BBBHandler{}
	cyborders, err := readBlock(cnum, &handler)
	if err != nil {
		log.Errorln(err)
		return
	}
	// log.Infoln(cyborders)
	// save cyborders
	for _, order := range cyborders {
		log.Infoln(order)
	}
}
func getHeadNum() (int, error) {
	s, err := api.GetDynamicGlobalProperties()
	if err != nil {
		return 0, err
	}
	return int(s.LastIrreversibleBlockNum), nil
}
func getlastBlock() (int, *model.Easy, error) {
	// is there a last
	s, err := model.EasyFind(&model.Easy{
		Key: "cybLastBlockNum",
	})
	if err != nil {
		return 0, nil, err
	}
	if len(s) > 0 {
		c := s[0]
		i, err := strconv.Atoi(c.Value)
		if err != nil {
			return 0, nil, err
		}
		return i, c, nil
	}
	// create it
	blockBegin := viper.GetInt("cybserver.blockBegin")
	if blockBegin == -1 {
		head, err := getHeadNum()
		if err != nil {
			return 0, nil, err
		}
		blockBegin = head
	}
	bstr := strconv.Itoa(blockBegin)
	newLast := model.Easy{
		Key:   "cybLastBlockNum",
		Value: bstr,
	}
	err = newLast.Save()
	if err != nil {
		return 0, nil, fmt.Errorf("cannot find")
	}

	return blockBegin, &newLast, nil
}
func handleBlock() {
	// get lastBlock
	lastBlockNum, easy, err := getlastBlock()
	if err != nil {
		log.Errorln("getlastBlock:", err)
		return
	}
	// get blockhead
	blockheadNum, err := getHeadNum()
	log.Debugln("last", lastBlockNum, "head", blockheadNum, err)
	if lastBlockNum >= blockheadNum {
		return
	}
	// for
	for cnum := lastBlockNum; cnum <= blockheadNum; cnum = cnum + 1 {
		handleBlockNum(cnum)
		updateLastBlock(cnum, easy)
	}
	//
}
func updateLastBlock(cnum int, easy *model.Easy) error {
	bstr := strconv.Itoa(cnum + 1)
	easy.Value = bstr
	err := easy.Save()
	return err
}

// BlockRead ...
func BlockRead() {
	i := 0
	for {
		i = i + 1
		log.Debugln("read round:", i)
		handleBlock()
		time.Sleep(time.Second * 3)
	}
}
