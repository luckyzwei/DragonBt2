package game

import (
	"encoding/json"
	"fmt"
	"log"
	"master/utils"
	"runtime/debug"
	"sync"
	"time"

	//"time"
)

type PassRecord struct {
	UID        int64                  `json:"uid"`
	Level      int                    `json:"level"`
	Icon       int                    `json:"icon"`
	Portrait   int                    `json:"portrait"`
	Name       string                 `json:"name"`
	Fight      int64                  `json:"fight"`
	Camp       int                    `json:"camp"`
	Vip        int                    `json:"vip"`
	PassTime   int64                  `json:"pass_time"`
	RecordInfo *San_TowerPlayerRecord `json:"battlerecord"`
}

type PassHero struct {
	HeroId     int `json:"heroid"`
	Level      int `json:"level"`
	Star       int `json:"star"`
	Color      int `json:"color"`
	MainTalent int `json:"maintalent"`
}

type San_PassRecord struct {
	KeyID       int64  `json:"keyid"`       // key值
	Name        string `json:"name"`        // 玩家名字
	Uid         int64  `json:"uid"`         // 玩家uid
	Icon        int    `json:"icon"`        // 头像
	Portrait    int    `json:"portrait"`    // 头像框
	Level       int    `json:"level"`       // 等级
	PlayerFight int64  `json:"playerfight"` // 玩家战力
	BattleFight int64  `json:"battlefight"` // 战斗参与的战力
	Time        int64  `json:"time"`        // 时间
}

type PassRecords struct {
	PassId     int    `json:"pass_id"`
	FirstTeam  string `json:"firstteam"`  //! 首次通关的玩家
	LowTeam    string `json:"lowteam"`    //! 最低战力三星通关玩家
	RecentTeam string `json:"recentteam"` //! 最近通关队伍

	First      *San_PassRecord   `json:"first"`  // 最早通关
	Low        *San_PassRecord   `json:"low"`    // 最低战力通关
	Recent     []*San_PassRecord `json:"recent"` // 最近通关
	InfoLocker *sync.RWMutex     `json:"-"`      // 互斥锁

	DataUpdate
}

type PassRecordMgr struct {
	Records    map[int]*PassRecords
	InfoLocker *sync.RWMutex
	migrateOK  bool
}

var passRecordMgr *PassRecordMgr

func GetPassRecordMgr() *PassRecordMgr {
	if passRecordMgr == nil {
		passRecordMgr = new(PassRecordMgr)
		passRecordMgr.Records = make(map[int]*PassRecords)
		passRecordMgr.InfoLocker = new(sync.RWMutex)
		passRecordMgr.migrateOK = false
	}
	return passRecordMgr
}

func (self *PassRecordMgr) GetTableName() string {
	return "san_passrecords"
}

// 创建国家任务
func (self *PassRecordMgr) GetData() {
	var pass PassRecords
	sql := fmt.Sprintf("select * from `%s`", self.GetTableName())
	res := GetServer().DBUser.GetAllData(sql, &pass)
	for i := 0; i < len(res); i++ {
		data := res[i].(*PassRecords)
		data.Init(self.GetTableName(), data, false)
		data.Decode()
		data.InfoLocker = new(sync.RWMutex)
		if data.Recent == nil {
			data.Recent = make([]*San_PassRecord, 0)
		}
		self.Records[data.PassId] = data
	}

	// 初始化关卡记录信息
	self.InitPass()
}

// 开启迁移数据协程，竞技场
func (self *PassRecordMgr) RunMigratePassRecord() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	infoName := "san_passbattleinfo"     //! info 表
	recordName := "san_passbattlerecord" //! record 表
	recordType := 6
	tableName := "san_battlerecord"

	len, err := GetRedisMgr().HLen(infoName)
	if err != nil {
		return
	}
	LogInfo("迁移数据：", infoName, len)

	migOK, err := GetRedisMgr().Exists(infoName + "_migrateOKNew")
	if err == nil && migOK == true {
		self.migrateOK = true
	}
	count := 0
	cursor := int64(0)
	for {
		if self.migrateOK == true {
			break
		}

		//! 迁移数据
		cursor1, num := MigrateDataOne(infoName, recordName, tableName, recordType, cursor)
		count += num
		cursor = cursor1

		if count >= len {
			GetRedisMgr().Set(infoName+"_migrateOKNew", "1")
			break
		}

		//! 延迟1ms
		time.Sleep(time.Millisecond)
	}

	LogInfo(infoName, "迁移数据OK")

}

