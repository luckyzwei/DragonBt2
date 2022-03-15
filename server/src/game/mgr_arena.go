package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"runtime/debug"
	"sort"
	"sync"
	"time"
)

const (
	ARENA_RANKS_MAX        = 50
	ARENA_ENEMY_MAX        = 5
	ARENA_BASE_POINT       = 1000
	ARENA_ADD_POINT        = 10
	ARENA_MIN_POINT        = 10
	ARENA_FIGHT_RECORD_MAX = 20
	ARENA_FIGHT_COUNT_MAX  = 100
)

// 竞技场战斗类型
const (
	ARENA_FIGHT_TYPE_NONE       = 0
	ARENA_FIGHT_TYPE_ENEMY      = 1
	ARENA_FIGHT_TYPE_FIGHT_BACK = 2
)

// 竞技场类型
const (
	ARENA_TYPE_NOMAL   = 0
	ARENA_TYPE_SPECIAL = 1
	ARENA_TYPE_HIGHEST = 2
	ARENA_TYPE_MAX     = 3
)

// 竞技场时间类型 取开启结束时间
const (
	ARENA_TIME_TYPE_NOMAL   = 4
	ARENA_TIME_TYPE_SPECIAL = 5
	//ARENA_TIME_TYPE_HIGHEST = 41
)

// 竞技场奖励类型 取奖励
const (
	ARENA_REWARD_TYPE_DAILY  = 1 // 每日
	ARENA_REWARD_TYPE_FINISH = 2 // 结算
)

type San_RankArena struct {
	Uid       int64
	Point     int64 // 积分
	Rank      int
	StartTime int64

	DataUpdate
}

// 竞技场战报
type ArenaFight struct {
	FightId    int64  `json:"fight_id"`      // 战斗Id
	Side       int    `json:"side"`          // 1 进攻方 0 防守方
	Result     int    `json:"attack_result"` // 0 进攻方成功 其他防守方胜利
	Point      int    `json:"point"`         // 积分增减
	Uid        int64  `json:"uid"`           // Uid
	IconId     int    `json:"icon_id"`       // 头像Id
	Name       string `json:"name"`          // 名字
	Level      int    `json:"level"`         // 等级
	Fight      int64  `json:"fight"`         // 战力
	Time       int64  `json:"time"`          // 发生的时间
	Portrait   int    `json:"portrait"`      //! 头像框
	Subsection int    `json:"subsection"`    //! 大段位
	Class      int    `json:"class"`         //! 小段位
}

// 竞技场战报
type ArenaFightList struct {
	Type    int           `json:"type"`
	FightId int64         `json:"fight_id"` // 战斗Id
	Random  int64         `json:"random"`
	Time    int64         `json:"time"`   // 发生的时间
	Attack  *JS_FightInfo `json:"attack"` // 攻击者
	Defend  *JS_FightInfo `json:"defend"` // 防御者
	BossId  int           `json:"bossid"` //
}

type topPlayer []*San_RankArena

func (s topPlayer) Len() int      { return len(s) }
func (s topPlayer) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s topPlayer) Less(i, j int) bool {
	if s[i].Point < s[j].Point {
		return false
	}

	if s[i].Point > s[j].Point {
		return true
	}

	if s[i].StartTime > s[j].StartTime {
		return false
	}

	if s[i].StartTime < s[j].StartTime {
		return true
	}

	return false
}

type ArenaEnemyLs struct {
	Step  int64
	Enemy *JS_FightInfo
}

type RedPoint struct {
	IsFight int
}

// 竞技场
type San_Arena struct {
	Uid        int64
	Rank       int
	Point      int64
	Name       string // 姓名
	Format     string // 阵容
	Count      int    // 今天挑战次数
	RandNum    int64  // 随机数
	ArenaFight string // 战报

	format     *JS_FightInfo   // 只保存防守阵容
	arenaFight []*ArenaFight   // 战报
	enemy      []*ArenaEnemyLs // 敌人
	redpoint   RedPoint        // 红点提醒

	DataUpdate
}

// 将数据库数据写入data
func (self *San_Arena) Decode() {
	json.Unmarshal([]byte(self.Format), &self.format)
	json.Unmarshal([]byte(self.ArenaFight), &self.arenaFight)

}

// 将data数据写入数据库
func (self *San_Arena) Encode() {
	self.Format = HF_JtoA(&self.format)
	self.ArenaFight = HF_JtoA(&self.arenaFight)
}

func (self *San_Arena) AddFight(attack, defend *JS_FightInfo, result int, side int, point int, battleInfo BattleInfo) *ArenaFight {
	var enemy *JS_FightInfo
	if side == 1 {
		enemy = defend
	} else {
		enemy = attack
	}

	fight := self.NewPvpFight(battleInfo.Id, enemy, result, side, point)
	if len(self.arenaFight) >= ARENA_FIGHT_RECORD_MAX {
		self.arenaFight = self.arenaFight[1:]
	}
	self.arenaFight = append(self.arenaFight, fight)

	data2 := BattleRecord{}
	data2.Level = 0
	data2.Side = side
	data2.Time = TimeServer().Unix()
	data2.Id = battleInfo.Id
	data2.LevelID = battleInfo.LevelID
	data2.Result = result
	data2.Type = BATTLE_TYPE_PVP
	data2.RandNum = battleInfo.Random
	data2.FightInfo[0] = attack
	data2.FightInfo[1] = defend

	HMSetRedisEx("san_arenabattleinfo", battleInfo.Id, &battleInfo, HOUR_SECS*12)
	HMSetRedisEx("san_arenabattlerecord", data2.Id, &data2, HOUR_SECS*12)

	GetServer().DBUser.SaveRecord(BATTLE_TYPE_ARENA, &battleInfo, &data2)

	if side == 0 {
		player := GetPlayerMgr().GetPlayer(self.Uid, false)
		if nil != player {
			var msg S2C_ArenaAddFightRecord
			msg.Cid = "arena_add_fight_record"
			player.Send(msg.Cid, msg)
		} else {
			self.redpoint.IsFight = 1
		}
	}

	return fight
}

func (self *San_Arena) DeleteFight(fightID int64) {
	nLen := len(self.arenaFight)
	for i := 0; i < nLen; i++ {
		if self.arenaFight[i].FightId == fightID {
			self.arenaFight = append(self.arenaFight[:i], self.arenaFight[i+1:]...)
			break
		}
	}
}

// 无回放战报: 主动-被动
func (self *San_Arena) NewPvpFight(FightID int64, enemy *JS_FightInfo, result int, side int, point int) *ArenaFight {

	p := &ArenaFight{}
	p.FightId = FightID
	p.Side = side
	p.Result = result
	p.Point = point
	if enemy != nil {
		p.Uid = enemy.Uid
		p.IconId = enemy.Iconid
		p.Name = enemy.Uname
		p.Level = enemy.Level
		p.Fight = enemy.Deffight
		p.Portrait = enemy.Portrait
	}
	p.Time = TimeServer().Unix()

	return p
}

// 更新阵容
func (self *ArenaMgr) UpdateFormat(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	uid := player.GetUid()
	_, ok := self.Sql_Uid[uid]
	if !ok {
		return
	}
	//posTypes := self.GetPosTypeByArenaType(false)
	self.Sql_Uid[uid].format = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_2)
	self.Sql_Uid[uid].Format = HF_JtoA(&self.Sql_Uid[uid].format)
	team := player.GetModule("team").(*ModTeam).getTeamPos(TEAMTYPE_ARENA_2)
	var herolst []int
	for _, p := range team.FightPos {
		herolst = append(herolst, p)
	}
	self.Sql_Uid[uid].Update(true)

	GetTopArenaMgr().UpdateFight(ARENA_TOP_TYPE_NOMAL, player)
}

// 机器人缓存
type ArenaRobot struct {
	Point int64
	Info  *JS_FightInfo
}

func (self *ArenaMgr) InitRobot() {
	//for i := 0; i < len(GetCsvMgr().JJCRobotConfig); i++ {
	//	cfg := GetCsvMgr().JJCRobotConfig[i]
	//	if cfg.Type != 1 {
	//		continue
	//	}
	//
	//	var node ArenaRobot
	//	node.Point = ARENA_BASE_POINT
	//	node.Info = GetCsvMgr().GetRobot(cfg)
	//	self.Robot = append(self.Robot, node)
	//}
}

type San_ArenaTime struct {
	Type        int
	StartTime   int64
	RefreshTime int64

	DataUpdate
}

