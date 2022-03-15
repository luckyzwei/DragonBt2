package game

import (
	"encoding/json"
	"fmt"
	"sync"
	//"time"
)

const (
	LOTTERY_DRAW_RECORD_MAX = 50
	LUCKY_POOL_RECORD_MAX   = 20
)

// 离线信息管理
type FriendPowerInfo struct {
	FriendUid int64 `json:"frienduid"`
	Time      int64 `json:"time"` // 赠送时间
}

type TeamFight struct {
	TeamType int   `json:"teamtype"`
	Fight    int64 `json:"fight"`
}

//用来显示基本信息,上线后同步,这个消息如果要和客户端交互，请过滤
type PlayerBaseInfo struct {
	Name         string       `json:"uname"`
	Face         int          `json:"face"`
	IconId       int          `json:"iconid"`
	Portrait     int          `json:"portrait"`
	LastUpTime   int64        `json:"lastuptime"`  //最后更新时间
	HeroSetInfo  []*NewHero   `json:"herosetinfo"` //个人信息里的hero
	Stage        int          `json:"stage"`       //
	Server       int          `json:"server"`      //
	StageTime    int64        `json:"stagetime"`   //
	MaxFight     int64        `json:"maxfight"`    //
	NewHeroLv    int          `json:"newherolv"`   //
	Level        int          `json:"level"`
	HeroMaxLevel []int        `json:"heromaxlevel"`
	Signature    string       `json:"signature"`
	TeamFight    []*TeamFight `json:"teamfight"`
	MaxLevel     int          `json:"maxlevel"`
}

// 离线信息管理
type OfflineInfo struct {
	Uid              int64
	FriendPowerInfos string //友情点信息
	MailInfos        string // 邮件信息
	BaseInfo         string // 显示信息

	friendPowerInfos map[int64]*FriendPowerInfo //
	//mailInfos        map[int]*PvpFight        // 邮件
	baseInfo *PlayerBaseInfo

	DataUpdate
}

//活动战报管理
type OfflineRecordInfo struct {
	Id      int //
	Period  int
	Records string

	records []*PveFight

	DataUpdate
}

type LotteryDrawRecord struct {
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Times  int    `json:"times"` //第几次拿到的
	ItemId int    `json:"itemid"`
	Num    int    `json:"num"`
}

// 将数据库数据写入dataf
func (self *OfflineInfo) Decode() {
	err := json.Unmarshal([]byte(self.FriendPowerInfos), &self.friendPowerInfos)
	if err != nil {
		LogError("OfflineInfo Decode error:", err.Error())
	}

	err1 := json.Unmarshal([]byte(self.BaseInfo), &self.baseInfo)
	if err1 != nil {
		LogError("OfflineInfo Decode error:", err1.Error())
	}
}

func (self *OfflineRecordInfo) Decode() {
	err := json.Unmarshal([]byte(self.Records), &self.records)
	if err != nil {
		LogError("OfflineRecordInfo Decode error:", err.Error())
	}
}

// 将data数据写入数据库
func (self *OfflineInfo) Encode() {
	self.FriendPowerInfos = HF_JtoA(&self.friendPowerInfos)
	self.BaseInfo = HF_JtoA(&self.baseInfo)
}

func (self *OfflineRecordInfo) Encode() {
	self.Records = HF_JtoA(&self.records)
}

// 竞技场管理器
type OfflineInfoMgr struct {
	OfflineInfo   map[int64]*OfflineInfo     //
	OfflineRecord map[int]*OfflineRecordInfo //
	Lock          *sync.RWMutex              // 数据操作锁
	//有固定时效的数据，存redis就好
	LotteryDrawRecordLow  []*LotteryDrawRecord
	LotteryDrawRecordHigh []*LotteryDrawRecord
	LuckyFindRecord       []*LotteryDrawRecord
}

var offlineinfomgr *OfflineInfoMgr = nil

func GetOfflineInfoMgr() *OfflineInfoMgr {
	if offlineinfomgr == nil {
		offlineinfomgr = new(OfflineInfoMgr)
		offlineinfomgr.OfflineInfo = make(map[int64]*OfflineInfo)
		offlineinfomgr.OfflineRecord = make(map[int]*OfflineRecordInfo)
		offlineinfomgr.Lock = new(sync.RWMutex)
	}
	return offlineinfomgr
}

