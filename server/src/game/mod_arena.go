package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

//
//// 斗技场对手
//type ArenaEnemy struct {
//	Enemy []*JS_FightInfo
//	Index int
//}

const (
	GET_ENEMY_FIGHT_INFO_TYPE_RECORD = 1 // 历史记录
	GET_ENEMY_FIGHT_INFO_TYPE_RANK   = 2 // 排行
	GET_ENEMY_FIGHT_INFO_TYPE_ENEMY  = 3 // 敌人
)

type ArenaEnemy struct {
	Type  int   // 类型
	Index int64 // id
}

// 斗技场
type ModArena struct {
	player *Player // player指针

	enemy      *ArenaEnemy // 敌人index
	updateTime int64       // 战力更新排行榜cd
}

func (self *ModArena) OnGetData(player *Player) {
	self.player = player
}

func (self *ModArena) OnGetOtherData() {

	self.enemy = &ArenaEnemy{0, -1}

	GetArenaMgr().UpdateFormat(self.player)
	data := GetArenaMgr().GetPlayerArenaData(self.player.GetUid())
	if data != nil {
		if data.redpoint.IsFight == 1 {
			var msg S2C_ArenaAddFightRecord
			msg.Cid = "arena_add_fight_record"
			self.player.Send(msg.Cid, msg)
			data.redpoint.IsFight = 0
		}
	}
}

func (self *ModArena) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "armsarenainfo": // 获取玩家斗技场信息
		var msg C2S_ArenaInfo
		json.Unmarshal(body, &msg)
		self.GetBaseInfo(true)
		return true
	case "getarmsarenafight": // 得到对手信息
		var msg C2S_GetArenaFight
		json.Unmarshal(body, &msg)
		self.GetFightInfo()
		return true
	case "armsarenafight": // 战斗
		var msg C2S_ArenaBegin
		json.Unmarshal(body, &msg)
		self.Fight(msg.Index)
		return true
	case "armsarenafightback": // 战斗
		var msg C2S_ArenaFightBack
		json.Unmarshal(body, &msg)
		self.FightBack(msg.FightID)
		return true
	case "get_pvp_fights": // 获得玩家全部战报
		var msg C2S_GetArenaFightInfo
		json.Unmarshal(body, &msg)
		self.GetFights()
		return true
	case "get_enemy_fight_info": // 获得玩家防守阵容
		var msg C2S_GetEnemyFightInfo
		json.Unmarshal(body, &msg)
		self.GetEnemyFightInfo(msg.Type, msg.FindUid)
		return true
	case "buy_arena_count": // 购买竞技场次数
		var msg C2S_BuyArenaCount
		json.Unmarshal(body, &msg)
		self.BuyArenaCount(msg.Count)
		return true
	case "arena_fight_result": // 客户端发来结果
		var msg C2S_ArenaFightResult
		json.Unmarshal(body, &msg)
		self.ArenaFightResult(msg.Type, msg.BattleInfo)
		return true
	}

	return false
}

func (self *ModArena) OnSave(sql bool) {
}

// 获取玩家斗技场信息
func (self *ModArena) GetBaseInfo(add bool) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	// 检查更新
	self.player.CheckRefresh()

	self.player.GetModule("team").(*ModTeam).CountArenaTopFight(TEAMTYPE_ARENA_2)
	GetArenaMgr().UpdateFormat(self.player)
	// 玩家uid
	playerUid := self.player.GetUid()
	// 如果玩家数据则添加数据
	var msg S2C_ArenaInfo
	msg.Cid = "arenainfo"
	GetArenaMgr().AddPlayer(self.player)
	self.CheckTeam()

	info := GetArenaMgr().GetPlayerArenaData(playerUid)
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}

	now := TimeServer().Unix()
	if self.updateTime == 0 || now-self.updateTime > HOUR_SECS {
		GetTopArenaMgr().UpdateFight(ARENA_TOP_TYPE_NOMAL, self.player)
		self.updateTime = now
	}

	startTime, endTime := GetCsvMgr().GetNowStartAndEnd(ARENA_TIME_TYPE_NOMAL)
	msg.Rank = info.Rank
	msg.StartTime = startTime
	msg.EndTime = endTime

	self.player.SendMsg("armsarenainfo", HF_JtoB(&msg))
}

