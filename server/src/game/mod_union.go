package game

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

const UNION_ACTIVITY_MAX_LIMIT = 100
const UNION_LIST_MAX = 10
const UNION_ACTIVITY_LIMIT_CLEAN = 8

const (
	UNION_ALERT_TYPE_NAME = 1
	UNION_ALERT_TYPE_ICON = 2
	UNION_ALERT_TYPE_BOTH = 3
)

type San_UserUnionInfo struct {
	Uid         int64  //! 用户ID
	Position    int    //! 用户等级，1-会长，2-副会长，4-会员
	Donation    int    //! 捐赠
	Givecount   int    //! 次数
	LastUpdTime int64  //! 上次退会时间
	Unionid     int    //! 军团id
	ApplyInfo   string //! 申请信息
	CopyNum     int    //! 今日副本次数
	CopyVer     int64  //! 副本版本
	CopyAward   string //! 副本奖励
	StateAward  int    // 是否领取周宝箱
	HuntLimit   string //!狩猎限制
	HuntStart   int    //!开启的狩猎

	applyInfo []int
	copyaward []int
	huntLimit []*UserHuntLimit //! 狩猎限制
	DataUpdate
}

type UserHuntLimit struct {
	Type         int           `json:"type"`       // 类型
	JoinCount    int           `json:"joincount"`  // 参与次数
	SweepCount   int           `json:"sweepcount"` // 参与次数
	EndTime      int64         `json:"endtime"`    // 结束时间
	MaxDamage    int64         `json:"maxdamage"`  // 历史最大伤害
	AwardCount   int           `json:"awardcount"` // 获得奖励个数
	BattleInfo   *BattleInfo   `json:"battleinfo"`
	BattleRecord *BattleRecord `json:"battlerecord"`
	GemItems     []PassItem    `json:"gemitems"` // 钻石物品
}

type JS_UserUnionInfo struct {
	Uid              int64 `json:"uid"`
	Position         int   `json:"position"`
	Donation         int   `json:"donation"`
	Givecount        int   `json:"givecount"`
	Unionid          int   `json:"unionid"`
	ApplyInfo        []int `json:"applyinfo"`
	CopyNum          int   `json:"copynum"`
	StateAward       int   `json:"state_award"`
	ActivityToday    int   `json:"activitytoday"`    // 当日活跃度
	ActivitySevenday int   `json:"activitysevenday"` // 七日
	LastUpdTime      int64 `json:"lastupdtime"`      //! 上次退会时间
}

type JS_UnionHero struct {
	Heroid int `json:"heroid"`
	Star   int `json:"star"`
	Level  int `json:"level"`
	Color  int `json:"color"`
	Talent int `json:"talent"`
}

type ModUnion struct {
	player            *Player           //! 角色对象
	Sql_UserUnionInfo San_UserUnionInfo //! 数据库结构
	lastUnionCall     int64             //! 上一次招募时间
	Locker            *sync.RWMutex     //! 操作锁
}

const QUIT_CD = 3600
const JOIN_CD = 3600
const BRAVE_HAND_CD = 86400
const UNION_HUNTER_CD = 86400

func (self *ModUnion) OnGetData(player *Player) {
	self.player = player
	self.Locker = new(sync.RWMutex)
	sql := fmt.Sprintf("select * from `san_userunioninfo` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_UserUnionInfo, "san_userunioninfo", self.player.ID)

	if self.Sql_UserUnionInfo.Uid <= 0 {
		self.Sql_UserUnionInfo.Uid = self.player.ID
		self.Sql_UserUnionInfo.applyInfo = make([]int, 0)
		self.Encode()
		InsertTable("san_userunioninfo", &self.Sql_UserUnionInfo, 0, true)
	} else {
		self.Decode()
	}

	self.Sql_UserUnionInfo.Init("san_userunioninfo", &self.Sql_UserUnionInfo, true)
	self.lastUnionCall = 0
}

func (self *ModUnion) OnGetOtherData() {
	//unionid := self.Sql_UserUnionInfo.Unionid
	//if unionid > 0 {
	//	GetUnionMgr().UpdateMemberState(unionid, self.player.Sql_UserBase.Uid)
	//}
	self.Sql_UserUnionInfo.HuntStart = 0
}

func (self *ModUnion) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_UserUnionInfo.ApplyInfo), &self.Sql_UserUnionInfo.applyInfo)
	json.Unmarshal([]byte(self.Sql_UserUnionInfo.CopyAward), &self.Sql_UserUnionInfo.copyaward)
	//json.Unmarshal([]byte(self.Sql_UserUnionInfo.ActivityLimit), &self.Sql_UserUnionInfo.activityLimit)
	json.Unmarshal([]byte(self.Sql_UserUnionInfo.HuntLimit), &self.Sql_UserUnionInfo.huntLimit)
}

func (self *San_UserUnionInfo) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.ApplyInfo), &self.applyInfo)
	json.Unmarshal([]byte(self.CopyAward), &self.copyaward)
	//json.Unmarshal([]byte(self.ActivityLimit), &self.activityLimit)
	json.Unmarshal([]byte(self.HuntLimit), &self.huntLimit)
}

func (self *ModUnion) Encode() { //! 将data数据写入数据库
	self.Sql_UserUnionInfo.ApplyInfo = HF_JtoA(self.Sql_UserUnionInfo.applyInfo)
	self.Sql_UserUnionInfo.CopyAward = HF_JtoA(self.Sql_UserUnionInfo.copyaward)
	//self.Sql_UserUnionInfo.ActivityLimit = HF_JtoA(self.Sql_UserUnionInfo.activityLimit)
	self.Sql_UserUnionInfo.HuntLimit = HF_JtoA(self.Sql_UserUnionInfo.huntLimit)
}

func (self *ModUnion) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "unionalertname": //军团改名 --- 完成
		var c2s_msg C2S_AlertUnionName
		json.Unmarshal(body, &c2s_msg)
		self.AlertUnionName(c2s_msg.Unionid, c2s_msg.Newname, c2s_msg.IconId)
		return true
	case "alertnotice": //修改军团内部公告 --- 完成
		var c2s_msg C2S_AlertUnionNotice
		json.Unmarshal(body, &c2s_msg)
		self.AlertUnionNotice(c2s_msg.Unionid, c2s_msg.Notice)
		return true
	case "alertboard": //修改军团外部介绍 --- 完成
		var c2s_msg C2S_AlertUnionNotice
		json.Unmarshal(body, &c2s_msg)
		self.AlertUnionBoard(c2s_msg.Unionid, c2s_msg.Notice)
		return true
	case "alertset": // 修改军团设置 --- 完成
		var c2s_msg C2S_AlertUnionSet
		json.Unmarshal(body, &c2s_msg)
		self.AlertUnionSet(c2s_msg.Unionid, c2s_msg.Jointype, c2s_msg.Joinlevel)
		return true
	case "applyunion": //申请某军团 --- 完成
		var c2s_msg C2S_ApplyUnion
		json.Unmarshal(body, &c2s_msg)
		self.ApplyUnion(c2s_msg.Unionid)
		return true
	case "cancel_applyunion": //撤销申请某军团 --- 完成
		var c2s_msg C2S_Cancel_ApplyUnion
		json.Unmarshal(body, &c2s_msg)
		self.CancelApplyUnion(c2s_msg.Unionid, self.player.Sql_UserBase.Uid, true)
		return true
	case "createunion": //创建军团 --- 完成
		var c2s_msg C2S_CreateUnion
		json.Unmarshal(body, &c2s_msg)
		self.CreateUnion(c2s_msg.Name, c2s_msg.Icon)
		return true
	case "dissolveunion": // 解散军团 --- 完成
		var c2s_msg C2S_Dissolveunion
		json.Unmarshal(body, &c2s_msg)
		self.DissolveUnion(c2s_msg.Unionid)
		return true
	case "getunioninfo": // 查询工会信息  --- 完成
		var c2s_msg C2S_Getunioninfo
		json.Unmarshal(body, &c2s_msg)
		self.GetUnionInfo(c2s_msg.Unionid)
		return true
	case "getunionlist": // 获取工会列表 --- 完成
		self.GetUnionList()
		return true
	case "getuserunioninfo": // 获得玩家工会信息 --- 完成
		self.GetUserUnionInfo()
		return true
	case "getunionrecord": //! 获得军团日志 --- 完成
		var c2s_msg C2S_Getunioninfo
		json.Unmarshal(body, &c2s_msg)
		self.GetUnionRecord(c2s_msg.Unionid)
		return true
	case "joinunion": //加入军团  --- 完成
		var c2s_msg C2S_Joinunion
		json.Unmarshal(body, &c2s_msg)
		self.JoinUnion(c2s_msg.Unionid)
		return true
	case "masterfail": ////会长拒绝加入 --- 完成
		var c2s_msg C2S_Masterfail
		json.Unmarshal(body, &c2s_msg)
		self.MasterFail(c2s_msg.Unionid, c2s_msg.Applyuid)
		return true
	case "masterok": //会长同意加入  --- 完成
		var c2s_msg C2S_Masterok
		json.Unmarshal(body, &c2s_msg)
		self.MasterOk(c2s_msg.Unionid, c2s_msg.Applyuid)
		return true
	case "masterallok": //会长同意加入  --- 完成
		var c2s_msg C2S_MasterAllok
		json.Unmarshal(body, &c2s_msg)
		self.MasterAllOk(c2s_msg.Unionid)
		return true
	case "masteroutplayer": // 踢出玩家 --- 完成
		var c2s_msg C2S_Masteroutplayer
		json.Unmarshal(body, &c2s_msg)
		self.MasterOutPlayer(c2s_msg.Unionid, c2s_msg.Outuid)
		return true
	case "outunion": //离开军团 --- 完成
		var c2s_msg C2S_Outunion
		json.Unmarshal(body, &c2s_msg)
		self.OutUnion(c2s_msg.Unionid)
		return true
	case "unionmodify": //军团任命 --- 完成
		var c2s_msg C2S_UnionModify
		json.Unmarshal(body, &c2s_msg)
		self.UnionModify(c2s_msg.Unionid, c2s_msg.Destuid, c2s_msg.Op)
		return true
	case "setbravehand": //无畏之手 --- 完成
		var c2s_msg C2S_SetBraveHand
		json.Unmarshal(body, &c2s_msg)
		self.SetBraveHand(c2s_msg.Unionid, c2s_msg.Destuid, c2s_msg.Op)
		return true
	case "findunion": // 搜索军团  --- 完成
		var c2s_msg C2S_UnionFind
		json.Unmarshal(body, &c2s_msg)
		self.UnionFind(c2s_msg.Type, HF_Atoi(c2s_msg.Unionid), c2s_msg.Unionname)
		return true
	case "unioncall":
		self.SendCallInfo()
		return true
	case "unionsendmail":
		var msg C2S_UnionSendMail
		err := json.Unmarshal(body, &msg)
		if err != nil {
			self.player.SendErr(err.Error())
			return true
		}
		self.SendMail(msg.Title, msg.Text)
		return true
	case "start_hunt_fight": // 开始战斗 --- 完成
		var msg C2S_StartHuntFight
		json.Unmarshal(body, &msg)
		self.StartHuntFight(msg.Type)
		return true
	case "open_hunt_fight": // 会长开启活动 --- 完成
		var msg C2S_OpenHuntFight
		json.Unmarshal(body, &msg)
		self.OpenHuntFight(msg.Type)
		return true
	case "end_hunt_fight": // 结束战斗 --- 完成
		var msg C2S_EndHuntFight
		json.Unmarshal(body, &msg)
		self.EndHuntFight(msg.Type, msg.Damage, msg.BattleInfo)
		return true
	case "sweep_hunt_fight": // 扫荡 --- 完成
		var msg C2S_SweepHuntFight
		json.Unmarshal(body, &msg)
		self.SweepHuntFight(msg.Type)
		return true
	case "get_hunt_info":
		var msg C2S_GetHuntInfo
		json.Unmarshal(body, &msg)
		self.GetHuntInfo()
		return true
	case "get_hunt_dps_top":
		// 获得层数战报
		var msg C2S_HuntDpsTop
		json.Unmarshal(body, &msg)
		self.GetDpsTop(msg.Type)
		return true
	}

	return false
}

