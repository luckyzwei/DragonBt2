package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

//PVE战报
type PveFight struct {
	FightId  int64  `json:"fight_id"`      // 战斗Id
	Side     int    `json:"side"`          // 1 进攻方 0 防守方
	Result   int    `json:"attack_result"` // 0 进攻方成功 其他防守方胜利
	Score    int64  `json:"score"`         //
	Uid      int64  `json:"uid"`           // Uid
	IconId   int    `json:"icon"`          // 头像Id
	Portrait int    `json:"portrait"`      // 头像框
	Name     string `json:"name"`          // 名字
	Level    int    `json:"level"`         // 等级
	Fight    int64  `json:"fight"`         // 战力
	Time     int64  `json:"time"`          // 发生的时间
}

type JS_ActivityBossFestivalInfo struct {
	Id          int         `json:"id"`          // 奖励ID
	Period      int         `json:"period"`      // 当前期数
	StartTime   int64       `json:"starttime"`   // 开始时间
	EndTime     int64       `json:"endtime"`     // 结束时间
	RewardTimes int         `json:"rewardtimes"` // 是否发奖
	Records     []*PveFight `json:"records"`     // 战报
}

//! 任务数据库
type San_ActivityBossFestival struct {
	Uid                      int64
	ActivityBossFestivalInfo string

	activityBossFestivalInfo *JS_ActivityBossFestivalInfo //! 任务信息
	DataUpdate
}

//! 任务
type ModActivityBossFestival struct {
	player                   *Player
	Sql_ActivityBossFestival San_ActivityBossFestival
	bossInfo                 *JS_FightInfo
}

