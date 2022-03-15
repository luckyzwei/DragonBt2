package game

//! client2server
//! 版本验证
type C2S_CtrlHead struct {
	Ctrl  string `json:"ctrl"`
	Uid   int64  `json:"uid"`
	Os    string `json:"os"`
	Ver   int    `json:"ver"`
	MsgId int    `json:"msgid"`
}

//! 注册
type C2S_Reg struct {
	Account           string `json:"account"`
	Password          string `json:"password"`
	ServerId          int    `json:"serverid"`
	Platform_os       string `json:"os"`
	Platform_Brand    string `json:"brand"`
	Platform_DeviceId string `json:"deviceid"`
	Platform_Model    string `json:"model"`
}

//! sdk登陆
type C2S_SDKLogin struct {
	Password string `json:"password"`
	ServerId int    `json:"serverid"`
	Username string `json:"username"`
	Third    string `json:"third"`
}

type C2S_SDKIOSLogin struct {
	AppId     string                     `json:"appid"`
	MemId     string                     `json:"memid"`
	UserToken string                     `json:"token1"`
	ServerId  int                        `json:"serverid"`
	DevInfo   SendRZ_EnvInfo_DevInfo_Ios `json:"devInfo"`
	ChInfo    SendRZ_EnvInfo_ChInfo_Ios  `json:"chInfo"`
}

//! 重连
type C2S_AutoLogin struct {
	CheckCode string `json:"checkcode"`
}

//! 剧情id
type C2S_JqInfo struct {
	Jqid int `json:"jqid"`
}

type C2S_Encryption struct {
	EnInfo string `json:"eninfo"`
}

//! 指引id
type C2S_ZyInfo struct {
	Zyid int `json:"zyid"`
}

type C2S_CreateRole struct {
	Name string `json:"name"`
	Icon int    `json:"icon"`
	Face int    `json:"face"`
}

//! 取名
type C2S_AlertName struct {
	Newname string `json:"newname"`
}

//! 改icon
type C2S_AlertIcon struct {
	Icon int `json:"icon"`
}

//! 开始战斗
type C2S_BoxPass struct {
	Passid    int `json:"passid"`
	NoPass    int `json:"nopass"`
	MissionId int `json:"missionid"`
	Step      int `json:"step"`
	WarNum    int `json:"warnum"`
	WorkNum   int `json:"wroknum"`
}

//! 觐见
type C2S_JJPass struct {
	Passid int `json:"passid"`
}

//! 开始战斗
type C2S_BeginPass struct {
	Passid   int   `json:"passid"`
	Again    int   `json:"again"`
	AddFight int64 `json:"addfight"`
}

//! 跳过战斗
type C2S_PassSkip struct {
	Passid   int   `json:"passid"`
	UseArmy  int   `json:"usearmy"`
	AddFight int64 `json:"addfight"`
}

//! 结束战斗
type C2S_EndPass struct {
	Passid     int         `json:"passid"`
	Star       int         `json:"star"`
	Index      int         `json:"index"`
	MissionId  int         `json:"missionid"`
	Step       int         `json:"step"`
	WarNum     int         `json:"warnum"`
	WorkNum    int         `json:"wroknum"`
	FightTime  int         `json:"fighttime"`
	Again      int         `json:"again"`
	UseArmy    int         `json:"usearmy"`
	BattleInfo *BattleInfo `json:"battleinfo"`
}

//! 请求副本通关记录
type C2S_PassRecord struct {
	Passid int `json:"passid"`
}

//! 国战结束
type C2S_EndGz struct {
	Gzid int `json:"gzid"`
}

//! 扫荡
type C2S_SwapPass struct {
	Passid int `json:"passid"`
	Num    int `json:"num"`
}

//! 抽奖
type C2S_Find struct {
	Findtype int `json:"findtype"`
}

type C2S_FindPool struct {
	Findtype int `json:"findtype"`
	FindNum  int `json:"findnum"` //1   10
}

type C2S_FindAstrology struct {
	FindNum int `json:"findnum"` //1   10
}

type C2S_FindSelfSelection struct {
	FindNum int `json:"findnum"` //1   10
}

type C2S_GetVipRecharge struct {
	VipLevel int `json:"viplevel"` //对应哪个VIP的礼包  现在只有免费是需要手动领的
}

type C2S_FindSaveWish struct {
	Camp int   `json:"camp"` //1~4
	Ids  []int `json:"ids"`  //5个英雄ID
}

type C2S_FindOpenCamp struct {
	Findtype int `json:"findtype"`
}

type C2S_AstrologyHero struct {
	Id int `json:"id"`
}

type C2S_SelfSelectionHero struct {
	Id int `json:"id"`
}

//! 买玉符
type C2S_BuyYuFu struct {
	Num int `json:"num"`
}

//! 合成英雄
type C2S_HeroCompound struct {
	ItemId int `json:"heroid"` //! 物品Id
	Num    int `json:"num"`    //! 数量
}

// 自动分解开关
type C2S_HeroAutoFire struct {
	Action int `json:"action"` //! 0关闭 1开启
}

type C2S_HeroGetHandBook struct {
	HeroId int `json:"heroid"`
}

//! 使用道具
type C2S_UseItem struct {
	Itemid int `json:"itemid"`
	Num    int `json:"num"`
	Destid int `json:"destid"`
}

type C2S_UseItemSelect struct {
	Itemid    int   `json:"itemid"`
	Num       int   `json:"num"`
	Destid    []int `json:"destid"`
	DestidNum []int `json:"destidnum"`
}

//! 邮件领取
type C2S_GetMailItem struct {
	Mailid int64 `json:"mailid"`
}

//! 读邮件
type C2S_ReadMail struct {
	Mailid int64 `json:"mailid"`
}

//! 删邮件
type C2S_DelMail struct {
	Mailid int64 `json:"mailid"`
}

type C2S_InsertError struct {
	Stack     string `json:"stack"`
	ErrorInfo string `json:"errorinfo"`
	Param1    string `json:"param1"`
	Uid       int64  `json:"uid"`
}

//! 签到奖励
type C2S_CheckinAward struct {
	Index int `json:"index"`
}

//! 买金币
type C2S_BuyGold struct {
	Counts int `json:"counts"`
}

//! npc对话
type C2S_NpcAward struct {
	Id int `json:"id"`
}

//! 商店购买
type C2S_ShopBuy struct {
	Shoptype int `json:"shoptype"`
	Grid     int `json:"grid"`
	Refindex int `json:"refindex"`
}

type C2S_NewShopBuy struct {
	Shoptype int `json:"shoptype"`
	Grid     int `json:"grid"`
}

type C2S_NewShoprefresh struct {
	Shoptype int `json:"shoptype"`
}

//! 领取幸运商店
type C2S_LuckShopBuy struct {
	BoxId int `json:"boxid"`
}

type C2S_GetFund struct {
	Id int `json:"id"`
}



type C2S_OrderHeroDraw struct {
	Index int `json:"index"`
}

//! 宝箱购买
type C2S_EquBuy struct {
	Grid int `json:"grid"`
}

//! 商店刷新
type C2S_ShopRef struct {
	Type int `json:"type"`
}

//! 商店CD刷新
type C2S_ShopSysRef struct {
	Type int `json:"type"`
}

//! 加道具
type C2S_CreateItem struct {
	Itemid int `json:"itemid"`
	Num    int `json:"num"`
}

//! 合成道具
type C2S_MergeItem struct {
	Itemid int `json:"itemid"`
	Num    int `json:"num"`
}

//! 道具出售
type C2S_SellItem struct {
	Itemid string `json:"itemid"`
	Num    int    `json:"num"`
}

//! 道具批量出售
type C2S_SellItems struct {
	Item []PassItem `json:"item"`
}

//! 充值
type C2S_Recharge struct {
	Type  int `json:"type"`
	Money int `json:"money"`
}

//! 领取首充
type C2S_GetFirsetAward struct {
	Type int `json:"type"`
}

//! 领取VIP每日限购
type C2S_GetVipDailyReward struct {
	VipLevel int `json:"viplevel"`
}

//! 领取VIP每周购买福利
type C2S_BuyVipWeek struct {
	VipLevel int `json:"viplevel"`
}

//! 购买VIP特权礼包
type C2S_BuyVipBox struct {
	Index int `json:"index"`
}

//! 购买基金
type C2S_ButFund struct {
	FundType int `json:"fundtype"`
}

//! 领取基金奖励
type C2S_GetFundAward struct {
	Fundid int `json:"fundid"`
	Pageid int `json:"pageid"`
}

//! 领取福利奖励
type C2S_GetWeepPlanAward struct {
	Id int `json:"id"`
}

//! 完成任务
type C2S_TaskFinish struct {
	Taskid   string `json:"taskid"`
	Tasktype int    `json:"tasktype"`
}

type C2S_AccessCardTask struct {
	Taskid int `json:"taskid"`
}

type C2S_AccessCardAward struct {
	Id int `json:"id"`
}

//! 完成目标任务
type C2S_TargetTaskFinish struct {
	Taskid int `json:"taskid"`
}

type C2S_GetTargetLvReward struct {
	SystemId int `json:"systemid"`
}

type C2S_TaskFinish2 struct {
	Taskid int `json:"taskid"`
}

type C2S_CampTaskBox struct {
	CurBox int `json:"curbox"`
	Index  int `json:"index"`
}

//! 赏金社奖励领取
type C2S_GetWarOrderReward struct {
	Type int `json:"type"` //1皇家 2勇者
	Id   int `json:"id"`
}

