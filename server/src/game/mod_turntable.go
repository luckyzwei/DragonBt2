package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

// 进度
type TurnTableItem struct {
	Id    int        `json:"id"`    // Id
	Items []PassItem `json:"items"` // 物品
	IsGet int        `json:"isget"` // 0未获得  1已获得
}

//! 任务数据库
type San_TurnTable struct {
	Uid           int64
	TurnTableinfo string
	NowStage      int   //! 当前阶段
	NowCount      int   //! 当前阶段转到次数
	NextTime      int64 //! 下次可以转的时间

	turnTableinfo []*TurnTableItem //! 任务信息
	DataUpdate
}

//! 任务
type ModTurnTable struct {
	player        *Player
	Sql_TurnTable San_TurnTable
}

func (self *ModTurnTable) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_userturntable` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_TurnTable, "san_userturntable", self.player.ID)

	if self.Sql_TurnTable.Uid <= 0 {
		self.Sql_TurnTable.Uid = self.player.ID
		self.Encode()
		InsertTable("san_userturntable", &self.Sql_TurnTable, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_TurnTable.Init("san_userturntable", &self.Sql_TurnTable, true)
}

//! 将数据库数据写入data
func (self *ModTurnTable) Decode() {
	json.Unmarshal([]byte(self.Sql_TurnTable.TurnTableinfo), &self.Sql_TurnTable.turnTableinfo)
}

//! 将data数据写入数据库
func (self *ModTurnTable) Encode() {
	s, _ := json.Marshal(&self.Sql_TurnTable.turnTableinfo)
	self.Sql_TurnTable.TurnTableinfo = string(s)
}

func (self *ModTurnTable) OnGetOtherData() {

}

// 注册消息
func (self *ModTurnTable) onReg(handlers map[string]func(body []byte)) {
	handlers["doturntable"] = self.DoTurnTable
}

func (self *ModTurnTable) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModTurnTable) OnSave(sql bool) {
	self.Encode()
	self.Sql_TurnTable.Update(sql)
}

func (self *ModTurnTable) SendInfo() {
	if !self.player.GetModule("activity").(*ModActivity).IsActivityOpen(ACT_TURNTABLE) {
		return
	}

	self.check()
	var msg S2C_TurnTableInfo
	msg.Cid = "turntableinfo"
	msg.NowStage = self.Sql_TurnTable.NowStage
	msg.NowCount = self.Sql_TurnTable.NowCount
	msg.NextTime = self.Sql_TurnTable.NextTime
	msg.TurnTableinfo = self.Sql_TurnTable.turnTableinfo
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModTurnTable) DoTurnTable(body []byte) {
	if !self.player.GetModule("activity").(*ModActivity).IsActivityOpen(ACT_TURNTABLE) {
		return
	}

	weightAll := 0
	for _, v := range self.Sql_TurnTable.turnTableinfo {
		config := GetCsvMgr().TurnTableConfigMap[v.Id]
		if config == nil || config.Stage != self.Sql_TurnTable.NowStage {
			continue
		}
		if v.IsGet == LOGIC_TRUE {
			continue
		}
		weightAll += config.Weight
	}

	if weightAll <= 0 {
		return
	}

	rate := HF_GetRandom(weightAll)
	nowRate := 0

	var msgRel S2C_DoTurnTable
	msgRel.Cid = "doturntable"
	for _, v := range self.Sql_TurnTable.turnTableinfo {
		config := GetCsvMgr().TurnTableConfigMap[v.Id]
		if config == nil || config.Stage != self.Sql_TurnTable.NowStage {
			continue
		}
		if v.IsGet == LOGIC_TRUE {
			continue
		}
		nowRate += config.Weight
		if nowRate > rate {
			v.IsGet = LOGIC_TRUE
			items := self.player.AddObjectPassItem(v.Items, "转盘活动", v.Id, 0, 0)
			msgRel.TurnTableItem = v
			msgRel.GetItems = append(msgRel.GetItems, items...)
			//计算冷却
			self.Sql_TurnTable.NowCount++
			configTime, ok := GetCsvMgr().TurnTableTimeConfigMap[self.Sql_TurnTable.NowCount]
			if ok {
				self.Sql_TurnTable.NextTime = TimeServer().Unix() + configTime.Time
				//self.Sql_TurnTable.NextTime = TimeServer().Unix() + 10
			}
			//看看STAGE 要不要切换
			isNeed := true
			for _, vv := range self.Sql_TurnTable.turnTableinfo {
				configvv := GetCsvMgr().TurnTableConfigMap[vv.Id]
				if configvv == nil || configvv.Stage != self.Sql_TurnTable.NowStage {
					continue
				}
				if vv.IsGet == LOGIC_FALSE {
					isNeed = false
					break
				}
			}
			if isNeed {
				self.Sql_TurnTable.NowStage++
			}
			break
		}
	}

	msgRel.NowStage = self.Sql_TurnTable.NowStage
	msgRel.NowCount = self.Sql_TurnTable.NowCount
	msgRel.NextTime = self.Sql_TurnTable.NextTime
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModTurnTable) check() {

	if len(self.Sql_TurnTable.turnTableinfo) == 0 {
		self.Sql_TurnTable.NextTime = TimeServer().Unix()
		self.Sql_TurnTable.NowCount = 0
		self.Sql_TurnTable.NowStage = 0

		for _, v := range GetCsvMgr().TurnTableConfigMap {
			turnTable := new(TurnTableItem)
			turnTable.Id = v.Id
			for i := 0; i < len(v.Items); i++ {
				turnTable.Items = append(turnTable.Items, PassItem{ItemID: v.Items[i], Num: v.Nums[i]})
			}
			turnTable.IsGet = LOGIC_FALSE
			self.Sql_TurnTable.turnTableinfo = append(self.Sql_TurnTable.turnTableinfo, turnTable)
		}
	}

	for _, vv := range self.Sql_TurnTable.turnTableinfo {
		configvv := GetCsvMgr().TurnTableConfigMap[vv.Id]
		if configvv == nil {
			continue
		}
		if vv.IsGet == LOGIC_FALSE {
			if self.Sql_TurnTable.NowStage == 0 || configvv.Stage < self.Sql_TurnTable.NowStage {
				self.Sql_TurnTable.NowStage = configvv.Stage
			}
		}
	}
}
