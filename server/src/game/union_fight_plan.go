package game

import (
	"sync"
)

// 战斗安排：一轮信息, 存放玩家以及战斗安排的所有信息

type UnionRecord struct {
	FightId int64 `json:"fight_id"` // 战报Id
	Player1 int64 `json:"team1"`    // 队伍1
	Player2 int64 `json:"team2"`    // 队伍2
	Result  int   `json:"result"`   // 结果: 1攻击方胜利 2防守方胜利 0平局
	Time    int64 `json:"time"`     // 产生时间
}

type RoundPlan struct {
	GroupsA      [5][]int64        `json:"groups_a"`      // A队安排, ok
	GroupsB      [5][]int64        `json:"groups_b"`      // B队安排, ok
	Records      [5]int            `json:"records"`       // 比赛结果, ok
	TeamMatch    [2]int            `json:"team_match"`    // 队伍匹配编号(1-3）, ok
	UnionNames   [2]string         `json:"union_names"`   // 军团名字, ok
	UnionMasters [2]string         `json:"union_masters"` // 军团长名字, ok
	UnionID      [2]int            `json:"union_id"`      // 军团ID
	IndexFight   [2]int            `json:"index_fight"`   // 正在打的两个人的下标, ok
	Camp         [2]int            `json:"camp"`          // 两边阵营, ok
	Result       int               `json:"result"`        // 战斗结果,ok
	DeadA        [5][]int64        `json:"dead_a"`        // A队阵亡的玩家, ok
	DeadB        [5][]int64        `json:"dead_b"`        // B队阵亡的玩家, ok
	FightRecord  [5][]*UnionRecord `json:"fight_record"`  // 战报信息
	UnionIcons   [2]int            `json:"union_icons"`   // 军团ICON 20190428 by zy
	MVP          [2]int64          `json:"union_mvp"`     //! mvp   20190428 by zy
	MVPKill      [2]int            `json:"union_mvpkill"` //! mvp   20190428 by zy
}

func (m *UnionPlan) CheckTeamDead(teamId int) bool {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()

	for _, v := range m.RoundPlans {
		if v.TeamMatch[0] == teamId && v.Result == 2 {
			return true
		} else if v.TeamMatch[1] == teamId && v.Result == 1 {
			return true
		}
	}
	return false
}

//获取战绩
func (m *UnionPlan) GetWinLoseRecord(teamId int) int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()

	recordNum := 0
	for k, v := range m.RoundPlans {
		if v.TeamMatch[0] == teamId {
			if v.Result == 2 {
				recordNum = k + 1
			}
		} else if v.TeamMatch[1] == teamId {
			if v.Result == 1 {
				recordNum = k + 1
			}
		}
	}
	return recordNum
}

func (m *UnionPlan) GetGroupNum(round, group int) (int, int) {
	if !m.CheckParam(round, group) {
		return 0, 0
	}

	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	indexRound := round - 1
	indexGroup := group - 1
	num1 := len(m.RoundPlans[indexRound].GroupsA[indexGroup])
	num2 := len(m.RoundPlans[indexRound].GroupsB[indexGroup])

	return num1, num2
}

func (m *UnionPlan) HasRecord(round, group int) bool {
	if !m.CheckParam(round, group) {
		return false
	}
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	indexRound := round - 1
	indexGroup := group - 1
	record := m.RoundPlans[indexRound].Records[indexGroup]
	return record > 0
}

func (m *UnionPlan) SetUnionID(round int, id1 int, id2 int) {
	if !m.CheckRound(round) {
		return
	}

	m.InfoLocker.Lock()
	defer m.InfoLocker.Unlock()
	indexRound := round - 1
	m.RoundPlans[indexRound].UnionID[0] = id1
	m.RoundPlans[indexRound].UnionID[1] = id2
}

func (m *UnionPlan) SetUnionIcon(round int, icon1 int, icon2 int) {
	if !m.CheckRound(round) {
		return
	}

	m.InfoLocker.Lock()
	defer m.InfoLocker.Unlock()
	indexRound := round - 1
	m.RoundPlans[indexRound].UnionIcons[0] = icon1
	m.RoundPlans[indexRound].UnionIcons[1] = icon2
}

