package match

import (
	"encoding/json"
	"fmt"
	"game"
	"master/core"
	"master/db"
	"master/utils"
	"sync"
	"time"
)

const (
	BEST_SUBSECTION  = 1  //传奇段位
	RECORD_MAX_CROSS = 10 //最大战报数s

	BATTLE_TYPE_PVP  = 1 //
	BATTLE_TYPE_PVE  = 2 //
	BATTLE_TYPE_BOSS = 3 //

	REDIS_TABLE_NAME_CROSSARENA = "crossarena_fightinfo"

	REDIS_TABLE_NAME_BATTLEINFO   = "crossarenabattleinfo"
	REDIS_TABLE_NAME_BATTLERECORD = "crossarenabattlerecord"
)

type Js_CrossArenaUserDB struct {
	Uid   int64  `json:"uid"`
	KeyId int    `json:"keyid"`
	SvrId int    `json:"svrid"`
	Info  string `json:"info"`

	info *Js_CrossArenaUser
	db.DataUpdate
}

type Js_CrossArenaUser struct {
	Uid         int64              `json:"uid"`
	SvrId       int                `json:"svrid"`
	SvrName     string             `json:"svrname"`
	Subsection  int                `json:"subsection"` //大段位
	Class       int                `json:"class"`      //小段位
	UName       string             `json:"uname"`
	Level       int                `json:"level"`
	Vip         int                `json:"vip"`
	Icon        int                `json:"icon"`
	Portrait    int                `json:"portrait"`
	Fight       int64              `json:"fight"`
	Robot       int                `json:"robot"`
	FightRecord []*game.ArenaFight `json:"fightrecord"` //战报集
}

type CrossArenaInfo struct {
	KeyId int `json:"keyid"`

	Mu                *sync.RWMutex
	crossArenaUserTop map[int]map[int][]*Js_CrossArenaUser //map[大段位][小段位]玩家数据
	arenaFights       map[int64]*Js_CrossArenaUser         //所有人数据  (涉及指针切换到话需要所有变量一起同步)
	fightInfos        map[int64]*game.JS_FightInfo         //竞技场防守阵容
	db_list           map[int64]*Js_CrossArenaUserDB       //数据存储
}

type JS_CrossArenaBattleInfo struct {
	Id           int    `json:"id"`           //! 自增Id
	FightId      int64  `json:"fightid"`      //! 战斗Id
	RecordType   int    `json:"recordtype"`   //! 战报类型
	BattleInfo   string `json:"battleinfo"`   //! 简报
	BattleRecord string `json:"battlerecord"` //! 详细战报
	UpdateTime   int64  `json:"updatetime"`   //! 插入时间

	db.DataUpdate
}

type ArenaSubsectionConfig struct {
	Id         int    `json:"id"`
	Subsection int    `json:"subsection"`
	Class      int    `json:"class"`
	Name       string `json:"name"`
	Capacity   int    `json:"capacity"`
}

type ArenaMatchConfigs struct {
	Id            int `json:"id"`
	System        int `json:"system"`
	Subsection    int `json:"subsection"`
	Class         int `json:"class"`
	Opponent      int `json:"opponent"`
	FirstMax      int `json:"firstmax"`
	FirstMin      int `json:"firstmin"`
	SupplementMax int `json:"supplementmax"`
	SupplementMin int `json:"supplementmin"`
}

type JJCRobotConfig struct {
	Id        int   `json:"id"`
	Type      int   `json:"type"`
	Jjcclass  int   `json:"jjcclass"`
	Jjcdan    int   `json:"jjcdan"`
	ShowFight int64 `json:"showfight"`
	Level     int   `json:"robotlevel"`
	Hero      []int `json:"optionhero"`
}

type CrossArenaMgr struct {
	Locker         *sync.RWMutex
	CrossArenaInfo map[int]*CrossArenaInfo

	ArenaSubsectionConfigs    []*ArenaSubsectionConfig
	ArenaSubsectionConfigsMap map[int]*ArenaSubsectionConfig
	ArenaMatchConfigs         []*ArenaMatchConfigs
	JJCRobotConfig            []*JJCRobotConfig
	RecordId                  int64
}

var crossArenaMgr *CrossArenaMgr = nil

func GetCrossArenaMgr() *CrossArenaMgr {
	if crossArenaMgr == nil {
		crossArenaMgr = new(CrossArenaMgr)
		crossArenaMgr.CrossArenaInfo = make(map[int]*CrossArenaInfo)
		crossArenaMgr.Locker = new(sync.RWMutex)
		crossArenaMgr.LoadCsv()
	}
	return crossArenaMgr
}

