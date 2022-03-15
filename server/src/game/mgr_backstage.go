package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"strconv"
	"time"

	"crypto/md5"
	"encoding/base64"
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime/pprof"
	"sync"
)

type BackStageMgr struct {
	Locker            *sync.RWMutex
	NoticeList        []*NoticeInfo //! 定时邮件,定时发送
	TempNoticeList    []*NoticeInfo //! 临时邮件,临时发送完
	NoticeLocker      *sync.RWMutex
	TempNocticeLocker *sync.RWMutex
	NextSendTime      int64
}

const (
	TEST_RATE_HIGH_CALL      = 1 //高级召唤
	TEST_RATE_ENERGRY_LOOT   = 2 //能量召唤
	TEST_RATE_CONSUMER       = 3 //无双神将
	TEST_RATE_GENERAL        = 4 //限时神将
	TEST_RATE_DREAMLAND_GOLD = 5 //神格金币转盘
	TEST_RATE_DREAMLAND_GEM  = 6 //神格钻石转盘
	TEST_RATE_ASTROLOGY      = 7 //占星
	TEST_RATE_SYNTHETIC      = 8 //合成紫灵魂石
)

var backstagemgrsingleton *BackStageMgr = nil

//! public
func GetBackStageMgr() *BackStageMgr {
	if backstagemgrsingleton == nil {
		backstagemgrsingleton = new(BackStageMgr)
		backstagemgrsingleton.Locker = new(sync.RWMutex)
		backstagemgrsingleton.NoticeList = make([]*NoticeInfo, 0)
		backstagemgrsingleton.NoticeLocker = new(sync.RWMutex)
		backstagemgrsingleton.TempNoticeList = make([]*NoticeInfo, 0)
		backstagemgrsingleton.TempNocticeLocker = new(sync.RWMutex)
	}

	return backstagemgrsingleton
}

func (self *BackStageMgr) Init() {
	http.HandleFunc("/saveplayer", self.SavePlayer)
	http.HandleFunc("/gagplayer", self.GagPlayer)
	http.HandleFunc("/tickplayer", self.TickPlayer)
	http.HandleFunc("/ungagplayer", self.UnGagPlayer)
	http.HandleFunc("/sendmail", self.SendMail)
	http.HandleFunc("/senditem", self.SendItem)
	http.HandleFunc("/notice", self.Notice)
	http.HandleFunc("/reloadactivity", self.ReloadActivity)
	http.HandleFunc("/updateactivity", self.UpdateActivity)
	http.HandleFunc("/backstage", self.HttpResource)
	http.HandleFunc("/campfight", self.CampFight)
	http.HandleFunc("/serverstatus", self.ServerStatus)
	http.HandleFunc("/reload", self.Reload)
	http.HandleFunc("/recharge", self.Recharge)
	http.HandleFunc("/blockplayer", self.BlockPlayer)
	http.HandleFunc("/memcheck", self.MemCheck)
	http.HandleFunc("/checkunion", self.CheckUnion)
	http.HandleFunc("/queryplayer", self.QueryPlayer)
	http.HandleFunc("/removehero", self.RemoveHero)
	http.HandleFunc("/removeitem", self.AddItem)
	http.HandleFunc("/checkallluckshop", self.checkAllLuckShop)
	http.HandleFunc("/getBarrage", self.GetBarrage)
	http.HandleFunc("/addNotice", self.AddNotice)
	http.HandleFunc("/fitserverBefore", self.FitServerBefore) // 合服前操作
	http.HandleFunc("/fitserverAfter", self.FitServerAfter)   // 合服后操作
	http.HandleFunc("/shutdown", self.ShutDown)
	http.HandleFunc("/addattackers", self.AddAttacker)
	http.HandleFunc("/mineok", self.MineOk)
	http.HandleFunc("/gveok", self.GveOk)
	http.HandleFunc("/minestart", self.MineStart)
	http.HandleFunc("/gvestart", self.GveStart)
	http.HandleFunc("/gvefaketaken", self.FakeTaken)
	http.HandleFunc("/unionstart", self.UnionStart)
	http.HandleFunc("/unionfightstart", self.UnionFightStart)
	http.HandleFunc("/unionfightunittest", self.UnionFightUnionTest)
	http.HandleFunc("/unionfightA", self.UnionFightStartA)
	http.HandleFunc("/unionfightB", self.UnionFightStartB)
	http.HandleFunc("/show_union_plan", self.ShowUnionPlan)
	http.HandleFunc("/test_union_fight", self.TestUnionFight)
	http.HandleFunc("/test_rate_by_type", self.TestRateByType)
	http.HandleFunc("/addServerTime", self.AddServerTime)
	http.HandleFunc("/checkpass", self.CheckPass)
	http.HandleFunc("/changeunionnotice", self.ChangeUnionNotice)
	http.HandleFunc("/activitybossdelete", self.ActivityBossDelete)
	http.HandleFunc("/loadsensitiveword", self.LoadSensitiveWord)
	http.HandleFunc("/getplayerbydebug", self.GetPlayerByDebug)
	http.HandleFunc("/recharge_update", self.RechargeUpdate)
}

type BackStageMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type BackStageStatus struct {
	Code   int           `json:"code"`
	Status *ServerStatus "status"
	Msg    string        `json:"msg"`
}

type BackStageMail struct {
	Uid      []int64    `json:"uid"`
	Title    string     `json:"title"`
	Body     string     `json:"body"`
	Item     []PassItem `json:"item"`
	Sender   string     `json:"sender"` //! 发送者
	MinLevel int        `json:"minlevel"`
	MaxLevel int        `json:"maxlevel"`
}

type BackStageItem struct {
	Uid  int64      `json:"uid"`
	Item []PassItem `json:"item"`
}

type BackStageActivity struct {
	Acts []*Sql_ActivityMask `acts`
}

type JS_NoticeInfo struct {
	Content   string `json:"content"`
	Systime   string `json:"systime"`
	BeginTime string `json:"begin_time"`
	Num       string `json:"num"`
}

type JS_GetNoticeListMsg struct {
	NoticList []JS_NoticeInfo `json:"lists"`
}

type NoticeInfo struct {
	Id        int
	Interval  int    //! 每隔多久播一次
	BeginTime int64  //! 开始时间
	MaxNum    int    //! 最大次数
	Content   string //! 内容
	PlayerNum int    //! 已经播放的次数
	TicksNum  int    //! tick时间
}

