package game

// 商店模块
import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	RESET_TIME_TYPE_DIS = 34 //商店重置时间偏移值
	TARIFF_TYPE_DIS     = 7  //刷新价格配置
)

// 商店数据库
type San_Shop struct {
	Uid           int64
	Shoptype      int    // 商店类型
	Shopgood      string // 道具购买状态
	Refindex      int    // 刷新次数
	Shopnextgood  string // 下次需要随机出来的道具
	Todayrefcount int    // 每天购买的次数
	Sysreftime    int64  // 系统刷新时间
	Lastupdtime   int64  // 上次刷新时间

	shopgood     []*JS_NewShopInfo
	shopnextgood []JS_ShopInfo

	DataUpdate
}

// 商店物品信息
type JS_ShopInfo struct {
	Id    int    `json:"id"`
	Isbuy int    `json:"isbuy"`
	Info  [7]int `json:"info"`
}

type JS_NewShopInfo struct {
	Grid     int   `json:"grid"`     //格子
	State    int   `json:"state"`    //0未买  1已买
	ItemId   int   `json:"itemid"`   //商品ID
	ItemNum  int   `json:"itemnum"`  //商品数量
	CostId   []int `json:"costid"`   //消耗ID数组
	CostNum  []int `json:"costnum"`  //消耗数量数组 价格已计算折扣
	DisCount int   `json:"discount"` //折扣，显示用
}

//! 商店物品信息
//!
type JS_ShopItemInfo struct {
	Data [5]int `json:"data"`
}

// 商店
type ModShop struct {
	player   *Player
	Sql_Shop map[int]*San_Shop // 所有商店的结构
}

// 初始化每个商店的数据表结构, 参数是商店类型
func (self *ModShop) InitData(shoptype int) {
	shop := new(San_Shop)
	// 初始化商店数据
	self.Sql_Shop[shoptype] = shop
	// 查询商店数据
	sql := fmt.Sprintf("select * from `san_shop%d` where `uid` = %d", shoptype, self.player.ID)
	GetServer().DBUser.GetOneData(sql, shop, fmt.Sprintf("san_shop%d", shoptype), self.player.ID)
	// 初始化玩家数据
	// 新玩家
	if shop.Uid <= 0 {
		shop.Uid = self.player.ID
		shop.Shoptype = shoptype
		// 计算下一次系统的刷新时间
		self.RefreshShop(shoptype, false)
		self.Encode()
		InsertTable(fmt.Sprintf("san_shop%d", shoptype), shop, 0, true)
	} else { // 老玩家
		self.Decode(shoptype)
	}

	shop.Init(fmt.Sprintf("san_shop%d", shoptype), shop, true)
}

func (self *ModShop) OnGetData(player *Player) {
	self.player = player
	self.Sql_Shop = make(map[int]*San_Shop)
}

func (self *ModShop) OnGetOtherData() {
	self.InitData(SHOP_NEW_NORMAL)
	self.InitData(SHOP_NEW_UNION)
	self.InitData(SHOP_NEW_FIRE)
	self.InitData(SHOP_NEW_PIT)
	self.InitData(SHOP_NEW_PVP)
	self.InitData(SHOP_OLD_CONSUMERTOP)
}

func (self *ModShop) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "newshopbuy":
		var c2s_msg C2S_NewShopBuy
		json.Unmarshal(body, &c2s_msg)
		self.NewShopBuy(c2s_msg.Shoptype, c2s_msg.Grid)
		return true
	case "newshoprefresh":
		var c2s_msg C2S_NewShoprefresh
		json.Unmarshal(body, &c2s_msg)
		self.NewShopRefresh(c2s_msg.Shoptype)
		return true
	case "buymagicalhero":
		var c2s_msg C2S_ShopBuy
		json.Unmarshal(body, &c2s_msg)

		var msg S2C_TreasuryBuy
		msg.Cid = "magicalhero2buy"
		msg.Type = c2s_msg.Shoptype
		msg.Ret, msg.Item, msg.Info = self.BuyMagicalHero(c2s_msg.Grid, c2s_msg.Refindex)
		self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
		return true
	}

	return false
}

