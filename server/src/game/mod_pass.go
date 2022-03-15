package game

import (
	"encoding/json"
	"fmt"
	"sort"
	//"time"
)

// 关卡数据库
type San_UserPass struct {
	Uid         int64  `json:"uid"` // 玩家Id
	WarInfo     string `json:"warinfo"`
	MissionInfo string `json:"mission"`
	BoxInfo     string `json:"boxinfo"`
	StarBoxInfo string `json:"starboxinfo"`
	Losehero    int    `json:"losehero"` // 丢失英雄ID
	PassItem    string `json:"passitem"`
	PassInfo    string `json:"passinfo"`
	JJInfo      string `json:"jjinfo"`
	IsFight     int    `json:"isfight"`    // 国战开启标识
	Totalstars  int    `json:"totalstars"` // 总星级
	StartTime   int64  `json:"starttime"`  // 到达总星级时间

	missioninfo map[int]*JS_TaskInfo // 推图信息
	passitem    []PassItem           // 关卡奖励
	boxinfo     []int                // 宝箱信息
	starboxinfo []StarBoxInfo        // 关卡宝箱
	passinfo    []PassInfo           // 关卡信息
	jjinfo      []int                // 觐见信息

	DataUpdate
}

const (
	CAL_STAR_TYPE = 1
	EXPITEMID     = 91000005
)

type JS_War struct {
	Id   int    `json:"id"`   // 当前关卡索引
	Task [3]int `json:"task"` // 任务完成情况
}

type Mission struct {
	MissionId int `json:"missionid"` // 关卡Id
	Step      int `json:"step"`      // 步骤
	WarNum    int `json:"warnum"`    // 关卡次数
	WorkNum   int `json:"wroknum"`   // 其他次数
}

type PassInfo struct {
	Id   int `json:"id"`   // 关卡Id
	Star int `json:"star"` // 星级
	Num  int `json:"num"`  // 攻打次数
}

type StarBoxInfo struct {
	Id   int    `json:"id"`   // 宝箱Id
	Star int    `json:"star"` // 星级
	Get  [3]int `json:"get"`  // 领取状态
}

type DropGroupModify struct {
	Original int `json:"original"`
	New      int `json:"new"`
	Rate     int `json:"rate"`
}

type PassItem struct {
	ItemID int `json:"itemid"` // 道具ID
	Num    int `json:"num"`    // 道具数量
}

// 创建js passitem 对象
func NewPassItem(id, num int) PassItem {
	return PassItem{
		ItemID: id,
		Num:    num,
	}
}

// 道具排序
type lstPassItem []PassItem

func (s lstPassItem) Len() int      { return len(s) }
func (s lstPassItem) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByNum struct{ lstPassItem }

func (s ByNum) Less(i, j int) bool { return s.lstPassItem[i].Num > s.lstPassItem[j].Num }

// 关卡
type ModPass struct {
	player       *Player
	Sql_UserPass San_UserPass // 数据库结构
	pass         map[int]int  // 数据库结构
	chg          []JS_TaskInfo
}

