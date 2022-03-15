package game

import (
	"encoding/json"
	"fmt"
)

//! 阵容-
type ModTeam struct {
	player   *Player
	Sql_Team San_Team
}

// 布阵信息, 所有的布阵共用一套装备信息
type San_Team struct {
	Uid     int64
	Info    string
	PreInfo string

	info    []*TeamPos // 布阵状态
	preInfo []*TeamPos // 编组状态

	DataUpdate
}

func (self *ModTeam) Decode() {
	json.Unmarshal([]byte(self.Sql_Team.Info), &self.Sql_Team.info)
	json.Unmarshal([]byte(self.Sql_Team.PreInfo), &self.Sql_Team.preInfo)
}

func (self *ModTeam) Encode() {
	self.Sql_Team.Info = HF_JtoA(self.Sql_Team.info)
	self.Sql_Team.PreInfo = HF_JtoA(self.Sql_Team.preInfo)
}

// 离线拉取
func (self *ModTeam) OnGetData(player *Player) {
	self.player = player
	tableName := self.TableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", tableName, self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Team, tableName, self.player.ID)
	if self.Sql_Team.Uid <= 0 {
		self.Sql_Team.Uid = self.player.ID
		self.Encode()
		InsertTable(tableName, &self.Sql_Team, 0, true)
	} else {
		self.Decode()
	}
	self.Sql_Team.Init(tableName, &self.Sql_Team, true)
	self.checkTeamPos()
}

func NewTeamPos() *TeamPos {
	teamPos := &TeamPos{}
	teamPos.HydraId = 0
	return teamPos
}

func (self *ModTeam) checkTeamPos() {
	if self.Sql_Team.info == nil {
		self.Sql_Team.info = make([]*TeamPos, 0)
		for i := 0; i < TEAM_END-1; i++ {
			self.Sql_Team.info = append(self.Sql_Team.info, NewTeamPos())
		}
	}

	// 新增新的队伍
	size := len(self.Sql_Team.info)
	for i := 0; i < TEAM_END-1-size; i++ {
		self.Sql_Team.info = append(self.Sql_Team.info, NewTeamPos())
	}

	if self.Sql_Team.preInfo == nil {
		self.Sql_Team.preInfo = make([]*TeamPos, 0)
		for i := 0; i < PRETEAM_MAX-1; i++ {
			self.Sql_Team.preInfo = append(self.Sql_Team.preInfo, NewTeamPos())
		}
	}

	// 新增新的队伍
	sizePre := len(self.Sql_Team.preInfo)
	for i := 0; i < PRETEAM_MAX-1-sizePre; i++ {
		self.Sql_Team.preInfo = append(self.Sql_Team.preInfo, NewTeamPos())
	}
}

func (self *ModTeam) TableName() string {
	return "san_teampos"
}

func (self *ModTeam) OnGetOtherData() {
	self.player.GetModule("team").(*ModTeam).CountArenaTopFight(TEAMTYPE_ARENA_2)
	self.player.GetModule("team").(*ModTeam).CountArenaTopFight(TEAMTYPE_ARENA_SPECIAL_4)
}

func (self *ModTeam) OnGetOtherData2() {
	tableName := self.TableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", tableName, self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Team, tableName, self.player.ID)
	if self.Sql_Team.Uid <= 0 {
		self.Sql_Team.Uid = self.player.ID
		self.Encode()
	} else {
		self.Decode()
	}
}

func (self *ModTeam) OnSave(sql bool) {
	self.Encode()
	self.Sql_Team.Update(sql)

}

func (self *ModTeam) OnRefresh() {

}

// 注册消息
func (self *ModTeam) onReg(handlers map[string]func(body []byte)) {
	handlers["addteampos"] = self.addUIPos
	handlers["swapfightpos"] = self.swapFightPos
	handlers["saveteam"] = self.SaveTeam
	handlers["addarenaspcialteampos"] = self.AddArenaSpcialTeamUIPos
}

