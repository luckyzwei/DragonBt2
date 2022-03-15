//星座系统
package game

import (
	"encoding/json"
	"fmt"
)

const (
	VALUE_MAX     = 1
	VALUE_PER_ADD = 2
	VALUE__ADD    = 3
)

const (
	PRIVILEGE_GOLD     = 101
	PRIVILEGE_HERO_EXP = 102
	PRIVILEGE_POWDER   = 103
	PRIVILEGE_PIT      = 104
	PRIVILEGE_UNION    = 105

	PRIVILEGE_ONHOOK_MIN = 1
	PRIVILEGE_ONHOOK_MAX = 13

	PRIVILEGE_ADD_SELFFIND_COUNT     = 201
	PRIVILEGE_REDUCE_SELFFIND_OFFSET = 202
)

//
type Galaxy struct {
	GalaxyId   int             `json:"galaxyid"`   //星系id
	NebulaInfo map[int]*Nebula `json:"nebulainfo"` //星云信息
}

type Nebula struct {
	NebulaId      int                `json:"nebulaid"`      //星云id
	NebulaTask    []*JS_TaskInfo     `json:"nebulatask"`    //星云解锁任务
	NebulaState   int                `json:"nebulastate"`   //星云解锁状态  0未解锁  1解锁
	NebulaWarInfo map[int]*NebulaWar `json:"nebulawarinfo"` //星点信息
}

type NebulaWar struct {
	NebulaWarId      int          `json:"nebulawarid"`      //星点id
	NebulaWarTask    *JS_TaskInfo `json:"nebulawartask"`    //星点任务
	NebulaWarState   int          `json:"nebulawarstate"`   //星点状态  0未解锁  1解锁
	NebulaWarBoxSign map[int]int  `json:"nebulawarboxsign"` //星点奖励领取状态,全部领取才解锁
}

//! 任务数据库
type San_InterStellar struct {
	Uid            int64
	GalaxyInfo     string //!
	PrivilegeValue string //!  特权值
	StellarCount   int    //!  解锁的星点数量

	galaxyInfo     map[int]*Galaxy //! 任务信息
	privilegeValue map[int]int     //! 特权值
	DataUpdate
}

//!
type ModInterStellar struct {
	player           *Player
	Sql_InterStellar San_InterStellar
	chgNebula        []*JS_TaskInfo
	chgNebulaWar     []*JS_TaskInfo
}

func (self *ModInterStellar) OnGetData(player *Player) {
	self.player = player
	self.chgNebula = make([]*JS_TaskInfo, 0)
	self.chgNebulaWar = make([]*JS_TaskInfo, 0)
}

