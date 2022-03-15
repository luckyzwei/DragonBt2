package game

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	//"time"
)

const (
	UPSTAR_TYPE_SELF     = 1 //本体
	UPSTAR_TYPE_COLOR    = 2 //同阵营
	UPSTAR_TYPE_ALL      = 3 //任意阵营
	UPSTAR_TYPE_SPECTIAL = 4 //指定材料
)

const (
	COMPOUND_NUM                = 60
	BACKOPEN_LIMIT_NUM          = 2
	BACKOPEN_LIMIT_STAR         = 5
	BACKOPEN_LIMIT_SPECIAL_STAR = 4
	HIRE_STAR_LIMIT             = 6
)

const (
	HERO_USE_TYPE_SUPPORT           = 0 // 支援英雄
	HERO_USE_TYPE_REWARD            = 1 // 个人悬赏
	HERO_USE_TYPE_CRYSTAL_PRIESTS   = 2 // 共鸣水晶祭司
	HERO_USE_TYPE_CRYSTAL_RESONANCE = 3 // 共鸣水晶普通英雄
	HERO_USE_TYPE_REWARD_TEAM       = 4 // 团队悬赏
	HERO_USE_TYPE_MAX               = 5
)

const (
	HERO_ATTRIBUTE_WATER   = 1 // 水系英雄
	HERO_ATTRIBUTE_FIRE    = 2 // 火系英雄
	HERO_ATTRIBUTE_THUNDER = 3 // 雷系英雄
	HERO_ATTRIBUTE_EARTH   = 4 // 土系英雄
	HERO_ATTRIBUTE_LIGHT   = 5 // 光系英雄
	HERO_ATTRIBUTE_DARK    = 6 // 暗系英雄
	HERO_ATTRIBUTE_AIR     = 7 // 虚空英雄
)

// 英雄相关处理
const (
	activateStarAction  = "activatestar"
	upStarAction        = "upstar"
	upStarAutoAction    = "upstarauto"
	heroSynthesis       = "synthesis"
	interchange         = "interchange"
	checkteam           = "checkteam"
	upgradeTalentAction = "upgradedivinity"
	activatefate        = "activatefate"
	herodecompose       = "herodecompose"
	reBornAction        = "reborn"
	talentReset         = "divinityreset"
	compose             = "compose"
)

const (
	MSG_HERO_STAGE_TALENT_SET_SKILL = "msg_hero_stage_talent_set_skill"
)

const (
	HERO_LV_LIMIT = 0 //英雄等级不能超过玩家等级的偏移值
)

const (
	COMPOSE_TYPE_RUNE  = 1
	COMPOSE_TYPE_EQUIP = 2
)

const (
	HERO_STAR_TOTAL       = 0 // 武将总天赋
	HERO_STAR_TOTAL_CAMP1 = 1 // 武将总天赋1
	HERO_STAR_TOTAL_CAMP2 = 2 // 武将总天赋2
	HERO_STAR_TOTAL_CAMP3 = 3 // 武将总天赋3
	HERO_STAR_TOTAL_CAMP4 = 4 // 武将总天赋4
	HERO_STAR_TOTAL_MAX   = 5 // 最大值
)

type HeroTopInfo struct {
	Stars    int   // 总星级
	StarTime int64 // 达到总星级的时间
}

type HandBookInfo struct {
	State   int `json:"state"`   //
	StarMax int `json:"starmax"` // 达到的最大星级
}

type San_Hero struct {
	Uid               int64  // 角色ID
	Info              string // 武将数据
	Totalstars        string // 武将总天赋
	StarTime          int64  // 达到总天赋的时间
	Reborn            int    // 重生次数
	HeroTotalStars    int    // 英雄总星级
	HeroStarTime      int64  // 达到英雄总星级的时间
	MaxKey            int    // 英雄最大KEY
	BuyPosNum         int    // 英雄栏扩展次数
	AutoFire          int    // 是否自动分解
	BackOpen          int    // 回撤功能是否开放
	HandBook          string // 图鉴激活情况
	CompoundSignNum   string // 合成标记次数组
	CompoundSignScore string // 合成标记分值组

	info              map[int]*Hero // 英雄数据
	totalStars        []*HeroTopInfo
	handBook          map[int]*HandBookInfo
	compoundSignNum   map[int]int
	compoundSignScore map[int]int
	DataUpdate
}

// 新武将数据库
type NewHero struct {
	Uid              int64                  `json:"uid"`       // 玩家Id
	HeroId           int                    `json:"heroid"`    // 英雄Id
	Fight            int64                  `json:"fight"`     // 保留
	StarItem         *StarItem              `json:"staritem"`  // 升星
	HeroKeyId        int                    `json:"herokeyid"` // 英雄KeyId
	IsLock           int                    `json:"islock"`    // 上锁状态
	EquipIds         []*Equip               `json:"equipIds"`
	ArtifactEquipIds []*ArtifactEquip       `json:"artifactequipids"`
	HeroLv           int                    `json:"herolv"`         // 英雄等级
	Skin             int                    `json:"skin"`           // 皮肤
	UseType          [HERO_USE_TYPE_MAX]int `json:"usetype"`        // 使用类型
	ExclusiveEquip   *ExclusiveEquip        `json:"exclusiveequip"` // 专属装备
	Attr             map[int]int64          `json:"attr"`           // 部分模块的展示属性
	StageTalent      *StageTalent           `json:"talent"`         // 新天赋系统 和老天赋分开
}

// 武将数据库
type Hero struct {
	Uid              int64                  `json:"uid"`       // 玩家Id
	HeroId           int                    `json:"heroid"`    // 英雄Id
	Fight            int64                  `json:"fight"`     // 保留
	Horse            int                    `json:"horse"`     // 英雄上阵的战马信息
	StarItem         *StarItem              `json:"staritem"`  // 升星
	HeroKeyId        int                    `json:"herokeyid"` // 英雄KeyId
	IsLock           int                    `json:"islock"`    // 上锁状态
	EquipIds         []int                  `json:"equipIds"`
	ArtifactEquipIds []int                  `json:"artifactequipids"`
	HeroLv           int                    `json:"herolv"`         // 英雄等级
	UseType          [HERO_USE_TYPE_MAX]int `json:"usetype"`        // 使用类型
	ExclusiveEquip   *ExclusiveEquip        `json:"exclusiveequip"` // 专属装备
	OriginalLevel    int                    `json:"originallevel"`  // 原本等级
	Skin             int                    `json:"skin"`           // 设置皮肤
	VoidHero         int                    `json:"voidhero"`       // 是否是虚空英雄
	Resonance        int                    `json:"resonance"`      // 共鸣对象
	StageTalent      *StageTalent           `json:"talent"`         // 新天赋系统 和老天赋分开
	//Tiger    int   `json:"tiger"` //20190907 一下从阵位移过来的
	//ArmyId   int   `json:"army_id"`
	//FlagIds  []int `json:"flag_ids"`
	//TalentItem *TalentItem `json:"divinity"`  // 天赋
	//FateItem   *FateItem   `json:"fateitem"`  // 缘分
	//RuneItem   *RuneItem   `json:"runeitem"`  // 符文系统
}

func (h *Hero) GetStar() int {
	if h.StarItem == nil {
		return 0
	}
	return h.StarItem.UpStar
}

type LstHeroStar []*Hero

func (a LstHeroStar) Len() int {
	return len(a)
}

func (a LstHeroStar) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a LstHeroStar) Less(i, j int) bool {
	return a[j].StarItem.UpStar < a[i].StarItem.UpStar
}

type LstHero []*Hero

func (a LstHero) Len() int {
	return len(a)
}

func (a LstHero) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a LstHero) Less(i, j int) bool {
	return a[j].Fight < a[i].Fight
}

type LstHeroLevel []*Hero

func (a LstHeroLevel) Len() int {
	return len(a)
}

func (a LstHeroLevel) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a LstHeroLevel) Less(i, j int) bool {
	if a[i].HeroLv > a[j].HeroLv {
		return true
	}

	if a[i].HeroLv < a[j].HeroLv {
		return false
	}

	if a[i].Fight > a[j].Fight {
		return true
	}

	if a[i].Fight < a[j].Fight {
		return false
	}

	return false
}

// 英雄
type ModHero struct {
	player   *Player
	Sql_Hero San_Hero
}

