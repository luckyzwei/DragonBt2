package game

import (
	"encoding/json"
	"fmt"
	"strings"
	//"time"
)

// 时光之巅
type San_Instance struct {
	Uid              int64  `json:"uid"` // 玩家Id
	InstanceInfo     string `json:"instanceinfo"`
	NowInstanceState string `json:"nowinstancestate"`
	BuffStore        string `json:"buffstore"`

	instanceInfo     map[int]*InstanceInfo // 时光之巅
	nowInstanceState *NowInstanceState     //
	buffStore        map[int]map[int]int   // 记录BUFF库存，辅助计算  第一个KEY是品质  第2个KEY是BUFFID
	DataUpdate
}

type ModInstance struct {
	player           *Player
	Sql_UserInstance San_Instance
	IsCheck          int
	buffStore        map[int]map[int]int // 记录BUFF库存，辅助计算  第一个KEY是品质  第2个KEY是BUFFID
}

type InstanceInfo struct {
	InstanceId  int                 `json:"id"`          //ID
	RewardState map[int]int         `json:"rewardstate"` //奖励状态
	Shadow      map[int]map[int]int `json:"shadow"`      //迷雾路径
}

type InstanceHeroState struct {
	HeroKeyId int `json:"herokeyid"` //
	Hp        int `json:"hp"`        //血量
	Energy    int `json:"energy"`    //能量
}

type NowInstanceState struct {
	NowInstanceId int                        `json:"nowinstanceid"` //当前地图ID
	NowRow        int                        `json:"nowrow"`        //当前位置ROW
	NowCol        int                        `json:"nowcol"`        //当前位置COL
	Buff          []int                      `json:"buff"`          //累计有哪些BUFF
	HeroState     map[int]*InstanceHeroState `json:"herostate"`     //保存英雄血量和能量
	GetHeroList   []*NewHero                 `json:"getherolist"`   //获得英雄
	BuffNew       map[int]int                `json:"buffnew"`       //辅助计算,之前的结构涉及客户端计算所以不能更改
	ThingInfo     map[int]*ThingInfo         `json:"thinginfo"`
	Level         int                        `json:"level"` // 用户法阵核心等级
}

type ThingInfo struct {
	ThingId     int                        `json:"thingid"`     //关卡ID
	ThingType   int                        `json:"thingtype"`   //
	ThingState  int                        `json:"thingstate"`  //
	IsRemove    int                        `json:"isremove"`    //
	Time        int64                      `json:"time"`        //
	HeroState   map[int]*InstanceHeroState `json:"herostate"`   //保存怪物血量和能量
	BuffChoose  []int                      `json:"buffchoose"`  //
	ThingHero   []*NewHero                 `json:"thinghero"`   //马车的英雄信息
	SwitchState int                        `json:"switchstate"` //开关状态
	TakeCount   int                        `json:"takecount"`   //触发计数，目前只针对类型20
}

const (
	INSTANCE_START          = 1  //起点
	INSTANCE_START_STORY    = 2  //初始剧情
	INSTANCE_MONSTER_LOW    = 3  //普通小怪
	INSTANCE_MONSTER_MID    = 4  //普通精锐
	INSTANCE_MONSTER_HIGH   = 5  //普通首领
	INSTANCE_BOX_LOW        = 6  //普通包箱
	INSTANCE_BOX_HIGH       = 7  //高级宝箱
	INSTANCE_BUFF           = 8  //选择BUFF
	INSTANCE_FRIEND         = 9  //获得佣兵
	INSTANCE_ADD            = 10 //队伍恢复
	INSTANCE_REBORN         = 11 //随机复活
	INSTANCE_STOP           = 12 //阻挡
	INSTANCE_MONSTER_LOW_S  = 13 //特殊小怪
	INSTANCE_MONSTER_MID_S  = 14 //特殊精锐
	INSTANCE_MONSTER_HIGH_S = 15 //特殊首领
	INSTANCE_STORY          = 16 //剧情/日志
	INSTANCE_SINGLE_DOOR    = 17 //单向传送门
	INSTANCE_SWITCH         = 20 //开关
)

