package game

import (
	"encoding/json"
	"fmt"
	"time"
)

// 地牢系统
type San_UserPit struct {
	Uid       int64  `json:"uid"` // 玩家Id
	PitInfo   string `json:"pitinfo"`
	MaxKeyId  int    `json:"maxkeyid"`
	FirstInfo string `json:"firstinfo"`

	pitinfo   []*PitInfo // 推图信息
	firstInfo []int      // 推图信息
	DataUpdate
}

type ModPit struct {
	player      *Player
	Sql_UserPit San_UserPit
}

type PitInfo struct {
	Id          int           `json:"id"`          //ID
	PitId       int           `json:"pitid"`       //ID生成地牢关卡
	PitType     int           `json:"pittype"`     //TYPE地牢产生规则
	PitKeyId    int           `json:"pitkeyid"`    //地牢唯一标识
	PitEvent    []*PitEvent   `json:"pitevent"`    // 推图信息
	EndTime     int64         `json:"endtime"`     //结束时间
	State       int           `json:"state"`       //地牢状态
	FinishTimes int           `json:"finishtimes"` //完成
	AllTimes    int           `json:"alltimes"`    //总共次数
	FightInfo   *JS_FightInfo `json:"fightinfo"`   //战斗信息
	Buff        []int         `json:"buff"`        //BUFF
}

type PitEvent struct {
	Id    int `json:"id"`    //自己的ID
	State int `json:"state"` //! 事件状态
}

const (
	PIT_STATE_CANT_FINISH = 0
	PIT_STATE_CAN_FINISH  = 1
	PIT_STATE_FINISH      = 2
)

//对应dungeons_map表dungeons_type
const (
	PIT_TYPE_OPEN_ITEM         = 1 //通过挂机时随机获得道具后开启。持续时间2小时，通关后关闭
	PIT_TYPE_OPEN_TIME         = 2 //固定时间开启，需要消耗道具才可以进入，通关后地牢不关闭
	PIT_TYPE_OPEN_HERO         = 3 //固定时间开启，阵容中需要有特定英雄才可以进入，通关后地牢关闭）
	PIT_TYPE_OPEN_WEEK         = 4 //每周1凌晨5点开启，每周通关次数3次，下周一凌晨4点59分关闭
	PIT_TYPE_OPEN_WEEK_VIP     = 5 //VIP5每周1凌晨5点开启，每周通关次数3次，下周一凌晨4点59分关闭
	PIT_TYPE_OPEN_WEEK_VIP_EXT = 6 //VIP每周1凌晨5点开启，每周通关次数3次，下周一凌晨4点59分关闭
	PIT_TYPE_OPEN_MONTH        = 7 //每月1号凌晨5点开启，每月通关次数7次，下个月1号凌晨4:59分关闭
)

const (
	PIT_TASK_START    = 0
	PIT_TASK_DOOR     = 1
	PIT_TASK_MONSTER  = 2 //怪物
	PIT_TASK_BOX      = 3
	PIT_TASK_WISH     = 6 //许愿池
	PIT_TASK_LOCK_BOX = 8 //上锁宝箱
)