func (self *ModUnion) GetResidualTime() string {
	residualtime := QUIT_CD - TimeServer().Unix() + self.Sql_UserUnionInfo.LastUpdTime
	msg := &S2CUnionCDTime{}
	msg.Cid = "send_union_cd_time"
	msg.CDTime = residualtime
	self.player.Send(msg.Cid, msg)
	return fmt.Sprintf(GetCsvMgr().GetText("STR_UNION_TIME_LIMIT"), residualtime/3600, (residualtime%3600)/60, residualtime%60)
}

func (self *ModUnion) OnSave(sql bool) {
	self.Encode()
	self.Sql_UserUnionInfo.Update(sql)
}

func (self *ModUnion) OnRefresh() {
	self.Sql_UserUnionInfo.Givecount = 0
	self.Sql_UserUnionInfo.CopyNum = 0
	self.Sql_UserUnionInfo.StateAward = 0
}

//军团改名 ret：0成功 1只有会长才能操作 2钻石不足 3:已有军团名 4
func (self *ModUnion) AlertUnionName(unionid int, newname string, icon int) {
	union := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if union == nil {
		return
	}

	hasAt := strings.Contains(union.Unionname, "@")
	ret := 0

	if self.Sql_UserUnionInfo.Position > UNION_POSITION_VICE_MASTER {
		ret = 1
		self.player.SendRet("union_alertname", ret)
		return
	}

	cfg := GetCsvMgr().GetTariffConfig2(TariffChangeUnionName)
	if cfg == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_TOWER_RESET_ITEM"))
		return
	}

	oldname := union.Unionname
	oldicon := union.Icon

	if newname != union.Unionname {
		if err := self.player.HasObjectOk(cfg.ItemIds, cfg.ItemNums); err != nil {
			if !hasAt {
				ret = 2
				self.player.SendRet("union_alertname", ret)
				return
			} else {
				ret = 4
				self.player.SendRet("union_alertname", ret)
				return
			}
		}
	}

	success, alerttype := GetUnionMgr().AlertUnionName(self.player.GetUid(), unionid, newname, icon)
	if !success {
		ret = 3
		self.player.SendRet("union_alertname", ret)
		return
	}

	var item []PassItem
	if alerttype != UNION_ALERT_TYPE_ICON {
		if err := self.player.HasObjectOk(cfg.ItemIds, cfg.ItemNums); err != nil {
			if !hasAt {
				ret = 2
				self.player.SendRet("union_alertname", ret)
				return
			} else {
				ret = 4
				self.player.SendRet("union_alertname", ret)
				return
			}
		}

		if !hasAt {
			//self.player.AddGem(-500, 2, unionid, "军团改名")
			item = self.player.RemoveObjectLst(cfg.ItemIds, cfg.ItemNums, "军团改名", union.Id, 0, 0)
		} else {
			ret = 4
			self.player.SendRet("union_alertname", ret)
			return
		}
	}

	if alerttype == UNION_ALERT_TYPE_NAME {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_ALERT_NAME, unionid, union.Level, 0, "公会改名:"+oldname+"->"+newname, 0, 0, self.player)
	} else if alerttype == UNION_ALERT_TYPE_ICON {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_CHANGE_ICON, unionid, oldicon, icon, "公会改旗帜", 0, union.Level, self.player)
	} else if alerttype == UNION_ALERT_TYPE_BOTH {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_ALERT_NAME, unionid, union.Level, 0, "公会改名:"+oldname+"->"+newname, 0, 0, self.player)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_CHANGE_ICON, unionid, oldicon, icon, "公会改旗帜", 0, union.Level, self.player)
	}

	//! 修改排行榜军团名字
	//GetTopUnionMgr().ChangeUnionName(unionid, newname)
	var msg S2C_AlertUnionName
	msg.Cid = "union_alertname"
	msg.Ret = ret
	msg.Money = item
	self.player.SendMsg("union_alertname", HF_JtoB(msg))
}

//修改军团公告 ret：1只有会长才能操作
func (self *ModUnion) AlertUnionNotice(unionid int, notice string) {
	union := GetUnionMgr().GetUnion(unionid)
	if union == nil {
		return
	}

	ret := 0
	if self.Sql_UserUnionInfo.Position > UNION_POSITION_VICE_MASTER {
		ret = 1
	} else {
		GetUnionMgr().AlertUnionNotice(self.player.GetUid(), unionid, notice)
	}

	self.player.SendRet("union_alertnotice", ret)
}

//修改军团公告 ret：1只有会长才能操作
func (self *ModUnion) AlertUnionBoard(unionid int, board string) {
	union := GetUnionMgr().GetUnion(unionid)
	if union == nil {
		return
	}

	ret := 0
	if self.Sql_UserUnionInfo.Position > UNION_POSITION_VICE_MASTER {
		ret = 1
	} else {
		GetUnionMgr().AlertUnionBoard(self.player.GetUid(), unionid, board)

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_CHANGE_BOARD, unionid, 0, 0, "工会改宣言", 0, 0, self.player)
	}

	self.player.SendRet("union_alertboard", ret)
}

//修改军团设置
func (self *ModUnion) AlertUnionSet(unionid int, jointype int, joinlevel int) {
	union := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if union == nil {
		return
	}

	ret := 0
	if self.Sql_UserUnionInfo.Position > UNION_POSITION_VICE_MASTER {
		ret = 1
	} else {
		GetUnionMgr().AlertUnionSet(self.player.GetUid(), unionid, jointype, joinlevel)
		self.player.SendErrInfo("err", GetCsvMgr().GetText("设置成功"))
	}

	self.player.SendRet("union_alertset", ret)
}

//申请某军团
func (self *ModUnion) ApplyUnion(unionid int) {
	union := GetUnionMgr().GetUnion(unionid)
	if union == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_UNION_NOT"))
		return
	}

	for i := 0; i < len(self.Sql_UserUnionInfo.applyInfo); i++ {
		if self.Sql_UserUnionInfo.applyInfo[i] == unionid {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_UNION_ANAGIN"))
			return
		}
	}

	if self.player.Sql_UserBase.Level < union.Joinlevel {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_UNION_APPLY"), union.Joinlevel))
		return
	}

	if TimeServer().Unix()-self.Sql_UserUnionInfo.LastUpdTime < JOIN_CD {
		self.GetResidualTime()
		return
	}

	var member JS_UnionApply

	member.Uid = self.player.Sql_UserBase.Uid
	member.Uname = self.player.Sql_UserBase.UName
	member.Level = self.player.Sql_UserBase.Level
	member.Vip = self.player.Sql_UserBase.Vip
	member.Fight = self.player.Sql_UserBase.Fight
	member.Iconid = self.player.Sql_UserBase.IconId
	member.Portrait = self.player.Sql_UserBase.Portrait
	member.Applytime = TimeServer().Unix()

	if !GetUnionMgr().AddApply(unionid, member) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_UNION_FAIL_APPLY"))
		return
	}

	LogDebug("申请加入军团", self.player.Sql_UserBase.Uid, unionid)
	self.Sql_UserUnionInfo.applyInfo = append(self.Sql_UserUnionInfo.applyInfo, unionid)

	self.player.SendRet("applyunion", 0)
}

