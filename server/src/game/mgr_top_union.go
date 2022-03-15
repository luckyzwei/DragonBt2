package game

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	//"time"
)

//军团排行榜
type JS_TopUnion struct {
	Id         int    `json:"id"`
	Icon       int    `json:"icon"`
	Unionname  string `json:"unionname"`
	Masteruid  string `json:"masteruid"`
	Mastername string `json:"masterName"`
	Level      int    `json:"level"`
	Fight      int64  `json:"fight"`
	Member     int    `json:"member"`
	Camp       int    `json:"camp"`
	LastRank   int    `json:"-"` //! 原有排名
}

type TopUnionMgr struct {
	TopUnion       []*JS_TopUnion //! 军团排行榜(实时)
	TopUnionActual []*JS_TopUnion //! 实际前n名战斗力排行
	TopUnionCur    map[int64]int
	TopUnionOld    map[int64]int
	TopverUnion    int
	Rank           int64
	Locker         *sync.RWMutex
}

var topUnionMgr *TopUnionMgr = nil

func GetTopUnionMgr() *TopUnionMgr {
	if topUnionMgr == nil {
		topUnionMgr = new(TopUnionMgr)
		topUnionMgr.TopUnion = make([]*JS_TopUnion, 0)
		topUnionMgr.TopUnionActual = make([]*JS_TopUnion, 0)
		topUnionMgr.TopUnionCur = make(map[int64]int, 0)
		topUnionMgr.TopUnionOld = make(map[int64]int, 0)
		topUnionMgr.Locker = new(sync.RWMutex)
	}
	return topUnionMgr
}

func (self *TopUnionMgr) InitTopActual() {
	// 军团数据初始化
	unionnum := MaxRankNum
	if len(self.TopUnion) < MaxRankNum {
		unionnum = len(self.TopUnion)
	}

	self.TopUnionActual = make([]*JS_TopUnion, 0)
	for i := 0; i < unionnum; i++ {
		var top JS_TopUnion
		top.Id = self.TopUnion[i].Id
		top.Fight = self.TopUnion[i].Fight
		top.Icon = self.TopUnion[i].Icon
		top.Level = self.TopUnion[i].Level
		top.Mastername = self.TopUnion[i].Mastername
		top.Masteruid = self.TopUnion[i].Masteruid
		top.Unionname = self.TopUnion[i].Unionname
		top.Camp = self.TopUnion[i].Camp
		top.Member = self.TopUnion[i].Member
		top.LastRank = i + 1

		self.TopUnionActual = append(self.TopUnionActual, &top)
	}
}

func (self *TopUnionMgr) GetData() {
	//军团
	self.Locker.Lock()
	if self.TopverUnion > 0 {
		HF_DeepCopy(&self.TopUnionOld, &self.TopUnionCur)
		self.TopUnion = make([]*JS_TopUnion, 0)
		self.TopUnionCur = make(map[int64]int)
	}

	var top3 JS_TopUnion
	sql3 := fmt.Sprintf("select `id`, `icon`, `unionname`, `masteruid`, `mastername`, `level`, `fight`, " +
		"`fight`, `camp` from `san_unioninfo` order by `fight` desc  limit 0, 200")
	res3 := GetServer().DBUser.GetAllDataEx(sql3, &top3)
	for i := 0; i < len(res3); i++ {
		data := res3[i].(*JS_TopUnion)
		self.TopUnion = append(self.TopUnion, data)
	}

	// 军团数据初始化
	unionnum := MaxRankNum
	if len(self.TopUnion) < MaxRankNum {
		unionnum = len(self.TopUnion)
	}

	self.TopUnionActual = make([]*JS_TopUnion, 0)
	for i := 0; i < unionnum; i++ {
		var top JS_TopUnion
		top.Id = self.TopUnion[i].Id
		top.Fight = self.TopUnion[i].Fight
		top.Icon = self.TopUnion[i].Icon
		top.Level = self.TopUnion[i].Level
		top.Mastername = self.TopUnion[i].Mastername
		top.Masteruid = self.TopUnion[i].Masteruid
		top.Unionname = self.TopUnion[i].Unionname
		top.Camp = self.TopUnion[i].Camp
		top.Member = self.TopUnion[i].Member
		lastRank, ok := self.TopUnionCur[int64(top.Id)]
		if ok {
			top.LastRank = lastRank
		}
		self.TopUnionActual = append(self.TopUnionActual, &top)
	}
	// 先排序
	sort.Sort(lstUnionTop(self.TopUnionActual))
	// 再设置rank
	for i := 0; i < len(self.TopUnionActual); i++ {
		if self.TopUnionActual[i] != nil {
			self.TopUnionActual[i].LastRank = i + 1
		}
	}

	for i := 0; i < len(self.TopUnion); i++ {
		self.TopUnionCur[int64(self.TopUnion[i].Id)] = i + 1
	}

	self.TopverUnion++
	self.Locker.Unlock()

	//! 更新人数
	for i := 0; i < len(self.TopUnion); i++ {
		union := GetUnionMgr().GetUnion(self.TopUnion[i].Id)
		if union != nil {
			self.TopUnion[i].Member = len(union.member)
		}
	}
}

