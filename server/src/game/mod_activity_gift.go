package game

import (
	"encoding/json"
	"fmt"
	"time"
)

///////////////////////////////
/*
活动礼包模块 控制日礼包 周礼包和月礼包 新手礼包
活动始终开放，根据配置从玩家创建账号后的第X天开始启动每日/周/月刷新次数
*/
///////////////////////////////
const (
	MSG_ACTIVITY_GIFT_SEND_INFO = "activity_gift_send_info"
	MSG_ACTIVITY_GIFT_GET_AWARD = "activity_gift_get_award"
)
const (
	ACTIVITY_GIFT_TYPE_DAY       = 1 // 日
	ACTIVITY_GIFT_TYPE_WEEK      = 2 // 周
	ACTIVITY_GIFT_TYPE_MONTH     = 3 // 月
	ACTIVITY_GIFT_TYPE_NOVICE    = 4 // 新手
	ACTIVITY_GIFT_TYPE_BOX       = 5 // 老礼包类型 activitybox
	ACTIVITY_GIFT_TYPE_DISCOUNT  = 6 // 特惠礼包
	ACTIVITY_GIFT_TYPE_STAR      = 7 // 星辰限时礼包
	ACTIVITY_GIFT_TYPE_STAR_HERO = 8 // 星辰礼包
	ACTIVITY_GIFT_TYPE_GIFT_EX   = 9 // 活动礼包
)

//! 活动结构
type ActivityGiftItem struct {
	BoxId     int   `json:"boxid"`     //! 礼包Id
	Type      int   `json:"type"`      //! 礼包类型
	Pickup    int   `json:"pickup"`    //! 暂时不用
	Plan      int   `json:"plan"`      //! 充值金额
	Done      int   `json:"done"`      //! 完成状态, 当所有的完成了,才算完成
	Times     int   `json:"times"`     //! 已完成次数
	TimeStamp int64 `json:"timestamp"` //! 时间戳
}

//! 限时礼包 数据库
type San_ActivityGift struct {
	Uid  int64  //! UID
	Info string //! 保存数据

	info []*ActivityGiftItem
	DataUpdate
}

//! 限时礼包
type ModActivityGift struct {
	player           *Player
	San_ActivityGift San_ActivityGift

	chg []*ActivityGiftItem
}

