package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/websocket"
)

//! 游戏消耗日志
type SQL_Log struct {
	Id     int64
	Time   int64
	Type   int
	Value  int
	Param1 int
	Param2 int
	Uid    int64
	Dec    string
	Cur    int
	Param3 int
}

//! 游戏行为日志
type SQL_BeLog struct {
	Id     int64
	Time   int64
	Type   int
	Value  int
	Param1 int
	Param2 int
	Uid    int64
	Dec    string
	Cur    int
	Param3 int
	Level  int
	Vip    int
	Fight  int
	PassId int
}

type SQL_LineLog struct {
	Id      int64
	Uid     int64
	Time    int64
	Ip      string
	Line    int
	Creator string
}

type SQL_MailLog struct {
	Id   int64
	Uid  int64
	Msg  string
	Item string
	Time string
}

//! 在线聊天
type SQL_ChatLog struct {
	Id   int64  //! Id
	Uid  int64  //! 角色Id
	Type int    //! 聊天类型
	Msg  string //! 消息
	Time int64  //! 时间
}

type UP_Log struct {
	GameId  string //! 游戏Id
	MsgType int    //! 是否为聊天
	Msg     []byte //! 内容
}

type SDK_Log struct {
	Addr string //! 游戏Id
	Msg  []byte //! 内容
}

type SQL_Set struct {
	Sql string
}

type ServerStatus struct {
	SqlLogLen     int //! 日志长度
	SqlBeLogLen   int //! 日志长度
	SqlLineLogLen int //! 在线日志长度
	UpLogLen      int //! 上传日志
	SDKLogLen     int //! 上传日志
	DisconnectLen int //! 断开连接数量
	BroadCastLen  int //! 广播消息长度
	SaveSqlLen    int //! 写入数据库数据
	SaveBaseLen   int //! 写入userbase数据
	LinePlayer    int //! 排队人数
	SessionNum    int //! 连接数
	OnlinePlayer  int //! 在线人数
	AllPlayer     int //! 缓存人数
	WaitLock      int //! 等待锁
}

const (
	MAX_SQL_NUM = 200000 //! 最长的SQL长度
)

/////////////////////////////////////////////////////////////////////////////////////
type Server struct {
	Con      *Config         //! 配置
	Wait     *sync.WaitGroup //! 同步阻塞
	ShutDown bool            //! 是否正在执行关闭
	//IsMaster bool            //! 中心服务器
	Locker *sync.Mutex //! 共享锁
	//! 数据库
	DBUser *DBServer //！database 接口
	DBLog  *DBServer //! database log 接口
	Redis  *redis.Pool

	//! 服务器状态
	Status ServerStatus
	//! 日志数据
	LogChan        chan *UP_Log      //! 经分上传消息
	SDKLogChan     chan *SDK_Log     //! SDK上传
	SqlLogChan     chan *SQL_Log     //! 游戏道具日志
	SqlBeLogChan   chan *SQL_BeLog   //！游戏行为日志
	SqlLineLogChan chan *SQL_LineLog //! 在线数据SQL
	SqlChan        chan *SQL_Set     //! 模块数据SQL
	SqlBaseChan    chan *SQL_Set     //! userbase的缓存
	SqlMailLogChan chan *SQL_MailLog //! 在线数据SQL
	SqlChatLogChan chan *SQL_ChatLog //! 在线聊天SQL

	cache *CacheClient //! 中转
	//! session处理
	AddSession   chan *Session
	DelSession   chan *Session
	BroadCastMsg chan []byte

	Level           int   //! 服务器平均等级
	UpdateForDay    int64 //! 每天更新时间
	OpenServerTime  int64 //! 服务端启动时间
	UserBaseLogTime int64 //! userbase备份时间

	LogConnNum   int
	BeLogConnNum int
	SqlConnNum   int

	//! 敏感词
	SensitiveWord             Trie     //! 敏感词
	SensitiveWordJudgePattern []string //! 敏感词

	Event int

	serverAgent   *ServerAgent //! 连接跨服
	serverChat    *TCPClient   //!
	serverChatMsg []string     //!
}

var serverSingleton *Server = nil

//! 得到服务器指针
func GetServer() *Server {
	if serverSingleton == nil {
		serverSingleton = new(Server)
		serverSingleton.Con = new(Config)
		serverSingleton.Wait = new(sync.WaitGroup)
		serverSingleton.Locker = new(sync.Mutex)
		serverSingleton.DBUser = new(DBServer)
		serverSingleton.DBLog = new(DBServer)
		serverSingleton.LogChan = make(chan *UP_Log, 10000)
		serverSingleton.SDKLogChan = make(chan *SDK_Log, 10000)
		serverSingleton.SqlLogChan = make(chan *SQL_Log, MAX_SQL_NUM)
		serverSingleton.SqlBeLogChan = make(chan *SQL_BeLog, MAX_SQL_NUM)
		serverSingleton.SqlLineLogChan = make(chan *SQL_LineLog, MAX_SQL_NUM)
		serverSingleton.SqlChan = make(chan *SQL_Set, MAX_SQL_NUM)
		serverSingleton.SqlBaseChan = make(chan *SQL_Set, MAX_SQL_NUM)
		serverSingleton.SqlMailLogChan = make(chan *SQL_MailLog, MAX_SQL_NUM)
		serverSingleton.SqlChatLogChan = make(chan *SQL_ChatLog, MAX_SQL_NUM)

		serverSingleton.LogConnNum = 20
		serverSingleton.BeLogConnNum = 6
		serverSingleton.SqlConnNum = 6

		serverSingleton.cache = nil
		serverSingleton.OpenServerTime = TimeServer().Unix()
		serverSingleton.UserBaseLogTime = HF_GetNextTimeToLog()

		serverSingleton.AddSession = make(chan *Session, 10000)
		serverSingleton.DelSession = make(chan *Session, 10000)
		serverSingleton.BroadCastMsg = make(chan []byte, 5000)
		serverSingleton.UpdateForDay = 0

		//serverSingleton.IsMaster = false
		serverSingleton.serverChat = nil

		serverSingleton.SensitiveWord = NewTrie()
	}

	return serverSingleton
}

