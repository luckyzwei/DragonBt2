package game

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	//"time"
)

const (
	ITEM_TYPE_HERO       = 1  // 武将， itemsubtype=1, heroId := (itemid - 11000000) / 100
	ITEM_TYPE_SPECIAL    = 7  // 特殊道具, itemsubtype=6 英雄经验药水,0 技能值,2 体力,4 魂器,12 特殊掉落 , 8 宝箱
	ITEM_TYPE_MONEY      = 9  // 货币类型
	ITEM_TYPE_HORSE      = 13 // 战马
	ITEM_TYPE_HORSE_SOUL = 14 // 马魂
	ITEM_TYPE_MONTH_CARD = 15 // 月卡 1: //! 超值月卡 2: //! 至尊月卡
	ITEM_TYPE_TIGER      = 16 // 虎符
	ITEM_TYPE_RED_POCKET = 18 // 使用红包道具
	ITEM_TYPE_TREASURE   = 19 // 宝物
	//ITEM_TYPE_EQUIP       = 20 // 装备
	ITEM_TYPE_EQUIP_CHIP  = 21 // 装备碎片
	ITEM_TYPE_SOLDIER     = 22 // 佣兵（废弃）
	ITEM_TYPE_HEAD        = 23 // 头像（废弃）
	ITEM_TYPE_ARMY        = 24 // 佣兵
	ITEM_TYPE_ARMY_FLAG   = 25 // 佣兵军旗
	ITEM_TYPE_GEM         = 30 // 宝石
	ITEM_TYPE_SHOW        = 31 // 显示材料用
	ITEM_TYPE_ICON        = 32 // 头像物品
	ITEM_TYPE_PORTRAIT    = 33 // 头像框物品
	ITEM_TYPE_RANDOM_HERO = 34 // 随机英雄
	ITEM_TYPE_RUNE        = 40 // 符文
	ITEM_TYPE_ONHOOK_ITEM = 50 // 挂机相关道具
	//新魔龙
	ITEM_TYPE_EQUIP      = 66  // 装备
	ITEM_TYPE_ARTIFACT   = 68  // 神器
	ITEM_TYPE_CAN_REMOVE = 96  // 可回收道具
	ITEM_TYPE_LOTTERY    = 99  // 根据lottery去随机
	ITEM_TYPE_RAND_TIMES = 100 // 多随机组
)

const (
	ITEM_TYPE_SPECIAL_HERO_BOX = 12 //英雄许愿匣
	ITEM_TYPE_SPECIAL_ITEM_BOX = 13 //道具许愿匣
)

//! 背包数据库
type JS_Bag struct {
	Num int64 `json:"num"`
}

//! 背包数据库
type San_Bag struct {
	Uid  int64
	Info string

	info map[string]*JS_Bag

	DataUpdate
}

//! 背包
type ModBag struct {
	player  *Player
	Sql_Bag San_Bag
	Locker  *sync.RWMutex
	MaxKey  int
}

