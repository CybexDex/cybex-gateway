package jpselect

import (
	"cybex-gateway/controller/jp"
	"cybex-gateway/controller/sass"
	"cybex-gateway/types"

	"github.com/spf13/viper"
)

// DepositAddress ...
func DepositAddress(coin string) (address *types.JPAddressResult, err error) {
	isSass := viper.GetBool("useSass")
	if isSass {
		return sass.DepositAddress(coin)
	} else {
		return jp.DepositAddress(coin)
	}
}

// VerifyAddress ...
func VerifyAddress(asset string, address string) (res *types.VerifyRes, err error) {
	isSass := viper.GetBool("useSass")
	if isSass {
		return sass.VerifyAddress(asset, address)
	} else {
		return jp.VerifyAddress(asset, address)
	}
}
