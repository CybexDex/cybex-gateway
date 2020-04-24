package types

import (
	"cybex-gateway/utils/ss58"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/scale"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/centrifuge/go-substrate-rpc-client/types"
)

type OptionU128 struct {
	hasValue bool
	value    types.U128
}

func (o OptionU128) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionU128) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

type TradePair struct {
	Hash  types.Hash
	Base  types.Hash
	Quote types.Hash

	LatestMatchedPrice OptionU128

	OneDayTradeVolume  types.U128
	OneDayHighestPrice OptionU128
	OneDayLowestPrice  OptionU128
}

type TransferArgs struct {
	TokenHash types.Hash
	To        types.AccountID
	Amount    types.U128
	Memo      types.OptionBytes
}

func (args TransferArgs) String() string {
	ok, value := args.Memo.Unwrap()
	if ok {
		return fmt.Sprintf("(hash: %v), (toAccount: %v), (amount: %v), (memo: %v)", hexutil.Encode(args.TokenHash[:]), ss58.Encode(hexutil.Encode(args.To[:])), args.Amount, string(value[:]))
	}
	return fmt.Sprintf("(hash: %v), (toAccount: %v), (amount: %v)", hexutil.Encode(args.TokenHash[:]), ss58.Encode(hexutil.Encode(args.To[:])), args.Amount)
}