// 竞技场管理器
type ArenaMgr struct {
	Sql_Uid        map[int64]*San_Arena   // 具体数据
	Sql_Rank       map[int]*San_RankArena // 排序用数据
	Lock           *sync.RWMutex          // 数据操作锁
	MaxRank        int                    // 最低排名，30000开外
	Robot          []ArenaRobot           // 机器人数据
	ArenaTime      San_ArenaTime
	migrateOK      bool
	FightList      []*ArenaFightList
	FightListBoss  []*ArenaFightList //增加暗域入侵战斗队列 20201122
	FightListCross []*ArenaFightList //增加跨服竞技战斗队列 20210105
}

var arenamgr *ArenaMgr = nil

// public
func GetArenaMgr() *ArenaMgr {
	if arenamgr == nil {
		arenamgr = new(ArenaMgr)
		arenamgr.Sql_Uid = make(map[int64]*San_Arena)
		arenamgr.Sql_Rank = make(map[int]*San_RankArena)
		arenamgr.Lock = new(sync.RWMutex)
		arenamgr.migrateOK = false
	}
	return arenamgr
}

// 开启迁移数据协程，竞技场
func (self *ArenaMgr) RunMigrateArena() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	infoName := "san_arenabattleinfo"     //! info 表
	recordName := "san_arenabattlerecord" //! record 表
	recordType := 2
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

func (self *ArenaMgr) GetPlayerArenaData(uid int64) *San_Arena {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.Sql_Uid[uid]
	if ok {
		return self.Sql_Uid[uid]
	}
	return nil
}

func (self *ArenaMgr) GetPlayerRankData(rank int) *San_RankArena {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.Sql_Rank[rank]
	if ok {
		return self.Sql_Rank[rank]
	}
	return nil
}

func (self *ArenaMgr) GetData() {
	self.Sql_Uid = make(map[int64]*San_Arena)
	self.Sql_Rank = make(map[int]*San_RankArena)
	var info1 San_RankArena
	tableName1 := fmt.Sprintf("san_rankarena%d", ARENA_TYPE_NOMAL+1)
	sql1 := fmt.Sprintf("select * from `%s`", tableName1)
	res1 := GetServer().DBUser.GetAllData(sql1, &info1)
	temp := topPlayer{}
	for i := 0; i < len(res1); i++ {
		data := res1[i].(*San_RankArena)
		data.Init(tableName1, data, true)
		temp = append(temp, data)
	}

	var info2 San_Arena
	tableName2 := fmt.Sprintf("san_playerarena%d", ARENA_TYPE_NOMAL+1)
	sql2 := fmt.Sprintf("select * from `%s`", tableName2)
	res2 := GetServer().DBUser.GetAllData(sql2, &info2)
	for i := 0; i < len(res2); i++ {
		data := res2[i].(*San_Arena)
		data.Init(tableName2, data, true)
		data.Decode()
		self.Sql_Uid[data.Uid] = data
	}

	sort.Sort(topPlayer(temp))
	nLen := temp.Len()
	for i := 1; i <= nLen; i++ {
		temp[i-1].Rank = i
		self.Sql_Rank[i] = temp[i-1]
		self.Sql_Uid[self.Sql_Rank[i].Uid].Rank = i
		self.Sql_Uid[self.Sql_Rank[i].Uid].Point = temp[i-1].Point
		self.Sql_Rank[i].Update(true)
		self.Sql_Uid[self.Sql_Rank[i].Uid].Update(true)
	}

	self.MaxRank = len(self.Sql_Rank)

	for _, value := range self.Sql_Uid {
		if value.Rank < 0 {
			continue
		}
		fight := GetOfflineInfoMgr().GetTeamFight(value.Uid, TEAMTYPE_ARENA_2)
		if fight <= 0 && value.format != nil {
			GetOfflineInfoMgr().SetArenaFight(TEAMTYPE_ARENA_2, value.Uid, value.format.Deffight)
		}
	}

	sql := fmt.Sprintf("select * from `san_arenatime` where type = %d", ARENA_TYPE_NOMAL+1)
	res := GetServer().DBUser.GetAllData(sql, &self.ArenaTime)
	if len(res) > 0 {
		self.ArenaTime = *res[0].(*San_ArenaTime)
	}

	if self.ArenaTime.Type <= 0 {
		self.ArenaTime.Type = ARENA_TYPE_NOMAL + 1
		startTime, _ := GetCsvMgr().GetNowStartAndEnd(ARENA_TIME_TYPE_NOMAL)
		self.ArenaTime.StartTime = startTime
		now := TimeServer()
		timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
		if TimeServer().Hour() < 5 {
			timeStamp -= DAY_SECS
		}
		self.ArenaTime.RefreshTime = timeStamp
		InsertTable("san_arenatime", &self.ArenaTime, 0, false)
	}

	self.ArenaTime.Init("san_arenatime", &self.ArenaTime, false)

	self.InitRobot()

	self.CheckArenaEnd()
}

func (self *ArenaMgr) Save() {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	for _, value := range self.Sql_Uid {
		value.Encode()
		value.Update(true)
	}

	for _, value := range self.Sql_Rank {
		value.Update(true)
	}

	self.ArenaTime.Update(true)
}

// 添加玩家
func (self *ArenaMgr) AddPlayer(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	tableName1 := fmt.Sprintf("san_rankarena%d", ARENA_TYPE_NOMAL+1)
	tableName2 := fmt.Sprintf("san_playerarena%d", ARENA_TYPE_NOMAL+1)
	uid := player.GetUid()
	// 通过uid找玩家的自身数据
	_, ok := self.Sql_Uid[uid]
	// 找到了说明之前参与过
	if ok {
		// 查询有没有排行数据
		rank := self.Sql_Uid[uid].Rank
		if rank > 0 {
			// 有排行数据 且就是自身数据 说明是正常的 直接返回
			_, ok2 := self.Sql_Rank[rank]
			if ok2 {
				if self.Sql_Rank[rank].Uid == uid {
					if self.Sql_Rank[rank].Point != self.Sql_Uid[uid].Point {
						self.Sql_Uid[uid].Point = self.Sql_Rank[rank].Point

						LogDebug("AddPlayer ERRER !!!!!!!!!! self.Sql_Rank[rank].Point != self.Sql_Uid[uid].Point")
						// 是否有前一名
						_, ok3 := self.Sql_Rank[rank-1]
						if ok3 {
							// 有且积分比自身低 则排序返回
							if self.Sql_Rank[rank-1].Point < self.Sql_Rank[rank].Point {
								self.UpRank(rank)
								GetTopArenaMgr().UpdateArenaRank(self.Sql_Uid[uid].Point, uid)
								LogDebug("AddPlayer ERRER !!!!!!!!!! self.Sql_Rank[rank-1].Point < self.Sql_Rank[rank].Point")
								return
							}
						}
						_, ok4 := self.Sql_Rank[rank+1]
						if ok4 {
							// 有且积分比自身高 则排序返回
							if self.Sql_Rank[rank+1].Point > self.Sql_Rank[rank].Point {
								self.DownRank(rank)
								GetTopArenaMgr().UpdateArenaRank(self.Sql_Uid[uid].Point, uid)
								LogDebug("AddPlayer ERRER !!!!!!!!!! self.Sql_Rank[rank+1].Point > self.Sql_Rank[rank].Point")
								return
							}
						}
						GetTopArenaMgr().UpdateArenaRank(self.Sql_Uid[uid].Point, uid)
					}
					return
				}
			}
			// 没找到 或者是不是自身数据 则出现错误 循环查找
			find := -1
			for i, t := range self.Sql_Rank {
				if t.Uid == uid {
					find = i
					break
				}
			}
			// 找到了数据
			if find > 0 {
				LogDebug("AddPlayer ERRER !!!!!!!!!! find > 0")
				// 以排行积分为主
				self.Sql_Uid[uid].Point = self.Sql_Rank[find].Point
				self.Sql_Uid[uid].Rank = self.Sql_Rank[find].Rank
				// 是否有前一名
				_, ok3 := self.Sql_Rank[find-1]
				if ok3 {
					// 有且积分比自身低 则排序返回
					if self.Sql_Rank[find-1].Point < self.Sql_Rank[find].Point {
						self.UpRank(find)
						GetTopArenaMgr().UpdateArenaRank(self.Sql_Uid[uid].Point, uid)
						LogDebug("AddPlayer ERRER !!!!!!!!!! find > 0 self.Sql_Rank[find-1].Point < self.Sql_Rank[find].Point")
						return
					}
				}
				_, ok4 := self.Sql_Rank[find+1]
				if ok4 {
					// 有且积分比自身低 则排序返回
					if self.Sql_Rank[find+1].Point > self.Sql_Rank[find].Point {
						self.DownRank(find)
						GetTopArenaMgr().UpdateArenaRank(self.Sql_Uid[uid].Point, uid)
						LogDebug("AddPlayer ERRER !!!!!!!!!! find > 0 self.Sql_Rank[find+1].Point > self.Sql_Rank[find].Point")
						return
					}
				}
				return
			}
		}

		// 自然过期的 完全没找到数据的 则说明是参与新赛季 将排名置为最后一名 积分不变
		self.Sql_Uid[uid].Rank = self.MaxRank + 1
	} else {
		// 没找到 则说明是新参与玩家 初始化数据
		self.Sql_Uid[uid] = &San_Arena{}
		self.Sql_Uid[uid].Uid = uid
		self.Sql_Uid[uid].Rank = self.MaxRank + 1
		self.Sql_Uid[uid].Name = player.GetName()
		self.Sql_Uid[uid].Point = ARENA_BASE_POINT
		self.Sql_Uid[uid].Encode()
		InsertTable(tableName2, self.Sql_Uid[uid], 0, true)
		self.Sql_Uid[uid].Init(tableName2, self.Sql_Uid[uid], true)
	}

	// 初始化玩家防御阵容
	self.Sql_Uid[uid].format = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_2)
	self.Sql_Uid[uid].Format = HF_JtoA(&self.Sql_Uid[uid].format)
	self.Sql_Uid[uid].Count = 0

	// 走到这里两种情况 一种是前赛季参与 重新进入的 第二种是全进进入的新号 都是置为了最后一名
	self.Sql_Rank[self.MaxRank+1] = &San_RankArena{}
	self.Sql_Rank[self.MaxRank+1].Uid = uid
	self.Sql_Rank[self.MaxRank+1].Rank = self.MaxRank + 1
	self.Sql_Rank[self.MaxRank+1].Point = self.Sql_Uid[uid].Point
	self.Sql_Rank[self.MaxRank+1].StartTime = TimeServer().Unix()
	InsertTable(tableName1, self.Sql_Rank[self.MaxRank+1], 0, true)
	self.Sql_Rank[self.MaxRank+1].Init(tableName1, self.Sql_Rank[self.MaxRank+1], true)
	// 向上排序
	self.UpRank(self.MaxRank + 1)
	// 增加人数
	self.MaxRank += 1

	GetTopArenaMgr().UpdateArenaRank(self.Sql_Uid[uid].Point, uid)
}

