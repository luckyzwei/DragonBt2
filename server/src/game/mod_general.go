package game

import (
	"encoding/json"
	"errors"
	"fmt"
)

type GeneralAct struct {
	KeyId      int `json:"keyid"` //! 活动Id
	AwardState int `json:"state"` //! 奖励状态 0未获奖, 1领取奖励/活动结束, 2未上榜
	Rank       int `json:"rank"`  //! 玩家当前排名, 排名为0表示不能进榜
	Times      int `json:"times"` //! 抽奖次数
}

// 限时神将模块
type San_General struct {
	Uid       int64
	FreeTimes int                //! 免费次数
	Point     int                //! 积分
	LootTimes int                //! 抽取次数
	Info      string             //! 积分领取信息
	ActRecord string             //! 活动记录
	info      []*JS_GeneralAward //! 积分奖励状态
	actRecord *GeneralAct        //! 活动状态, keyid:活动信息

	DataUpdate
}

type JS_GeneralAward struct {
	Index int `json:"index"` //! 领取奖励索引
	State int `json:"state"` //! 领取状态 0 未领取 1 领取
}

type ModGeneral struct {
	player      *Player
	Sql_General San_General //! 数据库结构
}

func (self *ModGeneral) Decode() {
	json.Unmarshal([]byte(self.Sql_General.Info), &self.Sql_General.info)
	json.Unmarshal([]byte(self.Sql_General.ActRecord), &self.Sql_General.actRecord)
}

func (self *ModGeneral) Encode() {
	self.Sql_General.Info = HF_JtoA(self.Sql_General.info)
	self.Sql_General.ActRecord = HF_JtoA(self.Sql_General.actRecord)
}

func (self *ModGeneral) OnGetData(player *Player) {
	self.player = player
	if self.Sql_General.actRecord == nil {
		self.Sql_General.actRecord = &GeneralAct{}
	}

	if self.Sql_General.info == nil {
		self.Sql_General.info = make([]*JS_GeneralAward, 0)
	}
}

