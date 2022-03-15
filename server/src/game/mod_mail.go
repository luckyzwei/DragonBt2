package game

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	//"time"
)

const MAX_MAIL_LEN = 99

const (
	MAIL_ID_GET_ITEM_WARORDER_1    = 1
	MAIL_ID_GET_ITEM_WARORDER_2    = 2
	MAIL_ID_GET_ITEM               = 3
	MAIL_ID_GET_ITEM_MONEY         = 4 //数值型
	MAIL_ID_GET_UNION_HUNTER_AWARD = 5
	MAIL_ID_GET_UNION_HUNTER_OPEN  = 6
	MAIL_ID_GET_ARENA_DAY          = 9
	MAIL_ID_GET_ARENA_END          = 10
	MAIL_ID_GET_DAILY_RECHARGE     = 21 //每日首充
	MAIL_ID_GET_SEVEN              = 22 //7日发送的邮件
	MAIL_ID_GET_FUND               = 25 //超值好礼是过期领取
	MAIL_ID_ACTIVITYBOSS_RANK      = 26 //排行
	MAIL_ID_ACTIVITYBOSS_TASK      = 27 //任务
	MAIL_ID_ACTIVITYBOSS_GOLD      = 28 //过期兑换
	MAIL_ID_REMOVE_EXCHANEG        = 29 //代币回收
	MAIL_ID_CROSSARENA_SUB         = 30 //跨服竞技 段位奖励
	MAIL_ID_CROSSARENA_RANK        = 31 //跨服竞技 至尊奖励
	MAIL_ID_CROSSARENA_TASK        = 32 //跨服竞技 补领任务
	MAIL_ID_CONSUMERTOP_RANK       = 33 //
	MAIL_ID_CONSUMERTOP_SCORE      = 34 //
	MAIL_ID_CROSSARENA_3V3_SUB     = 35 //跨服竞技3v3 段位奖励
	MAIL_ID_CROSSARENA_3V3_RANK    = 36 //跨服竞技3v3 至尊奖励
	MAIL_ID_CROSSARENA_3V3_TASK    = 37 //跨服竞技3v3 补领任务
)

//! 客户端错误日志
type SQL_ErrorLog struct {
	Id        int64 //! 邮件ID
	Uid       int64 //! 角色ID
	Stack     string
	ErrorInfo string
	Param1    string
	Time      int64 //! 操作时间
}

//! 邮件数据库
type JS_OneMail struct {
	Dyid        int64      `json:"dyid"`
	Uid         int64      `json:"uid"`
	Mailid      int64      `json:"mailid"`
	Mailtype    int        `json:"mailtype"`
	Mailstate   int        `json:"mailstate"`
	Title       string     `json:"titile"`
	Content     string     `json:"content"`
	Sender      string     `json:"sender"`
	Item        []PassItem `json:"item"`
	Inserttime  int64      `json:"inserttime"`
	Mailsubtype int        `json:"mailsubtype"`
}

const (
	MAIL_NO_ITEM         = 0 // 没有物品邮件
	MAIL_CAN_ALL_GET     = 1 // 可一键领取邮件
	MAIL_CAN_NOT_ALL_GET = 2 // 不可一键领取邮件
)

//! 背包数据库
type San_Mail struct {
	Uid  int64  //! 角色UID
	Info string //! 邮件列表
	Recv string //! 全服邮件
	Camp string //! 阵营提示记录

	info map[string]*JS_OneMail //! 邮件信息
	recv map[int]int            //! 系统邮件收取
	camp map[int]int            //! 阵营消息收取

	DataUpdate
}

//! 背包
type ModMail struct {
	player     *Player
	Sql_Mail   San_Mail
	MaxId      int64
	InfoLocker *sync.RWMutex
}

