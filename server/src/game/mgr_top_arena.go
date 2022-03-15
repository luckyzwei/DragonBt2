package game

import (
	"fmt"
	"sort"
	"sync"
	//"time"
)

// 竞技场类型
const (
	ARENA_TOP_TYPE_NOMAL         = 0
	ARENA_TOP_TYPE_SPECIAL_RANK  = 1
	ARENA_TOP_TYPE_SPECIAL_POINT = 2
	ARENA_TOP_TYPE_HIGHEST       = 3
	ARENA_TOP_TYPE_MAX           = 4
)

type TopArenaMgr struct {
	Top    [ARENA_TOP_TYPE_MAX][]*Js_ActTop
	TopCur [ARENA_TOP_TYPE_MAX]map[int64]int
	TopOld [ARENA_TOP_TYPE_MAX]map[int64]int
	Topver [ARENA_TOP_TYPE_MAX]int
	Locker *sync.RWMutex
}

var topArenaMgr *TopArenaMgr = nil

func GetTopArenaMgr() *TopArenaMgr {
	if topArenaMgr == nil {
		topArenaMgr = new(TopArenaMgr)
		for i := 0; i < ARENA_TOP_TYPE_MAX; i++ {
			topArenaMgr.Top[i] = make([]*Js_ActTop, 0)
			topArenaMgr.TopCur[i] = make(map[int64]int)
			topArenaMgr.TopOld[i] = make(map[int64]int)
		}
		topArenaMgr.Locker = new(sync.RWMutex)
	}

	return topArenaMgr
}

