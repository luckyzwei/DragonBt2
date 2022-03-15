package game

import (
	"encoding/json"
	"fmt"
	"time"
)

//! 限时礼包 数据库
type San_TimeGift struct {
	Uid        int64  //! UID
	Step       int    //! 活动期数
	Info       string //! 保存数据
	UpdateTime int64  //! 刷新时间

	info []JS_TimeGiftItem
	DataUpdate
}

//! 活动结构
type JS_TimeGiftItem struct {
	BoxId  int `json:"boxid"`  //! 礼包Id
	Pickup int `json:"pickup"` //! 暂时不用
	Plan   int `json:"plan"`   //! 充值金额
	Done   int `json:"done"`   //! 完成状态, 当所有的完成了,才算完成
	Times  int `json:"times"`  //! 已完成次数
}

//! 限时礼包, Activity_LimitGift.csv
//type JS_ActivityTimeGift struct {
//	ID        int      `json:"id"`       //! 礼包唯一Id
//	Group     int      `json:"group"`    //! 分组
//	RMB       int      `json:"rmb"`      //! 充值RMB
//	Name      string   `json:"tab_name"` //! 宝箱名字
//	Descs     []string `json:"describe"` //! 子礼包类型
//	Pic       string   `json:"pic"`      //! 原价
//	TabPic    string   `json:"tab_pic"`  //! 实际价格
//	Efficacys []int    `json:"efficacy"` //! 高亮显示
//	Recharge  int      `json:"recharge"` //! 任务类型
//}

//! 限时礼包
type ModTimeGift struct {
	player   *Player
	Sql_Shop *San_TimeGift

	chg []JS_TimeGiftItem
}

func (self *ModTimeGift) OnGetData(player *Player) {
	self.player = player
}

func (self *ModTimeGift) OnGetOtherData() {
	self.Sql_Shop = new(San_TimeGift)
	sql := fmt.Sprintf("select * from `san_timegift` where `uid` = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, self.Sql_Shop, "san_timegift", self.player.ID)
	if self.Sql_Shop.Uid <= 0 {
		self.Sql_Shop.Uid = self.player.ID
		self.Sql_Shop.info = make([]JS_TimeGiftItem, 0)
		self.Sql_Shop.UpdateTime = self.GetNextTime()
		self.RefreshNow()
		self.Encode()
		InsertTable("san_timegift", self.Sql_Shop, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_Shop.Init("san_timegift", self.Sql_Shop, true)

	step := GetActivityMgr().getActN3(ACT_TIMEGIFT)
	if step > 0 && self.Sql_Shop.Step != step {
		//! 刷新
		self.RefreshNow()
	}

	actType := GetActivityMgr().GetActivity(ACT_TIMEGIFT)
	if actType != nil && actType.status.Status != ACTIVITY_STATUS_OPEN {
		self.checkAward()
	}

	//! 重新开服后，礼包自动刷新
	//if self.Sql_Shop.UpdateTime < GetServer().OpenServerTime {
	//	self.RefreshNow()
	//}
}

func (self *ModTimeGift) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	// 领取奖励
	case "drawtimegift":
		var c2s_msg C2S_LuckShopBuy
		json.Unmarshal(body, &c2s_msg)
		self.DrawItem(c2s_msg.BoxId)
		return true
		// 获取礼包信息
	case "gettimegiftlst":
		//self.RefreshShop()
		self.SendInfo()
		return true
	}

	return false
}

func (self *ModTimeGift) OnSave(sql bool) {
	if self.Sql_Shop != nil {
		self.Encode()
		self.Sql_Shop.Update(sql)
	}
}

func (self *ModTimeGift) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_Shop.Info), &self.Sql_Shop.info)
}

func (self *ModTimeGift) Encode() { //! 将data数据写入数据库
	self.Sql_Shop.Info = HF_JtoA(&self.Sql_Shop.info)
}

//! 完成但是没有领取奖励 Done = 1, 发送奖励
func (self *ModTimeGift) checkAward() {
	chg := false
	for i := 0; i < len(self.Sql_Shop.info); i++ {
		if self.Sql_Shop.info[i].Done == 1 {
			//shopitem, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
			shopitem, ok := GetCsvMgr().ActivityTimeGiftMap[self.Sql_Shop.info[i].BoxId]
			//[self.Sql_Shop.]
			if !ok {
				continue
			}

			mailId := 3002
			mailConfig, ok := GetCsvMgr().MailConfig[mailId]
			if !ok {
				LogError("邮件不存在, mailId=", mailId)
				return
			}

			out := make([]PassItem, 0)
			for j := 0; j < len(shopitem.Items); j++ {
				itemid := shopitem.Items[j]
				if itemid == 0 {
					continue
				}
				out = append(out, PassItem{itemid, shopitem.Nums[j]})
			}

			(self.player.GetModule("mail").(*ModMail)).AddMail(1, 1, 0,
				mailConfig.Mailtitle, mailConfig.Mailtxt, GetCsvMgr().GetText("STR_SYS"),
				out, true, TimeServer().Unix())
			self.Sql_Shop.info[i].Done = 2
			chg = true
		}
	}

	if chg == true {
		self.OnSave(true)
	}
}

