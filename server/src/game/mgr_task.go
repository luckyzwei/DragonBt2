package game

import "time"

// 返回任务 true记次任务 false记值任务
// n2 表示当前次数, n3,n4表示额外条件
func DoTask(task *TaskNode, player *Player, n2, n3, n4 int) (int, bool) {
	//LogDebug("DoTask, tasktype:", task.Tasktypes, ", n2:", n2, ", n3:", n3, ",n4:", n4)
	//LogDebug("DoTask, task.N1:", task.N1, ", task.N2:", task.N2, ", task.N3:", task.N3, ",task.N4:", task.N4)
	switch task.Tasktypes {
	case TASK_TYPE_PLAYER_LEVEL:
		return player.Sql_UserBase.Level, false
	case TASK_TYPE_CRYSTAL_FIGHT:
		return n2, false
	case TASK_TYPE_GET_HERO:
		if n3 < task.N2 {
			return 0, true
		}
		if task.N3 != 0 && task.N3 != n4 {
			return 0, true
		}
		return n2, true
	case TASK_TYPE_BIGGEST_LEVEL:
		return player.GetModule("hero").(*ModHero).getHeroBiggestLevel(), false
	case TASK_TYPE_LEVEL_UP_COUNT:
		return n2, true
	case TASK_TYPE_STAR_UP_COUNT:
		if n3 >= task.N2 {
			return n2, true
		} else {
			return 0, true
		}
	case TASK_TYPE_HAVE_HERO:
		return player.GetModule("hero").(*ModHero).getHeroNumByLevel(task.N2), false
	case TASK_TYPE_GET_HOOK_AWARD:
		return n2, true
	case TASK_TYPE_GET_HOOK_FAST_AWARD:
		return n2, true
	case TASK_TYPE_MONTH_STATE:
		return 1, true
	case TASK_TYPE_GET_HOOK_GET_ITEM:
		if task.N2 == 0 {
			return n3, true
		} else if task.N2 != 0 && n2 == task.N2 {
			return n3, true
		}
		return 0, true
	case TASK_TYPE_EQUIP_LEVEL_UP:
		return 1, true
	case TASK_TYPE_HAVE_EQUIP:
		return player.GetModule("equip").(*ModEquip).getEquipNumByQuality(task.N2, task.N3), false
	case TASK_TYPE_FINISH_PASS:
		if task.N2 != 0 && n2 != task.N2 {
			return 0, true
		}
		if n3 != task.N3 {
			return 0, true
		}
		return 1, true
	case TASK_TYPE_FINISH_MAIN_PASS:
		return n2, false
	case TASK_TYPE_FINISH_CHAPTER:
		return n2, false
	case TASK_TYPE_ARENA_COUNT:
		if task.N2 == 0 {
			return 1, true
		} else if task.N2 != 0 && task.N2 == n2 {
			return 1, true
		}
		return 0, true
	case TASK_TYPE_ARENA_POINT:
		return n2, false
	case TASK_TYPE_SPECIAL_ARENA_CLASS:
		if n2 == 0 || n3 == 0 {
			return 0, false
		}
		if task.N2 > n2 {
			return 1, false
		} else if task.N2 == n2 {
			if task.N3 >= n3 {
				return 1, false
			}
		}
		return 0, false
	case TASK_TYPE_PIT_PASS:
		if task.N2 != 0 && n2 != task.N2 {
			return 0, true
		}
		return 1, true
	case TASK_TYPE_PIT_KILL:
		if task.N2 != 0 && n2 != task.N2 {
			return 0, true
		}
		return 1, true
	case TASK_TYPE_PIT_KILL_PLAYER:
		if task.N2 != 0 && n2 != task.N2 {
			return 0, true
		}
		return 1, true
	case TASK_TYPE_RESONANCE_CRYSTAL_COUNT:
		return n2, false
	case TASK_TYPE_RESONANCE_CRYSTAL_LEVEL:
		return n2, false
	case TASK_TYPE_RESONANCE_CRYSTAL_SET:
		return 1, true
	case TASK_TYPE_WOTER_LEVEL:
		return n2, false
	case TASK_TYPE_WOTER_COUNT:
		if task.N2 != n2 {
			return 0, true
		}
		return 1, true
	case TASK_TYPE_CAMP_WOTER_LEVEL:
		return handleCampTower(task, player)
	case TASK_TYPE_ONE_CAMP_TOWER_LEVEL:
		if task.N2 == 0 {
			return handleCampTowerLevel(task, player)
		}
		if task.N2 != n3 {
			return 0, false
		}
		return n2, false
	case TASK_TYPE_UNION_HUNT_COUNT:
		if task.N2 != 0 && task.N2 != n2 {
			return 0, true
		}
		if task.N3 != n3 {
			return 0, true
		}
		return 1, true
	case TASK_TYPE_REWARD_TASK_SET:
		if task.N2 != 0 && task.N2 != n2 {
			return 0, true
		}
		return 1, true
	case TASK_TYPE_REWARD_TASK_GET:
		if task.N2 != 0 && task.N2 != n3 {
			return 0, true
		}
		if task.N3 != 0 && task.N3 > n4 {
			return 0, true
		}
		return n2, true
	case TASK_TYPE_REWARD_LEVLE:
		return n2, false
	case TASK_TYPE_FRIEND_POINT_COUNT:
		return n2, true
	case TASK_TYPE_FRIEND_POINT:
		return n2, true
	case TASK_TYPE_SUMMON_HEROS:
		if task.N2 == 0 {
			return n2, true
		} else if task.N2 != 0 && task.N2 == n3 {
			return n2, true
		}
		return 0, true
	case TASK_TYPE_SHOP_BUY_COUNT:
		if task.N2 == 0 {
			return n2, true
		} else if task.N2 != 0 && task.N2 == n3 {
			return n2, true
		}
		return 0, true
	case TASK_TYPE_ADD_LIVENESS:
		return n2, true
	case TASK_TYPE_JOIN_UNION:
		return 1, true
	case TASK_TYPE_LOGIN_TOTAL_COUNT:
		return player.Sql_UserBase.LoginDays, false
	case TASK_TYPE_REG_TOTAL_COUNT:
		now := TimeServer()
		today := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
		if now.Hour() < 5 {
			today -= DAY_SECS
		}
		regtime, _ := time.ParseInLocation(DATEFORMAT, player.Sql_UserBase.Regtime, time.Local)
		regday := time.Date(regtime.Year(), regtime.Month(), regtime.Day(), 5, 0, 0, 0, regtime.Location()).Unix()
		if regtime.Hour() < 5 {
			regday -= DAY_SECS
		}
		return int(today-regday)/DAY_SECS + 1, false
	case TASK_TYPE_DECOMPOSE_HEROS:
		return n2, true
	case TASK_TYPE_HERO_STAR_POINT:
		if task.N1 != 0 {
			if task.N2 == 0 { // N2为0则没有限制直接增加
				return n2, false
			} else if task.N2 == n3 { // N2不为0 则有限制 只有n3符合时才会增加
				return n2, false
			}
			return 0, false // 否则返回0
		}
		return 0, false
	case TASK_TYPE_LOGIN_DAY: // 登录游戏
		if n2 != 0 {
			return 1, true
		}
		return 0, true
	case TASK_TYPE_IS_LOGIN: // 登录游戏
		if n2 != 0 {
			return 1, false
		}
		return 0, false
	case TASK_TYPE_ASTROLOGY_COUNT:
		return n2, true
	case TASK_TYPE_SUMMON_ELITE_HEROS:
		if task.N2 != 0 && task.N2 != n2 {
			return 0, true
		}
		return n2, true
	case TASK_TYPE_RECHARGE_ONCE: // 指定充值档
		if task.N2 != 0 && task.N2 != n2 {
			return 0, true
		}
		return 1, true
	case TASK_TYPE_RECHARGE_MONEY_DAILY: // 指定充值档
		return n2, false
	case TASK_TYPE_RECHARGE_COST: // 累积消耗
		return n2, true
	case TASK_TYPE_RECHARGE_EQUAL_SINGLE: // 单笔充值
		if task.N2 != 0 && task.N2 != n2 {
			return 0, false
		}
		return 1, true
	case TASK_TYPE_RECHARGE_ALL: // 累积充值
		return n2, true
	case TASK_TYPE_JJC_SCORE:
		if task.N2 != 0 && task.N2 != n3 {
			return 0, true
		}
		return n2, true
	case TASK_TYPE_VIP_BUY:
		return player.Sql_UserBase.Vip, false
	case TASK_TYPE_RECHARGE_GOLD_WEEK:
		return 1, true
	case TASK_TYPE_RECHARGE_SINGLE:
		if task.N2 != 0 && n2 < task.N2 {
			return 0, true
		}
		return 1, true
	case TASK_TYPE_STAR_UP_COUNT_EQUAL:
		if n3 == task.N2 {
			return n2, true
		} else {
			return 0, true
		}
	case TASK_TYPE_GET_HERO_WIDE:
		if n3 >= task.N2 {
			return n2, true
		} else {
			return 0, true
		}
	case TASK_TYPE_REWARD_TASK_GET_EQUAL:
		if task.N2 != 0 && task.N2 != n3 {
			return 0, true
		}
		if task.N3 != 0 && task.N3 > n4 {
			return 0, true
		}
		return n2, true
	case TASK_TYPE_INSTANCE_PROCESS:
		if task.N2 != 0 && task.N2 != n2 {
			return 0, false
		}
		if task.N3 != 0 && task.N3 > n3 {
			return 0, false
		}
		return 1, false
	case TASK_TYPE_ACTIFACT_GET:
		if task.N2 != 0 && task.N2 != n2 {
			return 0, false
		}
		return 1, false
	case TASK_TYPE_LIFETREE_MAIN_LEVEL:
		return n2, false
	case TASK_TYPE_LIFETREE_OTHER_LEVEL:
		if task.N2 != 0 && task.N2 != n3 {
			return 0, false
		}
		return n2, false
	case TASK_TYPE_ITEM_GET_COUNT:
		if task.N2 != 0 && task.N2 != n3 {
			return 0, true
		}
		return n2, true
	case TASK_TYPE_INTERSTELLAR_NEBULA:
		if task.N2 != 0 && task.N2 != n2 {
			return 0, false
		}
		return 1, false
	case TASK_TYPE_INTERSTELLAR_NEBULAWAR:
		if task.N2 != 0 && task.N2 != n2 {
			return 0, false
		}
		return 1, false
	case TASK_TYPE_NOBILITY_LEVEL:
		return n2, false
	case TASK_TYPE_ACTIVITY_BOSS_HURT_SINGLE:
		if task.N2 != 0 && task.N2 != n3 {
			return 0, false
		}
		return n2, false
	case TASK_TYPE_ACTIVITY_BOSS_HURT_ALL:
		if task.N2 != 0 && task.N2 != n3 {
			return 0, true
		}
		return n2, true
	case TASK_TYPE_ACTIVITY_BOSS_COUNT:
		if task.N2 != 0 && task.N2 != n3 {
			return 0, true
		}
		return n2, true
	case TASK_TYPE_LUCKY_FIND:
		return n2, true
	}
	return 0, true
}

