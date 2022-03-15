package db

import (
	"master/utils"
	"sync"
)

// mysql db 管理器
type DBMgr struct {
	//! 数据库
	DBUser   *DBServer       //！database 接口
	DBLog    *DBServer       //! database log 接口
	Wait     *sync.WaitGroup //! 同步阻塞
	InitFlag bool            //! 初始化标志
	Level    int             //! 等级
}

type San_Num struct {
	Num int
}

//! 单例模式
var s_dbmgr *DBMgr = nil

func GetDBMgr() *DBMgr {
	if s_dbmgr == nil {
		s_dbmgr = new(DBMgr)
		s_dbmgr.InitFlag = false
		s_dbmgr.Wait = new(sync.WaitGroup)
	}

	return s_dbmgr
}

//! 初始化-连接数据库
func (self *DBMgr) Init(userDsn string, logDsn string) bool {
	if self.DBUser == nil {
		self.DBUser = new(DBServer)
	}

	//! 初始化 User 数据库连接
	if ret := self.DBUser.Init(userDsn); ret == false {
		utils.LogError("角色数据库初始化失败：%s", userDsn)
		return false
	}

	if self.DBLog == nil {
		self.DBLog = new(DBServer)
	}

	//! 初始化 Log 数据库连接
	if ret := self.DBLog.Init(logDsn); ret == false {
		utils.LogError("日志数据库初始化失败：%s", logDsn)
		return false
	}

	//! 数据库服务
	GetGameDataMgr().GoService()

	//! 日志服务
	GetLogMgr().GoService()

	return true
}

func (self *DBMgr) Close() bool {
	utils.LogInfo("start close db....")

	GetLogMgr().SqlLog(0, 91000001, 0, 0, 0, "", 0, 0, 0, 0, 0)
	GetGameDataMgr().SqlSet("")
	GetGameDataMgr().SqlBaseSet("")

	self.Wait.Wait()

	utils.LogInfo("close db ok....")
	self.DBUser.Close()
	self.DBLog.Close()

	return true
}

//! 获取世界等级
func (self *DBMgr) GetWorldLevel(refresh bool) int {
	if refresh {
		total := 0
		var num San_Num
		res := self.DBUser.GetAllData("select `level` as num from san_userbase order by `level` desc limit 30", &num)
		for i := 0; i < len(res); i++ {
			total += res[i].(*San_Num).Num
		}

		self.Level = utils.HF_MinInt(90, utils.HF_MaxInt(total/30-2, 10))
	}

	return self.Level
}