func (self *ModInstance) OnGetData(player *Player) {
	self.player = player
	sql := fmt.Sprintf("select * from `san_userinstance` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_UserInstance, "san_userinstance", self.player.ID)

	if self.Sql_UserInstance.Uid <= 0 {
		self.Sql_UserInstance.Uid = self.player.ID
		self.Encode()
		InsertTable("san_userinstance", &self.Sql_UserInstance, 0, true)
		self.Sql_UserInstance.Init("san_userinstance", &self.Sql_UserInstance, true)
	} else {
		self.Decode()
		self.Sql_UserInstance.Init("san_userinstance", &self.Sql_UserInstance, true)
	}

}

func (self *ModInstance) OnSave(sql bool) {
	self.Encode()
	self.Sql_UserInstance.Update(sql)
}

func (self *ModInstance) OnGetOtherData() {

}

func (self *ModInstance) Decode() { // 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_UserInstance.InstanceInfo), &self.Sql_UserInstance.instanceInfo)
	json.Unmarshal([]byte(self.Sql_UserInstance.NowInstanceState), &self.Sql_UserInstance.nowInstanceState)
	json.Unmarshal([]byte(self.Sql_UserInstance.BuffStore), &self.Sql_UserInstance.buffStore)
}

func (self *ModInstance) Encode() { // 将data数据写入数据库
	self.Sql_UserInstance.InstanceInfo = HF_JtoA(&self.Sql_UserInstance.instanceInfo)
	self.Sql_UserInstance.NowInstanceState = HF_JtoA(&self.Sql_UserInstance.nowInstanceState)
	self.Sql_UserInstance.BuffStore = HF_JtoA(&self.Sql_UserInstance.buffStore)
}

func (self *ModInstance) Check() bool { // 将data数据写入数据库
	if self.IsCheck == LOGIC_TRUE {
		return false
	}
	if self.Sql_UserInstance.instanceInfo == nil {
		self.Sql_UserInstance.instanceInfo = make(map[int]*InstanceInfo, 0)
	}

	if self.Sql_UserInstance.nowInstanceState == nil {
		self.Sql_UserInstance.nowInstanceState = new(NowInstanceState)
	}

	if self.Sql_UserInstance.buffStore == nil {
		self.Sql_UserInstance.buffStore = make(map[int]map[int]int)
	}

	//先检查盒子旧数据
	for _, v := range self.Sql_UserInstance.instanceInfo {
		_, ok := GetCsvMgr().InstanceConfig[v.InstanceId]
		if ok {
			//对比箱子
			data := make(map[int]int, 0)
			for _, box := range GetCsvMgr().InstanceBox[v.InstanceId] {
				_, boxOK := v.RewardState[box.BoxId]
				if boxOK {
					data[box.BoxId] = v.RewardState[box.BoxId]
				} else {
					data[box.BoxId] = LOGIC_FALSE
				}
			}
			v.RewardState = data
		}
	}
	//插入盒子新数据
	for _, v := range GetCsvMgr().InstanceConfig {
		_, ok := self.Sql_UserInstance.instanceInfo[v.MapId]
		if !ok {
			data := new(InstanceInfo)
			data.InstanceId = v.MapId
			data.RewardState = make(map[int]int, 0)
			data.Shadow = make(map[int]map[int]int, 0)
			for _, box := range GetCsvMgr().InstanceBox[v.MapId] {
				data.RewardState[box.BoxId] = LOGIC_FALSE
			}
			self.Sql_UserInstance.instanceInfo[v.MapId] = data
		}
	}
	//检查当前地图是否有更新
	if self.Sql_UserInstance.nowInstanceState != nil {
		config, ok := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
		if ok {
			isNeedReset := false
			for _, v := range config {
				_, thingOK := self.Sql_UserInstance.nowInstanceState.ThingInfo[v.Id]
				if !thingOK {
					isNeedReset = true
					break
				}
			}
			if isNeedReset {
				self.Sql_UserInstance.nowInstanceState = self.MakeMap(LOGIC_FALSE)
			}
		}
	}

	self.IsCheck = LOGIC_TRUE
	return true
}

func (self *ModInstance) MakeMap(instanceId int) *NowInstanceState {
	mapRel := new(NowInstanceState)
	mapRel.NowInstanceId = instanceId
	mapRel.BuffNew = make(map[int]int)
	mapRel.ThingInfo = make(map[int]*ThingInfo)
	mapRel.Level = GetOfflineInfoMgr().GetMaxLevel(self.player.Sql_UserBase.Uid)
	info, ok := GetCsvMgr().InstanceThing[instanceId]
	if !ok {
		return mapRel
	}
	for _, v := range info {
		if v.Type == 0 {
			continue
		}
		thing := new(ThingInfo)
		thing.ThingId = v.Id
		thing.ThingType = v.Type
		thing.HeroState = make(map[int]*InstanceHeroState)
		if v.Type == INSTANCE_START {
			mapRel.NowRow = v.Row
			mapRel.NowCol = v.Col
			thing.ThingState = LOGIC_TRUE
			thing.Time = TimeServer().Unix()
		}
		_, ok := self.Sql_UserInstance.instanceInfo[mapRel.NowInstanceId].RewardState[v.Event]
		if ok {
			thing.ThingState = self.Sql_UserInstance.instanceInfo[mapRel.NowInstanceId].RewardState[v.Event]
		}

		switch thing.ThingType {
		case INSTANCE_FRIEND:
			//随便取个地牢的配置来生成，以后改
			for _, v := range GetCsvMgr().NewPitConfigMap {
				if v.Element != NEWPIT_TASK_SOUL_CART {
					continue
				}
				thing.ThingHero = self.player.GetModule("newpit").(*ModNewPit).MakeSoulCart(v, 0, 0)
			}
		case INSTANCE_SWITCH:
			thing.SwitchState = v.Event
		}
		mapRel.ThingInfo[thing.ThingId] = thing
	}

	self.CheckBuffShore()
	return mapRel
}

func (self *ModInstance) CheckBuffShore() {
	if self.buffStore == nil {
		self.buffStore = make(map[int]map[int]int)
		for _, v := range GetCsvMgr().NewPitRelique {
			_, ok := self.buffStore[v.Quality]
			if !ok {
				self.buffStore[v.Quality] = make(map[int]int)
			}
			self.buffStore[v.Quality][v.Id] = LOGIC_TRUE
		}
	}
}

func (self *ModInstance) SendInfo() {
	//if !GetCsvMgr().IsLevelAndPassOpenNew(self.player.Sql_UserBase.Uid, OPEN_LEVEL_NEWPITINFO) {
	//	return
	//}
	//更新过图相关的任务配置
	self.Check()
	var msg S2C_InstanceInfo
	msg.Cid = "instanceinfo"
	msg.InstanceInfo = self.Sql_UserInstance.instanceInfo
	msg.NowInstanceState = self.Sql_UserInstance.nowInstanceState
	self.player.SendMsg("instanceinfo", HF_JtoB(&msg))
}

func (self *ModInstance) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModInstance) onReg(handlers map[string]func(body []byte)) {
	handlers["instancestart"] = self.InstanceStart
	handlers["instanceend"] = self.InstanceEnd
	handlers["instancemove"] = self.InstanceMove
	handlers["instancebattle"] = self.InstanceBattle
	handlers["instancechoosebuff"] = self.InstanceChooseBuff
	handlers["instancemakebuff"] = self.InstanceMakeBuff
	handlers["finishthing"] = self.FinishThing
	handlers["instancefriend"] = self.InstanceFriend
	handlers["instanceadd"] = self.InstanceAdd
	handlers["instancereborn"] = self.InstanceReborn
	handlers["instancereswitch"] = self.InstanceSwitch
}