// 通用关卡判断条件
func handlePassMission(task *TaskNode, param []int) (int, bool) {
	levelId := task.N2
	if levelId != 0 {
		if levelId == param[0] {
			return param[1], true
		} else {
			return 0, true
		}
	}

	return param[1], true
}

// 普通副本,精英副本, 通关次数, 关卡Id[0代表任意关卡]
// param 关卡Id, 次数, 类型[普通或者精英]
func handlePassLevel(task *TaskNode, param []int, levelType int) (int, bool) {
	if levelType != param[2] {
		return 0, true
	}
	return handlePassMission(task, param)
}

// 组队副本, 通关次数, 关卡Id[0代表任意关卡]
// param 参数关卡Id, 次数
func handlePassTeam(task *TaskNode, param []int) (int, bool) {
	return handlePassMission(task, param)
}

// 宝石副本, 通关次数, 关卡Id[0代表任意关卡]
func handlePassGem(task *TaskNode, param []int) (int, bool) {
	return handlePassMission(task, param)
}

// 远征副本, 通关次数, 关卡Id[0代表任意关卡]
func handlePassExpedition(task *TaskNode, param []int) (int, bool) {
	return handlePassMission(task, param)
}

// 科技   完成次数
// param 次数
func handleTechDone(param []int) (int, bool) {
	return param[0], true
}