func (self *ModMail) OnGetData(player *Player) {
	self.player = player
	self.InfoLocker = new(sync.RWMutex)
	sql := fmt.Sprintf("select * from `san_usermail` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Mail, "san_usermail", self.player.ID)

	if self.Sql_Mail.Uid <= 0 { // 创建新号
		self.Sql_Mail.Uid = self.player.ID
		self.Sql_Mail.info = make(map[string]*JS_OneMail)
		self.Sql_Mail.recv = make(map[int]int, 0)
		self.Sql_Mail.camp = make(map[int]int, 0)

		self.Encode()
		InsertTable("san_usermail", &self.Sql_Mail, 0, true)
		self.Sql_Mail.Init("san_usermail", &self.Sql_Mail, true)
	} else {
		self.Decode()
		if self.Sql_Mail.info == nil {
			LogError("Mail Info Nil Error:", self.player.Sql_UserBase.Uid)
			self.Sql_Mail.info = make(map[string]*JS_OneMail, 0)
		}
		if self.Sql_Mail.recv == nil {
			LogError("Mail Recv Nil Error:", self.player.Sql_UserBase.Uid)
			self.Sql_Mail.recv = make(map[int]int, 0)
		}
		if self.Sql_Mail.camp == nil {
			LogError("Mail Camp Nil Error:", self.player.Sql_UserBase.Uid)
			self.Sql_Mail.camp = make(map[int]int, 0)
		}
		self.Sql_Mail.Init("san_usermail", &self.Sql_Mail, true)

		for _, value := range self.Sql_Mail.info {
			if value.Dyid > self.MaxId {
				self.MaxId = value.Dyid
			}
		}
	}

	GetMailMgr().GetMail(self.player)
	self.CheckMailOut()
}

func (self *ModMail) CheckMailOut() {
	//删除过期邮件
	tNowTime := TimeServer().Unix()
	for mailId, value := range self.Sql_Mail.info {
		if tNowTime-value.Inserttime >= DAY_SECS*7 {
			delete(self.Sql_Mail.info, mailId)
			GetServer().SqlMailLog(self.player.Sql_UserBase.Uid, value.Title+":"+value.Content, HF_JtoA(value.Item))
		}
	}
}

func (self *ModMail) OnGetOtherData() {

}

func (self *ModMail) Decode() { //! 将数据库数据写入data
	mailErrDecode := json.Unmarshal([]byte(self.Sql_Mail.Info), &self.Sql_Mail.info)
	if mailErrDecode != nil {
		mailInfo := self.Sql_Mail.Info
		if strings.Contains(mailInfo, "\n") || strings.Contains(mailInfo, "\r") ||
			strings.Contains(mailInfo, "\t") || strings.Contains(mailInfo, "\v") {
			mailInfo = strings.Replace(mailInfo, "\n", "\\n", -1)
			mailInfo = strings.Replace(mailInfo, "\r", "\\r", -1)
			mailInfo = strings.Replace(mailInfo, "\t", "\\t", -1)
			mailInfo = strings.Replace(mailInfo, "\v", "\\v", -1)
		}

		mailErrDecode := json.Unmarshal([]byte(mailInfo), &self.Sql_Mail.info)
		if mailErrDecode != nil {
			LogError("邮件解析错误:", mailErrDecode.Error(), ", json:", self.Sql_Mail.Info)
		}
	}

	json.Unmarshal([]byte(self.Sql_Mail.Recv), &self.Sql_Mail.recv)
	json.Unmarshal([]byte(self.Sql_Mail.Camp), &self.Sql_Mail.camp)
}

func (self *ModMail) Encode() { //! 将data数据写入数据库
	self.Sql_Mail.Info = HF_JtoA(&self.Sql_Mail.info)
	self.Sql_Mail.Recv = HF_JtoA(&self.Sql_Mail.recv)
	self.Sql_Mail.Camp = HF_JtoA(&self.Sql_Mail.camp)
}

func (self *ModMail) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "getmailallitem":
		self.GetMailAllItem()
		return true
	case "detelemailall":
		self.DeleteMailAll()
		return true
	case "getmailitem":
		var s2c_msg C2S_GetMailItem
		json.Unmarshal(body, &s2c_msg)
		self.GetMailItem(s2c_msg.Mailid)
		return true
	case "readmail":
		var s2c_msg C2S_ReadMail
		json.Unmarshal(body, &s2c_msg)
		self.ReadMail(s2c_msg.Mailid)
		return true
	case "delmail":
		var s2c_msg C2S_DelMail
		json.Unmarshal(body, &s2c_msg)
		self.DelMail(s2c_msg.Mailid, true)
		return true
	case "inserterror":
		var s2c_msg C2S_InsertError
		json.Unmarshal(body, &s2c_msg)
		self.InsertError(&s2c_msg)
		return true
	}

	return false
}

