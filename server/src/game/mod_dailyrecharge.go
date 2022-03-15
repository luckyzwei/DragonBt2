package game

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	ACT_DAILYRECHARGE = 9022
)

type ModDailyRecharge struct {
	player *Player
	Data   San_DailyRecharge      //! 充值数据
	Conf   []*DailyrechargeConfig //! 配置数据-乱序
	Index  int                    //! 完成的天数
}

type RechargeTask struct {
	TaskId   int   `json:"taskid"`   //! 任务Id
	TaskType int   `json:"tasktype"` //! 任务类型
	TaskTime int64 `json:"tasktime"` //! 任务时间
	Process  int   `json:"process"`  //! 进度
	Status   int   `json:"status"`   //! 0未完成 1完成 2领取
}

// 连续充值
type San_DailyRecharge struct {
	Uid  int64  //! 玩家uid
	Step int    //! 第几期
	Info string //! 任务状态

	info []*RechargeTask
	DataUpdate
}

func (self *ModDailyRecharge) Decode() {
	json.Unmarshal([]byte(self.Data.Info), &self.Data.info)
}

func (self *ModDailyRecharge) Encode() {
	self.Data.Info = HF_JtoA(self.Data.info)
}

func (self *ModDailyRecharge) OnGetData(player *Player) {
	self.player = player

	self.Conf = []*DailyrechargeConfig{}
}

//! 表名
func (self *ModDailyRecharge) TableName() string {
	return "san_userdailyrecharge"
}

// 每次重置
func (self *ModDailyRecharge) initInfo() {
	self.initTask()
}

func (self *ModDailyRecharge) OnGetOtherData() {
	tableName := self.TableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", tableName, self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Data, tableName, self.player.ID)
	if self.Data.Uid <= 0 {
		self.Data.Uid = self.player.ID
		self.initInfo()
		self.Encode()
		InsertTable(tableName, &self.Data, 0, true)
	} else {
		self.Decode()
	}
	self.Data.Init(tableName, &self.Data, true)

}

func (self *ModDailyRecharge) OnSave(sql bool) {
	self.Encode()
	self.Data.Update(sql)
}

func (self *ModDailyRecharge) OnRefresh() {

}

func (self *ModDailyRecharge) checkDailyInfo() {
	actN4 := GetActivityMgr().getActN4(self.getActType())
	if self.Data.Step != actN4 {
		self.Data.Step = actN4
		self.initInfo()
	} else {
		self.Conf = GetCsvMgr().getDailyConfig(self.Data.Step)
	}
}

// 消息处理
func (self *ModDailyRecharge) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "dailyrechargeaward": // 领取充值奖励
		self.checkDailyInfo()
		self.dailyRechargeAward(body, ctrl)
		return true
	case "actdailyrechargeinfo":
		self.SendInfo(false)
		return true
	}
	return false
}

func (self *ModDailyRecharge) dailyRechargeAward(body []byte, ctrl string) {
	if !self.isActOpen() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_DAILYRECHARGE_CONTINUOUS_RECHARGE_ACTIVITY_DOES_NOT"))
		return
	}

	var info C2S_DailyRechargeAward
	json.Unmarshal(body, &info)

	taskId := info.TaskId

	var task *RechargeTask
	for index := range self.Data.info {
		if self.Data.info[index].TaskId == taskId {
			task = self.Data.info[index]
			break
		}
	}

	if task == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSTART_CURRENT_ACTIVE_SUBITEM_DOES_NOT"))
		return
	}

	if task.Status == 2 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSTART_ACTIVITY_AWARDS_HAVE_BEEN_RECEIVED"))
		return
	}

	if task.Status != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_LUCKSTART_CURRENT_ACTIVITIES_NOT_COMPLETED"))
		return
	}

	pConfig := GetCsvMgr().GetDailyRechargeConfig(self.Data.Step, task.TaskId)
	if pConfig == nil {
		LogError("step:", self.Data.Step, ", taskId:", task.TaskId)
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	task.Status = 2
	var items []PassItem
	for index := range pConfig.Items {
		itemId := pConfig.Items[index]
		itemNum := pConfig.Nums[index]
		if itemId == 0 {
			continue
		}

		if itemNum == 0 {
			continue
		}

		outItem, outNum := self.player.AddObject(itemId, itemNum, 0, 0, 0, "连续充值")
		items = append(items, NewPassItem(outItem, outNum))
	}

	self.player.SendMsg(ctrl, HF_JtoB(&S2C_DailyRechargeAward{
		Cid:   ctrl,
		Task:  task,
		Items: items,
	}))
}