//! 战令奖励领取
type C2S_GetWarOrderLimitReward struct {
	Type int `json:"type"` //1主线 2爬塔 3钻石累消
	Id   int `json:"id"`	// 领取id
}

type C2S_WarOrderBuy struct {
	Type   int `json:"type"` //1皇家 2勇者
	BuyNum int `json:"num"`  //购买数量
}


type C2S_CampTaskList struct {
	Index int `json:"index"`
}

//! 领取活跃度
type C2S_TaskLiveness struct {
	Id int `json:"id"`
}

//! 武将互换
type C2S_Interchange struct {
	Heroid [2]int `json:"heroid"`
}

type C2S_Gmaddexp struct {
	Cid    string `json:"cid"`
	Addexp int    `json:"addexp"`
}

type C2S_GMAutoLogin struct {
	Cid       string `json:"cid"`
	LoginType int    `json:"logintype"`
}

type C2S_SwitchServer struct {
	Cid      string `json:"cid"`
	ServerID int    `json:"serverid"`
}

//! 移动主公位置
type C2S_MoveCity struct {
	Cityid int `json:"cityid"`
}

//! 推图信息
type C2S_SetMission struct {
	MissionId int `json:"missionid"`
	Step      int `json:"step"`
	WarNum    int `json:"warnum"`
	WorkNum   int `json:"wroknum"`
}

//! 获取关卡的任务状态
type C2S_GetMission struct {
	MissionId int `json:"missionid"`
}

type C2S_Collection struct {
	Index int `json:"index"`
}

type C2S_ConsumerTopDraw struct {
	Id int `json:"id"` //! 奖励id
}

//! 天下大事
type C2S_WordEvent struct {
	Ver int `json:"ver"`
}

//! 天下大事领奖
type C2S_WordEventAward struct {
	Id int `json:"id"`
}

//! 部队改名
type C2S_CampTeamName struct {
	Index  int    `json:"index"`
	Name   string `json:"name"`
	Icon   int    `json:"icon"`
	UseSys int    `json:"usesys"`
}

//! 移动部队
type C2S_BigMap struct {
	Type int `json:"type"`
}

//! 完成事件
type C2S_FinishiEvent struct {
	Cityid int `json:"cityid"`
	Event  int `json:"event"`
}

//! 完成拜访
type C2S_FinishiVisit struct {
	Id    int `json:"id"`
	Index int `json:"index"`
	Stars int `json:"stars"`
}

// 完成偶遇
type C2S_FinishiVisitOuyu struct {
	Id int `json:"id"`
}

// 宝藏到期
type C2S_cangbaotimeover struct {
	Id int `json:"id"`
}

//! 开始特殊事件
type C2S_BeginTsFeel struct {
	Type int `json:"type"`
}

//! 完成特殊事件
type C2S_FinishTsFeel struct {
	Type       int `json:"type"`
	Spychildid int `json:"spychildid"`
}

//! 完成民生民情
type C2S_FinishiFeeling struct {
	Cityid int `json:"cityid"`
	Id     int `json:"id"`
	Step   int `json:"step"`
}

//! 修改当前阵型
type C2S_ChgCurTeam struct {
	CurTeam int `json:"curteam"`
}

// 确认可以偶遇
type C2S_msxfouyuok struct {
	Cityid int `json:"cityid"`
}

//! 购买地下城重置次数
type C2S_BuyiVisit struct {
	Id    int `json:"id"`
	Index int `json:"index"`
}

// 领取城池声望箱子
type C2S_GetSwBox struct {
	Cityid int `json:"cityid"`
	Index  int `json:"index"`
}

// 领取民事寻访星级宝箱
type C2S_GetMsxfBox struct {
	Msxfid int `json:"msxfid"`
	Index  int `json:"index"`
}

//! 领取城池声望箱子
type C2S_GetCityBox struct {
	Cityid int `json:"cityid"`
}

// 选择过关斩将类型
type C2S_choseGgzj struct {
	Index int `json:"index"`
	Max   int `json:"max"`
}

// 试练战斗结束
type C2S_ShilianOver struct {
	Win   int `json:"win"`
	Group int `json:"group"`
	Index int `json:"index"`
}

// 试练扫荡
type C2S_TrialSweep struct {
	Group int `json:"group"`
	Index int `json:"index"`
}

//! 得到城池数据
type C2S_GetCityInfo struct {
	Cityid int `json:"cityid"`
}

//!
type C2S_GetFightTeamInfo struct {
	Cityid int   `json:"cityid"`
	Pid    int64 `json:"pid"`
	Index  int   `json:"index"`
}

//! 移动部队
type C2S_MoveTeamBegin struct {
	Index  int `json:"index"`
	Cityid int `json:"cityid"`
	Begin  int `json:"begin"`
}

//! 直接传入部队
type C2S_MoveTeamFight struct {
	Index  int `json:"index"`
	Cityid int `json:"cityid"`
}

//! 移动部队
type C2S_MoveTeam struct {
	Index int   `json:"index"`
	Way   []int `json:"way"`
}

//! 宣战
type C2S_SayFight struct {
	Cityid int `json:"cityid"`
}

type C2S_CampFightWait struct {
	Cityid int `json:"cityid"`
	Page   int `json:"page"` //! 队列页码
	Pos    int `json:"pos"`  //! 0-attack，1-defence
}

//! 国战功勋排行榜
type C2S_CampFightTop struct {
	Cityid int `json:"cityid"` //! 城市Id
	Camp   int `json:"camp"`   //! 阵营
}

//! 进出阵营战
type C2S_CampFightMove struct {
	Cityid int `json:"cityid"`
	Type   int `json:"type"` //! 0进 1出
}

//! 国战召唤援军
type C2S_CampFightHelp struct {
	Cityid int `json:"cityid"`
}

//! 国战单挑-请求
type C2S_CampFightSoloReq struct {
	CityId    int `json:"cityid"`    //! 城市
	WarPlayId int `json:"warplayid"` //! 玩法Id
	Again     int `json:"again"`     //! 0初次 1继续单挑
	CityPart  int `json:"citypart"`  //! 目标选择 0-主城 1-5为据点
}

//! 国战单挑
type C2S_CampFightSolo struct {
	Cityid   int `json:"cityid"`
	Fast     int `json:"fast"`
	CityPart int `json:"citypart"` //! 目标选择 0-主城 1-5为据点
	Index    int `json:"index"`    //! 挑选的队伍顺序
}

//! 国战单挑
type C2S_CampFightSolo2 struct {
	Cityid   int   `json:"cityid"`  //! 城市Id
	MyIndex  int   `json:"myindex"` //! 自己所在Index
	Pid      int64 `json:"pid"`
	Index    int   `json:"index"`
	CityPart int   `json:"citypart"` //! 目标选择 0-主城 1-5为据点
}

//! 国战结果
type C2S_CampFightSoloEnd struct {
	Cityid int         `json:"cityid"`
	Soloid int         `json:"soloid"`
	Result int         `json:"result"`
	Info   []FightHero `json:"info"`
}

//! 国战使用技能
type C2S_CampFightUseSkill struct {
	Cityid   int `json:"cityid"`
	Skill    int `json:"skill"`
	Citypart int `json:"citypart"`
}

//! 国战撤退
type C2S_CampFightExit struct {
	Begin int `json:"begin"`
	Index int `json:"index"`
	End   int `json:"endcity"`
}

//! 国战突进
type C2S_CampFightEnter struct {
	Begin int   `json:"begin"`
	Index []int `json:"index"`
	End   int   `json:"endcity"`
}

//! 国战五虎争雄报名
type C2S_CampFight55Req struct {
	Cityid    int `json:"cityid"`
	WarPlayid int `json:"warplayid"`
}

//! 五虎争雄信息
type C2S_CampFight55Info struct {
	Cityid    int `json:"cityid"`
	WarPlayid int `json:"warplayid"`
}

//! 五虎争雄信息
type C2S_CampFight55Record struct {
	Cityid    int `json:"cityid"`
	WarPlayid int `json:"warplayid"`
	Page      int `json:"page"`
}

type C2S_CampFight56Final struct {
	Cityid int `json:"cityid"`
	Index  int `json:"index"`
	Page   int `json:"page"`
}

//! 退出观看
type C2S_CampFight55Exit struct {
	Cityid    int `json:"cityid"`
	WarPlayid int `json:"warplayid"`
}

//! 国战能突进
type C2S_CampFightCanEnter struct {
	Cityid    int `json:"cityid"`
	WarPlayid int `json:"warplayid"`
}

//! 国战同步回放
type C2S_CampFightSyncResult struct {
	FightId int64 `json:"fightid"`
}

//! 好友推荐
type C2S_FriendCommend struct {
	Refresh int `json:"refresh"`
}

//! 好友查找
type C2S_FriendFind struct {
	Pid       int64  `json:"pid"`
	FriendUid int64  `json:"frienduid"`
	Name      string `json:"name"`
}

//! 好友申请
type C2S_FriendApply struct {
	Pid int64 `json:"pid"`
}

//! 处理申请
type C2S_FriendOrder struct {
	Pid   int64 `json:"pid"`
	Agree int   `json:"agree"`
}

//! 加黑名单
type C2S_FriendBlack struct {
	Pid int64 `json:"pid"`
}

//! 删除好友
type C2S_FriendDel struct {
	Pid  int64 `json:"pid"`
	Type int   `json:"type"`
}

//! 删除好友
type C2S_FriendDelBatch struct {
	Pid  []int64 `json:"pid"`
	Type int     `json:"type"`
}

//! 好友体力
type C2S_FriendPower struct {
	Pid  int64 `json:"pid"`
	Type int   `json:"type"`
}

