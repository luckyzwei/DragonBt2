package game

import (
	"encoding/json"
	"fmt"
)

const (
	CrossArena3V3_NEXT_CD = 10 //刷新间隔10秒
)

// 限时神将模块
type San_CrossArena3V3 struct {
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

type ModCrossArena3V3 struct {
	player            *Player
	Sql_CrossArena3V3 San_CrossArena3V3
	nextTime          int64 //用来防止无限刷新
}

func (self *ModCrossArena3V3) OnGetData(player *Player) {
	self.player = player
	sql := fmt.Sprintf("select * from `san_usercrossarena3v3` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_CrossArena3V3, "san_usercrossarena3v3", self.player.ID)

	if self.Sql_CrossArena3V3.Uid <= 0 {
		self.Sql_CrossArena3V3.Uid = self.player.ID
		self.Sql_CrossArena3V3.taskAwardSign = make(map[int]int, 0)
		self.Encode()
		InsertTable("san_usercrossarena3v3", &self.Sql_CrossArena3V3, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_CrossArena3V3.Init("san_usercrossarena3v3", &self.Sql_CrossArena3V3, true)
}

//! 将数据库数据写入data
func (self *ModCrossArena3V3) Decode() {
	json.Unmarshal([]byte(self.Sql_CrossArena3V3.TaskAwardSign), &self.Sql_CrossArena3V3.taskAwardSign)
}

//! 将data数据写入数据库
func (self *ModCrossArena3V3) Encode() {
	self.Sql_CrossArena3V3.TaskAwardSign = HF_JtoA(&self.Sql_CrossArena3V3.taskAwardSign)
}

func (self *ModCrossArena3V3) OnGetOtherData() {

}

// 注册消息
func (self *ModCrossArena3V3) onReg(handlers map[string]func(body []byte)) {
	handlers["crossarena3v3add"] = self.CrossArena3V3Add
	handlers["crossarena3v3getdefencelist"] = self.GetDefenceList
	handlers["crossarena3v3getreward"] = self.GetReward
	handlers["crossarena3v3getrank"] = self.GetRank
	handlers["crossarena3v3buytimes"] = self.BuyTimes
	handlers["crossarena3v3startattack"] = self.Attack
	handlers["crossarena3v3getplayerinfo"] = self.CrossArena3V3GetPlayerInfo
	handlers["crossarena3v3fightok"] = self.FightOK
	handlers["crossarena3v3getnow"] = self.CrossArena3V3GetNow
}

func (self *ModCrossArena3V3) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModCrossArena3V3) OnSave(sql bool) {
	self.Encode()
	self.Sql_CrossArena3V3.Update(sql)
}

//每日任务刷新
func (self *ModCrossArena3V3) OnRefresh() {
	self.Sql_CrossArena3V3.Times = 0
	self.Sql_CrossArena3V3.BuyTimes = 0
}

func (self *ModCrossArena3V3) Check() {
	if self.Sql_CrossArena3V3.taskAwardSign == nil {
		self.Sql_CrossArena3V3.taskAwardSign = make(map[int]int)
	}
	if self.Sql_CrossArena3V3.taskAwardSign == nil {
		self.Sql_CrossArena3V3.taskAwardSign = make(map[int]int)
	}

	isOpen, _ := GetActivityMgr().JudgeOpenAllIndex(ACT_AREAN_CROSS_SERVER_3V3, ACT_AREAN_CROSS_SERVER_3V3)
	activity := GetActivityMgr().GetActivity(ACT_AREAN_CROSS_SERVER_3V3)
	if activity == nil {
		return
	}
	self.Sql_CrossArena3V3.StartTime = HF_CalTimeForConfig(activity.info.Start, self.player.Sql_UserBase.Regtime)
	self.Sql_CrossArena3V3.EndTime = self.Sql_CrossArena3V3.StartTime + int64(activity.info.Continued)
	self.Sql_CrossArena3V3.ShowTime = self.Sql_CrossArena3V3.StartTime + int64(activity.info.Continued) + int64(activity.info.Show)
	if !isOpen {
		//活动关闭 补发之前没有领取的奖励
		itemMap := make(map[int]*Item)
		for _, v := range GetCsvMgr().CrossArena3V3Subsection {
			if self.Sql_CrossArena3V3.taskAwardSign[v.Id] == LOGIC_TRUE {
				continue
			}
			if self.IsCanGet(v.Subsection, v.Class) {
				AddItemMapHelper(itemMap, v.Item, v.Num)
				self.Sql_CrossArena3V3.taskAwardSign[v.Id] = LOGIC_TRUE
			}
		}

		//发送邮件
		if len(itemMap) > 0 {
			mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_CROSSARENA_3V3_TASK]
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
	keyId := GetActivityMgr().getActN3(ACT_AREAN_CROSS_SERVER_3V3)

	if self.Sql_CrossArena3V3.KeyId != keyId {
		//刷新
		self.Sql_CrossArena3V3.KeyId = keyId
		self.Sql_CrossArena3V3.taskAwardSign = make(map[int]int)
		self.Sql_CrossArena3V3.Subsection = 0
		self.Sql_CrossArena3V3.Class = 0
		self.Sql_CrossArena3V3.SubsectionMax = 0
		self.Sql_CrossArena3V3.ClassMax = 0
		self.Sql_CrossArena3V3.Times = 0
	}

	herouse1 := []int{}
	teamType1 := []int{TEAMTYPE_CROSSARENA_ATTACK_3V3_1, TEAMTYPE_CROSSARENA_ATTACK_3V3_2, TEAMTYPE_CROSSARENA_ATTACK_3V3_3}
	refresh := false
	for _, team := range teamType1 {
		teamPos := self.player.getTeamPosByType(team)
		if nil == teamPos || teamPos.isUIPosEmpty() {
			refresh = true
			break
		}
	}

	if refresh {
		var msg C2S_AddTeamUIPos
		msg.TeamType = TEAMTYPE_CROSSARENA_ATTACK_3V3_1
		count := 0
		heros := self.player.GetModule("hero").(*ModHero).GetBestFormat2()
		if len(heros) <= 0 {
			return
		}
		for _, v := range heros {
			if count >= MAX_FIGHT_POS*3 {
				break
			}
			hero := self.player.GetModule("hero").(*ModHero).GetHero(v)
			if nil == hero {
				continue
			}
			// 如果是已经共鸣的虚空英雄
			if hero.VoidHero != 0 && hero.Resonance != 0 {
				continue
			}
			find := false
			for _, t := range herouse1 {
				if hero.HeroId == t {
					find = true
					break
				}
			}
			if find {
				continue
			}
			msg.FightPos = append(msg.FightPos, v)
			herouse1 = append(herouse1, hero.HeroId)
			count++
		}
		nLen := len(msg.FightPos)
		if nLen <= MAX_FIGHT_POS*2 {
			temp := msg.FightPos
			msg.FightPos = []int{}
			for i := 0; i < MAX_FIGHT_POS*3; i++ {
				if i%MAX_FIGHT_POS == 0 && len(temp) > 0 {
					msg.FightPos = append(msg.FightPos, temp[0])
					temp = temp[1:]
				} else {
					msg.FightPos = append(msg.FightPos, 0)
				}
			}
			if len(temp) > 0 {
				for i, g := range msg.FightPos {
					if g == 0 && len(temp) > 0 {
						msg.FightPos[i] = temp[0]
						temp = temp[1:]
					}
				}
			}
		} else if nLen < MAX_FIGHT_POS*3 {
			for i := nLen; i < MAX_FIGHT_POS*3; i++ {
				msg.FightPos = append(msg.FightPos, 0)
			}
		}
		smsg, _ := json.Marshal(&msg)
		self.player.GetModule("team").(*ModTeam).AddCrossArena3V3(smsg)
	}
	herouse2 := []int{}
	teamType2 := []int{TEAMTYPE_CROSSARENA_DEFENCE_3V3_1, TEAMTYPE_CROSSARENA_DEFENCE_3V3_2, TEAMTYPE_CROSSARENA_DEFENCE_3V3_3}
	refresh = false
	for _, team := range teamType2 {
		teamPos := self.player.getTeamPosByType(team)
		if nil == teamPos || teamPos.isUIPosEmpty() {
			refresh = true
			break
		}
	}
	if refresh {
		var msg C2S_AddTeamUIPos
		msg.TeamType = TEAMTYPE_CROSSARENA_DEFENCE_3V3_1
		heros := self.player.GetModule("hero").(*ModHero).GetBestFormat2()
		count := 0
		for _, v := range heros {
			if count >= MAX_FIGHT_POS*3 {
				break
			}
			hero := self.player.GetModule("hero").(*ModHero).GetHero(v)
			if nil == hero {
				continue
			}
			// 如果是已经共鸣的虚空英雄
			if hero.VoidHero != 0 && hero.Resonance != 0 {
				continue
			}
			find := false
			for _, t := range herouse2 {
				if hero.HeroId == t {
					find = true
					break
				}
			}
			if find {
				continue
			}
			msg.FightPos = append(msg.FightPos, v)
			herouse2 = append(herouse2, hero.HeroId)
			count++
		}
		nLen := len(msg.FightPos)
		if nLen <= MAX_FIGHT_POS*2 {
			temp := msg.FightPos
			msg.FightPos = []int{}
			for i := 0; i < MAX_FIGHT_POS*3; i++ {
				if i%MAX_FIGHT_POS == 0 && len(temp) > 0 {
					msg.FightPos = append(msg.FightPos, temp[0])
					temp = temp[1:]
				} else {
					msg.FightPos = append(msg.FightPos, 0)
				}
			}
			if len(temp) > 0 {
				for i, g := range msg.FightPos {
					if g == 0 && len(temp) > 0 {
						msg.FightPos[i] = temp[0]
						temp = temp[1:]
					}
				}
			}
		} else if nLen < MAX_FIGHT_POS*3 {
			for i := nLen; i < MAX_FIGHT_POS*3; i++ {
				msg.FightPos = append(msg.FightPos, 0)
			}
		}
		smsg, _ := json.Marshal(&msg)
		self.player.GetModule("team").(*ModTeam).AddCrossArena3V3(smsg)
	}
}

//判断自己的当前最大值是否超过参数
func (self *ModCrossArena3V3) IsCanGet(subsection int, class int) bool {
	if self.Sql_CrossArena3V3.SubsectionMax == 0 || self.Sql_CrossArena3V3.ClassMax == 0 {
		return false
	}
	if self.Sql_CrossArena3V3.SubsectionMax < subsection {
		return true
	}
	if self.Sql_CrossArena3V3.SubsectionMax > subsection {
		return false
	}
	return self.Sql_CrossArena3V3.ClassMax <= class
}

func (self *ModCrossArena3V3) SendInfo() {
	self.Check()

	//获取自己的信息和排行榜前50
	top, info := GetCrossArena3V3Mgr().GetSendInfo(self.player)
	if info != nil {
		self.Sql_CrossArena3V3.Subsection = info.Subsection
		self.Sql_CrossArena3V3.Class = info.Class
	}
	var msg S2C_CrossArena3V3Info
	msg.Cid = "crossarena3v3info"
	msg.Top = top
	msg.SelfInfo = info
	msg.SubsectionMax = self.Sql_CrossArena3V3.SubsectionMax
	msg.ClassMax = self.Sql_CrossArena3V3.ClassMax
	msg.Times = self.Sql_CrossArena3V3.Times
	msg.BuyTimes = self.Sql_CrossArena3V3.BuyTimes
	msg.StartTime = self.Sql_CrossArena3V3.StartTime
	msg.EndTime = self.Sql_CrossArena3V3.EndTime
	msg.ShowTime = self.Sql_CrossArena3V3.ShowTime
	msg.TaskAwardSign = self.Sql_CrossArena3V3.taskAwardSign
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

//参加竞技场
func (self *ModCrossArena3V3) CrossArena3V3Add(body []byte) {
	fightInfos := [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo{}
	fightInfos[0] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_DEFENCE_3V3_1)
	fightInfos[1] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_DEFENCE_3V3_2)
	fightInfos[2] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_DEFENCE_3V3_3)

	info := GetCrossArena3V3Mgr().AddInfo(self.player, fightInfos)

	var msg S2C_CrossArena3V3Add
	msg.Cid = "crossarena3v3add"
	msg.SelfInfo = info
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModCrossArena3V3) GetDefenceList(body []byte) {
	if self.IsClose() {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_CITY_NOT_OPENED"))
		return
	}

	if self.nextTime > TimeServer().Unix() {
		self.player.SendErr(GetCsvMgr().GetText("STR_DUNGEON_TEAM_CD"))
		return
	}

	info, fightInfo, ret := GetCrossArena3V3Mgr().GetDefenceList(self.player)
	if ret != RETCODE_OK {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_BEAUTY_DATA_ABNORMITY"))
		return
	}

	var msg S2C_CrossArena3V3GetDefenceList
	msg.Cid = "crossarena3v3getdefencelist"
	//验证下战斗数据,并根据需求生成机器人
	for _, v := range info {
		//机器人就生成
		if v.Robot == LOGIC_TRUE {
			fightInfos := self.GetRobotTeam3V3(v)
			fightLstTemp := make([]*JS_FightInfo, 0)
			for _, v := range fightInfos {
				fightLstTemp = append(fightLstTemp, v)
			}
			msg.Info = append(msg.Info, v)
			msg.FightInfo = append(msg.FightInfo, fightLstTemp)
		} else {
			for _, info := range fightInfo {
				fightLstTemp := make([]*JS_FightInfo, 0)
				for _, infoi := range info {
					if infoi.Uid == v.Uid {
						fightLstTemp = append(fightLstTemp, infoi)
					}
				}
				if len(fightLstTemp) == CROSSARENA3V3_TEAM_MAX {
					msg.Info = append(msg.Info, v)
					msg.FightInfo = append(msg.FightInfo, fightLstTemp)
				}
			}
		}
	}
	msg.NextTime = self.nextTime
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModCrossArena3V3) GetReward(body []byte) {
	var msg C2S_CrossArena3V3GetReward
	json.Unmarshal(body, &msg)

	config := GetCsvMgr().GetCrossArena3V3Subsection(msg.Taskid)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_DROP_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	if self.Sql_CrossArena3V3.taskAwardSign[msg.Taskid] == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_AWARD_FOR_THE_EVENT_HAS"))
		return
	}

	if !self.IsCanGet(config.Subsection, config.Class) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DIPLOMACY_THE_TASK_IS_NOT_COMPLETED"))
		return
	}

	self.Sql_CrossArena3V3.taskAwardSign[msg.Taskid] = LOGIC_TRUE
	item := self.player.AddObjectLst(config.Item, config.Num, "至尊竞技任务", config.Id, 0, 0)

	var msgRel S2C_CrossArena3V3GetReward
	msgRel.Cid = "crossarena3v3getreward"
	msgRel.TaskSign = self.Sql_CrossArena3V3.taskAwardSign
	msgRel.GetItems = item

	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModCrossArena3V3) GetRank(body []byte) {
	//获取自己的信息和排行榜前50
	top, info := GetCrossArena3V3Mgr().GetSendInfo(self.player)
	if info != nil {
		self.Sql_CrossArena3V3.Subsection = info.Subsection
		self.Sql_CrossArena3V3.Class = info.Class
	}
	var msg S2C_CrossArena3V3GetRank
	msg.Cid = "crossarena3v3getrank"
	msg.Top = top
	msg.SelfInfo = info
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModCrossArena3V3) BuyTimes(body []byte) {
	if self.IsClose() {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_CITY_NOT_OPENED"))
		return
	}

	configCost := GetCsvMgr().GetTariffConfig(TARIFF_TYPE_CROSS_ARENA_TIMES, self.Sql_CrossArena3V3.BuyTimes+1)
	if configCost == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	//看消耗够不够
	if err := self.player.HasObjectOk(configCost.ItemIds, configCost.ItemNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	item := self.player.RemoveObjectLst(configCost.ItemIds, configCost.ItemNums, "跨服竞技场33购买", 0, 0, 0)
	self.Sql_CrossArena3V3.BuyTimes += 1

	var msgRel S2C_CrossArena3V3BuyTimes
	msgRel.Cid = "crossarena3v3buytimes"
	msgRel.CostItem = item
	msgRel.BuyTimes = self.Sql_CrossArena3V3.BuyTimes
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModCrossArena3V3) Attack(body []byte) {
	if self.IsClose() {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_CITY_NOT_OPENED"))
		return
	}

	var msg C2S_CrossArena3V3Attack
	json.Unmarshal(body, &msg)

	//次数处理
	freeTimes := GetCsvMgr().getInitNum(CROSSARENA_FREETIMES)
	if self.Sql_CrossArena3V3.Times >= freeTimes+self.Sql_CrossArena3V3.BuyTimes {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_REACH_THE_UPPER_LIMIT"))
		return
	}

	info, fightInfo, _ := GetCrossArena3V3Mgr().GetInfo(msg.AttackUid)
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CrossArena3V3_USER_NOT_EXIST"))
		return
	}
	if info.Subsection != msg.AttackSubsection || info.Class != msg.AttackClass {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CrossArena3V3_USER_HAS_CHANGE"))
		var msg S2C_CrossArena3V3ArenaAttack
		msg.Cid = "crossarena3v3startattackfail"
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
		return
	}
	//如果是不是机器人，但是没防守阵容则报错
	if info.Robot == LOGIC_FALSE {
		for _, v := range fightInfo {
			if v == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA3V3_USER_TEAM_ERROR"))
				return
			}
		}
	}
	//开打
	if info.Robot == LOGIC_TRUE {
		fightInfo = self.GetRobotTeam3V3(info)
	}

	now := TimeServer().Unix()
	myFight := [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo{}
	myFight[0] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_ATTACK_3V3_1)
	myFight[1] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_ATTACK_3V3_2)
	myFight[2] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_ATTACK_3V3_3)
	fightId := GetArenaSpecialMgr().AddFightListForCross(self.player, myFight,
		fightInfo,
		now,
		now)

	if fightId[0] > 0 {
		//开始战斗
		var msgRel S2C_CrossArena3V3ArenaAttack
		msgRel.Cid = "crossarena3v3startattack"
		msgRel.FightId = fightId
		msgRel.Attack = myFight
		msgRel.Defence = fightInfo
		msgRel.RandNum = now
		self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	}
}