// 向前排序 提升名次
func (self *ArenaMgr) UpRank(rank int) {
	nCount := 0
	// 向第一名循环
	for i := rank - 1; i >= 1; i-- {
		// 没有前一名 则可以往上挪 一般不可能
		_, ok := self.Sql_Rank[i]
		if !ok {
			nCount++
			continue
		}

		// 找不到该玩家数据 一般也不可能
		_, ok2 := self.Sql_Uid[self.Sql_Rank[i].Uid]
		if !ok2 {
			nCount++
			continue
		}

		// 上一名比我分数低 则可以往上排序
		if self.Sql_Rank[i].Point < self.Sql_Rank[rank].Point {
			nCount++ //计算挪多少位
		} else {
			break // 和我分数相等或大于则返回
		}
	}

	// 获得原来的指针
	change := self.Sql_Rank[rank]
	// 从找到的位置 开始 交换指针
	for i := rank - nCount; i <= rank && i > 0; i++ {
		temp := self.Sql_Rank[i]
		self.Sql_Rank[i] = change
		self.Sql_Rank[i].Rank = i // 排名变化
		change = temp

		self.Sql_Uid[self.Sql_Rank[i].Uid].Rank = i // 排名变化
	}
}

func (self *ArenaMgr) DownRank(rank int) {
	nCount := 0
	// 向最后一名循环
	for i := rank + 1; i <= self.MaxRank; i++ {
		// 没有后名  一般不可能
		_, ok := self.Sql_Rank[i]
		if !ok {
			nCount++
			continue
		}

		// 对象不存在 一般不可能
		_, ok2 := self.Sql_Uid[self.Sql_Rank[i].Uid]
		if !ok2 {
			nCount++
			continue
		}

		// 后一名的分数比我高
		if self.Sql_Rank[i].Point >= self.Sql_Rank[rank].Point {
			nCount++
		} else {
			break
		}
	}
	// 获得原来的指针
	change := self.Sql_Rank[rank]
	// 从找到的位置 开始 交换指针
	for i := rank + nCount; i >= rank && i <= self.MaxRank; i-- {
		temp := self.Sql_Rank[i]
		self.Sql_Rank[i] = change
		self.Sql_Rank[i].Rank = i // 排名变化
		change = temp

		self.Sql_Uid[self.Sql_Rank[i].Uid].Rank = i // 排名变化
	}
}

// 加分
func (self *ArenaMgr) AddPoint(uid int64, rank int, addPonit int64) {
	data, ok := self.Sql_Rank[rank]
	if !ok {
		return
	}

	playerInfo, ok2 := self.Sql_Uid[data.Uid]
	if !ok2 {
		return
	}

	if data.Uid != uid {
		return
	}

	// 加分
	data.Point += addPonit
	data.StartTime = TimeServer().Unix()
	playerInfo.Point = data.Point

	// 向上排序
	self.UpRank(rank)

	// 跟新任务和排行榜
	player := GetPlayerMgr().GetPlayer(uid, false)
	if player != nil {
		player.HandleTask(TASK_TYPE_ARENA_POINT, int(data.Point), 0, 0)
	}
	GetTopArenaMgr().UpdateArenaRank(data.Point, uid)
}

// 减分
func (self *ArenaMgr) MinPoint(uid int64, rank int, minPonit int64) {
	data, ok := self.Sql_Rank[rank]
	if !ok {
		return
	}

	playerInfo, ok2 := self.Sql_Uid[data.Uid]
	if !ok2 {
		return
	}

	if data.Uid != uid {
		return
	}

	// 减分
	data.Point -= minPonit
	data.StartTime = TimeServer().Unix()
	playerInfo.Point = data.Point

	// 向下排序
	self.DownRank(rank)

	// 通知任务和排行榜
	//player := GetPlayerMgr().GetPlayer(uid, false)
	GetTopArenaMgr().UpdateArenaRank(data.Point, uid)
}

// 增加战斗记录
func (self *ArenaMgr) AddFight(uid int64, myinfo, enemy *JS_FightInfo, result int, side int, point int, battleInfo BattleInfo) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	data, ok := self.Sql_Uid[uid]
	if !ok {
		return
	}

	if data.Uid != uid {
		return
	}

	data.AddFight(myinfo, enemy, result, side, point, battleInfo)
}

// 删除任务记录
func (self *ArenaMgr) DeleteFight(uid int64, fightID int64) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	data, ok := self.Sql_Uid[uid]
	if !ok {
		return
	}

	if data.Uid != uid {
		return
	}

	data.DeleteFight(fightID)
}

//func (self *ArenaData) GetPosTypeByArenaType(attack bool) int {
//	if attack {
//		return TEAMTYPE_ARENA_1
//	} else {
//		return TEAMTYPE_ARENA_2
//	}
//}

// 获取机器人
func (self *ArenaMgr) GetRobot() *JS_FightInfo {
	var robot JS_FightInfo
	//if len(self.Robot) <= 0 {
	//	return nil
	//}

	var node ArenaRobot
	node.Info = GetCsvMgr().GetRobot(GetCsvMgr().JJCRobotConfig[0])
	HF_DeepCopy(&robot, node.Info)
	return &robot
}

