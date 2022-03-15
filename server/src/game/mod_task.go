package game

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

const (
	FINISH_TASK     = "finishtask"
	FINISH_TASK_RET = "finishtask2ret"
	LIVENESS_TASK   = "livenesstask"
	TASK_UPDATE     = "task2_update"
	TASK_INFO       = "task2info"
	FINISH_ALL_TASK = "finishalltask"
)

const (
	TASK_TYPE_MIN   = 0
	TASK_TYPE_DAILY = 1 // 每日任务
	TASK_TYPE_WEEK  = 2 // 周常任务
	TASK_TYPE_MAX   = 3
)

type JS_LivenessInfo struct {
	Type     int `json:"type"`     // 类型
	Liveness int `json:"liveness"` //! 活跃度值
	Award    int `json:"award"`    //! 活跃度奖励状态
}

//! 任务数据库
type San_Task struct {
	Uid             int64
	Taskinfo        string             //! 任务信息
	Liveness        string             //! 活跃度信息
	Award           int                //! 活跃度奖励状态 废弃
	WeekRefreshTime int64              //! 周常刷新时间
	taskinfo        []JS_TaskInfo      //! 任务信息
	liveness        []*JS_LivenessInfo //! 活跃度信息

	DataUpdate
}

//! 任务
type ModTask struct {
	player   *Player
	Sql_Task San_Task
	chg      []JS_TaskInfo
	init     bool
}

func (self *ModTask) OnGetData(player *Player) {
	self.player = player
	self.chg = make([]JS_TaskInfo, 0)
}

