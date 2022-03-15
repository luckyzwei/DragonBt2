package game

import (
	"encoding/json"
	"fmt"
)

const (
	MaxMoneyTaskNum = 6
)

const (
	MONEY_TASK_COST_TYPE_REPLACE_ITEM = 0 // 替代物品立刻完成
	MONEY_TASK_COST_TYPE_GEM          = 1 // 钻石立刻完成
)

// 新版赏金任务
type ModMoneyTask struct {
	player *Player
	Data   SanMoneyTask // 数据库数据
}

type SanMoneyTask struct {
	Uid        int64
	FlushTimes int    // 刷新次数
	Star       int    // 随机到哪个组
	Info       string // 赏金任务信息
	TaskRecord string // 赏金任务信息

	info       map[int]*MoneyTask   // 当前任务
	taskRecord map[int]*MoneyRecord // 记录任务信息

	DataUpdate
}

type MoneyTask struct {
	TaskId      int         `json:"taskid"`       // 任务Id
	Tasktype    int         `json:"tasktype"`     // 任务类型
	Plan        int         `json:"plan"`         // 进度
	Finish      int         `json:"finish"`       // 是否完成
	Open        int         `json:"open"`         // 是否解锁 0未解锁  1 解锁
	Pickup      int         `json:"pickup"`       // 是否领取奖励
	Sort        int         `json:"sort"`         // 排序
	BoxState    []int       `json:"box_state"`    // 是否领取
	ItemsRecord [4]PassItem `json:"items_record"` // 记录领取的是哪些东西
}

type MoneyRecord struct {
	Tasktype int `json:"tasktype"` // 任务类型
	Plan     int `json:"plan"`     // 进度
}

// 同步赏金任务信息
type S2C_MoneyTaskInfo struct {
	Cid        string             `json:"cid"`
	FlushTimes int                `json:"flush_times"` // 刷新次数
	Star       int                `json:"star"`        // 随机到哪个星级
	Info       map[int]*MoneyTask `json:"info"`        // 任务信息
	Items      []PassItem         `json:"items"`       // 道具信息
	AwardTypes []int              `json:"award_types"` // 奖励类型
}

type S2C_SynMoneyTasks struct {
	Cid   string       `json:"cid"`
	Tasks []*MoneyTask `json:"tasks"`
}

// 同步赏金任务信息
type S2C_AwardMoneyTask struct {
	Cid   string     `json:"cid"`
	Info  *MoneyTask `json:"info"`  // 任务信息
	Items []PassItem `json:"items"` // 道具信息
}

// 赏金任务
type S2C_MongyTaskInfo struct {
	Cid  string             `json:"cid"`
	Info map[int]*MoneyTask `json:"task"`
}

// 同步赏金任务信息
type C2S_MoneyTaskAction struct {
	Cid    string `json:"cid"`
	TaskId int    `json:"task_id"`
}

func (m *ModMoneyTask) Decode() {
	err := json.Unmarshal([]byte(m.Data.Info), &m.Data.info)
	if err != nil {
		LogError(err.Error())
	}

	err = json.Unmarshal([]byte(m.Data.TaskRecord), &m.Data.taskRecord)
	if err != nil {
		LogError(err.Error())
	}
}

func (m *ModMoneyTask) Encode() {
	m.Data.Info = HF_JtoA(m.Data.info)
	m.Data.TaskRecord = HF_JtoA(m.Data.taskRecord)
}

func (m *ModMoneyTask) getTableName() string {
	return "san_moneytask"
}

func (m *ModMoneyTask) init(uid int64) {
	m.Data.Uid = uid
	m.CheckInfo()
}

func (m *ModMoneyTask) CheckInfo() {
	m.MakeTasks()
	m.CheckRecord()
}

func (m *ModMoneyTask) OnGetData(player *Player) {
	m.player = player
}

func (m *ModMoneyTask) OnGetOtherData() {
	tableName := m.getTableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", tableName, m.player.ID)
	GetServer().DBUser.GetOneData(sql, &m.Data, tableName, m.player.ID)
	if m.Data.Uid <= 0 {
		m.init(m.player.ID)
		m.CheckInfo()
		m.Encode()
		InsertTable(tableName, &m.Data, 0, true)
	} else {
		m.Decode()
		m.CheckInfo()
	}

	m.Data.Init(tableName, &m.Data, true)
}

