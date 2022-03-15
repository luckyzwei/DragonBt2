package game

import (
	"encoding/json"
	"fmt"
	"strings"
	//"time"
)

type NewPitShop struct {
	Shoptype int               `json:"shoptype"` // 商店类型   1，2，3
	IsShow   int               `json:"isshow"`   //是否显示
	ShopGood []*JS_NewShopInfo `json:"shopgood"`
}

const (
	NEWPIT_SHOP_MAX = 3 //商店最大数量
)

//地牢商店规则对应类型，相当于配置文件
var NewPitShopType = [NEWPIT_SHOP_MAX]int{7, 8, 9}

// 地牢系统
type San_UserNewPit struct {
	Uid         int64  `json:"uid"` // 玩家Id
	NewPitInfo  string `json:"newpitinfo"`
	UserPitInfo string `json:"userpitinfo"`
	BuffStore   string `json:"buffstore"`
	WinLevel    string `json:"winlevel"`  //通过次数
	LoseLevel   string `json:"loselevel"` //通过第一关但未完整通关次数
	Shop        string `json:"shop"`      //商店

	newPitInfo  map[int]*NewPitInfo // 推图信息
	userPitInfo *UserPitInfo        //
	buffStore   map[int]map[int]int // 记录BUFF库存，辅助计算  第一个KEY是品质  第2个KEY是BUFFID
	winLevel    map[int]int         //
	loseLevel   map[int]int         //
	shop        []*NewPitShop
	DataUpdate
}

type ModNewPit struct {
	player         *Player
	Sql_UserNewPit San_UserNewPit
}

type NewPitInfo struct {
	Id             int                      `json:"id"`    //关卡唯一标识
	State          int                      `json:"state"` //关卡状态
	Element        int                      `json:"element"`
	FightInfoId    int64                    `json:"fightinfoid"`
	HeroState      map[int]*NewPitHeroState `json:"herostate"`      //保存英雄血量和能量
	Rewards        []PassItem               `json:"rewards"`        //奖励
	NewPitShopInfo []*NewPitGoodsInfo       `json:"newpitshopinfo"` //商店信息
	NewPitCartHero []*NewHero               `json:"newpitcarthero"` //马车的英雄信息
}

type NewPitHeroState struct {
	HeroKeyId int `json:"herokeyid"` //
	Hp        int `json:"hp"`        //血量
	Energy    int `json:"energy"`    //能量
}

type UserPitInfo struct {
	NowPitId    int                      `json:"nowpitid"`    //当前位置
	NowAimId    int                      `json:"nowaimid"`    //当前选择
	EndTime     int64                    `json:"endtime"`     //结束时间
	Buff        []int                    `json:"buff"`        //累计有哪些BUFF
	BuffChoose  []int                    `json:"buffchoose"`  //这个有值的时候需要选择BUFF才可以继续
	HeroState   map[int]*NewPitHeroState `json:"herostate"`   //保存英雄血量和能量
	GetHeroList []*NewHero               `json:"getherolist"` //获得英雄
	IsFirst     int                      `json:"isfirst"`     //是否首次
	BuffNew     map[int]int              `json:"buffnew"`     //新版地牢辅助计算,之前的结构涉及客户端计算所以不能更改
	IsFinish    int                      `json:"isfinish"`    //
}

type NewPitGoodsInfo struct {
	GoodId   int `json:"goodid"`
	GoodNum  int `json:"goodnum"`
	CostId   int `json:"costid"`
	CostNum  int `json:"costnum"`
	Discount int `json:"discount"`
	State    int `json:"state"`
}

const (
	NEWPIT_STATE_CANT_FINISH = 0 //不能完成
	NEWPIT_STATE_CAN_FINISH  = 1 //能完成
	NEWPIT_STATE_FINISHED    = 2 //已完成
)

const (
	NEWPIT_DIFFICULT_NORMAL = 1 //普通难度
	NEWPIT_DIFFICULT_HARD   = 2 //困难难度
)

const (
	NEWPIT_BUFF_SPRING_LOW  = 93009101
	NEWPIT_BUFF_SPRING_HIGH = 93009201
)

const (
	NEWPIT_TASK_START          = 1  //起点
	NEWPIT_TASK_MONSTER_NORMAL = 2  //普通怪
	NEWPIT_TASK_MONSTER_HARD   = 3  //精英怪
	NEWPIT_TASK_BOSS           = 4  //BOSS
	NEWPIT_TASK_SHOP           = 5  //商人
	NEWPIT_TASK_SOUL_CART      = 6  //灵魂马车
	NEWPIT_TASK_SPRING         = 7  //温泉
	NEWPIT_TASK_MYSTERY        = 8  //神秘人
	NEWPIT_TASK_EVIL_CART      = 9  //邪恶马车
	NEWPIT_TASK_TREASURE       = 10 //宝藏守卫
	NEWPIT_TASK_REWARD         = 11 //通关奖励
)

const (
	NEWPIT_TIME_SYSTEM    = 26 //时间ID
	NEWPIT_MAX_STAR_ROBOT = 37 //马车最大星级
)

