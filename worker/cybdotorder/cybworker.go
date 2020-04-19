package cybdotorder

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/spf13/viper"
)

var api *gsrpc.SubstrateAPI

// InitNode ...
func InitNode() {
	node := viper.GetString("cybserver.node")
	gapi, err := gsrpc.NewSubstrateAPI(node)
	api = gapi
	if err != nil {
		panic(err)
	}

}
