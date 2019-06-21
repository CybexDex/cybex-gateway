package main

import (
	"fmt"

	apim "github.com/CybexDex/cybex-go/api"
	"github.com/CybexDex/cybex-go/config"
	"github.com/CybexDex/cybex-go/types"
)

//CybSign ...
func CybSign(user string, pubkey string, password string, tosign string) (string, error) {
	api := apim.New("node", "")
	config.SetCurrentConfig("111")
	keyBag := apim.KeyBagByUserPass(user, password)
	pubs := keyBag.Publics()
	fmt.Println(pubs)
	pubA, _ := types.NewPublicKeyFromString(pubkey)
	pubkeys := types.PublicKeys{*pubA}
	priKeyA := keyBag.PrivatesByPublics(pubkeys)
	if len(priKeyA) == 0 {
		fmt.Println("fail")
		return "", fmt.Errorf("cannot find prikey")
	}
	wif := priKeyA[0].ToWIF()
	s, err := api.SignStr(tosign, wif)
	// x, err := api.VerifySign(user, tosign, s)
	// fmt.Println(x, err)
	return s, err
}

func main() {
	// timestamp := time.Now().Unix()
	timestamp := 1561100589
	fmt.Println(timestamp)
	user := "xiao01"
	password := "xiaoyongyong"
	pubkey := "CYB7rd2zuyyaiwhTbWK5JSjn9kuWQvhKBSFuKjomnSPjvYXXR49Yw"
	to := "1" //fmt.Sprintf("%d%s", timestamp, user)

	str, err := CybSign(user, pubkey, password, to)
	fmt.Println(str, err)
	api := apim.New("ws://18.136.140.223:38090/", "")
	if err := api.Connect(); err != nil {
		panic(err)
	}
	x, err := api.VerifySign(user, to, str)
	fmt.Println(x, err)
}

// func verify() {
// 	VerifySign("xiao01", toSign, sig)
// }
