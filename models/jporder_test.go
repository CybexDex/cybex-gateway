package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init () {
	db := GetDB()
	db.Exec("DELETE from jp_orders")
	db.Exec("DELETE from orders")
	db.Exec("DELETE from balances")
}

func TestJPOrderOne(t *testing.T) {
	//	test case 1
	// 		send "DONE" jporder
	// 		create "DONE" jporder, calculate TotalAmount, Amount and Fee of jporder
	// 		add InLocked & InLockedFee to balance
	// 		create "INIT" order, set TotalAmount, Amount and Fee of jporder

	err := db.DB().Ping()
	assert.Nil(t, err)

	// assert equality
	assert.Equal(t, 123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(t, 123, 456, "they should not be equal")
}

func TestOne(t *testing.T) {
	//	test case 1
	// 		send "DONE" jporder
	// 		create "DONE" jporder, calculate TotalAmount, Amount and Fee of jporder
	// 		add InLocked & InLockedFee to balance
	// 		create "INIT" order, set TotalAmount, Amount and Fee of jporder

	//	test case 2
	//		send "FAILED" jporder
	// 		create "FAILED" jporder, calculate TotalAmount, Amount and Fee of jporder
	//		add no data to balance, don't create order

	// test case 3
	//		send "PENDING" jporder
	// 		create "PENDING" jporder, calculate TotalAmount, Amount and Fee of jporder
	// 		add InLocked & InLockedFee to balance
	//		send "PENDING" jporder
	//		add no data to balance, don't create order
	// 		send "DONE" jporder
	// 		set jporder status to "DONE"
	// 		add no data to balance
	// 		create "INIT" order, set TotalAmount, Amount and Fee of jporder

	// test case 4
	//		send "PENDING" jporder
	// 		create "PENDING" jporder, calculate TotalAmount, Amount and Fee of jporder
	// 		add InLocked & InLockedFee to balance
	//		send "PENDING" jporder
	//		add no data to balance, don't create order
	// 		send "FAILED" jporder
	// 		set jporder status to "FAILED"
	// 		REVERT balance
	// 		don't create order
}