func (self *ModBag) OnGetData(player *Player) {
	self.player = player
	self.Locker = new(sync.RWMutex)

	sql := fmt.Sprintf("select * from `san_bag` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Bag, "san_bag", self.player.ID)

	if self.Sql_Bag.Uid <= 0 {
		self.Sql_Bag.Uid = self.player.ID
		self.Sql_Bag.info = make(map[string]*JS_Bag)
		self.Encode()
		InsertTable("san_bag", &self.Sql_Bag, 0, true)
		self.Sql_Bag.Init("san_bag", &self.Sql_Bag, true)
	} else {
		self.Decode()
		self.Sql_Bag.Init("san_bag", &self.Sql_Bag, true)
	}
}

func (self *ModBag) OnGetOtherData() {

}

func (self *ModBag) Decode() { //! 将数据库数据写入data
	self.Locker.Lock()
	defer self.Locker.Unlock()

	json.Unmarshal([]byte(self.Sql_Bag.Info), &self.Sql_Bag.info)
}

func (self *ModBag) Encode() { //! 将data数据写入数据库
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	self.Sql_Bag.Info = HF_JtoA(&self.Sql_Bag.info)
}

func (self *ModBag) onReg(handlers map[string]func(body []byte)) {
	handlers["useitem"] = self.UseItemHandle
	handlers["use_selectbox"] = self.SelectBoxHandle
	handlers["sellitem"] = self.SellItemHandle
	handlers["sellitems"] = self.SellItemsHandle
	handlers["createitem"] = self.CreateItemsHandle
	handlers["buyitem"] = self.BuyItemHandle
	handlers["mergeitem"] = self.MergeItemHandle
}

func (self *ModBag) UseItemHandle(body []byte) {
	var msg C2S_UseItem
	json.Unmarshal(body, &msg)

	if msg.Num < 0 {
		return
	}

	// 数量判断
	if self.GetItemNum(msg.Itemid) < msg.Num {
		return
	}

	// 配置判断
	config := GetCsvMgr().GetItemConfig(msg.Itemid)
	if config == nil {
		self.player.SendErr(fmt.Sprintf("item not found, itemId=%d", msg.Itemid))
		return
	}

	switch config.ItemType {
	case ITEM_TYPE_HERO:
		switch config.ItemSubType {
		case 1: //! 武将卡
			var outitem []PassItem
			outitem = append(outitem, PassItem{config.Spirititem, config.Spiritexp * msg.Num})
			outitem[0].ItemID, outitem[0].Num = self.player.AddObject(outitem[0].ItemID, outitem[0].Num, msg.Itemid, 0, 0, "使用道具")
			self.SendOnItem(outitem)
		}
	case ITEM_TYPE_SPECIAL:
		switch config.ItemSubType {
		case 0: //! 技能值
			self.player.AddSkillPoint(config.ExchangeExp*msg.Num, msg.Itemid, 0, "使用道具")
		case 2: //! 体力
			//if self.player.GetPower() >= POWERMAX {
			//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_POWERMAX"))
			//	return false
			//}
			self.player.AddPower(config.ExchangeExp*msg.Num, msg.Itemid, 0, "使用道具")
		case 4: //! 器魂
			self.player.AddObject(91000006, config.ExchangeExp*msg.Num, msg.Itemid, 0, 0, "使用道具")
			var outitem []PassItem
			outitem = append(outitem, PassItem{91000006, config.ExchangeExp * msg.Num})
			self.SendOnItem(outitem)
		case 12:

			bagitem := HF_DropForItemBagGroup(config.Special)

			if bagitem == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ITEM_ERR"))
				return
			}

			if msg.Destid+1 > len(bagitem) {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ITEM_ERR"))
				return
			}

			iteminfo := bagitem[msg.Destid]
			totalnum := msg.Num * iteminfo.Num
			var outitem []PassItem
			outitem = append(outitem, PassItem{iteminfo.ItemID, totalnum})
			self.player.AddObject(iteminfo.ItemID, totalnum, msg.Itemid, 0, 0, "使用道具")

			self.SendSelectbox(outitem)

		case 8, 15: //! 宝箱
			total := make(map[int]*Item)
			for index := 0; index < msg.Num; index++ {
				if config.LotteryId == 0 {
					continue
				}
				items := GetLootMgr().LootItem(config.LotteryId, self.player)
				AddItemMap(total, items)
			}
			outItems := self.player.AddObjectItemMap(total, "使用道具", msg.Itemid, 0, 0)

			self.SendOnItem(outItems)

			//!宝箱稀有提示
			if msg.Itemid == 81900071 && len(outItems) > 0 && outItems[0].ItemID == DEFAULT_GEM && outItems[0].Num > 6*msg.Num {
				//GetServer().Notice(fmt.Sprintf("%s的%s人品爆棚，开启商会钻石宝箱额外获得%d钻石",
				//	CAMP_NAME[self.player.Sql_UserBase.Camp],
				//	self.player.Sql_UserBase.UName, outitem[0].Num-6*num), 0, 0)
			}
		}
	case ITEM_TYPE_MONEY:
		self.player.AddObject(msg.Itemid, msg.Num, 5, msg.Itemid, 0, "使用道具")
	case ITEM_TYPE_MONTH_CARD:
		switch config.ItemSubType {
		case 1: //! 超值月卡
			if msg.Num > 1 {
				return
			}
			csv, ok := GetCsvMgr().GetMoney(config.Special)
			if !ok {
				return
			}
			if self.player.GetModule("activity").(*ModActivity).GetMonthCardDay(101) <= 2 {
				self.player.GetModule("recharge").(*ModRecharge).RechargeEx(config.Special, int(TimeServer().Unix()), 1, HF_Atoi(csv["rmb"])*100, 0)
			} else {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("月卡已激活，请在月卡倒计时结束后使用"))
				return
			}

		case 2: //! 至尊月卡
			if msg.Num > 1 {
				return
			}
			csv, ok := GetCsvMgr().GetMoney(config.Special)
			if !ok {
				return
			}
			if self.player.GetModule("activity").(*ModActivity).GetMonthCardDay(102) <= 2 {
				self.player.GetModule("recharge").(*ModRecharge).RechargeEx(config.Special, int(TimeServer().Unix()), 1, HF_Atoi(csv["rmb"])*100, 0)
			} else {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("月卡已激活，请在月卡倒计时结束后使用"))
				return
			}
		case 3:
			if msg.Num > 1 {
				return
			}
			csv, ok := GetCsvMgr().GetMoney(config.Special)
			if !ok {
				return
			}
			key := self.NewMaxKey()
			//150天后才有可能出现重复订单,受限INT
			orderId := int(TimeServer().Unix()%10000000)*10 + key%10
			self.player.GetModule("recharge").(*ModRecharge).RechargeEx(config.Special, orderId, 1, HF_Atoi(csv["rmb"])*100, 0)
			self.player.GetModule("activitygift").(*ModActivityGift).ActivityGiftSendInfo([]byte{})
		}
	case ITEM_TYPE_RED_POCKET: // 使用红包道具
		special := config.Special
		if special <= 0 {
			LogError("special <= 0!")
			return
		}

		// 增加红包
		self.player.GetModule("redpac").(*ModRedPac).CreateRedWait(special)
	case ITEM_TYPE_ICON:
		if !self.player.GetModule("head").(*ModHead).CheckUseItem(config.ItemId) {
			return
		}
		outItem := make([]PassItem, 0)
		outItem = append(outItem, PassItem{ItemID: msg.Itemid, Num: msg.Num})
		self.SendOnItem(outItem)
	case ITEM_TYPE_PORTRAIT:
		if !self.player.GetModule("head").(*ModHead).CheckUseItem(config.ItemId) {
			return
		}
		outItem := make([]PassItem, 0)
		outItem = append(outItem, PassItem{ItemID: msg.Itemid, Num: msg.Num})
		self.SendOnItem(outItem)
	case ITEM_TYPE_ONHOOK_ITEM:
		rewards := self.player.GetModule("onhook").(*ModOnHook).CalItem(msg.Itemid, msg.Num)
		outItems := self.player.AddObjectItemMap(rewards, "使用道具", msg.Itemid, msg.Num, 0)
		self.SendOnItem(outItems)
	case ITEM_TYPE_RAND_TIMES:
		lootGroup := strings.Split(config.LotDrop, "|")
		if msg.Destid > len(lootGroup) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ITEM_ERR"))
			return
		}
		lottery := HF_Atoi(lootGroup[msg.Destid-1])
		items := GetLootMgr().LootItem(lottery, self.player)
		if len(items) == 0 {
			return
		}
		var outitem []PassItem
		for _, v := range items {
			item := self.player.AddObjectSimple(v.ItemId, v.ItemNum, "种族随机卡", config.ItemId, 0, 0)
			outitem = append(outitem, item...)
		}
		self.SendSelectbox(outitem)

		if msg.Itemid == 91000052 {
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_FIND_CAMP_CHOOSE, 1, msg.Destid, 0, "自选阵营招募", 0, 0, self.player)
		}
	}

	self.AddItem(msg.Itemid, -msg.Num, msg.Itemid, 0, "使用道具")
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BAG_USE, msg.Itemid, -msg.Num, 0, "使用道具", 0, 0, self.player)
}

