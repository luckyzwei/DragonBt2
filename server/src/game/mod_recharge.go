package game

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

const (
	MAX_BASE_RECHARGE = 11 //! 基础11种充值
	MAX_BOX_RECHAEGE  = 25 //! 礼包25种
)

const (
	RECHARGE_FIRST  = 1
	RECHARGE_SECOND = 2
	RECHARGE_THIRD  = 3
)

const (
	WARORDER_1   = 1 //高级皇家犒赏令
	WARORDER_2   = 2 //高级勇者犒赏令
	WARORDER_END = 3
)

const (
	WARORDERLIMIT_1   = 1 //主线战令
	WARORDERLIMIT_2   = 2 //爬塔战令
	WARORDERLIMIT_3   = 3 //钻石累消令
	WARORDERLIMIT_END = 4
)

const (
	FIRST_RECHARGE_TYPE_OFFSET = 900100
)

//! 充值数据库
type San_UserRecharge struct {
	Uid            int64  `json:"uid"`
	Money          int    `json:"money"`
	Getallgem      int    `json:"getallgem"`
	Type1          int    `json:"type1"`
	Type2          int    `json:"type2"`
	Type3          int    `json:"type3"`
	Type4          int    `json:"type4"`
	Type5          int    `json:"type5"`
	Type6          int    `json:"type6"`
	Record         string `json:"record"`
	Firstaward     int    `json:"firstaward"`
	MoneyDay       int    `json:"moneyday"`
	MoneyWeek      int    `json:"moneyweek"`
	MonthCount1    int    `json:"monthcount1"`
	MonthCount2    int    `json:"monthcount2"`
	MonthCount3    int    `json:"monthcount3"`
	VipBox         int64  `json:"vipbox"`
	FundType       int64  `json:"fundtype"`
	FundGet        string `json:"fundget"`
	BaseCounts     string `json:"basecounts"`     //! 基础充值次数
	BoxCounts      string `json:"boxcounts"`      //! 礼包购买次数
	VipDailyReward int64  `json:"vipdailyreward"` //! vip每日领取
	VipWeekBuy     int64  `json:"vipweekbuy"`     //！ VIP 每周礼包购买
	WarOrder       string `json:"warorder"`       //！ 犒赏令
	WarOrderLimit  string `json:"warorderlimit"`  //！ 限时战令

	record        []JS_RechargeRecord
	baseCounts    []int
	boxCounts     []int
	fundget       []int64
	warOrder      []JS_WarOrder
	warOrderLimit []JS_WarOrderLimit

	DataUpdate
}

//!
type JS_RechargeRecord struct {
	Type      int   `json:"type"`
	Money     int   `json:"money"`
	Addgem    int   `json:"addgem"`
	ExtraGem  int   `json:"extragem"`
	BeforeGem int   `json:"befroegem"`
	AfterGem  int   `json:"aftergem"`
	Time      int64 `json:"time"`
	OrderId   int   `json:"orderid"`
	Isok      int   `json:"isok"`
}

type JS_WarOrder struct {
	Type         int               `json:"type"`
	StartTime    int64             `json:"starttime"`
	EndTime      int64             `json:"endtime"`
	BuyState     int               `json:"buystate"` //购买状态
	WarOrderTask []JS_WarOrderTask `json:"warordertask"`
	Plan         int               `json:"plan"` //当前任务进度
}

type JS_WarOrderLimit struct {
	Type         int               `json:"type"`
	StartTime    int64             `json:"starttime"`
	EndTime      int64             `json:"endtime"`
	BuyState     int               `json:"buystate"` //购买状态
	WarOrderTask []JS_WarOrderTask `json:"warordertask"`
	N3           int               `json:"n3"`
	N4           int               `json:"n4"`
}

type JS_FundInfo struct {
	Type         int               `json:"type"`
	StartTime    int64             `json:"starttime"`
	EndTime      int64             `json:"endtime"`
	BuyState     int               `json:"buystate"` //购买状态
	WarOrderTask []JS_WarOrderTask `json:"warordertask"`
	N3           int               `json:"n3"`
	N4           int               `json:"n4"`
}

type JS_WarOrderTask struct {
	Id     int `json:"id"`     // id
	State  int `json:"state"`  // 0未完成 1可领取  2已领取
	BuyGet int `json:"buyget"` //付费领取
	Plan   int `json:"plan"`   //当前任务进度
}

//! 充值
type ModRecharge struct {
	player           *Player
	Sql_UserRecharge San_UserRecharge  //! 数据库结构
	chg              []JS_WarOrderTask //用于发送限时战令同步
}

