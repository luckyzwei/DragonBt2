package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

const (
	ROBOR_ACTION_DRAW_GEM = iota
	ROBOR_ACTION_DRAW_GEM_TEN
	ROBOR_ACTION_NOBILITYTASKINFO
	ROBOR_ACTION_ARMSARENAINFO
	ROBOR_ACTION_GETARMSARENAINFO
	ROBOR_ACTION_CAMPFIGHTINFOLIST
	ROBOR_ACTION_CAMPFIGHTINFO
	ROBOR_ACTION_DREAMLANDINFO
	ROBOR_ACTION_ACTIVITYMASK
	ROBOR_ACTION_CHAT
	ROBOR_ACTION_PASS

	ROBOR_ACTION_DRAW_END
)

// 新版机器人 在线玩游戏
func (self *Robot) BrainRun() {
	ticker := time.NewTicker(time.Second)
	self.csvId = 0
	self.guildId = 0
	self.Logic() //首先登录
	for {
		<-ticker.C
		if self.Dead {
			break
		}
		if self.OperTime == 0 {
			continue
		}
		//fmt.Println("csvId =",self.csvId)
		if time.Now().Unix() >= self.OperTime {
			self.CheckResourse()
			self.Action()
			timedis:= HF_GetRandom(GetRobotMgr().RobotCon.ActionDis)
			self.OperTime = time.Now().Unix() + int64(timedis+30)
		}
	}
}

//检查基础资源确保有足够资源
func (self *Robot) CheckResourse() {
	if self.userbase.Gold <= 10000000 {
		self.CreateItem(91000001, 100000000)
	}

	if self.userbase.Gem <= 10000 {
		self.CreateItem(91000002, 100000)
	}

	if self.userbase.Power <= 100 {
		self.CreateItem(91000003, 2000)
	}

	if self.userbase.Level < 90 {
		self.CreateItem(91000005, 100000)
	}
	fmt.Println("金币=", self.userbase.Gold, "钻石=", self.userbase.Gem, "体力=", self.userbase.Power, "等级=", self.userbase.Level)
}

func (self *Robot) Action() {

	action := HF_GetRandom(ROBOR_ACTION_DRAW_END)

	switch action {
	case ROBOR_ACTION_DRAW_GEM:
		self.Find(3, 0)
		return
	case ROBOR_ACTION_DRAW_GEM_TEN:
		self.Find(5, 0)
		return
	case ROBOR_ACTION_NOBILITYTASKINFO:
		self.Send("", C2S_Uid{Ctrl: "nobilitytaskinfo", Uid: self.Uid})
		return
	case ROBOR_ACTION_ARMSARENAINFO:
		self.Send("", C2S_Uid{Ctrl: "armsarenainfo", Uid: self.Uid})
		return
	case ROBOR_ACTION_GETARMSARENAINFO:
		self.Send("", C2S_Uid{Ctrl: "getarmsarenafight", Uid: self.Uid})
		return
	case ROBOR_ACTION_CAMPFIGHTINFOLIST:
		self.Send("", C2S_Uid{Ctrl: "campfightinfolist", Uid: self.Uid})
		return
	case ROBOR_ACTION_CAMPFIGHTINFO:
		self.Send("", C2S_Uid{Ctrl: "campunioninfo", Uid: self.Uid})
		return
	case ROBOR_ACTION_DREAMLANDINFO:
		self.Send("", C2S_Uid{Ctrl: "get_dreamland_base_info", Uid: self.Uid})
		return
	case ROBOR_ACTION_ACTIVITYMASK:
		self.Send("", C2S_Uid{Ctrl: "getactivitymask", Uid: self.Uid})
		return
	case ROBOR_ACTION_CHAT:
		self.RobotChat()
		return
	case ROBOR_ACTION_PASS:
		self.SendPassId(110101)
		return
	}
}

//检查基础资源确保有足够资源
func (self *Robot) RobotChat() {
	channel := 1
	var chats = []string{
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

	websocket.Message.Send(self.Ws, HF_DecodeMsg("chat", msg))
}