// 获得敌人 每次都会刷新
func (self *ArenaMgr) GetEnemy(uid int64) ([]*JS_FightInfo, []int64) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	data, ok := self.Sql_Uid[uid]
	if !ok {
		return nil, nil
	}

	if data.Rank < 0 {
		return nil, nil
	}

	_, ok2 := self.Sql_Rank[data.Rank]
	if !ok2 {
		return nil, nil
	}

	if len(data.arenaFight) == 0 {
		fightinfo := make([]*JS_FightInfo, 0)
		fightenemy := make([]*ArenaEnemyLs, 0)
		point := make([]int64, 0)

		value := int64(10000)
		for i := 0; i < 3; i++ {
			robot := self.GetRobot()
			if robot != nil {
				robot.Deffight += value
				value -= int64(HF_GetRandom(10000))
				fightinfo = append(fightinfo, robot)
				point = append(point, ARENA_BASE_POINT)
				fightenemy = append(fightenemy, &ArenaEnemyLs{int64(i), robot})
			}
		}

		data.enemy = fightenemy
		return fightinfo, point
	}

	base := self.Sql_Rank[data.Rank].Point
	// 积分区间结构体
	type PointStep struct {
		MinPoint int64
		MaxPoint int64
	}

	configs := GetCsvMgr().ArenaParameter
	if len(configs) != 5 {
		return nil, nil
	}

	// 五个区间 向上找三个
	upPoint := []PointStep{PointStep{base * int64(configs[0].Matemin) / 100, base * int64(configs[0].Matemax) / 100},
		PointStep{base * int64(configs[1].Matemin) / 100, base * int64(configs[1].Matemax) / 100},
		PointStep{base * int64(configs[2].Matemin) / 100, base * int64(configs[2].Matemax) / 100}}
	// 向下找两个
	downPoint := []PointStep{PointStep{base * int64(configs[3].Matemin) / 100, base * int64(configs[3].Matemax) / 100},
		PointStep{base * int64(configs[4].Matemin) / 100, base * int64(configs[4].Matemax) / 100}}

	upLen := len(upPoint)
	downLen := len(downPoint)
	rank := make([][]int, upLen+downLen)
	nLen := len(self.Sql_Rank)
	for i := data.Rank - 1; i > 0 && i <= nLen; i-- {
		point := self.Sql_Rank[i].Point
		if uid == self.Sql_Rank[i].Uid {
			continue
		}

		for t, v := range upPoint {
			if point >= v.MinPoint && point < v.MaxPoint {
				rank[t] = append(rank[t], i)
				break
			}
		}

		// 超过最大则返回
		if point >= upPoint[0].MaxPoint {
			break
		}
	}

	//index := -1
	for i := data.Rank + 1; i > 0 && i <= nLen; i++ {
		point := self.Sql_Rank[i].Point
		if uid == self.Sql_Rank[i].Uid {
			continue
		}
		for t, v := range downPoint {
			if point >= v.MinPoint && point < v.MaxPoint {
				rank[t+3] = append(rank[t+3], i)
				break
			}
		}
		// 超过最低则返回
		if point < downPoint[1].MinPoint {
			//index = i
			break
		}
	}

	//// 随出五个对手循环
	//for i := 0; i < ARENA_ENEMY_MAX; i++ {
	//	find := false
	//	// 从对应的层级开始 往下找
	//	for t := i; t < ARENA_ENEMY_MAX; t++ {
	//		count := len(rank[t])
	//		if count <= 0 {
	//			continue
	//		}
	//
	//		temp := HF_GetRandom(count)
	//		fightid = append(fightid, rank[t][temp])
	//		rank[t] = append(rank[t][:temp], rank[t][temp+1:]...)
	//		find = true
	//		break
	//	}
	//
	//	if find {
	//		continue
	//	}
	//
	//	if index > 0 && nLen-index > 0 {
	//		temp := HF_GetRandom(nLen - index)
	//		fightid = append(fightid, temp+index+1)
	//		break
	//	}
	//}

	fightinfo := make([]*JS_FightInfo, 0)
	fightenemy := make([]*ArenaEnemyLs, 0)

	point := make([]int64, 0)

	// 规则暂时修改为每一段从自身找 不到下一层找
	for t := 0; t < ARENA_ENEMY_MAX; t++ {
		count := len(rank[t])
		if count <= 0 {
			continue
		}

		temp := HF_GetRandom(count)
		if self.Sql_Uid[self.Sql_Rank[rank[t][temp]].Uid].format != nil {
			fightinfo = append(fightinfo, self.Sql_Uid[self.Sql_Rank[rank[t][temp]].Uid].format)
			point = append(point, self.Sql_Uid[self.Sql_Rank[rank[t][temp]].Uid].Point)
			fightenemy = append(fightenemy, &ArenaEnemyLs{int64(t), self.Sql_Uid[self.Sql_Rank[rank[t][temp]].Uid].format})
		}
		rank[t] = append(rank[t][:temp], rank[t][temp+1:]...)
	}

	if len(fightinfo) == 0 {
		nLen := len(fightinfo)
		for i := 0; i < ARENA_ENEMY_MAX-nLen; i++ {
			robot := self.GetRobot()
			if robot != nil {
				fightinfo = append(fightinfo, robot)
				point = append(point, ARENA_BASE_POINT)

				fightenemy = append(fightenemy, &ArenaEnemyLs{int64(4), robot})
			}
		}
	}

	data.enemy = fightenemy
	return fightinfo, point
}

func (self *ArenaMgr) CheckArenaEnd() {
	dayrefresh := false
	// 每日结算
	now := TimeServer()
	timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
	if TimeServer().Hour() < 5 {
		timeStamp -= DAY_SECS
	}

	if timeStamp != self.ArenaTime.RefreshTime {
		dayrefresh = true
	}

	// 检测活动是否结束
	startTime, endTime := GetCsvMgr().GetNowStartAndEnd(ARENA_TIME_TYPE_NOMAL)
	isFinish := false
	// 开服设置
	if self.ArenaTime.StartTime == 0 {
		self.ArenaTime.StartTime = startTime
	} else if startTime != self.ArenaTime.StartTime { // 开始时间不同则直接过期
		self.ArenaTime.StartTime = startTime
		self.ArenaTime.RefreshTime = timeStamp
		dayrefresh = false
		isFinish = true
		self.ArenaTime.Update(true)
	} else if TimeServer().Unix() >= endTime { // 已经结束 全部过期
		self.ArenaTime.StartTime = startTime
		isFinish = true
		self.ArenaTime.Update(true)
	}

	if dayrefresh && !isFinish {
		nLen := len(self.Sql_Rank)
		for t := 1; t <= nLen; t++ {
			player := GetPlayerMgr().GetPlayer(self.Sql_Rank[t].Uid, true)
			if player == nil {
				continue
			}
			config := GetCsvMgr().GetArenaRewardConfig(ARENA_REWARD_TYPE_DAILY, t)
			if nil == config {
				continue
			}

			out := []PassItem{}

			itemLen := len(config.Items)
			if itemLen != len(config.Nums) {
				continue
			}

			for j := 0; j < itemLen; j++ {
				if config.Items[j] != 0 {
					out = append(out, PassItem{config.Items[j], config.Nums[j]})
				}
			}

			mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_GET_ARENA_DAY]
			if ok && mailConfig != nil {
				context := GetCsvMgr().GetText(fmt.Sprintf(mailConfig.Mailtxt, t))
				player.GetModule("mail").(*ModMail).AddMail(1, 1, 0, mailConfig.Mailtitle,
					context, GetCsvMgr().GetText("STR_SYS"),
					out, true, TimeServer().Unix())
			}
		}

		for _, p := range self.Sql_Uid {
			p.Count = 0
		}

		self.ArenaTime.RefreshTime = timeStamp
	}

	tableName1 := fmt.Sprintf("san_rankarena%d", ARENA_TYPE_NOMAL+1)
	//tableName2 := fmt.Sprintf("san_playerarena%d", ARENA_TYPE_NOMAL+1)
	if isFinish {
		nLen := len(self.Sql_Rank)
		for t := 1; t <= nLen; t++ {
			// 所有人排名置为负数
			self.Sql_Uid[self.Sql_Rank[t].Uid].Rank = -1

			player := GetPlayerMgr().GetPlayer(self.Sql_Rank[t].Uid, true)
			if player == nil {
				DeleteTable(tableName1, self.Sql_Rank[t], []int{0})
				continue
			}
			config := GetCsvMgr().GetArenaRewardConfig(ARENA_REWARD_TYPE_FINISH, t)
			if nil == config {
				DeleteTable(tableName1, self.Sql_Rank[t], []int{0})
				continue
			}

			out := []PassItem{}

			itemLen := len(config.Items)
			if itemLen != len(config.Nums) {
				DeleteTable(tableName1, self.Sql_Rank[t], []int{0})
				continue
			}

			for j := 0; j < itemLen; j++ {
				if config.Items[j] != 0 {
					out = append(out, PassItem{config.Items[j], config.Nums[j]})
				}
			}

			mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_GET_ARENA_END]
			if ok && mailConfig != nil {
				context := GetCsvMgr().GetText(fmt.Sprintf(mailConfig.Mailtxt, t))
				player.GetModule("mail").(*ModMail).AddMail(1, 1, 0, mailConfig.Mailtitle,
					context, GetCsvMgr().GetText("STR_SYS"),
					out, true, TimeServer().Unix())

				DeleteTable(tableName1, self.Sql_Rank[t], []int{0})
			}

		}

		self.Sql_Rank = map[int]*San_RankArena{}
		self.MaxRank = 0
	}
}

