package db

import (
	"log"
	"master/utils"
	"time"
)

const DATEFORMAT = "2006-01-02 15:04:05" // 时间格式化

//! 游戏消耗日志
type SQL_Log struct {
	Id     int64  //! 日志主键-自增Id
	Time   int64  //! 时间
	Type   int    //! 类型-道具Id
	Value  int    //! 数量-正数获得-负数
	Param1 int    //! 参数1
	Param2 int    //! 参数2
	Uid    int64  //! 角色Id
	Dec    string //! 来源-描述
	Cur    int    //! 剩余值
	Param3 int    //! 参数3
}

//! 游戏行为日志
type SQL_BeLog struct {
	Id     int64  //! 日志主键-自增Id
	Time   int64  //! 时间
	Type   int    //! 类型-道具Id
	Value  int    //! 数量-正数获得-负数
	Param1 int    //! 参数1
	Param2 int    //! 参数2
	Uid    int64  //! 角色Id
	Dec    string //! 来源-描述
	Cur    int    //! 剩余值
	Param3 int    //! 参数3
	Level  int    //! 等级
	Vip    int    //! VIP
	Fight  int    //! 战斗力
}

//! 在线状态日志
type SQL_LineLog struct {
	Id      int64  //! 日志主键
	Uid     int64  //! 角色Id
	Time    int64  //! 时间
	Ip      string //! Ip
	Line    int    //! 在线时间
	Creator string //! 来源-渠道
}

//! 邮件日志
type SQL_MailLog struct {
	Id   int64  //! 邮件ID
	Uid  int64  //! 角色ID
	Msg  string //! 内容
	Item string //! 道具
	Time string //! 操作时间
}

//! 在线聊天
type SQL_ChatLog struct {
	Id   int64  //! Id
	Uid  int64  //! 角色Id
	Type int    //! 聊天类型
	Msg  string //! 消息
	Time int64  //! 时间
}

const (
	MAX_LOG_NUM = 200000 //! 日志缓存最长长度
)

//! 游戏日志管理类
type GameLogMgr struct {
	LogConnNum   int //! 道具日志缓存长度
	BeLogConnNum int //! 行为日志
	SqlConnNum   int //! 数据库连接

	SqlMailLogChan chan *SQL_MailLog //! 在线数据SQL
	SqlChatLogChan chan *SQL_ChatLog //! 在线聊天SQL
	SqlLogChan     chan *SQL_Log     //! 游戏道具日志
	SqlBeLogChan   chan *SQL_BeLog   //！游戏行为日志
	SqlLineLogChan chan *SQL_LineLog //! 在线数据SQL
}

//! 游戏日志单例模式
var s_gamelogmgr *GameLogMgr

func GetLogMgr() *GameLogMgr {
	if s_gamelogmgr == nil {
		s_gamelogmgr = new(GameLogMgr)
		//! 初始化
		s_gamelogmgr.Init()
	}

	return s_gamelogmgr
}

func (self *GameLogMgr) Init() bool {
	//self.LogChan = make(chan *UP_Log, 10000)
	self.SqlLogChan = make(chan *SQL_Log, MAX_LOG_NUM)
	self.SqlBeLogChan = make(chan *SQL_BeLog, MAX_LOG_NUM)
	self.SqlLineLogChan = make(chan *SQL_LineLog, MAX_LOG_NUM)
	//self.SqlChan = make(chan *SQL_Set, MAX_LOG_NUM)
	//self.SqlBaseChan = make(chan *SQL_Set, MAX_LOG_NUM)
	self.SqlMailLogChan = make(chan *SQL_MailLog, MAX_LOG_NUM)
	self.SqlChatLogChan = make(chan *SQL_ChatLog, MAX_LOG_NUM)

	return true
}

func (self *GameLogMgr) GoService() {
	go self.RunSqlLog()
	//go self.RunSqlLineLog()
	//go self.RunSqlMailLog()
}

