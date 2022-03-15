package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	TAG_TYPE_TASK = 1 //任务
	TAG_TYPE_BUY  = 2 //购买
	TAG_TYPE_GIFT = 3 //礼包
)

const (
	AllDAY = 21 //活动总天数
)

// 七日活动玩家
type Js_WeekPlanPlayer struct {
	Uid int64 `json:"uid"` //! 玩家Id
}

// 当活动结束之后, 由系统自动发送给客户端。
type JS_WeekPlanStatus struct {
	IsOpen   bool  `json:"open"`     //! 是否开放
	CanPick  bool  `json:"canpick"`  //! 是否能领取
	OpenTime int64 `json:"opentime"` //! 开放时间
}

//! 7日活动管理类
type WeekPlanMgr struct {
	Locker         *sync.RWMutex //! 活动锁
	LastUpdateTime int64         //! 上次更新，每小时更新一次

	Players      map[int64]*Js_WeekPlanPlayer // 参与玩家表
	PlayerLocker *sync.RWMutex                // 玩家锁
}

var weekplanmgrsingleton *WeekPlanMgr = nil

//! 获取活动管理类
func GetWeekPlanMgr() *WeekPlanMgr {
	if weekplanmgrsingleton == nil {
		weekplanmgrsingleton = new(WeekPlanMgr)
		weekplanmgrsingleton.Players = make(map[int64]*Js_WeekPlanPlayer)
		weekplanmgrsingleton.Locker = new(sync.RWMutex)
		weekplanmgrsingleton.LastUpdateTime = 0
		weekplanmgrsingleton.PlayerLocker = new(sync.RWMutex)
	}

	return weekplanmgrsingleton
}

//! 获取数据
func (self *WeekPlanMgr) GetData() {
	// 未初始化则初始化
	if len(self.Players) <= 0 {
		self.Players = make(map[int64]*Js_WeekPlanPlayer)

		// 从表中选出所有参与玩家的数据
		var player Js_WeekPlanPlayer
		sql := "select uid from san_weekplan"
		res := GetServer().DBUser.GetAllDataEx(sql, &player)
		for i := 0; i < len(res); i++ {
			data := res[i].(*Js_WeekPlanPlayer)

			self.Players[data.Uid] = data
		}
	}
}

