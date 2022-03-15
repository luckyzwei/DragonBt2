package game

import (
	"fmt"
	"sort"
	"sync"
	//"time"
)

// 战力排行榜
type JS_Top struct {
	Uid       int64  `json:"uid"`
	Uname     string `json:"uname"`
	Iconid    int    `json:"iconid"`
	Portrait  int    `json:"portrait"` // 边框  20190412 by zy
	Level     int    `json:"level"`
	Camp      int    `json:"camp"`
	Fight     int64  `json:"fight"`
	Vip       int    `json:"vip"`
	UnionName string `json:"union_name"`
	LastRank  int    `json:"-"` //! 原有排名
	Exp       int    `json:"exp"`
}

type TopFightMgr struct {
	Top       []*JS_Top //! 战力排行榜(实时)
	TopActual []*JS_Top //! 实际前n名战斗力排行
	TopCur    map[int64]int
	TopOld    map[int64]int
	Topver    int
	Rank      int64
	Locker    *sync.RWMutex
}

var topFightMgr *TopFightMgr = nil

//! public
func GetTopFightMgr() *TopFightMgr {
	if topFightMgr == nil {
		topFightMgr = new(TopFightMgr)
		topFightMgr.Top = make([]*JS_Top, 0)
		topFightMgr.TopActual = make([]*JS_Top, 0)
		topFightMgr.TopCur = make(map[int64]int, 0)
		topFightMgr.TopOld = make(map[int64]int, 0)
		topFightMgr.Locker = new(sync.RWMutex)
	}

	return topFightMgr
}

func (self *TopFightMgr) GetData(unionName map[int64]string) {
	//! 战力获取
	self.Locker.Lock()
	if self.Topver > 0 {
		HF_DeepCopy(&self.TopOld, &self.TopCur)
		self.Top = make([]*JS_Top, 0)
		self.TopCur = make(map[int64]int)
	}

	//战斗力榜
	var top JS_Top
	fighText := GetCsvMgr().GetText("STR_TOP_NULL")
	//增加 portrait 20190412 by zy
	sql := fmt.Sprintf("select `uid`, `uname`, `iconid`, `portrait`, `level`, `camp`, `fight`, `vip`, "+
		"'%s' as union_name, `exp` from `san_userbase` where `fight` > 0 order by `fight` desc  limit 0, 2000", fighText)
	res := GetServer().DBUser.GetAllDataEx(sql, &top)
	for i := 0; i < len(res); i++ {
		data, ok := res[i].(*JS_Top)
		if !ok {
			continue
		}
		// 更新军团名字
		v, ok := unionName[data.Uid]
		if ok {
			data.UnionName = v
		}
		//fmt.Printf("%#v\n", data)
		self.Top = append(self.Top, data)
	}

	for i := 0; i < len(self.Top); i++ {
		self.TopCur[self.Top[i].Uid] = i + 1
	}

	self.Topver++
	self.Locker.Unlock()
}

//获取战力排行榜
func (self *TopFightMgr) GetTop() ([]*JS_Top, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.TopActual) > MaxRankNum {
		return self.TopActual[:MaxRankNum], self.Topver
	}
	return self.TopActual, self.Topver
}

//获取战力排行榜
func (self *TopFightMgr) GetTopShow() ([]*JS_Top, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.TopActual) > MaxRankShowNum {
		return self.TopActual[:MaxRankShowNum], self.Topver
	}
	return self.TopActual, self.Topver
}

func (self *TopFightMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	data, _ := self.TopCur[id]
	self.Locker.RUnlock()
	return data
}

func (self *TopFightMgr) GetTopOldNum(topType int, id int64) int {
	self.Locker.RLock()
	data, _ := self.TopOld[id]
	self.Locker.RUnlock()
	return data
}

type lstTop []*JS_Top

