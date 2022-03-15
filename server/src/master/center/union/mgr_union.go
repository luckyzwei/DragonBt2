/*
@Time : 2020/4/22 9:55
@Author : 96121
@File : mgr_union
@Software: GoLand
*/
package union

import (
	"fmt"
	"log"
	"master/db"
	"master/utils"
	"runtime/debug"
	"sync"
	"time"
)

//////////////////////////////
/*
CREATE TABLE IF NOT EXISTS tbl_union (
		id int(11) NOT NULL COMMENT '公会id',
		icon int(11) NOT NULL COMMENT 'icon',
		unionname text NOT NULL COMMENT '公会名',
		masteruid bigint(20) NOT NULL COMMENT '所有者ID',
		mastername text NOT NULL COMMENT '会长昵称',
		level int(11) NOT NULL COMMENT '公会等级',
		jointype int(11) NOT NULL COMMENT '加入类型',
		joinlevel int(11) NOT NULL COMMENT '加入等级',
		serverid int(11) NOT NULL COMMENT '服务器id',
		notice text NOT NULL COMMENT '公告',
		board text NOT NULL COMMENT '对外展示',
		createtime bigint(20) NOT NULL COMMENT '创建时间',
		lastupdtime bigint(20) NOT NULL COMMENT '更新时间',
		fight bigint(20) NOT NULL COMMENT '总战力',
		exp int(11) NOT NULL COMMENT '经验',
		dayexp int(11) NOT NULL COMMENT '每日经验',
		activitypoint int(11) NOT NULL COMMENT '活跃点数',
		acitvitylimit int(11) NOT NULL COMMENT '活跃度限额',
		mailcd bigint(20) NOT NULL COMMENT '邮件cd',
		member text NOT NULL COMMENT '成员列表',
		applys text NOT NULL COMMENT '申请列表',
		record text NOT NULL COMMENT '操作记录',
		huntinfo text NOT NULL COMMENT '军团狩猎记录',
		bravehand text NOT NULL COMMENT '无畏之手',
		PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/

const UNION_LIST_MAX = 100

type CommunityConfig struct {
	Level            int `json:"lv"`
	Exp              int `json:"exp"`
	Changeexp        int `json:"changeexp"`
	Membernum        int `json:"population"`
	Lively           int `json:"lively"`
	Changelively     int `json:"changelively"`
	Elder            int `json:"elder"`
	Fearless         int `json:"fearless"`
	Activelimit      int `json:"activelimit"`
	Warfare          int `json:"warfare"`
	GuildActiveLimit int `json:"guildactivelimit"`
	GuildExpLimit    int `json:"guildexplimit"`
}

//! 公会管理类
type UnionMgr struct {
	Sql_Union       map[int]*UnionInfo
	LastCallTime    map[int]int64 //! 招募时间保存
	LastCheckMaster map[int]int64 //! 检查军团长时间
	Locker          *sync.RWMutex //! 操作保护
	TimeLocker      *sync.RWMutex //! 时间保护
	UnionCount      int
	RefreshTime     int64

	CommunityConfigs       map[int]*CommunityConfig
}

func (self *UnionMgr) LoadCsv() {
	utils.GetCsvUtilMgr().LoadCsv("Guild_Lv", &self.CommunityConfigs)
}

var s_unionmgr *UnionMgr

func GetUnionMgr() *UnionMgr {
	if s_unionmgr == nil {
		s_unionmgr = new(UnionMgr)
		s_unionmgr.Sql_Union = make(map[int]*UnionInfo)
		s_unionmgr.LastCallTime = make(map[int]int64)
		s_unionmgr.LastCheckMaster = make(map[int]int64)
		s_unionmgr.Locker = new(sync.RWMutex)
		s_unionmgr.TimeLocker = new(sync.RWMutex)

		s_unionmgr.CommunityConfigs = make(map[int]*CommunityConfig)
		s_unionmgr.LoadCsv()
	}

	return s_unionmgr
}

//! 保存数据
func (self *UnionMgr) OnSave() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, v := range self.Sql_Union {
		v.onSave()
	}
}

//func (self *UnionMgr) GetData() {
//	var union UnionInfo
//	sql := fmt.Sprintf("select * from `tbl_union`")
//	res := db.GetDBMgr().DBUser.GetAllData(sql, &union)
//
//	for i := 0; i < len(res); i++ {
//		data := res[i].(*UnionInfo)
//		//data.record = make([]JS_UnionRecord, 0)
//		data.Init("tbl_union", data, false)
//		data.Locker = new(sync.RWMutex)
//		data.ChangeMaster.CheckPlayer = make([]int64, 0)
//		data.Decode()
//
//		self.Sql_Union[data.Id] = data
//
//		data.CheckChangeMasterState() //! 验证同步状态
//
//		//! 重新开服的时候，如果在线状态，则修改为离线
//		for j := 0; j < len(data.member); j++ {
//			if data.member[j].Lastlogintime == 0 {
//				data.member[j].Lastlogintime = time.Now().Unix()
//			}
//			member := player.GetPlayerMgr().GetPlayer(data.member[j].Uid, false)
//			passid := member.Data.PassId
//			if passid == 0 {
//				data.member[j].Stage = 110101
//			} else {
//				data.member[j].Stage = passid
//			}
//		}
//	}
//	//self.OnHunterRefresh()
//}
func (self *UnionMgr) GetAllData() {
	var union UnionInfo
	sql := fmt.Sprintf("select * from `tbl_union`")
	res := db.GetDBMgr().DBUser.GetAllData(sql, &union)

	unionMax := 0
	for i := 0; i < len(res); i++ {
		data := res[i].(*UnionInfo)
		if data.Id > 0 {
			data.record = make([]*JS_UnionRecord, 0)
			data.Locker = new(sync.RWMutex)
			data.ChangeMaster.CheckPlayer = make([]int64, 0)
			data.Decode()
			data.CheckChangeMasterState() //! 验证同步状态

			//! 重新开服的时候，如果在线状态，则修改为离线
			for j := 0; j < len(data.member); j++ {
				if data.member[j].Lastlogintime == 0 {
					data.member[j].Lastlogintime = time.Now().Unix()
				}
			}
			data.Init("tbl_union", data, false)

			self.Sql_Union[data.Id] = data

			tempMax := data.Id % 1000000
			if tempMax > unionMax {
				unionMax = tempMax
			}
		}
	}

	s_unionmgr.UnionCount = unionMax + 1
}

func (self *UnionMgr) GetData(union_uid int) *UnionInfo {
	var union UnionInfo
	sql := fmt.Sprintf("select * from `tbl_union` where id = %d", union_uid)
	db.GetDBMgr().DBUser.GetOneData(sql, &union, "", 0)

	if union.Id > 0 {
		data := &union
		data.record = make([]*JS_UnionRecord, 0)
		data.Locker = new(sync.RWMutex)
		data.ChangeMaster.CheckPlayer = make([]int64, 0)
		data.Decode()
		data.CheckChangeMasterState() //! 验证同步状态

		//! 重新开服的时候，如果在线状态，则修改为离线
		for j := 0; j < len(data.member); j++ {
			if data.member[j].Lastlogintime == 0 {
				data.member[j].Lastlogintime = time.Now().Unix()
			}
			//member := player.GetPlayerMgr().GetPlayer(data.member[j].Uid, false)
			//passid := member.Data.PassId
			//if passid == 0 {
			//	data.member[j].Stage = 110101
			//} else {
			//	data.member[j].Stage = passid
			//}
		}

		data.Init("tbl_union", data, false)
		return data
	} else {
		return nil
	}
}

func (self *UnionMgr) GetUnion(union_uid int) *UnionInfo {
	if union_uid == 0 {
		return nil
	}
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	union, ok := self.Sql_Union[union_uid]
	if !ok {
		union = self.GetData(union_uid)
		if union != nil {
			self.Sql_Union[union_uid] = union
		} else {
			return nil
		}
	}
	return union
}

func (self *UnionMgr) CheckName(unionname string, unionid int) bool {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	for _, value := range self.Sql_Union {
		if value.Unionname == unionname {
			if unionid < 0 {
				return true
			} else {
				if value.Id != unionid {
					return true
				}
			}
		}
	}
	return false
}

//创建军团
func (self *UnionMgr) CreateUnion(icon int, unionname string, uid int64, mastername string, huntInfo []*JS_UnionHunt, serverid int) int {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	union := new(UnionInfo)
	union.Icon = icon
	union.Unionname = unionname
	union.Masteruid = uid
	union.Mastername = mastername
	union.member = make([]*JS_UnionMember, 0)
	union.apply = make([]*JS_UnionApply, 0)
	union.Level = 1
	union.Locker = new(sync.RWMutex)
	union.Joinlevel = UNION_JOIN_LEVEL_BASE
	union.ServerID = serverid
	union.Exp = 0
	union.ChangeMaster.CheckPlayer = make([]int64, 0)
	union.huntInfo = huntInfo
	union.Id = serverid*1000000 + self.UnionCount
	union.Encode()

	_, ok := self.Sql_Union[union.Id]
	if ok {
		return 0
	}
	db.InsertTable("tbl_union", union, 0, false)
	union.Init("tbl_union", union, false)
	self.Sql_Union[union.Id] = union

	self.UnionCount++

	return union.Id
}

//创建军团
func (self *UnionMgr) DissolveUnion(unionid int) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	union, ok := self.Sql_Union[unionid]
	if ok {
		db.DeleteTable("tbl_union", union, []int{0})
		delete(self.Sql_Union, unionid)
	}
}

func (self *UnionMgr) GetUnionCallTime(unionid int) int64 {
	self.TimeLocker.RLock()
	defer self.TimeLocker.RUnlock()

	calltime, ok := self.LastCallTime[unionid]
	if !ok {
		return 0
	}

	return calltime
}

func (self *UnionMgr) GetUnionCheckMaster(unionid int) int64 {
	self.TimeLocker.RLock()
	defer self.TimeLocker.RUnlock()

	checktime, ok := self.LastCheckMaster[unionid]
	if !ok {
		return 0
	}

	return checktime
}

func (self *UnionMgr) SetUnionCallTime(unionid int, calltime int64) {
	self.TimeLocker.Lock()
	defer self.TimeLocker.Unlock()

	self.LastCallTime[unionid] = calltime
}

func (self *UnionMgr) SetUnionCheckMaster(unionid int, checktime int64) {
	self.TimeLocker.Lock()
	defer self.TimeLocker.Unlock()

	self.LastCheckMaster[unionid] = checktime
}

//军团检查是否满员
func (self *UnionMgr) CheckFull(unioninfo *UnionInfo) bool {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	csv_community := self.CommunityConfigs[unioninfo.Level]

	if len(unioninfo.member) >= csv_community.Membernum {
		return true
	}
	return false
}
func (self *UnionMgr) CheckBraveHand(data *UnionInfo) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	now := time.Now().Unix()
	nLen := len(data.braveHand)
	for i := nLen - 1; i >= 0; i-- {
		if data.braveHand[i].EndTime > 0 && now >= data.braveHand[i].EndTime {
			data.braveHand = append(data.braveHand[0:i], data.braveHand[i+1:]...)
		}
	}
}

// 开启战斗协程
func (self *UnionMgr) Run() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			utils.LogError(x, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			//self.CheckHunterFightEnd()
			//self.OnHunterRefresh()
			self.OnHunterRefresh()
			self.OnTime()
		}
	}
	ticker.Stop()
}

func (self *UnionMgr) OnTime() {
	tNow := time.Now()
	if tNow.Hour() == 5 && tNow.Minute() < 5 && tNow.Unix()-self.RefreshTime > 1800 {
		//!每天5：00检测，
		self.RefreshTime = tNow.Unix()
		self.OnRefresh()
	}
}

func (self *UnionMgr) OnRefresh() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, value := range self.Sql_Union {
		value.OnRefresh()
	}
}

func (self *UnionMgr) OnHunterRefresh() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	endtime := self.GetRefreshTime()
	for _, union := range self.Sql_Union {
		union.OnHunterRefresh(endtime)
		union.CheckHunterFightEnd(endtime)
		union.RefreshActivityLimit()
	}
}

// 开启战斗协程
func (self *UnionMgr) GetRefreshTime() int64 {
	if time.Now().Hour() < 5 {
		return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 5, 0, 0, 0, time.Now().Location()).Unix()
	} else {
		return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 5, 0, 0, 0, time.Now().Location()).Unix() + utils.DAY_SECS
	}
}