//! 载入配置文件
func (self *Server) InitConfig() {
	configFile, err := ioutil.ReadFile("./config.json") ///尝试打开配置文件
	if err != nil {
		log.Fatal("config err 1")
	}
	err = json.Unmarshal(configFile, self.Con)
	if err != nil {
		log.Fatal("server InitConfig err:", err.Error())
	}

	if serverSingleton.Con.ServerVer == 0 {
		serverSingleton.Con.ServerVer = 1000
	}

	GetLogMgr().SetLevel(self.Con.LogCon.LogLevel, self.Con.LogCon.LogConsole)
	if self.Con.GameID == "" {
		self.Con.GameID = "10000008"
	}
	if self.Con.AppKey == "" {
		self.Con.AppKey = "158df75271d6439abba7870df6b3c8a2"
	}

	if self.Con.GM == 0 {
		self.Con.GM = GM_ACCOUNT_ID
	}

	if self.Con.LogCon.MaxFileSize == 0 {
		self.Con.LogCon.MaxFileSize = 500
	}

	if self.Con.LogCon.MaxFileNum == 0 {
		self.Con.LogCon.MaxFileNum = 50
	}

	if serverSingleton.Con.NetworkCon.MaxPlayer == 0 {
		serverSingleton.Con.NetworkCon.MaxPlayer = 2000
	}

	// 默认五分钟
	if self.Con.ServerExtCon.UpdateTime == 0 {
		self.Con.ServerExtCon.UpdateTime = 5
	}
}

func (self *Server) ReloadConfig() {
	configFile, err := ioutil.ReadFile("./config.json") ///尝试打开配置文件
	if err != nil {
		log.Fatal("config err 1")
	}
	err = json.Unmarshal(configFile, self.Con)
	if err != nil {
		log.Fatal("server ReloadConfig config err,", err.Error())
	}
	//! 修改日志配置
	GetLogMgr().level = self.Con.LogCon.LogLevel
	GetLogMgr().showStd = self.Con.LogCon.LogConsole
	if self.Con.GameID == "" {
		self.Con.GameID = "10000008"
	}
}

func (self *Server) GetConfig() *Config {
	return self.Con
}

func (self *Server) GetSensitiveWord() *Trie {
	return &self.SensitiveWord
}

//! 连接数据库
func (self *Server) ConnectDB() {
	//! 连接redis
	self.Redis = NewPool(self.Con.DBCon.Redis, self.Con.DBCon.RedisDB, self.Con.DBCon.RedisAuth)
	if self.Redis == nil {
		log.Fatal("redis err")
		return
	}
	self.DBUser.Init(self.Con.DBCon.DBUser)
	self.DBLog.Init(self.Con.DBCon.DBLog)

	// 在加载数据之前, 检查数据库字段
	GetSqlMgr().CheckMysql()

	self.GetLevel(true)

	//! 读入全局数据
	GetHeroSupportMgr().GetData()
	//!获取活动，并且更新数据
	GetActivityMgr().GetData()
	GetActivityMgr().UpdateActivityStatus(true)
	//! 七天任务
	GetWeekPlanMgr().GetData()
	GetMailMgr().GetData()
	GetPassRecordMgr().GetData()
	GetOfflineInfoMgr().GetData()
	GetUnionMgr().GetData()
	GetHireHeroInfoMgr().GetData()

	//! 建筑数据获取
	GetTopBuildMgr().GetData()
	//! 读取红包信息
	GetRedPacMgr().GetData()
	//! 读取转盘消息
	GetDialMgr().GetData()
	//国战战报获取数据
	//GetFightMgr().GetData()   不再获取

	GetArenaMgr().GetData()

	// 支援英雄
	GetSupportHeroMgr().GetData()

	GetRankTaskMgr().GetData()

	GetTowerMgr().GetData()
	GetArenaSpecialMgr().GetData()
	GetAccessCardRecordMgr().GetData()
	GetWechatWarningMgr().GetData()

	//! 初始化RPC 中心服务器
	GetMasterMgr().InitService()

	LogInfo("读取数据OK")
}

//! 得到一个websocket处理句柄
func (self *Server) GetConnectHandler() websocket.Handler {
	connectHandler := func(ws *websocket.Conn) {
		if self.ShutDown {
			return
		}
		session := GetSessionMgr().GetNewSession(ws)
		if session == nil {
			return
		}
		LogDebug("add session:", session.ID)
		session.Run()
	}
	return websocket.Handler(connectHandler)
}

//! 连接中转服务器
func (self *Server) ConnectCacheServer() bool {
	if self.cache == nil {
		self.cache = new(CacheClient)
	}
	if self.cache.InitSocket(self.Con.ServerExtCon.Cache) {
		return true
	} else {
		return false
	}
}

//! run
func (self *Server) Run() {
	//self.GoService()

	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			self.onTime()
		case session := <-self.AddSession:
			GetSessionMgr().MapSession[session.ID] = session
		case session := <-self.DelSession:
			delete(GetSessionMgr().MapSession, session.ID)
			//case msg := <-self.BroadCastMsg:
			//	GetSessionMgr().OrderBroadCastMsg("1", msg)
		}
	}

	ticker.Stop()
}

//! 开启本地服务
func (self *Server) GoService() {
	log.Println("主程序启动， 程序版本：", self.Con.ServerVer)

	//! 游戏主循环
	go self.Run()

	//! 道具日志
	for i := 0; i < self.LogConnNum; i++ {
		go self.RunSqlLog()
	}

	//! 行为日志
	for i := 0; i < self.BeLogConnNum; i++ {
		go self.RunSqlSet()
		go self.RunSqlBeLog()
	}

	//! 中心服逻辑
	go GetMasterMgr().OnTimer()
	go GetMasterMgr().OnLogic()

	//! 在线状态记录
	go self.RunSqlLineLog()

	//! 邮件领取
	go self.RunSqlMailLog()

	//! 在线聊天
	go self.RunSqlChatLog()

	//! 数据保存
	go self.RunSqlBaseSet()

	go self.RunSDKLog()

	//! 排行榜逻辑
	go GetTopMgr().Run()
	go GetRankRewardMgr().Run()
	//! 收藏夹线程
	go GetAccessCardRecordMgr().Run()
	//! 活动BOSS线程
	go GetActivityBossMgr().Run()

	go GetLineUpMgr().Run()
	go GetSessionMgr().Run()
	go GetRechargeMgr().Run()
	go GetLoggerCheckMgr().Run()
	go GetSdkLoggerCheckMgr().Run()
	go self.ServerAgentRun()
	go self.ChatAgentRun()
	go GetBackStageMgr().RunGetNotice()
	go GetBackStageMgr().RunCheckNotice()

	go GetArenaMgr().StartFight()
	go GetArenaSpecialMgr().StartFight()
	go GetArenaSpecialMgr().CountAward()

	go GetGeneralMgr().OnTimer()
	go GetCrossArenaMgr().OnTimer()
	go GetCrossArena3V3Mgr().OnTimer()
	go GetConsumerTop().OnTimer()

	if self.Con.ServerExtCon.NumRecord != "" { //! 上报人数
		go self.NumRecord()
	}

	go GetRedPacMgr().OnTimer()

	go self.LoadSensitiveWord()

	//! 迁移Redis数据
	go GetPassRecordMgr().RunMigratePassRecord()
	go GetArenaMgr().RunMigrateArena()
	go GetArenaSpecialMgr().RunMigrateArenaSpecial()

}