func (self *ModActivityBossFestival) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_useractivitybossfestival` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_ActivityBossFestival, "san_useractivitybossfestival", self.player.ID)

	if self.Sql_ActivityBossFestival.Uid <= 0 {
		self.Sql_ActivityBossFestival.Uid = self.player.ID
		self.Sql_ActivityBossFestival.activityBossFestivalInfo = new(JS_ActivityBossFestivalInfo)
		self.Encode()
		InsertTable("san_useractivitybossfestival", &self.Sql_ActivityBossFestival, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_ActivityBossFestival.Init("san_useractivitybossfestival", &self.Sql_ActivityBossFestival, true)
}

//! 将数据库数据写入data
func (self *ModActivityBossFestival) Decode() {
	json.Unmarshal([]byte(self.Sql_ActivityBossFestival.ActivityBossFestivalInfo), &self.Sql_ActivityBossFestival.activityBossFestivalInfo)
}

//! 将data数据写入数据库
func (self *ModActivityBossFestival) Encode() {
	self.Sql_ActivityBossFestival.ActivityBossFestivalInfo = HF_JtoA(&self.Sql_ActivityBossFestival.activityBossFestivalInfo)
}

func (self *ModActivityBossFestival) OnGetOtherData() {

}

// 注册消息
func (self *ModActivityBossFestival) onReg(handlers map[string]func(body []byte)) {
	handlers["activitybossfestivalstart"] = self.ActivityBossFestivalStart
	handlers["activitybossfestivalresult"] = self.ActivityBossFestivalResult
	handlers["activitybossfestivalgetrecord"] = self.ActivityBossFestivalGetRecord
}

func (self *ModActivityBossFestival) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModActivityBossFestival) OnSave(sql bool) {
	self.Encode()
	self.Sql_ActivityBossFestival.Update(sql)
}

//每日任务刷新
func (self *ModActivityBossFestival) OnRefresh() {
	if self.Sql_ActivityBossFestival.activityBossFestivalInfo == nil {
		return
	}

	self.Sql_ActivityBossFestival.activityBossFestivalInfo.RewardTimes = LOGIC_FALSE
}

func (self *ModActivityBossFestival) Check() {
	activity := GetActivityMgr().GetActivity(ACT_BOSS_FESTIVAL)
	if activity == nil {
		return
	}
	self.Sql_ActivityBossFestival.activityBossFestivalInfo.Id = ACT_BOSS_FESTIVAL
	period := GetActivityMgr().getActN3(ACT_BOSS_FESTIVAL)

	if self.Sql_ActivityBossFestival.activityBossFestivalInfo.Period != period {
		//刷新
		self.Sql_ActivityBossFestival.activityBossFestivalInfo.Period = period
		self.Sql_ActivityBossFestival.activityBossFestivalInfo.RewardTimes = LOGIC_FALSE
		self.Sql_ActivityBossFestival.activityBossFestivalInfo.Records = make([]*PveFight, 0)
	}
	self.Sql_ActivityBossFestival.activityBossFestivalInfo.StartTime = HF_CalTimeForConfig(activity.info.Start, self.player.Sql_UserBase.Regtime)
	self.Sql_ActivityBossFestival.activityBossFestivalInfo.EndTime = self.Sql_ActivityBossFestival.activityBossFestivalInfo.StartTime + int64(activity.info.Continued) + int64(activity.info.Show)

	if self.bossInfo == nil {
		level := GetOfflineInfoMgr().GetMaxLevel(self.player.Sql_UserBase.Uid)
		self.bossInfo = GetRobotMgr().GetRobotByWorldLv(ACT_BOSS_FESTIVAL, level)
	}
}

func (self *ModActivityBossFestival) SendInfo() {
	activity := GetActivityMgr().GetActivity(ACT_BOSS_FESTIVAL)
	if activity == nil || activity.status.Status == ACTIVITY_STATUS_CLOSED {
		return
	}

	self.Check()

	var msg S2C_ActivityBossFestivalInfo
	msg.Cid = "activitybossfestivalinfo"
	msg.ActivityBossFestivalInfo = self.Sql_ActivityBossFestival.activityBossFestivalInfo
	msg.Level = GetOfflineInfoMgr().GetMaxLevel(self.player.Sql_UserBase.Uid)
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModActivityBossFestival) ActivityBossFestivalGetRecord(body []byte) {
	records := GetOfflineInfoMgr().GetRecord(self.Sql_ActivityBossFestival.activityBossFestivalInfo.Id, self.Sql_ActivityBossFestival.activityBossFestivalInfo.Period)

	var msgRel S2C_ActivityBossFestivalGetRecord
	msgRel.Cid = "activitybossfestivalgetrecord"
	msgRel.Records = records
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModActivityBossFestival) ActivityBossFestivalResult(body []byte) {

	if self.bossInfo == nil {
		return
	}

	now := TimeServer().Unix()
	if now < self.Sql_ActivityBossFestival.activityBossFestivalInfo.StartTime || now > self.Sql_ActivityBossFestival.activityBossFestivalInfo.EndTime {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ORDERHERO_THE_ACTIVITY_IS_OVER_AND"))
		return
	}

	if self.Sql_ActivityBossFestival.activityBossFestivalInfo.RewardTimes == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_AWARD_HAS_BEEN_RECEIVED"))
		return
	}

	var msg C2S_ActivityBossFestivalResult
	json.Unmarshal(body, &msg)

	config := GetCsvMgr().GetActivityBossConfig(self.Sql_ActivityBossFestival.activityBossFestivalInfo.Id, self.Sql_ActivityBossFestival.activityBossFestivalInfo.Period)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	if msg.Score > PER_BIT {
		msg.Score = PER_BIT
	}

	msg.BattleInfo.Id = GetFightMgr().GetFightInfoID()
	attackFightInfo := GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_ACTIVITY_BOSS_FESTIVAL)
	msg.BattleInfo.UserInfo[0].Name = attackFightInfo.Uname
	msg.BattleInfo.UserInfo[0].Icon = attackFightInfo.Iconid
	msg.BattleInfo.UserInfo[0].UnionName = attackFightInfo.UnionName

	fightRecord := &PveFight{}
	fightRecord.FightId = msg.BattleInfo.Id
	fightRecord.Side = 1
	fightRecord.Score = int64(msg.Score)
	if attackFightInfo != nil {
		fightRecord.Uid = attackFightInfo.Uid
		fightRecord.IconId = attackFightInfo.Iconid
		fightRecord.Portrait = attackFightInfo.Portrait
		fightRecord.Name = attackFightInfo.Uname
		fightRecord.Level = attackFightInfo.Level
		fightRecord.Fight = attackFightInfo.Deffight
	}
	fightRecord.Time = TimeServer().Unix()

	data2 := BattleRecord{}
	data2.Level = 0
	data2.Side = 1
	data2.Time = TimeServer().Unix()
	data2.Id = msg.BattleInfo.Id
	data2.LevelID = msg.BattleInfo.LevelID
	data2.Type = msg.BattleInfo.Type
	data2.RandNum = msg.BattleInfo.Random
	data2.FightInfo[0] = attackFightInfo
	data2.FightInfo[1] = self.bossInfo

	HMSetRedisEx("san_activitybossbattleinfo", msg.BattleInfo.Id, &msg.BattleInfo, DAY_SECS*15)
	HMSetRedisEx("san_activitybossbattlerecord", data2.Id, &data2, DAY_SECS*15)

	self.Sql_ActivityBossFestival.activityBossFestivalInfo.Records = append(self.Sql_ActivityBossFestival.activityBossFestivalInfo.Records, fightRecord)
	GetOfflineInfoMgr().AddRecord(self.Sql_ActivityBossFestival.activityBossFestivalInfo.Id, self.Sql_ActivityBossFestival.activityBossFestivalInfo.Period, fightRecord)

	self.Sql_ActivityBossFestival.activityBossFestivalInfo.RewardTimes = LOGIC_TRUE
	times := msg.Score / config.HealthUnit
	if times < 0 {
		times = 0
	}
	items := self.player.AddObjectSimple(config.GetItem, config.GetNum*times, "活动BOSS", msg.Score, 0, 0)

	var msgRel S2C_ActivityBossFestivalResult
	msgRel.Cid = "activitybossfestivalresult"
	msgRel.GetItems = items
	msgRel.RewardTimes = self.Sql_ActivityBossFestival.activityBossFestivalInfo.RewardTimes
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModActivityBossFestival) ActivityBossFestivalStart(body []byte) {

	now := TimeServer().Unix()
	if now < self.Sql_ActivityBossFestival.activityBossFestivalInfo.StartTime || now > self.Sql_ActivityBossFestival.activityBossFestivalInfo.EndTime {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ORDERHERO_THE_ACTIVITY_IS_OVER_AND"))
		return
	}

	if self.Sql_ActivityBossFestival.activityBossFestivalInfo.RewardTimes == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_AWARD_HAS_BEEN_RECEIVED"))
		return
	}

	level := GetOfflineInfoMgr().GetMaxLevel(self.player.Sql_UserBase.Uid)
	self.bossInfo = GetRobotMgr().GetRobotByWorldLv(ACT_BOSS_FESTIVAL, level)

	if self.bossInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_DATA_ABNORMITY"))
		return
	}

	var msgRel S2C_ActivityBossFestivalStart
	msgRel.Cid = "activitybossfestivalstart"
	msgRel.BossFightInfo = self.bossInfo
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

// 获取战报信息
func (self *ModActivityBossFestival) GetBattleInfo(fightID int64) *BattleInfo {
	var battleInfo BattleInfo
	value, flag, err := HGetRedisEx(`san_ActivityBossFestivalbattleinfo`, fightID, fmt.Sprintf("%d", fightID))
	if err != nil || !flag {
		return GetServer().DBUser.GetBattleInfo(fightID)
	}
	if flag {
		err := json.Unmarshal([]byte(value), &battleInfo)
		if err != nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
			return &battleInfo
		}
	}

	if battleInfo.Id != 0 {
		return &battleInfo
	}
	return nil
}

// 获取战报信息
func (self *ModActivityBossFestival) GetBattleRecord(fightID int64) *BattleRecord {
	var battleRecord BattleRecord
	value, flag, err := HGetRedisEx(`san_ActivityBossFestivalbattlerecord`, fightID, fmt.Sprintf("%d", fightID))
	if err != nil || !flag {
		return GetServer().DBUser.GetBattleRecord(fightID)
	}
	if flag {
		err := json.Unmarshal([]byte(value), &battleRecord)
		if err != nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
			return &battleRecord
		}
	}
	if battleRecord.Id != 0 {
		return &battleRecord
	}
	return nil
}

func (self *ModActivityBossFestival) GMReset() {
	activity := GetActivityMgr().GetActivity(ACT_BOSS_FESTIVAL)
	if activity == nil || activity.status.Status == ACTIVITY_STATUS_CLOSED {
		return
	}

	self.Check()
	self.Sql_ActivityBossFestival.activityBossFestivalInfo.RewardTimes = LOGIC_FALSE

	var msg S2C_ActivityBossFestivalInfo
	msg.Cid = "activitybossfestivalinfo"
	msg.ActivityBossFestivalInfo = self.Sql_ActivityBossFestival.activityBossFestivalInfo
	msg.Level = GetOfflineInfoMgr().GetMaxLevel(self.player.Sql_UserBase.Uid)
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}
