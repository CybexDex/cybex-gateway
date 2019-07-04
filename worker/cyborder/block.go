package cyborder

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"cybex-gateway/model"
	"cybex-gateway/types"
	"cybex-gateway/utils"
	"cybex-gateway/utils/log"

	apim "github.com/CybexDex/cybex-go/api"
	"github.com/CybexDex/cybex-go/operations"
	cybTypes "github.com/CybexDex/cybex-go/types"
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

// CheckUR ...
func (a *BBBHandler) CheckUR(op *operations.TransferOperation, tx *cybTypes.SignedTransaction, prefix string) {
	toid := op.To.String()
	toasset, err := model.AssetsFrist(&model.Asset{
		GatewayID: toid,
	})
	if err != nil {
		// 不是处理中的资产
		// log.Errorln(err)
		return
	}
	if toasset == nil {
		// return
	}
	// 相关至少是UR
	fromUsers, err := api.GetAccounts(op.From)
	fromUser := fromUsers[0]
	assetChain, err := api.GetAsset(op.Amount.Asset.String())
	if err != nil {
		log.Errorln(assetChain, err)
		return
	}
	realAmount := amountToReal(op.Amount.Amount, assetChain.Precision)
	memo := op.Memo
	memoout := ""
	if memo == nil {

	} else {
		account1, _ := api.GetAccountByName(toasset.GatewayAccount)
		gatewaypass := utils.SeedString(toasset.GatewayPass)
		gatewaykeyBag := apim.KeyBagByUserPass(toasset.GatewayAccount, gatewaypass)
		memokey := account1.Options.MemoKey
		pubkeys := cybTypes.PublicKeys{memokey}
		gatewayMemoPri := gatewaykeyBag.PrivatesByPublics(pubkeys)
		for _, prv := range gatewayMemoPri {
			s, err := memo.Decrypt(&prv)
			if err == nil {
				memoout = s
			}
		}
	}
	// memo格式

	jporder := &model.JPOrder{
		Asset: op.Amount.Asset.String(),

		// BNOrderID:  "",
		CybUser:     fromUser.Name,
		CYBHash2:    toasset.GatewayAccount,
		TotalAmount: realAmount,
		Type:        "UR",
		Memo:        memoout,
		// Status:       model.JPOrderStatusPending,
		// Current:      "order",
		// CurrentState: model.JPOrderStatusInit,
		Sig:     tx.Signatures[0].String(),
		CYBHash: &prefix,
	}
	err = jporder.Save()
	if err != nil {
		log.Errorln("save UR error", err)
	}
	log.Infof("order:%d,%s:%+v\n", jporder.ID, "save_ur", *jporder)
}

