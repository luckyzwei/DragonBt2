package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	CS_ACTIVITY_FUND_GET_INFO  = "activity_fund_get_info"
	CS_ACTIVITY_FUND_GET_AWARD = "activity_fund_get_award"
	CS_ACTIVITY_FUND_ACTIVATE  = "activity_fund_get_activate"
)

const (
	ACTIVITY_MODE_LIMIT   = 3 //限时活动
	ACTIVITY_MODE_LIMIT_4 = 4 //
)

const (
	ACTIVITY_MODE_TYPE_2 = 2
	ACTIVITY_MODE_TYPE_3 = 3
)

const (
	GemCostReason         = 1
	OnRefrsh              = 2
	CreateAct             = 3
	TimeOut               = 4
	CdClear               = 5
	MonthCard1            = 100901
	MonthCard2            = 100902
	ForeverCard1          = 100903
	ForeverCard2          = 100904
	GoldWeek              = 151 //黄金周卡
	ActSingleRecharge     = 9013
	ActTimeGift           = 9021 //限时礼包
	ActivityValueFundType = 9023 //超值基金 根据活动配置 配置的七天能买
	ActivityOpenFundType  = 9026 //开服基金
	TotleAwardType        = 2

	// 需要特殊处理的新活动 上面的废弃

	ACT_GROWTH_GIFT        = 1001 // 成长礼包
	ACT_DAY_GIFT           = 1002 // 日礼包
	ACT_WEEK_GIFT          = 1003 // 周礼包
	ACT_MONTH_GIFT         = 1004 // 月礼包
	ACT_WARORDER_1         = 1005 // 皇家犒赏令
	ACT_WARORDER_2         = 1006 // 勇者犒赏令
	ACT_NOVICE_GIFT        = 1008 // 新手礼包
	ACT_SEVEN_DAY          = 1011 // 七日礼包
	ACT_TURNTABLE          = 1016 // 转盘活动
	ACT_WARORDERLIMIT_1    = 1020 // 主线战令
	ACT_WARORDERLIMIT_2    = 1021 // 爬塔战令
	ACT_WARORDERLIMIT_3    = 1027 // 钻石累消战令
	ACT_SEASON_ACTIVITY    = 1023 // 特殊节日活动，送礼物，开放特殊礼包，做成可重置
	ACT_GENERAL            = 1024 // 魔法的跨服神将
	MAGICHERO_ACTIVITY_ID  = 1026 // 神将开启活动
	ACT_NEW_FUND           = 1040 // 无敌理财
	ACT_LUCKY_FIND         = 1041 // 福袋召唤
	ACT_OVERFLOW_GIFT_MIN      = 1401 // 超值好礼
	ACT_OVERFLOW_GIFT_MAX      = 1499 // 超值好礼
	ACT_STAR_GIFT_MIN          = 1500 // 星辰礼包
	ACT_STAR_GIFT_MAX          = 1505 // 星辰礼包
	ACT_ACCESSCARD_MIN         = 1600 // 收藏家
	ACT_ACCESSCARD_MAX         = 1610 // 收藏家
	ACT_FUND_MIN               = 1611 // 超值基金
	ACT_FUND_MAX               = 1620 // 超值基金
	ACT_NEWPIT_HALF_MIN        = 1621 // 地牢异变
	ACT_NEWPIT_HALF_MAX        = 1625 // 地牢异变
	ACT_BOSS_MIN               = 1700 // 暗域入侵
	ACT_BOSS_MAX               = 1709 // 暗域入侵
	ACT_AREAN_CROSS_SERVER     = 1730 // 跨服竞技
	ACT_AREAN_CROSS_SERVER_3V3 = 1731 // 跨服竞技3v3
	ACT_RANKREWARD_COST        = 1732 // 消费排行榜
	ACT_BOSS_FESTIVAL          = 1750 // 节日BOSS
	ACT_HERO_GROW_MIN          = 1800 // 英雄成长礼包
	ACT_HERO_GROW_MAX          = 1807 // 英雄成长礼包
	ACT_LOTTERY_DRAW           = 1850 // 抽奖
	ACT_REG_TOTAL_DAY          = 2001 // 连续登陆
	ACT_REG_HERO_GET           = 2002 // 英雄集结
	ACT_FESTIVAL_EXCHANGE      = 2201 // 英雄集结

	ActivityLuckShop             = 1007 //限时礼包 // 废弃 为了不报错而暂时保留
	ActivityStarGift             = 1013 //星辰限时礼包
	ActivityWaveBless            = 1014 //怒涛祝福礼包
	ActivityDiscountGift         = 1019 //特惠礼包
	ActivityGiftEx               = 1025 //礼包活动
	ACT_TIME_LIMIT_GIFT_START    = 1100 //限时礼包起始
	ACT_TIME_LIMIT_GIFT_END      = 1200 //限时礼包结束
	ACT_TIME_LIMIT_GIFT_EX_START = 5001 //限时礼包扩展开始
	ACT_TIME_LIMIT_GIFT_EX_END   = 5999 //限时礼包扩展结束

	ACT_ONHOOK_ACTIVITY_SPRING_FESTIVAL_LIVENESS = 6020 //春节活跃度活动
	ACT_ONHOOK_ACTIVITY_SPRING_FESTIVAL          = 6100 //春节掉落活动

	ACT_ELITE_UPGRADE           = 4110 // 精英直升
)

const (
	BECOME_STRONGER_HERO = 1 //我要变强领奖条件 拥有某英雄
)

type ActivityFund struct {
	Pay             int    `json:"pay"`
	Ver             int    `json:"ver"`
	ActivityEndTime int64  `json:"accendt"`
	StartTime       int64  `json:"starttime"`
	FundGetType     [7]int `json:"gettype"`
}

//! 活动数据库
type San_Activity struct {
	Uid               int64
	Info              string
	JJ                string
	Month             string
	BecomeStronger    string
	Fund              string
	ActivityResetSign string


	info              map[int]*JS_Activity //! 活动状态, key:活动唯一Id, value:活动信息,ex.{"100101":{"id":100101,"progress":0,"time":1527926667,"done":0}
	jj                []int                //! 购买基金
	fund              map[int][2]*ActivityFund
	month             []JS_MonthCard             //! 月卡信息
	becomeStronger    map[int]*JS_BecomeStronger //! 我要变强领取状态
	activityResetSign map[int]int64              //!

	DataUpdate
}

// 我要变强领取状态
type JS_BecomeStronger struct {
	Id    int `json:"id"`     //! Id
	IsGet int `json:"is_get"` //! 是否领取
}

// 活动Id, 进度, 结束时间, 完成状态, Done=0,表示未完成, Done=1表示已完成, Done=2表示已领取
type JS_Activity struct {
	Id       int   `json:"id"`
	Progress int   `json:"progress"`
	Time     int64 `json:"time"`
	Done     int   `json:"done"`
	Step     int   `json:"step"` //! 第几天
	Ver      int   `json:"ver"`  //! 活动版本
}

// 月卡信息
type JS_MonthCard struct {
	Id         int         `json:"id"`         //! 月卡Id
	StartTime  int64       `json:"starttime"`  //! 开始时间
	Day        int         `json:"day"`        //! 有效天数
	Get        int         `json:"get"`        //! 领取天数
	RewardSign map[int]int `json:"rewardsign"` //! 领取标记
	Stage      int         `json:"stage"`      //! 期数
}

// 活动Mod管理
type ModActivity struct {
	player         *Player
	Sql_Activity   San_Activity  //! 数据库结构
	LastUpdateTime int64         //! 上次更新时间
	DataLocker     *sync.RWMutex //! 多线程处理
}

// 加载玩家活动数据
func (self *ModActivity) OnGetData(player *Player) {
	self.player = player
	self.DataLocker = new(sync.RWMutex)

	self.Sql_Activity.info = make(map[int]*JS_Activity)
	sql := fmt.Sprintf("select * from `san_activity` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Activity, "san_activity", self.player.ID)

	if self.Sql_Activity.Uid <= 0 {
		// 插入玩家活动信息
		self.Sql_Activity.Uid = self.player.ID
		self.Sql_Activity.info = make(map[int]*JS_Activity)
		self.Sql_Activity.jj = make([]int, 0)
		self.Sql_Activity.fund = make(map[int][2]*ActivityFund)
		self.Sql_Activity.month = make([]JS_MonthCard, 0)
		self.Sql_Activity.becomeStronger = make(map[int]*JS_BecomeStronger)
		self.Encode()
		InsertTable("san_activity", &self.Sql_Activity, 0, true)
	} else {
		self.Decode()
	}

	if self.Sql_Activity.activityResetSign == nil {
		self.Sql_Activity.activityResetSign = make(map[int]int64)
	}
	// 活动数据库逻辑初始化, redis缓存以及mysql表设置
	self.Sql_Activity.Init("san_activity", &self.Sql_Activity, true)
}

// 获取玩家其他信息
func (self *ModActivity) OnGetOtherData() {
	self.ReloadActivity()
	self.HandleTask(0, 0, 0, 0)
	self.HandleTask(TASK_TYPE_VIP_BUY, 0, 0, 0)
	self.LastUpdateTime = TimeServer().Unix()
}

//! 重新载入活动
func (self *ModActivity) ReloadActivity() {
	nowtime := TimeServer().Unix()
	// 获取所有活动类型
	lst := GetActivityMgr().GetActivityType()
	// 检查活动状态
	for i := 0; i < len(lst); i++ {
		globalAct := GetActivityMgr().GetActivity(lst[i])
		if globalAct == nil {
			continue
		}

		//fmt.Printf("globalAct.info.Type----------------%+v\n", globalAct.info.Type)
		//! 单笔充值一直开放
		if globalAct.info.Type == ActSingleRecharge || globalAct.info.Type == ActTimeGift {
			globalAct.info.Status = ACTIVITY_STATUS_OPEN
			globalAct.status.Status = ACTIVITY_STATUS_OPEN
			globalAct.status.EndTime = nowtime + 864000000
			continue
		}

		if globalAct.info.Type == ActivityOpenFundType {
			globalAct.info.Status = ACTIVITY_STATUS_OPEN
			globalAct.status.Status = ACTIVITY_STATUS_OPEN
			globalAct.status.EndTime = int64(globalAct.info.Continued) + int64(globalAct.info.Show)
			continue
		}

		if globalAct.info.Type == ACT_REG_TOTAL_DAY { //|| globalAct.info.Type == ACT_GROWTH_GIFT {
			globalAct.info.Status = ACTIVITY_STATUS_OPEN
			globalAct.status.Status = ACTIVITY_STATUS_OPEN
			globalAct.status.EndTime = int64(globalAct.info.Continued) + int64(globalAct.info.Show)
			continue
		}

		// 活动处于开启状态
		if globalAct.status.Status > ACTIVITY_STATUS_CLOSED {
			// 道具以及活动Id信息等
			for j := 0; j < len(globalAct.items); j++ {
				self.DataLocker.RLock()
				node, ok := self.Sql_Activity.info[globalAct.items[j].Id]
				if ok == false {
					self.DataLocker.RUnlock()
					// 初始化活动状态
					node = new(JS_Activity)
					node.Id = globalAct.items[j].Id
					node.Progress = 0
					node.Done = 0
					node.Time = globalAct.status.EndTime
					self.DataLocker.Lock()
					self.Sql_Activity.info[globalAct.items[j].Id] = node
					self.DataLocker.Unlock()
					node.reset(CreateAct)
				} else {
					self.DataLocker.RUnlock()
					//! 活动已过期，则刷新活动
					/*
						if nowtime > node.Time {
							if globalAct.info.TaskType != 800 {
								node.Progress = 0
								node.Done = 0
								node.Time = globalAct.status.EndTime
								node.reset(TimeOut)
							}
						}
					*/
				}
			}
		} else {
			//! 活动结束，删除活动
			self.DataLocker.Lock()

			self.CheckAward(globalAct)
			for j := 0; j < len(globalAct.items); j++ {
				//! 判断是否领取，未领取则发送邮件
				delete(self.Sql_Activity.info, globalAct.items[j].Id)
			}
			self.DataLocker.Unlock()
		} //end if act close
	}
}

//! 检测奖励领取，过期的时候，通过邮件发送
func (self *ModActivity) CheckAward(actType *Sql_ActivityMask) {
	if actType.info.Type == ACT_DAILY_RECHARGE {
		for j := 0; j < len(actType.items); j++ {
			//! 判断是否领取，未领取则发送邮件
			actInfo, ok := self.Sql_Activity.info[actType.items[j].Id]
			if ok {
				if actInfo.Done == 1 {
					mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_GET_DAILY_RECHARGE]
					if ok {
						//title := actType.info.Name
						//text := fmt.Sprintf(GetCsvMgr().GetText("STR_ACT_NOT_GET_AWARD"))
						/*
							itemLst := []PassItem{}
							for k := 0; k < len(actType.items[j].Item); k++ {
								if actType.items[j].Item[k] > 0 {
									itemLst = append(itemLst, PassItem{
										ItemID: actType.items[j].Item[k],
										Num:    actType.items[j].Num[k],
									})
								}
							}
						*/
						itemLst := self.GetDailyItem()
						self.sendMail(mailConfig.Mailtitle, mailConfig.Mailtxt, itemLst)
					}
					actInfo.Done = 2
				}
			}
		}
	}
}

