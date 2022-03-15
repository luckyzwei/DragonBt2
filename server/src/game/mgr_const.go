package game

// 游戏日志
const (
	LOG_FIGHT_FINALEX      = 26  // 决战
	LOG_USER_PASS_OK       = 100 // 关卡通过
	LOG_USER_TITLE_ACTIVE  = 101 // 威名-称号激活
	LOG_USER_RECHARGE      = 102 // 充值
	LOG_USER_VIP_UP        = 103 // VIP升级
	LOG_USER_LEVEL_UP      = 104 // 玩家升级
	LOG_BEAUTY_UP_TREASURE = 208 // 圣物-宝物-升级
	LOG_BEAUTY_UP_LEVEL    = 209 // 圣物-升级(激活，升阶）
	//LOG_HERO_ACTIVE        = 211 // 激活武将

	LOG_FIGHT_SOLO   = 301 // 国战-单挑
	LOG_FIGHT_FINAL  = 302 // 决战
	LOG_FIGHT_SERIES = 303 // 国战单人玩法

	LOG_PVP_ARMSARENA_BU = 401 // 天下会武-补-挑战
	LOG_PVP_KING_FIGHT   = 404 // 王权战斗
	LOG_PVP_OFFICE_FIGHT = 405 // 官职战斗
	LOG_PVP_EXPEDITION   = 406 // 远征

	LOG_EVENT_VISIT_CELEBRITY = 600 // 事件-名士寻访
	LOG_EVENT_PEOPLE_FEELING  = 601 // 民生民情
	LOG_EVENT_PASS_HERO       = 602 // 过关斩将
	LOG_EVENT_TREASURE_MAP    = 603 // 藏宝图

	LOG_EVENT_SummonHorse_Normal = 701 // 普通招募
	LOG_EVENT_SummonHorse_Senior = 702 // 高级招募
	LOG_EVENT_Discernhorse       = 703 // 相马
	LOG_EVENT_Combinehorse       = 704 // 合成马
	LOG_EVENT_Uphorse_4          = 705 // 马升4星
	LOG_EVENT_Uphorse_5          = 706 // 马升5星
	LOG_EVENT_Awakehorse         = 707 // 马觉醒
	LOG_EVENT_Decomposehorse     = 708 // 分解战马
	LOG_EVENT_Decomposesoul      = 709 // 分解马魂
	LOG_EVENT_Exchangesoul       = 710 // 马魂兑换
	LOG_EVENT_Uphorsesoul        = 711 // 马魂强化
	LOG_EVENT_GET_HORSE          = 712 // 获得战马

	LOG_EVENT_FlushmiliTask  = 801 // 刷新
	LOG_EVENT_ExectemiliTask = 802 // 派遣执行
	LOG_EVENT_AwardmiliTask  = 803 // 派遣领取

	LOG_GENERAL_HERO      = 1001 //  限时神将日志
	LOG_ACT_DIAL          = 1003 //  活动转盘日志
	LOG_ACT_DRAW          = 1004 // 活动翻牌
	LOG_ACT_LUCKEGG       = 1005 // 金蛋活动
	LOG_ACT_LUCKSTART     = 1006 // 开工福利
	LOG_ACT_DAILYRECHARGE = 1007 // 连续充值

	LOG_TREASURE = 2001 // 宝物
)

const (
	LOG_FIGHT                  = 1
	LOG_SWEEP                  = 2
	LOG_SUCCESS                = 1
	LOG_FAIL                   = 2
	LOG_PASS_TYPE_FIGHT        = 1
	LOG_PASS_TYPE_PASS_WIN     = 2
	LOG_PASS_TYPE_SEARCH       = 3
	LOG_PASS_TYPE_SWEEP        = 4
	LOG_MONEY_TASK_ITEM        = 1
	LOG_MONEY_TASK_GEM         = 2
	LOG_MONEY_TASK_REPLACEITEM = 3
	LOG_GVE_BUILD              = 1
	LOG_GVE_DRAGON             = 2
	LOG_SMELT_FREE             = 1
	LOG_SMELT_GEM              = 2
	LOG_SMELT_GEM_TEN          = 3
)

