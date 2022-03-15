package main

type Son_CampFightInfo struct {
	Uid      int64  `json:"uid"`
	Index    int    `json:"index"`
	Name     string `json:"name"`
	Icon     int    `json:"icon"`
	Fight    int    `json:"fight"`
	Level    int    `json:"level"`
	Camp     int    `json:"camp"`
	Kill     int    `json:"kill"`
	Elite    int    `json:"elite"`
	Class    int    `json:"class"`
	ArmsType int    `json:"armstype"`
	Honor    int    `json:"honor"`
	Buffer   int    `json:"buffer"`
}

type FightHero struct {
	Heroid int `json:"id"`   //! 武将id
	Hp     int `json:"hp"`   //! hp
	Energy int `json:"rage"` //! 怒气
}

//!
type C2S_Uid struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
}

//! 注册
type C2S_Reg struct {
	Ctrl     string `json:"ctrl"`
	Uid      int64  `json:"uid"`
	Account  string `json:"account"`
	Password string `json:"password"`
	ServerId int    `json:"serverid"`
}

//! 指引id
type C2S_ZyInfo struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
	Zyid int    `json:"zyid"`
}

type C2S_CreateRole struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
	Name string `json:"name"`
	Icon int    `json:"icon"`
	Face int    `json:"face"`
}

//! 开始战斗
type C2S_BoxPass struct {
	Ctrl      string `json:"ctrl"`
	Uid       int64  `json:"uid"`
	Passid    int    `json:"passid"`
	NoPass    int    `json:"nopass"`
	MissionId int    `json:"missionid"`
	Step      int    `json:"step"`
	WarNum    int    `json:"warnum"`
	WorkNum   int    `json:"wroknum"`
}

//! 觐见
type C2S_JJPass struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Passid int    `json:"passid"`
}

//! 开始战斗
type C2S_BeginPass struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Passid int    `json:"passid"`
}

type C2S_TakeNobilitytask struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	TaskId int    `json:"taskid"`
}

type C2S_ChaptErevents struct {
	Ctrl    string `json:"ctrl"`
	Uid     int64  `json:"uid"`
	Chapter int    `json:"chapter"`
}

type C2S_WinPass struct {
	Ctrl      string `json:"ctrl"`
	Uid       int64  `json:"uid"`
	Missionid int    `json:"missionid"`
}

//! 结束战斗
type C2S_EndPass struct {
	Ctrl      string `json:"ctrl"`
	Uid       int64  `json:"uid"`
	Passid    int    `json:"passid"`
	Star      int    `json:"star"`
	Index     int    `json:"index"`
	MissionId int    `json:"missionid"`
	Step      int    `json:"step"`
	WarNum    int    `json:"warnum"`
	WorkNum   int    `json:"wroknum"`
	FightTime int    `json:"fighttime"`
}

//! 抽奖
type C2S_Find struct {
	Ctrl     string `json:"ctrl"`
	Uid      int64  `json:"uid"`
	Findtype int    `json:"findtype"`
	Free     int    `json:"free"`
}

//! 升阶
type C2S_UpColor struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Heroid int    `json:"heroid"`
}

//! 合成英雄
type C2S_Synthesis struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Heroid string `json:"heroid"`
}

//! 签到奖励
type C2S_CheckinAward struct {
	Ctrl  string `json:"ctrl"`
	Uid   int64  `json:"uid"`
	Index int    `json:"index"`
}

//! 推图信息
type C2S_SetMission struct {
	Ctrl      string `json:"ctrl"`
	Uid       int64  `json:"uid"`
	MissionId int    `json:"missionid"`
	Step      int    `json:"step"`
	WarNum    int    `json:"warnum"`
	WorkNum   int    `json:"wroknum"`
}

//! 天下大事
type C2S_WordEvent struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
	Ver  int    `json:"ver"`
}

//! 编队
type C2S_CampTeam struct {
	Ctrl string         `json:"ctrl"`
	Uid  int64          `json:"uid"`
	Team [3]JS_CampTeam `json:"team"`
}

//! 阵营
type JS_CampTeam struct {
	Hero     [5]int `json:"hero"`
	Beautyid int    `json:"beautyid"`
	State    int    `json:"state"` //! 0空闲 1进攻 2防守
	Cityid   int    `json:"cityid"`
	Name     string `json:"name"`
	UseSys   int    `json:"usesys"`
	Icon     int    `json:"icon"`
}

//! 移动部队
type C2S_BigMap struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
	Type int    `json:"type"`
}

//! 移动部队
type C2S_BossAction struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Id     int    `json:"id"`
	Action int    `json:"action"`
}

//! 得到城池数据
type C2S_GetCityInfo struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Cityid int    `json:"cityid"`
}

type C2S_GetGem struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	ItemId int    `json:"itemid"`
	Num    int    `json:"num"`
}

//! 移动部队
type C2S_MoveTeamBegin struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Index  int    `json:"index"`
	Cityid int    `json:"cityid"`
	Begin  int    `json:"begin"`
}

