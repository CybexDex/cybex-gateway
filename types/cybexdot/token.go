package types

import (
	"cybex-gateway/utils/ss58"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/centrifuge/go-substrate-rpc-client/types"
)

type TransferArgs struct {
	TokenHash types.Hash
	To        types.AccountID
	Amount    types.U128
}

func (args TransferArgs) String() string {
	return fmt.Sprintf("(hash: %v), (toAccount: %v), (amount: %v)", hexutil.Encode(args.TokenHash[:]), ss58.Encode(hexutil.Encode(args.To[:])), args.Amount)
}
