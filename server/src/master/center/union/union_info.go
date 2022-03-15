/*
@Time : 2020/4/22 9:55
@Author : 96121
@File : union
@Software: GoLand
*/
package union

import (
	"encoding/json"
	"fmt"
	"game"
	"master/center/chat"
	"master/center/match"
	"master/center/tower"
	"master/core"
	"master/db"
	"master/utils"
	"sort"
	"sync"
	"time"
)

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
	UNION_ALERT_TYPE_NAME = 1
	UNION_ALERT_TYPE_ICON = 2
	UNION_ALERT_TYPE_BOTH = 3
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
const QUIT_CD = 3600
const JOIN_CD = 3600
const BRAVE_HAND_CD = 86400
const UNION_ACTIVITY_LIMIT_CLEAN = 8

type MSG_UnionInfo struct {
	Id            int                  //! 公会ID
	Icon          int                  //! icon
	Unionname     string               //! 公会名
	Masteruid     int64                //! 所有者ID
	Mastername    string               //! 会长昵称
	Level         int                  //! 公会等级
	Jointype      int                  //! 加入类型
	Joinlevel     int                  //! 加入等级
	ServerID      int                  //! 服务器id
	Notice        string               //! 公告
	Board         string               //! 对外展示
	Createtime    int64                //! 创建时间
	Lastupdtime   int64                //! 更新时间
	Fight         int64                //! 总战力
	Exp           int                  //! 经验
	DayExp        int                  //! 每日经验
	ActivityPoint int                  //! 活跃点数
	AcitvityLimit int                  //! 活跃度限额
	MailCD        int64                //! 邮件cd
	Member        string               //! 成员列表
	Applys        string               //! 申请列表
	Record        string               //! 操作记录
	HuntInfo      string               //! 军团狩猎记录
	BraveHand     string               //! 无畏之手
	ChangeMaster  JS_UnionChangeMaster //! 军团长自动更换
}
type UnionMail struct {
	Title string `json:"titile"`
	Text  string `json:"text"`
}

type UnionMailSave struct {
	Title string `json:"titile"`
	Text  string `json:"text"`
	Time  int64  `json:"time"`
}

//! 公会信息
type UnionInfo struct {
	Id            int    //! 公会ID
	Icon          int    //! icon
	Unionname     string //! 公会名
	Masteruid     int64  //! 所有者ID
	Mastername    string //! 会长昵称
	Level         int    //! 公会等级
	Jointype      int    //! 加入类型
	Joinlevel     int    //! 加入等级
	ServerID      int    //! 服务器id
	Notice        string //! 公告
	Board         string //! 对外展示
	Createtime    int64  //! 创建时间
	Lastupdtime   int64  //! 更新时间
	Fight         int64  //! 总战力
	Exp           int    //! 经验
	DayExp        int    //! 每日经验
	ActivityPoint int    //! 活跃点数
	AcitvityLimit int    //! 活跃度限额
	MailCD        int64  //! 邮件cd
	Member        string //! 成员列表
	Applys        string //! 申请列表
	Record        string //! 操作记录
	HuntInfo      string //! 军团狩猎记录
	BraveHand     string //! 无畏之手
	LastMail      string //! 最后一封邮件

	member    []*JS_UnionMember    //! 成员列表
	apply     []*JS_UnionApply     //! 申请列表
	record    []*JS_UnionRecord    //! 操作记录
	huntInfo  []*JS_UnionHunt      //! 军团狩猎记录
	braveHand []*JS_UnionBraveHand //! 无畏之手

	Locker       *sync.RWMutex        //! 操作锁
	ChangeMaster JS_UnionChangeMaster //! 军团长自动更换
	db.DataUpdate
}

type JS_UnionChangeMaster struct {
	CheckPlayer []int64 //! 同步过的玩家
	CheckFlag   bool    //! 同步标志
	OldMaster   string  //! 老军团长
	NowMaster   string  //!
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

//! 军团记录
type JS_UnionRecord struct {
	Type  int    `json:"type"`  //! 0-加入 1-退出 2-改变职务 3-军团更名 6-更名 7-更改标志 8-军团长更替
	Time  int64  `json:"time"`  //! 操作时间
	Name  string `json:"name"`  //! 昵称
	Param string `json:"param"` //! 参数
}

type JS_UnionBraveHand struct {
	Uid     int64
	EndTime int64 `json:"endtime"` //! 结束时间
}

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

type JS_UnionHunt struct {
	Type    int                `json:"type"`    //! 类型
	TopDps  lstUnionHuntDpsTop `json:"topdps"`  //! 伤害排行
	EndTime int64              `json:"endtime"` //! 结束时间
}

//! 保存数据
func (self *UnionInfo) onSave() {
	self.Encode()
	self.Update(true, false)
}

func (self *UnionInfo) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Member), &self.member)
	json.Unmarshal([]byte(self.Applys), &self.apply)
	json.Unmarshal([]byte(self.Record), &self.record)
	json.Unmarshal([]byte(self.HuntInfo), &self.huntInfo)
	json.Unmarshal([]byte(self.BraveHand), &self.braveHand)
}