func (m *ModMoneyTask) OnSave(sql bool) {
	m.Encode()
	m.Data.Update(sql)
}

func (m *ModMoneyTask) OnRefresh() {
	m.ClearRecord()
	m.FlushTasks()
	m.synInfo("money_task_flush", []PassItem{})
}

// 老的消息处理
func (m *ModMoneyTask) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (m *ModMoneyTask) onReg(handlers map[string]func(body []byte)) {
	handlers["money_task_flush"] = m.MoneyTaskFlush
	handlers["award_money_task"] = m.AwardMoneyTask
	handlers["done_money_task"] = m.DoneMoneyTask
	handlers["award_all_money_task"] = m.AwardAllMoneyTask
	handlers["money_task_give_items"] = m.GiveItems
}

// 提交材料任务
func (m *ModMoneyTask) GiveItems(body []byte) {
	var msg C2S_MoneyTaskAction
	err := json.Unmarshal(body, &msg)
	if err != nil {
		m.player.SendErr(err.Error())
		return
	}

	task, ok := m.Data.info[msg.TaskId]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_DOES_NOT_EXIST"))
		return
	}

	if task.Finish == 1 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_COMPLETED"))
		return
	}

	configs, ok := GetCsvMgr().MoneyTaskMap[task.TaskId]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MONEY_TASK_CONFIG_ERROR") + fmt.Sprintf("%d", task.TaskId))
		return
	}

	if m.player.GetLv() < configs.OpenLevel && m.player.GetVip() < configs.Openvip {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_LACK_OF_RANK_TASK_NOT"))
		return
	}

	//if m.player.GetVip() < configs.Openvip {
	//	m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_LACK_OF_ARISTOCRATIC_RANK_MISSION"))
	//	return
	//}

	if configs.Tasktypes != 2 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_TASK_TYPE_ERROR"))
		return
	}

	if len(configs.Ns) != 4 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_CONFIGURATION_ERROR"))
		return
	}

	costId := configs.Ns[0]
	config := GetCsvMgr().GetTariffConfig3(TariffGiveItems, costId)
	if config == nil {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_TARIFF_CONFIGURATION_ERROR"))
		return
	}

	err = m.player.HasObjectOk(config.ItemIds, config.ItemNums)
	if err != nil {
		m.player.SendErr(err.Error())
		return
	}

	//! 触发赏金任务次数，3
	//m.player.HandleTask(FinishTask, 3, 0, 0)
	task.Finish = 1
	task.Open = 1
	items := m.player.RemoveObjectLst(config.ItemIds, config.ItemNums, "赏金任务完成", LOG_MONEY_TASK_ITEM, 0, 0)
	m.synInfo("money_task_give_items", items)

	// 1 提交材料完成
	GetServer().SqlLog(m.player.Sql_UserBase.Uid, LOG_MONEY_TASK_FINISH, task.TaskId, LOG_MONEY_TASK_ITEM, 0, "赏金任务完成", 0, 0, m.player)
}