func (self *ModPass) OnGetData(player *Player) {
	self.player = player
	sql := fmt.Sprintf("select * from `san_pass` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_UserPass, "san_pass", self.player.ID)

	if self.Sql_UserPass.Uid <= 0 {
		self.Sql_UserPass.Uid = self.player.ID
		self.Sql_UserPass.missioninfo = make(map[int]*JS_TaskInfo)
		self.Sql_UserPass.boxinfo = make([]int, 0)
		self.Sql_UserPass.starboxinfo = make([]StarBoxInfo, 0)
		self.Sql_UserPass.passitem = make([]PassItem, 0)
		self.Sql_UserPass.passinfo = make([]PassInfo, 0)
		self.Sql_UserPass.jjinfo = make([]int, 0)
		self.Sql_UserPass.Totalstars = 0
		self.Sql_UserPass.IsFight = 1
		self.Encode()
		InsertTable("san_pass", &self.Sql_UserPass, 0, true)
		self.Sql_UserPass.Init("san_pass", &self.Sql_UserPass, true)
	} else {
		self.Decode()
		self.Sql_UserPass.Init("san_pass", &self.Sql_UserPass, true)

		self.Sql_UserPass.IsFight = 1
		nTotalstars := 0
		self.pass = make(map[int]int)
		for i := 0; i < len(self.Sql_UserPass.passinfo); i++ {
			self.pass[self.Sql_UserPass.passinfo[i].Id] = LOGIC_TRUE
			if nTotalstars < self.Sql_UserPass.passinfo[i].Id {
				nTotalstars = self.Sql_UserPass.passinfo[i].Id
			}
		}

		if nTotalstars != self.Sql_UserPass.Totalstars {
			self.Sql_UserPass.Totalstars = nTotalstars
			self.Sql_UserPass.StartTime = TimeServer().Unix()
			//GetTopPassMgr().UpdateRank(nTotalstars, player)
		}
	}

	//更新过图相关的任务配置
	self.CheckTask()
}

func (self *ModPass) NewTaskInfo(taskid int, tasktypes int) *JS_TaskInfo {
	taskinfo := new(JS_TaskInfo)
	taskinfo.Taskid = taskid
	taskinfo.Tasktypes = tasktypes
	taskinfo.Plan = 0
	taskinfo.Finish = CANTFINISH

	return taskinfo
}

func (self *ModPass) CheckTask() {
	for _, v := range GetCsvMgr().LevelConfigMap {
		value, ok := self.Sql_UserPass.missioninfo[v.LevelId]
		if ok {
			if v.TaskTypes > 0 && value.Tasktypes != v.TaskTypes {
				value.Tasktypes = v.TaskTypes
				value.Plan = 0
			}
			continue
		}
		if v.TaskTypes > 0 {
			self.Sql_UserPass.missioninfo[v.LevelId] = self.NewTaskInfo(v.LevelId, v.TaskTypes)
		}
	}

	//CHECK部分任务
	if len(self.Sql_UserPass.missioninfo) > 0 {
		self.HandleTask(TASK_TYPE_PLAYER_LEVEL, 0, 0, 0)
		//self.HandleTask(PlayerFightTask, 0, 0, 0)
		self.HandleTask(HeroUpStarNumTask, 0, 0, 0)
		self.HandleTask(EquipUpNumTask, 0, 0, 0)
		self.HandleTask(OwnGemNumTask, 0, 0, 0)
		self.HandleTask(TigerOwnTask, 0, 0, 0)
	}
}

func (self *ModPass) HandleTask(taskType int, n1 int, n2 int, n3 int) {

	if len(self.Sql_UserPass.missioninfo) == 0 {
		self.CheckTask()
	}

	for _, pTask := range self.Sql_UserPass.missioninfo {
		if pTask == nil || pTask.Finish >= CANTAKE {
			continue
		}

		if pTask.Tasktypes != taskType {
			continue
		}

		config, ok := GetCsvMgr().LevelConfigMap[pTask.Taskid]
		if !ok {
			LogError("关卡任务配置错误")
			return
		}

		plan, add := DoTask(&TaskNode{Id: config.LevelId, Tasktypes: config.TaskTypes, N1: config.TaskConds[0], N2: config.TaskConds[1], N3: config.TaskConds[2], N4: config.TaskConds[3]}, self.player, n1, n2, n3)

		if plan == 0 {
			continue
		}

		if add {
			pTask.Plan += plan
		} else {
			if plan > pTask.Plan {
				pTask.Plan = plan
			}
		}

		if pTask.Plan >= config.TaskConds[0] {
			pTask.Plan = config.TaskConds[0]
			pTask.Finish = CANTAKE

			self.chg = append(self.chg, *pTask)
		}
	}
}

func (self *ModPass) SendUpdate() {
	if len(self.chg) == 0 {
		return
	}

	var msg S2C_UpdateMission
	msg.Cid = "updatemission"
	msg.TaskInfo = self.chg
	self.chg = make([]JS_TaskInfo, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("updatemission", smsg)
}

func (self *ModPass) OnGetOtherData() {
	nTotalstars := 0
	for i := 0; i < len(self.Sql_UserPass.passinfo); i++ {
		if self.Sql_UserPass.passinfo[i].Id > nTotalstars {
			nTotalstars = self.Sql_UserPass.passinfo[i].Id
		}
	}

	if nTotalstars != self.Sql_UserPass.Totalstars {
		self.Sql_UserPass.Totalstars = nTotalstars
		self.Sql_UserPass.StartTime = TimeServer().Unix()
		//GetTopPassMgr().UpdateRank(nTotalstars, self.player)
	}
}

// 设置关卡信息
func (self *ModPass) SetMission(missionId, step, warNum, workNum int) {

}

func (self *ModPass) GetMission(missionId int) {

	taskinfo := self.Sql_UserPass.missioninfo[missionId]
	if taskinfo == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_SMELT_MISSION_DOES_NOT_EXIST"))
		return
	}

	var msg S2C_GetMission
	msg.Cid = "getmission"
	msg.Taskid = taskinfo.Taskid
	msg.Tasktypes = taskinfo.Tasktypes
	msg.Plan = taskinfo.Plan
	msg.Finish = taskinfo.Finish
	self.player.SendMsg("getmission", HF_JtoB(&msg))
}

func (self *ModPass) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "passbegin": // 开始打关卡
		var msg C2S_BeginPass
		json.Unmarshal(body, &msg)
		self.PassBegin(msg.Passid, msg.Again, msg.AddFight)
		return true
	case "passskip": // 开始打关卡
		var msg C2S_PassSkip
		json.Unmarshal(body, &msg)
		self.PassSkip(msg.Passid, msg.UseArmy, msg.AddFight)
		return true
	case "passresult": // 关卡条件
		var msg C2S_EndPass
		json.Unmarshal(body, &msg)
		self.PassResult(msg.Passid, msg.Star, msg.Index, msg.FightTime, msg.Again, msg.UseArmy, msg.BattleInfo)
		if msg.Again == 0 {
			self.SetMission(msg.MissionId, msg.Step, msg.WarNum, msg.WorkNum)
		}
		return true
	case "passwin": // 谈判
		var msg C2S_SetMission
		json.Unmarshal(body, &msg)
		self.PassWin(msg.MissionId)
		self.SetMission(msg.MissionId, msg.Step, msg.WarNum, msg.WorkNum)
		return true
	case "passrecord":
		var msg C2S_PassRecord
		json.Unmarshal(body, &msg)
		GetPassRecordMgr().GetPassRecord(self.player, msg.Passid)
		return true
	case "passgrind":
		var c2s_msg C2S_SwapPass
		json.Unmarshal(body, &c2s_msg)
		self.PassSweep(c2s_msg.Passid, c2s_msg.Num)
		return true
	case "movecity":
		var msg C2S_MoveCity
		json.Unmarshal(body, &msg)
		self.player.Sql_UserBase.City = msg.Cityid
		return true
	case "setmission":
		var msg C2S_SetMission
		json.Unmarshal(body, &msg)
		self.SetMission(msg.MissionId, msg.Step, msg.WarNum, msg.WorkNum)
		return true
	case "boxstar":
		var msg C2S_BoxPass
		json.Unmarshal(body, &msg)
		self.StarBox(msg.Passid, msg.NoPass)
		return true
	case "jjpass":
		var msg C2S_JJPass
		json.Unmarshal(body, &msg)
		self.JJ(msg.Passid)
		return true
	case "getlosehero":
		self.GetLoseHero()
		return true
	case "gzopen":
		self.GzOpen()
		return true
	case "getmission":
		var msg C2S_GetMission
		json.Unmarshal(body, &msg)
		self.GetMission(msg.MissionId)
		return true
	}

	return false
}

// 处理丢失的英雄信息
func (self *ModPass) GetLoseHero() {
	if self.Sql_UserPass.Losehero == 0 {
		var lstret []PassItem
		lstret = append(lstret, PassItem{11400701, 1})
		self.player.AddObject(11400701, 1, 23, 0, 0, "关卡失败")
		self.Sql_UserPass.Losehero = 1
		var msg S2C_Getlosehero
		msg.Cid = "loseheroitem"
		msg.Ret = 0
		msg.Item = lstret
		smsg, _ := json.Marshal(&msg)
		self.player.SendMsg("loseheroitem", smsg)
	}
}

func (self *ModPass) OnSave(sql bool) {
	self.Encode()
	self.Sql_UserPass.Update(sql)
}

func (self *ModPass) OnRefresh() {
}

func (self *ModPass) Decode() { // 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_UserPass.MissionInfo), &self.Sql_UserPass.missioninfo)
	json.Unmarshal([]byte(self.Sql_UserPass.BoxInfo), &self.Sql_UserPass.boxinfo)
	json.Unmarshal([]byte(self.Sql_UserPass.StarBoxInfo), &self.Sql_UserPass.starboxinfo)
	json.Unmarshal([]byte(self.Sql_UserPass.PassItem), &self.Sql_UserPass.passitem)
	json.Unmarshal([]byte(self.Sql_UserPass.PassInfo), &self.Sql_UserPass.passinfo)
	json.Unmarshal([]byte(self.Sql_UserPass.JJInfo), &self.Sql_UserPass.jjinfo)
}

func (self *ModPass) Encode() { // 将data数据写入数据库
	self.Sql_UserPass.MissionInfo = HF_JtoA(&self.Sql_UserPass.missioninfo)
	self.Sql_UserPass.BoxInfo = HF_JtoA(&self.Sql_UserPass.boxinfo)
	self.Sql_UserPass.StarBoxInfo = HF_JtoA(&self.Sql_UserPass.starboxinfo)
	self.Sql_UserPass.PassItem = HF_JtoA(&self.Sql_UserPass.passitem)
	self.Sql_UserPass.PassInfo = HF_JtoA(&self.Sql_UserPass.passinfo)
	self.Sql_UserPass.JJInfo = HF_JtoA(&self.Sql_UserPass.jjinfo)
}