func (self *ModNewPit) OnGetData(player *Player) {
	self.player = player
	sql := fmt.Sprintf("select * from `san_usernewpit` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_UserNewPit, "san_usernewpit", self.player.ID)

	if self.Sql_UserNewPit.Uid <= 0 {
		self.Sql_UserNewPit.Uid = self.player.ID
		self.Sql_UserNewPit.newPitInfo = make(map[int]*NewPitInfo)
		self.Sql_UserNewPit.userPitInfo = new(UserPitInfo)
		self.Sql_UserNewPit.buffStore = make(map[int]map[int]int)
		self.Sql_UserNewPit.winLevel = make(map[int]int)
		self.Sql_UserNewPit.loseLevel = make(map[int]int)
		self.Encode()
		InsertTable("san_usernewpit", &self.Sql_UserNewPit, 0, true)
		self.Sql_UserNewPit.Init("san_usernewpit", &self.Sql_UserNewPit, true)
	} else {
		self.Decode()
		self.Sql_UserNewPit.Init("san_usernewpit", &self.Sql_UserNewPit, true)
	}

}

func (self *ModNewPit) OnSave(sql bool) {
	self.Encode()
	self.Sql_UserNewPit.Update(sql)
}

func (self *ModNewPit) OnGetOtherData() {

}

func (self *ModNewPit) Decode() { // 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_UserNewPit.NewPitInfo), &self.Sql_UserNewPit.newPitInfo)
	json.Unmarshal([]byte(self.Sql_UserNewPit.UserPitInfo), &self.Sql_UserNewPit.userPitInfo)
	json.Unmarshal([]byte(self.Sql_UserNewPit.BuffStore), &self.Sql_UserNewPit.buffStore)
	json.Unmarshal([]byte(self.Sql_UserNewPit.WinLevel), &self.Sql_UserNewPit.winLevel)
	json.Unmarshal([]byte(self.Sql_UserNewPit.LoseLevel), &self.Sql_UserNewPit.loseLevel)
	json.Unmarshal([]byte(self.Sql_UserNewPit.Shop), &self.Sql_UserNewPit.shop)
}

func (self *ModNewPit) Encode() { // 将data数据写入数据库
	self.Sql_UserNewPit.NewPitInfo = HF_JtoA(&self.Sql_UserNewPit.newPitInfo)
	self.Sql_UserNewPit.UserPitInfo = HF_JtoA(&self.Sql_UserNewPit.userPitInfo)
	self.Sql_UserNewPit.BuffStore = HF_JtoA(&self.Sql_UserNewPit.buffStore)
	self.Sql_UserNewPit.WinLevel = HF_JtoA(&self.Sql_UserNewPit.winLevel)
	self.Sql_UserNewPit.LoseLevel = HF_JtoA(&self.Sql_UserNewPit.loseLevel)
	self.Sql_UserNewPit.Shop = HF_JtoA(&self.Sql_UserNewPit.shop)
}

func (self *ModNewPit) Check() bool { // 将data数据写入数据库
	if self.Sql_UserNewPit.newPitInfo == nil {
		self.Sql_UserNewPit.newPitInfo = make(map[int]*NewPitInfo, 0)
	}

	if self.Sql_UserNewPit.userPitInfo == nil {
		self.Sql_UserNewPit.userPitInfo = new(UserPitInfo)
	}

	if self.Sql_UserNewPit.buffStore == nil {
		self.Sql_UserNewPit.buffStore = make(map[int]map[int]int)
	}

	if self.Sql_UserNewPit.winLevel == nil {
		self.Sql_UserNewPit.winLevel = make(map[int]int)
	}

	if self.Sql_UserNewPit.loseLevel == nil {
		self.Sql_UserNewPit.loseLevel = make(map[int]int)
	}

	//商店初始化
	size := len(self.Sql_UserNewPit.shop)
	for i := size; i < NEWPIT_SHOP_MAX; i++ {
		shop := new(NewPitShop)
		shop.Shoptype = i + 1
		shop.IsShow = LOGIC_FALSE
		stage, _ := GetOfflineInfoMgr().GetBaseInfo(self.player.Sql_UserBase.Uid)
		shop.ShopGood = GetCsvMgr().MakeNewShop(self.GetShopRefreshType(shop.Shoptype), stage, self.player)
		self.Sql_UserNewPit.shop = append(self.Sql_UserNewPit.shop, shop)
	}

	//检查数据
	rel := false
	now := TimeServer().Unix()
	if self.Sql_UserNewPit.userPitInfo.EndTime <= now {
		self.ResetNewPit()
		rel = true
	}

	isOpen, index := GetActivityMgr().JudgeOpenIndex(ACT_NEWPIT_HALF_MIN, ACT_NEWPIT_HALF_MAX)
	if isOpen {
		if self.player.GetModule("activity").(*ModActivity).IsActivityOpen(index) {
			if self.Sql_UserNewPit.userPitInfo.EndTime-DAY_SECS > now {
				self.Sql_UserNewPit.userPitInfo.EndTime = self.Sql_UserNewPit.userPitInfo.EndTime - DAY_SECS
				rel = true
			}
		}
	}
	return rel
}

func (self *ModNewPit) ResetNewPit() {
	//判断首次生成
	isFirst := LOGIC_FALSE
	if self.Sql_UserNewPit.userPitInfo == nil || self.Sql_UserNewPit.userPitInfo.EndTime == 0 {
		isFirst = LOGIC_TRUE
	} else {
		//看看上次的进度计算等级
		config := GetCsvMgr().NewPitConfigMap[self.Sql_UserNewPit.userPitInfo.NowPitId]
		if config != nil {
			//如果通关了
			if self.Sql_UserNewPit.userPitInfo.IsFinish == LOGIC_TRUE {
				self.Sql_UserNewPit.winLevel[config.Difficulty]++
			} else if config.NumberOfPlies > 1 {
				//通过了第一层，但是没通关
				self.Sql_UserNewPit.loseLevel[config.Difficulty]++
			}
		}
	}

	//重置
	self.Sql_UserNewPit.userPitInfo = new(UserPitInfo)
	self.Sql_UserNewPit.userPitInfo.Buff = make([]int, 0)
	self.Sql_UserNewPit.userPitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
	self.Sql_UserNewPit.userPitInfo.BuffNew = make(map[int]int, 0)
	self.Sql_UserNewPit.userPitInfo.EndTime = HF_GetNewPitEnd()
	self.Sql_UserNewPit.userPitInfo.IsFirst = isFirst

	self.Sql_UserNewPit.buffStore = make(map[int]map[int]int)

	//生成BUFF库存，用于计算
	stage := self.player.GetModule("onhook").(*ModOnHook).GetStage()
	for _, v := range GetCsvMgr().NewPitRelique {
		//判断关卡要求是否满足
		if v.LevelLimit > 0 && v.LevelLimit > stage {
			continue
		}
		_, ok := self.Sql_UserNewPit.buffStore[v.Quality]
		if !ok {
			self.Sql_UserNewPit.buffStore[v.Quality] = make(map[int]int)
		}
		self.Sql_UserNewPit.buffStore[v.Quality][v.Id] = LOGIC_TRUE
	}

	//先确认用哪张图
	//mapId := HF_GetRandom(2) + 1
	mapId := HF_GetNewPitMapId()
	self.MakeMap(1, mapId, NEWPIT_DIFFICULT_NORMAL)

	stageShop, _ := GetOfflineInfoMgr().GetBaseInfo(self.player.Sql_UserBase.Uid)
	for i := 0; i < len(self.Sql_UserNewPit.shop); i++ {
		self.Sql_UserNewPit.shop[i].IsShow = LOGIC_FALSE
		self.Sql_UserNewPit.shop[i].ShopGood = GetCsvMgr().MakeNewShop(self.GetShopRefreshType(self.Sql_UserNewPit.shop[i].Shoptype), stageShop, self.player)
	}
}

func (self *ModNewPit) GmResetPit(stage int) {

	if self.Sql_UserNewPit.userPitInfo == nil || self.Sql_UserNewPit.userPitInfo.EndTime == 0 {

	} else {
		//看看上次的进度计算等级
		config := GetCsvMgr().NewPitConfigMap[self.Sql_UserNewPit.userPitInfo.NowPitId]
		if config != nil {
			//如果通关了
			if self.Sql_UserNewPit.userPitInfo.IsFinish == LOGIC_TRUE {
				self.Sql_UserNewPit.winLevel[config.Difficulty]++
			} else if config.NumberOfPlies > 1 {
				//通过了第一层，但是没通关
				self.Sql_UserNewPit.loseLevel[config.Difficulty]++
			}
		}
	}

	self.Sql_UserNewPit.userPitInfo = new(UserPitInfo)

	self.Sql_UserNewPit.userPitInfo.Buff = make([]int, 0)
	self.Sql_UserNewPit.userPitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
	self.Sql_UserNewPit.userPitInfo.BuffNew = make(map[int]int, 0)
	self.Sql_UserNewPit.userPitInfo.EndTime = HF_GetNewPitEnd()

	self.Sql_UserNewPit.buffStore = make(map[int]map[int]int)

	level := self.player.GetModule("onhook").(*ModOnHook).GetStage()
	//生成BUFF库存，用于计算
	for _, v := range GetCsvMgr().NewPitRelique {
		//判断关卡要求是否满足
		if v.LevelLimit > 0 && v.LevelLimit > level {
			continue
		}
		_, ok := self.Sql_UserNewPit.buffStore[v.Quality]
		if !ok {
			self.Sql_UserNewPit.buffStore[v.Quality] = make(map[int]int)
		}
		self.Sql_UserNewPit.buffStore[v.Quality][v.Id] = LOGIC_TRUE
	}

	//先确认用哪张图
	mapId := HF_GetNewPitMapId()
	if stage == 1 {
		self.MakeMap(1, mapId, NEWPIT_DIFFICULT_NORMAL)
	} else {
		self.MakeMap(3, mapId, NEWPIT_DIFFICULT_HARD)
	}

	stageShop, _ := GetOfflineInfoMgr().GetBaseInfo(self.player.Sql_UserBase.Uid)
	for i := 0; i < len(self.Sql_UserNewPit.shop); i++ {
		self.Sql_UserNewPit.shop[i].IsShow = LOGIC_FALSE
		self.Sql_UserNewPit.shop[i].ShopGood = GetCsvMgr().MakeNewShop(self.GetShopRefreshType(self.Sql_UserNewPit.shop[i].Shoptype), stageShop, self.player)
	}

	self.SendInfo()
}

func (self *ModNewPit) MakeMap(numberOfPlies int, mapId int, difficulty int) {
	self.Sql_UserNewPit.newPitInfo = make(map[int]*NewPitInfo, 0)
	fight := GetOfflineInfoMgr().GetMaxFight(self.player.Sql_UserBase.Uid)
	ExclusiveLv, ExclusiveNum := self.player.GetModule("crystal").(*ModResonanceCrystal).CalExclusiveLvMax()
	//AlreadyPlayer := make(map[string]int)
	for _, v := range GetCsvMgr().NewPitConfigMap {
		if v.MapId == mapId && v.NumberOfPlies == numberOfPlies && v.Difficulty == difficulty {
			newPit := new(NewPitInfo)
			newPit.Id = v.Id
			newPit.State = NEWPIT_STATE_CAN_FINISH
			newPit.Element = v.Element
			if newPit.Element == NEWPIT_TASK_TREASURE {
				newPit.Rewards = HF_RewardCsvToPassItem(v.TreasureItem, v.TreasureNum)
			} else {
				newPit.Rewards = HF_RewardCsvToPassItem(v.Prize, v.Num)
				if v.BraveItem > 0 && self.player.GetModule("recharge").(*ModRecharge).WarOrderIsOpen(WARORDER_2) {
					newPit.Rewards = append(newPit.Rewards, PassItem{v.BraveItem, v.BraveNum})
				}
			}
			newPit.HeroState = make(map[int]*NewPitHeroState)

			upperParam := v.UpperParam
			limitParam := v.LimitParam
			limitRobot := v.LimitRobot
			minimunLinit := v.MinimumLimit

			//计算动态变更
			configLevel := GetCsvMgr().GetNewPitDifficulty(v.Difficulty, v.NumberOfPlies)
			if configLevel != nil {
				upperParam = upperParam + self.Sql_UserNewPit.winLevel[v.Difficulty]*configLevel.UpperParamAdd + self.Sql_UserNewPit.loseLevel[v.Difficulty]*configLevel.UpperParamReduction
				if upperParam < v.UpperMin {
					upperParam = v.UpperMin
				} else if upperParam > v.UpperParam {
					upperParam = v.UpperParam
				}

				limitParam = limitParam + self.Sql_UserNewPit.winLevel[v.Difficulty]*configLevel.LimitParamAdd + self.Sql_UserNewPit.loseLevel[v.Difficulty]*configLevel.LimitParamReduction
				if limitParam < v.LimitMin {
					limitParam = v.LimitMin
				} else if limitParam > v.LimitParam {
					limitParam = v.LimitParam
				}

				limitRobot = limitRobot + int64(self.Sql_UserNewPit.winLevel[v.Difficulty])*configLevel.LimitRobotAdd + int64(self.Sql_UserNewPit.loseLevel[v.Difficulty])*configLevel.LimitRobotReduction
				if limitRobot < v.LimitRobotMin {
					limitRobot = v.LimitRobotMin
				} else if limitRobot > v.LimitRobot {
					limitRobot = v.LimitRobot
				}

				minimunLinit = minimunLinit + self.Sql_UserNewPit.winLevel[v.Difficulty]*configLevel.MinimumLimitAdd + self.Sql_UserNewPit.loseLevel[v.Difficulty]*configLevel.MinimumLimitReduction
				if minimunLinit < v.MinimumMin {
					minimunLinit = v.MinimumMin
				} else if minimunLinit > v.MinimumLimit {
					minimunLinit = v.MinimumLimit
				}
			}

			maxFight := float64(fight) * float64(upperParam) / float64(PER_BIT)
			minFight := float64(fight) * float64(limitParam) / float64(PER_BIT)

			dynamicTimeStar := make([]int, 0)
			dynamicTimeLv := make([]int, 0)
			switch newPit.Element {
			case NEWPIT_TASK_MONSTER_NORMAL, NEWPIT_TASK_MONSTER_HARD, NEWPIT_TASK_BOSS:
				if self.Sql_UserNewPit.userPitInfo.IsFirst == LOGIC_TRUE {
					newPit.FightInfoId = GetFightMgr().GetFightInfoID()
					temp := self.MakeBattleFightInfoFirst(v)
					for _, v := range temp.Heroinfo {
						dynamicTimeStar = append(dynamicTimeStar, v.HeroQuality)
						dynamicTimeLv = append(dynamicTimeLv, v.Levels)
					}
					HMSetRedisEx("san_newpitfightinfo", newPit.FightInfoId, &temp, DAY_SECS*2)
				} else {
					fightInfoList := GetArenaMgr().GetPlayerByFight(fight, int64(minFight), int64(maxFight), 1, int64(minimunLinit))
					newPit.FightInfoId = GetFightMgr().GetFightInfoID()
					if len(fightInfoList) > 0 {
						for _, v := range fightInfoList[0].Heroinfo {
							dynamicTimeStar = append(dynamicTimeStar, v.HeroQuality)
							dynamicTimeLv = append(dynamicTimeLv, v.Levels)
						}
						HMSetRedisEx("san_newpitfightinfo", newPit.FightInfoId, &fightInfoList[0], DAY_SECS*2)
					} else {

						fightLimt := fight * limitRobot / 5000000
						temp := self.MakeBattleFightInfo(fightLimt)
						for _, v := range temp.Heroinfo {
							dynamicTimeStar = append(dynamicTimeStar, v.HeroQuality)
							dynamicTimeLv = append(dynamicTimeLv, v.Levels)
						}
						HMSetRedisEx("san_newpitfightinfo", newPit.FightInfoId, &temp, DAY_SECS*2)
					}
				}
				//增加动态奖励
				allItems := make(map[int]*Item)
				for i := 0; i < len(dynamicTimeStar); i++ {
					dynamicConfig := GetCsvMgr().NewPitParam[dynamicTimeStar[i]]
					if dynamicConfig == nil {
						continue
					}
					for j := 0; j < len(dynamicConfig.Item); j++ {
						num := dynamicConfig.QualityParam[j] + dynamicConfig.LvParam[j]*dynamicTimeLv[i]
						AddItemMapHelper3(allItems, dynamicConfig.Item[j], num)
					}
				}
				for _, v := range allItems {
					newPit.Rewards = append(newPit.Rewards, PassItem{v.ItemId, v.ItemNum})
				}
			case NEWPIT_TASK_SOUL_CART:
				newPit.NewPitCartHero = self.MakeSoulCart(v, ExclusiveLv, ExclusiveNum)
			case NEWPIT_TASK_EVIL_CART:
				_, star, _ := self.player.GetModule("crystal").(*ModResonanceCrystal).GetPriestsPitLv()
				level, _ := GetOfflineInfoMgr().GetHeroMaxLevel(self.player.GetUid())
				cartLv := level + v.CartLv
				//测试
				if cartLv > 400 {
					cartLv = 400
				}
				cartStart := star + v.CartQuality
				if cartStart > NEWPIT_MAX_STAR_ROBOT {
					cartStart = NEWPIT_MAX_STAR_ROBOT
				}
				fightInfo := GetRobotMgr().GetRobotByPitCart(cartLv, cartStart)
				newPit.FightInfoId = GetFightMgr().GetFightInfoID()
				if fightInfo != nil {
					newPit.NewPitCartHero = self.MakeEvilCart(fightInfo, cartLv, cartStart)
					HMSetRedisEx("san_newpitfightinfo", newPit.FightInfoId, &fightInfo, DAY_SECS*2)
				} else {
					fightInfo := GetArenaMgr().GetRobot()
					temp := self.GetOne(fightInfo)
					newPit.NewPitCartHero = self.MakeEvilCart(temp, cartLv, cartStart)
					HMSetRedisEx("san_newpitfightinfo", newPit.FightInfoId, &temp, DAY_SECS*2)
				}
			case NEWPIT_TASK_TREASURE:
				//level, _, _ := self.player.GetModule("crystal").(*ModResonanceCrystal).GetPriestsPitLv()
				level, _ := GetOfflineInfoMgr().GetHeroMaxLevel(self.player.GetUid())
				levelMin := level + v.Param[0]
				levelMax := level + v.Param[1]
				levelNow := HF_GetRandom(levelMax-levelMin+1) + levelMin

				config, ok := GetCsvMgr().NewPitTreasureCave[levelNow]
				if !ok {
					//旧版本仅用来做不兼容保护
					fightInfoList := GetArenaMgr().GetPlayerByFight(fight, int64(minFight), int64(maxFight), 1, int64(minimunLinit))
					newPit.FightInfoId = GetFightMgr().GetFightInfoID()
					if len(fightInfoList) > 0 {
						temp := self.GetOne(fightInfoList[0])
						HMSetRedisEx("san_newpitfightinfo", newPit.FightInfoId, &temp, DAY_SECS*2)
					} else {
						fightInfo := GetArenaMgr().GetRobot()
						temp := self.GetOne(fightInfo)
						HMSetRedisEx("san_newpitfightinfo", newPit.FightInfoId, &temp, DAY_SECS*2)
					}
				} else {
					fightInfo := GetRobotMgr().GetRobotByMonster(config.MonsterId)
					newPit.FightInfoId = GetFightMgr().GetFightInfoID()
					if fightInfo != nil {
						temp := self.GetOne(fightInfo)
						HMSetRedisEx("san_newpitfightinfo", newPit.FightInfoId, &temp, DAY_SECS*2)
					} else {
						fightInfo := GetArenaMgr().GetRobot()
						temp := self.GetOne(fightInfo)
						HMSetRedisEx("san_newpitfightinfo", newPit.FightInfoId, &temp, DAY_SECS*2)
					}
				}

			case NEWPIT_TASK_SHOP:
				stage, _ := GetOfflineInfoMgr().GetBaseInfo(self.player.Sql_UserBase.Uid)
				lst := GetCsvMgr().MakeNewShop(SHOP_NEW_PIT_SHOP, stage, self.player)
				for i := 0; i < 4; i++ {
					if i < len(lst) {
						info := new(NewPitGoodsInfo)
						info.GoodId = lst[i].ItemId
						info.GoodNum = lst[i].ItemNum
						info.CostId = lst[i].CostId[0]
						info.CostNum = lst[i].CostNum[0]
						info.Discount = lst[i].DisCount
						newPit.NewPitShopInfo = append(newPit.NewPitShopInfo, info)
					}
				}
			}
			self.Sql_UserNewPit.newPitInfo[newPit.Id] = newPit

			//找到起点并初始化玩家信息
			if v.Element == NEWPIT_TASK_START {
				self.Sql_UserNewPit.userPitInfo.NowPitId = v.Id
			}
		}
	}
}

func (self *ModNewPit) SendInfo() {
	lastpass := self.player.GetModule("pass").(*ModPass).GetLastPass()
	passId := ONHOOK_INIT_LEVEL
	if lastpass != nil {
		passId = lastpass.Id
	}
	if !GetCsvMgr().IsLevelAndPassOpenNew(self.player.Sql_UserBase.Level, passId, OPEN_LEVEL_NEWPITINFO) {
		return
	}

	fightInfoId := GetFightMgr().GetFightInfoID()
	temp := GetArenaMgr().GetRobot()
	HMSetRedisEx("san_newpitfightinfo", fightInfoId, &temp, DAY_SECS*2)

	//更新过图相关的任务配置
	self.Check()
	var msg S2C_NewPitInfo
	msg.Cid = "newpitinfo"
	msg.NewPitInfo = self.Sql_UserNewPit.newPitInfo
	msg.UserPitInfo = self.Sql_UserNewPit.userPitInfo
	msg.UserPitInfoFightInfo = make(map[int64]*JS_FightInfo)
	msg.Shop = self.Sql_UserNewPit.shop
	for _, v := range self.Sql_UserNewPit.newPitInfo {
		if v.FightInfoId > 0 {
			fightInfo := &JS_FightInfo{}
			value, _, err := HGetRedisEx(`san_newpitfightinfo`, v.FightInfoId, fmt.Sprintf("%d", v.FightInfoId))
			if err != nil {
				v.FightInfoId = fightInfoId
				msg.UserPitInfoFightInfo[v.FightInfoId] = temp
				continue
			}
			err1 := json.Unmarshal([]byte(value), &fightInfo)
			if err1 != nil {
				v.FightInfoId = fightInfoId
				msg.UserPitInfoFightInfo[v.FightInfoId] = temp
				continue
			}
			msg.UserPitInfoFightInfo[v.FightInfoId] = fightInfo
		}
	}
	self.player.SendMsg("newpitinfo", HF_JtoB(&msg))
}

func (self *ModNewPit) CheckOpen() {
	lastpass := self.player.GetModule("pass").(*ModPass).GetLastPass()
	passId := ONHOOK_INIT_LEVEL
	if lastpass != nil {
		passId = lastpass.Id
	}
	if !GetCsvMgr().IsLevelAndPassOpenNew(self.player.Sql_UserBase.Level, passId, OPEN_LEVEL_NEWPITINFO) {
		return
	}

	//更新过图相关的任务配置
	flag := self.Check()
	if flag {
		var msg S2C_NewPitInfo
		msg.Cid = "newpitinfo"
		msg.NewPitInfo = self.Sql_UserNewPit.newPitInfo
		msg.UserPitInfo = self.Sql_UserNewPit.userPitInfo
		msg.UserPitInfoFightInfo = make(map[int64]*JS_FightInfo)
		msg.Shop = self.Sql_UserNewPit.shop
		for _, v := range self.Sql_UserNewPit.newPitInfo {
			if v.FightInfoId > 0 {
				fightInfo := &JS_FightInfo{}
				value, _, err := HGetRedisEx(`san_newpitfightinfo`, v.FightInfoId, fmt.Sprintf("%d", v.FightInfoId))
				if err != nil {
					continue
				}
				err1 := json.Unmarshal([]byte(value), &fightInfo)
				if err1 != nil {
					continue
				}
				msg.UserPitInfoFightInfo[v.FightInfoId] = fightInfo
			}
		}
		self.player.SendMsg("newpitinfo", HF_JtoB(&msg))
	}
}

func (self *ModNewPit) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "newpitinfo":
		self.SendInfo()
		return true
	}
	return false
}
func (self *ModNewPit) onReg(handlers map[string]func(body []byte)) {
	handlers["newpitfinishevent"] = self.NewPitFinishEvent
	handlers["newpitfinishnow"] = self.NewPitFinishNow
	handlers["newpitnowaim"] = self.NewPitNowAim
	handlers["newpitchoosebuff"] = self.NewPitChooseBuff
	handlers["newpituseitem"] = self.NewPitUseItem
	//具体关卡C2S_NewPitFinishBattle
	handlers["newpitfinishbattle"] = self.NewPitFinishBattle
	handlers["newpitfinishsoul"] = self.NewPitFinishSoul
	handlers["newpitfinishshop"] = self.NewPitFinishShop
	handlers["newpitfinishspring"] = self.NewPitFinishSpring
	handlers["newpitfinishmystery"] = self.NewPitFinishMystery
	handlers["newpitfinishevil"] = self.NewPitFinishEvil
	handlers["newpitfinishtreasure"] = self.NewPitFinishTreasure
	handlers["newpitfinishreward"] = self.NewPitFinishReward
	handlers["newpitshopbuy"] = self.NewPitShopBuy
}

