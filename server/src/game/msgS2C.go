package game

//! server2client
//! msgid
type S2C_MsgId struct {
	Cid     string `json:"cid"`
	CurTime int64  `json:"curtime"`
}

//! 注册
type S2C_Reg struct {
	Cid      string `json:"cid"`
	Uid      int64  `json:"uid"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Creator  string `json:"creator"`
}

//! 获取平台信息
type S2C_PlatFromInfo struct {
	Cid string `json:"cid"`
}

//! 基本信息
type S2C_UserBaseInfo struct {
	Cid      string           `json:"cid"`
	Baseinfo Son_UserBaseInfo `json:"baseinfo"`
}
type Son_UserBaseInfo struct {
	Checkinaward       int    `json:"checkinaward"`
	Checkinnum         int    `json:"checkinnum"`
	Exp                int    `json:"exp"`
	Face               int    `json:"face"`
	Gem                int    `json:"gem"`
	Gold               int    `json:"gold"`
	Iconid             string `json:"iconid"`
	Ischeckin          bool   `json:"ischeckin"`
	Isrename           int    `json:"isrename"`
	Juqingid           int    `json:"juqingid"`
	Juqingid2          int    `json:"juqingid2"`
	Lastcheckintime    string `json:"lastcheckintime"`
	Lastlivetime       string `json:"lastlivetime"`
	Lastlogintime      string `json:"lastlogintime"`
	Level              int    `json:"level"`
	Levelaward         int    `json:"levelaward"`
	Loginaward         int    `json:"loginaward"`
	Logindays          int    `json:"logindays"`
	Morale             int    `json:"morale"`
	Partyid            int    `json:"partyid"`
	Position           int    `json:"position"`
	Regtime            string `json:"regtime"`
	Skillpoint         int    `json:"skillpoint"`
	Splastupdatatime   int64  `json:"splastupdatatime"`
	Tili               int    `json:"tili"`
	Tililastupdatatime int64  `json:"tililastupdatatime"`
	Uid                int64  `json:"uid"`
	Uname              string `json:"uname"`
	Vip                int    `json:"vip"`
	Vipexp             int    `json:"vipexp"`
	Worldaward         int    `json:"worldaward"`
	Zhiyinid           int    `json:"zhiyinid"`
	Zhiyinid1          int    `json:"zhiyinid1"`
	Citylevel          int    `json:"citylevel"`
	Camp               int    `json:"camp"`
	City               int    `json:"city"`
	Day                int    `json:"day"`
	Promotebox         int    `json:"promotebox"`
	OpenServer         int64  `json:"openserver"`
	Channelid          string `json:"channelid"`
	Account            string `json:"account"`
	FitServer          int    `json:"fitserver"`
	Soul               int    `json:"soul"`          //! 魂石
	TechPoint          int    `json:"techpoint"`     //! 科技点
	BossMoney          int    `json:"bossmoney"`     //! 水晶币
	TowerStone         int    `json:"towerstone"`    //! 镇魂石
	Portrait           int    `json:"portrait"`      //! 头像框
	CampOk             int    `json:"camp_ok"`       // 阵营ok
	NameOk             int    `json:"name_ok"`       // 名字ok
	GuildId            int    `json:"guildeid"`      // 指引Id
	RedIcon            int    `json:"red_icon"`      // 小红点
	UserSignature      string `json:"usersignature"` // 个人宣言
}

//! sp
type S2C_UserSP struct {
	Cid  string `json:"cid"`
	Sp   int    `json:"sp"`
	Time int    `json:"time"`
	Uid  int64  `json:"uid"`
}

//! 背包
type S2C_BagInfo struct {
	Cid    string     `json:"cid"`
	Baglst []PassItem `json:"baglst"`
}

//! 关卡
type S2C_PassInfo struct {
	Cid        string       `json:"cid"`
	Passinfo   Son_UserPass `json:"passinfo"`
	IsFight    int          `json:"isfight"`
	TotalStars int          `json:"totalstars"`
}
type Son_UserPass struct {
	WarInfo     string `json:"warinfo"`
	PassInfo    string `json:"passinfo"`
	MissionInfo string `json:"mission"`
	BoxInfo     string `json:"boxinfo"`
	StarBoxInfo string `json:"starboxinfo"`
	JJInfo      string `json:"jjinfo"`
}

//! 武将
type S2C_HeroInfo struct {
	Cid        string      `json:"cid"`
	Herolst    []*Hero     `json:"herolst"`
	Newhero    bool        `json:"newhero"`
	TotalStars int         `json:"totalstars"`
	Reborn     int         `json:"reborn"`
	BuyPosNum  int         `json:"buyposnum"` //英雄栏扩展次数
	AutoFire   int         `json:"autofire"`  //自动分解开关
	BackOpen   int         `json:"backopen"`  //回退功能开启开关
	HandBook   map[int]int `json:"handbook"`  //图鉴key:heroid   value:1已开启可领取  2已领取
}

//! 图鉴
type S2C_HandBook struct {
	Cid      string                `json:"cid"`
	HandBook map[int]*HandBookInfo `json:"handbook"`
}

//! 同步英雄信息
type S2C_SynHero struct {
	Cid  string `json:"cid"`  //! cid
	Hero *Hero  `json:"hero"` //!
}

//! 体力
type S2C_UserPower struct {
	Cid        string `json:"cid"`
	Tili       int    `json:"tili"`
	Time       int    `json:"time"`
	Uid        int64  `json:"uid"`
	WorldLevel int    `json:"worldLevel"`
}

//! 抽奖
type S2C_FindInfo struct {
	Cid             string               `json:"cid"`
	BaseFindInfo    []*FindPool          `json:"basefindinfo"`
	RewardInfo      *FindRewardInfo      `json:"rewardinfo"`
	FindWishInfo    []*FindWishInfo      `json:"findwishinfo"`
	Astrology       *FindAstrology       `json:"astrology"`
	SelfSelection   *SelfSelection       `json:"selfselection"`
	LuckyPoolConfig []LuckyPassItem      `json:"luckypoolconfig"`
	LuckyFindRecord []*LotteryDrawRecord `json:"luckyfindrecord"`
	FindGiftProcess	int					 `json:"drawgiftprocess"`
	IsGotFindGift	[]int				 `json:"isgotfindgift"`		// 是否已领取累积抽奖奖励宝箱 0-未领取 1-已领取
}

type S2C_GetFindReward struct {
	Cid             string               `json:"cid"`
	FindGiftProcess	int					 `json:"drawgiftprocess"`
	IsGotFindGift	[]int				 `json:"isgotfindgift"`		// 是否已领取累积抽奖奖励宝箱 0-未领取 1-已领取
	GetItems		[]PassItem			 `json:"getitems"`			// 获得奖励
}

//! 修改阵型
type S2C_ChgCurTeam struct {
	Cid     string `json:"cid"`
	CurTeam int    `json:"curteam"`
}

//! 充值记录
type S2C_RechargeInfo struct {
	Cid       string           `json:"cid"`
	Info      San_UserRecharge `json:"info"`
	Vip       int              `json:"vip"`
	Vipexp    int              `json:"vipexp"`
	CurGem    int              `json:"curgem"`
	Ret       int              `json:"ret"`
	Vipbox    int64            `json:"vipbox"`
	Fundtype  int64            `json:"fundtype"`
	Fundget   []int64          `json:"fundget"`
	Fundtotal int              `json:"fundtotal"`
	WarOrderLimit []JS_WarOrderLimit `json:"warorderlimit"`
}

//! 商会人数
type S2C_GetFundTotal struct {
	Cid       string `json:"cid"`
	FundTotal int    `json:"fundtotal"`
}

//! 任务
type S2C_TaskInfo struct {
	Cid      string             `json:"cid"`
	Info     []JS_TaskInfo      `json:"info"`     // 任务信息
	Liveness []*JS_LivenessInfo `json:"liveness"` // 活跃度
}

//目标系统 爵位
type S2C_NobilityTask struct {
	Cid       string                       `json:"cid"`
	TaskInfo  map[int]*JS_NobilityTaskInfo `json:"taskinfo"`
	Level     int                          `json:"level"`
	GetReward map[int]int                  `json:"getreward"`
}

//目标系统 爵位 任务同步
type S2C_UpdateNobilityTask struct {
	Cid      string                `json:"cid"`
	TaskInfo []JS_NobilityTaskInfo `json:"taskinfo"`
}

//! 目标任务
type S2C_TargetTaskInfo struct {
	Cid        string                        `json:"cid"`
	Info       []JS_TargetTaskInfo           `json:"info"`       // 任务信息
	SystemInfo map[int]int                   `json:"systeminfo"` // 系统组信息
	BuyLevel   map[int]int                   `json:"buylevel"`   // 徽章等级
	Buyrewards map[int]map[int]map[int]*Item `json:"buyrewards"` // 购买后补充的奖励,上线发一次，后面客户端自己算key顺序:system->lv->物品
}

type S2C_VipRechargeInfo struct {
	Cid       string                      `json:"cid"`
	TaskInfos map[int]*JS_VipRechargeInfo `json:"taskinfos"` //完成详情
}

//! 周期任务
type S2C_WeekPlanInfo struct {
	Cid       string                  `json:"cid"`
	Info      Son_WeekPlanInfo        `json:"info"`
	Ret       int                     `json:"ret"`
	Point     int                     `json:"point"`
	IsGetMark map[int]int             `json:"isgetmark"`
	BuyTime   int64                   `json:"buytime"`
	Config    map[int]*SevendayConfig `json:"config"`
	Stage     int                     `json:"stage"` //1:1~7   2"8~14
}

type S2CTakeNobilityTask struct {
	Cid      string                       `json:"cid"`
	TaskInfo map[int]*JS_NobilityTaskInfo `json:"taskinfo"`
	Level    int                          `json:"level"` //爵位
	Award    []PassItem                   `json:"award"`
}

//转盘
type S2C_TurnTableInfo struct {
	Cid           string           `json:"cid"`
	NowStage      int              `json:"nowstage"`      //! 当前阶段
	NowCount      int              `json:"nowcount"`      //! 当前阶段转到次数
	NextTime      int64            `json:"nexttime"`      //! 下次可以转的时间
	TurnTableinfo []*TurnTableItem `json:"turntableinfo"` //! 任务信息
}

type S2C_AccessCardInfo struct {
	Cid       string             `json:"cid"`
	Group     int                `json:"group"`     //!
	TaskInfo  []*JS_TaskInfo     `json:"taskinfo"`  //! 任务信息
	AwardInfo []*AccessCardAward `json:"awardinfo"` //! 积分奖励信息
	Point     int                `json:"point"`     //!
}

type S2C_AccessGetRank struct {
	Cid        string                     `json:"cid"`
	Rank       []*AccessCardRecordTopNode `json:"rank"`       //! 积分奖励信息
	Self       *AccessCardRecordTopNode   `json:"self"`       //!
	StartTime  int64                      `json:"starttime"`  //开始时间
	EndTime    int64                      `json:"endtime"`    //结束时间
	RewardTime int64                      `json:"rewardtime"` //发奖时间
	HasReward  int                        `json:"hasreward"`  //0没发 1已发
}

type S2C_GetFund struct {
	Cid      string           `json:"cid"`
	TaskInfo *JS_FundTaskInfo `json:"taskinfo"` //! 任务信息
	GetItems []PassItem       `json:"getitems"`
}



type S2C_DoTurnTable struct {
	Cid           string         `json:"cid"`
	TurnTableItem *TurnTableItem `json:"turntableitem"`
	NowStage      int            `json:"nowstage"` //! 当前阶段
	NowCount      int            `json:"nowcount"` //! 当前阶段转到次数
	NextTime      int64          `json:"nexttime"` //! 下次可以转的时间
	GetItems      []PassItem     `json:"getitems"`
}

type S2C_AccessCardTask struct {
	Cid      string       `json:"cid"`
	TaskInfo *JS_TaskInfo `json:"taskinfo"`
	GetItems []PassItem   `json:"getitems"`
	Point    int          `json:"point"`
}

type S2C_AccessCardAward struct {
	Cid             string           `json:"cid"`
	AccessCardAward *AccessCardAward `json:"accesscardaward"`
	GetItems        []PassItem       `json:"getitems"`
}

type S2C_GetAccessCardRecord struct {
	Cid    string              `json:"cid"`
	Record []*AccessCardRecord `json:"record"`
}

type S2CLevelUpNobility struct {
	Cid   string     `json:"cid"`
	Level int        `json:"level"` //爵位
	Award []PassItem `json:"award"`
}

type S2C_VipRechargeGift struct {
	Cid      string              `json:"cid"`
	TaskInfo *JS_VipRechargeInfo `json:"taskinfo"`
	GetItems []PassItem          `json:"getitems"`
}

type S2C_GetNobilityReward struct {
	Cid       string      `json:"cid"`
	Level     int         `json:"level"` //爵位
	Award     []PassItem  `json:"award"`
	GetReward map[int]int `json:"getreward"`
}

type Son_WeekPlanInfo struct {
	Regtime       int64               `json:"regtime"`  //角色创建时间
	Taskinfo      []JS_TaskInfo       `json:"taskinfo"` //任务列表
	Servercurtime int64               `json:"servercurtime"`
	TaskStatus    []JS_WeekPlanStatus `json:"taskstatus"`
}

//! 称号
type S2C_TitleInfo struct {
	Cid   string `json:"cid"`
	Level int    `json:"level"`
	Task  [3]int `json:"task"`
}

type S2C_WarTaskInfo struct {
	Cid  string        `json:"cid"`
	Info []JS_TaskInfo `json:"info"`
}

//! 任务
type S2C_TaskUpdate struct {
	Cid  string        `json:"cid"`
	Info []JS_TaskInfo `json:"info"`
}

type S2C_TaskUpdateBoss struct {
	Cid  string                 `json:"cid"`
	Info []*JS_ActivityBossInfo `json:"info"`
}

//! 任务
type S2C_TargetTaskUpdate struct {
	Cid  string              `json:"cid"`
	Info []JS_TargetTaskInfo `json:"info"`
}

//! 任务
type S2C_InterStellUpdate struct {
	Cid  string         `json:"cid"`
	Info []*JS_TaskInfo `json:"info"`
}

//! 任务
type S2C_LuckShopUpdate struct {
	Cid  string            `json:"cid"`
	Info []JS_LuckShopItem `json:"info"`
}

// 主线战令数据
type S2C_WarOrderLimit struct {
	Cid  string            `json:"cid"`
	Info []JS_WarOrderTask `json:"info"`
}

//! 限时礼包
type S2C_TimeGiftUpdate struct {
	Cid  string            `json:"cid"`
	Info []JS_TimeGiftItem `json:"info"`
}

//! 限时礼包
type S2C_ActivityGiftUpdate struct {
	Cid  string              `json:"cid"`
	Info []*ActivityGiftItem `json:"info"`
}

//! 限时礼包
type S2C_GrowthGiftUpdate struct {
	Cid  string            `json:"cid"`
	Info []*GrowthGiftItem `json:"info"`
}

type S2C_NewCampTaskInfo struct {
	Cid     string        `json:"cid"`
	Info    []JS_TaskInfo `json:"info"`
	CurBox  int           `json:"curbox"`
	Boxget  int64         `json:"boxget"`
	RefTime int64         `json:"reftime"`
	Coin    int           `json:"coin"`
}

//! 领取阵营任务宝箱
type S2C_GetCampBox struct {
	Cid    string     `json:"cid"`
	Award  []PassItem `json:"award"`
	CurBox int        `json:"curbox"`
	Boxget int64      `json:"boxget"`
}

type S2C_LuckShopInfo struct {
	Cid  string            `json:"cid"`
	Info []JS_LuckShopItem `json:"info"`
	Item []JS_ActivityBox  `json:"item"`
}

type S2C_FundInfo struct {
	Cid      string                 `json:"cid"`
	TaskInfo []JS_FundTaskInfo      `json:"taskinfo"` //任务
	FundInfo *JS_NewFundInfo        `json:"fundinfo"` //基金
	Config   map[int]*FundConfigMap `json:"config"`   //配置
}


type S2C_TimeGiftInfo struct {
	Cid            string            `json:"cid"`
	NextUpdatetime int64             `json:"nextupdatetime"`
	Info           []JS_TimeGiftItem `json:"info"`
	Item           []TimeGiftConfig  `json:"item"`
}

//! 商店
type S2C_ShopInfo struct {
	Cid  string         `json:"cid"`
	Info []Son_ShopInfo `json:"info"`
}
type Son_ShopInfo struct {
	Id            int64             `json:"id"`
	Lastupdtime   int64             `json:"Lastupdtime"`
	Refindex      int               `json:"refindex"`
	Shopgood      []*JS_NewShopInfo `json:"shopgood"`
	Shopnextgood  []JS_ShopInfo     `json:"shopnextgood"`
	Shoptype      int               `json:"shoptype"`
	Sysreftime    int64             `json:"sysreftime"`
	Todayrefcount int               `json:"todayrefcount"`
	Uid           int64             `json:"uid"`
}

type S2C_NewShopBuy struct {
	Cid         string          `json:"cid"`
	GetItems    []PassItem      `json:"getitems"`    //获得
	CostItems   []PassItem      `json:"costitems"`   //消耗
	NewShopInfo *JS_NewShopInfo `json:"newshopinfo"` //商品同步
}

type S2C_NewShopRefresh struct {
	Cid       string        `json:"cid"`
	CostItems []PassItem    `json:"costitems"` //消耗
	Info      *Son_ShopInfo `json:"info"`
}

//! 商店刷新
type S2C_ShopRef struct {
	Cid      string       `json:"cid"`
	Info     Son_ShopInfo `json:"info"`
	Shoptype int          `json:"shoptype"`
	Ret      int          `json:"ret"`
	Gem      int          `json:"gem"`  // 钻石个数
	Items    []PassItem   `json:"item"` // 刷新卷信息
}

//! 商店刷新
type S2C_ShopSysRef struct {
	Cid      string       `json:"cid"`
	Info     Son_ShopInfo `json:"info"`
	Shoptype int          `json:"shoptype"`
}

type S2C_TreasuryBuy struct {
	Cid  string       `json:"cid"`
	Ret  int          `json:"ret"`
	Item []PassItem   `json:"item"`
	Info *JS_ShopInfo `json:"info"`
	Type int          `json:"type"`
}

type S2C_TowerBuy struct {
	Cid  string     `json:"cid"`
	Ret  int        `json:"ret"`
	Item []PassItem `json:"item"`
	Num  int        `json:"num"`
}

//! 邮件
type S2C_MailInfo struct {
	Cid      string        `json:"cid"`
	Mailinfo []*JS_OneMail `json:"mailinfo"`
	Uid      int64         `json:"uid"`
	Gmail    []*JS_Mail    `json:"gmail"`
}

type S2C_MailAllItem struct {
	Cid  string     `json:"cid"`
	Item []PassItem `json:"item"`
}

//! 登陆成功
type S2C_LoginRet struct {
	Cid        string `json:"cid"`
	Ret        int    `json:"ret"`
	CheckCode  string `json:"checkcode"`
	Servertime int64  `json:"servertime"`
}

//! 进入排队
type S2C_LineUp struct {
	Cid        string `json:"cid"`
	WaitCount  int    `json:"waitcount"`
	PassMinute int    `json:"passminute"`
	PassSecond int    `json:"passsecond"`
}

//! 发送结果
type S2C_ResultMsg struct {
	Cid string `json:"cid"`
	Ret int    `json:"ret"`
}

//! 发送结果2
type S2C_Result2Msg struct {
	Cid string `json:"cid"`
}

//! 发送结果3
type S2C_Result3Msg struct {
	Cid string `json:"cid"`
	Ret bool   `json:"ret"`
}

//! 发送结果4
type S2C_Result4Msg struct {
	Cid string `json:"cid"`
	Ok  bool   `json:"ok"`
}

//! 开始战斗
type S2C_BeginPass struct {
	Cid     string     `json:"cid"`
	Passid  int        `json:"passid"`
	Outitem []PassItem `json:"outitem"`
	Tili    int        `json:"tili"`
}

//! 结束战斗
type S2C_EndPass struct {
	Cid     string       `json:"cid"`
	Passid  []int        `json:"passid"`
	Star    int          `json:"star"`
	War     JS_War       `json:"war"`
	Records *PassRecords `json:"records"` // 关卡记录
}

type S2C_PassSkip struct {
	Cid      string     `json:"cid"`
	Passid   []int      `json:"passid"`
	GetItems []PassItem `json:"getitem"`
}

//! 战斗记录
type S2C_PassRecord struct {
	Cid    string            `json:"cid"`
	Passid int               `json:"passid"`
	First  *San_PassRecord   `json:"first"`  // 最早通关
	Low    *San_PassRecord   `json:"low"`    // 最低战力通关
	Recent []*San_PassRecord `json:"recent"` // 最近通关
}

//! 国战开启
type S2C_GZOpen struct {
	Cid  string     `json:"cid"`
	Item []PassItem `json:"item"`
}

//! 直接胜利
type S2C_WinPass struct {
	Cid     string     `json:"cid"`
	Passid  []int      `json:"passid"`
	OutItem []PassItem `json:"item"`
	Box     int        `json:"box"`
}

//!
type S2C_JJPass struct {
	Cid     string     `json:"cid"`
	Passid  int        `json:"passid"`
	OutItem []PassItem `json:"item"`
}

//! 达成条件
type S2C_WarTask struct {
	Cid string `json:"cid"`
	War JS_War `json:"war"`
}

//! 扫荡关卡
type S2C_SwapPass struct {
	Cid     string     `json:"cid"`
	Uid     int64      `json:"uid"`
	Passid  int        `json:"passid"`
	Num     int        `json:"num"`
	Outitem []PassItem `json:"outitem"`
}

//! 更新体力
type S2C_UpdateTiLi struct {
	Cid  string `json:"cid"`
	Tili int    `json:"tili"`
	Time int    `json:"time"`
	Uid  int64  `json:"uid"`
}

//! 购买体力
type S2C_BuyTiLi struct {
	Cid    string `json:"cid"`
	Usegem int    `json:"usegem"`
}

//! 购买金币
type S2C_BuyGold struct {
	Cid         string `json:"cid"`
	Ret         int    `json:"ret"`
	Counts      int    `json:"counts"`
	Usegem      int    `json:"usegem"`
	Addgold     int    `json:"addgold"`
	Double      int    `json:"double"`
	Getfreegold int    `json:"getfreegold"`
}

//!
type S2C_GetSP struct {
	Cid string `json:"cid"`
	Ret int    `json:"ret"`
	Sp  string `json:"sp"`
}

//!
type S2C_NpcAward struct {
	Cid      string     `json:"cid"`
	NpcAward int        `json:"npcaward"`
	Item     []PassItem `json:"item"`
}

//! 更新技能点
type S2C_UpdateSP struct {
	Cid  string `json:"cid"`
	Sp   int    `json:"sp"`
	Time int    `json:"time"`
	Uid  int64  `json:"uid"`
}

//! 更新经验
type S2C_UpdateExp struct {
	Cid      string     `json:"cid"`
	Old      int        `json:"old"`
	New      int        `json:"new"`
	Newexp   int        `json:"newexp"`
	GetItems []PassItem `json:"getitems"`
}

//! 发送错误信息
type S2C_ErrInfo struct {
	Cid  string `json:"cid"`
	Info string `json:"info"`
}

//! 发送抽奖消息
type S2C_FindOK struct {
	Cid         string     `json:"cid"`
	FindType    int        `json:"findtype"`
	FreeTime    int64      `json:"freetime"`    //! 免费倒计时
	HgFreeTime  int64      `json:"hgfreetime"`  //! 后宫免费倒计时
	BoxFreeTime int64      `json:"boxfreetime"` //! 后宫免费倒计时
	Num         int        `json:"num"`         //! 神将招募多少次
	MJNum       int        `json:"mjnum"`       //! 名将招募次数
	GoldNum     int        `json:"goldnum"`
	GemNum      int        `json:"gemnum"`
	Item        []PassItem `json:"item"`
	Cost        []PassItem `json:"cost"`
	Gem         int        `json:"gem"`
	BoxNum      int        `json:"boxnum"`
	BeautyNum   int        `json:"hgnum"`
	Gold        int        `json:"gold"`
	YuFu        int        `json:"yufu"`
	GoldEndTime int64      `json:"goldendtime"`  // 免费道具抽结束时间
	LootEnergy  int        `json:"loot_energy"`  // 高级召唤能量
	SummonTimes int        `json:"summon_times"` // 高级抽奖次数
}

type S2C_FindPool struct {
	Cid             string               `json:"cid"`
	FindType        int                  `json:"findtype"`
	FindNum         int                  `json:"findnum"`
	FindNumToday    int                  `json:"findnumtoday"`
	Item            []PassItem           `json:"item"`         //掉落英雄，仅用来显示
	CostItems       []PassItem           `json:"costitems"`    //召唤消耗
	GetItems        []PassItem           `json:"getitems"`     //获得的物品
	GetItemsTran    []PassItem           `json:"getitemstran"` //分解获得
	RewardInfo      *FindRewardInfo      `json:"rewardinfo"`   //奖励进度更新
	TipNum          int                  `json:"tipnum"`       //距离保底
	FreeNextTime    int64                `json:"freenexttime"`
	FindTimes       int                  `json:"findtimes"`
	LuckyFindRecord []*LotteryDrawRecord `json:"luckyfindrecord"`
	GemFindProcess  int					 `json:"gemfindprocess"`	// 累计招募进度
}

type S2C_FindAstrology struct {
	Cid       string     `json:"cid"`
	FindNum   int        `json:"findnum"`
	Item      []PassItem `json:"item"`      //掉落
	CostItems []PassItem `json:"costitems"` //召唤消耗
	GetItems  []PassItem `json:"getitems"`  //获得的物品
}

type S2C_FindSelfSelection struct {
	Cid          string          `json:"cid"`
	FindNum      int             `json:"findnum"`
	Item         []PassItem      `json:"item"`         //掉落
	CostItems    []PassItem      `json:"costitems"`    //召唤消耗
	GetItems     []PassItem      `json:"getitems"`     //获得的物品
	GetTimes     int             `json:"gettimes"`     //次数
	GetItemsTran []PassItem      `json:"getitemstran"` //分解获得
	RewardInfo   *FindRewardInfo `json:"rewardinfo"`   //
}

type S2C_GetSelfSelection struct {
	Cid           string         `json:"cid"`
	SelfSelection *SelfSelection `json:"selfselection"`
}

type S2C_FindSaveWish struct {
	Cid      string        `json:"cid"`
	WishInfo *FindWishInfo `json:"wishinfo"`
}

type S2C_FindOpenCamp struct {
	Cid       string     `json:"cid"`
	FindPool  *FindPool  `json:"findpool"`
	CostItems []PassItem `json:"costitems"`
}

type S2C_UpdateInterStellarPos struct {
	Cid          string `json:"cid"`
	Pos          int    `json:"pos"`          //位置 万分比
	StellarCount int    `json:"stellarcount"` //当前解锁数量
}

type S2C_AstrologyHero struct {
	Cid       string         `json:"cid"`
	Astrology *FindAstrology `json:"astrology"`
}

type S2C_SelfSelectionHero struct {
	Cid           string         `json:"cid"`
	SelfSelection *SelfSelection `json:"selfselection"`
}

type S2C_BuyYuFu struct {
	Cid  string `json:"cid"`
	YuFu int    `json:"yufu"`
	Gem  int    `json:"gem"`
}

//! 合成英雄
type S2C_Synthesis struct {
	Cid  string     `json:"cid"`
	Item []PassItem `json:"item"`
	Hero *Hero      `json:"hero"`
}

type S2C_Interchange struct {
	Cid  string     `json:"cid"`
	Ret  int        `json:"ret"`
	Item []PassItem `json:"item"`
}

//!
type S2C_OnItem struct {
	Cid     string     `json:"cid"`
	Itemlst []PassItem `json:"itemlst"`
}

//!
type S2C_BuyItem struct {
	Cid     string     `json:"cid"`
	Ret     int        `json:"ret"`
	Itemlst []PassItem `json:"itemlst"`
}

type S2C_MergeItem struct {
	Cid  string     `json:"cid"`
	Item []PassItem `json:"item"`
}

//! 读邮件
type S2C_ReadMail struct {
	Cid   string `json:"cid"`
	Index int64  `json:"index"`
}

//! 签到奖励
type S2C_CheckinAward struct {
	Cid      string     `json:"cid"`
	Item     []PassItem `json:"item"`
	Newvalue int        `json:"newvalue"`
	Ret      int        `json:"ret"`
}


//! 签到奖励展示
type S2C_CheckinAwardInfo struct {
	Cid      string     `json:"cid"`
	Items 	[]PassItem  `json:"checkiniteminfo"`	// 签到奖励
	CheckinNum int		`json:"checkinnum"`		// 累计签到天数
	CheckinState int	`json:"checkinstate"`	// 签到状态
}

//! 签到奖励领取
type S2C_CheckinToday struct {
	Cid      string     `json:"cid"`
	Item 	PassItem `json:"item"`	// 签到奖励
	CheckinNum int		`json:"checkinnum"`		// 累计签到天数
	CheckinState int 	`json:"checkistate"`	// 签到状态
}

//! 首充奖励
type S2C_GetFirsetAward struct {
	Cid   string           `json:"cid"`
	Award []PassItem       `json:"award"`
	Info  San_UserRecharge `json:"info"`
	Ret   int              `json:"ret"`
	Type  int              `json:"type"`
}

//! VIP每日奖励
type S2C_GetVipDaily struct {
	Cid           string     `json:"cid"`
	Award         []PassItem `json:"award"`
	VipDailyState int64      `json:"vipdailystate"`
}

//! VIP每周购买
type S2C_BuyVipWeek struct {
	Cid             string     `json:"cid"`
	Award           []PassItem `json:"award"`
	Cost            []PassItem `json:"cost"`
	BuyVipWeekState int64      `json:"buyvipweekstate"`
}

type S2C_GetWarOrderReward struct {
	Cid          string          `json:"cid"`
	GetItem      []PassItem      `json:"getitem"`
	WarOrderTask JS_WarOrderTask `json:"warordertask"`
}

type S2C_GetWarOrderLimitReward struct {
	Cid          string          `json:"cid"`
	GetItem      []PassItem      `json:"getitem"`
	WarOrderTask JS_WarOrderTask `json:"warordertask"`
}

type S2C_GetWarOrderBuy struct {
	Cid      string     `json:"cid"`
	GetItem  []PassItem `json:"getitem"`
	CostItem []PassItem `json:"costitem"`
}

//! 购买基金
type S2C_BuyFund struct {
	Cid      string     `json:"cid"`
	Cost     []PassItem `json:"cost"`
	FundType int64      `json:"fundtype"`
}

//! 领取基金奖励
type S2C_GetFundAward struct {
	Cid   string     `json:"cid"`
	Award []PassItem `json:"award"`
	Index int        `json:"index"`
	Value int64      `json:"value"`
}

//! 购买VIP特权礼包
type S2C_BuyVipBox struct {
	Cid       string     `json:"cid"`
	Award     []PassItem `json:"award"`
	Cost      []PassItem `json:"cost"`
	CurVipBox int64      `json:"curvipbox"`
}

//! 领取福利奖励
type S2C_GetWeekPlanAward struct {
	Cid       string        `json:"cid"`
	Info      []JS_TaskInfo `json:"info"`
	Ret       int           `json:"ret"`
	Param     int           `json:"param"`
	Point     int           `json:"point"`
	GetItem   []PassItem    `json:"getitem"`
	CostItem  []PassItem    `json:"costitem"`
	IsGetMark map[int]int   `json:"isgetmark"`
}

//赏金社任务同步
type S2C_WarOrderTaskUpdate struct {
	Cid      string            `json:"cid"`
	Info     []JS_WarOrderTask `json:"info"`
	WarOrder JS_WarOrder       `json:"warorder"`
}

type S2C_WarOrderLimitInfo struct {
	Cid           string             `json:"cid"`
	WarOrderLimit []JS_WarOrderLimit `json:"warorderlimit"`
}

//! 完成任务
type S2C_TaskFinish struct {
	Cid      string             `json:"cid"`
	Item     []PassItem         `json:"item"`
	Info     []JS_TaskInfo      `json:"info"`
	Ret      int                `json:"ret"`
	Tasktype int                `json:"tasktype"`
	Liveness []*JS_LivenessInfo `json:"liveness"` // 活跃度
}

//! 完成目标任务
type S2C_TargetTaskFinish struct {
	Cid        string              `json:"cid"`
	TaskId     int                 `json:"taskid"`
	Item       []PassItem          `json:"item"`
	Info       []JS_TargetTaskInfo `json:"info"`
	SystemInfo map[int]int         `json:"systeminfo"`
}

type S2C_GetTargetLvReward struct {
	Cid      string              `json:"cid"`
	Item     []PassItem          `json:"item"`
	Info     []JS_TargetTaskInfo `json:"info"`
	SystemId int                 `json:"systemid"`
}

//! 领取活跃
type S2C_TaskLiveness struct {
	Type  int        `json:"type"`
	Cid   string     `json:"cid"`
	Item  []PassItem `json:"item"`
	Award int        `json:"award"`
}

type S2C_Getlosehero struct {
	Cid  string     `json:"cid"`
	Ret  int        `json:"ret"`
	Item []PassItem `json:"item"`
}

//! 同步行动力
type S2C_TeamPower struct {
	Cid       string `json:"cid"`
	Power     int    `json:"power"`
	PowerTime int64  `json:"powertime"`
}

//! 设置队伍位置
type S2C_CampFightTeam struct {
	Cid    string `json:"cid"`
	Index  int    `json:"index"`
	State  int    `json:"state"`
	Cityid int    `json:"cityid"`
	Finish int    `json:"finish"`
}

//! 移动部队
type S2C_CampFightCanEnter struct {
	Cid   string `json:"cid"`
	Index []int  `json:"index"`
}

//! 国战玩法通知
type S2C_CampFightPlayNotice struct {
	Cid       string `json:"cid"`
	Uid       int64  `json:"uid"`
	Name      string `json:"name"`
	Camp      int    `json:"camp"`
	WarPlayId int    `json:"warplayid"`
	Data      int    `json:"data"`
}

//! 鼓舞广播
type S2C_CampFightInspireNotice struct {
	Cid   string `json:"cid"`
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Camp  int    `json:"camp"`
	Level int    `json:"level"`
}

//! 消除溃败
type S2C_CampFightClearLose struct {
	Cid    string     `json:"cid"`
	Ret    int        `json:"ret"`
	Cityid int        `json:cityid`
	Cost   []PassItem `json:"cost"`
}

//! 五虎争雄状态
type S2C_CampFight55Info struct {
	Cid       string `json:"cid"` //
	WarPlayid int    `json:"warplayid"`
	State     int    `json:"state"`
	Index     int    `json:"index"`
	Camp      [2]int `json:"camp"`
	Win       [2]int `json:"win"`
	WinPlayer [2]int `json:"winplayer"`
	Rankid    int    `json:"rankid"`
	OpenTime  int64  `json:"opentime"`
	Condition int    `json:"condition"`
}

//! 五虎争雄初始化信息
type S2C_CampFight55Init struct {
	Cid       string                  `json:"cid"`          //! id
	WarPlayid int                     `json:"warplayid"`    //! 战斗玩法
	State     int                     `json:"state"`        //! 当前状态
	Rankid    int                     `json:"rankid"`       //! 队列排名
	Index     int                     `json:"index"`        //! 当前战斗序列，state = 3
	Team      [2][]*Son_CampFightInfo `json:"team"`         //! 对战信息
	TeamNum   [2]int                  `json:"teamnum"`      //! 报名人数
	HP        [2][]int                `json:"hp"`           //! 剩余HP
	InspireLv [2][]int                `json:"inspirelevel"` //! 鼓舞等级
	Result    []int                   `json:"result"`       //! 战斗结果，0-还未出结果，1左胜利，2右胜利
	TotalTime int                     `json:"totaltime"`    //! 当前战斗结束时间
	EndTime   int64                   `json:"endtime"`      //! 结束时间
}

//! 五雄争霸战斗-序列
type S2C_CampFight55Run struct {
	Cid       string `json:"cid"`
	Index     int    `json:"index"`     //! 战斗学列
	WarPlayid int    `json:"warplayid"` //! 战斗玩法
	Result    int    `json:"result"`    //! 战斗结果，0开始战斗，1-攻方胜利，2-守方胜利
	Hp        [2]int `json:"hp"`        //! 剩余血量
	Time      int64  `json:"time"`      //! 战斗到期时间
}

//! 五虎争雄战斗结束
type S2C_CampFight55End struct {
	Cid          string `json:"cid"`
	WarPlayid    int    `json:"warplayid"` //! 战斗玩法
	Win          [2]int `json:"win"`
	WinPlayer    [2]int `json:"winplayer"` //!胜利剩余人数，50v50专用
	Honor        int    `json:"honor"`
	Occupy       [2]int `json:"occupy"`
	State        int    `json:"state"`
	NextOpenTime int64  `json:"nextopentime"`
}

//! 五虎争雄状态
type S2C_CampFight56Info struct {
	Cid       string `json:"cid"` //
	WarPlayid int    `json:"warplayid"`
	State     int    `json:"state"`
	Index     int    `json:"index"`
	Camp      [2]int `json:"camp"`
	Win       [2]int `json:"win"`
	WinPlayer [2]int `json:"winplayer"`
	Rankid    int    `json:"rankid"`
	OpenTime  int64  `json:"opentime"`
	Condition int    `json:"condition"`
	CityId    int    `json:"cityid"`
}

//! 五虎争雄初始化信息
type S2C_CampFight56Init struct {
	Cid         string                 `json:"cid"`       //! id
	WarPlayid   int                    `json:"warplayid"` //! 战斗玩法
	State       int                    `json:"state"`     //! 当前状态
	Rankid      int                    `json:"rankid"`    //! 队列排名
	Index       int                    `json:"index"`     //! 当前战斗序列，state = 3
	IndexTeam   [5][2]int              `json:"indexteam"` //! 组内战斗
	IndexMax    [5][2]int              `json:"indexmax"`  //! 队列人数
	Team        [2][]*Son_CampFightReg `json:"team"`      //! 对战信息
	TeamNum     [2]int                 `json:"teamnum"`   //! 报名人数
	FightTeam   [2]*JS_FightBase       `json:"fightteam"` //! 战斗队列
	HP          [5][2][5]int           `json:"hp"`        //! 剩余HP
	Result      [5]int                 `json:"result"`    //! 战斗结果，0-还未出结果，1左胜利，2右胜利
	TotalTime   int                    `json:"totaltime"` //! 当前战斗结束时间
	EndTime     int64                  `json:"endtime"`   //! 结束时间
	MyTeam      int                    `json:"myteam"`    //! 自己所在队伍，没有进入队列为0
	EnterPlayer int64                  `json:"enterplayer"`
}

type S2C_CampFight56Upd struct {
	Cid       string `json:"cid"`    //! 消息ID
	CityId    int    `json:"cityid"` //! 城池ID
	ReqPlayer [2]int `json:"req"`    //! 报名人数
}

type Son_CampFightReg struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	Icon    int    `json:"icon"`
	Fight   int64  `json:"fight"`
	Level   int    `json:"level"`
	Camp    int    `json:"camp"`
	Kill    int    `json:"kill"`
	Class   int    `json:"class"`
	Inspire int    `json:"inspire"`
}

//! 五雄争霸战斗-序列
type S2C_CampFight56Run struct {
	Cid        string           `json:"cid"`
	Index      int              `json:"index"`      //! 队伍序列
	IndexTeam  [2]int           `json:"indexteam"`  //! 队列序列
	IndexMax   [2]int           `json:"indexmax"`   //! 队列上限
	WarPlayid  int              `json:"warplayid"`  //! 战斗玩法
	Result     int              `json:"result"`     //! 战斗结果，0开始战斗，1-攻方胜利，2-守方胜利
	Win        [2]int           `json:"win"`        //! 胜负
	Kill       [2]int           `json:"kill"`       //! 连斩数
	TeamResult [5]int           `json:"teamresult"` //! 队伍输赢
	Hp         [2][5]int        `json:"hp"`         //! 剩余血量
	BossHP     [2]int           `json:"bosshp"`     //! 巨兽的HP
	Time       int64            `json:"time"`       //! 战斗到期时间
	FightTeam  [2]*JS_FightBase `json:"fightteam"`  //! 战斗数据
	CityId     int              `json:"cityid"`
}

//! 五虎争雄战斗结束
type S2C_CampFight56End struct {
	Cid          string           `json:"cid"`
	WarPlayid    int              `json:"warplayid"` //! 战斗玩法
	Win          [2]int           `json:"win"`
	TeamResult   [5]int           `json:"teamresult"`
	WinPlayer    [2]int           `json:"winplayer"` //!胜利剩余人数，50v50专用
	Honor        int              `json:"honor"`
	Occupy       [2]int           `json:"occupy"`
	State        int              `json:"state"`
	MaxKill      [2]int           `json:"maxkill"`
	MVP          [2]int64         `json:"mvp"`
	NextOpenTime int64            `json:"nextopentime"`
	FightTeam    [2]*JS_FightBase `json:"fightteam"`
	CityId       int              `json:"cityid"`
}

type S2C_CampFight56Giveway struct {
	Cid string `json:"cid"`
	Uid int64  `json:"uid"`
}

//! 改变战报
type S2C_UpdFightRecord struct {
	Cid    string `json:"cid"`
	Id     int64  `json:"id"`
	Result int    `json:"result"`
}

//! 移动部队
type S2C_MoveTeamBegin struct {
	Cid    string `json:"cid"`
	Index  int    `json:"index"`
	Cityid int    `json:"cityid"`
	State  int    `json:"state"`
}

//!
type S2C_GetFightTeamInfo struct {
	Cid  string       `json:"cid"`
	Info JS_FightInfo `json:"info"`
}

//! 阵营战结果
type S2C_CampFightResult struct {
	Cid    string `json:"cid"`
	Result int    `json:"result"` //! 1攻方胜利 2守方胜利
	Kill   int    `json:"kill"`   //! 杀敌数
	Solo   int    `json:"solo"`   //! 诱敌数
	Help   int    `json:"help"`   //! 召唤援军数
	Cityid int    `json:"cityid"` //! 移动城市
}

type Son_CampFightInfo struct {
	Uid      int64  `json:"uid"`
	Index    int    `json:"index"`
	Name     string `json:"name"`
	Icon     int    `json:"icon"`
	Fight    int64  `json:"fight"`
	Level    int    `json:"level"`
	Camp     int    `json:"camp"`
	Kill     int    `json:"kill"`
	Elite    int    `json:"elite"`
	Class    int    `json:"class"`
	ArmsType int    `json:"armstype"`
	Honor    int    `json:"honor"`
	Buffer   int    `json:"buffer"`
	Inspire  int    `json:"inspire"`
}

//! 加删队伍
type S2C_CampFightAdd struct {
	Cid  string            `json:"cid"`
	Info Son_CampFightInfo `json:"info"`
}

type S2C_CampFightDel struct {
	Cid   string `json:"cid"`
	Uid   int64  `json:"uid"`
	Index int    `json:"index"`
}

//! 国战加等待的人
type S2C_CampFightWait struct {
	Cid  string        `json:"cid"`
	Pos  int           `json:"pos"`
	Wait *JS_FightInfo `json:"wait"` //! 上阵方
}

//! 国战等待队列
type S2C_CampFightWaitList struct {
	Cid    string               `json:"cid"`
	Pos    int                  `json:"pos"`  //! 0-攻击方，1-防守方
	Page   int                  `json:"page"` //! 请求页码
	AttNum int                  `json:"attnum"`
	DefNum int                  `json:"defnum"`
	Wait   []*Son_CampFightInfo `json:"wait"` //! 等待队列
}

//! 一骑讨触发玩法
type S2C_CampFightSoloPlay struct {
	Cid      string `json:"cid"`
	PlayList []int  `json:"playlist"` //一骑讨触发玩法
}

//! 国战solo
type S2C_CampFightSolo2Begin struct {
	Cid     string          `json:"cid"`
	SoloId  int             `json:"soloid"`
	Buffer  int             `json:"buffer"`
	Inspire [2]int          `json:"inspire"`
	Info    [2]JS_FightInfo `json:"info"`
}

//! 公告
type S2C_Notice struct {
	Cid     string `json:"cid"`
	Type    int    `json:"type"`
	Content string `json:"content"`
}

//! 聊天
type S2C_Chat struct {
	Cid      string `json:"cid"`
	Uid      int64  `json:"uid"`
	Channel  int32  `json:"channel"`
	Name     string `json:"name"`
	Icon     int    `json:"icon"`
	Portrait int    `json:"portrait"` //! 头像框
	Vip      int    `json:"vip"`
	Level    int    `json:"level"`
	Time     int64  `json:"time"`
	Content  string `json:"content"`
	Url      string `json:"url"`
	Camp     int    `json:"camp"`
	Param    int    `json:"param"`
	TeamId   int64  `json:"teamid"`
	MineId   int    `json:"mine_id"`
	BuffCd   int64  `json:"buffcd"`
}

type S2C_NewChat struct {
	Cid     string         `json:"cid"`
	Channel int32          `json:"channel"`
	MsgList []*ChatMessage `json:"msglist"` //! 消息列表
}

type S2C_NewChatGap struct {
	Cid     string `json:"cid"`
	Channel int    `json:"channel"`
	GapUid  int64  `json:"gapuid"`
}

//! 加一个好友消息
type S2C_AddFriendMsg struct {
	Cid  string     `json:"cid"`
	Info *JS_Friend `json:"info"`
}

//! 删除好友消息
type S2C_DelFriendMsg struct {
	Cid string `json:"cid"`
	Uid int64  `json:"uid"`
}

//! 好友上线下线
type S2C_FriendOnline struct {
	Cid    string `json:"cid"`
	Uid    int64  `json:"uid"`
	Online int    `json:"online"`
}

//! 好友消息
type S2C_Friend struct {
	Cid         string              `json:"cid"`
	Friend      []*JS_FriendNode    `json:"friend,omitempty"`
	Apply       []*JS_Friend        `json:"apply,omitempty"`
	Black       []*JS_Friend        `json:"black,omitempty"`
	Count       int                 `json:"count"`
	ApplyHire   []*HireHero         `json:"applyhire"`   //自己的申请
	HireHero    []*FriendHero       `json:"hirehero"`    //租到的英雄
	HireTime    int64               `json:"hiretime"`    //过期时间
	HireList    map[int][]*HireHero `json:"hirelist"`    //!  可租借的别人的  key  英雄ID
	SelfList    map[int]*HireHero   `json:"selflist"`    //!  自己的,从这个里面可以看到谁申请过  key  KEYID
	HeroSetInfo []*NewHero          `json:"herosetinfo"` //个人信息里的hero
	HireUseSign []int               `json:"hireusesign"` //使用标记
	GiftSign    []int64             `json:"giftsign"`    //今天历史赠送记录
}

type S2C_FriendUpdate struct {
	Cid    string           `json:"cid"`
	Friend []*JS_FriendNode `json:"friend,omitempty"`
}

type S2C_SetUseSign struct {
	Cid         string `json:"cid"`
	HireUseSign []int  `json:"hireusesign"` //使用标记
}

type S2C_UpdateHireList struct {
	Cid      string              `json:"cid"`
	HireList map[int][]*HireHero `json:"hirelist"` //!  可租借的别人的  key  英雄ID
}

//! 好友推荐
type S2C_FriendCommend struct {
	Cid     string       `json:"cid"`
	Commend []*JS_Friend `json:"commend"`
}

//! 查看
type S2C_Look struct {
	Cid       string     `json:"cid"`
	Uid       int64      `json:"uid"`
	Name      string     `json:"name"`
	Icon      int        `json:"icon"`
	Vip       int        `json:"vip"`
	Level     int        `json:"level"`
	Camp      int        `json:"camp"`
	Party     string     `json:"party"`
	UnionID   int        `json:"unionid"`
	Time      int64      `json:"time"`
	Fight     int64      `json:"fight"`
	Office    int        `json:"office"`
	BraveHand int        `json:"bravehand"`
	Hero      []*NewHero `json:"hero"` //英雄详情
	Portrait  int        `json:"portrait"`
	Stage     int        `json:"stage"`     //关卡进度
	Server    int        `json:"server"`    //服务器
	Signature string     `json:"signature"` //签名
	//TreeLevel int            `json:"treelevel"`
	//TreeInfo  []*JS_LifeTree `json:"treeinfo"`
}
type Son_Look struct {
	Heroid int `json:"heroid"`
	Color  int `json:"color"`
	Level  int `json:"level"`
	Star   int `json:"star"`
	Talent int `json:"talent"`
}

//! 好友体力
type S2C_FriendPower struct {
	Cid   string     `json:"cid"`
	Uid   []int64    `json:"uid"`
	Type  int        `json:"type"` //! 0赠送 1领取
	Value int        `json:"value"`
	Count int        `json:"count"`
	Item  []PassItem `json:"award"`
}

type S2C_SetCamp struct {
	Cid  string `json:"cid"`
	Camp int    `json:"camp"`
	City int    `json:"city"`
}

type S2C_GetCamp struct {
	Cid  string `json:"cid"`
	Camp int    `json:"camp"`
	Num  [3]int `json:"num"`
}

type S2C_CreateUnion struct {
	Cid            string           `json:"cid"`
	Ret            int              `json:"ret"`
	Info           JS_UserUnionInfo `json:"info"`
	Money          []PassItem       `json:"money"`
	Unioninfo      JS_Union         `json:"unioninfo"`
	CopyUpdateTime int64            `json:"copyupdatetime"`
}

type S2C_GetUnionList struct {
	Cid  string      `json:"cid"`
	Ret  int         `json:"ret"`
	List []JS_Union2 `json:"list"`
}

// 获得玩家军团信息
type S2C_GetUserUnionInfo struct {
	Cid            string           `json:"cid"`
	Ret            int              `json:"ret"`
	CopyUpdateTime int64            `json:"copyupdatetime"`
	Selfinfo       JS_UserUnionInfo `json:"selfinfo"`
	Unioninfo      JS_Union         `json:"unioninfo"`
	Unionlist      []JS_Union2      `json:"unionlist"`
	ChangeMaster   bool             `json:"changemaster"`
	OldMaster      string           `json:"oldmaster"`
}

//修改军团名字
type S2C_AlertUnionName struct {
	Cid   string     `json:"cid"`
	Ret   int        `json:"ret"`
	Money []PassItem `json:"money"`
}
type S2C_GetUnionInfo struct {
	Cid       string    `json:"cid"`
	Ret       int       `json:"ret"`
	Unioninfo *JS_Union `json:"unioninfo"`
}

type S2C_GetUnionRecord struct {
	Cid    string           `json:"cid"`
	Ret    int              `json:"ret"`
	Record []JS_UnionRecord `json:"record"`
}

type S2C_JoinUnion struct {
	Cid      string           `json:"cid"`
	Ret      int              `json:"ret"`
	Selfinfo JS_UserUnionInfo `json:"selfinfo"`
	Info     JS_Union         `json:"info"`
}

type S2C_FindUnion struct {
	Cid  string      `json:"cid"`
	Ret  int         `json:"ret"`
	Info []JS_Union2 `json:"info"`
}

type S2C_MemberInfo struct {
	Cid     string         `json:"cid"`
	Herolst []JS_UnionHero `json:"herolst"`
}

type S2C_UnionModify struct {
	Cid     string `json:"cid"`
	Uid     int64  `json:"uid"`
	Destuid int64  `json:"destuid"`
	Op      int    `json:"op"`
}

// 申请成功
type S2C_MasterAllok struct {
	Cid       string       `json:"cid"`
	Unionid   int          `json:"unionid"`
	AddPlayer []*JS_Member `json:"addplayer"`
}

type S2C_UnionSendMail struct {
	Cid string `json:"cid"`
	Ret int    `json:"ret"`
}

type S2C_SetBraveHand struct {
	Cid     string `json:"cid"`
	Uid     int64  `json:"uid"`
	Destuid int64  `json:"destuid"`
	Op      int    `json:"op"`
}

type S2C_UnionDonation struct {
	Cid       string     `json:"cid"`
	Type      int        `json:"type"`
	Value     int        `json:"value"`
	Donation  int        `json:"donation"`
	Givecount int        `json:"givecount"`
	Exp       int        `json:"exp"`
	DayExp    int        `json:"dayexp"`
	Level     int        `json:"level"`
	Item      []PassItem `json:"item"`
}

// 开始公会狩猎战斗
type S2C_StartHuntFight struct {
	Cid     string     `json:"cid"`
	AddItem []PassItem `json:"additem"`
}

// 会长开启公会狩猎
type S2C_OpenHuntFight struct {
	Cid  string        `json:"cid"`
	Type int           `json:"type"`
	Info *JS_UnionHunt `json:"info"`
}

type S2C_EndHuntFight struct {
	Cid               string         `json:"cid"`
	Type              int            `json:"type"`
	Damage            int64          `json:"damage"`
	Info              *UserHuntLimit `json:"info"`
	Item1             []PassItem     `json:"item1"`
	Item2             []PassItem     `json:"item2"`
	GetPrivilegeItems []PassItem     `json:"getprivilegeitems"` //特权物品 新增
}

type S2C_SweepHuntFight struct {
	Cid               string         `json:"cid"`
	Type              int            `json:"type"`
	Damage            int64          `json:"damage"`
	Info              *UserHuntLimit `json:"info"`
	Item1             []PassItem     `json:"item1"`
	Item2             []PassItem     `json:"item2"`
	GetPrivilegeItems []PassItem     `json:"getprivilegeitems"` //特权物品 新增
}

//! dps排行
type S2C_HuntDpsTop struct {
	Cid    string                `json:"cid"`
	Type   int                   `json:"type"`
	DpsTop []*JS_UnionHuntDpsTop `json:"dpstop"`
}

type S2C_GetHuntInfo struct {
	Cid           string           `json:"cid"`
	UserHuntLimit []*UserHuntLimit `json:"userhuntlimit"` //! 狩猎限制
	GuildHunting  []int            `json:"guild_hunting"`
}

type S2C_UnionCopyInfo struct {
	Cid        string              `json:"cid"`
	UpdateTime int64               `json:"updatetime"`
	Info       []Son_UnionCopyInfo `json:"info"`
	CopyAward  []int               `json:"copyaward"`
}
type Son_UnionCopyInfo struct {
	Id       int `json:"id"`
	Progress int `json:"progress"`
}

type S2C_UnionCopyBegin struct {
	Cid     string                `json:"cid"`
	Id      int                   `json:"id"`
	Monster []JS_UnionCopyMonster `json:"monster"`
}

type S2C_UnionCopyEnd struct {
	Cid      string     `json:"cid"`
	Id       int        `json:"id"`
	Progress int        `json:"progress"`
	Item     []PassItem `json:"item"`
}

type S2C_UnionCopyAward struct {
	Cid  string     `json:"cid"`
	Id   int        `json:"id"`
	Item []PassItem `json:"item"`
}

type S2C_UnionCopyTop1 struct {
	Cid  string             `json:"cid"`
	Info lstUnionCopyNumTop `json:"info"`
}

type S2C_UnionCopyTop2 struct {
	Cid  string             `json:"cid"`
	Info lstUnionCopyDpsTop `json:"info"`
}

//! 天赋星级
type S2C_Top struct {
	Cid  string       `json:"cid"`
	Type int          `json:"type"`
	Top  []*Js_ActTop `json:"top"`
	Ver  int          `json:"ver"`
	Cur  int          `json:"cur"`
	Old  int          `json:"old"`
	Num  int64        `json:"num"`
}

// 得到国家宝箱
type S2C_CountryBox struct {
	Cid  string     `json:"cid"`
	Item []PassItem `json:"item"`
	Box  []int      `json:"box"`
}

//调试消息
type S2C_DebugString struct {
	Cid      string `json:"cid"`
	Debugstr string `json:"debugstr"`
}

//! 获取外交数据
type S2C_DiplomacyInfo struct {
	Cid        string   `json:"cid"`        //! 消息ID
	Index      [3]int   `json:"index"`      //! 国战进程
	Align      [3]int   `json:"align"`      //! 联盟状态
	AttackCity [3][]int `json:"attackcity"` //! 攻击城池

	Time      int      `json:"time"`      //! 开战时间
	NextTime  int64    `json:"nexttime"`  //!  每4天一次，开服3天第一场
	FightTime [3]int   `json:"fighttime"` //! 开始战斗时间
	Status    [3]int   `json:"status"`    //! 当前状态
	CityArr   []int    `json:"cityarr"`   //! NPC城池序列
	UnionArr  [4]int64 `json:"unionarr"`  //! 军团序列
}

//! 国战信息
type JS_CampFightInfo struct {
	CityId      int   `json:"cityid"`      //! 城池ID
	StartTime   int64 `json:"starttime"`   //! 开战时间
	BattleTime  int64 `json:"battletime"`  //! 决战时间
	Rank        int   `json:"rank"`        //! 报名排名，是否入选，-1-未报名，0未入选(超过35也是未入选)，
	Status      int   `json:"status"`      //! 状态，0-等待开战，1-报名，决战，2-结束
	FightStatus int   `json:"fightstatus"` //! 国战状态，1-等待，2-报名，3-决战开启，4-决战结束
	RecordId    int   `json:"record"`      //! 战报Id
	Camp        int   `json:"camp"`        //! 目标阵营
	Kill        int   `json:"kill"`        //! 击杀数
}

//! 国家任务
type S2C_CampTaskChg struct {
	Cid      string `json:"cid"`
	Progress int    `json:"progress"` //! 进度
	Done     int    `json:"done"`
}

type S2C_TitleTask struct {
	Cid  string `json:"cid"`
	Task [3]int `json:"task"`
}

//! 获得竞技场信息
type S2C_PvPInfo struct {
	Cid  string      `json:"cid"`
	Info Son_PvPInfo `json:"info"`
}
type Son_PvPInfo struct {
	Rankid     int    `json:"rankid"`     //! 排名
	Uid        int64  `json:"uid"`        //! uid
	Name       string `json:"name"`       //! uname
	Point      int    `json:"point"`      //! 积分
	Format     []int  `json:"format"`     //! 出战阵容
	Num        int    `json:"num"`        //! 已战场次
	Buynum     int    `json:"buynum"`     //! 购买场次
	Best       int    `json:"best"`       //! 最好排名
	Award      int64  `json:"award"`      //! 奖励状态
	Worship    int    `json:"worship"`    //! 膜拜状态
	Time       int64  `json:"time"`       //! 下一场倒计时时间
	PointAward int64  `json:"pointaward"` //! 积分奖励状态
	Vip        int    `json:"vip"`        //! vip等级
	TimeNum    int    `json:"timenum"`    //! 清除CD时间次数
}

//! 得到天下会武对手
type S2C_GetPvPFight struct {
	Cid    string          `json:"cid"`
	RankId int             `json:"rankid"`
	Enemy  []*JS_FightInfo `json:"enemy"`
	Top    []*JS_FightBase `json:"top'`
}

