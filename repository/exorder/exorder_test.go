package exorder

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
	//fmt.Println(*(res[0]))

	// fetch by pagination
	res, _ = repo.Fetch(p.Page{
		Offset:  0,
		Amount:  1,
		OrderBy: "ID",
		Sort:    "asc",
	})
	fmt.Println(res)

	o := m.ExOrder{
		Status: "DONE",
	}
	err := repo.Update((*(res[0])).ID, &o)
	fmt.Println(err)

	// delete by ID
	// fmt.Println(repo.DeleteByID(4))
}
