package game

import (
	"encoding/json"
	"fmt"
	"sync"
	//"time"
)

const RANK_TASK_PLAYER_LIST = 5

type RankPlayerInfo struct {
	Uid      int64  `json:"uid"`      // uid
	Name     string `json:"name"`     // 名字
	Level    int    `json:"level"`    // 等级
	Icon     int    `json:"icon"`     // 头像
	Portrait int    `json:"portrait"` // 头像框
	Time     int64  `json:"time"`     // 时间
}

//!支援英雄
type San_MgrRankTask struct {
	ID             int    // 任务id
	RankPlayerInfo string // 玩家列表

	rankPlayerInfo []*RankPlayerInfo // 玩家列表

	DataUpdate
}

func (self *San_MgrRankTask) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.RankPlayerInfo), &self.rankPlayerInfo)
}

func (self *San_MgrRankTask) Encode() { //! 将data数据写入数据库
	self.RankPlayerInfo = HF_JtoA(&self.rankPlayerInfo)
}

type RankTaskMgr struct {
	Locker          *sync.RWMutex            //! 数据锁
	Sql_MgrRankTask map[int]*San_MgrRankTask //! 玩家支援英雄
}

var ranktaskmgrsingleton *RankTaskMgr = nil

func GetRankTaskMgr() *RankTaskMgr {
	if ranktaskmgrsingleton == nil {
		ranktaskmgrsingleton = new(RankTaskMgr)
		ranktaskmgrsingleton.Locker = new(sync.RWMutex)
		ranktaskmgrsingleton.Sql_MgrRankTask = make(map[int]*San_MgrRankTask)
	}

	return ranktaskmgrsingleton
}
func (self *RankTaskMgr) GetData() {
	var rankTask San_MgrRankTask
	sql := fmt.Sprintf("select * from `san_mgrranktask`")
	res := GetServer().DBUser.GetAllData(sql, &rankTask)
	for i := 0; i < len(res); i++ {
		data, ok1 := res[i].(*San_MgrRankTask)
		if !ok1 {
			continue
		}

		data.Init("san_mgrranktask", data, false)
		data.Decode()
		_, ok2 := self.Sql_MgrRankTask[data.ID]
		if !ok2 {
			self.Sql_MgrRankTask[data.ID] = data
		}
	}
}

// 获得数据
func (self *RankTaskMgr) GetRankTaskData(id int, create bool) *San_MgrRankTask {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	// 获得配置
	data, ok := self.Sql_MgrRankTask[id]
	if ok {
		return data
	}

	// 是否创建
	if create {
		self.Sql_MgrRankTask[id] = &San_MgrRankTask{}
		data, ok = self.Sql_MgrRankTask[id]
		if ok {
			self.Sql_MgrRankTask[id].ID = id
			self.Sql_MgrRankTask[id].rankPlayerInfo = make([]*RankPlayerInfo, 0)
			self.Sql_MgrRankTask[id].Encode()
			InsertTable("san_mgrranktask", self.Sql_MgrRankTask[id], 0, false)
			self.Sql_MgrRankTask[id].Init("san_mgrranktask", self.Sql_MgrRankTask[id], false)

			var msg S2C_RankTaskRedPoint
			msg.Cid = MSG_RANK_TASK_RED_POINT
			msg.ID = id
			GetSessionMgr().BroadCastMsg(MSG_RANK_TASK_RED_POINT, HF_JtoB(&msg))
			return data
		} else {
			return nil
		}
	}
	return nil
}

// 使用英雄
func (self *RankTaskMgr) SetRankTask(id int, player *Player) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	// 获得数据
	data, ok := self.Sql_MgrRankTask[id]
	if !ok {
		return false
	}

	if len(data.rankPlayerInfo) >= RANK_TASK_PLAYER_LIST {
		return false
	}

	index := -1
	for i, v := range data.rankPlayerInfo {
		if v.Uid == player.GetUid() {
			index = i
			break
		}
	}

	if index < 0 {
		data.rankPlayerInfo = append(data.rankPlayerInfo, &RankPlayerInfo{player.GetUid(),
			player.GetName(),
			player.Sql_UserBase.Level,
			player.Sql_UserBase.IconId,
			player.Sql_UserBase.Portrait,
			TimeServer().Unix()})

		data.Encode()
		data.Update(true)
	}
	return true
}

// 获得状态
func (self *RankTaskMgr) SetGetState() []int {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	finishState := []int{}
	for _, v := range GetCsvMgr().RankTaskConfig {
		// 获得数据
		data, ok := self.Sql_MgrRankTask[v.Id]
		if !ok || nil == data {
			finishState = append(finishState, 0)
		} else {
			finishState = append(finishState, 1)
		}
	}
	return finishState
}

// 存储数据库
func (self *RankTaskMgr) Save() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	for _, v := range self.Sql_MgrRankTask {
		v.Encode()
		v.Update(true)
	}
}

//同步改名
func (self *RankTaskMgr) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for _, v := range self.Sql_MgrRankTask {
		for _, t := range v.rankPlayerInfo {
			if t.Uid == player.GetUid() {
				t.Name = player.GetName()
			}
		}
	}
}
