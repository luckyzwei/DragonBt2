package game

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

const (
	ARENA_SPECIAL_BASE_CLASS = 7
	ARENA_SPECIAL_BASE_DAN   = 1
	ARENA_SPECIAL_ENEMY_BASE = 10
)

const (
	ARENA_SPECIAL_TEAM_1   = 0
	ARENA_SPECIAL_TEAM_2   = 1
	ARENA_SPECIAL_TEAM_3   = 2
	ARENA_SPECIAL_TEAM_MAX = 3
)

type ArenaSpecialEnemy struct {
	Class     int                                   `json:"class"`
	Dan       int                                   `json:"dan"`
	EnemyTeam [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo `json:"enemyteam"`
}

// 高阶竞技场玩家数据
type San_ArenaSpecialPlayer struct {
	Uid        int64
	Class      int
	Dan        int
	Coin       int
	Point      int64
	Name       string // 姓名
	Format     string // 阵容
	Count      int    // 今天挑战次数
	BuyCount   int    // 今天购买次数
	ArenaFight string // 战报
	State      int    // 状态

	format     [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo // 只保存防守阵容
	arenaFight []*ArenaSpecialFight                  // 战报
	enemy      []*ArenaSpecialEnemy                  // 敌人
	redPoint   RedPoint                              // 红点提醒

	DataUpdate
}

// 竞技场战报
type ArenaSpecialFightList struct {
	FightId [ARENA_SPECIAL_TEAM_MAX]int64         `json:"fight_id"` // 战斗Id
	Random  int64                                 `json:"random"`
	Time    int64                                 `json:"time"`   // 发生的时间
	Attack  [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo `json:"attack"` // 攻击者
	Defend  [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo `json:"defend"` // 防御者
	Result  [ARENA_SPECIAL_TEAM_MAX]BattleInfo    `json:"result"`
}

// 机器人缓存
type ArenaSpecialRobot struct {
	Class int
	Dan   int
	Info  []*JS_FightInfo
}

// 将数据库数据写入data
func (self *San_ArenaSpecialPlayer) Decode() {
	json.Unmarshal([]byte(self.Format), &self.format)
	json.Unmarshal([]byte(self.ArenaFight), &self.arenaFight)

}

// 将data数据写入数据库
func (self *San_ArenaSpecialPlayer) Encode() {
	self.Format = HF_JtoA(&self.format)
	self.ArenaFight = HF_JtoA(&self.arenaFight)
}

type San_ArenaSpecialRank struct {
	Uid       int64
	Class     int
	Dan       int
	Rank      int
	Point     int64
	StartTime int64

	DataUpdate
}

type ArenaSpecialMgr struct {
	Sql_Uid           map[int64]*San_ArenaSpecialPlayer       // 具体数据
	Sql_Rank          map[int]map[int][]*San_ArenaSpecialRank // 排序用数据
	Lock              *sync.RWMutex                           // 数据操作锁
	migrateOK         bool
	Robot             map[int]map[int][]ArenaSpecialRobot // 机器人数据
	ArenaTime         San_ArenaTime
	FightList         []*ArenaSpecialFightList
	FightListForCross []*ArenaSpecialFightList
}

var arenaSpecialMgr *ArenaSpecialMgr = nil

func GetArenaSpecialMgr() *ArenaSpecialMgr {
	if arenaSpecialMgr == nil {
		arenaSpecialMgr = new(ArenaSpecialMgr)
		arenaSpecialMgr.Sql_Uid = make(map[int64]*San_ArenaSpecialPlayer)
		arenaSpecialMgr.Sql_Rank = make(map[int]map[int][]*San_ArenaSpecialRank)
		arenaSpecialMgr.Lock = new(sync.RWMutex)
		arenaSpecialMgr.migrateOK = false
	}
	return arenaSpecialMgr
}

// 开启迁移数据协程，竞技场
func (self *ArenaSpecialMgr) RunMigrateArenaSpecial() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	infoName := "san_arenaspecialbattleinfo"     //! info 表
	recordName := "san_arenaspecialbattlerecord" //! record 表
	recordType := 5
	tableName := "san_battlerecord"

	len, err := GetRedisMgr().HLen(infoName)
	if err != nil {
		return
	}
	LogInfo("迁移数据：", infoName, len)

	migOK, err := GetRedisMgr().Exists(infoName + "_migrateOKNew")
	if err == nil && migOK == true {
		self.migrateOK = true
	}
	count := 0
	cursor := int64(0)
	for {
		if self.migrateOK == true {
			break
		}

		//! 迁移数据
		cursor1, num := MigrateDataOne(infoName, recordName, tableName, recordType, cursor)
		count += num
		cursor = cursor1

		if count >= len {
			GetRedisMgr().Set(infoName+"_migrateOKNew", "1")
			break
		}

		//! 延迟1ms
		time.Sleep(time.Millisecond)
	}

	LogInfo(infoName, "迁移数据OK")

}

// 更新阵容
func (self *ArenaSpecialMgr) UpdateFormat(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	uid := player.GetUid()
	_, ok := self.Sql_Uid[uid]
	if !ok {
		return
	}
	self.Sql_Uid[uid].format = [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo{}
	self.Sql_Uid[uid].format[ARENA_SPECIAL_TEAM_1] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_SPECIAL_4)
	self.Sql_Uid[uid].format[ARENA_SPECIAL_TEAM_2] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_SPECIAL_5)
	self.Sql_Uid[uid].format[ARENA_SPECIAL_TEAM_3] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_SPECIAL_6)
	self.Sql_Uid[uid].Format = HF_JtoA(&self.Sql_Uid[uid].format)

	self.Sql_Uid[uid].Update(true)

	GetTopArenaMgr().UpdateFight(ARENA_TOP_TYPE_SPECIAL_RANK, player)
	GetTopArenaMgr().UpdateFight(ARENA_TOP_TYPE_SPECIAL_POINT, player)
}

func (self *ArenaSpecialMgr) GetData() {
	arenaSpecialMgr.Sql_Uid = make(map[int64]*San_ArenaSpecialPlayer)
	var info2 San_ArenaSpecialPlayer
	tableName2 := fmt.Sprintf("san_playerarena%d", ARENA_TYPE_SPECIAL+1)
	sql2 := fmt.Sprintf("select * from `%s`", tableName2)
	res2 := GetServer().DBUser.GetAllData(sql2, &info2)
	for i := 0; i < len(res2); i++ {
		data := res2[i].(*San_ArenaSpecialPlayer)
		data.State = 0
		data.Init(tableName2, data, true)
		data.Decode()
		self.Sql_Uid[data.Uid] = data
	}

	arenaSpecialMgr.Sql_Rank = make(map[int]map[int][]*San_ArenaSpecialRank)
	var info1 San_ArenaSpecialRank
	tableName1 := fmt.Sprintf("san_rankarena%d", ARENA_TYPE_SPECIAL+1)
	sql1 := fmt.Sprintf("select * from `%s`", tableName1)
	res1 := GetServer().DBUser.GetAllData(sql1, &info1)

	for i := 0; i < len(res1); i++ {
		data := res1[i].(*San_ArenaSpecialRank)
		data.Init(tableName1, data, true)
		class := data.Class
		dan := data.Dan
		_, ok1 := arenaSpecialMgr.Sql_Rank[class]
		if !ok1 {
			arenaSpecialMgr.Sql_Rank[class] = make(map[int][]*San_ArenaSpecialRank)
		}

		_, ok2 := arenaSpecialMgr.Sql_Rank[class][dan]
		if !ok2 {
			arenaSpecialMgr.Sql_Rank[class][dan] = []*San_ArenaSpecialRank{}
		}

		arenaSpecialMgr.Sql_Rank[class][dan] = append(arenaSpecialMgr.Sql_Rank[class][dan], data)
		player, ok := self.Sql_Uid[data.Uid]
		if ok {
			player.Class = class
			player.Dan = dan
			player.Update(true)
		}
		data.Update(true)
	}

	now := TimeServer().Unix()
	for _, value := range GetCsvMgr().ArenaSpecialClassMap {
		for _, config := range value {
			class := config.Class
			dan := config.Dan
			_, ok1 := arenaSpecialMgr.Sql_Rank[class]
			if !ok1 {
				arenaSpecialMgr.Sql_Rank[class] = make(map[int][]*San_ArenaSpecialRank)
			}

			_, ok2 := arenaSpecialMgr.Sql_Rank[class][dan]
			if !ok2 {
				arenaSpecialMgr.Sql_Rank[class][dan] = []*San_ArenaSpecialRank{}
			}

			nLen := len(arenaSpecialMgr.Sql_Rank[class][dan])
			if nLen < config.Capacity {
				for i := 0; i < config.Capacity-nLen; i++ {
					data := &San_ArenaSpecialRank{}
					data.Uid = 0
					data.Class = class
					data.Dan = dan
					data.Rank = config.Ranking
					data.StartTime = now
					arenaSpecialMgr.Sql_Rank[class][dan] = append(arenaSpecialMgr.Sql_Rank[class][dan], data)
				}
			}
		}
	}

	for _, value := range self.Sql_Uid {
		if value.Point <= 0 {
			continue
		}
		fight := GetOfflineInfoMgr().GetTeamFight(value.Uid, TEAMTYPE_ARENA_SPECIAL_4)
		if fight <= 0 {
			fight = 0
			for _, format := range value.format {
				if format != nil {
					fight += format.Deffight
				}
			}
			GetOfflineInfoMgr().SetArenaFight(TEAMTYPE_ARENA_SPECIAL_4, value.Uid, fight)
		}
	}

	sql := fmt.Sprintf("select * from `san_arenatime` where type = %d", ARENA_TYPE_SPECIAL+1)
	res := GetServer().DBUser.GetAllData(sql, &self.ArenaTime)
	if len(res) > 0 {
		self.ArenaTime = *res[0].(*San_ArenaTime)
	}

	if self.ArenaTime.Type <= 0 {
		self.ArenaTime.Type = ARENA_TYPE_SPECIAL + 1
		startTime, _ := GetCsvMgr().GetNowStartAndEnd(ARENA_TIME_TYPE_SPECIAL)
		self.ArenaTime.StartTime = startTime
		now := TimeServer()
		timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
		if TimeServer().Hour() < 5 {
			timeStamp -= DAY_SECS
		}
		self.ArenaTime.RefreshTime = timeStamp
		InsertTable("san_arenatime", &self.ArenaTime, 0, true)
	}

	self.ArenaTime.Init("san_arenatime", &self.ArenaTime, true)

	self.InitRobot()

	self.CheckArenaEnd()
}

func (self *ArenaSpecialMgr) InitRobot() {
	self.Robot = make(map[int]map[int][]ArenaSpecialRobot)
	nLen := len(GetCsvMgr().JJCRobotConfig)
	for i := 0; i < nLen; i++ {
		cfg := GetCsvMgr().JJCRobotConfig[i]
		if cfg.Type != 2 {
			continue
		}
		if cfg.Teamnum != 1 {
			continue
		}

		if i+1 >= nLen || i+2 >= nLen {
			continue
		}

		name := GetCsvMgr().GetName()
		var node ArenaSpecialRobot
		node.Class = cfg.Jjcclass
		node.Dan = cfg.Jjcdan
		node.Info = append(node.Info, GetCsvMgr().GetRobot(cfg))
		node.Info = append(node.Info, GetCsvMgr().GetRobot(GetCsvMgr().JJCRobotConfig[i+1]))
		node.Info = append(node.Info, GetCsvMgr().GetRobot(GetCsvMgr().JJCRobotConfig[i+2]))
		node.Info[0].Uname = name

		for _, z := range node.Info {
			if z.Uname != node.Info[0].Uname {
				z.Uname = node.Info[0].Uname
			}
		}
		_, ok1 := self.Robot[cfg.Jjcclass]
		if !ok1 {
			self.Robot[cfg.Jjcclass] = make(map[int][]ArenaSpecialRobot)
		}
		_, ok2 := self.Robot[cfg.Jjcclass][cfg.Jjcdan]
		if !ok2 {
			self.Robot[cfg.Jjcclass][cfg.Jjcdan] = make([]ArenaSpecialRobot, 0)
		}
		self.Robot[cfg.Jjcclass][cfg.Jjcdan] = append(self.Robot[cfg.Jjcclass][cfg.Jjcdan], node)
	}
}

func (self *ArenaSpecialMgr) Save() {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	for _, value := range self.Sql_Uid {
		value.Encode()
		value.Update(true)
	}

	for _, value := range self.Sql_Rank {
		for _, rank := range value {
			for _, t := range rank {
				t.Update(true)
			}
		}
	}

	self.ArenaTime.Update(true)
}

func (self *ArenaSpecialMgr) AddPlayer(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	uid := player.GetUid()
	_, ok1 := self.Sql_Uid[uid]
	if !ok1 {
		self.Sql_Uid[uid] = &San_ArenaSpecialPlayer{}
		self.Sql_Uid[uid].Uid = uid
		self.Sql_Uid[uid].Class = ARENA_SPECIAL_BASE_CLASS
		self.Sql_Uid[uid].Dan = ARENA_SPECIAL_BASE_DAN
		self.Sql_Uid[uid].Name = player.GetName()

		self.Sql_Uid[uid].format = [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo{}
		self.Sql_Uid[uid].format[ARENA_SPECIAL_TEAM_1] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_SPECIAL_4)
		self.Sql_Uid[uid].format[ARENA_SPECIAL_TEAM_2] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_SPECIAL_5)
		self.Sql_Uid[uid].format[ARENA_SPECIAL_TEAM_3] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_SPECIAL_6)
		self.Sql_Uid[uid].Format = HF_JtoA(&self.Sql_Uid[uid].format)
		self.Sql_Uid[uid].Count = 0

		_, ok2 := self.Sql_Rank[ARENA_SPECIAL_BASE_CLASS]
		if !ok2 {
			self.Sql_Rank[ARENA_SPECIAL_BASE_CLASS] = make(map[int][]*San_ArenaSpecialRank)
		}
		_, ok3 := self.Sql_Rank[ARENA_SPECIAL_BASE_CLASS][ARENA_SPECIAL_BASE_DAN]
		if !ok3 {
			self.Sql_Rank[ARENA_SPECIAL_BASE_CLASS][ARENA_SPECIAL_BASE_DAN] = make([]*San_ArenaSpecialRank, 0)
		}

		data := San_ArenaSpecialRank{}
		data.Uid = uid
		data.Rank = 0
		data.Point = 0
		data.Class = ARENA_SPECIAL_BASE_CLASS
		data.Dan = ARENA_SPECIAL_BASE_DAN
		data.StartTime = TimeServer().Unix()

		self.Sql_Rank[ARENA_SPECIAL_BASE_CLASS][ARENA_SPECIAL_BASE_DAN] = append(self.Sql_Rank[ARENA_SPECIAL_BASE_CLASS][ARENA_SPECIAL_BASE_DAN], &data)

		tableName1 := fmt.Sprintf("san_rankarena%d", ARENA_TYPE_SPECIAL+1)
		tableName2 := fmt.Sprintf("san_playerarena%d", ARENA_TYPE_SPECIAL+1)

		InsertTable(tableName1, &data, 0, true)
		data.Init(tableName1, &data, true)
		InsertTable(tableName2, self.Sql_Uid[uid], 0, true)
		self.Sql_Uid[uid].Init(tableName2, self.Sql_Uid[uid], true)
	} else {
		nClass := self.Sql_Uid[uid].Class
		nDan := self.Sql_Uid[uid].Dan
		_, ok2 := self.Sql_Rank[nClass]
		if !ok2 {
			self.Sql_Rank[nClass] = make(map[int][]*San_ArenaSpecialRank)
		}
		_, ok3 := self.Sql_Rank[nClass][nDan]
		if !ok3 {
			self.Sql_Rank[nClass][nDan] = make([]*San_ArenaSpecialRank, 0)
		}

		find := false
		for _, t := range self.Sql_Rank[nClass][nDan] {
			if t.Uid == uid {
				find = true
				break
			}
		}

		if !find {
			LogDebug("ArenaSpecialMgr ERRER !!!!!!!!!! self.Sql_Rank[nClass][nDan] !find")

			findrank := false
			for _, p := range arenaSpecialMgr.Sql_Rank {
				for _, q := range p {
					for _, r := range q {
						if r.Uid == 0 {
							continue
						}

						if r.Uid == uid {
							self.Sql_Uid[uid].Class = r.Class
							self.Sql_Uid[uid].Dan = r.Dan
							findrank = true
							break
						}
					}
					if findrank {
						break
					}
				}
				if findrank {
					break
				}
			}
			if !findrank {
				LogDebug("ArenaSpecialMgr ERRER !!!!!!!!!!  !findrank")
				data := San_ArenaSpecialRank{}
				data.Uid = uid
				data.Rank = 0
				data.Point = 0
				data.Class = nClass
				data.Dan = nDan
				data.StartTime = TimeServer().Unix()

				self.Sql_Rank[nClass][nDan] = append(self.Sql_Rank[nClass][nDan], &data)
			}
		}
	}
}