// 该关卡是否存在
func (self *ModPass) GetPass(passid int) *PassInfo {
	for i := 0; i < len(self.Sql_UserPass.passinfo); i++ {
		if self.Sql_UserPass.passinfo[i].Id == passid {
			return &self.Sql_UserPass.passinfo[i]
		}
	}

	return nil
}

// 获取关卡总星数
func (self *ModPass) GetRegionStar(regionid int) *StarBoxInfo {
	for i := 0; i < len(self.Sql_UserPass.starboxinfo); i++ {
		if self.Sql_UserPass.starboxinfo[i].Id == regionid {
			return &self.Sql_UserPass.starboxinfo[i]
		}
	}

	return nil
}

// 完成普通关卡总次数
func (self *ModPass) GetDoneNum() int {
	num := 0
	for i := 0; i < len(self.Sql_UserPass.passinfo); i++ {
		num += self.Sql_UserPass.passinfo[i].Num
	}

	return num
}

func (self *ModPass) FirstDrop(passid int) []PassItem {
	var outitem []PassItem
	csv_first := GetCsvMgr().Data["Level_Firstitem"][passid]
	for i := 0; i < 6; i++ {
		itemid := HF_Atoi(csv_first[fmt.Sprintf("firstitem%d", i+1)])
		if itemid == 0 {
			continue
		}
		num := HF_Atoi(csv_first[fmt.Sprintf("num%d", i+1)])
		outitem = append(outitem, PassItem{itemid, num})
	}
	return outitem
}

func (self *ModPass) CommonDrop(passid int) []PassItem {
	var outitem []PassItem
	csv_item := GetCsvMgr().Data["Level_Item"][passid]
	// 掉落
	for i := 0; i < 2; i++ {
		itemId := HF_Atoi(csv_item[fmt.Sprintf("item%d", i+1)])
		if itemId == 0 {
			continue
		}
		probability := HF_Atoi(csv_item[fmt.Sprintf("probability%d", i+1)])
		if probability == 10000 {
			outitem = append(outitem, PassItem{itemId, 1})
			continue
		}
		protection := HF_Atoi(csv_item[fmt.Sprintf("protection%d", i+1)])
		if protection != 0 { // 有保护
			outitem = append(outitem, PassItem{itemId, 1})
			continue
		}
		if HF_GetRandom(10000) < probability { // 概率掉落
			outitem = append(outitem, PassItem{itemId, 1})
			continue
		}
	}

	for i := 2; i < 6; i++ {
		itemId := HF_Atoi(csv_item[fmt.Sprintf("item%d", i+1)])
		if itemId == 0 {
			continue
		}
		num := HF_Atoi(csv_item[fmt.Sprintf("num%d", i+1)])
		if num == 0 {
			continue
		}

		probability := HF_Atoi(csv_item[fmt.Sprintf("probability%d", i+1)])
		if probability == 10000 {
			outitem = append(outitem, PassItem{itemId, num})
			continue
		}
		if HF_GetRandom(10000) < probability { // 概率掉落
			outitem = append(outitem, PassItem{itemId, num})
			continue
		}
	}

	// 掉落包
	for i := 0; i < 5; i++ {
		bag := HF_Atoi(csv_item[fmt.Sprintf("bag%d", i+1)])
		bagprobability := HF_Atoi(csv_item[fmt.Sprintf("bagprobability%d", i+1)])
		if HF_GetRandom(10000) >= bagprobability {
			continue
		}
		itemid, num := HF_DropForItemBag(bag)
		if itemid != 0 {
			outitem = append(outitem, PassItem{itemid, num})
		}
	}

	config, ok := GetCsvMgr().LevelItemMap[passid]
	if ok {
		lootItems := GetLootMgr().LootItems(config.Lotteryids, self.player)
		for _, v := range lootItems {
			outitem = append(outitem, PassItem{v.ItemId, v.ItemNum})
		}
	}
	return outitem
}

// 得到掉落
func (self *ModPass) GetDropItem(passid int, first bool) []PassItem {
	// 判断是否首次打
	if first { // 首次打
		return self.FirstDrop(passid)
	} else {
		return self.CommonDrop(passid)
	}

	return []PassItem{}
}

func (self *ModPass) GetPassType(passid int) int {
	return 1
}

func (self *ModPass) GetPassRegoinId(passid int) int {
	return 10000 + ((passid / 100) % 100)
}

func (self *ModPass) GetRegoinPassInfo(regionid int) (int, int) {
	begin := 110000 + (regionid-10000)*100
	end := begin + 99
	return begin, end
}

