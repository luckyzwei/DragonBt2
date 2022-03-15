package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

// 进度
type JS_FundTaskInfo struct {
	Taskid    int   `json:"taskid"`    // 任务Id
	Tasktypes int   `json:"tasktypes"` // 任务类型
	Plan      int   `json:"plan"`      // 进度
	Finish    int   `json:"finish"`    // 是否完成
	Pickup    int   `json:"pickup"`    // 是否领取奖励
	StartTime int64 `json:"starttime"` // 领取开始时间
}

type JS_NewFundInfo struct {
	N3       int               `json:"n3"`     //! 期数
	NGroup   int               `json:"ngroup"` //! 组
	TaskInfo []JS_FundTaskInfo `json:"taskinfo"`
}

//!
type San_Fund struct {
	Uid      int64
	N3       int    //! 期数
	NGroup   int    //! 组
	TaskInfo string //! 任务信息
	FundInfo string //! 基金信息

	taskInfo []JS_FundTaskInfo
	fundInfo *JS_NewFundInfo
	DataUpdate
}

//! 超值基金
type ModFund struct {
	player   *Player
	Sql_Fund *San_Fund
}

func (self *ModFund) OnGetData(player *Player) {
	self.player = player
}

func (self *ModFund) OnGetOtherData() {
	self.Sql_Fund = new(San_Fund)
	sql := fmt.Sprintf("select * from `san_userfund` where `uid` = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, self.Sql_Fund, "san_userfund", self.player.ID)
	if self.Sql_Fund.Uid <= 0 {
		self.Sql_Fund.Uid = self.player.ID
		self.Sql_Fund.taskInfo = make([]JS_FundTaskInfo, 0)
		self.Encode()
		InsertTable("san_userfund", self.Sql_Fund, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_Fund.Init("san_userfund", self.Sql_Fund, true)
}

func (self *ModFund) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	// 领取奖励
	case "getfund":
		var msg C2S_GetFund
		json.Unmarshal(body, &msg)
		self.GetFund(&msg)
		return true
	}

	return false
}

func (self *ModFund) OnSave(sql bool) {
	if self.Sql_Fund != nil {
		self.Encode()
		self.Sql_Fund.Update(sql)
	}

}

func (self *ModFund) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_Fund.TaskInfo), &self.Sql_Fund.taskInfo)
	json.Unmarshal([]byte(self.Sql_Fund.FundInfo), &self.Sql_Fund.fundInfo)
}

func (self *ModFund) Encode() { //! 将data数据写入数据库
	self.Sql_Fund.TaskInfo = HF_JtoA(&self.Sql_Fund.taskInfo)
	self.Sql_Fund.FundInfo = HF_JtoA(&self.Sql_Fund.fundInfo)
}

func (self *ModFund) Check() {
	if self.Sql_Fund.taskInfo == nil {
		self.Sql_Fund.taskInfo = make([]JS_FundTaskInfo, 0)
	}

	isOpen, index := GetActivityMgr().JudgeOpenAllIndex(ACT_FUND_MIN, ACT_FUND_MAX)
	if !isOpen {
		self.SendRewards()
		return
	}

	N3 := GetActivityMgr().getActN3(index)
	group := GetActivityMgr().getActN4(index)
	if N3 != self.Sql_Fund.N3 || group != self.Sql_Fund.NGroup {
		self.SendRewards()

		self.Sql_Fund.taskInfo = make([]JS_FundTaskInfo, 0)
		self.Sql_Fund.N3 = N3
		self.Sql_Fund.NGroup = group

		startTime := self.player.GetModule("activity").(*ModActivity).GetActivityStart(index)
		for _, v := range GetCsvMgr().FundConfigMap {
			if v.Group != self.Sql_Fund.NGroup {
				continue
			}
			if v.Type != 1 && v.Type != 2 {
				continue
			}
			var task JS_FundTaskInfo
			task.Taskid = v.Id
			task.Tasktypes = v.TaskTypes
			task.StartTime = startTime + (v.Day-1)*DAY_SECS
			self.Sql_Fund.taskInfo = append(self.Sql_Fund.taskInfo, task)
		}
	}
}

func (self *ModFund) CheckFund() {
	if self.Sql_Fund.fundInfo == nil {
		self.Sql_Fund.fundInfo = new(JS_NewFundInfo)
	}

	isOpen, index := GetActivityMgr().JudgeOpenAllIndex(ACT_NEW_FUND, ACT_NEW_FUND)
	if !isOpen {
		return
	}

	N3 := GetActivityMgr().getActN3(index)
	group := GetActivityMgr().getActN4(index)
	if N3 != self.Sql_Fund.fundInfo.N3 || group != self.Sql_Fund.fundInfo.NGroup {
		self.Sql_Fund.fundInfo.TaskInfo = make([]JS_FundTaskInfo, 0)
		self.Sql_Fund.fundInfo.N3 = N3
		self.Sql_Fund.fundInfo.NGroup = group

		startTime := self.player.GetModule("activity").(*ModActivity).GetActivityStart(index)
		for _, v := range GetCsvMgr().FundConfigMap {
			if v.Group != self.Sql_Fund.fundInfo.NGroup {
				continue
			}
			if v.Type != 3 && v.Type != 4 {
				continue
			}
			var task JS_FundTaskInfo
			task.Taskid = v.Id
			task.Tasktypes = v.TaskTypes
			task.StartTime = startTime + (v.Day-1)*DAY_SECS
			self.Sql_Fund.fundInfo.TaskInfo = append(self.Sql_Fund.fundInfo.TaskInfo, task)
		}
	}
}