// 新日志
const (
	LOG_USER_LOGIN         = 100 //! 角色登录
	LOG_USER_CREATE_PLAYER = 101 //! 创角色
	LOG_GUIDE_START        = 102 //! 引导
	LOG_GUIDE_STORY        = 103 //! 故事剧情
	LOG_USER_SWITCH_SERVER = 104 //! 玩家切换服务器
	LOG_USER_AUTO_LOGIN    = 105 //! 断线重连

	LOG_PLAYER_UP_LEVEL        = 201 //! 升级
	LOG_PLAYER_CHANGE_NAME     = 202 //! 改名
	LOG_PLAYER_CHANGE_ICON     = 203 //! 改头像
	LOG_PLAYER_CHANGE_PORTRAIT = 204 //! 改头像框

	LOG_PASS_START_NORMAL  = 301 //! 普通关卡挑战
	LOG_PASS_FINISH_NORMAL = 302 //! 普通关卡通关

	LOG_ONHOOK_AWARD = 401 //! 领取挂机奖励
	LOG_ONHOOK_FAST  = 402 //! 快速挂机

	LOG_HERO_LEVE_UP    = 501 //!英雄升级
	LOG_HERO_LEVE_UP_TO = 502 //!英雄一键升级

	LOG_HANDBOOK_ACTIVATE = 601 //!激活英雄图鉴
	LOG_HANDBOOK_AWARD    = 602 //!领取图鉴奖励

	LOG_SKIN_ACTIVITE = 801 //!激活皮肤
	LOG_SKIN_SET      = 802 //!更换皮肤

	LOG_HERO_LOCK   = 901 //!锁定英雄
	LOG_HERO_UNLOCK = 902 //!解锁英雄

	LOG_EQUIP_WEAR_AUTO     = 1001 // 一键穿戴装备
	LOG_EQUIP_OFF_AUTO      = 1002 // 一键卸下装备
	LOG_EQUIP_LEVEL_UP      = 1003 // 装备强化
	LOG_EQUIP_UPGRADE       = 1004 // 装备升阶
	LOG_EQUIP_RECAST        = 1005 // 装备重铸
	LOG_EQUIP_RECAST_CHOOSE = 1006 // 装备重铸保存
	LOG_EQUIP_WEAR          = 1007 // 穿戴装备
	LOG_EQUIP_CHANGE        = 1008 // 装备替换
	LOG_EQUIP_OFF           = 1009 // 装备卸下

	LOG_EXCLUSIVE_UNLOCK   = 1101 // 专属激活
	LOG_EXCLUSIVE_LEVEL_UP = 1102 // 专属强化
	LOG_EXCLUSIVE_RESET    = 1103 // 英雄降阶专属回退
	LOG_ARTIFACT_CHANGE    = 1105 // 神器替换

	LOG_ARTIFACT_LEVEL_UP = 1201 //!神器强化

	LOG_HERO_UP_STAR = 1301 //!英雄升阶
	LOG_HERO_REBORN  = 1302 //!英雄重塑
	LOG_HERO_BACK    = 1303 //!英雄降阶
	LOG_HERO_FIRE    = 1304 //!英雄遣退

	LOG_SUPPORT_HERO_SET    = 1401 //!设置我的派遣英雄
	LOG_SUPPORT_HERO_CANCEL = 1402 //!解除我的派遣英雄

	LOG_ENTANGLEMENT_SET = 1501 //!设置羁绊英雄

	LOG_DAILY_TASK          = 1601 //!领取日常任务奖励
	LOG_DAILY_TASK_LIVENESS = 1602 //!领取日常活跃度宝箱
	LOG_WEEK_TASK           = 1603 //!领取周常任务奖励
	LOG_WEEK_TASK_LIVENESS  = 1604 //!领取周常活跃度宝箱
	LOG_GROWTH_TASK         = 1605 //!领取主线任务奖励

	LOG_BAG_USE      = 1701 //!背包使用
	LOG_BAG_COMPOUND = 1702 //!背包合成

	LOG_RESONANCE_CRYSTAL_ADD    = 1801 // 解锁法阵格子
	LOG_RESONANCE_CRYSTAL_SET    = 1802 // 设置法阵列表英雄
	LOG_RESONANCE_CRYSTAL_CANCEL = 1803 // 解除法阵列表英雄

	LOG_UNION_CREATE       = 1901 //创建公会
	LOG_UNION_DISSOLVE     = 1902 //解散公会
	LOG_UNION_ALERT_NAME   = 1903 //公会改名
	LOG_UNION_CHANGE_ICON  = 1904 //公会改旗帜
	LOG_UNION_SEND_MAIL    = 1905 //发送公会全员邮件
	LOG_UNION_CHANGE_BOARD = 1906 //工会改宣言

	LOG_UNION_HUNTER_OPEN  = 2001 // 激活公会挑战关卡
	LOG_UNION_HUNTER_END   = 2002 // 公会挑战
	LOG_UNION_HUNTER_SWEEP = 2003 // 狩猎扫荡

	LOG_HERO_FIND_GEM_ONE     = 2101 //!单次高级招募
	LOG_HERO_FIND_GEM_TEN     = 2102 //!十连高级招募
	LOG_HERO_FIND_FRIEND_ONE  = 2103 //!单次友情招募
	LOG_HERO_FIND_FRIEND_TEN  = 2104 //!十连友情招募
	LOG_HERO_FIND_CAMP_ONE    = 2105 //!单次阵营招募
	LOG_HERO_FIND_CAMP_TEN    = 2106 //!十连阵营招募
	LOG_HERO_FIND_CAMP_CHOOSE = 2107 //!自选阵营招募
	LOG_HERO_FIND_CAMP_OPEN   = 2108 //!阵营招募激活阵营
	LOG_HERO_FIND_SELF_ONE    = 2109 //!单次自选招募
	LOG_HERO_FIND_SELF_TEN    = 2110 //!十连自选招募
	LOG_HERO_FIND_SELF_CHANGE = 2111 //!更换目标英雄

	LOG_HERO_FIND_WISH_CANCEL = 2201 //解除心愿单英雄
	LOG_HERO_FIND_WISH_SET    = 2202 //设置心愿单英雄

	LOG_SHOP_BUY_GOODS = 2301 //!商店购买道具
	LOG_SHOP_REFRESH   = 2302 //!刷新商店商品

	LOG_RANK_GET_AWARD = 2401 //!领取排行榜奖励
	LOG_RANK_LOOK      = 2402 //!排行榜查看用户信息

	LOG_GET_MAIL_SINGLE = 2501 //!领取邮件
	LOG_GET_MAIL_ALL    = 2502 //!一键领取邮件

	LOG_FRIEND_POWER_SEND_GET = 2601 //!好友一键领取和赠送
	LOG_FRIEND_POWER_SEND     = 2602 //!好友赠送
	LOG_FRIEND_POWER_GET      = 2603 //!好友领取
	LOG_FRIEND_BATTLE         = 2604 //!好友切磋
	LOG_FRIEND_APPLY          = 2605 //!好友申请
	LOG_FRIEND_HIRE_APPLY     = 2606 //!好友佣兵申请
	LOG_FRIEND_ADD            = 2607 //!添加好友
	LOG_FRIEND_REMOVE         = 2608 //!删除好友

	LOG_NEWPIT_BATTLE    = 2701 //!地牢战斗
	LOG_NEWPIT_SOUL_CART = 2702 //!地牢号令火炬
	LOG_NEWPIT_STRING    = 2703 //!地牢复苏清泉
	LOG_NEWPIT_REBORN    = 2704 //!地牢神使
	LOG_NEWPIT_EVIL_CART = 2705 //!地牢试炼火炬
	LOG_NEWPIT_TREASURE  = 2706 //!地牢秘宝守卫
	LOG_NEWPIT_REWARD    = 2707 //!领取通关地牢整层奖励
	LOG_NEWPIT_ITEM      = 2708 //!地牢使用复活道具
	LOG_NEWPIT_NEW_HARD  = 2709 //!进入地牢困难模式
	LOG_NEWPIT_NEW_EASY  = 2710 //!进入地牢简单模式

	LOG_ARENA_FIGHT       = 2801 //!竞技场战斗
	LOG_ARENA_FIGHT_BACK  = 2802 //!竞技场反击战斗
	LOG_ARENA_LOOK        = 2803 //!竞技场查看用户信息
	LOG_ARENA_BATTLE_INFO = 2804 //!查看竞技场战报

	LOG_ARENA_SPECIAL_FIGHT       = 2901 //!高阶竞技场战斗
	LOG_ARENA_SPECIAL_GET_AWARD   = 2902 //!领取高阶竞技场积累奖励
	LOG_ARENA_SPECIAL_LOOK        = 2903 //!高阶竞技场查看用户信息
	LOG_ARENA_SPECIAL_BATTLE_INFO = 2904 //!查看高阶竞技场战报

	LOG_TOWER_FINISH_NORMAL = 3001 //!试炼之塔普通塔战斗
	LOG_TOWER_FINISH_CAMP   = 3002 //!试炼之塔种族塔战斗
	LOG_TOWER_BATTLE_INFO   = 3003 //!查看试炼之塔战报

	LOG_CONSUMER_FIGHT  = 3011 //!诸神黄昏战斗
	LOG_CONSUMER_CHANGE = 3012 //!诸神黄昏排名变化

	LOG_REWARD_SET_NO_TEAM     = 3101 //!个人赏金任务派遣
	LOG_REWARD_SET_TEAM        = 3102 //!团队赏金任务派遣
	LOG_REWARD_SET_REFRESH     = 3103 //!刷新赏金任务
	LOG_REWARD_GET_NO_TEAM     = 3104 //!领取个人赏金任务
	LOG_REWARD_GET_TEAM        = 3105 //!领取团队赏金任务
	LOG_REWARD_SET_NO_TEAM_ALL = 3106 //!一键派遣个人赏金任务
	LOG_REWARD_SET_TEAM_ALL    = 3107 //!一键派遣团队赏金任务
	LOG_REWARD_GET_NO_TEAM_ALL = 3108 //!一键领取个人赏金任务
	LOG_REWARD_GET_TEAM_ALL    = 3109 //!一键领取团队赏金任务

	LOG_SHOP_BUY = 3201 //!商店购买
	//LOG_SHOP_REFRESH = 3202 //!商店刷新

	LOG_FIRST_GET   = 3401 //!领取首充奖励
	LOG_HERO_ACTIVE = 3402 // 英雄激活
	LOG_HERO_DELETE = 3403 // 消耗英雄

	LOG_MONTH_CARD           = 3501 //购买普通月卡
	LOG_MONTH_CARD_HIGH      = 3502 //购买高级月卡
	LOG_MONTH_CARD_GET       = 3503 //领取普通月卡奖励
	LOG_MONTH_CARD_HIGH_GET  = 3504 //领取高级月卡奖励
	LOG_MONTH_CARD_SCORE_GET = 3505 //领取月卡积分奖励

	LOG_WARORDER_BUY_1 = 3601 //!购买皇家犒赏令
	LOG_WARORDER_BUY_2 = 3602 //!购买勇士犒赏令
	LOG_WARORDER_GET_1 = 3603 //!领取皇家犒赏令奖励
	LOG_WARORDER_GET_2 = 3604 //!领取勇士犒赏令奖励

	LOG_HERO_TALENT_LOOT_GEM_ONE  = 3605 //!神格高级抽取
	LOG_HERO_TALENT_LOOT_GOLD_TEN = 3606 //!神格普通十连
	LOG_HERO_TALENT_LOOT_GEM_TEN  = 3607 //!神格高级十连

	LOG_SPECIAL_PURCHASE_ACTIVATE = 3701 //!触发限时礼包
	LOG_SPECIAL_PURCHASE_BUY      = 3702 //!购买限时礼包

	LOG_ACTIVITY_LOGIN_TIMES = 3801 //!连续登录

	LOG_ACTIVITY_SEVEN_1 = 3901 //!新兵训练1
	LOG_ACTIVITY_SEVEN_2 = 3902 //!新兵训练2
	LOG_ACTIVITY_SEVEN_3 = 3903 //!新兵训练3

	LOG_VIP_BOX     = 4001 //!贵族礼包
	LOG_BUY_VIP_BOX = 4002 //!购买vip礼包

	LOG_CDKEY_REWARD = 4101 //!领取激活码奖励 (游戏服暂无法判断)
	LOG_GEM_WEAR     = 4102 //!宝石镶嵌
	LOG_GEM_TAKEOFF  = 4103 //!宝石取消镶嵌

	LOG_INSTANCE_OPEN        = 4201 //!开启副本
	LOG_INSTANCE_CLOSE       = 4202 //!关闭副本
	LOG_INSTANCE_GET_REWARD  = 4203 //!获得宝箱奖励
	LOG_INSTANCE_BATTLE_WIN  = 4204 //!战斗胜利
	LOG_INSTANCE_BATTLE_LOSE = 4205 //!战斗失败

	LOG_LIFE_TREE_MAIN_LEVEL        = 4301 //!创世神木升级
	LOG_LIFE_TREE_CAMP_LEVEL        = 4302 //!神木职业升级
	LOG_LIFE_TREE_CAMP_LEVEL_REBORN = 4303 //!神木职业重置

	LOG_ACTIVITY_3001  = 4401 //!显示活动-英雄传说
	LOG_TIGER_STEP     = 4402 //!纹章进阶
	LOG_TIGER_SKILL_UP = 4403 //!纹章技能升级

	LOG_ACTIVITY_3002      = 4501 //!显示活动-阵营传说
	LOG_ARMY_FLAG_UP       = 4502 //!佣兵军旗升级
	LOG_ARMY_EXCHANGE      = 4503 //!佣兵兑换
	LOG_ARMY_FLAG_EXCHANGE = 4504 //!佣兵军旗兑换

	LOG_ACTIVITY_3008 = 4601 //!月度活动-狂欢盛典

	LOG_MONTH_CARD_GOLD = 4701 //!月度活动荣耀金卡
	LOG_MONTH_CARD_BOX  = 4702 //!月卡箱子领取

	LOG_USER_CHARGE      = 4801 //!充值
	LOG_USER_CHARGE_GIFT = 4802 //!礼包充值

	LOG_LUCK_SHOP = 4901 //!幸运礼包

	LOG_BUY_GOLD = 5001 //!金币购买

	LOG_BUY_POWER = 5101 //!体力购买

	LOG_MAIL_GET = 5201 //!邮件领取

	LOG_BOSS_BUY = 5401 //!巨兽购买

	LOG_FIGHT_INSPIRE = 5501 //!鼓舞

	LOG_POWER_REPLACEMENT = 5701 //!补领体力

	//LOG_ONHOOK_AWARD = 5901 //!挂机奖励领取

	LOG_NOBILITY_AWARD = 6001 //!爵位奖励领取
	LOG_HYDRA_AWARD    = 6002 //!神器奖励领取

	LOG_ACTIVITY_FUND_AWARD = 6101 //!周基金奖励领取

	LOG_ASTROLOGY_ONE    = 6201 //!占星单抽
	LOG_ASTROLOGY_TEN    = 6202 //!占星十抽
	LOG_ASTROLOGY_CHANGE = 6203 //!更换占星目标英雄

	LOG_TARGETTASK_GET       = 6301 //!领取冒险任务奖励
	LOG_TARGETTASK_BADGE_GET = 6302 //!领取冒险任务积累奖励
	LOG_TARGETTASK_BADEG_BUY = 6303 //!激活冒险徽章

	LOG_TARGETTASK_TOWER_GET       = 6401 //!领取试炼之塔奖励
	LOG_TARGETTASK_BADGE_TOWER_GET = 6402 //!领取试炼之塔积累奖励
	LOG_TARGETTASK_BADEG_TOWER_BUY = 6403 //!激活试炼徽章

	LOG_INTERSTELLAR_NEBULA    = 6501 //!激活星云
	LOG_INTERSTELLAR_NEBULAWAR = 6502 //!激活星耀

	LOG_NOBILITY_UP     = 6601 //!爵位晋升
	LOG_NOBILITY_TASK   = 6602 //!领取爵位任务奖励
	LOG_NOBILITY_UP_EXT = 6603 //!领取爵位额外奖励

	LOG_RECHARGE_LIMIT_MAIN      = 6701 //!购买主线战令
	LOG_RECHARGE_LIMIT_MAIN_GET  = 6702 //!主线战令奖励
	LOG_RECHARGE_LIMIT_TOWER     = 6703 //!购买爬塔战令
	LOG_RECHARGE_LIMIT_TOWER_GET = 6704 //!爬塔战令奖励
	LOG_RECHARGE_LIMIT_DIAMOND     = 6705 //!购买钻石累消战令
	LOG_RECHARGE_LIMIT_DIAMOND_GET = 6706 //!钻石累消战令奖励

	LOG_FUND_1     = 6801 //!购买勇者基金
	LOG_FUND_2     = 6802 //!购买至尊基金
	LOG_FUND_1_GET = 6803 //!领取勇者基金奖励
	LOG_FUND_2_GET = 6804 //!领取至尊基金奖励
	LOG_FUND_3     = 6805 //!
	LOG_FUND_3_GET = 6806 //!

	LOG_ACCESSCARD_GET       = 6901 //!领取英雄收藏家奖励
	LOG_ACCESSCARD_SCORE_GET = 6902 //!领取英雄收藏家积分奖励

	LOG_GIFT_LOW = 7001 //!购买特惠礼包

	LOG_STAR_HERO_BUY       = 7101 //!购买星辰英雄
	LOG_STAR_HERO_BUY_LIMIT = 7102 //!购买星辰限时

	LOG_GENERAL_FIND_SINGLE = 7201 //!次元召唤单次召唤
	LOG_GENERAL_FIND_TEN    = 7202 //!次元召唤十连召唤
	LOG_GENERAL_FIND_BOX    = 7203 //!领取次元召唤积分宝箱
	LOG_GENERAL_FIND_RANK   = 7204 //!领取次元召唤排行奖励

	LOG_HERO_GROW_BOX      = 7301 //!领取免费英雄成长礼包
	LOG_HERO_GROW_RECHARGE = 7302 //!购买英雄成长礼包

	LOG_BEATUY_PASS_FINISH     = 7401 //!圣物关卡
	LOG_BEAUTY_TREASURE_ACTIVE = 7402 //!圣物部件激活
	LOG_BEAUTY_ACTIVE          = 7403 //!圣物激活
	LOG_BEAUTY_TREASURE_UP     = 7404 //!圣物部件升级
	LOG_BEAUTY_UP              = 7405 //!圣物升级
	LOG_BEAUTY_FIND_ONE        = 7406 //!圣物搜索单抽
	LOG_BEAUTY_FIND_FIVE       = 7407 //!圣物搜索五连抽

	LOG_ACT_CLIENT_SAY_OK = 7501 //!广告活动

	LOG_MONEY_TASK_FINISH      = 9901 //!赏金完成
	LOG_MONEY_TASK_DRAW        = 9902 //!赏金翻牌
	LOG_MONEY_TASK_REFRESH     = 9903 //!赏金刷新
	LOG_MONEY_TASK_DRAW_ONEKEY = 9904 //!赏金一键翻牌

	LOG_USER_DECOMPOUND      = 992001 //!分解
	LOG_HORSE_DECOMPOSE      = 992002 //!
	LOG_HORSE_SOUL_DECOMPOSE = 992003 //!
	LOG_HORSE_SOUL_EXCHANGE  = 992004 //!魔宠符文兑换
	LOG_HORSE_SOUL_UP        = 992005 //!魔宠符文强化
	LOG_HORSE_SWITCH         = 992006 //!魔宠转换
	LOG_HORSE_WASH           = 992007 //!魔宠洗练
	LOG_HORSE_BUY_NOMAL      = 992008 //!魔宠普通购买
	LOG_HORSE_BUY_SPECIAL    = 992009 //!魔宠高级购买
	LOG_HORSE_DISCERN        = 992010 //!魔宠召唤
	LOG_HORSE_COMBINE        = 992011 //!魔宠合成
	LOG_HORSE_UP             = 992012 //!魔宠升星
	LOG_HORSE_AWAKE          = 992013 //!魔宠觉醒

	LOG_USER_REBORN = 992101 //!重生

	LOG_KING_FIGHT_FINISH = 992201 //!王权争夺

	LOG_WAIT = 9999
)