func (self *Server) onTime() {
	if self.ShutDown {
		return
	}

	tNow := TimeServer()
	if tNow.Minute()%6 == 3 && tNow.Second() == 0 { //! 每6分钟
		self.Save()
	}

	if tNow.Hour() == 5 && tNow.Minute() < 5 && tNow.Second() > 0 && tNow.Unix()-self.UpdateForDay > 1800 {
		//!每天5：00检测，
		LogInfo("每天5点刷新")
		self.UpdateForDay = tNow.Unix()
		//排行榜活动更新-优先更新活动，再更新其他相关的模块
		GetActivityMgr().UpdateActivityStatus(true)
		GetActivityMgr().Save()

		GetTopBuildMgr().Refresh()
		GetTopBoxMgr().Refresh()
		//GetUnionMgr().OnRefresh()
		//GetUnionMgr().OnHunterRefresh()
		GetArenaMgr().CheckArenaEnd()
		GetArenaSpecialMgr().CheckArenaEnd()

		//! 州宝箱系统屏蔽
		//GetUnionFightMgr().CheckMail()

		self.GetLevel(true)

		GetConsumerTop().OnFresh()
	}

	if tNow.Hour() == 6 && tNow.Minute() < 5 && tNow.Unix()-self.UpdateForDay > 1800 {
		//!每天6：00检测，
		LogInfo("每天6点刷新")
		self.UpdateForDay = tNow.Unix()
		//GetExpeditionMgr().InitExpeditionFight()
	}

	if tNow.Hour() == 21 && tNow.Minute() < 5 && tNow.Unix()-self.UpdateForDay > 3600 {
		self.UpdateForDay = tNow.Unix()
	}

	if tNow.Minute()%5 == 0 && tNow.Second() == 0 { // 每5分钟统计在线人数
		self.sendLog_onTime()
		//self.SendLog_SDKUP_Online()
	}

	if tNow.Minute()%30 == 0 && tNow.Second() == 0 { // 每30分钟请求一次屏蔽字库
		self.LoadSensitiveWord()
	}

	if tNow.Second() == 0 { // 每分钟发送一次聊天记录
		self.SendMsgChatNew()
	}

	//写LOG  20201112
	if tNow.Unix() > self.UserBaseLogTime {
		var userbase San_UserBase
		stardTime := TimeServer().Unix() - DAY_SECS
		timeStr := time.Unix(stardTime, 0).Format(DATEFORMAT)
		sql := fmt.Sprintf("select * from `san_userbase` where lastlogintime > '%s'", timeStr)
		res := GetServer().DBUser.GetAllData(sql, &userbase)
		for i := 0; i < len(res); i++ {
			InsertTable("san_userbase2", res[i], 0, false)
		}
		self.UserBaseLogTime = HF_GetNextTimeToLog()
	}

	//! 更新在线人数
	if tNow.Second() == 0 {
		GetPlayerMgr().UpdateOnline()
	}

	GetSessionMgr().ClearRemoveSession()
}

//! 每小时上报人数
func (self *Server) NumRecord() {
	ticker := time.NewTicker(time.Second * 3600)
	for {
		<-ticker.C
		str := fmt.Sprintf("http://%s?m=Doc&c=period&a=count&db_id=%d&num=%d", self.Con.ServerExtCon.NumRecord,
			self.Con.DBCon.DBId, GetSessionMgr().GetSessionNum())
		res, _ := http.Get(str)
		if res != nil {
			res.Body.Close()
		}
	}

	ticker.Stop()
}