func (self *CrossArenaMgr) LoadCsv() {
	utils.GetCsvUtilMgr().LoadCsv("Activity_Arenasubsection", &self.ArenaSubsectionConfigs)
	self.ArenaSubsectionConfigsMap = make(map[int]*ArenaSubsectionConfig)
	for _, v := range self.ArenaSubsectionConfigs {
		self.ArenaSubsectionConfigsMap[v.Id] = v
	}

	utils.GetCsvUtilMgr().LoadCsv("Activity_Arenamatching", &self.ArenaMatchConfigs)

	JJCRobotConfigTemp := make([]*JJCRobotConfig, 0)
	utils.GetCsvUtilMgr().LoadCsv("Jjc_Robot", &JJCRobotConfigTemp)
	for _, v := range JJCRobotConfigTemp {
		if v.Type != 3 {
			continue
		}
		self.JJCRobotConfig = append(self.JJCRobotConfig, v)
	}
}

func (self *Js_CrossArenaUserDB) Encode() {
	self.Info = utils.HF_JtoA(self.info)
}

func (self *Js_CrossArenaUserDB) Decode() {
	json.Unmarshal([]byte(self.Info), &self.info)
}

// 存储数据库
func (self *CrossArenaMgr) OnSave() {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for _, v := range self.CrossArenaInfo {
		v.Save()
	}
}

func (self *CrossArenaInfo) Save() {
	for _, v := range self.db_list {
		if v.info.Robot == game.LOGIC_TRUE {
			continue
		}
		v.Encode()
		v.Update(true, false)

		//utils.LogDebug("保存玩家:", v.Uid, "期数:", v.KeyId)
	}
}

func (self *CrossArenaMgr) GetAllData() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	queryStr := fmt.Sprintf("select uid,keyid,svrid,info from `tbl_crossarenaex`;")
	var msg Js_CrossArenaUserDB
	res := db.GetDBMgr().DBUser.GetAllData(queryStr, &msg)

	//utils.LogDebug("载入数据:", len(res))

	for i := 0; i < len(res); i++ {
		data := res[i].(*Js_CrossArenaUserDB)
		if data.KeyId <= 0 {
			continue
		}

		_, ok := self.CrossArenaInfo[data.KeyId]
		if !ok {
			self.CrossArenaInfo[data.KeyId] = self.NewCrossArena(data.KeyId)
		}

		if self.CrossArenaInfo[data.KeyId] == nil {
			continue
		}
		data.Decode()
		if data.info == nil {
			continue
		}
		//玩家数据拷贝上去
		_, okDefenceSub := self.CrossArenaInfo[data.KeyId].crossArenaUserTop[data.info.Subsection]
		if okDefenceSub {
			_, okDefenceClass := self.CrossArenaInfo[data.KeyId].crossArenaUserTop[data.info.Subsection][data.info.Class]
			if okDefenceClass {
				index := -1
				for i := 0; i < len(self.CrossArenaInfo[data.KeyId].crossArenaUserTop[data.info.Subsection][data.info.Class]); i++ {
					if self.CrossArenaInfo[data.KeyId].crossArenaUserTop[data.info.Subsection][data.info.Class][i].Robot == game.LOGIC_TRUE {
						index = i
						break
					}
				}

				data.Init("tbl_crossarenaex", data, true)

				self.CrossArenaInfo[data.KeyId].arenaFights[data.Uid] = data.info
				self.CrossArenaInfo[data.KeyId].db_list[data.Uid] = data
				if index == -1 {
					//直接添加,到这里都是最后一组了
					self.CrossArenaInfo[data.KeyId].crossArenaUserTop[data.info.Subsection][data.info.Class] = append(self.CrossArenaInfo[data.KeyId].crossArenaUserTop[data.info.Subsection][data.info.Class], data.info)
				} else {
					self.CrossArenaInfo[data.KeyId].crossArenaUserTop[data.info.Subsection][data.info.Class][index] = data.info
				}

				//utils.LogDebug("载入玩家:", data.Uid, "期数:", data.KeyId)

				fightInfo := game.JS_FightInfo{}
				value, flag, err := db.HGetRedis(REDIS_TABLE_NAME_CROSSARENA, fmt.Sprintf("%d", data.Uid))
				if err != nil {
					continue
				}
				if flag {
					err := json.Unmarshal([]byte(value), &fightInfo)
					if err != nil {
						continue
					}
				}
				self.CrossArenaInfo[data.KeyId].fightInfos[data.Uid] = &fightInfo
			}
		}
	}
}