func (self *BackStageMgr) CheckToken(token string) bool {
	h := md5.New()
	h.Write([]byte(GetServer().Con.AdminCode))
	md5str := fmt.Sprintf("%x", h.Sum(nil))

	return md5str == token
}

func (self *BackStageMgr) SendMail(w http.ResponseWriter, r *http.Request) {

	context := r.FormValue("context")
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	text, err := base64.URLEncoding.DecodeString(context)
	if err != nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	var mail BackStageMail
	err = json.Unmarshal(text, &mail)
	if err != nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "json解析错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	LogInfo("邮件信息：%v", mail)
	//! 个人邮件
	if len(mail.Uid) > 0 {
		for i := 0; i < len(mail.Uid); i++ {
			player := GetPlayerMgr().GetPlayer(mail.Uid[i], true)
			if player == nil {
				continue
			}
			player.GetModule("mail").(*ModMail).AddMail(MAIL_CAN_ALL_GET, 2, 0, mail.Title, mail.Body, mail.Sender, mail.Item, true, 0)
		}
	} else {
		//! 全局邮件
		GetMailMgr().AddMail(mail.Title, mail.Body, mail.Item, mail.Sender, mail.MinLevel, mail.MaxLevel)
	}
	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) SendItem(w http.ResponseWriter, r *http.Request) {
	context := r.FormValue("context")
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	text, err := base64.URLEncoding.DecodeString(context)
	if err != nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	var mail BackStageItem
	err = json.Unmarshal(text, &mail)
	if err != nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "json解析错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	person := GetPlayerMgr().GetPlayer(mail.Uid, true)
	if person == nil {
		var msg BackStageMsg
		msg.Code = 2
		msg.Msg = "找不到uid"
		w.Write(HF_JtoB(&msg))
		return
	}

	var msgitem S2C_MailAllItem
	msgitem.Cid = "mailallitem"
	for i := 0; i < len(mail.Item); i++ {
		person.AddObject(mail.Item[i].ItemID, mail.Item[i].Num, 4, mail.Item[i].ItemID, 0, "gm赠送")
		msgitem.Item = append(msgitem.Item, PassItem{ItemID: mail.Item[i].ItemID, Num: mail.Item[i].Num})
	}

	person.SendMsg("senditem", HF_JtoB(&msgitem))

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) SavePlayer(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	player := GetPlayerMgr().GetPlayer(uid, false)
	if player == nil {
		var msg BackStageMsg
		msg.Code = 0
		msg.Msg = "成功"
		w.Write(HF_JtoB(&msg))
		return
	}

	player.Save(false, false)

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) RemoveHero(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	heroid := HF_Atoi(r.FormValue("heroid"))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	player := GetPlayerMgr().GetPlayer(uid, true)
	if player == nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "操作失败，找不到该玩家"
		w.Write(HF_JtoB(&msg))
		return
	}

	if player.GetSession() != nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "操作失败，玩家在线"
		w.Write(HF_JtoB(&msg))

		return
	}

	if player.GetModule("hero").(*ModHero).RemoveHeroByKeyId(heroid) {
		player.Save(false, false)
		var msg BackStageMsg
		msg.Code = 0
		msg.Msg = "成功"
		w.Write(HF_JtoB(&msg))
	} else {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "删除失败，可能是找不到KEYID"
		w.Write(HF_JtoB(&msg))
	}
}

func (self *BackStageMgr) AddItem(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	itemid := HF_Atoi(r.FormValue("itemid"))
	itemnum := HF_Atoi(r.FormValue("itemnum"))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	player := GetPlayerMgr().GetPlayer(uid, true)
	if player == nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "操作失败，找不到该玩家"
		w.Write(HF_JtoB(&msg))
		return
	}

	if player.GetSession() != nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "操作失败，玩家在线"
		w.Write(HF_JtoB(&msg))

		return
	}

	player.AddObject(itemid, -itemnum, 0, 0, 0, "GM修改物品")
	player.Save(false, false)

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))

}

type SimplePlayerData struct {
	Uid        int64  `json:"uid"`        //! 玩家ID
	UName      string `json:"uname"`      //! 名字
	Channelid  string `json:"channelid"`  //渠道id
	HeroNum    int    `json:"heronum"`    //武将总数
	Level      int    `json:"level"`      //! 玩家等级
	Unionid    int    `json:"unionid"`    //! 军团id
	FightNum   int64  `json:"fightnum"`   //! 当前战斗力
	Vip        int    `json:"vip"`        //! Vip等级
	TotalMoney int    `json:"totalmoney"` //总充值数量

	Gem    int `json:"gem"`    //! 钻石
	Gold   int `json:"gold"`   //! 金币
	TiLi   int `json:"tili"`   //! 体力
	Power  int `json:"power"`  //! 军令
	RongYu int `json:"rongyu"` //!荣誉

	Beauty         int `json:"beauty"`         //! 美人总战斗力
	ArtTotalLevel  int `json:"arttotallevel"`  //! 神器总等级
	ArtTotalStar   int `json:"arttotalstar"`   //! 神器总星级
	ArmsTotalLevel int `json:"armstotallevel"` //! 统兵总等级
	TotalHeroColor int `json:"totalherocolor"` //! 统兵阶级总等级
}

func (self *BackStageMgr) QueryPlayer(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	player := GetPlayerMgr().GetPlayer(uid, true)
	if player == nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "失败"
		w.Write(HF_JtoB(&msg))
		return
	}

	//player.Save(false, false)

	var playerdata SimplePlayerData

	playerdata.Uid = player.Sql_UserBase.Uid
	playerdata.UName = player.Sql_UserBase.UName
	playerdata.Channelid = player.Account.Channelid
	playerdata.HeroNum = player.GetModule("hero").(*ModHero).GetHeroNum()
	playerdata.Level = player.Sql_UserBase.Level
	playerdata.Unionid = player.GetModule("union").(*ModUnion).Sql_UserUnionInfo.Unionid
	playerdata.FightNum = player.Sql_UserBase.Fight
	playerdata.Vip = player.Sql_UserBase.Vip
	playerdata.TotalMoney = player.GetModule("recharge").(*ModRecharge).Sql_UserRecharge.Money

	playerdata.Gem = player.Sql_UserBase.Gem
	playerdata.Gold = player.Sql_UserBase.Gold
	playerdata.TiLi = player.Sql_UserBase.TiLi
	//playerdata.Power = player.GetModule("city").(*ModCity).Sql_City.Power
	//playerdata.RongYu = player.GetObjectNum(91000018)
	playerdata.Beauty = player.GetModule("beauty").(*ModBeauty).CountBeautyLevel()

	playerdata.ArtTotalLevel = 0
	playerdata.ArtTotalStar = 0

	playerdata.ArmsTotalLevel = 0
	playerdata.TotalHeroColor = player.GetModule("hero").(*ModHero).GetHeroTotalColor()

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = HF_JtoA(playerdata)
	w.Write(HF_JtoB(&msg))
}

