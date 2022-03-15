package match

/*
import "master/db"

import (
	"encoding/json"
	"fmt"
	"game"
	"master/core"
	"master/utils"
	"sync"
	"time"
)

const (
	FIGHTMATCH_TEAM_COUNT = 3  //
	BEST_SUBSECTION       = 1  //传奇段位
	RECORD_MAX_CROSS      = 10 //最大战报数s

	BATTLE_TYPE_PVP  = 1 //
	BATTLE_TYPE_PVE  = 2 //
	BATTLE_TYPE_BOSS = 3 //

	REDIS_TABLE_NAME_FightMatch   = "FightMatch_fightinfo"
	REDIS_TABLE_NAME_BATTLEINFO   = "FightMatchbattleinfo"
	REDIS_TABLE_NAME_BATTLERECORD = "FightMatchbattlerecord"
)

type Js_FightMatchUserDB struct {
	Uid    int64  `json:"uid"`
	KeyId  int    `json:"keyid"`
	Period int    `json:"period"`
	SvrId  int    `json:"svrid"`
	Info   string `json:"info"`

	info *Js_FightMatchUser
	db.DataUpdate
}

type Js_FightMatchUser struct {
	Uid         int64              `json:"uid"`
	SvrId       int                `json:"svrid"`
	SvrName     string             `json:"svrname"`
	Stage       int                `json:"stage"` //比赛阶段
	Group       int                `json:"stage"` //比赛小组
	Rank        int                `json:"rank"`  //比赛名次
	UName       string             `json:"uname"`
	Level       int                `json:"level"`
	Vip         int                `json:"vip"`
	Icon        int                `json:"icon"`
	Portrait    int                `json:"portrait"`
	Fight       int64              `json:"fight"`
	Robot       int                `json:"robot"`
	FightRecord []*game.ArenaFight `json:"fightrecord"` //战报集
}

type FightMatchTimeConfig struct {
	Uid         int64              `json:"uid"`
	SvrId       int                `json:"svrid"`
	SvrName     string             `json:"svrname"`
	Stage       int                `json:"stage"` //比赛阶段
	Group       int                `json:"stage"` //比赛小组
	Rank        int                `json:"rank"`  //比赛名次
	UName       string             `json:"uname"`
	Level       int                `json:"level"`
	Vip         int                `json:"vip"`
	Icon        int                `json:"icon"`
	Portrait    int                `json:"portrait"`
	Fight       int64              `json:"fight"`
	Robot       int                `json:"robot"`
	FightRecord []*game.ArenaFight `json:"fightrecord"` //战报集
}

type FightMatchInfo struct {
	KeyId      int    `json:"keyid"`      //赛区，对应分期表
	Period     int    `json:"period"`     //第几期
	TimeConfig string `json:"timeconfig"` //时间配置

	Mu          *sync.RWMutex
	arenaFights map[int64]*Js_FightMatchUser                        //所有人数据  (涉及指针切换到话需要所有变量一起同步)
	fightInfos  map[int64][FIGHTMATCH_TEAM_COUNT]*game.JS_FightInfo //竞技场防守阵容
	db_list     map[int64]*Js_FightMatchUserDB                      //数据存储
	timeConfig  *FightMatchTimeConfig

	db.DataUpdate
}

type FightMatchMgr struct {
	Locker         *sync.RWMutex
	FightMatchInfo map[int]*FightMatchInfo

	//ArenaSubsectionConfigs    []*ArenaSubsectionConfig
	//ArenaSubsectionConfigsMap map[int]*ArenaSubsectionConfig
	//ArenaMatchConfigs         []*ArenaMatchConfigs
	//JJCRobotConfig            []*JJCRobotConfig
}

var fightMatchMgr *FightMatchMgr = nil

func GetFightMatchMgr() *FightMatchMgr {
	if fightMatchMgr == nil {
		fightMatchMgr = new(FightMatchMgr)
		fightMatchMgr.FightMatchInfo = make(map[int]*FightMatchInfo)
		fightMatchMgr.Locker = new(sync.RWMutex)
		fightMatchMgr.LoadCsv()
	}
	return fightMatchMgr
}

func (self *FightMatchMgr) LoadCsv() {

}

func (self *Js_FightMatchUserDB) Encode() {
	self.Info = utils.HF_JtoA(self.info)
}

func (self *Js_FightMatchUserDB) Decode() {
	json.Unmarshal([]byte(self.Info), &self.info)
}

//
func (self *FightMatchInfo) Encode() {
	self.TimeConfig = utils.HF_JtoA(self.timeConfig)
}

func (self *FightMatchInfo) Decode() {
	json.Unmarshal([]byte(self.TimeConfig), &self.timeConfig)
}

// 存储数据库
func (self *FightMatchMgr) OnSave() {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for _, v := range self.FightMatchInfo {
		v.Save()
	}
}

func (self *FightMatchInfo) Save() {

	self.Encode()
	self.Update(true, false)

	for _, v := range self.db_list {
		if v.info.Robot == game.LOGIC_TRUE {
			continue
		}
		v.Encode()
		v.Update(true, false)
		//utils.LogDebug("保存玩家:", v.Uid, "期数:", v.KeyId)
	}
}

func (self *FightMatchMgr) GetAllData() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	queryStr := fmt.Sprintf("select * from `tbl_fightmatchinfo`;")
	var msg FightMatchInfo
	res := db.GetDBMgr().DBUser.GetAllData(queryStr, &msg)

	for i := 0; i < len(res); i++ {
		data := res[i].(*FightMatchInfo)
		if data.KeyId <= 0 {
			continue
		}

		_, ok := self.FightMatchInfo[data.KeyId]
		if !ok {
			self.FightMatchInfo[data.KeyId] = self.NewFightMatch(data.KeyId, data.Period)
		}

		if self.FightMatchInfo[data.KeyId] == nil {
			continue
		}

		data.Decode()
		if data.timeConfig == nil {
			data.timeConfig = self.CalTimeConfig()
		}
	}
}

func (self *FightMatchMgr) NewFightMatch(keyId int, Period int) *FightMatchInfo {
	data := new(FightMatchInfo)
	data.KeyId = keyId
	data.Period = Period
	data.arenaFights = make(map[int64]*Js_FightMatchUser)
	data.fightInfos = make(map[int64][FIGHTMATCH_TEAM_COUNT]*game.JS_FightInfo)
	data.db_list = make(map[int64]*Js_FightMatchUserDB)
	data.Mu = new(sync.RWMutex)
	return data
}

func (self *FightMatchMgr) CalTimeConfig() *FightMatchTimeConfig {
	data := new(FightMatchTimeConfig)

	return data
}

func (self *FightMatchMgr) GetInfo(keyId int) *FightMatchInfo {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	info, ok := self.FightMatchInfo[keyId]
	if !ok {
		info = self.NewFightMatch(keyId)
		self.FightMatchInfo[info.KeyId] = info
	}

	return self.FightMatchInfo[keyId]
}

func (self *FightMatchMgr) GetAllRank(req *RPC_FightMatchActionReqAll, res *RPC_FightMatchActionRes) {
	info := self.GetInfo(req.KeyId)
	info.GetAllRank(req, res)
}

func (self *FightMatchInfo) GetAllRank(req *RPC_FightMatchActionReqAll, res *RPC_FightMatchActionRes) {
	res.RetCode = RETCODE_OK
	self.Mu.RLock()
	defer self.Mu.RUnlock()
	res.RankInfo = self.FightMatchUserTop
}

func (self *FightMatchMgr) AddInfo(req *RPC_FightMatchActionReq, res *RPC_FightMatchActionRes) {
	info := self.GetInfo(req.KeyId)
	info.AddInfo(req, res)
}

func (self *FightMatchInfo) AddInfo(req *RPC_FightMatchActionReq, res *RPC_FightMatchActionRes) {
	self.Mu.Lock()
	defer self.Mu.Unlock()
	_, ok := self.arenaFights[req.Uid]
	if !ok {
		_, okSubsection := self.FightMatchUserTop[req.SelfInfo.Subsection]
		if !okSubsection {
			res.RetCode = RETCODE_DATA_ERROR
			return
		}

		_, okClass := self.FightMatchUserTop[req.SelfInfo.Subsection][req.SelfInfo.Class]
		if !okClass {
			res.RetCode = RETCODE_DATA_ERROR
			return
		}
		self.arenaFights[req.Uid] = req.SelfInfo
		dbData := new(Js_FightMatchUserDB)
		dbData.info = req.SelfInfo
		dbData.Uid = req.SelfInfo.Uid
		dbData.KeyId = self.KeyId
		dbData.SvrId = req.SelfInfo.SvrId
		dbData.Encode()
		self.db_list[req.Uid] = dbData
		db.InsertTable("tbl_FightMatchex", dbData, 0, true)
		self.FightMatchUserTop[req.SelfInfo.Subsection][req.SelfInfo.Class] = append(self.FightMatchUserTop[req.SelfInfo.Subsection][req.SelfInfo.Class], req.SelfInfo)
	} else if req.SelfInfo != nil {
		//这个地方同步
		self.arenaFights[req.Uid].Icon = req.SelfInfo.Icon
		self.arenaFights[req.Uid].Portrait = req.SelfInfo.Portrait
		self.arenaFights[req.Uid].Fight = req.SelfInfo.Fight
		self.arenaFights[req.Uid].Level = req.SelfInfo.Level
		self.arenaFights[req.Uid].UName = req.SelfInfo.UName
		self.arenaFights[req.Uid].SvrId = req.SelfInfo.SvrId
		self.db_list[req.Uid].SvrId = req.SelfInfo.SvrId
	}

	self.fightInfos[req.Uid] = req.FightInfo
	db.HMSetRedis(REDIS_TABLE_NAME_FightMatch, req.Uid, req.FightInfo, utils.DAY_SECS*10)
	res.RetCode = RETCODE_OK
	res.RankInfo = self.GetTopSafe()
	res.SelfInfo = self.arenaFights[req.Uid]
}

func (self *FightMatchInfo) GetTopSafe() map[int]map[int][]*Js_FightMatchUser {
	data := make(map[int]map[int][]*Js_FightMatchUser)

	_, okSub := self.FightMatchUserTop[BEST_SUBSECTION]
	if okSub {
		data[BEST_SUBSECTION] = make(map[int][]*Js_FightMatchUser)
		data[BEST_SUBSECTION] = self.FightMatchUserTop[BEST_SUBSECTION]
		return data
	}
	return nil
}

func (self *FightMatchMgr) GetDefence(req *RPC_FightMatchActionReq, res *RPC_FightMatchGetDefenceRes) {
	info := self.GetInfo(req.KeyId)
	info.GetDefence(req, res)
}

func (self *FightMatchInfo) GetDefence(req *RPC_FightMatchActionReq, res *RPC_FightMatchGetDefenceRes) {
	self.Mu.RLock()
	defer self.Mu.RUnlock()
	info, ok := self.arenaFights[req.Uid]
	if !ok {
		res.RetCode = RETCODE_DATA_ERROR
		return
	}
	//生成匹配对手的配置和自己的编号方便后面计算
	config := make([]*ArenaMatchConfigs, 0)
	for _, v := range GetFightMatchMgr().ArenaMatchConfigs {
		if v.Subsection == info.Subsection && v.Class == info.Class {
			config = append(config, v)
		}
	}
	//如果没有合适的配置则生成统一配置
	if len(config) == 0 {
		for _, v := range GetFightMatchMgr().ArenaMatchConfigs {
			if v.Subsection == 0 && v.Class == 0 {
				config = append(config, v)
			}
		}
	}
	//生成自己的配置编号方便计算
	selfConfigId := self.GetConfigId(info.Subsection, info.Class)
	//生成对手,生成豁免列表
	cantUse := make(map[int64]int)
	cantUse[info.Uid] = game.LOGIC_TRUE
	for _, v := range config {
		//备选列表
		lst := self.GetLst(selfConfigId-v.FirstMin, selfConfigId-v.FirstMax, cantUse)
		if len(lst) == 0 {
			lst = self.GetLst(selfConfigId-v.SupplementMin, selfConfigId-v.SupplementMax, cantUse)
		}
		size := len(lst)
		if size == 0 {
			return
		}
		rand := game.HF_GetRandom(size)
		if lst[rand].Robot == game.LOGIC_FALSE {
			fightInfo, okFightInfo := self.fightInfos[lst[rand].Uid]
			if okFightInfo {
				res.Info = append(res.Info, lst[rand])
				res.FightInfo = append(res.FightInfo, fightInfo)
				cantUse[lst[rand].Uid] = game.LOGIC_TRUE
			} else {
				print("中心服丢失战斗数据")
			}
		} else {
			//机器人加入
			res.Info = append(res.Info, lst[rand])
			cantUse[lst[rand].Uid] = game.LOGIC_TRUE
		}
	}

	res.RetCode = RETCODE_OK
}

func (self *FightMatchInfo) GetConfigId(subsection int, class int) int {
	for _, v := range GetFightMatchMgr().ArenaSubsectionConfigs {
		if v.Subsection == subsection && v.Class == class {
			return v.Id
		}
	}
	return 0
}

func (self *FightMatchInfo) GetLst(startId int, endId int, cantUse map[int64]int) []*Js_FightMatchUser {
	lst := make([]*Js_FightMatchUser, 0)

	for i := endId; i <= startId; i++ {
		config := GetFightMatchMgr().ArenaSubsectionConfigsMap[i]
		if config == nil {
			continue
		}
		_, okSub := self.FightMatchUserTop[config.Subsection]
		if okSub {
			_, okClass := self.FightMatchUserTop[config.Subsection][config.Class]
			if okClass {
				for _, v := range self.FightMatchUserTop[config.Subsection][config.Class] {
					if cantUse[v.Uid] != game.LOGIC_TRUE {
						lst = append(lst, v)
					}
				}
			}
		}
	}
	return lst
}

func (self *FightMatchMgr) GetPlayerInfo(req *RPC_FightMatchActionReq, res *RPC_FightMatchGetInfoRes) {
	info := self.GetInfo(req.KeyId)
	info.GetPlayerInfo(req, res)
}

func (self *FightMatchInfo) GetPlayerInfo(req *RPC_FightMatchActionReq, res *RPC_FightMatchGetInfoRes) {
	self.Mu.RLock()
	defer self.Mu.RUnlock()
	info, ok := self.arenaFights[req.Uid]
	if !ok {
		res.RetCode = RETCODE_DATA_ERROR
		return
	}

	res.Info = info
	res.FightInfo = self.fightInfos[info.Uid]
	if info.Robot == game.LOGIC_FALSE {
		player := core.GetPlayerMgr().GetCorePlayer(info.Uid, true)
		if player != nil {
			res.LifeTreeInfo = player.GetLifeTree()
		}
	}
	res.RetCode = RETCODE_OK
}

func (self *FightMatchMgr) FightEnd(req *RPC_FightMatchFightEndReq, res *RPC_FightMatchActionRes) {
	info := self.GetInfo(req.KeyId)
	info.FightEnd(req, res)
}

func (self *FightMatchMgr) GetRecordId() int64 {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	self.RecordId++
	RecordId := (time.Now().Unix()%10000000)*100 + self.RecordId%100

	return RecordId
}

func (self *FightMatchInfo) FightEnd(req *RPC_FightMatchFightEndReq, res *RPC_FightMatchActionRes) {
	self.Mu.Lock()
	defer self.Mu.Unlock()
	attack := self.arenaFights[req.Attack.Uid]
	defence := self.arenaFights[req.Defend.Uid]
	if attack == nil || defence == nil {
		res.RetCode = RETCODE_DATA_ERROR
		return
	}

	res.NewFightId = GetFightMatchMgr().GetRecordId()
	req.BattleInfo.Id = res.NewFightId

	if req.BattleInfo.Result == 0 {
		attack.AddFight(req.Attack, req.Defend, 0, 1, req.BattleInfo, defence.Subsection, defence.Class, true)
		if defence.Robot == game.LOGIC_FALSE {
			defence.AddFight(req.Attack, req.Defend, 0, 0, req.BattleInfo, attack.Subsection, attack.Class, false)
		}
		//如果需要更新位置
		if attack.IsCanUp(defence) {
			self.SwapPosSafe(attack, defence)
		}
	} else {
		attack.AddFight(req.Attack, req.Defend, 1, 1, req.BattleInfo, defence.Subsection, defence.Class, true)
		if defence.Robot == game.LOGIC_FALSE {
			defence.AddFight(req.Attack, req.Defend, 1, 0, req.BattleInfo, attack.Subsection, attack.Class, false)
		}
	}
	//处理排行问题
	res.RankInfo = self.GetTopSafe()
	res.SelfInfo = attack
	res.Result = req.BattleInfo.Result
	res.RetCode = RETCODE_OK

	//推送给防守方
	if defence.Robot == game.LOGIC_FALSE {
		core.GetCenterApp().AddEvent(defence.SvrId, core.MATCH_FightMatch_UPDATE, defence.Uid,
			0, 0, utils.HF_JtoA(defence))
	}
}

func (self *FightMatchInfo) SwapPosSafe(attack *Js_FightMatchUser, defence *Js_FightMatchUser) {

	//记录攻击方之前段位
	oldAttackSub := attack.Subsection
	oldAttackClass := attack.Class

	_, okDefenceSub := self.FightMatchUserTop[defence.Subsection]
	if okDefenceSub {
		_, okDefenceClass := self.FightMatchUserTop[defence.Subsection][defence.Class]
		if okDefenceClass {
			index := -1
			for i := 0; i < len(self.FightMatchUserTop[defence.Subsection][defence.Class]); i++ {
				if self.FightMatchUserTop[defence.Subsection][defence.Class][i].Uid == defence.Uid {
					index = i
					break
				}
			}
			if index == -1 {
				//没找到,数据异常
				return
			}
			//攻击方上位
			attack.Subsection = defence.Subsection
			attack.Class = defence.Class
			self.FightMatchUserTop[defence.Subsection][defence.Class][index] = attack
			self.arenaFights[attack.Uid] = attack
			self.db_list[attack.Uid].info = attack
		}
	}

	_, okAttackSub := self.FightMatchUserTop[oldAttackSub]
	if okAttackSub {
		_, okAttackClass := self.FightMatchUserTop[oldAttackSub][oldAttackClass]
		if okAttackClass {
			index := -1
			for i := 0; i < len(self.FightMatchUserTop[oldAttackSub][oldAttackClass]); i++ {
				if self.FightMatchUserTop[oldAttackSub][oldAttackClass][i].Uid == attack.Uid {
					index = i
					break
				}
			}
			if index == -1 {
				//没找到,数据异常
				return
			}
			if defence.Robot == game.LOGIC_TRUE {
				//生成对应位置的机器人数据
				for _, config := range GetFightMatchMgr().ArenaSubsectionConfigs {
					if config.Subsection != oldAttackSub || config.Class != oldAttackClass {
						continue
					}
					defence.UName = config.Name + "守卫"

					for _, configJJC := range GetFightMatchMgr().JJCRobotConfig {
						if configJJC.Type != 3 {
							continue
						}
						if configJJC.Jjcclass == oldAttackSub && configJJC.Jjcdan == oldAttackClass {
							defence.Level = configJJC.Level
							break
						}
					}
				}
			}
			defence.Subsection = oldAttackSub
			defence.Class = oldAttackClass
			self.FightMatchUserTop[oldAttackSub][oldAttackClass][index] = defence
			self.arenaFights[defence.Uid] = defence
			if defence.Robot == game.LOGIC_FALSE {
				_, ok := self.db_list[defence.Uid]
				if ok {
					self.db_list[defence.Uid].info = defence
				}
			}
		}
	}
}

func (self *Js_FightMatchUser) IsCanUp(defence *Js_FightMatchUser) bool {
	if self.Subsection < defence.Subsection {
		return false
	}
	if self.Subsection > defence.Subsection {
		return true
	}

	return self.Class > defence.Class
}

func (self *Js_FightMatchUser) AddFight(attack *game.JS_FightInfo, defend *game.JS_FightInfo, result int, side int,
	battleInfo game.BattleInfo, subscetion int, class int, needSaveSql bool) {

	var enemy *game.JS_FightInfo
	if side == 1 {
		enemy = defend
	} else {
		enemy = attack
	}

	fight := self.NewPvpFight(battleInfo.Id, enemy, result, side, subscetion, class)
	if len(self.FightRecord) >= RECORD_MAX_CROSS {
		self.FightRecord = self.FightRecord[1:]
	}
	self.FightRecord = append(self.FightRecord, fight)

	data2 := game.BattleRecord{}
	data2.Level = 0
	data2.Side = side
	data2.Time = time.Now().Unix()
	data2.Id = battleInfo.Id
	data2.LevelID = battleInfo.LevelID
	data2.Result = result
	data2.Type = BATTLE_TYPE_PVP
	data2.RandNum = battleInfo.Random
	data2.FightInfo[0] = attack
	data2.FightInfo[1] = defend

	if needSaveSql {
		var db_battleInfo JS_CrossArenaBattleInfo
		db_battleInfo.FightId = battleInfo.Id
		db_battleInfo.RecordType = core.BATTLE_TYPE_RECORD_FightMatch
		db_battleInfo.BattleInfo = utils.HF_CompressAndBase64(game.HF_JtoB(&battleInfo))
		db_battleInfo.BattleRecord = utils.HF_CompressAndBase64(game.HF_JtoB(&data2))
		db_battleInfo.UpdateTime = time.Now().Unix()
		db.InsertTable("tbl_FightMatchrecord", &db_battleInfo, 0, false)

		db.HMSetRedisEx(REDIS_TABLE_NAME_BATTLEINFO, battleInfo.Id, &battleInfo, 3600*12)
		db.HMSetRedisEx(REDIS_TABLE_NAME_BATTLERECORD, data2.Id, &data2, 3600*12)
	}
}

func (self *Js_FightMatchUser) NewPvpFight(FightID int64, enemy *game.JS_FightInfo, result int, side int, subscetion int, class int) *game.ArenaFight {

	p := &game.ArenaFight{}
	p.FightId = FightID
	p.Side = side
	p.Result = result
	p.Subsection = subscetion
	p.Class = class
	if enemy != nil {
		p.Uid = enemy.Uid
		p.IconId = enemy.Iconid
		p.Name = enemy.Uname
		p.Level = enemy.Level
		p.Fight = enemy.Deffight
		p.Portrait = enemy.Portrait
	}
	p.Time = time.Now().Unix()

	return p
}

func (self *FightMatchMgr) GetBattleInfo(req *RPC_FightMatchAction64ReqAll, res *RPC_FightMatchBattleInfoRes) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var battleInfo game.BattleInfo
	value, flag, err := db.HGetRedisEx(REDIS_TABLE_NAME_BATTLEINFO, req.KeyId, fmt.Sprintf("%d", req.KeyId))
	if err != nil || !flag {
		//redis优化    读取简报
		//! 全部获取
		var db_battleinfo JS_CrossArenaBattleInfo
		sql := fmt.Sprintf("select * from `tbl_FightMatchrecord` where fightid=%d limit 1;", req.KeyId)
		ret := db.GetDBMgr().DBUser.GetOneData(sql, &db_battleinfo, "", 0)
		if ret == true { //! 获取成功
			//! 进行处理
			err := json.Unmarshal(utils.HF_Base64AndDecompress(db_battleinfo.BattleInfo), &battleInfo)
			if err != nil {
				utils.LogDebug("Decode Error")
				return
			}

			if battleInfo.Id != 0 {
				res.BattleInfo = &battleInfo
				return
			}
			//! 详细处理
		}
		return
	}
	err1 := json.Unmarshal([]byte(value), &battleInfo)
	if err1 != nil {
		return
	}

	if battleInfo.Id != 0 {
		res.BattleInfo = &battleInfo
		return
	}
}

func (self *FightMatchMgr) GetBattleRecord(req *RPC_FightMatchAction64ReqAll, res *RPC_FightMatchBattleRecordRes) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var battleRecord game.BattleRecord
	value, flag, err := db.HGetRedisEx(REDIS_TABLE_NAME_BATTLERECORD, req.KeyId, fmt.Sprintf("%d", req.KeyId))
	if err != nil || !flag {
		//redis优化    读取简报
		//! 全部获取
		var db_battleinfo JS_CrossArenaBattleInfo
		sql := fmt.Sprintf("select * from `tbl_FightMatchrecord` where fightid=%d limit 1;", req.KeyId)
		ret := db.GetDBMgr().DBUser.GetOneData(sql, &db_battleinfo, "", 0)
		if ret == true { //! 获取成功
			//! 进行处理
			err := json.Unmarshal(utils.HF_Base64AndDecompress(db_battleinfo.BattleRecord), &battleRecord)
			if err != nil {
				utils.LogDebug("Decode Error")
				return
			}

			if battleRecord.Id != 0 {
				res.BattleRecord = &battleRecord
				return
			}
			//! 详细处理
		}
		return
	}
	err1 := json.Unmarshal([]byte(value), &battleRecord)
	if err1 != nil {
		return
	}
	if battleRecord.Id != 0 {
		res.BattleRecord = &battleRecord
		return
	}
}
*/
