package game

import (
	"encoding/json"
	"fmt"
	"sync"

	//"time"
)

type HeroGrowInfo struct {
	ActivityType    int            `json:"activity_type"`    // 奖励ID
	ActivityPeriods int            `json:"activity_periods"` // 期数  N3
	StartTime       int64          `json:"starttime"`        //
	EndTime         int64          `json:"endtime"`          //
	TaskInfo        []*JS_TaskInfo `json:"taskinfo"`         // 任务类型
}

//! 任务数据库
type San_HeroGrowInfo struct {
	Uid           int64
	HeroGrowInfos string

	heroGrowInfos map[int]*HeroGrowInfo
	DataUpdate
}

//! 任务
type ModHeroGrow struct {
	player       *Player
	Sql_HeroGrow San_HeroGrowInfo

	chg    []JS_TaskInfo
	Locker *sync.RWMutex
}

func (self *ModHeroGrow) OnGetData(player *Player) {
	self.player = player
	self.Locker = new(sync.RWMutex)

	sql := fmt.Sprintf("select * from `san_userherogrow` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_HeroGrow, "san_userherogrow", self.player.ID)

	if self.Sql_HeroGrow.Uid <= 0 {
		self.Sql_HeroGrow.Uid = self.player.ID
		self.Sql_HeroGrow.heroGrowInfos = make(map[int]*HeroGrowInfo, 0)
		self.Encode()
		InsertTable("san_userherogrow", &self.Sql_HeroGrow, 0, true)
	} else {
		self.Decode()
	}
	if self.Sql_HeroGrow.heroGrowInfos == nil {
		self.Sql_HeroGrow.heroGrowInfos = make(map[int]*HeroGrowInfo, 0)
	}
	self.Sql_HeroGrow.Init("san_userherogrow", &self.Sql_HeroGrow, true)
}

//! 将数据库数据写入data
func (self *ModHeroGrow) Decode() {
	json.Unmarshal([]byte(self.Sql_HeroGrow.HeroGrowInfos), &self.Sql_HeroGrow.heroGrowInfos)
}

//! 将data数据写入数据库
func (self *ModHeroGrow) Encode() {
	self.Sql_HeroGrow.HeroGrowInfos = HF_JtoA(&self.Sql_HeroGrow.heroGrowInfos)
}

func (self *ModHeroGrow) OnGetOtherData() {

}

// 注册消息
func (self *ModHeroGrow) onReg(handlers map[string]func(body []byte)) {
	handlers["herogrowfreetask"] = self.HeroGrowFreeTask
}

func (self *ModHeroGrow) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModHeroGrow) OnSave(sql bool) {
	self.Encode()
	self.Sql_HeroGrow.Update(sql)
}

func (self *ModHeroGrow) SendInfo() {
	if self.Locker == nil {
		self.Locker = new(sync.RWMutex)
	}
	self.CheckTask()

	self.Locker.RLock()
	defer self.Locker.RUnlock()

	now := TimeServer().Unix()
	var msg S2C_HeroGrowInfo
	msg.Cid = "herogrowinfo"
	msg.HeroGrowInfo = make([]*HeroGrowInfo, 0)
	for _, v := range self.Sql_HeroGrow.heroGrowInfos {
		if now < v.StartTime || now > v.EndTime {
			continue
		}
		msg.HeroGrowInfo = append(msg.HeroGrowInfo, v)
	}
	msg.Config = GetCsvMgr().HeroGrowConfig
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

//根据活动进行任务检查
func (self *ModHeroGrow) CheckTask() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	now := TimeServer().Unix()
	for i := ACT_HERO_GROW_MIN; i < ACT_HERO_GROW_MAX; i++ {
		activity := GetActivityMgr().GetActivity(i)
		if activity == nil {
			continue
		}

		n3 := GetActivityMgr().getActN3(i)
		_, ok := self.Sql_HeroGrow.heroGrowInfos[i]
		if !ok {
			data := new(HeroGrowInfo)
			data.ActivityType = i
			data.ActivityPeriods = n3
			data.ResetTaskSafe()
			self.Sql_HeroGrow.heroGrowInfos[i] = data
		}

		//时间同步
		self.Sql_HeroGrow.heroGrowInfos[i].StartTime = HF_CalTimeForConfig(activity.info.Start, self.player.Sql_UserBase.Regtime)
		self.Sql_HeroGrow.heroGrowInfos[i].EndTime = self.Sql_HeroGrow.heroGrowInfos[i].StartTime + int64(activity.info.Continued) + int64(activity.info.Show)

		if now < self.Sql_HeroGrow.heroGrowInfos[i].StartTime || now > self.Sql_HeroGrow.heroGrowInfos[i].EndTime {
			continue
		}

		if self.Sql_HeroGrow.heroGrowInfos[i].ActivityPeriods != n3 || len(self.Sql_HeroGrow.heroGrowInfos[i].TaskInfo) == 0 {
			self.Sql_HeroGrow.heroGrowInfos[i].ActivityPeriods = n3
			self.Sql_HeroGrow.heroGrowInfos[i].ResetTaskSafe()
		}
	}
}

