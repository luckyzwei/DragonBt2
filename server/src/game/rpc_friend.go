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
	RETCODE_OK                  = 0   //! 没有错误
	RETCODE_PLAYER_NOT_EXIST    = -1  //! 角色不存在
	RETCODE_FRIEND_NOT_EXIST    = -2  //! 好友不存在
	RETCODE_SELF_FRIEND_FULL    = -3  //! 自己好友已满
	RETCODE_TARGET_FRIEND_FULL  = -4  //! 目标好友已满
	RETCODE_FRIEND_APPLY_ERROR  = -5  //! 好友申请不存在
	RETCODE_FRIEND_APPLY_BLACK  = -6  //! 黑名单
	RETCODE_FRIEND_HIRE_ALREADY = -7  //! 佣兵租用中
	RETCODE_FRIEND_HIRE_MAX     = -8  //! 目标佣兵已满
	RETCODE_UNKNOWN             = -99 //! 连接错误
)

const (
	RPC_FINDFRIEND     = "RPC_Friend.FindFriend"     //! 寻找好友s
	RPC_APPLYFRIEND    = "RPC_Friend.ApplyFriend"    //! 增加好友 = Apply 申请好友
	RPC_DELFRIEND      = "RPC_Friend.DelFriend"      //! 删除好友
	RPC_AGREEFRIEND    = "RPC_Friend.AgreeFriend"    //! 同意好友请求
	RPC_REFUSEFRIEND   = "RPC_Friend.RefuseFriend"   //! 拒绝好友请求
	RPC_BLACKFRIEND    = "RPC_Friend.BlackFriend"    //! 拉黑好友
	RPC_BLACKOUTFRIEND = "RPC_Friend.BlackOutFriend" //! 移除黑名单好友
	RPC_POWERFRIEND    = "RPC_Friend.PowerFriend"    //! 友情点操作
	RPC_HIRELOSE       = "RPC_Friend.HireLose"       //! 放弃租用的英雄
	RPC_HIREAPPLY      = "RPC_Friend.HireApply"      //! 申请租用
	RPC_HIRECANCEL     = "RPC_Friend.HireCancel"     //! 放弃申请
	RPC_HIREAGREE      = "RPC_Friend.HireAgree"      //! 同意
	RPC_HIREREFUSE     = "RPC_Friend.HireRefuse"     //! 拒绝
	RPC_HIREAGREEALL   = "RPC_Friend.HireAgreeAll"   //! 一键同意
	RPC_ADDHIREHERO    = "RPC_Friend.AddHireHero"    //! 增加自己的可租借英雄
	RPC_DELETEHIREHERO = "RPC_Friend.DeleteHireHero" //! 减少自己的可租借英雄

	RPC_QUERYPRIVATEESSAGE = "RPC_Friend.QueryPrivateMessage" //! 获得全部私聊
	RPC_SAVEPRIVATEMESSAGE = "RPC_Friend.SavePrivateMessage"  //! 保存私聊

	RPC_GETFRIEND   = "RPC_Friend.GetFriend"   //! 获取好友,同步
	RPC_GETAPPLY    = "RPC_Friend.GetApply"    //! 获取好友,同步
	RPC_GETBLACK    = "RPC_Friend.GetBlack"    //! 获取好友,同步
	RPC_GETHIRELIST = "RPC_Friend.GetHireList" //! 获取租借列表,同步
	RPC_GETSELFLIST = "RPC_Friend.GetSelfList" //! 获取自己的英雄列表
	RPC_ROBOTACTION = "RPC_Friend.RobotAction" //! 压测机器人动作
	/////////////////////////////

	RPC_SUPPORT_HERO_ADD              = "RPC_Friend.AddSupportHero"    //! 添加支援英雄
	RPC_SUPPORT_HERO_REMOVE           = "RPC_Friend.RemoveSupportHero" //! 取消支援英雄
	RPC_SUPPORT_HERO_USE              = "RPC_Friend.UseHero"           //! 使用支援英雄
	RPC_SUPPORT_HERO_CANCEL_USE       = "RPC_Friend.CancelUseHero"     //! 取消使用支援英雄
	RPC_SUPPORT_HERO_GET_DATA         = "RPC_Friend.GetPlayerData"     //! 获得玩家数据
	RPC_SUPPORT_HERO_CLEAN_DATA       = "RPC_Friend.CleanPlayerData"   //! 清理数据
	RPC_SUPPORT_HERO_GET_CAN_USE_HERO = "RPC_Friend.GetCanUseHero"     //! 获得可用英雄
	RPC_SUPPORT_HERO_UPDATE_HERO      = "RPC_Friend.UpdateHero"        //! 更新
	RPC_SUPPORT_HERO_RENAME           = "RPC_Friend.Rename"            //! 改名
)

//! 角色消息主体
type RPC_Friend struct {
	Client       *rpc.Client
	PlayerLocker *sync.RWMutex //! 数据锁
}

