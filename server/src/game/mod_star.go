package game

const (
	maxStarSlots = 9 // 星格孔个数
	STAR_INIT    = 1
)

type StarItem struct {
	UpStar      int                `json:"upstar"`      //! 升星后的star
	AttrMap     map[int]*Attribute `json:"attr"`        //! 总属性, 登录以及操作成功后重新计算
	Skills      []int              `json:"skills"`      //! 升星以及登录时进行计算
	HeroBreakId int                `json:"herobreakid"` // 英雄突破ID
}

// 从配置读取初始星级
func (self *Hero) checkStarItem(star int) {
	if self.StarItem == nil {
		self.StarItem = &StarItem{}

		//初始标准
		if star <= 1 {
			//获取初始星级
			self.StarItem.UpStar = GetCsvMgr().GetHeroInitLv(self.HeroId)
		} else {
			self.StarItem.UpStar = star
		}
		//看看配置是否对
		configHero := GetCsvMgr().GetHeroMapConfig(self.HeroId, self.StarItem.UpStar)
		if configHero == nil {
			self.StarItem.UpStar = STAR_INIT
		}
		// 计算出技能
		self.cacStarAtt()
	}
}

func (self *Hero) cacStarAtt() {
	if self.StarItem != nil {
		//计算技能
		self.StarItem.Skills = make([]int, 0)
		skillBreakConfig := GetCsvMgr().HeroBreakConfigMap[self.HeroId][self.StarItem.HeroBreakId]
		if skillBreakConfig == nil {
			return
		}

		for i := 0; i < len(skillBreakConfig.Skill); i++ {
			if skillBreakConfig.Skill[i] > 0 {
				self.StarItem.Skills = append(self.StarItem.Skills, skillBreakConfig.Skill[i])
			}
		}

		//计算属性
		self.StarItem.AttrMap = make(map[int]*Attribute)
		//突破属性
		AddAttrHelperForTimes(self.StarItem.AttrMap, skillBreakConfig.BaseTypes, skillBreakConfig.BaseValues, 1)
		//星级属性
		configStar := GetCsvMgr().GetHeroMapConfig(self.HeroId, self.StarItem.UpStar)
		if configStar == nil {
			return
		}

		AddAttrHelperForTimes(self.StarItem.AttrMap, configStar.BaseTypes, configStar.BaseValues, 1)
		//成长率
		configGrowth := GetCsvMgr().HeroGrowthConfigMap[self.HeroLv]
		if configGrowth != nil {
			AddAttrHelperForGrowth(self.StarItem.AttrMap, configStar.GrowthTypes, configStar.GrowthValues, configGrowth.GrowthType, configGrowth.GrowthValue, configStar.QuaType, configStar.QuaValue)
		}
	}

	//计算技能属性
	skillAttr := GetSkillAttr(self.StarItem.Skills)
	if len(skillAttr) > 0 {
		self.addAttEx(skillAttr, self.StarItem.AttrMap)
		//这个地方需要把万分比属性计算掉
		attrNew := make(map[int]*Attribute)
		for _, v := range self.StarItem.AttrMap {
			if v.AttType > AttrDisExt && v.AttType <= AttrEnd+AttrDisExt {
				attr, ok := self.StarItem.AttrMap[v.AttType-AttrDisExt]
				if ok {
					attr.AttValue = attr.AttValue * (1.0 + v.AttValue/10000.0)
				}
			}
		}
		for _, v := range self.StarItem.AttrMap {
			if v.AttType > AttrDisExt && v.AttType <= AttrEnd+AttrDisExt {
				continue
			}
			attrNew[v.AttType] = v
		}
		self.StarItem.AttrMap = attrNew
	}
}

func (self *Hero) GetStarAttr() map[int]*Attribute {
	if self.StarItem != nil {
		return self.StarItem.AttrMap
	}
	return nil
}

