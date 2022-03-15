package game

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// 活动涉及配置
type ServRank struct {
	Uid  int64 `json:"uid"`
	Rank int   `json:"rank"`
}

const (
	UPDATE_RANK_TIME = 300
)

// 神将活动单服务器配置
type GeneralMgr struct {
	Sql_TopRank   []*San_GeneralUser        //! 服务器排行信息, 主(存放1000名玩家排行榜,玩家获取前50显示)
	TopCache      []*Js_GeneralUser         //! 前50名跨服数据, 如果配置小于配置则取配置名次
	ServRank      []*ServRank               //! 50名外玩家信息
	Locker        *sync.RWMutex             //! 数据锁
	CacheLock     *sync.RWMutex             //! cache数据锁, 区分主, 从
	ServRankLock  *sync.RWMutex             //! 50名外玩家信息锁
	UpdateTime    int64                     //! 初始化请求时间
	TopNew        []*Js_GeneralUser         //! 排行前50，不定时同步
	RankNewMap    map[int64]*Js_GeneralUser //! 这个服务器里所有人的排名
	GeneralRecord []*GeneralRecord          //! 公告
}

//! 限时神将积分排行榜
type San_GeneralUser struct {
	Uid     int64  `json:"uid"`
	KeyId   int    `json:"keyid"`
	SvrId   int    `json:"svrid"`
	SvrName string `json:"svrname"`
	UName   string `json:"uname"`
	Level   int    `json:"level"`
	Vip     int    `json:"vip"`
	Icon    int    `json:"icon"`
	Point   int    `json:"point"`
	Rank    int    `json:"rank"`
	Time    int64  `json:"time"`

	DataUpdate
}

func (self *San_GeneralUser) toJsGeneralUser() *Js_GeneralUser {
	return &Js_GeneralUser{
		Uid:     self.Uid,
		SvrId:   self.SvrId,
		SvrName: self.SvrName,
		UName:   self.UName,
		Level:   self.Level,
		Vip:     self.Vip,
		Icon:    self.Icon,
		Point:   self.Point,
		Rank:    self.Rank,
		KeyId:   self.KeyId,
		Time:    self.Time,
	}
}

type Js_GeneralUser struct {
	Uid     int64  `json:"uid"`
	SvrId   int    `json:"svrid"`
	SvrName string `json:"svrname"`
	UName   string `json:"uname"`
	Level   int    `json:"level"`
	Vip     int    `json:"vip"`
	Icon    int    `json:"icon"`
	Point   int    `json:"point"`
	Rank    int    `json:"rank"`
	KeyId   int    `json:"keyid"`
	Time    int64  `json:"time"`
}

type GeneralRecord struct {
	Uid        int64  `json:"uid"`
	SvrId      int    `json:"svrid"`
	UName      string `json:"uname"`
	Item       int    `json:"item"`
	Num        int    `json:"num"`
	Time       int64  `json:"time"`
	RecordType int    `json:"recordtype"`
}

////////////////////////////////////////////////////////////////////////////////
type lstGeneralTop []*San_GeneralUser

func (s lstGeneralTop) Len() int      { return len(s) }
func (s lstGeneralTop) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstGeneralTop) Less(i, j int) bool {
	if s[i].Point > s[j].Point {
		return true
	}

	if s[i].Point < s[j].Point {
		return false
	}

	if s[i].Time < s[j].Time {
		return true
	}

	if s[i].Time > s[j].Time {
		return false
	}

	return false

}

var generalServ *GeneralMgr = nil

func GetGeneralMgr() *GeneralMgr {
	if generalServ == nil {
		generalServ = new(GeneralMgr)
		generalServ.Sql_TopRank = make([]*San_GeneralUser, 0)
		generalServ.TopCache = make([]*Js_GeneralUser, 0)
		generalServ.ServRank = make([]*ServRank, 0)
		generalServ.Locker = new(sync.RWMutex)
		generalServ.CacheLock = new(sync.RWMutex)
		generalServ.ServRankLock = new(sync.RWMutex)
		generalServ.TopNew = make([]*Js_GeneralUser, 0)
		generalServ.RankNewMap = make(map[int64]*Js_GeneralUser)
		generalServ.GeneralRecord = make([]*GeneralRecord, 0)
		generalServ.UpdateTime = TimeServer().Unix() + 30
	}

	return generalServ
}

//! 定时执行
func (self *GeneralMgr) OnTimer() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		if self.UpdateTime < TimeServer().Unix() {
			self.ReqGeneralRankListNew()
			self.UpdateTime = TimeServer().Unix() + UPDATE_RANK_TIME
			//LogDebug("开始向主服务器请求数据:", TimeServer().Format(DATEFORMAT))
		}
	}
	ticker.Stop()
}

