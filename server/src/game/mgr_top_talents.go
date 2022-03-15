package game

import (
	"encoding/json"
	"sort"
	"sync"
	//"time"
)

type TopTalentMgr struct {
	Top    [HERO_STAR_TOTAL_MAX][]*Js_ActTop
	TopCur [HERO_STAR_TOTAL_MAX]map[int64]int
	TopOld [HERO_STAR_TOTAL_MAX]map[int64]int
	Topver [HERO_STAR_TOTAL_MAX]int
	Locker *sync.RWMutex
}

var toptalent *TopTalentMgr = nil

func GetTopTalentMgr() *TopTalentMgr {
	if toptalent == nil {
		toptalent = new(TopTalentMgr)
		toptalent.Locker = new(sync.RWMutex)
		for i := 0; i < HERO_STAR_TOTAL_MAX; i++ {
			toptalent.Top[i] = make([]*Js_ActTop, 0)
			toptalent.TopCur[i] = make(map[int64]int, 0)
			toptalent.TopOld[i] = make(map[int64]int, 0)
		}
	}
	return toptalent
}

// 天赋总星级
func (self *TopTalentMgr) GetData(unionName map[int64]string) {
	self.Locker.Lock()
	for i := 0; i < HERO_STAR_TOTAL_MAX; i++ {
		if self.Topver[i] > 0 {
			HF_DeepCopy(&self.TopOld[i], &self.TopCur[i])
			self.Top[i] = make([]*Js_ActTop, 0)
			self.TopCur[i] = make(map[int64]int)
		}
	}

	var top Js_ActTopLoad
	sql := "select a.uid,a.uname,a.iconid,a.portrait,a.level,a.camp, a.fight,a.vip,b.totalstars,'0' as union_name" +
		" from san_userbase as a JOIN  san_userhero2 as b where a.fight > 0 and a.uid = b.uid"

	res := GetServer().DBUser.GetAllDataEx(sql, &top)
	for i := 0; i < len(res); i++ {
		data := res[i].(*Js_ActTopLoad)
		totalStars := make([]*HeroTopInfo, 0)
		json.Unmarshal([]byte(data.Nums), &totalStars)
		if len(totalStars) != HERO_STAR_TOTAL_MAX {
			continue
		}

		for t := 0; t < HERO_STAR_TOTAL_MAX; t++ {
			temp := Js_ActTop{}
			temp.Uid = data.Uid
			temp.Uname = data.Uname
			temp.Iconid = data.Iconid
			temp.Portrait = data.Portrait
			temp.Level = data.Level
			temp.Camp = data.Camp
			temp.Fight = data.Fight
			temp.Vip = data.Vip
			temp.UnionName = data.UnionName
			temp.LastRank = data.LastRank
			temp.StartTime = totalStars[t].StarTime
			temp.Num = int64(totalStars[t].Stars)
			if temp.Num <= 0 {
				continue
			}

			lastRank, ok := self.TopCur[t][temp.Uid]
			if ok {
				temp.LastRank = lastRank
			}
			// 更新军团名字
			v, ok := unionName[temp.Uid]
			if ok {
				temp.UnionName = v
			}
			self.Top[t] = append(self.Top[t], &temp)
			self.Topver[t]++
		}
	}

	for i := 0; i < HERO_STAR_TOTAL_MAX; i++ {
		sort.Sort(lstJsActTop(self.Top[i]))
		for t := 0; t < len(self.Top[i]); t++ {
			self.Top[i][t].LastRank = t + 1
			self.TopCur[i][self.Top[i][t].Uid] = t + 1
		}
	}

	self.Locker.Unlock()
}

func (self *TopTalentMgr) GetTopShow(topType int) ([]*Js_ActTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	nType := self.GetTopType(topType)
	if len(self.Top[nType]) > MaxRankShowNum {
		return self.Top[nType][:MaxRankShowNum], self.Topver[nType]
	}
	return self.Top[nType], self.Topver[nType]
}

//获取钻石消耗
func (self *TopTalentMgr) GetTop(topType int) ([]*Js_ActTop, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	nType := self.GetTopType(topType)
	if len(self.Top[nType]) > MaxRankNum {
		return self.Top[nType][:MaxRankNum], self.Topver[nType]
	}
	return self.Top[nType], self.Topver[nType]
}

func (self *TopTalentMgr) GetTopCurNum(topType int, id int64) int {
	self.Locker.RLock()
	nType := self.GetTopType(topType)
	data, _ := self.TopCur[nType][id]
	self.Locker.RUnlock()
	return data
}

func (self *TopTalentMgr) GetTopOldNum(topType int, id int64) int {
	self.Locker.RLock()
	nType := self.GetTopType(topType)
	data, _ := self.TopOld[nType][id]
	self.Locker.RUnlock()
	return data
}

// 更新排行数据
func (self *TopTalentMgr) UpdateRank(topType int, count int, player *Player) {
	self.Locker.Lock()

	insert := true  // 是否插入新数据
	change := false // 是否重新排序
	for i := 0; i < len(self.Top[topType]); i++ {
		if self.Top[topType][i].Uid == player.Sql_UserBase.Uid {
			self.Top[topType][i].Num = int64(count)
			self.Top[topType][i].Level = player.Sql_UserBase.Level
			insert = false
			if i > 0 {
				if self.Top[topType][i-1].Num < int64(count) {
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

		data.Fight = player.Sql_UserBase.Fight
		data.Vip = player.Sql_UserBase.Vip
		data.UnionName = player.GetUnionName()
		data.StartTime = TimeServer().Unix()
		data.Num = int64(count)
		self.Top[topType] = append(self.Top[topType], &data)
		self.TopCur[topType][data.Uid] = 0
	}

	if change == true || insert == true {
		sort.Sort(lstJsActTop(self.Top[topType]))
		for i := 0; i < len(self.Top[topType]); i++ {
			if self.Top[topType][i] != nil {
				self.Top[topType][i].LastRank = i + 1
				self.TopCur[topType][self.Top[topType][i].Uid] = i + 1
			}
		}
		self.Topver[topType]++
	}

	if len(self.Top[topType]) > 200 {
		self.Top[topType] = self.Top[topType][:200]
	}

	self.Locker.Unlock()
}

func (self *TopTalentMgr) GetTopType(topType int) int {
	nType := 0
	switch topType {
	case TOP_RANK_HERO_TALENT:
		nType = HERO_STAR_TOTAL
	case TOP_RANK_HERO_TALENT_CAMP1:
		nType = HERO_STAR_TOTAL_CAMP1
	case TOP_RANK_HERO_TALENT_CAMP2:
		nType = HERO_STAR_TOTAL_CAMP2
	case TOP_RANK_HERO_TALENT_CAMP3:
		nType = HERO_STAR_TOTAL_CAMP3
	case TOP_RANK_HERO_TALENT_CAMP4:
		nType = HERO_STAR_TOTAL_CAMP4
	}
	return nType
}

//排行同步改名
func (self *TopTalentMgr) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for i := 0; i < HERO_STAR_TOTAL_MAX; i++ {
		for t := 0; t < len(self.Top[i]); t++ {
			if self.Top[i][t].Uid == player.Sql_UserBase.Uid {
				self.Top[i][t].Uname = player.Sql_UserBase.UName
				break
			}
		}
	}
}