// 初始斗技场排行
func (self *TopArenaMgr) GetData(unionName map[int64]string) {
	self.Locker.Lock()
	for i := 0; i < ARENA_TOP_TYPE_MAX; i++ {
		topArenaMgr.Top[i] = make([]*Js_ActTop, 0)
	}

	for i := 0; i < ARENA_TOP_TYPE_MAX; i++ {
		var top Js_ActTop
		if i == ARENA_TOP_TYPE_NOMAL {
			recIndex := ARENA_TYPE_NOMAL + 1
			text := GetCsvMgr().GetText("STR_TOP_NULL")
			sql := fmt.Sprintf("SELECT a.uid, a.uname, a.iconid,a.portrait , a.level, a.camp, a.fight, a.vip ,b.point, '%s' as union_name, b.starttime FROM san_userbase AS a INNER JOIN", text)
			sql += fmt.Sprintf(" (SELECT san_rankarena%d.uid,  san_rankarena%d.point,san_rankarena%d.rank, san_rankarena%d.starttime FROM san_rankarena%d WHERE san_rankarena%d.rank > 0 ORDER BY rank ASC LIMIT 200) AS b ",
				recIndex, recIndex, recIndex, recIndex, recIndex, recIndex)
			sql += "ON a.uid = b.uid"
			res := GetServer().DBUser.GetAllDataEx(sql, &top)
			for t := 0; t < len(res); t++ {
				data := res[t].(*Js_ActTop)
				data.LastRank = data.LastRank
				// 更新军团名字
				v, ok := unionName[data.Uid]
				if ok {
					data.UnionName = v
				}

				data.Fight = GetOfflineInfoMgr().GetTeamFight(data.Uid, TEAMTYPE_ARENA_2)

				topArenaMgr.Top[i] = append(topArenaMgr.Top[i], data)
			}
		} else if i == ARENA_TOP_TYPE_SPECIAL_RANK {
			recIndex := ARENA_TYPE_SPECIAL + 1
			text := GetCsvMgr().GetText("STR_TOP_NULL")
			sql := fmt.Sprintf("SELECT a.uid, a.uname, a.iconid,a.portrait , a.level, a.camp, a.fight,  a.vip ,b.rank, '%s' as union_name,  b.starttime FROM san_userbase AS a INNER JOIN", text)
			sql += fmt.Sprintf(" (SELECT san_rankarena%d.uid,  san_rankarena%d.point,san_rankarena%d.rank, san_rankarena%d.starttime FROM san_rankarena%d WHERE san_rankarena%d.rank > 0 ORDER BY rank ASC LIMIT 200) AS b ",
				recIndex, recIndex, recIndex, recIndex, recIndex, recIndex)
			sql += "ON a.uid = b.uid"
			res := GetServer().DBUser.GetAllDataEx(sql, &top)
			for t := 0; t < len(res); t++ {
				data := res[t].(*Js_ActTop)
				data.LastRank = data.LastRank
				// 更新军团名字
				v, ok := unionName[data.Uid]
				if ok {
					data.UnionName = v
				}
				data.Fight = GetOfflineInfoMgr().GetTeamFight(data.Uid, TEAMTYPE_ARENA_SPECIAL_4)

				topArenaMgr.Top[i] = append(topArenaMgr.Top[i], data)
			}
		} else if i == ARENA_TOP_TYPE_SPECIAL_POINT {
			recIndex := ARENA_TYPE_SPECIAL + 1
			text := GetCsvMgr().GetText("STR_TOP_NULL")
			sql := fmt.Sprintf("SELECT a.uid, a.uname, a.iconid,a.portrait , a.level, a.camp, a.fight,  a.vip ,b.point, '%s' as union_name,  b.starttime FROM san_userbase AS a INNER JOIN", text)
			sql += fmt.Sprintf(" (SELECT san_rankarena%d.uid,  san_rankarena%d.point,san_rankarena%d.rank, san_rankarena%d.starttime FROM san_rankarena%d WHERE san_rankarena%d.point > 0  ORDER BY point DESC LIMIT 200) AS b ",
				recIndex, recIndex, recIndex, recIndex, recIndex, recIndex)
			sql += "ON a.uid = b.uid"
			res := GetServer().DBUser.GetAllDataEx(sql, &top)
			for t := 0; t < len(res); t++ {
				data := res[t].(*Js_ActTop)
				data.LastRank = data.LastRank
				// 更新军团名字
				v, ok := unionName[data.Uid]
				if ok {
					data.UnionName = v
				}

				data.Fight = GetOfflineInfoMgr().GetTeamFight(data.Uid, TEAMTYPE_ARENA_SPECIAL_4)

				topArenaMgr.Top[i] = append(topArenaMgr.Top[i], data)
			}
		}

		self.Topver[i]++
	}

	for i := 0; i < ARENA_TOP_TYPE_MAX; i++ {
		if i == ARENA_TOP_TYPE_NOMAL {
			if len(self.Top[i]) <= 0 {
				continue
			}
			sort.Sort(lstJsArenaTop(self.Top[i]))
			for t := 0; t < len(self.Top[i]); t++ {
				self.Top[i][t].LastRank = t + 1
				self.TopCur[i][self.Top[i][t].Uid] = t + 1
			}
		} else if i == ARENA_TOP_TYPE_SPECIAL_RANK {
			//if len(self.Top[i]) <= 0 {
			//	continue
			//}
			//sort.Sort(lstJsArenaSpecialTop(self.Top[i]))
			data := self.Top[i]
			self.Top[i] = []*Js_ActTop{}
			config := GetCsvMgr().ArenaSpecialClassMap[1]
			for _, j := range config {
				if j.Ranking <= 0 {
					break
				}
				find := false
				dataLen := len(data)
				for f := dataLen - 1; f >= 0; f-- {
					if int64(j.Ranking) == data[f].Num {
						data[f].LastRank = j.Ranking
						self.Top[i] = append(self.Top[i], data[f])
						find = true
						self.TopCur[i][data[f].Uid] = j.Ranking
						data = append(data[:f], data[f+1:]...)
						break
					}
				}

				if !find {
					robot := GetArenaSpecialMgr().GetRobot(j.Class, j.Dan)
					temp := Js_ActTop{0, robot[ARENA_SPECIAL_TEAM_1].Uname,
						robot[ARENA_SPECIAL_TEAM_1].Iconid,
						robot[ARENA_SPECIAL_TEAM_1].Portrait,
						robot[ARENA_SPECIAL_TEAM_1].Level,
						robot[ARENA_SPECIAL_TEAM_1].Camp,
						(robot[ARENA_SPECIAL_TEAM_1].Deffight + robot[ARENA_SPECIAL_TEAM_2].Deffight + robot[ARENA_SPECIAL_TEAM_3].Deffight) / 100,
						robot[ARENA_SPECIAL_TEAM_1].Vip,
						int64(j.Ranking),
						robot[ARENA_SPECIAL_TEAM_1].UnionName,
						j.Ranking,
						0}
					self.Top[i] = append(self.Top[i], &temp)
				}
			}

		} else if i == ARENA_TOP_TYPE_SPECIAL_POINT {
			if len(self.Top[i]) <= 0 {
				continue
			}
			sort.Sort(lstJsActTop(self.Top[i]))
			for t := 0; t < len(self.Top[i]); t++ {
				self.Top[i][t].LastRank = t + 1
				self.TopCur[i][self.Top[i][t].Uid] = t + 1
			}
		}
	}
	self.Locker.Unlock()
}