// 获取活动掉落配置
func (self *GeneralMgr) GetLootConfig() (*TimeGeneralsConfig, int) {
	return GetCsvMgr().GetLootConfig()
}

func (self *GeneralMgr) GetRankAward(groupId int) []*TimeGeneralRank {
	return GetCsvMgr().GetRankAwardConf(groupId)
}

func (self *GeneralMgr) IsActOk() bool {
	activity := GetActivityMgr().GetActivity(ACT_GENERAL)
	if activity == nil {
		LogError("---- activity == nil")
		return false
	}

	info := activity.info
	startTime, err := self.GetActStartTime()
	if err != nil {
		LogError("startTime is error, start:", info.Start)
		return false
	}

	startSec := startTime
	now := TimeServer().Unix()
	if info.Continued == 0 {
		LogError("---- info.Continued == 0")
		return false
	}

	showTime := startSec + int64(info.Continued)
	res := now >= startSec && now <= showTime
	return res
}

func (self *GeneralMgr) GetActStartTime() (int64, error) {
	activity := GetActivityMgr().GetActivity(ACT_GENERAL)
	now := TimeServer().Unix()
	if activity == nil {
		return now, errors.New("activity is nil")
	}

	info := activity.info
	var startTime time.Time
	day, err := strconv.Atoi(info.Start)
	if err == nil {
		openTime, err := time.ParseInLocation(DATEFORMAT, GetServer().Con.OpenTime, time.Local)
		if err != nil {
			LogError("openTime is error, info.Start:", info.Start)
			return now, err
		}
		day -= 1
		startTime = openTime.Add(time.Duration(day*86400) * time.Second)
	} else {
		startTime, err = time.ParseInLocation(DATEFORMAT, info.Start, time.Local)
		if err != nil {
			LogError("startTime is error, info.Start:", info.Start)
			return now, err
		}
	}

	return startTime.Unix(), nil
}

func (self *GeneralMgr) GetActTime() (showTime, endTime int64, err error) {
	activity := GetActivityMgr().GetActivity(ACT_GENERAL)
	now := TimeServer().Unix()
	if activity == nil {
		//LogError("GetActTime ---- activity is not open.")
		return now, now, errors.New("activity is nil")
	}

	info := activity.info

	var startTime time.Time
	day, err := strconv.Atoi(info.Start)
	if err == nil {
		openTime, err := time.ParseInLocation(DATEFORMAT, GetServer().Con.OpenTime, time.Local)
		if err != nil {
			LogError("startTime is error, info.Start:", info.Start)
			return now, now, err
		}
		day -= 1
		startTime = openTime.Add(time.Duration(day*86400) * time.Second)
	} else {
		startTime, err = time.ParseInLocation(DATEFORMAT, info.Start, time.Local)
		if err != nil {
			LogError("startTime is error, info.Start:", info.Start)
			return now, now, err
		}
	}

	startSec := startTime.Unix()
	if now < startSec {
		return now, now, errors.New("activity is nil")
	}

	showTime = startSec + int64(info.Continued)
	endTime = startSec + int64(info.Continued) + int64(info.Show)
	if now > endTime {
		return now, now, errors.New("activity is nil")
	}

	err = nil
	return
}

func (self *GeneralMgr) GetCheckActTime() (endTime int64) {
	activity := GetActivityMgr().GetActivity(ACT_GENERAL)
	now := TimeServer().Unix()
	if activity == nil {
		return now
	}

	info := activity.info
	var startTime time.Time
	day, err := strconv.Atoi(info.Start)
	if err == nil {
		openTime, err := time.ParseInLocation(DATEFORMAT, GetServer().Con.OpenTime, time.Local)
		if err != nil {
			LogError("startTime is error, info.Start:", info.Start)
			return now
		}
		day -= 1
		startTime = openTime.Add(time.Duration(day*86400) * time.Second)
	} else {
		startTime, err = time.ParseInLocation(DATEFORMAT, info.Start, time.Local)
		if err != nil {
			LogError("startTime is error, info.Start:", info.Start)
			return now
		}
	}

	startSec := startTime.Unix()
	endTime = startSec + int64(info.Continued) + int64(info.Show)
	return
}

// 去掉不是当前活动keyId
// 如果有新的活动那就重新清理
func (self *GeneralMgr) CheckMasterKeyId(keyId int) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	// 得到当前的keyId
	needFlush := false
	for index := range self.Sql_TopRank {
		if self.Sql_TopRank[index] == nil {
			continue
		}

		if self.Sql_TopRank[index].KeyId != keyId {
			needFlush = true
			break
		}
	}

	if needFlush {
		// 将之前的数据立即存盘
		for _, value := range self.Sql_TopRank {
			value.UpdateEx("keyid", value.KeyId)
		}
		self.Sql_TopRank = make([]*San_GeneralUser, 0)
	}
}