func (self *PassRecordMgr) GetPass(id int) *PassRecords {
	self.InfoLocker.RLock()
	defer self.InfoLocker.RUnlock()

	value, ok := self.Records[id]
	if ok {
		return value
	}
	return nil
}

func (self *LevelConfig) needRecord() bool {
	if self.MainType != 1 { //! 只记录普通副本
		return false
	}

	if self.LevelType != 2 && self.LevelType != 3 {
		return false
	}

	return true
}

func (self *PassRecordMgr) InitPass() {
	data := GetCsvMgr().LevelConfigMap
	for _, v := range data {
		if !v.needRecord() {
			continue
		}

		pass := self.GetPass(v.LevelId)
		if pass != nil { //! 开服默认所有城市都不开战
			continue
		}

		pass = new(PassRecords)
		pass.PassId = v.LevelId

		pass.First = nil
		pass.Low = nil
		pass.Recent = make([]*San_PassRecord, 0)
		pass.InfoLocker = new(sync.RWMutex)

		pass.Encode()
		InsertTable(self.GetTableName(), pass, 0, false)
		self.Records[pass.PassId] = pass
	}
}

func (self *PassRecords) Decode() { //! 将数据库数据写入data
	err := json.Unmarshal([]byte(self.FirstTeam), &self.First)
	if err != nil {
		LogError(err.Error())
	}
	err = json.Unmarshal([]byte(self.LowTeam), &self.Low)
	if err != nil {
		LogError(err.Error())
	}
	err = json.Unmarshal([]byte(self.RecentTeam), &self.Recent)
	if err != nil {
		LogError(err.Error())
	}
}

func (self *PassRecords) Encode() { //! 将data数据写入数据库
	self.FirstTeam = HF_JtoA(&self.First)
	self.LowTeam = HF_JtoA(&self.Low)
	self.RecentTeam = HF_JtoA(&self.Recent)
}

func (self *PassRecords) NewPassRecord(player *Player, dataRecord *San_PassRecord, info *BattleInfo, recordBattle *BattleRecord) *San_PassRecord {
	HMSetRedisEx("san_passbattleinfo", info.Id, info, utils.HOUR_SECS*12)
	HMSetRedisEx("san_passbattlerecord", recordBattle.Id, recordBattle, utils.HOUR_SECS*12)

	GetServer().DBUser.SaveRecord(BATTLE_TYPE_NORMAL, info, recordBattle)
	return dataRecord
}

// 最先通关
func (self *PassRecords) AddFirst(player *Player, dataRecord *San_PassRecord, info *BattleInfo, record *BattleRecord) {
	if self == nil {
		return
	}

	if self.First == nil {
		self.First = self.NewPassRecord(player, dataRecord, info, record)
	}
}

// 最低战力通关
func (self *PassRecords) AddLowFight(player *Player, dataRecord *San_PassRecord, info *BattleInfo, record *BattleRecord) {
	if self == nil {
		return
	}

	if self.Low == nil {
		self.Low = self.NewPassRecord(player, dataRecord, info, record)
	} else {
		if self.Low.BattleFight > dataRecord.BattleFight {
			self.Low = self.NewPassRecord(player, dataRecord, info, record)
		}
	}
}

// 最近三个通关
func (self *PassRecords) AddRecentPass(player *Player, dataRecord *San_PassRecord, info *BattleInfo, record *BattleRecord) {
	if self == nil {
		return
	}

	if len(self.Recent) < 5 {
		self.Recent = append(self.Recent, self.NewPassRecord(player, dataRecord, info, record))
	} else {
		self.Recent = append(self.Recent[:0], self.Recent[1:]...)
		self.Recent = append(self.Recent, self.NewPassRecord(player, dataRecord, info, record))
	}
}