//排行同步改名
func (self *TopArenaMgr) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for i := 0; i < ARENA_TOP_TYPE_MAX; i++ {
		for t := 0; t < len(self.Top[i]); t++ {
			if self.Top[i][t].Uid == player.Sql_UserBase.Uid {
				self.Top[i][t].Uname = player.Sql_UserBase.UName
				break
			}
		}
	}
}

//排行同步改头像
func (self *TopArenaMgr) Rehead(player *Player) {
	if player == nil {
		return
	}
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for i := 0; i < ARENA_TOP_TYPE_MAX; i++ {
		for t := 0; t < len(self.Top[i]); t++ {
			if self.Top[i][t].Uid == player.Sql_UserBase.Uid {
				self.Top[i][t].Iconid = player.Sql_UserBase.IconId
				self.Top[i][t].Portrait = player.Sql_UserBase.Portrait
				break
			}
		}
	}
}

func (self *TopArenaMgr) GetTopType(topType int) int {
	nType := 0
	switch topType {
	case TOP_RANK_ARENA_NORMAL:
		nType = ARENA_TOP_TYPE_NOMAL
	case TOP_RANK_ARENA_SPECIAL_RANK:
		nType = ARENA_TOP_TYPE_SPECIAL_RANK
	case TOP_RANK_ARENA_SPECIAL_POINT:
		nType = ARENA_TOP_TYPE_SPECIAL_POINT
	case TOP_RANK_ARENA_HIGHEST:
		nType = ARENA_TOP_TYPE_HIGHEST
	}
	return nType
}

func (self *TopArenaMgr) GetTopShow(topType int) ([]*Js_ActTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	nType := self.GetTopType(topType)
	if len(self.Top[nType]) > MaxRankShowNum {
		return self.Top[nType][:MaxRankShowNum], self.Topver[nType]
	}
	return self.Top[nType], self.Topver[nType]
}

//获取钻石消耗
func (self *TopArenaMgr) GetTop(topType int) ([]*Js_ActTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	nType := self.GetTopType(topType)
	if len(self.Top[nType]) > MaxRankNum {
		return self.Top[nType][:MaxRankNum], self.Topver[nType]
	}
	return self.Top[nType], self.Topver[nType]
}

func (self *TopArenaMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	nType := self.GetTopType(topType)
	data, _ := self.TopCur[nType][id]
	self.Locker.RUnlock()
	return data
}

func (self *TopArenaMgr) GetTopOldNum(topType int, id int64) int {
	self.Locker.RLock()
	nType := self.GetTopType(topType)
	data, _ := self.TopOld[nType][id]
	self.Locker.RUnlock()
	return data
}

