package game

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	//"time"
)

const POS_ATTACK = 0  //! 攻击方
const POS_DEFENCE = 1 //! 防守方
const POS_ALIGN = 2   //! 盟友

const WIN_ATTCK = 1
const WIN_DEFENCE = 2
const WIN_NULL = 0
const WIN_DOGFALL = 3

const BASE_HONOR = 30

const FIGHT55_WAIT = 1
const FIGHT55_ENROLL = 2
const FIGHT55_FIGHT = 3
const FIGHT55_CD = 4

const FIGHT_SOLO = 10000
const FIGHT_SERIES = 11004
const FIGHT_COMBO = 11001
const FIGHT_LOSE_TIME = 30

//const FIGHT_55 = 12001
//var FIGHT_50 int = 12002

//! 战斗消息
type FightMsg struct {
	Head   string
	Param1 int
	Person *Player
}

//! 战报
type FightRecord struct {
	Id        int64              `json:"id"`        //! 战报ID
	Info      [2]Son_FightRecord `json:"info"`      //! 信息
	Result    int                `json:"result"`    //! -1单挑中, 0交战中, 1胜利 2失败
	FightId   int64              `json:"fightid"`   //! 战斗ID
	FightSeed int                `json:"fightseed"` //! 战斗随机种子
	IsSet     bool               `json:"isset"`     //! 是否战斗服务器
}

//! 战报
type Son_FightRecord struct {
	Uid     int64  `json:"uid"`     //! 角色ID
	Name    string `json:"name"`    //! 名字
	Icon    int    `json:"icon"`    //! 图标
	Level   int    `json:"level"`   //! 等级
	Fight   int64  `json:"fight"`   //! 战斗力
	Inspire int    `json:"inspire"` //! 鼓舞等级
	Kill    int    `json:"kill"`    //! 连斩次数
}

//! 排队节点
type FightNode struct {
	Uid        int64             //! uid >0玩家 -1城防军
	Name       string            //! 名字
	Icon       int               //! icon
	Level      int               //! 等级
	Fight      int64             //! 战斗力
	Index      int               //! 第几个队伍
	Hero       []int             //! 武将
	Beautyid   int               //! 美人
	Camp       int               //! 阵营
	Kill       int               //! 连斩-单挑
	KillFinal  int               //! 连斩-决战
	Class      int               //! 官职
	Honor      int               //! 累积功勋
	Honorstep  int               //! 本轮挑战功勋
	Buffer     int               //! 单挑buff层数
	Playlist   [6][]*JS_PlayNode //! 玩法记录
	InspireExp int               //! inspire exp
	InspireLv  int               //! inspire lv
	BuffList   []JS_Buffer       //! 个人Buff列表
	DropBuff   []*JS_Buffer      //! 掉落列表
	DropIndex  int               //! 选择掉落
	LoseTime   int64             //! 溃败时间
	CityPart   int               //! 参加据点
	SoloTime   int64             //! 上次参加Solo时间
	chg        bool              //! 修改了，决战同步
}

//! 上阵节点
type FightTeam struct {
	Node  *FightNode    //! 节点
	Team  *JS_FightInfo //! 详细信息
	Clone *JS_FightInfo //! 克隆信息
}

//! 单挑节点
type SoloFight struct {
	SoloId      int
	Solo        [2]*FightTeam
	Time        int64
	Result      int
	Fast        bool
	End         []FightHero
	CityPart    int         //! 单挑据点
	SelectIndex int         //! 攻击对象
	DropBuff    []JS_Buffer //! 掉落buff碎片
}

//! 连斩节点
type SeriesFightNode struct {
	Uid        int64 // uid
	WarPlayId  int   // 玩法Id
	Index      int   // 胜利场次
	MaxIndex   int   // 总场次
	Time       int64 // 开始时间
	TotalHonor int   // 累积获得Honor
	Result     int   // 战斗结果，0-战斗中，1-胜利，2-失败
	CityPart   int   //! 据点or主城
}

//! 玩法触发节点
type JS_PlayNode struct {
	PlayId   int          `json:"playid"`   //! 玩法
	Times    int          `json:"times"`    //! 剩余次数
	Max      int          `json:"max"`      //! 最大次数，0为不限制
	Open     int          `json:"open"`     //! 是否开放
	Pass     int          `json:"pass"`     //! 过关次数
	Show     int          `json:"show"`     //! 是否显示
	DropBuff []*JS_Buffer `json:"dropbuff"` //! 掉落Buff
}

//! Buffer结构
type JS_Buffer struct {
	BufferId  int   `json:"bufferid"`  // buffer id
	Time      int64 `json:"time"`      // 生效时间
	Overlying int   `json:"overlying"` // 叠加层数
	Effect    int   `json:"effect"`    // 持续时间
}

