package jpselect

import (
	"cybex-gateway/worker/jp"
	"cybex-gateway/worker/sass"

	"github.com/spf13/viper"
)

// HandleWorker ...
func HandleWorker(seconds int) {
	isSass := viper.GetBool("useSass")
	if isSass {
		sass.HandleWorker(5)
	} else {
		jp.HandleWorker(5)
	}
}