func (self *ModInterStellar) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_interstellar` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_InterStellar, "san_interstellar", self.player.ID)

	if self.Sql_InterStellar.Uid <= 0 {
		self.Sql_InterStellar.Uid = self.player.ID
		self.Sql_InterStellar.galaxyInfo = make(map[int]*Galaxy, 0)
		self.Sql_InterStellar.privilegeValue = make(map[int]int, 0)
		self.Encode()
		InsertTable("san_interstellar", &self.Sql_InterStellar, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_InterStellar.Init("san_interstellar", &self.Sql_InterStellar, true)

	self.Check()
}

func (self *ModInterStellar) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModInterStellar) onReg(handlers map[string]func(body []byte)) {
	handlers["unlocknebula"] = self.UnlockNebula
	handlers["getnebulawarbox"] = self.GetNebulaWarBox
}

func (self *ModInterStellar) OnSave(sql bool) {
	self.Encode()
	self.Sql_InterStellar.Update(sql)
}

//! 将数据库数据写入data
func (self *ModInterStellar) Decode() {
	json.Unmarshal([]byte(self.Sql_InterStellar.GalaxyInfo), &self.Sql_InterStellar.galaxyInfo)
	json.Unmarshal([]byte(self.Sql_InterStellar.PrivilegeValue), &self.Sql_InterStellar.privilegeValue)
}

//! 将data数据写入数据库
func (self *ModInterStellar) Encode() {
	self.Sql_InterStellar.GalaxyInfo = HF_JtoA(&self.Sql_InterStellar.galaxyInfo)
	self.Sql_InterStellar.PrivilegeValue = HF_JtoA(&self.Sql_InterStellar.privilegeValue)
}

//! 处理任务
func (self *ModInterStellar) HandleTask(tasktype, n2, n3, n4 int) {

	for _, galaxy := range self.Sql_InterStellar.galaxyInfo {
		for _, nebula := range galaxy.NebulaInfo {
			//判断星云任务
			for i := 0; i < len(nebula.NebulaTask); i++ {
				if nebula.NebulaTask[i].Tasktypes != tasktype {
					continue
				}

				if nebula.NebulaTask[i].Finish == LOGIC_TRUE {
					continue
				}

				taskNode := GetCsvMgr().InterstellarTaskNode[nebula.NebulaTask[i].Taskid]
				if taskNode == nil {
					continue
				}

				plan, add := DoTask(taskNode, self.player, n2, n3, n4)
				if plan == 0 {
					continue
				}

				chg := false
				if add {
					nebula.NebulaTask[i].Plan += plan
					chg = true
				} else {
					if plan > nebula.NebulaTask[i].Plan {
						nebula.NebulaTask[i].Plan = plan
						chg = true
					}
				}

				if nebula.NebulaTask[i].Plan >= taskNode.N1 {
					nebula.NebulaTask[i].Finish = LOGIC_TRUE
					chg = true
				}

				if chg {
					self.chgNebula = append(self.chgNebula, nebula.NebulaTask[i])
				}
			}

			//判断星点任务
			for _, war := range nebula.NebulaWarInfo {
				if war.NebulaWarTask == nil {
					continue
				}

				if war.NebulaWarTask.Tasktypes != tasktype {
					continue
				}

				if war.NebulaWarTask.Finish == LOGIC_TRUE {
					continue
				}

				taskNode := GetCsvMgr().InterstellarTaskNode[war.NebulaWarTask.Taskid]
				if taskNode == nil {
					continue
				}

				plan, add := DoTask(taskNode, self.player, n2, n3, n4)
				if plan == 0 {
					continue
				}

				chg := false
				if add {
					war.NebulaWarTask.Plan += plan
					chg = true
				} else {
					if plan > war.NebulaWarTask.Plan {
						war.NebulaWarTask.Plan = plan
						chg = true
					}
				}

				if war.NebulaWarTask.Plan >= taskNode.N1 {
					war.NebulaWarTask.Finish = LOGIC_TRUE
					chg = true
				}

				if chg {
					self.chgNebulaWar = append(self.chgNebulaWar, war.NebulaWarTask)
				}
			}
		}
	}
}

func (self *ModInterStellar) UnlockNebula(body []byte) {
	var msg C2S_UnlockNebula
	json.Unmarshal(body, &msg)

	configNebula := GetCsvMgr().InterstellarConfig[msg.NebulaId]
	if configNebula == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	_, okinterstellar := self.Sql_InterStellar.galaxyInfo[configNebula.Galaxy]
	if !okinterstellar {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	info, okNebula := self.Sql_InterStellar.galaxyInfo[configNebula.Galaxy].NebulaInfo[configNebula.Id]
	if !okNebula {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	if info.NebulaState == LOGIC_TRUE {
		self.player.SendErr(GetCsvMgr().GetText("STR_UIPUBLIC_IS_OPEN"))
		return
	}

	for i := 0; i < len(info.NebulaTask); i++ {
		if info.NebulaTask[i].Finish == LOGIC_FALSE {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_NOT_COMPLETED"))
			return
		}
	}

	info.NebulaState = LOGIC_TRUE

	var msgRel S2C_UnlockNebula
	msgRel.Cid = "unlocknebula"
	msgRel.NebulaId = info.NebulaId
	msgRel.NebulaState = info.NebulaState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_INTERSTELLAR_NEBULA, info.NebulaId, 0, 0, "激活星云", 0, 0, self.player)

}

func (self *ModInterStellar) GetNebulaWarBox(body []byte) {
	var msg C2S_GetNebulaWarBox
	json.Unmarshal(body, &msg)

	configNebula := GetCsvMgr().InterstellarConfig[msg.NebulaId]
	if configNebula == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	_, okinterstellar := self.Sql_InterStellar.galaxyInfo[configNebula.Galaxy]
	if !okinterstellar {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	info, okNebula := self.Sql_InterStellar.galaxyInfo[configNebula.Galaxy].NebulaInfo[configNebula.Id]
	if !okNebula {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	if info.NebulaState == LOGIC_FALSE {
		self.player.SendErr(GetCsvMgr().GetText("STR_UI_FUNCTION_NOT_PUBLIC"))
		return
	}

	group := info.NebulaWarInfo[msg.GroupId]
	if group == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	if group.NebulaWarTask.Finish == LOGIC_FALSE {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_NOT_COMPLETED"))
		return
	}

	if info.NebulaWarInfo[msg.GroupId].NebulaWarBoxSign[msg.BoxId] == LOGIC_TRUE {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MONEYTASK_MISSION_COMPLETED"))
		return
	}

	configNebulaWarTemp := GetCsvMgr().InterstellarWar[msg.NebulaId]
	if configNebulaWarTemp == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}
	configNebulaWar := configNebulaWarTemp[msg.GroupId]
	if configNebulaWar == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	configTemp := GetCsvMgr().InterstellarBox[msg.GroupId]
	if configTemp == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	config := configTemp[msg.BoxId]
	if config == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MILITARY_MISSION_DOES_NOT_EXIST_AND"))
		return
	}

	info.NebulaWarInfo[msg.GroupId].NebulaWarBoxSign[msg.BoxId] = LOGIC_TRUE

	items := self.player.AddObjectLst(config.Item, config.Num, "领取宝箱奖励", msg.BoxId, info.NebulaWarInfo[msg.GroupId].NebulaWarId, 0)

	isUp := true
	for _, v := range info.NebulaWarInfo[msg.GroupId].NebulaWarBoxSign {
		if v == LOGIC_FALSE {
			isUp = false
		}
	}

	if isUp {
		info.NebulaWarInfo[msg.GroupId].NebulaWarState = LOGIC_TRUE

		if self.Sql_InterStellar.privilegeValue == nil {
			self.Sql_InterStellar.privilegeValue = make(map[int]int, 0)
		}
		_, ok := self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType]

		self.SetStellarCount(self.Sql_InterStellar.StellarCount + 1)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_INTERSTELLAR_NEBULAWAR, info.NebulaWarInfo[msg.GroupId].NebulaWarId, self.Sql_InterStellar.StellarCount+1, 0, "激活星耀", 0, 0, self.player)

		switch configNebulaWar.Type {
		case VALUE_MAX:
			if ok {
				if configNebulaWar.PrivileGeNum > self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType] {
					self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType] = configNebulaWar.PrivileGeNum
				}
			} else {
				self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType] = configNebulaWar.PrivileGeNum
			}
			//增加激活奖励到挂机掉落中
			configHangup := GetCsvMgr().InterstellarHangup[configNebulaWar.PrivileGeNum]
			if configHangup != nil {
				self.player.GetModule("onhook").(*ModOnHook).AddPrivilegeExtItems(configHangup.AddItem, configHangup.AddNum)
			}
		case VALUE_PER_ADD, VALUE__ADD:
			self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType] += configNebulaWar.PrivileGeNum
			if configNebulaWar.PrivileGeType == PRIVILEGE_ADD_SELFFIND_COUNT || configNebulaWar.PrivileGeType == PRIVILEGE_REDUCE_SELFFIND_OFFSET {
				self.player.GetModule("find").(*ModFind).CalSelfFind(self.Sql_InterStellar.privilegeValue, true)
			}
		}
		self.player.HandleTask(TASK_TYPE_INTERSTELLAR_NEBULAWAR, info.NebulaWarInfo[msg.GroupId].NebulaWarId, 0, 0)

		//看是否所有星耀都完成，则触发星云完成
		isFinish := true
		for _, nebulaWar := range info.NebulaWarInfo {
			if nebulaWar.NebulaWarState == LOGIC_FALSE {
				isFinish = false
				break
			}
		}
		if isFinish {
			self.player.HandleTask(TASK_TYPE_INTERSTELLAR_NEBULA, info.NebulaId, 0, 0)
		}
	}

	var msgRel S2C_GetNebulaWarBox
	msgRel.Cid = "getnebulawarbox"
	msgRel.NebulaId = msg.NebulaId
	msgRel.GroupId = msg.GroupId
	msgRel.BoxId = msg.BoxId
	msgRel.GetItems = items
	msgRel.NebulaWar = info.NebulaWarInfo[msg.GroupId]
	msgRel.PrivilegeValue = self.Sql_InterStellar.privilegeValue
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	return
}

func (self *ModInterStellar) SendInfo() {
	var msg S2C_InterStellarInfo
	msg.Cid = "interstellarinfo"
	msg.GalaxyInfo = self.Sql_InterStellar.galaxyInfo
	msg.StellarCount = self.Sql_InterStellar.StellarCount
	msg.StellarPos = GetTopInterstellarMgr().GetCurPos(self.player.Sql_UserBase.Uid)
	msg.PrivilegeValue = self.Sql_InterStellar.privilegeValue
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModInterStellar) SetStellarCount(value int) {
	self.Sql_InterStellar.StellarCount = value
	GetTopInterstellarMgr().SetCurProgress(self.player, self.Sql_InterStellar.StellarCount)
}

func (self *ModInterStellar) Check() {
	if self.Sql_InterStellar.galaxyInfo == nil {
		self.Sql_InterStellar.galaxyInfo = make(map[int]*Galaxy, 0)
	}

	for _, v := range GetCsvMgr().InterstellarConfig {
		_, okGalaxy := self.Sql_InterStellar.galaxyInfo[v.Galaxy]
		if !okGalaxy {
			self.Sql_InterStellar.galaxyInfo[v.Galaxy] = new(Galaxy)
			self.Sql_InterStellar.galaxyInfo[v.Galaxy].NebulaInfo = make(map[int]*Nebula)
			self.Sql_InterStellar.galaxyInfo[v.Galaxy].GalaxyId = v.Galaxy
		}

		_, okNebula := self.Sql_InterStellar.galaxyInfo[v.Galaxy].NebulaInfo[v.Id]
		if !okNebula {
			self.Sql_InterStellar.galaxyInfo[v.Galaxy].NebulaInfo[v.Id] = self.NewNebula(v.Id)
		} else {
			//检查任务
			for i := 0; i < len(self.Sql_InterStellar.galaxyInfo[v.Galaxy].NebulaInfo[v.Id].NebulaTask); i++ {
				taskNode := GetCsvMgr().InterstellarTaskNode[self.Sql_InterStellar.galaxyInfo[v.Galaxy].NebulaInfo[v.Id].NebulaTask[i].Taskid]
				if taskNode == nil {
					continue
				}
				if self.Sql_InterStellar.galaxyInfo[v.Galaxy].NebulaInfo[v.Id].NebulaTask[i].Tasktypes != taskNode.Tasktypes {
					self.Sql_InterStellar.galaxyInfo[v.Galaxy].NebulaInfo[v.Id].NebulaTask[i].Tasktypes = taskNode.Tasktypes
					self.Sql_InterStellar.galaxyInfo[v.Galaxy].NebulaInfo[v.Id].NebulaTask[i].Plan = 0
				}
			}

			//判断星点任务
			for _, war := range self.Sql_InterStellar.galaxyInfo[v.Galaxy].NebulaInfo[v.Id].NebulaWarInfo {
				if war.NebulaWarTask == nil {
					continue
				}

				taskNode := GetCsvMgr().InterstellarTaskNode[war.NebulaWarTask.Taskid]
				if taskNode == nil {
					continue
				}

				if war.NebulaWarTask.Tasktypes != taskNode.Tasktypes {
					war.NebulaWarTask.Tasktypes = taskNode.Tasktypes
					war.NebulaWarTask.Plan = 0
				}
			}
		}
	}
	count := 0
	for _, galaxy := range self.Sql_InterStellar.galaxyInfo {
		for _, nebula := range galaxy.NebulaInfo {
			for _, nebulaWar := range nebula.NebulaWarInfo {
				if nebulaWar.NebulaWarState == LOGIC_TRUE {
					count++
				}
			}
		}
	}
	self.Sql_InterStellar.StellarCount = count
	GetTopInterstellarMgr().SetCurProgress(self.player, count)
}

func (self *ModInterStellar) NewNebula(id int) *Nebula {
	data := new(Nebula)
	data.NebulaId = id

	config := GetCsvMgr().InterstellarConfig[data.NebulaId]
	if config == nil {
		return data
	}
	//初始化任务
	for i := 0; i < len(config.TaskTypes); i++ {
		if config.TaskTypes[i] == 0 {
			continue
		}
		task := new(JS_TaskInfo)
		task.Taskid = data.NebulaId*100 + i + 1
		task.Tasktypes = config.TaskTypes[i]
		data.NebulaTask = append(data.NebulaTask, task)
	}
	//初始化星点
	configWar, ok := GetCsvMgr().InterstellarWar[data.NebulaId]
	if !ok {
		return data
	}
	data.NebulaWarInfo = make(map[int]*NebulaWar)
	for _, v := range configWar {
		war := new(NebulaWar)
		war.NebulaWarId = v.Id
		//生成任务
		war.NebulaWarTask = new(JS_TaskInfo)
		war.NebulaWarTask.Taskid = war.NebulaWarId
		war.NebulaWarTask.Tasktypes = v.TaskTypes
		//生成奖励领取标记
		war.NebulaWarBoxSign = make(map[int]int)
		temp := GetCsvMgr().InterstellarBox
		configBox := temp[war.NebulaWarId]
		if configBox != nil {
			for _, box := range configBox {
				war.NebulaWarBoxSign[box.BoxId] = LOGIC_FALSE
			}
		}
		data.NebulaWarInfo[war.NebulaWarId] = war
	}
	return data
}

func (self *ModInterStellar) SendUpdate() {
	if len(self.chgNebula) > 0 {
		var msg S2C_InterStellUpdate
		msg.Cid = "updateinterstellarnebula"
		msg.Info = self.chgNebula
		self.chgNebula = make([]*JS_TaskInfo, 0)
		smsg, _ := json.Marshal(&msg)
		self.player.SendMsg(msg.Cid, smsg)
	}

	if len(self.chgNebulaWar) > 0 {
		var msg S2C_InterStellUpdate
		msg.Cid = "updateinterstellarnebulawar"
		msg.Info = self.chgNebulaWar
		self.chgNebulaWar = make([]*JS_TaskInfo, 0)
		smsg, _ := json.Marshal(&msg)
		self.player.SendMsg(msg.Cid, smsg)
	}
}

func (self *ModInterStellar) GetPrivilegeValue(nType int) int {
	return self.Sql_InterStellar.privilegeValue[nType]
}

func (self *ModInterStellar) GetPrivilegeValues() map[int]int {
	return self.Sql_InterStellar.privilegeValue
}

func (self *ModInterStellar) GMOrder(msg *C2S_GMInterstellar) {
	type C2S_GMInterstellar struct {
		NebulaId int `json:"nebulaid"` //星云id
		GroupId  int `json:"groupid"`  //地图id
		TaskId   int `json:"taskid"`   //任务id
	}

	for _, galaxy := range self.Sql_InterStellar.galaxyInfo {
		for _, nebula := range galaxy.NebulaInfo {
			if nebula.NebulaId != msg.NebulaId {
				continue
			}

			for _, v := range nebula.NebulaTask {
				if v.Taskid == msg.TaskId {
					v.Finish = LOGIC_TRUE
					self.chgNebula = append(self.chgNebula, v)
					return
				}
			}

			//判断星点任务
			for _, war := range nebula.NebulaWarInfo {
				if war.NebulaWarId != msg.GroupId {
					continue
				}
				if war.NebulaWarTask == nil {
					continue
				}
				if war.NebulaWarTask.Taskid == msg.TaskId {
					war.NebulaWarTask.Finish = LOGIC_TRUE
					self.chgNebulaWar = append(self.chgNebulaWar, war.NebulaWarTask)
					return
				}
			}
		}
	}
}

func (self *ModInterStellar) GMOrderAll() {
	for _, galaxy := range self.Sql_InterStellar.galaxyInfo {
		for _, nebula := range galaxy.NebulaInfo {
			configNebulaWarTemp := GetCsvMgr().InterstellarWar[nebula.NebulaId]
			if configNebulaWarTemp == nil {
				continue
			}
			for _, NebulaWar := range nebula.NebulaWarInfo {
				if NebulaWar.NebulaWarState == LOGIC_TRUE {
					continue
				}

				configNebulaWar := configNebulaWarTemp[NebulaWar.NebulaWarId]
				if configNebulaWar == nil {
					continue
				}

				NebulaWar.NebulaWarState = LOGIC_TRUE
				if self.Sql_InterStellar.privilegeValue == nil {
					self.Sql_InterStellar.privilegeValue = make(map[int]int, 0)
				}
				_, ok := self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType]

				switch configNebulaWar.Type {
				case VALUE_MAX:
					if ok {
						if configNebulaWar.PrivileGeNum > self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType] {
							self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType] = configNebulaWar.PrivileGeNum
						}
					} else {
						self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType] = configNebulaWar.PrivileGeNum
					}
					//增加激活奖励到挂机掉落中
					configHangup := GetCsvMgr().InterstellarHangup[configNebulaWar.PrivileGeNum]
					if configHangup != nil {
						self.player.GetModule("onhook").(*ModOnHook).AddPrivilegeExtItems(configHangup.AddItem, configHangup.AddNum)
					}
				case VALUE_PER_ADD, VALUE__ADD:
					self.Sql_InterStellar.privilegeValue[configNebulaWar.PrivileGeType] += configNebulaWar.PrivileGeNum
					if configNebulaWar.PrivileGeType == PRIVILEGE_ADD_SELFFIND_COUNT || configNebulaWar.PrivileGeType == PRIVILEGE_REDUCE_SELFFIND_OFFSET {
						self.player.GetModule("find").(*ModFind).CalSelfFind(self.Sql_InterStellar.privilegeValue, false)
					}
				}
			}
		}
	}
	self.Check()
	self.SendInfo()
}
