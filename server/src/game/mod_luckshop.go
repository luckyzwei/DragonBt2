package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

//! 幸运商店数据库
type San_LuckShop struct {
	Uid         int64
	Lastupdtime int64  //! 刷新时间
	Info        string //! 保存数据
	Ver         int    //! 版本
	StarVer     int    //! 星辰礼包版本
	DiscountVer int    //! 特惠礼包版本
	VerGroup    string //! 版本号组

	info     []JS_LuckShopItem
	verGroup map[int]int
	DataUpdate
}

//! 活动结构
type JS_LuckShopItem struct {
	BoxId        int   `json:"boxid"`        //! 礼包Id
	Pickup       int   `json:"pickup"`       //! 领取次数
	Plan         int   `json:"plan"`         //! 充值金额
	Done         int   `json:"done"`         //! 完成状态, 当所有的完成了,才算完成
	Times        int   `json:"times"`        //! 已完成次数
	StartTime    int64 `json:"starttime"`    //! 开始时间
	EndTime      int64 `json:"endtime"`      //! 结束时间
	Group        int   `json:"group"`        //! 对应的组
	ActivityType int   `json:"activitytype"` //! 对应的活动类型
}

//! 幸运商店
type ModLuckShop struct {
	player   *Player
	Sql_Shop *San_LuckShop

	chg []JS_LuckShopItem
}

func (self *ModLuckShop) OnGetData(player *Player) {
	self.player = player
}

func (self *ModLuckShop) OnGetOtherData() {
	self.Sql_Shop = new(San_LuckShop)
	sql := fmt.Sprintf("select * from `san_luckshop` where `uid` = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, self.Sql_Shop, "san_luckshop", self.player.ID)
	if self.Sql_Shop.Uid <= 0 {
		self.Sql_Shop.Uid = self.player.ID
		self.Sql_Shop.info = make([]JS_LuckShopItem, 0)
		self.Sql_Shop.verGroup = make(map[int]int, 0)
		self.Encode()
		InsertTable("san_luckshop", self.Sql_Shop, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_Shop.Init("san_luckshop", self.Sql_Shop, true)
}

func (self *ModLuckShop) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	// 领取奖励
	case "drawluckshop":
		var c2s_msg C2S_LuckShopBuy
		json.Unmarshal(body, &c2s_msg)
		self.DrawItem(c2s_msg.BoxId)
		return true
		// 获取礼包信息
	case "getluckshoplst":
		self.SendInfo()
		return true
	}

	return false
}

func (self *ModLuckShop) OnSave(sql bool) {
	if self.Sql_Shop != nil {
		self.Encode()
		self.Sql_Shop.Update(sql)
	}

}

func (self *ModLuckShop) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_Shop.Info), &self.Sql_Shop.info)
	json.Unmarshal([]byte(self.Sql_Shop.VerGroup), &self.Sql_Shop.verGroup)
}

func (self *ModLuckShop) Encode() { //! 将data数据写入数据库
	self.Sql_Shop.Info = HF_JtoA(&self.Sql_Shop.info)
	self.Sql_Shop.VerGroup = HF_JtoA(&self.Sql_Shop.verGroup)
}

//
////! 完成但是没有领取奖励 Done = 1, 发送奖励
//func (self *ModLuckShop) checkAward() {
//	for i := 0; i < len(self.Sql_Shop.info); i++ {
//		if self.Sql_Shop.info[i].Done == 1 {
//			shopitem, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
//			if !ok {
//				continue
//			}
//
//			out := make([]PassItem, 0)
//			for j := 0; j < 6; j++ {
//				itemid := shopitem.Item[j]
//				if itemid == 0 {
//					continue
//				}
//				out = append(out, PassItem{itemid, shopitem.Num[j]})
//			}
//
//			context := GetCsvMgr().GetText("STR_LUCK_SHOP_FRESH")
//			(self.player.GetModule("mail").(*ModMail)).AddMail(1, 1, 0,
//				GetCsvMgr().GetText("STR_LUCK_SHOP_MAIL"), context, GetCsvMgr().GetText("STR_SYS"),
//				out, true, TimeServer().Unix())
//			self.Sql_Shop.info[i].Done = 2
//		}
//	}
//}

