package player

import (
	"encoding/json"
	"fmt"
	"master/center/match"
	"master/center/tower"
	"master/center/union"
	"master/core"
	"master/db"
	"master/utils"
	"strconv"
	"strings"
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

func (self *CVarList) AddData(data interface{}) {
	self.Data = append(self.Data, data)
}

//! 错误码定义
const (
	UNION_SUCCESS      = 0 //! 没有错误
	UNION_NOT_MASTER   = 1 //! 不是会长
	UNION_NOT_GEM      = 2 //! 钻石不足
	UNION_NAME_ILLEGAL = 3 //! 名字违规
	UNION_LEVEL_LESS   = 4 //! 等级不足
	UNION_NO_UNION     = 5 //! 公会不存在
	UNION_IS_MEMBER    = 6 //! 玩家已经是成员
	UNION_NO_PLAYER    = 7 //! 玩家不存在
)

const (
	UNION_GM_TYPE_EXPLIMIT = 0
	UNION_GM_TYPE_EXP      = 1
	UNION_GM_TYPE_LEVEL    = 2
	UNION_GM_TYPE_ACTIVITY = 3
)

//! 操作请求
type RPC_UnionAction struct {
	Data string
}

//! 操作响应
type RPC_UnionActionRet struct {
	RetCode int //! 结果码
	Data    string
}

///////////获得公会//////////////
type S2M_UnionGetUnion struct {
	Unionuid int
}

type M2S_UnionGetUnion struct {
	Data union.MSG_UnionInfo
}

///////////获得公会列表//////////////
type S2M_UnionGetUnionList struct {
	ServerID int
}

type M2S_UnionGetUnionList struct {
	UnionList []union.JS_Union2
}

///////////公会改名//////////////
type S2M_UnionAlertName struct {
	Uid      int64
	Unionuid int
	Icon     int
	Name     string
}

type M2S_UnionAlertName struct {
	Ret int
}

///////////改内部公告//////////////
type S2M_UnionAlertNotice struct {
	Uid      int64
	Unionuid int
	Content  string
}

type M2S_UnionAlertNotice struct {
}

///////////改外部公告//////////////
type S2M_UnionAlertBoard struct {
	Uid      int64
	Unionuid int
	Content  string
}

type M2S_UnionAlertBoard struct {
}

///////////改设置//////////////
type S2M_UnionAlertSet struct {
	Uid      int64
	Unionuid int
	Type     int
	Level    int
}

type M2S_UnionAlertSet struct {
}

///////////申请//////////////
type S2M_UnionApply struct {
	Uid      int64
	Unionuid int
}

type M2S_UnionApply struct {
}

///////////取消申请//////////////
type S2M_UnionCancelApply struct {
	Uid      int64
	Unionuid int
}
type M2S_UnionCancelApply struct {
}

///////////创建公会//////////////
type S2M_UnionCreateUnion struct {
	Icon       int
	Unionname  string
	Masteruid  int64
	Mastername string
	Joinlevel  int
	HuntInfo   []*union.JS_UnionHunt
}
type M2S_UnionCreateUnion struct {
	Unionuid int
}

////////// 更新或者添加玩家/////////////

type S2M_UnionUpdateMember struct {
	Unionid  int
	Uid      int64
	Position int
	IsAdd    bool
	IsCreate bool
}

type M2S_UnionUpdateMember struct {
	Uid      int64
	Level    int
	UName    string
	IconId   int
	Portrait int
	Vip      int
	Fight    int64
	Position int
	Stage    int
}

////////// 更新玩家状态/////////////

type S2M_UnionUpdateMemberState struct {
	Unionid int
	Uid     int64
}

type M2S_UnionUpdateMemberState struct {
	Uid    int64
	Fight  int64
	Vip    int
	Online bool
}

///////////解散//////////////
type S2M_UnionDissolve struct {
	Uid      int64
	Unionuid int
}

type M2S_UnionDissolve struct {
}

///////////获得时间///////////////

type S2M_UnionGetTime struct {
	Type     int
	Unionuid int
}

type M2S_UnionGetTime struct {
	Time int64
}

///////////设置时间///////////////

type S2M_UnionSetTime struct {
	Type     int
	Unionuid int
	Time     int64
}

type M2S_UnionSetTime struct {
}

///////////检查会长///////////////

type S2M_UnionCheckMaster struct {
	Unionuid int
}

type M2S_UnionCheckMaster struct {
}

///////////加入//////////////
type S2M_UnionJoin struct {
	Uid      int64
	Unionuid int
}

type M2S_UnionJoin struct {
}

///////////会长拒绝//////////////
type S2M_UnionMasterFail struct {
	Uid      int64
	Unionuid int
	ApplyUid int64
}

type M2S_UnionMasterFail struct {
}

///////////会长同意//////////////
type S2M_UnionMasterOK struct {
	Uid      int64
	Unionuid int
	ApplyUid int64
}

type M2S_UnionMasterOK struct {
	ApplyUid []int64
}

///////////会长踢人//////////////
type S2M_UnionKickPlayer struct {
	Uid      int64
	Unionuid int
	OutUid   int64
}

type M2S_UnionKickPlayer struct {
}

///////////离开公会//////////////
type S2M_UnionOutPlayer struct {
	Unionuid int
	OutUid   int64
	IsMaster bool
}

type M2S_UnionOutPlayer struct {
}

/////////////会长辞职//////////////
//type S2M_UnionMasterResign struct {
//	Unionuid int
//	Uid      int64
//	Op       int
//}
//
//type M2S_UnionMasterResign struct {
//	Destuid int64
//	Name    string
//}

///////////修改职位//////////////
type S2M_UnionUnionModify struct {
	Unionuid int
	Uid      int64
	Destuid  int64
	Op       int
}

type M2S_UnionUnionModify struct {
}

///////////公会换人//////////////
type S2M_UnionChange struct {
	Unionuid int
	Uid      int64
	Destuid  int64
}

type M2S_UnionChange struct {
	Destuid int64
	Name    string
}

///////////设置公会之手//////////////
type S2M_UnionSetBraveHand struct {
	Unionuid int
	Uid      int64
	Destuid  int64
	Op       int
}

type M2S_UnionSetBraveHand struct {
}

///////////开启狩猎//////////////
type S2M_UnionOpenHunter struct {
	Uid      int64
	Unionuid int
	Type     int
	Cost     int
	Time     int
}

type M2S_UnionOpenHunter struct {
	UnionHunter *union.JS_UnionHunt
}

///////////开始狩猎//////////////
type S2M_UnionStartHunter struct {
	Uid      int64
	Unionuid int
	Type     int
}

type M2S_UnionStartHunter struct {
	EndTime int64
}

///////////开启狩猎//////////////
type S2M_UnionEndHunter struct {
	Uid      int64
	Unionuid int
	Type     int
}

type M2S_UnionEndHunter struct {
	UnionHunter *union.JS_UnionHunt
}

///////////增加战斗记录//////////////
type S2M_UnionAddDamage struct {
	Uid      int64
	Unionuid int
	Dps      int64
	Type     int
	FightID  int64
	Info     string
	Record   string
}

type M2S_UnionAddDamage struct {
	UnionHunter *union.JS_UnionHunt
}

///////////获得信息//////////////
type S2M_UnionGetBattleInfo struct {
	FightID int64
}

type M2S_UnionGetBattleInfo struct {
	Info string
}

///////////获得结构//////////////
type S2M_UnionGetBattleRecord struct {
	FightID int64
}

type M2S_UnionGetBattleRecord struct {
	Record string
}

///////////检查狩猎结束//////////////
type S2M_UnionCheckHunterEnd struct {
	Unionuid int
}

type M2S_UnionCheckHunterEnd struct {
}

///////////狩猎每日刷新//////////////
type S2M_UnionHunterRefresh struct {
	Unionuid int
	EndTime  int64
}

type M2S_UnionHunterRefresh struct {
}

///////////刷新活跃度限制//////////////
type S2M_UnionRefreshActivityLimit struct {
	Unionuid int
}

type M2S_UnionRefreshActivityLimit struct {
}

///////////增加活跃度//////////////
type S2M_UnionAddUnionActivity struct {
	Unionuid int
	Uid      int64
	Count    int
}

type M2S_UnionAddUnionActivity struct {
	Count int
}

///////////增加经验//////////////
type S2M_UnionAddUnionExp struct {
	Unionuid int
	Count    int
}

type M2S_UnionAddUnionExp struct {
}

///////////检查退出//////////////
type S2M_UnionCheckOut struct {
	Uid      int64
	Unionuid int
}

type M2S_UnionCheckOut struct {
}

///////////公会gm命令//////////////
type S2M_UnionGMAdd struct {
	Unionuid int
	Count    int
	Type     int
}

type M2S_UnionGMAdd struct {
}

///////////公会邮件//////////////
type S2M_UnionSendMail struct {
	Unionuid int
	Uid      int64
	Title    string
	Text     string
}

type M2S_UnionSendMail struct {
}

///////////通过名字搜索//////////////
type S2M_UnionGetUnionByName struct {
	Name string
}

type M2S_UnionGetUnionByName struct {
	Data []union.JS_Union2
}

///////////添加狩猎奖励公会日志//////////////
type S2M_UnionAddGemAwardRecord struct {
	Unionuid int
	Name     string
	Position int
	Type     int
	Count    int
}

type M2S_UnionAddGemAwardRecord struct {
}

///////////GM改内部公告//////////////
type S2M_GMUnionAlertNotice struct {
	Unionuid int
	Content  string
}

type M2S_GMUnionAlertNotice struct {
}

///////////GM改外部公告//////////////
type S2M_GMUnionAlertBoard struct {
	Unionuid int
	Content  string
}

type M2S_GMUnionAlertBoard struct {
}

////////////////////////////////

type RPC_Union struct {
}

func (self *RPC_Union) GetUnion(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionGetUnion
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	myunion.Encode()

	player := GetPlayerMgr().GetPlayer(myunion.Masteruid, true)
	if player != nil {
		myunion.Mastername = player.Data.UName
	}
	//同步

	var backmsg M2S_UnionGetUnion
	backmsg.Data = union.MSG_UnionInfo{myunion.Id,
		myunion.Icon,
		myunion.Unionname,
		myunion.Masteruid,
		myunion.Mastername,
		myunion.Level,
		myunion.Jointype,
		myunion.Joinlevel,
		myunion.ServerID,
		myunion.Notice,
		myunion.Board,
		myunion.Createtime,
		myunion.Lastupdtime,
		myunion.Fight,
		myunion.Exp,
		myunion.DayExp,
		myunion.ActivityPoint,
		myunion.AcitvityLimit,
		myunion.MailCD,
		myunion.Member,
		myunion.Applys,
		myunion.Record,
		myunion.HuntInfo,
		myunion.BraveHand,
		myunion.ChangeMaster,
	}

	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)

	if ret.Data == "" {
		fmt.Println("!!!!!!!!!!!!!")
	}

	return nil
}