// 章节宝箱
func (self *ModPass) StarBox(chapterId int, box int) {
	csv, ok := GetCsvMgr().LevelMapConfig[chapterId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	if box < 1 || box > 3 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	regionid := chapterId
	starbox := self.GetRegionStar(regionid)
	if starbox == nil {
		var starbox1 StarBoxInfo
		starbox1.Id = regionid
		starbox1.Star = 0
		self.Sql_UserPass.starboxinfo = append(self.Sql_UserPass.starboxinfo, starbox1)
		starbox = &self.Sql_UserPass.starboxinfo[len(self.Sql_UserPass.starboxinfo)-1]

		passbegin, passend := self.GetRegoinPassInfo(regionid)
		//totalstars := 0
		for i := passbegin; i <= passend; {
			userpass := self.GetPass(i)
			csvLevelConfig := GetCsvMgr().LevelConfigMap[i]
			if userpass != nil && csvLevelConfig != nil && csvLevelConfig.LevelType != 2 {
				starbox.Star += userpass.Star
			}
			//i = i + 10
			i = i + 1
		}
		LogDebug("领取章节宝箱：", starbox.Id, starbox.Star, starbox.Get, box, passbegin, passend)
	} else {
		passbegin, passend := self.GetRegoinPassInfo(regionid)
		for i := passbegin; i <= passend; {
			userpass := self.GetPass(i)

			csvLevelConfig := GetCsvMgr().LevelConfigMap[i]
			if userpass != nil && csvLevelConfig != nil && csvLevelConfig.LevelType != 2 {
				starbox.Star += userpass.Star
			}
			i = i + 1
		}
		LogDebug("领取章节宝箱：", starbox.Id, starbox.Star, starbox.Get, box, passbegin, passend)
	}

	outitem := make([]PassItem, 0)
	needstar := csv.JNeedStar[box-1]
	if starbox.Star >= needstar && starbox.Get[box-1] == 0 {
		for i := 0; i < 3; i++ {
			itemid := csv.Prop[(box-1)*3+i]
			if itemid == 0 {
				continue
			}
			starbox.Get[box-1] = 1
			num := csv.PropNumber[(box-1)*3+i]
			itemid, num = self.player.AddObject(itemid, num, chapterId, 0, 0, "普通关卡星级宝箱领取")
			outitem = append(outitem, PassItem{itemid, num})
		}
	} else {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_CANT_GET_THE_TREASURE_BOX"))
		return
	}

	var msg S2C_WinPass
	msg.Cid = "starpass"
	msg.Passid = []int{chapterId}
	msg.OutItem = outitem
	msg.Box = box
	self.player.SendMsg("boxpass", HF_JtoB(&msg))
}

func (self *ModPass) JJ(passid int) {
	csv, ok := GetCsvMgr().LevelboxConfig[passid]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	for i := 0; i < len(self.Sql_UserPass.jjinfo); i++ {
		if self.Sql_UserPass.jjinfo[i] == passid {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_DONT_REPEAT_YOUR_OPINIONS"))
			return
		}
	}

	self.Sql_UserPass.jjinfo = append(self.Sql_UserPass.jjinfo, passid)

	outitem := make([]PassItem, 0)
	for i := 0; i < len(csv.Items); i++ {
		itemid := csv.Items[i]
		if itemid == 0 {
			continue
		}
		if len(csv.Items) != len(csv.Nums) {
			continue
		}
		num := csv.Nums[i]
		if num == 0 {
			continue
		}
		itemid, num = self.player.AddObject(itemid, num, 23, 0, 0, "关卡觐见")
		outitem = append(outitem, PassItem{itemid, num})
	}

	var msg S2C_JJPass
	msg.Cid = "jjpass"
	msg.Passid = passid
	msg.OutItem = outitem
	self.player.SendMsg("jjpass", HF_JtoB(&msg))
}

func (self *ModPass) GetPassRegionId(passid int) int {
	return 10000 + ((passid / 100) % 100)
}

func (self *ModPass) PassBegin(passid int, again int, addfight int64) {
	/*
		regionid := self.GetPassRegionId(passid)
		csvlevel, ok := GetCsvMgr().LevelMapConfig[regionid]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}

		if csvlevel.NeedLevel > self.player.Sql_UserBase.Level {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_INSUFFICIENT_GRADE"))
			return
		}
	*/
	ret, ok := GetCsvMgr().LevelConfigMap[passid]
	if !ok {
		return
	}

	preLevelIndex := self.player.GetModule("onhook").(*ModOnHook).GetStage()
	if preLevelIndex != passid {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_MISSION_CONDITIONS_NOT_ACHIEVED"))
		self.player.GetModule("onhook").(*ModOnHook).CheckPass()
		self.player.GetModule("onhook").(*ModOnHook).SendInfo(make([]byte, 0))
		return
	}

	if ret.Comat > self.player.GetTeamFight(TEAMTYPE_DEFAULT)+addfight {
		self.player.SendErrInfo("err", "战力不足")
		return
	}

	doubleitem, doubleexp := GetActivityMgr().GetDoubleStatus(DOUBLE_PASS)
	LogDebug("关卡掉落倍率:", doubleitem, doubleexp)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_PASS_START_NORMAL, passid, int(self.player.GetTeamFight(TEAMTYPE_DEFAULT)+addfight), 0, "挑战冒险关卡", 0, 0, self.player)
	if again == 1 {
		// 如果体力不足，返回
		csv_pass := GetCsvMgr().LevelConfigMap[passid]
		need := csv_pass.PhysicalStrength
		if self.player.GetPower() < need {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_LACK_OF_PHYSICAL_STRENGTH"))
			return
		}
		if self.GetPass(passid) == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_YOU_DIDNT_PASS_THE_BARRIER"))
			return
		}

		self.Sql_UserPass.passitem = self.GetDropItem(passid, false)
		if doubleitem > 1 {
			for i := 0; i < len(self.Sql_UserPass.passitem); i++ {
				self.Sql_UserPass.passitem[i].Num *= doubleitem
			}
		}
		var msg S2C_BeginPass
		msg.Cid = "passbegin"
		msg.Outitem = self.Sql_UserPass.passitem
		msg.Passid = passid
		msg.Tili = self.player.Sql_UserBase.TiLi
		self.player.SendMsg("passbegin", HF_JtoB(&msg))
	} else {
		//  20190702 第一次不消耗体力
		// 如果体力不足，返回
		/*
			csv_pass := GetCsvMgr().LevelConfigMap[passid]
			need := csv_pass.PhysicalStrength
			if self.player.GetPower() < need {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_LACK_OF_PHYSICAL_STRENGTH"))
				return
			}
		*/
		self.Sql_UserPass.passitem = self.GetDropItem(passid, true)
		var msg S2C_BeginPass
		msg.Cid = "passbegin"
		msg.Outitem = self.Sql_UserPass.passitem
		msg.Passid = passid
		msg.Tili = self.player.Sql_UserBase.TiLi
		self.player.SendMsg("passbegin", HF_JtoB(&msg))
	}
	self.player.HandleTask(TASK_TYPE_FINISH_PASS, passid, 0, 0)
}