// 战斗安排(两支队伍), 存盘，初始化
// 单场战斗信息记录
type UnionPlan struct {
	RoundPlans [3]*RoundPlan `json:"round_plans"` // 每轮信息, ok
	WinnerId   int           `json:"winner_id"`   // 最后胜利, ok
	RoundNum   int           `json:"round_num"`   // 有几轮, ok
	TheRound   int           `json:"the_round"`   // 打到哪一轮, ok
	TheGroup   int           `json:"the_group"`   // 打到哪一组, ok

	InfoLocker *sync.RWMutex `json:"-"` // 锁
}

func (m *UnionPlan) GetUnionID(round int, index int) int {
	if !m.CheckRound(round) {
		return 0
	}
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	indexRound := round - 1
	return m.RoundPlans[indexRound].UnionID[index]
}

func (m *UnionPlan) GetUnionCamp(round int, index int) int {
	if !m.CheckRound(round) {
		return 0
	}
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	indexRound := round - 1
	return m.RoundPlans[indexRound].Camp[index]
}

func (m *UnionPlan) GetUnionCamps(round int) [2]int {
	if !m.CheckRound(round) {
		return [2]int{}
	}
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	indexRound := round - 1
	return m.RoundPlans[indexRound].Camp
}

func (m *UnionPlan) GetPlayersA(round int) []int64 {
	if !m.CheckRound(round) {
		return []int64{}
	}
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	indexRound := round - 1
	var ids []int64
	for _, v := range m.RoundPlans[indexRound].GroupsA {
		for _, v2 := range v {
			ids = append(ids, v2)
		}
	}

	return ids
}

func (m *UnionPlan) GetPlayersB(round int) []int64 {
	if !m.CheckRound(round) {
		return []int64{}
	}
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	indexRound := round - 1
	var ids []int64
	for _, v := range m.RoundPlans[indexRound].GroupsB {
		for _, v2 := range v {
			ids = append(ids, v2)
		}
	}

	return ids
}

func (m *UnionPlan) GetUnionName(round int) [2]string {
	if !m.CheckRound(round) {
		return [2]string{}
	}
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()

	indexRound := round - 1
	return m.RoundPlans[indexRound].UnionNames
}

func (m *UnionPlan) GetUnionnIcon(round int) [2]int {
	if !m.CheckRound(round) {
		return [2]int{}
	}
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()

	indexRound := round - 1
	return m.RoundPlans[indexRound].UnionIcons
}

func (m *UnionPlan) GetMasterName(round int) [2]string {
	if !m.CheckRound(round) {
		return [2]string{}
	}
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()

	indexRound := round - 1
	return m.RoundPlans[indexRound].UnionMasters
}

func (m *UnionPlan) GetFightRecord(round int, group int) []int64 {
	if !m.CheckParam(round, group) {
		return []int64{}
	}

	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	var fightIds []int64
	indexRound := round - 1
	indexGroup := group - 1
	for _, v := range m.RoundPlans[indexRound].FightRecord[indexGroup] {
		fightIds = append(fightIds, v.FightId)
	}
	return fightIds
}

// 添加玩家UID到分组
func (m *UnionPlan) AddFightRecord(round int, group int, fightId int64, player1 int64, player2 int64, result int, time int64) {
	if !m.CheckParam(round, group) {
		return
	}

	m.InfoLocker.Lock()
	found := false
	indexRound := round - 1
	indexGroup := group - 1
	for _, v := range m.RoundPlans[indexRound].FightRecord[indexGroup] {
		if v.FightId == fightId {
			found = true
			break
		}
	}

	if !found {
		m.RoundPlans[indexRound].FightRecord[indexGroup] = append(m.RoundPlans[indexRound].FightRecord[indexGroup], &UnionRecord{fightId, player1, player2, result, time})
	}
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) ClearPlan() {
	for i := 0; i < 3; i++ {
		pPlan := &RoundPlan{}
		for j := 0; j < 5; j++ {
			pPlan.GroupsA[j] = make([]int64, 0)
			pPlan.GroupsB[j] = make([]int64, 0)

			pPlan.DeadA[j] = make([]int64, 0)
			pPlan.DeadB[j] = make([]int64, 0)
			pPlan.FightRecord[j] = make([]*UnionRecord, 0)
		}
		m.RoundPlans[i] = pPlan
	}
	m.WinnerId = 0
	m.RoundNum = 0
	m.TheRound = 0
	m.TheGroup = 0
	m.InfoLocker = new(sync.RWMutex)
}