//! 关服
func (self *Server) Close() {
	self.ShutDown = true
	LogInfo("设置关服标志量ok")

	LogInfo("剩余sqllog:", len(self.SqlLogChan))
	LogInfo("剩余sqlbelog:", len(self.SqlBeLogChan))
	LogInfo("剩余sqllinelog:", len(self.SqlLineLogChan))
	LogInfo("剩余sqlset:", len(self.SqlChan))
	LogInfo("剩余sqlbase:", len(self.SqlBaseChan))
	LogInfo("剩余数据上报:", len(self.LogChan))
	LogInfo("剩余event:", self.Event)
	LogInfo("剩余sqlmaillog:", len(self.SqlMailLogChan))
	LogInfo("剩余chatlog:", len(self.SqlChatLogChan))

	LogInfo("设置关服标志量ok")

	for i := 0; i < self.LogConnNum; i++ {
		//! 关闭操作日志
		self.SqlLog(0, DEFAULT_GOLD, 0, 0, 0, "", 0, 0, nil)
	}

	LogInfo("关闭sql日志ok")

	//self.Log(0, []byte("close"))
	//LogInfo("关闭数据上报ok")
	//! 保存全部玩家信息
	GetPlayerMgr().SaveAll(true)
	//! 保存全局函数
	self.Save()

	self.SqlLineLog(0, "", 0, "")
	for i := 0; i < self.BeLogConnNum; i++ {
		self.SqlSet("")
		//! 关闭道具日志
		self.SqlLog(0, 0, 0, 0, 0, "", 0, 0, nil)
	}

	self.SqlBaseSet("")
	self.CloseMailLog()
	self.CloseChatLog()
	self.CloseSDKLog()
	//self.Log(0, []byte(""))
	//if self.LogChan != nil {
	//	if len(self.LogChan) > 0 { //!如果满了，则清空
	//		if self.cache == nil || self.cache.Dead == true {
	//			LogInfo("上报日志出错，清空缓存")
	//			//close(self.LogChan)
	//			//self.LogChan = nil
	//		}
	//	}
	//	LogInfo("关闭数据上报ok")
	//	self.LogChan <- &UP_Log{"0", 0, []byte("")}
	//}

	LogInfo("剩余sqllog:", len(self.SqlLogChan))
	LogInfo("剩余sqlbelog:", len(self.SqlBeLogChan))
	LogInfo("剩余sqllinelog:", len(self.SqlLineLogChan))
	LogInfo("剩余sqlset:", len(self.SqlChan))
	LogInfo("剩余sqlbase:", len(self.SqlBaseChan))
	LogInfo("剩余数据上报:", len(self.LogChan))
	LogInfo("剩余event:", self.Event)
	LogInfo("剩余sqlmaillog:", len(self.SqlMailLogChan))
	LogInfo("剩余chatlog:", len(self.SqlChatLogChan))

	LogInfo("关闭数据库ok")

	// 这个地方会出现阻塞
	self.Wait.Wait()
	self.DBUser.Close()
	self.DBLog.Close()

	if self.SqlBeLogChan != nil {
		close(self.SqlBeLogChan)
		self.SqlBeLogChan = nil
		LogInfo("SqlBeLogChan 完成关闭")
	}

	if self.SqlChan != nil {
		close(self.SqlChan)
		self.SqlChan = nil
		LogInfo("SqlChan 完成关闭")
	}

	if self.SqlLogChan != nil {
		close(self.SqlLogChan)
		self.SqlLogChan = nil
		LogInfo("SqlLogChan 完成关闭")
	}

	if self.SqlMailLogChan != nil {
		close(self.SqlMailLogChan)
		self.SqlMailLogChan = nil
		LogInfo("SqlMailLogChan 完成关闭")
	}

	if self.SqlChatLogChan != nil {
		close(self.SqlChatLogChan)
		self.SqlChatLogChan = nil
		LogInfo("SqlChatLogChan 完成关闭")
	}

	LogInfo("服务器完成关闭")
	LogFatal("server shutdown")
}

//! 服务器全局保存
func (self *Server) Save() {
	GetRankRewardMgr().Save()
	LogInfo("保存排行榜OK")
	GetAccessCardRecordMgr().Save()
	LogInfo("保存收藏家ok")
	GetActivityBossMgr().Save()
	LogInfo("保存暗域入侵ok")
	GetCrossArenaMgr().Save()
	LogInfo("保存跨服竞技场OK")
	GetCrossArena3V3Mgr().Save()
	LogInfo("保存33跨服竞技场OK")
	GetUnionMgr().Save()
	LogInfo("保存军团ok")
	GetMailMgr().Save()
	LogInfo("保存邮件信息ok")
	//GetPvpMgr().Save()
	LogInfo("保存兵种天下会武ok")
	GetOfflineInfoMgr().Save()
	LogInfo("离线信息保存ok")
	GetHireHeroInfoMgr().Save()
	LogInfo("雇佣信息保存ok")
	GetTopBuildMgr().SaveData()
	LogInfo("保存内政厅排行榜ok")
	GetPassRecordMgr().Save()
	LogInfo("关卡保存ok")
	GetRedPacMgr().Save()
	LogInfo("服务器红包保存ok")
	GetDialMgr().Save()
	LogInfo("保存活动转盘信息ok")
	GetSupportHeroMgr().Save()
	LogInfo("图书馆ok")
	GetRankTaskMgr().Save()
	LogInfo("排行任务ok")
	GetArenaMgr().Save()
	GetArenaSpecialMgr().Save()
}

//! 得到开服第几天
func (self *Server) GetOpenTime() int {
	t, _ := time.ParseInLocation(DATEFORMAT, self.Con.OpenTime, time.Local)
	_t := TimeServer().Unix() - t.Unix()
	day := int(_t / 86400)
	if _t%86400 >= 0 {
		day++
	}

	return day
}

//! 得到开服时间戳
func (self *Server) GetOpenServer() int64 {
	t, _ := time.ParseInLocation(DATEFORMAT, self.Con.OpenTime, time.Local)
	return t.Unix()
}

//! 得到开服时间-时间格式
func (self *Server) GetOpenTimeInfo() time.Time {
	t, _ := time.ParseInLocation(DATEFORMAT, self.Con.OpenTime, time.Local)
	return t
}

//! 得到开服时间戳
func (self *Server) GetOpenDay() int64 {
	t, _ := time.ParseInLocation(DATEFORMAT, self.Con.OpenTime, time.Local)
	timeSet := time.Date(t.Year(), t.Month(), t.Day(), 5, 0, 0, 0, time.Local)
	if t.Unix() > timeSet.Unix() {
		return timeSet.Unix() + int64(DAY_SECS)
	} else {
		return timeSet.Unix()
	}
}

//! 得到开服第几天的5：00
func (self *Server) GetDayTime(openday int) int64 {
	t, _ := time.ParseInLocation(DATEFORMAT, self.Con.OpenTime, time.Local)
	return int64(time.Date(t.Year(), t.Month(), t.Day(), 5, 0, 0, 0, t.Location()).Unix() + int64(86400*openday))
}

//! 得到服务器平均等级
func (self *Server) GetLevel(refresh bool) int {
	if refresh {
		total := 0
		var num San_Num
		res := GetServer().DBUser.GetAllData("select `level` as num from san_userbase order by `level` desc limit 30", &num)
		for i := 0; i < len(res); i++ {
			total += res[i].(*San_Num).Num
		}

		self.Level = HF_MinInt(90, HF_MaxInt(total/30-2, 10))
	}

	return self.Level
}

//! PvPId
func (self *Server) GetPvPId() int64 {
	return TimeServer().UnixNano() / 10000
}

// 全服发送系统消息
func (self *Server) sendSysChat(contend string) {
	var msg S2C_Chat
	msg.Cid = "chat"
	msg.Uid = 0
	msg.Channel = CHAT_SYSTEM
	msg.Name = GetCsvMgr().GetText("STR_SYS_MSG")
	msg.Icon = 0
	msg.Vip = 0
	msg.Time = TimeServer().Unix()
	msg.Content = contend
	msg.Url = ""
	msg.Param = 0
	GetSessionMgr().BroadCastMsg("chat", HF_JtoB(&msg))
}

