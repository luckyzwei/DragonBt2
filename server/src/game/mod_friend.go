package game

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	//"time"
)

const MAX_FRIEND = 30
const MAX_FRIEND_POWER = 20

const (
	HIRE_HERO_MAX       = 3
	HIRE_HERO_APPLY_MAX = 5
)

const (
	ROBOT_EQUIP_SUBTYPE_HIRE = 1
	ROBOT_EQUIP_SUBTYPE_PIT  = 2
)

const (
	ACTION_APPLY        = 0 //申请
	ACTION_APPLY_CANCEL = 1 //放弃申请
	ACTION_HIRE_LOSE    = 2 //放弃租用
)

//标志索引
const (
	HIRE_MOD_PASS         = 1 //主线使用
	HIRE_MOD_TOWER        = 2 //塔使用
	HIRE_MOD_TOWER_CAMP_1 = 3 //种族塔使用
	HIRE_MOD_TOWER_CAMP_2 = 4 //种族塔使用
	HIRE_MOD_TOWER_CAMP_3 = 5 //种族塔使用
	HIRE_MOD_TOWER_CAMP_4 = 6 //种族塔使用
	HIRE_MOD_CAMP_END     = 7
)

type FriendHero struct {
	OwnUid    int64           `json:"ownuid"`    //拥有者ID
	OwnPlayer *HirePlayerBase `json:"ownplayer"` //拥有者
	KeyId     int             `json:"keyid"`
	Hero      *NewHero        `json:"hero"`
	SelfKeyId int             `json:"selfkey"`
}

//! 好友数据库
type San_Friend struct {
	Uid       int64
	Friend    string //! 好友
	Apply     string //! 申请列表
	Black     string //! 黑名单
	Commend   string //! 推荐列表
	Hasapply  string //! 已经申请的列表
	Count     int    //! 今日领取次数
	ApplyHire string //! 申请租用列表
	HireHero  string //! 租到的英雄
	HireTime  int64  //! 过期时间
	UseSign   string //!
	GiftSign  string //!

	friend    *sync.Map
	apply     *sync.Map
	black     *sync.Map
	commend   *sync.Map
	hasapply  *sync.Map
	applyHire []*HireHero   //! 申请租用列表
	hireHero  []*FriendHero //! 租到的英雄
	useSign   []int         //使用标记
	giftSign  []int64       //用来解决删好友再加还能送的问题

	DataUpdate
}

type JS_Friend struct {
	Uid        int64  `json:"uid"`
	Name       string `json:"name"`
	Icon       int    `json:"icon"`
	Level      int    `json:"level"`
	Vip        int    `json:"vip"`
	Fight      int64  `json:"fight"`
	Online     int    `json:"online"`
	Portrait   int    `json:"portrait"` //20190412 增加边框 by zy
	Stage      int    `json:"stage"`
	Server     int    `json:"server"`
	LastUpTime int64  `json:"lastuptime"`
}

type JS_FriendNode struct {
	Gift       int    `json:"gift"`
	Get        int    `json:"get"`
	Uid        int64  `json:"uid"`
	Name       string `json:"name"`
	Icon       int    `json:"icon"`
	Portrait   int    `json:"portrait"`
	Level      int    `json:"level"`
	Vip        int    `json:"vip"`
	Fight      int64  `json:"fight"`
	Online     int    `json:"online"`
	LastUpTime int64  `json:"lastuptime"`
	Stage      int    `json:"stage"`
	Server     int    `json:"server"`
}

//! 好友
type ModFriend struct {
	player     *Player
	Sql_Friend San_Friend //! 数据库结构
}

func NewJsFriend(player *Player) *JS_Friend {
	stage, server := GetOfflineInfoMgr().GetBaseInfo(player.Sql_UserBase.Uid)

	return &JS_Friend{
		Uid:      player.Sql_UserBase.Uid,
		Name:     player.Sql_UserBase.UName,
		Icon:     player.Sql_UserBase.IconId,
		Level:    player.Sql_UserBase.Level,
		Vip:      player.Sql_UserBase.Vip,
		Fight:    player.Sql_UserBase.Fight,
		Portrait: player.Sql_UserBase.Portrait,
		Online:   0,
		Stage:    stage,
		Server:   server,
	}
}

func NewJsFriendByJsPlayerData(player *JS_PlayerData) *JS_Friend {

	return &JS_Friend{
		Uid:      player.UId,
		Name:     player.UName,
		Icon:     player.IconId,
		Level:    player.Level,
		Vip:      player.Vip,
		Fight:    player.Fight,
		Portrait: player.Portrait,
		Online:   LOGIC_TRUE,
		Stage:    player.PassId,
		Server:   player.ServerId,
	}
}

func NewJsFriendNode(gift int, get int, uid int64, name string, icon int, portrait int,
	level int, vip int, fight int64, online int) *JS_FriendNode {
	return &JS_FriendNode{
		Gift:     gift,
		Get:      get,
		Uid:      uid,
		Name:     name,
		Icon:     icon,
		Portrait: portrait,
		Level:    level,
		Vip:      vip,
		Fight:    fight,
		Online:   online,
	}
}

func (self *JS_Friend) SetOnline() {
	self.Online = 1
}

func (self *JS_FriendNode) SetOnline() {
	self.Online = 0
	if TimeServer().Unix()-self.LastUpTime < 900 {
		self.Online = 1
	}
}

func (self *San_Friend) Decode() {
	// friend info
	friendInfo := make([]*JS_FriendNode, 0)
	json.Unmarshal([]byte(self.Friend), &friendInfo)
	self.friend = new(sync.Map)
	for _, v := range friendInfo {
		self.friend.Store(v.Uid, v)
	}

	// apply
	applyInfo := make([]*JS_Friend, 0)
	json.Unmarshal([]byte(self.Apply), &applyInfo)
	for _, v := range applyInfo {
		self.apply.Store(v.Uid, v)
	}

	// bloack
	blackInfo := make([]*JS_Friend, 0)
	json.Unmarshal([]byte(self.Black), &blackInfo)
	for _, v := range blackInfo {
		self.black.Store(v.Uid, v)
	}

	// commend
	commendInfo := make([]*JS_Friend, 0)
	json.Unmarshal([]byte(self.Commend), &commendInfo)
	for _, v := range commendInfo {
		self.commend.Store(v.Uid, v)
	}

	// has apply
	hasapplyInfo := make([]*JS_Friend, 0)
	json.Unmarshal([]byte(self.Hasapply), &hasapplyInfo)
	for _, v := range hasapplyInfo {
		self.hasapply.Store(v.Uid, v)
	}

	json.Unmarshal([]byte(self.ApplyHire), &self.applyHire)
	json.Unmarshal([]byte(self.HireHero), &self.hireHero)
	json.Unmarshal([]byte(self.UseSign), &self.useSign)
	json.Unmarshal([]byte(self.GiftSign), &self.giftSign)
}

