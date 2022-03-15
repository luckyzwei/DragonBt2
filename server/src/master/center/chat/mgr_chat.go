package chat

import "sync"

const ()

//! 聊天管理器
//! 聊天频道只会增加，不会删除，回收后，只限制不加入
//! 自定义聊天频道
type ChatMgr struct {
	WorldMap       map[int]*ModChannel   //! 世界频道，Key = 频道Id
	PlayerWorldMap map[int64]*ModChannel //! 用户所在频道
	UnionMap       map[int]*ModChannel   //! 工会聊天频道
	PlayerUnionMap map[int64]*ModChannel //! 用户所在频道
	DataLocker     *sync.RWMutex         //! 数据锁，频道锁
	PlayerLocker   *sync.RWMutex         //! 数据锁，用户频道锁
	WorldIdMax     int                   //! 世界频道Id利用
}

//! 全局唯一
var s_chatmgr *ChatMgr = nil

func GetChatMgr() *ChatMgr {
	if s_chatmgr == nil {
		s_chatmgr = new(ChatMgr)
		s_chatmgr.DataLocker = new(sync.RWMutex)
		s_chatmgr.PlayerLocker = new(sync.RWMutex)
		s_chatmgr.WorldIdMax = 0
		s_chatmgr.WorldMap = make(map[int]*ModChannel)
		s_chatmgr.PlayerWorldMap = make(map[int64]*ModChannel)

		s_chatmgr.UnionMap = make(map[int]*ModChannel)
		s_chatmgr.PlayerUnionMap = make(map[int64]*ModChannel)
	}

	return s_chatmgr
}

func (self *ChatMgr) InitService() {
	//! 默认初始化10个世界频道

}

//! 增加一个世界聊天频道
func (self *ChatMgr) AddWorldChannel() int {
	channel := new(ModChannel)
	channel.InitChannel(CHAT_WORLD)
	//! 只增加不销毁
	self.WorldIdMax += 1
	channel.ChannelId = self.WorldIdMax

	self.DataLocker.Lock()
	self.WorldMap[channel.ChannelId] = channel
	self.DataLocker.Unlock()

	return channel.ChannelId
}

//! 世界聊天频道
func (self *ChatMgr) GetWorldCount() int {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	return len(self.WorldMap)
}

//! 所有在线用户人数
func (self *ChatMgr) GetPlayerCount() int {
	self.PlayerLocker.RLock()
	defer self.PlayerLocker.RUnlock()

	return len(self.PlayerWorldMap)
}

//! 获取指定世界聊天频道
func (self *ChatMgr) GetWorldChannel(channelId int) *ModChannel {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	//if channel, ok := self.WorldMap[channelId]; ok {
	//	return channel
	//}

	if channel, ok := self.WorldMap[1]; ok {
		return channel
	}

	return nil
}

//! 增加公会频道，加入公会自动加入
//! 系统启动时，自动生成所有的公会频道
func (self *ChatMgr) AddUnionChannel(UnionId int) *ModChannel {
	channel := new(ModChannel)
	channel.InitChannel(CHAT_PARTY)
	channel.ChannelId = UnionId

	self.DataLocker.Lock()
	self.UnionMap[UnionId] = channel
	self.DataLocker.Unlock()
	return channel
}

//! 设置角色频道信息
func (self *ChatMgr) SetPlayerChannel(chType int, uid int64, ch *ModChannel) {
	self.PlayerLocker.Lock()
	defer self.PlayerLocker.Unlock()

	if chType == CHAT_PARTY {
		if ch == nil {
			delete(self.PlayerUnionMap, uid)
		} else {
			self.PlayerUnionMap[uid] = ch
		}
	} else {
		if ch == nil {
			delete(self.PlayerWorldMap, uid)
		} else {
			self.PlayerWorldMap[uid] = ch
		}
	}

}

func (self *ChatMgr) GetPlayerChannel(uid int64, chType int) *ModChannel {
	self.PlayerLocker.RLock()
	defer self.PlayerLocker.RUnlock()

	if chType == CHAT_PARTY {
		if ch, ok := self.PlayerUnionMap[uid]; ok {
			return ch
		} else {
			return nil
		}
	} else if chType == CHAT_WORLD {
		if ch, ok := self.PlayerWorldMap[uid]; ok {
			return ch
		} else {
			return nil
		}
	}

	return nil
}

//! 获取公会聊天频道
func (self *ChatMgr) GetUnionChannel(unionId int) *ModChannel {
	self.DataLocker.RLock()
	if channel, ok := self.UnionMap[unionId]; ok {
		self.DataLocker.RUnlock()
		return channel
	}
	self.DataLocker.RUnlock()

	return GetChatMgr().AddUnionChannel(unionId)
}
