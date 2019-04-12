package order

import (
	"bitbucket.org/woyoutlz/bbb-gateway/model"
)

// IsBlack ...
func IsBlack(order1 *model.JPOrder) (bool, []*model.Black, error) {
	blockchain := order1.BlockChain
	addressFrom := order1.From
	addressTo := order1.To
	user := order1.CybUser
	isblack, bs, err := model.BlackWith(user, blockchain, addressFrom, addressTo)
	return isblack, bs, err
}