// 平台相关
const (
	PLATFORM_DEFAULT = 0
	PLATFORM_WINDOWS = 1
	PLATFORM_ANDRIOD = 2
	PLATFORM_IOS     = 3
	PLATFORM_WP      = 4
	PLATFORM_MAC     = 5
)

// 排行榜
const (
	FightRank    = iota + 1 //1.战力排行榜
	TowerActRank            //2.爬塔排行榜
	PassStarRank            //3.关卡总星数排行榜
	TalentRank              //4.天赋总星级
	GemCostRank             //5.钻石消耗总星数排行榜

)

// 开启等级
const (
	LEVEL_OPEN_TITLE  = 6  // 称号开启等级
	LEVEL_OPEN_CASERN = 32 // 兵种开启等级
	LEVEL_OPEN_FATE   = 21 // 缘分开启等级
)

//! 阵营PVP
const (
	CAMP_PVP_FIGHT = 1 //! 国战
	CAMP_PVP_MINE  = 2 //! 矿点争夺
	CAMP_PVP_UNION = 3 //! 军团战
	CAMP_PVP_GVE   = 4 //! 资源争夺战-孤山夺宝
)

// 表里面的type要+1
const (
	TYPE_CASERN_GONG = 0 // 弓
	TYPE_CASERN_QI   = 1 // 骑
	TYPE_CASERN_BU   = 2 // 步
)

