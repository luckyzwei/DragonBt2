package game

import (
	"encoding/json"
	"fmt"
)

// 进度
type JS_TargetTaskInfo struct {
	Taskid    int `json:"taskid"`    // 任务Id
	Tasktypes int `json:"tasktypes"` // 任务类型
	Plan      int `json:"plan"`      // 进度
	Finish    int `json:"finish"`    // 是否完成
	Pickup    int `json:"pickup"`    // 是否领取奖励
	BuyGet    int `json:"buyget"`    // 付费领取状态 1表示领完了1级徽章
}

//! 任务数据库
type San_TargetTask struct {
	Uid           int64
	Taskinfo      string //! 任务组标记
	SystemInfo    string //! 任务信息
	NobilityLevel int    //!
	BuyLevel      string //! 徽章等级

	taskinfo   []JS_TargetTaskInfo //! 任务信息
	systemInfo map[int]int         //! 任务组标记
	buyLevel   map[int]int         //徽章等级

	DataUpdate
}

//! 任务
type ModTargetTask struct {
	player   *Player
	Sql_Task San_TargetTask
	chg      []JS_TargetTaskInfo
	init     bool
}

func (self *ModTargetTask) OnGetData(player *Player) {
	self.player = player
	self.chg = make([]JS_TargetTaskInfo, 0)
}