// 添加玩家UID到分组
func (m *UnionPlan) AddDeadAPlayers(round int, group int, UID int64) {
	if !m.CheckParam(round, group) {
		return
	}

	m.InfoLocker.Lock()
	found := false
	indexRound := round - 1
	indexGroup := group - 1
	for _, playerUID := range m.RoundPlans[indexRound].DeadA[indexGroup] {
		if playerUID == UID {
			found = true
			break
		}
	}

	if !found {
		m.RoundPlans[indexRound].DeadA[indexGroup] = append(m.RoundPlans[indexRound].DeadA[indexGroup], UID)
	}
	m.InfoLocker.Unlock()
}

// 添加所有没死玩家UID到分组
func (m *UnionPlan) AddAllDeadAPlayers(round int) {
	if !m.CheckRound(round) {
		return
	}

	m.InfoLocker.Lock()
	indexRound := round - 1
	deadA := m.RoundPlans[indexRound].DeadA
	players := m.RoundPlans[indexRound].GroupsA

	for i := 0; i < 5; i++ {
		for _, v := range players[i] {
			found := false
			for _, vd := range deadA[i] {
				if v == vd {
					found = true
					break
				}
			}

			if !found {
				m.RoundPlans[indexRound].DeadA[i] = append(m.RoundPlans[indexRound].DeadA[i], v)
			}
		}
	}
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) AddDeadBPlayers(round int, group int, UID int64) {
	if !m.CheckParam(round, group) {
		return
	}

	m.InfoLocker.Lock()
	found := false
	indexRound := round - 1
	indexGroup := group - 1
	for _, playerUID := range m.RoundPlans[indexRound].DeadB[indexGroup] {
		if playerUID == UID {
			found = true
			break
		}
	}

	if !found {
		m.RoundPlans[indexRound].DeadB[indexGroup] = append(m.RoundPlans[indexRound].DeadB[indexGroup], UID)
	}
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetRound(round int) {
	m.InfoLocker.Lock()
	defer m.InfoLocker.Unlock()
	m.RoundNum = round
}

func (m *UnionPlan) GetTheRound() int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	return m.TheRound
}

func (m *UnionPlan) GetTheGroup() int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	return m.TheGroup
}

func (m *UnionPlan) GetRoundNum() int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	return m.RoundNum
}

func (m *UnionPlan) SetTheRound(round int) {
	m.InfoLocker.Lock()
	defer m.InfoLocker.Unlock()
	m.TheRound = round
}

func (m *UnionPlan) SetTheGroup(group int) {
	m.InfoLocker.Lock()
	defer m.InfoLocker.Unlock()
	m.TheGroup = group
}

func (m *UnionPlan) GetRecords(round int) [5]int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckRound(round) {
		return [5]int{}
	}

	index := round - 1
	return m.RoundPlans[index].Records
}

func (m *UnionPlan) GetIndexGroup() int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	return m.TheGroup
}

func (m *UnionPlan) GetTeamMatch(round int) [2]int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckRound(round) {
		return [2]int{}
	}

	index := round - 1
	return m.RoundPlans[index].TeamMatch
}

func (m *UnionPlan) GetTeamNum(round int, group int, index int) int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckParam(round, group) {
		return 0
	}

	indexRound := round - 1
	indexGroup := group - 1
	if index == 0 {
		return len(m.RoundPlans[indexRound].GroupsA[indexGroup])
	} else if index == 1 {
		return len(m.RoundPlans[indexRound].GroupsB[indexGroup])
	}

	return 0
}

func (m *UnionPlan) GetPlayerIndexA(round int, UID int64) int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckRound(round) {
		return 0
	}
	index := round - 1
	pInfo := m.RoundPlans[index]
	for i := 0; i < 5; i++ {
		for _, v := range pInfo.GroupsA[i] {
			if v == UID {
				return i
			}
		}
	}
	return 0
}

func (m *UnionPlan) GetPlayerIndexB(round int, UID int64) int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckRound(round) {
		return 0
	}
	index := round - 1
	pInfo := m.RoundPlans[index]
	for i := 0; i < 5; i++ {
		for _, v := range pInfo.GroupsB[i] {
			if v == UID {
				return i
			}
		}
	}
	return 0
}

