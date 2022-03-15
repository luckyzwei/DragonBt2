package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

const (
	TEAMPOS_DIS = 1682
)

type JS_ExChangeInfo struct {
	Id    int `json:"taskid"`    // id
	Times int `json:"tasktypes"` // 当前兑换次数
}

// 收藏家积分奖励
type ActivityBossAward struct {
	AwardId    int   `json:"awardid"`    // 奖励ID
	Group      int   `json:"group"`      // 进度
	NeedPoint  int   `json:"needpoint"`  // 任务类型
	Pickup     int   `json:"pickup"`     // 是否领取奖励
	Items      []int `json:"items"`      // 奖励ID
	Nums       []int `json:"nums"`       // 奖励数量
	ChangeTime int64 `json:"changetime"` // 奖励更换时间， 0表示不更换
	Notice     int   `json:"notice"`     // 是否通知
}

type JS_ActivityBossInfo struct {
	Id           int                `json:"id"`           // 奖励ID
	Period       int                `json:"period"`       // 当前期数
	Times        int                `json:"times"`        // 当前攻打次数
	TaskInfo     []*JS_TaskInfo     `json:"taskinfo"`     // 每日任务进度
	StartTime    int64              `json:"starttime"`    // 开始时间
	EndTime      int64              `json:"endtime"`      // 结束时间
	RewardTime   int64              `json:"rewardtime"`   // 发奖时间
	RefreshTimes int                `json:"refreshtimes"` // 刷新次数
	ExChange     []*JS_ExChangeInfo `json:"exchange"`     // 兑换信息
}

//! 任务数据库
type San_ActivityBoss struct {
	Uid              int64
	ActivityBossInfo string

	activityBossInfo map[int]*JS_ActivityBossInfo //! 任务信息
	DataUpdate
}

//! 任务
type ModActivityBoss struct {
	player           *Player
	Sql_ActivityBoss San_ActivityBoss

	chg map[int]*JS_ActivityBossInfo
}