//! 查看详细信息
type C2S_Look struct {
	Pid int64 `json:"pid"`
}

//设置阵营
type C2S_SetCamp struct {
	Camp int `json:"camp"`
}

//设置指引
type C2S_SetGuild struct {
	GuildId int `json:"guildeid"`
}

//修改军团名字
type C2S_AlertUnionName struct {
	Unionid int    `json:"unionid"`
	Newname string `json:"newname"`
	IconId  int    `json:"iconid"`
}

//修改军团公告
type C2S_AlertUnionNotice struct {
	Unionid int    `json:"unionid"`
	Notice  string `json:"notice"`
}

//修改军团设置
type C2S_AlertUnionSet struct {
	Unionid   int `json:"unionid"`
	Jointype  int `json:"jointype"`
	Joinlevel int `json:"joinlevel"`
}

// 申请军团
type C2S_ApplyUnion struct {
	Unionid int `json:"unionid"`
}

// 取消申请军团
type C2S_Cancel_ApplyUnion struct {
	Unionid int `json:"unionid"`
}

// 创建军团
type C2S_CreateUnion struct {
	Name string `json:"name"`
	Icon int    `json:"icon"`
}

// 取消军团
type C2S_Dissolveunion struct {
	Unionid int `json:"unionid"`
}

// 军团捐献
type C2S_Donatemoney struct {
	Donatetype int `json:"donatetype"`
}

// 增加军团贡献度
type C2S_AddDonate struct {
	Donation int `json:"donation"`
}

//
type C2S_Getunioninfo struct {
	Unionid int `json:"unionid"`
}

// 加入军团
type C2S_Joinunion struct {
	Unionid int `json:"unionid"`
}

// 申请失败
type C2S_Masterfail struct {
	Unionid  int   `json:"unionid"`
	Applyuid int64 `json:"applyuid"`
}

// 申请成功
type C2S_Masterok struct {
	Unionid  int   `json:"unionid"`
	Applyuid int64 `json:"applyuid"`
}

// 申请成功
type C2S_MasterAllok struct {
	Unionid int `json:"unionid"`
}

//
type C2S_Masteroutplayer struct {
	Unionid int   `json:"unionid"`
	Outuid  int64 `json:"outuid"`
}

//
type C2S_Outunion struct {
	Unionid int `json:"unionid"`
}

type C2S_UnionModify struct {
	Unionid int   `json:"unionid"`
	Destuid int64 `json:"destuid"`
	Op      int   `json:"op"`
}

type C2S_SetBraveHand struct {
	Unionid int   `json:"unionid"`
	Destuid int64 `json:"destuid"`
	Op      int   `json:"op"`
}

type C2S_UnionFind struct {
	Type      int    `json:"type"`
	Unionname string `json:"unionname"`
	Unionid   string `json:"unionid"`
}

type C2S_MemberInfo struct {
	Unionid int   `json:"unionid"`
	Destuid int64 `json:"destuid"`
}

//! 军团贡献
type C2S_UnionDonation struct {
	Type int `json:"type"`
}

// 发送军团邮件
type C2S_UnionSendMail struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// 开始公会狩猎战斗
type C2S_StartHuntFight struct {
	Type int `json:"type"`
}

// 会长开启公会狩猎
type C2S_OpenHuntFight struct {
	Type int `json:"type"`
}

type C2S_EndHuntFight struct {
	Type       int         `json:"type"`
	Damage     int64       `json:"damage"`
	BattleInfo *BattleInfo `json:"battleinfo"`
}

type C2S_SweepHuntFight struct {
	Type int `json:"type"`
}

type C2S_GetHuntInfo struct {
}

//! dps排行
type C2S_HuntDpsTop struct {
	Cid  string `json:"cid"`
	Type int    `json:"type"`
}

//! 副本开始
type C2S_UnionCopyBegin struct {
	Id int `json:"id"` //! 副本开始
}

//! 副本结束
type C2S_UnionCopyEnd struct {
	Id      int                   `json:"id"`      //! 副本id
	Dps     int64                 `json:"dps"`     //! 伤害
	Monster []JS_UnionCopyMonster `json:"monster"` //! 怪物
}

//! 副本排行
type C2S_UnionCopyTop struct {
	Type int `json:"type"` //! 副本排行
}

type C2S_GMStr struct {
	Cid     string `json:"cid"`
	Gmstr   string `json:"gmstr"`
	Herolst []int  `json:"herolst"`
}

type C2S_GMTestMail struct {
	Cid     string `json:"cid"`
	ItemId  []int  `json:"itemid"`
	ItemNum []int  `json:"itemnum"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type C2S_GetTop struct {
	Cid   string `json:"cid"`
	Index int    `json:"index"`
	Ver   int    `json:"ver"`
}

type C2S_GetTopKing struct {
	Cid        string `json:"cid"`
	Camp       int    `json:"camp"`
	TopverKing int    `json:"topverking"`
}

//! 天下会武排行
type C2S_PvpTop struct {
	Cid string `json:"cid"`
	Ver int    `json:"ver"`
}

//! 设置防守阵容
type C2S_PvPDef struct {
	Format []int `json:"format"`
}

//! 天下会武战斗
type C2S_PvPFight struct {
	Index int `json:"index"`
}

//! 结束
type C2S_PvPEnd struct {
	Result int `json:"result"`
}

//! 得到官职战斗
type C2S_GetOfficeFight struct {
	Class int `json:"class"`
}

//! 设置防守阵容
type C2S_OfficeDef struct {
	Format [][]int `json:"format"`
}

//! 战斗
type C2S_OfficeFight struct {
	Index int `json:"index"`
}

//! 结束
type C2S_OfficeEnd struct {
	Result int `json:"result"`
}

//! 结束
type C2S_KingFightEnd struct {
	Result     int `json:"result"`
	Nobilityid int `json:"nobilityid"`
	Index      int `json:"index"`
	Overtime   int `json:"overtime"`
}

//! 领取
type C2S_OfficeAward struct {
	Class int `json:"class"`
}

// 赞国王
type C2S_CountryStepking struct {
	Index int `json:"index"`
}

//
type C2S_KingInfoReq struct {
	Cid     string `json:"cid"`
	Build   int    `json:"jzid"`
	Page    int    `json:"index"`
	KingVer int    `json:"kingver"`
}

type C2S_KingFightReq struct {
	Cid        string `json:"cid"`
	Target     int64  `json:"target"`
	Nobilityid int    `json:"nobility"`
}

//! 产业
type C2S_GetIndustry struct {
	Cityid int `json:"cityid"`
}

//! 产业
type C2S_Industry struct {
	Cityid int `json:"cityid"`
	Pos    int `json:"pos"`
}

type C2S_CampFightBox1 struct {
	Id int `json:"id"`
}

//!
type C2S_CampOneAward struct {
	Step int `json:"step"`
}

type C2S_FinishActivity struct {
	Id     int `json:"id"`
	Actver int `json:"actver"`
}

type C2S_GetPromoteBox struct {
	Id int `json:"id"`
}

type C2S_GetPromoteInfo struct {
}

type C2S_Barrage struct {
	Cityid int    `json:"cityid"`
	Size   int    `json:"size"`
	Text   string `json:"text"`
	Red    int    `json:"red"`
	Green  int    `json:"green"`
	Blue   int    `json:"blue"`
}

//! 激活宝箱碎片
type C2S_LevelUpBoxSP struct {
	Index int `json:"index"`
}

//! 激活藏宝图和挖掘
type C2S_TreasureProc struct {
	Id int `json:"id"`
}

//! 远征相关消息 ------------------------------------------------------------------------------
//! 获取当前节点消息
type C2S_GetExpeditionCurNodeInfo struct {
	CurIndex int `json:"curindex"` //! 当前节点1开始
}

//! 保存远征阵容
type C2S_SaveExpeditionTeamInfo struct {
	CurTeam  int      `json:"curteam"`  //! 上阵的当前阵容 步0 弓1 骑2
	Herolist [3][]int `json:"herolist"` //! 3个阵容的英雄ID数组 没有上阵给个空数组
}

//! 兑换物品
type C2S_ExchangeSoul struct {
	Index  int `json:"index"`
	Soulid int `json:"soulid"`
}

//! 刷新未执行的军演演练任务,发送下标,下标从1开始
type C2S_FlushMiliTask struct {
	TaskIndex int `json:"taskindex"`
}

//! 执行军演演练任务
type C2S_ExecuteMiliTask struct {
	TaskIndex int   `json:"taskindex"`
	HeroLst   []int `json:"herolst"`
}

//! 加速军演演练任务
type C2S_ClearMiliTaskCD struct {
	TaskIndex int `json:"taskindex"`
}

//! 领取军演演练任务奖励
type C2S_AwardMiliTask struct {
	TaskIndex int `json:"taskindex"`
}

//! 扫荡过关
type C2S_SweepGgzj struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
	Hard int    `json:"hard"`
}

//! 抽限时神将
type C2S_LootGeneral struct {
	LootType  int `json:"lootType"`  // 掉落类型:1.免费 2.单抽 3.10抽
	IsUseItem int `json:"isuseitem"` // 0不使用  1使用
}

//! 领取奖励
type C2S_GeneralAward struct {
	AwardIndex int `json:"awardIndex"` // 积分奖励索引1~5
}

//! 任命
type C2S_KingAppoint struct {
	Unionid int    `json:"unionid"` //! 军团Id
	Uid     int64  `json:"uid"`     //! 玩家Id
	Name    string `json:"name"`    //! 玩家姓名
	Icon    int    `json:"icon"`    //! 玩家头像
}

//! center服ctrl
type S2S_CenterCid struct {
	Cid string `json:"cid"`
}

//! 请求跨服数据
type S2Center_ConsumerTopGlobal struct {
	Cid      string `json:"cid"`
	ServerId int    `json:"serverid"`
}

//! 检测玩家武将信息
type C2S_CheckHero struct {
	FightType int               `json:"fightype"`
	Team      int               `json:"team"`
	Fight     int64             `json:"fight"`
	Rate      int               `json:"rate"`
	Param     [][]*JS_HeroParam `json:"param"`
}

//! 一键扫荡
type C2S_FinishiVisitAll struct {
	Id int `json:"id"`
}

//! 扫荡
type C2S_VisitSweep struct {
	Id    int `json:"id"`
	Index int `json:"index"`
	Flag  int `json:"flag"` //! 扫荡次数,非1表示全部扫荡
}

//! 扫荡
type C2S_DungeonSweep struct {
	Index int `json:"index"`
	Flag  int `json:"flag"` //! 扫荡次数,非1表示全部扫荡
}

//! 一键购买
type C2S_BuyiVisitAll struct {
	Id int `json:"id"`
}

//! 冶炼
type C2S_DoSmelt struct {
	SmeltType int `json:"smelttype"` //! 冶炼类型,1 免费 2.钻石
}

//! 穿戴虎符
type C2S_TigerAction struct {
	Action     int `json:"action"`     // action=1 穿装备 2 脱装备 3 精炼虎符 4 突破虎符 5 特技选择 6 特技升级 7 特技重置
	HeroId     int `json:"heroid"`     // 英雄Id
	KeyId      int `json:"keyid"`      // 虎符唯一Id
	StuntIndex int `json:"stuntindex"` // 操作的特技序号 1~3 4~6 从左往右, 从上到下
	StuntId    int `json:"stuntids"`   // 特技选择id, 对应配置表的特技id, 实际就是一个技能Id
	TeamType   int `json:"team_type"`  // 阵营类型
	Index      int `json:"index"`      // 上阵序号
}

//! 获取活动信息
type C2S_GetActTop struct {
	ActType int `json:"acttype"`
}

//GVG 驻扎
type C2S_CityGvGDefense struct {
	CityId int `json:"cityid"` //! 城池信息
	Team   int `json:"team"`   //! 出战队伍
	Flag   int `json:"flag"`   //! 1表示无条件驻扎
}

//GVG 驻扎离开
type C2S_CityGvGDefenseLeave struct {
	CityId int `json:"cityid"` //! 城池信息
}

//GVG 压制
type C2S_CityGvGAttack struct {
	CityId int `json:"cityid"` //! 城池信息
	Team   int `json:"team"`   //! 出战队伍
	Flag   int `json:"flag"`   //! 1表示无条件压制
}

//GVG 第三方清剿
type C2S_CityGvGOtherAttack struct {
	CityId int `json:"cityid"` //! 城池信息
	Team   int `json:"team"`   //! 出战队伍
}

//GVG 反击
type C2S_CityGvGAttackBack struct {
	CityId int `json:"cityid"` //! 城池信息
	Team   int `json:"team"`   //! 出战队伍
}

//GVG 战斗结果
type C2S_CityGvGFightResult struct {
	CityId int `json:"cityid"` //! 城池信息
	Result int `json:"result"` //! 战斗结果
}

//GVG 城池信息
type C2S_CityGvGCityInfo struct {
	CityId int `json:"cityid"` //! 城池信息
}

//GVG 城池队伍
type C2S_CityGvGCityQueue struct {
	CityId int `json:"cityid"` //! 城池信息
}

//GVG 城池方针设定
type C2S_CityGvGSetPolicy struct {
	CityId int `json:"cityid"` //! 城池信息
	Policy int `json:"policy"` //! 方针id
}

//GVG 城池领奖
type C2S_CityGvGGetPosAward struct {
	CityId int `json:"cityid"` //! 城池信息
}

//GVG 城池领取特殊奖励
type C2S_CityGvGGetSpecialAward struct {
	CityId int `json:"cityid"` //! 城池信息
}

//GVG 城池特殊奖励数据
type C2S_CityGvGGetSpecialAwardInfo struct {
	Idx    int `json:"idx""`
	CityId int `json:"cityid"`
}

