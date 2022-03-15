package chat

import (
	"errors"
	"master/core"
	"master/utils"
	"time"
)

//! 聊天消息主体
type RPC_Chat struct {
}

//! 聊天事件请求
type RPC_ChatActionReq struct {
	Uid    int64  //! 角色Id
	Param1 int    //! 请求参数1,渠道Id，公会Id等
	Param2 string //! 请求参数2 消息内容，可为空
}

//! 聊天事件返回
type RPC_ChatEnterChat struct {
	RetCode int //! 返回结果
	Channel int //! 频道ID
}

//! 聊天事件返回
type RPC_ChatActionRes struct {
	RetCode int            //! 返回结果
	MsgList []*ChatMessage //! 消息列表
}

//! 聊天记录，保留信息快照，根据显示添加
type ChatMessage struct {
	MsgId      int    `json:"msgid"`      //! 聊天消息Id
	Uid        int64  `json:"uid"`        //! 发布角色ID
	Uname      string `json:"uname"`      //! 角色名字
	Level      int    `json:"level"`      //! 角色等级
	IconId     int    `json:"iconid"`     //! 头像
	Portrait   int    `json:"portrait"`   //! 头像框
	Content    string `json:"content"`    //! 聊天内容
	SendTime   int    `json:"sendtime"`   //! 发布时间
	ToUid      int64  `json:"touid"`      //! 目标
	ToUname    string `json:"touname"`    //! 角色名字
	ToLevel    int    `json:"tolevel"`    //! 角色等级
	ToIconId   int    `json:"toiconid"`   //! 头像
	ToPortrait int    `json:"toportrait"` //! 头像框
	IsGag      int    `json:"isgag"`      //! 禁言标记
}

//! 进入聊天，同时进入世界聊天和军团
func (self *RPC_Chat) EnterChat(req RPC_ChatActionReq, res *RPC_ChatEnterChat) error {
	p := core.GetPlayerMgr().GetCorePlayer(req.Uid, true)
	if p == nil {
		//! 角色数据不存在
		return nil
	}

	//! 世界频道加入
	channel := GetChatMgr().GetWorldChannel(req.Param1)
	if channel != nil {
		channel.AddPlayer(p.GetUid(), p.GetUname(), p.GetLevel(), p.GetIconId(), p.GetPortrait(), p.GetServerId())
		res.Channel = channel.ChannelId
		GetChatMgr().SetPlayerChannel(CHAT_WORLD, req.Uid, channel)
	}

	if p.GetUnionId() > 0 {
		unionCh := GetChatMgr().GetUnionChannel(p.GetUnionId())
		if unionCh != nil {
			unionCh.AddPlayer(p.GetUid(), p.GetUname(), p.GetLevel(), p.GetIconId(), p.GetPortrait(), p.GetServerId())
		}
	}

	return nil
}

//! 离开聊天
func (self *RPC_Chat) ExitChat(req *RPC_ChatActionReq, res *RPC_ChatActionRes) error {
	p := core.GetPlayerMgr().GetCorePlayer(req.Uid, true)
	if p == nil {
		return nil
	}

	//! 世界频道
	channel := GetChatMgr().GetWorldChannel(req.Param1)
	if channel != nil {
		channel.DelPlayer(p.GetUid())
	}

	if p.GetUnionId() > 0 {
		unionCh := GetChatMgr().GetUnionChannel(p.GetUnionId())
		if unionCh != nil {
			unionCh.DelPlayer(p.GetUid())
		}
	}
	return nil
}

//! 获取聊天消息（历史）
//! msgId = 0 ，请求最新消息
//! msgId > 0 ， 请求历史信息
func (self *RPC_Chat) QueryWorldMessage(req *RPC_ChatActionReq, res *RPC_ChatActionRes) error {
	channel := GetChatMgr().GetPlayerChannel(req.Uid, CHAT_WORLD)
	if channel != nil {
		res.RetCode = 0
		msgList := channel.QueryNewMessage(req.Param1)
		if len(msgList) > 0 {
			//! 存在消息
			res.MsgList = msgList
		}
	} else {
		res.RetCode = 1
		return errors.New("Don't Have World Channel.")
	}

	return nil
}