// 其他
const (
	DATEFORMAT    = "2006-01-02 15:04:05" // 时间格式化
	POWERMAX      = 1500                  // 最大体力限制
	SKILLPOINTMAX = 100                   // 最大技能点限制
	LEVELMAX      = 300                   // 最大等级限制80-90
	LEVELVIPMAX   = 18                    // vip最大等级
	ADDPOWERTIME  = 300                   // 体力回复时间(秒)
	ADDSPTIME     = 20                    // 技能点回复时间(秒)

	DEFAULT_JJC_WORSHIP_MONEY = 2000

	DEFAULT_GOLD         = 91000001 //! 金币
	DEFAULT_GEM          = 91000002 //! 元宝-钻石
	DEFAULT_EXP          = 91000005 //! 经验
	DEFAULT_INSPIRE      = 79000001 //! 鼓舞道具
	ACT_KEY_ITEM_ID      = 88800028 //! 转盘活动道具
	MONEY_MHERO          = 91000029 //! 无双币
	GVGNEED_ITEM_ID      = 91000020 //! GVG 消耗道具
	TigerChipId          = 91000031 // 铁矿
	GveTakeNum           = 91000036 // 据点占领值
	GveMaterial          = 91000037 // 战场物资
	GveCampMaterial      = 91000038 // 阵营战场物资
	GveGlory             = 91000039 // 孤山荣誉点
	HorseMatrial         = 91000027 //!
	HorseMatrial2        = 91000030 //! 战马材料2
	HeroSoul             = 91000032 // 魂石
	TechPoint            = 91000033 // 科技点
	BossMoney            = 91000034 // 巨兽精魄
	TowerStone           = 91000035 // 镇魂石
	HeroExp              = 91000041 // 卡牌经验
	LIVENESS_DAILY_POINT = 91000046 // 日常活跃度
	LIVENESS_WEEK_POINT  = 91000047 // 周活跃度
	WARORDER_ITEM_1      = 91000053 // 赏金社1
	WARORDER_ITEM_2      = 91000054 // 赏金社2
	LIFE_TREE_ITEM       = 91000055 // 成长涓流
	UNION_ACTIVITY_POINT = 91000057 // 公会活跃度

	DAY_SECS  = 86400
	HOUR_SECS = 3600
	MIN_SECS  = 60

	CAMP_SHU        = 1 // 帝国
	CAMP_WEI        = 2 // 联邦
	CAMP_WU         = 3 // 圣堂
	CAMP_QUN        = 4
	CITY_SHU        = 1104
	CITY_WEI        = 1059
	CITY_WU         = 1018
	FIGHTTYPE_DEF   = -200    // 城防军
	FIGHTTYPE_ROBOT = -500    // 机器人
	GM_ACCOUNT_ID   = 1000000 // 正常GM账号 = 10000，压测帐号放开限制

	SHOP_GENERAL     = 1  // 杂货-勋章
	SHOP_UNION       = 2  // 军团贡献
	SHOP_BOX         = 5  // 宝箱-斗技 神器
	SHOP_HONOR       = 6  // 荣誉 功勋 国库 商行
	SHOP_MAGICALSHOP = 7  // 无双商店
	SHOP_EXPEDITION  = 9  // 远征 秘境
	SHOP_TOWER       = 10 // 镇魂塔
	SHOP_DINIVITY    = 12 // 神格商店

	SHOP_NEW_NORMAL      = 1  //普通商店
	SHOP_NEW_UNION       = 2  //公会商店
	SHOP_NEW_FIRE        = 3  //遣散商店
	SHOP_NEW_PIT         = 4  //地牢商店
	SHOP_NEW_PVP         = 5  //高阶竞技场商店
	SHOP_NEW_PIT_SHOP    = 6  //地牢关卡
	SHOP_OLD_CONSUMERTOP = 10 //无双商店

	ACTIVITY_ORDER_HERO = 9006 // 预约神将活动
	INVEST_MONEY        = 9008 // 招财猫
	FIT_SERVER_ACT_TYPE = 9009 // 合服预热活动
	ACT_DIAL            = 9011 // 转盘
	ACT_TIMEGIFT        = 9021 //! 限时礼包
	//ACT_DAILY_RECHARGE1   = 9024 //! 每日首充-1
	//ACT_DAILY_RECHARGE2   = 9025 //! 每日首充-2
	ACT_DAILY_RECHARGE = 1015 //! 每日首充

	MAX_HORSE_AWAKEN_LEVEL = 7 // 最高等级

	ACTIVITY_STATUS_CLOSED = 0 // 活动关闭
	ACTIVITY_STATUS_OPEN   = 1 // 活动开启
	ACTIVITY_STATUS_SHOW   = 2 // 活动结束，可领奖，不可完成, 发送奖励

	MAX_UNION_DONATION = 6000 // 军团相关

	AWARD_NUM        = 7
	MAX_GENERAL_RANK = 50 // 最大神将排名

	MAX_CITY_POWER      = 150 // 军令上限
	MAX_CITY_POWER_TIME = 720 // 军令恢复时间

	MaxRankNum     = 100 // 排行榜最大数量
	MaxRankShowNum = 50  // 排行榜显示数量
	ActRankNum     = 20  // 排行榜活动显示数量

	GVGNEED_ITEM_NUM  = 5
	GVGNEED_ITEM_DESC = "GVG消耗"

	FightAtt = 99

	DEFAULT_HEAD_ICON = 1000
)