//! Top结构
type JS_FightTopNode struct {
	Uid   int64  `json:"uid"`   //! uid >0玩家 -1城防军
	Name  string `json:"name"`  //! 名字
	Icon  int    `json:"icon"`  //! icon
	Level int    `json:"level"` //! 等级
	Fight int64  `json:"fight"` //! 战斗力
	Camp  int    `json:"camp"`  //! 阵营
	Kill  int    `json:"kill"`  //! 连斩
	Honor int    `json:"honor"` //! 累积功勋
	Rank  int    `json:"rank"`  //! 名次
}

//! Top union 结构
type JS_UnionTopNode struct {
	Unionid int    `json:"unionid"` //! 军团Id
	Name    string `json:"name"`    //! 名字
	Level   int    `json:"level"`   //! 等级
	Fight   int64  `json:"fight"`   //! 战斗力
	Camp    int    `json:"camp"`    //! 阵营
	Honor   int    `json:"honor"`   //! 荣誉
	Rank    int    `json:"rank"`    //! 排行
	Iconid  int    `json:"iconid"`
}

//! 据点结构
type StrongholdNode struct {
	Index       int               //! 据点编号
	OccupyTotal int               //! 总占领值
	Occupy      [2]int            //! 占领值
	MaxOccupy   [2]int            //! 最大占领值
	Attack      [2]int            //! 攻打次数
	PlayerList  [2]map[int64]int  //! 人数列表
	Camp        int               //! 占领阵营 0-争夺中，1，2，3-蜀魏吴
	BufferList  [2][]JS_Buffer    //! buff叠加
	Record      []JS_OccupyRecord //! 占领值记录

	PlayerLocker *sync.RWMutex
}

type JS_OccupyRecord struct {
	Uid    int64  `json:"uid"`
	Uname  string `json:"uname"`
	Occupy int    `json:"occupy"`
	Id     int    `json:"id"`
}
type Fight struct {
	//! 队列战斗基础数据
	FightId         int64                   //! 当前战斗id
	FightTime       int64                   //! 当前的战斗结束时间-队列战斗
	StartTime       int64                   //! 国战开始时间
	Result          int                     //! 0未结果 1攻方胜利 2守方胜利
	Camp            [3]int                  //! 攻击方
	AlignForce      int                     //! 国战同盟势力
	Occupy          [2]int                  //! 占领值
	OccupyTotal     int                     //! 总占领值
	OccupyMax       [2]int                  //! 最高到达占领值
	Attack          [2]int                  //! 攻击次数
	WaitTeam        [2][]*FightTeam         //! 等待队伍0攻 1防
	FormatTeam      [2][]*FightTeam         //! 0攻 1防
	BufferInfo      [2][]JS_Buffer          //! Buffer
	OccupyRecord    [6][2][]JS_OccupyRecord //! 占领值记录
	FightResult     *FightResult            //! 当前的战斗结果
	IsOnce          bool                    //! 是否初始化
	FightRandom     int                     //! 战斗随机数
	PlayFight55     int                     //! 55玩法Id
	PlayFight56     int                     //! 56玩法Id
	LockPvP         bool                    //! 进入PvP锁定状态单人玩法不再计算占领值,PvP结束后计算是否接触状态
	SeriesKillAward bool                    //! 是否发送连斩奖励
	BattleEndAward  bool                    //! 是否发送胜负奖励
	OccupyRecordId  int                     //! 当前增加Id
	KingSkillTime   int64                   //! 国王技能冷却时间
	FinalRecord     int                     //! 决战记录
	FightSeq        int                     //! 决战顺序
	//--------------------------------------------------------------------------
	NpcMapId    [2]map[int][]int    //! 对应的 NPC 数据
	DefTeam     []*JS_FightInfo     //! 城防军数据
	AtkTeam     []*JS_FightInfo     //! 攻城军数据
	DefNum      int                 //! 城防军数量
	AtkNum      int                 //! 攻城军数量
	TeamInfo    []Son_CampFightInfo //! 队伍数据
	AttackPlay  []int               //! 进攻开启玩法
	DefensePlay []int               //! 防守开启玩法
	//--------------------------------------------------------------------------
	HelpIndex int           //! 援军自增长
	HelpTeam  map[int64]int //! 助攻数据
	//--------------------------------------------------------------------------
	//! 国战数据
	FightPlayer map[int64]*FightTeam //! 参与国战玩家数据
	QuitPlayer  map[int64]*FightNode //! 已退出玩家数据
	MapPlayer   map[int64]int        //! 观战玩家
	HonorLog    map[int64]int        //! 退出的Honor数量，下次直接增加
	PlayerLock  *sync.RWMutex        //! 观战玩家锁
	MsgChan     chan *FightMsg       //! 消息管道
	ChanLock    *sync.RWMutex        //! 管道锁

	//--------------------------------------------------------------------------
	//! 单挑玩法
	SoloPlayer map[int64]int //! 单挑玩家
	SoloTeam   []*SoloFight  //! 单挑队伍
	SoloIndex  int           //! 单挑自增长

	//--------------------------------------------------------------------------
	//!先锋强袭战-战斗

	//--------------------------------------------------------------------------
	//! 连斩战
	SeriesFightTeam map[int64]*SeriesFightNode //! 连斩节点

	//TeamLock *sync.RWMutex //! 参战队伍锁
	//SoloLock *sync.RWMutex //! 单挑队伍锁

	//! 国战据点
	StrongholdSet [5]*StrongholdNode //!国战据点

	//--------------------------------------------------------------------------
	//! 排行榜
	TopTeam     [2][]*JS_FightTopNode //! 前十排行榜
	Topver      [2]int                //! 排行榜版本
	TopUnion    [2][]*JS_UnionTopNode //! 军团排行榜
	TopverUnion [2]int                //! 军团排行榜版本
	RankTime    [2]int64              //! 排序时间-每3秒排序一次

	NpcIndex  map[int]int //! npcid
	RobotNum  int         //! 增加攻击机器人数量
	RobotTime int         //! 增加攻击机器人次数

	Event_Kill  map[int64]int //! 操作事件
	Event_Solo  map[int64]int //! 操作事件
	Event_Help  map[int64]int //! 操作事件
	Event_Honor map[int64]int //! 功勋记录

	Record     []FightRecord //! 战报
	RecordLock *sync.RWMutex
}

