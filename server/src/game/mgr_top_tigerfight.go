package game

import (
	"fmt"
	"sort"
	"sync"
)

type JS_TigerFightTop struct {
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

type lstTopTigerFight []*JS_TigerFightTop
type TopTigerFightMgr struct {
	// 排行数据
	TopTigerFight    []*JS_TigerFightTop // 排行榜
	TopTigerFightCur map[int64]int
	// 版本号
	Topver         int
	Locker         *sync.RWMutex
	CampTigerFight [3][]*JS_TigerFightTop
	CampLocker     *sync.RWMutex
}

func (self *TopTigerFightMgr) initCampRank() {
	top := self.TopTigerFight
	for camp := 0; camp < 3; camp++ {
		self.CampTigerFight[camp] = make([]*JS_TigerFightTop, 0)
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
		self.CampTigerFight[camp] = append(self.CampTigerFight[camp], &JS_TigerFightTop{
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
func (self *TopTigerFightMgr) GetData(unionName map[int64]string) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	topTigerFightMgr.TopTigerFight = make([]*JS_TigerFightTop, 0)
	var top4 JS_TigerFightTop
	text := GetCsvMgr().GetText("STR_TOP_NULL")
	sql6 := fmt.Sprintf("SELECT a.uid, a.uname, a.iconid,a.portrait , a.level, a.camp, b.tigertotalfight, a.vip , '%s' as union_name FROM san_userbase  AS a INNER JOIN (SELECT san_tiger.uid, san_tiger.tigertotalfight FROM san_tiger WHERE san_tiger.tigertotalfight > 0 ORDER BY tigertotalfight DESC LIMIT 200) AS b ON a.uid = b.uid", text)
	res6 := GetServer().DBUser.GetAllDataEx(sql6, &top4)
	for i := 0; i < len(res6); i++ {
		data := res6[i].(*JS_TigerFightTop)
		data.LastRank = i + 1
		// 更新军团名字
		v, ok := unionName[data.Uid]
		if ok {
			data.UnionName = v
		}
		topTigerFightMgr.TopTigerFight = append(topTigerFightMgr.TopTigerFight, data)
	}

	sort.Sort(lstTopTigerFight(self.TopTigerFight))
	for i := 0; i < len(self.TopTigerFight); i++ {
		self.TopTigerFight[i].LastRank = i + 1
		self.TopTigerFightCur[self.TopTigerFight[i].Uid] = i + 1
	}

	self.Topver++
	self.ResetCampRank()
}

// 获取装备宝石排行榜
func (self *TopTigerFightMgr) getCampRank(camp int) []*JS_TigerFightTop {
	self.CampLocker.RLock()
	defer self.CampLocker.RUnlock()

	if camp < 1 || camp > 3 {
		return []*JS_TigerFightTop{}
	}

	return self.CampTigerFight[camp-1]
}

// 获取装备宝石排行榜
func (self *TopTigerFightMgr) GetTopTigerFight_L() []*JS_TigerFightTop {
	return self.TopTigerFight
}

func (self *TopTigerFightMgr) GetTopTigerFightShow() ([]*JS_TigerFightTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	if len(self.TopTigerFight) > MaxRankShowNum {
		return self.TopTigerFight[:MaxRankShowNum], self.Topver
	}
	return self.TopTigerFight, self.Topver
}

func (self *TopTigerFightMgr) GetTopTigerFightInfo() ([]*JS_TigerFightTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.TopTigerFight) > MaxRankNum {
		return self.TopTigerFight[:MaxRankNum], self.Topver
	}
	return self.TopTigerFight, self.Topver
}

var topTigerFightMgr *TopTigerFightMgr = nil

func GetTopTigerFightMgr() *TopTigerFightMgr {
	if topTigerFightMgr == nil {
		topTigerFightMgr = new(TopTigerFightMgr)
		topTigerFightMgr.TopTigerFight = make([]*JS_TigerFightTop, 0)
		topTigerFightMgr.TopTigerFightCur = make(map[int64]int, 0)
		topTigerFightMgr.Locker = new(sync.RWMutex)
		topTigerFightMgr.CampLocker = new(sync.RWMutex)
		topTigerFightMgr.initCampRank()
	}

	return topTigerFightMgr
}

//排行同步改名
func (self *TopTigerFightMgr) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for i := 0; i < len(self.TopTigerFight); i++ {
		if self.TopTigerFight[i].Uid == player.Sql_UserBase.Uid {
			self.TopTigerFight[i].Uname = player.Sql_UserBase.UName
			break
		}
	}
}

//
//// 更新排行数据
//func (self *TopTigerFightMgr) UpdateRank(count int64, player *Player) {
//	self.Locker.Lock()
//
//	insert := true  // 是否插入新数据
//	change := false // 是否重新排序
//	for i := 0; i < len(self.TopTigerFight); i++ {
//		if self.TopTigerFight[i].Uid == player.Sql_UserBase.Uid {
//			self.TopTigerFight[i].Fight = count
//			self.TopTigerFight[i].Level = player.Sql_UserBase.Level
//			insert = false
//			if i > 0 {
//				if self.TopTigerFight[i-1].Fight < count {
//					change = true
//				}
//			}
//			break
//		}
//	}
//
//	if insert == true {
//		var data JS_TigerFightTop
//		data.Uid = player.Sql_UserBase.Uid
//		data.Uname = player.Sql_UserBase.UName
//		data.Iconid = player.Sql_UserBase.IconId
//		data.Camp = player.Sql_UserBase.Camp
//		data.Vip = player.Sql_UserBase.Vip
//		data.Fight = count
//		self.TopTigerFight = append(self.TopTigerFight, &data)
//	}
//
//	if change == true || insert == true {
//		sort.Sort(lstTopTigerFight(self.TopTigerFight))
//		for i := 0; i < len(self.TopTigerFight); i++ {
//			if self.TopTigerFight[i] != nil {
//				self.TopTigerFight[i].LastRank = i + 1
//			}
//		}
//		self.Topver++
//	}
//
//	if len(self.TopTigerFight) > 200 {
//		self.TopTigerFight = self.TopTigerFight[:200]
//	}
//
//	self.Locker.Unlock()
//
//	if change == true || insert == true {
//		self.UpdateCampRank(count, player)
//	}
//
//}
func (self *TopTigerFightMgr) ResetCampRank() {
	self.CampLocker.Lock()
	defer self.CampLocker.Unlock()

	top := GetTopTigerFightMgr().GetTopTigerFight_L()

	for camp := 0; camp < 3; camp++ {
		self.CampTigerFight[camp] = make([]*JS_TigerFightTop, 0)
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
		self.CampTigerFight[camp] = append(self.CampTigerFight[camp], &JS_TigerFightTop{
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
//func (self *TopTigerFightMgr) UpdateCampRank(count int64, player *Player) {
//	camp := player.GetCamp() - 1
//	if camp >= 0 && camp < CAMP_WU {
//		self.CampLocker.Lock()
//		defer self.CampLocker.Unlock()
//
//		insert := true  // 是否插入新数据
//		change := false // 是否重新排序
//		for i := 0; i < len(self.CampTigerFight[camp]); i++ {
//			if self.CampTigerFight[camp][i].Uid == player.Sql_UserBase.Uid {
//				self.CampTigerFight[camp][i].Fight = count
//				self.CampTigerFight[camp][i].Level = player.Sql_UserBase.Level
//				insert = false
//				if i > 0 {
//					if self.CampTigerFight[camp][i-1].Fight < count {
//						change = true
//					}
//				}
//				break
//			}
//		}
//
//		if insert == true {
//			var data JS_TigerFightTop
//			data.Uid = player.Sql_UserBase.Uid
//			data.Uname = player.Sql_UserBase.UName
//			data.Iconid = player.Sql_UserBase.IconId
//			data.Camp = player.Sql_UserBase.Camp
//			data.Vip = player.Sql_UserBase.Vip
//			data.Fight = count
//			self.CampTigerFight[camp] = append(self.CampTigerFight[camp], &data)
//		}
//
//		if change == true || insert == true {
//			sort.Sort(lstTopTigerFight(self.CampTigerFight[camp]))
//			for i := 0; i < len(self.CampTigerFight[camp]); i++ {
//				if self.CampTigerFight[camp][i] != nil {
//					self.CampTigerFight[camp][i].LastRank = i + 1
//				}
//			}
//		}
//	}
//}

func (self *TopTigerFightMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	data, _ := self.TopTigerFightCur[id]
	return data
}

////////////////////////////////////////////////////////////////////////////////
func (s lstTopTigerFight) Len() int      { return len(s) }
func (s lstTopTigerFight) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstTopTigerFight) Less(i, j int) bool {
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