func (self *ModBag) SelectBoxHandle(body []byte) {
	var msg C2S_UseItemSelect
	json.Unmarshal(body, &msg)

	if msg.Num < 0 {
		return
	}

	// 数量判断
	if self.GetItemNum(msg.Itemid) < msg.Num {
		return
	}

	// 配置判断
	config := GetCsvMgr().GetItemConfig(msg.Itemid)
	if config == nil {
		self.player.SendErr(fmt.Sprintf("item not found, itemId=%d", msg.Itemid))
		return
	}

	switch config.ItemType {
	case ITEM_TYPE_HERO:
		switch config.ItemSubType {
		case 1: //! 武将卡
			var outitem []PassItem
			outitem = append(outitem, PassItem{config.Spirititem, config.Spiritexp * msg.Num})
			outitem[0].ItemID, outitem[0].Num = self.player.AddObject(outitem[0].ItemID, outitem[0].Num, msg.Itemid, 0, 0, "背包使用")
			self.SendOnItem(outitem)
		}
	case ITEM_TYPE_SPECIAL:
		switch config.ItemSubType {
		case 0: //! 技能值
			self.player.AddSkillPoint(config.ExchangeExp*msg.Num, msg.Itemid, 0, "背包使用")
		case 2: //! 体力
			//if self.player.GetPower() >= POWERMAX {
			//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_POWERMAX"))
			//	return false
			//}
			self.player.AddPower(config.ExchangeExp*msg.Num, msg.Itemid, 0, "背包使用")
		case 4: //! 器魂
			self.player.AddObject(91000006, config.ExchangeExp*msg.Num, msg.Itemid, 0, 0, "背包使用")
			var outitem []PassItem
			outitem = append(outitem, PassItem{91000006, config.ExchangeExp * msg.Num})
			self.SendOnItem(outitem)
		case ITEM_TYPE_SPECIAL_HERO_BOX:
			bagitem := HF_DropForItemBagGroup(config.Special)
			if bagitem == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ITEM_ERR"))
				return
			}

			//看看数量对不对
			all := 0
			for i := 0; i < len(msg.DestidNum); i++ {
				all += msg.DestidNum[i]
			}
			if all != msg.Num {
				return
			}

			outItems := make([]PassItem, 0)
			for i := 0; i < len(msg.Destid); i++ {
				if msg.Destid[i] > len(bagitem) {
					continue
				}

				self.player.AddObject(bagitem[msg.Destid[i]-1].ItemID, bagitem[msg.Destid[i]-1].Num*msg.DestidNum[i], msg.Itemid, msg.DestidNum[i], 0, "背包使用")
				outItems = append(outItems, PassItem{bagitem[msg.Destid[i]-1].ItemID, bagitem[msg.Destid[i]-1].Num * msg.DestidNum[i]})
			}
			self.SendOnItem(outItems)
		case ITEM_TYPE_SPECIAL_ITEM_BOX:
			bagitem := HF_DropForItemBagGroup(config.Special)
			if bagitem == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ITEM_ERR"))
				return
			}

			//看看数量对不对
			all := 0
			for i := 0; i < len(msg.DestidNum); i++ {
				all += msg.DestidNum[i]
			}
			if all != msg.Num {
				return
			}
			total := make([]int, 0)
			totalNum := make([]int, 0)
			for i := 0; i < len(msg.Destid); i++ {
				if msg.Destid[i] > len(bagitem) {
					continue
				}
				total = append(total, bagitem[msg.Destid[i]-1].ItemID)
				totalNum = append(totalNum, bagitem[msg.Destid[i]-1].Num*msg.DestidNum[i])
			}
			outItems := self.player.AddObjectLst(total, totalNum, "背包使用", 0, 0, 0)
			self.SendOnItem(outItems)
		case 8: //! 宝箱
			total := make(map[int]*Item)
			for index := 0; index < msg.Num; index++ {
				if config.LotteryId == 0 {
					continue
				}
				items := GetLootMgr().LootItem(config.LotteryId, self.player)
				AddItemMap(total, items)
			}
			outItems := self.player.AddObjectItemMap(total, "背包使用", msg.Itemid, 0, 0)

			self.SendOnItem(outItems)

			//!宝箱稀有提示
			if msg.Itemid == 81900071 && len(outItems) > 0 && outItems[0].ItemID == DEFAULT_GEM && outItems[0].Num > 6*msg.Num {
				//GetServer().Notice(fmt.Sprintf("%s的%s人品爆棚，开启商会钻石宝箱额外获得%d钻石",
				//	CAMP_NAME[self.player.Sql_UserBase.Camp],
				//	self.player.Sql_UserBase.UName, outitem[0].Num-6*num), 0, 0)
			}
		}
	case ITEM_TYPE_MONEY:
		self.player.AddObject(msg.Itemid, msg.Num, 5, msg.Itemid, 0, "背包使用")
	case ITEM_TYPE_MONTH_CARD:
		switch config.ItemSubType {
		case 1: //! 超值月卡
			if msg.Num > 1 {
				return
			}
			csv, ok := GetCsvMgr().GetMoney(101)
			if !ok {
				return
			}
			if self.player.GetModule("activity").(*ModActivity).GetMonthCardDay(101) <= 2 {
				self.player.GetModule("recharge").(*ModRecharge).Recharge(101, int(TimeServer().Unix()), 1, HF_Atoi(csv["rmb"])*100)
			} else {
				return
			}

		case 2: //! 至尊月卡
			if msg.Num > 1 {
				return
			}
			csv, ok := GetCsvMgr().GetMoney(102)
			if !ok {
				return
			}
			if self.player.GetModule("activity").(*ModActivity).GetMonthCardDay(102) <= 2 {
				self.player.GetModule("recharge").(*ModRecharge).Recharge(102, int(TimeServer().Unix()), 1, HF_Atoi(csv["rmb"])*100)
			} else {
				return
			}
		}
	case ITEM_TYPE_RED_POCKET: // 使用红包道具
		special := config.Special
		if special <= 0 {
			LogError("special <= 0!")
			return
		}

		// 增加红包
		self.player.GetModule("redpac").(*ModRedPac).CreateRedWait(special)
	case ITEM_TYPE_ICON:
		if !self.player.GetModule("head").(*ModHead).CheckUseItem(config.ItemId) {
			return
		}
		outItem := make([]PassItem, 0)
		outItem = append(outItem, PassItem{ItemID: msg.Itemid, Num: msg.Num})
		self.SendOnItem(outItem)
	case ITEM_TYPE_PORTRAIT:
		if !self.player.GetModule("head").(*ModHead).CheckUseItem(config.ItemId) {
			return
		}
		outItem := make([]PassItem, 0)
		outItem = append(outItem, PassItem{ItemID: msg.Itemid, Num: msg.Num})
		self.SendOnItem(outItem)
	}

	self.AddItem(msg.Itemid, -msg.Num, msg.Itemid, 0, "背包使用")

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BAG_USE, msg.Itemid, -msg.Num, 0, "背包使用", 0, 0, self.player)

}