// 得到对手信息
func (self *ModArena) GetFightInfo() {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	// 判断次数
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_DOES_NOT"))
		return
	}

	GetArenaMgr().AddPlayer(self.player)
	self.CheckTeam()
	info := GetArenaMgr().GetPlayerArenaData(self.player.GetUid())
	if nil == info {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_DOES_NOT"))
		return
	}

	var msg S2C_GetArenaFight
	msg.Cid = "getarmsarenafight"
	msg.Enemy, msg.Point = self.GetEnemy()
	msg.FreeCount = vipcsv.ArenaFree[ARENA_TYPE_NOMAL]
	msg.FightCount = info.Count
	//msg.Top = fightbase
	self.player.SendMsg("getarmsarenafight", HF_JtoB(&msg))
}

// 获取对手信息
func (self *ModArena) GetEnemy() ([]*JS_FightInfo, []int64) {
	fightinfo := make([]*JS_FightInfo, 0)
	point := make([]int64, 0)
	fightinfo, point = GetArenaMgr().GetEnemy(self.player.Sql_UserBase.Uid)

	return fightinfo, point
}

// 战斗
func (self *ModArena) Fight(index int) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	// 获得玩家数据
	info := GetArenaMgr().GetPlayerArenaData(self.player.Sql_UserBase.Uid)
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}
	// 找到敌人
	if len(info.enemy) <= 0 {
		LogError("敌人不存在:", index, ", ", len(info.enemy))
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_THE_ENEMY_DOES_NOT_EXIST"))
		self.GetFightInfo()
		return
	}

	if len(GetArenaMgr().FightList) >= ARENA_FIGHT_COUNT_MAX {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_THE_ENEMY_DOES_NOT_EXIST"))
		return
	}

	// 判断越界
	if index < 0 || index >= len(info.enemy) {
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
	fight := info.enemy[index]
	if fight == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_THE_OTHER_SIDE_IS_FIGHTING"))
		return
	}

	// 扣除物品
	//item := []PassItem{}
	maxnum := vipcsv.ArenaFree[ARENA_TYPE_NOMAL]
	// 免费次数不足
	if info.Count >= maxnum {
		// 判断消耗是否足够
		cost := TARIFF_TYPE_ARENA_NORMAL

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
		self.player.RemoveObjectLst(config.ItemIds, config.ItemNums, "竞技场", 0, 0, 0)
	}
	self.enemy.Type = ARENA_FIGHT_TYPE_ENEMY
	self.enemy.Index = int64(index)
	info.RandNum = TimeServer().Unix()

	myInfo := GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_ARENA_1)
	fightID := GetArenaMgr().AddFightList(ARENA_FIGHT_TYPE_ENEMY,
		myInfo,
		fight.Enemy,
		info.RandNum,
		info.RandNum)

	if fightID > 0 {
		msg := &S2C_ArenaStart{}
		msg.Cid = "arena_fight_start"
		msg.RandNum = info.RandNum
		msg.FightID = fightID
		msg.MyFightInfo = myInfo
		self.player.Send(msg.Cid, msg)
	}
}

