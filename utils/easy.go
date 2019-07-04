package utils

import (
	"cybex-gateway/utils/seed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// ResultToStruct ...
func ResultToStruct(input interface{}, output interface{}) (err error) {
	err = mapstructure.Decode(input, output)
	return err
}

// SeedString ...
func SeedString(nowstring string) string {
	if strings.Index(nowstring, "seed__") == 0 {
		s := strings.TrimPrefix(nowstring, "seed__")
		sdata := seed.GetSeedData(s)
		return sdata
	}
	return nowstring
}

// ErrorAdd ...
func ErrorAdd(errin error, msg string) error {
	return fmt.Errorf("%s,%v", msg, errin)
}

// V2S ...
func V2S(v interface{}, s interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	return nil
}
