package game

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"sort"
	"sync"
	"time"
)

type CVarList struct {
	Data []interface{}
}

func (self *CVarList) IntVal(index int) int {
	if index < 0 || index >= len(self.Data) {
		return 0
	}
	return self.Data[index].(int)
}
func (self *CVarList) Int64Val(index int) int64 {
	if index < 0 || index >= len(self.Data) {
		return 0
	}
	return self.Data[index].(int64)
}
func (self *CVarList) StringVal(index int) string {
	if index < 0 || index >= len(self.Data) {
		return ""
	}
	return self.Data[index].(string)
}

func (self *CVarList) AddData(data interface{}) *CVarList {
	self.Data = append(self.Data, data)
	return self
}

const (
	UNION_GET_TIME_TYPE_CALL_TIME  = 1
	UNION_GET_TIME_TYPE_CHECK_TIME = 2
)

const (
	UNION_HUNT_TYPE_NOMAL   = 1
	UNION_HUNT_TYPE_SPECIAL = 2
	UNION_HUNT_TYPE_MAX     = 3
)

const (
	UNION_POSITION_MASTER      = 1
	UNION_POSITION_VICE_MASTER = 2
	UNION_POSITION_ELITE       = 3
	UNION_POSITION_MEMBER      = 4
)

//! 0-加入 1-退出 2-改变职务 4-创建军团 5 解散 6-更名 7-更改标志 8-军团长更替 9军团boss过期 10狩猎获得钻石礼包
const (
	UNION_RECORD_TYPE_JOIN                = 0
	UNION_RECORD_TYPE_OUT                 = 1
	UNION_RECORD_TYPE_MODIFY              = 2
	UNION_RECORD_TYPE_CREATE              = 4
	UNION_RECORD_TYPE_DISSOLVE            = 5
	UNION_RECORD_TYPE_UNION_CHANGE_NAME   = 6
	UNION_RECORD_TYPE_UNION_CHANGE_ICON   = 7
	UNION_RECORD_TYPE_UNION_CHANGE_MASTER = 8
	UNION_RECORD_TYPE_HUNTER_BOSS_LEAVE   = 9
	UNION_RECORD_TYPE_HUNTER_BOSS_AWARD   = 10
)

const UNION_JOIN_LEVEL_BASE = 10

const (
	UNION_UPDATE_CD = 1800 //! 更新时间间隔 - 30分钟
)

// 离线邮件存放到一个字段
type San_Union struct {
	Id            int
	Icon          int
	Unionname     string
	Masteruid     int64
	Mastername    string
	Level         int
	Jointype      int
	Joinlevel     int
	ServerID      int //! 服务器id
	Member        string
	Apply         string
	Notice        string
	Createtime    int64
	Lastupdtime   int64
	Fight         int64
	Exp           int
	DayExp        int
	Record        string
	ActivityPoint int //活跃点数
	AcitvityLimit int //活跃度限额
	HuntInfo      string
	Board         string
	MailCD        int64
	BraveHand     string

	member    []JS_Member
	apply     []JS_UnionApply
	record    []JS_UnionRecord
	huntInfo  []*JS_UnionHunt
	braveHand []*JS_UnionBraveHand

	Locker *sync.RWMutex //! 操作锁

	ChangeMaster JS_UnionChangeMaster //! 军团长自动更换

	lastUpdate int64
	DataUpdate
}

type JS_UnionBraveHand struct {
	Uid     int64
	EndTime int64 `json:"endtime"` //! 结束时间
}

type JS_UnionHunt struct {
	Type    int                `json:"type"`    //! 类型
	TopDps  lstUnionHuntDpsTop `json:"topdps"`  //! 伤害排行
	EndTime int64              `json:"endtime"` //! 结束时间
}

//! 军团记录
type JS_UnionRecord struct {
	Type  int    `json:"type"` //! 0-加入 1-退出 2-改变职务 3-军团更名 6-更名 7-更改标志 8-军团长更替
	Time  int64  `json:"time"`
	Name  string `json:"name"`
	Param string `json:"param"`
}

type JS_UnionChangeMaster struct {
	CheckPlayer []int64 //! 同步过的玩家
	CheckFlag   bool    //! 同步标志
	OldMaster   string  //! 老军团长
	NowMaster   string  //!
}

//! 军团副本
type JS_UnionCopy struct {
	Time   int64              `json:"time"`   //! 刷新时间
	Son    []*JS_UnionCopySon `json:"son"`    //! 进度
	TopNum lstUnionCopyNumTop `json:"topnum"` //! 次数排行
}
type JS_UnionCopySon struct {
	Id      int                   `json:"id"`      //! 副本id
	MaxHp   int                   `json:"maxhp"`   //! 进度
	Time    int64                 `json:"time"`    //! 正在被打的时间
	Uid     int64                 `json:"uid"`     //! 正在被打的人
	Monster []JS_UnionCopyMonster `json:"monster"` //! 怪物
	TopDps  lstUnionCopyDpsTop    `json:"topdps"`  //! 伤害排行
}

func (self *JS_UnionCopySon) GetProgress() int {
	hp := 0
	for i := 0; i < len(self.Monster); i++ {
		hp += self.Monster[i].Hp
	}
	if hp > self.MaxHp {
		hp = self.MaxHp
	}

	ret := (self.MaxHp - hp) * 100 / self.MaxHp
	if ret == 0 && hp < self.MaxHp {
		return 1
	} else if ret == 100 && hp > 0 {
		return 99
	}

	return ret
}

type UnionMail struct {
	Title string `json:"titile"`
	Text  string `json:"text"`
}

type JS_UnionCopyMonster struct {
	Id int `json:"id"` //! 怪物
	Hp int `json:"hp"` //! hp
}
type JS_UnionCopyNumTop struct {
	Uid      int64  `json:"uid"`      //! id
	Name     string `json:"name"`     //! 名字
	Icon     int    `json:"icon"`     //! icon
	Portrait int    `json:"portrait"` // 边框  20190412 by zy
	Level    int    `json:"level"`    //! 等级
	Job      int    `json:"job"`      //! 职位
	Fight    int64  `json:"fight"`    //! 战斗力
	Vip      int    `json:"vip"`      //! vip
	Num      int    `json:"num"`      //! 次数
}
type lstUnionCopyNumTop []*JS_UnionCopyNumTop

func (s lstUnionCopyNumTop) Len() int           { return len(s) }
func (s lstUnionCopyNumTop) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstUnionCopyNumTop) Less(i, j int) bool { return s[i].Num > s[j].Num }

type JS_UnionCopyDpsTop struct {
	Uid      int64  `json:"uid"`      //! id
	Name     string `json:"name"`     //! 名字
	Icon     int    `json:"icon"`     //! icon
	Portrait int    `json:"portrait"` // 边框  20190412 by zy
	Level    int    `json:"level"`    //! 等级
	Job      int    `json:"job"`      //! 职位
	Fight    int64  `json:"fight"`    //! 战斗力
	Vip      int    `json:"vip"`      //! vip
	Dps      int64  `json:"dps"`      //! 伤害
}
type lstUnionCopyDpsTop []*JS_UnionCopyDpsTop

func (s lstUnionCopyDpsTop) Len() int           { return len(s) }
func (s lstUnionCopyDpsTop) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstUnionCopyDpsTop) Less(i, j int) bool { return s[i].Dps > s[j].Dps }

type JS_UnionHuntDpsTop struct {
	Uid      int64  `json:"uid"`      //! id
	Name     string `json:"name"`     //! 名字
	Icon     int    `json:"icon"`     //! icon
	Portrait int    `json:"portrait"` // 边框  20190412 by zy
	Level    int    `json:"level"`    //! 等级
	Job      int    `json:"job"`      //! 职位
	Fight    int64  `json:"fight"`    //! 战斗力
	Vip      int    `json:"vip"`      //! vip
	Dps      int64  `json:"dps"`      //! 伤害
	FightID  int64  `json:"fightid"`  //! 战斗id
	Time     int64  `json:"time"`     //! 时间
}

type lstUnionHuntDpsTop []*JS_UnionHuntDpsTop