func (self *ModInstance) InstanceStart(body []byte) {
	var msg C2S_InstanceStart
	json.Unmarshal(body, &msg)

	if msg.Force == LOGIC_TRUE {
		self.Sql_UserInstance.nowInstanceState = new(NowInstanceState)
	}

	if self.Sql_UserInstance.nowInstanceState != nil {
		_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
		if ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("有其他关卡正在进行中!"))
			return
		}
	}

	_, ok := GetCsvMgr().InstanceConfig[msg.Id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到对应地图!"))
		return
	}

	self.Sql_UserInstance.nowInstanceState = self.MakeMap(msg.Id)

	var msgRel S2C_InstanceStart
	msgRel.Cid = "instancestart"
	msgRel.NowInstanceState = self.Sql_UserInstance.nowInstanceState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	low, high, rate := self.CalBox(self.Sql_UserInstance.nowInstanceState.NowInstanceId)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_INSTANCE_OPEN, self.Sql_UserInstance.nowInstanceState.NowInstanceId, low, high, "开启副本", 0, 0, self.player)

	self.player.HandleTask(TASK_TYPE_INSTANCE_PROCESS, self.Sql_UserInstance.nowInstanceState.NowInstanceId, rate, 0)
	return
}

func (self *ModInstance) InstanceEnd(body []byte) {
	low, high, _ := self.CalBox(self.Sql_UserInstance.nowInstanceState.NowInstanceId)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_INSTANCE_CLOSE, self.Sql_UserInstance.nowInstanceState.NowInstanceId, low, high, "关闭副本", 0, 0, self.player)
	self.Sql_UserInstance.nowInstanceState = new(NowInstanceState)
	self.player.SendRet2("instanceend")
	return
}

