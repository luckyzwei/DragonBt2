package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

const (
	LOTTERY_DRAW_TYPE_NORMAL = 1 //普通奖
	LOTTERY_DRAW_TYPE_PRIZE  = 2 //大奖

	LOTTERY_DRAW_GROUP_NORMAL = 1 //普通大奖
	LOTTERY_DRAW_GROUP_HIGH   = 2 //终极大奖
)

//
type LotteryDrawItem struct {
	Id       int        `json:"id"`    // Id
	Items    []PassItem `json:"items"` // 物品
	Type     int        `json:"typez"` // 物品
	NowTimes int        `json:"nowtimes"`
	MaxTimes int        `json:"maxtimes"` // 最大获得次数
}

//! 任务数据库
type San_LotteryDraw struct {
	Uid             int64
	KeyId           int //! 下次可以转的时间
	LotteryDrawinfo string
	AlreadyGet      string
	NowStage        int //! 当前阶段
	NowCount        int //! 当前阶段转到次数
	LowChoose       int //! 普通大奖选择
	HighChoose      int //! 终极大奖选择
	LuckValue       int //  当前幸运值

	lotteryDrawinfo []*LotteryDrawItem //! 奖品信息
	alreadyGet      map[int]int
	DataUpdate
}

//! 任务
type ModLotteryDraw struct {
	player          *Player
	Sql_LotteryDraw San_LotteryDraw
}

