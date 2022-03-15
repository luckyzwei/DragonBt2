/*
@Time : 2020/4/22 10:02 
@Author : 96121
@File : unrivaled_generals 【无双神将】 - 【诸神黄昏】
@Software: GoLand
*/
package act

import (
	"encoding/json"
	"fmt"
	"master/db"
	"master/utils"
	"sort"
	"sync"
	"time"
)

//! 消费排行榜-跨服-个人
type San_ConsumerTopUser struct {
	Uid      int64  `json:"uid"`
	SvrId    int    `json:"svrid"`
	SvrName  string `json:"svrname"`
	UName    string `json:"uname"`
	Level    int    `json:"level"`
	Vip      int    `json:"vip"`
	Icon     int    `json:"icon"`
	Portrait int    `json:"portrait"`
	Point    int    `json:"point"`
	Rank     int    `json:"rank"`
	Step     int    `json:"step"`

	db.DataUpdate
}

//! 消费排行榜-跨服-服务器
type San_ConsumerTopServer struct {
	SvrId   int    `json:"svrid"`   //! 服务器Id
	SvrName string `json:"svrname"` //! 服务器名字
	Rank    int    `json:"rank"`    //! 排名
	Point   int    `json:"point"`   //! 总积分
	Kill    int    `json:"kill"`    //! 击杀次数
	Step    int    `json:"step"`    //! 期数

	db.DataUpdate
}

//! 消费排行榜-击杀记录-稀有掉落-个人榜第一
type San_ConsumerMsg struct {
	MsgId   int    `json:"msgid"`   //! 全局消息ID
	MsgType int    `json:"msgtype"` //! 消息类型
	SvrId   string `json:"svrid"`   //! 服务器ID
	SvrName string `json:"svrname"` //! 服务器名字
	Uid     int64  `json:"uid"`     //! 玩家Id
	UName   string `json:"uname"`   //! 角色名字
	HeroId  int    `json:"heroid"`  //! 神将ID
	Level   int    `json:"level"`   //! 神将等级
	Step    int    `json:"step"`    //! 期数

	db.DataUpdate
}

//! 消费排行榜-跨服-个人
type JS_ConsumerTopUser struct {
	Uid      int64  `json:"uid"`
	SvrId    int    `json:"svrid"`
	SvrName  string `json:"svrname"`
	UName    string `json:"uname"`
	Level    int    `json:"level"`
	Vip      int    `json:"vip"`
	Icon     int    `json:"icon"`
	Point    int    `json:"point"`
	Portrait int    `json:"portrait"` // 边框  20190412 by zy
	Rank     int    `json:"rank"`
	Kill     int    `json:"kill"`
	Step     int    `json:"step"`
}

//! 消耗者排行榜
type ConsumerTopSvrMgr struct {
	Sql_TopUser      map[int64]*San_ConsumerTopUser
	Sql_TopRank      []*San_ConsumerTopUser
	Sql_TopUserBySvr map[int][]*San_ConsumerTopUser
	Sql_TopServer    []*San_ConsumerTopServer
	GlobalUser       []*San_ConsumerTopUser
	LastUpdate       int64         //! 统计时间
	Locker           *sync.RWMutex //! 数据锁
	UserSrvLocker    *sync.RWMutex //! 玩家锁
	CurStep          int           //! 当前期数
}

var consumertopsvrsingleton *ConsumerTopSvrMgr = nil

func GetConsumerTopSvr() *ConsumerTopSvrMgr {
	if consumertopsvrsingleton == nil {
		consumertopsvrsingleton = new(ConsumerTopSvrMgr)
		consumertopsvrsingleton.Locker = new(sync.RWMutex)

		consumertopsvrsingleton.Sql_TopUser = make(map[int64]*San_ConsumerTopUser)
		consumertopsvrsingleton.Sql_TopRank = make([]*San_ConsumerTopUser, 0)
		consumertopsvrsingleton.Sql_TopServer = make([]*San_ConsumerTopServer, 0)
		consumertopsvrsingleton.Sql_TopUserBySvr = make(map[int][]*San_ConsumerTopUser, 0)
		consumertopsvrsingleton.GlobalUser = make([]*San_ConsumerTopUser, 0)
		consumertopsvrsingleton.LastUpdate = 0
		consumertopsvrsingleton.UserSrvLocker = new(sync.RWMutex)
		consumertopsvrsingleton.CurStep = 0
	}

	return consumertopsvrsingleton
}

func (self *ConsumerTopSvrMgr) GetStep() int {
	step := 0
	activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
	if activity != nil {

		if len(activity.items) > 0 {
			if len(activity.items[0].N) == 4 {
				step = activity.items[0].N[3]*1000 + activity.items[0].N[2]
			}
		}
	}

	return step
}

