/*
@Time : 2020/4/22 10:10
@Author : 96121
@File : data_friend
@Software: GoLand
*/
package player

import (
	"encoding/json"
	"fmt"
	"game"
	"master/center/chat"
	"master/center/union"
	"master/core"
	"master/db"
	"master/utils"
	"sync"
	"time"
)

const (
	FRIEND_TABLE_NAME   = "tbl_userfriend"
	FRIEND_TABLE_SQL    = "select * from `%s` where uid = %d"
	MAX_FRIEND_NUM      = 30
	MAX_FRIEND_HIRE_NUM = 3
)

//! 好友数据
type JS_Friend struct {
	UId        int64  `json:"uid"`   //! 好友Uid
	Name       string `json:"name"`  //! 好友昵称
	Icon       int    `json:"icon"`  //! 头像
	Level      int    `json:"level"` //! 好友等级
	Vip        int    `json:"vip"`   //! VIP
	Fight      int64  `json:"fight"` //! 战力
	Online     int    `json:"online"`
	Portrait   int    `json:"portrait"` //! 头像边框
	Stage      int    `json:"stage"`
	Server     int    `json:"server"`
	LastUpTime int64  `json:"lastuptime"`
}

//! 好友佣兵-可雇佣
type JS_FriendHero struct {
}

//! 好友数据库结构
type SQL_Friend struct {
	UId          int64  `json:"uid"`          //! 好友Uid
	Friends      string `json:"friends"`      //! 好友列表
	Applys       string `json:"applys"`       //! 申请列表（自己申请）
	Applieds     string `json:"applieds"`     //! 被申请列表（他人申请
	Black        string `json:"black"`        //! 黑名单
	HireHeroInfo string `json:"hireheroinfo"` //! 自己的英雄租借租借列表
	SelfHire     string `json:"selfhire"`     //! 自己的英雄租借租借列表
	SupportHero  string `json:"supporthero"`  //! 图书馆支援英雄列表
	ChatSendInfo string `json:"chatsendinfo"` //! 自己发送到私聊
	ChatGetInfo  string `json:"chatgetinfo"`  //! 自己收到的私聊

	friends       map[int64]*JS_Friend //! 好友数据
	applys        map[int64]*JS_Friend //! 申请数据
	applieds      map[int64]*JS_Friend //! 被申请数据
	black         map[int64]*JS_Friend //! 被申请数据
	hireHeroInfo  map[int]*HireHero
	selfHire      []*HireHero
	supportHero   []*SupportHero
	chatSendInfo  []*chat.ChatMessage
	chatGetInfo   []*chat.ChatMessage
	db.DataUpdate //! 数据库操作结构
}

//! 解出数据
func (self *SQL_Friend) Decode() {
	json.Unmarshal([]byte(self.Friends), &self.friends)
	json.Unmarshal([]byte(self.Applys), &self.applys)
	json.Unmarshal([]byte(self.Applieds), &self.applieds)
	json.Unmarshal([]byte(self.Black), &self.black)
	json.Unmarshal([]byte(self.HireHeroInfo), &self.hireHeroInfo)
	json.Unmarshal([]byte(self.SelfHire), &self.selfHire)
	json.Unmarshal([]byte(self.SupportHero), &self.supportHero)
	json.Unmarshal([]byte(self.ChatSendInfo), &self.chatSendInfo)
	json.Unmarshal([]byte(self.ChatGetInfo), &self.chatGetInfo)
}

//! 加密数据
func (self *SQL_Friend) Encode() {
	self.Friends = utils.HF_JtoA(&self.friends)
	self.Applys = utils.HF_JtoA(&self.applys)
	self.Applieds = utils.HF_JtoA(&self.applieds)
	self.Black = utils.HF_JtoA(&self.black)
	self.HireHeroInfo = utils.HF_JtoA(&self.hireHeroInfo)
	self.SelfHire = utils.HF_JtoA(&self.selfHire)
	self.SupportHero = utils.HF_JtoA(&self.supportHero)
	self.ChatSendInfo = utils.HF_JtoA(&self.chatSendInfo)
	self.ChatGetInfo = utils.HF_JtoA(&self.chatGetInfo)
}

//! 好友模块
type ModFriend struct {
	Data   SQL_Friend    //! 好友数据
	player *Player       //! 角色数据
	Locker *sync.RWMutex //! 数据锁
}