func (self *ModLotteryDraw) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_userlotterydraw` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_LotteryDraw, "san_userlotterydraw", self.player.ID)

	if self.Sql_LotteryDraw.Uid <= 0 {
		self.Sql_LotteryDraw.Uid = self.player.ID
		self.Sql_LotteryDraw.NowStage = 1
		self.Sql_LotteryDraw.NowCount = 0
		self.Sql_LotteryDraw.lotteryDrawinfo = make([]*LotteryDrawItem, 0)
		self.Sql_LotteryDraw.alreadyGet = make(map[int]int, 0)
		self.Sql_LotteryDraw.LuckValue = GetCsvMgr().getInitNum(LOTTERY_DRAW_INIT)
		self.Encode()
		InsertTable("san_userlotterydraw", &self.Sql_LotteryDraw, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_LotteryDraw.Init("san_userlotterydraw", &self.Sql_LotteryDraw, true)
}

//! 将数据库数据写入data
func (self *ModLotteryDraw) Decode() {
	json.Unmarshal([]byte(self.Sql_LotteryDraw.LotteryDrawinfo), &self.Sql_LotteryDraw.lotteryDrawinfo)
	json.Unmarshal([]byte(self.Sql_LotteryDraw.AlreadyGet), &self.Sql_LotteryDraw.alreadyGet)
}

//! 将data数据写入数据库
func (self *ModLotteryDraw) Encode() {
	self.Sql_LotteryDraw.LotteryDrawinfo = HF_JtoA(self.Sql_LotteryDraw.lotteryDrawinfo)
	self.Sql_LotteryDraw.AlreadyGet = HF_JtoA(self.Sql_LotteryDraw.alreadyGet)
}

func (self *ModLotteryDraw) OnGetOtherData() {

}

// 注册消息
func (self *ModLotteryDraw) onReg(handlers map[string]func(body []byte)) {
	handlers["dolotterydraw"] = self.DoLotteryDraw
	handlers["lotterydrawchangeprize"] = self.LotteryDrawChangePrize
	handlers["lotterydrawnext"] = self.LotteryDrawNext
	handlers["lotterydrawrecord"] = self.LotteryDrawRecord
}

func (self *ModLotteryDraw) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModLotteryDraw) OnSave(sql bool) {
	self.Encode()
	self.Sql_LotteryDraw.Update(sql)
}

func (self *ModLotteryDraw) SendInfo() {
	if !self.player.GetModule("activity").(*ModActivity).IsActivityOpen(ACT_LOTTERY_DRAW) {
		return
	}

	self.check()
	var msg S2C_LotteryDrawInfo
	msg.Cid = "lotterydrawinfo"
	msg.NowStage = self.Sql_LotteryDraw.NowStage
	msg.NowCount = self.Sql_LotteryDraw.NowCount
	msg.LotteryDrawInfo = self.Sql_LotteryDraw.lotteryDrawinfo
	msg.LowChoose = self.Sql_LotteryDraw.LowChoose
	msg.HighChoose = self.Sql_LotteryDraw.HighChoose
	msg.AlreadyGet = self.Sql_LotteryDraw.alreadyGet
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModLotteryDraw) DoLotteryDraw(body []byte) {
	if !self.player.GetModule("activity").(*ModActivity).IsActivityOpen(ACT_LOTTERY_DRAW) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ACT_NOT_OPEN"))
		return
	}

	//判断大奖是否已选择
	if self.IsNeedChoosePrize() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LOTTERYDRAW_CHOOSE_PRIZE"))
		return
	}

	var msg C2S_DoLotteryDraw
	json.Unmarshal(body, &msg)

	realTimes := 0
	isPrize := LOGIC_FALSE
	outItemMap := make(map[int]*Item)
	outItemDraw := make([]PassItem, 0)
	costItemMap := make(map[int]*Item)
	for i := 0; i < msg.Times; i++ {
		//看消耗够不够
		configCost := GetCsvMgr().GetTariffConfig2(TARIFF_TYPE_LOTTERY_DRAW_COST)
		if configCost == nil {
			break
		}
		if err := self.player.HasObjectOk(configCost.ItemIds, configCost.ItemNums); err != nil {
			if realTimes <= 0 {
				self.player.SendErrInfo("err", err.Error())
				return
			}
			break
		}
		AddItemMapHelper(costItemMap, configCost.ItemIds, configCost.ItemNums)
		AddItemMapHelper(outItemMap, configCost.GetItem, configCost.GetNum)
		//开始抽奖
		outItem, needBreak := self.CalDrop()
		realTimes++
		outItemDraw = append(outItemDraw, outItem...)
		if needBreak {
			isPrize = LOGIC_TRUE
			break
		}
	}

	//增加物品
	outItems := self.player.AddObjectItemMap(outItemMap, "活动抽奖", realTimes, 0, 0)
	outItemsDraw := self.player.AddObjectPassItem(outItemDraw, "活动抽奖", realTimes, 0, 0)
	//扣除物品
	costItems := self.player.RemoveObjectItemMap(costItemMap, "活动抽奖", realTimes, 0, 0)

	var msgRel S2C_DoLotteryDraw
	msgRel.Cid = "dolotterydraw"
	msgRel.RealTimes = realTimes
	msgRel.GetItems = outItems
	msgRel.GetItemsDraw = outItemsDraw
	msgRel.CostItems = costItems
	msgRel.LotteryDrawInfo = self.Sql_LotteryDraw.lotteryDrawinfo
	msgRel.NowCount = self.Sql_LotteryDraw.NowCount
	msgRel.AlreadyGet = self.Sql_LotteryDraw.alreadyGet
	msgRel.LowChoose = self.Sql_LotteryDraw.LowChoose
	msgRel.HighChoose = self.Sql_LotteryDraw.HighChoose
	msgRel.HasPrize = isPrize
	if isPrize == LOGIC_TRUE {
		msgRel.LotteryDrawRecordLow, msgRel.LotteryDrawRecordHigh = GetOfflineInfoMgr().GetLotteryDrawRecord()
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModLotteryDraw) LotteryDrawChangePrize(body []byte) {
	if !self.player.GetModule("activity").(*ModActivity).IsActivityOpen(ACT_LOTTERY_DRAW) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ACT_NOT_OPEN"))
		return
	}

	var msg C2S_LotteryDrawChangePrize
	json.Unmarshal(body, &msg)

	//先检查Id是否满足本层的大奖要求
	config, ok := GetCsvMgr().LotteryDrawConfigMap[msg.Id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	prizeGroup := 0
	if self.Sql_LotteryDraw.NowStage%10 == 0 {
		prizeGroup = LOTTERY_DRAW_GROUP_HIGH
	} else {
		prizeGroup = LOTTERY_DRAW_GROUP_NORMAL
	}

	if config.Group != prizeGroup {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LOTTERY_DRAW_PRIZE_NOT_OPEN"))
		return
	}

	if config.Layer != 0 && self.Sql_LotteryDraw.NowStage < config.Layer {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LOTTERY_DRAW_PRIZE_NOT_OPEN"))
		return
	}
	if config.Change != 0 && self.Sql_LotteryDraw.NowStage >= config.Change {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LOTTERY_DRAW_PRIZE_NOT_OPEN"))
		return
	}
	if self.Sql_LotteryDraw.alreadyGet[config.Id] >= config.Limit {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LOTTERY_DRAW_PRIZE_NOT_OPEN"))
		return
	}
	//如果勾选默认
	if msg.IsDefault == LOGIC_TRUE {
		if prizeGroup == LOTTERY_DRAW_GROUP_HIGH {
			self.Sql_LotteryDraw.HighChoose = msg.Id
		} else {
			self.Sql_LotteryDraw.LowChoose = msg.Id
		}
	}
	//检查之前大奖是否已经兑换,顺便替换成新的
	isFind := false
	for _, v := range self.Sql_LotteryDraw.lotteryDrawinfo {
		vConfig := GetCsvMgr().LotteryDrawConfigMap[v.Id]
		if vConfig == nil {
			continue
		}
		if vConfig.Type != LOTTERY_DRAW_TYPE_PRIZE {
			continue
		}
		if v.NowTimes > 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LOTTERY_DRAW_PRIZE_CANT_CHANEG"))
			return
		}

		newItem := make([]PassItem, 0)
		newItem = append(newItem, PassItem{ItemID: config.Items, Num: config.Nums})
		v.Items = newItem
		v.Id = config.Id
		isFind = true
		break
	}
	//没找到的话，新增
	if !isFind {
		prize := self.makeLotterDrawItem(config)
		if prize != nil {
			self.Sql_LotteryDraw.lotteryDrawinfo = append(self.Sql_LotteryDraw.lotteryDrawinfo, prize)
		}
	}

	var msgRel S2C_LotteryDrawChangePrize
	msgRel.Cid = "lotterydrawchangeprize"
	msgRel.LotteryDrawInfo = self.Sql_LotteryDraw.lotteryDrawinfo
	msgRel.HighChoose = self.Sql_LotteryDraw.HighChoose
	msgRel.LowChoose = self.Sql_LotteryDraw.LowChoose
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModLotteryDraw) LotteryDrawNext(body []byte) {
	if !self.player.GetModule("activity").(*ModActivity).IsActivityOpen(ACT_LOTTERY_DRAW) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ACT_NOT_OPEN"))
		return
	}

	//先看是否满足进入下层的条件:抽中过大奖，或无大奖的情况下全部抽光
	prizeOk := false
	allOK := true

	for _, v := range self.Sql_LotteryDraw.lotteryDrawinfo {
		config, ok := GetCsvMgr().LotteryDrawConfigMap[v.Id]
		if !ok {
			continue
		}
		if config.Type == LOTTERY_DRAW_TYPE_PRIZE && v.NowTimes >= v.MaxTimes {
			prizeOk = true
			break
		}
		if v.NowTimes < v.MaxTimes {
			allOK = false
		}
	}

	if !prizeOk && !allOK {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LOTTERY_DRAW_CANT_REFRESH"))
		return
	}

	self.Sql_LotteryDraw.NowStage++
	self.Sql_LotteryDraw.NowCount = 0
	self.Sql_LotteryDraw.lotteryDrawinfo = self.CalPrizePool()

	var msgRel S2C_LotteryDrawNext
	msgRel.Cid = "lotterydrawnext"
	msgRel.LotteryDrawInfo = self.Sql_LotteryDraw.lotteryDrawinfo
	msgRel.NowCount = self.Sql_LotteryDraw.NowCount
	msgRel.NowStage = self.Sql_LotteryDraw.NowStage
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModLotteryDraw) LotteryDrawRecord(body []byte) {
	if !self.player.GetModule("activity").(*ModActivity).IsActivityOpen(ACT_LOTTERY_DRAW) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ACT_NOT_OPEN"))
		return
	}

	var msgRel S2C_LotteryDrawRecord
	msgRel.Cid = "lotterydrawrecord"
	msgRel.LotteryDrawRecordLow, msgRel.LotteryDrawRecordHigh = GetOfflineInfoMgr().GetLotteryDrawRecord()
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModLotteryDraw) check() {

	activity := GetActivityMgr().GetActivity(ACT_LOTTERY_DRAW)
	if activity == nil {
		return
	}

	period := GetActivityMgr().getActN3(ACT_LOTTERY_DRAW)

	if self.Sql_LotteryDraw.KeyId != period {
		self.Sql_LotteryDraw.KeyId = period
		self.Sql_LotteryDraw.NowStage = 1
		self.Sql_LotteryDraw.NowCount = 0
		self.Sql_LotteryDraw.lotteryDrawinfo = make([]*LotteryDrawItem, 0)
		self.Sql_LotteryDraw.alreadyGet = make(map[int]int, 0)
	}

	if self.Sql_LotteryDraw.alreadyGet == nil {
		self.Sql_LotteryDraw.alreadyGet = make(map[int]int, 0)
	}

	if len(self.Sql_LotteryDraw.lotteryDrawinfo) == 0 {
		self.Sql_LotteryDraw.lotteryDrawinfo = self.CalPrizePool()
	}

	//容错检查
	for _, v := range self.Sql_LotteryDraw.lotteryDrawinfo {
		_, ok := GetCsvMgr().LotteryDrawConfigMap[v.Id]
		if !ok {
			self.Sql_LotteryDraw.lotteryDrawinfo = self.CalPrizePool()
			break
		}
	}
}

func (self *ModLotteryDraw) CalPrizePool() []*LotteryDrawItem {
	data := make([]*LotteryDrawItem, 0)

	prizeId := 0
	if self.Sql_LotteryDraw.NowStage%10 == 0 {
		prizeId = self.Sql_LotteryDraw.HighChoose
	} else {
		prizeId = self.Sql_LotteryDraw.LowChoose
	}

	for _, v := range GetCsvMgr().LotteryDrawConfigMap {
		if v.Type == LOTTERY_DRAW_TYPE_NORMAL || v.Id == prizeId {
			prize := self.makeLotterDrawItem(v)
			if prize != nil {
				data = append(data, prize)
			}
		}
	}
	return data
}

func (self *ModLotteryDraw) makeLotterDrawItem(config *LotteryDrawConfig) *LotteryDrawItem {

	switch config.Type {
	case LOTTERY_DRAW_TYPE_NORMAL:
		data := new(LotteryDrawItem)
		data.Id = config.Id
		data.NowTimes = 0
		data.MaxTimes = config.Limit
		data.Type = config.Type
		data.Items = append(data.Items, PassItem{ItemID: config.Items, Num: config.Nums})
		return data
	case LOTTERY_DRAW_TYPE_PRIZE:
		if config.Layer != 0 && self.Sql_LotteryDraw.NowStage < config.Layer {
			return nil
		}
		if config.Change != 0 && self.Sql_LotteryDraw.NowStage >= config.Change {
			return nil
		}
		if self.Sql_LotteryDraw.alreadyGet[config.Id] >= config.Limit {
			return nil
		}
		data := new(LotteryDrawItem)
		data.Id = config.Id
		data.NowTimes = 0
		data.MaxTimes = 1
		data.Type = config.Type
		data.Items = append(data.Items, PassItem{ItemID: config.Items, Num: config.Nums})
		return data
	}
	return nil
}

func (self *ModLotteryDraw) IsNeedChoosePrize() bool {

	prizeGroup := 0
	if self.Sql_LotteryDraw.NowStage%10 == 0 {
		prizeGroup = LOTTERY_DRAW_GROUP_HIGH
	} else {
		prizeGroup = LOTTERY_DRAW_GROUP_NORMAL
	}

	//如果奖池中已经有大将，则不需要在选择了
	for _, v := range self.Sql_LotteryDraw.lotteryDrawinfo {
		config := GetCsvMgr().LotteryDrawConfigMap[v.Id]
		if config != nil && config.Group == prizeGroup {
			return false
		}
	}

	//奖池中没有，但是配置里有可选的，则需要选择
	for _, v := range GetCsvMgr().LotteryDrawConfigMap {
		if v.Group != prizeGroup {
			continue
		}
		if v.Layer != 0 && self.Sql_LotteryDraw.NowStage < v.Layer {
			continue
		}
		if v.Change != 0 && self.Sql_LotteryDraw.NowStage >= v.Change {
			continue
		}
		if self.Sql_LotteryDraw.alreadyGet[v.Id] >= v.Limit {
			continue
		}
		return true
	}
	//配置中也没有，不需要选了
	return false
}

func (self *ModLotteryDraw) CalDrop() ([]PassItem, bool) {

	self.Sql_LotteryDraw.LuckValue += GetCsvMgr().getInitNum(LOTTERY_DRAW_ADD)

	allWeight := 0
	for _, v := range self.Sql_LotteryDraw.lotteryDrawinfo {
		if v.NowTimes >= v.MaxTimes {
			continue
		}
		allWeight += GetCsvMgr().LotteryDrawGetWeight(v.Id, self.Sql_LotteryDraw.LuckValue)
	}

	rand := self.player.GetModule("find").(*ModFind).GetRandInt(allWeight)
	nowRand := 0
	data := make([]PassItem, 0)
	needBreak := false
	for _, v := range self.Sql_LotteryDraw.lotteryDrawinfo {
		if v.NowTimes >= v.MaxTimes {
			continue
		}
		nowRand += GetCsvMgr().LotteryDrawGetWeight(v.Id, self.Sql_LotteryDraw.LuckValue)
		if nowRand > rand {
			self.Sql_LotteryDraw.LuckValue -= GetCsvMgr().LotteryDrawGetLucky(v.Id)
			self.Sql_LotteryDraw.NowCount++

			config, ok := GetCsvMgr().LotteryDrawConfigMap[v.Id]
			if !ok {
				continue
			}
			data = append(data, PassItem{ItemID: config.Items, Num: config.Nums})
			v.NowTimes++
			if config.Type == LOTTERY_DRAW_TYPE_PRIZE {
				self.Sql_LotteryDraw.alreadyGet[config.Id]++
				if self.Sql_LotteryDraw.alreadyGet[config.Id] >= config.Limit {
					if self.Sql_LotteryDraw.LowChoose == config.Id {
						self.Sql_LotteryDraw.LowChoose = 0
					}
					if self.Sql_LotteryDraw.HighChoose == config.Id {
						self.Sql_LotteryDraw.HighChoose = 0
					}
				}
				needBreak = true
				//处理公告
				record := new(LotteryDrawRecord)
				record.Uid = self.player.GetUid()
				record.Name = self.player.GetName()
				record.Times = self.Sql_LotteryDraw.NowCount
				record.ItemId = config.Items
				record.Num = config.Nums
				GetOfflineInfoMgr().AddLotteryDrawRecord(record, config.Notice)
			}
			break
		}
	}
	return data, needBreak
}