func (self *ModPass) PassSkip(passid int, useArmy int, addfight int64) {
	config, ok := GetCsvMgr().LevelConfigMap[passid]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_GATE_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	preLevelIndex := self.player.GetModule("onhook").(*ModOnHook).GetStage()
	if preLevelIndex != passid {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_MISSION_CONDITIONS_NOT_ACHIEVED"))
		return
	}

	if config.MainType != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_GATE_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	if config.LevelSkip == LOGIC_FALSE || config.SkipType == LOGIC_FALSE {
		self.player.SendErrInfo("err", "此关卡无法跳过")
		return
	}

	if config.SkipNum > self.player.GetTeamFight(TEAMTYPE_DEFAULT)+addfight {
		self.player.SendErrInfo("err", "战力不足")
		return
	}

	doubleitem, doubleexp := GetActivityMgr().GetDoubleStatus(DOUBLE_PASS)
	LogDebug("关卡掉落倍率:", doubleitem, doubleexp)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_PASS_START_NORMAL, passid, int(self.player.GetTeamFight(TEAMTYPE_DEFAULT)+addfight), 0, "挑战冒险关卡", 0, 0, self.player)

	self.Sql_UserPass.passitem = self.GetDropItem(passid, true)

	if len(self.Sql_UserPass.passitem) == 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PASS_THERE_WAS_NO_DROP_AT"))
		return
	}
	star := 3

	if useArmy == LOGIC_TRUE {
		self.player.GetModule("friend").(*ModFriend).SetUseSign(HIRE_MOD_PASS, LOGIC_TRUE)
	}

	logType := 0
	if config.LevelGroup == 2 {
		logType = LOG_PASS_TYPE_SEARCH
	} else {
		logType = LOG_PASS_TYPE_FIGHT
	}

	var msg S2C_PassSkip
	msg.Cid = "passskip"
	msg.Passid = append(msg.Passid, passid)
	addStar := star
	passId := passid
	passnode := self.GetPass(passId)
	if passnode != nil {
		passnode.Num++
	} else {
		var node PassInfo
		node.Id = passId
		node.Num = 1
		node.Star = addStar
		self.Sql_UserPass.passinfo = append(self.Sql_UserPass.passinfo, node)
		//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_PASS_OK, passid, self.player.Sql_UserBase.Fight, self.player.Sql_UserBase.Level, "通关关卡", 0, fighttime, self.player)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_PASS_FINISH_NORMAL, passId, config.LevelType, 0, "跳过冒险关卡", 0, 0, self.player)
	}

	if passid > self.Sql_UserPass.Totalstars {
		self.Sql_UserPass.Totalstars = passid
		self.Sql_UserPass.StartTime = TimeServer().Unix()
	}

	// 加经验和武将经验
	teamexp := GetCsvMgr().GetData2Int("Level_Firstitem", passid, "teamexp") * doubleexp
	self.player.AddExp(teamexp, passid, logType, "冒险关卡通关")

	//  加道具
	msg.GetItems = self.player.AddObjectPassItem(self.Sql_UserPass.passitem, "通关冒险关卡", passid, 0, 0)
	self.Sql_UserPass.passitem = make([]PassItem, 0)
	//20210420  0表示开始 1表示结束...所以需要给2个任务
	self.player.HandleTask(TASK_TYPE_FINISH_PASS, passid, 0, 0)
	self.player.HandleTask(TASK_TYPE_FINISH_PASS, passid, 1, 0)

	// 胜利
	self.player.HandleTask(TASK_TYPE_FINISH_MAIN_PASS, config.TaskIndex, 0, 0)
	self.player.GetModule("support").(*ModSupportHero).CheckPassLevel(passid)
	self.player.GetModule("lifetree").(*ModLifeTree).CheckPass()

	//msg.Star = star

	//设置新的挂机关卡
	configNext := GetCsvMgr().LevelConfigMap[config.NextLevel]
	if configNext != nil && config.LevelIndex == configNext.LevelIndex {
		self.player.GetModule("onhook").(*ModOnHook).SetStage(configNext.LevelId)
		self.player.HandleTask(TASK_TYPE_FINISH_CHAPTER, config.LevelIndex-1, 0, 0)
	} else {
		self.player.GetModule("onhook").(*ModOnHook).SetStage(passid)
		self.player.HandleTask(TASK_TYPE_FINISH_CHAPTER, config.LevelIndex, 0, 0)
	}

	self.player.GetModule("task").(*ModTask).SendUpdate()

	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
	self.player.SendInfo("updateuserinfo")
}

// 战斗结果
func (self *ModPass) PassResult(passid int, star int, index int, fighttime int, again int, usearmy int, battleInfo *BattleInfo) {
	if len(self.Sql_UserPass.passitem) == 0 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PASS_THERE_WAS_NO_DROP_AT"))
		return
	}

	preLevelIndex := self.player.GetModule("onhook").(*ModOnHook).GetStage()
	if preLevelIndex != passid {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_MISSION_CONDITIONS_NOT_ACHIEVED"))
		return
	}

	config, ok := GetCsvMgr().LevelConfigMap[passid]
	if !ok && again == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_GATE_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	if config.MainType != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_GATE_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	//if config.Comat > self.player.GetTeamFight(TEAMTYPE_DEFAULT) {
	//	self.player.SendErrInfo("err", "战力不足")
	//	return
	//}
	if usearmy == LOGIC_TRUE {
		self.player.GetModule("friend").(*ModFriend).SetUseSign(HIRE_MOD_PASS, LOGIC_TRUE)
	}

	logType := 0
	if config.LevelGroup == 2 {
		logType = LOG_PASS_TYPE_SEARCH
	} else {
		logType = LOG_PASS_TYPE_FIGHT
	}

	doubleitem, doubleexp := GetActivityMgr().GetDoubleStatus(DOUBLE_PASS)
	LogDebug("关卡掉落倍率:", doubleitem, doubleexp)

	var msg S2C_EndPass
	msg.Cid = "passresult"
	msg.Passid = append(msg.Passid, passid)
	if star > 0 { // 胜利
		addStar := star
		for i := 0; i < 1; i++ {
			passId := passid
			if passId == 0 {
				continue
			}
			passnode := self.GetPass(passId)
			if passnode != nil {
				passnode.Num++
			} else {
				var node PassInfo
				node.Id = passId
				node.Num = 1
				node.Star = addStar
				self.Sql_UserPass.passinfo = append(self.Sql_UserPass.passinfo, node)
				//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_PASS_OK, passid, self.player.Sql_UserBase.Fight, self.player.Sql_UserBase.Level, "通关关卡", 0, fighttime, self.player)
				GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_PASS_FINISH_NORMAL, passId, config.LevelType, 0, "通关冒险关卡", 0, 0, self.player)
			}
		}
		//self.Sql_UserPass.Totalstars += addStar
		if passid > self.Sql_UserPass.Totalstars {
			self.Sql_UserPass.Totalstars = passid
			self.Sql_UserPass.StartTime = TimeServer().Unix()
			//GetTopPassMgr().UpdateRank(passid, self.player)
		}
		if addStar != 0 {
			if battleInfo != nil && battleInfo.UserInfo[0] != nil {

				team := self.player.getTeamPosByType(TEAMTYPE_DEFAULT)
				fight := int64(0)

				for i := 0; i < len(team.FightPos); i++ {
					heroKey := team.FightPos[i]
					if heroKey == 0 {
						continue
					}
					hero := self.player.getHero(heroKey)
					if nil == hero {
						continue
					}

					fight += hero.Fight
				}
				var army *ArmyInfo = nil
				for _, p := range battleInfo.UserInfo {
					for _, q := range p.HeroInfo {
						if q.ArmyInfo != nil && q.ArmyInfo.Uid != 0 {
							army = q.ArmyInfo
							army.Lv = q.HeroLv
							for i := 0; i < len(q.ArmyInfo.Atts); i++ {
								if q.ArmyInfo.Atts[i].Type == AttrFight {
									fight += q.ArmyInfo.Atts[i].Value
									break
								}
							}
							break
						}
					}
				}
				recordId := GetFightMgr().GetFightInfoID()
				data1 := San_PassRecord{}
				data1.KeyID = recordId
				data1.Uid = self.player.GetUid()
				data1.Name = self.player.GetName()
				data1.Icon = self.player.Sql_UserBase.IconId
				data1.Level = self.player.Sql_UserBase.Level
				data1.Portrait = self.player.Sql_UserBase.Portrait
				data1.BattleFight = fight
				data1.PlayerFight = self.player.Sql_UserBase.Fight
				data1.Time = TimeServer().Unix()

				battleInfo1 := BattleInfo(*battleInfo)
				battleInfo1.Id = recordId
				battleInfo1.LevelID = passid
				battleInfo1.Type = BATTLE_TYPE_PVE
				battleInfo1.Time = TimeServer().Unix()
				battleInfo1.UserInfo[0].Uid = self.player.GetUid()
				battleInfo1.UserInfo[0].Level = self.player.Sql_UserBase.Level
				battleInfo1.UserInfo[0].Icon = self.player.Sql_UserBase.IconId
				battleInfo1.UserInfo[0].Portrait = self.player.Sql_UserBase.Portrait
				battleInfo1.UserInfo[0].UnionName = self.player.GetUnionName()
				battleInfo1.UserInfo[0].Name = self.player.GetName()

				battleRecord1 := BattleRecord{}
				battleRecord1.Id = recordId
				battleRecord1.Time = TimeServer().Unix()
				battleRecord1.Result = battleInfo.Result
				battleRecord1.Type = BATTLE_TYPE_PVE
				battleRecord1.RandNum = battleInfo.Random
				battleRecord1.Weaken = battleInfo.Weaken
				battleRecord1.Side = 1
				battleRecord1.Level = passid
				battleRecord1.LevelID = passid

				if army != nil {
					battleRecord1.FightInfo[0] = GetRobotMgr().GetPlayerFightInfoWithArmyByPos(self.player, 0, 0, TEAMTYPE_DEFAULT, army)
					GetPassRecordMgr().CheckRecord(passid, addStar, self.player, &data1, &battleInfo1, &battleRecord1)
				} else {
					battleRecord1.FightInfo[0] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_DEFAULT)
					GetPassRecordMgr().CheckRecord(passid, addStar, self.player, &data1, &battleInfo1, &battleRecord1)
				}
			}
		}
		// 加经验和武将经验
		teamexp := GetCsvMgr().GetData2Int("Level_Firstitem", passid, "teamexp") * doubleexp

		self.player.AddExp(teamexp, passid, logType, "冒险关卡通关")

		//  加道具
		for i := 0; i < len(self.Sql_UserPass.passitem); i++ {
			self.Sql_UserPass.passitem[i].ItemID, self.Sql_UserPass.passitem[i].Num =
				self.player.AddObject(self.Sql_UserPass.passitem[i].ItemID, self.Sql_UserPass.passitem[i].Num, passid, config.LevelType, 0, "通关冒险关卡")
		}
		self.Sql_UserPass.passitem = make([]PassItem, 0)
		self.player.HandleTask(TASK_TYPE_FINISH_PASS, passid, 1, 0)
	} else { // 失败
		self.Sql_UserPass.passitem = make([]PassItem, 0)
		//self.player.HandleTask(TASK_TYPE_FINISH_PASS, passid, 0, 0)
	}

	// 胜利
	if star > 0 {
		self.player.HandleTask(TASK_TYPE_FINISH_MAIN_PASS, config.TaskIndex, 0, 0)
		self.player.GetModule("support").(*ModSupportHero).CheckPassLevel(passid)
		self.player.GetModule("lifetree").(*ModLifeTree).CheckPass()
	}

	msg.Star = star

	if msg.Star > 0 {
		//设置新的挂机关卡
		configNext := GetCsvMgr().LevelConfigMap[config.NextLevel]
		if configNext != nil && config.LevelIndex == configNext.LevelIndex {
			self.player.GetModule("onhook").(*ModOnHook).SetStage(configNext.LevelId)
			self.player.HandleTask(TASK_TYPE_FINISH_CHAPTER, config.LevelIndex-1, 0, 0)
		} else {
			self.player.GetModule("onhook").(*ModOnHook).SetStage(passid)
			self.player.HandleTask(TASK_TYPE_FINISH_CHAPTER, config.LevelIndex, 0, 0)
		}
	}
	self.player.GetModule("task").(*ModTask).SendUpdate()

	self.player.SendMsg("passresult", HF_JtoB(&msg))
	self.player.SendInfo("updateuserinfo")
}

