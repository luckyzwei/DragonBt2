package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

const (
	SPECIAL_PURCHASE_LEVEL  = 1 // 等级条件
	SPECIAL_PURCHASE_PASS   = 2 // 通关条件
	SPECIAL_PURCHASE_TOWER0 = 3 // 爬塔1
	SPECIAL_PURCHASE_TOWER1 = 4 // 爬塔2
	SPECIAL_PURCHASE_TOWER2 = 5 // 爬塔3
	SPECIAL_PURCHASE_TOWER3 = 6 // 爬塔4
	SPECIAL_PURCHASE_TOWER4 = 7 // 爬塔5
)

const (
	SPECIAL_PURCHASE_TYPE_LIMIT = 1 // 限时礼包
	SPECIAL_PURCHASE_TYPE_2     = 2 //
)

const (
	MSG_SPECIAL_PURCHASE_INFO      = "special_purchase_info"      // 获得信息
	MSG_SPECIAL_PURCHASE_GET_AWARD = "special_purchase_get_award" // 获得奖励
	MSG_SPECIAL_PURCHASE_DONE      = "special_purchase_done"      // 获得奖励
)

const (
	SPECIAL_PURCHASE_STATE_NONE       = 0
	SPECIAL_PURCHASE_STATE_CAN_BUY    = 1 // 可以购买
	SPECIAL_PURCHASE_STATE_IS_GET     = 2 // 已经领取
	SPECIAL_PURCHASE_STATE_TIME_LIMIT = 3 // 已经过期
)

// 类数据结构
type SpecialPurchaseInfo struct {
	ID       int   `json:"id"`       // 第几条
	TaskType int   `json:"tasktype"` // 类型
	Plan     int   `json:"plan"`     // 数值
	Done     int   `json:"done"`     // 领取状态 0未达成 1 可以购买 2 已经领取
	EndTime  int64 `json:"endtime"`  // 结束时间
	GiftID   int   `json:"giftid"`   // 对应的礼包id
	Type     int   `json:"type"`
}

// 服务器结构体
type San_SpecialPurchase struct {
	Uid      int64  // 角色ID
	Info     string // 数据
	Recharge int    // 充值累计
	Sign     string // 礼包标记

	info []*SpecialPurchaseInfo // 数据
	sign map[int]int            // 礼包标记

	DataUpdate
}

// 限时抢购
type ModSpecialPurchase struct {
	player              *Player
	San_SpecialPurchase San_SpecialPurchase
}

