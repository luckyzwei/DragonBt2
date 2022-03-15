package game

// 任务需要经常判断，故单独写一个结构
type TaskNode struct {
	Id        int
	Tasktypes int
	N1        int
	N2        int
	N3        int
	N4        int
}

const (
	LOGIC_FALSE = 0
	LOGIC_TRUE  = 1
	LOGIC_ALL   = 2
)

const (
	TIME_RESET_TYPE_CREATE_ROLE = 1 //根据创建角色天数
	TIME_RESET_TYPE_TIME        = 3 //根据固定时间计算
	TIME_RESET_TYPE_TIME_BASE   = 4 //基准时间
)

//tariff:type
const (
	TARIFF_TYPE_HERO_BACK      = 6  //英雄回退
	TARIFF_TYPE_ONHOOK_FAST    = 7  //快速挂机
	TARIFF_TYPE_FIND_OPEN_CAMP = 17 //阵营池开启

	TARIFF_RESONANCE_CRYSTAL_RESONANCE_GEM = 13
	TARIFF_RESONANCE_CRYSTAL_RESONANCE     = 14
	TARIFF_RESONANCE_CRYSTAL_CLEAN_CD      = 18
	TARIFF_ASTROLOGY_GEN                   = 26
	TARIFF_ASTROLOGY_ITEM                  = 27

	TARIFF_TYPE_ARENA_NORMAL     = 1 //普通竞技场
	TARIFF_TYPE_ARENA_SPECIAL    = 2 //高阶竞技场
	TARIFF_TYPE_ARENA_NORMALBUY  = 3 //普通竞技场
	TARIFF_TYPE_ARENA_SPECIALBUY = 4 //高阶竞技场

	TARIFF_TYPE_VOID_HERO_CANCEL = 23 //虚空英雄取消

	TARIFF_TYPE_LIFE_TREE_RESET       = 25 //生命树重置
	TARIFF_TYPE_ACTIVITY_BOSS_RESET   = 34 //世界BOSS次数重置
	TARIFF_TYPE_CROSS_ARENA_TIMES     = 35 //跨服竞技场刷新次数
	TARIFF_TYPE_LOTTERY_DRAW_COST     = 36 //抽奖
	TARIFF_TYPE_CROSS_ARENA_3V3_TIMES = 37 //跨服竞技场刷新次数
)

const (
	SIMPLE_NUM_SELFFIND_OFFSET = 11
	CROSSARENA_FREETIMES       = 12
	ASTROLOGY_INIT_SCORE       = 13
	LOTTERY_DRAW_ADD           = 14
	LOTTERY_DRAW_INIT          = 15
	CROSSARENA_3V3_FREETIMES   = 16
)

const (
	CANTFINISH = 0
	CANTAKE    = 1
	TAKEN      = 2
)

//新魔龙 英雄战力更新的原因
const (
	ReasonPlayerLogin = 1 //玩加登录
	//英雄本身相关
	ReasonHeroNew              = 101 //英雄创建
	ReasonHeroCheck            = 102 //英雄检查
	ReasonHeroReborn           = 103 //英雄重生
	ReasonHeroBack             = 104 //英雄回退
	ReasonHeroLvUp             = 105 //英雄升级
	ReasonHeroStarUp           = 106 //英雄升阶
	ReasonHeroFate             = 107 //英雄羁绊
	ReasonHeroLifeTree         = 108 //英雄生命树
	ReasonHeroVoidSetResonance = 109 //虚空英雄
	ReasonHoly                 = 110 //圣物

	//装备
	ReasonEquipWear = 201 //穿装备
	ReasonEquipOff  = 202 //脱装备
	ReasonEquipLvUp = 203 //装备升级
	//神器
	ReasonArtifactWear = 211 //穿神器
	ReasonArtifactOff  = 212 //脱神器
	//专属
	ReasonExclusiveUnlock = 221 //解锁专属
	ReasonExclusiveLvUp   = 222 //专属升级

	ReasonStageTalentSetSkill = 231 //设置技能

	ReasonEmbedHorseSoul  = 251 // 镶嵌魔魂
	ReasonUpHorseSoul     = 252 // 升级魔魂
	ReasonMountHorse      = 253 // 上魔宠
	ReasonUnMountHorse    = 254 // 下魔宠
	ReasonAwakenHorse     = 255 // 觉醒魔宠
	ReasonSaveWashHorse   = 256 // 魔宠洗练保存
	ReasonSaveSwicthHorse = 257 // 魔宠转换保存
	ReasonRemoveHorseSoul = 258 // 魔宠删除马魂
)
const (
	SPECIAL_STOP = 1
	INIT_LV      = 0
)

const (
	ITEM_SOUL_HERO_STONE     = 13000002 //紫色魂石
	ITEM_NEW_PIT_REBORN      = 40000011 //女神之泪
	ITEM_FIND_GEM_ITEM       = 42010001 //钻石召唤物品
	ITEM_FIND_CAMP_ITEM      = 42010002 //阵营召唤卷
	ITEM_FIND_ASTROLOGY_ITEM = 42010003 //占星卷
	ITEM_FIND_GENERAL_ITEM   = 42010004 //次元召唤卷
	ITEM_EQUIP_LVUP_ITEM_LOW = 81000001 //黑铁铸币
	ITEM_GOLD                = 91000001
	ITEM_GEM                 = 91000002
	ITEM_PLAYER_EXP          = 91000005 //领主经验
	ITEM_UNION               = 91000025 //贡献币（公会币）
	ITEM_SCIENCE             = 91000033
	ITEM_HERO_EXP            = 91000041 //英雄经验
	ITEM_POWDER              = 91000044 //魔粉
	ITEM_BACK_COIN           = 91000045 //英灵币（遣散币）
	ITEM_FRIEND_POINT        = 91000048 //友情点
	ITEM_ARENA_SPECIAL_COIN  = 91000050 //高阶竞技场硬币
	ITEM_NEW_PIT_COIN        = 91000051 //地牢币
	ITEM_CAMP_HERO_RAND      = 91000052 //种族英雄选择卡
	ITEM_ACCESS_ITEM         = 91000058 //收藏家点数
	ITEM_MONTH_SCORE         = 91000059 //月卡积分
)

const (
	ATTR_TYPE_HP      = 1
	ATTR_TYPE_ATTACK  = 2
	ATTR_TYPE_DEFENSE = 3
	ATTR_TYPE_FIGHT   = 99
	PER_BIT           = 10000 //万分比
)

type CsvNode map[string]string

type Item struct {
	ItemId  int
	ItemNum int
}

type Cond struct {
	Type int
	Cond int
}

// 派遣任务配置
type DispatchConfig struct {
	Id        int   `json:"id"`        // 派遣任务Id
	MinLv     int   `json:"minlv"`     // 最低等级
	MaxLv     int   `json:"maxlv"`     // 最高等级
	Weight    int   `json:"weight"`    // 权重
	Time      int   `json:"time"`      // 时间秒
	Commander int   `json:"commander"` // 上几个武将
	CondTypes []int `json:"type"`      // 条件类型
	Conds     []int `json:"condition"` // 条件
	Items     []int `json:"item"`      // 物品奖励
	Nums      []int `json:"num"`       // 物品奖励
	Once      int   `json:"once"`      // 最多出现次数
	Class     int   `json:"class"`     // 星级
}

// 排行榜奖励
type TimeGeneralRank struct {
	Id          int        `json:"Id"`
	Group       int        `json:"TaskGroup"`   // 分组Id
	RankMin     int        `json:"RankMin"`     // 最小排名
	RankMax     int        `json:"RankMax"`     // 最大排名
	NormalAward []PassItem `json:"NormalAward"` // 普通奖励
	NeetPoint   int        `json:"NeetPoint"`   // 需要积分
	ExtraAward  []PassItem `json:"ExtraAward"`  // 额外积分
}

type GeneralRankMail struct {
	MailTitle string
	MailText1 string
	MailText2 string
	MailText3 string
}

type RedpacketmoneyConfig struct {
	Id     int `json:"id"`
	Money  int `json:"money"`
	Item   int `json:"item"`
	Num    int `json:"num"`
	People int `json:"people"`
}

type RedpacketConfig struct {
	Id             int   `json:"id"`
	Item           int   `json:"item"`
	Prices         []int `json:"price"`
	Costs          []int `json:"cost"`
	Vips           []int `json:"vip"`
	Peoples        []int `json:"people"`
	Ceiling        int   `json:"Ceiling"`
	Ceilingtime    int   `json:"Ceilingtime"`
	Upperlimit     int   `json:"Upperlimit"`
	Upperlimittime int   `json:"Upperlimittime"`
}

type ExpspeedupConfig struct {
	Leveldf1   int `json:"leveldf1" trim:"0"`
	Leveldf2   int `json:"leveldf2" trim:"0"`
	Speedup    int `json:"speedup"`
	Ocpdecline int `json:"ocpdecline"`
}

type TreasureSuitAttributeConfig struct {
	Group           int
	Suitnum         int
	Attributetypes  []int
	Attributevalues []int64
}

type TreasureAwakenConfig struct {
	Group           int
	Star            int
	Costitems       []int
	Costnums        []int
	Attributetypes  []int
	Attributevalues []int64
	Skill           int
	Resetitem       int
	Resetnum        int
	Returnitems     []int
	Returnnums      []int
	Quality         int
}

type TreasureClearAttributeConfig struct {
	Attributelistgroup int
	Weightattribute    int
	Attributetype      int
	Lv                 int
	Oddslvs            []int
	Attributevalve     int64
	Quality            int
}

type TreasureClearItemConfig struct {
	Itemid     int
	Effecttype int
	Parms      []int
}

type TreasureEquipConfig struct {
	Id              int
	Class           int
	Position        int
	Attributetypes  []int
	Attributevalues []int64
	Attributeholes  []int
	Addeffect       int
	Decomposeitem   int
	Decomposenum    int
	Quality         int
	DecomposeGroup  int
}

type TreasureHeroConfig struct {
	Heroid      int
	Class       int
	Attributes  []int
	Suitgroup   int
	Awakengroup int
}

type LuckyturntableConfig struct {
	N4          int   `json:"n"`
	Id          int   `json:"id"`
	Type        int   `json:"type"`
	Show        int   `json:"show"`
	Sort        int   `json:"sort"`
	Item        int   `json:"item"`
	Num         int   `json:"num"`
	Value       int   `json:"value"`
	Luckadd     int   `json:"luckadd"`
	Lucktargets []int `json:"lucktarget"`
	Luckvalues  []int `json:"luckvalue"`
}

type LuckyturntablelistConfig struct {
	N4    int   `json:"n"`
	Id    int   `json:"id"`
	Arget int   `json:"arget"`
	Items []int `json:"item"`
	Nums  []int `json:"num"`
}

type SmeltConfig struct {
	Level     int `json:"level"`
	Exp       int `json:"exp"`
	Freetimes int `json:"free_times"`
	Maxtime   int `json:"max_time"`
	Addtime   int `json:"add_time"`
	Getexp    int `json:"get_exp"`
	Getitem   int `json:"Get_item"`
	Basevalue int `json:"base_value"`
	Luckadd   int `json:"luck_add"`
	Mcrate    int `json:"mcrate"`
	Odds      int `json:"odds"`
	Dropgroup int `json:"drop_group"`
	Mcrit1    int `json:"mcrit1" trim:"0"`
	Mcrit2    int `json:"mcrit2" trim:"0"`
	Mcrit3    int `json:"mcrit3" trim:"0"`
}

type SmeltDropConfig struct {
	Id           int `json:"id"`
	Group        int `json:"group"`
	Weight       int `json:"weight"`
	Index        int `json:"index"`
	Item         int `json:"item"`
	Itemnum      int `json:"item_num"`
	Costitem     int `json:"cost_item"`
	Primecost    int `json:"prime_cost"`
	Discount     int `json:"discount"`
	Presentprice int `json:"present_price"`
}

type SmeltPurchaseConfig struct {
	Id       int `json:"id"`
	Minnum   int `json:"min_num"`
	Maxnum   int `json:"max_num"`
	Costitem int `json:"cost_item"`
	Costnum  int `json:"cost_num"`
}

type InvitingConfig struct {
	Group    int   `json:"group"`
	Id       int   `json:"id"`
	Sort     int   `json:"sort"`
	Needvip  int   `json:"needvip"`
	Costitem int   `json:"costitem"`
	Costnum  int   `json:"costnum"`
	Items    []int `json:"item"`
	Minnums  []int `json:"minnum"`
	Maxnums  []int `json:"maxnum"`
	Values   []int `json:"value"`
	Start    int   `json:"start"`
	End      int   `json:"end"`
}

type TigerAdvancedConfig struct {
	Group     int
	Costs     []int
	Nums      []int
	Costitems []int
	Costnums  []int
	Icon      int
}

type TigerAttributeConfig struct {
	Id             int     `json:"id"`
	Stage          int     `json:"stage"`
	Attributetypes []int   `json:"attribute_type"`
	Values         []int64 `json:"value"`
	Quality        int     `json:"quality"`
}

type TigerStuntConfig struct {
	Stuntid      int
	Upgrademodel int
	HufuLimits   int // 升级特技需要的虎符Id
}

type TigerStuntUpgradeConfig struct {
	Stuntgroup     int
	LimitUpgradeId int
	Level          int
	Attributetype  []int
	Value          []int64
	Costitem       int
	Costnum        int
	Resetitem      int
	Resetnum       int
	Returnitem     int
	Returnnum      int
}

type TigerSymbolConfig struct {
	ID    int
	Holes []int
}

type TigerUpgradeConfig struct {
	Id            int     `json:"id"`
	Lvlimits      int     `json:"lv_limits"`
	Stage         int     `json:"stage"`
	Level         int     `json:"level"`
	Advancedgroup int     `json:"advanced_group"`
	Costitems     []int   `json:"cost_item"`
	Costnums      []int   `json:"cost_num"`
	Slottypes     []int   `json:"slot_type"`
	Slotvalues    []int64 `json:"slot_value"`
	Fights        []int64 `json:"fight"`
}

type LuckdrawConfig struct {
	N4     int   `json:"n"`
	Id     int   `json:"id"`
	Values []int `json:"value"`
	Costs  []int `json:"cost"`
}

type LuckdrawgroupConfig struct {
	N4    int `json:"n"`
	Group int `json:"group"`
	Id    int `json:"id"`
	Item  int `json:"item"`
	Num   int `json:"num"`
}

type LuckdrawlistConfig struct {
	N4    int   `json:"n"`
	Id    int   `json:"id"`
	Arget int   `json:"arget"`
	Items []int `json:"item"`
	Nums  []int `json:"num"`
}

type TreasureSuitConfig struct {
	Heroid      int
	Class       int
	Attributes  []int
	Suitgroup   int
	Awakengroup int
}

type TreasureDecomposeConfig struct {
	Group      int
	Totallevel int
	Itemid     int
	Value      int
}

type LuckeggConfig struct {
	N4        int   `json:"n"`
	Id        int   `json:"id"`
	Values    []int `json:"value"`
	Costs     []int `json:"cost"`
	Costitems []int `json:"costitem"`
}

type LuckegggroupConfig struct {
	N4    int `json:"n"`
	Group int `json:"group"`
	Id    int `json:"id"`
	Item  int `json:"item"`
	Num   int `json:"num"`
	End   int `json:"end"`
}

type LuckstartConfig struct {
	N4        int   `json:"n"`
	Id        int   `json:"id"`
	Step      int   `json:"step"`
	Tasktypes int   `json:"tasktypes"`
	Ps        []int `json:"p"`
	Items     []int `json:"item"`
	Nums      []int `json:"num"`
}

type HalfMoonTrial struct {
	TaskId int // 任务Id
	Index  int // index
	Hard   int // 难度
	Sort   int // 第几天
}

type HalfMoonBeauty struct {
	TaskId int // 任务Id
	MaxLv  int // 宝物最大等级
	Sort   int // 第几天
	N2     int // 美人等级
}

type CitygvgConfig struct {
	Citysize     int    `json:"citysize"`
	Items        []int  `json:"item"`
	Nums         []int  `json:"num"`
	Conditions   []int  `json:"condition"`
	Position     int    `json:"position"`
	Name         string `json:"name"`
	Number       int    `json:"number"`
	Addition     int    `json:"addition"`
	Reduce       int    `json:"reduce"`
	Limits       []int  `json:"limit"`
	Limitspecial int    `json:"limitspecial"`
	Repress      int    `json:"repress"`
}

type CityrandomConfig struct {
	Id     int `json:"id"`
	Type   int `json:"type"`
	Limit  int `json:"limit"`
	Weight int `json:"weight"`
	Item   int `json:"item"`
	Num    int `json:"num"`
	Time   int `json:"time"`
	Min    int `json:"min"`
	Max    int `json:"max"`
}

