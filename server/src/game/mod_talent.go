package game

type TalentInfo struct {
	Id int `json:"id"`
	Lv int `json:"lv"`
}

type TalentItem struct {
	Talents    []*TalentInfo      `json:"talents"`    // 天赋信息
	AttrMap    map[int]*Attribute `json:"attr"`       //! 总属性, 登录以及操作成功后重新计算
	AwakeStep  int                `json:"awake_step"` // 该英雄觉醒的层级
	MainTalent int                `json:"maintalent"` // 英雄主天赋等级
}

func NewTalentItem(id int) *TalentInfo {
	return &TalentInfo{Id: id, Lv: 0}
}

type Skill struct {
	SkillId int
	SkillLv int
}

// 从配置读取初始天赋
func (self *Hero) checkTalentItem() {
	/*
		if self.TalentItem == nil {
			self.TalentItem = &TalentItem{}
			self.TalentItem.AwakeStep = 0
			self.TalentItem.MainTalent = 0
			heroId := self.getHeroId()
			heroConfig := GetCsvMgr().GetHeroConfig(heroId)
			if heroConfig == nil {
				return
			}

			nLen := len(heroConfig.Point)
			if nLen <= 0 {
				return
			}

			for i := 0; i < nLen; i++ {
				self.TalentItem.Talents = append(self.TalentItem.Talents,
					NewTalentItem(heroConfig.Point[i])) //初始化给英雄设置天赋
			}
		}

		// 计算出技能
		self.cacTalentAtt()

	*/
}

func (self *Hero) cacTalentAtt() {
	/*
		//heroId := self.getHeroId()
		if self.TalentItem != nil {
			heroId := self.getHeroId()
			self.TalentItem.CheckHeroTalentAwakeStep(heroId, self.Uid)
			self.TalentItem.AttrMap = self.TalentItem.getAttr(heroId)
		}

	*/
}

func (self *Hero) GetTalentAttr() map[int]*Attribute {
	/*
		if self.TalentItem != nil {
			return self.TalentItem.AttrMap
		}

	*/
	return nil
}

// 天赋激活或者升级
func (self *ModHero) upgradeTalent(pMsg *C2S_UpgradeTalent) {
	/*
		//用户等级不足
		flag, _ := GetCsvMgr().IsLevelOpen(self.player.Sql_UserBase.Level, OPEN_LEVEL_TALENT_RESET)
		if !flag {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_PLAYER_LOW_LEVEL"))
			return
		}

		heroId := pMsg.HeroId
		pHero := self.GetHero(heroId)
		if pHero == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_HEROES_DO_NOT_EXIST"))
			return
		}

		heroConfig := GetCsvMgr().GetHeroConfig(heroId)
		if heroConfig == nil {
			return
		}

		nLen := len(heroConfig.Point)
		if nLen <= 0 {
			return
		}

		index := pMsg.Index
		if index < 1 || index > nLen {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_GIFTED_SUBSCRIPTS_DO_NOT_EXIST"))
			return
		}

		pos := index - 1
		// 检查资质是否已经满星
		talentItem := pHero.TalentItem
		if talentItem == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_DATA_EXCEPTIONS_FOR_GIFTED_MODULES"))
			return
		}

		if len(talentItem.Talents) <= 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_NATURAL_DATA_ABNORMALITY"))
			return
		}

		pTalent := talentItem.Talents[pos]
		maxLevel, ok := GetCsvMgr().MaxTalentLv[pTalent.Id]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_SERVER_MAXIMUM_TALENT_DOES_NOT"))
			return
		}
		// 获取英雄的最大资质星级
		config := GetCsvMgr().GetTalentConfig(pTalent.Id, pTalent.Lv)
		if config == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_TALENT_ACTIVATION_CONSUMPTION_CONFIGURATION_DOES"))
			return
		}

		// 技能等级不足
		if self.player.Sql_UserBase.Level < config.NeedLevel {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_PLAYER_LOW_LEVEL"))
			return
		}

		if pTalent.Lv >= maxLevel {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_TALENT_HAS_REACHED_ITS_MAXIMUM"))
			return
		}

		// 检查道具消耗
		if err := self.player.HasObjectOk(config.CostItems, config.Costnums); err != nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_INSUFFICIENT_MATERIAL"))
			return
		}

		// 扣除道具
		items := self.player.RemoveObjectLst(config.CostItems, config.Costnums, "神格升级", pTalent.Lv+1, pHero.getHeroId(), 0)
		pTalent.Lv += 1
		// 重新计算技能和属性
		pHero.cacTalentAtt()
		self.player.countHeroFight(pHero, 0)
		self.player.GetModule("team").(*ModTeam).CacTalents()
		self.player.HandleTask(TalentTask, 0, 0, 0)
		self.player.HandleTask(HaveDinivityTask, 0, 0, 0)
		self.player.HandleTask(AllDinivityTask, 0, 0, 0)
		var msg S2C_UpgradeTalent
		msg.Cid = upgradeTalentAction
		msg.HeroId = pHero.getHeroId()
		msg.Items = items
		msg.TalentItem = pTalent
		msg.Step = pHero.TalentItem.AwakeStep
		msg.Attr = pHero.GetTalentAttr()
		msg.MainTelent = pHero.TalentItem.MainTalent
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_TALENT, pTalent.Lv, msg.HeroId, 0, "神格升级", 0, 0, self.player)


	*/
}

