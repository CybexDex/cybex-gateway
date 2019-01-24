package singleton

import (
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	address "git.coding.net/bobxuyang/cy-gateway-BN/repository/address"
	app "git.coding.net/bobxuyang/cy-gateway-BN/repository/app"
	asset "git.coding.net/bobxuyang/cy-gateway-BN/repository/asset"
)

// App ...
var App app.Repository

// Address ...
var Address address.Repository

// Asset ...
var Asset asset.Repository

func init() {
	App = app.NewRepo(m.GetDB())
	Address = address.NewRepo(m.GetDB())
	Asset = asset.NewRepo(m.GetDB())
}