// 消息处理
func (self *ModTeam) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "clearteampos":
		self.clearTeamPos()
		return true
	case "changehydra":
		//var info C2S_AddHydra
		//json.Unmarshal(body, &info)
		//self.player.GetModule("hydra").(*ModHydra).AddHydra(ctrl, info.HydraId, info.TeamType)
		return true
	}

	return false
}

// 登录发送信息
func (self *ModTeam) SendInfo() {
	self.checkTeamPos()
	//self.checkPos()
	self.sendTeamMsg()
}

func (self *ModTeam) SaveTeam(body []byte) {
	var msg C2S_SaveTeam
	json.Unmarshal(body, &msg)

	if msg.TeamType <= 0 || msg.TeamType > PRETEAM_MAX {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAMTYPE_ERROR"))
		return
	}

	if len(msg.FightPos) != MAX_FIGHT_POS {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_MSG_ERROR"))
		return
	}

	team := self.Sql_Team.preInfo[msg.TeamType-1]
	if team == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_DATA_ERROR"))
		return
	}

	copy(team.FightPos[0:MAX_FIGHT_POS], msg.FightPos)

	msgRel := &S2C_UpdateTeamPos{
		Cid: "saveteam",
		Team: &Js_TeamPos{
			TeamPos:  team,
			TeamType: msg.TeamType,
		},
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
}

// 登录发送信息
func (self *ModTeam) sendTeamMsg() {
	msg := &S2C_TeamPos{}
	msg.Cid = "teampos"
	for index, v := range self.Sql_Team.info {
		msg.TeamPos = append(msg.TeamPos, &Js_TeamPos{
			//UIPos:       v.UIPos,
			TeamType: index + 1,
			TeamPos:  v,
		})
	}
	for index, v := range self.Sql_Team.preInfo {
		msg.PreTeamPos = append(msg.PreTeamPos, &Js_TeamPos{
			TeamType: index + 1,
			TeamPos:  v,
		})
	}
	self.player.SendMsg(msg.Cid, HF_JtoB(msg))
}

// 布阵,2阵容直接从1阵容取装备数据
func (self *ModTeam) AddArenaSpcialTeamUIPos(body []byte) {
	var msg C2S_AddTeamUIPos
	json.Unmarshal(body, &msg)

	if msg.TeamType != TEAMTYPE_ARENA_SPECIAL_1 && msg.TeamType != TEAMTYPE_ARENA_SPECIAL_4 && msg.TeamType != TEAMTYPE_CROSSARENA_ATTACK_3V3_1 && msg.TeamType != TEAMTYPE_CROSSARENA_DEFENCE_3V3_1 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAMTYPE_ERROR"))
		return
	}

	if len(msg.FightPos) != MAX_FIGHT_POS*3 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_MSG_ERROR"))
		return
	}

	if !self.CheckVoidHero(msg.FightPos) {
		return
	}

	pos := []*Js_TeamPos{}
	index := 0
	if msg.TeamType == TEAMTYPE_ARENA_SPECIAL_1 {
		for i := TEAMTYPE_ARENA_SPECIAL_1; i <= TEAMTYPE_ARENA_SPECIAL_3; i++ {
			team := self.Sql_Team.info[i-1]
			if team == nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_DATA_ERROR"))
				return
			}

			copy(team.FightPos[0:MAX_FIGHT_POS], msg.FightPos[index*MAX_FIGHT_POS:(index+1)*MAX_FIGHT_POS])
			pos = append(pos, &Js_TeamPos{i, team})
			index++
		}
	} else if msg.TeamType == TEAMTYPE_ARENA_SPECIAL_4 {
		for i := TEAMTYPE_ARENA_SPECIAL_4; i <= TEAMTYPE_ARENA_SPECIAL_6; i++ {
			team := self.Sql_Team.info[i-1]
			if team == nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_DATA_ERROR"))
				return
			}

			copy(team.FightPos[0:MAX_FIGHT_POS], msg.FightPos[index*MAX_FIGHT_POS:(index+1)*MAX_FIGHT_POS])
			pos = append(pos, &Js_TeamPos{i, team})
			index++
		}
		self.player.GetModule("team").(*ModTeam).CountArenaTopFight(TEAMTYPE_ARENA_SPECIAL_4)
	} else if msg.TeamType == TEAMTYPE_CROSSARENA_ATTACK_3V3_1 {
		for i := TEAMTYPE_CROSSARENA_ATTACK_3V3_1; i <= TEAMTYPE_CROSSARENA_ATTACK_3V3_3; i++ {
			team := self.Sql_Team.info[i-1]
			if team == nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_DATA_ERROR"))
				return
			}

			copy(team.FightPos[0:MAX_FIGHT_POS], msg.FightPos[index*MAX_FIGHT_POS:(index+1)*MAX_FIGHT_POS])
			pos = append(pos, &Js_TeamPos{i, team})
			index++
		}
	} else if msg.TeamType == TEAMTYPE_CROSSARENA_DEFENCE_3V3_1 {
		//更新跨服战防守阵容
		fightInfo := [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo{}
		fightInfo[0] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_DEFENCE_3V3_1)
		fightInfo[1] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_DEFENCE_3V3_2)
		fightInfo[2] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_DEFENCE_3V3_3)
		GetCrossArena3V3Mgr().AddInfo(self.player, fightInfo)

		for i := TEAMTYPE_CROSSARENA_DEFENCE_3V3_1; i <= TEAMTYPE_CROSSARENA_DEFENCE_3V3_3; i++ {
			team := self.Sql_Team.info[i-1]
			if team == nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_DATA_ERROR"))
				return
			}

			copy(team.FightPos[0:MAX_FIGHT_POS], msg.FightPos[index*MAX_FIGHT_POS:(index+1)*MAX_FIGHT_POS])
			pos = append(pos, &Js_TeamPos{i, team})
			index++
		}
	}

	msgRel := &S2C_UpdateArenaSpecialTeamPos{}
	msgRel.Cid = "updatearenaspcialteampos"
	msgRel.Team = pos
	self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))

	GetArenaSpecialMgr().UpdateFormat(self.player)
}

