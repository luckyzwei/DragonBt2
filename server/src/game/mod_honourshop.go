package game

// 商店模块
import (
	"encoding/json"
	"fmt"
)

const (
	PURCHASETYPE_NO   = 0
	PURCHASETYPE_DAY  = 1
	PURCHASETYPE_WEEK = 2
	PURCHASETYPE_ALL  = 3
)

// 商店数据库
type San_HonourShop struct {
	Uid          int64
	Shopgood     string // 道具购买状态
	NextDayTime  int64  // 下次日刷新时间
	NextWeekTime int64  // 下次周刷新时间

	shopgood []*JS_HonourGoodInfo
	DataUpdate
}

type JS_HonourGoodInfo struct {
	Grid         int          `json:"grid"`          //格子
	Id           int          `json:"id"`            //id
	BuyState     int          `json:"buystate"`      //当前已购买次数
	ItemId       int          `json:"itemid"`        //商品ID
	ItemNum      int          `json:"itemnum"`       //商品数量
	CostId       []int        `json:"costid"`        //消耗ID数组
	CostNum      []int        `json:"costnum"`       //消耗数量数组 价格已计算折扣
	PurchaseType int          `json:"purchase_type"` //限购类型 0不限1每日2每周3永久
	PurchaseNum  int          `json:"purchase_num"`  //限购次数
	UnLockTask   *JS_TaskInfo `json:"unlocktask"`    //解锁条件
	UnLock       int          `json:"unlock"`        //解锁状态  0未解锁  1解锁
}

// 商店
type ModHonourShop struct {
	player         *Player
	Sql_HonourShop San_HonourShop // 所有商店的结构
	chg            []*JS_HonourGoodInfo
}

func (self *ModHonourShop) OnGetData(player *Player) {

	self.player = player
	sql := fmt.Sprintf("select * from `san_honourshop` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_HonourShop, "san_honourshop", self.player.ID)

	if self.Sql_HonourShop.Uid <= 0 {
		self.Sql_HonourShop.Uid = self.player.ID
		self.Sql_HonourShop.shopgood = make([]*JS_HonourGoodInfo, 0)
		self.Encode()
		InsertTable("san_honourshop", &self.Sql_HonourShop, 0, true)
		self.Sql_HonourShop.Init("san_honourshop", &self.Sql_HonourShop, true)
	} else {
		self.Decode()
		self.Sql_HonourShop.Init("san_honourshop", &self.Sql_HonourShop, true)
	}

	self.CheckGrid()
}

func (self *ModHonourShop) OnGetOtherData() {

}

func (self *ModHonourShop) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModHonourShop) OnSave(sql bool) {
	self.Encode()
	self.Sql_HonourShop.Update(sql)
}

func (self *ModHonourShop) Decode() { // 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_HonourShop.Shopgood), &self.Sql_HonourShop.shopgood)
}

func (self *ModHonourShop) Encode() { // 将data数据写入数据库
	self.Sql_HonourShop.Shopgood = HF_JtoA(&self.Sql_HonourShop.shopgood)
}

