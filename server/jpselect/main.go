package jpselect

import (
	"cybex-gateway/server/jp"
	"cybex-gateway/server/sass"

	"github.com/spf13/viper"
)

// StartServer ...
func StartServer() {
	isSass := viper.GetBool("useSass")
	if isSass {
		sass.StartServer()
	} else {
		jp.StartServer()
	}
}