//! 移动部队
type C2S_MoveTeam struct {
	Ctrl  string `json:"ctrl"`
	Uid   int64  `json:"uid"`
	Index int    `json:"index"`
	Way   []int  `json:"way"`
}

//设置阵营
type C2S_SetCamp struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
	Camp int    `json:"camp"`
}

type C2S_GetTopKing struct {
	Ctrl       string `json:"ctrl"`
	Uid        int64  `json:"uid"`
	Camp       int    `json:"camp"`
	TopverKing int    `json:"topverking"`
}

type C2S_FinishActivity struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
	Id   int    `json:"id"`
}

type C2S_Barrage struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Cityid int    `json:"cityid"`
	Size   int    `json:"size"`
	Text   string `json:"text"`
	Red    int    `json:"red"`
	Green  int    `json:"green"`
	Blue   int    `json:"blue"`
}

//! 激活宝箱碎片
type C2S_LevelUpBoxSP struct {
	Ctrl  string `json:"ctrl"`
	Uid   int64  `json:"uid"`
	Index int    `json:"index"`
}

//！好友推荐
type C2S_FriendCommend struct {
	Ctrl    string `json:"ctrl"`
	Uid     int64  `json:"uid"`
	Refresh int64  `json:"refresh"`
}

type C2S_StatisticsInfo struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
}

//！一键申请好友
type C2S_FrienDapply struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
	Pid  int64  `json:"pid"`
}

//！一键通过好友
type C2S_FrienDorder struct {
	Ctrl  string `json:"ctrl"`
	Uid   int64  `json:"uid"`
	Pid   int64  `json:"pid"`
	Agree int64  `json:"agree"`
}

//！一键通过好友
type C2S_GMPassId struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Passid int64  `json:"passid"`
}

//！神器升级
type C2S_ShenqilvUp struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Heroid int    `json:"heroid"`
	Id     int    `json:"id"`
	Type   int    `json:"type"`
}

//! 武将升星
type C2S_UpHeroStar struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Heroid int    `json:"heroid"`
}

//! 兵种操作
type C2S_SoldierType struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Heroid int    `json:"heroid"`
	Type   int    `json:"type"`
}

//! 进出阵营战
type C2S_CampFightMove struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Cityid int    `json:"cityid"`
	Type   int    `json:"type"` //! 0进 1出
}

//! 国战单挑-请求
type C2S_CampFightSoloReq struct {
	Ctrl      string `json:"ctrl"`
	Uid       int64  `json:"uid"`
	CityId    int    `json:"cityid"`    //! 城市
	WarPlayId int    `json:"warplayid"` //! 玩法Id
	Again     int    `json:"again"`     //! 0初次 1继续单挑
	CityPart  int    `json:"citypart"`  //! 目标选择 0-主城 1-5为据点
}

//! 一骑讨
type S2C_CampFightSoloList struct {
	Cid       string                `json:"cid"`
	FightInfo [3]*Son_CampFightInfo `json:"fightinfo"` //一骑讨战斗信息
}

//! 国战solo
type S2C_CampFightSolo2Begin struct {
	Cid    string `json:"cid"`
	SoloId int    `json:"soloid"`
}

//! 获取外交数据
type S2C_DiplomacyInfo struct {
	Cid        string   `json:"cid"`
	Index      [3]int   `json:"index"`
	Align      [3]int   `json:"align"`
	AttackCity [3][]int `json:"attackcity"`
	Time       int      `json:"time"`
	FightTime  [3]int   `json:"fighttime"`
	Status     [3]int   `json:"status"`
	CityArr    []int    `json:"cityarr"`
	UnionArr   [4]int64 `json:"unionarr"`
}