//
func (self *Hero) LvUp(lv int) {
	self.HeroLv += lv
	if self.HeroLv <= 0 {
		self.HeroLv = 1
	}

	config := GetCsvMgr().HeroBreakConfigMap[self.HeroId]
	if config == nil {
		return
	}
	self.StarItem.HeroBreakId = 0
	//计算突破等级
	for _, v := range config {
		if self.HeroLv >= v.Break && self.StarItem.HeroBreakId < v.Id {
			self.StarItem.HeroBreakId = v.Id
		}
	}

	GetOfflineInfoMgr().SetHeroMaxLevel(self.Uid, self.HeroLv)

	player := GetPlayerMgr().GetPlayer(self.Uid, false)
	if player != nil {
		player.HandleTask(TASK_TYPE_HAVE_HERO, 0, 0, 0)
	}
}

func (self *NewHero) LvUp(lv int) {
	self.HeroLv += lv

	config := GetCsvMgr().HeroBreakConfigMap[self.HeroId]
	if config == nil {
		return
	}

	//计算突破等级
	if self.StarItem.HeroBreakId == 0 {
		for _, v := range config {
			if self.HeroLv >= v.Break && self.StarItem.HeroBreakId < v.Id {
				self.StarItem.HeroBreakId = v.Id
			}
		}
	} else {
		nextConfig, ok := config[self.StarItem.HeroBreakId+1]
		if !ok {
			return
		}
		if self.HeroLv >= nextConfig.Break {
			self.StarItem.HeroBreakId = nextConfig.Id
		}
	}
	self.cacStarAtt()
}