type CityrepressConfig struct {
	Id    int    `json:"id"`
	Class int    `json:"class"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Type  int    `json:"type"`
	Drop  int    `json:"drop"`
	Num   int    `json:"num"`
	Time  int    `json:"time"`
	Txt   string `json:"txt"`
}

type DailyrechargeConfig struct {
	Id        int    `json:"id"`
	Group     int    `json:"group"`
	Type      int    `json:"type"`
	Index     int    `json:"index"`
	ShowValue int    `json:"show_value"`
	Time      int    `json:"time"`
	Diamond   int    `json:"numeral"`
	Items     []int  `json:"item"`
	Nums      []int  `json:"num"`
	AddText   string `json:"add_txt"`
	Efficacy  []int  `json:"efficacy"`
}

type GemsweeperConfig struct {
	Id              int
	Type            int
	Quality         int
	Gradelevel      int
	Activationitems []int
	Activationnums  []int
	Costitem        int
	Costnum         int
	Eventgroup      int
	Duration        int
	Itemshows       []int
}

type BecomeStronger struct {
	Id            int    `json:"ID"`
	ConditionType int    `json:"condition"`
	Condition     string `json:"systempicture"`
	Boxid         int    `json:"boxid"`
}

type GemsweepereventConfig struct {
	Group           int
	Cycle           int
	Step            int
	Next            int
	Level           int
	Event           int
	Items           []int
	Nums            []int
	Itembagnummin   int
	Itembagnummax   int
	Itembaggroup    int
	Itemshows       []int
	Completionitems []int
	Completionnums  []int
}
type CsvMgr struct {
	Data  map[string]map[int]CsvNode
	Data2 map[string][]CsvNode

	// 特殊结构
	Activity_Type  []int
	State_City_Num map[int][]int
	Herodebris_CSV map[int]int
	Chesttotal_CSV map[int][]CsvNode
	Chesttotal_SUM map[int]int
	Equchest_CSV   map[int][]CsvNode
	Equchest_SUM   map[int]int

	Soulchest_CSV map[int][]CsvNode
	Soulchest_SUM map[int]int
	Soulchest_Max int

	WealTask_CSV      map[int]*TaskNode
	WarTaskTask_CSV   map[int]*TaskNode
	SevenDayTask_CSV  map[int]*TaskNode
	SevenStatus       map[int]int
	HalfMoonTask_CSV  map[int]*TaskNode
	HalfMonnStatus    map[int]int
	WarTargetTask_CSV map[int]*TaskNode

	War_CSV       map[int]CsvNode
	WarTarget_CSV map[int]CsvNode

	Name1_CSV []string
	Name2_CSV []string
	Name3_CSV []string

	HeroSpirit_CSV  map[int]CsvNode
	Spy_CSV         map[int]CsvNode
	TrialExp_CSV    map[int]CsvNode
	EquAwake_CSV    map[int]CsvNode
	OfficeRobot_CSV map[int][]CsvNode
	Gestapo_CSV     map[int]CsvNode
	Exciting_CSV    map[int]CsvNode

	Conditions_CSV            map[int]CsvNode
	Diplomacy_CSV             map[int]CsvNode
	Honorkill_CSV             map[int]CsvNode
	Homeoffice_ranking_CSV    map[int]CsvNode
	Reward_CSV                map[int]CsvNode
	Gemsweeper_extradrop_CSV  map[int]CsvNode
	Gemsweeper_itembag_GG_CSV map[int][]CsvNode
	Gemsweeper_itembag_SUM    map[int]int
	Peoplecity_CSV            map[int][]CsvNode
	Peoplebox_CSV             map[int][]CsvNode

	Harem_TreasureaAdv_CSV map[int]CsvNode

	WarCity_CSV             [3]map[int]CsvNode
	WarCity_Map             [3]map[int]*WorldMapConfig
	PromoteBox_CSV          map[int]CsvNode
	Gemsweeper_event_CSV    map[int]CsvNode
	Money_CSV               map[int]CsvNode
	ExpeditionBuffGroup_CSV map[int]CsvNode
	SevenDay_CSV            map[int]CsvNode
	HalfMoon_CSV            map[int]CsvNode
	WarList_CSV             map[int]CsvNode
	//ActivityBox_CSV         map[int][]CsvNode

	// 无双神将
	ConsumerList_CSV map[int][]CsvNode
	ConsumerShop_CSV map[int][]CsvNode

	HorseParam_CSV map[string]CsvNode
	HorseAttr_CSV  map[int][]CsvNode
	DispatchConfig map[int]*DispatchConfig

	TimeGeneralRankLst  []*TimeGeneralRank
	GeneralRankMail     *GeneralRankMail
	HalfMoonLogin       map[int]int             // login day, taskId
	HalfMoonTrial       map[int]*HalfMoonTrial  // trial taskId
	HalfMoonBeauty      map[int]*HalfMoonBeauty // beauty taskId
	HalfMoonTrialSort   int                     // 半月试炼第几天
	SmeltConfig         map[int]*SmeltConfig
	SmeltDropConfig     []*SmeltDropConfig
	SmeltPurchaseConfig map[int]*SmeltPurchaseConfig
	InvitingConfig      []*InvitingConfig

	TigerAttributeConfig         []*TigerAttributeConfig
	TigerStuntConfig             []*TigerStuntConfig
	TigerStuntUpgradeConfig      []*TigerStuntUpgradeConfig
	TigerSymbolConfig            []*TigerSymbolConfig
	TigerUpgradeConfig           []*TigerUpgradeConfig
	TigerAdvancedConfig          []*TigerAdvancedConfig
	RedpacketConfig              []*RedpacketConfig
	RedpacketmoneyConfig         []*RedpacketmoneyConfig
	ExpspeedupConfig             []*ExpspeedupConfig
	TreasureAwakenConfig         []*TreasureAwakenConfig
	TreasureClearAttributeConfig []*TreasureClearAttributeConfig
	MaxTreasureHoleLv            map[int]int
	TreasureClearItemConfig      map[int]*TreasureClearItemConfig
	TreasureEquipConfig          map[int]*TreasureEquipConfig
	TreasureHeroConfig           map[int]*TreasureHeroConfig
	TreasureSuitAttributeConfig  []*TreasureSuitAttributeConfig
	MaxTreasureAwaken            map[int]int
	LuckyturntableConfig         []*LuckyturntableConfig
	LuckyturntablelistConfig     []*LuckyturntablelistConfig
	LuckdrawConfig               []*LuckdrawConfig
	LuckdrawgroupConfig          []*LuckdrawgroupConfig
	LuckdrawlistConfig           []*LuckdrawlistConfig
	TreasureDecomposeConfig      []*TreasureDecomposeConfig
	TreasureDecomposeMap         map[int]map[int]*TreasureDecomposeConfig // grop,level
	JJCRobotConfig               []*JJCRobotConfig
	RobotMap                     map[int]*JJCRobotConfig

	LevelMonsterMap map[int]*LevelMonsterConfig

	WorldLvLevelConfig  []*WorldLvLevelConfig
	WorldLvTpyeConfig   []*WorldLvTpyeConfig
	WorldShowBossConfig []*WorldShowBossConfig

	HonourShopConfigMap map[int]*HonourShopConfig

	CitygvgConfig     []*CitygvgConfig
	CityrandomConfig  []*CityrandomConfig
	CityrepressConfig map[int]*CityrepressConfig

	LuckeggConfig      []*LuckeggConfig
	LuckegggroupConfig []*LuckegggroupConfig
	LuckstartConfig    []*LuckstartConfig
	LuckStartMap       map[int]map[int]*TaskNode
	LuckStartConfigMap map[int]map[int]*LuckstartConfig

	//! 连续充值
	DailyrechargeConfig []*DailyrechargeConfig
	DailyrechargeMap    map[int][]*DailyrechargeConfig

	GemsweeperConfig      map[int]*GemsweeperConfig
	GemsweepereventConfig []*GemsweepereventConfig
	GemGroupCycle         map[int]int //group:cycel
	GemGroupStep          map[int]int //group:step
	MiliWeekTaskConfig    []*MiliWeekTaskConfig
	MiliWeekTaskMapConfig map[int]*MiliWeekTaskConfig
	HeroStar              []*HeroStar               // 英雄升星配置
	HeroStarMap           map[int]map[int]*HeroStar // heroId, star, config
	SimpleConfigMap       map[int]*SimpleConfig
	HeroStarAttrMap       map[int]map[int]*HeroStarAttr    // heroId, star, attr
	SkillConfigMap        map[int]*HeroSkill               // 技能配置
	TalentConfig          []*TalentConfig                  // 技能配置
	TalentMap             map[int]map[int]*TalentConfig    // 天赋配置, talentId, lv, config
	TalentAwakeConfig     []*TalentAwake                   // 天赋觉醒所有配置
	TalentAwakeMap        map[int][]*TalentAwake           // 天赋觉醒配置
	DreamLandConfig       []*DreamLand                     // 神格幻境所有配置
	DreamLandGroupMap     map[int][]*DreamLand             // 神格幻境组配置
	DreamLandItemMap      map[int]*DreamLand               // 神格幻境物品配置
	DreamLandSpendConfig  []*DreamLandSpend                // 幻境花费所有配置
	DreamLandSpendMap     map[int]map[int]*DreamLandSpend  // 幻境花费配置
	DreamLandCostMap      map[int]*DreamLandCost           // 幻境花费配置
	BecomeStrongerConfig  []*BecomeStronger                // 我要变强配置
	BecomeStrongerMap     map[int]*BecomeStronger          // 我要变强配置
	MaxTalentLv           map[int]int                      // 天赋最大等级
	FateConfig            []*FateConfig                    // 缘分信息
	FateMap               map[int][]*FateConfig            // 英雄拥有的缘分
	HeroConfig            []*HeroConfig                    // 武将配置
	HeroConfigMap         map[int]map[int]*HeroConfig      // 武将配置
	HeroNumMap            map[int]*HeroNumConfig           // 武将扩展配置
	HeroBreakConfig       []*HeroBreakConfig               // 武将配置
	HeroBreakConfigMap    map[int]map[int]*HeroBreakConfig // 武将配置
	HeroHandBookConfigMap map[int]*HeroHandBookConfig      // 图鉴
	HeroGrowthConfig      []*HeroGrowthConfig              // 成长表
	HeroGrowthConfigMap   map[int]*HeroGrowthConfig        // 成长表
	TeamDungeonMap        map[int]*TeamDungeonConfig       // 地下城组队关卡配置
	ItemMap               map[int]*ItemConfig              // 道具配置

	EquipValueGroupMap   map[int]*EquipValueGroup   //
	EquipBaseValueMap    map[int]*EquipBaseValue    //
	EquipSpecialValueMap map[int]*EquipSpecialValue //
	EquipHireConfig      []*EquipHireConfig         //

	EquipUpgrade      []*EquipUpgrade               // 强化装备
	EquipUpgradeMap   map[int]map[int]*EquipUpgrade // 强化装备map
	EquipUpgradeMaxLv map[int]int                   // 强化最大等级
	EquipStar         []*EquipStar                  // 附魔装备
	EquipStarMap      map[int]map[int]*EquipStar    // 附魔装备map
	EquipStarMaxLv    map[int]int                   // 附魔最大等级
	EquipGem          map[int]*EquipGem             // 附魔装备
	EquipSuit         map[int]*EquipSuit            // 套装信息
	GemLevelup        map[int]int                   // 宝石合成到下一级
	ItemBagGroupMap   map[int]*ItemBagGroup         // 掉落组配置
	ItemBagMap        map[int]*ItemBag              // 掉落包配置
	HeroDecomposeMap  map[int]*HeroDecompose        // 英雄分解
	LevelConfigMap    map[int]*LevelConfig          // 关卡配置
	CostMap           map[int]*CostConfig           // 消耗配置
	TechConfigMap     map[int]*TechConfig           // 科技升级消耗
	TechAttr          []*TechAttr                   // 科技属性
	TechTime          []*TechTime                   // 科技时间
	MaxTechLv         map[int]int                   // 科技满级
	NewUserItem       []*NewUserItem                // 新手玩家赠送道具

	AstrologyDropConfig      map[int][]*AstrologyDropConfig // 占星
	AstrologyDropGroupConfig map[int]*AstrologyDropConfig   // 占星

	SyntheticDropConfig      map[int][]*SyntheticDropConfig // 紫灵魂石合成算法
	SyntheticDropGroupConfig map[int]*SyntheticDropConfig   // 紫灵魂石合成算法

	WarplayConfig           []*WarplayConfig
	CamprewardConfig        []*CamprewardConfig
	TitleConfig             map[int]*TitleConfig
	BespeakConfig           map[int]*BespeakConfig
	BuildingingConfig       map[int]*BuildingingConfig
	ConsumetopattackConfig  map[int]*ConsumetopattackConfig
	ResourcelevelConfig     []*ResourcelevelConfig
	ActiveConfig            map[int]*ActiveConfig
	PubchestdropgroupConfig map[int]*PubchestdropgroupConfig
	PubchestdropgroupLst    map[int][]*PubchestdropgroupConfig
	PubchestdropgroupSum    map[int]int

	HorseJudgecallConfig map[int]*HorseJudgecallConfig
	HorseSoulConfig      map[int]*HorseSoulConfig
	PubchesttotalConfig  map[int]*PubchesttotalConfig
	PubchesttotalGroup   map[int][]*PubchesttotalConfig
	TimeResetConfig      []*TimeResetConfig

	WarOrderConfig map[int]*WarOrderConfig
	WarOrderParam  map[int]*WarOrderParam

	WarOrderLimitConfig map[int]*WarOrderLimitConfig

	AccessAwardConfig map[int]*AccessAwardConfig
	AccessTaskConfig  map[int]*AccessTaskConfig
	AccessRankConfig  []*AccessRankConfig

	ActivityBossRankConfig     map[int]map[int][]*ActivityBossRankConfig
	ActivityBossConfig         []*ActivityBossConfig
	ActivityBossTargetConfig   map[int]map[int][]*ActivityBossTargetConfig
	ActivityBossExchangeConfig map[int]map[int][]*ActivityBossExchangeConfig

	RankRewardConfigMap map[int]map[int][]*RankRewardConfig

	ActivityDailyRecharge map[int]*ActivityDailyRecharge // 每日礼包
	ActivityOverflowGifts []*ActivityOverflowGifts       // 超值好礼

	TimeGeneralsPointsConfig        []*TimeGeneralsPointsConfig
	HorseBattleSteedAttributeConfig map[int]*HorseBattleSteedAttributeConfig
	ExcitingConfig                  []*ExcitingConfig
	LeveltitletaskConfig            map[int]*LeveltitletaskConfig
	PlayernameConfig                []*PlayernameConfig
	WarcityConfig                   map[int]*WarcityConfig
	WartargetConfig                 map[int]*WartargetConfig
	CommunityConfig                 map[int]*CommunityConfig
	VipConfigMap                    map[int]*VipConfig
	PeoplecityConfig                map[int]*PeoplecityConfig
	BuyphysicalConfig               map[int]*BuyphysicalConfig
	ConsumetopluckConfig            map[int]*ConsumetopluckConfig
	WarcontributionConfig           map[int]*WarcontributionConfig
	WarbuffConfig                   []*WarBuffConfig
	GemchestConfig                  map[int]*GemchestConfig
	ShopConfig                      map[int]*ShopConfig
	ShopGrid                        map[int]map[int][]*ShopConfig
	ShopSumGrid                     map[int]map[int]int
	TreasuryConfig                  map[int]*TreasuryConfig
	HomeofficepurchaseConfig        []*HomeofficepurchaseConfig
	ExpeditionbuffConfig            map[int]*ExpeditionbuffConfig
	SkillConfig                     map[int]*SkillConfig
	BeautychestConfig               map[int]*BeautychestConfig
	BuyskillpointsConfig            map[int]*BuyskillpointsConfig
	GrowthtaskConfig                map[int]*GrowthtaskConfig
	TrialConfig                     map[int]*TrialConfig
	WorldlvtpyeConfig               map[int]*WorldlvtpyeConfig
	CitywayConfig                   map[int]*CitywayConfig
	GemsweeperrankConfig            map[int]*GemsweeperrankConfig
	NpcConfig                       map[int]*NpcConfig
	DailytaskConfig                 map[int]*DailytaskConfig
	TargetTaskConfig                map[int]*TargetTaskConfig
	BadgeTaskConfig                 []*BadgeTaskConfig
	HalfmoonConfig                  map[int]*HalfmoonConfig
	HolyUpgradeConfig               []*HolyUpgradeConfig
	HolyUpgradeMap                  map[int]map[int]*HolyUpgradeConfig
	GemsweeperitembagConfig         []*GemsweeperitembagConfig
	PubchestspecialConfig           map[int]*PubchestspecialConfig
	PubChestSpecialLst              map[int][]*PubchestspecialConfig
	SummonBoxConfig                 []*SummonBoxConfig
	//TigerUpgradeConfig              map[int]*TigerUpgradeConfig
	ConsumetophpConfig           map[int]*ConsumetophpConfig
	HorseJudgeDiscernConfig      map[int]*HorseJudgeDiscernConfig
	HorseBattleSteedAwakenConfig map[int]*HorseBattleSteedAwakenConfig
	DiplomacychangesConfig       []*DiplomacychangesConfig
	ConsumetoplistConfig         []*ConsumetoplistConfig
	GemsweeperawardConfig        map[int]*GemsweeperawardConfig
	HomeofficecityawardConfig    []*HomeofficecityawardConfig
	LegioncopyConfig             map[int]*LegioncopyConfig
	WartaskConfig                map[int]*WartaskConfig
	TimeGeneralsConfig           map[int]*TimeGeneralsConfig
	HorseSoulUpgradeConfig       map[int]*HorseSoulUpgradeConfig
	RobotConfig                  map[int]*RobotConfig
	NationalwarEnrollmentConfig  map[int]*NationalwarEnrollmentConfig
	CamptaskConfig               map[int]*CamptaskConfig
	ExpeditionConfig             map[int]*ExpeditionConfig
	HorseParmConfig              map[int]*HorseParmConfig
	WarlistConfig                map[int]*WarlistConfig
	HolyLegendLevelConfig        []*HolyLegendLevelConfig
	ExpeditionbuffgroupConfig    map[int]*ExpeditionbuffgroupConfig
	ShoprefreshConfig            map[int]*ShoprefreshConfig
	PromoteboxConfig             map[int]*PromoteboxConfig
	ActivitynewConfig            map[int]*ActivitynewConfig
	ChestshopConfig              map[int]*ChestshopConfig
	HonorkillConfig              map[int]*HonorkillConfig
	TimeGeneralsRankawardConfig  []*TimeGeneralsRankawardConfig
	SevendayConfig               map[int]*SevendayConfig
	SevendayAward                []*SevendayAward
	IndustryConfig               map[int]*IndustryConfig
	SoulchestConfig              map[int]*SoulchestConfig
	ConsumetopshopConfig         []*ConsumetopshopConfig
	HolyLegendConfig             []*HolyLegendConfig
	//TigerStuntConfig                map[int]*TigerStuntConfig
	HorseJudgelevelConfig map[int]*HorseJudgelevelConfig
	CrownfightConfig      map[int]*CrownfightConfig
	CrownfightConfigLst   []*CrownfightConfig

	FundbuyConfig          map[int]*FundbuyConfig
	HolyPartsUpgradeConfig []*HolyPartsUpgradeConfig
	HolyPartsMap           map[int]map[int]*HolyPartsUpgradeConfig
	VisitConfig            map[int]*VisitConfig
	//ActivityboxConfig      map[int]*ActivityboxConfig
	JjcdamoninfoConfig map[int]*JjcdamoninfoConfig
	//TigerAttributeConfig            map[int]*TigerAttributeConfig
	SpyrewardConfig           map[int]*SpyrewardConfig
	MaincityConfig            map[int]*MaincityConfig
	SignrewardConfig          map[int]*SignrewardConfig
	FundConfig                map[int]*FundConfig
	GemsweeperextradropConfig map[int]*GemsweeperextradropConfig
	ActivityFundConfig        []*ActivityFundConfig
	ActivityFundMap           map[int]*ActivityFundGroupMap
	//ActivityFundTypeMap       map[int]*ActivityFundTypeConfig

	OfficerobotConfig       map[int]*OfficerobotConfig
	HomeofficerankingConfig []*HomeofficerankingConfig
	SignConfig              map[int]*SignConfig
	SignMap                 map[int]*SignConfig
	ConsumetopbossConfig    []*ConsumetopbossConfig
	FitConfig               map[int]*FitConfig
	TeamexpConfig           map[int]*TeamexpConfig
	TreasureSuitConfig      map[int]*TreasureSuitConfig
	HorseBattleSteedConfig  map[int]*HorseBattleSteedConfig
	SpytreasureConfig       map[int]*SpytreasureConfig
	MoneyConfig             map[int]*MoneyConfig
	WarencourageConfig      map[int]*WarencourageConfig
	//RedpacketmoneyConfig            map[int]*RedpacketmoneyConfig
	EquipchestConfig      map[int]*EquipchestConfig
	LevelboxConfig        map[int]*LevelboxConfig
	LevelMapConfig        map[int]*LevelMapConfig
	HorseShopConfig       map[int]*HorseShopConfig
	VisitchanceConfig     []*VisitchanceConfig
	NationalwardConfig    map[int]*NationalwarawardConfig
	NationalwarParmConfig map[int]*NationalwarParmConfig

	ArenaAwardConfig    []*PvpAwardConfig
	ArenaAwardConfigMap map[int]*PvpAwardConfig

	TariffConfig []*TariffConfig
	BossConfig   map[int]*BossConfig
	HangUpConfig map[int]*HangUpConfig

	GemstoneChapterConfig    []*GemstoneChapterConfig
	GemstoneChapterConfigMap map[int]*GemstoneChapterConfig
	GemstoneLevelConfig      []*GemstoneLevelConfig
	GemstoneLevelConfigMap   map[int]*GemstoneLevelConfig
	HorseAward               []*HorseAward
	HorseAwardMap            map[int]*HorseAward
	HorseSwitchMap           map[int]*HorseSwitch
	HorseSwitchId            int
	HorseSwitchNum           int
	HorseBattleSteedMap      map[int]*HorseBattleSteed
	ClearWarHorseId          int
	ClearWarHorseNum         int
	SmeltAward               []*SmeltAward
	SmeltBuyAward            []*SmeltAward
	SmeltAwardMap            map[int]*SmeltAward
	SmeltBuyAwardMap         map[int]*SmeltAward
	MercenaryRandom          []*MercenaryRandom
	MercenaryLv              []*MercenaryLv
	MercenaryLvMap           map[int]map[int]*MercenaryLv
	MercenaryConfig          map[int]*MercenaryConfig
	MercenaryRandomGroup     map[int]map[int][]*MercenaryRandom
	WorldPowerMap            map[int]*WorldPower
	TaskKingConfig           []*GrowthtaskKingConfig
	TaskKingConfigMap        map[int]*GrowthtaskKingConfig
	TaskKingGroupMap         map[int][]int
	HeroAttribute            map[int]*HeroAttribute
	CrownBuildConfig         map[int]*CrownBuildConfig
	LevelItemMap             map[int]*LevelItemConfig

	StatisticsConfig          []*StatisticsConfig
	StatisticsConfigMap       map[int]*StatisticsConfig
	StatisticsRewardConfig    []*StatisticsRewardConfig
	StatisticsRewardConfigMap map[int]*StatisticsRewardConfig

	NobilityConfig    []*NobilityConfig
	NobilityConfigMap map[int]*NobilityConfig
	NobilityReward    []*NobilityReward
	NobilityRewardMap map[int]*NobilityReward

	ActivityTotleAward    []*ActivityTotleAward
	ActivityTotleAwardMap map[int]*ActivityTotleAward

	MonthCardTotleAwardMap map[int]*MonthCardTotleAwardMap
	MonthCard              map[int]*MonthCard

	WholeShopConfig    []*WholeShopConfig
	WholeShopConfigMap map[int]*WholeShopConfig

	WholeShopTimeConfig    []*WholeShopTimeConfig
	WholeShopConfigTimeMap map[int]*WholeShopTimeConfig

	TurnTableConfigMap     map[int]*TurnTableConfig
	TurnTableTimeConfigMap map[int]*TurnTableTimeConfig

	LotteryDrawConfigMap map[int]*LotteryDrawConfig

	RechargeConfig    []*RechargeConfig
	RechargeConfigMap map[int]*RechargeConfig

	HeroExpConfig    []*HeroExpConfig // 英雄升级相关配置
	HeroExpConfigMap map[int]*HeroExpConfig

	RuneCompose    []*RuneCompose //符文合成，包括装备合成
	RuneComposeMap map[int]*RuneCompose

	RuneConfig    []*RuneConfig //符文合成，包括装备合成
	RuneConfigMap map[int]*RuneConfig

	FormationConfig    []*FormationConfig //阵型配置
	FormationConfigMap map[int]*FormationConfig

	FundConfigMap map[int]*FundConfigMap //超值基金

	NewShopConfigMap map[int]map[int][]*NewShopConfig
	NewShopDiscount  []*NewShopDiscount

	EquipShopRate          map[int]int
	EquipConfigMap         map[int]*EquipConfig             // 装备类型
	EquipAdvancedConfigMap map[int]*EquipAdvancedConfig     // 升阶
	EquipStrengthenMap     map[int][]*EquipStrengthenConfig //
	EquipStrengthenUpLvMap []*EquipStrengthenLvUpConfig     //
	EquipRecastConfig      []*EquipRecastConfig             //

	ArtifactEquipConfigMap map[int]*ArtifactEquipConfig // 神器
	ArtifactStrengthen     []*ArtifactStrengthenConfig  //神器升级

	ExclusiveEquipConfigMap map[int]*ExclusiveEquipConfig              // 专属
	ExclusiveStrengthen     map[int]map[int]*ExclusiveStrengthenConfig //专属升级

	//HydraConfig    []*HydraConfig
	HydraConfigMap map[int]*HydraConfig
	HydraSkill     []*HydraSkill
	HydraSkillMap  map[int]map[int]*HydraSkill
	//HydraLevel    []*HydraLevel
	HydraLevelMap map[int]*HydraLevel
	HydraStepMap  map[int]*HydraStep
	//HydraStar     []*HydraStar
	//HydraStarMap  map[int]map[int]*HydraStar
	//HydraTask    []*HydraTask
	//HydraTaskMap map[int]*HydraTask

	PitConfig     []*PitConfig
	PitConfigMap  map[int]*PitConfig
	PitMap        []*PitMap //地牢地图
	PitMapMap     map[int][]*PitMap
	PitBuffMap    map[int]*PitBuff
	PitBoxMap     map[int]*PitBox
	PitMonsterMap map[int]*PitMonster

	//新地牢
	NewPitConfigMap            map[int]*NewPitConfig
	NewPitRelique              map[int]*NewPitRelique
	NewPitTreasureCave         map[int]*NewPitTreasureCave
	NewPitRobotConfig          []*NewPitRobotConfig
	NewPitRobotExclusiveConfig []*NewPitRobotExclusiveConfig
	NewPitRobotGroupMap        map[int][]int
	NewPitRobotGroupFirstMap   map[int][]int
	NewPitRobotQuality         []*NewPitRobotQuality
	NewPitRobotMonsterQuality  []*NewPitRobotMonsterQuality
	//NewPitRobotMonsterLv      []*NewPitRobotMonsterLv
	NewPitRobotAttr   []*NewPitRobotAttr
	NewPitExtraReward map[int]*NewPitExtraReward
	NewPitParam       map[int]*NewPitParam
	NewPitDifficulty  []*NewPitDifficulty

	//时光之巅
	InstanceConfig map[int]*InstanceConfig
	InstanceBox    map[int][]*InstanceBox
	InstanceThing  map[int]map[int]*InstanceThing

	PlayNotice       map[int]*PlayNoticeInfo //! 矿点玩法公告
	GveNotice        map[int]*PlayNoticeInfo //! 孤山多宝玩法公告
	UnionFightNotice map[int]*PlayNoticeInfo //! 军团战玩法公告
	CampNotice       map[int]*PlayNoticeInfo //! 国战玩法公告
	PlayTimeMap      map[int]*PlayTimeConfig //! 国战PVP玩家时间
	PlayRewardList   []*PlayRewardConfig     //! 国战奖励配置

	OpenLevelMap map[int]*OpenLevelConfig
	MailConfig   map[int]*MailConfig

	WarEncourageConfig map[int]*WarEncourageConfig
	EncourageMap       map[int]map[int]*WarEncourageConfig
	HeadConfigMap      map[int]*HeadConfig

	WorldLevel map[int]*WorldLevelConfig
	WorldMap   map[int]*WorldMapConfig

	MoneyTaskMap      map[int]*MongytaskListConfig
	MoneyTaskStarList []*MongytaskStarlistConfig
	MoneyTotal        int
	//HeadItemMap       map[int]int
	MissionMap map[int]*MissionInfo

	StateOwner   map[int][]int                 // 州有哪些城池
	StateName    map[int]string                // 州名字
	StateBox     map[int]int                   // 州宝箱
	BossLvConfig []*BossLvConfig               // 巨兽等级信息
	BossAttMap   map[int]map[int]*BossLvConfig // 巨兽属性信息
	StrConfig    []*StrStringConfig            // 中文配置
	StrMap       map[string]string             // 中文映射

	//! 独立活动相关
	ActivityTimeGiftMap   map[int]*TimeGiftConfig   //! 限时礼包
	ActivityTimeGiftGroup map[int][]*TimeGiftConfig //! 分组-限时礼包

	//! 阵容属性加成
	TeamAttrConfig []*TeamAttrConfig //! 阵容属性加成

	ActivityGiftConfig []*ActivityGiftConfig
	GrowthGiftConfig   []*GrowthGiftConfig

	EntanglementConfig    map[int]*EntanglementConfig // 羁绊系统
	EntanglementMapConfig map[int]*EntanglementFate   // 羁绊系统

	RewardForbarConfig         []*RewardForbarConfig                 //悬赏
	RewardForbarMapConfig      map[int][]*RewardForbarConfig         // 悬赏系统
	RewardForbarColorMapConfig map[int]map[int][]*RewardForbarConfig // 悬赏系统

	RewardForbarLvUpConfig []*RewardForbarLvUpConfig
	RewardForbarAward      []*RewardForbarPrize
	RewardForbarAwardMap   map[int]map[int][]*RewardForbarPrize

	RankListIntegral  []*RankListIntegral       // 凯旋丰碑
	RankTaskConfig    []*RankTaskConfig         // 排行榜任务
	RankTaskMapConfig map[int][]*RankTaskConfig // 排行榜任务

	ResonanceCrystalconfig []*ResonanceCrystalconfig // 共鸣水晶

	ArenaRewardConfig   []*ArenaRewardConfig // 竞技场奖励
	ArenaParameter      []*ArenaParameterConfig
	UnionHuntConfig     []*UnionHuntConfig     // 军团狩猎
	UnionHuntDropConfig []*UnionHuntDropConfig // 军团狩猎

	ArenaSpecialClassConfig []*ArenaSpecialClass         // 高阶竞技场段位
	ArenaSpecialClassMap    map[int][]*ArenaSpecialClass // 高阶竞技场段位

	HeroSkinConfig []*HeroSkin // 英雄皮肤

	HeroGrowConfig []*HeroGrowConfig //

	CrossArenaRewardConfig []*CrossArenaRewardConfig
	CrossArenaSubsection   []*CrossArenaSubsection

	CrossArena3V3RewardConfig []*CrossArena3V3RewardConfig
	CrossArena3V3Subsection   []*CrossArena3V3Subsection

	ActivityBuyLimit   []*ActivityBuyLimit        //限时抢购
	ActivityBuyItem    []*ActivityBuyItem         //奖励列表
	ActivityMapBuyItem map[int][]*ActivityBuyItem //奖励列表

	//TreeConfig                []*TreeConfig
	TreeLevelConfig           []*TreeLevel                // 生命树等级
	TreeProfessionalConfig    []*TreeProfessional         // 生命树专业等级
	TreeProfessionalMapConfig map[int][]*TreeProfessional // 生命树专业等级

	//星座
	InterstellarConfig   map[int]*InterstellarConfig
	InterstellarHangup   map[int]*InterstellarHangup
	InterstellarWar      map[int]map[int]*InterstellarWar
	InterstellarBox      map[int]map[int]*InterstellarBox
	InterstellarTaskNode map[int]*TaskNode

	StageTalentConfig []*StageTalentConfig
	StageTalentMap    map[int][]*StageTalentConfig
} // 所有结构定义

type MissionInfo struct {
	Chapter    int   `json:"chapter"`
	MissionIds []int `json:"mission_ids"`
}

type SmeltAward struct {
	Id      int   `json:"id"`
	Times   int   `json:"times"`
	Rewards []int `json:"reward"`
	Nums    []int `json:"num"`
}

// 英雄升星
type HeroStar struct {
	HeroId        int   `json:"id"`
	Star          int   `json:"hero_star"`
	SlotItemId1s  []int `json:"qualification1_id"`
	SlotItemNum1s []int `json:"qualification1_num"`
	SlotItemId2s  []int `json:"qualification2_id"`
	SlotItemNum2s []int `json:"qualification2_num"`
	SlotItemId3s  []int `json:"qualification3_id"`
	SlotItemNum3s []int `json:"qualification3_num"`
	SlotItemId4s  []int `json:"qualification4_id"`
	SlotItemNum4s []int `json:"qualification4_num"`
	SlotItemId5s  []int `json:"qualification5_id"`
	SlotItemNum5s []int `json:"qualification5_num"`
	SlotItemId6s  []int `json:"qualification6_id"`
	SlotItemNum6s []int `json:"qualification6_num"`
	//老代码，配置有嵌套，要想想怎么改，时间问题暂时无力反抗  20190927 by zy
	SlotItemId7s   []int `json:"qualification7_id"`
	SlotItemNum7s  []int `json:"qualification7_num"`
	SlotItemId8s   []int `json:"qualification8_id"`
	SlotItemNum8s  []int `json:"qualification8_num"`
	SlotItemId9s   []int `json:"qualification9_id"`
	SlotItemNum9s  []int `json:"qualification9_num"`
	SlotItemId10s  []int `json:"qualification10_id"`
	SlotItemNum10s []int `json:"qualification10_num"`

	StarLvIds    []int   `json:"starlv_id"`
	StarLvNums   []int   `json:"starid_num"`
	Attr1        []int   `json:"qualification1_type"`
	Value1       []int64 `json:"qualification1_value"`
	Attr2        []int   `json:"qualification2_type"`
	Value2       []int64 `json:"qualification2_value"`
	Attr3        []int   `json:"qualification3_type"`
	Value3       []int64 `json:"qualification3_value"`
	Attr4        []int   `json:"qualification4_type"`
	Value4       []int64 `json:"qualification4_value"`
	Attr5        []int   `json:"qualification5_type"`
	Value5       []int64 `json:"qualification5_value"`
	Attr6        []int   `json:"qualification6_type"`
	Value6       []int64 `json:"qualification6_value"`
	Attr7        []int   `json:"qualification7_type"`
	Value7       []int64 `json:"qualification7_value"`
	Attr8        []int   `json:"qualification8_type"`
	Value8       []int64 `json:"qualification8_value"`
	Attr9        []int   `json:"qualification9_type"`
	Value9       []int64 `json:"qualification9_value"`
	Attr10       []int   `json:"qualification10_type"`
	Value10      []int64 `json:"qualification10_value"`
	StarLvTypes  []int   `json:"starlv_type"`
	StarLvValues []int64 `json:"starlv_value"`
	SkillIds     []int   `json:"skillid"`
}

// 通用配置
type SimpleConfig struct {
	Id   int    `json:"id"`
	Num  int    `json:"num"`
	Text string `json:"text"`
	ItemsId []int	`json:"item_id_"`
	ItemsNum []int `json:"item_num_"`
}

// 属性
type Attribute struct {
	AttType  int   `json:"attrtype"`  //! 类型
	AttValue int64 `json:"attrvalue"` //! 值
}

// 升星配置
type HeroStarAttr struct {
	SlotAttr [maxStarSlots][]*Attribute // 资质属性
	StarAttr []*Attribute               // 升星属性
}

// 技能配置
type HeroSkill struct {
	SkillId        int   `json:"skill_id"`         // 技能Id
	SkillValueType []int `json:"skill_value_type"` // 技能数值类型
	SkillCount     []int `json:"skill_count"`      // 技能初始值
	SkillAddType   int   `json:"skill_add_type"`   // 技能增加战斗力类型
	SkillAddValue  int   `json:"skill_add_value"`  // 技能增加战斗力值
	SkillType      int   `json:"skill_type"`       // 技能类型
}

type SkillConfig struct {
	Skillid         int     `json:"skillid"`
	Staropen        int     `json:"staropen"`
	Learnlv         int     `json:"learnlv"`
	Maxlv           int     `json:"maxlv"`
	Skillvaluetypes []int   `json:"skillvaluetype"`
	Skillcounts     []int64 `json:"skillcount"`
	Skilladdvalues  []int64 `json:"skilladdvalue"`
	Skilltype       int     `json:"skilltype"`
}

// 天赋配置
type TalentConfig struct {
	TalentId  int     `json:"id"`
	TalentLv  int     `json:"lv"`
	NeedLevel int     `json:"needlv"`
	CostItems []int   `json:"item"`
	Costnums  []int   `json:"num"`
	AttTypes  []int   `json:"base_type"`
	AttValues []int64 `json:"base_value"`
	SkillId   int     `json:"skill"`
	SkillId2  int     `json:"skill_fuck"`
	//SkillLv     int     `json:"skill_lv"`
	ReturnItems []int `json:"return_type"`
	Returnnums  []int `json:"return_value"`
}

type TalentAwake struct {
	Group      int     `json:"group"`
	Step       int     `json:"step"`
	LevelLimit int     `json:"alllv"`
	AttTypes   []int   `json:"base_type"`
	AttValues  []int64 `json:"base_value"`
}

// 神格抽奖配置
type DreamLand struct {
	ID            int `json:"id"`
	Group         int `json:"group"`
	Item          int `json:"item"`
	Num           int `json:"num"`
	RefreshWeight int `json:"refreshweight"`
	ExtractWeight int `json:"extractweight"`
	Only          int `json:"only"`
	HaveHero      int `json:"havediv"`
	HaveChance    int `json:"havedivchance"`
	Notice        int `json:"notice"`
}

// 神格抽奖消耗配置
type DreamLandSpend struct {
	ID       int   `json:"id"`    // id
	Class    int   `json:"class"` //
	RefCost  int   `json:"refcost"`
	LootCost int   `json:"luckcost"`
	Minvip   int   `json:"minvip"` //
	Maxvip   int   `json:"maxvip"`
	Type     int   `json:"type"`
	Total    int   `json:"total"`
	Group    []int `json:"group"`
	//Num      []int `json:"num"`
	Chance    int `json:"chance"`
	TypeTimes int `json:"typetime"`
}

type DreamLandCost struct {
	Type      int `json:"type"`
	RefCost   int `json:"refcost"`
	LootCost  int `json:"luckcost"`
	TypeTimes int `json:"typetime"`
}

// 缘分配置
type FateConfig struct {
	FateId    int     `json:"fate_num"`
	HeroId    int     `json:"main_hero"`
	Heroes    []int   `json:"hero_"`
	FateType  int     `json:"fate_type"`
	FateParam []int   `json:"fate_parm"`
	AttType   []int   `json:"attribute"`
	AttValue  []int64 `json:"value"`
}

// 武将配置
type HeroConfig struct {
	HeroId         int     `json:"id"`           // 武将Id
	HeroName       string  `json:"heroname"`     // 武将名字
	FullName       string  `json:"fullname"`     // 武将名字
	Attribute      int     `json:"attribute"`    // 阵营
	AttackType     int     `json:"attacktype"`   // 职业
	HeroStar       int     `json:"hsrostar"`     // 武将星级
	FinalLevel     int     `json:"finallevel"`   // 武将最大等级
	HeroLvMax      int     `json:"hsrolvmax"`    // 武将等级上线
	FightIndex     int     `json:"fight_index"`  // 战斗index, 品质
	HeroCamp       int     `json:"hero_camp"`    // 武将类型, 0 物理 1 法术
	CardIds        []int   `json:"debris_type"`  // 英雄整卡转换道具类型
	CardNums       []int   `json:"debris_value"` // 英雄整卡转换道具数量
	HeroArms       int     `json:"hero_arms"`
	BaseTypes      []int   `json:"base_type"`
	BaseValues     []int64 `json:"base_value"`
	DivinityAwaken int     `json:"divinity_awaken"`
	Point          []int   `json:"point"`
	GrowthTypes    []int   `json:"growth_type"`
	GrowthValues   []int64 `json:"growth_value"`
	Career         int     `json:"career_info"`
	RuneSkill      []int   `json:"rune_skill"`
	UpStarType     []int   `json:"type"`        //升星类型
	UpStarStar     []int   `json:"star"`        //升星需求
	UpStarNum      []int   `json:"num"`         //升星需求数量
	DisbandId      []int   `json:"disband_id"`  //英雄遣散物品IDs
	DisbandNum     []int   `json:"disband_num"` //英雄遣散物品数量
	BackItem       []int   `json:"back_item"`   //英雄回退物品IDs
	BackNum        []int   `json:"back_num"`    //英雄回退物品数量
	WishHero       int     `json:"wishhero"`    //是否心愿单英雄
	QuaType        []int   `json:"quatype"`
	QuaValue       []int64 `json:"quavalue"`
	TreeID         int     `json:"tree_id"`
	TreeNum        int     `json:"tree_num"`
	TalentGroup    int     `json:"talent_group"`
}

type HeroNumConfig struct {
	BuyNum     int `json:"buynum"`     // 购买次数
	MaxNum     int `json:"hmaxnum"`    // 最大数量
	BuyNeed    int `json:"buyneed"`    // 扩展需求id
	BuyNeedNum int `json:"buyneednum"` // 扩展需求num
}

type AstrologyDropConfig struct {
	Id              int `json:"id"`               // 唯一ID
	AstrologyId     int `json:"astrologyid"`      // 大组ID
	AstrologyChance int `json:"astrology_chance"` // 大组概率
	AstrologyMid    int `json:"astrology_mid"`    // 达到分值概率
	ScoreLimit      int `json:"scorelim"`         // 分值阀值
	ItemId          int `json:"itemid"`           // 道具ID
	ItemNum         int `json:"itemnum"`          // 道具数量
	ItemWT          int `json:"itemwt"`           // 道具权重
	ScoreValue      int `json:"scorevalue"`       // 分数正负号 0正号 1负号
	ItemScore       int `json:"itemscore"`        // 物品分数
}

type SyntheticDropConfig struct {
	Id              int `json:"id"`               // 唯一ID
	SyntheticId     int `json:"syntheticid"`      // 大组ID
	SyntheticChance int `json:"synthetic_chance"` // 大组概率
	SyntheticMid    int `json:"synthetic_mid"`    // 达到分值概率
	ScoreLimit      int `json:"scorelim"`         // 分值阀值
	ItemId          int `json:"itemid"`           // 道具ID
	ItemNum         int `json:"itemnum"`          // 道具数量
	ItemWT          int `json:"itemwt"`           // 道具权重
	ScoreValue      int `json:"scorevalue"`       // 分数正负号 0正号 1负号
	ItemScore       int `json:"itemscore"`        // 物品分数
}

type HeroBreakConfig struct {
	Id         int     `json:"id"`     // Id
	HeroId     int     `json:"heroid"` // 武将Id
	Break      int     `json:"bresk"`  //
	BaseTypes  []int   `json:"type"`
	BaseValues []int64 `json:"value"`
	Skill      []int   `json:"skillhave"`
}

type HeroHandBookConfig struct {
	Id       int `json:"hid"`      // Id
	Prize    int `json:"prize"`    //
	PrizeNum int `json:"prizenum"` //
	Showfdw  int `json:"showfdw"`  //
}

type HeroGrowthConfig struct {
	Id          int     `json:"id"`          // Id
	GrowthLevel int     `json:"growthlevel"` //
	Type        int     `json:"type"`        //
	GrowthType  []int   `json:"growthtype"`  //
	GrowthValue []int64 `json:"growthvalue"` //
}

// 组队地下城关卡配置
type TeamDungeonConfig struct {
	LevelId          int    `json:"levelid"`          // 关卡id
	ChapterName      string `json:"chaptername"`      // 关卡名称
	Refresh          []int  `json:"refresh"`          // 刷新时间
	ServerNum        int    `json:"servernum"`        // 刷新次数
	ColdTime         int64  `json:"coldtime"`         // CD时间
	PhysicalStrength int    `json:"physicalstrength"` //需求体力
	ForceLimit       int    `json:"force_limit"`      //需求战力
	TeamExp          int    `json:"teamexp"`          //经验值
	DropId           int    `json:"drop"`             //掉落，取Lottery_New
	Win              int    `json:"win"`
	Lose             int    `json:"lose"`
	Reduce           int    `json:"reduce"`
	MaxChance        int    `json:"maxchance"`
	MinChance        int    `json:"minchance"`
	JudgeLevel       int    `json:"judgelevel"`
	LevelIndex       int    `json:"levelindex"`
	SweepCost        int    `json:"sweep_cost"`
	SweepExp         int    `json:"sweep_exp"`
	SweepDrop        int    `json:"sweep_drop"`
}

// 道具配置
type ItemConfig struct {
	ItemId          int    `json:"itemid"`          // 道具Id
	ItemName        string `json:"itemname"`        // 道具名字
	ItemType        int    `json:"itemtype"`        // 道具类型
	ItemSubType     int    `json:"itemsubtype"`     // 子类型
	ExchangeExp     int    `json:"exchangeexp"`     // 经验值
	ExchangeGold    int    `json:"exchangegold"`    // 金币
	WoodPrice       int    `json:"woodprice"`       // 木材价格
	ClothPrice      int    `json:"clothprice"`      // 衣服价格
	IronPrice       int    `json:"ironprice"`       // 铁价格
	GoldPrice       int    `json:"goldprice"`       // 元宝价格
	SericePrice     int    `json:"serviceprice"`    // 价格
	EquipPrice      int    `json:"equipprice"`      // 价格
	HonorPrice      int    `json:"honorprice"`      // 荣誉价值
	ExpeditionPrice int    `json:"expeditionprice"` // 远征价格
	DevotePrice     int    `json:"devoteprice"`     // 价格
	UniquePrice     int    `json:"uniqueprice"`     // 价格
	OfficePrice     int    `json:"officeprice"`     // 价格
	ItemCheck       int    `json:"itemcheek"`       // 道具检查
	NeedLv          int    `json:"needlv"`          // 需要等级
	BuyLv           int    `json:"buylv"`           // 购买等级
	Special         int    `json:"special"`         // 特殊
	Spirititem      int    `json:"spirititem"`      // ...
	Spiritexp       int    `json:"spiritexp"`       // ...
	Sort            int    `json:"sort"`            // ...
	Rareeffect      int    `json:"rareeffect"`      // ...
	Vip             int    `json:"vip"`             // ...
	GemPrice        int    `json:"gemprice"`        // 道具价格(钻石)
	CompoundId      int    `json:"cardcomposite"`   // 道具合成Id(装备Id)
	CompoundNum     int    `json:"cardnum"`         // 道具合成数量(碎片数量)
	MaxNum          int    `json:"overlapnum"`      // 道具最大上限
	LotteryId       int    `json:"lotteryid"`       //
	Overflow        int    `json:"overflow"`        //
	LotDrop         string `json:"lotdrop"`         //多系掉落
}

// 装备配置
type EquipConfig struct {
	EquipId         int     `json:"equip_id"`         // 装备eid
	EquipPosition   int     `json:"equip_position"`   // 1：武器 2：头盔 3：衣服 4：鞋子
	Quality         int     `json:"quality"`          // 品质
	EquipAttackType int     `json:"equip_attacktype"` //装备适用1~3
	EquipType       int     `json:"type"`             //装备适用1~6
	BaseTypes       []int   `json:"base_type"`        // 基础属性类型
	BaseValues      []int64 `json:"base_value"`       // 基础属性值
	CampExtAdd      int     `json:"camp_addition"`    //阵营加成
	Camp            int     `json:"camp"`             //阵营
	AdvanceId       int     `json:"advanceid"`        //进阶变化
	ShopClass       int     `json:"shopclass"`        //商店组
	ShopWeight      int     `json:"shopweight"`       //商店权重
	Price           []int   `json:"price"`            //价格
}

//装备升阶
type EquipAdvancedConfig struct {
	Quality       int   `json:"quality"`        // 装备eid
	IsAdvanced    int   `json:"isadvanced"`     // 1：武器 2：头盔 3：衣服 4：鞋子
	AdvancedNeed  []int `json:"advancedneed"`   // 需要物品
	AdvancedNum   []int `json:"advancednum"`    // 需要数量
	AdvancedValue []int `json:"advanced_value"` // 升阶属性
}

type EquipStrengthenConfig struct {
	Type          int   `json:"type"`           //装备适用
	EquipPosition int   `json:"equip_position"` // 1：武器 2：头盔 3：衣服 4：鞋子
	Quality       int   `json:"quality"`        // 品质
	Lv            int   `json:"lv"`             //等级
	Vaual         []int `json:"vaual"`          // 属性加成
}

type EquipStrengthenLvUpConfig struct {
	Quality     int `json:"quality"`         // 品质
	Lv          int `json:"lv"`              //等级
	UpLvNeedExp int `json:"strengthenlvexp"` //
	ExpBy       int `json:"strengthenexp"`   //自己提供的经验
}

type EquipRecastConfig struct {
	Id        int   `json:"id"`
	Quality   int   `json:"quality"` // 品质
	Attribute int   `json:"attribute"`
	CostTime  int   `json:"costtime"`
	CostNum   int   `json:"costnum"`
	Change    []int `json:"change"`
	Weight    []int `json:"weight"`
}

type EquipValueGroup struct {
	ValueGroup int   `json:"value_group"` //
	Type       int   `json:"type"`        //
	Group      []int `json:"group"`       //
	Chance     []int `json:"chance"`      //
}

type EquipBaseValue struct {
	Id         int     `json:"id"`         //
	Group      int     `json:"group"`      //
	Chance     int     `json:"chance"`     //
	BaseTypes  []int   `json:"base_type"`  // 基础属性类型
	BaseValues []int64 `json:"base_value"` // 基础属性值
}

type EquipSpecialValue struct {
	Id         int     `json:"id"`         //
	Group      int     `json:"group"`      //
	Chance     int     `json:"chance"`     //
	BaseTypes  []int   `json:"base_type"`  // 基础属性类型
	BaseValues []int64 `json:"base_value"` // 基础属性值
}

type EquipHireConfig struct {
	Num        int   `json:"num"`        //
	Type       int   `json:"type"`       //
	SubType    int   `json:"subtype"`    //
	Max        int   `json:"max"`        //
	Min        int   `json:"min"`        //
	Equip      []int `json:"equip"`      //
	Strengthen []int `json:"strengthen"` //
}

// 装备强化
type EquipUpgrade struct {
	Id         int   `json:"upgrad_id"`
	Lv         int   `json:"equip_lv"`
	CostIds    []int `json:"cost_type"`
	CostNums   []int `json:"cost_num"`
	RebornIds  []int `json:"reborn_item"`
	RebornNums []int `json:"reborn_num"`
}

// 装备附魔
type EquipStar struct {
	Id         int   `json:"star_id"`
	Lv         int   `json:"star_lv"`
	CostIds    []int `json:"cost_type"`
	CostNums   []int `json:"cost_num"`
	RebornItem map[int]*Item // 重生道具
}

// 宝石配置
type EquipGem struct {
	Id         int     `json:"gem_id"`
	Level      int     `json:"lv"`        // 等级
	NeedId     int     `json:"combo"`     // 合成当前宝石需要的宝石Id
	NeedNum    int     `json:"combo_num"` // 合成当前宝石需要的宝石数量
	GemType    int     `json:"gem_type"`  // 宝石类型
	BaseTypes  []int   `json:"base_type"`
	BaseValues []int64 `json:"base_value"`
}

// 装备套装配置
type EquipSuit struct {
	SuitId      int     `json:"suit_id"`       // 套装id
	SuitMark    []int   `json:"suit_mark"`     // 套装数量
	BaseTypes1  []int   `json:"reward1_type"`  // 1属性类型
	BaseValues1 []int64 `json:"reward1_value"` // 1属性数值
	BaseTypes2  []int   `json:"reward2_type"`  // 2属性类型
	BaseValues2 []int64 `json:"reward2_value"` // 2属性数值
	BaseTypes3  []int   `json:"reward3_type"`  // 3属性类型
	BaseValues3 []int64 `json:"reward3_value"` // 3属性数值
}

type ArtifactEquipConfig struct {
	ArtifactId    int     `json:"artifact_id"`    // id
	BaseTypes     []int   `json:"base_type"`      // 属性类型
	BaseValues    []int64 `json:"base_value"`     // 属性数值
	ArtifactSkill []int64 `json:"artifact_skill"` //技能
}

type ArtifactStrengthenConfig struct {
	Id    int     `json:"artifact_id"`      // 套装id
	Lv    int     `json:"strengthen_lv"`    // 属性类型
	Value []int64 `json:"strengthen_value"` // 属性数值这个要对应神器表中的索引
	Need  []int   `json:"strengthen_need"`  // 强化所需道具
	Num   []int   `json:"strengthen_num"`   // 数量
	Skill []int   `json:"strengthen_skill"` // 技能
}

type ExclusiveEquipConfig struct {
	Id         int     `json:"id"`              //Id
	HeroId     int     `json:"relevanceheroid"` //
	BaseType   []int   `json:"base_type"`       //属性类型
	BaseValue  []int64 `json:"base_value"`      //属性数值
	ActiveNeed int     `json:"activeneed"`      //激活需要
	ActiveNum  int     `json:"activenum"`       //激活所需
	Skill      int     `json:"skill"`           //
}

type ExclusiveStrengthenConfig struct {
	Id       int     `json:"id"`             // 套装id
	Lv       int     `json:"lv"`             // 属性类型
	Value    []int64 `json:"base_value"`     // 属性数值这个要对应神器表中的索引
	Need     int     `json:"strengthenneed"` // 强化所需道具
	Num      int     `json:"strengthennum"`  // 数量
	Skill    int     `json:"skillid"`        // 技能
	Replace  int     `json:"replace"`        // 替换ID
	BackItem []int   `json:"back_item"`      // 回退返还
	BackNum  []int   `json:"back_num"`       // 回退返还
}

// 掉落组配置
type ItemBagGroup struct {
	Id      int   `json:"itembag"` // 掉落组Id
	Type    int   `json:"type"`    // 掉落类型
	DropIds []int `json:"item"`    // 掉落包Id
	Weights []int `json:"weight"`  // 掉落包权重
	Sum     int                    // 总权重
}

// 掉落包配置
type ItemBag struct {
	Id      int   `json:"itembag"` // 掉落包Id
	ItemIds []int `json:"item"`    // 物品Id
	Nums    []int `json:"num"`     // 物品数量
	Weights []int `json:"weight"`  // 权重
	Sum     int                    // 总权重
}

// 英雄分解
type HeroDecompose struct {
	ItemId        int `json:"item_id"`        // 英雄碎片Id
	DecomposeItem int `json:"decompose_item"` // 分解后的道具Id
	DecomposeNum  int `json:"decompose_num"`  // 分解后的数量
}

type LevelConfig struct {
	LevelId            int    `json:"levelid"`
	MapId              int    `json:"mapid"`
	MainType           int    `json:"maintype"`
	LevelType          int    `json:"leveltype"`
	LevelIndex         int    `json:"levelindex"`
	ChapterIndex       int    `json:"chapterindex"`
	TaskTypes          int    `json:"tasktypes"`
	TaskConds          []int  `json:"n"`
	ALevel             []int  `json:"alevel"` // 前置关卡
	EverydayNum        int    `json:"everydaynum"`
	PhysicalStrength   int    `json:"physicalstrength"`
	TeamExp            int    `json:"teamexp"`
	NextLevel          int    `json:"nextlevel"` // 下一关卡
	LevelName          string `json:"levelname"`
	Chaptername        string `json:"chaptername"`
	LevelGroup         int    `json:"levelgroup"`
	GoldYield          int    `json:"gold_yield"`
	ScienceYield       int    `json:"science_yield"`
	Gjtime             []int  `json:"gj_time"`
	GjDrop             []int  `json:"gj_drop"`
	GjScore            int    `json:"gj_score"`
	HangUp             int    `json:"hang_up"`
	TaskIndex          int    `json:"taskindex"`
	Comat              int64  `json:"comat"`
	DungeonExtraDrop   []int  `json:"dungeon_extra_drop"`
	DailyRechargeGroup int    `json:"daily_recharge_group"`
	LevelSkip          int    `json:"level_skip"`
	SkipType           int    `json:"skip_type"`
	SkipNum            int64  `json:"skip_num"`
}

// 所有消耗
type CostConfig struct {
	Id       int   `json:"id"`
	ItemIds  []int `json:"costid"`    // 消耗Id
	ItemNums []int `json:"costcount"` // 消耗数量
}

// 科技消耗
type TechConfig struct {
	Id         int   `json:"id"`            // 唯一Id
	ProtectIds []int `json:"protect_id"`    // 前置科技ID
	Conds      []int `json:"condition"`     // 对应的前置科技等级
	TechGroup  int   `json:"science_group"` // 科技的升级组, 对应TechAttr的gourp
}

// 科技属性
type TechAttr struct {
	Id           int     `json:"id"`             // 唯一Id
	Group        int     `json:"gourp"`          // 对应TechConfig的science_group
	Level        int     `json:"level"`          // 科技等级
	AttType      []int   `json:"attribute_type"` // 属性类型,全队
	AttValue     []int64 `json:"value"`          // 属性值
	PlayerLv     int     `json:"player_level"`   // 需要的玩家等级
	CompleteTime int     `json:"complete_time"`  // 每次的时间
	Costid       int     `json:"cost_id"`        // 消耗Id
	Costnum      int     `json:"cost_num"`       // 消耗数量
}

// 科技时间消耗
type TechTime struct {
	Id         int   `json:"id"`         // 消耗Id
	TimeLimit  int   `json:"time_limit"` // 总时间
	CdTime     int   `json:"hourly"`     // 每次秒多长时间
	Costidids  []int `json:"costid_id"`
	Costidnums []int `json:"costid_num"`
}

// 新手道具
type NewUserItem struct {
	Sequence int `json:"sequence"`
	Group    int `json:"group"`
	Id       int `json:"id"`
	Type     int `json:"type"`
	Num      int `json:"num"`
	Place    int `json:"place"`
}

type SoulchestConfig struct {
	Sequence     int   `json:"sequence"`
	Bag          int   `json:"bag"`
	Probabilitys []int `json:"probability"`
	Item         int   `json:"item"`
	Num          int   `json:"num"`
	Starttime    int   `json:"starttime"`
	Endtime      int   `json:"endtime"`
}

type LegioncopyConfig struct {
	Id    int `json:"id"`
	Level int `json:"leve"`
	//Text     int   `json:"text"`
	Ns       []int `json:"n"`
	Xs       []int `json:"x"`
	Monsters []int `json:"monster"`
	Hps      []int `json:"hp1"`
	Atts     []int `json:"att"`
	Items    []int `json:"item"`
}
type FundConfig struct {
	Fundid    int   `json:"fundid"`
	Paging    int   `json:"paging"`
	Awards    []int `json:"award1"`
	Nums      []int `json:"num"`
	Level     int   `json:"level"`
	Buy       int   `json:"buy"`
	Viplevel  int   `json:"viplevel"`
	Fundlevel int   `json:"fundlevel"`
}

type ActivityFundConfig struct {
	GroupID int   `json:"groupid"`
	Type    int   `json:"type"`
	ID      int   `json:"id"`
	Pay     int   `json:"pay"`
	Worth   int   `json:"worth"`
	Day     int   `json:"day"`
	Items   []int `json:"item"`
	Nums    []int `json:"num"`
}

type ActivityFundGroupMap struct {
	GroupID   int `json:"groupid"`
	PayConfig map[int]*ActivityFundPayMap
}

type ActivityFundPayMap struct {
	Pay       int `json:"pay"`
	Type      int `json:"type"`
	Worth     int `json:"worth"`
	DayConfig map[int]*ActivityFundConfig
}

type ActivityFundTypeConfig struct {
	Type int `json:"type"`
	Pay  int `json:"pay"`
}

type ShoprefreshConfig struct {
	Shoptypy    int `json:"shoptypy"`
	Cost        int `json:"cost"`
	Currency    int `json:"currency"`
	Refreshitem int `json:"refresh_item"`
	Refreshcost int `json:"refresh_cost"`
}
type HorseParmConfig struct {
	Id     int    `json:"id"`
	Parms  string `json:"parm"`
	System string `json:"system"`
}
type HolyLegendLevelConfig struct {
	Group      int    `json:"group"`
	Levelindex int    `json:"level_index"`
	Level      int    `json:"level"`
	Items      []int  `json:"item"`
	Nums       []int  `json:"num"`
	LevelName  string `json:"level_name"`
}
type WarcontributionConfig struct {
	Id            int   `json:"id"`
	Min           int   `json:"min"`
	Max           int   `json:"max"`
	Contributions []int `json:"contribution"`
	Boxs          []int `json:"box"`
}

type LevelboxConfig struct {
	Id    int   `json:"id"`
	Items []int `json:"item"`
	Nums  []int `json:"num"`
}

type LevelMapConfig struct {
	RegionId   int   `json:"regionid"`
	JNeedStar  []int `json:"jneedstar"`
	NeedLevel  int   `json:"needlevel"`
	Prop       []int `json:"prop"`
	PropNumber []int `json:"propnumber"`
}

type GemsweeperrankConfig struct {
	Rank    int   `json:"rank"`
	Itemids []int `json:"item_id"`
	Nums    []int `json:"num"`
}
type GemsweeperextradropConfig struct {
	Id             int   `json:"id"`
	Group          int   `json:"group"`
	Week           int   `json:"week"`
	Rewardtimes    []int `json:"reward_time"`
	Overtimes      []int `json:"over_time"`
	Extradrops     []int `json:"extra_drop"`
	Extradropshows []int `json:"extra_drop1_show"`
}

type HorseSoulUpgradeConfig struct {
	Id                int   `json:"id"`
	Upgrademodelgroup int   `json:"upgrade_model_group"`
	Level             int   `json:"level"`
	Costitem          int   `json:"cost_item"`
	Costnum           int   `json:"cost_num"`
	Decomposeitem     int   `json:"decompose_item"`
	Decomposenum      int   `json:"decompose_num"`
	Value             int64 `json:"value"`
	ValueOne          int64 `json:"valueone"`
	AttType1          int   `json:"attribute_type1" trim:"0"`
	AttValue1         int64 `json:"attribute_value1" trim:"0"`
}

type ConsumetophpConfig struct {
	Id     int   `json:"id"`
	Heroid int   `json:"heroid"`
	Mapid  int   `json:"mapid"`
	Level  int   `json:"level"`
	Hp     int   `json:"hp"`
	Items  []int `json:"item"`
	Nums   []int `json:"num"`
}
type ConsumetopbossConfig struct {
	Id      int    `json:"id"`
	Picture string `json:"picture"`
	Boss    int    `json:"boss"`
	Name    string `json:"name"`
	Mapid   int    `json:"mapid"`
	List    int    `json:"list"`
	Shop    int    `json:"shop"`
}

type FitConfig struct {
	Server int   `json:"server"`
	Fit    int   `json:"fit"`
	Group  int   `json:"group"`
	Items  []int `json:"item"`
	Nums   []int `json:"num"`
}

type ConsumetoplistConfig struct {
	Id    int    `json:"id"`
	Sort  int    `json:"sort"`
	Type  int    `json:"type"`
	Ns    []int  `json:"n"`
	Text  string `json:"txt"`
	Items []int  `json:"item"`
	Nums  []int  `json:"num"`
	Group int    `json:"group"`
}
type SignConfig struct {
	Id     int `json:"id"`
	Month  int `json:"month"`
	Sign   int `json:"sign"`
	Reward int `json:"reward"`
	Number int `json:"number"`
	Vip    int `json:"vip"`
	Reset  int `json:"reset"`
}
type HalfmoonConfig struct {
	Id        int   `json:"id"`
	Sort      int   `json:"sort"`
	Step      int   `json:"step"`
	Tasktypes int   `json:"tasktypes"`
	Ns        []int `json:"n"`
	Items     []int `json:"item"`
	Nums      []int `json:"num"`
	Costitems []int `json:"costitem"`
	Costnums  []int `json:"costnum"`
	Cost      int   `json:"cost"`
	Sale      int   `json:"sale"`
	Fund      int   `json:"fund"`
	Power     int   `json:"power"`
	Start     int   `json:"start"`
	Continued int   `json:"continued"`
	Cd        int   `json:"cd"`
	Show      int   `json:"show"`
	Coctime   int   `json:"coctime"`
	Renovate  int   `json:"renovate"`
	Reset     int   `json:"reset"`
	Status    int   `json:"status"`
}
type PeoplecityConfig struct {
	Id int `json:"id"`
	//Name         int   `json:"name"`
	Levels       []int `json:"level"`
	Boxs         []int `json:"box"`
	Countrys     []int `json:"country"`
	Shuitems     []int `json:"shuitem"`
	Shunums      []int `json:"shunum"`
	Prestiges    []int `json:"prestige"`
	Prestigenums []int `json:"prestigenum"`
	Weiitems     []int `json:"weiitem"`
	Weinums      []int `json:"weinum"`
	Wuitems      []int `json:"wuitem"`
	Wunums       []int `json:"wunum"`
	Pngs         []int `json:"png"`
	Itemids      []int `json:"itemid"`
	Specials     []int `json:"special"`
	Chances      []int `json:"chance"`
}
type TitleConfig struct {
	Id         int   `json:"id"`
	Conditions []int `json:"condition"`
	Ns         []int `json:"n1"`
	Methodss   []int `json:"methods"`
	Ys         []int `json:"y"`
	Zs         []int `json:"z"`
}
type SignrewardConfig struct {
	Id          int   `json:"id"`
	Signnum     int   `json:"signnum"`
	Rewarditems []int `json:"rewarditem"`
	Rewardnums  []int `json:"rewardnum"`
}
type ExpeditionbuffConfig struct {
	Id         int `json:"id"`
	Buffidtype int `json:"buffidtype"`
	Type       int `json:"type"`
	Typevalue  int `json:"typevalue"`
	Value      int `json:"value"`
	Itemconfig int `json:"itemconfig"`
	Price      int `json:"price"`
}
type HomeofficerankingConfig struct {
	Type      int    `json:"type"`
	Rank      int    `json:"rank"`
	Itemids   []int  `json:"item_id"`
	Nums      []int  `json:"num"`
	Mailtitle string `json:"mail_title"`
	Mailtxt   string `json:"mail_txt"`
}

type NationalwarawardConfig struct {
	Index int   `json:"index"`
	Items []int `json:"item"`
	Nums  []int `json:"num"`
}
type HolyLegendConfig struct {
	Beautyid         int    `json:"beauty_id"`
	Chapter          int    `json:"chapter"`
	Chaptercondition int    `json:"chapter_condition"`
	Legendname       string `json:"legend_name"`
	Chaptername      string `json:"chapter_name"`
	Levelgroup       int    `json:"level_group"`
}
type HonorkillConfig struct {
	Honorid      int `json:"honorid"`
	Honortype    int `json:"honortype"`
	Honorkill    int `json:"honorkill"`
	Honoritem    int `json:"honoritem"`
	Honoritemnum int `json:"honoritemnum"`
	Maxhonor     int `json:"maxhonor"`
	Minhonor     int `json:"minhonor"`
	Killbonuses  int `json:"killbonuses"`
}
type GemsweeperawardConfig struct {
	Id    int   `json:"id"`
	Items []int `json:"item"`
	Nums  []int `json:"num"`
}
type WartargetConfig struct {
	Id           int   `json:"id"`
	Step         int   `json:"step"`
	Tasktypes    int   `json:"tasktypes"`
	Ns           []int `json:"n"`
	Contribution int   `json:"contribution"`
	Items        []int `json:"item"`
	Nums         []int `json:"num"`
}

type WarencourageConfig struct {
	Id        int `json:"id"`
	Min       int `json:"min"`
	Max       int `json:"max"`
	Crtchance int `json:"crtchance"`
	Crtnum    int `json:"crtnum"`
	Encpoints int `json:"encpoints"`
}

type DiplomacychangesConfig struct {
	ID        int   `json:"ID"`
	Forcetype int   `json:"Force_type"`
	Index     int   `json:"index"`
	Type      int   `json:"type"`
	Parms     []int `json:"parm"`
	Forceids  []int `json:"Force_id"`
	Storys    []int `json:"story"`
	Mystory   int   `json:"my_story"`
	Juqingid  int   `json:"juqingid"`
}
type ResourcelevelConfig struct {
	Group        int   `json:"group"`
	Index        int   `json:"index"`
	LevelId      int   `json:"level_id"`
	Levellimits  int   `json:"level_limits"`
	Items        []int `json:"first_item"`
	Nums         []int `json:"num"`
	Itembaggroup int   `json:"itembag_group"`
	Costids      []int `json:"cost_id"`
	Costnums     []int `json:"cost_num"`
}
type BuyphysicalConfig struct {
	Number   int `json:"number"`
	Money    int `json:"money"`
	Physical int `json:"physical"`
}
type ChestshopConfig struct {
	Grid     int `json:"grid"`
	Itemid   int `json:"itemid"`
	Num      int `json:"num"`
	Currency int `json:"currency"`
	Type     int `json:"type"`
}
type CommunityConfig struct {
	Level            int `json:"lv"`
	Exp              int `json:"exp"`
	Changeexp        int `json:"changeexp"`
	Membernum        int `json:"population"`
	Lively           int `json:"lively"`
	Changelively     int `json:"changelively"`
	Elder            int `json:"elder"`
	Fearless         int `json:"fearless"`
	Activelimit      int `json:"activelimit"`
	Warfare          int `json:"warfare"`
	GuildActiveLimit int `json:"guildactivelimit"`
	GuildExpLimit    int `json:"guildexplimit"`
}
type CrownbuildConfig struct {
	Type    int   `json:"type"`
	Rewards []int `json:"reward"`
	Items   []int `json:"item"`
	Open    int   `json:"open"`
}

//! 世界等级配置
type WorldLevelConfig struct {
	Id         int   `json:"id"`
	WorldLv    int   `json:"worldlv"`
	NpcLv      int   `json:"npclv"`
	Ragerate   int   `json:"ragerate"`
	BaseTypes  []int `json:"base_type"`
	BaseValues []int `json:"base_value"`
}
type WorldlvtpyeConfig struct {
	Id                    int    `json:"id"`
	Npctype               int    `json:"npctype"`
	Teamshow              int    `json:"teamshow"`
	TeamName              string `json:"teamname"`
	Typesub               int    `json:"typesub"`
	Typeweight            int    `json:"typeweight"`
	Timesprotect          int    `json:"timesprotect"`
	Npclvtable            int    `json:"npclvtable"`
	Npcteam               int    `json:"npcteam"`
	Npcicon               int    `json:"npcicon"`
	Npcmods               []int  `json:"npcmod"`
	Heronpcs              []int  `json:"heronpc"`
	Armmods               []int  `json:"armmod"`
	Hpcorrect             int    `json:"hpcorrect"`
	Attackcorrect         int    `json:"attackcorrect"`
	Armorcorrect          int    `json:"armorcorrect"`
	Magiccorrect          int    `json:"magiccorrect"`
	Rescorrect            int    `json:"rescorrect"`
	Attackspeedcorrect    int    `json:"attackspeedcorrect"`
	Movespeedcorrect      int    `json:"movespeedcorrect"`
	Hitcorrect            int    `json:"hitcorrect"`
	Dodgecorrect          int    `json:"dodgecorrect"`
	Blockcorrect          int    `json:"blockcorrect"`
	Phycritcorrect        int    `json:"phycritcorrect"`
	Magiccritcorrect      int    `json:"magiccritcorrect"`
	Phycrittimescorrect   int    `json:"phycrittimescorrect"`
	Armorbreakcorrect     int    `json:"armorbreakcorrect"`
	Magicbreakcorrect     int    `json:"magicbreakcorrect"`
	Vampirecorrect        int    `json:"vampirecorrect"`
	Healcorrect           int    `json:"healcorrect"`
	Rageratecorrect       int    `json:"rageratecorrect"`
	Cablerangecorrect     int    `json:"cablerangecorrect"`
	Initialangercorrect   int    `json:"initialangercorrect"`
	Physicaldamagecorrect int    `json:"physicaldamagecorrect"`
	Magichurtcorrect      int    `json:"magichurtcorrect"`
	Immunitydamagecorrect int    `json:"immunitydamagecorrect"`
	Exp                   int    `json:"exp"`
}

type CamprewardConfig struct {
	Id            int   `json:"id"`
	Step          int   `json:"step"`
	Personaltype  int   `json:"personaltype"`
	Personaljudge int   `json:"personaljudge"`
	Personalnum   int   `json:"personalnum"`
	Items         []int `json:"item"`
	Nums          []int `json:"num"`
}

type BeautychestConfig struct {
	Sequence     int   `json:"sequence"`
	Bag          int   `json:"bag"`
	Probabilitys []int `json:"probability"`
	Item         int   `json:"item"`
	Num          int   `json:"num"`
}

type GemchestConfig struct {
	Sequence     int   `json:"sequence"`
	Bag          int   `json:"bag"`
	Probabilitys []int `json:"probability"`
	Item         int   `json:"item"`
	Num          int   `json:"num"`
}
type TreasuryConfig struct {
	Id        int `json:"id"`
	ShopType  int `json:"shoptype"`
	Sort      int `json:"sort"`
	Type      int `json:"type"` // 1不限购,2每日限购,3每周限购
	Time      int `json:"time"`
	Item      int `json:"item"`
	Currency  int `json:"currency"`
	Num       int `json:"num"`
	Cost      int `json:"cost"`
	Costnum   int `json:"costnum"`
	PataLimit int `json:"pata_limit"`
}
type CitywayConfig struct {
	Id    int   `json:"id"`
	Citys []int `json:"city"`
}
type BuildingingConfig struct {
	Id     int   `json:"id"`
	Type   int   `json:"type"`
	Lv     int   `json:"lv"`
	Need   int   `json:"need"`
	Fight  int   `json:"fight"`
	Battle int   `json:"battle"`
	Items  []int `json:"item"`
	Nums   []int `json:"num"`
}

// vip配置表替换
type VipConfig struct {
	Viplevel            int   `json:"viplevel"`
	Needexp             int   `json:"need_exp"`
	Totleexp            int   `json:"totle_exp"`
	Copperfreeaccess    int   `json:"Copper_freeaccess"`
	Copperbuyaccess     int   `json:"Copper_buyaccess"`
	Physicalfreeaccess  int   `json:"physical_freeaccess"`
	Physicalbuyaccess   int   `json:"physical_buyaccess"`
	Proactivefreeaccess int   `json:"proactive_freeaccess"`
	Proactivebuyaccess  int   `json:"proactive_buyaccess"`
	Infantryfreeaccess  int   `json:"Infantry_freeaccess"`
	Infantrybuyaccess   int   `json:"Infantry_buyaccess"`
	Elite               int   `json:"elite" trim:"0"`
	Shoprefreshs        []int `json:"shoprefresh"`
	GuildHunting        []int `json:"guild_hunting"`
	GuildSweep          int   `json:"guild_sweep"`
	Visit               int   `json:"visit" trim:"0"` // 有两个字段
	Jjcs                []int `json:"jjc"`
	People              int   `json:"people"`
	Visittime           int   `json:"visittime"`
	Jjctime             int   `json:"jjctime"`
	Sweep               int   `json:"sweep" trim:"0"` // 有两个一样的字段
	Beautyclick         int   `json:"beautyclick"`
	Meethero            int   `json:"meethero"`
	Physical            int   `json:"physical"`
	Alchemys            []int `json:"alchemy"`
	Archerfreeaccess    int   `json:"Archer_freeaccess"`
	Archerbuyaccess     int   `json:"Archer_buyaccess"`
	Cavalryfreeaccess   int   `json:"cavalry_freeaccess"`
	Cavalrybuyaccess    int   `json:"cavalry_buyaccess"`
	Honorfreeaccess     int   `json:"honor_freeaccess"`
	Honorbuyaccess      int   `json:"honor_buyaccess"`
	Skillnumber         int   `json:"skillnumber"`
	Buyskill            int   `json:"buyskill"`
	Continuousalchemy   int   `json:"continuousalchemy"`
	Profiteer           int   `json:"profiteer"`
	Touristtrap         int   `json:"touristtrap"`
	Beauty              int   `json:"beauty"`
	Taxnum              int   `json:"taxnum"`
	Resourcechallenge   int   `json:"resource_challenge"`
	Herosent            int   `json:"hero_sent"`
	Horsehigtcall       int   `json:"horse_higt_call"`
	TimeGeneralsnum     int   `json:"time_generals_num"`
	Consumetopnum       int   `json:"consumetop_num"`
	TreasureClear       int   `json:"treasure_clear"`
	GemsMopping         int   `json:"gems_mopping"`
	GemsBuy             int   `json:"gems_buy"`
	JJcBuy              []int `json:"jjc_buy"`
	Beforenum           int   `json:"beforenum"`
	Nownum              int   `json:"nownum"`
	Items               []int `json:"item"`
	Nums                []int `json:"num"`
	ExtendGroup         int   `json:"extend_group"`
	Growthtask_King1    int   `json:"Growthtask_King1" trim:"0"`
	Growthtask_King2    int   `json:"Growthtask_King2" trim:"0"`
	Pata_buy1           int   `json:"pata_buy1" trim:"0"` // 爬塔重置次数
	Pata_buy2           int   `json:"pata_buy2" trim:"0"` // 爬塔精英副本次数
	Pata_buy3           int   `json:"pata_buy3" trim:"0"` // 爬塔buff刷新次数
	CrownFight          int   `json:"Crown_Fight"`
	SummonDiscount      int   `json:"summon_discount"` // 打折100
	SummonTimes         int   `json:"summon_times"`    // 打折次数
	ArmyBuy             int   `json:"Army_buy"`        // 佣兵购买次数
	DungeonReset        int   `json:"Dungeons_Reset"`  // 地下城重置次数
	RebornTimes         int   `json:"hero_reborn"`     //英雄重生次数
	GjLimit             int   `json:"gj_limit"`        //挂机时长限制
	GjGrow              int   `json:"gj_grow"`         //挂机加成
	DungeonsSweep       int   `json:"dungeons_sweep"`  //地下城扫荡
	FreeItems           []int `json:"free_item"`       //VIP每日奖励
	FreeNums            []int `json:"free_num"`
	WeekPrice           int   `json:"week_price"`      //价格
	WeekLimit           int   `json:"week_limit"`      //购买数量限制
	WeekItems           []int `json:"week_item"`       //VIP每周福利
	WeekNums            []int `json:"week_num"`        //VIP每周福利
	ArenaFree           []int `json:"arena_free"`      // 竞技场免费次数
	RewardItask         int   `json:"reward_itask"`    // 悬赏个人任务个数
	RewardTeamtask      int   `json:"reward_teamtask"` // 悬赏个人任务个数
	RewardOnekey        int   `json:"reward_onekey"`   // 悬赏个人一键功能
	Astrologer          int   `json:"astrologer"`      // 占星开启

	HangupGold    int `json:"hangup_gold"`    //挂机金币加成
	HangupHeroExp int `json:"hangup_heroexp"` //挂机经验加成
	HangupFast    int `json:"hangup_fast"`    //快速次数
	HangupTime    int `json:"hangup_time"`    //最大时长
	MazeFateGold  int `json:"maze_fate_gold"` //迷宫命运金币产出增加万分比
	MazeGold      int `json:"maze_gold"`      //迷宫金币产出增加万分比
	HeroList      int `json:"hero_list"`      //英雄上限特权增加
	CallOptional  int `json:"call_optional"`  //自选池特权次数
	BossReset     int `json:"boss_reset"`     //世界BOSS重置次数
	//每日礼包相关
	DailyTimes int   `json:"daily_times"` //
	TaskTypes  int   `json:"tasktypes"`   //
	Ns         []int `json:"n"`           //
	DailyItem  []int `json:"daily_item"`  //
	DailyNum   []int `json:"daily_num"`   //
}

type SpyrewardConfig struct {
	Id      int `json:"id"`
	Eventid int `json:"event_id"`
	Office  int `json:"office"`
	Spyid   int `json:"spyid"`
	Chance  int `json:"chance"`
	Drop    int `json:"drop"`
	Notice  int `json:"notice"`
}
type PromoteboxConfig struct {
	Boxid     int `json:"boxid"`
	Teamlevel int `json:"teamlevel"`
	Item      int `json:"item"`
	Num       int `json:"num"`
}
type IndustryConfig struct {
	City      int   `json:"city"`
	Png       int   `json:"png"`
	Map       int   `json:"map"`
	Prestiges []int `json:"prestige"`
	Items     []int `json:"item1"`
	Nums      []int `json:"num1"`
	Builds    []int `json:"build"`
	Xs        []int `json:"x"`
	Ys        []int `json:"y"`
	Revs      []int `json:"rev"`
}
type WarcityConfig struct {
	Cityid       int   `json:"cityid"`
	Size         int   `json:"size"`
	Npcoccupy    int   `json:"npcoccupy"`
	Occupys      []int `json:"occupy"`
	Sequences    []int `json:"sequence"`
	Buffs        []int `json:"buff"`
	Debuffs      []int `json:"debuff"`
	Npcpalys     []int `json:"npcpaly"`
	Attackplays  []int `json:"attackplay"`
	Defenseplays []int `json:"defenseplay"`
	Strongholds  []int `json:"stronghold"`
	Influences   []int `json:"influence"`
}
type MaincityConfig struct {
	ID            int    `json:"ID"`
	Cityname      string `json:"cityname"`
	Mainlevel     int    `json:"mainlevel"`
	Consumption   int    `json:"consumption"`
	Number        int    `json:"number"`
	Openbuildings []int  `json:"openbuilding"`
}
type VisitConfig struct {
	Id        int   `json:"id"`
	Hero      int   `json:"hero"`
	Judgement int   `json:"judgement"`
	Needlv    int   `json:"needlv"`
	Levelids  []int `json:"levelid"`
	Needs     []int `json:"need"`
	Times     []int `json:"time"`
	Stars     []int `json:"star"`
	Boxs      []int `json:"box"`
}
type EquipchestConfig struct {
	Sequence     int   `json:"sequence"`
	Bag          int   `json:"bag"`
	Probabilitys []int `json:"probability"`
	Item         int   `json:"item"`
	Num          int   `json:"num"`
}
type TimeGeneralsRankawardConfig struct {
	Id           int      `json:"id"`
	Group        int      `json:"group"`
	Rankmin      int      `json:"rank_min"`
	Rankmax      int      `json:"rank_max"`
	Normalawards []int    `json:"normal_award"`
	Normalnums   []int    `json:"normal_num"`
	Needpoint    int      `json:"need_point"`
	Extraawards  []int    `json:"extra_award"`
	Extranums    []int    `json:"extra_num"`
	MailTitle    string   `json:"mail_title"`
	MainTxt      []string `json:"mail_txt"`
}
type ExpeditionbuffgroupConfig struct {
	Buffgroup int   `json:"buffgroup"`
	Ids       []int `json:"buffid"`
	Weights   []int `json:"weight"`
}
type ExcitingConfig struct {
	Power          int    `json:"power"`
	Id             int    `json:"id"`
	Icon           int    `json:"icon"`
	Type           int    `json:"type"`
	Targets        []int  `json:"target"`
	Countrys       []int  `json:"country"`
	Items          []int  `json:"item"`
	Nums           []int  `json:"num"`
	Ranking        int    `json:"ranking"`
	Rankingcontent string `json:"rankingcontent"`
	Awards         []int  `json:"award"`
	Fronts         []int  `json:"front"`
	Openicon       int    `json:"openicon"`
	Openname       int    `json:"openname"`
	Opendec        int    `json:"opendec"`
}
type HomeofficecityawardConfig struct {
	Counttype     int `json:"count_type"`
	Index         int `json:"index"`
	Citymin       int `json:"city_min"`
	Citymax       int `json:"city_max"`
	Extraitemid   int `json:"extra_itemid"`
	Singlecitynum int `json:"single_city_num"`
	Sumcitynum    int `json:"sum_city_num"`
}
type HorseBattleSteedAwakenConfig struct {
	ID              int   `json:"ID"`
	Battlesteedid   int   `json:"battle_steed_id"`
	Awakenlevel     int   `json:"awaken_level"`
	Quality         int   `json:"quality"`
	Costhorsestar   int   `json:"cost_horse_star"`
	Costhorsenum    int   `json:"cost_horse_num"`
	Costitem        int   `json:"cost_item"`
	Costnum         int   `json:"cost_num"`
	Attributetypes  []int `json:"attribute_type"`
	Attributevalves []int `json:"attribute_valve"`
	Dragonholenum   int   `json:"dragon_hole_num"`
}
type TeamexpConfig struct {
	Teamlv              int   `json:"teamlv"`
	Teamexplv           int   `json:"teamexplv"`
	Teamtotolexp        int   `json:"teamtotolexp"`
	Getphysical         int   `json:"getphysical"`
	Getlimit            int   `json:"getlimit"`
	Physicallimit       int   `json:"physicallimit"`
	Recommendedstrength int   `json:"Recommendedstrength"`
	Powerlimit          int   `json:"powerlimit"`
	Items               []int `json:"item"`
	Nums                []int `json:"num"`
}
type TimeGeneralsConfig struct {
	KeyId             int      `json:"keyid"`
	Id                int      `json:"id"`
	ActType           int      `json:"type"`
	Activityid        int      `json:"activity_id"`
	CallSingleType    int      `json:"call_single_type"`
	CallSinglePoint   int      `json:"call_single_points"`
	CallTenType       int      `json:"call_ten_type"`
	CallTenPoint      int      `json:"call_ten_points"`
	CostSingleitem    int      `json:"cost_single_item"`
	CostSingleNum     int      `json:"cost_single_num"`
	Costtenitem       int      `json:"cost_ten_item"`
	CostTenNum        int      `json:"cost_ten_num"`
	RewardPointsGroup int      `json:"reward_points_group"`
	RankAwardGroup    int      `json:"rank_award_group"`
	Mainheroids       []int    `json:"main_hero_id"`
	HeroIds           []int    `json:"hero_id"`
	CallDesc          string   `json:"call_describe"`
	NewHero           string   `json:"new_hero"`
	MainHeroLocation  []string `json:"mainhero_location"`
	HeroLocation      []string `json:"hero_location"`
}

type PubchesttotalConfig struct {
	Pubtpye              int   `json:"pub_tpye"`
	Cardsvipmin          int   `json:"cards_vip_min"`
	Cardsvipmax          int   `json:"cards_vip_max"`
	Paypubtype           int   `json:"pay_pub_type"`
	Payitem              int   `json:"pay_item"`
	Payitemnum           int   `json:"pay_item_num"`
	Dropcardsnum         int   `json:"drop_cards_num"`
	Dropgroups           []int `json:"drop_group"`
	Dropgroupids         []int `json:"drop_group_id"`
	Chance               int   `json:"chance"`
	Certaintimesdroptype int   `json:"certain_times_drop_type"`
	AddItemType          int   `json:"additem_type"`
	AddItemNum           int   `json:"additem_num"`
}

type TimeResetConfig struct {
	Id       int     `json:"id"`
	System   int     `json:"system"`
	TimeType int     `json:"timetype"`
	Continue int64   `json:"continue"`
	Cd       int64   `json:"cd"`
	Time     []int64 `json:"time"`
}

type WarOrderConfig struct {
	Id        int   `json:"id"`
	Type      int   `json:"type"`
	NeedPoint int   `json:"need_point"`
	FreeAward []int `json:"free_award"`
	FreeNum   []int `json:"free_num"`
	GoldAward []int `json:"gold_award"`
	GoldNum   []int `json:"gold_num"`
}

type WarOrderParam struct {
	Id         int `json:"id"`
	BuyItem    int `json:"buy_item"`
	BuyNum     int `json:"buy_num"`
	SwitchItem int `json:"switch_item"`
	SwitchNum  int `json:"switch_num"`
}

type WarOrderLimitConfig struct {
	Id        int   `json:"id"`
	Type      int   `json:"type"`
	Group     int   `json:"group"`
	TaskTypes int   `json:"tasktypes"`
	Ns        []int `json:"n"`
	FreeAward []int `json:"free_award"`
	FreeNum   []int `json:"free_num"`
	GoldAward []int `json:"gold_award"`
	GoldNum   []int `json:"gold_num"`
}

type AccessAwardConfig struct {
	Id             int      `json:"id"`
	Group          int      `json:"group"`
	Point          int      `json:"point"`
	Countdown      int64    `json:"countdown"`
	Item           []int    `json:"item"`
	Num            []int    `json:"num"`
	ConversionItem []int    `json:"conversion_item"`
	ConversionNum  []int    `json:"conversion_num"`
	Notice         []string `json:"notice"`
}

type AccessRankConfig struct {
	Id           int      `json:"id"`
	Group        int      `json:"group"`
	RankMin      int      `json:"rank_min"`      //最小名次
	RankMax      int      `json:"rank_max"`      //最大名次
	NormalAward  []int    `json:"normal_award"`  //基础奖励
	NormalNum    []int    `json:"normal_num"`    //基础奖励数量
	NeedPoint    int      `json:"need_point"`    //额外奖励积分
	ExtraAward   []int    `json:"extra_award"`   //额外奖励
	ExtraNum     []int    `json:"extra_num"`     //额外奖励数量
	MailTitle    string   `json:"mail_title"`    //邮件标题
	MailTxt      []string `json:"mail_txt"`      //内容
	RankOverTime int64    `json:"rank_overtime"` //在展示期第几天结算排名奖励
}

type AccessTaskConfig struct {
	Id        int   `json:"id"`
	Group     int   `json:"group"`
	TaskTypes int   `json:"tasktypes"`
	Ns        []int `json:"n"`
	Item      []int `json:"item"`
	Num       []int `json:"num"`
	Notice    []int `json:"notice"`
}

type ActivityDailyRecharge struct {
	Id   int   `json:"id"`
	Item []int `json:"item"`
	Num  []int `json:"num"`
}

type ActivityOverflowGifts struct {
	Id    int   `json:"id"`
	Group int   `json:"group"`
	Item  []int `json:"item"`
	Num   []int `json:"num"`
}

//type ActivityboxConfig struct {
//	Boxid       int   `json:"boxid"`
//	Gearid      int   `json:"gearid"`
//	Group       int   `json:"group"`
//	Sort        int   `json:"sort"`
//	Backicon    int   `json:"backicon"`
//	Type        int   `json:"type"`
//	Look        int   `json:"look"`
//	Sale        int   `json:"sale"`
//	Tasktypes   int   `json:"tasktypes"`
//	Ns          []int `json:"n"`
//	Items       []int `json:"item"`
//	Nums        []int `json:"num"`
//	Start       int   `json:"start"`
//	End         int   `json:"end"`
//	Needlv      int   `json:"needlv"`
//	Needvips    []int `json:"needvip"`
//	Monthcard   int   `json:"monthcard"`
//	Redpacket   int   `json:"redpacket"`
//	Redpacketid int   `json:"redpacketid"`
//}
type GemsweeperitembagConfig struct {
	Group  int `json:"group"`
	Gid    int `json:"gid"`
	Gvalue int `json:"gvalue"`
	Level  int `json:"level"`
	Value  int `json:"value"`
	Itemid int `json:"itemid"`
	Count  int `json:"count"`
}
type WarlistConfig struct {
	Id    int   `json:"id"`
	Min   int   `json:"min"`
	Max   int   `json:"max"`
	Items []int `json:"item"`
	Nums  []int `json:"num"`
	Type  int   `json:"type"`
}
type ShopConfig struct {
	Id             int `json:"id"`
	Type           int `json:"type"`
	Grid           int `json:"grid"`
	Itemid         int `json:"itemid"`
	Itemnumber     int `json:"itemnumber"`
	Weightfunction int `json:"weightfunction"`
	Currency       int `json:"currency"`
	CostItem       int `json:"cost_item"`
	CostNum        int `json:"cost_num"`
	Teamlvmax      int `json:"teamlvmax"`
	Guildlimit     int `json:"guild_limit"`
	Prestigelimit  int `json:"prestige_limit"`
	Havehero       int `json:"havehero"`
}
type RobotConfig struct {
	Id             int   `json:"id"`
	Jjcrank        int   `json:"jjcrank"`
	Robotlevel     int   `json:"robotlevel"`
	Robotquality   int   `json:"robotquality"`
	Robotstar      int   `json:"robotstar"`
	Equipquality   int   `json:"equipquality"`
	Equiplevel     int   `json:"equiplevel"`
	Equipstar      int   `json:"equipstar"`
	Specialquality int   `json:"specialquality"`
	Speciallevel   int   `json:"speciallevel"`
	Specialstar    int   `json:"specialstar"`
	Optionheros    []int `json:"optionhero"`
}
type BuyskillpointsConfig struct {
	Number int `json:"number"`
	Money  int `json:"money"`
}
type MoneyConfig struct {
	Id            int    `json:"id"`    //! 充值ID
	Game          int    `json:"game"`  //! 游戏类型
	Grade         int    `json:"grade"` //! 充值档位
	Type          int    `json:"type"`
	Show          int    `json:"show"`
	Sort          int    `json:"sort"`
	Rmb           int    `json:"rmb"`
	Diamond       int    `json:"diamond"`
	Extra         int    `json:"extra"`
	Time          int    `json:"time"`
	Vipexp        int    `json:"vipexp"`
	Start         int    `json:"start"`
	Continued     int    `json:"continued"`
	Cd            int    `json:"cd"`
	Charatype     string `json:"chara_type"`
	Givemoneyshow int    `json:"givemoneyshow"`
}
type PubchestdropgroupConfig struct {
	Dropid             int `json:"drop_id"`
	Dropgroup          int `json:"drop_group"`
	Itemid             int `json:"item_id"`
	Dropcardsnum       int `json:"drop_cards_num"`
	Chance             int `json:"chance"`
	Alert              int `json:"alert"`
	Haveitem           int `json:"haveitem"`
	Havechance         int `json:"havechance"`
	WishChance         int `json:"wishchance"`
	TimeGeneralsNotice int `json:"time_generals_notice"` //跨服神将是否写入公告 默认0不写
}

type ActiveConfig struct {
	Id           int   `json:"id"`
	Type         int   `json:"type"`
	Sort         int   `json:"sort"`
	Active       int   `json:"active"`
	ItemIds      []int `json:"n"`
	ItemNums     []int `json:"x"`
	SpecialItem  int   `json:"special_item"`
	SpecialNum   int   `json:"special_num"`
	RecastItem   int   `json:"recast_item"`
	RecastNum    int   `json:"recast_num"`
	Rank         int   `json:"rank"`
	RankItem     int   `json:"rankitem"`
	RankNum      int   `json:"ranknum"`
	FestivalItem []int `json:"festival_item"` //活动6020 额外掉落
	FestivalNum  []int `json:"festival_num"`  //活动6020 额外掉落
}
type HolyPartsUpgradeConfig struct {
	Treasureaid             int     `json:"treasurea_id"`
	Name                    string  `json:"name"`
	Quality                 int     `json:"quality"`
	Stagelvshow             int     `json:"stage_lv_show"`
	Stagelv                 int     `json:"stage_lv"`
	Costitems               []int   `json:"cost_item"`
	Costnums                []int   `json:"cost_num"`
	Attributetypes          []int   `json:"attribute_type"`
	Attributes              []int64 `json:"attribute"`
	Extraattributestypes    []int   `json:"extra_attributes_type"`
	Extraattributeaddvalues []int64 `json:"extra_attribute_addvalue"`
	Announcement            int     `json:"announcement"`
	Maxlevel                int     `json:"maxlevel"`
}

type TrialConfig struct {
	Levelid   int   `json:"levelid"`
	Type      int   `json:"type"`
	Opens     []int `json:"open"`
	Speci     int   `json:"speci"`
	Diffculty int   `json:"diffculty"`
	Level     int   `json:"level"`
	Front     int   `json:"front"`
	Awards    []int `json:"award"`
	Items     []int `json:"item"`
}
type BespeakConfig struct {
	Id       int   `json:"id"`
	Hero     int   `json:"hero"`
	Day      int   `json:"day"`
	Discount int   `json:"discount"`
	Cost     int   `json:"cost"`
	Peoples  []int `json:"people"`
	Sale     int   `json:"sale"`
	Items    []int `json:"item"`
	Nums     []int `json:"num"`
	Open     int   `json:"open"`
	Time     int   `json:"time"`
}

type PlayernameConfig struct {
	Num    int      `json:"num"`
	Names  []string `json:"name"`
	Namenv string   `json:"namenv"`
}
type CamptaskConfig struct {
	Id        int   `json:"id"`
	Min       int   `json:"min"`
	Max       int   `json:"max"`
	Type      int   `json:"type"`
	World     int   `json:"world"`
	Intensity int   `json:"intensity"`
	Citytype  int   `json:"citytype"`
	Citynum   int   `json:"citynum"`
	Camps     []int `json:"camp"`
	Personals []int `json:"personal"`
	Time      int   `json:"time"`
}
type SpytreasureConfig struct {
	Id          int `json:"id"`
	Eventid     int `json:"event_id"`
	Office      int `json:"office"`
	Storyid     int `json:"storyid"`
	Spyid       int `json:"spyid"`
	Type        int `json:"type"`
	Prop        int `json:"prop"`
	Propexp     int `json:"propexp"`
	Currency    int `json:"currency"`
	Currencynum int `json:"currencynum"`
	Rebate      int `json:"rebate"`
	Chance      int `json:"chance"`
}
type VisitchanceConfig struct {
	Id       int `json:"id"`
	Hero     int `json:"hero"`
	Sound    int `json:"sound"`
	Levelid  int `json:"levelid"`
	Sounddes int `json:"sound_des"`
	Minlv    int `json:"minlv"`
	Maxlv    int `json:"maxlv"`
	Weight   int `json:"weight"`
	Time     int `json:"time"`
}
type NationalwarParmConfig struct {
	Type  int   `json:"type"`
	Parms []int `json:"parm"`
}
type ConsumetopshopConfig struct {
	Id      int `json:"id"`
	Sort    int `json:"sort"`
	Type    int `json:"type"`
	Time    int `json:"time"`
	Item    int `json:"item"`
	Num     int `json:"num"`
	Cost    int `json:"cost"`
	Costnum int `json:"costnum"`
	Group   int `json:"group"`
}
type ConsumetopluckConfig struct {
	Id     int   `json:"id"`
	Group  int   `json:"group"`
	Change int   `json:"change"`
	Items  []int `json:"item"`
	Nums   []int `json:"num"`
	Values []int `json:"value"`
}

type TimeGeneralsPointsConfig struct {
	Id     int   `json:"id"`
	Group  int   `json:"group"`
	Points int   `json:"points"`
	Award  []int `json:"award"`
	Nums   []int `json:"num"`
}

type ExpeditionConfig struct {
	Id         int    `json:"id"`
	PointName  string `json:"pointname"`
	Point      int    `json:"point"`
	Pointtype  int    `json:"pointtype"`
	Boxid      int    `json:"boxid"`
	Buffs      []int  `json:"buff"`
	Rewards    []int  `json:"reward"`
	Rewardnums []int  `json:"rewardnum"`
	Mapid      int    `json:"mapid"`
}
type WartaskConfig struct {
	Id       int   `json:"id"`
	Taskid   int   `json:"taskid"`
	Icon     int   `json:"icon"`
	Tasktype int   `json:"tasktype"`
	Ns       []int `json:"n"`
	Items    []int `json:"item"`
	Nums     []int `json:"num"`
}
type HorseSoulConfig struct {
	ID           int   `json:"ID"`
	Type         int   `json:"type"`
	Quality      int   `json:"quality"`
	Soultype     int   `json:"soul_type"`
	AttType      int   `json:"attribute_type"`
	AttTypeOne   int   `json:"attribute_typeone"`
	AttValue     int64 `json:"dragon_hole_attribute" trim:"0"`
	AttValueOne  int64 `json:"dragon_hole_attributeone" trim:"0"`
	DragonAtt    int   `json:"dragon_hole_attribute_type1" trim:"0"`
	DragonValue  int64 `json:"dragon_hole_attribute1" trim:"0"`
	Upgrademodel int   `json:"upgrade_model"`
	Soulmod      int   `json:"soul_mod"`
}
type FundbuyConfig struct {
	Memberid int `json:"memberid"`
	Fee      int `json:"fee"`
	Get      int `json:"get"`
	Bugvip   int `json:"bugvip"`
	Power    int `json:"power"`
}
type SevendayConfig struct {
	Id         int   `json:"id"`
	Items      []int `json:"item"`
	Nums       []int `json:"num"`
	Costitems  []int `json:"costitem"`
	Costnums   []int `json:"costnum"`
	RebateType int   `json:"rebate_type"`
	Param1     int   `json:"param1"`
}

type SevendayAward struct {
	Id        int    `json:"id"`         // 流水号
	Stage     int    `json:"stage"`      // 阶段
	NeedPoint int    `json:"need_point"` // 达到积分
	Items     []int  `json:"item_id"`    // 物品id
	Nums      []int  `json:"num"`        // 物品数量
	Name      string `json:"name"`       // 名字
}

type HorseBattleSteedAttributeConfig struct {
	Id                 int `json:"id"`
	Attributelistgroup int `json:"attribute_list_group"`
	Index              int `json:"index"`
	Weight             int `json:"weight"`
	Attributetype      int `json:"attribute_type"`
	Attributevalve     int `json:"attribute_valve"`
}
type HorseJudgecallConfig struct {
	ID                int   `json:"ID"`
	Type              int   `json:"type"`
	Callnum           int   `json:"call_num"`
	Costitems         []int `json:"cost_item"`
	Costnums          []int `json:"cost_num"`
	Weights           []int `json:"weight"`
	Dropsmallhorseids []int `json:"dropsmall_horse_id"`
	Dropitemnums      []int `json:"dropitem_num"`
	Getjudgeexp       int   `json:"get_judge_exp"`
}
type PubchestspecialConfig struct {
	Droptimemodifyid int    `json:"drop_time_modify_id"`
	SpType           int    `json:"sp_type"`
	Paytype          string `json:"pay_type"`
	Droptimemin      int    `json:"drop_time_min"`
	Droptimemax      int    `json:"drop_time_max"`
	DropGroupModify  string `json:"drop_group_modify"`
}

type SummonBoxConfig struct {
	Id    int   `json:"id"`
	Type  int   `json:"type"`
	Order int   `json:"order"`
	Scale int   `json:"scale"`
	Item  []int `json:"item"`
	Num   []int `json:"num"`
}

type JjcdamoninfoConfig struct {
	Id      int `json:"id"`
	Count   int `json:"count"`
	Price   int `json:"price"`
	Worship int `json:"worship"`
	Reset   int `json:"reset"`
	Change  int `json:"change"`
}
type ActivitynewConfig struct {
	Id        int    `json:"id"`
	Sort      int    `json:"sort"`
	Mode      int    `json:"mode"`
	Type      int    `json:"type"`
	Step      int    `json:"step"`
	Tasktypes int    `json:"tasktypes"`
	Ns        []int  `json:"n"`
	Items     []int  `json:"item"`
	Nums      []int  `json:"num"`
	Costitems []int  `json:"costitem"`
	Costnums  []int  `json:"costnum"`
	Start     string `json:"start"`
	Continued int    `json:"continued"`
	Cd        int    `json:"cd"`
	Show      int    `json:"show"`
	Renovate  int    `json:"renovate"`
	Reset     int    `json:"reset"`
	Status    int    `json:"status"`
	Button    int    `json:"button"`
}

type HorseJudgeDiscernConfig struct {
	ID                int `json:"ID"`
	Discernsmallhorse int `json:"discern_small_horse"`
	Costsmallhorse    int `json:"cost_small_horse"`
	Costitem          int `json:"cost_item"`
	Costnum           int `json:"cost_num"`
	Itembaggroup      int `json:"itembaggroup"`
	Odds              int `json:"odds"`
	Dropitemid        int `json:"dropitem_id"`
	Dropitemnum       int `json:"dropitem_num"`
	Getjudgeexp       int `json:"get_judge_exp"`
}

type WarplayConfig struct {
	Id                int   `json:"id"`
	Type              int   `json:"type"`
	Showcondition     int   `json:"showcondition"`
	Showcount         int   `json:"showcount"`
	Shownotice        int   `json:"shownotice"`
	Playcondition     int   `json:"playcondition"`
	Playcount         int   `json:"playcount"`
	Text              int   `json:"text"`
	Readytime         int   `json:"readytime"`
	Continuedtime     int   `json:"continuedtime"`
	Initialtime       int   `json:"initialtime"`
	Maxtime           int   `json:"maxtime"`
	Timecondition     int   `json:"timecondition"`
	Timecount         int   `json:"timecount"`
	Cue               int   `json:"cue"`
	Enrolltime        int   `json:"enrolltime"`
	Cdtime            int   `json:"cdtime"`
	Showtime          int   `json:"showtime"`
	Attackcost        int   `json:"attackcost"`
	Defensecost       int   `json:"defensecost"`
	Needclass         int   `json:"needclass"`
	Fighting          int   `json:"fighting"`
	Attacknumber      int   `json:"attacknumber"`
	Defensenumber     int   `json:"defensenumber"`
	Npctype           int   `json:"npctype"`
	Npcnum            int   `json:"npcnum"`
	Killhonor         int   `json:"killhonor"`
	Winhonor          int   `json:"winhonor"`
	Losehonor         int   `json:"losehonor"`
	Occupy            int   `json:"occupy"`
	Timeratio         int   `json:"timeratio"`
	Cityratio         int   `json:"cityratio"`
	Strongholdratio   int   `json:"strongholdratio"`
	Stageratio        int   `json:"stageratio"`
	Rateratio         int   `json:"rateratio"`
	Winbuff           int   `json:"winbuff"`
	Personalbuff      int   `json:"personalbuff"`
	Battlefieldbuffin int   `json:"battlefieldbuffin"`
	Personalbuffin    int   `json:"personalbuffin"`
	Rule              int   `json:"rule"`
	Armys             []int `json:"army"`
	Items             []int `json:"item"`
	Mode              int   `json:"mode"`
	Loseitem          int   `json:"loseitem"`
	Losenum           int   `json:"losenum"`
	Skills            []int `json:"skill"`
	Minnums           []int `json:"minnum"`
	Maxnums           []int `json:"maxnum"`
	Grandtotal        int   `json:"grandtotal"`
}
type OfficerobotConfig struct {
	Id           int   `json:"id"`
	Type         int   `json:"type"`
	Category     int   `json:"category"`
	Officer      int   `json:"officer"`
	Officernum   int   `json:"officernum"`
	Officertype  int   `json:"officertype"`
	Topmod       int   `json:"topmod"`
	Npclv        int   `json:"npclv"`
	Npcquality   int   `json:"npcquality"`
	Npcstar      int   `json:"npcstar"`
	Armquality   int   `json:"armquality"`
	Armlv        int   `json:"armlv"`
	Npctitle     int   `json:"npctitle"`
	Heroweaponlv int   `json:"heroweaponlv"`
	Npcs         []int `json:"npc"`
	Armtypes     []int `json:"armtype"`
	Tabletypes   []int `json:"tabletype"`
	Tablevalues  []int `json:"tablevalue"`
	Nums         []int `json:"num"`
}
type CrownfightConfig struct {
	Id           int    `json:"id"`
	Power        int    `json:"power"`
	PowerId      string `json:"powerid"`
	Type         int    `json:"type"`
	Page         int    `json:"page"`
	Place        int    `json:"place"`
	Holdreward   int    `json:"holdreward"`
	Hourget      int    `json:"hourget"`
	Kingadd      int    `json:"kingadd"`
	Officerclass int    `json:"officerclass"`
}
type HorseBattleSteedConfig struct {
	Id                     int   `json:"id"`
	Quality                int   `json:"quality"`
	Star                   int   `json:"star"`
	Holecount              int   `json:"hole_count"`
	Extractattributenum    int   `json:"extract_attribute_num"`
	Attributelistgroup     int   `json:"attribute_list_group"`
	Skilllv                int   `json:"skill_lv"`
	Skillid                int   `json:"skill_id"`
	Decomposeitembaggroups []int `json:"decompose_itembaggroup"`
	Holelvs                []int `json:"hole1_lv"`
}
type LeveltitletaskConfig struct {
	Id        int   `json:"id"`
	Tasktypes int   `json:"tasktypes"`
	Ns        []int `json:"n"`
	Methodss  []int `json:"methods"`
}
type DailytaskConfig struct {
	Taskid         int    `json:"taskid"`
	OrderNumber    int    `json:"ordernumber"`
	MainType       int    `json:"main_type"`
	Type           int    `json:"type"`
	Shinerank      int    `json:"shine_rank"`
	Level          int    `json:"level"`
	TaskName       string `json:"taskname"`
	Checkpoint     int    `json:"checkpoint"`
	Tasktypes      int    `json:"tasktypes"`
	Ns             []int  `json:"n"`
	Itemid         int    `json:"itemid"`
	Numberofcycles int    `json:"numberofcycles"`
	Starttime      int    `json:"starttime"`
	Endtime        int    `json:"endtime"`
	Teamexp        int    `json:"teamexp"`
	Rewards        []int  `json:"reward"`
	Numbers        []int  `json:"number"`
	Pretask        int    `json:"pretask"`
	Vip            int    `json:"vip"`
	TaskNode       *TaskNode
}

type TargetTaskConfig struct {
	Taskid       int   `json:"taskid"`
	System       int   `json:"system"`       //归属系统
	Group        int   `json:"group"`        //任务组
	Preposegroup int   `json:"preposegroup"` //前置任务组
	Condition    int   `json:"condition"`    //任务开启对关卡条件
	Tasktypes    int   `json:"tasktypes"`    //任务类型
	Ns           []int `json:"n"`            //条件
	Item         []int `json:"item"`         //奖励物品
	Num          []int `json:"num"`          //奖励数量
	ItemLv       []int `json:"itemlv"`       //奖励物品
	NumLv        []int `json:"numlv"`        //奖励数量
	Prepose      int   `json:"prepose"`      //前置任务
	TaskNode     *TaskNode
}

type BadgeTaskConfig struct {
	Id     int `json:"id"`
	System int `json:"system"` //归属系统
	Money  int `json:"money"`  //充值档
	Lv     int `json:"lv"`     //Lv
	Need   int `json:"need"`   //前置
}

type HolyUpgradeConfig struct {
	Beautyid                int     `json:"beauty_id"`
	Name                    string  `json:"name"`
	Quality                 int     `json:"quality"`
	Stagelvshow             int     `json:"stage_lv_show"`
	Stagelv                 int     `json:"stage_lv"`
	Treasureastagelvs       []int   `json:"treasurea1_stage_lv"`
	Treasureas              []int   `json:"treasurea"`
	Items                   []int   `json:"item"`
	Itemcosts               []int   `json:"item_cost"`
	Skillid                 int     `json:"skill_id"`
	Extraattributestypes    []int   `json:"extra_attributes_type"`
	Extraattributeaddvalues []int64 `json:"extra_attribute_addvalue"`
	Fight                   int     `json:"fight"`
	Announcement            int     `json:"announcement"`
	Beautypicture           int     `json:"beauty_picture"`
	Beautyicon              int     `json:"beauty_icon"`
	Rank                    int     `json:"rank"`
	Nicknamepicture         int     `json:"nickname_picture"`
	Scenepicture            int     `json:"scene_picture"`
	Text                    int     `json:"text"`
	Maxlevel                int     `json:"maxlevel"`
}
type HorseShopConfig struct {
	Id       int `json:"id"`
	Type     int `json:"type"`
	Itemid   int `json:"item_id"`
	Itemnum  int `json:"item_num"`
	Lv       int `json:"lv"`
	Costitem int `json:"cost_item"`
	Costnum  int `json:"cost_num"`
}
type ConsumetopattackConfig struct {
	Id    int   `json:"id"`
	Level int   `json:"level"`
	Cost  int   `json:"cost"`
	Items []int `json:"item"`
	Nums  []int `json:"num"`
}
type NpcConfig struct {
	Id          int `json:"id"`
	Location    int `json:"location"`
	Npcid       int `json:"npcid"`
	Type        int `json:"type"`
	Npcicon     int `json:"npcicon"`
	Background  int `json:"background"`
	Storyid     int `json:"storyid"`
	Nexte       int `json:"nexte"`
	Position    int `json:"position"`
	Soundtime   int `json:"soundtime"`
	Cdtime      int `json:"cdtime"`
	Probability int `json:"probability"`
	Item        int `json:"item"`
	Itemnum     int `json:"itemnum"`
	Open        int `json:"open"`
	Close       int `json:"close"`
}

type NationalwarEnrollmentConfig struct {
	Id            int   `json:"id"`
	Enrollmentnum int   `json:"enrollment_num"`
	Items         []int `json:"item"`
	Itemnums      []int `json:"item_num"`
}
type HorseJudgelevelConfig struct {
	ID                int `json:"ID"`
	Level             int `json:"level"`
	Needlv            int `json:"need_lv"`
	Exp               int `json:"exp"`
	Horseupperlimit   int `json:"horse_upper_limit"`
	Dailynormalcall   int `json:"daily_normal_call"`
	Dailyhigtcall     int `json:"daily_higt_call"`
	Recoveryhorsetime int `json:"recovery_horse_time"`
	Recoveryhorseid   int `json:"recovery_horse_id"`
	Recoveryhorsenum  int `json:"recovery_horse_num"`
}
type GrowthtaskConfig struct {
	Taskid         int    `json:"taskid"`
	TaskName       string `json:"taskname"`
	Ordernumber    int    `json:"ordernumber"`
	Level          int    `json:"level"`
	Checkpoint     int    `json:"checkpoint"`
	Pretask        int    `json:"pretask"`
	Tasktypes      int    `json:"tasktypes" trim:"0"`
	Ns             []int  `json:"n"`
	Itemid         int    `json:"itemid"`
	Numberofcycles int    `json:"numberofcycles"`
	Rewards        []int  `json:"reward"`
	Numbers        []int  `json:"number"`
	TaskNode       *TaskNode
}
type HomeofficepurchaseConfig struct {
	Counttype      int   `json:"count_type"`
	Accesstype     int   `json:"access_type"`
	Maxcount       int   `json:"max_count"`
	Countcost      int   `json:"count_cost"`
	Getprops       int   `json:"Get_props"`
	Basevalue      int   `json:"base_value"`
	Parms          []int `json:"parm"`
	Luckadd        int   `json:"luck_add"`
	Mcrate         int   `json:"mcrate"`
	Mcrits         []int `json:"mcrit"`
	Cityvaluegroup int   `json:"city_value_group"`
}

type WorldpowerConfig struct {
	Id      int   `json:"id"`
	Warlord int   `json:"warlord"`
	Flag    int   `json:"flag"`
	Name    int   `json:"name"`
	Citys   []int `json:"city"`
	Citynum int   `json:"citynum"`
	Parts   []int `json:"part"`
	Items   []int `json:"item"`
	Nums    []int `json:"num"`
}

type PvpAwardConfig struct {
	Id      int   `json:"id"`
	Type    int   `json:"type"`
	AwardId int   `json:"awardid"`
	Rank1   int   `json:"rank1" trim:"0"`
	Rank2   int   `json:"rank2" trim:"0"`
	Item    []int `json:"item"`
	Num     []int `json:"num"`
}

type TariffConfig struct {
	Id       int   `json:"id"`
	Type     int   `json:"type"`           // 功能类型
	Rank1    int   `json:"rank1" trim:"0"` // 排行上区间
	Rank2    int   `json:"rank2" trim:"0"` // 排行下区间
	ItemIds  []int `json:"costitem"`       // 消耗道具
	ItemNums []int `json:"costnum"`        // 道具数量
	GetItem  []int `json:"getitem"`        // 获得道具
	GetNum   []int `json:"getnum"`         // 道具数量
}

type BossConfig struct {
	Id            int     `json:"id"`            // 流水号
	Name          string  `json:"name"`          // 巨兽名称
	Quality       int     `json:"quality_color"` // 品质
	Picture       int     `json:"picture"`       // 巨兽图片
	Head          string  `json:"head"`          // 巨兽头像
	Rim           string  `json:"rim"`           // 巨兽边框
	Time          int64   `json:"time"`          // 购买时间
	Buy           int     `json:"buy"`           // 购买价格
	Vip           int     `json:"vip"`           // 打折vip等级
	Sale          float32 `json:"sale"`          // 打折价格
	Renew         float32 `json:"renew"`         // 自动续费价格
	HeroId        int     `json:"heroid"`        // 英雄Id
	Permanent     int     `json:"permanent"`     // 巨兽状态。0未开启，1限时激活，2永久激活
	ItemIds       []int   `json:"costitem"`      // 消耗Id
	ItemNums      []int   `json:"costnum"`       // 消耗num
	Skills        []int   `json:"skill"`         // 巨兽技能
	HeadportaitId int     `json:"headportaitId"` // 头像
	RimId         int     `json:"rimId"`         // 相框
	Openskill     int     `json:"openskill"`     // 巨兽开场技能
}

type HangUpConfig struct {
	Id           int   `json:"id"`           //流水号
	GoldYield    int   `json:"goldyield"`    //金币产量
	HeroExpYield int   `json:"heroexpyield"` //经验产量
	ExpYield     int   `json:"expyield"`     //领主经验
	PowderYield  int   `json:"powderyield"`  //魔粉产量
	BasicsTime   int64 `json:"basicstime"`   //基础库间隔 单位S
	BasicsDrop   []int `json:"basicsdrop"`   //基础掉落组
	SeniorTime   int64 `json:"seniortime"`   //高级库间隔 单位S
	SeniorDrop   []int `json:"seniordrop"`   //高级掉落组
}

type GemstoneChapterConfig struct {
	Id          int `json:"id"`
	ChaptrGroup int `json:"chaptr_group"`
	LimitLevel  int `json:"limit_level"`
}

type GemstoneLevelConfig struct {
	Id           int   `json:"id"`
	ChaptrGroup  int   `json:"chaptr_group"`
	LevelIndex   int   `json:"level_index"`
	LevelId      int   `json:"level"`
	FirstItem    []int `json:"first_item"`
	Num          []int `json:"num"`
	ItemBagGroup []int `json:"itembag_group"`
	Costids      []int `json:"cost_id"`
	Costnums     []int `json:"cost_num"`
}

type JJCRobotConfig struct {
	Id         int     `json:"id"`
	Type       int     `json:"type"`
	Jjcclass   int     `json:"jjcclass"`
	Jjcdan     int     `json:"jjcdan"`
	Category   int     `json:"jjccategory"`
	Teamnum    int     `json:"teamnum"`
	Rank1      int     `json:"jjcrankmin" trim:"0"`
	Rank2      int     `json:"jjcrankmax" trim:"0"`
	Name       string  `json:"name"`
	Head       int     `json:"head"`
	MonsterId  int     `json:"Monsterid"`
	NpcLv      []int   `json:"herolv"`
	Level      int     `json:"robotlevel"`
	NpcQuality int     `json:"npcquality"`
	NpcStar    []int   `json:"robotstar"`
	Fight      []int64 `json:"showfight"`
	//MaxFight   int    `json:"showfight1"`
	Hero []int `json:"optionhero"`
	//Arms       []int     `json:"armtype"`
	Hydra      int       `json:"hydraoption"` //! 巨兽
	Rage       int       `json:"rageoption"`
	BaseTypes  []int     `json:"base_type"`
	BaseValues []float64 `json:"base_value"`
}

type LevelMonsterConfig struct {
	MonsterId    int     `json:"monsterid"`
	HeroId       int     `json:"heroid"`
	MonsterIndex int     `json:"monster_index"`
	Level        int     `json:"level"`
	MainSkill    int     `json:"mainskill"`
	BaseType     []int   `json:"base_type"`
	BaseValue    []int64 `json:"base_value"`
}

type HorseAward struct {
	Id      int   `json:"id"`
	Times   int   `json:"higt_call"`
	Rewards []int `json:"reward"`
	Nums    []int `json:"num"`
}

type HorseSwitch struct {
	HorseId int `json:"warhorse_id"`
	Rate    int `json:"weight"`
}

// Horse_BattleSteed
type HorseBattleSteed struct {
	HorseId          int `json:"id"`
	Quality          int `json:"quality"`
	ExtractAttNum    int `json:"extract_attribute_num"`
	ExtractListGroup int `json:"attribute_list_group"`
}

type MercenaryRandom struct {
	Class     int `json:"class"`         // 阶级
	Index     int `json:"index"`         // 资质
	BaseType  int `json:"base_type"`     // 属性类型
	MaxValue  int `json:"max_value"`     // 属性数值上限
	InitValue int `json:"initial_value"` // 属性初始最高值
	Weight    int `json:"weight"`        // 属性的权重
	Limit     int `json:"limit"`         // 极限保护
	LowUp     int `json:"low_up"`        // 小于限定时升的权重
	LowDown   int `json:"low_down"`      // 小于限定时降的权重
	LowFlat   int `json:"low_flat"`      // 小于限定时平的权重
	Threshold int `json:"threshold"`     // 限定万分比
	HighUp    int `json:"high_up"`       // 大于等于限定时升的权重
	HighDown  int `json:"high_down"`     // 大于等于限定时降的权重
	HighFlat  int `json:"high_flat"`     // 大于等于限定时平的权重
	ChangeMax int `json:"change_max"`    // 单次升降最大值
	Unit      int `json:"unit"`          // 评分的数值单位
	Evaluate  int `json:"evaluate"`      // 对应评分数值
}

type MercenaryLv struct {
	Id        int     `json:"id"`         // 编号
	Lv        int     `json:"lv"`         // 佣兵等级
	NeedExp   int     `json:"needexp"`    // 需要经验值
	GetExp    int     `json:"getexp"`     // 转化经验值
	Fire      int     `json:"fire"`       // 转化钱币
	BaseType  []int   `json:"base_type"`  // 属性类型
	BaseValue []int64 `json:"base_value"` // 属性值
	GetChange int     `json:"get_chance"` // 升级时获得副属性概率,万分比格式
	Gets      []int   `json:"get"`        // 获得属性权重和
}

type MercenaryConfig struct {
	Id        int   `json:"mercenary_id"`    // id
	Level     int   `json:"mercenary_level"` // 阶级
	Camp      int   `json:"mercenary_camp"`  // 种族
	Type      int   `json:"mercenary_type"`  // 类型 1近战 2远程
	Hole      int   `json:"mercenary_hole"`  // 佣兵位置孔
	Index     int   `json:"mercenary_index"` // 资质 1B 2A 3S 4SS
	BaseType  []int `json:"base_type"`       // 相同上阵时佣兵激活的属性类型
	BaseValue []int `json:"base_value"`      // 相同上阵时佣兵激活的属性数值
	Skills    []int `json:"mercenary_skill"` // 技能
}

type WorldPower struct {
	Id       int    `json:"id"` // id
	KingName string `json:"king_name"`
	Name     string `json:"name"`
}

type GrowthtaskKingConfig struct {
	Taskid    int   `json:"taskid"`
	Sort      int   `json:"sort"`
	Group     int   `json:"group"`
	Type      int   `json:"type"`
	Awardtype int   `json:"awardtype"`
	Tasktypes int   `json:"tasktypes"`
	Ns        []int `json:"n"`
	Items     []int `json:"item"`
	Itemnums  []int `json:"itemnum"`
	Values    []int `json:"value"`
	Latertask int   `json:"latertask"`
}

type HeroAttribute struct {
	ValueType int    `json:"valuetype"`
	FightNum  int64  `json:"fightingnum"`
	Name      string `json:"name"`
}

type CrownBuildConfig struct {
	Type    int    `json:"type"`
	Name    string `json:"name"`
	Rewards []int  `json:"reward"`
	Items   []int  `json:"item"`
	Open    int    `json:"open"`
}

// 关卡掉落
type LevelItemConfig struct {
	Levelid         int   `json:"levelid"`
	Items           []int `json:"item"`
	Probabilitys    []int `json:"probability"`
	Protections     []int `json:"protection"`
	Nums            []int `json:"num"`
	Lotteryids      []int `json:"lotteryid"`
	Bags            []int `json:"bag"`
	Bagprobabilitys []int `json:"bagprobability"`
}

// 国战相关配置
type PlayTimeConfig struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`
	Openday          int64    `json:"openday"`
	Continuedday     int      `json:"continuedday"`
	Cdday            int64    `json:"cdday"`
	Declarestart     int      `json:"declarestart"`
	Declareend       int      `json:"declareend"`
	Enrollstart      int      `json:"enrollstart"`
	Enrollend        int      `json:"enrollend"`
	Battlestart      int      `json:"battlestart"`
	Battleend        int      `json:"battleend"`
	Battles          []int    `json:"battle"`
	Readytime        int64    `json:"readytime"`
	NoticeStart      []int    `json:"notice_start"`
	NoticeEnd        []int    `json:"notice_end"`
	Interval         []int64  `json:"interval"`
	Notice           []string `json:"notice"`
	BattleStartHour  int
	BattleStartMin   int
	BattleEndHour    int
	BattleEndMin     int
	DeclareStartHour int
	DeclareStartMin  int
	DeclareEndHour   int
	DeclareEndMin    int
	EnrollStartHour  int
	EnrollStartMin   int
	EnrollEndHour    int
	EnrollEndMin     int
	BattleInfo       [3]HourMinutes
}