func (self *ModActivityBoss) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_useractivityboss` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_ActivityBoss, "san_useractivityboss", self.player.ID)

	if self.Sql_ActivityBoss.Uid <= 0 {
		self.Sql_ActivityBoss.Uid = self.player.ID
		self.Sql_ActivityBoss.activityBossInfo = make(map[int]*JS_ActivityBossInfo, 0)
		self.Encode()
		InsertTable("san_useractivityboss", &self.Sql_ActivityBoss, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_ActivityBoss.Init("san_useractivityboss", &self.Sql_ActivityBoss, true)
}

//! 将数据库数据写入data
func (self *ModActivityBoss) Decode() {
	json.Unmarshal([]byte(self.Sql_ActivityBoss.ActivityBossInfo), &self.Sql_ActivityBoss.activityBossInfo)
}

//! 将data数据写入数据库
func (self *ModActivityBoss) Encode() {
	self.Sql_ActivityBoss.ActivityBossInfo = HF_JtoA(&self.Sql_ActivityBoss.activityBossInfo)
}

func (self *ModActivityBoss) OnGetOtherData() {

}

// 注册消息
func (self *ModActivityBoss) onReg(handlers map[string]func(body []byte)) {
	handlers["activitybosstask"] = self.ActivityBossTask
	handlers["activitybossstart"] = self.Fight
	handlers["activitybossstartex"] = self.FightEx
	handlers["activitybossresultex"] = self.ActivityBossResultEx
	handlers["activitybossresettimes"] = self.ActivityBossResetTimes
	handlers["activitybossgetrank"] = self.ActivityBossGetRank
	handlers["activitybossgetrecord"] = self.ActivityBossGetRecord
	handlers["activitybossexchange"] = self.ActivityBossExchange
}

func (self *ModActivityBoss) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModActivityBoss) OnSave(sql bool) {
	self.Encode()
	self.Sql_ActivityBoss.Update(sql)
}

//每日任务刷新
func (self *ModActivityBoss) OnRefresh() {

	itemMap := make(map[int]*Item)
	for _, v := range self.Sql_ActivityBoss.activityBossInfo {
		for _, task := range v.TaskInfo {
			if task.Finish != CANTAKE {
				continue
			}
			if task.Pickup == LOGIC_TRUE {
				continue
			}
			config := GetCsvMgr().GetActivityBossTargetConfig(v.Id, v.Period, task.Taskid)
			if config == nil {
				continue
			}
			AddItemMapHelper3(itemMap, config.Item, config.Num)
		}
	}
	if len(itemMap) > 0 {
		mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_ACTIVITYBOSS_TASK]
		if ok {
			itemLst := make([]PassItem, 0)
			for _, v := range itemMap {
				itemLst = append(itemLst, PassItem{ItemID: v.ItemId, Num: v.ItemNum})
			}
			self.player.GetModule("mail").(*ModMail).AddMail(1,
				1, 0, mailConfig.Mailtitle, mailConfig.Mailtxt, GetCsvMgr().GetText("STR_SYS"), itemLst, false, 0)
		}
	}

	for _, v := range self.Sql_ActivityBoss.activityBossInfo {
		v.Times = 0
		v.RefreshTimes = 0
		for _, task := range v.TaskInfo {
			task.Plan = 0
			task.Finish = CANTFINISH
			task.Pickup = LOGIC_FALSE
		}
	}
}

func (self *ModActivityBoss) Check() {
	exchangeItemUse := self.GetExchangeItem()
	for i := ACT_BOSS_MIN; i < ACT_BOSS_MAX; i++ {
		_, ok := self.Sql_ActivityBoss.activityBossInfo[i]
		if !ok {
			self.Sql_ActivityBoss.activityBossInfo[i] = new(JS_ActivityBossInfo)
			self.Sql_ActivityBoss.activityBossInfo[i].Id = i
		}

		isOpen, _ := GetActivityMgr().JudgeOpenAllIndex(self.Sql_ActivityBoss.activityBossInfo[i].Id, self.Sql_ActivityBoss.activityBossInfo[i].Id)
		if !isOpen {
			configBoss := GetCsvMgr().GetActivityBossConfig(self.Sql_ActivityBoss.activityBossInfo[i].Id, self.Sql_ActivityBoss.activityBossInfo[i].Period)
			if configBoss == nil {
				continue
			}
			mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_ACTIVITYBOSS_GOLD]
			if !ok {
				continue
			}

			pMail := self.player.GetModule("mail").(*ModMail)
			if pMail == nil {
				continue
			}

			itemMap := make(map[int]*Item, 0)
			for j := 0; j < len(configBoss.Items); j++ {
				if configBoss.Items[j] == 0 || exchangeItemUse[configBoss.Items[j]] == LOGIC_TRUE {
					continue
				}
				num := self.player.GetObjectNum(configBoss.Items[j])
				if num > 0 {
					configItem := GetCsvMgr().ItemMap[configBoss.Items[j]]
					if configItem != nil && configItem.Special > 0 {
						self.player.RemoveObjectSimple(configBoss.Items[j], num, "世界BOSS刷新", self.Sql_ActivityBoss.activityBossInfo[i].Id, self.Sql_ActivityBoss.activityBossInfo[i].Period, 0)
						AddItemMapHelper3(itemMap, ITEM_GOLD, num*configItem.Special)
					}
				}
			}

			var items []PassItem
			for _, v := range itemMap {
				if v.ItemId == 0 {
					continue
				}

				if v.ItemNum == 0 {
					continue
				}
				items = append(items, PassItem{v.ItemId, v.ItemNum})
			}
			if len(itemMap) > 0 {
				pMail.AddMail(1, 1, 0, mailConfig.Mailtitle, fmt.Sprintf(mailConfig.Mailtxt, configBoss.Name), GetCsvMgr().GetText("STR_SYS"), items, true, 0)
			}
			continue
		}

		activity := GetActivityMgr().GetActivity(self.Sql_ActivityBoss.activityBossInfo[i].Id)
		if activity == nil {
			continue
		}

		period := GetActivityMgr().getActN3(self.Sql_ActivityBoss.activityBossInfo[i].Id)

		if self.Sql_ActivityBoss.activityBossInfo[i].Period != period {
			//刷新
			self.Sql_ActivityBoss.activityBossInfo[i].Period = period
			self.Sql_ActivityBoss.activityBossInfo[i].Times = 0
			self.Sql_ActivityBoss.activityBossInfo[i].RefreshTimes = 0
			self.Sql_ActivityBoss.activityBossInfo[i].TaskInfo = make([]*JS_TaskInfo, 0)
			self.Sql_ActivityBoss.activityBossInfo[i].ExChange = make([]*JS_ExChangeInfo, 0)

			_, ok := GetCsvMgr().ActivityBossTargetConfig[self.Sql_ActivityBoss.activityBossInfo[i].Id]
			if !ok {
				continue
			}

			config := GetCsvMgr().ActivityBossTargetConfig[self.Sql_ActivityBoss.activityBossInfo[i].Id][self.Sql_ActivityBoss.activityBossInfo[i].Period]
			if config == nil {
				continue
			}

			for _, v := range config {
				task := new(JS_TaskInfo)
				task.Taskid = v.TaskId
				task.Tasktypes = v.TaskTypes
				self.Sql_ActivityBoss.activityBossInfo[i].TaskInfo = append(self.Sql_ActivityBoss.activityBossInfo[i].TaskInfo, task)
			}
		}
		self.Sql_ActivityBoss.activityBossInfo[i].StartTime = HF_CalTimeForConfig(activity.info.Start, self.player.Sql_UserBase.Regtime)
		self.Sql_ActivityBoss.activityBossInfo[i].EndTime = self.Sql_ActivityBoss.activityBossInfo[i].StartTime + int64(activity.info.Continued) + int64(activity.info.Show)
		self.Sql_ActivityBoss.activityBossInfo[i].RewardTime = self.Sql_ActivityBoss.activityBossInfo[i].StartTime + int64(activity.info.Continued)

		_, okExchange := GetCsvMgr().ActivityBossExchangeConfig[self.Sql_ActivityBoss.activityBossInfo[i].Id]
		if !okExchange {
			continue
		}

		configExchange := GetCsvMgr().ActivityBossExchangeConfig[self.Sql_ActivityBoss.activityBossInfo[i].Id][self.Sql_ActivityBoss.activityBossInfo[i].Period]
		if configExchange == nil {
			continue
		}

		for _, v := range configExchange {
			isFind := false
			for j := 0; j < len(self.Sql_ActivityBoss.activityBossInfo[i].ExChange); j++ {
				if self.Sql_ActivityBoss.activityBossInfo[i].ExChange[j].Id == v.Id {
					isFind = true
					break
				}
			}
			if !isFind {
				exchange := new(JS_ExChangeInfo)
				exchange.Id = v.Id
				self.Sql_ActivityBoss.activityBossInfo[i].ExChange = append(self.Sql_ActivityBoss.activityBossInfo[i].ExChange, exchange)
			}
		}
	}
}

//涉及东西扣除，调用需要在mode_bag之前
func (self *ModActivityBoss) SendInfo() {
	self.Check()

	var msg S2C_ActivityBossInfo
	msg.Cid = "activitybossinfo"
	msg.ActivityBossInfo = self.Sql_ActivityBoss.activityBossInfo
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModActivityBoss) ActivityBossGetRank(body []byte) {
	var msg C2S_ActivityBossTask
	json.Unmarshal(body, &msg)

	GetActivityBossMgr().GetRank(self.player, msg.Id)
}

func (self *ModActivityBoss) ActivityBossGetRecord(body []byte) {
	var msg C2S_ActivityBossGetRecord
	json.Unmarshal(body, &msg)

	GetActivityBossMgr().GetRecord(self.player, msg.Id, msg.TargetUid)
}

func (self *ModActivityBoss) ActivityBossExchange(body []byte) {
	var msg C2S_ActivityBossExchange
	json.Unmarshal(body, &msg)

	node := self.GetExchange(msg.ActivityId, msg.Id)
	if node == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	_, ok := self.Sql_ActivityBoss.activityBossInfo[msg.ActivityId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	config := GetCsvMgr().GetActivityBossExchangeConfig(msg.ActivityId, self.Sql_ActivityBoss.activityBossInfo[msg.ActivityId].Period, msg.Id)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	if node.Times >= config.Frequency {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_REACH_THE_UPPER_LIMIT"))
		return
	}

	//判断消耗
	if err := self.player.HasObjectOk(config.NeedItem, config.NeedNum); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	CostItems := self.player.RemoveObjectLst(config.NeedItem, config.NeedNum, "兑换活动", 0, 0, 0)
	//发送任务奖励
	GetItems := self.player.AddObjectSimple(config.Item, config.Num, "兑换活动", 0, 0, 0)
	node.Times++

	var msgRel S2C_ActivityBossExchange
	msgRel.Cid = "activitybossexchange"
	msgRel.GetItems = GetItems
	msgRel.CostItems = CostItems
	msgRel.Exchange = node
	msgRel.ActivityId = msg.ActivityId
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

}

func (self *ModActivityBoss) FightEndFail() {
	var msg S2C_ActivityBossFightFail
	msg.Cid = "activitybossfightfail"
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModActivityBoss) FightEndOK(player *Player, battleInfo BattleInfo, attack *JS_FightInfo, defence *JS_FightInfo, bossId int) {

	//获得伤害
	score := int64(0)
	if battleInfo.UserInfo[1] != nil {
		for _, v := range battleInfo.UserInfo[1].HeroInfo {
			score += v.TakeDamage
		}
	}

	//任务系统
	player.HandleTask(TASK_TYPE_ACTIVITY_BOSS_HURT_SINGLE, int(score), bossId, 0)
	player.HandleTask(TASK_TYPE_ACTIVITY_BOSS_HURT_ALL, int(score), bossId, 0)
	player.HandleTask(TASK_TYPE_ACTIVITY_BOSS_COUNT, 1, bossId, 0)

	GetActivityBossMgr().UpdatePoint(self.player, score, bossId, attack, defence, battleInfo)

	self.Sql_ActivityBoss.activityBossInfo[bossId].Times++

	var msg S2C_ActivityBossFightOK
	msg.Cid = "activitybossfightok"
	msg.RandNum = battleInfo.Random
	msg.FightId = battleInfo.Id
	msg.Score = score
	msg.FightInfo[0] = attack
	msg.FightInfo[1] = defence
	msg.SelfInfo = GetActivityBossMgr().GetPlayerInfo(self.player, bossId)
	msg.ActivityBossInfo = self.Sql_ActivityBoss.activityBossInfo[bossId]
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModActivityBoss) HandleTask(tasktype, n2, n3, n4 int) {

	now := TimeServer().Unix()
	for _, v := range self.Sql_ActivityBoss.activityBossInfo {
		if now > v.RewardTime {
			continue
		}
		for _, node := range v.TaskInfo {
			if node.Tasktypes != tasktype {
				continue
			}
			if node.Finish > CANTFINISH {
				continue
			}
			config := GetCsvMgr().GetActivityBossTargetConfig(v.Id, v.Period, node.Taskid)
			if config == nil {
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
				node.Plan += plan
				chg = true
			} else {
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

			if node.Plan >= config.Ns[0] {
				node.Finish = LOGIC_TRUE
			}

			if chg {
				if self.chg == nil {
					self.chg = make(map[int]*JS_ActivityBossInfo, 0)
				}
				_, ok := self.chg[v.Id]
				if !ok {
					self.chg[v.Id] = v
				}
			}
		}
	}
}

// 发送礼包更新信息
func (self *ModActivityBoss) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_TaskUpdateBoss
	msg.Cid = "activitybossupdate"
	for _, v := range self.chg {
		msg.Info = append(msg.Info, v)
	}
	self.chg = make(map[int]*JS_ActivityBossInfo, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}

func (self *ModActivityBoss) GetTask(id int, taskId int) *JS_TaskInfo {
	_, ok := self.Sql_ActivityBoss.activityBossInfo[id]
	if !ok {
		return nil
	}

	for _, v := range self.Sql_ActivityBoss.activityBossInfo[id].TaskInfo {
		if v.Taskid == taskId {
			return v
		}
	}
	return nil
}

func (self *ModActivityBoss) GetExchange(activityId int, id int) *JS_ExChangeInfo {
	_, ok := self.Sql_ActivityBoss.activityBossInfo[activityId]
	if !ok {
		return nil
	}

	for _, v := range self.Sql_ActivityBoss.activityBossInfo[activityId].ExChange {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (self *ModActivityBoss) ActivityBossTask(body []byte) {
	var msg C2S_ActivityBossTask
	json.Unmarshal(body, &msg)

	node := self.GetTask(msg.Id, msg.TaskId)
	if node == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	if node.Finish == CANTFINISH {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DIPLOMACY_THE_TASK_IS_NOT_COMPLETED"))
		return
	}

	if node.Pickup == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_AWARD_HAS_BEEN_RECEIVED"))
		return
	}

	_, ok := self.Sql_ActivityBoss.activityBossInfo[msg.Id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	config := GetCsvMgr().GetActivityBossTargetConfig(msg.Id, self.Sql_ActivityBoss.activityBossInfo[msg.Id].Period, msg.TaskId)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	//发送任务奖励
	items := self.player.AddObjectSimple(config.Item, config.Num, "世界BOSS", node.Taskid, 0, 0)
	node.Pickup = LOGIC_TRUE

	var msgRel S2C_ActivityBossTask
	msgRel.Cid = "activitybosstask"
	msgRel.GetItems = items
	msgRel.ActivityId = msg.Id
	msgRel.TaskInfo = node
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ActivityBoss_GET, node.Taskid, num, oldNum, LOG_LANGUAGE_ActivityBoss_GET, 0, 0, self.player)
}

func (self *ModActivityBoss) Fight(body []byte) {

	var msg C2S_ActivityBossFight
	json.Unmarshal(body, &msg)

	info, ok := self.Sql_ActivityBoss.activityBossInfo[msg.Id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	//次数验证
	config := GetCsvMgr().GetActivityBossConfig(info.Id, info.Period)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	if info.Times >= config.Times {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_TODAYS_FREE_CHALLENGES_ARE_EXHAUSTED"))
		return
	}

	attack := GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, msg.Id-TEAMPOS_DIS)
	defence := GetActivityBossMgr().GetBossFightInfo(msg.Id)
	randNum := TimeServer().Unix()
	//战报加入战斗服务器等待队列
	fightId := GetArenaMgr().AddFightListForBoss(0,
		attack,
		defence,
		randNum,
		randNum,
		msg.Id)

	if fightId > 0 {
		msg := &S2C_ActivityBossStart{}
		msg.Cid = "activitybossstart"
		msg.FightId = fightId
		self.player.Send(msg.Cid, msg)
	}
}

func (self *ModActivityBoss) FightEx(body []byte) {

	var msg C2S_ActivityBossFight
	json.Unmarshal(body, &msg)

	info, ok := self.Sql_ActivityBoss.activityBossInfo[msg.Id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	//次数验证
	config := GetCsvMgr().GetActivityBossConfig(info.Id, info.Period)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	if info.Times >= config.Times {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_TODAYS_FREE_CHALLENGES_ARE_EXHAUSTED"))
		return
	}

	msgRel := &S2C_ActivityBossStartEx{}
	msgRel.Cid = "activitybossstartex"
	msgRel.FightInfo = GetActivityBossMgr().GetBossFightInfo(msg.Id)
	msgRel.PlayerFightInfo = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, msg.Id-TEAMPOS_DIS)
	self.player.Send(msgRel.Cid, msgRel)
}

func (self *ModActivityBoss) ActivityBossResultEx(body []byte) {

	var msg C2S_ActivityBossResultEx
	json.Unmarshal(body, &msg)

	info, ok := self.Sql_ActivityBoss.activityBossInfo[msg.Id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	//次数验证
	config := GetCsvMgr().GetActivityBossConfig(info.Id, info.Period)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	if info.Times >= config.Times {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_TODAYS_FREE_CHALLENGES_ARE_EXHAUSTED"))
		return
	}

	msg.BattleInfo.Id = GetFightMgr().GetFightInfoID()
	attackFightInfo := GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, msg.Id-TEAMPOS_DIS)
	defendFightInfo := GetActivityBossMgr().GetBossFightInfo(msg.Id)

	msg.BattleInfo.UserInfo[0].Name = attackFightInfo.Uname
	msg.BattleInfo.UserInfo[0].Icon = attackFightInfo.Iconid
	msg.BattleInfo.UserInfo[0].UnionName = attackFightInfo.UnionName
	msg.BattleInfo.UserInfo[1].Name = defendFightInfo.Uname

	self.FightEndOK(self.player, *msg.BattleInfo, attackFightInfo, defendFightInfo, msg.Id)
}

func (self *ModActivityBoss) ActivityBossResetTimes(body []byte) {

	var msg C2S_ActivityBossResetTimes
	json.Unmarshal(body, &msg)

	info, ok := self.Sql_ActivityBoss.activityBossInfo[msg.Id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	//次数验证
	config := GetCsvMgr().GetActivityBossConfig(info.Id, info.Period)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	if info.Times < config.Times {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_TOWER_RESET_CURLEVEL"))
		return
	}

	configVip := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if configVip == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_DOES_NOT"))
		return
	}

	if info.RefreshTimes >= configVip.BossReset {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_DRAW_THE_NUMBER_OF_RESETS_EXCEEDED"))
		return
	}
	//判断消耗
	costConfig := GetCsvMgr().GetTariffConfig(TARIFF_TYPE_ACTIVITY_BOSS_RESET, info.RefreshTimes+1)
	if costConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	// 物品不足
	if err := self.player.HasObjectOk(costConfig.ItemIds, costConfig.ItemNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	// 移除物品
	costItem := self.player.RemoveObjectLst(costConfig.ItemIds, costConfig.ItemNums, "世界BOSS次数", 0, 0, 0)

	info.Times = 0
	info.RefreshTimes++

	msgRel := &S2C_ActivityBossResetTimes{}
	msgRel.Cid = "activitybossresettimes"
	msgRel.CostItems = costItem
	msgRel.ActivityBossInfo = info
	self.player.Send(msgRel.Cid, msgRel)
}

// 获取战报信息
func (self *ModActivityBoss) GetBattleInfo(fightID int64) *BattleInfo {
	var battleInfo BattleInfo
	value, flag, err := HGetRedisEx(`san_activitybossbattleinfo`, fightID, fmt.Sprintf("%d", fightID))
	if err != nil || !flag {
		return GetServer().DBUser.GetBattleInfo(fightID)
	}
	if flag {
		err := json.Unmarshal([]byte(value), &battleInfo)
		if err != nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
			return &battleInfo
		}
	}

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ARENA_BATTLE_INFO, int(fightID), 0, 0, LOG_LANGUAGE_ARENA_BATTLE_INFO, 0, 0, self.player)

	if battleInfo.Id != 0 {
		return &battleInfo
	}
	return nil
}

// 获取战报信息
func (self *ModActivityBoss) GetBattleRecord(fightID int64) *BattleRecord {
	var battleRecord BattleRecord
	value, flag, err := HGetRedisEx(`san_activitybossbattlerecord`, fightID, fmt.Sprintf("%d", fightID))
	if err != nil || !flag {
		return GetServer().DBUser.GetBattleRecord(fightID)
	}
	if flag {
		err := json.Unmarshal([]byte(value), &battleRecord)
		if err != nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
			return &battleRecord
		}
	}
	if battleRecord.Id != 0 {
		return &battleRecord
	}
	return nil
}

func (self *ModActivityBoss) GetExchangeItem() map[int]int {
	data := make(map[int]int)
	for i := ACT_BOSS_MIN; i < ACT_BOSS_MAX; i++ {
		_, ok := self.Sql_ActivityBoss.activityBossInfo[i]
		if !ok {
			self.Sql_ActivityBoss.activityBossInfo[i] = new(JS_ActivityBossInfo)
			self.Sql_ActivityBoss.activityBossInfo[i].Id = i
		}

		isOpen, _ := GetActivityMgr().JudgeOpenAllIndex(self.Sql_ActivityBoss.activityBossInfo[i].Id, self.Sql_ActivityBoss.activityBossInfo[i].Id)
		if isOpen {
			configBoss := GetCsvMgr().GetActivityBossConfig(self.Sql_ActivityBoss.activityBossInfo[i].Id, self.Sql_ActivityBoss.activityBossInfo[i].Period)
			if configBoss == nil {
				continue
			}
			for j := 0; j < len(configBoss.Items); j++ {
				data[configBoss.Items[j]] = LOGIC_TRUE
			}
		}
	}
	return data
}