// 布阵,2阵容直接从1阵容取装备数据
func (self *ModTeam) addUIPos(body []byte) {

	var msg C2S_AddTeamUIPos
	json.Unmarshal(body, &msg)

	if msg.TeamType <= 0 || msg.TeamType > TEAM_END {
		//self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAMTYPE_ERROR"))
		return
	}

	if msg.TeamType >= TEAMTYPE_ARENA_SPECIAL_1 && msg.TeamType <= TEAMTYPE_ARENA_SPECIAL_6 {
		// 填判断后三个是否需要判断重复
		if false {
			//self.player.GetModule("arena").(*ModArena).UpdateFormat()
			return
		}
	}
	//升星后会到这里报错，先注了
	if len(msg.FightPos) != MAX_FIGHT_POS {
		//self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_MSG_ERROR"))
		return
	}

	if !self.CheckVoidHero(msg.FightPos) {
		return
	}

	team := self.Sql_Team.info[msg.TeamType-1]
	if team == nil {
		//self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_DATA_ERROR"))
		return
	}

	copy(team.FightPos[0:MAX_FIGHT_POS], msg.FightPos)

	msgRel := &S2C_UpdateTeamPos{
		Cid: "updateteampos",
		Team: &Js_TeamPos{
			TeamPos:  team,
			TeamType: msg.TeamType,
		},
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))

	if msg.TeamType == TEAMTYPE_DEFAULT {
		self.player.updateFight()
	} else if msg.TeamType == TEAMTYPE_ARENA_2 {
		self.player.GetModule("team").(*ModTeam).CountArenaTopFight(TEAMTYPE_ARENA_2)
		GetArenaMgr().UpdateFormat(self.player)
	} else if msg.TeamType >= TEAMTYPE_ARENA_SPECIAL_4 && msg.TeamType <= TEAMTYPE_ARENA_SPECIAL_6 {
		if msg.TeamType == TEAMTYPE_ARENA_SPECIAL_4 {
			self.player.GetModule("team").(*ModTeam).CountArenaTopFight(TEAMTYPE_ARENA_SPECIAL_4)
		}
		GetArenaSpecialMgr().UpdateFormat(self.player)
	} else if msg.TeamType == TEAMTYPE_CROSSARENA_DEFENCE {
		//更新跨服战防守阵容
		fightInfo := GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_CROSSARENA_DEFENCE)
		GetCrossArenaMgr().AddInfo(self.player, fightInfo)
	}
}

