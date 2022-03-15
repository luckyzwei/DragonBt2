package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

const MsgKey = "7#0n29#%@ub9"

type San_UserBase struct {
	Uid                int64  // uid
	UName              string // 名字
	IconId             int    // icon
	Gem                int    // 钻石
	Gold               int    // 金币
	Exp                int    // 经验
	Level              int    // 等级
	Regtime            string // 注册时间
	LastLoginTime      string // 最后登陆时间
	LastUpdTime        string // 最后更新时间
	LastLiveTime       string // 最后活动时间
	Face               int    // 性别
	Vip                int    // Vip等级
	VipExp             int    // vip经验
	TiLi               int    // 体力
	PartyId            int    // 帮派id
	SkillPoint         int    // 技能点
	TiLiLastUpdataTime int64  // 体力最后更新时间
	SpLastUpdataTime   int64  //
	LastCheckinTime    string //
	CheckinNum         int    // 签到次数
	CheckinAward       int    // 累计签到奖励
	Position           int    //
	IsRename           int    // 是否改名
	LoginDays          int    // 登陆天数
	LoginAward         int    // 登陆奖励
	LevelAward         int    // 等级奖励
	WorldAward         int    // 世界奖励
	Morale             int    // 士气
	Citylevel          int    // 封地等级
	Camp               int    // 阵营
	City               int    // 所在城池
	Fight              int64  // 战斗力
	IsGag              int    // 是否禁言
	IsBlock            int    // 禁止登录
	BlockDay           int    // 冻结时间
	IP                 string // 最后登陆IP
	Promotebox         int    // 变强宝箱hao
	LineTime           int64  // 在线时长(秒)
	PassMax            int    // 最大过关
	FitServer          int    // 合服状态，0-表示未合服，1-合服，未领取奖励 2-合服，未领取奖励
	BlockTime          int64  // 冻结开始时间
	BlockReason        string // 冻结原因
	Soul               int    // 魂石
	TechPoint          int    // 科技点
	BossMoney          int    // 水晶币
	TowerStone         int    // 镇魂石
	Portrait           int    // 边框
	CampOk             int    // 阵营ok 0 不ok 1 ok
	NameOk             int    // 名字ok 0 不ok 1 ok
	GuildId            int    // 指引Id
	RedIcon            int    // 红色图标
	UserSignature      string // 用户签名
	GetAllGem          int    // 玩家获得的总钻石量
	PayGem             int    //! 付费获得元宝，元宝消耗时，优先消耗免费元宝，再消耗付费元宝

	DataUpdate
}

type CheckFightInfo struct {
	CheckId int64             `json:"checkid"` // id
	Info    map[int][]float64 `json:"info"`    //key:英雄key
}

const (
	CHECKFIGHT_ATTR_HP      = 1
	CHECKFIGHT_ATTR_ATTACK  = 2
	CHECKFIGHT_ATTR_DEFENCE = 3
	CHECKFIGHT_ATTR_END     = 4
)

const Interval = 900

// 玩家
type Player struct {
	ID             int64         // 唯一id
	SessionObj     *Session      // sessionid
	MsgTime        int64         // 上次消息收发时间
	TaskLock       *sync.RWMutex // 任务锁
	Sql_UserBase   San_UserBase  // userbase
	Module         *ModAll
	Nosend         bool
	OtherData      bool
	Account        *San_Account
	Platform       Platform_Info
	SaveTime       int64           // 上次保存时间
	MsgWaitSave    int64           // 等待保存时间
	SaveTimes      int             // 保存次数
	IsSave         bool            // 是否需要保存
	CheckCode      string          // 重连验证码
	CityPowerTime  int64           // 上次恢复军令时间
	ReloginTimes   int             // 重登次数
	NoticeBaseInfo bool            // 是否需要通知中心服更新基础信息
	NoticeChatTime int64           // 是否需要重连聊天服务器
	DebugHelp      string          // DEBUG
	CheckFight     *CheckFightInfo //战斗校验
}

func (self *Player) SafeClose() {
	LogInfo("安全踢下线")

	self.SendRet2("kickout01")

	self.SendMsg("shutdown", []byte(""))
}

func (self *Player) SetAccount(account *San_Account) {
	if account != nil {
		self.Account = account
		return
	}

	if self.Account != nil {
		return
	}

	account = new(San_Account)
	sql := fmt.Sprintf("select * from `san_account` where `uid` = %d", self.ID)
	GetServer().DBUser.GetOneData(sql, account, "", 0)
	self.Account = account
}

func (self *Player) SetSessionID(session *Session, login bool, ip string, os string) {
	if self.SessionObj == nil && session != nil { // 重新连接，记录登陆时间
		if !login { // 断线重连
			self.GetModule("task").(*ModTask).init = true
		}

		self.OtherPlayerData()
		self.LoginRefresh()
		self.Sql_UserBase.Update(true)
		self.GetModule("friend").(*ModFriend).SendOnline(1)

		self.CheckCode = GetRandomString(8)

		//unionid := self.GetModule("union").(*ModUnion).Sql_UserUnionInfo.Unionid
		//if unionid > 0 {
		//	GetUnionMgr().UpdateMemberState(unionid, self.Sql_UserBase.Uid)
		//}
	}
	self.SessionObj = session
	self.MsgTime = TimeServer().Unix()
	self.ReloginTimes = 0
	self.Sql_UserBase.IP = ip
	self.Platform.Platform = strings.ToLower(os)
}

func (self *Player) RunTime() {
	ticker := time.NewTicker(time.Minute)
	for {
		<-ticker.C
		if GetServer().ShutDown {
			ticker.Stop()
			return
		}
		// 15分钟没有活动并且连接断开
		GetPlayerMgr().Locker.Lock()
		if TimeServer().Unix()-self.MsgTime >= Interval && self.GetSession() == nil {
			//！ 关掉定时器
			self.Save(false, true)
			GetPlayerMgr().RemovePlayer(self.ID)
			GetPlayerMgr().Locker.Unlock() // break之前要解锁
			break
		}
		GetPlayerMgr().Locker.Unlock()

		// 超过10分钟无消息，并且未保存，则数据保存一次
		if TimeServer().Unix()-self.SaveTime > 600 {
			self.SaveTime = TimeServer().Unix()
			self.Save(false, true)
		}

		if TimeServer().Unix()-self.NoticeChatTime > 1200 {
			GetMasterMgr().EnterChat(self) //加入聊天频道

			var msg S2C_NewChat
			msg.Cid = "clearchatrecord"
			self.SendMsg(msg.Cid, HF_JtoB(&msg))

			GetMasterMgr().QueryWorldMessage(self)   //获得世界聊天
			GetMasterMgr().QueryUnionMessage(self)   //获得公会聊天
			GetMasterMgr().QueryPrivateMessage(self) //获得私聊
			self.NoticeChatTime = TimeServer().Unix()
		}
	}
	ticker.Stop()
}

func (self *Player) UpdateNoticeChatTime() {
	self.NoticeChatTime = TimeServer().Unix()
}

func (self *Player) OnTimer() {

}

func (self *Player) Save(shutdown bool, sql bool) {
	if !self.IsSave {
		return
	}

	if shutdown != GetServer().ShutDown {
		return
	}

	self.Sql_UserBase.Update(true)
	self.Module.Save(sql)
	self.IsSave = false
}

func (self *Player) InitPlayerData() {
	self.MsgTime = TimeServer().Unix()
	self.SaveTime = self.MsgTime
	self.ReloginTimes = 0
	self.Module.GetData()

	go self.RunTime()
}

func (self *Player) OtherPlayerData() {
	if self.OtherData {
		return
	}
	self.Module.GetOtherData()
	self.OtherData = true

	self.countTeamFight(ReasonPlayerLogin)
	//! 检查更新时间
	self.CheckLoginTime()
}

func (self *Player) CheckLoginTime() {
	updtime, _ := time.ParseInLocation(DATEFORMAT, self.Sql_UserBase.LastLiveTime, time.Local)
	if updtime.Unix()-TimeServer().Unix() > DAY_SECS {
		self.Sql_UserBase.LastLiveTime = self.GetNextRefreshTime()
	}
}

// 新建角色
func (self *Player) New() {
	// 送英雄和物品
	for _, value := range GetCsvMgr().NewUserItem {
		if value.Group == 1 {
			switch value.Type {
			case 1:
				heroId := value.Id
				(self.GetModule("hero").(*ModHero)).AddHero(heroId, 1, 0, 0, "创建角色赠送奖励")
				//self.GetModule("team").(*ModTeam).checkPos()
				self.countInit(ReasonPlayerLogin) // 计算战力

			case 2:
				self.AddObject(value.Id, value.Num, 0, 0, 0, "创建角色赠送奖励")
			}
		}
	}

	//20190624 by zy  为初始英雄分配阵容
	//self.GetModule("team").(*ModTeam).InitNew()

	// 默认给1000
	self.InitHead(DEFAULT_HEAD_ICON)

	// 送邮件，取消邮件赠送
	if false {
		lstItem := make([]PassItem, 0)
		lstItem = append(lstItem, PassItem{ITEM_GEM, 200})
		lstItem = append(lstItem, PassItem{ITEM_GOLD, 30000})
		(self.GetModule("mail").(*ModMail)).AddMail(1, 1, 0,
			GetCsvMgr().GetText("STR_LOGIN_MAIL_TITLE"),
			GetCsvMgr().GetText("STR_LOGIN_MAIL_CONTENT"),
			GetCsvMgr().GetText("STR_SYS"), lstItem, false, 0)

		self.HandleTask(TASK_TYPE_LOGIN_TOTAL_COUNT, 0, 0, 0)
	}

}

// 创建新的玩家
func NewPlayer(id int64) *Player {
	p := new(Player)
	p.ID = id
	p.Module = NewModAll(p)
	p.TaskLock = new(sync.RWMutex)
	p.OtherData = false

	return p
}

func (self *Player) InitUserBase(uid int64) {
	//player := NewPlayer(uid)
	self.Sql_UserBase.Uid = uid
	self.Sql_UserBase.UName = fmt.Sprintf(GetCsvMgr().GetText("STR_GUEST"), uid)
	self.Sql_UserBase.Portrait = 1000
	self.Sql_UserBase.IconId = 1002
	self.Sql_UserBase.Gold = 0
	self.Sql_UserBase.Gem = 0
	self.Sql_UserBase.Level = 1
	self.Sql_UserBase.Regtime = TimeServer().Format(DATEFORMAT)
	self.Sql_UserBase.LastLoginTime = "2017-01-01 00:00:00"
	self.Sql_UserBase.LastUpdTime = self.Sql_UserBase.Regtime
	self.Sql_UserBase.LastLiveTime = self.GetNextRefreshTime()
	self.Sql_UserBase.TiLi = GetCsvMgr().GetPhysicallimit(1)
	self.Sql_UserBase.SkillPoint = 20
	self.Sql_UserBase.TiLiLastUpdataTime = 0
	self.Sql_UserBase.SpLastUpdataTime = 0
	self.Sql_UserBase.LastCheckinTime = "2017-01-01 00:00:00"
	self.Sql_UserBase.Position = 10010
	self.Sql_UserBase.LoginDays = 0
	self.Sql_UserBase.Citylevel = 1
	self.Sql_UserBase.Camp = 1

	self.SetFight(0, 5)
}

func (self *Player) GetModule(name string) ModBase {
	return self.Module.GetModule(name)
}

func (self *Player) GetSession() *Session {
	return self.SessionObj
}

// 发送消息
func (self *Player) SendMsg(head string, body []byte) bool {
	if self.Nosend {
		return true
	}
	session := self.GetSession()
	if session == nil {
		return false
	}

	//log.Println("send msg:", string(body))

	session.SendMsg(head, body)

	// 发送任务
	self.SendTask()

	return true
}

// 得到下次刷新时间
func (self *Player) GetNextRefreshTime() string {
	now := TimeServer()
	if now.Hour() < 5 {
		return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local).Format(DATEFORMAT)
	} else {
		return time.Date(now.Year(), now.Month(), now.Day(), 24+5, 0, 0, 0, time.Local).Format(DATEFORMAT)
	}
}

func (self *Player) NoticeCenterBaseInfo() {
	self.NoticeBaseInfo = true
}

//相当于登录服务器后找中心服要的sendinfo
func (self *Player) NoticeCenter() {
	self.NoticeCenterBaseInfo()
	GetMasterMgr().UpdatePlayer(self)
	GetMasterMgr().EnterChat(self)           //加入聊天频道
	GetMasterMgr().QueryWorldMessage(self)   //获得世界聊天
	GetMasterMgr().QueryUnionMessage(self)   //获得公会聊天
	GetMasterMgr().QueryPrivateMessage(self) //获得私聊
}