func (self *ModDailyRecharge) getActLog() int {
	return 0
}

func (self *ModDailyRecharge) isActOpen() bool {
	activity := GetActivityMgr().GetActivity(self.getActType())
	if activity == nil {
		//self.player.SendErrInfo("err", "活动不存在")
		return false
	}

	startTime := activity.getActTime()
	endTime := startTime + int64(activity.info.Continued)
	now := TimeServer().Unix()
	if now >= startTime && now <= endTime {
		return true
	}

	return false
}

func (self *ModDailyRecharge) getActType() int {
	return ACT_DAILYRECHARGE
}

func NewRechargeTask(taskId int, taskType int) *RechargeTask {
	return &RechargeTask{TaskId: taskId, TaskType: taskType}
}

//! 初始化任务信息
func (self *ModDailyRecharge) initTask() {
	self.Data.info = make([]*RechargeTask, 0)
	self.Index = 1
	self.Conf = GetCsvMgr().getDailyConfig(self.Data.Step)
	if len(self.Conf) <= 0 {
		//LogError("len(config) <= 0")
		return
	}

	for _, elem := range self.Conf {
		pTask := NewRechargeTask(elem.Id, elem.Type)
		self.Data.info = append(self.Data.info, pTask)
	}
}

//! 任务进度检查
func (self *ModDailyRecharge) taskProcessCheck() {
	now := TimeServer()
	timeSet := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local)
	var timeStamp int64
	if now.Unix() < timeSet.Unix() {
		timeStamp = timeSet.Unix()
	} else {
		timeStamp = timeSet.Unix() + DAY_SECS
	}

	index := 0
	for _, pTask := range self.Data.info {
		if pTask == nil {
			continue
		}

		index++
		// 检查时间
		if self.isTaskReset(pTask, timeStamp) { // 筛选不是今天的
			if pTask.Status <= 0 {
				pTask.Process = 0
				pTask.TaskTime = timeStamp
			}
		} else if pTask.TaskTime > 0 && pTask.TaskTime == timeStamp { // 今天的
			self.Index = index
			break
		}
	}
}

//! 完成任务
//! 依次完成第一天，第二天，第三天
//! 当天没有完成，第二天重置
//! 累积完成任务，触发累积任务
func (self *ModDailyRecharge) doTask(gem int) {
	if !self.isActOpen() {
		return
	}

	now := TimeServer()
	timeSet := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local)
	var timeStamp int64
	if now.Unix() < timeSet.Unix() {
		timeStamp = timeSet.Unix()
	} else {
		timeStamp = timeSet.Unix() + DAY_SECS
	}

	// 先检查1
	self.taskProcessCheck()

	for _, pTask := range self.Data.info {
		if pTask == nil {
			continue
		}

		// 检查时间
		if pTask.TaskTime > 0 && pTask.TaskTime == timeStamp { // 今天的
			if pTask.Status > 0 { // 完成的
				if pTask.Process < gem {
					pTask.Process = gem
				}
				break
			} else { // 重置
				if pTask.Process < gem {
					pTask.Process = gem
				}
			}

			self.checkTask(pTask)
			break
		}
	}

	// 检查任务1完成的个数
	taskFinishNum := 0
	for _, pTask := range self.Data.info {
		if pTask.Status != 0 && pTask.TaskType == 1 {
			taskFinishNum += 1
		}
	}

	// 完成累加任务
	if taskFinishNum > 0 {
		for i := len(self.Data.info) - 1; i >= 0; i-- {
			if self.Data.info[i].TaskType == 1 {
				break
			}
			conf := self.Conf[i]
			if conf.Diamond <= taskFinishNum && self.Data.info[i].Status == 0 {
				self.Data.info[i].Status = 1
				self.Data.info[i].TaskTime = timeStamp
				self.Data.info[i].Process = 1
			}
		}
	}

	// 同步任务状态
	self.sendTaskInfo()
}

func (self *ModDailyRecharge) isTaskReset(pTask *RechargeTask, timeStamp int64) bool {
	//!
	if pTask.TaskTime > 0 && pTask.TaskTime < timeStamp {
		return true
	} else if pTask.TaskTime == 0 {
		return true
	} else if pTask.TaskTime > 0 && pTask.TaskTime > timeStamp {
		return true
	}

	return false
}

