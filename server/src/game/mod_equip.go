package game

import (
	"encoding/json"
	"fmt"
	"sync"
)

// 缓存英雄属性
type HeroAttr struct {
	HeroId int          `json:"heroid"` //! 英雄Id
	Total  []*Attribute `json:"total"`  //! 总属性
}

const (
	RAND_MAX                = 10000
	BASE_GROUP              = 1
	SPECIAL_GROUP           = 2
	EQUIP_LVUP_GOLD_RATE    = 100  //消耗金币等于经验的10倍
	EQUIP_EXCLUSIVE_OPEN_LV = 8    //专属装备开放限制
	EQUIP_PACK_NUM          = 5    //分页个数
	EQUIP_PACK_LEN          = 1000 //分页容量
)

//信息交互中的装备位置索引
const (
	EQUIPPOS_START = iota
	EQUIPPOS_1
	EQUIPPOS_2
	EQUIPPOS_3
	EQUIPPOS_4
	EQUIPPOS_END
)

//装备初始  强化等级=0，附魔等级=0
//最终属性:
//base_valueN+upgrade_valueN*(强化等级-1)+star_valueN*附魔等级
const (
	wearEquip         = 1 // 穿戴装备
	takeOffEquip      = 2 // 脱装备
	upgradeEquip      = 3 // 升阶
	recastEquip       = 4 // 重铸
	recastEquipChoose = 5 // 重铸确认
	recastEquipCancel = 6 // 重铸取消

	wearEquipAll    = 1 // 一键穿戴装备
	takeOffEquipAll = 2 // 一键脱装备

	unLockExclusive = 1 // 解锁专属
	lvUpExclusive   = 2 // 升级专属

	compoundEquip          = 1   // 装备合成
	decompoundEquip        = 2   // 装备分解
	upgradeEquipAuto       = 6   // 装备一键强化
	starEquip              = 7   // 装备附魔
	rebornEquip            = 8   // 装备重生
	upgradeEquipAutoAll    = 9   // 装备全体一键强化
	upgradeEquipAutoAllCal = 10  // 装备全体一键强化(只计算，避免前后显示不一致)
	maxEquipNum            = 500 // 最大数量
	EquipChipType          = 21  // 装备碎片
	wearAllEquip           = 30  // 装备一键穿戴
	takeOffAllEquip        = 31  // 装备一键卸下

	compoundGem     = 1 // 合成宝石
	wearGem         = 2 // 宝石镶嵌
	takeOffGem      = 3 // 宝石卸下
	takeOffGemAuto  = 4 // 宝石一键卸下
	upLevelGem      = 5 // 宝石强化
	autoCompound    = 6 // 一键合成
	maxEquipWearNum = 6 // 最大穿戴个数
	autoWearGem     = 7 // 宝石一键安装

)

// 服务器存放的装备信息
type Equip struct {
	KeyId     int `json:"keyid"`     //! 装备唯一Id
	Id        int `json:"id"`        //! 装备配置Id
	HeroKeyId int `json:"herokeyid"` //! 装备拥有者
	//AttrInfo  []*AttrInfo `json:"attrinfo"`  //! 属性信息
	IsUpGrade int `json:"isupgrade"` //! 是否升阶(废弃)
	Lv        int `json:"lv"`        //!
	Exp       int `json:"exp"`       //!
	Recast    int `json:"recast"`    //!
}

// 可以升级的属性
type AttrInfo struct {
	AttrId    int   `json:"attrid"`
	AttrType  int   `json:"attrtype"`
	AttrValue int64 `json:"attrvalue"`
}

// 装备
type San_Equip struct {
	Uid           int64  // 玩家Id
	Maxkey        int    // 装备最大keyId
	Info          string // 装备
	TotalGemLevel int    // 总宝石等级
	StartTime     int64  //! 时间戳
	Info2         string // 装备
	Info3         string // 装备
	Info4         string // 装备
	Info5         string // 装备

	equipItems [EQUIP_PACK_NUM]map[int]*Equip // key: 装备唯一Id, value:装备信息

	DataUpdate
}

// 装备系统
type ModEquip struct {
	player     *Player   //! 玩家
	Data       San_Equip //! 数据库数据
	DataLocker *sync.RWMutex
}

// 装备唯一Id
func (self *ModEquip) MaxKey() int {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	self.Data.Maxkey += 1
	return self.Data.Maxkey
}

// save
func (self *ModEquip) Decode() {
	json.Unmarshal([]byte(self.Data.Info), &self.Data.equipItems[0])
	json.Unmarshal([]byte(self.Data.Info2), &self.Data.equipItems[1])
	json.Unmarshal([]byte(self.Data.Info3), &self.Data.equipItems[2])
	json.Unmarshal([]byte(self.Data.Info4), &self.Data.equipItems[3])
	json.Unmarshal([]byte(self.Data.Info5), &self.Data.equipItems[4])
}

// read
func (self *ModEquip) Encode() {
	self.Data.Info = HF_JtoA(self.Data.equipItems[0])
	self.Data.Info2 = HF_JtoA(self.Data.equipItems[1])
	self.Data.Info3 = HF_JtoA(self.Data.equipItems[2])
	self.Data.Info4 = HF_JtoA(self.Data.equipItems[3])
	self.Data.Info5 = HF_JtoA(self.Data.equipItems[4])
}

// get player and init map
func (self *ModEquip) OnGetData(player *Player) {
	self.player = player
	self.DataLocker = new(sync.RWMutex)
	self.checkEquip()
	tableName := self.getTableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", tableName, self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Data, tableName, self.player.ID)
	if self.Data.Uid <= 0 {
		self.init(self.player.ID)
		self.Encode()
		InsertTable(tableName, &self.Data, 0, true)
	} else {
		self.Decode()
		self.checkEquip()
	}

	self.Data.Init(tableName, &self.Data, true)

	if self.Data.TotalGemLevel <= 0 {
		//self.GetAllEquipGemLevel()
	}
}

// get data from db
func (self *ModEquip) OnGetOtherData() {
	self.checkEquip()
}

func (self *ModEquip) getTableName() string {
	return "san_userequip"
}

func (self *ModEquip) init(uid int64) {
	self.Data.Uid = uid
	self.checkEquip()
}

// make treasure map case nil
func (self *ModEquip) checkEquip() {
	for i := 0; i < len(self.Data.equipItems); i++ {
		if self.Data.equipItems[i] == nil {
			self.Data.equipItems[i] = make(map[int]*Equip, 0)
		}
	}
}

// save db every five minutes by changes
func (self *ModEquip) OnSave(sql bool) {
	self.Encode()
	self.Data.Update(sql)
}

// fresh every 5:00
func (self *ModEquip) OnRefresh() {

}

func (self *ModEquip) onReg(handlers map[string]func(body []byte)) {
	handlers["equipaction"] = self.onEquipAction
	handlers["equipuplv"] = self.EquipUpLv
	handlers["equipactionall"] = self.onEquipActionAll
	handlers["exclusiveaction"] = self.onExclusiveAction //专属装备
}

// 消息处理
func (self *ModEquip) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "gemaction":
		self.onGemAction(ctrl, body)
		return true
	}

	return false
}

func (self *ModEquip) onGemAction(cid string, body []byte) {
	var msg C2S_GemAction
	json.Unmarshal(body, &msg)
	if msg.Action == compoundGem {
		self.compoundGem(cid, &msg)
	} else if msg.Action == wearGem {
		//self.wearGem(cid, &msg)
	} else if msg.Action == takeOffGem {
		//self.takeOffGem(cid, &msg)
	} else if msg.Action == takeOffGemAuto {
		//self.takeOffGemAuto(cid, &msg)
	} else if msg.Action == upLevelGem {
		//self.upLevelGem(cid, &msg)
	} else if msg.Action == autoCompound {
		self.autoCompound(cid, &msg)
	} else if msg.Action == autoWearGem {
		//self.autoWearGem(cid, &msg)
	}
}

// 一键合成
// 60个1级宝石，点一键合成得到20个2级宝石
// 当前2级宝石999个,最多合成0个
// 当前2级宝石n个, 最多合成 config.stacknum = n + m
func (self *ModEquip) autoCompound(cid string, msg *C2S_GemAction) {
	gemId := msg.GemId
	// 检查配置是否存在
	config := GetCsvMgr().GetGemConfig(gemId)
	if config == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_GEMSTONE_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	// 检查道具Id
	itemConfig := GetCsvMgr().GetItemConfig(gemId)
	if itemConfig == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_GEM_PROPS_CONFIGURATION_DOES_NOT"))
		return
	}

	currentNum := self.player.GetObjectNum(gemId)
	compoundNum := itemConfig.MaxNum - currentNum
	if compoundNum <= 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_BEYOND_THE_UPPER_LIMIT_OF"))
		return
	}

	needId := config.NeedId
	needConfig := GetCsvMgr().GetItemConfig(needId)
	if needConfig == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_NEED_CONFIG_GEM_PROJECTS_CONFIGURATION"))
		return
	}

	if config.NeedNum == 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_CONFIGURATION_ERROR:_SYNTHETIC_GEMSTONE_REQUIRED"))
		return
	}
	ownNum := self.player.GetObjectNum(needId)

	// 当前可以合成的个数
	canCompoundNUm := ownNum / config.NeedNum

	// 实际需要可以合成的个数
	realNum := HF_MinInt(canCompoundNUm, compoundNum) // 30, 20 = 20;  10,30 = 10;
	if realNum <= 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_GEM_SHORTAGE"))
		return
	}

	// 扣除宝石以及钻石
	var items []PassItem
	// 扣除合成的钻石
	res := self.player.RemoveObjectEasy(needId, realNum*config.NeedNum, "宝石合成", 0, 0, 0)
	items = append(items, res...)
	// 增加钻石
	addRes := self.player.AddObjectSimple(gemId, realNum, "宝石合成", 0, 0, 0)
	items = append(items, addRes...)

	// 合成成功
	self.player.SendMsg(cid, HF_JtoB(&S2C_GemAction{
		Cid:   cid,
		Items: items,
	}))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GEM_COMPOUND, config.Id, config.Level, 0, "宝石合成", 0, 0, self.player)
}