// 开启战斗协程
func (self *ArenaMgr) StartFight() {
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
func (self *ArenaMgr) OnTimer() {
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
			player := GetPlayerMgr().GetPlayer(data.Attack.Uid, false)
			if player == nil {
				continue
			}

			fightClient := &FightResult{}
			battleInfo := BattleInfo{}
			fightClient.Id = data.FightId
			fightClient.Time = int(TimeServer().Unix())
			fightClient.CityId = 0
			fightClient.Fight[0] = data.Attack
			fightClient.Fight[1] = data.Defend
			if data.Attack.Deffight > data.Defend.Deffight {
				fightClient.Result = 1
				battleInfo.Result = 0
			} else {
				fightClient.Result = 2
				battleInfo.Result = 1
			}
			fightClient.SecKill = 0

			battleInfo.Id = data.FightId
			var attack []*BattleHeroInfo
			var defend []*BattleHeroInfo
			for _, v := range data.Attack.Heroinfo {
				attack = append(attack, &BattleHeroInfo{v.Heroid, v.Levels, v.Stars, v.Skin, 123, 456, 789, 11, 12, nil, 0, nil})
			}
			for _, v := range data.Defend.Heroinfo {
				defend = append(defend, &BattleHeroInfo{v.Heroid, v.Levels, v.Stars, v.Skin, 987, 654, 321, 12, 45, nil, 0, nil})
			}
			battleInfo.UserInfo[POS_ATTACK] = &BattleUserInfo{
				data.Attack.Uid,
				data.Attack.Uname,
				data.Attack.Iconid,
				data.Attack.Portrait,
				data.Attack.UnionName,
				data.Attack.Level,
				attack}
			battleInfo.UserInfo[POS_DEFENCE] = &BattleUserInfo{
				data.Defend.Uid,
				data.Defend.Uname,
				data.Defend.Iconid,
				data.Defend.Portrait,
				data.Defend.UnionName,
				data.Defend.Level,
				defend}
			battleInfo.Type = BATTLE_TYPE_PVP
			battleInfo.Time = data.Time
			battleInfo.Random = data.Random

			if data.Type == ARENA_FIGHT_TYPE_ENEMY {
				self.FightEnd(player, fightClient.Result, battleInfo)
			} else if data.Type == ARENA_FIGHT_TYPE_FIGHT_BACK {
				self.FightBackEnd(player, fightClient.Result, battleInfo)
			}

			GetFightMgr().DelResult(data.FightId)
			self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
			continue
		}

		// 尝试获取战斗结果
		//LogDebug("尝试获取战斗结果")
		FightResult := GetFightMgr().GetResult(data.FightId)
		// 有结果,设置战斗时间
		if FightResult == nil {
			continue
		}

		player := GetPlayerMgr().GetPlayer(data.Attack.Uid, false)
		if player == nil {
			GetFightMgr().DelResult(FightResult.Id)
			self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
			continue
		}

		battleInfo := BattleInfo{}
		battleInfo.Id = data.FightId
		attackHeroInfo := []*BattleHeroInfo{}
		for i, v := range FightResult.Info[POS_ATTACK] {
			level, star, skin, exclusiveLv := 0, 0, 0, 0
			if i < len(FightResult.Fight[POS_ATTACK].Heroinfo) {
				level = FightResult.Fight[POS_ATTACK].Heroinfo[i].Levels
				star = FightResult.Fight[POS_ATTACK].Heroinfo[i].Stars
				skin = FightResult.Fight[POS_ATTACK].Heroinfo[i].Skin
				exclusiveLv = FightResult.Fight[POS_ATTACK].Heroinfo[i].HeroExclusiveLv
			}
			attackHeroInfo = append(attackHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusiveLv, nil})
		}
		defendHeroInfo := []*BattleHeroInfo{}
		for i, v := range FightResult.Info[POS_DEFENCE] {
			level, star, skin, exclusiveLv := 0, 0, 0, 0
			if i < len(FightResult.Fight[POS_DEFENCE].Heroinfo) {
				level = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Levels
				star = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Stars
				skin = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Skin
				exclusiveLv = FightResult.Fight[POS_DEFENCE].Heroinfo[i].HeroExclusiveLv
			}
			defendHeroInfo = append(defendHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusiveLv, nil})
		}

		battleInfo.UserInfo[POS_ATTACK] = &BattleUserInfo{
			data.Attack.Uid,
			data.Attack.Uname,
			data.Attack.Iconid,
			data.Attack.Portrait,
			data.Attack.UnionName,
			data.Attack.Level,
			attackHeroInfo}
		battleInfo.UserInfo[POS_DEFENCE] = &BattleUserInfo{
			data.Defend.Uid,
			data.Defend.Uname,
			data.Defend.Iconid,
			data.Defend.Portrait,
			data.Defend.UnionName,
			data.Defend.Level,
			defendHeroInfo}
		battleInfo.Type = BATTLE_TYPE_PVP
		battleInfo.Time = data.Time
		battleInfo.Random = data.Random
		if FightResult.Result == 1 {
			battleInfo.Result = 0
		} else {
			battleInfo.Result = 1
		}

		if data.Type == ARENA_FIGHT_TYPE_ENEMY {
			self.FightEnd(player, FightResult.Result, battleInfo)
		} else if data.Type == ARENA_FIGHT_TYPE_FIGHT_BACK {
			self.FightBackEnd(player, FightResult.Result, battleInfo)
		}

		GetFightMgr().DelResult(FightResult.Id)
		self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
	}

	self.OnTimerBoss()
	self.OnTimerCross()
}

func (self *ArenaMgr) OnTimerBoss() {

	now := TimeServer().Unix()

	nLen := len(self.FightListBoss)
	for i := nLen - 1; i >= 0; i-- {
		data := self.FightListBoss[i]
		if data.Time+10 <= now {
			//超时
			player := GetPlayerMgr().GetPlayer(data.Attack.Uid, false)
			if player == nil {
				continue
			}
			player.GetModule("activityboss").(*ModActivityBoss).FightEndFail()
			GetFightMgr().DelResult(data.FightId)
			self.FightListBoss = append(self.FightListBoss[0:i], self.FightListBoss[i+1:]...)
			continue
		}

		// 尝试获取战斗结果
		//LogDebug("尝试获取战斗结果")
		FightResult := GetFightMgr().GetResult(data.FightId)
		// 有结果,设置战斗时间
		if FightResult == nil {
			continue
		}

		player := GetPlayerMgr().GetPlayer(data.Attack.Uid, false)
		if player == nil {
			GetFightMgr().DelResult(FightResult.Id)
			self.FightListBoss = append(self.FightListBoss[0:i], self.FightListBoss[i+1:]...)
			continue
		}

		battleInfo := BattleInfo{}
		battleInfo.Id = data.FightId
		attackHeroInfo := []*BattleHeroInfo{}
		for i, v := range FightResult.Info[POS_ATTACK] {
			level, star, skin, exclusiveLv := 0, 0, 0, 0
			if i < len(FightResult.Fight[POS_ATTACK].Heroinfo) {
				level = FightResult.Fight[POS_ATTACK].Heroinfo[i].Levels
				star = FightResult.Fight[POS_ATTACK].Heroinfo[i].Stars
				skin = FightResult.Fight[POS_ATTACK].Heroinfo[i].Skin
				exclusiveLv = FightResult.Fight[POS_ATTACK].Heroinfo[i].HeroExclusiveLv
			}
			attackHeroInfo = append(attackHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusiveLv, nil})
		}
		defendHeroInfo := []*BattleHeroInfo{}
		for i, v := range FightResult.Info[POS_DEFENCE] {
			level, star, skin, exclusiveLv := 0, 0, 0, 0
			if i < len(FightResult.Fight[POS_DEFENCE].Heroinfo) {
				level = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Levels
				star = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Stars
				skin = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Skin
				exclusiveLv = FightResult.Fight[POS_DEFENCE].Heroinfo[i].HeroExclusiveLv
			}
			defendHeroInfo = append(defendHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusiveLv, nil})
		}

		battleInfo.UserInfo[POS_ATTACK] = &BattleUserInfo{
			data.Attack.Uid,
			data.Attack.Uname,
			data.Attack.Iconid,
			data.Attack.Portrait,
			data.Attack.UnionName,
			data.Attack.Level,
			attackHeroInfo}
		battleInfo.UserInfo[POS_DEFENCE] = &BattleUserInfo{
			data.Defend.Uid,
			data.Defend.Uname,
			data.Defend.Iconid,
			data.Defend.Portrait,
			data.Defend.UnionName,
			data.Defend.Level,
			defendHeroInfo}
		battleInfo.Type = BATTLE_TYPE_PVP
		battleInfo.Time = data.Time
		battleInfo.Random = data.Random
		if FightResult.Result == 1 {
			battleInfo.Result = 0
		} else {
			battleInfo.Result = 1
		}

		player.GetModule("activityboss").(*ModActivityBoss).FightEndOK(player, battleInfo, data.Attack, data.Defend, data.BossId)

		GetFightMgr().DelResult(FightResult.Id)
		self.FightListBoss = append(self.FightListBoss[0:i], self.FightListBoss[i+1:]...)
	}
}

