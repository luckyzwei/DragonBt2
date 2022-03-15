package game

import (
	"encoding/json"
	"fmt"
	"strings"
	//"time"
)

const (
	DRAW_TYPE_GEM    = 1 // 钻石池
	DRAW_TYPE_FRIEND = 2 // 友情池
	DRAW_TYPE_CAMP_1 = 3 // 阵营召唤1
	DRAW_TYPE_CAMP_2 = 4 // 阵营召唤2
	DRAW_TYPE_CAMP_3 = 5 // 阵营召唤3
	DRAW_TYPE_CAMP_4 = 6 // 阵营召唤4
	DRAW_TYPE_BEAUTY = 7 // 圣物召唤
	DRAW_TYPE_LUCKY  = 8 // 福袋召唤
	DRAW_TYPE_END    = 9 // 基础召唤END

	DRAW_SELFSELECTION     = 8   //自选召唤单抽
	DRAW_SELFSELECTION_TEN = 108 //自选召唤十连

	CONFIG_SYSTEM_DIS = 10    //配置类型偏差值
	CAMP_DIS          = 1     //配置类型偏差值
	CAMP_DIS_TEN      = 101   //配置类型偏差值
	REWARD_INIT_TYPE  = 1     //奖励初始化类型
	REWARD_INIT_ORDER = 1     //奖励初始化阶段
	WISH_MAX          = 5     //许愿的最大数量
	FIND_TODAY_LIMIT  = 50000 //每天限制  提审30  正式服50000

	ASTROLOGY_ADD    = 500
	HERO_GROUP       = 10001
	HERO_GROUP_SCORE = -10000
	RAND_FIND_MAX    = 2000
)

//对应Pub_Chest_Total表：pay_pub_type
const (
	REALTYPE_GEM_ITEM     = 1   // 钻石物品单抽
	REALTYPE_GEM_ITEM_TEN = 101 // 钻石物品十连
	REALTYPE_GEM          = 2   // 钻石单抽
	REALTYPE_GEM_TEN      = 102 // 钻石十连
	REALTYPE_FRIEND       = 3   // 友情单抽
	REALTYPE_FRIEND_TEM   = 103 // 友情十连
	REALTYPE_CAMP_1       = 4   // 阵营3单抽
	REALTYPE_CAMP_1_TEN   = 104 // 阵营3十连
	REALTYPE_CAMP_2       = 5   // 阵营4单抽
	REALTYPE_CAMP_2_TEN   = 105 // 阵营4十连
	REALTYPE_CAMP_3       = 6   // 阵营5单抽
	REALTYPE_CAMP_3_TEN   = 106 // 阵营5十连
	REALTYPE_CAMP_4       = 7   // 阵营6单抽
	REALTYPE_CAMP_4_TEN   = 107 // 阵营6十连
	REALTYPE_BEAUTY       = 17  // 圣物单抽
	REALTYPE_BEAUTY_FIVE  = 19  // 圣物五连
)

const (
	START_INDEX			  = 17 	// simpleConfig起始索引
)

type FindPool struct {
	Type           int         `json:"type"`
	FindTimes      int         `json:"findtimes"`
	StartTime      int64       `json:"starttime"`
	EndTime        int64       `json:"endtime"`
	FindTimesToday int         `json:"findtimestoday"`
	FindTimesCount map[int]int `json:"findtimescount"`
	FreeNextTime   int64       `json:"freenexttime"`
}

type FindRewardInfo struct {
	Type  int `json:"type"`
	Order int `json:"order"`
	Scale int `json:"scale"`
}

type LuckyPassItem struct {
	ItemID int `json:"itemid"` // 道具ID
	Num    int `json:"num"`    // 道具数量
	Type   int `json:"type"`   // 物品等级
}

type FindAstrology struct {
	Id        int   `json:"id"`
	FindTimes int   `json:"findtimes"` //本期次数
	StartTime int64 `json:"starttime"` //开始时间
	EndTime   int64 `json:"endtime"`   //结束时间
	HeroId    int   `json:"heroid"`    //当前选中英雄
	Score     int   `json:"score"`     //
}

type SelfSelection struct {
	Id             int         `json:"id"`
	FindTimes      int         `json:"findtimes"` //本期次数
	StartTime      int64       `json:"starttime"` //开始时间
	EndTime        int64       `json:"endtime"`   //结束时间
	HeroId         int         `json:"heroid"`    //当前选中英雄
	GetTimes       int         `json:"gettimes"`  //本期获得次数
	MaxTimes       int         `json:"maxtimes"`  //最大次数
	FindTimesCount map[int]int `json:"findtimescount"`
}

type FindWishInfo struct {
	Type     int           `json:"type"`
	WishList [WISH_MAX]int `json:"wishlist"`
}

// 高级招募累计奖励
type GemFindGift struct {
	Process    int           `json:"process"`		// 累计招募进度
	RewardList  []int 		 `json:"rewardlist"`	// 奖励领取情况
}

// 抽卡数据库
type San_Find struct {
	Uid           int64
	BaseFindInfo  string // 基础召唤信息
	RewardInfo    string // 奖励信息
	WishInfo      string //
	Astrology     string // 占星
	SelfSelection string // 自选召唤
	GemFindGift	  string	// 高级招募累计奖励情况:整体进度、领取情况

	baseFindInfo  []*FindPool
	rewardInfo    *FindRewardInfo
	wishInfo      []*FindWishInfo
	astrology     *FindAstrology
	selfSelection *SelfSelection
	gemfindgift	  *GemFindGift
	DataUpdate
}

// 抽卡
type ModFind struct {
	player    *Player
	Sql_Find  San_Find
	RandData  [RAND_FIND_MAX]int
	RandIndex int
}

func (self *ModFind) OnGetData(player *Player) {
	self.player = player
}