// 装备相关操作
func (self *ModEquip) onEquipAction(body []byte) {
	var msg C2S_EquipAction
	json.Unmarshal(body, &msg)
	if msg.Action == wearEquip {
		self.wearEquip(&msg)
	} else if msg.Action == takeOffEquip {
		self.takeOffEquip(&msg)
	} else if msg.Action == upgradeEquip {
		self.upgradeEquip(&msg)
	} else if msg.Action == recastEquip {
		self.recastEquip(&msg)
	} else if msg.Action == recastEquipChoose {
		self.recastEquipChoose(&msg)
	} else if msg.Action == recastEquipCancel {
		self.recastEquipCancel(&msg)
	}
}

//专属装备相关操作
func (self *ModEquip) onExclusiveAction(body []byte) {

	var msg C2S_ExclusiveAction
	json.Unmarshal(body, &msg)
	if msg.Action == unLockExclusive {
		self.UnLockExclusive(&msg)
	} else if msg.Action == lvUpExclusive {
		self.LvUpExclusive(&msg)
	}
}

//装备一键操作
func (self *ModEquip) onEquipActionAll(body []byte) {
	var msg C2S_EquipActionAll
	json.Unmarshal(body, &msg)
	if msg.Action == wearEquipAll {
		self.wearAllEquip(&msg)
	} else if msg.Action == takeOffEquipAll {
		self.takeOffAllEquip(&msg)
	}
}

func (self *ModEquip) EquipUpLv(body []byte) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	var msg C2S_EquipUpLv
	json.Unmarshal(body, &msg)

	if len(msg.ItemKeyId) == 0 || len(msg.ItemKeyId) != len(msg.ItemId) || len(msg.ItemKeyId) != len(msg.ItemNum) {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_UPLV_MSG_ERROR"))
		return
	}

	pEquip := self.GetEquipItem(msg.EquipKeyId)
	if pEquip == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_DOES_NOT_EXIST"))
		return
	}

	var msgRel S2C_EquipUpLv
	msgRel.Cid = "equipuplv"

	oldlevel := pEquip.Lv
	addExp := 0
	for i := 0; i < len(msg.ItemKeyId); i++ {
		//如果是装备
		if msg.ItemKeyId[i] > 0 {
			equipItem := self.GetEquipItem(msg.ItemKeyId[i])
			if equipItem == nil || equipItem.HeroKeyId != 0 {
				continue
			}
			equipConfig := GetCsvMgr().GetEquipStrengthenLvUpConfig(equipItem.Id, equipItem.Lv)
			if equipConfig == nil {
				continue
			}

			if equipConfig.ExpBy <= 0 {
				continue
			}

			costItemAll := make([]int, 0)
			costNumAll := make([]int, 0)
			costItemAll = append(costItemAll, ITEM_GOLD)
			costNumAll = append(costNumAll, (equipConfig.ExpBy+equipItem.Exp)*EQUIP_LVUP_GOLD_RATE)
			//检查消耗够不够
			if err := self.player.HasObjectOk(costItemAll, costNumAll); err != nil {
				continue
			}

			addExp += equipConfig.ExpBy
			addExp += equipItem.Exp
			self.deleteEquip(msg.ItemKeyId[i])
			msgRel.ItemKeyId = append(msgRel.ItemKeyId, msg.ItemKeyId[i])

			items := self.player.RemoveObjectLst(costItemAll, costNumAll, "装备强化", pEquip.Id, pEquip.Lv, 0)
			msgRel.CostItem = append(msgRel.CostItem, items...)
		} else {
			//道具
			if msg.ItemNum[i] <= 0 {
				continue
			}
			iConfig := GetCsvMgr().ItemMap[msg.ItemId[i]]
			if iConfig == nil {
				continue
			}
			costItemAll := make([]int, 0)
			costNumAll := make([]int, 0)
			costItemAll = append(costItemAll, msg.ItemId[i])
			costNumAll = append(costNumAll, msg.ItemNum[i])
			costItemAll = append(costItemAll, ITEM_GOLD)
			costNumAll = append(costNumAll, iConfig.Special*msg.ItemNum[i]*EQUIP_LVUP_GOLD_RATE)

			//检查消耗够不够
			if err := self.player.HasObjectOk(costItemAll, costNumAll); err != nil {
				continue
			}
			addExp += iConfig.Special * msg.ItemNum[i]
			items := self.player.RemoveObjectLst(costItemAll, costNumAll, "装备强化", pEquip.Id, pEquip.Lv, 0)
			msgRel.CostItem = append(msgRel.CostItem, items...)
		}
	}
	pEquip.AddExp(addExp)
	hero := self.player.getHero(pEquip.HeroKeyId)
	if hero != nil {
		self.player.countHeroFight(hero, ReasonEquipLvUp)
		if GetOfflineInfoMgr().IsBaseHero(self.player.Sql_UserBase.Uid, hero.HeroKeyId) {
			self.player.NoticeCenterBaseInfo()
		}
	}
	//看看是否有材料返还
	configMax := GetCsvMgr().GetEquipStrengthenLvUpConfig(pEquip.Id, pEquip.Lv+1)
	if configMax == nil && pEquip.Exp > 0 {
		//说明达到了最大值并且发生了材料返还
		itemId := ITEM_EQUIP_LVUP_ITEM_LOW
		iConfig := GetCsvMgr().ItemMap[itemId]
		if iConfig != nil && iConfig.Special > 0 {
			itemNum := pEquip.Exp / iConfig.Special
			pEquip.Exp = 0
			msgRel.ReturnItem = self.player.AddObjectSimple(itemId, itemNum, "装备强化MAX返还", 0, 0, 0)
		}
	}
	self.player.HandleTask(TASK_TYPE_EQUIP_LEVEL_UP, 0, 0, 0)
	self.player.HandleTask(TASK_TYPE_HAVE_EQUIP, 0, 0, 0)
	msgRel.Equip = pEquip
	msgRel.CostItem = HF_MergePassitem(msgRel.CostItem)
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_LEVEL_UP, pEquip.Id, oldlevel, pEquip.Lv, "装备强化", 0, msg.EquipKeyId, self.player)
}

//装备增加经验
func (self *Equip) AddExp(exp int) {

	self.Exp += exp
	for {
		if exp <= 0 {
			break
		}
		//看是否达到最大等级
		config := GetCsvMgr().GetEquipStrengthenLvUpConfig(self.Id, self.Lv+1)
		if config == nil {
			return
		}

		if self.Exp >= config.UpLvNeedExp {
			self.Exp -= config.UpLvNeedExp
			self.Lv++
		} else {
			break
		}
	}
}

// 装备重生
func (self *ModEquip) rebornEquip(cid string, msg *C2S_EquipAction) {
	/*
		keyIds := msg.RemoveKeyIds
		if len(keyIds) <= 0 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_CLIENT_RENAISSANCE_LEN_KEYIDS_="))
			return
		}

		if len(keyIds) > 1 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_CLIENT_RENAISSANCE_LEN_KEYIDS_>"))
			return
		}

		// check
		itemMap := make(map[int]*Item)
		for _, keyId := range keyIds {
			pEquip := self.GetEquipItem(keyId)
			if pEquip == nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_DOES_NOT_EXIST"))
				return
			}

			heroId := self.GetTeamHero(pEquip.KeyId)
			if heroId != 0 {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_WEARED_EQUIPMENT_CANNOT_BE_DECOMPOSED"))
				return
			}

			config := GetCsvMgr().GetEquipConfig(pEquip.Id)
			if config == nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_CONFIGURATION_DOES_NOT_EXIST"))
				return
			}

			// 强化材料
			upItemMap := pEquip.getUpgradeReborn()
			AddItemMap(itemMap, upItemMap)

			// 附魔材料
			starReborn := pEquip.getStarReborn()
			AddItemMap(itemMap, starReborn)

			// 宝石道具
			gemReborn := pEquip.getGemReborn()
			AddItemMap(itemMap, gemReborn)

			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_REBORN, pEquip.Id, 0, 0, "重生", 0, 0, self.player)
		}

		if len(itemMap) <= 0 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_WHITE_CLOTHES_CANNOT_BE_REBORN"))
			return
		}

		outItems := self.player.AddObjectItemMap(itemMap, "重生", 0, 0, 0)
		// 装备清零
		euqips := self.rebornEquips(keyIds)
		self.player.SendMsg(cid, HF_JtoB(&S2C_EquipAction{
			Cid:          cid,
			Action:       msg.Action,
			Items:        outItems,
			RemoveKeyIds: msg.RemoveKeyIds,
			EquipInfos:   euqips,
			HeroId:       msg.HeroId,
		}))

	*/
}

func (self *ModEquip) deleteEquip(keyId int) {
	for i := 0; i < len(self.Data.equipItems); i++ {
		delete(self.Data.equipItems[i], keyId)
	}
}

