package utils

import (
	"fmt"
	"strings"

	apim "github.com/CybexDex/cybex-go/api"
	"github.com/CybexDex/cybex-go/config"
	"github.com/CybexDex/cybex-go/crypto"
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

// KeyBagByUserSeedPass ...
func KeyBagByUserSeedPass(user string, pass string) (keybag *crypto.KeyBag) {
	passArr := strings.Split(pass, ",")
	lenpass := len(passArr)
	if lenpass == 1 {
		gatewaypass := SeedString(pass)
		keybag = apim.KeyBagByUserPass(user, gatewaypass)
	} else if lenpass >= 2 {
		keybag = crypto.NewKeyBag()
		for _, seedkey := range passArr {
			newkey := SeedString(seedkey)
			keybag.Add(newkey)
		}
	}
	return keybag
}