// 更新排行数据
func (self *TopArenaMgr) UpdateArenaRank(count int64, uid int64) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	player := GetPlayerMgr().GetPlayer(uid, false)
	topType := ARENA_TOP_TYPE_NOMAL
	insert := true  // 是否插入新数据
	change := false // 是否重新排序
	for i := 0; i < len(self.Top[topType]); i++ {
		if self.Top[topType][i].Uid == uid {
			self.Top[topType][i].Num = count
			if player != nil {
				self.Top[topType][i].Level = player.Sql_UserBase.Level
				if nil != player {
					self.Top[topType][i].Fight = GetOfflineInfoMgr().GetTeamFight(player.GetUid(), TEAMTYPE_ARENA_2)
					self.Top[topType][i].Iconid = player.Sql_UserBase.IconId
					self.Top[topType][i].Portrait = player.Sql_UserBase.Portrait
					self.Top[topType][i].Level = player.Sql_UserBase.Level
					self.Top[topType][i].Vip = player.Sql_UserBase.Vip
				}
			}

			insert = false
			if i > 0 {
				if self.Top[topType][i-1].Num <= count {
					change = true
				}
			}
			if i < len(self.Top[topType])-1 {
				if self.Top[topType][i+1].Num >= count {
					change = true
				}
			}
			break
		}
	}

	if insert == true {
		if player == nil {
			player = GetPlayerMgr().GetPlayer(uid, true)
		}

		var data Js_ActTop
		data.Uid = player.Sql_UserBase.Uid
		data.Uname = player.Sql_UserBase.UName
		data.Iconid = player.Sql_UserBase.IconId
		data.Portrait = player.Sql_UserBase.Portrait
		data.Level = player.Sql_UserBase.Level
		data.Camp = player.Sql_UserBase.Camp
		data.Vip = player.Sql_UserBase.Vip
		data.UnionName = player.GetUnionName()
		data.Fight = player.Sql_UserBase.Fight
		if nil != player {
			data.Fight = GetOfflineInfoMgr().GetTeamFight(player.GetUid(), TEAMTYPE_ARENA_2)
		}
		data.StartTime = TimeServer().Unix()
		data.Num = count
		self.Top[topType] = append(self.Top[topType], &data)
		self.TopCur[topType][data.Uid] = 0
	}

	if change == true || insert == true {
		sort.Sort(lstJsArenaTop(self.Top[topType]))
		for i := 0; i < len(self.Top[topType]); i++ {
			if self.Top[topType][i] != nil {
				self.Top[topType][i].LastRank = i + 1
				self.TopCur[topType][self.Top[topType][i].Uid] = i + 1
			}
		}
		self.Topver[topType]++
	}

	if len(self.Top[topType]) > 200 {
		self.Top[topType] = self.Top[topType][:200]
	}
}

// 竞技场特殊排序规则
type lstJsArenaTop []*Js_ActTop

func (s lstJsArenaTop) Len() int      { return len(s) }
func (s lstJsArenaTop) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstJsArenaTop) Less(i, j int) bool {
	playeri, oki := GetArenaMgr().Sql_Uid[s[i].Uid]
	if !oki {
		fmt.Println(playeri.Name)
		return false
	}

	playerj, okj := GetArenaMgr().Sql_Uid[s[j].Uid]
	if !okj {
		fmt.Println(playerj.Name)
		return false
	}

	if playeri.Rank < playerj.Rank {
		if playeri.Point < playerj.Point {
			LogDebug("lstJsArenaTop ERRER !!!!!!!!!! playeri.Rank < playerj.Rank playeri.Point < playerj.Point")
		}
		return true
	}

	return false
}

type lstJsArenaSpecialTop []*Js_ActTop

func (s lstJsArenaSpecialTop) Len() int      { return len(s) }
func (s lstJsArenaSpecialTop) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstJsArenaSpecialTop) Less(i, j int) bool {
	if s[i].Num == 0 {
		return false
	}

	if s[i].Num < s[j].Num { // 由小到大
		return true
	}

	if s[i].Num > s[j].Num {
		return false
	}

	return false
}

// 更新排行数据
func (self *TopArenaMgr) UpdateFight(topType int, player *Player) {
	self.Locker.Lock()
	for i := 0; i < len(self.Top[topType]); i++ {
		if self.Top[topType][i].Uid == player.Sql_UserBase.Uid {
			self.Top[topType][i].Level = player.Sql_UserBase.Level
			if nil != player {
				fight := int64(0)
				if topType == ARENA_TOP_TYPE_NOMAL {
					fight = GetOfflineInfoMgr().GetTeamFight(player.GetUid(), TEAMTYPE_ARENA_2)
				} else if topType == ARENA_TOP_TYPE_SPECIAL_RANK || topType == ARENA_TOP_TYPE_SPECIAL_POINT {
					fight = GetOfflineInfoMgr().GetTeamFight(player.GetUid(), TEAMTYPE_ARENA_SPECIAL_4)
				}
				self.Top[topType][i].Fight = fight
				self.Top[topType][i].Iconid = player.Sql_UserBase.IconId
				self.Top[topType][i].Portrait = player.Sql_UserBase.Portrait
				self.Top[topType][i].Level = player.Sql_UserBase.Level
				self.Top[topType][i].Vip = player.Sql_UserBase.Vip
			}
			break
		}
	}
	self.Locker.Unlock()
}