func (self *ModNewPit) NewPitHeroState(heroKeyId int, hp int, energy int) *NewPitHeroState {
	info := new(NewPitHeroState)
	info.HeroKeyId = heroKeyId
	info.Hp = hp
	info.Energy = energy
	return info
}

//! 完成关卡   弃用  消息已拆分 待观察
func (self *ModNewPit) NewPitFinishEvent(body []byte) {
	/*
		var msg C2S_NewPitFinishEvent
		json.Unmarshal(body, &msg)

		//看是否要先选择BUFF
		if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
			return
		}
		//验证是否连续
		nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
		msgLine := self.GetLine(msg.Id)
		if msgLine-nowLine != 1 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
			return
		}

		nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
		msgRow := self.GetRow(msg.Id)
		if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
			return
		}

		config := GetCsvMgr().NewPitConfigMap[msg.Id]
		if config == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
			return
		}

		pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
		if pitInfo == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
			return
		}

		//更新一下英雄血量
		if self.Sql_UserNewPit.userPitInfo.HeroState == nil {
			self.Sql_UserNewPit.userPitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
		}
		for i := 0; i < len(msg.HeroState); i++ {
			self.Sql_UserNewPit.userPitInfo.HeroState[msg.HeroState[i].HeroKeyId] = self.NewPitHeroState(msg.HeroState[i].HeroKeyId, msg.HeroState[i].Hp, msg.HeroState[i].Energy)
		}

		if msg.IsFail == LOGIC_TRUE {
			if pitInfo.HeroState == nil {
				pitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
			}
			//更新一下怪物英雄血量
			for i := 0; i < len(msg.MonsterHeroState); i++ {
				pitInfo.HeroState[msg.MonsterHeroState[i].HeroKeyId] = self.NewPitHeroState(msg.MonsterHeroState[i].HeroKeyId, msg.MonsterHeroState[i].Hp, msg.MonsterHeroState[i].Energy)
			}
			var msgRel S2C_NewPitFinishEvent
			msgRel.Cid = "pitfinishevent"
			msgRel.Id = msg.Id
			msgRel.UserPitInfo = self.Sql_UserNewPit.userPitInfo
			msgRel.PitInfo = append(msgRel.PitInfo, pitInfo)
			self.player.SendMsg("pitfinishevent", HF_JtoB(&msgRel))
			return
		}

		//通关成功
		costItem := make([]PassItem, 0)
		var msgRel S2C_NewPitFinishEvent
		msgRel.Cid = "pitfinishevent"
		msgRel.Id = msg.Id
		switch config.Element {
		case NEWPIT_TASK_MONSTER_NORMAL, NEWPIT_TASK_MONSTER_HARD, NEWPIT_TASK_BOSS:
			//生成BUFF
			self.MakeBuff(config)
		case NEWPIT_TASK_SHOP:
			msgRel.PitInfo = self.MoveTo()
			if msg.Param > 0 && msg.Param <= len(pitInfo.NewPitShopInfo) {
				//看看物品够不够
				if err := self.player.HasObjectOkEasy(pitInfo.NewPitShopInfo[msg.Param-1].CostId, pitInfo.NewPitShopInfo[msg.Param-1].CostNum); err != nil {
					break
				}
				cost := self.player.RemoveObjectEasy(ITEM_NEW_PIT_REBORN, 1, "地牢商店", 0, 0, 0)
				msgRel.Cost = append(msgRel.Cost, cost...)
				getItem := self.player.AddObjectSimple(pitInfo.NewPitShopInfo[msg.Param-1].GoodId, pitInfo.NewPitShopInfo[msg.Param-1].GoodNum, "地牢商店", 0, 0, 0)
				msgRel.Item = append(msgRel.Item, getItem...)
			}
		case NEWPIT_TASK_SOUL_CART:
			msgRel.PitInfo = self.MoveTo2(msg.Id)
			if msg.Param > 0 && msg.Param <= len(pitInfo.NewPitCartHero) {
				msgRel.GetHero = pitInfo.NewPitCartHero[msg.Param-1]
				self.Sql_UserNewPit.userPitInfo.GetHeroList = append(self.Sql_UserNewPit.userPitInfo.GetHeroList, msgRel.GetHero)
			}
		case NEWPIT_TASK_SPRING:
			msgRel.PitInfo = self.MoveTo()
			for _, v := range self.Sql_UserNewPit.userPitInfo.HeroState {
				if v.Hp > 0 {
					v.Hp += 5000
				}
				if v.Hp > 10000 {
					v.Hp = 10000
				}
			}
		case NEWPIT_TASK_MYSTERY:
			msgRel.PitInfo = self.MoveTo()
			keyId := 0
			min := 10000
			for _, v := range self.Sql_UserNewPit.userPitInfo.HeroState {
				if v.Hp <= min {
					min = v.Hp
					keyId = v.HeroKeyId
				}
			}
			self.Sql_UserNewPit.userPitInfo.HeroState[keyId].Hp = 10000
			if self.Sql_UserNewPit.userPitInfo.HeroState[keyId].Energy < 5000 {
				self.Sql_UserNewPit.userPitInfo.HeroState[keyId].Energy = 5000
			}
		case NEWPIT_TASK_EVIL_CART:
			msgRel.PitInfo = self.MoveTo()
			if len(pitInfo.NewPitCartHero) > 0 {
				msgRel.GetHero = pitInfo.NewPitCartHero[0]
				self.Sql_UserNewPit.userPitInfo.GetHeroList = append(self.Sql_UserNewPit.userPitInfo.GetHeroList, msgRel.GetHero)
			}
		case NEWPIT_TASK_TREASURE:
			msgRel.PitInfo = self.MoveTo()
		case NEWPIT_TASK_REWARD:
			msgRel.PitInfo = self.MoveTo()
			//看看当前层奖励
			for _, v := range GetCsvMgr().NewPitBoxPrize {
				if v.MapId == config.MapId && v.Difficulty == config.Difficulty && v.NumberOfPlies == config.NumberOfPlies {
					msgRel.Item = self.player.AddObjectLst(config.Prize, config.Num, "地牢宝箱", config.NumberOfPlies, 0, 0)
					break
				}
			}
		default:
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_MSG_ERROR"))
			return
		}
		getItem := self.player.AddObjectPassItem(pitInfo.Rewards, "地牢", msg.Id, 0, 0)
		if len(getItem) > 0 {
			msgRel.Item = append(msgRel.Item, getItem...)
		}
		msgRel.Cost = costItem
		msgRel.UserPitInfo = self.Sql_UserNewPit.userPitInfo
		self.player.SendMsg("pitfinishevent", HF_JtoB(&msgRel))
		return

	*/
}