func (self *ModFriend) onGetData(player *Player) {
	self.player = player
	self.Locker = new(sync.RWMutex)
	self.Locker.Lock()
	defer self.Locker.Unlock()

	sql := fmt.Sprintf(FRIEND_TABLE_SQL, FRIEND_TABLE_NAME, player.GetUId())
	ret := db.GetDBMgr().DBUser.GetOneData(sql, &self.Data, FRIEND_TABLE_NAME, self.player.GetUId())

	if ret == true { // 创建新号
		if self.Data.UId <= 0 {
			self.Data.UId = self.player.GetUId()

			self.Data.friends = make(map[int64]*JS_Friend)
			self.Data.applys = make(map[int64]*JS_Friend)
			self.Data.applieds = make(map[int64]*JS_Friend)
			self.Data.black = make(map[int64]*JS_Friend)
			self.Data.hireHeroInfo = make(map[int]*HireHero)
			self.Data.selfHire = make([]*HireHero, 0)
			self.Data.supportHero = make([]*SupportHero, 0)
			self.Data.chatSendInfo = make([]*chat.ChatMessage, 0)
			self.Data.chatGetInfo = make([]*chat.ChatMessage, 0)

			self.Data.Encode()
			db.InsertTable(FRIEND_TABLE_NAME, &self.Data, 0, true)
		} else {
			self.Data.Decode()
		}
	}

	if self.Data.friends == nil {
		self.Data.friends = make(map[int64]*JS_Friend)
	}
	if self.Data.applys == nil {
		self.Data.applys = make(map[int64]*JS_Friend)
	}
	if self.Data.applieds == nil {
		self.Data.applieds = make(map[int64]*JS_Friend)
	}
	if self.Data.black == nil {
		self.Data.black = make(map[int64]*JS_Friend)
	}
	if self.Data.hireHeroInfo == nil {
		self.Data.hireHeroInfo = make(map[int]*HireHero)
	}

	self.Data.Init(FRIEND_TABLE_NAME, &self.Data, true)
}

func (self *ModFriend) onSave() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	self.Data.Encode()
	self.Data.Update(true, false)
}

func (self *ModFriend) GetFriendNum() int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	return len(self.Data.friends)
}

func (self *ModFriend) UpdateData() {
	for _, v := range self.Data.hireHeroInfo {
		if v.OwnPlayer == nil {
			continue
		}
		v.OwnPlayer.Name = self.player.Data.UName
		v.OwnPlayer.IconId = self.player.Data.data.IconId
		v.OwnPlayer.Portrait = self.player.Data.data.Portrait
	}
}

//! 寻找好友
func (self *ModFriend) FindFriend(Uid int64, friendUid int64, friendName string) string {

	rel := make([]*JS_PlayerData, 0)
	if friendUid > 0 {
		friendPlayer := GetPlayerMgr().GetPlayer(friendUid, true)
		if friendPlayer != nil && friendPlayer.Data.UId != Uid {
			rel = append(rel, friendPlayer.Data.data)
		}
	} else if friendName != "" {
		var userbase SQL_PlayerData
		sql := fmt.Sprintf("select * from `%s` where `uname` like '%%%s%%' limit 10", PLAYER_TABLE_NAME, friendName)
		res := db.GetDBMgr().DBUser.GetAllData(sql, &userbase)

		if len(res) > 0 {
			for i := 0; i < len(res); i++ {
				data := res[i].(*SQL_PlayerData)
				if data.UId == Uid {
					continue
				}
				data.Decode()
				rel = append(rel, data.data)
			}
		}
	}
	return utils.HF_JtoA(rel)
}

//! 增加好友-申请
func (self *ModFriend) AddFriend(friendId int64, res *RPC_FriendActionRes) {
	friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
	if friendPlayer == nil {
		//! 角色不存在
		res.RetCode = RETCODE_FRIEND_NOT_EXIST
		return
	}

	if len(friendPlayer.DataFriend.Data.friends) >= MAX_FRIEND_NUM {
		res.RetCode = RETCODE_TARGET_FRIEND_FULL
		return
	}
	self.Locker.Lock()
	friend := self.Data.friends[friendId]
	if friend == nil {
		if len(self.Data.applys) >= MAX_FRIEND_NUM {
			res.RetCode = RETCODE_SELF_FRIEND_FULL
			self.Locker.Unlock()
			return
		}
		//看是否是黑名单
		black := friendPlayer.DataFriend.Data.black[self.Data.UId]
		if black != nil {
			res.RetCode = RETCODE_FRIEND_APPLY_BLACK
			self.Locker.Unlock()
			return
		}

		//! 没有好友，检查是否申请好友
		apply := self.Data.applys[friendId]
		if apply == nil {
			apply = new(JS_Friend)
			apply.UId = friendPlayer.Data.UId
			apply.Name = friendPlayer.Data.UName
			apply.Level = friendPlayer.Data.Level
			apply.Server = friendPlayer.Data.ServerId
			apply.Icon = friendPlayer.Data.data.IconId
			apply.Portrait = friendPlayer.Data.data.Portrait
			self.Data.applys[friendId] = apply

			res.Data = utils.HF_JtoA(apply) //加入到给游戏服的消息中进行同步
		}
		self.Locker.Unlock()

		friendPlayer.DataFriend.Locker.Lock()
		applied := friendPlayer.DataFriend.Data.applieds[self.Data.UId]
		if applied == nil {
			applied = new(JS_Friend)
			applied.UId = self.player.Data.UId
			applied.Name = self.player.Data.UName
			applied.Server = self.player.Data.ServerId
			applied.Level = self.player.Data.Level
			applied.Icon = self.player.Data.data.IconId
			applied.Portrait = self.player.Data.data.Portrait
			applied.Fight = int64(self.player.Data.Fight)
			applied.Vip = self.player.Data.data.Vip
			friendPlayer.DataFriend.Data.applieds[self.Data.UId] = applied
			friendPlayer.DataFriend.Locker.Unlock()

			//! 生成事件
			sid := friendPlayer.Data.ServerId
			core.GetCenterApp().AddEvent(sid, core.PLAYER_EVENT_ADD_FRIEND, friendPlayer.GetUId(),
				self.Data.UId, 0, utils.HF_JtoA(applied))
		} else {
			friendPlayer.DataFriend.Locker.Unlock()
		}

	} else {
		//! 解锁
		self.Locker.Unlock()
	}
	return
}