// 一键翻牌
func (m *ModMoneyTask) AwardAllMoneyTask(body []byte) {
	var msg C2S_MoneyTaskAction
	err := json.Unmarshal(body, &msg)
	if err != nil {
		m.player.SendErr(err.Error())
		return
	}
	task, ok := m.Data.info[msg.TaskId]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_DOES_NOT_EXIST"))
		return
	}

	if task.Finish == 0 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_NOT_COMPLETED"))
		return
	}

	config, ok := GetCsvMgr().MoneyTaskMap[task.TaskId]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MONEY_TASK_CONFIG_ERROR") + fmt.Sprintf("%d", task.TaskId))
		return
	}

	if len(config.Drawcosts) != 4 {
		m.player.SendErr("len(config.Drawcosts) != 4")
		return
	}

	// 检查一键翻牌的总消耗
	totalCost := 0
	for i := 0; i < 4; i++ {
		if task.BoxState[i] == 1 {
			continue
		}
		totalCost += config.Drawcosts[i]
	}

	if totalCost <= 0 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_THE_REWARD_OF_TURNING_THE"))
		return
	}

	//! 固定八折
	totalCost = totalCost * 8 / 10

	// 检测消耗
	if err := m.player.HasObjectOkEasy(DEFAULT_GEM, totalCost); err != nil {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_DIAMOND_SHORTAGE"))
		return
	}

	const logName = "赏金任务一键翻牌"
	items := m.player.RemoveObjectSimple(DEFAULT_GEM, totalCost, logName, task.TaskId, 0, 0)

	// 掉落信息
	//lootItems := GetLootMgr().LootItem(config.Group)
	//for _, item := range lootItems {
	//	items = append(items, PassItem{item.ItemId, item.ItemNum})
	//}
	var awardTypes []int
	for i := 0; i < 4; i++ {
		if task.BoxState[i] == 1 {
			continue
		}

		lootGroupId := config.Groups[i]
		//lootGroupId := config.Group
		var awardTypes []int
		DrawRate := config.Draws[i]
		curRand := HF_GetRandom(10000)
		if curRand < DrawRate {
			//! 掉落特殊掉落，否则掉落普通掉落
			//lootGroupId = config.Groups[i]
			lootGroupId = config.Group
			awardTypes = append(awardTypes, 1)
		} else {
			awardTypes = append(awardTypes, 2)
		}

		// 掉落信息
		var lootItems = make(map[int]*Item)
		if i < 4 {
			lootItems = GetLootMgr().LootItem(lootGroupId, nil)
		}

		//var awardItem PassItem
		for _, item := range lootItems {
			items = append(items, PassItem{item.ItemId, item.ItemNum})
			m.player.AddObjectSimple(item.ItemId, item.ItemNum, logName, task.TaskId, 0, 0)
			task.ItemsRecord[i] = PassItem{item.ItemId, item.ItemNum}
		}

		//must := m.player.AddObjectSimple(config.Mustitem, config.Num, logName)
		//task.ItemsRecord[i] = PassItem{config.Mustitem, config.Num}
		//worseRand := HF_GetRandom(10000)
		//if worseRand < 1 { // 中2等奖
		//	res := m.player.AddObjectSimple(config.Worstitem, config.Worstnum, logName)
		//	items = append(items, res...)
		//	awardTypes = append(awardTypes, 2)
		//	m.CheckEquipNotice(config.Worstitem)
		//	task.ItemsRecord[i] = PassItem{config.Worstitem, config.Worstnum}
		//
		//} else {
		//	surpriseRand := HF_GetRandom(100000)
		//	if surpriseRand < config.Surprisevalue { // 中1等奖
		//		res := m.player.AddObjectSimple(config.Surpriseitem, config.Surprisenum, logName)
		//		items = append(items, res...)
		//		awardTypes = append(awardTypes, 1)
		//		m.CheckEquipNotice(config.Surpriseitem)
		//		task.ItemsRecord[i] = PassItem{config.Surpriseitem, config.Surprisenum}
		//
		//	} else { // 不中奖
		//		items = append(items, must...)
		//		awardTypes = append(awardTypes, 3)
		//	}
		//}

		task.BoxState[i] = 1
	}

	task.Pickup = 1
	var msgRet S2C_MoneyTaskInfo
	msgRet.Cid = "award_all_money_task"
	msgRet.FlushTimes = m.Data.FlushTimes
	msgRet.Star = m.Data.Star
	msgRet.Info = m.Data.info
	msgRet.Items = items
	msgRet.AwardTypes = awardTypes
	m.player.SendMsg(msgRet.Cid, HF_JtoB(&msgRet))

	GetServer().SqlLog(m.player.Sql_UserBase.Uid, LOG_MONEY_TASK_DRAW_ONEKEY, task.TaskId, 0, 0, "赏金任务一键翻牌", 0, 0, m.player)
}