func NewTimeGiftItem(boxId int) *JS_TimeGiftItem {
	return &JS_TimeGiftItem{
		BoxId:  boxId,
		Pickup: 0,
		Plan:   0,
		Done:   0,
		Times:  0,
	}
}

//! 充值回调
func (self *ModTimeGift) HandleTask(tasktype, n2, n3, n4 int) {
	if self.Sql_Shop == nil {
		return
	}
	for i := 0; i < len(self.Sql_Shop.info); i++ {
		item, ok := GetCsvMgr().ActivityTimeGiftMap[self.Sql_Shop.info[i].BoxId]
		if !ok {
			continue
		}
		if item.TaskTypes != tasktype {
			continue
		}

		var tasknode TaskNode
		tasknode.Tasktypes = item.TaskTypes
		tasknode.N1 = item.N[0]
		tasknode.N2 = item.N[1]
		tasknode.N3 = item.N[2]
		tasknode.N4 = item.N[3]
		plan, add := DoTask(&tasknode, self.player, n2, n3, n4)
		if plan == 0 {
			continue
		}

		chg := false
		if add {
			self.Sql_Shop.info[i].Plan += plan
			chg = true
		} else {
			if plan > self.Sql_Shop.info[i].Plan {
				self.Sql_Shop.info[i].Plan = plan
				chg = true
			}
		}

		if self.Sql_Shop.info[i].Plan >= item.N[0] {
			if self.Sql_Shop.info[i].Done == 1 { // 完成了没有进行领取
				self.checkMail(i, item)
				continue
			}
			self.Sql_Shop.info[i].Done = 1
		}

		if chg {
			self.chg = append(self.chg, self.Sql_Shop.info[i])
		}
	}
}

// 发送礼包信息
func (self *ModTimeGift) SendInfo() {
	var msg S2C_TimeGiftInfo
	msg.Cid = "gettimegift2lst"
	msg.NextUpdatetime = self.GetNextTime()

	if len(self.Sql_Shop.info) <= 0 {
		self.RefreshNow()
	}

	for i := 0; i < len(self.Sql_Shop.info); i++ {
		//shopitem, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
		shopitem, ok := GetCsvMgr().ActivityTimeGiftMap[self.Sql_Shop.info[i].BoxId]
		if !ok {
			continue
		}

		msg.Info = append(msg.Info, self.Sql_Shop.info[i])
		msg.Item = append(msg.Item, *shopitem)
	}

	// 防止客户端异常
	if len(msg.Info) == 0 {
		msg.Info = []JS_TimeGiftItem{}
	}

	if len(msg.Item) == 0 {
		msg.Item = []TimeGiftConfig{}
	}

	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("gettimegift2lst", smsg)
}

//! 每天晚上24点刷新
func (self *ModTimeGift) GetNextTime() int64 {
	now := TimeServer()
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), 24, 0, 0, 0, now.Location())
	return nextTime.Unix()
}

