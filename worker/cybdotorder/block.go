package cybdotorder

import (
	"cybex-gateway/model"
	types "cybex-gateway/types/cybexdot"
	"cybex-gateway/utils"
	"cybex-gateway/utils/log"
	"cybex-gateway/utils/ss58"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/spf13/viper"

	gsTypes "github.com/centrifuge/go-substrate-rpc-client/types"
)

type TransferOperation struct {
	from     string
	to       string
	token    string
	amount   string
	memo     string
	hash     string
	identity string
}

func amountToReal(amountin int64, prercison int64) decimal.Decimal {
	d := decimal.New(amountin, int32(-1*prercison))
	// log.Infoln(d.String())
	return d
}

func BlockRead() {
	i := 0
	for {
		i = i + 1
		// log.Debugln("read round:", i)
		handleBlock()
		time.Sleep(time.Second * 3)
	}
}

func HandleBlockNum(cnum uint64) {
	_, err := readBlock(cnum)
	if err != nil {
		log.Errorln(err)
	}
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
	// log.Debugln("next block", lastBlockNum, "head", blockheadNum, err)
	if lastBlockNum >= blockheadNum {
		return
	}
	// for
	for cnum := lastBlockNum; cnum <= blockheadNum; cnum = cnum + 1 {
		t, err := readBlock(uint64(cnum))
		if err != nil {
			log.Errorln("handleBlockNum", cnum, err)
			return
		}
		err = updateLastBlock(cnum, t, easy)
		if err != nil {
			log.Errorln("block updateLastBlock", cnum, err)
			return
		}
	}
	//
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

func getHeadNum() (int, error) {
	hash, err := api.RPC.Chain.GetFinalizedHead()
	if err != nil {
		return 0, err
	}

	header, err := api.RPC.Chain.GetHeader(hash)
	if err != nil {
		return 0, err
	}
	return int(header.Number), nil
}

func updateLastBlock(cnum int, t time.Time, easy *model.Easy) error {
	bstr := strconv.Itoa(cnum + 1)
	easy.Value = bstr
	easy.RecordTime = t
	err := easy.Save()
	return err
}

func readBlock(num uint64) (time.Time, error) {
	hash, err := api.RPC.Chain.GetBlockHash(num)
	if err != nil {
		return time.Time{}, err
	}

	block, err := api.RPC.Chain.GetBlock(hash)

	if err != nil {
		return time.Time{}, err
	}

	meta, err := api.RPC.State.GetMetadata(hash)
	if err != nil {
		return time.Time{}, err
	}

	var transferExtrinsicsIndexes []int
	operations := make(map[int]TransferOperation)
	var timestamp gsTypes.UCompact
	for i, ext := range block.Block.Extrinsics {
		moduleName, fmeta := functionMeta(meta, ext.Method)
		if moduleName == "Timestamp" && fmeta.Name == "set" {
			err := gsTypes.DecodeFromBytes(ext.Method.Args, &timestamp)
			if err != nil {
				return time.Time{}, err
			}
		}
		if moduleName == "TokenModule" && fmeta.Name == "transfer" {
			var args types.TransferArgs
			err := gsTypes.DecodeFromBytes(ext.Method.Args, &args)
			if err != nil {
				return time.Time{}, err
			}
			transferExtrinsicsIndexes = append(transferExtrinsicsIndexes, i)
			if ext.Signature.Signer.IsAccountID {
				from := ss58.Encode(hexutil.Encode(ext.Signature.Signer.AsAccountID[:]))
				to := ss58.Encode(hexutil.Encode(args.To[:]))
				ok, value := args.Memo.Unwrap()
				var memo string
				if ok {
					memo = string(value)
				} else {
					memo = ""
				}

				encodedBytes, err := gsTypes.EncodeToBytes(ext)
				if err != nil {
					return time.Time{}, err
				}

				checksum, _ := blake2b.New(32, []byte{})
				checksum.Write(encodedBytes)
				h := checksum.Sum(nil)
				extrinsicHash := "0x" + utils.BytesToHex(h)
				operations[i] = TransferOperation{from, to, hexutil.Encode(args.TokenHash[:]), args.Amount.String(), memo, extrinsicHash, fmt.Sprintf("%d:%d:%d", num, i, len(block.Block.Extrinsics))}
			}
			log.Debugf("\tTransfer Extrinsics:: (index=%#v)(module=%#v)(method=%#v)(args=%v)\n", i, moduleName, fmeta.Name, args)
		}
	}

	key, err := gsTypes.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		return time.Time{}, err
	}

	var er gsTypes.EventRecordsRaw
	err = api.RPC.State.GetStorage(key, &er, hash)
	if err != nil {
		return time.Time{}, err
	}

	e := types.Events{}
	err = er.DecodeEventRecords(meta, &e)
	if err != nil {
		return time.Time{}, err
	}

	for _, e := range e.System_ExtrinsicSuccess {
		for _, v := range transferExtrinsicsIndexes {
			if int(e.Phase.AsApplyExtrinsic) == v {
				operation := operations[v]
				handlerTransfer(operation)
				log.Debugln("transfer finalized index: ", v)
			}
		}
	}

	return time.Unix(int64(timestamp)/1000, 0), err
}