func (self *ModActivity) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	// 限时兑换
	case "exchangeitem":
		var msg C2S_FinishActivity
		json.Unmarshal(body, &msg)
		self.ExChangeItem(msg.Id)
		return true
		// 领取奖励
	case "finishactivity":
		var msg C2S_FinishActivity
		json.Unmarshal(body, &msg)
		self.FinishActivity(msg.Id)
		return true
		// 购买基金
	case "buyjj":
		var msg C2S_FinishActivity
		json.Unmarshal(body, &msg)
		self.BuyJJ(msg.Id)
		return true
		// 我要变强活动领取
	case "promotebox":
		var msg C2S_GetPromoteBox
		json.Unmarshal(body, &msg)
		self.GetPromoteBox(msg.Id)
		return true
		// 获取活动信息
	case "getactivitymask":
		var msg C2S_FinishActivity
		json.Unmarshal(body, &msg)
		self.GetActivityMask(msg.Actver)
		return true
	case "dailydiscount":
		var msg C2S_DailyDiscount
		json.Unmarshal(body, &msg)
		self.dailyDiscount(msg.Id, ctrl)
		return true
	case "getmonthcard":
		self.MonthCardCumulativeReward(false)
		self.SendMonthCard()
		return true
	case "getcumulativereward":
		var msg C2S_GetCumulativeReward
		json.Unmarshal(body, &msg)
		self.GetMonthCardCumulativeRewardNew(msg.Id)
		return true
	case CS_ACTIVITY_FUND_GET_INFO:
		var msg C2S_ActivityFundInfo
		json.Unmarshal(body, &msg)
		self.SendActivityFundInfo(msg.Ver)
		return true
	case CS_ACTIVITY_FUND_GET_AWARD:
		var msg C2S_ActivityFundAward
		json.Unmarshal(body, &msg)
		self.GetActivityFundAward(msg.Ver, msg.Pay, msg.Day)
		return true
	case "getoverflowgift":
		var msg C2S_GetOverflow
		json.Unmarshal(body, &msg)
		self.GetOverflow(msg.Id, ctrl)
		return true
	}

	return false
}

// 存盘逻辑
func (self *ModActivity) OnSave(sql bool) {
	self.Encode()
	self.Sql_Activity.Update(sql)
}

//! 序列化
func (self *ModActivity) Decode() {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	json.Unmarshal([]byte(self.Sql_Activity.Info), &self.Sql_Activity.info)
	json.Unmarshal([]byte(self.Sql_Activity.JJ), &self.Sql_Activity.jj)
	json.Unmarshal([]byte(self.Sql_Activity.Fund), &self.Sql_Activity.fund)
	json.Unmarshal([]byte(self.Sql_Activity.Month), &self.Sql_Activity.month)
	json.Unmarshal([]byte(self.Sql_Activity.BecomeStronger), &self.Sql_Activity.becomeStronger)
	json.Unmarshal([]byte(self.Sql_Activity.ActivityResetSign), &self.Sql_Activity.activityResetSign)
}

//! 反序列化
func (self *ModActivity) Encode() {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	self.Sql_Activity.Info = HF_JtoA(self.Sql_Activity.info)
	self.Sql_Activity.JJ = HF_JtoA(self.Sql_Activity.jj)
	self.Sql_Activity.Fund = HF_JtoA(self.Sql_Activity.fund)
	self.Sql_Activity.Month = HF_JtoA(self.Sql_Activity.month)
	self.Sql_Activity.BecomeStronger = HF_JtoA(self.Sql_Activity.becomeStronger)
	self.Sql_Activity.ActivityResetSign = HF_JtoA(self.Sql_Activity.activityResetSign)
}

// 获取单个活动信息
func (self *ModActivity) GetActivity(id int) *JS_Activity {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	return self.Sql_Activity.info[id]
}

// 是否是购买基金活动
func (self *ModActivity) IsBuyJJ(id int) bool {
	for i := 0; i < len(self.Sql_Activity.jj); i++ {
		if self.Sql_Activity.jj[i] == id {
			return true
		}
	}
	return false
}

// 获取月卡剩余天数
func (self *ModActivity) GetMonthCardDay(moneytype int) int {
	nowtime := TimeServer().Unix()
	for i := 0; i < len(self.Sql_Activity.month); i++ {
		if self.Sql_Activity.month[i].Id == moneytype {
			allday := self.Sql_Activity.month[i].Day
			useday := int((nowtime - self.Sql_Activity.month[i].StartTime) / 86400)
			return allday - useday
		}
	}

	return 0
}

// 领取宝箱奖励
func (self *ModActivity) GetPromoteBox(id int) {
	if len(self.Sql_Activity.becomeStronger) <= 0 {

		self.Sql_Activity.becomeStronger = make(map[int]*JS_BecomeStronger)
	}
	_, ok := self.Sql_Activity.becomeStronger[id]
	if ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_TREASURE_BOX_CAN_NOT"))
		return
	}

	config, ok := GetCsvMgr().BecomeStrongerMap[id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	if config.ConditionType == BECOME_STRONGER_HERO {
		hero := self.player.GetModule("hero").(*ModHero).GetHero(HF_Atoi(config.Condition))
		if hero == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FATE_HEROES_DO_NOT_EXIST"))
			return
		}
	} else {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	BoxConfig, ok := GetCsvMgr().LevelboxConfig[config.Boxid]
	if !ok {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_UNION_STATE_TREASURE_BOX_CONFIGURATION_DOES"))
		return
	}

	if len(BoxConfig.Items) <= 0 || len(BoxConfig.Items) != len(BoxConfig.Nums) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	//! 添加物品
	out := make([]PassItem, 0)

	for index, _ := range BoxConfig.Items {
		out = append(out, PassItem{BoxConfig.Items[index], BoxConfig.Nums[index]})
	}

	for i := 0; i < len(out); i++ {
		out[i].ItemID, out[i].Num = self.player.AddObject(out[i].ItemID, out[i].Num, 0, 0, 0, "活动宝箱领取")
	}
	//! 修改数据
	self.Sql_Activity.becomeStronger[id] = &JS_BecomeStronger{id, 1}
	//! 发消息
	var msg S2C_GetPromoteBox
	msg.Cid = "promotebox"
	msg.Id = id
	msg.Item = out
	self.player.SendMsg("promotebox", HF_JtoB(&msg))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_GET_BOX, 0, 0, 0, "活动宝箱领取", 0, 0, self.player)
}

// 特殊处理 是否有前置条件
func (self *ModActivity) HaveFirstCondition(act_type int) bool {
	if act_type == ActivityWaveBless {
		//for _, v := range self.player.GetModule("recharge").(*ModRecharge).Sql_UserRecharge.warOrder {
		//	if v.BuyState == LOGIC_TRUE {
		//		return true
		//	}
		//}
		//
		//return false
	}

	return true
}