func (self *ModPass) CheckHero(config *LevelConfig) bool {
	//新版 20190624  支持所有的任务类型
	if config.TaskTypes == 0 {
		return true
	}

	taskinfo := self.Sql_UserPass.missioninfo[config.LevelId]
	if taskinfo == nil {
		return true
	}

	if taskinfo.Finish >= CANTAKE {
		return true
	}

	return false

	/*
		if config.TaskTypes == 0 {
			return true
		}

		if len(config.TaskConds) != 4 {
			return true
		}
		number := config.TaskConds[0]
		heroId := config.TaskConds[1]
		star := config.TaskConds[3]

		if config.TaskTypes == 7 {
			heroConfig := GetCsvMgr().GetHeroConfig(heroId)
			if heroConfig == nil {
				return false
			}

			if heroConfig.HeroType == 2 {
				if self.player.hasBoss(heroId) {
					return true
				}
			}

			if heroId != 0 {
				hero := self.player.getHero(heroId)
				if hero == nil {
					return false
				}

				if hero.StarItem != nil && hero.StarItem.UpStar >= star {
					return true
				}
			} else {
				heros := self.player.getHeroes()
				total := 0
				for _, hero := range heros {
					if hero == nil {
						continue
					}

					if hero.StarItem != nil && hero.StarItem.UpStar >= star {
						total += 1
					}
				}
				if total >= number {
					return true
				}
			}
		} else if config.TaskTypes == 14 {
			if self.player.Sql_UserBase.Fight >= int64(number) {
				return true
			}
		}

		return false
	*/
}

// 直接通关
func (self *ModPass) PassWin(missionId int) {
	// 检查关卡的前置关卡是否完成
	config, ok := GetCsvMgr().LevelConfigMap[missionId]
	if !ok {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PASS_GATE_CONFIGURATION_ERROR"))
		return
	}

	if !self.CheckHero(config) {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PASS_MISSION_CONDITIONS_NOT_ACHIEVED"))
		return
	}

	// 判断前置关卡是否完成
	for i := 0; i < len(config.ALevel); i++ {
		if config.ALevel[i] == 0 {
			continue
		}

		pass := self.GetPass(config.ALevel[i])
		if pass == nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PASS_PRE-CHECKPOINT_NOT_COMPLETED"))
			return
		}
	}

	// 检查关卡谈判条件是否完成

	costitem := make(map[int]int)
	outitem := make([]PassItem, 0)
	var msg S2C_WinPass
	msg.Cid = "passwin"
	for i := 0; i < 1; i++ {
		passid := missionId
		msg.Passid = append(msg.Passid, passid)
		passnode := self.GetPass(passid)
		if passnode != nil {
			passnode.Num++
		} else {
			var node PassInfo
			node.Id = passid
			node.Num = 1

			if config.LevelType != CAL_STAR_TYPE {
				node.Star = 0
			} else {
				node.Star = 3
			}

			self.Sql_UserPass.passinfo = append(self.Sql_UserPass.passinfo, node)

			//通关行为log
			//curteamfight := HF_GetCurTeamFightNum(self.player.Sql_UserBase.Uid)
			//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_PASS_OK, passid, curteamfight, self.player.Sql_UserBase.Level, "直通关卡", 0, self.player.Sql_UserBase.Vip, self.player)
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_PASS_FINISH_NORMAL, passid, config.LevelType, 0, "冒险关卡通关", 0, 0, self.player)
		}

		// 加经验和武将经验
		teamexp := 0
		if passnode == nil {
			teamexp = GetCsvMgr().GetData2Int("Level_Firstitem", passid, "teamexp")
			outitem = self.GetDropItem(passid, true)
		} else {
			teamexp = config.TeamExp
			outitem = self.GetDropItem(passid, false)
		}

		self.player.AddExp(teamexp, passid, LOG_PASS_TYPE_PASS_WIN, "冒险关卡通关")
		for i := 0; i < len(outitem); i++ {
			outitem[i].ItemID, outitem[i].Num = self.player.AddObject(outitem[i].ItemID, outitem[i].Num, passid, LOG_PASS_TYPE_PASS_WIN, 0, "冒险关卡通关")
		}
	}

	// 扣除贿赂
	for key, value := range costitem {
		self.player.AddObject(key, -value, missionId, LOG_PASS_TYPE_PASS_WIN, 0, "普通关卡通关")
		outitem = append(outitem, PassItem{key, -value})
	}

	self.player.HandleTask(CommonLevelTask, config.LevelId, 1, 1)
	msg.OutItem = outitem
	self.player.SendMsg("passwin", HF_JtoB(&msg))

	self.player.SendInfo("updateuserinfo")
}