func (self *UnionInfo) Encode() { //! 将data数据写入数据库
	self.Locker.RLock()
	self.Member = utils.HF_JtoA(&self.member)
	self.Applys = utils.HF_JtoA(&self.apply)
	self.Record = utils.HF_JtoA(&self.record)
	self.HuntInfo = utils.HF_JtoA(&self.huntInfo)
	self.BraveHand = utils.HF_JtoA(&self.braveHand)
	self.Locker.RUnlock()
}

func (self *UnionInfo) GetMember(uid int64) *JS_UnionMember {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, v := range self.member {
		if v.Uid == uid {
			return v
		}
	}
	return nil
}

//! 获得等级人数
func (self *UnionInfo) GetMemberCount(position int) int {
	count := 0
	for i := 0; i < len(self.member); i++ {
		if self.member[i].Position == position {
			count++
		}
	}

	return count
}

func (self *UnionInfo) GetMemberLen() int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	return len(self.member)
}
func (self *UnionInfo) GetBraveHandLen() int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	return len(self.braveHand)
}

func (self *UnionInfo) GetHunter(nType int) *JS_UnionHunt {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	for _, v := range self.huntInfo {
		if v.Type == nType {
			return v
		}
	}
	return nil
}

//! 统计军团战力
func (self *UnionInfo) CalcFight() {
	self.Fight = 0
	for i := 0; i < len(self.member); i++ {
		self.Fight += self.member[i].Fight
	}
}

func (self *UnionInfo) UpdateMemberState(uid int64, fight int64, vip int, online bool) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

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
				self.member[i].Lastlogintime = time.Now().Unix()
			}

			break
		}
	}

	return true
}

func (self *UnionInfo) AlertUnionName(newname string, icon int) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	if self.Unionname != newname && self.Icon == icon {
		self.record = append(self.record, &JS_UnionRecord{UNION_RECORD_TYPE_UNION_CHANGE_NAME, time.Now().Unix(), self.Mastername, newname})
	} else if self.Unionname == newname && self.Icon != icon {
		self.record = append(self.record, &JS_UnionRecord{UNION_RECORD_TYPE_UNION_CHANGE_ICON, time.Now().Unix(), self.Mastername, fmt.Sprintf("%d", icon)})
	} else if self.Unionname != newname && self.Icon != icon {
		self.record = append(self.record, &JS_UnionRecord{UNION_RECORD_TYPE_UNION_CHANGE_NAME, time.Now().Unix(), self.Mastername, newname})
		self.record = append(self.record, &JS_UnionRecord{UNION_RECORD_TYPE_UNION_CHANGE_ICON, time.Now().Unix(), self.Mastername, fmt.Sprintf("%d", icon)})
	}

	self.Unionname = newname
	self.Icon = icon
	self.Lastupdtime = time.Now().Unix()

	return true
}

func (self *UnionInfo) AlertUnionNotice(content string) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.Notice = content
	self.Lastupdtime = time.Now().Unix()

	return true
}

func (self *UnionInfo) AlertUnionBoard(content string) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.Board = content
	self.Lastupdtime = time.Now().Unix()

	return true
}

func (self *UnionInfo) AlertUnionSet(jointype int, joinlevel int) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.Jointype = jointype
	self.Joinlevel = joinlevel
	self.Lastupdtime = time.Now().Unix()

	return true

}

