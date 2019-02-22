package cybsrv

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	rep "git.coding.net/bobxuyang/cy-gateway-BN/help/singleton"
	m "git.coding.net/bobxuyang/cy-gateway-BN/models"
	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	cybtypes "coding.net/yundkyy/cybexgolib/types"
	"github.com/cockroachdb/apd"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func Test() {
	readBlock(10311203)
}

// BlockRead ...
func BlockRead() {
	i := 0
	for {
		i = i + 1
		fmt.Println("read round:", i)
		handleBlock()
		time.Sleep(time.Second * 3)
	}
}
func getlastBlock() (int, *m.Easy, error) {
	// is there a last
	s, err := rep.Easy.FetchWith(&m.Easy{
		Key: "cybLastBlockNum",
	})
	if err != nil {
		return 0, nil, err
	}
	if len(s) > 0 {
		c := s[0]
		i, err := strconv.Atoi(c.Value)
		if err != nil {
			return 0, nil, err
		}
		return i, c, nil
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
		return 0, nil, fmt.Errorf("cannot find")
	}

	return blockBegin, newLast, nil
}
func getHeadNum() (int, error) {
	s, err := api.GetDynamicGlobalProperties()
	if err != nil {
		return 0, err
	}
	return int(s.LastIrreversibleBlockNum), nil
}
func readBlock(cnum int) (orders []*m.CybOrder) {
	// fmt.Println("block", cnum)
	// do read transfers
	blockInfo, err := api.GetBlock(uint64(cnum))
	b, err := json.Marshal(blockInfo)
	if err != nil {
		fmt.Println("error json block:", cnum)
		return orders
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
				amount := rawop.Get("amount.amount").Int()
				if toid == gatewayAccount.ID.ID() {
					assetObj, err := rep.Asset.FetchWith(&m.Asset{
						CybID: asset,
					})
					if err != nil {
						fmt.Println("asset error", err)
						continue
					}
					if len(assetObj) < 1 {
						continue
					}
					assetNow := assetObj[0]

					UserID := cybtypes.NewGrapheneID(fromid)
					fromusers, err := api.GetAccounts(UserID)
					if err != nil {
						fmt.Println("fromusers error", err)
						continue
					}
					userFrom := fromusers[0]
					cybasset, err := api.GetAsset(assetNow.CybName)
					app, err := rep.App.FindAppOrCreate(userFrom.Name)
					realAmount := rep.Asset.RawAmountToReal(amount, cybasset.Precision)
					f1, _ := realAmount.Float64()
					f2, _ := assetNow.WithdrawFee.Float64()
					amountNow := f1 - f2
					amountStr := fmt.Sprintf("%f", amountNow)
					amountA, _, _ := apd.NewFromString(amountStr)
					order := &m.CybOrder{
						AppID:       app.ID,
						AssetID:     assetNow.ID,
						From:        userFrom.Name,
						To:          gatewayAccount.Name,
						TotalAmount: realAmount,
						Amount:      amountA,
						Fee:         assetNow.WithdrawFee,
					}
					// is recharge
					if fromid == coldAccount.ID.ID() {
						fmt.Println("Recharge", rawop)
						order.Type = m.CybOrderTypeRecharge
						orders = append(orders, order)
						continue
					}
					memostr := rawop.Get("memo").String()
					if memostr == "" {
						// UR
						fmt.Println("UR", rawop)
						order.Type = m.CybOrderTypeUR
						orders = append(orders, order)
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
						if amountNow < 0 {
							fmt.Println("UR", rawop)
							order.Type = m.CybOrderTypeUR
							orders = append(orders, order)
							continue
						}
						// is a withdraw
						if !strings.HasPrefix(memoout, gatewayPrefix) {
							fmt.Println("UR:1 ", cnum, index, fromid, asset, amount, "memo:", memoout)
							order.Type = m.CybOrderTypeUR
							orders = append(orders, order)
							continue
						}
						s := strings.TrimPrefix(memoout, gatewayPrefix)
						s2 := strings.Split(s, ":")
						if len(s2) < 3 {
							fmt.Println("UR:2 ", cnum, index, fromid, asset, amount, "memo:", memoout, s)
							order.Type = m.CybOrderTypeUR
							orders = append(orders, order)
							continue
						}
						addr := strings.Join(s2[2:], ":")
						fmt.Println("withdraw:", cnum, index, fromid, asset, amount, "memo:", memoout, "add:", addr)
						order.Type = m.CybOrderTypeWithdraw
						order.Memo = memoout
						order.WithdrawAddr = addr
						orders = append(orders, order)
						continue
					}
				}
				if fromid == gatewayAccount.ID.ID() {
					// is from gateway => confirm
				}

			}
		}
	}
	return orders
}
func updateLastBlock(cnum int, easy *m.Easy) error {
	bstr := strconv.Itoa(cnum + 1)
	easy.Value = bstr
	err := easy.Save()
	return err
}
func handleBlock() {
	// get lastBlock
	lastBlockNum, easy, err := getlastBlock()
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
		cyborders := readBlock(cnum)
		// fmt.Println(cyborders)
		// save cyborders
		for _, order := range cyborders {
			saveCYBOrder(order)
		}
		// updateLastBlock
		updateLastBlock(cnum, easy)
	}
	//
}
func saveCYBOrder(order *m.CybOrder) error {
	tx := m.GetDB().Begin()
	defer func() {
		tx.Commit()
		if r := recover(); r != nil {
			utils.Errorf("%v, stack: %s", r, debug.Stack())
			tx.Rollback()
		}
	}()
	tx.Save(order)
	err := order.CreateOrder(tx)
	return err
}
