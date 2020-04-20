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

	meta, err := api.RPC.State.GetMetadata(hash)
	if err != nil {
		return err
	}

	var transferExtrinsicsIndexes []int

	for i, ext := range block.Block.Extrinsics {
		moduleName, fmeta := functionMeta(meta, ext.Method)
		if moduleName == "TokenModule" && fmeta.Name == "transfer" {
			var args types.TransferArgs
			err := gsTypes.DecodeFromBytes(ext.Method.Args, &args)
			if err != nil {
				return err
			}
			transferExtrinsicsIndexes = append(transferExtrinsicsIndexes, i)
			fmt.Printf("\tExtrinsics:: (module=%#v)(method=%#v)(args=%v)\n", moduleName, fmeta.Name, args)
		}
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

	for successExtrinsicIndex, e := range e.System_ExtrinsicSuccess {
		fmt.Printf("\tSystem:ExtrinsicSuccess:: (phase=%#v)\n", e.Phase)
		for _, v := range transferExtrinsicsIndexes {
			if int(e.Phase.AsApplyExtrinsic) == successExtrinsicIndex && e.Phase.IsFinalization && e.Phase.IsApplyExtrinsic {
				log.Debugln("transfer finalized index: ", v)
			}
		}
	}

	return nil
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
