package game

import (
	"fmt"
	"sort"
	"sync"
)

type JS_HeroStarsTop struct {
	Uid          int64  `json:"uid"`
	Uname        string `json:"uname"`
	Iconid       int    `json:"iconid"`
	Portrait     int    `json:"portrait"` // 边框  20190412 by zy
	Level        int    `json:"level"`
	Camp         int    `json:"camp"`
	Count        int    `json:"star"`
	Vip          int    `json:"vip"`
	LastRank     int    `json:"-"`
	UnionName    string `json:"union_name"`
	HeroStarTime int64  `json:"herostartime"`
}

type lstTopHeroStars []*JS_HeroStarsTop
type TopHeroStarsMgr struct {
	// 排行数据
	TopHeroStars    []*JS_HeroStarsTop // 排行榜
	TopHeroStarsCur map[int64]int
	// 版本号
	Topver        int
	Locker        *sync.RWMutex
	CampHeroStars [3][]*JS_HeroStarsTop
	CampLocker    *sync.RWMutex
}

func (self *TopHeroStarsMgr) initCampRank() {
	top := self.TopHeroStars
	for camp := 0; camp < 3; camp++ {
		self.CampHeroStars[camp] = make([]*JS_HeroStarsTop, 0)
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
		self.CampHeroStars[camp] = append(self.CampHeroStars[camp], &JS_HeroStarsTop{
			Uid:    topInfo.Uid,
			Uname:  topInfo.Uname,
			Iconid: topInfo.Iconid,
			Level:  topInfo.Level,
			Camp:   topInfo.Camp,
			Vip:    topInfo.Vip,
			Count:  topInfo.Count,
		})
	}
}

// 初始排行
func (self *TopHeroStarsMgr) GetData(unionName map[int64]string) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	topHeroStarsMgr.TopHeroStars = make([]*JS_HeroStarsTop, 0)
	var top4 JS_HeroStarsTop
	text := GetCsvMgr().GetText("STR_TOP_NULL")
	sql6 := fmt.Sprintf("SELECT a.uid, a.uname, a.iconid,a.portrait , a.level, a.camp, b.herototalstars, a.vip , '%s' as union_name,"+
		" b.herostartime FROM san_userbase  AS a INNER JOIN "+
		"(SELECT san_userhero2.uid, san_userhero2.herototalstars, san_userhero2.herostartime FROM san_userhero2 "+
		"WHERE san_userhero2.herototalstars > 0 ORDER BY herototalstars DESC, herostartime ASC LIMIT 200)"+
		" AS b ON a.uid = b.uid", text)
	res6 := GetServer().DBUser.GetAllDataEx(sql6, &top4)
	for i := 0; i < len(res6); i++ {
		data := res6[i].(*JS_HeroStarsTop)
		data.LastRank = i + 1
		// 更新军团名字
		v, ok := unionName[data.Uid]
		if ok {
			data.UnionName = v
		}
		topHeroStarsMgr.TopHeroStars = append(topHeroStarsMgr.TopHeroStars, data)
	}

	sort.Sort(lstTopHeroStars(self.TopHeroStars))
	for i := 0; i < len(self.TopHeroStars); i++ {
		self.TopHeroStars[i].LastRank = i + 1
		self.TopHeroStarsCur[self.TopHeroStars[i].Uid] = i + 1
	}

	self.Topver++
	self.ResetCampRank()
}

// 获取英雄星级排行榜
func (self *TopHeroStarsMgr) getCampRank(camp int) []*JS_HeroStarsTop {
	self.CampLocker.RLock()
	defer self.CampLocker.RUnlock()

	if camp < 1 || camp > 3 {
		return []*JS_HeroStarsTop{}
	}

	return self.CampHeroStars[camp-1]
}

// 获取英雄星级排行榜
func (self *TopHeroStarsMgr) GetTopHeroStars_L() []*JS_HeroStarsTop {
	return self.TopHeroStars
}

func (self *TopHeroStarsMgr) GetTopHeroStarsShow() ([]*JS_HeroStarsTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	if len(self.TopHeroStars) > MaxRankShowNum {
		return self.TopHeroStars[:MaxRankShowNum], self.Topver
	}
	return self.TopHeroStars, self.Topver
}

func (self *TopHeroStarsMgr) GetTopHeroStarsInfo() ([]*JS_HeroStarsTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.TopHeroStars) > MaxRankNum {
		return self.TopHeroStars[:MaxRankNum], self.Topver
	}
	return self.TopHeroStars, self.Topver
}

var topHeroStarsMgr *TopHeroStarsMgr = nil

func GetTopHeroStarsMgr() *TopHeroStarsMgr {
	if topHeroStarsMgr == nil {
		topHeroStarsMgr = new(TopHeroStarsMgr)
		topHeroStarsMgr.TopHeroStars = make([]*JS_HeroStarsTop, 0)
		topHeroStarsMgr.TopHeroStarsCur = make(map[int64]int, 0)
		topHeroStarsMgr.Locker = new(sync.RWMutex)
		topHeroStarsMgr.CampLocker = new(sync.RWMutex)
		topHeroStarsMgr.initCampRank()
	}

	return topHeroStarsMgr
}

