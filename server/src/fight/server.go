package main

import (
	"encoding/json"
	"fmt"
	"game"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	INDEX_ALL    = 0
	INDEX_MINUTE = 1
	INDEX_HOUR   = 2
	INDEX_DAY    = 3
	INDEX_MAX    = 4
)

type FightHero struct {
	Heroid int `json:"id"`   //! 武将id
	Hp     int `json:"hp"`   //! hp
	Energy int `json:"rage"` //! 怒气
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

type JS_HeroInfo struct {
	Heroid          int              `json:"heroid"`    // 武将id
	HeroKeyId       int              `json:"herokeyid"` // 武将keyid
	Color           int              `json:"color"`     //
	Stars           int              `json:"stars"`
	Levels          int              `json:"levels"`
	Soldiercolor    int              `json:"soldiercolor"` // 士兵星级
	Soldierid       int              `json:"soldierid"`    // 士兵id
	Skilllevel1     int              `json:"skilllevel1"`  // 技能等级
	Skilllevel2     int              `json:"skilllevel2"`
	Skilllevel3     int              `json:"skilllevel3"`
	Skilllevel4     int              `json:"skilllevel4"`
	Skilllevel5     int              `json:"skilllevel5"`
	Skilllevel6     int              `json:"skilllevel6"`
	Fervor1         int              `json:"fervor1"` // 战魂
	Fervor2         int              `json:"fervor2"`
	Fervor3         int              `json:"fervor3"`
	Fervor4         int              `json:"fervor4"`
	Fight           int64            `json:"fight"` // 战斗力
	ArmsSkill       []JS_ArmsSkill   `json:"armsskill"`
	TalentSkill     []Js_TalentSkill `json:"talentskill"` // 天赋技能
	ArmyId          int              `json:"army_id"`
	MainTalent      int              `json:"maintalent"`      // 英雄主天赋等级
	HeroQuality     int              `json:"heroquality"`     //英雄品质  (不变的属性)
	HeroArtifactId  int              `json:"heroartifactid"`  //英雄神器等级 (不变的属性)
	HeroArtifactLv  int              `json:"heroartifactlv"`  //英雄神器等级 (不变的属性)
	HeroExclusiveLv int              `json:"heroexclusivelv"` //英雄专属等级 (不变的属性)
	Skin            int              `json:"skin"`            //皮肤
	ExclusiveUnLock int              `json:"exclusiveunlock"` //专属等级是否解锁
	IsArmy          int              `json:"isarmy"`          //是否是佣兵
}

type JS_HeroExtAttr struct {
	Key   int     `json:"k"` //! 属性类型
	Value float32 `json:"v"` //! 属性值
}

type JS_HeroParam struct {
	Heroid  int              `json:"heroid"` //! 武将id
	Param   []float32        `json:"param"`  //! 武将属性
	Hp      float32          `json:"hp"`     //! 当前hp
	Energy  int              `json:"energy"` //! 当前怒气
	Pos     int              `json:"pos"`
	ExtAttr []JS_HeroExtAttr `json:"ext"` //扩展属性
}

const (
	MAX_UI_POS          = 6
	MAX_FIGHT_POS       = 9
	CURRENT_UI_POS      = 6
	TEAMTYPE_DEFAULT    = 1
	TEAMTYPE_EXPEDITION = 2
	MAX_ARMY_NUM        = 5
	MAX_EQUIP_NUM       = 6
	MAX_TEAM_TYPE       = 2
)

// 魔宠, 虎符, 佣兵, 军旗, 装备, 绑在阵容上
type TeamAttr struct {
	HorseKeyId int   `json:"horse_key_id"`
	TigerKeyId int   `json:"tiger_key_id"`
	ArmyId     int   `json:"army_id"`
	FlagIds    []int `json:"flag_ids"`
	EquipIds   []int `json:"equipIds"`
}

// 兵法: 已废弃
type JS_ArmsSkill struct {
	Id    int `json:"id"`    // 兵法Id
	Level int `json:"level"` // 等级
}

type Js_TalentSkill struct {
	SkillId int `json:"skillid"` // 技能Id
	Level   int `json:"level"`   // 技能等级
}

type TeamPos struct {
	UIPos    [MAX_UI_POS]int    `json:"uipos"`
	FightPos [MAX_FIGHT_POS]int `json:"fightpos"`
	TeamAttr []*TeamAttr        `json:"team_attr"` // 属性阵位绑定
}

// 巨兽信息
type JsBoss struct {
	HeroId    int            `json:"heroid"`
	HeroParam *JS_HeroParam  `json:"heroinfo"`
	ArmsSkill []JS_ArmsSkill `json:"armsskill"`
}

type JS_FightInfo struct {
	Rankid       int            `json:"rankid"`       // 类型
	Uid          int64          `json:"uid"`          // id
	Uname        string         `json:"uname"`        // 名字
	UnionName    string         `json:"union"`        //! 军团名字
	Iconid       int            `json:"iconid"`       // icon
	Camp         int            `json:"camp"`         // 阵营
	Level        int            `json:"level"`        // 等级
	Vip          int            `json:"vip"`          // Vip 等级
	Inspire      int            `json:"inspire"`      //! 鼓舞等级，新增
	Beautyid     int            `json:"beautyid"`     // 出战后宫
	Defhero      []int          `json:"defhero"`      // 出战武将
	Heroinfo     []JS_HeroInfo  `json:"heroinfo"`     // 各武将信息
	Morale       int            `json:"morale"`       // 士气
	HeroParam    []JS_HeroParam `json:"heroparam"`    // 各武将属性
	Deffight     int64          `json:"deffight"`     // 玩家总战力
	FightTeam    int            `json:"fightteam"`    // 出战兵营
	FightTeamPos TeamPos        `json:"fightteampos"` // 出战兵营
	BossInfo     *JsBoss        `json:"bossinfo"`     // 巨兽信息
	Portrait     int            `json:"portrait"`     // 边框  20190412 by zy
}

type FightClient struct {
	Id     int64            //! 战斗id
	Result int              //! 结果
	Fight  [2]*JS_FightInfo //! 交战双方数据
	Info   []FightHero      //! 战斗结果数据
	Random int              //! 随机数
	Time   int              //! 交战时间
}

type FightServerExtendNode struct {
	NewId    int64
	status   int //状态 0从游戏服务器获取之后的状态 1
	serverid int
	addTime  int64
	FightServerNode
}

type FightResultExtendNode struct {
	node       *FightServerExtendNode
	resultNode *FightResultNode
}

/////////////////////////////////////////////////////////////////////////////////////

type HandlerPrototype func(r *http.Request) interface{}

type FightServer struct {
	Con              *Config         //! 配置
	Wait             *sync.WaitGroup //! 同步阻塞
	ShutDown         bool            //! 是否正在执行关闭
	LockerCon        *sync.RWMutex
	StartTime        int64            //开始运行时间
	NextWriteLogTime int64            //下次写日志时间
	SuccessCount     [INDEX_MAX]int   //计算成功场次
	NeedCount        [INDEX_MAX]int   //请求数
	NextCalTimes     [INDEX_MAX]int64 //下次统计时间
}

var fightServerSingleton *FightServer = nil

//! public
func GetServer() *FightServer {
	if fightServerSingleton == nil {
		fightServerSingleton = new(FightServer)
		fightServerSingleton.Con = new(Config)
		fightServerSingleton.LockerCon = new(sync.RWMutex)
		fightServerSingleton.Wait = new(sync.WaitGroup)
	}

	return fightServerSingleton
}

//! 载入配置文件
func (self *FightServer) InitConfig() {
	configFile, err := ioutil.ReadFile("./config.json") ///尝试打开配置文件
	if err != nil {
		log.Fatal("config err 1")
	}

	err = json.Unmarshal(configFile, self.Con)
	if err != nil {
		log.Fatal("chat InitConfig err", err.Error())
	}

	GetLogMgr().SetLevel(self.Con.LogCon.LogLevel, self.Con.LogCon.LogConsole)
}

func (self *FightServer) ReloadConfig() {
	configFile, err := ioutil.ReadFile("./config.json") ///尝试打开配置文件
	if err != nil {
		log.Fatal("config err 1")
		return
	}

	newCon := new(Config)
	err = json.Unmarshal(configFile, newCon)
	if err != nil {
		log.Fatal("chat InitConfig err", err.Error())
		return
	}

	GetFightMgr().ReloadFightServers(newCon)

	self.LockerCon.Lock()
	self.Con = newCon
	self.LockerCon.Unlock()
}

func (self *FightServer) GetServerCon(serverid int) *ServerConfig {

	self.LockerCon.RLock()
	defer self.LockerCon.RUnlock()

	for i := 0; i < len(self.Con.ServersCon); i++ {
		if self.Con.ServersCon[i].Id == serverid {
			return &(self.Con.ServersCon[i])
		}
	}
	return nil
}

func (self *FightServer) Close() {
	self.ShutDown = true
	self.Wait.Wait()
	LogFatal("server shutdown")
}

func (self *FightServer) FightServer(w http.ResponseWriter, r *http.Request) {
	msgtype := r.FormValue("msgtype")
	if msgtype == "get" {
		client := GetFightMgr().GetTopFight()
		if client == nil {
			//log.Println("无战斗数据")
			w.Write([]byte("false"))
		} else {
			//LogDebug("有战斗数据")
			var node FightServerNode
			node.Id = client.NewId
			node.Att = client.Att
			node.Def = client.Def
			node.Random = client.Random
			w.Write(HF_JtoB(&node))
		}
		return
	} else if msgtype == "set" {
		data := r.FormValue("data")
		var node FightResultNode
		err := json.Unmarshal([]byte(data), &node)
		if err == nil {
			GetFightMgr().SetResult(&node)
			for i := 0; i < len(self.SuccessCount); i++ {
				self.SuccessCount[i]++
			}
		} else {
			LogDebug("有战报解析失败")
		}
		client := GetFightMgr().GetTopFight()
		if client == nil {
			//log.Println("无战斗数据")
			w.Write([]byte("false"))
		} else {
			LogDebug("有战斗数据")
			var node FightServerNode
			node.Id = client.NewId
			node.Att = client.Att
			node.Def = client.Def
			node.Random = client.Random
			w.Write(HF_JtoB(&node))
		}
		return
	}

}

//! run
func (self *FightServer) Run() {

	self.StartTime = time.Now().Unix()

	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			self.onTime()
		}
	}

	ticker.Stop()
}

