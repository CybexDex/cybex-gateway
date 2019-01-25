package address

import (
	"fmt"
	"testing"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

func TestOne(t *testing.T) {
	repo := NewRepo(m.GetDB())

	// fetch all
	res, _ := repo.FetchAll()
	fmt.Println(res)
	fmt.Println(*(res[0]))

	a, _ := repo.GetByID(1)
	fmt.Println(a.Asset)
}