//! 重置活动状态
func (self *BackStageMgr) ReloadActivity(w http.ResponseWriter, r *http.Request) {
	op := int(HF_Atoi(r.FormValue("op")))
	//token := r.FormValue("token")
	/*
		if self.CheckToken(token) == false {
			var msg BackStageMsg
			msg.Code = 1
			msg.Msg = "参数错误"
			w.Write(HF_JtoB(&msg))
			return
		}
	*/
	startTime := TimeServer().Unix()
	LogInfo("Reload Activity Backstage Start:", startTime)
	if op == 1 {
		GetActivityMgr().ReloadActivity()
		GetActivityMgr().Save()
	}

	LogInfo("Reload Activity Backstage End:", TimeServer().Unix()-startTime)
	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

//! 重置活动状态
func (self *BackStageMgr) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	op := int(HF_Atoi(r.FormValue("op")))
	context := r.FormValue("context")
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	text, err := base64.URLEncoding.DecodeString(context)
	if err != nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	textstr := string(text)
	LogDebug("update activity mask:", textstr)
	var actmask Sql_ActivityMask
	err = json.Unmarshal(text, &actmask)
	if err != nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "json解析错误" + textstr
		w.Write(HF_JtoB(&msg))
		return
	}
	actmask.Decode()
	if op == 1 {
		GetActivityMgr().UpdateActivity(&actmask)

		//! 刷新活动状态
		GetActivityMgr().ReloadActivity()

		GetActivityMgr().Save()
	} else if op == 2 {
		GetActivityMgr().UpdateActivity(&actmask)
	}

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

//! 改变国战状态
func (self *BackStageMgr) CampFight(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) Reload(w http.ResponseWriter, r *http.Request) {
	status := int(HF_Atoi(r.FormValue("op")))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	if status == 1 {
		GetCsvMgr().Reload()
	} else if status == 2 {
		GetServer().ReloadConfig()
	}

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) Recharge(w http.ResponseWriter, r *http.Request) {
	recharge := int(HF_Atoi(r.FormValue("recharge")))
	uid := int64(HF_Atoi(r.FormValue("uid")))
	order_id := int(HF_Atoi(r.FormValue("order_id")))
	status := int(HF_Atoi(r.FormValue("status")))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	player := GetPlayerMgr().GetPlayer(uid, true)
	if player != nil {
		//player.GetModule("recharge").(*ModRecharge).Recharge(recharge)
		player.GetModule("recharge").(*ModRecharge).BackRecharge(order_id, status, recharge)

		var msg BackStageMsg
		msg.Code = 0
		msg.Msg = "成功"
		w.Write(HF_JtoB(&msg))
	}

}

//! 服务器状态
func (self *BackStageMgr) ServerStatus(w http.ResponseWriter, r *http.Request) {
	//status := int(HF_Atoi(r.FormValue("status")))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	status := GetServer().GetServerStatus()

	var msg BackStageStatus
	msg.Code = 0
	msg.Msg = "OK"
	msg.Status = status

	w.Write(HF_JtoB(&msg))
}

//! 服务器状态
func (self *BackStageMgr) MemCheck(w http.ResponseWriter, r *http.Request) {
	//status := int(HF_Atoi(r.FormValue("status")))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	fm, err := os.OpenFile("./mem.out", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(fm)
	fm.Close()

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "完成"
	w.Write(HF_JtoB(&msg))
}

type BarrageMsg struct {
	MsgLst []string `json:"msglst"`
}

//! 服务器状态
func (self *BackStageMgr) CheckUnion(w http.ResponseWriter, r *http.Request) {
	//status := int(HF_Atoi(r.FormValue("status")))
	token := r.FormValue("token")
	unionid := r.FormValue("unionid")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	GetUnionMgr().CheckMaster(HF_Atoi(unionid))

	var msg1 BarrageMsg
	msg1.MsgLst = GetServer().GetBarrage(100, 1)

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = HF_JtoA(&msg1)
	w.Write(HF_JtoB(&msg))
}