//! 活动任务触发
func (self *ModActivity) HandleTask(condition int, n1 int, n2 int, n3 int) bool {
	//LogDebug("condition:", condition, ", n1:", n1, ", n2:", n2, ", n3:", n3)
	//检查服务器状态
	out := make([]JS_Activity, 0)
	lst := GetActivityMgr().GetActivityType()
	for i := 0; i < len(lst); i++ {
		// 获取活动状态
		act_type := GetActivityMgr().GetActivity(lst[i])
		// 活动如果开启
		if act_type != nil && act_type.status.Status > ACTIVITY_STATUS_CLOSED {

			startday := HF_Atoi(act_type.info.Start)
			if startday < 0 {
				if act_type.info.Type == ACT_REG_TOTAL_DAY {
					if !self.isActivityRegTotalDayOpen() {
						continue
					}
					//} else if act_type.info.Type == ACT_GROWTH_GIFT {
					//if !self.isActivityGrowthGiftOpen() {
					//	continue
					//}
				} else {
					rtime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
					firstTime := HF_CalPlayerCreateTime(rtime.Unix(), 0)
					closeTime := firstTime + int64(act_type.info.Continued) + int64(-(HF_Atoi(act_type.info.Start) + 1)*DAY_SECS)

					if TimeServer().Unix() > closeTime {
						continue
					}
				}

			}
			// 特殊处理 是否有前置条件
			if !self.HaveFirstCondition(act_type.Id) {
				continue
			}

			if act_type.info.Mode == ACTIVITY_MODE_LIMIT {
			}
			// 过滤任务类型
			if act_type.info.TaskType == condition || condition == 0 {
				// 道具信息遍历
				for j := 0; j < len(act_type.items); j++ {
					self.DataLocker.RLock()
					// 读取并插入数据
					actitem, ok := self.Sql_Activity.info[act_type.items[j].Id]
					if !ok {
						self.DataLocker.RUnlock()
						// 创建新的活动
						newinfo := new(JS_Activity)
						newinfo.Id = act_type.items[j].Id
						newinfo.Progress = 0
						newinfo.Done = 0
						newinfo.Time = act_type.status.EndTime
						self.DataLocker.Lock()
						self.Sql_Activity.info[act_type.items[j].Id] = newinfo
						self.DataLocker.Unlock()
						actitem = newinfo
					} else {
						self.DataLocker.RUnlock()
					}

					// 过滤已经完成和领取奖励的任务
					if actitem.Done != 0 {
						// 循环任务需要记录进度
						if !self.IsNeedSavePlan(act_type) {
							continue
						}
					}

					// 创建任务节点
					var tasknode TaskNode
					tasknode.Tasktypes = act_type.info.TaskType
					tasknode.N1 = act_type.items[j].N[0]
					tasknode.N2 = act_type.items[j].N[1]
					tasknode.N3 = act_type.items[j].N[2]
					tasknode.N4 = act_type.items[j].N[3]

					// 检查记次和记值任务
					process, add := DoTask(&tasknode, self.player, n1, n2, n3)

					if process == 0 && condition != TASK_TYPE_VIP_BUY {
						continue
					}

					if add {
						// 如果是记次任务
						actitem.Progress += process
						actitem.Time = act_type.status.EndTime
					} else if act_type.info.TaskType == PvpRankNow {
						if process != 0 { // 新排名为0则直接不处理
							if actitem.Progress == 0 { // 进度为0则说明未初始化 直接赋值
								actitem.Progress = process
							} else { // 进度不为0 则需要判断 获得的新名次比之前要高 则赋值
								if process < actitem.Progress {
									actitem.Progress = process
								}
							}
						}
					} else {
						// 如果是记值任务
						actitem.Progress = HF_MaxInt(process, actitem.Progress)
						actitem.Time = act_type.status.EndTime
					}

					if act_type.info.TaskType == PvpRankNow {
						if actitem.Progress <= act_type.items[j].N[0] {
							actitem.Progress = act_type.items[j].N[0]
							if actitem.Done == 0 {
								actitem.Done = 1
							}
						}
					} else {
						// 判断活动是否完成,n0表示任务的最大值
						if actitem.Progress >= act_type.items[j].N[0] {
							if !self.IsNeedSavePlan(act_type) {
								actitem.Progress = act_type.items[j].N[0]
							}
							if actitem.Done == 0 {
								actitem.Done = 1
								if act_type.Id == ACT_ELITE_UPGRADE{
									actitem.Time = TimeServer().Unix()
								}
							}
						}

						//限时的活动奖励直接发邮件
						if act_type.info.Mode == ACTIVITY_MODE_LIMIT && act_type.info.ModeType == 2 && actitem.Done == 1 {
							activityitem := GetActivityMgr().GetActivityItem(actitem.Id)
							if activityitem != nil {
								actitem.Done = 2
								item := make([]PassItem, 0)
								for j := 0; j < len(activityitem.Item); j++ {
									itemid := activityitem.Item[j]
									if itemid == 0 {
										break
									}
									item = append(item, PassItem{itemid, activityitem.Num[j]})
								}

								mailConfig, ok := GetCsvMgr().MailConfig[19]
								pMail := self.player.GetModule("mail").(*ModMail)
								text := fmt.Sprintf(mailConfig.Mailtxt, act_type.info.Name)
								title := fmt.Sprintf(mailConfig.Mailtitle, act_type.info.Name)
								if ok && pMail != nil {
									// 发送邮件
									pMail.AddMailWithItems(MAIL_CAN_ALL_GET, title, text, item)
								}
							}
						}
					}

					// 任务奖励
					out = append(out, *actitem)
					// 匪夷所思,结束时间-当前时间=持续时间
					//if out[len(out)-1].Time > 0 {
					//	out[len(out)-1].Time -= TimeServer().Unix()
					//}
				}
			}
		}
	}

	// 活动完成时发送信息, 活动或者月卡
	if len(out) > 0 {
		var msg S2C_ActivityUpd
		msg.Cid = "activityupd"
		msg.Info = out
		msg.Month = self.Sql_Activity.month
		self.player.SendMsg("1", HF_JtoB(&msg))
	}

	return len(out) > 0
}

//! 取得全局活动状态
func (self *ModActivity) GetActivityMask(ver int) {
	self.player.CheckRefresh()
	GetActivityMgr().SendInfo(ver, self.player)

	lst := GetActivityMgr().GetActivityType()
	for i := 0; i < len(lst); i++ {
		if self.IsLimitGiftType(lst[i]) {
			// 获取活动状态
			act_type := GetActivityMgr().GetActivity(lst[i])
			// 活动如果关闭
			if act_type.status.Status == ACTIVITY_STATUS_CLOSED {
				// 道具信息遍历
				for j := 0; j < len(act_type.items); j++ {
					actitem, _ := self.Sql_Activity.info[act_type.items[j].Id]
					if actitem != nil {
						actitem.Done = 0
						actitem.Progress = 0
					}
				}
			}
		}
	}
}

//! 限时兑换活动
func (self *ModActivity) ExChangeItem(id int) {

	//defer self.DataLocker.RUnlock()
	self.DataLocker.RLock()
	actitem, ok := self.Sql_Activity.info[id]
	if !ok {
		self.DataLocker.RUnlock()
		var msg S2C_ActivityExChangeItem
		msg.Cid = "exchangeitem"
		msg.Id = 0
		msg.Item = make([]PassItem, 0)
		msg.Cost = make([]PassItem, 0)
		self.player.SendMsg("1", HF_JtoB(&msg))

		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THERE_IS_NO_TIME_LIMIT"))
		return
	}
	self.DataLocker.RUnlock()

	activitytype := GetActivityMgr().GetActivity(actitem.Id)
	activityitem := GetActivityMgr().GetActivityItem(actitem.Id)

	if activityitem.N[0] <= actitem.Progress {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_AWARD_FOR_THE_EVENT_HAS"))
		return
	}
	// 活动异常检查
	if activitytype == nil || activityitem == nil || activitytype.status.Status == ACTIVITY_STATUS_CLOSED {
		return
	}

	// 限时兑换任务类型
	if activitytype.info.TaskType != TASK_TYPE_ACTIVITY_EXCHANGE {
		return
	}

	out := make([]PassItem, 0)
	for i := 0; i < len(activityitem.Item); i++ {
		itemid := activityitem.Item[i]
		if itemid == 0 {
			break
		}
		out = append(out, PassItem{itemid, activityitem.Num[i]})
	}

	cost := make([]PassItem, 0)
	for i := 0; i < 4; i++ {
		itemid := activityitem.CostItem[i]
		if itemid == 0 {
			break
		}
		// 检查道具是否充足
		if self.player.GetObjectNum(itemid) < activityitem.CostNum[i] {
			var msg S2C_ActivityExChangeItem
			msg.Cid = "exchangeitem"
			msg.Id = 0
			msg.Progress = 0
			msg.Done = 0
			msg.Item = out
			msg.Cost = cost
			self.player.SendMsg("1", HF_JtoB(&msg))

			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_INADEQUATE_EXCHANGE_CONDITIONS"))
			return
		}
		cost = append(cost, PassItem{itemid, activityitem.CostNum[i]})
	}

	for i := 0; i < len(out); i++ {
		if out[i].ItemID == 91000003 {
			out[i].Num *= 2
			//if self.player.GetPower() >= POWERMAX {
			//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_POWERMAX"))
			//	return
			//}
		}
	}

	for i := 0; i < len(out); i++ {
		out[i].ItemID, out[i].Num = self.player.AddObject(out[i].ItemID, out[i].Num, activitytype.info.Type, 0, 0, "活动")
	}
	for i := 0; i < len(cost); i++ {
		cost[i].ItemID, cost[i].Num = self.player.AddObject(cost[i].ItemID, cost[i].Num*-1, activitytype.info.Type, 0, 0, "活动")
	}
	actitem.Progress += 1

	var msg S2C_ActivityExChangeItem
	msg.Cid = "exchangeitem"
	msg.Id = id
	msg.Progress = actitem.Progress
	msg.Done = actitem.Done
	msg.Item = out
	msg.Cost = cost
	self.player.SendMsg("1", HF_JtoB(&msg))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_ACTIVITY, activitytype.info.Type, activitytype.info.Mode, activitytype.info.ModeType, "活动", 0, 0, self.player)

}

//! 领取活动奖励
func (self *ModActivity) FinishActivity(id int) {
	doubleitem, doubleexp := GetActivityMgr().GetDoubleStatus(DOUBLE_POWER)
	LogDebug("体力双倍倍率:", doubleitem, doubleexp)

	self.DataLocker.RLock()
	actitem, ok := self.Sql_Activity.info[id]
	if !ok {
		self.DataLocker.RUnlock()
		var msg S2C_ActivityGet
		msg.Cid = "activityget"
		msg.Id = 0
		msg.Item = make([]PassItem, 0)
		self.player.SendMsg("1", HF_JtoB(&msg))

		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_ACTIVITIES_DO_NOT_EXIST_RECEIVE"))
		return
	}
	self.DataLocker.RUnlock()

	// 任务未完成
	if actitem.Done == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_ACTIVITY_WAS_NOT_COMPLETED"))
		return
	}

	// 奖励已领取
	if actitem.Done == 2 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_AWARD_FOR_THE_EVENT_HAS"))
		return
	}

	activitytype := GetActivityMgr().GetActivity(actitem.Id)
	activityitem := GetActivityMgr().GetActivityItem(actitem.Id)
	if activitytype == nil || activityitem == nil || activitytype.status.Status == ACTIVITY_STATUS_CLOSED {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_ACTIVITY_DOES_NOT_EXIST"))
		return
	}

	//检查消耗够不够
	if err := self.player.HasObjectOk(activityitem.CostItem, activityitem.CostNum); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	// 扣除物品
	costItem := self.player.RemoveObjectLst(activityitem.CostItem, activityitem.CostNum, "活动购买", activityitem.Id, 0, 1)

	out := make([]PassItem, 0)
	if activitytype.Id == ACT_DAILY_RECHARGE {
		out = self.GetDailyItem()
	} else if activitytype.Id == ACT_ELITE_UPGRADE{
		for j := 0; j < len(activityitem.Item); j++ {
			itemid := activityitem.Item[j]
			itemnum := activityitem.Num[j]
			if itemid == 0 {
				break
			}
			if itemid == 92000041 {
				itemnum *= self.player.Sql_UserBase.Level	// 经验奖励为经验*玩家等级
			}
			out = append(out, PassItem{itemid, itemnum})
		}
	} else{
		for j := 0; j < len(activityitem.Item); j++ {
			itemid := activityitem.Item[j]
			if itemid == 0 {
				break
			}
			out = append(out, PassItem{itemid, activityitem.Num[j]})
		}
	}

	for j := 0; j < len(out); j++ {
		if out[j].ItemID == 91000003 {
			out[j].Num *= 2
		}
	}
	num := GetGemNum(costItem)
	var item []PassItem
	param3 := 0
	if num > 0 {
		param3 = -1
	}
	day := 0
	if id == MonthCard1 {
		day = self.player.GetModule("activity").(*ModActivity).GetMonthCardDay(101)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD_GET, 1, day, 0, "领取普通月卡奖励", 0, 0, self.player)
	} else if id == MonthCard2 {
		day = self.player.GetModule("activity").(*ModActivity).GetMonthCardDay(102)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD_HIGH_GET, 1, day, 0, "领取高级月卡奖励", 0, 0, self.player)
	} else if id == ForeverCard1 {
		day = self.player.GetModule("activity").(*ModActivity).GetMonthCardDay(103)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD_HIGH_GET, 1, day, 0, "领取高级月卡奖励", 0, 0, self.player)
	} else if id == ForeverCard2 {
		day = self.player.GetModule("activity").(*ModActivity).GetMonthCardDay(104)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD_HIGH_GET, 1, day, 0, "领取高级月卡奖励", 0, 0, self.player)
	}
	for j := 0; j < len(out); j++ {
		out[j].ItemID, out[j].Num = self.player.AddObject(out[j].ItemID, out[j].Num, activitytype.info.Type, day, param3, "活动购买")

		item = append(item, PassItem{out[j].ItemID, out[j].Num})
	}

	if num > 0 {
		AddSpecialSdkItemListLog(self.player, num, item, "活动购买")
	}

	//循环任务
	actitem.Done = 2
	isNeed := LOGIC_TRUE
	actType := actitem.Id / 100
	updateInfo := make([]*JS_Activity, 0)
	if activitytype.info.Mode == ACTIVITY_MODE_LIMIT && activitytype.info.Refresh == LOGIC_TRUE && actitem.Step < activitytype.info.RefreshMax {
		maxValue := 0
		act_type := GetActivityMgr().GetActivity(actType)
		for num := 0; num < len(act_type.items); num++ {
			if maxValue < act_type.items[num].N[0] {
				maxValue = act_type.items[num].N[0]
			}
		}

		self.DataLocker.RLock()
		for _, v := range self.Sql_Activity.info {
			if v.Id/100 != actType {
				continue
			}
			if v.Done < 2 {
				isNeed = LOGIC_FALSE
				break
			}
		}

		if isNeed == LOGIC_TRUE {
			for _, v := range self.Sql_Activity.info {
				if v.Id/100 != actType {
					continue
				}
				v.Done = 0
				v.Step++
				v.Progress -= maxValue
				if v.Progress <= 0 {
					v.Progress = 0
				}
				for num := 0; num < len(act_type.items); num++ {
					if act_type.items[num].Id != v.Id {
						continue
					}
					if v.Progress >= act_type.items[num].N[0] {
						v.Done = 1
					}
				}

				updateInfo = append(updateInfo, v)
			}
		}
		self.DataLocker.RUnlock()
	}

	if activitytype.info.Mode == ACTIVITY_MODE_LIMIT_4 && activitytype.info.ModeType == ACTIVITY_MODE_TYPE_3 {
		if activityitem.CostNum[0] == 0 {
			actitem.Step++
			actitem.Done = 0
			if actitem.Progress > actitem.Step {
				actitem.Done = 1
			}
		} else if actitem.Step < activityitem.CostNum[0] {
			actitem.Step++
			if actitem.Step < activityitem.CostNum[0] {
				actitem.Done = 0
				if actitem.Progress > actitem.Step {
					actitem.Done = 1
				}
			}
		}
	}

	//日志内容
	if actType == 3001 {
		//英雄传说
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ACTIVITY_3001, actitem.Id, actitem.Step, isNeed, "领取英雄传说奖励", 0, 0, self.player)
	} else if actType == 3002 {
		//阵营传说
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ACTIVITY_3002, actitem.Id, actitem.Step, isNeed, "领取阵营传说奖励", 0, 0, self.player)
	} else if actType == 3008 {
		//狂欢盛典
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ACTIVITY_3008, actitem.Id, self.player.Sql_UserBase.Vip, 0, "购买礼包", 0, 0, self.player)
	} else if actType == 2001 {
		//连续登录
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ACTIVITY_LOGIN_TIMES, actitem.Id, self.player.Sql_UserBase.Vip, 0, "领取连续登陆奖励", 0, 0, self.player)
	}
	var msg S2C_ActivityGet
	msg.Cid = "activityget"
	msg.Id = id
	msg.Item = out
	msg.CostItem = costItem
	msg.UpdateInfo = updateInfo
	msg.GetInfo = actitem
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_ACTIVITY, activitytype.info.Type, activitytype.info.Mode, activitytype.info.ModeType, "活动", 0, 0, self.player)

	//处理月卡
	//self.player.GetModule("activity").(*ModActivity).SendMonthCard()
	//if id == MonthCard1 || id == MonthCard2 {
	//	self.MonthCardCumulativeReward(true)
	//	self.player.GetModule("activity").(*ModActivity).SendMonthCard()
	//}
}

