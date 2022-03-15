/*
@Time : 2020/5/10 11:10
@Author : 96121
@File : proto_player
@Software: GoLand
*/
package game

import (
	"net/rpc"
	"sync"
)

//! 错误码定义
const (
	UNION_SUCCESS      = 0 //! 没有错误
	UNION_NOT_MASTER   = 1 //! 不是会长
	UNION_NOT_GEM      = 2 //! 钻石不足
	UNION_NAME_ILLEGAL = 3 //! 名字违规
	UNION_LEVEL_LESS   = 4 //! 等级不足
	UNION_NO_UNION     = 5 //! 公会不存在
)

const (
	UNION_GM_TYPE_EXPLIMIT = 0
	UNION_GM_TYPE_EXP      = 1
	UNION_GM_TYPE_LEVEL    = 2
	UNION_GM_TYPE_ACTIVITY = 3
)
const (
	RPC_UNION_GET_UNION            = "RPC_Union.GetUnion"             //! 获得公会
	RPC_UNION_GET_UNION_LIST       = "RPC_Union.GetUnionList"         //! 获得公会
	RPC_UNION_GM_ADD               = "RPC_Union.UnionGmAdd"           //! gm命令
	RPC_UNION_ALERT_NAME           = "RPC_Union.AlertUnionName"       //! 修改公会名
	RPC_UNION_ALERT_NOTICE         = "RPC_Union.AlertUnionNotice"     //! 修改内部公告
	RPC_UNION_ALERT_BOARD          = "RPC_Union.AlertUnionBoard"      //! 修改外部公告
	RPC_UNION_ALERT_SET            = "RPC_Union.AlertUnionSet"        //! 修改公会设置
	RPC_UNION_APPLY                = "RPC_Union.ApplyUnion"           //! 申请公会
	RPC_UNION_CANCEL_APPLY         = "RPC_Union.CancelApply"          //! 取消申请
	RPC_UNION_CREATE               = "RPC_Union.CreateUnion"          //! 创建公会
	RPC_UNION_UPDATE_MEMBER        = "RPC_Union.UpdateMember"         //! 更新公会成员
	RPC_UNION_UPDATE_MEMBER_STATE  = "RPC_Union.UpdateMemberState"    //! 更新公会成员状态
	RPC_UNION_DISSOLVE             = "RPC_Union.DissolveUnion"        //! 解散公会
	RPC_UNION_GET_TIME             = "RPC_Union.GetUnionTime"         //! 获得时间
	RPC_UNION_SET_TIME             = "RPC_Union.SetUnionTime"         //! 获得时间
	RPC_UNION_CHECK_MASTER         = "RPC_Union.CheckMasterOffline"   //! 检查离线
	RPC_UNION_JOIN                 = "RPC_Union.JoinUnion"            //! 检测加入
	RPC_UNION_MASTER_FAIL          = "RPC_Union.MasterFail"           //! 会长拒绝
	RPC_UNION_MASTER_OK            = "RPC_Union.MasterOk"             //! 会长同意
	RPC_UNION_KICK_PLAYER          = "RPC_Union.KickPlayer"           //! 会长踢人
	RPC_UNION_OUT_PLAYER           = "RPC_Union.OutPlayer"            //! 执行离开函数
	RPC_UNION_MASTER_RESIGN        = "RPC_Union.MasterResign"         //! 会长辞职
	RPC_UNION_MODIFY               = "RPC_Union.UnionModify"          //! 修改职位
	RPC_UNION_CHANGE               = "RPC_Union.UnionChange"          //! 会长换人
	RPC_UNION_SET_BRAVE_HAND       = "RPC_Union.SetBraveHand"         //! 设置无畏之手
	RPC_UNION_OPEN_HUNTER          = "RPC_Union.OpenHuntFight"        //! 开启狩猎
	RPC_UNION_START_HUNTER         = "RPC_Union.StartHuntFight"       //! 开始狩猎
	RPC_UNION_END_HUNTER           = "RPC_Union.EndHuntFight"         //! 结束狩猎
	RPC_UNION_ADD_DAMAGE           = "RPC_Union.AddUnionHuntDamage"   //! 结束狩猎
	RPC_UNION_GET_INFO             = "RPC_Union.GetBattleInfo"        //! 获得info
	RPC_UNION_GET_RECORD           = "RPC_Union.GetBattleRecord"      //! 获得record
	RPC_UNION_CHECK_HUNTER_END     = "RPC_Union.CheckHunterFightEnd"  //! 检测狩猎结束
	RPC_UNION_ACTIVITY_REFRESH     = "RPC_Union.RefreshActivityLimit" //! 刷新活跃度限制
	RPC_UNION_HUNTER_REFRESH       = "RPC_Union.OnHunterRefresh"      //! 狩猎刷新
	RPC_UNION_ADD_ACTIVITY         = "RPC_Union.AddUnionActivity"     //! 增加活跃度
	RPC_UNION_ADD_EXP              = "RPC_Union.AddUnionExp"          //! 增加经验
	RPC_UNION_CHECK_OUT            = "RPC_Union.UnionCheckOut"        //! 检查退出
	RPC_UNION_SEND_MAIL            = "RPC_Union.UnionSendMail"        //! 发邮件
	RPC_UNION_GET_UNION_BY_NAME    = "RPC_Union.GetUnionByName"       //! 通过名字获得公会
	RPC_UNION_ADD_GEM_AWARD_RECORD = "RPC_Union.AddGemAwardRecord"    //! 添加狩猎奖励日志
	RPC_UNION_GM_ALERT_NOTICE      = "RPC_Union.GMAlertUnionNotice"   //! 修改内部公告
	RPC_UNION_GM_ALERT_BOARD       = "RPC_Union.GMAlertUnionBoard"    //! 修改外部公告
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

///////////获得公会//////////////
type S2M_UnionGetUnion struct {
	Unionuid int
}

type M2S_UnionGetUnion struct {
	Data MSG_UnionInfo
}

///////////获得公会列表//////////////
type S2M_UnionGetUnionList struct {
	ServerID int
}

type M2S_UnionGetUnionList struct {
	UnionList []JS_Union2
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
	HuntInfo   []*JS_UnionHunt
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
	UnionHunter *JS_UnionHunt
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

///////////结束狩猎//////////////
type S2M_UnionEndHunter struct {
	Uid      int64
	Unionuid int
	Type     int
}

type M2S_UnionEndHunter struct {
	UnionHunter *JS_UnionHunt
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
	UnionHunter *JS_UnionHunt
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
	Data []JS_Union2
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

//////////////////////////

type RPC_Union struct {
	Client       *rpc.Client
	PlayerLocker *sync.RWMutex //! 数据锁
}

func (self *RPC_Union) Init() bool {
	self.PlayerLocker = new(sync.RWMutex)
	return true
}

//func (self *RPC_Union) AlertUnionName(uid int64, union_uid int64, newname string, newicon int) int {
//	if self.Client != nil {
//		var req RPC_UnionAction
//		req.Uid = uid
//		req.UnionUid = union_uid
//		req.Content = newname
//		req.Number1 = newicon
//
//		var res RPC_UnionActionRet
//
//		self.Client.Call(RPC_ALERT_UNION_NAME, req, &res)
//		return res.RetCode
//	}
//
//	return RETCODE_UNKNOWN
//}
//
////修改信息
//func (self *RPC_Union) AlertUnionNotice(uid int64, union_uid int64, notice string) int {
//	if self.Client != nil {
//		var req RPC_UnionAction
//		req.Uid = uid
//		req.UnionUid = union_uid
//		req.Content = notice
//
//		var res RPC_UnionActionRet
//
//		self.Client.Call(RPC_ALERT_UNION_NOTICE, req, &res)
//		return res.RetCode
//	}
//
//	return RETCODE_UNKNOWN
//}
//
////修改信息
//func (self *RPC_Union) AlertUnionBoard(uid int64, union_uid int64, notice string) int {
//	if self.Client != nil {
//		var req RPC_UnionAction
//		req.Uid = uid
//		req.UnionUid = union_uid
//		req.Content = notice
//
//		var res RPC_UnionActionRet
//
//		self.Client.Call(RPC_ALERT_UNION_BOARD, req, &res)
//		return res.RetCode
//	}
//
//	return RETCODE_UNKNOWN
//}
//
////修改军团设置
//func (self *RPC_Union) AlertUnionSet(uid int64, union_uid int64, jointype int, joinlevel int) int {
//	if self.Client != nil {
//		var req RPC_UnionAction
//		req.Uid = uid
//		req.UnionUid = union_uid
//		req.Number1 = jointype
//		req.Number2 = joinlevel
//
//		var res RPC_UnionActionRet
//
//		self.Client.Call(RPC_ALERT_UNION_SET, req, &res)
//		return res.RetCode
//	}
//
//	return RETCODE_UNKNOWN
//}
//
////申请军团
//func (self *RPC_Union) ApplyUnion(uid int64, union_uid int64) int {
//	if self.Client != nil {
//		var req RPC_UnionAction
//		req.Uid = uid
//		req.UnionUid = union_uid
//
//		var res RPC_UnionActionRet
//
//		self.Client.Call(RPC_ALERT_UNION_APPLY, req, &res)
//		return res.RetCode
//	}
//
//	return RETCODE_UNKNOWN
//}
//
////取消申请
//func (self *RPC_Union) CancelApply(uid int64, union_uid int64) int {
//	if self.Client != nil {
//		var req RPC_UnionAction
//		req.Uid = uid
//		req.UnionUid = union_uid
//
//		var res RPC_UnionActionRet
//
//		self.Client.Call(RPC_ALERT_UNION_CANCEL_APPLY, req, &res)
//		return res.RetCode
//	}
//
//	return RETCODE_UNKNOWN
//}

//操作
func (self *RPC_Union) UnionAction(action string, data interface{}) *RPC_UnionActionRet {
	if self.Client != nil {
		var req RPC_UnionAction
		req.Data = HF_JtoA(data)

		var ret RPC_UnionActionRet
		GetMasterMgr().CallEx(self.Client,action, req, &ret)
		return &ret
	}

	return nil
}