//撤销申请某军团
func (self *ModUnion) CancelApplyUnion(unionid int, uid int64, send bool) {
	self.Locker.Lock()
	for j := 0; j < len(self.Sql_UserUnionInfo.applyInfo); j++ {
		if self.Sql_UserUnionInfo.applyInfo[j] == unionid {
			copy(self.Sql_UserUnionInfo.applyInfo[j:], self.Sql_UserUnionInfo.applyInfo[j+1:])
			self.Sql_UserUnionInfo.applyInfo = self.Sql_UserUnionInfo.applyInfo[:len(self.Sql_UserUnionInfo.applyInfo)-1]
			break
		}
	}
	self.Locker.Unlock()

	GetUnionMgr().CancelApply(unionid, uid)

	if send {
		self.player.SendRet("cancel_applyunion", 0)
	}

}

//! 清除军团的申请
func (self *ModUnion) ClearApplyUnion() {
	for j := 0; j < len(self.Sql_UserUnionInfo.applyInfo); j++ {
		GetUnionMgr().CancelApply(self.Sql_UserUnionInfo.applyInfo[j], self.player.Sql_UserBase.Uid)
	}
	self.Sql_UserUnionInfo.applyInfo = make([]int, 0)
}

//创建军团 ret:1已有军团，不能创建 2:钻石不足，需要500钻石 3:军团重名
func (self *ModUnion) CreateUnion(name string, icon int) {
	name = HF_FilterEmoji(name)

	if name == "" || !HF_IsLicitName([]byte(name)) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_CANT"))
		return
	}
	cfg := GetCsvMgr().GetTariffConfig2(TariffCreateUnion)
	if cfg == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_TOWER_RESET_ITEM"))
		return
	}

	var ret = 0
	if self.Sql_UserUnionInfo.Unionid != 0 {
		ret = 1
	} else {
		if err := self.player.HasObjectOk(cfg.ItemIds, cfg.ItemNums); err != nil {
			ret = 2
		} else {
			if GetUnionMgr().CheckName(name) {
				ret = 3
			} else {
				//self.player.AddGem(-500, 26, 0, "军团创建")
				// 创建 创建成功
				unionid := GetUnionMgr().CreateUnion(icon, name, self.player.Sql_UserBase.Uid, self.player.Sql_UserBase.UName, self.player.Sql_UserBase.Camp)
				if unionid > 0 {
					// 添加玩家数据
					if GetUnionMgr().UpdateMember(unionid, self.player.GetUid(), UNION_POSITION_MASTER, true, true) {
						self.player.RemoveObjectLst(cfg.ItemIds, cfg.ItemNums, "创建公会", unionid, 0, 0)
						self.Sql_UserUnionInfo.Unionid = unionid
						self.Sql_UserUnionInfo.Position = UNION_POSITION_MASTER
						self.Sql_UserUnionInfo.Donation = 0
						self.player.HandleTask(TASK_TYPE_JOIN_UNION, 0, 0, 0)
						self.Sql_UserUnionInfo.applyInfo = make([]int, 0)
						// 清理所有数据
						self.player.GetModule("union").(*ModUnion).ClearApplyUnion()
						self.Sql_UserUnionInfo.Update(true)

						GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_CREATE, unionid, 0, 0, "创建公会", 0, 0, self.player)

					} else {
						ret = 1
					}
				} else {
					ret = 3
				}
			}
		}
	}

	var msg S2C_CreateUnion
	msg.Cid = "createunion"
	msg.Ret = ret
	msg.Info = self.GetJsUserUnionInfo()

	if ret == 0 {
		msg.Unioninfo = *GetUnionMgr().GetUnionJsInfo(self.Sql_UserUnionInfo.Unionid)
	}

	msg.Money = make([]PassItem, 0)
	msg.Money = append(msg.Money, PassItem{91000002, -500})
	msg.CopyUpdateTime = msg.Unioninfo.CopyUpdate
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("createunion", smsg)

}

//解散军团 ret:1只有会长才能操作 2:军团人员大于1
func (self *ModUnion) DissolveUnion(unionid int) {
	union := GetUnionMgr().GetUnion(unionid)
	if union == nil {
		return
	}
	level := union.Level
	ret := GetUnionMgr().DissolveUnion(self.Sql_UserUnionInfo.Unionid, self.player.Sql_UserBase.Uid)
	if ret == 0 {
		self.Sql_UserUnionInfo.Unionid = 0
		self.Sql_UserUnionInfo.LastUpdTime = TimeServer().Unix()

		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_DISSOLVE, unionid, level, 0, "解散公会", 0, 0, self.player)

	}

	self.player.SendRet("dissolveunion", ret)
}

func (self *ModUnion) GetUnionInfo(unionid int) {
	GetUnionMgr().UpdateUnion(unionid)
	var msg S2C_GetUnionInfo
	msg.Cid = "getunioninfo"
	msg.Unioninfo = GetUnionMgr().GetUnionJsInfo(unionid)

	if msg.Unioninfo == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_UNION_NOT2"))
		//return
	}

	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("getunioninfo", smsg)
}

//获取军团列表
func (self *ModUnion) GetUnionList() {
	var msg S2C_GetUnionList
	msg.Cid = "getunionlist"
	msg.List = GetUnionMgr().GetUnionList(self.player.GetServerId())

	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("getunionlist", smsg)
}

//! 获得军团日志
func (self *ModUnion) GetUnionRecord(unionid int) {
	if self.Sql_UserUnionInfo.Unionid == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_YOU_DONT_HAVE_A_LEGION"))
		return
	}
	var msg S2C_GetUnionRecord
	msg.Cid = "getunionrecord"
	msg.Ret = 1
	lstRecord := GetUnionMgr().GetUnionRecord(self.Sql_UserUnionInfo.Unionid)
	for i := 0; i < len(lstRecord); i++ {
		msg.Record = append(msg.Record, lstRecord[i])
	}

	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("getunionrecord", smsg)
}

func (self *ModUnion) GetJsUserUnionInfo() JS_UserUnionInfo {
	var ret JS_UserUnionInfo

	ret.Uid = self.Sql_UserUnionInfo.Uid
	ret.Position = self.Sql_UserUnionInfo.Position
	ret.Donation = self.Sql_UserUnionInfo.Donation
	ret.Givecount = self.Sql_UserUnionInfo.Givecount
	ret.Unionid = self.Sql_UserUnionInfo.Unionid
	ret.ApplyInfo = self.Sql_UserUnionInfo.applyInfo
	ret.CopyNum = self.Sql_UserUnionInfo.CopyNum
	ret.StateAward = self.Sql_UserUnionInfo.StateAward
	ret.LastUpdTime = self.Sql_UserUnionInfo.LastUpdTime

	now := TimeServer()
	timeStamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
	if TimeServer().Hour() < 5 {
		timeStamp -= DAY_SECS
	}
	data := GetUnionMgr().GetUnion(ret.Unionid)
	if data != nil {
		temp, success := self.GetUnionMember(data, ret.Uid)
		if success {
			nCount := 0
			for _, v := range temp.ActivityRecord {
				nCount += v.AddCount

				if v.Time == timeStamp {
					ret.ActivityToday = v.AddCount
				}
			}
			ret.ActivitySevenday = nCount
		}
	}

	return ret
}

func (self *ModUnion) UpdateUnionInfo() {
	GetUnionMgr().UpdateMember(self.Sql_UserUnionInfo.Unionid, self.player.GetUid(), 0, false, false)
}

func (self *ModUnion) GetUserUnionInfo() {
	//! 军团长自动转移-有BUG，待测试
	if self.Sql_UserUnionInfo.Unionid > 0 {
		GetUnionMgr().CheckMasterOffline(self.Sql_UserUnionInfo.Unionid)
		GetUnionMgr().UpdateUnion(self.Sql_UserUnionInfo.Unionid)
		GetUnionMgr().UpdateMember(self.Sql_UserUnionInfo.Unionid, self.player.GetUid(), 0, false, false)
	}

	var msg S2C_GetUserUnionInfo
	msg.Cid = "getuserunioninfo"
	msg.Selfinfo = self.GetJsUserUnionInfo()
	msg.Unioninfo.Member = make([]JS_Member, 0)
	msg.Unioninfo.Apply = make([]JS_UnionApply, 0)

	if self.Sql_UserUnionInfo.Unionid == 0 {
		msg.Unionlist = GetUnionMgr().GetUnionList(self.player.GetServerId())
	} else {
		data := GetUnionMgr().GetUnionJsInfo(self.Sql_UserUnionInfo.Unionid)
		if data == nil {
			LogDebug("获得数据失败：", self.Sql_UserUnionInfo.Unionid)
			return
		}

		//LogDebug(data.Level)
		HF_DeepCopy(&msg.Unioninfo, data)
		//if self.Sql_UserUnionInfo.Position > UNION_POSITION_VICE_MASTER {
		//	msg.Unioninfo.Apply = make([]JS_UnionApply, 0)
		//}
		msg.CopyUpdateTime = data.CopyUpdate

		msg.ChangeMaster = false
		union := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
		if union != nil {
			if union.IsChangeMaster(self.player.Sql_UserBase.Uid) == true {
				union.ChangeMaster.CheckPlayer = append(union.ChangeMaster.CheckPlayer, self.player.Sql_UserBase.Uid)

				msg.ChangeMaster = true
				msg.OldMaster = union.ChangeMaster.OldMaster
			}

			member, ok := self.GetUnionMember(union, self.player.GetUid())
			if ok {
				if self.Sql_UserUnionInfo.Position != member.Position {
					self.Sql_UserUnionInfo.Position = member.Position
				}

				if member.Lastlogintime != 0 {
					GetUnionMgr().UpdateMemberState(union.Id, self.Sql_UserUnionInfo.Uid)
				}
			} else {
				var msg S2M_UnionCheckOut
				msg.Uid = self.player.GetUid()
				msg.Unionuid = self.Sql_UserUnionInfo.Unionid
				ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_CHECK_OUT, &msg)
				if ret == nil || ret.RetCode != UNION_SUCCESS {
					self.OutUnionData(0)
				}
			}
		}
	}

	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("getuserunioninfo", smsg)

	self.GetHuntInfo()
}

