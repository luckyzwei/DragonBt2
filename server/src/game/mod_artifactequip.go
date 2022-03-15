package game

import (
	"encoding/json"
	"fmt"
)

//信息交互中的装备位置索引
const (
	ARTIFACT_EQUIPPOS_START = iota
	ARTIFACT_EQUIPPOS_1
	ARTIFACT_EQUIPPOS_END
)

const (
	wearArtifactEquip    = 1 // 穿戴
	takeOffArtifactEquip = 2 // 脱
)

// 服务器存放的装备信息
type ArtifactEquip struct {
	KeyId     int         `json:"keyid"`     //! 装备唯一Id
	Id        int         `json:"id"`        //! 装备配置Id
	HeroKeyId int         `json:"herokeyid"` //! 装备拥有者
	AttrInfo  []*AttrInfo `json:"attrinfo"`  //! 属性信息
	Lv        int         `json:"lv"`        //!
}

// 神器
type San_ArtifactEquip struct {
	Uid                int64                  // 玩家Id
	Maxkey             int                    // 神器最大keyId
	Info               string                 // 神器
	ArtifactEquipItems map[int]*ArtifactEquip // key: 神器唯一Id, value:神器信息

	DataUpdate
}

// 神器系统
type ModArtifactEquip struct {
	player *Player           //! 玩家
	Data   San_ArtifactEquip //! 数据库数据
}

// 神器唯一Id
func (self *ModArtifactEquip) MaxKey() int {
	self.Data.Maxkey += 1
	return self.Data.Maxkey
}

// save
func (self *ModArtifactEquip) Decode() {
	json.Unmarshal([]byte(self.Data.Info), &self.Data.ArtifactEquipItems)
}

// read
func (self *ModArtifactEquip) Encode() {
	self.Data.Info = HF_JtoA(self.Data.ArtifactEquipItems)
}

// get player and init map
func (self *ModArtifactEquip) OnGetData(player *Player) {
	self.player = player
	self.checkArtifactEquip()
	tableName := self.getTableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", tableName, self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Data, tableName, self.player.ID)
	if self.Data.Uid <= 0 {
		self.init(self.player.ID)
		self.Encode()
		InsertTable(tableName, &self.Data, 0, true)
	} else {
		self.Decode()
		self.checkArtifactEquip()
	}

	self.Data.Init(tableName, &self.Data, true)
}

// get data from db
func (self *ModArtifactEquip) OnGetOtherData() {
	self.checkArtifactEquip()
}

func (self *ModArtifactEquip) getTableName() string {
	return "san_userartifactequip"
}

func (self *ModArtifactEquip) init(uid int64) {
	self.Data.Uid = uid
	self.checkArtifactEquip()
}

func (self *ModArtifactEquip) checkArtifactEquip() {
	if self.Data.ArtifactEquipItems == nil {
		self.Data.ArtifactEquipItems = make(map[int]*ArtifactEquip, 0)
	}

	//神器初始化
	/*
		for _, v := range GetCsvMgr().ArtifactEquipConfigMap {
			isHas := false
			for _, value := range self.Data.ArtifactEquipItems {
				if value.Id == v.ArtifactId {
					isHas = true
					break
				}
			}
			if !isHas {
				self.AddArtifactWithParam(v.ArtifactId, 1, 0, 0, "神器初始化")
			}
		}
	*/
}

// save db every five minutes by changes
func (self *ModArtifactEquip) OnSave(sql bool) {
	self.Encode()
	self.Data.Update(sql)
}

func (self *ModArtifactEquip) onReg(handlers map[string]func(body []byte)) {
	handlers["artifactequipaction"] = self.onArtifactEquipAction
	handlers["artifactequipuplv"] = self.ArtifactEquipUpLv
}

func (self *ModArtifactEquip) ArtifactEquipUpLv(body []byte) {
	var msg C2S_ArtifactEquipUpLv
	json.Unmarshal(body, &msg)

	pEquip := self.GetArtifactEquipItem(msg.ArtifactEquipKeyId)
	if pEquip == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_ARTIFACT_EQUIPMENT_DOES_NOT_EXIST"))
		return
	}
	//看看是否达到上限
	config := GetCsvMgr().GetArtifactStrengthenLvUpConfig(pEquip.Id, pEquip.Lv+1)
	if config == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_ARTIFACT_LV_MAX"))
		return
	}

	//看看消耗是否充足
	err := self.player.HasObjectOk(config.Need, config.Num)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	items := self.player.RemoveObjectLst(config.Need, config.Num, "神器强化", pEquip.Id, 0, 0)
	oldlevel := pEquip.Lv
	pEquip.Lv += 1
	pEquip.CalAttr()

	hero := self.player.getHero(pEquip.HeroKeyId)
	if hero != nil {
		self.player.countHeroFight(hero, ReasonArtifactWear)
	}

	var msgRel S2C_ArtifactEquipUpLv
	msgRel.Cid = "artifactequipuplv"
	msgRel.ArtifactEquipItem = pEquip
	msgRel.CostItems = items
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ARTIFACT_LEVEL_UP, pEquip.Lv, oldlevel, pEquip.KeyId, "神器强化", 0, 0, self.player)
}