//任务初始化
func (self *HeroGrowInfo) ResetTaskSafe() {
	self.TaskInfo = make([]*JS_TaskInfo, 0)

	for _, v := range GetCsvMgr().HeroGrowConfig {
		if v.ActivityType == self.ActivityType && v.ActivityPeriods == self.ActivityPeriods {
			task := new(JS_TaskInfo)
			task.Taskid = v.Id
			self.TaskInfo = append(self.TaskInfo, task)
		}
	}
}

func (self *ModHeroGrow) HeroGrowFreeTask(body []byte) {

	var msg C2S_HeroGrowFreeTask
	json.Unmarshal(body, &msg)

	node := self.GetTask(msg.ActivityType, msg.Taskid)
	if node == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}
	//不等于0说明是付费的礼包
	if node.Tasktypes != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_MISSION_DOES_NOT_EXIST"))
		return
	}

	config := GetCsvMgr().GetHeroGrowConfig(node.Taskid)
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ARMY_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	nowStar := self.player.GetModule("hero").(*ModHero).GetHandBookStarMax(config.Hero)
	if nowStar < config.NeedQuality {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_ACTIVITY_THE_ACTIVITY_WAS_NOT_COMPLETED"))
		return
	}

	//免费礼包验证限制
	if node.Pickup >= config.Limit {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SHOP_REACH_THE_UPPER_LIMIT"))
		return
	}

	//发送任务奖励
	items := self.player.AddObjectLst(config.Item, config.Num, "领取免费英雄成长礼包", node.Taskid, self.player.Sql_UserBase.Vip, 0)
	node.Pickup++

	var msgRel S2C_HeroGrowFreeTask
	msgRel.Cid = "herogrowfreetask"
	msgRel.GetItems = items
	msgRel.TaskInfo = node
	msgRel.ActivityType = msg.ActivityType
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_GROW_BOX, node.Taskid, self.player.Sql_UserBase.Vip, 0, "领取免费英雄成长礼包", 0, 0, self.player)
}

func (self *ModHeroGrow) GetTask(activityType int, id int) *JS_TaskInfo {
	_, ok := self.Sql_HeroGrow.heroGrowInfos[activityType]
	if ok {
		for _, v := range self.Sql_HeroGrow.heroGrowInfos[activityType].TaskInfo {
			if v.Taskid == id {
				return v
			}
		}
	}
	return nil
}

func (self *ModHeroGrow) HandleRecharge(grade int, oldVip int) {

	task, activityType := self.GetTaskByGrade(grade)
	if task == nil {
		return
	}

	config := GetCsvMgr().GetHeroGrowConfig(task.Taskid)
	if config == nil {
		return
	}

	//发送任务奖励
	items := self.player.AddObjectLst(config.Item, config.Num, "购买英雄成长礼包", grade, self.player.Sql_UserBase.Vip, 0)
	task.Pickup++

	var msgRel S2C_HeroGrowFreeTask
	msgRel.Cid = "herogrowpaytask"
	msgRel.GetItems = items
	msgRel.TaskInfo = task
	msgRel.ActivityType = activityType
	self.player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HERO_GROW_RECHARGE, grade, oldVip, self.player.Sql_UserBase.Vip, "购买英雄成长礼包", 0, 0, self.player)
}

func (self *ModHeroGrow) GetTaskByGrade(grade int) (*JS_TaskInfo, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	now := TimeServer().Unix()
	for _, v := range self.Sql_HeroGrow.heroGrowInfos {
		if now < v.StartTime || now > v.EndTime {
			continue
		}
		for _, task := range v.TaskInfo {
			config := GetCsvMgr().GetHeroGrowConfig(task.Taskid)
			if config == nil || config.Ns[1] != grade {
				continue
			}
			return task, v.ActivityType
		}
	}
	return nil, 0
}
