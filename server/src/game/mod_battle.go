package game

// 战斗简报战报回放相关模块
import "encoding/json"

const (
	BATTLE_TYPE_PVP  = 1 //
	BATTLE_TYPE_PVE  = 2 //
	BATTLE_TYPE_BOSS = 3 //
)

const (
	BATTLE_TYPE_TOWER                 = 1  // 爬塔
	BATTLE_TYPE_ARENA                 = 2  // 普通竞技场
	BATTLE_TYPE_UNION_HUNT_NORMAl     = 3  // 公会狩猎1
	BATTLE_TYPE_UNION_HUNT_SPECIAL    = 4  // 公会狩猎2
	BATTLE_TYPE_ARENA_SPECIAL         = 5  // 高阶竞技场
	BATTLE_TYPE_NORMAL                = 6  // 日常副本
	BATTLE_TYPE_RECORD_BOSS           = 7  // 暗域入侵
	BATTLE_TYPE_RECORD_CROSSARENA     = 8  // 跨服竞技战报
	BATTLE_TYPE_RECORD_BOSS_FESTIVAL  = 9  // 节日BOSS
	BATTLE_TYPE_RECORD_CROSSARENA_3V3 = 10 // 跨服竞技33战报
)

const (
	MSG_GET_BATTLE_INFO   = "get_battle_info"
	MSG_GET_BATTLE_INFO2  = "get_battle_info2"
	MSG_GET_BATTLE_INFO3  = "get_battle_info3"
	MSG_GET_BATTLE_RECORD = "get_battle_record"
)

type BattleSkill struct {
	Time  int64 `json:"time"`  // 使用时间
	Skill int   `json:"skill"` // 使用技能
}

//! 战报
type BattleRecord struct {
	Id        int64            `json:"id"`            //! 战报ID
	Type      int              `json:"type"`          //! 战报类型
	Side      int              `json:"side"`          // 自己是 1 进攻方 0 防守方
	Result    int              `json:"attack_result"` // 0成功 其他失败
	Level     int              `json:"level"`         // 之前名次
	LevelID   int              `json:"levelid"`       //! 关卡id
	Time      int64            `json:"time"`          // 发生的时间
	RandNum   int64            `json:"rand_num"`      // 随机数
	Weaken    []*WeakenInfo    `json:"weaken"`        // 压制
	FightInfo [2]*JS_FightInfo `json:"fight_info"`    // 双方数据 一个是进攻 第二个是防守
}

type ArmyInfo struct {
	Uid      int64       `json:"uid"`      //! uid
	Uname    string      `json:"uname"`    //! 名字
	Iconid   int         `json:"iconid"`   //! 头像
	SelfKey  int         `json:"selfkey"`  //! 对应key
	Face     int         `json:"face"`     //! 李四
	Portrait int         `json:"portrait"` //! 头像框
	Pos      int         `json:"pos"`      //!
	Atts     []ArmyAttr  `json:"atts"`     //! 佣兵属性
	Skills   []ArmySkill `json:"skills"`   //! 佣兵技能
	Lv       int
}

type ArmyAttr struct {
	Type  int   `json:"t"` //!
	Value int64 `json:"v"` //!
}

type ArmySkill struct {
	Id    int `json:"id"`    //!
	Level int `json:"level"` //!
}

type BattleHeroInfo struct {
	HeroID      int       `json:"heroid"`     //! 英雄id
	HeroLv      int       `json:"herolv"`     //! 英雄等级
	HeroStar    int       `json:"herostar"`   //! 英雄星级
	HeroSkin    int       `json:"skin"`       //! 英雄皮肤
	Hp          int       `json:"hp"`         // hp
	Energy      int       `json:"rage"`       // 怒气
	Damage      int64     `json:"damage"`     //! 伤害
	TakeDamage  int64     `json:"takedamage"` //! 承受伤害
	Healing     int64     `json:"healing"`    //! 治疗
	ArmyInfo    *ArmyInfo `json:"ownplayer"`
	ExclusiveLv int       `json:"exclusivelv"` //! 专属等级
	UseSkill    []int     `json:"skilltime"`   //! 使用的技能 pve专用  竞技场不使用
}

