package cybdotorder

import (
	"cybex-gateway/model"
	types "cybex-gateway/types/cybexdot"
	"cybex-gateway/utils/log"
	"cybex-gateway/utils/ss58"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/spf13/viper"

	gsTypes "github.com/centrifuge/go-substrate-rpc-client/types"
)

type TransferOperation struct {
	from   string
	to     string
	token  string
	amount string
	memo   string
	sig    gsTypes.MultiSignature
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

				operations[i] = TransferOperation{from, to, hexutil.Encode(args.TokenHash[:]), args.Amount.String(), "", ext.Signature.Signature}
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

func handlerTransfer(op TransferOperation) (bool, error) {
	log.Debugf("handler transfer op (%#v)", op)
	findasset, err := model.AssetsFrist(&model.Asset{
		CYBID: op.token,
	})

	if err != nil {
		return false, err
	}

	if findasset.GatewayID == op.from { // deposit
		return true, nil

	} else if findasset.GatewayID == op.to { // withdraw
		return true, nil
	}

	return false, nil
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