func (self *Server) sendWorldChat(content string) {
	var msg S2C_Chat
	msg.Cid = "chat"
	msg.Uid = 0
	msg.Channel = CHAT_WORLD
	msg.Name = GetCsvMgr().GetText("STR_SYS_MSG")
	msg.Icon = 0
	msg.Vip = 0
	msg.Time = TimeServer().Unix()
	msg.Content = content
	msg.Url = ""
	msg.Camp = 0
	msg.Param = 0
	GetSessionMgr().BroadCastMsg("chat", HF_JtoB(&msg))
}

func (self *Server) sendCampChat(camp int, content string) {
	var msg S2C_Chat
	msg.Cid = "chat"
	msg.Uid = 0
	msg.Channel = CHAT_CAMP
	msg.Name = GetCsvMgr().GetText("STR_SYS_MSG")
	msg.Icon = 0
	msg.Vip = 0
	msg.Time = TimeServer().Unix()
	msg.Content = content
	msg.Url = ""
	msg.Camp = camp
	msg.Param = 0
	GetPlayerMgr().BroadCastMsgToCamp(camp, "chat", HF_JtoB(&msg))
}

func (self *Server) sendUnionChat(sanUnion *San_Union, content string) {
	var msg S2C_Chat
	msg.Cid = "chat"
	msg.Uid = 0
	msg.Channel = CHAT_PARTY
	msg.Name = GetCsvMgr().GetText("STR_SYS_MSG")
	msg.Icon = 0
	msg.Vip = 0
	msg.Time = TimeServer().Unix()
	msg.Content = content
	msg.Url = ""
	msg.Camp = 0
	msg.Param = 0
	sanUnion.BroadCastMsg("chat", HF_JtoB(&msg))
}

//! 全服公告
func (self *Server) Notice(content string, camp int, notype int) { // 内容，阵营(0全阵营)，类型(0全地图，1仅仅大地图，2仅仅战斗场景)

	var msg S2C_Notice
	msg.Cid = "notice"
	msg.Type = notype
	msg.Content = content
	if camp == 0 {
		GetSessionMgr().BroadCastMsg("1", HF_JtoB(&msg))
		GetServer().sendSysChat(content)
	} else {
		GetPlayerMgr().BroadCastMsgToCamp(camp, "1", HF_JtoB(&msg))
		GetServer().sendCampChat(camp, content)
	}
}

//! 新增数据上报按照游戏进行上报
func (self *Server) Log(gameId string, logs []byte) {
	//return
	//log.Println("加入日志:", string(logs))
	if self.ShutDown {
		return
	}

	if self.LogChan == nil {
		return
	}
	if len(self.LogChan) < 10000 {
		self.LogChan <- &UP_Log{GameId: gameId, MsgType: 2, Msg: logs}
	}
}

func (self *Server) SDKLog(addr string, logs []byte) {
	//return
	//log.Println("加入日志:", string(logs))
	if self.ShutDown {
		return
	}

	if self.SDKLogChan == nil {
		return
	}
	if len(self.SDKLogChan) < 10000 {
		self.SDKLogChan <- &SDK_Log{Addr: addr, Msg: logs}
	}
}

func (self *Server) LogChat(gameId string, logs []byte) {
	if self.ShutDown {
		return
	}

	if self.LogChan == nil {
		return
	}
	if len(self.LogChan) < 10000 {
		self.LogChan <- &UP_Log{GameId: gameId, MsgType: 1, Msg: logs}
	}
}

//! 数据上报:阿里经分
func (self *Server) RunLog() {
	self.Wait.Add(1)
	for msg := range self.LogChan {
		if string(msg.Msg) == "" {
			break
		}
		if self.cache == nil || self.cache.Dead == true {
			self.ConnectCacheServer()

			if self.ShutDown == true {
				break
			}
			//break
		}
		for i := 0; i < 10; i++ {
			if self.SendLog2Cache(msg.GameId, fmt.Sprintf("%d", msg.MsgType), msg.Msg) {
				break
			}
		}
	}
	self.Wait.Done()
	log.Println("logchan 上报完毕")
	close(self.LogChan)
	self.LogChan = nil
}

func (self *Server) RunSDKLog() {
	self.Wait.Add(1)
	for msg := range self.SDKLogChan {
		if string(msg.Msg) == "" {
			break
		}
		if self.ShutDown == true {
			break
		}
		body := bytes.NewBuffer(msg.Msg)
		resp, err := http.Post(msg.Addr, "application/json", body)
		if resp == nil {
			log.Println("RunSDKLog 错误resp")
			continue
		}
		if err != nil {
			log.Println("RunSDKLog 错误err")
			continue
		}
		if resp.StatusCode != 200 {
			log.Println("RunSDKLog 错误StatusCode")
			resp.Body.Close()
			continue
		}
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println("result=", string(result))
		resp.Body.Close()
	}
	self.Wait.Done()
	log.Println("sdklogchan 上报完毕")
	close(self.SDKLogChan)
	self.SDKLogChan = nil
}

func (self *Server) GetServerStatus() *ServerStatus {
	self.Status.BroadCastLen = len(self.BroadCastMsg)
	self.Status.UpLogLen = len(self.LogChan)
	self.Status.SDKLogLen = len(self.SDKLogChan)
	self.Status.DisconnectLen = len(self.DelSession)
	self.Status.SqlLogLen = len(self.SqlLogChan)
	self.Status.SqlBeLogLen = len(self.SqlBeLogChan)
	self.Status.SqlLineLogLen = len(self.SqlLineLogChan)
	self.Status.SaveSqlLen = len(self.SqlChan)
	self.Status.SaveBaseLen = len(self.SqlBaseChan)
	self.Status.OnlinePlayer = GetPlayerMgr().GetPlayerOnline()
	self.Status.AllPlayer = self.Event
	self.Status.SessionNum = len(GetSessionMgr().MapSession)
	self.Status.LinePlayer = GetLineUpMgr().WaitTotal
	return &self.Status
}

//! 测试数据上报
func (self *Server) TestRunLog() {
	//self.Wait.Add(1)
	i := 0
	for i < 10 {
		time.Sleep(time.Microsecond * 5)
		self.sendLog_onTime()
	}

	//self.Wait.Done()
	//close(self.LogChan)
}