func (self *ModTeam) InitNew() {
	pTeam := self.Sql_Team.info[0]

	heroes := self.player.getHeroes()
	for _, value := range GetCsvMgr().NewUserItem {
		if value.Group == 1 {
			switch value.Type {
			case 1:
				if value.Place > 0 {
					for _, v := range heroes {
						if v.HeroId == value.Id {
							pTeam.addUIPos(value.Place-1, v.HeroKeyId)
							break
						}
					}
				}
			}
		}
	}

	self.player.countTeamFight(0)
}

// 交换战斗位置
func (self *ModTeam) swapFightPos(body []byte) {

	var msg C2S_SwapFightPos
	json.Unmarshal(body, &msg)

	if msg.TeamType < 1 || msg.TeamType >= TEAM_END {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAM_TYPE_ERROR"))
		return
	}

	if msg.Index1 == msg.Index2 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAM_TYPE_ERROR"))
		return
	}

	teamIndex := msg.TeamType - 1
	if teamIndex < 0 || teamIndex >= len(self.Sql_Team.info) {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAM_INFORMATION_DOES_NOT_EXIST"))
		return
	}

	pTeam := self.Sql_Team.info[teamIndex]
	if pTeam == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAM_INFORMATION_IS_EMPTY"))
		return
	}

	err := pTeam.swapFightPosByIndex(msg.Index1-1, msg.Index2-1)
	if err != nil {
		self.player.SendErr(err.Error())
		return
	}

	msgRel := &S2C_UpdateTeamPos{
		Cid: "updateteampos",
		Team: &Js_TeamPos{
			TeamPos:  pTeam,
			TeamType: msg.TeamType,
		},
	}
	self.player.SendMsg(msgRel.Cid, HF_JtoB(msgRel))
}

// 检查第一个空站位
func (self *ModTeam) checkPos() {
	if self.Sql_Team.info == nil {
		return
	}
	teamIndex := 0
	pTeam := self.Sql_Team.info[teamIndex]
	if pTeam != nil && pTeam.isUIPosEmpty() {
		// 随机选一个英雄上阵
		heroKeyId := self.player.GetModule("hero").(*ModHero).randHero()
		pTeam.addUIPos(teamIndex, heroKeyId)
	}
}

// 下阵巨兽
func (self *ModTeam) offBoss(id int) bool {
	/*
		if len(self.Sql_Team.info) <= 0 {
			return false
		}

		bossConfig := GetCsvMgr().GetBossConfig(id)
		if bossConfig == nil {
			return false
		}

		bossId := bossConfig.HeroId
		for i, pTeam := range self.Sql_Team.info {
			if pTeam == nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAM_INFORMATION_IS_EMPTY"))
				return false
			}

			res := pTeam.removeUIPos(bossId)
			if !res {
				return false
			}

			msg := &S2C_UpdateTeamPos{
				Cid: "bossembattle",
				TeamPos: &Js_TeamPos{
					//UIPos:       pTeam.UIPos,
					FightPos:    pTeam.FightPos,
					TeamType:    i + 1,
					FormationId: pTeam.FormationId,
					HydraId:     pTeam.HydraId,
				},
			}
			self.player.SendMsg(msg.Cid, HF_JtoB(msg))
			return true
		}
	*/
	return false
}

func (self *ModTeam) getTeamPos(teamType int) *TeamPos {
	if teamType < 1 || teamType > len(self.Sql_Team.info) {
		return nil
	}

	return self.Sql_Team.info[teamType-1]
}

