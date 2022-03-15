package game

import (
	"encoding/json"
	"fmt"
)

const (
	//MSG_RANK_TASK_TOP_INFO  = "msg_rank_task_top_info"  // 获得排行榜信息
	MSG_RANK_TASK_GET_STATE       = "msg_rank_task_get_state"       // 获得领取信息
	MSG_RANK_TASK_GET_PLAYER      = "msg_rank_task_get_player"      // 获得玩家
	MSG_RANK_TASK_GET_TYPE_PLAYER = "msg_rank_task_get_type_player" // 获得该类型所有第一名玩家
	MSG_RANK_TASK_GET_AWARD       = "msg_rank_task_get_award"       // 完成任务
	MSG_RANK_TASK_RED_POINT       = "msg_rank_task_red_point"       // 红点推送
)

type GetState struct {
	Taskid int `json:"taskid"` // 任务Id
	IsGet  int `json:"isget"`
}

// 服务器结构体
type San_RankTask struct {
	Uid      int64 // 角色ID
	Taskinfo string
	GetState string

	taskinfo []*JS_TaskInfo // 任务
	getState []*GetState

	DataUpdate
}

//! 心愿礼物
type ModRankTask struct {
	player       *Player
	San_RankTask San_RankTask
	chg          []JS_TaskInfo //更新的任务
	init         bool          // 是否初始化
}