//type S2C_GetPvPEnd struct {
//	Cid    string    `json:"cid"`
//	Arms   int       `json:"arms"`
//	Rank   int       `json:"rank"`
//	Point  int       `json:"point"`
//	Result int       `json:"result"`
//	Fight  *PvpFight `json:"fight"`
//}
//
//// 天下会武排行榜
//type S2C_PvPTop struct {
//	Cid string       `json:"cid"`
//	Top []*JS_TopJJC `json:"top"`
//	Ver int          `json:"ver"`
//}

//! 神将预约
type S2C_OrderHero struct {
	Cid      string `json:"cid"`
	State    int    `json:"state"`    //! 0-未预约 1-已预约
	Hero     int    `json:"hero"`     //! 武将Id
	OrderNum int    `json:"ordernum"` //! 预约人数
	Day      int    `json:"day"`      //! 当前活动天数 0~2-预约，3~7购买
	Mask     []int  `json:"mask"`     //! 购买天数
}

// 玩家点赞
type S2C_UserCountryStep struct {
	Cid   string `json:"cid"`
	Index int    `json:"index"`
	Gold  int    `json:"gold"`
}

// 国王点赞
type S2C_CountryStep struct {
	Cid    string `json:"cid"`
	Camp   int    `json:"camp"`   // 赞
	Praise int    `json:"praise"` // 赞
	Step   int    `json:"step"`   // 踩
}