func (self *ModActivity) GetMonthCard(id int) {
	if id != MonthCard1 && id != MonthCard2 {
		return
	}
	self.DataLocker.RLock()
	actitem, ok := self.Sql_Activity.info[id]
	if !ok {
		self.DataLocker.RUnlock()
		return
	}
	self.DataLocker.RUnlock()

	if actitem.Done == 0 {
		return
	}

	if actitem.Done == 2 {
		return
	}

	//看是不是月卡用户
	if !self.IsMonthCard(id - 100800) {
		return
	}

	activitytype := GetActivityMgr().GetActivity(actitem.Id)
	activityitem := GetActivityMgr().GetActivityItem(actitem.Id)
	if activitytype == nil || activityitem == nil || activitytype.status.Status == ACTIVITY_STATUS_CLOSED {
		return
	}

	out := make([]PassItem, 0)
	for j := 0; j < len(activityitem.Item); j++ {
		itemid := activityitem.Item[j]
		if itemid == 0 {
			break
		}
		out = append(out, PassItem{itemid, activityitem.Num[j]})
	}

	actitem.Done = 2

	mailId := 0
	switch id {
	case MonthCard1:
		mailId = 16
	case MonthCard2:
		mailId = 17
	}
	if len(out) > 0 {
		mailConfig, ok := GetCsvMgr().MailConfig[mailId]
		pMail := self.player.GetModule("mail").(*ModMail)
		text := fmt.Sprintf(mailConfig.Mailtxt, self.GetMonthCardDay(id-100800))
		if ok && pMail != nil {
			// 发送邮件
			pMail.AddMailWithItems(MAIL_CAN_ALL_GET, mailConfig.Mailtitle, text, out)
		}
	}
}

func (self *ModActivity) GetGoldWeek() {
	self.DataLocker.RLock()
	actitem, ok := self.Sql_Activity.info[300701]
	if !ok {
		self.DataLocker.RUnlock()
		return
	}
	self.DataLocker.RUnlock()

	if actitem.Done == 0 {
		return
	}

	if actitem.Done == 2 {
		return
	}

	//看是不是月卡用户
	if !self.IsMonthCard(GoldWeek) {
		return
	}

	activitytype := GetActivityMgr().GetActivity(actitem.Id)
	activityitem := GetActivityMgr().GetActivityItem(actitem.Id)
	if activitytype == nil || activityitem == nil || activitytype.status.Status == ACTIVITY_STATUS_CLOSED {
		return
	}

	out := make([]PassItem, 0)
	for j := 0; j < len(activityitem.Item); j++ {
		itemid := activityitem.Item[j]
		if itemid == 0 {
			break
		}
		out = append(out, PassItem{itemid, activityitem.Num[j]})
	}

	actitem.Done = 2

	mailId := 20

	if len(out) > 0 {
		mailConfig, ok := GetCsvMgr().MailConfig[mailId]
		pMail := self.player.GetModule("mail").(*ModMail)
		text := fmt.Sprintf(mailConfig.Mailtxt, self.GetMonthCardDay(GoldWeek))
		if ok && pMail != nil {
			// 发送邮件
			pMail.AddMailWithItems(MAIL_CAN_ALL_GET, mailConfig.Mailtitle, text, out)
		}
	}
}

//! 领取月卡累计领取奖励
func (self *ModActivity) GetMonthCardCumulativeReward(Id int) {

	item := make([]PassItem, 0)

	for i := 0; i < len(self.Sql_Activity.month); i++ {
		if self.Sql_Activity.month[i].Id != 101 && self.Sql_Activity.month[i].Id != 102 {
			return
		}
		//先拿配置
		config := GetCsvMgr().ActivityTotleAwardMap
		if config == nil {
			LogError("Activity_TotleAward表配置没拿到")
			return
		}

		if self.Sql_Activity.month[i].RewardSign == nil {
			self.Sql_Activity.month[i].RewardSign = make(map[int]int)
		}

		_, ok := self.Sql_Activity.month[i].RewardSign[Id]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_ACTIVITIES_DO_NOT_EXIST"))
			return
		}
		if self.Sql_Activity.month[i].RewardSign[Id] == CANTFINISH {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_INADEQUATE_EXCHANGE_CONDITIONS"))
			return
		}
		if self.Sql_Activity.month[i].RewardSign[Id] == TAKEN {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_AWARD_HAS_BEEN_RECEIVED"))
			return
		}
		self.Sql_Activity.month[i].RewardSign[Id] = TAKEN
		//发送奖励
		for i := 0; i < len(config[Id].Item); i++ {
			if config[Id].Item[i] == 0 {
				break
			}
			self.player.AddObject(config[Id].Item[i], config[Id].Num[i], Id, 0, 0, "月卡箱子领取")
			item = append(item, PassItem{config[Id].Item[i], config[Id].Num[i]})
		}

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD_BOX, Id, 0, 0, "月卡箱子领取", 0, 0, self.player)

		var msg S2C_GetCumulativeReward
		msg.Cid = "getcumulativereward"
		msg.Id = Id
		msg.Item = item
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
	}

	self.MonthCardCumulativeReward(false)
	self.SendMonthCard()
}

//! 领取月卡累计领取奖励 新版本
func (self *ModActivity) GetMonthCardCumulativeRewardNew(Id int) {

	for i := 0; i < len(self.Sql_Activity.month); i++ {
		if self.Sql_Activity.month[i].Id != 101 && self.Sql_Activity.month[i].Id != 102 {
			continue
		}
		//先拿配置
		config := GetCsvMgr().MonthCardTotleAwardMap[Id]
		if config == nil {
			continue
		}

		if self.Sql_Activity.month[i].RewardSign == nil {
			self.Sql_Activity.month[i].RewardSign = make(map[int]int)
		}

		if self.Sql_Activity.month[i].RewardSign[Id] == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_AWARD_HAS_BEEN_RECEIVED"))
			return
		}

		if err := self.player.HasObjectOkEasy(config.PointItem, config.PointNum); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}

		self.Sql_Activity.month[i].RewardSign[Id] = LOGIC_TRUE
		//发送奖励
		item := self.player.AddObjectLst(config.Item, config.Num, "月卡箱子领取", Id, 0, 0)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD_BOX, Id, 0, 0, "月卡箱子领取", 0, 0, self.player)

		isRet := true

		for _, v := range GetCsvMgr().MonthCardTotleAwardMap {
			if self.Sql_Activity.month[i].RewardSign[v.Id] == LOGIC_FALSE {
				isRet = false
				break
			}
		}

		temp := new(MonthCardTotleAwardMap)
		costItem := make([]PassItem, 0)
		if isRet {
			for _, v := range GetCsvMgr().MonthCardTotleAwardMap {
				if v.PointNum > temp.PointNum {
					temp = v
				}
			}
			costItem = self.player.RemoveObjectSimple(temp.PointItem, temp.PointNum, "月卡积分刷新", 0, 0, 0)
			self.Sql_Activity.month[i].RewardSign = make(map[int]int)
		}

		var msg S2C_GetCumulativeRewardNew
		msg.Cid = "getcumulativereward"
		msg.Id = Id
		msg.Item = item
		msg.CostItem = costItem
		msg.Month = self.Sql_Activity.month
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

		num := self.player.GetObjectNum(ITEM_MONTH_SCORE)

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MONTH_CARD_SCORE_GET, Id, num, 0, "领取月卡积分奖励", 0, 0, self.player)
		return
	}

	//self.MonthCardCumulativeReward(false)
	//self.SendMonthCard()
}