func (self *ModMail) OnSave(sql bool) {
	self.InfoLocker.RLock()
	defer self.InfoLocker.RUnlock()
	self.Encode()
	self.Sql_Mail.Update(sql)
}

func (self *ModMail) AddMailWithItems(mailtype int, title, content string, item []PassItem) {
	//理不清消息顺序还是强行同步好  20190724  by zy
	//self.AddMail(mailtype, 1, 0, title, content, GetCsvMgr().GetText("STR_SYS"), item, false, 0)

	self.AddMail(mailtype, 1, 0, title, content, GetCsvMgr().GetText("STR_SYS"), item, true, 0)
}

func (self *ModMail) AddMail(mailtype, mailsubtype, mailstate int, title, content, sender string, item []PassItem, send bool, _time int64) {
	// 判断邮件长度,超过长度直接丢弃
	self.InfoLocker.Lock()
	if len(self.Sql_Mail.info) >= MAX_MAIL_LEN {
		self.ClearOverlapMail(true)
	}
	self.InfoLocker.Unlock()

	data := new(JS_OneMail)
	data.Uid = self.player.ID
	data.Mailtype = mailtype
	data.Mailsubtype = mailsubtype
	data.Mailstate = mailstate
	data.Title = title
	data.Content = content
	data.Sender = sender
	//以前用得太随意，要改的地方太多，这里修正下
	if len(item) == 0 {
		data.Mailtype = MAIL_NO_ITEM
	}
	data.Item = item
	if _time == 0 {
		data.Inserttime = TimeServer().Unix()
	} else {
		data.Inserttime = _time
	}
	data.Mailsubtype = mailsubtype
	self.MaxId++
	data.Dyid = self.MaxId
	data.Mailid = data.Dyid

	self.InfoLocker.Lock()
	self.Sql_Mail.info[fmt.Sprintf("%d", data.Dyid)] = data

	if send {
		self.SendInfo()
	}
	self.InfoLocker.Unlock()
}

func (self *ModMail) AddGlobalMail(mailid, mailtype, mailsubtype, mailstate int, title, content, sender string,
	item []PassItem, send bool, _time int64) {
	// 判断邮件长度,超过长度直接丢弃
	self.InfoLocker.Lock()
	if len(self.Sql_Mail.info) >= MAX_MAIL_LEN {
		self.ClearOverlapMail(true)
	}
	self.InfoLocker.Unlock()

	data := new(JS_OneMail)
	data.Uid = self.player.ID
	data.Mailtype = mailtype
	data.Mailsubtype = mailsubtype
	data.Mailstate = mailstate
	data.Title = title
	data.Content = content
	data.Sender = sender
	data.Item = item
	if _time == 0 {
		data.Inserttime = TimeServer().Unix()
	} else {
		data.Inserttime = _time
	}
	if len(item) == 0 {
		data.Mailtype = MAIL_NO_ITEM
	}
	self.MaxId++
	data.Dyid = self.MaxId
	data.Mailid = data.Dyid
	self.InfoLocker.Lock()
	if self.Sql_Mail.info == nil {
		LogError("Mail Info Nil Error:", self.player.Sql_UserBase.Uid)
		self.Sql_Mail.info = make(map[string]*JS_OneMail, 0)
	}
	if self.Sql_Mail.recv == nil {
		LogError("Mail Recv Nil Error:", self.player.Sql_UserBase.Uid)
		self.Sql_Mail.recv = make(map[int]int, 0)
	}
	self.Sql_Mail.info[fmt.Sprintf("%d", data.Dyid)] = data
	self.Sql_Mail.recv[mailid] = 1

	if send {
		self.SendInfo()
	}
	self.InfoLocker.Unlock()
}

func (self *ModMail) ClearOverlapMail(force bool) {
	//self.InfoLocker.Lock()
	//defer self.InfoLocker.Unlock()

	//删除空邮件
	minMailId := self.MaxId
	tNowTime := TimeServer().Unix()
	for mailId, value := range self.Sql_Mail.info {
		if len(value.Item) == 0 && tNowTime-value.Inserttime >= DAY_SECS*3 {
			delete(self.Sql_Mail.info, mailId)
			GetServer().SqlMailLog(self.player.Sql_UserBase.Uid, value.Title+":"+value.Content, HF_JtoA(value.Item))
		}

		if value.Mailid < minMailId {
			minMailId = value.Mailid
		}
	}

	if force == true {
		//delete(self.Sql_Mail.info, fmt.Sprint("%d", minMailId))
		value, ok := self.Sql_Mail.info[fmt.Sprintf("%d", minMailId)]
		if !ok {
			return
		}
		delete(self.Sql_Mail.info, fmt.Sprintf("%d", minMailId))

		GetServer().SqlMailLog(self.player.Sql_UserBase.Uid, value.Title+":"+value.Content, HF_JtoA(value.Item))
	}
}