type BattleUserInfo struct {
	Uid       int64             `json:"uid"`        //! 名字
	Name      string            `json:"name"`       //! 名字
	Icon      int               `json:"icon"`       //! 头像
	Portrait  int               `json:"portrait"`   // 头像框
	UnionName string            `json:"union_name"` //! 军团名字
	Level     int               `json:"level"`      //! 等级
	HeroInfo  []*BattleHeroInfo `json:"heroinfo"`   // 双方数据
}

type WeakenInfo struct {
	Att   int     `json:"att"`   //!
	Value float64 `json:"value"` //!
}

type BattleInfo struct {
	Id       int64              `json:"id"`       //! 战报ID
	LevelID  int                `json:"levelid"`  //! 关卡id
	Type     int                `json:"type"`     //! 战报类型
	Time     int64              `json:"time"`     // 发生的时间
	Result   int                `json:"result"`   // 结果
	Random   int64              `json:"random"`   // 随机数
	UserInfo [2]*BattleUserInfo `json:"userinfo"` // 己方玩家数据
	Weaken   []*WeakenInfo      `json:"weaken"`   // 压制
}

// 战斗回放
type ModBattle struct {
	player *Player
}

func (self *ModBattle) OnGetData(player *Player) {
	self.player = player
}

func (self *ModBattle) OnGetOtherData() {
}

func (self *ModBattle) OnSave(sql bool) {
}

// 老的消息处理
func (self *ModBattle) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (self *ModBattle) onReg(handlers map[string]func(body []byte)) {
	handlers[MSG_GET_BATTLE_INFO] = self.GetBattleInfo
	handlers[MSG_GET_BATTLE_INFO2] = self.GetBattleInfo2
	handlers[MSG_GET_BATTLE_INFO3] = self.GetBattleInfo3
	handlers[MSG_GET_BATTLE_RECORD] = self.GetBattleRecord
}

func (self *ModBattle) GetBattleInfo(body []byte) {
	var msg C2S_GetBattleInfo
	json.Unmarshal(body, &msg)

	var battleInfo *BattleInfo = nil
	switch msg.Type {
	case BATTLE_TYPE_TOWER:
		battleInfo = self.player.GetModule("tower").(*ModTower).GetBattleInfo(msg.FightID)
	case BATTLE_TYPE_ARENA:
		battleInfo = self.player.GetModule("arena").(*ModArena).GetBattleInfo(msg.FightID)
	case BATTLE_TYPE_UNION_HUNT_NORMAl:
		battleInfo = self.player.GetModule("union").(*ModUnion).GetBattleInfo(msg.FightID, UNION_HUNT_TYPE_NOMAL)
	case BATTLE_TYPE_UNION_HUNT_SPECIAL:
		battleInfo = self.player.GetModule("union").(*ModUnion).GetBattleInfo(msg.FightID, UNION_HUNT_TYPE_SPECIAL)
	case BATTLE_TYPE_NORMAL:
		battleInfo = self.player.GetModule("pass").(*ModPass).GetBattleInfo(msg.FightID)
	case BATTLE_TYPE_RECORD_BOSS, BATTLE_TYPE_RECORD_BOSS_FESTIVAL:
		battleInfo = self.player.GetModule("activityboss").(*ModActivityBoss).GetBattleInfo(msg.FightID)
	case BATTLE_TYPE_RECORD_CROSSARENA, BATTLE_TYPE_RECORD_CROSSARENA_3V3:
		battleInfo = self.player.GetModule("crossarena").(*ModCrossArena).GetBattleInfo(msg.FightID)
	}
	if battleInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("战报录像版本过低，无法播放"))
		return
	}

	var backmsg S2C_GetBattleInfo
	backmsg.Cid = MSG_GET_BATTLE_INFO
	backmsg.Uid = msg.Uid
	backmsg.Type = msg.Type
	backmsg.FightID = msg.FightID
	backmsg.BattleInfo = battleInfo
	self.player.Send(backmsg.Cid, backmsg)
}
func (self *ModBattle) GetBattleInfo2(body []byte) {
	var msg C2S_GetBattleInfo2
	json.Unmarshal(body, &msg)

	battleInfo := self.player.GetModule("arenaspecial").(*ModArenaSpecial).GetBattleInfo(msg.FightID)

	if battleInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIGHT_PLAYBACK_ERRER"))
		return
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_ARENA_SPECIAL_BATTLE_INFO, int(msg.FightID), 0, 0, "查看高阶竞技场战报", 0, 0, self.player)

	var backmsg S2C_GetBattleInfo2
	backmsg.Cid = MSG_GET_BATTLE_INFO2
	backmsg.FightID = msg.FightID
	backmsg.BattleInfo = battleInfo
	self.player.Send(backmsg.Cid, backmsg)
}