func (self *ConsumerTopSvrMgr) GetData() {
	self.CurStep = 0
	//activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
	//if activity != nil {
	//	if len(activity.items) > 0 {
	//		if len(activity.items[0].N) > 0 {
	//			self.CurStep = activity.items[0].N[0]
	//		}
	//	}
	//}
	self.CurStep = self.GetStep()

	//! 全服数据初始化
	self.Sql_TopRank = make([]*San_ConsumerTopUser, 0)
	var topuser San_ConsumerTopUser
	sql := fmt.Sprintf("select * from `san_consumertopuser` where step = %d order by point desc", self.CurStep)
	res := db.GetDBMgr().DBUser.GetAllData(sql, &topuser)
	self.Locker.Lock()
	for i := 0; i < len(res); i++ {
		data := res[i].(*San_ConsumerTopUser)
		data.Init("san_consumertopuser", data, false)
		self.Sql_TopUser[data.Uid] = data
		self.Sql_TopRank = append(self.Sql_TopRank, data)
	}
	self.Locker.Unlock()

	self.Sql_TopServer = make([]*San_ConsumerTopServer, 0)
	var topsvr San_ConsumerTopServer
	sql1 := fmt.Sprintf("select * from `san_consumertopserver` where step = %d", self.CurStep)
	res1 := db.GetDBMgr().DBUser.GetAllData(sql1, &topsvr)
	for i := 0; i < len(res1); i++ {
		data := res1[i].(*San_ConsumerTopServer)
		data.Init("san_consumertopserver", data, false)
		self.Sql_TopServer = append(self.Sql_TopServer, data)
	}

	self.UpdateRank()
}

func (self *ConsumerTopSvrMgr) Save() {
	self.Locker.RLock()
	for _, value := range self.Sql_TopUser {
		value.UpdateEx("step", value.Step)
	}
	self.Locker.RUnlock()

	for i := 0; i < len(self.Sql_TopServer); i++ {
		self.Sql_TopServer[i].UpdateEx("step", self.Sql_TopServer[i].Step)
	}
}

func (self *ConsumerTopSvrMgr) ReloadData(step int) {
	self.Save()

	self.CurStep = step
	self.GetData()
}

//! 上传战斗数据，每次更新后上传
func (self *ConsumerTopSvrMgr) UploadDamage(top *JS_ConsumerTopUser) {
	LogDebug("上传排行数据...")
	var topsvr *San_ConsumerTopServer = nil
	for i := 0; i < len(self.Sql_TopServer); i++ {
		if self.Sql_TopServer[i].SvrId == top.SvrId {
			topsvr = self.Sql_TopServer[i]
			break
		}
	}

	if topsvr == nil {
		newsvr := new(San_ConsumerTopServer)
		newsvr.SvrId = top.SvrId
		newsvr.SvrName = top.SvrName
		newsvr.Point = top.Point
		newsvr.Step = top.Step
		db.InsertTable("san_consumertopserver", newsvr, 0, false)
		newsvr.Init("san_consumertopserver", newsvr, false)

		topsvr = newsvr
		self.Sql_TopServer = append(self.Sql_TopServer, newsvr)
	}

	self.Locker.Lock()
	finduser, ok := self.Sql_TopUser[top.Uid]
	if !ok {
		newuser := new(San_ConsumerTopUser)
		newuser.Uid = top.Uid
		newuser.Rank = 0
		newuser.SvrName = top.SvrName
		newuser.SvrId = top.SvrId
		newuser.Icon = top.Icon
		newuser.Portrait = top.Portrait
		newuser.Level = top.Level
		newuser.Point = top.Point
		newuser.Vip = top.Vip
		newuser.UName = top.UName
		newuser.Step = top.Step

		db.InsertTable("san_consumertopuser", newuser, 0, false)
		newuser.Init("san_consumertopuser", newuser, false)

		if topsvr != nil && top.Kill == 1 {
			topsvr.Kill += 1
		}

		self.Sql_TopUser[newuser.Uid] = newuser
		self.Sql_TopRank = append(self.Sql_TopRank, newuser)
		self.Locker.Unlock()
	} else {
		self.Locker.Unlock()
		if topsvr != nil {
			addpoint := top.Point - finduser.Point
			if addpoint <= 0 {
				addpoint = 0
			}

			topsvr.Point += addpoint
			if top.Kill == 1 {
				topsvr.Kill += 1
			}
		}

		finduser.Level = top.Level
		finduser.UName = top.UName
		finduser.Point = top.Point
		finduser.Icon = top.Icon
		finduser.Portrait = top.Portrait
		finduser.Vip = top.Vip
	}

	if time.Now().Unix() > self.LastUpdate+60 {
		self.LastUpdate = time.Now().Unix()

		self.UpdateRank()
	}
}