func (self *ModHonourShop) CheckGrid() {
	//需要比对配置
	//CHECK现有的数据,并生成用于排重的记录
	map_temp := make(map[int]map[int]map[int]int)
	size := len(self.Sql_HonourShop.shopgood)
	for i := size - 1; i >= 0; i-- {
		config, ok := GetCsvMgr().HonourShopConfigMap[self.Sql_HonourShop.shopgood[i].Id]
		if !ok {
			self.Sql_HonourShop.shopgood = append(self.Sql_HonourShop.shopgood[:i], self.Sql_HonourShop.shopgood[i+1:]...)
			continue
		}
		_, okP := map_temp[config.ParentLabel]
		if !okP {
			map_temp[config.ParentLabel] = make(map[int]map[int]int)
		}
		_, okS := map_temp[config.ParentLabel][config.SubTab]
		if !okS {
			map_temp[config.ParentLabel][config.SubTab] = make(map[int]int)
		}
		map_temp[config.ParentLabel][config.SubTab][config.Grid] = LOGIC_TRUE
		//更新信息
		if self.Sql_HonourShop.shopgood[i].UnLockTask.Tasktypes != config.TaskTypes {
			self.Sql_HonourShop.shopgood[i].UnLockTask.Tasktypes = config.TaskTypes
			self.Sql_HonourShop.shopgood[i].UnLockTask.Plan = 0
			self.Sql_HonourShop.shopgood[i].UnLockTask.Finish = 0
			self.Sql_HonourShop.shopgood[i].UnLock = LOGIC_FALSE
		}
	}
	//检查配置，生成格子
	for _, config := range GetCsvMgr().HonourShopConfigMap {
		_, okP := map_temp[config.ParentLabel]
		if !okP {
			map_temp[config.ParentLabel] = make(map[int]map[int]int)
		}
		_, okS := map_temp[config.ParentLabel][config.SubTab]
		if !okS {
			map_temp[config.ParentLabel][config.SubTab] = make(map[int]int)
		}
		_, okG := map_temp[config.ParentLabel][config.SubTab][config.Grid]
		if !okG {
			map_temp[config.ParentLabel][config.SubTab][config.Grid] = LOGIC_TRUE
			data := new(JS_HonourGoodInfo)
			data.Id = config.Id
			data.Grid = config.Grid
			data.PurchaseType = config.PurchaseType
			data.PurchaseNum = config.PurchaseNum
			data.UnLockTask = new(JS_TaskInfo)
			data.UnLockTask.Tasktypes = config.TaskTypes
			data.UnLockTask.Taskid = config.Id
			self.Sql_HonourShop.shopgood = append(self.Sql_HonourShop.shopgood, data)
		}
	}
}

func (self *ModHonourShop) OnRefresh() {
	self.Sql_HonourShop.NextDayTime = HF_GetNextDayStart()
	//处理周刷新
	isNeedWeekRefresh := false
	if self.Sql_HonourShop.NextWeekTime <= TimeServer().Unix() {
		isNeedWeekRefresh = true
		self.Sql_HonourShop.NextWeekTime = HF_GetNextWeekStart()
	}
	//刷新物品
	preLevelIndex := self.player.GetModule("onhook").(*ModOnHook).GetStage()
	for _, v := range self.Sql_HonourShop.shopgood {
		config, ok := GetCsvMgr().HonourShopConfigMap[v.Id]
		if !ok {
			continue
		}
		goodConfig := GetCsvMgr().GetHonourShopConfig(config.ParentLabel, config.SubTab, config.Grid, preLevelIndex)
		if goodConfig == nil {
			continue
		}
		if goodConfig.Id == 110101 || goodConfig.Id == 110102 {
			print(111)
		}
		v.Id = goodConfig.Id
		v.ItemId = goodConfig.ItemId
		v.ItemNum = goodConfig.ItemNumber
		v.CostId = goodConfig.CostItem
		v.CostNum = goodConfig.CostNum
		if v.PurchaseType != config.PurchaseType {
			v.PurchaseType = config.PurchaseType
			v.BuyState = 0
		}
		v.PurchaseNum = config.PurchaseNum
		switch v.PurchaseType {
		case PURCHASETYPE_DAY:
			v.BuyState = 0
		case PURCHASETYPE_WEEK:
			if isNeedWeekRefresh {
				v.BuyState = 0
			}
		}
	}
}

func (self *ModHonourShop) SendInfo() {
	now := TimeServer().Unix()
	if now >= self.Sql_HonourShop.NextDayTime {
		self.OnRefresh()
	}

	var msg S2C_HonourShopInfo
	msg.Cid = "honourshopinfo"
	msg.HonourGoodInfo = self.Sql_HonourShop.shopgood
	msg.NextDayTime = self.Sql_HonourShop.NextDayTime
	msg.NextWeekTime = self.Sql_HonourShop.NextWeekTime
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}

func (self *ModHonourShop) onReg(handlers map[string]func(body []byte)) {
	handlers["honourshopbuy"] = self.HonourShopBuy
}