//! 踢人
func (self *BackStageMgr) TickPlayer(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	player := GetPlayerMgr().GetPlayer(uid, false)
	if player == nil {
		var msg BackStageMsg
		msg.Code = 0
		msg.Msg = "成功"
		w.Write(HF_JtoB(&msg))
		return
	}

	player.SafeClose()

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

//! 踢人
func (self *BackStageMgr) BlockPlayer(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	block := int64(HF_Atoi(r.FormValue("block")))
	blockday := HF_Atoi(r.FormValue("day"))
	blockreason := r.FormValue("reason")
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	player := GetPlayerMgr().GetPlayer(uid, true)
	if player == nil {

		var msg BackStageMsg
		msg.Code = 0
		msg.Msg = "没有用户"
		w.Write(HF_JtoB(&msg))
		return
	}

	if block == 0 {
		player.Sql_UserBase.IsBlock = 0
		player.Sql_UserBase.BlockDay = 0
	} else {
		player.Sql_UserBase.IsBlock = 1
		player.Sql_UserBase.BlockDay = blockday
	}
	player.Sql_UserBase.BlockTime = TimeServer().Unix()
	player.Sql_UserBase.BlockReason = blockreason
	player.Sql_UserBase.Update(true)
	player.SafeClose()

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) GagPlayer(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	blockday := HF_Atoi(r.FormValue("day"))
	blockreason := r.FormValue("reason")
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	player := GetPlayerMgr().GetPlayer(uid, true)
	if player == nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "找不到玩家"
		w.Write(HF_JtoB(&msg))
		return
	}

	if player.Sql_UserBase.IsBlock == 1 {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "该玩家已被封号"
		w.Write(HF_JtoB(&msg))
		return
	}

	player.Sql_UserBase.IsGag = 1
	player.Sql_UserBase.BlockDay = blockday
	player.Sql_UserBase.BlockTime = TimeServer().Unix()
	player.Sql_UserBase.BlockReason = blockreason
	player.Sql_UserBase.Update(true)

	//通知中心服清理聊天记录
	unionId := player.GetUnionId()
	GetMasterMgr().ChatRPC.GagPlayer(player, unionId)

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) UnGagPlayer(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	player := GetPlayerMgr().GetPlayer(uid, true)
	if player == nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "找不到玩家"
		w.Write(HF_JtoB(&msg))
		return
	}

	if player.Sql_UserBase.IsBlock == 1 {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "该玩家已被封号"
		w.Write(HF_JtoB(&msg))
		return
	}

	player.Sql_UserBase.IsGag = 0
	player.Sql_UserBase.BlockDay = 0
	player.Sql_UserBase.BlockTime = 0
	player.Sql_UserBase.BlockReason = ""
	player.Sql_UserBase.Update(true)

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) Notice(w http.ResponseWriter, r *http.Request) {
	context := r.FormValue("context")
	text, err := base64.URLEncoding.DecodeString(context)
	if err != nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	_type := HF_Atoi(r.FormValue("type"))

	if _type == 0 { //! 聊天公告
		GetServer().sendSysChat(string(text))
	} else { //! 跑马灯
		GetServer().Notice(string(text), 0, 0)
	}

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

// 合服前操作,token=7d68b043c22d770f333dc8a7326eadb8
// 发送购买排行奖励
func (self *BackStageMgr) FitServerBefore(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	self.setMerge(1)

	GetPlayerMgr().CheckAct()
	// 购买排行奖励
	GetTopBuildMgr().Refresh()
	// 踢掉所有玩家
	GetPlayerMgr().KickoutPlayers()

	// 合服
	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "合服前操作成功:" + self.getServInfo()
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) FitServerAfter(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	// 合服
	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "合服后操作成功:" + self.getServInfo()
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) getServInfo() string {
	return GetServer().Con.ServerName + fmt.Sprintf(", ServerId =%d", GetServer().Con.ServerId)
}

//! 全局检测幸运礼包
func (self *BackStageMgr) checkAllLuckShop(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	var luckshop San_LuckShop
	sql := fmt.Sprintf("select * from `san_luckshop`")
	res := GetServer().DBUser.GetAllData(sql, &luckshop)
	for i := 0; i < len(res); i++ {
		data := res[i].(*San_LuckShop)

		player := GetPlayerMgr().GetPlayer(data.Uid, true)

		if player != nil {
			player.GetModule("luckshop").(*ModLuckShop).OnGetOtherData()
			player.GetModule("luckshop").(*ModLuckShop).CheckBox()
		}

	}
	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

//! 检测幸运礼包
func (self *BackStageMgr) checkLuckShop(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	player := GetPlayerMgr().GetPlayer(uid, false)
	if player == nil {
		var msg BackStageMsg
		msg.Code = 0
		msg.Msg = "成功"
		w.Write(HF_JtoB(&msg))
		return
	}

	//player.GetModule("luckshop").(*ModLuckShop).checkAward()

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) HttpResource(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		LogError("参数错误")
		return
	}
	op := r.FormValue("op")
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	if op == "reloadcsv" {
		cmd := exec.Command("svn", "update")
		cmd.Dir = "./csv"
		err := cmd.Run()
		if err != nil {
			log.Println("reloadcsv:", err)
			w.Write([]byte("reload fail"))
			return
		}

		GetCsvMgr().Reload()

		w.Write([]byte("reload ok"))
	} else if op == "sessionnum" {
		w.Write([]byte(fmt.Sprintf("%d", GetSessionMgr().GetSessionNum())))
	}
}

func (self *BackStageMgr) checkGeneral(token string, actId, startId, endId int) error {
	if self.CheckToken(token) == false {
		return errors.New("token error")
	}

	// check act
	if actId == 0 {
		return errors.New("actId = 0")

	}

	// check start
	if startId == 0 {
		return errors.New("starId = 0")
	}

	// check end
	if endId == 0 {
		return errors.New("endId = 0")
	}

	return nil
}

//! 弹幕请求
func (self *BackStageMgr) GetBarrage(w http.ResponseWriter, r *http.Request) {

	token := r.FormValue("token")
	size := int(HF_Atoi(r.FormValue("size")))
	page := int(HF_Atoi(r.FormValue("page")))
	if size <= 0 {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "size <= 0 参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	if page <= 0 {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "page <= 0 参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	var msg1 BarrageMsgs
	msg1.Msglst = GetServer().GetBarrage(size, page)

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = HF_JtoA(&msg1)
	//msg.Msg = "123"
	w.Write(HF_JtoB(&msg))

}

func (self *BackStageMgr) RunGetNotice() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(time.Minute * 60)
	for {
		<-ticker.C
		//self.getNoticeList()
	}
	ticker.Stop()
}

func (self *BackStageMgr) getNoticeList() {
	url := "http://106.15.137.174/Alirefor/index.php?m=Doc&c=BlackWord&a=getLists&db_id=%d"
	dbId := GetServer().Con.DBCon.DBId
	str := fmt.Sprintf(url, dbId)
	LogInfo("get notice request:", str)
	res, err := http.Get(str)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		LogError("RunGetNotice http.Get error:", err.Error())
		return
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		LogError("RunGetNotice read body error:", err.Error())
		return
	}

	//LogDebug("result=", string(result))
	var msg JS_GetNoticeListMsg
	err = json.Unmarshal(result, &msg)
	if err != nil {
		LogError("RunGetNotice marshal error:", err.Error())
		return
	}

	//LogInfo(fmt.Sprintf("JS_GetNoticeListMsg=%+v", msg))
	// 检查数据的合法性
	for index := range msg.NoticList {
		value := msg.NoticList[index]
		if value.Content == "" {
			LogError("get notice content is empty")
			continue
		}

		// 间隔时间
		interlVal, err := strconv.Atoi(value.Systime)
		if err != nil {
			LogError("get notice systemtime error, systemtime:", value.Systime)
			continue
		}

		if interlVal <= 0 {
			LogError("get notice interlVal <= 0", interlVal)
			continue
		}

		// 播放次数
		playNum, err := strconv.Atoi(value.Num)
		if err != nil {
			LogError("get notice playNum error, systemtime:", value.Num)
			continue
		}

		// 设置时间格式错误
		t, err := time.ParseInLocation(DATEFORMAT, value.BeginTime, time.Local)
		if err != nil {
			LogError("get notice time error, beginTime:", value.BeginTime)
			continue
		}
		self.addNotice(interlVal, playNum, t.Unix(), value.Content)
	}

	self.NoticeLocker.RLock()
	for index := range self.NoticeList {
		LogInfo(fmt.Sprintf("self.NoticeList=%+v", self.NoticeList[index]))
	}
	self.NoticeLocker.RUnlock()
}

