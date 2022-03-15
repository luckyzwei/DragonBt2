package game

import "encoding/json"

// 单条的每一层
type StageTalentIndex struct {
	ID    int `json:"id"`    //! id
	Index int `json:"index"` //! 层
	Pos   int `json:"pos"`   //! 附带的选择技能
}

//天赋组
type StageTalent struct {
	Group    int                 `json:"group"`    //! 组别
	AllSkill []*StageTalentIndex `json:"allskill"` // 全部技能
}

// 设置天赋技能
func (self *ModHero) MsgStageTalentSetSkill(body []byte) {
	var msg C2S_SetStageTalentSkill
	json.Unmarshal(body, &msg)

	hero := self.GetHero(msg.HeroKeyId)
	if nil == hero {
		return
	}

	hero.CheckStageTalent()

	skill := hero.StageTalent.GetTalentSkill(msg.Index)
	if nil == skill {
		return
	}

	config := GetCsvMgr().GetStageTalent(skill.ID)
	if nil == config {
		return
	}

	if msg.Pos <= 0 || msg.Pos > len(config.Skill) {
		return
	}

	skill.Pos = msg.Pos

	self.player.countHeroFight(hero, ReasonStageTalentSetSkill)

	var backmsg S2C_SetStageTalentSkill
	backmsg.Cid = MSG_HERO_STAGE_TALENT_SET_SKILL
	backmsg.Index = msg.Index
	backmsg.Pos = msg.Pos
	backmsg.HeroKeyId = msg.HeroKeyId
	backmsg.Info = skill
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 当英雄升星 检测天赋解锁
func (self *Hero) CheckStageTalent() {
	star := self.StarItem.UpStar

	//获取升级的配置
	heroConfig := GetCsvMgr().GetHeroMapConfig(self.HeroId, star)
	if heroConfig == nil {
		return
	}

	talentConfig := GetCsvMgr().GetStageTalentMap(heroConfig.TalentGroup)
	if nil == talentConfig {
		return
	}

	if self.StageTalent == nil {
		self.StageTalent = &StageTalent{heroConfig.TalentGroup, []*StageTalentIndex{}}
	}

	if heroConfig.TalentGroup != self.StageTalent.Group {
		self.StageTalent.Group = heroConfig.TalentGroup
	}

	// 循环配置
	for _, config := range talentConfig {
		// 技能错误
		if len(config.Skill) <= 0 {
			continue
		}
		// 未开启跳出
		if star >= config.Open {
			// 技能已经解锁
			skill := self.StageTalent.GetTalentSkill(config.Index)
			if skill != nil {
				continue
			}

			// 添加
			self.StageTalent.AddTalentSkill(config, 1)
		} else {
			// 技能未解锁
			skill := self.StageTalent.GetTalentSkill(config.Index)
			if skill == nil {
				continue
			}

			// 删除
			self.StageTalent.RemoveTalentSkill(config.Index)
		}
	}
}

// 获得技能信息
func (self *StageTalent) GetTalentSkill(index int) *StageTalentIndex {
	for _, skill := range self.AllSkill {
		if skill.Index == index {
			return skill
		}
	}
	return nil
}

func (self *StageTalent) AddTalentSkill(config *StageTalentConfig, pos int) {
	// 添加技能 默认给第一个技能
	self.AllSkill = append(self.AllSkill, &StageTalentIndex{config.ID, config.Index, pos})
}

// 删除某个技能
func (self *StageTalent) RemoveTalentSkill(index int) bool {
	for i, skill := range self.AllSkill {
		if skill.Index == index {
			self.AllSkill = append(self.AllSkill[:i], self.AllSkill[i+1:]...)
			return true
		}
	}
	return false
}

func (self *Hero) GetStageTalentAttr() map[int]*Attribute {
	attrNew := make(map[int]*Attribute)
	if self.StageTalent == nil {
		return attrNew
	}
	fightAtt := make(map[int]*Attribute)
	skill := []int{}
	for _, value := range self.StageTalent.AllSkill {
		config := GetCsvMgr().GetStageTalent(value.ID)
		if nil == config {
			continue
		}

		if value.Pos <= 0 || value.Pos > len(config.Skill) {
			continue
		}

		skill = append(skill, config.Skill[value.Pos-1])

		_, ok := fightAtt[config.Type]
		if ok {
			fightAtt[config.Type].AttValue += int64(config.Value)
		} else {
			fightAtt[config.Type] = &Attribute{config.Type, int64(config.Value)}
		}
	}
	//计算技能属性
	skillAttr := GetSkillAttr(skill)
	self.addAttEx(skillAttr, attrNew)
	self.addAttEx(fightAtt, attrNew)
	//if len(skillAttr) > 0 {
	//	//这个地方需要把万分比属性计算掉
	//	for _, v := range attrNew {
	//		if v.AttType > AttrDisExt && v.AttType <= AttrEnd+AttrDisExt {
	//			attr, ok := attrNew[v.AttType-AttrDisExt]
	//			if ok {
	//				attr.AttValue = attr.AttValue * (1.0 + v.AttValue/10000.0)
	//			}
	//		}
	//	}
	//	for _, v := range attrNew {
	//		if v.AttType > AttrDisExt && v.AttType <= AttrEnd+AttrDisExt {
	//			continue
	//		}
	//		attrNew[v.AttType] = v
	//	}
	//}
	return attrNew
}
