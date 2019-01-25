package ordersrv

import (
	"database/sql"
	"fmt"

	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

func findOrders() (*sql.Rows, error) {
	rows, _ := rep.Order.MDB().Where(&m.Order{
		Status: "INIT",
	}).Rows() // (*sql.Rows, error)
	return rows, nil
}
func handleOrders(rows *sql.Rows) {
	fmt.Println(rows)
	for rows.Next() {
		var order1 = &m.Order{}
		m.GetDB().ScanRows(rows, order1)
		// do something
		fmt.Println(order1)
	}
	fmt.Println("finished")

}

// HandleOneTime ...
func HandleOneTime() {
	fmt.Println("handle one")
	// check order to process it
	orders, err := findOrders()
	defer orders.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	// //handle orders
	handleOrders(orders)
}
func WorkerStart() {

}
func WorkerStop() {

}
