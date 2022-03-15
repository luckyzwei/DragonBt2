package game

import "encoding/json"

const (
	MSG_HERO_SET_VOID_HERO_RESONANCE    = "msg_hero_set_void_hero_resonance"
	MSG_HERO_CANCEL_VOID_HERO_RESONANCE = "msg_hero_cancel_void_hero_resonance"
	MSG_HERO_UPDATE_VOID_HERO_RESONANCE = "msg_hero_update_void_hero_resonance"
)

func (self *ModHero) CheckTeamVoidHero(heroKeyId int, voidHeroKeyId int) {
	info := self.player.GetModule("team").(*ModTeam).Sql_Team.info
	// 共鸣前期处理 将产生了共鸣的两个英雄下阵
	for i, team := range info {
		teamtype := i + 1
		switch teamtype {
		case TEAMTYPE_ARENA_SPECIAL_1, TEAMTYPE_ARENA_SPECIAL_4:
			{
				var msg C2S_AddTeamUIPos
				msg.TeamType = teamtype

				isHero := false
				isVoidHero := false
				for teampos := i; teampos <= i+2; teampos++ {
					for _, heroid := range info[teampos].FightPos {
						if heroid == 0 {
							continue
						}

						hero := self.GetHero(heroid)
						if hero == nil {
							continue
						}

						if heroid == heroKeyId {
							isHero = true
						}

						if heroid == voidHeroKeyId {
							isVoidHero = true
						}
					}
				}

				if isHero && isVoidHero {
					for teampos := i; teampos <= i+2; teampos++ {
						for _, heroid := range info[teampos].FightPos {
							if heroid == heroKeyId || heroid == voidHeroKeyId {
								msg.FightPos = append(msg.FightPos, 0)
							} else {
								msg.FightPos = append(msg.FightPos, heroid)
							}

						}
					}
					smsg, _ := json.Marshal(&msg)
					self.player.GetModule("team").(*ModTeam).AddArenaSpcialTeamUIPos(smsg)
				}
			}
		case TEAMTYPE_ARENA_SPECIAL_2, TEAMTYPE_ARENA_SPECIAL_3, TEAMTYPE_ARENA_SPECIAL_5, TEAMTYPE_ARENA_SPECIAL_6:
			{
				continue
			}
		default:
			{
				heroIndex := -1
				voidheroIndex := -1
				for index, heroid := range team.FightPos {
					if heroid == 0 {
						continue
					}

					hero := self.GetHero(heroid)
					if hero == nil {
						continue
					}

					if heroid == heroKeyId {
						heroIndex = index
					}

					if heroid == voidHeroKeyId {
						voidheroIndex = index
					}
				}

				if heroIndex >= 0 && voidheroIndex >= 0 {
					temp := team.FightPos
					temp[heroIndex] = 0
					temp[voidheroIndex] = 0

					var msg C2S_AddTeamUIPos
					msg.TeamType = teamtype
					for _, hero := range temp {
						msg.FightPos = append(msg.FightPos, hero)
					}
					smsg, _ := json.Marshal(&msg)
					self.player.GetModule("team").(*ModTeam).addUIPos(smsg)
				}
			}
		}

		// 后期处理 当下阵后有的队伍需要自动上阵
		switch teamtype {
		case TEAMTYPE_ARENA_2:
			{
				self.player.GetModule("arena").(*ModArena).CheckTeam()
			}
		case TEAMTYPE_ARENA_SPECIAL_4:
			{
				self.player.GetModule("arenaspecial").(*ModArenaSpecial).CheckTeam()
			}
		}
	}
}

// 设置虚空英雄共鸣
func (self *ModHero) MsgSetVoidHeroResonance(body []byte) {
	var msg C2S_VoidHeroResonanceSet
	json.Unmarshal(body, &msg)

	// 获得普通英雄
	hero := self.GetHero(msg.HeroKeyId)
	if hero == nil {
		return
	}

	// 获得共鸣英雄
	voidhero := self.GetHero(msg.VoidHeroKeyId)
	if voidhero == nil {
		return
	}
	// 是共鸣英雄 或者已经参加了共鸣则返回
	if hero.VoidHero != 0 || hero.Resonance != 0 {
		return
	}
	// 不是虚空英雄 或者已经参加了共鸣则返回
	if voidhero.VoidHero == 0 || voidhero.Resonance != 0 {
		return
	}

	if voidhero.UseType[HERO_USE_TYPE_CRYSTAL_RESONANCE] == 1 {
		modcrystal := self.player.GetModule("crystal").(*ModResonanceCrystal)
		if modcrystal != nil {
			index := modcrystal.GetResonanceIndex(voidhero.HeroKeyId)
			if index >= 0 {
				// 取消共鸣英雄
				modcrystal.CancelResonanceHeros(index, voidhero.HeroKeyId, false)
			}
		}
	}

	// 如果有等级 则设置回退
	config := GetCsvMgr().HeroExpConfigMap[voidhero.HeroLv]
	getItem := self.player.AddObjectLst(config.ResetItems, config.ResetNums, "虚空英雄共鸣回退", 0, 0, 0)

	// 设置星级和等级
	voidhero.StarItem.UpStar = hero.StarItem.UpStar
	voidhero.LvUp(hero.HeroLv - voidhero.HeroLv)
	self.player.countHeroFight(voidhero, ReasonHeroVoidSetResonance)

	// 互相设置共鸣绑定
	voidhero.Resonance = hero.HeroKeyId
	hero.Resonance = voidhero.HeroKeyId
	voidhero.CheckStageTalent()

	self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(voidhero.HeroKeyId)
	self.CheckTeamVoidHero(msg.HeroKeyId, msg.VoidHeroKeyId)

	var backmsg S2C_VoidHeroResonanceSet
	backmsg.Cid = MSG_HERO_SET_VOID_HERO_RESONANCE
	backmsg.HeroKeyId = msg.HeroKeyId
	backmsg.VoidHeroKeyId = msg.VoidHeroKeyId
	backmsg.GetItem = getItem
	backmsg.VoidHero = self.GetHero(msg.VoidHeroKeyId)
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

}