func (self *ModBag) SellItemHandle(body []byte) {
	var c2s_msg C2S_SellItem
	json.Unmarshal(body, &c2s_msg)
	self.player.SendRet4("sellitem", self.SellItem(HF_Atoi(c2s_msg.Itemid), c2s_msg.Num))
}

func (self *ModBag) SellItemsHandle(body []byte) {
	var msg C2S_SellItems
	json.Unmarshal(body, &msg)
	self.SellItems(msg.Item)
}

func (self *ModBag) CreateItemsHandle(body []byte) {

	//self.player.GetModule("recharge").(*ModRecharge).GetFundAward()
	//return

	if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
		return
	}
	var c2s_msg C2S_CreateItem
	json.Unmarshal(body, &c2s_msg)
	if c2s_msg.Num >= math.MaxInt32 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_BAG_OVER_THE_MAXIMUM_VALUE_21"))
		return
	}
	id, num := self.player.AddObject(c2s_msg.Itemid, c2s_msg.Num, 4, c2s_msg.Itemid, 0, "gm指令")
	if num > 0 {
		self.SendOnItem([]PassItem{{id, num}})
	}
	return
}

func (self *ModBag) BuyItemHandle(body []byte) {
	var c2s_msg C2S_CreateItem
	json.Unmarshal(body, &c2s_msg)
	var msg S2C_BuyItem
	msg.Cid = "buyitem"
	msg.Ret, msg.Itemlst = self.BuyItem(c2s_msg.Itemid, c2s_msg.Num)
	self.player.SendMsg("buyitem", HF_JtoB(&msg))
}

