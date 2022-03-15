package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"code.google.com/p/go.net/websocket"
	"time"
)

//! 发送错误信息
type S2C_ErrInfo struct {
	Cid  string `json:"cid"`
	Info string `json:"info"`
}

/////////////////////////////////////////////////////////////////////////////////////

type HandlerPrototype func(r *http.Request) interface{};

type LoginServer struct {
	Con      *Config //! 配置
	ShutDown bool    //! 是否正在执行关闭
	DBUser *DBServer //！database 接口
	Redis  *redis.Pool
	handler map[string]HandlerPrototype;
}

var loginServerSingleton *LoginServer = nil

//! public
func GetServer() *LoginServer {
	if loginServerSingleton == nil {
		loginServerSingleton = new(LoginServer)
		loginServerSingleton.Con = new(Config)
		loginServerSingleton.DBUser = new(DBServer)
		loginServerSingleton.handler = make(map[string]HandlerPrototype);
	}

	return loginServerSingleton
}

//! 载入配置文件
func (self *LoginServer) InitConfig() {
	configFile, err := ioutil.ReadFile("./config.json") ///尝试打开配置文件
	if err != nil {
		log.Fatal("config err 1")
	}

	err = json.Unmarshal(configFile, self.Con)
	if err != nil {
		log.Fatal("chat InitConfig err", err.Error())
	}

	GetLogMgr().SetLevel(self.Con.LogCon.LogLevel, self.Con.LogCon.LogConsole)
}

func (self *LoginServer) ConnectDB() {
	//! 连接redis
	self.Redis = NewPool(self.Con.DBCon.Redis, self.Con.DBCon.RedisDB, self.Con.DBCon.RedisAuth)
	if self.Redis == nil {
		log.Fatal("redis err")
		return
	}
	self.DBUser.Init(self.Con.DBCon.DBUser)
}

func (self *LoginServer) GetRedisConn() redis.Conn {
	return self.Redis.Get()
}

func (self *LoginServer) Close() {

}

func (self *LoginServer) SqlSet(sql string)  {
	for i := 0; i < 10; i++ {
		_, _, ok := GetServer().DBUser.Exec(sql)
		if ok {
			break
		}
	}
}

func (self *LoginServer) Handler(w http.ResponseWriter, r *http.Request) {

	cmd := r.FormValue("c");
	if (cmd == "") {
		return ;
	}

	handler, ok := self.handler[cmd];
	if (!ok) {
		return ;
	}

	w.Write(HF_JtoB(handler(r)));
}

func (self *LoginServer) RegisterHandler(cmd string, handler HandlerPrototype)  {

	if _, ok := self.handler[cmd]; ok {

	} else {
		self.handler[cmd] = handler;
	}
}

func (self *LoginServer) GetErrInfoResponse(cid string, info string) interface{}  {
	var msg S2C_ErrInfo
	msg.Cid = cid
	msg.Info = info
	return &msg;
}

func (self *LoginServer) GetErrorInfo(ret int) string  {
	return "Error"
}

//! 得到一个websocket处理句柄
func (self *LoginServer) GetConnectHandler() websocket.Handler {
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

//! run
func (self *LoginServer) Run() {

	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			self.onTime()
		}
	}

	ticker.Stop()
}

func (self *LoginServer) onTime() {
	GetSessionMgr().ClearRemoveSession()
}