//! 合成装备
func (self *ModEquip) compoundEquip(cid string, msg *C2S_EquipAction) {
	/*
		itemId := msg.Itemid
		compoundNum := msg.CompundNum
		if compoundNum <= 0 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TREASURE_ILLEGAL_NUMBER_OF_PROPS_SENT"))
			return
		}

		// 最多合成一个
		if compoundNum > 1 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_NUMBER_EXCEEDING_SYNTHETIC_UPPER_LIMIT"))
			return
		}

		itemConfig, ok := GetCsvMgr().ItemMap[itemId]
		if itemConfig.ItemType != EquipChipType {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_THE_TYPE_OF_EQUIPMENT_DEBRIS"))
			return
		}

		if !ok {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_SHOP_PROJECTS_CONFIGURATION_DOES_NOT_EXIST"))
			return
		}

		// 合成后的道具
		if itemConfig.CompoundId == 0 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TREASURE_THE_PROP_CANNOT_BE_SYNTHESIZED"))
			return
		}

		needNum := itemConfig.CompoundNum * compoundNum
		if self.player.GetObjectNum(itemId) < needNum {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TREASURE_LACK_OF_PROPS"))
			return
		}

		if len(self.Data.EquipItems) >= maxEquipNum {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_UP_TO_THE_MAXIMUM_NUMBER"))
			return
		}

		var res []*Equip
		for i := 1; i <= compoundNum; i++ {
			pItem := self.NewEquipItem(itemConfig.CompoundId)
			self.Data.EquipItems[pItem.KeyId] = pItem
			res = append(res, pItem)
		}

		// 扣除道具
		items := self.player.RemoveObjectSimple(itemId, needNum, "装备合成", 0, 0, 0)
		self.player.HandleTask(EquipColorTask, 0, 0, 0)
		self.player.SendMsg(cid, HF_JtoB(&S2C_EquipAction{
			Cid:         cid,
			Action:      msg.Action,
			ComposeInfo: res,
			Items:       items,
			HeroId:      msg.HeroId,
		}))

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_COMPOUND, itemConfig.CompoundId, 0, 0, "装备合成", 0, 0, self.player)

	*/
}

// 创建一个装备
func (self *ModEquip) NewEquipItem(equipId int) *Equip {
	pItem := &Equip{}
	pItem.KeyId = self.MaxKey()
	pItem.Id = equipId

	config, ok := GetCsvMgr().EquipConfigMap[equipId]
	if !ok {
		LogError("equip is not found, equipId=", equipId)
		return nil
	}
	//生成属性
	for i := 0; i < len(config.BaseTypes); i++ {
		if config.BaseTypes[i] > 0 {
			attr := new(AttrInfo)
			attr.AttrId = i + 1
			//pItem.AttrInfo = append(pItem.AttrInfo, attr)
		}
	}

	return pItem
}

func (self *ModEquip) NewSuperEquipItem(equipId int) *Equip {
	pItem := &Equip{}
	pItem.KeyId = self.MaxKey()
	pItem.Id = equipId
	pItem.Lv = 5

	config, ok := GetCsvMgr().EquipConfigMap[equipId]
	if !ok {
		LogError("equip is not found, equipId=", equipId)
		return nil
	}
	//生成属性
	for i := 0; i < len(config.BaseTypes); i++ {
		if config.BaseTypes[i] > 0 {
			attr := new(AttrInfo)
			attr.AttrId = i + 1
		}
	}

	for j := 0; j < len(self.Data.equipItems); j++ {
		if len(self.Data.equipItems[j]) < EQUIP_PACK_LEN {
			self.Data.equipItems[j][pItem.KeyId] = pItem
			break
		}
	}
	return pItem
}

// 装备升阶
func (self *ModEquip) upgradeEquip(msg *C2S_EquipAction) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	pEquip := self.GetEquipItem(msg.EquipKeyId)
	if pEquip == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_DOES_NOT_EXIST"))
		return
	}

	config := GetCsvMgr().EquipConfigMap[pEquip.Id]
	if config == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_CONFIG_NOT_EXIST"))
		return
	}

	if GetCsvMgr().EquipConfigMap[config.AdvanceId] == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_CONFIG_NOT_EXIST"))
		return
	}

	configAdv := GetCsvMgr().EquipAdvancedConfigMap[config.Quality]
	if configAdv == nil || len(configAdv.AdvancedNeed) <= config.EquipAttackType-1 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_CONFIG_NOT_EXIST"))
		return
	}

	// 检查消耗是否正常
	if err := self.player.HasObjectOkEasy(configAdv.AdvancedNeed[config.EquipAttackType-1], configAdv.AdvancedNum[config.EquipAttackType-1]); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	oldid := pEquip.Id
	// 进行强化操作
	items := self.player.RemoveObjectSimple(configAdv.AdvancedNeed[config.EquipAttackType-1], configAdv.AdvancedNum[config.EquipAttackType-1], "装备升阶", pEquip.Id, 0, 0)
	pEquip.Id = config.AdvanceId
	if pEquip.Recast > 0 {
		pEquip.Recast = (pEquip.Id/10)*10 + pEquip.Recast%10
	}

	equipHeroId := self.GetTeamHero(pEquip.KeyId)
	hero := self.player.getHero(equipHeroId)
	if hero != nil {
		self.player.countHeroFight(hero, ReasonEquipWear)
		if GetOfflineInfoMgr().IsBaseHero(self.player.Sql_UserBase.Uid, hero.HeroKeyId) {
			self.player.NoticeCenterBaseInfo()
		}
	}
	self.player.SendMsg("equipaction", HF_JtoB(&S2C_EquipAction{
		Cid:       "equipaction",
		Action:    msg.Action,
		EquipItem: pEquip,
		Items:     items,
	}))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_UPGRADE, pEquip.Id, oldid, msg.EquipKeyId, "装备升阶", 0, 0, self.player)
}

func (self *ModEquip) recastEquip(msg *C2S_EquipAction) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	pEquip := self.GetEquipItem(msg.EquipKeyId)
	if pEquip == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_DOES_NOT_EXIST"))
		return
	}

	config := GetCsvMgr().EquipConfigMap[pEquip.Id]
	if config == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_CONFIG_NOT_EXIST"))
		return
	}

	configRecast := GetCsvMgr().GetEquipRecastConfig(config.Quality, config.Camp)
	if configRecast == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_CONFIG_NOT_EXIST"))
		return
	}

	// 检查消耗是否正常
	if err := self.player.HasObjectOkEasy(configRecast.CostTime, configRecast.CostNum); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	// 进行强化操作
	items := self.player.RemoveObjectSimple(configRecast.CostTime, configRecast.CostNum, "装备重铸", pEquip.Id, 0, 0)
	allRate := 0
	for i := 0; i < len(configRecast.Weight); i++ {
		allRate += configRecast.Weight[i]
	}
	rand := HF_GetRandom(allRate)
	rateNow := 0
	for i := 0; i < len(configRecast.Weight); i++ {
		rateNow += configRecast.Weight[i]
		if rateNow > rand {
			pEquip.Recast = (pEquip.Id/10)*10 + configRecast.Change[i]
			break
		}
	}

	self.player.SendMsg("equipaction", HF_JtoB(&S2C_EquipAction{
		Cid:       "equipaction",
		Action:    msg.Action,
		EquipItem: pEquip,
		Items:     items,
	}))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_RECAST, pEquip.Id, msg.EquipKeyId, 0, "装备重铸", 0, 0, self.player)
}

func (self *ModEquip) recastEquipChoose(msg *C2S_EquipAction) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	pEquip := self.GetEquipItem(msg.EquipKeyId)
	if pEquip == nil || pEquip.Recast == 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_DOES_NOT_EXIST"))
		return
	}

	oldid := pEquip.Id
	pEquip.Id = pEquip.Recast
	pEquip.Recast = 0

	equipHeroId := self.GetTeamHero(pEquip.KeyId)
	hero := self.player.getHero(equipHeroId)
	if hero != nil {
		self.player.countHeroFight(hero, ReasonEquipWear)
		if GetOfflineInfoMgr().IsBaseHero(self.player.Sql_UserBase.Uid, hero.HeroKeyId) {
			self.player.NoticeCenterBaseInfo()
		}
	}
	self.player.SendMsg("equipaction", HF_JtoB(&S2C_EquipAction{
		Cid:       "equipaction",
		Action:    msg.Action,
		EquipItem: pEquip,
	}))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_RECAST_CHOOSE, pEquip.Id, oldid, 0, "装备重铸保存", 0, 0, self.player)
}

func (self *ModEquip) recastEquipCancel(msg *C2S_EquipAction) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	pEquip := self.GetEquipItem(msg.EquipKeyId)
	if pEquip == nil || pEquip.Recast == 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_DOES_NOT_EXIST"))
		return
	}

	pEquip.Recast = 0

	equipHeroId := self.GetTeamHero(pEquip.KeyId)
	hero := self.player.getHero(equipHeroId)
	if hero != nil {
		self.player.countHeroFight(hero, ReasonEquipWear)
	}
	self.player.SendMsg("equipaction", HF_JtoB(&S2C_EquipAction{
		Cid:       "equipaction",
		Action:    msg.Action,
		EquipItem: pEquip,
	}))
}

func (self *ModEquip) GetEquipItem(keyId int) *Equip {
	for i := 0; i < len(self.Data.equipItems); i++ {
		pItem, ok := self.Data.equipItems[i][keyId]
		if ok {
			return pItem
		}
	}
	return nil
}

/*
func (self *ModEquip) CheckTeam(msg *C2S_EquipAction) (*TeamAttr, int, error) {
	index := msg.Index
	if index < 1 || index > 5 {
		return nil, 0, errors.New(GetCsvMgr().GetText("STR_MOD_EQUIP_SUBSCRIPT_ERROR"))
	}

	teamType := msg.TeamType
	if teamType < 1 || teamType > MAX_TEAM_TYPE {
		return nil, 0, errors.New(GetCsvMgr().GetText("STR_MOD_TIGER_NO_BATTLE_TYPE_EXISTS"))
	}

	teamPos := self.player.getTeamPosByType(teamType)
	if teamPos == nil {
		return nil, 0, errors.New("teamPos == nil")
	}

	heroId := teamPos.getHeroId(index - 1)
	if heroId == 0 {
		return nil, 0, errors.New(GetCsvMgr().GetText("STR_MOD_TIGER_HEROES_WHO_ARE_NOT_IN"))
	}

	teamAttr := teamPos.TeamAttr
	if len(teamAttr) != 5 {
		return nil, 0, errors.New("len(teamAttr) != 5")
	}

	pTeamAttr := teamAttr[index-1]
	return pTeamAttr, heroId, nil
}
*/

/*
func (self *ModEquip) GetTeamAttr(keyId int, teamType int, pos int) (*TeamAttr, int, int) {
	teamPos := self.player.getTeamPosByType(teamType)
	if teamPos == nil {
		return nil, -1, 0
	}

	teamAttr := teamPos.TeamAttr
	if len(teamAttr) != 5 {
		return nil, -1, 0
	}

	for i, v := range teamAttr {
		if v.EquipIds[pos-1] == keyId {
			return teamAttr[i], i, teamPos.UIPos[i]
		}
	}

	return nil, -1, 0
}
*/

