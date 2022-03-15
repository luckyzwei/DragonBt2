package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

const (
	MSG_ARENA_SPECIAL_ENTER            = "arena_special_enter"
	MSG_ARENA_SPECIAL_GET_ENEMY        = "arena_special_get_enemy"
	MSG_ARENA_SPECIAL_START_FIGHT      = "arena_special_start_fight"
	MSG_ARENA_SPECIAL_GET_AWARD        = "arena_special_get_award"
	MSG_ARENA_SPECIAL_GET_FIGHTS       = "arena_special_get_fights"
	MSG_ARENA_SPECIAL_BUY_COUNT        = "arena_special_buy_count"
	MSG_ARENA_SPECIAL_GET_FIGHT_INFO   = "arena_special_get_fight_info"
	MSG_ARENA_SPECIAL_FIGHT_START      = "arena_special_fight_start"
	MSG_ARENA_SPECIAL_FIGHT_RESULT     = "arena_special_fight_result"
	MSG_ARENA_SPECIAL_ADD_FIGHT_RECORD = "arena_special_add_fight_record"
)

// 斗技场
type ModArenaSpecial struct {
	player *Player // player指针
	Index  int
}

// 竞技场战报
type ArenaSpecialFight struct {
	FightId  [ARENA_SPECIAL_TEAM_MAX]int64 `json:"fight_id"`      // 战斗Id
	Side     int                           `json:"side"`          // 1 进攻方 0 防守方
	Result   int                           `json:"attack_result"` // 0 进攻方成功 其他防守方胜利
	Class    int                           `json:"class"`         // 阶
	Dan      int                           `json:"dan"`           // 段位
	Uid      int64                         `json:"uid"`           // Uid
	IconId   int                           `json:"icon_id"`       // 头像Id
	Portrait int                           `json:"portrait"`      // 头像框
	Name     string                        `json:"name"`          // 名字
	Level    int                           `json:"level"`         // 等级
	Time     int64                         `json:"time"`          // 发生的时间
	Fight    int64                         `json:"fight"`         // 战力
}

func (self *ModArenaSpecial) OnGetData(player *Player) {
	self.player = player
	self.Index = -1
}

func (self *ModArenaSpecial) OnGetOtherData() {
	data := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
	if data != nil {
		data.State = 0

		if data.redPoint.IsFight == 1 {
			var msg S2C_ArenaSpecialAddFightRecord
			msg.Cid = MSG_ARENA_SPECIAL_ADD_FIGHT_RECORD
			self.player.Send(msg.Cid, msg)
			data.redPoint.IsFight = 0
		}
	}
	GetArenaSpecialMgr().UpdateFormat(self.player)
}

func (self *ModArenaSpecial) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModArenaSpecial) OnSave(sql bool) {
}

// 注册消息
func (self *ModArenaSpecial) onReg(handlers map[string]func(body []byte)) {
	handlers[MSG_ARENA_SPECIAL_ENTER] = self.ArenaSpecialEnter
	handlers[MSG_ARENA_SPECIAL_GET_ENEMY] = self.ArenaSpecialGetEnemy
	handlers[MSG_ARENA_SPECIAL_START_FIGHT] = self.ArenaSpecialStartFight
	handlers[MSG_ARENA_SPECIAL_GET_AWARD] = self.ArenaSpecialGetAward
	handlers[MSG_ARENA_SPECIAL_GET_FIGHTS] = self.ArenaSpecialGetFights
	handlers[MSG_ARENA_SPECIAL_BUY_COUNT] = self.ArenaSpecialBuyCount
	handlers[MSG_ARENA_SPECIAL_GET_FIGHT_INFO] = self.ArenaSpecialGetFightInfo
	handlers[MSG_ARENA_SPECIAL_FIGHT_RESULT] = self.ArenaSpecialFightResult
}