//! 道具批量出售
type S2C_SellItems struct {
	Cid  string     `json:"cid"`
	Item []PassItem `json:"item"`
	Gold int        `json:"gold"`
}

// 国王争夺信息
type S2C_Kinginfo struct {
	Cid       string `json:"cid"`
	Id        int    `json:"id"`       // 蜀魏吴123
	Kinguid   int64  `json:"kinguid"`  // 国王uid，0表示无国王
	Kingname  string `json:"kingname"` // 国王名字
	Kingicon  int    `json:"kingicon"` // 国王图标
	KingDay   int    `json:"kingday"`  // 国王在位时间
	Isfight   int    `json:"isfight"`  // 战斗状态，0未开启，1开启状态，2攻胜，3守胜
	Openuid   int64  `json:"openuid"`  // 开启战斗UID，0表示没有
	Opentime  int64  `json:"opentime"` // 开启战斗事件，0表示没有
	Unionname string `json:"unionname"`
	Face      int    `json:"face"` //性别，用于半身像   20190510 by zy
}

//! 战斗信息
type S2C_KingFightBegin struct {
	Cid  string         `json:"cid"`
	Ret  int            `json:"ret"`
	Info []JS_FightInfo `json:"info"`
}

type S2C_KingFightEnd struct {
	Cid      string     `json:"cid"`
	Time     int64      `json:"time"`
	Item     []PassItem `json:"item"`
	Nobility int        `json:"nobility"`
	Result   int        `json:"result"`
	Num      int        `json:"num"`
	Direct   int        `json:"direct"`
}

type S2C_KingAttacked struct {
	Cid        string `json:"cid"`
	NobilityId int    `json:"nobilityid"`
}

//! 领取奖励
type S2C_KingAward struct {
	Cid     string     `json:"cid"`
	Award   int64      `json:"award"`
	OutItem []PassItem `json:"outitem"`
	LastGet int64      `json:""`
}

//! 购买次数
type S2C_KingBuyNum struct {
	Cid    string `json:"cid"`
	Num    int    `json:"num"`
	Buynum int    `json:"buynum"`
	Gem    int    `json:"gem"`
}

//! 阵营战info
type S2C_CampFightInfo struct {
	Cid  string `json:"cid"`
	Sign int    `json:"sign"`
	Num  int    `json:"num"`
	Box1 []int  `json:"box1"`
	Time int64  `json:"time"`
	Box2 int    `json:"box2"`
}

//! 阵营战报名奖励
type S2C_CampFightBox1 struct {
	Cid  string     `json:"cid"`
	Id   int        `json:"id"`
	Item []PassItem `json:"item"`
}

//! 阵营战军备奖励
type S2C_CampFightBox2 struct {
	Cid  string     `json:"cid"`
	Item []PassItem `json:"item"`
	Time int64      `json:"time"`
	Box2 int        `json:"box2"`
}

//! 个人奖励
type S2C_CampOneAward struct {
	Cid  string     `json:"cid"`
	Step int        `json:"step"`
	Item []PassItem `json:"item"`
}

//! 活动
type S2C_ActivityInfo struct {
	Cid    string         `json:"cid"`
	Info   []JS_Activity  `json:"info,omitempty"`
	JJ     []int          `json:"jj,omitempty"`
	Month  []JS_MonthCard `json:"month,omitempty"`
	MsgVer string         `json:"msgver"` //! 防止服务器编译不进去
}

type S2C_ActivityUpd struct {
	Cid   string         `json:"cid"`
	Info  []JS_Activity  `json:"info"`
	Month []JS_MonthCard `json:"month"`
}

//! 活动领取
type S2C_ActivityGet struct {
	Cid        string          `json:"cid"`
	Id         int             `json:"id"`
	Item       []PassItem      `json:"item"`
	CostItem   []PassItem      `json:"costitem"`
	NextInfo   JS_LuckShopItem `json:"info"`
	NextItem   JS_ActivityBox  `json:"nextitem"`
	UpdateInfo []*JS_Activity  `json:"updateinfo"`
	GetInfo    *JS_Activity    `json:"getinfo"`
}

type S2C_TimeGiftGet struct {
	Cid      string          `json:"cid"`
	Id       int             `json:"id"`
	Item     []PassItem      `json:"item"`
	NextInfo JS_TimeGiftItem `json:"info"`
	NextItem TimeGiftConfig  `json:"nextitem"`
}

//! 预约活动报名
type S2C_OrderHeroSign struct {
	Cid  string     `json:"cid"`
	Day  int        `json:"day"`
	Ret  int        `json:"ret"`
	Item []PassItem `json:"item"`
}

//! 预约活动购买
type S2C_OrderHeroDraw struct {
	Cid  string     `json:"cid"`
	Day  int        `json:"day"`
	Ret  int        `json:"ret"`
	Item []PassItem `json:"item"`
	Mask []int      `json:"mask"`
}

//! 活动物品兑换
type S2C_ActivityExChangeItem struct {
	Cid      string     `json:"cid"`
	Id       int        `json:"id"`
	Progress int        `json:"progress"`
	Done     int        `json:"done"`
	StepMax  int        `json:"stepmax"`
	Item     []PassItem `json:"item"`
	Cost     []PassItem `json:"cost"`
}

//! 活动开启状态
type S2C_ActivityMask struct {
	Cid          string            `json:"cid"`
	Ver          int               `json:"ver"`
	Info         []JS_ActivityType `json:"info,omitempty"`
	Items        []JS_ActivityItem `json:"items,omitempty"`
	PassIdRecord int               `json:"passidrecord"`
	MsgVer       string            `json:"msgver"` //! 防止服务器编不进代码
}

//! 变强宝箱
type S2C_GetPromoteBox struct {
	Cid  string     `json:"cid"`
	Id   int        `json:"id"`
	Item []PassItem `json:"item"`
}

//! 变强信息
type S2C_GetPromoteInfo struct {
	Cid  string              `json:"cid"`
	Info []JS_BecomeStronger `json:"info"`
}

//! 领取基金
type S2C_ActivityFundAward struct {
	Cid  string     `json:"cid"`
	Ver  int        `json:"ver"`
	Pay  int        `json:"pay"`
	Day  int        `json:"day"`
	Item []PassItem `json:"item"`
}

//! 激活基金
type S2C_ActivityFundActivate struct {
	Cid string `json:"cid"`
	Ver int    `json:"ver"`
	Pay int    `json:"pay"`
}

//! 获得基金信息
type S2C_ActivityFundInfo struct {
	Cid     string            `json:"cid"`
	State   int               `json:"state"`
	EndTime int64             `json:"endtime"`
	Ver     int               `json:"ver"`
	Fund    [2]*ActivityFund  `json:"fund"`
	Items   [2][7][]*PassItem `json:"items"`
	Cost    [2]int            `json:"cost"`
	Wroth   [2]int            `json:"worth"`
}

type S2C_ActivityJJ struct {
	Cid  string     `json:"cid"`
	JJ   []int      `json:"jj"`
	Item []PassItem `json:"item"`
}

type S2C_Barrage struct {
	Cid   string `json:"cid"`
	Size  int    `json:"size"`
	Text  string `json:"text"`
	Red   int    `json:"red"`
	Green int    `json:"green"`
	Blue  int    `json:"blue"`
	Uid   int64  `json:"uid"`
}

type S2C_GetTopBox struct {
	Cid  string       `json:"cid"`
	Info []*JS_TopBox `json:"info"`
}

// 选择过关斩将
type S2C_ChoseGgzj struct {
	Cid           string `json:"cid"`
	Index         int    `json:"index"`
	Max           int    `json:"max"`
	Ggzjhard      []int  `json:"ggzjhard"`
	Ggzjlevelnum  []int  `json:"ggzjlevelnum"`
	Ggzjchosehard []int  `json:"ggzjchosehard"`
}

type S2C_ItemMax struct {
	Cid  string `json:"cid"`
	Type int    `json:"type"` //!
}

//! 过关斩将扫荡
type S2C_GgzjSweep struct {
	Cid           string     `json:"cid"`           //! cid
	Items         []PassItem `json:"items"`         //! 道具信息
	Ggzjchosehard []int      `json:"ggzjchosehard"` //! 剩余次数
}

//! 任命
type S2C_KingAppoint struct {
	Cid  string `json:"cid"`  //! cid
	Uid  int64  `json:"uid"`  //! 被任命的玩家Id
	Name string `json:"name"` //! 玩家姓名
	Icon int    `json:"icon"` //! 玩家头像
}

//! 任命
type S2C_KingCounsellor struct {
	Cid     string `json:"cid"`     //! cid
	Uid     int64  `json:"uid"`     //! 玩家Id, 如果uid=0不用显示
	Name    string `json:"name"`    //! 玩家姓名
	Icon    int    `json:"icon"`    //! 玩家头像
	EndTime int64  `json:"endtime"` //! 结束时间
}

type S2C_BuyAnchor struct {
	Cid string `json:"cid"` //! cid
	Gem int    `json:"gem"` //! 当前玩家剩余钻石
}

// 日志发送用-------------------------------------------------------------
type SendRZ_EnvInfo_AccInfo struct {
	UserType  string `json:"userType"`
	Creator   string `json:"creator"`
	AccountId string `json:"accountId"`
}

type SendRZ_EnvInfo_DevInfo struct {
	DeviceId string `json:"deviceId"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Os       string `json:"os"`
}

type SendRZ_EnvInfo_DevInfo_Ios struct {
	DeviceId string `json:"deviceId"`
	UUID     string `json:"uuid"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Os       string `json:"os"`
	Fr       string `json:"Fr"`
	Res      string `json:"res"`
	Net      string `json:"Net"`
	Mac      string `json:"mac"`
	Operator string `json:"operator"`
	Ip       string `json:"ip"`
}

type SendRZ_EnvInfo_GmInfo struct {
	PkgName  string `json:"pkgName"`
	AppVer   string `json:"appVer"`
	BuildVer string `json:"buildVer"`
	ResVer   string `json:"resVer"`
}

//渠道信息
type SendRZ_EnvInfo_ChInfo_Ios struct {
	Ch    string `json:"ch"` //Sdk来源，如大圣，火速，大禹，天拓
	SubCh string `json:"subCh"`
}

type SendRZ_envInfo_ios struct {
	AccInfo SendRZ_EnvInfo_AccInfo     `json:"accInfo"`
	DevInfo SendRZ_EnvInfo_DevInfo_Ios `json:"devInfo"`
	ChInfo  SendRZ_EnvInfo_ChInfo_Ios  `json:"chInfo"`
	RunId   string                     `json:"runId"`
}

type SendRZ_envInfo struct {
	AccInfo SendRZ_EnvInfo_AccInfo `json:"accInfo"`
	DevInfo SendRZ_EnvInfo_DevInfo `json:"devInfo"`
	GmInfo  SendRZ_EnvInfo_GmInfo  `json:"gmInfo"`
	RunId   string                 `json:"runId"`
}

//===========================================
type SendRZ_Create struct {
	Ts      string               `json:"ts"`
	AppId   string               `json:"appId"`
	Event   string               `json:"event"`
	Params  SendRZ_Create_params `json:"params"`
	EnvInfo SendRZ_envInfo       `json:"envInfo"`
}

type SendRZ_CreateIOS struct {
	Ts      string               `json:"ts"`
	AppId   string               `json:"appId"`
	Event   string               `json:"event"`
	Params  SendRZ_Create_params `json:"params"`
	EnvInfo SendRZ_envInfo_ios   `json:"envInfo"`
}

type SendRZ_Create_params struct {
	// 创建帐号用
	ServerId   string `json:"serverId"`   // 服务器ID 1
	ServerName string `json:"serverName"` // 服务器名称 测试
	RoleId     string `json:"roleId"`     // 角色ID
	RoleName   string `json:"roleName"`   // 角色名称
}

//======================================================
type SendRZ_Loginok struct {
	Ts      string                `json:"ts"`
	AppId   string                `json:"appId"`
	Event   string                `json:"event"`
	Params  SendRZ_Loginok_params `json:"params"`
	EnvInfo SendRZ_envInfo        `json:"envInfo"`
}

// 角色登录成功-参数
type SendRZ_Loginok_params struct {
	ServerId   string `json:"serverId"`   // 服务器ID
	ServerName string `json:"serverName"` // 服务器名称
	RoleId     string `json:"roleId"`     // 角色ID
	RoleName   string `json:"roleName"`   // 角色名称
	RoleLevel  int    `json:"roleLevel"`  //! 角色等级
	VipLevel   int    `json:"vipLevel"`   //! 角色vip等级
	VipExp     int    `json:"vipExp"`     //! 角色vip等级
	Force      int64  `json:"force"`      //! 战斗力
	Diamond    int    `json:"diamond"`    //! 钻石存量
	Gold       int    `json:"gold"`       //! 金币存量
	Power      int    `json:"power"`      //! 体力存量
	Feats      int    `json:"feats"`      //! 国战功勋
}

//=====================================================
type SendRZ_Offline struct {
	Ts      string                `json:"ts"`
	AppId   string                `json:"appId"`
	Event   string                `json:"event"`
	Params  SendRZ_Offline_params `json:"params"`
	EnvInfo SendRZ_envInfo        `json:"envInfo"`
}

type SendRZ_HeroInfo_params struct {
	HeroId    int `json:"heroId"`
	HeroLevel int `json:"heroLevel"`
	HeroStar  int `json:"heroStar"`
	HeroRank  int `json:"heroRank"`
}

type SendRZ_Offline_params struct {
	//! 角色退出
	ServerId   string                   `json:"serverId"`   // 服务器ID
	ServerName string                   `json:"serverName"` // 服务器名称
	RoleId     string                   `json:"roleId"`     // 角色ID
	RoleName   string                   `json:"roleName"`   // 角色名称
	RoleLevel  int                      `json:"roleLevel"`  //角色等级
	VipLevel   int                      `json:"vipLevel"`   //! 角色vip等级
	Vip_exp    int                      `json:"vipExp"`     //! 贵族经验值
	LoginTime  int64                    `json:"loginTime"`  //! 角色登录时间
	LogoutTime int64                    `json:"logoutTime"` //! 角色退出时间
	Duration   int64                    `json:"duration"`   //! 在线时长
	Force      int64                    `json:"force"`      //! 战斗力
	Diamond    int                      `json:"diamond"`    //! 钻石存量
	Gold       int                      `json:"gold"`       //! 金币存量
	Power      int                      `json:"power"`      //! 体力存量
	Feats      int                      `json:"feats"`      //! 国战功勋
	HeroInfo   []SendRZ_HeroInfo_params `json:"heroInfo"`
}

//===========================================================
type SendRZ_Levelup struct {
	Ts      string                `json:"ts"`
	AppId   string                `json:"appId"`
	Event   string                `json:"event"`
	Params  SendRZ_Levelup_params `json:"params"`
	EnvInfo SendRZ_envInfo        `json:"envInfo"`
}

type SendRZ_Levelup_params struct {
	// 角色升级成功
	ServerId   string `json:"serverId"`   // 服务器ID
	ServerName string `json:"serverName"` // 服务器名称
	RoleId     string `json:"roleId"`     // 角色ID
	RoleName   string `json:"roleName"`   // 角色名称
	RoleLevel  int    `json:"roleLevel"`  // 角色等级
	VipLevel   int    `json:"vipLevel"`   //! 角色VIP等级
}

//==============================================================
type SendRZ_GetMoney struct {
	Ts      string                 `json:"ts"`
	AppId   string                 `json:"appId"`
	Event   string                 `json:"event"`
	Params  SendRZ_GetMoney_params `json:"params"`
	EnvInfo SendRZ_envInfo         `json:"envInfo"`
}
type SendRZ_GetMoney_params struct {
	// 获得虚拟币
	ServerId   string `json:"serverId"`   // 服务器ID
	ServerName string `json:"serverName"` // 服务器名称
	RoleId     string `json:"roleId"`     // 角色ID
	RoleName   string `json:"roleName"`   // 角色名称
	RoleLevel  int    `json:"roleLevel"`  // 角色等级
	VipLevel   int    `json:"vipLevel"`   //! 角色VIP等级
	Number     int    `json:"number"`     // 虚拟币获得数量
	Source     string `json:"source"`     // 虚拟币来源
	CoinType   string `json:"coinType"`   // 虚拟币类型
	Coins      int    `json:"coins"`      // 获得后玩家虚拟币总数
}

//=================================================================
type SendRZ_UseMoney struct {
	Ts      string                 `json:"ts"`
	AppId   string                 `json:"appId"`
	Event   string                 `json:"event"`
	Params  SendRZ_UseMoney_params `json:"params"`
	EnvInfo SendRZ_envInfo         `json:"envInfo"`
}
type SendRZ_UseMoney_params struct {
	// 消耗虚拟币
	ServerId    string `json:"serverId"`    // 服务器ID
	ServerName  string `json:"serverName"`  // 服务器名称
	RoleId      string `json:"roleId"`      // 角色ID
	RoleName    string `json:"roleName"`    // 角色名称
	RoleLevel   int    `json:"roleLevel"`   // 角色等级
	VipLevel    int    `json:"vipLevel"`    //! 角色VIP等级
	Number      int    `json:"number"`      // 虚拟币消费数量
	Destination string `json:"destination"` // 消费去向
	CoinType    string `json:"coinType"`    // 货币类型
	Coins       int    `json:"coins"`       // 消费后玩家剩余货币数量
}

//====================================================================
type SendRZ_GetItem struct {
	Ts      string                `json:"ts"`
	AppId   string                `json:"appId"`
	Event   string                `json:"event"`
	Params  SendRZ_GetItem_params `json:"params"`
	EnvInfo SendRZ_envInfo        `json:"envInfo"`
}
type SendRZ_GetItem_params struct {
	// 获得物品
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"` //! 角色VIP等级

	ItemId    string `json:"itemId"`    //! 道具Id
	ItemType  string `json:"itemType"`  // 物品类型
	Item      string `json:"item"`      // 物品名称
	Number    int    `json:"number"`    // 物品数量
	Source    string `json:"source"`    // 获取途径
	ItemCount int    `json:"itemCount"` //! 当前物品总量
}

//======================================================================
type SendRZ_UseItem struct {
	Ts      string                `json:"ts"`
	AppId   string                `json:"appId"`
	Event   string                `json:"event"`
	Params  SendRZ_UseItem_params `json:"params"`
	EnvInfo SendRZ_envInfo        `json:"envInfo"`
}
type SendRZ_UseItem_params struct {
	// 消耗物品
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"` //! 角色VIP等级

	Destination string `json:"destination"` // 消耗去向
	ItemId      string `json:"itemId"`      //! 道具ID
	ItemType    string `json:"itemType"`    // 物品类型
	Item        string `json:"item"`        // 物品名称
	Number      int    `json:"number"`      // 消耗数量
	ItemCount   int    `json:"itemCount"`   //! 当前物品总量
}

