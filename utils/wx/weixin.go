package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var users, corpid, corpsecret, tokenNow string
var agentid int

type tokenMsg struct {
	Code  int    `json:"errcode"`
	Token string `json:"access_token"`
	Msg   string `json:"errmsg"`
}

// InitWeixin ...
func InitWeixin(corpid1 string, corpsecret1 string, agentid1 int, users1 string) {
	users = users1
	corpid = corpid1
	corpsecret = corpsecret1
	agentid = agentid1
}

// GetToken ...
func GetToken() (token string, err error) {
	if tokenNow != "" {
		return tokenNow, nil
	}
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpid, corpsecret)
	resp, _ := http.Get(url)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	data := &tokenMsg{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return "", err
	}
	return data.Token, nil
}

type sendMsg struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	Text    text   `json:"text"`
	Safe    int    `json:"safe"`
}
type text struct {
	Content string `json:"content"`
}

func sendOneTime(token string, msgTo string) (*tokenMsg, error) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + token
	sendData := &sendMsg{
		Touser:  users,
		Msgtype: "text",
		Agentid: agentid,
		Text: text{
			Content: msgTo,
		},
		Safe: 0,
	}
	bs, _ := json.Marshal(sendData)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	data := &tokenMsg{}
	err = json.Unmarshal(bodyBytes, &data)
	return data, err
}

// WeixinSend ...
func WeixinSend(msgTo string) error {
	token, err := GetToken()
	if err != nil {

	}
	one, err := sendOneTime(token, msgTo)
	if err != nil {
		return err
	}
	if one.Code == 40014 {
		fmt.Println("reget token")
		tokenNow = ""
		token, err = GetToken()
		one, err = sendOneTime(token, msgTo)
	}
	if one.Msg == "ok" {
		return nil
	}
	return fmt.Errorf("sendweixin error: %v", *one)
}