func (self *ModShop) OnSave(sql bool) {
	self.Encode()
	for _, value := range self.Sql_Shop {
		value.Update(sql)
	}
}

func (self *ModShop) Decode(shoptype int) { // 将数据库数据写入data
	value, ok := self.Sql_Shop[shoptype]
	if ok {
		json.Unmarshal([]byte(value.Shopgood), &value.shopgood)
		json.Unmarshal([]byte(value.Shopnextgood), &value.shopnextgood)
	}
}

func (self *ModShop) Encode() { // 将data数据写入数据库
	for _, value := range self.Sql_Shop {
		s, _ := json.Marshal(&value.shopgood)
		value.Shopgood = string(s)

		s, _ = json.Marshal(&value.shopnextgood)
		value.Shopnextgood = string(s)
	}
}

// 根据商店类型获取商店刷新的下一次时间
func (self *ModShop) GetNextTime(shoptype int) int64 {
	if shoptype == SHOP_OLD_CONSUMERTOP {
		now := TimeServer()
		h := now.Hour()
		if h < 5 {
			return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
		} else {
			return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix() + DAY_SECS
		}
	}
	resetType := shoptype + RESET_TIME_TYPE_DIS
	_, endTime := GetCsvMgr().GetNowStartAndEndByRoleDays(resetType, self.player.Sql_UserBase.Regtime)
	return endTime
}

// 获取下一次的周刷新时间
func (self *ModShop) GetNextWeekTime(shoptype int) int64 {
	now := TimeServer()
	h := now.Hour()
	if shoptype == SHOP_OLD_CONSUMERTOP {
		if h < 5 {
			return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
		} else {
			return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix() + DAY_SECS
		}
	}
	return 0
}

// 刷新商店,参数是商店类型,是否是登录时刷新
func (self *ModShop) RefreshShop(shopType int, login bool) bool {

	//特殊商店不走shop生成
	if shopType == SHOP_OLD_CONSUMERTOP {
		self.RefreshMagicalHero(login)
		return true
	}
	// 检查类型
	shop, ok := self.Sql_Shop[shopType]
	if !ok {
		return false
	}

	//先拿到用户的关卡进度
	stage, _ := GetOfflineInfoMgr().GetBaseInfo(self.player.Sql_UserBase.Uid)
	shop.Sysreftime = self.GetNextTime(shopType)
	shop.shopgood = GetCsvMgr().MakeNewShop(shopType, stage, self.player)
	if len(shop.shopgood) <= 0 {
		return false
	}
	return true
}

// 获得某个格子的刷新
func (self *ModShop) GetShopGood(shopType int, grid int) *JS_NewShopInfo {
	// 检查类型
	_, ok := self.Sql_Shop[shopType]
	if !ok {
		return nil
	}

	//先拿到用户的关卡进度
	stage, _ := GetOfflineInfoMgr().GetBaseInfo(self.player.Sql_UserBase.Uid)

	return GetCsvMgr().GetGoodByGrid(shopType, stage, self.player, grid)
}

func (self *ModShop) Refresh() {
	for _, value := range self.Sql_Shop {
		value.Refindex = 0
		if TimeServer().Unix() >= value.Sysreftime {
			self.RefreshShop(value.Shoptype, false)
		}
	}
}