//! 计算
func (self *ModActivity) MonthCardCumulativeReward(cal bool) {

	for i := 0; i < len(self.Sql_Activity.month); i++ {
		if self.Sql_Activity.month[i].Id != 101 && self.Sql_Activity.month[i].Id != 102 {
			continue
		}
		if i == 0 {
			//先拿配置
			config := GetCsvMgr().ActivityTotleAwardMap
			if config == nil {
				LogError("Activity_TotleAward表配置没拿到")
				return
			}

			//有值的话要检查是不是箱子领完了
			isRet := true
			maxGet := 0
			if self.Sql_Activity.month[i].RewardSign != nil {
				for key, value := range self.Sql_Activity.month[i].RewardSign {
					if value != TAKEN {
						isRet = false
						break
					} else {
						_, ok := config[key]
						if ok && config[key].TotalNum > maxGet {
							maxGet = config[key].TotalNum
						}
					}
				}
			}
			//重置
			if isRet {
				self.Sql_Activity.month[i].Get -= maxGet
				self.Sql_Activity.month[i].RewardSign = make(map[int]int)
			}

			if self.Sql_Activity.month[i].RewardSign == nil {
				self.Sql_Activity.month[i].RewardSign = make(map[int]int)
			}

			for _, value := range config {
				if value.Type == TotleAwardType {
					_, ok := self.Sql_Activity.month[i].RewardSign[value.Id]
					if !ok {
						self.Sql_Activity.month[i].RewardSign[value.Id] = CANTFINISH
					}
				}
			}
			if cal {
				self.Sql_Activity.month[i].Get++
			}
			for key, value := range self.Sql_Activity.month[i].RewardSign {
				_, ok := config[key]
				//check去除掉有可能存在的脏数据
				if !ok {
					delete(self.Sql_Activity.month[i].RewardSign, key)
					continue
				} else if value == CANTFINISH && self.Sql_Activity.month[i].Get >= config[key].TotalNum {
					self.Sql_Activity.month[i].RewardSign[key] = CANTAKE
				}
			}
		} else {
			if self.Sql_Activity.month[i].RewardSign == nil {
				self.Sql_Activity.month[i].RewardSign = make(map[int]int)
			}
			if self.Sql_Activity.month[i-1].RewardSign == nil {
				self.Sql_Activity.month[i-1].RewardSign = make(map[int]int)
			}
			self.Sql_Activity.month[i].Get = self.Sql_Activity.month[i-1].Get
			self.Sql_Activity.month[i].RewardSign = self.Sql_Activity.month[i-1].RewardSign
		}
	}
}

//! 购买基金
func (self *ModActivity) BuyJJ(id int) {
	if self.IsBuyJJ(id) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_CANT_REPEAT_PURCHASE"))
		return
	}

	csv, ok := GetCsvMgr().Data["Activitynew"][id]
	if !ok {
		csv, ok = GetCsvMgr().Data["Activitynew1"][id]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}
		//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		//return
	}

	if self.player.Sql_UserBase.Vip < HF_Atoi(csv["n2"]) {
		//self.player.SendErrInfo("err", "VIP不足,无法购买基金!")
		//return
	}

	if self.player.Sql_UserBase.Gem < HF_Atoi(csv["n3"]) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_DIAMOND_SHORTAGE"))
		return
	}

	self.Sql_Activity.jj = append(self.Sql_Activity.jj, id)

	out := make([]PassItem, 0)
	for i := 0; i < 6; i++ {
		itemid := HF_Atoi(csv[fmt.Sprintf("item%d", i+1)])
		if itemid == 0 {
			break
		}
		out = append(out, PassItem{itemid, HF_Atoi(csv[fmt.Sprintf("num%d", i+1)])})
	}

	for i := 0; i < len(out); i++ {
		out[i].ItemID, out[i].Num = self.player.AddObject(out[i].ItemID, out[i].Num, HF_Atoi(csv["type"]), 0, 0, "活动")
	}

	var msg S2C_ActivityJJ
	msg.Cid = "activityjj"
	msg.JJ = self.Sql_Activity.jj
	msg.Item = out
	self.player.SendMsg("activityjj", HF_JtoB(&msg))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_ACTIVITY, HF_Atoi(csv["type"]), HF_Atoi(csv["mode"]), HF_Atoi(csv["modetype"]), "活动", 0, 0, self.player)
}

//! 刷新活动
func (self *ModActivity) OnRefresh() {
	lst := GetActivityMgr().GetActivityType()
	for i := 0; i < len(lst); i++ {
		act_type := GetActivityMgr().GetActivity(lst[i])
		if act_type != nil && act_type.status.Status > ACTIVITY_STATUS_CLOSED {
			self.DataLocker.RLock()
			for j := 0; j < len(act_type.items); j++ {
				actitem, ok := self.Sql_Activity.info[act_type.items[j].Id]
				if ok {
					// 检查是否刷新
					if act_type.info.Renovate == 1 {
						self.CheckAward(act_type)
						actitem.Done = 0
						actitem.Step = 0
						actitem.Progress = 0
						actitem.Time = act_type.status.EndTime
						actitem.reset(OnRefrsh)
					}
				}
			}
			self.DataLocker.RUnlock()
		}
	}
	//self.GetMonthCard(MonthCard1)
	//self.GetMonthCard(MonthCard2)
	self.GetGoldWeek()

	self.player.SendInfo("updateuserinfo")
}

// 是否是月卡类型
func (self *ModActivity) IsMonthCard(cardtype int) bool {
	for i := 0; i < len(self.Sql_Activity.month); i++ {
		monthcard := self.Sql_Activity.month[i]
		//LogDebug("month card adjust:", monthcard.Id, monthcard.Day, monthcard.StartTime)
		if monthcard.Id == cardtype {
			if monthcard.StartTime+int64(86400*monthcard.Day) > TimeServer().Unix() {
				return true
			}
		}
	}

	return false
}

////! 处理月卡, 作废
//func (self *ModActivity) ProcMonthCard(moneytype int, day int) {
//	for i := 0; i < len(self.Sql_Activity.month); i++ {
//		if self.Sql_Activity.month[i].Id == moneytype {
//			self.Sql_Activity.month[i].Day += day
//
//			return
//		}
//	}
//
//	self.Sql_Activity.month = append(self.Sql_Activity.month, JS_MonthCard{moneytype, TimeServer().Unix(), day, 0, make(map[int]int)})
//
//	upd := self.HandleTask(50, 0, 0, 0)
//	if upd == false {
//		var msg S2C_ActivityUpd
//		msg.Cid = "activityupd"
//		msg.Info = make([]JS_Activity, 0)
//		msg.Month = self.Sql_Activity.month
//		self.player.SendMsg("1", HF_JtoB(&msg))
//	}
//}

func (self *ModActivity) SendMonthCard() {
	//self.GetMonthCard(MonthCard1)
	//self.GetMonthCard(MonthCard2)
	//self.GetGoldWeek()
	var msg S2C_ActivityUpd
	msg.Cid = "activityupd"
	msg.Info = make([]JS_Activity, 0)
	msg.Month = self.Sql_Activity.month
	self.player.SendMsg("1", HF_JtoB(&msg))
}

//! 发送活动数据到客户端
func (self *ModActivity) SendInfo() {
	self.CheckOldData()

	//获得开放活动的代币需求，回收时增加一次豁免
	exchangeItemUse := self.GetExchangeItem()

	var msg S2C_ActivityInfo
	msg.Cid = "activityinfo"
	msg.Info = make([]JS_Activity, 0)
	lst := GetActivityMgr().GetActivityType()
	for i := 0; i < len(lst); i++ {
		activityMask := GetActivityMgr().GetActivity(lst[i])
		if activityMask == nil {
			continue
		}
		if activityMask.status.Status == ACTIVITY_STATUS_CLOSED {
			if activityMask.info.Mode == ACTIVITY_MODE_LIMIT_4 && activityMask.info.ModeType == ACTIVITY_MODE_TYPE_2 {
				self.DataLocker.RLock()
				for j := 0; j < len(activityMask.items); j++ {
					if activityMask.items[j].CostItem[0] == 0 {
						continue
					}

					itemConfig := GetCsvMgr().ItemMap[activityMask.items[j].CostItem[0]]
					if itemConfig == nil || itemConfig.ItemType != ITEM_TYPE_CAN_REMOVE || exchangeItemUse[itemConfig.ItemId] == LOGIC_TRUE {
						continue
					}

					num := self.player.GetObjectNum(itemConfig.ItemId)
					if num > 0 {
						mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_REMOVE_EXCHANEG]
						if !ok {
							continue
						}

						pMail := self.player.GetModule("mail").(*ModMail)
						if pMail == nil {
							continue
						}

						var items []PassItem
						items = append(items, PassItem{ITEM_GOLD, num * itemConfig.Special})
						if len(items) > 0 {
							pMail.AddMail(1, 1, 0, mailConfig.Mailtitle, fmt.Sprintf(mailConfig.Mailtxt, num, itemConfig.ItemName), GetCsvMgr().GetText("STR_SYS"), items, true, 0)
						}
						self.player.RemoveObjectSimple(itemConfig.ItemId, num, "代币回收", activityMask.items[j].Id, 0, 0)
					}
				}
				self.DataLocker.RUnlock()
			}
			continue
		}

		activity := GetActivityMgr().GetActivity(lst[i])
		if activity == nil {
			continue
		}
		startday := HF_Atoi(activity.info.Start)
		if startday < 0 {
			if activityMask.info.Type == ACT_REG_TOTAL_DAY {
				if !self.isActivityRegTotalDayOpen() {
					continue
				}
				//} else if activityMask.info.Type == ACT_GROWTH_GIFT {
				//	if !self.isActivityGrowthGiftOpen() {
				//		continue
				//	}
			} else {
				if !self.player.GetModule("activity").(*ModActivity).isActivityGiftOpen(lst[i]) {
					continue
				}
			}
		}

		if lst[i] >= ACT_OVERFLOW_GIFT_MIN && lst[i] <= ACT_OVERFLOW_GIFT_MAX {
			month := self.player.GetModule("activity").(*ModActivity).Sql_Activity.month
			newMonth := make([]JS_MonthCard, 0)
			for j := 0; j < len(month); j++ {
				if month[j].Id == 502 && (month[j].Stage != lst[i]) {
					continue
				}
				newMonth = append(newMonth, month[j])
			}
			self.player.GetModule("activity").(*ModActivity).Sql_Activity.month = newMonth
		}

		//if activityMask.info.Type == ActSingleRecharge {
		//	if !self.isSingleRegOpen() {
		//		continue
		//	}
		//}
		//
		//if activityMask.info.Type == ActTimeGift {
		//	if !self.isTimeGiftOpen() {
		//		continue
		//	}
		//}
		//if activityMask.info.Type == ActivityOpenFundType {
		//	if !self.isActivityOpenServerFundOpen() {
		//		continue
		//	}
		//}

		start := HF_CalTimeForConfig(activity.info.Start, self.player.Sql_UserBase.Regtime)
		now := TimeServer().Unix()
		self.DataLocker.RLock()
		for j := 0; j < len(activityMask.items); j++ {
			useract, ok := self.Sql_Activity.info[activityMask.items[j].Id]
			if !ok {
				continue
			}

			//fmt.Println("id=", activityMask.items[j].Id, ", ok :", ok)
			if len(activityMask.items[j].N) >= 1 {
				//新增对活动登录的支持
				if activityMask.info.TaskType == TASK_TYPE_ACTIVITY_LOGIN {
					if useract.Ver != activityMask.items[j].N[2] {
						useract.Ver = activityMask.items[j].N[2]
						useract.Done = 0
						useract.Progress = 0
					}
					if useract.Done == 0 {
						useract.Progress = (int(now-start) / DAY_SECS) + 1
						if useract.Progress < 0 {
							useract.Progress = 0
						}
						if useract.Progress >= activityMask.items[j].N[0] {
							useract.Done = 1
						}
					}
				} else {
					self.actCheck(useract, activityMask.info.TaskType, activityMask.items[j].N[0])
				}
			}

			self.makeDiscount(activityMask, useract, &activityMask.items[j])
			if activityMask.info.Renovate == 0 && activityMask.info.CD > 0 &&
				useract.Time > 0 && useract.Time != activityMask.status.EndTime {
				useract.Time = activityMask.status.EndTime
				useract.Progress = 0
				useract.Done = 0
				useract.reset(CdClear)
			}

			msg.Info = append(msg.Info, *useract)
			/*
				if activityMask.info.TaskType != DiscountTask { // 非打折兑换
					if msg.Info[len(msg.Info)-1].Time > 0 {
						msg.Info[len(msg.Info)-1].Time -= activityMask.status.EndTime
					}
				}
			*/
		}
		self.DataLocker.RUnlock()
	}

	msg.JJ = self.Sql_Activity.jj
	msg.Month = self.Sql_Activity.month
	msg.MsgVer = "20180702"
	self.player.SendMsg("activityinfo", HF_JtoB(&msg))

	self.GetActivityMask(0)

	self.SendMonthCard()
}