//! 领取礼包奖励
func (self *ModTimeGift) DrawItem(id int) {
	// 先领取奖励
	index := -1
	for i := 0; i < len(self.Sql_Shop.info); i++ {
		if self.Sql_Shop.info[i].BoxId != id {
			continue
		}

		if self.Sql_Shop.info[i].Done == 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THE_GIFT_PACKAGE_IS_NOT"))
			return
		}

		if self.Sql_Shop.info[i].Done == 2 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THE_GIFT_PACKAGE_HAS_BEEN"))
			return
		}
		index = i
		break
	}

	if index == -1 {
		self.sendNoItem()
		return
	}

	boxId := self.Sql_Shop.info[index].BoxId
	config, ok := GetCsvMgr().ActivityTimeGiftMap[boxId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THE_GIFT_PACKAGE_WAS_MISCONFIGURED"))
		return
	}

	out := make([]PassItem, 0)
	for j := 0; j < len(config.Items); j++ {
		itemid := config.Items[j]
		if itemid == 0 {
			break
		}
		out = append(out, PassItem{itemid, config.Nums[j]})
	}

	for j := 0; j < len(out); j++ {
		out[j].ItemID, out[j].Num = self.player.AddObject(out[j].ItemID, out[j].Num, id, config.Sale, 0,"幸运礼包")
	}

	pInfo := &self.Sql_Shop.info[index]
	//if config.Sort == 6 && config.Type == 1 {
	//	pInfo.Done = 2
	//} else if config.Type == 3 {
	//	pInfo.Times += 1
	//	if pInfo.Times >= config.Times {
	//		pInfo.Done = 2
	//	}
	//} else if config.Type == 2 {
	//	pInfo.Done = 2
	//}
	pInfo.Done = 2
	pInfo.Times += 1

	if self.Sql_Shop.info[index].Done != 2 {
		//if config.Type == 1 { // 刷新下一个
		//	nextInfo := self.refreshNextSort(config.Type, config.GearId, config.Sort)
		//	if nextInfo != nil {
		//		pInfo.BoxId = nextInfo.BoxId
		//		pInfo.Times = 0
		//		pInfo.Done = 0
		//		pInfo.Plan = 0
		//	} else {
		//		LogError("刷新失败, config type:", config.Type, ", boxId:", pInfo.BoxId)
		//		pInfo = nil
		//	}
		//
		//} else if config.Type == 3 {
		//	nextInfo := self.refreshNextRandom(config.Type, config.GearId)
		//	if nextInfo != nil {
		//		pInfo.BoxId = nextInfo.BoxId
		//		pInfo.Done = 0
		//		pInfo.Plan = 0
		//	} else {
		//		LogError("刷新失败, config type:", config.Type, ", boxId:", pInfo.BoxId)
		//		pInfo = nil
		//	}
		//}
	}

	if pInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_LUCK_FLUSH_ERROR")+
			GetCsvMgr().GetText("STR_LUCK_FLUSH_BOX")+fmt.Sprintf("%d", config.Id))
		return
	}

	var msg S2C_TimeGiftGet
	msg.Cid = "timegiftget"
	msg.Id = id
	msg.Item = out
	msg.NextInfo = *pInfo
	//shopitem, ok := GetActivityMgr().Sql_ActivityBox[pInfo.BoxId]
	shopitem, ok := GetCsvMgr().ActivityTimeGiftMap[pInfo.BoxId]
	if ok {
		msg.NextItem = *shopitem
	}

	self.player.SendMsg("1", HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_LUCK_SHOP, id, config.Sale, 0, "幸运礼包", 0, 0, self.player)
}

func (self *ModTimeGift) sendNoItem() {
	var msg S2C_TimeGiftGet
	msg.Cid = "timegiftget"
	msg.Id = 0
	msg.Item = make([]PassItem, 0)
	self.player.SendMsg("1", HF_JtoB(&msg))

	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THIS_COMMODITY_DOES_NOT_EXIST"))
}

// 发送礼包更新信息
func (self *ModTimeGift) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_TimeGiftUpdate
	msg.Cid = "timegift_update"
	msg.Info = self.chg
	self.chg = make([]JS_TimeGiftItem, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("timegift_update", smsg)
}

// 刷新配置, 对于gearid = 1,2 按照梯度来计算
// type = 1 按照sort来
// type = 3 按照次数来
func (self *ModTimeGift) RefreshShop() {
	if TimeServer().Unix() < self.Sql_Shop.UpdateTime {
		return
	}
	self.RefreshNow()
}

func (self *ModTimeGift) RefreshNow() {
	self.checkAward()
	// 重置充值以及礼包信息
	self.Sql_Shop.info = make([]JS_TimeGiftItem, 0)
	// 下一次刷新时间
	self.Sql_Shop.UpdateTime = self.GetNextTime()

	drop := GetActivityMgr().getActN4(ACT_TIMEGIFT)
	if drop == 0 {
		//! 没有数据，则默认采用初始化1
		drop = 1
	}

	step := GetActivityMgr().getActN3(ACT_TIMEGIFT)
	if self.Sql_Shop.Step == step {
		return
	}

	//! 是否存在配置
	_, ok := GetCsvMgr().ActivityTimeGiftGroup[drop]
	if !ok {
		return
	}

	self.Sql_Shop.Step = drop

	for _, giftItem := range GetCsvMgr().ActivityTimeGiftGroup[drop] {
		// type = 1时取第一个
		giftId := giftItem.Id
		_, ok := GetCsvMgr().ActivityTimeGiftMap[giftId]
		if !ok {
			continue
		}

		shopinfo := NewTimeGiftItem(giftId)
		self.Sql_Shop.info = append(self.Sql_Shop.info, *shopinfo)

	}
}

/*
func (self *ModTimeGift) HasMonthCard() bool {
	return self.player.GetModule("activity").(*ModActivity).HasMonthCard()
}

*/