//! 发送红包池的红包
type C2S_SendPool struct {
	KeyId int64  `json:"keyid"`
	Msg   string `json:"msg"`
}

//! 抢红包
type C2S_GotRedPac struct {
	KeyId int64 `json:"keyid"` //! 红包Id
}

//! 查看红包
type C2S_LookRedPac struct {
	KeyId int64 `json:"keyid"` //! 红包Id
}

//! 答谢红包
type C2S_ThankRedPac struct {
	KeyId int64 `json:"keyid"` //! 红包Id
}

//! 发送军团/国家红包
type C2S_SendGlobalRed struct {
	RedId       int    `json:"redid"`     //! 选择的配置Id
	SelectIndex int    `json:"selectnum"` //! 选择的档位 1,2,3
	Msg         string `json:"msg"`       //! 红包留言
	RedType     int    `json:"redtype"`   //! 1国家 2军团
}

//! 获取红包信息
type C2S_GetRedPacs struct {
	RedType int `json:"redtype"` //! 1国家页签 2军团页签
}

//! 抽奖
type C2S_DoDial struct {
	LootType int  `json:"lootType"` //! 掉落类型:1.免费 2.单抽 3.10抽
	UseKey   bool `json:"usekey"`   //! 是否使用钥匙[2,3时发送, 不使用发false]
}

//! 转盘领取箱子奖励
type C2S_DialAward struct {
	AwardIndex int `json:"awardIndex"` // 积分奖励索引1~5
}

//! 翻牌领取箱子奖励
type C2S_DrawAward struct {
	AwardIndex int `json:"awardIndex"` // 积分奖励索引1~5
}

//! 翻牌币购买
type C2S_DrawBuy struct {
	Cid string `json:"cid"` //! cid
	Num int    `json:"num"` //! 购买数量
}

//! 翻牌币购买
type C2S_DrawAction struct {
	Cid      string `json:"cid"`      //! cid
	DrawType int    `json:"drawtype"` //! 翻牌类型 1.开始翻牌 2.重置翻牌 3.进行翻牌
	Index    int    `json:"index"`    //! index=1~8
}

//! 后宫宝物升级
type C2S_Beauty_TreasureaUpLevel struct {
	Beautyid       int `json:"beautyid"`
	TreasureaIndex int `json:"index"`
}

//! 后宫升级
type C2S_Beauty_UpLevel struct {
	Beautyid int `json:"beautyid"`
}

//! 后宫传奇战斗结束
type C2S_Beauty_LegendOver struct {
	Beautyid int `json:"beautyid"`
	Chapter  int `json:"chapter"`
	Index    int `json:"index"`
}

//! 后宫
type C2S_BeautyOpen struct {
	BeautyId int `json:"beautyid"`
}

//! 后宫升级
type C2S_BeautyUpLevel struct {
	BeautyId int `json:"beautyid"`
	Type     int `json:"type"`
}

//! 后宫升级
type C2S_BeautyUpLevel2 struct {
	BeautyId int `json:"beautyid"`
	Type     int `json:"type"`
}

//! 后宫升星
type C2S_BeautyUpStar struct {
	BeautyId int `json:"beautyid"`
}

//! 后宫
type C2S_BeautyGetEvent struct {
	BeautyId int `json:"beautyid"`
}

//! 后宫
type C2S_BeautyFinishEvent struct {
	BeautyId int `json:"beautyid"`
}

//! 美人排行
type C2S_BeautyTop struct {
	Cid string `json:"cid"`
	Ver int    `json:"ver"`
}

//! 宝物相关操作
type C2S_TreasureAction struct {
	Action int `json:"action"` //! action=1.宝物合成 2.宝物分解 3.宝物穿戴 4一键穿戴 5.脱一件装备 6.宝物一键卸载
	// 7.宝物洗练:免费 8.宝物洗练:初级 9.宝物洗练:中级 10.宝物洗练:高级 11.宝物觉醒 12.宝物重置 13.一键洗练
	Itemid          int   `json:"itemid"`          //! 合成的碎片Id[合成时发送]
	CompundNum      int   `json:"compoundnum"`     //! 合成的宝物数量[合成时发送]
	KeyId           int   `json:"keyid"`           //! 宝物keyid[分解时发送,穿1件装备时发送]
	HeroId          int   `json:"heroid"`          //! 装备时武将Id
	Pos             int   `json:"pos"`             //! 第几个位置宝物[脱宝物时发送]
	WashIndex       int   `json:"washindex"`       //! 一键洗练, washindex=1,2,3,4
	DecomposeKeyIds []int `json:"decomposekeyids"` //! 多个宝物
}

type C2S_DailyDiscount struct {
	Id int `json:"id"` //! 活动Id
}

type C2S_GetOverflow struct {
	Id int `json:"id"` //! 活动Id
}

type C2S_GetCumulativeReward struct {
	Id int `json:"id"` //! 奖励ID
}

type C2S_LuckStartAward struct {
	TaskId int `json:"taskid"` //! 领取任务Id
}

//! 激活所有藏宝图和挖掘
type C2S_FinishAllTreasure struct {
	Id int `json:"id"`
}

//! 领取基金
type C2S_ActivityFundAward struct {
	Cid string `json:"cid"`
	Ver int    `json:"fundver"`
	Pay int    `json:"pay"`
	Day int    `json:"day"`
}

//! 获得基金信息
type C2S_ActivityFundInfo struct {
	Ver int `json:"fundver"`
}