func (self *CrossArenaMgr) NewCrossArena(keyId int) *CrossArenaInfo {
	data := new(CrossArenaInfo)
	data.KeyId = keyId
	data.crossArenaUserTop = make(map[int]map[int][]*Js_CrossArenaUser)
	data.arenaFights = make(map[int64]*Js_CrossArenaUser)
	data.fightInfos = make(map[int64]*game.JS_FightInfo)
	data.db_list = make(map[int64]*Js_CrossArenaUserDB)
	data.Mu = new(sync.RWMutex)

	uid := int64(0)
	//初始化排行
	for _, v := range self.ArenaSubsectionConfigs {
		_, okSubsection := data.crossArenaUserTop[v.Subsection]
		if !okSubsection {
			data.crossArenaUserTop[v.Subsection] = make(map[int][]*Js_CrossArenaUser)
		}

		_, okClass := data.crossArenaUserTop[v.Subsection][v.Class]
		if !okClass {
			data.crossArenaUserTop[v.Subsection][v.Class] = make([]*Js_CrossArenaUser, 0)
			for i := 0; i < v.Capacity; i++ {
				uid++
				userData := new(Js_CrossArenaUser)
				userData.Uid = uid
				userData.Subsection = v.Subsection
				userData.Class = v.Class
				userData.Robot = game.LOGIC_TRUE
				userData.UName = v.Name + "守卫"
				for _, config := range self.JJCRobotConfig {
					if config.Type != 3 {
						continue
					}
					if config.Jjcclass == userData.Subsection && config.Jjcdan == userData.Class {
						userData.Fight = config.ShowFight
						rand := game.HF_GetRandom(10000)
						if rand >= 5000 {
							userData.Icon = 1002
						} else {
							userData.Icon = 1003
						}
						userData.Level = config.Level
						break
					}
				}

				data.crossArenaUserTop[v.Subsection][v.Class] = append(data.crossArenaUserTop[v.Subsection][v.Class], userData)
				data.arenaFights[userData.Uid] = userData
			}
		} else {
			utils.LogDebug("表段位填写重复!")
		}
	}
	return data
}

func (self *CrossArenaMgr) GetInfo(keyId int) *CrossArenaInfo {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	info, ok := self.CrossArenaInfo[keyId]
	if !ok {
		info = self.NewCrossArena(keyId)
		self.CrossArenaInfo[info.KeyId] = info
	}

	return self.CrossArenaInfo[keyId]
}

func (self *CrossArenaMgr) GetAllRank(req *RPC_CrossArenaActionReqAll, res *RPC_CrossArenaActionRes) {
	info := self.GetInfo(req.KeyId)
	info.GetAllRank(req, res)
}

func (self *CrossArenaInfo) GetAllRank(req *RPC_CrossArenaActionReqAll, res *RPC_CrossArenaActionRes) {
	res.RetCode = RETCODE_OK
	self.Mu.RLock()
	defer self.Mu.RUnlock()
	res.RankInfo = self.crossArenaUserTop
}

func (self *CrossArenaMgr) AddInfo(req *RPC_CrossArenaActionReq, res *RPC_CrossArenaActionRes) {
	info := self.GetInfo(req.KeyId)
	info.AddInfo(req, res)
}

