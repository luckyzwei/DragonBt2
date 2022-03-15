package game

import (
	"encoding/json"
	"fmt"
)

type ModGuide struct {
	player *Player
	Data   San_GUide //! 引导Id
}

// 限时神将模块
type San_GUide struct {
	Uid  int64
	Info string

	info []int //! 引导Id
	DataUpdate
}

func (self *ModGuide) OnGetData(player *Player) {
	self.player = player
}

func (self *ModGuide) initInfo() {
	self.Data.info = make([]int, 0)
}

func (self *ModGuide) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_guides` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Data, "san_guides", self.player.ID)
	if self.Data.Uid <= 0 {
		self.Data.Uid = self.player.ID
		self.initInfo()
		self.Encode()
		InsertTable("san_guides", &self.Data, 0, true)
	} else {

		self.Decode()
	}
	self.Data.Init("san_guides", &self.Data, true)
}

func (self *ModGuide) Decode() {
	json.Unmarshal([]byte(self.Data.Info), &self.Data.info)
}

func (self *ModGuide) Encode() {
	self.Data.Info = HF_JtoA(self.Data.info)
}

func (self *ModGuide) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "add_guide":
		var c2s_msg C2S_AddGuide
		json.Unmarshal(body, &c2s_msg)
		self.AddGuide(c2s_msg.GuideId)
		return true
	case "get_guide":
		var msg S2C_Guides
		msg.Cid = "guides"
		msg.Guides = self.Data.info
		self.player.Send(msg.Cid, msg)
		return true
	case "add_story":
		var c2s_msg C2S_AddStory
		json.Unmarshal(body, &c2s_msg)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GUIDE_STORY, c2s_msg.StoryID, c2s_msg.StoryType, 0, "剧情", 0, 0, self.player)
		return true
	}
	return false
}

func (self *ModGuide) OnSave(sql bool) {
	self.Encode()
	self.Data.Update(sql)
}

func (self *ModGuide) SendInfo() {
	var msg S2C_Guides
	msg.Cid = "guides"
	msg.Guides = self.Data.info
	self.player.Send(msg.Cid, msg)
}

func (self *ModGuide) AddGuide(guidID int) {
	if guidID == 0 {
		return
	}

	found := false
	for _, v := range self.Data.info {
		if v == guidID {
			found = true
			break
		}
	}

	if !found {
		self.Data.info = append(self.Data.info, guidID)

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GUIDE_START, guidID, 0, 0, "引导", 0, 0, self.player)
	}
}