func (self *ModHonourShop) HonourShopBuy(body []byte) {

	var msg C2S_HonourShopBuy
	json.Unmarshal(body, &msg)

	for i := 0; i < len(self.Sql_HonourShop.shopgood); i++ {
		if self.Sql_HonourShop.shopgood[i].Id == msg.Id {
			//判断是否解锁
			if self.Sql_HonourShop.shopgood[i].UnLock == LOGIC_FALSE {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_INADEQUATE_EXCHANGE_CONDITIONS"))
				return
			}
			//判断限购
			if self.Sql_HonourShop.shopgood[i].PurchaseType != PURCHASETYPE_NO && self.Sql_HonourShop.shopgood[i].BuyState >= self.Sql_HonourShop.shopgood[i].PurchaseNum {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_ALREADY_BUY"))
				return
			}

			if err := self.player.HasObjectOk(self.Sql_HonourShop.shopgood[i].CostId, self.Sql_HonourShop.shopgood[i].CostNum); err != nil {
				self.player.SendErrInfo("err", err.Error())
				return
			}

			costItem := self.player.RemoveObjectLst(self.Sql_HonourShop.shopgood[i].CostId, self.Sql_HonourShop.shopgood[i].CostNum, "荣誉商店购买", msg.Id, 0, 1)
			num := GetGemNum(costItem)
			param3 := 0
			if num > 0 {
				param3 = -1
			}
			getItems := self.player.AddObjectSimple(self.Sql_HonourShop.shopgood[i].ItemId, self.Sql_HonourShop.shopgood[i].ItemNum, "荣誉商店购买", msg.Id, 0, param3)
			//CheckAddItemLog(self.player, "商店购买", costItem, getItems)

			if num > 0 {
				AddSpecialSdkItemListLog(self.player, num, getItems, "商店购买")
			}

			self.Sql_HonourShop.shopgood[i].BuyState += 1

			//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SHOP_BUY, shop.shopgood[i].ItemId, shoptype, shop.shopgood[i].ItemNum, "商店购买", 0, 0, self.player)

			var msgRel S2C_HonourShopBuy
			msgRel.Cid = "honourshopbuy"
			msgRel.GetItems = getItems
			msgRel.CostItems = costItem
			msgRel.HonourShopInfo = self.Sql_HonourShop.shopgood[i]
			self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
			self.player.HandleTask(TASK_TYPE_SHOP_BUY_COUNT, 1, msg.Id, 0)
			self.player.GetModule("task").(*ModTask).SendUpdate()

			//for i := 0; i < len(msgRel.GetItems); i++ {
			//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SHOP_BUY_GOODS, msgRel.GetItems[i].ItemID, shoptype, 0, "商店购买道具", 0, 0, self.player)
			//}
			return
		}
	}

	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_GOODS_NOT_FIND"))
	return
}

func (self *ModHonourShop) HandleTask(tasktype, n2, n3, n4 int) {
	for _, v := range self.Sql_HonourShop.shopgood {
		if v.UnLock == LOGIC_TRUE {
			continue
		}
		if v.UnLockTask == nil {
			continue
		}
		if v.UnLockTask.Tasktypes != tasktype {
			continue
		}
		if v.UnLockTask.Finish > CANTFINISH {
			continue
		}
		config := GetCsvMgr().HonourShopConfigMap[v.Id]
		if config == nil {
			continue
		}

		var tasknode TaskNode
		tasknode.Tasktypes = config.TaskTypes
		tasknode.N1 = config.Ns[0]
		tasknode.N2 = config.Ns[1]
		tasknode.N3 = config.Ns[2]
		tasknode.N4 = config.Ns[3]
		plan, add := DoTask(&tasknode, self.player, n2, n3, n4)
		if plan == 0 {
			continue
		}

		if add {
			v.UnLockTask.Plan += plan
		} else {
			if tasktype == PvpRankNow {
				if plan != 0 { // 新排名为0则直接不处理
					if v.UnLockTask.Plan == 0 { // 进度为0则说明未初始化 直接赋值
						v.UnLockTask.Plan = plan
					} else { // 进度不为0 则需要判断 获得的新名次比之前要高 则赋值
						if plan < v.UnLockTask.Plan {
							v.UnLockTask.Plan = plan
						}
					}
				}
			} else {
				if plan > v.UnLockTask.Plan {
					v.UnLockTask.Plan = plan
				}
			}
		}

		if v.UnLockTask.Plan >= config.Ns[0] {
			v.UnLockTask.Finish = LOGIC_TRUE
			v.UnLock = LOGIC_TRUE //解锁
			//同步状态
			self.chg = append(self.chg, v)
		}
	}
}

func (self *ModHonourShop) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_UpdateHonourShop
	msg.Cid = "updatehonourshop"
	msg.HonourShopInfo = self.chg
	self.chg = make([]*JS_HonourGoodInfo, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}
