package main

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
)

func main() {
	priKeyBytes, err := hex.DecodeString("bf12996feeaa2977b6ca0d33a0e8bd2ccfc4844c6f8a7e6d15c099f8da4a255d")
	if err != nil {
		fmt.Println(err)
	}
	_, pub := btcec.PrivKeyFromBytes(btcec.S256(), priKeyBytes)
	if err != nil {
		fmt.Println(err)
	}
	hexStr := hex.EncodeToString(pub.SerializeCompressed())
	fmt.Println(hexStr)
}