//! 同意好友
//! -1 超过好友上限
func (self *ModFriend) AgreeFriend(friendId int64, res *RPC_FriendActionRes) {
	newNum := 1
	if friendId == 0 && len(self.Data.applieds) > 0 {
		newNum = len(self.Data.applieds)
	}
	friendNum := len(self.Data.friends)
	if friendNum+newNum > MAX_FRIEND_NUM {
		res.RetCode = RETCODE_SELF_FRIEND_FULL
		return
	}
	uidArr := make([]*JS_Friend, 0)
	if friendId > 0 {
		friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
		if friendPlayer == nil {
			res.RetCode = RETCODE_FRIEND_NOT_EXIST
			return
		}

		if friendPlayer.DataFriend.GetFriendNum() >= MAX_FRIEND_NUM {
			res.RetCode = RETCODE_TARGET_FRIEND_FULL
			return
		}

		self.Locker.Lock()
		applied, ok := self.Data.applieds[friendId]
		if ok {
			delete(self.Data.applieds, friendId)
			delete(self.Data.applys, friendId)
			self.Data.friends[friendId] = applied
			uidArr = append(uidArr, applied)
			self.Locker.Unlock()
		} else {
			//! 不存在
			self.Locker.Unlock()
			res.RetCode = RETCODE_FRIEND_APPLY_ERROR
			return
		}

		friendPlayer.DataFriend.Locker.Lock()
		apply, ok := friendPlayer.DataFriend.Data.applys[self.Data.UId]
		if ok {
			friendPlayer.DataFriend.Data.friends[self.Data.UId] = apply
			delete(friendPlayer.DataFriend.Data.applys, self.Data.UId)
			delete(friendPlayer.DataFriend.Data.applieds, self.Data.UId)
			friendPlayer.DataFriend.Locker.Unlock()

			core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_AGREEE_FRIEND, friendPlayer.GetUId(),
				self.Data.UId, 0, utils.HF_JtoA(apply))
		} else {
			friendPlayer.DataFriend.Locker.Unlock()
			res.RetCode = RETCODE_FRIEND_APPLY_ERROR
			return
		}
	} else {
		for _, v := range self.Data.applieds {
			friendPlayer := GetPlayerMgr().GetPlayer(v.UId, true)
			if friendPlayer == nil {
				continue
			}

			if friendPlayer.DataFriend.GetFriendNum() >= MAX_FRIEND_NUM {
				continue
			}

			self.Data.friends[v.UId] = v
			uidArr = append(uidArr, v)

			friendPlayer.DataFriend.Locker.Lock()
			apply, ok := friendPlayer.DataFriend.Data.applys[self.Data.UId]
			if ok {
				friendPlayer.DataFriend.Data.friends[self.Data.UId] = apply
				delete(friendPlayer.DataFriend.Data.applys, self.Data.UId)
				friendPlayer.DataFriend.Locker.Unlock()

				core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_AGREEE_FRIEND, friendPlayer.GetUId(),
					self.Data.UId, 0, utils.HF_JtoA(apply))
			} else {
				friendPlayer.DataFriend.Locker.Unlock()
				continue
			}
		}
		self.Data.applieds = make(map[int64]*JS_Friend)
	}
	res.Data = utils.HF_JtoA(uidArr)
}

//! 拒绝好友
func (self *ModFriend) RefuseFriend(friendId int64, res *RPC_FriendActionRes) {
	uidArr := make([]*JS_Friend, 0)
	if friendId > 0 {
		friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
		if friendPlayer == nil {
			res.RetCode = RETCODE_FRIEND_NOT_EXIST
			return
		}

		self.Locker.Lock()
		applied, ok := self.Data.applieds[friendId]
		if ok {
			delete(self.Data.applieds, friendId)
			uidArr = append(uidArr, applied)
			self.Locker.Unlock()
		} else {
			//! 不存在
			self.Locker.Unlock()
		}

		friendPlayer.DataFriend.Locker.Lock()
		apply, ok := friendPlayer.DataFriend.Data.applys[self.Data.UId]
		if ok {
			delete(friendPlayer.DataFriend.Data.applys, self.Data.UId)
			friendPlayer.DataFriend.Locker.Unlock()

			core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_REFUSE_FRIEND, friendPlayer.GetUId(),
				self.Data.UId, 0, utils.HF_JtoA(apply))
		} else {
			friendPlayer.DataFriend.Locker.Unlock()
		}
	} else {
		for _, v := range self.Data.applieds {
			uidArr = append(uidArr, v)
			friendPlayer := GetPlayerMgr().GetPlayer(v.UId, true)
			if friendPlayer == nil {
				continue
			}

			friendPlayer.DataFriend.Locker.Lock()
			apply, ok := friendPlayer.DataFriend.Data.applys[self.Data.UId]
			if ok {
				delete(friendPlayer.DataFriend.Data.applys, self.Data.UId)
				friendPlayer.DataFriend.Locker.Unlock()

				core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_REFUSE_FRIEND, friendPlayer.GetUId(),
					self.Data.UId, 0, utils.HF_JtoA(apply))
			} else {
				friendPlayer.DataFriend.Locker.Unlock()
				continue
			}
		}
		self.Data.applieds = make(map[int64]*JS_Friend, 0)
	}
	res.Data = utils.HF_JtoA(uidArr)
}