func (self *ModLuckShop) CheckBox() {
	if self.Sql_Shop.info == nil {
		self.Sql_Shop.info = make([]JS_LuckShopItem, 0)
	}

	if self.Sql_Shop.verGroup == nil {
		self.Sql_Shop.verGroup = make(map[int]int)
	}

	//self.checkAward()
	ver := GetActivityMgr().getActN3(ActivityLuckShop)
	group := GetActivityMgr().getActN4(ActivityLuckShop)
	if ver != self.Sql_Shop.Ver {
		size := len(self.Sql_Shop.info)
		for i := size - 1; i >= 0; i-- {
			if self.Sql_Shop.info[i].ActivityType != ActivityLuckShop {
				continue
			}
			self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
		}
		self.Sql_Shop.Ver = ver
	}

	boxlist := GetActivityMgr().BoxGroup[group]
	size := len(self.Sql_Shop.info)
	for i := size - 1; i >= 0; i-- {
		value := self.Sql_Shop.info[i]
		if value.ActivityType != ActivityLuckShop {
			continue
		}
		find := false
		for _, boxId := range boxlist {
			config, ok := GetActivityMgr().Sql_ActivityBox[boxId]
			if !ok {
				continue
			}

			if config.ActivityType != ActivityLuckShop {
				continue
			}

			if config.BoxId != value.BoxId {
				continue
			}
			find = true
			break
		}

		if !find {
			self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
			continue
		}

		config, _ := GetActivityMgr().Sql_ActivityBox[value.BoxId]
		starttime := HF_CalTimeForConfig(config.Start, self.player.Sql_UserBase.Regtime)
		if value.StartTime != starttime || value.Group != config.GearId {
			self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
			continue
		}
	}

	starver := GetActivityMgr().getActN3(ActivityStarGift)
	stargroup := GetActivityMgr().getActN4(ActivityStarGift)
	if starver != self.Sql_Shop.StarVer {
		size := len(self.Sql_Shop.info)
		for i := size - 1; i >= 0; i-- {
			if self.Sql_Shop.info[i].ActivityType != ActivityStarGift {
				continue
			}
			self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
		}
		self.Sql_Shop.StarVer = starver
	}

	starboxlist := GetActivityMgr().BoxGroup[stargroup]
	size = len(self.Sql_Shop.info)
	for i := size - 1; i >= 0; i-- {
		value := self.Sql_Shop.info[i]
		if value.ActivityType != ActivityStarGift {
			continue
		}
		find := false
		for _, boxId := range starboxlist {
			config, ok := GetActivityMgr().Sql_ActivityBox[boxId]
			if !ok {
				continue
			}

			if config.ActivityType != ActivityStarGift {
				continue
			}

			if config.BoxId != value.BoxId {
				continue
			}
			find = true
			break
		}

		if !find {
			self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
			continue
		}

		config, _ := GetActivityMgr().Sql_ActivityBox[value.BoxId]
		if value.Group != config.GearId {
			self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
			continue
		}
	}

	discountver := GetActivityMgr().getActN3(ActivityDiscountGift)
	discountgroup := GetActivityMgr().getActN4(ActivityDiscountGift)
	if discountver != self.Sql_Shop.DiscountVer {
		size := len(self.Sql_Shop.info)
		for i := size - 1; i >= 0; i-- {
			if self.Sql_Shop.info[i].ActivityType != ActivityDiscountGift {
				continue
			}
			self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
		}
		self.Sql_Shop.DiscountVer = discountver
	}

	discountboxlist := GetActivityMgr().BoxGroup[discountgroup]
	size = len(self.Sql_Shop.info)
	for i := size - 1; i >= 0; i-- {
		value := self.Sql_Shop.info[i]
		if value.ActivityType != ActivityDiscountGift {
			continue
		}
		find := false
		for _, boxId := range discountboxlist {
			config, ok := GetActivityMgr().Sql_ActivityBox[boxId]
			if !ok {
				continue
			}

			if config.ActivityType != ActivityDiscountGift {
				continue
			}

			if config.BoxId != value.BoxId {
				continue
			}
			find = true
			break
		}

		if !find {
			self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
			continue
		}

		config, _ := GetActivityMgr().Sql_ActivityBox[value.BoxId]
		if value.Group != config.GearId {
			self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
			continue
		}
	}

	for id := ACT_STAR_GIFT_MIN; id < ACT_STAR_GIFT_MAX; id++ {
		ver := GetActivityMgr().getActN3(id)
		group := GetActivityMgr().getActN4(id)
		if ver != self.Sql_Shop.verGroup[id] {
			size := len(self.Sql_Shop.info)
			for i := size - 1; i >= 0; i-- {
				if self.Sql_Shop.info[i].ActivityType != id {
					continue
				}
				self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
			}
			self.Sql_Shop.verGroup[id] = ver
		}

		boxlist := GetActivityMgr().BoxGroup[group]
		size = len(self.Sql_Shop.info)
		for i := size - 1; i >= 0; i-- {
			value := self.Sql_Shop.info[i]
			if value.ActivityType != id {
				continue
			}
			find := false
			for _, boxId := range boxlist {
				config, ok := GetActivityMgr().Sql_ActivityBox[boxId]
				if !ok {
					continue
				}

				if config.ActivityType != id {
					continue
				}

				if config.BoxId != value.BoxId {
					continue
				}

				find = true

				configTemp, _ := GetActivityMgr().Sql_ActivityBox[value.BoxId]
				starttime := HF_CalTimeForConfig(configTemp.Start, self.player.Sql_UserBase.Regtime)
				self.Sql_Shop.info[i].StartTime = starttime
				self.Sql_Shop.info[i].EndTime = starttime + int64(configTemp.Continue)
				break
			}

			if !find {
				self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
				continue
			}

			config, _ := GetActivityMgr().Sql_ActivityBox[value.BoxId]
			if value.Group != config.GearId {
				self.Sql_Shop.info = append(self.Sql_Shop.info[:i], self.Sql_Shop.info[i+1:]...)
				continue
			}
		}
	}

	for _, boxlist := range GetActivityMgr().BoxGroup {
		buylist := make([]*JS_ActivityBox, 0)
		for i := 0; i < len(boxlist); i++ {
			boxId := boxlist[i]
			config, ok := GetActivityMgr().Sql_ActivityBox[boxId]
			if !ok {
				continue
			}

			buylist = append(buylist, config)
		}

		if len(buylist) <= 0 {
			LogError("len(buylist <=0, config must error!")
			continue
		}

		for index, _ := range buylist {
			find := false
			for _, v := range self.Sql_Shop.info {
				if v.BoxId == buylist[index].BoxId {
					find = true
					break
				}
			}
			if !find {
				shopinfo := NewLuckShopItem(buylist[index].GearId, buylist[index].BoxId, buylist[index].Start, buylist[index].Continue, self.player.Sql_UserBase.Regtime, buylist[index].ActivityType)
				self.Sql_Shop.info = append(self.Sql_Shop.info, *shopinfo)
			}
		}
	}
}