// 穿装备, teamType, teamIndex, keyId
func (self *ModEquip) wearEquip(msg *C2S_EquipAction) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	// 需要穿戴的装备
	keyId := msg.EquipKeyId
	pItem := self.GetEquipItem(keyId)
	if pItem == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_DOES_NOT_EXIST"))
		return
	}

	itemId := pItem.Id
	pConfig, ok := GetCsvMgr().EquipConfigMap[itemId]
	if !ok {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	// 通过需要穿戴的装备的配置取位置
	pos := pConfig.EquipPosition
	if pos <= EQUIPPOS_START || pos >= EQUIPPOS_END {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_ERROR_IN_EQUIPMENT_TYPE"))
		return
	}

	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroKeyId))
		return
	}

	//验证职业
	heroConfig := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if heroConfig == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_CAREER_ERROR"))
		return
	}

	if heroConfig.AttackType != pConfig.EquipAttackType {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_CAREER_ERROR"))
		return
	}

	if len(hero.EquipIds) != EQUIPPOS_END-EQUIPPOS_1 {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST"))
		return
	}

	if pItem.HeroKeyId != 0 {
		oldHero := self.player.getHero(pItem.HeroKeyId)
		if oldHero == nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroKeyId))
			return
		}
		if pItem.HeroKeyId == msg.HeroKeyId {
			self.player.SendErr(GetCsvMgr().GetText("STR_HERO_ALREADY_WEAR"))
			return
		}
	}

	// 获得当前位置的装备
	oldKeyId := hero.EquipIds[pos-1]
	// 如果是换装备 当前位置上有装备
	if oldKeyId != 0 {
		hero.EquipIds[pos-1] = 0
	}

	// 是否换下的装备装备到另一个英雄上
	bExchange := false
	oldHeroId := 0
	if pItem.HeroKeyId != 0 {
		oldHeroId = pItem.HeroKeyId
		oldHero := self.player.getHero(pItem.HeroKeyId)
		if oldKeyId != 0 {
			oldHero.EquipIds[pos-1] = oldKeyId
			bExchange = true
		} else {
			oldHero.EquipIds[pos-1] = 0
		}
	}

	// 如果是简单的换装备
	if hero.EquipIds[pos-1] == 0 {
		hero.EquipIds[pos-1] = keyId
	}

	oldEquip := self.GetEquipItem(oldKeyId)
	if oldEquip != nil {
		// 换下来的装备需要换位置装备时
		if bExchange {
			oldEquip.HeroKeyId = oldHeroId
			oldHero := self.player.getHero(pItem.HeroKeyId)
			self.player.countHeroFight(oldHero, ReasonEquipWear)
			if GetOfflineInfoMgr().IsBaseHero(self.player.Sql_UserBase.Uid, oldEquip.HeroKeyId) {
				self.player.NoticeCenterBaseInfo()
			}
		} else {
			oldEquip.HeroKeyId = 0
		}
	}

	self.player.countHeroFight(hero, ReasonEquipWear)
	pItem.HeroKeyId = msg.HeroKeyId
	if GetOfflineInfoMgr().IsBaseHero(self.player.Sql_UserBase.Uid, pItem.HeroKeyId) {
		self.player.NoticeCenterBaseInfo()
	}

	self.player.SendMsg("equipaction", HF_JtoB(&S2C_EquipAction{
		Cid:          "equipaction",
		Action:       msg.Action,
		EquipItem:    pItem,
		ExchangeItem: self.GetEquipItem(oldKeyId),
	}))

	if oldKeyId == 0 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_WEAR, itemId, pos, 0, "穿戴装备", 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_CHANGE, itemId, pos, oldEquip.Id, "装备替换", 0, 0, self.player)
	}

}

func (self *ModEquip) wearAllEquip(msg *C2S_EquipActionAll) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroKeyId))
		return
	}

	heroConfig := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if heroConfig == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_CAREER_ERROR"))
		return
	}

	equip := make([]*Equip, 0)
	var equipBest = [EQUIPPOS_END - EQUIPPOS_1]*Equip{nil, nil, nil, nil}
	for i := 0; i < len(hero.EquipIds); i++ {
		if hero.EquipIds[i] == 0 {
			continue
		}
		pItem := self.GetEquipItem(hero.EquipIds[i])
		if pItem == nil {
			continue
		}
		pConfig, ok := GetCsvMgr().EquipConfigMap[pItem.Id]
		if !ok {
			continue
		}

		pos := pConfig.EquipPosition
		if pos <= EQUIPPOS_START || pos >= EQUIPPOS_END {
			continue
		}
		equipBest[pos-1] = pItem
	}

	for i := 0; i < len(hero.EquipIds); i++ {
		//先找到每个部位战斗力最高并且无人穿戴的装备
		for j := 0; j < len(self.Data.equipItems); j++ {
			for _, v := range self.Data.equipItems[j] {
				if v.HeroKeyId != 0 && v.HeroKeyId != hero.HeroKeyId {
					continue
				}

				pItem := self.GetEquipItem(v.KeyId)
				if pItem == nil {
					continue
				}

				pConfig, ok := GetCsvMgr().EquipConfigMap[pItem.Id]
				if !ok {
					continue
				}

				pos := pConfig.EquipPosition
				if pos <= EQUIPPOS_START || pos >= EQUIPPOS_END {
					continue
				}

				if heroConfig.AttackType != pConfig.EquipAttackType {
					continue
				}

				//和之前的选择进行比较
				if equipBest[pos-1] == nil {
					equipBest[pos-1] = pItem
				} else {
					//先比较品质
					preConfig := GetCsvMgr().EquipConfigMap[equipBest[pos-1].Id]
					if preConfig == nil {
						continue
					}
					if pConfig.Quality > preConfig.Quality {
						equipBest[pos-1] = pItem
						continue
					} else if pConfig.Quality < preConfig.Quality {
						continue
					}
					//相等的情况比较对应种族
					if pConfig.Camp == heroConfig.Attribute && preConfig.Camp != heroConfig.Attribute {
						equipBest[pos-1] = pItem
						continue
					} else if pConfig.Camp != heroConfig.Attribute && preConfig.Camp == heroConfig.Attribute {
						continue
					}
					if pItem.Lv > equipBest[pos-1].Lv {
						equipBest[pos-1] = pItem
					} else if pItem.Lv == equipBest[pos-1].Lv {
						if pConfig.Camp != 0 && preConfig.Camp == 0 {
							equipBest[pos-1] = pItem
							continue
						}
					}
				}
			}
		}
	}

	//开始更换装备
	isHasBest := false
	for i := 0; i < len(equipBest); i++ {
		if equipBest[i] == nil {
			continue
		}
		isHasBest = true
		//看看当前玩家身上的装备,同一件也不更新
		nowKeyId := hero.EquipIds[i]
		if nowKeyId == equipBest[i].KeyId {
			continue
		}
		//开始更换
		// 先拆
		if nowKeyId != 0 {
			nowItem := self.GetEquipItem(nowKeyId)
			if nowItem == nil {
				continue
			}
			nowItem.HeroKeyId = 0
			equip = append(equip, nowItem)
		}
		//然后装
		equipBest[i].HeroKeyId = hero.HeroKeyId
		hero.EquipIds[i] = equipBest[i].KeyId
		equip = append(equip, equipBest[i])
	}

	if !isHasBest {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_WEAR_EQUIP_ALL_NOT_EQUIP"))
		return
	}

	if len(equip) == 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_WEAR_EQUIP_ALL_BEST"))
		return
	}

	self.player.countHeroFight(hero, ReasonEquipWear)
	if GetOfflineInfoMgr().IsBaseHero(self.player.Sql_UserBase.Uid, hero.HeroKeyId) {
		self.player.NoticeCenterBaseInfo()
	}
	self.player.SendMsg("equipactionall", HF_JtoB(&S2C_EquipActionAll{
		Cid:       "equipactionall",
		Action:    msg.Action,
		EquipItem: equip,
	}))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_WEAR_AUTO, hero.HeroId, hero.HeroKeyId, 0, "一键穿戴装备", 0, 0, self.player)
}

func (self *ModEquip) takeOffAllEquip(msg *C2S_EquipActionAll) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroKeyId))
		return
	}

	heroConfig := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if heroConfig == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST"))
		return
	}

	equip := make([]*Equip, 0)
	for i := 0; i < len(hero.EquipIds); i++ {
		pItem := self.GetEquipItem(hero.EquipIds[i])
		if pItem == nil {
			continue
		}
		pItem.HeroKeyId = 0
		hero.EquipIds[i] = 0
		equip = append(equip, pItem)
	}

	if len(equip) > 0 {
		self.player.countHeroFight(hero, ReasonEquipOff)
		if GetOfflineInfoMgr().IsBaseHero(self.player.Sql_UserBase.Uid, hero.HeroKeyId) {
			self.player.NoticeCenterBaseInfo()
		}
	}

	self.player.SendMsg("equipactionall", HF_JtoB(&S2C_EquipActionAll{
		Cid:       "equipactionall",
		Action:    msg.Action,
		EquipItem: equip,
	}))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_OFF_AUTO, hero.HeroId, hero.HeroKeyId, 0, "一键卸下装备", 0, 0, self.player)
}