//! 赠送友情点
func (self *ModFriend) PowerFriend(friendId []int64, res *RPC_FriendActionRes) {
	for i := 0; i < len(friendId); i++ {
		fid := friendId[i]
		friendPlayer := GetPlayerMgr().GetPlayer(fid, true)
		if friendPlayer == nil {
			continue
		}
		utils.LogDebug(self.Data.UId, "friend power to", friendPlayer.GetUId())
		core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_POWER_FRIEND, friendPlayer.GetUId(),
			self.Data.UId, 0, "")
	}
}

//拉黑好友
func (self *ModFriend) BlackFriend(friendId int64, res *RPC_FriendActionRes) {
	friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
	if friendPlayer == nil {
		res.RetCode = RETCODE_FRIEND_NOT_EXIST
		return
	}

	self.Locker.Lock()
	delete(self.Data.applieds, friendId)
	delete(self.Data.applys, friendId)
	delete(self.Data.friends, friendId)
	nodeToPlayer := new(JS_Friend)
	nodeToPlayer.UId = friendPlayer.Data.UId
	nodeToPlayer.Name = friendPlayer.Data.UName
	nodeToPlayer.Server = friendPlayer.Data.ServerId
	nodeToPlayer.Level = friendPlayer.Data.Level
	nodeToPlayer.Icon = friendPlayer.Data.data.IconId
	nodeToPlayer.Portrait = friendPlayer.Data.data.Portrait
	nodeToPlayer.Fight = int64(friendPlayer.Data.Fight)
	nodeToPlayer.Vip = friendPlayer.Data.data.Vip
	self.Data.black[friendId] = nodeToPlayer
	self.Locker.Unlock()

	friendPlayer.DataFriend.Locker.Lock()
	delete(friendPlayer.DataFriend.Data.applys, self.Data.UId)
	delete(friendPlayer.DataFriend.Data.applieds, self.Data.UId)
	delete(friendPlayer.DataFriend.Data.friends, self.Data.UId)
	nodeToBlack := new(JS_Friend)
	nodeToBlack.UId = self.player.Data.UId
	nodeToBlack.Name = self.player.Data.UName
	nodeToBlack.Server = self.player.Data.ServerId
	nodeToBlack.Level = self.player.Data.Level
	nodeToBlack.Icon = self.player.Data.data.IconId
	nodeToBlack.Portrait = self.player.Data.data.Portrait
	nodeToBlack.Fight = int64(self.player.Data.Fight)
	nodeToBlack.Vip = self.player.Data.data.Vip
	friendPlayer.DataFriend.Locker.Unlock()

	res.Data = utils.HF_JtoA(nodeToPlayer)
	core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_BLACK_FRIEND, friendPlayer.GetUId(),
		self.Data.UId, 0, utils.HF_JtoA(nodeToBlack))
}

//! 删除好友
func (self *ModFriend) DelFriend(friendId int64, res *RPC_FriendActionRes) {
	//! 删除自己的好友数据
	self.Locker.Lock()
	delete(self.Data.friends, friendId)
	res.Data = utils.HF_JtoA(friendId)
	self.Locker.Unlock()

	//! 删除对方的好友数据
	friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
	if friendPlayer != nil {
		friendPlayer.DataFriend.Locker.Lock()
		delete(friendPlayer.DataFriend.Data.friends, self.Data.UId)
		nodeToDetele := new(JS_Friend)
		nodeToDetele.UId = self.player.Data.UId
		nodeToDetele.Name = self.player.Data.UName
		nodeToDetele.Server = self.player.Data.ServerId
		nodeToDetele.Level = self.player.Data.Level
		nodeToDetele.Icon = self.player.Data.data.IconId
		nodeToDetele.Portrait = self.player.Data.data.Portrait
		nodeToDetele.Fight = int64(self.player.Data.Fight)
		nodeToDetele.Vip = self.player.Data.data.Vip
		core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_DEL_FRIEND, friendPlayer.GetUId(),
			self.Data.UId, 0, utils.HF_JtoA(nodeToDetele))
		friendPlayer.DataFriend.Locker.Unlock()
	}
	return
}

//! 移除黑名单
func (self *ModFriend) BlackOutFriend(friendId int64, res *RPC_FriendActionRes) {
	//! 删除自己的好友数据
	self.Locker.Lock()
	delete(self.Data.black, friendId)
	res.Data = utils.HF_JtoA(friendId)
	self.Locker.Unlock()
	return
}