func (self *ArenaSpecialMgr) GetRobot(class, dan int) [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo {
	var robot [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo
	_, ok1 := self.Robot[class]
	if !ok1 {
		return robot
	}
	value, ok2 := self.Robot[class][dan]
	if !ok2 {
		return robot
	}

	nLen := len(value)
	if nLen <= 0 {
		return robot
	}
	temp := HF_GetRandom(nLen)
	for i := 0; i < ARENA_SPECIAL_TEAM_MAX; i++ {
		tempdata := JS_FightInfo{}
		HF_DeepCopy(&tempdata, self.Robot[class][dan][temp].Info[i])
		robot[i] = &tempdata
	}
	return robot
}

func (self *ArenaSpecialMgr) GetPlayerData(uid int64) *San_ArenaSpecialPlayer {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.Sql_Uid[uid]
	if ok {
		return self.Sql_Uid[uid]
	}
	return nil
}
func (self *ArenaSpecialMgr) GetRankData(uid int64) *San_ArenaSpecialRank {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.Sql_Uid[uid]
	if !ok {
		return nil
	}

	nClass := self.Sql_Uid[uid].Class
	nDan := self.Sql_Uid[uid].Dan
	_, ok1 := self.Sql_Rank[nClass]
	if !ok1 {
		return nil
	}
	_, ok2 := self.Sql_Rank[nClass][nDan]
	if !ok2 {
		return nil
	}

	for _, t := range self.Sql_Rank[nClass][nDan] {
		if t.Uid == uid {
			return t
		}
	}
	return nil
}

func (self *ArenaSpecialMgr) RandomEnemy(nRange []int) ([]*ArenaSpecialEnemy, []int) {
	fightinfo := make([]*ArenaSpecialEnemy, 0)
	idRet := []int{}
	for _, id := range nRange {
		if id <= 0 {
			continue
		}

		pIDConfig := GetCsvMgr().GetArenaSpecialClassConfigByID(id)
		if pIDConfig == nil {
			continue
		}

		nIDClass := pIDConfig.Class
		nIDDan := pIDConfig.Dan
		_, ok1 := self.Sql_Rank[nIDClass]
		if !ok1 {
			continue
		}
		_, ok2 := self.Sql_Rank[nIDClass][nIDDan]
		if !ok2 {
			continue
		}

		nIDLen := len(self.Sql_Rank[nIDClass][nIDDan])
		if nIDLen <= 0 {
			continue
		}

		if nIDLen == 1 {
			uid := self.Sql_Rank[nIDClass][nIDDan][0].Uid
			if uid != 0 {
				value, ok3 := self.Sql_Uid[uid]
				if ok3 {
					fightinfo = append(fightinfo, &ArenaSpecialEnemy{nIDClass, nIDDan, value.format})
					idRet = append(idRet, id)
					continue
				}
			}
			fightinfo = append(fightinfo, &ArenaSpecialEnemy{nIDClass, nIDDan, self.GetRobot(nIDClass, nIDDan)})
			idRet = append(idRet, id)
			continue
		}

		temp := HF_GetRandom(nIDLen)
		uid := self.Sql_Rank[nIDClass][nIDDan][temp].Uid
		if uid != 0 {
			value, ok3 := self.Sql_Uid[uid]
			if ok3 {
				fightinfo = append(fightinfo, &ArenaSpecialEnemy{nIDClass, nIDDan, value.format})
				idRet = append(idRet, id)
				continue
			}
		}
		fightinfo = append(fightinfo, &ArenaSpecialEnemy{nIDClass, nIDDan, self.GetRobot(nIDClass, nIDDan)})
		idRet = append(idRet, id)
		continue
	}

	return fightinfo, idRet
}

func (self *ArenaSpecialMgr) GetEnemy(uid int64) ([]*ArenaSpecialEnemy, []int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	data, ok := self.Sql_Uid[uid]
	if !ok {
		return nil, nil
	}

	if data.Class < 0 {
		return nil, nil
	}

	if data.Dan < 0 {
		return nil, nil
	}

	class := data.Class
	dan := data.Dan

	config := GetCsvMgr().GetArenaSpecialClassConfig(class, dan)
	if config == nil {
		return nil, nil
	}

	param1 := config.Id - 1
	param2 := config.Id - (2 + HF_GetRandom(3))
	param3 := config.Id - (5 + HF_GetRandom(5))

	if param1 <= 0 {
		param1 = config.Id + 1
	}
	if param2 <= 0 {
		param2 = config.Id + (2 + HF_GetRandom(3))
	}
	if param3 <= 0 {
		param2 = config.Id + (5 + HF_GetRandom(5))
	}

	enemyID := []int{param1, param2, param3}
	fightinfo, idRet := self.RandomEnemy(enemyID)

	rankFight := make([]*ArenaSpecialEnemy, 0)
	nLen := len(fightinfo)
	if nLen > 0 {
		for i := nLen - 1; i >= 0; i-- {
			rankFight = append(rankFight, fightinfo[i])
		}
	} else {
		enemyID = []int{config.Id + 1, config.Id + (2 + HF_GetRandom(3)), config.Id + (5 + HF_GetRandom(5))}
		fightinfo, idRet = self.RandomEnemy(enemyID)
		nLen = len(fightinfo)
		for i := nLen - 1; i >= 0; i-- {
			rankFight = append(rankFight, fightinfo[i])
		}
	}

	data.enemy = rankFight
	return rankFight, idRet
}

func (self *ArenaSpecialMgr) AddFightList(player *Player, attack [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo, defend [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo, time int64, random int64) [ARENA_SPECIAL_TEAM_MAX]int64 {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	fightID := [ARENA_SPECIAL_TEAM_MAX]int64{}
	if len(self.FightList) >= ARENA_FIGHT_COUNT_MAX {
		return fightID
	}

	data, ok := self.Sql_Uid[player.GetUid()]
	if !ok {
		player.SendErr(GetCsvMgr().GetText("数据丢失"))
		return fightID
	}

	if data.State == 1 {
		player.SendErr(GetCsvMgr().GetText("有玩家正在攻击你"))
		return fightID
	}

	uid := defend[0].Uid
	if uid != 0 {
		_, ok := self.Sql_Uid[uid]
		if ok {
			if self.Sql_Uid[uid].State == 1 {
				player.SendErr(GetCsvMgr().GetText("敌人战斗中"))
				return fightID
			} else {
				self.Sql_Uid[uid].State = 1
			}
		}
	}
	data.State = 1
	// 打的次数+1
	data.Count++

	for i, _ := range fightID {
		ret := GetFightMgr().AddArenaFightID(attack[i], defend[i], random)
		fightID[i] = ret.Id
	}
	fightlist := &ArenaSpecialFightList{fightID, random, time, attack, defend, [ARENA_SPECIAL_TEAM_MAX]BattleInfo{}}
	self.FightList = append(self.FightList, fightlist)

	return fightID
}

func (self *ArenaSpecialMgr) AddFightListForCross(player *Player, attack [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo, defend [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo, time int64, random int64) [CROSSARENA3V3_TEAM_MAX]int64 {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	fightID := [ARENA_SPECIAL_TEAM_MAX]int64{}
	if len(self.FightListForCross) >= ARENA_FIGHT_COUNT_MAX {
		return fightID
	}

	for i, _ := range fightID {
		ret := GetFightMgr().AddArenaFightID(attack[i], defend[i], random)
		fightID[i] = ret.Id
	}
	fightlist := &ArenaSpecialFightList{fightID, random, time, attack, defend, [CROSSARENA3V3_TEAM_MAX]BattleInfo{}}
	self.FightListForCross = append(self.FightListForCross, fightlist)

	return fightID
}

func (self *ArenaSpecialMgr) CheckArenaEnd() {
	now := TimeServer()
	timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
	if TimeServer().Hour() < 5 {
		timeStamp -= DAY_SECS
	}

	if self.ArenaTime.RefreshTime != timeStamp {
		for _, p := range self.Sql_Uid {
			p.Count = 0
			p.BuyCount = 0
		}

		self.ArenaTime.RefreshTime = timeStamp
	}
}

// 十分钟结算一次
func (self *ArenaSpecialMgr) CountAward() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			self.OnCountAward()
		}
	}
	ticker.Stop()
}

// 定时器
func (self *ArenaSpecialMgr) OnCountAward() {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	for _, class := range self.Sql_Rank {
		for _, dan := range class {
			for _, rank := range dan {
				data, ok := self.Sql_Uid[rank.Uid]
				if ok {
					data.CheckArenaAward(rank)
				}
			}
		}
	}
}

// 开启战斗协程
func (self *ArenaSpecialMgr) StartFight() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			self.OnTimer()
		}
	}
	ticker.Stop()
}