// 一键升星
/*
func (self *ModHero) upStarAuto(pMsg *C2S_UpStarAuto) {
	pHero := self.GetHero(pMsg.HeroId)
	if pHero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TIGER_HEROES_DO_NOT_EXIST"))
		return
	}

	heroId := pHero.getHeroId()
	starItem := pHero.StarItem
	if starItem == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_ASCENDING_STAR_MODULE_DATA_ABNORMALITY"))
		return
	}

	if len(starItem.Slots) != maxStarSlots {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_ASCENDING_STAR_QUALIFICATION_DATA_ABNORMALITY"))
		return
	}

	// 检查材料是否充足
	itemMap := make(map[int]*Item)
	heroConfig, ok := GetCsvMgr().HeroStarMap[heroId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_CURRENT_HERO_ASCENSION_CONFIGURATION_DOES"))
		return
	}

	maxStar := len(heroConfig)

	// 资质激活消耗
	maxNum := 0
	count := 0
	for index, slotStar := range starItem.Slots {
		// 获取英雄的最大资质星级
		if slotStar >= maxStar {
			maxNum += 1
			continue
		}

		if slotStar >= starItem.UpStar+1 {
			continue
		}

		starConfig, ok := heroConfig[slotStar+1]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_CURRENT_ASCENDING_CONFIGURATION_DOES_NOT"))
			return
		}

		slotItemIds, slotItemNums := starConfig.getCost(index + 1)
		if len(slotItemIds) != len(slotItemNums) {
			continue
		}

		for costIndex := range slotItemIds {
			itemId := slotItemIds[costIndex]
			itemNum := slotItemNums[costIndex]
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

		//检查是否够
		isHas := true
		for itemId, item := range itemMap {
			if itemId == 0 || item.ItemNum == 0 {
				continue
			}
			if self.player.GetObjectNum(itemId) < item.ItemNum {
				isHas = false
				break
			}
		}

		if !isHas {
			for costIndex := range slotItemIds {
				itemId := slotItemIds[costIndex]
				itemNum := slotItemNums[costIndex]
				if itemId == 0 || itemNum == 0 {
					continue
				}
				itemMap[itemId].ItemNum -= itemNum
			}
			break
		} else {
			count++
		}
	}

	// 一键升级操作
	for index := range starItem.Slots {
		if count <= 0 {
			break
		}
		slotStar := starItem.Slots[index]
		if slotStar >= maxStar {
			continue
		}

		if slotStar >= starItem.UpStar+1 {
			continue
		}

		starItem.Slots[index] += 1
		count--
	}

	// 扣除道具
	items := self.player.RemoveObjectItemMap(itemMap, "英雄升星", starItem.UpStar, 0, 0)

	// 重新计算技能和属性
	pHero.cacStarAtt()
	self.player.countHeroFight(pHero, ReasonUpStar)

	self.player.HandleTask(HeroUpStarNumTask, 0, 0, 0)
	self.player.HandleTask(HeroUpStarTimesTask, 1, 0, 0)
	self.player.HandleTask(7, 0, 0, 0)

	self.GetAllHeroStars()
	//GetTopHeroStarsMgr().UpdateRank(self.GetAllHeroStars(), self.player)

	var msg S2C_UpStarAuto
	msg.Cid = activateStarAction
	msg.HeroId = pHero.getHeroId()
	msg.Items = items
	msg.StarItem = starItem
	msg.Attr = pHero.GetStarAttr()
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_UP_STAR, starItem.UpStar, heroId, 0, "英雄升星", 0, 0, self.player)

}


*/
func (self *NewHero) cacStarAtt() {
	if self.StarItem != nil {
		//计算技能
		self.StarItem.Skills = make([]int, 0)
		skillBreakConfig := GetCsvMgr().HeroBreakConfigMap[self.HeroId][self.StarItem.HeroBreakId]
		if skillBreakConfig == nil {
			return
		}

		for i := 0; i < len(skillBreakConfig.Skill); i++ {
			if skillBreakConfig.Skill[i] > 0 {
				self.StarItem.Skills = append(self.StarItem.Skills, skillBreakConfig.Skill[i])
			}
		}

		//计算属性
		self.StarItem.AttrMap = make(map[int]*Attribute)
		//突破属性
		AddAttrHelperForTimes(self.StarItem.AttrMap, skillBreakConfig.BaseTypes, skillBreakConfig.BaseValues, 1)
		//星级属性
		configStar := GetCsvMgr().GetHeroMapConfig(self.HeroId, self.StarItem.UpStar)
		if configStar == nil {
			return
		}

		AddAttrHelperForTimes(self.StarItem.AttrMap, configStar.BaseTypes, configStar.BaseValues, 1)
		//成长率
		configGrowth := GetCsvMgr().HeroGrowthConfigMap[self.HeroLv]
		if configGrowth != nil {
			AddAttrHelperForGrowth(self.StarItem.AttrMap, configStar.GrowthTypes, configStar.GrowthValues, configGrowth.GrowthType, configGrowth.GrowthValue, configStar.QuaType, configStar.QuaValue)
		}
	}

	//计算技能属性
	skillAttr := GetSkillAttr(self.StarItem.Skills)
	if len(skillAttr) > 0 {
		self.addAttEx(skillAttr, self.StarItem.AttrMap)
		//这个地方需要把万分比属性计算掉
		attrNew := make(map[int]*Attribute)
		for _, v := range self.StarItem.AttrMap {
			if v.AttType > AttrDisExt && v.AttType <= AttrEnd+AttrDisExt {
				attr, ok := self.StarItem.AttrMap[v.AttType-AttrDisExt]
				if ok {
					attr.AttValue = attr.AttValue * (1.0 + v.AttValue/10000.0)
				}
			}
		}
		for _, v := range self.StarItem.AttrMap {
			if v.AttType > AttrDisExt && v.AttType <= AttrEnd+AttrDisExt {
				continue
			}
			attrNew[v.AttType] = v
		}
		self.StarItem.AttrMap = attrNew
	}
}

func (self *NewHero) checkStarItem(star int) {
	if self.StarItem == nil {
		self.StarItem = &StarItem{}

		//初始标准
		if star <= 1 {
			//获取初始星级
			self.StarItem.UpStar = GetCsvMgr().GetHeroInitLv(self.HeroId)
		} else {
			self.StarItem.UpStar = star
		}
		//看看配置是否对
		configHero := GetCsvMgr().GetHeroMapConfig(self.HeroId, self.StarItem.UpStar)
		if configHero == nil {
			self.StarItem.UpStar = STAR_INIT
		}
		// 计算出技能
		self.cacStarAtt()
	}
}