//======================================================================
type SendRZ_BuyItem struct {
	Ts      string                `json:"ts"`
	AppId   string                `json:"appId"`
	Event   string                `json:"event"`
	Params  SendRZ_BuyItem_params `json:"params"`
	EnvInfo SendRZ_envInfo        `json:"envInfo"`
}
type SendRZ_BuyItem_params struct {
	// 购买道具
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"` //! 角色VIP等级

	Place  string `json:"place"`  // 购买位置
	Number int    `json:"number"` // 消耗虚拟币数量
	//Destination string `json:"destination"` // 虚拟币消费去向
	CoinType string `json:"coinType"` // 虚拟币类型
	ItemId   string `json:"itemId"`   // 道具Id
	ItemType string `json:"itemType"` // 物品类型
	Item     string `json:"item"`     // 物品名称
	Amount   int    `json:"amount"`   // 购买数量
}

//======================================================================
type SendSDK_Offline struct {
	OpenId              int    `json:"open_id"`
	AccountType         int    `json:"account_type"`
	RoleId              string `json:"role_id"`
	RoleName            string `json:"role_name"`
	RoleCreateTime      int    `json:"role_create_time"`
	RoleLevel           int    `json:"role_level"`
	PartyId             int    `json:"party_id"`
	VipLevel            int    `json:"vip_level"`
	VipRemain           int    `json:"vip_remain"`
	SVipRemain          int    `json:"svip_remain"`
	HistoryRecharge     int    `json:"history_recharge"`
	RechargeType        string `json:"recharge_type"`
	EventId             int    `json:"event_id"`
	AppId               int    `json:"app_id"`
	ServerId            int    `json:"server_id"`
	EventTime           int    `json:"event_time"`
	PlatformId          int    `json:"platform_id"`
	ChannelId           int    `json:"channel_id"`
	VipType             int    `json:"vip_type"`
	CurrencyBalanceInfo string `json:"currency_balance_info"`
	Distrubutor         int    `json:"distrubutor"`
	StoreName           string `json:"store_name"`

	RoleExp        int64  `json:"role_exp"`
	OnlineTime     int    `json:"online_time"`
	LastOperation  int    `json:"last_operation"`
	ReconnectCount int    `json:"reconnect_count"`
	SysOpCount     string `json:"sys_op_count"`
	PeriodItems    string `json:"period_items"`
}

type SendSDK_MoneyChange struct {
	OpenId              int    `json:"open_id"`
	AccountType         int    `json:"account_type"`
	RoleId              string `json:"role_id"`
	RoleName            string `json:"role_name"`
	RoleCreateTime      int    `json:"role_create_time"`
	RoleLevel           int    `json:"role_level"`
	PartyId             int    `json:"party_id"`
	VipLevel            int    `json:"vip_level"`
	VipRemain           int    `json:"vip_remain"`
	SVipRemain          int    `json:"svip_remain"`
	HistoryRecharge     int    `json:"history_recharge"`
	RechargeType        string `json:"recharge_type"`
	EventId             int    `json:"event_id"`
	AppId               int    `json:"app_id"`
	ServerId            int    `json:"server_id"`
	EventTime           int    `json:"event_time"`
	PlatformId          int    `json:"platform_id"`
	ChannelId           int    `json:"channel_id"`
	VipType             int    `json:"vip_type"`
	CurrencyBalanceInfo string `json:"currency_balance_info"`
	Distrubutor         int    `json:"distrubutor"`
	StoreName           string `json:"store_name"`

	CurrencyType    string `json:"currency_type"`
	ChangeType      string `json:"change_type"`
	Reason          string `json:"reason"`
	CurrencyCount   int    `json:"currency_count"`
	CurrencyBalance int    `json:"currency_balance"`
	ItemsId         int    `json:"items_id"`
	ItemsName       string `json:"items_name"`
	ItemsNum        int    `json:"items_num"`
}

type SendSDK_ItemChange struct {
	OpenId              int    `json:"open_id"`
	AccountType         int    `json:"account_type"`
	RoleId              string `json:"role_id"`
	RoleName            string `json:"role_name"`
	RoleCreateTime      int    `json:"role_create_time"`
	RoleLevel           int    `json:"role_level"`
	PartyId             int    `json:"party_id"`
	VipLevel            int    `json:"vip_level"`
	VipRemain           int    `json:"vip_remain"`
	SVipRemain          int    `json:"svip_remain"`
	HistoryRecharge     int    `json:"history_recharge"`
	RechargeType        string `json:"recharge_type"`
	EventId             int    `json:"event_id"`
	AppId               int    `json:"app_id"`
	ServerId            int    `json:"server_id"`
	EventTime           int    `json:"event_time"`
	PlatformId          int    `json:"platform_id"`
	ChannelId           int    `json:"channel_id"`
	VipType             int    `json:"vip_type"`
	CurrencyBalanceInfo string `json:"currency_balance_info"`
	Distrubutor         int    `json:"distrubutor"`
	StoreName           string `json:"store_name"`

	ChangeType   int    `json:"change_type"`
	Reason       string `json:"reason"`
	ItemsId      int    `json:"items_id"`
	ItemsType    int    `json:"items_type"`
	ItemsNum     int    `json:"items_num"`
	ItemsBalance int    `json:"items_balance"`
}

type SendSDK_Online struct {
	AppId           int   `json:"app_id"`
	ServerId        int   `json:"server_id"`
	OnlineRoleCount int   `json:"online_role_count"`
	EventId         int   `json:"event_id"`
	EventTime       int64 `json:"event_time"`
	DistributorId   int   `json:"distributor_id"`
}

//======================================================================
type SendRZ_PVE struct {
	Ts      string            `json:"ts"`
	AppId   string            `json:"appId"`
	Event   string            `json:"event"`
	Params  SendRZ_PVE_params `json:"params"`
	EnvInfo SendRZ_envInfo    `json:"envInfo"`
}
type SendRZ_PVE_params struct {
	// PVE战斗
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	VipLevel   int    `json:"vipLevel"` //! 角色VIP等级

	RoleLevel     int    `json:"roleLevel"`
	Begin_time    int64  `json:"beginTime"`    // 进入玩法时间
	Force         int64  `json:"force"`        // 角色战力
	TeamForce     int    `json:"teamForce"`    // 队伍战力
	MissionType   string `json:"missionType"`  // 关卡类型
	Mission       string `json:"mission"`      // 关卡名称
	Passed        string `json:"passed"`       // 完成情况(胜利失败)
	PassedType    string `json:"passedType"`   // 战斗类型(战斗/劝降/免战/扫荡)
	Duration      int    `json:"duration"`     // 完成耗时
	End_roleLevel int    `json:"endRoleLevel"` // 角色等级
}

//=====================================================================
type SendRZ_Chat struct {
	Ts      string             `json:"ts"`
	AppId   string             `json:"appId"`
	Event   string             `json:"event"`
	Params  SendRZ_Chat_params `json:"params"`
	EnvInfo SendRZ_envInfo     `json:"envInfo"`
}
type SendRZ_Chat_params struct {
	// 聊天
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	Type       string `json:"type"` // 所有聊天频道
	Text       string `json:"text"` // 文本聊天内容
}

//====================================================================
type SendRZ_PVP struct {
	Ts      string            `json:"ts"`
	AppId   string            `json:"appId"`
	Event   string            `json:"event"`
	Params  SendRZ_PVP_params `json:"params"`
	EnvInfo SendRZ_envInfo    `json:"envInfo"`
}
type SendRZ_PVP_params struct {
	// PVP战斗（不要）
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	VipLevel   int    `json:"vipLevel"`
	RoleLevel  int    `json:"roleLevel"`

	Begin_time    int64  `json:"beginTime"`    // 进入玩法时间
	Force         int64  `json:"force"`        // 角色战力
	TeamForce     int64  `json:"teamForce"`    //! 队伍战力
	MissionType   string `json:"missionType"`  // 关卡类型
	Mission       string `json:"mission"`      // 关卡名称
	Passed        string `json:"passed"`       // 完成情况(胜利失败)
	Duration      int    `json:"duration"`     // 完成耗时
	End_roleLevel int    `json:"endRoleLevel"` // 角色等级

	Pvpid   string `json:"pvpId"`   // 唯一PVP——ID
	Pvpname string `json:"pvpName"` // PVP名称
	Teamid  string `json:"teamId"`  // 团队ID
}

//=======================================================================
type SendRZ_onTime struct {
	Ts      string               `json:"ts"`
	AppId   string               `json:"appId"`
	Event   string               `json:"event"`
	Params  SendRZ_onTime_params `json:"params"`
	EnvInfo SendRZ_envInfo       `json:"envInfo"`
}
type SendRZ_onTime_params struct {
	// 5分钟在线统计
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	Timevalue  string `json:"timeValue"` // 五分钟一次记录时间：201705010915
	UserCnt    int    `json:"userCnt"`   // 在线人数
}

//===========================================================================
type SendRZ_HeroGain struct {
	Ts      string                 `json:"ts"`
	AppId   string                 `json:"appId"`
	Event   string                 `json:"event"`
	Params  SendRZ_HeroGain_params `json:"params"`
	EnvInfo SendRZ_envInfo         `json:"envInfo"`
}
type SendRZ_HeroGain_params struct {
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"`

	HeroType  string `json:"heroType"`
	HeroId    string `json:"heroId"`
	HeroName  string `json:"heroName"`
	HeroLevel int    `json:"heroLevel"`
	HeroStar  int    `json:"heroStar"`
	HeroRank  int    `json:"heroRank"`
	HeroForce int64  `json:"heroForce"`

	Source     string `json:"source"`
	ChipsCount int    `json:"chipsCount"`
}

//===========================================================================
type SendRZ_HeroLevelUp struct {
	Ts      string                    `json:"ts"`
	AppId   string                    `json:"appId"`
	Event   string                    `json:"event"`
	Params  SendRZ_HeroLevelUp_params `json:"params"`
	EnvInfo SendRZ_envInfo            `json:"envInfo"`
}
type SendRZ_HeroLevelUp_params struct {
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"`

	Force int64 `json:"force"`

	HeroType  string `json:"heroType"`
	HeroId    string `json:"heroId"`
	HeroName  string `json:"heroName"`
	HeroLevel int    `json:"heroLevel"`
	HeroStar  int    `json:"heroStar"`
	HeroRank  int    `json:"heroRank"`
	HeroForce int64  `json:"heroForce"`

	Source     string `json:"source"`
	ChipsCount int    `json:"chipsCount"`
}

//===========================================================================
type SendRZ_HeroStarup struct {
	Ts      string                   `json:"ts"`
	AppId   string                   `json:"appId"`
	Event   string                   `json:"event"`
	Params  SendRZ_HeroStarup_params `json:"params"`
	EnvInfo SendRZ_envInfo           `json:"envInfo"`
}
type SendRZ_HeroStarup_params struct {
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"`
	Force      int64  `json:"force"`

	HeroType  string `json:"heroType"`
	HeroId    string `json:"heroId"`
	HeroName  string `json:"heroName"`
	HeroLevel int    `json:"heroLevel"`
	HeroStar  int    `json:"heroStar"`
	HeroRank  int    `json:"heroRank"`
	HeroForce int64  `json:"heroForce"`

	ChipsCount int `json:"chipsCount"`
}

//===========================================================================
type SendRZ_HeroJoin struct {
	Ts      string                 `json:"ts"`
	AppId   string                 `json:"appId"`
	Event   string                 `json:"event"`
	Params  SendRZ_HeroJoin_params `json:"params"`
	EnvInfo SendRZ_envInfo         `json:"envInfo"`
}
type SendRZ_HeroJoin_params struct {
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"`
	Force      int64  `json:"force"`

	HeroType  string `json:"heroType"`
	HeroId    string `json:"heroId"`
	HeroName  string `json:"heroName"`
	HeroLevel int    `json:"heroLevel"`
	HeroStar  int    `json:"heroStar"`
	HeroRank  int    `json:"heroRank"`
	HeroForce int64  `json:"heroForce"`

	ChipsCount int    `json:"chipsCount"`
	CampName   string `json:"campName"`
}

//===========================================================================
type SendRZ_ArmyLevelUp struct {
	Ts      string                    `json:"ts"`
	AppId   string                    `json:"appId"`
	Event   string                    `json:"event"`
	Params  SendRZ_ArmyLevelUp_params `json:"params"`
	EnvInfo SendRZ_envInfo            `json:"envInfo"`
}
type SendRZ_ArmyLevelUp_params struct {
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"`
	Force      int64  `json:"force"`

	HeroType  string `json:"heroType"`
	HeroId    string `json:"heroId"`
	HeroName  string `json:"heroName"`
	HeroLevel int    `json:"heroLevel"`
	HeroStar  int    `json:"heroStar"`
	HeroRank  int    `json:"heroRank"`
	HeroForce int64  `json:"heroForce"`

	ArmyLevel int   `json:"armyLevel"`
	ArmyForce int64 `json:"armyForce"`

	Book      int `json:"book"`
	BookCount int `json:"bookCount"`
}

//===========================================================================
type SendRZ_FameRank struct {
	Ts      string                 `json:"ts"`
	AppId   string                 `json:"appId"`
	Event   string                 `json:"event"`
	Params  SendRZ_FameRank_params `json:"params"`
	EnvInfo SendRZ_envInfo         `json:"envInfo"`
}

type SendRZ_FameRank_params struct {
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"`
	Force      int64  `json:"force"`

	Rank int `json:"rank"`
}

//===========================================================================
type SendRZ_HeroDisplace struct {
	Ts      string                     `json:"ts"`
	AppId   string                     `json:"appId"`
	Event   string                     `json:"event"`
	Params  SendRZ_HeroDisplace_params `json:"params"`
	EnvInfo SendRZ_envInfo             `json:"envInfo"`
}

type SendRZ_HeroDisplace_params struct {
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"`
	Force      int64  `json:"force"`

	HeroIdBase        int `json:"heroIdBase"`
	HeroLevelBase     int `json:"heroLevelBase"`
	ArmyLevelBase     int `json:"armyLevelBase"`
	WeaponLevelBase   int `json:"weaponLevelBase"`
	ClothLevelBase    int `json:"clothLevelBase"`
	OrnamentLevelBase int `json:"ornamentLevelBase"`

	HeroIdTarget        int `json:"heroIdTarget"`
	HeroLevelTarget     int `json:"heroLevelTarget"`
	ArmyLevelTarget     int `json:"armyLevelTarget"`
	WeaponLevelTarget   int `json:"weaponLevelTarget"`
	ClothLevelTarget    int `json:"clothLevelTarget"`
	OrnamentLevelTarget int `json:"ornamentLevelTarget"`
}

//===========================================================================
type SendRZ_TaskFinish struct {
	Ts      string                   `json:"ts"`
	AppId   string                   `json:"appId"`
	Event   string                   `json:"event"`
	Params  SendRZ_TaskFinish_params `json:"params"`
	EnvInfo SendRZ_envInfo           `json:"envInfo"`
}

type SendRZ_TaskFinish_params struct {
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"`
	Force      int64  `json:"force"`

	Task_Type int    `json:"task_type"`
	Task_Id   int    `json:"task_id"`
	Task_Name string `json:"task_name"`
	Task_Star int    `json:"task_star"`
}

//===========================================================================
// 3.1.31 	客户端激活-参数
type SendRZ_Activation struct {
	Ts      string                   `json:"ts"`
	AppId   string                   `json:"appId"`
	Event   string                   `json:"event"`
	Params  SendRZ_Activation_params `json:"params"`
	EnvInfo SendRZ_envInfo           `json:"envInfo"`
}

// 3.1.31 	客户端激活-参数
type SendRZ_ActivationIOS struct {
	Ts      string                   `json:"ts"`
	AppId   string                   `json:"appId"`
	Event   string                   `json:"event"`
	Params  SendRZ_Activation_params `json:"params"`
	EnvInfo SendRZ_envInfo_ios       `json:"envInfo"`
}

type SendRZ_Activation_params struct {
	ServerId   string `json:"serverId"`   // 服务器ID
	ServerName string `json:"serverName"` // 服务器名称
}

//===========================================================
type SendRZ_LevelupIOS struct {
	Ts      string                   `json:"ts"`
	AppId   string                   `json:"appId"`
	Event   string                   `json:"event"`
	Params  SendRZ_LevelupIOS_params `json:"params"`
	EnvInfo SendRZ_envInfo_ios       `json:"envInfo"`
}

type SendRZ_LevelupIOS_params struct {
	// 角色升级成功
	ServerId   string `json:"serverId"`   // 服务器ID
	ServerName string `json:"serverName"` // 服务器名称
	RoleId     string `json:"roleId"`     // 角色ID
	RoleName   string `json:"roleName"`   // 角色名称
	RoleLevel  int    `json:"roleLevel"`  // 角色等级
	VipLevel   int    `json:"vipLevel"`   //! 角色VIP等级
}

//===========================================================
type SendRZ_AccountChargeSuccess struct {
	Ts      string                             `json:"ts"`
	AppId   string                             `json:"appId"`
	Event   string                             `json:"event"`
	Params  SendRZ_AccountChargeSuccess_params `json:"params"`
	EnvInfo SendRZ_envInfo_ios                 `json:"envInfo"`
}

type SendRZ_AccountChargeSuccess_params struct {
	// 角色升级成功
	ServerId   string  `json:"serverId"`   // 服务器ID
	ServerName string  `json:"serverName"` // 服务器名称
	RoleId     string  `json:"roleId"`     // 角色ID
	RoleName   string  `json:"roleName"`   // 角色名称
	RoleLevel  int     `json:"roleLevel"`  // 角色等级
	PayMent    string  `json:"payment"`    //! 支付方式
	DsorderId  string  `json:"dsorderId"`  //! 订单号
	Currency   string  `json:"currency"`   //! 币种
	Amount     float64 `json:"amount"`     //! 订单金额
	GoodsInfo  string  `json:"goodsInfo"`  //! 商品信息
}

//======================================================

type SendRZ_LoginokIOSEX struct {
	Ts      string                `json:"ts"`
	AppId   string                `json:"appId"`
	Event   string                `json:"event"`
	Params  SendRZ_Loginok_params `json:"params"`
	EnvInfo SendRZ_envInfo_ios    `json:"envInfo"`
}

type SendRZ_LoginokIOS struct {
	Ts      string                   `json:"ts"`
	AppId   string                   `json:"appId"`
	Event   string                   `json:"event"`
	Params  SendRZ_LoginokIOS_params `json:"params"`
	EnvInfo SendRZ_envInfo_ios       `json:"envInfo"`
}

// 角色登录成功-参数
type SendRZ_LoginokIOS_params struct {
	ServerId   string `json:"serverId"`   // 服务器ID
	ServerName string `json:"serverName"` // 服务器名称
}

//===========================================================================
type SendRZ_CommerceJoin struct {
	Ts      string                     `json:"ts"`
	AppId   string                     `json:"appId"`
	Event   string                     `json:"event"`
	Params  SendRZ_CommerceJoin_params `json:"params"`
	EnvInfo SendRZ_envInfo             `json:"envInfo"`
}

type SendRZ_CommerceJoin_params struct {
	ServerId   string `json:"serverId"`
	ServerName string `json:"serverName"`
	RoleId     string `json:"roleId"`
	RoleName   string `json:"roleName"`
	RoleLevel  int    `json:"roleLevel"`
	VipLevel   int    `json:"vipLevel"`
	Force      int64  `json:"force"`

	Type_Id   int    `json:"type_id"`
	Type_Name string `json:"type_name"`
}

//! 钻石消耗
type S2C_GemCostTop struct {
	Cid string       `json:"cid"`
	Top []*Js_ActTop `json:"top"`
	Ver int          `json:"ver"`
	Cur int          `json:"cur"`
	Old int          `json:"old"`
	Num int          `json:"num"`
}

//! 活动领取
type S2C_GetCityNum struct {
	Cid  string     `json:"cid"`
	Id   int        `json:"id"`
	Item []PassItem `json:"item"`
}

type S2C_GetActTop struct {
	Cid      string       `json:"cid"`
	TopInfo  []*Js_ActTop `json:"top"` //! key:taskType, value: topInfo
	TaskType int          `json:"tasktype"`
}

//GVG 压制
type S2C_CityGvGAttack struct {
	Cid string `json:"cid"`
	Tag int    `json:"tag"`
	Msg string `json:"msg"`
}

//! GVG 请求自己的城池信息
type S2C_CityGVGSelf struct {
	Cid    string `json:"cid"`
	CityId int    `json:"cityid"` //! 城市ID
}

//GVG 驻扎
type S2C_CityGvGDefense struct {
	Cid string `json:"cid"`
	Tag int    `json:"tag"`
	Msg string `json:"msg"`
}

//! 购买行动力
type S2C_CityGvGBuyPower struct {
	Cid      string     `json:"cid"`
	Power    int        `json:"power"`
	BuyPower int        `json:"buypower"`
	Items    []PassItem `json:items`
}

//GVG 开始战斗
type S2C_CityGvGFightBegin struct {
	Cid     string         `json:"cid"`
	CityId  int            `json:"cityid"`  //! 城池信息
	Targets []JS_FightInfo `json:"targets"` //!
}

//GVG 战斗结果通知
type S2C_CityGvGFightResult struct {
	Cid    string       `json:"cid"`
	Tag    int          `json:"tag"`
	Target JS_FightInfo `json:"target"` //!
}

type S2C_GetRedPac struct {
	Cid     string     `json:"cid"`
	RedWait []*RedWait `json:"redwait"` //! 待发送的红包
}

//! 红包池红包发送
type S2C_SendPool struct {
	Cid     string     `json:"cid"`     //! 命令
	RedWait []*RedWait `json:"redwait"` //! 同步红包
}

//! 广播红包信息
type S2C_RedPac struct {
	Cid   string `json:"cid"`   //! 命令
	KeyId int64  `json:"keyid"` //! 红包Id
}

//! 抢红包
type S2C_GotRedPac struct {
	Cid        string      `json:"cid"`        //! 命令
	UserRedPac *UserRedPac `json:"userredpac"` //! 需要更新的红包信息
	ItemId     int         `json:"itemid"`     //! 道具Id
}

type RedInfo struct {
	Name string `json:"name"` //! 玩家姓名
	Num  int    `json:"num"`  //! 领取的数量
}

//! 查看红包
type S2C_LookRedPac struct {
	Cid         string     `json:"cid"`         //! 命令
	IconId      int        `json:"iconid"`      //! 头像id
	Name        string     `json:"name"`        //! 名字
	Msg         string     `json:"msg"`         //! 留言
	Num         int        `json:"num"`         //! 获得道具
	ItemId      int        `json:"itemid"`      //! 道具Id
	TakenPacNum int        `json:"takenpacnum"` //! 已经领取红包数量
	AllPacNum   int        `json:"allpacnum"`   //! 总红包数量
	TakenNum    int        `json:"takennum"`    //! 已经领取道具数量
	AllNum      int        `json:"allnum"`      //! 总的道具数量
	RedInfo     []*RedInfo `json:"redinfo"`     //! 抢红包玩家信息
	Uid         int64      `json:"uid"`         //! 发送人uid
	Thank       int        `json:"thank"`       //! 是否答谢 0未答谢 1答谢
}

//! 答谢红包
type S2C_ThankRedPac struct {
	Cid        string      `json:"cid"`        //! 命令
	UserRedPac *UserRedPac `json:"userredpac"` //! 需要更新的红包状态
}

//! 发送军团红包
type S2C_SendUnionRed struct {
	Cid        string `json:"cid"`        //! 命令
	GemNum     int    `json:"gemnum"`     //! 玩家钻石信息
	LeftNum    int    `json:"leftgem"`    //! 剩余数量
	LeftItemId int    `json:"leftitemid"` //! 对应的道具信息
}

//! 发送/抢得的红包历史
type S2C_GetRedHis struct {
	Cid   string       `json:"cid"`   //! 命令
	Today []*RedRecord `json:"today"` //! 今天红包历史
	Total []*RedRecord `json:"total"` //! 历史红包历史
}

//! 获取红包信息
type S2C_GetRedPacs struct {
	Cid        string        `json:"cid"`        //! 命令
	RedPacShow []*RedPacShow `json:"redpacshow"` //! 玩家红包信息
	GemNum     int           `json:"gemnum"`     //! 剩余钻石额度
	SendLimit  []*RedRefresh `json:"limitnum"`   //! 发送上限
	GotLimit   []*RedRefresh `json:"gotnum"`     //! 接收上限
}

//! 活动物品兑换
type S2C_ActDailyDiscount struct {
	Cid      string     `json:"cid"`
	Id       int        `json:"id"`
	Progress int        `json:"progress"`
	Done     int        `json:"done"`
	Item     []PassItem `json:"item"`
}

type S2C_GetOverflow struct {
	Cid      string        `json:"cid"`
	GetItems []PassItem    `json:"getitems"`
	Month    *JS_MonthCard `json:"month"`
}

//! 领取月卡累计领取奖励
type S2C_GetCumulativeReward struct {
	Cid  string     `json:"cid"`
	Id   int        `json:"id"`
	Item []PassItem `json:"item"`
}

type S2C_GetCumulativeRewardNew struct {
	Cid      string         `json:"cid"`
	Id       int            `json:"id"`
	Item     []PassItem     `json:"item"`
	CostItem []PassItem     `json:"costitem"`
	Month    []JS_MonthCard `json:"month"`
}

//! 转盘掉落信息
type S2C_DialLootInfo struct {
	Cid        string      `json:"cid"`       //! cid
	Ids        []int       `json:"ids"`       //! 抽中的格子Id
	ReturnItem []PassItem  `json:"items"`     //! 掉落物品,以及扣除的钻石或者钥匙
	FreeTimes  int         `json:"freetimes"` //! 剩余免费次数
	Times      int         `json:"times"`     //! 抽奖次数
	Msg        []*DialInfo `json:"msg"`       //! 新增提示消息
}

//! 转盘排行榜
type S2C_DialRank struct {
	Cid  string         `json:"cid"`  //! cid
	Rank []*Js_DialRank `json:"rank"` //! 排行榜
}

//! 翻牌币购买
type S2C_DrawBuy struct {
	Cid   string     `json:"cid"`   //! cid
	Gem   int        `json:"gem"`   //! 当前钻石
	Items []PassItem `json:"items"` //! 获得物品
}