//军团改名 ret：0成功 1只有会长才能操作 2钻石不足 3:已有军团名 4
func (self *RPC_Union) AlertUnionName(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionAlertName
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if nil == master {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if master.Position > union.UNION_POSITION_VICE_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if utils.HF_IsLicitName([]byte(msg.Name)) == false {
		ret.RetCode = UNION_NAME_ILLEGAL
		return nil
	}

	if union.GetUnionMgr().CheckName(msg.Name, myunion.Id) {
		ret.RetCode = UNION_NAME_ILLEGAL
		return nil
	}

	if myunion.Unionname == msg.Name && myunion.Icon == msg.Icon {
		ret.RetCode = UNION_NAME_ILLEGAL
		return nil
	}

	alerttype := 0
	if myunion.Unionname != msg.Name && myunion.Icon == msg.Icon {
		alerttype = union.UNION_ALERT_TYPE_NAME
	} else if myunion.Unionname == msg.Name && myunion.Icon != msg.Icon {
		alerttype = union.UNION_ALERT_TYPE_ICON
	} else if myunion.Unionname != msg.Name && myunion.Icon != msg.Icon {
		alerttype = union.UNION_ALERT_TYPE_BOTH
	}

	myunion.AlertUnionName(msg.Name, msg.Icon)

	var backmsg M2S_UnionAlertName
	backmsg.Ret = alerttype
	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)
	myunion.UpdateUnion()
	return nil
}