func (self *ModSpecialPurchase) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_specialpurchase` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.San_SpecialPurchase, "san_specialpurchase", self.player.ID)

	if self.San_SpecialPurchase.Uid <= 0 {
		self.San_SpecialPurchase.Uid = self.player.ID

		self.Encode()
		InsertTable("san_specialpurchase", &self.San_SpecialPurchase, 0, true)
		self.San_SpecialPurchase.Init("san_specialpurchase", &self.San_SpecialPurchase, true)
	} else {
		self.Decode()
		self.San_SpecialPurchase.Init("san_specialpurchase", &self.San_SpecialPurchase, true)
	}

	if self.San_SpecialPurchase.sign == nil {
		self.San_SpecialPurchase.sign = make(map[int]int)
	}
}

func (self *ModSpecialPurchase) OnGetOtherData() {
	self.CheckOldTask()
}

func (self *ModSpecialPurchase) Decode() {
	json.Unmarshal([]byte(self.San_SpecialPurchase.Info), &self.San_SpecialPurchase.info)
	json.Unmarshal([]byte(self.San_SpecialPurchase.Sign), &self.San_SpecialPurchase.sign)
}

func (self *ModSpecialPurchase) Encode() {
	self.San_SpecialPurchase.Info = HF_JtoA(&self.San_SpecialPurchase.info)
	self.San_SpecialPurchase.Sign = HF_JtoA(&self.San_SpecialPurchase.sign)
}

// 存盘逻辑
func (self *ModSpecialPurchase) OnSave(sql bool) {
	self.Encode()
	self.San_SpecialPurchase.Update(sql)
}

// 消息处理
func (self *ModSpecialPurchase) OnMsg(ctrl string, body []byte) bool {
	return false
}

//发送信息
func (self *ModSpecialPurchase) SendInfo() {
	self.CheckOldTask()

	var msg S2C_SpecialPurchaseInfo

	msg.Cid = MSG_SPECIAL_PURCHASE_INFO
	now := TimeServer().Unix()

	send := make(map[int]int)
	for _, v := range self.San_SpecialPurchase.info {
		config := GetCsvMgr().GetSpecialPurchaseItemConfig(v.GiftID)
		if v.GiftID > 0 && nil == config {
			continue
		}

		taskConfig := GetCsvMgr().GetSpecialPurchaseConfig(v.ID)
		if nil == taskConfig {
			continue
		}
		if v.EndTime <= now {
			if v.Done == SPECIAL_PURCHASE_STATE_CAN_BUY {
				v.Done = SPECIAL_PURCHASE_STATE_TIME_LIMIT
				if config != nil {
					self.San_SpecialPurchase.Recharge -= config.MoneyReduce
					if self.San_SpecialPurchase.Recharge < 0 {
						self.San_SpecialPurchase.Recharge = 0
					}
				}
			}
			continue
		}

		if v.Done != SPECIAL_PURCHASE_STATE_CAN_BUY {
			continue
		}

		_, ok := send[taskConfig.Subtype]
		if ok {
			continue
		}

		send[taskConfig.Subtype] = 1
		msg.SpecialPurchaseInfo = append(msg.SpecialPurchaseInfo, v)
		msg.SpecialPurchaseConfig = append(msg.SpecialPurchaseConfig, GetCsvMgr().GetSpecialPurchaseItemConfig(v.GiftID))
	}

	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(MSG_SPECIAL_PURCHASE_INFO, smsg)
}

// 领取奖励
func (self *ModSpecialPurchase) GetAward(nID int) {
	data := self.GetSpecialPurchaseData(nID)

	if nil == data {
		return
	}

	if data.Done != SPECIAL_PURCHASE_STATE_CAN_BUY {
		return
	}

	config := GetCsvMgr().GetSpecialPurchaseItemConfig(data.GiftID)
	if nil == config {
		return
	}

	//nowTime := TimeServer().Unix()

	//// 已经过期
	//if nowTime >= data.EndTime {
	//	return
	//}

	var msg S2C_SpecialPurchaseGetAward

	msg.Cid = MSG_SPECIAL_PURCHASE_GET_AWARD
	msg.ID = nID
	msg.Items = self.player.AddObjectLst(config.Items, config.Nums, "购买限时礼包", nID, data.GiftID, 0)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SPECIAL_PURCHASE_BUY, config.MoneyID, config.Money, 0, "购买限时礼包", 0, self.player.GetVip(), self.player)

	data.Done = SPECIAL_PURCHASE_STATE_IS_GET

	msg.Progress = data

	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(MSG_SPECIAL_PURCHASE_GET_AWARD, smsg)

}

// 任务处理
func (self *ModSpecialPurchase) HandleTask(tasktype, n2, n3, n4 int) {
	// 不是升级或过关则 充值返回
	nowTime := TimeServer().Unix()
	configs := GetCsvMgr().ActivityBuyLimit

	var msg S2C_SpecialPurchaseDone

	msg.Cid = MSG_SPECIAL_PURCHASE_DONE

	for _, v := range configs {
		if v.TaskTypes != tasktype {
			continue
		}

		if v.Type == SPECIAL_PURCHASE_TYPE_2 {
			if v.TriggerNum <= self.San_SpecialPurchase.sign[v.ID] {
				continue
			}
			//检查时间
			startTime := GetServer().GetOpenServer() + int64(v.OpenTime-1)*86400
			endTime := GetServer().GetOpenServer() + int64(v.OverTime-1)*86400
			if nowTime < startTime || nowTime > endTime {
				continue
			}
		}

		switch tasktype {
		case TASK_TYPE_PLAYER_LEVEL:
			if self.player.GetLv() != v.N[0] {
				continue
			}
		case TASK_TYPE_FINISH_PASS:
			if n2 != v.N[1] || n3 == 0 {
				continue
			}
		case TASK_TYPE_WOTER_LEVEL:
			if n2 != v.N[0] {
				continue
			}
		case TASK_TYPE_ONE_CAMP_TOWER_LEVEL:
			if n2 != v.N[0] || n3 != v.N[1] {
				continue
			}
		case TASK_TYPE_GET_HOOK_FAST_AWARD:
			if n2 != v.N[0] || n3 != v.N[1] {
				continue
			}
		case TASK_TYPE_IS_LOGIN:
			if n2 != v.N[0] || n3 != v.N[1] {
				continue
			}
		}

		find := false
		for _, info := range self.San_SpecialPurchase.info {
			tempConfig := GetCsvMgr().GetSpecialPurchaseConfig(info.ID)
			if tempConfig != nil {
				if info.ID != v.ID &&
					tempConfig.Subtype == v.Subtype &&
					info.Done == SPECIAL_PURCHASE_STATE_CAN_BUY &&
					nowTime < info.EndTime {
					find = true
					break
				}
			}
		}
		if find {
			continue
		}

		data := self.GetSpecialPurchaseData(v.ID)
		// 获得数据
		if data == nil {
			tempData := SpecialPurchaseInfo{v.ID, v.TaskTypes, 0, SPECIAL_PURCHASE_STATE_NONE, 0, 0, v.Type}
			self.San_SpecialPurchase.info = append(self.San_SpecialPurchase.info, &tempData)
		}

		data = self.GetSpecialPurchaseData(v.ID)
		if data == nil {
			continue
		}

		if data.Done != SPECIAL_PURCHASE_STATE_NONE {
			continue
		}

		var tasknode TaskNode
		tasknode.Tasktypes = v.TaskTypes
		tasknode.N1 = v.N[0]
		tasknode.N2 = v.N[1]
		tasknode.N3 = v.N[2]
		tasknode.N4 = v.N[3]
		plan, add := DoTask(&tasknode, self.player, n2, n3, n4)
		if plan == 0 {
			continue
		}

		chg := false
		if add {
			data.Plan += plan
		} else {
			if plan > data.Plan {
				data.Plan = plan
			}
		}

		if data.Plan >= v.N[0] {
			data.Done = SPECIAL_PURCHASE_STATE_CAN_BUY
			data.EndTime = nowTime + v.LimitTime
			data.GiftID = self.GetGiftID(v.Group)
			chg = true

			if v.Type == SPECIAL_PURCHASE_TYPE_2 {
				self.San_SpecialPurchase.sign[v.ID]++
			}

			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SPECIAL_PURCHASE_ACTIVATE, data.ID, 0, 0, "触发限时礼包", 0, 0, self.player)
		}

		if chg {
			msg.Progress = append(msg.Progress, data)
			msg.Config = append(msg.Config, GetCsvMgr().GetSpecialPurchaseItemConfig(data.GiftID))
		}
	}

	if len(msg.Progress) > 0 {
		smsg, _ := json.Marshal(&msg)
		self.player.SendMsg(MSG_SPECIAL_PURCHASE_DONE, smsg)
	}

	return
}

// 获取配置
func (self *ModSpecialPurchase) GetSpecialPurchaseData(nIndex int) *SpecialPurchaseInfo {
	for _, v := range self.San_SpecialPurchase.info {
		if v.ID == nIndex {
			return v
		}
	}

	return nil
}

func (self *ModSpecialPurchase) HandleRecharge(grade int, giftId int) {
	timeNow := TimeServer().Unix()
	//虚拟充值先判断是不是所属礼包
	if giftId > 0 {
		for _, v := range GetCsvMgr().ActivityBuyItem {
			if v.ID == giftId && v.MoneyID == grade {
				for _, config := range GetCsvMgr().ActivityBuyLimit {
					if config.Group == v.Group {
						tempData := SpecialPurchaseInfo{config.ID, config.TaskTypes, 0, SPECIAL_PURCHASE_STATE_CAN_BUY, timeNow + DAY_SECS, giftId, SPECIAL_PURCHASE_TYPE_LIMIT}
						var msg S2C_SpecialPurchaseDone
						msg.Cid = MSG_SPECIAL_PURCHASE_DONE
						msg.Progress = append(msg.Progress, &tempData)
						msg.Config = append(msg.Config, GetCsvMgr().GetSpecialPurchaseItemConfig(tempData.GiftID))
						smsg, _ := json.Marshal(&msg)
						self.player.SendMsg(MSG_SPECIAL_PURCHASE_DONE, smsg)

						var msgAdd S2C_SpecialPurchaseGetAward
						msgAdd.Cid = MSG_SPECIAL_PURCHASE_GET_AWARD
						msgAdd.ID = config.ID
						msgAdd.Items = self.player.AddObjectLst(v.Items, v.Nums, "购买限时礼包", config.ID, tempData.GiftID, 0)
						msgAdd.Progress = &tempData
						GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SPECIAL_PURCHASE_BUY, v.MoneyID, v.Money, 0, "购买限时礼包", 0, self.player.GetVip(), self.player)
						smsgAdd, _ := json.Marshal(&msgAdd)
						self.player.SendMsg(MSG_SPECIAL_PURCHASE_GET_AWARD, smsgAdd)
						self.SendInfo()
						return
					}
				}
			}
		}
	}

	for _, v := range self.San_SpecialPurchase.info {
		config := GetCsvMgr().GetSpecialPurchaseItemConfig(v.GiftID)
		if config != nil {
			if v.Done == SPECIAL_PURCHASE_STATE_CAN_BUY && config.MoneyID == grade && timeNow < v.EndTime {
				self.GetAward(v.ID)
				self.SendInfo()
				break
			}
		}
	}
}

func (self *ModSpecialPurchase) GetGiftID(groupid int) int {
	count := self.San_SpecialPurchase.Recharge
	if count < 0 {
		count = 0
	}
	configs, ok := GetCsvMgr().ActivityMapBuyItem[groupid]
	if !ok {
		return 0
	}

	len := len(configs)
	for i := len - 1; i >= 0; i-- {
		if count >= configs[i].Money {
			return configs[i].ID
		}
	}

	return 0
}

// 老号保留下的特殊状态任务 激活没过期但是没有giftid 则删除 等他重新触发
func (self *ModSpecialPurchase) CheckOldTask() {
	now := TimeServer().Unix()
	len := len(self.San_SpecialPurchase.info)
	for i := len - 1; i >= 0; i-- {
		info := self.San_SpecialPurchase.info[i]
		// 不是可购买激活状态则跳过
		if info.Done != SPECIAL_PURCHASE_STATE_CAN_BUY {
			if info.Type == SPECIAL_PURCHASE_TYPE_2 {
				self.San_SpecialPurchase.info = append(self.San_SpecialPurchase.info[:i], self.San_SpecialPurchase.info[i+1:]...)
			}
			continue
		}

		if info.Type != SPECIAL_PURCHASE_TYPE_2 {
			if info.GiftID != 0 {
				continue
			}
			// 如果过期则跳过
			if info.EndTime <= now {
				continue
			}

			self.San_SpecialPurchase.info = append(self.San_SpecialPurchase.info[:i], self.San_SpecialPurchase.info[i+1:]...)
		}
	}

	return
}

func (self *ModSpecialPurchase) OnRefresh() {
	self.San_SpecialPurchase.sign = make(map[int]int)
	len := len(self.San_SpecialPurchase.info)
	for i := len - 1; i >= 0; i-- {
		info := self.San_SpecialPurchase.info[i]
		if info.Type == SPECIAL_PURCHASE_TYPE_2 {
			self.San_SpecialPurchase.info = append(self.San_SpecialPurchase.info[:i], self.San_SpecialPurchase.info[i+1:]...)
		}
	}
}