type HourMinutes struct {
	Hour    int
	Minutes int
}

type WorldMapConfig struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	First          int    `json:"first"`
	CityState      int    `json:"cityState"`
	CityInitial    int    `json:"cityinitial"`
	Upper          int    `json:"upper"`
	WeiRange       int    `json:"weirange"`
	WeiArmy        int    `json:"weiarmy"`
	ShuRange       int    `json:"shurange"`
	ShuArmy        int    `json:"shuarmy"`
	WuRange        int    `json:"wurange"`
	WuArmy         int    `json:"wuarmy"`
	NpcLV          int    `json:"npclv"`
	NpcArmy        int    `json:"npcarmy"`
	Size           int    `json:"size"`
	SizeDescribe   string `json:"sizedescribe"`
	Terrain        int    `json:"terrain"`
	InitialId      int    `json:"initialid"`
	InitialDean    string `json:"initialdean"`
	CharacterIstic []int  `json:"characteristic"`
	Item           int    `json:"item"`
	Num            int    `json:"num"`
	State          int    `json:"state"`
	StateName      string `json:"statename"`
	StateDay       int    `json:"stateday"`
	StateBox       int    `json:"statebox"`
	StateNum       int    `json:"statenum"`
	Sequence       []int  `json:"sequence"`
}