func (self *GeneralMgr) getKeyId() int {
	lootConfig, step := self.GetLootConfig()
	if lootConfig == nil {
		return 0
	}

	return lootConfig.KeyId*1000 + step
}

func (self *GeneralMgr) UploadScoreNew(player *Player, top *Js_GeneralUser, records []*GeneralRecord) {
	self.CheckMasterKeyId(top.KeyId)

	res := GetMasterMgr().MatchGeneralUpdate(top, records)
	if res != nil {
		self.Locker.Lock()
		defer self.Locker.Unlock()

		self.TopNew = res.RankInfo
		self.GeneralRecord = res.GeneralRecord
		if res.SelfInfo != nil {
			self.RankNewMap[res.SelfInfo.Uid] = res.SelfInfo
		}

		for _, v := range self.TopNew {
			if v.SvrId == GetServer().Con.ServerId {
				self.RankNewMap[v.Uid] = v
			}
		}
	}
}

func (self *GeneralMgr) ReqGeneralRankListNew() {
	config, step := GetCsvMgr().GetLootConfig()
	if config == nil {
		return
	}

	keyId := config.KeyId*1000 + step

	res := GetMasterMgr().MatchGeneralGetAllRank(keyId, GetServer().Con.ServerId)
	if res != nil {
		self.Locker.Lock()
		defer self.Locker.Unlock()

		self.TopNew = res.RankInfo
		self.GeneralRecord = res.GeneralRecord
		for _, v := range self.TopNew {
			if v.SvrId == GetServer().Con.ServerId {
				self.RankNewMap[v.Uid] = v
			}
		}
	}
}

func (self *GeneralMgr) GetRankNew(player *Player) {

	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var msg S2C_GeneralRank
	msg.Cid = "getgeneralrank"
	if len(self.TopNew) > 50 {
		msg.Rank = self.TopNew[0:50]
	} else {
		msg.Rank = self.TopNew
	}
	_, ok := self.RankNewMap[player.Sql_UserBase.Uid]
	if ok {
		msg.MyRank = self.RankNewMap[player.Sql_UserBase.Uid].Rank
	}
	msg.GeneralRecord = self.GeneralRecord
	player.SendMsg(msg.Cid, HF_JtoB(msg))
}

func (self *GeneralMgr) getRankAwardNew(uid int64) (rank int, err error) {
	_, ok := self.RankNewMap[uid]
	if !ok {
		return 0, errors.New(GetCsvMgr().GetText("STR_MGR_GENERAL_YOU_ARE_NOT_ON_THE"))
	}

	return self.RankNewMap[uid].Rank, nil
}

func (self *GeneralMgr) CheckMailNew(rank, point int) (title string, text string, item []PassItem, err error) {
	pRankConfig := GetCsvMgr().GetGeneralAward(rank, point)
	if pRankConfig == nil {
		LogError("pRankConfig == nil, rank :=", rank, ", point := ", point)
		return "1", "1", []PassItem{}, errors.New("no config")
	}

	rankTitle := fmt.Sprintf(GetCsvMgr().GeneralRankMail.MailTitle, rank)
	//修复邮件异常 20190808 by zy
	//mailText1 := fmt.Sprintf(GetCsvMgr().GeneralRankMail.MailText1, rank, point, pRankConfig.NeetPoint)
	mailText1 := fmt.Sprintf(GetCsvMgr().GeneralRankMail.MailText1, rank, point, pRankConfig.NeetPoint)
	mailText2 := fmt.Sprintf(GetCsvMgr().GeneralRankMail.MailText2, rank, point)
	mailText3 := fmt.Sprintf(GetCsvMgr().GeneralRankMail.MailText3, rank, point)
	if pRankConfig.NeetPoint == 0 { // 固定奖励
		return rankTitle, mailText3, pRankConfig.GetAward(rank, point), nil
	} else {
		if point >= pRankConfig.NeetPoint { // 达到
			return rankTitle, mailText2, pRankConfig.GetAward(rank, point), nil
		} else { // 未达到
			return rankTitle, mailText1, pRankConfig.GetAward(rank, point), nil
		}
	}
}

// 获取Top1000当前玩家排行, 检查邮件时需要
func (self *GeneralMgr) GetUserRankNew(uid int64) (rank int) {
	_, ok := self.RankNewMap[uid]
	if !ok {
		return 0
	}

	return self.RankNewMap[uid].Rank
}