// 脱除装备, 顺便把宝石也脱下来
func (self *ModEquip) takeOffEquip(msg *C2S_EquipAction) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroKeyId))
		return
	}

	if len(hero.EquipIds) != EQUIPPOS_END-EQUIPPOS_1 {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST"))
		return
	}

	pItem := self.GetEquipItem(msg.EquipKeyId)
	if pItem == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_NO_EQUIPMENT_DETACHABLE"))
		return
	}

	itemConfig := GetCsvMgr().EquipConfigMap[pItem.Id]

	hero.EquipIds[itemConfig.EquipPosition-1] = 0

	self.player.countHeroFight(hero, ReasonEquipOff)
	pItem.HeroKeyId = 0
	if GetOfflineInfoMgr().IsBaseHero(self.player.Sql_UserBase.Uid, hero.HeroKeyId) {
		self.player.NoticeCenterBaseInfo()
	}

	self.player.SendMsg("equipaction", HF_JtoB(&S2C_EquipAction{
		Cid:       "equipaction",
		Action:    msg.Action,
		EquipItem: pItem,
	}))
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_OFF, pItem.Id, itemConfig.EquipPosition, 0, "装备卸下", 0, 0, self.player)
}

// 装备一键强化, 只能强化一键装备,有多少强化多少
func (self *ModEquip) upgradeEquipAuto(cid string, msg *C2S_EquipAction) {
	/*
		pEquip := self.GetEquipItem(msg.KeyId)
		if pEquip == nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_DOES_NOT_EXIST"))
			return
		}

		// 检查是否强化到满级
		pItem := pEquip.UpGrade
		if pItem == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_EQUIP_ENHANCED_MODULE_DATA_EXCEPTION"))
			return
		}

		// 检查是否达到最大等级
		currentLv := pItem.Lv
		maxLv, ok := GetCsvMgr().EquipUpgradeMaxLv[pItem.Id]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_EQUIP_CURRENT_ENHANCED_CONFIGURATION_DOES_NOT"))
			return
		}

		// 可以强化到玩家*2
		reachLv := self.player.Sql_UserBase.Level * 2
		maxLv = HF_MinInt(reachLv, maxLv)

		if currentLv >= maxLv {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_EQUIP_HAS_REACHED_THE_MAXIMUM_LEVEL"))
			return
		}

		// 一级一级判断直到钱不够为止
		config := GetCsvMgr().GetEquipUpGrade(pItem.Id, currentLv)
		if config == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_EQUIP_CURRENT_ENHANCED_CONFIGURATION_DOES_NOT"))
			return
		}

		itemMap := make(map[int]*Item)
		targetLv := currentLv
		// 一级一级检查配置
		for level := currentLv; level < maxLv; level++ {
			config := GetCsvMgr().GetEquipUpGrade(pItem.Id, level)
			if config == nil {
				continue
			}

			costIds := config.CostIds
			costNums := config.CostNums
			if len(costIds) != len(costNums) {
				LogError("len(costIds) != len(costNums)")
				continue
			}

			for costIndex := range costIds {
				itemId := costIds[costIndex]
				itemNum := costNums[costIndex]
				if itemId == 0 || itemNum == 0 {
					continue
				}

				pItem, ok := itemMap[itemId]
				if !ok {
					itemMap[itemId] = &Item{itemId, itemNum}
				} else {
					pItem.ItemNum += itemNum
				}
			}

			err := self.player.hasItemMapOk(itemMap)
			if err != nil {
				break
			} else {
				targetLv = level + 1
			}
		}

		// 不能再升级了表示材料不够
		if targetLv == currentLv {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_INSUFFICIENT_MATERIAL"))
			return
		}

		if targetLv < currentLv {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_EQUIP_LOGIC_EXCEPTION"))
			return
		}

		// 计算实际消耗
		totalCost := make(map[int]*Item)
		for level := currentLv; level < targetLv; level++ {
			config := GetCsvMgr().GetEquipUpGrade(pItem.Id, level)
			if config == nil {
				continue
			}
			AddItemMapHelper(totalCost, config.CostIds, config.CostNums)
		}

		// 进行强化操作
		items := self.player.RemoveObjectItemMap(totalCost, "装备强化", targetLv, pEquip.Id, 0)
		pItem.Lv = targetLv
		self.player.HandleTask(EquipUpNumTask, 0, 0, 0)
		self.player.HandleTask(EquipUpTimesTask, targetLv-currentLv, 0, 0)

		equipHeroId := self.GetTeamHero(pEquip.KeyId)
		self.player.countHeroIdFight(equipHeroId, ReasonUpgradeEquipAuto)
		self.player.SendMsg(cid, HF_JtoB(&S2C_EquipAction{
			Cid:     cid,
			Action:  msg.Action,
			PItem:   pEquip,
			Items:   items,
			HeroAtt: self.getHeroAtt(equipHeroId),
			HeroId:  msg.HeroId,
		}))

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_UPGRADE, pItem.Lv, pEquip.Id, 0, "装备强化", 0, 0, self.player)

	*/
}

// 装备一键强化全体位置  20190723 by zy
// 20190729  by zy 改版  等级阶梯追赶，并且平均施行
func (self *ModEquip) upgradeEquipAutoAll(cid string, msg *C2S_EquipAction, isUp bool) {
	/*
		hero := self.player.getHero(msg.HeroId)

		if hero == nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroId))
			return
		}

		if len(hero.EquipIds) != EQUIPPOS_END-EQUIPPOS_1 {
			return
		}

		//变量初始化，可以提高效率，空间换时间
		//强化完之后用来判断是否需要发送变化给客户端
		oldLevel := make([]int, 6)
		//阶梯等级  需要补充1个最高等级
		levelStage := make([]int, 7)
		//6个装备
		equips := make([]*Equip, 6)
		for stageIndex := 0; stageIndex <= 5; stageIndex++ {
			keyId := hero.EquipIds[stageIndex]
			pEquip := self.GetEquipItem(keyId)

			if pEquip == nil {
				continue
			}

			itemId := pEquip.Id
			pConfig, ok := GetCsvMgr().EquipMap[itemId]
			if !ok {
				continue
			}

			// 通过需要穿戴的装备的配置取位置
			pos := pConfig.EquipType
			if pos < 1 || pos > 6 {
				continue
			}

			// 检查是否强化到满级
			pItem := pEquip.UpGrade
			if pItem == nil {
				continue
			}

			// 检查是否达到最大等级
			currentLv := pItem.Lv
			maxLv, ok := GetCsvMgr().EquipUpgradeMaxLv[pItem.Id]
			if !ok {
				continue
			}
			if currentLv >= maxLv {
				continue
			}

			equips[stageIndex] = pEquip
			oldLevel[stageIndex] = pEquip.UpGrade.Lv
			levelStage[stageIndex] = pEquip.UpGrade.Lv
		}
		levelStage[6] = self.player.Sql_UserBase.Level * 2
		sort.Ints(levelStage)

		totalCost := make(map[int]*Item)
		for stageIndex := 0; stageIndex < len(levelStage); stageIndex++ {
			//循环强化,每个阶段结束的标志就是所有装备达到当前阶段的标准
			for {
				//是否有装备强化的标记
				isHas := false
				//循环一轮
				for index := 0; index <= 5; index++ {
					//没有穿装备
					if equips[index] == nil {
						continue
					}
					//是否已经达标
					if equips[index].UpGrade.Lv >= levelStage[stageIndex] {
						continue
					}

					config := GetCsvMgr().GetEquipUpGrade(equips[index].UpGrade.Id, equips[index].UpGrade.Lv)
					if config == nil {
						continue
					}

					costIds := config.CostIds
					costNums := config.CostNums
					if len(costIds) != len(costNums) {
						continue
					}

					for costIndex := range costIds {
						itemId := costIds[costIndex]
						itemNum := costNums[costIndex]
						if itemId == 0 || itemNum == 0 {
							continue
						}
						pItem, ok := totalCost[itemId]
						if !ok {
							totalCost[itemId] = &Item{itemId, itemNum}
						} else {
							pItem.ItemNum += itemNum
						}
					}

					err := self.player.hasItemMapOk(totalCost)
					if err != nil {
						for costIndex := range costIds {
							itemId := costIds[costIndex]
							itemNum := costNums[costIndex]
							if itemId == 0 || itemNum == 0 {
								continue
							}
							totalCost[itemId].ItemNum -= itemNum
						}
						continue
					} else {
						equips[index].UpGrade.Lv++
						isHas = true
						if isUp {
							self.player.HandleTask(EquipUpNumTask, 0, 0, 0)
							self.player.HandleTask(EquipUpTimesTask, 1, 0, 0)
						}
					}
				}

				//没有装备强化了则跳出这个阶段
				if !isHas {
					break
				}
			}
		}

		//区分强化和计算，避免前后不一致
		if isUp {
			//扣除物品
			cost := self.player.RemoveObjectItemMap(totalCost, "一键全体装备强化", 0, 0, 0)
			//比较等级
			var euqipsInfo []*Equip
			for oldLevelIndex := 0; oldLevelIndex < len(oldLevel); oldLevelIndex++ {
				if equips[oldLevelIndex] == nil {
					continue
				}
				if equips[oldLevelIndex].UpGrade.Lv != oldLevel[oldLevelIndex] {
					euqipsInfo = append(euqipsInfo, equips[oldLevelIndex])
				}
			}

			self.player.countHeroIdFight(hero.HeroId, ReasonUpgradeEquipAuto)
			self.player.SendMsg(cid, HF_JtoB(&S2C_EquipAction{
				Cid:        cid,
				Action:     msg.Action,
				EquipInfos: euqipsInfo,
				Items:      cost,
				HeroAtt:    self.getHeroAtt(hero.HeroId),
				HeroId:     msg.HeroId,
			}))

			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_UPGRADE, 0, 0, 0, "一键全体装备强化", 0, 0, self.player)
		} else {
			//只是计算等级要还原
			var euqipsInfo []*Equip
			for oldLevelIndex := 0; oldLevelIndex < len(oldLevel); oldLevelIndex++ {
				if equips[oldLevelIndex] == nil {
					continue
				}
				if equips[oldLevelIndex].UpGrade.Lv != oldLevel[oldLevelIndex] {
					equips[oldLevelIndex].UpGrade.Lv = oldLevel[oldLevelIndex]
				}
			}

			//兼容消息
			itemRel := make([]PassItem, 0)
			for _, item := range totalCost {
				itemRel = append(itemRel, PassItem{item.ItemId, item.ItemNum})
			}

			self.player.SendMsg(cid, HF_JtoB(&S2C_EquipAction{
				Cid:        cid,
				Action:     msg.Action,
				EquipInfos: euqipsInfo,
				Items:      itemRel,
				HeroAtt:    self.getHeroAtt(hero.HeroId),
				HeroId:     msg.HeroId,
			}))
		}


	*/
}