//! 开工福利锤子个数更新
type S2C_LuckEggUpdate struct {
	Cid          string     `json:"cid"`
	RechargeNum  int        `json:"rechargenum"`  //! 充值获得锤子个数
	RechargeLeft int        `json:"rechargeleft"` //! 充值剩余
	Items        []PassItem `json:"items"`        //! 同步锤子道具
}

//! 连续充值任务状态
type S2C_DailyRecharge struct {
	Cid       string                 `json:"cid"`
	Step      int                    `json:"step"`         // 是哪一期活动
	Index     int                    `json:"index"`        //! 第几天
	Info      []*RechargeTask        `json:"rechargetask"` // 充值状态
	TodayNum  int                    `json:"todaynum"`     // 今日充值
	TodayLeft int                    `json:"todayleft"`    // 剩余充值
	Data      []*DailyrechargeConfig `json:"data"`         //! 对应的配置数据
}

//! 连续充值
type S2C_DailyRechargeAward struct {
	Cid   string        `json:"cid"`
	Task  *RechargeTask `json:"task"`
	Items []PassItem    `json:"items"`
}

//! 抢红包状态
type S2C_RedPacCode struct {
	Cid  string `json:"cid"`  //! 命令
	Code int    `json:"code"` //! 错误码 1 表示今日可抢%s达到上限
}

//! 新增红包数量同步
type S2C_SynRedGemNum struct {
	Cid    string `json:"cid"`
	GemNum int    `json:"gemnum"` //! 同步gemnum
}

//! 布阵信息
type S2C_UpdateTeamPos struct {
	Cid  string      `json:"cid"`
	Team *Js_TeamPos `json:"teampos"` //! 布阵信息
}

//! 布阵信息
type S2C_UpdateArenaSpecialTeamPos struct {
	Cid  string        `json:"cid"`
	Team []*Js_TeamPos `json:"teampos"` //! 布阵信息
}

//! 布阵信息
type S2C_SaveTeam struct {
	Cid  string      `json:"cid"`
	Team *Js_TeamPos `json:"teampos"` //! 布阵信息
}

//! 所有的布阵信息
type S2C_TeamPos struct {
	Cid        string        `json:"cid"`
	TeamPos    []*Js_TeamPos `json:"teampos"`    //! 所有布阵信息
	PreTeamPos []*Js_TeamPos `json:"preteampos"` //! 所有布阵信息
}

//! 激活资质
type S2C_ActivateStar struct {
	Cid      string             `json:"cid"`
	HeroId   int                `json:"heroid"`   //! 英雄Id
	StarItem *StarItem          `json:"staritem"` //! 升星信息
	Items    []PassItem         `json:"items"`    //! 需要的道具
	Attr     map[int]*Attribute `json:"attr"`     //! 属性
}

//! 升星
type S2C_HeroUpStar struct {
	Cid           string           `json:"cid"`
	HeroKeyId     int              `json:"herokeyid"`     //! 英雄Id
	StarItem      *StarItem        `json:"staritem"`      //! 升星信息
	CostHeroKeyId []int            `json:"costherokeyid"` //!
	OffEquip      []*Equip         `json:"offequip"`      //! 有变化的装备同步
	OffArtifact   []*ArtifactEquip `json:"offartifact"`   //! 有变化的神器同步
	OffHorseInfo  []*JS_HorseInfo  `json:"offhorseinfo"`  //! 有变化的魔宠同步
	GetItems      []PassItem       `json:"getitems"`      //英雄材料可能存在的返还
	Item          *PassItem        `json:"item"`          //能获得的成长涓流
	StageTalent   *StageTalent     `json:"stagetalent"`   // 天赋
}

//! 升星
type S2C_HeroUpStarAll struct {
	Cid           string           `json:"cid"`
	HeroKeyId     []int            `json:"herokeyid"`     //! 英雄Id
	StarItem      []*StarItem      `json:"staritem"`      //! 升星信息
	CostHeroKeyId [][]int          `json:"costherokeyid"` //!
	OffEquip      []*Equip         `json:"offequip"`      //! 有变化的装备同步
	OffArtifact   []*ArtifactEquip `json:"offartifact"`   //! 有变化的神器同步
	OffHorseInfo  []*JS_HorseInfo  `json:"offhorseinfo"`  //! 有变化的魔宠同步
	GetItems      []PassItem       `json:"getitems"`      //英雄材料可能存在的返还
	StageTalent   []*StageTalent   `json:"stagetalent"`   // 天赋
}

type S2C_HeroLock struct {
	Cid       string `json:"cid"`
	HeroKeyId int    `json:"herokeyid"` //! 英雄keyId
	IsLock    int    `json:"islock"`    //! 上锁状态
}

type S2C_BuyPos struct {
	Cid       string     `json:"cid"`
	BuyPosNum int        `json:"buyposnum"` //! 购买次数
	CostItems []PassItem `json:"costitems"` //! 花费
}

//! 一键升星
type S2C_UpStarAuto struct {
	Cid      string             `json:"cid"`
	HeroId   int                `json:"heroid"`   //! 英雄Id
	StarItem *StarItem          `json:"staritem"` //! 升星信息
	Items    []PassItem         `json:"items"`    //! 道具
	Attr     map[int]*Attribute `json:"attr"`     //! 属性
}

//! 激活天赋
type S2C_UpgradeTalent struct {
	Cid        string             `json:"cid"`
	HeroId     int                `json:"heroid"`     //! 英雄Id
	TalentItem *TalentInfo        `json:"talentitem"` //! 天赋信息
	Items      []PassItem         `json:"items"`      //! 道具
	Step       int                `json:"step"`       //! 觉醒层级
	Attr       map[int]*Attribute `json:"attr"`       //! 属性
	MainTelent int                `json:"maintalent"` //! 主神格等级
}

//! 重置天赋
type S2C_ResetTalent struct {
	Cid      string             `json:"cid"`
	HeroId   int                `json:"heroid"`   //! 英雄Id
	Items    []PassItem         `json:"items"`    //! 道具
	Step     int                `json:"step"`     //! 觉醒层级
	Attr     map[int]*Attribute `json:"attr"`     //! 属性
	OutItems []PassItem         `json:"outitems"` //! 道具
}

//! 获得幻境基础信息
type S2C_DreamLandBaseInfo struct {
	Cid        string `json:"cid"`
	FreeTimes1 int    `json:"freetimes1` // 免费次数1
	FreeTimes2 int    `json:"freetimes2` // 免费次数2

	LuckyTimes1 int `json:"luckytimes1` // 幸运值 以总抽奖次数% 固定数值
	TypeTimes1  int `json:"typetimes1`  // 幸运值上限 固定数值

	LuckyTimes2 int `json:"luckytimes2` // 幸运值 以总抽奖次数% 固定数值
	TypeTimes2  int `json:"typetimes2`  // 幸运值上限 固定数值

	RefreshCount1 int `json:"refreshcount1"` // 当天刷新次数
	RefreshCount2 int `json:"refreshcount2"` // 当天刷新次数
}

//! 获得幻境物品信息
type S2C_DreamLandItemInfo struct {
	Cid          string `json:"cid"`
	Type         int    `json:"type"`         //! 类型
	RefreshCount int    `json:"refreshcount"` // 当天刷新次数
	FreeTimes    int    `json:"freetimes`     // 免费次数
	LuckyTimes   int    `json:"luckytimes`    // 幸运值 以总抽奖次数% 固定数值
	TypeTimes    int    `json:"typetimes`     // 幸运值上限 固定数值
	Items        []PassItem
}

//! 幻境抽奖
type S2C_DreamLandLoot struct {
	Cid        string `json:"cid"`
	Type       int    `json:"type"`      //! 类型
	FreeTime   int    `json:"freetime"`  //免费次数
	LuckyTimes int    `json:"luckytimes` // 幸运值 以总抽奖次数% 固定数值
	TypeTimes  int    `json:"typetimes`  // 幸运值上限 固定数值
	Indexs     []int  `json:"indexs`     // 抽取的位置
	OutItems   []PassItem
	InItems    []PassItem
}

//! 幻境刷新
type S2C_DreamLandRefresh struct {
	Cid          string `json:"cid"`
	Type         int    `json:"type"`         //! 类型
	RefreshCount int    `json:"refreshcount"` // 当天刷新次数
	Items        []PassItem
	NewItem      []PassItem
}

//! 激活缘分
type S2C_ActivateFate struct {
	Cid      string             `json:"cid"`
	HeroId   int                `json:"heroid"`   //! 英雄Id
	FateInfo *FateInfo          `json:"fateitem"` //! 缘分信息
	Attr     map[int]*Attribute `json:"attr"`     //! 属性
	Chanegs  map[int]*Attribute `json:"changes"`  //! 属性
}

//重生
type S2C_HeroReborn struct {
	Cid           string           `json:"cid"`
	HeroKeyId     int              `json:"herokeyid"`     //! 英雄Id
	StarItem      *StarItem        `json:"staritem"`      //! 升星信息
	CostItem      []PassItem       `json:"costitem"`      //! 消耗物品
	GetItem       []PassItem       `json:"getitem"`       //!  获得材料
	OffEquip      []*Equip         `json:"offequip"`      //! 有变化的装备同步
	OffArtifact   []*ArtifactEquip `json:"offartifact"`   //! 有变化的神器同步
	OffHorseInfo  []*JS_HorseInfo  `json:"offhorseinfo"`  //! 有变化的魔宠同步
	HeroLv        int              `json:"herolv"`        //!
	OriginalLevel int              `json:"originallevel"` //!
}

//英雄回退
type S2C_HeroBack struct {
	Cid            string           `json:"cid"`
	HeroKeyId      int              `json:"herokeyid"`      //! 英雄Id
	StarItem       *StarItem        `json:"staritem"`       //! 升星信息
	CostItem       []PassItem       `json:"costitem"`       //! 消耗物品
	GetItem        []PassItem       `json:"getitem"`        //!  获得材料
	OffEquip       []*Equip         `json:"offequip"`       //! 有变化的装备同步
	OffArtifact    []*ArtifactEquip `json:"offartifact"`    //! 有变化的神器同步
	OffHorseInfo   []*JS_HorseInfo  `json:"offhorseinfo"`   //! 有变化的魔宠同步
	HeroLv         int              `json:"herolv"`         //!
	OriginalLevel  int              `json:"originallevel"`  //!
	ExclusiveEquip *ExclusiveEquip  `json:"exclusiveequip"` // 专属装备
	StageTalent    *StageTalent     `json:"stagetalent"`    // 天赋
}

//遣散
type S2C_HeroFire struct {
	Cid          string           `json:"cid"`
	HeroKeyId    []int            `json:"herokeyid"`    //! 英雄Id
	GetItem      []PassItem       `json:"getitem"`      //!  获得材料
	OffEquip     []*Equip         `json:"offequip"`     //! 有变化的装备同步
	OffArtifact  []*ArtifactEquip `json:"offartifact"`  //! 有变化的神器同步
	OffHorseInfo []*JS_HorseInfo  `json:"offhorseinfo"` //! 有变化的魔宠同步
}

//! 虚空英雄共鸣
type S2C_VoidHeroResonanceSet struct {
	Cid           string     `json:"cid"`
	HeroKeyId     int        `json:"herokeyid"`     //! 英雄keyId
	VoidHeroKeyId int        `json:"voidherokeyid"` //! 英雄keyId
	GetItem       []PassItem `json:"getitem"`       //! 获得物品
	VoidHero      *Hero      `json:"voidhero"`
}

//! 取消虚空英雄共鸣
type S2C_VoidHeroResonanceCancel struct {
	Cid           string     `json:"cid"`
	VoidHeroKeyId int        `json:"herokeyid"` //! 英雄keyId
	CostItem      []PassItem `json:"costitem"`  //! 消耗物品
	VoidHero      *Hero      `json:"voidhero"`
}

//! 更新虚空英雄共鸣
type S2C_VoidHeroResonanceUpdate struct {
	Cid      string `json:"cid"`
	VoidHero *Hero  `json:"voidhero"`
}

// 设置天赋技能
type S2C_SetStageTalentSkill struct {
	Cid       string            `json:"cid"`
	HeroKeyId int               `json:"herokeyid"` //! 英雄keyId
	Index     int               `json:"index"`     //! 第几层
	Pos       int               `json:"pos"`       //! 技能位置
	Info      *StageTalentIndex `json:"info"`      //! 技能结构
}

//! 装备操作信息
type S2C_EquipAction struct {
	Cid             string     `json:"cid"`             //! cid
	Action          int        `json:"action"`          //! 装备操作
	EquipItem       *Equip     `json:"equipitem"`       //! 装备信息
	ExchangeItem    *Equip     `json:"exchangeitem"`    //! 替换装备
	Items           []PassItem `json:"items"`           //! 道具信息[合成, 强化道具信息]
	RecastTempKeyId int        `json:"recasttempkeyid"` //! 重铸临时保存key
	RecastTempAim   int        `json:"recasttempaim"`   //! 重置临时保存aim
}

type S2C_EquipActionAll struct {
	Cid       string   `json:"cid"`       //! cid
	Action    int      `json:"action"`    //! 装备操作
	EquipItem []*Equip `json:"equipitem"` //! 产生变动的装备信息
}

//! 神器操作信息
type S2C_ArtifactEquipAction struct {
	Cid               string         `json:"cid"`               //! cid
	Action            int            `json:"action"`            //! 装备操作
	ArtifactEquipItem *ArtifactEquip `json:"artifactequipitem"` //! 装备信息
	ExchangeItem      *ArtifactEquip `json:"exchangeitem"`      //! 替换装备
	Items             []PassItem     `json:"items"`             //! 道具信息[合成, 强化道具信息]
}

type S2C_ArtifactEquipUpLv struct {
	Cid               string         `json:"cid"`               //! cid
	ArtifactEquipItem *ArtifactEquip `json:"artifactequipitem"` //! 装备信息
	CostItems         []PassItem     `json:"costitems"`         //!
}

type S2C_ExclusiveAction struct {
	Cid           string          `json:"cid"`           //! cid
	Action        int             `json:"action"`        //! 装备操作
	ExclusiveItem *ExclusiveEquip `json:"exclusiveitem"` //! 装备信息
	CostItems     []PassItem      `json:"costitems"`     //! 道具消耗
	HeroKeyId     int             `json:"herokeyid"`     //! 英雄KEYID
}

//! 装备信息
type S2C_EquipInfo struct {
	Cid     string                         `json:"cid"`     //! cid
	Equips  [EQUIP_PACK_NUM]map[int]*Equip `json:"equips"`  //! 宝物信息
	HeroAtt []*HeroAttr                    `json:"heroatt"` //! 英雄属性
	Ver     int                            `json:"ver"`     //! 开启分页功能
}

type S2C_ArtifactEquipInfo struct {
	Cid            string                 `json:"cid"`            //! cid
	ArtifactEquips map[int]*ArtifactEquip `json:"artifactequips"` //! 宝物信息
}

//! 同步装备信息
type S2C_SynEquip struct {
	Cid    string   `json:"cid"`    //! cid
	Equips []*Equip `json:"equips"` //! 装备信息
}

type S2C_SynActifactEquip struct {
	Cid            string           `json:"cid"`            //! cid
	ArtifactEquips []*ArtifactEquip `json:"artifactequips"` //! 装备信息
}

//! 宝石合成
type S2C_GemAction struct {
	Cid     string     `json:"cid"`       //! cid
	Items   []PassItem `json:"items"`     //! 道具信息[宝石合成, 宝石进背包]
	PItem   *Equip     `json:"equipitem"` //! 装备信息[脱装备, 一键脱装备, 穿装备]
	HeroAtt *HeroAttr  `json:"heroatt"`   //! 英雄属性
}

type S2C_HeroDecompose struct {
	Cid   string     `json:"cid"`
	Items []PassItem `json:"items"`
}

//! 同步战斗力
type S2C_SynFight struct {
	Cid       string `json:"cid"`
	HeroKeyId int    `json:"herokeyid"`
	Fight     int64  `json:"fight"`
	Reason    int    `json:"reason"`
	HeroLv    int    `json:"herolv"`
}

//! 全队同步战斗力
type S2C_SynAllFight struct {
	Cid       string  `json:"cid"`
	HeroId    []int   `json:"heroids"`
	Fight     []int64 `json:"fights"`
	Reason    int     `json:"reason"`
	BossId    int     `json:"bossid"`
	BossFight int64   `json:"bossfight"`
}

//!镇魂塔相关
//数据获取
type S2C_TowerInfo struct {
	Cid  string      `json:"cid"`
	Data []*JS_Tower `json:"data"`
}

//排行榜
type S2C_TowerRank struct {
	Cid   string       `json:"cid"`
	Items []*Js_ActTop `json:"items"`
}

type S2C_TowerFightBegin struct {
	Cid string `json:"cid"`
}

type S2C_TowerFightResult struct {
	Cid        string     `json:"cid"`
	MaxLevel   int        `json:"maxlevel"`   // 历史最大关卡
	CurLevel   int        `json:"curlevel"`   // 当前关卡
	LevelCount int        `json:"levelcount"` // 当前关卡
	Items      []PassItem `json:"items"`      // 道具信息
}

type S2C_TowerFightSkip struct {
	Cid        string     `json:"cid"`
	MaxLevel   int        `json:"maxlevel"`   // 历史最大关卡
	CurLevel   int        `json:"curlevel"`   // 当前关卡
	LevelCount int        `json:"levelcount"` // 当前关卡
	Items      []PassItem `json:"items"`      // 道具信息
}

//! 楼层信息
type S2C_TowerFloorInfo struct {
	Cid     string                 `json:"cid"`
	LevelId int                    `json:"levelid"` //! 关卡id
	Record  []*Js_TowerFightRecord `json:"record"`  //三条当前楼层的挑战记录
}

//! 宝石副本相关操作
//数据获取
type S2C_GemStoneInfo struct {
	Cid           string `json:"cid"`
	ChapterIdx    int    `json:"chapterid"`     //! 章节id
	LevelIdx      int    `json:"levelidx"`      //! 关卡id
	SweepTimes    []int  `json:"sweeptimes"`    //! 扫荡次数
	BuySweepTimes []int  `json:"buysweeptimes"` //! 购买扫荡次数
}

type S2C_GemStoneFightResult struct {
	Cid        string     `json:"cid"`
	ChapterIdx int        `json:"chapterid"` //! 章节id
	LevelIdx   int        `json:"levelidx"`  //! 关卡id
	LevelId    int        `json:"levelid"`   //! 关卡流水id
	Items      []PassItem `json:"items"`     //! 道具信息
}

type S2C_GemStoneSweep struct {
	Cid           string     `json:"cid"`
	Items         []PassItem `json:"items"`         //! 道具信息
	SweepTimes    []int      `json:"sweeptimes"`    //! 扫荡次数
	BuySweepTimes []int      `json:"buysweeptimes"` //! 购买扫荡次数
}

type S2C_GemStoneBuySweepTimes struct {
	Cid           string `json:"cid"`
	SweepTimes    []int  `json:"sweeptimes"`    //! 扫荡次数
	BuySweepTimes []int  `json:"buysweeptimes"` //! 购买扫荡次数
}

//! 一键领取所有任务
type S2C_FinishAllTask struct {
	Cid      string         `json:"cid"`
	TaskInfo []*JS_TaskInfo `json:"taskinfo"`
	OutItem  []PassItem     `json:"outitem"`
}

type S2C_TowerBuffRefresh struct {
	Cid          string     `json:"cid"`
	Buffs        []int      `json:"buffs"`
	CurBuff      int        `json:"curbuff"`
	BuffBuyTimes int        `json:"buffbuytimes"` //重置购买次数
	BuffTimes    int        `json:"bufftimes"`    // 免费次数
	Items        []PassItem `json:"items"`        // 道具
}

type S2C_SetBuff struct {
	Cid     string `json:"cid"`
	CurBuff int    `json:"curbuff"`
}

type S2C_TowerReset struct {
	Cid           string     `json:"cid"`
	CurLevel      int        `json:"curlevel"`
	CurFailBox    int        `json:"curfailbox"`
	BoxState      []int      `json:"boxstate"`
	CurBuff       int        `json:"curbuff"`
	Buff          []int      `json:"buff"`
	ResetTimes    int        `json:"resettimes"`
	ResetBuyTimes int        `json:"resetbuytimes"`
	Items         []PassItem `json:"items"`
	TowerBox      []int      `json:"tower_box"`
}

//! 镇魂宝箱奖励
type S2C_FailedBoxPrize struct {
	Cid      string     `json:"cid"`
	BoxState []int      `json:"boxstate"` // 已经领取的宝箱Id
	Items    []PassItem `json:"item"`
}

//! 爬塔宝箱奖励
type S2C_TowerBoxPrize struct {
	Cid      string     `json:"cid"`
	BoxState []int      `json:"boxstate"` // 已经领取的宝箱Id
	Items    []PassItem `json:"items"`    // 道具
}

//! 爬塔精英关卡次数
type S2C_BuyAdvanceTimes struct {
	Cid             string     `json:"cid"`
	Items           []PassItem `json:"items"`           // 已经领取的宝箱Id
	AdvanceBuyTimes int        `json:"advancebuytimes"` // 高级购买次数
}

//! 爬塔精英关卡次数
type S2C_ResetBuyTimes struct {
	Cid           string     `json:"cid"`
	Items         []PassItem `json:"items"`         // 已经领取的宝箱Id
	ResetBuyTimes int        `json:"resetbuytimes"` // 高级购买次数
}

//// 斗技场国家排行榜
//type S2C_PvpTop struct {
//	Cid string       `json:"cid"`
//	Top []*JS_TopJJC `json:"top"`
//}

type S2C_SynTalentNum struct {
	Cid      string `json:"cid"`
	TotalNum int    `json:"total_num"`
}

// 攻城提示
type S2C_AttackMention struct {
	Cid    string `json:"cid"`
	Camp   int    `json:"camp"`
	Icon   int    `json:"icon"`
	Name   string `json:"name"`
	MineId int    `json:"mine_id"`
}

// 当玩家胜利后，被攻击的据点上会显示"驻守人数-1"
type S2C_MineEvents struct {
	Cid    string `json:"cid"`
	MineId int    `json:"mine_id"`
	Action int    `json:"action"` // 1.action=1, "驻守人数-1" 2.发生战斗 3.战斗结束
}

type S2C_SyncMineResult struct {
	Cid          string        `json:"cid"`
	Reason       int           `json:"reason"`        // 1 击杀 2死亡
	ReturnMineId int           `json:"mineId"`        // 回程Id
	KillNum      int           `json:"kill_num"`      // 击杀次数
	DeadNum      int           `json:"dead_num"`      // 被击杀次数
	FightId      int64         `json:"fight_id"`      // 战报Id
	Camp         int           `json:"camp"`          // 敌方阵营
	Icon         int           `json:"icon"`          // 敌方头像
	MineId       int           `json:"mine_id"`       // 战斗完后所在城池
	Side         int           `json:"side"`          // 1 攻击 2 防守
	EnemyName    string        `json:"enemy_name"`    // 敌方name
	AttackTimes  int           `json:"attack_times"`  // 奔袭buff次数
	DefenceTimes int           `json:"defence_times"` // 驰援buff次数
	FightCd      int64         `json:"fight_cd"`      // 增加攻击cd
	Items        map[int]*Item `json:"items"`         // 道具信息
	CollectCd    int64         `json:"collect_cd"`    // 征收CD
	Hp           [2][5]int     `json:"hp"`            //! 剩余血量
	BossHP       [2]int        `json:"bosshp"`        //! 巨兽的HP
	HeroId       [2][5]int     `json:"heroid"`        //! 剩余血量对应的英雄ID
}

type PlayerRatio struct {
	Name  string  `json:"name"`
	Ratio float64 `json:"ratio"`
}

type GveDragonMsg struct {
	BuildId int     `json:"build_id"` // 建筑Id
	LevelId int     `json:"level_id"` // 关卡Id
	Ratio   float64 `json:"ratio"`    // 当前被占领值百分比, 最大可占领值读取配置
}

// 阵营变更
type S2C_CampChangeMention struct {
	Cid     string `json:"cid"`
	OldCamp int    `json:"old_camp"`
	Camp    int    `json:"camp"`
	MineId  int    `json:"mine_id"`
	Name    string `json:"name"`
	IconId  int    `json:"icon_id"`
	Level   int    `json:"level"`
}

// 阵营变更
type S2C_MineMatchFail struct {
	Cid          string `json:"cid"`
	MineId       int    `json:"mine_id"`
	ReturnMineId int    `json:"return_mine_id"`
	FightCd      int64  `json:"fight_cd"`
}

// 匹配成功
type S2C_MineMatchOk struct {
	Cid         string `json:"cid"`
	Uid         int64  `json:"uid"`
	MineId      int    `json:"mine_id"`      // 矿点ID
	IconId      int    `json:"icon_id"`      // 头像
	Portrait    int    `json:"portrait"`     // 边框
	Name        string `json:"name"`         // 姓名
	Fight       int64  `json:"fight"`        // 战力
	EncourageLv int    `json:"encourage_lv"` // 鼓舞等级
	Team        []int  `json:"team"`         // 阵容
	Star        []int  `json:"star"`         // 星级
	MainTalent  []int  `json:"maintalent"`   // 主天赋等级
	Camp        int    `json:"camp"`         // 阵营
	Level       int    `json:"level"`        // 等级
}

//// 斗技场战报
//type S2C_PvpFights struct {
//	Cid       string            `json:"cid"`
//	FightInfo map[int]*PvpFight `json:"fight_info,omitempty"` // 战报信息
//}

// 开始pvp战斗
type S2C_BeginPvp struct {
	Cid       string        `json:"cid"`
	RandNum   int64         `json:"rand_num"` // 随机数
	FightInfo *JS_FightInfo `json:"fightinfo"`
}

// 删除多余战报Id
type S2C_RemoveMineFight struct {
	Cid     string `json:"cid"`
	FightId int64  `json:"fight_id"`
}

//战斗开始倒计时
type S2C_ArmyTeamFightCountdown struct {
	Cid  string `json:"cid"`
	Flag int    `json:"flag"` //!倒计时状态 flag为1表示倒计时开始 flag为0表示结束
}

//队伍解散
type S2C_ArmyTeamDismiss struct {
	Cid string `json:"cid"`
}

