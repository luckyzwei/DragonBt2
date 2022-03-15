package game

import (
	"fmt"
	"sort"
	"sync"
)

type JS_HorseFightTop struct {
	Uid       int64  `json:"uid"`
	Uname     string `json:"uname"`
	Iconid    int    `json:"iconid"`
	Portrait  int    `json:"portrait"`
	Level     int    `json:"level"`
	Camp      int    `json:"camp"`
	Fight     int64  `json:"fight"`
	Vip       int    `json:"vip"`
	LastRank  int    `json:"-"`
	UnionName string `json:"union_name"`
}

type lstTopHorseFight []*JS_HorseFightTop
type TopHorseFightMgr struct {
	// 排行数据
	TopHorseFight    []*JS_HorseFightTop // 排行榜
	TopHorseFightCur map[int64]int
	// 版本号
	Topver         int
	Locker         *sync.RWMutex
	CampHorseFight [3][]*JS_HorseFightTop
	CampLocker     *sync.RWMutex
}

func (self *TopHorseFightMgr) initCampRank() {
	top := self.TopHorseFight
	for camp := 0; camp < 3; camp++ {
		self.CampHorseFight[camp] = make([]*JS_HorseFightTop, 0)
	}

	for index := range top {
		topInfo := top[index]
		if topInfo == nil {
			continue
		}

		if topInfo.Camp < 1 || topInfo.Camp > 3 {
			continue
		}

		camp := topInfo.Camp - 1
		self.CampHorseFight[camp] = append(self.CampHorseFight[camp], &JS_HorseFightTop{
			Uid:    topInfo.Uid,
			Uname:  topInfo.Uname,
			Iconid: topInfo.Iconid,
			Level:  topInfo.Level,
			Camp:   topInfo.Camp,
			Vip:    topInfo.Vip,
			Fight:  topInfo.Fight,
		})
	}
}

// 初始排行
func (self *TopHorseFightMgr) GetData(unionName map[int64]string) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	topHorseFightMgr.TopHorseFight = make([]*JS_HorseFightTop, 0)
	var top4 JS_HorseFightTop
	text := GetCsvMgr().GetText("STR_TOP_NULL")
	sql6 := fmt.Sprintf("SELECT a.uid, a.uname, a.iconid,a.portrait , a.level, a.camp, b.horsetotalfight, a.vip , '%s' as union_name FROM san_userbase  AS a INNER JOIN (SELECT san_userhorse.uid, san_userhorse.horsetotalfight FROM san_userhorse WHERE san_userhorse.horsetotalfight > 0 ORDER BY horsetotalfight DESC LIMIT 200) AS b ON a.uid = b.uid", text)
	res6 := GetServer().DBUser.GetAllDataEx(sql6, &top4)
	for i := 0; i < len(res6); i++ {
		data := res6[i].(*JS_HorseFightTop)
		data.LastRank = i + 1
		// 更新军团名字
		v, ok := unionName[data.Uid]
		if ok {
			data.UnionName = v
		}
		topHorseFightMgr.TopHorseFight = append(topHorseFightMgr.TopHorseFight, data)
	}

	sort.Sort(lstTopHorseFight(self.TopHorseFight))
	for i := 0; i < len(self.TopHorseFight); i++ {
		self.TopHorseFight[i].LastRank = i + 1
		self.TopHorseFightCur[self.TopHorseFight[i].Uid] = i + 1
	}

	self.Topver++
	self.ResetCampRank()
}

// 获取装备宝石排行榜
func (self *TopHorseFightMgr) getCampRank(camp int) []*JS_HorseFightTop {
	self.CampLocker.RLock()
	defer self.CampLocker.RUnlock()

	if camp < 1 || camp > 3 {
		return []*JS_HorseFightTop{}
	}

	return self.CampHorseFight[camp-1]
}

// 获取装备宝石排行榜
func (self *TopHorseFightMgr) GetTopHorseFight_L() []*JS_HorseFightTop {
	return self.TopHorseFight
}

func (self *TopHorseFightMgr) GetTopHorseFightShow() ([]*JS_HorseFightTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	if len(self.TopHorseFight) > MaxRankShowNum {
		return self.TopHorseFight[:MaxRankShowNum], self.Topver
	}
	return self.TopHorseFight, self.Topver
}

func (self *TopHorseFightMgr) GetTopHorseFightInfo() ([]*JS_HorseFightTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.TopHorseFight) > MaxRankNum {
		return self.TopHorseFight[:MaxRankNum], self.Topver
	}
	return self.TopHorseFight, self.Topver
}

var topHorseFightMgr *TopHorseFightMgr = nil

func GetTopHorseFightMgr() *TopHorseFightMgr {
	if topHorseFightMgr == nil {
		topHorseFightMgr = new(TopHorseFightMgr)
		topHorseFightMgr.TopHorseFight = make([]*JS_HorseFightTop, 0)
		topHorseFightMgr.TopHorseFightCur = make(map[int64]int, 0)
		topHorseFightMgr.Locker = new(sync.RWMutex)
		topHorseFightMgr.CampLocker = new(sync.RWMutex)
		topHorseFightMgr.initCampRank()
	}

	return topHorseFightMgr
}

//排行同步改名
func (self *TopHorseFightMgr) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for i := 0; i < len(self.TopHorseFight); i++ {
		if self.TopHorseFight[i].Uid == player.Sql_UserBase.Uid {
			self.TopHorseFight[i].Uname = player.Sql_UserBase.UName
			break
		}
	}
}

