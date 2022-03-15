package game

type FateInfo struct {
	Id     int `json:"id"`     //! 缘分Id
	Status int `json:"status"` //! 状态, 是否激活, 0 未激活, 1激活
}

// 缘分
type FateItem struct {
	Fates   []*FateInfo        `json:"fate"` //! 缘分信息
	AttrMap map[int]*Attribute `json:"attr"` //! 总属性, 登录以及操作成功后重新计算
}

func NewFateItem(id int) *FateInfo {
	return &FateInfo{Id: id, Status: 0}
}

// 从配置读取初始缘分
func (self *Hero) checkFateItem() {
	/*
		if self.FateItem == nil {
			self.FateItem = &FateItem{}
			heroId := self.getHeroId()
			fates, ok := GetCsvMgr().FateMap[heroId]
			if !ok {
				return
			}

			for i := 0; i < len(fates); i++ {
				self.FateItem.Fates = append(self.FateItem.Fates,
					NewFateItem(fates[i].FateId))
			}
		}

		// 计算出技能
		self.cacFateAtt()

	*/
}

func (self *Hero) cacFateAtt() {
	/*
		if self.FateItem != nil {
			heroId := self.getHeroId()
			self.FateItem.AttrMap = self.FateItem.getAttr(heroId)
		}

	*/
}

func (self *Hero) GetFateAttr() map[int]*Attribute {
	/*
		if self.FateItem != nil {
			return self.FateItem.AttrMap
		}

	*/
	return nil
}

// 获取缘分类型和属性
func (self *FateItem) getAttr(heroId int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	for _, v := range self.Fates {
		if v.Status == 0 {
			continue
		}
		config := GetCsvMgr().getFateConfig(heroId, v.Id)
		if config == nil {
			continue
		}

		AddAttrDirect(attMap, config.AttType, config.AttValue)
	}
	return attMap
}

// 缘分
// 只判断fateType等于0的缘分
func (self *ModHero) activateFate(pMsg *C2S_ActivateFate) {
	/*
		pHero := self.GetHero(pMsg.HeroId)
		if pHero == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TIGER_HEROES_DO_NOT_EXIST"))
			return
		}

		heroId := pHero.getHeroId()
		index := pMsg.Index
		if index < 1 || index > 4 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FATE_THE_PREDESTINATION_SUBSCRIPT_DOES_NOT"))
			return
		}

		pos := index - 1
		// 检查资质是否已经满星
		fateItem := pHero.FateItem
		if fateItem == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FATE_DATA_EXCEPTION_OF_FATE_MODULE"))
			return
		}

		if len(fateItem.Fates) <= 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FATE_MARGIN_DATA_ANOMALY"))
			return
		}

		pFate := fateItem.Fates[pos]
		if pFate.Status == 1 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FATE_FATE_HAS_BEEN_ACTIVATED"))
			return
		}

		config := GetCsvMgr().getFateConfig(heroId, pFate.Id)
		if config == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FATE_MARGINAL_ALLOCATION_DOES_NOT_EXIST"))
			return
		}

		heroIds := config.Heroes
		if config.FateType != 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FATE_FATE_ACTIVATION_TYPE_ERROR"))
			return
		}

		// 英雄是否存在
		for _, id := range heroIds {
			if id == 0 {
				continue
			}

			has := self.GetHero(id)
			if has == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FATE_HEROES_DO_NOT_EXIST_AND"))
				return
			}
		}

		// 可以激活
		pFate.Status = 1

		changes := make(map[int]*Attribute)
		AddAttrDirect(changes, config.AttType, config.AttValue)

		// 重新计算技能和属性
		pHero.cacFateAtt()
		self.player.countHeroIdFight(pHero.HeroId, 0)
		var msg S2C_ActivateFate
		msg.Cid = activatefate
		msg.HeroId = pHero.getHeroId()
		msg.FateInfo = pFate
		msg.Attr = pHero.GetFateAttr()
		msg.Chanegs = changes
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FATE, pFate.Id, heroId, 0, "英雄宿命", 0, 0, self.player)

	*/
}
