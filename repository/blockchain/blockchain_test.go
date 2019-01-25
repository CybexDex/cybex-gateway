package blockchain

import (
	"fmt"
	"testing"

	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	p "git.coding.net/bobxuyang/cy-gateway-BN/repository"
)

func TestOne(t *testing.T) {
	repo := NewRepo(m.GetDB())

	// fetch all
	res, _ := repo.FetchAll()
	fmt.Println(res)
	fmt.Println(*(res[0]))
	fmt.Println(*(res[1]))

	// fetch by pagination
	res, _ = repo.Fetch(p.Page{
		Offset:  1,
		Amount:  1,
		OrderBy: "ID",
		Sort:    "asc",
	})
	fmt.Println(*(res[0]))

	// delete by ID
	//fmt.Println(repo.DeleteByID(4))
}