func (self *Player) onReceive(head string, ctrl string, body []byte) {
	self.MsgTime = TimeServer().Unix()
	self.MsgWaitSave++

	switch ctrl {
	case "loginok": // 登陆成功
		self.GetModule("friend").(*ModFriend).SendInfo()
		self.GetModule("pass").(*ModPass).SendInfo()
		self.GetModule("equip").(*ModEquip).SendInfo()
		self.GetModule("beauty").(*ModBeauty).SendInfo()
		self.GetModule("horse").(*ModHorse).SendHorseInfo()
		self.GetModule("horse").(*ModHorse).SendHorseSoulInfo()
		self.GetModule("artifactequip").(*ModArtifactEquip).SendInfo()
		self.GetModule("crystal").(*ModResonanceCrystal).SendInfo([]byte{})
		self.GetModule("entanglement").(*ModEntanglement).SendInfo([]byte{})
		self.GetModule("lifetree").(*ModLifeTree).SendInfo([]byte{})
		self.GetModule("lifetree").(*ModLifeTree).SendRedPointInfo()
		self.GetModule("hero").(*ModHero).SendInfo()
		self.GetModule("find").(*ModFind).SendInfo()
		self.GetModule("recharge").(*ModRecharge).SendInfo("getrechargeinfo")
		self.GetModule("task").(*ModTask).SendInfo()
		self.GetModule("weekplan").(*ModWeekPlan).SendInfo()
		self.GetModule("shop").(*ModShop).SendInfo()
		self.GetModule("honourshop").(*ModHonourShop).SendInfo()
		self.GetModule("union").(*ModUnion).GetUserUnionInfo()
		self.GetModule("activity").(*ModActivity).SendInfo()
		self.GetModule("team").(*ModTeam).SendInfo()
		self.GetModule("luckshop").(*ModLuckShop).SendInfo()
		self.GetModule("fund").(*ModFund).SendInfo()
		self.GetModule("dailyrecharge").(*ModDailyRecharge).SendInfo(true)
		self.GetModule("newpit").(*ModNewPit).SendInfo()
		self.GetModule("viprecharge").(*ModVipRecharge).SendInfo()
		self.GetModule("instance").(*ModInstance).SendInfo()
		self.GetModule("head").(*ModHead).SendInfo()
		self.GetModule("moneytask").(*ModMoneyTask).SendInfo()
		self.GetModule("guide").(*ModGuide).SendInfo()
		self.GetModule("tower").(*ModTower).sendInfo()
		self.GetModule("reward").(*ModReward).SendInfo([]byte{})
		self.GetModule("timegift").(*ModTimeGift).SendInfo()
		//self.GetModule("hydra").(*ModHydra).SendInfo()
		self.GetModule("support").(*ModSupportHero).HeroSupportMyHero([]byte{})
		self.GetModule("support").(*ModSupportHero).HeroSupportInfo([]byte{})
		self.GetModule("skin").(*ModSkin).SendInfo([]byte{})
		self.GetModule("specialpurchase").(*ModSpecialPurchase).SendInfo()
		self.GetModule("onhook").(*ModOnHook).SendInfoAutoSend()
		self.GetModule("targettask").(*ModTargetTask).SendInfo()
		self.GetModule("nobilitytask").(*ModNobilityTask).SendInfo()
		self.GetModule("turntable").(*ModTurnTable).SendInfo()
		self.GetModule("accesscard").(*ModAccessCard).SendInfo()
		self.GetModule("clientsign").(*ModClientSign).SendInfo()
		self.GetModule("interstellar").(*ModInterStellar).SendInfo()
		self.GetModule("activityboss").(*ModActivityBoss).SendInfo()
		self.GetModule("general").(*ModGeneral).SendInfo()
		self.GetModule("herogrow").(*ModHeroGrow).SendInfo()
		self.GetModule("crossarena").(*ModCrossArena).SendInfo()
		self.GetModule("crossarena3v3").(*ModCrossArena3V3).SendInfo()
		self.GetModule("activitybossfestival").(*ModActivityBossFestival).SendInfo()
		self.GetModule("lotterydraw").(*ModLotteryDraw).SendInfo()
		self.GetModule("consumertop").(*ModConsumerTop).SendInfo()
		self.GetModule("onhook").(*ModOnHook).CheckPass()
		self.GetModule("mail").(*ModMail).SendInfo()
		self.GetModule("bag").(*ModBag).SendInfo()

		GetOfflineInfoMgr().GetInfo(self)
		GetServer().sendLog_LoginOK(self)

		GetArenaMgr().UpdateFormat(self)
		GetArenaSpecialMgr().UpdateFormat(self)

		// 刷新钻石累消令
		self.GetModule("recharge").(*ModRecharge).CalWarOrderLimit(WARORDERLIMIT_3)

		// 发送签到奖励信息
		self.GetCheckinAwardInfo()

		if self.Sql_UserBase.NameOk == LOGIC_TRUE {
			self.NoticeCenter()
		}

		var msg S2C_LoginRet
		msg.Cid = "loginokret"
		msg.Ret = 0
		msg.CheckCode = self.CheckCode
		msg.Servertime = TimeServer().Unix()
		smsg, _ := json.Marshal(&msg)
		self.SendMsg("loginokret", smsg)

		self.CheckTask()
		GetServer().SendLog_SDKUP_AIWAN_LOGIN(self, SKDUP_ADDR_URL_AIWAN_SDK_EVENT_ENTERSVR)
		return
	case "onready":
		self.GetModule("entanglement").(*ModEntanglement).OnReady()
		return
	case "savejq":
		var msg C2S_JqInfo
		json.Unmarshal(body, &msg)
		self.SetJqzy(msg.Jqid, -1)
		return
	case "savezy":
		var msg C2S_ZyInfo
		json.Unmarshal(body, &msg)
		self.SetJqzy(-1, msg.Zyid)
		return
	case "savezy1":
		var msg C2S_ZyInfo
		json.Unmarshal(body, &msg)
		self.SetJqzy1(-1, msg.Zyid)
	case "alertname":
		var msg C2S_AlertName
		json.Unmarshal(body, &msg)
		self.Alertname(msg.Newname)
		return
	case "alerticon":
		var msg C2S_AlertIcon
		json.Unmarshal(body, &msg)
		self.Alerticon(msg.Icon)
		return
	case "createrole":
		var msg C2S_CreateRole
		json.Unmarshal(body, &msg)
		self.CreateRole(msg.Name, msg.Icon, msg.Face)
		return
	case "addonesp":
		var msg S2C_UpdateSP
		if self.Sql_UserBase.SpLastUpdataTime != 0 {
			msg.Cid = "autoaddsp"
		} else {
			msg.Cid = "updatasp"
		}
		msg.Uid = self.ID
		msg.Sp = self.GetSkillPoint()
		if self.Sql_UserBase.SpLastUpdataTime == 0 {
			msg.Time = 0
		} else {
			msg.Time = ADDSPTIME - int(TimeServer().Unix()-self.Sql_UserBase.SpLastUpdataTime)
		}
		smsg, _ := json.Marshal(&msg)
		self.SendMsg("autoaddsp", smsg)
		return
	case "addone":
		var msg S2C_UpdateTiLi
		if self.Sql_UserBase.TiLiLastUpdataTime != 0 {
			msg.Cid = "autoaddenergy"
		} else {
			msg.Cid = "updataenergy"
		}
		msg.Uid = self.ID
		msg.Tili = self.GetPower()
		if self.Sql_UserBase.TiLiLastUpdataTime == 0 {
			msg.Time = 0
		} else {
			msg.Time = ADDPOWERTIME - int(TimeServer().Unix()-self.Sql_UserBase.TiLiLastUpdataTime)
		}
		smsg, _ := json.Marshal(&msg)
		self.SendMsg("autoaddenergy", smsg)
		return
	case "checkin":
		//self.SendRet4("checkin", self.Checkin())
		self.Checkin()
		return
	case "checkinaward":
		var c2s_msg C2S_CheckinAward
		json.Unmarshal(body, &c2s_msg)
		var msg S2C_CheckinAward
		msg.Cid = "checkinaward"
		msg.Ret, msg.Item = self.CheckinAward(c2s_msg.Index)
		msg.Newvalue = self.Sql_UserBase.CheckinAward
		smsg, _ := json.Marshal(&msg)
		self.SendMsg("checkinaward", smsg)

		self.SendInfo("updateuserinfo")
		return
	case "upcitylevel":
		nextlevel := self.Sql_UserBase.Citylevel + 1
		citycsv, ok := GetCsvMgr().MaincityConfig[nextlevel]
		if !ok {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return
		}

		if self.Sql_UserBase.Level < citycsv.Mainlevel {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_LEVEL_FAIL"))
			return
		}

		needid := citycsv.Consumption
		num := citycsv.Number
		if self.GetObjectNum(needid) < num {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_COST_FAIL"))
			return
		}

		self.AddObject(needid, -num, 17, 0, 0, "封地提升等级")
		self.Sql_UserBase.Citylevel = nextlevel
		self.SendRet("upcitylevel", 0)
		return
	case "checkfresh":
		self.Checkfresh()
		return
	case "fitaward":
		self.FitAward()
		return
	case "updateuserinfo":
		self.SendInfo("updateuserinfo")
		return
	case "set_guild_id":
		var msg C2S_SetGuild
		json.Unmarshal(body, &msg)
		self.Sql_UserBase.GuildId = msg.GuildId
		GetServer().SqlLog(self.Sql_UserBase.Uid, LOG_GUIDE_START, msg.GuildId, 0, 0, "引导", 0, 0, self)
		return
	case "set_redicon":
		var msg C2S_SetRedIcon
		json.Unmarshal(body, &msg)
		self.SetRedIcon(msg.Id)
		return
	case "get_redicon":
		var msg S2C_SetGuild
		msg.RedIcon = self.Sql_UserBase.RedIcon
		self.SendMsg("get_redicon", HF_JtoB(&msg))
		return
	case "syn_redicon":
		var msg S2C_SynGuild
		msg.RedIcon = self.Sql_UserBase.RedIcon
		self.SendMsg("get_redicon", HF_JtoB(&msg))
		return
	case "set_user_signature":
		var msg S2C_SetUserSignature
		json.Unmarshal(body, &msg)
		self.Sql_UserBase.UserSignature = msg.Signature
		GetOfflineInfoMgr().SetPlayerSignature(self.Sql_UserBase.Uid, self.Sql_UserBase.UserSignature)
		msg.Cid = "set_user_signature"
		self.SendMsg("set_user_signature", HF_JtoB(&msg))
		return
	case "checkfight":
		self.CheckFightMsg(body)
		return
	case "getrankrewardrank":
		self.GetRankRewardRank(body)
		return
	case "getrankrewardreward":
		self.GetRankRewardReward(body)
		return
	}

	self.Module.OnMsg(ctrl, body)

	//LogDebug("调试Msg信息：", self.MsgWaitSave, self.MsgTime, self.SaveTime)
	if self.MsgWaitSave > 100 || self.MsgTime-self.SaveTime > 180 {
		self.SaveTime = self.MsgTime
		self.MsgWaitSave = 0
		self.SaveTimes++
		if self.SaveTimes < 3 {
			self.Save(false, false)
		} else {
			self.SaveTimes = 0
			self.Save(false, true)
		}
		if self.NoticeBaseInfo {
			GetMasterMgr().UpdatePlayer(self)
		}
	}
}

// 设置小红点
func (self *Player) SetRedIcon(id int) {
	self.Sql_UserBase.RedIcon = id
}

// 客户端断开连接，即时保存数据
func (self *Player) onClose() {
	self.SessionObj = nil
	//通知中心服设置为离线状态。
	GetMasterMgr().ExitChat(self)
	GetMasterMgr().SetPlayerOffline(self)

	LogInfo("副本退出通知ok")
	self.GetModule("friend").(*ModFriend).SendOnline(0)
	//self.SessionObj = nil
	LogInfo("设置玩家离线ok")
	//GetUnionMgr().UpdateMemberState(self.GetModule("union").(*ModUnion).Sql_UserUnionInfo.Unionid, self.Sql_UserBase.Uid)
	LogInfo("更新军团数据ok")
	//self.GetModule("armsarena").(*ModPvp).UpdateFightInfo(true)
	LogInfo("更新竞技场战斗数据ok")
	self.Sql_UserBase.LastUpdTime = TimeServer().Format(DATEFORMAT)
	GetOfflineInfoMgr().SetPlayerOffTime(self.Sql_UserBase.Uid, TimeServer().Unix())
	tll, _ := time.ParseInLocation(DATEFORMAT, self.Sql_UserBase.LastLoginTime, time.Local)
	self.Sql_UserBase.LineTime += TimeServer().Unix() - tll.Unix()
	self.Save(false, true)
	LogInfo("保存玩家数据ok")
	GetServer().SqlLineLog(self.Sql_UserBase.Uid, self.Sql_UserBase.IP, int(TimeServer().Unix()-tll.Unix()), self.Account.Creator)
	LogInfo("保存在线数据库日志ok")
	GetServer().sendLog_Offline(self, tll.UnixNano()/1e6, TimeServer().UnixNano()/1e6, TimeServer().Unix()-tll.Unix())
	LogInfo("保存在线经分日志ok")
	LogInfo("玩家onClose成功")
	AddSdkOfflineLog(self)
	//GetServer().SendLog_SDKUP_Offline(self)
	self.OtherData = false
}

