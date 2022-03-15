//目标系统玩法（爵位） 20190604 by zy
package game

import (
	"encoding/json"
	"fmt"
)

// 进度
type JS_NobilityTaskInfo struct {
	Taskid    int `json:"taskid"`    // 任务Id
	Tasktypes int `json:"tasktypes"` // 任务类型
	Plan      int `json:"plan"`      // 进度
	State     int `json:"finish"`    // 状态
}

//! 任务数据库
type San_NobilityTask struct {
	Uid       int64
	Taskinfo  string
	Level     int    //! 爵位等级
	GetReward string //! 爵位等级

	taskinfo  map[int]*JS_NobilityTaskInfo //! 任务信息
	getReward map[int]int                  //! 爵位奖励的领取标记
	DataUpdate
}

//! 任务
type ModNobilityTask struct {
	player           *Player
	Sql_NobilityTask San_NobilityTask
	init             bool
	chg              []JS_NobilityTaskInfo
}

func (self *ModNobilityTask) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_usernobilitytask` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_NobilityTask, "san_usernobilitytask", self.player.ID)

	if self.Sql_NobilityTask.Uid <= 0 {
		self.Sql_NobilityTask.Uid = self.player.ID
		self.Sql_NobilityTask.taskinfo = make(map[int]*JS_NobilityTaskInfo)
		self.Sql_NobilityTask.getReward = make(map[int]int)
		self.Encode()
		InsertTable("san_usernobilitytask", &self.Sql_NobilityTask, 0, true)
	} else {
		self.Decode()
	}

	if self.Sql_NobilityTask.getReward == nil {
		self.Sql_NobilityTask.getReward = make(map[int]int)
	}

	self.Sql_NobilityTask.Init("san_usernobilitytask", &self.Sql_NobilityTask, true)
}

//! 将数据库数据写入data
func (self *ModNobilityTask) Decode() {
	json.Unmarshal([]byte(self.Sql_NobilityTask.Taskinfo), &self.Sql_NobilityTask.taskinfo)
	json.Unmarshal([]byte(self.Sql_NobilityTask.GetReward), &self.Sql_NobilityTask.getReward)
}

//! 将data数据写入数据库
func (self *ModNobilityTask) Encode() {
	self.Sql_NobilityTask.Taskinfo = HF_JtoA(&self.Sql_NobilityTask.taskinfo)
	self.Sql_NobilityTask.GetReward = HF_JtoA(&self.Sql_NobilityTask.getReward)
}

func (self *ModNobilityTask) OnGetOtherData() {

}

// 注册消息
func (self *ModNobilityTask) onReg(handlers map[string]func(body []byte)) {
	handlers["takenobilitytask"] = self.TakeNobilityTask
	handlers["levelupnobility"] = self.LevelUpNobility
	handlers["getnobilityreward"] = self.GetNobilityReward
}

func (self *ModNobilityTask) TakeNobilityTask(body []byte) {
	var msg C2STakeNobilityAward
	json.Unmarshal(body, &msg)
	node, nodeOK := self.Sql_NobilityTask.taskinfo[msg.TaskId]
	if !nodeOK || node.State != CANTAKE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NOBILITYTASK_TAKE_ERROR"))
		return
	}
	config, ok := GetCsvMgr().NobilityConfigMap[msg.TaskId]
	if !ok {
		LogError("爵位配置错误")
		return
	}

	item := self.player.AddObjectSimple(config.Items, config.Nums, "爵位奖励领取", msg.TaskId, self.Sql_NobilityTask.Level, 0)
	node.State = TAKEN

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NOBILITY_AWARD, msg.TaskId, 0, 0, "领取爵位任务奖励", 0, 0, self.player)

	var sendmsg S2CTakeNobilityTask
	sendmsg.Cid = "takenobilitytask"
	sendmsg.TaskInfo = self.Sql_NobilityTask.taskinfo
	sendmsg.Level = self.Sql_NobilityTask.Level
	sendmsg.Award = item
	self.player.SendMsg(sendmsg.Cid, HF_JtoB(&sendmsg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NOBILITY_TASK, node.Taskid, 0, 0, "领取爵位奖励", 0, 0, self.player)

}

func (self *ModNobilityTask) LevelUpNobility(body []byte) {

	var msg C2SLevelUpNobility
	json.Unmarshal(body, &msg)

	if self.Sql_NobilityTask.Level+1 != msg.TaskId {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TECH_INSUFFICIENT_UPGRADING_CONDITIONS"))
		return
	}

	group := 0
	config, ok := GetCsvMgr().NobilityRewardMap[self.Sql_NobilityTask.Level]
	if ok {
		group = config.Task
	}

	//看能不能升级
	isUp := true
	for _, v := range self.Sql_NobilityTask.taskinfo {
		taskConfig := GetCsvMgr().NobilityConfigMap[v.Taskid]
		if taskConfig == nil {
			continue
		}
		if taskConfig.TaskGroup != group {
			continue
		}
		if v.State != TAKEN {
			isUp = false
			break
		}
	}

	if !isUp {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TECH_INSUFFICIENT_UPGRADING_CONDITIONS"))
		return
	}

	//看等级是否达到最大
	_, okNext := GetCsvMgr().NobilityRewardMap[self.Sql_NobilityTask.Level+1]
	if !okNext {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NOBILITYTASK_CONFIG_ERROR"))
		return
	}

	self.Sql_NobilityTask.Level++
	self.player.HandleTask(TASK_TYPE_NOBILITY_LEVEL, self.Sql_NobilityTask.Level, 0, 0)
	//发送奖励
	items := make([]PassItem, 0)
	/*
		item := self.player.AddObjectLst(configNext.Items, configNext.Nums, "爵位晋升", self.Sql_NobilityTask.Level, 0, 0)
		items = append(items, item...)
		if self.player.Sql_UserBase.Vip >= configNext.Vip {
			itemVip := self.player.AddObjectLst(configNext.VipItems, configNext.VipNums, "爵位晋升", self.Sql_NobilityTask.Level, 0, 0)
			items = append(items, itemVip...)
		}
	*/

	var sendmsg S2CLevelUpNobility
	sendmsg.Cid = "levelupnobility"
	sendmsg.Level = self.Sql_NobilityTask.Level
	sendmsg.Award = items
	self.player.SendMsg(sendmsg.Cid, HF_JtoB(&sendmsg))

	//继续领奖励
	var msgNext C2S_GetNobilityReward
	msgNext.TaskId = self.Sql_NobilityTask.Level
	msgNext.Belog = LOGIC_TRUE
	self.GetNobilityReward(HF_JtoB(msgNext))
}

func (self *ModNobilityTask) GetNobilityReward(body []byte) {

	var msg C2S_GetNobilityReward
	json.Unmarshal(body, &msg)

	if msg.TaskId > self.Sql_NobilityTask.Level {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NOBILITYTASK_CONFIG_ERROR"))
		return
	}

	config, ok := GetCsvMgr().NobilityRewardMap[msg.TaskId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NOBILITYTASK_CONFIG_ERROR"))
		return
	}

	items := make([]PassItem, 0)
	nowSign := 0
	_, okSign := self.Sql_NobilityTask.getReward[msg.TaskId]
	if okSign {
		nowSign = self.Sql_NobilityTask.getReward[msg.TaskId]
	}
	//先计算普通奖励
	if nowSign < 1 {
		for i := 0; i < len(config.Items); i++ {
			if config.Items[i] > 0 {
				items = append(items, PassItem{
					ItemID: config.Items[i],
					Num:    config.Nums[i],
				})
			}
		}
		self.Sql_NobilityTask.getReward[msg.TaskId] = 1
	}

	if nowSign < 2 {
		if self.player.Sql_UserBase.Vip >= config.Vip {
			for i := 0; i < len(config.VipItems); i++ {
				if config.VipItems[i] > 0 {
					items = append(items, PassItem{
						ItemID: config.VipItems[i],
						Num:    config.VipNums[i],
					})
				}
			}
			self.Sql_NobilityTask.getReward[msg.TaskId] = 2
		}
	}

	if len(items) <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_CONSUMERTOP_NOT_AWARD"))
		return
	}

	logId := 0
	logDec := ""
	if msg.Belog == LOGIC_TRUE {
		logId = LOG_NOBILITY_UP
		logDec = "爵位晋升"
	} else {
		logId = LOG_NOBILITY_UP_EXT
		logDec = "领取爵位额外奖励"
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, logId, self.Sql_NobilityTask.Level, 0, 0, logDec, 0, 0, self.player)

	var sendmsg S2C_GetNobilityReward
	sendmsg.Cid = "getnobilityreward"
	sendmsg.Level = self.Sql_NobilityTask.Level
	sendmsg.Award = self.player.AddObjectPassItem(items, logDec, self.Sql_NobilityTask.getReward[msg.TaskId], 0, 0)
	sendmsg.GetReward = self.Sql_NobilityTask.getReward
	self.player.SendMsg(sendmsg.Cid, HF_JtoB(&sendmsg))
}

func (self *ModNobilityTask) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModNobilityTask) OnSave(sql bool) {
	self.Encode()
	self.Sql_NobilityTask.Update(sql)
}

func (self *ModNobilityTask) SendInfo() {
	self.checkTask()
	var msg S2C_NobilityTask
	msg.Cid = "nobilitytask"
	msg.TaskInfo = self.Sql_NobilityTask.taskinfo
	msg.Level = self.Sql_NobilityTask.Level
	msg.GetReward = self.Sql_NobilityTask.getReward
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModNobilityTask) checkTask() {

	for _, v := range GetCsvMgr().NobilityConfigMap {
		_, ok := self.Sql_NobilityTask.taskinfo[v.Id]
		if ok {
			continue
		}
		self.Sql_NobilityTask.taskinfo[v.Id] = self.NewTaskInfo(v.Id, v.TaskType)
	}
}

func (self *ModNobilityTask) HandleTask(taskType int, n1 int, n2 int, n3 int) {

	if len(self.Sql_NobilityTask.taskinfo) == 0 {
		self.checkTask()
	}

	for _, pTask := range self.Sql_NobilityTask.taskinfo {
		if pTask == nil || pTask.State >= CANTAKE {
			continue
		}

		if pTask.Tasktypes != taskType {
			continue
		}

		config, ok := GetCsvMgr().NobilityConfigMap[pTask.Taskid]
		if !ok {
			LogError("爵位配置错误")
			return
		}

		plan, add := DoTask(&TaskNode{Id: config.Id, Tasktypes: config.TaskType, N1: config.Ns[0], N2: config.Ns[1], N3: config.Ns[2], N4: config.Ns[3]}, self.player, n1, n2, n3)

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
			pTask.State = CANTAKE
		}
		self.chg = append(self.chg, *pTask)
	}
}

func (self *ModNobilityTask) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_UpdateNobilityTask
	msg.Cid = "updatenobilitytask"
	msg.TaskInfo = self.chg
	self.chg = make([]JS_NobilityTaskInfo, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("updatenobilitytask", smsg)

}

func (self *ModNobilityTask) NewTaskInfo(taskid int, tasktypes int) *JS_NobilityTaskInfo {
	taskinfo := new(JS_NobilityTaskInfo)
	taskinfo.Taskid = taskid
	taskinfo.Tasktypes = tasktypes
	taskinfo.Plan = 0
	taskinfo.State = CANTFINISH

	return taskinfo
}

func (self *ModNobilityTask) GetNobilityPrivilege(nobilityType int) bool {
	//return false //爵位干掉了
	//先找对应特权的爵位等级
	for _, v := range GetCsvMgr().NobilityRewardMap {
		if v.Privilege == nobilityType {
			return self.Sql_NobilityTask.Level >= v.Id
		}
	}
	return true
}

func (self *ModNobilityTask) GmLevelUp(level int) {
	for _, v := range self.Sql_NobilityTask.taskinfo {
		taskConfig := GetCsvMgr().NobilityConfigMap[v.Taskid]
		if taskConfig == nil {
			continue
		}
		if v.State != CANTFINISH {
			continue
		}
		if taskConfig.TaskGroup%100 > level {
			continue
		}
		v.State = CANTAKE
		self.chg = append(self.chg, *v)
	}
}