func (self *ModFriend) AddHireHero(hire map[int]*HireHero, isSend int, res *RPC_FriendActionRes) {
	self.Locker.Lock()
	if self.Data.hireHeroInfo == nil {
		self.Data.hireHeroInfo = make(map[int]*HireHero)
	}

	for _, v := range hire {
		hireHero, ok := self.Data.hireHeroInfo[v.HeroKeyId]
		if ok {
			hireHero.HeroQuality = v.HeroQuality
			hireHero.HeroExclusiveLv = v.HeroExclusiveLv
			hireHero.HeroExclusiveUnLock = v.HeroExclusiveUnLock
			hireHero.HeroArtifactLv = v.HeroArtifactLv
			hireHero.HeroArtifactId = v.HeroArtifactId
			if hireHero.ReSetTime <= time.Now().Unix() {
				hireHero.HirePlayer = nil
				hireHero.ApplyPlayer = make([]*HirePlayerBase, 0)
				hireHero.ReSetTime = game.HF_GetNextWeekStart()
			}
			if v.Talent != nil {
				hireHero.Talent = new(game.StageTalent)
				utils.HF_DeepCopy(hireHero.Talent, v.Talent)
			}
		} else {
			self.Data.hireHeroInfo[v.HeroKeyId] = v
		}
		if isSend == game.LOGIC_TRUE {
			core.GetCenterApp().AddEvent(self.player.Data.ServerId, core.PLAYER_EVENT_UPDATE_HIRE_HERO_SINGLE, self.player.Data.UId,
				self.Data.UId, 0, utils.HF_JtoA(v))
		}
	}
	self.Locker.Unlock()
}

func (self *ModFriend) DeleteHireHero(keyId int, res *RPC_FriendActionRes) {
	self.Locker.Lock()
	if self.Data.hireHeroInfo == nil {
		self.Data.hireHeroInfo = make(map[int]*HireHero)
	}

	node, ok := self.Data.hireHeroInfo[keyId]
	if ok {
		for _, friend := range self.Data.friends {
			core.GetCenterApp().AddEvent(friend.Server, core.PLAYER_EVENT_DELETE_HIRE, friend.UId,
				self.Data.UId, 0, utils.HF_JtoA(node))
		}
		delete(self.Data.hireHeroInfo, keyId)
		core.GetCenterApp().AddEvent(self.player.Data.ServerId, core.PLAYER_EVENT_DELETE_HIRE, self.player.Data.UId,
			self.Data.UId, 0, utils.HF_JtoA(node))
	}
	self.Locker.Unlock()
}

func (self *ModFriend) HireLose(friendId int64, keyId int, res *RPC_FriendActionRes) {
	friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
	if friendPlayer == nil {
		res.RetCode = RETCODE_FRIEND_NOT_EXIST
		return
	}
	friendPlayer.DataFriend.Locker.Lock()
	info, ok := friendPlayer.DataFriend.Data.hireHeroInfo[keyId]
	if ok {
		if info.HirePlayer != nil && info.HirePlayer.Uid == self.Data.UId {
			info.HirePlayer = new(HirePlayerBase)
			res.Data = utils.HF_JtoA(info)

			self.UpdateHireMsg(friendPlayer, res.Data, core.PLAYER_EVENT_UPDATE_HIRE_HERO_SINGLE)
		}
	}
	friendPlayer.DataFriend.Locker.Unlock()

	newList := make([]*HireHero, 0)
	for i := 0; i < len(self.Data.selfHire); i++ {
		if self.Data.selfHire[i].OwnPlayer.Uid == friendId && self.Data.selfHire[i].HeroKeyId == keyId {
			continue
		}
		newList = append(newList, self.Data.selfHire[i])
	}
	self.Data.selfHire = newList
}

func (self *ModFriend) HireApply(friendId int64, keyId int, res *RPC_FriendActionRes) {
	friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
	if friendPlayer == nil {
		res.RetCode = RETCODE_FRIEND_NOT_EXIST
		return
	}
	friendPlayer.DataFriend.Locker.Lock()
	info, ok := friendPlayer.DataFriend.Data.hireHeroInfo[keyId]
	if ok {
		isHas := false
		for i := 0; i < len(info.ApplyPlayer); i++ {
			if info.ApplyPlayer[i].Uid == self.Data.UId {
				isHas = true
			}
		}
		if !isHas {
			playerBase := new(HirePlayerBase)
			playerBase.Uid = self.Data.UId
			playerBase.Name = self.player.Data.UName
			playerBase.IconId = self.player.Data.data.IconId
			playerBase.Portrait = self.player.Data.data.Portrait
			info.ApplyPlayer = append(info.ApplyPlayer, playerBase)
			res.Data = utils.HF_JtoA(info)
			self.UpdateHireMsg(friendPlayer, res.Data, core.PLAYER_EVENT_UPDATE_HIRE_HERO_SINGLE)
		}
	}
	friendPlayer.DataFriend.Locker.Unlock()
}

func (self *ModFriend) HireCancel(friendId int64, keyId int, res *RPC_FriendActionRes) {
	friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
	if friendPlayer == nil {
		res.RetCode = RETCODE_FRIEND_NOT_EXIST
		return
	}
	friendPlayer.DataFriend.Locker.Lock()
	info, ok := friendPlayer.DataFriend.Data.hireHeroInfo[keyId]
	if ok {
		newApply := make([]*HirePlayerBase, 0)
		for i := 0; i < len(info.ApplyPlayer); i++ {
			if info.ApplyPlayer[i].Uid != self.Data.UId {
				newApply = append(newApply, info.ApplyPlayer[i])
			}
		}
		info.ApplyPlayer = newApply
		res.Data = utils.HF_JtoA(info)

		self.UpdateHireMsg(friendPlayer, res.Data, core.PLAYER_EVENT_UPDATE_HIRE_HERO_SINGLE)
	}
	friendPlayer.DataFriend.Locker.Unlock()
}