//军团检查是否满员
func (self *ModUnion) CheckFull(unioninfo *San_Union) bool {
	csv_community := GetCsvMgr().CommunityConfig[unioninfo.Level]

	if len(unioninfo.member) >= csv_community.Membernum {
		return true
	}
	return false
}

func (self *ModUnion) AddPlayer(unionid int) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.Sql_UserUnionInfo.Uid = self.player.Sql_UserBase.Uid
	self.Sql_UserUnionInfo.Unionid = unionid
	self.Sql_UserUnionInfo.Position = UNION_POSITION_MEMBER
	self.Sql_UserUnionInfo.Donation = 0
	self.Sql_UserUnionInfo.Givecount = 0
	self.Sql_UserUnionInfo.applyInfo = make([]int, 0)
	self.player.HandleTask(TASK_TYPE_JOIN_UNION, 0, 0, 0)
	self.Sql_UserUnionInfo.Update(true)
}

//加入军团 ret:1:无此军团 2:不满足军团加入等级 3:禁止加入 5:军团需要申请，不能直接加入 4军团已满
func (self *ModUnion) JoinUnion(unionid int) {
	union := GetUnionMgr().GetUnion(unionid)

	if union == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_THE_LEGION_DOES_NOT_EXIST"))
		return
	}

	if self.Sql_UserUnionInfo.Unionid != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_IF_YOU_HAVE_JOINED_THE"))
		return
	}

	if union == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	} else {
		if union.Jointype == 0 { //! 允许任何人加入
			if self.CheckFull(union) {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_THE_LEGION_IS_FULL"))
				return
			}

			if self.player.Sql_UserBase.Level < union.Joinlevel {
				self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_UNION_APPLY"), union.Joinlevel))
				return
			}

			if TimeServer().Unix()-self.Sql_UserUnionInfo.LastUpdTime < JOIN_CD {
				self.GetResidualTime()
				return
			}

			var msg S2M_UnionJoin
			msg.Unionuid = unionid
			msg.Uid = self.player.GetUid()
			ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_JOIN, msg)
			if ret == nil || ret.RetCode != UNION_SUCCESS {
				return
			}

			LogDebug("加入军团：", unionid, union.Unionname, self.player.Sql_UserBase.Uid)
			success, _ := GetUnionMgr().AddPlayer(unionid, self.player.Sql_UserBase.Uid, false)
			if success {
				self.player.GetModule("union").(*ModUnion).ClearApplyUnion()
				self.player.GetModule("union").(*ModUnion).AddPlayer(unionid)
			}

			//更新在线军团排行榜
			//GetTopMgr().SyncUnionFight(data.Fight, data)

		} else if union.Jointype == 1 { //! 禁止任何人加入
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_THE_LEGION_FORBIDS_ANYONE_TO"))
			return
		} else if union.Jointype == 2 { //! 需要申请
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_THE_LEGION_NEEDS_AN_APPLICATION"))
			return
		}
	}

	var msg S2C_JoinUnion
	msg.Cid = "joinunion"
	msg.Selfinfo = self.GetJsUserUnionInfo()
	msg.Info = *GetUnionMgr().GetUnionJsInfo(unionid)
	msg.Ret = 0
	self.player.SendMsg("joinunion", HF_JtoB(&msg))
}

////会长拒绝加入
func (self *ModUnion) MasterFail(unionid int, applyuid int64) {
	data := GetUnionMgr().GetUnion(unionid)
	if data == nil {
		self.player.SendRet("masterfail", -1)
		return
	}
	msgret := 0

	var msg S2M_UnionMasterFail
	msg.Uid = self.player.GetUid()
	msg.Unionuid = unionid
	msg.ApplyUid = applyuid
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_MASTER_FAIL, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		self.player.SendRet("masterfail", -1)
		return
	}

	GetUnionMgr().CleanApply(unionid, applyuid)

	self.player.SendRet("masterfail", msgret)

}

//会长同意加入 ret:1只有会长才能操作 2:该玩家已经加入别人军团
func (self *ModUnion) MasterOk(unionid int, applyuid int64) {
	data := GetUnionMgr().GetUnion(unionid)
	if data == nil {
		self.player.SendRet("masterok", -1)
		return
	}

	msgret := 0

	if self.Sql_UserUnionInfo.Position > UNION_POSITION_VICE_MASTER {
		msgret = 1
	} else {
		if self.CheckFull(data) == true {
			self.player.SendErrInfo("err", "军团已满，提升等级可以增加更多会员")
			msgret = 3
		} else {
			//GetUnionMgr().AddPlayer(unionid, applyuid, true)

			var msg S2M_UnionMasterOK
			msg.Uid = self.player.GetUid()
			msg.Unionuid = unionid
			msg.ApplyUid = applyuid
			ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_MASTER_OK, &msg)
			if ret == nil || ret.RetCode != UNION_SUCCESS {
				return
			}

			retFlag, ret1 := GetUnionMgr().AddPlayer(unionid, applyuid, true)
			if retFlag == false {
				msgret = ret1
			} else {
				//self.CancelApplyUnion(unionid, applyuid, false)
				GetUnionMgr().CancelApply(unionid, applyuid)
				//更新在线军团排行榜
				//GetTopMgr().SyncUnionFight(data.Fight, data)
				self.GetUserUnionInfo()
			}
			GetUnionMgr().CleanApply(unionid, msg.ApplyUid)
		}
	}

	self.player.SendRet("masterok", msgret)
}

//会长同意加入 ret:1只有会长才能操作 2:该玩家已经加入别人军团
func (self *ModUnion) MasterAllOk(unionid int) {
	data := GetUnionMgr().GetUnion(unionid)
	if data == nil {
		self.player.SendRet("masterallok", -1)
		return
	}

	if self.Sql_UserUnionInfo.Position > UNION_POSITION_VICE_MASTER {
		self.player.SendRet("masterallok", -1)
		return
	}

	if self.CheckFull(data) == true {
		self.player.SendErrInfo("err", "军团已满，提升等级可以增加更多会员")
		return
	}

	var msg S2M_UnionMasterOK
	msg.Uid = self.player.GetUid()
	msg.Unionuid = unionid
	msg.ApplyUid = 0
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_MASTER_OK, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}

	var backmsg M2S_UnionMasterOK
	json.Unmarshal([]byte(ret.Data), &backmsg)

	nLen := len(backmsg.ApplyUid)

	var addPlayer []int64
	for i := nLen - 1; i >= 0; i-- {
		addPlayer = append(addPlayer, backmsg.ApplyUid[i])
		retFlag, _ := GetUnionMgr().AddPlayer(unionid, backmsg.ApplyUid[i], true)
		if retFlag == false {
			addPlayer = addPlayer[:len(addPlayer)-1]
			continue
		} else {
			//self.CancelApplyUnion(unionid, applyuid, false)
			//data.CancelApply(data.apply[i].Uid)
			//更新在线军团排行榜
			//GetTopMgr().SyncUnionFight(data.Fight, data)
		}
	}

	GetUnionMgr().CleanApply(unionid, 0)

	if len(addPlayer) <= 0 {
		self.player.SendErrInfo("err", "未添加新成员")
		return
	}

	var retmsg S2C_MasterAllok
	retmsg.Cid = "masterallok"
	retmsg.Unionid = unionid

	for _, v := range addPlayer {
		temp, success := self.GetUnionMember(data, v)
		if success {
			retmsg.AddPlayer = append(retmsg.AddPlayer, &temp)
		}
	}

	smsg, _ := json.Marshal(&retmsg)
	self.player.SendMsg("masterallok", smsg)
}