func (self *ModRankTask) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_userranktask` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.San_RankTask, "san_userranktask", self.player.ID)

	if self.San_RankTask.Uid <= 0 {
		self.San_RankTask.Uid = self.player.ID

		self.Encode()
		InsertTable("san_userranktask", &self.San_RankTask, 0, true)
		self.San_RankTask.Init("san_userranktask", &self.San_RankTask, true)
	} else {
		self.Decode()
		self.San_RankTask.Init("san_userranktask", &self.San_RankTask, true)
	}
}

func (self *ModRankTask) OnGetOtherData() {

}

func (self *ModRankTask) onReg(handlers map[string]func(body []byte)) {
	//handlers[MSG_RANK_TASK_TOP_INFO] = self.SendTopInfo
	handlers[MSG_RANK_TASK_GET_STATE] = self.SendGetStateInfo
	handlers[MSG_RANK_TASK_GET_PLAYER] = self.SendGetPlayerInfo
	handlers[MSG_RANK_TASK_GET_TYPE_PLAYER] = self.SendGetTypePlayerInfo
	handlers[MSG_RANK_TASK_GET_AWARD] = self.GetAward // 领取奖励
}

//// 发送信息
//func (self *ModRankTask) SendTopInfo(body []byte) {
//	var msg C2S_RankTaskTopInfo
//	json.Unmarshal(body, &msg)
//
//	self.player.GetModule("top").(*ModTop).GetTop(msg.Index, msg.Ver)
//}

// 发送信息
func (self *ModRankTask) SendGetStateInfo(body []byte) {
	var msg C2S_RankTaskGetStateInfo
	json.Unmarshal(body, &msg)

	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_RANK_TASK)
	if !flag {
		return
	}

	var back S2C_RankTaskGetStateInfo
	back.Cid = MSG_RANK_TASK_GET_STATE
	back.GetState = self.San_RankTask.getState
	back.FinishState = GetRankTaskMgr().SetGetState()
	self.player.SendMsg(back.Cid, HF_JtoB(&back))
}

// 该类型首先完成玩家列表
func (self *ModRankTask) SendGetTypePlayerInfo(body []byte) {
	var msg C2S_RankTaskGetTypePlayerInfo
	json.Unmarshal(body, &msg)

	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_RANK_TASK)
	if !flag {
		return
	}

	var back S2C_RankTaskGetTypePlayerInfo
	back.Cid = MSG_RANK_TASK_GET_TYPE_PLAYER
	back.Type = msg.Type

	configs := GetCsvMgr().RankTaskConfig
	for _, config := range configs {
		if config.Type != msg.Type {
			continue
		}

		data := GetRankTaskMgr().GetRankTaskData(config.Id, false)
		if data == nil {
			continue
		}

		if len(data.rankPlayerInfo) <= 0 {
			continue
		}

		back.PlayerRank = append(back.PlayerRank, data.rankPlayerInfo[0])
		back.ID = append(back.ID, data.ID)
	}

	self.player.SendMsg(back.Cid, HF_JtoB(&back))
}

// 发送信息
func (self *ModRankTask) SendGetPlayerInfo(body []byte) {
	var msg C2S_RankTaskGetPlayerInfo
	json.Unmarshal(body, &msg)

	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_RANK_TASK)
	if !flag {
		//self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MINE_FAILING_TO_MEET_THE_OPENING"))
		return
	}

	// 获得配置
	config := GetCsvMgr().GetRankTaskConfigByID(msg.ID)
	if nil == config {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_CONFIGURATION_ERROR"))
		return
	}

	data := GetRankTaskMgr().GetRankTaskData(config.Id, false)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_WAR_ORDER_TASK_ONT_FINISH"))
		return
	}

	var back S2C_RankTaskGetPlayerInfo
	back.Cid = MSG_RANK_TASK_GET_PLAYER
	back.PlayerRank = data.rankPlayerInfo
	back.ID = msg.ID
	self.player.SendMsg(back.Cid, HF_JtoB(&back))
}

// 领取奖励
func (self *ModRankTask) GetAward(body []byte) {
	var msg C2S_RankTaskAward
	json.Unmarshal(body, &msg)

	flag := GetCsvMgr().IsLevelOpen2(self.player, OPEN_LEVEL_RANK_TASK)
	if !flag {
		//self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MINE_FAILING_TO_MEET_THE_OPENING"))
		return
	}

	// 获得配置
	config := GetCsvMgr().GetRankTaskConfigByID(msg.ID)
	if nil == config {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_CONFIGURATION_ERROR"))
		return
	}

	//flag, configLv := GetCsvMgr().IsLevelOpen(self.player.Sql_UserBase.Level, OPEN_LEVEL_35)
	//if !flag {
	//	self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_PASS_OPEN"), configLv))
	//	return
	//}

	data := GetRankTaskMgr().GetRankTaskData(config.Id, false)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_WAR_ORDER_TASK_ONT_FINISH"))
		return
	}

	for _, v := range self.San_RankTask.getState {
		if config.Id == v.Taskid {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_WAR_ORDER_IS_GET"))
			return
		}
	}

	self.San_RankTask.getState = append(self.San_RankTask.getState, &GetState{config.Id, 1})

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_RANK_GET_AWARD, config.Id, 0, 0, "领取排行榜奖励", 0, 0, self.player)

	var back S2C_RankTaskAward
	back.Cid = MSG_RANK_TASK_GET_AWARD
	back.ID = msg.ID
	back.Items = self.player.AddObjectLst(config.Items, config.Nums, "排行达标奖励", msg.ID, 0, 0)
	self.player.SendMsg(back.Cid, HF_JtoB(&back))
}

func (self *ModRankTask) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModRankTask) OnSave(sql bool) {
	self.Encode()
	self.San_RankTask.Update(sql)
}

func (self *ModRankTask) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.San_RankTask.Taskinfo), &self.San_RankTask.taskinfo)
	json.Unmarshal([]byte(self.San_RankTask.GetState), &self.San_RankTask.getState)
}

func (self *ModRankTask) Encode() { //! 将data数据写入数据库
	self.San_RankTask.Taskinfo = HF_JtoA(&self.San_RankTask.taskinfo)
	self.San_RankTask.GetState = HF_JtoA(&self.San_RankTask.getState)
}

// 任务处理
func (self *ModRankTask) HandleTask(tasktype, n2, n3, n4 int) {

	// 获得任务配置
	configs := GetCsvMgr().RankTaskConfig

	// 完成任务
	for _, config := range configs {
		if config.TaskTypes != tasktype {
			continue
		}

		data := GetRankTaskMgr().GetRankTaskData(config.Id, false)
		if data != nil {
			if len(data.rankPlayerInfo) >= RANK_TASK_PLAYER_LIST {
				continue
			}
		}

		node := self.GetTask(config.Id, true)

		if node != nil && node.Finish == 1 {
			continue
		}

		taskNode := TaskNode{config.Id, config.TaskTypes, config.Ns[0], config.Ns[1], config.Ns[2], config.Ns[3]}

		plan, add := DoTask(&taskNode, self.player, n2, n3, n4)
		if plan == 0 {
			continue
		}

		if node == nil {
			node = self.GetTask(config.Id, true)
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
			} else if plan > node.Plan {
				node.Plan = plan
				chg = true
			}
		}

		//! 任务完成
		if tasktype == PvpRankNow {
			if node.Plan != 0 && node.Plan <= config.Ns[0] {
				node.Finish = 1
				node.Plan = config.Ns[0]
				chg = true
				GetRankTaskMgr().GetRankTaskData(config.Id, true)
				GetRankTaskMgr().SetRankTask(config.Id, self.player)
			}
		} else if node.Plan >= config.Ns[0] {
			node.Finish = 1
			node.Plan = config.Ns[0]
			chg = true
			GetRankTaskMgr().GetRankTaskData(config.Id, true)
			GetRankTaskMgr().SetRankTask(config.Id, self.player)
		}

		if chg {
			self.chg = append(self.chg, *node)
		}
	}
}

//! 得到任务
func (self *ModRankTask) GetTask(taskid int, add bool) *JS_TaskInfo {
	for i := 0; i < len(self.San_RankTask.taskinfo); i++ {
		if self.San_RankTask.taskinfo[i].Taskid == taskid {
			return self.San_RankTask.taskinfo[i]
		}
	}

	if add {
		var node JS_TaskInfo
		node.Taskid = taskid
		node.Plan = 0
		node.Pickup = 0
		node.Finish = 0
		nIndex := -1
		for i, v := range GetCsvMgr().RankTaskConfig {
			if v.Id == taskid {
				nIndex = i
				break
			}
		}
		if nIndex < 0 {
			LogError("taskid not found, taskid:", taskid)
			return nil
		}
		node.Tasktypes = GetCsvMgr().RankTaskConfig[nIndex].TaskTypes
		self.San_RankTask.taskinfo = append(self.San_RankTask.taskinfo, &node)
		return self.San_RankTask.taskinfo[len(self.San_RankTask.taskinfo)-1]
	}

	return nil
}