func (self *ModHero) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_userhero2` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Hero, "san_userhero2", self.player.ID)

	if self.Sql_Hero.Uid <= 0 {
		self.Sql_Hero.Uid = self.player.ID
		self.Sql_Hero.info = make(map[int]*Hero)
		self.Sql_Hero.handBook = make(map[int]*HandBookInfo)
		self.Sql_Hero.compoundSignNum = make(map[int]int)
		self.Sql_Hero.compoundSignScore = make(map[int]int)
		self.Sql_Hero.AutoFire = LOGIC_TRUE
		self.Encode()
		InsertTable("san_userhero2", &self.Sql_Hero, 0, true)
		self.Sql_Hero.Init("san_userhero2", &self.Sql_Hero, true)
	} else {
		self.Decode()
		self.Sql_Hero.Init("san_userhero2", &self.Sql_Hero, true)
		self.checkHero()
	}

	if self.Sql_Hero.compoundSignNum == nil {
		self.Sql_Hero.compoundSignNum = make(map[int]int)
	}

	if self.Sql_Hero.compoundSignScore == nil {
		self.Sql_Hero.compoundSignScore = make(map[int]int)
	}

	nLen := len(self.Sql_Hero.totalStars)
	if nLen < HERO_STAR_TOTAL_MAX {
		for i := 0; i < HERO_STAR_TOTAL_MAX-nLen; i++ {
			self.Sql_Hero.totalStars = append(self.Sql_Hero.totalStars, &HeroTopInfo{})
		}
	}
}

func (self *ModHero) OnGetOtherData() {
	self.GetAllHeroStars()
}

func (self *ModHero) Decode() {
	json.Unmarshal([]byte(self.Sql_Hero.Info), &self.Sql_Hero.info)
	json.Unmarshal([]byte(self.Sql_Hero.Totalstars), &self.Sql_Hero.totalStars)
	json.Unmarshal([]byte(self.Sql_Hero.HandBook), &self.Sql_Hero.handBook)
	json.Unmarshal([]byte(self.Sql_Hero.CompoundSignNum), &self.Sql_Hero.compoundSignNum)
	json.Unmarshal([]byte(self.Sql_Hero.CompoundSignScore), &self.Sql_Hero.compoundSignScore)
}

func (self *ModHero) Encode() {
	self.Sql_Hero.Info = HF_JtoA(&self.Sql_Hero.info)
	self.Sql_Hero.Totalstars = HF_JtoA(&self.Sql_Hero.totalStars)
	self.Sql_Hero.HandBook = HF_JtoA(&self.Sql_Hero.handBook)
	self.Sql_Hero.CompoundSignNum = HF_JtoA(&self.Sql_Hero.compoundSignNum)
	self.Sql_Hero.CompoundSignScore = HF_JtoA(&self.Sql_Hero.compoundSignScore)
}

// 注册消息
func (self *ModHero) onReg(handlers map[string]func(body []byte)) {
	handlers["herouplv"] = self.HeroUpLv
	handlers["heroupstar"] = self.HeroUpStar
	handlers["heroupstarall"] = self.HeroUpStarAll
	handlers["herolock"] = self.HeroLock
	handlers["herobuypos"] = self.HeroBuyPos
	handlers["heroreborn"] = self.HeroReborn
	handlers["herofire"] = self.HeroFire
	handlers["heroback"] = self.HeroBack
	handlers["herocompound"] = self.HeroCompound
	handlers["heroautofire"] = self.HeroAutoFire
	handlers["herogethandbook"] = self.HeroGetHandBook
	handlers["herouplvto"] = self.HeroUpLvTo
	handlers[MSG_HERO_SET_VOID_HERO_RESONANCE] = self.MsgSetVoidHeroResonance
	handlers[MSG_HERO_CANCEL_VOID_HERO_RESONANCE] = self.MsgCancelVoidHeroResonance
	handlers[MSG_HERO_STAGE_TALENT_SET_SKILL] = self.MsgStageTalentSetSkill
}

// 英雄合成
func (self *ModHero) HeroCompound(body []byte) {

	var msg C2S_HeroCompound
	json.Unmarshal(body, &msg)

	itemConfig := GetCsvMgr().ItemMap[msg.ItemId]
	if itemConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ITEM_CONFIG_ERROR"))
		return
	}

	if itemConfig.CompoundNum <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ITEM_CONFIG_ERROR"))
		return
	}

	if msg.Num <= 0 || msg.Num%itemConfig.CompoundNum != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_MSG_NUM_ERROR"))
		return
	}
	count := msg.Num / itemConfig.CompoundNum
	//看看数量问题
	if !self.CheckHeroBuyPos() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_NUM_MAX"))
		return
	}

	err := self.player.HasObjectOkEasy(msg.ItemId, msg.Num)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	//开始合成
	var msgRel S2C_HeroCompound
	msgRel.Cid = "herocompound"
	msgRel.Costs = self.player.RemoveObjectSimple(msg.ItemId, msg.Num, "合成道具", msg.ItemId, msg.Num, 0)

	for i := 0; i < (count); i++ {
		if itemConfig.CompoundId == 0 {
			//紫色魂石附加逻辑
			if itemConfig.ItemId == ITEM_SOUL_HERO_STONE {
				items := self.CalSoulHeroStone()
				if len(items) == 0 {
					continue
				}
				for _, v := range items {
					item := self.player.AddObjectSimple(v.ItemId, v.ItemNum, "合成道具", 0, 0, 0)
					msgRel.GetHero = append(msgRel.GetHero, item...)
				}
			} else {
				items := GetLootMgr().LootItem(itemConfig.LotteryId, self.player)
				if len(items) == 0 {
					continue
				}
				for _, v := range items {
					item := self.player.AddObjectSimple(v.ItemId, v.ItemNum, "合成道具", 0, 0, 0)
					msgRel.GetHero = append(msgRel.GetHero, item...)
				}
			}
		} else {
			item := self.player.AddObjectSimple(itemConfig.CompoundId, 1, "合成道具", 0, 0, 0)
			msgRel.GetHero = append(msgRel.GetHero, item...)
		}
	}

	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BAG_COMPOUND, msg.ItemId, msg.Num, 0, "合成道具", 0, 0, self.player)
}

//
func (self *ModHero) CalSoulHeroStone() map[int]*Item {
	data := make(map[int]*Item)
	//先看保底组
	itemId := self.FindSpecial(ITEM_SOUL_HERO_STONE, self.Sql_Hero.compoundSignNum[ITEM_SOUL_HERO_STONE]+1)
	if itemId > 0 {
		self.Sql_Hero.compoundSignNum[ITEM_SOUL_HERO_STONE]++
		AddItemMapHelper3(data, itemId, 1)
		return data
	}
	//保底组没出现走权值计算
	item := self.CalSyntheticDrop()
	if item.ItemID > 0 {
		AddItemMapHelper3(data, item.ItemID, item.Num)
	}
	return data
}

func (self *ModHero) CalSyntheticDrop() PassItem {
	var item PassItem
	//增加自己的权值分
	self.Sql_Hero.compoundSignScore[ITEM_SOUL_HERO_STONE] += ASTROLOGY_ADD

	//根据权值分数计算大组的总权值
	groupRateAll := 0
	config := GetCsvMgr().SyntheticDropGroupConfig
	for _, v := range config {
		if v.ScoreLimit != 0 && self.Sql_Hero.compoundSignScore[ITEM_SOUL_HERO_STONE] >= v.ScoreLimit {
			groupRateAll += v.SyntheticMid
		} else {
			groupRateAll += v.SyntheticChance
		}
	}
	if groupRateAll <= 0 {
		return item
	}
	//获得随机数值
	groupRand := self.player.GetModule("find").(*ModFind).GetRandInt(groupRateAll)
	groupRateCal := 0
	//计算选中组
	for _, v := range config {
		if v.ScoreLimit != 0 && self.Sql_Hero.compoundSignScore[ITEM_SOUL_HERO_STONE] >= v.ScoreLimit {
			groupRateCal += v.SyntheticMid
		} else {
			groupRateCal += v.SyntheticChance
		}
		//选中
		if groupRateCal > groupRand {
			//如果是目标掉落组
			info, ok := GetCsvMgr().SyntheticDropConfig[v.SyntheticId]
			if ok {
				rateAll := 0
				for _, vv := range info {
					rateAll += vv.ItemWT
				}
				if rateAll <= 0 {
					return item
				}
				//获得随机数值
				rand := self.player.GetModule("find").(*ModFind).GetRandInt(rateAll)
				rateCal := 0
				for _, vv := range info {
					rateCal += vv.ItemWT
					if rateCal > rand {
						item.Num = vv.ItemNum
						item.ItemID = vv.ItemId
						if vv.ScoreValue == LOGIC_FALSE {
							self.Sql_Hero.compoundSignScore[ITEM_SOUL_HERO_STONE] += vv.ItemScore
						} else {
							self.Sql_Hero.compoundSignScore[ITEM_SOUL_HERO_STONE] -= vv.ItemScore
						}
						return item
					}
				}
			}
		}
	}

	return item
}

func (self *ModHero) GetSyntheticScore() int {
	return self.Sql_Hero.compoundSignScore[ITEM_SOUL_HERO_STONE]
}

func (self *ModHero) FindSpecial(itemId int, num int) int {
	lst := GetCsvMgr().PubChestSpecialLst[itemId]
	for i := 0; i < len(lst); i++ {
		//合成类型默认为3
		if lst[i].SpType != 3 {
			continue
		}
		if num >= lst[i].Droptimemin && num <= lst[i].Droptimemax {
			specail := lst[i].DropGroupModify
			items := strings.Split(specail, "|")
			if len(items) > 0 {
				groups := make([]DropGroupModify, 0)
				items := strings.Split(specail, "|")
				if len(items) > 0 {
					for j := 0; j < len(items); j++ {
						item := strings.Split(items[j], ":")
						var group DropGroupModify
						group.Original = HF_Atoi(item[0])
						group.Rate = HF_Atoi(item[1])
						groups = append(groups, group)
					}
				}

				allRate := 0
				for j := 0; j < len(groups); j++ {
					allRate += groups[j].Rate
				}

				value := self.player.GetModule("find").(*ModFind).GetRandInt(allRate)

				nowRate := 0
				for j := 0; j < len(groups); j++ {
					nowRate += groups[j].Rate
					if nowRate > value {
						return groups[j].Original
					}
				}
			}
			break
		}
	}

	return 0
}

//自动分解普通英雄
func (self *ModHero) HeroAutoFire(body []byte) {

	var msg C2S_HeroAutoFire
	json.Unmarshal(body, &msg)

	if msg.Action != LOGIC_FALSE && msg.Action != LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_MSG_ACTION_ERROR"))
		return
	}

	//不验证之前的状态，直接更新并同步最新状态
	self.Sql_Hero.AutoFire = msg.Action

	var msgRel S2C_HeroAutoFire
	msgRel.Cid = "heroautofire"
	msgRel.AutoFire = self.Sql_Hero.AutoFire
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModHero) HeroGetHandBook(body []byte) {

	var msg C2S_HeroGetHandBook
	json.Unmarshal(body, &msg)

	_, ok := self.Sql_Hero.handBook[msg.HeroId]

	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_CONFIGURATION_ERROR"))
		return
	}

	config := GetCsvMgr().HeroHandBookConfigMap[msg.HeroId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_CONFIGURATION_ERROR"))
		return
	}

	if self.Sql_Hero.handBook[msg.HeroId].State != CANTAKE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REWARD_CANT_GET"))
		return
	}

	item := self.player.AddObjectSimple(config.Prize, config.PrizeNum, "领取英雄图鉴奖励", msg.HeroId, 0, 0)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HANDBOOK_AWARD, msg.HeroId, 0, 0, "领取图鉴奖励", 0, 0, self.player)
	self.Sql_Hero.handBook[msg.HeroId].State = TAKEN

	var msgRel S2C_HeroGetHandBook
	msgRel.Cid = "herogethandbook"
	msgRel.HeroId = msg.HeroId
	msgRel.State = self.Sql_Hero.handBook[msg.HeroId].State
	msgRel.StarMax = self.Sql_Hero.handBook[msg.HeroId].StarMax
	msgRel.GetItems = item
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModHero) CheckItem(id int, num int) ([]int, []int) {

	flag, _ := GetCsvMgr().IsLevelOpen(self.player.GetLv(), OPEN_LEVEL_HERO_REBORN)
	if !flag {
		return nil, nil
	}

	if self.Sql_Hero.AutoFire == LOGIC_FALSE {
		return nil, nil
	}

	heroId := (id - 11000000) / 100
	initLv := GetCsvMgr().GetHeroInitLv(heroId)
	if initLv != LOGIC_TRUE {
		return nil, nil
	}
	config := GetCsvMgr().GetHeroMapConfig(heroId, initLv)
	if config == nil {
		return nil, nil
	}
	//转换成功，需要开启图鉴
	//处理图鉴
	_, ok := self.Sql_Hero.handBook[heroId]
	if !ok {
		self.Sql_Hero.handBook[heroId] = new(HandBookInfo)
		self.Sql_Hero.handBook[heroId].State = CANTAKE
		self.Sql_Hero.handBook[heroId].StarMax = config.HeroStar
	}
	return config.DisbandId, config.DisbandNum
}

func (self *ModHero) HeroReborn(body []byte) {
	flag, configLv := GetCsvMgr().IsLevelOpen(self.player.GetLv(), OPEN_LEVEL_HERO_REBORN)
	if !flag {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN"), configLv))
		return
	}

	var msg C2S_HeroReborn
	json.Unmarshal(body, &msg)

	hero := self.GetHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_HERO_NOT_EXIST"))
		return
	}

	config := GetCsvMgr().HeroExpConfigMap[hero.HeroLv]
	if self.player.GetModule("crystal").(*ModResonanceCrystal).IsResonanceHero(msg.HeroKeyId) {
		config = GetCsvMgr().HeroExpConfigMap[hero.OriginalLevel]
	}
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	//检查消耗够不够
	if err := self.player.HasObjectOkEasy(config.ResetCostId, config.ResetCostNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	oldlevel := hero.HeroLv
	if self.player.GetModule("crystal").(*ModResonanceCrystal).IsResonanceHero(msg.HeroKeyId) {
		hero.OriginalLevel = 1
	} else {
		hero.LvUp(0 - hero.HeroLv)
	}
	equipRel, artifactRel, horse := self.OffEquipAll(hero)

	costItem := self.player.RemoveObjectSimple(config.ResetCostId, config.ResetCostNums, "英雄重塑", hero.HeroId, 0, 0)
	getItem := self.player.AddObjectLst(config.ResetItems, config.ResetNums, "英雄重塑", hero.HeroId, 0, 0)

	self.player.countHeroFight(hero, ReasonHeroReborn)
	self.player.GetModule("support").(*ModSupportHero).UpdataHero(hero.HeroKeyId)
	self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(hero.HeroKeyId)
	self.ChangeVoidHeroResonance(hero.HeroKeyId, 0)
	self.GetAllHeroStars()

	var msgRel S2C_HeroReborn
	msgRel.Cid = "heroreborn"
	msgRel.HeroKeyId = msg.HeroKeyId
	msgRel.StarItem = hero.StarItem
	msgRel.CostItem = costItem
	msgRel.GetItem = getItem
	msgRel.OffEquip = equipRel
	msgRel.OffArtifact = artifactRel
	msgRel.OffHorseInfo = horse
	msgRel.HeroLv = hero.HeroLv
	msgRel.OriginalLevel = hero.OriginalLevel
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_REBORN, hero.HeroId, oldlevel, hero.HeroLv, "英雄重塑", 0, hero.HeroKeyId, self.player)
}

func (self *ModHero) HeroBack(body []byte) {
	var msg C2S_HeroBack
	json.Unmarshal(body, &msg)

	hero := self.GetHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_HERO_NOT_EXIST"))
		return
	}
	//看功能开了没有
	if self.Sql_Hero.BackOpen == LOGIC_FALSE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_BACK_STAR_NOT_OPEN"))
		return
	}

	//初始必须大于4星英雄
	if GetCsvMgr().GetHeroInitLv(hero.HeroId) < BACKOPEN_LIMIT_SPECIAL_STAR {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_BACK_THIS_HERO_CANT"))
		return
	}

	if hero.StarItem.UpStar <= BACKOPEN_LIMIT_STAR {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_BACK_STAR_NOT_ENOUGH"))
		return
	}

	config := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	if GetCsvMgr().GetHeroInitLv(hero.HeroId) < 4 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_MSG_ERROR"))
		return
	}

	configCost := GetCsvMgr().GetTariffConfig2(TARIFF_TYPE_HERO_BACK)
	if configCost == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	//检查消耗够不够
	if err := self.player.HasObjectOk(configCost.ItemIds, configCost.ItemNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	//处理等级
	configLv := GetCsvMgr().HeroExpConfigMap[hero.HeroLv]
	if self.player.GetModule("crystal").(*ModResonanceCrystal).IsResonanceHero(msg.HeroKeyId) {
		configLv = GetCsvMgr().HeroExpConfigMap[hero.OriginalLevel]
	}
	if configLv == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	// 扣除物品
	costItem := self.player.RemoveObjectLst(configCost.ItemIds, configCost.ItemNums, "英雄降阶", hero.HeroId, 0, 0)

	//返还材料 从多行改成单行读取 20200108
	itemMap := make(map[int]*Item)
	configHero := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if configHero != nil {
		AddItemMapHelper(itemMap, configHero.BackItem, configHero.BackNum)
	}
	AddItemMapHelper(itemMap, configLv.ResetItems, configLv.ResetNums)

	equipRel, artifactRel, horse := self.OffEquipAll(hero)

	oldstar := hero.StarItem.UpStar
	hero.StarItem.UpStar = BACKOPEN_LIMIT_STAR

	if self.player.GetModule("crystal").(*ModResonanceCrystal).IsResonanceHero(msg.HeroKeyId) {
		hero.OriginalLevel = 1
	} else {
		hero.LvUp(0 - hero.HeroLv)
	}

	//处理专属装备
	if hero.ExclusiveEquip.UnLock == LOGIC_TRUE {
		oldlevel := hero.ExclusiveEquip.Lv
		exclusiveActiveConfig := GetCsvMgr().ExclusiveEquipConfigMap[hero.ExclusiveEquip.Id]
		if exclusiveActiveConfig != nil {
			AddItemMapHelper3(itemMap, exclusiveActiveConfig.ActiveNeed, exclusiveActiveConfig.ActiveNum)
		}
		if hero.ExclusiveEquip.Lv > INIT_LV {
			exclusiveUpConfig := GetCsvMgr().ExclusiveStrengthen[hero.ExclusiveEquip.Id][hero.ExclusiveEquip.Lv]
			if exclusiveUpConfig != nil {
				AddItemMapHelper(itemMap, exclusiveUpConfig.BackItem, exclusiveUpConfig.BackNum)
			}
		}
		hero.ExclusiveEquip = hero.NewExclusiveEquipItem()
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_EXCLUSIVE_RESET, hero.ExclusiveEquip.Id, hero.HeroId, oldlevel, "英雄降阶专属回退", 0, hero.ExclusiveEquip.Lv, self.player)
	}

	// 检测天赋
	hero.CheckStageTalent()

	getItem := self.player.AddObjectItemMap(itemMap, "英雄降阶", hero.HeroId, 0, 0)
	self.player.countHeroFight(hero, ReasonHeroBack)

	var msgRel S2C_HeroBack
	msgRel.Cid = "heroback"
	msgRel.HeroKeyId = msg.HeroKeyId
	msgRel.StarItem = hero.StarItem
	msgRel.CostItem = costItem
	msgRel.GetItem = getItem
	msgRel.OffEquip = equipRel
	msgRel.OffArtifact = artifactRel
	msgRel.OffHorseInfo = horse
	msgRel.HeroLv = hero.HeroLv
	msgRel.ExclusiveEquip = hero.ExclusiveEquip
	msgRel.StageTalent = hero.StageTalent
	msgRel.OriginalLevel = hero.OriginalLevel
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_BACK, hero.HeroId, oldstar, hero.StarItem.UpStar, "英雄降阶", 0, hero.HeroKeyId, self.player)

}

func (self *ModHero) HeroFire(body []byte) {
	var msg C2S_HeroFire
	json.Unmarshal(body, &msg)
	itemMap := make(map[int]*Item)
	var msgRel S2C_HeroFire
	msgRel.Cid = "herofire"
	for i := 0; i < len(msg.HeroKeyId); i++ {
		hero := self.GetHero(msg.HeroKeyId[i])
		if hero == nil {
			continue
		}
		config := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
		if config == nil {
			continue
		}
		if GetCsvMgr().GetHeroInitLv(hero.HeroId) != 1 {
			continue
		}

		configExp := GetCsvMgr().HeroExpConfigMap[hero.HeroLv]
		if self.player.GetModule("crystal").(*ModResonanceCrystal).IsResonanceHero(msg.HeroKeyId[i]) {
			configExp = GetCsvMgr().HeroExpConfigMap[hero.OriginalLevel]
		}
		if configExp == nil {
			continue
		}
		AddItemMapHelper(itemMap, config.DisbandId, config.DisbandNum)
		AddItemMapHelper(itemMap, configExp.ResetItems, configExp.ResetNums)
		offequip, offartifact, horse := self.OffEquipAll(hero)
		msgRel.OffEquip = append(msgRel.OffEquip, offequip...)
		msgRel.OffArtifact = append(msgRel.OffArtifact, offartifact...)
		msgRel.OffHorseInfo = append(msgRel.OffHorseInfo, horse...)
		msgRel.HeroKeyId = append(msg.HeroKeyId, msg.HeroKeyId[i])
	}
	msgRel.GetItem = self.player.AddObjectItemMap(itemMap, "英雄遣退", 0, 0, 0)
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	for i := 0; i < len(msg.HeroKeyId); i++ {
		hero := self.GetHero(msg.HeroKeyId[i])
		if hero == nil {
			continue
		}
		config := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
		if config == nil {
			continue
		}
		if GetCsvMgr().GetHeroInitLv(hero.HeroId) != 1 {
			continue
		}

		configExp := GetCsvMgr().HeroExpConfigMap[hero.HeroLv]
		if self.player.GetModule("crystal").(*ModResonanceCrystal).IsResonanceHero(msg.HeroKeyId[i]) {
			configExp = GetCsvMgr().HeroExpConfigMap[hero.OriginalLevel]
		}
		if configExp == nil {
			continue
		}
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_DELETE, hero.HeroId, hero.StarItem.UpStar, 0, "消耗英雄", 1, 0, self.player)
		self.DeleteHero(msg.HeroKeyId[i])
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIRE, len(msg.HeroKeyId), 0, 0, "英雄遣退", 0, 0, self.player)
}

//英雄升级
func (self *ModHero) HeroUpLv(body []byte) {

	var msg C2S_HeroLvUp
	json.Unmarshal(body, &msg)

	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THERE_IS_NO_HERO"))
		return
	}
	// 共鸣中的英雄无法升级
	if hero.VoidHero != 0 && hero.Resonance != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THERE_IS_NO_HERO"))
		return
	}

	if !self.player.GetModule("crystal").(*ModResonanceCrystal).CanUpLevel(msg.HeroKeyId) {
		return
	}

	//看看阶级够不够
	configHero := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if configHero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_GRADE_CONFIG_ERROR"))
		return
	}
	//最终最大等级
	if hero.HeroLv >= configHero.FinalLevel {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_LEVEL_MAX"))
		return
	}
	//当前最大等级
	if hero.HeroLv >= configHero.HeroLvMax {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_LEVEL_MAX"))
		return
	}
	configUpLv := GetCsvMgr().HeroExpConfigMap[hero.HeroLv+1]
	if configUpLv == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_LEVEL_CONFIG_ERROR"))
		return
	}

	// 检查消耗是否正常
	err := self.player.HasObjectOk(configUpLv.CostItems, configUpLv.CostNums)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	oldlevel := hero.HeroLv
	//计算突破ID
	hero.LvUp(1)

	items := self.player.RemoveObjectLst(configUpLv.CostItems, configUpLv.CostNums, "英雄升级", hero.HeroLv, hero.HeroId, 0)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_LEVE_UP, hero.HeroId, oldlevel, hero.HeroLv, "英雄升级", 0, hero.HeroId, self.player)

	self.player.GetModule("support").(*ModSupportHero).UpdataHero(hero.HeroKeyId)
	self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(hero.HeroKeyId)
	self.ChangeVoidHeroResonance(hero.HeroKeyId, 0)
	self.GetAllHeroStars()

	self.player.HandleTask(TASK_TYPE_BIGGEST_LEVEL, 0, 0, 0)
	self.player.HandleTask(TASK_TYPE_LEVEL_UP_COUNT, 1, 0, 0)

	self.player.countHeroFight(hero, ReasonHeroLvUp)

	var msgRel S2C_HeroLvUp
	msgRel.Cid = "herouplv"
	msgRel.HeroKeyId = hero.HeroKeyId
	msgRel.HeroLv = hero.HeroLv
	msgRel.HeroBreakId = hero.StarItem.HeroBreakId
	msgRel.Items = items
	msgRel.StarItem = hero.StarItem
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModHero) HeroUpLvTo(body []byte) {

	var msg C2S_HeroLvUpTo
	json.Unmarshal(body, &msg)

	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THERE_IS_NO_HERO"))
		return
	}
	// 共鸣中的英雄无法升级
	if hero.VoidHero != 0 && hero.Resonance != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THERE_IS_NO_HERO"))
		return
	}

	if !self.player.GetModule("crystal").(*ModResonanceCrystal).CanUpLevelByAim(msg.HeroKeyId, msg.AimLv) {
		return
	}

	//看看阶级够不够
	configHero := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if configHero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_GRADE_CONFIG_ERROR"))
		return
	}
	//最终最大等级
	if hero.HeroLv >= configHero.FinalLevel || msg.AimLv > configHero.FinalLevel || msg.AimLv <= hero.HeroLv {
		//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_LEVEL_MAX"))
		return
	}
	//当前最大等级
	if hero.HeroLv >= configHero.HeroLvMax || msg.AimLv > configHero.HeroLvMax {
		//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_LEVEL_MAX"))
		return
	}

	config := GetCsvMgr().HeroBreakConfigMap[hero.HeroId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THERE_IS_NO_HERO"))
		return
	}

	//计算消耗
	costAll := make(map[int]*Item)
	for i := hero.HeroLv + 1; i <= msg.AimLv; i++ {
		configUpLv := GetCsvMgr().HeroExpConfigMap[i]
		if configUpLv == nil {
			continue
		}
		AddItemMapHelper(costAll, configUpLv.CostItems, configUpLv.CostNums)
	}

	// 检查消耗是否正常
	err := self.player.HasObjectMapItemOk(costAll)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	oldlevel := hero.HeroLv
	//计算突破ID
	hero.LvUp(msg.AimLv - hero.HeroLv)

	items := self.player.RemoveObjectItemMap(costAll, "英雄一键升级", hero.HeroLv, hero.HeroId, 0)
	//计算突破等级
	for _, v := range config {
		if hero.HeroLv >= v.Break && hero.StarItem.HeroBreakId < v.Id {
			hero.StarItem.HeroBreakId = v.Id
		}
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_LEVE_UP_TO, hero.HeroId, oldlevel, hero.HeroLv, "英雄一键升级", 0, hero.HeroKeyId, self.player)

	self.player.HandleTask(TASK_TYPE_BIGGEST_LEVEL, 0, 0, 0)
	self.player.HandleTask(TASK_TYPE_LEVEL_UP_COUNT, 1, 0, 0)

	self.player.countHeroFight(hero, ReasonHeroLvUp)

	self.player.GetModule("support").(*ModSupportHero).UpdataHero(hero.HeroKeyId)
	self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(hero.HeroKeyId)
	self.ChangeVoidHeroResonance(hero.HeroKeyId, 0)
	self.GetAllHeroStars()

	var msgRel S2C_HeroLvUpTo
	msgRel.Cid = "herouplvto"
	msgRel.HeroKeyId = hero.HeroKeyId
	msgRel.HeroLv = hero.HeroLv
	msgRel.HeroBreakId = hero.StarItem.HeroBreakId
	msgRel.StarItem = hero.StarItem
	msgRel.Items = items
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

//! 英雄升星
func (self *ModHero) HeroUpStar(body []byte) {

	var msg C2S_HeroUpStar
	json.Unmarshal(body, &msg)

	//消耗英雄是否存在
	for i := 0; i < len(msg.CostHeroKeyId); i++ {
		costHero := self.player.getHero(msg.CostHeroKeyId[i])
		if costHero == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
			self.SendInfoSyn()
			return
		}

		if costHero.IsLock == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_IS_LOCK"))
			return
		}
	}

	//目标英雄是否存在
	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
		self.SendInfoSyn()
		return
	}

	// 虚空英雄 无法升阶
	if hero.VoidHero == 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
		return
	}

	//获取升级的配置
	heroConfig := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar+1)
	if heroConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_MAX_UPSTAR"))
		return
	}

	isUse := make([]bool, 0)
	costHero := make([]*Hero, 0)
	needCount := 0
	for i := 0; i < len(msg.CostHeroKeyId); i++ {
		isUse = append(isUse, false)
		costHero = append(costHero, self.player.getHero(msg.CostHeroKeyId[i]))
	}

	//
	for i := 0; i < len(heroConfig.UpStarType); i++ {
		count := heroConfig.UpStarNum[i]
		needCount += count
		switch heroConfig.UpStarType[i] {
		case UPSTAR_TYPE_SELF:
			for j := 0; j < len(msg.CostHeroKeyId); j++ {
				if count <= 0 {
					break
				}
				if isUse[j] {
					continue
				}
				if costHero[j].HeroId == hero.HeroId && costHero[j].StarItem.UpStar == heroConfig.UpStarStar[i] {
					isUse[j] = true
					count--
				}
			}
		case UPSTAR_TYPE_COLOR:
			for j := 0; j < len(msg.CostHeroKeyId); j++ {
				if count <= 0 {
					break
				}
				if isUse[j] {
					continue
				}
				if costHero[j].StarItem.UpStar == heroConfig.UpStarStar[i] {
					isUse[j] = true
					count--
				}
			}
		case UPSTAR_TYPE_ALL:
			for j := 0; j < len(msg.CostHeroKeyId); j++ {
				if count <= 0 {
					break
				}
				if isUse[j] {
					continue
				}
				costHeroConfig := GetCsvMgr().GetHeroMapConfig(costHero[j].HeroId, costHero[j].StarItem.UpStar)
				if costHeroConfig == nil {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
					self.SendInfoSyn()
					return
				}
				if costHero[j].StarItem.UpStar == heroConfig.UpStarStar[i] {
					isUse[j] = true
					count--
				}
			}
		case UPSTAR_TYPE_SPECTIAL:
			for j := 0; j < len(msg.CostHeroKeyId); j++ {
				if count <= 0 {
					break
				}
				if isUse[j] {
					continue
				}
				costHeroConfig := GetCsvMgr().HeroConfig[costHero[j].HeroId]
				if costHeroConfig == nil {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
					self.SendInfoSyn()
					return
				}
				if costHero[j].StarItem.UpStar == heroConfig.UpStarStar[i] {
					isUse[j] = true
					count--
				}
			}
		}
	}

	if needCount != len(msg.CostHeroKeyId) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_UPSTAR_COST_IS_ERROR"))
		return
	}

	for i := 0; i < len(isUse); i++ {
		if !isUse[i] {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_UPSTAR_COST_IS_ERROR"))
			return
		}
	}

	var msgRel S2C_HeroUpStar
	msgRel.Cid = "heroupstar"

	returnItem := make(map[int]*Item)
	//扣除英雄
	for i := 0; i < len(msg.CostHeroKeyId); i++ {
		costHero := self.player.getHero(msg.CostHeroKeyId[i])
		offEquip, artifact, horse := self.OffEquipAll(costHero)
		msgRel.OffEquip = append(msgRel.OffEquip, offEquip...)
		msgRel.OffArtifact = append(msgRel.OffArtifact, artifact...)
		msgRel.OffHorseInfo = append(msgRel.OffHorseInfo, horse...)

		configExp := GetCsvMgr().HeroExpConfigMap[costHero.HeroLv]
		if self.player.GetModule("crystal").(*ModResonanceCrystal).IsResonanceHero(msg.CostHeroKeyId[i]) {
			configExp = GetCsvMgr().HeroExpConfigMap[costHero.OriginalLevel]
		}
		if configExp == nil {
			continue
		}
		AddItemMapHelper(returnItem, configExp.ResetItems, configExp.ResetNums)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_DELETE, costHero.HeroId, costHero.StarItem.UpStar, hero.HeroId, "消耗英雄", 1, 0, self.player)
	}
	//处理返还
	if len(returnItem) > 0 {
		msgRel.GetItems = self.player.AddObjectItemMap(returnItem, "升阶材料返还", 0, 0, 0)
	}
	//处理升星
	oldstar := hero.StarItem.UpStar
	hero.StarItem.UpStar++
	self.player.HandleTask(TASK_TYPE_STAR_UP_COUNT, 1, hero.StarItem.UpStar, 0)
	self.player.HandleTask(TASK_TYPE_GET_HERO_WIDE, 1, hero.StarItem.UpStar, 0)
	self.player.HandleTask(TASK_TYPE_STAR_UP_COUNT_EQUAL, 1, hero.StarItem.UpStar, 0)
	//处理图鉴
	info, ok := self.Sql_Hero.handBook[hero.HeroId]
	if ok && info.StarMax < hero.StarItem.UpStar {
		info.StarMax = hero.StarItem.UpStar
	}
	// 检测天赋
	hero.CheckStageTalent()
	self.player.countHeroFight(hero, ReasonHeroStarUp)
	msgRel.HeroKeyId = msg.HeroKeyId
	msgRel.StarItem = hero.StarItem
	msgRel.StageTalent = hero.StageTalent
	msgRel.CostHeroKeyId = msg.CostHeroKeyId
	msgRel.Item = self.player.GetModule("lifetree").(*ModLifeTree).GetItemCount(hero.HeroId, hero.StarItem.UpStar)
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	//客户端要求后发这个消息
	for i := 0; i < len(msg.CostHeroKeyId); i++ {
		self.DeleteHero(msg.CostHeroKeyId[i])
	}

	//检查回退功能开启
	self.player.GetModule("support").(*ModSupportHero).UpdataHero(hero.HeroKeyId)
	self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(hero.HeroKeyId)
	self.player.GetModule("task").(*ModTask).SendUpdate()
	self.player.GetModule("entanglement").(*ModEntanglement).UpdateHero(hero.HeroKeyId)
	self.player.GetModule("entanglement").(*ModEntanglement).SendInfo([]byte{})
	self.player.GetModule("lifetree").(*ModLifeTree).StarUp(hero.HeroId, hero.StarItem.UpStar)
	self.player.GetModule("lifetree").(*ModLifeTree).SendRedPointInfo()
	self.ChangeVoidHeroResonance(hero.HeroKeyId, 0)
	self.GetAllHeroStars()
	GetOfflineInfoMgr().UpdateHeroSkin(self.player)

	self.CheckBackOpen(hero)
	self.CheckHireUpdate(hero)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_UP_STAR, hero.HeroId, oldstar, hero.StarItem.UpStar, "英雄升阶", 0, msg.HeroKeyId, self.player)
}

func (self *ModHero) HeroUpStarAll(body []byte) {

	var msg C2S_HeroUpStarAll
	json.Unmarshal(body, &msg)

	if len(msg.HeroKeyId) != len(msg.CostHeroKeyId) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_MSG_ERROR"))
		return
	}

	for _, v := range msg.CostHeroKeyId {
		//消耗英雄是否存在
		for i := 0; i < len(v); i++ {
			costHero := self.player.getHero(v[i])
			if costHero == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
				self.SendInfoSyn()
				return
			}

			if costHero.IsLock == LOGIC_TRUE {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_IS_LOCK"))
				return
			}
		}
	}

	var msgRel S2C_HeroUpStarAll
	msgRel.Cid = "heroupstarall"
	for size := 0; size < len(msg.HeroKeyId); size++ {
		//目标英雄是否存在
		hero := self.player.getHero(msg.HeroKeyId[size])
		if hero == nil {
			//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
			//return
			continue
		}

		// 虚空英雄 无法升阶
		if hero.VoidHero == 1 {
			continue
		}

		//获取升级的配置
		heroConfig := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar+1)
		if heroConfig == nil {
			//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_MAX_UPSTAR"))
			//return
			continue
		}

		isUse := make([]bool, 0)
		costHero := make([]*Hero, 0)
		for i := 0; i < len(msg.CostHeroKeyId[size]); i++ {
			isUse = append(isUse, false)
			costHero = append(costHero, self.player.getHero(msg.CostHeroKeyId[size][i]))

			tempcostHero := self.player.getHero(msg.CostHeroKeyId[size][i])
			if tempcostHero != nil {
				GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_DELETE, tempcostHero.HeroId, tempcostHero.StarItem.UpStar, hero.HeroId, "消耗英雄", 1, 0, self.player)
			}
		}

		//
		for i := 0; i < len(heroConfig.UpStarType); i++ {
			count := heroConfig.UpStarNum[i]
			switch heroConfig.UpStarType[i] {
			case UPSTAR_TYPE_SELF:
				for j := 0; j < len(msg.CostHeroKeyId[size]); j++ {
					if count <= 0 {
						break
					}
					if isUse[j] {
						continue
					}
					if costHero[j].HeroId == hero.HeroId && costHero[j].StarItem.UpStar == heroConfig.UpStarStar[i] {
						isUse[j] = true
						count--
					}
				}
			case UPSTAR_TYPE_COLOR:
				for j := 0; j < len(msg.CostHeroKeyId[size]); j++ {
					if count <= 0 {
						break
					}
					if isUse[j] {
						continue
					}
					if costHero[j].StarItem.UpStar == heroConfig.UpStarStar[i] {
						isUse[j] = true
						count--
					}
				}
			case UPSTAR_TYPE_ALL:
				for j := 0; j < len(msg.CostHeroKeyId[size]); j++ {
					if count <= 0 {
						break
					}
					if isUse[j] {
						continue
					}
					costHeroConfig := GetCsvMgr().HeroConfig[costHero[j].HeroId]
					if costHeroConfig == nil {
						self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
						self.SendInfoSyn()
						return
					}
					if costHero[j].StarItem.UpStar == heroConfig.UpStarStar[i] {
						isUse[j] = true
						count--
					}
				}
			case UPSTAR_TYPE_SPECTIAL:
				for j := 0; j < len(msg.CostHeroKeyId[size]); j++ {
					if count <= 0 {
						break
					}
					if isUse[j] {
						continue
					}
					costHeroConfig := GetCsvMgr().HeroConfig[costHero[j].HeroId]
					if costHeroConfig == nil {
						//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
						continue
					}
					if costHero[j].StarItem.UpStar == heroConfig.UpStarStar[i] {
						isUse[j] = true
						count--
					}
				}
			}
		}

		isCan := true
		for i := 0; i < len(isUse); i++ {
			if !isUse[i] {
				isCan = false
				break
			}
		}

		returnItem := make(map[int]*Item)

		if isCan {
			//扣除英雄
			for i := 0; i < len(msg.CostHeroKeyId[size]); i++ {
				costHero := self.player.getHero(msg.CostHeroKeyId[size][i])
				offEquip, artifact, horse := self.OffEquipAll(costHero)
				msgRel.OffEquip = append(msgRel.OffEquip, offEquip...)
				msgRel.OffArtifact = append(msgRel.OffArtifact, artifact...)
				msgRel.OffHorseInfo = append(msgRel.OffHorseInfo, horse...)

				configExp := GetCsvMgr().HeroExpConfigMap[costHero.HeroLv]
				if self.player.GetModule("crystal").(*ModResonanceCrystal).IsResonanceHero(msg.CostHeroKeyId[size][i]) {
					configExp = GetCsvMgr().HeroExpConfigMap[costHero.OriginalLevel]
				}
				if configExp == nil {
					continue
				}
				AddItemMapHelper(returnItem, configExp.ResetItems, configExp.ResetNums)

				self.player.GetModule("hero").(*ModHero).DeleteHero(msg.CostHeroKeyId[size][i])
			}
			oldstar := hero.StarItem.UpStar
			hero.StarItem.UpStar++
			self.player.HandleTask(TASK_TYPE_STAR_UP_COUNT, 1, hero.StarItem.UpStar, 0)
			self.player.HandleTask(TASK_TYPE_GET_HERO_WIDE, 1, hero.StarItem.UpStar, 0)
			self.player.HandleTask(TASK_TYPE_STAR_UP_COUNT_EQUAL, 1, hero.StarItem.UpStar, 0)
			//处理图鉴
			info, ok := self.Sql_Hero.handBook[hero.HeroId]
			if ok && info.StarMax < hero.StarItem.UpStar {
				info.StarMax = hero.StarItem.UpStar
			}
			// 检测天赋
			hero.CheckStageTalent()
			self.player.countHeroFight(hero, ReasonHeroStarUp)

			msgRel.HeroKeyId = append(msgRel.HeroKeyId, msg.HeroKeyId[size])
			msgRel.StarItem = append(msgRel.StarItem, hero.StarItem)
			msgRel.StageTalent = append(msgRel.StageTalent, hero.StageTalent)
			msgRel.CostHeroKeyId = append(msgRel.CostHeroKeyId, msg.CostHeroKeyId[size])

			//处理返还
			if len(returnItem) > 0 {
				items := self.player.AddObjectItemMap(returnItem, "升阶材料返还", hero.HeroId, 0, 0)
				msgRel.GetItems = append(msgRel.GetItems, items...)
			}

			self.player.GetModule("support").(*ModSupportHero).UpdataHero(hero.HeroKeyId)
			self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(hero.HeroKeyId)
			self.player.GetModule("entanglement").(*ModEntanglement).UpdateHero(hero.HeroKeyId)
			self.player.GetModule("lifetree").(*ModLifeTree).StarUp(hero.HeroId, hero.StarItem.UpStar)
			self.ChangeVoidHeroResonance(hero.HeroKeyId, 0)
			self.GetAllHeroStars()

			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_UP_STAR, hero.HeroId, oldstar, hero.StarItem.UpStar, "英雄升阶", 0, hero.HeroKeyId, self.player)

		}
	}
	self.player.GetModule("entanglement").(*ModEntanglement).SendInfo([]byte{})
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetOfflineInfoMgr().UpdateHeroSkin(self.player)
	self.player.GetModule("task").(*ModTask).SendUpdate()
	self.player.GetModule("lifetree").(*ModLifeTree).SendRedPointInfo()

}

func (self *ModHero) HeroLock(body []byte) {

	var msg C2S_HeroLock
	json.Unmarshal(body, &msg)

	//目标英雄是否存在
	hero := self.player.getHero(msg.HeroKeyId)
	if hero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_NOT_EXIST"))
		return
	}

	if msg.LockAction != LOGIC_FALSE && msg.LockAction != LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_MSG_ERROR"))
		return
	}

	//已经是这个状态
	if msg.LockAction == hero.IsLock {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STAR_HERO_ALREADY_STATE"))
		return
	}

	hero.IsLock = msg.LockAction

	var msgRel S2C_HeroLock
	msgRel.Cid = "herolock"
	msgRel.HeroKeyId = msg.HeroKeyId
	msgRel.IsLock = hero.IsLock
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	if msg.LockAction == LOGIC_TRUE {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_LOCK, hero.HeroId, hero.HeroKeyId, 0, "锁定英雄", 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_UNLOCK, hero.HeroId, hero.HeroKeyId, 0, "解锁英雄", 0, 0, self.player)
	}

}

func (self *ModHero) DeleteHero(heroKeyId int) {
	//从阵位上删除
	self.player.GetModule("team").(*ModTeam).DeleteHeroFromTeam(heroKeyId)
	//从租借信息中删除
	//GetHireHeroInfoMgr().DeteleHire(self.player, heroKeyId)
	GetMasterMgr().FriendPRC.DeleteHireFriend(self.player.GetUid(), heroKeyId)
	if GetOfflineInfoMgr().IsBaseHero(self.player.Sql_UserBase.Uid, heroKeyId) {
		self.player.NoticeCenterBaseInfo()
	}
	GetOfflineInfoMgr().DeleteBaseHero(self.player.Sql_UserBase.Uid, heroKeyId)
	self.player.GetModule("friend").(*ModFriend).BaseHeroUpdata(heroKeyId)
	voidherokey := 0
	hero := self.GetHero(heroKeyId)
	if hero != nil && hero.VoidHero == 0 && hero.Resonance != 0 {
		voidherokey = hero.Resonance
	}

	delete(self.Sql_Hero.info, heroKeyId)
	self.player.GetModule("support").(*ModSupportHero).UpdataHero(heroKeyId)
	self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(heroKeyId)
	self.player.GetModule("arena").(*ModArena).UpdateFormat(heroKeyId)
	self.player.GetModule("arenaspecial").(*ModArenaSpecial).UpdateFormat(heroKeyId)
	self.player.GetModule("entanglement").(*ModEntanglement).DeleteHero(heroKeyId)
	self.player.GetModule("crossarena").(*ModCrossArena).UpdateFormat()
	self.player.GetModule("crossarena3v3").(*ModCrossArena3V3).UpdateFormat()
	self.ChangeVoidHeroResonance(heroKeyId, voidherokey)
	self.GetAllHeroStars()
}

func (self *ModHero) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case interchange:
		var msg C2S_Interchange
		json.Unmarshal(body, &msg)
		var ret S2C_Interchange
		ret.Cid = interchange
		ret.Ret, ret.Item = self.Interchange(msg.Heroid[0], msg.Heroid[1])
		self.player.SendMsg(interchange, HF_JtoB(&ret))
		return true
	case upgradeTalentAction:
		var msg C2S_UpgradeTalent
		json.Unmarshal(body, &msg)
		self.upgradeTalent(&msg)
		return true
	case activatefate:
		var msg C2S_ActivateFate
		json.Unmarshal(body, &msg)
		self.activateFate(&msg)
		return true
	case herodecompose:
		var msg C2S_HeroDecompose
		json.Unmarshal(body, &msg)
		self.decompose(&msg)
		return true
	case talentReset:
		var msg C2S_ResetTalent
		json.Unmarshal(body, &msg)
		self.TalentReset(&msg)
		return true
	case compose:
		var msg C2S_Compose
		json.Unmarshal(body, &msg)
		self.Compose(&msg)
		return true
	}

	return false
}

func (self *ModHero) OnSave(sql bool) {
	self.Encode()
	self.Sql_Hero.Update(sql)
}

func (self *ModHero) CheckHorse() {
	for _, herodata := range self.Sql_Hero.info {
		if herodata == nil {
			continue
		}

		if herodata.Horse <= 0 {
			continue
		}

		if self.player.GetModule("horse").(*ModHorse).GetHorse(herodata.Horse) == nil {
			herodata.Horse = 0
		}
	}
}

func (self *ModHero) doAddHeroTask() {

}

// 加一个英雄
// 若star = 0，则送的星级读heroid里
func (self *ModHero) AddHero(heroId int, param1, param2, param3 int, dec string) []*Hero {

	heroRel := make([]*Hero, 0)
	for i := 0; i < param1; i++ {
		if i > 100 {
			return heroRel
		}
		heroInfo := new(Hero)
		heroInfo.HeroKeyId = self.MaxKey()
		heroInfo.HeroId = heroId
		heroInfo.Uid = self.player.ID
		heroInfo.checkStarItem(param2)
		heroInfo.checkEquip()
		heroInfo.LvUp(1)
		heroInfo.CheckStageTalent()

		self.Sql_Hero.info[heroInfo.HeroKeyId] = heroInfo

		// 任务
		self.doAddHeroTask()
		self.player.countHeroFight(heroInfo, ReasonHeroNew)
		self.player.CheckHero(heroId)

		config := GetCsvMgr().GetHeroMapConfig(heroId, heroInfo.StarItem.UpStar)
		if nil != config {
			if config.Attribute == HERO_ATTRIBUTE_AIR {
				heroInfo.VoidHero = 1
			}
		}

		self.player.HandleTask(HeroUpStarNumTask, 0, 0, 0)

		self.GetAllHeroStars()
		//GetTopHeroStarsMgr().UpdateRank(self.GetAllHeroStars(), self.player)
		self.SynHero(heroInfo)
		self.CheckBackOpen(heroInfo)
		self.CheckHireUpdate(heroInfo)
		heroRel = append(heroRel, heroInfo)
		//处理图鉴
		info, ok := self.Sql_Hero.handBook[heroInfo.HeroId]
		if !ok {
			self.Sql_Hero.handBook[heroInfo.HeroId] = new(HandBookInfo)
			self.Sql_Hero.handBook[heroInfo.HeroId].State = CANTAKE
			self.Sql_Hero.handBook[heroInfo.HeroId].StarMax = heroInfo.StarItem.UpStar

			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HANDBOOK_ACTIVATE, heroInfo.HeroId, 0, 0, "激活英雄图鉴", 0, 0, self.player)
		} else {
			if info.StarMax < heroInfo.StarItem.UpStar {
				info.StarMax = heroInfo.StarItem.UpStar
			}
		}

		if dec != "创建角色" {
			self.player.HandleTask(TASK_TYPE_GET_HERO, 1, heroInfo.StarItem.UpStar, config.Attribute)
			self.player.HandleTask(TASK_TYPE_GET_HERO_WIDE, 1, heroInfo.StarItem.UpStar, 0)
		}

		if len(self.Sql_Hero.info) <= RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX {
			self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(heroInfo.HeroKeyId)
		}

		self.player.GetModule("lifetree").(*ModLifeTree).StarUp(heroId, heroInfo.StarItem.UpStar)

		self.makeAddHeroLog(heroId, param1, param3, dec, heroInfo)

		//商店CHECK
		if config.Attribute == HERO_ATTRIBUTE_AIR {
			self.player.GetModule("shop").(*ModShop).CheckHeroItem(heroId)
		}
	}
	self.player.GetModule("lifetree").(*ModLifeTree).SendRedPointInfo()
	return heroRel
}

func (self *ModHero) makeAddHeroLog(heroid int, param1, param2 int, dec string, value *Hero) {
	// 写入行为数据库，不再写入游戏日志
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_ACTIVE, heroid, value.StarItem.UpStar, param2%1000, dec, 1, param2/1000, self.player)
	GetServer().sendLog_herogain(self.player, value, dec)
}

//
func (self *ModHero) SynHero(hero *Hero) {
	if hero == nil {
		return
	}

	cid := "synhero"
	self.player.SendMsg(cid, HF_JtoB(&S2C_SynHero{
		Cid:  cid,
		Hero: hero,
	}))
}

// 有相同英雄,则转换为英雄碎片
func (self *ModHero) changeItem(heroId int, param1, param2 int, dec string) (int, int) {
	config := GetCsvMgr().GetHeroConfig(heroId)
	itemId := 0
	itemNum := 0
	if config != nil {
		if len(config.CardIds) == 1 {
			itemId = config.CardIds[0]
			self.player.AddObjectLst(config.CardIds, config.CardNums, dec, param1, param2, 0)
			self.player.CheckHero(heroId)
		}

		if len(config.CardNums) == 1 {
			itemNum = config.CardNums[0]
		}
	}
	return itemId, itemNum
}

func (self *ModHero) GetHero(heroKeyId int) *Hero {
	value, ok := self.Sql_Hero.info[heroKeyId]
	if ok {
		return value
	}

	return nil
}

func (self *ModHero) GetHeroes() map[int]*Hero {
	return self.Sql_Hero.info
}

func (self *ModHero) RemoveHero(heroid int) bool {
	_, ok := self.Sql_Hero.info[heroid]
	if ok {
		delete(self.Sql_Hero.info, heroid)
		return true
	}

	return false
}

func (self *ModHero) RemoveHeroByKeyId(keyId int) bool {
	_, ok := self.Sql_Hero.info[keyId]
	if ok {
		self.DeleteHero(keyId)
		return true
	}

	return false
}

// 得到最强hero阵容
func (self *ModHero) GetBestFormat() ([]int, int64) {
	lsthero := make(LstHero, 0)
	for _, value := range self.Sql_Hero.info {
		lsthero = append(lsthero, value)
	}

	sort.Sort(LstHero(lsthero))

	fight := int64(0)
	lst := make([]int, 0)
	for i := 0; i < HF_MinInt(5, len(lsthero)); i++ {
		lst = append(lst, lsthero[i].HeroKeyId)
		//LogDebug(lsthero[i].HeroId, ", 战斗力:", lsthero[i].Fight)
		fight += lsthero[i].Fight
	}

	return lst, fight
}

// 获得虚空英雄
func (self *ModHero) GetVoidHero(heroid int) *Hero {
	for _, value := range self.Sql_Hero.info {
		if value.HeroId == heroid && value.VoidHero == 1 {
			return value
		}
	}
	return nil
}

// 得到最强hero阵容
func (self *ModHero) GetBestFormat2() []int {
	lsthero := make(LstHero, 0)
	for _, value := range self.Sql_Hero.info {
		lsthero = append(lsthero, value)
	}

	sort.Sort(LstHero(lsthero))
	lst := make([]int, 0)
	for i := 0; i < len(lsthero); i++ {
		lst = append(lst, lsthero[i].HeroKeyId)
	}

	return lst
}

// 得到最强hero阵容 按星级排序
func (self *ModHero) GetBestFormat3() []int {
	lsthero := make(LstHeroStar, 0)
	for _, value := range self.Sql_Hero.info {
		lsthero = append(lsthero, value)
	}

	sort.Sort(LstHeroStar(lsthero))
	lst := make([]int, 0)
	for i := 0; i < len(lsthero); i++ {
		lst = append(lst, lsthero[i].HeroKeyId)
	}

	return lst
}

// 得到最强hero阵容 按等级排序
func (self *ModHero) GetBestFormat4() []int {
	lsthero := make(LstHeroLevel, 0)
	for _, value := range self.Sql_Hero.info {
		lsthero = append(lsthero, value)
	}

	sort.Sort(LstHeroLevel(lsthero))
	lst := make([]int, 0)
	for i := 0; i < len(lsthero); i++ {
		lst = append(lst, lsthero[i].HeroKeyId)
	}

	return lst
}

// 得到上阵战斗力
func (self *ModHero) GetFight(herolst []int) int64 {
	fight := int64(0)
	for i := 0; i < len(herolst); i++ {
		hero := self.GetHero(herolst[i])
		if hero == nil {
			continue
		}
		fight += hero.Fight
	}

	return fight
}

// 合成
// 0：成功，1：已经拥有，2：金币不足，3：碎片不足 4 配置不存在 5 配置异常
func (self *ModHero) Compound(heroId int) {
	/*
		config := GetCsvMgr().GetHeroConfig(heroId)
		if config == nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_HEROIC_CONFIGURATION_DOES_NOT_EXIST"))
			return
		}

		if len(config.CompoundIds) != len(config.CompoundNums) {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HERO_CONFIGURATION_EXCEPTION"))
			return
		}

		// 检查配置消耗
		if err := self.player.HasObjectOk(config.CompoundIds, config.CompoundNums); err != nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_INSUFFICIENT_MATERIAL"))
			return
		}

		var res []PassItem
		items := self.player.RemoveObjectLst(config.CompoundIds, config.CompoundNums, "英雄合成", heroId, 0, 0)
		res = append(res, items...)
		self.AddHero(heroId, heroId, 0, "英雄合成")

		hero := self.player.getHero(heroId)

		var msg S2C_Synthesis
		msg.Cid = heroSynthesis
		msg.Item = res
		msg.Hero = hero

		self.player.SendMsg(heroSynthesis, HF_JtoB(msg))

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_ENERGY_LOOT, heroId, 0, 0, "英雄合成", 0, 0, self.player)
	*/
}

// 武将互换
// 武将互换，等级，经验，神器等级，统兵等级
func (self *ModHero) Interchange(heroid1 int, heroid2 int) (int, []PassItem) {
	costitem := make([]PassItem, 0)
	hero := self.GetHero(heroid1)
	if hero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NOPARAM"))
		return -1, costitem
	}

	hero1 := self.GetHero(heroid2)
	if hero1 == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NOPARAM"))
		return -1, costitem
	}

	gemNeed := 200
	if self.player.GetObjectNum(DEFAULT_GEM) < gemNeed {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_NOT_ENOUGH_GEM_FAIL"))
		return -1, costitem
	}

	self.player.AddObject(DEFAULT_GEM, -gemNeed, heroid1, heroid2, 0, "武将互换")
	costitem = append(costitem, PassItem{ItemID: DEFAULT_GEM, Num: -gemNeed})

	self.player.countHeroFight(hero, 0)
	self.player.countHeroFight(hero1, 0)

	GetServer().sendLog_herodisplace(self.player, hero, hero1)
	return 0, costitem
}

func (self *ModHero) SendInfo() {
	// 登录检查
	self.player.HandleTask(TalentTask, 0, 0, 0)
	//self.player.HandleTask(HaveDinivityTask, 0, 0, 0)
	//self.player.HandleTask(AllDinivityTask, 0, 0, 0)
	var msg S2C_HeroInfo
	msg.Cid = "herolst"
	msg.Newhero = false
	//msg.TotalStars = self.Sql_Hero.Totalstars
	msg.Reborn = self.Sql_Hero.Reborn
	msg.BuyPosNum = self.Sql_Hero.BuyPosNum
	msg.AutoFire = self.Sql_Hero.AutoFire
	msg.BackOpen = self.Sql_Hero.BackOpen
	for _, value := range self.Sql_Hero.info {
		if value != nil {
			value.checkStarItem(0)
			value.checkEquip()
		}

		self.player.checkHeroFight(value, ReasonHeroCheck)
		msg.Herolst = append(msg.Herolst, value)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("herolst", smsg)

	self.SendHandBookInfo()

	updateHire := make(map[int]*HireHero)
	for _, value := range self.Sql_Hero.info {
		if value == nil {
			continue
		}
		if value.StarItem.UpStar < HIRE_STAR_LIMIT {
			continue
		}
		updateHire[value.HeroKeyId] = self.NewHireHero(value)
	}

	GetMasterMgr().FriendPRC.UpdateHireAll(self.player.GetUid(), updateHire, LOGIC_FALSE)
}

//增加一条消息来同步图鉴，为了避免更新旧版本出现问题。新开项目可以合并这个  20200604
func (self *ModHero) SendHandBookInfo() {
	var msg S2C_HandBook
	msg.Cid = "herohandbook"
	msg.HandBook = self.Sql_Hero.handBook
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}

func (self *ModHero) GetShowHeros() ([]*JS_PlayerHero, [][]*JS_HeroEquip) {
	dataHero := make([]*JS_PlayerHero, 0)
	dataEquip := make([][]*JS_HeroEquip, 0)

	info := GetOfflineInfoMgr().GetBaseHero(self.player)
	if info == nil {
		return dataHero, dataEquip
	}

	heroIds := make([]int, 0)
	for i := 0; i < len(info); i++ {
		if info[i] == nil {
			continue
		}
		heroData := self.player.getHero(info[i].HeroKeyId)
		if heroData == nil {
			continue
		}
		heroIds = append(heroIds, heroData.HeroId)
	}
	teamAttr := GetRobotMgr().GetTeamAttr(heroIds)

	for i := 0; i < len(info); i++ {
		if info[i] == nil {
			hero := &JS_PlayerHero{
				HeroId:          0,
				Star:            0,
				Level:           0,
				ArtifactId:      0,
				ArtifactLv:      0,
				ExclusiveId:     0,
				ExclusiveLv:     0,
				ExclusiveUnLock: LOGIC_FALSE,
				Skin:            0,
				Attr:            make(map[int]int64, 0),
				Talent:          nil,
			}
			equip := make([]*JS_HeroEquip, 0)
			for j := 0; j < 4; j++ {
				equip = append(equip, &JS_HeroEquip{ItemId: 0, Level: 0})
			}
			dataHero = append(dataHero, hero)
			dataEquip = append(dataEquip, equip)
		} else {
			hero := new(JS_PlayerHero)
			hero.HeroId = info[i].HeroId
			hero.Star = info[i].StarItem.UpStar
			hero.Level = info[i].HeroLv
			hero.Skin = info[i].Skin
			if len(info[i].ArtifactEquipIds) > 0 && info[i].ArtifactEquipIds[0] != nil {
				hero.ArtifactId = info[i].ArtifactEquipIds[0].Id
				hero.ArtifactLv = info[i].ArtifactEquipIds[0].Lv
			}
			if info[i].ExclusiveEquip != nil {
				hero.ExclusiveId = info[i].ExclusiveEquip.Id
				hero.ExclusiveLv = info[i].ExclusiveEquip.Lv
				hero.ExclusiveUnLock = info[i].ExclusiveEquip.UnLock
			}
			if info[i].StageTalent != nil {
				hero.Talent = new(StageTalent)
				HF_DeepCopy(hero.Talent, info[i].StageTalent)
			}

			equip := make([]*JS_HeroEquip, 0)
			if info[i].EquipIds == nil {
				for j := 0; j < 4; j++ {
					equip = append(equip, &JS_HeroEquip{ItemId: 0, Level: 0})
				}
			} else {
				for j := 0; j < len(info[i].EquipIds); j++ {
					if info[i].EquipIds[j] != nil {
						equip = append(equip, &JS_HeroEquip{ItemId: info[i].EquipIds[j].Id, Level: info[i].EquipIds[j].Lv})
					} else {
						equip = append(equip, &JS_HeroEquip{ItemId: 0, Level: 0})
					}
				}
			}
			hero.Attr = make(map[int]int64, 0)
			heroData := self.player.getHero(info[i].HeroKeyId)
			if heroData != nil {
				pAttrWrapper := heroData.GetAttr2(self.player, 0, teamAttr)
				for i := 0; i < len(pAttrWrapper.Base); i++ {
					value := int64(pAttrWrapper.Base[i])
					if value > 0 {
						hero.Attr[i] = value
					}
				}
				hero.Attr[ATTR_TYPE_FIGHT] = pAttrWrapper.FightNum
			}
			dataHero = append(dataHero, hero)
			dataEquip = append(dataEquip, equip)
		}
	}

	return dataHero, dataEquip
}

func (self *ModHero) ShopJudge() (bool, bool) {
	relJudgeCamp7 := false
	relJudeg7001 := false
	for heroId, sign := range self.Sql_Hero.handBook {
		//正常不会出现 除非改版之前的老号
		if sign == nil {
			continue
		}
		if sign.State == CANTFINISH {
			continue
		}
		//if heroId == 7001 {
		if heroId == 7002 {
			relJudeg7001 = true
		}
		//看看是否是虚空英雄
		if !relJudgeCamp7 {
			heroConfig := GetCsvMgr().HeroConfigMap[heroId][1]
			if heroConfig != nil && heroConfig.Attribute == HERO_ATTRIBUTE_AIR {
				relJudgeCamp7 = true
			}
		}
	}
	if !relJudeg7001 {
		//判断魂石是否超过60  如果后面增加多个虚空英雄 这里需要改
		if self.player.GetObjectNum(12700101) >= 60 {
			relJudeg7001 = true
		}
	}
	return relJudgeCamp7, relJudeg7001
}

func (self *ModHero) GetHeroNum() int {

	hero_num := len(self.Sql_Hero.info)
	return hero_num
}

//武将阶级总数
func (self *ModHero) GetHeroTotalColor() int {
	return 0
}

func (self *San_Hero) Decode() { // 将数据库数据写入data
	err := json.Unmarshal([]byte(self.Info), &self.info)
	if err != nil {
		log.Println("解析英雄失败:", err)
	}
}

// 检查属性
func (self *ModHero) checkHero() {
	for _, value := range self.Sql_Hero.info {
		if value == nil {
			continue
		}
		value.checkStarItem(0)
		value.checkEquip()
	}

	//新增加的图鉴
	if self.Sql_Hero.handBook == nil {
		self.Sql_Hero.handBook = make(map[int]*HandBookInfo)
	}
}

// 获得英雄Id
func (self *Hero) getHeroId() int {
	return self.HeroId
}

// 获得一个英雄
func (self *ModHero) randHero() int {
	for _, v := range self.Sql_Hero.info {
		return v.HeroKeyId
	}
	return 0
}

// 英雄碎片分解
func (self *ModHero) decompose(pMsg *C2S_HeroDecompose) {
	itemIds := pMsg.ItemIds
	itemNums := pMsg.ItemNums
	if len(itemIds) != len(itemNums) {
		self.player.SendErr("len(itemIds) != len(itemNums)!")
		return
	}

	if len(itemIds) <= 0 {
		self.player.SendErr("len(itemIds) <= 0")
		return
	}

	if len(itemNums) <= 0 {
		self.player.SendErr("len(itemNums) <= 0")
		return
	}

	for _, num := range itemNums {
		if num <= 0 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HERO_THE_NUMBER_OF_PROPS_SENT"))
			return
		}
	}

	// 检查道具是否充足
	err := self.player.HasObjectOk(itemIds, itemNums)
	if err != nil {
		self.player.SendErr(err.Error())
		return
	}

	// 计算能获得的道具
	itemMap := make(map[int]*Item)
	for index := range itemIds {
		itemId := itemIds[index]
		itemNum := itemNums[index]
		if itemId <= 0 || itemNum <= 0 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HERO_THE_NUMBER_OF_PROPS_SENT"))
			return
		}

		config := GetCsvMgr().getHeroDecompose(itemId)
		if config == nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HERO_PROPS_DO_NOT_EXIST"))
			return
		}

		pItem, ok := itemMap[config.DecomposeItem]
		if !ok {
			itemMap[config.DecomposeItem] = &Item{config.DecomposeItem, config.DecomposeNum * itemNum}
		} else {
			pItem.ItemNum += config.DecomposeNum * itemNum
		}
	}

	var outItems []PassItem
	items := self.player.AddObjectItemMap(itemMap, "分解", 0, 0, 0)
	outItems = append(outItems, items...)
	removeItems := self.player.RemoveObjectLst(itemIds, itemNums, "分解", 0, 0, 0)
	outItems = append(outItems, removeItems...)

	data := &S2C_HeroDecompose{
		Cid:   herodecompose,
		Items: outItems,
	}
	self.player.SendMsg(data.Cid, HF_JtoB(data))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_DECOMPOUND, 1, 0, 0, "分解", 0, 0, self.player)
}

func (self *ModHero) getHeroStarNum(n2, n3, n4 int) int {
	heroInfo := self.Sql_Hero.info
	num := 0
	for _, value := range heroInfo {
		if value == nil {
			continue
		}

		if value.StarItem == nil {
			continue
		}

		config := GetCsvMgr().GetHeroConfig(value.HeroId)
		if config == nil {
			continue
		}

		bStar, bHeroID, bStep := false, false, false

		if value.StarItem.UpStar >= n2 {
			bStar = true
		}

		if n3 == 0 {
			bHeroID = true
		} else {
			if n3 == value.HeroId {
				bHeroID = true
			}
		}

		if config.FightIndex >= n4 {
			bStep = true
		}

		if bStar && bHeroID && bStep {
			num++
		}
	}
	return num
}

func (self *ModHero) getHeroBiggestLevel() int {
	heroInfo := self.Sql_Hero.info
	nLevel := 0
	for _, value := range heroInfo {
		if value == nil {
			continue
		}

		if nLevel < value.HeroLv {
			nLevel = value.HeroLv
		}
	}
	return nLevel
}

func (self *ModHero) getHeroNumByLevel(level int) int {
	heroInfo := self.Sql_Hero.info
	num := 0
	for _, value := range heroInfo {
		if value == nil {
			continue
		}

		if level <= value.HeroLv {
			num++
		}
	}
	return num
}

//! 刷新处理  每日重置次数
func (self *ModHero) OnRefresh() {
	self.Sql_Hero.Reborn = 0
}

//! 英雄重生  20190413 by zy
/*
func (self *ModHero) reBorn(pMsg *C2S_ReBorn) {
	var msg S2C_Reborn

	pHero := self.GetHero(pMsg.HeroId)
	if pHero == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TIGER_HEROES_DO_NOT_EXIST"))
		return
	}

	//检查该英雄是否满足重生条件   1.星级大于1  2.星级等于1但激活过星格 3.任意天赋大于0
	hasStarSlots := false
	hasTalent := false

	for _, star := range pHero.StarItem.Slots {
		if star > 1 {
			hasStarSlots = true
			break
		}
	}

	for _, tal := range pHero.TalentItem.Talents {
		if tal.Lv > 0 {
			hasTalent = true
			break
		}
	}
	//不满足重生条件
	if pHero.StarItem.UpStar == 1 && !hasStarSlots && !hasTalent {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REBORN_STAR_NOT_ENOUGH"))
		return
	}
	//用户等级不足
	flag, _ := GetCsvMgr().IsLevelOpen(self.player.Sql_UserBase.Level, OPEN_LEVEL_80)
	if !flag {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REBORN_LV_NOT_ENOUGH"))
		return
	}
	//今日重生次数达到VIP上限
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REBORN_ERROR_CONFIG"))
		return
	}
	if self.Sql_Hero.Reborn > vipcsv.RebornTimes {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REBORN_TIMES_LIMIT"))
		return
	}

	_, configLv := GetCsvMgr().IsLevelOpen(self.player.Sql_UserBase.Level, OPEN_LEVEL_REBORN_FREE)

	if configLv == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REBORN_ERROR_CONFIG"))
		return
	} else if self.player.Sql_UserBase.Level >= configLv {
		//检测消耗是否足够
		costConfig := GetCsvMgr().GetTariffConfig(TariffReBorn, self.Sql_Hero.Reborn+1)
		if costConfig == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_REBORN_HERO_IN_TEAM"))
			return
		}
		if self.player.HasObjectOk(costConfig.ItemIds, costConfig.ItemNums) != nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_NOT_ENOUGH_GEM_FAIL"))
			return
		}

		//扣除物品
		msg.Items = self.player.RemoveObjectLst(costConfig.ItemIds, costConfig.ItemNums, "重生", 0, 0, 0)
	}

	totalCost := make(map[int]*Item)
	heroConfig := GetCsvMgr().HeroStarMap[pHero.HeroId]
	//返还升星材料
	if pHero.StarItem.UpStar > 1 {
		for i := pHero.StarItem.UpStar; i > 1; i-- {
			starConfig := heroConfig[i]
			AddItemMapHelper(totalCost, starConfig.StarLvIds, starConfig.StarLvNums)
		}
	}

	//返还星格材料
	if hasStarSlots {
		for i := 0; i < maxStarSlots; i++ {
			if pHero.StarItem.Slots[i] > 1 {
				for j := 1; j <= pHero.StarItem.Slots[i]; j++ {
					starConfig := heroConfig[j]
					slotItemIds, slotItemNums := starConfig.getCost(i + 1)
					AddItemMapHelper(totalCost, slotItemIds, slotItemNums)
				}
			}
		}
	}

	//返还天赋材料
	if hasTalent {
		for _, pTalent := range pHero.TalentItem.Talents {
			if pTalent.Lv > 0 {
				config := GetCsvMgr().GetTalentConfig(pTalent.Id, pTalent.Lv)
				if config == nil {
					continue
				}
				AddItemMapHelper(totalCost, config.ReturnItems, config.Returnnums)
			}
		}
	}

	//重生属性还原
	pHero.reborn()
	self.player.countHeroFight(pHero, ReasonReborn)

	//返还物品
	outItems := self.player.AddObjectItemMap(totalCost, "重生", 0, 0, 0)
	self.Sql_Hero.Reborn++

	//发送消息
	msg.Cid = reBornAction
	msg.HeroId = pHero.getHeroId()
	msg.ItemsReturn = outItems
	msg.StarItem = pHero.StarItem
	msg.TalentItem = pHero.TalentItem
	msg.Attr = pHero.GetStarAttr()
	msg.ReBorn = self.Sql_Hero.Reborn
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_REBORN, pMsg.HeroId, 0, 0, "重生", 0, 0, self.player)

	self.GetAllHeroStars()
}

*/

// 英雄重生
/*
func (self *Hero) reborn() {

	//initStar:= GetCsvMgr().SimpleConfigMap[1]
	for i := 0; i < maxStarSlots; i++ {
		self.StarItem.Slots[i] = 1
	}
	self.StarItem.UpStar = 1

	for i := 0; i < len(self.TalentItem.Talents); i++ {
		self.TalentItem.Talents[i].Lv = 0
	}

	self.cacStarAtt()
	self.cacTalentAtt()
}

*/

// 计算所有英雄总星级 并传入活动
func (self *ModHero) GetAllHeroStars() [HERO_STAR_TOTAL_MAX]int {
	heroInfo := self.GetBestFormat3()

	herouse := []int{}
	count := [HERO_STAR_TOTAL_MAX]int{}
	for _, v := range heroInfo {
		if v == 0 {
			continue
		}

		value := self.GetHero(v)
		if value == nil {
			continue
		}

		if value.StarItem == nil {
			continue
		}

		if value.StarItem.UpStar > 0 {
			count[HERO_STAR_TOTAL] += value.StarItem.UpStar
		}
		// 该类英雄已经使用
		find := false
		for _, t := range herouse {
			if value.HeroId == t {
				find = true
			}
		}
		if find {
			continue
		}

		heroConfig := GetCsvMgr().GetHeroMapConfig(value.HeroId, value.StarItem.UpStar)
		if nil != heroConfig {
			if heroConfig.Attribute > HERO_STAR_TOTAL && heroConfig.Attribute < HERO_STAR_TOTAL_MAX {
				rankConfig := GetCsvMgr().GetRankListIntegralConfig(value.StarItem.UpStar)
				if nil != rankConfig {
					count[heroConfig.Attribute] += rankConfig.RankValue
				}
			}
		}

		herouse = append(herouse, value.HeroId)
	}
	time := TimeServer().Unix()
	for i := 0; i < HERO_STAR_TOTAL_MAX; i++ {
		if count[i] != self.Sql_Hero.totalStars[i].Stars {
			self.Sql_Hero.totalStars[i].Stars = count[i]
			self.Sql_Hero.totalStars[i].StarTime = time
			self.player.HandleTask(TASK_TYPE_HERO_STAR_POINT, count[i], i, 0)
			GetTopTalentMgr().UpdateRank(i, count[i], self.player)
		}
	}
	return count
}

func (self *ModHero) Compose(msg *C2S_Compose) {
	config := GetCsvMgr().RuneComposeMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("compose001"))
		return
	}
	//测试中，最多只给10次
	if msg.Num > 10 {
		msg.Num = 10
	}
	numAll := msg.Num
	inItemMap := make(map[int]*Item)
	outItemMap := make(map[int]*Item)
	for {
		if numAll > 0 {
			ok, inItem, outItem := self.ComposeOnce(config)
			if !ok {
				break
			}
			AddItemMapHelper2(inItemMap, inItem)
			AddItemMapHelper2(outItemMap, outItem)
			numAll--
		} else {
			LogDebug("结束")
			break
		}
	}

	//材料不足
	if msg.Num-numAll == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_INSUFFICIENT_MATERIAL"))
		return
	}

	var msgRel S2C_RuneCompose
	msgRel.Cid = compose
	msgRel.Id = msg.Id
	msgRel.Num = numAll
	msgRel.InItem = inItemMap
	msgRel.OutItem = outItemMap
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModHero) ComposeOnce(config *RuneCompose) (bool, []PassItem, []PassItem) {
	// 检查道具消耗
	if err := self.player.HasObjectOk(config.Item, config.Num); err != nil {
		return false, nil, nil
	}

	// 扣除道具
	inItems := self.player.RemoveObjectLst(config.Item, config.Num, "合成", config.Id, config.Type, 0)
	//产出
	outitems := make([]PassItem, 0)

	if config.Type == COMPOSE_TYPE_RUNE {
		itemid := 0
		allRate := 0
		for i := 0; i < len(config.Rune); i++ {
			allRate += config.Probability[i]
		}

		randNum := HF_GetRandom(allRate)
		for i := 0; i < len(config.Rune); i++ {
			if randNum < config.Probability[i] {
				itemid = config.Rune[i]
			}
		}

		itemIdRel, itemNumRel := self.player.AddObject(itemid, 1, 9, 0, 0, "事件宝箱碎片")

		outitems = append(outitems, PassItem{itemIdRel, itemNumRel})
	} else if config.Type == COMPOSE_TYPE_EQUIP {
		groupId := 0
		allRate := 0
		for i := 0; i < len(config.Rune); i++ {
			allRate += config.Probability[i]
		}

		randNum := HF_GetRandom(allRate)
		for i := 0; i < len(config.Rune); i++ {
			if randNum < config.Probability[i] {
				groupId = config.Random[i]
				lootItems := GetLootMgr().LootItem(groupId, nil)
				items := self.player.AddObjectItemMap(lootItems, "合成", 0, 0, 0)
				outitems = append(outitems, items...)
				break
			}
		}
	}

	return true, inItems, outitems
}

func (self *ModHero) MaxKey() int {
	self.Sql_Hero.MaxKey += 1
	return self.Sql_Hero.MaxKey
}

func (self *ModHero) HeroBuyPos(body []byte) {
	config := GetCsvMgr().HeroNumMap[self.Sql_Hero.BuyPosNum+1]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_BUY_MAX"))
		return
	}

	// 检查消耗是否正常
	err := self.player.HasObjectOkEasy(config.BuyNeed, config.BuyNeedNum)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	self.Sql_Hero.BuyPosNum += 1
	items := self.player.RemoveObjectSimple(config.BuyNeed, config.BuyNeedNum, "购买英雄位置", 0, 0, 0)

	var msgRel S2C_BuyPos
	msgRel.Cid = "herobuypos"
	msgRel.BuyPosNum = self.Sql_Hero.BuyPosNum
	msgRel.CostItems = items
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModHero) CheckHeroBuyPos() bool {
	config := GetCsvMgr().HeroNumMap[self.Sql_Hero.BuyPosNum]
	if config == nil {
		return true
	}
	maxNum := config.MaxNum
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv != nil {
		maxNum += vipcsv.HeroList
	}
	if len(self.Sql_Hero.info) >= maxNum {
		return false
	}

	return true
}

// 创建一个装备
func (self *Hero) NewExclusiveEquipItem() *ExclusiveEquip {
	pItem := &ExclusiveEquip{}
	pItem.Lv = INIT_LV

	for _, v := range GetCsvMgr().ExclusiveEquipConfigMap {
		if v.HeroId == self.HeroId {
			pItem.Id = v.Id
			for i := 0; i < len(v.BaseType); i++ {
				if v.BaseType[i] > 0 {
					attr := new(AttrInfo)
					attr.AttrId = i + 1
					pItem.AttrInfo = append(pItem.AttrInfo, attr)
				}
			}
		}
	}

	return pItem
}

func (self *ModHero) CheckBackOpen(hero *Hero) {

	if self.Sql_Hero.BackOpen == LOGIC_TRUE || hero.StarItem.UpStar <= BACKOPEN_LIMIT_STAR {
		return
	}
	//看初始星级是否
	if GetCsvMgr().GetHeroInitLv(hero.HeroId) < BACKOPEN_LIMIT_SPECIAL_STAR {
		return
	}

	count := 0
	for _, v := range self.Sql_Hero.info {
		if v.StarItem.UpStar > BACKOPEN_LIMIT_STAR && v.HeroId == hero.HeroId {
			count++
			if count >= BACKOPEN_LIMIT_NUM {
				self.Sql_Hero.BackOpen = LOGIC_TRUE

				var msgRel S2C_HeroBackOpen
				msgRel.Cid = "herobackopen"
				msgRel.BackOpen = self.Sql_Hero.BackOpen
				self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
				return
			}
		}
	}
}

func (self *ModHero) CheckHireUpdate(hero *Hero) {
	if hero.StarItem.UpStar < HIRE_STAR_LIMIT {
		return
	}

	info := GetHireHeroInfoMgr().GetInfo(self.player.Sql_UserBase.Uid)
	if info == nil {
		return
	}

	_, ok := info.hireHeroInfo[hero.HeroKeyId]
	if !ok {
		info.hireHeroInfo[hero.HeroKeyId] = self.NewHireHero(hero)
	} else {
		info.hireHeroInfo[hero.HeroKeyId].HeroQuality = hero.StarItem.UpStar
		//神器
		artifact := self.player.GetModule("artifactequip").(*ModArtifactEquip).GetArtifactEquipItem(hero.ArtifactEquipIds[0])
		if artifact != nil {
			info.hireHeroInfo[hero.HeroKeyId].HeroArtifactId = artifact.Id
			info.hireHeroInfo[hero.HeroKeyId].HeroArtifactLv = artifact.Lv
		} else {
			info.hireHeroInfo[hero.HeroKeyId].HeroArtifactId = 0
			info.hireHeroInfo[hero.HeroKeyId].HeroArtifactLv = 0
		}
		//专属
		if hero.ExclusiveEquip != nil {
			info.hireHeroInfo[hero.HeroKeyId].HeroExclusiveLv = hero.ExclusiveEquip.Lv
			info.hireHeroInfo[hero.HeroKeyId].HeroExclusiveUnLock = hero.ExclusiveEquip.UnLock
		}

		if hero.StageTalent != nil {
			info.hireHeroInfo[hero.HeroKeyId].Talent = new(StageTalent)
			HF_DeepCopy(info.hireHeroInfo[hero.HeroKeyId].Talent, hero.StageTalent)
		}
	}

	var msgRel S2C_HireStateUpdate
	msgRel.Cid = "hirestateupdate"
	msgRel.HireHero = info.hireHeroInfo[hero.HeroKeyId]
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	//if GetCsvMgr().IsLevelAndPassOpenNew(self.player.Sql_UserBase.Level, self.player.Sql_UserBase.PassMax, OPEN_LEVEL_HIRE) {
	GetMasterMgr().FriendPRC.AddHireFriend(self.player.GetUid(), info.hireHeroInfo[hero.HeroKeyId], LOGIC_TRUE)
	//} else {
	//	GetMasterMgr().FriendPRC.AddHireFriend(self.player.GetUid(), info.hireHeroInfo[hero.HeroKeyId], LOGIC_FALSE)
	//}

	//通知给好友其他人
	//if GetCsvMgr().IsLevelAndPassOpenNew(self.player.Sql_UserBase.Uid, OPEN_LEVEL_HIRE) {
	//	self.player.GetModule("friend").(*ModFriend).NoticeFriend(info.hireHeroInfo[hero.HeroKeyId])
	//}
}

func (self *ModHero) NewHireHero(hero *Hero) *HireHero {
	hireHero := new(HireHero)
	hireHero.ReSetTime = HF_GetNextWeekStart()
	hireHero.HeroKeyId = hero.HeroKeyId
	hireHero.HeroId = hero.HeroId
	hireHero.HeroQuality = hero.StarItem.UpStar
	//神器
	artifact := self.player.GetModule("artifactequip").(*ModArtifactEquip).GetArtifactEquipItem(hero.ArtifactEquipIds[0])
	if artifact != nil {
		hireHero.HeroArtifactId = artifact.Id
		hireHero.HeroArtifactLv = artifact.Lv
	}
	//专属
	if hero.ExclusiveEquip != nil {
		hireHero.HeroExclusiveLv = hero.ExclusiveEquip.Lv
		hireHero.HeroExclusiveUnLock = hero.ExclusiveEquip.UnLock
	}

	if hero.StageTalent != nil {
		hireHero.Talent = new(StageTalent)
		HF_DeepCopy(hireHero.Talent, hero.StageTalent)
	}

	hireHero.OwnPlayer = GetHireHeroInfoMgr().NewPlayerBase(self.player)
	return hireHero
}

func (self *ModHero) HeroLvToMax() {
	isSend := false
	for _, v := range self.Sql_Hero.info {
		if self.player.GetModule("crystal").(*ModResonanceCrystal).IsResonanceHero(v.HeroKeyId) {
			continue
		}
		if v.StarItem.UpStar < 15 {
			v.StarItem.UpStar = 15
		}
		configHero := GetCsvMgr().GetHeroMapConfig(v.HeroId, v.StarItem.UpStar)
		if configHero == nil {
			continue
		}
		config := GetCsvMgr().HeroBreakConfigMap[v.HeroId]
		if config == nil {
			continue
		}
		if v.HeroLv < configHero.HeroLvMax {
			v.LvUp(configHero.HeroLvMax - v.HeroLv)
			self.player.GetModule("support").(*ModSupportHero).UpdataHero(v.HeroKeyId)
			self.player.GetModule("crystal").(*ModResonanceCrystal).UpdatePriestsHeros(v.HeroKeyId)
			self.ChangeVoidHeroResonance(v.HeroKeyId, 0)
			self.GetAllHeroStars()
			self.player.countHeroFight(v, ReasonHeroLvUp)
			isSend = true
		}
	}

	if isSend {
		self.SendInfo()
	}
}

func (self *ModHero) HeroQuToMax() {
	isSend := false
	for _, v := range self.Sql_Hero.info {
		if v.StarItem.UpStar < 15 {
			v.StarItem.UpStar = 15
			self.player.countHeroFight(v, ReasonHeroStarUp)
			isSend = true
		}
	}
	if isSend {
		self.SendInfo()
	}
}

func (self *ModHero) HeroEquipToBest() {
	isSend := false
	//装备ID=前缀+部位+阶级+阵营
	idPre := 60000000
	idGrade := 100
	for _, v := range self.Sql_Hero.info {
		config := GetCsvMgr().HeroConfigMap[v.HeroId][1]
		if config == nil {
			continue
		}
		equipId := make([]int, 0)
		switch config.AttackType {
		case 1:
			equipId = append(equipId, idPre+21*10000+idGrade+config.Attribute)
			equipId = append(equipId, idPre+42*10000+idGrade+config.Attribute)
			equipId = append(equipId, idPre+43*10000+idGrade+config.Attribute)
			equipId = append(equipId, idPre+44*10000+idGrade+config.Attribute)
		case 2:
			equipId = append(equipId, idPre+11*10000+idGrade+config.Attribute)
			equipId = append(equipId, idPre+62*10000+idGrade+config.Attribute)
			equipId = append(equipId, idPre+63*10000+idGrade+config.Attribute)
			equipId = append(equipId, idPre+64*10000+idGrade+config.Attribute)
		case 3:
			equipId = append(equipId, idPre+31*10000+idGrade+config.Attribute)
			equipId = append(equipId, idPre+52*10000+idGrade+config.Attribute)
			equipId = append(equipId, idPre+53*10000+idGrade+config.Attribute)
			equipId = append(equipId, idPre+54*10000+idGrade+config.Attribute)
		default:
			continue
		}

		for i := 0; i < len(v.EquipIds); i++ {
			item := self.player.GetModule("equip").(*ModEquip).GetEquipItem(v.EquipIds[i])
			if item != nil && item.Id == equipId[i] {
				continue
			}
			if item != nil {
				item.HeroKeyId = 0
			}

			equip := self.player.GetModule("equip").(*ModEquip).NewSuperEquipItem(equipId[i])
			if equip != nil {
				equip.HeroKeyId = v.HeroKeyId
				v.EquipIds[i] = equip.KeyId
				self.player.countHeroFight(v, ReasonEquipWear)
				isSend = true
			}
		}
	}
	//先装备同步，后英雄同步
	if isSend {
		self.player.GetModule("equip").(*ModEquip).SendInfo()
		self.SendInfo()
	}
}

func (self *ModHero) GetHandBookStarMax(heroId int) int {
	if self.Sql_Hero.handBook == nil {
		return 0
	}
	_, ok := self.Sql_Hero.handBook[heroId]
	if ok {
		return self.Sql_Hero.handBook[heroId].StarMax
	}
	return 0
}

func (self *ModHero) GmClearHero() {
	self.Sql_Hero.info = make(map[int]*Hero)

	self.SendInfo()
}

func (self *ModHero) SendInfoSyn() {
	var msg S2C_HeroInfo
	msg.Cid = "herolstsyn"
	msg.Newhero = false
	msg.Reborn = self.Sql_Hero.Reborn
	msg.BuyPosNum = self.Sql_Hero.BuyPosNum
	msg.AutoFire = self.Sql_Hero.AutoFire
	msg.BackOpen = self.Sql_Hero.BackOpen
	for _, value := range self.Sql_Hero.info {
		if value != nil {
			value.checkStarItem(0)
			value.checkEquip()
		}
		msg.Herolst = append(msg.Herolst, value)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("herolst", smsg)
}