// 客户端要求 检测刷新
func (self *Player) Checkfresh() {

	// 服务器检测是否需要 刷新
	rtime, _ := time.ParseInLocation(DATEFORMAT, self.Sql_UserBase.LastLiveTime, time.Local)
	if TimeServer().Unix() < rtime.Unix() {
		return
	}

	//! 登录相关更新-不需要重新登录，直接更新
	self.LoginRefresh()

	self.Refresh()

	self.SendInfo("updateuserinfo")

	self.GetModule("friend").(*ModFriend).SendInfo()
	self.GetModule("pass").(*ModPass).SendInfo()
	self.GetModule("hero").(*ModHero).SendInfo()
	self.GetModule("find").(*ModFind).SendInfo()
	self.GetModule("recharge").(*ModRecharge).SendInfo("getrechargeinfo")
	self.GetModule("task").(*ModTask).SendInfo()
	self.GetModule("weekplan").(*ModWeekPlan).SendInfo()
	self.GetModule("shop").(*ModShop).SendInfo()
	self.GetModule("honourshop").(*ModHonourShop).SendInfo()
	self.GetModule("union").(*ModUnion).GetUserUnionInfo()
	self.GetModule("activity").(*ModActivity).SendInfo()
	self.GetModule("luckshop").(*ModLuckShop).SendInfo()
	self.GetModule("fund").(*ModFund).SendInfo()
	self.GetModule("reward").(*ModReward).SendInfo([]byte{})
	self.GetModule("activityboss").(*ModActivityBoss).SendInfo()
	self.GetModule("specialpurchase").(*ModSpecialPurchase).SendInfo()
	self.GetModule("general").(*ModGeneral).SendInfo()
	self.GetModule("herogrow").(*ModHeroGrow).SendInfo()
	self.GetModule("crossarena").(*ModCrossArena).SendInfo()
	self.GetModule("crossarena3v3").(*ModCrossArena3V3).SendInfo()
	self.GetModule("activitybossfestival").(*ModActivityBossFestival).SendInfo()
	(self.GetModule("horse").(*ModHorse)).SendHorseInfo()
	(self.GetModule("horse").(*ModHorse)).SendHorseSoulInfo()

	self.GetModule("mail").(*ModMail).SendInfo()
	self.GetModule("bag").(*ModBag).SendInfo()

	self.SendRet("checkfresh", 0)
}

// 合服奖励
func (self *Player) FitAward() {
	if self.Sql_UserBase.FitServer == 2 {
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_PLAYER_CONFORMITY_AWARD_HAS_BEEN_RECEIVED"))
		self.SendRet("fitaward", 1)
	} else if self.Sql_UserBase.FitServer == 1 {
		if self.Sql_UserBase.Level < 30 {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_PLAYER_PLAYER_LEVEL_IS_INSUFFICIENT_TO"))
			return
		}

		var msg S2C_MailAllItem
		msg.Cid = "mailallitem"

		csv, ok := GetCsvMgr().FitConfig[1]
		if !ok {
			self.SendRet("fitaward", 2)
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_PLAYER_CANT_FIND_A_CONFORMITY_AWARD"))
			return
		}

		msg.Item = make([]PassItem, 0)
		for i := 0; i < len(csv.Items); i++ {
			itemid := csv.Items[i]
			itemnum := csv.Nums[i]
			if itemid > 0 && itemnum > 0 {
				itemid, itemnum = self.AddObject(itemid, itemnum, 0, 0, 0, "合服奖励")
				msg.Item = append(msg.Item, PassItem{ItemID: itemid, Num: itemnum})
			}
		}

		self.Sql_UserBase.FitServer = 2
		self.SendMsg("mailallitem", HF_JtoB(&msg))
		self.SendRet("fitaward", 0)
	}
}

func (self *Player) CheckRefresh() bool {
	rtime1, _ := time.ParseInLocation(DATEFORMAT, self.Sql_UserBase.LastLiveTime, time.Local)
	nextTime := HF_GetNextDayStart()
	if rtime1.Unix() > nextTime {
		self.Sql_UserBase.LastLiveTime = self.GetNextRefreshTime()
	}

	rtime, _ := time.ParseInLocation(DATEFORMAT, self.Sql_UserBase.LastLiveTime, time.Local)
	if TimeServer().Unix() >= rtime.Unix() {
		self.SendRet("loginrefresh", 1)
		return true
	}
	return false
}

// 刷新,凌晨5点
func (self *Player) Refresh() {
	rtime, _ := time.ParseInLocation(DATEFORMAT, self.Sql_UserBase.LastLiveTime, time.Local)
	if TimeServer().Unix() < rtime.Unix() {
		return
	}

	self.GetModule("pass").(*ModPass).OnRefresh()
	self.GetModule("task").(*ModTask).OnRefresh()
	self.GetModule("activityboss").(*ModActivityBoss).OnRefresh()
	self.GetModule("activitybossfestival").(*ModActivityBossFestival).OnRefresh()
	self.GetModule("friend").(*ModFriend).Refresh()
	self.GetModule("shop").(*ModShop).Refresh()
	self.GetModule("viprecharge").(*ModVipRecharge).OnRefresh()
	self.GetModule("honourshop").(*ModHonourShop).OnRefresh()
	self.GetModule("union").(*ModUnion).OnRefresh()
	self.GetModule("onhook").(*ModOnHook).OnRefresh()
	self.GetModule("activity").(*ModActivity).OnRefresh()
	self.GetModule("recharge").(*ModRecharge).OnRefresh()
	self.GetModule("redpac").(*ModRedPac).OnRefresh()
	self.GetModule("dailyrecharge").(*ModDailyRecharge).OnRefresh()
	self.GetModule("equip").(*ModEquip).OnRefresh()
	self.GetModule("team").(*ModTeam).SendInfo()
	self.GetModule("moneytask").(*ModMoneyTask).OnRefresh()
	self.GetModule("find").(*ModFind).OnRefresh()
	self.GetModule("hero").(*ModHero).OnRefresh() //刷新每日重生次数   20190413
	self.GetModule("reward").(*ModReward).OnRefresh(false)
	self.GetModule("activitygift").(*ModActivityGift).OnRefresh()
	self.GetModule("tower").(*ModTower).OnRefresh()
	self.GetModule("specialpurchase").(*ModSpecialPurchase).OnRefresh()
	self.GetModule("general").(*ModGeneral).OnRefresh()
	self.GetModule("crossarena").(*ModCrossArena).OnRefresh()
	self.GetModule("crossarena3v3").(*ModCrossArena3V3).OnRefresh()
	self.GetModule("consumertop").(*ModConsumerTop).OnRefresh()
	self.GetModule("beauty").(*ModBeauty).OnRefresh()
	self.GetModule("horse").(*ModHorse).OnRefresh()

	self.HandleTask(0, 0, 0, 0)

	self.Sql_UserBase.LastLiveTime = self.GetNextRefreshTime()

	//LogDebug("玩家每日凌晨5点刷新")
}

// 用户登录时候,做刷新用
func (self *Player) LoginRefresh() {
	now := TimeServer()

	//刷新签到，每月循环
	if HF_IsNewDate(self.Sql_UserBase.LastCheckinTime) == true {
		lct, _ := time.ParseInLocation(DATEFORMAT, self.Sql_UserBase.LastCheckinTime, time.Local)
		if lct.Year() != now.Year() || lct.Month() != now.Month() {
			self.Sql_UserBase.CheckinNum = 0
			self.Sql_UserBase.CheckinAward = 0
		}
	}

	// 刷新登录时间
	if HF_IsNewDate(self.Sql_UserBase.LastLoginTime) == true {
		self.Sql_UserBase.LoginDays += 1
		self.HandleTask(TASK_TYPE_LOGIN_TOTAL_COUNT, 0, 0, 0)
		self.HandleTask(TASK_TYPE_REG_TOTAL_COUNT, 0, 0, 0)
		//self.HandleTask(126, 0, 0, 0)
		self.HandleTask(TASK_TYPE_LOGIN_DAY, 1, 0, 0)
		//self.HandleTask(50, 0, 0, 0)
	}
	self.Sql_UserBase.LastLoginTime = now.Format(DATEFORMAT)

	// 每日刷新
	self.Refresh()
}

// 是否签到, 发给客户端显示
func (self *Player) IsCheckin() bool {
	tNow := TimeServer()
	checkInToday := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 5, 0, 0, 0, TimeServer().Location()).Unix()	// 应该签到的时间
	if tNow.Hour() < 5 { // 5点前则判定前一天
		checkInToday -= DAY_SECS
	}

	lct, _ := time.ParseInLocation(DATEFORMAT, self.Sql_UserBase.LastCheckinTime, time.Local)
	if lct.Unix() >= checkInToday {
		return true
	} else {
		return false
	}
}


// 签到奖励展示
func (self *Player) GetCheckinAwardInfo()  {

	// 当前奖励
	var items []PassItem
	for _, config := range GetCsvMgr().SignConfig {
		if config != nil {
			items = append(items,PassItem{config.Reward,config.Number})
		}else{
			return
		}
	}
	//for i :=0 ;i < 30; i++{
	//	items = append(items,PassItem{30020088,1})
	//}
	// 签到超过30次，则重置
	if self.Sql_UserBase.CheckinNum > 30{
		self.Sql_UserBase.CheckinNum = 0
	}
	var msg S2C_CheckinAwardInfo
	msg.Cid = "checkinawardinfo"
	msg.Items = items
	msg.CheckinNum = self.Sql_UserBase.CheckinNum
	if self.IsCheckin(){
		msg.CheckinState = LOGIC_TRUE
	}
	smsg, _ := json.Marshal(&msg)
	self.SendMsg(msg.Cid, smsg)

}

// 签到
func (self *Player) Checkin() bool {
	if self.IsCheckin() {	// 今日已签到
		return false
	}

	//// 获取下一次签到奖励配置
	//config, ok := GetCsvMgr().GetSign(int(TimeServer().Month()), self.Sql_UserBase.CheckinNum+1)
	//if !ok {
	//	return false
	//}
	// 获取签到奖励配置
	var config *SignConfig
	for _, config = range GetCsvMgr().SignConfig {
		if config.Sign == self.Sql_UserBase.CheckinNum+1 {
			break
		}
	}

	bs := 1		// 奖励基础倍率
	needvip := config.Vip

	if needvip > 0 && self.Sql_UserBase.Vip >= needvip {
		bs = 2		// 签到翻倍
	}

	// 增加签到奖励
	self.AddObject(config.Reward, config.Number*bs, config.Sign, 0, 0, "签到")

	// 签到时间
	self.Sql_UserBase.LastCheckinTime = TimeServer().Format(DATEFORMAT)
	self.Sql_UserBase.CheckinNum++

	// 纠正签到天数 10号那天的签到次数不能大于10吧
	day := TimeServer().Day()
	if self.Sql_UserBase.CheckinNum > day+1 && day > 1 {
		// 签到未刷新，直接重置为当前天,奖励领取情况重置
		if TimeServer().Hour() < 5 {
			self.Sql_UserBase.CheckinNum = day - 1
			self.Sql_UserBase.CheckinAward = 0
		} else {
			self.Sql_UserBase.CheckinNum = day
			self.Sql_UserBase.CheckinAward = 0
		}
	}

	var msg S2C_CheckinToday
	msg.Cid = "checkintoday"
	msg.Item = PassItem{config.Reward,config.Number*bs}
	msg.CheckinNum = self.Sql_UserBase.CheckinNum
	msg.CheckinState = LOGIC_TRUE
	smsg, _ := json.Marshal(&msg)
	self.SendMsg(msg.Cid, smsg)

	//self.SendInfo("updateuserinfo")

	//GetServer().SqlLog(self.Sql_UserBase.Uid, LOG_USER_SIGN, config.Sign, 0, 0, "签到", 0, 0, self)

	return true
}

// 领取签到奖励
// 1:已经领取过  2：签到不足 -1:操作失败
func (self *Player) CheckinAward(index int) (int, []PassItem) {
	var outitem []PassItem

	base := 1
	if index > 0 {
		base = 1 << uint(index-1)
	}

	if self.Sql_UserBase.CheckinAward&base != 0 {
		return 1, outitem
	}

	csv_signreward, ok := GetCsvMgr().SignrewardConfig[index]
	if !ok {
		return -1, outitem
	}

	if self.Sql_UserBase.CheckinNum < csv_signreward.Signnum {
		return 2, outitem
	}

	for i := 0; i < len(csv_signreward.Rewarditems); i++ {
		itemid := csv_signreward.Rewarditems[i]
		if itemid == 0 {
			continue
		}
		num := csv_signreward.Rewardnums[i]
		itemid, num = self.AddObject(itemid, num, 16, index, 0, "签到宝箱")
		outitem = append(outitem, PassItem{itemid, num})
	}
	self.Sql_UserBase.CheckinAward += base

	return 0, outitem
}

// 加金币
func (self *Player) AddGold(num int, param1, param2 int, dec string) {
	if num == 0 {
		return
	}

	self.Sql_UserBase.Gold = HF_MinInt(HF_MaxInt(self.Sql_UserBase.Gold+num, 0), 2100000000)
	GetServer().SqlLog(self.Sql_UserBase.Uid, DEFAULT_GOLD, num, param1, param2, dec, self.Sql_UserBase.Gold, 0, self)
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "金币", self.Sql_UserBase.Gold) // _,_,来源，类型，剩余数量
	} else if num < 0 {
		self.HandleTask(CostEnergeTask, -num, 0, DEFAULT_GOLD)
		GetServer().sendLog_UseMoney(self, num, dec, "金币", self.Sql_UserBase.Gold) // _,_,来源，类型，剩余数量
	}
}

//增加能量召唤点数
func (self *Player) AddLootEnergy(num int) {
	if num == 0 {
		return
	}

	//self.GetModule("find").(*ModFind).AddLootEnergy(num)
}

