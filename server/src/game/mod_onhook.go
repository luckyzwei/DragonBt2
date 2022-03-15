package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

const (
	ONHOOK_DROP_SOURCE = 0 //资源掉落 配置60S 间隔20S
	ONHOOK_DROP_BASE   = 1 //基础掉落  配置=间隔
	ONHOOK_DROP_SENIOR = 2 //高级掉落   配置=间隔
	ONHOOK_DROP_END    = 3 //
)

const (
	ONHOOK_INIT_CONFIG_INDEX = 1 //默认奖励
	ONHOOK_INIT_LEVEL        = 110101
	ONHOOK_VIP_RATE          = 2 //需要计算加成
)

//! 挂机数据库 TODAY
type San_OnHook struct {
	Uid              int64
	GetTime          int64  //	领取的时间
	OnHookStage      int    //	当前挂机的关卡
	CalTime          string //计算掉落组的多余时间
	HangUp           int    //用户挂机掉落组
	FastTimes        int    //快速挂机次数
	OnHookStageTime  int64  //当前挂机的关卡到达时间
	PassIdRecord     int    //记录当日刷新时的关卡进度，用来判断礼包奖励
	PassIdRecordTime int64  //
	CalTimePrivilege string //特权掉落组的多余时间
	CalTimeExtItems  string //特权掉落额外奖励

	calTime          []int64       //计算掉落组的多余时间
	calTimePrivilege map[int]int64 //特权掉落组
	calTimeExtItems  map[int]*Item //特权掉落额外奖励
	DataUpdate
}

//! 挂机
type ModOnHook struct {
	player     *Player
	Sql_OnHook San_OnHook //! 数据库结构
}