func (self *Server) SendLog2Cache(gameId string, msgType string, logs []byte) bool {
	LogDebug("上报日志:", string(logs))
	if self.cache != nil {
		self.cache.Sendstring(gameId, msgType, string(logs))
		return true
	}

	return false
}

//! redis更新数据-备份数据
func (self *Server) SqlSet(sql string) {
	if self.SqlChan == nil {
		return
	}
	self.SqlChan <- &SQL_Set{sql}
}

//! 数据保存-多线程
func (self *Server) RunSqlSet() {
	self.Wait.Add(1)
	var flag = false
	for msg := range self.SqlChan {
		if msg.Sql == "" {
			break
		}
		flag = false
		for i := 0; i < 10; i++ {
			LogDebug("SQL BASE EXEC:", msg.Sql)
			_, _, ok := GetServer().DBUser.Exec(msg.Sql)
			if ok {
				flag = true
				break
			}
		}

		if !flag {
			LogError("sql execute failed, sql = ", msg.Sql)
		}
	}
	self.Wait.Done()
	log.Println("sqlset完毕")
	//if self.SqlChan != nil {
	//	close(self.SqlChan)
	//	self.SqlChan = nil
	//}
}

//! sql语句-非Reis数据，数据保存-全局数据
func (self *Server) SqlBaseSet(sql string) {
	if self.SqlBaseChan == nil {
		return
	}
	self.SqlBaseChan <- &SQL_Set{sql}
}

//! 全局数据保存
func (self *Server) RunSqlBaseSet() {
	self.Wait.Add(1)
	for msg := range self.SqlBaseChan {
		if msg.Sql == "" {
			break
		}
		LogDebug("SQL BASE EXEC:", msg.Sql)
		for i := 0; i < 10; i++ {
			_, _, ok := GetServer().DBUser.Exec(msg.Sql)
			if ok {
				break
			}
		}
	}
	self.Wait.Done()
	log.Println("sqlbaseset完毕")
	close(self.SqlBaseChan)
	self.SqlBaseChan = nil
}

//! 在线日志
func (self *Server) SqlLineLog(uid int64, ip string, line int, creator string) {
	if self.SqlLineLogChan == nil {
		//LogError("self.SqlLineLogChan == nil!!!")
		return
	}

	if len(self.SqlLineLogChan) >= 200000 {
		LogError("丢了一个日志")
		return
	}

	self.SqlLineLogChan <- &SQL_LineLog{0, uid, TimeServer().Unix(), ip, line, creator}
}

//! 在线数据日志
func (self *Server) RunSqlLineLog() {
	//return
	self.Wait.Add(1)
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
	self.Wait.Done()
	log.Println("sqllinelog完毕")
	if self.SqlLineLogChan != nil {
		close(self.SqlLineLogChan)
		self.SqlLineLogChan = nil
	}
}

//! 数据库日志
// 玩家id，道具id，增加/减少值，变量1(系统出处)，变量2，描述
func (self *Server) SqlLog(uid int64, _type int, value int, param1 int, param2 int, dec string, cur int, param3 int, player *Player, args ...int) {
	if _type/10000000 > 0 {
		if self.SqlLogChan == nil {
			LogInfo("self.SqlLogChan == nil")
			return
		}
		if len(self.SqlLogChan) >= 200000 {
			LogError("丢了一个日志")
			return
		}
		AddSdkItemLog(uid, _type, value, param1, param2, dec, cur, param3, player)
		costGem := 0
		if len(args) > 0 {
			costGem = args[0]
		}
		self.SqlLogChan <- &SQL_Log{0, TimeServer().Unix(), _type, value, param1, param2, uid, dec, cur, costGem}
		if uid == 0 {
			//LogInfo("向管道SqlLogChan发送了停止消息")
		}
	} else {
		if self.SqlBeLogChan == nil {
			LogInfo("self.SqlBeLogChan == nil")
			return
		}
		if len(self.SqlBeLogChan) >= 200000 {
			LogError("丢了一个日志")
			return
		}
		if player != nil {
			self.SqlBeLogChan <- &SQL_BeLog{0, TimeServer().Unix(), _type, value, param1,
				param2, uid, dec, cur, param3, player.Sql_UserBase.Level,
				player.Sql_UserBase.Vip, int(player.Sql_UserBase.Fight / 100), player.Sql_UserBase.PassMax}
		} else {
			self.SqlBeLogChan <- &SQL_BeLog{0, TimeServer().Unix(), _type, value, param1,
				param2, uid, dec, cur, param3, 0, 0, 0, 0}
		}
	}
}

// 玩家id，道具id，增加/减少值，变量1(系统出处)，变量2，描述
func (self *Server) SqlLogEx(uid int64, _type int, value int, param1 int64, param2 int, dec string, cur int, param3 int, fightinfo *JS_FightInfo) {
	if _type/10000000 > 0 {
		if self.SqlLogChan == nil {
			return
		}
		if len(self.SqlLogChan) >= 200000 {
			LogError("丢了一个日志")
			return
		}
		self.SqlLogChan <- &SQL_Log{0, TimeServer().Unix(), _type, value, int(param1),
			param2, uid, dec, cur, param3}
	} else {
		if self.SqlBeLogChan == nil {
			return
		}
		if len(self.SqlBeLogChan) >= 200000 {
			LogError("丢了一个日志")
			return
		}
		self.SqlBeLogChan <- &SQL_BeLog{0, TimeServer().Unix(), _type, value, int(param1),
			param2, uid, dec, cur, param3, fightinfo.Level, fightinfo.Vip, int(fightinfo.Deffight / 100), 0}
	}
}

//! 游戏日志保存[没有完全关闭]
func (self *Server) RunSqlLog() {
	self.Wait.Add(1)
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
	self.Wait.Done()
	log.Println("sqllog完毕")
}

func (self *Server) CloseMailLog() {
	if self.SqlMailLogChan == nil {
		return
	}
	self.SqlMailLogChan <- &SQL_MailLog{0, 0, "", "", ""}
}

//! 邮件日志保存
func (self *Server) RunSqlMailLog() {
	self.Wait.Add(1)
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
	self.Wait.Done()
	log.Println("sql maillog完毕")
}