func (self *BackStageMgr) addNotice(interVal int, playNum int, time int64, content string) {
	self.NoticeLocker.Lock()
	msgNotice := &NoticeInfo{
		Interval:  interVal,
		MaxNum:    playNum,
		BeginTime: time,
		Content:   content,
	}

	foundIndex := -1
	for noticeIndex := range self.NoticeList {
		notice := self.NoticeList[noticeIndex]
		if notice.BeginTime == time {
			foundIndex = noticeIndex
			break
		}
	}

	if foundIndex == -1 {
		self.NoticeList = append(self.NoticeList, msgNotice)
	}

	self.NoticeLocker.Unlock()
}

func (self *BackStageMgr) RunCheckNotice() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		self.checkNotice()
		//self.simpleNotice()
		self.checkTempNotice()
	}
	ticker.Stop()
}

func (self *BackStageMgr) checkNotice() {
	self.NoticeLocker.Lock()
	defer self.NoticeLocker.Unlock()

	now := TimeServer().Unix()
	var invalidIndex []int
	for index := range self.NoticeList {
		value := self.NoticeList[index]
		if value.PlayerNum >= value.MaxNum {
			invalidIndex = append(invalidIndex, index)
			continue
		}

		//log.Println("beginTime:", time.Unix(value.BeginTime, 0).Format(DATEFORMAT), ", now :", TimeServer().Format(DATEFORMAT))
		if value.BeginTime > now {
			continue
		}

		if value.TicksNum%value.Interval == 0 {
			GetServer().Notice(self.NoticeList[index].Content, 0, 0)
			self.NoticeList[index].PlayerNum += 1
			//LogDebug("send temp notice!")
		}
		self.NoticeList[index].TicksNum += 1
	}

	for index := range invalidIndex {
		self.NoticeList = append(self.NoticeList[:index], self.NoticeList[index+1:]...)
	}

}

func (self *BackStageMgr) simpleNotice() {
	context := "亲，若在体验游戏中遇到问题，可先通过游戏主页-左下角设置-客服-自助查看问题哦。"
	t := TimeServer()
	hourOk := t.Hour() == 12 || t.Hour() == 18 || t.Hour() == 21
	minuteOk := t.Minute() == 0
	if hourOk && minuteOk && t.Second() < 3 {
		if self.NextSendTime < TimeServer().Unix() {
			self.NextSendTime = TimeServer().Unix() + int64(60)
			self.addNotice(10, 3, TimeServer().Unix()+3, context)
		}
	}
}

func (self *BackStageMgr) AddNotice(w http.ResponseWriter, r *http.Request) {
	context := r.FormValue("context")
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	num := r.FormValue("num")
	systime := r.FormValue("systime")
	id := r.FormValue("id")
	open := r.FormValue("op")
	if context == "" {
		LogError("get notice content is empty")
		return
	}

	// 间隔时间
	interlVal, err := strconv.Atoi(systime)
	if err != nil {
		LogError("get notice systemtime error, systime:", systime)
		self.makeError(w, r)
		return
	}

	if interlVal <= 0 {
		LogError("get notice interlVal <= 0", interlVal)
		self.makeError(w, r)
		return
	}

	// 播放次数
	playNum, err := strconv.Atoi(num)
	if err != nil {
		LogError("get notice playNum error, num:", num)
		self.makeError(w, r)
		return
	}

	// id
	noticeId, err := strconv.Atoi(id)
	if err != nil {
		LogError("get notice noticeId error, noticeId:", id)
		self.makeError(w, r)
		return
	}

	// open  1增加  2删除
	openAction, err := strconv.Atoi(open)
	if err != nil {
		LogError("get notice open error, op:", open)
		self.makeError(w, r)
		return
	}

	if openAction == 1 {
		self.addTempNotice(interlVal, playNum, TimeServer().Unix()+5, context, noticeId)
	} else if openAction == 2 {
		self.removeTempNotice(noticeId)
	} else {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "成功"
	w.Write(HF_JtoB(&msg))

}

func (self *BackStageMgr) makeError(w http.ResponseWriter, r *http.Request) {
	var msg BackStageMsg
	msg.Code = 1
	msg.Msg = "参数错误"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) checkTempNotice() {
	self.TempNocticeLocker.Lock()
	defer self.TempNocticeLocker.Unlock()

	var invalidIndex []int
	for index := range self.TempNoticeList {
		value := self.TempNoticeList[index]
		if TimeServer().Unix() < value.BeginTime {
			continue
		}
		if value.PlayerNum >= value.MaxNum {
			invalidIndex = append(invalidIndex, index)
			continue
		}

		if value.TicksNum%value.Interval == 0 {
			GetServer().Notice(self.TempNoticeList[index].Content, 0, 0)
			self.TempNoticeList[index].PlayerNum += 1
			//LogDebug("send temp notice!")
			LogDebug("send temp notice!", self.TempNoticeList[index].PlayerNum)
		}
		self.TempNoticeList[index].TicksNum += 1
	}

	// 清除多余的临时公告
	for index := range invalidIndex {
		self.TempNoticeList = append(self.TempNoticeList[:index], self.TempNoticeList[index+1:]...)
	}
}

// 增加临时文件
func (self *BackStageMgr) addTempNotice(interVal int, playNum int, time int64, content string, id int) {
	self.TempNocticeLocker.Lock()
	defer self.TempNocticeLocker.Unlock()

	isFind := false
	for _, v := range self.TempNoticeList {
		if v.Id == id {
			isFind = true
			break
		}
	}

	if !isFind {
		msgNotice := &NoticeInfo{
			Id:        id,
			Interval:  interVal,
			MaxNum:    playNum,
			BeginTime: time,
			Content:   content,
		}
		self.TempNoticeList = append(self.TempNoticeList, msgNotice)
	}
	/*
		foundIndex := -1
			for noticeIndex := range self.TempNoticeList {
				notice := self.TempNoticeList[noticeIndex]
				if notice.BeginTime == time {
					foundIndex = noticeIndex
					break
				}
			}

			if foundIndex == -1 {
				self.TempNoticeList = append(self.TempNoticeList, msgNotice)
			}
	*/
}

func (self *BackStageMgr) finishvisitall(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) buyVisitall(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) buyPower(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) doSmelt(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) buySmelt(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) UpdateMockCity(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) GemCostTest(w http.ResponseWriter, r *http.Request) {
	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "消费榜测试"
	w.Write(HF_JtoB(&msg))
}

// 模拟玩家进行钻石消费
func (self *BackStageMgr) doGemCost(w http.ResponseWriter, r *http.Request) {
	player := GetPlayerMgr().GetPlayer(18897, true)
	if player == nil {
		return
	}

	player.GetModule("actop").(*ModActop).OnGetOtherData()
	player.AddObject(DEFAULT_GEM, 100, 0, 0, 0, "钻石消费测试")
	player.AddObject(DEFAULT_GEM, -100, 0, 0, 0, "钻石消费测试")
	s := fmt.Sprintf("当前玩家钻石数量:%d", player.GetObjectNum(DEFAULT_GEM))
	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = s
	w.Write(HF_JtoB(&msg))
}

func (m *BackStageMgr) GetKey(key string) (string, error) {
	conn := GetServer().GetRedisConn()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return "0", fmt.Errorf("error getting key %s: %v", key, err)
	}
	return string(data), err
}