func (self *ModFriend) HireAgree(friendId int64, keyId int, res *RPC_FriendActionRes) {
	friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
	if friendPlayer == nil {
		res.RetCode = RETCODE_FRIEND_NOT_EXIST
		return
	}
	now := time.Now().Unix()
	self.Locker.Lock()
	info, ok := self.Data.hireHeroInfo[keyId]
	if ok {
		if info.HirePlayer != nil && info.HirePlayer.Uid > 0 {
			self.Locker.Unlock()
			res.RetCode = RETCODE_FRIEND_HIRE_ALREADY
			return
		}
		if len(friendPlayer.DataFriend.Data.selfHire) > 0 && friendPlayer.DataFriend.Data.selfHire[0].ReSetTime <= now {
			friendPlayer.DataFriend.Data.selfHire = make([]*HireHero, 0)
		}
		if len(friendPlayer.DataFriend.Data.selfHire) >= MAX_FRIEND_HIRE_NUM {
			self.Locker.Unlock()
			res.RetCode = RETCODE_FRIEND_HIRE_MAX
			return
		}
		isHas := false
		for i := 0; i < len(info.ApplyPlayer); i++ {
			if info.ApplyPlayer[i].Uid == friendPlayer.Data.UId {
				info.HirePlayer = new(HirePlayerBase)
				info.HirePlayer.Uid = friendPlayer.Data.UId
				info.HirePlayer.Name = friendPlayer.Data.UName
				info.HirePlayer.IconId = friendPlayer.Data.data.IconId
				info.HirePlayer.Portrait = friendPlayer.Data.data.Portrait
				isHas = true
				break
			}
		}

		if isHas {
			info.ApplyPlayer = make([]*HirePlayerBase, 0)
			res.Data = utils.HF_JtoA(info)
			selfHire := new(HireHero)
			utils.HF_DeepCopy(&selfHire, &info)
			friendPlayer.DataFriend.Data.selfHire = append(friendPlayer.DataFriend.Data.selfHire, selfHire)
			core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_AGREE_HIRE_HERO, friendPlayer.Data.UId,
				0, 0, utils.HF_JtoA(info))
			self.UpdateHireMsg(self.player, res.Data, core.PLAYER_EVENT_UPDATE_HIRE_HERO_SINGLE)
		}
	}
	self.Locker.Unlock()
}

func (self *ModFriend) HireAgreeAll(res *RPC_FriendActionRes) {
	rel := make([]*HireHero, 0)
	now := time.Now().Unix()
	self.Locker.Lock()
	for _, v := range self.Data.hireHeroInfo {
		if len(v.ApplyPlayer) <= 0 {
			continue
		}
		if v.HirePlayer != nil && v.HirePlayer.Uid > 0 {
			continue
		}
		friendId := int64(0)
		for i := 0; i < len(v.ApplyPlayer); i++ {
			friendId = v.ApplyPlayer[i].Uid
			friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
			if friendPlayer == nil {
				friendId = 0
				continue
			}
			//检查下过期
			if len(friendPlayer.DataFriend.Data.selfHire) > 0 && friendPlayer.DataFriend.Data.selfHire[0].ReSetTime <= now {
				friendPlayer.DataFriend.Data.selfHire = make([]*HireHero, 0)
			}
			if len(friendPlayer.DataFriend.Data.selfHire) >= MAX_FRIEND_HIRE_NUM {
				friendId = 0
				continue
			}
			break
		}
		if friendId == 0 {
			v.ApplyPlayer = make([]*HirePlayerBase, 0)
			rel = append(rel, v)
			continue
		}
		friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
		v.HirePlayer = new(HirePlayerBase)
		v.HirePlayer.Uid = friendPlayer.Data.UId
		v.HirePlayer.Name = friendPlayer.Data.UName
		v.HirePlayer.IconId = friendPlayer.Data.data.IconId
		v.HirePlayer.Portrait = friendPlayer.Data.data.Portrait
		v.ApplyPlayer = make([]*HirePlayerBase, 0)
		rel = append(rel, v)
		selfHire := new(HireHero)
		utils.HF_DeepCopy(&selfHire, &v)
		friendPlayer.DataFriend.Data.selfHire = append(friendPlayer.DataFriend.Data.selfHire, selfHire)
		core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_AGREE_HIRE_HERO, friendPlayer.Data.UId,
			self.Data.UId, 0, utils.HF_JtoA(v))
	}

	res.Data = utils.HF_JtoA(rel)
	self.UpdateHireMsg(self.player, res.Data, core.PLAYER_EVENT_UPDATE_HIRE_HERO)
	self.Locker.Unlock()
}