// 活跃度 完成次数 达标数值
// param 达标数值
func handleLiveness(task *TaskNode, param []int) (int, bool) {
	n2 := task.N2
	if param[0] == n2 {
		return 1, true
	}
	return 0, true
}

// 竞技场 完成次数 达标排名 0代表任意排名
// param 次数 排名
func handlePvpTimes(task *TaskNode, param []int) (int, bool) {
	n2 := task.N2
	if n2 != 0 {
		if param[0] == n2 {
			return 1, true
		} else {
			return 0, true
		}
	} else {
		return 1, true
	}

	return 0, true
}

// 竞技场 累计积分
// param 增加的积分
func handlePvpScore(param []int) (int, bool) {
	return param[0], true
}

// 爬塔 重置次数
func handleTowerReset(param []int) (int, bool) {
	return param[0], true
}

// 爬塔 通关次数 关卡 0代表任意关卡
// param 关卡Id 次数 关卡类型
func handleTowerWin(task *TaskNode, param []int) (int, bool) {
	if task.N2 != 0 {
		if param[0] == task.N2 {
			return param[1], true
		} else {
			return 0, true
		}
	} else {
		return param[1], true
	}
	return 0, true
}

// 爬塔 通关次数 精英关卡 0代表任意关卡
func handleTowerWin2(task *TaskNode, param []int) (int, bool) {
	if task.N2 != 0 {
		if param[0] == task.N2 {
			return param[1], true
		} else {
			return 0, true
		}
	} else {
		return param[1], true
	}
	return 0, true
}