// 修改内部公告
func (self *RPC_Union) AlertUnionNotice(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionAlertNotice
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if nil == master {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if master.Position > union.UNION_POSITION_VICE_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	myunion.AlertUnionNotice(msg.Content)
	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil
}

// 修改外部公告
func (self *RPC_Union) AlertUnionBoard(req RPC_UnionAction, res *RPC_UnionActionRet) error {
	var msg S2M_UnionAlertBoard
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		res.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if nil == master {
		res.RetCode = UNION_NOT_MASTER
		return nil
	}

	if master.Position > union.UNION_POSITION_VICE_MASTER {
		res.RetCode = UNION_NOT_MASTER
		return nil
	}

	myunion.AlertUnionBoard(msg.Content)
	res.RetCode = UNION_SUCCESS
	return nil
}

// 修改公会设置
func (self *RPC_Union) AlertUnionSet(req RPC_UnionAction, res *RPC_UnionActionRet) error {
	var msg S2M_UnionAlertSet
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		res.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if nil == master {
		res.RetCode = UNION_NOT_MASTER
		return nil
	}

	if master.Position > union.UNION_POSITION_VICE_MASTER {
		res.RetCode = UNION_NOT_MASTER
		return nil
	}

	myunion.AlertUnionSet(msg.Type, msg.Level)
	res.RetCode = UNION_SUCCESS
	return nil
}

// 申请公会
func (self *RPC_Union) ApplyUnion(req RPC_UnionAction, res *RPC_UnionActionRet) error {
	var msg S2M_UnionApply
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		res.RetCode = UNION_NO_UNION
		return nil
	}

	me := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if nil == me {
		res.RetCode = UNION_NO_UNION
		return nil
	}

	if myunion.Jointype != 2 {
		res.RetCode = UNION_NO_UNION
		return nil
	}

	if me.Data.Level < myunion.Joinlevel {
		res.RetCode = UNION_LEVEL_LESS
		return nil
	}

	player := &union.JS_UnionApply{me.Data.data.UId,
		me.Data.data.Level,
		me.Data.data.UName,
		me.Data.data.IconId,
		me.Data.data.Portrait,
		me.Data.data.Vip,
		me.Data.data.Fight,
		time.Now().Unix(),
		me.Data.data.ServerId}
	myunion.AddApply(player)
	myunion.UpdateUnion()
	return nil
}

// 取消申请
func (self *RPC_Union) CancelApply(req RPC_UnionAction, res *RPC_UnionActionRet) error {
	var msg S2M_UnionCancelApply
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		res.RetCode = UNION_NO_UNION
		return nil
	}

	if !myunion.CancelApply(msg.Uid) {
		res.RetCode = UNION_NO_UNION
		return nil
	}
	return nil
}

