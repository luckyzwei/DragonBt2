package game

import (
	"encoding/json"
	"sort"
	"sync"
	//"time"
)

type TopTowerMgr struct {
	Top    [TOWER_TYPE_MAX][]*Js_ActTop
	TopCur [TOWER_TYPE_MAX]map[int64]int
	TopOld [TOWER_TYPE_MAX]map[int64]int
	Topver [TOWER_TYPE_MAX]int
	Locker *sync.RWMutex
}

var toptower *TopTowerMgr = nil

func GetTopTowerMgr() *TopTowerMgr {
	if toptower == nil {
		toptower = new(TopTowerMgr)
		for i := 0; i < TOWER_TYPE_MAX; i++ {
			toptower.Top[i] = make([]*Js_ActTop, 0)
			toptower.TopCur[i] = make(map[int64]int, 0)
			toptower.TopOld[i] = make(map[int64]int, 0)
		}
		toptower.Locker = new(sync.RWMutex)
	}
	return toptower
}

// 钻石消耗
func (self *TopTowerMgr) GetData(unionName map[int64]string) {
	self.Locker.Lock()
	for i := 0; i < TOWER_TYPE_MAX; i++ {
		if self.Topver[i] > 0 {
			HF_DeepCopy(&self.TopOld[i], &self.TopCur[i])
			self.Top[i] = make([]*Js_ActTop, 0)
			self.TopCur[i] = make(map[int64]int)
		}
	}

	var top Js_ActTopLoad
	sql := "select a.uid,a.uname,a.iconid,a.portrait,a.level,a.camp, a.fight,a.vip,b.info,'0' as union_name " +
		"from san_userbase as a JOIN  san_kingtower as b " +
		"where a.fight > 0 and a.uid = b.uid"
	res := GetServer().DBUser.GetAllDataEx(sql, &top)
	for i := 0; i < len(res); i++ {
		data := res[i].(*Js_ActTopLoad)

		info := []*JS_Tower{}
		json.Unmarshal([]byte(data.Nums), &info)

		if len(info) != TOWER_TYPE_MAX {
			continue
		}

		for t := 0; t < TOWER_TYPE_MAX; t++ {
			temp := Js_ActTop{}
			temp.Uid = data.Uid
			temp.Uname = data.Uname
			temp.Iconid = data.Iconid
			temp.Portrait = data.Portrait
			temp.Level = data.Level
			temp.Camp = data.Camp
			temp.Fight = data.Fight
			temp.Vip = data.Vip
			temp.UnionName = data.UnionName
			temp.LastRank = data.LastRank
			temp.StartTime = info[t].MaxLevelTs
			temp.Num = int64(info[t].MaxLevel)
			if temp.Num <= 0 {
				continue
			}

			lastRank, ok := self.TopCur[t][temp.Uid]
			if ok {
				temp.LastRank = lastRank
			}
			// 更新军团名字
			v, ok := unionName[temp.Uid]
			if ok {
				temp.UnionName = v
			}
			self.Top[t] = append(self.Top[t], &temp)
			self.Topver[t]++
		}

	}

	for i := 0; i < TOWER_TYPE_MAX; i++ {
		sort.Sort(lstJsActTop(self.Top[i]))
		for t := 0; t < len(self.Top[i]); t++ {
			self.Top[i][t].LastRank = t + 1
			self.TopCur[i][self.Top[i][t].Uid] = t + 1
		}
	}
	self.Locker.Unlock()
}

func (self *TopTowerMgr) GetTopShow(topType int) ([]*Js_ActTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	nType := self.GetTopType(topType)

	if len(self.Top[nType]) > MaxRankShowNum {
		return self.Top[nType][:MaxRankShowNum], self.Topver[nType]
	}
	return self.Top[nType], self.Topver[nType]
}

//获取钻石消耗
func (self *TopTowerMgr) GetTop(topType int) ([]*Js_ActTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	nType := self.GetTopType(topType)

	if len(self.Top[nType]) > MaxRankNum {
		return self.Top[nType][:MaxRankNum], self.Topver[nType]
	}
	return self.Top[nType], self.Topver[nType]
}