func (s lstUnionHuntDpsTop) Len() int           { return len(s) }
func (s lstUnionHuntDpsTop) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstUnionHuntDpsTop) Less(i, j int) bool { return s[i].Dps > s[j].Dps }

type JS_Union struct {
	Id            int              `json:"id"`
	Icon          int              `json:"icon"`
	Unionname     string           `json:"unionname"`
	Masteruid     int64            `json:"masteruid"`
	Mastername    string           `json:"mastername"`
	Level         int              `json:"level"`
	Jointype      int              `json:"jointype"`
	Joinlevel     int              `json:"joinlevel"`
	Exp           int              `json:"exp"`
	DayExp        int              `json:"dayexp"`
	Member        []JS_Member      `json:"member"`
	Apply         []JS_UnionApply  `json:"apply"`
	State         int              `json:"state"`
	Cityinfo      string           `json:"cityinfo"`
	Notice        string           `json:"notice"`
	Createtime    int64            `json:"createtime"`
	Lastupdtime   int64            `json:"lastupdtime"`
	Record        []JS_UnionRecord `json:"record"`
	CopyUpdate    int64            `json:"copyupdate"`
	TotalFight    int64            `json:"total_fight"`
	Rank          int              `json:"rank"`
	ActivityPoint int              `json:"activitypoint"`
	HuntInfo      []*JS_UnionHunt  `json:"huntInfo"`
	Board         string           `json:"board"`
}

type JS_Union2 struct {
	Id            int    `json:"id"`
	Icon          int    `json:"icon"`
	Unionname     string `json:"unionname"`
	Masteruid     int64  `json:"masteruid"`
	Mastername    string `json:"mastername"`
	Level         int    `json:"level"`
	Jointype      int    `json:"jointype"`
	Joinlevel     int    `json:"joinlevel"`
	Money         int    `json:"money"`
	Member        int    `json:"member"`
	State         int    `json:"state"`
	Camp          int    `json:"camp"`
	Fight         int64  `json:"fight"`
	Exp           int    `json:"exp"`
	ActivityPoint int    `json:"activitypoint"`
}
type UserActivityRecord struct {
	Time     int64 `json:"time"`     // 时间
	AddCount int   `json:"addcount"` // 数量
}

type JS_Member struct {
	Uid            int64                 `json:"uid"`
	Level          int                   `json:"level"`
	Uname          string                `json:"uname"`
	Iconid         int                   `json:"iconid"`
	Portrait       int                   `json:"portrait"`
	Vip            int                   `json:"vip"`
	Position       int                   `json:"position"`
	Camp           int                   `json:"camp"`
	Fight          int64                 `json:"fight"`
	Donation       int                   `json:"donation"`
	TodayDon       int                   `json:"todaydon"`
	CopyNum        int                   `json:"copynum"`
	Lastlogintime  int64                 `json:"lastlogintime"`
	ActivityRecord []*UserActivityRecord `json:"activityrecord"` // 记录
	BraveHand      int                   `json:"bravehand"`      // 无畏之手
	Stage          int                   `json:"stage"`          // 关卡进度
}

type JS_UnionApply struct {
	Uid       int64  `json:"uid"`
	Level     int    `json:"level"`
	Uname     string `json:"uname"`
	Iconid    int    `json:"iconid"`
	Portrait  int    `json:"portrait"`
	Vip       int    `json:"vip"`
	Fight     int64  `json:"fight"`
	Applytime int64  `json:"lastlogintime"`
}

func (self *JS_Union) Init() {

}

func (self *San_Union) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Member), &self.member)
	json.Unmarshal([]byte(self.Apply), &self.apply)
	json.Unmarshal([]byte(self.Record), &self.record)
	json.Unmarshal([]byte(self.HuntInfo), &self.huntInfo)
	json.Unmarshal([]byte(self.BraveHand), &self.braveHand)
}

func (self *San_Union) Encode() { //! 将data数据写入数据库
	self.Locker.RLock()
	self.Member = HF_JtoA(&self.member)
	self.Apply = HF_JtoA(&self.apply)
	self.Record = HF_JtoA(&self.record)
	self.HuntInfo = HF_JtoA(&self.huntInfo)
	self.BraveHand = HF_JtoA(&self.braveHand)
	self.Locker.RUnlock()
}

type UnionMgr struct {
	Sql_Union     map[int]*San_Union
	Sql_UnionName map[int]string
	Locker        *sync.RWMutex //! 操作保护
	NameLocker    *sync.RWMutex //! 操作保护
}

var unionmgrsingleton *UnionMgr = nil

//! public
func GetUnionMgr() *UnionMgr {
	if unionmgrsingleton == nil {
		unionmgrsingleton = new(UnionMgr)
		unionmgrsingleton.Sql_Union = make(map[int]*San_Union)
		unionmgrsingleton.Sql_UnionName = make(map[int]string)
		unionmgrsingleton.Locker = new(sync.RWMutex)
		unionmgrsingleton.NameLocker = new(sync.RWMutex)
	}

	return unionmgrsingleton
}

func (self *UnionMgr) GetData() {
	//var city San_Union
	//sql := fmt.Sprintf("select * from `san_unioninfo`")
	//res := GetServer().DBUser.GetAllData(sql, &city)
	//
	//for i := 0; i < len(res); i++ {
	//	data := res[i].(*San_Union)
	//	//data.record = make([]JS_UnionRecord, 0)
	//	data.Init("san_unioninfo", data, false)
	//	data.Locker = new(sync.RWMutex)
	//	data.CopyLock = new(sync.RWMutex)
	//	data.ChangeMaster.CheckPlayer = make([]int64, 0)
	//	data.Decode()
	//
	//	self.Sql_Union[data.Id] = data
	//
	//	data.CalcFight()              //! 重新计算战斗力
	//	data.CheckChangeMasterState() //! 验证同步状态
	//
	//	//! 重新开服的时候，如果在线状态，则修改为离线
	//	for j := 0; j < len(data.member); j++ {
	//		if data.member[j].Lastlogintime == 0 {
	//			data.member[j].Lastlogintime = TimeServer().Unix()
	//		}
	//		passid, _ := GetOfflineInfoMgr().GetBaseInfo(data.member[j].Uid)
	//		if passid == 0 {
	//			data.member[j].Stage = ONHOOK_INIT_LEVEL
	//		} else {
	//			data.member[j].Stage = passid
	//		}
	//	}
	//
	//	if TimeServer().Unix() > data.copy.Time {
	//		data.RefreshCopy()
	//	}
	//
	//	if len(data.huntInfo) <= 0 {
	//		for i := UNION_HUNT_TYPE_NOMAL; i < UNION_HUNT_TYPE_MAX; i++ {
	//			config := GetCsvMgr().GetUnionHuntConfigByID(i)
	//			if nil == config {
	//				continue
	//			}
	//
	//			temp := JS_UnionHunt{}
	//			temp.Type = i
	//			temp.TopDps = lstUnionHuntDpsTop{}
	//			if i == UNION_HUNT_TYPE_NOMAL {
	//				temp.EndTime = data.GetRefreshTime()
	//			}
	//			data.huntInfo = append(data.huntInfo, &temp)
	//		}
	//	}
	//}
	//self.OnHunterRefresh()
}

func (self *UnionMgr) Save() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, value := range self.Sql_Union {
		value.Encode()
		value.Update(true)
	}
}