type S2C_CompareTeamFight struct {
	Cid        string             `json:"cid"`
	HeroFight  []int64            `json:"herofight"`
	HeroKeyId  []int              `json:"herokeyid"`
	HydraFight int64              `json:"hydrafight"`
	Fight      int64              `json:"fight"`
	attr       map[int]*Attribute `json:"attr"`
	FightInfo  *JS_FightInfo      `json:"fightinfo"`
}

type S2CUnionTime struct {
	Cid       string `json:"cid"`
	StartTime int64  `json:"start_time"` // 开启时间
	EndTime   int64  `json:"end_time"`   // 结束时间
}

type UnionAttendInfo struct {
	UnionId    int    `json:"union_id"`
	IconId     int    `json:"icon_id"`
	UnionName  string `json:"union_name"`
	Fight      int64  `json:"fight"`
	MasterName string `json:"master_name"`
	Camp       int    `json:"camp"`
}

type S2C_UnionAttendInfo struct {
	Cid          string                `json:"cid"`
	StateId      int                   `json:"state_id"`      // 州Id
	Info         []*UnionAttendInfo    `json:"info"`          // 参加信息
	StateProcess map[int]*StateProcess `json:"state_process"` // 占领进度
	TotalProcess int                   `json:"total_process"` // 总进度
}

type StateProcess struct {
	Camp    int `json:"camp"`
	Process int `json:"process"`
}

type S2CStateAward struct {
	Cid        string     `json:"cid"`
	AwardState int        `json:"award_state"`
	Items      []PassItem `json:"items"`
}

type S2CUnionCDTime struct {
	Cid    string `json:"cid"`
	CDTime int64  `json:"cdtime"`
}

type StateTakenInfo struct {
	StateId   int    `json:"state_id"`
	UnionId   int    `json:"union_id"`
	Icon      int    `json:"icon"`
	UnionName string `json:"union_name"`
	State     int    `json:"state"`
	Attenders []int  `json:"attenders"`
	Rank      int    `json:"rank"`
}

// 获取领地占领情况
type S2CStateTakenInfo struct {
	Cid  string            `json:"cid"`
	Info []*StateTakenInfo `json:"info"`
}

// 报名信息
type S2CAttendInfo struct {
	Cid           string `json:"cid"`
	StateId       int    `json:"state_id"`        // 当前州Id
	Master        int    `json:"master"`          // 是否是军团长
	CallState     int    `json:"call_state"`      // 州是否宣战
	AttendState   int    `json:"attend_state"`    // 是否参战
	SelfCallState int    `json:"self_call_state"` // 我的军团是否宣战
	StateAction   int    `json:"state_action"`    // 州的状态
}

type UnionFinalInfo struct {
	Team       int    `json:"team"`        // 队伍编号
	Icon       int    `json:"icon"`        // 军团Id
	Name       string `json:"name"`        // 军团名字
	Lv         int    `json:"lv"`          // 军团等级
	MasterName string `json:"master_name"` // 军团长名字
	Camp       int    `json:"camp"`        // 阵营
	Fight      int64  `json:"fight"`       // 战力值
	Fail       int    `json:"fail"`        // 0 成功 1失败
	HasData    int    `json:"has_data"`    // 0 没参加 1有参加
}

type S2CUnionBattleInfo struct {
	Cid         string            `json:"cid"`
	Info        []*UnionFinalInfo `json:"info"`
	StateId     int               `json:"state_id"`
	ActionState int               `json:"action_state"`
	Team        int               `json:"team"`
	Round       int               `json:"round"`       // 当前打到第几轮
	Group       int               `json:"group"`       // 当前打到第几组
	RoundNum    int               `json:"round_num"`   // 当前打几轮
	WinnerId    int               `json:"winner_id"`   // 占领的军团ID
	WinnerName  string            `json:"winner_name"` // 占领的军团名字
	WaitTime    int64             `json:"wait_time"`   // 下场的开战时间
}

type S2CUnionFightChange struct {
	Cid         string `json:"cid"`
	WinTeam     int    `json:"win_team"`     // 成功队伍
	LoseTeam    int    `json:"lose_team"`    // 失败队伍
	StateId     int    `json:"state_id"`     // 州Id
	Final       int    `json:"final"`        // 是否是最终结果
	TheRound    int    `json:"the_round"`    // 当前第几轮
	Reason      int    `json:"reason"`       // 1.结算 2.状态切换
	WinMVP      int64  `json:"win_mvp"`      // 胜利方MVP
	LoseMVP     int64  `json:"lose_mvp"`     // 失败方MVP
	WinMVPKill  int    `json:"win_mvpkill"`  // 胜利方MVP
	LoseMVPKill int    `json:"lose_mvpkill"` // 失败方MVP
}

type UnionFightMsg struct {
	UID       int64         `json:"uid"`       // 玩家UID, ok
	Icon      int           `json:"icon"`      // 图标, ok
	Fight     int64         `json:"fight"`     // 战力, ok
	Level     int           `json:"level"`     // 等级, ok
	Camp      int           `json:"camp"`      // 阵营, ok
	Encourage int           `json:"encourage"` // 鼓舞次数, ok
	FightBase *JS_FightBase `json:"fightteam"` // 战斗数据, 有可能为空, ok
}

// 战斗界面
// UnignFight56Init
type S2C_UnionFightInit struct {
	Cid         string              `json:"cid"`          // ok
	StateId     int                 `json:"state_id"`     // 州Id, ok
	ActionState int                 `json:"action_state"` // 当前军团战状态, ok
	EnterPlayer int64               `json:"enterplayer"`  // 玩家UID, ok
	IndexTeam   [2]int              `json:"indexteam"`    // 左右两边是谁和谁打的下标, ok
	FightTeam   [2]int              `json:"fight_team"`   // 当前是哪个队在打, ok
	GroupNum    [2]int              `json:"group_num"`    // 两边总人数, ok
	Result      [5]int              `json:"result"`       // 战斗结果，0-还未出结果，1左胜利，2右胜利, ok
	MatchGroup  int                 `json:"match_group"`  // 自己所在组, 没有进入队列为0, ok
	Msg         [2][]*UnionFightMsg `json:"msg"`          // 两边数据, ok
	UnionName   [2]string           `json:"union_name"`   // 军团名字, ok
	UnionMaster [2]string           `json:"union_master"` // 军团长名字, ok
	WaitTime    int64               `json:"wait_time"`    // 等待时间
	Camp        [2]int              `json:"camp"`         // 阵营信息
	PkInfo      [2]*JS_FightBase    `json:"pk_info"`      // pk信息
	HP          [2][5]int           `json:"hp"`           //! 剩余HP
	UnionIcon   [2]int              `json:"union_icon"`   // 军团ICON    //20190428 by zy
}

//! 战报
type UnionFightFinal struct {
	Id     int64                  `json:"id"`
	Info   [2]*SonUnionFightFinal `json:"info"`
	Result int                    `json:"result"`
}

//! 战报
type SonUnionFightFinal struct {
	Name      string `json:"name"`
	Icon      int    `json:"icon"`
	Level     int    `json:"level"`
	Fight     int64  `json:"fight"`
	Encourage int    `json:"encourage"`
}

// 战报界面
type S2C_UnionGetFights struct {
	Cid            string             `json:"cid"`
	Camp           [2]int             `json:"camp"`
	UnionId        [2]int             `json:"unionid"`
	UnionIcon      [2]int             `json:"unionicon"`
	Results        [5]int             `json:"results"`
	PlayerNum      [2]int             `json:"playernum"`
	Final          []*UnionFightFinal `json:"final"`
	PlayerNumGroup [2]int             `json:"playernumgroup"`
}

type S2C_UnionTeamRecord struct {
	Cid       string    `json:"cid"`
	TeamMatch [3][2]int `json:"team_match"` // 队伍匹配编号
}

type S2C_SetGuild struct {
	RedIcon int `json:"red_icon"`
}

type S2C_SynGuild struct {
	RedIcon int `json:"red_icon"`
}

type S2C_SetUserSignature struct {
	Cid       string `json:"cid"`
	Signature string `json:"signature"`
}

type S2C_SetUserLanguage struct {
	Cid      string `json:"cid"`
	Language int    `json:"language"`
}

type S2C_SetUserNationality struct {
	Cid         string `json:"cid"`
	Nationality int    `json:"nationality"`
}

//! 军团战战斗-序列
type S2C_UnionFightRun struct {
	Cid        string           `json:"cid"`
	Index      int              `json:"index"`      // 当前是哪个组, ok
	IndexTeam  [2]int           `json:"indexteam"`  // 哪两个队伍, ok
	IndexMax   [2]int           `json:"indexmax"`   // 队列人数, ok
	TeamResult [5]int           `json:"teamresult"` // 队伍输赢, ok
	Time       int64            `json:"time"`       // 战斗到期时间, ok
	Result     int              `json:"result"`     // 战斗结果，0开始战斗，1-攻方胜利，2-守方胜利
	Win        [2]int           `json:"win"`        // 胜负, ok
	Hp         [2][5]int        `json:"hp"`         // 剩余血量, ok
	FightTeam  [2]*JS_FightBase `json:"fightteam"`  // 战斗数据, ok
	TheRound   int              `json:"the_round"`  // 当前轮数
	MVP        [2]*JS_FightBase `json:"mvp"`        // mvp
	MVPKill    [2]int           `json:"mvpkill"`    // mvpkill
	BossHP     [2]int           `json:"bosshp"`     //! 巨兽的HP
	HeroId     [2][5]int        `json:"heroid"`     // 剩余血量对应的英雄ID
}

//! 战斗序列结束(56ok: 发送奖励, 56end: 战斗结束)
type S2C_UnionFightEnd struct {
	Cid        string `json:"cid"`
	Win        [2]int `json:"win"`        // 胜负
	TeamResult [5]int `json:"teamresult"` // 队伍结果
	State      int    `json:"state"`      // 战斗状态
	FightTime  int64  `json:"fighttime"`  // 结束时间
}

//调试消息
type S2C_Guides struct {
	Cid    string `json:"cid"`
	Guides []int  `json:"guides"`
}

//加入队伍
type S2C_SetFlag struct {
	Cid       string `json:"cid"`
	TeamId    int64  `json:"team_id"`
	AutoEnter int    `json:"auto_enter"`
}

//挂机系统
type S2C_OnHookInfo struct {
	Cid       string `json:"cid"`
	Time      int64  `json:"time"`      //当前累计时间
	Stage     int    `json:"stage"`     //当前挂机关卡
	HangUp    int    `json:"hangup"`    //挂机掉落组
	FastTimes int    `json:"fasttimes"` //快速挂机次数
	StageTime int64  `json:"stagetime"` //
}

/*
type S2C_UpdateHydraTask struct {
	Cid      string          `json:"cid"`
	TaskInfo []HydraTaskInfo `json:"taskinfo"`
}

*/

//单个系统的同步
type S2C_StatistcsOne struct {
	Cid        string        `json:"cid"`
	Type       int           `json:"type"`       //唯一标识
	State      int           `json:"state"`      //当前状态进度 初始为0 表示不可用
	Value      int           `json:"value"`      //当前值
	Score      int           `json:"score"`      //目前积分
	RewardInfo map[int]int64 `json:"rewardinfo"` //积分奖励领取情况
}

type S2CAwardStatistics struct {
	Cid        string        `json:"cid"`
	RewardInfo map[int]int64 `json:"rewardinfo"` //积分奖励领取情况
	Award      []PassItem    `json:"award"`
}

type S2C_AwardOnHook struct {
	Cid               string     `json:"cid"`
	Time              int64      `json:"time"`              //挂机时间
	GetItem           []PassItem `json:"getitem"`           //获得了哪些物品
	GetTime           int64      `json:"gettime"`           //本次奖励时长
	GetPrivilegeItems []PassItem `json:"getprivilegeitems"` //特权物品
	GetMonthItems     []PassItem `json:"getmonthitems"`     //月卡特权物品
	GetActivityItems  []PassItem `json:"getactivityitems"`  //活动奖励
}

type S2C_OnHookFast struct {
	Cid               string     `json:"cid"`
	GetItem           []PassItem `json:"getitem"`           //获得了哪些物品
	CostItem          []PassItem `json:"costitem"`          //消耗
	FastTimes         int        `json:"fasttimes"`         //
	GetPrivilegeItems []PassItem `json:"getprivilegeitems"` //特权物品
	GetMonthItems     []PassItem `json:"getmonthitems"`     //月卡特权物品
	GetActivityItems  []PassItem `json:"getactivityitems"`  //活动奖励
}

type S2C_OnHookStage struct {
	Cid       string `json:"cid"`
	Stage     int    `json:"stage"`     //选择哪个关卡
	HangUp    int    `json:"hangup"`    //挂机掉落组
	StageTime int64  `json:"stagetime"` //
}

//! 军团鼓舞
type S2C_UnionInspireNotice struct {
	Cid     string `json:"cid"`
	Uid     int64  `json:"uid"`
	UnionId int    `json:"unionid"`
	Level   int    `json:"level"`
}

//! 获得关卡任务
type S2C_GetMission struct {
	Cid       string `json:"cid"`
	Taskid    int    `json:"taskid"`    // 任务Id
	Tasktypes int    `json:"tasktypes"` // 任务类型
	Plan      int    `json:"plan"`      // 进度
	Finish    int    `json:"finish"`    // 是否完成
}

//目标系统 爵位 任务同步
type S2C_UpdateMission struct {
	Cid      string        `json:"cid"`
	TaskInfo []JS_TaskInfo `json:"taskinfo"`
}

type S2C_EquipUpLv struct {
	Cid        string     `json:"cid"`
	ItemKeyId  []int      `json:"itemkeyid"`  //!消耗的装备类key
	Equip      *Equip     `json:"equip"`      //!
	CostItem   []PassItem `json:"costitem"`   //!  消耗的物品类
	ReturnItem []PassItem `json:"returnitem"` //!  返还的物品
}

//英雄升级
type S2C_HeroLvUp struct {
	Cid         string     `json:"cid"`
	Items       []PassItem `json:"items"`
	HeroKeyId   int        `json:"herokeyid"`
	HeroLv      int        `json:"herolv"`
	HeroBreakId int        `json:"herobreakid"` //突破编号   hero_break表：id
	StarItem    *StarItem  `json:"staritem"`    //! 升星信息
}

type S2C_HeroLvUpTo struct {
	Cid         string     `json:"cid"`
	Items       []PassItem `json:"items"`
	HeroKeyId   int        `json:"herokeyid"`
	HeroLv      int        `json:"herolv"`
	HeroBreakId int        `json:"herobreakid"` //突破编号   hero_break表：id
	StarItem    *StarItem  `json:"staritem"`    //! 升星信息
}

type S2C_HeroCompound struct {
	Cid     string     `json:"cid"`
	Costs   []PassItem `json:"costs"`   //! 消耗物品
	GetHero []PassItem `json:"gethero"` //! （仅显示用）
}

type S2C_HeroAutoFire struct {
	Cid      string `json:"cid"`
	AutoFire int    `json:"autofire"` //!
}

type S2C_OnlineTip struct {
	Cid   string `json:"cid"`
	Param int    `json:"param"` //!   0被挤掉的号    1挤的号
}

type S2C_HeroGetHandBook struct {
	Cid      string     `json:"cid"`
	HeroId   int        `json:"heroid"`   //!
	State    int        `json:"state"`    //!
	StarMax  int        `json:"starmax"`  //!
	GetItems []PassItem `json:"getitems"` //!
}

type S2C_BaseHeroSet struct {
	Cid         string     `json:"cid"`
	BaseHeroSet []*NewHero `json:"baseheroset"` //!
}

type S2C_GetHireList struct {
	Cid      string              `json:"cid"`
	HireList map[int][]*HireHero `json:"hirelist"` //!  可租借的别人的  key  英雄ID
	SelfList map[int]*HireHero   `json:"selflist"` //!  自己的,从这个里面可以看到谁申请过  key  KEYID
}

type S2C_HireApply struct {
	Cid       string        `json:"cid"`
	ApplyHero []*HireHero   `json:"applyhero"`
	HireHero  []*FriendHero `json:"hirehero"`
}

type S2C_DeleteHireList struct {
	Cid            string      `json:"cid"`
	DeleteHireList []*HireHero `json:"deletehire"`
}

type S2C_HireStateUpdate struct {
	Cid      string    `json:"cid"`
	HireHero *HireHero `json:"hirehero"`
}

type S2C_GetPlayerTeam struct {
	Cid       string        `json:"cid"`
	FightInfo *JS_FightInfo `json:"fightinfo"`
}

type S2C_HeroBackOpen struct {
	Cid      string `json:"cid"`
	BackOpen int    `json:"backopen"` //!
}

//合成
type S2C_RuneCompose struct {
	Cid     string        `json:"cid"`
	Id      int           `json:"id"`      //! 配方ID
	Num     int           `json:"num"`     // 实际合成次数
	InItem  map[int]*Item `json:"initem"`  //消耗物品
	OutItem map[int]*Item `json:"outitem"` //产出物品
}

/*
//! 同步神兽信息
type S2C_HydraInfo struct {
	Cid        string             `json:"cid"`
	HydraInfos map[int]*HydraInfo `json:"hydrainfo"` //! 神兽信息
}

type S2CTakeHydraTask struct {
	Cid      string                 `json:"cid"`
	TaskInfo map[int]*HydraTaskInfo `json:"taskinfo"`
	Level    int                    `json:"level"`
	Award    []PassItem             `json:"award"`
}

type S2CTakeHydra struct {
	Cid        string             `json:"cid"`
	HydraInfos map[int]*HydraInfo `json:"hydrainfo"`
	Items      []PassItem         `json:"items"`
}

//英雄升级
type S2C_HydraLvUp struct {
	Cid         string     `json:"cid"`
	Action      int        `json:"action"`
	Items       []PassItem `json:"items"`
	HydraInfo   *HydraInfo `json:"hydrainfo"`
	SpecialStop int        `json:"specialstop"`
}

*/

//! 关卡
type S2C_ClientSign struct {
	Cid  string      `json:"cid"`
	Sign map[int]int `json:"sign"`
}

type S2C_PitInfo struct {
	Cid      string   `json:"cid"`
	PitKeyId int      `json:"pitkeyid"`
	Pitinfo  *PitInfo `json:"pitinfo"`
}

type S2C_PitInfoAll struct {
	Cid         string        `json:"cid"`
	PitInfoShow []PitInfoShow `json:"pitinfoshow"`
}

type PitInfoShow struct {
	Id          int   `json:"id"`
	PitId       int   `json:"pitid"`       //ID生成地牢关卡
	PitType     int   `json:"pittype"`     //TYPE地牢产生规则
	PitKeyId    int   `json:"pitkeyid"`    //地牢唯一标识
	EndTime     int64 `json:"endtime"`     //结束时间
	State       int   `json:"state"`       //地牢状态
	FinishTimes int   `json:"finishtimes"` //完成
	AllTimes    int   `json:"alltimes"`    //总共次数
}

//! 完成事件
type S2C_PitFinishEvent struct {
	Cid      string     `json:"cid"`
	PitKeyId int        `json:"pitkeyid"`
	PitEvent *PitEvent  `json:"pitevent"`
	Item     []PassItem `json:"item"`
	Cost     []PassItem `json:"cost"`
	Buff     []int      `json:"buff"`
}

//时光之巅
type S2C_InstanceInfo struct {
	Cid              string                `json:"cid"`
	InstanceInfo     map[int]*InstanceInfo `json:"instanceinfo"`     //整体进度
	NowInstanceState *NowInstanceState     `json:"nowinstancestate"` //当前攻略的地图
}

//新地牢消息
type S2C_NewPitInfo struct {
	Cid                  string                  `json:"cid"`
	NewPitInfo           map[int]*NewPitInfo     `json:"newpitinfo"`           // 关卡信息
	UserPitInfo          *UserPitInfo            `json:"userpitinfo"`          //
	UserPitInfoFightInfo map[int64]*JS_FightInfo `json:"userpitinfofightinfo"` //
	Shop                 []*NewPitShop           `json:"shop"`
}

type S2C_NewPitFinishEvent struct {
	Cid         string        `json:"cid"`
	Id          int           `json:"id"`          //!  关卡ID
	Item        []PassItem    `json:"item"`        //获得物品
	Cost        []PassItem    `json:"cost"`        //消耗物品
	PitInfo     []*NewPitInfo `json:"pitinfo"`     //更新关卡
	UserPitInfo *UserPitInfo  `json:"userpitinfo"` //更新玩家信息
	GetHero     *NewHero      `json:"gethero"`     //获得英雄
}

type S2C_NewPitFinishReward struct {
	Cid               string        `json:"cid"`
	Id                int           `json:"id"`                //!  关卡ID
	Item              []PassItem    `json:"item"`              //获得物品
	PitInfo           []*NewPitInfo `json:"pitinfo"`           //更新关卡
	NowPitId          int           `json:"nowpitid"`          //当前位置
	NowAimId          int           `json:"nowaimid"`          //当前选择
	GetPrivilegeItems []PassItem    `json:"getprivilegeitems"` //特权物品 新增
}

type S2C_NewPitShopBuy struct {
	Cid         string          `json:"cid"`
	GetItems    []PassItem      `json:"getitems"`    //获得
	CostItems   []PassItem      `json:"costitems"`   //消耗
	NewShopInfo *JS_NewShopInfo `json:"newshopinfo"` //商品同步
}

type S2C_NewPitFinishSpring struct {
	Cid       string                   `json:"cid"`
	Id        int                      `json:"id"`        //!  关卡ID
	PitInfo   []*NewPitInfo            `json:"pitinfo"`   //更新关卡
	NowPitId  int                      `json:"nowpitid"`  //当前位置
	NowAimId  int                      `json:"nowaimid"`  //当前选择
	HeroState map[int]*NewPitHeroState `json:"herostate"` //保存英雄血量和能量
}

type S2C_NewPitFinishMystery struct {
	Cid       string                   `json:"cid"`
	Id        int                      `json:"id"`        //!  关卡ID
	PitInfo   []*NewPitInfo            `json:"pitinfo"`   //更新关卡
	NowPitId  int                      `json:"nowpitid"`  //当前位置
	NowAimId  int                      `json:"nowaimid"`  //当前选择
	HeroState map[int]*NewPitHeroState `json:"herostate"` //保存英雄血量和能量
}

type S2C_NewPitFinishShop struct {
	Cid      string        `json:"cid"`
	Id       int           `json:"id"`       //!  关卡ID
	Item     []PassItem    `json:"item"`     //获得物品
	Cost     []PassItem    `json:"cost"`     //消耗物品
	PitInfo  []*NewPitInfo `json:"pitinfo"`  //更新关卡
	NowPitId int           `json:"nowpitid"` //当前位置
	NowAimId int           `json:"nowaimid"` //当前选择
}

type S2C_NewPitFinishBattle struct {
	Cid               string                   `json:"cid"`
	Id                int                      `json:"id"`                //!  关卡ID
	Item              []PassItem               `json:"item"`              //获得物品
	PitInfo           []*NewPitInfo            `json:"pitinfo"`           //更新关卡
	NowPitId          int                      `json:"nowpitid"`          //当前位置
	NowAimId          int                      `json:"nowaimid"`          //当前选择
	BuffChoose        []int                    `json:"buffchoose"`        //这个有值的时候需要选择BUFF才可以继续
	HeroState         map[int]*NewPitHeroState `json:"herostate"`         //保存英雄血量和能量
	GetPrivilegeItems []PassItem               `json:"getprivilegeitems"` //特权物品
	Shop              []*NewPitShop            `json:"shop"`              //商店
}

type S2C_NewPitFinishTreasure struct {
	Cid       string                   `json:"cid"`
	Id        int                      `json:"id"`        //!  关卡ID
	Item      []PassItem               `json:"item"`      //获得物品
	PitInfo   []*NewPitInfo            `json:"pitinfo"`   //更新关卡
	NowPitId  int                      `json:"nowpitid"`  //当前位置
	NowAimId  int                      `json:"nowaimid"`  //当前选择
	HeroState map[int]*NewPitHeroState `json:"herostate"` //保存英雄血量和能量
}

type S2C_NewPitFinishSoul struct {
	Cid      string        `json:"cid"`
	Id       int           `json:"id"`       //!  关卡ID
	PitInfo  []*NewPitInfo `json:"pitinfo"`  //更新关卡
	NowPitId int           `json:"nowpitid"` //当前位置
	NowAimId int           `json:"nowaimid"` //当前选择
	GetHero  *NewHero      `json:"gethero"`  //获得英雄
}

type S2C_NewPitFinishEvil struct {
	Cid       string                   `json:"cid"`
	Id        int                      `json:"id"`        //!  关卡ID
	PitInfo   []*NewPitInfo            `json:"pitinfo"`   //更新关卡
	NowPitId  int                      `json:"nowpitid"`  //当前位置
	NowAimId  int                      `json:"nowaimid"`  //当前选择
	GetHero   *NewHero                 `json:"gethero"`   //获得英雄
	HeroState map[int]*NewPitHeroState `json:"herostate"` //保存英雄血量和能量
}

type S2C_NewPitNowAim struct {
	Cid      string `json:"cid"`
	NowAimId int    `json:"nowaimid"` //只同步选择
}

type S2C_NewPitChooseBuff struct {
	Cid      string        `json:"cid"`
	Buff     []int         `json:"buff"`     //选择完成注意清空choose列表，感觉不需要通知这个NIL
	PitInfo  []*NewPitInfo `json:"pitinfo"`  //更新关卡
	NowPitId int           `json:"nowpitid"` //当前位置
	NowAimId int           `json:"nowaimid"` //当前选择
}

type S2C_NewPitUseItem struct {
	Cid       string                   `json:"cid"`
	Cost      []PassItem               `json:"cost"`
	HeroState map[int]*NewPitHeroState `json:"herostate"`
}

type S2C_NewPitFinishNow struct {
	Cid                  string                  `json:"cid"`
	PitInfo              map[int]*NewPitInfo     `json:"pitinfo"`              //更新这个关卡
	UserPitInfo          *UserPitInfo            `json:"userpitinfo"`          //更新玩家信息
	Item                 []PassItem              `json:"item"`                 //获得物品
	UserPitInfoFightInfo map[int64]*JS_FightInfo `json:"userpitinfofightinfo"` //
}

type S2C_SetClientSign struct {
	Cid   string `json:"cid"`
	Key   int    `json:"key"`   //!
	Value int    `json:"value"` //!
}

type S2C_InstanceStart struct {
	Cid              string            `json:"cid"`
	NowInstanceState *NowInstanceState `json:"nowinstancestate"`
}