type C2S_DailyRechargeAward struct {
	TaskId int `json:"taskid"` //! 领取任务Id
}

//! 翻牌币购买
type C2S_DoLuckEgg struct {
	Index int `json:"index"` //! index=1~3
}

//! 上阵英雄
type C2S_AddTeamUIPos struct {
	TeamType int   `json:"teamtype"` //! 1~6
	FightPos []int `json:"fightpos"` //上阵英雄KeyId 必须5个数
}

//! 交换战斗位置
type C2S_SwapFightPos struct {
	Index1   int `json:"index1"`   //! 战斗位置1
	Index2   int `json:"index2"`   //! 战斗位置2
	TeamType int `json:"teamtype"` //! 布阵类型, 1副本 2竞技场
}

//! 切换阵型
type C2S_ChangeFormation struct {
	Id       int `json:"id"`       //! 阵型ID
	TeamType int `json:"teamtype"` //! 布阵类型, 1副本 2竞技场
}

//! 激活资质
type C2S_ActivateStar struct {
	HeroId int `json:"heroid"` //! 英雄Id
	Index  int `json:"index"`  //! 资质孔的序号
}

//! 升星
type C2S_HeroUpStar struct {
	HeroKeyId     int   `json:"herokeyid"`     //! 英雄Id
	CostHeroKeyId []int `json:"costherokeyid"` //消耗选择的物品
}

type C2S_HeroUpStarAll struct {
	HeroKeyId     []int   `json:"herokeyid"`     //! 英雄Id
	CostHeroKeyId [][]int `json:"costherokeyid"` //消耗选择的物品
}

type C2S_HeroLock struct {
	HeroKeyId  int `json:"herokeyid"`  //! 英雄keyId
	LockAction int `json:"lockaction"` //! 0解锁  1锁定
}

//! 保存阵容
type C2S_SaveTeam struct {
	TeamType int   `json:"teamtype"` //! 1~6
	FightPos []int `json:"fightpos"` //上阵英雄KeyId
}

type C2S_BaseHeroSet struct {
	FightPos []int `json:"fightpos"` //上阵英雄KeyId
}

type C2S_HireApply struct {
	HireAction    int   `json:"hireaction"` //0申请  1取消申请  2放弃租用
	HireOwnUid    int64 `json:"hireownuid"`
	HireHeroKeyId int   `json:"hireherokeyid"`
}

type C2S_HireStateSet struct {
	HireState     int   `json:"hirestate"`     //0拒绝 1接受
	HireHeroKeyId int   `json:"hireherokeyid"` //英雄
	HireUid       int64 `json:"hireuid"`       //针对谁的申请操作
}

type C2S_GetPlayerTeam struct {
	FriendUid int64 `json:"frienduid"`
	TeamType  int   `json:"teamtype"`
}

//! 一键升星
type C2S_UpStarAuto struct {
	HeroId int `json:"heroid"` //! 英雄Id
}

//! 激活天赋
type C2S_UpgradeTalent struct {
	HeroId int `json:"heroid"` //! 英雄Id
	Index  int `json:"index"`  //! 天赋的序号
}

//! 重置天赋
type C2S_ResetTalent struct {
	HeroId int `json:"heroid"` //! 英雄Id
}

//! 英雄升级
type C2S_HeroLvUp struct {
	HeroKeyId int `json:"herokeyid"` //! 英雄keyId
}

type C2S_HeroLvUpTo struct {
	HeroKeyId int `json:"herokeyid"` //! 英雄keyId
	AimLv     int `json:"aimlv"`     //! 目标等级
}

//! 合成
type C2S_Compose struct {
	Id  int `json:"id"`  //! 配方ID
	Num int `json:"num"` // 次数
}

//! 获得幻境信息
type C2S_DreamLandBaseInfo struct {
}

//! 获得幻境物品信息
type C2S_DreamLandItemInfo struct {
	Type int `json:"type"` //! 类型
}

//! 幻境抽奖
type C2S_DreamLandLoot struct {
	Type  int `json:"type"`  //! 类型
	Times int `json:"times"` //! 次数
}

//! 幻境刷新
type C2S_DreamLandRefresh struct {
	Type int `json:"type"` //! 类型
}

//! 激活缘分
type C2S_ActivateFate struct {
	HeroId int `json:"heroid"` //! 英雄Id
	Index  int `json:"index"`  //! 缘分index,第几个缘分
}

//! 英雄重生
type C2S_HeroReborn struct {
	HeroKeyId int `json:"herokeyid"` //! 英雄keyId
}

type C2S_HeroBack struct {
	HeroKeyId int `json:"herokeyid"` //! 英雄keyId
}

type C2S_HeroFire struct {
	HeroKeyId []int `json:"herokeyid"` //! 英雄keyId
}

//! 虚空英雄共鸣
type C2S_VoidHeroResonanceSet struct {
	HeroKeyId     int `json:"herokeyid"`     //! 英雄keyId
	VoidHeroKeyId int `json:"voidherokeyid"` //! 英雄keyId
}

//! 取消虚空英雄共鸣
type C2S_VoidHeroResonanceCancel struct {
	VoidHeroKeyId int `json:"herokeyid"` //! 英雄keyId
}

// 设置天赋技能
type C2S_SetStageTalentSkill struct {
	HeroKeyId int `json:"herokeyid"` //! 英雄keyId
	Index     int `json:"index"`     //! 第几层
	Pos       int `json:"pos"`       //! 技能位置
}

//! 装备相关操作
type C2S_EquipAction struct {
	Action     int `json:"action"`     //! action=1.穿 2.脱
	HeroKeyId  int `json:"herokeyid"`  //! 英雄keyId
	EquipKeyId int `json:"equipkeyid"` //! 装备KeyId
}

type C2S_EquipActionAll struct {
	Action    int `json:"action"`    //! action=1.一键穿 2.一键脱
	HeroKeyId int `json:"herokeyid"` //! 英雄keyId
}

type C2S_EquipUpLv struct {
	EquipKeyId int   `json:"equipkeyid"` //! 装备KEYID
	ItemKeyId  []int `json:"itemkeyid"`  //!    5   6    0
	ItemId     []int `json:"itemid"`     //!    1   1    11111111
	ItemNum    []int `json:"itemnum"`    //!    3   1     5
}

//! 神器相关操作
type C2S_ArtifactEquipAction struct {
	Action             int `json:"action"`             //! action=1.穿 2.脱
	HeroKeyId          int `json:"herokeyid"`          //! 英雄keyId
	ArtifactEquipKeyId int `json:"artifactequipkeyid"` //! 装备KeyId
}

type C2S_ArtifactEquipUpLv struct {
	ArtifactEquipKeyId int `json:"artifactequipkeyid"` //!
}

type C2S_ExclusiveAction struct {
	Action    int `json:"action"`    //! action=1.解锁 2.升级
	HeroKeyId int `json:"herokeyid"` //! 英雄keyId
}

//符文相关操作
type C2S_RuneAction struct {
	Action int `json:"action"` //!
	HeroId int `json:"heroid"` //! 英雄Id
	KeyId  int `json:"keyid"`  //! 符文Id
	Pos    int `json:"pos"`    //! 第几个位置
}

//! 宝石相关操作
type C2S_GemAction struct {
	Action    int   `json:"action"`    //! action=1.宝石合成 2.宝石镶嵌 3.宝石卸载 4.宝石一键卸下
	GemId     int   `json:"gemid"`     //! 合成的宝石Id, 穿戴的宝石Id
	KeyId     int   `json:"keyid"`     //! 装备唯一Id
	Pos       int   `json:"pos"`       //! 装备孔位置 1~6
	AutoGemId []int `json:"autogemid"` // 一键镶嵌时宝石id
	AutoPos   []int `json:"autopos"`   // 一键镶嵌时位置
}

//! 分解英雄碎片
type C2S_HeroDecompose struct {
	ItemIds  []int `json:"itemid"`  //! 道具Id
	ItemNums []int `json:"itemnum"` //! 道具数量
}

//! 地下城组队关卡相关操作
//关卡信息
type C2S_DungeonServerInfo struct {
}

//关卡信息
type C2S_DungeonUserInfo struct {
}

//创建队伍
type C2S_DungeonCreateTeam struct {
	InstanceId int   `json:"levelid"`   //! 关卡id
	AutoEnter  int   `json:"autoenter"` //! 满员是否自动进入战斗 1表示自动加入
	JoinFlag   []int `json:"joinflag"`  //!禁入国家限制1表示禁入
}

//队伍信息
type C2S_DungeonTeamInfo struct {
	TeamId int64 `json:"teamid"` //! teamid
}

//加入队伍
type C2S_DungeonJoinTeam struct {
	TeamId int64 `json:"teamid"` //! 队伍id
}

//离开队伍
type C2S_DungeonLeaveTeam struct {
}

//踢出成员
type C2S_DungeonKickPlayer struct {
	PlayerId int64 `json:"playerid"` //! 队伍id
}

//呼唤
type C2S_DungeonCall struct {
}

//队长开始战斗
type C2S_DungeonFightBegin struct {
}

//队伍列表
type C2S_DungeonTeamList struct {
	InstanceId int `json:"levelid"` //! 关卡id
}

//! 科技相关操作
type C2S_TechAction struct {
	Action int `json:"action"` //! action=1.开始研究科技 2.立即研究 3.科技完成请求[无参数] 4.取消科技研究队列 5.科技一键完成[无参数] 6.单科技一键完成
	TechId int `json:"techid"` //! 科技Id[开始研究,立即研究]
	KeyId  int `json:"keyid"`  //! 科技队列keyId[取消研究科技,单个科技一键完成]
}

