package cybdotorder

import (
	"encoding/hex"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/ethereum/go-ethereum/common/hexutil"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
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

func MakeTransfer(toAccountAddress string, amount uint64) (*gsTypes.Extrinsic, error) {
	keyringPair, err := keyringPairFromSecret(toAccountAddress)
	if err != nil {
		return nil, err
	}

	to, err := gsTypes.NewAddressFromHexAccountID(hexutil.Encode(keyringPair.PublicKey[:]))

	if err != nil {
		return nil, err
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	c, err := gsTypes.NewCall(meta, "Balances.transfer", to, gsTypes.UCompact(amount))
	if err != nil {
		return nil, err
	}

	// Create the extrinsic
	ext := gsTypes.NewExtrinsic(c)

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