//! 初始化7日活动状态
func (self *ModWeekPlan) InitWeekPlanStatus() {
	if self.Sql_WeekPlan.Stage == 0 {
		self.Sql_WeekPlan.Stage = 1
	}
	self.WeekPlanStatus = make([]JS_WeekPlanStatus, 0)
	newStage := 0
	//! 初始化7天每天的开放和领取情况
	for i := 0; i < AllDAY; i++ {
		//var id = (i+1)*100000000 + 10000 + 1
		var id = 1000000 + (i+1)*10000 + 101
		csv, ok := GetCsvMgr().GetSevenDay(id)
		if !ok {
			return
		}

		startDay := HF_Atoi(csv["start"])
		continued := HF_Atoi(csv["continued"])
		show := HF_Atoi(csv["show"])
		//! 开服事件在开始之前 或者 逾期之后
		rTime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)

		if startDay >= 0 {
			LogError("seven startday is error")
			return
		}
		correctTime := HF_CalPlayerCreateTime(rTime.Unix(), 0)

		endTime := correctTime + int64(continued+show) + int64(-(startDay + 1)*DAY_SECS)
		now := TimeServer().Unix()
		starttime := correctTime + int64(-(startDay + 1)*DAY_SECS)
		if now < starttime || now >= endTime {
			self.WeekPlanStatus = append(self.WeekPlanStatus, JS_WeekPlanStatus{false, false, starttime})
		} else {
			if now < starttime+int64(continued) {
				self.WeekPlanStatus = append(self.WeekPlanStatus, JS_WeekPlanStatus{true, true, starttime})
				var idNext = 1000000 + (self.Sql_WeekPlan.Stage*7+1)*10000 + 101
				if id >= idNext {
					newStage = self.Sql_WeekPlan.Stage + 1
				}
			} else {
				self.WeekPlanStatus = append(self.WeekPlanStatus, JS_WeekPlanStatus{false, true, starttime})
			}
		}
	}

	for _, value := range GetCsvMgr().SevenDayTask_CSV {
		self.GetTask(value.Id, true)
	}
	if newStage > self.Sql_WeekPlan.Stage {
		itemMap := make(map[int]*Item)
		var powerok = false
		actitem, ok := self.player.GetModule("activity").(*ModActivity).Sql_Activity.info[GROWTH_GIFT_BASE]
		if ok {
			if actitem.Done == 1 {
				powerok = true
			}
		}
		//发放没领取的奖励
		for i := 0; i < len(self.Sql_WeekPlan.taskinfo); i++ {
			csv_weal, ok := GetCsvMgr().GetSevenDay(self.Sql_WeekPlan.taskinfo[i].Taskid)
			if !ok {
				continue
			}
			stageJudge := (self.Sql_WeekPlan.taskinfo[i].Taskid - 1000000) / 10000
			if stageJudge > self.Sql_WeekPlan.Stage*7 {
				continue
			}
			if self.Sql_WeekPlan.taskinfo[i].Finish == 0 {
				//看是不是可以免费领取TYPE3
				isCan := false
				rebate_type := HF_Atoi(csv_weal["rebate_type"])
				vipLimit := HF_Atoi(csv_weal["param1"])
				if (rebate_type == 1 && powerok) || (rebate_type == 2 && self.player.Sql_UserBase.Vip >= vipLimit) {
					isCan = true
				}
				if HF_Atoi(csv_weal["tag_type"]) == 3 && isCan {

				} else {
					continue
				}
			}
			if self.Sql_WeekPlan.taskinfo[i].Pickup == 1 {
				continue
			}

			self.Sql_WeekPlan.taskinfo[i].Pickup = 1
			items := make([]int, 0)
			nums := make([]int, 0)
			for i := 0; i < 4; i++ {
				itemid := HF_Atoi(csv_weal[fmt.Sprintf("item%d", i+1)])
				if itemid == 0 {
					continue
				}
				num := HF_Atoi(csv_weal[fmt.Sprintf("num%d", i+1)])
				items = append(items, itemid)
				nums = append(nums, num)
			}
			AddItemMapHelper(itemMap, items, nums)
		}
		for _, pointConfig := range GetCsvMgr().SevendayAward {
			if self.Sql_WeekPlan.isGetMark[pointConfig.Id] == LOGIC_TRUE {
				continue
			}
			if pointConfig.Stage != self.Sql_WeekPlan.Stage {
				continue
			}
			if self.Sql_WeekPlan.Point < pointConfig.NeedPoint {
				continue
			}
			AddItemMapHelper(itemMap, pointConfig.Items, pointConfig.Nums)
			self.Sql_WeekPlan.isGetMark[pointConfig.Id] = LOGIC_TRUE
		}

		//发送邮件
		if len(itemMap) > 0 {
			mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_GET_SEVEN]
			if ok {
				itemLst := make([]PassItem, 0)
				for _, v := range itemMap {
					itemLst = append(itemLst, PassItem{ItemID: v.ItemId, Num: v.ItemNum})
				}
				self.player.GetModule("mail").(*ModMail).AddMail(1,
					1, 0, mailConfig.Mailtitle, mailConfig.Mailtxt, GetCsvMgr().GetText("STR_SYS"), itemLst, false, 0)
			}
		}

		self.Sql_WeekPlan.Stage = newStage
		self.Sql_WeekPlan.Point = 0
	}
}

// 将参与玩家加入临时表 临时表每次开服读取数据库 所以不需要保存
func (self *WeekPlanMgr) AddPlayer(player *Player) {
	_, ok := self.Players[player.Sql_UserBase.Uid]
	if !ok {
		self.PlayerLocker.Lock()
		self.Players[player.Sql_UserBase.Uid] = &Js_WeekPlanPlayer{player.Sql_UserBase.Uid}
		defer self.PlayerLocker.Unlock()
	}
}

//! 周期计划数据库
type San_WeekPlan struct {
	Uid           int64
	Finishcount   int
	Type1         string
	Type2         string
	Type3         string
	Type4         string
	Type5         string
	Type6         string
	Type7         string
	Type8         string
	Completeaward int
	Lastupdtime   int64
	Taskinfo      string
	BoxNum        int    // 宝箱数量
	Point         int    // 获得的点数
	IsGet         int    // 是否领取
	IsGetMark     string // 是否领取
	Stage         int    // 活动阶段

	taskinfo  []JS_TaskInfo
	isGetMark map[int]int

	DataUpdate
}