func (self *ModArenaSpecial) ArenaSpecialEnter(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_SPECIAL_ARENA)
	if !flag {
		//self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	var msg C2S_ArenaSpecialEnter
	json.Unmarshal(body, &msg)
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_DOES_NOT"))
		return
	}

	self.player.GetModule("team").(*ModTeam).CountArenaTopFight(TEAMTYPE_ARENA_SPECIAL_4)
	GetArenaSpecialMgr().UpdateFormat(self.player)

	uid := self.player.GetUid()
	data := GetArenaSpecialMgr().GetPlayerData(uid)
	if data == nil {
		// 如果玩家数据则添加数据
		GetArenaSpecialMgr().AddPlayer(self.player)
		self.CheckTeam()
		data = GetArenaSpecialMgr().GetPlayerData(uid)
	}
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}

	rank := GetArenaSpecialMgr().GetRankData(uid)
	if nil == rank {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}

	startTime, endTime := GetCsvMgr().GetNowStartAndEnd(ARENA_TIME_TYPE_SPECIAL)
	var backmsg S2C_ArenaSpecialEnter
	backmsg.Cid = MSG_ARENA_SPECIAL_ENTER
	backmsg.Class = data.Class
	backmsg.Dan = data.Dan
	backmsg.StartTime = startTime
	backmsg.EndTime = endTime
	backmsg.FreeCount = vipcsv.ArenaFree[ARENA_TYPE_SPECIAL]
	backmsg.FightCount = data.Count
	backmsg.BuyCount = data.BuyCount
	backmsg.Point = data.Point
	backmsg.Coin = data.Coin
	backmsg.ClassTime = rank.StartTime
	self.player.Send(backmsg.Cid, backmsg)
}
func (self *ModArenaSpecial) ArenaSpecialGetEnemy(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_SPECIAL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	var msg C2S_ArenaSpecialGetEnemy
	json.Unmarshal(body, &msg)

	data := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}

	var backmsg S2C_ArenaSpecialGetEnemy
	backmsg.Cid = MSG_ARENA_SPECIAL_GET_ENEMY
	enemy, ids := GetArenaSpecialMgr().GetEnemy(self.player.GetUid())
	backmsg.Enemy = enemy
	for _, id := range ids {
		config := GetCsvMgr().GetArenaSpecialClassConfigByID(id)
		if config != nil {
			backmsg.Class = append(backmsg.Class, config.Class)
			backmsg.Dan = append(backmsg.Dan, config.Dan)
		}
	}
	self.player.Send(backmsg.Cid, backmsg)
}
func (self *ModArenaSpecial) ArenaSpecialStartFight(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_SPECIAL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	var msg C2S_ArenaSpecialStartFight
	json.Unmarshal(body, &msg)

	data := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}

	// 找到敌人
	if len(data.enemy) <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_THE_ENEMY_DOES_NOT_EXIST"))
		return
	}

	if len(GetArenaSpecialMgr().FightList) >= ARENA_FIGHT_COUNT_MAX {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_THE_ENEMY_DOES_NOT_EXIST"))
		return
	}

	// 判断越界
	if msg.Index < 0 || msg.Index >= len(data.enemy) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_THE_ENEMY_DOES_NOT_EXIST"))
		return
	}
	// 判断次数
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_DOES_NOT"))
		return
	}
	// 数据错误
	if len(vipcsv.ArenaFree) < 0 {
		self.player.SendErr("len(vipcsv.JJcBuy) < 0")
		return
	}
	// 获得敌人
	fight := data.enemy[msg.Index]
	if fight.Class == 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_THE_OTHER_SIDE_IS_FIGHTING"))
		return
	}

	// 扣除物品
	//item := []PassItem{}
	maxnum := vipcsv.ArenaFree[ARENA_TYPE_SPECIAL]
	removeItem := []PassItem{}
	// 免费次数不足
	if data.Count >= maxnum {
		// 判断消耗是否足够
		cost := TARIFF_TYPE_ARENA_SPECIAL

		// 获得消耗配置
		config := GetCsvMgr().GetTariffConfig2(cost)
		if config == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}
		// 是否足够
		if err := self.player.HasObjectOk(config.ItemIds, config.ItemNums); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
		removeItem = self.player.RemoveObjectLst(config.ItemIds, config.ItemNums, "高阶竞技场", 0, 0, 0)
	}
	self.Index = msg.Index

	now := TimeServer().Unix()
	myFight := [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo{}
	myFight[ARENA_SPECIAL_TEAM_1] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_ARENA_SPECIAL_1)
	myFight[ARENA_SPECIAL_TEAM_2] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_ARENA_SPECIAL_2)
	myFight[ARENA_SPECIAL_TEAM_3] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_ARENA_SPECIAL_3)
	fightID := GetArenaSpecialMgr().AddFightList(self.player, myFight,
		fight.EnemyTeam,
		now,
		now)

	if fightID[0] > 0 {
		var backmsg S2C_ArenaSpecialStart
		backmsg.RandNum = now
		backmsg.Cid = MSG_ARENA_SPECIAL_FIGHT_START
		backmsg.MyFightInfo = myFight
		backmsg.FightID = fightID
		backmsg.Items = removeItem
		self.player.Send(backmsg.Cid, backmsg)
	}
}