// 获取天赋类型和属性
/*
func (self *TalentItem) getAttr(heroid int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	// 资质属性加成
	var skills []*Skill
	for _, pTalent := range self.Talents {
		config := GetCsvMgr().GetTalentConfig(pTalent.Id, pTalent.Lv)
		if config == nil {
			continue
		}

		if pTalent.Lv <= 0 {
			continue
		}
		AddAttrDirect(attMap, config.AttTypes, config.AttValues)
		skills = append(skills, &Skill{config.SkillId, 1}) // 暂时技能没有等级 特殊处理为1
		skills = append(skills, &Skill{config.SkillId2, 1})
	}

	// 技能加成
	if len(skills) > 0 {
		skillAtr := GetSkillAttr(skills)
		AddAttrMapHelper(attMap, skillAtr)
	}

	// 获得英雄属性
	heroConfig := GetCsvMgr().GetHeroConfig(heroid)
	if heroConfig != nil {
		configs, ok := GetCsvMgr().TalentAwakeMap[heroConfig.DivinityAwaken] // 根据配置的觉醒模板选出第几组的觉醒配置
		if ok {
			nStep := self.AwakeStep // 获得当前英雄是第几层觉醒
			if nStep >= 0 && nStep < len(configs) {
				AddAttrDirect(attMap, configs[nStep].AttTypes, configs[nStep].AttValues)
			}
		}
	}

	return attMap
}
 */

// 检查该英雄当前的觉醒层级
func (self *TalentItem) CheckHeroTalentAwakeStep(heroid int, uid int64) {
	// 获得英雄属性
	heroConfig := GetCsvMgr().GetHeroConfig(heroid)
	player := GetPlayerMgr().GetPlayer(uid, false)

	if heroConfig == nil || player == nil {
		return
	}

	// 根据配置的觉醒模板选出第几组的觉醒配置
	configs, ok := GetCsvMgr().TalentAwakeMap[heroConfig.DivinityAwaken]
	if !ok {
		return
	}

	// 设置主天赋等级属性
	if len(self.Talents) > 0 {
		self.MainTalent = self.Talents[0].Lv
	} else {
		self.MainTalent = 0
	}

	// 判断当前是第几层觉醒
	nStep := 0
	for _, config := range configs {
		nLevelLimit := config.LevelLimit
		bIsHave := true
		for _, pTalent := range self.Talents {
			if nLevelLimit > pTalent.Lv {
				bIsHave = false
			}
		}

		if bIsHave == true {
			nStep = config.Step
		} else {
			break
		}
	}

	// 设置当前英雄是第几层觉醒

	if self.AwakeStep != nStep {
		self.AwakeStep = nStep
		//GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_HERO_TALENT_AWAKE, nStep, 0, 0, "神格觉醒", 0, 0, player)
	}
}

// 天赋重置
func (self *ModHero) TalentReset(pMsg *C2S_ResetTalent) {
	/*
		//用户等级不足
		flag, _ := GetCsvMgr().IsLevelOpen(self.player.Sql_UserBase.Level, OPEN_LEVEL_TALENT_RESET)
		if !flag {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REBORN_LV_NOT_ENOUGH"))
			return
		}

		pHero := self.GetHero(pMsg.HeroId)
		if pHero == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TIGER_HEROES_DO_NOT_EXIST"))
			return
		}

		costConfig := GetCsvMgr().GetTariffConfig2(TariffTalentReset)

		// 扣除道具
		items := self.player.RemoveObjectLst(costConfig.ItemIds, costConfig.ItemNums, "神格重置", pHero.getHeroId(), 0, 0)

		var additem []int
		var addnum []int

		for i := 0; i < len(pHero.TalentItem.Talents); i++ {
			nLevel := pHero.TalentItem.Talents[i].Lv
			nID := pHero.TalentItem.Talents[i].Id

			pHero.TalentItem.Talents[i].Lv = 0

			talentconfig := GetCsvMgr().GetTalentConfig(nID, nLevel)
			if talentconfig == nil {
				continue
			}

			if len(talentconfig.ReturnItems) <= 0 || len(talentconfig.ReturnItems) != len(talentconfig.Returnnums) {
				continue
			}

			for _, v := range talentconfig.ReturnItems {
				additem = append(additem, v)
			}

			for _, v := range talentconfig.Returnnums {
				addnum = append(addnum, v)
			}
		}

		outitem := self.player.AddObjectLst(additem, addnum, "神格重置", pHero.getHeroId(), 0, 0)

		// 重新计算技能和属性
		pHero.cacTalentAtt()
		self.player.countHeroFight(pHero, 0)
		self.player.GetModule("team").(*ModTeam).CacTalents()
		self.player.HandleTask(TalentTask, 0, 0, 0)
		self.player.HandleTask(HaveDinivityTask, 0, 0, 0)
		self.player.HandleTask(AllDinivityTask, 0, 0, 0)

		var msg S2C_ResetTalent
		msg.Cid = talentReset
		msg.HeroId = pHero.getHeroId()
		msg.Items = items
		msg.Attr = pHero.GetTalentAttr()
		msg.Step = pHero.TalentItem.AwakeStep
		msg.OutItems = outitem
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_TALENT_RESET, msg.HeroId, 0, 0, "神格重置", 0, 0, self.player)

	*/
}