// 加钻石
func (self *Player) AddGem(num int, param1, param2, param3 int, dec string) {
	if num == 0 {
		return
	}

	//! 付费元宝扣除
	costPayGem := 0
	if self.Sql_UserBase.PayGem > 0 && num < 0 {
		costPayGem = HF_MaxInt(self.Sql_UserBase.PayGem-self.Sql_UserBase.Gem-num, 0)
		self.Sql_UserBase.PayGem -= costPayGem
	}

	if num < 0 {
		GetRankRewardMgr().UpdateScore(ACT_RANKREWARD_COST, self, int64(-num))
	}

	self.Sql_UserBase.Gem = HF_MinInt(HF_MaxInt(self.Sql_UserBase.Gem+num, 0), 2100000000)
	GetServer().SqlLog(self.Sql_UserBase.Uid, DEFAULT_GEM, num, param1, param2, dec, self.Sql_UserBase.Gem, param3, self, costPayGem)
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "钻石", self.Sql_UserBase.Gem) // _,_,来源，类型，剩余数量
		self.Sql_UserBase.GetAllGem += num
		//GetServer().SendLog_SDKUP_MoneyChange(self, "1", "1", dec, num, self.Sql_UserBase.Gem, 0, "", 0)
	} else if num < 0 {
		GetServer().sendLog_UseMoney(self, num, dec, "钻石", self.Sql_UserBase.Gem) // _,_,来源，类型，剩余数量
	}
	// 处理太特殊了 只有这一个活动不计入消耗任务中 但是活动又是写在通用接口里的 为它单独写一个太费事了
	if num < 0 && dec != "招财猫" {
		// 任务
		self.HandleTask(TASK_TYPE_RECHARGE_COST, -num, 0, 0)
		//self.HandleTask(CostEnergeTask, -num, 0, DEFAULT_GEM)
	}
}

// 加功勋
func (self *Player) AddFeats(num int, param1, param2 int, dec string) {

}

func (self *Player) AddHonor(num int, param1, param2 int, dec string) {
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "功勋", self.GetObjectNum(91000018)+num)
	} else if num < 0 {
		GetServer().sendLog_UseMoney(self, num, dec, "功勋", self.GetObjectNum(91000018)+num)
	}
}

// 得到体力,使用函数,需要计算自动恢复
func (self *Player) GetPower() int {
	self.AutoPower()
	self.CheckPowerTime()

	return self.Sql_UserBase.TiLi
}

// 加体力
func (self *Player) AddPower(num int, param1, param2 int, dec string) {
	if num == 0 {
		return
	}

	if num < 0 {
		self.HandleTask(CostEnergeTask, -num, 0, 91000003)
	}

	//if num > POWERMAX {
	//self.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_POWERMAX"))
	//return
	//}

	tmp := self.Sql_UserBase.TiLi
	tmptime := self.Sql_UserBase.TiLiLastUpdataTime

	self.AutoPower()
	//修改体力可以突破一次1500，第二次失效
	//if num > 0 && self.Sql_UserBase.TiLi >= POWERMAX {
	//self.SendErrInfo("err", GetCsvMgr().GetText("STR_FRIEND_POWERMAX"))
	//return
	//} else {
	self.Sql_UserBase.TiLi = HF_MaxInt(self.Sql_UserBase.TiLi+num, 0)
	//}
	self.CheckPowerTime()

	if self.Sql_UserBase.TiLi != tmp || self.Sql_UserBase.TiLiLastUpdataTime != tmptime {
		var msg S2C_UpdateTiLi
		msg.Cid = "updataenergy"
		msg.Uid = self.ID
		msg.Tili = self.Sql_UserBase.TiLi
		if self.Sql_UserBase.TiLiLastUpdataTime == 0 {
			msg.Time = 0
		} else {
			msg.Time = ADDPOWERTIME - int(TimeServer().Unix()-self.Sql_UserBase.TiLiLastUpdataTime)
		}
		smsg, _ := json.Marshal(&msg)
		self.SendMsg("updataenergy", smsg)
	}
	GetServer().SqlLog(self.Sql_UserBase.Uid, 91000003, num, param1, param2, dec, self.Sql_UserBase.TiLi,
		0, self)
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "体力", self.Sql_UserBase.TiLi) // _,_,来源，类型，剩余数量
	} else if num < 0 {
		GetServer().sendLog_UseMoney(self, num, dec, "体力", self.Sql_UserBase.TiLi) // _,_,来源，类型，剩余数量
	}
}

// 计算体力回复
func (self *Player) AutoPower() {
	// 满了，不刷新
	if self.Sql_UserBase.TiLiLastUpdataTime <= 0 {
		return
	}

	// 经过时间
	stime := int(TimeServer().Unix() - self.Sql_UserBase.TiLiLastUpdataTime)

	// 回复
	add := stime / ADDPOWERTIME
	sub := stime % ADDPOWERTIME

	self.Sql_UserBase.TiLiLastUpdataTime = TimeServer().Unix() - int64(sub)

	if add <= 0 {
		return
	}

	max := GetCsvMgr().GetPhysicallimit(self.Sql_UserBase.Level)
	if self.Sql_UserBase.TiLi >= max {
		return
	}

	self.Sql_UserBase.TiLi = HF_MinInt(max, self.Sql_UserBase.TiLi+add)
}

func (self *Player) CheckPowerTime() {
	max := GetCsvMgr().GetPhysicallimit(self.Sql_UserBase.Level)
	if self.Sql_UserBase.TiLiLastUpdataTime > 0 && self.Sql_UserBase.TiLi >= max {
		self.Sql_UserBase.TiLiLastUpdataTime = 0
	} else if self.Sql_UserBase.TiLiLastUpdataTime <= 0 && self.Sql_UserBase.TiLi < max {
		self.Sql_UserBase.TiLiLastUpdataTime = TimeServer().Unix()
	}
}

// 得到技能点,使用函数,需要计算自动恢复
func (self *Player) GetSkillPoint() int {
	self.AutoSkillPoint()
	self.CheckSkillPointTime()

	return self.Sql_UserBase.SkillPoint
}

// 加技能点
func (self *Player) AddSkillPoint(num int, param1, param2 int, dec string) {
	if num == 0 {
		return
	}

	tmp := self.Sql_UserBase.SkillPoint
	tmptime := self.Sql_UserBase.SpLastUpdataTime

	self.AutoSkillPoint()
	self.Sql_UserBase.SkillPoint = HF_MinInt(HF_MaxInt(self.Sql_UserBase.SkillPoint+num, 0), SKILLPOINTMAX)
	self.CheckSkillPointTime()

	if self.Sql_UserBase.SkillPoint != tmp || self.Sql_UserBase.SpLastUpdataTime != tmptime {
		var msg S2C_UpdateSP
		msg.Cid = "updatasp"
		msg.Uid = self.ID
		msg.Sp = self.Sql_UserBase.SkillPoint
		if self.Sql_UserBase.SpLastUpdataTime == 0 {
			msg.Time = 0
		} else {
			msg.Time = ADDSPTIME - int(TimeServer().Unix()-self.Sql_UserBase.SpLastUpdataTime)
		}
		//smsg, _ := json.Marshal(&msg)
		//self.SendMsg("updatasp", smsg)
	}
	GetServer().SqlLog(self.Sql_UserBase.Uid, 91000004, num, param1, param2, dec, self.Sql_UserBase.SkillPoint, 0, self)
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "技能点", self.Sql_UserBase.SkillPoint) // _,_,来源，类型，剩余数量
	} else if num < 0 {
		GetServer().sendLog_UseMoney(self, num, dec, "技能点", self.Sql_UserBase.SkillPoint) // _,_,来源，类型，剩余数量
	}
}

// 计算技能点回复
func (self *Player) AutoSkillPoint() {
	if self.Sql_UserBase.SpLastUpdataTime <= 0 {
		return
	}

	// 经过时间
	stime := int(TimeServer().Unix() - self.Sql_UserBase.SpLastUpdataTime)

	// 回复
	add := stime / ADDSPTIME
	sub := stime % ADDSPTIME

	self.Sql_UserBase.SpLastUpdataTime = TimeServer().Unix() - int64(sub)

	if add <= 0 {
		return
	}

	max := GetCsvMgr().GetVipSkillNumber(self.Sql_UserBase.Vip)
	if self.Sql_UserBase.SkillPoint >= max {
		return
	}

	self.Sql_UserBase.SkillPoint = HF_MinInt(max, self.Sql_UserBase.SkillPoint+add)
}

func (self *Player) CheckSkillPointTime() {
	skillNumber := GetCsvMgr().GetVipSkillNumber(self.Sql_UserBase.Vip)
	if self.Sql_UserBase.SpLastUpdataTime > 0 && self.Sql_UserBase.SkillPoint >= skillNumber {
		self.Sql_UserBase.SpLastUpdataTime = 0
	} else if self.Sql_UserBase.SpLastUpdataTime <= 0 && self.Sql_UserBase.SkillPoint < skillNumber {
		self.Sql_UserBase.SpLastUpdataTime = TimeServer().Unix()
	}
}

func (self *Player) AddVipExp(num int, param1, param2 int, desc string) {
	if num <= 0 {
		return
	}

	if self.Sql_UserBase.Vip >= LEVELVIPMAX {
		return
	}

	oldlevel := self.Sql_UserBase.Vip
	self.Sql_UserBase.VipExp += num

	for {
		if self.Sql_UserBase.Vip >= LEVELVIPMAX {
			self.Sql_UserBase.Vip = LEVELVIPMAX
			self.Sql_UserBase.VipExp = 0
			break
		}

		needExp := GetCsvMgr().GetVipNeedExp(self.Sql_UserBase.Vip + 1)
		if self.Sql_UserBase.VipExp < needExp {
			break
		}
		self.Sql_UserBase.Vip++
		self.Sql_UserBase.VipExp -= needExp
		self.HandleTask(TASK_TYPE_VIP_BUY, self.Sql_UserBase.Vip, 0, 0)
	}

	if self.Sql_UserBase.Vip != oldlevel {
		self.GetModule("moneytask").(*ModMoneyTask).CheckOpen()
		self.GetModule("reward").(*ModReward).VipLevelChange(oldlevel)
		privilegeValues := self.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValues()
		self.GetModule("find").(*ModFind).CalSelfFind(privilegeValues, true)
		//GetServer().SqlLog(self.Sql_UserBase.Uid, LOG_USER_VIP_UP, 0, 0, int(self.Sql_UserBase.Fight/100), "玩家VIP升级", self.Sql_UserBase.Vip, self.Sql_UserBase.Level, self)
	}

	var msg S2C_UpdateExp
	msg.Cid = "addvipexp"
	msg.New = self.Sql_UserBase.Vip
	msg.Old = oldlevel
	msg.Newexp = self.Sql_UserBase.VipExp
	smsg, _ := json.Marshal(&msg)
	self.SendMsg("addvipexp", smsg)

	GetServer().SqlLog(self.Sql_UserBase.Uid, 91000022, num, param1, param2, desc, self.Sql_UserBase.Vip, self.Sql_UserBase.VipExp, self)
	GetServer().sendLog_GetMoney(self, num, desc, "贵族经验", self.Sql_UserBase.VipExp)
}

// 加经验
func (self *Player) AddExp(num int, param1, param2 int, dec string) {
	if num <= 0 {
		return
	}

	if self.Sql_UserBase.Level >= LEVELMAX {
		return
	}

	oldlevel := self.Sql_UserBase.Level
	tili := 0
	//计算升级奖励
	getItem := make(map[int]*Item)

	num = self.getExpNum(num)
	self.Sql_UserBase.Exp += num
	for {
		if self.Sql_UserBase.Level >= LEVELMAX {
			self.Sql_UserBase.Level = LEVELMAX
			self.Sql_UserBase.Exp = 0
			break
		}

		csv_teamexp := GetCsvMgr().GetTeamExp(self.Sql_UserBase.Level)
		if csv_teamexp == nil {
			self.Sql_UserBase.Exp = 0
			break
		}
		need := csv_teamexp.Teamexplv
		if self.Sql_UserBase.Exp >= need {
			self.Sql_UserBase.Exp -= need
			self.Sql_UserBase.Level += 1
			tili += csv_teamexp.Getphysical
			//计算升级奖励
			AddItemMapHelper(getItem, csv_teamexp.Items, csv_teamexp.Nums)
		} else {
			break
		}
	}

	if self.Sql_UserBase.Level != oldlevel {
		self.AddPower(tili, self.Sql_UserBase.Level, oldlevel, "用户升级")

		self.GetModule("union").(*ModUnion).UpdateUnionInfo()
		self.GetModule("moneytask").(*ModMoneyTask).CheckOpen()
		GetTopLevelMgr().SyncLevel(self.Sql_UserBase.Level, self)
		GetServer().sendLog_LevelupOk(self)

		getItems := self.AddObjectItemMap(getItem, "玩家升级", self.Sql_UserBase.Exp, 0, 0)

		var msg S2C_UpdateExp
		msg.Cid = "zdsj"
		msg.New = self.Sql_UserBase.Level
		msg.Old = oldlevel
		msg.Newexp = self.Sql_UserBase.Exp
		msg.GetItems = getItems
		smsg, _ := json.Marshal(&msg)
		self.SendMsg("zdsj", smsg)

		self.HandleTask(TASK_TYPE_PLAYER_LEVEL, 0, 0, 0)

		self.GetModule("recharge").(*ModRecharge).CheckOpen()
		self.GetModule("recharge").(*ModRecharge).CheckOpenLimit()
		self.NoticeCenterBaseInfo()

		//升级后战斗力变化
		self.countTeamFight(0)

		GetArenaSpecialMgr().Relevel(self.GetUid(), self.Sql_UserBase.Level)

		GetServer().SqlLog(self.Sql_UserBase.Uid, LOG_PLAYER_UP_LEVEL, self.Sql_UserBase.Level, oldlevel, 0, "用户升级", 0, 0, self)

		GetServer().SqlLog(self.Sql_UserBase.Uid, DEFAULT_EXP, num, self.Sql_UserBase.Level, oldlevel, "用户升级", self.Sql_UserBase.Exp, 0, self)

		GetServer().SendLog_SDKUP_AIWAN_LOGIN(self, SKDUP_ADDR_URL_AIWAN_SDK_EVENT_LEVELUP)
	}

	GetServer().SqlLog(self.Sql_UserBase.Uid, DEFAULT_EXP, num, param1, param2, dec, self.Sql_UserBase.Exp, 0, self)
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "军团经验", self.Sql_UserBase.Exp) // _,_,来源，类型，剩余数量
	} else if num < 0 {
		GetServer().sendLog_UseMoney(self, num, dec, "军团经验", self.Sql_UserBase.Exp) // _,_,来源，类型，剩余数量
	}

	// 检查玩家等级

}