// 定时器
func (self *ArenaSpecialMgr) OnTimer() {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	if len(self.Sql_Uid) <= 0 {
		return
	}

	now := TimeServer().Unix()

	nLen := len(self.FightList)
	for i := nLen - 1; i >= 0; i-- {
		data := self.FightList[i]
		if data.Time+300 <= now {
			player := GetPlayerMgr().GetPlayer(data.Attack[0].Uid, false)
			if player == nil {
				continue
			}
			battleInfoAdd := [ARENA_SPECIAL_TEAM_MAX]BattleInfo{}
			for t := ARENA_SPECIAL_TEAM_1; t < ARENA_SPECIAL_TEAM_MAX; t++ {
				battleInfo := BattleInfo{}
				battleInfo.Id = data.FightId[t]
				var attack []*BattleHeroInfo
				var defend []*BattleHeroInfo
				for _, v := range data.Attack[t].Heroinfo {
					attack = append(attack, &BattleHeroInfo{v.Heroid, v.Levels, v.Stars, v.Skin, 123, 456, 789, 11, 12, nil, 0, nil})
				}
				for _, v := range data.Defend[t].Heroinfo {
					defend = append(defend, &BattleHeroInfo{v.Heroid, v.Levels, v.Stars, v.Skin, 987, 654, 321, 12, 45, nil, 0, nil})
				}

				battleInfo.UserInfo[POS_ATTACK] = &BattleUserInfo{
					data.Attack[t].Uid,
					data.Attack[t].Uname,
					data.Attack[t].Iconid,
					data.Attack[t].Portrait,
					data.Attack[t].UnionName,
					data.Attack[t].Level,
					attack}
				battleInfo.UserInfo[POS_DEFENCE] = &BattleUserInfo{
					data.Defend[t].Uid,
					data.Defend[t].Uname,
					data.Defend[t].Iconid,
					data.Defend[t].Portrait,
					data.Defend[t].UnionName,
					data.Defend[t].Level,
					defend}
				battleInfo.Type = BATTLE_TYPE_PVP
				battleInfo.Time = data.Time
				battleInfo.Random = data.Random
				if data.Attack[t].Deffight > data.Defend[t].Deffight {
					battleInfo.Result = 0
				} else {
					battleInfo.Result = 1
				}

				battleInfoAdd[t] = battleInfo
			}
			data.Result = battleInfoAdd
			result := 0
			for _, p := range data.Result {
				if p.Result == 0 {
					result++
				}
			}
			if result >= 2 {
				self.FightEnd(player, 1, data.Result)
			} else {
				self.FightEnd(player, 2, data.Result)
			}

			for _, id := range data.FightId {
				GetFightMgr().DelResult(id)
			}
			self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
		}

		// 尝试获取战斗结果
		//LogDebug("尝试获取战斗结果")
		for t := ARENA_SPECIAL_TEAM_1; t < ARENA_SPECIAL_TEAM_MAX; t++ {
			FightResult := GetFightMgr().GetResult(data.FightId[t])
			// 有结果,设置战斗时间
			if FightResult == nil {
				continue
			}

			player := GetPlayerMgr().GetPlayer(data.Attack[t].Uid, false)
			if player == nil {
				continue
			}

			battleInfo := BattleInfo{}
			battleInfo.Id = data.FightId[t]
			attackHeroInfo := []*BattleHeroInfo{}
			for i, v := range FightResult.Info[POS_ATTACK] {
				level, star, skin, exclusivelv := 0, 0, 0, 0
				if i < len(FightResult.Fight[POS_ATTACK].Heroinfo) {
					level = FightResult.Fight[POS_ATTACK].Heroinfo[i].Levels
					star = FightResult.Fight[POS_ATTACK].Heroinfo[i].Stars
					skin = FightResult.Fight[POS_ATTACK].Heroinfo[i].Skin
					exclusivelv = FightResult.Fight[POS_ATTACK].Heroinfo[i].HeroExclusiveLv
				}
				attackHeroInfo = append(attackHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusivelv, nil})
			}
			defendHeroInfo := []*BattleHeroInfo{}
			for i, v := range FightResult.Info[POS_DEFENCE] {
				level, star, skin, exclusivelv := 0, 0, 0, 0
				if i < len(FightResult.Fight[POS_DEFENCE].Heroinfo) {
					level = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Levels
					star = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Stars
					skin = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Skin
					exclusivelv = FightResult.Fight[POS_DEFENCE].Heroinfo[i].HeroExclusiveLv
				}
				defendHeroInfo = append(defendHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusivelv, nil})
			}
			battleInfo.UserInfo[POS_ATTACK] = &BattleUserInfo{
				data.Attack[t].Uid,
				data.Attack[t].Uname,
				data.Attack[t].Iconid,
				data.Attack[t].Portrait,
				data.Attack[t].UnionName,
				data.Attack[t].Level,
				attackHeroInfo}
			battleInfo.UserInfo[POS_DEFENCE] = &BattleUserInfo{
				data.Defend[t].Uid,
				data.Defend[t].Uname,
				data.Defend[t].Iconid,
				data.Defend[t].Portrait,
				data.Defend[t].UnionName,
				data.Defend[t].Level,
				defendHeroInfo}
			battleInfo.Type = BATTLE_TYPE_PVP
			battleInfo.Time = data.Time
			battleInfo.Random = data.Random
			if FightResult.Result == 1 {
				battleInfo.Result = 0
			} else {
				battleInfo.Result = 1
			}

			data.Result[t] = battleInfo

			over := true
			for _, q := range data.Result {
				if q.Id <= 0 {
					over = false
				}
			}
			if over {
				result := 0
				for _, p := range data.Result {
					if p.Result == 0 {
						result++
					}
				}
				if result >= 2 {
					self.FightEnd(player, 1, data.Result)
				} else {
					self.FightEnd(player, 2, data.Result)
				}
				for _, id := range data.FightId {
					GetFightMgr().DelResult(id)
				}
				self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
				continue
			}
		}
	}


	self.ForCrossFightRun()
}

