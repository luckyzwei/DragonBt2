package tower

import (
	"encoding/json"
	"fmt"
	"game"
	//"github.com/garyburd/redigo/redis"
	"log"
	"master/center/match"
	"master/core"
	"master/db"
	"master/utils"
	"runtime/debug"
	"sort"
	"sync"
	"time"
)

type Js_TowerFightRecord struct {
	Key         int64  `json:"key"`      // key值
	Name        string `json:"name"`     // 玩家名字
	Uid         int64  `json:"uid"`      // 玩家uid
	Icon        int    `json:"icon"`     // 头像
	Portrait    int    `json:"portrait"` //
	Level       int    `json:"level"`
	PlayerFight int64  `json:"playerfight"` // 玩家战力
	BattleFight int64  `json:"battlefight"` // 战斗参与的战力
	Time        int64  `json:"time"`        // 时间
}

type ArmyInfo struct {
	Uid      int64  `json:"uid"`      //! uid
	Uname    string `json:"uname"`    //! 名字
	Iconid   int    `json:"iconid"`   //! 头像
	SelfKey  int    `json:"selfkey"`  //! 张三
	Face     int    `json:"face"`     //! 李四
	Portrait int    `json:"portrait"` //! 头像框
	Pos      int    `json:"pos"`      //! 队伍位置
}

type WeakenInfo struct {
	Att   int     `json:"att"`   //!
	Value float64 `json:"value"` //!
}

//! 战报
type BattleRecord struct {
	Id        int64                 `json:"id"`            //! 战报ID
	Type      int                   `json:"type"`          //! 战报类型
	Side      int                   `json:"side"`          // 自己是 1 进攻方 0 防守方
	Result    int                   `json:"attack_result"` // 0成功 其他失败
	Level     int                   `json:"level"`         // 之前名次
	LevelID   int                   `json:"levelid"`       //! 关卡id
	Time      int64                 `json:"time"`          // 发生的时间
	RandNum   int64                 `json:"rand_num"`      // 随机数
	Weaken    []*WeakenInfo         `json:"weaken"`        // 压制
	FightInfo [2]*game.JS_FightInfo `json:"fight_info"`    // 双方数据 一个是进攻 第二个是防守
}

type BattleHeroInfo struct {
	HeroID      int       `json:"heroid"`     //! 英雄id
	HeroLv      int       `json:"herolv"`     //! 英雄等级
	HeroStar    int       `json:"herostar"`   //! 英雄星级
	HeroSkin    int       `json:"skin"`       //! 英雄皮肤
	Hp          int       `json:"hp"`         // hp
	Energy      int       `json:"rage"`       // 怒气
	Damage      int64     `json:"damage"`     //! 伤害
	TakeDamage  int64     `json:"takedamage"` //! 承受伤害
	Healing     int64     `json:"healing"`    //! 治疗
	ArmyInfo    *ArmyInfo `json:"ownplayer"`
	ExclusiveLv int       `json:"exclusivelv"` //! 专属等级
	UseSkill    []int     `json:"skilltime"`   //! 使用的技能 pve专用  竞技场不使用
}

type BattleUserInfo struct {
	Uid       int64             `json:"uid"`        //! 名字
	Name      string            `json:"name"`       //! 名字
	Icon      int               `json:"icon"`       //! 头像
	Portrait  int               `json:"portrait"`   // 头像框
	UnionName string            `json:"union_name"` //! 军团名字
	Level     int               `json:"level"`      //! 等级
	HeroInfo  []*BattleHeroInfo `json:"heroinfo"`   // 双方数据
}

type BattleInfo struct {
	Id       int64              `json:"id"`       //! 战报ID
	LevelID  int                `json:"levelid"`  //! 关卡id
	Type     int                `json:"type"`     //! 战报类型
	Time     int64              `json:"time"`     // 发生的时间
	Result   int                `json:"result"`   // 结果
	Random   int64              `json:"random"`   // 随机数
	UserInfo [2]*BattleUserInfo `json:"userinfo"` // 己方玩家数据
	Weaken   []*WeakenInfo      `json:"weaken"`   // 压制
}

// 爬塔排行榜
type TowerMgr struct {
	migrateOK       bool                             //! 迁移OK
	Locker          *sync.RWMutex                    //! 数据锁
	Sql_TowerPlayer map[int64]*San_TowerPlayerRecord //! 玩家战斗记录

}