// 踢除玩家
func (self *ModUnion) MasterOutPlayer(unionid int, outuid int64) {
	// 获得公会
	data := GetUnionMgr().GetUnion(unionid)
	if data == nil {
		return
	}
	// 不能踢会长
	if data.Masteruid == outuid {
		self.player.SendRet("masteroutplayer", 2)
		return
	}
	// 不能踢自己
	if self.player.GetUid() == outuid {
		self.player.SendRet("masteroutplayer", 1)
		return
	}
	// 获得操作者
	selfdata, ok := self.GetUnionMember(data, self.player.Sql_UserBase.Uid)
	if !ok {
		return
	}
	// 权限不足
	if selfdata.Position > UNION_POSITION_VICE_MASTER {
		self.player.SendRet("masteroutplayer", 1)
		return
	}
	// 获得被操作者
	destdata, ok2 := self.GetUnionMember(data, outuid)
	if !ok2 {
		return
	}
	// 不能操作比自己权限高或同级的人
	if selfdata.Position >= destdata.Position {
		self.player.SendRet("masteroutplayer", 1)
		return
	}

	// 中心服检测能否踢人
	var msg S2M_UnionKickPlayer
	msg.Uid = self.player.GetUid()
	msg.Unionuid = unionid
	msg.OutUid = outuid
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_KICK_PLAYER, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		self.player.SendRet("masteroutplayer", 1)
		return
	}
	// 执行退会操作
	GetUnionMgr().OutPlayer(unionid, outuid, true)

	self.player.SendRet("masteroutplayer", 0)

	//更新在线军团排行榜
	//GetTopMgr().SyncUnionFight(data.Fight, data)
}

//操作数据
func (self *ModUnion) OutUnionData(unionid int) {

	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.Sql_UserUnionInfo.Unionid = 0
	self.Sql_UserUnionInfo.LastUpdTime = TimeServer().Unix()
}

//离开军团 ret:1会长不能退出军团
func (self *ModUnion) OutUnion(unionid int) {
	// 获得公会
	data := GetUnionMgr().GetUnion(unionid)
	if data == nil {
		return
	}

	ret := 0
	// 会长不能退会
	if data.Masteruid == self.player.Sql_UserBase.Uid {
		ret = 1
	} else {
		// 执行退会操作
		GetUnionMgr().OutPlayer(unionid, self.player.Sql_UserBase.Uid, false)
		// 设置cd时间
		//self.Sql_UserUnionInfo.LastUpdTime = TimeServer().Unix()
	}

	self.player.SendRet("outunion", ret)

	//更新在线军团排行榜
	//GetTopMgr().SyncUnionFight(data.Fight, data)
}

func (self *ModUnion) GetUnionMember(data *San_Union, uid int64) (JS_Member, bool) {
	for i := 0; i < len(data.member); i++ {
		if data.member[i].Uid == uid {
			return data.member[i], true
		}
	}

	var ret JS_Member

	return ret, false
}

//军团任命 提升发负数，降级发正数
func (self *ModUnion) UnionModify(unionid int, destuid int64, op int) {
	// 公会不存在
	data := GetUnionMgr().GetUnion(unionid)
	if data == nil {
		return
	}
	// 权限越界
	if op < UNION_POSITION_MASTER || op > UNION_POSITION_MEMBER {
		return
	}
	// 获得操作者
	selfdata, ok := self.GetUnionMember(data, self.player.Sql_UserBase.Uid)
	if !ok {
		return
	}
	// 获取配置
	csv, _ := GetCsvMgr().CommunityConfig[data.Level]
	if csv == nil {
		return
	}
	// 获得被操作者
	destdata, ok2 := self.GetUnionMember(data, destuid)
	if !ok2 {
		return
	}
	// 会长以下没有权限
	if selfdata.Position > UNION_POSITION_MASTER {
		return
	}
	// 不能操作自己
	if destuid == self.player.Sql_UserBase.Uid {
		return
	}
	// 不能操作同级或比自己权限高的人
	if selfdata.Position >= destdata.Position {
		return
	}
	// 不能给别人设置比自己更高的权限
	if op < selfdata.Position {
		return
	}
	// 设置为会长 则是转让工会
	if op == UNION_POSITION_MASTER {
		// 我得是会长
		if self.player.Sql_UserBase.Uid != data.Masteruid {
			return
		}
		// 并且只能转让给副会长
		if destdata.Position != UNION_POSITION_VICE_MASTER {
			return
		}
		// 转让工会
		if !GetUnionMgr().UnionChange(unionid, self.player.Sql_UserBase.Uid, destuid) {
			return
		}
		// 会长自己成为普通成员
		self.ModifyPosition(UNION_POSITION_MEMBER)
	} else { // 设置副会长及降级为普通成员
		// 设置副会长
		if op == UNION_POSITION_VICE_MASTER {
			// 判断副会长个数
			party := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
			if party.GetMemberCount(UNION_POSITION_VICE_MASTER) >= csv.Elder {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_THE_LEGION_CAN_ONLY_HAVE"))
				return
			}
		}
		// 设置
		GetUnionMgr().UnionModify(self.player.GetUid(), unionid, destuid, op)
	}

	var msg S2C_UnionModify
	msg.Cid = "unionmodify"
	msg.Uid = self.player.Sql_UserBase.Uid
	msg.Destuid = destuid
	msg.Op = op
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("unionmodify", smsg)
}

// 搜索军团
func (self *ModUnion) UnionFind(_type int, unionid int, name string) {

	if _type == 1 {
		data := GetUnionMgr().GetUnion(unionid)

		ret := 0

		if data == nil {
			ret = 1
		} else if data.GetCamp() != self.player.Sql_UserBase.Camp {
			ret = 1
		}

		var msg S2C_FindUnion
		msg.Cid = "union_find"

		if ret == 0 {
			msg.Info = append(msg.Info, *GetUnionMgr().GetUnionJsInfo2(unionid))
		}

		msg.Ret = ret

		smsg, _ := json.Marshal(&msg)
		self.player.SendMsg("union_find", smsg)
	} else {
		data := GetUnionMgr().GetUnionByName(name)

		var msg S2C_FindUnion
		msg.Cid = "union_find"
		//camp := self.player.Sql_UserBase.Camp
		//for _, v := range data {
		//	if v == nil {
		//		continue
		//	}
		//	msg.Info = append(msg.Info, *GetUnionMgr().GetUnionJsInfo2(v.Id))
		//	ret = 0
		//}
		ret := 0
		if data == nil {
			msg.Info = []JS_Union2{}
			ret = 1
		} else {
			msg.Info = data
			ret = 0
		}
		msg.Ret = ret

		smsg, _ := json.Marshal(&msg)
		self.player.SendMsg("union_find", smsg)
	}
}

//撤销申请某军团
func (self *ModUnion) ModifyPosition(pos int) {
	self.Locker.Lock()
	self.Sql_UserUnionInfo.Position = pos
	self.Locker.Unlock()
}

//! 发布召集信息
func (self *ModUnion) SendCallInfo() {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_YOU_DONT_HAVE_A_LEGION"))
		return
	}

	if data.GetPlayerUnionLevel(self.player.Sql_UserBase.Uid) > 2 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_YOU_DONT_HAVE_PERMISSION_TO"))
		return
	}

	calltime := GetUnionMgr().GetUnionCallTime(self.Sql_UserUnionInfo.Unionid)
	if TimeServer().Unix() < calltime+300 {
		leavetime := 300 - (TimeServer().Unix() - calltime)
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_UNION_HIRE_ERROR"), leavetime))
		return
	}
	GetUnionMgr().SetUnionCallTime(self.Sql_UserUnionInfo.Unionid, TimeServer().Unix())
	var msg S2C_Chat
	msg.Cid = "unioncall"
	msg.Uid = self.player.Sql_UserBase.Uid
	msg.Channel = CHAT_CAMP_SYSTEM
	msg.Name = self.player.Sql_UserBase.UName
	msg.Icon = self.player.Sql_UserBase.IconId
	msg.Portrait = self.player.Sql_UserBase.Portrait
	msg.Camp = self.player.Sql_UserBase.Camp
	msg.Vip = self.player.Sql_UserBase.Vip
	msg.Level = self.player.Sql_UserBase.Level
	msg.Time = TimeServer().Unix()
	msg.Content = fmt.Sprintf(GetCsvMgr().GetText("STR_UNION_ATTEND"), data.Unionname)
	msg.Url = ""
	msg.Param = data.Id

	GetPlayerMgr().BroadCastMsgToCamp(self.player.Sql_UserBase.Camp, "unioncall", HF_JtoB(&msg))
	//GetSessionMgr().BroadCastMsg("unioncall", HF_JtoB(&msg))
}

func (self *San_UserUnionInfo) Encode() { //! 将data数据写入数据库
	self.ApplyInfo = HF_JtoA(self.applyInfo)
	self.CopyAward = HF_JtoA(self.copyaward)
	//self.ActivityLimit = HF_JtoA(self.activityLimit)
	self.HuntLimit = HF_JtoA(self.huntLimit)
}

func (self *ModUnion) GetStatisticsValue1080() int {
	return 3 - self.Sql_UserUnionInfo.CopyNum
}

