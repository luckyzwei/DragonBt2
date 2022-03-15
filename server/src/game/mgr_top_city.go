package game

import (
	"encoding/json"
	"sync"
)

//! 城池数量
type CityNum struct {
	Camp int //! 阵营
	Num  int //! 数量
	Rank int //! 原来的排名
}

type TopSave struct {
	CityNum []*CityNum `json:"citynum"`
}

type LstCityNum []*CityNum

func (s LstCityNum) Len() int      { return len(s) }
func (s LstCityNum) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s LstCityNum) Less(i, j int) bool {
	if s[i].Num > s[j].Num {
		return true
	}

	if s[i].Num < s[j].Num {
		return false
	}

	if s[i].Rank < s[j].Rank {
		return true
	}

	if s[i].Rank > s[j].Rank {
		return false
	}

	if s[i].Camp < s[j].Camp {
		return true
	}

	if s[i].Camp > s[j].Camp {
		return false
	}

	return false
}

type TopCityMgr struct {
	ListCityLock *sync.RWMutex
	ListCityNum  []*CityNum
}

var topcitymgr *TopCityMgr = nil

func GetTopCityMgr() *TopCityMgr {
	if topcitymgr == nil {
		topcitymgr = new(TopCityMgr)
		topcitymgr.ListCityNum = make([]*CityNum, 0)
		topcitymgr.ListCityLock = new(sync.RWMutex)
	}
	return topcitymgr
}

func (self *TopCityMgr) GetCityNum() []*CityNum {
	self.ListCityLock.RLock()
	defer self.ListCityLock.RUnlock()

	return self.ListCityNum
}

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// 设置cityNum
func (self *TopCityMgr) setTopRedis() {
	Try(func() {
		if len(self.ListCityNum) <= 0 {
			return
		}

		var topSave TopSave
		topSave.CityNum = self.ListCityNum
		topInfo := HF_JtoA(topSave)
		c := GetServer().GetRedisConn()
		defer c.Close()
		_, err := c.Do("SET", "rankData:cityNum", topInfo)
		if err != nil {
			LogError(err.Error())
			return
		}
	}, func(e interface{}) {
		LogError(e)
	})
}

// 获取国家城池排行榜, 只可能是全服排行, 需要从缓存和实际拿
func (self *TopCityMgr) getCityNumRank() []*Js_ActTop {
	var actTop []*Js_ActTop
	listCityNum := GetTopCityMgr().GetCityNum()
	for i := 0; i < 3; i++ {
		actTop = append(actTop, &Js_ActTop{
			Uid:      0,
			Iconid:   0,
			Portrait: 0,
			Level:    1,
			Camp:     listCityNum[i].Camp,
			Num:      int64(listCityNum[i].Num),
			Uname:    "",
			Vip:      0,
		})
	}

	return actTop
}

// 获取cityNum
func (self *TopCityMgr) getTopRedis() {
	Try(func() {
		c := GetServer().GetRedisConn()
		defer c.Close()
		reply, err := c.Do("GET", "rankData:cityNum")
		if err != nil {
			LogError(err.Error())
			return
		}

		var topSave TopSave
		topSave.CityNum = make([]*CityNum, 0)
		if b, ok := reply.([]byte); ok {
			json.Unmarshal(b, &topSave)
			if len(topSave.CityNum) > 0 {
				self.ListCityNum = make([]*CityNum, 0)
				self.ListCityNum = append(self.ListCityNum, topSave.CityNum...)
			}
		}
	}, func(e interface{}) {
		LogError(e)
	})
}

func (self *TopCityMgr) getCityNum() map[int]int {
	GetTopCityMgr().getTopRedis()
	listCityNum := GetTopCityMgr().GetCityNum()
	res := make(map[int]int)
	for i := 0; i < len(listCityNum); i++ {
		res[listCityNum[i].Camp] = listCityNum[i].Num
	}
	return res
}