func (self *ModFriend) HireRefuse(friendId int64, keyId int, res *RPC_FriendActionRes) {
	friendPlayer := GetPlayerMgr().GetPlayer(friendId, true)
	if friendPlayer == nil {
		res.RetCode = RETCODE_FRIEND_NOT_EXIST
		return
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()
	info, ok := self.Data.hireHeroInfo[keyId]
	if ok {
		if info.HirePlayer != nil && info.HirePlayer.Uid > 0 {
			res.RetCode = RETCODE_FRIEND_HIRE_ALREADY
			return
		}
		newApply := make([]*HirePlayerBase, 0)
		for i := 0; i < len(info.ApplyPlayer); i++ {
			if info.ApplyPlayer[i].Uid != friendPlayer.Data.UId {
				newApply = append(newApply, info.ApplyPlayer[i])
			}
		}
		info.ApplyPlayer = newApply
		res.Data = utils.HF_JtoA(info)
		core.GetCenterApp().AddEvent(friendPlayer.Data.ServerId, core.PLAYER_EVENT_REFUSE_HIRE_HERO, friendPlayer.Data.UId,
			0, 0, utils.HF_JtoA(info))
		self.UpdateHireMsg(self.player, res.Data, core.PLAYER_EVENT_UPDATE_HIRE_HERO_SINGLE)
	}
}

func (self *ModFriend) GetFriend(res *RPC_FriendActionRes) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	for _, v := range self.Data.friends {
		player := GetPlayerMgr().GetPlayer(v.UId, true)
		if player != nil {
			v.Name = player.Data.UName
			v.Icon = player.Data.data.IconId
			v.Portrait = player.Data.data.Portrait
			v.Fight = player.Data.data.Fight
			v.Stage = player.Data.data.PassId
			v.LastUpTime = player.Data.data.LastUpdate
			v.Online = player.Online
			v.Level = player.Data.data.Level
			v.Vip = player.Data.data.Vip
		}
	}
	res.Data = utils.HF_JtoA(self.Data.friends)
}

func (self *ModFriend) GetApply(res *RPC_FriendActionRes) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	rel := []*JS_Friend{}
	for _, v := range self.Data.applieds {
		player := GetPlayerMgr().GetPlayer(v.UId, true)
		if player != nil {
			v.Name = player.Data.UName
			v.Icon = player.Data.data.IconId
			v.Portrait = player.Data.data.Portrait
			v.Fight = player.Data.data.Fight
			v.Stage = player.Data.data.PassId
			v.LastUpTime = player.Data.data.LastUpdate
			v.Online = player.Online
			v.Level = player.Data.data.Level
			v.Vip = player.Data.data.Vip
			rel = append(rel, v)
		}
	}
	res.Data = utils.HF_JtoA(rel)
}

func (self *ModFriend) GetBlack(res *RPC_FriendActionRes) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	rel := []*JS_Friend{}
	for _, v := range self.Data.black {
		player := GetPlayerMgr().GetPlayer(v.UId, true)
		if player != nil {
			v.Name = player.Data.UName
			v.Icon = player.Data.data.IconId
			v.Portrait = player.Data.data.Portrait
			v.Fight = player.Data.data.Fight
			v.Stage = player.Data.data.PassId
			v.LastUpTime = player.Data.data.LastUpdate
			v.Online = player.Online
			v.Level = player.Data.data.Level
			v.Vip = player.Data.data.Vip
			rel = append(rel, v)
		}
	}
	res.Data = utils.HF_JtoA(rel)
}

func (self *ModFriend) GetHireList(res *RPC_FriendActionRes) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	rel := make(map[int][]*HireHero)
	now := time.Now().Unix()

	uidMap := make(map[int64]int, 0)
	for _, v := range self.Data.friends {
		uidMap[v.UId]++
	}

	unionInfo := union.GetUnionMgr().GetUnion(self.player.Data.data.UnionId)
	if unionInfo != nil {
		braveHandUId := unionInfo.GetBraveHandUid()
		for _, v := range braveHandUId {
			uidMap[v] = game.LOGIC_TRUE
		}
	}

	for key := range uidMap {
		if key == self.player.Data.UId {
			continue
		}
		player := GetPlayerMgr().GetPlayer(key, true)
		if player != nil {
			for _, hire := range player.DataFriend.Data.hireHeroInfo {
				if hire.ReSetTime <= now {
					hire.HirePlayer = nil
					hire.ApplyPlayer = make([]*HirePlayerBase, 0)
					hire.ReSetTime = game.HF_GetNextWeekStart()
				}
				rel[hire.HeroId] = append(rel[hire.HeroId], hire)
			}
		}
	}
	res.Data = utils.HF_JtoA(rel)
}

func (self *ModFriend) GetSelfList(res *RPC_FriendActionRes) {

	now := time.Now().Unix()
	for _, v := range self.Data.hireHeroInfo {
		if v.ReSetTime <= now {
			v.ApplyPlayer = make([]*HirePlayerBase, 0)
			v.HirePlayer = new(HirePlayerBase)
			v.ReSetTime = game.HF_GetNextWeekStart()
		} else {
			for _, playerNode := range v.ApplyPlayer {
				player := GetPlayerMgr().GetPlayer(playerNode.Uid, true)
				if player != nil {
					playerNode.Name = player.Data.UName
					playerNode.Portrait = player.Data.data.Portrait
					playerNode.IconId = player.Data.data.IconId
				}
			}

			if v.HirePlayer != nil && v.HirePlayer.Uid > 0 {
				player := GetPlayerMgr().GetPlayer(v.HirePlayer.Uid, true)
				if player != nil {
					v.HirePlayer.Name = player.Data.UName
					v.HirePlayer.Portrait = player.Data.data.Portrait
					v.HirePlayer.IconId = player.Data.data.IconId
				}
			}

			if v.OwnPlayer == nil {
				v.OwnPlayer = new(HirePlayerBase)
				v.OwnPlayer.Uid = self.player.Data.UId
				v.OwnPlayer.Name = self.player.Data.UName
				v.OwnPlayer.Portrait = self.player.Data.data.Portrait
				v.OwnPlayer.IconId = self.player.Data.data.IconId
			} else {
				v.OwnPlayer.Uid = self.player.Data.UId
				v.OwnPlayer.Name = self.player.Data.UName
				v.OwnPlayer.Portrait = self.player.Data.data.Portrait
				v.OwnPlayer.IconId = self.player.Data.data.IconId
			}
		}
	}
	res.Data = utils.HF_JtoA(self.Data.hireHeroInfo)
}