func (self *San_TowerPlayerRecord) Decode() { //! 将数据库数据写入data
}

func (self *San_TowerPlayerRecord) Encode() { //! 将data数据写入数据库
}

// 记录第一层 用来排序 key值是玩家uid加上关卡id
type San_TowerPlayerRecord struct {
	KeyID       int64  `json:"keyid"`       // key值
	Name        string `json:"name"`        // 玩家名字
	Uid         int64  `json:"uid"`         // 玩家uid
	Icon        int    `json:"icon"`        // 头像
	Portrait    int    `json:"portrait"`    //
	Level       int    `json:"level"`       // 等级
	PlayerFight int64  `json:"playerfight"` // 玩家战力
	BattleFight int64  `json:"battlefight"` // 战斗参与的战力
	Time        int64  `json:"time"`        // 时间

	db.DataUpdate
}

var towerMgrSingleton *TowerMgr = nil

func GetTowerMgr() *TowerMgr {
	if towerMgrSingleton == nil {
		towerMgrSingleton = new(TowerMgr)
		towerMgrSingleton.Locker = new(sync.RWMutex)
		towerMgrSingleton.migrateOK = false
		towerMgrSingleton.Sql_TowerPlayer = make(map[int64]*San_TowerPlayerRecord)
	}

	return towerMgrSingleton
}
func (self *TowerMgr) GetAllData() {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	var playerRecord San_TowerPlayerRecord
	sql := fmt.Sprintf("select * from `tbl_towerrecord`")
	res := db.GetDBMgr().DBUser.GetAllData(sql, &playerRecord)
	for i := 0; i < len(res); i++ {
		data, ok1 := res[i].(*San_TowerPlayerRecord)
		if !ok1 {
			continue
		}

		data.Init("tbl_towerrecord", data, false)
		data.Decode()
		_, ok2 := self.Sql_TowerPlayer[data.KeyID]
		if !ok2 {
			self.Sql_TowerPlayer[data.KeyID] = data
		}
	}
}

//! 数据存储
func (self *TowerMgr) OnSave() {

}

// 开启迁移数据携程
func (self *TowerMgr) RunMigrate() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			utils.LogError(x, string(debug.Stack()))
		}
	}()

	infoName := "san_towerbattleinfo"
	recordName := "san_towerbattlerecord"
	recordType := core.BATTLE_TYPE_TOWER
	tableName := "tbl_crossarenarecord"

	len, err := db.GetRedisMgr().HLen(infoName)
	if err != nil {
		return
	}
	utils.LogInfo("迁移数据：", infoName, len)

	migOK, err := db.GetRedisMgr().Exists(infoName + "_migrateOKNew")
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
		cursor1, num := self.MigrateDataOne(infoName, recordName, tableName, recordType, cursor)
		count += num
		cursor = cursor1

		if count >= len {
			db.GetRedisMgr().Set(infoName+"_migrateOKNew", "1")
			break
		}

		//! 延迟1ms
		time.Sleep(time.Millisecond)
	}

	utils.LogInfo(infoName, "迁移数据OK")

}

//infoName san_towerbattleinfo
//recordName san_towerbattlerecord
//recordType core.BATTLE_TYPE_TOWER
// tableName tbl_crossarenarecord
func (self *TowerMgr) MigrateDataOne(infoName, recordName, tableName string, recordType int, cursor int64) (int64, int) {
	var info BattleInfo
	var record BattleRecord
	cursor, dataSlice, err := db.GetRedisMgr().HScan(infoName, cursor, "*", 20) //c.Do("hscan", infoName, 1, "match *", "count 1")
	if err != nil {
		//! 读取配置
		self.migrateOK = true
		return 0, 0
	}

	lenArr := len(dataSlice) / 2

	if lenArr > 0 {
		for i := 0; i < lenArr; i++ {
			json.Unmarshal([]byte(dataSlice[i*2+1]), &info)

			v1, ret, err := db.HGetRedis(recordName, fmt.Sprintf("%d", info.Id))
			if err != nil {
				self.migrateOK = true
				return 0, 0
			}

			//! 获取成功
			if ret {
				json.Unmarshal([]byte(v1), &record)
			}

			var db_battleInfo match.JS_CrossArenaBattleInfo
			sql := fmt.Sprintf("select * from `%s` where fightid=%d limit 1;", tableName, info.Id)
			ret1 := db.GetDBMgr().DBUser.GetOneData(sql, &db_battleInfo, "", 0)
			if ret1 == true && db_battleInfo.Id > 0 { //! 获取成功
				utils.LogInfo("已存在，更新：", cursor, dataSlice[0], info.Id)
				db_battleInfo.BattleInfo = utils.HF_CompressAndBase64(game.HF_JtoB(&info))
				db_battleInfo.BattleRecord = utils.HF_CompressAndBase64(game.HF_JtoB(&record))
				db_battleInfo.UpdateTime = time.Now().Unix()
				db_battleInfo.Update(true, false)
			} else {
				db_battleInfo.FightId = info.Id
				db_battleInfo.RecordType = recordType
				db_battleInfo.BattleInfo = utils.HF_CompressAndBase64(game.HF_JtoB(&info))
				db_battleInfo.BattleRecord = utils.HF_CompressAndBase64(game.HF_JtoB(&record))
				db_battleInfo.UpdateTime = time.Now().Unix()
				db.InsertTable(tableName, &db_battleInfo, 0, false)
			}
		}
	}

	return cursor, lenArr
}