// 增加神器
func (self *ModArtifactEquip) AddArtifactWithParam(itemId int, num int, param1, param2 int, dec string) {
	if num <= 0 {
		return
	}

	if num > maxEquipNum {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_EQUIP_ONLY_500_ADDITIONAL_EQUIPMENT_CAN"))
		return
	}

	isHas := false
	for _, value := range self.Data.ArtifactEquipItems {
		if value.Id == itemId {
			isHas = true
			break
		}
	}
	if isHas {
		return
	}

	var res []*ArtifactEquip
	for i := 1; i <= num; i++ {
		pItem := self.NewArtifactEquipItem(itemId)
		self.Data.ArtifactEquipItems[pItem.KeyId] = pItem
		res = append(res, pItem)
	}
	curNum := len(self.Data.ArtifactEquipItems)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, itemId, num, param1, param2, dec, curNum, 0, self.player)

	self.player.HandleTask(TASK_TYPE_ACTIFACT_GET, itemId, 0, 0)
	self.synArtifactEquip(res)
}

// 同步宝物信息
func (self *ModArtifactEquip) synArtifactEquip(res []*ArtifactEquip) {
	if len(res) <= 0 {
		return
	}

	cid := "synartifactequip"
	self.player.SendMsg(cid, HF_JtoB(&S2C_SynActifactEquip{
		Cid:            cid,
		ArtifactEquips: res,
	}))
}

// 创建一个装备
func (self *ModArtifactEquip) NewArtifactEquipItem(equipId int) *ArtifactEquip {
	pItem := &ArtifactEquip{}
	pItem.KeyId = self.MaxKey()
	pItem.Id = equipId

	config, ok := GetCsvMgr().ArtifactEquipConfigMap[equipId]
	if !ok {
		LogError("artifact is not found, artifact=", equipId)
		return nil
	}
	//生成属性
	for i := 0; i < len(config.BaseTypes); i++ {
		if config.BaseTypes[i] > 0 {
			attr := new(AttrInfo)
			attr.AttrId = i + 1
			pItem.AttrInfo = append(pItem.AttrInfo, attr)
		}
	}

	pItem.CalAttr()
	return pItem
}

// 装备计算属性
func (self *ArtifactEquip) CalAttr() {
	config := GetCsvMgr().ArtifactEquipConfigMap[self.Id]
	if config == nil {
		return
	}
	for _, v := range self.AttrInfo {
		if v.AttrId <= 0 || v.AttrId > len(config.BaseTypes) {
			continue
		}
		index := v.AttrId - 1
		v.AttrType = config.BaseTypes[index]
		v.AttrValue = config.BaseValues[index]

		if self.Lv > 0 {
			configLvUp := GetCsvMgr().GetArtifactStrengthenLvUpConfig(self.Id, self.Lv)
			if configLvUp != nil {
				v.AttrValue += configLvUp.Value[index]
			}
		}
	}
}

