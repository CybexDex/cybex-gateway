package main

import (
	"cybex-gateway/config"
	"cybex-gateway/controller/jp"
	"cybex-gateway/controller/sass"
	"fmt"
)

func init() {
	config.LoadConfig("uat")
}
func address() {
	// config.LoadConfig("uat")
	s, e := jp.DepositAddress("XTZ")
	fmt.Println(s, e)
	s, e = sass.DepositAddress("ETH")
	fmt.Println(s, e)
}
func verify() {
	s, e := jp.VerifyAddress("XTZ", "tz1bzLpsQsQMrCKEMmsoJjwZgotQ2N9y3CA9")
	fmt.Println(s, e)
	s, e = sass.VerifyAddress("ETH", "0xA703E5727929f90CB7d872BF4520AD68368CDd4F")
	fmt.Println(s, e)
}
func withdraw() {
	s, e := jp.Withdraw("XTZ", "tz1bzLpsQsQMrCKEMmsoJjwZgotQ2N9y3CA9", "0.01", 500)
	fmt.Println(s, e)
	// s, e = sass.Withdraw("ETH", "0xA703E5727929f90CB7d872BF4520AD68368CDd4F", "0.01", 1001)
	s, e = sass.Withdraw("ETH", "0xB0Da211fDd0c63d33493E14E05A7C06522fb1390", "1.0787", 5942200)
	fmt.Println(s, e)
}
func main() {
	// address()
	// verify()
	withdraw()
}