func (self *CrossArenaInfo) AddInfo(req *RPC_CrossArenaActionReq, res *RPC_CrossArenaActionRes) {
	self.Mu.Lock()
	defer self.Mu.Unlock()
	_, ok := self.arenaFights[req.Uid]
	if !ok {
		_, okSubsection := self.crossArenaUserTop[req.SelfInfo.Subsection]
		if !okSubsection {
			res.RetCode = RETCODE_DATA_ERROR
			return
		}

		_, okClass := self.crossArenaUserTop[req.SelfInfo.Subsection][req.SelfInfo.Class]
		if !okClass {
			res.RetCode = RETCODE_DATA_ERROR
			return
		}
		self.arenaFights[req.Uid] = req.SelfInfo
		dbData := new(Js_CrossArenaUserDB)
		dbData.info = req.SelfInfo
		dbData.Uid = req.SelfInfo.Uid
		dbData.KeyId = self.KeyId
		dbData.SvrId = req.SelfInfo.SvrId
		dbData.Encode()
		self.db_list[req.Uid] = dbData
		db.InsertTable("tbl_crossarenaex", dbData, 0, true)
		self.crossArenaUserTop[req.SelfInfo.Subsection][req.SelfInfo.Class] = append(self.crossArenaUserTop[req.SelfInfo.Subsection][req.SelfInfo.Class], req.SelfInfo)
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
	db.HMSetRedis(REDIS_TABLE_NAME_CROSSARENA, req.Uid, req.FightInfo, utils.DAY_SECS*10)
	res.RetCode = RETCODE_OK
	res.RankInfo = self.GetTopSafe()
	res.SelfInfo = self.arenaFights[req.Uid]
}

func (self *CrossArenaInfo) GetTopSafe() map[int]map[int][]*Js_CrossArenaUser {
	data := make(map[int]map[int][]*Js_CrossArenaUser)

	_, okSub := self.crossArenaUserTop[BEST_SUBSECTION]
	if okSub {
		data[BEST_SUBSECTION] = make(map[int][]*Js_CrossArenaUser)
		data[BEST_SUBSECTION] = self.crossArenaUserTop[BEST_SUBSECTION]
		return data
	}
	return nil
}

func (self *CrossArenaMgr) GetDefence(req *RPC_CrossArenaActionReq, res *RPC_CrossArenaGetDefenceRes) {
	info := self.GetInfo(req.KeyId)
	info.GetDefence(req, res)
}

func (self *CrossArenaInfo) GetDefence(req *RPC_CrossArenaActionReq, res *RPC_CrossArenaGetDefenceRes) {
	self.Mu.RLock()
	defer self.Mu.RUnlock()
	info, ok := self.arenaFights[req.Uid]
	if !ok {
		res.RetCode = RETCODE_DATA_ERROR
		return
	}
	//生成匹配对手的配置和自己的编号方便后面计算
	config := make([]*ArenaMatchConfigs, 0)
	for _, v := range GetCrossArenaMgr().ArenaMatchConfigs {
		if v.Subsection == info.Subsection && v.Class == info.Class {
			config = append(config, v)
		}
	}
	//如果没有合适的配置则生成统一配置
	if len(config) == 0 {
		for _, v := range GetCrossArenaMgr().ArenaMatchConfigs {
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

func (self *CrossArenaInfo) GetConfigId(subsection int, class int) int {
	for _, v := range GetCrossArenaMgr().ArenaSubsectionConfigs {
		if v.Subsection == subsection && v.Class == class {
			return v.Id
		}
	}
	return 0
}

func (self *CrossArenaInfo) GetLst(startId int, endId int, cantUse map[int64]int) []*Js_CrossArenaUser {
	lst := make([]*Js_CrossArenaUser, 0)

	for i := endId; i <= startId; i++ {
		config := GetCrossArenaMgr().ArenaSubsectionConfigsMap[i]
		if config == nil {
			continue
		}
		_, okSub := self.crossArenaUserTop[config.Subsection]
		if okSub {
			_, okClass := self.crossArenaUserTop[config.Subsection][config.Class]
			if okClass {
				for _, v := range self.crossArenaUserTop[config.Subsection][config.Class] {
					if cantUse[v.Uid] != game.LOGIC_TRUE {
						lst = append(lst, v)
					}
				}
			}
		}
	}
	return lst
}

func (self *CrossArenaMgr) GetPlayerInfo(req *RPC_CrossArenaActionReq, res *RPC_CrossArenaGetInfoRes) {
	info := self.GetInfo(req.KeyId)
	info.GetPlayerInfo(req, res)
}

func (self *CrossArenaInfo) GetPlayerInfo(req *RPC_CrossArenaActionReq, res *RPC_CrossArenaGetInfoRes) {
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

func (self *CrossArenaMgr) FightEnd(req *RPC_CrossArenaFightEndReq, res *RPC_CrossArenaActionRes) {
	info := self.GetInfo(req.KeyId)
	info.FightEnd(req, res)
}

func (self *CrossArenaMgr) GetRecordId() int64 {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	self.RecordId++
	RecordId := (time.Now().Unix()%10000000)*100 + self.RecordId%100

	return RecordId
}

func (self *CrossArenaInfo) FightEnd(req *RPC_CrossArenaFightEndReq, res *RPC_CrossArenaActionRes) {
	self.Mu.Lock()
	defer self.Mu.Unlock()
	attack := self.arenaFights[req.Attack.Uid]
	defence := self.arenaFights[req.Defend.Uid]
	if attack == nil || defence == nil {
		res.RetCode = RETCODE_DATA_ERROR
		return
	}

	res.NewFightId = GetCrossArenaMgr().GetRecordId()
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
		core.GetCenterApp().AddEvent(defence.SvrId, core.MATCH_CROSSARENA_UPDATE, defence.Uid,
			0, 0, utils.HF_JtoA(defence))
	}
}

func (self *CrossArenaInfo) SwapPosSafe(attack *Js_CrossArenaUser, defence *Js_CrossArenaUser) {

	//记录攻击方之前段位
	oldAttackSub := attack.Subsection
	oldAttackClass := attack.Class

	_, okDefenceSub := self.crossArenaUserTop[defence.Subsection]
	if okDefenceSub {
		_, okDefenceClass := self.crossArenaUserTop[defence.Subsection][defence.Class]
		if okDefenceClass {
			index := -1
			for i := 0; i < len(self.crossArenaUserTop[defence.Subsection][defence.Class]); i++ {
				if self.crossArenaUserTop[defence.Subsection][defence.Class][i].Uid == defence.Uid {
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
			self.crossArenaUserTop[defence.Subsection][defence.Class][index] = attack
			self.arenaFights[attack.Uid] = attack
			self.db_list[attack.Uid].info = attack
		}
	}

	_, okAttackSub := self.crossArenaUserTop[oldAttackSub]
	if okAttackSub {
		_, okAttackClass := self.crossArenaUserTop[oldAttackSub][oldAttackClass]
		if okAttackClass {
			index := -1
			for i := 0; i < len(self.crossArenaUserTop[oldAttackSub][oldAttackClass]); i++ {
				if self.crossArenaUserTop[oldAttackSub][oldAttackClass][i].Uid == attack.Uid {
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
				for _, config := range GetCrossArenaMgr().ArenaSubsectionConfigs {
					if config.Subsection != oldAttackSub || config.Class != oldAttackClass {
						continue
					}
					defence.UName = config.Name + "守卫"

					for _, configJJC := range GetCrossArenaMgr().JJCRobotConfig {
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
			self.crossArenaUserTop[oldAttackSub][oldAttackClass][index] = defence
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

func (self *Js_CrossArenaUser) IsCanUp(defence *Js_CrossArenaUser) bool {
	if self.Subsection < defence.Subsection {
		return false
	}
	if self.Subsection > defence.Subsection {
		return true
	}

	return self.Class > defence.Class
}

func (self *Js_CrossArenaUser) AddFight(attack *game.JS_FightInfo, defend *game.JS_FightInfo, result int, side int,
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
		db_battleInfo.RecordType = core.BATTLE_TYPE_RECORD_CROSSARENA
		db_battleInfo.BattleInfo = utils.HF_CompressAndBase64(game.HF_JtoB(&battleInfo))
		db_battleInfo.BattleRecord = utils.HF_CompressAndBase64(game.HF_JtoB(&data2))
		db_battleInfo.UpdateTime = time.Now().Unix()
		db.InsertTable("tbl_crossarenarecord", &db_battleInfo, 0, false)

		db.HMSetRedisEx(REDIS_TABLE_NAME_BATTLEINFO, battleInfo.Id, &battleInfo, 3600*12)
		db.HMSetRedisEx(REDIS_TABLE_NAME_BATTLERECORD, data2.Id, &data2, 3600*12)
	}
}

func (self *Js_CrossArenaUser) NewPvpFight(FightID int64, enemy *game.JS_FightInfo, result int, side int, subscetion int, class int) *game.ArenaFight {

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

func (self *CrossArenaMgr) GetBattleInfo(req *RPC_CrossArenaAction64ReqAll, res *RPC_CrossArenaBattleInfoRes) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var battleInfo game.BattleInfo
	value, flag, err := db.HGetRedisEx(REDIS_TABLE_NAME_BATTLEINFO, req.KeyId, fmt.Sprintf("%d", req.KeyId))
	if err != nil || !flag {
		//redis优化    读取简报
		//! 全部获取
		var db_battleinfo JS_CrossArenaBattleInfo
		sql := fmt.Sprintf("select * from `tbl_crossarenarecord` where fightid=%d limit 1;", req.KeyId)
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

func (self *CrossArenaMgr) GetBattleRecord(req *RPC_CrossArenaAction64ReqAll, res *RPC_CrossArenaBattleRecordRes) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var battleRecord game.BattleRecord
	value, flag, err := db.HGetRedisEx(REDIS_TABLE_NAME_BATTLERECORD, req.KeyId, fmt.Sprintf("%d", req.KeyId))
	if err != nil || !flag {
		//redis优化    读取简报
		//! 全部获取
		var db_battleinfo JS_CrossArenaBattleInfo
		sql := fmt.Sprintf("select * from `tbl_crossarenarecord` where fightid=%d limit 1;", req.KeyId)
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