//保存私聊，不需要返回
func (self *ModFriend) SavePrivateMessage(aim int64, content string) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	plTo := GetPlayerMgr().GetPlayer(aim, true)
	if plTo == nil {
		return
	}

	msgInfo := &chat.ChatMessage{
		MsgId:      0,
		Uid:        self.player.Data.data.UId,
		Uname:      self.player.Data.data.UName,
		Level:      self.player.Data.data.Level,
		IconId:     self.player.Data.data.IconId,
		Portrait:   self.player.Data.data.Portrait,
		Content:    content,
		ToUid:      plTo.Data.data.UId,
		ToUname:    plTo.Data.data.UName,
		ToLevel:    plTo.Data.data.Level,
		ToIconId:   plTo.Data.data.IconId,
		ToPortrait: plTo.Data.data.Portrait,
		SendTime:   int(time.Now().Unix()),
	}
	self.Data.chatSendInfo = append(self.Data.chatSendInfo, msgInfo)
	plTo.DataFriend.Data.chatGetInfo = append(plTo.DataFriend.Data.chatGetInfo, msgInfo)
}

//上线获得私聊
func (self *ModFriend) QueryPrivateMessage(res *chat.RPC_ChatActionRes) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	now := time.Now().Unix()
	indexSend := -1
	for i := 0; i < len(self.Data.chatSendInfo); i++ {
		if int64(self.Data.chatSendInfo[i].SendTime)+utils.DAY_SECS*3 < now {
			indexSend = i
		} else {
			break
		}
	}
	if indexSend >= 0 {
		self.Data.chatSendInfo = self.Data.chatSendInfo[indexSend+1:]
	}

	indexGet := 0
	for i := 0; i < len(self.Data.chatGetInfo); i++ {
		if int64(self.Data.chatGetInfo[i].SendTime)+utils.DAY_SECS*3 < now {
			indexGet = i
		} else {
			break
		}
	}
	if indexGet > 0 {
		self.Data.chatGetInfo = self.Data.chatGetInfo[indexGet:]
	}

	res.MsgList = append(res.MsgList, self.Data.chatSendInfo...)
	res.MsgList = append(res.MsgList, self.Data.chatGetInfo...)
}

func (self *ModFriend) UpdateHireMsg(player *Player, str string, msg int) {

	uidMap := make(map[int64]int, 0)
	for _, v := range player.DataFriend.Data.friends {
		uidMap[v.UId]++
	}
	uidMap[player.Data.data.UId]++
	unionInfo := union.GetUnionMgr().GetUnion(player.Data.data.UnionId)
	if unionInfo != nil {
		member := unionInfo.GetMemberList()
		for _, v := range member {
			uidMap[v] = game.LOGIC_TRUE
		}
	}
	for key := range uidMap {
		noticePlayer := GetPlayerMgr().GetPlayer(key, true)
		if noticePlayer != nil {
			core.GetCenterApp().AddEvent(noticePlayer.Data.ServerId, msg, noticePlayer.Data.UId,
				player.Data.UId, 0, str)
		}
	}
}

func (self *ModFriend) RobotAction(res *RPC_FriendActionRes) {

	if len(self.Data.friends) > 25 {
		return
	}
	//获得5个人加好友
	uid := GetPlayerMgr().RandPlayer(self.Data.UId, 5)
	for _, v := range uid {
		player := GetPlayerMgr().GetPlayer(v, true)
		if player != nil {
			self.Locker.Lock()
			friend, ok := self.Data.friends[player.GetUid()]
			if !ok {
				friend = new(JS_Friend)
				friend.UId = player.Data.UId
				friend.Name = player.Data.UName
				friend.Server = player.Data.ServerId
				friend.Level = player.Data.Level
				friend.Icon = player.Data.data.IconId
				friend.Portrait = player.Data.data.Portrait
				friend.Fight = int64(player.Data.Fight)
				friend.Vip = player.Data.data.Vip
				self.Data.friends[player.GetUid()] = friend
			}
			self.Locker.Unlock()

			player.DataFriend.Locker.Lock()
			_, ok1 := player.DataFriend.Data.friends[self.Data.UId]
			if !ok1 {
				friend = new(JS_Friend)
				friend.UId = self.player.Data.UId
				friend.Name = self.player.Data.UName
				friend.Server = self.player.Data.ServerId
				friend.Level = self.player.Data.Level
				friend.Icon = self.player.Data.data.IconId
				friend.Portrait = self.player.Data.data.Portrait
				friend.Fight = int64(self.player.Data.Fight)
				friend.Vip = self.player.Data.data.Vip
				self.Data.friends[self.player.GetUid()] = friend
			}
			player.DataFriend.Locker.Unlock()
		}
	}
}
