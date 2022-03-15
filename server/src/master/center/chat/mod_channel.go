package chat

import (
	"game"
	"master/core"
	"master/utils"
	"sync"
	"time"
)

//! 跨服聊天信息
const (
	CHAT_WORLD   = 1 // 世界(阵营)
	CHAT_PARTY   = 2 // 帮派(职位)
	CHAT_PRIVATE = 4 // 私聊(阵营)

	MAX_MESSAGE_NUM = 1000 //! 暂时保留最长记录 --> 1000条,环形结构
	PER_MESSAGE_NUM = 10   //! 每次获取聊天记录 --> 10条
	MAX_IDLE_TIME   = 1800 //! 30分钟不更新则不再同步消息
	MAX_PLAYER_NUM  = 500  //! 一个聊天频道最高500人
)

type ModChannel struct {
	ChannelId   int                      //! 渠道Id，世界频道=频道Id，公会频道=公会Id
	ChannelType int                      //! 频道类型
	DataLocker  *sync.RWMutex            //! 数据锁
	PlayerList  map[int64]*ChannelPlayer //! 角色列表

	MaxMsgId int            //! 当前占用的最大消息Id
	MsgArr   []*ChatMessage //! 聊天消息队列-环形队列
}

//! 渠道用户
type ChannelPlayer struct {
	Uid        int64  //! 用户Id
	Uname      string //! 用户名字
	Level      int    //! 等级
	IconId     int    //! 头像
	Portrait   int    //! 头像框
	ServerId   int    //! 服务器ID
	LastUpdate int64  //! 上次更新时间，超出15分钟没有刷新，视为掉线，忽略处理
}

//! 初始化频道
func (self *ModChannel) InitChannel(channelType int) {
	self.MaxMsgId = 0
	self.ChannelType = channelType

	self.DataLocker = new(sync.RWMutex)
	self.PlayerList = make(map[int64]*ChannelPlayer)
	self.MsgArr = make([]*ChatMessage, MAX_MESSAGE_NUM)
}

//! 获得频道角色
func (self *ModChannel) GetPlayer(uid int64) *ChannelPlayer {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	if uid == 0 {
		return nil
	}

	if p, ok := self.PlayerList[uid]; ok {
		return p
	}

	return nil
}

func (self *ModChannel) GetPlayerCount() int {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	return len(self.PlayerList)
}

//! 加入频道
func (self *ModChannel) AddPlayer(uid int64, uname string, level, iconId, portrait, serverid int) {
	p := self.GetPlayer(uid)
	if p == nil {
		utils.LogDebug("加入频道：", self.ChannelType, self.ChannelId, uid, uname)
		//! 不存在则添加
		p = new(ChannelPlayer)
		p.Uid = uid
		p.Uname = uname
		p.Level = level
		p.IconId = iconId
		p.Portrait = portrait
		p.ServerId = serverid

		self.DataLocker.Lock()
		self.PlayerList[uid] = p
		self.DataLocker.Unlock()

		//! 设置聊天频道
		GetChatMgr().SetPlayerChannel(self.ChannelType, uid, self)
	} else {
		//! 存在的话，则更新数据
		//! uid 全局唯一
		p.Uname = uname
		p.Level = level
		p.IconId = iconId
		p.Portrait = portrait
		p.ServerId = serverid

		ch := GetChatMgr().GetPlayerChannel(uid, self.ChannelType)
		if ch != nil {
			ch.DelPlayer(uid)
		}

		//! 设置聊天频道
		GetChatMgr().SetPlayerChannel(self.ChannelType, uid, self)
	}
}

//! 离开频道
func (self *ModChannel) DelPlayer(uid int64) {
	p := self.GetPlayer(uid)
	if p != nil {
		utils.LogDebug("离开频道：", self.ChannelType, self.ChannelId, uid, p.Uname)

		//! 存在角色，则删除
		self.DataLocker.Lock()
		delete(self.PlayerList, uid)
		self.DataLocker.Unlock()

		//! 删除聊天频道信息
		GetChatMgr().SetPlayerChannel(self.ChannelType, uid, nil)
	}
}

//! 发送聊天消息
func (self *ModChannel) SendMessage(req *RPC_ChatActionReq, res *RPC_ChatActionRes, event int) bool {
	p := self.GetPlayer(req.Uid)
	if p != nil {
		//! 发送消息
		self.DataLocker.Lock()
		self.MaxMsgId++
		msg := &ChatMessage{
			MsgId:    self.MaxMsgId,
			Uid:      req.Uid,
			Uname:    p.Uname,
			Level:    p.Level,
			IconId:   p.IconId,
			Portrait: p.Portrait,
			Content:  req.Param2,
			SendTime: int(time.Now().Unix()),
		}
		self.DataLocker.Unlock()
		//! 简易的环形结构，回收利用，读取时注意
		self.SetMessage(msg)
		res.MsgList = append(res.MsgList, msg)

		self.BroadcastMessage(self.MaxMsgId, event)

		return true
	}

	return false
}