func (self *ModTeam) resetTeamPos(teamType int) {
	team := self.getTeamPos(teamType)
	for i := 0; i < len(team.FightPos); i++ {
		team.FightPos[i] = 0
	}
	//self.player.GetModule("arena").(*ModArena).UpdateFormat()
}

func (self *ModTeam) getBossId() int {
	if len(self.Sql_Team.info) <= 0 {
		return 0
	}

	firstTeam := self.Sql_Team.info[0]
	if firstTeam == nil {
		return 0
	}
	//return firstTeam.UIPos[MAX_UI_POS-1]
	return 0
}

func (self *ModTeam) getFirstTeam() []int {
	if len(self.Sql_Team.info) <= 0 {
		return []int{0, 0, 0, 0, 0, 0}
	}

	firstTeam := self.Sql_Team.info[0]
	if firstTeam == nil {
		return []int{0, 0, 0, 0, 0, 0}
	}

	var res []int
	for _, heroId := range firstTeam.FightPos {
		res = append(res, heroId)
	}

	return res
}

func (self *ModTeam) clearTeamPos() {

	return

	/*
		if len(self.Sql_Team.info) <= 0 {
			return
		}
		pTeam := self.Sql_Team.info[0]
		if pTeam == nil {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAM_INFORMATION_IS_EMPTY"))
			return
		}

		for i := range pTeam.UIPos {
			if i == 0 {
				continue
			}
			pTeam.UIPos[i] = 0
		}

		for i := range pTeam.FightPos {
			if pTeam.FightPos[i] == pTeam.UIPos[0] {
				continue
			}
			pTeam.FightPos[i] = 0
		}

		msg := &S2C_UpdateTeamPos{
			Cid: "gm_clear_teampos",
			TeamPos: &Js_TeamPos{
				UIPos:       pTeam.UIPos,
				FightPos:    pTeam.FightPos,
				TeamType:    1,
				TeamAttr:    pTeam.TeamAttr,
				FormationId: pTeam.FormationId,
				HydraId:     pTeam.HydraId,
			},
		}
		self.player.SendMsg(msg.Cid, HF_JtoB(msg))
		self.checkEmbattle()
		self.player.countTeamFight(ReasonGmTeam)

	*/
}

// 获取当前玩家队伍总天赋
func (self *ModTeam) CacTalents() {
	/*
		team := self.getFirstTeam()
		total := 0
		for i := 0; i < 5; i++ {
			heroId := team[i]
			if heroId == 0 {
				continue
			}

			hero := self.player.getHero(heroId)
			if hero == nil {
				continue
			}

			if hero.TalentItem == nil {
				continue
			}

			for _, v := range hero.TalentItem.Talents {
				total += v.Lv
			}
		}

		pInfo := &self.player.GetModule("hero").(*ModHero).Sql_Hero
		if pInfo.Totalstars < total {
			pInfo.Totalstars = total
			pInfo.StarTime = TimeServer().Unix()
		}

		//self.player.Send("syntalentnum", HF_JtoB(&S2C_SynTalentNum{
		//	Cid:      "syntalentnum",
		//	TotalNum: pInfo.Totalstars,
		//}))

	*/
}

//检查英雄是否存在于阵位上
func (self *ModTeam) CheckHeroInTeam(heroKeyId int) bool {

	for _, v := range self.Sql_Team.info {
		for i := 0; i < len(v.FightPos); i++ {
			if v.FightPos[i] == heroKeyId {
				return true
			}
		}
	}
	for _, v := range self.Sql_Team.preInfo {
		for i := 0; i < len(v.FightPos); i++ {
			if v.FightPos[i] == heroKeyId {
				return true
			}
		}
	}
	return false
}

func (self *ModTeam) DeleteHeroFromTeam(heroKeyId int) {

	for _, v := range self.Sql_Team.info {
		for i := 0; i < len(v.FightPos); i++ {
			if v.FightPos[i] == heroKeyId {
				v.FightPos[i] = 0
			}
		}
	}
	for _, v := range self.Sql_Team.preInfo {
		for i := 0; i < len(v.FightPos); i++ {
			if v.FightPos[i] == heroKeyId {
				v.FightPos[i] = 0
			}
		}
	}
}