func (self *UnionMgr) UpdateUnion(id int) *San_Union {
	if id == 0 {
		fmt.Printf("!!!!!!")
	}
	//! 超时，或者未找到，则重新请求
	if GetMasterMgr().UnionRPC.Client != nil {
		var msg S2M_UnionGetUnion
		msg.Unionuid = id

		ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GET_UNION, msg)
		if ret == nil || ret.RetCode != UNION_SUCCESS {
			return nil
		}

		if ret.Data == "" {
			return nil
		}

		var backmsg M2S_UnionGetUnion
		json.Unmarshal([]byte(ret.Data), &backmsg)
		self.Locker.RLock()
		value, ok := self.Sql_Union[id]
		self.Locker.RUnlock()
		if ok {
			value.Id = backmsg.Data.Id
			value.Icon = backmsg.Data.Icon
			value.Unionname = backmsg.Data.Unionname
			value.Masteruid = backmsg.Data.Masteruid
			value.Mastername = backmsg.Data.Mastername
			value.Level = backmsg.Data.Level
			value.Jointype = backmsg.Data.Jointype
			value.Joinlevel = backmsg.Data.Joinlevel
			value.ServerID = backmsg.Data.ServerID
			value.Member = backmsg.Data.Member
			value.Apply = backmsg.Data.Applys
			value.Notice = backmsg.Data.Notice
			value.Createtime = backmsg.Data.Createtime
			value.Lastupdtime = backmsg.Data.Lastupdtime
			value.Fight = backmsg.Data.Fight
			value.Exp = backmsg.Data.Exp
			value.DayExp = backmsg.Data.DayExp
			value.Record = backmsg.Data.Record
			value.ActivityPoint = backmsg.Data.ActivityPoint
			value.AcitvityLimit = backmsg.Data.AcitvityLimit
			value.HuntInfo = backmsg.Data.HuntInfo
			value.Board = backmsg.Data.Board
			value.MailCD = backmsg.Data.MailCD
			value.BraveHand = backmsg.Data.BraveHand
			value.Decode()

			return value
		} else {
			value := new(San_Union)
			value.Id = backmsg.Data.Id
			value.Icon = backmsg.Data.Icon
			value.Unionname = backmsg.Data.Unionname
			value.Masteruid = backmsg.Data.Masteruid
			value.Mastername = backmsg.Data.Mastername
			value.Level = backmsg.Data.Level
			value.Jointype = backmsg.Data.Jointype
			value.Joinlevel = backmsg.Data.Joinlevel
			value.ServerID = backmsg.Data.ServerID
			value.Member = backmsg.Data.Member
			value.Apply = backmsg.Data.Applys
			value.Notice = backmsg.Data.Notice
			value.Createtime = backmsg.Data.Createtime
			value.Lastupdtime = backmsg.Data.Lastupdtime
			value.Fight = backmsg.Data.Fight
			value.Exp = backmsg.Data.Exp
			value.DayExp = backmsg.Data.DayExp
			value.Record = backmsg.Data.Record
			value.ActivityPoint = backmsg.Data.ActivityPoint
			value.AcitvityLimit = backmsg.Data.AcitvityLimit
			value.HuntInfo = backmsg.Data.HuntInfo
			value.Board = backmsg.Data.Board
			value.MailCD = backmsg.Data.MailCD
			value.BraveHand = backmsg.Data.BraveHand
			value.member = make([]JS_Member, 0)
			value.record = make([]JS_UnionRecord, 0)
			value.apply = make([]JS_UnionApply, 0)
			value.huntInfo = make([]*JS_UnionHunt, 0)
			value.braveHand = make([]*JS_UnionBraveHand, 0)
			value.Locker = new(sync.RWMutex)
			value.ChangeMaster.CheckPlayer = make([]int64, 0)

			value.Decode()

			self.Locker.Lock()
			self.Sql_Union[id] = value
			self.Locker.Unlock()

			self.NameLocker.Lock()
			self.Sql_UnionName[id] = value.Unionname
			self.NameLocker.Unlock()
			return value
		}
	}
	return nil
}

func (self *UnionMgr) GetUnion(id int) *San_Union {
	if id == 0 {
		return nil
	}
	self.Locker.RLock()
	value, ok := self.Sql_Union[id]
	self.Locker.RUnlock()
	if ok {
		if value.lastUpdate >= TimeServer().Unix()-UNION_UPDATE_CD {
			//! 一小时内更新过，则直接返回
			return value
		}
		return value
	}

	return self.UpdateUnion(id)
}

func (self *UnionMgr) GetUnionNum(id int) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	value, ok := self.Sql_Union[id]
	if ok {
		return len(value.member)
	}
	return 0
}

func (self *UnionMgr) GetUnionByName(name string) []JS_Union2 {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	unionList := []JS_Union2{}
	var mastermsg S2M_UnionGetUnionByName
	mastermsg.Name = name
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GET_UNION_BY_NAME, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return unionList
	}

	var backmsg M2S_UnionGetUnionByName
	json.Unmarshal([]byte(ret.Data), &backmsg)
	return backmsg.Data
}

func (self *UnionMgr) GetUnionCallTime(unionid int) int64 {
	var msg S2M_UnionGetTime
	msg.Unionuid = unionid
	msg.Type = UNION_GET_TIME_TYPE_CALL_TIME
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GET_TIME, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return 0
	}

	var backmsg M2S_UnionGetTime
	json.Unmarshal([]byte(ret.Data), &backmsg)
	return backmsg.Time
}

func (self *UnionMgr) GetUnionName(unionid int) string {
	self.NameLocker.RLock()
	defer self.NameLocker.RUnlock()

	value, ok := self.Sql_UnionName[unionid]
	if !ok {
		return ""
	}
	return value

	//self.Locker.RLock()
	//defer self.Locker.RUnlock()
	//
	//value, ok := self.Sql_Union[unionid]
	//if !ok {
	//	return ""
	//}
	//return value.Unionname
}

func (self *UnionMgr) GetMasterName(unionid int) string {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	value, ok := self.Sql_Union[unionid]
	if !ok {
		return ""
	}
	return value.Mastername
}

func (self *UnionMgr) GetUnionIcon(unionid int) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	value, ok := self.Sql_Union[unionid]
	if !ok {
		return 0
	}
	return value.Icon
}

func (self *UnionMgr) GetUnionLv(unionid int) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	value, ok := self.Sql_Union[unionid]
	if !ok {
		return 0
	}
	return value.Level
}

func (self *UnionMgr) SetUnionCallTime(unionid int, calltime int64) {
	var msg S2M_UnionSetTime
	msg.Unionuid = unionid
	msg.Type = UNION_GET_TIME_TYPE_CALL_TIME
	msg.Time = TimeServer().Unix()
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_SET_TIME, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}
}

func (self *UnionMgr) GetUnionCheckMaster(unionid int) int64 {
	var msg S2M_UnionGetTime
	msg.Unionuid = unionid
	msg.Type = UNION_GET_TIME_TYPE_CHECK_TIME
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GET_TIME, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return 0
	}

	var backmsg M2S_UnionGetTime
	json.Unmarshal([]byte(ret.Data), &backmsg)
	return backmsg.Time
}

func (self *UnionMgr) SetUnionCheckMaster(unionid int, checktime int64) {
	var msg S2M_UnionSetTime
	msg.Unionuid = unionid
	msg.Type = UNION_GET_TIME_TYPE_CHECK_TIME
	msg.Time = TimeServer().Unix()
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_SET_TIME, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}
}

func (self *UnionMgr) GetUnionJsInfo(id int) *JS_Union {
	union := self.GetUnion(id)
	if union == nil {
		return nil
	}

	var ret JS_Union
	ret.Id = union.Id
	ret.Icon = union.Icon
	ret.Unionname = union.Unionname
	ret.Masteruid = union.Masteruid
	ret.Mastername = union.Mastername
	ret.Level = union.Level
	ret.Jointype = union.Jointype
	ret.Joinlevel = union.Joinlevel
	ret.Exp = union.Exp
	ret.DayExp = union.DayExp
	ret.Member = union.member
	ret.Apply = union.apply
	ret.Notice = union.Notice
	ret.Board = union.Board
	ret.Createtime = union.Createtime
	ret.Lastupdtime = union.Lastupdtime
	ret.Record = union.record
	ret.HuntInfo = union.huntInfo
	ret.TotalFight = 0
	for _, v := range union.member {
		ret.TotalFight += v.Fight
	}
	ret.Rank = GetTopUnionMgr().GetUnionRank(union.Id)
	ret.ActivityPoint = union.ActivityPoint

	return &ret
}

func (self *UnionMgr) GetUnionJsInfo2(id int) *JS_Union2 {
	union := self.GetUnion(id)

	if union == nil {
		return nil
	}

	var ret JS_Union2
	ret.Id = union.Id
	ret.Icon = union.Icon
	ret.Unionname = union.Unionname
	ret.Masteruid = union.Masteruid
	ret.Mastername = union.Mastername
	ret.Level = union.Level
	ret.Jointype = union.Jointype
	ret.Joinlevel = union.Joinlevel
	ret.Member = len(union.member)
	ret.Camp = union.GetCamp()
	ret.Fight = union.Fight
	ret.Exp = union.Exp
	ret.ActivityPoint = union.ActivityPoint

	return &ret
}