func (self *TopTowerMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	nType := self.GetTopType(topType)
	data, _ := self.TopCur[nType][id]
	self.Locker.RUnlock()
	return data
}

func (self *TopTowerMgr) GetTopOldNum(topType int, id int64) int {
	self.Locker.RLock()
	nType := self.GetTopType(topType)
	data, _ := self.TopOld[nType][id]
	self.Locker.RUnlock()
	return data
}
func (self *TopTowerMgr) GetTopType(topType int) int {
	nType := 0
	switch topType {
	case TOP_RANK_TOWER:
		nType = TOWER_TYPE_0
	case TOP_RANK_TOWER1:
		nType = TOWER_TYPE_1
	case TOP_RANK_TOWER2:
		nType = TOWER_TYPE_2
	case TOP_RANK_TOWER3:
		nType = TOWER_TYPE_3
	case TOP_RANK_TOWER4:
		nType = TOWER_TYPE_4
	}
	return nType
}

//排行同步改名
func (self *TopTowerMgr) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for i := 0; i < TOWER_TYPE_MAX; i++ {
		for t := 0; t < len(self.Top[i]); t++ {
			if self.Top[i][t].Uid == player.Sql_UserBase.Uid {
				self.Top[i][t].Uname = player.Sql_UserBase.UName
				break
			}
		}
	}
}

// 更新排行数据
func (self *TopTowerMgr) updateRank(nType int, num int64, player *Player) {
	self.Locker.Lock()

	insert := true  // 是否插入新数据
	change := false // 是否重新排序
	for i := 0; i < len(self.Top[nType]); i++ {
		if self.Top[nType][i].Uid == player.Sql_UserBase.Uid {
			self.Top[nType][i].Level = player.Sql_UserBase.Level

			if self.Top[nType][i].Num < num {
				self.Top[nType][i].Num = num
				self.Top[nType][i].StartTime = TimeServer().Unix()
			}

			insert = false
			if i > 0 {
				if self.Top[nType][i-1].Num < self.Top[nType][i].Num {
					change = true
				}
			}
			break
		}
	}

	if insert == true {
		var data Js_ActTop
		data.Uid = player.Sql_UserBase.Uid
		data.Uname = player.Sql_UserBase.UName
		data.Iconid = player.Sql_UserBase.IconId
		data.Portrait = player.Sql_UserBase.Portrait
		data.Level = player.Sql_UserBase.Level
		data.Fight = player.Sql_UserBase.Fight
		data.Camp = player.Sql_UserBase.Camp
		data.Vip = player.Sql_UserBase.Vip
		data.StartTime = TimeServer().Unix()
		data.Num = num
		self.Top[nType] = append(self.Top[nType], &data)
		self.TopCur[nType][data.Uid] = 0
	}

	if change == true || insert == true {
		sort.Sort(lstJsActTop(self.Top[nType]))
		for i := 0; i < len(self.Top[nType]); i++ {
			if self.Top[nType][i] != nil {
				self.Top[nType][i].LastRank = i + 1
				self.TopCur[nType][self.Top[nType][i].Uid] = i + 1
			}
		}
		self.Topver[nType]++
	}

	if len(self.Top[nType]) > 200 {
		self.Top[nType] = self.Top[nType][:200]
	}

	self.Locker.Unlock()

}

//////////////////////////////////////////////////////////////////////////////////
//func (s lstTopTower) Len() int      { return len(s) }
//func (s lstTopTower) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
//func (s lstTopTower) Less(i, j int) bool {
//	if s[i].Num > s[j].Num {
//		return true
//	}
//
//	if s[i].Num < s[j].Num {
//		return false
//	}
//
//	if s[i].StartTime < s[j].StartTime {
//		return true
//	}
//
//	if s[i].StartTime > s[j].StartTime {
//		return false
//	}
//
//	if s[i].LastRank < s[j].LastRank {
//		return true
//	}
//
//	if s[i].LastRank > s[j].LastRank {
//		return false
//	}
//
//	if s[i].Uid < s[j].Uid {
//		return true
//	}
//
//	if s[i].Uid > s[j].Uid {
//		return false
//	}
//
//	return true
//}