func (self *ArenaMgr) ArenaFightResult(nType int, battleInfo *BattleInfo) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	if len(self.Sql_Uid) <= 0 {
		return
	}

	for i, v := range self.FightList {
		data := v

		if data.Type != nType {
			continue
		}

		if data.FightId != battleInfo.Id {
			continue
		}

		player := GetPlayerMgr().GetPlayer(data.Attack.Uid, false)
		if player == nil {
			GetFightMgr().DelResult(battleInfo.Id)
			self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
			return
		}

		battleInfo.UserInfo[POS_ATTACK].Uid = data.Attack.Uid
		battleInfo.UserInfo[POS_ATTACK].Name = data.Attack.Uname
		battleInfo.UserInfo[POS_ATTACK].Icon = data.Attack.Iconid
		battleInfo.UserInfo[POS_ATTACK].Portrait = data.Attack.Portrait
		battleInfo.UserInfo[POS_ATTACK].UnionName = data.Attack.UnionName
		battleInfo.UserInfo[POS_ATTACK].Level = data.Attack.Level

		battleInfo.UserInfo[POS_DEFENCE].Uid = data.Defend.Uid
		battleInfo.UserInfo[POS_DEFENCE].Name = data.Defend.Uname
		battleInfo.UserInfo[POS_DEFENCE].Icon = data.Defend.Iconid
		battleInfo.UserInfo[POS_DEFENCE].Portrait = data.Defend.Portrait
		battleInfo.UserInfo[POS_DEFENCE].UnionName = data.Defend.UnionName
		battleInfo.UserInfo[POS_DEFENCE].Level = data.Defend.Level

		battleInfo.Type = BATTLE_TYPE_PVP
		battleInfo.Time = data.Time
		battleInfo.Random = data.Random

		result := 0
		if battleInfo.Result == 0 {
			result = 1
		}

		if data.Type == ARENA_FIGHT_TYPE_ENEMY {
			self.FightEnd(player, result, *battleInfo)
		} else if data.Type == ARENA_FIGHT_TYPE_FIGHT_BACK {
			self.FightBackEnd(player, result, *battleInfo)
		}

		GetFightMgr().DelResult(battleInfo.Id)
		self.FightList = append(self.FightList[0:i], self.FightList[i+1:]...)
		return
	}
}

// 斗技场结束
func (self *ArenaMgr) FightEnd(player *Player, result int, battleInfo BattleInfo) {
	modArena := player.GetModule("arena").(*ModArena)
	if modArena.enemy.Type != ARENA_FIGHT_TYPE_ENEMY {
		return
	}

	if modArena.enemy.Index < 0 {
		return
	}

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

	// 敌人信息
	fight := info.enemy[modArena.enemy.Index]
	enemyPoint := int64(ARENA_BASE_POINT)
	if fight.Enemy.Uid != 0 {
		// 获得敌人数据
		enemy, ok := self.Sql_Uid[fight.Enemy.Uid]
		if ok {
			enemyPoint = enemy.Point
		}
	}
	myPoint := info.Point

	msg := &S2C_ArenaEnd{}

	myInfo := GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_1)
	backRet := 0
	// 打成功
	if result == 1 {
		player.HandleTask(TASK_TYPE_JJC_SCORE, 3, 1, 0)
		// 我方胜利 计算我方加分 主体是我方
		addPoint := self.CountAddPoint(true, fight.Step+1, myPoint, enemyPoint)
		// 敌方减分 主体是敌方
		minPoint := self.CountMinPoint(true, fight.Step+1, myPoint, enemyPoint)
		// 通知客户端我方加多少分 敌方减多少
		msg.MyAddPoint = addPoint
		msg.EnemyAddPoint = minPoint

		// 添加自己的积分和战报
		msg.OldRank = info.Rank
		self.AddPoint(uid, info.Rank, addPoint)
		msg.MyPoint = info.Point

		info.AddFight(myInfo, fight.Enemy, 0, 1, int(addPoint), battleInfo)
		backRet = 0

		// 打的如果不是机器人
		if fight.Enemy.Uid != 0 {
			// 获得敌人数据
			enemy, ok := self.Sql_Uid[fight.Enemy.Uid]
			if !ok || enemy == nil {
				player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
				return
			}
			// 减分并增加战报
			self.MinPoint(fight.Enemy.Uid, enemy.Rank, -minPoint)
			enemy.AddFight(myInfo, fight.Enemy, 0, 0, int(-minPoint), battleInfo)

			msg.EnemyPoint = enemy.Point
		} else {
			msg.EnemyPoint = ARENA_BASE_POINT + minPoint
		}
		msg.NewRank = info.Rank
		config := GetCsvMgr().GetArenaParameterConfig(modArena.enemy.Index + 1)
		if nil != config {
			items := GetLootMgr().LootItem(config.Drop, player)
			out := player.AddObjectItemMap(items, "竞技场胜利奖励", 0, 0, 0)
			for _, v := range out {
				msg.Item = append(msg.Item, v)
			}
		}
		player.HandleTask(TASK_TYPE_ARENA_COUNT, 1, 0, 0)
		GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_ARENA_FIGHT, 1, int(fight.Enemy.Uid), int(minPoint), "竞技场战斗", 0, int(addPoint), player)
	} else {
		player.HandleTask(TASK_TYPE_JJC_SCORE, 1, 1, 0)
		// 失败 计算我方减分和敌方加分 主体是我方
		minPoint := self.CountMinPoint(false, fight.Step+1, myPoint, enemyPoint)
		// 主体是敌方
		addPoint := self.CountAddPoint(false, fight.Step+1, myPoint, enemyPoint)
		// 通知客户端我方减多少分 敌方加多少分
		msg.MyAddPoint = minPoint
		msg.EnemyAddPoint = addPoint

		// 添加自己的积分和战报
		msg.OldRank = info.Rank
		// 失败减分 并且增加战报
		self.MinPoint(uid, info.Rank, -minPoint)
		msg.MyPoint = info.Point
		info.AddFight(myInfo, fight.Enemy, 1, 1, int(-minPoint), battleInfo)
		// 打的不是机器人
		if fight.Enemy.Uid != 0 {
			// 获得数据
			enemy, ok := self.Sql_Uid[fight.Enemy.Uid]
			if !ok || enemy == nil {
				player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
				return
			}
			// 加分并增加战报
			self.AddPoint(fight.Enemy.Uid, enemy.Rank, addPoint)
			enemy.AddFight(myInfo, fight.Enemy, 1, 0, int(addPoint), battleInfo)
			msg.EnemyPoint = enemy.Point
		} else {
			msg.EnemyPoint = ARENA_BASE_POINT + addPoint
		}
		msg.NewRank = info.Rank
		player.HandleTask(TASK_TYPE_ARENA_COUNT, 0, 0, 0)
		backRet = 1
		GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_ARENA_FIGHT, 0, int(fight.Enemy.Uid), int(addPoint), "竞技场战斗", 0, int(minPoint), player)
	}

	player.GetModule("task").(*ModTask).SendUpdate()
	maxnum := vipcsv.ArenaFree[ARENA_TYPE_NOMAL]
	if info.Count > maxnum {
		cost := TARIFF_TYPE_ARENA_NORMAL

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
	msg.Cid = "armsarenafight"
	msg.RandNum = info.RandNum
	msg.FightID = battleInfo.Id
	//msg.BattleInfo = &battleInfo
	msg.Index = modArena.enemy.Index
	msg.Result = backRet
	msg.FightInfo[0] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_1)
	msg.FightInfo[1] = fight.Enemy
	player.Send(msg.Cid, msg)
	modArena.enemy.Index = -1
}

