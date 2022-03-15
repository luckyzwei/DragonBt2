package game

import (
	"fmt"
	"sort"
	"sync"
)

type JS_EquipGemTop struct {
	Uid       int64  `json:"uid"`
	Uname     string `json:"uname"`
	Iconid    int    `json:"iconid"`
	Portrait  int    `json:"portrait"`
	Level     int    `json:"level"`
	Camp      int    `json:"camp"`
	Count     int    `json:"num"`
	Vip       int    `json:"vip"`
	LastRank  int    `json:"-"`
	UnionName string `json:"union_name"`
	StartTime int64  `json:"starttime"`
}

type lstTopEquipGem []*JS_EquipGemTop
type TopEquipGemMgr struct {
	// 排行数据
	TopEquipGem    []*JS_EquipGemTop // 排行榜
	TopEquipGemCur map[int64]int
	// 版本号
	Topver       int
	Locker       *sync.RWMutex
	CampEquipGem [3][]*JS_EquipGemTop
	CampLocker   *sync.RWMutex
}

func (self *TopEquipGemMgr) initCampRank() {
	top := self.TopEquipGem
	for camp := 0; camp < 3; camp++ {
		self.CampEquipGem[camp] = make([]*JS_EquipGemTop, 0)
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
		self.CampEquipGem[camp] = append(self.CampEquipGem[camp], &JS_EquipGemTop{
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
func (self *TopEquipGemMgr) GetData(unionName map[int64]string) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	topEquipGemMgr.TopEquipGem = make([]*JS_EquipGemTop, 0)
	var top4 JS_EquipGemTop
	text := GetCsvMgr().GetText("STR_TOP_NULL")
	sql6 := fmt.Sprintf("SELECT a.uid, a.uname, a.iconid,a.portrait , a.level, a.camp, b.totalgemlevel, a.vip , "+
		"'%s' as union_name, b.starttime FROM san_userbase  AS a INNER JOIN (SELECT san_userequip.uid, san_userequip.totalgemlevel, san_userequip.starttime FROM san_userequip "+
		"WHERE san_userequip.totalgemlevel > 0 ORDER BY totalgemlevel DESC , starttime ASC LIMIT 200) AS b ON a.uid = b.uid", text)
	res6 := GetServer().DBUser.GetAllDataEx(sql6, &top4)
	for i := 0; i < len(res6); i++ {
		data := res6[i].(*JS_EquipGemTop)
		data.LastRank = i + 1
		// 更新军团名字
		v, ok := unionName[data.Uid]
		if ok {
			data.UnionName = v
		}
		topEquipGemMgr.TopEquipGem = append(topEquipGemMgr.TopEquipGem, data)
	}

	sort.Sort(lstTopEquipGem(self.TopEquipGem))
	for i := 0; i < len(self.TopEquipGem); i++ {
		self.TopEquipGem[i].LastRank = i + 1
		self.TopEquipGemCur[self.TopEquipGem[i].Uid] = i + 1
	}

	self.Topver++
	self.ResetCampRank()
}

// 获取装备宝石排行榜
func (self *TopEquipGemMgr) getCampRank(camp int) []*JS_EquipGemTop {
	self.CampLocker.RLock()
	defer self.CampLocker.RUnlock()

	if camp < 1 || camp > 3 {
		return []*JS_EquipGemTop{}
	}

	return self.CampEquipGem[camp-1]
}

// 获取装备宝石排行榜
func (self *TopEquipGemMgr) GetTopEquipGem_L() []*JS_EquipGemTop {
	return self.TopEquipGem
}

func (self *TopEquipGemMgr) GetTopEquipGemShow() ([]*JS_EquipGemTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	if len(self.TopEquipGem) > MaxRankShowNum {
		return self.TopEquipGem[:MaxRankShowNum], self.Topver
	}
	return self.TopEquipGem, self.Topver
}

func (self *TopEquipGemMgr) GetTopEquipGemInfo() ([]*JS_EquipGemTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.TopEquipGem) > MaxRankNum {
		return self.TopEquipGem[:MaxRankNum], self.Topver
	}
	return self.TopEquipGem, self.Topver
}

var topEquipGemMgr *TopEquipGemMgr = nil

func GetTopEquipGemMgr() *TopEquipGemMgr {
	if topEquipGemMgr == nil {
		topEquipGemMgr = new(TopEquipGemMgr)
		topEquipGemMgr.TopEquipGem = make([]*JS_EquipGemTop, 0)
		topEquipGemMgr.TopEquipGemCur = make(map[int64]int, 0)
		topEquipGemMgr.Locker = new(sync.RWMutex)
		topEquipGemMgr.CampLocker = new(sync.RWMutex)
		topEquipGemMgr.initCampRank()
	}

	return topEquipGemMgr
}

//排行同步改名
func (self *TopEquipGemMgr) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for i := 0; i < len(self.TopEquipGem); i++ {
		if self.TopEquipGem[i].Uid == player.Sql_UserBase.Uid {
			self.TopEquipGem[i].Uname = player.Sql_UserBase.UName
			break
		}
	}
}

//// 更新排行数据
//func (self *TopEquipGemMgr) UpdateRank(count int, player *Player) {
//	self.Locker.Lock()
//
//	insert := true  // 是否插入新数据
//	change := false // 是否重新排序
//	for i := 0; i < len(self.TopEquipGem); i++ {
//		if self.TopEquipGem[i].Uid == player.Sql_UserBase.Uid {
//			self.TopEquipGem[i].Count = count
//			self.TopEquipGem[i].Level = player.Sql_UserBase.Level
//			insert = false
//			if i > 0 {
//				if self.TopEquipGem[i-1].Count < count {
//					change = true
//				}
//			}
//			break
//		}
//	}
//
//	if insert == true {
//		var data JS_EquipGemTop
//		data.Uid = player.Sql_UserBase.Uid
//		data.Uname = player.Sql_UserBase.UName
//		data.Iconid = player.Sql_UserBase.IconId
//		data.Camp = player.Sql_UserBase.Camp
//		data.Vip = player.Sql_UserBase.Vip
//		data.Count = count
//		self.TopEquipGem = append(self.TopEquipGem, &data)
//	}
//
//	if change == true || insert == true {
//		sort.Sort(lstTopEquipGem(self.TopEquipGem))
//		for i := 0; i < len(self.TopEquipGem); i++ {
//			if self.TopEquipGem[i] != nil {
//				self.TopEquipGem[i].LastRank = i + 1
//			}
//		}
//		self.Topver++
//	}
//
//	if len(self.TopEquipGem) > 200 {
//		self.TopEquipGem = self.TopEquipGem[:200]
//	}
//
//	self.Locker.Unlock()
//
//	if change == true || insert == true {
//		self.UpdateCampRank(count, player)
//	}
//
//}
func (self *TopEquipGemMgr) ResetCampRank() {
	self.CampLocker.Lock()
	defer self.CampLocker.Unlock()

	top := GetTopEquipGemMgr().GetTopEquipGem_L()

	for camp := 0; camp < 3; camp++ {
		self.CampEquipGem[camp] = make([]*JS_EquipGemTop, 0)
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
		self.CampEquipGem[camp] = append(self.CampEquipGem[camp], &JS_EquipGemTop{
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
//func (self *TopEquipGemMgr) UpdateCampRank(count int, player *Player) {
//	camp := player.GetCamp() - 1
//	if camp >= 0 && camp < CAMP_WU {
//		self.CampLocker.Lock()
//		defer self.CampLocker.Unlock()
//
//		insert := true  // 是否插入新数据
//		change := false // 是否重新排序
//		for i := 0; i < len(self.CampEquipGem[camp]); i++ {
//			if self.CampEquipGem[camp][i].Uid == player.Sql_UserBase.Uid {
//				self.CampEquipGem[camp][i].Count = count
//				self.CampEquipGem[camp][i].Level = player.Sql_UserBase.Level
//				insert = false
//				if i > 0 {
//					if self.CampEquipGem[camp][i-1].Count < count {
//						change = true
//					}
//				}
//				break
//			}
//		}
//
//		if insert == true {
//			var data JS_EquipGemTop
//			data.Uid = player.Sql_UserBase.Uid
//			data.Uname = player.Sql_UserBase.UName
//			data.Iconid = player.Sql_UserBase.IconId
//			data.Camp = player.Sql_UserBase.Camp
//			data.Vip = player.Sql_UserBase.Vip
//			data.Count = count
//			self.CampEquipGem[camp] = append(self.CampEquipGem[camp], &data)
//		}
//
//		if change == true || insert == true {
//			sort.Sort(lstTopEquipGem(self.CampEquipGem[camp]))
//			for i := 0; i < len(self.CampEquipGem[camp]); i++ {
//				if self.CampEquipGem[camp][i] != nil {
//					self.CampEquipGem[camp][i].LastRank = i + 1
//				}
//			}
//		}
//	}
//}

func (self *TopEquipGemMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	data, _ := self.TopEquipGemCur[id]
	return data
}

////////////////////////////////////////////////////////////////////////////////
func (s lstTopEquipGem) Len() int      { return len(s) }
func (s lstTopEquipGem) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstTopEquipGem) Less(i, j int) bool {
	if s[i].Count > s[j].Count {
		return true
	}

	if s[i].Count < s[j].Count {
		return false
	}

	if s[i].StartTime < s[j].StartTime {
		return true
	}

	if s[i].StartTime > s[j].StartTime {
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