//处理移动逻辑和某些特殊的因移动而触发的关卡
func (self *ModInstance) InstanceMove(body []byte) {
	var msg C2S_InstanceMove
	json.Unmarshal(body, &msg)

	if self.Sql_UserInstance.nowInstanceState == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	info, okInfo := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !okInfo {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}
	self.Sql_UserInstance.nowInstanceState.NowRow = msg.Row
	self.Sql_UserInstance.nowInstanceState.NowCol = msg.Col

	//触发型移动
	if msg.ThingId > 0 {
		if self.Sql_UserInstance.nowInstanceState == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
			return
		}

		_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
			return
		}

		_, ok = GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
			return
		}

		instance, okInstance := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
		if !okInstance {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
			return
		}

		thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingId]
		if !okThing {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
			return
		}

		config := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thing.ThingId]
		if config == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
			return
		}

		if thing.ThingState == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("已经打过!"))
			return
		}

		var msgRel S2C_FinishThing
		msgRel.Cid = "finishthing"
		msgRel.Row = self.Sql_UserInstance.nowInstanceState.NowRow
		msgRel.Col = self.Sql_UserInstance.nowInstanceState.NowCol

		if info.Shadow == nil {
			info.Shadow = make(map[int]map[int]int, 0)
		}
		if len(msg.RowShadow) > 0 && len(msg.RowShadow) == len(msg.ColShadow) {
			for i := 0; i < len(msg.RowShadow); i++ {
				_, rowOK := info.Shadow[msg.RowShadow[i]]
				if !rowOK {
					info.Shadow[msg.RowShadow[i]] = make(map[int]int)
				}
				info.Shadow[msg.RowShadow[i]][msg.ColShadow[i]] = LOGIC_TRUE
				msgRel.RowShadow = append(msgRel.RowShadow, msg.RowShadow[i])
				msgRel.ColShadow = append(msgRel.ColShadow, msg.ColShadow[i])
			}
		}
		msgRel.Id = info.InstanceId
		thing.ThingState = LOGIC_TRUE
		thing.Time = TimeServer().Unix()

		//看是不是宝箱关卡
		switch config.Type {
		case INSTANCE_BOX_LOW, INSTANCE_BOX_HIGH:
			_, ok := instance.RewardState[config.Event]
			if !ok {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
				return
			}
			if instance.RewardState[config.Event] == LOGIC_TRUE {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("已经领过!"))
				return
			}
			rewardConfig := GetCsvMgr().GetInstanceBoxConfig(instance.InstanceId, config.Event)
			if rewardConfig == nil {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
				return
			}
			instance.RewardState[config.Event] = LOGIC_TRUE
			msgRel.Item = self.player.AddObjectLst(rewardConfig.Item, rewardConfig.Num, "获得宝箱奖励", config.Event, 0, 0)

			low, high, rate := self.CalBox(self.Sql_UserInstance.nowInstanceState.NowInstanceId)
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_INSTANCE_GET_REWARD, thing.ThingId, low, high, "获得宝箱奖励", 0, 0, self.player)
			self.player.HandleTask(TASK_TYPE_INSTANCE_PROCESS, instance.InstanceId, rate, 0)
			msgRel.InstanceInfo = instance
		case INSTANCE_SINGLE_DOOR:
			//门不设定 ,客户端需求
			thing.ThingState = LOGIC_FALSE
		}
		msgRel.ThingInfo = thing
		self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
		self.CheckThingAfter(config, instance)
	} else {
		//普通移动
		var msgRel S2C_InstanceMove
		msgRel.Cid = "instancemove"
		msgRel.Row = self.Sql_UserInstance.nowInstanceState.NowRow
		msgRel.Col = self.Sql_UserInstance.nowInstanceState.NowCol

		if info.Shadow == nil {
			info.Shadow = make(map[int]map[int]int, 0)
		}
		if len(msg.RowShadow) > 0 && len(msg.RowShadow) == len(msg.ColShadow) {
			for i := 0; i < len(msg.RowShadow); i++ {
				_, rowOK := info.Shadow[msg.RowShadow[i]]
				if !rowOK {
					info.Shadow[msg.RowShadow[i]] = make(map[int]int)
				}
				info.Shadow[msg.RowShadow[i]][msg.ColShadow[i]] = LOGIC_TRUE
				msgRel.RowShadow = append(msgRel.RowShadow, msg.RowShadow[i])
				msgRel.ColShadow = append(msgRel.ColShadow, msg.ColShadow[i])
			}
		}
		msgRel.Id = info.InstanceId
		self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	}
	return
}

func (self *ModInstance) FinishThing(body []byte) {
	var msg C2S_FinishThing
	json.Unmarshal(body, &msg)

	if self.Sql_UserInstance.nowInstanceState == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok = GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	instance, okInstance := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !okInstance {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingId]
	if !okThing {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	config := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thing.ThingId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	if thing.ThingState == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("已经打过!"))
		return
	}
	thing.ThingState = LOGIC_TRUE
	thing.Time = TimeServer().Unix()

	var msgRel S2C_FinishThing
	msgRel.Cid = "finishthing"
	//看是不是宝箱关卡
	switch config.Type {
	case INSTANCE_BOX_LOW, INSTANCE_BOX_HIGH:
		_, ok := instance.RewardState[config.Event]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
			return
		}
		if instance.RewardState[config.Event] == LOGIC_TRUE {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("已经领过!"))
			return
		}
		rewardConfig := GetCsvMgr().GetInstanceBoxConfig(instance.InstanceId, config.Event)
		if rewardConfig == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
			return
		}
		instance.RewardState[config.Event] = LOGIC_TRUE
		//self.CalBox(instance.InstanceId)
		msgRel.Item = self.player.AddObjectLst(rewardConfig.Item, rewardConfig.Num, "时光之巅宝箱", config.Event, 0, 0)
		msgRel.InstanceInfo = instance
	}
	msgRel.ThingInfo = thing
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	self.CheckThingAfter(config, instance)
	return
}

