package game

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

import (
	"encoding/csv"
	"github.com/sanity-io/litter"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

var CAMP_NAME []string
var PART_NAME []string

func (self *TimeGeneralRank) GetAward(rank int, point int) []PassItem {
	var res []PassItem
	res = append(res, self.NormalAward...)
	if self.NeetPoint > 0 && point >= self.NeetPoint {
		res = append(res, self.ExtraAward...)
	}

	return res
}

func (self *TimeGeneralRank) IsHasExt(rank int, point int) int {
	if self.NeetPoint > 0 && point >= self.NeetPoint {
		return LOGIC_TRUE
	}
	return LOGIC_FALSE
}

var csvmgrsingleton *CsvMgr = nil

//! public
func GetCsvMgr() *CsvMgr {
	if csvmgrsingleton == nil {
		csvmgrsingleton = new(CsvMgr)
		csvmgrsingleton.InitStruct()
	}

	return csvmgrsingleton
}

func (self *CsvMgr) InitStruct() {
	self.Data = make(map[string]map[int]CsvNode)
	self.Data2 = make(map[string][]CsvNode)
	self.Equchest_CSV = make(map[int][]CsvNode)
	self.Equchest_SUM = make(map[int]int)
	self.Soulchest_CSV = make(map[int][]CsvNode)
	self.Soulchest_SUM = make(map[int]int)
	self.WarTaskTask_CSV = make(map[int]*TaskNode)
	self.SevenDayTask_CSV = make(map[int]*TaskNode)
	self.SevenStatus = make(map[int]int)
	self.HalfMonnStatus = make(map[int]int)
	self.HalfMoonTask_CSV = make(map[int]*TaskNode)
	self.WarTargetTask_CSV = make(map[int]*TaskNode)
	self.WarTarget_CSV = make(map[int]CsvNode)
	self.WarList_CSV = make(map[int]CsvNode)
	self.War_CSV = make(map[int]CsvNode)
	self.Exciting_CSV = make(map[int]CsvNode)
	self.Gestapo_CSV = make(map[int]CsvNode)
	self.Conditions_CSV = make(map[int]CsvNode)
	self.Diplomacy_CSV = make(map[int]CsvNode)
	self.Honorkill_CSV = make(map[int]CsvNode)
	self.Homeoffice_ranking_CSV = make(map[int]CsvNode)
	self.Reward_CSV = make(map[int]CsvNode)
	self.Gemsweeper_extradrop_CSV = make(map[int]CsvNode)
	self.Gemsweeper_itembag_GG_CSV = make(map[int][]CsvNode)
	self.Gemsweeper_itembag_SUM = make(map[int]int)
	self.Peoplecity_CSV = make(map[int][]CsvNode)
	self.Peoplebox_CSV = make(map[int][]CsvNode)
	self.Harem_TreasureaAdv_CSV = make(map[int]CsvNode)
	self.PromoteBox_CSV = make(map[int]CsvNode)
	self.Gemsweeper_event_CSV = make(map[int]CsvNode)
	self.Money_CSV = make(map[int]CsvNode)
	self.SevenDay_CSV = make(map[int]CsvNode)
	self.HalfMoon_CSV = make(map[int]CsvNode)

	self.Activity_Type = make([]int, 0)
	self.State_City_Num = make(map[int][]int, 0)
	//self.ActivityBox_CSV = make(map[int][]CsvNode, 0)

	self.HorseParam_CSV = make(map[string]CsvNode)
	self.HorseAttr_CSV = make(map[int][]CsvNode)

	self.TimeGeneralRankLst = make([]*TimeGeneralRank, 0)
	self.GeneralRankMail = &GeneralRankMail{}
	self.HalfMoonLogin = make(map[int]int)
	self.HalfMoonTrial = make(map[int]*HalfMoonTrial)
	self.HalfMoonBeauty = make(map[int]*HalfMoonBeauty)

	self.SmeltPurchaseConfig = make(map[int]*SmeltPurchaseConfig)
	self.TigerStuntConfig = make([]*TigerStuntConfig, 0)
	self.TigerStuntUpgradeConfig = make([]*TigerStuntUpgradeConfig, 0)
	self.TigerSymbolConfig = make([]*TigerSymbolConfig, 0)
	self.TigerUpgradeConfig = make([]*TigerUpgradeConfig, 0)
	self.TigerAdvancedConfig = make([]*TigerAdvancedConfig, 0)
	self.TigerAttributeConfig = make([]*TigerAttributeConfig, 0)

	self.ConsumerList_CSV = make(map[int][]CsvNode, 0)
	self.ConsumerShop_CSV = make(map[int][]CsvNode, 0)

	self.CitygvgConfig = make([]*CitygvgConfig, 0)
	self.CityrandomConfig = make([]*CityrandomConfig, 0)
	self.CityrepressConfig = make(map[int]*CityrepressConfig)
	self.TreasureEquipConfig = make(map[int]*TreasureEquipConfig)
	self.TreasureHeroConfig = make(map[int]*TreasureHeroConfig)
	self.TreasureSuitAttributeConfig = make([]*TreasureSuitAttributeConfig, 0)
	self.TreasureAwakenConfig = make([]*TreasureAwakenConfig, 0)
	self.TreasureClearAttributeConfig = make([]*TreasureClearAttributeConfig, 0)
	self.MaxTreasureHoleLv = make(map[int]int)
	self.TreasureClearItemConfig = make(map[int]*TreasureClearItemConfig)
	self.MaxTreasureAwaken = make(map[int]int)
	self.LuckyturntablelistConfig = make([]*LuckyturntablelistConfig, 0)
	self.TreasureDecomposeConfig = make([]*TreasureDecomposeConfig, 0)
	self.TreasureDecomposeMap = make(map[int]map[int]*TreasureDecomposeConfig)
	self.LuckegggroupConfig = make([]*LuckegggroupConfig, 0)
	self.LuckeggConfig = make([]*LuckeggConfig, 0)
	self.LuckstartConfig = make([]*LuckstartConfig, 0)
	self.LuckStartMap = make(map[int]map[int]*TaskNode)
	self.LuckStartConfigMap = make(map[int]map[int]*LuckstartConfig)
	self.GemsweeperConfig = make(map[int]*GemsweeperConfig)
	self.GemsweepereventConfig = make([]*GemsweepereventConfig, 0)
	self.GemGroupCycle = make(map[int]int)
	self.GemGroupStep = make(map[int]int)
	self.WorldLevel = make(map[int]*WorldLevelConfig)
	self.WorldMap = make(map[int]*WorldMapConfig)

	self.ActivityTimeGiftMap = make(map[int]*TimeGiftConfig)
	self.ActivityTimeGiftGroup = make(map[int][]*TimeGiftConfig)

	self.DailyrechargeConfig = make([]*DailyrechargeConfig, 0)

	self.EntanglementConfig = make(map[int]*EntanglementConfig)
	self.EntanglementMapConfig = make(map[int]*EntanglementFate)

	self.RewardForbarMapConfig = make(map[int][]*RewardForbarConfig)
	self.RewardForbarColorMapConfig = make(map[int]map[int][]*RewardForbarConfig)

	self.RankTaskMapConfig = make(map[int][]*RankTaskConfig)
}

type MiliWeekTaskConfig struct {
	Id       int   `json:"id"`
	TaskType int   `json:"tasktypes"`
	Conds    []int `json:"n"`
	ItemIds  []int `json:"item"`
	Nums     []int `json:"num"`
}

// 已经替换的绝大多数,后面涉及到维护的模块用LoadOther里面的配置替换
// 方便后续维护
func (self *CsvMgr) InitData() {
	self.ReadData("1GM")
	self.ReadData("Level_Firstitem")
	self.ReadData("Level_Item")
	self.ReadData("title")     // 暂时不换,有点恶心,n12,n22,n23
	self.ReadData("Level_Map") // 关卡
	self.ReadData2("City_Way") // 城防
	self.ReadData2("equipchest")
	self.ReadData2("soulchest")
	self.ReadData("warcity")  // 国战
	self.ReadData("warplay")  // 国战
	self.ReadData("War_Buff") // 国战
	self.ReadData("warencourage")

	self.ReadData2("Expedition_Buff")      // 远征
	self.ReadData2("Expedition_Buffgroup") // 远征
	self.ReadData2("wartask")              // 军功任务
	self.ReadData2("trial")
	self.ReadData("worldlvtpye") // 战斗
	self.ReadData2("exciting")
	self.ReadData2("World_Event")

	self.ReadData2("Diplomacy_changes")
	self.ReadData2("honorkill")
	self.ReadData("Gemsweeper")
	self.ReadData("Gemsweeper_award")
	self.ReadData("visit")
	self.ReadData("industry")
	self.ReadData2("Homeoffice_Purchase")
	self.ReadData2("Homeoffice_Ranking")
	self.ReadData2("Homeoffice_Cityaward")
	self.ReadData("Guild_ActiveCopy")
	self.ReadData2("camptask")
	self.ReadData("camptask")
	self.ReadData2("campreward")
	self.ReadData("Activitynew1")
	self.ReadData("Activitynew")
	self.ReadData2("Activitynew1")
	self.ReadData2("Activitynew")
	//self.ReadData2("Activitybox")
	self.ReadData("Money")
	self.ReadData("Gemsweeper_rank")
	self.ReadData2("Gemsweeper_extradrop")
	self.ReadData("Nationalwar_Parm")
	self.ReadData("Nationalwar_Enrollment")
	self.ReadData2("Gemsweeper_itembag")
	self.ReadData("peoplecity")
	self.ReadData("spyreward")
	self.ReadData("spytreasure")
	self.ReadData2("visitchance")
	self.ReadData("promotebox")
	self.ReadData2("Gemsweeper_event")
	self.ReadData2("Money")

	self.ReadData("Fund")
	self.ReadData("Fund_Buy")
	self.ReadData2("Sevenday")
	self.ReadData2("Halfmoon")
	self.ReadData2("wartarget")
	self.ReadData("warcontribution")
	self.ReadData("warlist")

	//! 坐骑
	self.ReadData("Horse_BattleSteed")
	self.ReadData2("Horse_BattleSteed_Attribute")
	self.ReadData("Horse_BattleSteed_Awaken")
	self.ReadData2("Horse_Judge_call")
	self.ReadData("Horse_Judge_call")
	self.ReadData("Horse_Judge_Discern")
	self.ReadData("Horse_Judge_level")
	self.ReadData("Horse_Parm")
	self.ReadData("Horse_Shop")
	self.ReadData2("Horse_Parm")

	//! 消费排行榜
	self.ReadData("Consumetop_Shop")
	self.ReadData("Consumetop_Attack")
	self.ReadData("Consumetop_List")
	self.ReadData("Consumetop_Hp")
	self.ReadData2("Consumetop_Luck")
	self.ReadData("Consumetop_Boss")

	//! 预约神将
	self.ReadData2("bespeak")

	// 虎符Start
	self.ReadData("Tiger_Upgrade")
	self.ReadData("Tiger_Advanced")
	self.ReadData("Tiger_Attribute")
	self.ReadData("Tiger_Stunt")
	self.ReadData2("Tiger_StuntUpgrade")
	self.ReadData("Tiger_Symbol")
	// 虎符End

	// 珍宝Start
	self.ReadData2("Treasure_Awaken")
	self.ReadData2("Treasure_ClearAttribute")
	self.ReadData2("Treasure_ClearItem")
	self.ReadData2("Treasure_Equip")
	self.ReadData("Treasure_Hero")
	self.ReadData2("Treasure_SuitAttribute")
	self.ReadData2("Treasure_Suit")
	self.ReadData2("Treasure_Decompose")
	// 珍宝End

	self.ReadData2("citygvg")
	self.ReadData2("cityrandom")
	self.ReadData("cityrepress")

	//! 新的城池配置表
	self.ReadData("World_Map")
	self.ReadData("World_Power")
	self.ReadData("World_Lvtpye")

	//! 王权争夺
	self.ReadData("Crown_Fight")

	self.ReadEquchest()
	self.ReadSoulchest()
	self.ReadWarTask()

	self.ReadExciting()
	self.ReadDiplomacy()
	self.ReadHonorkill()
	self.ReadHomeoffice_ranking()
	self.ReadCampreward()
	self.ReadGemsweeper_extradrop()
	self.ReadGemsweeper_itembag()
	self.ReadPromotebox()
	self.ReadGemsweeper_event()
	self.ReadMoney()
	self.ReadSevenDay()
	self.ReadSevenDayTask()
	self.ReadHalfMoon()
	self.ReadHalfMoonTask()
	//self.ReadWarCity()
	self.ReadWarTargetTask()
	self.ReadActivityType()
	//self.ReadActivityBox()
	self.ReadStateCity()

	self.ReadHorseParam()
	self.ReadHorseAttr()

	self.ReadTigerAttributeConfig()
	self.ReadTigerStuntConfig()
	self.ReadTigerStuntUpgradeConfig()
	self.ReadTigerSymbolConfig()
	self.ReadTigerUpgradeConfig()
	self.ReadTigerAdvancedConfig()

	//self.ReadCitygvgConfig()
	//self.ReadCityrandomConfig()
	//self.ReadCityrepressConfig()
	self.ReadTreasureEquipConfig()
	self.ReadTreasureHeroConfig()
	self.ReadTreasureSuitAttributeConfig()
	self.ReadTreasureAwakenConfig()
	self.ReadTreasureClearAttributeConfig()
	self.ReadTreasureClearItemConfig()
	self.ReadTreasureDecomposeConfig()
	self.ReadGemsweeperConfig()
	self.ReadGemsweepereventConfig()
	self.LoadCsv()
}

func (self *CsvMgr) Reload() {
	self.InitStruct()
	self.InitData()
	log.Println("reload csv")
}

func (self *CsvMgr) ReadData(name string) {
	_, ok := self.Data[name]
	if ok {
		//log.Println("重复读入csv:", name)
		LogError("重复读入csv:", name)
		return
	}

	file, err := os.Open("csv/" + name + ".csv")
	if err != nil {
		log.Fatalln("csv err1:", name, err)
		return
	}
	defer file.Close()

	header := make([]string, 0)
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln("csv err2:", name, err, strings.Join(record, ","))
			return
		}

		//log.Println(record)

		if len(header) == 0 {
			header = record
		} else {
			id, err := strconv.Atoi(record[0])
			if err != nil {
				log.Fatalln("csv err3:", name, err)
				return
			}

			_, ok := self.Data[name]
			if !ok {
				self.Data[name] = make(map[int]CsvNode)
			}

			_, ok = self.Data[name][id]
			if !ok {
				self.Data[name][id] = make(CsvNode)
			}

			for i := 0; i < len(record); i++ {
				self.Data[name][id][header[i]] = record[i]
			}
		}
	}
}

func (self *CsvMgr) ReadData2(name string) {
	_, ok := self.Data2[name]
	if ok {
		//log.Println("重复读入csv:", name)
		LogError("重复读入csv:", name)
		return
	}

	file, err := os.Open("csv/" + name + ".csv")
	if err != nil {
		log.Fatalln("csv err4:", name, err)
		return
	}
	defer file.Close()

	header := make([]string, 0)
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln("csv err5:", name, err)
			return
		}

		//log.Println(record)

		if len(header) == 0 {
			header = record
		} else {
			node := make(CsvNode)
			for i := 0; i < len(record); i++ {
				node[header[i]] = record[i]
			}
			self.Data2[name] = append(self.Data2[name], node)
		}
	}
}

// 这个函数要重构掉
func (self *CsvMgr) GetData2Int(table string, id int, field string) int {
	value, _ := strconv.Atoi(self.Data[table][id][field])
	return value
}

// 这个函数要重构掉
func (self *CsvMgr) GetData2String(table string, id int, field string) string {
	return self.Data[table][id][field]
}

