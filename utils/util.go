package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/btcsuite/btcd/btcec"
)

//Message ...
func Message(status bool, message string, d ...interface{}) map[string]interface{} {
	if len(d) == 1 {
		return map[string]interface{}{"status": status, "message": message, "data": d[0]}
	}

	return map[string]interface{}{"status": status, "message": message}
}

//Respond ...
func Respond(w http.ResponseWriter, data map[string]interface{}, s ...int) {
	status := http.StatusOK
	if len(s) == 1 {
		status = s[0]
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
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
		keys := make([]string, len(obj))

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
		for i, v := range arr {
			itemMsg := BuildMsg(v)
			msg += fmt.Sprintf("%d%s", i, itemMsg)
		}
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
