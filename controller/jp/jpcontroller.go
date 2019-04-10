package jp

import (
	"fmt"

	"bitbucket.org/woyoutlz/bbb-gateway/types"
)

// HandleWithdraw ...
func HandleWithdraw(result types.JPOrderResult) error {
	fmt.Println(2333, result)
	return nil
}

// HandleDeposit ...
func HandleDeposit(result types.JPOrderResult) error {
	return nil
}