func (self *ModMail) GetMailAllItem() {
	self.InfoLocker.Lock()
	defer self.InfoLocker.Unlock()

	heroItem := make([]PassItem, 0)
	var msg S2C_MailAllItem
	msg.Cid = "mailallitem"
	times := 0
	for _, value := range self.Sql_Mail.info {
		if value.Mailsubtype != 2 && value.Mailtype == MAIL_CAN_ALL_GET && value.Mailstate != 2 {
			//检查道具是否可以领取
			isLimit := false
			for i := 0; i < len(value.Item); i++ {
				if self.player.CheckItemLimit(value.Item[i].ItemID, value.Item[i].Num) {
					self.player.SendErr(GetCsvMgr().GetText("STR_MOD_SHOP_REACH_THE_UPPER_LIMIT"))
					isLimit = true
					break
				}
			}
			if isLimit {
				continue
			}
			//! 加道具  英雄的处理放到最后，从而调整消息顺序给客户端展示
			times++
			for i := 0; i < len(value.Item); i++ {
				itemConfig := GetCsvMgr().GetItemConfig(value.Item[i].ItemID)
				if itemConfig != nil && itemConfig.ItemType == ITEM_TYPE_HERO {
					heroItem = append(heroItem, PassItem{value.Item[i].ItemID, value.Item[i].Num})
				} else {
					value.Item[i].ItemID, value.Item[i].Num = self.player.AddObject(value.Item[i].ItemID, value.Item[i].Num,
						0, 0, 0, GetCsvMgr().GetText("STR_MAIL_AWARD_GOT"))
				}
				msg.Item = append(msg.Item, PassItem{value.Item[i].ItemID, value.Item[i].Num})
			}
			//self.DelMail(value.Mailid, false)
			//delete(self.Sql_Mail.info, fmt.Sprintf("%d", value.Mailid))
			value.Mailstate = 2

			GetServer().SqlMailLog(self.player.Sql_UserBase.Uid, value.Title+":"+value.Content, HF_JtoA(value.Item))
		}
	}

	if times > 0 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GET_MAIL_ALL, times, 0, 0, "一键领取邮件", 0, 0, self.player)
	}
	self.player.SendMsg("mailallitem", HF_JtoB(&msg))

	if len(self.Sql_Mail.info) > 0 {
		self.SendInfo()
	} else {
		self.player.SendRet2("nomail")
	}

	if len(heroItem) > 0 {
		for i := 0; i < len(heroItem); i++ {
			self.player.AddObject(heroItem[i].ItemID, heroItem[i].Num, 0, 0, 0, GetCsvMgr().GetText("STR_MAIL_AWARD_GOT"))
		}
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_MAIL_GET, 1, 0, 0, "邮件领取", 0, 0, self.player)
}

func (self *ModMail) DeleteMailAll() {
	self.InfoLocker.Lock()
	defer self.InfoLocker.Unlock()

	for _, value := range self.Sql_Mail.info {
		if (value.Mailtype == MAIL_NO_ITEM && value.Mailstate != 0) || value.Mailstate == 2 {
			delete(self.Sql_Mail.info, fmt.Sprintf("%d", value.Mailid))
		}
	}

	if len(self.Sql_Mail.info) > 0 {
		self.SendInfo()
	} else {
		self.player.SendRet2("nomail")
	}
}