// HandleTR ...
func (a *BBBHandler) HandleTR(op *operations.TransferOperation, tx *cybTypes.SignedTransaction, prefix string) bool {
	// log.Infoln("HandleTX", op.To, tx.Signatures)
	// 是否在币种中，不是的话，是否是gateway账号的UR。
	toid := op.To.String()
	fromid := op.From.String()
	assetid := op.Amount.Asset.String()
	findasset, err := model.AssetsFrist(&model.Asset{
		CYBID: assetid,
	})
	if err != nil {
		// log.Errorln(err)
	}
	if findasset == nil {
		return false
	}
	// 先看From,是充值或者Inner订单
	if findasset.GatewayID == fromid {
		account1, _ := api.GetAccountByName(findasset.GatewayAccount)
		gatewaypass := utils.SeedString(findasset.GatewayPass)
		gatewaykeyBag := apim.KeyBagByUserPass(findasset.GatewayAccount, gatewaypass)
		memokey := account1.Options.MemoKey
		pubkeys := cybTypes.PublicKeys{memokey}
		gatewayMemoPri := gatewaykeyBag.PrivatesByPublics(pubkeys)
		gatewayFrom := types.GatewayAccount{
			Account: account1,
			Type:    "DEPOSIT",
			Asset:   findasset.Name,
			MemoPri: gatewayMemoPri,
		}
		log.Infoln(findasset.Name)
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
			return false
		}
		if len(orders) == 1 {
			order := orders[0]
			if order.Type == model.JPOrderTypeDeposit {
				order.CYBHash = &prefix
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
		return true
	}
	//to gatewayTo 的话就是提现订单
	if findasset.GatewayID == toid {
		account1, _ := api.GetAccountByName(findasset.GatewayAccount)
		gatewaypass := utils.SeedString(findasset.GatewayPass)
		gatewaykeyBag := apim.KeyBagByUserPass(findasset.GatewayAccount, gatewaypass)
		memokey := account1.Options.MemoKey
		pubkeys := cybTypes.PublicKeys{memokey}
		gatewayMemoPri := gatewaykeyBag.PrivatesByPublics(pubkeys)
		gatewayTo := types.GatewayAccount{
			Account: account1,
			Type:    "WITHDRAW",
			Asset:   findasset.Name,
			MemoPri: gatewayMemoPri,
		}
		log.Infof("HandleTX:,to:%s,op:%+v,tx.sig:%v\n", gatewayTo.Account.Name, op, tx.Signatures)
		fromUsers, err := api.GetAccounts(op.From)
		fromUser := fromUsers[0]
		assetChain, err := api.GetAsset(op.Amount.Asset.String())
		if err != nil {
			log.Errorln(assetChain, err)
			return false
		}
		assetConf, err := model.AssetsFrist(&model.Asset{
			CYBName: assetChain.Symbol,
		})
		if err != nil {
			log.Infoln("UR 不是合法币种", assetChain.Symbol)
			return false
		}
		if assetConf == nil || assetConf.GatewayAccount != gatewayTo.Account.Name {
			// UR
			log.Infoln("UR 不是该账户合法币种", assetConf, assetChain.Symbol)
			return false
		}
		// 合法转账
		// 锁定期
		extensions := op.Extensions
		if len(extensions) > 0 {
			log.Infoln("UR extensions", op)
			return false
		}
		// 没有memo
		memo := op.Memo
		if memo == nil {
			log.Infoln("UR Memo", *op)
			return false
		}
		// memo格式
		memoout := ""
		for _, prv := range gatewayMemoPri {
			s, err := memo.Decrypt(&prv)
			if err == nil {
				memoout = s
			}
		}
		// log.Infoln(memoout, *op)
		gatewayPrefix := assetConf.WithdrawPrefix
		if !strings.HasPrefix(memoout, gatewayPrefix) {
			log.Infoln("UR:1 ", "memo:", memoout)
			return false
		}
		s := strings.TrimPrefix(memoout, gatewayPrefix)
		s2 := strings.Split(s, ":")
		if len(s2) < 3 {
			log.Infoln("UR:2", "memo:", memoout)
			return false
		}
		addr := strings.Join(s2[2:], ":")
		// log.Infoln("withdraw:", addr, *op)
		// 创建jporder对象
		realAmount := amountToReal(op.Amount.Amount, assetChain.Precision)
		jporder := &model.JPOrder{
			Asset:      assetConf.Name,
			BlockChain: assetConf.Blockchain,
			// BNOrderID:  "",
			CybUser:      fromUser.Name,
			OutAddr:      addr,
			Memo:         memoout,
			TotalAmount:  realAmount,
			Type:         "WITHDRAW",
			Status:       model.JPOrderStatusPending,
			Current:      "order",
			CurrentState: model.JPOrderStatusInit,
			Sig:          tx.Signatures[0].String(),
			CYBHash:      &prefix,
		}
		err = jporder.Save()
		if err != nil {
			log.Errorln("save jporder error", err)
		}
		log.Infof("order:%d,%s:%+v\n", jporder.ID, "save_withdraw", *jporder)
		return true
	}
	return false
}

var allgateways map[string]*types.GatewayAccount

// InitAsset ...初始化asset gateway 账户
func InitAsset() {
	log.Infoln("开始检查资产...")
	// allAssets = make(map[string]*types.AssetConfig)
	allgateways = make(map[string]*types.GatewayAccount)
	assets, err := model.AssetsAll()
	if err != nil {
		log.Errorln("InitAsset", err)
	}
	for _, asset := range assets {
		account1, _ := api.GetAccountByName(asset.GatewayAccount)
		assetcyb, _ := api.GetAsset(asset.CYBName)
		changed := false
		if asset.CYBID != assetcyb.ID.String() {
			log.Infoln("更新cybid", asset.Name, asset.CYBID, "=>", assetcyb.ID.String())
			asset.CYBID = assetcyb.ID.String()
			changed = true
		}
		if account1 == nil {
			log.Errorln("gateway account 不存在", asset.GatewayAccount)
			continue
			// panic("")
		}
		gatewaypass := utils.SeedString(asset.GatewayPass)
		gatewaykeyBag := apim.KeyBagByUserPass(asset.GatewayAccount, gatewaypass)
		memokey := account1.Options.MemoKey
		pubkeys := cybTypes.PublicKeys{memokey}
		gatewayMemoPri := gatewaykeyBag.PrivatesByPublics(pubkeys)
		g1 := types.GatewayAccount{
			Account: account1,
			Type:    "DEPOSIT",
			Asset:   asset.Name,
			MemoPri: gatewayMemoPri,
		}
		allgateways[g1.Account.ID.String()] = &g1
		if asset.GatewayID != account1.ID.String() {
			log.Infoln("更新gatewayid", asset.Name, asset.GatewayID, "=>", account1.ID.String())
			asset.GatewayID = account1.ID.String()
			changed = true
		}
		if changed {
			err := asset.Save()
			if err != nil {
				log.Errorln("更新asset失败", asset.Name, err)
			}
		}
	}
	// log.Infoln(allgateways, allAssets)
}

// Test ...
func Test() {
	handleBlockNum(7571221) // 6993893
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
			if op.Type() == opt.Type() {
				var opt operations.TransferOperation
				err = json.Unmarshal(opbyte, &opt)
				if err != nil {
					log.Infoln("非交易", cnum, txnum, opnum, op.Type(), err)
					continue
				}
				if opt.From.String() != "" {
					ok := handler.HandleTR(&opt, &tx, fmt.Sprintf("%d:%d:%d", cnum, txnum, opnum))
					if !ok {
						handler.CheckUR(&opt, &tx, fmt.Sprintf("%d:%d:%d", cnum, txnum, opnum))
					}
				}
			}
		}
	}
	return nil, nil
}

