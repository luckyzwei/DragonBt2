package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

const (
	TOWER_TYPE_0   = 0 // 塔总表
	TOWER_TYPE_1   = 1 // 塔1
	TOWER_TYPE_2   = 2 // 塔2
	TOWER_TYPE_3   = 3 // 塔3
	TOWER_TYPE_4   = 4 // 塔4
	TOWER_TYPE_MAX = 5 // 最大值
)
const (
	TOWER_MAIN_TYPE_0 = 30 // 塔总表
	TOWER_MAIN_TYPE_1 = 31 // 塔1
	TOWER_MAIN_TYPE_2 = 32 // 塔2
	TOWER_MAIN_TYPE_3 = 33 // 塔3
	TOWER_MAIN_TYPE_4 = 34 // 塔4
)

const TOWER_LOOK_FLOOR = 2
const TOWER_FIGHT_MAX = 10

type Js_TowerFightRecord struct {
	Key         int64  `json:"key"`      // key值
	Name        string `json:"name"`     // 玩家名字
	Uid         int64  `json:"uid"`      // 玩家uid
	Icon        int    `json:"icon"`     // 头像
	Portrait    int    `json:"portrait"` //
	Level       int    `json:"level"`
	PlayerFight int64  `json:"playerfight"` // 玩家战力
	BattleFight int64  `json:"battlefight"` // 战斗参与的战力
	Time        int64  `json:"time"`        // 时间
}

// 爬塔结构
type JS_Tower struct {
	Type       int   `json:"type"`       //塔类型
	MaxLevel   int   `json:"maxlevel"`   //历史最大关卡
	MaxLevelTs int64 `json:"maxlevelts"` //历史最大关卡获得时间
	CurLevel   int   `json:"curlevel"`   //当前关卡
	LevelCount int   `json:"levelcount"` //当天打了多少层
}

// 爬塔相关逻辑
type San_Tower struct {
	Uid  int64 `json:"uid"`
	Info string

	info []*JS_Tower
	DataUpdate
}

func (self *San_Tower) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Info), &self.info)

}

func (self *San_Tower) Encode() { //! 将data数据写入数据库
	self.Info = HF_JtoA(&self.info)
}

// 爬塔数据
type ModTower struct {
	player *Player
	data   San_Tower
}

func (self *ModTower) GetData() *San_Tower {
	return &self.data
}

// 获取玩家数据
func (self *ModTower) OnGetData(player *Player) {
	self.player = player
}

