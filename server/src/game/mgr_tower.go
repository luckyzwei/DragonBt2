package game

import (
	"encoding/json"
	"sync"
)

// 爬塔排行榜
type TowerMgr struct {
	Locker *sync.RWMutex //! 数据锁
	//Sql_TowerPlayer map[int64]*San_TowerPlayerRecord //! 玩家战斗记录
}

//func (self *San_TowerPlayerRecord) Decode() { //! 将数据库数据写入data
//}
//
//func (self *San_TowerPlayerRecord) Encode() { //! 将data数据写入数据库
//}

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

	DataUpdate
}

var towerMgrSingleton *TowerMgr = nil

func GetTowerMgr() *TowerMgr {
	if towerMgrSingleton == nil {
		towerMgrSingleton = new(TowerMgr)
		towerMgrSingleton.Locker = new(sync.RWMutex)
		//towerMgrSingleton.Sql_TowerPlayer = make(map[int64]*San_TowerPlayerRecord)
	}

	return towerMgrSingleton
}
func (self *TowerMgr) GetData() {
	//var playerRecord San_TowerPlayerRecord
	//sql := fmt.Sprintf("select * from `san_towerplayerrecord`")
	//res := GetServer().DBUser.GetAllData(sql, &playerRecord)
	//for i := 0; i < len(res); i++ {
	//	data, ok1 := res[i].(*San_TowerPlayerRecord)
	//	if !ok1 {
	//		continue
	//	}
	//
	//	data.Init("san_towerplayerrecord", data, false)
	//	data.Decode()
	//	_, ok2 := self.Sql_TowerPlayer[data.KeyID]
	//	if !ok2 {
	//		self.Sql_TowerPlayer[data.KeyID] = data
	//	}
	//}
}

// 添加支援英雄
func (self *TowerMgr) AddPlayerRecord(key int64, data *San_TowerPlayerRecord, info *BattleInfo, record *BattleRecord) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	var msg S2M_TowerAddPlayerRecord
	msg.Key = key
	msg.Data = data
	msg.Info = info
	msg.Record = record
	ret := GetMasterMgr().TowerRPC.TowerAction(RPC_TOWER_ADD_PLAYER_RECORD, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	var backmsg M2S_TowerAddPlayerRecord
	json.Unmarshal([]byte(ret.Data), &backmsg)

	//_, ok := self.Sql_TowerPlayer[key]
	//if ok {
	//	if data.BattleFight < self.Sql_TowerPlayer[key].BattleFight {
	//		self.Sql_TowerPlayer[key].Uid = data.Uid
	//		self.Sql_TowerPlayer[key].Name = data.Name
	//		self.Sql_TowerPlayer[key].Icon = data.Icon
	//		self.Sql_TowerPlayer[key].Level = data.Level
	//		self.Sql_TowerPlayer[key].PlayerFight = data.PlayerFight
	//		self.Sql_TowerPlayer[key].BattleFight = data.BattleFight
	//		self.Sql_TowerPlayer[key].Time = data.Time
	//		HMSetRedis("san_towerbattleinfo", info.Id, info, DAY_SECS*10)
	//		HMSetRedis("san_towerbattlerecord", record.Id, record, DAY_SECS*10)
	//	}
	//} else {
	//	InsertTable("san_towerplayerrecord", data, 0, false)
	//	data.Init("san_towerplayerrecord", data, false)
	//	self.Sql_TowerPlayer[key] = data
	//	HMSetRedis("san_towerbattleinfo", info.Id, info, DAY_SECS*10)
	//	HMSetRedis("san_towerbattlerecord", record.Id, record, DAY_SECS*10)
	//}
	return true
}

// 获取关卡配置
func (self *TowerMgr) GetLevelConfig(levelid int) *LevelConfig {
	ret, ok := GetCsvMgr().LevelConfigMap[levelid]
	if !ok {
		return nil
	}
	return ret
}

