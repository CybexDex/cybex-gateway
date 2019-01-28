package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"

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
	buf, _ := json.Marshal(data)
	decoder := json.NewDecoder(bytes.NewReader(buf))
	decoder.UseNumber()
	obj := make(map[string]interface{})
	err := decoder.Decode(&obj)
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

	_sig := &ECCSig{
		R: base64.StdEncoding.EncodeToString(sig.R.Bytes()),
		S: base64.StdEncoding.EncodeToString(sig.S.Bytes()),
	}
	return _sig, nil
}

// VerifyECCSign ...
func VerifyECCSign(data interface{}, sign *ECCSig, pubkey string) (bool, error) {
	buf, _ := json.Marshal(data)
	decoder := json.NewDecoder(bytes.NewReader(buf))
	decoder.UseNumber()
	obj := make(map[string]interface{})
	err := decoder.Decode(&obj)
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
