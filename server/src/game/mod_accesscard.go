package game

import (
	"encoding/json"
	"fmt"
	"strconv"

	//"time"
)

// 收藏家积分奖励
type AccessCardAward struct {
	AwardId    int   `json:"awardid"`    // 奖励ID
	Group      int   `json:"group"`      // 进度
	NeedPoint  int   `json:"needpoint"`  // 任务类型
	Pickup     int   `json:"pickup"`     // 是否领取奖励
	Items      []int `json:"items"`      // 奖励ID
	Nums       []int `json:"nums"`       // 奖励数量
	ChangeTime int64 `json:"changetime"` // 奖励更换时间， 0表示不更换
	Notice     int   `json:"notice"`     // 是否通知
}

//! 任务数据库
type San_AccessCard struct {
	Uid       int64
	TaskInfo  string
	AwardInfo string
	N3        int //!
	NGroup    int //!

	taskInfo  []*JS_TaskInfo     //! 任务信息
	awardInfo []*AccessCardAward //! 积分奖励信息

	DataUpdate
}

//! 任务
type ModAccessCard struct {
	player         *Player
	Sql_AccessCard San_AccessCard

	chg []JS_TaskInfo
}

func (self *ModAccessCard) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_useraccesscard` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_AccessCard, "san_useraccesscard", self.player.ID)

	if self.Sql_AccessCard.Uid <= 0 {
		self.Sql_AccessCard.Uid = self.player.ID
		self.Sql_AccessCard.taskInfo = make([]*JS_TaskInfo, 0)
		self.Sql_AccessCard.awardInfo = make([]*AccessCardAward, 0)
		self.Encode()
		InsertTable("san_useraccesscard", &self.Sql_AccessCard, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_AccessCard.Init("san_useraccesscard", &self.Sql_AccessCard, true)
}

//! 将数据库数据写入data
func (self *ModAccessCard) Decode() {
	json.Unmarshal([]byte(self.Sql_AccessCard.TaskInfo), &self.Sql_AccessCard.taskInfo)
	json.Unmarshal([]byte(self.Sql_AccessCard.AwardInfo), &self.Sql_AccessCard.awardInfo)
}

//! 将data数据写入数据库
func (self *ModAccessCard) Encode() {
	self.Sql_AccessCard.TaskInfo = HF_JtoA(&self.Sql_AccessCard.taskInfo)
	self.Sql_AccessCard.AwardInfo = HF_JtoA(&self.Sql_AccessCard.awardInfo)
}

func (self *ModAccessCard) OnGetOtherData() {

}

// 注册消息
func (self *ModAccessCard) onReg(handlers map[string]func(body []byte)) {
	handlers["accesscardtask"] = self.AccessCardTask
	handlers["accesscardaward"] = self.AccessCardAward
	handlers["getaccesscardrecord"] = self.GetAccessCardRecord
	handlers["accesscardall"] = self.AccessCardAll
	handlers["accessgetrank"] = self.AccessGetRank
}

func (self *ModAccessCard) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModAccessCard) OnSave(sql bool) {
	self.Encode()
	self.Sql_AccessCard.Update(sql)
}

func (self *ModAccessCard) SendInfo() {

	/*
		config1 := GetCsvMgr().AccessAwardConfig
		config := config1[100106]
		string := fmt.Sprintf(config.Notice[0], self.player.Sql_UserBase.UName, 1, "测试物品")
		GetServer().Notice(string, 0, 0)
	*/

	isOpen, index := GetActivityMgr().JudgeOpenAllIndex(ACT_ACCESSCARD_MIN, ACT_ACCESSCARD_MAX)
	if !isOpen {
		return
	}

	N3 := GetActivityMgr().getActN3(index)
	group := GetActivityMgr().getActN4(index)
	now := TimeServer().Unix()
	//重置
	if self.Sql_AccessCard.N3 != N3 || self.Sql_AccessCard.NGroup != group {
		self.Sql_AccessCard.N3 = N3
		self.Sql_AccessCard.NGroup = group

		self.Sql_AccessCard.taskInfo = make([]*JS_TaskInfo, 0)
		self.Sql_AccessCard.awardInfo = make([]*AccessCardAward, 0)

		//添加积分
		for _, v := range GetCsvMgr().AccessAwardConfig {
			if v.Group != self.Sql_AccessCard.NGroup {
				continue
			}
			data := new(AccessCardAward)
			data.AwardId = v.Id
			data.Group = v.Group
			data.NeedPoint = v.Point
			data.Items = v.Item
			data.Nums = v.Num
			if v.Countdown > 0 {
				startTime := self.player.GetModule("activity").(*ModActivity).GetActivityStart(index)
				data.ChangeTime = startTime + v.Countdown
				if data.ChangeTime < now {
					data.Items = v.ConversionItem
					data.Nums = v.ConversionNum
					data.ChangeTime = 0
				}
			}
			self.Sql_AccessCard.awardInfo = append(self.Sql_AccessCard.awardInfo, data)
		}

		//历史记录重置
		num := self.player.GetObjectNum(ITEM_ACCESS_ITEM)
		if num > 0 {
			self.player.RemoveObjectSimple(ITEM_ACCESS_ITEM, num, "收藏家刷新", self.Sql_AccessCard.N3, self.Sql_AccessCard.NGroup, 0)
		}
	}

	//做配置的同步处理
	nowMap := make(map[int]*JS_TaskInfo)
	for _, v := range self.Sql_AccessCard.taskInfo {
		nowMap[v.Taskid] = v
	}
	//添加任务
	self.Sql_AccessCard.taskInfo = make([]*JS_TaskInfo, 0)
	for _, v := range GetCsvMgr().AccessTaskConfig {
		if v.Group != self.Sql_AccessCard.NGroup {
			continue
		}
		data := new(JS_TaskInfo)
		data.Taskid = v.Id
		data.Tasktypes = v.TaskTypes
		_, ok := nowMap[v.Id]
		if ok {
			data.Plan = nowMap[v.Id].Plan
			data.Finish = nowMap[v.Id].Finish
			data.Pickup = nowMap[v.Id].Pickup
		}
		self.Sql_AccessCard.taskInfo = append(self.Sql_AccessCard.taskInfo, data)
	}

	for _, v := range self.Sql_AccessCard.awardInfo {
		if v.ChangeTime != 0 && v.ChangeTime < now {
			config := GetCsvMgr().AccessAwardConfig[v.AwardId]
			if config != nil {
				v.Items = config.ConversionItem
				v.Nums = config.ConversionNum
				v.ChangeTime = 0
			}
		}
	}

	var msg S2C_AccessCardInfo
	msg.Cid = "accesscardinfo"
	msg.Group = self.Sql_AccessCard.NGroup
	msg.TaskInfo = self.Sql_AccessCard.taskInfo
	msg.AwardInfo = self.Sql_AccessCard.awardInfo
	msg.Point = self.player.GetObjectNum(ITEM_ACCESS_ITEM)
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	if msg.Point > 0 {
		GetAccessCardRecordMgr().UpdatePoint(self.player, msg.Point)
	}
}

func (self *ModAccessCard) HandleTask(tasktype, n2, n3, n4 int) {
	isOpen, _ := GetActivityMgr().JudgeOpenIndex(ACT_ACCESSCARD_MIN, ACT_ACCESSCARD_MAX)
	if !isOpen {
		return
	}
	for _, node := range self.Sql_AccessCard.taskInfo {
		if node.Tasktypes != tasktype {
			continue
		}
		if node.Finish == LOGIC_TRUE {
			continue
		}
		config := GetCsvMgr().AccessTaskConfig[node.Taskid]
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

		chg := false
		if add {
			node.Plan += plan
			chg = true
		} else {
			if tasktype == PvpRankNow {
				if plan != 0 { // 新排名为0则直接不处理
					if node.Plan == 0 { // 进度为0则说明未初始化 直接赋值
						node.Plan = plan
						chg = true
					} else { // 进度不为0 则需要判断 获得的新名次比之前要高 则赋值
						if plan < node.Plan {
							node.Plan = plan
							chg = true
						}
					}
				}
			} else {
				if plan > node.Plan {
					node.Plan = plan
					chg = true
				}
			}
		}

		if node.Plan >= config.Ns[0] {
			node.Finish = LOGIC_TRUE
		}

		if chg {
			self.chg = append(self.chg, *node)
		}
	}
}

// 发送礼包更新信息
func (self *ModAccessCard) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_TaskUpdate
	msg.Cid = "accesscardupdate"
	msg.Info = self.chg
	self.chg = make([]JS_TaskInfo, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg(msg.Cid, smsg)
}

func (self *ModAccessCard) AccessCardTask(body []byte) {

	oldNum := self.player.GetObjectNum(ITEM_ACCESS_ITEM)
	var msg C2S_AccessCardTask
	json.Unmarshal(body, &msg)

	node := self.GetTask(msg.Taskid)
	if node == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	if node.Pickup == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_AWARD_HAS_BEEN_RECEIVED"))
		return
	}

	if node.Finish != LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_NOT_COMPLETED"))
		return
	}

	config := GetCsvMgr().AccessTaskConfig[msg.Taskid]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	newItem := make([]int, 0)
	newNum := make([]int, 0)
	for i := 0; i < len(config.Item); i++ {
		//策划说写死的  我不同意
		if i >= 2 && !GetAccessCardRecordMgr().IsCanGet() {
			break
		}
		if config.Item[i] == 0 {
			continue
		}
		newItem = append(newItem, config.Item[i])
		newNum = append(newNum, config.Num[i])
	}

	//发送任务奖励
	items := self.player.AddObjectLst(newItem, newNum, "领取英雄收藏家奖励", node.Taskid, 0, 0)
	node.Pickup = LOGIC_TRUE

	num := self.player.GetObjectNum(ITEM_ACCESS_ITEM)
	self.CheckNotice(num, config.Item[0])

	var msgRel S2C_AccessCardTask
	msgRel.Cid = "accesscardtask"
	msgRel.GetItems = items
	msgRel.TaskInfo = node
	msgRel.Point = num
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ACCESSCARD_GET, node.Taskid, num, oldNum, "领取英雄收藏家奖励", 0, 0, self.player)
}

func (self *ModAccessCard) AccessCardAward(body []byte) {
	var msg C2S_AccessCardAward
	json.Unmarshal(body, &msg)

	node := self.GetAward(msg.Id)
	if node == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	if node.Pickup == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_AWARD_HAS_BEEN_RECEIVED"))
		return
	}

	config := GetCsvMgr().AccessAwardConfig[msg.Id]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	num := self.player.GetObjectNum(ITEM_ACCESS_ITEM)
	if num < node.NeedPoint {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_STATISTICS_SCORE_NOT_ENOUGH"))
		return
	}

	if node.ChangeTime != 0 && node.ChangeTime < TimeServer().Unix() {
		node.Items = config.ConversionItem
		node.Nums = config.ConversionNum
		node.ChangeTime = 0
	}
	//发送任务奖励
	items := self.player.AddObjectLst(node.Items, node.Nums, "收藏家积分奖励", node.AwardId, num, 0)
	node.Pickup = LOGIC_TRUE

	self.CheckNoticeByAward(msg.Id, node.Items[0])

	var msgRel S2C_AccessCardAward
	msgRel.Cid = "accesscardaward"
	msgRel.GetItems = items
	msgRel.AccessCardAward = node
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ACCESSCARD_SCORE_GET, node.AwardId, num, 0, "领取英雄收藏家积分奖励", 0, 0, self.player)
}

func (self *ModAccessCard) GetAccessCardRecord(body []byte) {
	isOpen, _ := GetActivityMgr().JudgeOpenAllIndex(ACT_ACCESSCARD_MIN, ACT_ACCESSCARD_MAX)
	if !isOpen {
		return
	}

	var msgRel S2C_GetAccessCardRecord
	msgRel.Cid = "getaccesscardrecord"
	msgRel.Record = GetAccessCardRecordMgr().GetRecord()
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModAccessCard) AccessGetRank(body []byte) {
	GetAccessCardRecordMgr().GetRank(self.player)
}

func (self *ModAccessCard) AccessCardAll(body []byte) {
	self.SendInfo()
}

func (self *ModAccessCard) GetTask(id int) *JS_TaskInfo {
	for _, v := range self.Sql_AccessCard.taskInfo {
		if v.Taskid == id {
			return v
		}
	}
	return nil
}

func (self *ModAccessCard) GetAward(id int) *AccessCardAward {
	for _, v := range self.Sql_AccessCard.awardInfo {
		if v.AwardId == id {
			return v
		}
	}
	return nil
}

func (self *ModAccessCard) CheckNotice(num int, itemId int) {
	for _, v := range self.Sql_AccessCard.awardInfo {
		if v.Notice == LOGIC_TRUE {
			continue
		}
		if v.NeedPoint > num {
			continue
		}
		config := GetCsvMgr().AccessAwardConfig[v.AwardId]
		if config == nil {
			continue
		}
		if config.Notice[1] == "" {
			continue
		}
		v.Notice = LOGIC_TRUE

		itemConfig := GetCsvMgr().ItemMap[itemId]
		name := ""
		if itemConfig != nil {
			name = itemConfig.ItemName
		}
		GetAccessCardRecordMgr().AddRecord(self.player, num, name, 1)
	}
	return
}

func (self *ModAccessCard) CheckNoticeByAward(id int, itemId int) {
	for _, v := range self.Sql_AccessCard.awardInfo {
		if v.Notice == LOGIC_TRUE {
			continue
		}
		if v.AwardId != id {
			continue
		}
		config := GetCsvMgr().AccessAwardConfig[v.AwardId]
		if config == nil {
			continue
		}
		if config.Notice[0] == "0" {
			return
		}
		v.Notice = LOGIC_TRUE

		itemConfig := GetCsvMgr().ItemMap[itemId]
		name := ""
		if itemConfig != nil {
			name = itemConfig.ItemName
		}
		rank := GetAccessCardRecordMgr().AddRecord(self.player, config.Point, name, 0)
		GetServer().Notice(fmt.Sprintf(config.Notice[0], self.player.Sql_UserBase.UName, strconv.Itoa(rank), name), 0, 0)
	}
	return
}