func (self *ModBag) MergeItemHandle(body []byte) {
	var msg C2S_MergeItem
	json.Unmarshal(body, &msg)

	config := GetCsvMgr().GetItemConfig(msg.Itemid)
	if config == nil {
		return
	}

	if config.ItemType == ITEM_TYPE_RANDOM_HERO {
		self.MergeHero(msg.Itemid, msg.Num)
	} else {
		self.MergeItem(msg.Itemid, msg.Num)
	}

}

func (self *ModBag) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModBag) OnSave(sql bool) {
	self.Encode()
	self.Sql_Bag.Update(sql)
}

func (self *ModBag) AddItem(itemid int, num int, param1, param2 int, dec string) {
	curnum := int64(0)
	self.Locker.RLock()
	value, ok := self.Sql_Bag.info[fmt.Sprintf("%d", itemid)]
	self.Locker.RUnlock()
	if num > 0 { //! 加道具
		if !ok { //! 该道具不存在
			value = new(JS_Bag)
			value.Num = int64(num)
			self.Locker.Lock()
			self.Sql_Bag.info[fmt.Sprintf("%d", itemid)] = value
			self.Locker.Unlock()
			curnum = int64(num)
		} else {
			value.Num += int64(num)
			limit := GetCsvMgr().GetItemLimit(itemid)
			if limit > 0 {
				value.Num = HF_MinInt64(value.Num, limit)
			}
			curnum = value.Num
		}

		// 检查道具类型
		itemConfig := GetCsvMgr().GetItemConfig(itemid)
		if itemConfig != nil && itemConfig.ItemType == ITEM_TYPE_GEM {
			self.player.HandleTask(OwnGemNumTask, 0, 0, 0)
		}

		if itemConfig != nil && itemConfig.ItemType == ITEM_TYPE_EQUIP {
			self.player.HandleTask(EquipUpNumTask, 0, 0, 0)
		}

		if itemid == ITEM_ACCESS_ITEM {
			GetAccessCardRecordMgr().UpdatePoint(self.player, int(curnum))
		}
		//GetServer().SendLog_SDKUP_ItemChange(self.player, 1, dec, itemid, 0, num, int(curnum))
	} else if num < 0 { //! 删道具
		if !ok {
			return
		} else {
			value.Num += int64(num)
			curnum = value.Num
			if value.Num <= 0 {
				self.Locker.Lock()
				delete(self.Sql_Bag.info, fmt.Sprintf("%d", itemid))
				self.Locker.Unlock()
				curnum = 0
			}
		}
		//GetServer().SendLog_SDKUP_ItemChange(self.player, -1, dec, itemid, 0, num, int(curnum))
	}
	// 获得/使用物品
	itemname := GetCsvMgr().GetItemName(itemid)
	if itemid < 90000000 {
		if num > 0 {
			GetServer().sendLog_GetItem(self.player, num, dec, itemid, itemname, int(curnum)) // _,_,来源，类型，名称
		} else if num < 0 {
			GetServer().sendLog_UseItem(self.player, num, dec, itemid, itemname, int(curnum)) // _,_,来源，类型，名称
		}
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, itemid, num, param1, param2, dec, int(curnum), 0, self.player)
}