// 调公会限制经验
func (self *ModUnion) GMAddUnionExpLimit(add_num int) {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	var msg S2M_UnionGMAdd
	msg.Unionuid = self.Sql_UserUnionInfo.Unionid
	msg.Type = UNION_GM_TYPE_EXPLIMIT
	msg.Count = add_num
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GM_ADD, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}

	data.DayExp += add_num
	return
}
func (self *ModUnion) GMAddUnionExp(add_num int) {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	var msg S2M_UnionGMAdd
	msg.Unionuid = self.Sql_UserUnionInfo.Unionid
	msg.Type = UNION_GM_TYPE_EXP
	msg.Count = add_num
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GM_ADD, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}

	data.Exp += add_num
	return
}
func (self *ModUnion) GMAddUnionLevel() {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	var msg S2M_UnionGMAdd
	msg.Unionuid = self.Sql_UserUnionInfo.Unionid
	msg.Type = UNION_GM_TYPE_LEVEL
	msg.Count = len(GetCsvMgr().CommunityConfig)
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GM_ADD, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}
	data.Level = len(GetCsvMgr().CommunityConfig)
	return
}
func (self *ModUnion) GMAddUnionActivity(add_num int) {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	var msg S2M_UnionGMAdd
	msg.Unionuid = self.Sql_UserUnionInfo.Unionid
	msg.Type = UNION_GM_TYPE_ACTIVITY
	msg.Count = add_num
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GM_ADD, msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}
	data.ActivityPoint += add_num
	return
}

// 增加公会活跃度
func (self *ModUnion) AddUnionActivity(add_num int) int {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		return 0
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()

	csv, _ := GetCsvMgr().CommunityConfig[data.Level]
	if csv == nil {
		return 0
	}

	nCount := add_num * csv.Changelively
	nCount = data.AddActivityPoint(self.player.GetUid(), nCount)

	addExp := add_num * csv.Changeexp
	data.AddExp(addExp, 1)

	GetUnionMgr().UpdateMember(data.Id, self.player.GetUid(), 0, false, false)

	return nCount
}
func (self *ModUnion) GetHunterGemItems(config *UnionHuntDropConfig) []PassItem {
	ret := []PassItem{}
	if len(config.DiamondChance) != len(config.Diamond) {
		return ret
	}

	for i := 0; i < config.DiamondTime; i++ {
		randNum := HF_GetRandom(10000)
		check := 0

		for t, v := range config.DiamondChance {
			itemId := config.Diamond[t]
			if itemId == 0 {
				continue
			}
			check += v
			if randNum < check {
				ret = append(ret, PassItem{itemId, 1})
				break
			}
		}
	}

	return ret
}

// 开始狩猎战斗
func (self *ModUnion) StartHuntFight(nType int) bool {
	config := GetCsvMgr().GetUnionHuntConfigByID(nType)
	if nil == config {
		return false
	}

	dropConfig := GetCsvMgr().GetUnionHuntDropConfig(config.Group, 0)
	if nil == dropConfig {
		return false
	}

	if TimeServer().Unix()-self.Sql_UserUnionInfo.LastUpdTime < UNION_HUNTER_CD && nType == UNION_HUNT_TYPE_SPECIAL {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("离开公会24小时内不能参与狩猎"))
		return false
	}

	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	csv_vip := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if csv_vip == nil || len(csv_vip.GuildHunting) < 2 {
		return false
	}

	var mastermsg S2M_UnionStartHunter
	mastermsg.Uid = self.player.GetUid()
	mastermsg.Unionuid = self.Sql_UserUnionInfo.Unionid
	mastermsg.Type = nType
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_START_HUNTER, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_UnionStartHunter
	json.Unmarshal([]byte(ret.Data), &backmsg)

	var huntlimit *UserHuntLimit = nil
	for _, v := range self.Sql_UserUnionInfo.huntLimit {
		if v.Type == nType {
			huntlimit = v
			break
		}
	}

	if nil == huntlimit {
		huntlimit = &UserHuntLimit{nType, 0, 0, backmsg.EndTime, 0, 0, nil, nil, nil}
		self.Sql_UserUnionInfo.huntLimit = append(self.Sql_UserUnionInfo.huntLimit, huntlimit)
	} else if huntlimit.EndTime != backmsg.EndTime {
		huntlimit.JoinCount = 0
		huntlimit.SweepCount = 0
		huntlimit.GemItems = nil
		huntlimit.EndTime = backmsg.EndTime
	}

	maxCount := 0
	if nType == UNION_HUNT_TYPE_NOMAL {
		maxCount = csv_vip.GuildHunting[0]
	} else {
		maxCount = csv_vip.GuildHunting[1]
	}
	if huntlimit.SweepCount+huntlimit.JoinCount >= maxCount {
		return false
	}

	self.Sql_UserUnionInfo.HuntStart = nType

	var msg S2C_StartHuntFight
	if huntlimit.GemItems == nil {
		huntlimit.GemItems = self.GetHunterGemItems(dropConfig)
	}

	msg.Cid = "start_hunt_fight"
	msg.AddItem = huntlimit.GemItems
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("start_hunt_fight", smsg)
	return true
}

// 会长狩猎开启
func (self *ModUnion) OpenHuntFight(nType int) bool {
	if nType == UNION_HUNT_TYPE_NOMAL {
		return false
	}

	config := GetCsvMgr().GetUnionHuntConfigByID(nType)
	if nil == config {
		return false
	}

	if config.Cost <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	member, success := self.GetUnionMember(data, self.player.GetUid())
	if !success {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	if member.Position > UNION_POSITION_VICE_MASTER {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	nowTime := TimeServer().Unix()

	var mastermsg S2M_UnionOpenHunter
	mastermsg.Uid = self.player.GetUid()
	mastermsg.Unionuid = self.Sql_UserUnionInfo.Unionid
	mastermsg.Type = nType
	mastermsg.Cost = config.Cost
	mastermsg.Time = config.Time
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_OPEN_HUNTER, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_UnionOpenHunter
	json.Unmarshal([]byte(ret.Data), &backmsg)

	unionHunter := backmsg.UnionHunter
	data.MinActivityPoint(config.Cost)
	if !data.OpenHunterFight(nType, nowTime+int64(config.Time)) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}
	var msg S2C_OpenHuntFight
	msg.Cid = "open_hunt_fight"
	msg.Type = nType
	msg.Info = unionHunter
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("open_hunt_fight", smsg)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, UNION_ACTIVITY_POINT, -config.Cost, self.Sql_UserUnionInfo.Unionid, data.Level, "激活公会挑战", data.ActivityPoint, 0, self.player)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_HUNTER_OPEN, config.Level, self.Sql_UserUnionInfo.Unionid, data.Level, "激活公会挑战关卡", 0, 0, self.player)

	return true
}

