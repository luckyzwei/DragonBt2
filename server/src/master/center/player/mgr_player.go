/*
@Time : 2020/4/22 10:18
@Author : 96121
@File : mgr_player 角色管理器
@Software: GoLand
*/
package player

import (
	"fmt"
	"game"
	"master/core"
	"master/db"
	"master/utils"
	"strings"
	"sync"
)

//! 角色列表
type PlayerMgr struct {
	MapPlayer map[int64]*Player //! 角色存档
	Locker    *sync.RWMutex     //! 数据锁

	TestPlayer core.IPlayer
}

//! 单例模式
var s_playermgr *PlayerMgr

func GetPlayerMgr() *PlayerMgr {
	if s_playermgr == nil {
		s_playermgr = new(PlayerMgr)
		s_playermgr.MapPlayer = make(map[int64]*Player)
		s_playermgr.Locker = new(sync.RWMutex)

		core.PlayerMgr = s_playermgr
	}

	return s_playermgr
}

func (self *PlayerMgr) Init() {
	//! 初始化 
}

//! 保存数据
func (self *PlayerMgr) OnSave() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, v := range self.MapPlayer {
		v.onSave(false)
		//v.DataFriend.onSave()
	}
}

//! 获取新玩家
func (self *PlayerMgr) GetPlayer(uid int64, create bool) *Player {
	//! 优化读写锁
	//self.Locker.Lock()
	//defer self.Locker.Unlock()

	self.Locker.RLock()
	player, ok := self.MapPlayer[uid]
	self.Locker.RUnlock()
	if ok {
		if player.Data.data == nil {
			player.Data.data = new(JS_PlayerData)
		}

		if player.Data.lifetree == nil {
			player.Data.lifetree = new(game.JS_LifeTreeInfo)
		}
		return player
	} else {
		if create { //! 不存在则载入
			player = new(Player)
			player.onGetData(uid)
			self.Locker.Lock()
			self.MapPlayer[uid] = player
			self.Locker.Unlock()
		}
	}

	return player
}

func (self *PlayerMgr) GetPlayerArena(uid int64, res *RPC_PlayerData_Res) {
	player := GetPlayerMgr().GetPlayer(uid, true)
	if player != nil {
		res.Data = utils.HF_JtoA(player.Data.data.ArenaFightInfo)
	}
}

//! 获取新玩家
func (self *PlayerMgr) GetCorePlayer(uid int64, create bool) core.IPlayer {
	//self.Locker.Lock()
	//defer self.Locker.Unlock()

	if uid == 0 {
		return nil
	}

	self.Locker.RLock()
	player, ok := self.MapPlayer[uid]
	self.Locker.RUnlock()
	if ok {
		return player
	} else {
		if create { //! 不存在则载入
			player := new(Player)
			player.onGetData(uid)

			self.Locker.Lock()
			self.MapPlayer[uid] = player
			self.Locker.Unlock()
		}
	}

	return nil
}

func (self *PlayerMgr) GetOnline() int {
	return 0
}

//! 发布消息
func (self *PlayerMgr) BroadcastMsg(head string, body []byte) {

}

const (
	PLAYER_TABLE_ALL_SQL = "select uid,uname from %s limit %d,%d;"
)

type PlayerNameByUid struct {
	Uid   int64
	Uname string
}

func (self *PlayerMgr) LoadAllName() {
	var playerData PlayerNameByUid
	for i := 0; i < 100; i++ {
		sql := fmt.Sprintf(PLAYER_TABLE_ALL_SQL, PLAYER_TABLE_NAME, i*1000, (i+1)*1000)
		res := db.GetDBMgr().DBUser.GetAllData(sql, &playerData)
		num := len(res)
		if num == 0 {
			break
		}
		utils.LogDebug("set player name num : %d", num)
		for j := 0; j < num; j++ {
			p := res[j].(*PlayerNameByUid)
			if p != nil {
				db.GetRedisMgr().Set(fmt.Sprintf("%s_%d", p.Uname, p.Uid), fmt.Sprintf("%d", p.Uid))
			}
		}
	}
}

func (self *PlayerMgr) AddNameMap(uname string, uid int64) {
	db.GetRedisMgr().Set(fmt.Sprintf("%s_%d", uname, uid), fmt.Sprintf("%d", uid))
}

func (self *PlayerMgr) GetUidByName(uname string) []int64 {
	ret := make([]int64, 0)

	mapName, err := db.GetRedisMgr().Keys(fmt.Sprintf("*%s*", uname))
	if err == nil {
		for _, name := range mapName {
			arrName := strings.Split(name, "_")
			if len(arrName) >= 2 {
				ret = append(ret, utils.HF_AtoI64(arrName[1]))
			}
		}
	}

	return ret
}

//! 获取新玩家
func (self *PlayerMgr) RandPlayer(uid int64, num int) []int64 {

	uids := make([]int64, 0)

	self.Locker.RLock()
	for key := range self.MapPlayer {
		if key == uid {
			continue
		}
		uids = append(uids, key)
		if len(uids) >= num {
			break
		}
	}
	self.Locker.RUnlock()

	return uids

}
