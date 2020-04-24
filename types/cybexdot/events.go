package types

import "github.com/centrifuge/go-substrate-rpc-client/types"

// Events holds thee default events and custom events
type Events struct {
	types.EventRecords
	TokenModule_Issued    []EventTokenIssued    //nolint:stylecheck,golint
	TokenModule_Transferd []EventTokenTransferd //nolint:stylecheck,golint
	TokenModule_Freezed   []EventTokenFreezed   //nolint:stylecheck,golint
	TokenModule_UnFreezed []EventTokenUnFreezed //nolint:stylecheck,golint

	TradeModule_TradePairCreated []EventTradePairCreated //nolint:stylecheck,golint
	TradeModule_OrderCreated     []EventOrderCreated     //nolint:stylecheck,golint
	TradeModule_TradeCreated     []EventTradeCreated     //nolint:stylecheck,golint
	TradeModule_OrderCanceled    []EventOrderCanceled    //nolint:stylecheck,golint
}

type EventTokenIssued struct {
	Phase     types.Phase
	From      types.AccountID
	TokenHash types.Hash
	Supply    types.U128
	Topics    []types.Hash
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

type EventTokenUnFreezed struct {
	Phase     types.Phase
	AccountID types.AccountID
	TokenHash types.Hash
	Balance   types.U128
	Topics    []types.Hash
}

type EventTradePairCreated struct {
	Phase         types.Phase
	AccountID     types.AccountID
	TradePairHash types.Hash
	Pair          TradePair
	Topics        []types.Hash
}

type EventTradeCreated struct {
	Phase          types.Phase
	AccountID      types.AccountID
	BaseTokenHash  types.Hash
	QuoteTokenHash types.Hash
	OrderHash      types.Hash
	Trade          Trade
	Topics         []types.Hash
}

type EventOrderCreated struct {
	Phase          types.Phase
	AccountID      types.AccountID
	BaseTokenHash  types.Hash
	QuoteTokenHash types.Hash
	OrderHash      types.Hash
	Order          LimitOrder
	Topics         []types.Hash
}

type EventOrderCanceled struct {
	Phase     types.Phase
	AccountID types.AccountID
	OrderHash types.Hash
	Topics    []types.Hash
}
