package player

import (
	"encoding/json"
	"master/center/tower"
	"master/utils"
)

//! 错误码定义
const (
	TOWER_SUCCESS = 0 //! 没有错误
	TOWER_ERRER   = 1 //! 错误
)

///////////添加表//////////////
type S2M_TowerAddPlayerRecord struct {
	Key    int64
	Data   *tower.San_TowerPlayerRecord
	Info   *tower.BattleInfo
	Record *tower.BattleRecord
}

type M2S_TowerAddPlayerRecord struct {
}

///////////获取列表//////////////
type S2M_TowerGetRecordList struct {
	Keys map[int64]int64
}

type M2S_TowerGetRecordList struct {
	Data *tower.Js_TowerFightRecord
}

///////////获取单条数据//////////////
type S2M_TowerGetRecord struct {
	Key int64
}

type M2S_TowerGetRecord struct {
	Data *tower.Js_TowerFightRecord
}

///////////获得信息//////////////
type S2M_TowerGetBattleInfo struct {
	Key int64
}

type M2S_TowerGetBattleInfo struct {
	Data *tower.BattleInfo
}

///////////获得具体数据//////////////
type S2M_TowerGetBattleRecord struct {
	Key int64
}

type M2S_TowerGetBattleRecord struct {
	Data *tower.BattleRecord
}

//! 操作请求
type RPC_TowerAction struct {
	Data string
}

//! 操作响应
type RPC_TowerActionRet struct {
	RetCode int //! 结果码
	Data    string
}

type RPC_Tower struct {
}

func (self *RPC_Tower) AddPlayerRecord(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_TowerAddPlayerRecord
	json.Unmarshal([]byte(req.Data), &msg)

	if !tower.GetTowerMgr().AddPlayerRecord(msg.Key, msg.Data, msg.Info, msg.Record) {
		ret.RetCode = TOWER_ERRER
		return nil
	}

	ret.RetCode = TOWER_SUCCESS
	return nil
}

func (self *RPC_Tower) GetRecordList(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_TowerGetRecordList
	json.Unmarshal([]byte(req.Data), &msg)

	var backmsg M2S_TowerGetRecordList
	backmsg.Data = tower.GetTowerMgr().GetRecordList(msg.Keys)

	if nil == backmsg.Data {
		ret.RetCode = TOWER_ERRER
		return nil
	}

	ret.Data = utils.HF_JtoA(backmsg)
	ret.RetCode = TOWER_SUCCESS
	return nil
}

func (self *RPC_Tower) GetRecord(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_TowerGetRecord
	json.Unmarshal([]byte(req.Data), &msg)

	var backmsg M2S_TowerGetRecord
	backmsg.Data = tower.GetTowerMgr().GetRecord(msg.Key)

	if nil == backmsg.Data {
		ret.RetCode = TOWER_ERRER
		return nil
	}

	ret.Data = utils.HF_JtoA(backmsg)
	ret.RetCode = TOWER_SUCCESS
	return nil
}
func (self *RPC_Tower) GetBattleInfo(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_TowerGetBattleInfo
	json.Unmarshal([]byte(req.Data), &msg)

	var backmsg M2S_TowerGetBattleInfo
	backmsg.Data = tower.GetTowerMgr().GetBattleInfo(msg.Key)

	if nil == backmsg.Data {
		ret.RetCode = TOWER_ERRER
		return nil
	}

	ret.Data = utils.HF_JtoA(backmsg)
	ret.RetCode = TOWER_SUCCESS
	return nil
}

func (self *RPC_Tower) GetBattleRecord(req RPC_UnionAction, ret *RPC_UnionActionRet) error {
	var msg S2M_TowerGetBattleRecord
	json.Unmarshal([]byte(req.Data), &msg)

	var backmsg M2S_TowerGetBattleRecord
	backmsg.Data = tower.GetTowerMgr().GetBattleRecord(msg.Key)

	if nil == backmsg.Data {
		ret.RetCode = TOWER_ERRER
		return nil
	}

	ret.Data = utils.HF_JtoA(backmsg)
	ret.RetCode = TOWER_SUCCESS
	return nil
}