func (self *OfflineInfoMgr) GetData() {
	var info OfflineInfo
	tableName := "san_offlineinfo"
	sql := fmt.Sprintf("select * from `%s`", tableName)
	res := GetServer().DBUser.GetAllData(sql, &info)
	for i := 0; i < len(res); i++ {
		data := res[i].(*OfflineInfo)
		data.Init(tableName, data, false)
		data.Decode()
		self.OfflineInfo[data.Uid] = data
	}

	var infoRecord OfflineRecordInfo
	tableName = "san_offlinerecordinfo"
	sql = fmt.Sprintf("select * from `%s`", tableName)
	res = GetServer().DBUser.GetAllData(sql, &infoRecord)
	for i := 0; i < len(res); i++ {
		data := res[i].(*OfflineRecordInfo)
		data.Init(tableName, data, false)
		data.Decode()
		self.OfflineRecord[data.Id] = data
	}

	value, _, err := HGetRedisEx(`san_lotterydrawrecordlow`, 0, fmt.Sprintf("%d", 0))
	if err == nil {
		json.Unmarshal([]byte(value), &self.LotteryDrawRecordLow)
	}

	value1, _, err1 := HGetRedisEx(`san_lotterydrawrecordhigh`, 0, fmt.Sprintf("%d", 0))
	if err1 == nil {
		json.Unmarshal([]byte(value1), &self.LotteryDrawRecordHigh)
	}
}

func (self *OfflineInfoMgr) Save() {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	for _, value := range self.OfflineInfo {
		value.Encode()
		value.Update(true)
	}

	for _, value := range self.OfflineRecord {
		value.Encode()
		value.Update(true)
	}

	HMSetRedisEx("san_lotterydrawrecordlow", 0, &self.LotteryDrawRecordLow, DAY_SECS*10)
	HMSetRedisEx("san_lotterydrawrecordhigh", 0, &self.LotteryDrawRecordHigh, DAY_SECS*10)
}

//以前的单服离线系统，因为有部分数据的保存，还是需要的，下次改版解决下
func (self *OfflineInfoMgr) GetInfo(player *Player) *OfflineInfo {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.OfflineInfo[player.Sql_UserBase.Uid]
	if !ok {
		info = new(OfflineInfo)
		info.Uid = player.Sql_UserBase.Uid
		info.friendPowerInfos = make(map[int64]*FriendPowerInfo)
		info.Encode()
		tableName := "san_offlineinfo"
		InsertTable(tableName, info, 0, false)
		info.Init(tableName, info, false)
		self.OfflineInfo[info.Uid] = info
	}
	if info.baseInfo == nil {
		info.baseInfo = self.NewBaseInfo(player)
	}
	if len(info.baseInfo.HeroSetInfo) < MAX_FIGHT_POS {
		size := MAX_FIGHT_POS - len(info.baseInfo.HeroSetInfo)
		for i := 0; i < size; i++ {
			info.baseInfo.HeroSetInfo = append(info.baseInfo.HeroSetInfo, nil)
		}
	}
	return self.OfflineInfo[player.Sql_UserBase.Uid]
}

func (self *OfflineInfoMgr) GetRecord(id int, period int) []*PveFight {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.OfflineRecord[id]
	if !ok {
		info = new(OfflineRecordInfo)
		info.Id = id
		info.Period = period
		info.Encode()
		tableName := "san_offlinerecordinfo"
		InsertTable(tableName, info, 0, false)
		info.Init(tableName, info, false)
		self.OfflineRecord[info.Id] = info
	}

	if info.Period != period {
		info.Period = period
		info.records = make([]*PveFight, 0)
	}
	return info.records
}

func (self *OfflineInfoMgr) UpdateHeroSkin(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.OfflineInfo[player.Sql_UserBase.Uid]
	if !ok {
		return
	}
	if info.baseInfo != nil {
		for _, v := range info.baseInfo.HeroSetInfo {
			if v == nil {
				continue
			}
			hero := player.getHero(v.HeroKeyId)
			if nil != hero {
				if v.Skin != hero.Skin {
					v.Skin = hero.Skin
				}
			}
		}
	}

	return
}

func (self *OfflineInfoMgr) GetName(uid int64) string {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return ""
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return ""
	}
	return self.OfflineInfo[uid].baseInfo.Name
}

func (self *OfflineInfoMgr) GetIconId(uid int64) int {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return 0
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return 0
	}
	return self.OfflineInfo[uid].baseInfo.IconId
}

func (self *OfflineInfoMgr) GetPortrait(uid int64) int {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return 0
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return 0
	}
	return self.OfflineInfo[uid].baseInfo.Portrait
}