func (self *ModFind) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_userfindpool` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Find, "san_userfindpool", self.player.ID)
	if self.Sql_Find.Uid <= 0 {
		self.Sql_Find.Uid = self.player.ID
		self.Encode()
		InsertTable("san_userfindpool", &self.Sql_Find, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_Find.Init("san_userfindpool", &self.Sql_Find, true)
	self.RandIndex = 0
}

func (self *ModFind) GetRandInt(num int) int {
	if num <= 0 {
		LogError("GetRandInt:0")
		return 0
	}

	if self.RandIndex == 0 || self.RandData[0] == 0 || self.RandIndex >= RAND_FIND_MAX {
		for i := 0; i < RAND_FIND_MAX; i++ {
			self.RandData[i] = HF_GetRandom(10000)
		}
		self.RandIndex = 0
	}

	if (num > 10000) {
		return HF_GetRandom(num)
	}

	ret := self.RandData[self.RandIndex] % num
	self.RandIndex++
	//LogInfo("GetRand : ", self.RandData[self.RandIndex-1], self.RandIndex, num, ret)
	return ret
}

func (self *ModFind) CheckPool() {
	if len(self.Sql_Find.baseFindInfo) < DRAW_TYPE_END-1 {
		size := len(self.Sql_Find.baseFindInfo)
		for i := size; i < DRAW_TYPE_END-1; i++ {
			newPool := self.NewFindPool(i)
			self.Sql_Find.baseFindInfo = append(self.Sql_Find.baseFindInfo, newPool)
		}
	}

	if self.Sql_Find.rewardInfo == nil {
		self.Sql_Find.rewardInfo = new(FindRewardInfo)
		self.Sql_Find.rewardInfo.Type = REWARD_INIT_TYPE
		self.Sql_Find.rewardInfo.Order = REWARD_INIT_ORDER
	}

	//
	if len(self.Sql_Find.wishInfo) != HERO_ATTRIBUTE_EARTH-HERO_ATTRIBUTE_WATER+1 {
		self.Sql_Find.wishInfo = make([]*FindWishInfo, 0)
		for i := HERO_ATTRIBUTE_WATER; i <= HERO_ATTRIBUTE_EARTH; i++ {
			wish := new(FindWishInfo)
			wish.Type = i
			index := 0
			for _, v := range GetCsvMgr().HeroConfigMap {
				info, ok := v[1]
				if ok && info.Attribute == i && info.WishHero == LOGIC_TRUE {
					wish.WishList[index] = info.HeroId
					index++
					if index >= WISH_MAX {
						break
					}
				}
			}
			self.Sql_Find.wishInfo = append(self.Sql_Find.wishInfo, wish)
		}
	}
	//占星
	if self.Sql_Find.astrology == nil {
		self.Sql_Find.astrology = new(FindAstrology)
		self.Sql_Find.astrology.Score = GetCsvMgr().getInitNum(ASTROLOGY_INIT_SCORE)
	}
	//自选召唤
	if self.Sql_Find.selfSelection == nil {
		self.Sql_Find.selfSelection = new(SelfSelection)
		self.Sql_Find.selfSelection.FindTimesCount = make(map[int]int)
	}

	if self.Sql_Find.selfSelection.EndTime < TimeServer().Unix() {
		PrivilegeValue := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValues()
		self.CalSelfFind(PrivilegeValue, false)
	}

	// 累计招募奖励初始化
	if self.Sql_Find.gemfindgift == nil {
		self.Sql_Find.gemfindgift = new(GemFindGift)
		self.Sql_Find.gemfindgift.Process = 0
		self.Sql_Find.gemfindgift.RewardList = []int{0,0,0}
	}

	//如果4个阵营池子都没开放，则重新生成时间
	needCal := true
	now := TimeServer().Unix()
	for _, v := range self.Sql_Find.baseFindInfo {
		if now >= v.StartTime && now <= v.EndTime {
			needCal = false
			break
		}
	}
	if needCal {
		for _, v := range self.Sql_Find.baseFindInfo {
			v.CalTime()
		}
	}
}

func (self *ModFind) NewFindPool(index int) *FindPool {
	pool := new(FindPool)
	pool.Type = index + 1
	pool.FindTimesCount = make(map[int]int, 0)
	pool.CalTime()
	return pool
}

func (self *FindPool) CalTime() {
	switch self.Type {
	case DRAW_TYPE_CAMP_1, DRAW_TYPE_CAMP_2, DRAW_TYPE_CAMP_3, DRAW_TYPE_CAMP_4:
		system := self.Type + CONFIG_SYSTEM_DIS
		self.StartTime, self.EndTime = GetCsvMgr().GetNowStartAndEnd(system)
	}
}

func (self *FindPool) getTimesAll(realType int) int {
	times := 0
	for _, v := range GetCsvMgr().PubchestspecialConfig {
		items := strings.Split(v.Paytype, "|")
		isHas := false
		for i := 0; i < len(items); i++ {
			pay_type := HF_Atoi(items[i])
			if pay_type == realType {
				isHas = true
			}
		}
		if isHas {
			for i := 0; i < len(items); i++ {
				pay_type := HF_Atoi(items[i])
				times += self.FindTimesCount[pay_type]
			}
			break
		}
	}
	return times
}

func (self *ModFind) onReg(handlers map[string]func(body []byte)) {
	handlers["findpool"] = self.FindPool
	handlers["findsavewish"] = self.FindSaveWish
	handlers["findopencamp"] = self.FindOpenCamp
	handlers["astrologyhero"] = self.AstrologyHero
	handlers["findastrology"] = self.FindAstrology
	handlers["selfselectionhero"] = self.SelfSelectionHero
	handlers["findselfselection"] = self.FindSelfSelection
	handlers["getselfselection"] = self.GetSelfSelection
	handlers["getfindreward"] = self.GetFindReward
}

// 领取累积招募奖励
func (self *ModFind) GetFindReward(body []byte)  {
	var msg C2S_GetFindReward
	json.Unmarshal(body, &msg)

	// 获取奖励配置
	config := GetCsvMgr().SimpleConfigMap[START_INDEX + msg.RewardId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	// 若不够资格或已领取,直接返回
	if self.Sql_Find.gemfindgift.RewardList[msg.RewardId] != 0 ||
		self.Sql_Find.gemfindgift.Process < config.Num{
			return
		}



	// 奖励加进背包
	getItems := self.player.AddObjectLst(config.ItemsId, config.ItemsNum, "累计招募奖励", msg.RewardId, 0, 0)

	// 同步奖励领取状态
	self.Sql_Find.gemfindgift.RewardList[msg.RewardId] = LOGIC_TRUE

	var res S2C_GetFindReward
	res.Cid = "getfindreward"
	res.FindGiftProcess = self.Sql_Find.gemfindgift.Process
	res.IsGotFindGift = self.Sql_Find.gemfindgift.RewardList
	res.GetItems = getItems
	self.player.SendMsg(res.Cid, HF_JtoB(&res))
}

// 招募入口
func (self *ModFind) FindPool(body []byte) {

	var msg C2S_FindPool
	json.Unmarshal(body, &msg)

	if !self.player.GetModule("hero").(*ModHero).CheckHeroBuyPos() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_NUM_MAX"))
		return
	}

	switch msg.Findtype {
	case DRAW_TYPE_GEM:
		self.FindGemPool(&msg)		// 高级招募
	case DRAW_TYPE_FRIEND:
		self.FindFriendPool(&msg)	// 友情招募
	case DRAW_TYPE_CAMP_1, DRAW_TYPE_CAMP_2, DRAW_TYPE_CAMP_3, DRAW_TYPE_CAMP_4:	// 阵营招募
		self.FindCampPool(&msg)
	case DRAW_TYPE_BEAUTY:
		self.FindBeautyPool(&msg)
	case DRAW_TYPE_LUCKY:
		self.FindLuckyPool(&msg)
	default:
		return
	}
}

//占星
func (self *ModFind) FindAstrology(body []byte) {

	//看开启条件是否满足
	lastpass := self.player.GetModule("pass").(*ModPass).GetLastPass()
	passId := ONHOOK_INIT_LEVEL
	if lastpass != nil {
		passId = lastpass.Id
	}
	if !GetCsvMgr().IsLevelAndPassOpenNew(self.player.Sql_UserBase.Level, passId, OPEN_ASTROLOGY) {
		vipcsv, ok := GetCsvMgr().VipConfigMap[self.player.Sql_UserBase.Vip]
		if !ok || vipcsv.Astrologer == LOGIC_FALSE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_UIPUBLIC_SYS_NOT_OPEN"))
			return
		}
	}

	var msg C2S_FindAstrology
	json.Unmarshal(body, &msg)

	if msg.FindNum != 1 && msg.FindNum != 10 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TYPE_ERROR"))
		return
	}

	if !self.player.GetModule("hero").(*ModHero).CheckHeroBuyPos() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_NUM_MAX"))
		return
	}

	if !self.IsAstrologyHero() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}
	//扣除物品
	items := make([]PassItem, 0)
	getItems := make([]PassItem, 0)
	costItems := make([]PassItem, 0)
	//检查消耗够不够  先看占星卷  在看钻石
	configItem := GetCsvMgr().GetTariffConfig2(TARIFF_ASTROLOGY_ITEM)
	if configItem == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	configGem := GetCsvMgr().GetTariffConfig2(TARIFF_ASTROLOGY_GEN)
	if configGem == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	calCostItem := make([]int, 0)
	calCostNum := make([]int, 0)
	for i := 0; i < len(configItem.ItemIds); i++ {
		if configItem.ItemIds[i] == 0 {
			break
		}
		calCostItem = append(calCostItem, configItem.ItemIds[i])
		calCostNum = append(calCostNum, configItem.ItemNums[i]*msg.FindNum)
	}
	if err := self.player.HasObjectOk(calCostItem, calCostNum); err != nil {
		calCostItem = make([]int, 0)
		calCostNum = make([]int, 0)
		for i := 0; i < len(configGem.ItemIds); i++ {
			if configGem.ItemIds[i] == 0 {
				break
			}
			calCostItem = append(calCostItem, configGem.ItemIds[i])
			calCostNum = append(calCostNum, configGem.ItemNums[i]*msg.FindNum)
		}
		if err := self.player.HasObjectOk(calCostItem, calCostNum); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
	}

	costItems = self.player.RemoveObjectLst(calCostItem, calCostNum, "占星", msg.FindNum, 0, 0)

	//计算产出
	for i := 0; i < msg.FindNum; i++ {
		item := self.CalAstrologyDrop()
		self.Sql_Find.astrology.FindTimes += 1
		items = append(items, item)
	}

	param2 := self.Sql_Find.astrology.Score * 1000
	getItemTran := make(map[int]*Item)
	for i := 0; i < len(items); i++ {
		configItem := GetCsvMgr().ItemMap[items[i].ItemID]
		if configItem == nil {
			continue
		}
		if configItem.ItemType == ITEM_TYPE_HERO {
			self.player.AddObject(items[i].ItemID, items[i].Num, msg.FindNum, param2, 0, "占星")
		} else {
			AddItemMapHelper3(getItemTran, items[i].ItemID, items[i].Num)
		}
	}
	self.player.AddObjectItemMap(getItemTran, "占星", msg.FindNum, param2, 0)

	self.player.HandleTask(TASK_TYPE_ASTROLOGY_COUNT, msg.FindNum, 0, 0)

	var msgRel S2C_FindAstrology
	msgRel.Cid = "findastrology"
	msgRel.FindNum = msg.FindNum
	msgRel.Item = items
	msgRel.CostItems = costItems
	msgRel.GetItems = getItems
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	logDec := ""
	logId := 0
	if msg.FindNum == 1 {
		logDec = "单次占星"
		logId = LOG_ASTROLOGY_ONE
	} else if msg.FindNum == 10 {
		logDec = "十连占星"
		logId = LOG_ASTROLOGY_TEN
	}
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, logId, 1, 0, 0, logDec, 0, 0, self.player)
}

func (self *ModFind) FindSelfSelection(body []byte) {

	var msg C2S_FindSelfSelection
	json.Unmarshal(body, &msg)

	realType := 0
	dec := ""
	if msg.FindNum == 1 {
		realType = DRAW_SELFSELECTION
		dec = "单次自选招募"
	} else if msg.FindNum == 10 {
		realType = DRAW_SELFSELECTION_TEN
		dec = "十连自选招募"
	} else {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TYPE_ERROR"))
		return
	}

	if !self.player.GetModule("hero").(*ModHero).CheckHeroBuyPos() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HERO_NUM_MAX"))
		return
	}

	if !self.IsSelfSelectionHero() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	// 找到符合条件的掉落组
	lst := GetCsvMgr().PubchesttotalGroup[realType]
	if len(lst) == 0 {
		LogDebug("找不到抽卡类型：", realType)
		return
	}

	//检查消耗够不够
	payitem, payitemnum := lst[0].Payitem, lst[0].Payitemnum
	if payitem != 0 {
		if err := self.player.HasObjectOkEasy(payitem, payitemnum); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
	}

	// 计算必掉掉落
	bag := make([]PassItem, 0)
	certaintimes := lst[0].Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(lst[0].Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := lst[0].Dropgroups[3]
		tempCer := self.Sql_Find.selfSelection.FindTimes % certaintimes
		tempCer += msg.FindNum
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
	param1 := LOGIC_FALSE
	if self.Sql_Find.selfSelection.HeroId != 0 {
		param1 = LOGIC_TRUE
	}
	if self.Sql_Find.selfSelection.FindTimesCount == nil {
		self.Sql_Find.selfSelection.FindTimesCount = make(map[int]int, 0)
	}
	self.Sql_Find.selfSelection.FindTimes += msg.FindNum
	self.Sql_Find.selfSelection.FindTimesCount[realType]++

	//times := findPool.getTimesAll(realType)
	item := self.LootItems(lst, realType, self.Sql_Find.selfSelection.FindTimes, certaintimes, bag, self.Sql_Find.selfSelection.FindTimes)
	//扣除消耗
	costItems := self.player.RemoveObjectSimple(payitem, payitemnum, dec, param1, 0, 1)
	num := GetGemNum(costItems)
	//发送英雄，如果开启了转换开关，就自动分解
	getItem := make(map[int]*Item)
	getItemTran := make(map[int]*Item) //转换
	//计算任务进度
	self.CalRewardTimes(getItem, msg.FindNum)

	param3 := 0
	if num > 0 {
		param3 = -1
	}
	nCount := 0
	nDecompose := 0
	param2 := self.Sql_Find.selfSelection.FindTimesCount[realType]*msg.FindNum*1000 + realType

	//这个地方是自选池特性
	for i := 0; i < len(item); i++ {
		if self.Sql_Find.selfSelection.GetTimes >= self.Sql_Find.selfSelection.MaxTimes {
			break
		}
		config := GetCsvMgr().GetItemConfig(item[i].ItemID)
		if nil != config {
			if config.ItemCheck >= 4 {
				nCount++
			} else {
				continue
			}
		}
		newItem := self.Sql_Find.selfSelection.HeroId*100 + 11000001
		item[i].ItemID = newItem
		self.Sql_Find.selfSelection.GetTimes++
	}

	for i := 0; i < len(item); i++ {
		//如果用户打开了自动分解，需要判断这个物品是否转化为碎片
		itemId, itemNum := self.player.GetModule("hero").(*ModHero).CheckItem(item[i].ItemID, item[i].Num)
		if len(itemId) > 0 {
			AddItemMapHelper(getItemTran, itemId, itemNum)
			nDecompose++
		} else {
			item[i].ItemID, item[i].Num = self.player.AddObject(item[i].ItemID, item[i].Num, param1, param2, 0, dec)
		}
		config := GetCsvMgr().GetItemConfig(item[i].ItemID)
		if nil != config {
			if config.ItemCheck >= 4 {
				nCount++
			}
		}
	}
	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, msg.FindNum, 8, 0) //
	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, msg.FindNum, 1, 0)
	self.player.HandleTask(TASK_TYPE_SUMMON_ELITE_HEROS, nCount, realType, 0)
	self.player.HandleTask(TASK_TYPE_DECOMPOSE_HEROS, nDecompose, 0, 0)
	getItems := self.player.AddObjectItemMap(getItem, dec, realType, self.Sql_Find.selfSelection.FindTimes, param3)
	getItemsTran := self.player.AddObjectItemMap(getItemTran, dec, realType, self.Sql_Find.selfSelection.FindTimes, param3)

	if num > 0 {
		var logitem []PassItem
		logitem = append(getItems, getItemsTran...)
		AddSpecialSdkItemListLog(self.player, num, logitem, dec)
	}

	if msg.FindNum == 1 {
		num := self.Sql_Find.selfSelection.MaxTimes - self.Sql_Find.selfSelection.GetTimes
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_SELF_ONE, msg.FindNum, num, self.Sql_Find.selfSelection.HeroId, dec, 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_SELF_TEN, msg.FindNum, num, self.Sql_Find.selfSelection.HeroId, dec, 0, 0, self.player)
	}

	var msgRel S2C_FindSelfSelection
	msgRel.Cid = "findselfselection"
	msgRel.FindNum = msg.FindNum
	msgRel.Item = item
	msgRel.CostItems = costItems
	msgRel.GetItems = getItems
	msgRel.GetTimes = self.Sql_Find.selfSelection.GetTimes
	msgRel.RewardInfo = self.Sql_Find.rewardInfo
	msgRel.GetItemsTran = getItemsTran
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFind) GetSelfSelection(body []byte) {
	if self.Sql_Find.selfSelection.EndTime < TimeServer().Unix() {
		PrivilegeValue := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValues()
		self.CalSelfFind(PrivilegeValue, false)
	}

	var msgRel S2C_GetSelfSelection
	msgRel.Cid = "getselfselection"
	msgRel.SelfSelection = self.Sql_Find.selfSelection
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFind) CalAstrologyDrop() PassItem {
	if self.Sql_Find.astrology == nil {
		self.Sql_Find.astrology = new(FindAstrology)
		self.Sql_Find.astrology.Score = GetCsvMgr().getInitNum(ASTROLOGY_INIT_SCORE)
	}

	var item PassItem

	//增加自己的权值分
	self.Sql_Find.astrology.Score += ASTROLOGY_ADD

	config := GetCsvMgr().AstrologyDropGroupConfig
	//增加占星保底设置   20200108
	lst := GetCsvMgr().PubChestSpecialLst[10000]
	for i := 0; i < len(lst); i++ {
		if self.Sql_Find.astrology.FindTimes+1 >= lst[i].Droptimemin && self.Sql_Find.astrology.FindTimes+1 <= lst[i].Droptimemax {
			config = make(map[int]*AstrologyDropConfig, 0)
			specail := lst[i].DropGroupModify
			items := strings.Split(specail, "|")

			for j := 0; j < len(items); j++ {
				item := strings.Split(items[j], ":")
				if len(item) >= 2 {
					config[HF_Atoi(item[0])] = GetCsvMgr().AstrologyDropGroupConfig[HF_Atoi(item[0])]
				}
			}
		}
	}
	//根据权值分数计算大组的总权值
	groupRateAll := 0
	for _, v := range config {
		if v.ScoreLimit != 0 && self.Sql_Find.astrology.Score >= v.ScoreLimit {
			groupRateAll += v.AstrologyMid
		} else {
			groupRateAll += v.AstrologyChance
		}
	}
	if groupRateAll <= 0 {
		return item
	}
	//获得随机数值
	groupRand := self.GetRandInt(groupRateAll)
	groupRateCal := 0
	//计算选中组
	for _, v := range config {
		if v.ScoreLimit != 0 && self.Sql_Find.astrology.Score >= v.ScoreLimit {
			groupRateCal += v.AstrologyMid
		} else {
			groupRateCal += v.AstrologyChance
		}
		//选中
		if groupRateCal > groupRand {
			//如果是目标掉落组
			if v.AstrologyId == HERO_GROUP {
				item.Num = v.ItemNum
				item.ItemID = 11000000 + self.Sql_Find.astrology.HeroId*100 + 1

				score := self.GetAstrologyHeroScore(item.ItemID)
				self.Sql_Find.astrology.Score += score
				return item
			} else {
				info, ok := GetCsvMgr().AstrologyDropConfig[v.AstrologyId]
				if ok {
					rateAll := 0
					for _, vv := range info {
						rateAll += vv.ItemWT
					}
					if rateAll <= 0 {
						return item
					}
					//获得随机数值
					rand := self.GetRandInt(rateAll)
					rateCal := 0
					for _, vv := range info {
						rateCal += vv.ItemWT
						if rateCal > rand {
							item.Num = vv.ItemNum
							item.ItemID = vv.ItemId
							if vv.ScoreValue == LOGIC_FALSE {
								self.Sql_Find.astrology.Score += vv.ItemScore
							} else {
								self.Sql_Find.astrology.Score -= vv.ItemScore
							}
							return item
						}
					}
				}
			}
		}
	}

	return item
}

func (self *ModFind) FindSaveWish(body []byte) {
	var msg C2S_FindSaveWish
	json.Unmarshal(body, &msg)

	wishInfo := self.GetWishInfo(msg.Camp)
	if wishInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_WISHINFO_NOT_EXIST"))
		return
	}

	if len(msg.Ids) != WISH_MAX {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_MSG_ERROR"))
		return
	}

	count := 0
	for _, v := range wishInfo.WishList {
		if v != 0 {
			count++
		}
	}

	for _, v := range wishInfo.WishList {
		if v == 0 {
			continue
		}
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_WISH_CANCEL, v, msg.Camp, count, "解除心愿单英雄", 0, 0, self.player)
	}

	count = 0
	for _, v := range msg.Ids {
		if v != 0 {
			count++
		}
	}

	for i := 0; i < len(msg.Ids); i++ {
		wishInfo.WishList[i] = msg.Ids[i]
		if msg.Ids[i] != 0 {
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_WISH_SET, msg.Ids[i], msg.Camp, count, "设置心愿单英雄", 0, 0, self.player)
		}
	}

	var msgRel S2C_FindSaveWish
	msgRel.Cid = "findsavewish"
	msgRel.WishInfo = wishInfo
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFind) FindOpenCamp(body []byte) {
	var msg C2S_FindOpenCamp
	json.Unmarshal(body, &msg)

	if msg.Findtype < DRAW_TYPE_CAMP_1 || msg.Findtype > DRAW_TYPE_CAMP_4 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_MSG_ERROR"))
		return
	}

	findPool := self.GetFindPool(msg.Findtype)
	if findPool == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TYPE_ERROR"))
		return
	}

	now := TimeServer().Unix()
	if now > findPool.StartTime && now < findPool.EndTime {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_POOL_ALREADY_OPEN"))
		return
	}

	configCost := GetCsvMgr().GetTariffConfig2(TARIFF_TYPE_FIND_OPEN_CAMP)
	if configCost == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	isFree := false
	if self.player.GetModule("nobilitytask").(*ModNobilityTask).GetNobilityPrivilege(3) {
		isFree = true
	}

	//检查消耗够不够
	if !isFree {
		if err := self.player.HasObjectOk(configCost.ItemIds, configCost.ItemNums); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
	}
	costItem := make([]PassItem, 0)
	// 扣除物品
	if !isFree {
		costItem = self.player.RemoveObjectLst(configCost.ItemIds, configCost.ItemNums, "阵营招募激活阵营", msg.Findtype, 0, 0)
	}
	//结束时间为了更灵活配置，硬核处理
	for _, v := range self.Sql_Find.baseFindInfo {
		if now > v.StartTime && now < v.EndTime {
			findPool.StartTime = now
			findPool.EndTime = v.EndTime
		}
	}

	var msgRel S2C_FindOpenCamp
	msgRel.Cid = "findopencamp"
	msgRel.FindPool = findPool
	msgRel.CostItems = costItem
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_CAMP_OPEN, msg.Findtype, 0, 0, "阵营招募激活阵营", 0, 0, self.player)
}

func (self *ModFind) AstrologyHero(body []byte) {
	var msg C2S_AstrologyHero
	json.Unmarshal(body, &msg)

	oldHeroId := self.Sql_Find.astrology.HeroId
	self.Sql_Find.astrology.HeroId = msg.Id
	if !self.IsAstrologyHero() {
		self.Sql_Find.astrology.HeroId = 0
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ASTROLOGY_CHANGE, self.Sql_Find.astrology.HeroId, oldHeroId, 0, "更换占星目标英雄", 0, 0, self.player)

	var msgRel S2C_AstrologyHero
	msgRel.Cid = "astrologyhero"
	msgRel.Astrology = self.Sql_Find.astrology
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFind) SelfSelectionHero(body []byte) {
	var msg C2S_SelfSelectionHero
	json.Unmarshal(body, &msg)

	old := self.Sql_Find.selfSelection.HeroId
	self.Sql_Find.selfSelection.HeroId = msg.Id
	if !self.IsSelfSelectionHero() {
		self.Sql_Find.astrology.HeroId = 0
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	var msgRel S2C_SelfSelectionHero
	msgRel.Cid = "selfselectionhero"
	msgRel.SelfSelection = self.Sql_Find.selfSelection
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_SELF_CHANGE, self.Sql_Find.selfSelection.HeroId, old, 0, "更换自选目标英雄", 0, 0, self.player)

}

func (self *ModFind) IsAstrologyHero() bool {
	config, ok := GetCsvMgr().AstrologyDropConfig[HERO_GROUP]
	if !ok {
		return false
	}
	itemId := 11000000 + self.Sql_Find.astrology.HeroId*100 + 1
	for _, v := range config {
		if v.ItemId == itemId {
			return true
		}
	}
	return false
}

func (self *ModFind) IsSelfSelectionHero() bool {
	config, ok := GetCsvMgr().HeroConfigMap[self.Sql_Find.selfSelection.HeroId][1]
	if !ok {
		return false
	}
	if config.Attribute >= HERO_ATTRIBUTE_WATER && config.Attribute <= HERO_ATTRIBUTE_EARTH {
		return true
	}
	return false
}

func (self *ModFind) GetAstrologyHeroScore(itemId int) int {
	config, ok := GetCsvMgr().AstrologyDropConfig[HERO_GROUP]
	if !ok {
		return HERO_GROUP_SCORE
	}

	for _, v := range config {
		if v.ItemId == itemId {
			if v.ScoreValue == LOGIC_FALSE {
				return v.ItemScore
			} else {
				return -v.ItemScore
			}
		}
	}
	return HERO_GROUP_SCORE
}

func (self *ModFind) GetWishInfo(wishType int) *FindWishInfo {
	for _, v := range self.Sql_Find.wishInfo {
		if v.Type == wishType {
			return v
		}
	}
	return nil
}

func (self *ModFind) GetFindPool(findType int) *FindPool {
	for _, v := range self.Sql_Find.baseFindInfo {
		if v.Type == findType {
			return v
		}
	}
	return nil
}

func (self *ModFind) CalRewardTimes(items map[int]*Item, times int) {
	if items == nil {
		return
	}

	config := GetCsvMgr().GetSummonBoxConfig(self.Sql_Find.rewardInfo.Type, self.Sql_Find.rewardInfo.Order)
	if config == nil {
		return
	}
	self.Sql_Find.rewardInfo.Scale += times
	//发奖判断
	if self.Sql_Find.rewardInfo.Scale >= config.Scale {
		self.Sql_Find.rewardInfo.Scale -= config.Scale
		AddItemMapHelper(items, config.Item, config.Num)
		//更新任务
		//判断有限任务下一阶段
		nextOrderConfig := GetCsvMgr().GetSummonBoxConfig(self.Sql_Find.rewardInfo.Type, self.Sql_Find.rewardInfo.Order+1)
		if nextOrderConfig != nil {
			self.Sql_Find.rewardInfo.Type = nextOrderConfig.Type
			self.Sql_Find.rewardInfo.Order = nextOrderConfig.Order
			return
		}
		//判断无限任务下一阶段
		nextTypeConfig := GetCsvMgr().GetSummonBoxConfig(self.Sql_Find.rewardInfo.Type+1, REWARD_INIT_ORDER)
		if nextTypeConfig != nil {
			self.Sql_Find.rewardInfo.Type = nextTypeConfig.Type
			self.Sql_Find.rewardInfo.Order = nextTypeConfig.Order
			return
		}
		//判断无限任务循环
		nextConfig := GetCsvMgr().GetSummonBoxConfig(self.Sql_Find.rewardInfo.Type, REWARD_INIT_ORDER)
		if nextConfig != nil {
			self.Sql_Find.rewardInfo.Type = nextConfig.Type
			self.Sql_Find.rewardInfo.Order = nextConfig.Order
			return
		}
	}

	return
}

// 高级招募
func (self *ModFind) FindGemPool(msg *C2S_FindPool) {

	realType := 0
	findPool := self.GetFindPool(msg.Findtype)
	if findPool == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TYPE_ERROR"))
		return
	}

	if findPool.FindTimesToday+msg.FindNum > FIND_TODAY_LIMIT {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TODAY_LIMIT"))
		return
	}

	dec := ""
	//条件检查
	if msg.FindNum == 1 { //单抽
		if self.player.GetObjectNum(ITEM_FIND_GEM_ITEM) < 1 { // 时间
			realType = REALTYPE_GEM
		} else {
			realType = REALTYPE_GEM_ITEM
		}
		dec = "单次高级招募"
	} else if msg.FindNum == 10 {
		if self.player.GetObjectNum(ITEM_FIND_GEM_ITEM) < 10 { // 时间
			realType = REALTYPE_GEM_TEN
		} else {
			realType = REALTYPE_GEM_ITEM_TEN
		}
		dec = "十连高级招募"
	} else {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_NUM_ERROR"))
		return
	}

	// 找到符合条件的掉落组
	lst := GetCsvMgr().PubchesttotalGroup[realType]
	if len(lst) == 0 {
		LogDebug("找不到抽卡类型：", realType)
		return
	}

	//检查消耗够不够
	payitem, payitemnum := lst[0].Payitem, lst[0].Payitemnum
	if payitem != 0 {
		if err := self.player.HasObjectOkEasy(payitem, payitemnum); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
	}

	// 计算必掉掉落
	bag := make([]PassItem, 0)
	certaintimes := lst[0].Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(lst[0].Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := lst[0].Dropgroups[3]
		tempCer := findPool.FindTimes % certaintimes
		tempCer += msg.FindNum
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
	if findPool.FindTimesCount == nil {
		findPool.FindTimesCount = make(map[int]int, 0)
	}
	findPool.FindTimes += msg.FindNum
	self.Sql_Find.gemfindgift.Process += msg.FindNum 	// 同步累计招募进度
	findPool.FindTimesCount[realType]++

	times := findPool.getTimesAll(realType)
	item := self.LootItems(lst, realType, findPool.FindTimes, certaintimes, bag, times)
	//扣除消耗
	costItems := self.player.RemoveObjectSimple(payitem, payitemnum, dec, realType, 0, 1)
	num := GetGemNum(costItems)
	//发送英雄，如果开启了转换开关，就自动分解
	getItem := make(map[int]*Item)
	getItemTran := make(map[int]*Item) //转换
	//计算任务进度
	findPool.FindTimesToday += msg.FindNum
	self.CalRewardTimes(getItem, msg.FindNum)

	param3 := 0
	if num > 0 {
		param3 = -1
	}

	//self.player.HandleTask(FindTask, 2, 1, 0)
	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, msg.FindNum, 1, 0)
	nCount := 0
	nDecompose := 0
	param2 := findPool.FindTimesCount[realType]*msg.FindNum*1000 + realType
	for i := 0; i < len(item); i++ {
		//如果用户打开了自动分解，需要判断这个物品是否转化为碎片
		itemId, itemNum := self.player.GetModule("hero").(*ModHero).CheckItem(item[i].ItemID, item[i].Num)
		if len(itemId) > 0 {
			AddItemMapHelper(getItemTran, itemId, itemNum)
			nDecompose++
		} else {
			item[i].ItemID, item[i].Num = self.player.AddObject(item[i].ItemID, item[i].Num, realType, param2, 0, dec)
		}
		config := GetCsvMgr().GetItemConfig(item[i].ItemID)
		if nil != config {
			if config.ItemCheck >= 4 {
				nCount++
			}
		}
	}

	self.player.HandleTask(TASK_TYPE_SUMMON_ELITE_HEROS, nCount, realType, 0)
	self.player.HandleTask(TASK_TYPE_DECOMPOSE_HEROS, nDecompose, 0, 0)
	getItems := self.player.AddObjectItemMap(getItem, dec, realType, findPool.FindTimes, param3)
	getItemsTran := self.player.AddObjectItemMap(getItemTran, dec, realType, findPool.FindTimes, param3)

	if num > 0 {
		var logitem []PassItem
		logitem = append(getItems, getItemsTran...)
		AddSpecialSdkItemListLog(self.player, num, logitem, dec)
	}
	var msgRel S2C_FindPool
	msgRel.Cid = "findpool"
	msgRel.FindType = msg.Findtype
	msgRel.FindNum = msg.FindNum
	msgRel.FindNumToday = findPool.FindTimesToday
	msgRel.Item = item
	msgRel.CostItems = costItems
	msgRel.GetItems = getItems
	msgRel.GetItemsTran = getItemsTran
	msgRel.RewardInfo = self.Sql_Find.rewardInfo
	msgRel.GemFindProcess = self.Sql_Find.gemfindgift.Process
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	if msg.FindNum == 1 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_GEM_ONE, 1, 0, 0, dec, 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_GEM_TEN, msg.FindNum, 0, 0, dec, 0, 0, self.player)
	}
}

// 友情招募
func (self *ModFind) FindFriendPool(msg *C2S_FindPool) {

	realType := 0
	findPool := self.GetFindPool(msg.Findtype)
	if findPool == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TYPE_ERROR"))
		return
	}
	if findPool.FindTimesToday+msg.FindNum > FIND_TODAY_LIMIT {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TODAY_LIMIT"))
		return
	}
	dec := ""
	//条件检查
	if msg.FindNum == 1 { //单抽
		realType = REALTYPE_FRIEND
		dec = "单次友情招募"
	} else if msg.FindNum == 10 { //
		realType = REALTYPE_FRIEND_TEM
		dec = "十连友情招募"
	} else {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_NUM_ERROR"))
		return
	}

	// 找到符合条件的掉落组
	lst := GetCsvMgr().PubchesttotalGroup[realType]
	if len(lst) == 0 {
		LogDebug("找不到抽卡类型：", realType)
		return
	}

	//检查消耗够不够
	payitem, payitemnum := lst[0].Payitem, lst[0].Payitemnum
	if payitem != 0 {
		if err := self.player.HasObjectOkEasy(payitem, payitemnum); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
	}

	// 计算必掉掉落
	bag := make([]PassItem, 0)
	certaintimes := lst[0].Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(lst[0].Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := lst[0].Dropgroups[3]
		tempCer := findPool.FindTimes % certaintimes
		tempCer += msg.FindNum
		if tempCer/certaintimes > 0 {
			itemid = certainitem
		}

		if itemid != 0 {
			var dropitem PassItem
			dropitem.ItemID = itemid
			dropitem.Num = 1
			bag = append(bag, dropitem)
		}
	}
	if findPool.FindTimesCount == nil {
		findPool.FindTimesCount = make(map[int]int, 0)
	}
	findPool.FindTimes += msg.FindNum
	findPool.FindTimesCount[realType]++
	times := findPool.getTimesAll(realType)
	item := self.LootItems(lst, realType, findPool.FindTimes, certaintimes, bag, times)
	//扣除消耗
	costItems := self.player.RemoveObjectSimple(payitem, payitemnum, dec, realType, 0, 0)
	//发送英雄，如果开启了转换开关，就自动分解
	getItem := make(map[int]*Item)
	getItemTran := make(map[int]*Item)
	//计算任务进度     20200302友情不计算
	//self.CalRewardTimes(getItem, findPool, msg.FindNum)

	//self.player.HandleTask(FindTask, 2, 1, 0)
	findPool.FindTimes += msg.FindNum
	findPool.FindTimesToday += msg.FindNum
	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, msg.FindNum, 2, 0)
	nCount := 0
	nDecompose := 0
	param2 := findPool.FindTimesCount[realType]*msg.FindNum*1000 + realType
	for i := 0; i < len(item); i++ {
		//如果用户打开了自动分解，需要判断这个物品是否转化为碎片
		itemId, itemNum := self.player.GetModule("hero").(*ModHero).CheckItem(item[i].ItemID, item[i].Num)
		if len(itemId) > 0 {
			AddItemMapHelper(getItemTran, itemId, itemNum)
			nDecompose++
		} else {
			item[i].ItemID, item[i].Num = self.player.AddObject(item[i].ItemID, item[i].Num, realType, param2, 0, dec)
		}
		config := GetCsvMgr().GetItemConfig(item[i].ItemID)
		if nil != config {
			if config.ItemCheck >= 4 {
				nCount++
			}
		}
	}
	self.player.HandleTask(TASK_TYPE_SUMMON_ELITE_HEROS, nCount, realType, 0)
	self.player.HandleTask(TASK_TYPE_DECOMPOSE_HEROS, nDecompose, 0, 0)
	getItems := self.player.AddObjectItemMap(getItem, dec, 0, param2, 0)
	getItemsTran := self.player.AddObjectItemMap(getItemTran, dec, 0, 0, 0)
	var msgRel S2C_FindPool
	msgRel.Cid = "findpool"
	msgRel.FindType = msg.Findtype
	msgRel.FindNum = msg.FindNum
	msgRel.FindNumToday = findPool.FindTimesToday
	msgRel.Item = item
	msgRel.CostItems = costItems
	msgRel.GetItems = getItems
	msgRel.GetItemsTran = getItemsTran
	msgRel.RewardInfo = self.Sql_Find.rewardInfo
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	if msg.FindNum == 1 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_FRIEND_ONE, 1, 0, 0, dec, 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_FRIEND_TEN, msg.FindNum, 0, 0, dec, 0, 0, self.player)
	}
}

// 阵营招募
func (self *ModFind) FindCampPool(msg *C2S_FindPool) {

	//先看对应池子是否开放
	realType := 0
	findPool := self.GetFindPool(msg.Findtype)
	if findPool == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TYPE_ERROR"))
		return
	}
	/*
		if findPool.FindTimesToday+msg.FindNum > FIND_TODAY_LIMIT {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TODAY_LIMIT"))
			return
		}
	*/
	now := TimeServer().Unix()
	if now > findPool.EndTime || now < findPool.StartTime {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_NOT_OPEN"))
		return
	}
	dec := ""
	//条件检查
	if msg.FindNum == 1 { //单抽
		realType = msg.Findtype + CAMP_DIS
		dec = "单次阵营招募"
	} else if msg.FindNum == 10 { //
		realType = msg.Findtype + CAMP_DIS_TEN
		dec = "十连阵营招募"
	} else {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_NUM_ERROR"))
		return
	}

	// 找到符合条件的掉落组
	lst := GetCsvMgr().PubchesttotalGroup[realType]
	if len(lst) == 0 {
		LogDebug("找不到抽卡类型：", realType)
		return
	}

	//检查消耗够不够
	payitem, payitemnum := lst[0].Payitem, lst[0].Payitemnum
	if payitem != 0 {
		if err := self.player.HasObjectOkEasy(payitem, payitemnum); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
	}

	// 计算必掉掉落
	bag := make([]PassItem, 0)
	certaintimes := lst[0].Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(lst[0].Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := lst[0].Dropgroups[3]
		tempCer := findPool.FindTimes % certaintimes
		tempCer += msg.FindNum
		if tempCer/certaintimes > 0 {
			itemid = certainitem
		}

		if itemid != 0 {
			var dropitem PassItem
			dropitem.ItemID = itemid
			dropitem.Num = 1
			bag = append(bag, dropitem)
		}
	}
	if findPool.FindTimesCount == nil {
		findPool.FindTimesCount = make(map[int]int, 0)
	}
	findPool.FindTimes += msg.FindNum
	findPool.FindTimesCount[realType]++
	times := findPool.getTimesAll(realType)
	item := self.LootItems(lst, realType, findPool.FindTimes, certaintimes, bag, times)
	//扣除消耗
	costItems := self.player.RemoveObjectSimple(payitem, payitemnum, dec, realType, 0, 0)
	//发送英雄，如果开启了转换开关，就自动分解
	getItem := make(map[int]*Item)
	getItemTran := make(map[int]*Item)
	findPool.FindTimesToday += msg.FindNum
	self.CalRewardTimes(getItem, msg.FindNum)

	//self.player.HandleTask(FindTask, 2, 1, 0)
	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, msg.FindNum, msg.Findtype, 0)
	self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, msg.FindNum, 7, 0)
	nCount := 0
	nDecompose := 0
	param2 := findPool.FindTimesCount[realType]*msg.FindNum*1000 + realType
	for i := 0; i < len(item); i++ {
		//如果用户打开了自动分解，需要判断这个物品是否转化为碎片
		itemId, itemNum := self.player.GetModule("hero").(*ModHero).CheckItem(item[i].ItemID, item[i].Num)
		if len(itemId) > 0 {
			AddItemMapHelper(getItemTran, itemId, itemNum)
			nDecompose++
		} else {
			item[i].ItemID, item[i].Num = self.player.AddObject(item[i].ItemID, item[i].Num, realType, param2, 0, dec)
		}
		config := GetCsvMgr().GetItemConfig(item[i].ItemID)
		if nil != config {
			if config.ItemCheck >= 4 {
				nCount++
			}
		}
	}
	self.player.HandleTask(TASK_TYPE_SUMMON_ELITE_HEROS, nCount, realType, 0)
	self.player.HandleTask(TASK_TYPE_DECOMPOSE_HEROS, nDecompose, 0, 0)
	getItems := self.player.AddObjectItemMap(getItem, dec, 0, param2, 0)
	getItemsTran := self.player.AddObjectItemMap(getItemTran, dec, 0, 0, 0)
	var msgRel S2C_FindPool
	msgRel.Cid = "findpool"
	msgRel.FindType = msg.Findtype
	msgRel.FindNum = msg.FindNum
	msgRel.FindNumToday = findPool.FindTimesToday
	msgRel.Item = item
	msgRel.CostItems = costItems
	msgRel.GetItems = getItems
	msgRel.GetItemsTran = getItemsTran
	msgRel.RewardInfo = self.Sql_Find.rewardInfo
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	if msg.FindNum == 1 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_CAMP_ONE, 1, msg.Findtype, 0, dec, 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_CAMP_TEN, msg.FindNum, msg.Findtype, 0, dec, 0, 0, self.player)
	}
}

func (self *ModFind) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModFind) LootItems(lst []*PubchesttotalConfig, findtype int, findtimes int, certaintimes int, bag []PassItem, specialTimes int) []PassItem {
	//没有必出掉落，计算随机掉落
	temp := make([]PassItem, 0)
	sum := 0
	vip := self.player.Sql_UserBase.Vip
	for i := 0; i < len(lst); i++ {
		if vip < lst[i].Cardsvipmin {
			continue
		}

		if vip > lst[i].Cardsvipmax {
			continue
		}
		sum += lst[i].Chance
	}

	//随机出掉落组
	pro := self.GetRandInt(sum)
	//LogInfo("DROP PRO : =============================>>>>>>>>>>>>>>>>>", pro)
	cur := 0
	for i := 0; i < len(lst); i++ {
		if vip < lst[i].Cardsvipmin {
			continue
		}

		if vip > lst[i].Cardsvipmax {
			continue
		}
		cur += lst[i].Chance
		if pro < cur {
			for j := 0; j < 3; j++ {
				//LogDebug("j:", j, ", len Dropgroups:", len(lst[i].Dropgroups), ", len Dropgroupids", len(lst[i].Dropgroupids))
				bagitem, bagnum := lst[i].Dropgroups[j], lst[i].Dropgroupids[j]
				//LogDebug("drop item item:", bagitem, bagnum)
				if bagitem > 0 && bagnum > 0 {
					var dropitem PassItem
					dropitem.ItemID = bagitem
					dropitem.Num = bagnum
					temp = append(temp, dropitem)
				}
			}
			break
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
		dropitem.ItemID = lst[0].Dropgroups[3]
		dropitem.Num = 1
		temp = append(temp, dropitem)
		bag = temp
	}

	/*
		if len(bag) == 1 && bag[0].Num == 1 && lst[0].Dropcardsnum > 1 {
			var dropitem PassItem
			dropitem.ItemID = lst[0].Dropgroups[2]
			dropitem.Num = lst[0].Dropcardsnum - 1
			bag = append(bag, dropitem)
		}
	*/

	//首抽替换
	lstSpecail, spType := self.FindSpecial(findtype, specialTimes)
	if len(lstSpecail) > 0 {
		if spType == 1 {
			for i := 0; i < len(lstSpecail); i++ {
				for j := 0; j < len(bag); j++ {
					if bag[j].ItemID == lstSpecail[i].Original {
						if lstSpecail[i].Rate >= self.GetRandInt(10000) {
							bag[j].ItemID = lstSpecail[i].New
						}
					}
				}
			}
		} else if spType == 2 {
			bag = make([]PassItem, 0)
			for i := 0; i < len(lstSpecail); i++ {
				var dropitem PassItem
				dropitem.ItemID = lstSpecail[i].Original
				dropitem.Num = lstSpecail[i].New
				bag = append(bag, dropitem)
			}
		}
	}

	// 计算掉落物品
	item := make([]PassItem, 0)
	for i := 0; i < len(bag); i++ {
		droplist, _ := self.FindDrop(bag[i].ItemID, bag[i].Num)

		for j := 0; j < len(droplist); j++ {
			item = append(item, droplist[j])
		}
	}
	additemType, additemNum := lst[0].AddItemType, lst[0].AddItemNum
	if additemType != 0 && additemNum != 0 {
		item = append(item, PassItem{additemType, additemNum})
	}

	return item
}

func (self *ModFind) OnSave(sql bool) {
	self.Encode()
	self.Sql_Find.Update(sql)
}

func (self *ModFind) Decode() { // 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_Find.BaseFindInfo), &self.Sql_Find.baseFindInfo)
	json.Unmarshal([]byte(self.Sql_Find.RewardInfo), &self.Sql_Find.rewardInfo)
	json.Unmarshal([]byte(self.Sql_Find.WishInfo), &self.Sql_Find.wishInfo)
	json.Unmarshal([]byte(self.Sql_Find.Astrology), &self.Sql_Find.astrology)
	json.Unmarshal([]byte(self.Sql_Find.SelfSelection), &self.Sql_Find.selfSelection)
	json.Unmarshal([]byte(self.Sql_Find.GemFindGift), &self.Sql_Find.gemfindgift)
}

func (self *ModFind) Encode() { // 将data数据写入数据库
	self.Sql_Find.BaseFindInfo = HF_JtoA(self.Sql_Find.baseFindInfo)
	self.Sql_Find.RewardInfo = HF_JtoA(self.Sql_Find.rewardInfo)
	self.Sql_Find.WishInfo = HF_JtoA(self.Sql_Find.wishInfo)
	self.Sql_Find.Astrology = HF_JtoA(self.Sql_Find.astrology)
	self.Sql_Find.SelfSelection = HF_JtoA(self.Sql_Find.selfSelection)
	self.Sql_Find.GemFindGift = HF_JtoA(self.Sql_Find.gemfindgift)
}

func (self *ModFind) SendInfo() {
	self.CheckPool()
	var msg S2C_FindInfo
	msg.Cid = "findinfo"
	msg.BaseFindInfo = self.Sql_Find.baseFindInfo
	msg.RewardInfo = self.Sql_Find.rewardInfo
	msg.FindWishInfo = self.Sql_Find.wishInfo
	msg.Astrology = self.Sql_Find.astrology
	msg.SelfSelection = self.Sql_Find.selfSelection
	msg.LuckyPoolConfig = self.GetLuckyPoolConfig()
	msg.LuckyFindRecord = GetOfflineInfoMgr().GetLuckyFindRecord()
	msg.FindGiftProcess = self.Sql_Find.gemfindgift.Process		// 累计招募进度
	msg.IsGotFindGift = self.Sql_Find.gemfindgift.RewardList	// 宝箱领取情况
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModFind) OnRefresh() {
	for _, v := range self.Sql_Find.baseFindInfo {
		v.FindTimesToday = 0
		v.CalTime()
	}
	now := TimeServer().Unix()
	if self.Sql_Find.selfSelection == nil {
		self.Sql_Find.selfSelection = new(SelfSelection)
		self.Sql_Find.selfSelection.FindTimesCount = make(map[int]int)
	}

	if self.Sql_Find.selfSelection.EndTime < now {
		//获得特权并计算
		PrivilegeValue := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValues()
		self.CalSelfFind(PrivilegeValue, false)
	}
	self.SendInfo()
}

//掉落组抽取
func (self *ModFind) FindDrop(groupid, num int) ([]PassItem, []*GeneralRecord) {
	records := make([]*GeneralRecord, 0)

	//修正存在英雄的概率问题，之前的概率是读取的时候算好的，在这个地方并不准确
	correct := 0
	nodetemp := GetCsvMgr().GetDropgroup(groupid)
	for i := 0; i < len(nodetemp); i++ {
		if nodetemp[i].WishChance == 0 || nodetemp[i].Chance == nodetemp[i].WishChance {
			correct += nodetemp[i].Chance * 10
			continue
		}
		//检测心愿单
		if self.IsWish(nodetemp[i].Itemid) {
			correct += nodetemp[i].WishChance * 10
		} else {
			correct += nodetemp[i].Chance * 10
		}
	}

	droplist := make([]PassItem, 0)
	for i := 0; i < num; i++ {
		node := GetCsvMgr().GetDropgroup(groupid)
		//temp:=GetCsvMgr().PubchestdropgroupSum
		//temp[888]=0
		//if GetCsvMgr().PubchestdropgroupSum[groupid] + correct==0{
		//	fmt.Printf("%d,correct:%d",GetCsvMgr().PubchestdropgroupSum[groupid],correct)
		//}
		pro := self.GetRandInt(correct)

		//print(correct, "||", pro, "\n")

		cur := 0
		for i := 0; i < len(node); i++ {
			if nodetemp[i].WishChance == 0 || nodetemp[i].Chance == nodetemp[i].WishChance {
				cur += node[i].Chance * 10
			} else {
				if self.IsWish(nodetemp[i].Itemid) {
					cur += nodetemp[i].WishChance * 10
				} else {
					cur += nodetemp[i].Chance * 10
				}
			}
			if pro < cur {
				itemid := node[i].Itemid
				itemnum := node[i].Dropcardsnum
				if node[i].TimeGeneralsNotice > 0 {
					data := new(GeneralRecord)
					data.Item = itemid
					data.Num = itemnum
					data.RecordType = node[i].TimeGeneralsNotice
					records = append(records, data)
				}
				droplist = append(droplist, PassItem{itemid, itemnum})
				break
			}
		}
	}

	return droplist, records
}

// 触发条件 - 替换必掉
func (self *ModFind) FindSpecial(findtype, droptimes int) ([]DropGroupModify, int) {
	groups := make([]DropGroupModify, 0)
	SpType := 0
	lst := GetCsvMgr().PubChestSpecialLst[findtype]
	for i := 0; i < len(lst); i++ {
		if droptimes >= lst[i].Droptimemin && droptimes <= lst[i].Droptimemax {
			SpType = lst[i].SpType
			if lst[i].SpType == 1 {
				specail := lst[i].DropGroupModify
				items := strings.Split(specail, "|")
				if len(items) >= 3 {
					for j := 0; j < len(items); j++ {
						item := strings.Split(items[j], ":")
						var group DropGroupModify
						group.Original = HF_Atoi(item[0])
						group.New = HF_Atoi(item[1])
						group.Rate = HF_Atoi(item[2])

						groups = append(groups, group)
					}
				}
			} else {
				specail := lst[i].DropGroupModify
				items := strings.Split(specail, "|")
				for j := 0; j < len(items); j++ {
					item := strings.Split(items[j], ":")
					var group DropGroupModify
					group.Original = HF_Atoi(item[0])
					group.New = HF_Atoi(item[1])
					groups = append(groups, group)
				}
			}
		}
	}

	return groups, SpType
}

func (self *ModFind) IsWish(itemIt int) bool {

	heroId := (itemIt - 11000000) / 100
	for i := 0; i < len(self.Sql_Find.wishInfo); i++ {
		for j := 0; j < len(self.Sql_Find.wishInfo[i].WishList); j++ {
			if self.Sql_Find.wishInfo[i].WishList[j] == 0 {
				return false
			}
		}
	}

	for i := 0; i < len(self.Sql_Find.wishInfo); i++ {
		for j := 0; j < len(self.Sql_Find.wishInfo[i].WishList); j++ {
			if self.Sql_Find.wishInfo[i].WishList[j] == heroId {
				return true
			}
		}
	}

	return false
}

func (self *ModFind) CalSelfFind(value map[int]int, isSend bool) {
	if value == nil {
		return
	}

	if self.Sql_Find.selfSelection == nil {
		self.Sql_Find.selfSelection = new(SelfSelection)
		self.Sql_Find.selfSelection.FindTimesCount = make(map[int]int)
	}

	//总开关
	if value[PRIVILEGE_ADD_SELFFIND_COUNT] > 0 {
		self.Sql_Find.selfSelection.MaxTimes = value[PRIVILEGE_ADD_SELFFIND_COUNT]
		//激活的情况下增加VIP特权效果
		vipcsv := GetCsvMgr().VipConfigMap[self.player.Sql_UserBase.Vip]
		if vipcsv != nil {
			self.Sql_Find.selfSelection.MaxTimes += vipcsv.CallOptional
		}

		initTimes := GetCsvMgr().getInitNum(SIMPLE_NUM_SELFFIND_OFFSET)
		if initTimes > 0 {
			cd := int64(initTimes - value[PRIVILEGE_REDUCE_SELFFIND_OFFSET])
			now := TimeServer().Unix()
			if self.Sql_Find.selfSelection.StartTime == 0 {
				self.Sql_Find.selfSelection.StartTime = now
				self.Sql_Find.selfSelection.GetTimes = 0
			}
			self.Sql_Find.selfSelection.EndTime = self.Sql_Find.selfSelection.StartTime + cd
			if now > self.Sql_Find.selfSelection.EndTime {
				count := (now - self.Sql_Find.selfSelection.StartTime) / cd
				self.Sql_Find.selfSelection.StartTime = self.Sql_Find.selfSelection.StartTime + count*cd
				self.Sql_Find.selfSelection.EndTime = self.Sql_Find.selfSelection.StartTime + cd
				self.Sql_Find.selfSelection.GetTimes = 0
			}
		}
	}
	if isSend {
		self.SendInfo()
	}
}

func (self *ModFind) FindBeautyPool(msg *C2S_FindPool) {

	realType := 0
	findPool := self.GetFindPool(msg.Findtype)
	if findPool == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TYPE_ERROR"))
		return
	}
	dec := ""
	//条件检查
	if msg.FindNum == 1 { //单抽
		realType = REALTYPE_BEAUTY
		dec = "单次圣物招募"
	} else if msg.FindNum == 5 { //
		realType = REALTYPE_BEAUTY_FIVE
		dec = "五连圣物招募"
	} else {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_NUM_ERROR"))
		return
	}

	// 找到符合条件的掉落组
	lst := GetCsvMgr().PubchesttotalGroup[realType]
	if len(lst) == 0 {
		LogDebug("找不到抽卡类型：", realType)
		return
	}

	//检查消耗够不够
	payitem, payitemnum := lst[0].Payitem, lst[0].Payitemnum
	if payitem != 0 {
		if msg.FindNum == 1 && findPool.FreeNextTime < TimeServer().Unix() {

		} else {
			if err := self.player.HasObjectOkEasy(payitem, payitemnum); err != nil {
				self.player.SendErrInfo("err", err.Error())
				return
			}
		}
	}

	// 计算必掉掉落
	bag := make([]PassItem, 0)
	certaintimes := lst[0].Certaintimesdroptype
	if certaintimes > 0 {
		itemid := 0
		if len(lst[0].Dropgroups) != 4 {
			self.player.SendErr("len(lst[0].Dropgroups) != 4")
			return
		}
		certainitem := lst[0].Dropgroups[3]
		tempCer := findPool.FindTimes % certaintimes
		tempCer += msg.FindNum
		if tempCer/certaintimes > 0 {
			itemid = certainitem
		}

		if itemid != 0 {
			var dropitem PassItem
			dropitem.ItemID = itemid
			dropitem.Num = 1
			bag = append(bag, dropitem)
		}
	}
	if findPool.FindTimesCount == nil {
		findPool.FindTimesCount = make(map[int]int, 0)
	}
	findPool.FindTimes += msg.FindNum
	findPool.FindTimesCount[realType]++
	times := findPool.getTimesAll(realType)
	item := self.LootItems(lst, realType, findPool.FindTimes, certaintimes, bag, times)
	//扣除消耗
	costItems := make([]PassItem, 0)
	if msg.FindNum == 1 && findPool.FreeNextTime < TimeServer().Unix() {
		findPool.FreeNextTime = HF_GetNextDayStart()
	} else {
		costItems = self.player.RemoveObjectSimple(payitem, payitemnum, dec, realType, 0, 0)
	}
	//计算任务进度     20200302友情不计算
	//self.CalRewardTimes(getItem, findPool, msg.FindNum)
	//self.player.HandleTask(FindTask, 2, 1, 0)
	findPool.FindTimesToday += msg.FindNum
	//self.player.HandleTask(TASK_TYPE_SUMMON_HEROS, msg.FindNum, 2, 0)
	param2 := findPool.FindTimesCount[realType]*msg.FindNum*1000 + realType
	//self.player.HandleTask(TASK_TYPE_SUMMON_ELITE_HEROS, nCount, realType, 0)
	//self.player.HandleTask(TASK_TYPE_DECOMPOSE_HEROS, nDecompose, 0, 0)
	getItems := self.player.AddObjectPassItem(item, dec, realType, param2, 0, )
	//getItemsTran := self.player.AddObjectItemMap(getItemTran, dec, 0, 0, 0)
	var msgRel S2C_FindPool
	msgRel.Cid = "drawok"
	msgRel.FindType = msg.Findtype
	msgRel.FindNum = msg.FindNum
	msgRel.FindNumToday = findPool.FindTimesToday
	msgRel.FreeNextTime = findPool.FreeNextTime
	msgRel.Item = item
	msgRel.CostItems = costItems
	msgRel.GetItems = getItems
	msgRel.RewardInfo = self.Sql_Find.rewardInfo
	msgRel.FindTimes = findPool.FindTimes
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	if msg.FindNum == 1 {
		//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_FRIEND_ONE, 1, 0, 0, dec, 0, 0, self.player)
	} else {
		//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_FRIEND_TEN, msg.FindNum, 0, 0, dec, 0, 0, self.player)
	}
}

func (self *ModFind) FindLuckyPool(msg *C2S_FindPool) {

	isOpen, index := GetActivityMgr().JudgeOpenAllIndex(ACT_LUCKY_FIND, ACT_LUCKY_FIND)
	if !isOpen {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ACT_NOT_OPEN"))
		return
	}

	findPool := self.GetFindPool(msg.Findtype)
	if findPool == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TYPE_ERROR"))
		return
	}
	dec := ""
	costItem := ITEM_GEM
	costNum := 10000
	//条件检查
	if msg.FindNum == 1 { //单抽
		dec = "单次福袋召唤"
	} else if msg.FindNum == 10 { //
		dec = "十连福袋召唤"
		costNum = 88888
	} else {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_NUM_ERROR"))
		return
	}
	if msg.FindNum == 1 && findPool.FreeNextTime < TimeServer().Unix() {

	} else {
		if err := self.player.HasObjectOkEasy(costItem, costNum); err != nil {
			self.player.SendErrInfo("err", err.Error())
			return
		}
	}
	lottery := GetActivityMgr().getActN4(index)
	if findPool.FindTimesCount == nil {
		findPool.FindTimesCount = make(map[int]int, 0)
	}
	findPool.FindTimes += msg.FindNum
	findPool.FindTimesCount[lottery]++
	//扣除消耗
	costItems := make([]PassItem, 0)
	if msg.FindNum == 1 && findPool.FreeNextTime < TimeServer().Unix() {
		findPool.FreeNextTime = HF_GetNextDayStart()
	} else {
		costItems = self.player.RemoveObjectSimple(costItem, costNum, dec, ACT_LUCKY_FIND, lottery, 0)
	}
	item := make([]PassItem, 0)
	for i := 0; i < msg.FindNum; i++ {
		lotteryConfig := self.GetLuckyPoolItem(lottery)
		if lotteryConfig == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIND_TYPE_ERROR"))
			return
		}
		item = append(item, PassItem{ItemID: lotteryConfig.Itemid, Num: lotteryConfig.Min})
		//增加广播显示
		if lotteryConfig.Havelot > 0 {
			record := new(LotteryDrawRecord)
			record.Uid = self.player.GetUid()
			record.Name = self.player.GetName()
			record.Times = lotteryConfig.Havelot
			record.ItemId = lotteryConfig.Itemid
			record.Num = lotteryConfig.Min
			GetOfflineInfoMgr().AddLuckyFindRecord(record)
		}
	}
	findPool.FindTimesToday += msg.FindNum
	param2 := findPool.FindTimesCount[lottery]*msg.FindNum*1000 + lottery
	getItems := self.player.AddObjectPassItem(item, dec, lottery, param2, 0, )

	self.player.HandleTask(TASK_TYPE_LUCKY_FIND, msg.FindNum, 0, 0)
	var msgRel S2C_FindPool
	msgRel.Cid = "findpool"
	msgRel.FindType = msg.Findtype
	msgRel.FindNum = msg.FindNum
	msgRel.FindNumToday = findPool.FindTimesToday
	msgRel.FreeNextTime = findPool.FreeNextTime
	msgRel.Item = item
	msgRel.CostItems = costItems
	msgRel.GetItems = getItems
	msgRel.RewardInfo = self.Sql_Find.rewardInfo
	msgRel.FindTimes = findPool.FindTimes
	msgRel.LuckyFindRecord = GetOfflineInfoMgr().GetLuckyFindRecord()
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}

func (self *ModFind) GetLuckyPoolConfig() []LuckyPassItem {

	rel := make([]LuckyPassItem, 0)
	isOpen, index := GetActivityMgr().JudgeOpenAllIndex(ACT_LUCKY_FIND, ACT_LUCKY_FIND)
	if !isOpen {
		return rel
	}

	lottery := GetActivityMgr().getActN4(index)

	group, ok := GetLootMgr().LootGroupMap[lottery]
	if !ok {
		return rel
	}
	for _, v := range group.Configs {
		if v.Itemwt == 0 {
			continue
		}
		rel = append(rel, LuckyPassItem{ItemID: v.Itemid, Num: v.Min, Type: v.Havelot})
	}
	return rel
}

func (self *ModFind) GetLuckyPoolItem(lootId int) *LotteryConfig {

	lootGrop, ok := GetLootMgr().LootGroupMap[lootId]
	if !ok {
		LogError("掉落Id配置不存在, lootId=", lootId)
		return nil
	}

	// 先判断掉不掉这个掉落包(0~4999)
	randNum := HF_GetRandom(RatioNum)
	// 什么都不掉落
	if lootGrop.Chance != RatioNum && randNum >= lootGrop.Chance {
		return nil
	}

	// 计算组, 权重
	targetGroupId := lootGrop.LootGroupId()
	if targetGroupId == 0 {
		return nil
	}

	lootItems := lootGrop.GetItems(targetGroupId)
	total := 0
	for _, v := range lootItems {
		total += v.Rate
	}
	if total <= 0 {
		LogError("随机0值，掉落组:", targetGroupId)
		return nil
	}
	rand := HF_GetRandom(total)

	check := 0
	id := 0
	for _, v := range lootItems {
		check += v.Rate
		if rand < check {
			id = v.Id
			break
		}
	}

	if id == 0 {
		return nil
	}

	lootConfig, ok := GetLootMgr().LotteryMap[id]
	if !ok {
		LogError("掉落不存在, Id:", id)
		return nil
	}
	return lootConfig
}