//! 增加玩法列表
func (self *FightNode) AddPlay(citypart int, playnode *JS_PlayNode) {
	if citypart < 0 || citypart > 5 {
		return
	}

	for i := 0; i < len(self.Playlist[citypart]); i++ {
		if self.Playlist[citypart][i].PlayId == playnode.PlayId {
			return
		}
	}

	self.Playlist[citypart] = append(self.Playlist[citypart], playnode)
}

//! 返回用户的玩法列表
func (self *FightNode) GetPlay(citypart int, playid int) *JS_PlayNode {
	for i := 0; i < len(self.Playlist[citypart]); i++ {
		if self.Playlist[citypart][i].PlayId == playid {
			return self.Playlist[citypart][i]
		}
	}

	return nil
}

func (self *FightNode) IsLosed() bool { //! 30秒溃败
	if self.LoseTime == 0 {
		return false
	}

	if TimeServer().Unix()-self.LoseTime < FIGHT_LOSE_TIME {
		return true
	}

	return false
}

func (self *FightNode) GetLoseTime() int {
	losetime := int(FIGHT_LOSE_TIME + self.LoseTime - TimeServer().Unix())
	if losetime < 0 || losetime > FIGHT_LOSE_TIME {
		losetime = 0
	}
	return losetime
}

func (self *FightNode) AddBuffer(bufferid int, overlying int) {
	find := false
	for i := 0; i < len(self.BuffList); i++ {
		if bufferid == self.BuffList[i].BufferId {
			self.BuffList[i].Overlying += overlying
			if warbuffcsv, ok := GetCsvMgr().Data["War_Buff"][bufferid]; ok {
				if self.BuffList[i].Overlying > HF_Atoi(warbuffcsv["synthetise"]) {
					self.BuffList[i].Overlying = HF_Atoi(warbuffcsv["synthetise"])
				}
			}

			find = true
		}
	}

	if find == false {
		var buffer JS_Buffer
		buffer.BufferId = bufferid
		buffer.Overlying = overlying

		self.BuffList = append(self.BuffList, buffer)
	}
}

func (self *FightNode) GetBufferLevel(bufferid int) int {
	for i := 0; i < len(self.BuffList); i++ {
		if self.BuffList[i].BufferId == bufferid {
			return self.BuffList[i].Overlying
		}
	}

	return 0
}

func (self *FightNode) RemoveBuffer(bufferid int, level int) {
	for i := 0; i < len(self.BuffList); i++ {
		if self.BuffList[i].BufferId == bufferid {
			if self.BuffList[i].Overlying >= level {
				self.BuffList[i].Overlying -= level
			} // else {
			//}
			//copy(self.BuffList[i:], self.BuffList[i+1:])
			//self.BuffList = self.BuffList[:len(self.BuffList)-1]
		}
	}
}

func (self *StrongholdNode) AddPlayer(pos int, uid int64) {
	if pos < POS_ATTACK || pos > POS_DEFENCE {
		return
	}

	self.PlayerLocker.Lock()
	defer self.PlayerLocker.Unlock()

	if _, ok := self.PlayerList[pos][uid]; !ok {
		self.PlayerList[pos][uid] = 1
	}
}

func (self *StrongholdNode) RemovePlayer(pos int, uid int64) {
	if pos < POS_ATTACK || pos > POS_DEFENCE {
		return
	}

	self.PlayerLocker.Lock()
	defer self.PlayerLocker.Unlock()

	delete(self.PlayerList[pos], uid)
}

func (self *StrongholdNode) GetPlayerNum() [2]int {
	self.PlayerLocker.RLock()
	defer self.PlayerLocker.RUnlock()

	return [2]int{len(self.PlayerList[0]), len(self.PlayerList[1])}
}

