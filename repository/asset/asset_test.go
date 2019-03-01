package asset

import (
	"fmt"
	"testing"

	m "coding.net/bobxuyang/cy-gateway-BN/models"
)

func TestOne(t *testing.T) {
	repo := NewRepo(m.GetDB())

	// fetch all
	res, _ := repo.FetchAll()
	fmt.Println(res)
	fmt.Println(*(res[0]))

	var asset m.Asset
	m.GetDB().Model(&m.Asset{}).Where("ID=?", 1).Preload("Blockchain").Find(&asset)
	fmt.Println(asset)
}