//---------------------------  Player Data  ---------------------
//! RPC 角色数据
type RPC_FriendData_Req struct {
	UId    int64             //! 角色ID
	Online int               //! 在线状态状态
	Data   *JS_PlayerData    //! 基础数据
	Heros  []*JS_PlayerHero  //! 武将列表
	Equips [][]*JS_HeroEquip //! 装备列表，最大5个，可为空
}

type RPC_FriendData_Res struct {
	RetCode int //! 操作结果
}

//! 操作请求消息
type RPC_FriendActionReq struct {
	UId int64 //! 自己的Id
	FId int64 //! 好友Id
}

type RPC_FriendActionFindReq struct {
	UId        int64  //! 自己的Id
	FriendId   int64  //! 好友Id
	FriendName string //! 搜索名字
}

type RPC_FriendActionAgreeReq struct {
	UId   int64 //! 自己的Id
	FId   int64 //! 好友Id
	Agree int   //! 0拒绝  1同意
}

type RPC_FriendActionPowerReq struct {
	UId int64   //! 自己的Id
	FId []int64 //! 好友Id
}

type RPC_FriendActionHireReq struct {
	UId    int64             //! 自己的Id
	Hire   map[int]*HireHero //租借数据
	IsSend int               //是否通知给其他人知道
}

type RPC_FriendHireActionReq struct {
	UId   int64 //! 自己的Id
	FId   int64 //! 好友Id
	KeyId int   //! 英雄KEY
}

//! 操作响应消息
type RPC_FriendActionRes struct {
	RetCode int //! 结果码
	Data    string
}

//! ***********************************************************

func (self *RPC_Friend) Init() bool {
	self.PlayerLocker = new(sync.RWMutex)
	return true
}

//! 寻找好友
//! int 寻找成功
func (self *RPC_Friend) FindFriend(Uid int64, friendUid int64, friendName string) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionFindReq
		req.UId = Uid
		req.FriendId = friendUid
		req.FriendName = friendName

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_FINDFRIEND, req, &res)
		if err == nil {
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

//! 增加好友  Apply
//! int 添加成功
func (self *RPC_Friend) ApplyFriend(uid int64, fid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid
		req.FId = fid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_APPLYFRIEND, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			return &res
		}
	}

	return nil
}

func (self *RPC_Friend) GetFriend(uid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_GETFRIEND, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) AgreeFriend(uid int64, fid int64, agree int) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionAgreeReq
		req.UId = uid
		req.FId = fid
		req.Agree = agree

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_AGREEFRIEND, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) PowerFriend(uid int64, fid []int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionPowerReq
		req.UId = uid
		req.FId = fid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_POWERFRIEND, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) BlackFriend(uid int64, fid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid
		req.FId = fid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_BLACKFRIEND, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) AddHireFriend(uid int64, hire *HireHero, isSend int) *RPC_FriendActionRes {
	if self.Client != nil && hire != nil {
		var req RPC_FriendActionHireReq
		req.UId = uid
		req.Hire = make(map[int]*HireHero, 0)
		req.Hire[hire.HeroKeyId] = hire
		req.IsSend = isSend

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_ADDHIREHERO, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) UpdateHireAll(uid int64, hireList map[int]*HireHero, isSend int) *RPC_FriendActionRes {
	if self.Client != nil && len(hireList) > 0 {
		var req RPC_FriendActionHireReq
		req.UId = uid
		req.Hire = make(map[int]*HireHero, 0)
		req.Hire = hireList
		req.IsSend = isSend

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_ADDHIREHERO, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) DeleteHireFriend(uid int64, keyId int) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendHireActionReq
		req.UId = uid
		req.KeyId = keyId

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_DELETEHIREHERO, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) GetFriendInfo(uid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid
		//req.FId = fid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_APPLYFRIEND, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil

}

func (self *RPC_Friend) DelFriend(uid int64, fid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid
		req.FId = fid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_DELFRIEND, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Friend) BlackOutFriend(uid int64, fid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid
		req.FId = fid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_BLACKOUTFRIEND, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}
	return nil
}

//! 拒绝好友
func (self *RPC_Friend) RefuseFriend(uid int64, fid int64) int {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid
		req.FId = fid

		var res RPC_FriendActionRes

		err := GetMasterMgr().CallEx(self.Client,RPC_REFUSEFRIEND, req, &res)
		if err == nil {
			//! 添加成功
			return RETCODE_OK
		} else {
			//! 添加失败
			return res.RetCode
		}
	}

	return RETCODE_UNKNOWN
}

