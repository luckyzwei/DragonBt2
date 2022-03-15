package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	//发送消息使用导的url
	sendurl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`
	//获取token使用导的url
	get_token = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`
)

var requestError = errors.New("request error,check url or network")

type access_token struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

//定义一个简单的文本消息格式
type send_msg struct {
	Touser  string            `json:"touser"`
	Toparty string            `json:"toparty"`
	Totag   string            `json:"totag"`
	Msgtype string            `json:"msgtype"`
	Agentid int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

type send_msg_error struct {
	Errcode int    `json:"errcode`
	Errmsg  string `json:"errmsg"`
}

type WechatWarningMgr struct {
	ToUser     string
	AgentID    int
	CorpID     string
	CorpSecret string
}

var wechatwarningsingleton *WechatWarningMgr = nil

//! public
func GetWechatWarningMgr() *WechatWarningMgr {
	if wechatwarningsingleton == nil {
		wechatwarningsingleton = new(WechatWarningMgr)
		wechatwarningsingleton.ToUser = "@all"
		wechatwarningsingleton.AgentID = 1000003
		wechatwarningsingleton.CorpID = "wwebdcba2b00260eaa"
		wechatwarningsingleton.CorpSecret = "ay2bwZVegzKqOest46OhJoh9p7GsHTuPfBdGRYGFq8o"
	}

	return wechatwarningsingleton
}

func (self *WechatWarningMgr) GetData() {
}

func (self *WechatWarningMgr) SendWarning(warning string) {
	if self.CorpID == "" || self.CorpSecret == "" {
		//flag.Usage()
		return
	}

	var m send_msg = send_msg{Touser: self.ToUser, Toparty: "@all", Msgtype: "text", Agentid: self.AgentID, Text: map[string]string{"content": warning}}

	token, err := self.Get_token(self.CorpID, self.CorpSecret)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println("获取token成功：", token)
	buf, err := json.Marshal(m)
	if err != nil {
		return
	}
	err = self.Send_msg(token.Access_token, buf)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Println("发送消息成功", string(buf))
	}
}

//发送消息.msgbody 必须是 API支持的类型
func (self *WechatWarningMgr) Send_msg(Access_token string, msgbody []byte) error {
	body := bytes.NewBuffer(msgbody)
	resp, err := http.Post(sendurl+Access_token, "application/json", body)
	if resp.StatusCode != 200 {
		return requestError
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var e send_msg_error
	err = json.Unmarshal(buf, &e)
	if err != nil {
		return err
	}
	if e.Errcode != 0 && e.Errmsg != "ok" {
		return errors.New(string(buf))
	}
	return nil
}

//通过corpid 和 corpsecret 获取token
func (self *WechatWarningMgr) Get_token(corpid, corpsecret string) (at access_token, err error) {
	resp, err := http.Get(get_token + corpid + "&corpsecret=" + corpsecret)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = requestError
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &at)
	if at.Access_token == "" {
		err = errors.New("corpid or corpsecret error.")
	}
	return
}