func NewLuckShopItem(group int, boxId int, startDay string, continueTime int, rTime string, activityType int) *JS_LuckShopItem {
	starttime := HF_CalTimeForConfig(startDay, rTime)

	return &JS_LuckShopItem{
		BoxId:        boxId,
		Pickup:       0,
		Plan:         0,
		Done:         0,
		Times:        0,
		StartTime:    starttime,
		EndTime:      starttime + int64(continueTime),
		Group:        group,
		ActivityType: activityType,
	}
}

//! 充值回调
func (self *ModLuckShop) HandleTask(tasktype, n2, n3, n4 int) {
	if self.Sql_Shop == nil {
		return
	}

	now := TimeServer().Unix()
	group := GetActivityMgr().getActN4(ActivityLuckShop)
	for i := 0; i < len(self.Sql_Shop.info); i++ {
		item, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
		if !ok {
			continue
		}
		if ActivityLuckShop != item.ActivityType {
			continue
		}
		if item.TaskTypes != tasktype {
			continue
		}

		if group != self.Sql_Shop.info[i].Group {
			continue
		}

		if now < self.Sql_Shop.info[i].StartTime || now >= self.Sql_Shop.info[i].EndTime {
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
				//self.checkMail(i, item)
				continue
			} else if self.Sql_Shop.info[i].Done == 2 {
				if self.Sql_Shop.info[i].Times < item.Times {
					self.Sql_Shop.info[i].Done = 1
					chg = true
				} else {
					continue
				}
			} else {
				self.Sql_Shop.info[i].Done = 1
			}
		}

		if chg {
			self.chg = append(self.chg, self.Sql_Shop.info[i])
		}
	}
}