// 获取玩家数据
func (self *ModTower) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_kingtower` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.data, "san_kingtower", self.player.ID)
	if self.data.Uid <= 0 {
		self.data.Uid = self.player.ID

		for i := 0; i < TOWER_TYPE_MAX; i++ {
			temp := &JS_Tower{}
			temp.Type = i
			self.data.info = append(self.data.info, temp)
		}

		self.data.Encode()
		InsertTable("san_kingtower", &self.data, 0, true)
	} else {
		self.data.Decode()
	}
	self.data.Init("san_kingtower", &self.data, true)

	nLen := len(self.data.info)
	if nLen < TOWER_TYPE_MAX {
		for i := 0; i < TOWER_TYPE_MAX-nLen; i++ {
			temp := JS_Tower{}
			temp.Type = nLen + i
			self.data.info = append(self.data.info, &temp)
		}
	}
}

func (self *ModTower) OnRefresh() {
	for _, v := range self.data.info {
		v.LevelCount = 0
	}
}

// 消息处理
func (self *ModTower) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "towerinfo":
		// 获取爬塔数据
		self.sendInfo()
		return true
	case "towerfightskip":
		// 发送爬塔战斗结果
		var msg C2S_TowerFightSkip
		json.Unmarshal(body, &msg)
		self.FightSkip(msg.LevelId)
		return true
	case "towerfightbegin":
		// 发送爬塔战斗结果
		var msg C2S_TowerFightBegin
		json.Unmarshal(body, &msg)
		self.FightBegin(msg.LevelId)
		return true
	case "towerfightresult":
		// 发送爬塔战斗结果
		var msg C2S_TowerFightResult
		json.Unmarshal(body, &msg)
		self.FightResult(msg.LevelId, msg.Result, msg.BattleInfo)
		return true
	case "towerrank":
		// 获取爬塔排行榜
		var msg C2S_TowerRank
		json.Unmarshal(body, &msg)
		self.GetRank(msg.Type)
		return true
	case "towerfloorinfo":
		// 获得层数信息
		var msg C2S_TowerFloorInfo
		json.Unmarshal(body, &msg)
		self.GetFloor(msg.LevelId)
		return true
	}
	return false
}

// 数据存盘
func (self *ModTower) OnSave(sql bool) {
	self.data.Encode()
	self.data.Update(sql)
}

// 同步爬塔消息
func (self *ModTower) sendInfo() {
	var msg S2C_TowerInfo
	msg.Cid = "towerinfo"
	msg.Data = self.data.info

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModTower) GetMainLevel() int {
	return self.data.info[0].CurLevel
}

func (self *ModTower) GetRank(nType int) {
	player := self.player
	var msg S2C_TowerRank
	msg.Cid = "towerrank"

	rank, _ := GetTopTowerMgr().GetTop(nType)
	for i := 0; i < len(rank); i++ {
		msg.Items = append(msg.Items, rank[i])
	}

	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModTower) GetBattleInfo(id int64) *BattleInfo {
	var battleInfo *BattleInfo
	//value, flag, err := HGetRedis(`san_towerbattleinfo`, fmt.Sprintf("%d", id))
	//if err != nil {
	//	self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
	//	return &battleInfo
	//}
	//if flag {
	//	err := json.Unmarshal([]byte(value), &battleInfo)
	//	if err != nil {
	//		self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
	//		return &battleInfo
	//	}
	//}

	battleInfo = GetTowerMgr().GetBattleInfo(id)

	if battleInfo != nil && battleInfo.Id != 0 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_TOWER_BATTLE_INFO, int(id/1000000), int(id%1000000), 0, "查看试炼之塔战报", 0, 0, self.player)
		return battleInfo
	}
	return nil
}

func (self *ModTower) GetBattleRecord(id int64) *BattleRecord {
	var battleRecord *BattleRecord
	//value, flag, err := HGetRedis(`san_towerbattlerecord`, fmt.Sprintf("%d", id))
	//if err != nil {
	//	self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
	//	return &battleRecord
	//}
	//if flag {
	//	err := json.Unmarshal([]byte(value), &battleRecord)
	//	if err != nil {
	//		self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
	//		return &battleRecord
	//	}
	//}

	battleRecord = GetTowerMgr().GetBattleRecord(id)
	if battleRecord != nil && battleRecord.Id != 0 {
		return battleRecord
	}
	return nil
}

func (self *ModTower) GetFloor(levelId int) {
	player := self.player

	// 关卡信息
	levelConfig := GetTowerMgr().GetLevelConfig(levelId)
	if levelConfig == nil {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_TOWER_LEVEL_ERROR"))
		return
	}

	nType := self.GetTowerType(levelConfig.MainType)
	// 检查当前关卡等级信息
	if levelConfig.ChapterIndex-TOWER_LOOK_FLOOR-1 > self.data.info[nType].CurLevel {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TOWER_CURRENT_LAYER_ERROR"))
		return
	}

	var msg S2C_TowerFloorInfo
	msg.Cid = "towerfloorinfo"
	msg.LevelId = levelId
	msg.Record = append(msg.Record, GetTowerMgr().GetRecord(int64(levelId)))

	friendkeys := make(map[int64]int64)
	friends := self.player.GetModule("friend").(*ModFriend).getFriend()
	for _, v := range friends {
		key := HF_AtoI64(HF_I64toA(v.Uid) + HF_ItoA(levelId))
		friendkeys[key] = key
	}
	msg.Record = append(msg.Record, GetTowerMgr().GetRecordList(friendkeys))

	unionkeys := make(map[int64]int64)
	unionid := self.player.GetUnionId()
	if unionid > 0 {
		data := GetUnionMgr().GetUnion(unionid)
		if data != nil {
			for _, v := range data.member {
				key := HF_AtoI64(HF_I64toA(v.Uid) + HF_ItoA(levelId))
				unionkeys[key] = key
			}
			msg.Record = append(msg.Record, GetTowerMgr().GetRecordList(unionkeys))
		}
	} else {
		msg.Record = append(msg.Record, nil)
	}

	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

// 爬塔战斗结束,310006之后首次通关有策划bug
func (self *ModTower) FightResult(levelid int, result int, battleInfo *BattleInfo) {
	player := self.player

	// 关卡信息
	levelConfig := GetTowerMgr().GetLevelConfig(levelid)
	if levelConfig == nil {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_TOWER_LEVEL_ERROR"))
		return
	}

	nType := self.GetTowerType(levelConfig.MainType)
	//teamType := self.GetTowerTeamPos(nType)
	//if levelConfig.Comat > self.player.GetTeamFight(teamType) {
	//	player.SendErrInfo("err", "战力不足")
	//	return
	//}
	// 检查当前关卡等级信息
	if levelConfig.ChapterIndex != self.data.info[nType].CurLevel+1 {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TOWER_CURRENT_LAYER_ERROR"))
		return
	}

	if nType != TOWER_TYPE_0 {
		if self.data.info[nType].LevelCount >= TOWER_FIGHT_MAX {
			player.SendErrInfo("err", GetCsvMgr().GetText("超过当天最高次数"))
			return
		}
	}

	var msg S2C_TowerFightResult
	msg.Cid = "towerfightresult"
	if result != 0 {
		// 如果结果>0,说明胜利
		self.data.info[nType].CurLevel = self.data.info[nType].CurLevel + 1
		// 检查是否首次通关
		// 首次通关
		if self.data.info[nType].CurLevel > self.data.info[nType].MaxLevel {
			self.data.info[nType].MaxLevel = self.data.info[nType].CurLevel
			self.data.info[nType].MaxLevelTs = TimeServer().Unix() // 达成的时间
		}
		outitem := self.player.GetModule("pass").(*ModPass).GetDropItem(levelid, true)
		for j := 0; j < len(outitem); j++ {
			outitem[j].ItemID, outitem[j].Num = self.player.AddObject(outitem[j].ItemID, outitem[j].Num, levelid, 0, 0, "试炼之塔通关奖励")
		}
		self.data.info[nType].LevelCount++
		msg.Items = outitem
		player.HandleTask(TowerWinTask, levelid, 1, 0)
		GetTopTowerMgr().updateRank(nType, int64(self.data.info[nType].MaxLevel), player)

		self.AddFightRecord(levelid, battleInfo)

		fight := int(self.player.GetTeamFight(self.GetTowerTeamPos(nType)))
		if nType == TOWER_TYPE_0 {
			self.player.HandleTask(TASK_TYPE_WOTER_LEVEL, self.data.info[nType].CurLevel, 0, 0)
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_TOWER_FINISH_NORMAL, levelid, fight, 0, "试炼之塔普通塔战斗", 0, 0, self.player)
		} else {
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_TOWER_FINISH_CAMP, levelid, fight, self.data.info[nType].LevelCount, "试炼之塔种族塔战斗", 0, 0, self.player)
		}
	}

	self.player.HandleTask(TASK_TYPE_WOTER_COUNT, 1, 0, 0)
	if nType != TOWER_TYPE_0 {
		self.player.HandleTask(TASK_TYPE_CAMP_WOTER_LEVEL, 0, 0, 0)
		self.player.HandleTask(TASK_TYPE_ONE_CAMP_TOWER_LEVEL, self.data.info[nType].CurLevel, nType, 0)
	}

	self.player.GetModule("task").(*ModTask).SendUpdate()
	msg.CurLevel = self.data.info[nType].CurLevel
	msg.MaxLevel = self.data.info[nType].MaxLevel
	msg.LevelCount = self.data.info[nType].LevelCount
	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModTower) MaxLevel() int {
	return 0
	//return self.data.MaxLevel
}

func (self *ModTower) FightBegin(levelid int) {
	player := self.player

	// 关卡信息
	levelConfig := GetTowerMgr().GetLevelConfig(levelid)
	if levelConfig == nil {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_TOWER_LEVEL_ERROR"))
		return
	}

	nType := self.GetTowerType(levelConfig.MainType)
	teamType := self.GetTowerTeamPos(nType)
	if levelConfig.Comat > self.player.GetTeamFight(teamType) {
		player.SendErrInfo("err", "战力不足")
		return
	}

	// 检查当前关卡等级信息
	if levelConfig.ChapterIndex != self.data.info[nType].CurLevel+1 {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TOWER_CURRENT_LAYER_ERROR"))
		return
	}

	self.player.HandleTask(TASK_TYPE_WOTER_COUNT, 0, 0, 0)
	if nType != TOWER_TYPE_0 {
		if self.data.info[nType].LevelCount >= TOWER_FIGHT_MAX {
			player.SendErrInfo("err", GetCsvMgr().GetText("超过当天最高次数"))
			return
		}
	}

	var msg S2C_TowerFightBegin
	msg.Cid = "towerfightbegin"
	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModTower) FightSkip(levelid int) {
	player := self.player

	// 关卡信息
	levelConfig := GetTowerMgr().GetLevelConfig(levelid)
	if levelConfig == nil {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_TOWER_LEVEL_ERROR"))
		return
	}

	nType := self.GetTowerType(levelConfig.MainType)
	teamType := self.GetTowerTeamPos(nType)
	if levelConfig.LevelSkip == LOGIC_FALSE || levelConfig.SkipType == LOGIC_FALSE {
		self.player.SendErrInfo("err", "此关卡不能跳过")
		return
	}

	if levelConfig.SkipNum > self.player.GetTeamFight(teamType) {
		self.player.SendErrInfo("err", "战力不足")
		return
	}

	// 检查当前关卡等级信息
	if levelConfig.ChapterIndex != self.data.info[nType].CurLevel+1 {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TOWER_CURRENT_LAYER_ERROR"))
		return
	}

	self.player.HandleTask(TASK_TYPE_WOTER_COUNT, 0, 0, 0)
	if nType != TOWER_TYPE_0 {
		if self.data.info[nType].LevelCount >= TOWER_FIGHT_MAX {
			player.SendErrInfo("err", GetCsvMgr().GetText("超过当天最高次数"))
			return
		}
	}

	var msg S2C_TowerFightSkip
	msg.Cid = "towerfightskip"
	self.data.info[nType].CurLevel = self.data.info[nType].CurLevel + 1
	// 检查是否首次通关
	// 首次通关
	if self.data.info[nType].CurLevel > self.data.info[nType].MaxLevel {
		self.data.info[nType].MaxLevel = self.data.info[nType].CurLevel
		self.data.info[nType].MaxLevelTs = TimeServer().Unix() // 达成的时间
	}
	outitem := self.player.GetModule("pass").(*ModPass).GetDropItem(levelid, true)
	for j := 0; j < len(outitem); j++ {
		outitem[j].ItemID, outitem[j].Num = self.player.AddObject(outitem[j].ItemID, outitem[j].Num, levelid, 0, 0, "试炼之塔通关奖励")
	}
	self.data.info[nType].LevelCount++
	msg.Items = outitem
	player.HandleTask(TowerWinTask, levelid, 1, 0)
	GetTopTowerMgr().updateRank(nType, int64(self.data.info[nType].MaxLevel), player)

	fight := int(self.player.GetTeamFight(self.GetTowerTeamPos(nType)))
	if nType == TOWER_TYPE_0 {
		self.player.HandleTask(TASK_TYPE_WOTER_LEVEL, self.data.info[nType].CurLevel, 0, 0)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_TOWER_FINISH_NORMAL, levelid, fight, 0, "试炼之塔普通塔战斗", 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_TOWER_FINISH_CAMP, levelid, fight, self.data.info[nType].LevelCount, "试炼之塔种族塔战斗", 0, 0, self.player)
	}

	self.player.HandleTask(TASK_TYPE_WOTER_COUNT, 1, 0, 0)
	if nType != TOWER_TYPE_0 {
		self.player.HandleTask(TASK_TYPE_CAMP_WOTER_LEVEL, 0, 0, 0)
		self.player.HandleTask(TASK_TYPE_ONE_CAMP_TOWER_LEVEL, self.data.info[nType].CurLevel, nType, 0)
	}

	self.player.GetModule("task").(*ModTask).SendUpdate()
	msg.CurLevel = self.data.info[nType].CurLevel
	msg.MaxLevel = self.data.info[nType].MaxLevel
	msg.LevelCount = self.data.info[nType].LevelCount
	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModTower) AddFightRecord(levelid int, battleInfo *BattleInfo) {
	// 关卡信息
	levelConfig := GetTowerMgr().GetLevelConfig(levelid)
	if levelConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_TOWER_LEVEL_ERROR"))
		return
	}

	nType := self.GetTowerType(levelConfig.MainType)
	teamType := self.GetTowerTeamPos(nType)
	team := self.player.getTeamPosByType(teamType)
	fight := int64(0)

	for i := 0; i < len(team.FightPos); i++ {
		heroKey := team.FightPos[i]
		if heroKey == 0 {
			continue
		}
		hero := self.player.getHero(heroKey)
		if nil == hero {
			continue
		}

		fight += hero.Fight
	}
	var army *ArmyInfo = nil
	for _, p := range battleInfo.UserInfo {
		for _, q := range p.HeroInfo {
			if q.ArmyInfo != nil && q.ArmyInfo.Uid != 0 {
				army = q.ArmyInfo
				army.Lv = q.HeroLv
				for i := 0; i < len(q.ArmyInfo.Atts); i++ {
					if q.ArmyInfo.Atts[i].Type == AttrFight {
						fight += q.ArmyInfo.Atts[i].Value
						break
					}
				}
				break
			}
		}
	}

	data1 := San_TowerPlayerRecord{}
	data1.KeyID = int64(levelid)
	data1.Uid = self.player.GetUid()
	data1.Name = self.player.GetName()
	data1.Icon = self.player.Sql_UserBase.IconId
	data1.Portrait = self.player.Sql_UserBase.Portrait
	data1.Level = self.player.Sql_UserBase.Level
	data1.BattleFight = fight
	data1.PlayerFight = self.player.Sql_UserBase.Fight
	data1.Time = TimeServer().Unix()

	data2 := data1
	data2.KeyID = HF_AtoI64(HF_I64toA(self.player.GetUid()) + HF_ItoA(levelid))

	battleInfo1 := BattleInfo(*battleInfo)
	battleInfo1.Id = int64(levelid)
	battleInfo1.LevelID = levelid
	battleInfo1.Type = BATTLE_TYPE_PVE
	battleInfo1.Time = TimeServer().Unix()
	battleInfo1.UserInfo[0].Uid = self.player.GetUid()
	battleInfo1.UserInfo[0].Level = self.player.Sql_UserBase.Level
	battleInfo1.UserInfo[0].Icon = self.player.Sql_UserBase.IconId
	battleInfo1.UserInfo[0].Portrait = self.player.Sql_UserBase.Portrait
	battleInfo1.UserInfo[0].UnionName = self.player.GetUnionName()
	battleInfo1.UserInfo[0].Name = self.player.GetName()
	battleInfo2 := BattleInfo(battleInfo1)
	battleInfo2.Id = data2.KeyID

	battleRecord1 := BattleRecord{}
	battleRecord1.Id = int64(levelid)
	battleRecord1.Time = TimeServer().Unix()
	battleRecord1.Result = battleInfo.Result
	battleRecord1.Type = BATTLE_TYPE_PVE
	battleRecord1.RandNum = battleInfo.Random
	if army != nil {
		battleRecord1.FightInfo[0] = GetRobotMgr().GetPlayerFightInfoWithArmyByPos(self.player, 0, 0, teamType, army)
		self.player.GetModule("friend").(*ModFriend).SetUseSign(self.GetTowerArmyType(nType), 1)
	} else {
		battleRecord1.FightInfo[0] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, teamType)
	}

	battleRecord1.Side = 1
	battleRecord1.Level = levelid
	battleRecord1.LevelID = levelid

	battleRecord2 := BattleRecord(battleRecord1)
	battleRecord2 = battleRecord1
	battleRecord2.Id = data2.KeyID

	GetTowerMgr().AddPlayerRecord(int64(levelid), &data1, &battleInfo1, &battleRecord1)
	GetTowerMgr().AddPlayerRecord(data2.KeyID, &data2, &battleInfo2, &battleRecord2)
}

func (self *ModTower) GetStatisticsValue1110() (lastValue int, buyValue int) {
	//vipConfig := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	//if vipConfig == nil {
	//	self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TOWER_ARISTOCRACY_DOES_NOT_EXIST"))
	//	return
	//}
	//return self.data.ResetTimes, vipConfig.Pata_buy1 - self.data.ResetBuyTimes

	return 0, 0
}

func (self *ModTower) GetStatisticsValue1111() (lastValue int, buyValue int) {
	//vipConfig := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	//if vipConfig == nil {
	//	self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TOWER_ARISTOCRACY_DOES_NOT_EXIST"))
	//	return 3 - self.data.AdvanceTimes, 0
	//}
	//return 3 - self.data.AdvanceTimes, vipConfig.Pata_buy2 - self.data.AdvanceBuyTimes
	return 0, 0
}

func (self *ModTower) GetTowerType(mainType int) int {
	nType := 0
	switch mainType {
	case TOWER_MAIN_TYPE_0:
		nType = TOWER_TYPE_0
	case TOWER_MAIN_TYPE_1:
		nType = TOWER_TYPE_1
	case TOWER_MAIN_TYPE_2:
		nType = TOWER_TYPE_2
	case TOWER_MAIN_TYPE_3:
		nType = TOWER_TYPE_3
	case TOWER_MAIN_TYPE_4:
		nType = TOWER_TYPE_4
	}
	return nType
}

func (self *ModTower) GetTowerArmyType(mainType int) int {
	nType := 0
	switch mainType {
	case TOWER_TYPE_0:
		nType = HIRE_MOD_TOWER
	case TOWER_TYPE_1:
		nType = HIRE_MOD_TOWER_CAMP_1
	case TOWER_TYPE_2:
		nType = HIRE_MOD_TOWER_CAMP_2
	case TOWER_TYPE_3:
		nType = HIRE_MOD_TOWER_CAMP_3
	case TOWER_TYPE_4:
		nType = HIRE_MOD_TOWER_CAMP_4
	}
	return nType
}

func (self *ModTower) GetTowerTeamPos(towerType int) int {
	nType := 0
	switch towerType {
	case TOWER_TYPE_0:
		nType = TEAMTYPE_TOWER_MAIN
	case TOWER_TYPE_1:
		nType = TEAMTYPE_TOWER_1
	case TOWER_TYPE_2:
		nType = TEAMTYPE_TOWER_2
	case TOWER_TYPE_3:
		nType = TEAMTYPE_TOWER_3
	case TOWER_TYPE_4:
		nType = TEAMTYPE_TOWER_4
	}
	return nType
}
func (self *ModTower) GmSetToweLevel(nType, level int) {
	for _, v := range self.data.info {
		if v.Type == nType {
			v.MaxLevel = level
			v.CurLevel = level
			v.MaxLevelTs = TimeServer().Unix()
			return
		}
	}
	return
}

func (self *ModTower) CheckTask() {
	for _, v := range self.data.info {
		if v.Type == TOWER_TYPE_0 {
			self.player.HandleTask(TASK_TYPE_WOTER_LEVEL, self.data.info[v.Type].CurLevel, 0, 0)
		} else {
			self.player.HandleTask(TASK_TYPE_CAMP_WOTER_LEVEL, 0, 0, 0)
			self.player.HandleTask(TASK_TYPE_ONE_CAMP_TOWER_LEVEL, self.data.info[v.Type].CurLevel, v.Type, 0)
		}
	}
}