func (self *ArenaSpecialMgr) ForCrossFightRun() {
	now := TimeServer().Unix()
	nLen := len(self.FightListForCross)
	for i := nLen - 1; i >= 0; i-- {
		data := self.FightListForCross[i]
		if data.Time+300 <= now {
			player := GetPlayerMgr().GetPlayer(data.Attack[0].Uid, false)
			if player == nil {
				continue
			}
			battleInfoAdd := [CROSSARENA3V3_TEAM_MAX]BattleInfo{}
			for t := 0; t < 3; t++ {
				battleInfo := BattleInfo{}
				battleInfo.Id = data.FightId[t]
				var attack []*BattleHeroInfo
				var defend []*BattleHeroInfo
				for _, v := range data.Attack[t].Heroinfo {
					attack = append(attack, &BattleHeroInfo{v.Heroid, v.Levels, v.Stars, v.Skin, 123, 456, 789, 11, 12, nil, 0, nil})
				}
				for _, v := range data.Defend[t].Heroinfo {
					defend = append(defend, &BattleHeroInfo{v.Heroid, v.Levels, v.Stars, v.Skin, 987, 654, 321, 12, 45, nil, 0, nil})
				}

				battleInfo.UserInfo[POS_ATTACK] = &BattleUserInfo{
					data.Attack[t].Uid,
					data.Attack[t].Uname,
					data.Attack[t].Iconid,
					data.Attack[t].Portrait,
					data.Attack[t].UnionName,
					data.Attack[t].Level,
					attack}
				battleInfo.UserInfo[POS_DEFENCE] = &BattleUserInfo{
					data.Defend[t].Uid,
					data.Defend[t].Uname,
					data.Defend[t].Iconid,
					data.Defend[t].Portrait,
					data.Defend[t].UnionName,
					data.Defend[t].Level,
					defend}
				battleInfo.Type = BATTLE_TYPE_PVP
				battleInfo.Time = data.Time
				battleInfo.Random = data.Random
				if data.Attack[t].Deffight > data.Defend[t].Deffight {
					battleInfo.Result = 0
				} else {
					battleInfo.Result = 1
				}

				battleInfoAdd[t] = battleInfo
			}
			data.Result = battleInfoAdd
			result := 0
			for _, p := range data.Result {
				if p.Result == 0 {
					result++
				}
			}
			if result >= 2 {
				player.GetModule("crossarena3v3").(*ModCrossArena3V3).FightEndOK(player, 1, data)
			} else {
				player.GetModule("crossarena3v3").(*ModCrossArena3V3).FightEndOK(player, 2, data)
			}

			for _, id := range data.FightId {
				GetFightMgr().DelResult(id)
			}
			self.FightListForCross = append(self.FightListForCross[0:i], self.FightListForCross[i+1:]...)
		}

		// 尝试获取战斗结果
		//LogDebug("尝试获取战斗结果")
		for t := ARENA_SPECIAL_TEAM_1; t < ARENA_SPECIAL_TEAM_MAX; t++ {
			FightResult := GetFightMgr().GetResult(data.FightId[t])
			// 有结果,设置战斗时间
			if FightResult == nil {
				continue
			}

			player := GetPlayerMgr().GetPlayer(data.Attack[t].Uid, false)
			if player == nil {
				continue
			}

			battleInfo := BattleInfo{}
			battleInfo.Id = data.FightId[t]
			attackHeroInfo := []*BattleHeroInfo{}
			for i, v := range FightResult.Info[POS_ATTACK] {
				level, star, skin, exclusivelv := 0, 0, 0, 0
				if i < len(FightResult.Fight[POS_ATTACK].Heroinfo) {
					level = FightResult.Fight[POS_ATTACK].Heroinfo[i].Levels
					star = FightResult.Fight[POS_ATTACK].Heroinfo[i].Stars
					skin = FightResult.Fight[POS_ATTACK].Heroinfo[i].Skin
					exclusivelv = FightResult.Fight[POS_ATTACK].Heroinfo[i].HeroExclusiveLv
				}
				attackHeroInfo = append(attackHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusivelv, nil})
			}
			defendHeroInfo := []*BattleHeroInfo{}
			for i, v := range FightResult.Info[POS_DEFENCE] {
				level, star, skin, exclusivelv := 0, 0, 0, 0
				if i < len(FightResult.Fight[POS_DEFENCE].Heroinfo) {
					level = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Levels
					star = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Stars
					skin = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Skin
					exclusivelv = FightResult.Fight[POS_DEFENCE].Heroinfo[i].HeroExclusiveLv
				}
				defendHeroInfo = append(defendHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusivelv, nil})
			}
			battleInfo.UserInfo[POS_ATTACK] = &BattleUserInfo{
				data.Attack[t].Uid,
				data.Attack[t].Uname,
				data.Attack[t].Iconid,
				data.Attack[t].Portrait,
				data.Attack[t].UnionName,
				data.Attack[t].Level,
				attackHeroInfo}
			battleInfo.UserInfo[POS_DEFENCE] = &BattleUserInfo{
				data.Defend[t].Uid,
				data.Defend[t].Uname,
				data.Defend[t].Iconid,
				data.Defend[t].Portrait,
				data.Defend[t].UnionName,
				data.Defend[t].Level,
				defendHeroInfo}
			battleInfo.Type = BATTLE_TYPE_PVP
			battleInfo.Time = data.Time
			battleInfo.Random = data.Random
			if FightResult.Result == 1 {
				battleInfo.Result = 0
			} else {
				battleInfo.Result = 1
			}

			data.Result[t] = battleInfo

			over := true
			for _, q := range data.Result {
				if q.Id <= 0 {
					over = false
				}
			}
			if over {
				result := 0
				for _, p := range data.Result {
					if p.Result == 0 {
						result++
					}
				}
				if result >= 2 {
					player.GetModule("crossarena3v3").(*ModCrossArena3V3).FightEndOK(player, 1, data)
				} else {
					player.GetModule("crossarena3v3").(*ModCrossArena3V3).FightEndOK(player, 2, data)
				}
				for _, id := range data.FightId {
					GetFightMgr().DelResult(id)
				}
				self.FightListForCross = append(self.FightListForCross[0:i], self.FightListForCross[i+1:]...)
				continue
			}
		}
	}

}