func (self *StrongholdNode) AddBuffer(pos int, bufferid int) {
	find := false
	for i := 0; i < len(self.BufferList[pos]); i++ {
		if bufferid == self.BufferList[pos][i].BufferId {
			self.BufferList[pos][i].Overlying++
			find = true
		}
	}

	if find == false {
		var buffer JS_Buffer
		buffer.BufferId = bufferid
		buffer.Overlying = 1

		self.BufferList[pos] = append(self.BufferList[pos], buffer)
	}
}

func (self *StrongholdNode) GetBufferLevel(pos int, bufferid int) int {
	for i := 0; i < len(self.BufferList[pos]); i++ {
		if self.BufferList[pos][i].BufferId == bufferid {
			return self.BufferList[pos][i].Overlying
		}
	}

	return 0
}

func (self *StrongholdNode) RemoveBuffer(pos int, bufferid int) {
	for i := 0; i < len(self.BufferList[pos]); i++ {
		if self.BufferList[pos][i].BufferId == bufferid {
			copy(self.BufferList[pos][i:], self.BufferList[pos][i+1:])
			self.BufferList[pos] = self.BufferList[pos][:len(self.BufferList[pos])-1]
		}
	}
}

func (self *StrongholdNode) AddAttack(pos int) {
	if pos < 0 || pos > POS_DEFENCE {
		return
	}
	self.Attack[pos]++
}

func (self *Fight) AddOccupyRecord(citypart int, pos int, uid int64, uname string, occupy int) {
	if citypart < 0 || citypart > 5 {
		return
	}

	if pos < POS_ATTACK || pos > POS_DEFENCE {
		return
	}

	self.OccupyRecord[citypart][pos] = append(self.OccupyRecord[citypart][pos],
		JS_OccupyRecord{uid, uname, occupy, self.OccupyRecordId})
	recordlen := len(self.OccupyRecord[citypart][pos])
	if recordlen > 10 {
		self.OccupyRecord[citypart][pos] = self.OccupyRecord[citypart][pos][recordlen-10:]
	}
	self.OccupyRecordId++
}

func (self *Fight) GetIndex(uid int) int {
	self.NpcIndex[uid]++

	return self.NpcIndex[uid]
}

//! 获得决战连斩数据
func (self *Fight) GetPlayerFinalKill(uid int64) int {
	self.PlayerLock.RLock()
	defer self.PlayerLock.RUnlock()

	if playerTeam, ok := self.FightPlayer[uid]; ok == true {
		return playerTeam.Node.KillFinal
	}

	return 0
}

//! 设置决战连斩-增加，重置，返回值为当前的值
func (self *Fight) AddPlayerFinalKill(uid int64, reset bool) int {
	self.PlayerLock.RLock()
	defer self.PlayerLock.RUnlock()

	if playerTeam, ok := self.FightPlayer[64]; ok == true {
		if reset == true {
			playerTeam.Node.KillFinal = 0
		} else {
			playerTeam.Node.KillFinal += 1
		}

		return playerTeam.Node.KillFinal
	}

	return 0
}

func (self *Fight) GetRecord() []FightRecord {
	self.RecordLock.RLock()
	defer self.RecordLock.RUnlock()

	return self.Record
}

func (self *Fight) AddMsg(msg *FightMsg) bool {
	self.ChanLock.Lock()
	defer self.ChanLock.Unlock()

	if self.MsgChan == nil {
		return false
	}

	//log.Println("addmsg:", msg.Head)
	LogDebug("addmsg:", msg.Head)
	self.MsgChan <- msg

	return true
}

//! 加助攻
func (self *Fight) AddHelp(uid int64) {
	if uid < -1000 { //! 玩家的援军
		uid += 1000
		uid = -uid
	}

	if uid <= 0 {
		return
	}

	_, ok := self.HelpTeam[uid]
	if ok {
		return
	}

	self.HelpTeam[uid] = 1
}

//! 结算
func (self *Fight) CountHelp() {
	for key := range self.HelpTeam {
		player := GetPlayerMgr().GetPlayer(key, true)
		if player == nil {
			continue
		}
		if player.Sql_UserBase.Camp == self.Camp[POS_DEFENCE] {
			continue
		}

		//player.GetModule("feats").(*Mod_Feats).AddHelpack()
		//player.HandleTask(101, 1, 1, self.Cityid)
	}
}