func (self *Player) LvToMax() {
	if self.Sql_UserBase.Level >= LEVELMAX {
		return
	}
	oldlevel := self.Sql_UserBase.Level
	self.Sql_UserBase.Level = LEVELMAX
	self.Sql_UserBase.Exp = 0

	self.HandleTask(TASK_TYPE_PLAYER_LEVEL, 0, 0, 0)
	self.GetModule("union").(*ModUnion).UpdateUnionInfo()
	self.GetModule("moneytask").(*ModMoneyTask).CheckOpen()
	GetTopLevelMgr().SyncLevel(self.Sql_UserBase.Level, self)
	GetServer().sendLog_LevelupOk(self)

	var msg S2C_UpdateExp
	msg.Cid = "zdsj"
	msg.New = self.Sql_UserBase.Level
	msg.Old = oldlevel
	msg.Newexp = self.Sql_UserBase.Exp
	smsg, _ := json.Marshal(&msg)
	self.SendMsg("zdsj", smsg)

	self.GetModule("recharge").(*ModRecharge).CheckWarOrder()
	self.GetModule("recharge").(*ModRecharge).CheckWarOrderLimit()
	//升级后战斗力变化
	self.countTeamFight(0)
	self.NoticeCenterBaseInfo()
}

// 世界等级加速功能, 增加openLevel
func (self *Player) getExpNum(num int) int {
	day := GetServer().GetOpenTime()
	if day < 30 {
		return num
	}

	worldLevel := GetServer().Level
	if self.Sql_UserBase.Level >= worldLevel {
		return num
	}

	flag := GetCsvMgr().IsLevelOpen2(self, 55)
	if !flag {
		return num
	}

	diffLv := HF_AbsInt(worldLevel - self.Sql_UserBase.Level)
	factor := GetCsvMgr().getExpFactor(diffLv)
	return int(factor * float32(num))
}

// 加一个道具,必须是itemconfig里有的,num可负
//  普通货币：金币    高级货币：钻石
func (self *Player) AddObject(id, num int, param1, param2, param3 int, dec string) (int, int) {
	if num == 0 {
		return id, num
	}

	itemConfig := GetCsvMgr().GetItemConfig(id)
	if itemConfig == nil {
		tempId := (id/100)*100 + 1
		tempConfig := GetCsvMgr().GetItemConfig(tempId)
		if tempConfig == nil {
			//非物品类，不走配置
			switch id {
			case 1: // 能量召唤
				self.AddLootEnergy(num)
			}
			return id, num
		}
	}

	//看物品是否超过限制
	if itemConfig != nil && itemConfig.Overflow == LOGIC_TRUE {
		nowNum := self.GetObjectNum(id)
		if nowNum+num > itemConfig.MaxNum {
			//发邮件
			realNum := nowNum + num - itemConfig.MaxNum
			//..................
			num = num - realNum

			mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_GET_ITEM]
			if itemConfig.ItemType == ITEM_TYPE_MONEY {
				mailConfig, ok = GetCsvMgr().MailConfig[MAIL_ID_GET_ITEM_MONEY]
			}
			if !ok {
				return id, num
			}

			pMail := self.GetModule("mail").(*ModMail)
			if pMail == nil {
				return id, num
			}

			// 获得奖励
			var mailItems []PassItem
			mailItems = append(mailItems, PassItem{id, realNum})
			// 发送邮件
			if itemConfig.ItemType == ITEM_TYPE_MONEY {
				pMail.AddMailWithItems(MAIL_CAN_ALL_GET, mailConfig.Mailtitle, mailConfig.Mailtxt, mailItems)
				var msg S2C_ItemMax
				msg.Cid = "itemmax"
				msg.Type = 1
				self.SendMsg("itemmax", HF_JtoB(&msg))
			} else {
				pMail.AddMailWithItems(MAIL_CAN_ALL_GET, mailConfig.Mailtitle, mailConfig.Mailtxt, mailItems)
				var msg S2C_ItemMax
				msg.Cid = "itemmax"
				msg.Type = 2
				self.SendMsg("itemmax", HF_JtoB(&msg))
			}
			if num == 0 {
				return id, num
			}
		}
	}

	if num > 0 {
		self.HandleTask(TASK_TYPE_ITEM_GET_COUNT, num, id, 0)
	}

	itemid := id
	itemtype := 0
	isHelpMake := false
	if itemConfig == nil {
		tempId := (id/100)*100 + 1
		itemConfig = GetCsvMgr().GetItemConfig(tempId)
		itemtype = itemConfig.ItemType
		isHelpMake = true
	} else {
		itemtype = itemConfig.ItemType
	}

	if itemid == LIVENESS_DAILY_POINT { // 日活跃度
		self.GetModule("task").(*ModTask).CheckLivenessinfo()
		info := self.GetModule("task").(*ModTask).GetLivenessinfo(TASK_TYPE_DAILY)
		if info != nil {
			info.Liveness += num
			self.GetModule("union").(*ModUnion).AddUnionActivity(num)
			self.HandleTask(TASK_TYPE_ADD_LIVENESS, num, 0, 0)
			GetServer().SqlLog(self.Sql_UserBase.Uid, id, num, param1, param2, dec, info.Liveness, 0, self)
			if num > 0 {
				GetServer().sendLog_GetMoney(self, num, dec, "活跃度", info.Liveness) // _,_,来源，类型，剩余数量
			} else if num < 0 {
				GetServer().sendLog_UseMoney(self, num, dec, "活跃度", info.Liveness) // _,_,来源，类型，剩余数量
			}
			return id, num
		}
		return id, 0
	}

	if itemid == LIVENESS_WEEK_POINT { // 周活跃度
		self.GetModule("task").(*ModTask).CheckLivenessinfo()
		info := self.GetModule("task").(*ModTask).GetLivenessinfo(TASK_TYPE_WEEK)
		if info != nil {
			info.Liveness += num
			GetServer().SqlLog(self.Sql_UserBase.Uid, id, num, param1, param2, dec, info.Liveness, 0, self)
			if num > 0 {
				GetServer().sendLog_GetMoney(self, num, dec, "活跃度", info.Liveness) // _,_,来源，类型，剩余数量
			} else if num < 0 {
				GetServer().sendLog_UseMoney(self, num, dec, "活跃度", info.Liveness) // _,_,来源，类型，剩余数量
			}
			return id, num
		}
		return id, 0
	}

	if itemtype == ITEM_TYPE_MONEY {
		switch itemid {
		case 91000001: // 加金币
			self.AddGold(num, param1, param2, dec)
			return id, num
		case 91000002: // 钻石
			self.AddGem(num, param1, param2, param3, dec)
			return id, num
		case 91000003: // 体力
			self.AddPower(num, param1, param2, dec)
			return id, num
		case 91000004: // 技能点
			self.AddSkillPoint(num, param1, param2, dec)
			return id, num
		case 91000005: // 玩家主公经验
			self.AddExp(num, param1, param2, dec)
			return id, num
		case 91000017: // 战勋
			self.AddFeats(num, param1, param2, dec)
			return id, num
		case 91000018:
			self.AddHonor(num, param1, param2, dec)
			break
		case 91000022: // VIP经验
			self.AddVipExp(num, param1, param2, dec)
			return id, num
		case HeroSoul: // 魂石
			self.AddSoul(num, param1, param2, dec)
			return id, num
		case TechPoint: // 科技点
			self.AddTechPoint(num, param1, param2, dec)
			return id, num
		case BossMoney: // 巨兽精魄
			self.AddBossMoney(num, param1, param2, dec)
			return id, num
		case TowerStone: // 镇魂石
			self.AddTowerStone(num, param1, param2, dec)
			return id, num
		case WARORDER_ITEM_1:
			self.GetModule("recharge").(*ModRecharge).CalWarOrder(WARORDER_1, num)
			return id, num
		case WARORDER_ITEM_2:
			self.GetModule("recharge").(*ModRecharge).CalWarOrder(WARORDER_2, num)
			return id, num
		}
	}

	if itemtype == ITEM_TYPE_HERO {
		itemsubtype := itemConfig.ItemSubType
		switch itemsubtype {
		case 1: // 武将
			heroId := (itemid - 11000000) / 100
			param1 := itemid % 100
			if isHelpMake {
				self.GetModule("hero").(*ModHero).AddHero(heroId, num, param1, param2, dec)
			} else {
				self.GetModule("hero").(*ModHero).AddHero(heroId, num, itemConfig.Special, param2, dec)
			}
			return id, num
		}
	} else if itemtype == ITEM_TYPE_HORSE { // 战马
		for k := 0; k < num; k++ {
			(self.GetModule("horse").(*ModHorse)).AddHorseSafe(id, true, "道具获取")
		}

		csv_horse, ok := GetCsvMgr().Data["Horse_BattleSteed"][id]
		if ok {
			GetServer().SqlLog(self.GetUid(), id, num, HF_Atoi(csv_horse["quality"]), HF_Atoi(csv_horse["star"]), dec, param1, 1, self)
		}
		/*
			GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_MOD_ADD_OBJIECT_NOTICE"),
				self.Sql_UserBase.UName,
				dec, HF_GetColorByQuality(itemConfig.ItemCheck),itemConfig.ItemName, HF_GetColorByQuality(itemConfig.ItemCheck), num), 0, 1)
		*/
		return id, num
	} else if itemtype == ITEM_TYPE_HORSE_SOUL { // 马魂
		(self.GetModule("horse").(*ModHorse)).AddHorseSoulSafe(id, 1, num)

		csv_horse, ok := GetCsvMgr().HorseSoulConfig[id]
		if ok {
			GetServer().SqlLog(self.GetUid(), id, num, csv_horse.Quality, 0, dec, param1, 1, self)
		}
		/*
			GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_MOD_ADD_OBJIECT_NOTICE"),
				self.Sql_UserBase.UName,
				dec, HF_GetColorByQuality(itemConfig.ItemCheck),itemConfig.ItemName, HF_GetColorByQuality(itemConfig.ItemCheck), num), 0, 1)
		*/
		return id, num
	} else if itemtype == ITEM_TYPE_EQUIP { // 装备
		self.GetModule("equip").(*ModEquip).AddEquipWithParam(itemid, num, param1, param2, dec)
		self.HandleTask(EquipColorTask, 0, 0, 0)
		return id, num
	} else if itemtype == ITEM_TYPE_ARTIFACT { // 神器
		self.GetModule("artifactequip").(*ModArtifactEquip).AddArtifactWithParam(itemid, num, param1, param2, dec)
		self.HandleTask(EquipColorTask, 0, 0, 0)
		return id, num
	} else if itemtype == ITEM_TYPE_LOTTERY { // 箱子
		total := make(map[int]*Item)
		for index := 0; index < num; index++ {
			if itemConfig.LotteryId == 0 {
				continue
			}
			items := GetLootMgr().LootItem(itemConfig.LotteryId, self)
			AddItemMap(total, items)
		}
		outItems := self.AddObjectItemMap(total, "随机箱子", itemid, 0, 0)
		//self.GetModule("bag").(*ModBag).SendOnItem(outItems)
		for _, v := range outItems {
			return v.ItemID, v.Num
		}
		return id, num
	} else if itemtype == ITEM_TYPE_ICON && itemConfig.ItemSubType == 2 { // 箱子
		self.GetModule("head").(*ModHead).CheckUseItem(itemConfig.ItemId)
		return id, num
	} else if itemtype == ITEM_TYPE_PORTRAIT && itemConfig.ItemSubType == 2 { // 箱子
		self.GetModule("head").(*ModHead).CheckUseItem(itemConfig.ItemId)
		return id, num
	}

	self.CheckUseItem(itemid)
	// 加入背包
	(self.GetModule("bag").(*ModBag)).AddItem(id, num, param1, param2, dec)
	return id, num
}

func (self *Player) CheckItemLimit(id, num int) bool {
	itemConfig := GetCsvMgr().GetItemConfig(id)
	//看物品是否超过限制
	if itemConfig != nil && itemConfig.Overflow == LOGIC_TRUE {
		nowNum := self.GetObjectNum(id)
		if nowNum+num > itemConfig.MaxNum {
			return true
		}
	}

	return false
}