type S2C_InstanceMove struct {
	Cid       string `json:"cid"`
	Id        int    `json:"id"`
	Row       int    `json:"row"`       //!
	Col       int    `json:"col"`       //!
	RowShadow []int  `json:"rowshadow"` //!
	ColShadow []int  `json:"colshadow"` //!
}

type S2C_InstanceBattle struct {
	Cid       string                     `json:"cid"`
	ThingInfo *ThingInfo                 `json:"thinginfo"` //更新关卡
	HeroState map[int]*InstanceHeroState `json:"herostate"` //保存英雄血量和能量
}

type S2C_InstanceFriend struct {
	Cid       string     `json:"cid"`
	ThingInfo *ThingInfo `json:"thinginfo"` //更新关卡
	GetHero   *NewHero   `json:"gethero"`   //获得英雄
}

type S2C_InstanceAdd struct {
	Cid       string                     `json:"cid"`
	ThingInfo *ThingInfo                 `json:"thinginfo"` //更新关卡
	HeroState map[int]*InstanceHeroState `json:"herostate"` //保存英雄血量和能量
}

type S2C_InstanceReborn struct {
	Cid       string                     `json:"cid"`
	ThingInfo *ThingInfo                 `json:"thinginfo"` //更新关卡
	HeroState map[int]*InstanceHeroState `json:"herostate"` //保存英雄血量和能量
}

type S2C_InstanceSwitch struct {
	Cid         string       `json:"cid"`
	ThingInfo   *ThingInfo   `json:"thinginfo"`   //更新关卡
	ThingInfoEx []*ThingInfo `json:"thinginfoex"` //更新关卡
}

type S2C_FinishThing struct {
	Cid          string        `json:"cid"`
	ThingInfo    *ThingInfo    `json:"thinginfo"` //更新关卡
	InstanceInfo *InstanceInfo `json:"instanceinfo"`
	Item         []PassItem    `json:"item"` //获得物品
	Id           int           `json:"id"`
	Row          int           `json:"row"`       //!
	Col          int           `json:"col"`       //!
	RowShadow    []int         `json:"rowshadow"` //!
	ColShadow    []int         `json:"colshadow"` //!
}

type S2C_InstanceChooseBuff struct {
	Cid       string     `json:"cid"`
	ThingInfo *ThingInfo `json:"thinginfo"` //更新关卡
	Buff      []int      `json:"buff"`      //累计有哪些BUFF
}

type S2C_InstanceMakeBuff struct {
	Cid       string     `json:"cid"`
	ThingInfo *ThingInfo `json:"thinginfo"` //更新关卡
}

type S2C_InstanceUpdate struct {
	Cid       string     `json:"cid"`
	ThingInfo *ThingInfo `json:"thinginfo"` //更新关卡
}

type S2C_GetPvPFightInfo struct {
	Cid       string           `json:"cid"`       //! cid
	FightId   int              `json:"fightid"`   //! 战报ID
	FightInfo [2]*JS_FightInfo `json:"fightinfo"` //战斗数据
}

type S2C_SupportHeroSet struct {
	Cid     string `json:"cid"`     //! cid
	Index   int    `json:"index"`   // index
	HeroKey int    `json:"herokey"` // herokey
	CDTime  int64  `json:"cdtime"`  // cd时间
}

type S2C_SupportHeroCancel struct {
	Cid     string `json:"cid"`     //! cid
	HeroKey int    `json:"herokey"` // herokey
}

type S2C_SupportHeroInfo struct {
	Cid    string         `json:"cid"`    //! cid
	Info   []*SupportHero `json:"info"`   // info
	HeroID int            `json:"heroid"` // 英雄id
}
type S2C_SupportMyHeroInfo struct {
	Cid  string           `json:"cid"`  //! cid
	Info []*MySupportHero `json:"info"` // info
}

type S2C_EntanglementUse struct {
	Cid       string                   `json:"cid"`       // cid
	Index     int                      `json:"index"`     // 栏位
	Type      int                      `json:"type"`      // 羁绊类型
	MasterUid int64                    `json:"masteruid"` // 被借玩家的uid
	HeroKey   int                      `json:"herokey"`   // herokey
	Info      []*JS_Entanglement       `json:"info"`      // 数据
	Property  *JS_EntanglementProperty `json:"property"`
}

type S2C_EntanglementInfo struct {
	Cid      string                     `json:"cid"` // cid
	Info     []*JS_Entanglement         `json:"info"`
	Property []*JS_EntanglementProperty `json:"property"`
}

type S2C_EntanglementAutoUse struct {
	Cid       string             `json:"cid"`       // cid
	Index     []int              `json:"index"`     // 栏位
	Type      int                `json:"type"`      // 羁绊类型
	MasterUid []int64            `json:"masteruid"` // 被借玩家的uid
	HeroKey   []int              `json:"herokey"`   // herokey
	Info      []*JS_Entanglement `json:"info"`      // 数据
	//Property  []*JS_EntanglementProperty `json:"property"`
}

//type S2C_EntanglementCancel struct {
//	Cid     string `json:"cid"`     //! cid
//	Uid     int64  `json:"uid"`     // 被借玩家的uid
//	HeroKey int    `json:"herokey"` // herokey
//}

type S2C_RewardSet struct {
	Cid      string       `json:"cid"`      // cid
	ID       int          `json:"id"`       // 任务id
	Uids     []int64      `json:"uids"`     // 被借玩家的uid
	HeroKeys []int        `json:"herokeys"` // herokey
	Info     []*JS_Reward `json:"info"`     // info
}

type S2C_RewardGet struct {
	Cid       string                    `json:"cid"`       // cid
	ID        int                       `json:"id"`        // 任务id
	Item      []PassItem                `json:"item"`      // 获得道具
	Level     int                       `json:"level"`     // 等级
	TaskCount [REWARD_TASK_TYPE_MAX]int `json:"taskcount"` // 任务进度
}

type S2C_RewardGetAll struct {
	Cid       string                    `json:"cid"`       // cid
	ID        []int                     `json:"id"`        // 任务id
	Item      []PassItem                `json:"item"`      // 获得道具
	Level     int                       `json:"level"`     // 等级
	TaskCount [REWARD_TASK_TYPE_MAX]int `json:"taskcount"` // 任务进度
}

type S2C_RewardRedPoint struct {
	Cid      string `json:"cid"`      // cid
	RedPoint bool   `json:"redpoint"` // 红点
}

type S2C_RewardInfo struct {
	Cid       string                    `json:"cid"`       // cid
	Info      []*JS_Reward              `json:"info"`      //
	Level     int                       `json:"level"`     // 等级
	TaskCount [REWARD_TASK_TYPE_MAX]int `json:"taskcount"` // 任务进度
}

type S2C_RewardRefresh struct {
	Cid       string       `json:"cid"`  // cid
	Info      []*JS_Reward `json:"info"` //
	CostItems []PassItem   `json:"costitems"`
}

//! 排行领取信息
type S2C_RankTaskGetStateInfo struct {
	Cid         string      `json:"cid"`
	ID          int         `json:"id"`
	GetState    []*GetState `json:"getstate"`
	FinishState []int       `json:"finishstate"`
}

//! 完成id的玩家数据
type S2C_RankTaskGetPlayerInfo struct {
	Cid        string            `json:"cid"`
	ID         int               `json:"id"`
	PlayerRank []*RankPlayerInfo `json:"playerrank"`
}

//! 完成type的玩家数据
type S2C_RankTaskGetTypePlayerInfo struct {
	Cid        string            `json:"cid"`
	Type       int               `json:"type"`
	ID         []int             `json:"id"`
	PlayerRank []*RankPlayerInfo `json:"playerrank"`
}

//! 排行任务领取奖励
type S2C_RankTaskAward struct {
	Cid   string     `json:"cid"`
	ID    int        `json:"id"`
	Items []PassItem `json:"items"`
}

//! 排行任务红点推送
type S2C_RankTaskRedPoint struct {
	Cid string `json:"cid"`
	ID  int    `json:"id"`
}

type S2C_ResonanceCrystalSet struct {
	Cid     string `json:"cid"`     // cid
	Index   int    `json:"index"`   // index
	HeroKey int    `json:"herokey"` // herokey
}

type S2C_ResonanceCrystalFight struct {
	Cid      string `json:"cid"`      // cid
	FightAll int64  `json:"fightall"` // index
}

type S2C_ResonanceCrystaCancel struct {
	Cid     string `json:"cid"`     // cid
	Index   int    `json:"index"`   // index
	EndTime int64  `json:"endtime"` // endtime
}

type S2C_ResonanceCrystaInfo struct {
	Cid            string            `json:"cid"`            // cid
	ResonanceCount int               `json:"resonancecount"` // 最大共鸣人数
	PriestsHeros   []int             `json:"priestsheros"`   // 祭司英雄
	ResonanceHeros []*ResonanceHeros `json:"resonanceheros"` // 共鸣英雄
	FightAll       int64             `json:"fightall"`       // 共鸣总战力
	Level          int               `json:"level"`          // 用户法阵核心等级
}
type S2C_ResonanceCrystaAddResonanceCount struct {
	Cid            string     `json:"cid"`  // cid
	Type           int        `json:"type"` // type
	Items          []PassItem `json:"items"`
	ResonanceCount int        `json:"resonancecount"` // 最大共鸣人数
}

type S2C_ResonanceCrystaUpdateResonance struct {
	Cid   string  `json:"cid"`   // cid
	Heros []*Hero `json:"heros"` // 英雄
}

type S2C_ResonanceCrystaUpdatePriests struct {
	Cid   string `json:"cid"`   // cid
	Heros []int  `json:"heros"` // 英雄
}
type S2C_ResonanceCrystalCleanCD struct {
	Cid   string          `json:"cid"`   // cid
	Index int             `json:"index"` // index
	Items []PassItem      `json:"items"`
	Data  *ResonanceHeros `json:"data"`
}

//! 获得竞技场信息
type S2C_ArenaInfo struct {
	Cid string `json:"cid"`
	//Info *JS_FightInfo `json:"info"`
	Rank      int   `json:"rank"`
	StartTime int64 `json:"starttime"`
	EndTime   int64 `json:"endtime"`
}

//! 获得竞技场
type S2C_GetArenaFight struct {
	Cid        string          `json:"cid"`
	Enemy      []*JS_FightInfo `json:"enemy"`
	Point      []int64         `json:"point"`
	FreeCount  int             `json:"freecount"`  // 免费次数
	FightCount int             `json:"fightcount"` // 已战斗次数
}

//! 竞技场排行榜
type S2C_ArenaTop struct {
	Cid string       `json:"cid"`
	Top []*Js_ActTop `json:"top"`
}

// 斗技场战报
type S2C_ArenaFights struct {
	Cid        string        `json:"cid"`
	FightInfo  []*ArenaFight `json:"fight_info,omitempty"` // 战报信息
	FreeCount  int           `json:"freecount"`            // 免费次数
	FightCount int           `json:"fightcount"`           // 已战斗次数
}

// 斗技场战报
type S2C_ArenaBattleDefend struct {
	Cid       string        `json:"cid"`
	FightInfo *JS_FightInfo `json:"fight_info"` // 双方数据
}

type S2C_BuyArenaCount struct {
	Cid   string     `json:"cid"`
	Items []PassItem `json:"items"`
}

type S2C_GetEnemyFightInfo struct {
	Cid       string           `json:"cid"`
	Type      int              `json:"type"`
	Uid       int64            `json:"uid"`
	FightInfo *JS_FightInfo    `json:"fight_info"` // 双方数据
	LifeTree  *JS_LifeTreeInfo `json:"lifetree"`   // 生命树
}

//! 竞技场
type S2C_ArenaStart struct {
	Cid         string        `json:"cid"`
	RandNum     int64         `json:"rand_num"` // 随机数
	FightID     int64         `json:"fightid"`
	MyFightInfo *JS_FightInfo `json:"fight_info"`
}

//! 竞技场
type S2C_ArenaBackStart struct {
	Cid         string        `json:"cid"`
	RandNum     int64         `json:"rand_num"` // 随机数
	FightID     int64         `json:"fightid"`
	MyFightInfo *JS_FightInfo `json:"fight_info"`
}

//! 战斗结束
type S2C_ArenaEnd struct {
	Cid           string           `json:"cid"`
	RandNum       int64            `json:"rand_num"`      // 随机数
	FightInfo     [2]*JS_FightInfo `json:"fightinfo"`     // 我方战力结构
	Index         int64            `json:"index"`         // index
	Result        int              `json:"result"`        //结果
	FightID       int64            `json:"fightid"`       // 战斗id
	OldRank       int              `json:"oldrank"`       // 旧排名
	NewRank       int              `json:"newrank"`       // 新排名
	MyAddPoint    int64            `json:"myaddpoint"`    // 我加的点数
	EnemyAddPoint int64            `json:"enemyaddpoint"` // 敌方加的点数
	MyPoint       int64            `json:"mypoint"`       // 我当前点数
	EnemyPoint    int64            `json:"enemypoint"`    // 敌方当前点数
	Item          []PassItem       `json:"item"`
}

//! 增加战报
type S2C_ArenaAddFightRecord struct {
	Cid string `json:"cid"`
}

//! 反击战斗结束
type S2C_ArenaFightBackEnd struct {
	Cid           string           `json:"cid"`
	RandNum       int64            `json:"rand_num"`      // 随机数
	FightInfo     [2]*JS_FightInfo `json:"fightinfo"`     // 我方战力结构
	Index         int64            `json:"index"`         // index
	Result        int              `json:"result"`        // 结果
	FightID       int64            `json:"fightid"`       // 战斗id
	OldRank       int              `json:"oldrank"`       // 旧排名
	NewRank       int              `json:"newrank"`       // 新排名
	MyAddPoint    int64            `json:"myaddpoint"`    // 我加的点数
	EnemyAddPoint int64            `json:"enemyaddpoint"` // 敌方加的点数
	MyPoint       int64            `json:"mypoint"`       // 我当前点数
	EnemyPoint    int64            `json:"enemypoint"`    // 敌方当前点数
	Item          []PassItem       `json:"item"`
}

type S2C_GetBattleInfo struct {
	Cid        string      `json:"cid"`
	FightID    int64       `json:"fightid"`
	Type       int         `json:"type"`
	Uid        int         `json:"uid"`
	BattleInfo *BattleInfo `json:"battleinfo"`
}

type S2C_GetBattleInfo2 struct {
	Cid        string        `json:"cid"`
	FightID    int64         `json:"fightid"`
	BattleInfo []*BattleInfo `json:"battleinfo"`
}

type S2C_GetBattleRecord struct {
	Cid          string        `json:"cid"`
	FightID      int64         `json:"fightid"`
	Type         int           `json:"type"`
	Uid          int           `json:"uid"`
	BattleRecord *BattleRecord `json:"battlerecord"`
}

//! 增加战报
type S2C_ArenaSpecialAddFightRecord struct {
	Cid string `json:"cid"`
}

//! 进入高阶竞技场
type S2C_ArenaSpecialEnter struct {
	Cid        string `json:"cid"`
	Class      int    `json:"class"`
	Dan        int    `json:"dan"`
	StartTime  int64  `json:"starttime"`
	EndTime    int64  `json:"endtime"`
	FreeCount  int    `json:"freecount"`  // 免费次数
	FightCount int    `json:"fightcount"` // 已战斗次数
	BuyCount   int    `json:"buycount"`   // 购买次数
	Coin       int    `json:"coin"`
	Point      int64  `json:"point"`
	ClassTime  int64  `json:"classtime"`
}