func (self *ArenaSpecialMgr) ArenaSpecialFightResult(battleInfo [ARENA_SPECIAL_TEAM_MAX]BattleInfo) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	if len(self.Sql_Uid) <= 0 {
		return
	}

	for i, v := range self.FightList {
		data := v

		if battleInfo[0].Id != data.FightId[0] {
			continue
		}

		for p, t := range data.FightId {
			if battleInfo[p].Id != t {
				for _, z := range data.FightId {
					GetFightMgr().DelResult(z)
				}
				self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
				return
			}
		}

		player := GetPlayerMgr().GetPlayer(data.Attack[0].Uid, false)
		if player == nil {
			for _, z := range data.FightId {
				GetFightMgr().DelResult(z)
			}
			self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
			return
		}

		for d := 0; d < ARENA_SPECIAL_TEAM_MAX; d++ {
			battleInfo[d].UserInfo[POS_ATTACK].Uid = data.Attack[d].Uid
			battleInfo[d].UserInfo[POS_ATTACK].Name = data.Attack[d].Uname
			battleInfo[d].UserInfo[POS_ATTACK].Icon = data.Attack[d].Iconid
			battleInfo[d].UserInfo[POS_ATTACK].Portrait = data.Attack[d].Portrait
			battleInfo[d].UserInfo[POS_ATTACK].UnionName = data.Attack[d].UnionName
			battleInfo[d].UserInfo[POS_ATTACK].Level = data.Attack[d].Level

			battleInfo[d].UserInfo[POS_DEFENCE].Uid = data.Defend[d].Uid
			battleInfo[d].UserInfo[POS_DEFENCE].Name = data.Defend[d].Uname
			battleInfo[d].UserInfo[POS_DEFENCE].Icon = data.Defend[d].Iconid
			battleInfo[d].UserInfo[POS_DEFENCE].Portrait = data.Defend[d].Portrait
			battleInfo[d].UserInfo[POS_DEFENCE].UnionName = data.Defend[d].UnionName
			battleInfo[d].UserInfo[POS_DEFENCE].Level = data.Defend[d].Level

			battleInfo[d].Type = BATTLE_TYPE_PVP
			battleInfo[d].Time = data.Time
			battleInfo[d].Random = data.Random
		}

		result := 0
		for _, p := range battleInfo {
			if p.Result == 0 {
				result++
			}
		}
		if result >= 2 {
			self.FightEnd(player, 1, battleInfo)
		} else {
			self.FightEnd(player, 2, battleInfo)
		}

		for _, z := range data.FightId {
			GetFightMgr().DelResult(z)
		}
		self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
		return
	}
}
func (self *ArenaSpecialMgr) Swap(player *Player, myInfo *San_ArenaSpecialPlayer, enemy *ArenaSpecialEnemy) {
	myClass := myInfo.Class
	myDan := myInfo.Dan
	enemyClass := enemy.Class
	enemyDan := enemy.Dan
	config1 := GetCsvMgr().GetArenaSpecialClassConfig(myClass, myDan)
	if nil == config1 {
		return
	}

	config2 := GetCsvMgr().GetArenaSpecialClassConfig(enemy.Class, enemy.Dan)
	if nil == config2 {
		return
	}

	if config1.Id <= config2.Id {
		return
	}

	if enemy.EnemyTeam[0].Uid != 0 {
		enemyInfo, ok := self.Sql_Uid[enemy.EnemyTeam[0].Uid]
		if ok {
			for i, v := range self.Sql_Rank[enemyInfo.Class][enemyInfo.Dan] {
				if enemyInfo.Uid == v.Uid {
					self.Sql_Uid[enemy.EnemyTeam[0].Uid].CheckArenaAward(v)

					data := self.Sql_Rank[enemyInfo.Class][enemyInfo.Dan][i]
					self.Sql_Rank[enemyInfo.Class][enemyInfo.Dan] = append(self.Sql_Rank[enemyInfo.Class][enemyInfo.Dan][:i], self.Sql_Rank[enemyInfo.Class][enemyInfo.Dan][i+1:]...)
					self.Sql_Uid[enemy.EnemyTeam[0].Uid].Class = myClass
					self.Sql_Uid[enemy.EnemyTeam[0].Uid].Dan = myDan

					player := GetPlayerMgr().GetPlayer(enemy.EnemyTeam[0].Uid, false)
					if player != nil {
						player.HandleTask(TASK_TYPE_SPECIAL_ARENA_CLASS, myClass, myDan, 0)
					}

					data.Class = myClass
					data.Dan = myDan
					data.Rank = config1.Ranking
					data.StartTime = TimeServer().Unix()
					self.Sql_Rank[myClass][myDan] = append(self.Sql_Rank[myClass][myDan], data)
					break
				}
			}
		}
	} else {
		for i, v := range self.Sql_Rank[enemy.Class][enemy.Dan] {
			if 0 == v.Uid {
				self.Sql_Rank[enemy.Class][enemy.Dan] = append(self.Sql_Rank[enemy.Class][enemy.Dan][:i], self.Sql_Rank[enemy.Class][enemy.Dan][i+1:]...)
				newEnemy := &San_ArenaSpecialRank{}
				newEnemy.Uid = 0
				newEnemy.Class = myClass
				newEnemy.Dan = myDan
				newEnemy.Rank = config1.Ranking
				newEnemy.StartTime = TimeServer().Unix()
				self.Sql_Rank[myClass][myDan] = append(self.Sql_Rank[myClass][myDan], newEnemy)
				break
			}
		}
	}

	for i, v := range self.Sql_Rank[myClass][myDan] {
		if myInfo.Uid == v.Uid {
			self.Sql_Uid[v.Uid].CheckArenaAward(v)
			self.Sql_Uid[v.Uid].Class = enemyClass
			self.Sql_Uid[v.Uid].Dan = enemyDan
			player := GetPlayerMgr().GetPlayer(v.Uid, false)
			if player != nil {
				player.HandleTask(TASK_TYPE_SPECIAL_ARENA_CLASS, enemyClass, enemyDan, 0)
			}
			data := self.Sql_Rank[myClass][myDan][i]
			self.Sql_Rank[myClass][myDan] = append(self.Sql_Rank[myClass][myDan][:i], self.Sql_Rank[myClass][myDan][i+1:]...)
			data.Class = self.Sql_Uid[v.Uid].Class
			data.Dan = self.Sql_Uid[v.Uid].Dan
			data.Rank = config2.Ranking
			data.StartTime = TimeServer().Unix()
			self.Sql_Rank[self.Sql_Uid[v.Uid].Class][self.Sql_Uid[v.Uid].Dan] = append(self.Sql_Rank[self.Sql_Uid[v.Uid].Class][self.Sql_Uid[v.Uid].Dan], data)
			break
		}
	}

	if config1.Ranking > 0 || config2.Ranking > 0 {
		GetTopArenaMgr().UpdateArenaSpecialRank(int64(config2.Ranking), player, int64(config1.Ranking), enemy.EnemyTeam[0].Uid)
	}
}