//! 镇魂塔相关操作
//数据获取
type C2S_TowerInfo struct {
}

//排行榜
type C2S_TowerRank struct {
	Cid  string `json:"cid"`
	Type int    `json:"type"`
}

//开始
type C2S_TowerFightBegin struct {
	LevelId int `json:"levelid"` //! 关卡id
}

type C2S_TowerFightSkip struct {
	LevelId int `json:"levelid"` //! 关卡id
}

//通关
type C2S_TowerFightResult struct {
	LevelId    int         `json:"levelid"` //! 关卡id
	Result     int         `json:"result"`  //! 战斗结果 1胜利 2失败
	BattleInfo *BattleInfo `json:"battleinfo"`
}

//! 楼层信息
type C2S_TowerFloorInfo struct {
	Cid     string `json:"cid"`
	LevelId int    `json:"levelid"` //! 关卡id
}

//巨兽购买
type C2S_BossAction struct {
	Action int `json:"action"` //! action=1.购买巨兽
	Id     int `json:"id"`     //! 巨兽Id
}

//! 宝石副本相关操作
//数据获取
type C2S_GemStoneInfo struct {
}

//通关
type C2S_GemStoneFightResult struct {
	LevelId int `json:"levelid"` //! 关卡id
	Result  int `json:"result"`  //! 战斗结果 1胜利 2失败
}

type C2S_GemStoneSweep struct {
	Id int `json:"id"` //! 章节id
}

type C2S_GemStoneBuySweepTimes struct {
	Id int `json:"id"` //! 章节id
}

//! 领取高级召唤
type C2S_AwardHorseTask struct {
	Id int `json:"id"` //! 任务Id
}

//! 冶炼次数任务奖励
type C2S_AwardSmeltTask struct {
	Id int `json:"id"` //! 任务Id
}

//! 购买冶炼次数任务奖励
type C2S_AwardBuySmeltTask struct {
	Id int `json:"id"` //! 任务Id
}

//! 战马转换
type C2S_SwitchHorse struct {
	KeyId int `json:"keyid"` //! 发送唯一Id,不要发送horseId
}

//! 战马洗练
type C2S_SwitchWash struct {
	KeyId int `json:"keyid"` //! 发送唯一Id,不要发送horseId
}

//! 战马保存
type C2S_SaveWash struct {
	KeyId int `json:"keyid"` //! 发送唯一Id,不要发送horseId
}

//! 战马保存
type C2S_SaveSwitch struct {
	KeyId int `json:"keyid"` //! 发送唯一Id,不要发送horseId
}

//! 领取任务奖励
type C2S_AwardMiliWeekTask struct {
	TaskId int `json:"taskid"` //! 任务Id
}

// 佣兵操作
type C2S_SoldierAction struct {
	Action    int   `json:"action"`    // action=1 穿佣兵
	HeroId    int   `json:"heroid"`    // 英雄Id
	KeyId     int   `json:"keyid"`     // 佣兵唯一Id
	RemoveIds []int `json:"removeIds"` // 佣兵唯一Id(多个)
	Pos       int   `json:"pos"`       // 位置
	WashNum   int   `json:"washnum"`   // 洗几次
	LockIndex int   `json:"lockindex"` // 第几个属性(上锁或者解锁)
	SaveIndex int   `json:"saveindex"` // 继承第几个属性
	Poses     []int `json:"poses"`     // 一键穿的位置
}

// 刷新任务
type C2S_FlushKingTask struct {
	HardMode int `json:"hard_mode"` // 难度
}

//! 斗技场排行
type C2S_PvpCountryTop struct {
	Cid  string `json:"cid"`
	Camp int    `json:"camp"`
}

// 购买任务宝箱
type C2S_BuyKingBox struct {
	BoxIndex int `json:"index"` // 宝箱下标
}

// 矿点移动
type C2S_MineAction struct {
	Action int `json:"action"`  // 1.进入矿点 2.移动矿点
	MineId int `json:"mine_id"` // 从当前矿点移动到这个矿点
}

// 查看矿点战报回放
type C2S_GetMineFight struct {
	FightId int64 `json:"fight_id"` // 战斗Id
}

// 孤山多宝操作
type C2S_GveAction struct {
	Action  int `json:"action"`   // 1.进入建筑 2.移动 3.开始战斗 4.战斗结束
	BuildId int `json:"build_id"` // 从当前矿点移动到这个矿点
	Result  int `json:"result"`   // gve战斗结果
	LevelId int `json:"level_id"` // 关卡Id
}

type C2S_PassMission struct {
	ChpaterId int `json:"chpater_id"`
}

type C2S_GMTechUp struct {
	TechId int `json:"tech_id"`
	TechLv int `json:"tech_lv"`
}

type C2S_AddUnionExpLimit struct {
	AddNum int `json:"add_num"`
}

type C2S_SetTowerLevel struct {
	Type  int `json:"type"`
	Level int `json:"level"`
}

type C2S_GMSuperHelp struct {
	Index int `json:"index"`
}

type C2S_GMNobilityUp struct {
	Level int `json:"level"`
}

type C2S_GMInterstellar struct {
	NebulaId int `json:"nebulaid"` //星云id
	GroupId  int `json:"groupid"`  //地图id
	TaskId   int `json:"taskid"`   //任务id
}

type C2S_GMInstancePass struct {
	Id int `json:"id"` //对应ID，0表示全部通关
}

// 新佣兵操作
type C2S_ArmyAction struct {
	Action     int `json:"action"`      // action=1 部署佣兵
	TeamType   int `json:"team_type"`   // 阵营类型 默认填1
	Index      int `json:"index"`       // 上阵下标1
	Index2     int `json:"index2"`      // 上阵下标2
	FlagIndex  int `json:"flag_index"`  // 军旗下标
	ArmyId     int `json:"army_id"`     // 佣兵Id
	FlagId     int `json:"flag_id"`     // 军旗Id
	ExchangeId int `json:"exchange_id"` // 兑换Id
	BuyType    int `json:"buy_type"`    // 挑战类型 1.免费 2.钻石
	Id         int `json:"id"`          // 配置Id
	HeroId     int `json:"heroid"`      // 英雄Id
}

// 创建队伍
type C2S_ArmyCreateTeam struct {
	InstanceId int   `json:"levelid"`   //! 关卡id
	AutoEnter  int   `json:"autoenter"` //! 满员是否自动进入战斗 1表示自动加入
	JoinFlag   []int `json:"joinflag"`  //!禁入国家限制1表示禁入
}

// 队伍列表
type C2S_ArmyTeamList struct {
	InstanceId int `json:"levelid"` //! 关卡id
}

// 加入队伍
type C2S_ArmyJoinTeam struct {
	TeamId int64 `json:"teamid"` //! 队伍id
}

// 踢出成员
type C2S_ArmyKickPlayer struct {
	PlayerId int64 `json:"playerid"` //! 玩家Id
}

// 拉取单个州的报名情况
type C2SGetUnionInfo struct {
	StateId int `json:"state_id"`
}

// 拉取宣战列表情况
type C2SStateAttendInfo struct {
	StateId int `json:"state_id"`
}

// 宣战
type C2SUnionOpenWar struct {
	State int `json:"state"` // 州类型
}

// 报名
type C2SUnionAttendWar struct {
	State int `json:"state"` // 州类型
}

// 获取州宝箱奖励
type C2SAwardState struct {
	StateId int `json:"stateId"`
}

// 拉取宣战列表情况
type C2SGetUnionFight struct {
	FightId int64 `json:"fight_id"`
}

// 查询报名信息
type C2SQueryAttend struct {
	State int `json:"state"` // 州类型
}

// 查询决赛信息
type C2SQueryFinal struct {
	State int `json:"state"` // 州类型
	Round int `json:"round"` // 回合
}

type C2SUnignGroupBattle struct {
	State int `json:"state"` // 州类型
	Round int `json:"round"` // 第几轮
}

type C2S_UnionGetFights struct {
	StateId int `json:"stateid"` // 州Id
	Round   int `json:"round"`   // 第几轮
	Group   int `json:"group"`   // 组
}

type C2S_GetTeamRecord struct {
	StateId int `json:"stateid"` // 州Id
}

//设置红色标记
type C2S_SetRedIcon struct {
	Id int `json:"id"` // 功能Id
}

type C2S_Language struct {
	Id int `json:"id"` // 功能Id
}

type C2S_Nationality struct {
	Id int `json:"id"` // 功能Id
}

type C2S_CheckFight struct {
	CheckId  int64             `json:"checkid"`  // 功能Id
	TeamType int               `json:"teamtype"` // 阵容Id
	Info     map[int][]float64 `json:"info"`     //key:英雄key
}

type C2S_AddGuide struct {
	GuideId int `json:"guide_id"`
}

type C2S_AddStory struct {
	StoryID   int `json:"story_id"`
	StoryType int `json:"story_type"`
}

// 快速进入
type C2S_FastEnter struct {
	InstanceId int `json:"instance_id"` //! 关卡id
}

// 设置自动进入
type C2S_SetAutoEnter struct {
	TeamId int64 `json:"teamid"` //! 队伍id
}

type C2SAwardStatistics struct {
	Id int `json:"id"` //!  积分
}

type C2STakeNobilityAward struct {
	TaskId int `json:"taskid"` //!  任务ID
}

type C2SLevelUpNobility struct {
	TaskId int `json:"taskid"` //!  目标爵位等级 对应表里ID
}

type C2S_GetNobilityReward struct {
	TaskId int `json:"taskid"` //!  目标爵位等级 对应表里ID
	Belog  int `json:"belog"`
}