// 斗技场结束
func (self *ArenaMgr) FightBackEnd(player *Player, result int, battleInfo BattleInfo) {
	modArena := player.GetModule("arena").(*ModArena)
	if modArena.enemy.Type != ARENA_FIGHT_TYPE_FIGHT_BACK {
		return
	}

	uid := player.GetUid()
	info, ok := self.Sql_Uid[uid]
	if !ok {
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}

	fightID := modArena.enemy.Index

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
		player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
		return
	}
	// 不是失败
	if (fight.Result == 0 && fight.Side == 1) || (fight.Result == 1 && fight.Side == 0) {
		player.SendErrInfo("err", GetCsvMgr().GetText("不是失败"))
		return
	}

	var oldRecord BattleRecord
	value, flag, err := HGetRedisEx(`san_arenabattlerecord`, fightID, fmt.Sprintf("%d", fightID))
	if err != nil {
		player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
		return
	}
	if flag {
		err := json.Unmarshal([]byte(value), &oldRecord)
		if err != nil {
			player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
			return
		}
	}

	if oldRecord.Id == 0 {
		player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
		return
	}

	record := oldRecord

	vipcsv := GetCsvMgr().GetVipConfig(player.Sql_UserBase.Vip)
	if vipcsv == nil {
		player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_ERROR"))
		return
	}

	enemyPoint := int64(ARENA_BASE_POINT)
	if fight.Uid != 0 {
		// 获得敌人数据
		enemy, ok := self.Sql_Uid[fight.Uid]
		if ok {
			enemyPoint = enemy.Point
		}
	}

	myPoint := info.Point

	msg := &S2C_ArenaFightBackEnd{}
	// 打的次数+1
	//info.Count++
	myInfo := GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_1)
	backRet := 0
	// 打成功
	if result == 1 {
		player.HandleTask(TASK_TYPE_JJC_SCORE, 3, 1, 0)
		var fightInfo *JS_FightInfo = nil
		// 我方胜利 计算我方加分 主体是我方
		addPoint := self.CountAddPoint(true, 3, myPoint, enemyPoint)
		// 敌方减分 主体是敌方
		minPoint := self.CountMinPoint(true, 3, myPoint, enemyPoint)
		// 通知客户端我方加多少分 敌方减多少
		msg.MyAddPoint = addPoint
		msg.EnemyAddPoint = minPoint

		// 记录我方旧排名 加分后计算新排名和新分属
		msg.OldRank = info.Rank
		self.AddPoint(uid, info.Rank, addPoint)
		msg.MyPoint = info.Point
		if fight.Uid != 0 {
			// 获得敌人数据
			enemy, ok := self.Sql_Uid[fight.Uid]
			if !ok || enemy == nil {
				player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
				return
			}
			// 敌方减分并计算新排名新分数 加战报
			self.MinPoint(fight.Uid, enemy.Rank, -minPoint)
			enemy.AddFight(myInfo, enemy.format, 0, 0, int(-minPoint), battleInfo)
			fightInfo = enemy.format
			msg.EnemyPoint = enemy.Point
		} else {
			fightInfo = record.FightInfo[1]
			msg.EnemyPoint = ARENA_BASE_POINT + minPoint
		}
		msg.NewRank = info.Rank
		info.AddFight(myInfo, fightInfo, 0, 1, int(addPoint), battleInfo)
		msg.FightInfo[1] = fightInfo
		backRet = 0

		config := GetCsvMgr().GetArenaParameterConfig(3)
		if nil != config {
			items := GetLootMgr().LootItem(config.Drop, player)
			out := player.AddObjectItemMap(items, "竞技场胜利奖励", 0, 0, 0)
			for _, v := range out {
				msg.Item = append(msg.Item, v)
			}
		}
		player.HandleTask(TASK_TYPE_ARENA_COUNT, 1, 0, 0)
		GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_ARENA_FIGHT_BACK, 1, int(fight.Uid), int(minPoint), "竞技场反击战斗", 0, int(addPoint), player)
	} else {
		player.HandleTask(TASK_TYPE_JJC_SCORE, 1, 1, 0)
		// 失败 计算我方减分和敌方加分 主体是我方
		minPoint := self.CountMinPoint(false, 3, myPoint, enemyPoint)
		// 主体是敌方
		addPoint := self.CountAddPoint(false, 3, myPoint, enemyPoint)
		// 通知客户端我方减多少分 敌方加多少分
		msg.MyAddPoint = minPoint
		msg.EnemyAddPoint = addPoint

		// 记录我方旧排名 减分后计算新排名和新分属
		msg.OldRank = info.Rank
		self.MinPoint(uid, info.Rank, -minPoint)
		msg.MyPoint = info.Point

		var fightInfo *JS_FightInfo = nil
		if fight.Uid != 0 {
			// 获得敌人数据
			enemy, ok := self.Sql_Uid[fight.Uid]
			if !ok || enemy == nil {
				player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_MY_DATA_IS_EMPTY"))
				return
			}
			self.AddPoint(fight.Uid, enemy.Rank, addPoint)
			enemy.AddFight(myInfo, enemy.format, 1, 0, int(addPoint), battleInfo)
			fightInfo = enemy.format
			msg.EnemyPoint = enemy.Point
		} else {
			fightInfo = record.FightInfo[1]
			msg.EnemyPoint = ARENA_BASE_POINT + addPoint
		}
		msg.NewRank = info.Rank
		info.AddFight(myInfo, fightInfo, 1, 1, int(-minPoint), battleInfo)
		msg.FightInfo[1] = fightInfo
		player.HandleTask(TASK_TYPE_ARENA_COUNT, 0, 0, 0)
		backRet = 1
		GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_ARENA_FIGHT_BACK, 0, int(fight.Uid), int(addPoint), "竞技场反击战斗", 0, int(minPoint), player)
	}
	info.DeleteFight(fightID)
	player.GetModule("task").(*ModTask).SendUpdate()
	maxnum := vipcsv.ArenaFree[ARENA_TYPE_NOMAL]
	if info.Count > maxnum {
		cost := TARIFF_TYPE_ARENA_NORMAL

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
	msg.Cid = "armsarenafightback"
	msg.RandNum = info.RandNum
	msg.FightID = battleInfo.Id
	//msg.BattleInfo = &battleInfo
	msg.Index = modArena.enemy.Index
	msg.Result = backRet
	msg.FightInfo[0] = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_1)
	player.Send(msg.Cid, msg)

	modArena.enemy.Index = -1
}

func (self *ArenaMgr) AddFightList(nType int, attack *JS_FightInfo, defend *JS_FightInfo, time int64, random int64) int64 {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	// 打的次数+1
	uid := attack.Uid
	info, ok := self.Sql_Uid[uid]
	if !ok {
		return 0
	}
	info.Count++
	if len(self.FightList) >= ARENA_FIGHT_COUNT_MAX {
		return 0
	}

	ret := GetFightMgr().AddArenaFightID(attack, defend, random)
	self.FightList = append(self.FightList, &ArenaFightList{nType, ret.Id, random, time, attack, defend, 0})

	return ret.Id
}

func (self *ArenaMgr) CountAddPoint(win bool, index int64, myPoint int64, enemyPoint int64) int64 {
	config := GetCsvMgr().GetArenaParameterConfig(index)
	if nil == config {
		return 0
	}

	if config.Expected <= 0 {
		return 0
	}

	// 我方胜利敌方失败我方获得积分 = integral的值 *（winparameter的值/100 - 1/（1+10^((对手积分 - 我的积分）/expected的值）
	// 我方失败敌方胜利敌方获得积分 = integral的值 *（winparameter的值/100 - 1/（1+10^(( 我的积分-对手积分）/expected的值）
	min := float64(enemyPoint - myPoint)
	if !win {
		min = float64(myPoint - enemyPoint)
	}
	diff := min / float64(config.Expected)
	mult := float64(config.Winparameter)/100 - float64(1)/(math.Pow(10, diff)+float64(1))
	addCount := int64(float64(config.Integral) * mult)
	return addCount
}
func (self *ArenaMgr) CountMinPoint(win bool, index int64, myPoint int64, enemyPoint int64) int64 {
	config := GetCsvMgr().GetArenaParameterConfig(index)
	if nil == config {
		return 0
	}

	if config.Expected <= 0 {
		return 0
	}

	// 我方胜利敌方失败敌方减少积分 = integral的值 *（loseparameter的值/100 - 1/（1+10^((我的积分 - 对手积分）/expected的值）
	// 我方失败敌方胜利我方减少积分 = integral的值 *（loseparameter的值/100 - 1/（1+10^((对手积分 - 我的积分）/expected的值）
	min := float64(myPoint - enemyPoint)
	if !win {
		min = float64(enemyPoint - myPoint)
	}
	diff := min / float64(config.Expected)
	mult := float64(config.Loseparameter)/100 - float64(1)/(math.Pow(10, diff)+float64(1))
	addCount := int64(float64(config.Integral) * mult)
	return addCount
}

