package db

import (
	"log"
	"master/utils"
)

type SQL_Set struct {
	Sql string
}

type GameDataMgr struct {
	SqlChan     chan *SQL_Set //! 模块数据SQL
	SqlBaseChan chan *SQL_Set //! userbase的缓存
}

var s_gamedatamgr *GameDataMgr

func GetGameDataMgr() *GameDataMgr {
	if s_gamedatamgr == nil {
		s_gamedatamgr = new(GameDataMgr)
		s_gamedatamgr.SqlChan = make(chan *SQL_Set, MAX_LOG_NUM)
		s_gamedatamgr.SqlBaseChan = make(chan *SQL_Set, MAX_LOG_NUM)
	}

	return s_gamedatamgr
}

func (self *GameDataMgr) GoService() {
	go self.RunSqlSet()
	go self.RunSqlBaseSet()
}

//! redis更新数据-备份数据
func (self *GameDataMgr) SqlSet(sql string) {
	if self.SqlChan == nil {
		return
	}
	self.SqlChan <- &SQL_Set{sql}
}

//! 数据保存-多线程
func (self *GameDataMgr) RunSqlSet() {
	log.Println("sqlset start...")
	GetDBMgr().Wait.Add(1)
	var flag = false
	for msg := range self.SqlChan {
		if msg.Sql == "" {
			break
		}
		flag = false
		for i := 0; i < 10; i++ {
			utils.LogDebug("SQL BASE EXEC:", msg.Sql)
			_, _, ok := GetDBMgr().DBUser.Exec(msg.Sql)
			if ok {
				flag = true
				break
			}
		}

		if !flag {
			utils.LogError("sql execute failed, sql = ", msg.Sql)
		}
	}
	GetDBMgr().Wait.Done()
	log.Println("sqlset完毕")
	//if self.SqlChan != nil {
	//	close(self.SqlChan)
	//	self.SqlChan = nil
	//}
}

//! sql语句-非Reis数据，数据保存-全局数据
func (self *GameDataMgr) SqlBaseSet(sql string) {
	if self.SqlBaseChan == nil {
		return
	}
	self.SqlBaseChan <- &SQL_Set{sql}
}

//! 全局数据保存
func (self *GameDataMgr) RunSqlBaseSet() {
	log.Println("sqlbaseset start...")
	GetDBMgr().Wait.Add(1)
	for msg := range self.SqlBaseChan {
		if msg.Sql == "" {
			break
		}
		utils.LogDebug("SQL BASE EXEC:", msg.Sql)
		for i := 0; i < 10; i++ {
			_, _, ok := GetDBMgr().DBUser.Exec(msg.Sql)
			if ok {
				break
			}
		}
	}
	GetDBMgr().Wait.Done()
	log.Println("sqlbaseset完毕")
	close(self.SqlBaseChan)
	self.SqlBaseChan = nil
}
