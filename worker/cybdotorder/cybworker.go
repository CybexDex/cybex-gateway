package cybdotorder

import (
	"cybex-gateway/model"
	"cybex-gateway/utils/log"
	"encoding/hex"
	"fmt"
	"math/big"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	gsTypes "github.com/centrifuge/go-substrate-rpc-client/types"

	"github.com/centrifuge/go-substrate-rpc-client/signature"
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

func HandleWorker(seconds int) {
	for {
		for {
			ret := HandleDepositOneTime()
			if ret != 0 {
				break
			}
		}

		time.Sleep(time.Second * time.Duration(seconds))
	}
}

func HandleDepositOneTime() int {
	order1, _ := HoldOne()
	if order1.ID == 0 {
		return 1
	}
	if order1.Current == "cyborder" {
		handleOrders(order1)
	}
	order1.Save()
	return 0
	// check order to process it
}

// HoldOne ...
func HoldOne() (*model.JPOrder, error) {
	order, err := model.HoldCYBOrderOne()
	return order, err
}

func handleOrders(order *model.JPOrder) (err error) {
	// 是否可处理的asset
	log.Infoln("start handle order", order.ID)
	assetC, err := model.AssetsFind(order.Asset)
	if err != nil {
		log.Errorln("AssetsFind", err)
		return fmt.Errorf("AssetsFind %v", err)
	}
	gatewayAccount := assetC.GatewayAccount

	if gatewayAccount == order.CybUser {
		order.SetCurrent("cyborder", model.JPOrderStatusTerminate, "网关账号不能发给自己")
		return nil
	}

	// TODO:
	if viper.GetBool("cybserver.sendMemo") {
	}

	amount, err := strconv.ParseInt(order.Amount.String(), 10, 64)
	if err != nil {
		log.Errorln(err)
	}
	hash, err := gsTypes.NewHashFromHexString(assetC.CYBID)
	if err != nil {
		log.Errorln(err)
	}
	extrinsic, err := CreateTransfer(order.CybUser, amount, hash, "")
	if err != nil {
		log.Errorln(err)
	}

	txHash, err := SignAndSendTransfer(extrinsic, assetC.GatewayPass)
	if err != nil {
		log.Errorln(err)
		order.SetCurrent("cyborder", model.JPOrderStatusFailed, "send error")
		errmsg := fmt.Sprintf("id:%d\nerr:%v", order.ID, err)
		model.WxSendTaskCreate("cybex充值失败", errmsg)
	}

	log.Infof("order:%d,%s:%+v\n", order.ID, "sendorder", txHash)
	order.SetCurrent("cyborder", model.JPOrderStatusPending, "")
	return nil
}

func CreateTransfer(toAccountAddress string, amount int64, tokenHash types.Hash, memo string) (*gsTypes.Extrinsic, error) {
	keyringPair, err := keyringPairFromSecret(toAccountAddress)
	if err != nil {
		return nil, err
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	var ext gsTypes.Extrinsic
	if (tokenHash == types.Hash{}) {
		to, err := gsTypes.NewAddressFromHexAccountID(hexutil.Encode(keyringPair.PublicKey[:]))

		if err != nil {
			return nil, err
		}

		c, err := gsTypes.NewCall(meta, "Balance.transfer", to, gsTypes.UCompact(amount))
		if err != nil {
			return nil, err
		}

		// Create the extrinsic
		ext = gsTypes.NewExtrinsic(c)
	} else {
		to := gsTypes.NewAccountID(keyringPair.PublicKey[:])
		var memoBytes [8]byte
		copy(memoBytes[:], memo)
		c, err := gsTypes.NewCall(meta, "TokenModule.transfer", tokenHash, to, gsTypes.NewU128(*big.NewInt(amount)), gsTypes.NewOptionBytes(gsTypes.NewBytes([]byte(memo))))

		if err != nil {
			return nil, err
		}

		// Create the extrinsic
		ext = gsTypes.NewExtrinsic(c)
	}

	return &ext, nil
}

func SignAndSendTransfer(ext *gsTypes.Extrinsic, seedOrPhrase string) (string, error) {
	keyringPair, err := keyringPairFromSecret(seedOrPhrase)
	if err != nil {
		return "", err
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return "", err
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return "", err
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return "", err
	}

	key, err := gsTypes.CreateStorageKey(meta, "System", "AccountNonce", keyringPair.PublicKey, nil)
	if err != nil {
		return "", err
	}

	var nonce uint32
	err = api.RPC.State.GetStorageLatest(key, &nonce)
	if err != nil {
		return "", err
	}

	o := gsTypes.SignatureOptions{
		BlockHash:   genesisHash,
		Era:         gsTypes.ExtrinsicEra{IsMortalEra: false},
		GenesisHash: genesisHash,
		Nonce:       gsTypes.UCompact(nonce),
		SpecVersion: rv.SpecVersion,
		Tip:         0,
	}
	err = ext.Sign(keyringPair, o)
	if err != nil {
		return "", err
	}

	// Send the extrinsic
	hash, err := api.RPC.Author.SubmitExtrinsic(*ext)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(hash[:]), nil
}

var rePubKey = regexp.MustCompile(`Public key \(hex\):\s+0x([a-f0-9]*)\n`)
var reAddressNew = regexp.MustCompile(`SS58 Address:\s+([a-zA-Z0-9]*)\n`)

// in: mnemonic,seed,pubkey uri,address
func keyringPairFromSecret(in string) (signature.KeyringPair, error) {
	cmd := exec.Command("subkey", "inspect", in)

	// execute the command, get the output
	out, err := cmd.Output()
	if err != nil {
		return signature.KeyringPair{}, fmt.Errorf("failed to generate keyring pair from secret: %v", err.Error())
	}

	if string(out) == "Invalid phrase/URI given" {
		return signature.KeyringPair{}, fmt.Errorf("failed to generate keyring pair from secret: invalid phrase/URI given")
	}

	// find the pub key
	resPk := rePubKey.FindStringSubmatch(string(out))
	if len(resPk) != 2 {
		return signature.KeyringPair{}, fmt.Errorf("failed to generate keyring pair from secret, pubkey not found in output: %v", resPk)
	}
	pk, err := hex.DecodeString(resPk[1])
	if err != nil {
		return signature.KeyringPair{}, fmt.Errorf("failed to generate keyring pair from secret, could not hex decode pubkey: %v",
			resPk[1])
	}

	// find the address
	addr := reAddressNew.FindStringSubmatch(string(out))

	if len(addr) != 2 {
		return signature.KeyringPair{}, fmt.Errorf("failed to generate keyring pair from secret, address not found in output: %v", addr)
	}

	return signature.KeyringPair{
		URI:       in,
		Address:   addr[1],
		PublicKey: pk,
	}, nil
}