func (s lstTop) Len() int      { return len(s) }
func (s lstTop) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstTop) Less(i, j int) bool {
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

//! 实时更新战斗力
func (self *TopFightMgr) SyncFight(fight int64, player *Player) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	insert := true  //! 是否插入新数据
	change := false //! 是否重新排序
	if len(self.TopActual) >= MaxRankNum {
		tailplayer := self.TopActual[len(self.TopActual)-1]
		if tailplayer.Fight > fight {
			return change
		}
	}

	for i := 0; i < len(self.TopActual); i++ {
		if self.TopActual[i].Uid == player.Sql_UserBase.Uid {
			self.TopActual[i].Level = player.Sql_UserBase.Level
			self.TopActual[i].Vip = player.Sql_UserBase.Vip
			self.TopActual[i].UnionName = player.GetUnionName()
			self.TopActual[i].Camp = player.Sql_UserBase.Camp
			self.TopActual[i].Iconid = player.Sql_UserBase.IconId
			self.TopActual[i].Portrait = player.Sql_UserBase.Portrait

			insert = false
			if self.TopActual[i].Fight == fight {
				return change
			}
			self.TopActual[i].Uname = player.Sql_UserBase.UName
			self.TopActual[i].Fight = fight
			if i > 0 {
				if self.TopActual[i-1].Fight < fight {
					change = true
				}

				if i < len(self.TopActual)-1 {
					if self.TopActual[i+1].Fight > fight {
						change = true
					}
				}
			}

			break
		}
	}

	if insert == true {
		// 过滤阵营为0的玩家
		if player.Sql_UserBase.Camp < 1 || player.Sql_UserBase.Camp > 3 {
			return false
		}
		top := new(JS_Top)
		top.Uid = player.Sql_UserBase.Uid
		top.Camp = player.Sql_UserBase.Camp

		top.Level = player.Sql_UserBase.Level
		top.Iconid = player.Sql_UserBase.IconId
		top.Portrait = player.Sql_UserBase.Portrait // 边框  20190412 by zy
		top.Fight = fight
		top.Uname = player.Sql_UserBase.UName
		top.Vip = player.Sql_UserBase.Vip
		top.UnionName = player.GetUnionName()

		self.TopActual = append(self.TopActual, top)
	}

	if change == true || insert == true {
		if TimeServer().Unix() > self.Rank {
			self.Rank = TimeServer().Unix()
			sort.Sort(lstTop(self.TopActual))
			// 重新设置上次排名
			for i := 0; i < len(self.TopActual); i++ {
				self.TopActual[i].LastRank = i + 1
			}
			self.Topver++

			if len(self.TopActual) > MaxRankNum {
				self.TopActual = self.TopActual[:MaxRankNum]
			}
		}

	}

	return change || insert
	//LogDebug("更新排行榜:", insert, change, player.Sql_UserBase.UName, player.Sql_UserBase.Fight, len(self.TopActual))
}

func (self *TopFightMgr) InitTopActual() {
	// 战斗力数据初始化
	topnum := MaxRankNum
	if len(self.Top) < MaxRankNum {
		topnum = len(self.Top)
	}
	self.TopActual = make([]*JS_Top, 0)
	for i := 0; i < topnum; i++ {
		var top JS_Top
		top.Uid = self.Top[i].Uid
		top.Camp = self.Top[i].Camp
		top.Level = self.Top[i].Level
		top.Iconid = self.Top[i].Iconid
		top.Portrait = self.Top[i].Portrait // 边框  20190412 by zy
		top.Fight = self.Top[i].Fight
		top.Uname = self.Top[i].Uname
		top.Vip = self.Top[i].Vip
		top.UnionName = self.Top[i].UnionName
		top.LastRank = i + 1

		self.TopActual = append(self.TopActual, &top)
	}
}

func (self *TopFightMgr) SyncPlayerName(player *Player) {
	//实时排行榜更新名字
	self.Locker.Lock()
	for i := 0; i < len(self.TopActual); i++ {
		if self.TopActual[i].Uid == player.Sql_UserBase.Uid {
			self.TopActual[i].Uname = player.Sql_UserBase.UName
			break
		}
	}
	self.Locker.Unlock()
}
