package types

import (
	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

type LimitOrder struct {
	Hash               types.Hash
	Base               types.Hash
	Quote              types.Hash
	Owner              types.AccountID
	Price              types.U128
	SellAmount         types.U128
	BuyAmount          types.U128
	RemainedSellAmount types.U128
	RemainedBuyAmount  types.U128
	Type               OrderType
	Status             OrderStatus
}

type OrderType struct {
	Buy  types.Bool
	Sell types.Bool
}

func (m *OrderType) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	if b == 0 {
		m.Buy = true
	} else if b == 1 {
		m.Sell = true
	}

	return nil
}

func (m OrderType) Encode(encoder scale.Encoder) error {
	var err1 error
	if m.Buy {
		err1 = encoder.PushByte(0)
	} else if m.Sell {
		err1 = encoder.PushByte(1)
	}

	if err1 != nil {
		return err1
	}

	return nil
}

type OrderStatus struct {
	Created       types.Bool
	PartialFilled types.Bool
	Filled        types.Bool
	Canceled      types.Bool
}

func (m *OrderStatus) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	if b == 0 {
		m.Created = true
	} else if b == 1 {
		m.PartialFilled = true
	} else if b == 2 {
		m.Filled = true
	} else if b == 3 {
		m.Canceled = true
	}

	return nil
}

func (m OrderStatus) Encode(encoder scale.Encoder) error {
	var err1 error
	if m.Created {
		err1 = encoder.PushByte(0)
	} else if m.PartialFilled {
		err1 = encoder.PushByte(1)
	} else if m.Filled {
		err1 = encoder.PushByte(2)
	} else if m.Canceled {
		err1 = encoder.PushByte(3)
	}

	if err1 != nil {
		return err1
	}

	return nil
}