type PlayRewardConfig struct {
	Play       int    `json:"play"`
	Id         int    `json:"id"`
	Type       int    `json:"type"`
	Minorder   int    `json:"min_order"`
	Maxorder   int    `json:"max_order"`
	Txt        string `json:"txt"`
	Items      []int  `json:"item"`
	Nums       []int  `json:"num"`
	Baseitemid int    `json:"base_itemid"`
	Basenum    int    `json:"base_num"`
	Mailtype   int    `json:"mail_type"`
}

// 矿战斗
type MineBattleMapConfig struct {
	ID               int    `json:"ID"`
	Type             int    `json:"type"`
	Score            int    `json:"score"`
	Defaultcamp      int    `json:"default_camp"`
	Baseitemid       int    `json:"base_itemid"`
	Basenum          int    `json:"base_num"`
	Capacity         int    `json:"capacity"`
	Doublepoint      int    `json:"double_point"`
	Pointdescription string `json:"point_description"`
	Links            []int  `json:"link_id"`
}

type MineBattleOrderConfig struct {
	Id        int    `json:"id"`
	Ordertype int    `json:"order_type"`
	Type      int    `json:"type"`
	Params    []int  `json:"param"`
	Cdtime    int64  `json:"cd_time"`
	Getweight int    `json:"get_weight"`
	Name      string `json:"name"`
	Describe  string `json:"describe"`
}