// 装备附魔
func (self *ModEquip) starEquip(cid string, msg *C2S_EquipAction) {
	/*
		pEquip := self.GetEquipItem(msg.KeyId)
		if pEquip == nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_EQUIPMENT_DOES_NOT_EXIST"))
			return
		}

		// 检查是否强化到满级
		pItem := pEquip.StarInfo
		if pItem == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_EQUIP_DATA_EXCEPTION_OF_ENCHANTMENT_MODULE"))
			return
		}

		// 检查是否达到最大等级
		currentLv := pItem.Lv
		maxLv, ok := GetCsvMgr().EquipStarMaxLv[pItem.Id]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_EQUIP_CURRENT_ENCHANTMENT_CONFIGURATION_DOES_NOT"))
			return
		}

		if currentLv >= maxLv {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_EQUIP_HAS_REACHED_THE_MAXIMUM_LEVEL"))
			return
		}

		config := GetCsvMgr().GetEquipStar(pItem.Id, currentLv)
		if config == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_EQUIP_CURRENT_ENCHANTMENT_CONFIGURATION_DOES_NOT"))
			return
		}

		// 检查消耗是否正常
		costIds := config.CostIds
		costNums := config.CostNums
		if err := self.player.HasObjectOk(costIds, costNums); err != nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_INSUFFICIENT_MATERIAL"))
			return
		}

		// 进行升星操作
		items := self.player.RemoveObjectLst(costIds, costNums, "装备附魔", pItem.Lv+1, pEquip.Id, 0)
		pItem.Lv += 1
		heroId := self.GetTeamHero(pEquip.KeyId)
		self.player.countHeroIdFight(heroId, ReasonStarEquip)
		self.player.HandleTask(EquipStarNumTask, 0, 0, 0)
		self.player.HandleTask(EquipStarTimesTask, 1, 0, 0)

		equipHeroId := self.GetTeamHero(pEquip.KeyId)
		self.player.SendMsg(cid, HF_JtoB(&S2C_EquipAction{
			Cid:     cid,
			Action:  msg.Action,
			PItem:   pEquip,
			Items:   items,
			HeroAtt: self.getHeroAtt(equipHeroId),
			HeroId:  msg.HeroId,
		}))

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EQUIP_STAR, pItem.Lv, pEquip.Id, 0, "装备附魔", 0, 0, self.player)

	*/
}

// 增加装备
func (self *ModEquip) AddEquipWithParam(itemId int, num int, param1, param2 int, dec string) {
	if num <= 0 {
		return
	}

	var res []*Equip
	for i := 1; i <= num; i++ {
		pItem := self.NewEquipItem(itemId)
		for j := 0; j < len(self.Data.equipItems); j++ {
			if len(self.Data.equipItems[j]) < EQUIP_PACK_LEN {
				self.Data.equipItems[j][pItem.KeyId] = pItem
				res = append(res, pItem)
				break
			}
		}
	}
	curNum := 0
	for j := 0; j < len(self.Data.equipItems); j++ {
		curNum += len(self.Data.equipItems[j])
	}
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, itemId, num, param1, param2, dec, curNum, 0, self.player)
	self.synEquip(res)

	self.player.HandleTask(TASK_TYPE_HAVE_EQUIP, 0, 0, 0)
}

// 同步宝物信息
func (self *ModEquip) synEquip(res []*Equip) {
	if len(res) <= 0 {
		return
	}

	cid := "synequip"
	self.player.SendMsg(cid, HF_JtoB(&S2C_SynEquip{
		Cid:    cid,
		Equips: res,
	}))
}

// 合成背包中的宝石
func (self *ModEquip) compoundGem(cid string, msg *C2S_GemAction) {
	gemId := msg.GemId
	// 检查配置是否存在
	config := GetCsvMgr().GetGemConfig(gemId)
	if config == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_GEMSTONE_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	itemConfig := GetCsvMgr().GetItemConfig(gemId)
	if itemConfig == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_GEM_PROPS_CONFIGURATION_DOES_NOT"))
		return
	}

	currentNum := self.player.GetObjectNum(gemId)
	if currentNum >= itemConfig.MaxNum {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_BEYOND_THE_UPPER_LIMIT_OF"))
		return
	}

	needId := config.NeedId
	needNum := config.NeedNum
	ownNum := self.player.GetObjectNum(needId)
	if ownNum <= 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_THERE_IS_NO_SUPERIOR_GEMSTONE"))
		return
	}

	leftNum := 0
	if ownNum < needNum {
		leftNum = needNum - ownNum
	}

	needConfig := GetCsvMgr().GetItemConfig(needId)
	if needConfig == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_NEED_CONFIG_GEM_PROJECTS_CONFIGURATION"))
		return
	}

	// 检查需要的钻石
	if leftNum > 0 {
		needMoney := needConfig.GemPrice * leftNum
		if self.player.GetObjectNum(DEFAULT_GEM) < needMoney {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MINE_DIAMOND_SHORTAGE"))
			return
		}
	}

	// 扣除宝石以及钻石
	var items []PassItem
	// 检查是否需要扣除钻石
	if leftNum > 0 {
		needMoney := needConfig.GemPrice * leftNum
		money := self.player.RemoveObjectEasy(DEFAULT_GEM, needMoney, "宝石合成", 0, 0, 0)
		items = append(items, money...)
	}
	// 扣除合成的钻石
	res := self.player.RemoveObjectEasy(needId, needNum, "宝石合成", 0, 0, 0)
	items = append(items, res...)
	// 增加钻石
	addRes := self.player.AddObjectSimple(gemId, 1, "宝石合成", 0, 0, 0)
	items = append(items, addRes...)

	// 合成成功
	self.player.SendMsg(cid, HF_JtoB(&S2C_GemAction{
		Cid:   cid,
		Items: items,
	}))

	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GEM_COMPOUND, config.Id, config.Level, 0, "宝石合成", 0, 0, self.player)
}

func (self *ModEquip) GetTeamHero(keyId int) int {

	info := self.player.GetModule("hero").(*ModHero).Sql_Hero.info

	for _, v := range info {
		for _, key := range v.EquipIds {
			if key == 0 {
				continue
			}
			if keyId == key {
				return v.HeroId
			}
		}
	}
	return 0
}

// 获取英雄装备
func (self *ModEquip) getHeroEquips(heroId int) map[int]*Equip {
	equipItems := self.Data.equipItems
	res := make(map[int]*Equip)
	hero := self.player.getHero(heroId)
	if hero == nil {
		return res
	}
	for _, keyId := range hero.EquipIds {
		for j := 0; j < len(self.Data.equipItems); j++ {
			item, ok := equipItems[j][keyId]
			if !ok {
				continue
			}
			pConfig := GetCsvMgr().GetEquipConfig(item.Id)
			if pConfig == nil {
				continue
			}
			res[pConfig.EquipPosition] = item
			break
		}
	}

	return res
}

// 基础属性
func (self *ModEquip) getBaseAttr(pEquip *Equip) map[int]*Attribute {
	res := make(map[int]*Attribute)

	/*
		for i := 0; i < len(pEquip.BaseAttrInfo); i++ {
			config := GetCsvMgr().EquipBaseValueMap[pEquip.BaseAttrInfo[i]]
			if config != nil {
				for index := range config.BaseTypes {
					attrType := config.BaseTypes[index]
					attrValue := config.BaseValues[index]
					if attrValue == 0 {
						continue
					}
					v, ok := res[attrType]
					if !ok {
						res[attrType] = &Attribute{AttType: attrType, AttValue: attrValue}
					} else {
						v.AttValue += attrValue
					}
				}
			}
		}

		for i := 0; i < len(pEquip.SpecialAttrInfo); i++ {
			config := GetCsvMgr().EquipBaseValueMap[pEquip.SpecialAttrInfo[i]]
			if config != nil {
				for index := range config.BaseTypes {
					attrType := config.BaseTypes[index]
					attrValue := config.BaseValues[index]
					if attrValue == 0 {
						continue
					}
					v, ok := res[attrType]
					if !ok {
						res[attrType] = &Attribute{AttType: attrType, AttValue: attrValue}
					} else {
						v.AttValue += attrValue
					}
				}
			}
		}

	*/
	return res
}

func (self *ModEquip) addAttEx(inparam map[int]*Attribute, attMap map[int]*Attribute) {
	for _, att := range inparam {
		_, ok := attMap[att.AttType]
		if !ok {
			attMap[att.AttType] = &Attribute{
				AttType:  att.AttType,
				AttValue: att.AttValue,
			}
		} else {
			attMap[att.AttType].AttValue += att.AttValue
		}
	}
}

// 强化属性
func (self *ModEquip) getUpgradeAttr(pEquip *Equip) map[int]*Attribute {
	res := make(map[int]*Attribute)
	/*
		upgrade := pEquip.UpGrade
		if upgrade == nil {
			return res
		}

		pConfig := GetCsvMgr().GetEquipConfig(pEquip.Id)
		if pConfig == nil {
			return res
		}

			for index := range pConfig.UpgradeTypes {
				attrType := pConfig.UpgradeTypes[index]
				attrValue := pConfig.UpgradeValues[index]
				if attrValue == 0 {
					continue
				}
				v, ok := res[attrType]
				if !ok {
					res[attrType] = &Attribute{AttType: attrType, AttValue: attrValue * int64(upgrade.Lv-1)}
				} else {
					v.AttValue += attrValue * int64(upgrade.Lv-1)
				}
			}

	*/

	return res
}