func (self *ModInstance) InstanceSwitch(body []byte) {
	var msg C2S_InstanceSwitch
	json.Unmarshal(body, &msg)

	if self.Sql_UserInstance.nowInstanceState == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok = GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	instance, okInstance := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !okInstance {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingId]
	if !okThing {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	config := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thing.ThingId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	//类型支持验证
	if config.Type != INSTANCE_SWITCH {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}
	thing.SwitchState = msg.SwitchState
	thing.TakeCount++

	var msgRel S2C_InstanceSwitch
	msgRel.Cid = "instancereswitch"
	msgRel.ThingInfo = thing

	//处理附加逻辑
	configList := make([]*InstanceThing, 0)
	if len(msg.ThingIdEx) > 0 && len(msg.ThingIdEx) == len(msg.SwitchStateEx) {
		for i := 0; i < len(msg.ThingIdEx); i++ {
			thingEx, okThingEx := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingIdEx[i]]
			if !okThingEx {
				continue
			}

			configEx := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thingEx.ThingId]
			if configEx == nil {
				continue
			}
			//类型支持验证
			if configEx.Type != INSTANCE_SWITCH {
				continue
			}
			thingEx.SwitchState = msg.SwitchStateEx[i]
			thingEx.TakeCount++
			msgRel.ThingInfoEx = append(msgRel.ThingInfoEx, thingEx)
			configList = append(configList, configEx)
		}
	}

	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	self.CheckThingAfter(config, instance)
	if len(msg.ThingIdEx) > 0 && len(msg.ThingIdEx) == len(msg.SwitchStateEx) {
		for i := 0; i < len(configList); i++ {
			self.CheckThingAfter(configList[i], instance)
		}
	}
	return
}

func (self *ModInstance) InstanceHeroState(heroKeyId int, hp int, energy int) *InstanceHeroState {
	info := new(InstanceHeroState)
	info.HeroKeyId = heroKeyId
	info.Hp = hp
	info.Energy = energy
	return info
}

func (self *ModInstance) InstanceBattle(body []byte) {
	var msg C2S_InstanceBattle
	json.Unmarshal(body, &msg)

	_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok = GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	instance, okInstance := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !okInstance {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingId]
	if !okThing {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	config := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thing.ThingId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	if thing.ThingState == LOGIC_TRUE || len(thing.BuffChoose) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("已经打过!"))
		return
	}

	//类型支持验证
	if config.Type != INSTANCE_MONSTER_LOW && config.Type != INSTANCE_MONSTER_MID && config.Type != INSTANCE_MONSTER_HIGH {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	//更新一下英雄血量
	if self.Sql_UserInstance.nowInstanceState.HeroState == nil {
		self.Sql_UserInstance.nowInstanceState.HeroState = make(map[int]*InstanceHeroState, 0)
	}
	for i := 0; i < len(msg.HeroState); i++ {
		self.Sql_UserInstance.nowInstanceState.HeroState[msg.HeroState[i].HeroKeyId] = self.InstanceHeroState(msg.HeroState[i].HeroKeyId, msg.HeroState[i].Hp, msg.HeroState[i].Energy)
	}

	var msgRel S2C_InstanceBattle
	msgRel.Cid = "instancebattle"

	if msg.IsFail == LOGIC_TRUE {
		if thing.HeroState == nil {
			thing.HeroState = make(map[int]*InstanceHeroState, 0)
		}
		//更新一下怪物英雄血量
		for i := 0; i < len(msg.MonsterHeroState); i++ {
			thing.HeroState[msg.MonsterHeroState[i].HeroKeyId] = self.InstanceHeroState(msg.MonsterHeroState[i].HeroKeyId, msg.MonsterHeroState[i].Hp, msg.MonsterHeroState[i].Energy)
		}
		msgRel.ThingInfo = thing
		msgRel.HeroState = self.Sql_UserInstance.nowInstanceState.HeroState
		self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_INSTANCE_BATTLE_LOSE, thing.ThingId, 0, 0, "时光之巅战斗失败", 0, int(msg.Fight), self.player)
		return
	}
	thing.ThingState = LOGIC_TRUE

	low, high, _ := self.CalBox(self.Sql_UserInstance.nowInstanceState.NowInstanceId)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_INSTANCE_BATTLE_WIN, thing.ThingId, low, high, "时光之巅战斗胜利", 0, int(msg.Fight), self.player)
	//生成BUFF 待逻辑补
	//thing.BuffChoose = self.MakeBuff(config.Type - 2)
	msgRel.ThingInfo = thing
	msgRel.HeroState = self.Sql_UserInstance.nowInstanceState.HeroState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	self.CheckThingAfter(config, instance)
	return
}