type S2C_ArenaSpecialStart struct {
	Cid         string                                `json:"cid"`
	RandNum     int64                                 `json:"rand_num"` // 随机数
	FightID     [ARENA_SPECIAL_TEAM_MAX]int64         `json:"fightid"`
	MyFightInfo [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo `json:"fight_info"`
	Items       []PassItem                            `json:"items"`
}

//! 获得竞技场敌人
type S2C_ArenaSpecialGetEnemy struct {
	Cid   string               `json:"cid"`
	Enemy []*ArenaSpecialEnemy `json:"enemy"`
	Class []int                `json:"class"`
	Dan   []int                `json:"dan"`
}

//! 开始战斗
type S2C_ArenaSpecialStartFight struct {
	Cid        string                                   `json:"cid"`
	RandNum    int64                                    `json:"rand_num"`  // 随机数
	FightInfo  [2][ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo `json:"fightinfo"` // 我方战力结构
	Index      int                                      `json:"index"`     // index
	Result     int                                      `json:"result"`    //结果
	FightID    [ARENA_SPECIAL_TEAM_MAX]int64            `json:"fightid"`   // 战斗id
	BattleInfo [ARENA_SPECIAL_TEAM_MAX]BattleInfo       `json:"battleinfo"`
	Item       []PassItem                               `json:"item"`
	MyClass    int                                      `json:"myclass"`
	MyDan      int                                      `json:"mydan"`
	EnemyClass int                                      `json:"enemyclass"`
	EnemyDan   int                                      `json:"enemydan"`
}

//! 领取奖励
type S2C_ArenaSpecialGetAward struct {
	Cid       string     `json:"cid"`
	Items     []PassItem `json:"items"`
	IsFull    bool       `json:"isfull"`
	Point     int64      `json:"point"`
	ClassTime int64      `json:"classtime"`
}

// 斗技场战报
type S2C_ArenaSpecialGetFights struct {
	Cid       string               `json:"cid"`
	FightInfo []*ArenaSpecialFight `json:"fight_info"` // 战报信息
}

type S2C_ArenaSpecialBuyCount struct {
	Cid      string     `json:"cid"`
	Items    []PassItem `json:"items"`
	BuyCount int        `json:"buycount"`
}

type S2C_ArenaSpecialGetFightInfo struct {
	Cid       string                                `json:"cid"`
	Type      int                                   `json:"type"`
	Uid       int64                                 `json:"uid"`
	FightInfo [ARENA_SPECIAL_TEAM_MAX]*JS_FightInfo `json:"fight_info"`
	Class     int                                   `json:"class"`
	Dan       int                                   `json:"dan"`
	LifeTree  *JS_LifeTreeInfo                      `json:"lifetree"` // 生命树
}

type S2C_ActivityGiftSendInfo struct {
	Cid    string                `json:"cid"`    // cid
	Info   []*ActivityGiftItem   `json:"info"`   //
	Config []*ActivityGiftConfig `json:"config"` //
}

type S2C_ActivityGiftGetAward struct {
	Cid   string     `json:"cid"`
	ID    int        `json:"id"`
	Times int        `json:"time"`
	Items []PassItem `json:"items"`
}

type S2C_GrowthGiftSendInfo struct {
	Cid        string              `json:"cid"` // cid
	IsRecharge bool                `json:"isrecharge"`
	Info       []*GrowthGiftItem   `json:"info"`   //
	Config     []*GrowthGiftConfig `json:"config"` //
}

type S2C_GrowthGiftGetAward struct {
	Cid   string     `json:"cid"`
	ID    int        `json:"id"`
	Items []PassItem `json:"items"`
}

//! 皮肤激活
type S2C_ActivateSkin struct {
	Cid   string     `json:"cid"`   //消息头
	ID    int        `json:"id"`    // 皮肤id
	Items []PassItem `json:"items"` // 消耗物品
	Data  *JS_Skin   `json:"data"`  // 激活的皮肤
}

//! 皮肤设置
type S2C_SetSkin struct {
	Cid       string `json:"cid"`       //消息头
	ID        int    `json:"id"`        // 皮肤id
	HeroIndex int    `json:"heroindex"` // 英雄index
}

//! 皮肤信息
type S2C_SendInfo struct {
	Cid  string     `json:"cid"`  //消息头
	Info []*JS_Skin `json:"info"` //激活的皮肤结构体
}

//! 获得限时抢购信息
type S2C_SpecialPurchaseInfo struct {
	Cid                   string                 `json:"cid"`
	SpecialPurchaseInfo   []*SpecialPurchaseInfo `json:"info"`
	SpecialPurchaseConfig []*ActivityBuyItem     `json:"config"`
}

//! 获得限时抢购奖励
type S2C_SpecialPurchaseGetAward struct {
	Cid      string               `json:"cid"`
	ID       int                  `json:"id"`       //! 类型
	Items    []PassItem           `json:"items"`    // 道具
	Progress *SpecialPurchaseInfo `json:"progress"` // 当前领取数据
}

// 完成状态通知
type S2C_SpecialPurchaseDone struct {
	Cid      string                 `json:"cid"`
	Progress []*SpecialPurchaseInfo `json:"progress"` //! 进度
	Config   []*ActivityBuyItem     `json:"config"`
}

type S2C_LifeTreeUpMainLevel struct {
	Cid   string     `json:"cid"`   // cid
	Items []PassItem `json:"items"` // 道具
	Level int        `json:"level"` // 升级后等级
}

type S2C_LifeTreeUpTypeLevel struct {
	Cid   string     `json:"cid"`   // cid
	Type  int        `json:"type"`  // 类型
	Items []PassItem `json:"items"` // 道具
	Level int        `json:"level"` // 升级后等级
}

type S2C_LifeTreeResetTypeLevel struct {
	Cid   string     `json:"cid"`   // cid
	Items []PassItem `json:"items"` // 道具
	Level int        `json:"level"` // 升级后等级
	Type  int        `json:"type"`  // 类型
}

type S2C_LifeTreeGetAward struct {
	Cid   string     `json:"cid"`   // cid
	Items []PassItem `json:"items"` // 道具
}

type S2C_LifeTreeSendInfo struct {
	Cid       string         `json:"cid"`       // cid
	MainLevel int            `json:"mainlevel"` // 主等级
	Info      []*JS_LifeTree `json:"info"`      // 职业分支
}

type S2C_LifeTreeSendAwardInfo struct {
	Cid   string          `json:"cid"`   // cid
	Award []*JS_HeroAward `json:"award"` // 累计奖励
}

type S2C_LifeTreeSendRedpointInfo struct {
	Cid   string     `json:"cid"`   // cid
	Items []PassItem `json:"items"` // 道具
}

type S2C_InterStellarInfo struct {
	Cid            string          `json:"cid"`
	GalaxyInfo     map[int]*Galaxy `json:"galaxyinfo"`     //! 任务信息
	StellarCount   int             `json:"stellarcount"`   //! 解锁的星星数量
	StellarPos     int             `json:"stellarpos"`     //! 当前位置 万分比
	PrivilegeValue map[int]int     `json:"privilegevalue"` //! 特权值
}

type S2C_UnlockNebula struct {
	Cid         string `json:"cid"`
	NebulaId    int    `json:"nebulaid"`
	NebulaState int    `json:"nebulastate"` //星云解锁状态  0未解锁  1解锁
}

type S2C_GetNebulaWarBox struct {
	Cid            string      `json:"cid"`
	NebulaId       int         `json:"nebulaid"` //星云id
	GroupId        int         `json:"groupid"`  //地图id
	BoxId          int         `json:"boxid"`
	NebulaWar      *NebulaWar  `json:"nebulawar"`
	GetItems       []PassItem  `json:"getitems"`
	PrivilegeValue map[int]int `json:"privilegevalue"` //! 特权值
}

type S2C_ActivityBossInfo struct {
	Cid              string                       `json:"cid"`
	ActivityBossInfo map[int]*JS_ActivityBossInfo `json:"activitybossinfo"`
}

type S2C_ActivityBossFightFail struct {
	Cid string `json:"cid"`
}

type S2C_ActivityBossFightOK struct {
	Cid              string                `json:"cid"`
	RandNum          int64                 `json:"rand_num"`         // 随机数
	FightInfo        [2]*JS_FightInfo      `json:"fightinfo"`        // 我方战力结构
	Score            int64                 `json:"score"`            // 分数
	FightId          int64                 `json:"fightid"`          // 战斗id
	ActivityBossInfo *JS_ActivityBossInfo  `json:"activitybossinfo"` //
	SelfInfo         *ActivityBossSelfInfo `json:"SelfInfo"`         //自己的详情
}

type S2C_ActivityBossTask struct {
	Cid        string       `json:"cid"`
	GetItems   []PassItem   `json:"getitems"`
	TaskInfo   *JS_TaskInfo `json:"taskinfo"`
	ActivityId int          `json:"activityid"`
}

type S2C_ActivityBossExchange struct {
	Cid        string           `json:"cid"`
	GetItems   []PassItem       `json:"getitems"`
	CostItems  []PassItem       `json:"costitems"`
	Exchange   *JS_ExChangeInfo `json:"exchange"`
	ActivityId int              `json:"activityid"`
}

type S2C_ActivityBossStart struct {
	Cid     string `json:"cid"`
	FightId int64  `json:"fightid"`
}

type S2C_ActivityBossStartEx struct {
	Cid             string        `json:"cid"`
	PlayerFightInfo *JS_FightInfo `json:"playerfightinfo"`
	FightInfo       *JS_FightInfo `json:"fightinfo"`
}

type S2C_ActivityBossResetTimes struct {
	Cid              string               `json:"cid"`
	ActivityBossInfo *JS_ActivityBossInfo `json:"activitybossinfo"`
	CostItems        []PassItem           `json:"costitems"`
}

type S2C_ActivityBossGetRank struct {
	Cid              string                `json:"cid"`
	ActivityBossRank ActivityBossRank      `json:"activitybossrank"`
	Id               int                   `json:"id"`
	SelfInfo         *ActivityBossSelfInfo `json:"SelfInfo"` //自己的详情
}

type S2C_ActivityBossGetRecord struct {
	Cid         string               `json:"cid"`
	FightRecord []*ActivityBossFight `json:"fightrecordid"` //战报集
	Target      int64                `json:"target"`
}

//! 获取限时神将信息
type S2C_GetGeneralInfo struct {
	Cid              string             `json:"cid"`           //! cid
	Score            int                `json:"score"`         //! 积分
	LootTimes        int                `json:"loottimes"`     //! 已经抽取次数, 再抽多少次 = 10 - lootTimes
	GeneralAward     []*JS_GeneralAward `json:"generalaward"`  //! 积分奖励领取状态, index = 1 ~ 5 对应40 ~ 2400
	FreeTimes        int                `json:"freetimes"`     //! 剩余免费次数
	HeroIds          []int              `json:"heroids"`       //! 展示的武将Id
	RankConfig       []*TimeGeneralRank `json:"rankconf"`      //! 排行榜奖励
	ShowTime         int64              `json:"showtime"`      //! 活动展示开始时间,这个时候不能抽奖了
	EndTime          int64              `json:"endtime"`       //! 活动结束时间,这个时候活动结束
	RankAward        int                `json:"rankaward"`     //! 0 未领奖 1 已领奖
	ScoreAward       [][]PassItem       `json:"scoreaward"`    //! 积分奖励
	ScorePoints      []int              `json:"scorepoints"`   //! 需要积分
	CostSingleNum    int                `json:"costsinglenum"` //! 单抽价格
	CostTenNum       int                `json:"costtennum"`    //! 10抽价格
	ActRecord        *GeneralAct        `json:"rankstate"`     //! 排行榜领取奖励
	ServerId         int                `json:"serverid"`      //! 服务器Id
	ServerName       string             `json:"servername"`    //! 服务器名称
	ActType          int                `json:"acttype"`       //! 跨服类型, 增加跨服类型
	CallDesc         string             `json:"calldesc"`      //! 活动描述
	NewHero          string             `json:"new_hero"`
	MainHeroLocation []string           `json:"mainhero_location"`
	HeroLocation     []string           `json:"hero_location"`
}

//! 限时神将掉落
type S2C_GeneralLootInfo struct {
	Cid          string     `json:"cid"`          //! cid
	Score        int        `json:"score"`        //! 积分
	LootTimes    int        `json:"loottimes"`    //! 已经抽取次数, 再抽多少次 = 10 - lootTimes
	LootItem     []PassItem `json:"items"`        //! 掉落物品
	Gem          int        `json:"gem"`          //! 钻石
	FreeTimes    int        `json:"freetimes"`    //! 剩余免费次数
	Times        int        `json:"times"`        //! 抽奖次数
	CostItems    []PassItem `json:"costitems"`    //召唤消耗
	GetItemsTran []PassItem `json:"getitemstran"` //分解获得
}

//! 积分奖励领取信息状态
type S2C_GeneralAward struct {
	Cid    string             `json:"cid"`          //! cid
	Items  []PassItem         `json:"items"`        //! 获得物品
	Status []*JS_GeneralAward `json:"generalaward"` //! 积分奖励领取状态, index = 1 ~ 5 对应40 ~ 2400
}

// 获得限时神将排行榜
type S2C_GeneralRank struct {
	Cid           string            `json:"cid"`           //! cid
	Rank          []*Js_GeneralUser `json:"rank"`          //! rank
	MyRank        int               `json:"myrank"`        //! myrank
	GeneralRecord []*GeneralRecord  `json:"generalrecord"` //! 公告
}

// 获得限时神将排行榜奖励
type S2C_GeneralRankAward struct {
	Cid   string     `json:"cid"`   //! cid
	Items []PassItem `json:"items"` //! items
	State int        `json:"state"` //! state = 0 未领取 state = 1已经领取
}

type S2C_HeroGrowInfo struct {
	Cid          string            `json:"cid"`
	HeroGrowInfo []*HeroGrowInfo   `json:"herogrowinfo"`
	Config       []*HeroGrowConfig `json:"config"`
}

type S2C_HeroGrowFreeTask struct {
	Cid          string       `json:"cid"`
	TaskInfo     *JS_TaskInfo `json:"taskinfo"`
	GetItems     []PassItem   `json:"getitems"`
	ActivityType int          `json:"activitytype"`
}

type S2C_CrossArenaInfo struct {
	Cid           string               `json:"cid"`
	Top           []*Js_CrossArenaUser `json:"top"`
	SelfInfo      *Js_CrossArenaUser   `json:"selfinfo"`
	SubsectionMax int                  `json:"subsectionmax"` //最高大段位
	ClassMax      int                  `json:"classmax"`      //最高小段位
	Times         int                  `json:"times"`         //挑战次数
	BuyTimes      int                  `json:"buytimes"`      //购买次数
	StartTime     int64                `json:"starttime"`
	EndTime       int64                `json:"endtime"`
	ShowTime      int64                `json:"showtime"`
	TaskAwardSign map[int]int          `json:"taskawardsign"`
}

type S2C_CrossArenaGetRank struct {
	Cid      string               `json:"cid"`
	Top      []*Js_CrossArenaUser `json:"top"`
	SelfInfo *Js_CrossArenaUser   `json:"selfinfo"`
}

type S2C_CrossArenaAdd struct {
	Cid      string             `json:"cid"`
	SelfInfo *Js_CrossArenaUser `json:"selfinfo"`
}

type S2C_CrossArenaGetDefenceList struct {
	Cid       string               `json:"cid"`
	Info      []*Js_CrossArenaUser `json:"info"`
	FightInfo []*JS_FightInfo      `json:"fightinfo"`
	NextTime  int64                `json:"nexttime"` //下次可刷新时间
}

type S2C_CrossArenaGetReward struct {
	Cid      string      `json:"cid"`
	GetItems []PassItem  `json:"getitems"` //获得奖励
	TaskSign map[int]int `json:"tasksign"` //任务领取标记 key对应表中的id
}

type S2C_CrossArenaArenaAttack struct {
	Cid       string               `json:"cid"`
	FightId   int64                `json:"fightid"`  //这个ID最后会被中心服修正便于存储，避免重复,流程做完了在测
	Attack    *JS_FightInfo        `json:"attack"`   //获得奖励
	Defence   *JS_FightInfo        `json:"defence"`  //任务领取标记 key对应表中的id
	RandNum   int64                `json:"rand_num"` // 随机数
	Info      []*Js_CrossArenaUser `json:"info"`
	FightInfo []*JS_FightInfo      `json:"fightinfo"`
}

type S2C_CrossArenaGetPlayerInfo struct {
	Cid          string             `json:"cid"`
	Info         *Js_CrossArenaUser `json:"info"`
	FightInfo    *JS_FightInfo      `json:"fightinfo"`
	LifeTreeInfo *JS_LifeTreeInfo   `json:"lifetreeinfo"`
}

type S2C_CrossArenaBuyTimes struct {
	Cid      string     `json:"cid"`
	CostItem []PassItem `json:"costitem"` //消耗物品
	BuyTimes int        `json:"buytimes"` //购买次数
}

type S2C_CrossArenaFightOK struct {
	Cid           string               `json:"cid"`
	Top           []*Js_CrossArenaUser `json:"top"`
	SelfInfo      *Js_CrossArenaUser   `json:"selfinfo"`
	OldFightId    int64                `json:"oldfightid"`
	NewFightId    int64                `json:"newfightid"`
	SubsectionMax int                  `json:"subsectionmax"` //最高大段位
	ClassMax      int                  `json:"classmax"`      //最高小段位
	Times         int                  `json:"times"`         //挑战次数
	Result        int                  `json:"result"`        //挑战次数
}

type S2C_CrossArenaUpdate struct {
	Cid      string             `json:"cid"`
	SelfInfo *Js_CrossArenaUser `json:"selfinfo"`
}

type S2C_ActivityBossFestivalInfo struct {
	Cid                      string                       `json:"cid"`
	ActivityBossFestivalInfo *JS_ActivityBossFestivalInfo `json:"info"`
	Level                    int                          `json:"level"` // 用户法阵核心等级
}

type S2C_ActivityBossFestivalGetRecord struct {
	Cid     string      `json:"cid"`
	Records []*PveFight `json:"records"`
}

type S2C_ActivityBossFestivalResult struct {
	Cid         string     `json:"cid"`
	GetItems    []PassItem `json:"getitems"`    //获得奖励
	RewardTimes int        `json:"rewardtimes"` // 是否发奖
}

type S2C_ActivityBossFestivalStart struct {
	Cid           string        `json:"cid"`
	BossFightInfo *JS_FightInfo `json:"fightinfo"` //获得奖励
}

type S2C_LotteryDrawInfo struct {
	Cid             string             `json:"cid"`
	NowStage        int                `json:"nowstage"`        //! 当前阶段
	NowCount        int                `json:"nowcount"`        //! 当前阶段次数
	LotteryDrawInfo []*LotteryDrawItem `json:"lotterydrawinfo"` //!
	LowChoose       int                `json:"lowchoose"`       //! 普通大奖默认选择
	HighChoose      int                `json:"highchoose"`      //! 终极大奖默认选择
	AlreadyGet      map[int]int        `json:"alreadyget"`      //! 大奖领取次数
}

type S2C_DoLotteryDraw struct {
	Cid                   string               `json:"cid"`
	RealTimes             int                  `json:"realtimes"`            //实际抽奖次数
	GetItems              []PassItem           `json:"getitems"`             //获得奖励
	GetItemsDraw          []PassItem           `json:"getitemsdraw"`         //抽奖获得
	CostItems             []PassItem           `json:"costitems"`            //物品消耗
	NowCount              int                  `json:"nowcount"`             //! 当前阶段次数
	LotteryDrawInfo       []*LotteryDrawItem   `json:"lotterydrawinfo"`      //!
	AlreadyGet            map[int]int          `json:"alreadyget"`           //! 大奖领取次数
	LowChoose             int                  `json:"lowchoose"`            //! 普通大奖默认选择
	HighChoose            int                  `json:"highchoose"`           //! 终极大奖默认选择
	HasPrize              int                  `json:"hasprize"`             //! 是否抽中大奖
	LotteryDrawRecordLow  []*LotteryDrawRecord `json:"lotterydrawrecordlow"` //记录，仅在抽中大奖才有
	LotteryDrawRecordHigh []*LotteryDrawRecord `json:"lotterydrawrecordhigh"`
}

type S2C_LotteryDrawChangePrize struct {
	Cid             string             `json:"cid"`
	LotteryDrawInfo []*LotteryDrawItem `json:"lotterydrawinfo"` //!
	LowChoose       int                `json:"lowchoose"`       //! 普通大奖默认选择
	HighChoose      int                `json:"highchoose"`      //! 终极大奖默认选择
}

type S2C_LotteryDrawNext struct {
	Cid             string             `json:"cid"`
	LotteryDrawInfo []*LotteryDrawItem `json:"lotterydrawinfo"` //!
	NowStage        int                `json:"nowstage"`        //! 当前阶段
	NowCount        int                `json:"nowcount"`        //! 当前阶段次数
}

type S2C_LotteryDrawRecord struct {
	Cid                   string               `json:"cid"`
	LotteryDrawRecordLow  []*LotteryDrawRecord `json:"lotterydrawrecordlow"`
	LotteryDrawRecordHigh []*LotteryDrawRecord `json:"lotterydrawrecordhigh"`
}

type S2C_HonourShopInfo struct {
	Cid            string               `json:"cid"`
	HonourGoodInfo []*JS_HonourGoodInfo `json:"honourgoodinfo"`
	NextDayTime    int64                `json:"nextdaytime"`
	NextWeekTime   int64                `json:"nextchecktime"`
}

type S2C_HonourShopBuy struct {
	Cid            string             `json:"cid"`
	GetItems       []PassItem         `json:"getitems"`       //获得
	CostItems      []PassItem         `json:"costitems"`      //消耗
	HonourShopInfo *JS_HonourGoodInfo `json:"honourshopinfo"` //商品同步
}

type S2C_UpdateHonourShop struct {
	Cid            string               `json:"cid"`
	HonourShopInfo []*JS_HonourGoodInfo `json:"honourshopinfo"`
}

type SendSDK_AIWAN struct {
	EventType string `json:"eventtype"`
	AppId     string `json:"appid"`
	ServerId  string `json:"serverid"`
	OpenId    string `json:"openid"`
	RoleId    string `json:"roleid"`
	NickName  string `json:"nickname"`
	RegTime   int    `json:"regtime"`
	PostTime  int    `json:"posttime"`
	Level     int    `json:"level"`
	Ext       string `json:"ext"`
	Sign      string `json:"sign"`
}

//! 消费榜数据
type S2C_ConsumerTopInfo struct {
	Cid           string                  `json:"cid"`
	Level         int                     `json:"level"`
	Point         int                     `json:"point"`
	Rank          int                     `json:"rank"`
	ServerRank    int                     `json:"serverrank"`
	ServerName    string                  `json:"servername"`
	RankAward     int                     `json:"rankaward"`
	Record        []JS_DamageRecord       `json:"record"`
	Award         []int                   `json:"award"`
	MHero         []JS_MagicialHero       `json:"mhero"`
	Msg           []*JS_ConsumerMsg       `json:"msg"`
	BossConfig    []*ConsumetopbossConfig `json:"bossconfig"`
	ListConfig1   []*ConsumetoplistConfig `json:"listconfig1"`
	ListConfig2   []*ConsumetoplistConfig `json:"listconfig2"`
	ShopConfig    []*ConsumetopshopConfig `json:"shopconfig"`
	BossFightInfo *JS_FightInfo           `json:"bossfightinfo"`
}

//! 个人消费榜
type S2C_ConsumerTopUser struct {
	Cid     string                `json:"cid"`
	Ver     int                   `json:"ver"`
	Rank    int                   `json:"rank"`
	TopUser []*JS_ConsumerTopUser `json:"topuser"`
}

//! 服务器消费榜
type S2C_ConsumerTopServer struct {
	Cid    string                  `json:"cid"`
	Ver    int                     `json:"ver"`
	Rank   int                     `json:"rank"`
	SvrId  int                     `json:"serverid"`
	TopSvr []*JS_ConsumerTopServer `json:"topsvr"`
}

//! 攻击
type S2C_ConsumerAttackHero struct {
	Cid        string     `json:"cid"`
	Ret        int        `json:"ret"`
	Kill       int        `json:"kill"`
	Damage     int        `json:"damage"`
	Point      int        `json:"point"`
	Level      int        `json:"level"`
	Hp         int        `json:"hp"`
	MaxHP      int        `json:"maxhp"`
	MHeroLevel int        `json:"mherolevel"`
	Item       []PassItem `json:"item"`
}

type S2C_ConsumerTopKillBoss struct {
	Cid       string `json:"cid"`
	Uname     string `json:"uname"`
	Level     int    `json:"level"`
	Boss      string `json:"boss"`
	HeroId    int    `json:"heroid"`
	Kill      int    `json:"kill"`
	Damage    int    `json:"damage"`
	BossLevel int    `json:"bosslevel"`
	HP        int    `json:"hp"`
	MaxHP     int    `json:"maxhp"`
}

//! 领取奖励
type S2C_ConsumerTopDraw struct {
	Cid   string     `json:"cid"`
	Id    int        `json:"id"`
	Award []int      `json:"award"`
	Item  []PassItem `json:"item"`
}

type S2C_SendRefreshInfo struct {
	Cid  string `json:"cid"`
	Type int    `json:"type"` // 商店类型
}

//! 后宫战斗结束
type S2C_Beauty_LegendOver struct {
	Cid      string     `json:"cid"`
	BeautyId int        `json:"beautyid"`
	Chapter  int        `json:"chapter"`
	Index    int        `json:"index"`
	Finish   bool       `json:"finish"`
	Award    []PassItem `json:"award"`
}

//! 圣物宝物升级
type S2C_Beauty_TreasureaUpLevel struct {
	Cid        string     `json:"cid"`
	Cost       []PassItem `json:"cost"`
	Index      int        `json:"index"`
	CurLevel   int        `json:"curlevel"`
	Fight      int64      `json:"fight"`
	TotalFight int64      `json:"totalfight"`
}

//! 圣物升级
type S2C_Beauty_UpLevel struct {
	Cid        string     `json:"cid"`
	Cost       []PassItem `json:"cost"`
	BeautyId   int        `json:"beautyid"`
	CurLevel   int        `json:"curlevel"`
	Fight      int64      `json:"fight"`
	TotalFight int64      `json:"totalfight"`
}

//! 后宫排行
type S2C_BeautyTop struct {
	Cid string          `json:"cid"`
	Top []*JS_BeautyTop `json:"top"`
	Ver int             `json:"ver"`
}

//! 后宫消息
type S2C_BeautyInfo struct {
	Cid  string         `json:"cid"`
	Info Son_BeautyInfo `json:"info"`
}

//! 子消息
type Son_BeautyInfo struct {
	Beautyinfo  []BeautyInfo `json:"beautyinfo"`
	Fight       int64        `json:"count"`
	Lastupdtime int64        `json:"lastupdtime"`
	Uid         int64        `json:"uid"`
}

//
type S2C_CheckFight struct {
	Cid  string `json:"cid"`
	Code int    `json:"code"` //0正常
}

//! 坐骑数据-同步结构
type S2C_HorseInfo struct {
	Cid          string          `json:"cid"`          //! cid
	Level        int             `json:"level"`        //! 等级
	Exp          int             `json:"exp"`          //! 经验
	SummonNormal int             `json:"summonnormal"` //! 召唤次数
	SummonSenior int             `json:"summonsenior"` //! 高级召唤次数
	SummonTime   int64           `json:"summontime"`   //! 上次获得时间
	Discern      int             `json:"discern"`      //! 鉴定次数
	Combine      int             `json:"combine"`      //! 合成次数
	Summonlst    [2]int          `json:"summonlst"`    //! 召唤野马列表
	Horselst     []*JS_HorseInfo `json:"horselst"`     //! 当前坐骑列表
	HorseTask    []*HorseTask    `json:"horsetask"`    //! 高级召唤任务
}

type S2C_SummonHorse struct {
	Cid          string       `json:"cid"`
	Level        int          `json:"level"`
	Exp          int          `json:"exp"`
	Gold         int          `json:"gold"` //! 金币
	Gem          int          `json:"gem"`  //! 钻石
	Summontime   int64        `json:"summontime"`
	SummonNormal int          `json:"summonnormal"`
	SummonSenior int          `json:"summonsenior"`
	Summonlst    []int        `json:"summonlst"`
	Item         []PassItem   `json:"item"`
	HorseTask    []*HorseTask `json:"horsetask"` //! 高级召唤任务
}

type S2C_DiscernHorse struct {
	Cid        string          `json:"cid"`
	Level      int             `json:"level"`
	Exp        int             `json:"exp"`
	Gold       int             `json:"gold"`       //! 金币
	Gem        int             `json:"gem"`        //! 钻石
	Discern    int             `json:"discern"`    //! 鉴定次数
	Summonlst  [2]int          `json:"summonlst"`  //!
	SummonTime int64           `json:"summontime"` //! 召唤刷新时间
	Item       []PassItem      `json:"item"`
	Horselst   []*JS_HorseInfo `json:"horselst"` //! 当前坐骑列表
}

//! 马魂-同步结构
type S2C_HorseSoulInfo struct {
	Cid          string              `json:"cid"`
	Soullst      []*JS_HorseSoulInfo `json:"soullst"`
	Horselst     []*JS_HorseInfo     `json:"horselst"`
	Summontime   int64               `json:"summontime"`
	SummonNormal int                 `json:"summonnormal"`
	SummonSenior int                 `json:"summonsenior"`
	NewSoul      bool                `json:"newsoul"`
}

//! 分解马魂
type S2C_DecomposeHorse struct {
	Cid      string             `json:"cid"`
	Horselst []int              `json:"horselst"`
	Soullst  []JS_HorseSoulInfo `json:"soullst"`
	Item     []PassItem         `json:"item"`
}

//! 分解马魂
type S2C_DecomposeSoul struct {
	Cid      string             `json:"cid"`
	CostSoul []JS_HorseSoulInfo `json:"costsoul"`
	Item     []PassItem         `json:"item"`
}

//! 合成橙色坐骑
type S2C_CombineHorse struct {
	Cid     string       `json:"cid"`
	Combine int          `json:"combine"`
	Cost    []PassItem   `json:"item"`
	Horse   JS_HorseInfo `json:"horse"`
}

//! 合成攻击橙马
type S2C_UpHorse struct {
	Cid      string             `json:"cid"`
	Cost     []PassItem         `json:"item"`
	Soullst  []JS_HorseSoulInfo `json:"soullst"`
	Material []int              `json:"material"`
	Horse    JS_HorseInfo       `json:"horse"`
}

//! 使用坐骑
type S2C_MountHero struct {
	Cid       string             `json:"cid"`
	Heroid    int                `json:"heroid"`
	Horseid   int                `json:"horseid"`
	HorseSoul []JS_HorseSoulInfo `json:"horsesoul"`
	Oldhorse  int                `json:"oldhorse"`
	Soullst   []JS_HorseSoulInfo `json:"soullst"`
	Horse     JS_HorseInfo       `json:"horse"`
	TeamAttr  *TeamAttr          `json:"team_attr,omitempty"` // 阵位属性
	TeamType  int                `json:"team_type"`           // 阵营类型
	Index     int                `json:"index"`               // 上阵序号
}

//! 镶嵌马魂
type S2C_EmbedHorseSoul struct {
	Cid     string             `json:"cid"`
	Horseid int                `json:"horseid"`
	Soulid  int                `json:"soulid"`
	Soullst []JS_HorseSoulInfo `json:"soullst"`
	Horse   JS_HorseInfo       `json:"horse"`
}

//! 升级马魂
type S2C_UpHorseSoul struct {
	Cid     string       `json:"cid"`
	Horseid int          `json:"horseid"`
	Index   int          `json:"index"`
	Soulid  int          `json:"soulid"`
	Item    []PassItem   `json:"item"`
	Horse   JS_HorseInfo `json:"horse"`
}

//! 兑换马魂
type S2C_ExchangeSoul struct {
	Cid   string             `json:"cid"`
	Index int                `json:"index"`
	Soul  []JS_HorseSoulInfo `json:"soul"`
	Item  []PassItem         `json:"item"`
}

//! 高级召唤任务奖励
type S2C_AwardHorseTask struct {
	Cid   string     `json:"cid"`
	Info  *HorseTask `json:"info"`  //! 任务信息
	Items []PassItem `json:"items"` //! 道具信息
}

//! 转换坐骑
type S2C_SwitchHero struct {
	Cid        string        `json:"cid"`
	Horseid    int           `json:"horseid"`
	OldHorseId int           `json:"oldhorseid"`
	Horse      *JS_HorseInfo `json:"horse"`
	Items      []PassItem    `json:"items"` //! 道具信息
}

//! 洗练战马
type S2C_WashHero struct {
	Cid   string        `json:"cid"`
	Horse *JS_HorseInfo `json:"horse"`
	Items []PassItem    `json:"items"` //! 道具信息
}

//! 洗练保留
type S2C_SaveWash struct {
	Cid   string        `json:"cid"`
	Horse *JS_HorseInfo `json:"horse"`
}

//! 转换保留
type S2C_SaveSwitch struct {
	Cid   string        `json:"cid"`
	Horse *JS_HorseInfo `json:"horse"`
}

type S2C_CrossArena3V3Info struct {
	Cid           string                  `json:"cid"`
	Top           []*Js_CrossArena3V3User `json:"top"`
	SelfInfo      *Js_CrossArena3V3User   `json:"selfinfo"`
	SubsectionMax int                     `json:"subsectionmax"` //最高大段位
	ClassMax      int                     `json:"classmax"`      //最高小段位
	Times         int                     `json:"times"`         //挑战次数
	BuyTimes      int                     `json:"buytimes"`      //购买次数
	StartTime     int64                   `json:"starttime"`
	EndTime       int64                   `json:"endtime"`
	ShowTime      int64                   `json:"showtime"`
	TaskAwardSign map[int]int             `json:"taskawardsign"`
}

type S2C_CrossArena3V3GetRank struct {
	Cid      string                  `json:"cid"`
	Top      []*Js_CrossArena3V3User `json:"top"`
	SelfInfo *Js_CrossArena3V3User   `json:"selfinfo"`
}

type S2C_CrossArena3V3Add struct {
	Cid      string                `json:"cid"`
	SelfInfo *Js_CrossArena3V3User `json:"selfinfo"`
}

type S2C_CrossArena3V3GetDefenceList struct {
	Cid       string                  `json:"cid"`
	Info      []*Js_CrossArena3V3User `json:"info"`
	FightInfo [][]*JS_FightInfo       `json:"fightinfo"`
	NextTime  int64                   `json:"nexttime"` //下次可刷新时间
}

type S2C_CrossArena3V3GetReward struct {
	Cid      string      `json:"cid"`
	GetItems []PassItem  `json:"getitems"` //获得奖励
	TaskSign map[int]int `json:"tasksign"` //任务领取标记 key对应表中的id
}

type S2C_CrossArena3V3ArenaAttack struct {
	Cid     string                                `json:"cid"`
	FightId [CROSSARENA3V3_TEAM_MAX]int64         `json:"fightid"`  //这个ID最后会被中心服修正便于存储，避免重复,流程做完了在测
	Attack  [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo `json:"attack"`   //获得奖励
	Defence [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo `json:"defence"`  //任务领取标记 key对应表中的id
	RandNum int64                                 `json:"rand_num"` // 随机数
}

type S2C_CrossArena3V3GetPlayerInfo struct {
	Cid          string                                `json:"cid"`
	Info         *Js_CrossArena3V3User                 `json:"info"`
	FightInfo    [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo `json:"fightinfo"`
	LifeTreeInfo *JS_LifeTreeInfo                      `json:"lifetreeinfo"`
}

type S2C_CrossArena3V3BuyTimes struct {
	Cid      string     `json:"cid"`
	CostItem []PassItem `json:"costitem"` //消耗物品
	BuyTimes int        `json:"buytimes"` //购买次数
}

type S2C_CrossArena3V3FightOK struct {
	Cid           string                        `json:"cid"`
	Top           []*Js_CrossArena3V3User       `json:"top"`
	SelfInfo      *Js_CrossArena3V3User         `json:"selfinfo"`
	OldFightId    [CROSSARENA3V3_TEAM_MAX]int64 `json:"oldfightid"`
	NewFightId    [CROSSARENA3V3_TEAM_MAX]int64 `json:"newfightid"`
	SubsectionMax int                           `json:"subsectionmax"` //最高大段位
	ClassMax      int                           `json:"classmax"`      //最高小段位
	Times         int                           `json:"times"`         //挑战次数
	Result        int                           `json:"result"`        //挑战次数
}

type S2C_CrossArena3V3Update struct {
	Cid      string                `json:"cid"`
	SelfInfo *Js_CrossArena3V3User `json:"selfinfo"`
}

type S2C_GetRankRewardRank struct {
	Cid      string               `json:"cid"`
	Id       int                  `json:"id"`
	Top      []*RankRewardTopNode `json:"top"`
	SelfInfo *RankRewardTopNode   `json:"selfinfo"`
	Config   []*RankRewardConfig  `json:"config"`
}

type S2C_GetRankRewardReward struct {
	Cid      string                `json:"cid"`
	Id       int                   `json:"id"`
	GetItems []PassItem            `json:"getitems"`
	SelfInfo *RankRewardPlayerInfo `json:"selfinfo"`
}