func (self *OfflineInfoMgr) GetBaseHero(player *Player) []*NewHero {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[player.Sql_UserBase.Uid]
	if !ok {
		return nil
	}
	if self.OfflineInfo[player.Sql_UserBase.Uid].baseInfo == nil {
		self.OfflineInfo[player.Sql_UserBase.Uid].baseInfo = self.NewBaseInfo(player)
	}
	return self.OfflineInfo[player.Sql_UserBase.Uid].baseInfo.HeroSetInfo
}

func (self *OfflineInfoMgr) IsBaseHero(uid int64, keyId int) bool {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	info, ok := self.OfflineInfo[uid]
	if !ok {
		return false
	}
	if info.baseInfo == nil {
		return false
	}
	for i := 0; i < len(info.baseInfo.HeroSetInfo); i++ {
		if info.baseInfo.HeroSetInfo[i] != nil && info.baseInfo.HeroSetInfo[i].HeroKeyId == keyId {
			return true
		}
	}
	return false
}

func (self *OfflineInfoMgr) DeleteBaseHero(uid int64, keyId int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	info, ok := self.OfflineInfo[uid]
	if !ok {
		return
	}
	if info.baseInfo == nil {
		return
	}
	for i := 0; i < len(info.baseInfo.HeroSetInfo); i++ {
		if info.baseInfo.HeroSetInfo[i] != nil && info.baseInfo.HeroSetInfo[i].HeroKeyId == keyId {
			info.baseInfo.HeroSetInfo[i] = nil
			return
		}
	}
}

func (self *OfflineInfoMgr) GetBaseInfo(uid int64) (int, int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return ONHOOK_INIT_LEVEL, 0
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return ONHOOK_INIT_LEVEL, 0
	}
	if self.OfflineInfo[uid].baseInfo.Stage == 0 {
		return ONHOOK_INIT_LEVEL, self.OfflineInfo[uid].baseInfo.Server
	} else {
		return self.OfflineInfo[uid].baseInfo.Stage, self.OfflineInfo[uid].baseInfo.Server
	}
}

func (self *OfflineInfoMgr) SetHeroMaxLevel(uid int64, level int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return
	}

	nlen := len(self.OfflineInfo[uid].baseInfo.HeroMaxLevel)
	if nlen > 0 {
		for i := 0; i < nlen; i++ {
			if level <= self.OfflineInfo[uid].baseInfo.HeroMaxLevel[i] {
				continue
			} else {
				self.OfflineInfo[uid].baseInfo.HeroMaxLevel = append(self.OfflineInfo[uid].baseInfo.HeroMaxLevel[0:i], append([]int{level}, self.OfflineInfo[uid].baseInfo.HeroMaxLevel[i:]...)...)
				break
			}
		}
	} else {
		if !self.CheckHeroMaxLevel(uid) {
			self.OfflineInfo[uid].baseInfo.HeroMaxLevel = append(self.OfflineInfo[uid].baseInfo.HeroMaxLevel, level)
		}
	}

	if len(self.OfflineInfo[uid].baseInfo.HeroMaxLevel) > 5 {
		self.OfflineInfo[uid].baseInfo.HeroMaxLevel = self.OfflineInfo[uid].baseInfo.HeroMaxLevel[:5]
	}
}

func (self *OfflineInfoMgr) CheckHeroMaxLevel(uid int64) bool {
	player := GetPlayerMgr().GetPlayer(uid, false)
	if player != nil {
		heros, _ := player.GetModule("hero").(*ModHero).GetBestFormat()
		nlen := len(heros)
		for i := 0; i < nlen; i++ {
			hero := player.getHero(heros[i])
			if hero != nil {
				self.OfflineInfo[uid].baseInfo.HeroMaxLevel = append(self.OfflineInfo[uid].baseInfo.HeroMaxLevel, hero.HeroLv)
			}
		}
		return true
	}
	return false
}

func (self *OfflineInfoMgr) GetHeroMaxLevel(uid int64) (int, int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return 0, 0
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return 0, 0
	}
	nlen := len(self.OfflineInfo[uid].baseInfo.HeroMaxLevel)
	if nlen <= 0 {
		self.CheckHeroMaxLevel(uid)
		nlen = len(self.OfflineInfo[uid].baseInfo.HeroMaxLevel)
	}
	allLevel := 0
	for _, v := range self.OfflineInfo[uid].baseInfo.HeroMaxLevel {
		allLevel += v
	}
	if nlen > 0 {
		return allLevel / nlen, self.OfflineInfo[uid].baseInfo.Server
	} else {
		return 0, self.OfflineInfo[uid].baseInfo.Server
	}
}

