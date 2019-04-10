package jp

import (
	"bitbucket.org/woyoutlz/bbb-gateway/types"
)

// HandleWithdraw ...
func HandleWithdraw(result types.JPOrderResult) error {
	// 事务
	// 更新提现订单
	jporderID := result.ID
	// 寻找提现订单
	// jporder.update
	// 告知record
	// record.afterJPWithdraw
	// 记录Done事件
	// 或者抛出错误
	return nil
}

// HandleDeposit ...
func HandleDeposit(result types.JPOrderResult) error {
	// 事务
	// 新建充值订单
	// jporder.create
	// 或者更新充值订单
	// jporder.update
	// 告知record
	// record.afterJPDeposit
	// 记录Done事件
	// 或者抛出错误
	return nil
}