// 检查是否为月卡用户
func (self *ModActivity) HasMonthCard() bool {
	for i := 0; i < len(self.Sql_Activity.month); i++ {
		monthcard := self.Sql_Activity.month[i]
		if monthcard.StartTime+int64(86400*monthcard.Day) > TimeServer().Unix() {
			return true
		}
	}

	return false
}

// 活动检查
func (self *ModActivity) actCheck(act *JS_Activity, taskType int, n1 int) {
	if taskType != 129 {
		return
	}

	// 检查是否为累计登录类的活动,如果当前活动进度为0,则进度+1, 如果活动进度不是0,则不管
	if act.Progress != 0 {
		return
	}

	if act.Done != 0 {
		return
	}

	act.Progress += 1
	if act.Progress >= n1 {
		act.Done = 1
	}
}

// 初始化当前打折兑换活动
func (self *ModActivity) makeDiscount(pMask *Sql_ActivityMask, userAct *JS_Activity, pItem *JS_ActivityItem) {
	if pMask == nil || userAct == nil || pItem == nil {
		return
	}
	if pMask.info.TaskType != DiscountTask {
		return
	}
	step := pItem.Step
	startTime := pMask.getActTime()
	userAct.Time = startTime + int64(step)*DAY_SECS
	userAct.Step = step
}

func (self *ModActivity) getActCamp(act_type *Sql_ActivityMask, taskType int) int {
	//fmt.Sprintf("info taskType = %d, taskType = %d\n", act_type.info.TaskType, taskType)
	if act_type.info.TaskType != taskType {
		return -1
	}

	items := act_type.items
	if items == nil {
		return -1
	}

	if len(items) <= 0 {
		return -1
	}

	firstItem := items[0]
	if firstItem.N == nil {
		return -1
	}

	if len(firstItem.N) <= 3 {
		return -1
	}

	campType := firstItem.N[2]

	return campType
}

// 更新国家城池榜活动信息
func (self *ModActivity) sendCityAct(item *JS_Activity) {
	var msg S2C_ActivityUpd
	msg.Cid = "activityupd"
	msg.Info = append(msg.Info, *item)
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

// 限时兑换操作
func (self *ModActivity) dailyDiscount(id int, cid string) {
	self.DataLocker.RLock()
	actitem, ok := self.Sql_Activity.info[id]
	if !ok {
		self.DataLocker.RUnlock()
		var msg S2C_ActDailyDiscount
		msg.Cid = cid
		msg.Id = 0
		msg.Item = make([]PassItem, 0)
		self.player.SendMsg(cid, HF_JtoB(&msg))
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKEGG_ACTIVITIES_DO_NOT_EXIST"))
		return
	}
	self.DataLocker.RUnlock()

	if actitem.Done == 2 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_AWARD_HAS_BEEN_RECEIVED"))
		return
	}

	activitytype := GetActivityMgr().GetActivity(actitem.Id)
	if activitytype == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKEGG_ACTIVITIES_DO_NOT_EXIST"))
		return
	}

	activityitem := GetActivityMgr().GetActivityItem(actitem.Id)
	if activityitem == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_SUBACTIVITY_DOES_NOT_EXIST"))
		return
	}

	startTime := activitytype.getActTime()
	realStart := startTime + int64(86400*(activityitem.Step-1))
	now := TimeServer().Unix()
	if now < realStart {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_CAMPAIGN_HASNT_STARTED_YET"))
		return
	}

	// 限时折扣任务类型
	if activitytype.info.TaskType != DiscountTask {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_INCORRECT_TYPES_OF_REWARD_TASKS"))
		return
	}

	param := activityitem.N
	if len(param) <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKEGG_CONFIGURATION_ERROR"))
		return
	}

	needVip := param[0]
	if self.player.Sql_UserBase.Vip < needVip {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TASK_LACK_OF_ARISTOCRATIC_RANK"))
		return
	}

	addItems := make([]PassItem, 0)
	for i := 0; i < 1; i++ {
		itemid := activityitem.Item[i]
		if itemid == 0 {
			continue
		}
		addItems = append(addItems, PassItem{itemid, activityitem.Num[i]})
	}

	cost := make([]PassItem, 0)
	var items []PassItem
	for i := 0; i < 1; i++ {
		itemid := activityitem.CostItem[i]
		if itemid == 0 {
			continue
		}
		// 检查道具是否充足
		if self.player.GetObjectNum(itemid) < activityitem.CostNum[i] {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TIGER_LACK_OF_PROPS"))
			return
		}
		cost = append(cost, PassItem{itemid, activityitem.CostNum[i]})
	}

	// 扣除道具
	for i := 0; i < len(cost); i++ {
		itemId, itemNum := self.player.AddObject(cost[i].ItemID, -cost[i].Num, activitytype.info.Type, 0, 0, "活动")
		items = append(items, NewPassItem(itemId, itemNum))
	}

	// 增加道具
	for i := 0; i < len(addItems); i++ {
		itemId, itemNum := self.player.AddObject(addItems[i].ItemID, addItems[i].Num, activitytype.info.Type, 0, 0, "活动")
		items = append(items, NewPassItem(itemId, itemNum))
	}

	actitem.Progress += 1
	actitem.Done = 2

	var msg S2C_ActDailyDiscount
	msg.Cid = cid
	msg.Id = id
	msg.Progress = actitem.Progress
	msg.Done = actitem.Done
	msg.Item = items
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_ACTIVITY, activitytype.info.Type, activitytype.info.Mode, activitytype.info.ModeType, "活动", 0, 0, self.player)
}