func (self *ModNewPit) NewPitFinishReward(body []byte) {
	var msg C2S_NewPitFinishReward
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}
	//验证是否连续
	nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgLine := self.GetLine(msg.Id)
	if msgLine-nowLine != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgRow := self.GetRow(msg.Id)
	if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	config := GetCsvMgr().NewPitConfigMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
		return
	}

	pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
		return
	}
	if pitInfo.State == NEWPIT_STATE_FINISHED {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_REWARD_HAS_GET"))
		return
	}

	self.player.HandleTask(TASK_TYPE_PIT_PASS, config.NumberOfPlies, 0, 0)
	if config.NumberOfPlies == 3 {
		self.Sql_UserNewPit.userPitInfo.IsFinish = LOGIC_TRUE
	}

	var msgRel S2C_NewPitFinishReward
	msgRel.Cid = "newpitfinishreward"
	msgRel.Id = msg.Id
	msgRel.PitInfo = self.MoveToReward(msg.Id)
	//增加一个额外掉落  20200521
	extLottery := self.player.GetModule("pass").(*ModPass).CalPitExt(config)
	msgRel.Item, msgRel.GetPrivilegeItems = self.GetReward(pitInfo.Rewards, "地牢通关奖励", msg.Id, 0, 0, config.Lottery, extLottery)
	msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModNewPit) NewPitShopBuy(body []byte) {
	var msg C2S_NewPitShopBuy
	json.Unmarshal(body, &msg)

	shop := self.GetShop(msg.Shoptype)
	if shop == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_NOT_EXIST"))
		return
	}

	if msg.Grid <= 0 || msg.Grid > len(shop.ShopGood) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_GRID_ERROR"))
		return
	}

	for i := 0; i < len(shop.ShopGood); i++ {
		if shop.ShopGood[i].Grid == msg.Grid {
			if shop.ShopGood[i].State == LOGIC_TRUE {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_ALREADY_BUY"))
				return
			}

			if err := self.player.HasObjectOk(shop.ShopGood[i].CostId, shop.ShopGood[i].CostNum); err != nil {
				self.player.SendErrInfo("err", err.Error())
				return
			}

			costItem := self.player.RemoveObjectLst(shop.ShopGood[i].CostId, shop.ShopGood[i].CostNum, "商店购买", msg.Shoptype, 0, 1)
			num := GetGemNum(costItem)
			param3 := 0
			if num > 0 {
				param3 = -1
			}
			getItems := self.player.AddObjectSimple(shop.ShopGood[i].ItemId, shop.ShopGood[i].ItemNum, "地牢商店购买", msg.Shoptype, 0, param3)
			//CheckAddItemLog(self.player, "商店购买", costItem, getItems)

			if num > 0 {
				AddSpecialSdkItemListLog(self.player, num, getItems, "商店购买")
			}

			//旧商店是值拷贝，这个地方需要整体引用
			shop.ShopGood[i].State = LOGIC_TRUE

			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SHOP_BUY, shop.ShopGood[i].ItemId, msg.Shoptype, shop.ShopGood[i].ItemNum, "地牢商店购买", 0, 0, self.player)

			var msgRel S2C_NewPitShopBuy
			msgRel.Cid = "newpitshopbuy"
			msgRel.GetItems = getItems
			msgRel.CostItems = costItem
			msgRel.NewShopInfo = shop.ShopGood[i]
			self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
			self.player.HandleTask(TASK_TYPE_SHOP_BUY_COUNT, 1, msg.Shoptype, 0)
			self.player.GetModule("task").(*ModTask).SendUpdate()

			for i := 0; i < len(msgRel.GetItems); i++ {
				GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SHOP_BUY_GOODS, msgRel.GetItems[i].ItemID, msg.Shoptype, 0, "地牢商店购买", 0, 0, self.player)
			}
			return
		}
	}

	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_GOODS_NOT_FIND"))
	return
}

func (self *ModNewPit) NewPitFinishSoul(body []byte) {
	var msg C2S_NewPitFinishSoul
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}
	//验证是否连续
	nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgLine := self.GetLine(msg.Id)
	if msgLine-nowLine != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgRow := self.GetRow(msg.Id)
	if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	config := GetCsvMgr().NewPitConfigMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
		return
	}

	pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
		return
	}

	//类型支持验证
	if pitInfo.Element != NEWPIT_TASK_SOUL_CART {
		return
	}

	//通关成功
	var msgRel S2C_NewPitFinishSoul
	msgRel.Cid = "newpitfinishsoul"
	msgRel.Id = msg.Id
	msgRel.PitInfo = self.MoveTo2(msg.Id)
	if msg.Param > 0 && msg.Param <= len(pitInfo.NewPitCartHero) {
		hero := pitInfo.NewPitCartHero[msg.Param-1]
		hero.HeroKeyId = self.player.GetModule("hero").(*ModHero).MaxKey()
		msgRel.GetHero = hero
		self.Sql_UserNewPit.userPitInfo.GetHeroList = append(self.Sql_UserNewPit.userPitInfo.GetHeroList, hero)
	}
	msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_SOUL_CART, pitInfo.Id, pitInfo.NewPitCartHero[msg.Param-1].HeroId, 0, "地牢号令火炬", 0, 0, self.player)

	return
}

func (self *ModNewPit) NewPitFinishEvil(body []byte) {
	var msg C2S_NewPitFinishEvil
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}
	//验证是否连续
	nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgLine := self.GetLine(msg.Id)
	if msgLine-nowLine != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgRow := self.GetRow(msg.Id)
	if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	config := GetCsvMgr().NewPitConfigMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
		return
	}

	pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
		return
	}

	//类型支持验证
	if pitInfo.Element != NEWPIT_TASK_EVIL_CART {
		return
	}

	//更新一下英雄血量
	if self.Sql_UserNewPit.userPitInfo.HeroState == nil {
		self.Sql_UserNewPit.userPitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
	}
	for i := 0; i < len(msg.HeroState); i++ {
		self.Sql_UserNewPit.userPitInfo.HeroState[msg.HeroState[i].HeroKeyId] = self.NewPitHeroState(msg.HeroState[i].HeroKeyId, msg.HeroState[i].Hp, msg.HeroState[i].Energy)
	}

	//通关成功
	var msgRel S2C_NewPitFinishEvil
	msgRel.Cid = "newpitfinishevil"
	msgRel.Id = msg.Id
	if msg.IsFail == LOGIC_TRUE {
		if pitInfo.HeroState == nil {
			pitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
		}
		//更新一下怪物英雄血量
		for i := 0; i < len(msg.MonsterHeroState); i++ {
			pitInfo.HeroState[msg.MonsterHeroState[i].HeroKeyId] = self.NewPitHeroState(msg.MonsterHeroState[i].HeroKeyId, msg.MonsterHeroState[i].Hp, msg.MonsterHeroState[i].Energy)
		}
		msgRel.PitInfo = append(msgRel.PitInfo, pitInfo)
		msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
		msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
		msgRel.HeroState = self.Sql_UserNewPit.userPitInfo.HeroState
		self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
		return
	}
	msgRel.PitInfo = self.MoveTo2(msg.Id)
	if len(pitInfo.NewPitCartHero) > 0 {
		pitInfo.NewPitCartHero[0].HeroKeyId = self.player.GetModule("hero").(*ModHero).MaxKey()
		msgRel.GetHero = pitInfo.NewPitCartHero[0]
		self.Sql_UserNewPit.userPitInfo.GetHeroList = append(self.Sql_UserNewPit.userPitInfo.GetHeroList, msgRel.GetHero)
	}
	msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	msgRel.HeroState = self.Sql_UserNewPit.userPitInfo.HeroState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_EVIL_CART, pitInfo.Id, 0, 0, "地牢试炼火炬", 0, 0, self.player)

	return
}

func (self *ModNewPit) NewPitFinishTreasure(body []byte) {
	var msg C2S_NewPitFinishTreasure
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}
	//验证是否连续
	nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgLine := self.GetLine(msg.Id)
	if msgLine-nowLine != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgRow := self.GetRow(msg.Id)
	if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	config := GetCsvMgr().NewPitConfigMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
		return
	}

	pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
		return
	}

	//类型支持验证
	if pitInfo.Element != NEWPIT_TASK_TREASURE {
		return
	}

	//更新一下英雄血量
	deadNum := 0
	if self.Sql_UserNewPit.userPitInfo.HeroState == nil {
		self.Sql_UserNewPit.userPitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
	}
	for i := 0; i < len(msg.HeroState); i++ {
		self.Sql_UserNewPit.userPitInfo.HeroState[msg.HeroState[i].HeroKeyId] = self.NewPitHeroState(msg.HeroState[i].HeroKeyId, msg.HeroState[i].Hp, msg.HeroState[i].Energy)
		if msg.HeroState[i].Hp == 0 {
			deadNum++
		}
	}

	//通关成功
	var msgRel S2C_NewPitFinishTreasure
	msgRel.Cid = "newpitfinishtreasure"
	msgRel.Id = msg.Id
	msgRel.PitInfo = self.MoveTo2(msg.Id)
	msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	if msg.IsFail == LOGIC_FALSE {
		//成功才发奖励
		if len(pitInfo.Rewards) > 0 {
			index := HF_GetRandom(len(pitInfo.Rewards))
			item := make([]PassItem, 0)
			item = append(item, pitInfo.Rewards[index])
			msgRel.Item = self.player.AddObjectPassItem(item, "地牢", msg.Id, 0, 0)
		}
	}
	msgRel.HeroState = self.Sql_UserNewPit.userPitInfo.HeroState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_TREASURE, pitInfo.Id, msg.IsFail, deadNum, "地牢秘宝守卫", 0, 0, self.player)

	return
}