func (self *ModOnHook) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_onhook` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_OnHook, "san_onhook", self.player.ID)

	if self.Sql_OnHook.Uid <= 0 {
		self.Sql_OnHook.Uid = self.player.ID
		self.Sql_OnHook.calTimePrivilege = make(map[int]int64, 0)
		self.Sql_OnHook.calTimeExtItems = make(map[int]*Item, 0)
		self.Encode()
		InsertTable("san_onhook", &self.Sql_OnHook, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_OnHook.Init("san_onhook", &self.Sql_OnHook, true)
	self.Check()
}

func (self *ModOnHook) Check() {
	if self.Sql_OnHook.GetTime == 0 {
		self.Sql_OnHook.GetTime = TimeServer().Unix()
	}
	if self.Sql_OnHook.calTimePrivilege == nil {
		self.Sql_OnHook.calTimePrivilege = make(map[int]int64, 0)
	}
	if self.Sql_OnHook.calTimeExtItems == nil {
		self.Sql_OnHook.calTimeExtItems = make(map[int]*Item, 0)
	}
	if self.Sql_OnHook.OnHookStage == 0 {
		self.Sql_OnHook.OnHookStage = ONHOOK_INIT_LEVEL
		self.Sql_OnHook.OnHookStageTime = TimeServer().Unix()
		self.Sql_OnHook.HangUp = ONHOOK_INIT_CONFIG_INDEX
	}
	if len(self.Sql_OnHook.calTime) < ONHOOK_DROP_END {
		size := len(self.Sql_OnHook.calTime)
		for i := size; i < ONHOOK_DROP_END; i++ {
			self.Sql_OnHook.calTime = append(self.Sql_OnHook.calTime, 0)
		}
	}
}

func (self *ModOnHook) OnRefresh() {
	self.Sql_OnHook.FastTimes = 0
}

func (self *ModOnHook) Decode() {
	json.Unmarshal([]byte(self.Sql_OnHook.CalTime), &self.Sql_OnHook.calTime)
	json.Unmarshal([]byte(self.Sql_OnHook.CalTimePrivilege), &self.Sql_OnHook.calTimePrivilege)
	json.Unmarshal([]byte(self.Sql_OnHook.CalTimeExtItems), &self.Sql_OnHook.calTimeExtItems)
}

func (self *ModOnHook) Encode() {
	self.Sql_OnHook.CalTime = HF_JtoA(self.Sql_OnHook.calTime)
	self.Sql_OnHook.CalTimePrivilege = HF_JtoA(self.Sql_OnHook.calTimePrivilege)
	self.Sql_OnHook.CalTimeExtItems = HF_JtoA(self.Sql_OnHook.calTimeExtItems)
}

func (self *ModOnHook) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModOnHook) OnGetData(player *Player) {
	self.player = player
}

func (self *ModOnHook) OnSave(sql bool) {
	self.Encode()
	self.Sql_OnHook.Update(sql)
}

//计算奖励  isCalPlayer 是否计算领主经验
func (self *ModOnHook) CalValue(sec int64, isCalPlayer bool, isFast bool) (map[int]*Item, map[int]*Item, map[int]*Item, map[int]*Item) {
	rewards := make(map[int]*Item, 0)
	privilegeRewards := make(map[int]*Item, 0)
	monthCardRewards := make(map[int]*Item, 0)
	activityRewards := make(map[int]*Item, 0)
	config, ok := GetCsvMgr().HangUpConfig[self.Sql_OnHook.HangUp]
	if !ok {
		config = GetCsvMgr().HangUpConfig[ONHOOK_INIT_CONFIG_INDEX]
		if config == nil {
			return rewards, privilegeRewards, monthCardRewards, activityRewards
		}
	}
	//计算资源掉落
	self.CalSource(rewards, sec, config, isCalPlayer, privilegeRewards, monthCardRewards)
	//计算基础掉落
	self.CalBase(rewards, sec, config)
	//计算高级掉落
	self.CalSenior(rewards, sec, config, activityRewards, isFast)
	//计算特权掉落
	PrivilegeValue := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValues()
	self.CalPrivilege(privilegeRewards, sec, PrivilegeValue)
	return rewards, privilegeRewards, monthCardRewards, activityRewards
}

func (self *ModOnHook) CalSource(rewards map[int]*Item, sec int64, config *HangUpConfig, isCalPlayer bool, privilegeRewards map[int]*Item, monthCardRewards map[int]*Item) {
	if len(self.Sql_OnHook.calTime) < ONHOOK_DROP_SOURCE {
		return
	}
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		return
	}
	realtime := sec + self.Sql_OnHook.calTime[ONHOOK_DROP_SOURCE]
	self.Sql_OnHook.calTime[ONHOOK_DROP_SOURCE] = realtime % 20
	times := realtime / 20
	//表里配置的分钟以兼容显示 计算需要/3
	//计算金币
	if config.GoldYield > 0 {
		num := config.GoldYield * int(times) * (PER_BIT + vipcsv.HangupGold) / (3 * PER_BIT)
		AddItemMapHelper3(rewards, ITEM_GOLD, num)
		value := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValue(PRIVILEGE_GOLD)
		if value > 0 {
			baseNum := config.GoldYield * int(times) / 3
			AddItemMapHelper3(privilegeRewards, ITEM_GOLD, baseNum*value/100)
		}

		valueMonth := self.player.GetModule("activity").(*ModActivity).GetMonthGold()
		if valueMonth > 0 {
			baseNum := config.GoldYield * int(times) / 3
			AddItemMapHelper3(monthCardRewards, ITEM_GOLD, baseNum*valueMonth/PER_BIT)
		}
	}

	//计算英雄经验
	if config.HeroExpYield > 0 {
		num := config.HeroExpYield * int(times) * (PER_BIT + vipcsv.HangupHeroExp) / (3 * PER_BIT)
		AddItemMapHelper3(rewards, ITEM_HERO_EXP, num)
		value := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValue(PRIVILEGE_HERO_EXP)
		if value > 0 {
			baseNum := config.HeroExpYield * int(times) / 3
			AddItemMapHelper3(privilegeRewards, ITEM_HERO_EXP, baseNum*value/100)
		}

		valueMonth := self.player.GetModule("activity").(*ModActivity).GetMonthHeroExp()
		if valueMonth > 0 {
			baseNum := config.HeroExpYield * int(times) / 3
			AddItemMapHelper3(monthCardRewards, ITEM_HERO_EXP, baseNum*valueMonth/PER_BIT)
		}
	}
	//计算领主经验
	if isCalPlayer && config.ExpYield > 0 {
		AddItemMapHelper3(rewards, ITEM_PLAYER_EXP, config.ExpYield*int(times)/3)
	}
	//计算魔粉
	if config.PowderYield > 0 {
		num := config.PowderYield * int(times) / (3 * 60)
		AddItemMapHelper3(rewards, ITEM_POWDER, num)
		value := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValue(PRIVILEGE_POWDER)
		if value > 0 {
			AddItemMapHelper3(privilegeRewards, ITEM_POWDER, num*value/100)
		}
	}
}

func (self *ModOnHook) CalBase(rewards map[int]*Item, sec int64, config *HangUpConfig) {
	if len(self.Sql_OnHook.calTime) < ONHOOK_DROP_BASE {
		return
	}

	if config.BasicsTime <= 0 {
		return
	}

	realtime := sec + self.Sql_OnHook.calTime[ONHOOK_DROP_BASE]
	self.Sql_OnHook.calTime[ONHOOK_DROP_BASE] = realtime % config.BasicsTime
	times := realtime / config.BasicsTime

	for i := 0; i < int(times); i++ {
		for j := 0; j < len(config.BasicsDrop); j++ {
			if config.BasicsDrop[j] == 0 {
				break
			}
			outitem := GetLootMgr().LootItem(config.BasicsDrop[j], self.player)
			AddItemMapHelper4(rewards, outitem)
		}
	}
}

//activityRewards:新增活动掉落组，和高级掉落组共用算法,返回活动掉落组的实际计算时间
func (self *ModOnHook) CalSenior(rewards map[int]*Item, sec int64, config *HangUpConfig, activityRewards map[int]*Item, isFast bool) int64 {
	if len(self.Sql_OnHook.calTime) < ONHOOK_DROP_SENIOR {
		return 0
	}

	if config.SeniorTime <= 0 {
		return 0
	}

	realtime := sec + self.Sql_OnHook.calTime[ONHOOK_DROP_SENIOR]
	self.Sql_OnHook.calTime[ONHOOK_DROP_SENIOR] = realtime % config.SeniorTime
	times := realtime / config.SeniorTime

	for i := 0; i < int(times); i++ {
		for j := 0; j < len(config.SeniorDrop); j++ {
			if config.SeniorDrop[j] == 0 {
				break
			}
			outitem := GetLootMgr().LootItem(config.SeniorDrop[j], self.player)
			AddItemMapHelper4(rewards, outitem)
		}
	}

	//活动掉落
	activity := GetActivityMgr().GetActivity(ACT_ONHOOK_ACTIVITY_SPRING_FESTIVAL)
	if activity == nil || activity.status.Status != ACTIVITY_STATUS_OPEN {
		return 0
	}
	realtimeActivity := realtime
	now := TimeServer().Unix()
	//不是快速挂机的话 需要修正时间
	if !isFast {
		if activity.status.Status == ACTIVITY_STATUS_OPEN {
			//如果活动中
			activityAlreadyStartTime := now - HF_CalTimeForConfig(activity.info.Start, self.player.Sql_UserBase.Regtime)
			if activityAlreadyStartTime > 0 && realtime > activityAlreadyStartTime {
				realtimeActivity = activityAlreadyStartTime
			}
		} else {
			//活动未开放的情况下
			activityAlreadyEndTime := now - (HF_CalTimeForConfig(activity.info.Start, self.player.Sql_UserBase.Regtime) + int64(activity.info.Continued))
			if activityAlreadyEndTime < 0 || realtimeActivity < activityAlreadyEndTime {
				realtimeActivity = 0
			} else {
				realtimeActivity -= activityAlreadyEndTime
			}
		}
	}
	timesActivity := realtimeActivity / config.SeniorTime
	for i := 0; i < int(timesActivity); i++ {
		for _, itemConfig := range activity.items {
			if itemConfig.N[0] > 0 {
				outitem := GetLootMgr().LootItem(itemConfig.N[0], self.player)
				AddItemMapHelper4(activityRewards, outitem)
			}
		}
	}
	return realtimeActivity
}

func (self *ModOnHook) CalPrivilege(privilegeRewards map[int]*Item, sec int64, privilegeValue map[int]int) {

	for k, v := range privilegeValue {
		if k < PRIVILEGE_ONHOOK_MIN || k > PRIVILEGE_ONHOOK_MAX {
			continue
		}

		config := GetCsvMgr().InterstellarHangup[v]
		if config == nil {
			continue
		}

		realtime := sec + self.Sql_OnHook.calTimePrivilege[config.PrivileGeType]

		switch config.Type {
		case 1:
			//先得出产量
			num := realtime * int64(config.Num) / config.InterstellarTime
			if num <= 0 {
				continue
			} else {
				AddItemMapHelper3(privilegeRewards, config.Item, int(num))
				realtime -= int64(num) * config.InterstellarTime / int64(config.Num)
			}
		case 2:
			self.Sql_OnHook.calTimePrivilege[config.PrivileGeType] = realtime % config.InterstellarTime
			times := realtime / config.InterstellarTime
			//先获得掉落组
			DropGroup := config.InterstellarDrop
			for i := 0; i < int(times); i++ {
				outitem := GetLootMgr().LootItem(DropGroup, self.player)
				AddItemMapHelper4(privilegeRewards, outitem)
			}
		case 3:
			self.Sql_OnHook.calTimePrivilege[config.PrivileGeType] = realtime % config.InterstellarTime
			times := realtime / config.InterstellarTime
			//先获得掉落组
			DropGroup := config.InterstellarDrop
			for i := len(config.Judge) - 1; i > 0; i-- {
				if config.Judge[i] == 0 {
					continue
				}
				if privilegeValue[config.Judge[i]] > 0 {
					DropGroup = config.ChangeDrop[i]
					break
				}
			}

			for i := 0; i < int(times); i++ {
				outitem := GetLootMgr().LootItem(DropGroup, self.player)
				AddItemMapHelper4(privilegeRewards, outitem)
			}
		}
	}
}

func (self *ModOnHook) SendInfo(body []byte) {
	var msg S2C_OnHookInfo
	msg.Cid = "onhookinfo"
	msg.Time = TimeServer().Unix() - self.Sql_OnHook.GetTime
	msg.Stage = self.Sql_OnHook.OnHookStage
	msg.HangUp = self.Sql_OnHook.HangUp
	msg.FastTimes = self.Sql_OnHook.FastTimes
	msg.StageTime = self.Sql_OnHook.OnHookStageTime
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModOnHook) SendInfoAutoSend() {
	var msg S2C_OnHookInfo
	msg.Cid = "onhookinfo"
	msg.Time = TimeServer().Unix() - self.Sql_OnHook.GetTime
	msg.Stage = self.Sql_OnHook.OnHookStage
	msg.HangUp = self.Sql_OnHook.HangUp
	msg.FastTimes = self.Sql_OnHook.FastTimes
	msg.StageTime = self.Sql_OnHook.OnHookStageTime
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModOnHook) onReg(handlers map[string]func(body []byte)) {
	handlers["awardonhook"] = self.AwardOnHook
	handlers["onhookstage"] = self.OnHookStage
	handlers["onhookinfo"] = self.SendInfo
	handlers["onhookfast"] = self.OnHookFast
}

// 设置当前目标关卡
func (self *ModOnHook) SetStage(passId int) {
	if self.Sql_OnHook.OnHookStage == passId {
		return
	}
	//之前的表示刚过关,更新掉落库
	levelConfig := GetCsvMgr().LevelConfigMap[self.Sql_OnHook.OnHookStage]
	if levelConfig != nil {
		self.Sql_OnHook.HangUp = levelConfig.HangUp
	}
	GetServer().SendLog_SDKUP_AIWAN_LOGIN(self.player, SKDUP_ADDR_URL_AIWAN_SDK_EVENT_LEVELUP)
	self.Sql_OnHook.OnHookStage = passId
	self.Sql_OnHook.OnHookStageTime = TimeServer().Unix()
	self.player.Sql_UserBase.PassMax = passId
	self.player.NoticeCenterBaseInfo()
	//更新玩家信息管理
	GetOfflineInfoMgr().SetPlayerStage(self.player.Sql_UserBase.Uid, self.Sql_OnHook.OnHookStage)
	if self.player.GetUnionId() > 0 {
		GetUnionMgr().UpdateMemberPassID(self.player.GetUnionId(), self.player.Sql_UserBase.Uid, self.Sql_OnHook.OnHookStage)
	}
	GetTopPassMgr().UpdateRank(int64(self.Sql_OnHook.OnHookStage), self.player)
	self.player.GetModule("recharge").(*ModRecharge).CheckOpen()
	self.player.GetModule("recharge").(*ModRecharge).CheckOpenLimit()
	self.player.GetModule("newpit").(*ModNewPit).CheckOpen()

	var msgRel S2C_OnHookStage
	msgRel.Cid = "setonhookstage"
	msgRel.Stage = self.Sql_OnHook.OnHookStage
	msgRel.HangUp = self.Sql_OnHook.HangUp
	msgRel.StageTime = self.Sql_OnHook.OnHookStageTime
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModOnHook) GetStage() int {
	return self.Sql_OnHook.OnHookStage
}

func (self *ModOnHook) GetDailyPassId() int {
	if self.Sql_OnHook.PassIdRecordTime < TimeServer().Unix() {
		self.Sql_OnHook.PassIdRecordTime = HF_GetNextDayStart()
		self.Sql_OnHook.PassIdRecord = self.Sql_OnHook.OnHookStage
	}
	return self.Sql_OnHook.PassIdRecord
}

func (self *ModOnHook) AwardOnHook(body []byte) {

	sec := TimeServer().Unix() - self.Sql_OnHook.GetTime
	if sec <= 0 {
		return
	}

	//先看看挂机时长是否超过了限制  单位小时
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REBORN_ERROR_CONFIG"))
		return
	}
	if sec > int64(vipcsv.HangupTime) {
		sec = int64(vipcsv.HangupTime)
	}
	//计算奖励
	rewards, privilegeRewards, monthCardRewards, activityRewards := self.CalValue(sec, true, false)
	self.Sql_OnHook.GetTime = TimeServer().Unix()
	getItems := self.player.AddObjectItemMap(rewards, "领取挂机奖励", self.Sql_OnHook.HangUp, int(sec), 0)
	if len(self.Sql_OnHook.calTimeExtItems) > 0 {
		AddItemMapHelper4(privilegeRewards, self.Sql_OnHook.calTimeExtItems)
		self.Sql_OnHook.calTimeExtItems = make(map[int]*Item, 0)
	}
	getPrivilegeItems := self.player.AddObjectItemMap(privilegeRewards, "领取挂机科技奖励", 0, int(sec), 0)
	getMonthCardItems := self.player.AddObjectItemMap(monthCardRewards, "领取挂机月卡奖励", 0, int(sec), 0)
	getActivityItems := self.player.AddObjectItemMap(activityRewards, "领取挂机活动奖励", 0, int(sec), 0)
	//self.player.HandleTask(GetOnHookAwardTask, 0, 0, 0)
	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ONHOOK_AWARD, OPEN_LEVEL_ON_HOOK, 0, 0, "挂机奖励领取", 0, 0, self.player)

	var msg S2C_AwardOnHook
	msg.Cid = "awardonhook"
	msg.Time = 0
	msg.GetItem = getItems
	msg.GetTime = sec
	msg.GetPrivilegeItems = getPrivilegeItems
	msg.GetMonthItems = getMonthCardItems
	msg.GetActivityItems = getActivityItems
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ONHOOK_AWARD, self.Sql_OnHook.HangUp, int(sec), 0, "领取挂机奖励", 0, 0, self.player)

	self.player.HandleTask(TASK_TYPE_GET_HOOK_AWARD, 1, 0, 0)

	for _, v := range getItems {
		self.player.HandleTask(TASK_TYPE_GET_HOOK_GET_ITEM, v.ItemID, v.Num, 0)
	}
	for _, v := range getPrivilegeItems {
		self.player.HandleTask(TASK_TYPE_GET_HOOK_GET_ITEM, v.ItemID, v.Num, 0)
	}
	for _, v := range getMonthCardItems {
		self.player.HandleTask(TASK_TYPE_GET_HOOK_GET_ITEM, v.ItemID, v.Num, 0)
	}

	self.player.SendTask()
}

// 玩家选择挂机关卡
func (self *ModOnHook) OnHookStage(body []byte) {
	var msg C2S_OnHookStage
	json.Unmarshal(body, &msg)

	//看前置关卡是否通过 之前没通过也没验前置关卡 这个地方暂时放着
	config := GetCsvMgr().LevelConfigMap[self.Sql_OnHook.OnHookStage]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_GATE_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	if config.NextLevel != msg.Stage {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_MSG_STAGE_ERROR"))
		return
	}
	configNext := GetCsvMgr().LevelConfigMap[msg.Stage]
	if configNext == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_GATE_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	if config.LevelIndex == configNext.LevelIndex {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_MSG_STAGE_ERROR"))
		return
	}

	self.Sql_OnHook.OnHookStage = msg.Stage
	self.Sql_OnHook.OnHookStageTime = TimeServer().Unix()

	var msgRel S2C_OnHookStage
	msgRel.Cid = "onhookstage"
	msgRel.Stage = self.Sql_OnHook.OnHookStage
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

//快速挂机
func (self *ModOnHook) OnHookFast(body []byte) {

	//是否还有次数
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REBORN_ERROR_CONFIG"))
		return
	}

	numFast := self.player.GetModule("activity").(*ModActivity).GetMonthFreeFast()
	calFastTimes := self.Sql_OnHook.FastTimes - numFast
	if calFastTimes < 0 {
		calFastTimes = 0
	}

	if calFastTimes >= vipcsv.HangupFast {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ONHOOK_VIP_NOT_ENOUGH"))
		return
	}

	buyConfig := GetCsvMgr().GetTariffConfig(TARIFF_TYPE_ONHOOK_FAST, calFastTimes+1)
	if buyConfig == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_ARMY_MAXIMUM_NUMBER_OF_PURCHASES"))
		return
	}

	//检查消耗够不够
	if err := self.player.HasObjectOk(buyConfig.ItemIds, buyConfig.ItemNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	self.Sql_OnHook.FastTimes += 1
	// 扣除物品
	costItem := self.player.RemoveObjectLst(buyConfig.ItemIds, buyConfig.ItemNums, "使用快速挂机", self.Sql_OnHook.HangUp, 0, 1)

	//2小时收益
	sec := int64(2 * HOUR_SECS)
	//计算奖励
	//快速挂机不适用活动掉落 修改于20210204
	rewards, privilegeRewards, monthCardRewards, _ := self.CalValue(sec, false, true)
	num := GetGemNum(costItem)
	param3 := 0
	if num > 0 {
		param3 = -1
	}
	getItems := self.player.AddObjectItemMap(rewards, "使用快速挂机", self.Sql_OnHook.HangUp, 0, param3)
	getPrivilegeItems := self.player.AddObjectItemMap(privilegeRewards, "快速挂机科技奖励", 0, 0, param3)
	getMonthCardItems := self.player.AddObjectItemMap(monthCardRewards, "使用快速挂机", 0, int(sec), 0)
	//快速挂机不适用活动掉落 修改于20210204
	getActivityItems := make([]PassItem, 0)
	//getActivityItems := self.player.AddObjectItemMap(activityRewards, "使用快速挂机", 0, int(sec), 0)
	//CheckAddItemLog(self.player, "快速挂机", costItem, getItems)
	if num > 0 {
		AddSpecialSdkItemListLog(self.player, num, getItems, "使用快速挂机")
	}

	for _, v := range getItems {
		self.player.HandleTask(TASK_TYPE_GET_HOOK_GET_ITEM, v.ItemID, v.Num, 0)
	}
	for _, v := range getPrivilegeItems {
		self.player.HandleTask(TASK_TYPE_GET_HOOK_GET_ITEM, v.ItemID, v.Num, 0)
	}
	for _, v := range getMonthCardItems {
		self.player.HandleTask(TASK_TYPE_GET_HOOK_GET_ITEM, v.ItemID, v.Num, 0)
	}
	var msg S2C_OnHookFast
	msg.Cid = "onhookfast"
	msg.GetItem = getItems
	msg.CostItem = costItem
	msg.FastTimes = self.Sql_OnHook.FastTimes
	msg.GetPrivilegeItems = getPrivilegeItems
	msg.GetMonthItems = getMonthCardItems
	msg.GetActivityItems = getActivityItems
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	self.player.HandleTask(TASK_TYPE_GET_HOOK_FAST_AWARD, 1, 0, 0)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ONHOOK_FAST, self.Sql_OnHook.HangUp, 0, 0, "使用快速挂机", 0, 0, self.player)
}

func (self *ModOnHook) CalItem(itemId int, itemNum int) map[int]*Item {
	rewards := make(map[int]*Item, 0)

	if itemNum <= 0 {
		return rewards
	}

	config := GetCsvMgr().ItemMap[itemId]
	if config == nil {
		return rewards
	}

	configOnHook := GetCsvMgr().HangUpConfig[self.Sql_OnHook.HangUp]
	if configOnHook == nil {
		return rewards
	}

	if config.ItemSubType == 1 || config.ItemSubType == 2 {
		switch config.CompoundId {
		case ITEM_GOLD:
			rate := PER_BIT
			if config.ItemSubType == ONHOOK_VIP_RATE {
				vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
				if vipcsv != nil {
					rate += vipcsv.HangupGold
				}
			}
			if configOnHook.GoldYield > 0 {
				AddItemMapHelper3(rewards, ITEM_GOLD, itemNum*configOnHook.GoldYield*config.Special*rate/(60*PER_BIT))
			}
		case ITEM_HERO_EXP:
			rate := PER_BIT
			if config.ItemSubType == ONHOOK_VIP_RATE {
				vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
				if vipcsv != nil {
					rate += vipcsv.HangupHeroExp
				}
			}
			if configOnHook.HeroExpYield > 0 {
				AddItemMapHelper3(rewards, ITEM_HERO_EXP, itemNum*configOnHook.HeroExpYield*config.Special*rate/(60*PER_BIT))
			}
		case ITEM_POWDER:
			if configOnHook.PowderYield > 0 {
				AddItemMapHelper3(rewards, ITEM_POWDER, itemNum*configOnHook.PowderYield*config.Special/(60*60))
			}
		case 1:
			rate := PER_BIT
			if config.ItemSubType == ONHOOK_VIP_RATE {
				vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
				if vipcsv != nil {
					rate += vipcsv.HangupGold
				}
			}
			if configOnHook.GoldYield > 0 {
				AddItemMapHelper3(rewards, ITEM_GOLD, itemNum*configOnHook.GoldYield*config.Special*rate/(60*PER_BIT))
			}

			rate = PER_BIT
			if config.ItemSubType == ONHOOK_VIP_RATE {
				vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
				if vipcsv != nil {
					rate += vipcsv.HangupHeroExp
				}
			}
			if configOnHook.HeroExpYield > 0 {
				AddItemMapHelper3(rewards, ITEM_HERO_EXP, itemNum*configOnHook.HeroExpYield*config.Special*rate/(60*PER_BIT))
			}

			if configOnHook.PowderYield > 0 {
				AddItemMapHelper3(rewards, ITEM_POWDER, itemNum*configOnHook.PowderYield*config.Special/(60*60))
			}
		}
	}
	return rewards
}

func (self *ModOnHook) AddPrivilegeExtItems(itemId []int, itemNum []int) {
	if self.Sql_OnHook.calTimeExtItems == nil {
		self.Sql_OnHook.calTimeExtItems = make(map[int]*Item, 0)
	}
	AddItemMapHelper(self.Sql_OnHook.calTimeExtItems, itemId, itemNum)
}

func (self *ModOnHook) CheckPass() {
	lastpass := self.player.GetModule("pass").(*ModPass).GetLastPass()
	if lastpass != nil {
		config := GetCsvMgr().LevelConfigMap[lastpass.Id]
		if config != nil {
			configNext := GetCsvMgr().LevelConfigMap[config.NextLevel]
			if configNext != nil && config.LevelIndex == configNext.LevelIndex {
				self.SetStage(configNext.LevelId)
				self.player.HandleTask(TASK_TYPE_FINISH_CHAPTER, config.LevelIndex-1, 0, 0)
			} else {
				self.SetStage(lastpass.Id)
				self.player.HandleTask(TASK_TYPE_FINISH_CHAPTER, config.LevelIndex, 0, 0)
			}
		}
	}
}