// 国战开启
func (self *ModPass) GzOpen() {
	csv, ok := GetCsvMgr().Data["Nationalwar_Parm"][3]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	if self.player.Sql_UserBase.Level < HF_Atoi(csv["parm1"]) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_INSUFFICIENT_GRADE"))
		return
	}

	self.Sql_UserPass.IsFight = 1

	out := make([]PassItem, 0)
	for i := 2; i <= 6; i += 2 {
		itemid := HF_Atoi(csv[fmt.Sprintf("parm%d", i)])
		if itemid == 0 {
			break
		}
		itemnum := HF_Atoi(csv[fmt.Sprintf("parm%d", i+1)])
		itemid, itemnum = self.player.AddObject(itemid, itemnum, 15, 0, 0, "阵营战开启")
		out = append(out, PassItem{itemid, itemnum})
	}

	var msg S2C_GZOpen
	msg.Cid = "gzopen"
	msg.Item = out
	self.player.SendMsg("gzopen", HF_JtoB(&msg))
}

// 关卡扫荡 1:体力不足 2：不是3星关卡 3：vip等级不够 4:次数不足
func (self *ModPass) PassSweep(passid int, num int) {
	if num <= 0 {
		self.player.SendRet("passgrind", -1)
		return
	}

	doubleitem, doubleexp := GetActivityMgr().GetDoubleStatus(DOUBLE_PASS)
	LogDebug("关卡掉落倍率:", doubleitem, doubleexp)

	passnode := self.GetPass(passid)
	if passnode == nil {
		lastpass := self.Sql_UserPass.passinfo[len(self.Sql_UserPass.passinfo)-1]
		if lastpass.Id > passid {
			var newpass PassInfo
			newpass.Id = passid
			newpass.Num = 1
			self.Sql_UserPass.passinfo = append(self.Sql_UserPass.passinfo, newpass)
		} else {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_NO_BARRIER_EXISTS"))
			return
		}

	} else {
		passnode.Num += num
	}

	csv_pass := GetCsvMgr().LevelConfigMap[passid]
	need := csv_pass.PhysicalStrength * num
	if self.player.GetPower() < need {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_LACK_OF_PHYSICAL_STRENGTH"))
		return
	}

	self.player.AddPower(-need, passid, LOG_PASS_TYPE_SWEEP, "冒险关卡通关")

	mapitem := make(map[int]int)
	var outitem []PassItem
	// 获得物品
	for i := 0; i < num; i++ {
		item := self.GetDropItem(passid, false)
		for j := 0; j < len(item); j++ {
			mapitem[item[j].ItemID] += item[j].Num * doubleitem
		}
	}

	//经验改成物品发送
	teamexp := csv_pass.TeamExp
	mapitem[EXPITEMID] = teamexp * num * doubleexp

	for key, value := range mapitem {
		key, value = self.player.AddObject(key, value, passid, LOG_PASS_TYPE_SWEEP, 0, "冒险关卡通关")
		outitem = append(outitem, PassItem{key, value})
	}

	sort.Sort(ByNum{outitem})
	//teamexp := csv_pass.TeamExp
	//self.player.AddExp(teamexp*num*doubleexp, passid, LOG_PASS_TYPE_SWEEP, "普通关卡通关")

	self.player.HandleTask(CommonLevelTask, passid, num, 1)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_PASS_FINISH_NORMAL, passid, csv_pass.LevelType, 0, "冒险关卡通关", 0, 0, self.player)

	var msg S2C_SwapPass
	msg.Uid = self.player.ID
	msg.Cid = "passgrind"
	msg.Passid = passid
	msg.Num = num
	msg.Outitem = outitem
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("passgrind", smsg)

	self.player.SendInfo("updateuserinfo")
}

////////////////////////////////
func (self *ModPass) SendInfo() {
	self.Encode()

	var msg S2C_PassInfo
	msg.Cid = "passinfo"
	msg.Passinfo.WarInfo = self.Sql_UserPass.WarInfo
	msg.Passinfo.PassInfo = self.Sql_UserPass.PassInfo
	msg.Passinfo.MissionInfo = self.Sql_UserPass.MissionInfo
	msg.Passinfo.BoxInfo = self.Sql_UserPass.BoxInfo
	msg.Passinfo.StarBoxInfo = self.Sql_UserPass.StarBoxInfo

	msg.Passinfo.JJInfo = self.Sql_UserPass.JJInfo
	msg.IsFight = self.Sql_UserPass.IsFight
	self.player.SendMsg("passinfo", HF_JtoB(&msg))
}

func (self *ModPass) GetPassByType(passid int, passType int) *PassInfo {
	for i := 0; i < len(self.Sql_UserPass.passinfo); i++ {
		if self.Sql_UserPass.passinfo[i].Id == passid {
			config, ok := GetCsvMgr().LevelConfigMap[passid]
			if !ok {
				continue
			}
			if config.MainType == passType {
				return &self.Sql_UserPass.passinfo[i]
			}
		}
	}

	return nil
}