func (self *San_Union) AddApply(member JS_UnionApply) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	for i := 0; i < len(self.apply); i++ {
		if self.apply[i].Uid == member.Uid {
			self.apply[i].Level = member.Level
			self.apply[i].Vip = member.Vip
			self.apply[i].Fight = member.Fight
			return true
		}
	}
	self.apply = append(self.apply, member)

	return true
}

func (self *San_Union) GetCamp() int {
	//self.Locker.RLock()
	//defer self.Locker.RUnlock()

	if len(self.member) > 0 {
		return self.member[0].Camp
	}

	//for i := 0; i < len(self.member); i++ {
	//	if self.member[i].Position == 1 {
	//		return self.member[i].Camp
	//	}
	//}

	return CAMP_SHU
}

func (self *San_Union) AddMailCD() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	self.MailCD = TimeServer().Unix()
}

func (self *San_Union) GetCampSafe() int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	return self.GetCamp()
}

func (self *UnionMgr) AddApply(unionid int, member JS_UnionApply) bool {
	union := self.GetUnion(unionid)
	if union == nil {
		return false
	}

	var msg S2M_UnionApply
	msg.Uid = member.Uid
	msg.Unionuid = unionid

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_APPLY, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	union.AddApply(member)

	return true
}

func (self *San_Union) CancelApply(uid int64) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	for j := 0; j < len(self.apply); j++ {
		if self.apply[j].Uid == uid {
			copy(self.apply[j:], self.apply[j+1:])
			self.apply = self.apply[:len(self.apply)-1]
			return true
		}
	}

	return false
}

//撤销军团申请
func (self *UnionMgr) CancelApply(unionid int, uid int64) bool {
	union := self.GetUnion(unionid)
	if union == nil {
		return false
	}

	var msg S2M_UnionCancelApply
	msg.Uid = uid
	msg.Unionuid = unionid

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_CANCEL_APPLY, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	return union.CancelApply(uid)
}

//删除军团申请
func (self *UnionMgr) CleanApply(unionid int, uid int64) bool {
	union := self.GetUnion(unionid)
	if union == nil {
		return false
	}

	if uid == 0 {
		union.apply = []JS_UnionApply{}
	} else {
		for i, v := range union.apply {
			if v.Uid == uid {
				union.apply = append(union.apply[0:i], union.apply[i+1:]...)
				break
			}
		}
	}

	return true
}

//检查军团名字
func (self *UnionMgr) CheckName(unionname string) bool {
	self.NameLocker.RLock()
	defer self.NameLocker.RUnlock()

	for _, name := range self.Sql_UnionName {
		if name == unionname {
			return true
		}
	}

	/*self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, value := range self.Sql_Union {
		if value.Unionname == unionname {
			return true
		}
	}*/

	return false
}

//! 检查会长数据丢失
func (self *UnionMgr) CheckMaster(unionid int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	find := false
	for _, value := range self.Sql_Union {
		if value.Id == unionid {
			for i := 0; i < len(value.member); i++ {
				if value.member[i].Uid == value.Masteruid {
					if value.member[i].Position != UNION_POSITION_MASTER {
						value.member[i].Position = UNION_POSITION_MASTER
						find = true
						break
					}
					find = true
				}
			}

			if find == false {
				value.Masteruid = value.member[0].Uid
				value.member[0].Position = UNION_POSITION_MASTER
			}
		}
	}
}

func (self *San_Union) CheckMasterOffline() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	change := false
	hasmaster := false
	masteridx := -1
	for j := 0; j < len(self.member); j++ {
		if self.member[j].Uid == self.Masteruid {
			hasmaster = true
			offtime := TimeServer().Unix() - self.member[j].Lastlogintime
			if self.member[j].Lastlogintime > 0 && offtime > DAY_SECS*3 {
				masteridx = j
				change = true
				break
			}
		}
	}

	//! 触发军团长转移-太久没有上线，或者没有军团长
	if change == true || hasmaster == false {
		maxlevel := 0
		idx := 0
		var fight int64 = 0
		for i := 0; i < len(self.member); i++ {
			if self.member[i].Uid == self.Masteruid {
				continue
			}
			if self.member[i].Lastlogintime == 0 {
				if maxlevel == 0 {
					maxlevel = self.member[i].Position
					fight = self.member[i].Fight
					idx = i
				}
				if self.member[i].Position < maxlevel {
					maxlevel = self.member[i].Position
					fight = self.member[i].Fight
					idx = i
				} else if self.member[i].Position == maxlevel {
					if self.member[i].Fight > fight {
						idx = i
						fight = self.member[i].Fight
					}
				}
			}
		}

		if maxlevel > 0 {
			if hasmaster {
				self.member[masteridx].Position = UNION_POSITION_MEMBER
			}
			self.Masteruid = self.member[idx].Uid
			self.Mastername = self.member[idx].Uname
			self.member[idx].Position = UNION_POSITION_MASTER

			//! 增加record
			if hasmaster {
				self.record = append(self.record, JS_UnionRecord{
					Type:  UNION_RECORD_TYPE_UNION_CHANGE_MASTER,
					Time:  TimeServer().Unix(),
					Name:  self.member[idx].Uname,
					Param: self.member[masteridx].Uname})
			} else {
				self.record = append(self.record, JS_UnionRecord{
					Type:  UNION_RECORD_TYPE_UNION_CHANGE_MASTER,
					Time:  TimeServer().Unix(),
					Name:  self.member[idx].Uname,
					Param: ""})
			}

			self.ChangeMaster.CheckFlag = true
			self.ChangeMaster.NowMaster = self.Mastername
			if hasmaster {
				self.ChangeMaster.OldMaster = self.member[masteridx].Uname
			}
			self.ChangeMaster.CheckPlayer = make([]int64, 0)

			//newmaster := GetPlayerMgr().GetPlayer(self.Masteruid, true)
			//if newmaster != nil {
			//	newmaster.GetModule("union").(*ModUnion).Sql_UserUnionInfo.Position = UNION_POSITION_MASTER
			//}
			//
			//if hasmaster {
			//	oldmaster := GetPlayerMgr().GetPlayer(self.member[masteridx].Uid, true)
			//	if oldmaster != nil {
			//		oldmaster.GetModule("union").(*ModUnion).Sql_UserUnionInfo.Position = UNION_POSITION_MEMBER
			//	}
			//}

		}
	}

	//! 触发军团长缺失
	if hasmaster == false {

	}
}

//! 检查会长离线时间
func (self *UnionMgr) CheckMasterOffline(unionid int) {
	union := self.GetUnion(unionid)
	if union == nil {
		return
	}

	if TimeServer().Unix()-self.GetUnionCheckMaster(unionid) > 600 {
		self.SetUnionCheckMaster(unionid, TimeServer().Unix())

		var msg S2M_UnionCheckMaster
		msg.Unionuid = unionid

		ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_CHECK_MASTER, &msg)
		if ret == nil || ret.RetCode != UNION_SUCCESS {
			return
		}
		union.CheckMasterOffline()
	}
}