func (self *ModChannel) SetMessage(msg *ChatMessage) {
	arrIdx := msg.MsgId % MAX_MESSAGE_NUM
	self.MsgArr[arrIdx] = msg
}

func (self *ModChannel) GetMessage(msgId int) []*ChatMessage {
	rel := make([]*ChatMessage, 0)
	arrIdx := msgId % MAX_MESSAGE_NUM
	rel = append(rel, self.MsgArr[arrIdx])
	return rel
}

//! 刷新角色更新时间
func (self *ModChannel) UpatePlayer(uid int64) {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	p := self.GetPlayer(uid)
	if p != nil {
		p.LastUpdate = time.Now().Unix()
	}
}

//! 获取最新的聊天消息
func (self *ModChannel) QueryNewMessage(LastReqId int) []*ChatMessage {
	resArr := make([]*ChatMessage, 0)

	start := 1
	if self.MaxMsgId > 50 {
		start = self.MaxMsgId - 50
	}

	for i := start; i <= self.MaxMsgId; i++ {
		msg := self.GetMessage(i)
		if len(msg) > 0 && msg[0].IsGag == game.LOGIC_FALSE {
			pl := core.GetPlayerMgr().GetCorePlayer(msg[0].Uid, false)
			if pl != nil {
				msg[0].Uname = pl.GetUname()
				msg[0].Level = pl.GetLevel()
				msg[0].IconId = pl.GetIconId()
				msg[0].Portrait = pl.GetPortrait()
			}
			resArr = append(resArr, msg...)
		}
	}

	/*
		//! 如果是最新的数据，则不更新
			if LastReqId >= self.MaxMsgId {
				return resArr
			}

			msgNum := PER_MESSAGE_NUM
			if LastReqId+PER_MESSAGE_NUM > self.MaxMsgId {
				msgNum = self.MaxMsgId - LastReqId
			}

			for i := 0; i < msgNum; i++ {
				msg := self.GetMessage(LastReqId + i)
				if msg != nil {
					resArr = append(resArr, msg...)
				}
			}
	*/

	return resArr
}

//! 取最新的10条聊天记录
func (self *ModChannel) QueryLastMessage() []*ChatMessage {
	return self.QueryNewMessage(self.MaxMsgId)
}

//! 广播消息，发送聊天事件
func (self *ModChannel) BroadcastMessage(msgId int, event int) bool {
	if msgId <= 0 {
		return false
	}

	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	msg := self.GetMessage(msgId)
	if msg == nil {
		return false
	}
	for i := 0; i < len(msg); i++ {
		pl := core.GetPlayerMgr().GetCorePlayer(msg[i].Uid, false)
		if pl != nil {
			msg[i].Uname = pl.GetUname()
			msg[i].Level = pl.GetLevel()
			msg[i].IconId = pl.GetIconId()
			msg[i].Portrait = pl.GetPortrait()
		}
	}
	utils.LogDebug("BroadcastMessage:", utils.HF_JtoA(msg))
	//tNowTime := time.Now().Unix()
	for uid, p := range self.PlayerList {
		//! 超过30秒未更新，则不同步信息
		//if p.LastUpdate < tNowTime-MAX_IDLE_TIME {
		//	continue
		//}
		core.GetCenterApp().AddEvent(p.ServerId, event, uid, 0,
			0, utils.HF_JtoA(msg))
	}

	return true
}

//处理因为中心服断线导致的聊天频道丢失问题
func (self *ModChannel) CheckPlayer(uid int64) {
	nowP := self.GetPlayer(uid)
	if nowP == nil {
		p := core.GetPlayerMgr().GetCorePlayer(uid, true)
		if p != nil {
			self.AddPlayer(p.GetUid(), p.GetUname(), p.GetLevel(), p.GetIconId(), p.GetPortrait(), p.GetServerId())
		}
	}
}

//! 删除禁言聊天
func (self *ModChannel) DeleteMessageByUid(uid int64) bool {
	p := self.GetPlayer(uid)
	if p != nil {
		self.DataLocker.Lock()
		defer self.DataLocker.Unlock()

		start := 1
		if self.MaxMsgId > 1000 {
			start = self.MaxMsgId - 1000
		}

		isNeed := false
		for i := start; i <= self.MaxMsgId; i++ {
			msg := self.GetMessage(i)
			if len(msg) > 0 && msg[0].Uid == uid {
				msg[0].IsGag = game.LOGIC_TRUE
				isNeed = true
			}
		}

		//广播封禁信息
		if isNeed {
			for uidPlayer, p := range self.PlayerList {
				//! 超过30秒未更新，则不同步信息
				//if p.LastUpdate < tNowTime-MAX_IDLE_TIME {
				//	continue
				//}
				core.GetCenterApp().AddEvent(p.ServerId, core.CHAT_GAP_PLAYER, uidPlayer, uid,
					self.ChannelType, "")
			}
		}

		return true
	}

	return false
}
