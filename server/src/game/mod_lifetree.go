package game

import (
	"encoding/json"
	"fmt"
)

const (
	MSG_LIFE_TREE_UP_MAIN_LEVEL      = "life_tree_up_main_level"      // 升级主等级
	MSG_LIFE_TREE_UP_TYPE_LEVEL      = "life_tree_up_type_level"      // 升级子等级
	MSG_LIFE_TREE_RESET_TYPE_LEVEL   = "life_tree_reset_type_level"   // 重置子等级
	MSG_LIFE_TREE_GET_AWARD          = "life_tree_get_award"          // 领取奖励
	MSG_LIFE_TREE_SEND_INFO          = "life_tree_send_info"          // 发送信息
	MSG_LIFE_TREE_SEND_AWARD_INFO    = "life_tree_send_award_info"    // 发送奖励累计信息
	MSG_LIFE_TREE_SEND_REDPOINT_INFO = "life_tree_send_redpoint_info" // 发送红点信息
)
const HERO_ALL_FINALL_LEVEL = 300
const HERO_RESET_ITEM = 41000003
const (
	LIFE_TREE_TYPE_1 = 41 // 41 坦克
	LIFE_TREE_TYPE_2 = 42 // 42 战士
	LIFE_TREE_TYPE_3 = 43 // 43 法师
	LIFE_TREE_TYPE_4 = 44 // 44 游侠
	LIFE_TREE_TYPE_5 = 45 // 45 辅助
)

type JS_LifeTreeInfo struct {
	MainLevel int            `json:"mainlevel"`
	Info      []*JS_LifeTree `json:"info"`
}

// 生命树层级
type JS_LifeTree struct {
	Type  int `json:"type"`  // 类型
	Level int `json:"level"` // 等级
}

// 英雄领取记录
type JS_HeroGet struct {
	HeroID int `json:"heroid"` // 英雄星级
	Star   int `json:"star"`   // 星级
}

// 英雄奖励物品记录
type JS_HeroAward struct {
	HeroID int        `json:"heroid"` // 英雄星级
	Star   int        `json:"star"`   // 星级
	Award  []PassItem `json:"award"`  // 奖励物品
}

// 数据结构
type San_LifeTree struct {
	Uid       int64  // 角色ID
	Info      string // 职业分支
	Award     string // 奖励记录
	HeroGet   string // 获取记录
	MainLevel int    // 主等级
	IsGet     int    // 是否获得初次结算

	info    []*JS_LifeTree // 职业分支
	award   []*JS_HeroAward
	heroget []*JS_HeroGet // 领取记录
	DataUpdate
}

// 羁绊
type ModLifeTree struct {
	player       *Player
	San_LifeTree San_LifeTree
}