func (self *UnionInfo) AddApply(player *JS_UnionApply) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	for i := 0; i < len(self.apply); i++ {
		if self.apply[i].Uid == player.Uid {
			self.apply[i].Level = player.Level
			self.apply[i].Iconid = player.Iconid
			self.apply[i].Portrait = player.Portrait
			self.apply[i].Uname = player.Uname
			self.apply[i].Vip = player.Vip
			self.apply[i].Fight = player.Fight
			self.apply[i].Applytime = time.Now().Unix()
			self.apply[i].ServerID = player.ServerID
			return true
		}
	}
	self.apply = append(self.apply, &JS_UnionApply{player.Uid,
		player.Level,
		player.Uname,
		player.Iconid,
		player.Portrait,
		player.Vip,
		player.Fight,
		time.Now().Unix(),
		player.ServerID})

	return true

}

func (self *UnionInfo) CancelApply(uid int64) bool {
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

func (self *UnionInfo) IsMember(uid int64) bool {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, v := range self.member {
		if v.Uid == uid {
			return true
		}
	}
	return false
}

func (self *UnionInfo) GetApply() []*JS_UnionApply {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	return self.apply
}

func (self *UnionInfo) GetBraveHandUid() []int64 {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	rel := make([]int64, 0)
	for _, v := range self.member {
		if v.BraveHand == game.LOGIC_TRUE {
			rel = append(rel, v.Uid)
		}
	}
	return rel
}

func (self *UnionInfo) GetMemberList() []int64 {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	rel := make([]int64, 0)
	for _, v := range self.member {
		rel = append(rel, v.Uid)
	}
	return rel
}

func (self *UnionInfo) UpdateMember(player *JS_UnionMember, isadd bool, iscreate bool) int {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	position := -1
	index := -1
	for i, v := range self.member {
		if v.Uid == player.Uid {
			index = i
			break
		}
	}

	change := false
	if index < 0 {
		if isadd {
			self.member = append(self.member, player)
			position = player.Position
		}
	} else {
		self.member[index].Uid = player.Uid
		self.member[index].Level = player.Level
		if self.member[index].Uname != player.Uname {
			self.member[index].Uname = player.Uname
			change = true
		}

		if self.member[index].Iconid != player.Iconid {
			self.member[index].Iconid = player.Iconid
			change = true
		}
		if self.member[index].Portrait != player.Portrait {
			self.member[index].Portrait = player.Portrait
			change = true
		}
		self.member[index].Vip = player.Vip
		//self.member[index].Position = player.Position
		if player.Fight != 0 && self.member[index].Fight != player.Fight {
			self.member[index].Fight = player.Fight
		}

		if self.member[index].Stage != player.Stage {
			self.member[index].Stage = player.Stage
		}
		position = self.member[index].Position
	}

	if isadd {
		record_type := UNION_RECORD_TYPE_JOIN
		if iscreate {
			record_type = UNION_RECORD_TYPE_CREATE
		}
		self.record = append(self.record, &JS_UnionRecord{record_type, time.Now().Unix(), player.Uname, "0"})

		unionCh := chat.GetChatMgr().GetUnionChannel(self.Id)
		if unionCh != nil {
			unionCh.AddPlayer(player.Uid, player.Uname, player.Level, player.Iconid, player.Portrait, player.ServerID)
		}
	}

	if change {
		for _, v := range self.huntInfo {
			for _, dps := range v.TopDps {
				if dps.Uid == player.Uid {
					dps.Icon = player.Iconid
					dps.Portrait = player.Portrait
					dps.Name = player.Uname
				}
			}
		}
	}

	return position
}

//! 检查是否有军团转移
func (self *UnionInfo) CheckChangeMasterState() {
	tNowTime := time.Now().Unix()
	for i := 0; i < len(self.record); i++ {
		if self.record[i].Type == UNION_RECORD_TYPE_UNION_CHANGE_MASTER && tNowTime < self.record[i].Time+86400 {
			self.ChangeMaster.CheckFlag = true
			self.ChangeMaster.NowMaster = self.Mastername
			self.ChangeMaster.OldMaster = self.record[i].Param
			break
		}
	}
}

func (self *UnionInfo) CheckMasterOffline() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	change := false
	hasmaster := false
	masteridx := -1
	for j := 0; j < len(self.member); j++ {
		if self.member[j].Uid == self.Masteruid {
			hasmaster = true
			offtime := time.Now().Unix() - self.member[j].Lastlogintime
			if self.member[j].Lastlogintime > 0 && offtime > 86400*3 {
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
				self.record = append(self.record, &JS_UnionRecord{
					Type:  UNION_RECORD_TYPE_UNION_CHANGE_MASTER,
					Time:  time.Now().Unix(),
					Name:  self.member[idx].Uname,
					Param: self.member[masteridx].Uname})
			} else {
				self.record = append(self.record, &JS_UnionRecord{
					Type:  UNION_RECORD_TYPE_UNION_CHANGE_MASTER,
					Time:  time.Now().Unix(),
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

////会长拒绝加入
func (self *UnionInfo) MasterFail(applyuid int64) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	if applyuid == 0 {
		var apply []JS_UnionApply
		apply = make([]JS_UnionApply, 0)
		utils.HF_DeepCopy(&apply, &self.apply)
		for i := 0; i < len(apply); i++ {
			//self.CancelApply(apply[i].Uid)
			core.GetCenterApp().AddEvent(apply[i].ServerID, core.UNION_EVENT_MASTER_FAIL, apply[i].Uid,
				0, self.Id, "")
		}
	} else {
		index := -1
		for i := 0; i < len(self.apply); i++ {
			if self.apply[i].Uid == applyuid {
				index = i
				break
			}
		}

		if index < 0 {
			return
		}

		//self.CancelApply(applyuid)
		core.GetCenterApp().AddEvent(self.apply[index].ServerID, core.UNION_EVENT_MASTER_FAIL, self.apply[index].Uid,
			0, self.Id, "")

	}
}

func (self *UnionInfo) CleanApply(applyuid int64) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	if applyuid > 0 {
		for i := 0; i < len(self.apply); i++ {
			if self.apply[i].Uid == applyuid {
				self.apply = append(self.apply[:i], self.apply[i+1:]...)
				break
			}
		}
	} else {
		self.apply = []*JS_UnionApply{}
	}
}

func (self *UnionInfo) MasterOK(applyuid int64) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	index := -1
	for i := 0; i < len(self.apply); i++ {
		if self.apply[i].Uid == applyuid {
			index = i
			break
		}
	}

	if index < 0 {
		return
	}

	core.GetCenterApp().AddEvent(self.apply[index].ServerID, core.UNION_EVENT_MASTER_OK, self.apply[index].Uid,
		0, self.Id, "")
}

func (self *UnionInfo) KickPlayer(outid int64) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	index := -1
	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == outid {
			index = i
			break
		}
	}

	if index < 0 {
		return
	}

	//core.GetCenterApp().AddEvent(self.member[index].ServerID, core.UNION_EVENT_OUT_PLAYER, self.member[index].Uid,
	//0, self.Id, "")
}

func (self *UnionInfo) OutPlayer(outid int64, isMaster bool) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	index := -1
	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == outid {
			index = i
			break
		}
	}

	if index < 0 {
		return
	}

	ismaster := int64(0)
	if isMaster {
		ismaster = 1
	}

	core.GetCenterApp().AddEvent(self.member[index].ServerID, core.UNION_EVENT_OUT_PLAYER, self.member[index].Uid,
		ismaster, self.Id, "")

	var record JS_UnionRecord
	if isMaster {
		record.Type = UNION_RECORD_TYPE_DISSOLVE
	} else {
		record.Type = UNION_RECORD_TYPE_OUT
	}

	record.Time = time.Now().Unix()
	record.Name = self.member[index].Uname
	record.Param = "0"

	unionCh := chat.GetChatMgr().GetUnionChannel(self.Id)
	if unionCh != nil {
		unionCh.DelPlayer(outid)
	}

	self.record = append(self.record, &record)
	self.Fight -= self.member[index].Fight
	self.member = append(self.member[:index], self.member[index+1:]...)
}