func (self *RPC_Chat) QueryUnionMessage(req *RPC_ChatActionReq, res *RPC_ChatActionRes) error {
	channel := GetChatMgr().GetPlayerChannel(req.Uid, CHAT_PARTY)
	if channel != nil {
		res.RetCode = 0
		msgList := channel.QueryNewMessage(req.Param1)
		if len(msgList) > 0 {
			//! 存在消息
			res.MsgList = msgList
		}
	} else {
		res.RetCode = 1
		return errors.New("Don't Have Union Channel.")
	}

	return nil
}

//! 发送世界聊天
func (self *RPC_Chat) SendWorldMessage(req *RPC_ChatActionReq, res *RPC_ChatActionRes) error {
	channel := GetChatMgr().GetWorldChannel(req.Param1)
	if channel == nil {
		return nil
	}
	channel.CheckPlayer(req.Uid)

	//! 发送聊天内容
	channel.SendMessage(req, res, core.CHAT_NEW_WORLD_MESSAGE)
	return nil
}

//! 发送公会消息
func (self *RPC_Chat) SendUnionMessage(req *RPC_ChatActionReq, res *RPC_ChatActionRes) error {
	channel := GetChatMgr().GetUnionChannel(req.Param1)
	if channel == nil {
		return nil
	}

	channel.CheckPlayer(req.Uid)

	//! 发送聊天内容
	channel.SendMessage(req, res, core.CHAT_NEW_UNION_MESSAGE)
	return nil
}

//! 发送私聊
func (self *RPC_Chat) SendPrivateMessage(req *RPC_ChatActionReq, res *RPC_ChatActionRes) error {

	pl := core.GetPlayerMgr().GetCorePlayer(req.Uid, true)
	if pl == nil {
		return nil
	}

	plTo := core.GetPlayerMgr().GetCorePlayer(int64(req.Param1), true)
	if plTo == nil {
		return nil
	}
	msg := make([]*ChatMessage, 0)
	msgInfo := &ChatMessage{
		MsgId:      0,
		Uid:        pl.GetUid(),
		Uname:      pl.GetUname(),
		Level:      pl.GetLevel(),
		IconId:     pl.GetIconId(),
		Portrait:   pl.GetPortrait(),
		Content:    req.Param2,
		ToUid:      plTo.GetUid(),
		ToUname:    plTo.GetUname(),
		ToLevel:    plTo.GetLevel(),
		ToIconId:   plTo.GetIconId(),
		ToPortrait: plTo.GetPortrait(),
		SendTime:   int(time.Now().Unix()),
	}
	msg = append(msg, msgInfo)

	core.GetCenterApp().AddEvent(pl.GetServerId(), core.CHAT_NEW_PRIVATE_MESSAGE, pl.GetUid(), 0,
		0, utils.HF_JtoA(msg))

	core.GetCenterApp().AddEvent(plTo.GetServerId(), core.CHAT_NEW_PRIVATE_MESSAGE, plTo.GetUid(), 0,
		0, utils.HF_JtoA(msg))
	return nil
}

func (self *RPC_Chat) GagPlayer(req *RPC_ChatActionReq, res *RPC_ChatActionRes) error {
	utils.LogDebug("GagPlayer :", req.Uid)
	//先处理世界
	channelWorld := GetChatMgr().GetWorldChannel(0)
	if channelWorld != nil {
		channelWorld.DeleteMessageByUid(req.Uid)
	}

	//处理工会
	channelUnion := GetChatMgr().GetUnionChannel(req.Param1)
	if channelUnion == nil {
		return nil
	}
	//! 删除聊天记录
	channelUnion.DeleteMessageByUid(req.Uid)
	return nil
}