//! 机器人逻辑
func (self *Fight) OnRobot() {
	//return

	//self.City.moveLock.Lock()
	//defer self.City.moveLock.Unlock()

	//if self.Result != WIN_NULL {
	//	return
	//}

	//if len(self.MapPlayer) == 0 { //! 没有观战玩家，不增加机器人
	//	return
	//}

	//if self.RobotTime < 10 { //! 每10秒钟
	//	self.RobotTime++
	//	return
	//} else {
	//	self.RobotTime = 5
	//}

	//random := HF_GetRandom(self.RobotNum + 1)
	//if random != 0 {
	//	return
	//}

	self.RobotTime = 0
	self.RobotNum++

	fightteam := new(FightTeam)

	fightteam.Team = GetRobotMgr().GetDefNpc(GetCsvMgr().GetName(), 1, 0, self.Camp[POS_ATTACK])
	fightteam.Team.Rankid = FIGHTTYPE_ROBOT

	fightnode := new(FightNode)
	fightnode.Uid = fightteam.Team.Uid
	fightnode.Name = fightteam.Team.Uname
	fightnode.Icon = fightteam.Team.Iconid
	fightnode.Level = fightteam.Team.Level
	fightnode.Index = self.GetIndex(FIGHTTYPE_ROBOT)
	for i := 0; i < len(fightteam.Team.Defhero); i++ {
		fightnode.Hero = append(fightnode.Hero, fightteam.Team.Defhero[i])
	}
	fightnode.Fight = fightteam.Team.Deffight
	fightnode.Camp = self.Camp[POS_ATTACK]

	fightteam.Node = fightnode

	if len(self.FormatTeam[POS_ATTACK]) < 3 {
		self.FormatTeam[POS_ATTACK] = append(self.FormatTeam[POS_ATTACK], fightteam)

		var msg S2C_CampFightWait
		msg.Cid = "campfightwait"
		msg.Pos = POS_ATTACK
		msg.Wait = fightteam.Team
		self.BroadCastMsg("1", HF_JtoB(&msg))
	} else {
		self.WaitTeam[POS_ATTACK] = append(self.WaitTeam[POS_ATTACK], fightteam)
	}

	var info Son_CampFightInfo
	info.Uid = fightnode.Uid
	info.Index = fightnode.Index
	info.Name = fightnode.Name
	info.Icon = fightnode.Icon
	info.Fight = fightnode.Fight
	info.Level = fightnode.Level
	info.Camp = fightnode.Camp
	info.Kill = fightnode.Kill
	info.Elite = 0
	self.TeamInfo = append(self.TeamInfo, info)

	var msg S2C_CampFightAdd
	msg.Cid = "campfightadd"
	msg.Info = info
	self.BroadCastMsg("1", HF_JtoB(&msg))
}

func (self *Fight) IsFightNode(player *Player) bool {
	if self.Result != 0 {
		return false
	}

	_, ok := self.FightPlayer[player.Sql_UserBase.Uid]

	return ok
}

func (self *Fight) GetFightNode(player *Player) *FightTeam {
	if self.Result != 0 {
		return nil
	}

	playerteam, ok := self.FightPlayer[player.Sql_UserBase.Uid]
	if !ok {
		return nil
	}
	return playerteam
}

func (self *Fight) CanJoinFight(camp int) bool {
	if camp == 0 {
		return false
	}

	for i := 0; i < len(self.Camp); i++ {
		if self.Camp[i] == camp {
			return true
		}
	}

	return false
}