//! 得到道具数量
func (self *ModBag) GetItemNum(itemid int) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	value, ok := self.Sql_Bag.info[fmt.Sprintf("%d", itemid)]
	if ok {
		return self.GetShowNum(value.Num)
	}

	return 0
}

//! 使用道具

//! 卖道具
func (self *ModBag) SellItem(itemid, num int) bool {
	if num < 0 {
		return false
	}

	config := GetCsvMgr().GetItemConfig(itemid)
	if config == nil {
		return false
	}

	num = HF_MinInt(num, self.GetItemNum(itemid))
	gold := num * config.ExchangeGold

	self.AddItem(itemid, -num, 5, 0, "背包出售")
	self.player.AddGold(gold, 5, 0, "背包出售")

	self.player.SendInfo("updateuserinfo")

	return true
}

//! 卖道具
func (self *ModBag) SellItems(item []PassItem) {
	gold := 0
	for i := 0; i < len(item); i++ {
		config := GetCsvMgr().GetItemConfig(item[i].ItemID)
		if config == nil {
			continue
		}

		if item[i].Num < 0 {
			continue
		}

		item[i].Num = HF_MinInt(item[i].Num, self.GetItemNum(item[i].ItemID))
		gold += item[i].Num * config.ExchangeGold
		self.AddItem(item[i].ItemID, -item[i].Num, 5, 0, "背包一键出售")
	}
	self.player.AddGold(gold, 5, 0, "背包一键出售")

	var msg S2C_SellItems
	msg.Cid = "sellitems"
	msg.Gold = gold
	msg.Item = item
	self.player.SendMsg("sellitems", HF_JtoB(&msg))
}