// 斗技场结束
func (self *ArenaSpecialMgr) FightEnd(player *Player, result int, battleInfo [ARENA_SPECIAL_TEAM_MAX]BattleInfo) {
	modArena := player.GetModule("arenaspecial").(*ModArenaSpecial)
	if modArena.Index < 0 {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}

	backRet := 0

	uid := player.GetUid()
	info, ok := self.Sql_Uid[uid]
	if !ok {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}

	vipcsv := GetCsvMgr().GetVipConfig(player.Sql_UserBase.Vip)
	if vipcsv == nil {
		player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_ERROR"))
		return
	}

	msg := &S2C_ArenaSpecialStartFight{}
	info.State = 0
	// 敌人信息
	fight := info.enemy[modArena.Index]
	myFight := [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo{}
	myFight[ARENA_SPECIAL_TEAM_1] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_SPECIAL_1)
	myFight[ARENA_SPECIAL_TEAM_2] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_SPECIAL_2)
	myFight[ARENA_SPECIAL_TEAM_3] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_SPECIAL_3)
	// 打成功
	if result == 1 {
		info.AddFight(myFight, fight.EnemyTeam, 0, 1, battleInfo, fight.Class, fight.Dan, true)

		msg.MyClass = info.Class
		msg.MyDan = info.Dan
		msg.EnemyClass = fight.Class
		msg.EnemyDan = fight.Dan
		// 打的如果不是机器人
		if fight.EnemyTeam[0].Uid != 0 {
			// 获得敌人数据
			enemy, ok := self.Sql_Uid[fight.EnemyTeam[0].Uid]
			if !ok || enemy == nil {
				player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
				return
			}
			enemy.AddFight(myFight, fight.EnemyTeam, 0, 0, battleInfo, info.Class, info.Dan, false)
			enemy.State = 0
		}
		self.Swap(player, info, fight)
		player.HandleTask(TASK_TYPE_ARENA_COUNT, 1, 0, 0)
		player.HandleTask(TASK_TYPE_JJC_SCORE, 3, 2, 0)

		msg.FightInfo[1] = fight.EnemyTeam

		config := GetCsvMgr().GetArenaSpecialClassConfig(msg.MyClass, msg.MyDan)
		enemyid := 0
		if config != nil {
			enemyid = config.Id
		}

		config = GetCsvMgr().GetArenaSpecialClassConfig(msg.EnemyClass, msg.EnemyDan)
		myid := 0
		if config != nil {
			myid = config.Id
		}

		GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_ARENA_SPECIAL_FIGHT, 1, int(fight.EnemyTeam[0].Uid), enemyid, "高阶竞技场战斗", 0, myid, player)
		backRet = 0
	} else {
		msg.MyClass = info.Class
		msg.MyDan = info.Dan
		msg.EnemyClass = fight.Class
		msg.EnemyDan = fight.Dan
		info.AddFight(myFight, fight.EnemyTeam, 1, 1, battleInfo, fight.Class, fight.Dan, true)

		// 打的如果不是机器人
		if fight.EnemyTeam[0].Uid != 0 {
			// 获得敌人数据
			enemy, ok := self.Sql_Uid[fight.EnemyTeam[0].Uid]
			if !ok || enemy == nil {
				player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
				return
			}
			enemy.AddFight(myFight, fight.EnemyTeam, 1, 0, battleInfo, info.Class, info.Dan, false)
			enemy.State = 0
		}

		msg.FightInfo[1] = fight.EnemyTeam
		backRet = 1

		player.HandleTask(TASK_TYPE_JJC_SCORE, 1, 2, 0)

		config := GetCsvMgr().GetArenaSpecialClassConfig(msg.MyClass, msg.MyDan)
		enemyid := 0
		if config != nil {
			enemyid = config.Id
		}

		config = GetCsvMgr().GetArenaSpecialClassConfig(msg.EnemyClass, msg.EnemyDan)
		myid := 0
		if config != nil {
			myid = config.Id
		}
		GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_ARENA_SPECIAL_FIGHT, 0, int(fight.EnemyTeam[0].Uid), enemyid, "高阶竞技场战斗", 0, myid, player)
	}

	var fightID [ARENA_SPECIAL_TEAM_MAX]int64
	for i, v := range battleInfo {
		fightID[i] = v.Id
	}

	maxnum := vipcsv.ArenaFree[ARENA_TYPE_NOMAL]
	if info.Count > maxnum {
		cost := TARIFF_TYPE_ARENA_SPECIAL

		// 获得消耗配置
		config := GetCsvMgr().GetTariffConfig2(cost)
		if config != nil {
			for i, v := range config.ItemIds {
				if v != 0 {
					msg.Item = append(msg.Item, PassItem{v, -config.ItemNums[i]})
				}
			}
		}
	}
	msg.Cid = MSG_ARENA_SPECIAL_START_FIGHT
	msg.RandNum = battleInfo[0].Random
	msg.FightID = fightID
	msg.BattleInfo = battleInfo
	msg.Index = modArena.Index
	msg.Result = backRet
	msg.FightInfo[0] = myFight
	player.Send(msg.Cid, msg)
	modArena.Index = -1
}

