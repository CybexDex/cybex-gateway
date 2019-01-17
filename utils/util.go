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

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// BuildMsg 将对象按照key字符顺序序列化字符串
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

// PriToPub 私钥串导出公钥串
func PriToPub(prikey string) string {
	pkBytes, err := hex.DecodeString(prikey)
	if err != nil {
		return ""
	}
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	return string(pubKey.SerializeCompressed())
}