// gm工具
func (self *ModPass) GMPassChapter(chapter int) {
	if chapter > 100000 {
		csvLevelConfig := GetCsvMgr().LevelConfigMap[chapter]
		if csvLevelConfig == nil {
			return
		}
		for i := 1; i <= csvLevelConfig.LevelIndex; i++ {
			config, ok := GetCsvMgr().MissionMap[i]
			if !ok {
				return
			}
			for _, missionId := range config.MissionIds {
				if missionId >= chapter {
					continue
				}
				csvLevelConfig := GetCsvMgr().LevelConfigMap[missionId]
				if csvLevelConfig == nil {
					return
				}
				addStar := 3
				if csvLevelConfig.LevelType != CAL_STAR_TYPE {
					addStar = 0
				}

				mission := self.GetPass(missionId)
				if mission == nil {
					var node PassInfo
					node.Id = missionId
					node.Num = 1
					node.Star = addStar
					self.Sql_UserPass.passinfo = append(self.Sql_UserPass.passinfo, node)
				} else {
					mission.Num += 1
					mission.Star = addStar
				}
				self.player.HandleTask(TASK_TYPE_FINISH_PASS, missionId, 1, 0)

				if csvLevelConfig.MainType == 1 {
					self.player.HandleTask(TASK_TYPE_FINISH_MAIN_PASS, csvLevelConfig.TaskIndex, 0, 0)
				}

				//完成章节任务
				configNext := GetCsvMgr().LevelConfigMap[csvLevelConfig.NextLevel]
				if configNext == nil || csvLevelConfig.LevelIndex != configNext.LevelIndex {
					self.player.HandleTask(TASK_TYPE_FINISH_CHAPTER, csvLevelConfig.LevelIndex, 0, 0)
				}

				self.player.GetModule("lifetree").(*ModLifeTree).CheckPass()
			}
		}
		self.SendInfo()
		self.player.SendInfo("updateuserinfo")
		configTemp := GetCsvMgr().LevelConfigMap[chapter+1]
		if configTemp == nil {
			self.player.GetModule("onhook").(*ModOnHook).SetStage(chapter)
		} else {
			self.player.GetModule("onhook").(*ModOnHook).SetStage(chapter + 1)
		}
		return
	}
	totalExp := 0
	passid := 0
	for i := 1; i <= chapter; i++ {
		config, ok := GetCsvMgr().MissionMap[i]
		if !ok {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PASS_CHAPTERS_DO_NOT_EXIST"))
			return
		}

		for _, missionId := range config.MissionIds {
			csvLevelConfig := GetCsvMgr().LevelConfigMap[missionId]
			if csvLevelConfig == nil {
				return
			}
			addStar := 3
			if csvLevelConfig.LevelType != CAL_STAR_TYPE {
				addStar = 0
			}

			mission := self.GetPass(missionId)
			if mission == nil {
				var node PassInfo
				node.Id = missionId
				node.Num = 1
				node.Star = addStar
				self.Sql_UserPass.passinfo = append(self.Sql_UserPass.passinfo, node)
				levelConfig, hasMission := GetCsvMgr().LevelConfigMap[missionId]
				if hasMission {
					totalExp += levelConfig.TeamExp
				}
			} else {
				mission.Num += 1
				mission.Star = addStar
			}
			if passid < missionId {
				passid = missionId
			}
			self.player.HandleTask(TASK_TYPE_FINISH_PASS, missionId, 1, 0)

			if csvLevelConfig.MainType == 1 {
				self.player.HandleTask(TASK_TYPE_FINISH_MAIN_PASS, csvLevelConfig.TaskIndex, 0, 0)
			}

			//完成章节任务
			configNext := GetCsvMgr().LevelConfigMap[csvLevelConfig.NextLevel]
			if configNext == nil || csvLevelConfig.LevelIndex != configNext.LevelIndex {
				self.player.HandleTask(TASK_TYPE_FINISH_CHAPTER, csvLevelConfig.LevelIndex, 0, 0)
			} else {
				self.player.HandleTask(TASK_TYPE_FINISH_CHAPTER, csvLevelConfig.LevelIndex-1, 0, 0)
			}

			self.player.GetModule("lifetree").(*ModLifeTree).CheckPass()
		}
	}

	totalstars := 0
	for i := 0; i < len(self.Sql_UserPass.passinfo); i++ {
		if self.Sql_UserPass.passinfo[i].Id > totalstars {
			totalstars = self.Sql_UserPass.passinfo[i].Id
		}
	}
	self.Sql_UserPass.Totalstars = totalstars
	self.Sql_UserPass.StartTime = TimeServer().Unix()
	//GetTopPassMgr().UpdateRank(totalstars, self.player)

	self.player.AddExp(totalExp, 0, 0, "gm命令")

	self.SendInfo()

	self.player.SendInfo("updateuserinfo")

	self.player.GetModule("onhook").(*ModOnHook).SetStage(passid)
}

func (self *ModPass) GetLastPass() *PassInfo {
	if len(self.Sql_UserPass.passinfo) > 0 {
		i := 0
		max := 0
		for k, v := range self.Sql_UserPass.passinfo {
			if v.Id > max {
				max = v.Id
				i = k
			}
		}
		return &self.Sql_UserPass.passinfo[i]
	}
	return nil
}

func (self *ModPass) CalPitExt(config *NewPitConfig) int {
	//根据关卡进度拿到关卡配置
	stage := self.player.GetModule("onhook").(*ModOnHook).GetStage()
	levelConfig := GetCsvMgr().LevelConfigMap[stage]
	if levelConfig == nil {
		return 0
	}
	//计算
	index := -1
	if config.NumberOfPlies == 1 {
		index = 0
	} else if config.NumberOfPlies == 2 {
		index = 1
	} else if config.Difficulty == 1 {
		index = 2
	} else if config.Difficulty == 2 {
		index = 3
	}
	if index >= 0 {
		return levelConfig.DungeonExtraDrop[index]
	} else {
		return 0
	}
}

// 该关卡是否存在
func (self *ModPass) AddPass(passid int, star int) {

	_, ok := GetCsvMgr().LevelConfigMap[passid]
	if !ok {
		return
	}

	find := false
	for i := 0; i < len(self.Sql_UserPass.passinfo); i++ {
		if self.Sql_UserPass.passinfo[i].Id == passid {
			find = true
			break
		}
	}

	if find == false {
		var pass PassInfo
		pass.Id = passid
		pass.Num = 1
		pass.Star = star

		self.Sql_UserPass.passinfo = append(self.Sql_UserPass.passinfo, pass)
		self.player.HandleTask(TASK_TYPE_FINISH_PASS, passid, 1, 0)
	}
}

func (self *ModPass) GetBattleInfo(id int64) *BattleInfo {
	var battleInfo BattleInfo
	value, flag, err := HGetRedisEx(`san_passbattleinfo`, id, fmt.Sprintf("%d", id))
	if err != nil || !flag {
		return GetServer().DBUser.GetBattleInfo(id)
	}
	if flag {
		err := json.Unmarshal([]byte(value), &battleInfo)
		if err != nil {
			return &battleInfo
		}
	}

	if battleInfo.Id != 0 {
		return &battleInfo
	}
	return nil
}

func (self *ModPass) GetBattleRecord(id int64) *BattleRecord {
	var battleRecord BattleRecord

	value, flag, err := HGetRedisEx(`san_passbattlerecord`, id, fmt.Sprintf("%d", id))
	if err != nil || !flag {
		return GetServer().DBUser.GetBattleRecord(id)
	}
	if flag {
		err := json.Unmarshal([]byte(value), &battleRecord)
		if err != nil {
			return &battleRecord
		}
	}

	if battleRecord.Id != 0 {
		return &battleRecord
	}
	return nil
}