//! 放弃已经租借到的雇佣
func (self *RPC_Friend) HireApplyLose(uid int64, fid int64, keyId int) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendHireActionReq
		req.UId = uid
		req.FId = fid
		req.KeyId = keyId

		var res RPC_FriendActionRes

		err := GetMasterMgr().CallEx(self.Client,RPC_HIRELOSE, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

//申请雇佣
func (self *RPC_Friend) HireApplyApply(uid int64, fid int64, keyId int) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendHireActionReq
		req.UId = uid
		req.FId = fid
		req.KeyId = keyId

		var res RPC_FriendActionRes

		err := GetMasterMgr().CallEx(self.Client,RPC_HIREAPPLY, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

//取消雇佣
func (self *RPC_Friend) HireApplyCancel(uid int64, fid int64, keyId int) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendHireActionReq
		req.UId = uid
		req.FId = fid
		req.KeyId = keyId

		var res RPC_FriendActionRes

		err := GetMasterMgr().CallEx(self.Client,RPC_HIRECANCEL, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) HireApplyAgree(uid int64, fid int64, keyId int) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendHireActionReq
		req.UId = uid
		req.FId = fid
		req.KeyId = keyId

		var res RPC_FriendActionRes

		err := GetMasterMgr().CallEx(self.Client,RPC_HIREAGREE, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) HireApplyRefuse(uid int64, fid int64, keyId int) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendHireActionReq
		req.UId = uid
		req.FId = fid
		req.KeyId = keyId

		var res RPC_FriendActionRes

		err := GetMasterMgr().CallEx(self.Client,RPC_HIREREFUSE, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) HireApplyAgreeAll(uid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid
		req.FId = 0

		var res RPC_FriendActionRes

		err := GetMasterMgr().CallEx(self.Client,RPC_HIREAGREEALL, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) GetApply(uid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_GETAPPLY, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) GetBlack(uid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_GETBLACK, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) GetHireList(uid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_GETHIRELIST, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

func (self *RPC_Friend) GetSelfList(uid int64) *RPC_FriendActionRes {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid

		var res RPC_FriendActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_GETSELFLIST, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//! 添加失败
			return nil
		}
	}

	return nil
}

/////////////////////////////////////////////////

//! 操作请求
type RPC_SupportHeroAction struct {
	Data string
}

//! 操作响应
type RPC_SupportHeroActionRet struct {
	RetCode int //! 结果码
	Data    string
}

///////////添加支援英雄//////////////
type S2M_SupportHeroAdd struct {
	Uid       int64
	Index     int
	HeroKeyId int
	HeroID    int
	HeroStar  int
	HeroLv    int
	Skin      int
}

type M2S_SupportHeroAdd struct {
}

///////////取消支援英雄//////////////
type S2M_SupportHeroRemove struct {
	Uid       int64
	HeroKeyId int
}

type M2S_SupportHeroRemove struct {
}

///////////使用支援英雄//////////////
type S2M_SupportHeroUse struct {
	Uid       int64
	HeroKeyId int
	Useruid   int64
	Username  string
	Type      int
	Endtime   int64
}

type M2S_SupportHeroUse struct {
}

///////////取消使用支援英雄//////////////
type S2M_SupportHeroCancelUse struct {
	Uid       int64
	HeroKeyId int
	Useruid   int64
}

type M2S_SupportHeroCancelUse struct {
}

///////////获得玩家数据//////////////
type S2M_SupportHeroGetPlayerData struct {
	Uid int64
}

type M2S_SupportHeroGetPlayerData struct {
	Data []*SupportHero
}

///////////清理玩家数据//////////////
type S2M_SupportHeroCleanPlayerData struct {
	Uid int64
}

type M2S_SupportHeroCleanPlayerData struct {
}

///////////获得玩家数据//////////////
type S2M_SupportHeroGetCanUseHero struct {
	Uids   map[int64]int64
	HeroID int
}

type M2S_SupportHeroGetCanUseHero struct {
	Data []*SupportHero
}

///////////使用支援英雄//////////////
type S2M_SupportHeroUpdate struct {
	Uid       int64
	HeroKeyId int
	HeroStar  int
	HeroLv    int
	Skin      int
}

type M2S_SupportHeroUpdate struct {
}

///////////改名//////////////
type S2M_SupportHeroRename struct {
	Uid  int64
	Name string
}

type M2S_SupportHeroRename struct {
}

//操作
func (self *RPC_Friend) SupportHeroAction(action string, data interface{}) *RPC_SupportHeroActionRet {
	if self.Client != nil {
		var req RPC_SupportHeroAction
		req.Data = HF_JtoA(data)

		var ret RPC_SupportHeroActionRet
		GetMasterMgr().CallEx(self.Client,action, req, &ret)
		return &ret
	}

	return nil
}

func (self *RPC_Friend) SavePrivateMessage(player *Player, content string, channel int) *RPC_ChatActionRes {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.Param1 = channel
		req.Param2 = content

		var res RPC_ChatActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_SAVEPRIVATEMESSAGE, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Friend) QueryPrivateMessage(player *Player) *RPC_ChatActionRes {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.Param1 = player.GetUnionId()

		var res RPC_ChatActionRes
		err := GetMasterMgr().CallEx(self.Client,RPC_QUERYPRIVATEESSAGE, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Friend) RobotAction(uid int64) {
	if self.Client != nil {
		var req RPC_FriendActionReq
		req.UId = uid

		var res RPC_FriendActionRes
		GetMasterMgr().CallEx(self.Client,RPC_ROBOTACTION, req, &res)
	}
}