//! 游戏-操作日志-保存
func (self *Server) RunSqlBeLog() {
	self.Wait.Add(1)
	for msg := range self.SqlBeLogChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 10; i++ {
			if InsertLogTable("san_belog", msg, 1) > 0 {
				break
			}
		}
	}
	self.Wait.Done()
	log.Println("san_belog完毕")
	//if self.SqlBeLogChan != nil {
	//	close(self.SqlBeLogChan)
	//	self.SqlBeLogChan = nil
	//}
}

func (self *Server) GetRedisConn() redis.Conn {
	return self.Redis.Get()
}

func (self *Server) IsBlackIP(ip string) bool {
	for i := 0; i < len(self.Con.NetworkCon.BlackIP); i++ {
		if self.Con.NetworkCon.BlackIP[i] == ip {
			return true
		}
	}

	return false
}

func (self *Server) IsWhiteID(uid int64) bool {
	for i := 0; i < len(self.Con.NetworkCon.WhiteID); i++ {
		if self.Con.NetworkCon.WhiteID[i] == uid {
			return true
		}
	}

	return false
}

type SensitiveWord_Msg struct {
	Words []string "words"
}

type SensitiveWord_Msg_New struct {
	Code       int      "code"
	Data       []string "data"
	Updated_at int      "updated_at"
}

func (self *Server) LoadSensitiveWord() {
	LogInfo("请求黑词数据...")

	url := self.GetSensitiveAddr()
	res, err := http.Get(url)
	if err != nil {
		log.Println("请求敏感词失败:", err)
		return
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println("敏感词请求失败:", err)
		return
	}
	//LogDebug("result=", string(result))

	//var msg SensitiveWord_Msg
	var msg SensitiveWord_Msg_New
	json.Unmarshal(result, &msg)
	if msg.Updated_at == 0 {
		return
	}
	if msg.Updated_at != self.SensitiveWord.Ver {
		LogDebug("同步屏蔽词库完成")
		self.SensitiveWord.Mu.Lock()
		defer self.SensitiveWord.Mu.Unlock()

		self.SensitiveWordJudgePattern = make([]string, 0)
		self.SensitiveWord.Root = NewTrieNode()
		self.SensitiveWord.Ver = msg.Updated_at
		for i := 0; i < len(msg.Data); i++ {
			temp := []rune(msg.Data[i])
			if len(temp) > 0 && string(temp[0]) == "[" {
				self.SensitiveWordJudgePattern = append(self.SensitiveWordJudgePattern, msg.Data[i])
				continue
			}
			self.SensitiveWord.Inster(msg.Data[i])
		}

		GetServer().IsSensitiveWordTest()
	}
}

type BarrageLog struct {
	Uid     int64  `json:"uid"`
	Uname   string `json:"uname"`
	Context string `json:"context"`
}

// 插入 redis
func (self *Server) InsertBarrage(uid int64, uname string, context string) {
	c := GetServer().GetRedisConn()
	defer c.Close()

	//查询总长度
	lens, _ := c.Do("LLEN", "barragelist")

	if lens.(int64) >= 2000 {
		c.Do("Rpop", "barragelist")
	}

	var logitem BarrageLog
	logitem.Uid = uid
	logitem.Uname = uname
	logitem.Context = context
	_, err := c.Do("LPUSH", "barragelist", HF_JtoA(&logitem))
	if err != nil {
		LogError("redis fail! query:", "SET", ",err:", err)
	}
}

type BarrageMsgs struct {
	Msglst []string `json:"msglst"`
}

// 读取redis 列表
func (self *Server) GetBarrage(size int, page int) []string {
	c := GetServer().GetRedisConn()
	defer c.Close()

	start := (page - 1) * size
	end := page * size
	values, err := redis.Values(c.Do("LRANGE", "barragelist", start, end))
	//
	if err != nil {
		LogError("redis fail! query:", "get", ",err:", err)
	}
	dmlist := make([]string, 0)
	for _, v := range values {
		dmlist = append(dmlist, string(v.([]byte)))
		//fmt.Println(string(v.([]byte)))
	}
	return dmlist
}

func (self *Server) NoticeInfo(data []byte, camp int, notype int) { // 内容，阵营(0全阵营)，类型(0全地图，1仅仅大地图，2仅仅战斗场景)
	if camp == 0 {
		GetSessionMgr().BroadCastMsg("1", data)
	} else {
		GetPlayerMgr().BroadCastMsgToCamp(camp, "1", data)
	}
}

// 按照不同的平台来划分gameId
func (self *Config) GetGameIdByAppId(appId string) string {
	if appId == "" {
		return self.GetAndroidGameId(appId)
	} else {
		return self.GetIosGameId(appId)
	}
}

// Ios专用
func (self *Config) GetAppKeyByAppId(appId string) string {
	if len(self.AppCon) <= 0 {
		LogError("len(self.AppConfig) <= 0")
		return "6d0f09494f4e2935ceb971f4048db965"
	}

	for _, v := range self.AppCon {
		if v.AppId == appId {
			return v.AppKey
		}
	}
	return self.AppCon[0].AppKey
}

func (self *Config) GetCheckUrlByAppId(appId string) string {
	if len(self.AppCon) <= 0 {
		LogError("len(self.AppConfig) <= 0")
		return "https://aliapi.1tsdk.com/api/v7/cp/user/check?"
	}

	if HF_Atoi(appId) > 60120 {
		return "https://union.huoyx.cn/api/v7/cp/user/check?"
	} else {
		return "https://aliapi.1tsdk.com/api/v7/cp/user/check?"
	}
}

// 获取安卓服对应的gameId
func (self *Config) GetAndroidGameId(appId string) string {
	return self.GameID
}

// 获取Ios服对应的gameId
func (self *Config) GetIosGameId(appId string) string {
	if len(self.AppCon) <= 0 {
		LogError("len(self.AppConfig) <= 0")
		return "10000008"
	}

	for _, v := range self.AppCon {
		if v.AppId == appId {
			return v.GameId
		}
	}
	return self.AppCon[0].GameId
}

// 获取角色ServerId

//! 连接跨服
func (self *Server) ConnectCenterServer() bool {
	if self.serverAgent == nil {
		self.serverAgent = new(ServerAgent)
	}

	res := self.serverAgent.InitSocket(self.Con.ServerExtCon.MasterSvr)

	return res
}