func (self *ModInstance) InstanceFriend(body []byte) {
	var msg C2S_InstanceFriend
	json.Unmarshal(body, &msg)

	_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok = GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	instance, okInstance := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !okInstance {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingId]
	if !okThing {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	config := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thing.ThingId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	if thing.ThingState == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("已经打过!"))
		return
	}

	//类型支持验证
	if config.Type != INSTANCE_FRIEND {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	if msg.Param <= 0 || msg.Param > len(thing.ThingHero) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	thing.ThingState = LOGIC_TRUE
	thing.Time = TimeServer().Unix()

	var msgRel S2C_InstanceFriend
	msgRel.Cid = "instancefriend"
	msgRel.ThingInfo = thing

	hero := thing.ThingHero[msg.Param-1]
	hero.HeroKeyId = self.player.GetModule("hero").(*ModHero).MaxKey()
	msgRel.GetHero = hero
	self.Sql_UserInstance.nowInstanceState.GetHeroList = append(self.Sql_UserInstance.nowInstanceState.GetHeroList, hero)

	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	self.CheckThingAfter(config, instance)
	return
}

func (self *ModInstance) InstanceAdd(body []byte) {
	var msg C2S_InstanceAdd
	json.Unmarshal(body, &msg)

	_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok = GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	instance, okInstance := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !okInstance {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingId]
	if !okThing {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	config := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thing.ThingId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	if thing.ThingState == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("已经打过!"))
		return
	}

	//类型支持验证
	if config.Type != INSTANCE_ADD {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	thing.ThingState = LOGIC_TRUE
	thing.Time = TimeServer().Unix()

	valueHp := 5000
	valueMp := 0
	for i := 0; i < len(self.Sql_UserInstance.nowInstanceState.Buff); i++ {
		if self.Sql_UserInstance.nowInstanceState.Buff[i] == NEWPIT_BUFF_SPRING_LOW {
			valueHp += 2500
			valueMp += 2500
		} else if self.Sql_UserInstance.nowInstanceState.Buff[i] == NEWPIT_BUFF_SPRING_HIGH {
			valueHp += 5000
			valueMp += 5000
		}
	}
	for _, v := range self.Sql_UserInstance.nowInstanceState.HeroState {
		if v.Hp > 0 {
			v.Hp += valueHp
			v.Energy += valueMp
		}
		if v.Hp > 10000 {
			v.Hp = 10000
		}
		if v.Energy > 10000 {
			v.Energy = 10000
		}
	}

	var msgRel S2C_InstanceAdd
	msgRel.Cid = "instanceadd"
	msgRel.ThingInfo = thing
	msgRel.HeroState = self.Sql_UserInstance.nowInstanceState.HeroState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	self.CheckThingAfter(config, instance)
	return
}

func (self *ModInstance) InstanceReborn(body []byte) {
	var msg C2S_InstanceReborn
	json.Unmarshal(body, &msg)

	_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok = GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	instance, okInstance := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !okInstance {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingId]
	if !okThing {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	config := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thing.ThingId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	if thing.ThingState == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("已经打过!"))
		return
	}

	//类型支持验证
	if config.Type != INSTANCE_REBORN {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}
	isCan := false
	for _, v := range self.Sql_UserInstance.nowInstanceState.HeroState {
		if v.Hp <= 0 {
			isCan = true
		}
	}
	if !isCan {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("没有英雄阵亡!"))
		return
	}

	thing.ThingState = LOGIC_TRUE
	thing.Time = TimeServer().Unix()

	keyId := 0
	min := 10000
	for _, v := range self.Sql_UserInstance.nowInstanceState.HeroState {
		if v.Hp <= min {
			min = v.Hp
			keyId = v.HeroKeyId
		}
	}
	if keyId > 0 {
		self.Sql_UserInstance.nowInstanceState.HeroState[keyId].Hp = 10000
		if self.Sql_UserInstance.nowInstanceState.HeroState[keyId].Energy < 5000 {
			self.Sql_UserInstance.nowInstanceState.HeroState[keyId].Energy = 5000
		}
	}

	var msgRel S2C_InstanceReborn
	msgRel.Cid = "instanceadd"
	msgRel.ThingInfo = thing
	msgRel.HeroState = self.Sql_UserInstance.nowInstanceState.HeroState
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	self.CheckThingAfter(config, instance)
	return
}