// 附魔属性
func (self *ModEquip) getStarAttr(pEquip *Equip) map[int]*Attribute {
	res := make(map[int]*Attribute)
	/*
		star := pEquip.StarInfo
		if star == nil {
			return res
		}

		pConfig := GetCsvMgr().GetEquipConfig(pEquip.Id)
		if pConfig == nil {
			return res
		}

			for index := range pConfig.StarTypes {
				attrType := pConfig.StarTypes[index]
				attrValue := pConfig.StarValues[index]
				if attrValue == 0 {
					continue
				}
				v, ok := res[attrType]
				if !ok {
					res[attrType] = &Attribute{AttType: attrType, AttValue: attrValue * int64(star.Lv)}
				} else {
					v.AttValue += attrValue * int64(star.Lv)
				}
			}

	*/
	return res
}

// 装备基础 + 强化 + 附魔 + 宝石
func (self *ModEquip) getEquipAttr(pItem *Equip, camp int, isCalCamp bool) map[int]*Attribute {
	attMap := make(map[int]*Attribute)

	equipConfig := GetCsvMgr().EquipConfigMap[pItem.Id]
	if equipConfig == nil {
		return attMap
	}

	for i := 0; i < len(equipConfig.BaseTypes); i++ {
		if equipConfig.BaseTypes[i] > 0 && equipConfig.BaseValues[i] != 0 {
			attrType := equipConfig.BaseTypes[i]
			attrValue := int64(0)
			rate := PER_BIT
			if isCalCamp && equipConfig.Camp == camp {
				rate += 3000
			}
			attrValue = equipConfig.BaseValues[i] * int64(rate) / PER_BIT
			if pItem.Lv > 0 {
				configLvUp := GetCsvMgr().GetEquipStrengthenConfig(equipConfig.EquipType, equipConfig.EquipPosition, equipConfig.Quality, pItem.Lv)
				if configLvUp != nil {
					attrValue += int64(configLvUp.Vaual[i])
				}
			}

			_, ok := attMap[attrType]
			if !ok {
				attMap[attrType] = &Attribute{attrType, attrValue}
			} else {
				if attMap[attrType] != nil {
					attMap[attrType].AttValue += attrValue
				} else {
					LogError("nil ptr in getEquipAttr")
				}
			}
		}
	}

	return attMap
}

func (self *ModEquip) getEquipsAttr(heroKeyId int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	// 找出当前武将所有装备, 并计算所有的属性
	heroEquip := self.getHeroEquips(heroKeyId)
	if len(heroEquip) <= 0 {
		return attMap
	}
	hero := self.player.getHero(heroKeyId)
	if hero == nil {
		return attMap
	}
	config := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if config == nil {
		return attMap
	}
	// 装备属性
	for _, pItem := range heroEquip {
		if pItem == nil {
			continue
		}
		totalAtt := self.getEquipAttr(pItem, config.Attribute, true)
		self.addAttEx(totalAtt, attMap)
	}
	return attMap
}

func (self *ModEquip) getEquipsAttrNoCamp(heroKeyId int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	// 找出当前武将所有装备, 并计算所有的属性
	heroEquip := self.getHeroEquips(heroKeyId)
	if len(heroEquip) <= 0 {
		return attMap
	}
	hero := self.player.getHero(heroKeyId)
	if hero == nil {
		return attMap
	}
	config := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if config == nil {
		return attMap
	}
	// 装备属性
	for _, pItem := range heroEquip {
		if pItem == nil {
			continue
		}
		totalAtt := self.getEquipAttr(pItem, config.Attribute, false)
		self.addAttEx(totalAtt, attMap)
	}
	return attMap
}

func (self *ModEquip) getSuitAttr(heroId int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	// 找出当前武将所有装备, 并计算所有的属性
	heroEquip := self.getHeroEquips(heroId)
	if len(heroEquip) <= 0 {
		return attMap
	}

	// 套装属性, suitId, num
	suitMap := make(map[int]int)
	for _, pItem := range heroEquip {
		if pItem == nil {
			continue
		}
		config := GetCsvMgr().GetEquipConfig(pItem.Id)
		if config == nil {
			continue
		}

		//suitMap[config.SuitId] += 1
	}

	for suitId, num := range suitMap {
		totalAtt := GetCsvMgr().GetEquipSuitAttr(suitId, num)
		if len(totalAtt) <= 0 {
			continue
		}
		self.addAttEx(totalAtt, attMap)
	}

	return attMap
}

// 装备基础 + 强化 + 附魔 + 宝石 + 套装属性
func (self *ModEquip) getAttr(heroKeyId int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	equipAttr := self.getEquipsAttr(heroKeyId)
	self.addAttEx(equipAttr, attMap)
	//suitAttr := self.getSuitAttr(heroKeyId)
	//self.addAttEx(suitAttr, attMap)

	return attMap
}

func (self *ModEquip) getAttrNoCamp(heroKeyId int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	equipAttr := self.getEquipsAttrNoCamp(heroKeyId)
	self.addAttEx(equipAttr, attMap)
	//suitAttr := self.getSuitAttr(heroKeyId)
	//self.addAttEx(suitAttr, attMap)

	return attMap
}

//func (self *ModEquip) getAttrExclusive(heroKeyId int, attStar map[int]*Attribute, pAttr *AttrWrapper) {
func (self *ModEquip) getAttrExclusive(heroKeyId int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	hero := self.player.getHero(heroKeyId)
	if hero == nil {
		return attMap
	}

	if hero.ExclusiveEquip.UnLock == LOGIC_FALSE {
		return attMap
	}

	// 如果是虚空英雄但是星级不足 则不增加属性
	if hero.StarItem.UpStar < EQUIP_EXCLUSIVE_OPEN_LV {
		return attMap
	}

	for _, v := range hero.ExclusiveEquip.AttrInfo {
		_, ok := attMap[v.AttrType]
		if !ok {
			attMap[v.AttrType] = &Attribute{
				AttType:  v.AttrType,
				AttValue: v.AttrValue,
			}
		} else {
			attMap[v.AttrType].AttValue += v.AttrValue
		}
	}

	skills := make([]int, 0)
	skills = append(skills, hero.ExclusiveEquip.Skill)
	skillAttr := GetSkillAttr(skills)
	if len(skillAttr) > 0 {
		for _, v := range skillAttr {
			_, ok := attMap[v.AttType]
			if !ok {
				attMap[v.AttType] = &Attribute{
					AttType:  v.AttType,
					AttValue: v.AttValue,
				}
			} else {
				attMap[v.AttType].AttValue += v.AttValue
			}
		}
	}
	return attMap
}