// 反击
func (self *ModArena) FightBack(fightID int64) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	// 获得玩家数据
	info := GetArenaMgr().GetPlayerArenaData(self.player.Sql_UserBase.Uid)
	if info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}
	// 获得战报
	var fight *ArenaFight = nil

	nLen := len(info.arenaFight)
	for i := 0; i < nLen; i++ {
		if fightID == info.arenaFight[i].FightId {
			fight = info.arenaFight[i]
			break
		}
	}

	// 没有此战报
	if nil == fight {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}

	// 不是失败
	if (fight.Result == 0 && fight.Side == 1) || (fight.Result == 1 && fight.Side == 0) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("不是失败"))
		return
	}

	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_DOES_NOT"))
		return
	}

	if len(vipcsv.ArenaFree) < 0 {
		self.player.SendErr("len(vipcsv.JJcBuy) < 0")
		return
	}

	//item := []PassItem{}
	maxnum := vipcsv.ArenaFree[ARENA_TYPE_NOMAL]
	if info.Count >= maxnum {
		cost := TARIFF_TYPE_ARENA_NORMAL

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
		self.player.RemoveObjectLst(config.ItemIds, config.ItemNums, "竞技场", 0, 0, 0)
	}
	// 设置随机数
	info.RandNum = TimeServer().Unix()
	self.enemy.Type = ARENA_FIGHT_TYPE_FIGHT_BACK
	self.enemy.Index = fightID

	var enemyinfo *JS_FightInfo = nil
	if fight.Uid != 0 {
		// 获得敌人数据 如果是机器人 则直接返回
		enemy := GetArenaMgr().GetPlayerArenaData(fight.Uid)
		if enemy != nil {
			enemyinfo = enemy.format
		}
	} else {
		var oldRecord BattleRecord
		value, flag, err := HGetRedisEx(`san_arenabattlerecord`, fightID, fmt.Sprintf("%d", fightID))
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

		enemyinfo = oldRecord.FightInfo[1]
	}
	myInfo := GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_ARENA_1)
	ret := GetArenaMgr().AddFightList(ARENA_FIGHT_TYPE_FIGHT_BACK,
		myInfo,
		enemyinfo,
		info.RandNum,
		info.RandNum)

	if ret > 0 {
		msg := &S2C_ArenaBackStart{}
		msg.Cid = "arena_fight_back_start"
		msg.RandNum = info.RandNum
		msg.FightID = ret
		msg.MyFightInfo = myInfo
		self.player.Send(msg.Cid, msg)
	}
}

// 获取战报信息, 需要进行修改
func (self *ModArena) GetFights() {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	fightInfos := GetArenaMgr().GetPlayerArenaData(self.player.GetUid())
	if fightInfos == nil {
		return
	}
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_DOES_NOT"))
		return
	}
	msg := &S2C_ArenaFights{}
	msg.Cid = "get_pvp_fights"
	msg.FightInfo = fightInfos.arenaFight
	msg.FreeCount = vipcsv.ArenaFree[ARENA_TYPE_NOMAL]
	msg.FightCount = fightInfos.Count
	self.player.Send(msg.Cid, msg)
}

// 获取战报信息
func (self *ModArena) GetBattleInfo(fightID int64) *BattleInfo {
	var battleInfo BattleInfo
	value, flag, err := HGetRedisEx(`san_arenabattleinfo`, fightID, fmt.Sprintf("%d", fightID))
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

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ARENA_BATTLE_INFO, int(fightID), 0, 0, "查看竞技场战报", 0, 0, self.player)

	if battleInfo.Id != 0 {
		return &battleInfo
	}
	return nil
}