// 添加支援英雄
func (self *TowerMgr) AddPlayerRecord(key int64, data *San_TowerPlayerRecord, info *BattleInfo, record *BattleRecord) bool {
	self.Locker.Lock()
	value, ok := self.Sql_TowerPlayer[key]
	self.Locker.Unlock()
	if ok {
		if data.BattleFight < value.BattleFight {
			value.Uid = data.Uid
			value.Name = data.Name
			value.Icon = data.Icon
			value.Portrait = data.Portrait
			value.Level = data.Level
			value.PlayerFight = data.PlayerFight
			value.BattleFight = data.BattleFight
			value.Time = data.Time
			db.HMSetRedisEx("san_towerbattleinfo", info.Id, info, utils.HOUR_SECS*12)
			db.HMSetRedisEx("san_towerbattlerecord", record.Id, record, utils.HOUR_SECS*12)

			var db_battleInfo match.JS_CrossArenaBattleInfo
			sql := fmt.Sprintf("select * from `tbl_crossarenarecord` where fightid=%d limit 1;", info.Id)
			ret := db.GetDBMgr().DBUser.GetOneData(sql, &db_battleInfo, "", 0)
			if ret == true { //! 获取成功
				db_battleInfo.BattleInfo = utils.HF_CompressAndBase64(game.HF_JtoB(&info))
				db_battleInfo.BattleRecord = utils.HF_CompressAndBase64(game.HF_JtoB(&record))
				db_battleInfo.UpdateTime = time.Now().Unix()
				db_battleInfo.Update(true, false)
			} else {
				db_battleInfo.FightId = info.Id
				db_battleInfo.RecordType = core.BATTLE_TYPE_TOWER
				db_battleInfo.BattleInfo = utils.HF_CompressAndBase64(game.HF_JtoB(&info))
				db_battleInfo.BattleRecord = utils.HF_CompressAndBase64(game.HF_JtoB(&record))
				db_battleInfo.UpdateTime = time.Now().Unix()
				db.InsertTable("tbl_crossarenarecord", &db_battleInfo, 0, false)
			}
		}
	} else {
		db.InsertTable("tbl_towerrecord", data, 0, false)
		data.Init("tbl_towerrecord", data, false)
		self.Locker.Lock()
		self.Sql_TowerPlayer[key] = data
		self.Locker.Unlock()
		db.HMSetRedisEx("san_towerbattleinfo", info.Id, info, utils.HOUR_SECS*12)
		db.HMSetRedisEx("san_towerbattlerecord", record.Id, record, utils.HOUR_SECS*12)

		var db_battleInfo match.JS_CrossArenaBattleInfo
		db_battleInfo.FightId = info.Id
		db_battleInfo.RecordType = core.BATTLE_TYPE_TOWER
		db_battleInfo.BattleInfo = utils.HF_CompressAndBase64(game.HF_JtoB(&info))
		db_battleInfo.BattleRecord = utils.HF_CompressAndBase64(game.HF_JtoB(&record))
		db_battleInfo.UpdateTime = time.Now().Unix()
		db.InsertTable("tbl_crossarenarecord", &db_battleInfo, 0, false)
	}
	return true
}