// 发送消息到center服务器
func (self *Server) SendMsg2Center(msg []byte) {
	if GetServer().Con.ServerExtCon.IsMaster {
		return
	}

	if self.ShutDown {
		return
	}

	if self.serverAgent == nil || self.serverAgent.Dead == true {
		return
	}

	info := HF_DecodeCenterMsg("center", msg)

	self.serverAgent.SendMsg(info)
}

//! run
func (self *Server) ServerAgentRun() {
	if GetServer().Con.ServerExtCon.IsMaster {
		return
	}

	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		if self.ShutDown {
			return
		}

		if self.serverAgent == nil || self.serverAgent.Dead == true {
			self.ConnectCenterServer()
		}

	}

	ticker.Stop()
}

//! 连接跨服
func (self *Server) ConnectChatServer() bool {
	if self.serverChat == nil {
		self.serverChat = new(TCPClient)
	}

	res := self.serverChat.InitSocket(self.Con.ServerExtCon.MasterSvr)

	return res
}

// 发送消息到center服务器
func (self *Server) SendMsg2Chat(msg []byte) {
	//if GetServer().Con.ServerExtCon.IsMaster {
	//	return
	//}
	if self.ShutDown {
		return
	}

	if self.serverChat == nil {
		return
	}

	if len(self.serverChatMsg) >= 1000 {
		return
	}
	self.serverChatMsg = append(self.serverChatMsg, string(msg))

	//self.serverChat.Send(string(msg))
}

//! run
func (self *Server) ChatAgentRun() {
	//if GetServer().Con.ServerExtCon.IsMaster {
	//	return
	//}

	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		if self.ShutDown {
			return
		}

		if self.serverChat == nil || self.serverChat.Dead == true {
			if self.serverChat != nil && self.serverChat.Conn != nil {
				self.serverChat.Conn.Close()
			}

			LogDebug("connect chat server ..")
			self.ConnectChatServer()
		}
	}

	ticker.Stop()
}

//! 在线日志
func (self *Server) SqlMailLog(uid int64, msg string, item string) {
	if self.SqlMailLogChan == nil {
		return
	}

	if len(self.SqlMailLogChan) >= 200000 {
		LogError("丢了一个日志")
		return
	}

	self.SqlMailLogChan <- &SQL_MailLog{0, uid, msg, item, TimeServer().Format(DATEFORMAT)}
}

//---------------------------Chat Log-----------------------------------------------------------------------------------
//! 在线日志
func (self *Server) SqlChatLog(uid int64, msg string, chatType int) {
	if self.SqlChatLogChan == nil {
		return
	}

	if len(self.SqlChatLogChan) >= MAX_SQL_NUM {
		LogError("丢了一个日志")
		return
	}

	self.SqlChatLogChan <- &SQL_ChatLog{0, uid, chatType, msg, TimeServer().Unix()}
}

func (self *Server) CloseChatLog() {
	if self.SqlChatLogChan == nil {
		return
	}
	self.SqlChatLogChan <- &SQL_ChatLog{0, 0, 0, "", 0}
}

func (self *Server) CloseSDKLog() {
	if self.SDKLogChan == nil {
		return
	}
	self.SDKLogChan <- &SDK_Log{"", []byte("")}
}

//! 聊天日志保存
func (self *Server) RunSqlChatLog() {
	self.Wait.Add(1)
	for msg := range self.SqlChatLogChan {
		if msg.Uid == 0 {
			break
		}
		for i := 0; i < 10; i++ {
			if InsertLogTable("san_chatlog", msg, 1) > 0 {
				break
			}
		}
	}
	self.Wait.Done()
	log.Println("sql chatlog完毕")
}

func (self *Server) GetSensitiveAddr() string {
	str := fmt.Sprintf("http://laucher.s15.q-dazzle.com/Dragon/index.php?m=api&c=options&a=show&key=ban_words&chk_at=%d", self.SensitiveWord.Ver)
	return str
}

func (self *Server) SendMsgChatNew() {

	if len(self.serverChatMsg) == 0 {
		return
	}

	self.serverChat.InitSocket("")

	endIndex := -1
	for i := 0; i < len(self.serverChatMsg); i++ {
		if self.serverChat.Conn == nil {
			break
		}
		blen, err := self.serverChat.Conn.Write([]byte(self.serverChatMsg[i]))
		if err != nil {
			fmt.Println("err = ", err.Error())
			endIndex = i
			break
		}
		LogDebug("blen = %d", blen)
	}
	if endIndex >= 0 {
		self.serverChatMsg = self.serverChatMsg[endIndex:]
	} else {
		self.serverChatMsg = make([]string, 0)
	}
}

func (self *Server) IsSensitiveWord(txt string) bool {
	if len(txt) < 1 {
		return true
	}

	for _, v := range self.SensitiveWordJudgePattern {
		match, _ := regexp.MatchString(v, txt)
		if match {
			return true
		}
	}
	return GetServer().SensitiveWord.IsSensitiveWord(txt)
}

func (self *Server) IsSensitiveWordTest() {
	return

	str := "[0123456789①②③④⑤⑥⑦⑧⑨⑩⑴⑵⑶⑷⑸⑹⑺⑻⑼⑽⑾⑿⒀⒁⒂⒃⒄⒅⒆⒇⒈⒉⒊⒋⒌⒍⒎⒏⒐⒑⒒⒓⒔⒕⒖⒗⒘⒙⒚⒛〡ニミヨョローㄙㄧ㈠㈡㈢㈣㈤㈥㈦㈧㈨㈩㎎㎏㎜㎝㎞㎡㏄㏎㏑㏒㏕０１２３４５６７８９＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫⅰⅱⅲⅳⅴⅵⅶⅷⅸⅹ℃℅℉№€ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρστυφχψωЁАБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдежзийклмнопрстуфхцчшщъыьэюяёabcdefghijklmnopqrstuvwxyzaoàáèéêìíDòó×ùúüYTàáaèéêìíeòó÷ùúüytāāēēěěīīńňōōūū∥ǎǎǐǐǒǒǔǔǖǖǘǘǚǚǜǜɑɡ=@ABCDEFGHIJKLMNOPQRSTUVWXYZ壹贰叁肆伍陆柒捌玖拾一二三四五六七八九十ᵂˣ]{6,}"
	str1 := "12344"

	match, _ := regexp.MatchString(str, str1)
	print(match)
}
