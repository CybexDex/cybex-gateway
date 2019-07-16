package ecc

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/sha3"
)

// ECCSig ...
type ECCSig struct {
	R string `json:"r"`
	S string `json:"s"`
	V int64  `json:"v"`
}

// SignECCData ...
func SignECCData(prikey string, data interface{}) (*ECCSig, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(buf))
	decoder.UseNumber()
	obj := make(map[string]interface{})
	err = decoder.Decode(&obj)
	if err != nil {
		return nil, err
	}

	priKeyBytes, err := hex.DecodeString(prikey)
	if err != nil {
		return nil, err
	}
	priKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), priKeyBytes)
	if err != nil {
		return nil, err
	}

	msgStr := BuildMsg(obj)
	sha3Hash := sha3.NewLegacyKeccak256()
	_, err = sha3Hash.Write([]byte(msgStr))
	if err != nil {
		return nil, err
	}
	msgBuf := sha3Hash.Sum(nil)
	sig, err := priKey.Sign(msgBuf)
	if err != nil {
		return nil, err
	}
	var R [32]byte
	var S [32]byte
	copy(R[32-len(sig.R.Bytes()):], sig.R.Bytes())
	copy(S[32-len(sig.S.Bytes()):], sig.S.Bytes())
	// _sig := &ECCSig{
	// 	R: hex.EncodeToString(R[:]),
	// 	S: hex.EncodeToString(S[:]),
	// }
	_sig := &ECCSig{
		R: base64.StdEncoding.EncodeToString(R[:]),
		S: base64.StdEncoding.EncodeToString(S[:]),
	}
	return _sig, nil
}
func TestECCSign() {
	pubkey := "04ace32532c90652e1bae916248e427a7ab10aeeea1067949669a3f4da10965ef90d7297f538f23006a31f94fdcfaed9e8dd38c85ba7e285f727430332925aefe5"
	pubKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {

	}
	fmt.Println(42, pubkey)
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	fmt.Println(1, pubKey, err)
}

// VerifyECCSign ...
func VerifyECCSign(data interface{}, sign *ECCSig, pubkey string) (bool, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	decoder := json.NewDecoder(bytes.NewReader(buf))
	decoder.UseNumber()
	obj := make(map[string]interface{})
	err = decoder.Decode(&obj)
	if err != nil {
		return false, err
	}
	pubKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		return false, err
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return false, err
	}
	msgStr := BuildMsg(obj)
	sha3Hash := sha3.NewLegacyKeccak256()
	_, err = sha3Hash.Write([]byte(msgStr))
	if err != nil {
		return false, err
	}
	msgBuf := sha3Hash.Sum(nil)

	// Decode hex-encoded serialized signature.
	decodedR, err := base64.StdEncoding.DecodeString(sign.R)
	if err != nil {
		return false, err
	}
	decodedS, err := base64.StdEncoding.DecodeString(sign.S)
	if err != nil {
		return false, err
	}
	signature := btcec.Signature{
		R: new(big.Int).SetBytes(decodedR),
		S: new(big.Int).SetBytes(decodedS),
	}
	return signature.Verify(msgBuf, pubKey), nil
}

// BuildMsg ...
func BuildMsg(val interface{}) string {
	if val == nil {
		return ""
	}

	msg := ""
	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		obj := val.(map[string]interface{})
		keyVals := make(map[string]string)
		keys := make([]string, 0, len(obj))

		for k, v := range obj {
			_msg := BuildMsg(v)
			keyVals[k] = _msg
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			msg += key + keyVals[key]
		}
	case reflect.Slice:
		arr := val.([]interface{})
		keyVals := make(map[string]string)
		keys := make([]string, 0, len(arr))

		for i, v := range arr {
			key := strconv.Itoa(i)
			keys = append(keys, key)
			keyVals[key] = BuildMsg(v)
		}
		sort.Strings(keys)

		groupStrs := make([]string, 0, len(keys))
		for _, key := range keys {
			groupStrs = append(groupStrs, key+keyVals[key])
		}
		msg += strings.Join(groupStrs, "")

	default:
		msg = fmt.Sprintf("%v", val)
	}

	return msg
}

// PriToPub ...
func PriToPub(prikey string) string {
	pkBytes, err := hex.DecodeString(prikey)
	if err != nil {
		return ""
	}
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	return hex.EncodeToString(pubKey.SerializeCompressed())
}

// NewPriPub ...
func NewPriPub() (string, string) {
	pri, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", ""
	}
	return hex.EncodeToString(pri.Serialize()), hex.EncodeToString(pri.PubKey().SerializeCompressed())
}