//排行同步改名
func (self *TopHeroStarsMgr) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for i := 0; i < len(self.TopHeroStars); i++ {
		if self.TopHeroStars[i].Uid == player.Sql_UserBase.Uid {
			self.TopHeroStars[i].Uname = player.Sql_UserBase.UName
			break
		}
	}
}

//// 更新排行数据
//func (self *TopHeroStarsMgr) UpdateRank(count int, player *Player) {
//	self.Locker.Lock()
//
//	insert := true  // 是否插入新数据
//	change := false // 是否重新排序
//	for i := 0; i < len(self.TopHeroStars); i++ {
//		if self.TopHeroStars[i].Uid == player.Sql_UserBase.Uid {
//			self.TopHeroStars[i].Count = count
//			self.TopHeroStars[i].Level = player.Sql_UserBase.Level
//			insert = false
//			if i > 0 {
//				if self.TopHeroStars[i-1].Count < count {
//					change = true
//				}
//			}
//			break
//		}
//	}
//
//	if insert == true {
//		var data JS_HeroStarsTop
//		data.Uid = player.Sql_UserBase.Uid
//		data.Uname = player.Sql_UserBase.UName
//		data.Iconid = player.Sql_UserBase.IconId
//		data.Camp = player.Sql_UserBase.Camp
//		data.Vip = player.Sql_UserBase.Vip
//		data.Count = count
//		self.TopHeroStars = append(self.TopHeroStars, &data)
//	}
//
//	if change == true || insert == true {
//		sort.Sort(lstTopHeroStars(self.TopHeroStars))
//		for i := 0; i < len(self.TopHeroStars); i++ {
//			if self.TopHeroStars[i] != nil {
//				self.TopHeroStars[i].LastRank = i + 1
//			}
//		}
//		self.Topver++
//	}
//
//	if len(self.TopHeroStars) > 200 {
//		self.TopHeroStars = self.TopHeroStars[:200]
//	}
//
//	self.Locker.Unlock()
//
//	if change == true || insert == true {
//		self.UpdateCampRank(count, player)
//	}
//
//}
func (self *TopHeroStarsMgr) ResetCampRank() {
	self.CampLocker.Lock()
	defer self.CampLocker.Unlock()

	top := GetTopHeroStarsMgr().GetTopHeroStars_L()

	for camp := 0; camp < 3; camp++ {
		self.CampHeroStars[camp] = make([]*JS_HeroStarsTop, 0)
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
		self.CampHeroStars[camp] = append(self.CampHeroStars[camp], &JS_HeroStarsTop{
			Uid:    topInfo.Uid,
			Uname:  topInfo.Uname,
			Iconid: topInfo.Iconid,
			Level:  topInfo.Level,
			Camp:   topInfo.Camp,
			Vip:    topInfo.Vip,
			Count:  topInfo.Count,
		})
	}
}

//
//// 更新每个国家的排行
//func (self *TopHeroStarsMgr) UpdateCampRank(count int, player *Player) {
//	camp := player.GetCamp() - 1
//	if camp >= 0 && camp < CAMP_WU {
//		self.CampLocker.Lock()
//		defer self.CampLocker.Unlock()
//
//		insert := true  // 是否插入新数据
//		change := false // 是否重新排序
//		for i := 0; i < len(self.CampHeroStars[camp]); i++ {
//			if self.CampHeroStars[camp][i].Uid == player.Sql_UserBase.Uid {
//				self.CampHeroStars[camp][i].Count = count
//				self.CampHeroStars[camp][i].Level = player.Sql_UserBase.Level
//				insert = false
//				if i > 0 {
//					if self.CampHeroStars[camp][i-1].Count < count {
//						change = true
//					}
//				}
//				break
//			}
//		}
//
//		if insert == true {
//			var data JS_HeroStarsTop
//			data.Uid = player.Sql_UserBase.Uid
//			data.Uname = player.Sql_UserBase.UName
//			data.Iconid = player.Sql_UserBase.IconId
//			data.Camp = player.Sql_UserBase.Camp
//			data.Vip = player.Sql_UserBase.Vip
//			data.Count = count
//			self.CampHeroStars[camp] = append(self.CampHeroStars[camp], &data)
//		}
//
//		if change == true || insert == true {
//			sort.Sort(lstTopHeroStars(self.CampHeroStars[camp]))
//			for i := 0; i < len(self.CampHeroStars[camp]); i++ {
//				if self.CampHeroStars[camp][i] != nil {
//					self.CampHeroStars[camp][i].LastRank = i + 1
//				}
//			}
//		}
//	}
//}

func (self *TopHeroStarsMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	data, _ := self.TopHeroStarsCur[id]
	return data
}

////////////////////////////////////////////////////////////////////////////////
func (s lstTopHeroStars) Len() int      { return len(s) }
func (s lstTopHeroStars) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstTopHeroStars) Less(i, j int) bool {
	if s[i].Count > s[j].Count {
		return true
	}

	if s[i].Count < s[j].Count {
		return false
	}

	if s[i].HeroStarTime < s[j].HeroStarTime {
		return true
	}

	if s[i].HeroStarTime > s[j].HeroStarTime {
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