//! 广播消息
func (self *Fight) BroadCastMsg(head string, body []byte) {
	self.PlayerLock.RLock()
	defer self.PlayerLock.RUnlock()
	LogDebug("fight broadcast:", head, "....", string(body))
	var buffer bytes.Buffer
	buffer.Write(HF_DecodeMsg(head, body))
	for key := range self.MapPlayer {
		player := GetPlayerMgr().GetPlayer(key, false)
		if player != nil && player.SessionObj != nil {
			player.SessionObj.SendMsgBatch(buffer.Bytes())
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
type lstHonorTop []*JS_FightTopNode

func (s lstHonorTop) Len() int           { return len(s) }
func (s lstHonorTop) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstHonorTop) Less(i, j int) bool { return s[i].Honor > s[j].Honor }

//军团功勋排行榜
////////////////////////////////////////////////////////////////////////////////
type lstUnionHonorTop []*JS_UnionTopNode

func (s lstUnionHonorTop) Len() int           { return len(s) }
func (s lstUnionHonorTop) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstUnionHonorTop) Less(i, j int) bool { return s[i].Honor > s[j].Honor }

type mapFight map[int]*Fight

// 战报信息
type FightResult struct {
	Id           int64            // 战斗id
	Result       int              // 结果[0无结果 1攻击方胜利 2防守方胜利]
	Fight        [2]*JS_FightInfo // 交战双方数据
	Info         [2][]FightHero   // 战斗结果数据
	Random       int              // 随机数
	Time         int              // 交战时间
	ResultDetail *FightResultNode // 战斗结果详情
	CityId       int              // 战斗发起城池/或者矿点Id
	SecKill      int              // 是否秒杀, 0 不是秒杀 1是秒杀 [战力悬殊过大，直接秒杀]
	TeamA        int              // 军团战队伍1
	TeamB        int              // 军团战队伍2
	Group        int              // 分组
	IsSet        bool             // 是否通过SET拿到的战报结果
}

type FightMgr struct {
	MapFight mapFight
	Locker   *sync.RWMutex

	FightId      int64                  //! 战斗id
	MapAutoFight map[int64]*FightResult //! 自动战斗
	FightLock    *sync.RWMutex
}

var fightmgrsingleton *FightMgr = nil

//! public
func GetFightMgr() *FightMgr {
	if fightmgrsingleton == nil {
		fightmgrsingleton = new(FightMgr)
		fightmgrsingleton.Locker = new(sync.RWMutex)
		fightmgrsingleton.MapFight = make(mapFight)
		fightmgrsingleton.FightLock = new(sync.RWMutex)
		fightmgrsingleton.MapAutoFight = make(map[int64]*FightResult)
	}

	return fightmgrsingleton
}

func (self *FightMgr) GetData() {
	res, err := GetAllRedis("san_fightresult")
	if res == nil {
		return
	}

	if err != nil {
		LogError(err.Error())
		return
	}

	self.FightLock.Lock()
	for _, v := range res {
		var data FightResult
		err = json.Unmarshal([]byte(v), &data)

		if data.Fight[0] == nil || data.Fight[1] == nil {
			continue
		}
		if len(data.Fight[0].Heroinfo) == 0 || len(data.Fight[1].Heroinfo) == 0 {
			continue
		}

		if err == nil {
			if self.FightId < data.Id {
				self.FightId = data.Id
			}
			if data.Id > 0 {
				self.MapAutoFight[data.Id] = &data
			}
		}
	}
	self.FightLock.Unlock()
}

func (self *FightMgr) GetFight(id int) *Fight {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	fight, ok := self.MapFight[id]
	if ok {
		return fight
	} else {
		return nil
	}
}

// 这里生成一个战斗序列, 提供给战斗服务器进行get请求处理
func (self *FightMgr) AddFightID(attack *JS_FightInfo, defence *JS_FightInfo, random int, cityid int) *FightResult {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	self.FightId++
	record := (TimeServer().Unix()%10000000)*100 + self.FightId%100

	client := new(FightResult)
	client.Id = record
	client.Fight[0] = attack
	client.Fight[1] = defence
	client.Random = random
	client.CityId = cityid
	self.MapAutoFight[record] = client

	//国战战报存储调试
	LogDebug("战报ID：", client.Id, "--数据：", HF_JtoA(client))

	//HMSetRedis("san_fightresult", client.Id, client, DAY_SECS*10)
	return client
}

func (self *FightMgr) AddUnionFightID(attack *JS_FightInfo, defence *JS_FightInfo, random int, cityid int) *FightResult {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	self.FightId++
	record := (TimeServer().Unix()%10000000)*100 + self.FightId%100
	LogDebug("fightRecordId = ", record)

	client := new(FightResult)
	client.Id = record
	client.Fight[0] = attack
	client.Fight[1] = defence
	client.Random = random
	client.CityId = cityid
	self.MapAutoFight[record] = client

	//HMSetRedis("san_fightresult", client.Id, client, DAY_SECS*10)
	return client
}

func (self *FightMgr) GetFightID(serverid int) int64 {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	self.FightId++
	record := (TimeServer().Unix()%100000)*100 + 100000000*int64(serverid) + self.FightId%100
	return record
}

func (self *FightMgr) GetFightInfoID() int64 {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	self.FightId++
	record := (TimeServer().Unix()%1000000)*1000 + self.FightId%1000
	return record
}

func (self *FightMgr) GetResult(fightid int64) *FightResult {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	result, ok := self.MapAutoFight[fightid]
	if !ok || result.Result == WIN_NULL { //! 没有结果
		return nil
	}
	return result
}

func (self *FightMgr) GetFightClient(fightid int64) *FightResult {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	result, ok := self.MapAutoFight[fightid]
	if !ok { //! 没有结果
		return nil
	} else {
		return result
	}
}

func (self *FightMgr) DelResult(fightid int64) {
	self.FightLock.Lock()
	delete(self.MapAutoFight, fightid)
	self.FightLock.Unlock()
}

//定期清除战报 20190808 by zy
func (self *FightMgr) DelResultExceed() {
	self.FightLock.Lock()

	nowTime := TimeServer().Unix()
	for key, value := range self.MapAutoFight {
		if int64(value.Time) < nowTime-10*DAY_SECS {
			delete(self.MapAutoFight, key)
		}
	}
	self.FightLock.Unlock()
}

func (self *FightMgr) GetTopFight() *FightResult {
	self.FightLock.RLock()
	defer self.FightLock.RUnlock()

	for _, value := range self.MapAutoFight {
		if value.Result != WIN_NULL {
			continue
		}
		return value
	}

	return nil
}

// 战报存储数据库
func (self *FightMgr) SetResult(node *FightResultNode) {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	result, ok := self.MapAutoFight[node.Fightid]
	if !ok {
		return
	}

	result.Info = node.Info
	result.Time = int(TimeServer().Unix())
	result.IsSet = true //标记是通过战斗服务器拿到的结果
	if node.Winner == 0 {
		result.Result = WIN_ATTCK
	} else {
		result.Result = WIN_DEFENCE
	}
	result.ResultDetail = node
	LogDebug("战斗结果改变:", result.Id, ",", result.Result, node)

	//写入数据库
	self.MapAutoFight[node.Fightid] = result
	//HMSetRedis("san_fightresult", result.Id, result, DAY_SECS*10)

	//attack_player := GetPlayerMgr().GetPlayer(result.Fight[0].Uid, false)
	//defence_player := GetPlayerMgr().GetPlayer(result.Fight[1].Uid, false)
	//if attack_player != nil && defence_player != nil {
	//	if result.Result == WIN_ATTCK {
	//		GetServer().SqlLog(attack_player.GetUid(), LOG_MINE_FIGHT_FINISH, LOG_SUCCESS, int(defence_player.Sql_UserBase.Uid),
	//			int(defence_player.Sql_UserBase.Fight/100), "矿点争夺战斗", 0, result.CityId, attack_player)
	//
	//		GetServer().SqlLog(defence_player.GetUid(), LOG_MINE_FIGHT_FINISH, LOG_FAIL, int(attack_player.Sql_UserBase.Uid),
	//			int(attack_player.Sql_UserBase.Fight/100), "矿点争夺战斗", 0, result.CityId, attack_player)
	//	} else {
	//		GetServer().SqlLog(attack_player.GetUid(), LOG_MINE_FIGHT_FINISH, LOG_FAIL, int(defence_player.Sql_UserBase.Uid),
	//			int(defence_player.Sql_UserBase.Fight/100), "矿点争夺战斗", 0, result.CityId, attack_player)
	//
	//		GetServer().SqlLog(defence_player.GetUid(), LOG_MINE_FIGHT_FINISH, LOG_SUCCESS, int(attack_player.Sql_UserBase.Uid),
	//			int(attack_player.Sql_UserBase.Fight/100), "矿点争夺战斗", 0, result.CityId, attack_player)
	//	}
	//}
}

// 加入战报
func (self *FightMgr) AddResult(result *FightResult) {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	_, ok := self.MapAutoFight[result.Id]
	if ok {
		return
	}

	self.MapAutoFight[result.Id] = result
	LogDebug("战斗结果改变:", result.Id, ",", result.Result)

	//写入数据库
	//HMSetRedis("san_fightresult", result.Id, result, DAY_SECS*10)
}

type FightServerNode struct {
	Id     int64         `json:"id"`
	Att    *JS_FightInfo `json:"att"`
	Def    *JS_FightInfo `json:"def"`
	Random int           `json:"random"`
}

type FightResultNode struct {
	Fightid int64          `json:"fightid"`
	Info    [2][]FightHero `json:"info"`
	Time    int            `json:"time"`
	Winner  int            `json:"winner"`
}

func FightServer(w http.ResponseWriter, r *http.Request) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()
	msgtype := r.FormValue("msgtype")
	if msgtype == "get" {
		client := GetFightMgr().GetTopFight()
		if client == nil {
			//log.Println("无战斗数据")
			w.Write([]byte("false"))

			//var node FightServerNode
			//configFile, err := ioutil.ReadFile("./test.json") ///尝试打开配置文件
			//if err != nil {
			//	log.Fatal("config err 1")
			//}
			//var fightContent FightResult
			//err = json.Unmarshal(configFile, &fightContent)
			//if err != nil {
			//	log.Fatal("server InitConfig err:", err.Error())
			//}
			//
			//node.Id = fightContent.Id
			//node.Random = fightContent.Random
			//node.Att = fightContent.Fight[0]
			//node.Def = fightContent.Fight[1]
			//
			//w.Write(HF_JtoB(&node))

		} else {
			LogDebug("有战斗数据")
			var node FightServerNode
			node.Id = client.Id
			node.Att = client.Fight[0]
			node.Def = client.Fight[1]
			node.Random = client.Random
			LogDebug("有战斗数据", *node.Att, *node.Def)
			w.Write(HF_JtoB(&node))
		}

		return
	} else if msgtype == "set" {
		data := r.FormValue("data")
		LogDebug(data)
		var node FightResultNode
		err := json.Unmarshal([]byte(data), &node)
		if err == nil {
			GetFightMgr().SetResult(&node)
		}
		client := GetFightMgr().GetTopFight()
		if client == nil {
			//log.Println("无战斗数据")
			w.Write([]byte("false"))
		} else {
			LogDebug("有战斗数据", client.Id)
			var node FightServerNode
			node.Id = client.Id
			node.Att = client.Fight[0]
			node.Def = client.Fight[1]
			node.Random = client.Random
			w.Write(HF_JtoB(&node))
		}
		return
	}
}