//! Buffer结构
type JS_Buffer struct {
	BufferId  int   `json:"bufferid"`  // buffer id
	Time      int64 `json:"time"`      // 生效时间
	Overlying int   `json:"overlying"` // 叠加层数
	Effect    int   `json:"effect"`    // 持续时间
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

//! 玩法数据结构
type S2C_CampFightPlayInfo struct {
	Cid           string            `json:"cid"`
	Cityid        int               `json:"cityid"`
	InspireLevel  int               `json:"inspirelevel"`
	InspireExp    int               `json:"inspireexp"`
	PlayerNum     [2]int            `json:"playernum"`
	Occupy        [2]int            `json:"occupy"`
	CityPart      int               `json:"citypart"`
	ReviveTime    int               `json:"revivetime"`
	Playlist      []JS_PlayNode     `json:"playlist"`
	BuffList      [2][]JS_Buffer    `json:"bufflist"`
	SelfBuff      []JS_Buffer       `json:"selfbuff"`
	PartBuffList  [5][2][]JS_Buffer `json:"partbufflist"`
	PartPlayerNum [5][2]int         `json:"partplayernum"`
	PartOccupy    [5][2]int         `json:"partoccupy"`
	PartCamp      [5]int            `json:"partcamp"`
}

//! 国战单挑
type C2S_CampFightSolo2 struct {
	Ctrl    string `json:"ctrl"`
	Uid     int64  `json:"uid"`
	Cityid  int    `json:"cityid"`
	MyIndex int    `json:"myindex"`
	Pid     int64  `json:"pid"`
	Index   int    `json:"index"`
}

//! 国战结果
type C2S_CampFightSoloEnd struct {
	Ctrl   string      `json:"ctrl"`
	Uid    int64       `json:"uid"`
	Cityid int         `json:"cityid"`
	Soloid int         `json:"soloid"`
	Result int         `json:"result"`
	Info   []FightHero `json:"info"`
}

//! 国战五虎争雄报名
type C2S_CampFight55Req struct {
	Ctrl      string `json:"ctrl"`
	Uid       int64  `json:"uid"`
	Cityid    int    `json:"cityid"`
	WarPlayid int    `json:"warplayid"`
}

//! 使用道具
type C2S_UseItem struct {
	Ctrl   string `json:"ctrl"`
	Uid    int64  `json:"uid"`
	Itemid string `json:"itemid"`
	Num    int    `json:"num"`
	Destid int    `json:"destid"`
}

//! 恢复行动力
type C2S_TimePower struct {
	Ctrl  string `json:"ctrl"`
	Uid   int64  `json:"uid"`
	Index int    `json:"index"`
}

type C2S_SetGuild struct {
	GuildId int    `json:"guildeid"`
	Ctrl    string `json:"ctrl"`
	Uid     int64  `json:"uid"`
}

//! 获取关卡的任务状态
type C2S_GetMission struct {
	Ctrl      string `json:"ctrl"`
	MissionId int    `json:"missionid"`
	Uid       int64  `json:"uid"`
}

//! 购买军令
type C2S_Collection struct {
	Ctrl  string `json:"ctrl"`
	Uid   int64  `json:"uid"`
	Index int    `json:"index"`
}

type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type C2S_FinishEvents struct {
	Ctrl    string `json:"ctrl"`
	ThingId int    `json:"thing_id"` // 物件Id
	EventId int    `json:"event_id"` // 事件Id
	Uid     int64  `json:"uid"`
}

//设置红色标记
type C2S_SetRedIcon struct {
	Ctrl string `json:"ctrl"`
	Id   int    `json:"id"` // 功能Id
}

//! 上阵英雄
type C2S_AddTeamUIPos struct {
	Ctrl     string `json:"ctrl"`
	Index    int    `json:"index"`    //! UI布阵索引1~5
	HeroId   int    `json:"heroid"`   //! 英雄Id
	TeamType int    `json:"teamtype"` //! 布阵类型, 1副本 2竞技场
	Uid      int64  `json:"uid"`
}

//! 激活资质
type C2S_ActivateStar struct {
	Ctrl   string `json:"ctrl"`
	HeroId int    `json:"heroid"` //! 英雄Id
	Index  int    `json:"index"`  //! 资质孔的序号
	Uid    int64  `json:"uid"`
}

//! 一键升星
type C2S_UpStarAuto struct {
	Ctrl   string `json:"ctrl"`
	HeroId int    `json:"heroid"` //! 英雄Id
	Uid    int64  `json:"uid"`
}

//! 交换战斗位置
type C2S_SwapFightPos struct {
	Ctrl     string `json:"ctrl"`
	Index1   int    `json:"index1"`   //! 战斗位置1
	Index2   int    `json:"index2"`   //! 战斗位置2
	TeamType int    `json:"teamtype"` //! 布阵类型, 1副本 2竞技场
	Uid      int64  `json:"uid"`
}

//! 完成任务
type C2S_TaskFinish struct {
	Ctrl     string `json:"ctrl"`
	Taskid   string `json:"taskid"`
	Tasktype int    `json:"tasktype"`
	Uid      int64  `json:"uid"`
}

type C2S_AddGuide struct {
	Ctrl    string `json:"ctrl"`
	GuideId int    `json:"guide_id"`
	Uid     int64  `json:"uid"`
}

//! 装备相关操作
type C2S_EquipAction struct {
	Ctrl         string `json:"ctrl"`
	Action       int    `json:"action"`       //! action=1.装备合成 2.装备分解 3.装备穿戴
	Itemid       int    `json:"itemid"`       //! 合成的碎片Id[合成时发送]
	CompundNum   int    `json:"compoundnum"`  //! 合成的装备数量[合成时发送]
	HeroId       int    `json:"heroid"`       //! 英雄Id
	KeyId        int    `json:"keyid"`        //! 装备Id
	Pos          int    `json:"pos"`          //! 第几个位置装备[脱装备时发送]
	GemId        int    `json:"gemid"`        //! 合成的宝石Id
	RemoveKeyIds []int  `json:"removekeyids"` //! 分解的装备Id
	TeamType     int    `json:"team_type"`    // 阵营类型 默认填1
	Index        int    `json:"index"`        // 上阵下标
	Uid          int64  `json:"uid"`
}

type C2S_PassMission struct {
	Ctrl      string `json:"ctrl"`
	Uid       int64  `json:"uid"`
	ChpaterId int    `json:"chpater_id"`
}
