package cybdotorder

import (
	types "cybex-gateway/types/cybexdot"
	"cybex-gateway/utils/log"
	"fmt"

	gsTypes "github.com/centrifuge/go-substrate-rpc-client/types"
)

func HandleBlockNum(cnum uint64) {
	err := readBlock(cnum)
	if err != nil {
		log.Errorln(err)
	}
}

func readBlock(num uint64) error {
	hash, err := api.RPC.Chain.GetBlockHash(num)
	if err != nil {
		return err
	}
	block, err := api.RPC.Chain.GetBlock(hash)
	if err != nil {
		return err
	}

	for _, ext := range block.Block.Extrinsics {
		fmt.Printf("\tExtrinsics:: (method=%#v)\n", ext.Method)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	key, err := gsTypes.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		return err
	}

	var er gsTypes.EventRecordsRaw
	err = api.RPC.State.GetStorage(key, &er, hash)
	if err != nil {
		return err
	}

	e := types.Events{}
	err = er.DecodeEventRecords(meta, &e)
	if err != nil {

		return err
	}

	for _, e := range e.System_ExtrinsicSuccess {
		fmt.Printf("\tSystem:ExtrinsicSuccess:: (phase=%#v)\n", e.Phase)
	}
	//log.Debugln(j)

	return nil
}
