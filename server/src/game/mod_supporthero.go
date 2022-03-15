package game

import (
	"encoding/json"
	//"time"
)

const (
	HERO_SUPPORT_SET    = "herosupport_set"    // 设置
	HERO_SUPPORT_CANCEL = "herosupport_cancel" // 取消设置
	//HERO_SUPPORT_USE        = "herosupportuse"       // 使用
	//HERO_SUPPORT_CANCEL_USE = "herosupportcanceluse" // 取消使用
	HERO_SUPPORT_INFO    = "herosupport_info"   // 获取所有可用英雄信息
	HERO_SUPPORT_MY_HERO = "herosupport_myhero" // 获取自己信息
)

// 英雄支援
type ModSupportHero struct {
	player *Player
}

// 获得配置
func (self *ModSupportHero) OnGetData(player *Player) {
	self.player = player
}

func (self *ModSupportHero) OnGetOtherData() {
	self.CheckHeroType()
}
func (self *ModSupportHero) Decode()         {}
func (self *ModSupportHero) Encode()         {}
func (self *ModSupportHero) OnSave(sql bool) {}
func (self *ModSupportHero) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (self *ModSupportHero) onReg(handlers map[string]func(body []byte)) {
	handlers[HERO_SUPPORT_SET] = self.HeroSupportSet       // 设置
	handlers[HERO_SUPPORT_CANCEL] = self.HeroSupportCancel // 取消设置
	//handlers[HERO_SUPPORT_USE] = self.HeroSupportUse              // 使用
	//handlers[HERO_SUPPORT_CANCEL_USE] = self.HeroSupportCancelUse // 取消使用
	handlers[HERO_SUPPORT_INFO] = self.HeroSupportInfo      // 获取信息
	handlers[HERO_SUPPORT_MY_HERO] = self.HeroSupportMyHero // 获取信息
}

