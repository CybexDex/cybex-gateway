package jp

import (
	"testing"
)

func TestHandleDeposit(t *testing.T) {
	t.Log(1, 2)
}
func TestDepositAddress(t *testing.T) {
	re, _ := DepositAddress("ETH")
	t.Log(re)
}