//! 买道具
func (self *ModBag) BuyItem(itemid, num int) (int, []PassItem) {
	var outitem []PassItem

	if num < 0 {
		return -1, outitem
	}

	config := GetCsvMgr().GetItemConfig(itemid)
	if config == nil {
		return -1, outitem
	}

	need := config.GemPrice * num
	if self.player.Sql_UserBase.Gem < need {
		return 1, outitem
	}

	if self.player.Sql_UserBase.Level < config.BuyLv {
		return 2, outitem
	}

	self.player.AddGem(-need, 5, 0, 0, "背包购买")
	self.AddItem(itemid, num, 5, 0, "背包购买")

	outitem = append(outitem, PassItem{itemid, num})
	outitem = append(outitem, PassItem{DEFAULT_GEM, -need})

	return 0, outitem
}

//! 合成道具
func (self *ModBag) MergeItem(itemid int, num int) {
	if num < 0 {
		return
	}

	config := GetCsvMgr().GetItemConfig(itemid)
	if config == nil {
		return
	}

	cardcomposite := config.CompoundId
	if cardcomposite == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BAG_THE_PROP_CANNOT_BE_SYNTHESIZED"))
		return
	}

	cardnum := config.CompoundNum
	if self.player.GetObjectNum(itemid) < cardnum*num {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TIGER_LACK_OF_PROPS"))
		return
	}

	self.player.AddObject(itemid, -cardnum*num, itemid, -cardnum*num, 0, "合成道具")
	self.player.AddObject(cardcomposite, num, itemid, num, 0, "合成道具")

	var msg S2C_MergeItem
	msg.Cid = "mergeitem"
	msg.Item = append(msg.Item, PassItem{itemid, -cardnum * num})
	msg.Item = append(msg.Item, PassItem{cardcomposite, num})
	self.player.SendMsg("mergeitem", HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BAG_COMPOUND, itemid, num, 0, "合成道具", 0, 0, self.player)
}