//! 全民商店购买
type C2S_WholeShopBuy struct {
	Id int `json:"id"`
}

//! 请求全民商店信息
type C2S_WholeShopInfo struct {
	Edition int `json:"edition"`
}

type C2S_OnHookStage struct {
	Stage int `json:"stage"` //切换章节的时候需要手动通过关卡
}

type C2S_OnHookTeams struct {
	Id []int `json:"id"` // 英雄ID
}

type C2STakeHydraAward struct {
	TaskId int `json:"taskid"` //!  任务ID
}

//领取神兽，发送任务等级
type C2STakeHydra struct {
	HydraID int `json:"hydraid"` //!  任务ID
}

//神兽
type C2S_HydraAction struct {
	Action int `json:"action"` //! action=1.神兽升级
	Id     int `json:"id"`     //!
	Index  int `json:"index"`  //!
}

//! 上阵神兽
type C2S_AddHydra struct {
	HydraId  int `json:"hydraid"`  //! 神兽Id
	TeamType int `json:"teamtype"` //! 布阵类型, 1副本 2竞技场
}

type C2S_PitFinishEvent struct {
	PitKeyId int `json:"pitkeyid"`
	EventId  int `json:"eventid"` //!  关卡ID
	Option   int `json:"option"`  //!  选项
}

type C2S_PitInfo struct {
	PitKeyId int `json:"pitkeyid"` //!  关卡ID
}

type C2S_PitStart struct {
	PitKeyId int `json:"pitkeyid"` //!  key
	TeamType int `json:"teamtype"` //!  用哪一套阵容
}

type C2S_NewPitFinishEvent struct {
	Id               int                `json:"id"` //!  关卡ID
	HeroState        []*NewPitHeroState `json:"herostate"`
	IsFail           int                `json:"isfail"` //!  1表示攻打失败，只更新英雄血量和怪物血量
	MonsterHeroState []*NewPitHeroState `json:"monsterherostate"`
	Param            int                `json:"param"`
}

type C2S_NewPitFinishReward struct {
	Id int `json:"id"` //!  关卡ID
}

type C2S_NewPitShopBuy struct {
	Shoptype int `json:"shoptype"` //1,2,3
	Grid     int `json:"grid"`
}

type C2S_NewPitFinishSoul struct {
	Id    int `json:"id"` //!  关卡ID
	Param int `json:"param"`
}

type C2S_NewPitFinishEvil struct {
	Id               int                `json:"id"` //!  关卡ID
	HeroState        []*NewPitHeroState `json:"herostate"`
	IsFail           int                `json:"isfail"` //!  1表示攻打失败，只更新英雄血量和怪物血量
	MonsterHeroState []*NewPitHeroState `json:"monsterherostate"`
}

type C2S_NewPitFinishTreasure struct {
	Id               int                `json:"id"` //!  关卡ID
	HeroState        []*NewPitHeroState `json:"herostate"`
	IsFail           int                `json:"isfail"` //!  1表示攻打失败，只更新英雄血量和怪物血量
	MonsterHeroState []*NewPitHeroState `json:"monsterherostate"`
}

type C2S_NewPitFinishShop struct {
	Id    int `json:"id"` //!  关卡ID
	Param int `json:"param"`
}

type C2S_NewPitFinishSpring struct {
	Id int `json:"id"` //!  关卡ID
}

type C2S_NewPitFinishMystery struct {
	Id int `json:"id"` //!  关卡ID
}

//战斗关卡，对应类型2，3，4
type C2S_NewPitFinishBattle struct {
	Id               int                `json:"id"` //!  关卡ID
	HeroState        []*NewPitHeroState `json:"herostate"`
	IsFail           int                `json:"isfail"` //!  1表示攻打失败，只更新英雄血量和怪物血量
	MonsterHeroState []*NewPitHeroState `json:"monsterherostate"`
}

type C2S_NewPitNowAim struct {
	Id int `json:"id"` //!  关卡ID
}

type C2S_NewPitChooseBuff struct {
	Index int `json:"index"` //!  1~3
}

type C2S_NewPitFinishNow struct {
	Difficult int `json:"difficult"` //!  选择难度1，2  服务器找不到这个难度就给默认难度1
}

type C2S_SetClientSign struct {
	Key   int `json:"key"`   //!
	Value int `json:"value"` //!
}

//时光之巅
type C2S_InstanceStart struct {
	Id    int `json:"id"`    //!
	Force int `json:"force"` //!   0正常 1强制开启
}

type C2S_InstanceMove struct {
	Row       int   `json:"row"`       //!
	Col       int   `json:"col"`       //!
	RowShadow []int `json:"rowshadow"` //!
	ColShadow []int `json:"colshadow"` //!
	ThingId   int   `json:"thingid"`   //!
}

type C2S_FinishThing struct {
	ThingId int `json:"thingid"` //!
}

type C2S_InstanceSwitch struct {
	ThingId       int   `json:"thingid"`       //!
	SwitchState   int   `json:"switchstate"`   //!
	ThingIdEx     []int `json:"thingidex"`     //!
	SwitchStateEx []int `json:"switchstateex"` //!
}

type C2S_InstanceBattle struct {
	ThingId          int                  `json:"thingid"` //!  关卡ID
	HeroState        []*InstanceHeroState `json:"herostate"`
	IsFail           int                  `json:"isfail"` //!  1表示攻打失败，只更新英雄血量和怪物血量
	MonsterHeroState []*InstanceHeroState `json:"monsterherostate"`
	Fight            int64                `json:"fight"`
}

type C2S_InstanceFriend struct {
	ThingId int `json:"thingid"` //!  关卡ID
	Param   int `json:"param"`
}

type C2S_InstanceAdd struct {
	ThingId int `json:"thingid"` //!  关卡ID
}

type C2S_InstanceReborn struct {
	ThingId int `json:"thingid"` //!  关卡ID
}

type C2S_InstanceChooseBuff struct {
	ThingId int `json:"thingid"` //!  关卡ID
	Index   int `json:"index"`   //!  1~3
}

type C2S_InstanceMakeBuff struct {
	ThingId int `json:"thingid"` //!  关卡ID
}

//
type C2S_GetPvPFightInfo struct {
	FightId int `json:"fightid"` // 战报ID
}

type C2S_SupportHeroSet struct {
	Cid     string `json:"cid"`     // cid
	Index   int    `json:"index"`   // index
	HeroKey int    `json:"herokey"` // herokey
}

type C2S_SupportHeroCancel struct {
	Cid     string `json:"cid"`     // cid
	HeroKey int    `json:"herokey"` // herokey
}

type C2S_EntanglementUse struct {
	Cid       string `json:"cid"`       // cid
	Index     int    `json:"index"`     // 栏位
	Type      int    `json:"type"`      // 羁绊类型
	MasterUid int64  `json:"masteruid"` // 被借玩家的uid
	HeroKey   int    `json:"herokey"`   // herokey
}

//type C2S_EntanglementCancel struct {
//	Cid     string `json:"cid"`     // cid
//	Index   int    `json:"index"`   // 栏位
//	Type    int    `json:"type"`    // 羁绊类型
//}

type C2S_EntanglementAutoUse struct {
	Cid  string `json:"cid"`  // cid
	Type int    `json:"type"` // 羁绊类型
}

type C2S_SupportHeroInfo struct {
	Cid    string `json:"cid"`    //! cid
	HeroID int    `json:"heroid"` //! 英雄id
}

type C2S_RewardSet struct {
	Cid      string  `json:"cid"`      // cid
	ID       int     `json:"id"`       // 任务id
	Uids     []int64 `json:"uids"`     // 玩家的uid
	HeroKeys []int   `json:"herokeys"` // herokey
}

type C2S_RewardGet struct {
	Cid string `json:"cid"` // cid
	ID  int    `json:"id"`  // 任务id
}

type C2S_RewardGetAll struct {
	Cid  string `json:"cid"`  // cid
	Type int    `json:"type"` // cid
}

////! 排行信息
//type C2S_RankTaskTopInfo struct {
//	Cid   string `json:"cid"`
//	Index int    `json:"index"`
//	Ver   int    `json:"ver"`
//}

//! 获得领取信息
type C2S_RankTaskGetStateInfo struct {
	Cid string `json:"cid"`
}

//! 获得id类型完成的玩家数据
type C2S_RankTaskGetPlayerInfo struct {
	Cid string `json:"cid"`
	ID  int    `json:"id"`
}

//! 获得type类型完成的玩家数据
type C2S_RankTaskGetTypePlayerInfo struct {
	Cid  string `json:"cid"`
	Type int    `json:"type"`
}

//! 排行任务领取奖励
type C2S_RankTaskAward struct {
	Cid string `json:"cid"`
	ID  int    `json:"id"`
}

type C2S_ResonanceCrystalSet struct {
	Cid     string `json:"cid"`     // cid
	Index   int    `json:"index"`   // index
	HeroKey int    `json:"herokey"` // herokey
}

type C2S_ResonanceCrystalCancel struct {
	Cid     string `json:"cid"`     // cid
	HeroKey int    `json:"herokey"` // herokey
}
type C2S_ResonanceCrystalAddResonanceCount struct {
	Cid  string `json:"cid"`  // cid
	Type int    `json:"type"` // type
}

type C2S_ResonanceCrystalCleanCD struct {
	Cid   string `json:"cid"`   // cid
	Index int    `json:"index"` // index
}

//! 获得竞技场信息
type C2S_ArenaInfo struct {
	Cid string `json:"cid"`
}

