package cyborder

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"bitbucket.org/woyoutlz/bbb-gateway/model"
	"bitbucket.org/woyoutlz/bbb-gateway/types"
	"bitbucket.org/woyoutlz/bbb-gateway/utils"
	"bitbucket.org/woyoutlz/bbb-gateway/utils/log"
	apim "coding.net/yundkyy/cybexgolib/api"
	"coding.net/yundkyy/cybexgolib/operations"
	cybTypes "coding.net/yundkyy/cybexgolib/types"
	"github.com/spf13/viper"
)

// BBBHandler ...
type BBBHandler struct {
}

func amountToReal(amountin cybTypes.Int64, prercison int) decimal.Decimal {
	d := decimal.New(int64(amountin), int32(-1*prercison))
	log.Infoln(d.String())
	return d
}

// HandleTR ...
func (a *BBBHandler) HandleTR(op *operations.TransferOperation, tx *cybTypes.SignedTransaction) {
	log.Infoln("bbb", op.To, tx.Signatures)
	// 是否在币种中，没有的话，是否是gateway账号的UR。
	gatewayTo := allgateways[op.To.String()]
	//to gatewayin 的话就是充值,排除from gatewayout, from 特殊账户的
	if gatewayTo != nil {
		// log.Infoln("gatewayTo", *gatewayTo)
		fromUsers, err := api.GetAccounts(op.From)
		fromUser := fromUsers[0]
		assetChain, err := api.GetAsset(op.Amount.Asset.String())
		if err != nil {
			log.Errorln(assetChain, err)
			return
		}
		var assetConf *types.AssetConfig
		for _, assetC := range allAssets {
			if assetC.Withdraw.Coin == assetChain.Symbol {
				assetConf = assetC
			}
		}
		if assetConf == nil || assetConf.Withdraw.Gateway != gatewayTo.Account.Name {
			// UR
			log.Infoln("UR 不是合法币种", assetConf, assetChain.Symbol)
			return
		}
		// 合法转账
		// 锁定期
		extensions := op.Extensions
		if len(extensions) > 0 {
			log.Infoln("UR extensions", op)
			return
		}
		// 没有memo
		memo := op.Memo
		if memo == nil {
			log.Infoln("UR Memo", *op)
			return
		}
		// memo格式
		memoout := ""
		gatewayMemoPri := gatewayTo.MemoPri
		for _, prv := range gatewayMemoPri {
			s, err := memo.Decrypt(&prv)
			if err == nil {
				memoout = s
			}
		}
		// log.Infoln(memoout, *op)
		gatewayPrefix := assetConf.Withdraw.Memopre
		if !strings.HasPrefix(memoout, gatewayPrefix) {
			log.Infoln("UR:1 ", "memo:", memoout)
			return
		}
		s := strings.TrimPrefix(memoout, gatewayPrefix)
		s2 := strings.Split(s, ":")
		if len(s2) < 3 {
			log.Infoln("UR:2", "memo:", memoout)
			return
		}
		addr := strings.Join(s2[2:], ":")
		log.Infoln("withdraw:", addr, *op)
		// 创建jporder对象
		realAmount := amountToReal(op.Amount.Amount, assetChain.Precision)
		jporder := &model.JPOrder{
			Asset:      assetConf.Name,
			BlockChain: "",
			BNOrderID:  "",
			CybUser:    fromUser.Name,

			From: fromUser.Name,
			To:   addr,
			Memo: memoout,

			TotalAmount:  realAmount,
			Type:         "WITHDRAW",
			Status:       "PENDING",
			Current:      "cybinner",
			CurrentState: "INIT",
		}
		log.Infoln(*jporder)
	}
}

var allgateways map[string]*types.GatewayAccount
var allAssets map[string]*types.AssetConfig
var assetsOfChain map[string]*cybTypes.Asset

// InitAsset ...初始化asset gateway 账户
func InitAsset() {
	allAssets = make(map[string]*types.AssetConfig)
	allgateways = make(map[string]*types.GatewayAccount)
	assetsOfChain = make(map[string]*cybTypes.Asset)
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
		gatewaykeyBag := apim.KeyBagByUserPass(assetC.Deposit.Gateway, assetC.Deposit.Gatewaypass)
		memokey := account1.Options.MemoKey
		pubkeys := cybTypes.PublicKeys{memokey}
		gatewayMemoPri := gatewaykeyBag.PrivatesByPublics(pubkeys)
		g1 := types.GatewayAccount{
			Account: account1,
			Type:    "DEPOSIT",
			Asset:   keyName,
			MemoPri: gatewayMemoPri,
		}
		account2, err := api.GetAccountByName(assetC.Withdraw.Gateway)
		if account2 == nil {
			log.Errorln("gateway account 不存在", assetC.Withdraw.Gateway)
			panic("")
		}
		gatewaykeyBag2 := apim.KeyBagByUserPass(assetC.Withdraw.Gateway, assetC.Withdraw.Gatewaypass)
		memokey2 := account2.Options.MemoKey
		pubkeys2 := cybTypes.PublicKeys{memokey2}
		gatewayMemoPri2 := gatewaykeyBag2.PrivatesByPublics(pubkeys2)
		g2 := types.GatewayAccount{
			Account: account2,
			Type:    "WITHDRAW",
			Asset:   keyName,
			MemoPri: gatewayMemoPri2,
		}
		allgateways[g1.Account.ID.String()] = &g1
		allgateways[g2.Account.ID.String()] = &g2
	}
	log.Infoln(allgateways, allAssets)
}

// Test ...
func Test() {
	handleBlockNum(6993893) // 6993893
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