func (self *Player) GetObjectNum(id int) int {
	switch id {
	case 91000001: // 加金币
		return self.Sql_UserBase.Gold
	case 91000002: // 钻石
		return self.Sql_UserBase.Gem
	case 91000003: // 体力
		return self.GetPower()
	case 91000004: // 技能点
		return self.GetSkillPoint()
	case 91000005: // 军团经验
		return self.Sql_UserBase.Exp
	case 91000022:
		return self.Sql_UserBase.VipExp
	case LIVENESS_DAILY_POINT: // 活跃度
		self.GetModule("task").(*ModTask).CheckLivenessinfo()
		info := self.GetModule("task").(*ModTask).GetLivenessinfo(TASK_TYPE_DAILY)
		if nil != info {
			return info.Liveness
		} else {
			return 0
		}
	case LIVENESS_WEEK_POINT: // 活跃度
		self.GetModule("task").(*ModTask).CheckLivenessinfo()
		info := self.GetModule("task").(*ModTask).GetLivenessinfo(TASK_TYPE_WEEK)
		if nil != info {
			return info.Liveness
		} else {
			return 0
		}
	case HeroSoul: // 魂石
		return self.Sql_UserBase.Soul
	case BossMoney: // 巨兽精魄
		return self.Sql_UserBase.BossMoney
	case TowerStone: // 镇魂晶石
		return self.Sql_UserBase.TowerStone

	}

	return (self.GetModule("bag").(*ModBag)).GetItemNum(id)
}

// 设置剧情指引
// 为-1表示不变
func (self *Player) SetJqzy(jq, zy int) {
	self.SendRet4("jqzy", true)
}

// 设置辅助指引
// 为-1表示不变
func (self *Player) SetJqzy1(jq, zy int) {
	self.SendRet4("jqzy1", true)
}

// 取名
func (self *Player) Alertname(name string) {

	if self.Sql_UserBase.Gem < 500 && self.Sql_UserBase.IsRename == 0 {
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_NOT_ENOUGH_GEM_FAIL"))
		return
	}

	name = HF_FilterEmoji(name)

	if name == "" || !HF_IsLicitName([]byte(name)) || GetServer().IsSensitiveWord(name) {
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_CANT"))
		return
	}

	oldName := self.Sql_UserBase.UName
	ok := !HF_IsHasName(name)
	if ok {
		self.Sql_UserBase.UName = name
		GetUnionMgr().UpdateMember(self.GetModule("union").(*ModUnion).Sql_UserUnionInfo.Unionid, self.Sql_UserBase.Uid, 0, false, false)
		self.Sql_UserBase.Update(true)
		self.NoticeCenterBaseInfo()
	} else {
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_HASNAME"))
		return
	}
	self.GetModule("friend").(*ModFriend).Rename(true)
	GetTopPvpMgr().Rename(self)
	GetArenaMgr().Rename(self)
	GetArenaSpecialMgr().Rename(self)

	GetOfflineInfoMgr().Rename(self)
	GetHireHeroInfoMgr().Rename(self)
	GetActivityBossMgr().Rename(self)
	GetSupportHeroMgr().Rename(self)
	self.GetModule("entanglement").(*ModEntanglement).Rename()

	//同步实时排行榜
	GetTopMgr().SyncPlayerName(self)

	GetRankTaskMgr().Rename(self)

	GetPassRecordMgr().Rename(self)
	// 内政厅排行榜
	GetTopBuildMgr().Rename(self.Sql_UserBase.Uid, self.Sql_UserBase.UName)

	//同步副本排行名字 20190429 by zy
	GetUnionMgr().Rename(self.Sql_UserBase.Uid, self.Sql_UserBase.UName)

	if self.Sql_UserBase.IsRename == 0 {
		//! 保存老名字，备查
		self.AddGem(-500, 0, 0, 0, "玩家改名"+oldName)
	} else {
		self.Sql_UserBase.IsRename = 0
	}

	GetServer().SqlLog(self.Sql_UserBase.Uid, LOG_PLAYER_CHANGE_NAME, 0, 0, 0, "用户改名"+oldName+"->"+name, 0, 0, self)

	self.SendRet4("alertname", true)

	// 同步玩家信息
	self.SendInfo("updateuserinfo")
}

// 改icon
func (self *Player) Alerticon(icon int) {
	self.Sql_UserBase.IconId = icon
	GetOfflineInfoMgr().ReIconId(self)

	self.SendRet4("alerticon", true)
}

func (self *Player) CreateRole(name string, icon int, face int) {
	name = HF_FilterEmoji(name)

	if name == "" || !HF_IsLicitName([]byte(name)) || GetServer().IsSensitiveWord(name) {
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_CANT"))
		return
	}

	ok := !HF_IsHasName(name)
	if !ok {
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_HASNAME"))
		return
	}

	if self.Sql_UserBase.NameOk == 0 {
		self.Sql_UserBase.UName = name
		self.Sql_UserBase.IconId = icon
		self.Sql_UserBase.NameOk = 1
		self.NoticeCenterBaseInfo()
		GetMasterMgr().UpdatePlayer(self)
		self.NoticeCenter()
		GetServer().SendLog_SDKUP_AIWAN_LOGIN(self, SKDUP_ADDR_URL_AIWAN_SDK_EVENT_ENTERSVR)
	}

	self.Sql_UserBase.Face = face

	self.Sql_UserBase.Update(true)
	self.AddHead2(icon)

	self.SendRet2("createrole")

	// 发送创建角色log
	GetServer().sendLog_createOK(self)

	//阵营初始化 20190923 by zy
	self.Sql_UserBase.Camp = 1
	self.Sql_UserBase.CampOk = 1
	self.Sql_UserBase.City = HF_GetMainCityID(self.Sql_UserBase.Camp)

	GetServer().SqlLog(self.Sql_UserBase.Uid, LOG_USER_CREATE_PLAYER, 1, 0, 0, "创建角色", 0, 0, self)

}

// 处理任务
func (self *Player) HandleTask(tasktype, param1, param2, param3 int) {
	self.TaskLock.Lock()
	defer self.TaskLock.Unlock()

	//LogDebug("触发任务条件：", tasktype, param1, param2, param3)
	// 处理日常和主线任务
	self.GetModule("task").(*ModTask).HandleTask(tasktype, param1, param2, param3)
	// 处理试炼任务
	self.GetModule("targettask").(*ModTargetTask).HandleTask(tasktype, param1, param2, param3)
	// 处理福利任务
	self.GetModule("weekplan").(*ModWeekPlan).HandleTask(tasktype, param1, param2, param3)
	// 活动
	self.GetModule("activity").(*ModActivity).HandleTask(tasktype, param1, param2, param3)
	// 幸运商店
	self.GetModule("luckshop").(*ModLuckShop).HandleTask(tasktype, param1, param2, param3)
	self.GetModule("luckshop").(*ModLuckShop).StarHandleTask(tasktype, param1, param2, param3)
	self.GetModule("luckshop").(*ModLuckShop).DiscountHandleTask(tasktype, param1, param2, param3)
	self.GetModule("luckshop").(*ModLuckShop).StarLimitHandleTask(tasktype, param1, param2, param3)
	// 王国任务第二版
	self.GetModule("moneytask").(*ModMoneyTask).HandleTask(tasktype, param1, param2, param3)
	//// 神兽任务
	//self.GetModule("hydra").(*ModHydra).HandleTask(tasktype, param1, param2, param3)
	// 关卡主线
	self.GetModule("pass").(*ModPass).HandleTask(tasktype, param1, param2, param3)
	//! 限时礼包
	self.GetModule("timegift").(*ModTimeGift).HandleTask(tasktype, param1, param2, param3)
	// 排行任务
	self.GetModule("ranktask").(*ModRankTask).HandleTask(tasktype, param1, param2, param3)
	// 活动礼包
	self.GetModule("activitygift").(*ModActivityGift).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("growthgift").(*ModGrowthGift).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("specialpurchase").(*ModSpecialPurchase).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("nobilitytask").(*ModNobilityTask).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("viprecharge").(*ModVipRecharge).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("accesscard").(*ModAccessCard).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("recharge").(*ModRecharge).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("interstellar").(*ModInterStellar).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("activityboss").(*ModActivityBoss).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("honourshop").(*ModHonourShop).HandleTask(tasktype, param1, param2, param3)
}

func (self *Player) HandleTask2(tasktype, param1, param2, param3 int) {
	// 处理日常和主线任务
	self.GetModule("task").(*ModTask).HandleTask(tasktype, param1, param2, param3)
	// 处理试炼任务
	self.GetModule("targettask").(*ModTargetTask).HandleTask(tasktype, param1, param2, param3)
	// 处理福利任务
	self.GetModule("weekplan").(*ModWeekPlan).HandleTask(tasktype, param1, param2, param3)
	// 处理半月任务
	//self.GetModule("halfmoon").(*ModHalfMoon).HandleTask(tasktype, param1, param2, param3)
	// 活动
	self.GetModule("activity").(*ModActivity).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("interstellar").(*ModInterStellar).HandleTask(tasktype, param1, param2, param3)

	self.GetModule("activityboss").(*ModActivityBoss).HandleTask(tasktype, param1, param2, param3)
}

// 发送任务
func (self *Player) SendTask() {
	// 处理日常和主线任务
	self.GetModule("task").(*ModTask).SendUpdate()
	self.GetModule("targettask").(*ModTargetTask).SendUpdate()
	self.GetModule("interstellar").(*ModInterStellar).SendUpdate()
	// 处理福利任务
	self.GetModule("weekplan").(*ModWeekPlan).SendUpdate()
	// 幸运商店
	self.GetModule("luckshop").(*ModLuckShop).SendUpdate()
	//! 限时礼包
	self.GetModule("timegift").(*ModTimeGift).SendUpdate()
	self.GetModule("pass").(*ModPass).SendUpdate()
	self.GetModule("activitygift").(*ModActivityGift).SendUpdate()
	self.GetModule("growthgift").(*ModGrowthGift).SendUpdate()
	self.GetModule("nobilitytask").(*ModNobilityTask).SendUpdate()
	self.GetModule("accesscard").(*ModAccessCard).SendUpdate()
	self.GetModule("activityboss").(*ModActivityBoss).SendUpdate()
	self.GetModule("recharge").(*ModRecharge).SendUpdate()
	self.GetModule("honourshop").(*ModHonourShop).SendUpdate()
}

// 得到注册天数
func (self *Player) GetRegDays() int {
	t, _ := time.ParseInLocation(DATEFORMAT, self.Sql_UserBase.Regtime, time.Local)
	now := TimeServer()

	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	now1 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	subtime := now1.Unix() - t1.Unix()

	return int(subtime/86400) + 1
}

func (self *Player) GetCamp() int {
	return self.Sql_UserBase.Camp
}

func (self *Player) GetName() string {
	return self.Sql_UserBase.UName
}

func (self *Player) GetLv() int {
	return self.Sql_UserBase.Level
}

func (self *Player) GetUid() int64 {
	return self.Sql_UserBase.Uid
}

func (self *Player) GetUname() string {
	return self.Sql_UserBase.UName
}

func (self *Player) IsOnline() bool {
	return self.SessionObj != nil
}

// 是否是回流用户
func (self *Player) IsBack() bool {
	if self.Account == nil {
		return false
	}

	var account San_Account
	sql := fmt.Sprintf("select * from `san_accountaward` where `account` = '%s' and `creator` = '%s'", self.Account.Account, self.Account.Creator)
	GetServer().DBUser.GetOneData(sql, &account, "", 0)
	return account.Uid > 0
}

// 消息
func (self *Player) GetUserBaseInfo() Son_UserBaseInfo {
	var ret Son_UserBaseInfo
	ret.Checkinaward = self.Sql_UserBase.CheckinAward
	ret.Checkinnum = self.Sql_UserBase.CheckinNum
	ret.Exp = self.Sql_UserBase.Exp
	ret.Face = self.Sql_UserBase.Face
	ret.Gem = self.Sql_UserBase.Gem
	ret.Gold = self.Sql_UserBase.Gold
	ret.Iconid = fmt.Sprintf("%d", self.Sql_UserBase.IconId)
	ret.Ischeckin = self.IsCheckin()
	ret.Isrename = self.Sql_UserBase.IsRename
	ret.Lastcheckintime = self.Sql_UserBase.LastCheckinTime
	ret.Lastlivetime = self.Sql_UserBase.LastLiveTime
	ret.Lastlogintime = self.Sql_UserBase.LastLoginTime
	ret.Level = self.Sql_UserBase.Level
	ret.Levelaward = self.Sql_UserBase.LevelAward
	ret.Loginaward = self.Sql_UserBase.LoginAward
	ret.Logindays = self.Sql_UserBase.LoginDays
	ret.Morale = self.Sql_UserBase.Morale
	ret.Partyid = self.Sql_UserBase.PartyId
	ret.Position = self.Sql_UserBase.Position
	ret.Regtime = self.Sql_UserBase.Regtime
	ret.Skillpoint = self.GetSkillPoint()
	ret.Splastupdatatime = self.Sql_UserBase.SpLastUpdataTime
	ret.Tili = self.GetPower()
	ret.Tililastupdatatime = self.Sql_UserBase.TiLiLastUpdataTime
	ret.Uid = self.Sql_UserBase.Uid
	ret.Uname = self.Sql_UserBase.UName
	ret.Vip = self.Sql_UserBase.Vip
	ret.Vipexp = self.Sql_UserBase.VipExp
	ret.Worldaward = self.Sql_UserBase.WorldAward
	ret.Citylevel = self.Sql_UserBase.Citylevel
	ret.Camp = self.Sql_UserBase.Camp
	ret.City = self.Sql_UserBase.City
	ret.Day = GetServer().GetOpenTime()
	ret.Promotebox = self.Sql_UserBase.Promotebox
	ret.OpenServer = GetServer().GetOpenServer()
	ret.Channelid = self.Account.Channelid
	ret.Account = self.Account.Account
	ret.FitServer = self.GetFitServer()
	ret.Soul = self.Sql_UserBase.Soul
	ret.TechPoint = self.Sql_UserBase.TechPoint
	ret.BossMoney = self.Sql_UserBase.BossMoney
	ret.TowerStone = self.Sql_UserBase.TowerStone
	ret.Portrait = self.Sql_UserBase.Portrait
	ret.CampOk = self.Sql_UserBase.CampOk
	ret.NameOk = self.Sql_UserBase.NameOk
	ret.GuildId = self.Sql_UserBase.GuildId
	ret.RedIcon = self.Sql_UserBase.RedIcon
	ret.UserSignature = self.Sql_UserBase.UserSignature

	return ret
}