//! 得到会长权限外最高的人
func (self *UnionInfo) GetIdWithoutMaster() int64 {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	lst1 := make([]int64, 0)
	lst2 := make([]int64, 0)

	for i := 0; i < len(self.member); i++ {
		if self.member[i].Position == UNION_POSITION_VICE_MASTER {
			return self.member[i].Uid
		} else if self.member[i].Position == UNION_POSITION_ELITE {
			lst1 = append(lst1, self.member[i].Uid)
		} else if self.member[i].Position == UNION_POSITION_MEMBER {
			lst2 = append(lst2, self.member[i].Uid)
		}
	}

	if len(lst1) > 0 {
		return lst1[utils.HF_GetRandom(len(lst1))]
	}

	if len(lst2) > 0 {
		return lst2[utils.HF_GetRandom(len(lst2))]
	}

	return 0
}
func (self *UnionInfo) UnionChange(uid int64, destuid int64, name string) bool {

	self.Locker.Lock()
	defer self.Locker.Unlock()

	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == destuid {
			self.Masteruid = destuid
			self.Mastername = name
			self.member[i].Position = UNION_POSITION_MASTER
		} else if self.member[i].Uid == uid {
			self.member[i].Position = UNION_POSITION_MEMBER
		}
	}

	return true
}

func (self *UnionInfo) UnionModify(destuid int64, op int) bool {
	csv, _ := GetUnionMgr().CommunityConfigs[self.Level]
	if csv == nil {
		return false
	}
	// 设置副会长
	if op == UNION_POSITION_VICE_MASTER {
		if self.GetMemberCount(UNION_POSITION_VICE_MASTER) >= csv.Elder {
			return false
		}
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()
	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == destuid {
			self.member[i].Position = op
			self.record = append(self.record, &JS_UnionRecord{UNION_RECORD_TYPE_MODIFY, time.Now().Unix(), self.member[i].Uname, fmt.Sprintf("%d", op)})
			return true
		}
	}

	return false
}
func (self *UnionInfo) SetBraveHand(destuid int64, op int) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	for i := 0; i < len(self.member); i++ {
		if self.member[i].Uid == destuid {
			if op == 1 {
				self.braveHand = append(self.braveHand, &JS_UnionBraveHand{destuid, 0})
			} else {
				find := false
				for _, v := range self.braveHand {
					if v.Uid == destuid {
						v.Uid = 0
						v.EndTime = time.Now().Unix() + BRAVE_HAND_CD
						find = true
						break
					}
				}

				if !find {
					return false
				}
			}

			self.member[i].BraveHand = op
			//value.record = append(value.record, JS_UnionRecord{2, time.Now().Unix(), value.member[i].Uname, fmt.Sprintf("%d", op)})
			return true
		}
	}

	return false
}

