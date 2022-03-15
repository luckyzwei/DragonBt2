package game

import (
	"fmt"
	"sort"
	"sync"
)

type JS_PvpTop struct {
	Uid       int64  `json:"uid"`
	Uname     string `json:"uname"`
	Iconid    int    `json:"iconid"`
	Portrait  int    `json:"portrait"` // 边框  20190412 by zy
	Level     int    `json:"level"`
	Camp      int    `json:"camp"`
	RankID    int    `json:"rankid"`
	Vip       int    `json:"vip"`
	LastRank  int    `json:"-"`
	UnionName string `json:"union_name"`
}

type lstTopPvp []*JS_PvpTop
type TopPvpMgr struct {
	// 排行数据
	TopPvp [3][]*JS_PvpTop // 斗技场排行榜 分为三个阵营
	// 版本号
	Topver int
	Locker *sync.RWMutex
}

// 初始斗技场排行
func (self *TopPvpMgr) GetData(unionName map[int64]string) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for camp := 0; camp < 3; camp++ {
		topPvpMgr.TopPvp[camp] = make([]*JS_PvpTop, 0)
	}

	for camp := CAMP_SHU; camp < CAMP_QUN; camp++ {
		var top5 JS_PvpTop
		text := GetCsvMgr().GetText("STR_TOP_NULL")
		sql5 := fmt.Sprintf("SELECT a.uid, a.uname, a.iconid,a.portrait , a.level, a.camp, b.rankid, a.vip , '%s' as ", text)
		sql5 += fmt.Sprintf("union_name FROM san_userbase AS a INNER JOIN (SELECT san_armsarena%d.uid, san_armsarena%d.rankid, san_armsarena%d.point FROM san_armsarena%d WHERE san_armsarena%d.rankid < 20 ORDER BY rankid DESC LIMIT 200) AS b ", camp, camp, camp, camp, camp)
		sql5 += "ON a.uid = b.uid"
		res5 := GetServer().DBUser.GetAllDataEx(sql5, &top5)
		for i := 0; i < len(res5); i++ {
			data := res5[i].(*JS_PvpTop)
			data.LastRank = data.RankID
			// 更新军团名字
			v, ok := unionName[data.Uid]
			if ok {
				data.UnionName = v
			}
			topPvpMgr.TopPvp[camp-1] = append(topPvpMgr.TopPvp[camp-1], data)
		}
	}

	for camp := 0; camp < 3; camp++ {
		sort.Sort(lstTopPvp(self.TopPvp[camp]))
		for i := 0; i < len(self.TopPvp[camp]); i++ {
			if self.TopPvp[camp][i] != nil {

				self.TopPvp[camp][i].LastRank = self.TopPvp[camp][i].RankID
			}
		}
	}

	self.Topver++
}

// 获取国家斗技场排行榜
func (self *TopPvpMgr) getCampRank(camp int) []*JS_PvpTop {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if camp < 1 || camp > 3 {
		return []*JS_PvpTop{}
	}

	return self.TopPvp[camp-1]
}

//// 获取斗技场排行榜
//func (self *TopPvpMgr) GetTopPvp_L() []*JS_PvpTop {
//	self.Locker.RLock()
//	defer self.Locker.RUnlock()
//
//	return self.TopPvp
//}

//func (self *TopPvpMgr) GetTopPvpInfo() ([]*JS_PvpTop, int) {
//	self.Locker.RLock()
//	defer self.Locker.RUnlock()
//
//	if len(self.TopPvp) > MaxRankNum {
//		return self.TopPvp[:MaxRankNum], self.Topver
//	}
//	return self.TopPvp, self.Topver
//}

var topPvpMgr *TopPvpMgr = nil

func GetTopPvpMgr() *TopPvpMgr {
	if topPvpMgr == nil {
		topPvpMgr = new(TopPvpMgr)
		for camp := 0; camp < 3; camp++ {
			topPvpMgr.TopPvp[camp] = make([]*JS_PvpTop, 0)
		}
		topPvpMgr.Locker = new(sync.RWMutex)
	}

	return topPvpMgr
}

