package singleton

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	blockchain "git.coding.net/bobxuyang/cy-gateway-BN/repository/Blockchain"
	address "git.coding.net/bobxuyang/cy-gateway-BN/repository/address"
	app "git.coding.net/bobxuyang/cy-gateway-BN/repository/app"
	asset "git.coding.net/bobxuyang/cy-gateway-BN/repository/asset"
	cyborder "git.coding.net/bobxuyang/cy-gateway-BN/repository/cyborder"
	jporder "git.coding.net/bobxuyang/cy-gateway-BN/repository/jporder"
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

// Blockchain ...
var Blockchain blockchain.Repository

// JPOrder ...
var JPOrder jporder.Repository

func init() {
	App = app.NewRepo(m.GetDB())
	Address = address.NewRepo(m.GetDB())
	Asset = asset.NewRepo(m.GetDB())
	Order = order.NewRepo(m.GetDB())
	CybOrder = cyborder.NewRepo(m.GetDB())
	Blockchain = blockchain.NewRepo(m.GetDB())
	JPOrder = jporder.NewRepo(m.GetDB())
}