// 天赋最高等级 天赋数量 等级, 记值任务
func handleTalentNum(task *TaskNode, player *Player) (int, bool) {
	//heroes := player.getHeroes()
	num := 0
	/*
		for _, hero := range heroes {
			if hero == nil {
				continue
			}

			if hero.TalentItem == nil {
				continue
			}

			if len(hero.TalentItem.Talents) <= 0 {
				continue
			}

			for _, t := range hero.TalentItem.Talents {
				if t == nil {
					continue
				}
				if t.Lv >= task.N2 {
					num += 1
				}
			}
		}


	*/
	return num, false
}

// 拥有英雄信息
func handleOwnHero(task *TaskNode, player *Player) (int, bool) {
	total := 0
	for _, value := range player.GetModule("hero").(*ModHero).Sql_Hero.info {
		config := GetCsvMgr().GetHeroConfig(value.HeroId)
		if config == nil {
			continue
		}
		if config.FightIndex < task.N2 {
			continue
		}
		if task.N3 != 0 && config.HeroCamp != task.N3 {
			continue
		}
		if task.N4 != 0 && config.HeroArms != task.N4 {
			continue
		}
		total++
	}
	return total, false
}

// 英雄升星
func handleHeroStar(task *TaskNode, player *Player) (int, bool) {
	hero := player.GetModule("hero").(*ModHero).GetHero(task.N2)
	if hero == nil {
		return 0, false
	}
	if hero.StarItem.UpStar < task.N4 {
		return 0, false
	}
	return 1, true
}

// 招财
func handleBuyMoney(task *TaskNode, param []int) (int, bool) {
	if len(param) <= 0 {
		return 0, true
	}
	n2 := param[0]
	if task.N2 == 0 || n2 == task.N2 {
		return 1, true
	}
	return 0, true
}

// 酒馆购买
func handleFind(task *TaskNode, param []int) (int, bool) {
	if len(param) <= 0 {
		return 0, true
	}
	n2 := param[0]
	n3 := param[1]
	if task.N2 != 0 && task.N2 != n2 {
		return 0, true
	}
	return n3, true
}

// 商店任务
func handleShop(task *TaskNode, n2 int) (int, bool) {
	if task.N2 != 0 && task.N2 != n2 {
		return 0, true
	}
	return 1, true
}

//215 军团捐献
func handleUnionDonation(task *TaskNode, param []int) (int, bool) {
	if len(param) <= 0 {
		return 0, true
	}

	nDonateMode := param[0]

	if task.N2 == 0 {
		return 1, true
	} else {
		if nDonateMode == task.N2 {
			return 1, true
		}
	}

	return 0, true
}

//217 佣兵挑战
func handleArmyFightCount(task *TaskNode, param []int) (int, bool) {
	if len(param) <= 0 {
		return 0, true
	}

	nFightMode := param[0]

	if task.N2 == 0 {
		return 1, true
	} else {
		if nFightMode == task.N2 {
			return 1, true
		}
	}

	return 0, true
}