// 获取战报信息, 需要进行修改
func (self *ModArena) GetBattleRecord(fightID int64) *BattleRecord {
	var battleRecord BattleRecord
	value, flag, err := HGetRedisEx(`san_arenabattlerecord`, fightID, fmt.Sprintf("%d", fightID))
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

// 获得defend
func (self *ModArena) GetBattleDefend(uid int64) {
	fightInfos := GetArenaMgr().GetPlayerArenaData(uid)
	if fightInfos == nil {
		return
	}
	msg := &S2C_ArenaBattleDefend{}
	msg.Cid = "get_player_defend"
	msg.FightInfo = fightInfos.format
	self.player.Send(msg.Cid, msg)
}

// 获得defend
func (self *ModArena) GetEnemyFightInfo(nType int, uid int64) {
	msg := &S2C_GetEnemyFightInfo{}
	msg.Cid = "get_enemy_fight_info"
	msg.Type = nType
	if nType == GET_ENEMY_FIGHT_INFO_TYPE_RECORD {
		myUid := self.player.GetUid()
		fightInfos := GetArenaMgr().GetPlayerArenaData(myUid)
		if fightInfos == nil {
			return
		}

		// 默认发防守方
		index := 1

		var oldRecord BattleRecord
		value, flag, err := HGetRedisEx(`san_arenabattlerecord`, uid, fmt.Sprintf("%d", uid))
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

		if oldRecord.FightInfo[index].Uid == 0 {
			msg.FightInfo = oldRecord.FightInfo[index]
		} else {
			enemyInfo := GetArenaMgr().GetPlayerArenaData(oldRecord.FightInfo[index].Uid)
			if enemyInfo == nil {
				return
			}
			msg.FightInfo = enemyInfo.format
			data := GetMasterMgr().GetPlayer(oldRecord.FightInfo[index].Uid)
			if data != nil && data.Data != nil {
				msg.LifeTree = data.LifeTree
			}
		}
	} else if nType == GET_ENEMY_FIGHT_INFO_TYPE_ENEMY {
		index := uid
		fightInfos := GetArenaMgr().GetPlayerArenaData(self.player.GetUid())
		if fightInfos == nil {
			return
		}
		if index <= 0 || index > int64(len(fightInfos.enemy)) {
			return
		}

		msg.FightInfo = fightInfos.enemy[index-1].Enemy

		if fightInfos.enemy[index-1].Enemy.Uid > 0 {
			data := GetMasterMgr().GetPlayer(fightInfos.enemy[index-1].Enemy.Uid)
			if data != nil && data.Data != nil {
				msg.LifeTree = data.LifeTree
			}
		}
	} else {
		fightInfos := GetArenaMgr().GetPlayerArenaData(uid)
		if fightInfos == nil {
			return
		}
		msg.FightInfo = fightInfos.format
		data := GetMasterMgr().GetPlayer(uid)
		if data != nil && data.Data != nil {
			msg.LifeTree = data.LifeTree
		}
	}

	self.player.Send(msg.Cid, msg)
}

// 获得defend
func (self *ModArena) BuyArenaCount(count int) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_ARENA)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN")))
		return
	}
	if count <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 获得消耗配置
	config := GetCsvMgr().GetTariffConfig2(TARIFF_TYPE_ARENA_NORMALBUY)
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

	removeitems := self.player.RemoveObjectLst(config.ItemIds, costnum, "竞技场购买", 0, 0, 0)
	additems := self.player.AddObjectLst(config.GetItem, getNum, "竞技场购买", 0, 0, 0)
	msg := &S2C_BuyArenaCount{}
	msg.Cid = "buy_arena_count"
	msg.Items = removeitems

	for _, v := range additems {
		msg.Items = append(msg.Items, v)
	}
	self.player.Send(msg.Cid, msg)
}

func (self *ModArena) CheckTeam() {
	// 检查玩家身上的队伍 如果直接为空
	teamPos := self.player.getTeamPosByType(TEAMTYPE_ARENA_2)
	if nil == teamPos || teamPos.isUIPosEmpty() {
		herouse := []int{}
		var msg C2S_AddTeamUIPos
		msg.TeamType = TEAMTYPE_ARENA_2
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
	// 默认上阵最强阵容 且不能有重复英雄
	teamPos = self.player.getTeamPosByType(TEAMTYPE_ARENA_1)
	if nil == teamPos || teamPos.isUIPosEmpty() {
		herouse := []int{}
		var msg C2S_AddTeamUIPos
		msg.TeamType = TEAMTYPE_ARENA_1
		nCount := 0
		heros := self.player.GetModule("hero").(*ModHero).GetBestFormat2()
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

func (self *ModArena) UpdateFormat(heroUid int) {
	// 检查玩家身上的队伍 如果直接为空
	teamPos := self.player.getTeamPosByType(TEAMTYPE_ARENA_2)
	if nil == teamPos || teamPos.isUIPosEmpty() {
		self.CheckTeam()
		return
	}
	var msg C2S_AddTeamUIPos
	index := -1
	count := 0
	herouse := []int{}
	for i, v := range teamPos.FightPos {
		if v == 0 {
			continue
		}

		if v == heroUid {
			index = i
		}
		msg.FightPos = append(msg.FightPos, v)
		hero := self.player.GetModule("hero").(*ModHero).GetHero(v)
		if nil == hero {
			continue
		}
		herouse = append(herouse, hero.HeroId)
		count++
	}
	if index < 0 {
		return
	}

	if count <= 1 {
		teamPos.FightPos[index] = 0
		self.CheckTeam()
	}
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
	//msg.TeamType = TEAMTYPE_ARENA_2
	//smsg, _ := json.Marshal(&msg)
	//self.player.GetModule("team").(*ModTeam).addUIPos(smsg)
}

func (self *ModArena) ArenaFightResult(nType int, battleInfo *BattleInfo) {
	GetArenaMgr().ArenaFightResult(nType, battleInfo)
}