func (self *ModInstance) InstanceChooseBuff(body []byte) {
	var msg C2S_InstanceChooseBuff
	json.Unmarshal(body, &msg)

	_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok = GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	instance, okInstance := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !okInstance {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingId]
	if !okThing {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	if msg.Index <= 0 || msg.Index > len(thing.BuffChoose) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_NEWPIT_INDEX_ERROR"))
		return
	}

	config := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thing.ThingId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	if thing.ThingState == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("已经打过!"))
		return
	}

	//类型支持验证
	if config.Type != INSTANCE_MONSTER_LOW && config.Type != INSTANCE_MONSTER_MID && config.Type != INSTANCE_MONSTER_HIGH && config.Type != INSTANCE_BUFF {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	buffId := thing.BuffChoose[msg.Index-1]
	self.Sql_UserInstance.nowInstanceState.Buff = append(self.Sql_UserInstance.nowInstanceState.Buff, buffId)
	if self.Sql_UserInstance.nowInstanceState.BuffNew == nil {
		self.Sql_UserInstance.nowInstanceState.BuffNew = make(map[int]int, 0)
	}
	_, hasOk := self.Sql_UserInstance.nowInstanceState.BuffNew[buffId]
	if hasOk {
		self.Sql_UserInstance.nowInstanceState.BuffNew[buffId] += 1
	} else {
		self.Sql_UserInstance.nowInstanceState.BuffNew[buffId] = 1
	}
	thing.ThingState = LOGIC_TRUE
	thing.Time = TimeServer().Unix()
	var msgRel S2C_InstanceChooseBuff
	msgRel.Cid = "instancechoosebuff"
	msgRel.ThingInfo = thing
	msgRel.Buff = self.Sql_UserInstance.nowInstanceState.Buff
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	self.CheckThingAfter(config, instance)
	return
}

func (self *ModInstance) InstanceMakeBuff(body []byte) {
	var msg C2S_InstanceMakeBuff
	json.Unmarshal(body, &msg)

	_, ok := GetCsvMgr().InstanceConfig[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	_, ok = GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	instance, okInstance := self.Sql_UserInstance.instanceInfo[self.Sql_UserInstance.nowInstanceState.NowInstanceId]
	if !okInstance {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[msg.ThingId]
	if !okThing {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到关卡!"))
		return
	}

	config := GetCsvMgr().InstanceThing[self.Sql_UserInstance.nowInstanceState.NowInstanceId][thing.ThingId]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	if thing.ThingState == LOGIC_TRUE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("已经打过!"))
		return
	}

	//类型支持验证
	if config.Type != INSTANCE_BUFF {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("找不到配置!"))
		return
	}

	//类型支持验证
	if len(thing.BuffChoose) > 3 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("已经打过!"))
		return
	}

	level := config.Event
	if level >= 1 && level <= 3 {
		thing.BuffChoose = self.MakeBuff(level)
	} else {
		thing.BuffChoose = append(thing.BuffChoose, 1)
		thing.BuffChoose = append(thing.BuffChoose, 2)
		thing.BuffChoose = append(thing.BuffChoose, 3)
	}

	thing.Time = TimeServer().Unix()
	var msgRel S2C_InstanceMakeBuff
	msgRel.Cid = "instancemakebuff"
	msgRel.ThingInfo = thing
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
	self.CheckThingAfter(config, instance)
	return
}

func (self *ModInstance) MakeBuff(level int) []int {

	self.CheckBuffShore()

	buff := make([]int, 0)

	_, ok := self.buffStore[level]
	if !ok {
		buff = append(buff, 1)
		buff = append(buff, 2)
		buff = append(buff, 3)
		return buff
	}
	for times := 0; times < 3; times++ {
		for k, _ := range self.buffStore[level] {
			//验证身上存在这个BUFF时候不重复出现
			buffConfig := GetCsvMgr().NewPitRelique[k]
			if buffConfig == nil || buffConfig.Interesting == LOGIC_TRUE {
				continue
			}
			if buffConfig.OnlyOwned == LOGIC_TRUE {
				_, ok := self.Sql_UserInstance.nowInstanceState.BuffNew[k]
				if ok {
					continue
				}
			}
			//验证身上存在这个BUFF不能已经是3个之一
			isCan := true
			for j := 0; j < len(buff); j++ {
				if buff[j] == k {
					isCan = false
					break
				}
			}
			if !isCan {
				continue
			}

			buff = append(buff, k)
			break
		}
	}

	size := len(buff)
	for i := size; i < 3; i++ {
		for _, v := range self.buffStore {
			for kk, _ := range v {
				//验证身上存在这个BUFF时候不重复出现
				buffConfig := GetCsvMgr().NewPitRelique[kk]
				if buffConfig == nil || buffConfig.Interesting == LOGIC_TRUE {
					continue
				}
				if buffConfig.OnlyOwned == LOGIC_TRUE {
					_, ok := self.Sql_UserInstance.nowInstanceState.BuffNew[kk]
					if ok {
						continue
					}
				}
				//验证身上存在这个BUFF不能已经是3个之一
				isCan := true
				for j := 0; j < len(buff); j++ {
					if buff[j] == kk {
						isCan = false
						break
					}
				}
				if !isCan {
					continue
				}
				buff = append(buff, kk)
				break
			}
		}
	}
	if len(buff) <= 0 {
		buff = append(buff, 1)
		buff = append(buff, 2)
		buff = append(buff, 3)
	}
	return buff
}