func (m *ModMoneyTask) DoneMoneyTask(body []byte) {
	var msg C2S_MoneyTaskAction
	err := json.Unmarshal(body, &msg)
	if err != nil {
		m.player.SendErr(err.Error())
		return
	}

	task, ok := m.Data.info[msg.TaskId]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_DOES_NOT_EXIST"))
		return
	}

	if task.Finish == 1 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_COMPLETED"))
		return
	}

	config, ok := GetCsvMgr().MoneyTaskMap[task.TaskId]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MONEY_TASK_CONFIG_ERROR") + fmt.Sprintf("%d", task.TaskId))
		return
	}

	if m.player.GetLv() < config.OpenLevel && m.player.GetVip() < config.Openvip {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_LACK_OF_RANK_TASK_NOT"))
		return
	}

	//if m.player.GetVip() < config.Openvip {
	//	m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_LACK_OF_ARISTOCRATIC_RANK_MISSION"))
	//	return
	//}

	nCostType := MONEY_TASK_COST_TYPE_REPLACE_ITEM

	//判断物品够不够 足够就是替代物品完成
	if err := m.player.HasObjectOkEasy(config.Replaceitem, config.Replacenum); err != nil {
		// 不够则判断钻石够不够
		if err := m.player.HasObjectOkEasy(DEFAULT_GEM, config.Moment); err != nil {
			//不够则返回
			m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_DIAMOND_SHORTAGE"))
			return
		}

		// 足够则是钻石完成
		nCostType = MONEY_TASK_COST_TYPE_GEM
	}

	//! 触发赏金任务次数，3
	//m.player.HandleTask(FinishTask, 3, 0, 0)

	var items []PassItem
	if nCostType == MONEY_TASK_COST_TYPE_GEM {
		items = m.player.RemoveObjectSimple(DEFAULT_GEM, config.Moment, "赏金任务完成", task.TaskId, LOG_MONEY_TASK_GEM, 0)
		// 2 钻石完成
		GetServer().SqlLog(m.player.Sql_UserBase.Uid, LOG_MONEY_TASK_FINISH, task.TaskId, LOG_MONEY_TASK_GEM, 0, "赏金任务完成", 0, 0, m.player)
	} else {
		items = m.player.RemoveObjectSimple(config.Replaceitem, config.Replacenum, "赏金任务完成", task.TaskId, LOG_MONEY_TASK_REPLACEITEM, 0)
		// 3 替代物品完成
		GetServer().SqlLog(m.player.Sql_UserBase.Uid, LOG_MONEY_TASK_FINISH, task.TaskId, LOG_MONEY_TASK_REPLACEITEM, 0, "赏金任务完成", 0, 0, m.player)
	}

	task.Finish = 1
	task.Open = 1
	m.synInfo("done_money_task", items)
}

// 检查任务是不是首次翻牌
func (m *MoneyTask) DrawTimes() int {
	num := 0
	for _, v := range m.BoxState {
		if v == 1 {
			num += 1
		}
	}
	return num
}