func (self *ModCrossArena3V3) CrossArena3V3GetPlayerInfo(body []byte) {
	var msg C2S_CrossArena3V3GetPlayerInfo
	json.Unmarshal(body, &msg)

	info, fightInfo, liffTreeInfo := GetCrossArena3V3Mgr().GetInfo(msg.PlayerUid)
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CrossArena3V3_USER_NOT_EXIST"))
		return
	}
	//如果是不是机器人，但是没防守阵容则报错
	if info.Robot == LOGIC_FALSE {
		for _, v := range fightInfo {
			if v == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA3V3_USER_TEAM_ERROR"))
				return
			}
		}
	}
	//开打
	if info.Robot == LOGIC_TRUE {
		fightInfo = self.GetRobotTeam3V3(info)
	}

	var msgRel S2C_CrossArena3V3GetPlayerInfo
	msgRel.Cid = "crossarena3v3getplayerinfo"
	msgRel.Info = info
	msgRel.FightInfo = fightInfo
	msgRel.LifeTreeInfo = liffTreeInfo
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModCrossArena3V3) CrossArena3V3GetNow(body []byte) {
	var msg C2S_CrossArena3V3GetPlayerInfo
	json.Unmarshal(body, &msg)

	info, fightInfo, liffTreeInfo := GetCrossArena3V3Mgr().GetInfo(msg.PlayerUid)
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CrossArena3V3_USER_NOT_EXIST"))
		return
	}
	//如果是不是机器人，但是没防守阵容则报错
	if info.Robot == LOGIC_FALSE {
		for _, v := range fightInfo {
			if v == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CROSSARENA3V3_USER_TEAM_ERROR"))
				return
			}
		}
	}
	//开打
	if info.Robot == LOGIC_TRUE {
		fightInfo = self.GetRobotTeam3V3(info)
	}

	var msgRel S2C_CrossArena3V3GetPlayerInfo
	msgRel.Cid = "crossarena3v3getnow"
	msgRel.Info = info
	msgRel.FightInfo = fightInfo
	msgRel.LifeTreeInfo = liffTreeInfo
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModCrossArena3V3) FightEndOK(player *Player, result int, battleInfo *ArenaSpecialFightList) {

	res := GetCrossArena3V3Mgr().FightEnd(self.player, result, battleInfo)

	top, info := GetCrossArena3V3Mgr().GetSendInfo(self.player)

	if info != nil {
		if self.Sql_CrossArena3V3.SubsectionMax == 0 || self.Sql_CrossArena3V3.ClassMax == 0 {
			self.Sql_CrossArena3V3.SubsectionMax = info.Subsection
			self.Sql_CrossArena3V3.ClassMax = info.Class
		}
		if !self.IsCanGet(info.Subsection, info.Class) {
			self.Sql_CrossArena3V3.SubsectionMax = info.Subsection
			self.Sql_CrossArena3V3.ClassMax = info.Class
		}
	}

	if res != nil {
		self.Sql_CrossArena3V3.Times += 1

		var msg S2C_CrossArena3V3FightOK
		msg.Cid = "crossarena3v3fightok"
		msg.Top = top
		msg.SelfInfo = info
		msg.OldFightId = battleInfo.FightId
		msg.NewFightId = res.NewFightId
		msg.Result = res.Result
		msg.SubsectionMax = self.Sql_CrossArena3V3.SubsectionMax
		msg.ClassMax = self.Sql_CrossArena3V3.ClassMax
		msg.Times = self.Sql_CrossArena3V3.Times
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
	}
}