// 进度
type JS_TaskInfo struct {
	Taskid    int `json:"taskid"`    // 任务Id
	Tasktypes int `json:"tasktypes"` // 任务类型
	Plan      int `json:"plan"`      // 进度
	Finish    int `json:"finish"`    // 是否完成
	Pickup    int `json:"pickup"`    // 是否领取奖励
}

//! 周期计划
type ModWeekPlan struct {
	player         *Player
	Sql_WeekPlan   San_WeekPlan //! 数据库结构
	chg            []JS_TaskInfo
	init           bool
	WeekPlanStatus []JS_WeekPlanStatus //! 活动状态
}

func (self *ModWeekPlan) OnGetData(player *Player) {
	self.player = player
	self.chg = make([]JS_TaskInfo, 0)
}

func (self *ModWeekPlan) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_weekplan` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_WeekPlan, "san_weekplan", self.player.ID)

	if self.Sql_WeekPlan.Uid <= 0 {
		self.Sql_WeekPlan.Uid = self.player.ID
		self.Sql_WeekPlan.taskinfo = make([]JS_TaskInfo, 0)
		self.Sql_WeekPlan.isGetMark = make(map[int]int, 0)
		self.Sql_WeekPlan.Lastupdtime = 0
		self.Encode()
		InsertTable("san_weekplan", &self.Sql_WeekPlan, 0, true)
	} else {
		self.Decode()
	}

	if self.Sql_WeekPlan.isGetMark == nil {
		self.Sql_WeekPlan.isGetMark = make(map[int]int, 0)
	}

	self.Sql_WeekPlan.Init("san_weekplan", &self.Sql_WeekPlan, true)

	self.chg = make([]JS_TaskInfo, 0)

	//self.InitWeekPlanStatus()  移到sendinfo里，可以兼容5点刷新
	//! 判断积分奖励
	//self.CheckAward()
}

func (self *ModWeekPlan) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "getweekplanaward":
		var c2s_msg C2S_GetWeepPlanAward
		json.Unmarshal(body, &c2s_msg)
		self.GetAward(c2s_msg.Id)
		return true
	case "weekplan_buyitem":
		var c2s_msg C2S_GetWeepPlanAward
		json.Unmarshal(body, &c2s_msg)
		self.BuyItem(c2s_msg.Id)
		return true
	case "weekplan_funditem":
		var c2s_msg C2S_GetWeepPlanAward
		json.Unmarshal(body, &c2s_msg)
		self.FundItem(c2s_msg.Id)
		return true
	case "getweekplanpoint":
		var c2s_msg C2S_GetWeepPlanAward
		json.Unmarshal(body, &c2s_msg)
		self.GetWeekplanPoint(c2s_msg.Id)
		return true
	}

	return false
}

func (self *ModWeekPlan) OnSave(sql bool) {
	self.Encode()
	self.Sql_WeekPlan.Update(sql)
}

func (self *ModWeekPlan) Decode() {
	//! 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_WeekPlan.Taskinfo), &self.Sql_WeekPlan.taskinfo)
	json.Unmarshal([]byte(self.Sql_WeekPlan.IsGetMark), &self.Sql_WeekPlan.isGetMark)
}

func (self *ModWeekPlan) Encode() {
	self.Sql_WeekPlan.Taskinfo = HF_JtoA(self.Sql_WeekPlan.taskinfo)
	self.Sql_WeekPlan.IsGetMark = HF_JtoA(self.Sql_WeekPlan.isGetMark)
}

func (self *ModWeekPlan) CheckAward() {

	// 已领过
	if self.Sql_WeekPlan.IsGet == 1 {
		return
	}

	if GetServer().GetOpenTime() < 8 {
		return
	}

	mailConfig, ok := GetCsvMgr().MailConfig[3001] //特殊处理 3001 为邮件模板
	if !ok {
		LogError("mail table need check")
		return
	}

	pMail := self.player.GetModule("mail").(*ModMail)
	if pMail == nil {
		return
	}

	nPoint := self.Sql_WeekPlan.Point
	if nPoint == 0 {
		self.Sql_WeekPlan.IsGet = 1
		return
	}

	// 获得奖励配置
	configs := GetCsvMgr().SevendayAward

	index := -1
	for i, config := range configs {
		if nPoint >= config.NeedPoint {
			index = i
		} else {
			break
		}
	}

	if index < 0 {
		return
	}

	if len(configs[index].Items) != len(configs[index].Nums) {
		LogError("len(configs[indxe].Items) = len(configs[indxe].Nums), ", len(configs[index].Items), len(configs[index].Nums))
		return
	}

	// 获得奖励
	var Award []PassItem
	for i := 0; i < len(configs[index].Items); i++ {
		Award = append(Award, PassItem{configs[index].Items[i], configs[index].Nums[i]})
	}

	text := fmt.Sprintf(mailConfig.Mailtxt, nPoint)

	// 发送邮件
	pMail.AddMailWithItems(MAIL_CAN_ALL_GET, mailConfig.Mailtitle, text, Award)

	self.Sql_WeekPlan.IsGet = 1

}

//! 得到任务
func (self *ModWeekPlan) GetTask(taskid int, add bool) *JS_TaskInfo {
	for i := 0; i < len(self.Sql_WeekPlan.taskinfo); i++ {
		if self.Sql_WeekPlan.taskinfo[i].Taskid == taskid {
			return &self.Sql_WeekPlan.taskinfo[i]
		}
	}

	if add {
		var node JS_TaskInfo
		node.Taskid = taskid
		node.Plan = 0
		node.Pickup = 0
		node.Finish = 0
		//node.Tasktypes = GetCsvMgr().WealTask_CSV[taskid].Tasktypes
		node.Tasktypes = GetCsvMgr().SevenDayTask_CSV[taskid].Tasktypes
		self.Sql_WeekPlan.taskinfo = append(self.Sql_WeekPlan.taskinfo, node)
		return &(self.Sql_WeekPlan.taskinfo[len(self.Sql_WeekPlan.taskinfo)-1])
	}

	return nil
}

func (self *ModWeekPlan) HandleTask(tasktype, n2, n3, n4 int) {
	//return
	activity := GetActivityMgr().GetActivity(ACT_SEVEN_DAY)
	if activity == nil {
		return
	}
	rtime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
	firstTime := HF_CalPlayerCreateTime(rtime.Unix(), 0)
	closeTime := firstTime + int64(activity.info.Continued) + int64(-(HF_Atoi(activity.info.Start) + 1)*DAY_SECS)

	if TimeServer().Unix() > closeTime {
		return
	}

	for _, value := range GetCsvMgr().SevenDayTask_CSV {
		if value.Tasktypes != tasktype {
			continue
		}

		node := self.GetTask(value.Id, true)
		if node != nil && node.Finish == 1 {
			continue
		}

		plan, add := DoTask(value, self.player, n2, n3, n4)
		if plan == 0 {
			continue
		}

		if node == nil {
			node = self.GetTask(value.Id, true)
		}

		chg := false
		if add {
			node.Plan += plan
			chg = true
		} else if tasktype == PvpRankNow {
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

		if tasktype == PvpRankNow {
			if node.Plan != 0 && node.Plan <= value.N1 {
				node.Finish = 1
				chg = true
			}
		} else {
			//! 任务完成
			if node.Plan >= value.N1 {
				node.Finish = 1
				node.Plan = value.N1
				chg = true
				v, ok := GetCsvMgr().SevenStatus[value.Id]
				if ok && v == 1 {
					self.Sql_WeekPlan.BoxNum += 1
				}
			}
		}

		if chg {
			self.chg = append(self.chg, *node)
		}
	}
}

//! 获取奖励 ret:1未知任务，2:已经领取 3:时间未到 4:未完成 5:活动过期 6.客户端索引错误
func (self *ModWeekPlan) GetAward(id int) {
	if !self.IsCanPickUp(id) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_MILITARY_ITS_NOT_TIME_TO_COLLECT"))
		return
	}
	info := make([]JS_TaskInfo, 0)

	csv_weal, ok := GetCsvMgr().GetSevenDay(id)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_CONFIG_ERROR"))
		return
	}
	node := self.GetTask(HF_Atoi(csv_weal["id"]), true)
	if node == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_CONFIG_ERROR"))
		return
	}
	if node.Finish == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_CANT_FINISH"))
		return
	}
	if node.Pickup == 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_ALREADY_GET"))
		return
	}

	node.Pickup = 1
	items := make([]int, 0)
	nums := make([]int, 0)
	for i := 0; i < 5; i++ {
		//itemid := HF_Atoi(csv_weal[fmt.Sprintf("reward%d", i+1)])
		itemid := HF_Atoi(csv_weal[fmt.Sprintf("item%d", i+1)])
		if itemid == 0 {
			continue
		}
		//num := HF_Atoi(csv_weal[fmt.Sprintf("itemnum%d", i+1)])
		num := HF_Atoi(csv_weal[fmt.Sprintf("num%d", i+1)])
		items = append(items, itemid)
		nums = append(nums, num)
	}
	getItems := self.player.AddObjectLst(items, nums, "七日活动", 0, 0, 0)
	info = append(info, *node)
	self.Sql_WeekPlan.Point += HF_Atoi(csv_weal["point"])

	var msg S2C_GetWeekPlanAward
	msg.Cid = "getweekplanaward"
	msg.Info = info
	msg.Point = self.Sql_WeekPlan.Point
	msg.GetItem = getItems
	msg.IsGetMark = self.Sql_WeekPlan.isGetMark
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)

	logId := 0
	logDec := ""

	week := (node.Taskid%1000000 - 1) / (10000 * 7)
	if week == 0 {
		logId = LOG_ACTIVITY_SEVEN_1
		logDec = "领取新兵训练奖励"
	} else if week == 1 {
		logId = LOG_ACTIVITY_SEVEN_2
		logDec = "领取半月庆典奖励"
	} else {
		logId = LOG_ACTIVITY_SEVEN_3
		logDec = "领取嘉年华奖励"
	}
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, logId, node.Taskid, 0, 0, logDec, 0, 0, self.player)
}

func (self *ModWeekPlan) GetWeekplanPoint(id int) {

	//先看看是否领取
	_, ok := self.Sql_WeekPlan.isGetMark[id]
	if ok && self.Sql_WeekPlan.isGetMark[id] == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_AWARD_FOR_THE_EVENT_HAS"))
		return
	}

	config := GetCsvMgr().GetSevendayAward(id, self.Sql_WeekPlan.Stage)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_CONFIG_ERROR"))
		return
	}

	if self.Sql_WeekPlan.Point < config.NeedPoint {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STATISTICS_SCORE_NOT_ENOUGH"))
		return
	}

	getItems := self.player.AddObjectLst(config.Items, config.Nums, "七日活动点数奖励", config.NeedPoint, 0, 0)
	self.Sql_WeekPlan.isGetMark[id] = LOGIC_TRUE

	var msg S2C_GetWeekPlanAward
	msg.Cid = "getweekplanpoint"
	msg.Point = self.Sql_WeekPlan.Point
	msg.GetItem = getItems
	msg.IsGetMark = self.Sql_WeekPlan.isGetMark
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)

}

func (self *ModWeekPlan) IsOpen(id int) bool {
	index := (id%1000000)/10000 - 1
	if index < 0 || index >= len(self.WeekPlanStatus) {
		return false
	}
	return self.WeekPlanStatus[index].IsOpen
}

func (self *ModWeekPlan) IsCanPickUp(id int) bool {
	index := (id%1000000)/10000 - 1
	if index < 0 || index >= len(self.WeekPlanStatus) {
		return false
	}
	return self.WeekPlanStatus[index].CanPick
}

func (self *ModWeekPlan) GetNodeOpenTime(id int) int64 {
	index := (id%1000000)/10000 - 1
	if index < 0 || index >= len(self.WeekPlanStatus) {
		return 0
	}
	return self.WeekPlanStatus[index].OpenTime
}

//! 条件购买处理 ret:1未知活动 2:钻石不足 3:未开放 5已过期 4:不能重复购买 6:领取条件不足 7: 8:vip等级不足
func (self *ModWeekPlan) BuyItem(id int) {
	info := make([]JS_TaskInfo, 0)
	//! 检测是否能领取
	if !self.IsOpen(id) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_MILITARY_ITS_NOT_TIME_TO_COLLECT"))
		return
	}

	csv_weal, ok := GetCsvMgr().GetSevenDay(id)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_CONFIG_ERROR"))
		return
	}

	if HF_Atoi(csv_weal["tag_type"]) != TAG_TYPE_BUY {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_CONFIG_ERROR"))
		return
	}

	node := self.GetTask(HF_Atoi(csv_weal["id"]), true)
	if node == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_CONFIG_ERROR"))
		return
	}

	if node.Pickup >= HF_Atoi(csv_weal["times"]) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CITY_SHORTAGE_OF_PURCHASES"))
		return
	}

	if HF_Atoi(csv_weal["vip_limit"]) > 0 {
		if self.player.Sql_UserBase.Vip < HF_Atoi(csv_weal["vip_limit"]) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_INVESTMONEY_LACK_OF_ARISTOCRATIC_RANK"))
			return
		}
	}

	item := make([]int, 0)
	num := make([]int, 0)
	costItem := make([]int, 0)
	costNum := make([]int, 0)
	for i := 1; i < 5; i++ {
		item = append(item, HF_Atoi(csv_weal[fmt.Sprintf("item%d", i)]))
		num = append(num, HF_Atoi(csv_weal[fmt.Sprintf("num%d", i)]))
		costItem = append(costItem, HF_Atoi(csv_weal[fmt.Sprintf("costitem%d", i)]))
		costNum = append(costNum, HF_Atoi(csv_weal[fmt.Sprintf("costnum%d", i)]))
	}
	//检查消耗
	if err := self.player.HasObjectOk(costItem, costNum); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	node.Pickup++

	if HF_Atoi(csv_weal["point"]) > 0 {
		self.Sql_WeekPlan.Point += HF_Atoi(csv_weal["point"])
		// 增加临时表的玩家数据
		//GetWeekPlanMgr().AddPlayer(self.player)
	}
	info = append(info, *node)

	removeitems := self.player.RemoveObjectLst(costItem, costNum, "七日礼包消耗", 0, 0, 0)
	additems := self.player.AddObjectLst(item, num, "七日礼包购买", 0, 0, 0)

	var msg S2C_GetWeekPlanAward
	msg.Cid = "weekplan_buyitem"
	msg.Info = info
	msg.Point = self.Sql_WeekPlan.Point
	msg.GetItem = additems
	msg.CostItem = removeitems
	msg.IsGetMark = self.Sql_WeekPlan.isGetMark
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)

	return
}

//! 商会条件购买
func (self *ModWeekPlan) FundItem(id int) {

	info := make([]JS_TaskInfo, 0)

	if !self.IsOpen(id) {
		LogDebug("活动已经过期")
		return
	}

	csv_weal, ok := GetCsvMgr().GetSevenDay(id)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_CONFIG_ERROR"))
		return
	}

	node := self.GetTask(HF_Atoi(csv_weal["id"]), false)
	if node == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SEVEN_CONFIG_ERROR"))
		return
	}

	if node.Pickup == 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CITY_SHORTAGE_OF_PURCHASES"))
		return
	}

	getItems := make([]PassItem, 0)
	costItems := make([]PassItem, 0)

	item := make([]int, 0)
	num := make([]int, 0)
	costItem := make([]int, 0)
	costNum := make([]int, 0)
	for i := 1; i < 5; i++ {
		item = append(item, HF_Atoi(csv_weal[fmt.Sprintf("item%d", i)]))
		num = append(num, HF_Atoi(csv_weal[fmt.Sprintf("num%d", i)]))
		costItem = append(costItem, HF_Atoi(csv_weal[fmt.Sprintf("costitem%d", i)]))
		costNum = append(costNum, HF_Atoi(csv_weal[fmt.Sprintf("costnum%d", i)]))
	}

	powerok := self.player.GetModule("recharge").(*ModRecharge).IsBuy(WARORDERLIMIT_1)

	rebate_type := HF_Atoi(csv_weal["rebate_type"])
	vipLimit := HF_Atoi(csv_weal["param1"])

	if (rebate_type == 1 && powerok) || (rebate_type == 2 && self.player.Sql_UserBase.Vip >= vipLimit) {
		//! 可以免费领取
		node.Pickup = 1
		getItems = self.player.AddObjectLst(item, num, "七日礼包领取", node.Taskid, 0, 0)
		info = append(info, *node)
	} else { //! 需要购买
		//检查消耗
		if err := self.player.HasObjectOk(costItem, costNum); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
		node.Pickup = 1
		costItems = self.player.RemoveObjectLst(costItem, costNum, "七日礼包领取购买", node.Taskid, 0, 0)
		getItems = self.player.AddObjectLst(item, num, "七日礼包领取", node.Taskid, 0, 0)
		info = append(info, *node)
	}

	if HF_Atoi(csv_weal["point"]) > 0 {
		self.Sql_WeekPlan.Point += HF_Atoi(csv_weal["point"])
		// 增加临时表的玩家数据
		//GetWeekPlanMgr().AddPlayer(self.player)
	}

	var msg S2C_GetWeekPlanAward
	msg.Cid = "weekplan_funditem"
	msg.Info = info
	msg.Point = self.Sql_WeekPlan.Point
	msg.IsGetMark = self.Sql_WeekPlan.isGetMark
	msg.GetItem = getItems
	msg.CostItem = costItems
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("weekplan_funditem", smsg)

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SEVEN_DAY, HF_Atoi(csv_weal["id"]), HF_Atoi(csv_weal["sort"]), HF_Atoi(csv_weal["step"]), "七日活动", 0, 0, self.player)
}

///////////////////////
func (self *ModWeekPlan) SendInfo() {
	self.InitWeekPlanStatus()
	var msg S2C_WeekPlanInfo
	msg.Cid = "weekplaninfo"
	t, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
	msg.Info.Regtime = t.Unix()
	msg.Info.Servercurtime = TimeServer().Unix()
	msg.Info.Taskinfo = self.Sql_WeekPlan.taskinfo
	msg.Info.TaskStatus = self.WeekPlanStatus
	msg.Point = self.Sql_WeekPlan.Point
	msg.IsGetMark = self.Sql_WeekPlan.isGetMark
	msg.BuyTime = self.Sql_WeekPlan.Lastupdtime
	msg.Config = GetCsvMgr().SevendayConfig
	msg.Stage = self.Sql_WeekPlan.Stage
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("weekplaninfo", smsg)

	self.init = true

	//TODO 屏蔽七日宝箱
	//self.checkBoxItem()
}

func (self *ModWeekPlan) UpdateBuyTime() {
	self.Sql_WeekPlan.Lastupdtime = TimeServer().Unix()
	self.SendInfo()
}

func (self *ModWeekPlan) SendUpdate() {
	if !self.init {
		return
	}

	if len(self.chg) == 0 {
		return
	}

	var msg S2C_TaskUpdate
	msg.Cid = "weekplan_update"
	msg.Info = self.chg
	self.chg = make([]JS_TaskInfo, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}

// 计算宝箱道具数量
func (self *ModWeekPlan) checkBoxItem() {
	//if !GetWeekPlanMgr().isTimeOver() {
	//	return
	//}
	//
	//if self.Sql_WeekPlan.BoxNum <= 0 {
	//	return
	//}
	//
	//boxConfig := GetTowerMgr().GetBoxConfig(41)
	//if boxConfig == nil {
	//	LogError("BoxItem 41 is nil")
	//	return
	//}
	//
	//if len(boxConfig.Item) != len(boxConfig.Num) {
	//	LogError("len(boxConfig.Item) != len(boxConfig.Num)")
	//	return
	//}
	//
	//title := GetCsvMgr().GetText("STR_SEVEN_TITLE")
	//text := GetCsvMgr().GetText("STR_SEVEN_CONTENT")
	//pMail := self.player.GetModule("mail").(*ModMail)
	//var res []PassItem
	//for i := range boxConfig.Item {
	//	res = append(res, PassItem{boxConfig.Item[i], boxConfig.Num[i] * self.Sql_WeekPlan.BoxNum})
	//}
	//
	//if pMail != nil {
	//	pMail.AddMail(1, 1, 0, title, text, GetCsvMgr().GetText("STR_SYS"), res, true, 0)
	//}
	//
	//self.Sql_WeekPlan.BoxNum = 0

}

//func (self *ModWeekPlan) CheckTaskDone()  {
//	for _, pTask := range self.chg {
//		config := GetCsvMgr().SevenDay_CSV[pTask.Taskid]
//		if config == nil {
//			continue
//		}
//
//		if HF_Atoi(config["tasktypes"]) == HeroUpStarNumTask {
//			if pTask.Finish == 1 {
//				continue
//			}
//
//			if pTask.Pickup == 1 {
//				continue
//			}
//
//			pass := self.player.GetModule("pass").(*ModPass).GetPass(passId)
//			if pass != nil {
//				self.chg.taskinfo[i].Finish = 1
//			}
//
//		}
//	}
//}
