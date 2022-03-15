package game

import (
	"encoding/json"
	"fmt"
)

const (
	DREAM_LAND_GOLD = 101
	DREAM_LAND_GEM  = 201
)

const (
	CS_GET_DREAMLAND_BASE_INFO = "get_dreamland_base_info"
	CS_GET_DREAMLAND_ITEM_INFO = "get_dreamland_item_info"
	CS_GET_DREAMLAND_START     = "get_dreamland_start"
	CS_GET_DREAMLAND_REFRESH   = "get_dreamland_refresh"
)

const (
	REFRESH_ITEM_TYPE_SPECIAL = 0
	REFRESH_ITEM_TYPE_GUIDE   = 4
	REFRESH_ITEM_TYPE_MAX     = 12
)

// 神格抽奖配置
type DreamLandItem struct {
	ID     int `json:"id"` //在配置中的id 非物品id
	ItemID int `json:"itemid"`
	Num    int `json:"num"`
}

type DreamLandLootItems struct {
	RefreshCount int `json:"refreshcount"` // 当天刷新次数
	LootCount    int `json:"lootcount"`    // 总抽奖次数 当数值过大时会归零
	FreeTimes    int `json:"freetimes`     // 免费次数
	//Class        int `json:"class`         // 抽奖类型

	//TypeTimes    int `json:"typetimes"`    // 隔多少次出一次必出
	//RefeshCost   int `json:"refreshcost"`  // 刷新花费id
	//LootCost     int `json:"lootcost"`     // 抽奖花费id

	Items []*DreamLandItem
}

type San_DreamLand struct {
	Uid       int64                       // 角色ID
	Info      string                      // 抽奖信息
	info      map[int]*DreamLandLootItems // 刷出的物品
	GuideType int

	DataUpdate
}

// 神格幻境抽奖
type ModDreamLand struct {
	player        *Player
	Sql_DreamLand San_DreamLand
}