func (self *Player) SendInfo(cid string) {
	var msg S2C_UserBaseInfo
	msg.Cid = cid
	msg.Baseinfo = self.GetUserBaseInfo()
	self.SendMsg(cid, HF_JtoB(&msg))

}

func (self *Player) GetPlatfromInfo(cid string) {
	var msg S2C_PlatFromInfo
	msg.Cid = cid
	self.SendMsg(cid, HF_JtoB(&msg))
}

// 发送一个结果
func (self *Player) SendRet(cid string, ret int) {
	var msg S2C_ResultMsg
	msg.Cid = cid
	msg.Ret = ret
	smsg, _ := json.Marshal(&msg)
	self.SendMsg(cid, smsg)
}

// 发送一个结果2
func (self *Player) SendRet2(cid string) {
	var msg S2C_Result2Msg
	msg.Cid = cid
	smsg, _ := json.Marshal(&msg)
	self.SendMsg(cid, smsg)
}

// 发送一个结果3
func (self *Player) SendRet3(cid string, ret bool) {
	var msg S2C_Result3Msg
	msg.Cid = cid
	msg.Ret = ret
	smsg, _ := json.Marshal(&msg)
	self.SendMsg(cid, smsg)
}

// 发送一个结果4
func (self *Player) SendRet4(cid string, ok bool) {
	var msg S2C_Result4Msg
	msg.Cid = cid
	msg.Ok = ok
	smsg, _ := json.Marshal(&msg)
	self.SendMsg(cid, smsg)
}

// 发送错误
func (self *Player) SendErrInfo(cid string, info string) {
	if info == "错误" {
		LogDebug("出现错误，调用位置")
		_, file, line, ok := runtime.Caller(1)
		if ok {
			LogDebug(file, line)
		}

		_, file, line, ok = runtime.Caller(2)
		if ok {
			LogDebug(file, line)
		}
	}
	var msg S2C_ErrInfo
	msg.Cid = cid
	msg.Info = info
	smsg, _ := json.Marshal(&msg)
	self.SendMsg(cid, smsg)
	//DumpStacks()
}

// 玩家平台检查
func (self *Player) getPlatNo() int {
	s := self.Platform.Platform
	switch s {
	case "windows":
		return PLATFORM_WINDOWS
	case "android":
		return PLATFORM_ANDRIOD
	case "ios":
		return PLATFORM_IOS
	case "wp":
		return PLATFORM_WP
	case "mac":
		return PLATFORM_MAC
	default:
		return PLATFORM_DEFAULT
	}
}

func (self *Player) setPlayrInfo(DevInfo *SendRZ_EnvInfo_DevInfo_Ios, ChInfo *SendRZ_EnvInfo_ChInfo_Ios, AccountId string, AccountAppid string) {

	self.Platform.Platform = strings.ToLower(DevInfo.Os)
	self.Platform.DeviceId = DevInfo.DeviceId
	self.Platform.Brand = DevInfo.Brand
	self.Platform.Model = DevInfo.Model
	self.Platform.UUID = DevInfo.UUID
	self.Platform.Fr = DevInfo.Fr
	self.Platform.Res = DevInfo.Res
	self.Platform.Net = DevInfo.Net
	self.Platform.Mac = DevInfo.Mac
	self.Platform.Operator = DevInfo.Operator
	self.Platform.Ip = DevInfo.Ip
	self.Platform.Ch = ChInfo.Ch
	self.Platform.SubCh = ChInfo.SubCh
	self.Platform.AccountId = AccountId
	self.Platform.Account_AppId = AccountAppid
}

func (self *Player) GetAppleId() string {
	return self.Platform.Account_AppId
}

// 获取钻石个数
func (self *Player) GetGem() int {
	return self.Sql_UserBase.Gem
}

// 获取合服状态, 如果没有合服活动，则发送0
// 否则发送实际值
func (self *Player) GetFitServer() int {
	activityInfo := GetActivityMgr().GetActivity(FIT_SERVER_ACT_TYPE)
	if activityInfo == nil {
		return 0
	}

	if activityInfo.status.Status == 0 {
		return 0
	}

	if self.Sql_UserBase.FitServer == 0 { // 检查初始化状态
		self.Sql_UserBase.FitServer = 1
	}

	return self.Sql_UserBase.FitServer
}

/// 兼容合服后的Uid
func (self *Player) GetServerId() int {
	if self.Sql_UserBase.Uid < 10000000 {
		return GetServer().Con.ServerId
	}

	return int(self.Sql_UserBase.Uid / 10000000)
}

// 物品扣除和增加单个
func (self *Player) AddObjectSimple(id int, n int, reason string, param1, param2, param3 int) []PassItem {
	var items []PassItem

	itemConfig := GetCsvMgr().GetItemConfig(id)
	if itemConfig != nil && itemConfig.ItemType == ITEM_TYPE_LOTTERY {
		for i := 0; i < n; i++ {
			itemId, itemNum := self.AddObject(id, 1, param1, param2, param3, reason)
			items = append(items, PassItem{ItemID: itemId, Num: itemNum})
		}
	} else {
		itemId, itemNum := self.AddObject(id, n, param1, param2, param3, reason)
		items = append(items, PassItem{ItemID: itemId, Num: itemNum})
	}

	return items
}

// 物品扣除和增加单个
func (self *Player) RemoveObjectSimple(id int, n int, reason string, param1, param2, param3 int) []PassItem {
	var items []PassItem
	itemId, itemNum := self.AddObject(id, -n, param1, param2, param3, reason)
	if itemNum != 0 {
		items = append(items, PassItem{ItemID: itemId, Num: itemNum})
	}

	return items
}

// 增加多个道具
func (self *Player) AddObjectLst(ids []int, nums []int, reason string, param1, param2, param3 int) []PassItem {
	var res []PassItem
	if len(ids) != len(nums) {
		LogError("len(ids) != len(nums)", reason)
		return res
	}

	for index := range ids {
		id := ids[index]
		if id == 0 {
			continue
		}
		num := nums[index]
		if num == 0 {
			continue
		}
		itemConfig := GetCsvMgr().GetItemConfig(id)
		if itemConfig != nil && itemConfig.ItemType == ITEM_TYPE_LOTTERY {
			for i := 0; i < num; i++ {
				itemId, itemNum := self.AddObject(id, 1, param1, param2, param3, reason)
				res = append(res, PassItem{ItemID: itemId, Num: itemNum})
			}
		} else {
			itemId, itemNum := self.AddObject(id, num, param1, param2, param3, reason)
			res = append(res, PassItem{ItemID: itemId, Num: itemNum})
		}
	}

	return res
}

// 扣除多个道具
func (self *Player) RemoveObjectLst(ids []int, nums []int, reason string, param1, param2, param3 int) []PassItem {
	var res []PassItem
	if len(ids) != len(nums) {
		LogError("len(ids) != len(nums)", reason)
		return res
	}

	for index := range ids {
		id := ids[index]
		if id == 0 {
			continue
		}
		num := nums[index]
		if num == 0 {
			continue
		}
		itemId, itemNum := self.AddObject(id, -num, param1, param2, param3, reason)
		res = append(res, PassItem{ItemID: itemId, Num: itemNum})
	}

	return res
}

func (self *Player) RemoveObjectEasy(id int, num int, reason string, param1, param2, param3 int) []PassItem {
	var res []PassItem
	if id == 0 {
		return res
	}

	if num == 0 {
		return res
	}

	itemId, itemNum := self.AddObject(id, -num, param1, param2, param3, reason)
	res = append(res, PassItem{ItemID: itemId, Num: itemNum})

	return res
}

func (self *Player) AddObjectItemMap(itemMap map[int]*Item, reason string, param1, param2, param3 int) []PassItem {
	var res []PassItem
	for _, item := range itemMap {
		if item == nil {
			continue
		}
		id := item.ItemId
		num := item.ItemNum
		if id == 0 || num == 0 {
			continue
		}
		itemConfig := GetCsvMgr().GetItemConfig(id)
		if itemConfig != nil && itemConfig.ItemType == ITEM_TYPE_LOTTERY {
			for i := 0; i < num; i++ {
				itemId, itemNum := self.AddObject(id, 1, param1, param2, param3, reason)
				res = append(res, PassItem{ItemID: itemId, Num: itemNum})
			}
		} else {
			itemId, itemNum := self.AddObject(id, num, param1, param2, param3, reason)
			res = append(res, PassItem{ItemID: itemId, Num: itemNum})
		}
	}

	return res
}

func (self *Player) AddObjectPassItem(passItem []PassItem, reason string, param1, param2, param3 int) []PassItem {
	var res []PassItem
	for _, v := range passItem {
		id := v.ItemID
		num := v.Num
		if id == 0 || num == 0 {
			continue
		}
		itemId, itemNum := self.AddObject(id, num, param1, param2, param3, reason)
		res = append(res, PassItem{ItemID: itemId, Num: itemNum})
	}
	return res
}

func (self *Player) RemoveObjectItemMap(itemMap map[int]*Item, reason string, param1, param2, param3 int) []PassItem {
	var res []PassItem
	for _, item := range itemMap {
		if item == nil {
			continue
		}
		id := item.ItemId
		num := item.ItemNum
		if id == 0 || num == 0 {
			continue
		}
		itemId, itemNum := self.AddObject(id, -num, param1, param2, param3, reason)
		res = append(res, PassItem{ItemID: itemId, Num: itemNum})
	}

	return res
}

// 获得英雄
func (self *Player) getHero(heroKeyId int) *Hero {
	hero := self.GetModule("hero").(*ModHero).GetHero(heroKeyId)
	return hero
}

func (self *Player) getHeroes() map[int]*Hero {
	heroes := self.GetModule("hero").(*ModHero).GetHeroes()
	return heroes
}

func (self *Player) getEquips() [EQUIP_PACK_NUM]map[int]*Equip {
	mod := self.GetModule("equip").(*ModEquip)
	return mod.Data.equipItems
}

// 检查道具是否充足
func (self *Player) HasObjectOk(ids []int, nums []int) error {
	if len(ids) != len(nums) {
		return errors.New("len(ids) != len(nums)")
	}

	for index := range ids {
		costId := ids[index]
		costNum := nums[index]
		if costId == 0 {
			continue
		}
		if costNum < 0 {
			return errors.New("cost error")
		}

		if self.GetObjectNum(costId) < costNum {
			text := fmt.Sprintf(GetCsvMgr().GetText("STR_SHOP_NOT_ENOUGH_NUM"), GetCsvMgr().GetItemName(costId))
			return errors.New(text)
		}
	}

	return nil
}

func (self *Player) HasObjectMapItemOk(item map[int]*Item) error {

	for _, v := range item {
		costId := v.ItemId
		costNum := v.ItemNum
		if costId == 0 {
			continue
		}
		if costNum < 0 {
			return errors.New("cost error")
		}

		if self.GetObjectNum(costId) < costNum {
			text := fmt.Sprintf(GetCsvMgr().GetText("STR_SHOP_NOT_ENOUGH_NUM"), GetCsvMgr().GetItemName(costId))
			return errors.New(text)
		}
	}

	return nil
}

// 检查道具是否充足
func (self *Player) HasObjectOkEasy(costId int, costNum int) error {
	if costNum == 0 {
		return nil
	}

	if costId == 0 || costNum < 0 {
		return errors.New("cost error")
	}

	if self.GetObjectNum(costId) < costNum {
		text := fmt.Sprintf(GetCsvMgr().GetText("STR_SHOP_NOT_ENOUGH_NUM"), GetCsvMgr().GetItemName(costId))
		return errors.New(text)
	}

	return nil
}

func (self *Player) SendErr(info string) {
	var msg S2C_ErrInfo
	msg.Cid = "err"
	msg.Info = info
	smsg, _ := json.Marshal(&msg)
	self.SendMsg(msg.Cid, smsg)
}

// 检查道具是否充足
func (self *Player) hasItemMapOk(itemMap map[int]*Item) error {
	for _, item := range itemMap {
		costId := item.ItemId
		costNum := item.ItemNum
		if costId == 0 {
			continue
		}

		if costNum == 0 {
			continue
		}
		if costNum < 0 {
			return errors.New("cost error")
		}
		if self.GetObjectNum(costId) < costNum {
			text := fmt.Sprintf(GetCsvMgr().GetText("STR_SHOP_NOT_ENOUGH_NUM"), GetCsvMgr().GetItemName(costId))
			return errors.New(text)
		}
	}

	return nil
}

func (self *Player) getTeam() map[int]bool {
	//team := self.GetModule("team").(*ModTeam).getTeam()
	//return team
	return nil
}

func (self *Player) getFirstTeam() []int {
	team := self.GetModule("team").(*ModTeam).getFirstTeam()
	return team
}

func (self *Player) getTeamPos() *TeamPos {
	team := self.GetModule("team").(*ModTeam).getTeamPos(TEAMTYPE_DEFAULT)
	return team
}

func (self *Player) getTeamPosByType(teamType int) *TeamPos {
	team := self.GetModule("team").(*ModTeam).getTeamPos(teamType)
	return team
}