func (self *ArenaMgr) GetPlayerByFight(myfight int64, minFight int64, maxFight int64, needCount int, minlimit int64) []*JS_FightInfo {
	data := []*JS_FightInfo{}
	count := 0

	secondfight := int64(0)
	uid := int64(0)

	lv := 0
	for _, v := range self.Sql_Uid {
		if v.format == nil {
			continue
		}

		if v.Rank <= 0 || v.Point <= 0 {
			continue
		}

		if len(v.format.Defhero) < MAX_FIGHT_POS {
			continue
		}

		fight := v.format.Deffight
		if fight >= minFight && fight <= maxFight {
			data = append(data, v.format)
			count++
			if count >= needCount {
				if v.format != nil {
					for _, value := range v.format.Heroinfo {
						lv += value.Levels
					}
				}
				return data
			}
		} else {
			if fight < minFight && fight >= minlimit*myfight/10000 {
				if fight > secondfight {
					secondfight = fight
					uid = v.Uid
				}
			}
		}
	}

	if count <= 0 && uid > 0 {
		player, ok := self.Sql_Uid[uid]
		if ok {
			data = append(data, player.format)
			if player.format != nil {
				for _, value := range player.format.Heroinfo {
					lv += value.Levels
				}
			}
		}
	}

	return data
}

func (self *ArenaMgr) Rename(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.Sql_Uid[player.GetUid()]
	if ok {
		self.Sql_Uid[player.GetUid()].Name = player.GetName()
		self.Sql_Uid[player.GetUid()].format.Uname = player.GetName()
	}
}

func (self *ArenaMgr) Rehead(player *Player) {
	if player == nil {
		return
	}
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.Sql_Uid[player.Sql_UserBase.Uid]
	if ok {
		self.Sql_Uid[player.Sql_UserBase.Uid].format.Iconid = player.Sql_UserBase.IconId
		self.Sql_Uid[player.Sql_UserBase.Uid].format.Portrait = player.Sql_UserBase.Portrait
	}
}

func (self *ArenaMgr) AddFightListForBoss(nType int, attack *JS_FightInfo, defend *JS_FightInfo, time int64, random int64, bossId int) int64 {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	ret := GetFightMgr().AddArenaFightID(attack, defend, random)
	self.FightListBoss = append(self.FightListBoss, &ArenaFightList{nType, ret.Id, random, time, attack, defend, bossId})

	return ret.Id
}

func (self *ArenaMgr) AddFightListForCross(nType int, attack *JS_FightInfo, defend *JS_FightInfo, time int64, random int64) int64 {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	ret := GetFightMgr().AddArenaFightID(attack, defend, random)
	self.FightListCross = append(self.FightListCross, &ArenaFightList{nType, ret.Id, random, time, attack, defend, 0})

	return ret.Id
}

func (self *ArenaMgr) OnTimerCross() {

	now := TimeServer().Unix()
	nLen := len(self.FightListCross)
	for i := nLen - 1; i >= 0; i-- {
		data := self.FightListCross[i]
		if data.Time+300 <= now {
			GetFightMgr().DelResult(data.FightId)
			self.FightListCross = append(self.FightListCross[0:i], self.FightListCross[i+1:]...)
			continue
		}

		// 尝试获取战斗结果
		//LogDebug("尝试获取战斗结果")
		FightResult := GetFightMgr().GetResult(data.FightId)
		// 有结果,设置战斗时间
		if FightResult == nil {
			continue
		}

		player := GetPlayerMgr().GetPlayer(data.Attack.Uid, false)
		if player == nil {
			GetFightMgr().DelResult(FightResult.Id)
			self.FightListCross = append(self.FightListCross[0:i], self.FightListCross[i+1:]...)
			continue
		}

		battleInfo := BattleInfo{}
		battleInfo.Id = data.FightId
		attackHeroInfo := []*BattleHeroInfo{}
		for i, v := range FightResult.Info[POS_ATTACK] {
			level, star, skin, exclusiveLv := 0, 0, 0, 0
			if i < len(FightResult.Fight[POS_ATTACK].Heroinfo) {
				level = FightResult.Fight[POS_ATTACK].Heroinfo[i].Levels
				star = FightResult.Fight[POS_ATTACK].Heroinfo[i].Stars
				skin = FightResult.Fight[POS_ATTACK].Heroinfo[i].Skin
				exclusiveLv = FightResult.Fight[POS_ATTACK].Heroinfo[i].HeroExclusiveLv
			}
			attackHeroInfo = append(attackHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusiveLv, nil})
		}
		defendHeroInfo := []*BattleHeroInfo{}
		for i, v := range FightResult.Info[POS_DEFENCE] {
			level, star, skin, exclusiveLv := 0, 0, 0, 0
			if i < len(FightResult.Fight[POS_DEFENCE].Heroinfo) {
				level = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Levels
				star = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Stars
				skin = FightResult.Fight[POS_DEFENCE].Heroinfo[i].Skin
				exclusiveLv = FightResult.Fight[POS_DEFENCE].Heroinfo[i].HeroExclusiveLv
			}
			defendHeroInfo = append(defendHeroInfo, &BattleHeroInfo{v.Heroid, level, star, skin, v.Hp, v.Energy, v.Damage, v.TakeDamage, v.Healing, nil, exclusiveLv, nil})
		}

		battleInfo.UserInfo[POS_ATTACK] = &BattleUserInfo{
			data.Attack.Uid,
			data.Attack.Uname,
			data.Attack.Iconid,
			data.Attack.Portrait,
			data.Attack.UnionName,
			data.Attack.Level,
			attackHeroInfo}
		battleInfo.UserInfo[POS_DEFENCE] = &BattleUserInfo{
			data.Defend.Uid,
			data.Defend.Uname,
			data.Defend.Iconid,
			data.Defend.Portrait,
			data.Defend.UnionName,
			data.Defend.Level,
			defendHeroInfo}
		battleInfo.Type = BATTLE_TYPE_PVP
		battleInfo.Time = data.Time
		battleInfo.Random = data.Random
		if FightResult.Result == 1 {
			battleInfo.Result = 0
		} else {
			battleInfo.Result = 1
		}

		player.GetModule("crossarena").(*ModCrossArena).FightEndOK(player, battleInfo, data.Attack, data.Defend)

		GetFightMgr().DelResult(FightResult.Id)
		self.FightListCross = append(self.FightListCross[0:i], self.FightListCross[i+1:]...)
	}
}

func (self *ArenaMgr) ArenaFightResultByCross(nType int, battleInfo *BattleInfo) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	if len(self.Sql_Uid) <= 0 {
		return
	}

	for i, v := range self.FightListCross {
		data := v

		if data.Type != nType {
			continue
		}

		if data.FightId != battleInfo.Id {
			continue
		}

		player := GetPlayerMgr().GetPlayer(data.Attack.Uid, false)
		if player == nil {
			GetFightMgr().DelResult(battleInfo.Id)
			self.FightListCross = append(self.FightListCross[0:i], self.FightListCross[i+1:]...)
			return
		}

		battleInfo.UserInfo[POS_ATTACK].Uid = data.Attack.Uid
		battleInfo.UserInfo[POS_ATTACK].Name = data.Attack.Uname
		battleInfo.UserInfo[POS_ATTACK].Icon = data.Attack.Iconid
		battleInfo.UserInfo[POS_ATTACK].Portrait = data.Attack.Portrait
		battleInfo.UserInfo[POS_ATTACK].UnionName = data.Attack.UnionName
		battleInfo.UserInfo[POS_ATTACK].Level = data.Attack.Level

		battleInfo.UserInfo[POS_DEFENCE].Uid = data.Defend.Uid
		battleInfo.UserInfo[POS_DEFENCE].Name = data.Defend.Uname
		battleInfo.UserInfo[POS_DEFENCE].Icon = data.Defend.Iconid
		battleInfo.UserInfo[POS_DEFENCE].Portrait = data.Defend.Portrait
		battleInfo.UserInfo[POS_DEFENCE].UnionName = data.Defend.UnionName
		battleInfo.UserInfo[POS_DEFENCE].Level = data.Defend.Level

		battleInfo.Type = BATTLE_TYPE_PVP
		battleInfo.Time = data.Time
		battleInfo.Random = data.Random

		player.GetModule("crossarena").(*ModCrossArena).FightEndOK(player, *battleInfo, data.Attack, data.Defend)

		GetFightMgr().DelResult(battleInfo.Id)
		self.FightListCross = append(self.FightListCross[0:i], self.FightListCross[i+1:]...)
		return
	}
}