type FightInfo struct {
	FightId   int64  `json:"fight_id"`   // 战报Id
	Result    int    `json:"result"`     // 战斗方式
	EnemyName string `json:"enemy_name"` // 敌方
	Icon      int    `json:"icon"`       // 敌方Icon
	Fight     int64  `json:"fight"`      // 敌方战力
	Time      int    `json:"time"`       // 时间
	Side      int    `json:"side"`       // 1 进攻方  2 防守方
	SecKill   int    `json:"sec_kill"`   // 0 不是秒杀 1 秒杀
	TeamA     int    `json:"team_a"`     // 军团战队伍1
	TeamB     int    `json:"team_b"`     // 军团战队伍2
	Group     int    `json:"group"`      // 分组
	IsSet     bool   `json:"isset"`      // 是否通过战斗服务器
}

// 通过战报Id获取战报
func (self *FightMgr) GetFightResults(fightIds []int64) []*FightResult {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()
	var res []*FightResult
	for _, fightId := range fightIds {
		data, ok := self.MapAutoFight[fightId]
		if !ok {
			continue
		}
		res = append(res, data)
	}
	return res
}

func (self *FightMgr) GetFightResults2(fightIds []int64) map[int64]*FightResult {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()
	var res = make(map[int64]*FightResult)
	for _, fightId := range fightIds {
		data, ok := self.MapAutoFight[fightId]
		if !ok {
			continue
		}
		res[fightId] = data
	}
	return res
}