// 检查记录
func (self *PassRecords) CheckRecord(passId int, star int, player *Player, dataRecord *San_PassRecord, info *BattleInfo, record *BattleRecord) {
	if self.InfoLocker == nil {
		LogError("CheckRecord: self.InfoLocker == nil")
		return
	}

	self.InfoLocker.Lock()
	defer self.InfoLocker.Unlock()
	self.AddLowFight(player, dataRecord, info, record)
	if star >= 1 {
		self.AddFirst(player, dataRecord, info, record)
	}
	self.AddRecentPass(player, dataRecord, info, record)
}

func (self *PassRecordMgr) GetRecord(passId int) *PassRecords {
	self.InfoLocker.RLock()
	defer self.InfoLocker.RUnlock()
	pRecord, ok := self.Records[passId]
	if !ok {
		return nil
	}
	return pRecord
}

func (self *PassRecordMgr) AddRecord(passId int, pRecord *PassRecords) {
	self.InfoLocker.Lock()
	defer self.InfoLocker.Unlock()
	self.Records[passId] = pRecord
}

// 检查记录
func (self *PassRecordMgr) CheckRecord(passId int, star int, player *Player, dataRecord *San_PassRecord, info *BattleInfo, record *BattleRecord) {
	data := GetCsvMgr().LevelConfigMap
	config, ok := data[passId]
	if !ok {
		return
	}

	if !config.needRecord() {
		return
	}

	pRecord := self.GetRecord(passId)
	if pRecord == nil {
		return
	}

	pRecord.CheckRecord(passId, star, player, dataRecord, info, record)
}

func (self *PassRecordMgr) Rename(player *Player) {
	LogDebug("修改名字：", player.Sql_UserBase.UName)
	data := GetCsvMgr().LevelConfigMap
	for _, v := range data {
		if !v.needRecord() {
			continue
		}

		pass := self.GetPass(v.LevelId)
		if pass == nil {
			continue
		}

		LogDebug("修改名字：", v.LevelId, player.Sql_UserBase.UName)
		pass.AlterName(player)
	}
}

func (self *PassRecords) AlterName(player *Player) {
	if self.InfoLocker == nil {
		LogError("self.InfoLocker == nil")
		return
	}
	self.InfoLocker.Lock()
	defer self.InfoLocker.Unlock()
	if self.First != nil && self.First.Uid == player.GetUid() {
		self.First.Name = player.Sql_UserBase.UName
	}

	if self.Low != nil && self.Low.Uid == player.GetUid() {
		self.Low.Name = player.Sql_UserBase.UName
	}

	for i := 0; i < len(self.Recent); i++ {
		if self.Recent[i].Uid == player.GetUid() {
			self.Recent[i].Name = player.Sql_UserBase.UName
		}
	}
}

func (self *PassRecordMgr) GetRecords() []*PassRecords {
	self.InfoLocker.RLock()
	defer self.InfoLocker.RUnlock()
	var res []*PassRecords
	for _, v := range self.Records {
		res = append(res, v)
	}
	return res
}

func (self *PassRecordMgr) Save() {
	records := self.GetRecords()
	for _, value := range records {
		value.Encode()
		value.Update(true)
	}
}

func (self *PassRecordMgr) GetPassRecord(player *Player, passId int) {
	pass := self.GetPass(passId)
	if pass == nil {
		player.SendErr(GetCsvMgr().GetText("STR_MGR_PASS_RECORD_NO_BARRIER_EXISTS"))
		return
	}

	var msg S2C_PassRecord
	msg.Cid = "passrecord"
	msg.First = pass.First
	msg.Low = pass.Low
	msg.Recent = pass.Recent
	msg.Passid = passId
	player.SendMsg("passrecord", HF_JtoB(&msg))
}
