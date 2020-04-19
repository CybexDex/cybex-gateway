package types

import "github.com/centrifuge/go-substrate-rpc-client/types"

// Events holds thee default events and custom events
type Events struct {
	types.EventRecords
	TokenModule_Transferd    []EventTokenTransferd    //nolint:stylecheck,golint
	TokenModule_Freezed      []EventTokenFreezed      //nolint:stylecheck,golint
	TradeModule_OrderCreated []EventTradeOrderCreated //nolint:stylecheck,golint
}

type EventTokenTransferd struct {
	Phase     types.Phase
	From      types.AccountID
	To        types.AccountID
	TokenHash types.Hash
	Value     types.U128
	Topics    []types.Hash
}

type EventTokenFreezed struct {
	Phase     types.Phase
	AccountID types.AccountID
	TokenHash types.Hash
	Balance   types.U128
	Topics    []types.Hash
}

type EventTradeOrderCreated struct {
	Phase          types.Phase
	AccountID      types.AccountID
	BaseTokenHash  types.Hash
	QuoteTokenHash types.Hash
	OrderHash      types.Hash
	Order          LimitOrder
	Topics         []types.Hash
}