// 获得数据
func (self *ModLifeTree) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_lifetree` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.San_LifeTree, "san_lifetree", self.player.ID)

	if self.San_LifeTree.Uid <= 0 {
		self.San_LifeTree.Uid = self.player.ID
		self.San_LifeTree.info = make([]*JS_LifeTree, 0)
		self.San_LifeTree.award = make([]*JS_HeroAward, 0)
		self.San_LifeTree.heroget = make([]*JS_HeroGet, 0)
		self.San_LifeTree.MainLevel = 0
		self.Encode()
		InsertTable("san_lifetree", &self.San_LifeTree, 0, true)
		self.San_LifeTree.Init("san_lifetree", &self.San_LifeTree, true)
	} else {
		self.Decode()
		self.San_LifeTree.Init("san_lifetree", &self.San_LifeTree, true)
	}
}

// 上线
func (self *ModLifeTree) OnGetOtherData() {
	self.CheckInfo()
}

// save
func (self *ModLifeTree) Decode() {
	json.Unmarshal([]byte(self.San_LifeTree.Info), &self.San_LifeTree.info)
	json.Unmarshal([]byte(self.San_LifeTree.Award), &self.San_LifeTree.award)
	json.Unmarshal([]byte(self.San_LifeTree.HeroGet), &self.San_LifeTree.heroget)
}

// read
func (self *ModLifeTree) Encode() {
	self.San_LifeTree.Info = HF_JtoA(self.San_LifeTree.info)
	self.San_LifeTree.Award = HF_JtoA(self.San_LifeTree.award)
	self.San_LifeTree.HeroGet = HF_JtoA(self.San_LifeTree.heroget)
}

// 储存
func (self *ModLifeTree) OnSave(sql bool) {
	self.Encode()
	self.San_LifeTree.Update(sql)
}

// 检查数据结构
func (self *ModLifeTree) CheckInfo() {
	// 循环配置
	if len(self.San_LifeTree.info) != len(GetCsvMgr().TreeProfessionalMapConfig) {
		for i, _ := range GetCsvMgr().TreeProfessionalMapConfig {
			data := self.GetInfoData(i)
			if data == nil {
				self.San_LifeTree.info = append(self.San_LifeTree.info, &JS_LifeTree{i, 0})
			}
		}
	}
}

// 获得专业数据
func (self *ModLifeTree) GetInfoData(nType int) *JS_LifeTree {
	for _, v := range self.San_LifeTree.info {
		if v.Type == nType {
			return v
		}
	}
	return nil
}

// 获得英雄领取数据
func (self *ModLifeTree) GetHeroGet(heroID int, star int) *JS_HeroGet {
	for _, v := range self.San_LifeTree.heroget {
		if v.HeroID == heroID && v.Star == star {
			return v
		}
	}
	return nil
}

// 获得累计奖励物品
func (self *ModLifeTree) GetHeroAward(heroID int) *JS_HeroAward {
	for _, v := range self.San_LifeTree.award {
		if v.HeroID == heroID {
			return v
		}
	}
	return nil
}

//更新英雄
func (self *ModLifeTree) UpdateHero() {
	for _, value := range self.player.GetModule("hero").(*ModHero).Sql_Hero.info {
		self.player.checkHeroFight(value, ReasonHeroLifeTree)
	}
}

// 关卡解锁结算
func (self *ModLifeTree) CheckPass() {
	// 是否解锁
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return
	}

	// 已经结算了
	if self.San_LifeTree.IsGet != 0 {
		return
	}

	heromax := make(map[int]int)
	for _, value := range self.player.GetModule("hero").(*ModHero).Sql_Hero.info {
		_, ok := heromax[value.HeroId]
		if ok {
			if value.StarItem.UpStar > heromax[value.HeroId] {
				heromax[value.HeroId] = value.StarItem.UpStar
			}
		} else {
			heromax[value.HeroId] = value.StarItem.UpStar
		}
	}

	for heroid, value := range heromax {
		itemMap := make(map[int]*Item, 0)
		for i := value; i >= 0; i-- {
			self.San_LifeTree.heroget = append(self.San_LifeTree.heroget, &JS_HeroGet{heroid, i})

			// 获得配置
			config := GetCsvMgr().GetHeroMapConfig(heroid, i)
			if config == nil {
				continue
			}

			if config.TreeNum == 0 {
				continue
			}

			// 物品累计
			AddItemMapHelper3(itemMap, config.TreeID, config.TreeNum)

		}

		if len(itemMap) <= 0 {
			continue
		}

		itemArray := []PassItem{}
		for _, v := range itemMap {
			itemArray = append(itemArray, PassItem{v.ItemId, v.ItemNum})
		}
		self.San_LifeTree.award = append(self.San_LifeTree.award, &JS_HeroAward{heroid, value, itemArray}) // 设置奖励
	}

	self.San_LifeTree.IsGet = 1 // 设置结算完成
}

// 获得升级后获得物品的数量
func (self *ModLifeTree) GetItemCount(heroID, star int) *PassItem {
	// 是否解锁
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return nil
	}

	// 取出之前的累计
	itemMap := make(map[int]*Item, 0)
	for i := star; i >= 0; i-- {
		// 未领过
		get := self.GetHeroGet(heroID, i)
		if get == nil {
			config := GetCsvMgr().GetHeroMapConfig(heroID, i)
			if config == nil {
				continue
			}

			if config.TreeNum == 0 {
				continue
			}

			AddItemMapHelper3(itemMap, config.TreeID, config.TreeNum)
		}
	}

	_, ok := itemMap[LIFE_TREE_ITEM]
	if ok {
		return &PassItem{itemMap[LIFE_TREE_ITEM].ItemId, itemMap[LIFE_TREE_ITEM].ItemNum}
	}

	return nil
}

// 升星结算
func (self *ModLifeTree) StarUp(heroID, star int) {
	// 是否解锁
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return
	}

	// 取出之前的累计
	itemMap := make(map[int]*Item, 0)
	award := self.GetHeroAward(heroID)
	if award != nil {
		AddItemMapHelper2(itemMap, award.Award)
	}

	for i := star; i >= 0; i-- {
		// 未领过
		get := self.GetHeroGet(heroID, i)
		if get == nil {
			self.San_LifeTree.heroget = append(self.San_LifeTree.heroget, &JS_HeroGet{heroID, i})
			config := GetCsvMgr().GetHeroMapConfig(heroID, i)
			if config == nil {
				continue
			}

			if config.TreeNum == 0 {
				continue
			}

			AddItemMapHelper3(itemMap, config.TreeID, config.TreeNum)
		}
	}

	if len(itemMap) <= 0 {
		return
	}

	if award == nil {
		award = &JS_HeroAward{heroID, star, nil}
		self.San_LifeTree.award = append(self.San_LifeTree.award, award) // 设置奖励
	}

	if star > award.Star {
		award.Star = star
	}

	itemArray := []PassItem{}
	for _, v := range itemMap {
		if v.ItemNum == 0 {
			continue
		}
		itemArray = append(itemArray, PassItem{v.ItemId, v.ItemNum})
	}
	award.Award = itemArray // 设置奖励
}

// 获得所有属性
func (self *ModLifeTree) GetAllProperty(heroID int, star int) map[int]*Attribute {
	ret := map[int]*Attribute{}
	mainConfig := GetCsvMgr().GetTreeLevelConfig(self.San_LifeTree.MainLevel)
	if nil == mainConfig {
		return ret
	}

	nLen := len(mainConfig.Value)
	if nLen == len(mainConfig.Num) {
		for i := 0; i < nLen; i++ {
			_, ok := ret[mainConfig.Value[i]]
			if ok {
				ret[mainConfig.Value[i]].AttValue += mainConfig.Num[i]
			} else {
				ret[mainConfig.Value[i]] = &Attribute{mainConfig.Value[i], mainConfig.Num[i]}
			}
		}
	}

	heroConfig := GetCsvMgr().GetHeroMapConfig(heroID, star)
	if heroConfig == nil {
		return ret
	}

	// 只有最高级英雄享受加成
	if heroConfig.FinalLevel != HERO_ALL_FINALL_LEVEL {
		return ret
	}

	bookConfig := GetCsvMgr().HeroHandBookConfigMap[heroID]
	if bookConfig == nil {
		return ret
	}

	data := self.GetInfoData(bookConfig.Showfdw)
	if data == nil {
		return ret
	}

	typeConfig := GetCsvMgr().GetTreeProfessionalConfig(bookConfig.Showfdw, data.Level)
	if nil == typeConfig {
		return ret
	}

	nLen = len(typeConfig.Value)
	if nLen == len(typeConfig.Num) {
		for i := 0; i < nLen; i++ {
			_, ok := ret[typeConfig.Value[i]]
			if ok {
				ret[typeConfig.Value[i]].AttValue += typeConfig.Num[i]
			} else {
				ret[typeConfig.Value[i]] = &Attribute{typeConfig.Value[i], typeConfig.Num[i]}
			}
		}
	}

	return ret
}

// 获得所有技能
func (self *ModLifeTree) GetAllSkill(heroID int, star int) []int {
	ret := []int{}
	bookConfig := GetCsvMgr().HeroHandBookConfigMap[heroID]
	if bookConfig == nil {
		return ret
	}

	heroConfig := GetCsvMgr().GetHeroMapConfig(heroID, star)
	if heroConfig == nil {
		return ret
	}

	// 只有最高级英雄享受加成
	if heroConfig.FinalLevel != HERO_ALL_FINALL_LEVEL {
		return ret
	}

	data := self.GetInfoData(bookConfig.Showfdw)
	if data == nil {
		return ret
	}

	level := data.Level
	for i := level; i >= 0; i-- {
		typeConfig := GetCsvMgr().GetTreeProfessionalConfig(bookConfig.Showfdw, i)
		if nil == typeConfig {
			continue
		} else {
			for j := 0; j < len(typeConfig.Skill); j++ {
				if typeConfig.Skill[j] == 0 {
					continue
				}
				ret = append(ret, typeConfig.Skill[j])
			}
			break
		}
	}
	return ret
}

func (self *ModLifeTree) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (self *ModLifeTree) onReg(handlers map[string]func(body []byte)) {
	handlers[MSG_LIFE_TREE_UP_MAIN_LEVEL] = self.MsgUpMainLevel       // 提升主等级
	handlers[MSG_LIFE_TREE_UP_TYPE_LEVEL] = self.MsgUpTypeLevel       // 提升专业等级
	handlers[MSG_LIFE_TREE_RESET_TYPE_LEVEL] = self.MsgResetTypeLevel // 重置专业等级
	handlers[MSG_LIFE_TREE_GET_AWARD] = self.MsgGetAward              // 领取奖励
	handlers[MSG_LIFE_TREE_SEND_AWARD_INFO] = self.MsgSendAwardInfo   // 发送奖励累计信息
}

// 主等级提升
func (self *ModLifeTree) SendInfo(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return
	}

	var backmsg S2C_LifeTreeSendInfo
	backmsg.Cid = MSG_LIFE_TREE_SEND_INFO
	backmsg.MainLevel = self.San_LifeTree.MainLevel
	backmsg.Info = self.San_LifeTree.info

	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 主等级提升
func (self *ModLifeTree) MsgUpMainLevel(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return
	}

	config := GetCsvMgr().GetTreeLevelConfig(self.San_LifeTree.MainLevel)
	if nil == config {
		return
	}

	if err := self.player.HasObjectOkEasy(config.ItemID, config.ItemNum); err != nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("生命水晶不足"))
		return
	}

	oldLevel := self.San_LifeTree.MainLevel
	self.San_LifeTree.MainLevel++
	costItem := self.player.RemoveObjectSimple(config.ItemID, config.ItemNum, "生命树主等级", self.San_LifeTree.MainLevel, oldLevel, 0)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_LIFE_TREE_MAIN_LEVEL, self.San_LifeTree.MainLevel, oldLevel, 0, "创世神木升级", 0, 0, self.player)
	self.player.HandleTask(TASK_TYPE_LIFETREE_MAIN_LEVEL, self.San_LifeTree.MainLevel, 0, 0)
	self.UpdateHero()

	self.player.NoticeCenterBaseInfo()

	var backmsg S2C_LifeTreeUpMainLevel
	backmsg.Cid = MSG_LIFE_TREE_UP_MAIN_LEVEL
	backmsg.Items = costItem
	backmsg.Level = self.San_LifeTree.MainLevel
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 次等级提升
func (self *ModLifeTree) MsgUpTypeLevel(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return
	}

	var msg C2S_LifeTreeUpTypeLevel
	json.Unmarshal(body, &msg)

	mainConfig := GetCsvMgr().GetTreeLevelConfig(self.San_LifeTree.MainLevel)
	if nil == mainConfig {
		return
	}

	data := self.GetInfoData(msg.Type)
	if data == nil {
		return
	}

	if data.Level+1 > mainConfig.ProfessionalLevel {
		self.player.SendErrInfo("err", "达到最高等级")
		return
	}

	typeConfig := GetCsvMgr().GetTreeProfessionalConfig(msg.Type, data.Level)
	if nil == typeConfig {
		return
	}

	if err := self.player.HasObjectOkEasy(typeConfig.ItemID, typeConfig.ItemNum); err != nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("生命水晶不足"))
		return
	}

	oldLevel := data.Level
	data.Level++
	costItem := self.player.RemoveObjectSimple(typeConfig.ItemID, typeConfig.ItemNum, "神木职业升级", data.Level, data.Type, 0)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_LIFE_TREE_CAMP_LEVEL, data.Level, oldLevel, 0, "神木职业升级", 0, 0, self.player)

	self.player.HandleTask(TASK_TYPE_LIFETREE_OTHER_LEVEL, data.Level, data.Type, 0)

	self.UpdateHero()

	self.player.NoticeCenterBaseInfo()

	var backmsg S2C_LifeTreeUpTypeLevel
	backmsg.Cid = MSG_LIFE_TREE_UP_TYPE_LEVEL
	backmsg.Items = costItem
	backmsg.Level = data.Level
	backmsg.Type = msg.Type
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 次等级重置
func (self *ModLifeTree) MsgResetTypeLevel(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return
	}

	var msg C2S_LifeTreeResetTypeLevel
	json.Unmarshal(body, &msg)

	data := self.GetInfoData(msg.Type)
	if data == nil {
		return
	}

	if data.Level <= 0 {
		return
	}
	oldLevel := data.Level
	item := []PassItem{}
	err := self.player.HasObjectOkEasy(HERO_RESET_ITEM, 1)
	if err != nil {
		// 获得消耗配置
		config := GetCsvMgr().GetTariffConfig2(TARIFF_TYPE_LIFE_TREE_RESET)
		if config == nil {
			self.player.SendErrInfo("err", "生命树配置不存在")
			return
		}

		err2 := self.player.HasObjectOk(config.ItemIds, config.ItemNums)
		if err2 != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}

		item = self.player.RemoveObjectLst(config.ItemIds, config.ItemNums, "神木职业重置", 0, data.Type, 0)
	} else {
		item = self.player.RemoveObjectEasy(HERO_RESET_ITEM, 1, "神木职业重置", 0, data.Type, 0)
	}

	itemMap := make(map[int]*Item, 0)
	for i := data.Level - 1; i >= 0; i-- {
		typeConfig := GetCsvMgr().GetTreeProfessionalConfig(msg.Type, i)
		if nil == typeConfig {
			continue
		}
		AddItemMapHelper3(itemMap, typeConfig.ItemID, typeConfig.ItemNum)
	}

	additem := self.player.AddObjectItemMap(itemMap, "生命树重置", msg.Type, 0, 0)
	item = append(item, additem...)
	// 重置等级和属性
	data.Level = 0
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_LIFE_TREE_CAMP_LEVEL_REBORN, data.Level, oldLevel, 0, "神木职业重置", 0, 0, self.player)
	self.UpdateHero()

	self.player.NoticeCenterBaseInfo()

	var backmsg S2C_LifeTreeResetTypeLevel
	backmsg.Cid = MSG_LIFE_TREE_RESET_TYPE_LEVEL
	backmsg.Items = item
	backmsg.Level = data.Level
	backmsg.Type = data.Type
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 领取奖励
func (self *ModLifeTree) MsgGetAward(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return
	}

	itemMap := make(map[int]*Item, 0)
	for _, v := range self.San_LifeTree.award {
		AddItemMapHelper2(itemMap, v.Award)
	}

	item := self.player.AddObjectItemMap(itemMap, "成长涓流", 0, 0, 0)
	self.San_LifeTree.award = []*JS_HeroAward{}

	var backmsg S2C_LifeTreeGetAward
	backmsg.Cid = MSG_LIFE_TREE_GET_AWARD
	backmsg.Items = item
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	self.SendRedPointInfo()
}

// 发送累计奖励信息
func (self *ModLifeTree) MsgSendAwardInfo(body []byte) {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return
	}

	var backmsg S2C_LifeTreeSendAwardInfo
	backmsg.Cid = MSG_LIFE_TREE_SEND_AWARD_INFO
	backmsg.Award = self.San_LifeTree.award
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

func (self *ModLifeTree) SendRedPointInfo() {
	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_LIFE_TREE)
	if !flag {
		return
	}

	var backmsg S2C_LifeTreeSendRedpointInfo
	backmsg.Cid = MSG_LIFE_TREE_SEND_REDPOINT_INFO
	itemMap := make(map[int]*Item, 0)
	for _, v := range self.San_LifeTree.award {
		AddItemMapHelper2(itemMap, v.Award)
	}

	for _, v := range itemMap {
		backmsg.Items = append(backmsg.Items, PassItem{v.ItemId, v.ItemNum})
	}

	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

func (self *ModLifeTree) CheckTask() {
	self.player.HandleTask(TASK_TYPE_LIFETREE_MAIN_LEVEL, self.San_LifeTree.MainLevel, 0, 0)
	for _, v := range self.San_LifeTree.info {
		self.player.HandleTask(TASK_TYPE_LIFETREE_OTHER_LEVEL, v.Level, v.Type, 0)
	}
}