// 获得我的支援英雄列表
func (self *ModSupportHero) HeroSupportMyHero(body []byte) {
	heroList := GetSupportHeroMgr().GetMyHero(self.player.GetUid())

	var backmsg S2C_SupportMyHeroInfo
	backmsg.Cid = HERO_SUPPORT_MY_HERO
	for _, v := range heroList {
		temp := MySupportHero{v.Index, v.HeroKey, v.Type, v.UserUid, v.UserName, v.CDTime}
		backmsg.Info = append(backmsg.Info, &temp)
	}
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 获得我能使用的支援英雄列表
func (self *ModSupportHero) HeroSupportInfo(body []byte) {
	var msg C2S_SupportHeroInfo
	json.Unmarshal(body, &msg)

	// 获得uid列表
	playerList := make(map[int64]int64)
	// 好友列表
	friend := self.player.GetModule("friend").(*ModFriend).getFriend()
	for _, v := range friend {
		_, ok := playerList[v.Uid]
		if !ok {
			playerList[v.Uid] = v.Uid
		}
	}

	// 公会列表
	unionid := self.player.GetUnionId()
	if unionid > 0 {
		uniondata := GetUnionMgr().GetUnion(unionid)
		if uniondata != nil {
			for _, v := range uniondata.member {
				if v.Uid == self.player.GetUid() {
					continue
				}
				_, ok := playerList[v.Uid]
				if !ok {
					if v.Lastlogintime != 0 && TimeServer().Unix()-v.Lastlogintime >= SUPPORT_HERO_END_TIME*DAY_SECS {
						GetSupportHeroMgr().CleanPlayerData(v.Uid)
						continue
					}
					playerList[v.Uid] = v.Uid
				}
			}
		}
	}

	//获得列表
	heroList := GetSupportHeroMgr().GetCanUseHero1(playerList, msg.HeroID)
	var backmsg S2C_SupportHeroInfo
	backmsg.Cid = HERO_SUPPORT_INFO
	if heroList != nil {
		backmsg.Info = heroList
	} else {
		backmsg.Info = []*SupportHero{}
	}

	backmsg.HeroID = msg.HeroID
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

func (self *ModSupportHero) CheckPassLevel(passid int) {
	config, ok := GetCsvMgr().OpenLevelMap[OPEN_LEVEL_HERO_SUPPORT]
	if !ok {
		return
	}
	if passid == config.Passid {
		self.HeroSupportAutoSet()
	}
}

//自动设置支援英雄
func (self *ModSupportHero) HeroSupportAutoSet() {
	uid := self.player.GetUid()
	// 取得支援英雄数据
	data := GetSupportHeroMgr().GetPlayerData(uid, true)
	if nil == data {
		self.player.SendErrInfo("err", "取得支援英雄数据")
		return
	}

	// 超过数量
	if len(data) != 0 {
		self.player.SendErrInfo("err", "超过数量")
		return
	}

	// 英雄排序
	heros := self.player.GetModule("hero").(*ModHero).GetBestFormat4()

	nCount := 0
	for _, t := range heros {
		// 获得英雄
		hero := self.player.getHero(t)
		if nil == hero {
			continue
		}

		// 设置支援英雄
		config := GetCsvMgr().GetHeroMapConfig(hero.getHeroId(), hero.GetStar())
		if config == nil {
			continue
		}

		var msg S2M_SupportHeroAdd
		msg.Uid = uid
		msg.Index = nCount + 1
		msg.HeroKeyId = hero.HeroKeyId
		msg.HeroID = config.HeroId
		msg.HeroLv = hero.HeroLv
		msg.HeroStar = hero.StarItem.UpStar
		msg.Skin = hero.Skin

		ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_ADD, &msg)
		if ret == nil || ret.RetCode != UNION_SUCCESS {
			continue
		}

		hero.UseType[HERO_USE_TYPE_SUPPORT] = 1

		nCount++
		// 个数达到
		if nCount >= SUPPORT_HERO_AUTO_MAX {
			break
		}
	}
	self.HeroSupportMyHero([]byte{})
}

//设置支援英雄
func (self *ModSupportHero) HeroSupportSet(body []byte) {
	var msg C2S_SupportHeroSet
	json.Unmarshal(body, &msg)

	// 判断index
	if msg.Index <= 0 || msg.Index > SUPPORT_HERO_MAX {
		self.player.SendErrInfo("err", "判断index")
		return
	}

	// 获得英雄
	hero := self.player.getHero(msg.HeroKey)
	if nil == hero {
		self.player.SendErrInfo("err", "获得英雄错误")
		return
	}

	// 取得支援英雄数据
	data := GetSupportHeroMgr().GetPlayerData(self.player.GetUid(), true)
	if nil == data {
		self.player.SendErrInfo("err", " 取得支援英雄数据错误")
		return
	}

	// 超过数量
	if len(data) > SUPPORT_HERO_MAX {
		self.player.SendErrInfo("err", "超过数量错误")
		return
	}

	// 判断是否已经设置
	for _, v := range data {
		if v.HeroKey == msg.HeroKey {
			self.player.SendErrInfo("err", "判断是否已经设置key错误")
			return
		}

		if v.Index == msg.Index {
			self.player.SendErrInfo("err", "判断是否已经设置index错误")
			return
		}
	}

	config := GetCsvMgr().GetHeroMapConfig(hero.getHeroId(), hero.GetStar())
	if config == nil {
		return
	}

	var mastermsg S2M_SupportHeroAdd
	mastermsg.Uid = self.player.Sql_UserBase.Uid
	mastermsg.Index = msg.Index
	mastermsg.HeroKeyId = hero.HeroKeyId
	mastermsg.HeroID = config.HeroId
	mastermsg.HeroLv = hero.HeroLv
	mastermsg.HeroStar = hero.StarItem.UpStar
	mastermsg.Skin = hero.Skin

	ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_ADD, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		self.player.SendErrInfo("err", "RPC_SUPPORT_HERO_ADD中心服连接问题")
		return
	}

	//// 设置支援英雄
	//if !GetSupportHeroMgr().AddSupportHero(msg.Index, hero, self.player) {
	//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
	//	return
	//}

	hero.UseType[HERO_USE_TYPE_SUPPORT] = 1

	var backmsg S2C_SupportHeroSet
	backmsg.Cid = HERO_SUPPORT_SET
	backmsg.HeroKey = msg.HeroKey
	backmsg.Index = msg.Index
	backmsg.CDTime = TimeServer().Unix() + SUPPORT_HERO_CD_TIME*HOUR_SECS
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SUPPORT_HERO_SET, hero.HeroId, msg.Index, msg.HeroKey, "设置我的派遣英雄", 0, 0, self.player)
}