func (self *ModMail) GetMailItem(mailid int64) {
	self.InfoLocker.RLock()
	node, ok := self.Sql_Mail.info[fmt.Sprintf("%d", mailid)]
	self.InfoLocker.RUnlock()
	if !ok {
		self.player.SendRet2("nomail")
		return
	}

	if node.Mailtype != MAIL_CAN_ALL_GET && node.Mailtype != MAIL_CAN_NOT_ALL_GET {
		self.player.SendRet2("nomail")
		return
	}

	if node.Mailstate == 2 { //! 已经领取
		self.player.SendRet2("nomail")
		return
	}

	//检查道具是否可以领取
	for i := 0; i < len(node.Item); i++ {
		if self.player.CheckItemLimit(node.Item[i].ItemID, node.Item[i].Num) {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_SHOP_REACH_THE_UPPER_LIMIT"))
			return
		}
	}

	//! 加道具
	for i := 0; i < len(node.Item); i++ {
		node.Item[i].ItemID, node.Item[i].Num = self.player.AddObject(node.Item[i].ItemID, node.Item[i].Num, 0,
			0, 0, GetCsvMgr().GetText("STR_MAIL_AWARD_GOT"))
	}

	//self.DelMail(mailid, true)
	node.Mailstate = 2
	self.player.SendRet("getmailitem", 1)

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_GET_MAIL_SINGLE, int(node.Mailid), 0, 0, "领取邮件奖励", 0, 0, self.player)
}

func (self *ModMail) InsertError(msg *C2S_InsertError) {
	var errorInfo SQL_ErrorLog
	errorInfo.Uid = msg.Uid
	errorInfo.Stack = strings.Replace(msg.Stack, "%", "%%", -1)
	errorInfo.ErrorInfo = strings.Replace(msg.ErrorInfo, "%", "%%", -1)
	errorInfo.Param1 = strings.Replace(msg.Param1, "%", "%%", -1)
	errorInfo.Time = TimeServer().Unix()
	InsertLogTable("san_errorinfo", &errorInfo, 1)
}

func (self *ModMail) DelMail(mailid int64, sendmsg bool) bool {
	self.InfoLocker.Lock()
	defer self.InfoLocker.Unlock()

	value, ok := self.Sql_Mail.info[fmt.Sprintf("%d", mailid)]
	if !ok {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_MAIL_MAIL_DOES_NOT_EXIST"))
		return false
	}

	delete(self.Sql_Mail.info, fmt.Sprintf("%d", mailid))

	GetServer().SqlMailLog(self.player.Sql_UserBase.Uid, value.Title+":"+value.Content, HF_JtoA(value.Item))

	if sendmsg {
		self.SendInfo()
	}
	return true
}

func (self *ModMail) IsRecvSystemMail(mailid int) bool {
	self.InfoLocker.RLock()
	defer self.InfoLocker.RUnlock()

	_, ok := self.Sql_Mail.recv[mailid]
	return ok
}

func (self *ModMail) SetRecvSystemMail(mailid int) {
	self.InfoLocker.Lock()
	defer self.InfoLocker.Unlock()

	self.Sql_Mail.recv[mailid] = 1
}

func (self *ModMail) IsRecvCampRecord(recordid int) bool {
	self.InfoLocker.RLock()
	defer self.InfoLocker.RUnlock()

	_, ok := self.Sql_Mail.camp[recordid]
	return ok
}

func (self *ModMail) SetRecvCampRecord(recordid int) {
	self.InfoLocker.Lock()
	defer self.InfoLocker.Unlock()

	self.Sql_Mail.camp[recordid] = 1

}

func (self *ModMail) ReadMail(mailid int64) {
	self.InfoLocker.RLock()
	defer self.InfoLocker.RUnlock()
	node, ok := self.Sql_Mail.info[fmt.Sprintf("%d", mailid)]
	if !ok {
		self.player.SendRet2("nomail")
		return
	}

	if node.Mailstate != 0 {
		return
	}

	node.Mailstate = 1

	var msg S2C_ReadMail
	msg.Cid = "readmail"
	msg.Index = mailid
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("readmail", smsg)
}

//////////////////////
func (self *ModMail) SendInfo() {
	var msg S2C_MailInfo
	msg.Cid = "maininfo"
	for _, value := range self.Sql_Mail.info {
		msg.Mailinfo = append(msg.Mailinfo, value)
	}
	msg.Gmail = make([]*JS_Mail, 0)
	msg.Uid = self.player.ID
	self.player.SendMsg("maininfo", HF_JtoB(&msg))
}

func (self *San_Mail) Encode() { //! 将data数据写入数据库
	self.Info = HF_JtoA(&self.info)
	self.Recv = HF_JtoA(&self.recv)
	self.Camp = HF_JtoA(&self.camp)
}