// 检查虚空英雄是否合规则
func (self *ModTeam) CheckVoidHero(fightPos []int) bool {
	heromod := self.player.GetModule("hero").(*ModHero)
	if nil == heromod {
		return false
	}
	for _, v := range fightPos {
		hero := heromod.GetHero(v)
		if hero == nil {
			continue
		}

		if hero.VoidHero != 0 {
			continue
		}

		if hero.Resonance == 0 {
			continue
		}

		for _, p := range fightPos {
			if p == hero.Resonance {
				return false
			}
		}
	}

	return true
}

func (self *ModTeam) GmClearItem() {
	self.Sql_Team.info = nil
	self.Sql_Team.preInfo = nil
	self.SendInfo()
}

func (self *ModTeam) CountArenaTopFight(teamType int) {
	var herolst []int
	if teamType == TEAMTYPE_ARENA_2 {
		team := self.player.GetModule("team").(*ModTeam).getTeamPos(TEAMTYPE_ARENA_2)
		for _, p := range team.FightPos {
			herolst = append(herolst, p)
		}
	} else if teamType == TEAMTYPE_ARENA_SPECIAL_4 {
		team := self.player.GetModule("team").(*ModTeam).getTeamPos(TEAMTYPE_ARENA_SPECIAL_4)
		for _, p := range team.FightPos {
			herolst = append(herolst, p)
		}
		team = self.player.GetModule("team").(*ModTeam).getTeamPos(TEAMTYPE_ARENA_SPECIAL_5)
		for _, p := range team.FightPos {
			herolst = append(herolst, p)
		}
		team = self.player.GetModule("team").(*ModTeam).getTeamPos(TEAMTYPE_ARENA_SPECIAL_6)
		for _, p := range team.FightPos {
			herolst = append(herolst, p)
		}
	}
	fight := self.player.GetModule("hero").(*ModHero).GetFight(herolst) / 100

	GetOfflineInfoMgr().SetArenaFight(teamType, self.player.GetUid(), fight)
}

func (self *ModTeam) AddCrossArena3V3(body []byte) {
	var msg C2S_AddTeamUIPos
	json.Unmarshal(body, &msg)

	if msg.TeamType != TEAMTYPE_CROSSARENA_ATTACK_3V3_1 && msg.TeamType != TEAMTYPE_CROSSARENA_DEFENCE_3V3_1 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_TEAMTYPE_ERROR"))
		return
	}

	if len(msg.FightPos) != MAX_FIGHT_POS*3 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_MSG_ERROR"))
		return
	}

	if !self.CheckVoidHero(msg.FightPos) {
		return
	}

	pos := []*Js_TeamPos{}
	index := 0
	if msg.TeamType == TEAMTYPE_CROSSARENA_ATTACK_3V3_1 {
		for i := TEAMTYPE_CROSSARENA_ATTACK_3V3_1; i <= TEAMTYPE_CROSSARENA_ATTACK_3V3_3; i++ {
			team := self.Sql_Team.info[i-1]
			if team == nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_DATA_ERROR"))
				return
			}

			copy(team.FightPos[0:MAX_FIGHT_POS], msg.FightPos[index*MAX_FIGHT_POS:(index+1)*MAX_FIGHT_POS])
			pos = append(pos, &Js_TeamPos{i, team})
			index++
		}
	} else if msg.TeamType == TEAMTYPE_CROSSARENA_DEFENCE_3V3_1 {
		for i := TEAMTYPE_CROSSARENA_DEFENCE_3V3_1; i <= TEAMTYPE_CROSSARENA_DEFENCE_3V3_3; i++ {
			team := self.Sql_Team.info[i-1]
			if team == nil {
				self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TEAM_DATA_ERROR"))
				return
			}

			copy(team.FightPos[0:MAX_FIGHT_POS], msg.FightPos[index*MAX_FIGHT_POS:(index+1)*MAX_FIGHT_POS])
			pos = append(pos, &Js_TeamPos{i, team})
			index++
		}
	}
}