func (self *ModArenaSpecial) ArenaSpecialGetFights(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_SPECIAL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	var msg C2S_ArenaSpecialGetFights
	json.Unmarshal(body, &msg)

	data := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}
	var backmsg S2C_ArenaSpecialGetFights
	backmsg.Cid = MSG_ARENA_SPECIAL_GET_FIGHTS
	backmsg.FightInfo = data.arenaFight
	self.player.Send(backmsg.Cid, backmsg)
}

func (self *ModArenaSpecial) ArenaSpecialGetAward(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_SPECIAL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	rank := GetArenaSpecialMgr().GetRankData(self.player.GetUid())
	if nil == rank {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}
	var msg C2S_ArenaSpecialGetAward
	json.Unmarshal(body, &msg)

	var backmsg S2C_ArenaSpecialGetAward
	backmsg.Cid = MSG_ARENA_SPECIAL_GET_AWARD
	backmsg.Items, backmsg.IsFull = GetArenaSpecialMgr().GetArenaAward(self.player)
	backmsg.Point = rank.Point
	backmsg.ClassTime = rank.StartTime
	self.player.Send(backmsg.Cid, backmsg)
}

// 获取战报信息
func (self *ModArenaSpecial) GetBattleInfo(fightID int64) []*BattleInfo {
	infos := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
	if infos == nil {
		return nil
	}

	var fight *ArenaSpecialFight
	for _, v := range infos.arenaFight {
		if fightID == v.FightId[0] {
			fight = v
			break
		}
	}

	if nil == fight {
		return nil
	}

	var ret []*BattleInfo
	for _, id := range fight.FightId {
		var battleInfo BattleInfo
		value, flag, err := HGetRedisEx(`san_arenaspecialbattleinfo`, id, fmt.Sprintf("%d", id))
		if err != nil || !flag {
			res := GetServer().DBUser.GetBattleInfo(id)
			if res != nil {
				ret = append(ret, res)
			}
			continue
		}
		if flag {
			err := json.Unmarshal([]byte(value), &battleInfo)
			if err != nil {
				continue
			}
		}

		if battleInfo.Id == 0 {
			continue
		}
		ret = append(ret, &battleInfo)
	}

	return ret
}

