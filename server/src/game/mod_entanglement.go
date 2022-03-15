package game

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	ENTANGLEMENT_SET = "entanglement_set" // 设置
	//ENTANGLEMENT_CANCEL = "entanglementcancel" // 取消设置
	ENTANGLEMENT_SEND_INFO        = "entanglement_send_info"        // 消息
	ENTANGLEMENT_SEND_HERO_DELETE = "entanglement_send_hero_delete" // 删除提示
	ENTANGLEMENT_AUTO_SET         = "entanglement_auto_set"         // 一键设置
)

// 羁绊index
const ENTANGLEMENT_INDEX_MAX = 3

// 羁绊层级状态
type Fate struct {
	Fate  int `json:"fate"`  // 命运id
	State int `json:"state"` // 状态
}

// 羁绊英雄
type FateHero struct {
	Index     int    `json:"index"`     // index
	Uid       int64  `json:"uid"`       // 谁的
	Name      string `json:"name"`      // 谁的
	HeroKey   int    `json:"herokey"`   // 英雄key值
	HeroID    int    `json:"heroid"`    // 英雄id
	HeroStar  int    `json:"herostar"`  // 英雄star
	HeroLevel int    `json:"herolevel"` // 英雄level
	HeroSkin  int    `json:"skin"`      // 英雄皮肤
}

// 羁绊结构
type JS_Entanglement struct {
	Type     int         `json:"type"`     // group
	Fate     []*Fate     `json:"fate"`     // 羁绊层级
	FateHero []*FateHero `json:"fatehero"` // 上阵英雄
}

// 检查属性
func (self *JS_Entanglement) CheckProperty() {
	// 重置
	self.Fate = []*Fate{}
	// 获得配置
	config := GetCsvMgr().GetEntanglementConfig(self.Type)
	if nil == config {
		return
	}
	// 检测英雄
	for i, t := range config.HeroId {
		if t == 0 {
			continue
		}
		// 检测
		find := false
		for _, g := range self.FateHero {
			if g.Index-1 == i && g.HeroID == t {
				find = true
				break
			}
		}
		// 没装备对应的英雄
		if !find {
			return
		}
	}

	// 检测属性层级
	for _, v := range config.Property {
		isDone := true
		for _, k := range self.FateHero {
			if k.HeroStar < v.MinQuality {
				isDone = false
				break
			}
		}
		if !isDone {
			continue
		}
		self.Fate = append(self.Fate, &Fate{v.FateNum, 1})
	}
}

// 属性结构
type JS_EntanglementProperty struct {
	HeroID int   `json:"heroid"` // 英雄id
	FateID []int `json:"fateid"` // 羁绊子条目id
}

// 数据结构
type San_Entanglement struct {
	Uid  int64  // 角色ID
	Info string // 羁绊

	info     []*JS_Entanglement         // 羁绊属性
	property []*JS_EntanglementProperty // 属性 每次上线修改重新算 不需要储存
	DataUpdate
}

// 羁绊
type ModEntanglement struct {
	player           *Player
	Sql_Entanglement San_Entanglement

	hint bool
}

