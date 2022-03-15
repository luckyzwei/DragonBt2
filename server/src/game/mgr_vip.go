package game

func (self *CsvMgr) GetVipConfig(vip int) *VipConfig {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return nil
	}
	return config
}

func (self *CsvMgr) GetVipPeople(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.People
}

func (self *CsvMgr) GetVipMeetHero(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.Meethero
}

func (self *CsvMgr) GetVipElite(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.Elite
}

func (self *CsvMgr) GetVipRc(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.Resourcechallenge
}

func (self *CsvMgr) GetVipBuySkill(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.Buyskill
}

func (self *CsvMgr) GetVipGp(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.GemsMopping
}

func (self *CsvMgr) GetVipGb(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.GemsBuy
}

func (self *CsvMgr) GetVipVisit(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.Visit
}

func (self *CsvMgr) GetVipKing(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.CrownFight
}

func (self *CsvMgr) GetVipVisitTime(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.Visittime
}

func (self *CsvMgr) GetVipHeroSent(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.Herosent
}

func (self *CsvMgr) GetVipNeedExp(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.Needexp
}

func (self *CsvMgr) GetVipSummonDiscount(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 100
	}
	return config.SummonDiscount
}

func (self *CsvMgr) GetVipSummonTimes(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.SummonTimes
}

func (self *CsvMgr) GetArmyBuy(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.ArmyBuy
}

func (self *CsvMgr) GetVipSkillNumber(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.Skillnumber
}

func (self *CsvMgr) GetVipConsumeTop(level int) int {
	config, ok := self.VipConfigMap[level]
	if !ok {
		return 0
	}
	return config.Consumetopnum
}

func (self *CsvMgr) GetVipPhysical(level int) int {
	config, ok := self.VipConfigMap[level]
	if !ok {
		return 0
	}
	return config.Physical
}

func (self *CsvMgr) GetVipAlchemy(level int) int {
	config, ok := self.VipConfigMap[level]
	if !ok {
		return 0
	}
	if len(config.Alchemys) < 0 {
		return 0
	}
	return config.Alchemys[0]
}

func (self *CsvMgr) GetVipHorseCall(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}

	return config.Horsehigtcall
}

func (self *CsvMgr) GetDungeonReset(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.DungeonReset
}

func (self *CsvMgr) GetDungeonSweep(vip int) int {
	config, ok := self.VipConfigMap[vip]
	if !ok {
		return 0
	}
	return config.DungeonsSweep
}