func (self *UnionInfo) MinActivityPoint(count int) int {
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

func (self *UnionInfo) OpenHunterFight(nType int, endTime int64) bool {
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
func (self *UnionInfo) AddUnionHuntDamage(player *JS_UnionMember, nType int, nCount int64, uid int64, fightid int64, battleInfo string, battleRecord string) int {
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
		&JS_UnionHuntDpsTop{uid,
			player.Uname,
			player.Iconid,
			player.Portrait,
			player.Level,
			player.Position,
			player.Fight,
			player.Vip,
			nCount,
			fightid,
			time.Now().Unix()})

	sort.Sort(lstUnionHuntDpsTop(info.TopDps))

	var js_battleinfo tower.BattleInfo
	var js_battleRecord tower.BattleRecord
	json.Unmarshal([]byte(battleInfo), &js_battleinfo)
	json.Unmarshal([]byte(battleRecord), &js_battleRecord)

	db.HMSetRedisEx("san_huntbattleinfo", fightid, js_battleinfo, utils.HOUR_SECS*12)
	db.HMSetRedisEx("san_huntbattlerecord", fightid, js_battleRecord, utils.HOUR_SECS*12)

	var db_battleInfo match.JS_CrossArenaBattleInfo
	db_battleInfo.FightId = js_battleinfo.Id
	db_battleInfo.RecordType = js_battleinfo.Type
	db_battleInfo.BattleInfo = utils.HF_CompressAndBase64(game.HF_JtoB(&js_battleinfo))
	db_battleInfo.BattleRecord = utils.HF_CompressAndBase64(game.HF_JtoB(&js_battleRecord))
	db_battleInfo.UpdateTime = time.Now().Unix()
	db.InsertTable("tbl_crossarenarecord", &db_battleInfo, 0, false)

	return 0
}
func (self *UnionInfo) CheckHunterFightEnd(endtime int64) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	change := false
	for _, v := range self.huntInfo {
		if UNION_HUNT_TYPE_NOMAL == v.Type {
			if v.EndTime == 0 {
				v.EndTime = endtime
				change = true
			}
			continue
		} else {
			if v.EndTime > 0 && time.Now().Unix() >= v.EndTime {
				v.EndTime = 0
				self.record = append(self.record, &JS_UnionRecord{UNION_RECORD_TYPE_HUNTER_BOSS_LEAVE, time.Now().Unix(), self.Mastername, fmt.Sprintf("%d", v.Type)})
				if v.TopDps.Len() <= 0 {
					continue
				}

				change = true
				for _, g := range self.member {
					find := false
					for _, u := range v.TopDps {
						if u.Uid == g.Uid {
							find = true
							break
						}
					}
					if !find {
						continue
					}

					core.GetCenterApp().AddEvent(g.ServerID, core.UNION_EVENT_UNION_HUNTER_AWARD, g.Uid,
						v.TopDps[0].Dps, v.Type, "")
				}
			}
		}
	}

	return change
}

