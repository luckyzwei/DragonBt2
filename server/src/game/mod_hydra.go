package game

/*
import (
	"encoding/json"
	"fmt"
)

const (
	HYDRA_ID       = 5000
	HydraLvUp      = 1
	HydraSkillLvUp = 2
	HydraSetpUp    = 3
	HydraReborn    = 4
	//HydraCheck  = 3
	//HydraOff    = 4
	HydraAction = "hydraaction"
)

const HYDRA_BASE_STEP = 1

type HydraSkillAtt struct {
	HydraSkillID int `json:"hydraskillid"` // 神兽技能id
	HydraSkillLv int `json:"hydraskilllv"` // 神兽技能等级，影响技能
}

// 神兽
type HydraInfo struct {
	Id            int              `json:"id"`         // 神兽Id
	HydraLv       int              `json:"hydralv"`    // 神兽等级，影响属性
	HydraStep     int              `json:"hydrastep"`  // 神兽阶
	HydraSkillAtt []*HydraSkillAtt `json:"hydraskill"` // 技能
	Fight         int64
}

// 神兽任务
type HydraTaskInfo struct {
	Taskid    int `json:"taskid"`    // 任务Id
	Tasktypes int `json:"tasktypes"` // 任务类型
	Plan      int `json:"plan"`      // 进度
	State     int `json:"finish"`    // 状态
}

type SanHydra struct {
	Uid       int64
	HydraInfo string

	hydraInfos map[int]*HydraInfo // 神兽信息

	DataUpdate
}

type ModHydra struct {
	player *Player
	Data   SanHydra
	chg    []HydraTaskInfo
}

func (self *ModHydra) Decode() {
	err := json.Unmarshal([]byte(self.Data.HydraInfo), &self.Data.hydraInfos)
	if err != nil {
		LogError(err.Error())
	}
}

func (self *ModHydra) Encode() {
	self.Data.HydraInfo = HF_JtoA(self.Data.hydraInfos)
}

func (self *ModHydra) OnGetData(player *Player) {
	self.player = player
	tableName := self.getTableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", tableName, self.player.ID)
	//fmt.Println(sql)
	GetServer().DBUser.GetOneData(sql, &self.Data, tableName, self.player.ID)
	//litter.Dump(self.Data)
	if self.Data.Uid <= 0 {
		self.init(self.player.ID)
		self.Encode()
		InsertTable(tableName, &self.Data, 0, true)
	} else {
		self.Decode()
	}

	self.Data.Init(tableName, &self.Data, true)
	for _, v := range self.Data.hydraInfos {
		self.updateFight(v.Id)
	}
}

func (self *ModHydra) OnRefresh() {

}

func (self *ModHydra) OnMsg(ctrl string, body []byte) bool {
	return false
}

func (self *ModHydra) onReg(handlers map[string]func(body []byte)) {
	handlers["hydraaction"] = self.onHydraAction
	//handlers["takehydratask"] = self.TakeHydraTask
	handlers["takehydra"] = self.TakeHydra
}

// 神兽相关操作
func (self *ModHydra) onHydraAction(body []byte) {
	var msg C2S_HydraAction
	err := json.Unmarshal(body, &msg)
	if err != nil {
		LogError(err.Error())
	}
	if msg.Action == HydraLvUp {
		self.HydraLvUp(&msg)
	} else if msg.Action == HydraSkillLvUp {
		self.HydraSkillLvUp(&msg)
	} else if msg.Action == HydraSetpUp {
		self.HydraStepUp(&msg)
	} else if msg.Action == HydraReborn {
		self.HydraReborn(&msg)
	}
}

func (self *ModHydra) getTableName() string {
	return "san_userhydra"
}

func (self *ModHydra) init(uid int64) {
	self.Data.Uid = uid
	self.Data.hydraInfos = make(map[int]*HydraInfo)
}

func (self *ModHydra) OnGetOtherData() {
	config, ok := GetCsvMgr().HydraStepMap[HYDRA_BASE_STEP]
	if !ok || config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("ModHydra_CONFIG_ERROR"))
	}
}

func (self *ModHydra) OnSave(sql bool) {
	self.Encode()
	self.Data.Update(sql)
}

// 同步Hydra信息
func (self *ModHydra) SendInfo() {
	var msg S2C_HydraInfo
	msg.Cid = "hydrainfo"
	msg.HydraInfos = self.Data.hydraInfos
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

// 获取巨兽战斗力  增加个神兽参数比较好
func (self *ModHydra) updateFight(hydraId int) {

	hydra := self.Data.hydraInfos[hydraId]
	if hydra == nil {
		return
	}

	hydra.Fight = 0

	for _, v := range hydra.HydraSkillAtt {
		config, ok1 := GetCsvMgr().HydraSkillMap[v.HydraSkillID]
		if !ok1 || config == nil {
			continue
		}

		skill, ok2 := config[v.HydraSkillLv]
		if !ok2 || skill == nil {
			return
		}

		for i := 0; i < len(skill.Base_type); i++ {
			if skill.Base_type[i] == ATTR_TYPE_FIGHT {
				hydra.Fight += skill.Base_value[i]
			}
		}
	}
	team := self.player.GetModule("team").(*ModTeam).getTeamPos(TEAMTYPE_DEFAULT)
	if team != nil && team.HydraId == hydraId {
		self.player.updateFight()
	}
}

// 激活神兽
func (self *ModHydra) TakeHydra(body []byte) {
	var msg C2STakeHydra
	json.Unmarshal(body, &msg)

	hydraId := msg.HydraID

	hydraConfig, ok := GetCsvMgr().HydraConfigMap[hydraId]
	if !ok || hydraConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("ModHydra_CONFIG_ERROR"))
	}

	config, ok := GetCsvMgr().HydraStepMap[HYDRA_BASE_STEP]
	if !ok || config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("ModHydra_CONFIG_ERROR"))
	}

	if err := self.player.HasObjectOk(config.Items, config.Nums); err != nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("cost is not enough"))
		return
	}

	items := self.player.RemoveObjectLst(config.Items, config.Nums, "激活巨兽", hydraId, 0, 0)

	self.AddHydro(hydraId)
	self.updateFight(hydraId)

	var sendmsg S2CTakeHydra
	sendmsg.Cid = "takehydra"
	sendmsg.HydraInfos = self.Data.hydraInfos
	sendmsg.Items = items
	self.player.SendMsg(sendmsg.Cid, HF_JtoB(&sendmsg))
}

func (self *ModHydra) AddHydro(hydraId int) {
	//看是否已经拥有
	_, ok := self.Data.hydraInfos[hydraId]
	//已经拥有
	if ok {
		return
	}

	hydra := self.NewHydra(hydraId)
	self.Data.hydraInfos[hydraId] = hydra
}

func (self *ModHydra) GetHydroFight(hydraId int) int64 {
	//看是否已经拥有
	hydra, ok := self.Data.hydraInfos[hydraId]
	//已经拥有
	if !ok {
		return 0
	}

	return hydra.Fight
}

func (self *ModHydra) NewHydra(hydraId int) *HydraInfo {
	hydra := new(HydraInfo)
	hydra.Id = hydraId
	hydra.HydraLv = 1
	hydra.HydraStep = 1
	hydra.AddSkill(hydra.HydraStep)
	hydra.Fight = 0
	return hydra
}

func (self *HydraInfo) AddSkill(step int) {
	config, ok := GetCsvMgr().HydraConfigMap[self.Id]
	if !ok || config == nil {
		return
	}

	if step <= 0 || len(config.Skill) < step {
		return
	}

	skill := config.Skill[step-1]
	if skill == 0 {
		return
	}

	self.HydraSkillAtt = append(self.HydraSkillAtt, &HydraSkillAtt{skill, 1})
}

// 神兽升级
func (self *ModHydra) HydraLvUp(msg *C2S_HydraAction) {

	hydra := self.Data.hydraInfos[msg.Id]
	if hydra == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("Hydra_IS_NO"))
		msgRel := &S2C_HydraLvUp{
			Cid:         HydraAction,
			Action:      msg.Action,
			HydraInfo:   hydra,
			SpecialStop: SPECIAL_STOP,
		}
		self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
		return
	}

	config := GetCsvMgr().HydraLevelMap[hydra.HydraLv]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("config is error"))
		msgRel := &S2C_HydraLvUp{
			Cid:         HydraAction,
			Action:      msg.Action,
			HydraInfo:   hydra,
			SpecialStop: SPECIAL_STOP,
		}
		self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
		return
	}
	//看等级限制
	configMax := GetCsvMgr().HydraStepMap[hydra.HydraStep]
	if configMax == nil || hydra.HydraLv >= configMax.MaxLevel {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("level is max"))
		msgRel := &S2C_HydraLvUp{
			Cid:         HydraAction,
			Action:      msg.Action,
			HydraInfo:   hydra,
			SpecialStop: SPECIAL_STOP,
		}
		self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
		return
	}

	// 检查消耗是否正常
	err := self.player.HasObjectOk(config.Items, config.Nums)
	if err != nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("货币不足"))
		msgRel := &S2C_HydraLvUp{
			Cid:         HydraAction,
			Action:      msg.Action,
			HydraInfo:   hydra,
			SpecialStop: SPECIAL_STOP,
		}
		self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
		return
	}

	items := self.player.RemoveObjectLst(config.Items, config.Nums, "神兽升级", 0, 0, 0)
	hydra.HydraLv++
	self.updateFight(hydra.Id)

	msgRel := &S2C_HydraLvUp{
		Cid:       HydraAction,
		Action:    msg.Action,
		Items:     items,
		HydraInfo: hydra,
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
}

// 神兽技能升级
func (self *ModHydra) HydraSkillLvUp(msg *C2S_HydraAction) {

	hydra := self.Data.hydraInfos[msg.Id]
	if hydra == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("Hydra_IS_NO"))
		return
	}

	if msg.Index > len(hydra.HydraSkillAtt) || msg.Index <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("Hydra_IS_NO"))
		return
	}

	skill := hydra.HydraSkillAtt[msg.Index-1]
	if nil == skill {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("skill is error"))
		return
	}

	config := GetCsvMgr().HydraSkillMap[skill.HydraSkillID][skill.HydraSkillLv+1]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("config is error"))
		return
	}

	if hydra.HydraStep < config.Unlock {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("config is error"))
		return
	}

	// 检查消耗是否正常
	err := self.player.HasObjectOk(config.Items, config.Nums)
	if err != nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("cost is not enough"))
		return
	}
	items := self.player.RemoveObjectLst(config.Items, config.Nums, "神兽升级", 0, 0, 0)
	skill.HydraSkillLv++

	msgRel := &S2C_HydraLvUp{
		Cid:       HydraAction,
		Action:    msg.Action,
		Items:     items,
		HydraInfo: hydra,
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
}

// 神兽技能升级
func (self *ModHydra) HydraStepUp(msg *C2S_HydraAction) {
	hydra := self.Data.hydraInfos[msg.Id]
	if hydra == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("Hydra_IS_NO"))
		return
	}

	config := GetCsvMgr().HydraStepMap[hydra.HydraStep+1]
	if config == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("config is error"))
		return
	}

	configMax := GetCsvMgr().HydraStepMap[hydra.HydraStep]
	if configMax == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("config is error"))
		return
	}

	if hydra.HydraLv < configMax.MaxLevel {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("config is error"))
		return
	}

	// 检查消耗是否正常
	err := self.player.HasObjectOk(config.Items, config.Nums)
	if err != nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("cost is not enough"))
		return
	}
	items := self.player.RemoveObjectLst(config.Items, config.Nums, "神兽升级", 0, 0, 0)
	hydra.HydraStep++
	hydra.AddSkill(hydra.HydraStep)

	msgRel := &S2C_HydraLvUp{
		Cid:       HydraAction,
		Action:    msg.Action,
		Items:     items,
		HydraInfo: hydra,
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
}

func (self *ModHydra) HydraReborn(msg *C2S_HydraAction) {
	hydra := self.Data.hydraInfos[msg.Id]
	if hydra == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("Hydra_IS_NO"))
		return
	}

	//看看等级是否大于1
	if hydra.HydraLv <= 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("HYDRA_CANT_REBORN"))
		return
	}

	getItem := make(map[int]*Item)
	//返还等级材料
	configLv := GetCsvMgr().HydraLevelMap[hydra.HydraLv]
	if configLv == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("config is error"))
		return
	}
	AddItemMapHelper(getItem, configLv.Items, configLv.ReturnNums)
	hydra.HydraLv = 1
	//返还进阶材料
	for i := hydra.HydraStep; i > 1; i-- {
		configStep := GetCsvMgr().HydraStepMap[hydra.HydraStep]
		if configStep == nil {
			continue
		}
		AddItemMapHelper(getItem, configStep.Items, configStep.Nums)
	}
	hydra.HydraStep = 1
	//返还技能材料
	for _, v := range hydra.HydraSkillAtt {
		if v == nil {
			continue
		}

		if v.HydraSkillLv <= 1 {
			continue
		}

		for i := v.HydraSkillLv; i > 1; i-- {
			configSkill := GetCsvMgr().HydraSkillMap[v.HydraSkillID][v.HydraSkillLv]
			if configSkill == nil {
				continue
			}
			AddItemMapHelper(getItem, configSkill.Items, configSkill.RetuenNums)
		}
	}

	hydra.HydraSkillAtt = make([]*HydraSkillAtt, 0)
	hydra.AddSkill(hydra.HydraStep)
	items := self.player.AddObjectItemMap(getItem, "巨兽重生", 0, 0, 0)

	msgRel := &S2C_HydraLvUp{
		Cid:       HydraAction,
		Action:    msg.Action,
		Items:     items,
		HydraInfo: hydra,
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
}

func (self *ModHydra) AddHydra(ctrl string, id int, teamType int) {
	if teamType < 1 || teamType >= TEAM_END {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAM_TYPE_ERROR"))
		return
	}

	pTeam := self.player.GetModule("team").(*ModTeam).getTeamPos(teamType)
	if pTeam == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAM_INFORMATION_IS_EMPTY"))
		return
	}

	hydra := self.Data.hydraInfos[id]
	if hydra == nil {
		self.player.SendErr(GetCsvMgr().GetText("hydra is not exist"))
		return
	}

	if pTeam.HydraId == id {
		self.player.SendErr(GetCsvMgr().GetText("this is now hydra"))
		return
	}

	pTeam.HydraId = id

	msg := &S2C_UpdateTeamPos{
		Cid: "updateteampos",
		Team: &Js_TeamPos{
			TeamPos:  pTeam,
			TeamType: teamType,
		},
	}
	self.player.SendMsg(msg.Cid, HF_JtoB(msg))
}

*/