//取消支援英雄
func (self *ModSupportHero) HeroSupportCancel(body []byte) {
	var msg C2S_SupportHeroCancel
	json.Unmarshal(body, &msg)

	// 获得英雄
	hero := self.player.getHero(msg.HeroKey)
	if nil == hero {
		self.player.SendErrInfo("err", "取消英雄获得英雄错误")
		return
	}

	// 取得支援英雄数据
	data := GetSupportHeroMgr().GetPlayerData(self.player.GetUid(), false)
	if nil == data {
		self.player.SendErrInfo("err", "取消取得支援英雄数据")
		return
	}

	// 判断是否设置了
	find := false
	array := -1
	for i, v := range data {
		if v.HeroKey == msg.HeroKey {
			find = true
			array = i
			break
		}
	}
	if !find {
		self.player.SendErrInfo("err", "判断是否设置了未找到")
		return
	}

	// cd中
	if data[array].CDTime > TimeServer().Unix() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR")+"下阵cd中")
		return
	}

	var mastermsg S2M_SupportHeroRemove
	mastermsg.Uid = self.player.Sql_UserBase.Uid
	mastermsg.HeroKeyId = hero.HeroKeyId

	ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_REMOVE, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR")+"取消中心服连接问题")
		return
	}

	//// 取消英雄
	//if !GetSupportHeroMgr().RemoveSupportHero(hero.HeroKeyId, self.player.GetUid()) {
	//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
	//	return
	//}

	hero.UseType[HERO_USE_TYPE_SUPPORT] = 0

	var backmsg S2C_SupportHeroCancel
	backmsg.Cid = HERO_SUPPORT_CANCEL
	backmsg.HeroKey = msg.HeroKey
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SUPPORT_HERO_CANCEL, hero.HeroId, array, msg.HeroKey, "解除我的派遣英雄", 0, 0, self.player)
	return
}

////消息使用英雄
//func (self *ModSupportHero) HeroSupportUse(body []byte) {
//	var msg C2S_SupportHeroUse
//	json.Unmarshal(body, &msg)
//	self.HeroUse(msg.Uid, msg.HeroKey, msg.Type)
//}
//
////使用英雄接口 其他模块也要用
//func (self *ModSupportHero) HeroUse(uid int64, nHeroKey int, nType int) bool {
//	data := GetSupportHeroMgr().GetPlayerData(uid)
//	if nil == data {
//		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
//		return false
//	}
//
//	// 不能使用自己的支援英雄
//	if uid == self.player.GetUid() {
//		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
//		return false
//	}
//
//	// 超出最大借出个数
//	nCount := self.GetUseCount(uid)
//	if nCount < 0 || nCount >= SUPPORT_HERO_USE_MAX {
//		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
//		return false
//	}
//
//	find := false
//	for _, v := range data.supportHero {
//		// 英雄id判断
//		if v.HeroKey != nHeroKey {
//			continue
//		}
//
//		find = true
//	}
//
//	if !find {
//		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
//		return false
//	}
//
//	//if GetSupportHeroMgr().UseHero(uid, nHeroKey, self.player, nType) {
//	//	var backmsg S2C_SupportHeroUse
//	//	backmsg.Cid = HERO_SUPPORT_USE
//	//	backmsg.Uid = uid
//	//	backmsg.HeroKey = nHeroKey
//	//	backmsg.Type = nType
//	//
//	//	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
//	//	return true
//	//}
//
//	if nType == HERO_SUPPORT_TYPE_ENTANGLEMENT {
//		self.player.GetModule("entanglement").(*ModEntanglement).
//	}
//
//	return true
//}
//
////消息取消使用英雄
//func (self *ModSupportHero) HeroSupportCancelUse(body []byte) {
//	var msg C2S_SupportHeroCancelUse
//	json.Unmarshal(body, &msg)
//
//	self.CancelUseHero(msg.Uid, msg.HeroKey,msg.Type)
//
//}
//
////取消使用英雄接口 其他模块也要用
//func (self *ModSupportHero) CancelUseHero(uid int64, nHeroKey int,nType int) bool {
//	data := GetSupportHeroMgr().GetPlayerData(uid)
//	if nil == data {
//		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
//		return false
//	}
//
//	// 不能取消自己的支援英雄
//	if uid == self.player.GetUid() {
//		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
//		return false
//	}
//
//	find:= false
//
//	for _, v := range data.supportHero {
//		// 英雄id判断
//		if v.HeroKey != nHeroKey {
//			continue
//		}
//
//		find = true
//	}
//
//	if !find{
//		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
//		return false
//	}
//
//	//if GetSupportHeroMgr().CancelUseHero(uid, nHeroKey, self.player) {
//	//	var backmsg S2C_SupportHeroCancelUse
//	//	backmsg.Cid = HERO_SUPPORT_CANCEL_USE
//	//	backmsg.Uid = uid
//	//	backmsg.HeroKey = nHeroKey
//	//
//	//	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
//	//	return true
//	//}
//
//	if nType == HERO_SUPPORT_TYPE_ENTANGLEMENT {
//		self.player.GetModule("entanglement").(*ModEntanglement).
//	}
//
//	return true
//}
//