//223 十连抽
func handleSearchTenTime(task *TaskNode, param []int) (int, bool) {
	if len(param) <= 0 {
		return 0, true
	}

	nSearchMode := param[0]

	if task.N2 == 0 {
		return 1, true
	} else {
		if nSearchMode == task.N2 {
			return 1, true
		}
	}

	return 0, true
}

// 神格数量 神格等级 英雄资质 神格点
func handleDinivityNum(task *TaskNode, player *Player) (int, bool) {
	/*
		heroes := player.getHeroes()
		num := 0
		for _, hero := range heroes {
			if hero == nil {
				continue
			}

			config := GetCsvMgr().GetHeroConfig(hero.HeroId)

			if config == nil {
				continue
			}

			if hero.TalentItem == nil {
				continue
			}

			if len(hero.TalentItem.Talents) <= 0 {
				continue
			}

			for i, t := range hero.TalentItem.Talents {
				if t == nil {
					continue
				}

				bLevel, bFightIndex, bPointX := false, false, false

				if task.N2 > 0 {
					if t.Lv >= task.N2 {
						bLevel = true
					}
				} else {
					bLevel = true
				}

				if task.N3 > 0 {
					if config.FightIndex >= task.N3 {
						bFightIndex = true
					}
				} else {
					bFightIndex = true
				}

				if task.N4 > 0 {
					if len(config.Point) >= task.N4 {
						if i+1 == task.N3 {
							bPointX = true
						}
					}
				} else {
					bPointX = true
				}

				if bLevel && bFightIndex && bPointX {
					num++
				}
			}
		}

		return num, false

	*/
	return 0, false
}

// 英雄数量 神格总等级 英雄资质
func handleHeroDinivity(task *TaskNode, player *Player) (int, bool) {
	//heroes := player.getHeroes()
	num := 0
	/*
		for _, hero := range heroes {
			if hero == nil {
				continue
			}

			config := GetCsvMgr().GetHeroConfig(hero.HeroId)

			if config == nil {
				continue
			}

			if hero.TalentItem == nil {
				continue
			}

			if len(hero.TalentItem.Talents) <= 0 {
				continue
			}

			nLevel := 0
			for _, t := range hero.TalentItem.Talents {
				if t == nil {
					continue
				}

				nLevel += t.Lv
			}

			bLevel, bFightIndex := false, false

			if task.N2 > 0 {
				if nLevel >= task.N2 {
					bLevel = true
				}
			} else {
				bLevel = true
			}

			if task.N3 > 0 {
				if config.FightIndex >= task.N3 {
					bFightIndex = true
				}
			} else {
				bFightIndex = true
			}

			if bLevel && bFightIndex {
				num++
			}
		}


	*/
	return num, false
}

// 神格幻境抽
func handleDinivityDreamLand(task *TaskNode, param []int) (int, bool) {
	if len(param) <= 0 {
		return 0, true
	}

	nSearchMode := param[0]
	nTimes := param[1]

	if task.N2 == 0 {
		return nTimes, true
	} else {
		if nSearchMode == task.N2 {
			return nTimes, true
		}
	}

	return 0, true
}

// 神格幻境十连抽
func handleDinivityDreamLandTen(task *TaskNode, param []int) (int, bool) {
	if len(param) <= 0 {
		return 0, true
	}

	nSearchMode := param[0]
	nTimes := param[1]

	if nTimes != 10 {
		return 0, true
	}

	if task.N2 == 0 {
		return 1, true
	} else {
		if nSearchMode == task.N2 {
			return 1, true
		}
	}

	return 0, true
}

// 获得达标的塔的数量
func handleCampTower(task *TaskNode, player *Player) (int, bool) {
	nCount := 0
	info := player.GetModule("tower").(*ModTower).data.info
	for _, v := range info {
		if v.Type == TOWER_TYPE_0 {
			continue
		}

		if v.MaxLevel >= task.N2 {
			nCount++
		}
	}

	return nCount, false
}

func handleCampTowerLevel(task *TaskNode, player *Player) (int, bool) {
	nCount := 0
	info := player.GetModule("tower").(*ModTower).data.info
	for _, v := range info {
		if v.Type == TOWER_TYPE_0 {
			continue
		}

		if v.MaxLevel > nCount {
			nCount = v.MaxLevel
		}
	}

	return nCount, false
}