func (self *ModActivity) GetMonth(id int) *JS_MonthCard {
	for _, v := range self.Sql_Activity.month {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

func (self *ModActivity) GetOverflow(id int, cid string) {
	info := self.GetMonth(502)
	if info == nil {
		//没购买
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_ACTIVITY_WAS_NOT_COMPLETED"))
		return
	}

	_, ok := info.RewardSign[id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	if info.RewardSign[id] == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_AWARD_FOR_THE_EVENT_HAS"))
		return
	}

	config := GetActivityMgr().GetActivityItem(id)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	day := int(((TimeServer().Unix() - info.StartTime) / 86400) + 1)
	if day < id%10 {
		self.player.SendErrInfo("err", "STR_MOD_CAMPTASK_FAILURE_TO_MEET_THE_CONDITIONS")
		return
	}

	//发送奖励
	info.RewardSign[id] = LOGIC_TRUE
	getItems := self.player.AddObjectLst(config.Item, config.Num, "超值有礼", id, 0, 0)

	var msg S2C_GetOverflow
	msg.Cid = cid
	msg.Month = info
	msg.GetItems = getItems
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_ACTIVITY, activitytype.info.Type, activitytype.info.Mode, activitytype.info.ModeType, "活动", 0, 0, self.player)
}

// 重置活动
func (self *JS_Activity) reset(reason int) {
	//LogError("done=0, id:", self.Id, ",reason:", reason)
}

//! 独立充值
func (self *ModActivity) isSingleRegOpen() bool {
	rtime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
	activity := GetActivityMgr().GetActivity(ActSingleRecharge)
	if activity == nil {
		LogError("ActSingleRecharge is nil")
		return false
	}
	endTime := rtime.Unix() + int64(activity.info.Continued)
	now := TimeServer().Unix()
	if now >= rtime.Unix() && now <= endTime {
		return true
	}
	return false
}

//! 限时礼包
func (self *ModActivity) isTimeGiftOpen() bool {
	rtime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
	activity := GetActivityMgr().GetActivity(ActTimeGift)
	if activity == nil {
		LogError("ActTimeGift is nil")
		return false
	}
	endTime := rtime.Unix() + int64(activity.info.Continued)
	now := TimeServer().Unix()
	if now >= rtime.Unix() && now <= endTime {
		return true
	}
	return false
}

//!开服基金
func (self *ModActivity) isActivityOpenServerFundOpen() bool {
	rtime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
	activity := GetActivityMgr().GetActivity(ActivityOpenFundType)
	if activity == nil {
		LogError("ActOpenServerFund is nil")
		return false
	}
	endTime := rtime.Unix() + int64(activity.info.Continued) + int64(activity.info.Show)
	now := TimeServer().Unix()
	if now >= rtime.Unix() && now <= endTime {
		return true
	}
	return false
}

// 计算开服基金的状态
func (self *ModActivity) GetOpenFundStatus(act *Sql_ActivityMask) int {
	nStatus := ACTIVITY_STATUS_CLOSED
	if act == nil {
		return nStatus
	}

	regday, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
	rtime := time.Date(regday.Year(), regday.Month(), regday.Day(), 5, 0, 0, 0, regday.Location()).Unix()
	if regday.Hour() < 5 {
		rtime -= DAY_SECS
	}

	endTime := rtime + act.status.EndTime
	startTime := rtime
	timeNow := TimeServer().Unix()

	// 在开启时间之前 或者大于结束时间则为关闭
	if timeNow < startTime || timeNow >= endTime {
		return nStatus
	}

	// 在开始时间后 并在持续时间中则为开启
	if timeNow < startTime+int64(act.info.Continued) && timeNow >= startTime {
		return ACTIVITY_STATUS_OPEN
	}

	// 在持续时间后 并在结束时间中则为开启
	if timeNow < endTime && timeNow >= startTime+int64(act.info.Continued) {
		return ACTIVITY_STATUS_SHOW
	}

	return nStatus
}

func (self *ModActivity) isActivityGiftOpen(nType int) bool {
	regday, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
	rtime := time.Date(regday.Year(), regday.Month(), regday.Day(), 5, 0, 0, 0, regday.Location()).Unix()
	if regday.Hour() < 5 {
		rtime -= DAY_SECS
	}

	activity := GetActivityMgr().GetActivity(nType)
	if activity == nil {
		LogError("ActSingleRecharge is nil")
		return false
	}
	startday := HF_Atoi(activity.info.Start)
	if startday >= 0 {
		LogError("ActSingleRecharge is nil")
		return false
	}

	endTime := rtime + int64(activity.info.Continued) + int64(-(startday + 1)*DAY_SECS)
	now := TimeServer().Unix()
	if activity.info.Continued == 0 && now >= rtime+int64(-(startday + 1)*DAY_SECS) {
		return true
	} else if now >= rtime+int64(-(startday + 1)*DAY_SECS) && now <= endTime {
		return true
	}

	return false
}

func (self *ModActivity) isActivityGrowthGiftOpen() bool {
	regday, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
	rtime := time.Date(regday.Year(), regday.Month(), regday.Day(), 5, 0, 0, 0, regday.Location()).Unix()
	if regday.Hour() < 5 {
		rtime -= DAY_SECS
	}
	activity := GetActivityMgr().GetActivity(ACT_GROWTH_GIFT)
	if activity == nil {
		LogError("ActSingleRecharge is nil")
		return false
	}
	startday := HF_Atoi(activity.info.Start)
	if startday >= 0 {
		LogError("ActSingleRecharge is nil")
		return false
	}

	id := activity.items[0].Id

	actitem, ok := self.Sql_Activity.info[id]
	if ok {
		if actitem.Done == 1 {
			return true
		}
	}

	endTime := rtime + int64(activity.info.Continued) + int64(-(startday + 1)*DAY_SECS)
	now := TimeServer().Unix()
	if now >= rtime+int64(-(startday + 1)*DAY_SECS) && now <= endTime {
		return true
	}

	return false
}

func (self *ModActivity) isActivityRegTotalDayOpen() bool {
	regday, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
	rtime := time.Date(regday.Year(), regday.Month(), regday.Day(), 5, 0, 0, 0, regday.Location()).Unix()
	if regday.Hour() < 5 {
		rtime -= DAY_SECS
	}
	activity := GetActivityMgr().GetActivity(ACT_REG_TOTAL_DAY)
	if activity == nil {
		LogError("isActivityRegTotalDayOpen is nil")
		return false
	}
	startday := HF_Atoi(activity.info.Start)
	if startday >= 0 {
		LogError("isActivityRegTotalDayOpen is nil")
		return false
	}

	//endTime := rtime.Unix() + int64(activity.info.Continued) + int64(-(startday+1)*DAY_SECS)
	now := TimeServer().Unix()
	if now < rtime+int64(-(startday + 1)*DAY_SECS) {
		return false
	}

	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()
	for _, v := range activity.items {
		id := v.Id

		actitem, ok := self.Sql_Activity.info[id]
		if ok {
			if actitem.Done != 2 {
				return true
			}
		} else {
			return true
		}
	}

	return false
}

func (self *ModActivity) sendMail(title, text string, itemlst []PassItem) {
	if len(itemlst) <= 0 {
		return
	}

	pMail := self.player.GetModule("mail").(*ModMail)
	if pMail == nil {
		LogError("checkMail in Daily Activity, pMail == nil!")
		return
	}

	//title := fmt.Sprintf(GetCsvMgr().GetText("STR_DAILY_RECHARGE_AWARD"))
	//text := fmt.Sprintf(GetCsvMgr().GetText("STR_DAILY_RECHARGE_AWARD"))
	pMail.AddMail(1, 1, 0, title, text, GetCsvMgr().GetText("STR_SYS"), itemlst, false, 0)
}

// 发送信息
func (self *ModActivity) SendActivityFundInfo(ver int) bool {
	// 请求基金类型
	nType := ActivityValueFundType
	// 如果为0表示开服基金
	if ver == 0 {
		nType = ActivityOpenFundType
	}

	// 获取活动状态
	actType := GetActivityMgr().GetActivity(nType)

	if actType == nil {
		return false
	}

	// 状态
	nStatus := actType.status.Status
	// 如果是开服基金
	if ver == 0 {
		nStatus = self.GetOpenFundStatus(actType)
	}

	if nStatus == ACTIVITY_STATUS_CLOSED {
		return false
	}

	// 结束时间
	endtime := actType.status.EndTime

	// 如果是开服基金
	if ver == 0 {
		// 获得时间
		t, err := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
		if err != nil {
			return false
		}

		endtime = t.Unix() + actType.status.EndTime
	}

	// 初始化
	if self.Sql_Activity.fund == nil {
		self.Sql_Activity.fund = make(map[int][2]*ActivityFund)
	}

	// 返回配置
	back := [2]*ActivityFund{}
	nVer := actType.getTaskN3()
	nN4 := actType.getTaskN4()

	// 查找有没有相关的配置
	fund, ok := self.Sql_Activity.fund[nVer]

	//! 发消息
	var msg S2C_ActivityFundInfo
	msg.Cid = CS_ACTIVITY_FUND_GET_INFO
	msg.State = nStatus
	msg.EndTime = endtime
	msg.Ver = nVer

	// 当还是活动阶段时 不管有没有都发过去
	if nStatus == ACTIVITY_STATUS_OPEN {
		if ok {
			back = fund
		}
	} else {
		// 当在展示阶段时 如果有就做判断
		if ok {
			timeNow := TimeServer().Unix()

			for _, v := range fund {
				if v == nil {
					continue
				}

				// 查找配置
				groupConfig, ok1 := GetCsvMgr().ActivityFundMap[nN4]
				if !ok1 || groupConfig == nil {
					continue
				}

				config, ok2 := groupConfig.PayConfig[v.Pay]
				if !ok2 || config == nil {
					continue
				}

				//typeConfig, ok2 := GetCsvMgr().ActivityFundTypeMap[v.Pay]
				//if !ok2 || typeConfig == nil {
				//	continue
				//}

				// 查看基金活动自身是否结束
				if v.StartTime > timeNow {
					continue
				}

				// 是否超过自身活动时间
				day := (timeNow-v.StartTime)/DAY_SECS + 1
				if day > (int64)(len(config.DayConfig)) {
					continue
				}

				// 是否有东西没领
				isSend := false
				for _, t := range v.FundGetType {
					if t == 0 {
						isSend = true
						break
					}
				}

				if !isSend {
					continue
				}

				back[config.Type-1] = &ActivityFund{v.Pay, v.Ver, v.ActivityEndTime, v.StartTime, v.FundGetType}
			}
		}
	}

	msg.Fund = back

	backItem := [2][7][]*PassItem{}
	backCost := [2]int{}
	backWroth := [2]int{}

	// 发送配置
	groupConfig, ok := GetCsvMgr().ActivityFundMap[nN4]
	if ok {
		for _, v1 := range groupConfig.PayConfig {
			for _, v := range v1.DayConfig {
				nCount := len(v.Items)
				if nCount == len(v.Nums) {
					for i := 0; i < nCount; i++ {
						backItem[v.Type-1][v.Day-1] = append(backItem[v.Type-1][v.Day-1], &PassItem{v.Items[i], v.Nums[i]})
					}
				}
			}

			backCost[v1.Type-1] = v1.Pay
			backWroth[v1.Type-1] = v1.Worth
		}
	}

	msg.Items = backItem
	msg.Cost = backCost
	msg.Wroth = backWroth
	self.player.SendMsg(CS_ACTIVITY_FUND_GET_INFO, HF_JtoB(&msg))

	return true
}

// 购买活动基金
func (self *ModActivity) BuyActivityFund(pay int) bool {
	// 默认是超值基金 如果是 301 302 特殊处理为开服基金
	nType := 0
	if pay == 301 || pay == 302 {
		nType = ActivityOpenFundType
	} else if pay == 303 || pay == 304 {
		nType = ActivityValueFundType
	} else {
		return false
	}

	// 获取活动状态
	actType := GetActivityMgr().GetActivity(nType)

	if actType == nil {
		return false
	}

	nN3 := actType.getTaskN3()
	nN4 := actType.getTaskN4()

	// 查找配置
	groupConfig, ok1 := GetCsvMgr().ActivityFundMap[nN4]
	if !ok1 || groupConfig == nil {
		return false
	}

	config, ok2 := groupConfig.PayConfig[pay]
	if !ok2 || config == nil {
		return false
	}

	nStatus := actType.status.Status
	if nType == ActivityOpenFundType {
		nStatus = self.GetOpenFundStatus(actType)
	}

	// 活动已结束 则无法参与
	if nStatus == ACTIVITY_STATUS_CLOSED || nStatus == ACTIVITY_STATUS_SHOW {
		return false
	}

	temp := [2]*ActivityFund{}

	// 找到了数据
	fund, ok3 := self.Sql_Activity.fund[nN3]
	if ok3 {
		// 已经买过了
		for i, v := range fund {
			if v == nil {
				continue
			}

			temp[i] = v

			if v.Pay == pay {
				return false
			}
		}
	}

	// 添加数据 表示已经购买了基金
	now := TimeServer()
	timeSet := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local).Unix()
	//if nType == ActivityOpenFundType {
	//	timeSet = now.Unix()
	//}

	// 结束时间
	endtime := actType.status.EndTime

	// 如果是开服基金
	if nType == ActivityOpenFundType {
		// 获得时间
		t, err := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
		if err == nil {
			endtime = t.Unix() + actType.status.EndTime
		}
	}

	temp[config.Type-1] = &ActivityFund{pay, nN3, endtime, timeSet, [7]int{}}
	self.Sql_Activity.fund[nN3] = temp
	//! 发消息
	var msg S2C_ActivityFundActivate
	msg.Cid = CS_ACTIVITY_FUND_ACTIVATE
	msg.Ver = nN3
	msg.Pay = pay

	self.player.SendMsg(CS_ACTIVITY_FUND_ACTIVATE, HF_JtoB(&msg))

	return true
}