//排行同步改名
func (self *TopPvpMgr) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	camp := player.GetCamp() - 1
	for i := 0; i < len(self.TopPvp[camp]); i++ {
		if self.TopPvp[camp][i].Uid == player.Sql_UserBase.Uid {
			self.TopPvp[camp][i].Uname = player.Sql_UserBase.UName
			break
		}
	}
}

//
//// 更新排行数据
//func (self *TopPvpMgr) updateRank(rankid int, player *Player) {
//	self.Locker.Lock()
//	defer self.Locker.Unlock()
//
//	camp := player.GetCamp() - 1
//
//	insert := true  // 是否插入新数据
//	change := false // 是否重新排序
//	for i := 0; i < len(self.TopPvp[camp]); i++ {
//		if self.TopPvp[camp][i].Uid == player.Sql_UserBase.Uid {
//
//			self.TopPvp[camp][i].Level = player.Sql_UserBase.Level
//
//			if rankid != self.TopPvp[camp][i].RankID {
//				self.TopPvp[camp][i].RankID = rankid
//			}
//
//			if i > 0 {
//				if self.TopPvp[camp][i-1].RankID > rankid {
//					change = true
//				}
//			}
//
//			insert = false
//		}
//	}
//
//	if insert == true {
//		var data JS_PvpTop
//		data.Uid = player.Sql_UserBase.Uid
//		data.Uname = player.Sql_UserBase.UName
//		data.Iconid = player.Sql_UserBase.IconId
//		data.Camp = player.Sql_UserBase.Camp
//		data.Vip = player.Sql_UserBase.Vip
//		data.RankID = rankid
//		self.TopPvp[camp] = append(self.TopPvp[camp], &data)
//	}
//
//	if change == true || insert == true {
//
//		sort.Sort(lstTopPvp(self.TopPvp[camp]))
//		for i := 0; i < len(self.TopPvp[camp]); i++ {
//			if self.TopPvp[camp][i] != nil {
//				self.TopPvp[camp][i].LastRank = self.TopPvp[camp][i].RankID
//			}
//		}
//		self.Topver++
//	}
//
//	if len(self.TopPvp[camp]) > 200 {
//		self.TopPvp[camp] = self.TopPvp[camp][:200]
//	}
//}

//// 更新每个国家的排行
//func (self *TopPvpMgr) updateCampRank() {
//	top := GetTopPvpMgr().GetTopPvp_L()
//	self.Locker.Lock()
//	for camp := 0; camp < 3; camp++ {
//		self.CampPvp[camp] = make([]*JS_PvpTop, 0)
//	}
//
//	for index := range top {
//		topInfo := top[index]
//		if topInfo == nil {
//			continue
//		}
//
//		if topInfo.Camp < 1 || topInfo.Camp > 3 {
//			continue
//		}
//
//		camp := topInfo.Camp - 1
//		self.CampPvp[camp] = append(self.CampPvp[camp], &JS_PvpTop{
//			Uid:    topInfo.Uid,
//			Uname:  topInfo.Uname,
//			Iconid: topInfo.Iconid,
//			Level:  topInfo.Level,
//			Camp:   topInfo.Camp,
//			Vip:    topInfo.Vip,
//			Fight:  topInfo.Fight,
//		})
//	}
//	self.CampLocker.Unlock()
//}

func (s lstTopPvp) Len() int      { return len(s) }
func (s lstTopPvp) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstTopPvp) Less(i, j int) bool {
	if s[i].RankID < s[j].RankID {
		return true
	}

	if s[i].RankID > s[j].RankID {
		return false
	}

	if s[i].LastRank > s[j].LastRank {
		return true
	}

	if s[i].LastRank < s[j].LastRank {
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