func (self *ModTargetTask) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_usertargettask` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Task, "san_usertargettask", self.player.ID)

	if self.Sql_Task.Uid <= 0 {
		self.Sql_Task.Uid = self.player.ID
		self.Sql_Task.taskinfo = make([]JS_TargetTaskInfo, 0)
		self.Sql_Task.systemInfo = make(map[int]int, 0)
		self.Sql_Task.buyLevel = make(map[int]int, 0)
		self.Encode()
		InsertTable("san_usertargettask", &self.Sql_Task, 0, true)
	} else {
		self.Decode()
	}

	if self.Sql_Task.buyLevel == nil {
		self.Sql_Task.buyLevel = make(map[int]int, 0)
	}

	self.Sql_Task.Init("san_usertargettask", &self.Sql_Task, true)
}

func (self *ModTargetTask) OnMsg(ctrl string, body []byte) bool {
	return false
}

/*
LOG_TARGETTASK_GET       = 6301 //!领取冒险任务奖励
	LOG_TARGETTASK_BADGE_GET = 6302 //!领取冒险任务积累奖励
	LOG_TARGETTASK_BADEG_BUY = 6303 //!激活冒险徽章
*/

func (self *ModTargetTask) onReg(handlers map[string]func(body []byte)) {
	handlers["finishtargettask"] = self.FinishTargetTask
	handlers["gettargetlvreward"] = self.GetTargetLvReward
}

func (self *ModTargetTask) OnSave(sql bool) {
	self.Encode()
	self.Sql_Task.Update(sql)
}

//! 将数据库数据写入data
func (self *ModTargetTask) Decode() {
	json.Unmarshal([]byte(self.Sql_Task.Taskinfo), &self.Sql_Task.taskinfo)
	json.Unmarshal([]byte(self.Sql_Task.SystemInfo), &self.Sql_Task.systemInfo)
	json.Unmarshal([]byte(self.Sql_Task.BuyLevel), &self.Sql_Task.buyLevel)
}

//! 将data数据写入数据库
func (self *ModTargetTask) Encode() {
	self.Sql_Task.Taskinfo = HF_JtoA(&self.Sql_Task.taskinfo)
	self.Sql_Task.SystemInfo = HF_JtoA(&self.Sql_Task.systemInfo)
	self.Sql_Task.BuyLevel = HF_JtoA(&self.Sql_Task.buyLevel)
}

func (self *ModTargetTask) AddTask(taskid int) {
	var node JS_TargetTaskInfo
	node.Taskid = taskid
	node.Plan = 0
	node.Pickup = 0
	node.Finish = 0
	config := GetCsvMgr().TargetTaskConfig[taskid]
	if config == nil {
		LogError("TargetTaskConfig taskId:", taskid, " is nil")
	} else {
		node.Tasktypes = config.Tasktypes
	}
	self.Sql_Task.taskinfo = append(self.Sql_Task.taskinfo, node)
}

//! 得到任务
func (self *ModTargetTask) GetTask(taskid int, add bool) *JS_TargetTaskInfo {
	for i := 0; i < len(self.Sql_Task.taskinfo); i++ {
		if self.Sql_Task.taskinfo[i].Taskid == taskid {
			return &self.Sql_Task.taskinfo[i]
		}
	}

	if add {
		var node JS_TargetTaskInfo
		node.Taskid = taskid
		node.Plan = 0
		node.Pickup = 0
		node.Finish = 0

		config := GetCsvMgr().TargetTaskConfig[taskid]
		if config == nil {
			LogError("TargetTaskConfig taskId:", taskid, " is nil")
		} else {
			node.Tasktypes = config.Tasktypes
		}
		self.Sql_Task.taskinfo = append(self.Sql_Task.taskinfo, node)
		return &(self.Sql_Task.taskinfo[len(self.Sql_Task.taskinfo)-1])
	}

	return nil
}

//! 处理任务
func (self *ModTargetTask) HandleTask(tasktype, n2, n3, n4 int) {
	for _, value := range GetCsvMgr().TargetTaskConfig {
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
			if plan > node.Plan {
				node.Plan = plan
				chg = true
			}
		}

		if node.Plan >= value.TaskNode.N1 {
			node.Finish = LOGIC_TRUE
			chg = true
		}

		if chg {
			//只同步需要更新的
			config := GetCsvMgr().TargetTaskConfig[node.Taskid]
			if config == nil {
				continue
			}
			_, ok := self.Sql_Task.systemInfo[config.System]
			if !ok {
				continue
			}
			if config.Group != self.Sql_Task.systemInfo[config.System] {
				continue
			}
			self.chg = append(self.chg, *node)
		}
	}
}

func (self *ModTargetTask) FinishTargetTask(body []byte) {

	var msg C2S_TargetTaskFinish
	json.Unmarshal(body, &msg)

	pConfig, ok := GetCsvMgr().TargetTaskConfig[msg.Taskid]
	if !ok || len(pConfig.Item) != len(pConfig.Num) {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	if pConfig.Condition != 0 && self.player.Sql_UserBase.PassMax < pConfig.Condition {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	_, okSysyem := self.Sql_Task.systemInfo[pConfig.System]
	if !okSysyem {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	if self.Sql_Task.systemInfo[pConfig.System] != pConfig.Group {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	node := self.GetTask(msg.Taskid, false)
	if node == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}
	if node.Finish == 0 {
		self.chg = append(self.chg, *node)
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_NOT_COMPLETED"))
		return
	}
	if node.Pickup == 1 {
		self.chg = append(self.chg, *node)
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_KINGTASK_TASKS_HAVE_BEEN_RECEIVED"))
		return
	}

	if pConfig.Prepose != 0 { //! 前置任务未完成
		pretask := self.GetTask(pConfig.Prepose, false)
		if pretask == nil || pretask.Pickup == 0 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
			return
		}
	}
	taskUpdate := make([]JS_TargetTaskInfo, 0)
	node.Pickup = 1
	//增加徽章奖励
	systemLv := self.Sql_Task.buyLevel[pConfig.System]
	items := pConfig.ItemLv[node.BuyGet:systemLv]
	nums := pConfig.NumLv[node.BuyGet:systemLv]
	node.BuyGet = systemLv
	extItem := self.player.AddObjectLst(items, nums, "领取目标任务", 0, 0, 0)

	taskUpdate = append(taskUpdate, *node)
	//更新system
	oldSign := self.Sql_Task.systemInfo[pConfig.System]
	self.Sql_Task.systemInfo[pConfig.System]++
	for _, v := range self.Sql_Task.taskinfo {
		config := GetCsvMgr().TargetTaskConfig[v.Taskid]
		if config == nil {
			continue
		}
		if config.System != pConfig.System {
			continue
		}
		if v.Pickup != 1 && config.Group < self.Sql_Task.systemInfo[pConfig.System] {
			self.Sql_Task.systemInfo[config.System] = config.Group
		}
	}
	//如果进行了组切换，则需要同步更多任务
	if oldSign != self.Sql_Task.systemInfo[pConfig.System] {
		for _, v := range self.Sql_Task.taskinfo {
			config := GetCsvMgr().TargetTaskConfig[v.Taskid]
			if config == nil {
				continue
			}
			if config.System != pConfig.System {
				continue
			}
			if config.Group != self.Sql_Task.systemInfo[pConfig.System] {
				continue
			}
			taskUpdate = append(taskUpdate, v)
		}
	}
	logDec := ""
	logId := 0
	if pConfig.System == 1 {
		logId = LOG_TARGETTASK_GET
		logDec = "领取冒险任务奖励"
	} else {
		logId = LOG_TARGETTASK_TOWER_GET
		logDec = "领取试炼之塔任务奖励"
	}
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, logId, 1, systemLv, 0, logDec, 0, 0, self.player)

	getItem := self.player.AddObjectLst(pConfig.Item, pConfig.Num, logDec, 0, 0, 0)

	var msgRel S2C_TargetTaskFinish
	msgRel.Cid = "targettaskfinish"
	msgRel.Item = getItem
	msgRel.Item = append(msgRel.Item, extItem...)
	msgRel.TaskId = msg.Taskid
	msgRel.SystemInfo = self.Sql_Task.systemInfo
	msgRel.Info = taskUpdate
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

}

func (self *ModTargetTask) GetTargetLvReward(body []byte) {

	var msg C2S_GetTargetLvReward
	json.Unmarshal(body, &msg)

	_, okSysyem := self.Sql_Task.systemInfo[msg.SystemId]
	if !okSysyem {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	reward := make(map[int]*Item)
	taskUpdate := make([]JS_TargetTaskInfo, 0)

	for k, v := range self.Sql_Task.taskinfo {
		config := GetCsvMgr().TargetTaskConfig[v.Taskid]
		if config == nil {
			continue
		}
		if msg.SystemId != config.System {
			continue
		}
		_, ok := self.Sql_Task.systemInfo[config.System]
		if !ok {
			continue
		}
		if v.Pickup != LOGIC_TRUE {
			continue
		}
		systemLv := self.Sql_Task.buyLevel[config.System]
		if v.BuyGet >= systemLv {
			continue
		}

		if config.Group > self.Sql_Task.systemInfo[config.System] {
			continue
		}

		if config.Group <= self.Sql_Task.systemInfo[config.System] {
			items := config.ItemLv[v.BuyGet:systemLv]
			nums := config.NumLv[v.BuyGet:systemLv]
			AddItemMapHelper(reward, items, nums)
			self.Sql_Task.taskinfo[k].BuyGet = systemLv
		}
		if config.Group != self.Sql_Task.systemInfo[config.System] {
			continue
		}
		taskUpdate = append(taskUpdate, v)
	}

	if len(reward) == 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_NOT_COMPLETED"))
		return
	}
	logDec := ""
	logId := 0
	if msg.SystemId == 1 {
		logId = LOG_TARGETTASK_BADGE_GET
		logDec = "领取冒险任务积累奖励"
	} else {
		logId = LOG_TARGETTASK_BADGE_TOWER_GET
		logDec = "领取试炼之塔积累奖励"
	}
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, logId, 1, 0, 0, logDec, 0, 0, self.player)

	getItem := self.player.AddObjectItemMap(reward, logDec, msg.SystemId, 0, 0)

	var msgRel S2C_GetTargetLvReward
	msgRel.Cid = "gettargetlvreward"
	msgRel.Item = getItem
	msgRel.SystemId = msg.SystemId
	msgRel.Info = taskUpdate
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

///////////////////////////
func (self *ModTargetTask) SendInfo() {
	self.checkTask()

	//计算出累积物品
	buyRewards := make(map[int]map[int]map[int]*Item)

	var msg S2C_TargetTaskInfo
	msg.Cid = "targettaskinfo"
	msg.SystemInfo = self.Sql_Task.systemInfo
	for _, v := range self.Sql_Task.taskinfo {
		config := GetCsvMgr().TargetTaskConfig[v.Taskid]
		if config == nil {
			continue
		}
		_, ok := self.Sql_Task.systemInfo[config.System]
		if !ok {
			continue
		}
		if config.Group <= self.Sql_Task.systemInfo[config.System] {
			//累计下以前没购买徽章的奖励
			for _, con := range GetCsvMgr().BadgeTaskConfig {
				if v.Finish != LOGIC_TRUE {
					continue
				}
				if v.Pickup != LOGIC_TRUE {
					continue
				}
				if con.System != config.System {
					continue
				}
				if v.BuyGet >= con.Lv {
					continue
				}
				_, okSysten := buyRewards[con.System]
				if !okSysten {
					buyRewards[con.System] = make(map[int]map[int]*Item)
				}

				_, okLv := buyRewards[con.System][con.Lv]
				if !okLv {
					buyRewards[con.System][con.Lv] = make(map[int]*Item)
				}
				items := config.ItemLv[v.BuyGet:con.Lv]
				nums := config.NumLv[v.BuyGet:con.Lv]
				AddItemMapHelper(buyRewards[con.System][con.Lv], items, nums)
			}
		}
		if config.Group != self.Sql_Task.systemInfo[config.System] {
			continue
		}
		msg.Info = append(msg.Info, v)
	}
	msg.BuyLevel = self.Sql_Task.buyLevel
	msg.Buyrewards = buyRewards
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	self.init = true
}

func (self *ModTargetTask) checkTask() {
	if self.Sql_Task.systemInfo == nil {
		self.Sql_Task.systemInfo = make(map[int]int)
	}

	taskMap := make(map[int]JS_TargetTaskInfo)
	for _, v := range self.Sql_Task.taskinfo {
		taskMap[v.Taskid] = v
	}

	//20201021版本大改，数据重做，如果是之前的数据需要全部清除  208020用来判断版本切换,不存在则删除记录
	_, okExist := taskMap[208020]
	if !okExist {
		taskMap = make(map[int]JS_TargetTaskInfo)
		self.Sql_Task.taskinfo = make([]JS_TargetTaskInfo, 0)
		self.Sql_Task.systemInfo = make(map[int]int, 0)
	}

	for _, v := range GetCsvMgr().TargetTaskConfig {
		_, ok := taskMap[v.Taskid]
		if ok {
			continue
		}
		self.AddTask(v.Taskid)
	}

	for _, v := range self.Sql_Task.taskinfo {
		config := GetCsvMgr().TargetTaskConfig[v.Taskid]
		if config == nil {
			continue
		}
		group, ok := self.Sql_Task.systemInfo[config.System]
		if ok {
			if v.Pickup != 1 && config.Group < group {
				self.Sql_Task.systemInfo[config.System] = config.Group
			}
		} else {
			self.Sql_Task.systemInfo[config.System] = config.Group
		}
	}
}

func (self *ModTargetTask) SendUpdate() {
	if !self.init {
		return
	}

	if len(self.chg) == 0 {
		return
	}

	var msg S2C_TargetTaskUpdate
	msg.Cid = "updatetargettask"
	msg.Info = self.chg
	self.chg = make([]JS_TargetTaskInfo, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}

func (self *ModTargetTask) IsCanBuy(money int) bool {
	config := GetCsvMgr().GetBadgeTaskConfig(money)
	if config == nil {
		return false
	}

	if config.Need != self.Sql_Task.buyLevel[config.System] {
		return false
	}

	return true
}

func (self *ModTargetTask) BuyLevel(money int) {
	config := GetCsvMgr().GetBadgeTaskConfig(money)
	if config == nil {
		return
	}
	self.Sql_Task.buyLevel[config.System] = config.Lv

	self.SendInfo()
}