//创建军团
func (self *UnionMgr) CreateUnion(icon int, unionname string, uid int64, mastername string, camp int) int {
	union := new(San_Union)
	union.Icon = icon
	union.Unionname = unionname
	union.Masteruid = uid
	union.Mastername = mastername
	union.member = make([]JS_Member, 0)
	union.apply = make([]JS_UnionApply, 0)
	union.record = make([]JS_UnionRecord, 0)
	union.huntInfo = make([]*JS_UnionHunt, 0)
	union.braveHand = make([]*JS_UnionBraveHand, 0)
	union.Level = 1
	union.Locker = new(sync.RWMutex)
	union.Joinlevel = UNION_JOIN_LEVEL_BASE
	union.Exp = 0
	union.ChangeMaster.CheckPlayer = make([]int64, 0)
	for i := UNION_HUNT_TYPE_NOMAL; i < UNION_HUNT_TYPE_MAX; i++ {

		config := GetCsvMgr().GetUnionHuntConfigByID(i)
		if nil == config {
			continue
		}
		temp := JS_UnionHunt{}
		temp.Type = i
		temp.TopDps = lstUnionHuntDpsTop{}
		if i == UNION_HUNT_TYPE_NOMAL {
			temp.EndTime = union.GetRefreshTime()
		}
		union.huntInfo = append(union.huntInfo, &temp)
	}
	var msg S2M_UnionCreateUnion
	msg.Icon = icon
	msg.Mastername = mastername
	msg.Masteruid = uid
	msg.Unionname = unionname
	msg.Joinlevel = UNION_JOIN_LEVEL_BASE
	msg.HuntInfo = union.huntInfo

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_CREATE, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return -1
	}

	union.Encode()
	var backmsg M2S_UnionCreateUnion
	json.Unmarshal([]byte(ret.Data), &backmsg)
	union.Id = backmsg.Unionuid

	self.Locker.Lock()
	self.Sql_Union[union.Id] = union
	self.Locker.Unlock()

	self.NameLocker.Lock()
	self.Sql_UnionName[union.Id] = unionname
	self.NameLocker.Unlock()

	return union.Id
}

//解散军团 ret:1只有会长才能操作 2:军团人员大于1
func (self *UnionMgr) DissolveUnion(unionid int, uid int64) int {
	union := self.GetUnion(unionid)
	if union == nil {
		return -1
	}

	if union.Masteruid != uid {
		return 1
	}

	if len(union.member) > 1 {
		return 2
	}

	var msg S2M_UnionDissolve
	msg.Uid = uid
	msg.Unionuid = unionid

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_DISSOLVE, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return 1
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()
	//DeleteTable("san_unioninfo", union, []int{0})

	delete(self.Sql_Union, unionid)

	return 0
}

//! 统计军团战力
func (self *San_Union) CalcFight() {
	self.Fight = 0
	for i := 0; i < len(self.member); i++ {
		self.Fight += self.member[i].Fight
	}
}

//! 检查是否有军团转移
func (self *San_Union) CheckChangeMasterState() {
	tNowTime := TimeServer().Unix()
	for i := 0; i < len(self.record); i++ {
		if self.record[i].Type == UNION_RECORD_TYPE_UNION_CHANGE_MASTER && tNowTime < self.record[i].Time+DAY_SECS {
			self.ChangeMaster.CheckFlag = true
			self.ChangeMaster.NowMaster = self.Mastername
			self.ChangeMaster.OldMaster = self.record[i].Param
			break
		}
	}
}

func (self *San_Union) IsChangeMaster(uid int64) bool {
	if self.ChangeMaster.CheckFlag == false {
		return false
	}

	for i := 0; i < len(self.ChangeMaster.CheckPlayer); i++ {
		if self.ChangeMaster.CheckPlayer[i] == uid {
			return false
		}
	}

	return true
}

//! 加经验
func (self *San_Union) AddExp(exp int, addtype int) (int, int) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	var msg S2M_UnionAddUnionExp
	msg.Unionuid = self.Id
	msg.Count = exp
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_ADD_EXP, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return 0, 0
	}

	csv, _ := GetCsvMgr().CommunityConfig[self.Level]
	addexp := exp
	if addtype != 2 {
		if self.DayExp+exp > csv.GuildExpLimit {
			addexp = csv.GuildExpLimit - self.DayExp
		}
		self.Exp += addexp
		self.DayExp += addexp
	} else {
		self.Exp += addexp
	}

	need := csv.Exp
	for self.Exp >= need {
		nextcsv, ok := GetCsvMgr().CommunityConfig[self.Level+1]
		if !ok {
			self.Exp = need
			break
		}

		self.Level++
		self.Exp -= need

		need = nextcsv.Exp
	}

	return self.Level, self.Exp
}

func (self *San_Union) UpdateMember(userinfo *JS_Member, position int, todaydon int) bool {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == userinfo.Uid {
			self.member[i].Uid = userinfo.Uid
			self.member[i].Level = userinfo.Level
			self.member[i].Uname = userinfo.Uname
			self.member[i].Iconid = userinfo.Iconid
			self.member[i].Portrait = userinfo.Portrait
			self.member[i].Vip = userinfo.Vip
			self.member[i].Position = position
			self.member[i].Donation = 0
			self.member[i].TodayDon += todaydon
			self.member[i].CopyNum = 0
			self.member[i].Fight = userinfo.Fight
			return true
		}
	}

	return false
}

// 增加狩猎钻石奖励公会记录
func (self *San_Union) AddGemAwardRecord(name string, position, nType, count int) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.record = append(self.record, JS_UnionRecord{UNION_RECORD_TYPE_HUNTER_BOSS_AWARD,
		TimeServer().Unix(),
		name,
		fmt.Sprintf("%d", nType*100000+position*1000+count)})

	return false
}

func (self *San_Union) FreshMember(userinfo *San_UserBase, userunioninfo *San_UserUnionInfo) bool {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == userinfo.Uid {
			self.member[i].Iconid = userinfo.IconId
			self.member[i].Portrait = userinfo.Portrait
			return true
		}
	}

	return false
}

func (self *San_Union) AddMember(userinfo *JS_Member, position int, iscreate bool) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	var data JS_Member
	data.Uid = userinfo.Uid
	data.Level = userinfo.Level
	data.Uname = userinfo.Uname
	data.Iconid = userinfo.Iconid
	data.Portrait = userinfo.Portrait
	data.Vip = userinfo.Vip
	data.Camp = userinfo.Camp
	if position > 0 {
		data.Position = position
	}
	data.Donation = 0
	data.TodayDon = 0
	data.CopyNum = 0
	data.Fight = userinfo.Fight

	self.Fight += userinfo.Fight

	data.Stage = userinfo.Stage

	self.member = append(self.member, data)

	record_type := UNION_RECORD_TYPE_JOIN

	if iscreate {
		record_type = UNION_RECORD_TYPE_CREATE
	}

	self.record = append(self.record, JS_UnionRecord{record_type, TimeServer().Unix(), userinfo.Uname, "0"})

	//! 排行榜更新
	//GetTopMgr().SyncUnionFight(self.Fight, self)

	return true
}

// 更新成员 没有就加进去 最后的是否创建是指是否是创建公会只是为了记录日志使用
func (self *UnionMgr) UpdateMember(unionid int, uid int64, position int, iscreate bool, isadd bool) bool {
	value := self.GetUnion(unionid)
	if value == nil {
		return false
	}

	var msg S2M_UnionUpdateMember
	msg.Unionid = unionid
	msg.Uid = uid
	msg.Position = position
	msg.IsAdd = isadd
	msg.IsCreate = iscreate

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_UPDATE_MEMBER, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_UnionUpdateMember
	json.Unmarshal([]byte(ret.Data), &backmsg)

	var member JS_Member
	member.Uid = backmsg.Uid
	member.Level = backmsg.Level
	member.Uname = backmsg.UName
	member.Iconid = backmsg.IconId
	member.Portrait = backmsg.Portrait
	member.Vip = backmsg.Vip
	member.Position = backmsg.Position
	member.Stage = backmsg.Stage
	member.Fight = backmsg.Fight

	if !isadd {
		value.UpdateMember(&member, backmsg.Position, 0)
	} else {
		value.AddMember(&member, position, iscreate)
	}
	return true
}

func (self *UnionMgr) FreshMember(unionid int, userinfo *San_UserBase, userunioninfo *San_UserUnionInfo) bool {
	//self.Locker.RLock()
	//value, ok := self.Sql_Union[unionid]
	//self.Locker.RUnlock()
	//
	//if !ok {
	//	return false
	//}
	//
	//value.FreshMember(userinfo, userunioninfo)
	self.UpdateMember(unionid, userinfo.Uid, 0, false, false)

	return true
}