func (self *ModFund) SendRewards() {
	itemMap := make(map[int]*Item)
	for i := 0; i < len(self.Sql_Fund.taskInfo); i++ {
		if self.Sql_Fund.taskInfo[i].Finish == LOGIC_FALSE {
			continue
		}
		if self.Sql_Fund.taskInfo[i].Pickup == LOGIC_TRUE {
			continue
		}
		config := GetCsvMgr().FundConfigMap[self.Sql_Fund.taskInfo[i].Taskid]
		if config == nil {
			continue
		}
		self.Sql_Fund.taskInfo[i].Pickup = LOGIC_TRUE
		AddItemMapHelper(itemMap, config.Item, config.Num)
	}
	//发送邮件
	if len(itemMap) > 0 {
		mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_GET_FUND]
		if ok {
			itemLst := make([]PassItem, 0)
			for _, v := range itemMap {
				itemLst = append(itemLst, PassItem{ItemID: v.ItemId, Num: v.ItemNum})
			}
			self.player.GetModule("mail").(*ModMail).AddMail(1,
				1, 0, mailConfig.Mailtitle, mailConfig.Mailtxt, GetCsvMgr().GetText("STR_SYS"), itemLst, false, 0)
		}
	}
}

// 发送礼包信息
func (self *ModFund) SendInfo() {
	self.Check()
	self.CheckFund()
	var msg S2C_FundInfo
	msg.Cid = "fundinfo"
	msg.TaskInfo = self.Sql_Fund.taskInfo
	msg.Config = GetCsvMgr().FundConfigMap
	msg.FundInfo = self.Sql_Fund.fundInfo
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}

func (self *ModFund) HandleRecharge(grade int) int {
	rel := 0
	for i := 0; i < len(self.Sql_Fund.taskInfo); i++ {
		task := &self.Sql_Fund.taskInfo[i]
		config := GetCsvMgr().FundConfigMap[task.Taskid]
		if config == nil {
			continue
		}
		if config.Pay != grade {
			continue
		}
		task.Finish = LOGIC_TRUE
		rel = config.Type
	}
	for i := 0; i < len(self.Sql_Fund.fundInfo.TaskInfo); i++ {
		task := &self.Sql_Fund.fundInfo.TaskInfo[i]
		config := GetCsvMgr().FundConfigMap[task.Taskid]
		if config == nil {
			continue
		}
		if config.Pay != grade {
			continue
		}
		task.Finish = LOGIC_TRUE
		rel = config.Type
	}
	self.SendInfo()
	return rel
}

func (self *ModFund) GetFund(msg *C2S_GetFund) {
	task := self.GetTask(msg.Id)
	if task == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	if task.Pickup == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_AWARD_HAS_BEEN_RECEIVED"))
		return
	}

	if task.Finish != LOGIC_TRUE || task.StartTime > TimeServer().Unix() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_MILITARY_ITS_NOT_TIME_TO_COLLECT"))
		return
	}

	config := GetCsvMgr().FundConfigMap[task.Taskid]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	items := make([]PassItem, 0)
	if config.Type == 1 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FUND_1_GET, task.Taskid, 0, 0, "领取勇者基金奖励", 0, self.player.Sql_UserBase.Vip, self.player)
		items = self.player.AddObjectLst(config.Item, config.Num, "领取勇者基金奖励", task.Taskid, 0, 0)
	} else if config.Type == 2 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FUND_2_GET, task.Taskid, 0, 0, "领取至尊基金奖励", 0, self.player.Sql_UserBase.Vip, self.player)
		items = self.player.AddObjectLst(config.Item, config.Num, "领取至尊基金奖励", task.Taskid, 0, 0)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FUND_3_GET, task.Taskid, 0, 0, "领取无敌基金奖励", 0, self.player.Sql_UserBase.Vip, self.player)
		items = self.player.AddObjectLst(config.Item, config.Num, "领取无敌基金奖励", task.Taskid, 0, 0)
	}

	task.Pickup = LOGIC_TRUE
	//发送奖励

	var msgRel S2C_GetFund
	msgRel.Cid = "getfund"
	msgRel.TaskInfo = task
	msgRel.GetItems = items
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFund) GetTask(id int) *JS_FundTaskInfo {

	for i := 0; i < len(self.Sql_Fund.taskInfo); i++ {
		if self.Sql_Fund.taskInfo[i].Taskid == id {
			return &self.Sql_Fund.taskInfo[i]
		}
	}
	for i := 0; i < len(self.Sql_Fund.fundInfo.TaskInfo); i++ {
		if self.Sql_Fund.fundInfo.TaskInfo[i].Taskid == id {
			return &self.Sql_Fund.fundInfo.TaskInfo[i]
		}
	}
	return nil
}