func (self *FightServer) onTime() {
	now := time.Now().Unix()
	//每分钟输入到控制台
	if now >= self.NextCalTimes[INDEX_MINUTE] {
		self.WriteControl()
	}
	//十分钟写日志
	if now >= self.NextWriteLogTime {
		self.WriteLog()
	}
	//清空数据
	for i := 0; i < len(self.NextCalTimes); i++ {
		if now >= self.NextCalTimes[i] {
			self.GetNextConfig(i)
		}
	}
}

func (self *FightServer) GetNextConfig(index int) {
	now := time.Now()
	switch index {
	case INDEX_MINUTE:
		self.NextCalTimes[index] = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location()).Unix() + game.MIN_SECS
	case INDEX_HOUR:
		self.NextCalTimes[index] = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).Unix() + game.HOUR_SECS
	case INDEX_DAY:
		self.NextCalTimes[index] = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix() + game.DAY_SECS
	default:
		return
	}

	self.SuccessCount[index] = 0
	self.NeedCount[index] = 0
}

func (self *FightServer) WriteControl() {
	now := time.Now().Unix()
	str := ""
	str = fmt.Sprintf("开启时间:%s  当前时间:%s  总运行时间:%d秒 ", time.Unix(self.StartTime, 0).Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), now-self.StartTime)
	log.Println(str)
	str = fmt.Sprintf("累计返回:%d  当前剩余:%d  累计需求计算:%d ", self.SuccessCount[INDEX_ALL], GetFightMgr().GetNeedNow(), self.NeedCount[INDEX_ALL])
	log.Println(str)
	str = fmt.Sprintf("分钟请求:%d  分钟返回:%d  ", self.NeedCount[INDEX_MINUTE], self.SuccessCount[INDEX_MINUTE])
	log.Println(str)
	str = fmt.Sprintf("小时请求:%d  小时返回:%d  ", self.NeedCount[INDEX_HOUR], self.SuccessCount[INDEX_HOUR])
	log.Println(str)
	str = fmt.Sprintf("每日请求:%d  每日返回:%d  ", self.NeedCount[INDEX_DAY], self.SuccessCount[INDEX_DAY])
	log.Println(str)
}

func (self *FightServer) WriteLog() {
	now := time.Now().Unix()
	str := ""
	str = fmt.Sprintf("开启时间:%s  当前时间:%s  总运行时间:%d秒 ", time.Unix(self.StartTime, 0).Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), now-self.StartTime)
	LogDebug(str)
	str = fmt.Sprintf("累计返回:%d  当前剩余:%d  累计需求计算:%d ", self.SuccessCount[INDEX_ALL], GetFightMgr().GetNeedNow(), self.NeedCount[INDEX_ALL])
	LogDebug(str)
	str = fmt.Sprintf("分钟请求:%d  分钟返回:%d  ", self.NeedCount[INDEX_MINUTE], self.SuccessCount[INDEX_MINUTE])
	LogDebug(str)
	str = fmt.Sprintf("小时请求:%d  小时返回:%d  ", self.NeedCount[INDEX_HOUR], self.SuccessCount[INDEX_HOUR])
	LogDebug(str)
	str = fmt.Sprintf("每日请求:%d  每日返回:%d  ", self.NeedCount[INDEX_DAY], self.SuccessCount[INDEX_DAY])
	LogDebug(str)

	self.NextWriteLogTime = time.Now().Unix() + 600
}