// Tariff类型
const (
	TariffPvpBuy               = 1 // pvp挑战次数购买
	TariffPvpRest              = 2 // pvp重置次数
	TariffTowerSet             = 3
	TariffTowerAdvanceBuy      = 4
	TariffTowerBuffBuy         = 5
	TariffGemStone             = 6
	TariffbuffKingTaskAction   = 7
	TariffbuffKingFlush        = 8
	TariffTaskDone             = 9
	TariffBuyBox               = 10
	TariffKing                 = 11
	TariffEditName             = 12
	TariffGvgPower             = 13
	TariffCreateUnion          = 15
	TariffChangeUnionName      = 16
	TariffGiveItems            = 20
	TariffSoldierUpgrade       = 21
	TariffSoldierWash          = 22
	TariffTakePower            = 23
	TariffLootGold             = 24
	TariffLootGem              = 25
	TariffArmy                 = 27
	TariffBuyArmy              = 28
	TariffBuyDungeonReset      = 29
	TariffReBorn               = 30
	TariffTalentReset          = 31
	TariffDreamlandGoldRefresh = 32
	TariffDreamlandGemRefresh  = 33
	TariffDreamlandGoldLoot    = 34
	TariffDreamlandGemLoot     = 35
)

// OpenLevel 类型
const (
	OPEN_LEVEL_TALENT_RESET = 82
	OPEN_LEVEL_REBORN_FREE  = 86
	OPEN_LEVEL_57           = 57
	OPEN_LEVEL_35           = 35
	OPEN_LEVEL_80           = 80
	OPEN_LEVEL_ON_HOOK      = 89
	//新魔龙
	OPEN_LEVEL_ARENA             = 4
	OPEN_LEVEL_SPECIAL_ARENA     = 5
	OPEN_LEVEL_HERO_REBORN       = 8
	OPEN_LEVEL_HERO_SUPPORT      = 9
	OPEN_LEVEL_RESONANCE_CRYSTAL = 15
	OPEN_LEVEL_NEWPITINFO        = 17
	OPEN_LEVEL_RANK_TASK         = 21
	OPEN_LEVEL_HIRE              = 40
	OPEN_LEVEL_WARORDER_2        = 46
	OPEN_ASTROLOGY               = 52
	OPEN_LEVEL_LAST_AWARD        = 57
	OPEN_LEVEL_LIFE_TREE         = 62
	OPEN_LEVEL_WARORDER_1        = 82
	OPEN_LEVEL_WARORDERLIMIT_1   = 90
	OPEN_LEVEL_WARORDERLIMIT_2   = 91
	OPEN_LEVEL_NEWPITSHOP        = 107
)