// 领取赏金任务奖励, 第n个箱子
func (m *ModMoneyTask) AwardMoneyTask(body []byte) {
	var msg C2S_MoneyTaskAction
	err := json.Unmarshal(body, &msg)
	if err != nil {
		m.player.SendErr(err.Error())
		return
	}
	task, ok := m.Data.info[msg.TaskId]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_DOES_NOT_EXIST"))
		return
	}

	if task.Finish == 0 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_NOT_COMPLETED"))
		return
	}

	config, ok := GetCsvMgr().MoneyTaskMap[task.TaskId]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MONEY_TASK_CONFIG_ERROR") + fmt.Sprintf("%d", task.TaskId))
		return
	}

	drawTimes := task.DrawTimes()
	if drawTimes >= 4 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_THE_NUMBER_OF_DOUBLES_HAS"))
		return
	}

	if len(task.BoxState) <= 0 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_TREASURE_BOX_STATUS_ERROR"))
		return
	}

	drawIndex := drawTimes
	// 检测消耗
	cost := config.Drawcosts[drawIndex]
	if err := m.player.HasObjectOkEasy(DEFAULT_GEM, cost); err != nil {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_DIAMOND_SHORTAGE"))
		return
	}

	const logName = "赏金任务翻牌"
	items := m.player.RemoveObjectSimple(DEFAULT_GEM, cost, logName, task.TaskId, drawTimes, 0)

	lootGroupId := config.Groups[drawIndex]
	//lootGroupId := config.Group
	var awardTypes []int
	DrawRate := config.Draws[drawIndex]
	curRand := HF_GetRandom(10000)
	if curRand < DrawRate {
		//! 掉落特殊掉落，否则掉落普通掉落
		//lootGroupId = config.Groups[drawIndex]
		lootGroupId = config.Group
		awardTypes = append(awardTypes, 1)
	} else {
		awardTypes = append(awardTypes, 2)
	}

	task.BoxState[drawIndex] = 1

	// 掉落信息
	var lootItems = make(map[int]*Item)
	if drawTimes < 4 {
		lootItems = GetLootMgr().LootItem(lootGroupId, nil)
	}

	//var awardItem PassItem
	for _, item := range lootItems {
		items = append(items, PassItem{item.ItemId, item.ItemNum})
		m.player.AddObjectSimple(item.ItemId, item.ItemNum, logName, task.TaskId, drawTimes, 0)
		task.ItemsRecord[drawIndex] = PassItem{item.ItemId, item.ItemNum}
	}

	//must := m.player.AddObjectSimple(config.Mustitem, config.Num, logName)
	//task.ItemsRecord[drawIndex] = PassItem{config.Surpriseitem, config.Surprisenum}
	//items = append(items, must...)

	//var awardTypes []int
	//worseRand := HF_GetRandom(10000)
	//if worseRand < 1 { // 中2等奖
	//	res := m.player.AddObjectSimple(config.Worstitem, config.Worstnum, logName)
	//	items = append(items, res...)
	//	awardTypes = append(awardTypes, 2)
	//	m.CheckEquipNotice(config.Worstitem)
	//	task.ItemsRecord[drawIndex] = PassItem{config.Worstitem, config.Worstnum}
	//
	//} else {
	//	surpriseRand := HF_GetRandom(100000)
	//	if surpriseRand < config.Surprisevalue { // 中1等奖
	//		res := m.player.AddObjectSimple(config.Surpriseitem, config.Surprisenum, logName)
	//		items = append(items, res...)
	//		awardTypes = append(awardTypes, 1)
	//		m.CheckEquipNotice(config.Surpriseitem)
	//		task.ItemsRecord[drawIndex] = PassItem{config.Surpriseitem, config.Surprisenum}
	//	} else { // 不中奖
	//		items = append(items, must...)
	//		awardTypes = append(awardTypes, 3)
	//	}
	//}

	// 赋值
	if task.Pickup == 0 {
		task.Pickup = 1
	}
	task.BoxState[drawIndex] = 1
	var msgRet S2C_MoneyTaskInfo
	msgRet.Cid = "award_money_task"
	msgRet.FlushTimes = m.Data.FlushTimes
	msgRet.Star = m.Data.Star
	msgRet.Info = m.Data.info
	msgRet.Items = items
	msgRet.AwardTypes = awardTypes
	m.player.SendMsg(msgRet.Cid, HF_JtoB(&msgRet))

	GetServer().SqlLog(m.player.Sql_UserBase.Uid, LOG_MONEY_TASK_DRAW, task.TaskId, awardTypes[0], drawTimes, "赏金任务翻牌", 0, 0, m.player)
}

// 刷新任务
func (m *ModMoneyTask) MoneyTaskFlush(body []byte) {
	tConfig := GetCsvMgr().GetTariffConfig2(TariffbuffKingFlush)
	if tConfig == nil {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_CONFIGURATION_ERROR"))
		return
	}

	if err := m.player.HasObjectOk(tConfig.ItemIds, tConfig.ItemNums); err != nil {
		m.player.SendErr(err.Error())
		return
	}

	vipConfig := GetCsvMgr().GetVipConfig(m.player.Sql_UserBase.Vip)
	if vipConfig == nil {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_ARISTOCRATIC_MISALLOCATION"))
		return
	}

	if m.Data.FlushTimes >= vipConfig.Growthtask_King2 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_INSUFFICIENT_REFRESH_TIMES"))
		return
	}

	m.Data.FlushTimes += 1
	m.FlushTasks()
	items := m.player.RemoveObjectLst(tConfig.ItemIds, tConfig.ItemNums, "赏金任务刷新", 0, 0, 0)

	// 同步任务信息
	m.synInfo("money_task_flush", items)

	GetServer().SqlLog(m.player.Sql_UserBase.Uid, LOG_MONEY_TASK_REFRESH, 1, 0, 0, "赏金任务刷新", 0, 0, m.player)
}