func (self *ModNewPit) NewPitFinishShop(body []byte) {
	var msg C2S_NewPitFinishShop
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}
	//验证是否连续
	nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgLine := self.GetLine(msg.Id)
	if msgLine-nowLine != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgRow := self.GetRow(msg.Id)
	if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	config := GetCsvMgr().NewPitConfigMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
		return
	}

	pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
		return
	}

	//类型支持验证
	if pitInfo.Element != NEWPIT_TASK_SHOP {
		return
	}
	//通关成功
	var msgRel S2C_NewPitFinishShop
	msgRel.Cid = "newpitfinishshop"
	if msg.Param > 0 && msg.Param <= len(pitInfo.NewPitShopInfo) {
		if pitInfo.NewPitShopInfo[msg.Param-1].State == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_SHOP_IS_BUY"))
			return
		}
		self.Sql_UserNewPit.userPitInfo.NowAimId = msg.Id
		//看看物品够不够
		if err := self.player.HasObjectOkEasy(pitInfo.NewPitShopInfo[msg.Param-1].CostId, pitInfo.NewPitShopInfo[msg.Param-1].CostNum); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
		cost := self.player.RemoveObjectEasy(pitInfo.NewPitShopInfo[msg.Param-1].CostId, pitInfo.NewPitShopInfo[msg.Param-1].CostNum, "地牢商店", 0, 0, 0)
		msgRel.Cost = append(msgRel.Cost, cost...)
		getItem := self.player.AddObjectSimple(pitInfo.NewPitShopInfo[msg.Param-1].GoodId, pitInfo.NewPitShopInfo[msg.Param-1].GoodNum, "地牢商店", 0, 0, 0)
		msgRel.Item = append(msgRel.Item, getItem...)
		pitInfo.NewPitShopInfo[msg.Param-1].State = LOGIC_TRUE
		msgRel.PitInfo = append(msgRel.PitInfo, pitInfo)
	} else {
		msgRel.PitInfo = self.MoveTo()
	}
	msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModNewPit) NewPitFinishSpring(body []byte) {
	var msg C2S_NewPitFinishSpring
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}
	//验证是否连续
	nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgLine := self.GetLine(msg.Id)
	if msgLine-nowLine != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgRow := self.GetRow(msg.Id)
	if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	config := GetCsvMgr().NewPitConfigMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
		return
	}

	pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
		return
	}

	//类型支持验证
	if pitInfo.Element != NEWPIT_TASK_SPRING {
		return
	}
	//通关成功
	var msgRel S2C_NewPitFinishSpring
	msgRel.Cid = "newpitfinishspring"
	msgRel.Id = msg.Id
	msgRel.PitInfo = self.MoveTo()
	valueHp := 5000
	valueMp := 0

	for i := 0; i < len(self.Sql_UserNewPit.userPitInfo.Buff); i++ {
		if self.Sql_UserNewPit.userPitInfo.Buff[i] == NEWPIT_BUFF_SPRING_LOW {
			valueHp += 2500
			valueMp += 2500
		} else if self.Sql_UserNewPit.userPitInfo.Buff[i] == NEWPIT_BUFF_SPRING_HIGH {
			valueHp += 5000
			valueMp += 5000
		}
	}

	for _, v := range self.Sql_UserNewPit.userPitInfo.HeroState {
		if v.Hp > 0 {
			v.Hp += valueHp
			v.Energy += valueMp
		}
		if v.Hp > 10000 {
			v.Hp = 10000
		}
		if v.Energy > 10000 {
			v.Energy = 10000
		}
	}
	msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	msgRel.HeroState = self.Sql_UserNewPit.userPitInfo.HeroState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_STRING, pitInfo.Id, 0, 0, "地牢复苏清泉", 0, 0, self.player)
	return
}

func (self *ModNewPit) NewPitFinishMystery(body []byte) {
	var msg C2S_NewPitFinishMystery
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}
	//验证是否连续
	nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgLine := self.GetLine(msg.Id)
	if msgLine-nowLine != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgRow := self.GetRow(msg.Id)
	if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	config := GetCsvMgr().NewPitConfigMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
		return
	}

	pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
		return
	}

	//类型支持验证
	if pitInfo.Element != NEWPIT_TASK_MYSTERY {
		return
	}
	//通关成功
	var msgRel S2C_NewPitFinishMystery
	msgRel.Cid = "newpitfinishmystery"
	msgRel.Id = msg.Id
	msgRel.PitInfo = self.MoveTo()
	keyId := 0
	min := 10000
	for _, v := range self.Sql_UserNewPit.userPitInfo.HeroState {
		if v.Hp <= min {
			min = v.Hp
			keyId = v.HeroKeyId
		}
	}
	if keyId > 0 {
		self.Sql_UserNewPit.userPitInfo.HeroState[keyId].Hp = 10000
		self.Sql_UserNewPit.userPitInfo.HeroState[keyId].Energy += 5000
		if self.Sql_UserNewPit.userPitInfo.HeroState[keyId].Energy > 10000 {
			self.Sql_UserNewPit.userPitInfo.HeroState[keyId].Energy = 10000
		}
	}
	msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	msgRel.HeroState = self.Sql_UserNewPit.userPitInfo.HeroState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_REBORN, pitInfo.Id, keyId, 0, "地牢神使", 0, 0, self.player)
	return
}

func (self *ModNewPit) NewPitFinishBattle(body []byte) {
	var msg C2S_NewPitFinishBattle
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}
	//验证是否连续
	nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgLine := self.GetLine(msg.Id)
	if msgLine-nowLine != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgRow := self.GetRow(msg.Id)
	if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	config := GetCsvMgr().NewPitConfigMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
		return
	}

	pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
		return
	}

	//类型支持验证
	if pitInfo.Element != NEWPIT_TASK_MONSTER_NORMAL && pitInfo.Element != NEWPIT_TASK_MONSTER_HARD && pitInfo.Element != NEWPIT_TASK_BOSS {
		return
	}

	//更新一下英雄血量
	if self.Sql_UserNewPit.userPitInfo.HeroState == nil {
		self.Sql_UserNewPit.userPitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
	}
	deadNum := 0
	for i := 0; i < len(msg.HeroState); i++ {
		self.Sql_UserNewPit.userPitInfo.HeroState[msg.HeroState[i].HeroKeyId] = self.NewPitHeroState(msg.HeroState[i].HeroKeyId, msg.HeroState[i].Hp, msg.HeroState[i].Energy)
		if msg.HeroState[i].Hp == 0 {
			deadNum++
		}
	}
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_BATTLE, pitInfo.Id, config.Id, config.Difficulty, "地牢战斗", 0, deadNum, self.player)

	var msgRel S2C_NewPitFinishBattle
	msgRel.Cid = "newpitfinishbattle"
	msgRel.Id = msg.Id

	if msg.IsFail == LOGIC_TRUE {
		if pitInfo.HeroState == nil {
			pitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
		}
		//更新一下怪物英雄血量
		for i := 0; i < len(msg.MonsterHeroState); i++ {
			pitInfo.HeroState[msg.MonsterHeroState[i].HeroKeyId] = self.NewPitHeroState(msg.MonsterHeroState[i].HeroKeyId, msg.MonsterHeroState[i].Hp, msg.MonsterHeroState[i].Energy)
		}
		msgRel.PitInfo = append(msgRel.PitInfo, pitInfo)
		msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
		msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
		msgRel.HeroState = self.Sql_UserNewPit.userPitInfo.HeroState
		self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
		return
	}

	if pitInfo.Element == NEWPIT_TASK_BOSS {
		self.player.HandleTask(TASK_TYPE_PIT_KILL, config.NumberOfPlies, 0, 0)
	} else if pitInfo.Element == NEWPIT_TASK_MONSTER_NORMAL || pitInfo.Element == NEWPIT_TASK_MONSTER_HARD {
		self.player.HandleTask(TASK_TYPE_PIT_KILL_PLAYER, config.NumberOfPlies, 0, 0)
	}

	//生成BUFF  如果是第三层BOSS 则不生成
	if config.NumberOfPlies == 3 && config.Element == NEWPIT_TASK_BOSS {
		msgRel.PitInfo = self.MoveTo()
	} else {
		self.MakeBuff(config)
		msgRel.PitInfo = append(msgRel.PitInfo, pitInfo)
	}
	msgRel.Item, msgRel.GetPrivilegeItems = self.GetReward(pitInfo.Rewards, "地牢", msg.Id, 0, 0, config.Lottery, 0)
	msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	msgRel.BuffChoose = self.Sql_UserNewPit.userPitInfo.BuffChoose
	msgRel.HeroState = self.Sql_UserNewPit.userPitInfo.HeroState

	lastpass := self.player.GetModule("pass").(*ModPass).GetLastPass()
	passId := ONHOOK_INIT_LEVEL
	if lastpass != nil {
		passId = lastpass.Id
	}
	if GetCsvMgr().IsLevelAndPassOpenNew(self.player.Sql_UserBase.Level, passId, OPEN_LEVEL_NEWPITSHOP) {
		//生成商店
		msgRel.Shop = make([]*NewPitShop, 0)
		weight := 0
		for i := 0; i < len(config.VShop_p); i++ {
			weight += config.VShop_p[i]
		}
		if weight > 0 {
			pro := HF_GetRandom(weight)
			cur := 0
			for i := 0; i < len(config.VShop_p); i++ {
				cur += config.VShop_p[i]
				if pro < cur {
					//找当前的商店
					isFind := false
					for _, shop := range self.Sql_UserNewPit.shop {
						if shop.Shoptype != i+1 {
							continue
						}
						if shop.IsShow == LOGIC_FALSE {
							shop.IsShow = LOGIC_TRUE
							msgRel.Shop = append(msgRel.Shop, shop)
						}
						isFind = true
						break
					}
					if isFind {
						break
					}
				}
			}
		}
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModNewPit) GetReward(passItem []PassItem, reason string, param1, param2, param3 int, lottery []int, extLottery int) ([]PassItem, []PassItem) {

	privilegeItems := make(map[int]*Item)
	//计算VIP加成
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		items := self.player.AddObjectPassItem(passItem, reason, param1, 0, 0)
		return items, make([]PassItem, 0)
	}

	realItems := make(map[int]*Item)
	value := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValue(PRIVILEGE_PIT)
	for i := 0; i < len(passItem); i++ {
		switch passItem[i].ItemID {
		case ITEM_NEW_PIT_COIN:
			rate := float32(PER_BIT+vipcsv.MazeFateGold) / float32(PER_BIT)
			num := int(float32(passItem[i].Num) * rate)
			AddItemMapHelper3(realItems, passItem[i].ItemID, num)

			if value > 0 {
				pNum := passItem[i].Num * value / 100
				AddItemMapHelper3(privilegeItems, ITEM_NEW_PIT_COIN, pNum)
			}
		case ITEM_GOLD:
			rate := float32(PER_BIT+vipcsv.MazeGold) / float32(PER_BIT)
			num := int(float32(passItem[i].Num) * rate)
			AddItemMapHelper3(realItems, passItem[i].ItemID, num)
		default:
			AddItemMapHelper3(realItems, passItem[i].ItemID, passItem[i].Num)
		}
	}
	for i := 0; i < len(lottery); i++ {
		if lottery[i] == 0 {
			continue
		}
		item := GetLootMgr().LootItem(lottery[i], self.player)
		AddItemMap(realItems, item)
	}
	if extLottery > 0 {
		extConfig := GetCsvMgr().NewPitExtraReward[extLottery]
		if extConfig != nil {
			for i := 0; i < len(extConfig.Lottery); i++ {
				if extConfig.Lottery[i] == 0 {
					continue
				}
				item := GetLootMgr().LootItem(extConfig.Lottery[i], self.player)
				AddItemMap(realItems, item)
			}
		}
	}
	items := self.player.AddObjectItemMap(realItems, reason, param1, param2, param3)
	if len(privilegeItems) > 0 {
		pItems := self.player.AddObjectItemMap(privilegeItems, reason, param1, value, param3)
		return items, pItems
	} else {
		return items, nil
	}
}