// 更新排行数据
func (self *TopArenaMgr) UpdateArenaSpecialRank(count1 int64, player1 *Player, count2 int64, uid int64) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	topType := ARENA_TOP_TYPE_SPECIAL_RANK
	insert := true  // 是否插入新数据
	change := false // 是否重新排序

	insert2 := false // 是否插入新数据
	if uid != 0 {
		insert2 = true
		for i := 0; i < len(self.Top[topType]); i++ {
			if self.Top[topType][i].Uid == uid {
				insert2 = false
				self.Top[topType][i].Num = count2
				player2 := GetPlayerMgr().GetPlayer(uid, false)
				if nil != player2 {
					self.Top[topType][i].Iconid = player2.Sql_UserBase.IconId
					self.Top[topType][i].Portrait = player2.Sql_UserBase.Portrait
					self.Top[topType][i].Level = player2.Sql_UserBase.Level
					self.Top[topType][i].Vip = player2.Sql_UserBase.Vip
				}

				if i > 0 {
					if self.Top[topType][i-1].Num <= count2 {
						change = true
					}
				}
				if i < len(self.Top[topType])-1 {
					if self.Top[topType][i+1].Num >= count2 {
						change = true
					}
				}
				break
			}
		}

		if insert2 {
			player2 := GetPlayerMgr().GetPlayer(uid, true)
			if player2 != nil {
				var data Js_ActTop
				data.Uid = player2.Sql_UserBase.Uid
				data.Uname = player2.Sql_UserBase.UName
				data.Iconid = player2.Sql_UserBase.IconId
				data.Portrait = player2.Sql_UserBase.Portrait
				data.Level = player2.Sql_UserBase.Level
				data.Camp = player2.Sql_UserBase.Camp
				data.Vip = player2.Sql_UserBase.Vip
				data.UnionName = player2.GetUnionName()
				data.Fight = player2.Sql_UserBase.Fight
				if nil != player2 {
					data.Fight = GetOfflineInfoMgr().GetTeamFight(player2.GetUid(), TEAMTYPE_ARENA_SPECIAL_4)
				}
				data.StartTime = TimeServer().Unix()
				data.Num = count2
				self.Top[topType] = append(self.Top[topType], &data)
				self.TopCur[topType][data.Uid] = 0
			}
		}
	} else {
		for i := 0; i < len(self.Top[topType]); i++ {
			if self.Top[topType][i].Num == count1 {
				config := GetCsvMgr().GetArenaSpecialClassConfigByID(int(count2))
				if nil != config {
					robot := GetArenaSpecialMgr().GetRobot(config.Class, config.Dan)
					self.Top[topType][i].Uname = robot[ARENA_SPECIAL_TEAM_1].Uname
					self.Top[topType][i].Iconid = robot[ARENA_SPECIAL_TEAM_1].Iconid
					self.Top[topType][i].Portrait = robot[ARENA_SPECIAL_TEAM_1].Portrait
					self.Top[topType][i].Level = robot[ARENA_SPECIAL_TEAM_1].Level
					self.Top[topType][i].Camp = robot[ARENA_SPECIAL_TEAM_1].Camp
					self.Top[topType][i].Fight = (robot[ARENA_SPECIAL_TEAM_1].Deffight + robot[ARENA_SPECIAL_TEAM_2].Deffight + robot[ARENA_SPECIAL_TEAM_3].Deffight) / 100
					self.Top[topType][i].Vip = robot[ARENA_SPECIAL_TEAM_1].Vip
					self.Top[topType][i].UnionName = robot[ARENA_SPECIAL_TEAM_1].UnionName
					self.Top[topType][i].Num = count2
				}

				break
			}
		}
		change = true
	}

	for i := 0; i < len(self.Top[topType]); i++ {
		if self.Top[topType][i].Uid == player1.Sql_UserBase.Uid {
			self.Top[topType][i].Num = count1
			self.Top[topType][i].Level = player1.Sql_UserBase.Level
			if nil != player1 {
				self.Top[topType][i].Fight = GetOfflineInfoMgr().GetTeamFight(player1.GetUid(), TEAMTYPE_ARENA_SPECIAL_4)
				self.Top[topType][i].Iconid = player1.Sql_UserBase.IconId
				self.Top[topType][i].Portrait = player1.Sql_UserBase.Portrait
				self.Top[topType][i].Level = player1.Sql_UserBase.Level
				self.Top[topType][i].Vip = player1.Sql_UserBase.Vip
			}
			insert = false
			if i > 0 {
				if self.Top[topType][i-1].Num <= count1 {
					change = true
				}
			}
			if i < len(self.Top[topType])-1 {
				if self.Top[topType][i+1].Num >= count1 {
					change = true
				}
			}
			break
		}
	}

	if insert == true {
		var data Js_ActTop
		data.Uid = player1.Sql_UserBase.Uid
		data.Uname = player1.Sql_UserBase.UName
		data.Iconid = player1.Sql_UserBase.IconId
		data.Portrait = player1.Sql_UserBase.Portrait
		data.Level = player1.Sql_UserBase.Level
		data.Camp = player1.Sql_UserBase.Camp
		data.Vip = player1.Sql_UserBase.Vip
		data.UnionName = player1.GetUnionName()
		data.Fight = player1.Sql_UserBase.Fight
		if nil != player1 {
			data.Fight = GetOfflineInfoMgr().GetTeamFight(player1.GetUid(), TEAMTYPE_ARENA_SPECIAL_4)
		}
		data.StartTime = TimeServer().Unix()
		data.Num = count1
		self.Top[topType] = append(self.Top[topType], &data)
		self.TopCur[topType][data.Uid] = 0
	}

	if change == true || insert == true || insert2 == true {
		sort.Sort(lstJsArenaSpecialTop(self.Top[topType]))
		nLen := len(self.Top[topType])
		for i := nLen - 1; i >= 0; i-- {
			if self.Top[topType][i].Num == 0 {
				self.Top[topType] = append(self.Top[topType][:i], self.Top[topType][i+1:]...)
			}
		}

		for i := 0; i < len(self.Top[topType]); i++ {
			if self.Top[topType][i] != nil {
				self.Top[topType][i].LastRank = i + 1
				self.TopCur[topType][self.Top[topType][i].Uid] = i + 1
			}
		}
		self.Topver[topType]++
	}

	if len(self.Top[topType]) > 200 {
		self.Top[topType] = self.Top[topType][:200]
	}
}