// 获取战报信息, 需要进行修改
func (self *ModArenaSpecial) GetBattleRecord(fightID int64) *BattleRecord {
	var battleRecord BattleRecord
	value, flag, err := HGetRedisEx(`san_arenaspecialbattlerecord`, fightID, fmt.Sprintf("%d", fightID))
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

// 购买次数
func (self *ModArenaSpecial) ArenaSpecialBuyCount(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_SPECIAL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	var msg C2S_ArenaSpecialBuyCount
	json.Unmarshal(body, &msg)

	count := 1
	if count <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	info := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 获得消耗配置
	config := GetCsvMgr().GetTariffConfig(TARIFF_TYPE_ARENA_SPECIALBUY, info.BuyCount+1)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	var costnum []int
	for _, v := range config.ItemNums {
		costnum = append(costnum, v*count)
	}
	// 是否足够
	if err := self.player.HasObjectOk(config.ItemIds, costnum); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	var getNum []int
	for _, v := range config.GetNum {
		getNum = append(getNum, v*count)
	}

	removeitems := self.player.RemoveObjectLst(config.ItemIds, config.ItemNums, "高阶竞技场购买", 0, 0, 0)
	additems := self.player.AddObjectLst(config.GetItem, getNum, "高阶竞技场购买", 0, 0, 0)

	info.BuyCount++
	var backmsg S2C_ArenaSpecialBuyCount
	backmsg.Cid = MSG_ARENA_SPECIAL_BUY_COUNT
	backmsg.Items = removeitems
	for _, v := range additems {
		backmsg.Items = append(backmsg.Items, v)
	}
	backmsg.BuyCount = info.BuyCount
	self.player.Send(backmsg.Cid, backmsg)
}

func (self *ModArenaSpecial) CheckTeam() {
	herouse1 := []int{}
	teamType1 := []int{TEAMTYPE_ARENA_SPECIAL_1, TEAMTYPE_ARENA_SPECIAL_2, TEAMTYPE_ARENA_SPECIAL_3}
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
		msg.TeamType = TEAMTYPE_ARENA_SPECIAL_1
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
		self.player.GetModule("team").(*ModTeam).AddArenaSpcialTeamUIPos(smsg)

	}
	herouse2 := []int{}
	teamType2 := []int{TEAMTYPE_ARENA_SPECIAL_4, TEAMTYPE_ARENA_SPECIAL_5, TEAMTYPE_ARENA_SPECIAL_6}
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
		msg.TeamType = TEAMTYPE_ARENA_SPECIAL_4
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
		self.player.GetModule("team").(*ModTeam).AddArenaSpcialTeamUIPos(smsg)
	}
}

// 获得战斗配置
func (self *ModArenaSpecial) ArenaSpecialGetFightInfo(body []byte) {
	var msg C2S_ArenaSpecialGetFightInfo
	json.Unmarshal(body, &msg)

	nType := msg.Type
	uid := msg.FindUid
	backmsg := &S2C_ArenaSpecialGetFightInfo{}
	backmsg.Cid = MSG_ARENA_SPECIAL_GET_FIGHT_INFO
	backmsg.Type = msg.Type
	backmsg.LifeTree = new(JS_LifeTreeInfo)
	if nType == GET_ENEMY_FIGHT_INFO_TYPE_RECORD {
		myUid := self.player.GetUid()
		fightInfos := GetArenaSpecialMgr().GetPlayerData(myUid)
		if fightInfos == nil {
			return
		}
		// 默认发防守方
		index := 1
		var oldRecord BattleRecord
		value, flag, err := HGetRedisEx(`san_arenaspecialbattlerecord`, uid, fmt.Sprintf("%d", uid))
		if err != nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
			return
		}
		if flag {
			err := json.Unmarshal([]byte(value), &oldRecord)
			if err != nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
				return
			}
		}
		if oldRecord.Id == 0 {
			self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
			return
		}
		// 如果自己是防守方 那就发进攻方
		if oldRecord.FightInfo[1].Uid == myUid {
			index = 0
		}

		class := 0
		dan := 0
		if oldRecord.FightInfo[index].Uid == 0 {
			find := false
			var fightid [ARENA_SPECIAL_TEAM_MAX]int64
			for _, v := range fightInfos.arenaFight {
				if v.FightId[0] == uid {
					find = true
					fightid = v.FightId
					class = v.Class
					dan = v.Dan
					break
				}
			}
			if !find {
				return
			}

			for i, v := range fightid {
				var record BattleRecord
				value, flag, err := HGetRedisEx(`san_arenaspecialbattlerecord`, v, fmt.Sprintf("%d", v))
				if err != nil {
					continue
				}
				if flag {
					err := json.Unmarshal([]byte(value), &record)
					if err != nil {
						continue
					}
				}
				if record.Id == 0 {
					continue
				}

				backmsg.FightInfo[i] = record.FightInfo[index]
				backmsg.Class = class
				backmsg.Dan = dan
			}
		} else {
			enemyInfo := GetArenaSpecialMgr().GetPlayerData(oldRecord.FightInfo[index].Uid)
			if enemyInfo == nil {
				return
			}
			backmsg.FightInfo = enemyInfo.format
			backmsg.Class = enemyInfo.Class
			backmsg.Dan = enemyInfo.Dan

			data := GetMasterMgr().GetPlayer(oldRecord.FightInfo[index].Uid)
			if data != nil && data.Data != nil {
				backmsg.LifeTree = data.LifeTree
			}
		}

	} else if nType == GET_ENEMY_FIGHT_INFO_TYPE_ENEMY {
		index := uid
		fightInfos := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
		if fightInfos == nil {
			return
		}
		if index < 0 || index >= int64(len(fightInfos.enemy)) {
			return
		}

		backmsg.FightInfo = fightInfos.enemy[index].EnemyTeam
		backmsg.Class = fightInfos.enemy[index].Class
		backmsg.Dan = fightInfos.enemy[index].Dan

		if fightInfos.enemy[index].EnemyTeam[0].Uid != 0 {
			data := GetMasterMgr().GetPlayer(fightInfos.enemy[index].EnemyTeam[0].Uid)
			if data != nil && data.Data != nil {
				backmsg.LifeTree = data.LifeTree
			}
		}
	} else {
		if uid != 0 {
			fightInfos := GetArenaSpecialMgr().GetPlayerData(uid)
			if fightInfos == nil {
				return
			}
			backmsg.FightInfo = fightInfos.format
			backmsg.Class = fightInfos.Class
			backmsg.Dan = fightInfos.Dan
			data := GetMasterMgr().GetPlayer(uid)
			if data != nil && data.Data != nil {
				backmsg.LifeTree = data.LifeTree
			}
		} else {
			config := GetCsvMgr().GetArenaSpecialClassConfigByID(msg.Rank)
			if config != nil {
				fightInfos := GetArenaSpecialMgr().GetRobot(config.Class, config.Dan)
				backmsg.FightInfo = fightInfos
				backmsg.Class = config.Class
				backmsg.Dan = config.Dan
			}
		}
	}

	self.player.Send(backmsg.Cid, backmsg)
}

