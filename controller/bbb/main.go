package bbb

import (
	"bytes"
	"cybex-gateway/types"
	"github.com/spf13/viper"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"cybex-gateway/utils/log"
)
// InfoResponse ...
type InfoResponse struct {
	BlockID    string               `json:"block_id"`
	BlockNum uint32                 `json:"block_num"`
}
// broad
type broadResponse struct {
	Code    int               `json:"code"`
}
// Info ...
func Info() (out *InfoResponse,err error){
	url  := viper.GetString("bbb.info_url")
	resp, err := http.Get(url)
	if err != nil {
		return nil,fmt.Errorf("post error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil,fmt.Errorf("ReadAll error: %v", err)
	}
	respData := InfoResponse{}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return nil,fmt.Errorf("Unmarshal error: %v, body: %s", err, string(bodyBytes))
	}
	return &respData,nil
}
// BroadcastTransaction ...
func BroadcastTransaction(tx string)(err error) {
	outbody := fmt.Sprintf(`{"transaction":%s}`,tx)
	bs := []byte(outbody)
	url  := viper.GetString("bbb.broad_url")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return fmt.Errorf("post error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("ReadAll error: %v", err)
	}
	respData := broadResponse{}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return fmt.Errorf("Unmarshal error: %v, body: %s", err, string(bodyBytes))
	}
	if respData.Code != 0 {
		return fmt.Errorf("code error: %d",respData.Code)
	}
	return nil
}
func bbbResult(urlPath string, sendData *types.JPSendData, v interface{}) (err error) {
	bs, _ := json.Marshal(sendData)
	strsend := string(bs)
	log.Debugln("send jp json", strsend)
	jadepoolAddr := viper.GetString("jpserver.bnhost")
	url := jadepoolAddr + urlPath
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return fmt.Errorf("post error: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("ReadAll error: %v", err)
	}
	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return fmt.Errorf("Unmarshal error: %v, body: %s", err, string(bodyBytes))
	}
	return nil
}

