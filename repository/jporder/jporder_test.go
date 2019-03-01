package jporder

import (
	"fmt"
	"testing"

	m "coding.net/bobxuyang/cy-gateway-BN/models"
	p "coding.net/bobxuyang/cy-gateway-BN/repository"
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

	o := m.JPOrder{
		Status: "PENDING",
	}
	order := *(res[0])
	err := order.UpdateColumns(&o)
	fmt.Println(err)

	eo, _ := repo.GetByJPID(o.JadepoolID)
	fmt.Println(*eo)

	o = m.JPOrder{
		Status: "PENDING",
	}
	os, _ := repo.FetchWith(&o)
	fmt.Println(*(os[0]))

	// delete by ID
	// fmt.Println(repo.DeleteByID(4))
}