func (self *ModBag) MergeHero(itemid int, num int) {
	if num < 0 {
		return
	}

	config := GetCsvMgr().GetItemConfig(itemid)
	if config == nil {
		return
	}

	if config.ItemType != ITEM_TYPE_RANDOM_HERO {
		return
	}

	cardnum := config.CompoundNum
	if self.player.GetObjectNum(itemid) < cardnum*num {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BAG_LACK_OF_PROPS"))
		return
	}

	var outItems []PassItem

	lootID := config.LotteryId

	if lootID <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BAG_LOTTERY_ID_NULL"))
		return
	}

	var ids, nums []int

	self.player.AddObject(itemid, -cardnum*num, itemid, -cardnum*num, 0, "合成道具")

	for i := 0; i < num; i++ {
		items := GetLootMgr().LootItem(lootID, self.player)

		for _, v := range items {
			ids = append(ids, v.ItemId)
			nums = append(nums, v.ItemNum)
		}
	}

	tempItems := self.player.AddObjectLst(ids, nums, "合成道具", itemid, -cardnum*num, 0)

	for _, v := range tempItems {
		index := -1

		for i, _ := range outItems {
			if v.ItemID == outItems[i].ItemID {
				index = i
				break
			}
		}

		if index < 0 {
			outItems = append(outItems, v)
		} else {
			outItems[index].Num += v.Num
		}
	}

	var msg S2C_MergeItem
	msg.Cid = "mergeitem"
	msg.Item = outItems
	msg.Item = append(msg.Item, PassItem{itemid, -cardnum * num})

	self.player.SendMsg("mergeitem", HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BAG_COMPOUND, itemid, num, 0, "合成道具", 0, 0, self.player)
}

//////////////////////
func (self *ModBag) SendInfo() {
	var msg S2C_BagInfo
	msg.Cid = "baglst"
	self.Locker.RLock()
	for itemid, value := range self.Sql_Bag.info {
		showNum := self.GetShowNum(value.Num)
		msg.Baglst = append(msg.Baglst, PassItem{HF_Atoi(itemid), showNum})
	}
	self.Locker.RUnlock()
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("baglst", smsg)
}

func (self *ModBag) SendOnItem(itemlst []PassItem) {
	var msg S2C_OnItem
	msg.Cid = "onitem"
	msg.Itemlst = itemlst
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("onitem", smsg)
}

func (self *ModBag) SendSelectbox(itemlst []PassItem) {
	var msg S2C_OnItem
	msg.Cid = "selectbox"
	msg.Itemlst = itemlst
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("selectbox", smsg)
}

func (self *ModBag) GetGemNum() int {
	num := 0
	for itemId, itemNum := range self.Sql_Bag.info {
		id := HF_Atoi(itemId)
		itemConfig := GetCsvMgr().GetGemConfig(id)
		if itemConfig == nil {
			continue
		}

		num += int(itemNum.Num)
	}
	return num
}

func (self *ModBag) GetGemLvNum(lv int) int {
	num := 0
	for itemId, itemNum := range self.Sql_Bag.info {
		id := HF_Atoi(itemId)
		itemConfig := GetCsvMgr().GetGemConfig(id)
		if itemConfig == nil {
			continue
		}

		if itemConfig.Level < lv {
			continue
		}

		num += int(itemNum.Num)
		//LogDebug("itemNum num:", itemNum.Num)
	}
	return num
}

func (self *ModBag) GetGemNumByType(gemType int) int {
	num := 0
	for itemId, itemNum := range self.Sql_Bag.info {
		id := HF_Atoi(itemId)
		itemConfig := GetCsvMgr().GetGemConfig(id)
		if itemConfig == nil {
			continue
		}

		if itemConfig.GemType != gemType {
			continue
		}

		num += int(itemNum.Num)
	}
	return num
}

func (self *ModBag) GetGemLvNumByType(lv int, gemType int) int {
	num := 0
	for itemId, itemNum := range self.Sql_Bag.info {
		id := HF_Atoi(itemId)
		itemConfig := GetCsvMgr().GetGemConfig(id)
		if itemConfig == nil {
			continue
		}

		if itemConfig.GemType != gemType {
			continue
		}

		if itemConfig.Level < lv {
			continue
		}

		num += int(itemNum.Num)
	}
	return num
}

func (self *ModBag) GmClearItem() {
	self.Locker.RLock()
	self.Sql_Bag.info = make(map[string]*JS_Bag)
	self.Locker.RUnlock()
	self.SendInfo()
}

func (self *ModBag) GetShowNum(num int64) int {
	rel := 2100000000
	if num > int64(rel) {
		return rel
	}
	return int(num)
}

func (self *ModBag) NewMaxKey() int {
	self.MaxKey++
	return self.MaxKey
}
