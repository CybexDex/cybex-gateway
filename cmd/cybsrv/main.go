package main

import "git.coding.net/bobxuyang/cy-gateway-BN/controllers/cybsrv"

func main() {
	// cybsrv.Test()
	go cybsrv.BlockRead()
	go cybsrv.HandleWorker()
	select {} // block forever
}