//! 充值回调
func (self *ModLuckShop) StarHandleTask(tasktype, n2, n3, n4 int) {
	if self.Sql_Shop == nil {
		return
	}

	activity := GetActivityMgr().GetActivity(ActivityStarGift)
	if activity == nil || activity.status.Status != ACTIVITY_STATUS_OPEN {
		return
	}

	group := GetActivityMgr().getActN4(ActivityStarGift)
	for i := 0; i < len(self.Sql_Shop.info); i++ {
		item, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
		if !ok {
			continue
		}
		if ActivityStarGift != item.ActivityType {
			continue
		}

		if item.TaskTypes != tasktype {
			continue
		}

		if item.PicType == 1 && item.StarHero != 0 {
			have := false
			hero := self.player.GetModule("hero").(*ModHero).GetVoidHero(item.StarHero)
			if hero != nil {
				have = true
			}

			// 特殊处理 12 + 英雄id + 01 是英雄的魂石id
			cardid := item.StarHero*100 + 12000000 + 1
			carditem := GetCsvMgr().GetItemConfig(cardid)
			if carditem != nil {
				if err := self.player.HasObjectOkEasy(cardid, carditem.CompoundNum); err == nil {
					have = true
				}
			}

			if have {
				continue
			}
		}

		if item.Look > self.player.GetModule("recharge").(*ModRecharge).Sql_UserRecharge.Money {
			continue
		}

		if group != self.Sql_Shop.info[i].Group {
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
				//self.checkMail(i, item)
				continue
			} else if self.Sql_Shop.info[i].Done == 2 {
				if self.Sql_Shop.info[i].Times < item.Times {
					self.Sql_Shop.info[i].Done = 1
					chg = true
				}
			} else {
				self.Sql_Shop.info[i].Done = 1
			}
		}

		if chg {
			self.chg = append(self.chg, self.Sql_Shop.info[i])
		}
	}
}

func (self *ModLuckShop) StarLimitHandleTask(tasktype, n2, n3, n4 int) {
	if self.Sql_Shop == nil {
		return
	}

	for id := ACT_STAR_GIFT_MIN; id < ACT_STAR_GIFT_MAX; id++ {
		activity := GetActivityMgr().GetActivity(id)
		if activity == nil || activity.status.Status != ACTIVITY_STATUS_OPEN {
			continue
		}

		group := GetActivityMgr().getActN4(id)
		for i := 0; i < len(self.Sql_Shop.info); i++ {
			item, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
			if !ok {
				continue
			}
			if id != item.ActivityType {
				continue
			}

			if item.TaskTypes != tasktype {
				continue
			}

			if item.PicType == 1 && item.StarHero != 0 {
				have := false
				hero := self.player.GetModule("hero").(*ModHero).GetVoidHero(item.StarHero)
				if hero != nil {
					have = true
				}

				// 特殊处理 12 + 英雄id + 01 是英雄的魂石id
				cardid := item.StarHero*100 + 12000000 + 1
				carditem := GetCsvMgr().GetItemConfig(cardid)
				if carditem != nil {
					if err := self.player.HasObjectOkEasy(cardid, carditem.CompoundNum); err == nil {
						have = true
					}
				}

				if have {
					continue
				}
			}

			if item.Look > self.player.GetModule("recharge").(*ModRecharge).Sql_UserRecharge.Money {
				continue
			}

			if group != self.Sql_Shop.info[i].Group {
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
					//self.checkMail(i, item)
					continue
				} else if self.Sql_Shop.info[i].Done == 2 {
					if self.Sql_Shop.info[i].Times < item.Times {
						self.Sql_Shop.info[i].Done = 1
						chg = true
					}
				} else {
					self.Sql_Shop.info[i].Done = 1
				}
			}

			if chg {
				self.chg = append(self.chg, self.Sql_Shop.info[i])
			}
		}
	}
}