//func (self *San_Union) UpdateMemberTimeAndName(userinfo San_UserBase, userunioninfo San_UserUnionInfo) bool {
//	self.Locker.Lock()
//	defer self.Locker.Unlock()
//
//	for i := 0; i < len(self.member); i++ {
//		if self.member[i].Uid == userinfo.Uid {
//
//			//_t, _ := time.ParseInLocation(DATEFORMAT, userinfo.LastLoginTime, time.Local)
//			//self.member[i].Lastlogintime = _t.Unix()
//			self.member[i].Uname = userinfo.UName
//			self.member[i].Level = userinfo.Level
//			self.member[i].Fight = userinfo.Fight
//			self.member[i].Vip = userinfo.Vip
//			break
//		}
//	}
//
//	self.CalcFight()
//
//	if userinfo.Uid == self.Masteruid {
//		self.Mastername = userinfo.UName
//	}
//	return true
//}
//
//func (self *UnionMgr) UpdateMemberTimeAndName(userinfo San_UserBase, userunioninfo San_UserUnionInfo) bool {
//
//	if userunioninfo.Unionid == 0 {
//		return false
//	}
//
//	self.Locker.RLock()
//	value, ok := self.Sql_Union[userunioninfo.Unionid]
//	self.Locker.RUnlock()
//
//	if !ok {
//		return false
//	}
//
//	value.UpdateMemberTimeAndName(userinfo, userunioninfo)
//
//	//log.Println(self.Sql_Union[unionid])
//	return true
//}

func (self *San_Union) UpdateMemberState(uid int64, fight int64, vip int, online bool) bool {
	//self.Locker.Lock()
	//defer self.Locker.Unlock()

	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == uid {
			if self.member[i].Fight != fight {
				self.member[i].Fight = fight
				self.member[i].Vip = vip

				self.CalcFight()
			}

			if online {
				self.member[i].Lastlogintime = 0
			} else {
				self.member[i].Lastlogintime = TimeServer().Unix()
			}

			break
		}
	}

	return true
}

func (self *UnionMgr) UpdateMemberState(unionid int, uid int64) bool {

	if unionid == 0 {
		return false
	}
	value := self.GetUnion(unionid)
	if nil == value {
		return false
	}

	var msg S2M_UnionUpdateMemberState
	msg.Uid = uid
	msg.Unionid = unionid
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_UPDATE_MEMBER_STATE, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_UnionUpdateMemberState
	json.Unmarshal([]byte(ret.Data), &backmsg)

	value.UpdateMemberState(uid, backmsg.Fight, backmsg.Vip, backmsg.Online)

	//log.Println(self.Sql_Union[unionid])
	return true
}

func (self *UnionMgr) UpdateMemberPassID(unionid int, uid int64, passid int) bool {

	if unionid == 0 {
		return false
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()

	value, ok := self.Sql_Union[unionid]
	if !ok {
		return false
	}

	for i := 0; i < len(value.member); i++ {
		if value.member[i].Uid == uid {
			if value.member[i].Stage != passid {
				value.member[i].Stage = passid
			}
			break
		}
	}

	//log.Println(self.Sql_Union[unionid])
	return true
}

type topUnion []*San_Union

func (s topUnion) Len() int      { return len(s) }
func (s topUnion) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s topUnion) Less(i, j int) bool {
	if s[i].Level > s[j].Level {
		return true
	}

	return false
}

func (self *UnionMgr) GetUnionList(serverid int) []JS_Union2 {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	var lst []JS_Union2

	var msg S2M_UnionGetUnionList
	msg.ServerID = serverid
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GET_UNION_LIST, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return lst
	}

	var backmsg M2S_UnionGetUnionList
	json.Unmarshal([]byte(ret.Data), &backmsg)
	lst = backmsg.UnionList

	//for _, value := range self.Sql_Union {
	//	var data JS_Union2
	//	data.Id = value.Id
	//	data.Icon = value.Icon
	//	data.Unionname = value.Unionname
	//	data.Masteruid = value.Masteruid
	//	data.Mastername = value.Mastername
	//	data.Level = value.Level
	//	data.Jointype = value.Jointype
	//	data.Joinlevel = value.Joinlevel
	//	data.Member = len(value.member)
	//	data.Camp = value.GetCamp()
	//	data.Fight = value.Fight
	//	data.Exp = value.Exp
	//	data.ActivityPoint = value.ActivityPoint
	//
	//	lst = append(lst, data)
	//}

	return lst
}

func (self *UnionMgr) AddPlayer(unionid int, uid int64, send bool) (bool, int) {
	value := self.GetUnion(unionid)
	if value == nil {
		return false, 1
	}

	//player := GetPlayerMgr().GetPlayer(uid, false)
	//if nil != player {
	//	player.GetModule("union").(*ModUnion).ClearApplyUnion()
	//}
	if self.UpdateMember(unionid, uid, UNION_POSITION_MEMBER, false, true) {
		self.UpdateMemberState(unionid, uid)
	}

	//! 军团保存
	value.Encode()
	value.Update(true)

	if send {

	}
	return true, 0

	return false, -1
}

func (self *San_Union) OutPlayer(outuid int64, ismaster bool) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	for j := 0; j < len(self.member); j++ {
		if self.member[j].Uid == outuid {
			var record JS_UnionRecord
			if ismaster {
				record.Type = UNION_RECORD_TYPE_DISSOLVE
			} else {
				record.Type = UNION_RECORD_TYPE_OUT
			}

			record.Time = TimeServer().Unix()
			record.Name = self.member[j].Uname
			record.Param = "0"

			self.record = append(self.record, record)
			self.Fight -= self.member[j].Fight

			copy(self.member[j:], self.member[j+1:])
			self.member = self.member[:len(self.member)-1]

			//! 即时保存军团
			//self.Encode()
			//self.Update(true)

			return true
		}
	}

	return false
}

func (self *UnionMgr) OutPlayer(unionid int, outuid int64, ismaster bool) bool {
	value := self.GetUnion(unionid)
	if value == nil {
		return false
	}
	// 会长不能退出
	if value.Masteruid == outuid {
		return false
	}
	var msg S2M_UnionOutPlayer
	msg.Unionuid = unionid
	msg.OutUid = outuid
	msg.IsMaster = ismaster
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_OUT_PLAYER, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	value.OutPlayer(outuid, ismaster)
	value.Encode()

	return true
}

func (self *San_Union) AlertUnionName(newname string, icon int) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	if self.Unionname == newname && self.Icon == icon {
		return false
	}

	if self.Unionname != newname && self.Icon == icon {
		self.record = append(self.record, JS_UnionRecord{UNION_RECORD_TYPE_UNION_CHANGE_NAME, TimeServer().Unix(), self.Mastername, newname})
	} else if self.Unionname == newname && self.Icon != icon {
		self.record = append(self.record, JS_UnionRecord{UNION_RECORD_TYPE_UNION_CHANGE_ICON, TimeServer().Unix(), self.Mastername, fmt.Sprintf("%d", icon)})
	} else if self.Unionname != newname && self.Icon != icon {
		self.record = append(self.record, JS_UnionRecord{UNION_RECORD_TYPE_UNION_CHANGE_NAME, TimeServer().Unix(), self.Mastername, newname})
		self.record = append(self.record, JS_UnionRecord{UNION_RECORD_TYPE_UNION_CHANGE_ICON, TimeServer().Unix(), self.Mastername, fmt.Sprintf("%d", icon)})
	}

	self.Unionname = newname
	self.Icon = icon
	self.Lastupdtime = TimeServer().Unix()

	return true
}

//修改军团名字
func (self *UnionMgr) AlertUnionName(uid int64, unionid int, newname string, icon int) (bool, int) {
	value := self.GetUnion(unionid)
	if value == nil {
		return false, 0
	}

	if HF_IsLicitName([]byte(newname)) == false {
		return false, 0
	}

	var msg S2M_UnionAlertName
	msg.Icon = icon
	msg.Unionuid = unionid
	msg.Uid = uid
	msg.Name = newname

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_ALERT_NAME, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		GetPlayerMgr().GetPlayer(uid, false).SendErrInfo("err", fmt.Sprintf("中心服返回错误%d", ret.RetCode))
		return false, 0
	}

	var backmsg M2S_UnionAlertName
	json.Unmarshal([]byte(ret.Data), &backmsg)

	self.NameLocker.Lock()
	self.Sql_UnionName[unionid] = newname
	self.NameLocker.Unlock()

	return value.AlertUnionName(newname, icon), backmsg.Ret

}