//
//// 更新排行数据
//func (self *TopHorseFightMgr) UpdateRank(count int64, player *Player) {
//	self.Locker.Lock()
//
//	insert := true  // 是否插入新数据
//	change := false // 是否重新排序
//	for i := 0; i < len(self.TopHorseFight); i++ {
//		if self.TopHorseFight[i].Uid == player.Sql_UserBase.Uid {
//			self.TopHorseFight[i].Fight = count
//			self.TopHorseFight[i].Level = player.Sql_UserBase.Level
//			insert = false
//			if i > 0 {
//				if self.TopHorseFight[i-1].Fight < count {
//					change = true
//				}
//			}
//			break
//		}
//	}
//
//	if insert == true {
//		var data JS_HorseFightTop
//		data.Uid = player.Sql_UserBase.Uid
//		data.Uname = player.Sql_UserBase.UName
//		data.Iconid = player.Sql_UserBase.IconId
//		data.Camp = player.Sql_UserBase.Camp
//		data.Vip = player.Sql_UserBase.Vip
//		data.Fight = count
//		self.TopHorseFight = append(self.TopHorseFight, &data)
//	}
//
//	if change == true || insert == true {
//		sort.Sort(lstTopHorseFight(self.TopHorseFight))
//		for i := 0; i < len(self.TopHorseFight); i++ {
//			if self.TopHorseFight[i] != nil {
//				self.TopHorseFight[i].LastRank = i + 1
//			}
//		}
//		self.Topver++
//	}
//
//	if len(self.TopHorseFight) > 200 {
//		self.TopHorseFight = self.TopHorseFight[:200]
//	}
//
//	self.Locker.Unlock()
//
//	if change == true || insert == true {
//		self.UpdateCampRank(count, player)
//	}
//
//}
func (self *TopHorseFightMgr) ResetCampRank() {
	self.CampLocker.Lock()
	defer self.CampLocker.Unlock()

	top := GetTopHorseFightMgr().GetTopHorseFight_L()

	for camp := 0; camp < 3; camp++ {
		self.CampHorseFight[camp] = make([]*JS_HorseFightTop, 0)
	}

	for index := range top {
		topInfo := top[index]
		if topInfo == nil {
			continue
		}

		if topInfo.Camp < 1 || topInfo.Camp > 3 {
			continue
		}

		camp := topInfo.Camp - 1
		self.CampHorseFight[camp] = append(self.CampHorseFight[camp], &JS_HorseFightTop{
			Uid:    topInfo.Uid,
			Uname:  topInfo.Uname,
			Iconid: topInfo.Iconid,
			Level:  topInfo.Level,
			Camp:   topInfo.Camp,
			Vip:    topInfo.Vip,
			Fight:  topInfo.Fight,
		})
	}
}

//
//// 更新每个国家的排行
//func (self *TopHorseFightMgr) UpdateCampRank(count int64, player *Player) {
//	camp := player.GetCamp() - 1
//	if camp >= 0 && camp < CAMP_WU {
//		self.CampLocker.Lock()
//		defer self.CampLocker.Unlock()
//
//		insert := true  // 是否插入新数据
//		change := false // 是否重新排序
//		for i := 0; i < len(self.CampHorseFight[camp]); i++ {
//			if self.CampHorseFight[camp][i].Uid == player.Sql_UserBase.Uid {
//				self.CampHorseFight[camp][i].Fight = count
//				self.CampHorseFight[camp][i].Level = player.Sql_UserBase.Level
//				insert = false
//				if i > 0 {
//					if self.CampHorseFight[camp][i-1].Fight < count {
//						change = true
//					}
//				}
//				break
//			}
//		}
//
//		if insert == true {
//			var data JS_HorseFightTop
//			data.Uid = player.Sql_UserBase.Uid
//			data.Uname = player.Sql_UserBase.UName
//			data.Iconid = player.Sql_UserBase.IconId
//			data.Camp = player.Sql_UserBase.Camp
//			data.Vip = player.Sql_UserBase.Vip
//			data.Fight = count
//			self.CampHorseFight[camp] = append(self.CampHorseFight[camp], &data)
//		}
//
//		if change == true || insert == true {
//			sort.Sort(lstTopHorseFight(self.CampHorseFight[camp]))
//			for i := 0; i < len(self.CampHorseFight[camp]); i++ {
//				if self.CampHorseFight[camp][i] != nil {
//					self.CampHorseFight[camp][i].LastRank = i + 1
//				}
//			}
//		}
//	}
//}

func (self *TopHorseFightMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	data, _ := self.TopHorseFightCur[id]
	return data
}

////////////////////////////////////////////////////////////////////////////////
func (s lstTopHorseFight) Len() int      { return len(s) }
func (s lstTopHorseFight) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstTopHorseFight) Less(i, j int) bool {
	if s[i].Fight > s[j].Fight {
		return true
	}

	if s[i].Fight < s[j].Fight {
		return false
	}

	if s[i].LastRank < s[j].LastRank {
		return true
	}

	if s[i].LastRank > s[j].LastRank {
		return false
	}

	if s[i].Uid < s[j].Uid {
		return true
	}

	if s[i].Uid > s[j].Uid {
		return false
	}

	return true
}