func (m *BackStageMgr) SetKey(key string, v string) error {
	conn := GetServer().GetRedisConn()
	defer conn.Close()

	value := []byte(v)
	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

// 设置合服状态
func (m *BackStageMgr) setMerge(n int) {
	m.SetKey("merge-flag", fmt.Sprintf("%d", n))
}

// 获取合服状态
func (m *BackStageMgr) getMerge() int {
	flag, err := m.GetKey("merge-flag")
	if err != nil {
		return 0
	}

	return HF_Atoi(flag)
}

// 发送军团信息以及全国信息
func (self *BackStageMgr) sendrecnotice1(w http.ResponseWriter, r *http.Request) {
	player := GetPlayerMgr().GetPlayer(18882, true)
	if player == nil {
		return
	}
	GetRedPacMgr().sendSystemChat(player.Sql_UserBase.UName, player.Sql_UserBase.Camp, 79500603)

}

func (self *BackStageMgr) sendrecnotice2(w http.ResponseWriter, r *http.Request) {
	player := GetPlayerMgr().GetPlayer(18882, true)
	if player == nil {
		return
	}

	player.GetModule("union").OnGetOtherData()
	party := GetUnionMgr().GetUnion(player.GetModule("union").(*ModUnion).Sql_UserUnionInfo.Unionid)
	if party != nil {
		GetRedPacMgr().sendUnionChat(party, player.Sql_UserBase.Camp, player.Sql_UserBase.UName, 91000002)
	}
}

func (self *BackStageMgr) setworldLevel(w http.ResponseWriter, r *http.Request) {
	level := r.FormValue("level")
	GetServer().Level = HF_Atoi(level)
}

func (self *BackStageMgr) AddServerTime(w http.ResponseWriter, r *http.Request) {
	timeParam := r.FormValue("time")
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	param := HF_AtoI64(timeParam)
	if param >= 2000000000 {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数过大"
		w.Write(HF_JtoB(&msg))
		return
	}

	if param != 0 && param < time.Now().Unix()+ServerTimeOffset {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "时光无法倒流"
		w.Write(HF_JtoB(&msg))
		return
	}

	if param == 0 {
		ServerTimeOffset += param
	} else {
		//防手滑
		if GetServer().GetConfig().LogCon.LogLevel == debugLevel {
			ServerTimeOffset = param - time.Now().Unix()
		} else {
			var msg BackStageMsg
			msg.Code = 1
			msg.Msg = "非debug模式下，限制修改功能"
			w.Write(HF_JtoB(&msg))
			return
		}
	}

	w.Write([]byte(fmt.Sprintf("当前系统时间:%s,unixtime:%d", TimeServer().String(), TimeServer().Unix())))
}

func (self *BackStageMgr) ActDial(w http.ResponseWriter, r *http.Request) {

}

// 测试军演外挂行为
func (self *BackStageMgr) TestMilitary(w http.ResponseWriter, r *http.Request) {

}

// 关闭服务器
func (self *BackStageMgr) ShutDown(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if token != "mbBJFN2m:IXG`}!i2vD)" {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	GetServer().Close()
}

// http://192.168.10.155:1241/addattackers?uid=100&times=10
func (self *BackStageMgr) AddAttacker(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) MineOk(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) GveOk(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) MineStart(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) GveStart(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) FakeTaken(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) UnionStart(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) UnionFightUnionTest(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) UnionFightStart(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) UnionFightStartA(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) UnionFightStartB(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) ShowUnionPlan(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) TestUnionFight(w http.ResponseWriter, r *http.Request) {

}

func (self *BackStageMgr) TestRateByType(w http.ResponseWriter, r *http.Request) {

	testType, err1 := strconv.Atoi(r.FormValue("type"))
	if err1 != nil {
		w.Write([]byte("testType:" + err1.Error()))
		return
	}

	count, err2 := strconv.Atoi(r.FormValue("count"))
	if err2 != nil {
		count = 1
	}

	subType, err4 := strconv.Atoi(r.FormValue("subtype"))
	if err4 != nil {
		subType = 1
	}

	uid, err5 := strconv.Atoi(r.FormValue("uid"))
	if err5 != nil {
		uid = 0
	}

	switch testType {
	case TEST_RATE_HIGH_CALL:

		if subType < 1 || subType > 16 {
			subType = 1
		}

		player := GetPlayerMgr().GetPlayer(int64(uid), true)
		if player == nil {
			w.Write([]byte(fmt.Sprintf("用户不存在")))
			return
		}
		realtype := 0
		findtimes := 0

		// 计算掉落物品 最多10000次循环
		if count >= 10000 {
			count = 10000
		}

		switch subType {
		case 1:
			w.Write([]byte(fmt.Sprintf("（上限10000）钻石物品单抽累积进行%d次：\n", count)))
			realtype = REALTYPE_GEM_ITEM
			findtimes = 1
		case 2:
			w.Write([]byte(fmt.Sprintf("（上限10000）钻石物品十连累积进行%d次：\n", count)))
			realtype = REALTYPE_GEM_ITEM_TEN
			findtimes = 10
		case 3:
			w.Write([]byte(fmt.Sprintf("（上限10000）钻石单抽累积进行%d次：\n", count)))
			realtype = REALTYPE_GEM
			findtimes = 1
		case 4:
			w.Write([]byte(fmt.Sprintf("（上限10000）钻石十连累积进行%d次：\n", count)))
			realtype = REALTYPE_GEM_TEN
			findtimes = 10
		case 5:
			w.Write([]byte(fmt.Sprintf("（上限10000）友情单抽累积进行%d次：\n", count)))
			realtype = REALTYPE_FRIEND
			findtimes = 1
		case 6:
			w.Write([]byte(fmt.Sprintf("（上限10000）友情十连累积进行%d次：\n", count)))
			realtype = REALTYPE_FRIEND_TEM
			findtimes = 10
		case 7, 8, 9, 10, 11, 12, 13, 14:
			camp := (subType - 5) / 2
			if subType%2 == 1 {
				w.Write([]byte(fmt.Sprintf("（上限10000）阵营%d单抽累计进行%d次：\n", camp, count)))
				realtype = camp + 2 + CAMP_DIS
				findtimes = 1
			} else {
				w.Write([]byte(fmt.Sprintf("（上限10000）阵营%d十连累计进行%d次：\n", camp, count)))
				realtype = camp + 2 + CAMP_DIS_TEN
				findtimes = 10
			}
		case 15:
			w.Write([]byte(fmt.Sprintf("（上限10000）自选单抽累计进行%d次：\n", count)))
			realtype = 8
			findtimes = 1
		case 16:
			w.Write([]byte(fmt.Sprintf("（上限10000）自选十连累计进行%d次：\n", count)))
			realtype = 108
			findtimes = 10
		}

		temp := 0
		itemRel := make(map[int]*Item, 0)
		countTen := make(map[int]int, 0)
		preThreeIndex := 0
		preThreeCount := 0

		for countTimes := 0; countTimes < count; countTimes++ {
			bag := make([]PassItem, 0)
			lst := GetCsvMgr().PubchesttotalGroup[realtype]

			certaintimes := lst[0].Certaintimesdroptype
			if certaintimes != 0 {
				itemid := 0
				if len(lst[0].Dropgroups) != 4 {
					return
				}
				certainitem := lst[0].Dropgroups[3]
				temp++
				findtimes = temp
				if temp >= certaintimes {
					itemid = certainitem
					temp = 0
				}

				if itemid != 0 {
					var dropitem PassItem
					dropitem.ItemID = itemid
					dropitem.Num = 1
					bag = append(bag, dropitem)
				}
			}
			item := player.GetModule("find").(*ModFind).LootItems(lst, realtype, findtimes, certaintimes, bag, 0)

			count := 0
			for _, v := range item {
				ccc := GetCsvMgr().ItemMap[v.ItemID]
				if ccc != nil && ccc.ItemCheck >= 4 {
					count++
				}
			}
			if count == 3 {
				if preThreeIndex > 0 && countTimes-preThreeIndex < 5 {
					preThreeCount++
				}
				preThreeIndex = countTimes
			}

			countTen[count]++
			AddItemMapHelper2(itemRel, item)
		}

		dropNum := 0
		for _, item := range itemRel {
			config := GetCsvMgr().GetItemConfig(item.ItemId)
			if config != nil {
				w.Write([]byte(fmt.Sprintf("[%d,%s,%d]\n", item.ItemId, config.ItemName, item.ItemNum)))
				dropNum += item.ItemNum
			} else {
				w.Write([]byte(fmt.Sprintf("异常配置ID：%d，个数：%d\n", item.ItemId, item.ItemNum)))
			}
		}
		w.Write([]byte(fmt.Sprintf("累计掉落个数：%d\n", dropNum)))
		/*
			for k, v := range countTen {
				w.Write([]byte(fmt.Sprintf("[%d紫次数：%d]\n", k, v)))
			}
			w.Write([]byte(fmt.Sprintf("[3紫之间间隔小于5的次数：%d]\n", preThreeCount)))
		*/
		return
	case TEST_RATE_CONSUMER:
		player := GetPlayerMgr().GetPlayer(int64(uid), true)
		if player == nil {
			w.Write([]byte(fmt.Sprintf("用户不存在")))
			return
		}
		realtype := 0
		findtimes := 0

		// 计算掉落物品 最多10000次循环
		if count >= 100000 {
			count = 100000
		}
		lootConfig, _ := GetCsvMgr().GetLootConfig()
		if lootConfig == nil {
			return
		}

		w.Write([]byte(fmt.Sprintf("（上限100000）次元召唤十连累积进行%d次：\n", count)))
		realtype = lootConfig.CallTenType
		findtimes = 10

		temp := 0
		itemRel := make(map[int]*Item, 0)
		countTen := make(map[int]int, 0)
		preThreeIndex := 0
		preThreeCount := 0

		for countTimes := 0; countTimes < count; countTimes++ {
			bag := make([]PassItem, 0)
			lst := GetCsvMgr().PubchesttotalGroup[realtype]

			certaintimes := lst[0].Certaintimesdroptype
			if certaintimes != 0 {
				itemid := 0
				if len(lst[0].Dropgroups) != 4 {
					return
				}
				certainitem := lst[0].Dropgroups[3]
				temp += 10
				findtimes = temp
				if temp >= certaintimes {
					itemid = certainitem
					temp = temp % certaintimes
				}

				if itemid != 0 {
					var dropitem PassItem
					dropitem.ItemID = itemid
					dropitem.Num = 1
					bag = append(bag, dropitem)
				}
			}
			item := player.GetModule("find").(*ModFind).LootItems(lst, realtype, findtimes, certaintimes, bag, 0)

			count := 0
			for _, v := range item {
				ccc := GetCsvMgr().ItemMap[v.ItemID]
				if v.ItemID == 11500701 {
					print(countTimes)
				}
				if ccc != nil && ccc.ItemCheck >= 4 {
					count++
				}
			}
			if count == 3 {
				if preThreeIndex > 0 && countTimes-preThreeIndex < 5 {
					preThreeCount++
				}
				preThreeIndex = countTimes
			}

			countTen[count]++
			AddItemMapHelper2(itemRel, item)
		}

		for _, item := range itemRel {
			config := GetCsvMgr().GetItemConfig(item.ItemId)
			if config != nil {
				w.Write([]byte(fmt.Sprintf("[%d,%s,%d]\n", item.ItemId, config.ItemName, item.ItemNum)))
			}
		}
		return
	case TEST_RATE_ASTROLOGY:
		player := GetPlayerMgr().GetPlayer(int64(uid), true)
		if player == nil {
			w.Write([]byte(fmt.Sprintf("用户不存在")))
			return
		}
		if !player.GetModule("find").(*ModFind).IsAstrologyHero() {
			w.Write([]byte(fmt.Sprintf("占星英雄未设置")))
			return
		}

		// 计算掉落物品 最多10000次循环
		if count >= 10000 {
			count = 10000
		}

		itemRel := make(map[int]*Item, 0)
		for i := 0; i < count; i++ {
			item := player.GetModule("find").(*ModFind).CalAstrologyDrop()
			if item.ItemID > 0 {
				AddItemMapHelper3(itemRel, item.ItemID, item.Num)
			}
		}

		for _, item := range itemRel {
			config := GetCsvMgr().GetItemConfig(item.ItemId)
			if config != nil {
				w.Write([]byte(fmt.Sprintf("[%d,%s,%d]\n", item.ItemId, config.ItemName, item.ItemNum)))
			}
		}
		return
	case TEST_RATE_SYNTHETIC:
		player := GetPlayerMgr().GetPlayer(int64(uid), true)
		if player == nil {
			w.Write([]byte(fmt.Sprintf("用户不存在")))
			return
		}
		// 计算掉落物品 最多10000次循环
		if count >= 10000 {
			count = 10000
		}

		itemRel := make(map[int]*Item, 0)
		for i := 0; i < count; i++ {
			item := player.GetModule("hero").(*ModHero).CalSyntheticDrop()
			if item.ItemID > 0 {
				AddItemMapHelper3(itemRel, item.ItemID, item.Num)
			}
		}

		for _, item := range itemRel {
			config := GetCsvMgr().GetItemConfig(item.ItemId)
			if config != nil {
				w.Write([]byte(fmt.Sprintf("[%d,%s,%d]\n", item.ItemId, config.ItemName, item.ItemNum)))
			}
		}
		score := player.GetModule("hero").(*ModHero).GetSyntheticScore()
		w.Write([]byte(fmt.Sprintf("剩余积分[%d]\n", score)))
		return
	default:
		w.Write([]byte(fmt.Sprintf("type找不到匹配的功能\n")))
		return
	}
}

// 刷新出转盘物品配置
func (self *BackStageMgr) RandLootItem(items []*DreamLandItem) (bool, *DreamLandItem, int) {
	nTotalChance := 0
	for _, obj := range items {
		objConfig := GetCsvMgr().DreamLandItemMap[obj.ID]

		if obj.Num > 0 {
			// 是否有英雄限制
			nTotalChance += objConfig.ExtractWeight
		}
	}

	if nTotalChance <= 0 {
		return false, nil, -1
	}

	// 随机配置
	nRandNum := HF_GetRandom(nTotalChance)
	total := 0
	// 根据权重返回配置
	for k, v := range items {
		if v.Num <= 0 {
			continue
		}

		objConfig, _ := GetCsvMgr().DreamLandItemMap[v.ID]

		total += objConfig.ExtractWeight
		if nRandNum < total {
			return true, v, k
		}
	}
	return false, nil, -1
}

func (self *BackStageMgr) CheckPass(w http.ResponseWriter, r *http.Request) {
	uid := int64(HF_Atoi(r.FormValue("uid")))
	token := r.FormValue("token")
	passid := HF_Atoi(r.FormValue("passid"))
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}
	player := GetPlayerMgr().GetPlayer(uid, true)
	if player == nil {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "找不到玩家"
		w.Write(HF_JtoB(&msg))
		return
	}

	if player.GetModule("pass").(*ModPass).GetPass(passid) == nil {
		player.GetModule("pass").(*ModPass).AddPass(passid, 3)
	}

	w.Write([]byte("处理成功"))
}

func (self *BackStageMgr) ChangeUnionNotice(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	unionid := HF_Atoi(r.FormValue("unionid"))
	content := r.FormValue("content")
	nType := HF_Atoi(r.FormValue("type"))
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	if nType == 1 {
		GetUnionMgr().GMAlertUnionNotice(unionid, content)
	} else {
		GetUnionMgr().GMAlertUnionBoard(unionid, content)
	}

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "处理成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) ActivityBossDelete(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	id := HF_Atoi(r.FormValue("id"))
	uid := HF_AtoI64(r.FormValue("uid"))
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	if !GetActivityBossMgr().DeleteUserRecord(id, uid) {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "删除失败"
		w.Write(HF_JtoB(&msg))
		return
	}

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "处理成功"
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) LoadSensitiveWord(w http.ResponseWriter, r *http.Request) {
	GetServer().LoadSensitiveWord()
}

func (self *BackStageMgr) GetPlayerByDebug(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = fmt.Sprintf("内存人数:%d", GetPlayerMgr().GetPlayerByDebug())
	w.Write(HF_JtoB(&msg))
}

func (self *BackStageMgr) removeTempNotice(id int) {
	self.TempNocticeLocker.Lock()
	defer self.TempNocticeLocker.Unlock()

	newList := make([]*NoticeInfo, 0)
	for _, v := range self.TempNoticeList {
		if v.Id == id {
			continue
		}
		newList = append(newList, v)
	}
	self.TempNoticeList = newList
}

func (self *BackStageMgr) RechargeUpdate(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if self.CheckToken(token) == false {
		var msg BackStageMsg
		msg.Code = 1
		msg.Msg = "参数错误"
		w.Write(HF_JtoB(&msg))
		return
	}

	GetRechargeMgr().NeedUpdate = true

	var msg BackStageMsg
	msg.Code = 0
	msg.Msg = "通知服务器获取订单OK"
	w.Write(HF_JtoB(&msg))
}