func (self *San_ArenaSpecialPlayer) AddFight(attack, defend [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo, result int, side int, battleInfo [ARENA_SPECIAL_TEAM_MAX]BattleInfo, class int, dan int, needSave bool) *ArenaSpecialFight {
	var enemy [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo
	if side == 1 {
		enemy = defend
	} else {
		enemy = attack
	}

	var fightID [ARENA_SPECIAL_TEAM_MAX]int64

	for i, v := range battleInfo {
		fightID[i] = v.Id
	}

	fight := self.NewPvpFight(fightID, enemy, result, side, class, dan)
	if len(self.arenaFight) >= ARENA_FIGHT_RECORD_MAX {
		self.arenaFight = self.arenaFight[1:]
	}
	self.arenaFight = append(self.arenaFight, fight)

	for i, _ := range battleInfo {

		data2 := BattleRecord{}
		data2.Level = 0
		data2.Side = side
		data2.Time = TimeServer().Unix()
		data2.Id = battleInfo[i].Id
		data2.LevelID = battleInfo[i].LevelID
		data2.Result = result
		data2.Type = BATTLE_TYPE_PVP
		data2.RandNum = battleInfo[i].Random
		data2.FightInfo[0] = attack[i]
		data2.FightInfo[1] = defend[i]

		if needSave {
			HMSetRedisEx("san_arenaspecialbattleinfo", battleInfo[i].Id, &battleInfo[i], HOUR_SECS*12)
			HMSetRedisEx("san_arenaspecialbattlerecord", battleInfo[i].Id, &data2, HOUR_SECS*12)
			GetServer().DBUser.SaveRecord(BATTLE_TYPE_ARENA_SPECIAL, &battleInfo[i], &data2)
		}
	}

	if side == 0 {
		player := GetPlayerMgr().GetPlayer(self.Uid, false)
		if nil != player {
			var msg S2C_ArenaSpecialAddFightRecord
			msg.Cid = MSG_ARENA_SPECIAL_ADD_FIGHT_RECORD
			player.Send(msg.Cid, msg)
		} else {
			self.redPoint.IsFight = 1
		}
	}

	return fight
}

// 无回放战报: 主动-被动
func (self *San_ArenaSpecialPlayer) NewPvpFight(FightID [ARENA_SPECIAL_TEAM_MAX]int64, enemy [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo, result int, side int, class int, dan int) *ArenaSpecialFight {

	p := &ArenaSpecialFight{}
	p.FightId = FightID
	p.Side = side
	p.Result = result
	p.Class = class
	p.Dan = dan
	if enemy[0] != nil {
		p.Uid = enemy[0].Uid
		p.IconId = enemy[0].Iconid
		p.Name = enemy[0].Uname
		p.Level = enemy[0].Level
		p.Portrait = enemy[0].Portrait
		p.Fight = enemy[0].Deffight + enemy[1].Deffight + enemy[2].Deffight
	}
	p.Time = TimeServer().Unix()

	return p
}

func (self *San_ArenaSpecialPlayer) CheckArenaAward(rank *San_ArenaSpecialRank) {
	class := self.Class
	dan := self.Dan

	config := GetCsvMgr().GetArenaSpecialClassConfig(class, dan)
	if config == nil {
		return
	}

	if config.Income == 0 {
		return
	}

	if rank.StartTime == 0 {
		return
	}

	now := TimeServer().Unix()

	addCoint := int(now-rank.StartTime) * config.Income / HOUR_SECS

	if self.Coin+addCoint >= config.Limit {
		addCoint = config.Limit - self.Coin
	}

	self.Coin += addCoint

	if self.Coin < 0 {
		self.Coin = 0
	}

	if config.Integral > 0 {
		addPoint := (now - rank.StartTime) * int64(config.Integral) / HOUR_SECS
		self.Point += addPoint
		rank.Point += addPoint
		GetTopArenaMgr().UpdateArenaSpecialPoint(self.Point, self.Uid)
	}

	rank.StartTime = now
}

func (self *ArenaSpecialMgr) GetArenaAward(player *Player) ([]PassItem, bool) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	uid := player.GetUid()
	data, ok := self.Sql_Uid[uid]
	if !ok {
		return nil, false
	}

	class := data.Class
	dan := data.Dan
	config := GetCsvMgr().GetArenaSpecialClassConfig(class, dan)
	if config == nil {
		return nil, false
	}

	addCoint := 0
	addPoint := int64(0)
	now := TimeServer().Unix()
	for _, v := range self.Sql_Rank[class][dan] {
		if uid == v.Uid {
			addCoint = int(now-v.StartTime) * config.Income / HOUR_SECS
			if config.Integral > 0 {
				addPoint = (now - v.StartTime) * int64(config.Integral) / HOUR_SECS
			}

			if data.Coin+addCoint >= config.Limit {
				addCoint = config.Limit - data.Coin
			}
			v.StartTime = now
			if addPoint > 0 {
				v.Point += addPoint
				data.Point += addPoint
				GetTopArenaMgr().UpdateArenaSpecialPoint(data.Point, uid)
			}
			break
		}
	}
	allAdd := addCoint + data.Coin
	if allAdd < 0 {
		allAdd = 0
	}
	item, num := player.AddObject(ITEM_ARENA_SPECIAL_COIN, allAdd, config.Id, 0, 0, "领取高阶竞技场积累奖励")
	data.Coin = 0

	isfull := false
	if num != allAdd {
		isfull = true
	}

	GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_ARENA_SPECIAL_GET_AWARD, 1, config.Id, num, "领取高阶竞技场积累奖励", 0, 0, player)

	return []PassItem{PassItem{item, num}}, isfull
}

func (self *ArenaSpecialMgr) Rename(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.Sql_Uid[player.GetUid()]
	if ok {
		self.Sql_Uid[player.GetUid()].Name = player.GetName()
		for _, v := range self.Sql_Uid[player.GetUid()].format {
			v.Uname = player.GetName()
		}
	}
}

func (self *ArenaSpecialMgr) Rehead(player *Player) {
	if player == nil {
		return
	}
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.Sql_Uid[player.Sql_UserBase.Uid]
	if ok {
		for _, v := range self.Sql_Uid[player.Sql_UserBase.Uid].format {
			v.Iconid = player.Sql_UserBase.IconId
			v.Portrait = player.Sql_UserBase.Portrait
		}
	}
}

func (self *ArenaSpecialMgr) Relevel(uid int64, newlevel int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.Sql_Uid[uid]
	if ok {
		for _, v := range self.Sql_Uid[uid].format {
			v.Level = newlevel
		}
	}
}

func (self *ArenaSpecialMgr) ArenaFightResultByCross(nType int, battleInfoLst [CROSSARENA3V3_TEAM_MAX]BattleInfo) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	for i, v := range self.FightListForCross {
		data := v

		if battleInfoLst[0].Id != data.FightId[0] {
			continue
		}

		for p, t := range data.FightId {
			if battleInfoLst[p].Id != t {
				for _, z := range data.FightId {
					GetFightMgr().DelResult(z)
				}
				self.FightListForCross = append(self.FightListForCross[0:i], self.FightListForCross[i+1:]...)
				return
			}
		}

		player := GetPlayerMgr().GetPlayer(data.Attack[0].Uid, false)
		if player == nil {
			for _, z := range data.FightId {
				GetFightMgr().DelResult(z)
			}
			self.FightListForCross = append(self.FightListForCross[0:i], self.FightListForCross[i+1:]...)
			return
		}

		for d := 0; d < CROSSARENA3V3_TEAM_MAX; d++ {
			battleInfoLst[d].UserInfo[POS_ATTACK].Uid = data.Attack[d].Uid
			battleInfoLst[d].UserInfo[POS_ATTACK].Name = data.Attack[d].Uname
			battleInfoLst[d].UserInfo[POS_ATTACK].Icon = data.Attack[d].Iconid
			battleInfoLst[d].UserInfo[POS_ATTACK].Portrait = data.Attack[d].Portrait
			battleInfoLst[d].UserInfo[POS_ATTACK].UnionName = data.Attack[d].UnionName
			battleInfoLst[d].UserInfo[POS_ATTACK].Level = data.Attack[d].Level

			battleInfoLst[d].UserInfo[POS_DEFENCE].Uid = data.Defend[d].Uid
			battleInfoLst[d].UserInfo[POS_DEFENCE].Name = data.Defend[d].Uname
			battleInfoLst[d].UserInfo[POS_DEFENCE].Icon = data.Defend[d].Iconid
			battleInfoLst[d].UserInfo[POS_DEFENCE].Portrait = data.Defend[d].Portrait
			battleInfoLst[d].UserInfo[POS_DEFENCE].UnionName = data.Defend[d].UnionName
			battleInfoLst[d].UserInfo[POS_DEFENCE].Level = data.Defend[d].Level

			battleInfoLst[d].Type = BATTLE_TYPE_PVP
			battleInfoLst[d].Time = data.Time
			battleInfoLst[d].Random = data.Random


			data.Result[d].UserInfo=battleInfoLst[d].UserInfo
			data.Result[d].Result=battleInfoLst[d].Result
			data.Result[d].Random=battleInfoLst[d].Random
		}

		result := 0
		for _, p := range battleInfoLst {
			if p.Result == 0 {
				result++
			}
		}

		if result >= 2 {
			player.GetModule("crossarena3v3").(*ModCrossArena3V3).FightEndOK(player, 1, data)
		} else {
			player.GetModule("crossarena3v3").(*ModCrossArena3V3).FightEndOK(player, 2, data)
		}

		for _, z := range data.FightId {
			GetFightMgr().DelResult(z)
		}
		self.FightListForCross = append(self.FightListForCross[0:i], self.FightListForCross[i+1:]...)
		return
	}
}