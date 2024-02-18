package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type WxMsg struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func SendOrderMsg(content string) error {

	v := url.Values{}
	v.Add("title", "重要通知")
	v.Add("content", content)
	v.Add("user_token", "f7d51095aeb1d0acacd6c28937ddc6a21712679569")

	resp, err := http.Post("https://push.showdoc.com.cn/server/api/push/f1f44627566f111984732a18c4a2ae90944797916", "application/x-www-form-urlencoded", strings.NewReader(v.Encode()))

	//client := &http.Client{}
	//
	//data := &WxMsg{
	//	MsgType: "text",
	//	Text: struct {
	//		Content string "json:\"content\""
	//	}{
	//		content,
	//	},
	//}
	//
	//d, err := json.Marshal(data)
	//if err != nil {
	//	return err
	//}
	//req, _ := http.NewRequest("POST", WXROBOTURL, bytes.NewBuffer(d))
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("charset", "UTF-8")
	//resp, err := client.Do(req)
	if err != nil {
		fmt.Println("send msg err", err)
		return err
	}
	defer resp.Body.Close()
	d, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("resp", string(d))
	return nil

}