func (self *ModPit) OnGetData(player *Player) {
	self.player = player
	sql := fmt.Sprintf("select * from `san_userpit` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_UserPit, "san_userpit", self.player.ID)

	if self.Sql_UserPit.Uid <= 0 {
		self.Sql_UserPit.Uid = self.player.ID
		self.Sql_UserPit.pitinfo = make([]*PitInfo, 0)
		self.Encode()
		InsertTable("san_userpit", &self.Sql_UserPit, 0, true)
		self.Sql_UserPit.Init("san_userpit", &self.Sql_UserPit, true)
	} else {
		self.Decode()
		self.Sql_UserPit.Init("san_userpit", &self.Sql_UserPit, true)
	}

	//更新过图相关的任务配置
	self.Check()
}

func (self *ModPit) OnSave(sql bool) {
	self.Encode()
	self.Sql_UserPit.Update(sql)
}

func (self *ModPit) OnGetOtherData() {

}

func (self *ModPit) Decode() { // 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_UserPit.PitInfo), &self.Sql_UserPit.pitinfo)
	json.Unmarshal([]byte(self.Sql_UserPit.FirstInfo), &self.Sql_UserPit.firstInfo)
}

func (self *ModPit) Encode() { // 将data数据写入数据库
	self.Sql_UserPit.PitInfo = HF_JtoA(&self.Sql_UserPit.pitinfo)
	self.Sql_UserPit.FirstInfo = HF_JtoA(&self.Sql_UserPit.firstInfo)
}

func (self *ModPit) Check() { // 将data数据写入数据库
	if self.Sql_UserPit.firstInfo == nil {
		self.Sql_UserPit.firstInfo = make([]int, 0)
	}

	//检查常驻地牢的更新
	self.CheckUpdate()

	//检查过期地牢
	isHasPassed := false
	now := TimeServer().Unix()
	for _, v := range self.Sql_UserPit.pitinfo {
		if v.EndTime < now {
			isHasPassed = true
			break
		}
	}

	if isHasPassed {
		newInfo := make([]*PitInfo, 0)
		for _, v := range self.Sql_UserPit.pitinfo {
			if v.EndTime > now {
				newInfo = append(newInfo, v)
			}
		}
		self.Sql_UserPit.pitinfo = newInfo
	}

	//判断限时地牢
	self.CheckTimePit()
}

//检查常驻地牢
func (self *ModPit) CheckUpdate() {

	now := TimeServer().Unix()
	for i := PIT_TYPE_OPEN_WEEK; i <= PIT_TYPE_OPEN_MONTH; i++ {
		//先检查数据里是否有这个本
		isFind := false
		for j := 0; j < len(self.Sql_UserPit.pitinfo); j++ {
			if self.Sql_UserPit.pitinfo[j].PitType == i {
				isFind = true
				//判断是否需要更新
				if self.Sql_UserPit.pitinfo[j].EndTime < now {
					self.Sql_UserPit.pitinfo[j].updateInfo()
				}
			}
		}

		if !isFind {
			switch i {
			case PIT_TYPE_OPEN_WEEK:
				pitInfo := self.getNewPitInfo(PIT_TYPE_OPEN_WEEK)
				self.Sql_UserPit.pitinfo = append(self.Sql_UserPit.pitinfo, pitInfo)
			case PIT_TYPE_OPEN_WEEK_VIP:
				if self.player.Sql_UserBase.Vip >= 5 {
					pitInfo := self.getNewPitInfo(PIT_TYPE_OPEN_WEEK_VIP)
					self.Sql_UserPit.pitinfo = append(self.Sql_UserPit.pitinfo, pitInfo)
				}
			case PIT_TYPE_OPEN_WEEK_VIP_EXT:
				if self.player.Sql_UserBase.Vip >= 5 {
					pitInfo := self.getNewPitInfo(PIT_TYPE_OPEN_WEEK_VIP_EXT)
					self.Sql_UserPit.pitinfo = append(self.Sql_UserPit.pitinfo, pitInfo)
				}
			case PIT_TYPE_OPEN_MONTH:
				pitInfo := self.getNewPitInfo(PIT_TYPE_OPEN_MONTH)
				self.Sql_UserPit.pitinfo = append(self.Sql_UserPit.pitinfo, pitInfo)
			}
		}
	}
}

//检查限时地牢
func (self *ModPit) CheckTimePit() {

	now := TimeServer().Unix()
	for i := PIT_TYPE_OPEN_TIME; i <= PIT_TYPE_OPEN_HERO; i++ {
		//先检查数据里是否有这个本
		isFind := false
		for j := 0; j < len(self.Sql_UserPit.pitinfo); j++ {
			if self.Sql_UserPit.pitinfo[j].PitType == i {
				isFind = true
			}
		}

		if !isFind {
			//必须先判断活动是否有
			config := GetCsvMgr().PitMapMap[i]
			if config == nil {
				return
			}

			for j := 0; j < len(config); j++ {
				endTime := config[j].OpenTime + config[j].OpenDuration
				if now >= config[j].OpenTime && now < endTime {
					pitInfo := self.getNewPitInfo(i)
					self.Sql_UserPit.pitinfo = append(self.Sql_UserPit.pitinfo, pitInfo)
					break
				}
			}
		}
	}
}

func (self *ModPit) getNewPitInfo(pitType int) *PitInfo {
	pitInfo := new(PitInfo)
	pitInfo.PitKeyId = self.getMaxPitKeyId()
	pitInfo.PitType = pitType
	pitInfo.updateInfo()
	return pitInfo
}

func (self *ModPit) getMaxPitKeyId() int {
	self.Sql_UserPit.MaxKeyId++
	return self.Sql_UserPit.MaxKeyId
}

func (self *ModPit) SendInfo(body []byte) {
	var msg C2S_PitInfo
	json.Unmarshal(body, &msg)

	var msgRel S2C_PitInfo
	msgRel.Cid = "pitinfo"
	msgRel.PitKeyId = msg.PitKeyId
	msgRel.Pitinfo = self.getPitInfo(msg.PitKeyId)
	self.player.SendMsg("pitinfo", HF_JtoB(&msgRel))
}

func (self *ModPit) SendInfoAll(body []byte) {
	self.Check()

	var msgRel S2C_PitInfoAll
	msgRel.Cid = "pitinfoall"

	for i := 0; i < len(self.Sql_UserPit.pitinfo); i++ {
		var pitInfoShow PitInfoShow
		pitInfoShow.Id = self.Sql_UserPit.pitinfo[i].Id
		pitInfoShow.PitId = self.Sql_UserPit.pitinfo[i].PitId
		pitInfoShow.State = self.Sql_UserPit.pitinfo[i].State
		pitInfoShow.PitType = self.Sql_UserPit.pitinfo[i].PitType
		pitInfoShow.EndTime = self.Sql_UserPit.pitinfo[i].EndTime
		pitInfoShow.PitKeyId = self.Sql_UserPit.pitinfo[i].PitKeyId
		pitInfoShow.FinishTimes = self.Sql_UserPit.pitinfo[i].FinishTimes
		pitInfoShow.AllTimes = self.Sql_UserPit.pitinfo[i].AllTimes

		msgRel.PitInfoShow = append(msgRel.PitInfoShow, pitInfoShow)
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModPit) OnMsg(ctrl string, body []byte) bool {
	return false
}
func (self *ModPit) onReg(handlers map[string]func(body []byte)) {
	handlers["pitinfo"] = self.SendInfo
	handlers["pitfinishevent"] = self.PitFinishEvent
	handlers["pitinfoall"] = self.SendInfoAll
	handlers["pitstart"] = self.StartPit
}

func (self *ModPit) getPitInfo(keyId int) *PitInfo {

	for i := 0; i < len(self.Sql_UserPit.pitinfo); i++ {
		if self.Sql_UserPit.pitinfo[i].PitKeyId == keyId {
			return self.Sql_UserPit.pitinfo[i]
		}
	}
	return nil
}

//! 完成城池事件
func (self *ModPit) PitFinishEvent(body []byte) {
	var msg C2S_PitFinishEvent
	json.Unmarshal(body, &msg)

	config := GetCsvMgr().PitConfigMap[msg.EventId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_CONFIGURATION_ERROR"))
		return
	}

	pitInfo := self.getPitInfo(msg.PitKeyId)
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PIT_NOT_EXIST"))
		return
	}

	for _, v := range pitInfo.PitEvent {
		if v.Id != msg.EventId {
			continue
		}
		if v.State == PIT_STATE_FINISH {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_KINGTASK_THE_TASK_HAS_BEEN_COMPLETED"))
			return
		}

		outitem := make([]PassItem, 0)
		costitem := make([]PassItem, 0)
		switch config.ThingType {
		case PIT_TASK_DOOR: //! 传送门
			pitInfo.GoThrough()
			item := self.player.AddObjectSimple(91000002, 500, "宝石合成", 0, 0, 0)
			outitem = append(outitem, item...)
		case PIT_TASK_MONSTER: //! 怪物
			configLevel := GetCsvMgr().LevelConfigMap[config.BattleId]
			if configLevel == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PIT_CONFIG_ERROR"))
				return
			}
			items := self.player.GetModule("pass").(*ModPass).GetDropItem(configLevel.LevelId, false)
			for i := 0; i < len(items); i++ {
				if items[i].Num > 0 {
					itemID, num := self.player.AddObject(items[i].ItemID, items[i].Num, v.Id, 0, 0, "地牢怪物")
					outitem = append(outitem, PassItem{ItemID: itemID, Num: num})
				}
			}
		case PIT_TASK_BOX: //! 宝箱
			for i := 0; i < len(config.Items); i++ {
				if config.Items[i] > 0 {
					itemID, num := self.player.AddObject(config.Items[i], config.Nums[i], v.Id, 0, 0, "地牢宝箱开启")
					outitem = append(outitem, PassItem{ItemID: itemID, Num: num})
				}
			}
		case PIT_TASK_WISH: //! 许愿池
			configLevel := GetCsvMgr().PitBuffMap[v.Id]
			if configLevel == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PIT_CONFIG_ERROR"))
				return
			}
			if msg.Option <= 0 || msg.Option > len(configLevel.BuffId) {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PIT_MSG_OPTION_ERROR"))
				return
			}
			pitInfo.Buff = append(pitInfo.Buff, configLevel.BuffId[msg.Option-1])
		case PIT_TASK_LOCK_BOX: //! 上锁宝箱
			configLevel := GetCsvMgr().PitBoxMap[v.Id]
			if configLevel == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PIT_CONFIG_ERROR"))
				return
			}
			if msg.Option <= 0 || msg.Option > len(configLevel.ItemId) {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PIT_MSG_OPTION_ERROR"))
				return
			}
			realIndex := msg.Option - 1
			if realIndex == 0 {
				if err := self.player.HasObjectOkEasy(configLevel.BoxId[realIndex], 1); err != nil {
					self.player.SendErr(err.Error())
					return
				}
				cost := self.player.RemoveObjectSimple(configLevel.BoxId[realIndex], 1, "地牢开箱", 0, 0, 0)
				costitem = append(costitem, cost...)

				for i := 0; i < len(configLevel.ItemId); i++ {
					if configLevel.ItemId[i] > 0 {
						itemID, num := self.player.AddObject(configLevel.ItemId[i], configLevel.ItemNum[i], v.Id, 0, 0, "地牢宝箱开启")
						outitem = append(outitem, PassItem{ItemID: itemID, Num: num})
					}
				}
			} else if realIndex == 1 {
				rand := HF_GetRandom(3)
				itemID, num := self.player.AddObject(configLevel.ItemId[rand], configLevel.ItemNum[rand], v.Id, 0, 0, "地牢宝箱开启")
				outitem = append(outitem, PassItem{ItemID: itemID, Num: num})
			} else if realIndex == 2 {
				for i := 0; i < len(configLevel.ItemId); i++ {
					if configLevel.ItemId[i] > 0 {
						itemID, num := self.player.AddObject(configLevel.ItemId[i], configLevel.ItemNum[i], v.Id, 0, 0, "地牢宝箱开启")
						outitem = append(outitem, PassItem{ItemID: itemID, Num: num})
					}
				}
			}
		}
		if config.ThingType != PIT_TASK_DOOR {
			v.State = PIT_STATE_FINISH
		}
		var msgRel S2C_PitFinishEvent
		msgRel.Cid = "pitfinishevent"
		msgRel.PitKeyId = msg.PitKeyId
		msgRel.PitEvent = v
		msgRel.Item = outitem
		msgRel.Cost = costitem
		msgRel.Buff = pitInfo.Buff
		self.player.SendMsg("pitfinishevent", HF_JtoB(&msgRel))
		return
	}
}