// 同步消息
func (m *ModMoneyTask) synInfo(cid string, items []PassItem) {
	var msg S2C_MoneyTaskInfo
	msg.Cid = cid
	msg.FlushTimes = m.Data.FlushTimes
	msg.Star = m.Data.Star
	msg.Info = m.Data.info
	msg.Items = items
	m.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

// 登录时同步消息
func (m *ModMoneyTask) SendInfo() {
	const cid = "moneytask"
	msg := &S2C_MongyTaskInfo{
		Cid:  cid,
		Info: m.Data.info,
	}
	m.player.Send(cid, msg)
}

func (m *ModMoneyTask) NewMoneyTask(taskId int, sort int) *MoneyTask {
	config, ok := GetCsvMgr().MoneyTaskMap[taskId]
	if !ok {
		LogError("money task is not exist, taskId=", taskId)
		return nil
	}
	task := &MoneyTask{
		TaskId:   config.Taskid,
		Tasktype: config.Tasktypes,
		Plan:     m.FindPlanByType(config.Tasktypes),
		Finish:   0,
		Pickup:   0,
		Sort:     sort,
	}

	if task.Plan >= config.Ns[0] && m.IsTaskTypeOk(config.Tasktypes) {
		task.Finish = 1
		if m.player.GetLv() >= config.OpenLevel || m.player.GetVip() >= config.Openvip {
			task.Open = 1
		}
	}

	task.BoxState = []int{}
	for i := 0; i < 4; i++ {
		task.BoxState = append(task.BoxState, 0)
	}

	return task
}

// 刷出六个任务
func (m *ModMoneyTask) MakeTasks() {
	if m.Data.info == nil {
		m.Data.info = make(map[int]*MoneyTask)
	}

	if len(m.Data.info) >= MaxMoneyTaskNum {
		return
	}

	m.Data.info = make(map[int]*MoneyTask)
	for i := 0; i < MaxMoneyTaskNum; i++ {
		taskId := m.MakeTask()
		if taskId == 0 {
			LogError("taskId == 0, error may occur.")
			continue
		}
		task := m.NewMoneyTask(taskId, i+1)
		m.Data.info[taskId] = task
	}
}

func (m *ModMoneyTask) FlushTasks() {
	if m.Data.info == nil {
		m.Data.info = make(map[int]*MoneyTask)
	}

	m.Data.info = make(map[int]*MoneyTask)
	for i := 0; i < MaxMoneyTaskNum; i++ {
		taskId := m.MakeTask()
		if taskId == 0 {
			LogError("taskId == 0, error may occur.")
			continue
		}
		task := m.NewMoneyTask(taskId, i+1)
		m.Data.info[taskId] = task
	}
}

// 通过当前任务以及星级刷出一个没有的
func (m *ModMoneyTask) FilterTask() []*MongytaskListConfig {
	var taskConfigs []*MongytaskListConfig
	for _, v := range GetCsvMgr().MoneyTaskMap {
		_, ok := m.Data.info[v.Taskid]
		if ok {
			continue
		}

		//if m.player.GetLv() >= v.OpenLevel || m.player.GetVip() >= v.Openvip {
		taskConfigs = append(taskConfigs, v)
		//}
	}
	return taskConfigs
}

// 刷出一个任务
func (m *ModMoneyTask) MakeTask() int {
	taskConfigs := m.FilterTask()
	total := 0
	for _, v := range taskConfigs {
		total += v.Value
	}

	if total == 0 {
		LogError("total is 0")
		return 0
	}

	randNum := HF_GetRandom(total)
	check := 0
	for _, v := range taskConfigs {
		check += v.Value
		if randNum < check {
			return v.Taskid
		}
	}

	LogError("loot taskId 0")
	return 0
}

// 生成任务记录, 每天清0
func (m *ModMoneyTask) CheckRecord() {
	if m.Data.taskRecord == nil {
		m.Data.taskRecord = make(map[int]*MoneyRecord)
	}

	var taskTypes = []int{CommonLevelTask}
	for _, taskType := range taskTypes {
		_, ok := m.Data.taskRecord[taskType]
		if !ok {
			m.Data.taskRecord[taskType] = &MoneyRecord{taskType, 0}
		}
	}
}

// 增加记录
func (m *ModMoneyTask) AddRecord(taskType int, num int) {
	pRecord, ok := m.Data.taskRecord[taskType]
	if !ok {
		return
	}
	pRecord.Plan += num
}

// 清空记录
func (m *ModMoneyTask) ClearRecord() {
	for _, v := range m.Data.taskRecord {
		v.Plan = 0
	}
}

func (m *ModMoneyTask) IsTaskTypeOk(taskType int) bool {
	return taskType == CommonLevelTask
}

// 任务触发
// 活动任务触发, 活动记录加工
func (m *ModMoneyTask) HandleTask(taskType int, n1 int, n2 int, n3 int) {
	if !m.IsTaskTypeOk(taskType) {
		return
	}

	// 检查活动是否开启
	for _, pTask := range m.Data.taskRecord {
		if pTask == nil {
			return
		}

		if pTask.Tasktype != taskType {
			return
		}

		pTaskNode := &TaskNode{
			Tasktypes: pTask.Tasktype,
		}

		process, add := DoTask(pTaskNode, m.player, n1, n2, n3)
		if process == 0 {
			continue
		}

		if add { // 如果是记次任务
			pTask.Plan += process
		} else { // 如果是记值任务
			if process > pTask.Plan {
				pTask.Plan = process
			}
		}
	}
	m.RealTask(taskType)
}

func (m *ModMoneyTask) CheckTask() {
	for _, v := range m.Data.info {
		if v.Finish == 1 || v.Pickup == 1 {
			continue
		}

		pRecord, ok := m.Data.taskRecord[v.TaskId]
		if !ok {
			continue
		}

		v.Plan = pRecord.Plan
		config, ok := GetCsvMgr().MoneyTaskMap[v.TaskId]
		if !ok {
			LogError("config is error, taskId:", v.TaskId)
			continue
		}

		if len(config.Ns) <= 0 {
			continue
		}

		if v.Plan >= config.Ns[0] {
			// 需要判断玩家等级和vip
			v.Finish = 1
			if m.player.GetLv() >= config.OpenLevel || m.player.GetVip() >= config.Openvip {
				v.Open = 1
			}
		}
	}
}

// 通过加工的record检查活动任务是否完成
func (m *ModMoneyTask) RealTask(taskType int) {
	var tasks []*MoneyTask
	for _, v := range m.Data.info {
		if v.Tasktype != taskType {
			continue
		}

		if v.Finish == 1 || v.Pickup == 1 {
			continue
		}

		pRecord, ok := m.Data.taskRecord[v.Tasktype]
		if !ok {
			continue
		}

		config, ok := GetCsvMgr().MoneyTaskMap[v.TaskId]
		if !ok {
			LogError("task not exist, taskId:", v.TaskId)
			continue
		}

		if len(config.Ns) != 4 {
			continue
		}

		if taskType == CommonLevelTask {
			if pRecord.Plan >= config.Ns[0] {
				v.Plan = config.Ns[0]
				v.Finish = 1
				if m.player.GetLv() >= config.OpenLevel || m.player.GetVip() >= config.Openvip {
					v.Open = 1
				}
			} else {
				v.Plan = pRecord.Plan
			}
			tasks = append(tasks, v)
		}
	}
	m.SynTask(tasks)
}

func (m *ModMoneyTask) SynTask(tasks []*MoneyTask) {
	if len(tasks) <= 0 {
		return
	}

	msg := &S2C_SynMoneyTasks{"syn_money_task", tasks}
	m.player.Send(msg.Cid, msg)
}

func (m *ModMoneyTask) FindPlanByType(taskType int) int {
	pRecord, ok := m.Data.taskRecord[taskType]
	if !ok {
		return 0
	}

	return pRecord.Plan
}

// 升级vip和等级的时候检查open
func (m *ModMoneyTask) CheckOpen() {
	var tasks []*MoneyTask
	for _, v := range m.Data.info {
		if v.Open == 1 {
			continue
		}

		configs, ok := GetCsvMgr().MoneyTaskMap[v.TaskId]
		if !ok {
			m.player.SendErr(GetCsvMgr().GetText("STR_MONEY_TASK_CONFIG_ERROR") + fmt.Sprintf("%d", v.TaskId))
			return
		}

		if m.player.GetLv() < configs.OpenLevel && m.player.GetVip() < configs.Openvip {
			//m.player.SendErr("等级不足，任务没解锁")
			return
		}

		//if m.player.GetVip() < configs.Openvip {
		//	//m.player.SendErr("贵族等级不足，任务没解锁")
		//	return
		//}

		v.Open = 1
		tasks = append(tasks, v)
	}
	m.SynTask(tasks)
}

func (m *ModMoneyTask) CheckEquipNotice(equipId int) {
	config := GetCsvMgr().GetEquipConfig(equipId)
	if config == nil {
		return
	}

	if config.Quality <= 4 {
		return
	}
	//GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_MONEY_TASK_NOTICE"), m.player.GetName(), config.EquipName), 0, 0)
}