//! 数据库日志
// 玩家id，道具id，增加/减少值，变量1(系统出处)，变量2，描述
func (self *GameLogMgr) SqlLog(uid int64, _type int, value int, param1 int, param2 int, dec string, cur int, param3 int, level, vip, fight int) {
	if _type/10000000 > 0 {
		if self.SqlLogChan == nil {
			utils.LogInfo("self.SqlLogChan == nil")
			return
		}
		if len(self.SqlLogChan) >= 200000 {
			utils.LogError("丢了一个日志")
			return
		}
		self.SqlLogChan <- &SQL_Log{0, time.Now().Unix(), _type, value, param1, param2, uid, dec, cur, param3}
		if uid == 0 {
			//LogInfo("向管道SqlLogChan发送了停止消息")
		}
	} else {
		if self.SqlBeLogChan == nil {
			utils.LogInfo("self.SqlBeLogChan == nil")
			return
		}
		if len(self.SqlBeLogChan) >= 200000 {
			utils.LogError("丢了一个日志")
			return
		}
		//if player != nil {
		//	self.SqlBeLogChan <- &SQL_BeLog{0, time.Now().Unix(), _type, value, param1,
		//		param2, uid, dec, cur, param3, player.GetLevel(),
		//		player.GetVip(), int(player.GetFight() / 100)}
		//} else {
		self.SqlBeLogChan <- &SQL_BeLog{0, time.Now().Unix(), _type, value, param1,
			param2, uid, dec, cur, param3, level, vip, fight}
		//}
	}
}

// 玩家id，道具id，增加/减少值，变量1(系统出处)，变量2，描述
func (self *GameLogMgr) SqlLogEx(uid int64, _type int, value int, param1 int64, param2 int, dec string, cur int, param3 int, level, vip, fight int) {
	if _type/10000000 > 0 {
		if self.SqlLogChan == nil {
			return
		}
		if len(self.SqlLogChan) >= 200000 {
			utils.LogError("丢了一个日志")
			return
		}
		self.SqlLogChan <- &SQL_Log{0, time.Now().Unix(), _type, value, int(param1),
			param2, uid, dec, cur, param3}
	} else {
		if self.SqlBeLogChan == nil {
			return
		}
		if len(self.SqlBeLogChan) >= 200000 {
			utils.LogError("丢了一个日志")
			return
		}
		self.SqlBeLogChan <- &SQL_BeLog{0, time.Now().Unix(), _type, value, int(param1),
			param2, uid, dec, cur, param3, level, vip, fight}
	}
}

//! 游戏日志保存[没有完全关闭]
func (self *GameLogMgr) RunSqlLog() {
	log.Println("sqllog start...")
	GetDBMgr().Wait.Add(1)
	for msg := range self.SqlLogChan {
		if msg.Uid == 0 {
			//LogDebug("RunSqlLog is Closed, len:", len(self.SqlLogChan))
			break
		}
		for i := 0; i < 10; i++ {
			if InsertLogTable("san_log", msg, 1) > 0 {
				break
			}
		}
	}
	GetDBMgr().Wait.Done()
	log.Println("sqllog 完毕")
}

func (self *GameLogMgr) CloseMailLog() {
	if self.SqlMailLogChan == nil {
		return
	}
	self.SqlMailLogChan <- &SQL_MailLog{0, 0, "", "", ""}
}

//! 邮件日志保存
func (self *GameLogMgr) RunSqlMailLog() {
	GetDBMgr().Wait.Add(1)
	for msg := range self.SqlMailLogChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 10; i++ {
			if InsertLogTable("san_mail", msg, 1) > 0 {
				break
			}
		}
	}
	GetDBMgr().Wait.Done()
	log.Println("sql maillog完毕")
}

//! 在线日志
func (self *GameLogMgr) SqlMailLog(uid int64, msg string, item string) {
	if self.SqlMailLogChan == nil {
		return
	}

	if len(self.SqlMailLogChan) >= 200000 {
		utils.LogError("丢了一个日志")
		return
	}

	self.SqlMailLogChan <- &SQL_MailLog{0, uid, msg, item, time.Now().Format(DATEFORMAT)}
}

//! 在线日志
func (self *GameLogMgr) SqlLineLog(uid int64, ip string, line int, creator string) {
	if self.SqlLineLogChan == nil {
		//utils.LogError("self.SqlLineLogChan == nil!!!")
		return
	}

	if len(self.SqlLineLogChan) >= 200000 {
		utils.LogError("丢了一个日志")
		return
	}

	self.SqlLineLogChan <- &SQL_LineLog{0, uid, time.Now().Unix(), ip, line, creator}
}

//! 在线数据日志
func (self *GameLogMgr) RunSqlLineLog() {
	//return
	GetDBMgr().Wait.Add(1)
	for msg := range self.SqlLineLogChan {
		if msg.Uid == 0 {
			//LogDebug("RunSqlLineLog is Closed!!!")
			break
		}
		for i := 0; i < 10; i++ {
			if InsertLogTable("san_linelog", msg, 1) > 0 {
				break
			}
		}
	}
	GetDBMgr().Wait.Done()
	log.Println("sqllinelog完毕")
	if self.SqlLineLogChan != nil {
		close(self.SqlLineLogChan)
		self.SqlLineLogChan = nil
	}
}
