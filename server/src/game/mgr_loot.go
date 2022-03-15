package game

// 新随机掉落
const (
	LootByRate  = 1     // 权重
	LootByRatio = 2     // 比例
	RatioNum    = 10000 // 万分比
	PercentNum  = 100   // 百分比
)

// 掉落管理器
type LootMgr struct {
	LotteryConfig []*LotteryConfig
	LotteryMap    map[int]*LotteryConfig
	LootGroupMap  map[int]*LootGroup
}

var lootMgr *LootMgr

func GetLootMgr() *LootMgr {
	if lootMgr == nil {
		lootMgr = new(LootMgr)
	}
	return lootMgr
}

// 掉落配置
type LotteryConfig struct {
	Id          int    `json:"id"`
	Lotteryid   int    `json:"lotteryid"`
	Description string `json:"description"`
	Chance      int    `json:"chance"`
	Type        int    `json:"type"`
	Groupid     int    `json:"groupid"`
	Groupwt     int    `json:"groupwt"`
	Itemid      int    `json:"itemid"`
	Itemdes     string `json:"itemdes"`
	Itemwt      int    `json:"itemwt"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Havelot     int    `json:"havelot"`
	Haveitemwt  int    `json:"haveitemwt"`
}

// 掉落组
type LootGroup struct {
	Lotteryid int
	LootType  int
	Chance    int
	Configs   []*LotteryConfig
}

type GroupRate struct {
	GroupId int
	Rate    int
}

type LootItem struct {
	Id         int
	ItemId     int
	Rate       int
	Min        int
	Max        int
	Havelot    int
	Haveitemwt int
}

// 加载配置
func (self *LootMgr) LoadLottery() {
	GetCsvUtilMgr().LoadCsv("Lottery_New", &self.LotteryConfig)
	self.LotteryMap = make(map[int]*LotteryConfig)
	GetCsvUtilMgr().LoadCsv("Lottery_New", &self.LotteryMap)
	self.LootGroupMap = make(map[int]*LootGroup)
	for _, v := range self.LotteryConfig {
		value, ok := self.LootGroupMap[v.Lotteryid]
		if !ok {
			info := &LootGroup{}
			info.Lotteryid = v.Lotteryid
			info.LootType = v.Type
			info.Chance = v.Chance
			info.Configs = append(info.Configs, v)
			self.LootGroupMap[v.Lotteryid] = info
		} else {
			if value.LootType != v.Type {
				LogError("掉落类型不一致, LotteryId:", v.Lotteryid)
			}
			value.Configs = append(value.Configs, v)
		}
	}
}

func (self *LootMgr) LootItemLst(lootId int, player *Player) []PassItem {
	var info = self.LootItem(lootId, player)
	var lst []PassItem
	for _, v := range info {
		lst = append(lst, PassItem{v.ItemId, v.ItemNum})
	}
	return lst
}

// 掉落道具
func (self *LootMgr) LootItem(lootId int, player *Player) map[int]*Item {
	var res = make(map[int]*Item)
	lootGrop, ok := self.LootGroupMap[lootId]
	if !ok {
		LogError("掉落Id配置不存在, lootId=", lootId)
		return res
	}

	// 先判断掉不掉这个掉落包(0~4999)
	randNum := HF_GetRandom(RatioNum)
	// 什么都不掉落
	if lootGrop.Chance != RatioNum && randNum >= lootGrop.Chance {
		return res
	}

	// 再判断掉落类型
	if lootGrop.LootType == LootByRate {
		return lootGrop.LootByRate(player)
	} else if lootGrop.LootType == LootByRatio {
		return lootGrop.LootByPercent(player)
	} else {
		LogError("掉落类型错误")
		return res
	}
}

// 掉落组
func (self *LootMgr) LootItems(lootIds []int, player *Player) map[int]*Item {
	total := make(map[int]*Item)
	for _, lootId := range lootIds {
		if lootId == 0 {
			continue
		}
		result := self.LootItem(lootId, player)
		AddItemMap(total, result)
	}
	return total
}

// 获取掉落组
func (self *LootGroup) GetGroupRates() []*GroupRate {
	var groupRates []*GroupRate
	for _, v := range self.Configs {
		if v.Groupwt == 0 {
			continue
		}

		found := false
		for _, r := range groupRates {
			if r.GroupId == v.Groupid {
				//LogError("已经有相同的组了, id:", v.Lotteryid)
				found = true
				break
			}
		}

		if found {
			continue
		}
		groupRates = append(groupRates, &GroupRate{v.Groupid, v.Groupwt})
	}
	return groupRates
}

// 随机一个组
func (self *LootGroup) LootGroupId() int {
	var groupRates = self.GetGroupRates()
	// 按照权重随机出一个组
	total := 0
	for _, v := range groupRates {
		total += v.Rate
	}

	if total == 0 {
		LogError("权重之和为0")
		return 0
	}
	rand := HF_GetRandom(total)

	check := 0
	targetGroupId := 0
	for _, v := range groupRates {
		check += v.Rate
		if rand < check {
			targetGroupId = v.GroupId
			break
		}
	}

	return targetGroupId
}

// 获取这个组里面道具
func (self *LootGroup) GetItems(groupId int) []*LootItem {
	var lootItems []*LootItem
	for _, v := range self.Configs {
		if v.Groupid != groupId {
			continue
		}
		for _, r := range lootItems {
			if v.Id == r.Id {
				LogError("已经有相同的Id了, id=", v.Id)
				return lootItems
			}
		}
		lootItems = append(lootItems, &LootItem{v.Id, v.Itemid, v.Itemwt, v.Min, v.Max, v.Havelot, v.Haveitemwt})
	}
	return lootItems
}

// 根据权重掉落物品
func (self *LootGroup) LootItems(groupId int, player *Player) map[int]*Item {
	var res = make(map[int]*Item)
	lootItems := self.GetItems(groupId)
	total := 0
	for _, v := range lootItems {
		if v.Havelot == 0 || player == nil {
			total += v.Rate
		} else if player != nil {
			hero := player.getHero(v.Havelot)
			if hero != nil {
				total += v.Haveitemwt
			}
		}
	}
	if total <= 0 {
		LogError("随机0值，掉落组:", groupId)
		return res
	}
	rand := HF_GetRandom(total)

	check := 0
	id := 0
	for _, v := range lootItems {
		if v.Havelot == 0 || player == nil {
			check += v.Rate
		} else if player != nil {
			hero := player.getHero(v.Havelot)
			if hero != nil {
				check += v.Haveitemwt
			}
		}
		if rand < check {
			id = v.Id
			break
		}
	}

	if id == 0 {
		return res
	}

	lootConfig, ok := GetLootMgr().LotteryMap[id]
	if !ok {
		LogError("掉落不存在, Id:", id)
		return res
	}

	lootNum := HF_RandInt(lootConfig.Min, lootConfig.Max+1)
	res[lootConfig.Itemid] = &Item{lootConfig.Itemid, lootNum}

	return res
}

// 按照权重掉落
func (self *LootGroup) LootByRate(player *Player) map[int]*Item {
	var res = make(map[int]*Item)
	// 计算组, 权重
	targetGroupId := self.LootGroupId()
	if targetGroupId == 0 {
		return res
	}

	// 再从组里面随机几个道具
	lootItems := self.LootItems(targetGroupId, player)
	return lootItems
}

// 万分比随机
func (self *LootGroup) LootGroupId2() []int {
	var groupRates = self.GetGroupRates()
	var groupIds []int
	for _, v := range groupRates {
		rand := HF_GetRandom(RatioNum)
		if rand < v.Rate {
			groupIds = append(groupIds, v.GroupId)
		}
	}

	return groupIds
}

// 按照万分比掉落
func (self *LootGroup) LootByPercent(player *Player) map[int]*Item {
	var res = make(map[int]*Item)
	groupIds := self.LootGroupId2()
	for _, groupId := range groupIds {
		lootItems := self.LootItems(groupId, player)
		AddItemMap(res, lootItems)
	}
	return res
}