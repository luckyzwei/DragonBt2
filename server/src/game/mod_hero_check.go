package game

func (self *ModHero) CheckHeroOffice(heros [][]*JS_HeroParam, rate int) bool {
	return false
}

func (self *ModHero) CheckHeroAreas(team int, heros [][]*JS_HeroParam, rate int) bool {
	//info := GetPvpMgr().GetInfoFromUid(self.player.Sql_UserBase.Uid, self.player.GetCamp())
	//if info == nil {
	//	return false
	//}
	//
	//if len(heros) <= 0 {
	//	return false
	//}
	//return self.CheckParamOne(info.fightinfo.HeroParam, heros[0], rate)
	return false
}

func (self *ModHero) CheckParamOne(org []JS_HeroParam, target []*JS_HeroParam, rate int) bool {
	if len(org) != len(target) {
		return false
	}

	for i := 0; i < len(org); i++ {
		for j := 0; j < len(org[i].Param); j++ {
			if target[i].Heroid != org[i].Heroid {
				continue
			}
			if target[i].Param[j] > 0 {
				if j == 31 {
					continue
				}
				rate1 := org[i].Param[j] / target[i].Param[j]

				if int(rate1*100) > 100+rate {
					return false
				}

				if rate1 < 1 {
					if j == 11 {
						if rate1 < 0.5 {
							return false
						}
					} else {
						if int(rate1*100) < 10000/(100+rate) {
							return false
						}
					}

				}
			}

		}
	}

	return true
}

func (self *ModHero) CheckParam(org [][]JS_HeroParam, target [][]*JS_HeroParam, rate int) bool {

	if len(org) != len(target) {
		return false
	}

	for i := 0; i < len(org); i++ {
		if self.CheckParamOne(org[i], target[i], rate) == false {
			return false
		}
	}

	return true
}
