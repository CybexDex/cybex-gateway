package types

import (
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

type Trade struct {
	Hash        types.Hash
	Base        types.Hash
	Quote       types.Hash
	Buyer       types.AccountID
	Seller      types.AccountID
	Maker       types.AccountID
	Taker       types.AccountID
	Type        OrderType
	Price       types.U128
	BaseAmount  types.U128
	QuoteAmount types.U128
}