func (self *ModTask) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_usertask` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Task, "san_usertask", self.player.ID)

	if self.Sql_Task.Uid <= 0 {
		self.Sql_Task.Uid = self.player.ID
		self.Sql_Task.taskinfo = make([]JS_TaskInfo, 0)
		self.Encode()
		InsertTable("san_usertask", &self.Sql_Task, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_Task.Init("san_usertask", &self.Sql_Task, true)

	if self.Sql_Task.WeekRefreshTime == 0 {
		self.Sql_Task.WeekRefreshTime = HF_GetWeekTime()
	}

	//修复竞技场52任务  反复提示可领取
	self.chg = make([]JS_TaskInfo, 0)
}

func (self *ModTask) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case FINISH_TASK:
		var c2s_msg C2S_TaskFinish
		json.Unmarshal(body, &c2s_msg)
		var msg S2C_TaskFinish
		msg.Cid = FINISH_TASK_RET
		msg.Ret, msg.Item, msg.Info = self.Finish(HF_Atoi(c2s_msg.Taskid), c2s_msg.Tasktype)
		self.SendUpdate()
		msg.Tasktype = c2s_msg.Tasktype
		msg.Liveness = self.Sql_Task.liveness
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
		return true
	case FINISH_ALL_TASK:
		self.FinishAllTask()
		return true
	case LIVENESS_TASK:
		var msg C2S_TaskLiveness
		json.Unmarshal(body, &msg)
		self.Liveness(msg.Id)
		return true
	}

	return false
}

//! 领取所有奖励
func (self *ModTask) FinishAllTask() {
	self.player.CheckRefresh()
	outitem := make([]PassItem, 0)
	itemMap := make(map[int]*PassItem)
	// 先赛选出所有的已经完成的日常任务
	var taskInfo []*JS_TaskInfo
	for index := range self.Sql_Task.taskinfo {
		pTask := &self.Sql_Task.taskinfo[index]
		if pTask.Finish == 1 && pTask.Pickup == 0 {
			config, ok := GetCsvMgr().DailytaskConfig[pTask.Taskid]
			if !ok {
				continue
			}

			taskType := config.Tasktypes
			if taskType == 56 {
				continue
			}

			if config.Pretask != 0 { //! 前置任务未完成
				pretask := self.GetTask(config.Pretask, false)
				if pretask == nil || pretask.Finish == 0 {
					continue
				}
			}

			if self.player.Sql_UserBase.Vip < config.Vip {
				continue
			}

			if self.player.Sql_UserBase.Level < config.Level {
				continue
			}

			for i := 0; i < 3; i++ {
				itemid := config.Rewards[i]
				if itemid == 0 {
					continue
				}
				num := config.Numbers[i]
				itemid, num = self.player.AddObject(itemid, num, 26, 0, 0, "任务领取")
				pItem, ok := itemMap[itemid]
				if !ok {
					itemMap[itemid] = &PassItem{itemid, num}
				} else {
					pItem.Num += num
				}
			}

			//! 完成日常任务-1日常，2成长，3赏金
			//self.player.HandleTask(FinishTask, 1, 0, 0)
			pTask.Pickup = 1
			taskInfo = append(taskInfo, pTask)
			//GetServer().sendLog_TaskFinish(self.player, pTask.Taskid, taskType, config["taskname"], HF_Atoi(config["number3"]))
		}
	}

	if len(taskInfo) <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TASK_NO_REWARD_IS_PAYABLE"))
		return
	}

	for key := range itemMap {
		outitem = append(outitem, *itemMap[key])
	}

	self.player.SendMsg(FINISH_ALL_TASK, HF_JtoB(&S2C_FinishAllTask{
		Cid:      FINISH_ALL_TASK,
		TaskInfo: taskInfo,
		OutItem:  outitem,
	}))

}

func (self *ModTask) OnSave(sql bool) {
	self.Encode()
	self.Sql_Task.Update(sql)
}

// 日常任务、奖励、活跃度刷新
func (self *ModTask) OnRefresh() {
	weekNextTime := HF_GetWeekTime()
	for i := 0; i < len(self.Sql_Task.taskinfo); {
		if self.Sql_Task.taskinfo[i].Taskid >= 200000 && self.Sql_Task.taskinfo[i].Taskid < 300000 {
			copy(self.Sql_Task.taskinfo[i:], self.Sql_Task.taskinfo[i+1:])
			self.Sql_Task.taskinfo = self.Sql_Task.taskinfo[:len(self.Sql_Task.taskinfo)-1]
		} else {
			config, ok := GetCsvMgr().DailytaskConfig[self.Sql_Task.taskinfo[i].Taskid]
			// 取日常任务表中 判断是否是周任务 且时间戳不一致
			if ok && config != nil && config.MainType == TASK_TYPE_WEEK && weekNextTime != self.Sql_Task.WeekRefreshTime {
				copy(self.Sql_Task.taskinfo[i:], self.Sql_Task.taskinfo[i+1:])
				self.Sql_Task.taskinfo = self.Sql_Task.taskinfo[:len(self.Sql_Task.taskinfo)-1]
			} else {
				i++
			}
		}
	}

	if weekNextTime != self.Sql_Task.WeekRefreshTime {
		self.Sql_Task.WeekRefreshTime = weekNextTime
		week := self.GetLivenessinfo(TASK_TYPE_WEEK)
		if nil != week {
			week.Award = 0
			week.Liveness = 0
		}
	}

	daily := self.GetLivenessinfo(TASK_TYPE_DAILY)
	if nil != daily {
		daily.Award = 0
		daily.Liveness = 0
	}
	self.HandleTask(TASK_TYPE_PLAYER_LEVEL, 0, 0, 0)
}

//! 将数据库数据写入data
func (self *ModTask) Decode() {
	json.Unmarshal([]byte(self.Sql_Task.Taskinfo), &self.Sql_Task.taskinfo)
	json.Unmarshal([]byte(self.Sql_Task.Liveness), &self.Sql_Task.liveness)
}

//! 将data数据写入数据库
func (self *ModTask) Encode() {
	s, _ := json.Marshal(&self.Sql_Task.taskinfo)
	self.Sql_Task.Taskinfo = string(s)

	s2, _ := json.Marshal(&self.Sql_Task.liveness)
	self.Sql_Task.Liveness = string(s2)
}

func (self *ModTask) AddTask(taskid int) {
	var node JS_TaskInfo
	node.Taskid = taskid
	node.Plan = 0
	node.Pickup = 0
	node.Finish = 0
	if taskid >= 200000 {
		config := GetCsvMgr().DailytaskConfig[taskid]
		if config == nil {
			LogError("DailytaskConfig taskId:", taskid, " is nil")
		} else {
			node.Tasktypes = config.Tasktypes
		}
	} else {
		config := GetCsvMgr().GrowthtaskConfig[taskid]
		if config == nil {
			LogError("GrowthtaskConfig taskId:", taskid, " is nil")
		} else {
			node.Tasktypes = config.Tasktypes
		}
	}
	self.Sql_Task.taskinfo = append(self.Sql_Task.taskinfo, node)
}

//! 得到任务
func (self *ModTask) GetTask(taskid int, add bool) *JS_TaskInfo {
	for i := 0; i < len(self.Sql_Task.taskinfo); i++ {
		if self.Sql_Task.taskinfo[i].Taskid == taskid {
			return &self.Sql_Task.taskinfo[i]
		}
	}

	if add {
		var node JS_TaskInfo
		node.Taskid = taskid
		node.Plan = 0
		node.Pickup = 0
		node.Finish = 0
		if taskid >= 200000 {
			config := GetCsvMgr().DailytaskConfig[taskid]
			if config == nil {
				LogError("DailytaskConfig taskId:", taskid, " is nil")
			} else {
				node.Tasktypes = config.Tasktypes
			}

		} else {
			config := GetCsvMgr().GrowthtaskConfig[taskid]
			if config == nil {
				LogError("GrowthtaskConfig taskId:", taskid, " is nil")
			} else {
				node.Tasktypes = config.Tasktypes
			}
		}
		self.Sql_Task.taskinfo = append(self.Sql_Task.taskinfo, node)
		return &(self.Sql_Task.taskinfo[len(self.Sql_Task.taskinfo)-1])
	}

	return nil
}

//! 处理任务
func (self *ModTask) HandleTask(tasktype, n2, n3, n4 int) {
	for _, value := range GetCsvMgr().DailytaskConfig {
		if value.Tasktypes != tasktype {
			continue
		}
		node := self.GetTask(value.Taskid, false)
		if node != nil && node.Finish == 1 {
			continue
		}
		if value.TaskNode == nil {
			continue
		}
		plan, add := DoTask(value.TaskNode, self.player, n2, n3, n4)
		if plan == 0 {
			continue
		}

		if node == nil {
			node = self.GetTask(value.Taskid, true)
		}

		chg := false
		if add {
			node.Plan += plan
			chg = true
		} else {
			//if tasktype == 22 || tasktype == 122 {
			//	if plan < node.Plan {
			//		node.Plan = plan
			//		chg = true
			//	}
			//} else
			if tasktype == PvpRankNow {
				if plan != 0 { // 新排名为0则直接不处理
					if node.Plan == 0 { // 进度为0则说明未初始化 直接赋值
						node.Plan = plan
						chg = true
					} else { // 进度不为0 则需要判断 获得的新名次比之前要高 则赋值
						if plan < node.Plan {
							node.Plan = plan
							chg = true
						}
					}
				}
			} else {
				if plan > node.Plan {
					node.Plan = plan
					chg = true
				}
			}
		}

		//! 任务完成
		//if tasktype == 22 || tasktype == 122 {
		//	if node.Plan <= value.TaskNode.N1 {
		//		node.Finish = 1
		//		chg = true
		//	}
		//} else
		if tasktype == PvpRankNow {
			if node.Plan != 0 && node.Plan <= value.TaskNode.N1 {
				node.Finish = 1
				chg = true
			}
		} else {
			if node.Plan >= value.TaskNode.N1 {
				node.Finish = 1
				chg = true
			}
		}

		if chg {
			self.chg = append(self.chg, *node)
		}
	}

	for _, value := range GetCsvMgr().GrowthtaskConfig {
		if value.Tasktypes != tasktype {
			continue
		}

		node := self.GetTask(value.Taskid, false)
		if node != nil && node.Finish == 1 {
			continue
		}

		if value.TaskNode == nil {
			continue
		}

		if value.Level > 0 && self.player.Sql_UserBase.Level < value.Level {
			continue
		}

		plan, add := DoTask(value.TaskNode, self.player, n2, n3, n4)
		if plan == 0 {
			continue
		}

		if node == nil {
			node = self.GetTask(value.Taskid, true)
		}

		chg := false
		if add {
			node.Plan += plan
			chg = true
		} else {
			//if tasktype == 22 || tasktype == 122 {
			//	if plan < node.Plan {
			//		node.Plan = plan
			//		chg = true
			//	}
			//} else
			if tasktype == PvpRankNow {
				if plan != 0 { // 新排名为0则直接不处理
					if node.Plan == 0 { // 进度为0则说明未初始化 直接赋值
						node.Plan = plan
						chg = true
					} else { // 进度不为0 则需要判断 获得的新名次比之前要高 则赋值
						if plan < node.Plan {
							node.Plan = plan
							chg = true
						}
					}
				}
			} else {
				if plan > node.Plan {
					node.Plan = plan
					chg = true
				}
			}
		}

		//! 任务完成
		//if tasktype == 22 || tasktype == 122 {
		//	if node.Plan <= value.TaskNode.N1 {
		//		node.Finish = 1
		//		chg = true
		//	}
		//} else
		if tasktype == PvpRankNow {
			if node.Plan != 0 && node.Plan <= value.TaskNode.N1 {
				node.Finish = 1
				chg = true
			}
		} else {
			if node.Plan >= value.TaskNode.N1 {
				node.Finish = 1
				chg = true
			}
		}

		if chg {
			self.chg = append(self.chg, *node)
		}
	}
}

//! //提交任务 $tasktype任务类型1日常 2成长  ret 1:等级不足  2:已经领取 3：条件未完成 4:未知任务 5:任务有效时间已过
func (self *ModTask) Finish(taskid, tasktype int) (int, []PassItem, []JS_TaskInfo) {
	if tasktype == 1 {
		return self.finishDailyTask(taskid, tasktype)
	} else if tasktype == 2 {
		return self.finishGrowTask(taskid, tasktype)
	} else {
		return -1, []PassItem{}, []JS_TaskInfo{}
	}
}

func (self *ModTask) finishDailyTask(taskid, tasktype int) (int, []PassItem, []JS_TaskInfo) {
	self.player.CheckRefresh()

	doubleitem, _ := GetActivityMgr().GetDoubleStatus(DOUBLE_POWER)
	LogDebug("体力双倍倍率:", doubleitem)
	info := make([]JS_TaskInfo, 0)
	outitem := make([]PassItem, 0)

	pConfig, ok := GetCsvMgr().DailytaskConfig[taskid]
	if !ok {
		self.player.SendErr(fmt.Sprintf(GetCsvMgr().GetText("STR_TASK_NO_CONFIG2"), taskid))
		return -1, outitem, info
	}

	if !ok {
		return -1, outitem, info
	}

	if self.player.Sql_UserBase.Vip < pConfig.Vip {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TASK_LACK_OF_ARISTOCRATIC_RANK"))
		return -1, outitem, info
	}

	if self.player.Sql_UserBase.Level < pConfig.Level {
		return 1, outitem, info
	}

	checkpoint := pConfig.Checkpoint
	if checkpoint > 0 {
		passnode := self.player.GetModule("pass").(*ModPass).GetPass(checkpoint)
		if passnode == nil {
			return 1, outitem, info
		}
	}

	//! 限时任务
	if pConfig.Tasktypes == 56 {
		retCode := self.checkPower(pConfig, taskid, doubleitem, &outitem, &info)
		if retCode != 0 {
			return retCode, outitem, info
		}

		return 0, outitem, info
	}

	node := self.GetTask(taskid, false)
	if node == nil {
		return 4, outitem, info
	}
	if node.Finish == 0 {
		self.chg = append(self.chg, *node)
		return 3, outitem, info
	}
	if node.Pickup == 1 {
		self.chg = append(self.chg, *node)
		return 2, outitem, info
	}
	if pConfig.Pretask != 0 { //! 前置任务未完成
		pretask := self.GetTask(pConfig.Pretask, false)
		if pretask == nil || pretask.Finish == 0 {
			return 3, outitem, info
		}
	}

	dec := "领取日常任务奖励"
	if pConfig.MainType == TASK_TYPE_WEEK {
		dec = "领取周常任务奖励"
	}

	node.Pickup = 1
	for i := 0; i < len(pConfig.Rewards); i++ {
		itemid := pConfig.Rewards[i]
		if itemid == 0 {
			continue
		}

		if len(pConfig.Rewards) != len(pConfig.Numbers) {
			continue
		}

		num := pConfig.Numbers[i]
		if num == 0 {
			continue
		}
		if itemid == 91000017 {
			//self.player.GetModule("feats").(*Mod_Feats).OverTask(num)
		} else {
			itemid, num = self.player.AddObject(itemid, num, taskid, 0, 0, dec)
		}
		outitem = append(outitem, PassItem{itemid, num})
	}
	info = append(info, *node)

	if pConfig.MainType == TASK_TYPE_WEEK {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_WEEK_TASK, 1, taskid, 0, dec, 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_DAILY_TASK, 1, taskid, 0, dec, 0, 0, self.player)
	}

	//self.player.HandleTask(FinishTask, tasktype, 0, 0)
	self.HandleTask(TASK_TYPE_PLAYER_LEVEL, 0, 0, 0)

	number := 0
	if len(pConfig.Numbers) >= 3 {
		number = pConfig.Numbers[2]
	}
	GetServer().sendLog_TaskFinish(self.player, taskid, tasktype, pConfig.TaskName, number)

	return 0, outitem, info
}

func (self *ModTask) finishGrowTask(taskid, tasktype int) (int, []PassItem, []JS_TaskInfo) {
	self.player.CheckRefresh()

	doubleitem, _ := GetActivityMgr().GetDoubleStatus(DOUBLE_POWER)
	LogDebug("体力双倍倍率:", doubleitem)
	info := make([]JS_TaskInfo, 0)
	outitem := make([]PassItem, 0)

	pConfig, ok := GetCsvMgr().GrowthtaskConfig[taskid]
	if !ok {
		self.player.SendErr(fmt.Sprintf(GetCsvMgr().GetText("STR_TASK_NO_CONFIG"), taskid, tasktype))
		return -1, outitem, info
	}

	if self.player.Sql_UserBase.Level < pConfig.Level {
		return 1, outitem, info
	}

	checkpoint := pConfig.Checkpoint
	if checkpoint > 0 {
		passnode := self.player.GetModule("pass").(*ModPass).GetPass(checkpoint)
		if passnode == nil {
			return 1, outitem, info
		}
	}

	node := self.GetTask(taskid, false)
	if node == nil {
		return 4, outitem, info
	}
	if node.Finish == 0 {
		self.chg = append(self.chg, *node)
		return 3, outitem, info
	}
	if node.Pickup == 1 {
		self.chg = append(self.chg, *node)
		return 2, outitem, info
	}

	if pConfig.Pretask != 0 { //! 前置任务未完成
		pretask := self.GetTask(pConfig.Pretask, false)
		if pretask == nil || pretask.Finish == 0 {
			return 3, outitem, info
		}
	}

	node.Pickup = 1
	for i := 0; i < len(pConfig.Rewards); i++ {
		itemid := pConfig.Rewards[i]
		if itemid == 0 {
			continue
		}

		if len(pConfig.Rewards) != len(pConfig.Numbers) {
			continue
		}

		num := pConfig.Numbers[i]
		if num == 0 {
			continue
		}

		if itemid == 91000017 {
			//self.player.GetModule("feats").(*Mod_Feats).OverTask(num)
		} else {
			itemid, num = self.player.AddObject(itemid, num, taskid, 0, 0, "领取主线任务奖励")
		}
		outitem = append(outitem, PassItem{itemid, num})
	}
	info = append(info, *node)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GROWTH_TASK, 1, taskid, 0, "领取主线任务奖励", 0, 0, self.player)
	//self.player.HandleTask(FinishTask, tasktype, 0, 0)
	self.HandleTask(TASK_TYPE_PLAYER_LEVEL, 0, 0, 0)
	number := 0
	if len(pConfig.Numbers) >= 3 {
		number = pConfig.Numbers[2]
	}
	GetServer().sendLog_TaskFinish(self.player, taskid, tasktype, pConfig.TaskName, number)
	return 0, outitem, info
}

func (self *ModTask) Liveness(id int) {
	csv, ok := GetCsvMgr().ActiveConfig[id]
	if !ok {
		return
	}

	nType := csv.Type
	nSort := csv.Sort

	liveness := self.GetLivenessinfo(nType)
	if nil == liveness {
		return
	}

	tmp := int(math.Pow(2.0, float64(nSort-1)))
	if (liveness.Award & tmp) != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TASK_CANT_RECEIVE_REPEATEDLY"))
		return
	}

	liveness.Award |= tmp

	dec := "领取日常活跃度宝箱"
	if nType == TASK_TYPE_WEEK {
		dec = "领取周常活跃度宝箱"
	}

	// 计算活跃度获取的任务
	outitem := make([]PassItem, 0)
	for i := 0; i < len(csv.ItemIds); i++ {
		if len(csv.ItemIds) != len(csv.ItemNums) {
			LogError("len(csv.ItemIds) != len(csv.ItemNums)")
			continue
		}
		n := csv.ItemIds[i]
		if n == 0 {
			continue
		}
		x := csv.ItemNums[i]
		n, x = self.player.AddObject(n, x, id, 0, 0, dec)
		outitem = append(outitem, PassItem{n, x})
	}

	if self.player.GetModule("recharge").(*ModRecharge).WarOrderIsOpen(WARORDER_1) && GetCsvMgr().IsLevelAndPassOpenNew(self.player.Sql_UserBase.Level, self.player.Sql_UserBase.PassMax, OPEN_LEVEL_WARORDER_1) {
		n, x := self.player.AddObject(csv.SpecialItem, csv.SpecialNum, id, 0, 0, dec)
		outitem = append(outitem, PassItem{n, x})
	}

	if GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LAST_AWARD) {
		n, x := self.player.AddObject(csv.RecastItem, csv.RecastNum, id, 0, 0, dec)
		outitem = append(outitem, PassItem{n, x})
	}
	//爵位额外箱子
	if csv.Rank > 0 {
		if self.player.GetModule("nobilitytask").(*ModNobilityTask).GetNobilityPrivilege(csv.Rank) {
			n, x := self.player.AddObject(csv.RankItem, csv.RankNum, id, 0, 0, dec)
			outitem = append(outitem, PassItem{n, x})
		}
	}
	//活动额外奖励
	activity := GetActivityMgr().GetActivity(ACT_ONHOOK_ACTIVITY_SPRING_FESTIVAL_LIVENESS)
	if activity != nil && activity.status.Status == ACTIVITY_STATUS_OPEN {
		item := self.player.AddObjectLst(csv.FestivalItem, csv.FestivalNum, dec, id, 0, 0)
		outitem = append(outitem, item...)
	}

	self.player.HandleTask(LivenessTask, csv.Active, 0, 0)

	var msg S2C_TaskLiveness
	msg.Type = csv.Type
	msg.Cid = LIVENESS_TASK
	msg.Item = outitem
	msg.Award = liveness.Award
	self.player.SendMsg(LIVENESS_TASK, HF_JtoB(&msg))

	if nType == TASK_TYPE_WEEK {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_WEEK_TASK_LIVENESS, 1, 0, 0, dec, 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_DAILY_TASK_LIVENESS, 1, 0, 0, dec, 0, 0, self.player)
	}

}

func (self *ModTask) checkTask() {
	taskMap := make(map[int]JS_TaskInfo)
	for _, v := range self.Sql_Task.taskinfo {
		taskMap[v.Taskid] = v
	}

	for _, v := range GetCsvMgr().DailytaskConfig {
		if v.Tasktypes == 56 {
			continue
		}

		if v.Level > self.player.Sql_UserBase.Level {
			continue
		}
		_, ok := taskMap[v.Taskid]
		if ok {
			continue
		}
		self.AddTask(v.Taskid)
	}

	for _, v := range GetCsvMgr().GrowthtaskConfig {
		if v.Tasktypes == 56 {
			continue
		}

		if v.Level > self.player.Sql_UserBase.Level {
			continue
		}
		_, ok := taskMap[v.Taskid]
		if ok {
			continue
		}
		self.AddTask(v.Taskid)
	}
}

///////////////////////////
func (self *ModTask) SendInfo() {
	self.checkTask()
	self.CheckTaskDone()
	self.CheckLivenessinfo()
	var msg S2C_TaskInfo
	msg.Cid = TASK_INFO
	var tasks []JS_TaskInfo
	for _, v := range self.Sql_Task.taskinfo {
		//if v.Tasktypes == 1 && v.Pickup == 1 {
		//	continue
		//}
		//
		//if v.Tasktypes == 42 && v.Pickup == 1 {
		//	continue
		//}

		tasks = append(tasks, v)
	}
	msg.Info = tasks
	msg.Liveness = self.Sql_Task.liveness
	self.player.SendMsg("task2info", HF_JtoB(&msg))

	self.init = true
}

func (self *ModTask) SendUpdate() {
	if !self.init {
		return
	}

	if len(self.chg) == 0 {
		return
	}

	var msg S2C_TaskUpdate
	msg.Cid = TASK_UPDATE
	msg.Info = self.chg
	self.chg = make([]JS_TaskInfo, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(TASK_UPDATE, smsg)
}

// 检查体力领取
func (self *ModTask) checkPower(config *DailytaskConfig, taskId int, doubleitem int, outitem *[]PassItem, info *[]JS_TaskInfo) (retCode int) {
	node := self.GetTask(taskId, true)
	if node.Pickup == 1 {
		return 2
	}

	if len(config.Ns) < 3 {
		LogError("len(config.Ns) < 3")
		return -1
	}

	begintime := config.Ns[1]
	endtime := config.Ns[2]

	beginhour := begintime / 100
	beginmin := begintime % 100
	endhour := endtime / 100
	endmin := endtime % 100

	now := TimeServer()
	beginTime := time.Date(now.Year(), now.Month(), now.Day(), beginhour, beginmin, 0, 0, time.Local).Unix()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), endhour, endmin, 0, 0, time.Local).Unix()

	// 时间没到就无法领取
	if now.Unix() < beginTime {
		return 5
	}

	// 体力值不足
	//if self.player.GetPower() >= POWERMAX {
	//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_POWERMAX"))
	//	return -1
	//}

	// 时间过了,就消耗元宝领取
	costConfig := GetCsvMgr().GetTariffConfig3(TariffTakePower, 0)
	if costConfig == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TASK_CONSUMPTION_CONFIGURATION_ERROR"))
		return
	}
	if now.Unix() > endTime {
		if err := self.player.HasObjectOk(costConfig.ItemIds, costConfig.ItemNums); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
		res := self.player.RemoveObjectLst(costConfig.ItemIds, costConfig.ItemNums, "补领体力", 0, 0, 0)
		*outitem = append(*outitem, res...)
	}

	nAddPower := 0
	for i := 0; i < 2; i++ {
		itemid := config.Rewards[i]
		if itemid == 0 {
			continue
		}
		num := config.Numbers[i]
		if itemid == 91000003 {
			num = num * doubleitem
			nAddPower = num
		}
		self.player.AddObject(itemid, num, 0, 0, 0, "补领体力")
		*outitem = append(*outitem, PassItem{itemid, num})
	}

	node.Finish = 1
	node.Pickup = 1
	*info = append(*info, *node)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_POWER_REPLACEMENT, nAddPower, costConfig.ItemNums[0], 0, "补领体力", 0, 0, self.player)
	return 0
}

func (self *ModTask) CheckTaskDone() {
	for i, pTask := range self.Sql_Task.taskinfo {
		config := GetCsvMgr().DailytaskConfig[pTask.Taskid]
		if config == nil {
			continue
		}

		if config.Tasktypes == CommonLevelTask {
			if pTask.Finish == 1 {
				continue
			}

			if pTask.Pickup == 1 {
				continue
			}

			if len(config.Ns) != 4 {
				continue
			}

			passId := config.Ns[1]
			pass := self.player.GetModule("pass").(*ModPass).GetPass(passId)
			if pass != nil {
				self.Sql_Task.taskinfo[i].Finish = 1
			}

		}
	}
}

// 获得数据
func (self *ModTask) GetLivenessinfo(nType int) *JS_LivenessInfo {
	for _, v := range self.Sql_Task.liveness {
		if v.Type == nType {
			return v
		}
	}
	return nil
}

// 检测数据结构
func (self *ModTask) CheckLivenessinfo() {
	for i := TASK_TYPE_MIN; i < TASK_TYPE_MAX; i++ {
		temp := self.GetLivenessinfo(i)
		if nil == temp {
			info := JS_LivenessInfo{i, 0, 0}
			self.Sql_Task.liveness = append(self.Sql_Task.liveness, &info)
		}
	}
}

func (self *ModTask) GMFinishAllTask() {
	for i, _ := range self.Sql_Task.taskinfo {
		if self.Sql_Task.taskinfo[i].Finish != 1 {
			self.Sql_Task.taskinfo[i].Finish = 1
		}
	}
}
