package game

//属性定义
const (
	AttrHp      = 1
	AttrAttack  = 2
	AttrDefence = 3
	AttrFight   = 99
	AttrEnd     = 64

	AttrDisExt = 500 //属性百分比的偏移值
)

// 属性包装类
type AttrWrapper struct {
	Base     []float64        // 固定属性
	Ext      map[int]float64  // 扩展属性
	Per      map[int]float64  // 百分比属性
	Energy   int              // 怒气恢复
	FightNum int64            // 战斗力值
	ExtRet   []JS_HeroExtAttr // 最终扩展属性
}

func (self *Hero) GetBase() map[int]*Attribute {
	res := GetCsvMgr().getHeroBase(self.HeroId)
	return res
}

func (self *Hero) GetGrowth(playerLv int) map[int]*Attribute {
	res := GetCsvMgr().getHeroGrowth(self.HeroId, playerLv)
	return res
}

func (self *Hero) addAttEx(inparam map[int]*Attribute, attMap map[int]*Attribute) {
	if inparam == nil || attMap == nil {
		return
	}

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

func ProcAtt(attMap map[int]*Attribute, pAttr *AttrWrapper) int {
	for _, pAttribute := range attMap {
		valuetype := pAttribute.AttType
		value := pAttribute.AttValue
		if valuetype < 0 || value == 0 {
			continue
		}

		if valuetype >= 0 && valuetype < len(pAttr.Base) { // 基础属性
			pAttr.Base[valuetype] += float64(value)
		} else if valuetype == 99 { // 战斗力
			pAttr.FightNum += value //这里不好处理int64 暂时搁置 20190506 by zy
		} else if valuetype == 2500 {
			pAttr.Energy += int(value)
		} else if valuetype >= AttrDisExt && valuetype < AttrEnd+AttrDisExt { // 百分比
			pAttr.Per[valuetype] += float64(value)
		} else {
			pAttr.Ext[valuetype] += float64(value)
		}
	}
	return 0
}

// 基础属性
func ProcLast(playerLv int, pAttr *AttrWrapper) {
	// 最后计算
	for i := 0; i < AttrEnd; i++ {
		value, ok := pAttr.Per[i+AttrDisExt]
		if !ok {
			continue
		}
		//0.005用于消除误差
		pAttr.Base[i] = pAttr.Base[i] * (1.0 + value/10000.0)
	}

	//! 计算基础属性
	/*
		heroLv := float32(playerLv - 1)
		pAttr.Base[0] = pAttr.Base[0] + pAttr.Base[3]*heroLv
		// 统率=增加生命百分比
		pAttr.Base[1] = pAttr.Base[1] + pAttr.Base[4]*heroLv
		pAttr.Base[2] = pAttr.Base[2] + pAttr.Base[5]*heroLv

		for i := 0; i <= 5; i++ {
			pAttr.Base[i] /= 100
		}

		//! 生命
		pAttr.Base[6] = pAttr.Base[6] + pAttr.Base[1]*10
		//! 物理攻击力
		pAttr.Base[7] = pAttr.Base[7] + pAttr.Base[0]*1.5
		//! 魔法强度
		pAttr.Base[9] = pAttr.Base[9] + pAttr.Base[2]*1.5
		//! 物理护甲
		pAttr.Base[8] = pAttr.Base[8] + pAttr.Base[0]*0.5
		//! 魔法抗性
		pAttr.Base[10] = pAttr.Base[10] + pAttr.Base[2]*0.5

		for i := 6; i < 32; i++ {
			value, ok := pAttr.Per[i+100]
			if !ok {
				continue
			}
			pAttr.Base[i] = pAttr.Base[i] * (1.0 + value/10000.0)
		}

	*/
}

// 扩展属性
func ProcExtAttr(pAttr *AttrWrapper) []JS_HeroExtAttr {
	// 额外属性
	extRet := make([]JS_HeroExtAttr, 0)
	for k := range pAttr.Ext {
		value, ok := pAttr.Per[k+100]
		if ok {
			pAttr.Ext[k] += pAttr.Ext[k] * (1.0 + value/10000.0)
		}
		extRet = append(extRet, JS_HeroExtAttr{k, pAttr.Ext[k]})
	}
	return extRet
}

// 重写战斗计算
// 计算英雄战斗力就是英雄战斗力, 属性类型=99的加战斗力
// 计算战斗力应该从玩家身上发出来
func (self *Hero) CountFight(player *Player) int64 {
	pAttr := self.GetAttr(player)
	self.Fight = pAttr.FightNum
	return self.Fight
}

// 计算玩家战斗力, 从玩家身上发出来, 而不是英雄
func (self *Player) countHeroFight(hero *Hero, reason int) {
	if hero == nil {
		LogError("hero is nil")
		return
	}

	oldFight := hero.Fight
	//LogDebug("oldFight:", oldFight)
	hero.CountFight(self)
	self.updateFight()
	self.GetModule("friend").(*ModFriend).UpdateHeroSet(hero.HeroKeyId)
	self.GetModule("crystal").(*ModResonanceCrystal).UpdateMaxFightAll()
	afterFight := hero.Fight
	//LogDebug("afterFight:", afterFight)

	//ReasonOutTeam 需要强行同步
	if oldFight != afterFight || reason == 0 {
		self.synFight(hero.HeroKeyId, afterFight, reason, hero.HeroLv)
	}
}

func (self *Player) checkHeroFight(hero *Hero, reason int) {
	if hero == nil {
		LogError("hero is nil")
		return
	}
	hero.CountFight(self)
	//self.checkUpdateFight(hero.HeroId)

}

func (self *Player) checkUpdateFight(updateId int) {
	team := self.getTeam()
	_, ok := team[updateId]
	if !ok {
		return
	}
	var fightValue int64 = 0
	for heroId := range team {
		heroConfig := GetCsvMgr().GetHeroConfig(heroId)
		if heroConfig == nil {
			continue
		}
		/*
			if heroConfig.HeroType == 1 {
				hero := self.GetModule("hero").(*ModHero).GetHero(heroId)
				if hero == nil {
					LogError("hero is nil")
					continue
				}
				fightValue += hero.Fight
			} else if heroConfig.HeroType == 2 {
				fightValue += self.GetModule("boss").(*ModBoss).Fight
			}

		*/
	}
	//self.Sql_UserBase.Fight = fightValue
	self.SetFight(fightValue, 1)
}

func (self *Player) GetTeamFight(teamType int) int64 {
	team := self.GetModule("team").(*ModTeam).getTeamPos(teamType)
	if team == nil {
		return 0
	}
	var fightValue int64 = 0
	for _, heroKeyId := range team.FightPos {
		hero := self.getHero(heroKeyId)
		if hero == nil {
			continue
		}
		fightValue += hero.Fight
	}
	//hydraFight := self.GetModule("hydra").(*ModHydra).GetHydroFight(team.HydraId)
	//fightValue += hydraFight
	return fightValue
}

// 更新玩家战斗力, 如果英雄在战队里面,则进行更新, 不在, 就不需要更新
func (self *Player) updateFight() {

	var msg S2C_CompareTeamFight
	msg.Cid = "compareteamfight"

	team := self.GetModule("team").(*ModTeam).getTeamPos(TEAMTYPE_DEFAULT)
	if team == nil {
		return
	}
	var fightValue int64 = 0
	heroIds := make([]int, 0)
	for _, heroKeyId := range team.FightPos {
		hero := self.getHero(heroKeyId)
		if hero == nil {
			continue
		}
		fightValue += hero.Fight
		msg.HeroFight = append(msg.HeroFight, hero.Fight)
		msg.HeroKeyId = append(msg.HeroKeyId, heroKeyId)
		heroIds = append(heroIds, hero.HeroId)
	}
	msg.Fight = fightValue
	msg.attr = GetRobotMgr().GetTeamAttr(heroIds)
	msg.FightInfo = GetRobotMgr().GetPlayerFightInfoByPos(self, 0, 0, TEAMTYPE_DEFAULT)
	self.SendMsg(msg.Cid, HF_JtoB(msg))
	if fightValue != self.Sql_UserBase.Fight {
		self.SetFight(fightValue, 4)
	}
	//self.HandleTask(PlayerFightTask, 0, 0, 0)
	if self.Sql_UserBase.Fight > 0 {
		res := GetTopFightMgr().SyncFight(self.Sql_UserBase.Fight, self)
		if res {
			//GetTopActMgr().updateFight()
		}
	}
}

// 更新全队战力 + 巨兽 + 科技
func (self *Player) countTeamFight(reason int) {
	team := self.getTeamPos()
	if team == nil {
		return
	}
	for heroKeyId := range team.FightPos {
		if heroKeyId == 0 {
			continue
		}
		hero := self.GetModule("hero").(*ModHero).GetHero(heroKeyId)
		if hero == nil {
			continue
		}
		hero.CountFight(self)
	}

	self.updateFight()
}

func (self *Player) countAllHero() {
	heroes := self.GetModule("hero").(*ModHero).GetHeroes()
	for _, v := range heroes {
		v.CountFight(self)
	}
	self.updateFight()
	self.GetModule("crystal").(*ModResonanceCrystal).UpdateMaxFightAll()
}

//注意:生成fightinfo的时候需要加上阵容加成，计算英雄战力的时候不需要
func (self *Hero) GetHeroAttr(player *Player, teamAttr map[int]*Attribute) (map[int]*Attribute, map[int]*Attribute) {
	attMap := make(map[int]*Attribute)
	// 可以理解为英雄属性
	upStar := self.GetStarAttr()
	self.addAttEx(upStar, attMap)
	// 新天赋
	talent := self.GetStageTalentAttr()
	self.addAttEx(talent, attMap)
	// 装备属性
	equip := player.GetModule("equip").(*ModEquip).getAttr(self.HeroKeyId)
	self.addAttEx(equip, attMap)
	articactEquip := player.GetModule("artifactequip").(*ModArtifactEquip).getAttr(self.HeroKeyId)
	self.addAttEx(articactEquip, attMap)
	// 美人属性
	beauty := player.GetModule("beauty").(*ModBeauty).getAttr()
	self.addAttEx(beauty, attMap)
	// 战马属性
	horse := player.GetModule("horse").(*ModHorse).GetHorseAttrInfo(self.HeroKeyId)
	self.addAttEx(horse, attMap)
	//专属属性计算放最后,并且不享受其他加成
	exclusiveEquip := player.GetModule("equip").(*ModEquip).getAttrExclusive(self.HeroKeyId)
	self.addAttEx(exclusiveEquip, attMap)
	fate := player.GetModule("entanglement").(*ModEntanglement).GetAllProperty(self.HeroId)
	self.addAttEx(fate, attMap)
	tree := player.GetModule("lifetree").(*ModLifeTree).GetAllProperty(self.HeroId, self.StarItem.UpStar)
	self.addAttEx(tree, attMap)
	self.addAttEx(teamAttr, attMap)

	return attMap, upStar
}

//! 计算属性值
func (self *Hero) GetAttr(player *Player) *AttrWrapper {
	pAttr := &AttrWrapper{
		Base:     make([]float64, 32),
		Ext:      make(map[int]float64),
		Per:      make(map[int]float64),
		Energy:   0,
		FightNum: 0,
	}
	self.cacStarAtt()
	playerLv := player.Sql_UserBase.Level
	attMap, _ := self.GetHeroAttr(player, nil)

	ProcAtt(attMap, pAttr)
	ProcLast(playerLv, pAttr)
	//专属放最后
	//player.GetModule("equip").(*ModEquip).getAttrExclusive(self.HeroKeyId, attStar, pAttr)
	//pAttr.ExtRet = ProcExtAttr(pAttr)

	return pAttr
}

func (self *Player) countInit(reason int) {
	team := self.getFirstTeam()
	if len(team) <= 0 {
		return
	}

	var fightValue int64 = 0
	for _, heroKeyId := range team {
		if heroKeyId == 0 {
			continue
		}

		hero := self.GetModule("hero").(*ModHero).GetHero(heroKeyId)
		if hero == nil {
			continue
		}
		hero.CountFight(self)
		fightValue += hero.Fight
	}
	//self.Sql_UserBase.Fight = fightValue
	self.SetFight(fightValue, 2)
}

//! 鼓舞
func (self *Hero) GetAttr2(player *Player, param int, teamAttr map[int]*Attribute) *AttrWrapper {
	pAttr := &AttrWrapper{
		Base:     make([]float64, AttrEnd),
		Ext:      make(map[int]float64),
		Per:      make(map[int]float64),
		Energy:   0,
		FightNum: 0,
	}
	self.cacStarAtt()
	playerLv := player.Sql_UserBase.Level
	attMap, _ := self.GetHeroAttr(player, teamAttr)
	ProcAtt(attMap, pAttr)
	ProcLast(playerLv, pAttr)
	//专属放最后
	//player.GetModule("equip").(*ModEquip).getAttrExclusive(self.HeroKeyId, attStar, pAttr)
	//pAttr.ExtRet = ProcExtAttr(pAttr)
	return pAttr
}

//! 新结构
func (self *NewHero) CalAttr() *AttrWrapper {
	pAttr := &AttrWrapper{
		Base:     make([]float64, 32),
		Ext:      make(map[int]float64),
		Per:      make(map[int]float64),
		Energy:   0,
		FightNum: 0,
	}
	self.cacStarAtt()
	attMap := self.GetHeroAttr()

	ProcAtt(attMap, pAttr)
	ProcLast(0, pAttr)
	pAttr.ExtRet = ProcExtAttr(pAttr)

	self.Fight = pAttr.FightNum

	return pAttr
}

func (self *NewHero) addAttEx(inparam map[int]*Attribute, attMap map[int]*Attribute) {
	if inparam == nil || attMap == nil {
		return
	}

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

func (self *NewHero) GetHeroAttr() map[int]*Attribute {
	attMap := make(map[int]*Attribute)

	// 可以理解为英雄属性
	if self.StarItem != nil {
		self.addAttEx(self.StarItem.AttrMap, attMap)
	}

	config := GetCsvMgr().GetHeroMapConfig(self.HeroId, self.StarItem.UpStar)
	if config == nil {
		return attMap
	}
	//神器
	if len(self.ArtifactEquipIds) > 0 && self.ArtifactEquipIds[0] != nil {
		for _, v := range self.ArtifactEquipIds[0].AttrInfo {
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
	}
	//专属
	if self.ExclusiveEquip != nil && self.ExclusiveEquip.UnLock == LOGIC_TRUE {
		for _, v := range self.ExclusiveEquip.AttrInfo {
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
	}

	if self.StageTalent != nil {
		// 天赋
		for _, value := range self.StageTalent.AllSkill {
			if value == nil {
				continue
			}
			config := GetCsvMgr().GetStageTalent(value.ID)
			if nil == config {
				continue
			}

			if value.Pos <= 0 || value.Pos > len(config.Skill) {
				continue
			}

			_, ok := attMap[config.Type]
			if ok {
				attMap[config.Type].AttValue += int64(config.Value)
			} else {
				attMap[config.Type] = &Attribute{config.Type, int64(config.Value)}
			}
		}
	}

	return attMap
}
