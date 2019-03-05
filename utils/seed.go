package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

var cmdkey string
var server string

func init() {
	InitConfig()
	cmdkey = viper.GetString("cybsrv.seedCMDKEY")
	server = viper.GetString("cybsrv.seedServer")
}

// GetSeedData ...
func GetSeedData(name string) string {
	resp, err := http.Get(server + "/data/" + name)
	if err != nil {
		Errorln("Unable to make request GetSeedData: ", name, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	result := KeyDecrypt(cmdkey, bodyString)
	return result
}

//KeyDecrypt ...
func KeyDecrypt(keyStr string, cryptoText string) string {
	keyBytes := sha256.Sum256([]byte(keyStr))
	return decrypt(keyBytes[:], cryptoText)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return fmt.Sprintf("%s", ciphertext)
}