// 更新排行数据
func (self *TopArenaMgr) UpdateArenaSpecialPoint(count int64, uid int64) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	player := GetPlayerMgr().GetPlayer(uid, false)
	topType := ARENA_TOP_TYPE_SPECIAL_POINT
	insert := true  // 是否插入新数据
	change := false // 是否重新排序
	for i := 0; i < len(self.Top[topType]); i++ {
		if self.Top[topType][i].Uid == uid {
			self.Top[topType][i].Num = count
			if player != nil {
				self.Top[topType][i].Level = player.Sql_UserBase.Level
				if nil != player {
					self.Top[topType][i].Fight = GetOfflineInfoMgr().GetTeamFight(player.GetUid(), TEAMTYPE_ARENA_SPECIAL_4)
					self.Top[topType][i].Iconid = player.Sql_UserBase.IconId
					self.Top[topType][i].Portrait = player.Sql_UserBase.Portrait
					self.Top[topType][i].Level = player.Sql_UserBase.Level
					self.Top[topType][i].Vip = player.Sql_UserBase.Vip
				}
			}

			insert = false
			if i > 0 {
				if self.Top[topType][i-1].Num <= count {
					change = true
				}
			}
			if i < len(self.Top[topType])-1 {
				if self.Top[topType][i+1].Num >= count {
					change = true
				}
			}
			break
		}
	}

	if insert == true {
		if player == nil {
			return
		}
		var data Js_ActTop
		data.Uid = player.Sql_UserBase.Uid
		data.Uname = player.Sql_UserBase.UName
		data.Iconid = player.Sql_UserBase.IconId
		data.Portrait = player.Sql_UserBase.Portrait
		data.Level = player.Sql_UserBase.Level
		data.Camp = player.Sql_UserBase.Camp
		data.Vip = player.Sql_UserBase.Vip
		data.UnionName = player.GetUnionName()
		data.Fight = player.Sql_UserBase.Fight
		if nil != player {
			data.Fight = GetOfflineInfoMgr().GetTeamFight(player.GetUid(), TEAMTYPE_ARENA_SPECIAL_4)
		}
		data.StartTime = TimeServer().Unix()
		data.Num = count
		self.Top[topType] = append(self.Top[topType], &data)
	}

	if change == true || insert == true {
		sort.Sort(lstJsActTop(self.Top[topType]))
		for i := 0; i < len(self.Top[topType]); i++ {
			if self.Top[topType][i] != nil {
				self.Top[topType][i].LastRank = i + 1
			}
		}
		self.Topver[topType]++
	}

	if len(self.Top[topType]) > 200 {
		self.Top[topType] = self.Top[topType][:200]
	}
}
