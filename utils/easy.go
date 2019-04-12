package utils

import (
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
		return s
	}
	return nowstring
}

// ErrorAdd ...
func ErrorAdd(errin error, msg string) error {
	return fmt.Errorf("%s,%v", msg, errin)
}
