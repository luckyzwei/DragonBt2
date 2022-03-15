//目标系统玩法（爵位） 20190604 by zy
package game

import (
	"encoding/json"
	"fmt"
)

// 进度
type JS_VipRechargeInfo struct {
	Taskid    int `json:"taskid"`    // 任务Id
	Tasktypes int `json:"tasktypes"` // 任务类型
	Plan      int `json:"plan"`      // 进度
	Finish    int `json:"finish"`    // 完成次数
	PickUp    int `json:"pickup"`    // 领取次数
}

//! 任务数据库
type San_VipRecharge struct {
	Uid      int64
	Taskinfo string
	NextTime int64

	taskinfo map[int]*JS_VipRechargeInfo //! 任务信息
	DataUpdate
}

//! 任务
type ModVipRecharge struct {
	player          *Player
	Sql_VipRecharge San_VipRecharge
	init            bool
}

func (self *ModVipRecharge) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_userviprecharge` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_VipRecharge, "san_userviprecharge", self.player.ID)

	if self.Sql_VipRecharge.Uid <= 0 {
		self.Sql_VipRecharge.Uid = self.player.ID
		self.Sql_VipRecharge.taskinfo = make(map[int]*JS_VipRechargeInfo)
		self.Encode()
		InsertTable("san_userviprecharge", &self.Sql_VipRecharge, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_VipRecharge.Init("san_userviprecharge", &self.Sql_VipRecharge, true)
}

//! 将数据库数据写入data
func (self *ModVipRecharge) Decode() {
	json.Unmarshal([]byte(self.Sql_VipRecharge.Taskinfo), &self.Sql_VipRecharge.taskinfo)
}

//! 将data数据写入数据库
func (self *ModVipRecharge) Encode() {
	self.Sql_VipRecharge.Taskinfo = HF_JtoA(&self.Sql_VipRecharge.taskinfo)
}

func (self *ModVipRecharge) OnGetOtherData() {

}

// 注册消息
func (self *ModVipRecharge) onReg(handlers map[string]func(body []byte)) {
	handlers["getviprecharge"] = self.GetVipRecharge
}

func (self *ModVipRecharge) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModVipRecharge) OnSave(sql bool) {
	self.Encode()
	self.Sql_VipRecharge.Update(sql)
}

func (self *ModVipRecharge) SendInfo() {
	self.checkTask()
	var msg S2C_VipRechargeInfo
	msg.Cid = "viprecharge"
	msg.TaskInfos = self.Sql_VipRecharge.taskinfo
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModVipRecharge) GetVipRecharge(body []byte) {

	var msg C2S_GetVipRecharge
	json.Unmarshal(body, &msg)

	for _, v := range self.Sql_VipRecharge.taskinfo {
		config, ok := GetCsvMgr().VipConfigMap[v.Taskid]
		if !ok {
			continue
		}
		if v.Taskid == msg.VipLevel && v.PickUp < v.Finish {
			v.PickUp++
			//发放奖励
			item := self.player.AddObjectLst(config.DailyItem, config.DailyNum, "VIP每日礼包", v.Taskid, 0, 0)
			var sendmsg S2C_VipRechargeGift
			sendmsg.Cid = "viprechargegift"
			sendmsg.TaskInfo = v
			sendmsg.GetItems = item
			self.player.SendMsg(sendmsg.Cid, HF_JtoB(&sendmsg))
			return
		}
	}
}

func (self *ModVipRecharge) checkTask() {
	for _, v := range GetCsvMgr().VipConfigMap {
		_, ok := self.Sql_VipRecharge.taskinfo[v.Viplevel]
		if ok {
			continue
		}
		self.Sql_VipRecharge.taskinfo[v.Viplevel] = self.NewTaskInfo(v.Viplevel, v.TaskTypes)
	}

	if self.Sql_VipRecharge.NextTime <= TimeServer().Unix() {
		for _, v := range self.Sql_VipRecharge.taskinfo {
			v.Finish = 0
		}
		self.Sql_VipRecharge.NextTime = HF_GetNextDayStart()
	}

	for _, v := range self.Sql_VipRecharge.taskinfo {
		config, ok := GetCsvMgr().VipConfigMap[v.Taskid]
		if !ok {
			continue
		}
		if v.Tasktypes == 0 {
			v.Finish = config.DailyTimes
		}
	}
}

func (self *ModVipRecharge) HandleTask(taskType int, n1 int, n2 int, n3 int) {

	if len(self.Sql_VipRecharge.taskinfo) == 0 {
		self.checkTask()
		if len(self.Sql_VipRecharge.taskinfo) == 0 {
			return
		}
	}

	for _, pTask := range self.Sql_VipRecharge.taskinfo {
		if pTask == nil {
			continue
		}

		if pTask.Tasktypes != taskType {
			continue
		}

		config, ok := GetCsvMgr().VipConfigMap[pTask.Taskid]
		if !ok {
			LogError("error:vipconfig")
			return
		}

		if config.DailyTimes <= pTask.Finish {
			continue
		}

		plan, add := DoTask(&TaskNode{Id: config.Viplevel, Tasktypes: config.TaskTypes, N1: config.Ns[0], N2: config.Ns[1], N3: config.Ns[2], N4: config.Ns[3]}, self.player, n1, n2, n3)

		if plan == 0 {
			continue
		}

		if add {
			pTask.Plan += plan
		} else {
			if plan > pTask.Plan {
				pTask.Plan = plan
			}
		}

		if pTask.Plan >= config.Ns[0] {
			pTask.Plan = 0
			pTask.Finish++
			return
		}
	}
}

func (self *ModVipRecharge) NewTaskInfo(taskid int, tasktypes int) *JS_VipRechargeInfo {
	taskinfo := new(JS_VipRechargeInfo)
	taskinfo.Taskid = taskid
	taskinfo.Tasktypes = tasktypes
	taskinfo.Plan = 0
	taskinfo.Finish = 0

	return taskinfo
}

func (self *ModVipRecharge) OnRefresh() {
	for _, pTask := range self.Sql_VipRecharge.taskinfo {
		pTask.Finish = 0
		pTask.PickUp = 0
	}
}

func (self *ModVipRecharge) HandleRecharge(grade int) {
	for _, v := range self.Sql_VipRecharge.taskinfo {
		config, ok := GetCsvMgr().VipConfigMap[v.Taskid]
		if !ok {
			continue
		}
		if config.TaskTypes != 0 && v.PickUp < v.Finish {
			v.PickUp++
			//发放奖励
			item := self.player.AddObjectLst(config.DailyItem, config.DailyNum, "VIP每日礼包", v.Taskid, 0, 0)
			var sendmsg S2C_VipRechargeGift
			sendmsg.Cid = "viprechargegift"
			sendmsg.TaskInfo = v
			sendmsg.GetItems = item
			self.player.SendMsg(sendmsg.Cid, HF_JtoB(&sendmsg))

			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BUY_VIP_BOX, grade, config.Viplevel, 0, "购买VIP礼包", 0, 0, self.player)
			return
		}
	}
}