const (
	TYPE_RED_EMBATTLE = 1  // 上阵
	TYPE_RED_HERO     = 2  // 英雄
	TYPE_RED_EQUIP    = 3  // 装备
	TYPE_RED_ARMY     = 4  // 佣兵
	TYPE_RED_STAR     = 5  //星界、魔宠
	TYPE_RED_BOSS     = 6  // 巨兽
	TYPE_RED_BAG      = 7  // 背包
	TYPE_RED_WORLD    = 8  // 世界
	TYPE_RED_SHOP     = 9  // 市场
	TYPE_RED_TECH     = 10 // 学院
	TYPE_RED_BATTLE   = 11 // 战役
	TYPE_RED_BEAUTY   = 12 // 圣物
	TYPE_RED_GOD      = 13 // 神殿
	TYPE_RED_UNION    = 14 // 军团
	TYPE_RED_PVP      = 15 // 竞技场
	TYPE_RED_MAIL     = 16 // 邮件
	TYPE_RED_FRIEND   = 17 // 好友
)

//const (
//	RED_EMBATTLE = 1 << iota
//	RED_HERO
//	RED_EQUIP
//	RED_ARMY
//	RED_STAR
//	RED_BOSS
//	RED_BAG
//	RED_WORLD
//	RED_SHOP
//	RED_TECH
//	RED_BATTLE
//	RED_BEAUTY
//	RED_GOD
//	RED_UNION
//	RED_PVP
//	RED_MAIL
//	RED_FRIEND
//)
