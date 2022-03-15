package game

import (
	"fmt"
	"sort"
	"sync"
	//"time"
)

//
////星级排行榜,武将，关卡共用
//type JS_TopStar struct {
//	Uid       int64  `json:"uid"`
//	Uname     string `json:"uname"`
//	Iconid    int    `json:"iconid"`
//	Portrait  int    `json:"portrait"` // 边框  20190412 by zy
//	Level     int    `json:"level"`
//	Camp      int    `json:"camp"`
//	Fight     int    `json:"fight"`
//	Vip       int    `json:"vip"`
//	Star      int    `json:"star"`
//	LastRank  int    `json:"-"` //! 原有排名
//	UnionName string `json:"union_name"`
//	StartTime int64  `json:"starttime"` //! 时间戳
//}
//
//// 星级
//type lstJsTopStar []*JS_TopStar
//
//func (s lstJsTopStar) Len() int      { return len(s) }
//func (s lstJsTopStar) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
//func (s lstJsTopStar) Less(i, j int) bool {
//	if s[i].Star > s[j].Star { // 由大到小
//		return true
//	}
//
//	if s[i].Star < s[j].Star {
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
//	if s[i].LastRank < s[j].LastRank { // 由大到小
//		return true
//	}
//
//	if s[i].LastRank > s[j].LastRank {
//		return false
//	}
//
//	if s[i].Uid < s[j].Uid { // 由小到大
//		return true
//	}
//
//	if s[i].Uid > s[j].Uid {
//		return false
//	}
//
//	return false
//}

type TopPassMgr struct {
	TopPass    []*Js_ActTop //! 关卡总星级排行榜(从数据库非实时拉取)
	TopPassCur map[int64]int
	TopPassOld map[int64]int
	TopverPass int
	PassLocker *sync.RWMutex
}

var topPassMgr *TopPassMgr = nil

func GetTopPassMgr() *TopPassMgr {
	if topPassMgr == nil {
		topPassMgr = new(TopPassMgr)
		topPassMgr.TopPass = make([]*Js_ActTop, 0)
		topPassMgr.TopPassCur = make(map[int64]int, 0)
		topPassMgr.TopPassOld = make(map[int64]int, 0)
		topPassMgr.PassLocker = new(sync.RWMutex)
	}
	return topPassMgr
}

func (self *TopPassMgr) GetData(unionName map[int64]string) {
	// 副本星级
	self.PassLocker.Lock()
	if self.TopverPass > 0 {
		HF_DeepCopy(&self.TopPassOld, &self.TopPassCur)
		self.TopPass = make([]*Js_ActTop, 0)
		self.TopPassCur = make(map[int64]int)
	}

	var top4 Js_ActTop
	starText := GetCsvMgr().GetText("STR_TOP_NULL")
	sql4 := fmt.Sprintf("select a.uid,a.uname,a.iconid,a.portrait,a.level,a.camp,a.fight,a.vip,b.onhookstage, "+
		"'%s' as union_name ,b.onhookstagetime from san_userbase as a JOIN san_onhook as b where a.uid = b.uid and b.onhookstage > 0 "+
		"ORDER BY b.onhookstage desc , b.onhookstagetime ASC limit 200", starText)
	res4 := GetServer().DBUser.GetAllDataEx(sql4, &top4)
	for i := 0; i < len(res4); i++ {
		data := res4[i].(*Js_ActTop)
		lastRank, ok := self.TopPassCur[data.Uid]
		if ok {
			data.LastRank = lastRank
		}

		// 更新军团名字
		v, ok := unionName[data.Uid]
		if ok {
			data.UnionName = v
		}

		self.TopPass = append(self.TopPass, data)
	}
	sort.Sort(lstJsActTop(self.TopPass))
	for i := 0; i < len(self.TopPass); i++ {
		self.TopPass[i].LastRank = i + 1
		self.TopPassCur[self.TopPass[i].Uid] = i + 1
	}
	self.TopverPass++
	self.PassLocker.Unlock()

}

//获取副本排行榜
func (self *TopPassMgr) GetTopPass() ([]*Js_ActTop, int) {
	self.PassLocker.RLock()
	defer self.PassLocker.RUnlock()

	if len(self.TopPass) > MaxRankNum {
		return self.TopPass[:MaxRankNum], self.TopverPass
	}
	return self.TopPass, self.TopverPass
}

func (self *TopPassMgr) GetTopPassShow() ([]*Js_ActTop, int) {
	self.PassLocker.RLock()
	defer self.PassLocker.RUnlock()

	if len(self.TopPass) > MaxRankShowNum {
		return self.TopPass[:MaxRankShowNum], self.TopverPass
	}
	return self.TopPass, self.TopverPass
}

func (self *TopPassMgr) GetTopCurNum(topType int, id int64) int {
	self.PassLocker.RLock()
	data, _ := self.TopPassCur[id]
	self.PassLocker.RUnlock()
	return data
}

func (self *TopPassMgr) GetTopOldNum(topType int, id int64) int {
	self.PassLocker.RLock()
	data, _ := self.TopPassOld[id]
	self.PassLocker.RUnlock()
	return data
}

func (self *TopPassMgr) SyncPlayerName(player *Player) {
	self.PassLocker.Lock()
	for i := 0; i < len(self.TopPass); i++ {
		if self.TopPass[i].Uid == player.Sql_UserBase.Uid {
			self.TopPass[i].Uname = player.Sql_UserBase.UName
			break
		}
	}
	self.PassLocker.Unlock()
}

// 更新排行数据
func (self *TopPassMgr) UpdateRank(count int64, player *Player) {
	self.PassLocker.Lock()

	insert := true  // 是否插入新数据
	change := false // 是否重新排序
	for i := 0; i < len(self.TopPass); i++ {
		if self.TopPass[i].Uid == player.Sql_UserBase.Uid {
			self.TopPass[i].Num = count
			self.TopPass[i].Level = player.Sql_UserBase.Level
			insert = false
			if i > 0 {
				if self.TopPass[i-1].Num < count {
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
		data.Camp = player.Sql_UserBase.Camp
		data.Vip = player.Sql_UserBase.Vip
		data.Fight = player.Sql_UserBase.Fight
		data.Num = count
		data.StartTime = TimeServer().Unix()
		data.UnionName = player.GetUnionName()
		self.TopPass = append(self.TopPass, &data)
	}

	if change == true || insert == true {
		sort.Sort(lstJsActTop(self.TopPass))
		for i := 0; i < len(self.TopPass); i++ {
			if self.TopPass[i] != nil {
				self.TopPass[i].LastRank = i + 1
			}
		}
		self.TopverPass++
	}

	if len(self.TopPass) > 200 {
		self.TopPass = self.TopPass[:200]
	}

	self.PassLocker.Unlock()
}