// 登录发送
func (self *ModEquip) SendInfo() {
	self.checkEquip()
	self.player.HandleTask(EquipColorTask, 0, 0, 0)

	//这个CHECK
	for i := 0; i < len(self.Data.equipItems); i++ {
		for _, v := range self.Data.equipItems[i] {
			if v.HeroKeyId == 0 {
				continue
			}
			hero := self.player.getHero(v.HeroKeyId)
			if hero == nil {
				v.HeroKeyId = 0
			}
		}
	}

	var msg S2C_EquipInfo
	msg.Cid = "equipinfo"
	msg.Equips = self.Data.equipItems
	msg.Ver = LOGIC_TRUE
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

// 可以分解
func (self *Equip) isDecompound() error {
	/*
		if self.UpGrade == nil {
			return errors.New(GetCsvMgr().GetText("STR_MOD_EQUIP_UNABLE_TO_DECOMPOSE:_ENHANCE_DATA"))
		}

		if self.StarInfo == nil {
			return errors.New(GetCsvMgr().GetText("STR_MOD_EQUIP_UNABLE_TO_DECOMPOSE:_ASCENDING_STAR"))
		}

	*/

	return nil
}

// 强化重生
func (self *Equip) getUpgradeReborn() map[int]*Item {
	itemMap := make(map[int]*Item)
	/*
		upgrade := self.UpGrade
		if upgrade == nil {
			return itemMap
		}

		if upgrade.Lv <= 1 {
			return itemMap
		}

		config := GetCsvMgr().GetEquipUpGrade(upgrade.Id, upgrade.Lv)
		if config == nil {
			return itemMap
		}
		AddItemMapHelper(itemMap, config.RebornIds, config.RebornNums)

	*/
	return itemMap
}

// 附魔重生, 通过计算
func (self *Equip) getStarReborn() map[int]*Item {
	itemMap := make(map[int]*Item)
	/*
		// 附魔材料
		starInfo := self.StarInfo
		if starInfo == nil {
			return itemMap
		}

		if starInfo.Lv < 1 {
			return itemMap
		}

		config := GetCsvMgr().GetEquipStar(starInfo.Id, starInfo.Lv)
		if config == nil {
			return itemMap
		}

	*/
	return itemMap
}

// 宝石卸载
func (self *Equip) getGemReborn() map[int]*Item {
	itemMap := make(map[int]*Item)
	/*
		for _, gemId := range self.Gems {
			if gemId == 0 {
				continue
			}

			pItem, ok := itemMap[gemId]
			if !ok {
				itemMap[gemId] = &Item{gemId, 1}
			} else {
				pItem.ItemNum += 1
			}
		}

	*/
	return itemMap
}

// 装备重置
func (self *ModEquip) rebornEquips(keyIds []int) []*Equip {
	var euqips []*Equip
	/*
		for _, keyId := range keyIds {
			pInfo := self.GetEquipItem(keyId)
			if pInfo == nil {
				continue
			}

			if pInfo.UpGrade != nil {
				pInfo.UpGrade.Lv = 1
			}

			if pInfo.StarInfo != nil {
				pInfo.StarInfo.Lv = 0
			}

			for i := range pInfo.Gems {
				pInfo.Gems[i] = 0
			}
			euqips = append(euqips, pInfo)
		}

	*/
	return euqips
}

func (self *Equip) ClearGem() {

}

//设置属性
func (self *Equip) setAttr() {
	/*
		self.BaseAttrInfo = make([]int, 0)
		self.SpecialAttrInfo = make([]int, 0)

		configEquip := GetCsvMgr().EquipMap[self.Id]
		if configEquip == nil {
			return
		}

		//计算基础属性
		configBaseAttr := GetCsvMgr().EquipValueGroupMap[configEquip.RandomValue]
		if configBaseAttr == nil {
			return
		}
		for i := 0; i < len(configBaseAttr.Group); i++ {
			if configBaseAttr.Group[i] == 0 {
				continue
			}
			rand := HF_GetRandom(RAND_MAX)
			if rand < configBaseAttr.Chance[i] {
				self.addAttr(configBaseAttr.Type, configBaseAttr.Group[i])
			}
		}

		//特殊属性
		configSpecialAttr := GetCsvMgr().EquipValueGroupMap[configEquip.SpecialValue]
		if configSpecialAttr == nil {
			return
		}
		for i := 0; i < len(configSpecialAttr.Group); i++ {
			if configSpecialAttr.Group[i] == 0 {
				continue
			}
			rand := HF_GetRandom(RAND_MAX)
			if rand < configSpecialAttr.Chance[i] {
				self.addAttr(configSpecialAttr.Type, configSpecialAttr.Group[i])
			}
		}

	*/
}

// 装备强化 装备数量, 达标等级
func (self *ModEquip) getEquipUpNum(n2, n3, n4 int) int {
	/*
		info := self.Data.EquipItems
		num := 0
		for _, value := range info {
			bLevel, bStep, bType := false, false, false
			if value == nil {
				continue
			}

			if value.UpGrade == nil {
				continue
			}

			config := GetCsvMgr().GetEquipConfig(value.Id)
			if config == nil {
				continue
			}

			if n2 == 0 {
				bLevel = true
			} else {
				if value.UpGrade.Lv >= n2 {
					bLevel = true
				}
			}

			if n3 == 0 {
				bStep = true
			} else {
				if config.Quality >= n3 {
					bStep = true
				}
			}

			if n4 == 0 {
				bType = true
			} else {
				if config.EquipType >= n4 {
					bType = true
				}
			}

			if bLevel && bStep && bType {
				num += 1
			}

		}
		return num

	*/
	return 0
}

// 装备附魔 装备数量,达标等级
func (self *ModEquip) getEquipStarNum(n int) int {
	//info := self.Data.EquipItems
	num := 0
	/*
		for _, value := range info {
			if value == nil {
				continue
			}

			if value.StarInfo == nil {
				continue
			}

			if value.StarInfo.Lv >= n {
				num += 1
			}
		}
	*/
	return num
}

func (self *Hero) checkEquip() {
	if len(self.EquipIds) == 0 {
		for i := 0; i < EQUIPPOS_END-EQUIPPOS_1; i++ {
			self.EquipIds = append(self.EquipIds, 0)
		}
	}

	if len(self.ArtifactEquipIds) == 0 {
		for i := 0; i < ARTIFACT_EQUIPPOS_END-ARTIFACT_EQUIPPOS_1; i++ {
			self.ArtifactEquipIds = append(self.ArtifactEquipIds, 0)
		}
	}
	//专属装备
	if self.ExclusiveEquip == nil {
		self.ExclusiveEquip = self.NewExclusiveEquipItem()
	}
}

func (self *ModHero) OffEquipAll(hero *Hero) ([]*Equip, []*ArtifactEquip, []*JS_HorseInfo) {
	equipRel := make([]*Equip, 0)
	artifactRel := make([]*ArtifactEquip, 0)
	horseInfoRel := make([]*JS_HorseInfo, 0)
	if hero == nil {
		return equipRel, artifactRel, horseInfoRel
	}
	equips := self.player.getEquips()
	for i := 0; i < len(equips); i++ {
		if equips[i] == nil {
			return equipRel, artifactRel, horseInfoRel
		}
	}
	for i := 0; i < len(hero.EquipIds); i++ {
		if hero.EquipIds[i] == 0 {
			continue
		}
		for j := 0; j < len(equips); j++ {
			if equips[j][hero.EquipIds[i]] != nil {
				equips[j][hero.EquipIds[i]].HeroKeyId = 0
				equipRel = append(equipRel, equips[j][hero.EquipIds[i]])
			}
			hero.EquipIds[i] = 0
			break
		}
	}
	for i := 0; i < len(hero.ArtifactEquipIds); i++ {
		if hero.ArtifactEquipIds[i] == 0 {
			continue
		}
		artifact := self.player.GetModule("artifactequip").(*ModArtifactEquip).GetArtifactEquipItem(hero.ArtifactEquipIds[i])
		if artifact != nil {
			artifact.HeroKeyId = 0
			artifactRel = append(artifactRel, artifact)
		}
		hero.ArtifactEquipIds[i] = 0
	}
	if hero.Horse > 0 {
		horse := self.player.GetModule("horse").(*ModHorse).GetHorse(hero.Horse)
		if horse != nil {
			horseInfoRel = append(horseInfoRel, horse)

			horse.chg = true
			horse.Heroid = 0
			//! 卸下马魂
			soullst := make([]JS_HorseSoulInfo, 0)
			for j := 0; j < len(horse.SoulLst); j++ {
				if horse.SoulLst[j] > 0 {
					soulid := horse.SoulLst[j] / 100
					soulrank := horse.SoulLst[j] % 100
					soullst = append(soullst, JS_HorseSoulInfo{soulid, soulrank, 1})
					self.player.GetModule("horse").(*ModHorse).AddHorseSoul(soulid, soulrank, 1)
				}
			}
			soulnum := len(horse.SoulLst)
			horse.SoulLst = make([]int, soulnum)
		}
		hero.Horse = 0
	}
	return equipRel, artifactRel, horseInfoRel
}

func (self *ModEquip) UnLockExclusive(msg *C2S_ExclusiveAction) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroKeyId))
		return
	}

	if hero.ExclusiveEquip.UnLock == LOGIC_TRUE {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EXCLUSIVE_ALREADY_UNLOCK"))
		return
	}

	if hero.StarItem.UpStar < EQUIP_EXCLUSIVE_OPEN_LV {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EXCLUSIVE_UPSTAR_NOT_ENOUGH"))
		return
	}

	config := GetCsvMgr().ExclusiveEquipConfigMap[hero.ExclusiveEquip.Id]
	if config == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EXCLUSIVE_CONFIG_ERROR"))
		return
	}

	//检查消耗
	if err := self.player.HasObjectOkEasy(config.ActiveNeed, config.ActiveNum); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	hero.ExclusiveEquip.UnLock = LOGIC_TRUE
	hero.ExclusiveEquip.CalAttr()
	self.player.countHeroFight(hero, ReasonExclusiveUnlock)
	self.player.GetModule("hero").(*ModHero).CheckHireUpdate(hero)
	costItem := self.player.RemoveObjectSimple(config.ActiveNeed, config.ActiveNum, "专属激活", hero.ExclusiveEquip.Id, hero.HeroId, 0)

	var msgRel S2C_ExclusiveAction
	msgRel.Cid = "exclusiveaction"
	msgRel.Action = msg.Action
	msgRel.ExclusiveItem = hero.ExclusiveEquip
	msgRel.CostItems = costItem
	msgRel.HeroKeyId = msg.HeroKeyId
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EXCLUSIVE_UNLOCK, config.Id, hero.HeroId, hero.HeroKeyId, "专属激活", 0, 0, self.player)
}

func (self *ModEquip) LvUpExclusive(msg *C2S_ExclusiveAction) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()
	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroKeyId))
		return
	}

	if hero.ExclusiveEquip.UnLock == LOGIC_FALSE {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EXCLUSIVE_NOT_UNLOCK"))
		return
	}

	config := GetCsvMgr().ExclusiveStrengthen[hero.ExclusiveEquip.Id][hero.ExclusiveEquip.Lv+1]
	if config == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EXCLUSIVE_LVUP_MAX"))
		return
	}

	//生成消耗
	costId := make([]int, 0)
	costNum := make([]int, 0)
	if config.Replace > 0 {
		nowNum := self.player.GetObjectNum(config.Need)
		if nowNum >= config.Num {
			costId = append(costId, config.Need)
			costNum = append(costNum, config.Num)
		} else {
			costId = append(costId, config.Need)
			costNum = append(costNum, nowNum)

			costId = append(costId, config.Replace)
			costNum = append(costNum, config.Num-nowNum)
		}
	} else {
		costId = append(costId, config.Need)
		costNum = append(costNum, config.Num)
	}

	//检查消耗
	if err := self.player.HasObjectOk(costId, costNum); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	costItem := self.player.RemoveObjectLst(costId, costNum, "专属强化", hero.ExclusiveEquip.Id, hero.HeroId, 0)

	oldlevel := hero.ExclusiveEquip.Lv
	hero.ExclusiveEquip.Lv += 1
	hero.ExclusiveEquip.CalAttr()
	self.player.countHeroFight(hero, ReasonExclusiveLvUp)
	self.player.GetModule("hero").(*ModHero).CheckHireUpdate(hero)

	var msgRel S2C_ExclusiveAction
	msgRel.Cid = "exclusiveaction"
	msgRel.Action = msg.Action
	msgRel.ExclusiveItem = hero.ExclusiveEquip
	msgRel.CostItems = costItem
	msgRel.HeroKeyId = msg.HeroKeyId
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EXCLUSIVE_LEVEL_UP, config.Id, hero.HeroId, oldlevel, "专属强化", 0, hero.ExclusiveEquip.Lv, self.player)
}

func (self *ModEquip) getEquipNumByQuality(quality int, level int) int {
	num := 0
	for j := 0; j < len(self.Data.equipItems); j++ {
		for _, v := range self.Data.equipItems[j] {
			pConfig, ok := GetCsvMgr().EquipConfigMap[v.Id]
			if !ok {
				continue
			}

			if pConfig.Quality < quality {
				continue
			}

			if v.Lv < level {
				continue
			}

			num++
		}
	}
	return num
}

func (self *ModEquip) GmClearItem() {
	for i := 0; i < len(self.Data.equipItems); i++ {
		self.Data.equipItems[i] = make(map[int]*Equip, 0)
	}
	self.SendInfo()
}