func (self *ModShop) NewShopBuy(shoptype int, grid int) {

	// 先尝试刷新
	self.Refresh()
	shop, ok := self.Sql_Shop[shoptype]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_NOT_EXIST"))
		return
	}

	if grid <= 0 || grid > len(shop.shopgood) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_GRID_ERROR"))
		return
	}

	for i := 0; i < len(shop.shopgood); i++ {
		if shop.shopgood[i].Grid == grid {
			if shop.shopgood[i].State == LOGIC_TRUE {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_ALREADY_BUY"))
				return
			}

			if err := self.player.HasObjectOk(shop.shopgood[i].CostId, shop.shopgood[i].CostNum); err != nil {
				self.player.SendErrInfo("err", err.Error())
				return
			}

			costItem := self.player.RemoveObjectLst(shop.shopgood[i].CostId, shop.shopgood[i].CostNum, "商店购买", shoptype, 0, 1)
			num := GetGemNum(costItem)
			param3 := 0
			if num > 0 {
				param3 = -1
			}
			getItems := self.player.AddObjectSimple(shop.shopgood[i].ItemId, shop.shopgood[i].ItemNum, "商店购买", shoptype, 0, param3)
			//CheckAddItemLog(self.player, "商店购买", costItem, getItems)

			if num > 0 {
				AddSpecialSdkItemListLog(self.player, num, getItems, "商店购买")
			}

			//旧商店是值拷贝，这个地方需要整体引用
			self.Sql_Shop[shoptype].shopgood[i].State = LOGIC_TRUE

			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SHOP_BUY, shop.shopgood[i].ItemId, shoptype, shop.shopgood[i].ItemNum, "商店购买", 0, 0, self.player)

			var msgRel S2C_NewShopBuy
			msgRel.Cid = "newshopbuy"
			msgRel.GetItems = getItems
			msgRel.CostItems = costItem
			msgRel.NewShopInfo = shop.shopgood[i]
			self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
			self.player.HandleTask(TASK_TYPE_SHOP_BUY_COUNT, 1, shoptype, 0)
			self.player.GetModule("task").(*ModTask).SendUpdate()

			for i := 0; i < len(msgRel.GetItems); i++ {
				GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SHOP_BUY_GOODS, msgRel.GetItems[i].ItemID, shoptype, 0, "商店购买道具", 0, 0, self.player)
			}
			return
		}
	}

	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_GOODS_NOT_FIND"))
	return
}

func (self *ModShop) NewShopRefresh(shoptype int) {

	shop, ok := self.Sql_Shop[shoptype]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_NOT_EXIST"))
		return
	}

	//看消耗够不够
	realType := TARIFF_TYPE_DIS + shoptype
	configCost := GetCsvMgr().GetTariffConfig(realType, shop.Refindex+1)
	if configCost == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CONFIG_NOT_EXIST"))
		return
	}

	//看消耗够不够
	if err := self.player.HasObjectOk(configCost.ItemIds, configCost.ItemNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	costItem := self.player.RemoveObjectLst(configCost.ItemIds, configCost.ItemNums, "商店刷新", shoptype, 0, 0)
	self.RefreshShop(shoptype, false)
	shop.Refindex++
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SHOP_REFRESH, shop.Refindex, shoptype, 0, "刷新商店商品", 0, 0, self.player)

	var msgRel S2C_NewShopRefresh
	msgRel.Cid = "newshoprefresh"
	msgRel.CostItems = costItem
	msgRel.Info = new(Son_ShopInfo)
	msgRel.Info.Lastupdtime = shop.Lastupdtime
	msgRel.Info.Refindex = shop.Refindex
	msgRel.Info.Shopgood = shop.shopgood
	if len(msgRel.Info.Shopgood) <= 0 {
		msgRel.Info.Shopgood = []*JS_NewShopInfo{}
	}
	msgRel.Info.Shoptype = shop.Shoptype
	msgRel.Info.Sysreftime = shop.Sysreftime
	msgRel.Info.Uid = shop.Uid

	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SHOP_REFRESH, 1, shoptype, 0, "商店刷新", 0, 0, self.player)
	return
}