func (self *OfflineInfoMgr) NewBaseInfo(player *Player) *PlayerBaseInfo {
	rel := new(PlayerBaseInfo)
	rel.IconId = player.Sql_UserBase.IconId
	rel.Face = player.Sql_UserBase.Face
	rel.Portrait = player.Sql_UserBase.Portrait
	rel.Name = player.Sql_UserBase.UName
	rel.Server = GetServer().Con.ServerId
	heros, _ := player.GetModule("hero").(*ModHero).GetBestFormat()
	nlen := len(heros)
	for i := 0; i < nlen; i++ {
		hero := player.getHero(heros[i])
		if hero != nil {
			rel.HeroMaxLevel = append(rel.HeroMaxLevel, hero.HeroLv)
		}
	}
	if len(rel.HeroMaxLevel) > 5 {
		rel.HeroMaxLevel = rel.HeroMaxLevel[:5]
	}

	if len(rel.TeamFight) <= 0 {
		teams := []int{TEAMTYPE_ARENA_2, TEAMTYPE_ARENA_SPECIAL_4}
		for team := range teams {
			rel.TeamFight = append(rel.TeamFight, &TeamFight{team, 0})
		}
	}

	return rel
}

func (self *OfflineInfoMgr) SetFriendInfo(aim int64, uid int64) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[aim]
	if !ok {
		return
	}
	_, uidOk := self.OfflineInfo[aim].friendPowerInfos[uid]
	if uidOk {
		self.OfflineInfo[aim].friendPowerInfos[uid].Time = TimeServer().Unix()
	} else {
		info := new(FriendPowerInfo)
		info.FriendUid = uid
		info.Time = TimeServer().Unix()
		self.OfflineInfo[aim].friendPowerInfos[info.FriendUid] = info
	}
}

func (self *OfflineInfoMgr) SetPlayerSignature(uid int64, signature string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return
	}
	self.OfflineInfo[uid].baseInfo.Signature = signature
}

func (self *OfflineInfoMgr) SetPlayerOffTime(uid int64, time int64) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return
	}
	self.OfflineInfo[uid].baseInfo.LastUpTime = time
}

func (self *OfflineInfoMgr) SetPlayerStage(uid int64, stage int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return
	}
	self.OfflineInfo[uid].baseInfo.Stage = stage
	self.OfflineInfo[uid].baseInfo.StageTime = TimeServer().Unix()
}

func (self *OfflineInfoMgr) SetMaxFight(uid int64, fight int64) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return
	}
	self.OfflineInfo[uid].baseInfo.MaxFight = fight
}

func (self *OfflineInfoMgr) SetNewHeroLv(uid int64, lv int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return
	}
	self.OfflineInfo[uid].baseInfo.NewHeroLv = lv
}

func (self *OfflineInfoMgr) GetMaxFight(uid int64) int64 {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return 0
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return 0
	}
	return self.OfflineInfo[uid].baseInfo.MaxFight
}

func (self *OfflineInfoMgr) GetPlayerOffTime(uid int64) int64 {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return 0
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return 0
	}
	return self.OfflineInfo[uid].baseInfo.LastUpTime
}

func (self *OfflineInfoMgr) GetPlayerSignature(uid int64) string {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return ""
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return ""
	}
	return self.OfflineInfo[uid].baseInfo.Signature
}

func (self *OfflineInfoMgr) GetNewHeroLv(uid int64) int {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return 1
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return 1
	}
	return self.OfflineInfo[uid].baseInfo.NewHeroLv
}

func (self *OfflineInfoMgr) Rename(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[player.Sql_UserBase.Uid]
	if !ok {
		return
	}
	if self.OfflineInfo[player.Sql_UserBase.Uid].baseInfo == nil {
		return
	}
	self.OfflineInfo[player.Sql_UserBase.Uid].baseInfo.Name = player.Sql_UserBase.UName
}

func (self *OfflineInfoMgr) ReIconId(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[player.Sql_UserBase.Uid]
	if !ok {
		return
	}
	if self.OfflineInfo[player.Sql_UserBase.Uid].baseInfo == nil {
		return
	}
	self.OfflineInfo[player.Sql_UserBase.Uid].baseInfo.IconId = player.Sql_UserBase.IconId
}

func (self *OfflineInfoMgr) RePortrait(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[player.Sql_UserBase.Uid]
	if !ok {
		return
	}
	if self.OfflineInfo[player.Sql_UserBase.Uid].baseInfo == nil {
		return
	}
	self.OfflineInfo[player.Sql_UserBase.Uid].baseInfo.Portrait = player.Sql_UserBase.Portrait
}

