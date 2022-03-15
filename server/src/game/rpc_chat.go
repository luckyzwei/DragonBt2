/*
@Time : 2020/5/10 11:10
@Author : 96121
@File : proto_player
@Software: GoLand
*/
package game

import (
	"net/rpc"
)

const (
	RPC_ENTERCHAT          = "RPC_Chat.EnterChat"          //! 进入聊天频道
	RPC_EXITCHAT           = "RPC_Chat.ExitChat"           //! 离开聊天频道
	RPC_SENDWORLDMESSAGE   = "RPC_Chat.SendWorldMessage"   //! 发送世界聊天
	RPC_QUERYWORLDMESSAGE  = "RPC_Chat.QueryWorldMessage"  //! 获得世界聊天
	RPC_SENDUNIONMESSAGE   = "RPC_Chat.SendUnionMessage"   //! 发送公会聊天
	RPC_QUERYUNIONMESSAGE  = "RPC_Chat.QueryUnionMessage"  //! 获得公会聊天
	RPC_SENDPRIVATEMESSAGE = "RPC_Chat.SendPrivateMessage" //! 发送私聊
	RPC_GAGPLAYER          = "RPC_Chat.GagPlayer"          //! 禁言用户
)

//! 角色消息主体
type RPC_Chat struct {
	Client *rpc.Client
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
}

func (self *RPC_Chat) Init() bool {
	return true
}

func (self *RPC_Chat) EnterChat(Uid int64) *RPC_ChatEnterChat {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = Uid

		var res RPC_ChatEnterChat
		err := GetMasterMgr().CallEx(self.Client, RPC_ENTERCHAT, req, &res)

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

func (self *RPC_Chat) ExitChat(Uid int64) *RPC_ChatEnterChat {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = Uid

		var res RPC_ChatEnterChat
		err := GetMasterMgr().CallEx(self.Client, RPC_EXITCHAT, req, &res)

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

func (self *RPC_Chat) SendWorldMessage(player *Player, content string, channel int) *RPC_ChatActionRes {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.Param1 = channel
		req.Param2 = content

		var res RPC_ChatActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_SENDWORLDMESSAGE, req, &res)
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

func (self *RPC_Chat) SendUnionMessage(player *Player, content string, channel int) *RPC_ChatActionRes {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.Param1 = channel
		req.Param2 = content

		var res RPC_ChatActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_SENDUNIONMESSAGE, req, &res)
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

func (self *RPC_Chat) SendPrivateMessage(player *Player, content string, channel int) *RPC_ChatActionRes {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.Param1 = channel
		req.Param2 = content

		var res RPC_ChatActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_SENDPRIVATEMESSAGE, req, &res)
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

func (self *RPC_Chat) QueryWorldMessage(player *Player) *RPC_ChatActionRes {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.Param1 = player.GetModule("chat").(*ModChat).GetWorldChannel()

		var res RPC_ChatActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_QUERYWORLDMESSAGE, req, &res)
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

func (self *RPC_Chat) QueryUnionMessage(player *Player) *RPC_ChatActionRes {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.Param1 = player.GetUnionId()

		var res RPC_ChatActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_QUERYUNIONMESSAGE, req, &res)
		if err == nil {
			//! 添加成功
			return &res
		} else {
			//print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Chat) GagPlayer(player *Player, union int) *RPC_ChatActionRes {
	if self.Client != nil {
		var req RPC_ChatActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.Param1 = union

		var res RPC_ChatActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_GAGPLAYER, req, &res)
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