//! 充值回调
func (self *ModLuckShop) DiscountHandleTask(tasktype, n2, n3, n4 int) {
	if self.Sql_Shop == nil {
		return
	}

	activity := GetActivityMgr().GetActivity(ActivityDiscountGift)
	if activity == nil || activity.status.Status != ACTIVITY_STATUS_OPEN {
		return
	}

	group := GetActivityMgr().getActN4(ActivityDiscountGift)
	for i := 0; i < len(self.Sql_Shop.info); i++ {
		item, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
		if !ok {
			continue
		}
		if ActivityDiscountGift != item.ActivityType {
			continue
		}

		if item.TaskTypes != tasktype {
			continue
		}

		if item.Look > self.player.GetModule("recharge").(*ModRecharge).Sql_UserRecharge.Money {
			continue
		}

		if group != self.Sql_Shop.info[i].Group {
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
				//self.checkMail(i, item)
				continue
			} else if self.Sql_Shop.info[i].Done == 2 {
				if self.Sql_Shop.info[i].Times < item.Times {
					self.Sql_Shop.info[i].Done = 1
					chg = true
				}
			} else {
				self.Sql_Shop.info[i].Done = 1
			}
		}

		if chg {
			self.chg = append(self.chg, self.Sql_Shop.info[i])
		}
	}
}

// 发送礼包信息
func (self *ModLuckShop) SendInfo() {
	self.CheckBox()
	var msg S2C_LuckShopInfo
	msg.Cid = "getluckshop2lst"

	group := GetActivityMgr().getActN4(ActivityLuckShop)
	now := TimeServer().Unix()

	for i := 0; i < len(self.Sql_Shop.info); i++ {
		shopItem, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
		if !ok {
			continue
		}

		if shopItem.ActivityType != ActivityLuckShop {
			continue
		}

		if group != self.Sql_Shop.info[i].Group {
			continue
		}

		if now < self.Sql_Shop.info[i].StartTime || now >= self.Sql_Shop.info[i].EndTime {
			continue
		}

		msg.Info = append(msg.Info, self.Sql_Shop.info[i])
		msg.Item = append(msg.Item, *shopItem)
	}

	group = GetActivityMgr().getActN4(ActivityStarGift)
	for i := 0; i < len(self.Sql_Shop.info); i++ {
		shopItem, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
		if !ok {
			continue
		}

		if shopItem.ActivityType != ActivityStarGift {
			continue
		}

		if group != self.Sql_Shop.info[i].Group {
			continue
		}

		if shopItem.PicType == 1 && shopItem.StarHero != 0 {
			have := false
			hero := self.player.GetModule("hero").(*ModHero).GetVoidHero(shopItem.StarHero)
			if hero != nil {
				have = true
			}

			// 特殊处理 12 + 英雄id + 01 是英雄的魂石id
			cardid := shopItem.StarHero*100 + 12000000 + 1
			carditem := GetCsvMgr().GetItemConfig(cardid)
			if carditem != nil {
				if err := self.player.HasObjectOkEasy(cardid, carditem.CompoundNum); err == nil {
					have = true
				}
			}

			if have {
				continue
			}
		}

		msg.Info = append(msg.Info, self.Sql_Shop.info[i])
		msg.Item = append(msg.Item, *shopItem)
	}

	group = GetActivityMgr().getActN4(ActivityDiscountGift)
	for i := 0; i < len(self.Sql_Shop.info); i++ {
		shopItem, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
		if !ok {
			continue
		}

		if shopItem.ActivityType != ActivityDiscountGift {
			continue
		}

		if group != self.Sql_Shop.info[i].Group {
			continue
		}

		msg.Info = append(msg.Info, self.Sql_Shop.info[i])
		msg.Item = append(msg.Item, *shopItem)
	}

	for id := ACT_STAR_GIFT_MIN; id < ACT_STAR_GIFT_MAX; id++ {
		for i := 0; i < len(self.Sql_Shop.info); i++ {
			shopItem, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
			if !ok {
				continue
			}

			if shopItem.ActivityType != id {
				continue
			}

			if group != self.Sql_Shop.info[i].Group {
				continue
			}

			msg.Info = append(msg.Info, self.Sql_Shop.info[i])
			msg.Item = append(msg.Item, *shopItem)
		}
	}

	// 防止客户端异常
	if len(msg.Info) == 0 {
		msg.Info = []JS_LuckShopItem{}
	}

	if len(msg.Item) == 0 {
		msg.Item = []JS_ActivityBox{}
	}

	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("getluckshop2lst", smsg)
}

