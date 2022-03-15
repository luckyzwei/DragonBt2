/*
@Time : 2020/4/22 10:27 
@Author : 96121
@File : mgr_server 游戏服节点
@Software: GoLand
*/
package server

import (
	"fmt"
	"master/db"
	"sync"
	"time"
)

const (
	SERVER_TABLE_NAME = "tbl_server"
	SERVER_TBALE_SQL  = "select * from %s where serverid = %d"
)

type ServerMgr struct {
	MapServer map[int]*GameServer //! 服务器节点
	Locker    *sync.RWMutex       //! 数据锁
}

var s_servermgr *ServerMgr = nil

func GetServerMgr() *ServerMgr {
	if s_servermgr == nil {
		s_servermgr = new(ServerMgr)

		s_servermgr.MapServer = make(map[int]*GameServer)
		s_servermgr.Locker = new(sync.RWMutex)
	}

	return s_servermgr
}



func (self *ServerMgr) OnLogic() {

}

func (self *ServerMgr) OnSave() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, server := range self.MapServer {
		server.OnSave()
	}
}

func (self *ServerMgr) GetServer(sid int, create bool) *GameServer {
	self.Locker.RLock()
	if server, ok := self.MapServer[sid]; ok {
		self.Locker.RUnlock()
		return server
	} else {
		self.Locker.RUnlock()
		if create {
			var sqlGameServer SQL_GameServer
			//sql := fmt.Sprintf(SERVER_TBALE_SQL, SERVER_TABLE_NAME)
			db.GetDBMgr().DBUser.GetOneData(fmt.Sprintf(SERVER_TBALE_SQL, SERVER_TABLE_NAME, sid), &sqlGameServer, "", 0)
			if sqlGameServer.ServerId <= 0 { //! 数据库内找不到,插入新的
				server = new(GameServer)
				server.Data.ServerId = sid
				server.Data.Id = 0
				server.EventArr = make(chan *ServerEvent, Max_Server_Event)

				db.InsertTable(SERVER_TABLE_NAME, &server.Data, 0, false)
				server.Data.Init(SERVER_TABLE_NAME, &server.Data, false)
				return server
			} else {
				server = new(GameServer)
				server.Data = sqlGameServer
				server.EventArr = make(chan *ServerEvent, Max_Server_Event)
				server.Data.Init(SERVER_TABLE_NAME, &server.Data, false)

				self.Locker.Lock()
				self.MapServer[sid] = server
				self.Locker.Unlock()

				return server
			}
		}
	}

	return nil
}

func (self *ServerMgr) ResServer(sid int, name string, online int) bool {
	if server, ok := self.MapServer[sid]; ok {
		self.Locker.RUnlock()
		server.Data.Online = online
		server.Data.LastUpdate = time.Now().Unix()
	} else {
		server = new(GameServer)
		server.Data.ServerId = sid
		server.Data.Online = online
		server.Data.Name = name

		server.Data.Init(SERVER_TABLE_NAME, &server.Data, false)
		db.InsertTable(SERVER_TABLE_NAME, &server.Data, 0, false)

		self.Locker.Lock()
		self.MapServer[sid] = server
		self.Locker.Unlock()
	}

	return false
}
