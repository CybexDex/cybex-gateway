package cybsrv

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	cybtypes "git.coding.net/yundkyy/cybexgolib/types"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func Test() {
	handleBlock(10311203)
}

// BlockRead ...
func BlockRead() {
	i := 0
	for {
		i = i + 1
		fmt.Println("read round:", i)
		readBlock()
		time.Sleep(time.Second * 3)
	}
}
func getlastBlock() (int, error) {
	// is there a last
	s, err := rep.Easy.FetchWith(&m.Easy{
		Key: "cybLastBlockNum",
	})
	if err != nil {
		return 0, err
	}
	if len(s) > 0 {
		c := s[0]
		i, err := strconv.Atoi(c.Value)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	// create it
	blockBegin := viper.GetInt("cybsrv.blockBegin")
	bstr := strconv.Itoa(blockBegin)
	newLast := &m.Easy{
		Key:   "cybLastBlockNum",
		Value: string(bstr),
	}
	err = newLast.Save()
	if err != nil {
		return 0, fmt.Errorf("cannot find")
	}

	return blockBegin, nil
}
func getHeadNum() (int, error) {
	s, err := api.GetDynamicGlobalProperties()
	if err != nil {
		return 0, err
	}
	return int(s.LastIrreversibleBlockNum), nil
}
func handleBlock(cnum int) {
	fmt.Println("block", cnum)
	// do read transfers
	blockInfo, err := api.GetBlock(uint64(cnum))
	b, err := json.Marshal(blockInfo)
	if err != nil {
		fmt.Println("error json block:", cnum)
		return
	}
	ts := gjson.GetBytes(b, "transactions")

	for index, t := range ts.Array() {
		ops := t.Get("operations")
		for _, op := range ops.Array() {
			if op.Get("0").Int() == 0 {
				rawop := op.Get("1")
				toid := rawop.Get("to").String()
				fromid := rawop.Get("from").String()
				asset := rawop.Get("amount.asset_id").String()
				amount := rawop.Get("amount.amount").String()
				if toid == gatewayAccount.ID.ID() {
					// is recharge
					if fromid == coldAccount.ID.ID() {
						fmt.Println("Recharge", rawop)
						continue
					}
					memostr := rawop.Get("memo").String()
					if memostr == "" {
						// UR
						fmt.Println("UR", rawop)
						continue
					} else {
						memo1 := &cybtypes.Memo{}
						json.Unmarshal([]byte(memostr), memo1)
						memoout := ""
						for _, prv := range gatewayMemoPri {
							s, err := memo1.Decrypt(&prv)
							if err == nil {
								memoout = s
							}
						}
						// is a withdraw
						if !strings.HasPrefix(memoout, gatewayPrefix) {
							fmt.Println("UR:1 ", cnum, index, fromid, asset, amount, "memo:", memoout)
							continue
						}
						s := strings.TrimPrefix(memoout, gatewayPrefix)
						s2 := strings.Split(s, ":")
						if len(s2) < 3 {
							fmt.Println("UR:2 ", cnum, index, fromid, asset, amount, "memo:", memoout, s)
							continue
						}
						addr := strings.Join(s2[2:], ":")
						fmt.Println("withdraw:", cnum, index, fromid, asset, amount, "memo:", memoout, "add:", addr)
						continue
					}
					// gen cyborder
				}
				if fromid == gatewayAccount.ID.ID() {
					// is from gateway => confirm
				}

			}
		}
	}
}
func readBlock() {
	// get lastBlock
	lastBlockNum, err := getlastBlock()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(lastBlockNum)
	// get blockhead
	blockheadNum, err := getHeadNum()
	fmt.Println(blockheadNum, err)
	if lastBlockNum >= blockheadNum {
		return
	}
	// for
	for cnum := lastBlockNum; cnum <= blockheadNum; cnum = cnum + 1 {
		handleBlock(cnum)
	}
	//
}