func (self *ModShop) SendInfo() {
	var msg S2C_ShopInfo
	msg.Cid = "getshop2lst"
	for _, value := range self.Sql_Shop {
		var shop Son_ShopInfo
		shop.Lastupdtime = value.Lastupdtime
		shop.Refindex = value.Refindex
		shop.Shopgood = value.shopgood
		if len(value.shopgood) <= 0 {
			shop.Shopgood = []*JS_NewShopInfo{}
		}
		shop.Shopnextgood = value.shopnextgood
		if len(value.shopnextgood) <= 0 {
			shop.Shopnextgood = []JS_ShopInfo{}
		}
		shop.Shoptype = value.Shoptype
		shop.Sysreftime = value.Sysreftime
		shop.Todayrefcount = value.Todayrefcount
		shop.Uid = value.Uid
		msg.Info = append(msg.Info, shop)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("getshop2lst", smsg)
}

func (self *ModShop) CheckHeroItem(heroId int) {
	heroItem := 12000000 + (heroId * 100) + 1

	for _, shop := range self.Sql_Shop {
		for index, shopgood := range shop.shopgood {
			if shopgood.ItemId == heroItem {
				shop.shopgood[index] = self.GetShopGood(shop.Shoptype, shopgood.Grid)
			}
		}
	}

	self.SendInfo()
}

// 刷新无双商店-每天更新
func (self *ModShop) RefreshMagicalHero(login bool) {
	shop, ok := self.Sql_Shop[SHOP_OLD_CONSUMERTOP]
	if !ok {
		return
	}

	// 没有数据直接初始化
	if len(shop.shopnextgood) == 0 {
		for _, value := range GetCsvMgr().Data["Consumetop_Shop"] {
			var shopinfo JS_ShopInfo
			shopinfo.Id = HF_Atoi(value["id"])
			shopinfo.Isbuy = 0
			shopinfo.Info = [7]int{
				HF_Atoi(value["item"]),
				HF_Atoi(value["num"]),
				HF_Atoi(value["type"]),
				HF_Atoi(value["cost"]),
				HF_Atoi(value["costnum"]),
				HF_Atoi(value["type"]),
				HF_Atoi(value["time"]),
			}

			shop.shopnextgood = append(shop.shopnextgood, shopinfo)
		}
		shop.Refindex = 0
		shop.Sysreftime = self.GetNextTime(SHOP_OLD_CONSUMERTOP)
		shop.Todayrefcount = int(self.GetNextWeekTime(SHOP_OLD_CONSUMERTOP))
	} else {
		if login { // 登录重置，添加不存在的道具
			for _, value := range GetCsvMgr().Data["Consumetop_Shop"] {
				find := false
				for k := 0; k < len(shop.shopnextgood); k++ {
					if shop.shopnextgood[k].Id == HF_Atoi(value["id"]) {
						find = true
						break
					}
				}
				if find == false {
					var shopinfo JS_ShopInfo
					shopinfo.Id = HF_Atoi(value["id"])
					shopinfo.Isbuy = 0
					shopinfo.Info = [7]int{
						HF_Atoi(value["item"]),
						HF_Atoi(value["num"]),
						HF_Atoi(value["type"]),
						HF_Atoi(value["cost"]),
						HF_Atoi(value["costnum"]),
						HF_Atoi(value["type"]),
						HF_Atoi(value["time"]),
					}
					shop.shopnextgood = append(shop.shopnextgood, shopinfo)
				}
			}
		}

		if TimeServer().Unix() < shop.Sysreftime {
			return
		}

		week_update := false
		if shop.Todayrefcount == 0 {
			shop.Todayrefcount = int(self.GetNextWeekTime(SHOP_OLD_CONSUMERTOP))
		}

		if TimeServer().Unix() >= int64(shop.Todayrefcount) {
			shop.Todayrefcount = int(self.GetNextWeekTime(SHOP_OLD_CONSUMERTOP))
			week_update = true
		}

		for i := 0; i < len(shop.shopnextgood); i++ {
			if shop.shopnextgood[i].Id > 0 {
				csv, ok := GetCsvMgr().Data["Consumetop_Shop"][shop.shopnextgood[i].Id]
				if !ok {
					continue
				}

				if csv["type"] == "2" {
					shop.shopnextgood[i].Isbuy = 0
					shop.shopnextgood[i].Info = [7]int{
						HF_Atoi(csv["item"]),
						HF_Atoi(csv["num"]),
						HF_Atoi(csv["type"]),
						HF_Atoi(csv["cost"]),
						HF_Atoi(csv["costnum"]),
						HF_Atoi(csv["type"]),
						HF_Atoi(csv["time"]),
					}
					shop.Sysreftime = self.GetNextTime(SHOP_OLD_CONSUMERTOP)
					shop.Refindex++
				} else if csv["type"] == "3" {
					if week_update {
						shop.Refindex = 0
						shop.shopnextgood[i].Isbuy = 0
						shop.shopnextgood[i].Info = [7]int{
							HF_Atoi(csv["item"]),
							HF_Atoi(csv["num"]),
							HF_Atoi(csv["type"]),
							HF_Atoi(csv["cost"]),
							HF_Atoi(csv["costnum"]),
							HF_Atoi(csv["type"]),
							HF_Atoi(csv["time"]),
						}
					}
				}
			}
		}
	}
	self.SendRefreshInfo(SHOP_OLD_CONSUMERTOP)
}

func (self *ModShop) BuyMagicalHero(grid int, num int) (int, []PassItem, *JS_ShopInfo) {
	if num < 0 {
		return 0, []PassItem{}, nil
	}
	// 尝试刷新
	self.RefreshMagicalHero(false)

	item := make([]PassItem, 0)
	shop, ok := self.Sql_Shop[SHOP_OLD_CONSUMERTOP]
	if !ok {
		return 3, item, nil
	}

	csv_shop, ok := GetCsvMgr().Data["Consumetop_Shop"][grid]
	if !ok {
		return 3, item, nil
	}

	config := GetCsvMgr().GetItemConfig(HF_Atoi(csv_shop["item"]))
	if config == nil {
		return 3, item, nil
	}

	index := 0
	for i := 0; i < len(shop.shopnextgood); i++ {
		if shop.shopnextgood[i].Id == grid {
			index = i
			break
		}
	}

	if csv_shop["type"] == "2" || csv_shop["type"] == "3" {
		if shop.shopnextgood[index].Isbuy+num > HF_Atoi(csv_shop["time"]) {
			return 3, item, nil
		}
	}

	shop.shopnextgood[index].Isbuy += num
	need := HF_Atoi(csv_shop["costnum"]) * num
	if self.player.GetObjectNum(MONEY_MHERO) < need {
		return 3, item, &shop.shopnextgood[index]
	}

	costitemid := HF_Atoi(csv_shop["cost"])
	self.player.AddObject(costitemid, -need, SHOP_MAGICALSHOP, 0, 0, "商店购买")
	self.player.HandleTask(CostEnergeTask, need, 0, costitemid)

	item = append(item, PassItem{ItemID: costitemid, Num: -need})
	item = append(item, PassItem{ItemID: HF_Atoi(csv_shop["item"]), Num: HF_Atoi(csv_shop["num"]) * num})

	self.player.AddObject(HF_Atoi(csv_shop["item"]), HF_Atoi(csv_shop["num"])*num, SHOP_MAGICALSHOP, 0, 0, "商店购买")

	costitem_csv := GetCsvMgr().GetItemConfig(costitemid)
	itemId := HF_Atoi(csv_shop["item"])
	getitem_csv := GetCsvMgr().GetItemConfig(itemId)

	if costitem_csv != nil && getitem_csv != nil {
		GetServer().sendLog_BuyItem(self.player, need, "商店购买", getitem_csv.ItemId, getitem_csv.ItemName,
			costitem_csv.ItemName, HF_Atoi(csv_shop["itemnumber"]), "商店购买")
	}

	// 任务
	self.player.HandleTask(ShopTask, 6, 0, 0)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SHOP_BUY, itemId, SHOP_MAGICALSHOP, 0, "商店购买", 0, 0, self.player)

	return 0, item, &shop.shopnextgood[index]
}

func (self *ModShop) SendRefreshInfo(shoptype int) {
	var msg S2C_SendRefreshInfo
	msg.Cid = "sendrefreshinfo"
	msg.Type = shoptype
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}
