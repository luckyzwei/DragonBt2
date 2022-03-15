/*
@Time : 2020/5/10 11:10
@Author : 96121
@File : proto_player
@Software: GoLand
*/
package game

import (
	"net/rpc"
	"sync"
)

const (
	RPC_TOWER_ADD_PLAYER_RECORD = "RPC_Tower.AddPlayerRecord" //! 添加记录
	RPC_TOWER_GET_RECORD_LIST   = "RPC_Tower.GetRecordList"   //! 获得列表
	RPC_TOWER_GET_RECORD        = "RPC_Tower.GetRecord"       //! 获得单条数据
	RPC_TOWER_GET_BATTLE_INFO   = "RPC_Tower.GetBattleInfo"   //! 获得单条数据
	RPC_TOWER_GET_BATTLE_RECORD = "RPC_Tower.GetBattleRecord" //! 获得单条数据
)

///////////添加表//////////////
type S2M_TowerAddPlayerRecord struct {
	Key    int64
	Data   *San_TowerPlayerRecord
	Info   *BattleInfo
	Record *BattleRecord
}

type M2S_TowerAddPlayerRecord struct {
}

///////////获取列表//////////////
type S2M_TowerGetRecordList struct {
	Keys map[int64]int64
}

type M2S_TowerGetRecordList struct {
	Data *Js_TowerFightRecord
}

///////////获取单条数据//////////////
type S2M_TowerGetRecord struct {
	Key int64
}

type M2S_TowerGetRecord struct {
	Data *Js_TowerFightRecord
}

///////////获得信息//////////////
type S2M_TowerGetBattleInfo struct {
	Key int64
}

type M2S_TowerGetBattleInfo struct {
	Data *BattleInfo
}

///////////获得具体数据//////////////
type S2M_TowerGetBattleRecord struct {
	Key int64
}

type M2S_TowerGetBattleRecord struct {
	Data *BattleRecord
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
	Client       *rpc.Client
	PlayerLocker *sync.RWMutex //! 数据锁
}

func (self *RPC_Tower) Init() bool {
	self.PlayerLocker = new(sync.RWMutex)
	return true
}

//操作
func (self *RPC_Tower) TowerAction(action string, data interface{}) *RPC_TowerActionRet {
	if self.Client != nil {
		var req RPC_TowerAction
		req.Data = HF_JtoA(data)

		var ret RPC_TowerActionRet
		GetMasterMgr().CallEx(self.Client,action, req, &ret)
		return &ret
	}

	return nil
}