func (self *UnionMgr) GetUnionRecord(unionid int) []JS_UnionRecord {
	value := self.GetUnion(unionid)
	if value != nil {
		return value.record
	} else {
		return make([]JS_UnionRecord, 0)
	}

}

func (self *San_Union) AlertUnionNotice(newname string) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.Notice = newname
	self.Lastupdtime = TimeServer().Unix()

	return true
}

func (self *San_Union) AlertUnionBoard(newname string) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.Board = newname
	self.Lastupdtime = TimeServer().Unix()

	return true
}

//修改信息
func (self *UnionMgr) AlertUnionNotice(uid int64, unionid int, notice string) bool {
	value := self.GetUnion(unionid)
	if nil == value {
		return false
	}

	var msg S2M_UnionAlertNotice
	msg.Unionuid = unionid
	msg.Uid = uid
	msg.Content = notice

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_ALERT_NOTICE, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	return value.AlertUnionNotice(notice)
}

//修改信息
func (self *UnionMgr) AlertUnionBoard(uid int64, unionid int, board string) bool {
	value := self.GetUnion(unionid)
	if nil == value {
		return false
	}

	var msg S2M_UnionAlertBoard
	msg.Unionuid = unionid
	msg.Uid = uid
	msg.Content = board

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_ALERT_BOARD, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}
	return value.AlertUnionBoard(board)
}

func (self *San_Union) AlertUnionSet(jointype int, joinlevel int) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.Jointype = jointype
	self.Joinlevel = joinlevel
	self.Lastupdtime = TimeServer().Unix()

	return true
}

//修改军团设置
func (self *UnionMgr) AlertUnionSet(uid int64, unionid int, jointype int, joinlevel int) bool {
	value := self.GetUnion(unionid)
	if nil == value {
		return false
	}

	var msg S2M_UnionAlertSet
	msg.Unionuid = unionid
	msg.Uid = uid
	msg.Type = jointype
	msg.Level = joinlevel

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_ALERT_SET, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	value.AlertUnionSet(jointype, joinlevel)

	return true
}

func (self *San_Union) UnionModify(destuid int64, op int) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == destuid {
			self.member[i].Position = op
			self.record = append(self.record, JS_UnionRecord{UNION_RECORD_TYPE_MODIFY, TimeServer().Unix(), self.member[i].Uname, fmt.Sprintf("%d", op)})
			return true
		}
	}

	return false
}

//军团任命
func (self *UnionMgr) UnionModify(uid int64, unionid int, destuid int64, op int) bool {
	self.Locker.RLock()
	value, ok := self.Sql_Union[unionid]
	self.Locker.RUnlock()

	if !ok {
		return false
	}

	var msg S2M_UnionUnionModify
	msg.Unionuid = unionid
	msg.Uid = uid
	msg.Destuid = destuid
	msg.Op = op

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_MODIFY, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_UnionUnionModify
	json.Unmarshal([]byte(ret.Data), &backmsg)

	//GetPlayerMgr().GetPlayer(destuid, true).GetModule("union").(*ModUnion).ModifyPosition(op)

	return value.UnionModify(destuid, op)
}
func (self *UnionMgr) CheckBraveHand(data *San_Union) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	now := TimeServer().Unix()
	nLen := len(data.braveHand)
	for i := nLen - 1; i >= 0; i-- {
		if data.braveHand[i].EndTime > 0 && now >= data.braveHand[i].EndTime {
			data.braveHand = append(data.braveHand[0:i], data.braveHand[i+1:]...)
		}
	}
}

//军团任命
func (self *UnionMgr) SetBraveHand(uid int64, unionid int, destuid int64, op int) bool {
	value := self.GetUnion(unionid)
	if nil == value {
		return false
	}

	csv, _ := GetCsvMgr().CommunityConfig[value.Level]
	if csv == nil {
		return false
	}

	if op == 1 && len(value.braveHand) >= csv.Fearless {
		return false
	}

	var msg S2M_UnionSetBraveHand
	msg.Uid = uid
	msg.Unionuid = unionid
	msg.Destuid = destuid
	msg.Op = op
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_SET_BRAVE_HAND, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_UnionSetBraveHand
	json.Unmarshal([]byte(ret.Data), &backmsg)

	for i := 0; i < len(value.member); i++ {
		if value.member[i].Uid == destuid {
			if op == 1 {
				value.braveHand = append(value.braveHand, &JS_UnionBraveHand{destuid, 0})
			} else {
				find := false
				for _, v := range value.braveHand {
					if v.Uid == destuid {
						v.Uid = 0
						v.EndTime = TimeServer().Unix() + BRAVE_HAND_CD
						find = true
						break
					}
				}

				if !find {
					return false
				}
			}

			value.member[i].BraveHand = op
			return true
		}
	}

	return false
}

func (self *San_Union) UnionChange(uid int64, destuid int64, name string) bool {

	self.Locker.Lock()
	defer self.Locker.Unlock()

	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == destuid {
			self.Masteruid = destuid
			self.member[i].Position = UNION_POSITION_MASTER
		} else if self.member[i].Uid == uid {
			self.member[i].Position = UNION_POSITION_MEMBER
		}
	}

	return true
}

//军团转让
func (self *UnionMgr) UnionChange(unionid int, uid int64, destuid int64) bool {
	value := self.GetUnion(unionid)
	if nil == value {
		return false
	}

	var msg S2M_UnionChange
	msg.Unionuid = unionid
	msg.Uid = uid
	msg.Destuid = destuid

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_CHANGE, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_UnionChange
	json.Unmarshal([]byte(ret.Data), &backmsg)

	return value.UnionChange(uid, destuid, backmsg.Name)
}

////! 刷新每日贡献
//func (self *UnionMgr) OnRefresh() {
//	self.Locker.RLock()
//	defer self.Locker.RUnlock()
//
//	for _, value := range self.Sql_Union {
//		if value.ServerID != GetServer().Con.ServerId {
//			continue
//		}
//		//value.OnRefresh()
//		value.RefreshActivityLimit()
//	}
//}

//获取军团玩家的军团等级 1-4 会长 长老 精英 成员
func (self *San_Union) OnRefresh() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.DayExp = 0
	for i := 0; i < len(self.member); i++ {
		self.member[i].TodayDon = 0
		self.member[i].CopyNum = 0
	}
}

//获取军团玩家的军团等级 1-4 会长 长老 精英 成员
func (self *San_Union) GetPlayerUnionLevel(uid int64) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == uid {
			return self.member[i].Position
		}
	}

	return 0
}

//! 军团广播消息
func (self *San_Union) BroadCastMsg(head string, body []byte) {

	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for i := 0; i < len(self.member); i++ {
		value := GetPlayerMgr().GetPlayer(self.member[i].Uid, false)
		if value != nil {
			value.SendMsg(head, body)
		}
	}
}

//! 刷新副本
func (self *San_Union) RefreshActivityLimit() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var msg S2M_UnionRefreshActivityLimit
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_ACTIVITY_REFRESH, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}

	now := TimeServer()
	timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()

	for _, v := range self.member {
		nSize := len(v.ActivityRecord)
		for i := nSize - 1; i >= 0; i-- {
			// 清理过期的限制表
			if (v.ActivityRecord[i].Time-timeStamp)/DAY_SECS > UNION_ACTIVITY_LIMIT_CLEAN {
				v.ActivityRecord = append(v.ActivityRecord[:i],
					v.ActivityRecord[i+1:]...)
			}
		}
	}
	self.AcitvityLimit = 0
}

//! 获得等级人数
func (self *San_Union) GetMemberCount(position int) int {
	count := 0
	for i := 0; i < len(self.member); i++ {
		if self.member[i].Position == position {
			count++
		}
	}

	return count
}

//! 副本排行
func (self *UnionMgr) Rename(uid int64, newname string) {

}

// 先获得
func (self *UnionMgr) GetUnionNameMap() map[int64]string {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	// 全服发送军团红包
	res := make(map[int64]string)
	for _, value := range self.Sql_Union {
		for _, member := range value.member {
			res[member.Uid] = value.Unionname
		}
	}
	return res
}