// 刷新下一个type = 1的BoxId
func (self *ModTimeGift) refreshNextSort(actionType, gearId, sortId int) *JS_LuckShopItem {
	// 根据当前类型,刷新下一个类型
	//vip := self.player.Sql_UserBase.Vip
	//level := self.player.Sql_UserBase.Level

	boxlist, ok := GetActivityMgr().BoxGroup[gearId]
	if !ok {
		LogError("gearId error:", gearId)
		return nil
	}

	for i := 0; i < len(boxlist); i++ {
		boxId := boxlist[i]
		_, ok := GetCsvMgr().ActivityTimeGiftMap[boxId]
		if !ok {
			continue
		}

		//// 切页类型
		//if config.GearId != gearId {
		//	continue
		//}
		//
		//// 玩法类型
		//if config.Type != actionType {
		//	continue
		//}
		//
		//// 等级
		//if config.NeedLv > level {
		//	continue
		//}
		//
		//// vip等级
		//if vip < config.NeedVip1 || vip > config.NeedVip2 {
		//	continue
		//}
		//
		//// 过滤月卡
		//if config.MonthCard == 1 && !self.HasMonthCard() {
		//	continue
		//}
		//
		//if config.Sort == sortId+1 {
		//	shopinfo := NewLuckShopItem(config.BoxId)
		//	return shopinfo
		//}
	}

	return nil
}

// 刷新下一个type = 3的BoxId
/*
func (self *ModTimeGift) refreshNextRandom(actionType, gearId int) *JS_LuckShopItem {
	// 根据当前类型,刷新下一个类型
	//vip := self.player.Sql_UserBase.Vip
	//level := self.player.Sql_UserBase.Level
	_, ok := GetActivityMgr().BoxGroup[gearId]
	if !ok {
		LogError("gearId error:", gearId)
		return nil
	}

	buylist := make([]*JS_ActivityBox, 0)
	//for i := 0; i < len(boxlist); i++ {
	//	boxId := boxlist[i]
	//	config, ok := GetCsvMgr().ActivityTimeGiftMap[boxId]
	//	if !ok {
	//		continue
	//	}
	//
	//	//// 切页类型
	//	//if config.GearId != gearId {
	//	//	continue
	//	//}
	//	//
	//	//// 玩法类型
	//	//if config.Type != actionType {
	//	//	continue
	//	//}
	//	//
	//	//// 等级
	//	//if config.NeedLv > level {
	//	//	continue
	//	//}
	//	//
	//	//// vip等级
	//	//if vip < config.NeedVip1 || vip > config.NeedVip2 {
	//	//	continue
	//	//}
	//	//
	//	//// 过滤月卡
	//	//if config.MonthCard == 1 && !self.HasMonthCard() {
	//	//	continue
	//	//}
	//
	//	//buylist = append(buylist, config)
	//}
	//
	//if len(buylist) <= 0 {
	//	LogError("len(buylist <=0, config must error!")
	//	return nil
	//}

	index := HF_GetRandom(len(buylist))
	shopinfo := NewLuckShopItem(buylist[index].BoxId, buylist[index].Start, buylist[index].Continue, self.player.Sql_UserBase.Regtime)
	return shopinfo
}
 */

// 合服发送幸运礼包
func (self *ModTimeGift) FitServer() {
	if self.Sql_Shop == nil {
		return
	}

	for i := 0; i < len(self.Sql_Shop.info); i++ {
		if self.Sql_Shop.info[i].Done == 1 {
			shopitem, ok := GetCsvMgr().ActivityTimeGiftMap[self.Sql_Shop.info[i].BoxId]
			if !ok {
				continue
			}

			out := make([]PassItem, 0)
			for j := 0; j < 6; j++ {
				itemid := shopitem.Items[j]
				if itemid == 0 {
					continue
				}
				out = append(out, PassItem{itemid, shopitem.Nums[j]})
			}

			context := GetCsvMgr().GetText("STR_LUCK_SHOP_LEFT_AWARD")
			(self.player.GetModule("mail").(*ModMail)).AddMail(1, 1, 0,
				GetCsvMgr().GetText("STR_LUCK_SHOP_MAIL"), context, GetCsvMgr().GetText("STR_SYS"),
				out, true, TimeServer().Unix())
			self.Sql_Shop.info[i].Done = 2
			//log.Println("玩家：", self.player.Sql_UserBase.Uid, " 领取礼物!")
		}
	}
}

func (self *ModTimeGift) checkMail(index int, shopitem *TimeGiftConfig) {
	if shopitem == nil {
		return
	}

	out := make([]PassItem, 0)
	for j := 0; j < len(shopitem.Items); j++ {
		itemid := shopitem.Items[j]
		if itemid == 0 {
			continue
		}
		out = append(out, PassItem{itemid, shopitem.Nums[j]})
	}

	context := GetCsvMgr().GetText("STR_LUCK_SHOP_FRESH")
	(self.player.GetModule("mail").(*ModMail)).AddMail(1, 1, 0, GetCsvMgr().GetText("STR_LUCK_SHOP_MAIL"), context, GetCsvMgr().GetText("STR_SYS"),
		out, true, TimeServer().Unix())
}