//// 开启战斗协程
//func (self *UnionInfo) GetRefreshTime() int64 {
//	if time.Now().Hour() < 5 {
//		return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 5, 0, 0, 0, time.Now().Location()).Unix()
//	} else {
//		return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 5, 0, 0, 0, time.Now().Location()).Unix() + utils.DAY_SECS
//	}
//}

func (self *UnionInfo) RefreshActivityLimit() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	now := time.Now()
	timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()

	for _, v := range self.member {
		nSize := len(v.ActivityRecord)
		for i := nSize - 1; i >= 0; i-- {
			// 清理过期的限制表
			if (v.ActivityRecord[i].Time-timeStamp)/utils.DAY_SECS > UNION_ACTIVITY_LIMIT_CLEAN {
				v.ActivityRecord = append(v.ActivityRecord[:i],
					v.ActivityRecord[i+1:]...)
			}
		}
	}
}

//获取军团玩家的军团等级 1-4 会长 长老 精英 成员
func (self *UnionInfo) OnRefresh() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.DayExp = 0
	self.AcitvityLimit = 0
}

// 开启战斗协程
func (self *UnionInfo) OnHunterRefresh(endtime int64) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	change := false
	for _, v := range self.huntInfo {
		if UNION_HUNT_TYPE_NOMAL == v.Type {
			if v.EndTime == endtime {
				continue
			} else {
				v.EndTime = endtime
				self.record = append(self.record, &JS_UnionRecord{UNION_RECORD_TYPE_HUNTER_BOSS_LEAVE, time.Now().Unix(), self.Mastername, fmt.Sprintf("%d", v.Type)})
			}

			if v.TopDps.Len() <= 0 {
				v.TopDps = lstUnionHuntDpsTop{}
				continue
			}

			for _, g := range self.member {
				find := false
				for _, u := range v.TopDps {
					if u.Uid == g.Uid {
						find = true
						break
					}
				}
				if !find {
					continue
				}

				core.GetCenterApp().AddEvent(g.ServerID, core.UNION_EVENT_UNION_HUNTER_AWARD, g.Uid,
					v.TopDps[0].Dps, v.Type, "")

			}
			change = true
			v.TopDps = lstUnionHuntDpsTop{}
		}
	}

	return change
}

//! 加活跃度
func (self *UnionInfo) AddActivityPoint(uid int64, count int) int {
	config, ok := GetUnionMgr().CommunityConfigs[self.Level]
	if !ok || config == nil {
		return 0
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()

	addCount := count

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

	now := time.Now()
	timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
	if time.Now().Hour() < 5 {
		timeStamp -= utils.DAY_SECS
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

//! 加经验
func (self *UnionInfo) AddExp(exp int, addtype int) (int, int) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	csv, _ := GetUnionMgr().CommunityConfigs[self.Level]
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
		nextcsv, ok := GetUnionMgr().CommunityConfigs[self.Level+1]
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

func (self *UnionInfo) UpdateUnion() {
	//send := make(map[int]int)
	//
	//for _, v := range self.member {
	//	_, ok := send[v.ServerID]
	//	if ok {
	//		continue
	//	}
	//	core.GetCenterApp().AddEvent(v.ServerID, core.UNION_EVENT_UNION_UPDATE, 0,
	//		0, self.Id, "")
	//
	//	send[v.ServerID] = v.ServerID
	//}

}
func (self *UnionInfo) SendMail(title string, text string) {
	content := utils.HF_JtoA(UnionMail{title, text})
	save := utils.HF_JtoA(UnionMailSave{title, text, time.Now().Unix()})
	self.LastMail = save
	for _, v := range self.member {
		core.GetCenterApp().AddEvent(v.ServerID, core.UNION_EVENT_UNION_SEND_MAIL, v.Uid,
			0, 0, content)
	}

}

// 增加狩猎钻石奖励公会记录
func (self *UnionInfo) AddGemAwardRecord(name string, position, nType, count int) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.record = append(self.record, &JS_UnionRecord{UNION_RECORD_TYPE_HUNTER_BOSS_AWARD,
		time.Now().Unix(),
		name,
		fmt.Sprintf("%d", nType*100000+position*1000+count)})

	return false
}
