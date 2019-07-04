package order

import (
	"cybex-gateway/model"
)

// IsBlack ...
func IsBlack(order1 *model.JPOrder) (bool, []*model.Black, error) {
	blockchain := order1.BlockChain
	addressFrom := order1.From
	addressTo := order1.To
	addressOut := order1.OutAddr
	user := order1.CybUser
	isblack, bs, err := model.BlackWith(user, blockchain, addressFrom, addressTo, addressOut)
	return isblack, bs, err
}