// 结束狩猎战斗
func (self *ModUnion) EndHuntFight(nType int, nDamage int64, battleInfo *BattleInfo) bool {
	config := GetCsvMgr().GetUnionHuntConfigByID(nType)
	if nil == config {
		return false
	}

	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	csv_vip := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if csv_vip == nil || len(csv_vip.GuildHunting) < 2 {
		return false
	}

	if self.Sql_UserUnionInfo.HuntStart != nType {
		return false
	}

	var mastermsg S2M_UnionEndHunter
	mastermsg.Uid = self.player.GetUid()
	mastermsg.Unionuid = self.Sql_UserUnionInfo.Unionid
	mastermsg.Type = nType
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_END_HUNTER, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_UnionEndHunter
	json.Unmarshal([]byte(ret.Data), &backmsg)
	unionHunter := backmsg.UnionHunter

	var huntlimit *UserHuntLimit = nil
	for _, v := range self.Sql_UserUnionInfo.huntLimit {
		if v.Type == nType {
			huntlimit = v
			break
		}
	}

	if nil == huntlimit {
		huntlimit = &UserHuntLimit{nType, 0, 0, unionHunter.EndTime, 0, 0, nil, nil, nil}
		self.Sql_UserUnionInfo.huntLimit = append(self.Sql_UserUnionInfo.huntLimit, huntlimit)
	} else if huntlimit.EndTime != unionHunter.EndTime {
		huntlimit.JoinCount = 0
		huntlimit.SweepCount = 0
		huntlimit.GemItems = nil
		huntlimit.EndTime = unionHunter.EndTime
	}

	maxCount := 0
	if nType == UNION_HUNT_TYPE_NOMAL {
		maxCount = csv_vip.GuildHunting[0]
	} else {
		maxCount = csv_vip.GuildHunting[1]
	}
	if huntlimit.JoinCount >= maxCount {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	huntlimit.JoinCount += 1

	self.Sql_UserUnionInfo.HuntStart = 0

	battleInfo.Id = GetFightMgr().GetFightID(self.player.GetServerId())
	battleInfo.Type = BATTLE_TYPE_PVE
	battleInfo.UserInfo[0].Uid = self.player.GetUid()
	battleInfo.UserInfo[0].Level = self.player.Sql_UserBase.Level
	battleInfo.UserInfo[0].Icon = self.player.Sql_UserBase.IconId
	battleInfo.UserInfo[0].Portrait = self.player.Sql_UserBase.Portrait
	battleInfo.UserInfo[0].UnionName = self.player.GetUnionName()
	battleInfo.UserInfo[0].Name = self.player.GetName()

	//var army *ArmyInfo = nil
	//for _, p := range battleInfo.UserInfo {
	//	for _, q := range p.HeroInfo {
	//		if q.ArmyInfo != nil {
	//			army = q.ArmyInfo
	//			break
	//		}
	//	}
	//}

	battleRecord := BattleRecord{}
	battleRecord.Id = battleInfo.Id
	battleRecord.Side = 1
	battleRecord.Result = 0
	battleRecord.Type = BATTLE_TYPE_PVE
	battleRecord.Time = TimeServer().Unix()
	//if army != nil {
	//	battleRecord.FightInfo[0] = GetRobotMgr().GetPlayerFightInfoWithArmyByPos(self.player, 0, 0, TEAMTYPE_UNION_HUNT, army.SelfKey, army.Pos)
	//	self.player.GetModule("friend").(*ModFriend).SetUseSign(HIRE_MOD_TOWER, 1)
	//} else {
	battleRecord.FightInfo[0] = GetRobotMgr().GetPlayerFightInfoByPos(self.player, 0, 0, TEAMTYPE_UNION_HUNT)
	//}
	battleRecord.RandNum = battleInfo.Random
	battleRecord.LevelID = battleInfo.LevelID

	data.AddUnionHuntDamage(nType, nDamage, self.player, self.Sql_UserUnionInfo.Position, battleInfo, &battleRecord)

	ret1 := make(map[int]*Item, 0)
	ret2 := make(map[int]*Item, 0)
	ret3 := make(map[int]*Item, 0)
	dropConfig := GetCsvMgr().GetUnionHuntDropConfig(config.Group, nDamage)
	if dropConfig != nil {
		temp := self.player.AddObjectLst(dropConfig.Personalitem, dropConfig.Personalnum, "公会挑战", int(nDamage), LOG_FIGHT, 0)
		AddItemMapHelper2(ret1, temp)
	}
	dropsID := GetCsvMgr().GetUnionHuntDrop(self.player, config.Group, nDamage)
	dropItem := GetLootMgr().LootItems(dropsID, self.player)
	temp := self.player.AddObjectItemMap(dropItem, "公会挑战", int(nDamage), LOG_FIGHT, 0)
	AddItemMapHelper2(ret2, temp)

	temp = self.player.AddObjectPassItem(huntlimit.GemItems, "公会挑战", int(nDamage), LOG_FIGHT, 0)
	self.AddGemAwardRecord(self.player.GetName(), self.Sql_UserUnionInfo.Position, nType, temp)
	AddItemMapHelper2(ret3, temp)
	huntlimit.GemItems = nil

	awardcount := 0
	for _, g := range dropItem {
		awardcount += g.ItemNum
	}

	if nDamage > huntlimit.MaxDamage {
		huntlimit.MaxDamage = nDamage
		huntlimit.AwardCount = awardcount
		huntlimit.BattleInfo = battleInfo
		huntlimit.BattleRecord = &battleRecord
	}

	self.player.HandleTask(TASK_TYPE_UNION_HUNT_COUNT, nType, 0, 0)
	self.player.HandleTask(TASK_TYPE_UNION_HUNT_COUNT, nType, 1, 0)
	self.player.GetModule("task").(*ModTask).SendUpdate()
	GetUnionMgr().UpdateMember(self.Sql_UserUnionInfo.Unionid, self.player.GetUid(), 0, false, false)

	privilegeItems := make(map[int]*Item)
	value := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValue(PRIVILEGE_UNION)

	var msg S2C_EndHuntFight
	msg.Cid = "end_hunt_fight"
	msg.Type = nType
	msg.Info = huntlimit
	msg.Damage = nDamage
	for _, v := range ret1 {
		msg.Item1 = append(msg.Item1, PassItem{v.ItemId, v.ItemNum})
		if value > 0 && v.ItemId == ITEM_UNION {
			AddItemMapHelper3(privilegeItems, ITEM_UNION, v.ItemNum*value/100)
		}
	}
	for _, v := range ret3 {
		msg.Item1 = append(msg.Item1, PassItem{v.ItemId, v.ItemNum})
		if value > 0 && v.ItemId == ITEM_UNION {
			AddItemMapHelper3(privilegeItems, ITEM_UNION, v.ItemNum*value/100)
		}
	}
	for _, v := range ret2 {
		msg.Item2 = append(msg.Item2, PassItem{v.ItemId, v.ItemNum})
		if value > 0 && v.ItemId == ITEM_UNION {
			AddItemMapHelper3(privilegeItems, ITEM_UNION, v.ItemNum*value/100)
		}
	}
	msg.GetPrivilegeItems = self.player.AddObjectItemMap(privilegeItems, "公会挑战", int(nDamage), LOG_FIGHT, 0)

	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("end_hunt_fight", smsg)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_HUNTER_END, config.Level, int(nDamage), LOG_FIGHT, "公会挑战", 0, data.Id, self.player)

	return true
}

// 获得狩猎信息
func (self *ModUnion) GetHuntInfo() bool {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		return false
	}

	csv_vip := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if csv_vip == nil || len(csv_vip.GuildHunting) < 2 {
		return false
	}

	for i := UNION_HUNT_TYPE_NOMAL; i < UNION_HUNT_TYPE_MAX; i++ {
		unionHunter := &JS_UnionHunt{}
		for _, v := range data.huntInfo {
			if v.Type == i {
				unionHunter = v
				break
			}
		}

		if nil == unionHunter {
			continue
		}

		var huntlimit *UserHuntLimit = nil

		for _, v := range self.Sql_UserUnionInfo.huntLimit {
			if v.Type == i {
				huntlimit = v
				break
			}
		}

		if nil == huntlimit {
			huntlimit = &UserHuntLimit{i, 0, 0, unionHunter.EndTime, 0, 0, nil, nil, nil}
			self.Sql_UserUnionInfo.huntLimit = append(self.Sql_UserUnionInfo.huntLimit, huntlimit)
		} else if huntlimit.EndTime != unionHunter.EndTime {
			huntlimit.JoinCount = 0
			huntlimit.SweepCount = 0
			huntlimit.GemItems = nil
			huntlimit.EndTime = unionHunter.EndTime
		}
	}

	var msg S2C_GetHuntInfo
	msg.Cid = "get_hunt_info"
	msg.UserHuntLimit = self.Sql_UserUnionInfo.huntLimit
	msg.GuildHunting = csv_vip.GuildHunting
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("get_hunt_info", smsg)

	return true
}

// 结束狩猎战斗
func (self *ModUnion) SweepHuntFight(nType int) bool {
	config := GetCsvMgr().GetUnionHuntConfigByID(nType)
	if nil == config {
		return false
	}

	dropGemConfig := GetCsvMgr().GetUnionHuntDropConfig(config.Group, 0)
	if nil == dropGemConfig {
		return false
	}

	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	if TimeServer().Unix()-self.Sql_UserUnionInfo.LastUpdTime < UNION_HUNTER_CD && nType == UNION_HUNT_TYPE_SPECIAL {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("离开公会24小时内不能参与狩猎"))
		return false
	}

	csv_vip := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if csv_vip == nil || len(csv_vip.GuildHunting) < 2 {
		return false
	}

	if csv_vip.GuildSweep != 1 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	var mastermsg S2M_UnionEndHunter
	mastermsg.Uid = self.player.GetUid()
	mastermsg.Unionuid = self.Sql_UserUnionInfo.Unionid
	mastermsg.Type = nType
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_END_HUNTER, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_UnionEndHunter
	json.Unmarshal([]byte(ret.Data), &backmsg)
	unionHunter := backmsg.UnionHunter
	if nil == unionHunter {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return false
	}

	//var playerRec *JS_UnionCopyDpsTop = nil
	//for _, v := range unionHunter.TopDps {
	//	if v.Uid == self.player.GetUid() {
	//		playerRec = v
	//		break
	//	}
	//}
	//if playerRec == nil {
	//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
	//	return false
	//}

	var huntlimit *UserHuntLimit = nil

	for _, v := range self.Sql_UserUnionInfo.huntLimit {
		if v.Type == nType {
			huntlimit = v
			break
		}
	}

	if nil == huntlimit {
		return false
	}

	if unionHunter.EndTime != huntlimit.EndTime {
		return false
	}

	if huntlimit.MaxDamage <= 0 || huntlimit.BattleInfo == nil {
		return false
	}

	maxCount := 0
	if nType == UNION_HUNT_TYPE_NOMAL {
		maxCount = csv_vip.GuildHunting[0]
	} else {
		maxCount = csv_vip.GuildHunting[1]
	}
	if huntlimit.JoinCount >= maxCount {
		return false
	}

	huntlimit.BattleInfo.Id = GetFightMgr().GetFightID(self.player.GetServerId())
	huntlimit.BattleRecord.Id = huntlimit.BattleInfo.Id

	data.AddUnionHuntDamage(nType, huntlimit.MaxDamage, self.player, self.Sql_UserUnionInfo.Position, huntlimit.BattleInfo, huntlimit.BattleRecord)

	ret1 := make(map[int]*Item, 0)
	ret2 := make(map[int]*Item, 0)
	ret3 := make(map[int]*Item, 0)
	dropConfig := GetCsvMgr().GetUnionHuntDropConfig(config.Group, huntlimit.MaxDamage)
	if dropConfig != nil {
		temp := self.player.AddObjectLst(dropConfig.Personalitem, dropConfig.Personalnum, "公会挑战", int(huntlimit.MaxDamage), LOG_SWEEP, 0)

		AddItemMapHelper2(ret1, temp)
	}
	dropsID := GetCsvMgr().GetUnionHuntDrop(self.player, config.Group, huntlimit.MaxDamage)
	dropItem := GetLootMgr().LootItems(dropsID, self.player)
	temp := self.player.AddObjectItemMap(dropItem, "公会挑战", int(huntlimit.MaxDamage), LOG_SWEEP, 0)
	AddItemMapHelper2(ret2, temp)

	if huntlimit.GemItems == nil {
		temp = self.player.AddObjectPassItem(self.GetHunterGemItems(dropGemConfig), "公会挑战", int(huntlimit.MaxDamage), LOG_SWEEP, 0)
		self.AddGemAwardRecord(self.player.GetName(), self.Sql_UserUnionInfo.Position, nType, temp)
		AddItemMapHelper2(ret3, temp)
	} else {
		temp = self.player.AddObjectPassItem(huntlimit.GemItems, "公会挑战", int(huntlimit.MaxDamage), LOG_SWEEP, 0)
		self.AddGemAwardRecord(self.player.GetName(), self.Sql_UserUnionInfo.Position, nType, temp)
		AddItemMapHelper2(ret3, temp)
	}
	huntlimit.GemItems = nil

	huntlimit.JoinCount += 1

	self.player.HandleTask(TASK_TYPE_UNION_HUNT_COUNT, nType, 0, 0)
	self.player.HandleTask(TASK_TYPE_UNION_HUNT_COUNT, nType, 1, 0)
	self.player.GetModule("task").(*ModTask).SendUpdate()
	GetUnionMgr().UpdateMember(self.Sql_UserUnionInfo.Unionid, self.player.GetUid(), 0, false, false)

	privilegeItems := make(map[int]*Item)
	value := self.player.GetModule("interstellar").(*ModInterStellar).GetPrivilegeValue(PRIVILEGE_UNION)

	var msg S2C_SweepHuntFight
	msg.Cid = "sweep_hunt_fight"
	msg.Type = nType
	msg.Info = huntlimit
	msg.Damage = huntlimit.MaxDamage
	for _, v := range ret1 {
		msg.Item1 = append(msg.Item1, PassItem{v.ItemId, v.ItemNum})
		if value > 0 && v.ItemId == ITEM_UNION {
			AddItemMapHelper3(privilegeItems, ITEM_UNION, v.ItemNum*value/100)
		}
	}
	for _, v := range ret3 {
		msg.Item1 = append(msg.Item1, PassItem{v.ItemId, v.ItemNum})
		if value > 0 && v.ItemId == ITEM_UNION {
			AddItemMapHelper3(privilegeItems, ITEM_UNION, v.ItemNum*value/100)
		}
	}
	for _, v := range ret2 {
		msg.Item2 = append(msg.Item2, PassItem{v.ItemId, v.ItemNum})
		if value > 0 && v.ItemId == ITEM_UNION {
			AddItemMapHelper3(privilegeItems, ITEM_UNION, v.ItemNum*value/100)
		}
	}
	msg.GetPrivilegeItems = self.player.AddObjectItemMap(privilegeItems, "公会挑战", int(huntlimit.MaxDamage), LOG_SWEEP, 0)
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("sweep_hunt_fight", smsg)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_HUNTER_END, config.Level, int(huntlimit.MaxDamage), LOG_SWEEP, "公会挑战", 0, data.Id, self.player)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_HUNTER_SWEEP, nType, 0, 0, "工会副本扫荡", 0, 0, self.player)
	return true
}

