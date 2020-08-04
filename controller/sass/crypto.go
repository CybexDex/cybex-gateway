package sass

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// SignHMACSHA256 ...
func SignHMACSHA256(data interface{}, secret string) (string, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(bytes.NewReader(buf))
	decoder.UseNumber()
	obj := make(map[string]interface{})
	err = decoder.Decode(&obj)
	if err != nil {
		return "", err
	}

	msgStr := BuildMsg(obj, "=", "&")
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msgStr))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha, nil
}

// BuildMsg ...
func BuildMsg(val interface{}, keyValSeparator, groupSeparator string) string {
	if val == nil {
		return ""
	}

	msg := ""
	switch reflect.TypeOf(val).Kind() {
	case reflect.Struct:
		buf, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		decoder := json.NewDecoder(bytes.NewReader(buf))
		decoder.UseNumber()
		m := make(map[string]interface{})
		err = decoder.Decode(&m)
		if err != nil {
			return ""
		}
		msg = BuildMsg(m, keyValSeparator, groupSeparator)
	case reflect.Map:
		buf, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		decoder := json.NewDecoder(bytes.NewReader(buf))
		decoder.UseNumber()
		obj := make(map[string]interface{})
		err = decoder.Decode(&obj)

		keyVals := make(map[string]string)
		keys := make([]string, 0, len(obj))

		for k, v := range obj {
			_msg := BuildMsg(v, keyValSeparator, groupSeparator)
			keyVals[k] = _msg
			keys = append(keys, k)
		}
		sort.Strings(keys)
		groupStrs := make([]string, 0, len(keys))
		for _, key := range keys {
			groupStrs = append(groupStrs, key+keyValSeparator+keyVals[key])
		}
		msg = strings.Join(groupStrs, groupSeparator)
	case reflect.Slice:
		arr := val.([]interface{})
		keyVals := make(map[string]string)
		keys := make([]string, 0, len(arr))

		for i, v := range arr {
			key := strconv.Itoa(i)
			keys = append(keys, key)
			keyVals[key] = BuildMsg(v, keyValSeparator, groupSeparator)
		}
		sort.Strings(keys)

		groupStrs := make([]string, 0, len(keys))
		for _, key := range keys {
			groupStrs = append(groupStrs, key+keyValSeparator+keyVals[key])
		}
		msg = strings.Join(groupStrs, groupSeparator)
	default:
		msg = fmt.Sprintf("%v", val)
	}

	return msg
}