// 取消虚空英雄共鸣
func (self *ModHero) MsgCancelVoidHeroResonance(body []byte) {
	var msg C2S_VoidHeroResonanceCancel
	json.Unmarshal(body, &msg)

	configCost := GetCsvMgr().GetTariffConfig2(TARIFF_TYPE_VOID_HERO_CANCEL)
	if configCost == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	//检查消耗够不够
	if err := self.player.HasObjectOk(configCost.ItemIds, configCost.ItemNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	// 扣除物品
	costItem := self.player.RemoveObjectLst(configCost.ItemIds, configCost.ItemNums, "虚空英雄取消共鸣", msg.VoidHeroKeyId, 0, 0)

	if !self.CancelVoidHeroResonance(msg.VoidHeroKeyId) {
		return
	}

	var backmsg S2C_VoidHeroResonanceCancel
	backmsg.Cid = MSG_HERO_CANCEL_VOID_HERO_RESONANCE
	backmsg.VoidHeroKeyId = msg.VoidHeroKeyId
	backmsg.CostItem = costItem
	backmsg.VoidHero = self.GetHero(msg.VoidHeroKeyId)
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 取消虚空英雄共鸣
func (self *ModHero) CancelVoidHeroResonance(herokey int) bool {
	// 获得共鸣英雄
	voidhero := self.GetHero(herokey)
	if voidhero == nil {
		return false
	}

	// 不是虚空英雄 或者没有参加了共鸣则返回
	if voidhero.VoidHero == 0 || voidhero.Resonance == 0 {
		return false
	}

	// 获得普通英雄
	hero := self.GetHero(voidhero.Resonance)
	// 找不到 可能被分解
	if hero != nil {
		// 是虚空英雄 或者没有参与共鸣
		if hero.VoidHero != 0 || hero.Resonance == 0 {
			return false
		}

		if hero.Resonance != herokey {
			return false
		}

		hero.Resonance = 0
	}

	self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(voidhero.HeroKeyId)

	// 设置回默认星级 等级
	voidhero.StarItem.UpStar = GetCsvMgr().GetHeroInitLv(voidhero.HeroId)
	voidhero.LvUp(1 - voidhero.HeroLv)
	voidhero.Resonance = 0
	voidhero.CheckStageTalent()
	self.player.countHeroFight(voidhero, ReasonHeroVoidSetResonance)
	return true
}

// 虚空英雄共鸣发生变化
func (self *ModHero) ChangeVoidHeroResonance(herokey int, voidherokey int) bool {
	// 获得普通英雄
	hero := self.GetHero(herokey)
	if hero == nil { // 找不到则说明是删除流程
		if voidherokey != 0 {
			self.CancelVoidHeroResonance(voidherokey)

			//var backmsg S2C_VoidHeroResonanceCancel
			//backmsg.Cid = MSG_HERO_CANCEL_VOID_HERO_RESONANCE
			//backmsg.VoidHeroKeyId = voidherokey
			//backmsg.VoidHero = self.GetHero(voidherokey)
			//self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

			var backmsg S2C_VoidHeroResonanceUpdate
			backmsg.Cid = MSG_HERO_UPDATE_VOID_HERO_RESONANCE
			backmsg.VoidHero = self.GetHero(voidherokey)
			self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
			return true
		}
		return false
	} else { // 更新流程
		// 如果是共鸣英雄 返回
		if hero.VoidHero != 0 {
			return false
		}
		// 如果没有共鸣则返回
		if hero.Resonance == 0 {
			return false
		}
		// 获得共鸣英雄
		voidhero := self.GetHero(hero.Resonance)
		if voidhero == nil {
			return false
		}

		// 设置星级和等级
		voidhero.StarItem.UpStar = hero.StarItem.UpStar
		voidhero.LvUp(hero.HeroLv - voidhero.HeroLv)
		voidhero.CheckStageTalent()
		self.player.countHeroFight(voidhero, ReasonHeroVoidSetResonance)

		var backmsg S2C_VoidHeroResonanceUpdate
		backmsg.Cid = MSG_HERO_UPDATE_VOID_HERO_RESONANCE
		backmsg.VoidHero = voidhero
		self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
		return true
	}
}
