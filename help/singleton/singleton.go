package singleton

import (
	m "coding.net/bobxuyang/cy-gateway-BN/models"
	address "coding.net/bobxuyang/cy-gateway-BN/repository/address"
	app "coding.net/bobxuyang/cy-gateway-BN/repository/app"
	asset "coding.net/bobxuyang/cy-gateway-BN/repository/asset"
	"coding.net/bobxuyang/cy-gateway-BN/repository/bigasset"
	black "coding.net/bobxuyang/cy-gateway-BN/repository/black"
	blockchain "coding.net/bobxuyang/cy-gateway-BN/repository/blockchain"
	cyborder "coding.net/bobxuyang/cy-gateway-BN/repository/cyborder"
	"coding.net/bobxuyang/cy-gateway-BN/repository/cybtoken"
	"coding.net/bobxuyang/cy-gateway-BN/repository/easy"
	jporder "coding.net/bobxuyang/cy-gateway-BN/repository/jporder"
	"coding.net/bobxuyang/cy-gateway-BN/repository/order"
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

// Black ...
var Black black.Repository

// BigAsset ...
var BigAsset bigasset.Repository

// CybToken ...
var CybToken cybtoken.Repository

// Easy ...
var Easy easy.Repository

func Init() {
	App = app.NewRepo(m.GetDB())
	Address = address.NewRepo(m.GetDB())
	Asset = asset.NewRepo(m.GetDB())
	Order = order.NewRepo(m.GetDB())
	CybOrder = cyborder.NewRepo(m.GetDB())
	Blockchain = blockchain.NewRepo(m.GetDB())
	JPOrder = jporder.NewRepo(m.GetDB())
	Black = black.NewRepo(m.GetDB())
	BigAsset = bigasset.NewRepo(m.GetDB())
	CybToken = cybtoken.NewRepo(m.GetDB())
	Easy = easy.NewRepo(m.GetDB())
}