func (self *ModArenaSpecial) UpdateFormat(heroUid int) {
	herouse := []int{}
	team := -1
	index := -1
	var msg C2S_AddTeamUIPos
	teamType := []int{TEAMTYPE_ARENA_SPECIAL_4, TEAMTYPE_ARENA_SPECIAL_5, TEAMTYPE_ARENA_SPECIAL_6}
	for _, g := range teamType {
		teamPos := self.player.getTeamPosByType(g)
		if nil == teamPos || teamPos.isUIPosEmpty() {
			self.CheckTeam()
			return
		}

		count := 0
		for i, v := range teamPos.FightPos {
			if v == heroUid {
				index = i
				team = g

				for _, t := range teamPos.FightPos {
					msg.FightPos = append(msg.FightPos, t)
				}
			}

			hero := self.player.GetModule("hero").(*ModHero).GetHero(v)
			if nil == hero {
				continue
			}
			herouse = append(herouse, hero.HeroId)
			count++
		}

		if count <= 1 && index >= 0 && team >= 0 {
			teamPos.FightPos[index] = 0
			self.CheckTeam()
			return
		}
	}

	//if index < 0 || team < 0 {
	//	return
	//}

	//// 默认上阵最强阵容 且不能有重复英雄
	//heros := self.player.GetModule("hero").(*ModHero).GetBestFormat2()
	//for _, v := range heros {
	//	// 找不到英雄
	//	hero := self.player.GetModule("hero").(*ModHero).GetHero(v)
	//	if nil == hero {
	//		continue
	//	}
	//	// 该类英雄已经使用
	//	find := false
	//	for _, t := range herouse {
	//		if hero.HeroId == t {
	//			find = true
	//		}
	//	}
	//	if find {
	//		continue
	//	}
	//	msg.FightPos[index] = v
	//	break
	//}
	//
	//msg.TeamType = team
	//smsg, _ := json.Marshal(&msg)
	//self.player.GetModule("team").(*ModTeam).addUIPos(smsg)
}

// 获得战斗结果
func (self *ModArenaSpecial) ArenaSpecialFightResult(body []byte) {
	var msg C2S_ArenaSpecialFightResult
	json.Unmarshal(body, &msg)

	GetArenaSpecialMgr().ArenaSpecialFightResult(msg.BattleInfo)
}

func (self *ModArenaSpecial) CheckTask() {
	data := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
	if data != nil {
		class := data.Class
		dan := data.Dan
		self.player.HandleTask(TASK_TYPE_SPECIAL_ARENA_CLASS, class, dan, 0)
	}
}