//开始地牢
func (self *ModPit) StartPit(body []byte) {
	var msg C2S_PitStart
	json.Unmarshal(body, &msg)

	pitInfo := self.getPitInfo(msg.PitKeyId)

	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PIT_NOT_EXIST"))
		return
	}

	fightInfo := GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, msg.TeamType)
	if fightInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PIT_TEAMTYPE_ERROR"))
		return
	}

	pitInfo.FightInfo = fightInfo
	pitInfo.State = PIT_STATE_CAN_FINISH

	var msgRel S2C_PitInfo
	msgRel.Cid = "pitstart"
	msgRel.PitKeyId = msg.PitKeyId
	msgRel.Pitinfo = self.getPitInfo(msg.PitKeyId)
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

//更新地牢信息,只在副本过期或新建副本的时候调用
func (self *PitInfo) updateInfo() {
	now := TimeServer().Unix()
	switch self.PitType {
	case PIT_TYPE_OPEN_TIME:
		config := GetCsvMgr().PitMapMap[PIT_TYPE_OPEN_TIME]
		if config == nil {
			return
		}

		for j := 0; j < len(config); j++ {
			endTime := config[j].OpenTime + config[j].OpenDuration
			if now >= config[j].OpenTime && now < endTime {
				self.init(config[j], endTime)
				break
			}
		}
	case PIT_TYPE_OPEN_HERO:
		config := GetCsvMgr().PitMapMap[PIT_TYPE_OPEN_HERO]
		if config == nil {
			return
		}

		for j := 0; j < len(config); j++ {
			endTime := config[j].OpenTime + config[j].OpenDuration
			if now >= config[j].OpenTime && now < endTime {
				self.init(config[j], endTime)
				break
			}
		}
	case PIT_TYPE_OPEN_WEEK:
		config := GetCsvMgr().PitMapMap[PIT_TYPE_OPEN_WEEK]
		if config == nil {
			return
		}
		//当前
		nowTime := TimeServer()
		offset := int(time.Monday - nowTime.Weekday())
		if offset > 0 {
			offset = -6
		}
		nowStart := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 5, 0, 0, 0, time.Local).AddDate(0, 0, offset)
		//开服
		openTime := time.Unix(GetServer().GetOpenServer(), 0)
		offset = int(time.Monday - openTime.Weekday())
		if offset > 0 {
			offset = -6
		}
		openStart := time.Date(openTime.Year(), openTime.Month(), openTime.Day(), 5, 0, 0, 0, time.Local).AddDate(0, 0, offset)

		dis := (nowStart.Unix() - openStart.Unix()) % (DAY_SECS * 7)
		index := int(dis) % len(config)

		startTime := nowStart.Unix()
		endTime := startTime + (DAY_SECS * 7)
		self.init(config[index], endTime)
	case PIT_TYPE_OPEN_WEEK_VIP:
		config := GetCsvMgr().PitMapMap[PIT_TYPE_OPEN_WEEK_VIP]
		if config == nil {
			return
		}
		//当前
		nowTime := TimeServer()
		offset := int(time.Monday - nowTime.Weekday())
		if offset > 0 {
			offset = -6
		}
		nowStart := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 5, 0, 0, 0, time.Local).AddDate(0, 0, offset)
		//开服
		openTime := time.Unix(GetServer().GetOpenServer(), 0)
		offset = int(time.Monday - openTime.Weekday())
		if offset > 0 {
			offset = -6
		}
		openStart := time.Date(openTime.Year(), openTime.Month(), openTime.Day(), 5, 0, 0, 0, time.Local).AddDate(0, 0, offset)

		dis := (nowStart.Unix() - openStart.Unix()) % (DAY_SECS * 7)
		index := int(dis) % len(config)

		startTime := nowStart.Unix()
		endTime := startTime + (DAY_SECS * 7)
		self.init(config[index], endTime)
	case PIT_TYPE_OPEN_WEEK_VIP_EXT:
		config := GetCsvMgr().PitMapMap[PIT_TYPE_OPEN_WEEK_VIP_EXT]
		if config == nil {
			return
		}
		//当前
		nowTime := TimeServer()
		offset := int(time.Monday - nowTime.Weekday())
		if offset > 0 {
			offset = -6
		}
		nowStart := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 5, 0, 0, 0, time.Local).AddDate(0, 0, offset)
		//开服
		openTime := time.Unix(GetServer().GetOpenServer(), 0)
		offset = int(time.Monday - openTime.Weekday())
		if offset > 0 {
			offset = -6
		}
		openStart := time.Date(openTime.Year(), openTime.Month(), openTime.Day(), 5, 0, 0, 0, time.Local).AddDate(0, 0, offset)

		dis := (nowStart.Unix() - openStart.Unix()) % (DAY_SECS * 7)
		index := int(dis) % len(config)

		startTime := nowStart.Unix()
		endTime := startTime + (DAY_SECS * 7)
		self.init(config[index], endTime)
	case PIT_TYPE_OPEN_MONTH:
		config := GetCsvMgr().PitMapMap[PIT_TYPE_OPEN_MONTH]
		if config == nil {
			return
		}
		//当前
		nowYear, nowMonth, _ := TimeServer().Date()
		nowMonthOne := time.Date(nowYear, nowMonth, 1, 0, 0, 0, 0, time.Local)

		//开服
		openYear, openMonth, _ := time.Unix(GetServer().GetOpenServer(), 0).Date()
		openMonthOne := time.Date(openYear, openMonth, 1, 0, 0, 0, 0, time.Local)

		dis := (nowMonthOne.Year()-openMonthOne.Year())*12 + nowMonthOne.Year() - openMonthOne.Year()
		index := int(dis) % len(config)

		endTime := nowMonthOne.AddDate(0, 1, 0).Unix()
		self.init(config[index], endTime)
	}
}

//更新地牢信息
func (self *PitInfo) reSetEvent() {
	config := GetCsvMgr().PitConfigMap
	if config == nil {
		LogError("reSetEvent is error")
		return
	}
	self.PitEvent = make([]*PitEvent, 0)
	for _, v := range GetCsvMgr().PitConfigMap {
		if v.PitId == self.PitId {
			self.PitEvent = append(self.PitEvent, &PitEvent{Id: v.Id, State: PIT_STATE_CANT_FINISH})
		}
	}
}

func (self *PitInfo) GoThrough() {
	for _, valueReset := range self.PitEvent {
		valueReset.State = PIT_STATE_CANT_FINISH
	}
	self.Buff = make([]int, 0)
	self.FinishTimes++
	self.State = CANTFINISH
}

func (self *PitInfo) init(config *PitMap, endTime int64) {
	self.Id = config.Id
	self.PitId = config.PitId
	self.State = PIT_STATE_CANT_FINISH
	self.AllTimes = config.PassTime
	self.EndTime = endTime
	self.Buff = make([]int, 0)

	self.reSetEvent()
}
