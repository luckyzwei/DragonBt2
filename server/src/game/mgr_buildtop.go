package game

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/garyburd/redigo/redis"
)

const (
	MAX_BUILD_TOP = 7
)

//! 征收排行
type JS_TopBuild struct {
	Uid       int64  `json:"uid"`
	Name      string `json:"name"`
	Camp      int    `json:"camp"`
	Office    int    `json:"office"`
	Num       int    `json:"num"`
	Icon      int    `json:"icon"`
	Portrait  int    `json:"portrait"`   // 边框  20190412 by zy
	Vip       int    `json:"vip"`
	Level     int    `json:"level"`
	UnionName string `json:"union_name"`
}
type lstTopBuild []*JS_TopBuild

func (s lstTopBuild) Len() int           { return len(s) }
func (s lstTopBuild) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstTopBuild) Less(i, j int) bool { return s[i].Num > s[j].Num }

type TopBuildMgr struct {
	TopBuild [MAX_BUILD_TOP]lstTopBuild //! 功勋排行
	Lock     [MAX_BUILD_TOP]*sync.RWMutex
}

var topbuildmgrsingleton *TopBuildMgr = nil

func GetTopBuildMgr() *TopBuildMgr {
	if topbuildmgrsingleton == nil {
		topbuildmgrsingleton = new(TopBuildMgr)
		for i := 0; i < MAX_BUILD_TOP; i++ {
			topbuildmgrsingleton.TopBuild[i] = make(lstTopBuild, 0)
			topbuildmgrsingleton.Lock[i] = new(sync.RWMutex)
		}
	}

	return topbuildmgrsingleton
}

//! 请求redis数据
func (self *TopBuildMgr) GetData() {
	for i := 0; i < MAX_BUILD_TOP; i++ {
		c := GetServer().GetRedisConn()
		defer c.Close()

		v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("%s_%d", "topbuild", i+1)))
		if err == nil {
			json.Unmarshal(v, &self.TopBuild[i])
			//LogDebug("从redis读取", string(v), self.TopBuild[i])
		}
	}

}

//! 数据保存到redis
func (self *TopBuildMgr) SaveData() {
	c := GetServer().GetRedisConn()
	defer c.Close()
	for i := 0; i < MAX_BUILD_TOP; i++ {
		_, err := c.Do("SETEX", fmt.Sprintf("%s_%d", "topbuild", i+1), 86400, HF_JtoB(self.TopBuild[i]))
		if err != nil {
			LogError("redis fail! query:", "SET", ",err:", err)
		}
	}
}

func (self *TopBuildMgr) Refresh() {
	// 排除挖矿排行的奖励
	for i := 0; i < MAX_BUILD_TOP; i++ {
		if i == MAX_BUILD_TOP-1 {
			continue
		}
		for j := 0; j < HF_MinInt(3, len(self.TopBuild[i])); j++ {
			csv, ok := GetCsvMgr().GetHomeoffice_ranking(i+1, j+1)
			if !ok {
				continue
			}
			player := GetPlayerMgr().GetPlayer(self.TopBuild[i][j].Uid, true)
			if player == nil {
				continue
			}
			lstItem := make([]PassItem, 0)
			for k := 0; k < 4; k++ {
				itemid := HF_Atoi(csv[fmt.Sprintf("item_id%d", k+1)])
				if itemid == 0 {
					break
				}
				lstItem = append(lstItem, PassItem{itemid, HF_Atoi(csv[fmt.Sprintf("num%d", k+1)])})
			}
			(player.GetModule("mail").(*ModMail)).AddMail(1, 1, 0, csv["mail_title"], csv["mail_txt"], GetCsvMgr().GetText("STR_SYS"), lstItem, false, 0)
		}

		self.TopBuild[i] = make(lstTopBuild, 0)
	}
}

//!
func (self *TopBuildMgr) Add(index int, player *Player, num int) {
	self.Lock[index].Lock()
	defer self.Lock[index].Unlock()

	find := false
	for i := 0; i < len(self.TopBuild[index]); i++ {
		if self.TopBuild[index][i].Uid == player.Sql_UserBase.Uid {
			self.TopBuild[index][i].Num = num
			self.TopBuild[index][i].Vip = player.Sql_UserBase.Vip
			self.TopBuild[index][i].Level = player.Sql_UserBase.Level
			self.TopBuild[index][i].UnionName = player.GetUnionName()
			if i > 0 && self.TopBuild[index][i].Num < self.TopBuild[index][i-1].Num {
				return
			}
			find = true
			break
		}
	}

	if !find {
		node := new(JS_TopBuild)
		node.Uid = player.Sql_UserBase.Uid
		node.Name = player.Sql_UserBase.UName
		node.Camp = player.Sql_UserBase.Camp
		node.Num = num
		node.Icon = player.Sql_UserBase.IconId
		node.Portrait = player.Sql_UserBase.Portrait
		node.Vip = player.Sql_UserBase.Vip
		node.Level = player.GetLv()
		node.UnionName = player.GetUnionName()
		node.Office = 0
		self.TopBuild[index] = append(self.TopBuild[index], node)
	}

	sort.Sort(lstTopBuild(self.TopBuild[index]))
}

func (self *TopBuildMgr) Get(index int) []*JS_TopBuild {
	self.Lock[index].RLock()
	defer self.Lock[index].RUnlock()

	return self.TopBuild[index][:HF_MinInt(10, len(self.TopBuild[index]))]
}
func (self *TopBuildMgr) Rename(uid int64, newname string) {
	for i := 0; i < MAX_BUILD_TOP; i++ {
		for j := 0; j < len(self.TopBuild[i]); j++ {
			if self.TopBuild[i][j].Uid == uid {
				self.TopBuild[i][j].Name = newname
				break
			}
		}
	}
}
