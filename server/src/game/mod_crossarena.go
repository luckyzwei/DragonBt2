package game

import (
	"encoding/json"
	"fmt"
)

const (
	CROSSARENA_NEXT_CD = 10 //刷新间隔10秒
)

// 限时神将模块
type San_CrossArena struct {
	Uid           int64
	KeyId         int //! 活动KEY
	Subsection    int //当前大段位
	Class         int //当前小段位
	SubsectionMax int //最高大段位
	ClassMax      int //最高小段位
	Times         int //挑战次数
	BuyTimes      int //购买的次数
	StartTime     int64
	EndTime       int64
	ShowTime      int64
	TaskAwardSign string

	taskAwardSign map[int]int

	DataUpdate
}

type ModCrossArena struct {
	player         *Player
	Sql_CrossArena San_CrossArena
	nextTime       int64 //用来防止无限刷新
}

func (self *ModCrossArena) OnGetData(player *Player) {
	self.player = player
	sql := fmt.Sprintf("select * from `san_usercrossarena` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_CrossArena, "san_usercrossarena", self.player.ID)

	if self.Sql_CrossArena.Uid <= 0 {
		self.Sql_CrossArena.Uid = self.player.ID
		self.Sql_CrossArena.taskAwardSign = make(map[int]int, 0)
		self.Encode()
		InsertTable("san_usercrossarena", &self.Sql_CrossArena, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_CrossArena.Init("san_usercrossarena", &self.Sql_CrossArena, true)
}

//! 将数据库数据写入data
func (self *ModCrossArena) Decode() {
	json.Unmarshal([]byte(self.Sql_CrossArena.TaskAwardSign), &self.Sql_CrossArena.taskAwardSign)
}

//! 将data数据写入数据库
func (self *ModCrossArena) Encode() {
	self.Sql_CrossArena.TaskAwardSign = HF_JtoA(&self.Sql_CrossArena.taskAwardSign)
}

func (self *ModCrossArena) OnGetOtherData() {

}

// 注册消息
func (self *ModCrossArena) onReg(handlers map[string]func(body []byte)) {
	handlers["crossarenaadd"] = self.CrossArenaAdd
	handlers["crossarenagetdefencelist"] = self.GetDefenceList
	handlers["crossarenagetreward"] = self.GetReward
	handlers["crossarenagetrank"] = self.GetRank
	handlers["crossarenabuytimes"] = self.BuyTimes
	handlers["crossarenastartattack"] = self.Attack
	handlers["crossarenagetplayerinfo"] = self.CrossArenaGetPlayerInfo
	handlers["crossarenafightok"] = self.FightOK
	handlers["crossarenagetnow"] = self.CrossArenaGetNow
}

func (self *ModCrossArena) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModCrossArena) OnSave(sql bool) {
	self.Encode()
	self.Sql_CrossArena.Update(sql)
}

//每日任务刷新
func (self *ModCrossArena) OnRefresh() {
	self.Sql_CrossArena.Times = 0
	self.Sql_CrossArena.BuyTimes = 0
}

func (self *ModCrossArena) Check() {
	if self.Sql_CrossArena.taskAwardSign == nil {
		self.Sql_CrossArena.taskAwardSign = make(map[int]int)
	}
	if self.Sql_CrossArena.taskAwardSign == nil {
		self.Sql_CrossArena.taskAwardSign = make(map[int]int)
	}

	isOpen, _ := GetActivityMgr().JudgeOpenAllIndex(ACT_AREAN_CROSS_SERVER, ACT_AREAN_CROSS_SERVER)
	activity := GetActivityMgr().GetActivity(ACT_AREAN_CROSS_SERVER)
	if activity == nil {
		return
	}
	self.Sql_CrossArena.StartTime = HF_CalTimeForConfig(activity.info.Start, self.player.Sql_UserBase.Regtime)
	self.Sql_CrossArena.EndTime = self.Sql_CrossArena.StartTime + int64(activity.info.Continued)
	self.Sql_CrossArena.ShowTime = self.Sql_CrossArena.StartTime + int64(activity.info.Continued) + int64(activity.info.Show)
	if !isOpen {
		//活动关闭 补发之前没有领取的奖励
		itemMap := make(map[int]*Item)
		for _, v := range GetCsvMgr().CrossArenaSubsection {
			if self.Sql_CrossArena.taskAwardSign[v.Id] == LOGIC_TRUE {
				continue
			}
			if self.IsCanGet(v.Subsection, v.Class) {
				AddItemMapHelper(itemMap, v.Item, v.Num)
				self.Sql_CrossArena.taskAwardSign[v.Id] = LOGIC_TRUE
			}
		}

		//发送邮件
		if len(itemMap) > 0 {
			mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_CROSSARENA_TASK]
			if ok {
				itemLst := make([]PassItem, 0)
				for _, v := range itemMap {
					itemLst = append(itemLst, PassItem{ItemID: v.ItemId, Num: v.ItemNum})
				}
				self.player.GetModule("mail").(*ModMail).AddMail(1,
					1, 0, mailConfig.Mailtitle, mailConfig.Mailtxt, GetCsvMgr().GetText("STR_SYS"), itemLst, false, 0)
			}
		}
		return
	}
	keyId := GetActivityMgr().getActN3(ACT_AREAN_CROSS_SERVER)

	if self.Sql_CrossArena.KeyId != keyId {
		//刷新
		self.Sql_CrossArena.KeyId = keyId
		self.Sql_CrossArena.taskAwardSign = make(map[int]int)
		self.Sql_CrossArena.Subsection = 0
		self.Sql_CrossArena.Class = 0
		self.Sql_CrossArena.SubsectionMax = 0
		self.Sql_CrossArena.ClassMax = 0
		self.Sql_CrossArena.Times = 0
	}

	// 补充跨服竞技场 切
	teamPos := self.player.getTeamPosByType(TEAMTYPE_CROSSARENA_DEFENCE)
	if nil == teamPos || teamPos.isUIPosEmpty() {
		herouse := []int{}
		var msg C2S_AddTeamUIPos
		msg.TeamType = TEAMTYPE_CROSSARENA_DEFENCE
		// 默认上阵最强阵容 且不能有重复英雄
		heros := self.player.GetModule("hero").(*ModHero).GetBestFormat2()
		if len(heros) <= 0 {
			return
		}
		nCount := 0
		for _, v := range heros {
			// 超过最大则跳出
			if nCount >= MAX_FIGHT_POS {
				break
			}
			// 找不到英雄
			hero := self.player.GetModule("hero").(*ModHero).GetHero(v)
			if nil == hero {
				continue
			}
			// 如果是已经共鸣的虚空英雄
			if hero.VoidHero != 0 && hero.Resonance != 0 {
				continue
			}
			// 该类英雄已经使用
			find := false
			for _, t := range herouse {
				if hero.HeroId == t {
					find = true
				}
			}
			if find {
				continue
			}
			msg.FightPos = append(msg.FightPos, v)
			herouse = append(herouse, hero.HeroId)
			nCount++
		}

		nLen := len(msg.FightPos)
		if nLen < MAX_FIGHT_POS {
			for i := nLen; i < MAX_FIGHT_POS; i++ {
				msg.FightPos = append(msg.FightPos, 0)
			}
		}
		smsg, _ := json.Marshal(&msg)
		self.player.GetModule("team").(*ModTeam).addUIPos(smsg)
	}
}

//判断自己的当前最大值是否超过参数
func (self *ModCrossArena) IsCanGet(subsection int, class int) bool {
	if self.Sql_CrossArena.SubsectionMax == 0 || self.Sql_CrossArena.ClassMax == 0 {
		return false
	}
	if self.Sql_CrossArena.SubsectionMax < subsection {
		return true
	}
	if self.Sql_CrossArena.SubsectionMax > subsection {
		return false
	}
	return self.Sql_CrossArena.ClassMax <= class
}

func (self *ModCrossArena) SendInfo() {
	self.Check()

	//获取自己的信息和排行榜前50
	top, info := GetCrossArenaMgr().GetSendInfo(self.player)
	if info != nil {
		self.Sql_CrossArena.Subsection = info.Subsection
		self.Sql_CrossArena.Class = info.Class
	}
	var msg S2C_CrossArenaInfo
	msg.Cid = "crossarenainfo"
	msg.Top = top
	msg.SelfInfo = info
	msg.SubsectionMax = self.Sql_CrossArena.SubsectionMax
	msg.ClassMax = self.Sql_CrossArena.ClassMax
	msg.Times = self.Sql_CrossArena.Times
	msg.BuyTimes = self.Sql_CrossArena.BuyTimes
	msg.StartTime = self.Sql_CrossArena.StartTime
	msg.EndTime = self.Sql_CrossArena.EndTime
	msg.ShowTime = self.Sql_CrossArena.ShowTime
	msg.TaskAwardSign = self.Sql_CrossArena.taskAwardSign
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

//参加竞技场
func (self *ModCrossArena) CrossArenaAdd(body []byte) {
	fightInfo := GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_DEFENCE)
	info := GetCrossArenaMgr().AddInfo(self.player, fightInfo)

	var msg S2C_CrossArenaAdd
	msg.Cid = "crossarenaadd"
	msg.SelfInfo = info
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModCrossArena) GetDefenceList(body []byte) {
	if self.IsClose() {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_CITY_NOT_OPENED"))
		return
	}

	if self.nextTime > TimeServer().Unix() {
		self.player.SendErr(GetCsvMgr().GetText("STR_DUNGEON_TEAM_CD"))
		return
	}

	info, fightInfo, ret := GetCrossArenaMgr().GetDefenceList(self.player)
	if ret != RETCODE_OK {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_BEAUTY_DATA_ABNORMITY"))
		return
	}

	//self.nextTime = TimeServer().Unix() + CROSSARENA_NEXT_CD

	var msg S2C_CrossArenaGetDefenceList
	msg.Cid = "crossarenagetdefencelist"
	//验证下战斗数据,并根据需求生成机器人
	for _, v := range info {
		//机器人就生成
		if v.Robot == LOGIC_TRUE {
			for _, config := range GetCsvMgr().JJCRobotConfig {
				if config.Type != 3 {
					continue
				}
				if config.Jjcclass == v.Subsection && config.Jjcdan == v.Class {
					robotFightInfo := GetCsvMgr().GetRobot(config)
					if robotFightInfo != nil {
						robotFightInfo.Uname = v.UName
						msg.Info = append(msg.Info, v)
						msg.FightInfo = append(msg.FightInfo, robotFightInfo)
					}
					break
				}
			}
		} else {
			for _, info := range fightInfo {
				if info.Uid == v.Uid {
					msg.Info = append(msg.Info, v)
					msg.FightInfo = append(msg.FightInfo, info)
					break
				}
			}
		}
	}
	msg.NextTime = self.nextTime
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModCrossArena) GetReward(body []byte) {
	var msg C2S_CrossArenaGetReward
	json.Unmarshal(body, &msg)

	config := GetCsvMgr().GetCrossArenaSubsection(msg.Taskid)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_DROP_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	if self.Sql_CrossArena.taskAwardSign[msg.Taskid] == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_AWARD_FOR_THE_EVENT_HAS"))
		return
	}

	if !self.IsCanGet(config.Subsection, config.Class) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DIPLOMACY_THE_TASK_IS_NOT_COMPLETED"))
		return
	}

	self.Sql_CrossArena.taskAwardSign[msg.Taskid] = LOGIC_TRUE
	item := self.player.AddObjectLst(config.Item, config.Num, "至尊竞技任务", config.Id, 0, 0)

	var msgRel S2C_CrossArenaGetReward
	msgRel.Cid = "crossarenagetreward"
	msgRel.TaskSign = self.Sql_CrossArena.taskAwardSign
	msgRel.GetItems = item

	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModCrossArena) GetRank(body []byte) {
	//获取自己的信息和排行榜前50
	top, info := GetCrossArenaMgr().GetSendInfo(self.player)
	if info != nil {
		self.Sql_CrossArena.Subsection = info.Subsection
		self.Sql_CrossArena.Class = info.Class
	}
	var msg S2C_CrossArenaGetRank
	msg.Cid = "crossarenagetrank"
	msg.Top = top
	msg.SelfInfo = info
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModCrossArena) BuyTimes(body []byte) {
	if self.IsClose() {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_CITY_NOT_OPENED"))
		return
	}

	configCost := GetCsvMgr().GetTariffConfig(TARIFF_TYPE_CROSS_ARENA_TIMES, self.Sql_CrossArena.BuyTimes+1)
	if configCost == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	//看消耗够不够
	if err := self.player.HasObjectOk(configCost.ItemIds, configCost.ItemNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	item := self.player.RemoveObjectLst(configCost.ItemIds, configCost.ItemNums, "跨服竞技场购买", 0, 0, 0)
	self.Sql_CrossArena.BuyTimes += 1

	var msgRel S2C_CrossArenaBuyTimes
	msgRel.Cid = "crossarenabuytimes"
	msgRel.CostItem = item
	msgRel.BuyTimes = self.Sql_CrossArena.BuyTimes
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModCrossArena) Attack(body []byte) {
	if self.IsClose() {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_CITY_NOT_OPENED"))
		return
	}

	var msg C2S_CrossArenaAttack
	json.Unmarshal(body, &msg)

	//次数处理
	freeTimes := GetCsvMgr().getInitNum(CROSSARENA_FREETIMES)
	if self.Sql_CrossArena.Times >= freeTimes+self.Sql_CrossArena.BuyTimes {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_REACH_THE_UPPER_LIMIT"))
		return
	}

	info, fightInfo, _ := GetCrossArenaMgr().GetInfo(msg.AttackUid)
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_NOT_EXIST"))
		return
	}
	if info.Subsection != msg.AttackSubsection || info.Class != msg.AttackClass {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_HAS_CHANGE"))
		var msg S2C_CrossArenaArenaAttack
		msg.Cid = "crossarenastartattackfail"
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
		return
	}
	//如果是不是机器人，但是没防守阵容则报错
	if info.Robot == LOGIC_FALSE && fightInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_TEAM_ERROR"))
		return
	}
	//开打
	if info.Robot == LOGIC_TRUE {
		//生成机器人阵容
		for _, config := range GetCsvMgr().JJCRobotConfig {
			if config.Type != 3 {
				continue
			}
			if config.Jjcclass == info.Subsection && config.Jjcdan == info.Class {
				fightInfo = GetCsvMgr().GetRobot(config)
				fightInfo.Uid = info.Uid
				fightInfo.Iconid = info.Icon

				if fightInfo != nil {
					fightInfo.Uname = info.UName
				} else {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_TEAM_ERROR"))
				}
				break
			}
		}
	}

	attack := GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_ATTACK)
	defence := fightInfo
	randNum := TimeServer().Unix()
	//战报加入战斗服务器等待队列
	fightId := GetArenaMgr().AddFightListForCross(BATTLE_TYPE_RECORD_CROSSARENA,
		attack,
		defence,
		randNum,
		randNum)

	if fightId > 0 {
		//开始战斗
		var msgRel S2C_CrossArenaArenaAttack
		msgRel.Cid = "crossarenastartattack"
		msgRel.FightId = fightId
		msgRel.Attack = attack
		msgRel.Defence = defence
		msgRel.RandNum = randNum
		self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	}
}