func (self *ModActivityGift) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_activitygift` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.San_ActivityGift, "san_activitygift", self.player.ID)

	if self.San_ActivityGift.Uid <= 0 {
		self.San_ActivityGift.Uid = self.player.ID
		self.Encode()
		InsertTable("san_activitygift", &self.San_ActivityGift, 0, true)
		self.San_ActivityGift.Init("san_activitygift", &self.San_ActivityGift, true)
	} else {
		self.Decode()
		self.San_ActivityGift.Init("san_activitygift", &self.San_ActivityGift, true)
	}
}

func (self *ModActivityGift) OnGetOtherData() {
	self.CheckOldTaskID()
}

func (self *ModActivityGift) OnSave(sql bool) {
	self.Encode()
	self.San_ActivityGift.Update(sql)
}

// 老的消息处理
func (m *ModActivityGift) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModActivityGift) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.San_ActivityGift.Info), &self.San_ActivityGift.info)
}

func (self *ModActivityGift) Encode() { //! 将data数据写入数据库
	self.San_ActivityGift.Info = HF_JtoA(&self.San_ActivityGift.info)
}

//! 充值回调
func (self *ModActivityGift) HandleTask(tasktype, n2, n3, n4 int) {
	configs := GetCsvMgr().ActivityGiftConfig
	for i := 0; i < len(configs); i++ {
		if configs[i].Group != 0 && configs[i].Type != ACTIVITY_GIFT_TYPE_GIFT_EX {
			continue
		}
		if configs[i].TaskTypes != tasktype {
			continue
		}

		if configs[i].Sale == 0 {
			continue
		}

		if configs[i].Pic2Type == 1 && configs[i].StarHero != 0 {
			have := false
			hero := self.player.GetModule("hero").(*ModHero).GetVoidHero(configs[i].StarHero)
			if hero != nil {
				have = true
			}

			// 特殊处理 12 + 英雄id + 01 是英雄的魂石id
			cardid := configs[i].StarHero*100 + 12000000 + 1
			carditem := GetCsvMgr().GetItemConfig(cardid)
			if carditem != nil {
				if err := self.player.HasObjectOkEasy(cardid, carditem.CompoundNum); err == nil {
					have = true
				}
			}

			if have {
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

		id := configs[i].ActivityType*100000 + configs[i].Group*100 + configs[i].Index
		var data *ActivityGiftItem = nil
		if plan > 0 {
			find := false
			for _, v := range self.San_ActivityGift.info {
				if v.BoxId == id {
					data = v
					find = true
					break
				}
			}
			if !find {
				data = &ActivityGiftItem{id, configs[i].Type, 0, 0, 0, 0, self.GetTimeStamp(configs[i].Id)}
				self.San_ActivityGift.info = append(self.San_ActivityGift.info, data)
			}
		}

		if data == nil {
			continue
		}

		if data.Done == 1 {
			continue
		}

		if data.Times >= configs[i].Times {
			continue
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
			if data.Done == 1 { // 完成了没有进行领取
				continue
			}
			data.Done = 1
		}

		if chg {
			self.chg = append(self.chg, data)
		}
	}
}

// 发送礼包更新信息
func (self *ModActivityGift) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_ActivityGiftUpdate
	msg.Cid = "activity_gift_update"
	msg.Info = self.chg
	self.chg = make([]*ActivityGiftItem, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("activity_gift_update", smsg)
}

// 注册消息
func (self *ModActivityGift) onReg(handlers map[string]func(body []byte)) {
	handlers[MSG_ACTIVITY_GIFT_SEND_INFO] = self.ActivityGiftSendInfo
	handlers[MSG_ACTIVITY_GIFT_GET_AWARD] = self.ActivityGiftGetAward
}
func (self *ModActivityGift) ActivityGiftSendInfo(body []byte) {
	var backmsg S2C_ActivityGiftSendInfo
	backmsg.Cid = MSG_ACTIVITY_GIFT_SEND_INFO
	backmsg.Info = self.San_ActivityGift.info
	backmsg.Config = GetCsvMgr().ActivityGiftConfig
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

func (self *ModActivityGift) GetAllAward(grade int) {
	for _, v := range self.San_ActivityGift.info {
		config := GetCsvMgr().GetActivityGiftConfig(v.BoxId)
		if nil == config {
			continue
		}

		//日周月礼包后台充值直接发
		if config.ActivityType == ACT_DAY_GIFT || config.ActivityType == ACT_WEEK_GIFT || config.ActivityType == ACT_MONTH_GIFT || config.ActivityType == ACT_NOVICE_GIFT {
			v.Done = 1
		}

		if v.Done == 1 && config.Sale != 0 && config.TaskTypes == TASK_TYPE_RECHARGE_ONCE && grade == config.N[1] {
			var msg C2S_ActivityGiftGetAward
			msg.ID = v.BoxId
			smsg, _ := json.Marshal(&msg)
			self.ActivityGiftGetAward(smsg)
		}
	}
}

func (self *ModActivityGift) ActivityGiftGetAward(body []byte) {
	var msg C2S_ActivityGiftGetAward
	json.Unmarshal(body, &msg)

	config := GetCsvMgr().GetActivityGiftConfig(msg.ID)
	if nil == config {
		self.player.SendErrInfo("err", "配置不存在")
		return
	}

	activity := GetActivityMgr().GetActivity(config.ActivityType)
	if activity == nil {
		return
	}

	now := TimeServer().Unix()
	startTime := HF_CalTimeForConfig(activity.info.Start, self.player.Sql_UserBase.Regtime)
	endTime := startTime + int64(activity.info.Continued) + int64(activity.info.Show)

	if startTime != endTime && (now < startTime || now > endTime) {
		return
	}

	var data *ActivityGiftItem = nil
	find := false
	for _, v := range self.San_ActivityGift.info {
		if v.BoxId == msg.ID {
			data = v
			find = true
			break
		}
	}
	if !find || nil == data {
		if config.Sale == 0 {
			id := config.ActivityType*100000 + config.Group*100 + config.Index
			data = &ActivityGiftItem{id, config.Type, 0, 0, 0, 0, self.GetTimeStamp(config.Id)}
			self.San_ActivityGift.info = append(self.San_ActivityGift.info, data)
		} else {
			self.player.SendErrInfo("err", "任务未找到")
			return
		}
	}

	if config.Pic2Type == 1 && config.StarHero != 0 {
		have := false
		hero := self.player.GetModule("hero").(*ModHero).GetVoidHero(config.StarHero)
		if hero != nil {
			have = true
		}

		// 特殊处理 12 + 英雄id + 01 是英雄的魂石id
		cardid := config.StarHero*100 + 12000000 + 1
		carditem := GetCsvMgr().GetItemConfig(cardid)
		if carditem != nil {
			if err := self.player.HasObjectOkEasy(cardid, carditem.CompoundNum); err == nil {
				have = true
			}
		}

		if have {
			self.player.SendErrInfo("err", "已拥有该虚空英雄")
			return
		}
	}

	if config.Sale != 0 && data.Done != 1 {
		self.player.SendErrInfo("err", "任务未完成")
		return
	}

	data.Done = 0

	if data.Times < config.Times {
		data.Times++
	}

	item := self.player.AddObjectLst(config.Items, config.Nums, "活动礼包", msg.ID, data.Times, 0)

	var backmsg S2C_ActivityGiftGetAward
	backmsg.Cid = MSG_ACTIVITY_GIFT_GET_AWARD
	backmsg.ID = msg.ID
	backmsg.Items = item
	backmsg.Times = data.Times
	self.player.Send(backmsg.Cid, backmsg)
}

//! 刷新处理  每日重置次数
func (self *ModActivityGift) OnRefresh() {
	nLen := len(self.San_ActivityGift.info)

	for i := nLen - 1; i >= 0; i-- {
		timeStamp := self.GetTimeStamp(self.San_ActivityGift.info[i].BoxId)
		if timeStamp != self.San_ActivityGift.info[i].TimeStamp {
			self.San_ActivityGift.info = append(self.San_ActivityGift.info[:i], self.San_ActivityGift.info[i+1:]...)
			continue
		}
	}

	self.ActivityGiftSendInfo([]byte{})
}

//! 刷新处理  每日重置次数
func (self *ModActivityGift) GetTimeStamp(id int) int64 {
	now := TimeServer()
	//获取刷新类型
	refreshTime := GetCsvMgr().GetActivityGiftConfigRefresh(id)
	switch refreshTime {
	case 1:
		timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
		if TimeServer().Hour() < 5 {
			timeStamp -= DAY_SECS
		}
		return timeStamp
	case 2:
		if now.Weekday() == time.Monday {
			if now.Hour() < 5 {
				return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix() - 7*DAY_SECS
			} else {
				return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
			}
		} else if now.Weekday() == time.Sunday {
			return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix() - 6*DAY_SECS
		} else {
			return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix() - int64((int(now.Weekday())-1)*DAY_SECS)
		}
	case 3:
		if now.Day() == 1 && now.Hour() < 5 {
			if now.Month() == 1 {
				return time.Date(now.Year()-1, 12, 1, 5, 0, 0, 0, now.Location()).Unix()
			} else {
				return time.Date(now.Year(), now.Month()-1, 1, 5, 0, 0, 0, now.Location()).Unix()
			}
		} else {
			return time.Date(now.Year(), now.Month(), 1, 5, 0, 0, 0, now.Location()).Unix()
		}
	}
	return 0
}

func (self *ModActivityGift) CheckOldTaskID() {
	nSize := len(self.San_ActivityGift.info)
	for i := nSize - 1; i >= 0; i-- {
		if self.San_ActivityGift.info[i].BoxId < 50 {
			self.San_ActivityGift.info = append(self.San_ActivityGift.info[:i], self.San_ActivityGift.info[i+1:]...)
		}
	}
}