func (self *ModNewPit) MakeBuff(config *NewPitConfig) {
	//计算品质
	allRate := 0
	for i := 0; i < len(config.Relique); i++ {
		allRate += config.Relique[i]
	}
	for times := 0; times < 3; times++ {
		rand := HF_GetRandom(allRate)
		rate := 0
		for i := 0; i < len(config.Relique); i++ {
			rate += config.Relique[i]
			if rate >= rand {
				realIndex := i + 1
				for k, _ := range self.Sql_UserNewPit.buffStore[realIndex] {
					//验证身上存在这个BUFF时候不重复出现
					buffConfig := GetCsvMgr().NewPitRelique[k]
					if buffConfig == nil {
						continue
					}
					if buffConfig.OnlyOwned == LOGIC_TRUE {
						_, ok := self.Sql_UserNewPit.userPitInfo.BuffNew[k]
						if ok {
							continue
						}
					}
					//验证身上存在这个BUFF不能已经是3个之一
					isCan := true
					for j := 0; j < len(self.Sql_UserNewPit.userPitInfo.BuffChoose); j++ {
						if self.Sql_UserNewPit.userPitInfo.BuffChoose[j] == k {
							isCan = false
							break
						}
					}
					if !isCan {
						continue
					}

					self.Sql_UserNewPit.userPitInfo.BuffChoose = append(self.Sql_UserNewPit.userPitInfo.BuffChoose, k)
					break
				}
				break
			}
		}
	}

	size := len(self.Sql_UserNewPit.userPitInfo.BuffChoose)
	for i := size; i < 3; i++ {
		for _, v := range self.Sql_UserNewPit.buffStore {
			for kk, _ := range v {
				//验证身上存在这个BUFF时候不重复出现
				buffConfig := GetCsvMgr().NewPitRelique[kk]
				if buffConfig == nil {
					continue
				}
				if buffConfig.OnlyOwned == LOGIC_TRUE {
					_, ok := self.Sql_UserNewPit.userPitInfo.BuffNew[kk]
					if ok {
						continue
					}
				}
				//验证身上存在这个BUFF不能已经是3个之一
				isCan := true
				for j := 0; j < len(self.Sql_UserNewPit.userPitInfo.BuffChoose); j++ {
					if self.Sql_UserNewPit.userPitInfo.BuffChoose[j] == kk {
						isCan = false
						break
					}
				}
				if !isCan {
					continue
				}
				self.Sql_UserNewPit.userPitInfo.BuffChoose = append(self.Sql_UserNewPit.userPitInfo.BuffChoose, kk)
				break
			}
		}
	}
}
func (self *ModNewPit) NewPitFinishNow(body []byte) {
	var msg C2S_NewPitFinishNow
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}

	//看当前层是否通过
	pitInfo := self.Sql_UserNewPit.newPitInfo[self.Sql_UserNewPit.userPitInfo.NowPitId]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_INFO_NIL"))
		return
	}

	if pitInfo.Element != NEWPIT_TASK_BOSS {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_THIS_NOT_FINISH"))
		return
	}

	//看看有没下层
	config := GetCsvMgr().NewPitConfigMap[pitInfo.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_ERROR"))
		return
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_REWARD, config.NumberOfPlies, config.Difficulty, 0, "领取通关地牢整层奖励", 0, 0, self.player)

	if config.NumberOfPlies < 3 {
		//生成新的一层
		self.MakeMap(config.NumberOfPlies+1, config.MapId, msg.Difficult)

		if config.NumberOfPlies == 2 {
			if msg.Difficult == 1 {
				GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_NEW_EASY, 0, 0, 0, "进入地牢简单模式", 0, 0, self.player)
			} else {
				GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_NEW_HARD, 0, 0, 0, "进入地牢困难模式", 0, 0, self.player)
			}
		}
	}

	var msgRel S2C_NewPitFinishNow
	msgRel.Cid = "newpitfinishnow"
	msgRel.PitInfo = self.Sql_UserNewPit.newPitInfo
	msgRel.UserPitInfo = self.Sql_UserNewPit.userPitInfo
	msgRel.UserPitInfoFightInfo = make(map[int64]*JS_FightInfo)
	for _, v := range self.Sql_UserNewPit.newPitInfo {
		if v.FightInfoId > 0 {
			fightInfo := &JS_FightInfo{}
			value, _, err := HGetRedisEx(`san_newpitfightinfo`, v.FightInfoId, fmt.Sprintf("%d", v.FightInfoId))
			if err != nil {
				continue
			}
			err1 := json.Unmarshal([]byte(value), &fightInfo)
			if err1 == nil {
				msgRel.UserPitInfoFightInfo[v.FightInfoId] = fightInfo
			}
		}
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModNewPit) NewPitNowAim(body []byte) {
	var msg C2S_NewPitNowAim
	json.Unmarshal(body, &msg)

	//看是否要先选择BUFF
	if len(self.Sql_UserNewPit.userPitInfo.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_BUFF_NOT_CHOOSE"))
		return
	}
	//验证是否连续
	nowLine := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgLine := self.GetLine(msg.Id)
	if msgLine-nowLine != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	nowRow := self.GetRow(self.Sql_UserNewPit.userPitInfo.NowPitId)
	msgRow := self.GetRow(msg.Id)
	if nowRow-msgRow > 1 || nowRow-msgRow < -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CANT_FINISH"))
		return
	}

	config := GetCsvMgr().NewPitConfigMap[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_CONFIG_CANT_EXIST"))
		return
	}

	pitInfo := self.Sql_UserNewPit.newPitInfo[msg.Id]
	if pitInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_PIT_ERROR"))
		return
	}

	self.Sql_UserNewPit.userPitInfo.NowAimId = msg.Id

	var msgRel S2C_NewPitNowAim
	msgRel.Cid = "newpitnowaim"
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModNewPit) NewPitChooseBuff(body []byte) {
	var msg C2S_NewPitChooseBuff
	json.Unmarshal(body, &msg)

	if msg.Index <= 0 || msg.Index > len(self.Sql_UserNewPit.userPitInfo.BuffChoose) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_INDEX_ERROR"))
		return
	}
	buffId := self.Sql_UserNewPit.userPitInfo.BuffChoose[msg.Index-1]
	self.Sql_UserNewPit.userPitInfo.Buff = append(self.Sql_UserNewPit.userPitInfo.Buff, buffId)
	//更新一下英雄血量
	if self.Sql_UserNewPit.userPitInfo.BuffNew == nil {
		self.Sql_UserNewPit.userPitInfo.BuffNew = make(map[int]int, 0)
	}
	_, hasOk := self.Sql_UserNewPit.userPitInfo.BuffNew[buffId]
	if hasOk {
		self.Sql_UserNewPit.userPitInfo.BuffNew[buffId] += 1
	} else {
		self.Sql_UserNewPit.userPitInfo.BuffNew[buffId] = 1
	}
	self.Sql_UserNewPit.userPitInfo.BuffChoose = make([]int, 0)

	var msgRel S2C_NewPitChooseBuff
	msgRel.Cid = "newpitchoosebuff"
	msgRel.Buff = self.Sql_UserNewPit.userPitInfo.Buff
	msgRel.PitInfo = self.MoveTo()
	msgRel.NowPitId = self.Sql_UserNewPit.userPitInfo.NowPitId
	msgRel.NowAimId = self.Sql_UserNewPit.userPitInfo.NowAimId
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModNewPit) NewPitUseItem(body []byte) {

	//看看物品够不够
	if err := self.player.HasObjectOkEasy(ITEM_NEW_PIT_REBORN, 1); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	cost := self.player.RemoveObjectEasy(ITEM_NEW_PIT_REBORN, 1, "地牢复活", 0, 0, 0)

	//更新一下英雄血量
	if self.Sql_UserNewPit.userPitInfo.HeroState == nil {
		self.Sql_UserNewPit.userPitInfo.HeroState = make(map[int]*NewPitHeroState, 0)
	}
	for _, v := range self.Sql_UserNewPit.userPitInfo.HeroState {
		v.Hp = 10000
		v.Energy = 10000
	}
	//将共鸣水晶中的人添加进入
	keyIds := self.player.GetModule("crystal").(*ModResonanceCrystal).GetHeroInThis()
	for _, v := range keyIds {
		self.Sql_UserNewPit.userPitInfo.HeroState[v] = self.NewPitHeroState(v, 10000, 10000)
	}
	//将马车的人加入
	if self.Sql_UserNewPit.userPitInfo != nil {
		for _, v := range self.Sql_UserNewPit.userPitInfo.GetHeroList {
			self.Sql_UserNewPit.userPitInfo.HeroState[v.HeroKeyId] = self.NewPitHeroState(v.HeroKeyId, 10000, 10000)
		}
	}

	var msgRel S2C_NewPitUseItem
	msgRel.Cid = "newpituseitem"
	msgRel.Cost = cost
	msgRel.HeroState = self.Sql_UserNewPit.userPitInfo.HeroState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	//看看有没下层
	config := GetCsvMgr().NewPitConfigMap[self.Sql_UserNewPit.userPitInfo.NowPitId]
	if config != nil {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_NEWPIT_ITEM, config.Difficulty, 0, 0, "地牢使用复活道具", 0, 0, self.player)
	}

	return
}

func (self *ModNewPit) MoveTo() []*NewPitInfo {
	self.Sql_UserNewPit.userPitInfo.NowPitId = self.Sql_UserNewPit.userPitInfo.NowAimId
	self.Sql_UserNewPit.userPitInfo.NowAimId = 0
	rel := make([]*NewPitInfo, 0)
	line := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	for _, v := range self.Sql_UserNewPit.newPitInfo {
		if v.Id == self.Sql_UserNewPit.userPitInfo.NowPitId {
			v.State = NEWPIT_STATE_FINISHED
			rel = append(rel, v)
			continue
		}
		if self.GetLine(v.Id) == line {
			v.State = NEWPIT_STATE_CANT_FINISH
			rel = append(rel, v)
			continue
		}
	}
	return rel
}

func (self *ModNewPit) MoveTo2(id int) []*NewPitInfo {
	self.Sql_UserNewPit.userPitInfo.NowPitId = id
	self.Sql_UserNewPit.userPitInfo.NowAimId = 0
	rel := make([]*NewPitInfo, 0)
	line := self.GetLine(self.Sql_UserNewPit.userPitInfo.NowPitId)
	for _, v := range self.Sql_UserNewPit.newPitInfo {
		if v.Id == self.Sql_UserNewPit.userPitInfo.NowPitId {
			v.State = NEWPIT_STATE_FINISHED
			rel = append(rel, v)
			continue
		}
		if self.GetLine(v.Id) == line {
			v.State = NEWPIT_STATE_CANT_FINISH
			rel = append(rel, v)
			continue
		}
	}
	return rel
}

func (self *ModNewPit) MoveToReward(id int) []*NewPitInfo {
	self.Sql_UserNewPit.userPitInfo.NowAimId = 0
	rel := make([]*NewPitInfo, 0)
	line := self.GetLine(id)
	for _, v := range self.Sql_UserNewPit.newPitInfo {
		if v.Id == id {
			v.State = NEWPIT_STATE_FINISHED
			rel = append(rel, v)
			continue
		}
		if self.GetLine(v.Id) == line {
			v.State = NEWPIT_STATE_CANT_FINISH
			rel = append(rel, v)
			continue
		}
	}
	return rel
}