type MineBattleParamConfig struct {
	Id       int    `json:"id"`
	Type     int    `json:"type"`
	Params   []int  `json:"param"`
	Describe string `json:"describe"`
}

// open level
type OpenLevelConfig struct {
	Id     int `json:"id"`
	Level  int `json:"level"`
	Passid int `json:"passid"`
}

type MineBattleDoublePointConfig struct {
	Id          int   `json:"id"`
	Doublestart int   `json:"double_start"`
	Continued   int64 `json:"continued"`
	StartHour   int
	StartMin    int
}

// 邮件配置
type MailConfig struct {
	Id        int    `json:"id"`
	Mailtitle string `json:"mail_title"`
	Mailtxt   string `json:"mail_txt"`
}

// Gve建筑配置
type GveBuildingConfig struct {
	Id           int    `json:"id"`           // 建筑Id
	Camp         int    `json:"camp"`         // 阵营
	Campname     string `json:"campname"`     // 阵营名字
	Name         string `json:"name"`         // 建筑名字
	Type         int    `json:"type"`         // 建筑类型
	Disappear    int    `json:"disappear"`    // 是否消失
	Connections  []int  `json:"connection"`   // 连接信息
	Takenum      int    `json:"take_num"`     // 占领值
	Enemygroups  []int  `json:"enemy_group"`  // 怪物级别1,2,3
	Enemyrandoms []int  `json:"enemy_random"` // 随机百分比,独立随机:2,3
	Enemytimes   []int  `json:"enemy_time"`   // 刷出来的次数,2,3
	IsHold       int    `json:"ishold"`       // 是否可以刷出恶龙
}