func (self *ModDailyRecharge) checkTask(pTask *RechargeTask) {
	// 过滤配置
	config := GetCsvMgr().GetDailyRechargeConfig(self.Data.Step, pTask.TaskId)
	if config == nil {
		return
	}

	// 增加进度状态
	if pTask.Status == 0 && pTask.Process >= config.Diamond {
		// 设置状态
		pTask.Status = 1
	}
}

func (self *ModDailyRecharge) sendTaskInfo() {
	if !self.isActOpen() {
		self.sendMail()
		return
	}

	pInfo := &self.Data
	self.checkDailyInfo()
	cid := "actdailyrechargeinfo"
	todayNum, todayLeft := self.getToday()
	self.player.SendMsg(cid, HF_JtoB(&S2C_DailyRecharge{
		Cid:       cid,
		Step:      pInfo.Step,
		Info:      self.Data.info,
		TodayNum:  todayNum,
		TodayLeft: todayLeft,
		Index:     self.Index,
		Data:      []*DailyrechargeConfig{},
	}))
}

// 登录发送信息
func (self *ModDailyRecharge) SendInfo(conf bool) {
	if !self.isActOpen() {
		self.sendMail()
		return
	}

	pInfo := &self.Data
	self.checkDailyInfo()
	self.taskProcessCheck()
	cid := "actdailyrechargeinfo"
	todayNum, todayLeft := self.getToday()

	msg := S2C_DailyRecharge{
		Cid:       cid,
		Step:      pInfo.Step,
		Info:      self.Data.info,
		TodayNum:  todayNum,
		TodayLeft: todayLeft,
		Index:     self.Index,
	}

	if conf {
		msg.Data = self.Conf
	}

	self.player.SendMsg(cid, HF_JtoB(&msg))
}

// 从前往后检查, 如果taskTime都是0,则取第一个
// 否则取等于今天timeStamp的任务
func (self *ModDailyRecharge) getToday() (int, int) {
	isZero := true
	for index := range self.Data.info {
		pTask := self.Data.info[index]
		if pTask != nil && pTask.TaskTime != 0 {
			isZero = false
			break
		}
	}

	if isZero && len(self.Data.info) > 0 {
		target := self.Data.info[0]
		if target == nil {
			return 0, 0
		}

		pConfig := GetCsvMgr().GetDailyRechargeConfig(self.Data.Step, target.TaskId)
		if pConfig == nil {
			return 0, 0
		}
		num := target.Process
		left := pConfig.Diamond - num
		if left <= 0 {
			left = 0
		}
		return num, left
	}

	now := TimeServer()
	timeSet := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local)
	var timeStamp int64
	if now.Unix() < timeSet.Unix() {
		timeStamp = timeSet.Unix()
	} else {
		timeStamp = timeSet.Unix() + DAY_SECS
	}

	for index := range self.Data.info {
		pTask := self.Data.info[index]
		if pTask != nil && pTask.TaskTime == timeStamp {
			pConfig := GetCsvMgr().GetDailyRechargeConfig(self.Data.Step, pTask.TaskId)
			if pConfig == nil {
				return 0, 0
			}
			num := pTask.Process
			left := pConfig.Diamond - num
			if left <= 0 {
				left = 0
			}
			return num, left
		}
	}
	return 0, 0

}

func (self *ModDailyRecharge) sendMail() {
	var total []PassItem
	totalMap := make(map[int]int)
	for _, pTask := range self.Data.info {
		if pTask.Status == 1 {
			pConfig := GetCsvMgr().GetDailyRechargeConfig(self.Data.Step, pTask.TaskId)
			if pConfig == nil {
				continue
			}
			pTask.Status = 2
			for index := range pConfig.Items {
				itemId := pConfig.Items[index]
				itemNum := pConfig.Nums[index]
				if itemId == 0 {
					continue
				}

				if itemNum == 0 {
					continue
				}

				totalMap[itemId] += itemNum
			}
		}
	}

	if len(totalMap) <= 0 {
		return
	}

	for itemId, itemNum := range totalMap {
		total = append(total, NewPassItem(itemId, itemNum))
	}

	if len(total) <= 0 {
		return
	}

	pMail := self.player.GetModule("mail").(*ModMail)
	if pMail == nil {
		LogError("checkMail in SendRankAward, pMail == nil!")
		return
	}

	mailConfig, ok := GetCsvMgr().MailConfig[3004]
	if ok {
		pMail.AddMail(1, 1, 0, mailConfig.Mailtitle, mailConfig.Mailtxt, GetCsvMgr().GetText("STR_SYS"), total, false, 0)
	}
}