func (self *OfflineInfoMgr) SetArenaFight(teamType int, uid int64, fight int64) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.OfflineInfo[uid]
	if !ok {
		return
	}

	if info.baseInfo == nil {
		return
	}

	teams := []int{TEAMTYPE_ARENA_2, TEAMTYPE_ARENA_SPECIAL_4}
	for _, team := range teams {
		find := false

		for _, value := range info.baseInfo.TeamFight {
			if value.TeamType == team {
				find = true
				break
			}
		}

		if !find {
			info.baseInfo.TeamFight = append(info.baseInfo.TeamFight, &TeamFight{team, 0})
		}
	}

	for _, value := range info.baseInfo.TeamFight {
		if value.TeamType == teamType {
			value.Fight = fight
			break
		}
	}
}

func (self *OfflineInfoMgr) GetTeamFight(uid int64, teamType int) int64 {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.OfflineInfo[uid]
	if !ok {
		return 0
	}
	if self.OfflineInfo[uid].baseInfo == nil {
		return 0
	}

	for _, value := range self.OfflineInfo[uid].baseInfo.TeamFight {
		if value.TeamType == teamType {
			return value.Fight
		}
	}
	return 0
}

//法阵核心等级
func (self *OfflineInfoMgr) UpdateMaxLevel(player *Player, level int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.OfflineInfo[player.Sql_UserBase.Uid]
	if !ok {
		return
	}
	if info.baseInfo == nil {
		info.baseInfo = self.NewBaseInfo(player)
	}

	if info.baseInfo.MaxLevel <= 0 {
		info.baseInfo.MaxLevel = 1
	}

	if info.baseInfo.MaxLevel < level {
		info.baseInfo.MaxLevel = level
	}

	return
}

func (self *OfflineInfoMgr) GetMaxLevel(uid int64) int {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.OfflineInfo[uid]
	if !ok {
		return 1
	}
	if info.baseInfo == nil {
		return 1
	}
	if info.baseInfo.MaxLevel <= 0 {
		info.baseInfo.MaxLevel = 1
	}

	return info.baseInfo.MaxLevel
}

func (self *OfflineInfoMgr) AddRecord(id int, period int, fight *PveFight) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.OfflineRecord[id]
	if !ok {
		info = new(OfflineRecordInfo)
		info.Id = id
		info.Period = period
		info.Encode()
		tableName := "san_offlinerecordinfo"
		InsertTable(tableName, info, 0, false)
		info.Init(tableName, info, false)
		self.OfflineRecord[info.Id] = info
	}
	if info.Period != period {
		info.records = make([]*PveFight, 0)
	}
	if len(info.records) > 20 {
		info.records = info.records[1:]
	}
	info.records = append(info.records, fight)
}

func (self *OfflineInfoMgr) AddLotteryDrawRecord(record *LotteryDrawRecord, notice int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	if notice == 1 {
		if len(self.LotteryDrawRecordLow) > LOTTERY_DRAW_RECORD_MAX {
			self.LotteryDrawRecordLow = self.LotteryDrawRecordLow[1:]
		}
		self.LotteryDrawRecordLow = append(self.LotteryDrawRecordLow, record)
	} else if notice == 2 {
		if len(self.LotteryDrawRecordHigh) > LOTTERY_DRAW_RECORD_MAX {
			self.LotteryDrawRecordHigh = self.LotteryDrawRecordHigh[1:]
		}
		self.LotteryDrawRecordHigh = append(self.LotteryDrawRecordHigh, record)
	}
}

func (self *OfflineInfoMgr) GetLotteryDrawRecord() ([]*LotteryDrawRecord, []*LotteryDrawRecord) {
	self.Lock.RLock()
	defer self.Lock.RUnlock()
	return self.LotteryDrawRecordLow, self.LotteryDrawRecordHigh
}

func (self *OfflineInfoMgr) AddLuckyFindRecord(record *LotteryDrawRecord) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	if len(self.LuckyFindRecord) > LUCKY_POOL_RECORD_MAX {
		self.LuckyFindRecord = self.LuckyFindRecord[1:]
	}
	self.LuckyFindRecord = append(self.LuckyFindRecord, record)
}

func (self *OfflineInfoMgr) GetLuckyFindRecord() []*LotteryDrawRecord {
	self.Lock.RLock()
	defer self.Lock.RUnlock()
	return self.LuckyFindRecord
}