func (self *ModInstance) CheckThingAfter(config *InstanceThing, instance *InstanceInfo) {
	if config == nil {
		return
	}

	group := strings.Split(config.Remove, "|")
	isNeedRemove := false
	if config.RemoveType == 0 {
		isNeedRemove = true
	} else {
		removeGroup := strings.Split(config.RemoveRelation, "|")
		count := 0
		for i := 0; i < len(removeGroup); i++ {
			index := HF_Atoi(removeGroup[i])
			if index == 0 {
				continue
			}
			thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[index]
			if !okThing {
				continue
			}
			if thing.ThingState == LOGIC_TRUE {
				count++
				if count >= config.RemoveType {
					isNeedRemove = true
					break
				}
			}
		}
	}
	if isNeedRemove {
		for i := 0; i < len(group); i++ {
			index := HF_Atoi(group[i])
			if index == 0 {
				continue
			}
			thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[index]
			if !okThing {
				continue
			}
			thing.IsRemove = LOGIC_TRUE
		}
	}

	establishGroup := strings.Split(config.Establish, "|")
	for i := 0; i < len(establishGroup); i++ {
		index := HF_Atoi(establishGroup[i])
		if index == 0 {
			continue
		}
		thing, okThing := self.Sql_UserInstance.nowInstanceState.ThingInfo[index]
		if !okThing {
			continue
		}
		switch thing.ThingType {
		case INSTANCE_BUFF:
			var msgRel S2C_InstanceUpdate
			msgRel.Cid = "instanceupdate"
			msgRel.ThingInfo = thing
			self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
		}
	}

	dispel := strings.Split(config.Dispel, "=")
	for i := 0; i < len(dispel); i++ {
		pointStr := dispel[i]
		if pointStr == "" {
			continue
		}
		point := strings.Split(pointStr, "|")
		if len(point) != 2 {
			continue
		}
		row := HF_Atoi(point[0])
		col := HF_Atoi(point[1])

		_, rowOK := instance.Shadow[row]
		if !rowOK {
			instance.Shadow[row] = make(map[int]int)
		}
		instance.Shadow[row][col] = LOGIC_TRUE
	}
}

func (self *ModInstance) GmShadowReset() {
	for _, v := range self.Sql_UserInstance.instanceInfo {
		v.Shadow = make(map[int]map[int]int, 0)
	}
	self.SendInfo()
}

func (self *ModInstance) GmGMInstancePass(msg *C2S_GMInstancePass) {

	if msg.Id == 0 {
		for _, instance := range self.Sql_UserInstance.instanceInfo {
			for k := range instance.RewardState {
				instance.RewardState[k] = LOGIC_TRUE
			}
			_, _, rate := self.CalBox(instance.InstanceId)
			self.player.HandleTask(TASK_TYPE_INSTANCE_PROCESS, instance.InstanceId, rate, 0)
		}
	} else {
		instance, ok := self.Sql_UserInstance.instanceInfo[msg.Id]
		if ok {
			for k := range instance.RewardState {
				instance.RewardState[k] = LOGIC_TRUE
			}
			_, _, rate := self.CalBox(instance.InstanceId)
			self.player.HandleTask(TASK_TYPE_INSTANCE_PROCESS, instance.InstanceId, rate, 0)
		}
	}

	self.SendInfo()
}

//计算宝箱进度
func (self *ModInstance) CalBox(instanceId int) (int, int, int) {
	lowCount := 0
	lowMax := 0
	highCount := 0
	highMax := 0
	rate := 0
	lowPer := 0
	highPer := 0
	info, okInfo := self.Sql_UserInstance.instanceInfo[instanceId]
	if okInfo {
		_, okBox := GetCsvMgr().InstanceBox[instanceId]
		if okBox {
			for _, box := range GetCsvMgr().InstanceBox[instanceId] {
				_, boxOK := info.RewardState[box.BoxId]
				if boxOK {
					if info.RewardState[box.BoxId] == LOGIC_TRUE {
						switch box.Type {
						case 1:
							lowCount++
							lowMax++
						case 2:
							highCount++
							highMax++
						}
					} else if info.RewardState[box.BoxId] == LOGIC_FALSE {
						switch box.Type {
						case 1:
							lowMax++
						case 2:
							highMax++
						}
					}
				}
			}
		}

		lowPer := 0
		highPer := 0
		if lowMax > 0 {
			lowPer = lowCount * 10000 / lowMax
		}

		if highMax > 0 {
			highPer = highCount * 10000 / highMax
		}
		if lowPer == 0 {
			rate = highPer
		} else if highPer == 0 {
			rate = lowPer
		} else {
			rate = (highPer + lowPer) / 2
		}
		if lowMax+highMax == 0 {
			return lowPer, highPer, 0
		}
		rate = (lowCount + highCount) * 100 / (lowMax + highMax)
	}

	return lowPer, highPer, rate
}

func (self *ModInstance) CheckTask() {
	for _, v := range self.Sql_UserInstance.instanceInfo {
		_, _, rate := self.CalBox(v.InstanceId)
		self.player.HandleTask(TASK_TYPE_INSTANCE_PROCESS, v.InstanceId, rate, 0)
	}
}