// 通过战报Id获取战报
func (self *FightMgr) GetFightResult(fightId int64) *FightResult {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()
	data, ok := self.MapAutoFight[fightId]
	if !ok {
		return nil
	}
	return data
}

func (self *FightMgr) GetFightInfos(results []*FightResult, uid int64) []*FightInfo {
	var fightInfos []*FightInfo
	for _, result := range results {
		fightInfos = append(fightInfos, NewFightInfo(result, uid))
	}
	return fightInfos
}

func (self *FightMgr) GetFightInfo(result *FightResult, uid int64) *FightInfo {
	return NewFightInfo(result, uid)
}

func NewFightInfo(result *FightResult, uid int64) *FightInfo {
	fightInfo := &FightInfo{}
	fightInfo.FightId = result.Id
	fightInfo.Result = result.Result
	fightInfo.Time = result.Time
	//fightInfo.Time = int(TimeServer().Unix())
	fightInfo.SecKill = result.SecKill
	attacker := result.Fight[0]
	defercer := result.Fight[1]
	if attacker != nil && attacker.Uid == uid && defercer != nil {
		fightInfo.Icon = defercer.Iconid
		fightInfo.Side = 1
		fightInfo.EnemyName = defercer.Uname
		fightInfo.Fight = defercer.Deffight

	} else if defercer != nil && defercer.Uid == uid && attacker != nil {
		fightInfo.Icon = attacker.Iconid
		fightInfo.Side = 2
		fightInfo.EnemyName = attacker.Uname
		fightInfo.Fight = attacker.Deffight
	}
	fightInfo.IsSet = result.IsSet
	return fightInfo
}

// 加入战报
func (self *FightMgr) AddMineResult(result *FightResult) {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	_, ok := self.MapAutoFight[result.Id]
	if ok {
		return
	}

	result.Time = int(TimeServer().Unix())
	self.MapAutoFight[result.Id] = result
	LogDebug("战斗结果改变:", result.Id, ",", result.Result)

	//写入数据库
	//HMSetRedis("san_fightresult", result.Id, result, DAY_SECS*1)
}

// 加入军团战战报
func (self *FightMgr) AddUnionResult(result *FightResult, isOutTime bool) {
	self.FightLock.Lock()
	//self.FightLock.Unlock()

	_, ok := self.MapAutoFight[result.Id]
	if ok {
		self.FightLock.Unlock()
		return
	}

	self.MapAutoFight[result.Id] = result
	LogDebug("战斗结果改变:", result.Id, ",", result.Result)
	self.FightLock.Unlock()
	//写入数据库
	//HMSetRedis("san_fightresult", result.Id, result, DAY_SECS*10)

	//20190727 by zy 虽然写了输赢但战报肯定是没有，需要去除掉，玩家也必然读不到战报
	//这里也个错误LOG，表明是战斗服务器超时造成战报丢失
	if isOutTime {
		LogError("战斗服务器超时造成战报丢失", result.Id)
		GetFightMgr().DelResult(result.Id)
	}
}

// 加入军团战战报
func (self *FightMgr) AddArenaFightID(attack *JS_FightInfo, defend *JS_FightInfo, random int64) *FightResult {
	self.FightLock.Lock()
	defer self.FightLock.Unlock()

	self.FightId++
	record := (TimeServer().Unix()%10000000)*100 + self.FightId%100

	client := new(FightResult)
	client.Id = record
	client.Fight[0] = attack
	client.Fight[1] = defend
	client.Random = int(random)
	self.MapAutoFight[record] = client

	//国战战报存储调试
	LogDebug("战报ID：", client.Id, "--数据：", HF_JtoA(client))

	//HMSetRedis("san_fightresult", client.Id, client, DAY_SECS*10)
	return client
}
