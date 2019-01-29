package singleton

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	address "git.coding.net/bobxuyang/cy-gateway-BN/repository/address"
	app "git.coding.net/bobxuyang/cy-gateway-BN/repository/app"
	asset "git.coding.net/bobxuyang/cy-gateway-BN/repository/asset"
	cyborder "git.coding.net/bobxuyang/cy-gateway-BN/repository/cyborder"
	order "git.coding.net/bobxuyang/cy-gateway-BN/repository/order"
)

// App ...
var App app.Repository

// Address ...
var Address address.Repository

// Asset ...
var Asset asset.Repository

// Order ...
var Order order.Repository

// CybOrder ...
var CybOrder cyborder.Repository

func init() {
	App = app.NewRepo(m.GetDB())
	Address = address.NewRepo(m.GetDB())
	Asset = asset.NewRepo(m.GetDB())
	Order = order.NewRepo(m.GetDB())
	CybOrder = cyborder.NewRepo(m.GetDB())
}