//发送邮件
func (self *ModUnion) SendMail(title string, text string) {
	unionID := self.player.GetUnionId()
	data := GetUnionMgr().GetUnion(unionID)
	if data == nil {
		return
	}

	if self.Sql_UserUnionInfo.Position > UNION_POSITION_VICE_MASTER {
		return
	}

	now := TimeServer().Unix()
	if data.MailCD != 0 {
		cd := now - data.MailCD
		if cd < DAY_SECS {
			left := DAY_SECS - cd
			self.player.SendErrInfo("err", fmt.Sprintf("%d小时%d分钟%d秒后可发送邮件", (left/3600), (left%3600)/60, (left%60)))
			return
		}
	}
	var msg S2C_UnionSendMail
	msg.Cid = "unionsendmail"
	if GetUnionMgr().SendMail(unionID, self.player.GetUid(), title, text) {
		msg.Ret = 1
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_UNION_SEND_MAIL, unionID, data.Level, 0, "发送公会全员邮件", 0, 0, self.player)
	} else {
		msg.Ret = 0
	}
	data.AddMailCD()
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("unionsendmail", smsg)
}

//设置无畏之手
func (self *ModUnion) SetBraveHand(unionid int, destuid int64, op int) {
	data := GetUnionMgr().GetUnion(unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	selfdata, ok := self.GetUnionMember(data, self.player.Sql_UserBase.Uid)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	destdata, ok2 := self.GetUnionMember(data, destuid)
	if !ok2 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	if selfdata.Position > UNION_POSITION_VICE_MASTER {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	if op == 1 {
		if destdata.BraveHand == 1 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("玩家已经是无畏之手"))
			return
		}
		if !GetCsvMgr().GuildIsLevelAndPassOpen(destdata.Level, destdata.Stage, OPEN_LEVEL_HIRE) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("玩家未开启佣兵，无法设置为无畏之手"))
			return
		}
	} else {
		if destdata.BraveHand == 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("玩家不是无畏之手"))
			return
		}
	}

	csv, _ := GetCsvMgr().CommunityConfig[data.Level]
	if csv == nil {
		return
	}

	GetUnionMgr().CheckBraveHand(data)

	if op == 1 && len(data.braveHand) >= csv.Fearless {
		endTime := int64(0)
		for _, v := range data.braveHand {
			if v.EndTime == 0 {
				continue
			}

			if endTime == 0 {
				endTime = v.EndTime
			} else {
				if v.EndTime < endTime {
					endTime = v.EndTime
				}
			}
		}

		if endTime > 0 {
			left := endTime - TimeServer().Unix()
			self.player.SendErrInfo("err", fmt.Sprintf("%d小时%d分钟%d秒后可设置无畏之手", left/3600, (left%3600)/60, (left%60)))
		} else {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("无畏之手已满"))
		}
		return
	}

	if !GetUnionMgr().SetBraveHand(self.player.GetUid(), unionid, destuid, op) {
		return
	}

	var msg S2C_SetBraveHand
	msg.Cid = "setbravehand"
	msg.Uid = self.player.Sql_UserBase.Uid
	msg.Destuid = destuid
	msg.Op = op
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("setbravehand", smsg)
}

func (self *ModUnion) GetBattleInfo(id int64, nType int) *BattleInfo {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return nil
	}

	var battleInfo BattleInfo
	//value, flag, err := HGetRedis(`san_huntbattleinfo`, fmt.Sprintf("%d", id))
	//if err != nil {
	//	self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
	//	return &battleInfo
	//}
	//if flag {
	//	err := json.Unmarshal([]byte(value), &battleInfo)
	//	if err != nil {
	//		self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
	//		return &battleInfo
	//	}
	//}

	var mastermsg S2M_UnionGetBattleInfo
	mastermsg.FightID = id
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GET_INFO, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
		return &battleInfo
	}

	var backmsg M2S_UnionGetBattleInfo
	json.Unmarshal([]byte(ret.Data), &backmsg)
	err := json.Unmarshal([]byte(backmsg.Info), &battleInfo)
	if err == nil && battleInfo.Id != 0 {
		return &battleInfo
	}

	return nil
}

func (self *ModUnion) GetBattleRecord(id int64, nType int) *BattleRecord {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return nil
	}

	var battleRecord BattleRecord
	//value, flag, err := HGetRedis(`san_huntbattlerecord`, fmt.Sprintf("%d", id))
	//if err != nil {
	//	self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
	//	return &battleRecord
	//}
	//if flag {
	//	err := json.Unmarshal([]byte(value), &battleRecord)
	//	if err != nil {
	//		self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
	//		return &battleRecord
	//	}
	//}

	var mastermsg S2M_UnionGetBattleRecord
	mastermsg.FightID = id
	ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_GET_RECORD, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
		return &battleRecord
	}

	var backmsg M2S_UnionGetBattleRecord
	json.Unmarshal([]byte(ret.Data), &backmsg)
	err := json.Unmarshal([]byte(backmsg.Record), &battleRecord)
	if err == nil && battleRecord.Id != 0 {
		return &battleRecord
	}

	return nil
}

func (self *ModUnion) GetDpsTop(nType int) {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	info := &JS_UnionHunt{}
	for _, v := range data.huntInfo {
		if nType == v.Type {
			info = v
		}
	}
	if nil == info {
		return
	}

	var msg S2C_HuntDpsTop
	msg.Cid = "get_hunt_dps_top"
	msg.Type = nType
	msg.DpsTop = info.TopDps
	self.player.SendMsg("get_hunt_dps_top", HF_JtoB(&msg))
}

func (self *ModUnion) UpdateFight() {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	GetUnionMgr().UpdateMember(self.Sql_UserUnionInfo.Unionid, self.player.GetUid(), 0, false, false)
}
func (self *ModUnion) AddGemAwardRecord(name string, position, nType int, items []PassItem) {
	data := GetUnionMgr().GetUnion(self.Sql_UserUnionInfo.Unionid)
	if data == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	nCount := 0
	for _, v := range items {
		if v.ItemID == 30020011 {
			nCount++
		}
	}

	if nCount > 0 {
		var mastermsg S2M_UnionAddGemAwardRecord
		mastermsg.Unionuid = self.Sql_UserUnionInfo.Unionid
		mastermsg.Name = name
		mastermsg.Position = position
		mastermsg.Type = nType
		mastermsg.Count = nCount
		ret := GetMasterMgr().UnionRPC.UnionAction(RPC_UNION_ADD_GEM_AWARD_RECORD, &mastermsg)
		if ret == nil || ret.RetCode != UNION_SUCCESS {
			self.player.SendErr(GetCsvMgr().GetText("STR_UIArena_Non_Report"))
			return
		}
		data.AddGemAwardRecord(name, position, nType, nCount)
	}
}