func (m *UnionPlan) GetTeamMatchs() [3][2]int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	var teamMatchs [3][2]int
	for i := range m.RoundPlans {
		for j, team := range m.RoundPlans[i].TeamMatch {
			teamMatchs[i][j] = team
		}
	}
	return teamMatchs
}

func (m *UnionPlan) GetTeamMatchRound(round int) [2]int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	roundIndex := round - 1
	return m.RoundPlans[roundIndex].TeamMatch
}

// round: [0,1,2]
// group: [0,1,2,3,4]
func (m *UnionPlan) CheckParam(round int, group int) bool {
	if round < 1 || round > 3 {
		LogError("CheckParam: round < 1 || round > 3, round:", round)
		return false
	}

	if group < 1 || group > 5 {
		LogError("CheckParam: group < 1 || group > 5, group:", group)
		return false
	}
	return true
}

func (m *UnionPlan) CheckRound(round int) bool {
	if round < 1 || round > 3 {
		LogError("CheckRound: round < 1 || round > 3, round:", round)
		return false
	}

	return true
}

func (m *UnionPlan) CheckGroup(group int) bool {
	if group < 1 || group > 5 {
		LogError("group < 1 || group > 5, group:", group)
		return false
	}

	return true
}

// 添加玩家UID到分组
func (m *UnionPlan) AddGroupAPlayers(round int, group int, UID int64) {
	if !m.CheckParam(round, group) {
		return
	}

	m.InfoLocker.Lock()
	found := false
	indexRound := round - 1
	indexGroup := group - 1
	for indexGroup := range m.RoundPlans[indexRound].GroupsA {
		for _, playerUID := range m.RoundPlans[indexRound].GroupsA[indexGroup] {
			if playerUID == UID {
				found = true
				break
			}
		}
	}

	if !found {
		m.RoundPlans[indexRound].GroupsA[indexGroup] = append(m.RoundPlans[indexRound].GroupsA[indexGroup], UID)
	}
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) AddGroupBPlayers(round int, group int, UID int64) {
	if !m.CheckParam(round, group) {
		return
	}

	m.InfoLocker.Lock()
	found := false
	indexRound := round - 1
	indexGroup := group - 1
	for indexGroup := range m.RoundPlans[indexRound].GroupsB {
		for _, playerUID := range m.RoundPlans[indexRound].GroupsB[indexGroup] {
			if playerUID == UID {
				found = true
				break
			}
		}
	}

	if !found {
		m.RoundPlans[indexRound].GroupsB[indexGroup] = append(m.RoundPlans[indexRound].GroupsB[indexGroup], UID)
	}
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetRecords(round, group, record int) {
	if !m.CheckParam(round, group) {
		return
	}

	indexRound := round - 1
	indexGroup := group - 1
	m.InfoLocker.Lock()
	m.RoundPlans[indexRound].Records[indexGroup] = record
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetTeamMatch(round, team1, team2 int) {
	if !m.CheckRound(round) {
		return
	}
	index := round - 1
	m.InfoLocker.Lock()
	m.RoundPlans[index].TeamMatch[0] = team1
	m.RoundPlans[index].TeamMatch[1] = team2
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetUnionMaster(round int, master1, master2 string) {
	if !m.CheckRound(round) {
		return
	}
	index := round - 1
	m.InfoLocker.Lock()
	m.RoundPlans[index].UnionMasters[0] = master1
	m.RoundPlans[index].UnionMasters[1] = master2
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetUnionName(round int, unionName1, unionName2 string) {
	if !m.CheckRound(round) {
		return
	}
	index := round - 1
	m.InfoLocker.Lock()
	m.RoundPlans[index].UnionNames[0] = unionName1
	m.RoundPlans[index].UnionNames[1] = unionName2
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetIndexFight(round int, index1, index2 int) {
	if !m.CheckRound(round) {
		return
	}
	index := round - 1
	m.InfoLocker.Lock()
	m.RoundPlans[index].IndexFight[0] = index1
	m.RoundPlans[index].IndexFight[1] = index2
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetIndexGroup(group int) {
	if !m.CheckGroup(group) {
		return
	}

	m.InfoLocker.Lock()
	defer m.InfoLocker.Unlock()
	m.TheGroup = group
}

func (m *UnionPlan) SetIndexCamp(round, camp1, camp2 int) {
	if !m.CheckRound(round) {
		return
	}
	index := round - 1
	m.InfoLocker.Lock()
	m.RoundPlans[index].Camp[0] = camp1
	m.RoundPlans[index].Camp[1] = camp2
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetRoundAndGroup(round, result int) {
	if !m.CheckRound(round) {
		return
	}
	m.InfoLocker.Lock()
	m.TheRound = round
	m.TheGroup = 1
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetResult(round, result int) {
	if !m.CheckRound(round) {
		return
	}
	index := round - 1
	m.InfoLocker.Lock()
	m.RoundPlans[index].Result = result
	m.InfoLocker.Unlock()
}

func (m *UnionPlan) SetWinnerID(winnerId int) {
	m.InfoLocker.Lock()
	m.WinnerId = winnerId
	m.InfoLocker.Unlock()
}

// 初始化军团安排以及战斗结果
func NewUnionPlan() *UnionPlan {
	plan := &UnionPlan{}
	for i := 0; i < 3; i++ {
		pPlan := &RoundPlan{}
		for j := 0; j < 5; j++ {
			pPlan.GroupsA[j] = make([]int64, 0)
			pPlan.GroupsB[j] = make([]int64, 0)

			pPlan.DeadA[j] = make([]int64, 0)
			pPlan.DeadB[j] = make([]int64, 0)
			pPlan.FightRecord[j] = make([]*UnionRecord, 0)
		}
		plan.RoundPlans[i] = pPlan
	}
	plan.InfoLocker = new(sync.RWMutex)
	return plan
}

func (m *UnionPlan) GetPeople(round int) [2]int {
	if !m.CheckParam(m.TheRound, m.TheGroup) {
		return [2]int{}
	}

	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()

	indexRound := round - 1
	num1 := 0
	for _, v := range m.RoundPlans[indexRound].GroupsA {
		num1 += len(v)
	}
	num2 := 0
	for _, v := range m.RoundPlans[indexRound].GroupsB {
		num2 += len(v)
	}

	return [2]int{num1, num2}
}

func (m *UnionPlan) GetPeopleGroup(round int, group int) [2]int {
	if !m.CheckParam(m.TheRound, m.TheGroup) {
		return [2]int{}
	}

	if group < 1 || group > 5 {
		return [2]int{}
	}

	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()

	indexRound := round - 1
	indexgroup := group - 1
	return [2]int{len(m.RoundPlans[indexRound].GroupsA[indexgroup]), len(m.RoundPlans[indexRound].GroupsB[indexgroup])}
}

func (m *UnionPlan) GetFightIndex(round int) [2]int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckParam(m.TheRound, m.TheGroup) {
		return [2]int{}
	}
	indexRound := round - 1
	return m.RoundPlans[indexRound].IndexFight
}

func (m *UnionPlan) GetFightTeam(round int) [2]int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckParam(m.TheRound, m.TheGroup) {
		return [2]int{}
	}
	indexRound := round - 1
	return m.RoundPlans[indexRound].TeamMatch
}

func (m *UnionPlan) GetFightResult(round int) [5]int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckParam(m.TheRound, m.TheGroup) {
		return [5]int{}
	}
	indexRound := round - 1
	return m.RoundPlans[indexRound].Records
}

func (m *UnionPlan) GetMyGroup(round int, UID int64) int {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckParam(round, m.TheGroup) {
		return 0
	}

	indexRound := round - 1
	for i := range m.RoundPlans[indexRound].GroupsA {
		for _, playerId := range m.RoundPlans[indexRound].GroupsA[i] {
			if playerId == UID {
				return i + 1
			}
		}
	}

	for i := range m.RoundPlans[indexRound].GroupsB {
		for _, playerId := range m.RoundPlans[indexRound].GroupsB[i] {
			if playerId == UID {
				return i + 1
			}
		}
	}

	return 0
}

func (m *UnionPlan) IsTeamGroupDead(round, group, side int) bool {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckParam(round, group) {
		return false
	}

	indexRound := round - 1
	indexGroup := group - 1

	if side == 1 {
		total := len(m.RoundPlans[indexRound].GroupsA[indexGroup])
		dead := len(m.RoundPlans[indexRound].DeadA[indexGroup])
		return dead >= total && total > 0
	} else if side == 2 {
		total := len(m.RoundPlans[indexRound].GroupsB[indexGroup])
		dead := len(m.RoundPlans[indexRound].DeadB[indexGroup])
		return dead >= total && total > 0
	}
	return true
}

func (m *UnionPlan) IsPlayerDead(round int, group int, side int, playerUID int64) bool {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	if !m.CheckParam(round, group) {
		return false
	}

	indexRound := round - 1
	indexGroup := group - 1

	if side == 1 {
		total := len(m.RoundPlans[indexRound].GroupsA[indexGroup])
		dead := len(m.RoundPlans[indexRound].DeadA[indexGroup])
		return total == dead
	} else if side == 2 {
		total := len(m.RoundPlans[indexRound].GroupsB[indexGroup])
		dead := len(m.RoundPlans[indexRound].DeadB[indexGroup])
		return total == dead
	}
	return true
}

func (m *UnionPlan) GetDeads(round, group, side int) map[int64]bool {
	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()
	var players = make(map[int64]bool)
	if !m.CheckParam(round, group) {
		return players
	}
	indexRound := round - 1
	pInfo := m.RoundPlans[indexRound]
	if side == 1 {
		for i := range pInfo.DeadA {
			for _, playerId := range pInfo.DeadA[i] {
				players[playerId] = true
			}
		}
	} else if side == 2 {
		for i := range pInfo.DeadB {
			for _, playerId := range pInfo.DeadB[i] {
				players[playerId] = true
			}
		}
	}

	return players
}

func (m *UnionPlan) GetAlive(round, group, side int) int64 {
	deadPlayers := m.GetDeads(round, group, side)
	if !m.CheckParam(round, group) {
		return 0
	}

	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()

	indexRound := round - 1
	pInfo := m.RoundPlans[indexRound]
	if side == 1 {
		for i := range pInfo.GroupsA {
			if i+1 != group {
				continue
			}
			for _, playerId := range pInfo.GroupsA[i] {
				_, ok := deadPlayers[playerId]
				if !ok {
					return playerId
				}
			}
		}
	} else if side == 2 {
		for i := range pInfo.GroupsB {
			if i+1 != group {
				continue
			}
			for _, playerId := range pInfo.GroupsB[i] {
				_, ok := deadPlayers[playerId]
				if !ok {
					return playerId
				}
			}
		}
	}
	return 0
}

// 结果全部都出来
// 如果攻击方num < 防守方，攻击方死亡
// 如果防守方num < 攻击方，防守方死亡

func (m *UnionPlan) IsTeamDead2(team int, round int, side int) bool {
	if round < 1 || round > 3 {
		return false
	}

	m.InfoLocker.RLock()
	defer m.InfoLocker.RUnlock()

	indexRound := round - 1
	records := m.RoundPlans[indexRound].Records

	num := 0
	winNum := 0
	defenceNum := 0
	for _, v := range records {
		if v > 0 {
			num += 1
		}
		if v == 1 {
			winNum += 1
		} else if v == 2 {
			defenceNum += 1
		}
	}

	if num < 5 {
		return false
	}

	//胜负相等的情况下比较额外条件 20190420 by zy
	if winNum == defenceNum {
		//比较2个军团的军团等级
		unionAUid := m.RoundPlans[indexRound].UnionID[0]
		unionBUid := m.RoundPlans[indexRound].UnionID[1]
		pUnionA := GetUnionMgr().GetUnion(unionAUid)
		pUnionB := GetUnionMgr().GetUnion(unionBUid)

		//先比较等级
		if pUnionA.Level != pUnionB.Level {
			if side == 1 {
				return pUnionA.Level < pUnionB.Level
			} else {
				return pUnionB.Level < pUnionA.Level
			}
		}
		//比较战斗力
		if side == 1 {
			return pUnionA.Fight < pUnionB.Fight
		} else {
			return pUnionB.Fight < pUnionA.Fight
		}
	}

	if side == 1 {
		return winNum <= defenceNum
	}

	return defenceNum < winNum
}