func handlerTransfer(op TransferOperation) bool {
	log.Debugf("handler transfer op (%#v)", op)
	findasset, err := model.AssetsFrist(&model.Asset{
		CYBID: op.token,
	})

	if err != nil {
		log.Errorln(op.token, err)
		return false
	}

	if findasset.GatewayAccount == op.from { // deposit
		log.Infof("HandleTX 从gateway账号打出:,from:%s, asset:%s, op:%+v\n", findasset.GatewayAccount, findasset.Name,
			op)
		// 直接看sig,能不能找到,找到就更新。没有找到的话。先不管
		orders, err := model.JPOrderFind(&model.JPOrder{
			Sig:          op.hash,
			CurrentState: model.JPOrderStatusPending,
		})
		if err != nil {
			log.Errorln("JPOrderFind error", err)
			return false
		}
		if len(orders) == 1 {
			order := orders[0]
			if order.Type == model.JPOrderTypeDeposit {
				order.CYBHash = &op.identity
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
		} else {
			log.Errorln(op.identity, "没有需要更新的充值订单")
		}
		return true

	} else if findasset.GatewayAccount == op.to { // withdraw
		memo := op.memo
		if memo == "" {
			log.Infoln("UR Memo", op)
			return false
		}

		gatewayPrefix := findasset.WithdrawPrefix
		if !strings.HasPrefix(memo, gatewayPrefix) {
			log.Infoln("UR:1 ", "memo:", memo)
			return false
		}
		s := strings.TrimPrefix(memo, gatewayPrefix)
		s2 := strings.Split(s, ":")
		if len(s2) < 2 {
			log.Infoln("UR:2", "memo:", memo)
			return false
		}
		addr := s2[1]

		n, err := strconv.ParseInt(op.amount, 10, 64)
		if err != nil {
			log.Errorln(op.amount, err)
			return false
		}

		np, err := strconv.ParseInt(findasset.Precision, 10, 64)
		if err != nil {
			log.Errorln(findasset.Precision, err)
			return false
		}

		realAmount := amountToReal(n, np)

		jporder := &model.JPOrder{
			Asset:      findasset.Name,
			CybAsset:   findasset.CYBName,
			BlockChain: findasset.Blockchain,
			// BNOrderID:  "",
			CybUser:      op.from,
			OutAddr:      addr,
			Memo:         memo,
			TotalAmount:  realAmount,
			Type:         "WITHDRAW",
			Status:       model.JPOrderStatusPending,
			Current:      "order",
			CurrentState: model.JPOrderStatusInit,
			Sig:          op.hash,
			CYBHash:      &op.identity,
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

func functionMeta(m *gsTypes.Metadata, call gsTypes.Call) (gsTypes.Text, gsTypes.FunctionMetadataV4) {
	var fmeta gsTypes.FunctionMetadataV4
	var moduleName gsTypes.Text

	mi := uint8(0)

	switch {
	case m.IsMetadataV4:
		for index, mod := range m.AsMetadataV4.Modules {
			if mi == call.CallIndex.SectionIndex {
				module := m.AsMetadataV4.Modules[index]
				fmeta = module.Calls[call.CallIndex.MethodIndex]
				moduleName = mod.Name

				break
			}
			if !mod.HasCalls {
				continue
			}
			mi++
		}
	case m.IsMetadataV7:
		for index, mod := range m.AsMetadataV4.Modules {
			if mi == call.CallIndex.SectionIndex {
				module := m.AsMetadataV7.Modules[index]
				fmeta = module.Calls[call.CallIndex.MethodIndex]
				moduleName = mod.Name

				break
			}
			if !mod.HasCalls {
				continue
			}
			mi++
		}
	case m.IsMetadataV8:
		for index, mod := range m.AsMetadataV8.Modules {
			if mi == call.CallIndex.SectionIndex {
				module := m.AsMetadataV8.Modules[index]
				fmeta = module.Calls[call.CallIndex.MethodIndex]
				moduleName = mod.Name

				break
			}
			if !mod.HasCalls {
				continue
			}
			mi++
		}
	case m.IsMetadataV9:
		for index, mod := range m.AsMetadataV9.Modules {
			if mi == call.CallIndex.SectionIndex {
				module := m.AsMetadataV9.Modules[index]
				fmeta = module.Calls[call.CallIndex.MethodIndex]
				moduleName = mod.Name

				break
			}
			if !mod.HasCalls {
				continue
			}
			mi++
		}
	case m.IsMetadataV10:
		for index, mod := range m.AsMetadataV10.Modules {
			if mi == call.CallIndex.SectionIndex {
				module := m.AsMetadataV10.Modules[index]
				fmeta = module.Calls[call.CallIndex.MethodIndex]
				moduleName = mod.Name

				break
			}
			if !mod.HasCalls {
				continue
			}
			mi++
		}
	}
	return moduleName, fmeta
}