func (self *RPC_Union) CreateUnion(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionCreateUnion
	json.Unmarshal([]byte(req.Data), &msg)

	player := GetPlayerMgr().GetPlayer(msg.Masteruid, true)
	if nil == player {
		ret.RetCode = UNION_NO_PLAYER
		return nil
	}

	var backmsg M2S_UnionCreateUnion
	if union.GetUnionMgr().CheckName(msg.Unionname, 0) {
		ret.RetCode = UNION_NAME_ILLEGAL
		return nil
	} else {
		backmsg.Unionuid = union.GetUnionMgr().CreateUnion(msg.Icon, msg.Unionname, msg.Masteruid, msg.Mastername, msg.HuntInfo, player.Data.data.ServerId)
		player.Data.data.UnionId = backmsg.Unionuid
	}
	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) UpdateMember(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionUpdateMember
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	me := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if me == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	if myunion.IsMember(msg.Uid) && msg.IsAdd {
		ret.RetCode = UNION_IS_MEMBER
		return nil
	}

	player := &union.JS_UnionMember{me.Data.data.UId,
		me.Data.data.Level,
		me.Data.data.UName,
		me.Data.data.IconId,
		me.Data.data.Portrait,
		me.Data.data.Vip,
		msg.Position,
		me.Data.data.Fight,
		time.Now().Unix(),
		[]*union.UserActivityRecord{},
		0,
		me.Data.data.PassId,
		me.Data.data.ServerId}

	position := myunion.UpdateMember(player, msg.IsAdd, msg.IsCreate)

	if msg.IsAdd {
		me.Data.data.UnionId = msg.Unionid
	}

	var backmsg M2S_UnionUpdateMember
	backmsg.Uid = me.Data.data.UId
	backmsg.Level = me.Data.data.Level
	backmsg.UName = me.Data.data.UName
	backmsg.IconId = me.Data.data.IconId
	backmsg.Portrait = me.Data.data.Portrait
	backmsg.Vip = me.Data.data.Vip
	backmsg.Fight = me.Data.data.Fight
	backmsg.Position = position
	backmsg.Stage = me.Data.data.PassId

	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) UpdateMemberState(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionUpdateMemberState
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	me := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if me == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	if !myunion.IsMember(msg.Uid) {
		ret.RetCode = UNION_IS_MEMBER
		return nil
	}

	myunion.UpdateMemberState(me.Data.data.UId, me.Data.data.Fight, me.Data.data.Vip, me.Online == 1)

	var backmsg M2S_UnionUpdateMemberState
	backmsg.Uid = me.Data.data.UId
	backmsg.Vip = me.Data.data.Vip
	backmsg.Fight = me.Data.data.Fight
	backmsg.Online = me.Online == 1

	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) DissolveUnion(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionDissolve
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	if myunion.Masteruid != msg.Uid {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if myunion.GetMemberLen() > 1 {
		ret.RetCode = UNION_NO_PLAYER
		return nil
	}

	union.GetUnionMgr().DissolveUnion(msg.Unionuid)
	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil
}

func (self *RPC_Union) GetUnionTime(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionGetTime
	json.Unmarshal([]byte(req.Data), &msg)

	var backmsg M2S_UnionGetTime
	if msg.Type == union.UNION_GET_TIME_TYPE_CALL_TIME {
		backmsg.Time = union.GetUnionMgr().GetUnionCallTime(msg.Unionuid)
	} else {
		backmsg.Time = union.GetUnionMgr().GetUnionCheckMaster(msg.Unionuid)
	}

	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) SetUnionTime(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionSetTime
	json.Unmarshal([]byte(req.Data), &msg)

	if msg.Type == union.UNION_GET_TIME_TYPE_CALL_TIME {
		union.GetUnionMgr().SetUnionCallTime(msg.Unionuid, msg.Time)
	} else {
		union.GetUnionMgr().SetUnionCheckMaster(msg.Unionuid, msg.Time)
	}

	ret.RetCode = UNION_SUCCESS
	return nil
}

func (self *RPC_Union) CheckMasterOffline(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionCheckMaster
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	myunion.CheckMasterOffline()
	ret.RetCode = UNION_SUCCESS
	return nil
}

// 申请公会
func (self *RPC_Union) JoinUnion(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionJoin
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	me := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if nil == me {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	if me.Data.data.UnionId != 0 {
		ret.RetCode = UNION_IS_MEMBER
		return nil
	}

	if myunion.Jointype != 0 {
		ret.RetCode = UNION_LEVEL_LESS
		return nil
	}

	if me.Data.data.Level < myunion.Joinlevel {
		ret.RetCode = UNION_LEVEL_LESS
		return nil
	}

	if union.GetUnionMgr().CheckFull(myunion) {
		ret.RetCode = UNION_LEVEL_LESS
		return nil
	}

	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil
}

func (self *RPC_Union) MasterFail(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionMasterFail
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if nil == master {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if master.Position > union.UNION_POSITION_VICE_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	myunion.MasterFail(msg.ApplyUid)
	myunion.CleanApply(msg.ApplyUid)

	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil
}

func (self *RPC_Union) MasterOk(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionMasterOK
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if nil == master {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if master.Position > union.UNION_POSITION_VICE_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if union.GetUnionMgr().CheckFull(myunion) {
		ret.RetCode = UNION_LEVEL_LESS
		return nil
	}

	csv_community := union.GetUnionMgr().CommunityConfigs[myunion.Level]
	maxCount := csv_community.Membernum - myunion.GetMemberLen()
	count := 0
	applys := []int64{}
	if msg.ApplyUid == 0 {
		//var applys []*union.JS_UnionApply
		unionapplys := myunion.GetApply()
		for _, v := range unionapplys {
			member := GetPlayerMgr().GetPlayer(v.Uid, true)
			if nil == member {
				continue
			}

			if member.Data.data.UnionId != 0 {
				continue
			}

			if count >= maxCount {
				break
			}

			myunion.MasterOK(v.Uid)
			applys = append(applys, v.Uid)
			count++
		}
	} else {
		member := GetPlayerMgr().GetPlayer(msg.ApplyUid, true)
		if nil == member {
			ret.RetCode = UNION_NO_UNION
			return nil
		}

		if member.Data.data.UnionId != 0 {
			ret.RetCode = UNION_IS_MEMBER
			return nil
		}

		myunion.MasterOK(msg.ApplyUid)
		applys = append(applys, msg.ApplyUid)
	}

	myunion.CleanApply(msg.ApplyUid)

	var backmsg M2S_UnionMasterOK
	backmsg.ApplyUid = applys
	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	ret.Data = utils.HF_JtoA(backmsg)

	return nil
}

func (self *RPC_Union) KickPlayer(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionKickPlayer
	json.Unmarshal([]byte(req.Data), &msg)
	// 获得公会
	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}
	// 不能踢会长
	if myunion.Masteruid == msg.OutUid {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 不能踢自己
	if msg.Uid == msg.OutUid {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 获得操作者
	master := myunion.GetMember(msg.Uid)
	if nil == master {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 权限不足
	if master.Position > union.UNION_POSITION_VICE_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 获得被操作者
	outplayer := myunion.GetMember(msg.OutUid)
	if nil == outplayer {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 不能操作比自己权限高或同级的人
	if master.Position >= outplayer.Position {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	//myunion.KickPlayer(msg.OutUid)

	ret.RetCode = UNION_SUCCESS
	//myunion.UpdateUnion()
	return nil
}

func (self *RPC_Union) OutPlayer(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionOutPlayer
	json.Unmarshal([]byte(req.Data), &msg)
	// 获得公会
	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}
	// 会长不能退出
	if myunion.Masteruid == msg.OutUid {
		ret.RetCode = UNION_NO_UNION
		return nil
	}
	// 获得被操作者
	outplayer := GetPlayerMgr().GetPlayer(msg.OutUid, true)
	if nil == outplayer {
		ret.RetCode = UNION_NO_UNION
		return nil
	}
	outmember := myunion.GetMember(msg.OutUid)
	if nil == outmember {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 检测被操作者职位
	if outmember.Position == union.UNION_POSITION_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 踢出玩家
	myunion.OutPlayer(msg.OutUid, msg.IsMaster)
	outplayer.Data.data.UnionId = 0

	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil
}

//
//func (self *RPC_Union) MasterResign(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
//	var msg S2M_UnionMasterResign
//	json.Unmarshal([]byte(req.Data), &msg)
//
//	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
//	if myunion == nil {
//		ret.RetCode = UNION_NO_UNION
//		return nil
//	}
//
//	if myunion.GetMemberLen() <= 1 {
//		ret.RetCode = UNION_NOT_MASTER
//		return nil
//	}
//
//	if myunion.Masteruid != msg.Uid {
//		ret.RetCode = UNION_NOT_MASTER
//		return nil
//	}
//
//	destuid := myunion.GetIdWithoutMaster()
//	if destuid == 0 {
//		ret.RetCode = UNION_NOT_MASTER
//		return nil
//	}
//
//	member := myunion.GetMember(destuid)
//	if member == nil {
//		ret.RetCode = UNION_NOT_MASTER
//		return nil
//	}
//
//	//myunion.UnionChange(msg.Uid, destuid, member.Uname)
//	//
//	//core.GetCenterApp().AddEvent(member.ServerID, core.UNION_EVENT_UNION_MODIFY, destuid,
//	//	msg.Uid, union.UNION_POSITION_MASTER, "")
//	//
//	//core.GetCenterApp().AddEvent(member.ServerID, core.UNION_EVENT_UNION_MODIFY, msg.Uid,
//	//	msg.Uid, union.UNION_POSITION_MEMBER, "")
//
//	var backmsg M2S_UnionMasterResign
//	backmsg.Destuid = destuid
//	backmsg.Name = member.Uname
//	ret.RetCode = UNION_SUCCESS
//	myunion.UpdateUnion()
//	ret.Data = utils.HF_JtoA(backmsg)
//
//	return nil
//}

func (self *RPC_Union) UnionModify(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionUnionModify
	json.Unmarshal([]byte(req.Data), &msg)
	// 公会不存在
	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}
	// 获得操作者
	master := myunion.GetMember(msg.Uid)
	if master == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 获得被操作者
	member := myunion.GetMember(msg.Destuid)
	if member == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 操作者只能是会长
	if master.Position > union.UNION_POSITION_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 操作对象不能是自己
	if msg.Uid == msg.Destuid {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 不能操作同级或比自己权限高的人
	if master.Position >= member.Position {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 不能给别人设置比自己更高的权限
	if msg.Op < master.Position {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	myunion.UnionModify(msg.Destuid, msg.Op)

	core.GetCenterApp().AddEvent(member.ServerID, core.UNION_EVENT_UNION_MODIFY, msg.Destuid,
		msg.Uid, msg.Op, strconv.FormatInt(msg.Destuid, 10))

	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil
}

func (self *RPC_Union) UnionChange(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionChange
	json.Unmarshal([]byte(req.Data), &msg)
	// 公会不存在
	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}
	// 操作对象不能是自己
	if msg.Uid == msg.Destuid {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 获得操作者
	master := myunion.GetMember(msg.Uid)
	if master == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 获得被操作者
	member := myunion.GetMember(msg.Destuid)
	if member == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 操作者只能是会长
	if master.Position > union.UNION_POSITION_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 判断公会会长uid
	if myunion.Masteruid != msg.Uid {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 并且只能转让给副会长
	if member.Position != union.UNION_POSITION_VICE_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	// 设置公会转让
	myunion.UnionChange(msg.Uid, msg.Destuid, member.Uname)
	// 通知被操作者 公会转让消息
	core.GetCenterApp().AddEvent(member.ServerID, core.UNION_EVENT_UNION_MODIFY, msg.Destuid,
		msg.Uid, union.UNION_POSITION_MASTER, strconv.FormatInt(msg.Destuid, 10))

	var backmsg M2S_UnionChange
	backmsg.Destuid = member.Uid
	backmsg.Name = member.Uname
	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) SetBraveHand(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionSetBraveHand
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if master == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	member := myunion.GetMember(msg.Destuid)
	if member == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if master.Position > union.UNION_POSITION_VICE_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if msg.Op == 1 {
		if member.BraveHand == 1 {
			ret.RetCode = UNION_NOT_MASTER
			return nil
		}
		//if !GetCsvMgr().IsLevelAndPassOpenNew(destuid, OPEN_LEVEL_HIRE) {
		//	ret.RetCode = UNION_NOT_MASTER
		//	return nil
		//}
	} else {
		if member.BraveHand == 0 {
			ret.RetCode = UNION_NOT_MASTER
			return nil
		}
	}

	csv, ok := union.GetUnionMgr().CommunityConfigs[myunion.Level]
	if !ok {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	union.GetUnionMgr().CheckBraveHand(myunion)

	if msg.Op == 1 && myunion.GetBraveHandLen() >= csv.Fearless {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	myunion.SetBraveHand(msg.Destuid, msg.Op)

	core.GetCenterApp().AddEvent(member.ServerID, core.UNION_EVENT_UNION_SET_BRAVE_HAND, msg.Destuid,
		msg.Uid, msg.Op, strconv.FormatInt(msg.Destuid, 10))

	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil

}

func (self *RPC_Union) OpenHuntFight(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionOpenHunter
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if master == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if master.Position > union.UNION_POSITION_VICE_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	unionHunter := myunion.GetHunter(msg.Type)
	if nil == unionHunter {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if 0 != unionHunter.EndTime {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	nowTime := time.Now().Unix()
	if nowTime <= unionHunter.EndTime {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if msg.Cost <= 0 {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if myunion.ActivityPoint < msg.Cost {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	myunion.MinActivityPoint(msg.Cost)
	if !myunion.OpenHunterFight(msg.Type, nowTime+int64(msg.Time)) {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}
	var backmsg M2S_UnionOpenHunter
	backmsg.UnionHunter = unionHunter
	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) StartHuntFight(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionStartHunter
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if master == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	unionHunter := myunion.GetHunter(msg.Type)
	if nil == unionHunter {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if 0 == unionHunter.EndTime {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	nowTime := time.Now().Unix()
	if nowTime >= unionHunter.EndTime {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	var backmsg M2S_UnionStartHunter
	backmsg.EndTime = unionHunter.EndTime
	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) EndHuntFight(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionEndHunter
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if master == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	unionHunter := myunion.GetHunter(msg.Type)
	if nil == unionHunter {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	//if 0 == unionHunter.EndTime {
	//	ret.RetCode = UNION_NOT_MASTER
	//	return nil
	//}
	//
	//nowTime := time.Now().Unix()
	//if nowTime <= unionHunter.EndTime {
	//	ret.RetCode = UNION_NOT_MASTER
	//	return nil
	//}

	var backmsg M2S_UnionEndHunter
	backmsg.UnionHunter = unionHunter
	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) AddUnionHuntDamage(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionAddDamage
	json.Unmarshal([]byte(req.Data), &msg)

	player := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if nil == player {
		ret.RetCode = UNION_NO_PLAYER
		return nil
	}

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if master == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	unionHunter := myunion.GetHunter(msg.Type)
	if nil == unionHunter {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	member := &union.JS_UnionMember{player.Data.data.UId,
		player.Data.data.Level,
		player.Data.data.UName,
		player.Data.data.IconId,
		player.Data.data.Portrait,
		player.Data.data.Vip,
		master.Position,
		player.Data.data.Fight,
		time.Now().Unix(),
		[]*union.UserActivityRecord{},
		0,
		player.Data.data.PassId,
		player.Data.data.ServerId}

	myunion.AddUnionHuntDamage(member, msg.Type, msg.Dps, msg.Uid, msg.FightID, msg.Info, msg.Record)
	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil
}

func (self *RPC_Union) GetBattleInfo(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionGetBattleInfo
	json.Unmarshal([]byte(req.Data), &msg)

	value, flag, err := db.HGetRedisEx(`san_huntbattleinfo`, msg.FightID, fmt.Sprintf("%d", msg.FightID))
	flag = false
	if err != nil || !flag {
		var db_battleInfo match.JS_CrossArenaBattleInfo
		sql := fmt.Sprintf("select * from `tbl_crossarenarecord` where fightid=%d limit 1;", msg.FightID)
		ret1 := db.GetDBMgr().DBUser.GetOneData(sql, &db_battleInfo, "", 0)
		if ret1 == true { //! 获取成功
			//! 进行处理
			var battleInfo tower.BattleInfo
			err1 := json.Unmarshal(utils.HF_Base64AndDecompress(db_battleInfo.BattleInfo), &battleInfo)
			if err1 != nil {
				ret.RetCode = UNION_NOT_MASTER
				return nil
			}

			if battleInfo.Id != 0 {
				value = utils.HF_JtoA(battleInfo)
			} else {
				ret.RetCode = UNION_NOT_MASTER
				return nil
			}
		} else {
			ret.RetCode = UNION_NOT_MASTER
			return nil
		}
	}

	var backmsg M2S_UnionGetBattleInfo
	backmsg.Info = value
	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) GetBattleRecord(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionGetBattleRecord
	json.Unmarshal([]byte(req.Data), &msg)

	value, flag, err := db.HGetRedisEx(`san_huntbattlerecord`, msg.FightID, fmt.Sprintf("%d", msg.FightID))
	if err != nil || !flag {
		var db_battleInfo match.JS_CrossArenaBattleInfo
		sql := fmt.Sprintf("select * from `tbl_crossarenarecord` where fightid=%d limit 1;", msg.FightID)
		ret1 := db.GetDBMgr().DBUser.GetOneData(sql, &db_battleInfo, "", 0)
		if ret1 == true { //! 获取成功
			//! 进行处理
			var battleRecord tower.BattleRecord
			err1 := json.Unmarshal(utils.HF_Base64AndDecompress(db_battleInfo.BattleRecord), &battleRecord)
			if err1 != nil {
				ret.RetCode = UNION_NOT_MASTER
				return nil
			}

			if battleRecord.Id != 0 {
				value = utils.HF_JtoA(battleRecord)
			} else {
				ret.RetCode = UNION_NOT_MASTER
				return nil
			}
		} else {
			ret.RetCode = UNION_NOT_MASTER
			return nil
		}
	}

	var backmsg M2S_UnionGetBattleRecord
	backmsg.Record = value
	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

//func (self *RPC_Union) CheckHunterFightEnd(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
//	var msg S2M_UnionCheckHunterEnd
//	json.Unmarshal([]byte(req.Data), &msg)
//
//	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
//	if myunion == nil {
//		ret.RetCode = UNION_NO_UNION
//		return nil
//	}
//
//	if !myunion.CheckHunterFightEnd() {
//		ret.RetCode = UNION_NO_UNION
//		return nil
//	}
//
//	var backmsg M2S_UnionCheckHunterEnd
//	ret.RetCode = UNION_SUCCESS
//	myunion.UpdateUnion()
//	ret.Data = utils.HF_JtoA(backmsg)
//	return nil
//}

//
//func (self *RPC_Union) OnHunterRefresh(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
//	var msg S2M_UnionHunterRefresh
//	json.Unmarshal([]byte(req.Data), &msg)
//
//	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
//	if myunion == nil {
//		ret.RetCode = UNION_NO_UNION
//		return nil
//	}
//
//	if !myunion.OnHunterRefresh(msg.EndTime) {
//		ret.RetCode = UNION_NO_UNION
//		return nil
//	}
//
//	ret.RetCode = UNION_SUCCESS
//	myunion.UpdateUnion()
//	return nil
//}
//func (self *RPC_Union) RefreshActivityLimit(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
//	var msg S2M_UnionRefreshActivityLimit
//	json.Unmarshal([]byte(req.Data), &msg)
//
//	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
//	if myunion == nil {
//		ret.RetCode = UNION_NO_UNION
//		return nil
//	}
//
//	myunion.RefreshActivityLimit()
//
//	ret.RetCode = UNION_SUCCESS
//	myunion.UpdateUnion()
//	return nil
//}

func (self *RPC_Union) AddUnionActivity(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionAddUnionActivity
	json.Unmarshal([]byte(req.Data), &msg)
	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	var backmsg M2S_UnionAddUnionActivity
	backmsg.Count = myunion.AddActivityPoint(msg.Uid, msg.Count)
	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)
	myunion.UpdateUnion()
	return nil
}

func (self *RPC_Union) AddUnionExp(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionAddUnionExp
	json.Unmarshal([]byte(req.Data), &msg)
	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	myunion.AddExp(msg.Count, 1)

	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil
}
func (self *RPC_Union) UnionCheckOut(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionCheckOut
	json.Unmarshal([]byte(req.Data), &msg)
	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	if !myunion.IsMember(msg.Uid) {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	ret.RetCode = UNION_SUCCESS
	return nil
}

func (self *RPC_Union) GetUnionList(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionGetUnionList
	json.Unmarshal([]byte(req.Data), &msg)

	var backmsg M2S_UnionGetUnionList
	nCount := 0
	allunion := union.GetUnionMgr().Sql_Union
	for _, v := range allunion {
		if nCount >= union.UNION_LIST_MAX {
			break
		}
		data := union.JS_Union2{v.Id,
			v.Icon,
			v.Unionname,
			v.Masteruid,
			v.Mastername,
			v.Level,
			v.Jointype,
			v.Joinlevel,
			0,
			v.GetMemberLen(),
			0,
			0,
			v.Fight,
			v.Exp,
			v.ActivityPoint}
		backmsg.UnionList = append(backmsg.UnionList, data)
		nCount++
	}

	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) UnionGmAdd(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionGMAdd
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	switch msg.Type {
	case UNION_GM_TYPE_EXPLIMIT:
		myunion.DayExp += msg.Count
	case UNION_GM_TYPE_EXP:
		myunion.Exp += msg.Count
	case UNION_GM_TYPE_LEVEL:
		myunion.Level = msg.Count
	case UNION_GM_TYPE_ACTIVITY:
		myunion.ActivityPoint += msg.Count
	}

	ret.RetCode = UNION_SUCCESS
	return nil
}
func (self *RPC_Union) UnionSendMail(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionSendMail
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	master := myunion.GetMember(msg.Uid)
	if master == nil {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	if master.Position > union.UNION_POSITION_VICE_MASTER {
		ret.RetCode = UNION_NOT_MASTER
		return nil
	}

	now := time.Now().Unix()
	if myunion.MailCD != 0 {
		cd := now - myunion.MailCD
		if cd < utils.DAY_SECS {
			ret.RetCode = UNION_NOT_MASTER
			return nil
		}
	}

	myunion.SendMail(msg.Title, msg.Text)
	myunion.MailCD = time.Now().Unix()

	ret.RetCode = UNION_SUCCESS
	return nil
}

func (self *RPC_Union) GetUnionByName(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionGetUnionByName
	json.Unmarshal([]byte(req.Data), &msg)

	var backmsg M2S_UnionGetUnionByName
	nCount := 0
	allunion := union.GetUnionMgr().Sql_Union
	for _, v := range allunion {
		if nCount >= union.UNION_LIST_MAX {
			break
		}
		if !strings.Contains(v.Unionname, msg.Name) {
			continue
		}
		data := union.JS_Union2{v.Id,
			v.Icon,
			v.Unionname,
			v.Masteruid,
			v.Mastername,
			v.Level,
			v.Jointype,
			v.Joinlevel,
			0,
			v.GetMemberLen(),
			0,
			0,
			v.Fight,
			v.Exp,
			v.ActivityPoint}
		backmsg.Data = append(backmsg.Data, data)
		nCount++
	}

	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}

func (self *RPC_Union) AddGemAwardRecord(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_UnionAddGemAwardRecord
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}
	myunion.AddGemAwardRecord(msg.Name, msg.Position, msg.Type, msg.Count)

	ret.RetCode = UNION_SUCCESS
	return nil
}

// 修改内部公告
func (self *RPC_Union) GMAlertUnionNotice(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_GMUnionAlertNotice
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		ret.RetCode = UNION_NO_UNION
		return nil
	}

	myunion.AlertUnionNotice(msg.Content)
	ret.RetCode = UNION_SUCCESS
	myunion.UpdateUnion()
	return nil
}

// 修改外部公告
func (self *RPC_Union) GMAlertUnionBoard(req RPC_UnionAction, res *RPC_UnionActionRet) error {
	var msg S2M_GMUnionAlertBoard
	json.Unmarshal([]byte(req.Data), &msg)

	myunion := union.GetUnionMgr().GetUnion(msg.Unionuid)
	if myunion == nil {
		res.RetCode = UNION_NO_UNION
		return nil
	}

	myunion.AlertUnionBoard(msg.Content)
	res.RetCode = UNION_SUCCESS
	return nil
}