// Gve掉落信息
type GveCheckpointGroupConfig struct {
	Checkpointgroup int `json:"checkpoint_group"` // 检查组
	Levelid         int `json:"levelid"`          // 关卡Id
	Victorydrop     int `json:"victory_drop"`     // 胜利掉落
	Failuredrop     int `json:"failure_drop"`     // 失败掉落
	Victoryroulette int `json:"victory_roulette"` // 翻牌奖励
}

// Gve技能配置
type GveskillsConfig struct {
	Id            int    `json:"id"`              // 技能时间
	Name          string `json:"name"`            // 名字
	Cdtime        int    `json:"cd_time"`         // 使用Cd
	Expend        int    `json:"expend"`          // 消耗道具Id
	Itemnum       int    `json:"itemnum"`         // 消耗道具数量
	Param1        int    `json:"param1" trim:"0"` // 阵营前五才能召唤恶龙
	Effectminimum int    `json:"effect_minimum"`  // 技能效果最小值
	Effectbiggest int    `json:"effect_biggest"`  // 技能效果最大值
	Notice        string `json:"notice"`          // 公告内容
	RewardMin     int    `json:"reward_minimum"`  // 获得道具最小值
	RewardMax     int    `json:"reward_biggest"`  // 获得道具最大值
}

type GvewayConfigConfig struct {
	Wayid         int `json:"way_id"`
	Startbuilding int `json:"start_building"`
	Endbuilding   int `json:"end_building"`
	Movetime      int `json:"movetime"`
}

