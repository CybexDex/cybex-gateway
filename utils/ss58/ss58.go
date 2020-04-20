package ss58

import (
	"cybex-gateway/utils"
	"cybex-gateway/utils/base58"

	"golang.org/x/crypto/blake2b"
)

func Decode(address string) string {
	checksumPrefix := []byte("SS58PRE")
	ss58Format := base58.Decode(address)
	if ss58Format[0] != 42 {
		return ""
	}
	var checksumLength int
	if utils.IntInSlice(len(ss58Format), []int{3, 4, 6, 10}) {
		checksumLength = 1
	} else if utils.IntInSlice(len(ss58Format), []int{5, 7, 11, 35}) {
		checksumLength = 2
	} else if utils.IntInSlice(len(ss58Format), []int{8, 12}) {
		checksumLength = 3
	} else if utils.IntInSlice(len(ss58Format), []int{9, 13}) {
		checksumLength = 4
	} else if utils.IntInSlice(len(ss58Format), []int{14}) {
		checksumLength = 5
	} else if utils.IntInSlice(len(ss58Format), []int{15}) {
		checksumLength = 6
	} else if utils.IntInSlice(len(ss58Format), []int{16}) {
		checksumLength = 7
	} else if utils.IntInSlice(len(ss58Format), []int{17}) {
		checksumLength = 8
	} else {
		return ""
	}
	bss := ss58Format[0 : len(ss58Format)-checksumLength]
	checksum, _ := blake2b.New(64, []byte{})
	w := append(checksumPrefix[:], bss[:]...)
	checksum.Write(w)
	h := checksum.Sum(nil)
	if utils.BytesToHex(h[0:checksumLength]) != utils.BytesToHex(ss58Format[len(ss58Format)-checksumLength:]) {
		return ""
	}
	return utils.BytesToHex(ss58Format[1:33])
}

func Encode(address string) string {
	checksumPrefix := []byte("SS58PRE")
	addressBytes := utils.HexToBytes(address)
	var checksumLength int
	if len(addressBytes) == 32 {
		checksumLength = 2
	} else if utils.IntInSlice(len(addressBytes), []int{1, 2, 4, 8}) {
		checksumLength = 1
	} else {
		return ""
	}
	addressFormat := append(utils.HexToBytes("2a")[:], addressBytes[:]...)
	checksum, _ := blake2b.New(64, []byte{})
	w := append(checksumPrefix[:], addressFormat[:]...)
	checksum.Write(w)
	h := checksum.Sum(nil)
	b := append(addressFormat[:], h[:checksumLength][:]...)
	return base58.Encode(b)
}
