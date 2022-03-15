package game

import (
	"fmt"
	"sort"
	"sync"
	//"time"
)

type TopLevelMgr struct {
	TopLevel       []*JS_Top //! 等级排行榜(实时)
	TopLevelActual []*JS_Top //! 实际前n名战斗力排行
	TopLevelCur    map[int64]int
	TopLevelOld    map[int64]int
	TopverLevel    int
	Rank           int64
	Locker         *sync.RWMutex
}

var toplevelmgr *TopLevelMgr = nil

func GetTopLevelMgr() *TopLevelMgr {
	if toplevelmgr == nil {
		toplevelmgr = new(TopLevelMgr)
		toplevelmgr.TopLevel = make([]*JS_Top, 0)
		toplevelmgr.TopLevelActual = make([]*JS_Top, 0)
		toplevelmgr.TopLevelCur = make(map[int64]int, 0)
		toplevelmgr.TopLevelOld = make(map[int64]int, 0)
		toplevelmgr.Locker = new(sync.RWMutex)
	}
	return toplevelmgr
}

//! 初始化动态排行榜[战力,军团,等级]
func (self *TopLevelMgr) InitTopActual() {
	self.TopLevelActual = make([]*JS_Top, 0)
	levelnum := MaxRankNum
	if len(self.TopLevel) < MaxRankNum {
		levelnum = len(self.TopLevel)
	}
	//self.TopLevelActual = make([]*JS_Top, 0)
	for i := 0; i < levelnum; i++ {
		var top JS_Top
		top.Uid = self.TopLevel[i].Uid
		top.Camp = self.TopLevel[i].Camp
		top.Level = self.TopLevel[i].Level
		top.Iconid = self.TopLevel[i].Iconid
		top.Portrait = self.TopLevel[i].Portrait
		top.Fight = self.TopLevel[i].Fight
		top.Uname = self.TopLevel[i].Uname
		top.Vip = self.TopLevel[i].Vip
		top.UnionName = self.TopLevel[i].UnionName
		top.LastRank = i + 1

		self.TopLevelActual = append(self.TopLevelActual, &top)
	}
	//LogDebug("实时等级排行榜:", len(self.TopLevelActual), len(self.TopLevel), len(self.TopActual))
}

func (self *TopLevelMgr) GetData(unionName map[int64]string) {
	//! 等级排行榜
	self.Locker.Lock()
	if self.TopverLevel > 0 {
		HF_DeepCopy(&self.TopLevelOld, &self.TopLevelCur)
		self.TopLevel = make([]*JS_Top, 0)
		self.TopLevelCur = make(map[int64]int)
	}

	var top2 JS_Top
	topText := GetCsvMgr().GetText("STR_TOP_NULL")
	sql2 := fmt.Sprintf("select `uid`, `uname`, `iconid`, `portrait`,`level`, `camp`, `fight`, `vip`, '%s' as union_name ,`exp` from `san_userbase` " +
		"where `fight` > 0 order by `level` desc, `exp` desc limit 0, 200", topText)
	res2 := GetServer().DBUser.GetAllDataEx(sql2, &top2)
	for i := 0; i < len(res2); i++ {
		data, ok := res2[i].(*JS_Top)
		if !ok {
			continue
		}
		// 更新军团名字
		v, ok := unionName[data.Uid]
		if ok {
			data.UnionName = v
		}
		self.TopLevel = append(self.TopLevel, data)
	}

	for i := 0; i < len(self.TopLevel); i++ {
		self.TopLevelCur[self.TopLevel[i].Uid] = i + 1
	}

	self.TopverLevel++
	self.Locker.Unlock()
}

// 获取等级排行榜
func (self *TopLevelMgr) GetTopLevel() ([]*JS_Top, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.TopLevelActual) > MaxRankNum {
		return self.TopLevelActual[:MaxRankNum], self.TopverLevel
	}
	return self.TopLevelActual, self.TopverLevel
}

//! 实时更新战斗力
func (self *TopLevelMgr) SyncLevel(level int, player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	if len(self.TopLevelActual) > 0 {
		tailplayer := self.TopLevelActual[len(self.TopLevelActual)-1]
		if tailplayer.Level > level {
			return
		}
	}

	insert := true  // 是否插入新数据
	change := false // 是否重新排序
	for i := 0; i < len(self.TopLevelActual); i++ {
		if self.TopLevelActual[i].Uid == player.Sql_UserBase.Uid {
			self.TopLevelActual[i].Fight = player.Sql_UserBase.Fight
			self.TopLevelActual[i].Vip = player.Sql_UserBase.Vip
			if self.TopLevelActual[i].Level == level {
				return
			}
			self.TopLevelActual[i].Level = level
			insert = false
			if i > 0 {
				if self.TopLevelActual[i-1].Level < level {
					change = true
				}
			}

			break
		}
	}

	if insert == true {
		top := new(JS_Top)
		top.Uid = player.Sql_UserBase.Uid
		top.Camp = player.Sql_UserBase.Camp
		top.Level = player.Sql_UserBase.Level
		top.Iconid = player.Sql_UserBase.IconId
		top.Portrait = player.Sql_UserBase.Portrait
		top.Fight = int64(level)
		top.Uname = player.Sql_UserBase.UName
		top.Vip = player.Sql_UserBase.Vip

		self.TopLevelActual = append(self.TopLevelActual, top)
	}

	LogDebug("更新等级排行榜:", insert, change, player.Sql_UserBase.UName, player.Sql_UserBase.Level)
	if change == true || insert == true {
		if TimeServer().Unix() > self.Rank+3 {
			self.Rank = TimeServer().Unix()
			sort.Sort(lstLevelTop(self.TopLevelActual))
			for i := 0; i < len(self.TopLevelActual); i++ {
				if self.TopLevelActual[i] != nil {
					self.TopLevelActual[i].LastRank = i + 1
				}
			}
			self.TopverLevel++

			if len(self.TopLevelActual) > MaxRankNum {
				self.TopLevelActual = self.TopLevelActual[:MaxRankNum]
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
type lstLevelTop []*JS_Top

func (s lstLevelTop) Len() int      { return len(s) }
func (s lstLevelTop) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstLevelTop) Less(i, j int) bool {
	if s[i].Level > s[j].Level {
		return true
	}

	if s[i].Level < s[j].Level {
		return false
	}

	if s[i].Exp > s[j].Exp {
		return true
	}

	if s[i].Exp < s[j].Exp {
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