// 鼓舞配置
type WarEncourageConfig struct {
	Id        int   `json:"id"`
	Lv        int   `json:"lv"`
	Type      int   `json:"type"`
	Costitems []int `json:"cost_item"`
	Costnums  []int `json:"cost_num"`
	Min       int   `json:"min"`
	Max       int   `json:"max"`
	Crtchance int   `json:"crtchance"`
	Crtnum    int   `json:"crtnum"`
	Encpoints int   `json:"encpoints"`
}

// 鼓舞掉落
type WarBuffConfig struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Lv            int     `json:"lv"`
	Dec           string  `json:"dec"`
	Group         int     `json:"group"`
	Type          int     `json:"type"`
	Superposition int     `json:"superposition"`
	Basetypes     []int   `json:"base_type"`
	Basevalues    []int64 `json:"base_value"`
}

// 头像配置
type HeadConfig struct {
	Id        int    `json:"id"`
	Resources string `json:"Resources"`
	Type      int    `json:"type"`
	Timetype  int    `json:"timetype"`
	Timevalue int    `json:"timevalue"`
	Open      int    `json:"open"`
	Condition int    `json:"condition"`
	Item      int    `json:"item"`
	Name      string `json:"name"`
	Dec       string `json:"dec"`
}

type MongytaskListConfig struct {
	Taskid        int    `json:"taskid"`
	Sort          int    `json:"sort"`
	Taskname      string `json:"taskname"`
	Star          int    `json:"star"`
	Type          int    `json:"type"`
	Tasktypes     int    `json:"tasktypes"`
	Ns            []int  `json:"n"`
	Content       string `json:"content"`
	Group         int    `json:"groupdef"`
	Value         int    `json:"value"`
	OpenLevel     int    `json:"openleve"`
	Openvip       int    `json:"openvip"`
	Showitems     []int  `json:"showitem"`
	Surprisevalue int    `json:"surprisevalue"`
	Surpriseitem  int    `json:"surpriseitem"`
	Surprisenum   int    `json:"surprisenum"`
	Worstitem     int    `json:"worstitem"`
	Worstnum      int    `json:"worstnum"`
	Mustitem      int    `json:"mustitem"`
	Num           int    `json:"num"`
	Moment        int    `json:"moment"`
	Draws         []int  `json:"draw"`
	Drawcosts     []int  `json:"drawcost"`
	Groups        []int  `json:"group"`
	Drawall       int    `json:"drawall"`
	Replaceitem   int    `json:"replaceitem"`
	Replacenum    int    `json:"replacenum"`
}

type MongytaskStarlistConfig struct {
	Id    int `json:"id"`
	Star  int `json:"star"`
	Value int `json:"value"`
}

type ArmyConfig struct {
	Id         int `json:"id"`
	Level      int `json:"level"`
	Camp       int `json:"camp"`
	Type       int `json:"type"`
	Index      int `json:"index"`
	Attacktype int `json:"attacktype"`
	Skill      int `json:"skill"`
	Sort       int `json:"sort"`
	Show       int `json:"show"`
}

type ArmyExchange struct {
	Id       int `json:"id"`
	Item     int `json:"item"`
	Type     int `json:"type"`
	Num      int `json:"num"`
	Needitem int `json:"needitem"`
	Neednum  int `json:"neednum"`
}

type ArmyLvConfig struct {
	KeyId      int     `json:"lazyid"`
	Id         int     `json:"id"`
	Lv         int     `json:"lv"`
	Needitem   int     `json:"needitem"`
	Neednum    int     `json:"neednum"`
	Basetypes  []int   `json:"base_type"`
	Basevalues []int64 `json:"base_value"`
}

type ArmyFlagConfig struct {
	Id             int `json:"id"`
	Index          int `json:"index"`
	Sort           int `json:"sort"`
	Identification int `json:"identification"`
}

type ArmyFlagLvConfig struct {
	KeyId      int     `json:"lazyid"`
	Id         int     `json:"id"`
	Lv         int     `json:"lv"`
	NeedLv     int     `json:"needlv"`
	Needitem   int     `json:"needitem"`
	Neednum    int     `json:"neednum"`
	Basetypes  []int   `json:"base_type"`
	Basevalues []int64 `json:"base_value"`
}

type ArmyChallengeConfig struct {
	Id        int   `json:"id"`
	Day       int   `json:"day"`
	Type      int   `json:"type"`
	Armyids   []int `json:"armyid"`
	Drops     []int `json:"drop"`
	Paydrops  []int `json:"paydrop"`
	Items     []int `json:"item"`
	Surprised int   `json:"surprised"`
}

type ArmyTeamConfig struct {
	Levelid     int    `json:"levelid"`
	Mapid       int    `json:"mapid"`
	ColdTime    int64  `json:"coldtime"`
	Refreshs    []int  `json:"refresh"`
	Time        int    `json:"time"`
	Chaptername string `json:"chaptername"`
	Win         int    `json:"win"`
	Lose        int    `json:"lose"`
	Reduce      int    `json:"reduce"`
	MaxChance   int    `json:"maxchance"`
	MinChance   int    `json:"minchance"`
	Drops       []int  `json:"drop"`
}

type BossLvConfig struct {
	Id         int     `json:"id"`
	Lv         int     `json:"lv"`
	Basetypes  []int   `json:"base_type"`
	Basevalues []int64 `json:"base_value"`
}

type StrStringConfig struct {
	Dec string `json:"dec"`
	Str string `json:"str"`
}

type StatisticsConfig struct {
	Id            int    `json:"id"`
	Lable         int    `json:"lable"`
	Type          int    `json:"type"`
	ShowType      int    `json:"show_type"`
	SubType       int    `json:"sub_type"`
	Openlevel     int    `json:"openlevel_id"`
	DescribeState string `json:"describe_state"`
	Percentage    int    `json:"percentage"`
}

type StatisticsRewardConfig struct {
	Id    int   `json:"id"`
	Point int   `json:"point"`
	Items []int `json:"item_id"`
	Nums  []int `json:"num"`
}

type NobilityConfig struct {
	Id        int   `json:"id"`
	TaskGroup int   `json:"task"`
	TaskType  int   `json:"tasktypes"`
	Ns        []int `json:"n"`
	Items     int   `json:"item"`
	Nums      int   `json:"num"`
}

type RechargeConfig struct {
	Id            int   `json:"id"`
	TotleRecharge int   `json:"totle_recharge"`
	AwardItems    []int `json:"award_item"`
	Nums          []int `json:"num"`
}

type WholeShopConfig struct {
	Id             int `json:"id"`
	Group          int `json:"group"`
	Type           int `json:"type"`
	Times          int `json:"time"`
	Item           int `json:"item"`
	Num            int `json:"num"`
	Cost           int `json:"cost"`
	CostNum        int `json:"costnum"`
	ShowValue      int `json:"show_value"`
	Recommend      int `json:"recommend"`
	Efficacy       int `json:"efficacy"`
	Sort           int `json:"sort"`
	Grid           int `json:"grid"`
	Weightfunction int `json:"weightfunction"`
	Personal_limit int `json:"personal_limit"`
}

type HeroExpConfig struct {
	HeroLv        int   `json:"herolv"`
	CostItems     []int `json:"costitem"`
	CostNums      []int `json:"costnum"`
	ResetItems    []int `json:"reset_item"`
	ResetNums     []int `json:"reset_num"`
	ResetCostId   int   `json:"reset_costid"`
	ResetCostNums int   `json:"reset_costnum"`
}

type RuneCompose struct {
	Id          int   `json:"id"`
	Type        int   `json:"type"`
	Item        []int `json:"item"`
	Num         []int `json:"num"`
	Rune        []int `json:"rune"`
	Random      []int `json:"random"`
	Probability []int `json:"probability"`
}

type RuneConfig struct {
	Id           int     `json:"id"`
	Level        int     `json:"type"`
	Index        []int   `json:"index"`         //品质
	Skill        []int   `json:"skill"`         //技能
	SkillV       []int   `json:"skill_v"`       //技能等级
	UpgradeType  []int   `json:"upgrade_type"`  // 属性类型
	UpgradeValue []int64 `json:"upgrade_value"` // 属性值
}

type FormationConfig struct {
	Id    int     `json:"id"`
	Index []int   `json:"index"` //位置开放
	Type  []int   `json:"type"`  // 属性类型
	Value []int64 `json:"value"` // 属性值
	Level int     `json:"level"` // 开放等级
}

type FundConfigMap struct {
	Id              int    `json:"id"`                //任务ID
	Group           int    `json:"group"`             //被activenew表调用组，用于配置不同期数
	Type            int    `json:"type"`              //1=勇者基金 2=至尊基金
	Pay             int    `json:"pay"`               //调用充值档位
	TaskTypes       int    `json:"tasktypes"`         //购买月卡前，需要判断条件
	Ns              []int  `json:"n"`                 //判断
	ConditionTxt    string `json:"condition_txt"`     //倍数显示
	Day             int64  `json:"day"`               //领取天数
	Item            []int  `json:"item"`              //奖励道具
	Num             []int  `json:"num"`               //奖励数量
	TotleShow       int    `json:"totle_show"`        //总奖励显示
	MultipleShow    int    `json:"multiple_show"`     //倍数显示
	Button          string `json:"button"`            //充值按钮显示
	EffectShow      string `json:"effect_show"`       //每日奖励对应的充值按钮特效
	TotleEffectShow string `json:"totle_effect_show"` //总奖励显示中，对应的ID道具闪光
}

type NewShopConfig struct {
	Id             int   `json:"id"`
	Type           int   `json:"type"`
	Grid           int   `json:"grid"`
	Group          int   `json:"group"`
	ItemId         int   `json:"itemid"`
	ItemNumber     int   `json:"itemnumber"`
	WeightFunction int   `json:"weightfunction"`
	CostItems      []int `json:"cost_item"`
	CostNums       []int `json:"cost_num"`
	LevelLimit     int   `json:"level_limit"`
	LevelShield    int   `json:"level_shield"`
	Judge          int   `json:"judge"`
	ReplaceGroup   int   `json:"replacegroup"`
	ReplaceItem    int   `json:"replaceitem"`
	ReplaceNum     int   `json:"replacenum"`
	ReplaceWeight  int   `json:"replaceweight"`
	ReplaceCost    []int `json:"replacecost"`
	ReplaceCostNum []int `json:"replacecostnum"`
}

type NewShopDiscount struct {
	Id          int   `json:"id"`
	Shop        int   `json:"shop"`
	Grid        []int `json:"grid"`
	LevelLimit  []int `json:"level_limit"`
	LevelShield []int `json:"level_shield"`
	Discount    []int `json:"discount"`
	Chance      []int `json:"chance"`
}

type WholeShopTimeConfig struct {
	Id          int `json:"id"`
	Refreshtime int `json:"refreshtime"`
}

type TurnTableConfig struct {
	Id     int   `json:"id"`
	Stage  int   `json:"stage"`
	Items  []int `json:"item"`
	Nums   []int `json:"num"`
	Weight int   `json:"weight"`
}

type TurnTableTimeConfig struct {
	Id    int   `json:"id"`
	Group int   `json:"group"`
	Index int   `json:"index"`
	Time  int64 `json:"time"`
}

type NobilityReward struct {
	Id        int   `json:"id"`
	Task      int   `json:"task"`
	Items     []int `json:"item"`
	Nums      []int `json:"num"`
	Vip       int   `json:"vip"`
	VipItems  []int `json:"vipitem"`
	VipNums   []int `json:"vipnum"`
	Privilege int   `json:"privilege"`
}

//! 限时礼包
type TimeGiftConfig struct {
	Id        int      `json:"id"`
	Group     int      `json:"group"`
	TaskTypes int      `json:"tasktypes"`
	Sale      int      `json:"sale"`
	N         []int    `json:"n"`
	Rmb       int      `json:"rmb"`
	Items     []int    `json:"item"`
	Nums      []int    `json:"num"`
	Describe  []string `json:"describe"`
	TabName   string   `json:"tab_name"`
	TabPic    string   `json:"tab_pic"`
	Pic       string   `json:"pic"`
	Efficacy  []int    `json:"efficacy"`
	EndDesc   string   `json:"end_describe"`
}

type TeamAttrConfig struct {
	Type       int     `json:"type"`
	Camp       []int   `json:"camp"`
	Base_type  []int   `json:"base_type"`
	Base_value []int64 `json:"base_value"`
}

//! 日周月新手礼包
type ActivityGiftConfig struct {
	Id             int      `json:"id"`
	ActivityType   int      `json:"activity_type"`
	Type           int      `json:"type"`
	RefreshTime    int      `json:"refresh_time"`
	Index          int      `json:"index"`
	Sort           int      `json:"sort"`
	Group          int      `json:"group"`
	TaskTypes      int      `json:"tasktypes"`
	N              []int    `json:"n"`
	Items          []int    `json:"item"`
	Nums           []int    `json:"num"`
	Sale           int      `json:"sale"`
	Times          int      `json:"times"`
	ShowValue      string   `json:"show_value"`
	Name           string   `json:"name"`
	Quality        int      `json:"quality"`
	Twinkle        []int    `json:"twinkle"`
	Pic            []string `json:"pic"`
	Start          string   `json:"start"`
	Continued      int      `json:"continued"`
	Cd             int      `json:"cd"`
	Pic2Type       int      `json:"pic2_type"`
	CostPrice      int      `json:"cost_price"`
	RechargeAmount int      `json:"recharge_amount"`
	StarHero       int      `json:"star_hero"`
	EffectShow     string   `json:"effect_show"`
	VipNum         string   `json:"vipnum"`
}

//! 成长礼包
type GrowthGiftConfig struct {
	Id        int   `json:"id"`
	TaskTypes int   `json:"tasktypes"`
	N         []int `json:"n"`
	Items     []int `json:"award_item"`
	Nums      []int `json:"num"`
}

type ActivityTotleAward struct {
	Id       int   `json:"id"`
	Type     int   `json:"type"`
	TotalNum int   `json:"total_num"`
	Item     []int `json:"item"`
	Num      []int `json:"num"`
}

type MonthCardTotleAwardMap struct {
	Id        int   `json:"id"`
	PointItem int   `json:"point_item"`
	PointNum  int   `json:"point_num"`
	Item      []int `json:"item"`
	Num       []int `json:"num"`
}

type MonthCard struct {
	Id            int `json:"id"`
	Type          int `json:"type"`
	FirstRecharge int `json:"first_recharge"`
	HangUpFast    int `json:"hangup_fast"`
	HangUpHeroExp int `json:"hangup_heroexp"`
	HangUpGold    int `json:"hangup_gold"`
}

type HydraConfig struct {
	HydraID    int    `json:"hydraid"`
	Icon       string `json:"icon"`
	FightSkill int    `json:"fightskill"`
	Skill      []int  `json:"skill"`
}

type HydraSkill struct {
	SkillID    int     `json:"skillid"`
	HydraLv    int     `json:"lv"`
	Unlock     int     `json:"unlock"`
	Items      []int   `json:"lvupneed"`
	Nums       []int   `json:"lvupnum"`
	RetuenNums []int   `json:"returnnum"`
	Base_type  []int   `json:"base_type"`
	Base_value []int64 `json:"base_value"`
}

type HydraLevel struct {
	Level      int   `json:"lv"`
	Items      []int `json:"lvupneed"`
	Nums       []int `json:"lvupneednum"`
	ReturnNums []int `json:"returnnum"`
}

type HydraStep struct {
	ID       int   `json:"id"`
	MaxLevel int   `json:"maxlv"`
	Items    []int `json:"lvupneed"`
	Nums     []int `json:"lvupneednum"`
}

//type HydraStar struct {
//	Id         int     `json:"id"`
//	HydraId    int     `json:"hydra_id"`
//	HydraLv    int     `json:"hydra_lv"`
//	Items      []int   `json:"qualification1_id"`
//	Nums       []int   `json:"qualification1_num"`
//	Exp        int     `json:"qualification_exp"`
//	Count      int     `json:"qualification"`
//	Base_type  []int   `json:"qualification1_type"`
//	Base_value []int64 `json:"qualification1_value"`
//}

type HydraTask struct {
	Id       int   `json:"id"`
	TaskId   int   `json:"taskid"`
	TaskType int   `json:"tasktypes"`
	Ns       []int `json:"n"`
	Items    []int `json:"item"`
	Nums     []int `json:"num"`
}