//! 获得竞技场信息
type C2S_GetArenaFight struct {
	Cid string `json:"cid"`
}

//! 天下会武排行
type C2S_ArenaTop struct {
	Cid string `json:"cid"`
}

//! 战斗开始
type C2S_ArenaBegin struct {
	Index int `json:"index"`
}

//! 战斗开始
type C2S_ArenaFightBack struct {
	FightID int64 `json:"fightid"`
}

type C2S_GetArenaFightInfo struct {
}

type C2S_GetArenaDefend struct {
	FindUid int64 `json:"finduid"`
}

type C2S_GetEnemyFightInfo struct {
	Type    int   `json:"type"`
	FindUid int64 `json:"finduid"`
}

type C2S_BuyArenaCount struct {
	Count int `json:"count"`
}

type C2S_ArenaFightResult struct {
	Type       int         `json:"type"`
	BattleInfo *BattleInfo `json:"battleinfo"`
}

type C2S_GetBattleInfo struct {
	Cid     string `json:"cid"`
	FightID int64  `json:"fightid"`
	Type    int    `json:"type"`
	Uid     int    `json:"uid"`
}

type C2S_GetBattleInfo2 struct {
	Cid     string `json:"cid"`
	FightID int64  `json:"fightid"`
}

type C2S_GetBattleRecord struct {
	Cid     string `json:"cid"`
	FightID int64  `json:"fightid"`
	Type    int    `json:"type"`
	Uid     int    `json:"uid"`
}

//! 进入高阶竞技场
type C2S_ArenaSpecialEnter struct {
	Cid string `json:"cid"`
}

//! 获得竞技场敌人
type C2S_ArenaSpecialGetEnemy struct {
	Cid string `json:"cid"`
}

//! 开始战斗
type C2S_ArenaSpecialStartFight struct {
	Cid   string `json:"cid"`
	Index int    `json:"index"`
}

//! 领取奖励
type C2S_ArenaSpecialGetAward struct {
	Cid string `json:"cid"`
}

//! 获得战报
type C2S_ArenaSpecialGetFights struct {
	Cid string `json:"cid"`
}

type C2S_ArenaSpecialBuyCount struct {
	Count int `json:"count"`
}

type C2S_ArenaSpecialGetFightInfo struct {
	Type    int   `json:"type"`
	FindUid int64 `json:"finduid"`
	Rank    int   `json:"rank"`
}

type C2S_ActivityGiftGetAward struct {
	Cid string `json:"cid"`
	ID  int    `json:"id"`
}

type C2S_GrowthGiftGetAward struct {
	Cid string `json:"cid"`
	ID  int    `json:"id"`
}

type C2S_ArenaSpecialFightResult struct {
	BattleInfo [ARENA_SPECIAL_TEAM_MAX]BattleInfo `json:"battleinfo"`
}

//! 皮肤激活
type C2S_ActivateSkin struct {
	ID int `json:"id"` // 皮肤id
}

//! 皮肤设置
type C2S_SetSkin struct {
	ID        int `json:"id"`        // 皮肤id
	HeroIndex int `json:"heroindex"` // 英雄index
}

//! 获得信息
type C2S_SendInfo struct {
}

//! 获得限时抢购信息
type C2S_SpecialPurchaseInfo struct {
}

//! 提升主等级
type C2S_LifeTreeUpMainLevel struct {
}

//! 提升专业等级
type C2S_LifeTreeUpTypeLevel struct {
	Type int `json:"type"` // 类型
}

// 重置专业等级
type C2S_LifeTreeResetTypeLevel struct {
	Type int `json:"type"` // 类型
}

type C2S_NewChat struct {
	Channel    int    `json:"channel"`
	Content    string `json:"content"`
	PrivateUid int64  `json:"privateuid"`
	Language   int    `json:"language"`
}

type C2S_UnlockNebula struct {
	NebulaId int `json:"nebulaid"`
}

type C2S_GetNebulaWarBox struct {
	NebulaId int `json:"nebulaid"` //星云id
	GroupId  int `json:"groupid"`  //地图id
	BoxId    int `json:"boxid"`
}

type C2S_ActivityBossTask struct {
	Id     int `json:"id"`     //活动ID  1700～1709
	TaskId int `json:"taskid"` //任务ID 表里的taskid
}

type C2S_ActivityBossGetRecord struct {
	Id        int   `json:"id"` //活动ID  1700～1709
	TargetUid int64 `json:"targetuid"`
}

type C2S_ActivityBossExchange struct {
	ActivityId int `json:"activityid"` //活动id
	Id         int `json:"id"`         //任务id
}

type C2S_ActivityBossFight struct {
	Id int `json:"id"` //活动ID  1700～1709
}

type C2S_ActivityBossResultEx struct {
	Id         int         `json:"id"` //活动ID  1700～1709
	BattleInfo *BattleInfo `json:"battleinfo"`
}

type C2S_ActivityBossGetRank struct {
	Id int `json:"id"` //活动ID  1700～1709
}

type C2S_ActivityBossResetTimes struct {
	Id int `json:"id"` //活动ID  1700～1709
}

//! 上传限时神将积分到center服务器
type S2Center_UploadGeneral struct {
	Cid string          `json:"cid"`
	Top *Js_GeneralUser `json:"top"`
}

//! 请求跨服限时神将数据
type S2Center_GeneralRank struct {
	Cid      string `json:"cid"`
	ServerId int    `json:"serverid"`
}

//! 发送跨服神将数据
type S2Center_GeneralTop struct {
	Cid      string            `json:"cid"`      //! cid
	Rank     []*Js_GeneralUser `json:"rank"`     //! 前50名信息
	ServRank []*ServRank       `json:"servrank"` //! 50名外玩家信息
}

type C2S_HeroGrowFreeTask struct {
	ActivityType int `json:"activitytype"`
	Taskid       int `json:"taskid"`
}

type C2S_CrossArenaGetReward struct {
	Taskid int `json:"taskid"`
}

type C2S_CrossArenaAttack struct {
	AttackUid        int64 `json:"attackuid"`        //目标uid
	AttackSubsection int   `json:"attacksubsection"` //目标大段位
	AttackClass      int   `json:"attackclass"`      //目标小段位
}

type C2S_CrossArenaGetPlayerInfo struct {
	PlayerUid int64 `json:"playeruid"` //目标uid
}

type C2S_CrossArenaFightResult struct {
	Type       int         `json:"type"`
	BattleInfo *BattleInfo `json:"battleinfo"`
}

type C2S_ActivityBossFestivalResult struct {
	Id         int         `json:"id"`
	BattleInfo *BattleInfo `json:"battleinfo"`
	Score      int         `json:"score"` //万分比
}

type C2S_DoLotteryDraw struct {
	Times int `json:"times"` //!  抽奖次数
}

type C2S_LotteryDrawChangePrize struct {
	IsDefault int `json:"isdefault"` //0 仅这次选择，1作为默认选择
	Id        int `json:"id"`        //  选择id
}

type C2S_HonourShopBuy struct {
	Id int `json:"id"`
}

//! 召唤战马
type C2S_SummonHorse struct {
	Index int `json:"index"` //! 召唤类型，1-普通，2-高级
	Num   int `json:"num"`   //! 数量
}

//! 相马
type C2S_IdentifyHorse struct {
	Num int `json:"num"` //! 鉴定次数
}

//! 镶嵌马魂
type C2S_EmbedHorseSoul struct {
	HorseId int `json:horseid`  //! 坐骑Id
	SoulId  int `json:"soulid"` //! 马魂Id
	Index   int `json:"index"`  //! 孔编号，0-n
}

//! 升级马魂
type C2S_UpHorseSoul struct {
	HorseId int `json:horseid`  //! 坐骑Id
	SoulId  int `json:"soulid"` //! 马魂Id
	Index   int `json:"index"`  //! 孔编号，0-n
}

//! 分解坐骑-批量
type C2S_DecomposeHorse struct {
	Horselst []int `json:horselst`
}

//! 分解坐骑-批量
type C2S_DecomposeSoul struct {
	Soullst []JS_HorseSoulInfo `json:soullst`
}

//! 使用坐骑
type C2S_MountHorse struct {
	Heroid   int `json:"heroid"`
	Horseid  int `json:"horseid"`
	TeamType int `json:"team_type"`
	Index    int `json:"index"`
}

//! 使用坐骑
type C2S_UpHorse struct {
	Star     int   `json:"star"`
	Material []int `json:"material"`
}

//! 觉醒坐骑
type C2S_AwakeHorse struct {
	HorseId  int   `json:"horseid"`
	Materail []int `json:"materail"`
}

type C2S_CrossArena3V3GetReward struct {
	Taskid int `json:"taskid"`
}

type C2S_CrossArena3V3Attack struct {
	AttackUid        int64 `json:"attackuid"`        //目标uid
	AttackSubsection int   `json:"attacksubsection"` //目标大段位
	AttackClass      int   `json:"attackclass"`      //目标小段位
}

type C2S_CrossArena3V3GetPlayerInfo struct {
	PlayerUid int64 `json:"playeruid"` //目标uid
}

type C2S_CrossArena3V3FightResult struct {
	Type       int                                `json:"type"`
	BattleInfo [CROSSARENA3V3_TEAM_MAX]BattleInfo `json:"battleinfo"`
}

type C2S_GetRankRewardRank struct {
	Id int `json:"id"` // 1732
}

type C2S_GetRankRewardReward struct {
	Id int `json:"id"` // 1732
}

// 领取累抽奖励
type C2S_GetFindReward struct {
	RewardId int `json:"rewardid"` // 奖励id 0 1 2
}