// 获得可使用的英雄
func (self *TowerMgr) GetRecordList(keys map[int64]int64) *Js_TowerFightRecord {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	ret := []*San_TowerPlayerRecord{}
	for _, key := range keys {
		data, ok := self.Sql_TowerPlayer[key]
		if !ok {
			continue
		}
		ret = append(ret, data)
	}

	nLen := len(ret)

	if nLen <= 0 {
		return nil
	}

	sort.Sort(lstTowerPlayer((ret)))

	temp := Js_TowerFightRecord{}
	temp.Time = ret[0].Time
	temp.Key = ret[0].KeyID
	temp.Name = ret[0].Name
	temp.Uid = ret[0].Uid
	temp.Icon = ret[0].Icon
	temp.Portrait = ret[0].Portrait
	temp.Level = ret[0].Level
	temp.PlayerFight = ret[0].PlayerFight
	temp.BattleFight = ret[0].BattleFight
	temp.Time = ret[0].Time

	return &temp
}

// 获得可使用的英雄
func (self *TowerMgr) GetRecord(key int64) *Js_TowerFightRecord {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	data, ok := self.Sql_TowerPlayer[key]
	if !ok {
		return nil
	}

	temp := Js_TowerFightRecord{}
	temp.Time = data.Time
	temp.Key = data.KeyID
	temp.Name = data.Name
	temp.Uid = data.Uid
	temp.Icon = data.Icon
	temp.Portrait = data.Portrait
	temp.Level = data.Level
	temp.PlayerFight = data.PlayerFight
	temp.BattleFight = data.BattleFight
	temp.Time = data.Time

	return &temp
}

type lstTowerPlayer []*San_TowerPlayerRecord

func (s lstTowerPlayer) Len() int      { return len(s) }
func (s lstTowerPlayer) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstTowerPlayer) Less(i, j int) bool {
	if s[i].BattleFight < s[j].BattleFight {
		return true
	}

	if s[i].Time < s[j].Time {
		return true
	}

	if s[i].PlayerFight > s[j].PlayerFight {
		return true
	}

	if s[i].Uid < s[j].Uid {
		return true
	}

	return true
}
func (self *TowerMgr) GetBattleInfo(id int64) *BattleInfo {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var battleInfo BattleInfo
	value, flag, err := db.HGetRedisEx(`san_towerbattleinfo`, id, fmt.Sprintf("%d", id))
	if err != nil || !flag {
		//读数据库
		var db_battleInfo match.JS_CrossArenaBattleInfo
		sql := fmt.Sprintf("select * from `tbl_crossarenarecord` where fightid=%d limit 1;", id)
		ret := db.GetDBMgr().DBUser.GetOneData(sql, &db_battleInfo, "", 0)
		if ret == true { //! 获取成功
			//! 进行处理
			err := json.Unmarshal(utils.HF_Base64AndDecompress(db_battleInfo.BattleInfo), &battleInfo)
			if err != nil {
				utils.LogDebug("Decode Error")
				return nil
			}

			if battleInfo.Id != 0 {
				return &battleInfo
			}
			//! 详细处理
		}

		return &battleInfo
	}
	if flag {
		err := json.Unmarshal([]byte(value), &battleInfo)
		if err != nil {
			return &battleInfo
		}
	}

	if battleInfo.Id != 0 {
		return &battleInfo
	}
	return nil
}

func (self *TowerMgr) GetBattleRecord(id int64) *BattleRecord {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var battleRecord BattleRecord
	value, flag, err := db.HGetRedisEx(`san_towerbattlerecord`, id, fmt.Sprintf("%d", id))
	if err != nil || !flag {
		var db_battleInfo match.JS_CrossArenaBattleInfo
		sql := fmt.Sprintf("select * from `tbl_crossarenarecord` where fightid=%d limit 1;", id)
		ret := db.GetDBMgr().DBUser.GetOneData(sql, &db_battleInfo, "", 0)
		if ret == true { //! 获取成功
			//! 进行处理
			err := json.Unmarshal(utils.HF_Base64AndDecompress(db_battleInfo.BattleRecord), &battleRecord)
			if err != nil {
				utils.LogDebug("Decode Error")
				return nil
			}

			if battleRecord.Id != 0 {
				return &battleRecord
			}
		}
	}
	if flag {
		err := json.Unmarshal([]byte(value), &battleRecord)
		if err != nil {
			return &battleRecord
		}
	}
	if battleRecord.Id != 0 {
		return &battleRecord
	}
	return nil
}