func (self *ModArtifactEquip) getAttr(heroKeyId int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	hero := self.player.getHero(heroKeyId)
	if hero == nil {
		return attMap
	}

	if len(hero.ArtifactEquipIds) == 0 {
		return attMap
	}
	realIndex := ARTIFACT_EQUIPPOS_1 - 1
	pItem := self.GetArtifactEquipItem(hero.ArtifactEquipIds[realIndex])
	if pItem == nil {
		return attMap
	}

	for _, v := range pItem.AttrInfo {
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

	configArt := GetCsvMgr().GetArtifactStrengthenLvUpConfig(pItem.Id, pItem.Lv)
	if configArt != nil {
		skillAttr := GetSkillAttr(configArt.Skill)
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
	}

	return attMap
}

// 消息处理
func (self *ModArtifactEquip) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 装备相关操作
func (self *ModArtifactEquip) onArtifactEquipAction(body []byte) {
	var msg C2S_ArtifactEquipAction
	json.Unmarshal(body, &msg)
	if msg.Action == wearArtifactEquip {
		self.wearArtifactEquip(&msg)
	} else if msg.Action == takeOffArtifactEquip {
		self.takeOffArtifactEquip(&msg)
	}
}

// 获取装备信息
func (self *ModArtifactEquip) GetArtifactEquipItem(keyId int) *ArtifactEquip {
	pItem, ok := self.Data.ArtifactEquipItems[keyId]
	if ok {
		return pItem
	}
	return nil
}

func (self *ModArtifactEquip) wearArtifactEquip(msg *C2S_ArtifactEquipAction) {

	// 需要穿戴的装备
	keyId := msg.ArtifactEquipKeyId
	pItem := self.GetArtifactEquipItem(keyId)
	if pItem == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_ARTIFACT_EQUIPMENT_DOES_NOT_EXIST"))
		return
	}

	itemId := pItem.Id
	_, ok := GetCsvMgr().ArtifactEquipConfigMap[itemId]
	if !ok {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_ARTIFACT_EQUIPMENT_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	// 通过需要穿戴的装备的配置取位置  目前只有1件装备
	pos := 1
	if pos <= EQUIPPOS_START || pos >= EQUIPPOS_END {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_ARTIFACT_ERROR_IN_EQUIPMENT_TYPE"))
		return
	}

	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroKeyId))
		return
	}

	if len(hero.ArtifactEquipIds) != ARTIFACT_EQUIPPOS_END-ARTIFACT_EQUIPPOS_1 {
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
	oldKeyId := hero.ArtifactEquipIds[pos-1]
	// 如果是换装备 当前位置上有装备
	if oldKeyId != 0 {
		hero.ArtifactEquipIds[pos-1] = 0
	}

	// 是否换下的装备装备到另一个英雄上
	bExchange := false
	oldHeroId := 0
	if pItem.HeroKeyId != 0 {
		oldHeroId = pItem.HeroKeyId
		oldHero := self.player.getHero(pItem.HeroKeyId)
		if oldKeyId != 0 {
			oldHero.ArtifactEquipIds[pos-1] = oldKeyId
			bExchange = true
		} else {
			oldHero.ArtifactEquipIds[pos-1] = 0
		}
	}

	// 如果是简单的换装备
	if hero.ArtifactEquipIds[pos-1] == 0 {
		hero.ArtifactEquipIds[pos-1] = keyId
	}

	oldEquip := self.GetArtifactEquipItem(oldKeyId)
	if oldEquip != nil {
		// 换下来的装备需要换位置装备时
		if bExchange {
			oldEquip.HeroKeyId = oldHeroId
			oldHero := self.player.getHero(pItem.HeroKeyId)
			self.player.countHeroFight(oldHero, ReasonArtifactWear)
		} else {
			oldEquip.HeroKeyId = 0
		}
	}

	self.player.countHeroFight(hero, ReasonArtifactWear)
	self.player.GetModule("hero").(*ModHero).CheckHireUpdate(hero)
	pItem.HeroKeyId = msg.HeroKeyId

	self.player.SendMsg("artifactequipaction", HF_JtoB(&S2C_ArtifactEquipAction{
		Cid:               "artifactequipaction",
		Action:            msg.Action,
		ArtifactEquipItem: pItem,
		ExchangeItem:      self.GetArtifactEquipItem(oldKeyId),
	}))

	if oldKeyId != 0 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ARTIFACT_CHANGE, pItem.Id, oldEquip.Id, 0, "神器替换", 0, 0, self.player)
	}
}

func (self *ModArtifactEquip) takeOffArtifactEquip(msg *C2S_ArtifactEquipAction) {
	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST") + fmt.Sprintf("%d", msg.HeroKeyId))
		return
	}

	if len(hero.ArtifactEquipIds) != ARTIFACT_EQUIPPOS_END-ARTIFACT_EQUIPPOS_1 {
		self.player.SendErr(GetCsvMgr().GetText("STR_HERO_NOT_EXIST"))
		return
	}

	pItem := self.GetArtifactEquipItem(msg.ArtifactEquipKeyId)
	if pItem == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_ARTIFACT_NO_EQUIPMENT_DETACHABLE"))
		return
	}

	//目前有且只有一件，方便以后扩展
	hero.ArtifactEquipIds[0] = 0

	pItem.HeroKeyId = 0
	self.player.countHeroFight(hero, ReasonArtifactOff)
	self.player.GetModule("hero").(*ModHero).CheckHireUpdate(hero)

	self.player.SendMsg("artifactequipaction", HF_JtoB(&S2C_ArtifactEquipAction{
		Cid:               "artifactequipaction",
		Action:            msg.Action,
		ArtifactEquipItem: pItem,
	}))
}

// 登录发送
func (self *ModArtifactEquip) SendInfo() {
	self.checkArtifactEquip()
	var msg S2C_ArtifactEquipInfo
	msg.Cid = "artifactequipinfo"
	msg.ArtifactEquips = self.Data.ArtifactEquipItems
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}
