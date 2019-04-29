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
	// log.Infoln(d.String())
	return d
}

// HandleTR ...
func (a *BBBHandler) HandleTR(op *operations.TransferOperation, tx *cybTypes.SignedTransaction, prefix string) {
	// log.Infoln("HandleTX", op.To, tx.Signatures)
	// 是否在币种中，没有的话，是否是gateway账号的UR。
	gatewayTo := allgateways[op.To.String()]
	gatewayFrom := allgateways[op.From.String()]
	// 先看From,是充值或者Inner订单
	if gatewayFrom != nil {
		log.Infof("HandleTX:,from:%s,op:%+v,tx.sig:%v\n", gatewayFrom.Account.Name, op, tx.Signatures)
		// 直接看sig,能不能找到,找到就更新。没有找到的话。先不管
		// if gatewayTo != nil {
		sig := tx.Signatures[0].String()
		orders, err := model.JPOrderFind(&model.JPOrder{
			Sig:          sig,
			CurrentState: model.JPOrderStatusPending,
		})
		if err != nil {
			log.Errorln("JPOrderFind error", err)
			return
		}
		if len(orders) == 1 {
			order := orders[0]
			if order.Type == model.JPOrderTypeDeposit {
				order.SetCurrent("done", model.JPOrderStatusDone, "")
				order.SetStatus(model.JPOrderStatusDone)
			}
			if order.Type == model.JPOrderTypeWithdraw {
				order.SetCurrent("order", model.JPOrderStatusInit, "")
			}
			err := order.Save()
			if err != nil {
				log.Errorln("save error", err)
			}
		} else if len(orders) > 1 {
			log.Errorln("sig len ", len(orders))
		}
		return
	}
	//to gatewayin 的话就是充值,排除from gatewayout, from 特殊账户的
	if gatewayTo != nil {
		log.Infof("HandleTX:,to:%s,op:%+v,tx.sig:%v\n", gatewayTo.Account.Name, op, tx.Signatures)
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
		// log.Infoln("withdraw:", addr, *op)
		// 创建jporder对象
		realAmount := amountToReal(op.Amount.Amount, assetChain.Precision)
		jporder := &model.JPOrder{
			Asset:      assetConf.Name,
			BlockChain: "",
			// BNOrderID:  "",
			CybUser: fromUser.Name,
			OutAddr: addr,

			TotalAmount:  realAmount,
			Type:         "WITHDRAW",
			Status:       model.JPOrderStatusPending,
			Current:      "cybinner",
			CurrentState: model.JPOrderStatusInit,
			CYBHash:      &prefix,
		}
		err = jporder.Save()
		if err != nil {
			log.Errorln("save jporder error", err)
		}
		log.Infof("order:%d,%s:%+v\n", jporder.ID, "save_withdraw", *jporder)
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
		// log.Infoln(assetC.Name)
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
	// log.Infoln(allgateways, allAssets)
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
	for txnum, tx := range block.Transactions {
		for opnum, op := range tx.Operations {
			if op == nil {
				// cancel all 等无法识别的type就会解析不到 op
				continue
			}
			opbyte, _ := json.Marshal(op)
			var opt operations.TransferOperation
			err = json.Unmarshal(opbyte, &opt)
			if err != nil {
				log.Errorln(cnum, txnum, opnum, op)
				continue
			}
			if opt.From.String() != "" {
				handler.HandleTR(&opt, &tx, fmt.Sprintf("%d:%d:%d", cnum, txnum, opnum))
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
		if err == apim.ErrShutdown {
			api.Connect()
		}
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
	if err !=nil {
		log.Errorln(err)
		return
	}
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
