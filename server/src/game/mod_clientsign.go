package game

import (
	"encoding/json"
	"fmt"
)

//
type San_ClientSign struct {
	Uid  int64  `json:"uid"` // 玩家Id
	Sign string `json:"sign"`

	sign map[int]int // 标记
	DataUpdate
}

type ModClientSign struct {
	player         *Player
	San_ClientSign San_ClientSign
}

func (self *ModClientSign) OnGetData(player *Player) {
	self.player = player
	sql := fmt.Sprintf("select * from `san_clientsign` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.San_ClientSign, "san_clientsign", self.player.ID)

	if self.San_ClientSign.Uid <= 0 {
		self.San_ClientSign.Uid = self.player.ID
		self.San_ClientSign.sign = make(map[int]int, 0)
		self.Encode()
		InsertTable("san_clientsign", &self.San_ClientSign, 0, true)
		self.San_ClientSign.Init("san_clientsign", &self.San_ClientSign, true)
	} else {
		self.Decode()
		self.San_ClientSign.Init("san_clientsign", &self.San_ClientSign, true)
	}

	if self.San_ClientSign.sign == nil {
		self.San_ClientSign.sign = make(map[int]int)
	}
}

func (self *ModClientSign) OnSave(sql bool) {
	self.Encode()
	self.San_ClientSign.Update(sql)
}

func (self *ModClientSign) OnGetOtherData() {

}

func (self *ModClientSign) Decode() { // 将数据库数据写入data
	json.Unmarshal([]byte(self.San_ClientSign.Sign), &self.San_ClientSign.sign)
}

func (self *ModClientSign) Encode() { // 将data数据写入数据库
	self.San_ClientSign.Sign = HF_JtoA(&self.San_ClientSign.sign)
}

func (self *ModClientSign) SendInfo() {

	var msgRel S2C_ClientSign
	msgRel.Cid = "clientsigninfo"
	msgRel.Sign = self.San_ClientSign.sign
	self.player.SendMsg("clientsigninfo", HF_JtoB(&msgRel))
}

func (self *ModClientSign) OnMsg(ctrl string, body []byte) bool {
	return false
}
func (self *ModClientSign) onReg(handlers map[string]func(body []byte)) {
	handlers["setclientsign"] = self.SetClientSign
}

func (self *ModClientSign) SetClientSign(body []byte) {
	var msg C2S_SetClientSign
	json.Unmarshal(body, &msg)

	if self.San_ClientSign.sign == nil {
		self.San_ClientSign.sign = make(map[int]int)
	}

	self.San_ClientSign.sign[msg.Key] = msg.Value

	var msgRel S2C_SetClientSign
	msgRel.Cid = "setclientsign"
	msgRel.Key = msg.Key
	msgRel.Value = self.San_ClientSign.sign[msg.Key]
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}