// 获得可使用的英雄
func (self *TowerMgr) GetRecordList(keys map[int64]int64) *Js_TowerFightRecord {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	//ret := []*San_TowerPlayerRecord{}
	//for _, key := range keys {
	//	data, ok := self.Sql_TowerPlayer[key]
	//	if !ok {
	//		continue
	//	}
	//	ret = append(ret, data)
	//}
	//
	//nLen := len(ret)
	//
	//if nLen <= 0 {
	//	return nil
	//}
	//
	//sort.Sort(lstTowerPlayer((ret)))
	//
	//temp := Js_TowerFightRecord{}
	//temp.Time = ret[0].Time
	//temp.Key = ret[0].KeyID
	//temp.Name = ret[0].Name
	//temp.Uid = ret[0].Uid
	//temp.Icon = ret[0].Icon
	//temp.Level = ret[0].Level
	//temp.PlayerFight = ret[0].PlayerFight
	//temp.BattleFight = ret[0].BattleFight
	//temp.Time = ret[0].Time

	var msg S2M_TowerGetRecordList
	msg.Keys = keys
	ret := GetMasterMgr().TowerRPC.TowerAction(RPC_TOWER_GET_RECORD_LIST, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return nil
	}

	var backmsg M2S_TowerGetRecordList
	json.Unmarshal([]byte(ret.Data), &backmsg)

	return backmsg.Data
}

// 获得可使用的英雄
func (self *TowerMgr) GetRecord(key int64) *Js_TowerFightRecord {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	//data, ok := self.Sql_TowerPlayer[key]
	//if !ok {
	//	return nil
	//}
	//
	//temp := Js_TowerFightRecord{}
	//temp.Time = data.Time
	//
	//temp.Key = data.KeyID
	//temp.Name = data.Name
	//temp.Uid = data.Uid
	//temp.Icon = data.Icon
	//temp.Level = data.Level
	//temp.PlayerFight = data.PlayerFight
	//temp.BattleFight = data.BattleFight
	//temp.Time = data.Time

	var msg S2M_TowerGetRecord
	msg.Key = key
	ret := GetMasterMgr().TowerRPC.TowerAction(RPC_TOWER_GET_RECORD, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return nil
	}

	var backmsg M2S_TowerGetRecord
	json.Unmarshal([]byte(ret.Data), &backmsg)

	return backmsg.Data
}

//
//type lstTowerPlayer []*San_TowerPlayerRecord
//
//func (s lstTowerPlayer) Len() int      { return len(s) }
//func (s lstTowerPlayer) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
//func (s lstTowerPlayer) Less(i, j int) bool {
//	if s[i].BattleFight < s[j].BattleFight {
//		return true
//	}
//
//	if s[i].Time < s[j].Time {
//		return true
//	}
//
//	if s[i].PlayerFight > s[j].PlayerFight {
//		return true
//	}
//
//	if s[i].Uid < s[j].Uid {
//		return true
//	}
//
//	return true
//}

func (self *TowerMgr) GetBattleInfo(key int64) *BattleInfo {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	var msg S2M_TowerGetBattleInfo
	msg.Key = key
	ret := GetMasterMgr().TowerRPC.TowerAction(RPC_TOWER_GET_BATTLE_INFO, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return nil
	}

	var backmsg M2S_TowerGetBattleInfo
	json.Unmarshal([]byte(ret.Data), &backmsg)

	return backmsg.Data
}

func (self *TowerMgr) GetBattleRecord(key int64) *BattleRecord {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	var msg S2M_TowerGetBattleRecord
	msg.Key = key
	ret := GetMasterMgr().TowerRPC.TowerAction(RPC_TOWER_GET_BATTLE_RECORD, &msg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return nil
	}

	var backmsg M2S_TowerGetBattleRecord
	json.Unmarshal([]byte(ret.Data), &backmsg)

	return backmsg.Data
}