func (self *ModFriend) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_friend` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Friend, "san_friend", self.player.ID)
	self.checkMapNil()
	if self.Sql_Friend.Uid <= 0 {
		self.Sql_Friend.Uid = self.player.ID
		//self.Encode() 好友模块getfriend第一次这个地方如果写数据,会因为其他模块没初始化完成而出现报错,流程需要优化 zy
		InsertTable("san_friend", &self.Sql_Friend, 0, true)
	} else {
		self.Decode()
	}
	//self.CheckOfflineFriend(player)
	//self.CheckHireHero(player)

	size := len(self.Sql_Friend.UseSign)
	for i := 0; i < HIRE_MOD_CAMP_END-size; i++ {
		self.Sql_Friend.useSign = append(self.Sql_Friend.useSign, LOGIC_FALSE)
	}

	self.Sql_Friend.Init("san_friend", &self.Sql_Friend, true)
}

func (self *ModFriend) checkMapNil() {
	if self.Sql_Friend.friend == nil {
		self.Sql_Friend.friend = new(sync.Map)
	}

	if self.Sql_Friend.apply == nil {
		self.Sql_Friend.apply = new(sync.Map)
	}

	if self.Sql_Friend.black == nil {
		self.Sql_Friend.black = new(sync.Map)
	}

	if self.Sql_Friend.commend == nil {
		self.Sql_Friend.commend = new(sync.Map)
	}

	if self.Sql_Friend.hasapply == nil {
		self.Sql_Friend.hasapply = new(sync.Map)
	}
}

//! 同步好友状态
func (self *ModFriend) OnGetOtherData() {
	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_FriendNode)
		if ok {
			friend := GetPlayerMgr().GetPlayer(pInfo.Uid, false)
			if friend != nil {
				pInfo.Vip = friend.Sql_UserBase.Vip
				pInfo.Level = friend.Sql_UserBase.Level
				pInfo.Fight = friend.Sql_UserBase.Fight
			}
		}
		return true
	})

	self.Sql_Friend.hasapply.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_Friend)
		if ok {
			friend := GetPlayerMgr().GetPlayer(pInfo.Uid, false)
			if friend != nil {
				pInfo.Vip = friend.Sql_UserBase.Vip
				pInfo.Level = friend.Sql_UserBase.Level
				pInfo.Fight = friend.Sql_UserBase.Fight
			}
		}
		return true
	})

	self.Sql_Friend.apply.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_Friend)
		if ok {
			friend := GetPlayerMgr().GetPlayer(pInfo.Uid, false)
			if friend != nil {
				pInfo.Vip = friend.Sql_UserBase.Vip
				pInfo.Level = friend.Sql_UserBase.Level
				pInfo.Fight = friend.Sql_UserBase.Fight
			}
		}
		return true
	})
}

func (self *ModFriend) onReg(handlers map[string]func(body []byte)) {
	handlers["baseheroset"] = self.BaseHeroSet
	handlers["hireapply"] = self.HireApply
	handlers["hirestateset"] = self.HireStateSet
	handlers["hirestatesetall"] = self.HireStateSetAll //1
	handlers["updatehirelist"] = self.UpdateHireList   //1
	handlers["updatehireusesign"] = self.UpdateHireUseSign
	handlers["getplayerteam"] = self.GetPlayerTeam
	handlers["getfriendinfo"] = self.GetFriendInfo
	handlers["checkbasehero"] = self.CheckBaseHero
	handlers["robotaction"] = self.RobotAction
}

func (self *ModFriend) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "friendcommend":
		var msg C2S_FriendCommend
		json.Unmarshal(body, &msg)
		self.GetCanFriend(msg.Refresh)
		return true
	case "friendorder":
		var msg C2S_FriendOrder
		json.Unmarshal(body, &msg)
		self.Order(msg.Pid, msg.Agree)
		return true
	case "friendblack":
		var msg C2S_FriendBlack
		json.Unmarshal(body, &msg)
		self.Black(msg.Pid)
		return true
	case "frienddel":
		var msg C2S_FriendDel
		json.Unmarshal(body, &msg)
		self.Delete(msg.Pid, msg.Type)
		return true
	case "frienddel_all":
		var msg C2S_FriendDelBatch
		json.Unmarshal(body, &msg)
		self.DeleteBatch(msg.Pid, msg.Type)
		return true
	case "friendapply":
		var msg C2S_FriendApply
		json.Unmarshal(body, &msg)
		self.Apply(msg.Pid)
		return true
	case "friendfind":
		var msg C2S_FriendFind
		json.Unmarshal(body, &msg)
		self.Find(msg.FriendUid, msg.Name)
		return true
	case "look":
		var msg C2S_Look
		json.Unmarshal(body, &msg)
		self.Look(msg.Pid)
		return true
	case "friendpower":
		var msg C2S_FriendPower
		json.Unmarshal(body, &msg)
		self.Power(msg.Pid, msg.Type)
		/*
			for i := 0; i < len(lst); i++ {
				friend := GetPlayerMgr().GetPlayer(lst[i], false)
				if friend == nil || friend.IsOnline() == false {
					//投递到离线信息
					GetOfflineInfoMgr().SetFriendInfo(lst[i], self.player.Sql_UserBase.Uid)
					continue
				}
				friend.GetModule("friend").(*ModFriend).SetGift(self.player.Sql_UserBase.Uid, 1)
			}
		*/
		return true
	}

	return false
}

func (self *ModFriend) OnSave(sql bool) {
	self.Encode()
	self.Sql_Friend.Update(sql)
}

func (self *ModFriend) Decode() {
	//! 将数据库数据写入data
	self.Sql_Friend.Decode()
}

func (self *ModFriend) Encode() {
	//! 将data数据写入数据库
	friendInfo := self.getFriend()
	self.Sql_Friend.Friend = HF_JtoA(&friendInfo)

	applyInfo := self.getApply()
	self.Sql_Friend.Apply = HF_JtoA(&applyInfo)

	blackInfo := self.getBlack()
	self.Sql_Friend.Black = HF_JtoA(&blackInfo)

	commendInfo := self.getCommend()
	self.Sql_Friend.Commend = HF_JtoA(&commendInfo)

	hasApplyInfo := self.getHasApplyInfo()
	self.Sql_Friend.Hasapply = HF_JtoA(&hasApplyInfo)

	self.Sql_Friend.ApplyHire = HF_JtoA(&self.Sql_Friend.applyHire)
	self.Sql_Friend.HireHero = HF_JtoA(&self.Sql_Friend.hireHero)
	self.Sql_Friend.UseSign = HF_JtoA(&self.Sql_Friend.useSign)
	self.Sql_Friend.GiftSign = HF_JtoA(&self.Sql_Friend.giftSign)
}

func (self *ModFriend) Refresh() {
	self.Sql_Friend.Count = 0
	self.Sql_Friend.giftSign = make([]int64, 0)
	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_FriendNode)
		if ok && pInfo != nil {
			if pInfo.Get != 1 {
				pInfo.Get = 0
			}
			pInfo.Gift = 0
		}
		return true
	})

}

func (self *ModFriend) IsHasFriend(uid int64) bool {
	_, ok := self.Sql_Friend.friend.Load(uid)

	return ok
}

func (self *ModFriend) IsHasBlack(uid int64) bool {
	_, ok := self.Sql_Friend.black.Load(uid)

	return ok
}

func (self *ModFriend) IsHasHasapply(uid int64) bool {
	_, ok := self.Sql_Friend.hasapply.Load(uid)
	return ok
}

func (self *ModFriend) IsHasCommend(uid int64) bool {
	_, ok := self.Sql_Friend.commend.Load(uid)
	return ok
}

func (self *ModFriend) AddFriend(player *Player) bool {
	uid := player.Sql_UserBase.Uid
	_, ok := self.Sql_Friend.friend.Load(uid)
	if ok {
		return false
	}

	user := player.Sql_UserBase
	node := NewJsFriendNode(0, 0, uid, user.UName, user.IconId, user.Portrait, user.Level, user.Vip, user.Fight, 0)

	self.Sql_Friend.friend.Store(uid, node)
	self.player.HandleTask(FriendCount, 0, 0, 0)
	stage, server := GetOfflineInfoMgr().GetBaseInfo(uid)
	friend := &JS_Friend{
		Uid:      uid,
		Name:     player.Sql_UserBase.UName,
		Icon:     player.Sql_UserBase.IconId,
		Portrait: player.Sql_UserBase.Portrait,
		Level:    player.Sql_UserBase.Level,
		Vip:      player.Sql_UserBase.Vip,
		Fight:    player.Sql_UserBase.Fight,
		Stage:    stage,
		Server:   server,
	}
	self.AddInfo("addfriend", friend)

	return true
}

func (self *ModFriend) AddFriendByCenter(node *JS_Friend) {
	//先看自己是否有这个好友
	isFind := false
	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_FriendNode)
		if !ok {
			return true
		}
		if pInfo.Uid == node.Uid {
			isFind = true
			return false
		}
		return true
	})

	if isFind {
		return
	}

	nodeSave := NewJsFriendNode(0, 0, node.Uid, node.Name, node.Icon, node.Portrait, node.Level, node.Vip, node.Fight, 0)

	nCount := 0
	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_FriendNode)
		if ok && pInfo != nil {
			nCount++
		}
		return true
	})

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FRIEND_ADD, int(node.Uid), nCount, 0, "添加好友", 0, 0, self.player)

	self.Sql_Friend.friend.Store(nodeSave.Uid, nodeSave)
	self.player.HandleTask(FriendCount, 0, 0, 0)
	self.AddInfo("addfriend", node)
}

func (self *ModFriend) deleteFriend(uid int64) {
	self.Sql_Friend.friend.Delete(uid)
}

func (self *ModFriend) DelFriend(uid int64, send bool, writeLog int) bool {
	self.Sql_Friend.friend.Delete(uid)
	if send == true {
		self.DelInfo("delfriend", uid)
	}
	nCount := 0
	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_FriendNode)
		if ok && pInfo != nil {
			nCount++
		}
		return true
	})

	if writeLog == LOGIC_TRUE {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FRIEND_REMOVE, int(uid), nCount, 0, "删除好友", 0, 0, self.player)
	}
	return true
}

func (self *ModFriend) addBlack(uid int64, node *JS_Friend) {
	self.Sql_Friend.black.Store(uid, node)
	self.AddInfo("addblack", node)
}

func (self *ModFriend) AddBlack(player *Player) bool {
	uid := player.Sql_UserBase.Uid
	node := NewJsFriend(player)
	self.addBlack(uid, node)

	return true
}

func (self *ModFriend) deleteBlack(uid int64) {
	self.Sql_Friend.black.Delete(uid)
	self.DelInfo("delblack", uid)
}

func (self *ModFriend) hasBlack(uid int64) bool {
	_, ok := self.Sql_Friend.black.Load(uid)
	return ok
}

func (self *ModFriend) DelBlack(uid int64) bool {
	if !self.hasBlack(uid) {
		return false
	}
	self.deleteBlack(uid)
	return true
}

func (self *ModFriend) hasApply(uid int64) bool {
	_, ok := self.Sql_Friend.apply.Load(uid)
	return ok
}

func (self *ModFriend) addApply(uid int64, node *JS_Friend) {
	self.Sql_Friend.apply.Store(uid, node)
}

func (self *ModFriend) getApplySize() int {
	size := 0
	self.Sql_Friend.apply.Range(func(key, value interface{}) bool {
		size += 1
		return true
	})

	return size
}

func (self *ModFriend) AddApply(info *JS_Friend) bool {
	if info == nil {
		return false
	}
	self.addApply(info.Uid, info)
	self.AddInfo("addfriendapply", info)
	return true
}

func (self *ModFriend) delApply(uid int64) {
	self.Sql_Friend.apply.Delete(uid)
	self.DelInfo("delfriendapply", uid)
}

func (self *ModFriend) DelApply(uid int64) bool {
	ok := self.hasApply(uid)
	if !ok {
		return false
	}

	self.delApply(uid)
	return true
}

func (self *ModFriend) hasHasApply(uid int64) bool {
	_, ok := self.Sql_Friend.hasapply.Load(uid)
	return ok
}

func (self *ModFriend) getHasApply(uid int64) int {
	size := 0
	self.Sql_Friend.hasapply.Range(func(key, value interface{}) bool {
		size += 1
		return true
	})
	return size
}

func (self *ModFriend) addHasApply(uid int64, node *JS_Friend) {
	self.Sql_Friend.hasapply.Store(uid, node)
}

func (self *ModFriend) AddHasapply(player *Player) bool {
	uid := player.Sql_UserBase.Uid
	ok := self.hasHasApply(uid)
	if ok {
		return false
	}

	size := self.getHasApply(uid)
	if size >= 50 {
		return false
	}

	node := NewJsFriend(player)
	self.addHasApply(uid, node)

	return true
}

func (self *ModFriend) DelHasapply(uid int64) {
	self.Sql_Friend.hasapply.Delete(uid)
	self.DelInfo("delfriendapply", uid)
}

func (self *ModFriend) deleteCommend(uid int64) {
	self.Sql_Friend.commend.Delete(uid)
}

func (self *ModFriend) addCommend(lst []*JS_Friend, clear bool) {
	if clear {
		self.Sql_Friend.commend = new(sync.Map)
	}
	for _, v := range lst {
		self.Sql_Friend.commend.Store(v.Uid, v)
	}
}

func (self *ModFriend) DelCanFriend(uid int64, refresh bool) {
	self.deleteCommend(uid)
	if refresh {
		size := self.getCommendSize()
		if size >= 10 {
			return
		}

		//! 不足10个凑齐10个
		lst := self.AddCanFriend(10-size, uid)
		self.addCommend(lst, false)
		self.SendCommend()
	}
}

func (self *ModFriend) AddCanFriend(num int, delete_uid int64) []*JS_Friend {
	lst := make([]*JS_Friend, 0)
	//! 从在线人里面随机
	online := GetPlayerMgr().GetOnlineRandom(num)
	for index := 0; index < len(online); index++ {
		player := online[index]
		uid := player.Sql_UserBase.Uid
		if uid == self.player.GetUid() {
			continue
		}

		if uid == delete_uid {
			continue
		}

		if self.IsHasFriend(uid) {
			continue
		}

		if self.IsHasBlack(uid) {
			continue
		}

		if self.IsHasHasapply(uid) {
			continue
		}

		if self.IsHasCommend(uid) {
			continue
		}

		node := NewJsFriend(player)
		lst = append(lst, node)
	}
	return lst
}

//! 得到推荐好友
func (self *ModFriend) GetCanFriend(refresh int) {
	if refresh == 1 {
		lst := self.AddCanFriend(10, 0)
		self.addCommend(lst, true)
	} else {
		size := self.getCommendSize()
		if size < 10 {
			lst := self.AddCanFriend(10-size, 0)
			self.addCommend(lst, false)
		}
	}
	self.SendCommend()
}

func (self *ModFriend) SetGift(uid int64, gift int) {
	pInfo, ok := self.Sql_Friend.friend.Load(uid)
	if !ok {
		return
	}

	friend, ok := pInfo.(*JS_FriendNode)
	friend.Get = gift

	var msg S2C_FriendPower
	msg.Cid = "friendpower"
	msg.Uid = append(msg.Uid, uid)
	msg.Type = 1
	msg.Value = gift
	msg.Count = self.Sql_Friend.Count
	msg.Item = make([]PassItem, 0)
	self.player.SendMsg("friendpower", HF_JtoB(&msg))
}

//! 申请好友
func (self *ModFriend) Apply(uid int64) {
	if uid == self.player.Sql_UserBase.Uid {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_HASFRIEND"))
		return
	}

	if self.IsHasFriend(uid) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_HASFRIEND"))
		return
	}

	rel := GetMasterMgr().FriendPRC.ApplyFriend(self.player.GetUid(), uid)
	if rel == nil {
		return
	}
	if self.CheckErr(rel.RetCode) {
		return
	}

	node := &JS_Friend{}
	json.Unmarshal([]byte(rel.Data), &node)
	self.addHasApply(uid, node)

	count := 0
	self.Sql_Friend.hasapply.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FRIEND_APPLY, count, 0, 0, "好友申请", 0, 0, self.player)

	//! 申请成功
	self.player.SendRet2("applyfriend")
}

//! 查找
func (self *ModFriend) Find(friendUid int64, friendName string) {
	rel := GetMasterMgr().FriendPRC.FindFriend(self.player.Sql_UserBase.Uid, friendUid, friendName)

	data := []*JS_Friend{}

	var friendData []*JS_PlayerData
	if rel != nil {
		json.Unmarshal([]byte(rel.Data), &friendData)
	}

	for _, v := range friendData {
		if v == nil {
			continue
		}
		node := NewJsFriendByJsPlayerData(v)
		data = append(data, node)
	}
	var msg S2C_FriendCommend
	msg.Cid = "friendcommend"
	msg.Commend = data
	self.player.SendMsg("friendcommend", HF_JtoB(&msg))
}

//! 查看
func (self *ModFriend) Look(uid int64) {
	data := GetMasterMgr().GetPlayer(uid)
	if data != nil {
		if data.Data == nil {
			return
		}
		var msg S2C_Look
		msg.Cid = "look"
		msg.Uid = data.Data.UId
		msg.Name = data.Data.UName
		msg.Icon = data.Data.IconId
		msg.Vip = data.Data.Vip
		msg.Level = data.Data.Level
		msg.Fight = data.Data.Fight
		msg.Portrait = data.Data.Portrait
		msg.UnionID = data.Data.UnionId
		if msg.UnionID == 0 {
			msg.Party = ""
			msg.Office = 4
			msg.BraveHand = 0
		} else {
			msg.Party = data.Data.UnionName
			msg.Office = data.Data.Position
			msg.BraveHand = data.Data.BraveHand
		}
		if data.Online == LOGIC_TRUE {
			msg.Time = 0
		} else {
			msg.Time = TimeServer().Unix() - data.Data.LastUpdate
		}
		msg.Hero = self.GetLookBaseHero(data)
		msg.Portrait = data.Data.Portrait
		msg.Stage = data.Data.PassId
		msg.Server = data.Data.ServerId
		msg.Signature = data.Data.UserSignature

		//msg.TreeLevel = data.LifeTree.MainLevel
		//msg.TreeInfo = data.LifeTree.Info

		self.player.SendMsg("look", HF_JtoB(&msg))

		if pInfo, ok := self.Sql_Friend.friend.Load(uid); ok {
			friend, ok2 := pInfo.(*JS_FriendNode)
			if ok2 {
				friend.Level = data.Data.Level
				friend.Vip = data.Data.Vip
				friend.Fight = data.Data.Fight
			}
		}
	}
}

//! 好友体力    返回值修改为用户UID，支持离线赠送 20200108 by zy
func (self *ModFriend) Power(uid int64, _type int) {

	lst := make([]int64, 0)
	if _type == 1 && self.Sql_Friend.Count >= MAX_FRIEND_POWER {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_POWER"))
		return
	}

	var msg S2C_FriendPower
	msg.Cid = "friendpower"
	msg.Uid = make([]int64, 0)
	msg.Item = make([]PassItem, 0)
	power := 0

	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		if _type == 0 && len(self.Sql_Friend.giftSign) >= MAX_FRIEND {
			return false
		}
		pInfo, ok := value.(*JS_FriendNode)
		if !ok {
			return true
		}

		if uid > 0 && pInfo.Uid != uid {
			return true
		}

		//容错??迷之代码 先放着 by zy
		if pInfo.Uid == self.player.Sql_UserBase.Uid {
			return true
		}

		if _type == 0 {
			//! 赠送
			if pInfo.Gift > 0 {
				if uid > 0 {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
					return false
				} else {
					return true
				}
			}

			pInfo.Gift = 1
			lst = append(lst, pInfo.Uid)
			self.Sql_Friend.giftSign = append(self.Sql_Friend.giftSign, pInfo.Uid)

			msg.Uid = append(msg.Uid, pInfo.Uid)
			msg.Type = 0
			msg.Value = 1
			msg.Count = self.Sql_Friend.Count
		} else {
			//! 领取
			if self.Sql_Friend.Count >= MAX_FRIEND_POWER {
				return false
			}

			if pInfo.Get != 1 {
				//! 不能领取
				if uid > 0 {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
					return false
				} else {
					return true
				}
			}

			pInfo.Get = 2
			self.Sql_Friend.Count++
			power += 1

			msg.Uid = append(msg.Uid, pInfo.Uid)
			msg.Type = 1
			msg.Value = 2
			msg.Count = self.Sql_Friend.Count
		}
		if uid > 0 {
			return false
		}

		return true
	})

	//修改为友情点
	if power > 0 {
		msg.Item = self.player.AddObjectSimple(ITEM_FRIEND_POINT, power, "领取好友赠送", 0, 0, 0)
	}
	//这个地方不知道哪个是送哪个是收，等测到在改

	if _type == 0 {
		self.player.HandleTask(TASK_TYPE_FRIEND_POINT_COUNT, len(lst), 0, 0)
		self.player.HandleTask(TASK_TYPE_FRIEND_POINT, len(lst), 0, 0)
	}

	self.player.SendMsg("friendpower", HF_JtoB(&msg))

	if len(lst) > 0 {
		if _type == 0 {
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FRIEND_POWER_SEND, len(lst), 0, 0, "好友赠送", 0, 0, self.player)
		} else {
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FRIEND_POWER_GET, len(lst), 0, 0, "好友领取", 0, 0, self.player)
		}
	}

	if len(lst) > 0 {
		GetMasterMgr().FriendPRC.PowerFriend(self.player.GetUid(), lst)
	}
	return
}

func (self *ModFriend) OrderEx(agree int) []*Player {
	lst := make([]*Player, 0)
	uid := self.player.Sql_UserBase.Uid

	self.Sql_Friend.apply.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_Friend)
		if !ok {
			return true
		}

		if pInfo.Uid == uid {
			return true
		}

		friend := GetPlayerMgr().GetPlayer(pInfo.Uid, true)
		if friend == nil {
			return true
		}

		friend.GetModule("friend").(*ModFriend).DelHasapply(uid)
		if agree == 1 {
			//! 互删申请列表
			self.DelHasapply(pInfo.Uid)
			lst = append(lst, friend)

			//! 互删推荐好友
			self.DelCanFriend(pInfo.Uid, false)
			friend.GetModule("friend").(*ModFriend).DelCanFriend(uid, true)

			if self.IsHasBlack(pInfo.Uid) {
				return true
			}
			if friend.GetModule("friend").(*ModFriend).IsHasBlack(uid) {
				return true
			}

			if self.getFriendSize() >= MAX_FRIEND {
				return false
			}

			if friend.GetModule("friend").(*ModFriend).getFriendSize() >= MAX_FRIEND {
				return true
			}
			//! 互相添加好友
			self.AddFriend(friend)
			friend.GetModule("friend").(*ModFriend).AddFriend(self.player)
		}

		return true
	})

	self.Sql_Friend.apply = new(sync.Map)

	return lst
}

func (self *ModFriend) getFriendSize() int {
	size := 0
	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		size += 1
		return true
	})
	return size
}

func (self *ModFriend) getCommendSize() int {
	size := 0
	self.Sql_Friend.commend.Range(func(key, value interface{}) bool {
		size += 1
		return true
	})
	return size
}

//! 处理好友申请
func (self *ModFriend) Order(uid int64, agree int) {
	if uid == 0 {
		var msg S2C_AddFriendMsg
		msg.Cid = "addfriendret"
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
	}

	rel := GetMasterMgr().FriendPRC.AgreeFriend(self.player.GetUid(), uid, agree)
	if rel == nil {
		return
	}
	if self.CheckErr(rel.RetCode) {
		return
	}
	node := make([]*JS_Friend, 0)
	json.Unmarshal([]byte(rel.Data), &node)

	if agree == LOGIC_FALSE {
		for i := 0; i < len(node); i++ {
			self.DelHasapply(node[i].Uid)
		}
	} else if agree == LOGIC_TRUE {
		for i := 0; i < len(node); i++ {
			self.DelHasapply(node[i].Uid)
			nodeSave := NewJsFriendNode(0, 0, node[i].Uid, node[i].Name, node[i].Icon, node[i].Portrait, node[i].Level, node[i].Vip, node[i].Fight, 0)
			self.Sql_Friend.friend.Store(node[i].Uid, nodeSave)
			self.player.HandleTask(FriendCount, 0, 0, 0)
			self.AddInfo("addfriend", node[i])

			nCount := 0
			self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
				pInfo, ok := value.(*JS_FriendNode)
				if ok && pInfo != nil {
					nCount++
				}
				return true
			})

			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FRIEND_ADD, int(node[i].Uid), nCount, 0, "添加好友", 0, 0, self.player)
		}
	}
}

func (self *ModFriend) Black(uid int64) {
	rel := GetMasterMgr().FriendPRC.BlackFriend(self.player.GetUid(), uid)
	if rel == nil {
		return
	}
	if self.CheckErr(rel.RetCode) {
		return
	}
	node := &JS_Friend{}
	json.Unmarshal([]byte(rel.Data), &node)

	if node != nil {
		self.DelHasapply(node.Uid)
		self.delApply(node.Uid)
		self.DelFriend(node.Uid, true, LOGIC_TRUE)
		self.addBlack(node.Uid, node)
	}
}

func (self *ModFriend) Delete(uid int64, _type int) {
	if _type == 0 {
		rel := GetMasterMgr().FriendPRC.DelFriend(self.player.GetUid(), uid)
		if rel == nil {
			return
		}
		if self.CheckErr(rel.RetCode) {
			return
		}
		var uid int64
		json.Unmarshal([]byte(rel.Data), &uid)

		self.DelFriend(uid, true, LOGIC_TRUE)
	} else {
		rel := GetMasterMgr().FriendPRC.BlackOutFriend(self.player.GetUid(), uid)
		if rel == nil {
			return
		}
		if self.CheckErr(rel.RetCode) {
			return
		}
		var uid int64
		json.Unmarshal([]byte(rel.Data), &uid)

		self.deleteBlack(uid)
	}
}

func (self *ModFriend) DeleteBatch(uid []int64, _type int) {
	if _type == 0 {
		//! 好友
		for i := 0; i < len(uid); i++ {
			friend := GetPlayerMgr().GetPlayer(uid[i], true)
			if friend == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
				return
			}
			self.DelFriend(uid[i], false, LOGIC_TRUE)
			friend.GetModule("friend").(*ModFriend).DelFriend(self.player.Sql_UserBase.Uid, true, LOGIC_FALSE)
		}
	} else {
		//! 黑名单
		for i := 0; i < len(uid); i++ {
			self.DelBlack(uid[i])
		}
	}
}

// 这个很耗时
func (self *ModFriend) Rename(loaddb bool) {
	if self.Sql_Friend.friend == nil {
		return
	}
	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_FriendNode)
		if !ok {
			return true
		}

		uid := pInfo.Uid
		player := GetPlayerMgr().GetPlayer(uid, loaddb)
		if player == nil {
			return true
		}
		player.GetModule("friend").(*ModFriend).UpdateFriend(self.player)

		pInfo.Level = player.Sql_UserBase.Level
		pInfo.Fight = player.Sql_UserBase.Fight
		pInfo.Vip = player.Sql_UserBase.Vip

		return true
	})
}

func (self *ModFriend) UpdateFriend(player *Player) {
	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_FriendNode)
		if !ok {
			return true
		}
		//不加判断的话所有好友数据都改掉了。20190509 by zy
		if pInfo.Uid == player.Sql_UserBase.Uid {
			pInfo.Name = player.Sql_UserBase.UName
			pInfo.Fight = player.Sql_UserBase.Fight
			pInfo.Level = player.Sql_UserBase.Level
			pInfo.Vip = player.Sql_UserBase.Vip
			pInfo.Icon = player.Sql_UserBase.IconId
			pInfo.Portrait = player.Sql_UserBase.Portrait
		}
		return true
	})
}

func (self *ModFriend) AddInfo(cid string, info *JS_Friend) {
	var msg S2C_AddFriendMsg
	msg.Cid = cid
	msg.Info = info
	msg.Info.SetOnline()
	self.player.SendMsg(cid, HF_JtoB(&msg))
}

func (self *ModFriend) getFriend() []*JS_FriendNode {
	res := []*JS_FriendNode{}
	//找中心服同步
	rel := GetMasterMgr().FriendPRC.GetFriend(self.player.GetUid())
	if rel != nil {
		friends := make(map[int64]*JS_Friend, 0)
		json.Unmarshal([]byte(rel.Data), &friends)
		if len(friends) > 0 {
			temp := new(sync.Map)
			for _, friend := range friends {
				newNode := new(JS_FriendNode)
				newNode.Uid = friend.Uid
				newNode.Name = friend.Name
				newNode.Icon = friend.Icon
				newNode.Stage = friend.Stage
				newNode.Server = friend.Server
				newNode.Portrait = friend.Portrait
				newNode.Level = friend.Level
				newNode.Vip = friend.Vip
				newNode.Fight = friend.Fight
				newNode.LastUpTime = friend.LastUpTime
				newNode.Online = friend.Online

				pInfo, ok := self.Sql_Friend.friend.Load(newNode.Uid)
				if ok {
					fNode, okNode := pInfo.(*JS_FriendNode)
					if okNode {
						newNode.Gift = fNode.Gift
						newNode.Get = fNode.Get
					}
				}

				//记录下之前的
				res = append(res, newNode)
				temp.Store(newNode.Uid, newNode)
			}
			self.Sql_Friend.friend = temp
			self.player.HandleTask(FriendCount, 0, 0, 0)
		}
	}
	return res
}

func (self *ModFriend) getApply() []*JS_Friend {

	//找中心服同步
	rel := GetMasterMgr().FriendPRC.GetApply(self.player.GetUid())

	apply := make([]*JS_Friend, 0)
	if rel != nil {
		json.Unmarshal([]byte(rel.Data), &apply)
	}

	return apply
}

func (self *ModFriend) getBlack() []*JS_Friend {

	//找中心服同步
	rel := GetMasterMgr().FriendPRC.GetBlack(self.player.GetUid())

	black := make([]*JS_Friend, 0)
	if rel != nil {
		json.Unmarshal([]byte(rel.Data), &black)
	}

	return black
}

func (self *ModFriend) getHasApplyInfo() []*JS_Friend {
	res := []*JS_Friend{}
	self.Sql_Friend.hasapply.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_Friend)
		if !ok {
			return true
		}
		res = append(res, pInfo)
		return true
	})

	if len(res) <= 0 {
		res = []*JS_Friend{}
	}

	return res
}

func (self *ModFriend) getCommend() []*JS_Friend {
	res := []*JS_Friend{}
	self.Sql_Friend.commend.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_Friend)
		if !ok {
			return true
		}
		res = append(res, pInfo)
		return true
	})

	if len(res) <= 0 {
		res = []*JS_Friend{}
	}

	return res
}

func (self *ModFriend) SendInfo() {
	self.player.HandleTask(FriendCount, 0, 0, 0)
	var msg S2C_Friend
	msg.Cid = "friend"
	msg.Friend = self.getFriend()
	msg.Apply = self.getApply()
	msg.Black = self.getBlack()
	msg.Count = self.Sql_Friend.Count
	msg.HireHero = self.GetHireHeroList()

	newList := make([]*HireHero, 0)
	for _, v := range self.Sql_Friend.applyHire {
		if v.OwnPlayer == nil {
			continue
		}
		if TimeServer().Unix() >= v.ReSetTime {
			continue
		}
		newList = append(newList, v)
	}
	self.Sql_Friend.applyHire = newList

	msg.ApplyHire = self.Sql_Friend.applyHire
	msg.HireTime = self.Sql_Friend.HireTime
	//看好友的
	msg.HireList = self.GetHireList()
	msg.SelfList = self.getSelfList()

	msg.HeroSetInfo = GetOfflineInfoMgr().GetBaseHero(self.player)
	msg.HireUseSign = self.Sql_Friend.useSign
	msg.GiftSign = self.Sql_Friend.giftSign

	self.player.SendMsg("friend", HF_JtoB(&msg))
}

func (self *ModFriend) UpdateHireList(body []byte) {
	var msg S2C_UpdateHireList
	msg.Cid = "updatehirelist"
	msg.HireList = self.GetHireList()
	self.player.SendMsg("updatehirelist", HF_JtoB(&msg))
}

func (self *ModFriend) UpdateHireUseSign(body []byte) {
	var msg S2C_UpdateHireList
	msg.Cid = "updatehirelist"
	msg.HireList = self.GetHireList()
	self.player.SendMsg("updatehirelist", HF_JtoB(&msg))
}

func (self *ModFriend) SendCommend() {
	var msg S2C_FriendCommend
	msg.Cid = "friendcommend"
	msg.Commend = self.getCommend()
	for i := 0; i < len(msg.Commend); i++ {
		msg.Commend[i].SetOnline()
	}
	self.player.SendMsg("friendcommend", HF_JtoB(&msg))
}

func (self *ModFriend) DelInfo(cid string, uid int64) {
	var msg S2C_DelFriendMsg
	msg.Cid = cid
	msg.Uid = uid
	self.player.SendMsg(cid, HF_JtoB(&msg))
}

func (self *ModFriend) SendOnline(online int) {
	var msg S2C_FriendOnline
	msg.Cid = "friendonline"
	msg.Uid = self.player.Sql_UserBase.Uid
	msg.Online = online
	_msg := HF_JtoB(&msg)

	friends := self.getFriend()
	for _, pInfo := range friends {
		player := GetPlayerMgr().GetPlayer(pInfo.Uid, false)
		if player == nil || player.SessionObj == nil {
			continue
		}
		player.SendMsg("friendonline", _msg)
	}
}

func (self *San_Friend) Encode() {
	self.Friend = HF_JtoA(&self.friend)
	self.Apply = HF_JtoA(&self.apply)
	self.Black = HF_JtoA(&self.black)
	self.Commend = HF_JtoA(&self.commend)
	self.Hasapply = HF_JtoA(&self.hasapply)
	self.ApplyHire = HF_JtoA(&self.applyHire)
	self.HireHero = HF_JtoA(&self.hireHero)
}

//检查
/*
func (self *ModFriend) CheckHireHero(player *Player) {

	//过期的话初始化数据
	if self.Sql_Friend.HireTime < TimeServer().Unix() {
		self.Sql_Friend.hireHero = make([]*FriendHero, 0)
		self.Sql_Friend.applyHire = make([]*HireHero, 0)
		self.Sql_Friend.HireTime = HF_GetNextWeekStart()
	} else {
		newApply := make([]*HireHero, 0)
		//检查自己的申请在离线时候是否通过从而影响租用英雄
		for _, v := range self.Sql_Friend.applyHire {
			//更新自己数据
			info := GetHireHeroInfoMgr().GetHireInfo(v.OwnPlayer.Uid, v.HeroKeyId)
			if info == nil {
				continue
			}
			if v.HirePlayer == nil || v.HirePlayer.Uid == 0 {
				newApply = append(newApply, v)
			} else if v.HirePlayer.Uid == self.player.Sql_UserBase.Uid {
				self.AddHireHero(v)
			}
		}
		self.Sql_Friend.applyHire = newApply
	}
	//改名同步
	for _, v := range self.Sql_Friend.applyHire {
		if v.OwnPlayer != nil {
			v.OwnPlayer.Name = GetOfflineInfoMgr().GetName(v.OwnPlayer.Uid)
			v.OwnPlayer.IconId = GetOfflineInfoMgr().GetIconId(v.OwnPlayer.Uid)
			v.OwnPlayer.Portrait = GetOfflineInfoMgr().GetPortrait(v.OwnPlayer.Uid)
		}
	}
	for _, v := range self.Sql_Friend.hireHero {
		if v.OwnPlayer != nil {
			v.OwnPlayer.Name = GetOfflineInfoMgr().GetName(v.OwnPlayer.Uid)
			v.OwnPlayer.IconId = GetOfflineInfoMgr().GetIconId(v.OwnPlayer.Uid)
			v.OwnPlayer.Portrait = GetOfflineInfoMgr().GetPortrait(v.OwnPlayer.Uid)
		}
	}
}
*/

func (self *ModFriend) AddHireHero(hireHero *HireHero) *FriendHero {
	if len(self.Sql_Friend.hireHero) >= HIRE_HERO_MAX {
		return nil
	}
	//看是否已经生成
	isNeedAdd := true
	for _, v := range self.Sql_Friend.hireHero {
		if v.OwnUid == hireHero.OwnPlayer.Uid && v.KeyId == hireHero.HeroKeyId {
			isNeedAdd = false
			break
		}
	}

	if isNeedAdd {
		hero := self.NewHeroForHire(hireHero)
		if hero != nil {
			friendHero := new(FriendHero)
			friendHero.OwnUid = hireHero.OwnPlayer.Uid
			friendHero.OwnPlayer = hireHero.OwnPlayer
			friendHero.KeyId = hireHero.HeroKeyId
			friendHero.Hero = hero
			friendHero.SelfKeyId = self.player.GetModule("hero").(*ModHero).MaxKey()

			self.Sql_Friend.hireHero = append(self.Sql_Friend.hireHero, friendHero)
			return friendHero
		}
	}
	return nil
}

func (self *ModFriend) RemoveApplyHire(hireHero *HireHero) *FriendHero {

	newHire := make([]*HireHero, 0)
	for _, v := range self.Sql_Friend.applyHire {
		if v.OwnPlayer == nil {
			continue
		}
		if TimeServer().Unix() >= v.ReSetTime {
			continue
		}
		if v.OwnPlayer.Uid != hireHero.OwnPlayer.Uid || v.HeroKeyId != hireHero.HeroKeyId {
			newHire = append(newHire, v)
		}
	}
	self.Sql_Friend.applyHire = newHire

	var msgRel S2C_HireApply
	msgRel.Cid = "hireapply"
	msgRel.ApplyHero = self.Sql_Friend.applyHire
	msgRel.HireHero = self.Sql_Friend.hireHero
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return nil
}

func (self *ModFriend) DeleteHire(hireHero *HireHero) {

	if hireHero == nil {
		return
	}
	var msgRel S2C_DeleteHireList
	msgRel.Cid = "detelehirelist"
	msgRel.DeleteHireList = append(msgRel.DeleteHireList, hireHero)
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFriend) NewHeroForHire(hireHero *HireHero) *NewHero {

	configHero := GetCsvMgr().GetHeroMapConfig(hireHero.HeroId, hireHero.HeroQuality)
	if configHero == nil {
		return nil
	}

	heroInfo := new(NewHero)
	heroInfo.HeroId = hireHero.HeroId
	//继承品质
	heroInfo.StarItem = &StarItem{}
	heroInfo.StarItem.UpStar = hireHero.HeroQuality
	heroLv := GetOfflineInfoMgr().GetNewHeroLv(self.player.Sql_UserBase.Uid)
	if heroLv <= 0 {
		heroLv = 1
	}
	heroInfo.LvUp(heroLv)
	//生成神器
	pItem := &ArtifactEquip{}
	pItem.Id = hireHero.HeroArtifactId
	pItem.Lv = hireHero.HeroArtifactLv
	config, ok := GetCsvMgr().ArtifactEquipConfigMap[pItem.Id]
	if ok {
		for i := 0; i < len(config.BaseTypes); i++ {
			if config.BaseTypes[i] > 0 {
				attr := new(AttrInfo)
				attr.AttrId = i + 1
				pItem.AttrInfo = append(pItem.AttrInfo, attr)
			}
		}

		pItem.CalAttr()
		heroInfo.ArtifactEquipIds = append(heroInfo.ArtifactEquipIds, pItem)
	}
	//生成专属
	heroInfo.ExclusiveEquip = &ExclusiveEquip{}
	heroInfo.ExclusiveEquip.Lv = hireHero.HeroExclusiveLv
	heroInfo.ExclusiveEquip.UnLock = hireHero.HeroExclusiveUnLock

	for _, v := range GetCsvMgr().ExclusiveEquipConfigMap {
		if v.HeroId == heroInfo.HeroId {
			heroInfo.ExclusiveEquip.Id = v.Id
			for i := 0; i < len(v.BaseType); i++ {
				if v.BaseType[i] > 0 {
					attr := new(AttrInfo)
					attr.AttrId = i + 1
					heroInfo.ExclusiveEquip.AttrInfo = append(heroInfo.ExclusiveEquip.AttrInfo, attr)
				}
			}
		}
	}
	if hireHero.Talent != nil {
		heroInfo.StageTalent = new(StageTalent)
		HF_DeepCopy(heroInfo.StageTalent, hireHero.Talent)
	}

	//生成装备,读取战力
	fight := self.player.GetModule("crystal").(*ModResonanceCrystal).GetPriestsEquipFight()
	configEquip := GetCsvMgr().GetHireEquipConfig(configHero.AttackType, fight/100, ROBOT_EQUIP_SUBTYPE_HIRE)
	if configEquip != nil {
		for i := 0; i < len(configEquip.Equip); i++ {
			equip := &Equip{}
			equip.Id = configEquip.Equip[i]
			equip.Lv = configEquip.Strengthen[i]
			//config, ok := GetCsvMgr().EquipConfigMap[equip.Id]
			//if ok {
			//生成属性
			/*
				for i := 0; i < len(config.BaseTypes); i++ {
					if config.BaseTypes[i] > 0 {
						attr := new(AttrInfo)
						attr.AttrId = i + 1
						equip.AttrInfo = append(equip.AttrInfo, attr)
					}
				}
				//计算属性
				for _, v := range equip.AttrInfo {
					if v.AttrId <= 0 || v.AttrId > len(config.BaseTypes) {
						continue
					}
					index := v.AttrId - 1
					rate := PER_BIT

					if equip.Lv > 0 {
						configLvUp := GetCsvMgr().GetEquipStrengthenConfig(config.EquipAttackType, config.EquipPosition, config.Quality, equip.Lv)
						if configLvUp != nil {
							rate += configLvUp.Vaual[index]
						}
					}

					if configHero != nil && configHero.AttackType == config.EquipAttackType {
						rate += config.CampExtAdd
					}
					v.AttrType = config.BaseTypes[index]
					v.AttrValue = config.BaseValues[index] * int64(rate) / PER_BIT
				}

			*/
			//}
			heroInfo.EquipIds = append(heroInfo.EquipIds, equip)
		}
	}
	return heroInfo
}

func (self *ModFriend) GetLookBaseHero(playerData *RPC_PlayerData_Req) []*NewHero {

	data := make([]*NewHero, 0)

	for i := 0; i < len(playerData.Heros); i++ {
		if playerData.Heros[i] == nil {
			//这个正常情况不会出现
			continue
		}

		configHero := GetCsvMgr().GetHeroMapConfig(playerData.Heros[i].HeroId, playerData.Heros[i].Star)
		if configHero == nil {
			continue
		}

		heroInfo := new(NewHero)
		heroInfo.HeroId = playerData.Heros[i].HeroId
		heroInfo.HeroLv = playerData.Heros[i].Level
		heroInfo.Skin = playerData.Heros[i].Skin
		heroInfo.Attr = playerData.Heros[i].Attr
		//继承品质
		heroInfo.StarItem = &StarItem{}
		heroInfo.StarItem.UpStar = playerData.Heros[i].Star
		//生成神器
		pItem := &ArtifactEquip{}
		pItem.Id = playerData.Heros[i].ArtifactId
		pItem.Lv = playerData.Heros[i].ArtifactLv
		config, ok := GetCsvMgr().ArtifactEquipConfigMap[pItem.Id]
		if ok {
			for i := 0; i < len(config.BaseTypes); i++ {
				if config.BaseTypes[i] > 0 {
					attr := new(AttrInfo)
					attr.AttrId = i + 1
					pItem.AttrInfo = append(pItem.AttrInfo, attr)
				}
			}

			pItem.CalAttr()
			heroInfo.ArtifactEquipIds = append(heroInfo.ArtifactEquipIds, pItem)
		}
		//生成专属
		heroInfo.ExclusiveEquip = &ExclusiveEquip{}
		heroInfo.ExclusiveEquip.Lv = playerData.Heros[i].ExclusiveLv
		heroInfo.ExclusiveEquip.UnLock = playerData.Heros[i].ExclusiveUnLock
		for _, v := range GetCsvMgr().ExclusiveEquipConfigMap {
			if v.HeroId == heroInfo.HeroId {
				heroInfo.ExclusiveEquip.Id = v.Id
				for i := 0; i < len(v.BaseType); i++ {
					if v.BaseType[i] > 0 {
						attr := new(AttrInfo)
						attr.AttrId = i + 1
						heroInfo.ExclusiveEquip.AttrInfo = append(heroInfo.ExclusiveEquip.AttrInfo, attr)
					}
				}
			}
		}
		if playerData.Heros[i].Talent != nil {
			heroInfo.StageTalent = new(StageTalent)
			HF_DeepCopy(heroInfo.StageTalent, playerData.Heros[i].Talent)
		}

		//生成装备,读取战力
		for j := 0; j < len(playerData.Equips[i]); j++ {
			equip := &Equip{}
			equip.Id = playerData.Equips[i][j].ItemId
			equip.Lv = playerData.Equips[i][j].Level
			heroInfo.EquipIds = append(heroInfo.EquipIds, equip)
		}
		data = append(data, heroInfo)
	}

	return data
}
func (self *ModFriend) CheckBaseHero(body []byte) {
	info := GetOfflineInfoMgr().GetBaseHero(self.player)
	if info != nil {
		isempty := true
		for _, v := range info {
			if v != nil {
				isempty = false
				break
			}
		}
		if !isempty {
			return
		}
	}

	var msg C2S_BaseHeroSet
	heros := self.player.GetModule("team").(*ModTeam).getTeamPos(TEAMTYPE_DEFAULT)
	for i := 0; i < MAX_FIGHT_POS; i++ {
		if i >= len(heros.FightPos) {
			msg.FightPos = append(msg.FightPos, 0)
		} else {
			msg.FightPos = append(msg.FightPos, heros.FightPos[i])
		}
	}
	self.BaseHeroSet(HF_JtoB(msg))
}

func (self *ModFriend) BaseHeroSet(body []byte) {

	var msg C2S_BaseHeroSet
	json.Unmarshal(body, &msg)

	info := GetOfflineInfoMgr().GetBaseHero(self.player)
	if info == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_DATA_ERROR"))
		return
	}

	if len(msg.FightPos) != MAX_FIGHT_POS || len(info) != MAX_FIGHT_POS {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_MSG_ERROR"))
		return
	}

	for i := 0; i < len(msg.FightPos); i++ {
		hero := self.player.getHero(msg.FightPos[i])
		if hero == nil {
			continue
		}
		newHero := GetRobotMgr().HeroToNewHero(self.player.Sql_UserBase.Uid, hero)
		info[i] = newHero
	}
	self.player.NoticeCenterBaseInfo()
	var msgRel S2C_BaseHeroSet
	msgRel.Cid = "baseheroset"
	msgRel.BaseHeroSet = info
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFriend) UpdateHeroSet(heroKeyId int) {

	info := GetOfflineInfoMgr().GetBaseHero(self.player)
	if info == nil {
		return
	}

	for i := 0; i < len(info); i++ {
		if info[i] != nil && info[i].HeroKeyId == heroKeyId {
			hero := self.player.getHero(heroKeyId)
			if hero == nil {
				continue
			}
			newHero := GetRobotMgr().HeroToNewHero(self.player.Sql_UserBase.Uid, hero)
			info[i] = newHero
			return
		}
	}
}

func (self *ModFriend) BaseHeroUpdata(heroKeyId int) {
	info := GetOfflineInfoMgr().GetBaseHero(self.player)
	if info == nil {
		return
	}

	for _, v := range info {
		if v != nil && v.HeroKeyId == heroKeyId {
			v = nil
			var msgRel S2C_BaseHeroSet
			msgRel.Cid = "baseheroset"
			msgRel.BaseHeroSet = info
			self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
			return
		}
	}
}

func (self *ModFriend) HireApply(body []byte) {
	var msg C2S_HireApply
	json.Unmarshal(body, &msg)

	//这个逻辑改动大，拎出来改，以后重写
	if msg.HireAction == ACTION_HIRE_LOSE {
		self.HireApplyLose(&msg)
		return
	} else if msg.HireAction == ACTION_APPLY {
		self.HireApplyApply(&msg)
		return
	} else if msg.HireAction == ACTION_APPLY_CANCEL {
		self.HireApplyCancel(&msg)
		return
	}
}

//中心服再次修改 这个功能是玩家退还已经成功租借到的英雄
func (self *ModFriend) HireApplyLose(msg *C2S_HireApply) {

	rel := GetMasterMgr().FriendPRC.HireApplyLose(self.player.GetUid(), msg.HireOwnUid, msg.HireHeroKeyId)
	if rel == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HIRE_INFO_NOT_EXIST"))
		return
	}

	//if self.CheckErr(rel.RetCode) {
	//	return
	//}
	hire := &HireHero{}
	json.Unmarshal([]byte(rel.Data), &hire)

	if hire != nil {
		newHire := make([]*FriendHero, 0)
		for _, v := range self.Sql_Friend.hireHero {
			if v == nil {
				continue
			}

			if (v.OwnPlayer == nil || msg.HireOwnUid != v.OwnPlayer.Uid) || v.KeyId != msg.HireHeroKeyId {
				newHire = append(newHire, v)
			}
		}
		self.Sql_Friend.hireHero = newHire

		newApply := make([]*HireHero, 0)
		for _, v := range self.Sql_Friend.applyHire {
			if v.OwnPlayer == nil {
				continue
			}
			if TimeServer().Unix() >= v.ReSetTime {
				continue
			}
			if v.OwnPlayer.Uid != msg.HireOwnUid || v.HeroKeyId != msg.HireHeroKeyId {
				newApply = append(newApply, v)
			}
		}
		self.Sql_Friend.applyHire = newApply

		var msgRel S2C_HireApply
		msgRel.Cid = "hireapply"
		msgRel.ApplyHero = self.Sql_Friend.applyHire
		msgRel.HireHero = self.Sql_Friend.hireHero
		self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	}
}

func (self *ModFriend) HireApplyApply(msg *C2S_HireApply) {

	rel := GetMasterMgr().FriendPRC.HireApplyApply(self.player.GetUid(), msg.HireOwnUid, msg.HireHeroKeyId)
	if rel == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HIRE_INFO_NOT_EXIST"))
		return
	}
	if self.CheckErr(rel.RetCode) {
		return
	}
	hire := &HireHero{}
	json.Unmarshal([]byte(rel.Data), &hire)
	if hire == nil || hire.OwnPlayer == nil {
		return
	}

	isApply := false
	for _, v := range self.Sql_Friend.applyHire {
		if v.OwnPlayer == nil {
			continue
		}
		if TimeServer().Unix() >= v.ReSetTime {
			continue
		}
		if v.OwnPlayer.Uid == hire.OwnPlayer.Uid && v.HeroKeyId == hire.HeroKeyId {
			isApply = true
			break
		}
	}

	if !isApply {
		self.Sql_Friend.applyHire = append(self.Sql_Friend.applyHire, hire)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_FRIEND_HIRE_APPLY, len(self.Sql_Friend.applyHire), 0, 0, "好友佣兵申请", 0, 0, self.player)
	}

	var msgRel S2C_HireApply
	msgRel.Cid = "hireapply"
	msgRel.ApplyHero = self.Sql_Friend.applyHire
	msgRel.HireHero = self.Sql_Friend.hireHero
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFriend) HireApplyCancel(msg *C2S_HireApply) {

	rel := GetMasterMgr().FriendPRC.HireApplyCancel(self.player.GetUid(), msg.HireOwnUid, msg.HireHeroKeyId)
	if rel == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HIRE_INFO_NOT_EXIST"))
		return
	}
	if self.CheckErr(rel.RetCode) {
		return
	}
	hire := &HireHero{}
	json.Unmarshal([]byte(rel.Data), &hire)

	newApply := make([]*HireHero, 0)
	for _, v := range self.Sql_Friend.applyHire {
		if v.OwnPlayer == nil {
			continue
		}
		if TimeServer().Unix() >= v.ReSetTime {
			continue
		}
		if v.OwnPlayer.Uid != hire.OwnPlayer.Uid || v.HeroKeyId != hire.HeroKeyId {
			newApply = append(newApply, v)
		}
	}
	self.Sql_Friend.applyHire = newApply

	var msgRel S2C_HireApply
	msgRel.Cid = "hireapply"
	msgRel.ApplyHero = self.Sql_Friend.applyHire
	msgRel.HireHero = self.Sql_Friend.hireHero
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFriend) HireStateSet(body []byte) {

	var msg C2S_HireStateSet
	json.Unmarshal(body, &msg)

	if msg.HireState == LOGIC_TRUE {
		self.HireApplyAgree(&msg)
		return
	} else if msg.HireState == LOGIC_FALSE {
		self.HireApplyRefuse(&msg)
		return
	}
}

func (self *ModFriend) HireApplyAgree(msg *C2S_HireStateSet) {

	rel := GetMasterMgr().FriendPRC.HireApplyAgree(self.player.GetUid(), msg.HireUid, msg.HireHeroKeyId)
	if rel == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HIRE_INFO_NOT_EXIST"))
		return
	}
	if self.CheckErr(rel.RetCode) {
		return
	}
	hire := &HireHero{}
	json.Unmarshal([]byte(rel.Data), &hire)
	if hire != nil {
		self.HireStateUpdate(hire)
	}
}

func (self *ModFriend) HireApplyRefuse(msg *C2S_HireStateSet) {

	rel := GetMasterMgr().FriendPRC.HireApplyRefuse(self.player.GetUid(), msg.HireUid, msg.HireHeroKeyId)
	if rel == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HIRE_INFO_NOT_EXIST"))
		return
	}
	//提示错误
	if self.CheckErr(rel.RetCode) {
		return
	}
	hire := &HireHero{}
	json.Unmarshal([]byte(rel.Data), &hire)
	if hire != nil {
		self.HireStateUpdate(hire)
	}
}

func (self *ModFriend) HireStateUpdate(hire *HireHero) {
	if hire == nil || hire.OwnPlayer == nil {
		return
	}
	newList := make([]*HireHero, 0)
	for _, v := range self.Sql_Friend.applyHire {
		if v.OwnPlayer == nil {
			continue
		}
		if TimeServer().Unix() >= v.ReSetTime {
			continue
		}
		if v.OwnPlayer.Uid == hire.OwnPlayer.Uid && v.HeroKeyId == hire.HeroKeyId {
			isFind := false
			for _, info := range hire.ApplyPlayer {
				if info.Uid == self.Sql_Friend.Uid {
					isFind = true
					break
				}
			}
			if !isFind {
				continue
			}
			newList = append(newList, hire)
		} else {
			newList = append(newList, v)
		}
	}
	self.Sql_Friend.applyHire = newList
	//同步自己的申请列表
	var msgRel S2C_HireStateUpdate
	msgRel.Cid = "hirestateupdate"
	msgRel.HireHero = hire
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFriend) HireStateSetAll(body []byte) {
	rel := GetMasterMgr().FriendPRC.HireApplyAgreeAll(self.player.GetUid())
	if rel == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HIRE_INFO_NOT_EXIST"))
		return
	}
	hire := make([]*HireHero, 0)
	json.Unmarshal([]byte(rel.Data), &hire)
	for _, v := range hire {
		self.HireStateUpdate(v)
	}
}

func (self *ModFriend) GetPlayerTeam(body []byte) {

	var msg C2S_GetPlayerTeam
	json.Unmarshal(body, &msg)

	var msgRel S2C_GetPlayerTeam
	msgRel.Cid = "getplayerteam"
	switch msg.TeamType {
	case TEAMTYPE_ARENA_2:
		rel := GetMasterMgr().PlayerPRC.GetPlayerArena(msg.FriendUid)
		fightInfo := new(JS_FightInfo)
		if rel != nil {
			json.Unmarshal([]byte(rel.Data), &fightInfo)
		}
		msgRel.FightInfo = fightInfo
	default:
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_FRIEND_MSG_ERROR"))
		return
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

}

func (self *ModFriend) GetFriendInfo(body []byte) {

	var msg S2C_FriendUpdate
	msg.Cid = "getfriendinfo"
	msg.Friend = self.getFriend()
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

//通知雇佣更新
func (self *ModFriend) NoticeFriend(hireHero *HireHero) {

	self.Sql_Friend.friend.Range(func(key, value interface{}) bool {
		pInfo, ok := value.(*JS_FriendNode)
		if ok {
			friend := GetPlayerMgr().GetPlayer(pInfo.Uid, false)
			if friend != nil {
				var msgRel S2C_HireStateUpdate
				msgRel.Cid = "hirestateupdate"
				msgRel.HireHero = hireHero
				friend.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
			}
		}
		return true
	})
}

func (self *ModFriend) CalHireInitLv() int {
	heroes := self.player.getHeroes()
	heroList := make([]int, 0)
	for _, v := range heroes {
		heroList = append(heroList, v.HeroLv)
	}

	sort.Ints(heroList)

	allLv := 0
	for i := 0; i < len(heroList) && i < 3; i++ {
		allLv += heroList[i]
	}

	return allLv / 3
}

func (self *ModFriend) GetHireList() map[int][]*HireHero {

	//找中心服同步
	rel := GetMasterMgr().FriendPRC.GetHireList(self.player.GetUid())

	hire := make(map[int][]*HireHero, 0)
	if rel != nil {
		json.Unmarshal([]byte(rel.Data), &hire)
	}

	return hire
}

func (self *ModFriend) GetHireHeroList() []*FriendHero {

	if self.Sql_Friend.HireTime < TimeServer().Unix() {
		self.Sql_Friend.hireHero = make([]*FriendHero, 0)
		self.Sql_Friend.HireTime = HF_GetNextWeekStart()
		self.Sql_Friend.useSign = make([]int, 0)
		for i := 0; i < HIRE_MOD_CAMP_END; i++ {
			self.Sql_Friend.useSign = append(self.Sql_Friend.useSign, LOGIC_FALSE)
		}
	}
	return self.Sql_Friend.hireHero
}

func (self *ModFriend) getSelfList() map[int]*HireHero {

	//找中心服同步
	rel := GetMasterMgr().FriendPRC.GetSelfList(self.player.GetUid())

	hire := make(map[int]*HireHero, 0)
	if rel != nil {
		json.Unmarshal([]byte(rel.Data), &hire)
	}

	return hire
}

//设置佣兵使用
func (self *ModFriend) SetUseSign(index int, state int) {

	if index <= 0 || index >= HIRE_MOD_CAMP_END {
		return
	}

	realIndex := index - 1

	self.Sql_Friend.useSign[realIndex] = state

	var msg S2C_SetUseSign
	msg.Cid = "setusesign"
	msg.HireUseSign = self.Sql_Friend.useSign

	self.player.SendMsg("friend", HF_JtoB(&msg))
}

//获得佣兵
func (self *ModFriend) GetHireHero(selfKeyId int) *NewHero {
	for _, v := range self.Sql_Friend.hireHero {
		if v.SelfKeyId == selfKeyId {
			//更新佣兵信息

			return v.Hero
		}
	}
	return nil
}

func (self *ModFriend) CheckErr(code int) bool {
	switch code {
	case RETCODE_PLAYER_NOT_EXIST:
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_PLAYER_NOT_EXIST"))
		return true
	case RETCODE_FRIEND_NOT_EXIST:
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_FRIEND_NOT_EXIST"))
		return true
	case RETCODE_SELF_FRIEND_FULL:
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_SELF_FRIEND_FULL"))
		return true
	case RETCODE_TARGET_FRIEND_FULL:
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_TARGET_FRIEND_FULL"))
		return true
	case RETCODE_FRIEND_APPLY_ERROR:
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_APPLY_ERROR"))
		return true
	case RETCODE_FRIEND_APPLY_BLACK:
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_IN_BLACK"))
		return true
	case RETCODE_FRIEND_HIRE_ALREADY:
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_HIRE_ALREADY"))
		return true
	case RETCODE_FRIEND_HIRE_MAX:
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_HIRE_MAX"))
		return true
	}
	return false
}

func (self *ModFriend) RobotAction(body []byte) {
	GetMasterMgr().FriendPRC.RobotAction(self.player.GetUid())
	self.SendInfo()
}