//! 加活跃度
func (self *San_Union) AddActivityPoint(uid int64, count int) int {
	config, ok := GetCsvMgr().CommunityConfig[self.Level]
	if !ok || config == nil {
		return 0
	}

	var msg S2M_UnionAddUnionActivity
	msg.Uid = uid
	msg.Unionuid = self.Id
	msg.Count = count
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_ADD_ACTIVITY, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return 0
	}

	var backmsg M2S_UnionAddUnionActivity
	json.Unmarshal([]byte(ret.Data), &backmsg)

	self.Locker.Lock()
	defer self.Locker.Unlock()

	addCount := backmsg.Count

	index := -1
	for i, v := range self.member {
		if v.Uid == uid {
			index = i
			break
		}
	}

	if index < 0 {
		return 0
	}

	now := TimeServer()
	timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
	if TimeServer().Hour() < 5 {
		timeStamp -= DAY_SECS
	}

	var record *UserActivityRecord = nil

	for _, v := range self.member[index].ActivityRecord {
		if v.Time == timeStamp {
			record = v
			break
		}
	}

	if nil == record {
		record = &UserActivityRecord{}
		record.Time = timeStamp
		self.member[index].ActivityRecord = append(self.member[index].ActivityRecord, record)
	}

	// 玩家每日提供上限检测
	if record.AddCount+addCount > config.Activelimit {
		addCount = config.Activelimit - record.AddCount
	}

	// 工会每日接受上限检测
	if self.AcitvityLimit+addCount >= config.GuildActiveLimit {
		addCount = config.GuildActiveLimit - self.AcitvityLimit
	}

	// 工会总上限检测
	if self.ActivityPoint+addCount >= config.Lively {
		addCount = config.Lively - self.ActivityPoint
	}

	if addCount <= 0 {
		return 0
	}

	self.ActivityPoint += addCount
	self.AcitvityLimit += addCount
	record.AddCount += addCount

	return addCount
}

func (self *San_Union) MinActivityPoint(count int) int {
	if self.ActivityPoint <= 0 {
		return 0
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()

	addCount := count

	if self.ActivityPoint < addCount {
		addCount = self.ActivityPoint
	}

	self.ActivityPoint -= addCount

	return addCount
}

func (self *San_Union) OpenHunterFight(nType int, endTime int64) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	info := &JS_UnionHunt{}
	for _, v := range self.huntInfo {
		if nType == v.Type {
			info = v
		}
	}

	if nil == info {
		return false
	}

	info.EndTime = endTime
	info.TopDps = lstUnionHuntDpsTop{}

	return true
}

//! 增加玩家战报
func (self *San_Union) AddUnionHuntDamage(nType int, nCount int64, player *Player, nJob int, battleInfo *BattleInfo, battleRecord *BattleRecord) int {
	config, ok := GetCsvMgr().CommunityConfig[self.Level]
	if !ok || config == nil {
		return 0
	}

	var mastermsg S2M_UnionAddDamage
	mastermsg.Uid = player.GetUid()
	mastermsg.Unionuid = self.Id
	mastermsg.Dps = nCount
	mastermsg.Type = nType
	mastermsg.FightID = battleInfo.Id
	mastermsg.Record = HF_JtoA(battleRecord)
	mastermsg.Info = HF_JtoA(battleInfo)
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_ADD_DAMAGE, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return 0
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()

	info := &JS_UnionHunt{}
	for _, v := range self.huntInfo {
		if nType == v.Type {
			info = v
		}
	}

	if nil == info {
		return 0
	}

	info.TopDps = append(info.TopDps,
		&JS_UnionHuntDpsTop{player.Sql_UserBase.Uid,
			player.Sql_UserBase.UName,
			player.Sql_UserBase.IconId,
			player.Sql_UserBase.Portrait,
			player.Sql_UserBase.Level,
			nJob,
			player.Sql_UserBase.Fight,
			player.Sql_UserBase.Vip,
			nCount,
			battleInfo.Id,
			TimeServer().Unix()})

	sort.Sort(lstUnionHuntDpsTop(info.TopDps))

	//HMSetRedis("san_huntbattleinfo", battleInfo.Id, battleInfo, DAY_SECS*10)
	//HMSetRedis("san_huntbattlerecord", battleRecord.Id, battleRecord, DAY_SECS*10)
	return 0
}

// 开启战斗协程
func (self *San_Union) CheckHunterFightEnd() {
	var msg S2M_UnionCheckHunterEnd
	msg.Unionuid = self.Id
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_CHECK_HUNTER_END, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}

	var backmsg M2S_UnionUpdateMember
	json.Unmarshal([]byte(ret.Data), &backmsg)

	for _, v := range self.huntInfo {
		if UNION_HUNT_TYPE_NOMAL == v.Type {
			if v.EndTime == 0 {
				v.EndTime = self.GetRefreshTime()
			}
			continue
		} else {
			if v.EndTime > 0 && TimeServer().Unix() >= v.EndTime {
				v.EndTime = 0
				self.record = append(self.record, JS_UnionRecord{UNION_RECORD_TYPE_HUNTER_BOSS_LEAVE, TimeServer().Unix(), self.Mastername, fmt.Sprintf("%d", v.Type)})
			}
		}
	}
}

// 开启战斗协程
func (self *UnionMgr) OnHunterRefresh() {
	for _, data := range self.Sql_Union {
		if data.ServerID != GetServer().Con.ServerId {
			continue
		}
		endtime := data.GetRefreshTime()

		var msg S2M_UnionHunterRefresh
		ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_HUNTER_REFRESH, msg)
		if ret == nil || ret.RetCode != UNION_SUCCESS {
			continue
		}

		for _, v := range data.huntInfo {
			if UNION_HUNT_TYPE_NOMAL == v.Type {
				if v.EndTime == endtime {
					continue
				} else {
					v.EndTime = endtime
					data.record = append(data.record, JS_UnionRecord{UNION_RECORD_TYPE_HUNTER_BOSS_LEAVE, TimeServer().Unix(), data.Mastername, fmt.Sprintf("%d", v.Type)})
				}

				v.TopDps = lstUnionHuntDpsTop{}
			}
		}
	}
}

// 开启战斗协程
func (self *San_Union) GetRefreshTime() int64 {
	if TimeServer().Hour() < 5 {
		return time.Date(TimeServer().Year(), TimeServer().Month(), TimeServer().Day(), 5, 0, 0, 0, TimeServer().Location()).Unix()
	} else {
		return time.Date(TimeServer().Year(), TimeServer().Month(), TimeServer().Day(), 5, 0, 0, 0, TimeServer().Location()).Unix() + DAY_SECS
	}
}

// 开启战斗协程
func (self *UnionMgr) Run() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			for _, data := range self.Sql_Union {
				if data.ServerID != GetServer().Con.ServerId {
					continue
				}
				data.CheckHunterFightEnd()
			}
		}
	}
	ticker.Stop()
}

//发邮件
func (self *UnionMgr) SendMail(unionid int, uid int64, title string, text string) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	value, ok := self.Sql_Union[unionid]
	if !ok {
		return false
	}

	var member *JS_Member = nil
	for _, v := range value.member {
		if v.Uid == uid {
			member = &v
			break
		}
	}

	if nil == member {
		return false
	}

	if member.Position > UNION_POSITION_VICE_MASTER {
		return false
	}

	var msg S2M_UnionSendMail
	msg.Unionuid = unionid
	msg.Uid = uid
	msg.Text = text
	msg.Title = title

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_SEND_MAIL, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	return true
}

//修改信息
func (self *UnionMgr) GMAlertUnionNotice(unionid int, notice string) bool {
	value := self.GetUnion(unionid)
	if nil == value {
		return false
	}

	var msg S2M_GMUnionAlertNotice
	msg.Unionuid = unionid
	msg.Content = notice

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GM_ALERT_NOTICE, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	return value.AlertUnionNotice(notice)
}

//修改信息
func (self *UnionMgr) GMAlertUnionBoard(unionid int, board string) bool {
	value := self.GetUnion(unionid)
	if nil == value {
		return false
	}

	var msg S2M_GMUnionAlertBoard
	msg.Unionuid = unionid
	msg.Content = board

	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GM_ALERT_BOARD, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}
	return value.AlertUnionBoard(board)
}
