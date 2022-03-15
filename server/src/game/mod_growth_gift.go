package game

import (
	"encoding/json"
	"fmt"
)

///////////////////////////////
/*
成长礼包
*/
///////////////////////////////
const (
	MSG_GROWTH_GIFT_SEND_INFO = "growth_gift_send_info"
	MSG_GROWTH_GIFT_GET_AWARD = "growth_gift_get_award"
	MSG_GROWTH_GIFT_UPDATE    = "growth_gift_update"
)

const GROWTH_GIFT_BASE = 100101

//! 活动结构
type GrowthGiftItem struct {
	Id        int `json:"id"`        //! 任务id
	Tasktypes int `json:"tasktypes"` // 任务类型
	Plan      int `json:"plan"`      //! 进度
	Done      int `json:"done"`      //! 完成状态, 当所有的完成了,才算完成
}

//! 限时礼包 数据库
type San_GrowthGift struct {
	Uid  int64  //! UID
	Info string //! 保存数据

	info []*GrowthGiftItem
	DataUpdate
}

//! 限时礼包
type ModGrowthGift struct {
	player         *Player
	San_GrowthGift San_GrowthGift

	chg []*GrowthGiftItem
}

func (self *ModGrowthGift) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_growthgift` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.San_GrowthGift, "san_growthgift", self.player.ID)

	if self.San_GrowthGift.Uid <= 0 {
		self.San_GrowthGift.Uid = self.player.ID
		self.Encode()
		InsertTable("san_growthgift", &self.San_GrowthGift, 0, true)
		self.San_GrowthGift.Init("san_growthgift", &self.San_GrowthGift, true)
	} else {
		self.Decode()
		self.San_GrowthGift.Init("san_growthgift", &self.San_GrowthGift, true)
	}
}

func (self *ModGrowthGift) OnGetOtherData() {
}

func (self *ModGrowthGift) OnSave(sql bool) {
	self.Encode()
	self.San_GrowthGift.Update(sql)
}

// 老的消息处理
func (m *ModGrowthGift) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModGrowthGift) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.San_GrowthGift.Info), &self.San_GrowthGift.info)
}

func (self *ModGrowthGift) Encode() { //! 将data数据写入数据库
	self.San_GrowthGift.Info = HF_JtoA(&self.San_GrowthGift.info)
}

//! 充值回调
func (self *ModGrowthGift) HandleTask(tasktype, n2, n3, n4 int) {
	configs := GetCsvMgr().GrowthGiftConfig
	for i := 0; i < len(configs); i++ {
		if configs[i].TaskTypes != tasktype {
			continue
		}

		var data *GrowthGiftItem = nil
		for _, v := range self.San_GrowthGift.info {
			if v.Id == configs[i].Id {
				data = v
				break
			}
		}

		if data != nil {
			if data.Done == 1 || data.Done == 2 {
				continue
			}
		}

		var tasknode TaskNode
		tasknode.Tasktypes = configs[i].TaskTypes
		tasknode.N1 = configs[i].N[0]
		tasknode.N2 = configs[i].N[1]
		tasknode.N3 = configs[i].N[2]
		tasknode.N4 = configs[i].N[3]
		plan, add := DoTask(&tasknode, self.player, n2, n3, n4)
		if plan == 0 {
			continue
		}

		if data == nil {
			data = &GrowthGiftItem{configs[i].Id, configs[i].TaskTypes, 0, 0}
			self.San_GrowthGift.info = append(self.San_GrowthGift.info, data)
		}

		chg := false
		if add {
			data.Plan += plan
			chg = true
		} else {
			if plan > data.Plan {
				data.Plan = plan
				chg = true
			}
		}

		if data.Plan >= configs[i].N[0] {
			data.Done = 1
		}

		if chg {
			self.chg = append(self.chg, data)
		}
	}
}

// 发送礼包更新信息
func (self *ModGrowthGift) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_GrowthGiftUpdate
	msg.Cid = MSG_GROWTH_GIFT_UPDATE
	msg.Info = self.chg
	self.chg = make([]*GrowthGiftItem, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(MSG_GROWTH_GIFT_UPDATE, smsg)
}

// 注册消息
func (self *ModGrowthGift) onReg(handlers map[string]func(body []byte)) {
	handlers[MSG_GROWTH_GIFT_SEND_INFO] = self.GrowthGiftSendInfo
	handlers[MSG_GROWTH_GIFT_GET_AWARD] = self.ActivityGiftGetAward
}
func (self *ModGrowthGift) GrowthGiftSendInfo(body []byte) {
	isRecharge := false
	actitem, ok := self.player.GetModule("activity").(*ModActivity).Sql_Activity.info[GROWTH_GIFT_BASE]
	if ok {
		if actitem.Done == 1 {
			isRecharge = true
		}
	}

	var backmsg S2C_GrowthGiftSendInfo
	backmsg.Cid = MSG_GROWTH_GIFT_SEND_INFO
	backmsg.IsRecharge = isRecharge
	backmsg.Info = self.San_GrowthGift.info
	backmsg.Config = GetCsvMgr().GrowthGiftConfig
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}
func (self *ModGrowthGift) ActivityGiftGetAward(body []byte) {
	var msg C2S_GrowthGiftGetAward
	json.Unmarshal(body, &msg)

	if !self.player.GetModule("activity").(*ModActivity).isActivityGrowthGiftOpen() {
		self.player.SendErrInfo("err", "活动未开启")
		return
	}

	actitem, ok := self.player.GetModule("activity").(*ModActivity).Sql_Activity.info[GROWTH_GIFT_BASE]
	if !ok {
		self.player.SendErrInfo("err", "活动未开启")
		return
	}

	if actitem.Done != 1 {
		self.player.SendErrInfo("err", "活动未开启")
		return
	}

	config := GetCsvMgr().GetGrowthGiftConfig(msg.ID)
	if nil == config {
		self.player.SendErrInfo("err", "配置不存在")
		return
	}

	var data *GrowthGiftItem = nil
	find := false
	for _, v := range self.San_GrowthGift.info {
		if v.Id == msg.ID {
			data = v
			find = true
			break
		}
	}
	if !find || nil == data {
		self.player.SendErrInfo("err", "任务未找到")
		return
	}

	if data.Done != 1 {
		self.player.SendErrInfo("err", "任务未完成")
		return
	}

	data.Done = 2

	item := self.player.AddObjectLst(config.Items, config.Nums, "活动礼包", msg.ID, 0, 0)

	var backmsg S2C_GrowthGiftGetAward
	backmsg.Cid = MSG_GROWTH_GIFT_GET_AWARD
	backmsg.ID = msg.ID
	backmsg.Items = item
	self.player.Send(backmsg.Cid, backmsg)
}

func (self *ModGrowthGift) GMFinishGrowthGift(id int) {
	configs := GetCsvMgr().GrowthGiftConfig
	for i := 0; i < len(configs); i++ {
		if configs[i].Id != id {
			continue
		}

		var data *GrowthGiftItem = nil
		for _, v := range self.San_GrowthGift.info {
			if v.Id == configs[i].Id {
				data = v
				break
			}
		}

		if data != nil {
			if data.Done == 1 || data.Done == 2 {
				continue
			}
		}

		if data == nil {
			data = &GrowthGiftItem{configs[i].Id, configs[i].TaskTypes, 0, 1}
			self.San_GrowthGift.info = append(self.San_GrowthGift.info, data)
		} else {
			data.Done = 1
		}

		break
	}
}
