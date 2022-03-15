/*
@Time : 2020/5/10 11:10
@Author : 96121
@File : proto_player
@Software: GoLand
*/
package player

import (
	"game"
	"master/center/chat"
	"master/utils"
)

//! 错误码定义
const (
	RETCODE_OK                  = 0  //! 没有错误
	RETCODE_PLAYER_NOT_EXIST    = -1 //! 角色不存在
	RETCODE_FRIEND_NOT_EXIST    = -2 //! 好友不存在
	RETCODE_SELF_FRIEND_FULL    = -3 //! 自己好友已满
	RETCODE_TARGET_FRIEND_FULL  = -4 //! 目标好友已满
	RETCODE_FRIEND_APPLY_ERROR  = -5 //! 好友申请不存在
	RETCODE_FRIEND_APPLY_BLACK  = -6 //! 黑名单
	RETCODE_FRIEND_HIRE_ALREADY = -7 //! 佣兵租用中
	RETCODE_FRIEND_HIRE_MAX     = -8 //! 目标佣兵已满
)

type HireHero struct {
	HeroKeyId           int               `json:"herokeyid"`
	HeroId              int               `json:"heroid"`
	HeroQuality         int               `json:"heroquality"`         //英雄品质  (不变的属性)
	HeroArtifactId      int               `json:"heroartifactid"`      //英雄神器等级 (不变的属性)
	HeroArtifactLv      int               `json:"heroartifactlv"`      //英雄神器等级 (不变的属性)
	HeroExclusiveLv     int               `json:"heroexclusivelv"`     //英雄专属等级 (不变的属性)
	HeroExclusiveUnLock int               `json:"heroexclusiveunlock"` //英雄专属解锁 (不变的属性)
	Talent              *game.StageTalent `json:"talent"`              //天赋
	OwnPlayer           *HirePlayerBase   `json:"ownplayer"`           //拥有者
	ApplyPlayer         []*HirePlayerBase `json:"applyplayer"`         //申请列表
	HirePlayer          *HirePlayerBase   `json:"hireplayer"`          //
	ReSetTime           int64             `json:"resettime"`           //
}

type HirePlayerBase struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"uname"`
	Face     int    `json:"face"`
	IconId   int    `json:"iconid"`
	Portrait int    `json:"portrait"`
}

//! 角色消息主体
type RPC_Friend struct {
}

type RPC_FriendData_Res struct {
	RetCode int //! 操作结果
}

//! 操作请求
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

//! 操作响应
type RPC_FriendActionRes struct {
	RetCode int //! 结果码
	Data    string
}

//! 查找好友
func (self *RPC_Friend) FindFriend(req RPC_FriendActionFindReq, res *RPC_FriendActionRes) error {
	utils.LogDebug("Find Friend :", req.FriendId, req.FriendName)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		res.Data = my.DataFriend.FindFriend(my.GetUId(), req.FriendId, req.FriendName)
	}
	return nil
}

//! 增加好友
func (self *RPC_Friend) ApplyFriend(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	utils.LogDebug("Add Friend :", req.UId, req.FId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.AddFriend(req.FId, res)
	}
	return nil
}

//! 删除好友
func (self *RPC_Friend) DelFriend(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("Del Friend :", req.UId, req.FId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.DelFriend(req.FId, res)
	}
	return nil
}

//! 移除黑名单
func (self *RPC_Friend) BlackOutFriend(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("BlackOutFriend :", req.UId, req.FId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.BlackOutFriend(req.FId, res)
	}
	return nil
}

//! 同意好友
func (self *RPC_Friend) AgreeFriend(req RPC_FriendActionAgreeReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("Agree Friend :", req.UId, req.FId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		if req.Agree == game.LOGIC_FALSE {
			my.DataFriend.RefuseFriend(req.FId, res)
		} else if req.Agree == game.LOGIC_TRUE {
			my.DataFriend.AgreeFriend(req.FId, res)
		}
	}
	return nil
}

func (self *RPC_Friend) PowerFriend(req RPC_FriendActionPowerReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("Power Friend :", req.UId, req.FId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.PowerFriend(req.FId, res)
	}
	return nil
}

func (self *RPC_Friend) BlackFriend(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("Black Friend :", req.UId, req.FId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.BlackFriend(req.FId, res)
	}
	return nil
}

func (self *RPC_Friend) AddHireHero(req RPC_FriendActionHireReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("AddHireHero :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.AddHireHero(req.Hire, req.IsSend, res)
	}
	return nil
}

func (self *RPC_Friend) DeleteHireHero(req RPC_FriendHireActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("DeleteHireHero :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.DeleteHireHero(req.KeyId, res)
	}
	return nil
}

func (self *RPC_Friend) HireLose(req RPC_FriendHireActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("HireLose :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.HireLose(req.FId, req.KeyId, res)
	}
	return nil
}

func (self *RPC_Friend) HireApply(req RPC_FriendHireActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("HireApply :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.HireApply(req.FId, req.KeyId, res)
	}
	return nil
}

func (self *RPC_Friend) HireCancel(req RPC_FriendHireActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("HireCancel :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.HireCancel(req.FId, req.KeyId, res)
	}
	return nil
}

func (self *RPC_Friend) HireAgree(req RPC_FriendHireActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("HireAgree :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.HireAgree(req.FId, req.KeyId, res)
	}
	return nil
}

func (self *RPC_Friend) HireAgreeAll(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("HireAgreeAll :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.HireAgreeAll(res)
	}
	return nil
}

func (self *RPC_Friend) HireRefuse(req RPC_FriendHireActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("HireRefuse :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.HireRefuse(req.FId, req.KeyId, res)
	}
	return nil
}

func (self *RPC_Friend) GetFriend(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("GetFriend :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.GetFriend(res)
	}
	return nil
}

func (self *RPC_Friend) GetApply(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("GetApply :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.GetApply(res)
	}
	return nil
}

func (self *RPC_Friend) GetBlack(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("GetBlack :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.GetBlack(res)
	}
	return nil
}

func (self *RPC_Friend) GetHireList(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("GetHireList :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.GetHireList(res)
	}
	return nil
}

func (self *RPC_Friend) GetSelfList(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	//utils.LogDebug("GetSelfList :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.GetSelfList(res)
	}
	return nil
}

func (self *RPC_Friend) SavePrivateMessage(req *chat.RPC_ChatActionReq, res *chat.RPC_ChatActionRes) error {
	//utils.LogDebug("SavePrivateMessage :", req.UId)
	my := GetPlayerMgr().GetPlayer(req.Uid, true)
	if my != nil {
		my.DataFriend.SavePrivateMessage(int64(req.Param1), req.Param2)
	}
	return nil
}

func (self *RPC_Friend) QueryPrivateMessage(req *chat.RPC_ChatActionReq, res *chat.RPC_ChatActionRes) error {

	my := GetPlayerMgr().GetPlayer(req.Uid, true)
	if my != nil {
		my.DataFriend.QueryPrivateMessage(res)
	}
	return nil
}

func (self *RPC_Friend) RobotAction(req RPC_FriendActionReq, res *RPC_FriendActionRes) error {
	my := GetPlayerMgr().GetPlayer(req.UId, true)
	if my != nil {
		my.DataFriend.RobotAction(res)
	}
	return nil
}