//! 领取礼包奖励
func (self *ModLuckShop) DrawItem(id int) {

	timenow := TimeServer().Unix()
	// 先领取奖励
	index := -1
	for i := 0; i < len(self.Sql_Shop.info); i++ {
		if self.Sql_Shop.info[i].BoxId != id {
			continue
		}

		config, ok := GetActivityMgr().Sql_ActivityBox[id]
		if !ok {
			continue
		}

		if config.ActivityType == ActivityLuckShop  {
			if timenow < self.Sql_Shop.info[i].StartTime || timenow >= self.Sql_Shop.info[i].EndTime {
				continue
			}
		}

		if config.ActivityType >= ACT_STAR_GIFT_MIN && config.ActivityType <= ACT_STAR_GIFT_MAX {
			if timenow < self.Sql_Shop.info[i].StartTime || (timenow >= self.Sql_Shop.info[i].EndTime && (self.Sql_Shop.info[i].StartTime != self.Sql_Shop.info[i].EndTime)) {
				continue
			}
		}

		if self.Sql_Shop.info[i].Done == 0 && config.N[1] != 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THE_GIFT_PACKAGE_IS_NOT"))
			continue
		}

		if self.Sql_Shop.info[i].Done == 2 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THE_GIFT_PACKAGE_HAS_BEEN"))
			continue
		}
		index = i
		break
	}

	if index == -1 {
		self.sendNoItem()
		return
	}

	boxId := self.Sql_Shop.info[index].BoxId
	config, ok := GetActivityMgr().Sql_ActivityBox[boxId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THE_GIFT_PACKAGE_WAS_MISCONFIGURED"))
		return
	}

	if config.ActivityType == ActivityStarGift {
		if config.PicType == 1 && config.StarHero != 0 {
			have := false
			hero := self.player.GetModule("hero").(*ModHero).GetVoidHero(config.StarHero)
			if hero != nil {
				have = true
			}

			// 特殊处理 12 + 英雄id + 01 是英雄的魂石id
			cardid := config.StarHero*100 + 12000000 + 1
			carditem := GetCsvMgr().GetItemConfig(cardid)
			if carditem != nil {
				if err := self.player.HasObjectOkEasy(cardid, carditem.CompoundNum); err == nil {
					have = true
				}
			}

			if have {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THE_GIFT_PACKAGE_WAS_MISCONFIGURED"))
				return
			}
		}
		if config.Look > self.player.GetModule("recharge").(*ModRecharge).Sql_UserRecharge.Money {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THE_GIFT_PACKAGE_WAS_MISCONFIGURED"))
			return
		}
	}

	if self.Sql_Shop.info[index].Times >= config.Times {
		//不拦
		//return
	}
	out := make([]PassItem, 0)
	for j := 0; j < 6; j++ {
		itemid := config.Item[j]
		if itemid == 0 {
			break
		}
		out = append(out, PassItem{itemid, config.Num[j]})
	}

	logDec := "幸运礼包"
	switch config.Type {
	case ACTIVITY_GIFT_TYPE_DISCOUNT:
		logDec = "领取特惠礼包奖励"
	case ACTIVITY_GIFT_TYPE_STAR:
		logDec = "领取星辰礼包奖励"
	case ACTIVITY_GIFT_TYPE_STAR_HERO:
		logDec = "领取星辰英雄奖励"
	}

	for j := 0; j < len(out); j++ {
		out[j].ItemID, out[j].Num = self.player.AddObject(out[j].ItemID, out[j].Num, id, config.Sale, 0, logDec)
	}

	pInfo := &self.Sql_Shop.info[index]
	pInfo.Times += 1
	pInfo.Done = 2

	if pInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_LUCK_FLUSH_ERROR")+fmt.Sprintf("%d", config.Type)+
			GetCsvMgr().GetText("STR_LUCK_FLUSH_BOX")+fmt.Sprintf("%d", config.BoxId))
		return
	}

	var msg S2C_ActivityGet
	msg.Cid = "luckshopget"
	msg.Id = id
	msg.Item = out
	msg.NextInfo = *pInfo
	shopitem, ok := GetActivityMgr().Sql_ActivityBox[pInfo.BoxId]
	if ok {
		msg.NextItem = *shopitem
	}

	self.player.SendMsg("1", HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_LUCK_SHOP, id, config.Sale, 0, "幸运礼包", 0, 0, self.player)
}

