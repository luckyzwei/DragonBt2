package game

import (
	"fmt"
	"sort"
	"sync"
)

type TopGemCostMgr struct {
	Top    []*Js_ActTop
	TopCur map[int64]int
	TopOld map[int64]int
	Topver int
	Locker *sync.RWMutex
}

var topgemcost *TopGemCostMgr = nil

func GetTopGemCostMgr() *TopGemCostMgr {
	if topgemcost == nil {
		topgemcost = new(TopGemCostMgr)
		topgemcost.Top = make([]*Js_ActTop, 0)
		topgemcost.TopCur = make(map[int64]int, 0)
		topgemcost.TopOld = make(map[int64]int, 0)
		topgemcost.Locker = new(sync.RWMutex)
	}
	return topgemcost
}

// 钻石消耗
func (self *TopGemCostMgr) GetData() {
	self.Locker.Lock()
	if self.Topver > 0 {
		HF_DeepCopy(&self.TopOld, &self.TopCur)
		self.Top = make([]*Js_ActTop, 0)
		self.TopCur = make(map[int64]int)
	}

	var topGemCost Js_ActTop
	//step := GetActivityMgr().getGemTaskN4()
	step := 0
	topGemCostSql := fmt.Sprintf("select a.uid,a.uname,a.iconid,a.level,a.camp, a.fight, a.vip, b.cost1 as num "+
		"from san_userbase as a JOIN  san_cost as b where a.uid = b.uid and b.cost1 > 0 and b.step1 = %d "+
		"ORDER BY b.cost1 desc limit 200", step)
	gemCostRes := GetServer().DBUser.GetAllDataEx(topGemCostSql, &topGemCost)
	for i := 0; i < len(gemCostRes); i++ {
		data := gemCostRes[i].(*Js_ActTop)
		lastRank, ok := self.TopCur[data.Uid]
		if ok {
			data.LastRank = lastRank
		}
		self.Top = append(self.Top, data)
	}
	sort.Sort(lstJsActTop(self.Top))
	for i := 0; i < len(self.Top); i++ {
		self.Top[i].LastRank = i + 1
		self.TopCur[self.Top[i].Uid] = i + 1
	}
	self.Topver++

	self.Locker.Unlock()
}

func (self *TopGemCostMgr) GetTopShow() ([]*Js_ActTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.Top) > MaxRankShowNum {
		return self.Top[:MaxRankShowNum], self.Topver
	}
	return self.Top, self.Topver
}

//获取钻石消耗
func (self *TopGemCostMgr) GetTop() ([]*Js_ActTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	if len(self.Top) > MaxRankNum {
		return self.Top[:MaxRankNum], self.Topver
	}
	return self.Top, self.Topver
}

func (self *TopGemCostMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	data, _ := self.TopCur[id]
	self.Locker.RUnlock()
	return data
}

func (self *TopGemCostMgr) GetTopOldNum(topType int, id int64) int {
	self.Locker.RLock()
	data, _ := self.TopOld[id]
	self.Locker.RUnlock()
	return data
}