// HandleBlockNum ...
func HandleBlockNum(cnum int) {
	handler := BBBHandler{}
	cyborders, err := readBlock(cnum, &handler)
	if err != nil {
		log.Errorln(err)
		// if err == apim.ErrShutdown {
		err = api.Connect()
		if err != nil {
			log.Errorln(err)
		}
		// }
		return
	}
	// log.Infoln(cyborders)
	// save cyborders
	for _, order := range cyborders {
		log.Infoln(order)
	}
}

// UpdateLastTime ...
func UpdateLastTime(cnum int) (t time.Time, error error) {
	blockInfo, err := api.GetBlock(uint64(cnum))
	if err != nil {
		return t, err
	}
	bs, _ := json.Marshal(blockInfo)
	var block cybTypes.Block
	err = json.Unmarshal(bs, &block)
	if err != nil {
		return t, err
	}
	return block.TimeStamp.Time, nil

}
func handleBlockNum(cnum int) error {
	handler := BBBHandler{}
	cyborders, err := readBlock(cnum, &handler)
	if err != nil {
		log.Errorln("readBlock", cnum, err)
		// if err == apim.ErrShutdown {
		err2 := api.Connect()
		if err2 != nil {
			log.Errorln(err2)
		}
		// }
		return err
	}
	// log.Infoln(cyborders)
	// save cyborders
	for _, order := range cyborders {
		log.Infoln(order)
	}
	return nil
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
	if err != nil {
		log.Errorln(err)
		return
	}
	easyhead, err := model.EasyFristOrCreate("cybHeadNum")
	if err != nil {
		log.Errorln(err)
		return
	}
	easyhead.Value = fmt.Sprintf("%d", blockheadNum)
	easyhead.Save()
	log.Debugln("last", lastBlockNum, "head", blockheadNum, err)
	if lastBlockNum >= blockheadNum {
		return
	}
	// for
	for cnum := lastBlockNum; cnum <= blockheadNum; cnum = cnum + 1 {
		err := handleBlockNum(cnum)
		if err != nil {
			log.Errorln("handleBlockNum", cnum, err)
			return
		}
		t, err := UpdateLastTime(cnum)
		if err != nil {
			log.Errorln("block UpdateLastTime", cnum, err)
			return
		}
		updateLastBlock(cnum, t, easy)
	}
	//
}
func updateLastBlock(cnum int, t time.Time, easy *model.Easy) error {
	bstr := strconv.Itoa(cnum + 1)
	easy.Value = bstr
	easy.RecordTime = t
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