//这个是客户端发送到结果
func (self *ModCrossArena3V3) FightOK(body []byte) {
	if self.IsClose() {
		return
	}

	var msg C2S_CrossArena3V3FightResult
	json.Unmarshal(body, &msg)
	GetArenaSpecialMgr().ArenaFightResultByCross(msg.Type, msg.BattleInfo)
}

func (self *ModCrossArena3V3) GetBattleInfo(key int64) []*BattleInfo {
	return GetCrossArena3V3Mgr().GetBattleInfo(key)
}

func (self *ModCrossArena3V3) GetBattleRecord(key int64) *BattleRecord {
	return GetCrossArena3V3Mgr().GetBattleRecord(key)
}

func (self *ModCrossArena3V3) UpdateFormat() {
	for i := TEAMTYPE_CROSSARENA_DEFENCE_3V3_1; i <= TEAMTYPE_CROSSARENA_DEFENCE_3V3_3; i++ {
		teamPos := self.player.getTeamPosByType(i)
		if nil == teamPos || teamPos.isUIPosEmpty() {
			self.Check()
			return
		}
	}
}

//判断活动机能的开启，不包括展示时间
func (self *ModCrossArena3V3) IsClose() bool {
	now := TimeServer().Unix()

	if now >= self.Sql_CrossArena3V3.StartTime && now <= self.Sql_CrossArena3V3.EndTime {
		return false
	}
	return true
}

func (self *ModCrossArena3V3) GetRobotTeam3V3(robot *Js_CrossArena3V3User) [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo {
	data := [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo{}

	//生成机器人阵容
	index := 0
	for _, config := range GetCsvMgr().JJCRobotConfig {
		if config.Type != 4 {
			continue
		}
		if config.Jjcclass == robot.Subsection && config.Jjcdan == robot.Class {
			data[index] = GetCsvMgr().GetRobot(config)
			data[index].Uid = robot.Uid
			data[index].Iconid = robot.Icon
			if data[index] != nil {
				data[index].Uname = robot.UName

				isSet := false
				for k, _ := range data[index].Heroinfo {
					if !isSet {
						data[index].Heroinfo[k].Fight = data[index].Deffight
						isSet = true
					} else {
						data[index].Heroinfo[k].Fight = 0
					}
				}
			}
			index++
		}
	}
	return data
}