func (self *ModCrossArena) CrossArenaGetPlayerInfo(body []byte) {
	var msg C2S_CrossArenaGetPlayerInfo
	json.Unmarshal(body, &msg)

	info, fightInfo, liffTreeInfo := GetCrossArenaMgr().GetInfo(msg.PlayerUid)
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_NOT_EXIST"))
		return
	}
	//如果是不是机器人，但是没防守阵容则报错
	if info.Robot == LOGIC_FALSE && fightInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_TEAM_ERROR"))
		return
	}
	//开打
	if info.Robot == LOGIC_TRUE {
		//生成机器人阵容
		for _, config := range GetCsvMgr().JJCRobotConfig {
			if config.Type != 3 {
				continue
			}
			if config.Jjcclass == info.Subsection && config.Jjcdan == info.Class {
				fightInfo = GetCsvMgr().GetRobot(config)
				if fightInfo != nil {
					fightInfo.Uname = info.UName
					fightInfo.Deffight = info.Fight
					fightInfo.Iconid = info.Icon
				} else {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_TEAM_ERROR"))
					return
				}
				break
			}
		}
	}

	var msgRel S2C_CrossArenaGetPlayerInfo
	msgRel.Cid = "crossarenagetplayerinfo"
	msgRel.Info = info
	msgRel.FightInfo = fightInfo
	msgRel.LifeTreeInfo = liffTreeInfo
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModCrossArena) CrossArenaGetNow(body []byte) {
	var msg C2S_CrossArenaGetPlayerInfo
	json.Unmarshal(body, &msg)

	info, fightInfo, liffTreeInfo := GetCrossArenaMgr().GetInfo(msg.PlayerUid)
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_NOT_EXIST"))
		return
	}
	//如果是不是机器人，但是没防守阵容则报错
	if info.Robot == LOGIC_FALSE && fightInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_TEAM_ERROR"))
		return
	}
	//开打
	if info.Robot == LOGIC_TRUE {
		//生成机器人阵容
		for _, config := range GetCsvMgr().JJCRobotConfig {
			if config.Type != 3 {
				continue
			}
			if config.Jjcclass == info.Subsection && config.Jjcdan == info.Class {
				fightInfo = GetCsvMgr().GetRobot(config)
				if fightInfo != nil {
					fightInfo.Uname = info.UName
					fightInfo.Deffight = info.Fight
					fightInfo.Iconid = info.Icon

					isSet := false
					for k, _ := range fightInfo.Heroinfo {
						if !isSet {
							fightInfo.Heroinfo[k].Fight = fightInfo.Deffight
							isSet = true
						} else {
							fightInfo.Heroinfo[k].Fight = 0
						}
					}
				} else {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA_USER_TEAM_ERROR"))
					return
				}
				break
			}
		}
	}

	var msgRel S2C_CrossArenaGetPlayerInfo
	msgRel.Cid = "crossarenagetnow"
	msgRel.Info = info
	msgRel.FightInfo = fightInfo
	msgRel.LifeTreeInfo = liffTreeInfo
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModCrossArena) FightEndOK(player *Player, battleInfo BattleInfo, attack *JS_FightInfo, defence *JS_FightInfo) {

	res := GetCrossArenaMgr().FightEnd(self.player, attack, defence, battleInfo)

	top, info := GetCrossArenaMgr().GetSendInfo(self.player)

	if info != nil {
		if self.Sql_CrossArena.SubsectionMax == 0 || self.Sql_CrossArena.ClassMax == 0 {
			self.Sql_CrossArena.SubsectionMax = info.Subsection
			self.Sql_CrossArena.ClassMax = info.Class
		}
		if !self.IsCanGet(info.Subsection, info.Class) {
			self.Sql_CrossArena.SubsectionMax = info.Subsection
			self.Sql_CrossArena.ClassMax = info.Class
		}
	}

	if res != nil {
		self.Sql_CrossArena.Times += 1

		var msg S2C_CrossArenaFightOK
		msg.Cid = "crossarenafightok"
		msg.Top = top
		msg.SelfInfo = info
		msg.OldFightId = battleInfo.Id
		msg.NewFightId = res.NewFightId
		msg.Result = res.Result
		msg.SubsectionMax = self.Sql_CrossArena.SubsectionMax
		msg.ClassMax = self.Sql_CrossArena.ClassMax
		msg.Times = self.Sql_CrossArena.Times
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
	}
}

//这个是客户端发送到结果
func (self *ModCrossArena) FightOK(body []byte) {
	if self.IsClose() {
		return
	}

	var msg C2S_CrossArenaFightResult
	json.Unmarshal(body, &msg)
	GetArenaMgr().ArenaFightResultByCross(msg.Type, msg.BattleInfo)
}

func (self *ModCrossArena) GetBattleInfo(key int64) *BattleInfo {
	return GetCrossArenaMgr().GetBattleInfo(key)
}

func (self *ModCrossArena) GetBattleRecord(key int64) *BattleRecord {
	return GetCrossArenaMgr().GetBattleRecord(key)
}

func (self *ModCrossArena) UpdateFormat() {
	teamPos := self.player.getTeamPosByType(TEAMTYPE_CROSSARENA_DEFENCE)
	if nil == teamPos || teamPos.isUIPosEmpty() {
		self.Check()
		return
	}
}

//判断活动机能的开启，不包括展示时间
func (self *ModCrossArena) IsClose() bool {
	now := TimeServer().Unix()

	if now >= self.Sql_CrossArena.StartTime && now <= self.Sql_CrossArena.EndTime {
		return false
	}
	return true
}