//和雇佣兵生成相比，差异在 装备和神器的生成走新的一套算法
func (self *ModNewPit) NewHeroForFightInfo(heroId int, level int, star int, exclusiveLv int, exclusiveNum int) *NewHero {

	//修正配置中可能出现的星级低于最低星级的情况
	realStar := star
	itemConfig := GetCsvMgr().ItemMap[11000000+(heroId*100)+1]
	if itemConfig != nil && realStar < itemConfig.Special {
		realStar = itemConfig.Special
	}

	heroInfo := new(NewHero)
	heroInfo.HeroId = heroId
	//继承品质
	heroInfo.StarItem = &StarItem{}
	//这里检测下实际星级
	heroInfo.StarItem.UpStar = realStar
	if level == 0 {
		level = 1
	}
	heroInfo.LvUp(level)
	config := GetCsvMgr().GetNewPitRobotByHeroId(heroId, realStar)
	if config != nil {
		//生成神器
		pItem := &ArtifactEquip{}
		id := 0
		if config.ArtifactType == LOGIC_FALSE {
			id = config.ArtifactId
		} else {
			rand := HF_GetRandom(len(GetCsvMgr().ArtifactEquipConfigMap))
			index := 0
			for _, v := range GetCsvMgr().ArtifactEquipConfigMap {
				index++
				if index > rand {
					id = v.ArtifactId
				}
			}
		}
		pItem.Id = id
		pItem.Lv = config.ArtifactLv
		configArt, ok := GetCsvMgr().ArtifactEquipConfigMap[pItem.Id]
		if ok {
			for i := 0; i < len(configArt.BaseTypes); i++ {
				if configArt.BaseTypes[i] > 0 {
					attr := new(AttrInfo)
					attr.AttrId = i + 1
					pItem.AttrInfo = append(pItem.AttrInfo, attr)
				}
			}

			pItem.CalAttr()
			heroInfo.ArtifactEquipIds = append(heroInfo.ArtifactEquipIds, pItem)
		}

		//生成装备,读取战力
		/*
			randIndex := HF_GetRandom(2)
			if config.ArmsId[randIndex] == 0 {
				randIndex = 0
			}
			equipArm := self.MakeOneEquip(config.ArmsId[randIndex], config.ArmsRange[randIndex])
			heroInfo.EquipIds = append(heroInfo.EquipIds, equipArm)

			randIndex = HF_GetRandom(2)
			if config.HeadId[randIndex] == 0 {
				randIndex = 0
			}
			equipHead := self.MakeOneEquip(config.HeadId[randIndex], config.HeadRange[randIndex])
			heroInfo.EquipIds = append(heroInfo.EquipIds, equipHead)

			randIndex = HF_GetRandom(2)
			if config.BodyId[randIndex] == 0 {
				randIndex = 0
			}
			equipBody := self.MakeOneEquip(config.BodyId[randIndex], config.BodyRange[randIndex])
			heroInfo.EquipIds = append(heroInfo.EquipIds, equipBody)

			randIndex = HF_GetRandom(2)
			if config.ShoesId[randIndex] == 0 {
				randIndex = 0
			}
			equipShoe := self.MakeOneEquip(config.ShoesId[randIndex], config.ShoesRange[randIndex])
			heroInfo.EquipIds = append(heroInfo.EquipIds, equipShoe)

		*/
		//生成专属 exclusiveLv
		configExclusive := GetCsvMgr().GetNewPitRobotExclusive(exclusiveLv, exclusiveNum) //
		heroInfo.ExclusiveEquip = &ExclusiveEquip{}
		for _, v := range GetCsvMgr().ExclusiveEquipConfigMap {
			if v.HeroId == heroInfo.HeroId {
				heroInfo.ExclusiveEquip.Id = v.Id
				for i := 0; i < len(v.BaseType); i++ {
					if v.BaseType[i] > 0 {
						attr := new(AttrInfo)
						attr.AttrId = i + 1
						heroInfo.ExclusiveEquip.AttrInfo = append(heroInfo.ExclusiveEquip.AttrInfo, attr)
					}
				}
			}
		}
		if configExclusive != nil {
			heroInfo.ExclusiveEquip.UnLock = LOGIC_TRUE
			nowLv := len(GetCsvMgr().ExclusiveStrengthen[heroInfo.ExclusiveEquip.Id])
			if nowLv > configExclusive.CartExclusiveLv {
				nowLv = configExclusive.CartExclusiveLv
			}
			heroInfo.ExclusiveEquip.Lv = nowLv
		} else {
			heroInfo.ExclusiveEquip.UnLock = LOGIC_FALSE
			heroInfo.ExclusiveEquip.Lv = 0
		}
		heroInfo.ExclusiveEquip.CalAttr()
	}
	//生成装备,读取战力
	fight := self.player.GetModule("crystal").(*ModResonanceCrystal).GetPriestsEquipFight()
	configHero := GetCsvMgr().GetHeroMapConfig(heroInfo.HeroId, heroInfo.StarItem.UpStar)
	if configHero != nil {
		configEquip := GetCsvMgr().GetHireEquipConfig(configHero.AttackType, fight/100, ROBOT_EQUIP_SUBTYPE_PIT)
		if configEquip != nil {
			for i := 0; i < len(configEquip.Equip); i++ {
				equip := &Equip{}
				equip.Id = configEquip.Equip[i]
				equip.Lv = configEquip.Strengthen[i]
				//configEq, ok := GetCsvMgr().EquipConfigMap[equip.Id]
				//if ok {
				//生成属性
				/*
					for i := 0; i < len(configEq.BaseTypes); i++ {
						if configEq.BaseTypes[i] > 0 {
							attr := new(AttrInfo)
							attr.AttrId = i + 1
							equip.AttrInfo = append(equip.AttrInfo, attr)
						}
					}
					//计算属性
					for _, v := range equip.AttrInfo {
						if v.AttrId <= 0 || v.AttrId > len(configEq.BaseTypes) {
							continue
						}
						index := v.AttrId - 1
						rate := PER_BIT

						if equip.Lv > 0 {
							configLvUp := GetCsvMgr().GetEquipStrengthenConfig(configEq.EquipAttackType, configEq.EquipPosition, configEq.Quality, equip.Lv)
							if configLvUp != nil {
								rate += configLvUp.Vaual[index]
							}
						}

						if config != nil && config.Attribute == configEq.EquipAttackType {
							rate += configEq.CampExtAdd
						}
						v.AttrType = configEq.BaseTypes[index]
						v.AttrValue = configEq.BaseValues[index] * int64(rate) / PER_BIT
					}
				*/
				//}
				heroInfo.EquipIds = append(heroInfo.EquipIds, equip)
			}
		}

		talentConfig := GetCsvMgr().GetStageTalentMap(configHero.TalentGroup)
		if nil != talentConfig {
			if heroInfo.StageTalent == nil {
				heroInfo.StageTalent = &StageTalent{configHero.TalentGroup, []*StageTalentIndex{}}
			}

			if configHero.TalentGroup != heroInfo.StageTalent.Group {
				heroInfo.StageTalent.Group = configHero.TalentGroup
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
					skill := heroInfo.StageTalent.GetTalentSkill(config.Index)
					if skill != nil {
						continue
					}
					rand := len(config.Skill)
					pos := HF_GetRandom(rand) + 1
					// 添加
					heroInfo.StageTalent.AddTalentSkill(config, pos)
				} else {
					// 技能未解锁
					skill := heroInfo.StageTalent.GetTalentSkill(config.Index)
					if skill == nil {
						continue
					}

					// 删除
					heroInfo.StageTalent.RemoveTalentSkill(config.Index)
				}
			}
		}
	}
	heroInfo.CalAttr()
	return heroInfo
}

func (self *ModNewPit) MakeOneEquip(equipId int, equipLv int) *Equip {
	if equipId == 0 {
		return nil
	}
	pItem := &Equip{}
	pItem.Id = equipId
	pItem.Lv = equipLv
	//config, ok := GetCsvMgr().EquipConfigMap[pItem.Id]
	//if ok {
	//生成属性
	/*
		for i := 0; i < len(config.BaseTypes); i++ {
			if config.BaseTypes[i] > 0 {
				attr := new(AttrInfo)
				attr.AttrId = i + 1
				pItem.AttrInfo = append(pItem.AttrInfo, attr)
			}
		}
		//计算属性
		for _, v := range pItem.AttrInfo {
			if v.AttrId <= 0 || v.AttrId > len(config.BaseTypes) {
				continue
			}
			index := v.AttrId - 1
			rate := PER_BIT

			if pItem.Lv > 0 {
				configLvUp := GetCsvMgr().GetEquipStrengthenConfig(config.EquipAttackType, config.EquipPosition, config.Quality, pItem.Lv)
				if configLvUp != nil {
					rate += configLvUp.Vaual[index]
				}
			}

			//if configHero.AttackType == config.EquipAttackType {
			//	rate += config.CampExtAdd
			//}
			v.AttrType = config.BaseTypes[index]
			v.AttrValue = config.BaseValues[index] * int64(rate) / PER_BIT
		}
	*/
	//}
	return pItem
}

func (self *ModNewPit) GetOne(fightInfo *JS_FightInfo) *JS_FightInfo {
	data := new(JS_FightInfo)
	for i := 0; i < len(fightInfo.FightTeamPos.FightPos); i++ {
		if fightInfo.FightTeamPos.FightPos[i] > 0 {
			//根据fightinfo生成的规则，注意索引。
			data.FightTeamPos.FightPos[i] = fightInfo.FightTeamPos.FightPos[i]
			data.Defhero = append(data.Defhero, fightInfo.Defhero[0])
			data.Heroinfo = append(data.Heroinfo, fightInfo.Heroinfo[0])
			data.HeroParam = append(data.HeroParam, fightInfo.HeroParam[0])
			data.Deffight = fightInfo.Deffight
			break
		}
	}
	return data
}

func (self *ModNewPit) MakeEvilCart(fightInfo *JS_FightInfo, level int, star int) []*NewHero {

	data := make([]*NewHero, 0)
	newHero := self.NewHeroForFightInfo(fightInfo.Heroinfo[0].Heroid, level, star, 0, 0)
	if newHero != nil {
		data = append(data, newHero)
	}
	return data
}

func (self *ModNewPit) MakeSoulCart(config *NewPitConfig, exclusiveLv int, exclusiveNum int) []*NewHero {
	level, _, star := self.player.GetModule("crystal").(*ModResonanceCrystal).GetPriestsPitLv()

	//生成品质
	starList := make([]int, 0)
	configStar := GetCsvMgr().GetNewPitRobotQuality(star)
	for i := 0; i < 4; i++ {
		if configStar == nil || i > len(configStar.Star) {
			starList = append(starList, 2)
		} else {
			starList = append(starList, configStar.Star[i])
		}
	}

	data := make([]*NewHero, 0)
	//生成备选列表
	heroList := make(map[int]map[int]int, 0)
	for _, v := range GetCsvMgr().NewPitRobotConfig {
		_, ok := heroList[v.Attribute]
		if !ok {
			heroList[v.Attribute] = make(map[int]int)
		}
		if configStar != nil && v.InitialQuality >= configStar.CommonCart {
			heroList[v.Attribute][v.HeroId] = LOGIC_FALSE
		}
	}

	for i := 0; i < len(config.Select); i++ {
		group := strings.Split(config.Select[i], "|")
		rand := HF_GetRandom(len(group))
		heroId := 0
		for hero, state := range heroList[HF_Atoi(group[rand])] {
			if state == LOGIC_FALSE {
				heroId = hero
				heroList[HF_Atoi(group[rand])][hero] = LOGIC_TRUE
				break
			}
		}
		if heroId == 0 {
			LogError("MakeSoulCart error")
			return data
		}

		newHero := self.NewHeroForFightInfo(heroId, level, starList[i], exclusiveLv, exclusiveNum)
		if newHero != nil {
			data = append(data, newHero)
		}
	}
	return data
}

