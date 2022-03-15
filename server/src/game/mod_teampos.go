package game

import (
	"errors"
)

const (
	MAX_UI_POS    = 6
	MAX_FIGHT_POS = 5
	MAX_GENE_POS  = 9
	PRETEAM_MAX   = 11 //最大编组数
)

const (
	TEAMTYPE_DEFAULT                  = 1
	TEAMTYPE_TOWER_MAIN               = 2
	TEAMTYPE_TOWER_1                  = 3
	TEAMTYPE_TOWER_2                  = 4
	TEAMTYPE_TOWER_3                  = 5
	TEAMTYPE_TOWER_4                  = 6
	TEAMTYPE_ARENA_1                  = 7  // 进攻
	TEAMTYPE_ARENA_2                  = 8  // 防御
	TEAMTYPE_UNION_HUNT               = 9  // 公会狩猎
	TEAMTYPE_ARENA_SPECIAL_1          = 10 // 进攻
	TEAMTYPE_ARENA_SPECIAL_2          = 11 // 进攻
	TEAMTYPE_ARENA_SPECIAL_3          = 12 // 进攻
	TEAMTYPE_ARENA_SPECIAL_4          = 13 // 防守
	TEAMTYPE_ARENA_SPECIAL_5          = 14 // 防守
	TEAMTYPE_ARENA_SPECIAL_6          = 15 // 防守
	TEAMTYPE_NEW_PIT                  = 16 // 地牢
	TEAMTYPE_INSTANCE                 = 17 // 时光之巅
	TEAMTYPE_ACTIVITY_BOSS_1700       = 18 // 暗域入侵 1700
	TEAMTYPE_ACTIVITY_BOSS_1701       = 19 // 暗域入侵 1701
	TEAMTYPE_ACTIVITY_BOSS_1702       = 20 // 暗域入侵 1702
	TEAMTYPE_ACTIVITY_BOSS_1703       = 21 // 暗域入侵 1703
	TEAMTYPE_ACTIVITY_BOSS_1704       = 22 // 暗域入侵 1704
	TEAMTYPE_ACTIVITY_BOSS_1705       = 23 // 暗域入侵 1705
	TEAMTYPE_ACTIVITY_BOSS_1706       = 24 // 暗域入侵 1706
	TEAMTYPE_ACTIVITY_BOSS_1707       = 25 // 暗域入侵 1707
	TEAMTYPE_ACTIVITY_BOSS_1708       = 26 // 暗域入侵 1708
	TEAMTYPE_ACTIVITY_BOSS_1709       = 27 // 暗域入侵 1709
	TEAMTYPE_CROSSARENA_ATTACK        = 28 // 跨服竞技场进攻阵容
	TEAMTYPE_CROSSARENA_DEFENCE       = 29 // 跨服竞技场防守阵容
	TEAMTYPE_ACTIVITY_BOSS_FESTIVAL   = 30 // 活动BOSS
	TEAMTYPE_ACTIVITY_CONSUMERTOP     = 31 // 无双神将
	TEAMTYPE_CROSSARENA_ATTACK_3V3_1  = 32 // 跨服竞技场3v3进攻阵容
	TEAMTYPE_CROSSARENA_ATTACK_3V3_2  = 33 // 跨服竞技场3v3进攻阵容
	TEAMTYPE_CROSSARENA_ATTACK_3V3_3  = 34 // 跨服竞技场3v3进攻阵容
	TEAMTYPE_CROSSARENA_DEFENCE_3V3_1 = 35 // 跨服竞技场3v3防守阵容
	TEAMTYPE_CROSSARENA_DEFENCE_3V3_2 = 36 // 跨服竞技场3v3防守阵容
	TEAMTYPE_CROSSARENA_DEFENCE_3V3_3 = 37 // 跨服竞技场3v3防守阵容
	TEAM_END                          = 38
)

type TeamPos struct {
	FightPos [MAX_FIGHT_POS]int `json:"fightpos"`
	HydraId  int                `json:"hydraid"`
	Gene     [MAX_GENE_POS]int  `json:"gene"`
}

// 魔宠, 虎符, 佣兵, 军旗, 装备, 绑在阵容上
type TeamAttr struct {
	//HorseKeyId int   `json:"horse_key_id"`
	//TigerKeyId int   `json:"tiger_key_id"`
	//ArmyId  int   `json:"army_id"`
	//FlagIds []int `json:"flag_ids"`
	//EquipIds   []int `json:"equipIds"`
}

type Js_TeamPos struct {
	TeamType int      `json:"teamtype"`
	TeamPos  *TeamPos `json:"teampos"`
}

//需要通过阵型找空位  20190923  by zy
func (self *TeamPos) getEmptyFightPos() int {

	for i := 0; i < MAX_FIGHT_POS; i++ {
		if self.FightPos[i] == 0 {
			return i
		}
	}
	return -1
}

func (self *TeamPos) getFightPos() int {
	// 首先上5
	if self.FightPos[4] == 0 {
		return 4
	}

	// 其次从1~4开始
	for i := 0; i < MAX_FIGHT_POS; i++ {
		if self.FightPos[i] != 0 {
			continue
		}

		if i == 4 {
			continue
		}

		return i
	}

	return -1
}

func (self *TeamPos) addFightPos(heroId int) error {
	//巨兽先屏蔽
	if heroId > 5000 {
		return nil
	}
	index := self.getEmptyFightPos()
	if index != -1 {
		self.FightPos[index] = heroId
	} else {
		return errors.New("no empty pos")
	}
	return nil
}

func (self *TeamPos) swapFightPosByIndex(index1 int, index2 int) error {
	if index1 < 0 || index1 >= MAX_FIGHT_POS {
		return errors.New("invalid index")
	}

	if index2 < 0 || index2 >= MAX_FIGHT_POS {
		return errors.New("invalid index")
	}

	self.FightPos[index1], self.FightPos[index2] = self.FightPos[index2], self.FightPos[index1]

	return nil
}

// 英雄上阵
func (self *TeamPos) addUIPos(index int, heroKeyId int) (error, int) {
	if index < 0 || index >= MAX_FIGHT_POS {
		return errors.New(GetCsvMgr().GetText("STR_MOD_TEAMPOS_INDEX_IS_NOT_VALID")), 0
	}

	// 已经上阵的英雄无法再上阵
	for i := 0; i < len(self.FightPos); i++ {
		if self.FightPos[i] == heroKeyId {
			return errors.New(GetCsvMgr().GetText("STR_MOD_TEAMPOS_ITS_ALREADY_IN_ACTION")), 0
		}
	}

	heroOld := self.FightPos[index]

	self.FightPos[index] = heroKeyId
	return nil, heroOld
}

// 阵容为空
func (self *TeamPos) isUIPosEmpty() bool {
	for _, heroId := range self.FightPos {
		if heroId != 0 {
			return false
		}
	}
	return true
}

// 检查英雄是否上阵
func (self *TeamPos) isEmbattle(heroId int) bool {
	for i := 0; i < MAX_UI_POS; i++ {
		if self.FightPos[i] == heroId {
			return true
		}
	}
	return false
}

// 根据上阵英雄返回位置
func (self *TeamPos) getHeroId(index int) int {
	if index < 0 || index >= MAX_FIGHT_POS {
		return 0
	}
	return self.FightPos[index]
}