// 取得使用的支援个数
func (self *ModSupportHero) GetUseCount(uid int64) int {
	nCount := 0
	nCount += self.player.GetModule("entanglement").(*ModEntanglement).GetUseCount(uid)
	//nCount += self.player.GetModule("reward").(*ModReward).GetUseCount(uid)
	return nCount
}

func (self *ModSupportHero) CheckHeroType() bool {
	data := GetSupportHeroMgr().GetPlayerData(self.player.GetUid(), false)
	if nil == data {
		return false
	}

	for _, v := range self.player.GetModule("hero").(*ModHero).Sql_Hero.info {
		if v.UseType[HERO_USE_TYPE_SUPPORT] == 1 {
			find := false
			for _, t := range data {
				if t.HeroKey == v.HeroKeyId {
					find = true
					break
				}
			}
			if !find {
				v.UseType[HERO_USE_TYPE_SUPPORT] = 0
			}
		}
	}
	return true
}

// 是否是好友或者是同一公会
func (self *ModSupportHero) IsFriendOrUnion(uid int64) bool {
	// 是否是好友
	if !self.player.GetModule("friend").(*ModFriend).IsHasFriend(uid) {
		// 公会id
		unionid := self.player.GetUnionId()
		if unionid > 0 {
			uniondata := GetUnionMgr().GetUnion(unionid)
			if uniondata != nil {
				// 是否是同公会
				for _, v := range uniondata.member {
					if v.Uid == uid {
						// 太久没上线 清理
						if v.Lastlogintime != 0 && TimeServer().Unix()-v.Lastlogintime >= SUPPORT_HERO_END_TIME*DAY_SECS {
							GetSupportHeroMgr().CleanPlayerData(uid)
							return false
						}
						return true
					}
				}
			}
		}
		return false
	} else {
		//
		//if TimeServer().Unix()-v.Lastlogintime >= SUPPORT_HERO_END_TIME*DAY_SECS {
		//	GetSupportHeroMgr().CleanPlayerData(uid)
		//	return false
		//}
	}
	return true
}

// 更新英雄信息
func (self *ModSupportHero) UpdataHero(heroKey int) bool {
	data := GetSupportHeroMgr().GetPlayerData(self.player.GetUid(), false)
	if nil == data {
		//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	// 获得英雄
	hero := self.player.getHero(heroKey)
	if nil == hero {
		var mastermsg S2M_SupportHeroRemove
		mastermsg.Uid = self.player.Sql_UserBase.Uid
		mastermsg.HeroKeyId = heroKey

		ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_REMOVE, &mastermsg)
		if ret == nil || ret.RetCode != UNION_SUCCESS {
			return false
		}
		//// 取消英雄
		//if !GetSupportHeroMgr().RemoveSupportHero(heroKey, self.player.GetUid()) {
		//	return false
		//}
		return false
	}

	if hero.UseType[HERO_USE_TYPE_SUPPORT] != 1 {
		return false
	}

	find := false
	for _, v := range data {
		// 英雄id判断
		if v.HeroKey != heroKey {
			continue
		}

		find = true
	}

	if !find {
		//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}
	GetSupportHeroMgr().UpdateHero(self.player.GetUid(), hero)

	return true
}

func (self *ModSupportHero) GetCanUseHero(heroid int) []*SupportHero {
	// 获得uid列表
	playerList := make(map[int64]int64)
	// 好友列表
	friend := self.player.GetModule("friend").(*ModFriend).getFriend()
	for _, v := range friend {
		_, ok := playerList[v.Uid]
		if !ok {
			playerList[v.Uid] = v.Uid
		}
	}

	// 公会列表
	unionid := self.player.GetUnionId()
	if unionid > 0 {
		uniondata := GetUnionMgr().GetUnion(unionid)
		if uniondata != nil {
			for _, v := range uniondata.member {
				if v.Uid == self.player.GetUid() {
					continue
				}
				_, ok := playerList[v.Uid]
				if !ok {
					if v.Lastlogintime != 0 && TimeServer().Unix()-v.Lastlogintime >= SUPPORT_HERO_END_TIME*DAY_SECS {
						GetSupportHeroMgr().CleanPlayerData(v.Uid)
						continue
					}
					playerList[v.Uid] = v.Uid
				}
			}
		}
	}

	//获得列表
	return GetSupportHeroMgr().GetCanUseHero1(playerList, heroid)
}