func (self *ModDreamLand) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_dreamland` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_DreamLand, "san_dreamland", self.player.ID)

	if self.Sql_DreamLand.Uid <= 0 {
		self.Sql_DreamLand.Uid = self.player.ID
		self.Sql_DreamLand.GuideType = 0
		self.Sql_DreamLand.info = make(map[int]*DreamLandLootItems)
		for key, _ := range GetCsvMgr().DreamLandSpendMap {
			self.Sql_DreamLand.info[key] = &DreamLandLootItems{0, 0, 0, make([]*DreamLandItem, 0)}
		}
		self.AllRefresh()
		// 暂时特殊处理
		initTimes1, ok := GetCsvMgr().SimpleConfigMap[9]
		if ok && initTimes1 != nil {
			self.Sql_DreamLand.info[DREAM_LAND_GEM].FreeTimes = initTimes1.Num
		}

		initTimes2, ok := GetCsvMgr().SimpleConfigMap[10]
		if ok && initTimes2 != nil {
			self.Sql_DreamLand.info[DREAM_LAND_GOLD].FreeTimes = initTimes2.Num
		}

		self.Encode()
		InsertTable("san_dreamland", &self.Sql_DreamLand, 0, true)
		self.Sql_DreamLand.Init("san_dreamland", &self.Sql_DreamLand, true)
	} else {
		self.Decode()
		self.Sql_DreamLand.Init("san_dreamland", &self.Sql_DreamLand, true)
	}

	if self.Sql_DreamLand.info == nil {
		self.Sql_DreamLand.info = make(map[int]*DreamLandLootItems)
		for key, _ := range GetCsvMgr().DreamLandSpendMap {
			self.Sql_DreamLand.info[key] = &DreamLandLootItems{0, 0, 0, make([]*DreamLandItem, 0)}
		}
	}
}

func (self *ModDreamLand) OnGetOtherData() {}

func (self *ModDreamLand) Decode() {
	err := json.Unmarshal([]byte(self.Sql_DreamLand.Info), &self.Sql_DreamLand.info)
	if err != nil {
		LogError(err)
	}
}

func (self *ModDreamLand) Encode() {
	self.Sql_DreamLand.Info = HF_JtoA(&self.Sql_DreamLand.info)
}

// 存盘逻辑
func (self *ModDreamLand) OnSave(sql bool) {
	self.Encode()
	self.Sql_DreamLand.Update(sql)
}

// 消息处理
func (self *ModDreamLand) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case CS_GET_DREAMLAND_BASE_INFO:
		var msg C2S_DreamLandBaseInfo
		json.Unmarshal(body, &msg)
		self.SendBaseInfo()
		return true
	case CS_GET_DREAMLAND_ITEM_INFO: // 获得抽奖物品相关信息
		var msg C2S_DreamLandItemInfo
		json.Unmarshal(body, &msg)
		self.SendItemInfo(msg.Type)
		return true
	case CS_GET_DREAMLAND_START: //开始抽奖
		var msg C2S_DreamLandLoot
		json.Unmarshal(body, &msg)
		self.StartLoot(msg.Type, msg.Times)
		return true
	case CS_GET_DREAMLAND_REFRESH: // 刷新
		var msg C2S_DreamLandRefresh
		json.Unmarshal(body, &msg)
		self.StartRefresh(msg.Type)
		return true
	}

	return false
}

// 发送物品信息
func (self *ModDreamLand) SendBaseInfo() bool {
	var msg S2C_DreamLandBaseInfo
	msg.Cid = CS_GET_DREAMLAND_BASE_INFO

	msg.FreeTimes1 = self.Sql_DreamLand.info[DREAM_LAND_GOLD].FreeTimes
	msg.FreeTimes2 = self.Sql_DreamLand.info[DREAM_LAND_GEM].FreeTimes

	msg.RefreshCount1 = self.Sql_DreamLand.info[DREAM_LAND_GOLD].RefreshCount
	msg.RefreshCount1 = self.Sql_DreamLand.info[DREAM_LAND_GEM].RefreshCount

	msg.LuckyTimes1 = 0
	msg.TypeTimes1 = 0
	config1, ok := GetCsvMgr().DreamLandCostMap[DREAM_LAND_GOLD]
	if ok {
		if config1.TypeTimes > 0 {
			msg.LuckyTimes1 = self.Sql_DreamLand.info[DREAM_LAND_GOLD].LootCount
			msg.TypeTimes1 = config1.TypeTimes
		}
	}

	msg.LuckyTimes2 = 0
	msg.TypeTimes2 = 0
	config2, ok := GetCsvMgr().DreamLandCostMap[DREAM_LAND_GEM]
	if ok {
		if config2.TypeTimes > 0 {
			msg.LuckyTimes2 = self.Sql_DreamLand.info[DREAM_LAND_GEM].LootCount
			msg.TypeTimes2 = config2.TypeTimes
		}
	}

	msg.RefreshCount1 = self.Sql_DreamLand.info[DREAM_LAND_GOLD].RefreshCount
	msg.RefreshCount2 = self.Sql_DreamLand.info[DREAM_LAND_GEM].RefreshCount

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	return true
}

// 发送物品信息
func (self *ModDreamLand) SendItemInfo(Type int) bool {
	items, ok := self.Sql_DreamLand.info[Type]
	if !ok || items == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_WRONG_TYPE"))
		return false
	}

	config, ok := GetCsvMgr().DreamLandCostMap[Type]
	if !ok || config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_WRONG_TYPE"))
		return false
	}

	var msg S2C_DreamLandItemInfo
	msg.Cid = CS_GET_DREAMLAND_ITEM_INFO
	msg.Type = Type
	msg.RefreshCount = items.RefreshCount
	msg.FreeTimes = items.FreeTimes
	if config.TypeTimes > 0 {
		msg.LuckyTimes = items.LootCount
	} else {
		msg.LuckyTimes = 0
	}

	msg.TypeTimes = config.TypeTimes

	for _, v := range items.Items {
		msg.Items = append(msg.Items, PassItem{v.ItemID, v.Num})
	}

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	return true
}

// 开始刷新
func (self *ModDreamLand) StartRefresh(Type int) bool {
	items, ok := self.Sql_DreamLand.info[Type]
	if !ok || items == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_WRONG_TYPE"))
		return false
	}

	config, ok := GetCsvMgr().DreamLandCostMap[Type]
	if !ok || config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_WRONG_TYPE"))
		return false
	}

	costConfig := GetCsvMgr().GetTariffConfig(config.RefCost, self.Sql_DreamLand.info[Type].RefreshCount+1)

	if costConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_COST_WEONG"))
		return false
	}

	// 物品不足
	if err := self.player.HasObjectOk(costConfig.ItemIds, costConfig.ItemNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return false
	}

	// 移除物品
	outitem := self.player.RemoveObjectLst(costConfig.ItemIds, costConfig.ItemNums, "幻境刷新", 0, 0, 0)

	// 次数增加
	self.Sql_DreamLand.info[Type].RefreshCount++

	// 刷新
	if self.Refresh(Type) {

		var msg S2C_DreamLandRefresh
		msg.Cid = CS_GET_DREAMLAND_REFRESH
		msg.Type = Type
		msg.Items = outitem
		msg.RefreshCount = self.Sql_DreamLand.info[Type].RefreshCount
		for _, v := range items.Items {
			msg.NewItem = append(msg.NewItem, PassItem{v.ItemID, v.Num})
		}

		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

		return true
	}

	return false
}

// 开始抽奖
func (self *ModDreamLand) StartLoot(Type, Time int) bool {
	items, ok := self.Sql_DreamLand.info[Type]
	if !ok || items == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_WRONG_TYPE"))
		return false
	}

	config, ok := GetCsvMgr().DreamLandCostMap[Type]
	if !ok || config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_WRONG_TYPE"))
		return false
	}

	for _, obj := range items.Items {
		objConfig, ok := GetCsvMgr().DreamLandItemMap[obj.ID]
		if !ok || objConfig == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_WRONG_ITEM"))
			return false
		}
	}

	// 数量不足
	if len(items.Items) != REFRESH_ITEM_TYPE_MAX {
		return false
	}

	log := 0
	dec := ""
	if Type == DREAM_LAND_GOLD {
		if Time == 1 {
			//log = LOG_HERO_TALENT_LOOT_GOLD_ONE
			dec = "神格普通抽取"
		} else {
			log = LOG_HERO_TALENT_LOOT_GOLD_TEN
			dec = "神格普通十连"
		}

	} else {
		if Time == 1 {
			log = LOG_HERO_TALENT_LOOT_GEM_ONE
			dec = "神格高级抽取"
		} else {
			log = LOG_HERO_TALENT_LOOT_GEM_TEN
			dec = "神格高级十连"
		}
	}

	var outitems []PassItem

	// 有免费次数 且是单抽
	if self.Sql_DreamLand.info[Type].FreeTimes > 0 && Time == 1 {
		// 扣除免费次数
		self.Sql_DreamLand.info[Type].FreeTimes--

	} else { // 计算并扣除货币

		costConfig := GetCsvMgr().GetTariffConfig2(config.LootCost)

		if costConfig == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_COST_WEONG"))
			return false
		}

		if len(costConfig.ItemIds) != len(costConfig.ItemNums) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_COST_WEONG"))
			return false
		}

		var res = make(map[int]*Item)

		for i := 0; i < Time; i++ {
			AddItemMapHelper(res, costConfig.ItemIds, costConfig.ItemNums)
		}

		if err := self.player.hasItemMapOk(res); err != nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DREAM_LAND_NO_COST"))
			return false
		}

		outitems = self.player.RemoveObjectItemMap(res, dec, 0, 0, 0)
	}

	var initems []PassItem
	var indexs []int

	for i := 0; i < Time; i++ {
		// 唯一一次 引导状态 第一次普通单抽
		if self.Sql_DreamLand.GuideType == 0 &&
			Type == DREAM_LAND_GOLD &&
			Time == 1 &&
			self.Sql_DreamLand.info[Type].LootCount == 0 {

			finalConfig := items.Items[REFRESH_ITEM_TYPE_GUIDE]

			itemId, itemNum := self.player.AddObject(finalConfig.ItemID, finalConfig.Num, Type, 0, 0, dec)

			initems = append(initems, PassItem{itemId, itemNum})

			indexs = append(indexs, REFRESH_ITEM_TYPE_GUIDE)

			self.Sql_DreamLand.info[Type].LootCount = 1

			self.Sql_DreamLand.GuideType = 1

		} else if config.TypeTimes != 0 && // 必出情况
			self.Sql_DreamLand.info[Type].LootCount != 0 &&
			self.Sql_DreamLand.info[Type].LootCount >= config.TypeTimes &&
			items.Items[REFRESH_ITEM_TYPE_SPECIAL].Num > 0 {

			finalConfig := items.Items[REFRESH_ITEM_TYPE_SPECIAL]
			finalObjConfig := GetCsvMgr().DreamLandItemMap[finalConfig.ID]

			itemId, itemNum := self.player.AddObject(finalConfig.ItemID, finalConfig.Num, Type, 0, 0, dec)

			initems = append(initems, PassItem{itemId, itemNum})

			if finalObjConfig.Notice == 1 {
				GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_MOD_DREAM_LAND_NOTICE"),
					HF_GetColorByCamp(self.player.Sql_UserBase.Camp),
					CAMP_NAME[self.player.Sql_UserBase.Camp-1], self.player.Sql_UserBase.UName, GetCsvMgr().GetItemName(finalConfig.ItemID), finalConfig.Num), 0, 1)
			}

			if finalObjConfig.Only == 1 {
				self.Sql_DreamLand.info[Type].Items[REFRESH_ITEM_TYPE_SPECIAL].Num = 0
			}

			indexs = append(indexs, REFRESH_ITEM_TYPE_SPECIAL)

			self.Sql_DreamLand.info[Type].LootCount = 1

		} else {
			// 给个默认最后一个
			index := REFRESH_ITEM_TYPE_MAX - 1
			finalConfig := items.Items[index]

			// 随机出物品
			bSuccess, tempFinalConfig, tempIndex := self.RandLootItem(items.Items)

			if bSuccess {
				index = tempIndex
				finalConfig = tempFinalConfig
			}

			itemId, itemNum := self.player.AddObject(finalConfig.ItemID, finalConfig.Num, Type, 0, 0, dec)

			initems = append(initems, PassItem{itemId, itemNum})

			finalObjConfig := GetCsvMgr().DreamLandItemMap[finalConfig.ID]

			// 唯一物品
			if finalObjConfig.Only == 1 {
				finalConfig.Num = 0
			}

			if finalObjConfig.Notice == 1 {
				GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_MOD_DREAM_LAND_NOTICE"),
					HF_GetColorByCamp(self.player.Sql_UserBase.Camp),
					CAMP_NAME[self.player.Sql_UserBase.Camp-1], self.player.Sql_UserBase.UName, GetCsvMgr().GetItemName(finalConfig.ItemID), finalConfig.Num), 0, 1)
			}

			indexs = append(indexs, index)

			if self.Sql_DreamLand.info[Type].LootCount < config.TypeTimes {
				self.Sql_DreamLand.info[Type].LootCount++
			}
		}
	}

	var msg S2C_DreamLandLoot
	msg.Cid = CS_GET_DREAMLAND_START
	msg.Type = Type
	msg.FreeTime = items.FreeTimes
	if config.TypeTimes > 0 {
		msg.LuckyTimes = items.LootCount
	} else {
		msg.LuckyTimes = 0
	}

	//self.player.HandleTask(DreamLandLootTask, Type, Time, 0)
	self.player.HandleTask(DreamLandLootTenTask, Type, Time, 0)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, log, 0, 0, 0, dec, 0, 0, self.player)

	msg.TypeTimes = config.TypeTimes
	msg.Indexs = indexs
	msg.OutItems = outitems
	msg.InItems = initems

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	return true
}

// 凌晨五点 转盘刷新
func (self *ModDreamLand) AllRefresh() bool {
	for key, _ := range GetCsvMgr().DreamLandSpendMap {
		self.Refresh(key)
	}

	return true
}

// 具体刷新
func (self *ModDreamLand) Refresh(Type int) bool {
	// 随机刷新类型
	ok, Config := self.RandRefreshType(Type)

	if !ok || Config == nil {
		return false
	}

	// 随机刷出转盘物品
	ok, Item := self.RandRefreshItem(Config)
	if !ok || len(Item) <= 0 {
		return false
	}

	//self.Sql_DreamLand.info[Type].Class = Config.Class
	self.Sql_DreamLand.info[Type].Items = Item

	return true
}

// 刷新出配置
func (self *ModDreamLand) RandRefreshType(Type int) (bool, *DreamLandSpend) {
	// 取配置 看是否有该配置类型
	configs, ok := GetCsvMgr().DreamLandSpendMap[Type]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false, nil
	}

	// 总权重值
	nTotalChance := 0

	randConfig := make([]*DreamLandSpend, 0)

	// 筛选出符合需求的配置
	for index, _ := range configs {
		if self.player.Sql_UserBase.Vip >= configs[index].Minvip && self.player.Sql_UserBase.Vip <= configs[index].Maxvip {
			nTotalChance += configs[index].Chance
			randConfig = append(randConfig, configs[index])
		}
	}

	// 没有配置
	if nTotalChance <= 0 {
		return false, nil
	}

	// 随机配置
	nRandNum := HF_GetRandom(nTotalChance)

	total := 0

	finalConfig := randConfig[0]

	// 根据权重返回配置
	for _, v := range randConfig {
		total += v.Chance
		if nRandNum < total {
			finalConfig = v
			break
		}
	}

	return true, finalConfig
}

// 刷新出转盘物品配置
func (self *ModDreamLand) RandRefreshItem(config *DreamLandSpend) (bool, []*DreamLandItem) {
	if config == nil {
		return false, nil
	}

	// 配置出错
	if len(config.Group) != REFRESH_ITEM_TYPE_MAX {
		return false, nil
	}

	// 返回的物品列表
	var items []*DreamLandItem

	// 循环12件物品
	for i, _ := range config.Group {
		nGroup := 0
		//nNum := 0

		//// 特殊物品默认一个
		//if i == REFRESH_ITEM_TYPE_SPECIAL {
		//	nGroup = config.Group[i]
		//	nNum = 1
		//} else {
		//	nGroup = config.Group[i]
		//	nNum = config.Num[i]
		//}

		// 物品默认一个
		nGroup = config.Group[i]
		//nNum = 1

		// 根据掉落组 找到对应掉落表
		configs, ok := GetCsvMgr().DreamLandGroupMap[nGroup]
		if !ok {
			continue
		}

		// 获得该掉落表的总权重值
		nTotalChance := 0
		for _, v := range configs {
			if v.HaveHero > 0 {
				hero := self.player.GetModule("hero").(*ModHero).GetHero(v.HaveHero)
				if hero == nil {
					continue
				}
			}

			nTotalChance += v.RefreshWeight
		}

		// 没有物品
		if nTotalChance <= 0 {
			continue
		}

		//// 随多个物品
		//for j := 0; j < nNum; j++ {
		nRandNum := HF_GetRandom(nTotalChance)

		finalConf := DreamLandItem{configs[0].ID, configs[0].Item, configs[0].Num}

		total := 0
		for _, v := range configs {
			if v.HaveHero > 0 {
				hero := self.player.GetModule("hero").(*ModHero).GetHero(v.HaveHero)
				if hero == nil {
					continue
				}
			}

			total += v.RefreshWeight
			if nRandNum < total {
				finalConf = DreamLandItem{v.ID, v.Item, v.Num}

				//nTotalChance -= v.RefreshWeight                         // 移除已随的物品权重
				//configs = append(configs[:index], configs[index+1:]...) // 移除已抽取物品
				break
			}
		}

		items = append(items, &finalConf) // 增加物品
	}

	return true, items
}

// 刷新出转盘物品配置
func (self *ModDreamLand) RandLootItem(items []*DreamLandItem) (bool, *DreamLandItem, int) {
	nTotalChance := 0
	for _, obj := range items {
		objConfig := GetCsvMgr().DreamLandItemMap[obj.ID]

		if obj.Num > 0 {
			// 是否有英雄限制
			if objConfig.HaveHero > 0 {
				hero := self.player.GetModule("hero").(*ModHero).GetHero(objConfig.HaveHero)
				if hero == nil {
					LogError("出现错误 没有英雄的情况下 转盘里随到了英雄")
					continue
				}
				nTotalChance += objConfig.HaveChance
			} else {
				nTotalChance += objConfig.ExtractWeight
			}
		}
	}

	if nTotalChance <= 0 {
		return false, nil, -1
	}

	// 随机配置
	nRandNum := HF_GetRandom(nTotalChance)

	total := 0

	// 根据权重返回配置
	for k, v := range items {
		if v.Num <= 0 {
			continue
		}

		objConfig, _ := GetCsvMgr().DreamLandItemMap[v.ID]

		if objConfig.HaveHero > 0 {
			hero := self.player.GetModule("hero").(*ModHero).GetHero(objConfig.HaveHero)
			if hero == nil {
				LogError("出现错误 没有英雄的情况下 转盘里随到了英雄")
				continue
			}
			total += objConfig.HaveChance
		} else {
			total += objConfig.ExtractWeight
		}

		if nRandNum < total {
			return true, v, k
		}
	}

	return false, nil, -1
}

// 凌晨5点重置回调
func (self *ModDreamLand) OnRefresh() {
	self.AllRefresh()
	self.Sql_DreamLand.info[DREAM_LAND_GOLD].RefreshCount = 0
	self.Sql_DreamLand.info[DREAM_LAND_GEM].RefreshCount = 0

	// 暂时特殊处理
	initTimes1, ok := GetCsvMgr().SimpleConfigMap[9]
	if ok && initTimes1 != nil {
		self.Sql_DreamLand.info[DREAM_LAND_GEM].FreeTimes = initTimes1.Num
	}

	initTimes2, ok := GetCsvMgr().SimpleConfigMap[10]
	if ok && initTimes2 != nil {
		self.Sql_DreamLand.info[DREAM_LAND_GOLD].FreeTimes = initTimes2.Num
	}
}

func (self *ModDreamLand) GetStatisticsValue2040() (freeValue int, itemNum int) {

	needItemId := 81400701 //为了提高效率不查表写死，这个一般不会动
	relItemNum := self.player.GetObjectNum(needItemId)
	return self.Sql_DreamLand.info[DREAM_LAND_GOLD].FreeTimes, relItemNum
}

func (self *ModDreamLand) GetStatisticsValue2050() (itemNum int) {
	return self.player.GetObjectNum(81400702)
}