func (self *TopUnionMgr) ChangeUnionName(unionid int, unionname string) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if self.TopUnionActual == nil || len(self.TopUnionActual) <= 0 {
		return
	}

	for index := range self.TopUnionActual {
		if self.TopUnionActual[index].Id == unionid {
			self.TopUnionActual[index].Unionname = unionname
			break
		}
	}
}

//获取军团排行榜
func (self *TopUnionMgr) GetTopUnion() ([]*JS_TopUnion, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.TopUnionActual) > MaxRankNum {
		return self.TopUnionActual[:MaxRankNum], self.TopverUnion
	}
	return self.TopUnionActual, self.TopverUnion
}

func (self *TopUnionMgr) GetTopUnionShow() ([]*JS_TopUnion, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.TopUnionActual) > MaxRankShowNum {
		return self.TopUnionActual[:MaxRankShowNum], self.TopverUnion
	}
	return self.TopUnionActual, self.TopverUnion
}

//! 实时更新战斗力,这里其实是注释掉的, 改成每小时更新
func (self *TopUnionMgr) SyncUnionFight(fight int64, union *San_Union) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	insert := true  // 是否插入新数据
	change := false // 是否重新排序
	for i := 0; i < len(self.TopUnionActual); i++ {
		if self.TopUnionActual[i].Id == union.Id {
			if self.TopUnionActual[i].Fight == fight {
				return
			}
			self.TopUnionActual[i].Fight = fight
			self.TopUnionActual[i].Member = len(union.member)
			self.TopUnionActual[i].Level = union.Level
			insert = false
			if i > 0 {
				if self.TopUnionActual[i-1].Fight < fight {
					change = true
				}
			}
			break
		}
	}

	if insert == true {
		top := new(JS_TopUnion)
		top.Id = union.Id
		top.Fight = fight
		top.Icon = union.Icon
		top.Level = union.Level
		top.Mastername = union.Mastername
		top.Masteruid = strconv.FormatInt(union.Masteruid, 10)
		top.Unionname = union.Unionname
		//top.Camp = union.Camp
		top.Member = len(union.member)

		self.TopUnionActual = append(self.TopUnionActual, top)
	}

	LogDebug("军团更新排行榜:", insert, change, union.Unionname, fight)
	if change == true || insert == true {
		if TimeServer().Unix() > self.Rank+3 {
			self.Rank = TimeServer().Unix()
			sort.Sort(lstUnionTop(self.TopUnionActual))
			for i := 0; i < len(self.TopUnionActual); i++ {
				if self.TopUnionActual[i] != nil {
					self.TopUnionActual[i].LastRank = i + 1
				}
			}
			self.TopverUnion++

			if len(self.TopUnionActual) > MaxRankNum {
				self.TopUnionActual = self.TopUnionActual[:MaxRankNum]
			}
		}
	}

}

//! 获得分阵营军团数据
func (self *TopUnionMgr) getCampUnion(camp int) []*JS_TopUnion {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var unionTop []*JS_TopUnion
	for i := 0; i < len(self.TopUnionActual); i++ {
		if self.TopUnionActual[i].Camp == camp {
			unionTop = append(unionTop, self.TopUnionActual[i])
		}
	}

	return unionTop
}

func (self *TopUnionMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	data, _ := self.TopUnionCur[id]
	self.Locker.RUnlock()
	return data
}

func (self *TopUnionMgr) GetTopOldNum(topType int, id int64) int {
	self.Locker.RLock()
	data, _ := self.TopUnionOld[id]
	self.Locker.RUnlock()
	return data
}

////////////////////////////////////////////////////////////////////////////////
type lstUnionTop []*JS_TopUnion

func (s lstUnionTop) Len() int      { return len(s) }
func (s lstUnionTop) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstUnionTop) Less(i, j int) bool {
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

	if s[i].Id < s[j].Id {
		return true
	}

	if s[i].Id > s[j].Id {
		return false
	}

	return true
}

// 根据军团ID返回排行榜名次
func (self *TopUnionMgr) GetUnionRank(unionId int) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	for i, v := range self.TopUnion {
		if v.Id == unionId {
			return i + 1
		}
	}
	return 0
}