func (self *Player) IsLevelPass(levelid int) bool {
	return true
}

func (self *Player) GetVisitPassLevelId(id int, idx int) int {
	return 200000 + id*100 + idx + 1
}

func (self *Player) GetVisitPassNum(id int) int {
	num := 0
	for _, v := range GetCsvMgr().LevelConfigMap {
		if v.MainType == 2 && v.LevelType == 1 && v.LevelIndex == id {
			num += 1
		}
	}

	return num
}

// 加魂石
func (self *Player) AddSoul(num int, param1, param2 int, dec string) {
	if num == 0 {
		return
	}

	self.Sql_UserBase.Soul = HF_MinInt(HF_MaxInt(self.Sql_UserBase.Soul+num, 0), 2100000000)
	GetServer().SqlLog(self.Sql_UserBase.Uid, HeroSoul, num, param1, param2, dec, self.Sql_UserBase.Soul, 0, self)
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "魂石", self.Sql_UserBase.Soul)
	} else if num < 0 {
		GetServer().sendLog_UseMoney(self, num, dec, "魂石", self.Sql_UserBase.Soul)
		self.HandleTask(CostEnergeTask, -num, 0, HeroSoul)
	}

	// 刷新本地魂石
	//self.SendRet("updatesoul", self.Sql_UserBase.Soul)
}

// 加科技点
func (self *Player) AddTechPoint(num int, param1, param2 int, dec string) {
	if num == 0 {
		return
	}

	self.Sql_UserBase.TechPoint = HF_MinInt(HF_MaxInt(self.Sql_UserBase.TechPoint+num, 0), 2100000000)
	GetServer().SqlLog(self.Sql_UserBase.Uid, TechPoint, num, param1, param2, dec, self.Sql_UserBase.TechPoint, 0, self)
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "科技点", self.Sql_UserBase.TechPoint)
	} else if num < 0 {
		GetServer().sendLog_UseMoney(self, num, dec, "科技点", self.Sql_UserBase.TechPoint)
		self.HandleTask(CostEnergeTask, -num, 0, TechPoint)
	}

	// 刷新本地科技点
	self.SendRet("addtechpoint", self.Sql_UserBase.TechPoint)
}

// 加镇魂石
func (self *Player) AddTowerStone(num int, param1, param2 int, dec string) {
	if num == 0 {
		return
	}

	self.Sql_UserBase.TowerStone = HF_MinInt(HF_MaxInt(self.Sql_UserBase.TowerStone+num, 0), 2100000000)
	GetServer().SqlLog(self.Sql_UserBase.Uid, TowerStone, num, param1, param2, dec, self.Sql_UserBase.TowerStone, 0, self)
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "镇魂石", self.Sql_UserBase.TowerStone)
	} else if num < 0 {
		GetServer().sendLog_UseMoney(self, num, dec, "镇魂石", self.Sql_UserBase.TowerStone)
		self.HandleTask(CostEnergeTask, -num, 0, TowerStone)
	}

	// 刷新本地科技点
	self.SendRet("addtowerstone", self.Sql_UserBase.TowerStone)
}

// 巨兽水晶
func (self *Player) AddBossMoney(num int, param1, param2 int, dec string) {
	if num == 0 {
		return
	}

	self.Sql_UserBase.BossMoney = HF_MinInt(HF_MaxInt(self.Sql_UserBase.BossMoney+num, 0), 2100000000)
	GetServer().SqlLog(self.Sql_UserBase.Uid, BossMoney, num, param1, param2, dec, self.Sql_UserBase.BossMoney, 0, self)
	if num > 0 {
		GetServer().sendLog_GetMoney(self, num, dec, "巨兽精魄", self.Sql_UserBase.BossMoney)
	} else if num < 0 {
		GetServer().sendLog_UseMoney(self, num, dec, "巨兽精魄", self.Sql_UserBase.BossMoney)
	}

	// 刷新本地科技点
	self.SendRet("addbossmoney", self.Sql_UserBase.BossMoney)
}

// 同步战斗力
func (self *Player) synFight(heroKeyId int, fight int64, reason int, heroLv int) {
	cid := "synfight"
	msg := &S2C_SynFight{
		Cid:       cid,
		HeroKeyId: heroKeyId,
		Fight:     fight,
		Reason:    reason,
		HeroLv:    heroLv,
	}
	self.SendMsg(cid, HF_JtoB(msg))
}

// 同步整个战队战斗力
func (self *Player) synAllFight(heroIds []int, fights []int64, reason int, bossId int, bossFight int64) {
	if len(heroIds) <= 0 && bossId == 0 {
		return
	}
	cid := "synallfight"
	msg := &S2C_SynAllFight{
		Cid:       cid,
		HeroId:    heroIds,
		Fight:     fights,
		Reason:    reason,
		BossId:    bossId,
		BossFight: bossFight,
	}
	self.SendMsg(cid, HF_JtoB(msg))
}

func (self *Player) GetVip() int {
	return self.Sql_UserBase.Vip
}

func (self *Player) Send(head string, data interface{}) bool {
	body, err := json.Marshal(data)
	if err != nil {
		LogError(err.Error())
		return false
	}

	if self.Nosend {
		return true
	}
	session := self.GetSession()
	if session == nil {
		return false
	}
	session.SendMsg(head, body)
	return true
}

func (self *Player) GetUnionName() string {
	modUnion := self.GetModule("union").(*ModUnion)
	if modUnion != nil {
		return GetUnionMgr().GetUnionName(modUnion.Sql_UserUnionInfo.Unionid)
	}
	return ""
}

func (self *Player) GetUnionIcon() int {
	modUnion := self.GetModule("union").(*ModUnion)
	if modUnion != nil {
		return GetUnionMgr().GetUnionIcon(modUnion.Sql_UserUnionInfo.Unionid)
	}
	return 0
}

func (self *Player) GetUnionId() int {
	modUnion := self.GetModule("union").(*ModUnion)
	if modUnion != nil {
		return modUnion.Sql_UserUnionInfo.Unionid
	}
	return 0
}

func (self *Player) GetUnionLv() int {
	modUnion := self.GetModule("union").(*ModUnion)
	if modUnion != nil {
		return GetUnionMgr().GetUnionLv(modUnion.Sql_UserUnionInfo.Unionid)
	}
	return 0
}

func (self *Player) MaxTowerLv() int {
	pTower := self.GetModule("tower").(*ModTower)
	if pTower == nil {
		return 0
	}
	return pTower.MaxLevel()
}

func (self *Player) SetFight(fight int64, reason int) {
	//! 战斗力异常
	//if fight > 2100000000 {
	//	fight = 0
	//}

	if fight < 0 {
		fight = 0
	}

	self.Sql_UserBase.Fight = fight
	self.NoticeCenterBaseInfo()

	if self.GetUnionId() > 0 {
		self.GetModule("union").(*ModUnion).UpdateFight()
	}
	//LogDebug("fight changes, reason:", reason)
}

func (self *Player) CheckHero(heroId int) {
	self.GetModule("head").(*ModHead).CheckHero(heroId)
}

func (self *Player) CheckUseItem(itemId int) {
	self.GetModule("head").(*ModHead).CheckHero(itemId)
}

func (self *Player) AddHead2(id int) {
	modHead := self.GetModule("head").(*ModHead)
	if modHead != nil {
		modHead.AddHead2(id)
	}
}

func (self *Player) InitHead(id int) {
	modHead := self.GetModule("head").(*ModHead)
	if modHead != nil {
		modHead.InitHead(id)
	}
}

func (self *Player) getHeroIndex(heroId int) int {
	teamPos := self.getTeamPos()
	index := -1
	for i := 0; i < len(teamPos.FightPos); i++ {
		if heroId == teamPos.FightPos[i] {
			index = i
			break
		}
	}

	return index
}

func (self *Player) getHeroId(index int) int {
	teamPos := self.getTeamPos()
	for i := 0; i < len(teamPos.FightPos); i++ {
		if i == index {
			return teamPos.FightPos[i]
		}
	}

	return 0
}

func (self *Player) VipToMax() {
	if self.Sql_UserBase.Vip >= LEVELVIPMAX {
		return
	}

	self.Sql_UserBase.Vip = LEVELVIPMAX
	self.Sql_UserBase.VipExp = 0
	self.GetModule("recharge").(*ModRecharge).SendInfo("recharge")
}

func (self *Player) CheckTask() {
	//人物等级  		OK
	self.HandleTask(TASK_TYPE_PLAYER_LEVEL, 0, 0, 0)
	//试炼之塔  		OK
	self.GetModule("tower").(*ModTower).CheckTask()
	//爵位  			OK
	self.HandleTask(TASK_TYPE_NOBILITY_LEVEL, self.GetModule("nobilitytask").(*ModNobilityTask).Sql_NobilityTask.Level, 0, 0)
	//总战力  法阵等级	法阵上阵英雄数量  OK
	self.GetModule("crystal").(*ModResonanceCrystal).CheckTask()
	//装备品质  		!
	//装备强化 			!
	//英雄等级			OK
	self.HandleTask(TASK_TYPE_BIGGEST_LEVEL, 0, 0, 0)
	//英雄品质  		!
	//铸时星域完成度	OK
	self.GetModule("instance").(*ModInstance).CheckTask()
	//创世神木等级		OK
	self.GetModule("lifetree").(*ModLifeTree).CheckTask()
	//高阶竞技场段位
	self.GetModule("arenaspecial").(*ModArenaSpecial).CheckTask()
	//登录
	self.HandleTask(TASK_TYPE_IS_LOGIN, 1, 0, 0)

	modcrystal := self.GetModule("crystal").(*ModResonanceCrystal)
	if modcrystal != nil {
		self.HandleTask(TASK_TYPE_RESONANCE_CRYSTAL_COUNT, modcrystal.San_ResonanceCrystal.ResonanceCount, 0, 0)
		modcrystal.CheckLevel()
	}
	arenaData := GetArenaMgr().GetPlayerArenaData(self.GetUid())
	if arenaData != nil {
		self.HandleTask(TASK_TYPE_ARENA_POINT, int(arenaData.Point), 0, 0)
	}

	self.HandleTask(TASK_TYPE_VIP_BUY, 0, 0, 0)
}

func (self *Player) GetRegStampTime() int64 {
	regTime, _ := Parse(self.Sql_UserBase.Regtime)
	return regTime.Unix()
}

func (self *Player) CheckFightMsg(body []byte) {
	if self.CheckFight == nil {
		self.CheckFight = new(CheckFightInfo)
		self.CheckFight.Info = make(map[int][]float64)
	}

	var msg C2S_CheckFight
	json.Unmarshal(body, &msg)

	if msg.CheckId == 0 || len(msg.Info) == 0 {
		self.SendCheckFightErr(LOGIC_TRUE, &msg)
		return
	}

	if self.CheckFight.CheckId != msg.CheckId {

		fightInfo := GetRobotMgr().GetPlayerFightInfoByPos(self, 0, 0, msg.TeamType)
		if fightInfo == nil {
			log.Println(fmt.Sprintf("CheckFightMsg 错误,msg.teamtype:%d", msg.TeamType))
			return
		}
		self.CheckFight.Info = make(map[int][]float64)
		self.CheckFight.CheckId = msg.CheckId
		for i := 0; i < len(fightInfo.Heroinfo); i++ {
			self.CheckFight.Info[fightInfo.Heroinfo[i].HeroKeyId] = fightInfo.HeroParam[i].Param
		}
		return
	}

	//开始验证
	for keyId, param := range self.CheckFight.Info {
		_, ok := msg.Info[keyId]
		if !ok {
			self.SendCheckFightErr(2, &msg)
			return
		}
		if len(msg.Info[keyId]) < CHECKFIGHT_ATTR_END {
			self.SendCheckFightErr(3, &msg)
			return
		}
		if msg.Info[keyId][CHECKFIGHT_ATTR_HP-1] >= param[AttrHp]*1.5 {
			self.SendCheckFightErr(4, &msg)
			return
		}
		if msg.Info[keyId][CHECKFIGHT_ATTR_ATTACK-1] >= param[AttrAttack]*1.5 {
			self.SendCheckFightErr(4, &msg)
			return
		}
		if msg.Info[keyId][CHECKFIGHT_ATTR_DEFENCE-1] >= param[AttrDefence]*1.5 {
			self.SendCheckFightErr(4, &msg)
			return
		}
	}

	self.SendCheckFightErr(LOGIC_FALSE, &msg)
}

func (self *Player) SendCheckFightErr(code int, msg *C2S_CheckFight) {
	var msgRel S2C_CheckFight
	msgRel.Cid = "checkfight"
	msgRel.Code = code
	self.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

	if code != 0 {
		var errorInfo SQL_ErrorLog
		errorInfo.Uid = self.Sql_UserBase.Uid
		errorInfo.Stack = HF_JtoA(self.CheckFight)
		errorInfo.ErrorInfo = HF_JtoA(msg)
		errorInfo.Param1 = HF_JtoA(code)
		errorInfo.Time = TimeServer().Unix()
		InsertLogTable("san_errorinfo", &errorInfo, 1)
	}
}

func (self *Player) GetRankRewardRank(body []byte) {
	var msg C2S_GetRankRewardRank
	json.Unmarshal(body, &msg)

	GetRankRewardMgr().GetRank(self, msg.Id)
}

func (self *Player) GetRankRewardReward(body []byte) {
	var msg C2S_GetRankRewardReward
	json.Unmarshal(body, &msg)

	GetRankRewardMgr().GetReward(self, msg.Id)
}