type PitConfig struct {
	Id        int   `json:"id"`
	PitId     int   `json:"dungeons_id"`
	ThingType int   `json:"thingtype"`
	BattleId  int   `json:"battle_id"`
	Condition []int `json:"condition"`
	Items     []int `json:"item_type"`
	Nums      []int `json:"item_number"`
}

type PitMap struct {
	Id           int   `json:"id"`
	PitId        int   `json:"dungeons_id"`
	PitType      int   `json:"dungeons_type"`
	PassTime     int   `json:"pass_time"`
	JoinNeed     int   `json:"join_need"`
	NeedItemType int   `json:"need_item_type"`
	NeedItem     int   `json:"need_item"`
	OpenTime     int64 `json:"open_time"`
	OpenDuration int64 `json:"open_duration"`
}

type PitBuff struct {
	Id     int   `json:"id"`
	BuffId []int `json:"buff_id"`
}

type PitBox struct {
	Id      int   `json:"id"`
	BoxId   []int `json:"box_id"`
	ItemId  []int `json:"item_id"`
	ItemNum []int `json:"item_number"`
}

type PitMonster struct {
	Id         int   `json:"id"`
	OptionInfo []int `json:"option_info"`
}

type NewPitConfig struct {
	Id                int      `json:"id"`
	MapId             int      `json:"mapid"`
	Line              int      `json:"line"`
	Row               int      `json:"row"`
	NumberOfPlies     int      `json:"numberofplies"`
	Difficulty        int      `json:"difficulty"`
	Element           int      `json:"element"`
	UpperParam        int      `json:"upper_param"`
	LimitParam        int      `json:"limit_param"`
	LimitRobot        int64    `json:"limit_Robot"`
	MinimumLimit      int      `json:"minimum_limit"`
	LimitMin          int      `json:"limit_min"`
	UpperMin          int      `json:"upper_min"`
	LimitRobotMin     int64    `json:"limitrobot_min"`
	MinimumMin        int      `json:"minimum_min"`
	Prize             []int    `json:"prize"`
	Num               []int    `json:"num"`
	Relique           []int    `json:"relique_p"`
	BraveItem         int      `json:"brave_item"`
	BraveNum          int      `json:"brave_num"`
	CartQuality       int      `json:"cart_quality"`
	CartLv            int      `json:"cart_lv"`
	Lottery           []int    `json:"lottery"`
	Param             []int    `json:"param"`
	TreasureItem      []int    `json:"treasure_item"`
	TreasureNum       []int    `json:"treasure_num"`
	Select            []string `json:"select"`
	FirstRobotQuality string   `json:"first_robot_quality"`
	FirstRobotLv      string   `json:"first_robot_lv"`
	FirstRobotGroup   string   `json:"first_robot_group"`
	VShop_p           []int    `json:"vshop_p"`
}

type NewPitExtraReward struct {
	Id      int   `json:"id"`
	Lottery []int `json:"lottery"`
}

type NewPitParam struct {
	Id           int   `json:"id"`
	Quality      int   `json:"quality"`
	Item         []int `json:"item"`
	QualityParam []int `json:"quality_param"`
	LvParam      []int `json:"lv_param"`
}

type NewPitDifficulty struct {
	Id                    int   `json:"id"`
	Difficulty            int   `json:"difficulty"`
	NumberOfPlies         int   `json:"numberofplies"`
	LimitParamReduction   int   `json:"limit_param_reduction"`
	UpperParamReduction   int   `json:"upper_param_reduction"`
	LimitRobotReduction   int64 `json:"limit_robot_reduction"`
	MinimumLimitReduction int   `json:"minimum_limit_reduction"`
	LimitParamAdd         int   `json:"limit_param_add"`
	UpperParamAdd         int   `json:"upper_param_add"`
	LimitRobotAdd         int64 `json:"limit_robot_add"`
	MinimumLimitAdd       int   `json:"minimum_limit_add"`
}

type NewPitRelique struct {
	Id          int `json:"id"`
	Quality     int `json:"quality"`
	OnlyOwned   int `json:"only_owned"`
	LevelLimit  int `json:"level_limit"`
	Interesting int `json:"interesting"`
}

type NewPitTreasureCave struct {
	Id        int `json:"id"`
	Lv        int `json:"lv"`
	MonsterId int `json:"monsterid"`
}

type NewPitRobotConfig struct {
	Id              int   `json:"id"`
	HeroId          int   `json:"hero_id"`
	FirstRobotGroup int   `json:"first_robot_group"`
	RobotGroup      int   `json:"robot_group"`
	Attribute       int   `json:"attribute"`
	HeroStar        int   `json:"hero_star"`
	ArmsId          []int `json:"arms_id"`
	ArmsRange       []int `json:"arms_range"`
	HeadId          []int `json:"head_id"`
	HeadRange       []int `json:"head_range"`
	BodyId          []int `json:"body_id"`
	BodyRange       []int `json:"body_range"`
	ShoesId         []int `json:"shoes_id"`
	ShoesRange      []int `json:"shoes_range"`
	ArtifactType    int   `json:"artifact_type"`
	ArtifactId      int   `json:"artifact_id"`
	ArtifactLv      int   `json:"artifact_lv"`
	ExclusiveEquip  int   `json:"exclusiveequip"`
	ExclusiveLv     int   `json:"exclusivelv"`
	Cart            int   `json:"cart"`
	GroupAttribute  int   `json:"group_attribute"`
	InitialQuality  int   `json:"initial_quality"`
	CartTalentNum   int   `json:"cart_talent_num"`
	NpcTalentNum    int   `json:"npc_talent_num"`
}

type NewPitRobotExclusiveConfig struct {
	Id              int   `json:"id"`
	ExclusiveLv     []int `json:"exclusive_lv"`
	ExclusiveNum    int   `json:"exclusive_num"`
	CartExclusiveLv int   `json:"cart_exclusivelv"`
}

type NewPitRobotQuality struct {
	Id         int   `json:"id"`
	MeanStar1  int   `json:"mean_star1"`
	MeanStar2  int   `json:"mean_star2"`
	Star       []int `json:"star"`
	CommonCart int   `json:"common_cart"` //过滤用
}

type NewPitRobotMonsterQuality struct {
	Id        int   `json:"id"`
	MeanStar1 int   `json:"mean_lv1"`
	MeanStar2 int   `json:"mean_lv2"`
	HeroLimit []int `json:"hero_limit"`
	HeroUpper []int `json:"hero_upper"`
	HeroGroup []int `json:"hero_group"`
}

type NewPitRobotMonsterLv struct {
	Num      int   `json:"num"`
	StarMin  int   `json:"star_min"`
	StarMax  int   `json:"star_max"`
	FightMin int64 `json:"fight_min"`
	FightMax int64 `json:"fight_max"`
	LvMin    int   `json:"lv_min"`
	LvMax    int   `json:"lv_max"`
}

type NewPitRobotAttr struct {
	Id          int     `json:"id"`
	Group       int     `json:"group"`
	Lv          int     `json:"lv"`
	GrowthType  []int   `json:"growth_type"`
	GrowthValue []int64 `json:"growth_value"`
}

//时光之巅
type InstanceConfig struct {
	MapId      int `json:"mapid"`
	LevelJudge int `json:"leveljudge"`
	Condition  int `json:"condition"`
	Judge      int `json:"judge"`
}

type InstanceBox struct {
	BoxId   int   `json:"boxid"`
	GroupId int   `json:"groupid"`
	Type    int   `json:"type"`
	Item    []int `json:"item"`
	Num     []int `json:"num"`
}

type InstanceThing struct {
	Id             int    `json:"id"`
	Row            int    `json:"x"`
	Col            int    `json:"y"`
	Event          int    `json:"event"`
	Type           int    `json:"type"`
	RemoveType     int    `json:"removetype"`
	RemoveRelation string `json:"removerelation"`
	Remove         string `json:"remove"`
	Establish      string `json:"establish"`
	Dispel         string `json:"dispel"`
}

//! 羁绊系统
type EntanglementConfig struct {
	FateNum    int     `json:"fate_num"`
	Group      int     `json:"group"`
	HeroId     []int   `json:"hero_id"`
	MinQuality int     `json:"min_quality"`
	HeroNum    int     `json:"hero_num"`
	BaseType   []int   `json:"base_type"`
	BaseValue  []int64 `json:"base_value"`
}

//! 羁绊系统
type EntanglementFate struct {
	Group    int                     `json:"group"`
	HeroId   []int                   `json:"hero_id"`
	HeroNum  int                     `json:"hero_num"`
	Property []*EntanglementProperty `json:"property"`
}

//! 羁绊系统
type EntanglementProperty struct {
	FateNum    int     `json:"fate_num"`
	MinQuality int     `json:"min_quality"`
	BaseType   []int   `json:"base_type"`
	BaseValue  []int64 `json:"base_value"`
}

//! 悬赏
type RewardForbarConfig struct {
	ID          int   `json:"id"`
	Group       int   `json:"group"`
	Color       int   `json:"color"`
	Star        int   `json:"star"`
	IsTeam      int   `json:"isteam"`
	SepndTime   int   `json:"time"`
	ElapsedTime int   `json:"task_elapsed_time"`
	NeedCamp    []int `json:"needtocamp"`
	NeedStar    int   `json:"needtostar"`
	StarNum     int   `json:"starnum"`
}

//! 悬赏升级
type RewardForbarLvUpConfig struct {
	LV             int   `json:"lv"`
	Renovatestar   []int `json:"renovatestar"`
	Renovatepro    []int `json:"renovatepro"`
	Persontaskstar int   `json:"persontaskstar"`
	Persontasknum  int   `json:"persontasknum"`
	Teamtaskstar   int   `json:"teamtaskstar"`
	Teamtasknum    int   `json:"teamtasknum"`
}

type RewardForbarPrize struct {
	ID       int `json:"id"`
	Group    int `json:"group"`
	Isteam   int `json:"isteam"`
	P        int `json:"p"`
	Prize    int `json:"prize"`
	Prizenum int `json:"prizenum"`
}

//! 凯旋丰碑
type RankListIntegral struct {
	ID          int `json:"id"`
	HeroQuality int `json:"hero_quality"`
	RankValue   int `json:"rank_value"`
}

// 排行任务
type RankTaskConfig struct {
	Id        int    `json:"id"`
	Type      int    `json:"type"`
	Sort      int    `json:"sort"`
	TaskTypes int    `json:"tasktypes"`
	Ns        []int  `json:"n"`
	Items     []int  `json:"item"`
	Nums      []int  `json:"num"`
	Dec       string `json:"dec"`
	Txt       string `json:"txt"`
	Judge     int    `json:"judge"`
}

type ResonanceCrystalconfig struct {
	Id           int `json:"id"`
	MaxLevel     int `json:"max_level"`
	CrystalLevel int `json:"crystal_level"`
}

type ArenaRewardConfig struct {
	Id    int   `json:"id"`
	Type  int   `json:"type"`
	Max   int   `json:"max"`
	Min   int   `json:"min"`
	Items []int `json:"item"`
	Nums  []int `json:"num"`
}

type ArenaParameterConfig struct {
	Id            int64 `json:"id"`
	Matemax       int   `json:"matemax"`
	Matemin       int   `json:"matemin"`
	Expected      int   `json:"expected"`
	Winparameter  int   `json:"winparameter"`
	Loseparameter int   `json:"loseparameter"`
	Integral      int   `json:"integral"`
	Drop          int   `json:"drop"`
}

// 军团狩猎
type UnionHuntConfig struct {
	Id    int `json:"id"`
	Level int `json:"level"`
	Cost  int `json:"cost"`
	Time  int `json:"time"`
	Group int `json:"group"`
}

// 军团狩猎掉落
type UnionHuntDropConfig struct {
	Id            int   `json:"id"`
	Group         int   `json:"group"`
	Lv            int   `json:"lv"`
	Hp            int64 `json:"hp"`
	Drop          []int `json:"drop"`
	Level         []int `json:"Level"`
	Personalitem  []int `json:"personalitem"`
	Personalnum   []int `json:"personalnum"`
	Guilditem     []int `json:"guilditem"`
	Guildnum      []int `json:"guildnum"`
	DiamondTime   int   `json:"diamond_time"`
	Diamond       []int `json:"diamond"`
	DiamondChance []int `json:"diamond_chance"`
}

// 高阶竞技场
type ArenaSpecialClass struct {
	Id       int    `json:"id"`
	Class    int    `json:"class"`
	Dan      int    `json:"dan"`
	Icon     string `json:"icon"`
	Name     string `json:"name"`
	Ranking  int    `json:"ranking"`
	Capacity int    `json:"capacity"`
	Integral int    `json:"integral"`
	Income   int    `json:"income"`
	Limit    int    `json:"limit"`
}

// 英雄皮肤
type HeroSkin struct {
	ID    int `json:"id"`
	Hero  int `json:"hero"`
	ModId int `json:"modid"`
	Item  int `json:"item"`
	Num   int `json:"num"`
}

type ActivityBuyLimit struct {
	ID         int   `json:"id"`
	Type       int   `json:"type"`
	OpenTime   int   `json:"open_time"`
	OverTime   int   `json:"over_time"`
	TaskTypes  int   `json:"tasktypes"`
	N          []int `json:"n"`
	Group      int   `json:"group"`
	LimitTime  int64 `json:"limit_time"`
	Subtype    int   `json:"subtype"`
	TriggerNum int   `json:"trigger_num"`
}

type ActivityBuyItem struct {
	ID          int    `json:"id"`
	Group       int    `json:"group"`
	Money       int    `json:"money"`
	Items       []int  `json:"item"`
	Nums        []int  `json:"num"`
	MoneyID     int    `json:"money_id"`
	MoneyReduce int    `json:"money_reduce"`
	Sale        int    `json:"sale"`
	RebateShow  string `json:"rebate_show"`
}

type TreeLevel struct {
	ID                int     `json:"id"`
	TreeLevel         int     `json:"tree_level"`
	ItemID            int     `json:"item_id"`
	ItemNum           int     `json:"item_num"`
	ProfessionalLevel int     `json:"professional_level"`
	Value             []int   `json:"value"`
	Num               []int64 `json:"num"`
	Increase          []int   `json:"increase"`
}

type TreeProfessional struct {
	ID       int     `json:"id"`
	Type     int     `json:"type"`
	Level    int     `json:"level"`
	ItemID   int     `json:"item_id"`
	ItemNum  int     `json:"item_num"`
	Skill    []int   `json:"open_skill"`
	Value    []int   `json:"value"`
	Num      []int64 `json:"num"`
	Increase []int   `json:"increase"`
}

//type TreeConfig struct {
//	ID      int `json:"id"`
//	Quality int `json:"first_quality"`
//	ItemID  int `json:"itemid"`
//	ItemNum int `json:"num"`
//	Sum     int `json:"sum_num"`
//}

type InterstellarConfig struct {
	Id        int   `json:"id"`        //编号
	Galaxy    int   `json:"galaxy"`    //所属星系
	TaskTypes []int `json:"tasktypes"` //解锁条件
	N         []int `json:"n"`         //N1
	M         []int `json:"m"`         //N2
	P         []int `json:"p"`         //N3
	Q         []int `json:"q"`         //N4
}

type InterstellarHangup struct {
	PrivileGeNum     int   `json:"privilegenum"`     //科技特权编号
	PrivileGeType    int   `json:"privilegetype"`    //科技特权类型
	Type             int   `json:"type"`             //类型
	InterstellarTime int64 `json:"interstellartime"` //时间
	Item             int   `json:"item"`             //道具
	Num              int   `json:"num"`              //数量
	InterstellarDrop int   `json:"interstellardrop"` //掉落库
	Judge            []int `json:"judge"`            //条件判断
	ChangeDrop       []int `json:"changedrop"`       //替换库
	AddItem          []int `json:"additem"`          //额外一次奖励
	AddNum           []int `json:"addnum"`           //额外一次奖励
}

type InterstellarWar struct {
	Id            int   `json:"id"`            //编号
	Nebula        int   `json:"nebula"`        //所属星云
	Front         int   `json:"front"`         //前置星点
	TaskTypes     int   `json:"tasktypes"`     //判断
	Ns            []int `json:"n"`             //目标
	Type          int   `json:"type"`          //科技处理类型0取最高  1百分比累加  2值累加
	PrivileGeType int   `json:"privilegetype"` //激活特权类型
	PrivileGeNum  int   `json:"privilegenum"`  //激活特权效果
}

type InterstellarBox struct {
	BoxId   int   `json:"boxid"`   //箱子唯一编号
	GroupId int   `json:"groupid"` //箱子组编号
	Type    int   `json:"type"`    //宝箱类型
	Item    []int `json:"item"`    //奖励道具
	Num     []int `json:"num"`     //奖励数量
}

type ActivityBossRankConfig struct {
	Id              int    `json:"id"`               //编号
	ActivityType    int    `json:"activity_type"`    //关联活动
	ActivityPeriods int    `json:"activity_periods"` //关联期数
	Subsection      int    `json:"subsection"`       //段位
	Ranking         int    `json:"ranking"`          //分段排名
	Name            string `json:"name"`             //名称
	Section         int64  `json:"section"`          //下区间
	Contain         int    `json:"contain"`          //容纳人数
	Item            []int  `json:"item"`             //奖励物品
	Num             []int  `json:"num"`              //奖励数量
}

type ActivityBossConfig struct {
	Id              int    `json:"id"`               //编号
	ActivityType    int    `json:"activity_type"`    //关联活动
	ActivityPeriods int    `json:"activity_periods"` //关联期数
	BossId          int    `json:"bossid"`           //
	Name            string `json:"name"`             //名称
	Level           int    `json:"level"`            //下区间
	Times           int    `json:"times"`            //每天免费次数
	Skill           []int  `json:"skill"`            //
	Hero            []int  `json:"hero"`             //
	Position        []int  `json:"position"`         //
	Items           []int  `json:"item"`             //
	HealthUnit      int    `json:"healthunit"`       //计量单位【万分比】
	GetItem         int    `json:"getitem"`          //单位道具
	GetNum          int    `json:"getnum"`           //血条数量
}

type ActivityBossTargetConfig struct {
	Id              int   `json:"id"`               //编号
	ActivityType    int   `json:"activity_type"`    //关联活动
	ActivityPeriods int   `json:"activity_periods"` //关联期数
	TaskId          int   `json:"taskid"`           //
	TaskTypes       int   `json:"tasktypes"`        //名称
	Ns              []int `json:"n"`                //
	Item            int   `json:"item"`             //
	Num             int   `json:"num"`              //
}

type ActivityBossExchangeConfig struct {
	Id              int   `json:"id"`               //编号
	ActivityType    int   `json:"activity_type"`    //关联活动
	ActivityPeriods int   `json:"activity_periods"` //关联期数
	Item            int   `json:"item"`             //兑换
	Num             int   `json:"num"`              //
	NeedItem        []int `json:"needitem"`         //消耗
	NeedNum         []int `json:"neednum"`          //
	Frequency       int   `json:"frequency"`        //次数
}

type HeroGrowConfig struct {
	Id              int   `json:"id"`
	ActivityType    int   `json:"activity_type"`
	ActivityPeriods int   `json:"activity_periods"`
	Hero            int   `json:"hero"`
	NeedQuality     int   `json:"needquality"`
	Limit           int   `json:"limit"`
	TaskTypes       int   `json:"tasktypes"`
	Ns              []int `json:"n"`
	Item            []int `json:"item"`
	Num             []int `json:"num"`
}

type CrossArenaRewardConfig struct {
	Id         int   `json:"id"`
	Type       int   `json:"type"`
	Subsection int   `json:"subsection"`
	Class      int   `json:"class"`
	Item       []int `json:"item"`
	Num        []int `json:"num"`
}

type CrossArenaSubsection struct {
	Id         int    `json:"id"`
	Subsection int    `json:"subsection"`
	Class      int    `json:"class"`
	Name       string `json:"name"`
	Item       []int  `json:"item"`
	Num        []int  `json:"num"`
}

type CrossArena3V3RewardConfig struct {
	Id         int   `json:"id"`
	Type       int   `json:"type"`
	Subsection int   `json:"subsection"`
	Class      int   `json:"class"`
	Item       []int `json:"item"`
	Num        []int `json:"num"`
}

type CrossArena3V3Subsection struct {
	Id         int    `json:"id"`
	Subsection int    `json:"subsection"`
	Class      int    `json:"class"`
	Name       string `json:"name"`
	Item       []int  `json:"item"`
	Num        []int  `json:"num"`
}

// 新天赋系统 旧系统废弃
type StageTalentConfig struct {
	ID    int   `json:"talent_id"`
	Group int   `json:"group"`
	Open  int   `json:"talent_star_open"`
	Index int   `json:"index"`
	Skill []int `json:"skill_id"`
	Type  int   `json:"type"`
	Value int   `json:"value"`
}

type LotteryDrawConfig struct {
	Id        int `json:"id"`
	Type      int `json:"type"`
	Group     int `json:"group"`
	Items     int `json:"item"`
	Nums      int `json:"num"`
	Limit     int `json:"limit"`
	Layer     int `json:"layer"`
	Change    int `json:"change"`
	Lucky     int `json:"lucky"`
	Weight    int `json:"weight"`
	MinLucky  int `json:"minlucky"`
	MinWeight int `json:"minweight"`
	MaxLucky  int `json:"maxlucky"`
	MaxWeight int `json:"maxweight"`
	Notice    int `json:"notice"`
}

type WorldLvLevelConfig struct {
	Id         int     `json:"id"`
	WorldLv    int     `json:"worldlv"`
	NpcLv      int     `json:"npclv"`
	RageRate   int     `json:"ragerate"`
	CableRange int     `json:"cablerange"`
	BaseType   []int   `json:"base_type"`
	BaseValue  []int64 `json:"base_value"`
}

type WorldLvTpyeConfig struct {
	Id         int     `json:"id"`
	WorldLv    int     `json:"worldlv"`
	NpcId      int     `json:"npcid"`
	NpclvTable int     `json:"npclvtable"`
	HeroNpc    int     `json:"heronpc"`
	HeroEquip  int     `json:"heroequip"`
	BaseType   []int   `json:"base_type"`
	BaseValue  []int64 `json:"base_value"`
}

type WorldShowBossConfig struct {
	Id          int `json:"id"`
	Crystallv   int `json:"crystallv"`
	ShowHeroqua int `json:"showheroqua"`
	ShowEquip   int `json:"showequip"`
}

type HonourShopConfig struct {
	Id           int   `json:"id"`
	ParentLabel  int   `json:"parent_label"`
	SubTab       int   `json:"subtab"`
	Grid         int   `json:"grid"`
	Index        int   `json:"index"`
	LevelStart   int   `json:"level_start"`
	LevelEnd     int   `json:"level_end"`
	TaskTypes    int   `json:"tasktypes"`
	Ns           []int `json:"n"`
	ItemId       int   `json:"itemid"`
	ItemNumber   int   `json:"itemnumber"`
	CostItem     []int `json:"cost_item"`
	CostNum      []int `json:"cost_num"`
	PurchaseType int   `json:"purchase_type"`
	PurchaseNum  int   `json:"purchase_num"`
}

type RankRewardConfig struct {
	Id           int   `json:"id"`            //编号
	ActivityType int   `json:"activity_type"` //关联活动
	Group        int   `json:"group"`         //掉落组
	RankHigh     int   `json:"rankhigh"`      //排行区间高
	RankLow      int   `json:"ranklow"`       //排行区间低
	Item         []int `json:"item"`          //奖励物品
	Num          []int `json:"num"`           //奖励数量
}