func (self *ModNewPit) MakeBattleFightInfo(fight int64) *JS_FightInfo {
	data := new(JS_FightInfo)

	baseLv := GetOfflineInfoMgr().GetNewHeroLv(self.player.Sql_UserBase.Uid)
	if baseLv <= 0 {
		baseLv = 1
	}
	formatConfig := GetCsvMgr().GetNewPitRobotMonsterQuality(baseLv)
	if formatConfig == nil {
		return data
	}

	data.Rankid = 0
	data.Uid = 0
	data.Uname = GetCsvMgr().GetName()

	//data.Iconid = cfg.Head
	data.Portrait = 1000
	data.Level = baseLv
	data.Defhero = make([]int, 0)
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)

	//排重结构
	heroExist := make(map[int]int)

	for i := 0; i < len(formatConfig.HeroGroup); i++ {
		err := data.FightTeamPos.addFightPos(i + 1)
		if err != nil {
			continue
		}
		data.Defhero = append(data.Defhero, i)

		//通过组取得英雄ID
		size := len(GetCsvMgr().NewPitRobotGroupMap[formatConfig.HeroGroup[i]])
		rand := HF_GetRandom(size)
		star := HF_GetRandom(formatConfig.HeroUpper[i]-formatConfig.HeroLimit[i]+1) + formatConfig.HeroLimit[i]
		heroIdNow := GetCsvMgr().NewPitRobotGroupMap[formatConfig.HeroGroup[i]][rand]
		//排重
		_, okExist := heroExist[heroIdNow]
		if okExist {
			for _, v := range GetCsvMgr().NewPitRobotGroupMap[formatConfig.HeroGroup[i]] {
				_, okNow := heroExist[v]
				if !okNow {
					heroIdNow = v
					break
				}
			}
		}
		heroExist[heroIdNow] = LOGIC_TRUE
		var hero JS_HeroInfo
		hero.Heroid = heroIdNow
		hero.Color = star
		hero.HeroKeyId = i + 1
		hero.Stars = star
		hero.HeroQuality = star

		//计算等级 如果改配置表属性  要取type99 这里必须对应改

		configHero := GetCsvMgr().HeroConfigMap[hero.Heroid][hero.Stars]
		if configHero == nil {
			hero.Levels = 1
		} else {
			levelParam := 100000000 * (100*fight - configHero.BaseValues[10]) / (configHero.GrowthValues[3] * configHero.QuaValue[3])
			hero.Levels = GetCsvMgr().CalNewPitMonsterLv(levelParam)
		}

		hero.Skin = 0
		hero.ArmsSkill = make([]JS_ArmsSkill, 0)
		hero.TalentSkill = []Js_TalentSkill{}
		hero.MainTalent = 0

		var param JS_HeroParam
		param.Heroid = hero.Heroid

		pAttrWrapper := &AttrWrapper{
			Base:     make([]float64, AttrEnd),
			Ext:      make(map[int]float64),
			Per:      make(map[int]float64),
			Energy:   0,
			FightNum: 0,
		}
		//先计算英雄属性  走英雄接口方便计算
		heroTemp := new(Hero)
		heroTemp.HeroId = hero.Heroid
		heroTemp.checkStarItem(star)
		heroTemp.LvUp(hero.Levels)
		skillBreakConfig := GetCsvMgr().HeroBreakConfigMap[heroTemp.HeroId][heroTemp.StarItem.HeroBreakId]
		if skillBreakConfig == nil {
			return data
		}

		for i := 0; i < len(skillBreakConfig.Skill); i++ {
			if skillBreakConfig.Skill[i] > 0 {
				hero.ArmsSkill = append(hero.ArmsSkill, JS_ArmsSkill{Id: skillBreakConfig.Skill[i] / 100, Level: skillBreakConfig.Skill[i] % 100})
			}
		}

		//计算属性
		heroTemp.StarItem.AttrMap = make(map[int]*Attribute)
		//突破属性
		AddAttrHelperForTimes(heroTemp.StarItem.AttrMap, skillBreakConfig.BaseTypes, skillBreakConfig.BaseValues, 1)
		//星级属性
		configStar := GetCsvMgr().GetHeroMapConfig(heroTemp.HeroId, heroTemp.StarItem.UpStar)
		if configStar == nil {
			return data
		}

		AddAttrHelperForTimes(heroTemp.StarItem.AttrMap, configStar.BaseTypes, configStar.BaseValues, 1)
		//成长率
		configGrowth := GetCsvMgr().HeroGrowthConfigMap[heroTemp.HeroLv]
		if configGrowth != nil {
			AddAttrHelperForGrowth(heroTemp.StarItem.AttrMap, configStar.GrowthTypes, configStar.GrowthValues, configGrowth.GrowthType, configGrowth.GrowthValue, configStar.QuaType, configStar.QuaValue)
		}
		ProcAtt(heroTemp.StarItem.AttrMap, pAttrWrapper)
		//如果星级属性不符合要求，需要增加装备属性
		/*
			if pAttrWrapper.FightNum<fight{

			}
			attMap := make(map[int]*Attribute)
			ProcAtt(attMap, pAttrWrapper)
		*/
		ProcLast(1, pAttrWrapper)
		pAttrWrapper.ExtRet = ProcExtAttr(pAttrWrapper)
		param.Param, param.ExtAttr, param.Energy = pAttrWrapper.Base, pAttrWrapper.ExtRet, pAttrWrapper.Energy
		param.Hp = param.Param[AttrHp]
		hero.Fight = pAttrWrapper.FightNum
		data.Deffight += pAttrWrapper.FightNum

		data.Heroinfo = append(data.Heroinfo, hero)
		data.HeroParam = append(data.HeroParam, param)
		data.FightTeamPos.FightPos[i] = hero.HeroKeyId
	}
	return data
}

func (self *ModNewPit) MakeBattleFightInfoFirst(config *NewPitConfig) *JS_FightInfo {
	data := new(JS_FightInfo)

	firstRobotQuality := strings.Split(config.FirstRobotQuality, "|")
	FirstRobotLv := strings.Split(config.FirstRobotLv, "|")
	FirstRobotGroup := strings.Split(config.FirstRobotGroup, "|")
	if len(firstRobotQuality) != len(FirstRobotLv) || len(firstRobotQuality) != len(FirstRobotGroup) {
		return data
	}

	data.Rankid = 0
	data.Uid = 0
	data.Uname = GetCsvMgr().GetName()

	//data.Iconid = cfg.Head
	data.Portrait = 1000
	data.Level = 1
	data.Defhero = make([]int, 0)
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)

	for i := 0; i < len(firstRobotQuality); i++ {
		err := data.FightTeamPos.addFightPos(i + 1)
		if err != nil {
			continue
		}
		data.Defhero = append(data.Defhero, i)

		//通过组取得英雄ID
		group := HF_Atoi(FirstRobotGroup[i])
		size := len(GetCsvMgr().NewPitRobotGroupFirstMap[group])
		rand := HF_GetRandom(size)
		star := HF_Atoi(firstRobotQuality[i])

		var hero JS_HeroInfo
		hero.Heroid = GetCsvMgr().NewPitRobotGroupFirstMap[group][rand]
		hero.Color = star
		hero.HeroKeyId = i + 1
		hero.Stars = star
		hero.HeroQuality = star

		hero.Levels = HF_Atoi(FirstRobotLv[i])
		hero.ArmsSkill = make([]JS_ArmsSkill, 0)
		hero.TalentSkill = []Js_TalentSkill{}
		hero.MainTalent = 0

		var param JS_HeroParam
		param.Heroid = hero.Heroid

		pAttrWrapper := &AttrWrapper{
			Base:     make([]float64, AttrEnd),
			Ext:      make(map[int]float64),
			Per:      make(map[int]float64),
			Energy:   0,
			FightNum: 0,
		}
		//先计算英雄属性  走英雄接口方便计算
		heroTemp := new(Hero)
		heroTemp.HeroId = hero.Heroid
		heroTemp.checkStarItem(star)
		heroTemp.LvUp(hero.Levels)
		skillBreakConfig := GetCsvMgr().HeroBreakConfigMap[heroTemp.HeroId][heroTemp.StarItem.HeroBreakId]
		if skillBreakConfig == nil {
			return data
		}

		for i := 0; i < len(skillBreakConfig.Skill); i++ {
			if skillBreakConfig.Skill[i] > 0 {
				hero.ArmsSkill = append(hero.ArmsSkill, JS_ArmsSkill{Id: skillBreakConfig.Skill[i] / 100, Level: skillBreakConfig.Skill[i] % 100})
			}
		}

		//计算属性
		heroTemp.StarItem.AttrMap = make(map[int]*Attribute)
		//突破属性
		AddAttrHelperForTimes(heroTemp.StarItem.AttrMap, skillBreakConfig.BaseTypes, skillBreakConfig.BaseValues, 1)
		//星级属性
		configStar := GetCsvMgr().GetHeroMapConfig(heroTemp.HeroId, heroTemp.StarItem.UpStar)
		if configStar == nil {
			return data
		}

		AddAttrHelperForTimes(heroTemp.StarItem.AttrMap, configStar.BaseTypes, configStar.BaseValues, 1)
		//成长率
		configGrowth := GetCsvMgr().HeroGrowthConfigMap[heroTemp.HeroLv]
		if configGrowth != nil {
			AddAttrHelperForGrowth(heroTemp.StarItem.AttrMap, configStar.GrowthTypes, configStar.GrowthValues, configGrowth.GrowthType, configGrowth.GrowthValue, configStar.QuaType, configStar.QuaValue)
		}

		ProcAtt(heroTemp.StarItem.AttrMap, pAttrWrapper)
		//如果星级属性不符合要求，需要增加装备属性
		/*
			if pAttrWrapper.FightNum<fight{

			}
			attMap := make(map[int]*Attribute)
			ProcAtt(attMap, pAttrWrapper)
		*/
		ProcLast(1, pAttrWrapper)
		pAttrWrapper.ExtRet = ProcExtAttr(pAttrWrapper)
		param.Param, param.ExtAttr, param.Energy = pAttrWrapper.Base, pAttrWrapper.ExtRet, pAttrWrapper.Energy
		param.Hp = param.Param[AttrHp]
		hero.Fight = pAttrWrapper.FightNum
		data.Deffight += pAttrWrapper.FightNum

		data.Heroinfo = append(data.Heroinfo, hero)
		data.HeroParam = append(data.HeroParam, param)
		data.FightTeamPos.FightPos[i] = hero.HeroKeyId
	}
	return data
}

func (self *ModNewPit) GetLine(id int) int {
	return (id / 100) % 100
}

func (self *ModNewPit) GetRow(id int) int {
	return (id / 10) % 10
}

func (self *ModNewPit) GetNumberOfPlies(id int) int {
	return id % 10
}

func (self *ModNewPit) GmSuperPass() {
	self.MakeMap(3, 1, 2)

	for _, v := range self.Sql_UserNewPit.newPitInfo {
		if v.Element != NEWPIT_TASK_REWARD {
			v.State = NEWPIT_STATE_FINISHED
		}
		if v.Element == NEWPIT_TASK_BOSS {
			self.Sql_UserNewPit.userPitInfo.NowPitId = v.Id
		}
	}
	self.SendInfo()
}

func (self *ModNewPit) GetShopRefreshType(shopType int) int {
	index := shopType - 1
	if index >= 0 && index < len(NewPitShopType) {
		return NewPitShopType[index]
	}

	return NewPitShopType[0]
}

func (self *ModNewPit) GetShop(shopType int) *NewPitShop {
	for _, v := range self.Sql_UserNewPit.shop {
		if v.Shoptype == shopType {
			return v
		}
	}
	return nil
}
