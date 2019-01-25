package ordersrv

import (
	"database/sql"
	"fmt"

	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
)

func findOrders() (*sql.Rows, error) {
	rows, err := rep.Order.Rows(&m.Order{
		Status: "INIT",
	})
	return rows, err
}
func handleOrders(rows *sql.Rows) {
	for rows.Next() {
		var order = &m.Order{}
		m.GetDB().ScanRows(rows, order)
		// do something
		fmt.Println(order)
	}
	fmt.Println("finished")

}

// HandleOneTime ...
func HandleOneTime() {
	fmt.Println("handle one")
	// check order to process it
	orders, err := findOrders()
	if err != nil {
		fmt.Println(err)
		return
	}
	//handle orders
	handleOrders(orders)
}
func WorkerStart() {

}
func WorkerStop() {

}