func (self *ModLuckShop) sendNoItem() {
	var msg S2C_ActivityGet
	msg.Cid = "luckshopget"
	msg.Id = 0
	msg.Item = make([]PassItem, 0)
	self.player.SendMsg("1", HF_JtoB(&msg))

	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSHOP_THIS_COMMODITY_DOES_NOT_EXIST"))
}

// 发送礼包更新信息
func (self *ModLuckShop) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_LuckShopUpdate
	msg.Cid = "luckshop_update"
	msg.Info = self.chg
	self.chg = make([]JS_LuckShopItem, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("luckshop_update", smsg)
}

func (self *ModLuckShop) HasMonthCard() bool {
	return self.player.GetModule("activity").(*ModActivity).HasMonthCard()
}

// 合服发送幸运礼包
func (self *ModLuckShop) FitServer() {
	if self.Sql_Shop == nil {
		return
	}

	for i := 0; i < len(self.Sql_Shop.info); i++ {
		if self.Sql_Shop.info[i].Done == 1 {
			shopitem, ok := GetActivityMgr().Sql_ActivityBox[self.Sql_Shop.info[i].BoxId]
			if !ok {
				continue
			}

			out := make([]PassItem, 0)
			for j := 0; j < 6; j++ {
				itemid := shopitem.Item[j]
				if itemid == 0 {
					continue
				}
				out = append(out, PassItem{itemid, shopitem.Num[j]})
			}

			context := GetCsvMgr().GetText("STR_LUCK_SHOP_LEFT_AWARD")
			(self.player.GetModule("mail").(*ModMail)).AddMail(1, 1, 0, GetCsvMgr().GetText("STR_LUCK_SHOP_MAIL"), context, GetCsvMgr().GetText("STR_SYS"),
				out, true, TimeServer().Unix())
			self.Sql_Shop.info[i].Done = 2
			//log.Println("玩家：", self.player.Sql_UserBase.Uid, " 领取礼物!")
		}
	}
}

func (self *ModLuckShop) checkMail(index int, shopitem *JS_ActivityBox) {
	if shopitem == nil {
		return
	}

	out := make([]PassItem, 0)
	for j := 0; j < 6; j++ {
		itemid := shopitem.Item[j]
		if itemid == 0 {
			continue
		}
		out = append(out, PassItem{itemid, shopitem.Num[j]})
	}

	context := GetCsvMgr().GetText("STR_LUCK_SHOP_FRESH")
	(self.player.GetModule("mail").(*ModMail)).AddMail(1, 1, 0, GetCsvMgr().GetText("STR_LUCK_SHOP_MAIL"), context, GetCsvMgr().GetText("STR_SYS"),
		out, true, TimeServer().Unix())
}

func (self *ModLuckShop) HandleRecharge(grade int) int {
	for _, v := range self.Sql_Shop.info {
		config, _ := GetActivityMgr().Sql_ActivityBox[v.BoxId]
		if config != nil {
			if v.Done == SPECIAL_PURCHASE_STATE_CAN_BUY && config.N[1] == grade {
				self.DrawItem(v.BoxId)
				return config.Type
			}
		}
	}
	return 0
}