func (self *CsvMgr) GetHolyLegend(beautyid int, chapter int) *HolyLegendConfig {
	for _, v := range self.HolyLegendConfig {
		if v.Beautyid == beautyid && v.Chapter == chapter {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetHolyLengendLevel(groupId int, index int) *HolyLegendLevelConfig {
	for _, v := range self.HolyLegendLevelConfig {
		if v.Group == groupId && v.Levelindex == index {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetExpedition(id int) (*ExpeditionConfig, bool) {
	value, ok := self.ExpeditionConfig[id]
	if !ok {
		return nil, false
	}
	return value, ok
}

func (self *CsvMgr) GetExpeditionBuff(id int) (*ExpeditionbuffConfig, bool) {
	value, ok := self.ExpeditionbuffConfig[id]
	return value, ok
}

func (self *CsvMgr) GetExpeditionBuffGroup(buffgroup int) (*ExpeditionbuffgroupConfig, bool) {
	value, ok := self.ExpeditionbuffgroupConfig[buffgroup]
	return value, ok
}

func (self *CsvMgr) ReadMoney() {
	data := self.Data2["Money"]
	for _, value := range data {
		self.Money_CSV[HF_Atoi(value["id"])] = value
	}
}

func (self *CsvMgr) GetMoney(id int) (CsvNode, bool) {
	value, ok := self.Money_CSV[id]
	return value, ok
}

func (self *CsvMgr) ReadPromotebox() {
	data := self.Data["promotebox"]
	for _, value := range data {
		self.PromoteBox_CSV[HF_Atoi(value["boxid"])] = value
	}
}

func (self *CsvMgr) GetHolyAdvance(beautyid, lv int) (*HolyUpgradeConfig, bool) {
	for _, v := range self.HolyUpgradeConfig {
		if v.Beautyid == beautyid && v.Stagelv == lv {
			return v, true
		}
	}
	return nil, false
}

func (self *CsvMgr) GetHarem_TreasureaAdv(treasurea_id, lv int) (CsvNode, bool) {
	value, ok := self.Harem_TreasureaAdv_CSV[treasurea_id*100+lv]
	return value, ok
}

func (self *CsvMgr) GetHolyPartsAdv(treasurea_id, lv int) (*HolyPartsUpgradeConfig, bool) {
	for _, v := range self.HolyPartsUpgradeConfig {
		if v.Treasureaid == treasurea_id && v.Stagelv == lv {
			return v, true
		}
	}
	return nil, false

}

func (self *CsvMgr) LoadShop() {
	self.ShopGrid = make(map[int]map[int][]*ShopConfig)
	self.ShopSumGrid = make(map[int]map[int]int)
	for _, v := range self.ShopConfig {
		_, ok := self.ShopGrid[v.Type]
		if !ok {
			self.ShopGrid[v.Type] = make(map[int][]*ShopConfig)
			self.ShopSumGrid[v.Type] = make(map[int]int)
		}

		self.ShopGrid[v.Type][v.Grid] = append(self.ShopGrid[v.Type][v.Grid], v)
		self.ShopSumGrid[v.Type][v.Grid] += v.Weightfunction
	}
}

func (self *CsvMgr) LoadPubChestTotal() {
	self.PubchesttotalGroup = make(map[int][]*PubchesttotalConfig)
	data := self.PubchesttotalConfig
	for _, value := range data {
		self.PubchesttotalGroup[value.Paypubtype] = append(self.PubchesttotalGroup[value.Paypubtype], value)
	}
	//litter.Dump(self.PubchesttotalGroup)
}

func (self *CsvMgr) LoadDropGroup() {
	self.PubchestdropgroupLst = make(map[int][]*PubchestdropgroupConfig)
	self.PubchestdropgroupSum = make(map[int]int)
	for _, value := range self.PubchestdropgroupConfig {
		self.PubchestdropgroupLst[value.Dropgroup] = append(self.PubchestdropgroupLst[value.Dropgroup], value)
		self.PubchestdropgroupSum[value.Dropgroup] += value.Chance
	}
}

func (self *CsvMgr) GetDropgroup(bag int) []*PubchestdropgroupConfig {
	return self.PubchestdropgroupLst[bag]
}

func (self *CsvMgr) LoadPubChestSpecial() {
	data := self.PubchestspecialConfig
	self.PubChestSpecialLst = make(map[int][]*PubchestspecialConfig)
	for _, value := range data {
		items := strings.Split(value.Paytype, "|")
		for i := 0; i < len(items); i++ {
			realType := HF_Atoi(items[i])
			self.PubChestSpecialLst[realType] = append(self.PubChestSpecialLst[realType], value)
		}
	}
}

func (self *CsvMgr) ReadEquchest() {
	data := self.Data2["equipchest"]
	for _, value := range data {
		self.Equchest_CSV[HF_Atoi(value["bag"])] = append(self.Equchest_CSV[HF_Atoi(value["bag"])], value)
		self.Equchest_SUM[HF_Atoi(value["bag"])] += HF_Atoi(value["probability2"])
	}
}

func (self *CsvMgr) GetEquchest(bag int) []CsvNode {
	return self.Equchest_CSV[bag]
}

func (self *CsvMgr) GetEquchestPro(bag int) int {
	node := self.Equchest_CSV[bag]
	for i := 0; i < len(node); i++ {
		if HF_Atoi(node[i]["probability1"]) != 0 {
			return HF_Atoi(node[i]["probability1"])
		}
	}

	return 0
}

func (self *CsvMgr) ReadSoulchest() {
	data := self.Data2["soulchest"]
	for _, value := range data {
		bag := HF_Atoi(value["bag"])
		if bag/100 == 1 && (self.Soulchest_Max == 0 || bag-self.Soulchest_Max == 1) {
			self.Soulchest_Max = bag
		}
		self.Soulchest_CSV[bag] = append(self.Soulchest_CSV[bag], value)
		self.Soulchest_SUM[bag] += HF_Atoi(value["probability2"])
	}
}

func (self *CsvMgr) GetSoulchest(bag int) []CsvNode {
	return self.Soulchest_CSV[bag]
}

func (self *CsvMgr) GetSoulchestPro(bag int) int {
	node := self.Soulchest_CSV[bag]
	for i := 0; i < len(node); i++ {
		if HF_Atoi(node[i]["probability1"]) != 0 {
			return HF_Atoi(node[i]["probability1"])
		}
	}

	return 0
}

func (self *CsvMgr) GetPromotebox(id int) (CsvNode, bool) {
	value, ok := self.PromoteBox_CSV[id]
	return value, ok
}

func (self *CsvMgr) GetSign(month, sign int) (*SignConfig, bool) {
	value, ok := self.SignMap[month*10000+sign]
	return value, ok
}

func (self *CsvMgr) ReadWarTask() {
	for _, value := range self.Data2["wartask"] {
		node := new(TaskNode)
		node.Id = HF_Atoi(value["taskid"])
		node.Tasktypes = HF_Atoi(value["tasktype"])
		node.N1 = HF_Atoi(value["n1"])
		node.N2 = HF_Atoi(value["n2"])
		node.N3 = HF_Atoi(value["n3"])
		node.N4 = HF_Atoi(value["n4"])
		self.WarTaskTask_CSV[node.Id] = node
		self.War_CSV[node.Id] = value
	}
}

func (self *CsvMgr) ReadWarTargetTask() {
	for _, value := range self.Data2["wartarget"] {
		node := new(TaskNode)
		node.Id = HF_Atoi(value["id"])
		node.Tasktypes = HF_Atoi(value["tasktypes"])
		node.N1 = HF_Atoi(value["n1"])
		node.N2 = HF_Atoi(value["n2"])
		node.N3 = HF_Atoi(value["n3"])
		node.N4 = HF_Atoi(value["n4"])
		self.WarTargetTask_CSV[node.Id] = node
		self.WarTarget_CSV[node.Id] = value
	}
}

func (self *CsvMgr) ReadSevenDayTask() {
	for i := 0; i < len(self.Data2["Sevenday"]); i++ {
		node := new(TaskNode)
		node.Id = HF_Atoi(self.Data2["Sevenday"][i]["id"])
		node.Tasktypes = HF_Atoi(self.Data2["Sevenday"][i]["tasktypes"])
		node.N1 = HF_Atoi(self.Data2["Sevenday"][i]["n1"])
		node.N2 = HF_Atoi(self.Data2["Sevenday"][i]["n2"])
		node.N3 = HF_Atoi(self.Data2["Sevenday"][i]["n3"])
		node.N4 = HF_Atoi(self.Data2["Sevenday"][i]["n4"])
		self.SevenDayTask_CSV[node.Id] = node
		self.SevenStatus[node.Id] = HF_Atoi(self.Data2["Sevenday"][i]["touch"])
	}
	//litter.Dump(self.SevenStatus)
}

func (self *CsvMgr) ReadSevenDay() {
	data := self.Data2["Sevenday"]
	for _, value := range data {
		self.SevenDay_CSV[HF_Atoi(value["id"])] = value
	}
}

func (self *CsvMgr) GetSevenDay(id int) (CsvNode, bool) {
	value, ok := self.SevenDay_CSV[id]
	return value, ok
}

func (self *CsvMgr) SplitStringToInt(s string) []int {

	str := strings.Split(s, "||")
	rel := make([]int, 0)
	for i := 0; i < len(str); i++ {
		num, err := strconv.Atoi(str[i])
		if err != nil {
			LogError(fmt.Sprintf("SplitStringToInt:%s", s))
			return nil
		}
		rel = append(rel, num)
	}

	return rel
}

func (self *CsvMgr) ReadHalfMoonTask() {
	for i := 0; i < len(self.Data2["Halfmoon"]); i++ {
		node := new(TaskNode)
		node.Id = HF_Atoi(self.Data2["Halfmoon"][i]["id"])
		node.Tasktypes = HF_Atoi(self.Data2["Halfmoon"][i]["tasktypes"])
		node.N1 = HF_Atoi(self.Data2["Halfmoon"][i]["n1"])
		node.N2 = HF_Atoi(self.Data2["Halfmoon"][i]["n2"])
		node.N3 = HF_Atoi(self.Data2["Halfmoon"][i]["n3"])
		node.N4 = HF_Atoi(self.Data2["Halfmoon"][i]["n4"])
		sortNum := HF_Atoi(self.Data2["Halfmoon"][i]["sort"])

		self.HalfMoonTask_CSV[node.Id] = node
		self.HalfMonnStatus[node.Id] = HF_Atoi(self.Data2["Halfmoon"][i]["touch"])

		if node.Tasktypes == TASK_TYPE_LOGIN_DAY {
			if _, ok := self.HalfMoonLogin[sortNum]; !ok {
				self.HalfMoonLogin[sortNum] = node.Id
			}
		} else if node.Tasktypes == TrialTask {
			self.HalfMoonTrial[node.Id] = &HalfMoonTrial{
				TaskId: node.Id,
				Index:  node.N2,
				Hard:   node.N3,
				Sort:   sortNum,
			}
		} else if node.Tasktypes == BeautyAdvance {
			if node.N2 == 0 { // 把美人等级筛选出去
				self.HalfMoonBeauty[node.Id] = &HalfMoonBeauty{
					TaskId: node.Id,
					MaxLv:  node.N3,
					Sort:   sortNum,
				}
			}

		}
	}
}

func (self *CsvMgr) ReadHalfMoon() {
	data := self.Data2["Halfmoon"]
	for _, value := range data {
		self.HalfMoon_CSV[HF_Atoi(value["id"])] = value
	}
}

func (self *CsvMgr) GetHalfMoon(id int) (CsvNode, bool) {
	value, ok := self.HalfMoon_CSV[id]
	return value, ok
}

func (self *CsvMgr) ReadGemsweeper_event() {
	data := self.Data2["Gemsweeper_event"]
	for _, value := range data {
		self.Gemsweeper_event_CSV[HF_Atoi(value["group"])*100000+HF_Atoi(value["cycle"])*100+HF_Atoi(value["step"])] = value
	}
}

func (self *CsvMgr) GetGemsweeper_event(group int, cycle int, step int) (CsvNode, bool) {
	value, ok := self.Gemsweeper_event_CSV[group*100000+cycle*100+step]

	return value, ok
}

func (self *CsvMgr) LoadPlayerName() {
	data := self.PlayernameConfig
	for _, value := range data {
		if len(value.Names) < 3 {
			continue
		}
		self.Name1_CSV = append(self.Name1_CSV, value.Names[0])
		self.Name2_CSV = append(self.Name2_CSV, value.Names[1])
		self.Name3_CSV = append(self.Name3_CSV, value.Names[2])
	}
}

func (self *CsvMgr) GetName() string {
	return self.Name1_CSV[HF_GetRandom(len(self.Name1_CSV))] + self.Name2_CSV[HF_GetRandom(len(self.Name2_CSV))]
}

func (self *CsvMgr) ReadExciting() {
	data := self.Data2["World_Event"]
	for _, value := range data {
		self.Exciting_CSV[HF_Atoi(value["power"])*100+HF_Atoi(value["id"])] = value
	}
}

func (self *CsvMgr) GetExciting(power int, id int) (CsvNode, bool) {
	value, ok := self.Exciting_CSV[power*100+id]

	return value, ok
}

func (self *CsvMgr) ReadDiplomacy() {
	data := self.Data2["Diplomacy_changes"]
	for _, value := range data {
		self.Diplomacy_CSV[HF_Atoi(value["Force_type"])*1000+HF_Atoi(value["index"])] = value
	}
}

func (self *CsvMgr) GetDiplomacy(camp int, index int) (CsvNode, bool) {
	value, ok := self.Diplomacy_CSV[camp*1000+index]

	return value, ok
}

func (self *CsvMgr) ReadHonorkill() {
	data := self.Data2["honorkill"]
	for _, value := range data {
		self.Honorkill_CSV[HF_Atoi(value["honortype"])*1000+HF_Atoi(value["honorkill"])] = value
	}
}

func (self *CsvMgr) GetHonorkill(honortype int, honorkill int) int {
	value, ok := self.Honorkill_CSV[honortype*1000+honorkill]
	if ok {
		return HF_Atoi(value["honoritemnum"])
	}

	return 0
}

func (self *CsvMgr) ReadHomeoffice_ranking() {
	data := self.Data2["Homeoffice_Ranking"]
	for _, value := range data {
		self.Homeoffice_ranking_CSV[HF_Atoi(value["type"])*1000+HF_Atoi(value["rank"])] = value
	}
}

func (self *CsvMgr) GetHomeoffice_ranking(_type int, rank int) (CsvNode, bool) {
	value, ok := self.Homeoffice_ranking_CSV[_type*1000+rank]
	return value, ok
}

func (self *CsvMgr) ReadCampreward() {
	data := self.Data2["campreward"]
	for _, value := range data {
		self.Reward_CSV[HF_Atoi(value["id"])*1000+HF_Atoi(value["step"])] = value
	}
}

func (self *CsvMgr) GetCampreward(id int, step int) (CsvNode, bool) {
	value, ok := self.Reward_CSV[id*1000+step]
	return value, ok
}

func (self *CsvMgr) ReadGemsweeper_extradrop() {
	data := self.Data2["Gemsweeper_extradrop"]
	for _, value := range data {
		week := HF_Atoi(value["week"])
		if week == 7 {
			week = 0
		}
		self.Gemsweeper_extradrop_CSV[HF_Atoi(value["group"])*1000+week] = value
	}
}

func (self *CsvMgr) ReadWarCity() {
	for i := 0; i < 3; i++ {
		self.WarCity_Map[i] = make(map[int]*WorldMapConfig)
	}

	for _, value := range self.WorldMap {
		for j := 0; j < 3; j++ {
			if value.Sequence[j] == 0 {
				continue
			}

			self.WarCity_Map[j][value.Sequence[j]] = value
		}
	}
}

func (self *CsvMgr) ReadActivityType() {
	for _, value := range self.Data2["Activitynew"] {
		find := false
		for i := 0; i < len(self.Activity_Type); i++ {
			if self.Activity_Type[i] == HF_Atoi(value["id"])/100 {
				find = true
			}
		}

		if find == false {
			self.Activity_Type = append(self.Activity_Type, HF_Atoi(value["id"])/100)
		}
	}

	for _, value := range self.Data2["Activitynew1"] {
		find := false
		for i := 0; i < len(self.Activity_Type); i++ {
			if self.Activity_Type[i] == HF_Atoi(value["id"])/100 {
				find = true
			}
		}

		if find == false {
			self.Activity_Type = append(self.Activity_Type, HF_Atoi(value["id"])/100)
		}
	}
}

//func (self *CsvMgr) ReadActivityBox() {
//	for _, value := range self.Data2["Activitybox"] {
//		group := HF_Atoi(value["gearid"])
//		_, ok := self.ActivityBox_CSV[group]
//		if !ok {
//			self.ActivityBox_CSV[group] = make([]CsvNode, 0)
//		}
//
//		self.ActivityBox_CSV[group] = append(self.ActivityBox_CSV[group], value)
//	}
//}

func (self *CsvMgr) ReadStateCity() {
	data, _ := GetCsvMgr().Data["World_Map"]
	for _, value := range data {
		stateid := HF_Atoi(value["state"])
		if stateid > 0 {
			_, ok := self.State_City_Num[stateid]
			if !ok {
				self.State_City_Num[stateid] = make([]int, 0)
			}
			self.State_City_Num[stateid] = append(self.State_City_Num[stateid], HF_Atoi(value["id"]))
		}
	}
}

func (self *CsvMgr) ReadGemsweeper_itembag() {
	data := self.Data2["Gemsweeper_itembag"]

	for _, value := range data {
		self.Gemsweeper_itembag_SUM[HF_Atoi(value["group"])] += HF_Atoi(value["gvalue"])
		self.Gemsweeper_itembag_GG_CSV[(HF_Atoi(value["group"])*100 + HF_Atoi(value["gid"]))] = append(self.Gemsweeper_itembag_GG_CSV[(HF_Atoi(value["group"])*100 + HF_Atoi(value["gid"]))], value)
	}
}

func (self *CsvMgr) GetGemsweeper_itembag_GG(group int, gid int) ([]CsvNode, bool) {
	value, ok := self.Gemsweeper_itembag_GG_CSV[group*100+gid]
	return value, ok
}

func (self *CsvMgr) GetGemsweeper_extradrop(group int, week int) (CsvNode, bool) {
	value, ok := self.Gemsweeper_extradrop_CSV[group*1000+week]
	return value, ok
}

func (self *CsvMgr) ReadHorseParam() {
	data := self.Data2["Horse_Parm"]
	for _, value := range data {
		self.HorseParam_CSV[value["name"]] = value

		if value["name"] == "Switch_WarHorse_1" {
			self.HorseSwitchId = HF_Atoi(value["parm1"])
			if self.HorseSwitchId == 0 {
				LogError("战马转换道具消耗不存在")
			}
		} else if value["name"] == "Switch_WarHorse_2" {
			self.HorseSwitchNum = HF_Atoi(value["parm1"])
		}

		if value["name"] == "Clear_WarHorse_1" {
			self.ClearWarHorseId = HF_Atoi(value["parm1"])
			if self.ClearWarHorseId == 0 {
				LogError("战马洗练道具消耗不存在")
			}
		} else if value["name"] == "Clear_WarHorse_2" {
			self.ClearWarHorseNum = HF_Atoi(value["parm1"])
		}

	}
}

func (self *CsvMgr) GetHorseParamInt(name string) int {
	value, ok := self.HorseParam_CSV[name]
	if ok {
		return HF_Atoi(value["parm1"])
	}

	return 0
}

func (self *CsvMgr) GetHorseParamString(name string) string {
	value, ok := self.HorseParam_CSV[name]
	if ok {
		return value["parm1"]
	}

	return ""
}

func (self *CsvMgr) ReadHorseAttr() {
	data, _ := GetCsvMgr().Data2["Horse_BattleSteed_Attribute"]
	for _, value := range data {
		group := HF_Atoi(value["attribute_list_group"])
		if group > 0 {
			_, ok := self.HorseAttr_CSV[group]
			if !ok {
				self.HorseAttr_CSV[group] = make([]CsvNode, 0)
			}
			self.HorseAttr_CSV[group] = append(self.HorseAttr_CSV[group], value)
		}
	}
}

// 加载排行榜奖励
func (self *CsvMgr) LoadTimeGeneralRank() {
	for _, value := range self.TimeGeneralsRankawardConfig {
		if self.GeneralRankMail.MailTitle == "" {
			self.GeneralRankMail.MailTitle = value.MailTitle
			if len(value.MainTxt) == 3 {
				self.GeneralRankMail.MailText1 = value.MainTxt[0]
				self.GeneralRankMail.MailText2 = value.MainTxt[1]
				self.GeneralRankMail.MailText3 = value.MainTxt[2]
			}
		}

		var normalAward []PassItem
		for i := 0; i < len(value.Normalawards); i++ {
			if len(value.Normalawards) != len(value.Normalnums) {
				LogError("len(value.Extraawards) = len(value.Extranums), ", len(value.Extraawards), len(value.Extranums))
				continue
			}

			normalAward = append(normalAward, PassItem{value.Normalawards[i], value.Normalnums[i]})
		}

		var extraAward []PassItem
		for j := 0; j < len(value.Extraawards); j++ {
			if len(value.Extraawards) != len(value.Extranums) {
				LogError("len(value.Extraawards) = len(value.Extranums), ", len(value.Extraawards), len(value.Extranums))
				continue
			}

			extraAward = append(extraAward, PassItem{value.Extraawards[j], value.Extranums[j]})
		}

		self.TimeGeneralRankLst = append(self.TimeGeneralRankLst, &TimeGeneralRank{
			Id:          value.Id,
			Group:       value.Group,
			RankMin:     value.Rankmin,
			RankMax:     value.Rankmax,
			NeetPoint:   value.Needpoint,
			NormalAward: normalAward,
			ExtraAward:  extraAward,
		})
	}

}

// 获取当前活动的配置
// 现在的逻辑是通过活动Id找到keyId,通过keyId找到对应的神将Id
func (self *CsvMgr) GetLootConfig() (*TimeGeneralsConfig, int) {
	// 没有这个活动
	activity := GetActivityMgr().GetActivity(ACT_GENERAL)
	//LogError("----------------222-------------------")
	if activity == nil {
		//LogError("GetLootConfig activity == nil")
		return nil, 0
	}
	//LogError("----------------111-------------------")
	keyId := activity.getTaskN4()
	step := activity.getTaskN3()
	if keyId == 0 {
		//LogError("GetLootConfig task n4 == 0")
		return nil, 0
	}
	var lootConfig *TimeGeneralsConfig
	timeConfig := self.TimeGeneralsConfig
	for _, value := range timeConfig {
		//LogError("---------------------value.KeyId:", value.KeyId)
		if value.KeyId == keyId {
			lootConfig = value
			break
		}
	}

	if lootConfig == nil {
		LogError("配置错误: lootConfig == nil, keyId:", keyId)
	}

	return lootConfig, step
}

// 获取排行榜奖励配置
func (self *CsvMgr) GetRankAwardConf(groupId int) []*TimeGeneralRank {
	rankConfig := self.TimeGeneralRankLst
	var res []*TimeGeneralRank
	for index := range rankConfig {
		if rankConfig[index] == nil {
			continue
		}
		if rankConfig[index].Group != groupId {
			continue
		}

		res = append(res, rankConfig[index])
	}
	return res
}

// 根据排行和积分获得奖励
func (self *CsvMgr) GetGeneralAward(rank, point int) *TimeGeneralRank {
	lootConfig, _ := self.GetLootConfig()
	if lootConfig == nil {
		LogError("GetGeneralAward lootConfig == nil")
		return &TimeGeneralRank{}
	}

	rankConfig := self.GetRankAwardConf(lootConfig.Id)
	if rankConfig == nil || len(rankConfig) < 0 {
		LogError("GetGeneralAward rankConfig == nil || len(rankConfig) < 0, id :", lootConfig.Id)
		return &TimeGeneralRank{}
	}

	var pRankConfig *TimeGeneralRank
	for index := range rankConfig {
		if rankConfig[index] == nil {
			LogError("rankConfig[index] == nil")
			continue
		}

		elem := rankConfig[index]
		if rank >= elem.RankMin && rank <= elem.RankMax {
			pRankConfig = elem
			break
		}
	}

	if pRankConfig == nil {
		LogError("限时神将排行榜配置错误")
		return &TimeGeneralRank{}
	}

	return pRankConfig
}

func (self *CsvMgr) ReadTigerAdvancedConfig() {
	data, ok := GetCsvMgr().Data["Tiger_Advanced"]
	if !ok {
		LogError("Tiger_Advanced is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		group := HF_Atoi(configValue["group"])
		var Costs []int
		for id := 1; id <= 2; id++ {
			Costs = append(Costs, HF_Atoi(configValue[fmt.Sprintf("cost%d", id)]))
		}
		var Nums []int
		for id := 1; id <= 2; id++ {
			Nums = append(Nums, HF_Atoi(configValue[fmt.Sprintf("num%d", id)]))
		}
		var Costitems []int
		for id := 1; id <= 6; id++ {
			Costitems = append(Costitems, HF_Atoi(configValue[fmt.Sprintf("cost_item%d", id)]))
		}
		var Costnums []int
		for id := 1; id <= 6; id++ {
			Costnums = append(Costnums, HF_Atoi(configValue[fmt.Sprintf("cost_num%d", id)]))
		}
		icon := HF_Atoi(configValue["icon"])
		self.TigerAdvancedConfig = append(self.TigerAdvancedConfig, &TigerAdvancedConfig{
			Group:     group,
			Costs:     Costs,
			Nums:      Nums,
			Costitems: Costitems,
			Costnums:  Costnums,
			Icon:      icon,
		})
	}
	//for index := range self.TigerAdvancedConfig {
	//	fmt.Printf("%+v\n", self.TigerAdvancedConfig[index])
	//}
}

func (self *CsvMgr) ReadTigerAttributeConfig() {
	GetCsvUtilMgr().LoadCsv("Tiger_Attribute", &self.TigerAttributeConfig)
}

func (self *CsvMgr) ReadTigerStuntConfig() {
	data, ok := GetCsvMgr().Data["Tiger_Stunt"]
	if !ok {
		LogError("Tiger_Stunt is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		stunt_id := HF_Atoi(configValue["stunt_id"])
		upgrade_model := HF_Atoi(configValue["upgrade_model"])
		hufu_limits := HF_Atoi(configValue["hufu_limits"])

		self.TigerStuntConfig = append(self.TigerStuntConfig, &TigerStuntConfig{
			Stuntid:      stunt_id,
			Upgrademodel: upgrade_model,
			HufuLimits:   hufu_limits,
		})
	}
}

func (self *CsvMgr) ReadTigerStuntUpgradeConfig() {
	data, ok := GetCsvMgr().Data2["Tiger_StuntUpgrade"]
	if !ok {
		LogError("Tiger_StuntUpgrade is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		stunt_group := HF_Atoi(configValue["stunt_group"])
		limitUpgradeId := HF_Atoi(configValue["limit_upgradeid"])
		level := HF_Atoi(configValue["level"])
		//attribute_type := HF_Atoi(configValue["attribute_type"])
		//value := HF_Atoi(configValue["value"])
		var attribute_type []int
		for id := 1; id <= 2; id++ {
			attribute_type = append(attribute_type, HF_Atoi(configValue[fmt.Sprintf("attribute_type%d", id)]))
		}
		var value []int64
		for id := 1; id <= 2; id++ {
			value = append(value, int64(HF_Atoi(configValue[fmt.Sprintf("value%d", id)])))
		}

		cost_item := HF_Atoi(configValue["cost_item"])
		cost_num := HF_Atoi(configValue["cost_num"])
		Reset_item := HF_Atoi(configValue["Reset_item"])
		Reset_num := HF_Atoi(configValue["Reset_num"])
		Return_item := HF_Atoi(configValue["Return_item"])
		Return_num := HF_Atoi(configValue["Return_num"])
		self.TigerStuntUpgradeConfig = append(self.TigerStuntUpgradeConfig, &TigerStuntUpgradeConfig{
			Stuntgroup:     stunt_group,
			LimitUpgradeId: limitUpgradeId,
			Level:          level,
			Attributetype:  attribute_type,
			Value:          value,
			Costitem:       cost_item,
			Costnum:        cost_num,
			Resetitem:      Reset_item,
			Resetnum:       Reset_num,
			Returnitem:     Return_item,
			Returnnum:      Return_num,
		})
	}
}

func (self *CsvMgr) ReadTigerSymbolConfig() {
	data, ok := GetCsvMgr().Data["Tiger_Symbol"]
	if !ok {
		LogError("Tiger_Symbol is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		ID := HF_Atoi(configValue["ID"])
		var Holes []int
		for id := 1; id <= 10; id++ {
			Holes = append(Holes, HF_Atoi(configValue[fmt.Sprintf("hole_%d", id)]))
		}
		self.TigerSymbolConfig = append(self.TigerSymbolConfig, &TigerSymbolConfig{
			ID:    ID,
			Holes: Holes,
		})
	}
	//litter.Dump(self.TigerSymbolConfig)
}

func (self *CsvMgr) ReadTigerUpgradeConfig() {
	GetCsvUtilMgr().LoadCsv("Tiger_Upgrade", &self.TigerUpgradeConfig)
	//litter.Dump(self.TigerUpgradeConfig)
}

func (self *CsvMgr) GetItemName(itemId int) string {
	itemName := ""
	config, ok := GetCsvMgr().ItemMap[itemId]
	if !ok {
		itemName = self.GetText("STR_ITEM_ERROR")
	} else {
		itemName = config.ItemName
	}
	return itemName
}

func (self *CsvMgr) ReadCitygvgConfig() {
	data, ok := GetCsvMgr().Data2["citygvg"]
	if !ok {
		LogError("citygvg is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		citysize := HF_Atoi(configValue["citysize"])
		var Items []int
		for id := 1; id <= 3; id++ {
			Items = append(Items, HF_Atoi(configValue[fmt.Sprintf("item%d", id)]))
		}
		var Nums []int
		for id := 1; id <= 3; id++ {
			Nums = append(Nums, HF_Atoi(configValue[fmt.Sprintf("num%d", id)]))
		}
		var Conditions []int
		for id := 1; id <= 3; id++ {
			Conditions = append(Conditions, HF_Atoi(configValue[fmt.Sprintf("condition%d", id)]))
		}
		position := HF_Atoi(configValue["position"])
		name := configValue["name"]
		number := HF_Atoi(configValue["number"])
		addition := HF_Atoi(configValue["addition"])
		reduce := HF_Atoi(configValue["reduce"])
		var Limits []int
		for id := 1; id <= 3; id++ {
			Limits = append(Limits, HF_Atoi(configValue[fmt.Sprintf("limit%d", id)]))
		}
		limitspecial := HF_Atoi(configValue["limitspecial"])
		repress := HF_Atoi(configValue["repress"])
		self.CitygvgConfig = append(self.CitygvgConfig, &CitygvgConfig{
			Citysize:     citysize,
			Items:        Items,
			Nums:         Nums,
			Conditions:   Conditions,
			Position:     position,
			Name:         name,
			Number:       number,
			Addition:     addition,
			Reduce:       reduce,
			Limits:       Limits,
			Limitspecial: limitspecial,
			Repress:      repress,
		})
	}
}

func (self *CsvMgr) ReadCityrandomConfig() {
	data, ok := GetCsvMgr().Data2["cityrandom"]
	if !ok {
		LogError("cityrandom is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		id := HF_Atoi(configValue["id"])
		typeV := HF_Atoi(configValue["type"])
		limit := HF_Atoi(configValue["limit"])
		weight := HF_Atoi(configValue["weight"])
		item := HF_Atoi(configValue["item"])
		num := HF_Atoi(configValue["num"])
		time := HF_Atoi(configValue["time"])
		min := HF_Atoi(configValue["min"])
		max := HF_Atoi(configValue["max"])
		self.CityrandomConfig = append(self.CityrandomConfig, &CityrandomConfig{
			Id:     id,
			Type:   typeV,
			Limit:  limit,
			Weight: weight,
			Item:   item,
			Num:    num,
			Time:   time,
			Min:    min,
			Max:    max,
		})
	}
}

func (self *CsvMgr) ReadCityrepressConfig() {
	data, ok := GetCsvMgr().Data["cityrepress"]
	if !ok {
		LogError("cityrepress is not exists!")
		return
	}
	for _, configValue := range data {
		id := HF_Atoi(configValue["id"])
		class := HF_Atoi(configValue["class"])
		name := configValue["name"]
		icon := configValue["icon"]
		typeV := HF_Atoi(configValue["type"])
		drop := HF_Atoi(configValue["drop"])
		num := HF_Atoi(configValue["num"])
		time := HF_Atoi(configValue["time"])
		txt := configValue["txt"]
		self.CityrepressConfig[id] = &CityrepressConfig{
			Id:    id,
			Class: class,
			Name:  name,
			Icon:  icon,
			Type:  typeV,
			Drop:  drop,
			Num:   num,
			Time:  time,
			Txt:   txt,
		}
	}
}

func (self *CsvMgr) getConfigById(redId int) *RedpacketmoneyConfig {
	config := self.RedpacketmoneyConfig
	var pConfig *RedpacketmoneyConfig
	for index := range config {
		if config[index].Id == redId {
			pConfig = config[index]
			break
		}
	}

	return pConfig
}

func (self *CsvMgr) getRedPacConfig(redId int) *RedpacketConfig {
	config := self.RedpacketConfig
	var pConfig *RedpacketConfig
	for index := range config {
		if config[index].Id == redId {
			pConfig = config[index]
			break
		}
	}

	return pConfig
}

func (self *CsvMgr) getRedPacConfByItemId(itemId int) *RedpacketConfig {
	config := self.RedpacketConfig
	var pConfig *RedpacketConfig
	for index := range config {
		if config[index].Item == itemId {
			pConfig = config[index]
			break
		}
	}

	return pConfig
}

// 获取经验加成
func (self *CsvMgr) getExpFactor(level int) float32 {
	// 遍历
	curFactor := float32(1.0)
	if len(self.ExpspeedupConfig) <= 0 {
		return curFactor
	}

	for index := range self.ExpspeedupConfig {
		elem := self.ExpspeedupConfig[index]
		if level >= elem.Leveldf1 && level <= elem.Leveldf2 {
			return float32(elem.Speedup / 10000.0)
		}
	}

	return curFactor
}

func (self *CsvMgr) getSendFreshTime(itemId int) int64 {
	config := GetCsvMgr().getRedPacConfByItemId(itemId)
	now := TimeServer()
	timeSet := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local)
	if config.Ceilingtime == 0 {
		return 0
	} else if config.Ceilingtime == 1 {
		if now.Unix() > timeSet.Unix() {
			return timeSet.Unix() + 86400
		} else {
			return timeSet.Unix()
		}
	} else if config.Ceilingtime == 2 {
		return self.getNextWeek().Unix()
	}
	return 0
}

func (self *CsvMgr) getGotFreshTime(itemId int) int64 {
	config := GetCsvMgr().getRedPacConfByItemId(itemId)
	now := TimeServer()
	timeSet := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local)
	if config.Upperlimittime == 0 {
		return 0
	} else if config.Upperlimittime == 1 {
		if now.Unix() > timeSet.Unix() {
			return timeSet.Unix() + 86400
		} else {
			return timeSet.Unix()
		}
	} else if config.Upperlimittime == 2 {
		return self.getNextWeek().Unix()
	}
	return 0
}

// 获取周一的时间
func (self *CsvMgr) getNextWeek() time.Time {
	now := TimeServer()
	week := NewTimeUtil(now).Monday()
	if now.Unix() < week.Unix() {
		return week.Add(time.Hour * 5)
	} else {
		return week.AddDate(0, 0, 7).Add(time.Hour * 5)
	}
}

func (self *CsvMgr) ReadTreasureAwakenConfig() {
	data, ok := GetCsvMgr().Data2["Treasure_Awaken"]
	if !ok {
		LogError("Treasure_Awaken is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		group := HF_Atoi(configValue["group"])
		star := HF_Atoi(configValue["star"])
		var Costitems []int
		for id := 1; id <= 2; id++ {
			Costitems = append(Costitems, HF_Atoi(configValue[fmt.Sprintf("cost_item%d", id)]))
		}
		var Costnums []int
		for id := 1; id <= 2; id++ {
			Costnums = append(Costnums, HF_Atoi(configValue[fmt.Sprintf("cost_num%d", id)]))
		}
		var Attributetypes []int
		for id := 1; id <= 6; id++ {
			Attributetypes = append(Attributetypes, HF_Atoi(configValue[fmt.Sprintf("attribute_type%d", id)]))
		}
		var Attributevalues []int64
		for id := 1; id <= 6; id++ {
			Attributevalues = append(Attributevalues, int64(HF_Atoi(configValue[fmt.Sprintf("attribute_value%d", id)])))
		}
		skill := HF_Atoi(configValue["skill"])
		reset_item := HF_Atoi(configValue["reset_item"])
		reset_num := HF_Atoi(configValue["reset_num"])
		var Returnitems []int
		for id := 1; id <= 2; id++ {
			Returnitems = append(Returnitems, HF_Atoi(configValue[fmt.Sprintf("return_item%d", id)]))
		}
		var Returnnums []int
		for id := 1; id <= 2; id++ {
			Returnnums = append(Returnnums, HF_Atoi(configValue[fmt.Sprintf("return_num%d", id)]))
		}
		quality := HF_Atoi(configValue["quality"])
		self.TreasureAwakenConfig = append(self.TreasureAwakenConfig, &TreasureAwakenConfig{
			Group:           group,
			Star:            star,
			Costitems:       Costitems,
			Costnums:        Costnums,
			Attributetypes:  Attributetypes,
			Attributevalues: Attributevalues,
			Skill:           skill,
			Resetitem:       reset_item,
			Resetnum:        reset_num,
			Returnitems:     Returnitems,
			Returnnums:      Returnnums,
			Quality:         quality,
		})

		awakenStar, ok := self.MaxTreasureAwaken[group]
		if !ok {
			self.MaxTreasureAwaken[group] = star
		} else {
			if awakenStar < star {
				self.MaxTreasureAwaken[group] = star
			}
		}
	}
}

func (self *CsvMgr) ReadTreasureClearAttributeConfig() {
	data, ok := GetCsvMgr().Data2["Treasure_ClearAttribute"]
	if !ok {
		LogError("Treasure_ClearAttribute is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		attribute_list_group := HF_Atoi(configValue["attribute_list_group"])
		weight_attribute := HF_Atoi(configValue["weight_attribute"])
		attribute_type := HF_Atoi(configValue["attribute_type"])
		lv := HF_Atoi(configValue["lv"])
		var Oddslvs []int
		for id := 1; id <= 2; id++ {
			Oddslvs = append(Oddslvs, HF_Atoi(configValue[fmt.Sprintf("odds_lv%d", id)]))
		}
		attribute_valve := int64(HF_Atoi(configValue["attribute_valve"]))
		quality := HF_Atoi(configValue["quality"])
		self.TreasureClearAttributeConfig = append(self.TreasureClearAttributeConfig, &TreasureClearAttributeConfig{
			Attributelistgroup: attribute_list_group,
			Weightattribute:    weight_attribute,
			Attributetype:      attribute_type,
			Lv:                 lv,
			Oddslvs:            Oddslvs,
			Attributevalve:     attribute_valve,
			Quality:            quality,
		})

		holeLv, ok := self.MaxTreasureHoleLv[lv]
		if !ok {
			self.MaxTreasureHoleLv[attribute_list_group] = lv
		} else {
			if holeLv < lv {
				self.MaxTreasureHoleLv[attribute_list_group] = lv
			}
		}

	}
}

func (self *CsvMgr) ReadTreasureClearItemConfig() {
	data, ok := GetCsvMgr().Data2["Treasure_ClearItem"]
	if !ok {
		LogError("Treasure_ClearItem is not exists!")
		return
	}
	for _, configValue := range data {
		item_id := HF_Atoi(configValue["item_id"])
		effect_type := HF_Atoi(configValue["effect_type"])
		var Parms []int
		for id := 1; id <= 3; id++ {
			Parms = append(Parms, HF_Atoi(configValue[fmt.Sprintf("parm%d", id)]))
		}
		self.TreasureClearItemConfig[item_id] = &TreasureClearItemConfig{
			Itemid:     item_id,
			Effecttype: effect_type,
			Parms:      Parms,
		}
	}
}

func (self *CsvMgr) ReadTreasureEquipConfig() {
	data, ok := GetCsvMgr().Data2["Treasure_Equip"]
	if !ok {
		LogError("Treasure_Equip is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		id := HF_Atoi(configValue["id"])
		class := HF_Atoi(configValue["class"])
		position := HF_Atoi(configValue["position"])

		var Attributetypes []int
		for id := 1; id <= 4; id++ {
			Attributetypes = append(Attributetypes, HF_Atoi(configValue[fmt.Sprintf("attribute_type%d", id)]))
		}
		var Attributevalues []int64
		for id := 1; id <= 4; id++ {
			_, err := strconv.Atoi(configValue[fmt.Sprintf("attribute_value%d", id)])
			if err != nil {
				LogError(err.Error())
				continue
			}
			Attributevalues = append(Attributevalues, int64(HF_Atoi(configValue[fmt.Sprintf("attribute_value%d", id)])))
		}
		var Attributeholes []int
		for id := 1; id <= 3; id++ {
			Attributeholes = append(Attributeholes, HF_Atoi(configValue[fmt.Sprintf("attribute_hole%d", id)]))
		}
		add_effect := HF_Atoi(configValue["add_effect"])
		decompose_item := HF_Atoi(configValue["decompose_item"])
		decompose_num := HF_Atoi(configValue["decompose_num"])
		quality := HF_Atoi(configValue["quality"])
		decompose_group := HF_Atoi(configValue["decompose_group"])

		self.TreasureEquipConfig[id] = &TreasureEquipConfig{
			Id:              id,
			Class:           class,
			Position:        position,
			Attributetypes:  Attributetypes,
			Attributevalues: Attributevalues,
			Attributeholes:  Attributeholes,
			Addeffect:       add_effect,
			Decomposeitem:   decompose_item,
			Decomposenum:    decompose_num,
			Quality:         quality,
			DecomposeGroup:  decompose_group,
		}
	}
}

func (self *CsvMgr) ReadTreasureHeroConfig() {
	data, ok := GetCsvMgr().Data["Treasure_Hero"]
	if !ok {
		LogError("Treasure_Hero is not exists!")
		return
	}
	for _, configValue := range data {
		hero_id := HF_Atoi(configValue["hero_id"])
		class := HF_Atoi(configValue["class"])
		var Attributes []int
		for id := 1; id <= 6; id++ {
			Attributes = append(Attributes, HF_Atoi(configValue[fmt.Sprintf("attribute%d", id)]))
		}
		suit_group := HF_Atoi(configValue["suit_group"])
		awaken_group := HF_Atoi(configValue["awaken_group"])
		self.TreasureHeroConfig[hero_id] = &TreasureHeroConfig{
			Heroid:      hero_id,
			Class:       class,
			Attributes:  Attributes,
			Suitgroup:   suit_group,
			Awakengroup: awaken_group,
		}
	}
}

func (self *CsvMgr) ReadTreasureSuitAttributeConfig() {
	data, ok := GetCsvMgr().Data2["Treasure_SuitAttribute"]
	if !ok {
		LogError("Treasure_SuitAttribute is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		group := HF_Atoi(configValue["group"])
		suit_num := HF_Atoi(configValue["suit_num"])
		var Attributetypes []int
		for id := 1; id <= 4; id++ {
			Attributetypes = append(Attributetypes, HF_Atoi(configValue[fmt.Sprintf("attribute_type%d", id)]))
		}
		var Attributevalues []int64
		for id := 1; id <= 4; id++ {
			Attributevalues = append(Attributevalues, int64(HF_Atoi(configValue[fmt.Sprintf("attribute_value%d", id)])))
		}
		self.TreasureSuitAttributeConfig = append(self.TreasureSuitAttributeConfig, &TreasureSuitAttributeConfig{
			Group:           group,
			Suitnum:         suit_num,
			Attributetypes:  Attributetypes,
			Attributevalues: Attributevalues,
		})
	}
}

// 获取转盘配置
func (self *CsvMgr) GetDialConfig() []*LuckyturntableConfig {
	// 没有这个活动
	activity := GetActivityMgr().GetActivity(ACT_DIAL)
	if activity == nil {
		return nil
	}

	n4 := activity.getTaskN4()
	if n4 == 0 {
		LogError("GetLootConfig task n4 == 0")
		return nil
	}

	var lootConfig []*LuckyturntableConfig
	dialConfig := self.LuckyturntableConfig
	for index, value := range dialConfig {
		if value.N4 == n4 {
			lootConfig = append(lootConfig, dialConfig[index])
		}
	}

	return lootConfig
}

// 通过幸运值进行转盘掉落
func (self *CsvMgr) DoDialLoot(luck int, configlst []*LuckyturntableConfig) (*LuckyturntableConfig, error) {
	// 获得总权重
	chance := 0
	for index := range configlst {
		pConfig := configlst[index]
		luckTargets := pConfig.Lucktargets
		luckValues := pConfig.Luckvalues

		if len(luckTargets) != 2 {
			return nil, errors.New(GetCsvMgr().GetText("STR_CSVMGR_TURNTABLE_CONFIGURATION_ERROR"))
		}

		if len(luckValues) != 2 {
			return nil, errors.New(GetCsvMgr().GetText("STR_CSVMGR_TURNTABLE_CONFIGURATION_ERROR"))
		}

		target1 := luckTargets[0]
		target2 := luckTargets[1]
		value1 := luckValues[0]
		value2 := luckValues[1]

		if luck < target1 {
			chance += pConfig.Value
		} else if luck >= target1 && luck < target2 {
			chance += value1
		} else if luck >= target2 {
			chance += value2
		}
	}

	if chance == 0 {
		return nil, errors.New(GetCsvMgr().GetText("STR_CSVMGR_THE_TOTAL_WEIGHT_OF_TURNTABLE"))
	}

	randNum := HF_GetRandom(chance) + 1

	check := 0
	var config *LuckyturntableConfig
	// 随机找到其中一项
	for index := range configlst {
		pConfig := configlst[index]
		luckTargets := pConfig.Lucktargets
		luckValues := pConfig.Luckvalues

		target1 := luckTargets[0]
		target2 := luckTargets[1]
		value1 := luckValues[0]
		value2 := luckValues[1]

		if luck < target1 {
			check += pConfig.Value
		} else if luck >= target1 && luck < target2 {
			check += value1
		} else if luck >= target2 {
			check += value2
		}

		if randNum <= check {
			config = configlst[index]
			break
		}
	}

	if config == nil {
		return nil, errors.New(GetCsvMgr().GetText("STR_CSVMGR_ROTARY_DISK_FALLING_CONFIGURATION_NOT"))
	}

	return config, nil
}

// 获取转盘配置
func (self *CsvMgr) GetDialBoxConfig() []*LuckyturntablelistConfig {
	// 没有这个活动
	activity := GetActivityMgr().GetActivity(ACT_DIAL)
	if activity == nil {
		return nil
	}

	n4 := activity.getTaskN4()
	if n4 == 0 {
		LogError("GetLootConfig task n4 == 0")
		return nil
	}

	var lootConfig []*LuckyturntablelistConfig
	dialConfig := self.LuckyturntablelistConfig
	for index, value := range dialConfig {
		if value.N4 == n4 {
			lootConfig = append(lootConfig, dialConfig[index])
		}
	}

	return lootConfig
}

// 获取翻牌配置
func (self *CsvMgr) GetDrawBoxConfig(step int) ([]*LuckdrawlistConfig, error) {
	var lootConfig []*LuckdrawlistConfig
	drawConfig := self.LuckdrawlistConfig
	for index, value := range drawConfig {
		if value.N4 == step {
			lootConfig = append(lootConfig, drawConfig[index])
		}
	}
	return lootConfig, nil
}

// 获取翻牌消耗
func (self *CsvMgr) GetDrawConfig(step int, group int) (*LuckdrawConfig, error) {
	var lootConfig []*LuckdrawConfig
	dialConfig := self.LuckdrawConfig
	for index, value := range dialConfig {
		if value.N4 == step {
			lootConfig = append(lootConfig, dialConfig[index])
		}
	}

	for index := range lootConfig {
		if lootConfig[index].Id == group {
			return lootConfig[index], nil
		}
	}

	return nil, errors.New(GetCsvMgr().GetText("STR_CSVMGR_UNABLE_TO_FIND_THE_CORRESPONDING"))
}

func (self *CsvMgr) GetCurDrawLoot(step int) []*LuckdrawConfig {
	var lootConfig []*LuckdrawConfig
	for _, v := range self.LuckdrawConfig {
		if v.N4 != step {
			continue
		}
		lootConfig = append(lootConfig, v)
	}
	return lootConfig
}

//! 找到分解的宝物配置
func (self *CsvMgr) GetTreasureEquip(itemId int) *TreasureEquipConfig {
	v, ok := self.TreasureEquipConfig[itemId]
	if ok {
		return v
	}
	return nil
}

//! 找到武将能穿的装备
func (self *CsvMgr) GetTreasureHero(heroid int) *TreasureHeroConfig {
	config, ok := self.TreasureHeroConfig[heroid]
	if !ok {
		return nil
	}

	return config
}

//! 先根据hole,lv 筛选配置
func (self *CsvMgr) GetTresureHole(holeId int, holeLv int) []*TreasureClearAttributeConfig {
	var res []*TreasureClearAttributeConfig
	for index := range self.TreasureClearAttributeConfig {
		pConfig := self.TreasureClearAttributeConfig[index]
		if pConfig == nil {
			continue
		}

		if pConfig.Attributelistgroup != holeId {
			continue
		}

		if pConfig.Lv != holeLv {
			continue
		}
		res = append(res, pConfig)
	}
	return res
}

func (self *CsvMgr) GetTresureHoleByAtt(holeId int, holeLv int, attr int) *TreasureClearAttributeConfig {
	for index := range self.TreasureClearAttributeConfig {
		pConfig := self.TreasureClearAttributeConfig[index]
		if pConfig == nil {
			continue
		}

		if pConfig.Attributelistgroup != holeId {
			continue
		}

		if pConfig.Lv != holeLv {
			continue
		}

		if pConfig.Attributetype != attr {
			continue
		}

		return pConfig
	}
	return nil
}

func (self *CsvMgr) ReadTreasureDecomposeConfig() {
	data, ok := GetCsvMgr().Data2["Treasure_Decompose"]
	if !ok {
		LogError("Treasure_Decompose is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		group := HF_Atoi(configValue["group"])
		total_level := HF_Atoi(configValue["total_level"])
		item_id := HF_Atoi(configValue["item_id"])
		value := HF_Atoi(configValue["value"])
		config := &TreasureDecomposeConfig{
			Group:      group,
			Totallevel: total_level,
			Itemid:     item_id,
			Value:      value,
		}
		self.TreasureDecomposeConfig = append(self.TreasureDecomposeConfig, config)
		infoMap, ok := self.TreasureDecomposeMap[group]
		if !ok {
			infoMap = make(map[int]*TreasureDecomposeConfig)
			infoMap[total_level] = config
			self.TreasureDecomposeMap[group] = infoMap
		} else {
			infoMap[total_level] = config
		}
	}
}

func (self *CsvMgr) GetTreasureDecomposeConfig(group int, level int) *TreasureDecomposeConfig {
	infoMap, ok := self.TreasureDecomposeMap[group]
	if !ok {
		return nil
	}

	pConfig, ok := infoMap[level]
	if !ok {
		return nil
	}

	return pConfig
}

// 获取觉醒消耗配置
func (self *CsvMgr) GetAwakenConfig(group int, star int) *TreasureAwakenConfig {
	for _, v := range self.TreasureAwakenConfig {
		if v.Group != group {
			continue
		}

		if v.Star != star {
			continue
		}
		return v
	}
	return nil
}

// 通过group和suitnum来获取套装属性
func (self *CsvMgr) GetSuitConfig(group int, suitnum int) *TreasureSuitAttributeConfig {
	for index := range self.TreasureSuitAttributeConfig {
		pConfig := self.TreasureSuitAttributeConfig[index]
		if pConfig == nil {
			continue
		}

		if pConfig.Group == group && pConfig.Suitnum == suitnum {
			return pConfig
		}
	}

	return nil
}

// 获取砸金蛋配置
func (self *CsvMgr) GetLuckEggGroup(actType int) (int, error) {
	// 没有这个活动
	activity := GetActivityMgr().GetActivity(actType)
	if activity == nil {
		return 0, errors.New(GetCsvMgr().GetText("STR_CSVMGR_GOLDEN_EGG_BREAKING_ACTIVITY_CONFIGURATION"))
	}

	n4 := activity.getTaskN4()
	if n4 == 0 {
		LogError("GetLuckEggGroup task n4 == 0")
		return 0, errors.New(GetCsvMgr().GetText("STR_CSVMGR_N4_IS_0"))
	}

	drawConfig := self.LuckegggroupConfig
	groupMap := make(map[int]bool)
	for _, value := range drawConfig {
		if value.N4 == n4 {
			groupMap[value.Group] = true
		}
	}

	var keys []int
	for key := range groupMap {
		keys = append(keys, key)
	}

	size := len(keys)
	if size == 0 {
		return 0, errors.New(GetCsvMgr().GetText("STR_CSVMGR_THE_GROUP_FIELD_COULD_NOT"))
	}

	index := HF_GetRandom(size)

	return keys[index], nil
}

// 获取翻牌配置
func (self *CsvMgr) GetLuckEggConfig(step int) []*LuckeggConfig {
	var lootConfig []*LuckeggConfig
	configSlice := self.LuckeggConfig
	for index, value := range configSlice {
		if value.N4 == step {
			lootConfig = append(lootConfig, configSlice[index])
		}
	}

	return lootConfig
}

// 找到砸金蛋道具
func (self *CsvMgr) GetLuckEggLoot(step int, group int, id int) *LuckegggroupConfig {
	// 获得总权重
	for index := range self.LuckegggroupConfig {
		pConfig := self.LuckegggroupConfig[index]
		if pConfig == nil {
			continue
		}
		if pConfig.N4 == step && pConfig.Group == group && pConfig.Id == id {
			return pConfig
		}
	}

	return nil
}

// 根据step找到对应的活动配置
func (self *CsvMgr) getLuckTask(step int) []*LuckstartConfig {
	var res []*LuckstartConfig
	for _, pTask := range self.LuckstartConfig {
		if pTask == nil {
			continue
		}

		if pTask.N4 != step {
			continue
		}
		res = append(res, pTask)
	}
	return res
}

// 获取开工福利任务
func (self *CsvMgr) getLuckStartTaskNode(step int, taskId int) *TaskNode {
	infoMap, ok := self.LuckStartMap[step]
	if !ok {
		return nil
	} else {
		taskNode, ok := infoMap[taskId]
		if !ok {
			return nil
		}
		return taskNode
	}
}

func (self *CsvMgr) getLuckStartConfig(step int, taskId int) *LuckstartConfig {
	infoMap, ok := self.LuckStartConfigMap[step]
	if !ok {
		return nil
	} else {
		config, ok := infoMap[taskId]
		if !ok {
			return nil
		}
		return config
	}
}

func (self *CsvMgr) LoadDailyrechargeConfig() {
	self.DailyrechargeMap = make(map[int][]*DailyrechargeConfig)
	for _, value := range self.DailyrechargeConfig {
		_, ok := self.DailyrechargeMap[value.Group]
		if !ok {
			self.DailyrechargeMap[value.Group] = make([]*DailyrechargeConfig, 0)
		}
		self.DailyrechargeMap[value.Group] = append(self.DailyrechargeMap[value.Group], value)
		//fmt.Printf("%+v\n", self.GetDailyRechargeConfig(n4, id))
	}
	//litter.Dump(self.DailyrechargeMap)

}

func (self *CsvMgr) GetDailyRechargeConfig(step int, id int) *DailyrechargeConfig {
	confArr, ok := self.DailyrechargeMap[step]
	if !ok {
		return nil
	}

	for i := 0; i <= len(confArr)-1; i++ {
		if confArr[i].Id == id {
			return confArr[i]
		}
	}

	return nil

	//infoMap, ok := self.DailyrechargeMap[step]
	//if !ok {
	//	return nil
	//}

	//config, ok := infoMap[id]
	//if !ok {
	//	return nil
	//}
	//return config
}

// 根据step找到对应的连续充值
func (self *CsvMgr) getDailyConfig(step int) []*DailyrechargeConfig {
	if confArr, ok := self.DailyrechargeMap[step]; ok == true {
		return confArr
	} else {
		return []*DailyrechargeConfig{}
	}

	//var res []*DailyrechargeConfig
	//for _, pTask := range self.DailyrechargeConfig {
	//	if pTask == nil {
	//		continue
	//	}
	//	if pTask.Group != step {
	//		continue
	//	}
	//	res = append(res, pTask)
	//}
	//return res
}

// 宝藏逻辑
func (self *CsvMgr) ReadGemsweeperConfig() {
	data, ok := GetCsvMgr().Data["Gemsweeper"]
	if !ok {
		LogError("Gemsweeper is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		id := HF_Atoi(configValue["id"])
		typeValue := HF_Atoi(configValue["type"])
		quality := HF_Atoi(configValue["quality"])
		gradelevel := HF_Atoi(configValue["gradelevel"])
		var Activationitems []int
		for id := 1; id <= 5; id++ {
			Activationitems = append(Activationitems, HF_Atoi(configValue[fmt.Sprintf("activation_item%d", id)]))
		}
		var Activationnums []int
		for id := 1; id <= 5; id++ {
			Activationnums = append(Activationnums, HF_Atoi(configValue[fmt.Sprintf("activation_num%d", id)]))
		}
		cost_item := HF_Atoi(configValue["cost_item"])
		cost_num := HF_Atoi(configValue["cost_num"])
		event_group := HF_Atoi(configValue["event_group"])
		duration := HF_Atoi(configValue["duration"])
		var Itemshows []int
		for id := 1; id <= 4; id++ {
			Itemshows = append(Itemshows, HF_Atoi(configValue[fmt.Sprintf("item_show%d", id)]))
		}
		self.GemsweeperConfig[id] = &GemsweeperConfig{
			Id:              id,
			Type:            typeValue,
			Quality:         quality,
			Gradelevel:      gradelevel,
			Activationitems: Activationitems,
			Activationnums:  Activationnums,
			Costitem:        cost_item,
			Costnum:         cost_num,
			Eventgroup:      event_group,
			Duration:        duration,
			Itemshows:       Itemshows,
		}
	}
}

func (self *CsvMgr) ReadGemsweepereventConfig() {
	data, ok := GetCsvMgr().Data2["Gemsweeper_event"]
	if !ok {
		LogError("Gemsweeper_event is not exists!")
		return
	}
	var keys []int
	for key := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		configValue := data[key]
		group := HF_Atoi(configValue["group"])
		cycle := HF_Atoi(configValue["cycle"])
		step := HF_Atoi(configValue["step"])
		next := HF_Atoi(configValue["next"])
		level := HF_Atoi(configValue["level"])
		event := HF_Atoi(configValue["event"])
		var Items []int
		for id := 1; id <= 5; id++ {
			Items = append(Items, HF_Atoi(configValue[fmt.Sprintf("item%d", id)]))
		}
		var Nums []int
		for id := 1; id <= 5; id++ {
			Nums = append(Nums, HF_Atoi(configValue[fmt.Sprintf("num%d", id)]))
		}
		itembagnum_min := HF_Atoi(configValue["itembagnum_min"])
		itembagnum_max := HF_Atoi(configValue["itembagnum_max"])
		itembaggroup := HF_Atoi(configValue["itembaggroup"])
		var Itemshows []int
		for id := 1; id <= 4; id++ {
			Itemshows = append(Itemshows, HF_Atoi(configValue[fmt.Sprintf("item_show%d", id)]))
		}
		var Descriptions []int
		for id := 1; id <= 1; id++ {
			Descriptions = append(Descriptions, HF_Atoi(configValue[fmt.Sprintf("description%d", id)]))
		}

		var Completionitems []int
		for id := 1; id <= 2; id++ {
			Completionitems = append(Completionitems, HF_Atoi(configValue[fmt.Sprintf("completion_item%d", id)]))
		}
		var Completionnums []int
		for id := 1; id <= 2; id++ {
			Completionnums = append(Completionnums, HF_Atoi(configValue[fmt.Sprintf("completion_num%d", id)]))
		}

		self.GemsweepereventConfig = append(self.GemsweepereventConfig, &GemsweepereventConfig{
			Group:           group,
			Cycle:           cycle,
			Step:            step,
			Next:            next,
			Level:           level,
			Event:           event,
			Items:           Items,
			Nums:            Nums,
			Itembagnummin:   itembagnum_min,
			Itembagnummax:   itembagnum_max,
			Itembaggroup:    itembaggroup,
			Itemshows:       Itemshows,
			Completionitems: Completionitems,
			Completionnums:  Completionnums,
		})

		cycleValue, ok := self.GemGroupCycle[group]
		if !ok {
			self.GemGroupCycle[group] = cycle
		} else {
			if cycleValue < cycle {
				self.GemGroupCycle[group] = cycle
			}
		}

		stepValue, ok := self.GemGroupStep[group]
		if !ok {
			self.GemGroupStep[group] = step
		} else {
			if stepValue < step {
				self.GemGroupStep[group] = step
			}
		}
	}
}

// 获取宝藏消耗
func (self *CsvMgr) getTreasureCost(group int, cycle int, step int) *GemsweepereventConfig {
	for index := range self.GemsweepereventConfig {
		pConfig := self.GemsweepereventConfig[index]
		if pConfig == nil {
			continue
		}

		if pConfig.Group == group && pConfig.Cycle == cycle && pConfig.Step == step {
			return pConfig
		}
	}
	return nil
}

//获得英雄初始星级1
func (self *CsvMgr) GetHeroInitLv(heroId int) int {
	item := heroId*100 + 11000001
	config := self.ItemMap[item]
	if config != nil {
		return config.Special
	}
	return 0
}

func (self *CsvMgr) LoadCsv() {
	_ = litter.Sdump()
	self.LoadHeroStar()
	self.LoadSimpleConfig()
	self.LoadSkill()
	self.LoadTalent()
	self.LoadTalentAwakeConfig()
	self.LoadDreamLandConfig()
	self.LoadDreamLandSpendConfig()
	self.LoadActivityFundConfig()
	self.LoadFate()
	self.LoadHero()
	self.LoadTeamDungeon()
	self.LoadTimeRest()
	self.LoadAstrologyDrop()
	self.LoadSyntheticDrop()
	self.LoadEquip()
	self.LoadActifactEquip()
	self.LoadExclusiveEquip()
	self.LoadItemDrop()
	self.LoadHeroDecompose()
	self.LoadLevel()
	self.LoadTech()
	self.LoadNewUser()
	self.LoadArena()
	self.LoadTariff()
	self.LoadOther()
	//self.LoadBoss()
	self.LoadGemStone()
	self.LoadHorse()
	self.LoadSmeltTask()
	self.LoadMercenary()
	self.LoadMiliTask()
	self.LoadWorldPower()
	self.LoadGrouwthTasking()
	self.LoadHeroAttribute()
	self.LoadCrownBuild()
	GetLootMgr().LoadLottery()
	self.LoadPass()
	self.LoadWar()
	GetEventsMgr().LoadConfig()
	self.LoadOpenLevel()
	self.LoadEncourage()
	self.LoadHead()
	self.LoadBecomeStronger()

	self.LoadWorldLevel()
	self.LoadWorldMap()
	self.LoadMongyTask()
	self.LoadStr()

	//! 今日玩法
	self.LoadStatistics()

	//! 活动相关载入
	self.LoadActivityCsv()
	self.LoadNobilityTask()
	self.LoadWholeShop()
	self.LoadTurnTable()
	self.LoadTotleAward()
	self.LoadRecharge()
	self.LoadActivityGift()
	self.LoadGrowthGift()

	//英雄升级相关
	self.LoadHeroExp()

	//符文相关
	self.LoadRuneConfig()
	//阵型
	self.LoadFormationConfig()
	//超值基金
	self.LoadFundConfig()
	//新商店
	self.LoadNewShopConfig()
	//神兽
	self.LoadHydra()
	//地牢
	//self.LoadPit()
	//新地牢 异界迷宫
	self.LoadNewPit()
	//时光之巅
	self.LoadInstanceConfig()
	// 羁绊
	self.LoadEntanglement()
	// 悬赏
	self.LoadRewardConfig()
	//赏金令
	self.LoadWarOrder()
	//收藏家
	self.LoadAccessCard()
	// 凯旋丰碑
	self.LoadRankListIntegralConfig()
	self.LoadRankTask()
	self.LoadResonanceCrystalconfig()
	self.LoadArenaRewardConfig()
	self.LoadHangUp()
	self.LoadUnionHunt()
	self.LoadArenaSpecialClass()
	self.LoadHeroSkinConfig()
	self.LoadCrossArenaConfig()
	self.LoadSpecialPurchase()
	self.LoadTreeLevelConfig()
	self.LoadTreeProfessionalConfig()
	self.LoadInterstellarConfig()
	self.LoadActivityBossConfig()
	self.LoadStageTalentConfig()
	self.LoadRankReward()
}

func (self *CsvMgr) LoadActivityCsv() {
	GetCsvUtilMgr().LoadCsv("Activity_LimitGift", &self.ActivityTimeGiftMap)

	for _, item := range self.ActivityTimeGiftMap {
		_, ok := self.ActivityTimeGiftGroup[item.Group]
		if !ok {
			self.ActivityTimeGiftGroup[item.Group] = make([]*TimeGiftConfig, 0)
		}

		self.ActivityTimeGiftGroup[item.Group] = append(self.ActivityTimeGiftGroup[item.Group], item)
	}

	//litter.Dump(self.ActivityTimeGiftMap)
}

func (self *CsvMgr) LoadActivityGift() {
	GetCsvUtilMgr().LoadCsv("Activity_CurrencyGift", &self.ActivityGiftConfig)
}

func (self *CsvMgr) GetActivityGiftConfig(id int) *ActivityGiftConfig {
	for _, v := range self.ActivityGiftConfig {
		activityid := v.ActivityType*100000 + v.Group*100 + v.Index
		if activityid == id {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetActivityGiftConfigRefresh(id int) int {
	for _, v := range self.ActivityGiftConfig {
		if v.Id == id {
			return v.RefreshTime
		}
	}
	return 0
}

func (self *CsvMgr) LoadGrowthGift() {
	GetCsvUtilMgr().LoadCsv("Activity_GrowGift", &self.GrowthGiftConfig)
}

func (self *CsvMgr) GetGrowthGiftConfig(id int) *GrowthGiftConfig {
	for _, v := range self.GrowthGiftConfig {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) LoadStr() {
	GetCsvUtilMgr().LoadCsv("Str_String_Srv", &self.StrConfig)
	self.StrMap = make(map[string]string)

	for _, v := range self.StrConfig {
		self.StrMap[v.Dec] = v.Str
	}
	CAMP_NAME = append(CAMP_NAME, self.StrMap["STR_CAMP_1"])
	CAMP_NAME = append(CAMP_NAME, self.StrMap["STR_CAMP_2"])
	CAMP_NAME = append(CAMP_NAME, self.StrMap["STR_CAMP_3"])

	PART_NAME = append(PART_NAME, self.StrMap["STR_CAMP_HOUSE"])
	PART_NAME = append(PART_NAME, self.StrMap["STR_SOUTH_MISSION"])
	PART_NAME = append(PART_NAME, self.StrMap["STR_ARMY_WEAPON"])
	PART_NAME = append(PART_NAME, self.StrMap["STR_EAST_MISSION"])
	PART_NAME = append(PART_NAME, self.StrMap["STR_WOOD_HOUSE"])

	//litter.Dump(self.StrMap)
}

func (self *CsvMgr) GetText(s string) string {
	info, ok := self.StrMap[s]
	if !ok {
		return s
	}
	return info
}

// 加载新版王国任务
func (self *CsvMgr) LoadMongyTask() {
	self.MoneyTaskMap = make(map[int]*MongytaskListConfig)
	GetCsvUtilMgr().LoadCsv("Mongytask_List", &self.MoneyTaskMap)
	GetCsvUtilMgr().LoadCsv("Mongytask_Starlist", &self.MoneyTaskStarList)
	self.InitTaskTotal()
	//litter.Dump(self.MoneyTaskMap)
	//litter.Dump(self.MoneyTaskStarList)
}

func (self *CsvMgr) InitTaskTotal() {
	self.MoneyTotal = 0
	for _, v := range self.MoneyTaskStarList {
		self.MoneyTotal += v.Value
	}
	//litter.Dump(self.MoneyTotal)
}

// 加载头像配置
func (self *CsvMgr) LoadHead() {
	self.HeadConfigMap = make(map[int]*HeadConfig)
	GetCsvUtilMgr().LoadCsv("Headportrait", &self.HeadConfigMap)
	//litter.Dump(self.HeadConfigMap)
	//self.LoadHeadItem()
}

// 加载头像配置
func (self *CsvMgr) LoadWorldLevel() {
	self.WorldLevel = make(map[int]*WorldLevelConfig)
	GetCsvUtilMgr().LoadCsv("World_Level", &self.WorldLevel)
	//litter.Dump(self.WorldLevel)
}

func (self *CsvMgr) GetWorldFight() int {
	worldLv := GetServer().GetLevel(false)
	config, ok := self.WorldLevel[worldLv]
	if !ok {
		LogError("GetWorldFight worldLv not exists:", worldLv)
		return 0
	}

	fight := 0
	if len(config.BaseTypes) != len(config.BaseValues) {
		LogError("GetWorldFight len(config.BaseTypes) != len(config.BaseValues)")
		return 0
	}

	for i, v := range config.BaseTypes {
		if v == 99 {
			fight = config.BaseValues[i]
			break
		}
	}

	return fight
}

//func (self *CsvMgr) LoadHeadItem() {
//	self.HeadItemMap = make(map[int]int)
//	for _, v := range self.HeadConfigMap {
//		if v.Open != 2 {
//			continue
//		}
//		self.HeadItemMap[v.Condition] = v.Id
//	}
//}

//! 加载国战配置
func (self *CsvMgr) LoadWorldMap() {
	self.WorldMap = make(map[int]*WorldMapConfig)
	GetCsvUtilMgr().LoadCsv("World_Map", &self.WorldMap)

	self.ReadWarCity()
	//litter.Dump(self.WorldMap)
	self.StateOwner = make(map[int][]int)
	self.StateName = make(map[int]string)
	self.StateBox = make(map[int]int)
	for _, v := range self.WorldMap {
		if v.State == 0 {
			continue
		}
		self.StateOwner[v.State] = append(self.StateOwner[v.State], v.Id)
		self.StateName[v.State] = v.StateName
		self.StateBox[v.State] = v.StateBox
	}
	//litter.Dump(self.StateBox)

}

func (self *CsvMgr) GetStateName(stateId int) string {
	stateName, ok := self.StateName[stateId]
	if ok {
		return stateName
	}
	return "unkown"
}

// 加载鼓舞
func (self *CsvMgr) LoadEncourage() {
	self.WarEncourageConfig = make(map[int]*WarEncourageConfig)
	GetCsvUtilMgr().LoadCsv("War_Encourage", &self.WarEncourageConfig)
	self.InitEncourage()
	//litter.Dump(self.WarEncourageConfig)
	GetCsvUtilMgr().LoadCsv("War_Buff", &self.WarbuffConfig)
	//litter.Dump(self.WarbuffConfig)
}

func (self *CsvMgr) InitEncourage() {
	self.EncourageMap = make(map[int]map[int]*WarEncourageConfig)
	for _, v := range self.WarEncourageConfig {
		mapEn, ok := self.EncourageMap[v.Type]
		if !ok {
			mapEn = make(map[int]*WarEncourageConfig)
			mapEn[v.Lv] = v
			self.EncourageMap[v.Type] = mapEn
		} else {
			mapEn[v.Lv] = v
		}
	}
}

func (self *CsvMgr) GetEnLevelConfig(playType int, level int) *WarEncourageConfig {
	mapEn, ok := self.EncourageMap[playType]
	if !ok {
		return nil
	}

	config, ok := mapEn[level]
	if !ok {
		return config
	}

	return config

}

func (self *CsvMgr) LoadOpenLevel() {
	self.OpenLevelMap = make(map[int]*OpenLevelConfig)
	GetCsvUtilMgr().LoadCsv("Open_Level", &self.OpenLevelMap)
	//litter.Dump(self.OpenLevelMap)
}

func (self *CsvMgr) IsLevelOpen(curLv int, id int) (bool, int) {
	config, ok := self.OpenLevelMap[id]
	if !ok {
		return false, 0
	}
	return curLv >= config.Level, config.Level
}

func (self *CsvMgr) IsLevelOpen2(player *Player, id int) bool {
	config, ok := self.OpenLevelMap[id]
	if !ok {
		return false
	}

	if config.Passid == 0 {
		return true
	}

	passtitle := player.GetModule("pass").(*ModPass).GetPass(config.Passid)
	if passtitle != nil && passtitle.Num > 0 {
		return true
	}
	return false
}

func (self *CsvMgr) IsLevelOpenNew(player *Player, id int) bool {
	config, ok := self.OpenLevelMap[id]
	if !ok {
		return false
	}

	lvFlag := false
	if player.Sql_UserBase.Level >= config.Level {
		lvFlag = true
	}

	passFlag := false
	stage, _ := GetOfflineInfoMgr().GetBaseInfo(player.Sql_UserBase.Uid)
	if stage >= config.Passid {
		passFlag = true
	}

	return lvFlag && passFlag
}

func (self *CsvMgr) GuildIsLevelAndPassOpen(level int, stage int, id int) bool {
	config, ok := self.OpenLevelMap[id]
	if !ok {
		return false
	}

	lvFlag := false
	if level >= config.Level {
		lvFlag = true
	}

	passFlag := false
	//关卡修正，计算实际通关
	tempStage := stage - 1
	configTemp := self.LevelConfigMap[tempStage]
	if configTemp != nil {
		stage = stage - 1
	}
	if stage >= config.Passid {
		passFlag = true
	}

	return lvFlag && passFlag
}

func (self *CsvMgr) IsLevelAndPassOpenNew(level int, passId int, id int) bool {
	config, ok := self.OpenLevelMap[id]
	if !ok {
		return false
	}

	lvFlag := false
	if level >= config.Level {
		lvFlag = true
	}

	//关卡修正，计算实际通关
	passFlag := false
	if passId >= config.Passid {
		passFlag = true
	}

	return lvFlag && passFlag
}

type PlayNoticeInfo struct {
	StartHour, StartMin int
	EndHour, EndMin     int
	Interval            int64
	Notice              string
}

func (self *CsvMgr) LoadWar() {
	self.PlayTimeMap = make(map[int]*PlayTimeConfig)
	GetCsvUtilMgr().LoadCsv("Play_Time", &self.PlayTimeMap)
	GetCsvUtilMgr().LoadCsv("Play_Reward", &self.PlayRewardList)
	//pInfo := self.GetPlayRewar2(3, 2)
	//fmt.Println(*pInfo)

	self.checkPlayNoticeConfig()
	//litter.Dump(self.PlayTimeMap)
	//litter.Dump(self.PlayRewardList)
	//litter.Dump(self.LevelmapBoxConfig)
	//litter.Dump(self.LevelmapMailConfig)
	self.MailConfig = make(map[int]*MailConfig)
	GetCsvUtilMgr().LoadCsv("Mail", &self.MailConfig)
	//litter.Dump(self.MailConfig)
}

// 根据类型和排名获得奖励
func (self *CsvMgr) GetPlayReward(playType int, rank int) *PlayRewardConfig {
	for _, v := range self.PlayRewardList {
		if v.Type == playType && rank >= v.Minorder && rank <= v.Maxorder {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetPlayRewar2(play int, typeInfo int) *PlayRewardConfig {
	for _, v := range self.PlayRewardList {
		if v.Play == play && typeInfo == v.Type {
			return v
		}
	}
	return nil
}

//! 读取国战系统公告配置
func (self *CsvMgr) checkPlayNoticeConfig() {
	//! 矿点争夺
	self.PlayNotice = self.GetPlayNotice(CAMP_PVP_MINE)

	//! 孤山夺宝
	self.GveNotice = self.GetPlayNotice(CAMP_PVP_GVE)

	//!军团公告
	self.UnionFightNotice = self.GetPlayNotice(CAMP_PVP_UNION)

	//! 国战公告
	self.CampNotice = self.GetPlayNotice(CAMP_PVP_FIGHT)
	//litter.Dump(self.GveNotice)
	//litter.Dump(self.CampNotice)
}

func (self *CsvMgr) GetPlayNotice(id int) map[int]*PlayNoticeInfo {
	res := make(map[int]*PlayNoticeInfo)
	config, ok := GetCsvMgr().PlayTimeMap[id]
	if !ok {
		LogError("playTime config error, id=2 not exists!")
		os.Exit(1)
		return res
	}

	if len(config.NoticeStart) != len(config.NoticeEnd) &&
		len(config.NoticeStart) != len(config.Notice) && len(config.Notice) != 2 {
		LogError("playTime config error, id=2 size error")
		os.Exit(1)
		return res
	}

	for i := range config.NoticeStart {
		hourStar := config.NoticeStart[i] / 10000
		minStart := (config.NoticeStart[i] - hourStar*10000) / 100
		hourEnd := config.NoticeEnd[i] / 10000
		minEnd := (config.NoticeEnd[i] - hourEnd*10000) / 100
		res[i] = &PlayNoticeInfo{
			StartHour: hourStar,
			StartMin:  minStart,
			EndHour:   hourEnd,
			EndMin:    minEnd,
			Interval:  int64(config.Interval[i]),
			Notice:    config.Notice[i],
		}
	}

	config.BattleStartHour = config.Battlestart / 10000
	config.BattleStartMin = (config.Battlestart - config.BattleStartHour*10000) / 100

	config.BattleEndHour = config.Battleend / 10000
	config.BattleEndMin = (config.Battleend - config.BattleEndHour*10000) / 100

	config.DeclareStartHour = config.Declarestart / 10000
	config.DeclareStartMin = (config.Declarestart - config.DeclareStartHour*10000) / 100

	config.DeclareEndHour = config.Declareend / 10000
	config.DeclareEndMin = (config.Declareend - config.DeclareEndHour*10000) / 100

	config.EnrollStartHour = config.Enrollstart / 10000
	config.EnrollStartMin = (config.Enrollstart - config.EnrollStartHour*10000) / 100

	config.EnrollEndHour = config.Enrollend / 10000
	config.EnrollEndMin = (config.Enrollend - config.EnrollEndHour*10000) / 100

	// 战斗时间
	for i := 0; i < 3 && i < len(config.Battles); i++ {
		hour := config.Battles[i] / 10000
		minutes := (config.Battles[i] - hour*10000) / 100
		config.BattleInfo[i] = HourMinutes{hour, minutes}
	}

	//litter.Dump(config)
	return res
}

func (self *CsvMgr) LoadPass() {
	self.LevelItemMap = make(map[int]*LevelItemConfig)
	GetCsvUtilMgr().LoadCsv("Level_Item", &self.LevelItemMap)
}

func (self *CsvMgr) LoadCrownBuild() {
	self.CrownBuildConfig = make(map[int]*CrownBuildConfig)
	GetCsvUtilMgr().LoadCsv("Crown_Build", &self.CrownBuildConfig)
	//litter.Dump(self.CrownBuildConfig)
}

func (self *CsvMgr) LoadHeroAttribute() {
	self.HeroAttribute = make(map[int]*HeroAttribute)
	GetCsvUtilMgr().LoadCsv("Hero_Attribute", &self.HeroAttribute)
}

// 加载王国任务
func (self *CsvMgr) LoadGrouwthTasking() {
	GetCsvUtilMgr().LoadCsv("Growthtask_King", &self.TaskKingConfig)
	//litter.Dump(self.TaskKingConfig)
	self.TaskKingConfigMap = make(map[int]*GrowthtaskKingConfig)
	self.TaskKingGroupMap = make(map[int][]int)
	for _, v := range self.TaskKingConfig {
		self.TaskKingConfigMap[v.Taskid] = v
		found := false
		for _, g := range self.TaskKingGroupMap[v.Type] {
			if g == v.Group {
				found = true
				break
			}
		}
		if !found {
			self.TaskKingGroupMap[v.Type] = append(self.TaskKingGroupMap[v.Type], v.Group)
		}
	}
	//litter.Dump(self.TaskKingGroupMap)
}

func (self *CsvMgr) GetKingTask(group int) *GrowthtaskKingConfig {
	for _, v := range self.TaskKingConfig {
		if v.Group == group && v.Sort == 1 {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) LoadWorldPower() {
	self.WorldPowerMap = make(map[int]*WorldPower)
	GetCsvUtilMgr().LoadCsv("World_Power", &self.WorldPowerMap)
}

func (self *HeroStar) getAttribute(i int) []*Attribute {
	var attrTypes []int
	var attrValues []int64
	if i == 0 {
		attrTypes = self.Attr1
		attrValues = self.Value1
	} else if i == 1 {
		attrTypes = self.Attr2
		attrValues = self.Value2
	} else if i == 2 {
		attrTypes = self.Attr3
		attrValues = self.Value3
	} else if i == 3 {
		attrTypes = self.Attr4
		attrValues = self.Value4
	} else if i == 4 {
		attrTypes = self.Attr5
		attrValues = self.Value5
	} else if i == 5 {
		attrTypes = self.Attr6
		attrValues = self.Value6
	} else if i == 6 {
		attrTypes = self.Attr7
		attrValues = self.Value7
	} else if i == 7 {
		attrTypes = self.Attr8
		attrValues = self.Value8
	} else if i == 8 {
		attrTypes = self.Attr9
		attrValues = self.Value9
	} else if i == 9 {
		attrTypes = self.Attr10
		attrValues = self.Value10
	}

	if len(attrTypes) != len(attrValues) {
		LogError("len(attrTypes) != len(attrValues)")
		return []*Attribute{}
	}

	var attr []*Attribute
	for index := range attrTypes {
		attr = append(attr, &Attribute{attrTypes[index], attrValues[index]})
	}

	return attr
}

func (self *HeroStar) getStarAttribute() []*Attribute {
	var attrTypes = self.StarLvTypes
	var attrValues = self.StarLvValues

	if len(attrTypes) != len(attrValues) {
		LogError("len(attrTypes) != len(attrValues)")
		return []*Attribute{}
	}

	var attr []*Attribute
	for index := range attrTypes {
		//这里不好处理int64 暂时搁置 20190506 by zy
		attr = append(attr, &Attribute{attrTypes[index], attrValues[index]})
	}

	return attr
}

// 加载升星
func (self *CsvMgr) LoadHeroStar() {
	GetCsvUtilMgr().LoadCsv("Hero_Star", &self.HeroStar)
	self.HeroStarMap = make(map[int]map[int]*HeroStar)
	self.HeroStarAttrMap = make(map[int]map[int]*HeroStarAttr)
	for _, pInfo := range self.HeroStar {
		starMap, ok := self.HeroStarMap[pInfo.HeroId]
		attrMap, attrOk := self.HeroStarAttrMap[pInfo.HeroId]
		// star map
		if !ok {
			starMap = make(map[int]*HeroStar)
			starMap[pInfo.Star] = pInfo
			self.HeroStarMap[pInfo.HeroId] = starMap
		} else {
			starMap[pInfo.Star] = pInfo
		}

		// hero attr
		if !attrOk {
			attrMap = make(map[int]*HeroStarAttr)
			pAttr := &HeroStarAttr{}
			attrMap[pInfo.Star] = pAttr
			for i := 0; i < maxStarSlots; i++ {
				var attr = pInfo.getAttribute(i)
				pAttr.SlotAttr[i] = append(pAttr.SlotAttr[i], attr...)
			}
			pAttr.StarAttr = pInfo.getStarAttribute()
			self.HeroStarAttrMap[pInfo.HeroId] = attrMap
		} else {
			pAttr := &HeroStarAttr{}
			attrMap[pInfo.Star] = pAttr
			for i := 0; i < maxStarSlots; i++ {
				var attr = pInfo.getAttribute(i)
				pAttr.SlotAttr[i] = append(pAttr.SlotAttr[i], attr...)
			}
			pAttr.StarAttr = pInfo.getStarAttribute()
		}

	}
	//litter.Dump(self.HeroStar)
}

func (self *CsvMgr) getHeroStarSkill(heroId int, star int) []int {
	skillMap, skillOk := self.HeroStarMap[heroId]
	if !skillOk {
		return []int{}
	}

	config, skillsOk := skillMap[star]
	if !skillsOk {
		return []int{}
	}

	return config.SkillIds
}

func (self *CsvMgr) LoadSimpleConfig() {
	self.SimpleConfigMap = make(map[int]*SimpleConfig)
	GetCsvUtilMgr().LoadCsv("Simple_Config", &self.SimpleConfigMap)
	//for _, v := range self.SimpleConfigMap {
	//	fmt.Println(*v)
	//}
}

// 获取资质所加得属性
func (self *CsvMgr) getQualityAtt(heroId int, star int, slot int) []*Attribute {
	attrMap, attrMapOk := self.HeroStarAttrMap[heroId]
	if !attrMapOk {
		return []*Attribute{}
	}
	pAttr, pAttrOk := attrMap[star]
	if !pAttrOk {
		return []*Attribute{}
	}

	if slot < 0 || slot >= 6 {
		return []*Attribute{}
	}

	return pAttr.SlotAttr[slot]
}

func (self *CsvMgr) LoadSkill() {
	self.SkillConfigMap = make(map[int]*HeroSkill)
	GetCsvUtilMgr().LoadCsv("Hero_Skill", &self.SkillConfigMap)
}

// 加载天赋配置
func (self *CsvMgr) LoadTalent() {
	GetCsvUtilMgr().LoadCsv("Divinity_Config", &self.TalentConfig)

	// 天赋配置
	self.TalentMap = make(map[int]map[int]*TalentConfig)
	// 最大天赋等级
	self.MaxTalentLv = make(map[int]int)
	for k, v := range self.TalentConfig {
		talent, ok := self.TalentMap[v.TalentId]

		nID := v.TalentId
		nLevel := v.TalentLv
		if !ok {
			talent = make(map[int]*TalentConfig)
			talent[nLevel] = self.TalentConfig[k]
			self.TalentMap[nID] = talent
		} else {
			talent[nLevel] = self.TalentConfig[k]
		}

		if self.MaxTalentLv[nID] < nLevel {
			self.MaxTalentLv[nID] = nLevel
		}
	}

	//litter.Dump(self.TalentConfig)

}

// 读取觉醒配置
func (self *CsvMgr) LoadTalentAwakeConfig() {
	GetCsvUtilMgr().LoadCsv("Divinity_Awaken", &self.TalentAwakeConfig)

	self.TalentAwakeMap = make(map[int][]*TalentAwake)

	for _, v := range self.TalentAwakeConfig {
		nGroup := v.Group
		_, ok := self.TalentAwakeMap[nGroup]
		if !ok {
			self.TalentAwakeMap[nGroup] = make([]*TalentAwake, 0)
			self.TalentAwakeMap[nGroup] = append(self.TalentAwakeMap[nGroup], v)
		} else {
			self.TalentAwakeMap[nGroup] = append(self.TalentAwakeMap[nGroup], v)
		}
	}
}

// 读取幻境配置
func (self *CsvMgr) LoadDreamLandConfig() {
	GetCsvUtilMgr().LoadCsv("Divinity_Dreamland", &self.DreamLandConfig)

	self.DreamLandGroupMap = make(map[int][]*DreamLand)
	self.DreamLandItemMap = make(map[int]*DreamLand)

	for _, v := range self.DreamLandConfig {
		nGroup := v.Group
		_, ok1 := self.DreamLandGroupMap[nGroup]
		if !ok1 {
			self.DreamLandGroupMap[nGroup] = make([]*DreamLand, 0)
			self.DreamLandGroupMap[nGroup] = append(self.DreamLandGroupMap[nGroup], v)
		} else {
			self.DreamLandGroupMap[nGroup] = append(self.DreamLandGroupMap[nGroup], v)
		}

		_, ok2 := self.DreamLandItemMap[v.ID]
		if !ok2 {
			self.DreamLandItemMap[v.ID] = v
		}
	}
}

// 读取幻境消耗配置
func (self *CsvMgr) LoadDreamLandSpendConfig() {
	GetCsvUtilMgr().LoadCsv("Divinity_Refresh", &self.DreamLandSpendConfig)

	self.DreamLandSpendMap = make(map[int]map[int]*DreamLandSpend)
	self.DreamLandCostMap = make(map[int]*DreamLandCost)

	for _, v := range self.DreamLandSpendConfig {
		nType := v.Type
		nClass := v.Class
		_, ok1 := self.DreamLandSpendMap[nType]
		if !ok1 {
			self.DreamLandSpendMap[nType] = make(map[int]*DreamLandSpend)
			self.DreamLandSpendMap[nType][nClass] = v
		} else {
			self.DreamLandSpendMap[nType][nClass] = v
		}

		_, ok2 := self.DreamLandCostMap[nType]
		if !ok2 {
			self.DreamLandCostMap[nType] = &DreamLandCost{v.Type, v.RefCost, v.LootCost, v.TypeTimes}
		}
	}
}

// 读取活动基金配置
func (self *CsvMgr) LoadActivityFundConfig() {
	GetCsvUtilMgr().LoadCsv("Activity_Fund", &self.ActivityFundConfig)

	self.ActivityFundMap = make(map[int]*ActivityFundGroupMap)
	//self.ActivityFundTypeMap =  make(map[int]*ActivityFundTypeConfig)

	for _, v := range self.ActivityFundConfig {
		nGroupID := v.GroupID
		nPay := v.Pay
		nDay := v.Day
		nType := v.Type
		nWroth := v.Worth
		group, ok1 := self.ActivityFundMap[nGroupID]
		if !ok1 {
			groupConfig := &ActivityFundGroupMap{nGroupID, make(map[int]*ActivityFundPayMap)}
			self.ActivityFundMap[nGroupID] = groupConfig

			payConfig := &ActivityFundPayMap{nPay, nType, nWroth, make(map[int]*ActivityFundConfig)}
			self.ActivityFundMap[nGroupID].PayConfig[nPay] = payConfig

			self.ActivityFundMap[nGroupID].PayConfig[nPay].DayConfig[nDay] = v

		} else {
			pay, ok2 := group.PayConfig[nPay]
			if !ok2 {
				payConfig := &ActivityFundPayMap{nPay, nType, nWroth, make(map[int]*ActivityFundConfig)}
				self.ActivityFundMap[nGroupID].PayConfig[nPay] = payConfig

				self.ActivityFundMap[nGroupID].PayConfig[nPay].DayConfig[nDay] = v
			} else {
				_, ok3 := pay.DayConfig[nDay]
				if !ok3 {
					self.ActivityFundMap[nGroupID].PayConfig[nPay].DayConfig[nDay] = v
				}
			}
		}
	}
}

func (self *CsvMgr) GetTalentConfig(talentId, talentLv int) *TalentConfig {
	talentMapId, ok := self.TalentMap[talentId]
	if !ok {
		return nil
	}

	talent, ok := talentMapId[talentLv]
	if !ok {
		return nil
	}

	return talent

}

// 加载我要变强配置
func (self *CsvMgr) LoadBecomeStronger() {
	GetCsvUtilMgr().LoadCsv("Become_Stronger", &self.BecomeStrongerConfig)

	// 只读取有效配置
	self.BecomeStrongerMap = make(map[int]*BecomeStronger)

	for _, v := range self.BecomeStrongerConfig {
		if v.ConditionType > 0 && v.Boxid > 0 {
			_, ok := self.BecomeStrongerMap[v.Id]
			if !ok {
				self.BecomeStrongerMap[v.Id] = v
			}
		}
	}
}

func (self *CsvMgr) LoadFate() {
	GetCsvUtilMgr().LoadCsv("Hero_Fate", &self.FateConfig)
	self.FateMap = make(map[int][]*FateConfig)
	for _, v := range self.FateConfig {
		_, ok := self.FateMap[v.HeroId]
		if !ok {
			self.FateMap[v.HeroId] = make([]*FateConfig, 0)
			self.FateMap[v.HeroId] = append(self.FateMap[v.HeroId], v)
		} else {
			self.FateMap[v.HeroId] = append(self.FateMap[v.HeroId], v)
		}
	}
	//litter.Dump(self.FateConfig)
}

func (self *CsvMgr) getFateConfig(heroId int, id int) *FateConfig {
	configs, ok := GetCsvMgr().FateMap[heroId]
	if !ok {
		return nil
	}

	var config *FateConfig
	for _, v := range configs {
		if v.FateId == id {
			config = v
			break
		}
	}
	return config
}

// 获得对应资质的消耗
func (self *HeroStar) getCost(index int) ([]int, []int) {
	if index == 1 {
		return self.SlotItemId1s, self.SlotItemNum1s
	} else if index == 2 {
		return self.SlotItemId2s, self.SlotItemNum2s
	} else if index == 3 {
		return self.SlotItemId3s, self.SlotItemNum3s
	} else if index == 4 {
		return self.SlotItemId4s, self.SlotItemNum4s
	} else if index == 5 {
		return self.SlotItemId5s, self.SlotItemNum5s
	} else if index == 6 {
		return self.SlotItemId6s, self.SlotItemNum6s
	} else if index == 7 {
		return self.SlotItemId7s, self.SlotItemNum7s
	} else if index == 8 {
		return self.SlotItemId8s, self.SlotItemNum8s
	} else if index == 9 {
		return self.SlotItemId9s, self.SlotItemNum9s
	} else if index == 10 {
		return self.SlotItemId10s, self.SlotItemNum10s
	}
	return []int{}, []int{}
}

// 加载英雄配置
func (self *CsvMgr) LoadHero() {
	GetCsvUtilMgr().LoadCsv("Hero_Config", &self.HeroConfig)
	self.HeroConfigMap = make(map[int]map[int]*HeroConfig)
	for _, v := range self.HeroConfig {
		config, ok := self.HeroConfigMap[v.HeroId]
		if ok {
			config[v.HeroStar] = v
		} else {
			config = make(map[int]*HeroConfig)
			config[v.HeroStar] = v
			self.HeroConfigMap[v.HeroId] = config
		}
	}

	GetCsvUtilMgr().LoadCsv("Hero_break", &self.HeroBreakConfig)
	self.HeroBreakConfigMap = make(map[int]map[int]*HeroBreakConfig)
	for _, v := range self.HeroBreakConfig {
		config, ok := self.HeroBreakConfigMap[v.HeroId]
		if ok {
			config[v.Id] = v
		} else {
			config = make(map[int]*HeroBreakConfig)
			config[v.Id] = v
			self.HeroBreakConfigMap[v.HeroId] = config
		}
	}

	self.HeroNumMap = make(map[int]*HeroNumConfig)
	GetCsvUtilMgr().LoadCsv("Hero_Num", &self.HeroNumMap)

	handBookTemp := make([]*HeroHandBookConfig, 0)
	GetCsvUtilMgr().LoadCsv("Hero_tujian", &handBookTemp)
	self.HeroHandBookConfigMap = make(map[int]*HeroHandBookConfig)
	for _, v := range handBookTemp {
		self.HeroHandBookConfigMap[v.Id] = v
	}

	GetCsvUtilMgr().LoadCsv("Hero_growth", &self.HeroGrowthConfig)
	self.HeroGrowthConfigMap = make(map[int]*HeroGrowthConfig)
	for _, v := range self.HeroGrowthConfig {
		self.HeroGrowthConfigMap[v.GrowthLevel] = v
	}

	GetCsvUtilMgr().LoadCsv("Lineup_Add", &self.TeamAttrConfig)
}

// 加载副本配置
func (self *CsvMgr) LoadTeamDungeon() {
	self.TeamDungeonMap = make(map[int]*TeamDungeonConfig)
	GetCsvUtilMgr().LoadCsv("Level_DungeonsTeam", &self.TeamDungeonMap)
	//litter.Dump(self.TeamDungeonMap)
}

func (self *CsvMgr) LoadTimeRest() {
	GetCsvUtilMgr().LoadCsv("Time_Reset", &self.TimeResetConfig)
}

func (self *CsvMgr) LoadAstrologyDrop() {

	AstrologyDropConfigTemp := make([]*AstrologyDropConfig, 0)
	GetCsvUtilMgr().LoadCsv("Astrology_Drop", &AstrologyDropConfigTemp)

	self.AstrologyDropConfig = make(map[int][]*AstrologyDropConfig, 0)
	self.AstrologyDropGroupConfig = make(map[int]*AstrologyDropConfig, 0)
	for _, v := range AstrologyDropConfigTemp {
		self.AstrologyDropConfig[v.AstrologyId] = append(self.AstrologyDropConfig[v.AstrologyId], v)
		if v.AstrologyChance > 0 {
			self.AstrologyDropGroupConfig[v.AstrologyId] = v
		}
	}

	//litter.Dump(self.AstrologyDropGroupConfig)
}

func (self *CsvMgr) LoadSyntheticDrop() {

	SyntheticDropConfigTemp := make([]*SyntheticDropConfig, 0)
	GetCsvUtilMgr().LoadCsv("synthetic_config", &SyntheticDropConfigTemp)

	self.SyntheticDropConfig = make(map[int][]*SyntheticDropConfig, 0)
	self.SyntheticDropGroupConfig = make(map[int]*SyntheticDropConfig, 0)
	for _, v := range SyntheticDropConfigTemp {
		self.SyntheticDropConfig[v.SyntheticId] = append(self.SyntheticDropConfig[v.SyntheticId], v)
		if v.SyntheticChance > 0 {
			self.SyntheticDropGroupConfig[v.SyntheticId] = v
		}
	}

	//litter.Dump(self.AstrologyDropGroupConfig)
}

// 加载装备配置
func (self *CsvMgr) LoadEquip() {
	self.ItemMap = make(map[int]*ItemConfig)
	GetCsvUtilMgr().LoadCsv("Itemconfig", &self.ItemMap)
	//litter.Dump(self.ItemMap)

	GetCsvUtilMgr().LoadCsv("Equip_Upgrade", &self.EquipUpgrade)
	//litter.Dump(self.EquipUpgrade)

	self.EquipUpgradeMap = make(map[int]map[int]*EquipUpgrade)
	self.EquipUpgradeMaxLv = make(map[int]int)
	for _, v := range self.EquipUpgrade {
		lvInfo, ok := self.EquipUpgradeMap[v.Id]
		if !ok {
			lvInfo = make(map[int]*EquipUpgrade)
			lvInfo[v.Lv] = v
			self.EquipUpgradeMap[v.Id] = lvInfo
		} else {
			lvInfo[v.Lv] = v
		}
		lv, ok2 := self.EquipUpgradeMaxLv[v.Id]
		if ok2 {
			self.EquipUpgradeMaxLv[v.Id] = v.Lv
		} else {
			if lv < v.Lv {
				self.EquipUpgradeMaxLv[v.Id] = v.Lv
			}
		}
	}

	//litter.Dump(self.EquipUpgradeMap)
	//litter.Dump(self.EquipUpgradeMaxLv)

	GetCsvUtilMgr().LoadCsv("Equip_Star", &self.EquipStar)
	//litter.Dump(self.EquipStar)
	self.EquipStarMap = make(map[int]map[int]*EquipStar)
	self.EquipStarMaxLv = make(map[int]int)
	for _, v := range self.EquipStar {
		lvInfo, ok := self.EquipStarMap[v.Id]
		if !ok {
			lvInfo = make(map[int]*EquipStar)
			lvInfo[v.Lv] = v
			self.EquipStarMap[v.Id] = lvInfo
		} else {
			lvInfo[v.Lv] = v
		}
		lv, ok2 := self.EquipStarMaxLv[v.Id]
		if ok2 {
			self.EquipStarMaxLv[v.Id] = v.Lv
		} else {
			if lv < v.Lv {
				self.EquipStarMaxLv[v.Id] = v.Lv
			}
		}
	}

	// 计算重生道具
	for _, lvInfo := range self.EquipStarMap {
		for level, pConfig := range lvInfo {
			pConfig.RebornItem = make(map[int]*Item)
			//20190411 by zhangyang
			//原来的逻辑是3星物品 返回1升2，2升3，3升4星的累计消耗，现在改为返回0升1，1升2，2升3的累计消耗。
			//这个是计算的逻辑，不是配置的问题，客户端也需要修改重生消耗的计算逻辑
			//for i := 1; i <= level; i++ {
			for i := 0; i <= level-1; i++ {
				config, ok := lvInfo[i]
				if !ok {
					continue
				}
				AddItemMapHelper(pConfig.RebornItem, config.CostIds, config.CostNums)
			}
		}
	}

	//litter.Dump(self.EquipStarMap)

	//litter.Dump(self.EquipStarMap)
	//litter.Dump(self.EquipStarMaxLv)
	self.EquipGem = make(map[int]*EquipGem)
	GetCsvUtilMgr().LoadCsv("Equip_Gem", &self.EquipGem)
	//litter.Dump(self.EquipGem)
	self.GemLevelup = make(map[int]int)
	for _, v := range self.EquipGem {
		self.GemLevelup[v.NeedId] = v.Id
	}

	self.EquipSuit = make(map[int]*EquipSuit)
	GetCsvUtilMgr().LoadCsv("Equip_Suit", &self.EquipSuit)
	//litter.Dump(self.EquipSuit)

	self.EquipValueGroupMap = make(map[int]*EquipValueGroup)
	GetCsvUtilMgr().LoadCsv("Equip_valuegroup", &self.EquipValueGroupMap)

	self.EquipBaseValueMap = make(map[int]*EquipBaseValue)
	GetCsvUtilMgr().LoadCsv("Equip_basevalue", &self.EquipBaseValueMap)

	self.EquipSpecialValueMap = make(map[int]*EquipSpecialValue)
	GetCsvUtilMgr().LoadCsv("Equip_specialvalue", &self.EquipSpecialValueMap)

	/*
		EquipStrengthenMap     map[int]*EquipStrengthenConfig     //
		EquipStrengthenUpLvMap map[int]*EquipStrengthenLvUpConfig //
	*/
	//新魔龙
	self.EquipConfigMap = make(map[int]*EquipConfig)
	GetCsvUtilMgr().LoadCsv("Equip_Config", &self.EquipConfigMap)

	self.EquipShopRate = make(map[int]int)
	for _, v := range self.EquipConfigMap {
		if v.ShopClass > 0 && v.ShopWeight > 0 {
			self.EquipShopRate[v.ShopClass] += v.ShopWeight
		}
	}

	self.EquipAdvancedConfigMap = make(map[int]*EquipAdvancedConfig)
	GetCsvUtilMgr().LoadCsv("Equip_Advanced", &self.EquipAdvancedConfigMap)

	equipStrengthen := make([]*EquipStrengthenConfig, 0)
	GetCsvUtilMgr().LoadCsv("Equip_Strengthen", &equipStrengthen)
	self.EquipStrengthenMap = make(map[int][]*EquipStrengthenConfig)
	for _, v := range equipStrengthen {
		self.EquipStrengthenMap[v.Type] = append(self.EquipStrengthenMap[v.Type], v)
	}

	GetCsvUtilMgr().LoadCsv("Equip_StrengthenLvUp", &self.EquipStrengthenUpLvMap)

	GetCsvUtilMgr().LoadCsv("Mercenary_config", &self.EquipHireConfig)

	GetCsvUtilMgr().LoadCsv("Equip_Recast", &self.EquipRecastConfig)
}

// 加载神器配置
func (self *CsvMgr) LoadActifactEquip() {
	self.ArtifactEquipConfigMap = make(map[int]*ArtifactEquipConfig)
	GetCsvUtilMgr().LoadCsv("Artifact_Config", &self.ArtifactEquipConfigMap)

	GetCsvUtilMgr().LoadCsv("Artifact_Strengthen", &self.ArtifactStrengthen)
}

// 加载专属配置
func (self *CsvMgr) LoadExclusiveEquip() {
	self.ExclusiveEquipConfigMap = make(map[int]*ExclusiveEquipConfig)
	GetCsvUtilMgr().LoadCsv("Exclusiveequip", &self.ExclusiveEquipConfigMap)

	self.ExclusiveStrengthen = make(map[int]map[int]*ExclusiveStrengthenConfig)
	tempConfig := make([]*ExclusiveStrengthenConfig, 0)
	GetCsvUtilMgr().LoadCsv("Exclusiveequip_Strengthen", &tempConfig)
	for _, v := range tempConfig {
		_, ok := self.ExclusiveStrengthen[v.Id]
		if !ok {
			config := make(map[int]*ExclusiveStrengthenConfig)
			config[v.Lv] = v
			self.ExclusiveStrengthen[v.Id] = config
		} else {
			self.ExclusiveStrengthen[v.Id][v.Lv] = v
		}
	}
	//litter.Dump(self.ExclusiveStrengthen)
}

//! 找到道具配置
func (self *CsvMgr) GetItemConfig(itemId int) *ItemConfig {
	v, ok := self.ItemMap[itemId]
	if ok {
		return v
	}
	return nil
}

func (self *CsvMgr) GetEquipStrengthenLvUpConfig(equipId int, lv int) *EquipStrengthenLvUpConfig {

	config := self.EquipConfigMap[equipId]
	if config == nil {
		return nil
	}
	for _, v := range self.EquipStrengthenUpLvMap {
		if v.Quality == config.Quality && v.Lv == lv {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetArtifactStrengthenLvUpConfig(equipId int, lv int) *ArtifactStrengthenConfig {

	for _, v := range self.ArtifactStrengthen {
		if v.Id == equipId && v.Lv == lv {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetEquipStrengthenConfig(equipType int, equipPos int, equipQuality int, lv int) *EquipStrengthenConfig {

	config := self.EquipStrengthenMap[equipType]
	if config == nil {
		return nil
	}
	for _, v := range config {
		if v.EquipPosition == equipPos && v.Quality == equipQuality && v.Lv == lv {
			return v
		}
	}
	return nil
}

//! 找到装备配置
func (self *CsvMgr) GetEquipConfig(itemId int) *EquipConfig {
	v, ok := self.EquipConfigMap[itemId]
	if ok {
		return v
	}
	return nil
}

//! 找到强化配置
func (self *CsvMgr) GetEquipUpGrade(id int, lv int) *EquipUpgrade {
	data, ok := self.EquipUpgradeMap[id]
	if !ok {
		return nil
	}

	config, ok := data[lv]
	if !ok {
		return nil
	}
	return config
}

//! 找到附魔配置
func (self *CsvMgr) GetEquipStar(id int, lv int) *EquipStar {
	data, ok := self.EquipStarMap[id]
	if !ok {
		return nil
	}

	config, ok := data[lv]
	if !ok {
		return nil
	}
	return config
}

//! 找到宝石配置
func (self *CsvMgr) GetGemConfig(itemId int) *EquipGem {
	v, ok := self.EquipGem[itemId]
	if ok {
		return v
	}
	return nil
}

//! 找到套装属性
func (self *CsvMgr) GetEquipSuitAttr(suitId int, num int) map[int]*Attribute {
	res := make(map[int]*Attribute)
	config, ok := self.EquipSuit[suitId]
	if !ok {
		return res
	}

	index := -1
	for i, mark := range config.SuitMark {
		if num >= mark {
			index = i
		}
	}

	if index == -1 {
		return res
	}

	var baseTypes []int
	var baseValues []int64
	//属性累加  20190531 by zy
	if index >= 0 {
		baseTypes = config.BaseTypes1
		baseValues = config.BaseValues1
		AddAttrDirect(res, baseTypes, baseValues)
	}
	if index >= 1 {
		baseTypes = config.BaseTypes2
		baseValues = config.BaseValues2
		AddAttrDirect(res, baseTypes, baseValues)
	}

	if index >= 2 {
		baseTypes = config.BaseTypes3
		baseValues = config.BaseValues3
		AddAttrDirect(res, baseTypes, baseValues)
	}

	/*
		if index == 0 {
			baseTypes = config.BaseTypes1
			baseValues = config.BaseValues1
		} else if index == 1 {
			baseTypes = config.BaseTypes2
			baseValues = config.BaseValues2
		} else if index == 2 {
			baseTypes = config.BaseTypes3
			baseValues = config.BaseValues3
		}
		AddAttrDirect(res, baseTypes, baseValues)
	*/
	return res
}

//! 找到英雄配置
func (self *CsvMgr) GetHeroConfig(id int) *HeroConfig {
	//v, ok := self.HeroConfig[id]
	//if ok {
	//	return v
	//}
	return nil
}

//! 找到英雄配置
func (self *CsvMgr) GetHeroMapConfig(id int, star int) *HeroConfig {
	v, ok := self.HeroConfigMap[id]
	if ok {
		t, ok2 := v[star]
		if ok2 {
			return t
		}
	}
	return nil
}

// 加载掉落包配置
func (self *CsvMgr) LoadItemDrop() {
	self.ItemBagGroupMap = make(map[int]*ItemBagGroup)
	GetCsvUtilMgr().LoadCsv("Itembaggroup", &self.ItemBagGroupMap)
	for _, v := range self.ItemBagGroupMap {
		total := 0
		for index := range v.Weights {
			total += v.Weights[index]
		}

		if len(v.Weights) != len(v.DropIds) {
			LogError("len(v.Weights) != len(v.DropIds)")
		}
		v.Sum = total
	}
	///litter.Dump(self.ItemBagGroupMap)

	self.ItemBagMap = make(map[int]*ItemBag)
	GetCsvUtilMgr().LoadCsv("Itembag", &self.ItemBagMap)
	for _, v := range self.ItemBagMap {
		total := 0
		for index := range v.Weights {
			total += v.Weights[index]
		}
		if len(v.Weights) != len(v.ItemIds) || len(v.Weights) != len(v.Nums) {
			LogError("len(v.Weights) != len(v.ItemIds) || len(v.Weights) != len(v.Nums)")
		}
		v.Sum = total
	}
	//litter.Dump(self.ItemBagMap)
}

//! 从掉落包里选出一个
func (self *CsvMgr) DropItemBag(dropId int) (int, int) {
	dropConfig, ok := self.ItemBagMap[dropId]
	if !ok {
		return 0, 0
	}

	if dropConfig.Sum == 0 {
		return 0, 0
	}

	randNum := HF_GetRandom(dropConfig.Sum)
	check := 0
	for i := 0; i < len(dropConfig.ItemIds); i++ {
		itemId := dropConfig.ItemIds[i]
		itemNum := dropConfig.Nums[i]
		if itemId == 0 {
			continue
		}

		if itemNum == 0 {
			continue
		}

		check += dropConfig.Weights[i]
		if randNum < check {
			return itemId, itemNum
		}
	}

	return 0, 0
}

//! 从掉落组里选出若干
func (self *CsvMgr) DropItem(groupId int) []PassItem {
	outItem := make([]PassItem, 0)
	itemBagGrop, ok := self.ItemBagGroupMap[groupId]
	if !ok {
		return outItem
	}

	if itemBagGrop.Type == 1 {
		// 得到总权重
		randNum := HF_GetRandom(itemBagGrop.Sum)
		check := 0
		for i := 0; i < len(itemBagGrop.Weights); i++ {
			check += itemBagGrop.Weights[i]
			if randNum < check {
				// 随机一个掉落Id, 再从掉落Id中随机一个物品
				itemId, itemNum := self.DropItemBag(itemBagGrop.DropIds[i])
				if itemId != 0 && itemNum != 0 {
					outItem = append(outItem, PassItem{itemId, itemNum})
				}
				break
			}
		}
	} else if itemBagGrop.Type == 3 {
		for i := 0; i < len(itemBagGrop.DropIds); i++ {
			itemId := itemBagGrop.DropIds[i]
			itemNum := itemBagGrop.Weights[i]
			if itemId != 0 && itemNum != 0 {
				outItem = append(outItem, PassItem{itemId, itemNum})
			}
		}
	} else if itemBagGrop.Type == 2 {
		for i := 0; i < len(itemBagGrop.Weights); i++ {
			weight := itemBagGrop.Weights[i]
			dropId := itemBagGrop.DropIds[i]
			if dropId == 0 {
				continue
			}
			randNum := HF_GetRandom(10000)
			if randNum < weight {
				itemId, itemNum := self.DropItemBag(dropId)
				if itemId != 0 && itemNum != 0 {
					outItem = append(outItem, PassItem{itemId, itemNum})
				}
			}
		}
	}

	return outItem
}

// 英雄碎片分解
func (self *CsvMgr) LoadHeroDecompose() {
	self.HeroDecomposeMap = make(map[int]*HeroDecompose)
	GetCsvUtilMgr().LoadCsv("Hero_Decompose", &self.HeroDecomposeMap)
	//LogDebug(litter.Sdump(self.HeroDecomposeMap))
}

// 加载关卡配置
func (self *CsvMgr) LoadLevel() {
	self.LevelConfigMap = make(map[int]*LevelConfig)
	GetCsvUtilMgr().LoadCsv("Level_Config", &self.LevelConfigMap)
	//litter.Dump(self.LevelConfigMap)
	self.LoadMissionInfo()
}

func (self *CsvMgr) LoadMissionInfo() {
	self.MissionMap = make(map[int]*MissionInfo)
	for _, v := range self.LevelConfigMap {
		if v.MainType != 1 {
			continue
		}
		_, ok := self.MissionMap[v.LevelIndex]
		if !ok {
			self.MissionMap[v.LevelIndex] = &MissionInfo{}
			self.MissionMap[v.LevelIndex].Chapter = v.LevelIndex
			self.MissionMap[v.LevelIndex].MissionIds = append(self.MissionMap[v.LevelIndex].MissionIds, v.LevelId)
		} else {
			self.MissionMap[v.LevelIndex].MissionIds = append(self.MissionMap[v.LevelIndex].MissionIds, v.LevelId)
		}
	}

	//litter.Dump(self.MissionMap)
}

func (self *CsvMgr) getHeroDecompose(itemId int) *HeroDecompose {
	config, ok := self.HeroDecomposeMap[itemId]
	if !ok {
		return nil
	}
	return config
}

// 加载科技
func (self *CsvMgr) LoadTech() {
	self.CostMap = make(map[int]*CostConfig)
	GetCsvUtilMgr().LoadCsv("Cost", &self.CostMap)
	//litter.Dump(self.CostMap)
	self.TechConfigMap = make(map[int]*TechConfig)
	GetCsvUtilMgr().LoadCsv("Science", &self.TechConfigMap)
	//litter.Dump(self.TechConfigMap)
	GetCsvUtilMgr().LoadCsv("Science_Attribute", &self.TechAttr)
	//litter.Dump(self.TechAttr)
	GetCsvUtilMgr().LoadCsv("Science_QueueTime", &self.TechTime)
	//litter.Dump(self.TechTime)

	// 科技等级信息
	self.MaxTechLv = make(map[int]int)
	for _, v := range self.TechAttr {
		num, ok := self.MaxTechLv[v.Group]
		if !ok {
			self.MaxTechLv[v.Group] = v.Level
		} else {
			if num < v.Level {
				self.MaxTechLv[v.Group] = v.Level
			}
		}
	}

}

func (self *CsvMgr) GetCostConfig(costId int) *CostConfig {
	v, ok := self.CostMap[costId]
	if ok {
		return v
	}
	return nil
}

func (self *CsvMgr) GetBossConfig(id int) *BossConfig {
	v, ok := self.BossConfig[id]
	if ok {
		return v
	}
	return nil
}

func (self *CsvMgr) GetTechConfig(techId int) *TechConfig {
	v, ok := self.TechConfigMap[techId]
	if ok {
		return v
	}
	return nil
}

func (self *CsvMgr) GetTechAttrConfig(techId int, techLv int) *TechAttr {
	config := self.GetTechConfig(techId)
	if config == nil {
		return nil
	}
	group := config.TechGroup
	for _, v := range self.TechAttr {
		if v.Group == group && v.Level == techLv {
			return v
		}
	}

	return nil
}

func (self *CsvMgr) getTechTime() *TechTime {
	for _, v := range self.TechTime {
		return v
	}

	return nil
}

func (self *CsvMgr) GetItemLimit(itemId int) int {
	config := self.GetItemConfig(itemId)
	if config == nil {
		return 0
	}

	return config.MaxNum
}

func (self *CsvMgr) GetItemGemPrice(itemId int) int {
	config := self.GetItemConfig(itemId)
	if config == nil {
		return 0
	}

	return config.GemPrice
}

func (self *CsvMgr) GetItemType(itemId int) string {
	config := self.GetItemConfig(itemId)
	if config == nil {
		return ""
	}

	return fmt.Sprintf("%d", config.ItemType)
}

// 加载科技
func (self *CsvMgr) LoadNewUser() {
	GetCsvUtilMgr().LoadCsv("Newuseritem", &self.NewUserItem)
	//litter.Dump(self.NewUserItem)
}

// 加载以前的配置, 老的配置全部替换掉,替换了的写done, 持续重构ing...
func (self *CsvMgr) LoadOther() {
	self.ExpeditionConfig = make(map[int]*ExpeditionConfig)
	GetCsvUtilMgr().LoadCsv("Expedition", &self.ExpeditionConfig)
	GetCsvUtilMgr().LoadCsv("Luckyturntable_list", &self.LuckyturntablelistConfig)
	self.BuyskillpointsConfig = make(map[int]*BuyskillpointsConfig)
	GetCsvUtilMgr().LoadCsv("buyskillpoints", &self.BuyskillpointsConfig)
	GetCsvUtilMgr().LoadCsv("visitchance", &self.VisitchanceConfig)
	self.GrowthtaskConfig = make(map[int]*GrowthtaskConfig)
	GetCsvUtilMgr().LoadCsv("Growthtask", &self.GrowthtaskConfig)

	GetCsvUtilMgr().LoadCsv("luckdraw", &self.LuckdrawConfig)
	self.HorseSoulConfig = make(map[int]*HorseSoulConfig)
	GetCsvUtilMgr().LoadCsv("Horse_Soul", &self.HorseSoulConfig)
	GetCsvUtilMgr().LoadCsv("Gemsweeper_itembag", &self.GemsweeperitembagConfig)
	self.PubchestspecialConfig = make(map[int]*PubchestspecialConfig)
	GetCsvUtilMgr().LoadCsv("Pub_Chest_Special", &self.PubchestspecialConfig)
	self.SpyrewardConfig = make(map[int]*SpyrewardConfig)
	GetCsvUtilMgr().LoadCsv("spyreward", &self.SpyrewardConfig)
	self.HorseSoulUpgradeConfig = make(map[int]*HorseSoulUpgradeConfig)
	GetCsvUtilMgr().LoadCsv("Horse_Soul_Upgrade", &self.HorseSoulUpgradeConfig)
	self.CrownfightConfig = make(map[int]*CrownfightConfig)
	GetCsvUtilMgr().LoadCsv("Crown_Fight", &self.CrownfightConfig)
	GetCsvUtilMgr().LoadCsv("Crown_Fight", &self.CrownfightConfigLst)

	//litter.Dump(self.CrownfightConfig)

	self.SmeltPurchaseConfig = make(map[int]*SmeltPurchaseConfig)
	GetCsvUtilMgr().LoadCsv("Smelt_Purchase", &self.SmeltPurchaseConfig)
	//self.ActivityboxConfig = make(map[int]*ActivityboxConfig)
	//GetCsvUtilMgr().LoadCsv("Activitybox", &self.ActivityboxConfig)
	self.CommunityConfig = make(map[int]*CommunityConfig)
	GetCsvUtilMgr().LoadCsv("Guild_Lv", &self.CommunityConfig)
	self.SignrewardConfig = make(map[int]*SignrewardConfig)
	GetCsvUtilMgr().LoadCsv("Sign_Reward", &self.SignrewardConfig)
	self.SignConfig = make(map[int]*SignConfig)
	GetCsvUtilMgr().LoadCsv("Sign", &self.SignConfig)
	GetCsvUtilMgr().LoadCsv("luckegg_group", &self.LuckegggroupConfig)
	self.WarcontributionConfig = make(map[int]*WarcontributionConfig)
	GetCsvUtilMgr().LoadCsv("warcontribution", &self.WarcontributionConfig)
	GetCsvUtilMgr().LoadCsv("Gemsweeper_event", &self.GemsweepereventConfig)
	self.TrialConfig = make(map[int]*TrialConfig)
	GetCsvUtilMgr().LoadCsv("trial", &self.TrialConfig)
	self.LevelboxConfig = make(map[int]*LevelboxConfig)
	GetCsvUtilMgr().LoadCsv("Level_Box", &self.LevelboxConfig)
	self.LevelMapConfig = make(map[int]*LevelMapConfig)
	GetCsvUtilMgr().LoadCsv("Level_Map", &self.LevelMapConfig)

	GetCsvUtilMgr().LoadCsv("Treasure_Decompose", &self.TreasureDecomposeConfig)
	self.TimeGeneralsConfig = make(map[int]*TimeGeneralsConfig)
	GetCsvUtilMgr().LoadCsv("Time_Generals", &self.TimeGeneralsConfig)
	//litter.Dump(self.TimeGeneralsConfig)

	//! 世界事件
	GetCsvUtilMgr().LoadCsv("World_Event", &self.ExcitingConfig)
	GetCsvUtilMgr().LoadCsv("Summon_Box", &self.SummonBoxConfig)
	//litter.Dump(self.ExcitingConfig)

	GetCsvUtilMgr().LoadCsv("luckdraw_group", &self.LuckdrawgroupConfig)
	GetCsvUtilMgr().LoadCsv("Time_Generals_Rankaward", &self.TimeGeneralsRankawardConfig)
	self.WorldlvtpyeConfig = make(map[int]*WorldlvtpyeConfig)

	GetCsvUtilMgr().LoadCsv("World_Lvtpye", &self.WorldlvtpyeConfig)
	self.HorseParmConfig = make(map[int]*HorseParmConfig)
	GetCsvUtilMgr().LoadCsv("Horse_Parm", &self.HorseParmConfig)
	self.NpcConfig = make(map[int]*NpcConfig)
	GetCsvUtilMgr().LoadCsv("npc", &self.NpcConfig)
	GetCsvUtilMgr().LoadCsv("Holy_Legend", &self.HolyLegendConfig)

	// done
	self.TeamexpConfig = make(map[int]*TeamexpConfig)
	GetCsvUtilMgr().LoadCsv("Team_Exp", &self.TeamexpConfig)

	self.SmeltConfig = make(map[int]*SmeltConfig)
	GetCsvUtilMgr().LoadCsv("Smelt", &self.SmeltConfig)
	//litter.Dump(self.SmeltConfig)
	GetCsvUtilMgr().LoadCsv("Consumetop_Boss", &self.ConsumetopbossConfig)
	GetCsvUtilMgr().LoadCsv("Consumetop_List", &self.ConsumetoplistConfig)
	GetCsvUtilMgr().LoadCsv("Consumetop_Shop", &self.ConsumetopshopConfig)
	self.TreasureSuitConfig = make(map[int]*TreasureSuitConfig)
	GetCsvUtilMgr().LoadCsv("Treasure_Suit", &self.TreasureSuitConfig)
	self.SpytreasureConfig = make(map[int]*SpytreasureConfig)
	GetCsvUtilMgr().LoadCsv("spytreasure", &self.SpytreasureConfig)
	self.PeoplecityConfig = make(map[int]*PeoplecityConfig)
	GetCsvUtilMgr().LoadCsv("peoplecity", &self.PeoplecityConfig)
	self.NationalwarParmConfig = make(map[int]*NationalwarParmConfig)
	GetCsvUtilMgr().LoadCsv("Nationalwar_Parm", &self.NationalwarParmConfig)
	GetCsvUtilMgr().LoadCsv("expspeedup", &self.ExpspeedupConfig)
	GetCsvUtilMgr().LoadCsv("Homeoffice_Cityaward", &self.HomeofficecityawardConfig)
	GetCsvUtilMgr().LoadCsv("luckdraw_list", &self.LuckdrawlistConfig)
	self.ShopConfig = make(map[int]*ShopConfig)
	GetCsvUtilMgr().LoadCsv("Shop", &self.ShopConfig)
	self.HonourShopConfigMap = make(map[int]*HonourShopConfig)
	GetCsvUtilMgr().LoadCsv("Honour_Shop", &self.HonourShopConfigMap)

	GetCsvUtilMgr().LoadCsv("Luckegg_Config", &self.LuckeggConfig)

	self.DispatchConfig = make(map[int]*DispatchConfig)
	GetCsvUtilMgr().LoadCsv("Dispatch_Config", &self.DispatchConfig)

	//! 连续充值
	GetCsvUtilMgr().LoadCsv("Activity_DayRecharge", &self.DailyrechargeConfig)
	//litter.Dump(self.DailyrechargeConfig)

	self.ActivitynewConfig = make(map[int]*ActivitynewConfig)
	GetCsvUtilMgr().LoadCsv("Activitynew", &self. ActivitynewConfig)
	//litter.Dump(self.ActivitynewConfig)

	self.BuyphysicalConfig = make(map[int]*BuyphysicalConfig)
	GetCsvUtilMgr().LoadCsv("buyphysical", &self.BuyphysicalConfig)
	GetCsvUtilMgr().LoadCsv("Smelt_Drop", &self.SmeltDropConfig)
	GetCsvUtilMgr().LoadCsv("Time_Generals_Points", &self.TimeGeneralsPointsConfig)
	self.HorseJudgecallConfig = make(map[int]*HorseJudgecallConfig)
	GetCsvUtilMgr().LoadCsv("Holy_LegendLevel", &self.HolyLegendLevelConfig)
	self.PubchestdropgroupConfig = make(map[int]*PubchestdropgroupConfig)
	GetCsvUtilMgr().LoadCsv("Pub_Chest_Dropgroup", &self.PubchestdropgroupConfig)
	GetCsvUtilMgr().LoadCsv("Treasure_SuitAttribute", &self.TreasureSuitAttributeConfig)
	self.MaincityConfig = make(map[int]*MaincityConfig)
	GetCsvUtilMgr().LoadCsv("maincity", &self.MaincityConfig)
	self.WarcityConfig = make(map[int]*WarcityConfig) 
	GetCsvUtilMgr().LoadCsv("warcity", &self.WarcityConfig)
	self.SevendayConfig = make(map[int]*SevendayConfig)
	GetCsvUtilMgr().LoadCsv("Sevenday", &self.SevendayConfig)
	GetCsvUtilMgr().LoadCsv("Sevenday_TotalAward", &self.SevendayAward)
	GetCsvUtilMgr().LoadCsv("Homeoffice_Ranking", &self.HomeofficerankingConfig)
	self.SkillConfig = make(map[int]*SkillConfig)
	GetCsvUtilMgr().LoadCsv("Skillconfig", &self.SkillConfig)
	GetCsvUtilMgr().LoadCsv("redpacket_money", &self.RedpacketmoneyConfig)
	GetCsvUtilMgr().LoadCsv("Tiger_Symbol", &self.TigerSymbolConfig)
	GetCsvUtilMgr().LoadCsv("campreward", &self.CamprewardConfig)
	GetCsvUtilMgr().LoadCsv("Diplomacy_changes", &self.DiplomacychangesConfig)
	self.PubchesttotalConfig = make(map[int]*PubchesttotalConfig)
	GetCsvUtilMgr().LoadCsv("Pub_Chest_Total", &self.PubchesttotalConfig)

	//! 充值表读取
	self.MoneyConfig = make(map[int]*MoneyConfig)
	GetCsvUtilMgr().LoadCsv("Money", &self.MoneyConfig)
	//litter.Dump(self.MoneyConfig)

	self.TreasuryConfig = make(map[int]*TreasuryConfig)
	GetCsvUtilMgr().LoadCsv("Treasury", &self.TreasuryConfig)
	self.WarencourageConfig = make(map[int]*WarencourageConfig)
	GetCsvUtilMgr().LoadCsv("warencourage", &self.WarencourageConfig)
	self.WartargetConfig = make(map[int]*WartargetConfig)
	GetCsvUtilMgr().LoadCsv("wartarget", &self.WartargetConfig)
	self.FundConfig = make(map[int]*FundConfig)
	GetCsvUtilMgr().LoadCsv("Fund", &self.FundConfig)
	self.GemsweeperConfig = make(map[int]*GemsweeperConfig)
	GetCsvUtilMgr().LoadCsv("Gemsweeper", &self.GemsweeperConfig)
	self.FitConfig = make(map[int]*FitConfig)
	GetCsvUtilMgr().LoadCsv("fit", &self.FitConfig)
	self.NationalwardConfig = make(map[int]*NationalwarawardConfig)
	GetCsvUtilMgr().LoadCsv("Nationalwar_award", &self.NationalwardConfig)
	GetCsvUtilMgr().LoadCsv("Luckyturntable", &self.LuckyturntableConfig)
	self.VipConfigMap = make(map[int]*VipConfig)
	GetCsvUtilMgr().LoadCsv("Vip", &self.VipConfigMap)
	//litter.Dump(self.VipConfigMap)
	self.LegioncopyConfig = make(map[int]*LegioncopyConfig)
	GetCsvUtilMgr().LoadCsv("Guild_ActiveCopy", &self.LegioncopyConfig)
	self.GemsweeperawardConfig = make(map[int]*GemsweeperawardConfig)
	GetCsvUtilMgr().LoadCsv("Gemsweeper_award", &self.GemsweeperawardConfig)
	self.HonorkillConfig = make(map[int]*HonorkillConfig)
	GetCsvUtilMgr().LoadCsv("honorkill", &self.HonorkillConfig)
	self.ChestshopConfig = make(map[int]*ChestshopConfig)
	GetCsvUtilMgr().LoadCsv("chestshop", &self.ChestshopConfig)
	self.ExpeditionbuffConfig = make(map[int]*ExpeditionbuffConfig)
	GetCsvUtilMgr().LoadCsv("Expedition_Buff", &self.ExpeditionbuffConfig)
	self.ExpeditionbuffgroupConfig = make(map[int]*ExpeditionbuffgroupConfig)
	GetCsvUtilMgr().LoadCsv("Expedition_Buffgroup", &self.ExpeditionbuffgroupConfig)

	GetCsvUtilMgr().LoadCsv("Tiger_StuntUpgrade", &self.TigerStuntUpgradeConfig)
	GetCsvUtilMgr().LoadCsv("Homeoffice_Purchase", &self.HomeofficepurchaseConfig)
	self.GemsweeperextradropConfig = make(map[int]*GemsweeperextradropConfig)
	GetCsvUtilMgr().LoadCsv("Gemsweeper_extradrop", &self.GemsweeperextradropConfig)
	GetCsvUtilMgr().LoadCsv("Resource_Level", &self.ResourcelevelConfig)
	self.ShoprefreshConfig = make(map[int]*ShoprefreshConfig)
	GetCsvUtilMgr().LoadCsv("Shoprefresh", &self.ShoprefreshConfig)

	self.CityrepressConfig = make(map[int]*CityrepressConfig)
	GetCsvUtilMgr().LoadCsv("Gvg_Repress", &self.CityrepressConfig)
	GetCsvUtilMgr().LoadCsv("Gvg_Random", &self.CityrandomConfig)
	GetCsvUtilMgr().LoadCsv("Gvg_City", &self.CitygvgConfig)
	//litter.Dump(self.CityrepressConfig)
	//litter.Dump(self.CityrandomConfig)
	//litter.Dump(self.CitygvgConfig)

	GetCsvUtilMgr().LoadCsv("redpacket", &self.RedpacketConfig)
	GetCsvUtilMgr().LoadCsv("Playername", &self.PlayernameConfig)

	self.IndustryConfig = make(map[int]*IndustryConfig)
	GetCsvUtilMgr().LoadCsv("industry", &self.IndustryConfig)
	self.VisitConfig = make(map[int]*VisitConfig)
	GetCsvUtilMgr().LoadCsv("visit", &self.VisitConfig)
	self.ActiveConfig = make(map[int]*ActiveConfig)
	GetCsvUtilMgr().LoadCsv("Active", &self.ActiveConfig)
	self.DailytaskConfig = make(map[int]*DailytaskConfig)
	GetCsvUtilMgr().LoadCsv("Dailytask", &self.DailytaskConfig)
	self.TargetTaskConfig = make(map[int]*TargetTaskConfig)
	GetCsvUtilMgr().LoadCsv("Task_Target", &self.TargetTaskConfig)
	GetCsvUtilMgr().LoadCsv("Task_Badge", &self.BadgeTaskConfig)
	self.PromoteboxConfig = make(map[int]*PromoteboxConfig)
	GetCsvUtilMgr().LoadCsv("promotebox", &self.PromoteboxConfig)
	GetCsvUtilMgr().LoadCsv("Holy_Upgrade", &self.HolyUpgradeConfig)
	self.GemsweeperrankConfig = make(map[int]*GemsweeperrankConfig)
	GetCsvUtilMgr().LoadCsv("Gemsweeper_rank", &self.GemsweeperrankConfig)
	self.BuildingingConfig = make(map[int]*BuildingingConfig)
	GetCsvUtilMgr().LoadCsv("buildinging", &self.BuildingingConfig)
	GetCsvUtilMgr().LoadCsv("luckstart", &self.LuckstartConfig)
	GetCsvUtilMgr().LoadCsv("Invitingmoney", &self.InvitingConfig)
	GetCsvUtilMgr().LoadCsv("Jjc_Robot", &self.JJCRobotConfig)
	self.RobotMap = make(map[int]*JJCRobotConfig)
	GetCsvUtilMgr().LoadCsv("Jjc_Robot", &self.RobotMap)
	//litter.Dump(self.RobotMap)

	self.LevelMonsterMap = make(map[int]*LevelMonsterConfig)
	GetCsvUtilMgr().LoadCsv("Level_Monster", &self.LevelMonsterMap)

	GetCsvUtilMgr().LoadCsv("World_Level", &self.WorldLvLevelConfig)
	GetCsvUtilMgr().LoadCsv("World_Lvtpye", &self.WorldLvTpyeConfig)
	GetCsvUtilMgr().LoadCsv("show_boss", &self.WorldShowBossConfig)

	self.LoadHoly()
	self.LoadShop()
	self.LoadDailyTask()
	self.LoadPubChestTotal()
	self.LoadDropGroup()
	self.LoadPubChestSpecial()
	self.LoadSign()
	self.LoadTimeGeneralRank()

	//! 连续充值 - 配置
	self.LoadDailyrechargeConfig()
	//! 角色随机名字
	self.LoadPlayerName()

}

func (self *CsvMgr) getInitNum(id int) int {
	initStar, ok := GetCsvMgr().SimpleConfigMap[id]
	if !ok {
		return 0
	}
	return initStar.Num
}

func (self *CsvMgr) LoadSign() {
	self.SignMap = make(map[int]*SignConfig)
	for _, value := range self.SignConfig {
		self.SignMap[value.Month*10000+value.Sign] = value
	}
}

func NewTaskNode(taskId int, taskType int, n []int) *TaskNode {
	if len(n) < 4 {
		return &TaskNode{taskId, taskType, 0, 0, 0, 0}
	}

	return &TaskNode{
		Id:        taskId,
		Tasktypes: taskType,
		N1:        n[0],
		N2:        n[1],
		N3:        n[2],
		N4:        n[3],
	}
}

func (self *CsvMgr) LoadDailyTask() {
	for _, v := range self.DailytaskConfig {
		v.TaskNode = NewTaskNode(v.Taskid, v.Tasktypes, v.Ns)
	}

	for k, v := range self.GrowthtaskConfig {
		self.GrowthtaskConfig[k].TaskNode = NewTaskNode(v.Taskid, v.Tasktypes, v.Ns)
	}

	for _, v := range self.TargetTaskConfig {
		v.TaskNode = NewTaskNode(v.Taskid, v.Tasktypes, v.Ns)
	}
	//litter.Dump(self.GrowthtaskConfig)
}

func (self *CsvMgr) GetPhysicallimit(level int) int {
	config, ok := self.TeamexpConfig[level]
	if !ok {
		return 150
	}

	return config.Physicallimit
}

func (self *CsvMgr) GetTeamExp(level int) *TeamexpConfig {
	config, ok := self.TeamexpConfig[level]
	if !ok {
		return nil
	}

	return config
}

func (self *CsvMgr) LoadHoly() {
	// 圣物升级
	GetCsvUtilMgr().LoadCsv("Holy_PartsUpgrade", &self.HolyPartsUpgradeConfig)
	//litter.Dump(self.HolyPartsUpgradeConfig)
	self.HolyUpgradeMap = make(map[int]map[int]*HolyUpgradeConfig)
	for _, v := range self.HolyUpgradeConfig {
		pInfo, ok := self.HolyUpgradeMap[v.Beautyid]
		if !ok {
			pInfo = make(map[int]*HolyUpgradeConfig)
			self.HolyUpgradeMap[v.Beautyid] = pInfo
			pInfo[v.Stagelv] = v
		} else {
			pInfo[v.Stagelv] = v
		}
	}

	// 圣物部件
	self.HolyPartsMap = make(map[int]map[int]*HolyPartsUpgradeConfig)
	for _, v := range self.HolyPartsUpgradeConfig {
		pInfo, ok := self.HolyPartsMap[v.Treasureaid]
		if !ok {
			pInfo = make(map[int]*HolyPartsUpgradeConfig)
			self.HolyPartsMap[v.Treasureaid] = pInfo
			pInfo[v.Stagelv] = v
		} else {
			pInfo[v.Stagelv] = v
		}
	}
	//litter.Dump(self.HolyPartsMap)
}

// 获取圣物Id
func (self *CsvMgr) GetHolyPartConfig(treasureId int, stageLv int) *HolyPartsUpgradeConfig {
	config, ok := self.HolyPartsMap[treasureId]
	if !ok {
		return nil
	}

	pInfo, ok := config[stageLv]
	if !ok {
		return nil
	}

	return pInfo
}

// 加载竞技场配置
func (self *CsvMgr) LoadArena() {
	GetCsvUtilMgr().LoadCsv("Jjc_Award", &self.ArenaAwardConfig)
	self.ArenaAwardConfigMap = make(map[int]*PvpAwardConfig)
	for _, v := range self.ArenaAwardConfig {
		self.ArenaAwardConfigMap[v.Id] = v
	}
	//litter.Dump(self.PvpAwardConfig)
}

// 加载挂机掉落
func (self *CsvMgr) LoadHangUp() {
	self.HangUpConfig = make(map[int]*HangUpConfig)
	GetCsvUtilMgr().LoadCsv("Hang_Up", &self.HangUpConfig)
}

// 加载限制
func (self *CsvMgr) LoadTariff() {
	GetCsvUtilMgr().LoadCsv("Tariff", &self.TariffConfig)

}

// 次数限制表消耗
func (self *CsvMgr) GetTariffConfig(t int, times int) *TariffConfig {
	for _, v := range self.TariffConfig {
		if v.Type != t {
			continue
		}

		if times >= v.Rank1 && times <= v.Rank2 {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetTariffConfig2(t int) *TariffConfig {
	for _, v := range self.TariffConfig {
		if v.Type == t {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetTariffConfig3(t int, rank1 int) *TariffConfig {
	for _, v := range self.TariffConfig {
		if v.Type == t && rank1 == v.Rank1 {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetTariMaxTimes(t int) int {
	maxNum := 0
	for i := 0; i < len(GetCsvMgr().TariffConfig); i++ {
		cfg := GetCsvMgr().TariffConfig[i]
		if cfg.Type == t && cfg.Rank2 >= maxNum {
			maxNum = cfg.Rank2
		}
	}

	return maxNum
}

// 加载巨兽配置
func (self *CsvMgr) LoadBoss() {
	self.BossConfig = make(map[int]*BossConfig)
	GetCsvUtilMgr().LoadCsv("Hydra", &self.BossConfig)

	GetCsvUtilMgr().LoadCsv("Hydra_Lv", &self.BossLvConfig)
	//litter.Dump(self.BossLvConfig)
	self.CheckBossAttMap()
}

// 检查属性map
func (self *CsvMgr) CheckBossAttMap() {
	self.BossAttMap = make(map[int]map[int]*BossLvConfig)
	for _, v := range self.BossLvConfig {
		pInfo, ok := self.BossAttMap[v.Id]
		if !ok {
			pInfo = make(map[int]*BossLvConfig)
			pInfo[v.Lv] = v
			self.BossAttMap[v.Id] = pInfo
		} else {
			pInfo[v.Lv] = v
		}
	}

	//litter.Dump(self.BossAttMap)
}

func (self *CsvMgr) GetBossAtt(id int, level int) *BossLvConfig {
	pInfo, ok := self.BossAttMap[id]
	if !ok {
		return nil
	}
	res, ok := pInfo[level]
	if !ok {
		return nil
	}
	return res
}

// 加载宝石副本
func (self *CsvMgr) LoadGemStone() {
	GetCsvUtilMgr().LoadCsv("Gemstone_Level", &self.GemstoneLevelConfig)
	self.GemstoneLevelConfigMap = make(map[int]*GemstoneLevelConfig)
	for _, v := range self.GemstoneLevelConfig {
		self.GemstoneLevelConfigMap[v.Id] = v
	}
	GetCsvUtilMgr().LoadCsv("Gemstone_Chapter", &self.GemstoneChapterConfig)
	self.GemstoneChapterConfigMap = make(map[int]*GemstoneChapterConfig)
	for _, v := range self.GemstoneChapterConfig {
		self.GemstoneChapterConfigMap[v.Id] = v
	}
}

// 加载系统奖励
func (self *CsvMgr) LoadStatistics() {
	GetCsvUtilMgr().LoadCsv("System_Explain", &self.StatisticsConfig)
	self.StatisticsConfigMap = make(map[int]*StatisticsConfig)
	for _, v := range self.StatisticsConfig {
		self.StatisticsConfigMap[v.Type*1000+v.SubType] = v
	}

	GetCsvUtilMgr().LoadCsv("System_Award", &self.StatisticsRewardConfig)
	self.StatisticsRewardConfigMap = make(map[int]*StatisticsRewardConfig)
	for _, v := range self.StatisticsRewardConfig {
		self.StatisticsRewardConfigMap[v.Id] = v
	}
}

// 加载爵位系统
func (self *CsvMgr) LoadNobilityTask() {
	GetCsvUtilMgr().LoadCsv("Rank_Task", &self.NobilityConfig)
	self.NobilityConfigMap = make(map[int]*NobilityConfig)
	for _, v := range self.NobilityConfig {
		self.NobilityConfigMap[v.Id] = v
	}

	GetCsvUtilMgr().LoadCsv("Rank_Config", &self.NobilityReward)
	self.NobilityRewardMap = make(map[int]*NobilityReward)
	for _, v := range self.NobilityReward {
		if v.Id > 0 {
			self.NobilityRewardMap[v.Id] = v
		}
	}
}

// 加载月卡连续领取奖励
func (self *CsvMgr) LoadTotleAward() {
	GetCsvUtilMgr().LoadCsv("Activity_TotleAward", &self.ActivityTotleAward)
	self.ActivityTotleAwardMap = make(map[int]*ActivityTotleAward)
	for _, v := range self.ActivityTotleAward {
		self.ActivityTotleAwardMap[v.Id] = v
	}

	self.MonthCard = make(map[int]*MonthCard)
	GetCsvUtilMgr().LoadCsv("MonthCard", &self.MonthCard)

	self.MonthCardTotleAwardMap = make(map[int]*MonthCardTotleAwardMap)
	GetCsvUtilMgr().LoadCsv("MonthCard_TotleAward", &self.MonthCardTotleAwardMap)

	//litter.Dump(self.MonthCardTotleAwardMap)
}

// 加载全民商店
func (self *CsvMgr) LoadWholeShop() {
	GetCsvUtilMgr().LoadCsv("Activity_TimeShop", &self.WholeShopConfig)
	self.WholeShopConfigMap = make(map[int]*WholeShopConfig)
	for _, v := range self.WholeShopConfig {
		self.WholeShopConfigMap[v.Id] = v
	}

	GetCsvUtilMgr().LoadCsv("Activity_TimeShopRefresh", &self.WholeShopTimeConfig)
	self.WholeShopConfigTimeMap = make(map[int]*WholeShopTimeConfig)
	for _, v := range self.WholeShopTimeConfig {
		self.WholeShopConfigTimeMap[v.Id] = v
	}

	//litter.dump()
}

// 加载转盘
func (self *CsvMgr) LoadTurnTable() {
	TurnTableConfigTemp := make([]*TurnTableConfig, 0)
	GetCsvUtilMgr().LoadCsv("ActivityTurntable", &TurnTableConfigTemp)
	self.TurnTableConfigMap = make(map[int]*TurnTableConfig)
	for _, v := range TurnTableConfigTemp {
		self.TurnTableConfigMap[v.Id] = v
	}

	TurnTableTimeConfigTemp := make([]*TurnTableTimeConfig, 0)
	GetCsvUtilMgr().LoadCsv("ActivityTurntableTime", &TurnTableTimeConfigTemp)
	self.TurnTableTimeConfigMap = make(map[int]*TurnTableTimeConfig)
	for _, v := range TurnTableTimeConfigTemp {
		self.TurnTableTimeConfigMap[v.Id] = v
	}

	self.LotteryDrawConfigMap = make(map[int]*LotteryDrawConfig)
	GetCsvUtilMgr().LoadCsv("Activity_Lottery", &self.LotteryDrawConfigMap)

	//litter.dump()
}

// 充值配置
func (self *CsvMgr) LoadRecharge() {
	GetCsvUtilMgr().LoadCsv("Activity_Recharge", &self.RechargeConfig)
	self.RechargeConfigMap = make(map[int]*RechargeConfig)
	for _, v := range self.RechargeConfig {
		self.RechargeConfigMap[v.Id] = v
	}

	self.ActivityDailyRecharge = make(map[int]*ActivityDailyRecharge)
	GetCsvUtilMgr().LoadCsv("Activity_Daily_Recharge", &self.ActivityDailyRecharge)

	GetCsvUtilMgr().LoadCsv("Overflow_Gifts", &self.ActivityOverflowGifts)
}

// 加载战马相关配置
func (self *CsvMgr) LoadHorse() {
	GetCsvUtilMgr().LoadCsv("Horse_Score", &self.HorseAward)
	self.HorseAwardMap = make(map[int]*HorseAward)
	for index := range self.HorseAward {
		pConfig := self.HorseAward[index]
		if pConfig == nil {
			continue
		}
		self.HorseAwardMap[pConfig.Id] = pConfig
	}

	self.HorseSwitchMap = make(map[int]*HorseSwitch)
	GetCsvUtilMgr().LoadCsv("Horse_Switch", &self.HorseSwitchMap)

	self.HorseBattleSteedMap = make(map[int]*HorseBattleSteed)
	GetCsvUtilMgr().LoadCsv("Horse_BattleSteed", &self.HorseBattleSteedMap)
}

// 随机战马id
func (self *CsvMgr) randHorseId(horseId int) (int, error) {
	total := 0
	for _, horse := range self.HorseSwitchMap {
		if horse.HorseId == horseId {
			continue
		}
		total += horse.Rate
	}

	if total == 0 {
		return 0, errors.New(GetCsvMgr().GetText("STR_CSVMGR_BATTLE_HORSE_CONVERSION_TOTAL_WEIGHT"))
	}

	randNum := HF_GetRandom(total) + 1
	//LogDebug("randHorseId total:", total)
	check := 0
	//LogDebug("randNum num:", randNum)
	// 随机找到其中一项
	for _, horse := range self.HorseSwitchMap {
		if horse.HorseId == horseId {
			continue
		}

		check += horse.Rate
		//LogDebug("check num:", check)
		if randNum <= check {
			return horse.HorseId, nil
		}
	}

	return 0, errors.New(GetCsvMgr().GetText("STR_CSVMGR_SERVER_LOGIC_EXCEPTION"))
}

func (self *CsvMgr) getHeroBase(heroId int) map[int]*Attribute {
	res := make(map[int]*Attribute)
	pConfig := GetCsvMgr().GetHeroConfig(heroId)
	if pConfig == nil {
		return res
	}

	if len(pConfig.BaseTypes) != len(pConfig.BaseValues) {
		return res
	}

	for index := range pConfig.BaseTypes {
		attrType := pConfig.BaseTypes[index]
		attrValue := pConfig.BaseValues[index]
		if attrValue == 0 {
			continue
		}
		v, ok := res[attrType]
		if !ok {
			res[attrType] = &Attribute{AttType: attrType, AttValue: attrValue}
		} else {
			v.AttValue += attrValue
		}
	}
	return res
}

func (self *CsvMgr) getHeroGrowth(heroId int, playerlv int) map[int]*Attribute {
	res := make(map[int]*Attribute)
	pConfig := GetCsvMgr().GetHeroConfig(heroId)
	if pConfig == nil {
		return res
	}
	if len(pConfig.GrowthTypes) != len(pConfig.GrowthValues) {
		return res
	}

	for index := range pConfig.GrowthTypes {
		attrType := pConfig.GrowthTypes[index]
		//英雄的成长属性计算等级=玩家等级-1  20190627 by zy
		attrValue := pConfig.GrowthValues[index] * int64(playerlv-1)
		if attrValue == 0 {
			continue
		}
		v, ok := res[attrType]
		if !ok {
			res[attrType] = &Attribute{AttType: attrType, AttValue: attrValue}
		} else {
			v.AttValue += attrValue
		}
	}
	return res
}

type TimeInfo struct {
	Hour   int
	Minute int
	Second int
}

func NewTimeInfo(timeStamp int) *TimeInfo {
	hour := timeStamp / 10000
	timeStamp = timeStamp % 10000
	minute := timeStamp / 100
	timeStamp = timeStamp % 100
	second := timeStamp
	return &TimeInfo{hour, minute, second}
}

func (self *CsvMgr) getNextDay() int64 {
	now := TimeServer()
	timeSet := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local)
	if now.Unix() > timeSet.Unix() {
		return timeSet.Unix() + int64(DAY_SECS)
	} else {
		return timeSet.Unix()
	}
}

func (self *CsvMgr) GetTodyDay() int64 {
	now := TimeServer()
	timeSet := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local)
	return timeSet.Unix()
}

func (self *CsvMgr) LoadSmeltTask() {
	GetCsvUtilMgr().LoadCsv("Smelt_Reward", &self.SmeltAward)
	GetCsvUtilMgr().LoadCsv("Smelt_BuyAward", &self.SmeltBuyAward)
	self.SmeltAwardMap = make(map[int]*SmeltAward)
	for index := range self.SmeltAward {
		pConfig := self.SmeltAward[index]
		if pConfig == nil {
			continue
		}
		self.SmeltAwardMap[pConfig.Id] = pConfig
	}

	self.SmeltBuyAwardMap = make(map[int]*SmeltAward)
	for index := range self.SmeltBuyAward {
		pConfig := self.SmeltBuyAward[index]
		if pConfig == nil {
			continue
		}
		self.SmeltBuyAwardMap[pConfig.Id] = pConfig
	}
}

// 加载佣兵配置
func (self *CsvMgr) LoadMercenary() {
	GetCsvUtilMgr().LoadCsv("Mercenary_Random", &self.MercenaryRandom)
	GetCsvUtilMgr().LoadCsv("Mercenary_Lv", &self.MercenaryLv)
	self.MercenaryConfig = make(map[int]*MercenaryConfig)
	GetCsvUtilMgr().LoadCsv("Mercenary_config", &self.MercenaryConfig)
	self.MercenaryLvMap = make(map[int]map[int]*MercenaryLv)
	for _, v := range self.MercenaryLv {
		info, ok := self.MercenaryLvMap[v.Id]
		if !ok {
			info = make(map[int]*MercenaryLv)
			info[v.Lv] = v
			self.MercenaryLvMap[v.Id] = info
		} else {
			info[v.Lv] = v
		}
	}
	//fmt.Printf("%# v", pretty.Formatter(self.MercenaryLvMap))
	self.MercenaryRandomGroup = make(map[int]map[int][]*MercenaryRandom)
	for _, v := range self.MercenaryRandom {
		info, ok := self.MercenaryRandomGroup[v.Class]
		if !ok {
			info = make(map[int][]*MercenaryRandom)
			info[v.Index] = append(info[v.Index], v)
			self.MercenaryRandomGroup[v.Class] = info
		} else {
			info[v.Index] = append(info[v.Index], v)
		}
	}
	//fmt.Printf("%# v", pretty.Formatter(self.MercenaryRandomGroup))
}

func (self *CsvMgr) LoadMiliTask() {
	// 派遣周任务
	GetCsvUtilMgr().LoadCsv("Dispatch_Task", &self.MiliWeekTaskConfig)
	self.MiliWeekTaskMapConfig = make(map[int]*MiliWeekTaskConfig)
	for index := range self.MiliWeekTaskConfig {
		pConfig := self.MiliWeekTaskConfig[index]
		if pConfig == nil {
			continue
		}
		self.MiliWeekTaskMapConfig[pConfig.Id] = pConfig
	}
}

// 配置1
func (self *CsvMgr) GetConsumetoplistConfig1(group int) []*ConsumetoplistConfig {
	var res []*ConsumetoplistConfig
	for _, v := range self.ConsumetoplistConfig {
		if v.Group != group {
			continue
		}
		if v.Type == 2 {
			continue
		}
		res = append(res, v)
	}

	return res
}

// 配置2
func (self *CsvMgr) GetConsumetoplistConfig2() []*ConsumetoplistConfig {
	var res []*ConsumetoplistConfig
	for _, v := range self.ConsumetoplistConfig {
		if v.Type != 2 {
			continue
		}
		res = append(res, v)
	}
	return res
}

// 获取配置
func (self *CsvMgr) GetConsumetopshopConfig(group int) []*ConsumetopshopConfig {
	var res []*ConsumetopshopConfig
	for _, v := range self.ConsumetopshopConfig {
		if v.Group != group {
			continue
		}
		res = append(res, v)
	}

	return res
}

func (self *CsvMgr) GetKingName(id int) string {
	config, ok := self.WorldPowerMap[id]
	if !ok {
		return "NULL"
	}
	return config.KingName
}

func (self *CsvMgr) GetUnionName(id int) string {
	config, ok := self.WorldPowerMap[id]
	if !ok {
		return "NULL"
	}
	return config.Name
}

// 获取第一个KingTask
func (self *CsvMgr) getFirstKingTask() *GrowthtaskKingConfig {
	if len(self.TaskKingConfig) > 0 {
		return self.TaskKingConfig[0]
	}
	return nil
}

func (self *CsvMgr) GetCampName(camp int) string {
	countryName := ""
	if camp == 1 {
		countryName = self.GetText("STR_CAMP_1")
	} else if camp == 2 {
		countryName = self.GetText("STR_CAMP_2")
	} else if camp == 3 {
		countryName = self.GetText("STR_CAMP_3")
	}
	return countryName
}

func (self *CsvMgr) GetWarMaxLv(playType int) int {
	maxLv := 0
	for _, v := range self.WarEncourageConfig {
		if playType != v.Type {
			continue
		}
		if maxLv < v.Lv {
			maxLv = v.Lv
		}
	}
	return maxLv
}

func (self *CsvMgr) GetEncourageConfig(playType int, curLv int) *WarEncourageConfig {
	for _, v := range self.WarEncourageConfig {
		if playType != v.Type {
			continue
		}

		if curLv != v.Lv {
			continue
		}

		return v
	}
	return nil
}

// team win
func (self *CsvMgr) IsTeamWin(fight int64, num int, winParam int, loseParam int, reduceParam int, maxChance int, minChance int) bool {
	firstRes := float32(fight) / float32(num)
	win := float32(winParam) / float32(100.0)
	lose := float32(loseParam) / float32(100.0)

	reduce := float32(reduceParam / 100.0)
	secondRes := float32(1.0-firstRes) * reduce
	n := HF_GetRandom(100)
	if firstRes > win {
		return true
	} else if firstRes < lose {
		return false
	} else {
		if reduce <= 0.0 {
			return false
		}
		left := float32(100.0) - secondRes
		if left >= float32(maxChance) {
			return n < maxChance
		} else if left <= float32(minChance) {
			return n < int(secondRes*float32(100))
		} else {
			return n < int(left)
		}
	}
	return false
}

func (self *CsvMgr) IsDungenTeamWin(fight int64, config *TeamDungeonConfig) bool {
	fightNum := fight / 100
	LogDebug("fightNum", fightNum, ", forcelimit:", config.ForceLimit)
	return self.IsTeamWin(fightNum, config.ForceLimit, config.Win, config.Lose, config.Reduce, config.MaxChance, config.MinChance)
}

func (self *CsvMgr) IsArmyTeamWin(fight int64, config *ArmyTeamConfig) bool {
	num := self.GetWorldFight()
	num = num / 100
	fightNum := fight / 100
	LogDebug("fightNum", fightNum, ", 世界战力值:", num)
	return self.IsTeamWin(fightNum, num, config.Win, config.Lose, config.Reduce, config.MaxChance, config.MinChance)
}

// 加载英雄升级
func (self *CsvMgr) LoadHeroExp() {
	GetCsvUtilMgr().LoadCsv("heroexp", &self.HeroExpConfig)
	self.HeroExpConfigMap = make(map[int]*HeroExpConfig)
	for _, v := range self.HeroExpConfig {
		self.HeroExpConfigMap[v.HeroLv] = v
	}
}

//符文系统相关配置
func (self *CsvMgr) LoadRuneConfig() {
	GetCsvUtilMgr().LoadCsv("rune_compose", &self.RuneCompose)
	self.RuneComposeMap = make(map[int]*RuneCompose)
	for _, v := range self.RuneCompose {
		var num = make([]int, len(v.Item))
		for i := 0; i < len(v.Item); i++ {
			num[i] = 1
		}
		v.Num = num
		self.RuneComposeMap[v.Id] = v
	}

	GetCsvUtilMgr().LoadCsv("rune_Config", &self.RuneConfig)
	self.RuneConfigMap = make(map[int]*RuneConfig)
	for _, v := range self.RuneConfig {
		self.RuneConfigMap[v.Id] = v
	}
}

//阵型系统相关配置
func (self *CsvMgr) LoadFormationConfig() {
	GetCsvUtilMgr().LoadCsv("formation", &self.FormationConfig)
	self.FormationConfigMap = make(map[int]*FormationConfig)
	for _, v := range self.FormationConfig {
		self.FormationConfigMap[v.Id] = v
	}
}

func (self *CsvMgr) LoadFundConfig() {
	self.FundConfigMap = make(map[int]*FundConfigMap)
	GetCsvUtilMgr().LoadCsv("Activity_Fund", &self.FundConfigMap)
}

//新商店
func (self *CsvMgr) LoadNewShopConfig() {

	GetCsvUtilMgr().LoadCsv("Shop_Discount", &self.NewShopDiscount)

	tempNewShopConfig := make([]*NewShopConfig, 0)
	GetCsvUtilMgr().LoadCsv("Shop", &tempNewShopConfig)
	self.NewShopConfigMap = make(map[int]map[int][]*NewShopConfig)
	for _, v := range tempNewShopConfig {
		_, okType := self.NewShopConfigMap[v.Type]
		if !okType {
			shopType := make(map[int][]*NewShopConfig)
			self.NewShopConfigMap[v.Type] = shopType
		}
		self.NewShopConfigMap[v.Type][v.Grid] = append(self.NewShopConfigMap[v.Type][v.Grid], v)
	}
}

// 加载神兽配置
func (self *CsvMgr) LoadHydra() {
	self.HydraConfigMap = make(map[int]*HydraConfig)
	GetCsvUtilMgr().LoadCsv("Hydra_Config", &self.HydraConfigMap)

	self.HydraSkillMap = make(map[int]map[int]*HydraSkill)
	GetCsvUtilMgr().LoadCsv("Hydra_Skill", &self.HydraSkill)
	for _, v := range self.HydraSkill {
		pInfo, ok := self.HydraSkillMap[v.SkillID]
		if !ok {
			pInfo = make(map[int]*HydraSkill)
			pInfo[v.HydraLv] = v
			self.HydraSkillMap[v.SkillID] = pInfo
		} else {
			pInfo[v.HydraLv] = v
		}
	}

	//litter.Dump(self.HydraSkillMap)

	self.HydraLevelMap = make(map[int]*HydraLevel)
	GetCsvUtilMgr().LoadCsv("Hydra_lvup", &self.HydraLevelMap)

	//litter.Dump(self.HydraLevelMap)

	self.HydraStepMap = make(map[int]*HydraStep)
	GetCsvUtilMgr().LoadCsv("Hydra_advance", &self.HydraStepMap)

	//litter.Dump(self.HydraStepMap)

	// self.HydraStarMap = make(map[int]map[int]*HydraStar)
	//GetCsvUtilMgr().LoadCsv("Hydra_Star", &self.HydraStar)
	//for _, v := range self.HydraStar {
	//	pInfo, ok := self.HydraStarMap[v.HydraId]
	//	if !ok {
	//		pInfo = make(map[int]*HydraStar)
	//		pInfo[v.HydraLv] = v
	//		self.HydraStarMap[v.HydraId] = pInfo
	//	} else {
	//		pInfo[v.HydraLv] = v
	//	}
	//}

	//self.HydraTaskMap = make(map[int]*HydraTask)
	//GetCsvUtilMgr().LoadCsv("Hydra_Task", &self.HydraTask)
	//for _, v := range self.HydraTask {
	//	self.HydraTaskMap[v.TaskId] = v
	//}
}

//加载地牢配置
func (self *CsvMgr) LoadPit() {
	GetCsvUtilMgr().LoadCsv("Dungeons_config", &self.PitConfig)
	self.PitConfigMap = make(map[int]*PitConfig)
	for _, v := range self.PitConfig {
		self.PitConfigMap[v.Id] = v
	}

	GetCsvUtilMgr().LoadCsv("Dungeons_map", &self.PitMap)
	self.PitMapMap = make(map[int][]*PitMap)
	for _, v := range self.PitMap {
		self.PitMapMap[v.PitType] = append(self.PitMapMap[v.PitType], v)
	}

	self.PitBuffMap = make(map[int]*PitBuff)
	GetCsvUtilMgr().LoadCsv("Dungeons_buff", &self.PitBuffMap)

	self.PitBoxMap = make(map[int]*PitBox)
	GetCsvUtilMgr().LoadCsv("Dungeons_box", &self.PitBoxMap)

	self.PitMonsterMap = make(map[int]*PitMonster)
	GetCsvUtilMgr().LoadCsv("Dungeons_monster", &self.PitMonsterMap)
}

//加载新地牢配置
func (self *CsvMgr) LoadNewPit() {

	NewPitConfigMapTemp := make([]*NewPitConfig, 0)
	GetCsvUtilMgr().LoadCsv("Newpit_Map", &NewPitConfigMapTemp)
	self.NewPitConfigMap = make(map[int]*NewPitConfig)
	for _, v := range NewPitConfigMapTemp {
		self.NewPitConfigMap[v.Id] = v
	}

	NewPitReliqueTemp := make([]*NewPitRelique, 0)
	GetCsvUtilMgr().LoadCsv("Newpit_Relique", &NewPitReliqueTemp)
	self.NewPitRelique = make(map[int]*NewPitRelique, 0)
	for _, v := range NewPitReliqueTemp {
		self.NewPitRelique[v.Id] = v
	}

	GetCsvUtilMgr().LoadCsv("Newpit_RobotCartExclusive", &self.NewPitRobotExclusiveConfig)

	GetCsvUtilMgr().LoadCsv("Newpit_Robot", &self.NewPitRobotConfig)
	self.NewPitRobotGroupMap = make(map[int][]int)
	self.NewPitRobotGroupFirstMap = make(map[int][]int)
	for i := 0; i < len(self.NewPitRobotConfig); i++ {
		if self.NewPitRobotConfig[i].HeroStar != 1 {
			continue
		}
		self.NewPitRobotGroupMap[self.NewPitRobotConfig[i].RobotGroup] = append(self.NewPitRobotGroupMap[self.NewPitRobotConfig[i].RobotGroup], self.NewPitRobotConfig[i].HeroId)
		self.NewPitRobotGroupFirstMap[self.NewPitRobotConfig[i].FirstRobotGroup] = append(self.NewPitRobotGroupFirstMap[self.NewPitRobotConfig[i].FirstRobotGroup], self.NewPitRobotConfig[i].HeroId)
	}
	GetCsvUtilMgr().LoadCsv("Newpit_RobotAttribute", &self.NewPitRobotAttr)
	GetCsvUtilMgr().LoadCsv("Newpit_RobotQuality", &self.NewPitRobotQuality)
	GetCsvUtilMgr().LoadCsv("Newpit_RobotMonsterQuality", &self.NewPitRobotMonsterQuality)
	//GetCsvUtilMgr().LoadCsv("Newpit_RobotMonsterLv", &self.NewPitRobotMonsterLv)

	NewPitTreasureCaveTemp := make([]*NewPitTreasureCave, 0)
	GetCsvUtilMgr().LoadCsv("Treasure_Caves", &NewPitTreasureCaveTemp)
	self.NewPitTreasureCave = make(map[int]*NewPitTreasureCave)
	for _, v := range NewPitTreasureCaveTemp {
		self.NewPitTreasureCave[v.Lv] = v
	}

	self.NewPitExtraReward = make(map[int]*NewPitExtraReward)
	GetCsvUtilMgr().LoadCsv("Newpit_ExtraReward", &self.NewPitExtraReward)

	NewPitParamTemp := make([]*NewPitParam, 0)
	GetCsvUtilMgr().LoadCsv("Newpit_Param", &NewPitParamTemp)
	self.NewPitParam = make(map[int]*NewPitParam)
	for _, v := range NewPitParamTemp {
		self.NewPitParam[v.Quality] = v
	}

	GetCsvUtilMgr().LoadCsv("Newpit_Difficulty", &self.NewPitDifficulty)
}

//加载时光之巅
func (self *CsvMgr) LoadInstanceConfig() {

	InstanceConfigMapTemp := make([]*InstanceConfig, 0)
	GetCsvUtilMgr().LoadCsv("Interesting_Config", &InstanceConfigMapTemp)
	self.InstanceConfig = make(map[int]*InstanceConfig)
	for _, v := range InstanceConfigMapTemp {
		self.InstanceConfig[v.MapId] = v
	}

	InstanceBoxTemp := make([]*InstanceBox, 0)
	GetCsvUtilMgr().LoadCsv("Interesting_Box", &InstanceBoxTemp)
	self.InstanceBox = make(map[int][]*InstanceBox, 0)
	for _, v := range InstanceBoxTemp {
		self.InstanceBox[v.GroupId] = append(self.InstanceBox[v.GroupId], v)
	}

	self.InstanceThing = make(map[int]map[int]*InstanceThing)
	for _, fn := range self.InstanceConfig {
		fileName := "Interesting_Thing_" + strconv.Itoa(fn.MapId)
		csvData := make(map[int]*InstanceThing)
		GetCsvUtilMgr().LoadEventsCsv(fileName, &csvData)
		if len(csvData) <= 0 {
			continue
		}
		self.InstanceThing[fn.MapId] = csvData
	}
	//litter.Dump(self.InstanceThing)
}

func (self *CsvMgr) LoadEntanglement() {
	GetCsvUtilMgr().LoadCsv("Fate_config", &self.EntanglementConfig)
	for _, v := range self.EntanglementConfig {
		nLen := len(v.HeroId)
		for i := nLen - 1; i >= 0; i-- {
			if v.HeroId[i] <= 0 {
				v.HeroId = append(v.HeroId[:i], v.HeroId[i+1:]...)
			}
		}

		if len(v.HeroId) != v.HeroNum {
			continue
		}
		if len(v.BaseValue) != len(v.BaseType) {
			continue
		}
		nLen = len(v.BaseValue)
		for i := nLen - 1; i >= 0; i-- {
			if v.BaseValue[i] <= 0 || v.BaseType[i] <= 0 {
				v.BaseValue = append(v.BaseValue[:i], v.BaseValue[i+1:]...)
				v.BaseType = append(v.BaseType[:i], v.BaseType[i+1:]...)
			}
		}

		_, ok := self.EntanglementMapConfig[v.Group]
		if !ok {
			data := EntanglementFate{}
			data.HeroId = v.HeroId
			data.Group = v.Group
			data.HeroNum = v.HeroNum
			data.Property = make([]*EntanglementProperty, 0)

			property := EntanglementProperty{}
			property.BaseType = v.BaseType
			property.BaseValue = v.BaseValue
			property.FateNum = v.FateNum
			property.MinQuality = v.MinQuality

			data.Property = append(data.Property, &property)

			self.EntanglementMapConfig[v.Group] = &data
		} else {
			property := EntanglementProperty{}
			property.BaseType = v.BaseType
			property.BaseValue = v.BaseValue
			property.FateNum = v.FateNum
			property.MinQuality = v.MinQuality

			self.EntanglementMapConfig[v.Group].Property = append(self.EntanglementMapConfig[v.Group].Property, &property)
		}
	}
}

func (self *CsvMgr) GetEntanglementConfig(nType int) *EntanglementFate {
	config, ok := GetCsvMgr().EntanglementMapConfig[nType]
	if ok {
		return config
	}
	return nil
}

func (self *CsvMgr) GetNewPitRobot(star int, cart int) []*NewPitRobotConfig {
	rel := make([]*NewPitRobotConfig, 0)
	for _, v := range self.NewPitRobotConfig {
		if v.Cart == cart && v.HeroStar == star {
			rel = append(rel, v)
		}
	}
	return rel
}

func (self *CsvMgr) GetNewPitRobotByHeroId(heroId int, star int) *NewPitRobotConfig {
	for _, v := range self.NewPitRobotConfig {
		if v.HeroId == heroId && v.HeroStar == star {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetNewPitRobotExclusive(lv int, num int) *NewPitRobotExclusiveConfig {
	for _, v := range self.NewPitRobotExclusiveConfig {
		if v.ExclusiveLv[0] <= lv && lv <= v.ExclusiveLv[1] && num >= v.ExclusiveNum {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetNewPitRobotAttr(group int, lv int) *NewPitRobotAttr {
	for _, v := range self.NewPitRobotAttr {
		if v.Group == group && v.Lv == lv {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetNewPitRobotQuality(star int) *NewPitRobotQuality {
	for _, v := range self.NewPitRobotQuality {
		if v.MeanStar1 <= star*10 && v.MeanStar2 > star*10 {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetNewPitRobotMonsterQuality(lv int) *NewPitRobotMonsterQuality {
	for _, v := range self.NewPitRobotMonsterQuality {
		if v.MeanStar1 <= lv && v.MeanStar2 > lv {
			return v
		}
	}
	return nil
}

//func (self *CsvMgr) GetNewPitRobotMonsterLv(star int, fight int64) *NewPitRobotMonsterLv {
//	for _, v := range self.NewPitRobotMonsterLv {
//		if v.StarMin <= star && v.StarMax >= star && v.FightMin <= fight && v.FightMax >= fight {
//			return v
//		}
//	}
//	return nil
//}

func (self *CsvMgr) LoadRewardConfig() {
	GetCsvUtilMgr().LoadCsv("Rewardforbar", &self.RewardForbarConfig)

	//for _, v := range self.RewardForbarConfig {
	//	nLen := len(v.NeedCamp)
	//	for i := nLen - 1; i >= 0; i-- {
	//		if v.NeedCamp[i] == 0 {
	//			v.NeedCamp = append(v.NeedCamp[:i], v.NeedCamp[i+1:]...)
	//		}
	//	}
	//}

	for _, v := range self.RewardForbarConfig {
		nTeam := v.IsTeam

		_, ok := self.RewardForbarMapConfig[nTeam]

		if ok {
			self.RewardForbarMapConfig[nTeam] = append(self.RewardForbarMapConfig[nTeam], v)
		} else {
			self.RewardForbarMapConfig[nTeam] = make([]*RewardForbarConfig, 0)
			self.RewardForbarMapConfig[nTeam] = append(self.RewardForbarMapConfig[nTeam], v)
		}
	}

	for _, v := range self.RewardForbarConfig {
		nTeam := v.IsTeam
		nGroup := v.Group

		_, ok1 := self.RewardForbarColorMapConfig[nTeam]

		if ok1 {
			_, ok2 := self.RewardForbarColorMapConfig[nTeam][nGroup]
			if ok2 {
				self.RewardForbarColorMapConfig[nTeam][nGroup] = append(self.RewardForbarColorMapConfig[nTeam][nGroup], v)
			} else {
				self.RewardForbarColorMapConfig[nTeam][nGroup] = make([]*RewardForbarConfig, 0)
				self.RewardForbarColorMapConfig[nTeam][nGroup] = append(self.RewardForbarColorMapConfig[nTeam][nGroup], v)
			}
		} else {
			self.RewardForbarColorMapConfig[nTeam] = make(map[int][]*RewardForbarConfig)
			self.RewardForbarColorMapConfig[nTeam][nGroup] = make([]*RewardForbarConfig, 0)
			self.RewardForbarColorMapConfig[nTeam][nGroup] = append(self.RewardForbarColorMapConfig[nTeam][nGroup], v)
		}
	}

	GetCsvUtilMgr().LoadCsv("Rewardforbarlvup", &self.RewardForbarLvUpConfig)

	for _, v := range self.RewardForbarLvUpConfig {
		nLen := len(v.Renovatestar)
		if nLen != len(v.Renovatepro) {
			LogError("Rewardforbarlvup errer nLen != len(v.Renovatepro)")
		}

		for i := nLen - 1; i >= 0; i-- {
			if v.Renovatestar[i] == 0 {
				v.Renovatestar = append(v.Renovatestar[:i], v.Renovatestar[i+1:]...)
				v.Renovatepro = append(v.Renovatepro[:i], v.Renovatepro[i+1:]...)
			}
		}
	}

	GetCsvUtilMgr().LoadCsv("Rewardforbarprize", &self.RewardForbarAward)
	self.RewardForbarAwardMap = make(map[int]map[int][]*RewardForbarPrize)
	for _, v := range self.RewardForbarAward {
		_, ok1 := self.RewardForbarAwardMap[v.Isteam]
		if !ok1 {
			self.RewardForbarAwardMap[v.Isteam] = make(map[int][]*RewardForbarPrize)
			self.RewardForbarAwardMap[v.Isteam][v.Group] = make([]*RewardForbarPrize, 0)
			self.RewardForbarAwardMap[v.Isteam][v.Group] = append(self.RewardForbarAwardMap[v.Isteam][v.Group], v)
		} else {
			_, ok2 := self.RewardForbarAwardMap[v.Isteam][v.Group]
			if !ok2 {
				self.RewardForbarAwardMap[v.Isteam][v.Group] = make([]*RewardForbarPrize, 0)
				self.RewardForbarAwardMap[v.Isteam][v.Group] = append(self.RewardForbarAwardMap[v.Isteam][v.Group], v)
			} else {
				self.RewardForbarAwardMap[v.Isteam][v.Group] = append(self.RewardForbarAwardMap[v.Isteam][v.Group], v)
			}
		}
	}
}

func (self *CsvMgr) GetRewardAwardConfig(isTeam, group int) []*RewardForbarPrize {
	data, ok1 := self.RewardForbarAwardMap[isTeam]
	if ok1 {
		_, ok2 := data[group]
		if ok2 {
			return data[group]
		}
	}
	return nil
}

func (self *CsvMgr) GetRewardConfig(isTeam, id int) *RewardForbarConfig {
	data, ok := self.RewardForbarMapConfig[isTeam]
	if ok {
		for _, v := range data {
			if v.ID == id {
				return v
			}
		}
	}
	return nil
}

func (self *CsvMgr) GetRewardForbarLvUpConfig(lv int) *RewardForbarLvUpConfig {
	for _, v := range self.RewardForbarLvUpConfig {
		if v.LV == lv {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetRewardForbarGroupConfig(isTeam, group int) []*RewardForbarConfig {
	configs, ok1 := self.RewardForbarColorMapConfig[isTeam]
	if ok1 {
		config, ok2 := configs[group]
		if ok2 {
			return config
		}
	}
	return nil
}

// 读取凯旋丰碑
func (self *CsvMgr) LoadRankListIntegralConfig() {
	GetCsvUtilMgr().LoadCsv("Ranklist_Integral", &self.RankListIntegral)
}

func (self *CsvMgr) GetRankListIntegralConfig(quality int) *RankListIntegral {
	for _, v := range self.RankListIntegral {
		if v.HeroQuality == quality {
			return v
		}
	}
	return nil
}

// 加载排行任务配置
func (self *CsvMgr) LoadRankTask() {
	GetCsvUtilMgr().LoadCsv("Ranklist_Task", &self.RankTaskConfig)

	for _, v := range self.RankTaskConfig {
		if len(v.Items) != len(v.Nums) {
			continue
		}

		if len(v.Ns) != 4 {
			continue
		}

		nType := v.Type

		_, ok := self.RankTaskMapConfig[nType]
		if ok {
			self.RankTaskMapConfig[nType] = append(self.RankTaskMapConfig[nType], v)
		} else {
			var data []*RankTaskConfig
			data = append(data, v)

			self.RankTaskMapConfig[nType] = data
		}
	}
}

func (self *CsvMgr) GetRankTaskConfig(nType, nSort int) *RankTaskConfig {
	configs, ok := GetCsvMgr().RankTaskMapConfig[nType]
	if !ok {
		return nil
	}

	for _, v := range configs {
		if v.Sort != nSort {
			continue
		}
		return v
	}
	return nil
}

func (self *CsvMgr) GetRankTaskConfigByID(nID int) *RankTaskConfig {

	configs := GetCsvMgr().RankTaskConfig
	for i := 0; i < len(configs); i++ {
		if nID == configs[i].Id {
			return configs[i]
		}
	}

	return nil
}

// 读取列表
func (self *CsvMgr) LoadResonanceCrystalconfig() {
	GetCsvUtilMgr().LoadCsv("Crystal_config", &self.ResonanceCrystalconfig)
}

func (self *CsvMgr) GetResonanceCrystalconfig(level int) *ResonanceCrystalconfig {
	configs := GetCsvMgr().ResonanceCrystalconfig
	config := &ResonanceCrystalconfig{}
	for i := 0; i < len(configs); i++ {
		if configs[i].MaxLevel < level {
			config = configs[i]
		} else {
			return config
		}
	}
	return nil
}

func (self *CsvMgr) GetHireEquipConfig(heroType int, fight int64, subType int) *EquipHireConfig {

	for _, v := range self.EquipHireConfig {
		if subType == v.SubType && heroType == v.Type && int(fight) <= v.Max && int(fight) >= v.Min {
			return v
		}
	}

	return nil
}

//常规时间计算,使用类型3
func (self *CsvMgr) GetNowStartAndEnd(systemId int) (int64, int64) {
	startTime := int64(0)
	endTime := int64(0)
	for _, config := range GetCsvMgr().TimeResetConfig {
		if config.System == systemId {
			switch config.TimeType {
			case TIME_RESET_TYPE_TIME:
				stage := config.Continue + config.Cd
				calCount := (TimeServer().Unix() - config.Time[0]) / stage
				calTime := (TimeServer().Unix() - config.Time[0]) % stage
				//处于冷却期则计算期数加1
				if calTime > config.Continue {
					calCount += 1
				}
				startTime = config.Time[0] + calCount*stage
				endTime = startTime + config.Continue
			}
		}
	}
	return startTime, endTime
}

//使用类型1,根据角色创建天数
func (self *CsvMgr) GetNowStartAndEndByRoleDays(systemId int, createTime string) (int64, int64) {
	startTime := int64(0)
	endTime := int64(0)
	rTime, _ := time.ParseInLocation(DATEFORMAT, createTime, time.Local)
	rDay := time.Date(rTime.Year(), rTime.Month(), rTime.Day(), 5, 0, 0, 0, TimeServer().Location()).Unix()
	for _, config := range GetCsvMgr().TimeResetConfig {
		if config.System == systemId {
			switch config.TimeType {
			case TIME_RESET_TYPE_CREATE_ROLE:
				stage := config.Continue + config.Cd
				calCount := (TimeServer().Unix() - rDay) / stage
				calTime := (TimeServer().Unix() - rDay) % stage
				//处于冷却期则计算期数加1
				if calTime > config.Continue {
					calCount += 1
				}
				startTime = rDay + calCount*stage
				endTime = startTime + config.Continue
			}
		}
	}
	return startTime, endTime
}

//基准时间计算,使用类型4，返回一个天数
func (self *CsvMgr) GetTimeBaseDay(systemId int) (int64, int64) {
	day := int64(1)
	startTime := int64(0)
	for _, config := range GetCsvMgr().TimeResetConfig {
		if config.System == systemId {
			switch config.TimeType {
			case TIME_RESET_TYPE_TIME_BASE:
				if config.Time[0] > 0 {
					day = (TimeServer().Unix()-config.Time[0])/DAY_SECS + 1
					startTime = config.Time[0]
				}
			}
		}
	}
	return day, startTime
}

func (self *CsvMgr) GetSummonBoxConfig(configType int, configOrder int) *SummonBoxConfig {
	for _, v := range self.SummonBoxConfig {
		if v.Type == configType && v.Order == configOrder {
			return v
		}
	}
	return nil
}

// 读取列表
func (self *CsvMgr) LoadArenaRewardConfig() {
	GetCsvUtilMgr().LoadCsv("Arena_Reward", &self.ArenaRewardConfig)
	GetCsvUtilMgr().LoadCsv("Arena_Parameter", &self.ArenaParameter)
}

func (self *CsvMgr) GetArenaRewardConfig(configType int, rank int) *ArenaRewardConfig {
	for _, v := range self.ArenaRewardConfig {
		if v.Type == configType && rank <= v.Min && rank >= v.Max {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetArenaParameterConfig(index int64) *ArenaParameterConfig {
	for _, v := range self.ArenaParameter {
		if index == v.Id {
			return v
		}
	}
	return nil
}

/*
func (self *CsvMgr) GetAstrologyConfig(id int) *AstrologyDropConfig {
	for _, v := range self.AstrologyDropConfig {
		if v.
		if id == v.Id {
			return v
		}
	}
	return nil
}
*/

// 加载军团狩猎配置
func (self *CsvMgr) LoadUnionHunt() {
	GetCsvUtilMgr().LoadCsv("Guild_Boss", &self.UnionHuntConfig)

	GetCsvUtilMgr().LoadCsv("Guild_Bossdrop", &self.UnionHuntDropConfig)
}

// 加载赏金令配置
func (self *CsvMgr) LoadWarOrder() {
	tempWarOrderConfig := make([]*WarOrderConfig, 0)
	GetCsvUtilMgr().LoadCsv("Warorder_Level", &tempWarOrderConfig)
	self.WarOrderConfig = make(map[int]*WarOrderConfig)
	for _, v := range tempWarOrderConfig {
		self.WarOrderConfig[v.Id] = v
	}

	tempWarOrderParam := make([]*WarOrderParam, 0)
	GetCsvUtilMgr().LoadCsv("Warorder_param", &tempWarOrderParam)
	self.WarOrderParam = make(map[int]*WarOrderParam)
	for _, v := range tempWarOrderParam {
		self.WarOrderParam[v.Id] = v
	}

	self.WarOrderLimitConfig = make(map[int]*WarOrderLimitConfig)
	GetCsvUtilMgr().LoadCsv("Activity_Warorder", &self.WarOrderLimitConfig)
}

func (self *CsvMgr) LoadAccessCard() {
	self.AccessAwardConfig = make(map[int]*AccessAwardConfig)
	GetCsvUtilMgr().LoadCsv("Activity_AccessCard_TotleAward", &self.AccessAwardConfig)

	self.AccessTaskConfig = make(map[int]*AccessTaskConfig)
	GetCsvUtilMgr().LoadCsv("Activity_AccessCard", &self.AccessTaskConfig)

	GetCsvUtilMgr().LoadCsv("Activity_AccessCard_Rankaward", &self.AccessRankConfig)
}

func (self *CsvMgr) GetUnionHuntConfigByID(nID int) *UnionHuntConfig {

	configs := GetCsvMgr().UnionHuntConfig
	for _, v := range configs {
		if nID == v.Id {
			return v
		}
	}

	return nil
}

func (self *CsvMgr) GetUnionHuntDrop(player *Player, nGroup int, nDamage int64) []int {
	ret := []int{}

	configs := GetCsvMgr().UnionHuntDropConfig
	for _, v := range configs {
		if nGroup != v.Group {
			continue
		}

		if len(v.Drop) != len(v.Level) {
			continue
		}

		if nDamage < v.Hp {
			break
		}

		drop := 0
		for i, t := range v.Drop {
			if t > 0 {
				if v.Level[i] == 0 {
					drop = t
				} else if player.GetModule("pass").(*ModPass).GetPass(v.Level[i]) != nil {
					drop = t
				}
			}
		}

		if drop != 0 {
			ret = append(ret, drop)
		}
	}

	return ret
}

func (self *CsvMgr) GetUnionHuntDropConfig(nGroup int, nDamage int64) *UnionHuntDropConfig {
	var ret *UnionHuntDropConfig = nil
	configs := GetCsvMgr().UnionHuntDropConfig
	for _, v := range configs {
		if nGroup != v.Group {
			continue
		}

		if ret == nil {
			ret = v
		}

		if nDamage < v.Hp {
			break
		}

		ret = v
	}

	return ret
}

// 加载军团狩猎配置
func (self *CsvMgr) LoadArenaSpecialClass() {
	self.ArenaSpecialClassConfig = make([]*ArenaSpecialClass, 0)
	GetCsvUtilMgr().LoadCsv("Arena_Class", &self.ArenaSpecialClassConfig)

	self.ArenaSpecialClassMap = make(map[int][]*ArenaSpecialClass)
	for _, v := range self.ArenaSpecialClassConfig {
		_, ok := self.ArenaSpecialClassMap[v.Class]
		if !ok {
			self.ArenaSpecialClassMap[v.Class] = []*ArenaSpecialClass{}
		}
		self.ArenaSpecialClassMap[v.Class] = append(self.ArenaSpecialClassMap[v.Class], v)
	}
}

func (self *CsvMgr) GetArenaSpecialClassConfigByID(ID int) *ArenaSpecialClass {
	for _, v := range self.ArenaSpecialClassConfig {
		if v.Id == ID {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetArenaSpecialClassConfig(nClass int, nDan int) *ArenaSpecialClass {
	var ret *ArenaSpecialClass = nil
	configs, ok := GetCsvMgr().ArenaSpecialClassMap[nClass]
	if !ok || nil == configs {
		return ret
	}

	for _, v := range configs {
		if nDan != v.Dan {
			continue
		}
		return v
	}
	return ret
}

func (self *CsvMgr) GetWarOrderConfig(id int) *WarOrderConfig {
	for _, v := range self.WarOrderConfig {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetEquipRecastConfig(quality int, camp int) *EquipRecastConfig {
	for _, v := range self.EquipRecastConfig {
		if v.Quality == quality && v.Attribute == camp {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) CalNewPitMonsterLv(param int64) int {
	rel := 1
	for _, v := range self.HeroGrowthConfig {
		if v.GrowthValue[3] < param {
			rel = v.GrowthLevel
		} else {
			break
		}
	}
	return rel
}

// 20200521 引入虚空英雄修改替换算法
func (self *CsvMgr) MakeNewShop(shopType int, stage int, player *Player) []*JS_NewShopInfo {

	newShopList := make([]*JS_NewShopInfo, 0)
	config := GetCsvMgr().NewShopConfigMap[shopType]
	if config == nil {
		return newShopList
	}

	//20200521 add
	//是否有虚空英雄
	isHasCamp7 := false
	isHas7001 := false
	if player != nil {
		isHasCamp7, isHas7001 = player.GetModule("hero").(*ModHero).ShopJudge()
	}

	//第一次遍历记录每个格子的最大权值,方便扩展多个区域重叠
	wight := make(map[int]int)
	for _, value := range config {
		for _, v := range value {
			if (v.LevelLimit == LOGIC_FALSE || stage >= v.LevelLimit) && (v.LevelShield == LOGIC_FALSE || stage < v.LevelShield) {
				if v.Judge == 1 && isHasCamp7 {
					wight[v.Grid] += v.ReplaceWeight
				} else if v.Judge == 2 && isHas7001 {
					wight[v.Grid] += v.ReplaceWeight
				} else {
					wight[v.Grid] += v.WeightFunction
				}
			}
		}
	}
	//第二次遍历生成物品
	for grid, value := range wight {
		if value <= 0 {
			continue
		}
		configGrid := config[grid]
		if len(configGrid) <= 0 {
			continue
		}
		rand := HF_GetRandom(value)
		rate := 0
		for _, v := range configGrid {
			if (v.LevelLimit == LOGIC_FALSE || stage >= v.LevelLimit) && (v.LevelShield == LOGIC_FALSE || stage < v.LevelShield) {

				realGroup := v.Group
				realItem := v.ItemId
				realNum := v.ItemNumber
				realWeight := v.WeightFunction
				realCost := v.CostItems
				realCostNum := v.CostNums

				if (v.Judge == 1 && isHasCamp7) || (v.Judge == 2 && isHas7001) {
					realGroup = v.ReplaceGroup
					realItem = v.ReplaceItem
					realNum = v.ReplaceNum
					realWeight = v.ReplaceWeight
					realCost = v.ReplaceCost
					realCostNum = v.ReplaceCostNum
				}

				rate += realWeight
				//选中
				if rate > rand {
					var shopinfo JS_NewShopInfo
					shopinfo.Grid = v.Grid
					shopinfo.State = LOGIC_FALSE
					shopinfo.DisCount = self.GetDiscountConfig(shopType, stage, v.Grid) //折扣

					if realGroup == 1 {
						shopinfo.ItemId = realItem
						shopinfo.ItemNum = realNum
						for i := 0; i < len(realCost); i++ {
							if realCost[i] > 0 {
								shopinfo.CostId = append(shopinfo.CostId, realCost[i])
								shopinfo.CostNum = append(shopinfo.CostNum, (realCostNum[i]*shopinfo.DisCount)/PercentNum)
							}
						}
					} else if realGroup == 2 {
						shopinfo.ItemNum = realNum
						if self.EquipShopRate[realItem] <= 0 {
							return newShopList
						}
						rand := HF_GetRandom(self.EquipShopRate[realItem])
						rate := 0
						for _, equip := range self.EquipConfigMap {
							if equip.ShopClass != realItem {
								continue
							}
							rate += equip.ShopWeight
							isFind := false
							if rate > rand {
								shopinfo.ItemId = equip.EquipId
								for i := 0; i < len(realCost); i++ {
									if realCost[i] > 0 {
										isFind = true
										shopinfo.CostId = append(shopinfo.CostId, realCost[i])
										costIndex := realCostNum[i] - 1
										if costIndex < 0 {
											return newShopList
										}
										shopinfo.CostNum = append(shopinfo.CostNum, (equip.Price[costIndex]*shopinfo.DisCount)/PercentNum)
									}
								}
							}
							if isFind {
								break
							}
						}
					} else {
						return newShopList
					}
					newShopList = append(newShopList, &shopinfo)
					break
				}
			}
		}
	}

	return newShopList
}

//获得某个格子的物品
func (self *CsvMgr) GetGoodByGrid(shopType int, stage int, player *Player, grid int) *JS_NewShopInfo {

	config := GetCsvMgr().NewShopConfigMap[shopType]
	if config == nil {
		return nil
	}

	//是否有虚空英雄
	isHasCamp7 := false
	isHas7001 := false
	if player != nil {
		isHasCamp7, isHas7001 = player.GetModule("hero").(*ModHero).ShopJudge()
	}

	//第一次遍历记录每个格子的最大权值,方便扩展多个区域重叠
	wight := make(map[int]int)
	for _, value := range config {
		for _, v := range value {
			if v.Grid == grid && (v.LevelLimit == LOGIC_FALSE || stage >= v.LevelLimit) && (v.LevelShield == LOGIC_FALSE || stage < v.LevelShield) {
				if v.Judge == 1 && isHasCamp7 {
					wight[v.Grid] += v.ReplaceWeight
				} else if v.Judge == 2 && isHas7001 {
					wight[v.Grid] += v.ReplaceWeight
				} else {
					wight[v.Grid] += v.WeightFunction
				}
			}
		}
	}
	//第二次遍历生成物品
	for grid, value := range wight {
		if value <= 0 {
			continue
		}
		configGrid := config[grid]
		if len(configGrid) <= 0 {
			continue
		}
		rand := HF_GetRandom(value)
		rate := 0
		for _, v := range configGrid {
			if (v.LevelLimit == LOGIC_FALSE || stage >= v.LevelLimit) && (v.LevelShield == LOGIC_FALSE || stage < v.LevelShield) {

				realGroup := v.Group
				realItem := v.ItemId
				realNum := v.ItemNumber
				realWeight := v.WeightFunction
				realCost := v.CostItems
				realCostNum := v.CostNums

				if (v.Judge == 1 && isHasCamp7) || (v.Judge == 2 && isHas7001) {
					realGroup = v.ReplaceGroup
					realItem = v.ReplaceItem
					realNum = v.ReplaceNum
					realWeight = v.ReplaceWeight
					realCost = v.ReplaceCost
					realCostNum = v.ReplaceCostNum
				}

				rate += realWeight
				//选中
				if rate > rand {
					var shopinfo JS_NewShopInfo
					shopinfo.Grid = v.Grid
					shopinfo.State = LOGIC_FALSE
					shopinfo.DisCount = self.GetDiscountConfig(shopType, stage, v.Grid) //折扣

					if realGroup == 1 {
						shopinfo.ItemId = realItem
						shopinfo.ItemNum = realNum
						for i := 0; i < len(realCost); i++ {
							if realCost[i] > 0 {
								shopinfo.CostId = append(shopinfo.CostId, realCost[i])
								shopinfo.CostNum = append(shopinfo.CostNum, (realCostNum[i]*shopinfo.DisCount)/PercentNum)
							}
						}
					} else if realGroup == 2 {
						shopinfo.ItemNum = realNum
						if self.EquipShopRate[realItem] <= 0 {
							return nil
						}
						rand := HF_GetRandom(self.EquipShopRate[realItem])
						rate := 0
						for _, equip := range self.EquipConfigMap {
							if equip.ShopClass != realItem {
								continue
							}
							rate += equip.ShopWeight
							isFind := false
							if rate > rand {
								shopinfo.ItemId = equip.EquipId
								for i := 0; i < len(realCost); i++ {
									if realCost[i] > 0 {
										isFind = true
										shopinfo.CostId = append(shopinfo.CostId, realCost[i])
										costIndex := realCostNum[i] - 1
										if costIndex < 0 {
											return nil
										}
										shopinfo.CostNum = append(shopinfo.CostNum, (equip.Price[costIndex]*shopinfo.DisCount)/PercentNum)
									}
								}
							}
							if isFind {
								break
							}
						}
					} else {
						return nil
					}
					return &shopinfo
				}
			}
		}
	}

	return nil
}

func (self *CsvMgr) GetDiscountConfig(shopType int, stage int, grid int) int {
	//print(stage)
	//print("\n")
	//print(grid)
	//print("\n")
	for _, v := range self.NewShopDiscount {
		if v.Shop == shopType {
			for i := 0; i < len(v.Grid); i++ {
				if v.Grid[i] == grid {
					//第一次循环计算权重
					rateAll := 0
					for j := 0; j < len(v.LevelLimit); j++ {
						if v.Chance[j] <= 0 {
							continue
						}
						if (v.LevelLimit[j] == LOGIC_FALSE || stage >= v.LevelLimit[j]) && (v.LevelShield[j] == LOGIC_FALSE || stage < v.LevelShield[j]) {
							rateAll += v.Chance[j]
						}
					}
					print(rateAll)
					print("\n")
					if rateAll <= 0 {
						return 100
					}
					rand := HF_GetRandom(rateAll)
					rate := 0
					for j := 0; j < len(v.LevelLimit); j++ {
						if v.Chance[j] <= 0 {
							continue
						}
						if (v.LevelLimit[j] == LOGIC_FALSE || stage >= v.LevelLimit[j]) && (v.LevelShield[j] == LOGIC_FALSE || stage < v.LevelShield[j]) {
							rate += v.Chance[j]
							if rate > rand {
								if v.Discount[j] > 0 {
									//print(v.Discount[j])
									//print("-------------------\n")
									return v.Discount[j]
								} else {
									return 100
								}
							}
						}
					}
				}
			}
		}
	}

	return 100
}

func (self *CsvMgr) LoadHeroSkinConfig() {
	self.HeroSkinConfig = make([]*HeroSkin, 0)
	GetCsvUtilMgr().LoadCsv("Hero_Mod", &self.HeroSkinConfig)

	self.HeroGrowConfig = make([]*HeroGrowConfig, 0)
	GetCsvUtilMgr().LoadCsv("Activity_Herogrow", &self.HeroGrowConfig)
}

func (self *CsvMgr) LoadCrossArenaConfig() {
	self.CrossArenaRewardConfig = make([]*CrossArenaRewardConfig, 0)
	GetCsvUtilMgr().LoadCsv("Activity_Arenareward", &self.CrossArenaRewardConfig)

	self.CrossArenaSubsection = make([]*CrossArenaSubsection, 0)
	GetCsvUtilMgr().LoadCsv("Activity_Arenasubsection", &self.CrossArenaSubsection)

	self.CrossArena3V3RewardConfig = make([]*CrossArena3V3RewardConfig, 0)
	GetCsvUtilMgr().LoadCsv("Activity_Arenareward3v3", &self.CrossArena3V3RewardConfig)

	self.CrossArena3V3Subsection = make([]*CrossArena3V3Subsection, 0)
	GetCsvUtilMgr().LoadCsv("Activity_Arenasubsection3v3", &self.CrossArena3V3Subsection)
}

func (self *CsvMgr) GetHeroSkinConfig(id int) *HeroSkin {
	// 获得配置
	configs := GetCsvMgr().HeroSkinConfig
	for _, v := range configs {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) LoadSpecialPurchase() {
	GetCsvUtilMgr().LoadCsv("Activity_BuyLimitCondition", &self.ActivityBuyLimit)

	GetCsvUtilMgr().LoadCsv("Activity_BuyLimitItem", &self.ActivityBuyItem)

	self.ActivityMapBuyItem = make(map[int][]*ActivityBuyItem)
	for _, v := range self.ActivityBuyItem {
		_, ok := self.ActivityMapBuyItem[v.Group]
		if !ok {
			self.ActivityMapBuyItem[v.Group] = []*ActivityBuyItem{}
		}
		self.ActivityMapBuyItem[v.Group] = append(self.ActivityMapBuyItem[v.Group], v)
	}
}

func (self *CsvMgr) GetSpecialPurchaseConfig(nID int) *ActivityBuyLimit {

	configs := GetCsvMgr().ActivityBuyLimit

	for _, v := range configs {
		if nID == v.ID {
			return v
		}
	}

	return nil
}

func (self *CsvMgr) GetSpecialPurchaseItemConfig(nGiftID int) *ActivityBuyItem {

	configs := GetCsvMgr().ActivityBuyItem

	for _, v := range configs {
		if nGiftID == v.ID {
			return v
		}
	}

	return nil
}

//func (self *CsvMgr) LoadTreeConfig() {
//	GetCsvUtilMgr().LoadCsv("Tree_config", &self.TreeConfig)
//}
//
//func (self *CsvMgr) GetTreeConfig(quality int) *TreeConfig {
//	configs := GetCsvMgr().TreeConfig
//	for _, v := range configs {
//		if quality == v.Quality {
//			return v
//		}
//	}
//	return nil
//}

func (self *CsvMgr) LoadTreeLevelConfig() {
	GetCsvUtilMgr().LoadCsv("Tree_level", &self.TreeLevelConfig)
}

func (self *CsvMgr) LoadInterstellarConfig() {

	self.InterstellarTaskNode = make(map[int]*TaskNode)

	self.InterstellarConfig = make(map[int]*InterstellarConfig)
	GetCsvUtilMgr().LoadCsv("Interstellar_Config", &self.InterstellarConfig)
	for _, v := range self.InterstellarConfig {
		for i := 0; i < len(v.TaskTypes); i++ {
			if v.TaskTypes[i] == 0 {
				continue
			}
			taskNode := new(TaskNode)
			taskNode.Id = v.Id*100 + i + 1
			taskNode.Tasktypes = v.TaskTypes[i]
			taskNode.N1 = v.N[i]
			taskNode.N2 = v.M[i]
			taskNode.N3 = v.P[i]
			taskNode.N4 = v.Q[i]
			self.InterstellarTaskNode[taskNode.Id] = taskNode
		}
	}

	self.InterstellarHangup = make(map[int]*InterstellarHangup)
	GetCsvUtilMgr().LoadCsv("Interstellar_Hangup", &self.InterstellarHangup)

	InterstellarWarTemp := make([]*InterstellarWar, 0)
	GetCsvUtilMgr().LoadCsv("Interstellar_War", &InterstellarWarTemp)
	self.InterstellarWar = make(map[int]map[int]*InterstellarWar)
	for _, v := range InterstellarWarTemp {
		_, ok := self.InterstellarWar[v.Nebula]
		if !ok {
			self.InterstellarWar[v.Nebula] = make(map[int]*InterstellarWar)
		}
		self.InterstellarWar[v.Nebula][v.Id] = v

		taskNode := new(TaskNode)
		taskNode.Id = v.Id
		taskNode.Tasktypes = v.TaskTypes
		taskNode.N1 = v.Ns[0]
		taskNode.N2 = v.Ns[1]
		taskNode.N3 = v.Ns[2]
		taskNode.N4 = v.Ns[3]
		self.InterstellarTaskNode[taskNode.Id] = taskNode
	}

	InterstellarBoxTemp := make([]*InterstellarBox, 0)
	GetCsvUtilMgr().LoadCsv("Interstellar_Box", &InterstellarBoxTemp)
	self.InterstellarBox = make(map[int]map[int]*InterstellarBox)
	for _, v := range InterstellarBoxTemp {
		_, ok := self.InterstellarBox[v.GroupId]
		if !ok {
			self.InterstellarBox[v.GroupId] = make(map[int]*InterstellarBox)
		}
		self.InterstellarBox[v.GroupId][v.BoxId] = v
	}

}

func (self *CsvMgr) LoadActivityBossConfig() {
	ActivityBossRankConfigTemp := make([]*ActivityBossRankConfig, 0)
	GetCsvUtilMgr().LoadCsv("Activity_Bossrank", &ActivityBossRankConfigTemp)
	self.ActivityBossRankConfig = make(map[int]map[int][]*ActivityBossRankConfig)
	for _, v := range ActivityBossRankConfigTemp {
		_, ok := self.ActivityBossRankConfig[v.ActivityType]
		if !ok {
			self.ActivityBossRankConfig[v.ActivityType] = make(map[int][]*ActivityBossRankConfig)
		}
		self.ActivityBossRankConfig[v.ActivityType][v.ActivityPeriods] = append(self.ActivityBossRankConfig[v.ActivityType][v.ActivityPeriods], v)
	}

	GetCsvUtilMgr().LoadCsv("Activity_Boss", &self.ActivityBossConfig)

	ActivityBossTargetConfigTemp := make([]*ActivityBossTargetConfig, 0)
	GetCsvUtilMgr().LoadCsv("Activity_Bosstarget", &ActivityBossTargetConfigTemp)
	self.ActivityBossTargetConfig = make(map[int]map[int][]*ActivityBossTargetConfig)
	for _, v := range ActivityBossTargetConfigTemp {
		_, ok := self.ActivityBossTargetConfig[v.ActivityType]
		if !ok {
			self.ActivityBossTargetConfig[v.ActivityType] = make(map[int][]*ActivityBossTargetConfig)
		}
		self.ActivityBossTargetConfig[v.ActivityType][v.ActivityPeriods] = append(self.ActivityBossTargetConfig[v.ActivityType][v.ActivityPeriods], v)
	}

	ActivityBossExchangeConfigTemp := make([]*ActivityBossExchangeConfig, 0)
	GetCsvUtilMgr().LoadCsv("Activity_Bossexchange", &ActivityBossExchangeConfigTemp)
	self.ActivityBossExchangeConfig = make(map[int]map[int][]*ActivityBossExchangeConfig)
	for _, v := range ActivityBossExchangeConfigTemp {
		_, ok := self.ActivityBossExchangeConfig[v.ActivityType]
		if !ok {
			self.ActivityBossExchangeConfig[v.ActivityType] = make(map[int][]*ActivityBossExchangeConfig)
		}
		self.ActivityBossExchangeConfig[v.ActivityType][v.ActivityPeriods] = append(self.ActivityBossExchangeConfig[v.ActivityType][v.ActivityPeriods], v)
	}
}

func (self *CsvMgr) GetTreeLevelConfig(level int) *TreeLevel {
	configs := GetCsvMgr().TreeLevelConfig
	for _, v := range configs {
		if level == v.TreeLevel {
			return v
		}
	}
	return nil
}
func (self *CsvMgr) LoadTreeProfessionalConfig() {
	GetCsvUtilMgr().LoadCsv("Tree_professional", &self.TreeProfessionalConfig)

	self.TreeProfessionalMapConfig = make(map[int][]*TreeProfessional)
	for _, v := range self.TreeProfessionalConfig {
		_, ok := self.TreeProfessionalMapConfig[v.Type]
		if !ok {
			self.TreeProfessionalMapConfig[v.Type] = []*TreeProfessional{}
		}
		self.TreeProfessionalMapConfig[v.Type] = append(self.TreeProfessionalMapConfig[v.Type], v)
	}
}
func (self *CsvMgr) GetTreeProfessionalConfig(nType, level int) *TreeProfessional {
	configs, ok := GetCsvMgr().TreeProfessionalMapConfig[nType]
	if !ok {
		return nil
	}

	for _, v := range configs {
		if level == v.Level {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetTeamAttrConfig(nType, camp1, camp2 int) *TeamAttrConfig {
	for _, v := range self.TeamAttrConfig {
		if v.Type == nType && v.Camp[0] == camp1 && v.Camp[1] == camp2 {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetInstanceBoxConfig(id int, event int) *InstanceBox {
	value, ok := self.InstanceBox[id]
	if ok {
		for _, v := range value {
			if v.BoxId == event {
				return v
			}
		}
	}
	return nil
}

func (self *CsvMgr) GetSevendayAward(id int, stage int) *SevendayAward {
	for _, v := range self.SevendayAward {
		if v.Id == id && v.Stage == stage {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetOverflowConfig(id int, group int) *ActivityOverflowGifts {
	for _, v := range self.ActivityOverflowGifts {
		if v.Id == id && v.Group == group {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetNewPitDifficulty(difficulty int, numberofplies int) *NewPitDifficulty {
	for _, v := range self.NewPitDifficulty {
		if v.Difficulty == difficulty && v.NumberOfPlies == numberofplies {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetBadgeTaskConfig(money int) *BadgeTaskConfig {
	for _, v := range self.BadgeTaskConfig {
		if v.Money == money {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetAccessTaskConfig(group int, rank int) *AccessRankConfig {
	for _, v := range self.AccessRankConfig {
		if v.Group == group && rank >= v.RankMin && rank <= v.RankMax {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetActivityBossRankConfig(id int, period int, subsection int, ranking int) *ActivityBossRankConfig {

	_, ok := self.ActivityBossRankConfig[id]
	if !ok {
		return nil
	}

	for _, v := range self.ActivityBossRankConfig[id][period] {
		if v.Subsection == subsection && v.Ranking == ranking {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetActivityBossConfig(activityId int, period int) *ActivityBossConfig {
	for _, v := range self.ActivityBossConfig {
		if v.ActivityType == activityId && v.ActivityPeriods == period {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetActivityBossTargetConfig(id int, period int, taskId int) *ActivityBossTargetConfig {

	_, ok := self.ActivityBossTargetConfig[id]
	if !ok {
		return nil
	}

	for _, v := range self.ActivityBossTargetConfig[id][period] {
		if v.TaskId == taskId {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetActivityBossExchangeConfig(id int, period int, exchangeId int) *ActivityBossExchangeConfig {

	_, ok := self.ActivityBossExchangeConfig[id]
	if !ok {
		return nil
	}

	for _, v := range self.ActivityBossExchangeConfig[id][period] {
		if v.Id == exchangeId {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetHeroGrowConfig(id int) *HeroGrowConfig {
	// 获得配置
	configs := GetCsvMgr().HeroGrowConfig
	for _, v := range configs {
		if v.Id == id {
			return v
		}
	}
	return nil
}

//获得大段位奖励
func (self *CsvMgr) GetCrossArenaRewardBySubsection(subsection int) map[int]*Item {
	data := make(map[int]*Item)
	if subsection == 0 {
		return data
	}
	for _, v := range self.CrossArenaRewardConfig {
		if v.Type != 1 {
			continue
		}
		if v.Subsection < subsection {
			continue
		}
		AddItemMapHelper(data, v.Item, v.Num)
	}
	return data
}

//获得至尊奖励
func (self *CsvMgr) GetCrossArenaRewardByRank(subsection int, class int) *CrossArenaRewardConfig {
	for _, v := range self.CrossArenaRewardConfig {
		if v.Type != 2 {
			continue
		}
		if v.Subsection == subsection && v.Class == class {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetCrossArenaSubsection(id int) *CrossArenaSubsection {
	for _, v := range self.CrossArenaSubsection {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetCrossArenaSubsectionName(sub int, class int) string {
	for _, v := range self.CrossArenaSubsection {
		if v.Subsection == sub && v.Class == class {
			return v.Name
		}
	}
	return ""
}

func (self *CsvMgr) LoadStageTalentConfig() {
	GetCsvUtilMgr().LoadCsv("Hero_StageTalent", &self.StageTalentConfig)

	self.StageTalentMap = make(map[int][]*StageTalentConfig)
	for _, value := range self.StageTalentConfig {
		_, ok := self.StageTalentMap[value.Group]
		if !ok {
			self.StageTalentMap[value.Group] = make([]*StageTalentConfig, 0)
		}
		self.StageTalentMap[value.Group] = append(self.StageTalentMap[value.Group], value)
	}
}

// 获得天赋
func (self *CsvMgr) GetStageTalent(id int) *StageTalentConfig {
	for _, v := range self.StageTalentConfig {
		if v.ID == id {
			return v
		}
	}
	return nil
}

// 获得天赋组
func (self *CsvMgr) GetStageTalentMap(group int) []*StageTalentConfig {
	_, ok := self.StageTalentMap[group]
	if !ok {
		return nil
	}
	return self.StageTalentMap[group]
}

func (self *CsvMgr) LotteryDrawGetWeight(id int, value int) int {
	config, ok := GetCsvMgr().LotteryDrawConfigMap[id]
	if !ok {
		return 0
	}
	if value >= config.MaxLucky {
		return config.MaxWeight
	}
	if value >= config.MinLucky {
		return config.MinWeight
	}
	return config.Weight
}

func (self *CsvMgr) LotteryDrawGetLucky(id int) int {
	config, ok := GetCsvMgr().LotteryDrawConfigMap[id]
	if !ok {
		return 0
	}
	return config.Lucky
}

func (self *CsvMgr) GetWorldLevelConfig(level int) *WorldLvLevelConfig {
	for _, v := range self.WorldLvLevelConfig {
		if v.WorldLv == level {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetWorldLevelTypeConfig(heroId int, level int) *WorldLvTpyeConfig {

	for _, v := range self.WorldLvTpyeConfig {
		if v.HeroNpc != heroId {
			continue
		}
		if level <= v.WorldLv {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetWorldShowBossConfig(level int) *WorldShowBossConfig {
	for _, v := range self.WorldShowBossConfig {
		if level <= v.Crystallv {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetHonourShopConfig(parent_label int, subtab int, grid int, stage int) *HonourShopConfig {
	for _, v := range self.HonourShopConfigMap {
		if v.ParentLabel != parent_label || v.SubTab != subtab || v.Grid != grid {
			continue
		}
		if v.LevelStart != 0 && stage < v.LevelStart {
			continue
		}
		if v.LevelEnd != 0 && stage >= v.LevelEnd {
			continue
		}
		return v
	}
	return nil
}

//获得大段位奖励
func (self *CsvMgr) GetCrossArena3V3RewardBySubsection(subsection int) map[int]*Item {
	data := make(map[int]*Item)
	if subsection == 0 {
		return data
	}
	for _, v := range self.CrossArenaRewardConfig {
		if v.Type != 1 {
			continue
		}
		if v.Subsection < subsection {
			continue
		}
		AddItemMapHelper(data, v.Item, v.Num)
	}
	return data
}

//获得至尊奖励
func (self *CsvMgr) GetCrossArena3V3RewardByRank(subsection int, class int) *CrossArena3V3RewardConfig {
	for _, v := range self.CrossArena3V3RewardConfig {
		if v.Type != 2 {
			continue
		}
		if v.Subsection == subsection && v.Class == class {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetCrossArena3V3Subsection(id int) *CrossArena3V3Subsection {
	for _, v := range self.CrossArena3V3Subsection {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (self *CsvMgr) GetCrossArena3V3SubsectionName(sub int, class int) string {
	for _, v := range self.CrossArena3V3Subsection {
		if v.Subsection == sub && v.Class == class {
			return v.Name
		}
	}
	return ""
}

func (self *CsvMgr) LoadRankReward() {
	RankRewardConfigTemp := make([]*RankRewardConfig, 0)
	GetCsvUtilMgr().LoadCsv("RankReward", &RankRewardConfigTemp)
	self.RankRewardConfigMap = make(map[int]map[int][]*RankRewardConfig)
	for _, v := range RankRewardConfigTemp {
		_, ok := self.RankRewardConfigMap[v.ActivityType]
		if !ok {
			self.RankRewardConfigMap[v.ActivityType] = make(map[int][]*RankRewardConfig)
		}
		self.RankRewardConfigMap[v.ActivityType][v.Group] = append(self.RankRewardConfigMap[v.ActivityType][v.Group], v)
	}
}

func (self *CsvMgr) GetRankRewardConfig(id int, keyId int, rank int) *RankRewardConfig {
	group := keyId / 1000
	_, ok := self.RankRewardConfigMap[id]
	if ok {
		for _, v := range self.RankRewardConfigMap[id][group] {
			if rank >= v.RankHigh && rank <= v.RankLow {
				return v
			}
		}
	}
	return nil
}

func (self *CsvMgr) GetRankRewardConfig2C(id int, keyId int) []*RankRewardConfig {
	group := keyId / 1000
	_, ok := self.RankRewardConfigMap[id]
	if ok {
		return self.RankRewardConfigMap[id][group]
	}
	return nil
}
