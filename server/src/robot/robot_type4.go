package main

import (
	"fmt"
	"time"
)

const TIMEFORMAT = "2006-01-02 15:04:05" //! 时间格式化


func (self *Robot) Logic() {
	fmt.Println("robot:", self.Account, " running, type = 4")
	fmt.Println("login start:", time.Now().Format(TIMEFORMAT))
	self.LoginGuest()

}

// pc登录
func (self *Robot) LoginGuest()  {
	fmt.Println("机器人发送登录消息:login_guest")
	var msg C2S_Reg
	msg.Ctrl = "login_guest"
	msg.Uid = self.Uid
	msg.Account = self.Account
	msg.Password = self.Password
	msg.ServerId = self.ServerId
	self.Send("passport.php", &msg)
}


func (self *Robot) SetCamp()  {
	fmt.Println("机器人发送登录消息:setcamp")
	var msg C2S_SetCamp
	msg.Uid = self.Uid
	msg.Camp = HF_GetRandom(3) + 1
	msg.Ctrl = "setcamp"
	self.Send("", &msg)
}

// 进行聊天
func (self *Robot) SendMsg()  {
	/*
	channel := 1
	var chats  =[]string {
		" 看少侠骨骼奇特，悟性非凡，相信一定可以做出一番惊天动地的大事。",
		"这个世界上并没有成功和失败，有的只是经验。正如聪明和愚昧，生和死，都是不同的经验。",
		"魔说虽然今天我死了，但是在3000年后，我的子子孙孙会穿上你们佛的衣服。",
	}
	content := fmt.Sprintf(chats[HF_GetRandom(len(chats))])
	name := self.Uname
	data := &Chat{
		Channel:   proto.Int32(int32(channel)),
		Content:   []byte(content),
		Name:      proto.String(name),
		Medianame: proto.String(""),
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	for {
		err := websocket.Message.Send(self.Ws, HF_DecodeMsg("chat", msg))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		time.Sleep(time.Second * 5)
	}
		*/
}