func (self *ModGeneral) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_general` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_General, "san_general", self.player.ID)
	if self.Sql_General.Uid <= 0 {
		self.Sql_General.Uid = self.player.ID
		self.initInfo()
		self.Encode()
		InsertTable("san_general", &self.Sql_General, 0, true)
	} else {
		if self.Sql_General.info == nil {
			self.initAward()
		}

		if self.Sql_General.actRecord == nil {
			self.Sql_General.actRecord = &GeneralAct{}
		}

		self.Decode()
	}
	self.Sql_General.Init("san_general", &self.Sql_General, true)
}

func (self *ModGeneral) OnSave(sql bool) {
	self.Encode()
	self.Sql_General.Update(sql)
}

func (self *ModGeneral) OnRefresh() {
	self.Sql_General.FreeTimes = 1
	self.Sql_General.actRecord.Times = 0

	self.getGeneral()
}

// 消息处理
func (self *ModGeneral) OnMsg(ctrl string, body []byte) bool {
	//log.Println("head:", head, ", ctrl:", ctrl)
	switch ctrl {
	case "generalLoot":
		self.loot(body)
		return true
	case "generalscore":
		self.scoreAward(body)
		return true
	case "getgeneralrank":
		self.getRankNew()
		return true
	case "rankaward":
		self.rankAward()
		return true
	}

	return false
}

func NewGeneralAward(index int) *JS_GeneralAward {
	return &JS_GeneralAward{
		Index: index,
		State: 0,
	}
}

func (self *ModGeneral) initAward() {
	self.Sql_General.info = make([]*JS_GeneralAward, 0)
	for i := 0; i < AWARD_NUM; i++ {
		self.Sql_General.info = append(self.Sql_General.info, NewGeneralAward(i))
	}
}

func (self *ModGeneral) checkActRecord(keyId int) {
	if self.Sql_General.actRecord == nil {
		self.Sql_General.actRecord = &GeneralAct{
			KeyId:      keyId,
			AwardState: 0,
			Rank:       0,
		}
	}

	if self.Sql_General.actRecord.KeyId == 0 {
		self.Sql_General.actRecord.KeyId = keyId
	}
}

// 同步消息
func (self *ModGeneral) getGeneral() {
	if self.Sql_General.info == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_TIME-LIMITED_GOD_WILL_NOT_EXIST"))
		return
	}

	lootConfig, _ := self.GetLootConfig()
	if lootConfig == nil {
		//self.player.SendErrInfo("err", "限时神将活动配置不存在")
		return
	}

	rankConfig := self.GetRankAward(lootConfig.RankAwardGroup)
	if rankConfig == nil || len(rankConfig) < 0 {
		LogError("GetGeneralAward rankConfig == nil || len(rankConfig) < 0, id :", lootConfig.Id)
		return
	}

	showTime, endTime, err := GetGeneralMgr().GetActTime()
	if err != nil {
		//LogError(err.Error())
		msg := &S2C_GetGeneralInfo{
			Cid:      "getgeneral",
			ShowTime: showTime,
			EndTime:  endTime,
			HeroIds:  []int{},
		}
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
		return
	}

	for i := len(self.Sql_General.info); i < AWARD_NUM; i++ {
		self.Sql_General.info = append(self.Sql_General.info, NewGeneralAward(i))
	}

	keyId := GetGeneralMgr().getKeyId()
	self.checkActRecord(keyId)
	scoreAward, scorePoint := self.getScoreAward(lootConfig.RewardPointsGroup)
	var heroIds []int
	heroIds = append(heroIds, lootConfig.Mainheroids...)
	heroIds = append(heroIds, lootConfig.HeroIds...)

	certaintimes := 10
	//lootType := config.CallSingleType
	config, err := self.getCheckCsv(lootConfig.CallTenType)
	if err == nil {
		certaintimes = config.Certaintimesdroptype
	}
	msg := &S2C_GetGeneralInfo{
		Cid:              "getgeneral",
		Score:            self.Sql_General.Point,
		LootTimes:        certaintimes - self.Sql_General.LootTimes%certaintimes,
		GeneralAward:     self.Sql_General.info,
		FreeTimes:        self.Sql_General.FreeTimes,
		HeroIds:          heroIds,
		RankConfig:       rankConfig,
		ShowTime:         showTime,
		EndTime:          endTime,
		RankAward:        self.Sql_General.actRecord.AwardState,
		ScoreAward:       scoreAward,
		ScorePoints:      scorePoint,
		CostSingleNum:    lootConfig.CostSingleNum,
		CostTenNum:       lootConfig.CostTenNum,
		ActRecord:        self.Sql_General.actRecord,
		ServerId:         GetServer().Con.ServerId,
		ServerName:       GetServer().Con.ServerName,
		ActType:          lootConfig.ActType,
		CallDesc:         lootConfig.CallDesc,
		NewHero:          lootConfig.NewHero,
		MainHeroLocation: lootConfig.MainHeroLocation,
		HeroLocation:     lootConfig.HeroLocation,
	}

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModGeneral) getScoreAward(groupId int) ([][]PassItem, []int) {
	var resConfig []*TimeGeneralsPointsConfig
	configLst := GetCsvMgr().TimeGeneralsPointsConfig
	for index, elem := range configLst {
		if elem.Group != groupId {
			continue
		}

		if configLst[index] == nil {
			continue
		}

		resConfig = append(resConfig, configLst[index])
	}

	var scoreAward [][]PassItem
	var scorePoint []int
	for index := range resConfig {
		elem := resConfig[index]
		if elem == nil {
			continue
		}
		var items []PassItem
		award := elem.Award
		nums := elem.Nums
		if len(award) != len(nums) {
			continue
		}

		for awardIndex := range award {
			if award[awardIndex] == 0 {
				continue
			}

			if nums[awardIndex] == 0 {
				continue
			}

			items = append(items, PassItem{ItemID: award[awardIndex], Num: nums[awardIndex]})
		}
		scoreAward = append(scoreAward, items)
		scorePoint = append(scorePoint, elem.Points)
	}
	return scoreAward, scorePoint
}

func (self *ModGeneral) loot(body []byte) {
	if !GetGeneralMgr().IsActOk() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_THE_EVENT_IS_OVER_AND"))
		return
	}

	var msg C2S_LootGeneral
	json.Unmarshal(body, &msg)
	lootType := msg.LootType
	if lootType < 1 || lootType > 3 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_WRONG_DROP_TYPE"))
		return
	}

	if lootType == 1 {
		self.freeLoot()
	} else if lootType == 2 {
		if msg.IsUseItem == LOGIC_TRUE {
			self.singleItemLoot()
		} else {
			self.singleLoot()
		}
	} else if lootType == 3 {
		if msg.IsUseItem == LOGIC_TRUE {
			self.tenItemLoot()
		} else {
			self.tenLoot()
		}
	}
}

// 免费抽
func (self *ModGeneral) freeLoot() {
	pGeneral := &self.Sql_General
	if pGeneral.FreeTimes <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_THE_NUMBER_OF_FREE_LOTTERY"))
		return
	}

	// 获取抽奖配置
	//timeConfig := GetCsvMgr().TimeGeneral_CSV
	//if len(timeConfig) != 1 {
	//	self.player.SendErrInfo("err", "掉落配置错误,出现两个活动!")
	//	return
	//}

	lootConfig, _ := self.GetLootConfig()
	if lootConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_FALLING_CONFIGURATION_IS_EMPTY"))
		return
	}

	config, err := self.getCheckCsv(lootConfig.CallSingleType)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	var allItem []PassItem
	var records []*GeneralRecord

	bag := make([]PassItem, 0)
	certaintimes := config.Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(config.Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := config.Dropgroups[3]
		tempCer := self.Sql_General.LootTimes % certaintimes
		tempCer += 1
		if tempCer/certaintimes > 0 {
			itemid = certainitem
		}

		if itemid != 0 {
			var dropitem PassItem
			dropitem.ItemID = itemid
			dropitem.Num = 1
			bag = append(bag, dropitem)
			LogDebug("drop baodi.......")
		}
	}

	res, record := self.LootItemNew(config, lootConfig.CallSingleType, 1, self.Sql_General.LootTimes, bag)
	allItem = append(allItem, res...)
	records = append(records, record...)

	var items []PassItem
	for _, v := range allItem {
		itemId, itemNum := self.player.AddObject(v.ItemID, v.Num, 0, 0, 0, "次元召唤单次召唤")
		items = append(items, PassItem{itemId, itemNum})
	}
	self.addLootTimes(1)
	self.addScore(lootConfig.CallSinglePoint, records)
	if pGeneral.FreeTimes > 0 {
		pGeneral.FreeTimes -= 1
	}
	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, 1, 1, 0)
	////限时神将单抽 行为日志
	//GetServer().SqlLog(self.player.GetUid(), LOG_ACT_DIAL, 1, 0, 0, "限时神将免费单抽", 0, 0, self.player)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GENERAL_FIND_SINGLE, 1, 0, 0, "次元召唤单次召唤", 0, 0, self.player)

	msg := &S2C_GeneralLootInfo{
		Cid:       "generalLoot",
		Score:     pGeneral.Point,
		LootTimes: certaintimes - self.Sql_General.LootTimes%certaintimes,
		LootItem:  items,
		Gem:       self.player.Sql_UserBase.Gem,
		FreeTimes: pGeneral.FreeTimes,
		Times:     self.getTimes(),
	}
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModGeneral) lootItems(groupId, num int) ([]PassItem, []*GeneralRecord) {
	var allItem []PassItem
	records := make([]*GeneralRecord, 0)
	modFind := self.player.GetModule("find").(*ModFind)
	if modFind == nil {
		LogError("module find is nil!")
		return []PassItem{}, records
	}

	itemLoot, record := modFind.FindDrop(groupId, num)
	if len(itemLoot) <= 0 {
		LogError("lootItems len(itemLoot) <= 0")
		return []PassItem{}, records
	}

	records = append(records, record...)
	// 添加到玩家身上
	for index := range itemLoot {
		item := itemLoot[index]
		config := GetCsvMgr().GetItemConfig(item.ItemID)
		if config == nil {
			LogError("config == nil, itemId:", item.ItemID)
			continue
		}
		allItem = append(allItem, PassItem{item.ItemID, item.Num})
	}

	return allItem, records
}

func (self *ModGeneral) GetLootConfig() (*TimeGeneralsConfig, int) {
	return GetCsvMgr().GetLootConfig()
}

func (self *ModGeneral) GetRankAward(groupId int) []*TimeGeneralRank {
	return GetCsvMgr().GetRankAwardConf(groupId)
}

// 单抽, 增加vip判断
func (self *ModGeneral) singleLoot() {
	pGeneral := &self.Sql_General
	// 获取抽奖配置
	lootConfig, _ := self.GetLootConfig()

	if lootConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_FALLING_CONFIGURATION_IS_EMPTY"))
		return
	}

	if lootConfig.CostSingleNum <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_LOTTERY_WASTAGE_IS_INCORRECT"))
		return
	}

	if self.player.Sql_UserBase.Gem < lootConfig.CostSingleNum {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_DIAMOND_SHORTAGE"))
		return
	}

	if !self.checkVipTimes(1) {
		return
	}

	config, err := self.getCheckCsv(lootConfig.CallSingleType)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	var allItem []PassItem
	var records []*GeneralRecord

	bag := make([]PassItem, 0)
	certaintimes := config.Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(config.Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := config.Dropgroups[3]
		tempCer := self.Sql_General.LootTimes % certaintimes
		tempCer += 1
		if tempCer/certaintimes > 0 {
			itemid = certainitem
		}

		if itemid != 0 {
			var dropitem PassItem
			dropitem.ItemID = itemid
			dropitem.Num = 1
			bag = append(bag, dropitem)
			LogDebug("drop baodi.......")
		}
	}

	res, record := self.LootItemNew(config, lootConfig.CallSingleType, 1, self.Sql_General.LootTimes, bag)
	allItem = append(allItem, res...)
	records = append(records, record...)

	//发送英雄，如果开启了转换开关，就自动分解
	getItemTran := make(map[int]*Item) //转换
	param3 := lootConfig.CostTenNum

	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, 1, 1, 0)
	nCount := 0
	nDecompose := 0
	for i := 0; i < len(allItem); i++ {
		//如果用户打开了自动分解，需要判断这个物品是否转化为碎片
		itemId, itemNum := self.player.GetModule("hero").(*ModHero).CheckItem(allItem[i].ItemID, allItem[i].Num)
		if len(itemId) > 0 {
			AddItemMapHelper(getItemTran, itemId, itemNum)
			nDecompose++
		} else {
			allItem[i].ItemID, allItem[i].Num = self.player.AddObject(allItem[i].ItemID, allItem[i].Num, 0, 0, 0, "众神降临单抽")
		}
		config := GetCsvMgr().GetItemConfig(allItem[i].ItemID)
		if nil != config {
			if config.ItemCheck >= 4 {
				nCount++
			}
		}
	}

	self.player.HandleTask(TASK_TYPE_SUMMON_ELITE_HEROS, nCount, 0, 0)
	self.player.HandleTask(TASK_TYPE_DECOMPOSE_HEROS, nDecompose, 0, 0)
	getItemsTran := self.player.AddObjectItemMap(getItemTran, "次元召唤单次召唤", 0, 0, param3)

	////限时神将单抽 行为日志
	//GetServer().SqlLog(self.player.GetUid(), LOG_GENERAL_HERO, 1, 0, 0, "限时神将单抽", 0, 0, self.player)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GENERAL_FIND_SINGLE, 1, 1, 0, "次元召唤单次召唤", 0, 0, self.player)

	self.addScore(lootConfig.CallSinglePoint, records)
	//self.player.Sql_UserBase.Gem -= lootConfig.CostSingleNum
	costItems := self.player.RemoveObjectSimple(lootConfig.CostSingleitem, lootConfig.CostSingleNum, "次元召唤单次召唤", 0, 0, 0)

	//限时神将单抽 钻石消耗日志
	//GetServer().SqlLog(self.player.GetUid(),DEFAULT_GEM,-lootConfig.CostSingleNum,0,0,"限时神将单抽",self.player.Sql_UserBase.Gem,0,self.player)

	self.addLootTimes(1)
	self.addTimes(1)
	msg := &S2C_GeneralLootInfo{
		Cid:          "generalLoot",
		Score:        pGeneral.Point,
		LootTimes:    certaintimes - self.Sql_General.LootTimes%certaintimes,
		LootItem:     allItem,
		Gem:          self.player.Sql_UserBase.Gem,
		FreeTimes:    pGeneral.FreeTimes,
		Times:        self.getTimes(),
		CostItems:    costItems,
		GetItemsTran: getItemsTran,
	}

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

}

func (self *ModGeneral) singleItemLoot() {
	pGeneral := &self.Sql_General
	// 获取抽奖配置
	lootConfig, _ := self.GetLootConfig()

	if lootConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_FALLING_CONFIGURATION_IS_EMPTY"))
		return
	}

	if lootConfig.CostSingleNum <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_LOTTERY_WASTAGE_IS_INCORRECT"))
		return
	}

	costItem := make([]int, 0)
	costNum := make([]int, 0)

	costItem = append(costItem, ITEM_FIND_GENERAL_ITEM)
	costNum = append(costNum, 1)

	if err := self.player.HasObjectOk(costItem, costNum); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	if !self.checkVipTimes(1) {
		return
	}

	config, err := self.getCheckCsv(lootConfig.CallSingleType)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	var allItem []PassItem
	var records []*GeneralRecord

	bag := make([]PassItem, 0)
	certaintimes := config.Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(config.Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := config.Dropgroups[3]
		tempCer := self.Sql_General.LootTimes % certaintimes
		tempCer += 1
		if tempCer/certaintimes > 0 {
			itemid = certainitem
		}

		if itemid != 0 {
			var dropitem PassItem
			dropitem.ItemID = itemid
			dropitem.Num = 1
			bag = append(bag, dropitem)
			LogDebug("drop baodi.......")
		}
	}

	res, record := self.LootItemNew(config, lootConfig.CallSingleType, 1, self.Sql_General.LootTimes, bag)
	allItem = append(allItem, res...)
	records = append(records, record...)

	//发送英雄，如果开启了转换开关，就自动分解
	getItemTran := make(map[int]*Item) //转换
	param3 := lootConfig.CostTenNum

	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, 1, 1, 0)
	nCount := 0
	nDecompose := 0
	for i := 0; i < len(allItem); i++ {
		//如果用户打开了自动分解，需要判断这个物品是否转化为碎片
		itemId, itemNum := self.player.GetModule("hero").(*ModHero).CheckItem(allItem[i].ItemID, allItem[i].Num)
		if len(itemId) > 0 {
			AddItemMapHelper(getItemTran, itemId, itemNum)
			nDecompose++
		} else {
			allItem[i].ItemID, allItem[i].Num = self.player.AddObject(allItem[i].ItemID, allItem[i].Num, 0, 0, 0, "众神降临单抽")
		}
		config := GetCsvMgr().GetItemConfig(allItem[i].ItemID)
		if nil != config {
			if config.ItemCheck >= 4 {
				nCount++
			}
		}
	}

	self.player.HandleTask(TASK_TYPE_SUMMON_ELITE_HEROS, nCount, 0, 0)
	self.player.HandleTask(TASK_TYPE_DECOMPOSE_HEROS, nDecompose, 0, 0)
	getItemsTran := self.player.AddObjectItemMap(getItemTran, "次元召唤单次召唤", 0, 0, param3)

	////限时神将单抽 行为日志
	//GetServer().SqlLog(self.player.GetUid(), LOG_GENERAL_HERO, 1, 0, 0, "限时神将单抽", 0, 0, self.player)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GENERAL_FIND_SINGLE, 1, 1, 0, "次元召唤单次召唤", 0, 0, self.player)

	self.addScore(lootConfig.CallSinglePoint, records)
	//self.player.Sql_UserBase.Gem -= lootConfig.CostSingleNum
	costItems := self.player.RemoveObjectLst(costItem, costNum, "次元召唤单次召唤", 0, 0, 0)

	//限时神将单抽 钻石消耗日志
	//GetServer().SqlLog(self.player.GetUid(),DEFAULT_GEM,-lootConfig.CostSingleNum,0,0,"限时神将单抽",self.player.Sql_UserBase.Gem,0,self.player)

	self.addLootTimes(1)
	self.addTimes(1)
	msg := &S2C_GeneralLootInfo{
		Cid:          "generalLoot",
		Score:        pGeneral.Point,
		LootTimes:    certaintimes - self.Sql_General.LootTimes%certaintimes,
		LootItem:     allItem,
		Gem:          self.player.Sql_UserBase.Gem,
		FreeTimes:    pGeneral.FreeTimes,
		Times:        self.getTimes(),
		CostItems:    costItems,
		GetItemsTran: getItemsTran,
	}

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

}

func (self *ModGeneral) addLootTimes(n int) {
	self.Sql_General.LootTimes += n
}

// 抽奖次数
func (self *ModGeneral) addTimes(n int) {
	self.Sql_General.actRecord.Times += n
}

func (self *ModGeneral) getTimes() int {
	return self.Sql_General.actRecord.Times
}

func (self *ModGeneral) mustLoot(config *PubchesttotalConfig) ([]PassItem, []*GeneralRecord) {
	groupId := config.Dropgroups[3]
	if groupId == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_THE_NUMBER_OF_CONFIGURATION_ERRORS"))
		return []PassItem{}, []*GeneralRecord{}
	}
	res, records := self.lootItems(groupId, 1)
	return res, records
}

func (self *ModGeneral) commomLoot(config *PubchesttotalConfig) ([]PassItem, []*GeneralRecord) {

	records := make([]*GeneralRecord, 0)

	var groupIds []int
	var groupNum []int
	for i := 0; i < 3; i++ {
		id := config.Dropgroups[i]
		if id == 0 {
			continue
		}
		num := config.Dropgroupids[i]
		if num == 0 {
			continue
		}
		groupIds = append(groupIds, id)
		groupNum = append(groupNum, num)
	}

	var allItem []PassItem
	for i := 0; i < len(groupIds); i++ {
		res, record := self.lootItems(groupIds[i], groupNum[i])
		if len(res) > 0 {
			allItem = append(allItem, res...)
			records = append(records, record...)
		}
	}

	return allItem, records
}

// 10抽, 增加vip判断
func (self *ModGeneral) tenLoot() {
	pGeneral := &self.Sql_General
	// 获取抽奖配置
	lootConfig, _ := self.GetLootConfig()
	if lootConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_FALLING_CONFIGURATION_IS_EMPTY"))
		return
	}

	if lootConfig.CostTenNum <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_LOTTERY_ALLOCATION_CONSUMPTION_IS_LESS"))
		return
	}

	if self.player.Sql_UserBase.Gem < lootConfig.CostTenNum {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_DIAMOND_SHORTAGE"))
		return
	}

	if !self.checkVipTimes(10) {
		return
	}

	//lootType := config.CallSingleType
	config, err := self.getCheckCsv(lootConfig.CallTenType)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	var allItem []PassItem
	var records []*GeneralRecord
	bag := make([]PassItem, 0)
	certaintimes := config.Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(config.Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := config.Dropgroups[3]
		tempCer := self.Sql_General.LootTimes % certaintimes
		tempCer += 10
		if tempCer/certaintimes > 0 {
			itemid = certainitem
		}

		if itemid != 0 {
			var dropitem PassItem
			dropitem.ItemID = itemid
			dropitem.Num = 1
			bag = append(bag, dropitem)
			LogDebug("drop baodi.......")
		}
	}

	res, record := self.LootItemNew(config, lootConfig.CallTenType, 10, self.Sql_General.LootTimes, bag)
	allItem = append(allItem, res...)
	records = append(records, record...)

	//发送英雄，如果开启了转换开关，就自动分解
	getItemTran := make(map[int]*Item) //转换
	param3 := lootConfig.CostTenNum

	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, 10, 1, 0)
	nCount := 0
	nDecompose := 0
	for i := 0; i < len(allItem); i++ {
		//如果用户打开了自动分解，需要判断这个物品是否转化为碎片
		itemId, itemNum := self.player.GetModule("hero").(*ModHero).CheckItem(allItem[i].ItemID, allItem[i].Num)
		if len(itemId) > 0 {
			AddItemMapHelper(getItemTran, itemId, itemNum)
			nDecompose++
		} else {
			allItem[i].ItemID, allItem[i].Num = self.player.AddObject(allItem[i].ItemID, allItem[i].Num, 0, 0, 0, "众神降临十连抽")
		}
		config := GetCsvMgr().GetItemConfig(allItem[i].ItemID)
		if nil != config {
			if config.ItemCheck >= 4 {
				nCount++
			}
		}
	}

	self.player.HandleTask(TASK_TYPE_SUMMON_ELITE_HEROS, nCount, 0, 0)
	self.player.HandleTask(TASK_TYPE_DECOMPOSE_HEROS, nDecompose, 0, 0)
	getItemsTran := self.player.AddObjectItemMap(getItemTran, "次元召唤十连召唤", 0, 0, param3)

	self.addScore(lootConfig.CallTenPoint, records)
	//self.player.Sql_UserBase.Gem -= lootConfig.CostTenNum
	costItems := self.player.RemoveObjectSimple(lootConfig.Costtenitem, lootConfig.CostTenNum, "次元召唤十连召唤", 0, 0, 0)

	////限时神将十连抽 行为日志
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GENERAL_FIND_TEN, 1, 0, 0, "次元召唤十连召唤", 0, 0, self.player)

	//限时神将十连抽 钻石消耗日志
	//GetServer().SqlLog(self.player.GetUid(), DEFAULT_GEM, -lootConfig.CostTenNum, 0, 0, "限时神将十连抽", self.player.Sql_UserBase.Gem, 0, self.player)
	self.addLootTimes(10)
	self.addTimes(10)
	msg := &S2C_GeneralLootInfo{
		Cid:          "generalLoot",
		Score:        pGeneral.Point,
		LootTimes:    certaintimes - self.Sql_General.LootTimes%certaintimes,
		LootItem:     allItem,
		Gem:          self.player.Sql_UserBase.Gem,
		FreeTimes:    pGeneral.FreeTimes,
		Times:        self.getTimes(),
		CostItems:    costItems,
		GetItemsTran: getItemsTran,
	}

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModGeneral) tenItemLoot() {
	pGeneral := &self.Sql_General
	// 获取抽奖配置
	lootConfig, _ := self.GetLootConfig()
	if lootConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_FALLING_CONFIGURATION_IS_EMPTY"))
		return
	}

	if lootConfig.CostTenNum <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_LOTTERY_ALLOCATION_CONSUMPTION_IS_LESS"))
		return
	}

	costItem := make([]int, 0)
	costNum := make([]int, 0)

	costItem = append(costItem, ITEM_FIND_GENERAL_ITEM)
	costNum = append(costNum, 10)

	if err := self.player.HasObjectOk(costItem, costNum); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	if !self.checkVipTimes(10) {
		return
	}

	//lootType := config.CallSingleType
	config, err := self.getCheckCsv(lootConfig.CallTenType)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	var allItem []PassItem
	var records []*GeneralRecord
	bag := make([]PassItem, 0)
	certaintimes := config.Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(config.Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := config.Dropgroups[3]
		tempCer := self.Sql_General.LootTimes % certaintimes
		tempCer += 10
		if tempCer/certaintimes > 0 {
			itemid = certainitem
		}

		if itemid != 0 {
			var dropitem PassItem
			dropitem.ItemID = itemid
			dropitem.Num = 1
			bag = append(bag, dropitem)
			LogDebug("drop baodi.......")
		}
	}

	res, record := self.LootItemNew(config, lootConfig.CallTenType, 10, self.Sql_General.LootTimes, bag)
	allItem = append(allItem, res...)
	records = append(records, record...)

	//发送英雄，如果开启了转换开关，就自动分解
	getItemTran := make(map[int]*Item) //转换
	param3 := lootConfig.CostTenNum

	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, 10, 1, 0)
	nCount := 0
	nDecompose := 0
	for i := 0; i < len(allItem); i++ {
		//如果用户打开了自动分解，需要判断这个物品是否转化为碎片
		itemId, itemNum := self.player.GetModule("hero").(*ModHero).CheckItem(allItem[i].ItemID, allItem[i].Num)
		if len(itemId) > 0 {
			AddItemMapHelper(getItemTran, itemId, itemNum)
			nDecompose++
		} else {
			allItem[i].ItemID, allItem[i].Num = self.player.AddObject(allItem[i].ItemID, allItem[i].Num, 0, 0, 0, "众神降临十连抽")
		}
		config := GetCsvMgr().GetItemConfig(allItem[i].ItemID)
		if nil != config {
			if config.ItemCheck >= 4 {
				nCount++
			}
		}
	}

	self.player.HandleTask(TASK_TYPE_SUMMON_ELITE_HEROS, nCount, 0, 0)
	self.player.HandleTask(TASK_TYPE_DECOMPOSE_HEROS, nDecompose, 0, 0)
	getItemsTran := self.player.AddObjectItemMap(getItemTran, "次元召唤十连召唤", 0, 0, param3)

	self.addScore(lootConfig.CallTenPoint, records)
	//self.player.Sql_UserBase.Gem -= lootConfig.CostTenNum
	costItems := self.player.RemoveObjectLst(costItem, costNum, "次元召唤十连召唤", 0, 0, 0)

	////限时神将十连抽 行为日志
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GENERAL_FIND_TEN, 1, 0, 0, "次元召唤十连召唤", 0, 0, self.player)

	//限时神将十连抽 钻石消耗日志
	//GetServer().SqlLog(self.player.GetUid(), DEFAULT_GEM, -lootConfig.CostTenNum, 0, 0, "限时神将十连抽", self.player.Sql_UserBase.Gem, 0, self.player)
	self.addLootTimes(10)
	self.addTimes(10)
	msg := &S2C_GeneralLootInfo{
		Cid:          "generalLoot",
		Score:        pGeneral.Point,
		LootTimes:    certaintimes - self.Sql_General.LootTimes%certaintimes,
		LootItem:     allItem,
		Gem:          self.player.Sql_UserBase.Gem,
		FreeTimes:    pGeneral.FreeTimes,
		Times:        self.getTimes(),
		CostItems:    costItems,
		GetItemsTran: getItemsTran,
	}

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

// 先按照paypubtype,再按照vip,最后按照权重来选取一个掉落包
func (self *ModGeneral) getCheckCsv(payPubType int) (*PubchesttotalConfig, error) {
	lst := GetCsvMgr().PubchesttotalGroup[payPubType]
	if len(lst) == 0 {
		return nil, errors.New(GetCsvMgr().GetText("STR_MOD_GENERAL_TIMELINE_GOD_WILL_DROP_CONFIGURATION"))
	}

	// 筛选vip
	var res []*PubchesttotalConfig
	vipLv := self.player.Sql_UserBase.Vip
	chance := 0
	for index := range lst {
		vipMin := lst[index].Cardsvipmin
		vipMax := lst[index].Cardsvipmax
		if vipLv < vipMin || vipLv > vipMax {
			continue
		}
		res = append(res, lst[index])
		chance += lst[index].Chance
	}

	if chance == 0 {
		return nil, errors.New(GetCsvMgr().GetText("STR_MOD_GENERAL_TIME-LIMITED_GOD_WILL_DROP_THE"))
	}

	randNum := HF_GetRandom(chance) + 1
	check := 0
	var config *PubchesttotalConfig
	// 随机找到其中一项
	for _, v := range res {
		check += v.Chance
		if randNum <= check {
			config = v
			break
		}
	}

	if config == nil {
		return nil, errors.New(GetCsvMgr().GetText("STR_MOD_GENERAL_TIME-LIMITED_GOD_WILL_FALL_CONFIGURATION"))
	}

	return config, nil
}

// 同时需要上传当前积分到中心服务器
func (self *ModGeneral) addScore(score int, records []*GeneralRecord) {
	self.Sql_General.Point += score
	player := &self.player.Sql_UserBase
	top := &Js_GeneralUser{
		Uid:     player.Uid,
		SvrId:   GetServer().Con.ServerId,
		SvrName: GetServer().Con.ServerName,
		UName:   player.UName,
		Level:   player.Level,
		Vip:     player.Vip,
		Icon:    player.IconId,
		Point:   self.Sql_General.Point,
		Rank:    0,
		KeyId:   self.getKeyId(),
		Time:    TimeServer().Unix(),
	}
	for _, v := range records {
		v.Uid = self.player.Sql_UserBase.Uid
		v.SvrId = GetServer().Con.ServerId
		v.UName = self.player.Sql_UserBase.UName
		v.Time = TimeServer().Unix()
	}
	//GetGeneralMgr().UploadScore(top)
	GetGeneralMgr().UploadScoreNew(self.player, top, records)
}

func (self *ModGeneral) getKeyId() int {
	lootConfig, step := self.GetLootConfig()
	if lootConfig == nil {
		return 0
	}
	return lootConfig.KeyId*1000 + step
}

func (self *ModGeneral) initInfo() {
	self.Sql_General.FreeTimes = 1
	self.Sql_General.Point = 0
	self.Sql_General.LootTimes = 0
	self.initAward()
	self.Sql_General.actRecord = &GeneralAct{}
}

func (self *ModGeneral) scoreAward(body []byte) {
	var msg C2S_GeneralAward
	json.Unmarshal(body, &msg)
	awardIndex := msg.AwardIndex
	if awardIndex < 1 || awardIndex > AWARD_NUM {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_DRAW_INDEX_ERROR")+fmt.Sprintf("%d", awardIndex))
		return
	}

	// 根据index和mode找到对应的配置
	lootConfig, _ := self.GetLootConfig()
	if lootConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_DROP_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	configLst := GetCsvMgr().TimeGeneralsPointsConfig
	var resConfig []*TimeGeneralsPointsConfig
	for index, elem := range configLst {
		if elem.Group != lootConfig.Id {
			continue
		}
		resConfig = append(resConfig, configLst[index])
	}

	if len(resConfig) != AWARD_NUM {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_INCORRECT_CONFIGURATION_DATA"))
		return
	}

	readIndex := awardIndex - 1
	targetConfig := resConfig[readIndex]
	needScore := targetConfig.Points
	PData := &self.Sql_General
	if PData.Point < needScore {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_PLAYERS_WITH_INSUFFICIENT_POINTS"))
		return
	}

	awardInfo := self.Sql_General.info
	// 兼容逻辑
	if readIndex >= len(awardInfo) {
		awardInfo = append(awardInfo, NewGeneralAward(awardIndex))
	}

	pInfo := awardInfo[readIndex]
	if pInfo.State != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_THE_TREASURE_BOX_AWARD_HAS"))
		return
	}

	// 拿到奖励配置
	award := targetConfig.Award
	nums := targetConfig.Nums
	if len(award) != len(nums) {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_GENERAL_CONFIGURATION_ERROR:_LEN_AWARD=_LEN"))
		return
	}

	var res []PassItem
	for i := range award {
		itemId, itemNum := self.player.AddObject(award[i], nums[i], needScore, 0, 0, "次元召唤积分宝箱")
		res = append(res, PassItem{itemId, itemNum})
	}

	////限时神将宝箱领取 行为日志
	//GetServer().SqlLog(self.player.GetUid(), LOG_GENERAL_HERO, 1, needScore, 0, "限时神将宝箱领取", 0, 0, self.player)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GENERAL_FIND_BOX, targetConfig.Id, needScore, 0, "次元召唤积分宝箱", 0, 0, self.player)

	pInfo.State = 1
	data := &S2C_GeneralAward{
		Cid:    "generalscore",
		Items:  res,
		Status: self.Sql_General.info,
	}

	self.player.SendMsg(data.Cid, HF_JtoB(data))
}

// 2.排行榜更新数据: now >= start + continue && now <= start + continue + show, 可以领取奖励
func (self *ModGeneral) rankAward() {
	showTime, endTime, err := GetGeneralMgr().GetActTime()
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	now := TimeServer().Unix()
	if now < showTime || now > endTime {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_YOU_CANT_GET_THE_LIST"))
		return
	}

	lootConfig, _ := self.GetLootConfig()
	if lootConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKEGG_CONFIGURATION_ERROR"))
		return
	}

	keyId := lootConfig.KeyId
	self.checkActRecord(keyId)

	if self.Sql_General.actRecord.AwardState == 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_LIST_AWARDS_HAVE_BEEN_RECEIVED"))
		return
	}

	rank, err := GetGeneralMgr().getRankAwardNew(self.player.Sql_UserBase.Uid)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	point := self.Sql_General.Point
	pRankConfig := GetCsvMgr().GetGeneralAward(rank, point)
	if pRankConfig == nil {
		LogError("pRankConfig == nil!")
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKEGG_CONFIGURATION_ERROR"))
		return
	}

	items := pRankConfig.GetAward(rank, point)
	hasExt := pRankConfig.IsHasExt(rank, point)
	var result []PassItem
	for _, v := range items {
		itemId, itemNum := self.player.AddObject(v.ItemID, v.Num, rank, 0, 0, "次元召唤排行奖励")
		result = append(result, PassItem{itemId, itemNum})
	}

	////限时神将排行奖励领取 行为日志
	//GetServer().SqlLog(self.player.GetUid(), LOG_GENERAL_HERO, 1, rank, 0, "限时神将排行奖励领取", 0, 0, self.player)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GENERAL_FIND_RANK, rank, hasExt, 0, "次元召唤排行奖励", 0, 0, self.player)

	self.Sql_General.actRecord.AwardState = 1
	self.Sql_General.actRecord.Rank = rank
	msg := &S2C_GeneralRankAward{
		Cid:   "rankaward",
		Items: result,
		State: 1,
	}
	self.player.SendMsg(msg.Cid, HF_JtoB(msg))

	//! 强制保存
	self.OnSave(true)
}

// 1.先进行领奖检查, 2.再进行活动重置检查
// 登陆时进行检查, 获取信息检查玩家是否发送邮件
// 活动结束后, 由主服务器发送消息给从服务器,进行增加,如果没有收到不继续发送
func (self *ModGeneral) checkMail() {
	if self.Sql_General.Point <= 0 {
		return
	}

	actRecord := self.Sql_General.actRecord
	if actRecord == nil {
		return
	}

	if actRecord.AwardState == 1 || actRecord.AwardState == 2 {
		return
	}

	// 检查活动有没有结束?
	endTime := GetGeneralMgr().GetCheckActTime()
	if TimeServer().Unix() <= endTime {
		return
	}

	rank := GetGeneralMgr().GetUserRankNew(self.Sql_General.Uid)
	if rank == 0 {
		actRecord.AwardState = 2
		return
	}

	rankTitle, mailText, allItems, err := GetGeneralMgr().CheckMailNew(rank, self.Sql_General.Point)
	if err != nil {
		LogError("check mail err, err:", err.Error())
		return
	}

	// 检查玩家有没有在排行榜里面
	pMail := self.player.GetModule("mail").(*ModMail)
	if pMail == nil {
		LogError("checkMail in mod mail, pMail == nil!")
		return
	}

	// 发送奖励
	pMail.AddMail(1, 1, 0, rankTitle, mailText, GetCsvMgr().GetText("STR_SYS"), allItems, false, 0)
	actRecord.AwardState = 1 // 活动完成
	actRecord.Rank = rank
}

// 过期检查
func (self *ModGeneral) resetInfo() {
	pInfo := &self.Sql_General
	lootConfig, _ := self.GetLootConfig()
	if lootConfig == nil {
		return
	}
	curKeyId := GetGeneralMgr().getKeyId()
	if pInfo.actRecord.KeyId != curKeyId {
		self.initInfo()
		pInfo.actRecord.KeyId = curKeyId
		pInfo.actRecord.AwardState = 0
		pInfo.actRecord.Rank = 0
		pInfo.actRecord.Times = 0
	}
}

func (self *ModGeneral) SendInfo() {
	self.checkBoxMail()
	self.checkMail()
	self.resetInfo()
	self.getGeneral()
}

func (self *ModGeneral) checkBoxMail() {
	if self.Sql_General.Point <= 0 {
		return
	}

	actRecord := self.Sql_General.actRecord
	if actRecord == nil {
		return
	}

	// 检查活动有没有结束?
	endTime := GetGeneralMgr().GetCheckActTime()
	if TimeServer().Unix() <= endTime {
		return
	}

	lootConfig, _ := self.GetLootConfig()
	if lootConfig == nil {
		return
	}

	configLst := GetCsvMgr().TimeGeneralsPointsConfig
	var resConfig []*TimeGeneralsPointsConfig
	for index, elem := range configLst {
		if elem.Group != lootConfig.Id {
			continue
		}
		resConfig = append(resConfig, configLst[index])
	}

	point := self.Sql_General.Point
	for index := range self.Sql_General.info {
		// 检查积分并且领取
		targetConfig := resConfig[index]
		needScore := targetConfig.Points
		if point < needScore {
			continue
		}

		pInfo := self.Sql_General.info[index]
		if pInfo == nil {
			continue
		}

		if pInfo.State != 0 {
			continue
		}

		// 拿到奖励配置
		award := targetConfig.Award
		if award == nil {
			continue
		}

		nums := targetConfig.Nums
		if len(award) != len(nums) {
			continue
		}

		var res []PassItem
		for i := range award {
			res = append(res, PassItem{award[i], nums[i]})
		}

		if res == nil || len(res) <= 0 {
			continue
		}

		pMail := self.player.GetModule("mail").(*ModMail)
		if pMail == nil {
			LogError("checkMail in mod mail, pMail == nil!")
			return
		}

		pMail.AddMail(1, 1, 0, GetCsvMgr().GetText("STR_GENERAL_MAIL_TITLE"),
			fmt.Sprintf(GetCsvMgr().GetText("STR_GENERAL_MAIL_CONTENT"), point),
			GetCsvMgr().GetText("STR_SYS"), res, false, 0)
		pInfo.State = 1
	}
}

func (self *ModGeneral) getVipNum() int {
	vip := self.player.Sql_UserBase.Vip
	config, ok := GetCsvMgr().VipConfigMap[vip]
	if !ok {
		LogError("vip error in getVipNum, vip:", vip)
		return 0
	}

	if config == nil {
		LogError("config is nil")
		return 0
	}

	return config.TimeGeneralsnum
}

// 抽奖次数判断
func (self *ModGeneral) checkVipTimes(n int) bool {
	limit := self.getVipNum()
	if limit == -1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_TIME-LIMITED_GOD_MISALLOCATED_THE_NUMBER"))
		return false
	}

	if self.Sql_General.actRecord.Times >= limit {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GENERAL_HAS_REACHED_THE_MAXIMUM_NUMBER"))
		return false
	}
	return true
}

func (self *ModGeneral) getRankNew() {
	GetGeneralMgr().GetRankNew(self.player)
}

func (self *ModGeneral) LootTimesDrop() int {
	value := 10

	return value
}

func (self *ModGeneral) LootItemNew(lst *PubchesttotalConfig, findtype int, findtimes int, certaintimes int, bag []PassItem) ([]PassItem, []*GeneralRecord) {
	records := make([]*GeneralRecord, 0)

	temp := make([]PassItem, 0)
	for j := 0; j < 3; j++ {
		//LogDebug("j:", j, ", len Dropgroups:", len(lst[i].Dropgroups), ", len Dropgroupids", len(lst[i].Dropgroupids))
		bagitem, bagnum := lst.Dropgroups[j], lst.Dropgroupids[j]
		//LogDebug("drop item item:", bagitem, bagnum)
		if bagitem > 0 && bagnum > 0 {
			var dropitem PassItem
			dropitem.ItemID = bagitem
			dropitem.Num = bagnum
			temp = append(temp, dropitem)
		}
	}

	if len(bag) == 0 {
		bag = temp
	} else {
		size := len(temp)
		//randIndex := self.GetRandInt(size)
		randIndex := size - 1
		temp[randIndex].Num -= 1
		var dropitem PassItem
		dropitem.ItemID = lst.Dropgroups[3]
		dropitem.Num = 1
		temp = append(temp, dropitem)
		bag = temp
	}

	modFind := self.player.GetModule("find").(*ModFind)
	if modFind == nil {
		LogError("module find is nil!")
		return []PassItem{}, records
	}

	// 计算掉落物品
	item := make([]PassItem, 0)
	for i := 0; i < len(bag); i++ {
		droplist, record := modFind.FindDrop(bag[i].ItemID, bag[i].Num)

		item = append(item, droplist...)
		records = append(records, record...)
	}

	return item, records
}