func (self *ModBattle) GetBattleInfo3(body []byte) {
	var msg C2S_GetBattleInfo2
	json.Unmarshal(body, &msg)

	battleInfo := self.player.GetModule("crossarena3v3").(*ModCrossArena3V3).GetBattleInfo(msg.FightID)

	if battleInfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIGHT_PLAYBACK_ERRER"))
		return
	}

	var backmsg S2C_GetBattleInfo2
	backmsg.Cid = MSG_GET_BATTLE_INFO3
	backmsg.FightID = msg.FightID
	backmsg.BattleInfo = battleInfo
	self.player.Send(backmsg.Cid, backmsg)
}

func (self *ModBattle) GetBattleRecord(body []byte) {
	var msg C2S_GetBattleRecord
	json.Unmarshal(body, &msg)

	var battleRecord *BattleRecord = nil
	switch msg.Type {
	case BATTLE_TYPE_TOWER:
		battleRecord = self.player.GetModule("tower").(*ModTower).GetBattleRecord(msg.FightID)
	case BATTLE_TYPE_ARENA:
		battleRecord = self.player.GetModule("arena").(*ModArena).GetBattleRecord(msg.FightID)
	case BATTLE_TYPE_UNION_HUNT_NORMAl:
		battleRecord = self.player.GetModule("union").(*ModUnion).GetBattleRecord(msg.FightID, UNION_HUNT_TYPE_NOMAL)
	case BATTLE_TYPE_UNION_HUNT_SPECIAL:
		battleRecord = self.player.GetModule("union").(*ModUnion).GetBattleRecord(msg.FightID, UNION_HUNT_TYPE_SPECIAL)
	case BATTLE_TYPE_ARENA_SPECIAL:
		battleRecord = self.player.GetModule("arenaspecial").(*ModArenaSpecial).GetBattleRecord(msg.FightID)
	case BATTLE_TYPE_NORMAL:
		battleRecord = self.player.GetModule("pass").(*ModPass).GetBattleRecord(msg.FightID)
	case BATTLE_TYPE_RECORD_BOSS, BATTLE_TYPE_RECORD_BOSS_FESTIVAL:
		battleRecord = self.player.GetModule("activityboss").(*ModActivityBoss).GetBattleRecord(msg.FightID)
	case BATTLE_TYPE_RECORD_CROSSARENA, BATTLE_TYPE_RECORD_CROSSARENA_3V3:
		battleRecord = self.player.GetModule("crossarena").(*ModCrossArena).GetBattleRecord(msg.FightID)
	}

	var backmsg S2C_GetBattleRecord
	backmsg.Cid = MSG_GET_BATTLE_RECORD
	backmsg.Uid = msg.Uid
	backmsg.Type = msg.Type
	backmsg.FightID = msg.FightID
	backmsg.BattleRecord = battleRecord
	self.player.Send(backmsg.Cid, &backmsg)
}