////////////////////////////////////////////////////////////////////////////////
type lstConsumerTop []*San_ConsumerTopUser

func (s lstConsumerTop) Len() int      { return len(s) }
func (s lstConsumerTop) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstConsumerTop) Less(i, j int) bool {
	if s[i].Point == s[j].Point {
		return s[i].Uid > s[j].Uid
	} else {
		return s[i].Point > s[j].Point
	}
}

type lstConsumerTopSvr []*San_ConsumerTopServer

func (s lstConsumerTopSvr) Len() int           { return len(s) }
func (s lstConsumerTopSvr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstConsumerTopSvr) Less(i, j int) bool { return s[i].Point > s[j].Point }

func (self *ConsumerTopSvrMgr) UpdateRank() {
	sort.Sort(lstConsumerTop(self.Sql_TopRank))

	self.UserSrvLocker.Lock()
	self.Sql_TopUserBySvr = make(map[int][]*San_ConsumerTopUser, 0)

	//! 写入排名数据
	for i := 0; i < len(self.Sql_TopRank); i++ {
		self.Sql_TopRank[i].Rank = i + 1
		self.Sql_TopUserBySvr[self.Sql_TopRank[i].SvrId] = append(self.Sql_TopUserBySvr[self.Sql_TopRank[i].SvrId], self.Sql_TopRank[i])
	}
	self.UserSrvLocker.Unlock()

	sort.Sort(lstConsumerTopSvr(self.Sql_TopServer))
	//! 写入服务器榜单
	for i := 0; i < len(self.Sql_TopServer); i++ {
		self.Sql_TopServer[i].Rank = i + 1
	}

	self.UpdateGlobalTopUser()
}

func (self *ConsumerTopSvrMgr) UpdateGlobalTopUser() []*San_ConsumerTopUser {
	if len(self.Sql_TopRank) < 10 {
		self.GlobalUser = make([]*San_ConsumerTopUser, 0)
		for i := 0; i < len(self.Sql_TopRank); i++ {
			self.GlobalUser = append(self.GlobalUser, self.Sql_TopRank[i])
		}
		utils.HF_DeepCopy(&self.GlobalUser, &self.Sql_TopRank)
	} else {
		self.GlobalUser = make([]*San_ConsumerTopUser, 0)
		for i := 0; i < 10; i++ {
			self.GlobalUser = append(self.GlobalUser, self.Sql_TopRank[i])
		}

		if len(self.Sql_TopRank) > 50 {
			self.GlobalUser = append(self.GlobalUser, self.Sql_TopRank[49])
		}

		if len(self.Sql_TopRank) > 100 {
			self.GlobalUser = append(self.GlobalUser, self.Sql_TopRank[99])
		}

		if len(self.Sql_TopRank) > 500 {
			self.GlobalUser = append(self.GlobalUser, self.Sql_TopRank[499])
		}

		if len(self.Sql_TopRank) > 1000 {
			self.GlobalUser = append(self.GlobalUser, self.Sql_TopRank[999])
		}

	}
	return self.GlobalUser
}

func (self *ConsumerTopSvrMgr) GetGlobalTopUser() []*San_ConsumerTopUser {
	return self.GlobalUser
}

//! 获取全服排名-个人
func (self *ConsumerTopSvrMgr) GetGlobalTopUserBySvr(serverid int) []*San_ConsumerTopUser {
	if time.Now().Unix() > self.LastUpdate+60 {
		self.LastUpdate = time.Now().Unix()

		if activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID); activity != nil {
			if activity.status.Status != ACTIVITY_STATUS_CLOSED {
				self.UpdateRank()
			}
		}
	}

	self.UserSrvLocker.RLock()
	res := self.Sql_TopUserBySvr[serverid]
	self.UserSrvLocker.RUnlock()

	return res
}

func (self *ConsumerTopSvrMgr) GetGlobalTopSvr() []*San_ConsumerTopServer {
	if time.Now().Unix() > self.LastUpdate+300 {
		self.LastUpdate = time.Now().Unix()
		self.UpdateRank()
	}

	return self.Sql_TopServer
}

func (self *ConsumerTopMgr) DoHeroTopinfo(body []byte) []byte {
	var msg S2Center_ConsumerTopGlobal
	json.Unmarshal(body, &msg)

	var data S2C_ConsumerTopGlobal
	data.Server = GetConsumerTopSvr().GetGlobalTopSvr()
	data.User = GetConsumerTopSvr().GetGlobalTopUserBySvr(msg.ServerId)
	data.Top = GetConsumerTopSvr().GetGlobalTopUser()

	return utils.HF_JtoB(&data)
}