// 获得数据
func (self *ModEntanglement) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_entanglement` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Entanglement, "san_entanglement", self.player.ID)

	if self.Sql_Entanglement.Uid <= 0 {
		self.Sql_Entanglement.Uid = self.player.ID
		self.Sql_Entanglement.info = make([]*JS_Entanglement, 0)
		self.Encode()
		InsertTable("san_entanglement", &self.Sql_Entanglement, 0, true)
		self.Sql_Entanglement.Init("san_entanglement", &self.Sql_Entanglement, true)
	} else {
		self.Decode()
		self.Sql_Entanglement.Init("san_entanglement", &self.Sql_Entanglement, true)
	}

	self.hint = false
}

// 上线
func (self *ModEntanglement) OnGetOtherData() {
	self.CheckAllEntanglementHero()     // 检查数据
	self.CheckAllEntanglementProperty() // 添加属性
}

func (self *ModEntanglement) OnReady() {
	if self.hint {
		self.hint = false
		self.player.SendRet2(ENTANGLEMENT_SEND_HERO_DELETE)

		//self.player.SendErrInfo("err", "失去好友的派遣英雄，羁绊失效")
	}
}

// save
func (self *ModEntanglement) Decode() {
	json.Unmarshal([]byte(self.Sql_Entanglement.Info), &self.Sql_Entanglement.info)
}

// read
func (self *ModEntanglement) Encode() {
	self.Sql_Entanglement.Info = HF_JtoA(self.Sql_Entanglement.info)
}
func (self *ModEntanglement) OnSave(sql bool) {
	self.Encode()
	self.Sql_Entanglement.Update(sql)
}

func (self *ModEntanglement) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (self *ModEntanglement) onReg(handlers map[string]func(body []byte)) {
	handlers[ENTANGLEMENT_SET] = self.EntanglementUse // 设置
	//handlers[ENTANGLEMENT_CANCEL] = self.EntanglementCancel // 取消设置
	handlers[ENTANGLEMENT_SEND_INFO] = self.SendInfo // 设置
	handlers[ENTANGLEMENT_AUTO_SET] = self.AutoSet   // 设置
}

// 获得借用该玩家的数量
func (self *ModEntanglement) GetUseCount(uid int64) int {
	nCount := 0
	for _, v := range self.Sql_Entanglement.info {
		for _, g := range v.FateHero {
			if g.Uid == uid {
				nCount++
			}
		}
	}
	return nCount
}

// 发送信息
func (self *ModEntanglement) SendInfo(body []byte) {
	var backmsg S2C_EntanglementInfo
	backmsg.Cid = ENTANGLEMENT_SEND_INFO
	backmsg.Info = self.Sql_Entanglement.info
	backmsg.Property = self.Sql_Entanglement.property
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 获得羁绊数据信息
func (self *ModEntanglement) GetEntanglementInfo(nType int, create bool) *JS_Entanglement {
	for _, v := range self.Sql_Entanglement.info {
		if nType == v.Type {
			return v
		}
	}

	if create {
		data := JS_Entanglement{}
		data.Type = nType
		self.Sql_Entanglement.info = append(self.Sql_Entanglement.info, &data)
		return &data
	}
	return nil
}

//消息使用英雄
func (self *ModEntanglement) EntanglementUse(body []byte) {
	var msg C2S_EntanglementUse
	json.Unmarshal(body, &msg)

	// 判断index
	if msg.Index <= 0 || msg.Index > ENTANGLEMENT_INDEX_MAX {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	// 获得配置
	config := GetCsvMgr().GetEntanglementConfig(msg.Type)
	if nil == config {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	// 英雄id和英雄星级
	heroID := 0
	heroStar := 0
	heroLevel := 0
	heroSkin := 0
	name := ""
	masterUid := int64(0)
	// 借用别人的英雄
	if self.player.GetUid() != msg.MasterUid {
		// 判断是不是好友或者一个军团
		if !self.player.GetModule("support").(*ModSupportHero).IsFriendOrUnion(msg.MasterUid) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}

		// 寻找被借人的data
		data := GetSupportHeroMgr().GetPlayerData(msg.MasterUid, false)
		if nil == data {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}

		// 判断有没有这个英雄提供援助
		find := false
		index := -1
		for i, v := range data {
			if v.HeroKey == msg.HeroKey {
				find = true
				index = i
				break
			}
		}
		if !find {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}

		// 超出借出数量
		if self.player.GetModule("support").(*ModSupportHero).GetUseCount(msg.MasterUid) >= SUPPORT_HERO_USE_MAX {
			self.player.SendErrInfo("err", "雇佣同一个好友或公会成员英雄数量不得超过3名")
			return
		}

		// 判断该英雄是不是属于该羁绊的
		heroID = data[index].HeroID
		heroStar = data[index].HeroStar
		name = data[index].MasterName
		heroLevel = data[index].HeroLv
		heroSkin = data[index].HeroSkin
		masterUid = data[index].MasterUid
	} else {
		// 获得自己的英雄
		hero := self.player.getHero(msg.HeroKey)
		if nil == hero {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}

		heroID = hero.HeroId
		heroStar = hero.GetStar()
		name = self.player.GetName()
		heroLevel = hero.HeroLv
		heroSkin = hero.Skin
		masterUid = self.player.Sql_UserBase.Uid
	}

	// 判断该英雄是不是属于该羁绊的
	heroFind := false
	if msg.Index-1 < len(config.HeroId) && msg.Index-1 >= 0 {
		if config.HeroId[msg.Index-1] == heroID {
			heroFind = true
		}
	}
	if !heroFind {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	// 判断位置是否被占用
	posFind := false
	info := self.GetEntanglementInfo(msg.Type, true)
	if nil == info {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	index := -1
	for i, v := range info.FateHero {
		if v.Index == msg.Index {
			posFind = true
			index = i
			break
		}
	}
	// 位置被占
	if posFind {
		if info.FateHero[index].Uid == msg.MasterUid && info.FateHero[index].HeroKey == msg.HeroKey {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}
		// 卸下属性
		if !self.EntanglementCancel(msg.Index, msg.Type) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}
	}

	// 满了
	if len(info.FateHero) >= config.HeroNum {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	// 不是我自己的
	if self.player.GetUid() != msg.MasterUid {
		// 简单设置提供者的数据
		timeNow := TimeServer().Unix()
		endtime := timeNow + SUPPORT_HERO_END_TIME*DAY_SECS

		var mastermsg S2M_SupportHeroUse
		mastermsg.Uid = msg.MasterUid
		mastermsg.HeroKeyId = msg.HeroKey
		mastermsg.Useruid = self.player.Sql_UserBase.Uid
		mastermsg.Username = self.player.Sql_UserBase.UName
		mastermsg.Type = HERO_SUPPORT_TYPE_ENTANGLEMENT
		mastermsg.Endtime = endtime

		ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_USE, &mastermsg)
		if ret == nil || ret.RetCode != UNION_SUCCESS {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}

		//if !GetSupportHeroMgr().UseHero(msg.MasterUid, msg.HeroKey, self.player, HERO_SUPPORT_TYPE_ENTANGLEMENT, endtime) {
		//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		//	return
		//}
	}

	// 清理全部属性
	self.DeleteTypeAllProperty(msg.Type)
	info.FateHero = append(info.FateHero, &FateHero{msg.Index, msg.MasterUid, name, msg.HeroKey, heroID, heroStar, heroLevel, heroSkin})
	// 检查属性
	info.CheckProperty()
	// 添加属性
	self.AddTypeAllProperty(msg.Type)

	var backmsg S2C_EntanglementUse
	backmsg.Cid = ENTANGLEMENT_SET
	backmsg.MasterUid = msg.MasterUid
	backmsg.Index = msg.Index
	backmsg.Type = msg.Type
	backmsg.HeroKey = msg.HeroKey
	backmsg.Info = append(backmsg.Info, info)
	backmsg.Property = self.GetHeroPropertyData(heroID)
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ENTANGLEMENT_SET, heroID, msg.Type, int(masterUid), "设置羁绊英雄", 0, len(info.Fate), self.player)
}

//取消使用英雄
func (self *ModEntanglement) EntanglementCancel(nIndex int, nType int) bool {
	// 判断index
	if nIndex <= 0 || nIndex > ENTANGLEMENT_INDEX_MAX {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	// 获得配置
	config := GetCsvMgr().GetEntanglementConfig(nType)
	if nil == config {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	// 获得数据
	info := self.GetEntanglementInfo(nType, false)
	if nil == info {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	// 检查英雄是否被使用
	index := -1
	for i, v := range info.FateHero {
		if v.Index == nIndex {
			index = i
			break
		}
	}
	if index < 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	// 借的英雄
	if info.FateHero[index].Uid != self.player.GetUid() {
		var mastermsg S2M_SupportHeroCancelUse
		mastermsg.Uid = info.FateHero[index].Uid
		mastermsg.HeroKeyId = info.FateHero[index].HeroKey
		mastermsg.Useruid = self.player.Sql_UserBase.Uid
		GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_CANCEL_USE, &mastermsg)

		//// 简单设置提供者的数据
		//GetSupportHeroMgr().CancelUseHero(info.FateHero[index].Uid, info.FateHero[index].HeroKey, self.player)
	}

	// 删除所有属性
	self.DeleteTypeAllProperty(nType)
	info.FateHero = append(info.FateHero[:index], info.FateHero[index+1:]...)
	// 检查属性
	info.CheckProperty()
	// 添加属性
	self.AddTypeAllProperty(nType)
	return true
}

// 删除所有属性
func (self *ModEntanglement) DeleteTypeAllProperty(nType int) {
	info := self.GetEntanglementInfo(nType, false)
	if nil == info {
		//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	for _, v := range info.FateHero {
		data := self.GetHeroPropertyData(v.HeroID)
		if data != nil {
			self.DeleteHeroProperty(data, info.Fate)
		}
	}
}

// 添加所有属性
func (self *ModEntanglement) AddTypeAllProperty(nType int) {
	info := self.GetEntanglementInfo(nType, false)
	if nil == info {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	for _, v := range info.FateHero {
		data := self.GetHeroPropertyData(v.HeroID)
		if data != nil {
			self.AddHeroProperty(data, info.Fate)
		}

		for _, value := range self.player.GetModule("hero").(*ModHero).Sql_Hero.info {
			if value.HeroId == v.HeroID {
				self.player.checkHeroFight(value, ReasonHeroFate)
			}
		}
	}
}

// 删除英雄的属性
func (self *ModEntanglement) DeleteHeroProperty(data *JS_EntanglementProperty, fates []*Fate) {
	for _, v := range fates {
		nLen := len(data.FateID)
		for i := nLen - 1; i >= 0; i-- {
			if v.Fate == data.FateID[i] {
				data.FateID = append(data.FateID[:i], data.FateID[i+1:]...)
			}
		}
	}
	return
}

// 添加英雄属性
func (self *ModEntanglement) AddHeroProperty(data *JS_EntanglementProperty, fates []*Fate) {
	for _, v := range fates {
		data.FateID = append(data.FateID, v.Fate)
	}
}

// 获得属性数据
func (self *ModEntanglement) GetHeroPropertyData(heroID int) *JS_EntanglementProperty {
	for _, v := range self.Sql_Entanglement.property {
		if v.HeroID == heroID {
			return v
		}
	}

	data := JS_EntanglementProperty{heroID, []int{}}
	self.Sql_Entanglement.property = append(self.Sql_Entanglement.property, &data)

	return &data
}

// 添加属性
func (self *ModEntanglement) CheckAllEntanglementProperty() {
	// 获得配置
	configs := GetCsvMgr().EntanglementMapConfig
	self.Sql_Entanglement.property = []*JS_EntanglementProperty{}
	for _, v := range configs {
		info := self.GetEntanglementInfo(v.Group, false)
		if nil == info {
			continue
		}
		// 设置英雄属性
		for _, v := range info.FateHero {
			data := self.GetHeroPropertyData(v.HeroID)
			if data != nil {
				self.AddHeroProperty(data, info.Fate)
			}
		}
	}
}

// 获得所有属性
func (self *ModEntanglement) GetAllProperty(heroID int) map[int]*Attribute {
	ret := map[int]*Attribute{}
	data := self.GetHeroPropertyData(heroID)
	if data != nil {
		for _, v := range data.FateID {
			config, ok := GetCsvMgr().EntanglementConfig[v]
			if ok {
				nLen := len(config.BaseType)
				for i := 0; i < nLen; i++ {
					_, ok := ret[config.BaseType[i]]
					if ok {
						ret[config.BaseType[i]].AttValue += config.BaseValue[i]
					} else {
						ret[config.BaseType[i]] = &Attribute{config.BaseType[i], config.BaseValue[i]}
					}
				}
			}
		}
	}

	return ret
}

// 检查上阵的所有的英雄
func (self *ModEntanglement) CheckAllEntanglementHero() {
	uid := self.player.GetUid()

	// 玩家自己是否离线了七天
	refresh := false
	tll, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.LastUpdTime, time.Local)
	if (TimeServer().Unix() - tll.Unix()) >= SUPPORT_HERO_END_TIME*DAY_SECS {
		refresh = true
	}

	delete := false
	for _, v := range self.Sql_Entanglement.info {
		nLen := len(v.FateHero)
		for i := nLen - 1; i >= 0; i-- {
			// 支援的英雄
			if v.FateHero[i].Uid != uid {
				if refresh {
					v.FateHero = append(v.FateHero[:i], v.FateHero[i+1:]...)
					delete = true
					continue
				}

				// 判断是不是好友或者一个军团
				if !self.player.GetModule("support").(*ModSupportHero).IsFriendOrUnion(v.FateHero[i].Uid) {
					v.FateHero = append(v.FateHero[:i], v.FateHero[i+1:]...)
					delete = true
					continue
				}

				// 寻找被借人的data
				data := GetSupportHeroMgr().GetPlayerData(v.FateHero[i].Uid, false)
				if nil == data {
					v.FateHero = append(v.FateHero[:i], v.FateHero[i+1:]...)
					delete = true
					continue
				}

				// 判断有没有这个英雄提供援助
				find := false
				for _, g := range data {
					if g.HeroKey == v.FateHero[i].HeroKey {
						find = true
						v.FateHero[i].Name = g.MasterName
						v.FateHero[i].HeroStar = g.HeroStar
						v.FateHero[i].HeroLevel = g.HeroLv
						break
					}
				}
				if !find {
					v.FateHero = append(v.FateHero[:i], v.FateHero[i+1:]...)
					delete = true
					continue
				}
			} else {
				hero := self.player.getHero(v.FateHero[i].HeroKey)
				if hero == nil {
					v.FateHero = append(v.FateHero[:i], v.FateHero[i+1:]...)
					continue
				}
				v.FateHero[i].HeroStar = hero.StarItem.UpStar
				v.FateHero[i].HeroLevel = hero.HeroLv
			}
		}
		v.CheckProperty()
	}

	if refresh || delete {
		self.hint = true
	}
}

// 更新英雄信息
func (self *ModEntanglement) DeleteHero(heroKey int) {
	uid := self.player.GetUid()

	for _, v := range self.Sql_Entanglement.info {
		nLen := len(v.FateHero)
		for i := nLen - 1; i >= 0; i-- {
			// 支援的英雄
			if v.FateHero[i].Uid != uid {
				continue
			} else {
				if v.FateHero[i].HeroKey == heroKey {
					self.EntanglementCancel(v.FateHero[i].Index, v.Type)
					break
				}
			}
		}
	}

	self.SendInfo([]byte{})
}
func (self *ModEntanglement) UpdateHero(heroKey int) {
	uid := self.player.GetUid()
	hero := self.player.getHero(heroKey)
	if nil == hero {
		return
	}

	heroid := hero.HeroId

	for _, v := range self.Sql_Entanglement.info {
		nLen := len(v.FateHero)
		for i := nLen - 1; i >= 0; i-- {
			// 支援的英雄
			if v.FateHero[i].Uid != uid {
				continue
			} else {
				if v.FateHero[i].HeroKey == heroKey {
					v.FateHero[i].HeroStar = hero.StarItem.UpStar
					v.FateHero[i].HeroLevel = hero.HeroLv
					v.CheckProperty()
					break
				} else {
					if v.FateHero[i].HeroID == heroid {
						if v.FateHero[i].HeroStar < hero.StarItem.UpStar {
							self.EntanglementUse(HF_JtoB(C2S_EntanglementUse{ENTANGLEMENT_SET, v.FateHero[i].Index, v.Type, uid, hero.HeroKeyId}))
						}

						break
					}
				}
			}
		}
	}

	//self.SendInfo([]byte{})
}

// 检查上阵的所有的英雄
func (self *ModEntanglement) GmAddAllFate() {

	configs := GetCsvMgr().EntanglementMapConfig
	for _, config := range configs {
		// 清理全部属性
		self.DeleteTypeAllProperty(config.Group)
		info := self.GetEntanglementInfo(config.Group, true)
		if nil == info {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}

		info.FateHero = []*FateHero{}

		for i, v := range config.HeroId {
			info.FateHero = append(info.FateHero, &FateHero{i + 1, self.player.GetUid(), self.player.GetName(), 0, v, 7, 70, 0})
		}
		// 检查属性
		info.CheckProperty()

		// 添加属性
		self.AddTypeAllProperty(config.Group)
	}
}

// 修改名字
func (self *ModEntanglement) Rename() bool {
	uid := self.player.GetUid()
	for _, v := range self.Sql_Entanglement.info {
		for _, hero := range v.FateHero {
			if hero.Uid == uid {
				hero.Name = self.player.GetName()
			}
		}
	}
	self.SendInfo([]byte{})
	return true
}

// 发送信息
func (self *ModEntanglement) AutoSet(body []byte) {
	var msg C2S_EntanglementAutoUse
	json.Unmarshal(body, &msg)

	// 获得配置
	config := GetCsvMgr().GetEntanglementConfig(msg.Type)
	if nil == config {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	info := self.GetEntanglementInfo(msg.Type, true)
	if nil == info {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	useheroUids := []int64{}
	useIndex := []int{}
	useheroKeys := []int{}
	//heros := []int{}
	change := false
	// 循环该类的英雄
	for i, heroid := range config.HeroId {
		// 用来获得阵位上英雄的信息
		useheroStar := 0       // 阵位上英雄星级
		useheroKey := 0        // 阵位上英雄key
		useheroUid := int64(0) // 阵位上英雄主人uid
		index := i + 1

		// 是否需要卸下属性
		isUse := false
		// 判断当前的阵位上是否有英雄
		for _, v := range info.FateHero {
			if v.Index == index {
				useheroStar = v.HeroStar
				useheroKey = v.HeroKey
				useheroUid = v.Uid
				isUse = true
				break
			}
		}

		// 是否用的自己英雄
		usemine := false
		myUid := self.player.GetUid()
		if useheroUid == myUid {
			usemine = true
		}
		// 是否有更高的
		haveHiger := false
		for _, hero := range self.player.GetModule("hero").(*ModHero).Sql_Hero.info {
			// 如果用自己的 判断是不是同一个key
			if usemine && useheroKey == hero.HeroKeyId {
				continue
			}

			if hero.HeroId != heroid {
				continue
			}

			if useheroStar >= hero.StarItem.UpStar {
				continue
			}

			useheroKey = hero.HeroKeyId
			useheroStar = hero.StarItem.UpStar
			useheroUid = myUid
			haveHiger = true
		}

		data := self.player.GetModule("support").(*ModSupportHero).GetCanUseHero(heroid)
		if len(data) <= 0 && !haveHiger {
			continue
		}

		// 找好友的借出英雄
		for _, friendhero := range data {
			if friendhero.HeroStar <= useheroStar {
				continue
			}

			// 超出借出数量
			if self.player.GetModule("support").(*ModSupportHero).GetUseCount(friendhero.MasterUid) >= SUPPORT_HERO_USE_MAX {
				continue
			}

			useheroKey = friendhero.HeroKey
			useheroStar = friendhero.HeroStar
			useheroUid = friendhero.MasterUid
			haveHiger = true
		}

		// 没有更新的英雄 跳过一轮
		if !haveHiger {
			continue
		}

		usemine = false
		if useheroUid == myUid {
			usemine = true
		}

		if isUse {
			// 卸下属性
			if !self.EntanglementCancel(index, msg.Type) {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
				return
			}
		}

		change = true
		// 清理全部属性
		self.DeleteTypeAllProperty(msg.Type)

		if usemine {
			usehero := self.player.getHero(useheroKey)
			if usehero == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
				return
			}
			info.FateHero = append(info.FateHero, &FateHero{
				index,
				useheroUid,
				self.player.GetName(),
				useheroKey,
				heroid,
				usehero.StarItem.UpStar,
				usehero.HeroLv,
				usehero.Skin})
		} else {
			// 寻找被借人的data
			data := GetSupportHeroMgr().GetPlayerData(useheroUid, false)
			if nil == data {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
				return
			}

			var usehero *SupportHero = nil
			for _, datahero := range data {
				if datahero.HeroKey == useheroKey && datahero.HeroStar == useheroStar {
					usehero = datahero
					break
				}
			}
			if usehero == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
				return
			}

			// 简单设置提供者的数据
			timeNow := TimeServer().Unix()
			endtime := timeNow + SUPPORT_HERO_END_TIME*DAY_SECS

			var mastermsg S2M_SupportHeroUse
			mastermsg.Uid = useheroUid
			mastermsg.HeroKeyId = useheroKey
			mastermsg.Useruid = self.player.Sql_UserBase.Uid
			mastermsg.Username = self.player.Sql_UserBase.UName
			mastermsg.Type = HERO_SUPPORT_TYPE_ENTANGLEMENT
			mastermsg.Endtime = endtime

			ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_USE, &mastermsg)
			if ret == nil || ret.RetCode != UNION_SUCCESS {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
				return
			}

			info.FateHero = append(info.FateHero, &FateHero{
				index,
				useheroUid,
				usehero.MasterName,
				usehero.HeroKey,
				heroid,
				usehero.HeroStar,
				usehero.HeroLv,
				usehero.HeroSkin})
		}

		useheroUids = append(useheroUids, useheroUid)
		useheroKeys = append(useheroKeys, useheroKey)
		useIndex = append(useIndex, index)
		//heros = append(heros, heroid)

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ENTANGLEMENT_SET, heroid, msg.Type, int(useheroUid), "设置羁绊英雄", 0, len(info.Fate), self.player)

	}

	if !change {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	// 检查属性
	info.CheckProperty()
	// 添加属性
	self.AddTypeAllProperty(msg.Type)

	var backmsg S2C_EntanglementAutoUse
	backmsg.Cid = ENTANGLEMENT_AUTO_SET
	backmsg.MasterUid = useheroUids
	backmsg.Index = useIndex
	backmsg.Type = msg.Type
	backmsg.HeroKey = useheroKeys
	backmsg.Info = append(backmsg.Info, info)
	//for _, value := range heros {
	//	backmsg.Property = self.GetHeroPropertyData(value)
	//}
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}