func (self *ModRecharge) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_userrecharge` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_UserRecharge, "san_userrecharge", self.player.ID)

	if self.Sql_UserRecharge.Uid <= 0 {
		self.Sql_UserRecharge.Uid = self.player.ID
		self.Sql_UserRecharge.record = make([]JS_RechargeRecord, 0)
		self.Sql_UserRecharge.fundget = make([]int64, 5)
		self.Sql_UserRecharge.baseCounts = make([]int, MAX_BASE_RECHARGE)
		self.Sql_UserRecharge.boxCounts = make([]int, MAX_BOX_RECHAEGE)
		self.Sql_UserRecharge.FundType = 0
		self.Encode()
		InsertTable("san_userrecharge", &self.Sql_UserRecharge, 0, true)
	} else {
		self.Decode()
		if self.Sql_UserRecharge.baseCounts == nil {
			self.Sql_UserRecharge.baseCounts = make([]int, MAX_BASE_RECHARGE)
		}
		if self.Sql_UserRecharge.boxCounts == nil {
			self.Sql_UserRecharge.boxCounts = make([]int, MAX_BOX_RECHAEGE)
		}
	}

	self.Sql_UserRecharge.Init("san_userrecharge", &self.Sql_UserRecharge, true)
}

func (self *ModRecharge) CheckWarOrder() bool {
	rel := false
	size := len(self.Sql_UserRecharge.warOrder)
	if size < WARORDER_END-1 {
		for i := size; i < WARORDER_END-1; i++ {
			var warorder JS_WarOrder
			warorder.Type = i + 1
			self.Sql_UserRecharge.warOrder = append(self.Sql_UserRecharge.warOrder, warorder)
		}
		rel = true
	}

	for i := 0; i < len(self.Sql_UserRecharge.warOrder); i++ {
		flag := self.Sql_UserRecharge.warOrder[i].Check(self.player)
		rel = rel || flag
	}
	return rel
}

// 检查活跃数据
func (self *ModRecharge) CheckWarOrderLimit() bool {
	rel := false
	size := len(self.Sql_UserRecharge.warOrderLimit)
	if size < WARORDERLIMIT_END-1 {
		for i := size; i < WARORDERLIMIT_END-1; i++ {
			var warOrderLimit JS_WarOrderLimit
			warOrderLimit.Type = i + 1
			self.Sql_UserRecharge.warOrderLimit = append(self.Sql_UserRecharge.warOrderLimit, warOrderLimit)
		}
		rel = true
	}

	for i := 0; i < len(self.Sql_UserRecharge.warOrderLimit); i++ {
		flag := self.Sql_UserRecharge.warOrderLimit[i].Check(self.player)
		rel = rel || flag
	}
	return rel
}

func (self *ModRecharge) CalWarOrder(nType int, add int) {
	if nType <= 0 || nType > len(self.Sql_UserRecharge.warOrder) {
		return
	}
	realIndex := nType - 1
	if TimeServer().Unix() < self.Sql_UserRecharge.warOrder[realIndex].StartTime ||
		TimeServer().Unix() > self.Sql_UserRecharge.warOrder[realIndex].EndTime {
		return
	}
	self.Sql_UserRecharge.warOrder[realIndex].Plan += add
	for i := 0; i < len(self.Sql_UserRecharge.warOrder[realIndex].WarOrderTask); i++ {
		if self.Sql_UserRecharge.warOrder[realIndex].WarOrderTask[i].State != CANTFINISH {
			continue
		}
		config := GetCsvMgr().WarOrderConfig[self.Sql_UserRecharge.warOrder[realIndex].WarOrderTask[i].Id]
		if config == nil {
			continue
		}
		if config.NeedPoint < self.Sql_UserRecharge.warOrder[realIndex].Plan {
			self.Sql_UserRecharge.warOrder[realIndex].WarOrderTask[i].State = CANTAKE
		}
	}

	var msg S2C_WarOrderTaskUpdate
	msg.Cid = "warordertaskupdate"
	msg.WarOrder = self.Sql_UserRecharge.warOrder[realIndex]
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}

//! 更新战令数据
func (self *ModRecharge) CalWarOrderLimit(nType int) {

	if nType <= 0 || nType > len(self.Sql_UserRecharge.warOrderLimit) {
		return
	}

	//! 判断战令是否过期
	realIndex := nType - 1
	if TimeServer().Unix() < self.Sql_UserRecharge.warOrderLimit[realIndex].StartTime ||
		TimeServer().Unix() > self.Sql_UserRecharge.warOrderLimit[realIndex].EndTime {
		return
	}
	//! 更新任务
	if nType == WARORDERLIMIT_3{
		for i := 0; i < len(self.Sql_UserRecharge.warOrderLimit[realIndex].WarOrderTask); i++ {		// 更新各个奖励的领取状态
			if self.Sql_UserRecharge.warOrderLimit[realIndex].WarOrderTask[i].State != CANTFINISH {
				continue
			}
			config := GetCsvMgr().WarOrderLimitConfig[self.Sql_UserRecharge.warOrderLimit[realIndex].WarOrderTask[i].Id]
			if config == nil {
				continue
			}
			//if config.Ns[0] <= self.player.GetModule("activity").(*ModActivity).Sql_Activity.DiamondConsum {
			//	self.Sql_UserRecharge.warOrderLimit[realIndex].WarOrderTask[i].State = CANTAKE
			//}
		}
	}

	//! 发送消息同步
	var msg S2C_WarOrderLimitInfo
	msg.Cid = "warorderlimitinfo"
	msg.WarOrderLimit = self.Sql_UserRecharge.warOrderLimit
	//msg.WarOrderLimit[2].Plan = self.player.GetModule("activity").(*ModActivity).Sql_Activity.DiamondConsum
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}


func (self *ModRecharge) GetLastRechargeTime() int64 {
	nLen := len(self.Sql_UserRecharge.record)
	if nLen > 0 {
		return self.Sql_UserRecharge.record[nLen-1].Time
	}
	return 0
}

func (self *ModRecharge) WarOrderIsOpen(nType int) bool {
	if nType <= 0 || nType > len(self.Sql_UserRecharge.warOrder) {
		return false
	}
	realIndex := nType - 1
	if TimeServer().Unix() < self.Sql_UserRecharge.warOrder[realIndex].StartTime ||
		TimeServer().Unix() > self.Sql_UserRecharge.warOrder[realIndex].EndTime {
		return false
	}
	return true
}

func (self *JS_WarOrder) Check(player *Player) bool {
	now := TimeServer().Unix()
	if now < self.EndTime {
		return false
	}

	switch self.Type {
	case 1:
		self.ResetWarOrder(player, ACT_WARORDER_1)
	case 2:
		self.ResetWarOrder(player, ACT_WARORDER_2)
	}
	return true
}


// 检查主线战令、钻石累消令
func (self *JS_WarOrderLimit) Check(player *Player) bool {
	now := TimeServer().Unix()
	//if len(self.WarOrderTask) > 0 && now < self.EndTime {
	//	return false
	//}
	var sum int
	for _,v := range GetCsvMgr().WarOrderLimitConfig{	//计算配置文件中，该任务一共多少项
		if v.Type ==self.Type{
			sum += 1
		}
	}
	if len(self.WarOrderTask) == sum && now < self.EndTime {
		if self.Type == 3 && self.N3 != GetActivityMgr().getActN3(ACT_WARORDERLIMIT_3){
		}else{
			return false
		}
	}

	switch self.Type {
	case 1:
		self.ResetWarOrder(player, ACT_WARORDERLIMIT_1)
		passId := player.GetModule("onhook").(*ModOnHook).GetStage()
		config := GetCsvMgr().LevelConfigMap[passId]
		if config != nil {
			configNext := GetCsvMgr().LevelConfigMap[config.NextLevel]
			if configNext != nil && config.LevelIndex != configNext.LevelIndex {
				player.GetModule("recharge").(*ModRecharge).HandleTask(TASK_TYPE_FINISH_CHAPTER, config.LevelIndex, 0, 0)
			} else {
				player.GetModule("recharge").(*ModRecharge).HandleTask(TASK_TYPE_FINISH_CHAPTER, config.LevelIndex-1, 0, 0)
			}
		}
	case 2:
		self.ResetWarOrder(player, ACT_WARORDERLIMIT_2)
		level := player.GetModule("tower").(*ModTower).GetMainLevel()
		player.GetModule("recharge").(*ModRecharge).HandleTask(TASK_TYPE_WOTER_LEVEL, level, 0, 0)
	case 3:
		self.ResetWarOrder(player, ACT_WARORDERLIMIT_3)
		//level := player.GetModule("tower").(*ModTower).GetMainLevel()
		//player.GetModule("recharge").(*ModRecharge).HandleTask(TASK_TYPE_DIAMOND_CONSUME, 0, 0, 0)
	}
	return true
}

func (self *JS_WarOrder) ResetWarOrder(player *Player, activityId int) {
	now := TimeServer().Unix()

	rTime, _ := time.ParseInLocation(DATEFORMAT, player.Sql_UserBase.Regtime, time.Local)
	activity := GetActivityMgr().GetActivity(activityId)
	if activity == nil {
		return
	}
	mailId := 0
	switch activityId {
	case ACT_WARORDER_1:
		startDay := HF_Atoi(activity.info.Start)
		if startDay >= 0 {
			return
		}
		mailId = 1
	case ACT_WARORDER_2:
		lastpass := player.GetModule("pass").(*ModPass).GetLastPass()
		passId := ONHOOK_INIT_LEVEL
		if lastpass != nil {
			passId = lastpass.Id
		}
		if !GetCsvMgr().IsLevelAndPassOpenNew(player.Sql_UserBase.Level, passId, OPEN_LEVEL_WARORDER_2) {
			return
		}
		mailId = 2
	}
	//重置之前需要发送上次没领完的奖励
	getItem := make(map[int]*Item)
	for i := 0; i < len(self.WarOrderTask); i++ {
		config := GetCsvMgr().GetWarOrderConfig(self.WarOrderTask[i].Id)
		if config == nil {
			continue
		}
		if self.WarOrderTask[i].State == CANTAKE {
			AddItemMapHelper(getItem, config.FreeAward, config.FreeNum)
			self.WarOrderTask[i].State = TAKEN
		}
		if self.WarOrderTask[i].State != CANTFINISH {
			if self.BuyState == LOGIC_TRUE && self.WarOrderTask[i].BuyGet == LOGIC_FALSE {
				AddItemMapHelper(getItem, config.GoldAward, config.GoldNum)
				self.WarOrderTask[i].BuyGet = LOGIC_TRUE
			}
		}
	}
	if len(getItem) > 0 {
		mailConfig, ok := GetCsvMgr().MailConfig[mailId]
		pMail := player.GetModule("mail").(*ModMail)
		if ok && pMail != nil {
			// 获得奖励
			var mailItems []PassItem
			for _, v := range getItem {
				mailItems = append(mailItems, PassItem{v.ItemId, v.ItemNum})
			}
			// 发送邮件
			pMail.AddMailWithItems(MAIL_CAN_ALL_GET, mailConfig.Mailtitle, mailConfig.Mailtxt, mailItems)
		}
	}

	switch activityId {
	case ACT_WARORDER_1:
		//开始重置
		startDay := HF_Atoi(activity.info.Start)
		//计算经过了几轮
		firstTime := HF_CalPlayerCreateTime(rTime.Unix(), -(startDay + 1))
		times := (now - firstTime) / int64(activity.info.Continued+activity.info.CD)
		self.StartTime = firstTime + int64(activity.info.Continued+activity.info.CD)*times
		self.EndTime = self.StartTime + int64(activity.info.Continued)
	case ACT_WARORDER_2:
		self.StartTime = self.EndTime + int64(activity.info.CD)
		self.EndTime = self.StartTime + int64(activity.info.Continued)
		if self.EndTime < now {
			self.StartTime = time.Date(TimeServer().Year(), TimeServer().Month(), TimeServer().Day(), 5, 0, 0, 0, TimeServer().Location()).Unix()
			self.EndTime = self.StartTime + int64(activity.info.Continued)
		}
	}

	self.BuyState = LOGIC_FALSE
	self.Plan = 0

	//初始化任务
	self.WarOrderTask = make([]JS_WarOrderTask, 0)
	for _, v := range GetCsvMgr().WarOrderConfig {
		if v.Type != self.Type {
			continue
		}
		var task JS_WarOrderTask
		task.Id = v.Id
		if v.NeedPoint <= self.Plan {
			task.State = CANTAKE
		} else {
			task.State = CANTFINISH
		}
		task.BuyGet = LOGIC_FALSE
		self.WarOrderTask = append(self.WarOrderTask, task)
	}
}

func (self *JS_WarOrderLimit) ResetWarOrder(player *Player, activityId int) {

	activity := GetActivityMgr().GetActivity(activityId)
	if activity == nil {
		return
	}

	lastpass := player.GetModule("pass").(*ModPass).GetLastPass()
	passId := ONHOOK_INIT_LEVEL
	if lastpass != nil {
		passId = lastpass.Id
	}

	//mailId := 0
	switch activityId {
	case ACT_WARORDERLIMIT_1:
		if !GetCsvMgr().IsLevelAndPassOpenNew(player.Sql_UserBase.Level, passId, OPEN_LEVEL_WARORDERLIMIT_1) {
			return
		}
		//mailId = 1
	case ACT_WARORDERLIMIT_2:
		if !GetCsvMgr().IsLevelAndPassOpenNew(player.Sql_UserBase.Level, passId, OPEN_LEVEL_WARORDERLIMIT_2) {
			return
		}
		//mailId = 2
	}
	//重置之前需要发送上次没领完的奖励
	/*
		getItem := make(map[int]*Item)
		for i := 0; i < len(self.WarOrderTask); i++ {
			config := GetCsvMgr().GetWarOrderConfig(self.WarOrderTask[i].Id)
			if config == nil {
				continue
			}
			if self.WarOrderTask[i].State == CANTAKE {
				AddItemMapHelper(getItem, config.FreeAward, config.FreeNum)
				self.WarOrderTask[i].State = TAKEN
			}
			if self.WarOrderTask[i].State != CANTFINISH {
				if self.BuyState == LOGIC_TRUE && self.WarOrderTask[i].BuyGet == LOGIC_FALSE {
					AddItemMapHelper(getItem, config.GoldAward, config.GoldNum)
					self.WarOrderTask[i].BuyGet = LOGIC_TRUE
				}
			}
		}
		if len(getItem) > 0 {
			mailConfig, ok := GetCsvMgr().MailConfig[mailId]
			pMail := player.GetModule("mail").(*ModMail)
			if ok && pMail != nil {
				// 获得奖励
				var mailItems []PassItem
				for _, v := range getItem {
					mailItems = append(mailItems, PassItem{v.ItemId, v.ItemNum})
				}
				// 发送邮件
				pMail.AddMailWithItems(MAIL_CAN_ALL_GET, mailConfig.Mailtitle, mailConfig.Mailtxt, mailItems)
			}
		}
	*/

	ver := GetActivityMgr().getActN3(activityId)
	group := GetActivityMgr().getActN4(activityId)

	if self.N3 == ver && self.N4 == group {
		return
	}

	switch activityId {
	case ACT_WARORDERLIMIT_1:
		//开始重置
		startDay := HF_Atoi(activity.info.Start)
		if startDay > 0 {
			self.StartTime = GetServer().GetOpenServer() + int64(startDay-1)*86400
		} else {
			rTime, _ := time.ParseInLocation(DATEFORMAT, player.Sql_UserBase.Regtime, time.Local)
			correctTime := HF_CalPlayerCreateTime(rTime.Unix(), 0)
			self.StartTime = correctTime + int64(-(startDay + 1)*DAY_SECS)
		}
		self.EndTime = self.StartTime + int64(activity.info.Continued)
		self.N3 = ver
		self.N4 = group
	case ACT_WARORDERLIMIT_2:
		//开始重置
		startDay := HF_Atoi(activity.info.Start)
		if startDay > 0 {
			self.StartTime = GetServer().GetOpenServer() + int64(startDay-1)*86400
		} else {
			rTime, _ := time.ParseInLocation(DATEFORMAT, player.Sql_UserBase.Regtime, time.Local)
			correctTime := HF_CalPlayerCreateTime(rTime.Unix(), 0)
			self.StartTime = correctTime + int64(-(startDay + 1)*DAY_SECS)
		}
		self.EndTime = self.StartTime + int64(activity.info.Continued)
		self.N3 = ver
		self.N4 = group
	case ACT_WARORDERLIMIT_3:
		//开始重置
		startDay := HF_Atoi(activity.info.Start)
		if startDay > 0 {
			self.StartTime = GetServer().GetOpenServer() + int64(startDay-1)*DAY_SECS
		} else if startDay < 0 {
			rTime, _ := time.ParseInLocation(DATEFORMAT, player.Sql_UserBase.Regtime, time.Local)
			correctTime := HF_CalPlayerCreateTime(rTime.Unix(), 0)
			self.StartTime = correctTime + int64(-(startDay+1)*DAY_SECS)
		} else {
			if activity.info.Start != "" && activity.info.Start != "0" {
				// 否则按照实际时间算
				t, err := time.ParseInLocation(DATEFORMAT, activity.info.Start, time.Local)
				if err == nil {
					self.StartTime = t.Unix()
				}
			} else if activity.info.Start == "" || activity.info.Start == "0" {
				self.StartTime = TimeServer().Unix()
			}
		}
		if activity.info.Continued == 0 {
			self.EndTime = self.StartTime + DAY_SECS*730
		} else {
			self.EndTime = self.StartTime + int64(activity.info.Continued)
		}
		self.N3 = ver	// 更新活动期数
		self.N4 = group
	}

	self.BuyState = LOGIC_FALSE

	//初始化任务
	self.WarOrderTask = make([]JS_WarOrderTask, 0)
	for _, v := range GetCsvMgr().WarOrderLimitConfig {
		if v.Type != self.Type || v.Group != self.N4 {
			continue
		}
		var task JS_WarOrderTask
		task.Id = v.Id
		task.State = CANTFINISH
		task.BuyGet = LOGIC_FALSE
		self.WarOrderTask = append(self.WarOrderTask, task)
	}
}

func (self *ModRecharge) HandleTask(tasktype, n2, n3, n4 int) {
	now := TimeServer().Unix()

	// 战令相关
	for i := 0; i < len(self.Sql_UserRecharge.warOrderLimit); i++ {
		for j := 0; j < len(self.Sql_UserRecharge.warOrderLimit[i].WarOrderTask); j++ {
			warOrder := &self.Sql_UserRecharge.warOrderLimit[i].WarOrderTask[j]
			config := GetCsvMgr().WarOrderLimitConfig[warOrder.Id]
			if config == nil {
				continue
			}
			if config.TaskTypes != tasktype {
				continue
			}
			if config.Group != self.Sql_UserRecharge.warOrderLimit[i].N4 {
				continue
			}

			if self.Sql_UserRecharge.warOrderLimit[i].BuyState == LOGIC_FALSE &&
				(now < self.Sql_UserRecharge.warOrderLimit[i].StartTime || now >= self.Sql_UserRecharge.warOrderLimit[i].EndTime) {
				continue
			}

			var tasknode TaskNode
			tasknode.Tasktypes = config.TaskTypes
			tasknode.N1 = config.Ns[0]
			tasknode.N2 = config.Ns[1]
			tasknode.N3 = config.Ns[2]
			tasknode.N4 = config.Ns[3]
			plan, add := DoTask(&tasknode, self.player, n2, n3, n4)
			if plan == 0 {
				continue
			}

			chg := false
			if add {
				warOrder.Plan += plan
				chg = true
			} else {
				if plan > warOrder.Plan {
					warOrder.Plan = plan
					chg = true
				}
			}

			if warOrder.Plan >= config.Ns[0] {
				if warOrder.State == CANTFINISH { // 完成了没有进行领取
					warOrder.State = CANTAKE
				}
			}

			if chg {
				self.chg = append(self.chg, *warOrder)
			}
		}
	}
}

// 发送礼包更新信息
func (self *ModRecharge) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_WarOrderLimit
	msg.Cid = "warorderlimitupdate"
	msg.Info = self.chg
	smsg, _ := json.Marshal(&msg)
	self.chg = make([]JS_WarOrderTask, 0)
	self.player.SendMsg(msg.Cid, smsg)
}

func (self *ModRecharge) OnGetOtherData() {

}

func (self *ModRecharge) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "recharge":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		// 检查重置
		flag := self.player.CheckRefresh()
		if flag {
			return true
		}
		var s2c_msg C2S_Recharge
		json.Unmarshal(body, &s2c_msg)
		if s2c_msg.Type > 1000 {
			self.RechargeEx(s2c_msg.Type, int(TimeServer().Unix()), 1, s2c_msg.Money, 0)
		} else {
			self.RechargeEx(1001000+s2c_msg.Type, int(TimeServer().Unix()), 1, s2c_msg.Money, 0)
		}

		return true
	case "getfirstaward":
		var c2s_msg C2S_GetFirsetAward
		json.Unmarshal(body, &c2s_msg)
		self.GetFirsetAward(c2s_msg.Type)
		return true
	case "buyvipbox":
		var msg C2S_BuyVipBox
		json.Unmarshal(body, &msg)
		self.BuyVipBox(msg.Index)
		return true
	case "buyfund":
		var msg C2S_ButFund
		json.Unmarshal(body, &msg)
		self.BuyFundType(msg.FundType)
		return true
	case "getfundaward":
		var msg C2S_GetFundAward
		json.Unmarshal(body, &msg)
		self.GetFundAward(msg.Fundid, msg.Pageid)
		return true
	case "getfundtotal":
		self.GetFundTotal()
		return true
	case "getvipdailyreward":
		var c2s_msg C2S_GetVipDailyReward
		json.Unmarshal(body, &c2s_msg)
		self.GetVipDaily(c2s_msg.VipLevel)
		return true
	case "buyvipweek":
		var c2s_msg C2S_BuyVipWeek
		json.Unmarshal(body, &c2s_msg)
		self.BuyVipWeek(c2s_msg.VipLevel)
		return true
	case "getwarorderreward":
		var c2s_msg C2S_GetWarOrderReward
		json.Unmarshal(body, &c2s_msg)
		self.GetWarOrderReward(&c2s_msg)
		return true
	case "getwarorderlimitreward":
		var c2s_msg C2S_GetWarOrderLimitReward
		json.Unmarshal(body, &c2s_msg)
		self.GetWarOrderLimitReward(&c2s_msg)
		return true
	case "warorderbuy":
		var c2s_msg C2S_WarOrderBuy
		json.Unmarshal(body, &c2s_msg)
		self.WarOrderBuy(&c2s_msg)
		return true
	}

	return false
}

func (self *ModRecharge) OnSave(sql bool) {
	self.Encode()
	self.Sql_UserRecharge.Update(sql)
}

func (self *ModRecharge) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_UserRecharge.Record), &self.Sql_UserRecharge.record)
	json.Unmarshal([]byte(self.Sql_UserRecharge.FundGet), &self.Sql_UserRecharge.fundget)
	json.Unmarshal([]byte(self.Sql_UserRecharge.BaseCounts), &self.Sql_UserRecharge.baseCounts)
	json.Unmarshal([]byte(self.Sql_UserRecharge.BoxCounts), &self.Sql_UserRecharge.boxCounts)
	json.Unmarshal([]byte(self.Sql_UserRecharge.WarOrder), &self.Sql_UserRecharge.warOrder)
	json.Unmarshal([]byte(self.Sql_UserRecharge.WarOrderLimit), &self.Sql_UserRecharge.warOrderLimit)
}

func (self *ModRecharge) Encode() { //! 将data数据写入数据库
	s, _ := json.Marshal(&self.Sql_UserRecharge.record)
	self.Sql_UserRecharge.Record = string(s)
	s, _ = json.Marshal(&self.Sql_UserRecharge.fundget)
	self.Sql_UserRecharge.FundGet = string(s)

	self.Sql_UserRecharge.BaseCounts = HF_JtoA(self.Sql_UserRecharge.baseCounts)
	self.Sql_UserRecharge.BoxCounts = HF_JtoA(self.Sql_UserRecharge.boxCounts)
	self.Sql_UserRecharge.WarOrder = HF_JtoA(self.Sql_UserRecharge.warOrder)
	self.Sql_UserRecharge.WarOrderLimit = HF_JtoA(self.Sql_UserRecharge.warOrderLimit)
}

func (self *ModRecharge) GetFundTotal() {
	var msg S2C_GetFundTotal
	msg.Cid = "getfundtotal"
	msg.FundTotal = GetHeroSupportMgr().GetJJ()

	self.player.SendMsg("getfundtotal", HF_JtoB(&msg))
}

func (self *ModRecharge) GetMaxFundType() int {
	var maxfundtype int = 0
	for i := 1; i <= 3; i++ {
		var base int64 = 1 << uint(i-1)
		if self.Sql_UserRecharge.FundType&base != 0 {
			maxfundtype = i
		}
	}
	return maxfundtype
}

//! 购买基金
func (self *ModRecharge) BuyFundType(fundtype int) {
	if fundtype <= 0 {
		self.player.SendErrInfo("err", "基金类型错误")
		return
	}
	//! 判断重复购买
	var base int64 = 1 << uint(fundtype-1)
	if self.Sql_UserRecharge.FundType&base != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_CANT_REPEAT_PURCHASE"))
		return
	}
	csv_fundbuy, ok := GetCsvMgr().Data["Fund_Buy"][fundtype]
	if !ok {
		self.player.SendErrInfo("err", "基金配置错误")
		return
	}
	//! 判断vip等级
	if self.player.Sql_UserBase.Vip < HF_Atoi(csv_fundbuy["buyvip"]) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TASK_LACK_OF_ARISTOCRATIC_RANK"))
		return
	}
	//! 判断钻石是否足够
	if self.player.GetObjectNum(91000002) < HF_Atoi(csv_fundbuy["fee"]) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_DIAMOND_SHORTAGE"))
		return
	}
	//! 修改数据
	self.player.AddObject(91000002, HF_Atoi(csv_fundbuy["fee"])*-1, fundtype, 0, 0, "商会购买")
	if self.Sql_UserRecharge.FundType == 0 {
		GetHeroSupportMgr().AddJJ()
	}
	self.Sql_UserRecharge.FundType += base
	cost := make([]PassItem, 0)
	cost = append(cost, PassItem{91000002, HF_Atoi(csv_fundbuy["fee"])})

	//type_name := "商会"
	//if fundtype == 1 {
	//	type_name = "聚财计划"
	//} else if fundtype == 2 {
	//	type_name = "红颜计划"
	//} else if fundtype == 3 {
	//	type_name = "纳贤计划"
	//}

	GetServer().sendLog_CommerceJoin(self.player, fundtype, "商会购买")

	var msg S2C_BuyFund
	msg.Cid = "buyfund"
	msg.Cost = cost
	msg.FundType = self.Sql_UserRecharge.FundType

	self.player.SendMsg("buyfund", HF_JtoB(&msg))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FUND_BUY, fundtype, 0, 0, "商会购买", 0, 0, self.player)
}

func (self *ModRecharge) OnRefresh() {
	if len(self.Sql_UserRecharge.fundget) > 3 {
		self.Sql_UserRecharge.fundget[3] = 0
	}
	self.Sql_UserRecharge.VipDailyReward = 0
	self.Sql_UserRecharge.MoneyDay = 0

	//周一重置 VIP周礼包
	if TimeServer().Weekday() == 1 {
		self.Sql_UserRecharge.VipWeekBuy = 0
	}

	self.CheckOpen()
	self.CheckOpenLimit()
}

//! 领取基金奖励
func (self *ModRecharge) GetFundAward(fundid int, pageid int) {

	if fundid == 0 {
		//读pageid一键领取  20190723 by zy

		fundconfig, ok := GetCsvMgr().Data["Fund"]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FUND_NO_CONFIG"))
			return
		}
		award := make([]PassItem, 0)
		isHas := false
		for _, config := range fundconfig {
			if pageid != HF_Atoi(config["paging"]) {
				continue
			}

			page := pageid - 1

			index := HF_Atoi(config["fundid"]) % 100
			//! 判断重复领取
			var base int64 = 1 << uint(index-1)
			if self.Sql_UserRecharge.fundget[page]&base != 0 {
				continue
			}
			//! 判断基金类型
			if HF_Atoi(config["fundlevel"]) > 0 {
				if HF_Atoi(config["fundlevel"]) == 4 {
					if self.Sql_UserRecharge.FundType == 0 {
						continue
					}
				} else {
					if self.Sql_UserRecharge.FundType&(1<<uint(HF_Atoi(config["fundlevel"])-1)) == 0 {
						continue
					}
				}
			}

			//! 判断等级
			if HF_Atoi(config["level"]) > 0 {
				if HF_Atoi(config["level"]) > self.player.Sql_UserBase.Level {
					continue
				}
			}
			//! 判断全服购买人数
			if HF_Atoi(config["buy"]) > GetHeroSupportMgr().JJNum {
				continue
			}
			//! 判断vip等级
			if HF_Atoi(config["viplevel"]) > 0 {
				if self.player.Sql_UserBase.Vip < HF_Atoi(config["viplevel"]) {
					continue
				}
			}
			//! 修改数据
			self.Sql_UserRecharge.fundget[page] += base
			isHas = true

			for i := 0; i < 3; i++ {
				var item PassItem
				item.ItemID, item.Num = HF_Atoi(config[fmt.Sprintf("award%d", i+1)]), HF_Atoi(config[fmt.Sprintf("num%d", i+1)])
				if item.ItemID > 0 {
					award = append(award, item)
				}
			}
			for i := 0; i < len(award); i++ {
				self.player.AddObject(award[i].ItemID, award[i].Num, fundid, 0, 0, "商会奖励领取")
			}
		}

		if isHas {
			var msg S2C_GetFundAward
			msg.Cid = "getfundaward"
			msg.Award = award
			msg.Index = pageid
			msg.Value = self.Sql_UserRecharge.fundget[pageid-1]

			self.player.SendMsg("getfundaward", HF_JtoB(&msg))

			//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FUND_AWARD, fundid, 0, 0, "商会奖励领取", 0, 0, self.player)
		}

	} else {
		csv_fund, ok := GetCsvMgr().Data["Fund"][fundid]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FUND_NO_CONFIG")+fmt.Sprintf("%d", fundid))
			return
		}
		page := HF_Atoi(csv_fund["paging"]) - 1
		index := fundid % 100
		//! 判断重复领取
		var base int64 = 1 << uint(index-1)
		if self.Sql_UserRecharge.fundget[page]&base != 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_CANT_REPEAT_PURCHASE"))
			return
		}
		//! 判断基金类型
		if HF_Atoi(csv_fund["fundlevel"]) > 0 {
			if HF_Atoi(csv_fund["fundlevel"]) == 4 {
				if self.Sql_UserRecharge.FundType == 0 {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PLEASE_JOIN_THE_CHAMBER_OF"))
					return
				}
			} else {
				if self.Sql_UserRecharge.FundType&(1<<uint(HF_Atoi(csv_fund["fundlevel"])-1)) == 0 {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PLEASE_JOIN_THE_CHAMBER_OF"))
					return
				}
			}
		}

		//! 判断等级
		if HF_Atoi(csv_fund["level"]) > 0 {
			if HF_Atoi(csv_fund["level"]) > self.player.Sql_UserBase.Level {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_LACK_OF_PROTAGONIST_RANK"))
				return
			}
		}
		//! 判断全服购买人数
		if HF_Atoi(csv_fund["buy"]) > GetHeroSupportMgr().JJNum {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_SHORTAGE_OF_FULL-SERVICE_PURCHASERS"))
			return
		}
		//! 判断vip等级
		if HF_Atoi(csv_fund["viplevel"]) > 0 {
			if self.player.Sql_UserBase.Vip < HF_Atoi(csv_fund["viplevel"]) {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_LACK_OF_ARISTOCRATIC_RANK"))
				return
			}
		}
		//! 修改数据
		self.Sql_UserRecharge.fundget[page] += base

		award := make([]PassItem, 0)
		for i := 0; i < 3; i++ {
			var item PassItem
			item.ItemID, item.Num = HF_Atoi(csv_fund[fmt.Sprintf("award%d", i+1)]), HF_Atoi(csv_fund[fmt.Sprintf("num%d", i+1)])
			if item.ItemID > 0 {
				award = append(award, item)
			}
		}
		for i := 0; i < len(award); i++ {
			self.player.AddObject(award[i].ItemID, award[i].Num, fundid, 0, 0, "商会奖励领取")
		}

		var msg S2C_GetFundAward
		msg.Cid = "getfundaward"
		msg.Award = award
		msg.Index = page + 1
		msg.Value = self.Sql_UserRecharge.fundget[page]

		self.player.SendMsg("getfundaward", HF_JtoB(&msg))

		//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FUND_AWARD, fundid, 0, 0, "商会奖励领取", 0, 0, self.player)
	}
}

//! 购买VIP特权礼包
func (self *ModRecharge) BuyVipBox(index int) {
	LogDebug("buybox", index)
	if index <= 0 {
		self.player.SendErrInfo("err", "索引错误")
		return
	}

	if self.player.Sql_UserBase.Vip < index {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TASK_LACK_OF_ARISTOCRATIC_RANK"))
		return
	}

	var base int64 = 1 << uint(index-1)
	if self.Sql_UserRecharge.VipBox&base != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_CANT_REPEAT_PURCHASE"))
		return
	}

	vipcsv, ok := GetCsvMgr().VipConfigMap[index]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_ARISTOCRACY_DOES_NOT_EXIST"))
		return
	}

	if self.player.GetObjectNum(DEFAULT_GEM) < vipcsv.Nownum {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_DIAMOND_SHORTAGE"))
		return
	}

	self.Sql_UserRecharge.VipBox += base

	award := make([]PassItem, 0)
	for i := 0; i < 6; i++ {
		var item PassItem
		item.ItemID, item.Num = vipcsv.Items[i], vipcsv.Nums[i]
		if item.ItemID != 0 {
			award = append(award, item)
		}
	}
	for i := 0; i < len(award); i++ {
		self.player.AddObject(award[i].ItemID, award[i].Num, self.player.Sql_UserBase.Vip, 0, 0, "领取贵族礼包奖励")
	}
	cost := make([]PassItem, 0)
	var item PassItem
	item.ItemID, item.Num = DEFAULT_GEM, vipcsv.Nownum*-1
	self.player.AddObject(item.ItemID, item.Num, self.player.Sql_UserBase.Vip, 0, 0, "领取贵族礼包奖励")
	cost = append(cost, item)

	var msg S2C_BuyVipBox
	msg.Cid = "buyvipbox"
	msg.Award = award
	msg.Cost = cost
	msg.CurVipBox = self.Sql_UserRecharge.VipBox

	self.player.SendMsg("buyvipbox", HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_VIP_BOX, self.player.Sql_UserBase.Vip, 0, 0, "领取贵族礼包奖励", 0, 0, self.player)

}

//! 后台充值调用
func (self *ModRecharge) BackRecharge(order_id int, bok int, moneytype int) {
	//! 如果充值失败
	if bok == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_FAILURE_OF_ORDER_PROCESSING"))
		return
	}
	//! 检查订单ID是否重复
	for i := 0; i < len(self.Sql_UserRecharge.record); i++ {
		if self.Sql_UserRecharge.record[i].OrderId == order_id {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_DUPLICATE_ORDER_NUMBER_NO_PROCESSING"))
			return
		}
	}

	csv, ok := GetCsvMgr().GetMoney(moneytype)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_RECHARGE_TYPE_ANOMALY"))
		return
	}
	addmoney := HF_Atoi(csv["rmb"])
	addgem := HF_Atoi(csv["diamond"])
	extragem := HF_Atoi(csv["extra"])
	logstr := "充值"
	switch moneytype {
	case 1:
		if self.Sql_UserRecharge.Type1 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type1++
		break
	case 2:
		if self.Sql_UserRecharge.Type2 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type2++
		break
	case 3:
		if self.Sql_UserRecharge.Type3 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type3++
		break
	case 4:
		if self.Sql_UserRecharge.Type4 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type4++
		break
	case 5:
		if self.Sql_UserRecharge.Type5 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type5++
		break
	case 6:
		if self.Sql_UserRecharge.Type6 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type6++
		break
	case 101:
		//! 检测月卡时间
		nowtime := TimeServer().Unix()
		for i := 0; i < len(self.player.GetModule("activity").(*ModActivity).Sql_Activity.month); i++ {
			if self.player.GetModule("activity").(*ModActivity).Sql_Activity.month[i].Id == moneytype {
				//! 过期处理
				allday := int64(self.player.GetModule("activity").(*ModActivity).Sql_Activity.month[i].Day)
				useday := (nowtime - self.player.GetModule("activity").(*ModActivity).Sql_Activity.month[i].StartTime) / 86400
				if allday-useday > 90 {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MONTH_CARD_LEFT"))
					return
				}
			}
		}
		self.Sql_UserRecharge.MonthCount1++
		extragem = 0
		logstr = "小月卡"
		self.ProcMonthCard(moneytype, HF_Atoi(csv["time"]))

		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	case 102:
		//! 检测月卡时间
		nowtime := TimeServer().Unix()
		for i := 0; i < len(self.player.GetModule("activity").(*ModActivity).Sql_Activity.month); i++ {
			if self.player.GetModule("activity").(*ModActivity).Sql_Activity.month[i].Id == moneytype {
				//! 过期处理
				allday := int64(self.player.GetModule("activity").(*ModActivity).Sql_Activity.month[i].Day)
				useday := (nowtime - self.player.GetModule("activity").(*ModActivity).Sql_Activity.month[i].StartTime) / 86400
				if allday-useday > 90 {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MONTH_CARD_LEFT"))
					return
				}
			}
		}
		self.Sql_UserRecharge.MonthCount2++
		extragem = 0
		logstr = "大月卡"
		self.ProcMonthCard(moneytype, HF_Atoi(csv["time"]))
		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	case 103:
		//! 检测永久卡
		if self.Sql_UserRecharge.MonthCount3 > 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			return
		}
		self.Sql_UserRecharge.MonthCount3++
		extragem = 0
		logstr = "永久月卡"
		self.ProcMonthCard(moneytype, HF_Atoi(csv["time"]))
		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	}

	if moneytype >= 200 {
		logstr = "幸运礼包"
	}

	before := self.player.Sql_UserBase.Gem
	//! 加钻石
	self.player.AddGem(addgem+extragem, 24, 0, 0, "充值")

	//! 更新任务-判断条件为增加基础钻石
	if addgem == 0 {
		self.player.HandleTask(RechargeGemTask, addmoney*1000, 0, 0)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, DEFAULT_GEM, addmoney*1000, 24, addgem, logstr,
			self.player.Sql_UserBase.Gem, extragem, self.player)
	} else {
		self.player.HandleTask(RechargeGemTask, addgem, 0, 0)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, DEFAULT_GEM, addgem+extragem, 24, addgem, logstr,
			self.player.Sql_UserBase.Gem, extragem, self.player)
	}
	if addgem != 0 {
		self.Sql_UserRecharge.MoneyDay += addgem
	} else {
		self.Sql_UserRecharge.MoneyDay += self.CalGemFroRMB(addmoney)
	}
	self.player.HandleTask(TASK_TYPE_RECHARGE_ONCE, moneytype, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_SINGLE, addmoney, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_MONEY_DAILY, self.Sql_UserRecharge.MoneyDay, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_EQUAL_SINGLE, addmoney, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_ALL, addmoney, 0, 0)
	//! 添加充值记录
	var record JS_RechargeRecord
	record.Money = addmoney
	record.Type = moneytype
	record.Addgem = addgem
	record.ExtraGem = extragem
	record.BeforeGem = before
	record.AfterGem = self.player.Sql_UserBase.Gem
	record.Time = TimeServer().Unix()
	record.OrderId = order_id
	record.Isok = bok

	self.Sql_UserRecharge.record = append(self.Sql_UserRecharge.record, record)
	self.Sql_UserRecharge.Money += addmoney
	self.Sql_UserRecharge.Getallgem += addgem

	//! 付费元宝
	self.player.Sql_UserBase.PayGem += addmoney * 1000

	self.Encode()
	//! 处理vip经验
	self.player.AddVipExp(HF_Atoi(csv["vipexp"]), 0, 0, "充值")

	self.SendInfo("recharge")

}

//! 充值
func (self *ModRecharge) Recharge(moneytype int, order_id int, bok int, money int) int {
	//! 如果充值失败
	if bok == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_FAILURE_OF_ORDER_PROCESSING"))
		return -1
	}
	//! 检查订单ID是否重复
	for i := 0; i < len(self.Sql_UserRecharge.record); i++ {
		if self.Sql_UserRecharge.record[i].OrderId == order_id {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_DUPLICATE_ORDER_NUMBER_NO_PROCESSING"))
			return -1
		}
	}

	csv, ok := GetCsvMgr().GetMoney(moneytype)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_RECHARGE_TYPE_ANOMALY"))
		return -2
	}
	addmoney := HF_Atoi(csv["rmb"])
	addgem := HF_Atoi(csv["diamond"])
	extragem := HF_Atoi(csv["extra"])
	if money != addmoney*100 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_RECHARGE_TYPE_ANOMALY"))
		return -3
	}
	logstr := "充值"
	switch moneytype {
	case 1:
		if self.Sql_UserRecharge.Type1 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type1++
		break
	case 2:
		if self.Sql_UserRecharge.Type2 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type2++
		break
	case 3:
		if self.Sql_UserRecharge.Type3 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type3++
		break
	case 4:
		if self.Sql_UserRecharge.Type4 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type4++
		break
	case 5:
		if self.Sql_UserRecharge.Type5 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type5++
		break
	case 6:
		if self.Sql_UserRecharge.Type6 >= HF_Atoi(csv["time"]) {
			extragem = 0
		}
		self.Sql_UserRecharge.Type6++
		break
	case 101:
		//! 检测月卡时间
		//ok, err := self.CheckMonthCard2(moneytype)
		//if !ok {
		//	self.player.SendErrInfo("err", err)
		//	return 2
		//}
		self.Sql_UserRecharge.MonthCount1++
		extragem = 0
		logstr = "小月卡"
		self.ProcMonthCard(moneytype, HF_Atoi(csv["time"]))

		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	case 102:
		//! 检测月卡时间
		//ok, err := self.CheckMonthCard2(moneytype)
		//if !ok {
		//	self.player.SendErrInfo("err", err)
		//	return 2
		//}

		self.Sql_UserRecharge.MonthCount2++
		extragem = 0
		logstr = "大月卡"
		self.ProcMonthCard(moneytype, HF_Atoi(csv["time"]))
		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	case 103:
		//! 检测永久卡
		if self.Sql_UserRecharge.MonthCount3 > 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			return 2
		}
		self.Sql_UserRecharge.MonthCount3++
		extragem = 0
		logstr = "永久月卡"
		self.ProcMonthCard(moneytype, HF_Atoi(csv["time"]))
		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	case 104:
		//! 检查是否有104了, 已经进行了首冲
		if self.hasMonthCard(104) {
			return 3
		}

		self.Sql_UserRecharge.MonthCount1++
		extragem = 0
		logstr = "小月卡首冲"
		self.ProcMonthCard(moneytype, HF_Atoi(csv["time"]))
		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	case 105:
		//! 检查是否有105了, 已经进行了首冲
		if self.hasMonthCard(105) {
			return 3
		}

		self.Sql_UserRecharge.MonthCount2++
		extragem = 0
		logstr = "大月卡首冲"
		self.ProcMonthCard(moneytype, HF_Atoi(csv["time"]))
		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	}

	if moneytype > 200 {
		logstr = "礼包充值"
	}

	before := self.player.Sql_UserBase.Gem
	//! 更新任务-判断条件为增加基础钻石
	if addgem == 0 {
		//	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_RECHARGE, addmoney*10, 24, addgem, logstr,
		//	//	self.player.Sql_UserBase.Gem+addgem+extragem, extragem, self.player)
		self.player.HandleTask(RechargeGemTask, addmoney*1000, 0, 0)
	} else {
		//	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_RECHARGE, addgem+extragem, 24, addgem, logstr,
		//	//	self.player.Sql_UserBase.Gem+addgem+extragem, extragem, self.player)
		self.player.HandleTask(RechargeGemTask, addgem, 0, 0)
	}
	//self.player.HandleTask(RechargeGemTask, addgem, 0, 0)
	if addgem != 0 {
		self.Sql_UserRecharge.MoneyDay += addgem
	} else {
		self.Sql_UserRecharge.MoneyDay += self.CalGemFroRMB(addmoney)
	}
	self.player.HandleTask(TASK_TYPE_RECHARGE_ONCE, moneytype, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_SINGLE, addmoney, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_MONEY_DAILY, self.Sql_UserRecharge.MoneyDay, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_EQUAL_SINGLE, addmoney, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_ALL, addmoney, 0, 0)

	//! 加钻石
	self.player.AddGem(addgem+extragem, moneytype, addmoney, 0, "充值")
	//! 添加充值记录
	var record JS_RechargeRecord
	record.Money = addmoney
	record.Type = moneytype
	record.Addgem = addgem
	record.ExtraGem = extragem
	record.BeforeGem = before
	record.AfterGem = self.player.Sql_UserBase.Gem
	record.Time = TimeServer().Unix()
	record.OrderId = order_id
	record.Isok = bok

	self.Sql_UserRecharge.record = append(self.Sql_UserRecharge.record, record)
	self.Sql_UserRecharge.Money += addmoney
	self.Sql_UserRecharge.Getallgem += addgem
	self.player.Sql_UserBase.PayGem += addmoney * 1000
	self.Encode()
	//! 处理vip经验
	self.player.AddVipExp(HF_Atoi(csv["vipexp"]), moneytype, addmoney, "充值")

	GetServer().sendLog_AccountChargeSuccess(self.player, "未知", strconv.Itoa(order_id), "CNY", HF_Itof64(addmoney), logstr)
	self.SendInfo("recharge")

	// 增加红包
	self.player.GetModule("redpac").(*ModRedPac).CreateRedWait(addgem)
	self.doRecharge(HF_Atoi(csv["rmb"]), HF_Atoi(csv["diamond"]))
	if moneytype == 101 || moneytype == 102 {
		//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD, moneytype, 0, 0, "月卡奖励领取", 0, 0, self.player)
	}

	self.player.GetModule("activity").(*ModActivity).BuyActivityFund(moneytype)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_CHARGE, moneytype, addmoney, 0, "充值", 0, 0, self.player)
	return 1
}

//! 充值-新版本充值-根据类型
func (self *ModRecharge) RechargeEx(moneytype int, order_id int, bok int, money int, paramLevel int) int {
	//! 如果充值失败
	if bok == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_FAILURE_OF_ORDER_PROCESSING"))
		return -1
	}
	//! 检查订单ID是否重复
	for i := 0; i < len(self.Sql_UserRecharge.record); i++ {
		if self.Sql_UserRecharge.record[i].OrderId == order_id {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_DUPLICATE_ORDER_NUMBER_NO_PROCESSING"))
			return -1
		}
	}

	//csv, ok := GetCsvMgr().GetMoney(moneytype)
	moneyConf, ok := GetCsvMgr().MoneyConfig[moneytype]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_RECHARGE_TYPE_ANOMALY"))
		return -2
	}
	addmoney := moneyConf.Rmb
	addgem := moneyConf.Diamond //HF_Atoi(csv["diamond"])
	extragem := moneyConf.Extra //HF_Atoi(csv["extra"])
	if money != addmoney*100 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_RECHARGE_TYPE_ANOMALY"))
		return -3
	}
	logstr := "充值"
	switch moneyConf.Grade {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11:
		//! 基础充值
		if self.Sql_UserRecharge.baseCounts[moneyConf.Grade-1] >= moneyConf.Time { //HF_Atoi(csv["time"]) {
			extragem = 0
		}
		//! 次数增加1
		self.Sql_UserRecharge.baseCounts[moneyConf.Grade-1] += 1
		break
	case 101, 102, 103, 104:
		//! 检测月卡时间
		//ok, err := self.CheckMonthCard2(moneyConf.Grade)
		//if !ok {
		//	self.player.SendErrInfo("err", err)
		//	return 2
		//}
		if moneyConf.Grade == 101 {
			logstr = "小月卡"
			self.Sql_UserRecharge.MonthCount1++
		} else if moneyConf.Grade == 102 {
			logstr = "大月卡"
			self.Sql_UserRecharge.MonthCount2++
		} else if moneyConf.Grade == 103 {
			logstr = "霸服爽抽"
		} else if moneyConf.Grade == 104 {
			logstr = "霸服永抽"
		}

		extragem = 0
		self.ProcMonthCard(moneyConf.Grade, moneyConf.Time)

		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	case 151:
		ok, err := self.CheckMonthCard2(moneyConf.Grade)
		if !ok {
			self.player.SendErrInfo("err", err)
			return 2
		}

		extragem = 0
		self.ProcMonthCard(moneyConf.Grade, moneyConf.Time)

		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	case 502:
		ok, err := self.CheckMonthCard2(moneyConf.Grade)
		if !ok {
			self.player.SendErrInfo("err", err)
			return 2
		}

		extragem = 0
		self.ProcMonthCard(moneyConf.Grade, 30)

		break
	case 105:
		//! 检查是否有105了, 已经进行了首冲
		//if self.hasMonthCard(105) {
		//	return 3
		//}

		self.Sql_UserRecharge.MonthCount2++
		extragem = 0
		logstr = "大月卡"
		self.ProcMonthCard(moneyConf.Grade, moneyConf.Time)
		//self.player.GetModule("activity").(*ModActivity).HandleTask(50, 0, 0, 0)
		break
	case 301:
		//! 高级皇家犒赏令
		if WARORDER_1 > len(self.Sql_UserRecharge.warOrder) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		realIndex := WARORDER_1 - 1
		if TimeServer().Unix() < self.Sql_UserRecharge.warOrder[realIndex].StartTime ||
			TimeServer().Unix() > self.Sql_UserRecharge.warOrder[realIndex].EndTime {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		if self.Sql_UserRecharge.warOrder[realIndex].BuyState == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		self.Sql_UserRecharge.warOrder[realIndex].BuyState = LOGIC_TRUE

		logstr = "高级皇家犒赏令"
		break
	case 302:
		//高级勇者犒赏令
		if WARORDER_2 > len(self.Sql_UserRecharge.warOrder) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		realIndex := WARORDER_2 - 1
		if TimeServer().Unix() < self.Sql_UserRecharge.warOrder[realIndex].StartTime ||
			TimeServer().Unix() > self.Sql_UserRecharge.warOrder[realIndex].EndTime {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		if self.Sql_UserRecharge.warOrder[realIndex].BuyState == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		self.Sql_UserRecharge.warOrder[realIndex].BuyState = LOGIC_TRUE

		logstr = "高级勇者犒赏令"
		break
	case 303:
		//! 主线战令
		if WARORDERLIMIT_1 > len(self.Sql_UserRecharge.warOrderLimit) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		realIndex := WARORDERLIMIT_1 - 1
		if TimeServer().Unix() < self.Sql_UserRecharge.warOrderLimit[realIndex].StartTime ||
			TimeServer().Unix() > self.Sql_UserRecharge.warOrderLimit[realIndex].EndTime {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_UIShop_Time_Error"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		if self.Sql_UserRecharge.warOrderLimit[realIndex].BuyState == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		self.Sql_UserRecharge.warOrderLimit[realIndex].BuyState = LOGIC_TRUE

		logstr = "主线战令"
		break
	case 304:
		//爬塔战令
		if WARORDERLIMIT_2 > len(self.Sql_UserRecharge.warOrderLimit) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		realIndex := WARORDERLIMIT_2 - 1
		if TimeServer().Unix() < self.Sql_UserRecharge.warOrderLimit[realIndex].StartTime ||
			TimeServer().Unix() > self.Sql_UserRecharge.warOrderLimit[realIndex].EndTime {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_UIShop_Time_Error"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		if self.Sql_UserRecharge.warOrderLimit[realIndex].BuyState == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		self.Sql_UserRecharge.warOrderLimit[realIndex].BuyState = LOGIC_TRUE

		logstr = "爬塔战令"
		break
	case 305:
		// 钻石累消
		if WARORDERLIMIT_3 > len(self.Sql_UserRecharge.warOrderLimit) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		realIndex := WARORDERLIMIT_3 - 1
		if TimeServer().Unix() < self.Sql_UserRecharge.warOrderLimit[realIndex].StartTime ||
			TimeServer().Unix() > self.Sql_UserRecharge.warOrderLimit[realIndex].EndTime {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_UIShop_Time_Error"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		if self.Sql_UserRecharge.warOrderLimit[realIndex].BuyState == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			self.SendBuyFail(moneyConf)
			return 2
		}
		self.Sql_UserRecharge.warOrderLimit[realIndex].BuyState = LOGIC_TRUE

		logstr = "钻石累消"
		break
	case 360, 380:
		if !self.player.GetModule("targettask").(*ModTargetTask).IsCanBuy(moneyConf.Grade) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_PERMANENT_CARDS_CANNOT_BE_PURCHASED"))
			return 2
		}
		self.player.GetModule("targettask").(*ModTargetTask).BuyLevel(moneyConf.Grade)

		logstr = "冒险徽章购买"
		break
	default:
		//打折月卡
		config1 := GetCsvMgr().MonthCard[1]
		config2 := GetCsvMgr().MonthCard[2]
		if config1 != nil && config1.FirstRecharge == moneyConf.Grade {
			logstr = "首次小月卡"
			self.Sql_UserRecharge.MonthCount1++
			extragem = 0
			tempConf := GetCsvMgr().MoneyConfig[10010101]
			self.ProcMonthCard(101, tempConf.Time)
		} else if config2 != nil && config2.FirstRecharge == moneyConf.Grade {
			logstr = "首次大月卡"
			self.Sql_UserRecharge.MonthCount2++
			extragem = 0
			//self.ProcMonthCard(102, moneyConf.Time)
		} else {
			break
		}
	}

	if moneyConf.Grade > 200 {
		logstr = "礼包充值"
	}

	before := self.player.Sql_UserBase.Gem
	//! 更新任务-判断条件为增加基础钻石
	if addgem == 0 {
		//	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_RECHARGE, addmoney*10, 24, addgem, logstr,
		//	//	self.player.Sql_UserBase.Gem+addgem+extragem, extragem, self.player)
		self.player.HandleTask(RechargeGemTask, addmoney*1000, 0, 0)
	} else {
		//	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_RECHARGE, addgem+extragem, 24, addgem, logstr,
		//	//	self.player.Sql_UserBase.Gem+addgem+extragem, extragem, self.player)
		self.player.HandleTask(RechargeGemTask, addgem, 0, 0)
	}
	//self.player.HandleTask(RechargeGemTask, addgem, 0, 0)
	if addgem != 0 {
		self.Sql_UserRecharge.MoneyDay += addgem
	} else {
		self.Sql_UserRecharge.MoneyDay += self.CalGemFroRMB(addmoney)
	}
	self.player.HandleTask(TASK_TYPE_RECHARGE_ONCE, moneyConf.Grade, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_SINGLE, addmoney, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_MONEY_DAILY, self.Sql_UserRecharge.MoneyDay, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_EQUAL_SINGLE, addmoney, 0, 0)
	self.player.HandleTask(TASK_TYPE_RECHARGE_ALL, addmoney, 0, 0)

	self.player.GetModule("specialpurchase").(*ModSpecialPurchase).San_SpecialPurchase.Recharge += addmoney

	//! 加钻石
	self.player.AddGem(addgem+extragem, moneyConf.Grade, addmoney, 0, "充值")
	//成长礼包的订单
	if moneyConf.Grade == 501 {
		self.player.GetModule("weekplan").(*ModWeekPlan).UpdateBuyTime()
	}
	//! 添加充值记录
	var record JS_RechargeRecord
	record.Money = addmoney
	record.Type = moneyConf.Grade
	record.Addgem = addgem
	record.ExtraGem = extragem
	record.BeforeGem = before
	record.AfterGem = self.player.Sql_UserBase.Gem
	record.Time = TimeServer().Unix()
	record.OrderId = order_id
	record.Isok = bok

	self.Sql_UserRecharge.record = append(self.Sql_UserRecharge.record, record)
	self.Sql_UserRecharge.Money += addmoney
	self.Sql_UserRecharge.Getallgem += addgem
	self.player.Sql_UserBase.PayGem += addmoney * 1000
	self.Encode()
	//! 处理vip经验
	oldVip := self.player.Sql_UserBase.Vip
	//self.player.AddVipExp(HF_Atoi(csv["vipexp"]), moneytype, addmoney, "充值")
	self.player.AddVipExp(moneyConf.Vipexp, moneyConf.Grade, addmoney, "充值")

	GetServer().sendLog_AccountChargeSuccess(self.player, "未知", strconv.Itoa(order_id), "CNY", HF_Itof64(addmoney), logstr)
	self.SendInfo("recharge")

	// 增加红包
	self.player.GetModule("redpac").(*ModRedPac).CreateRedWait(addgem)
	//self.doRecharge(HF_Atoi(csv["rmb"]), HF_Atoi(csv["diamond"]))
	self.doRecharge(moneyConf.Rmb, moneyConf.Diamond)
	if moneytype == 101 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买普通月卡", 0, self.player.Sql_UserBase.Vip, self.player)
	} else if moneytype == 102 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD_HIGH, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买高级月卡", 0, self.player.Sql_UserBase.Vip, self.player)
	}

	if moneytype == 151 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD_GOLD, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买金卡", 0, self.player.Sql_UserBase.Vip, self.player)
	} else if moneytype == 301 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_WARORDER_BUY_1, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买皇家犒赏令", 0, self.player.Sql_UserBase.Vip, self.player)
	} else if moneytype == 302 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_WARORDER_BUY_2, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买勇士犒赏令", 0, self.player.Sql_UserBase.Vip, self.player)
	} else if moneytype == 303 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_RECHARGE_LIMIT_MAIN, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买主线战令", 0, self.player.Sql_UserBase.Vip, self.player)
	} else if moneytype == 304 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_RECHARGE_LIMIT_TOWER, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买试炼之塔战令", 0, self.player.Sql_UserBase.Vip, self.player)
	} else if moneytype == 305 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_RECHARGE_LIMIT_DIAMOND, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买钻石累消战令", 0, self.player.Sql_UserBase.Vip, self.player)
	}else if moneytype == 360 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_TARGETTASK_BADEG_BUY, moneyConf.Grade, moneyConf.Rmb, oldVip, "激活冒险徽章", 0, self.player.Sql_UserBase.Vip, self.player)
	} else if moneytype == 380 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_TARGETTASK_BADEG_TOWER_BUY, moneyConf.Grade, moneyConf.Rmb, oldVip, "激活试炼徽章", 0, self.player.Sql_UserBase.Vip, self.player)
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_CHARGE, moneyConf.Grade, addmoney, 0, "充值", 0, 0, self.player)

	self.player.GetModule("activity").(*ModActivity).BuyActivityFund(moneyConf.Grade)
	self.player.GetModule("specialpurchase").(*ModSpecialPurchase).HandleRecharge(moneyConf.Grade, paramLevel)
	self.player.GetModule("activitygift").(*ModActivityGift).GetAllAward(moneyConf.Grade)
	self.player.GetModule("herogrow").(*ModHeroGrow).HandleRecharge(moneyConf.Grade, oldVip)

	nType := self.player.GetModule("luckshop").(*ModLuckShop).HandleRecharge(moneyConf.Grade)
	switch nType {
	case ACTIVITY_GIFT_TYPE_DISCOUNT:
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GIFT_LOW, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买特惠礼包", 0, self.player.Sql_UserBase.Vip, self.player)
	case ACTIVITY_GIFT_TYPE_STAR:
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_STAR_HERO_BUY_LIMIT, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买星辰限时", 0, self.player.Sql_UserBase.Vip, self.player)
	case ACTIVITY_GIFT_TYPE_STAR_HERO:
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_STAR_HERO_BUY, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买星辰英雄", 0, self.player.Sql_UserBase.Vip, self.player)
	}

	self.player.GetModule("viprecharge").(*ModVipRecharge).HandleRecharge(moneyConf.Grade)
	self.player.GetModule("activity").(*ModActivity).HandleRecharge(moneyConf.Grade)
	rel := self.player.GetModule("fund").(*ModFund).HandleRecharge(moneyConf.Grade)
	if rel == 1 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FUND_1, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买勇者基金", 0, self.player.Sql_UserBase.Vip, self.player)
	} else if rel == 2 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FUND_2, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买至尊基金", 0, self.player.Sql_UserBase.Vip, self.player)
	} else if rel == 3 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FUND_3, moneyConf.Grade, moneyConf.Rmb, oldVip, "购买无敌基金", 0, self.player.Sql_UserBase.Vip, self.player)
	}
	// 精英直升
	//self.player.GetModule("eliteupgrade").(*ModActivity).HandleRecharge(moneyConf.Grade)

	return 1
}

/*
func (self *ModRecharge) CheckMonthCard(moneytype int) (bool, string) {
	nowtime := TimeServer().Unix()
	month := self.player.GetModule("activity").(*ModActivity).Sql_Activity.month
	for i := 0; i < len(month); i++ {
		if month[i].Id == moneytype {
			//! 过期处理
			allday := int64(month[i].Day)
			useday := (nowtime - month[i].StartTime) / 86400
			if allday-useday >= 120 {
				return false, GetCsvMgr().GetText("STR_MONTH_CARD")
			}
		}
	}
	return true, ""
}

*/

/*
func (self *ModRecharge) CheckMonthCard2(moneytype int) (bool, string) {
	nowtime := TimeServer().Unix()
	month := self.player.GetModule("activity").(*ModActivity).Sql_Activity.month
	for i := 0; i < len(month); i++ {
		if month[i].Id == moneytype {
			//! 过期处理
			allday := int64(month[i].Day)
			useday := (nowtime - month[i].StartTime) / 86400
			if allday-useday > 120 {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MONTH_CARD_LEFT"))
				return false, GetCsvMgr().GetText("STR_MONTH_CARD_LEFT")
			}
		}
	}
	return true, ""
}

*/

//! 得到下次刷新时间
func (self *ModRecharge) GetProcMonthCardTime() int64 {
	now := TimeServer()
	if now.Hour() < 5 {
		return time.Date(now.Year(), now.Month(), now.Day(), -24+5, 0, 0, 0, time.Local).Unix()
	} else {
		return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local).Unix()
	}
}

func (self *ModRecharge) CalActivityStartTime(start string) int64 {
	startday := HF_Atoi(start)
	if startday > 0 {
		return GetServer().GetOpenServer() + int64((startday-1)*DAY_SECS)
	} else if startday < 0 {
		rtime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
		firstTime := HF_CalPlayerCreateTime(rtime.Unix(), 0)
		return firstTime
	}
	return 0
}

//! 处理月卡
func (self *ModRecharge) ProcMonthCard(moneytype int, day int) {
	//! 检查月卡开始时间
	nowtime := TimeServer().Unix()
	useractivity := self.player.GetModule("activity").(*ModActivity)
	for i := 0; i < len(useractivity.Sql_Activity.month); i++ {
		if useractivity.Sql_Activity.month[i].Id == moneytype {
			if moneytype == 502 {
				activity := GetActivityMgr().GetActivity(3012)
				if activity == nil {
					return
				}
				useractivity.Sql_Activity.month[i].StartTime = self.CalActivityStartTime(activity.info.Start)
				for _, v := range GetCsvMgr().ActivityOverflowGifts {
					useractivity.Sql_Activity.month[i].RewardSign[v.Id] = 0
				}
				self.player.GetModule("activity").(*ModActivity).SendMonthCard()
				return
			} else {
				if (nowtime-useractivity.Sql_Activity.month[i].StartTime)/86400 > int64(useractivity.Sql_Activity.month[i].Day) {
					useractivity.Sql_Activity.month[i].StartTime = nowtime
					useractivity.Sql_Activity.month[i].Day = day
				} else {
					useractivity.Sql_Activity.month[i].Day += day
				}
				self.player.GetModule("activity").(*ModActivity).SendMonthCard()
				return
			}
		}
	}

	if moneytype == 502 {
		activity := GetActivityMgr().GetActivityOverflow()
		if activity == nil {
			return
		}
		var month JS_MonthCard
		month.StartTime = self.CalActivityStartTime(activity.info.Start)
		month.Id = 502
		month.Day = day
		month.Stage = activity.Id
		month.RewardSign = make(map[int]int)
		for _, v := range activity.items {
			month.RewardSign[v.Id] = 0
		}
		useractivity.Sql_Activity.month = append(useractivity.Sql_Activity.month, month)
	} else {
		// 首次充值
		useractivity.Sql_Activity.month = append(useractivity.Sql_Activity.month,
			JS_MonthCard{moneytype, self.GetProcMonthCardTime(), day, 0, make(map[int]int), 0})
	}

	useractivity.ReloadActivity()
	self.player.GetModule("activity").(*ModActivity).SendMonthCard()
}

//! 获取首充奖励 ret 1:未充值 2:已经领取了奖励 -1:操作失败
//  新版 增加二充 三充
func (self *ModRecharge) GetFirsetAward(Type int) {

	if Type < RECHARGE_FIRST || Type > RECHARGE_THIRD {
		self.player.SendErr(GetCsvMgr().GetText("STR_DUNGEON_TEAM_FLAG_ERROR"))
		return
	}

	sign := (self.Sql_UserRecharge.Firstaward / int(math.Pow(10, float64(Type-1)))) % 10
	if sign != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_TIMES_ALREADY_GET"))
		return
	}

	activityId := FIRST_RECHARGE_TYPE_OFFSET + Type
	//读取奖励
	activityitem := GetActivityMgr().GetActivityItem(activityId)
	if activityitem == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_DUNGEON_TEAM_FLAG_ERROR"))
		return
	}

	//看是否满足
	if len(self.Sql_UserRecharge.record) == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_TIMES_NOT_ENOUGH"))
		return
	}

	//计算天数 是否满足
	nextTime := HF_CalPlayerCreateTime(self.Sql_UserRecharge.record[0].Time, Type-1)
	if nextTime > TimeServer().Unix() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NOBILITYTASK_TAKE_ERROR"))
		return
	}
	//这个奖励因为格子不够 ，所以COST也是填的奖励  20200908
	getItemId := make([]int, 0)
	getItemNum := make([]int, 0)
	for i := 0; i < len(activityitem.Item); i++ {
		if activityitem.Item[i] > 0 {
			getItemId = append(getItemId, activityitem.Item[i])
			getItemNum = append(getItemNum, activityitem.Num[i])
		}
	}
	for i := 0; i < len(activityitem.CostItem); i++ {
		if activityitem.CostItem[i] > 0 {
			getItemId = append(getItemId, activityitem.CostItem[i])
			getItemNum = append(getItemNum, activityitem.CostNum[i])
		}
	}

	outitem := self.player.AddObjectLst(getItemId, getItemNum, "首充奖励", 0, 0, 0)

	self.Sql_UserRecharge.Firstaward += int(math.Pow(10, float64(Type-1)))
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FIRST_GET, self.Sql_UserRecharge.Firstaward, 0, 0, "领取首充奖励", 0, 0, self.player)

	var msg S2C_GetFirsetAward
	msg.Cid = "getfirstaward"
	msg.Award = outitem
	msg.Info = self.Sql_UserRecharge
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)

	return
}

//! 基金领取
func (self *ModRecharge) SendInfo(cid string) {
	self.CheckWarOrder()
	self.CheckWarOrderLimit()

	self.Encode()

	var msg S2C_RechargeInfo
	msg.Cid = cid
	msg.Info = self.Sql_UserRecharge
	msg.Vip = self.player.Sql_UserBase.Vip
	msg.Vipexp = self.player.Sql_UserBase.VipExp
	msg.CurGem = self.player.Sql_UserBase.Gem
	msg.Vipbox = self.Sql_UserRecharge.VipBox
	msg.Fundtype = self.Sql_UserRecharge.FundType
	msg.Fundget = self.Sql_UserRecharge.fundget
	msg.Fundtotal = GetHeroSupportMgr().GetJJ()
	//发送主线战令状态、领取进度
	msg.WarOrderLimit = self.Sql_UserRecharge.warOrderLimit
	msg.Ret = 0
	//! 屏蔽充值记录
	msg.Info.record = []JS_RechargeRecord{}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(cid, smsg)
}

func (self *ModRecharge) CheckOpen() {
	isSend := self.CheckWarOrder()
	if !isSend {
		return
	}
	self.CalWarOrder(WARORDER_1, 0)
	self.CalWarOrder(WARORDER_2, 0)
}

func (self *ModRecharge) CheckOpenLimit() {
	isSend := self.CheckWarOrderLimit()
	if !isSend {
		return
	}
	var msg S2C_WarOrderLimitInfo
	msg.Cid = "warorderlimitinfo"
	msg.WarOrderLimit = self.Sql_UserRecharge.warOrderLimit
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}

func (self *ModRecharge) hasMonthCard(moneytype int) bool {
	//! 检查月卡开始时间
	useractivity := self.player.GetModule("activity").(*ModActivity)
	for i := 0; i < len(useractivity.Sql_Activity.month); i++ {
		if useractivity.Sql_Activity.month[i].Id == moneytype {
			return true
		}
	}
	return false
}

// 重置金额回调
func (self *ModRecharge) doRecharge(rmb int, diamond int) {
	gem := 0
	if diamond > 0 {
		gem = diamond
	} else {
		gem = rmb * 1000
	}
	self.player.GetModule("dailyrecharge").(*ModDailyRecharge).doTask(rmb)

	configNum := GetCsvMgr().getInitNum(5)
	if configNum != 0 {
		self.player.AddBossMoney(gem/configNum, 0, 0, "充值送水晶")
	}
}

//  VIP每日福利
func (self *ModRecharge) GetVipDaily(viplevel int) {

	outitem := make([]PassItem, 0)

	if self.Sql_UserRecharge.VipDailyReward != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_TIMES_ALREADY_GET"))
		return
	}

	if self.player.Sql_UserBase.Vip != viplevel {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_TIMES_ALREADY_GET"))
		return
	}
	//读取奖励
	vipConfig := GetCsvMgr().GetVipConfig(viplevel)
	if vipConfig == nil || len(vipConfig.FreeItems) != len(vipConfig.FreeNums) {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_KINGTASK_ARISTOCRATIC_MISALLOCATION"))
		return
	}

	for i := 0; i < len(vipConfig.FreeItems); i++ {
		if vipConfig.FreeNums[i] > 0 {
			self.player.AddObject(vipConfig.FreeItems[i], vipConfig.FreeNums[i], 24, 0, 0, fmt.Sprintf("VIP每日奖励%d", viplevel))
			outitem = append(outitem, PassItem{ItemID: vipConfig.FreeItems[i], Num: vipConfig.FreeNums[i]})
		} else {
			break
		}
	}

	self.Sql_UserRecharge.VipDailyReward = 1

	var msg S2C_GetVipDaily
	msg.Cid = "getvipdailyreward"
	msg.Award = outitem
	msg.VipDailyState = self.Sql_UserRecharge.VipDailyReward
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("recharge_firstaward", smsg)

}

//  VIP每周购买福利
func (self *ModRecharge) BuyVipWeek(viplevel int) {

	if viplevel > self.player.Sql_UserBase.Vip {
		return
	}

	//读取奖励配置
	vipConfig := GetCsvMgr().GetVipConfig(viplevel)
	if vipConfig == nil || len(vipConfig.WeekItems) != len(vipConfig.WeekNums) {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_KINGTASK_ARISTOCRATIC_MISALLOCATION"))
		return
	}

	sign := (self.Sql_UserRecharge.VipWeekBuy / int64(math.Pow(10, float64(viplevel)))) % 10
	if sign > int64(vipConfig.WeekLimit) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_TIMES_ALREADY_GET"))
		return
	}

	var res = make(map[int]*Item)
	res[91000002] = &Item{ItemId: 91000002, ItemNum: vipConfig.WeekPrice}

	if err := self.player.hasItemMapOk(res); err != nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_NO_COST"))
		return
	}

	costItem := self.player.RemoveObjectItemMap(res, "购买VIP每周礼包", viplevel, 0, 0)

	outitem := make([]PassItem, 0)

	for i := 0; i < len(vipConfig.WeekItems); i++ {
		if vipConfig.WeekItems[i] > 0 {
			self.player.AddObject(vipConfig.WeekItems[i], vipConfig.WeekNums[i], 24, 0, 0, fmt.Sprintf("VIP每周礼包购买%d", viplevel))
			outitem = append(outitem, PassItem{ItemID: vipConfig.WeekItems[i], Num: vipConfig.WeekNums[i]})
		} else {
			break
		}
	}

	self.Sql_UserRecharge.VipWeekBuy += int64(math.Pow(10, float64(viplevel)))

	var msg S2C_BuyVipWeek
	msg.Cid = "buyvipweek"
	msg.Award = outitem
	msg.Cost = costItem
	msg.BuyVipWeekState = self.Sql_UserRecharge.VipWeekBuy
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("buyvipweek", smsg)
}

func (self *ModRecharge) GetWarOrderReward(msg *C2S_GetWarOrderReward) {

	if msg.Type <= 0 || msg.Type > len(self.Sql_UserRecharge.warOrder) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_MSG_ERROR"))
		return
	}
	realIndex := msg.Type - 1
	warOrder := self.Sql_UserRecharge.warOrder[realIndex]
	if TimeServer().Unix() < warOrder.StartTime ||
		TimeServer().Unix() > warOrder.EndTime {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_WARORDER_NOT_OPEN"))
		return
	}

	var msgRel S2C_GetWarOrderReward
	msgRel.Cid = "getwarorderreward"
	for i := 0; i < len(warOrder.WarOrderTask); i++ {
		if warOrder.WarOrderTask[i].Id == msg.Id {
			config := GetCsvMgr().GetWarOrderConfig(msg.Id)
			if config == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_WARORDER_CONFIG_ERROR"))
				return
			}
			var getItem = make(map[int]*Item)
			if warOrder.WarOrderTask[i].State == CANTAKE {
				AddItemMapHelper(getItem, config.FreeAward, config.FreeNum)
				warOrder.WarOrderTask[i].State = TAKEN
			}
			if warOrder.WarOrderTask[i].State != CANTFINISH {
				if warOrder.BuyState == LOGIC_TRUE && warOrder.WarOrderTask[i].BuyGet == LOGIC_FALSE {
					AddItemMapHelper(getItem, config.GoldAward, config.GoldNum)
					warOrder.WarOrderTask[i].BuyGet = LOGIC_TRUE
				}
			}

			if len(getItem) == 0 {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_WARORDER_ALREADY_GET"))
				return
			}
			logId := 0
			logDec := ""
			if msg.Type == 1 {
				logId = LOG_WARORDER_GET_1
				logDec = "领取皇家犒赏令奖励"
			} else if msg.Type == 2 {
				logId = LOG_WARORDER_GET_2
				logDec = "领取勇士犒赏令奖励"
			}
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, logId, warOrder.WarOrderTask[i].Id, 0, 0, logDec, 0, 0, self.player)

			msgRel.GetItem = self.player.AddObjectItemMap(getItem, logDec, warOrder.WarOrderTask[i].Id, 0, 0)
			msgRel.WarOrderTask = warOrder.WarOrderTask[i]
			smsg, _ := json.Marshal(&msgRel)
			self.player.SendMsg(msgRel.Cid, smsg)
			return
		}
	}

	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_MSG_ERROR"))
	return
}

func (self *ModRecharge) GetWarOrderLimitReward(msg *C2S_GetWarOrderLimitReward) {

	// 战令类型不存在
	if msg.Type <= 0 || msg.Type > len(self.Sql_UserRecharge.warOrderLimit) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_MSG_ERROR"))
		return
	}
	realIndex := msg.Type - 1
	warOrderLimit := self.Sql_UserRecharge.warOrderLimit[realIndex]
	// 已购买或活动期间才可以领取奖励
	if warOrderLimit.BuyState == LOGIC_FALSE && (TimeServer().Unix() < warOrderLimit.StartTime || TimeServer().Unix() > warOrderLimit.EndTime) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_WARORDER_NOT_OPEN"))
		return
	}

	var msgRel S2C_GetWarOrderLimitReward
	msgRel.Cid = "getwarorderlimitreward"
	for i := 0; i < len(warOrderLimit.WarOrderTask); i++ {
		if warOrderLimit.WarOrderTask[i].Id == msg.Id {			// 领取的奖励id
			config := GetCsvMgr().WarOrderLimitConfig[msg.Id]
			if config == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_WARORDER_CONFIG_ERROR"))
				return
			}
			var getItem = make(map[int]*Item)
			if warOrderLimit.WarOrderTask[i].State == CANTAKE {
				AddItemMapHelper(getItem, config.FreeAward, config.FreeNum)
				warOrderLimit.WarOrderTask[i].State = TAKEN
			}
			if warOrderLimit.WarOrderTask[i].State != CANTFINISH {
				if warOrderLimit.BuyState == LOGIC_TRUE && warOrderLimit.WarOrderTask[i].BuyGet == LOGIC_FALSE {	// 已购买且未领取
					AddItemMapHelper(getItem, config.GoldAward, config.GoldNum)
					warOrderLimit.WarOrderTask[i].BuyGet = LOGIC_TRUE
				}
			}

			if len(getItem) == 0 {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_WARORDER_ALREADY_GET"))		// 奖励已领取
				return
			}
			logId := 0
			logDec := ""
			if msg.Type == 1 {
				logId = LOG_RECHARGE_LIMIT_MAIN_GET
				logDec = "领取主线战令奖励"
			} else if msg.Type == 2 {
				logId = LOG_RECHARGE_LIMIT_TOWER_GET
				logDec = "领取试炼之塔战令奖励"
			}else if msg.Type ==3 {
				logId = LOG_RECHARGE_LIMIT_DIAMOND_GET
				logDec = "领取钻石累消战令奖励"
			}
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, logId, warOrderLimit.WarOrderTask[i].Id, 0, 0, logDec, 0, 0, self.player)
			msgRel.GetItem = self.player.AddObjectItemMap(getItem, logDec, warOrderLimit.WarOrderTask[i].Id, 0, 0)
			msgRel.WarOrderTask = warOrderLimit.WarOrderTask[i]
			smsg, _ := json.Marshal(&msgRel)
			self.player.SendMsg(msgRel.Cid, smsg)
			return
		}
	}

	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_MSG_ERROR"))
	return
}

// 购买战令
func (self *ModRecharge) WarOrderBuy(msg *C2S_WarOrderBuy) {

	if msg.Type <= 0 || msg.Type > len(self.Sql_UserRecharge.warOrder) || msg.BuyNum <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_MSG_ERROR"))
		return
	}
	realIndex := msg.Type - 1
	warOrder := self.Sql_UserRecharge.warOrder[realIndex]
	if TimeServer().Unix() < warOrder.StartTime ||
		TimeServer().Unix() > warOrder.EndTime {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_WARORDER_NOT_OPEN"))
		return
	}

	config := GetCsvMgr().WarOrderParam[msg.Type]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_WARORDER_CONFIG_ERROR"))
		return
	}

	var getItem = make(map[int]*Item)
	var costItem = make(map[int]*Item)

	AddItemMapHelper3(getItem, config.BuyItem, config.BuyNum*msg.BuyNum)
	AddItemMapHelper3(costItem, config.SwitchItem, config.SwitchNum*msg.BuyNum)

	var msgRel S2C_GetWarOrderBuy
	msgRel.Cid = "warorderbuy"
	msgRel.GetItem = self.player.AddObjectItemMap(getItem, "购买犒赏令", msg.Type, msg.BuyNum, 0)
	msgRel.CostItem = self.player.RemoveObjectItemMap(costItem, "购买犒赏令", msg.Type, msg.BuyNum, 0)
	smsg, _ := json.Marshal(&msgRel)
	self.player.SendMsg(msgRel.Cid, smsg)
}

func (self *ModRecharge) CheckMonthCard2(moneytype int) (bool, string) {
	nowtime := TimeServer().Unix()
	month := self.player.GetModule("activity").(*ModActivity).Sql_Activity.month
	for i := 0; i < len(month); i++ {
		if month[i].Id == moneytype {
			//! 过期处理
			allday := int64(month[i].Day)
			useday := (nowtime - month[i].StartTime) / 86400
			if allday-useday > 0 {
				return false, GetCsvMgr().GetText("STR_MOD_ACTIVITY_CANT_REPEAT_PURCHASE")
			}
		}
	}
	return true, ""
}

func (self *ModRecharge) IsBuy(id int) bool {
	if id != WARORDERLIMIT_1 && id != WARORDERLIMIT_2 {
		return false
	}

	if WARORDERLIMIT_1 > len(self.Sql_UserRecharge.warOrderLimit) {
		return false
	}

	realIndex := id - 1
	if self.Sql_UserRecharge.warOrderLimit[realIndex].BuyState == LOGIC_TRUE {
		return true
	}

	return false
}

func (self *ModRecharge) SendBuyFail(config *MoneyConfig) {

	switch config.Grade {
	case 301:
	case 302:
	case 303:
	case 304:
	default:
		return
	}

	lstItem := make([]PassItem, 0)
	lstItem = append(lstItem, PassItem{ITEM_GEM, self.CalGemFroRMB(config.Rmb)})
	self.player.GetModule("bag").(*ModBag).SendOnItem(lstItem)
}

func (self *ModRecharge) CalGemFroRMB(rmb int) int {
	return rmb * 1000
}