// 领取活动基金奖励
func (self *ModActivity) GetActivityFundAward(ver, pay, index int) bool {
	nType := ActivityValueFundType
	if ver == 0 {
		nType = ActivityOpenFundType
	}
	// 获取活动状态
	actType := GetActivityMgr().GetActivity(nType)

	if actType == nil {
		return false
	}

	nN3 := actType.getTaskN3()

	nN4 := actType.getTaskN4()

	// 查找配置
	groupConfig, ok1 := GetCsvMgr().ActivityFundMap[nN4]
	if !ok1 || groupConfig == nil {
		return false
	}

	config, ok2 := groupConfig.PayConfig[pay]
	if !ok2 || config == nil {
		return false
	}

	// 客户端传过来 1-7
	// 领取的天数出错
	if index > len(config.DayConfig) || index <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_FUND_TIME_ERROR"))
		return false
	}

	// 没有该天数的配置
	//dayConfig, ok2 := config.DayConfig[index-1]

	// 没有该天数的配置
	//dayConfig, ok2 := config.DayConfig[index - 1]

	// 没有该天数的配置 储存的数据key为1-7 所以不减一
	dayConfig, ok2 := config.DayConfig[index]

	if !ok2 || dayConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_FUND_TIME_ERROR"))
		return false
	}

	// 是否有资格领取
	fund, ok3 := self.Sql_Activity.fund[nN3]
	if !ok3 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_FUND_DONT_HAVE"))
		return false
	}

	var temp *ActivityFund
	// 是否已经领取
	for _, v := range fund {
		if v == nil {
			continue
		}

		if pay == v.Pay {
			// 数组0-6 所以减一
			if v.FundGetType[index-1] == 1 {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_FUND_IS_GET"))
				return false
			}
			temp = v
			break
		}
	}

	// 查看基金活动自身是否结束
	timeNow := TimeServer().Unix()

	if temp.StartTime > timeNow {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_FUND_TIME_ERROR"))
		return false
	}

	day := (timeNow-temp.StartTime)/DAY_SECS + 1

	if day > (int64)(len(config.DayConfig)) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_FUND_RUN_OUT_TIME"))
		return false
	}

	// 领取的天数未到 直接计算
	if day < (int64)(index) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_FUND_TIME_ERROR"))
		return false
	}

	// 添加物品
	outItems := self.player.AddObjectLst(dayConfig.Items, dayConfig.Nums, "周基金奖励领取", dayConfig.ID, 0, 0)

	// 设置领取状态  数组0-6 所以减一
	temp.FundGetType[index-1] = 1

	//! 发消息
	var msg S2C_ActivityFundAward
	msg.Cid = CS_ACTIVITY_FUND_GET_AWARD
	msg.Pay = pay
	msg.Ver = nN3
	msg.Day = index
	msg.Item = outItems

	self.player.SendMsg(CS_ACTIVITY_FUND_GET_AWARD, HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ACTIVITY_FUND_AWARD, dayConfig.ID, 0, 0, "周基金奖励领取", 0, 0, self.player)

	return true
}
func (self *ModActivity) CheckTaskDone() {

	lst := GetActivityMgr().GetActivityType()
	for i := 0; i < len(lst); i++ {
		// 获取活动状态
		act_type := GetActivityMgr().GetActivity(lst[i])
		// 活动如果开启
		if act_type == nil || act_type.status.Status != ACTIVITY_STATUS_OPEN {
			continue
		}

		// 过滤任务类型
		if act_type.info.TaskType != CommonLevelTask && act_type.info.TaskType != EliteLevelTask {
			continue
		}

		// 道具信息遍历
		for j := 0; j < len(act_type.items); j++ {
			// 读取并插入数据
			actitem, ok := self.Sql_Activity.info[act_type.items[j].Id]
			if actitem == nil || !ok {
				continue
			}

			// 过滤已经完成和领取奖励的任务
			if actitem.Done != 0 {
				continue
			}

			// 创建任务节点
			var tasknode TaskNode
			tasknode.Tasktypes = act_type.info.TaskType
			tasknode.N1 = act_type.items[j].N[0]
			tasknode.N2 = act_type.items[j].N[1]
			tasknode.N3 = act_type.items[j].N[2]
			tasknode.N4 = act_type.items[j].N[3]

			passId := tasknode.N2

			if act_type.info.TaskType == CommonLevelTask {
				pass := self.player.GetModule("pass").(*ModPass).GetPass(passId)
				if pass != nil {
					actitem.Progress = 1
					actitem.Done = 1
				}
			}
		}
	}
}

// 没考虑指定日期
func (self *ModActivity) IsActivityOpen(nType int) bool {

	activity := GetActivityMgr().GetActivity(nType)
	if activity == nil {
		return false
	}

	startday := HF_Atoi(activity.info.Start)
	if startday >= 0 {
		if activity.status.Status > ACTIVITY_STATUS_CLOSED {
			return true
		}
	} else {
		rtime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
		endTime := rtime.Unix() + int64(activity.info.Continued) + int64(-(startday + 1)*DAY_SECS)
		now := TimeServer().Unix()
		if activity.info.Continued == 0 && now >= rtime.Unix()+int64(-(startday + 1)*DAY_SECS) {
			return true
		} else if now >= rtime.Unix()+int64(-(startday + 1)*DAY_SECS) && now <= endTime {
			return true
		}
	}

	return false
}

func (self *ModActivity) GetActivityStart(nType int) int64 {

	activity := GetActivityMgr().GetActivity(nType)
	if activity == nil {
		return 0
	}

	startday := HF_Atoi(activity.info.Start)
	if startday >= 0 {
		return GetServer().GetOpenServer() + int64((startday-1)*DAY_SECS)
	} else {
		rtime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
		correctTime := HF_CalPlayerCreateTime(rtime.Unix(), 0)
		return correctTime
	}

	return 0
}

func (self *ModActivity) GetDailyItem() []PassItem {
	items := make([]PassItem, 0)

	passId := self.player.GetModule("onhook").(*ModOnHook).GetDailyPassId()
	group := 1 //默认值

	config, ok := GetCsvMgr().LevelConfigMap[passId]
	if ok && config.DailyRechargeGroup > 0 {
		group = config.DailyRechargeGroup
	}

	configItem := GetCsvMgr().ActivityDailyRecharge[group]
	if configItem == nil {
		configItem = GetCsvMgr().ActivityDailyRecharge[1]
	}

	for i := 0; i < len(configItem.Item); i++ {
		if configItem.Item[i] > 0 {
			items = append(items, PassItem{
				ItemID: configItem.Item[i],
				Num:    configItem.Num[i],
			})
		}
	}
	return items
}

func (self *ModActivity) CalActivityStartTime(start int) int64 {
	if start > 0 {
		return GetServer().GetOpenServer() + int64((start-1)*DAY_SECS)
	} else if start < 0 {
		rtime, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.Regtime, time.Local)
		firstTime := HF_CalPlayerCreateTime(rtime.Unix(), 0)
		return firstTime
	}
	return 0
}

// 只处理限时礼包 因为需要完成就领取
func (self *ModActivity) HandleRecharge(grade int) {
	lst := GetActivityMgr().GetActivityType()
	for i := 0; i < len(lst); i++ {
		if self.IsLimitGiftType(lst[i]) {
			// 获取活动状态
			act_type := GetActivityMgr().GetActivity(lst[i])
			// 活动如果开启
			if act_type != nil && act_type.status.Status > ACTIVITY_STATUS_CLOSED && act_type.getTaskType() == TASK_TYPE_RECHARGE_ONCE {
				// 道具信息遍历
				for j := 0; j < len(act_type.items); j++ {
					actitem, ok := self.Sql_Activity.info[act_type.items[j].Id]
					if ok && actitem != nil && actitem.Done == 1 && act_type.items[j].N[1] == grade {
						self.FinishActivity(act_type.items[j].Id)
						return
					}
				}
			}
		}
	}
}

//获得月卡特权的免费快速挂机次数
func (self *ModActivity) GetMonthFreeFast() int {
	value := 0
	for i := 0; i < len(self.Sql_Activity.month); i++ {
		if self.Sql_Activity.month[i].Id != 101 && self.Sql_Activity.month[i].Id != 102 {
			continue
		}

		if self.Sql_Activity.month[i].StartTime+int64(86400*self.Sql_Activity.month[i].Day) < TimeServer().Unix() {
			continue
		}

		//先拿配置
		config := GetCsvMgr().MonthCard[self.Sql_Activity.month[i].Id-100]
		if config == nil {
			continue
		}

		value += config.HangUpFast
	}
	return value
}

//获得月卡特权的英雄经验加成
func (self *ModActivity) GetMonthHeroExp() int {
	value := 0
	for i := 0; i < len(self.Sql_Activity.month); i++ {
		if self.Sql_Activity.month[i].Id != 101 && self.Sql_Activity.month[i].Id != 102 {
			continue
		}

		if self.Sql_Activity.month[i].StartTime+int64(86400*self.Sql_Activity.month[i].Day) < TimeServer().Unix() {
			continue
		}

		//先拿配置
		config := GetCsvMgr().MonthCard[self.Sql_Activity.month[i].Id-100]
		if config == nil {
			continue
		}

		value += config.HangUpHeroExp
	}
	return value
}

//获得月卡特权的金币加成
func (self *ModActivity) GetMonthGold() int {
	value := 0
	for i := 0; i < len(self.Sql_Activity.month); i++ {
		if self.Sql_Activity.month[i].Id != 101 && self.Sql_Activity.month[i].Id != 102 {
			continue
		}

		if self.Sql_Activity.month[i].StartTime+int64(86400*self.Sql_Activity.month[i].Day) < TimeServer().Unix() {
			continue
		}

		//先拿配置
		config := GetCsvMgr().MonthCard[self.Sql_Activity.month[i].Id-100]
		if config == nil {
			continue
		}

		value += config.HangUpGold
	}
	return value
}

func (self *ModActivity) IsLimitGiftType(id int) bool {
	if ACT_TIME_LIMIT_GIFT_START <= id && id <= ACT_TIME_LIMIT_GIFT_END {
		return true
	}
	if ACT_TIME_LIMIT_GIFT_EX_START <= id && id <= ACT_TIME_LIMIT_GIFT_EX_END {
		return true
	}
	return false
}

func (self *ModActivity) IsNeedSavePlan(act_type *Sql_ActivityMask) bool {
	if act_type.info.Mode == ACTIVITY_MODE_LIMIT && act_type.info.Refresh == LOGIC_TRUE {
		return true
	} else if act_type.info.Mode == ACTIVITY_MODE_LIMIT_4 && act_type.info.ModeType == ACTIVITY_MODE_TYPE_3 {
		return true
	}
	return false
}

func (self *ModActivity) CheckOldData() {
	self.CheckResetSign()

	now := TimeServer().Unix()
	for k, v := range self.Sql_Activity.activityResetSign {
		//先看对应活动是否开放
		activityMask := GetActivityMgr().GetActivity(k)
		if activityMask == nil {
			continue
		}
		if activityMask.status.Status == ACTIVITY_STATUS_CLOSED {
			continue
		}

		if v < now {
			startTime := activityMask.getActTime()
			endTime := startTime + int64(activityMask.info.Continued) + int64(activityMask.info.CD)
			self.Sql_Activity.activityResetSign[k] = endTime

			self.DataLocker.RLock()
			for j := 0; j < len(activityMask.items); j++ {
				useract, ok := self.Sql_Activity.info[activityMask.items[j].Id]
				if !ok {
					continue
				}
				useract.Time = activityMask.status.EndTime
				useract.Progress = 0
				useract.Done = 0
				useract.Step = 0
			}
			self.DataLocker.RUnlock()
		}
	}
}

func (self *ModActivity) GetExchangeItem() map[int]int {
	data := make(map[int]int)
	lst := GetActivityMgr().GetActivityType()
	for i := 0; i < len(lst); i++ {
		activityMask := GetActivityMgr().GetActivity(lst[i])
		if activityMask == nil {
			continue
		}
		if activityMask.status.Status == ACTIVITY_STATUS_OPEN {
			if activityMask.info.Mode == ACTIVITY_MODE_LIMIT_4 && activityMask.info.ModeType == ACTIVITY_MODE_TYPE_2 {
				self.DataLocker.RLock()
				for j := 0; j < len(activityMask.items); j++ {
					if activityMask.items[j].CostItem[0] == 0 {
						continue
					}

					itemConfig := GetCsvMgr().ItemMap[activityMask.items[j].CostItem[0]]
					if itemConfig == nil || itemConfig.ItemType != ITEM_TYPE_CAN_REMOVE {
						continue
					}
					data[itemConfig.ItemId] = LOGIC_TRUE
				}
				self.DataLocker.RUnlock()
			}
			continue
		}
	}
	return data
}

func (self *ModActivity) CheckResetSign() {
	if self.Sql_Activity.activityResetSign == nil {
		self.Sql_Activity.activityResetSign = make(map[int]int64)
	}

	_, ok := self.Sql_Activity.activityResetSign[ACT_SEASON_ACTIVITY]
	if !ok {
		self.Sql_Activity.activityResetSign[ACT_SEASON_ACTIVITY] = 0
	}

	_, ok = self.Sql_Activity.activityResetSign[ACT_FESTIVAL_EXCHANGE]
	if !ok {
		self.Sql_Activity.activityResetSign[ACT_FESTIVAL_EXCHANGE] = 0
	}
}
