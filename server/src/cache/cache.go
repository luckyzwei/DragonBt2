package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type LogMsg struct {
	GameId  string
	MsgType int
	MsgBuf  []byte
}

//!消息组
type MsgGroup struct {
	Wait     *sync.WaitGroup
	logChan  chan *LogMsg
	Client   *http.Client
	MaxMsg   [][]byte
	DeadLine int64 // 最后通牒时间
}

//! 初始化消息组
func (self *MsgGroup) initGroup() {
	self.Wait = new(sync.WaitGroup)
	self.logChan = make(chan *LogMsg, GetCacheServer().Con.MaxMsgChan)
	self.Client = nil
	go self.RunLog()
}

//! 获取当前队列消息数量
func (self *MsgGroup) getLogNum() int {
	return len(self.logChan)
}

//! 加入新的日志
func (self *MsgGroup) Log(logs *LogMsg) {
	//log.Println("加入日志:", self.groupId, string(logs.MsgBuf), len(self.logChan))
	self.logChan <- logs
}

//! 数据上报
func (self *MsgGroup) RunLog() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println("runLog:", x, string(debug.Stack()))
			LogDebug("runLog:", x, string(debug.Stack()))
		}
	}()

	self.Wait.Add(1)
	for msg := range self.logChan {
		if string(msg.MsgBuf) == "close" {
			LogDebug("receive close signal!")
			break
		}

		self.MaxMsg = append(self.MaxMsg, msg.MsgBuf)
		now := time.Now().Unix()
		if self.DeadLine == 0 {
			self.DeadLine = now + 10
		}

		timeDiff := self.DeadLine - now
		maxMsg := GetCacheServer().Con.MsgNum
		lenCurrent := len(self.MaxMsg)
		if lenCurrent < maxMsg && timeDiff > 0 {
			continue
		} else {
			// 打包发送
			var b bytes.Buffer
			w := gzip.NewWriter(&b)
			for _, data := range self.MaxMsg {
				// 字节长度
				w.Write(data)
			}
			w.Close()
			for i := 0; i < 3; i++ {
				if self.SendLog(b.Bytes()) {
					break
				}
			}
			GetCacheServer().updateMsgCount(int64(lenCurrent))
			self.MaxMsg = [][]byte{}
			self.DeadLine = now + 10
		}
	}

	self.Wait.Done()
	close(self.logChan)
	LogDebug("RunLog is close chan!")
}

//! 发送普通日志
func (self *MsgGroup) SendLog(logs []byte) bool {
	if GetCacheServer().Con.LogFlag == 1 {
		s, _ := HF_GzipDecode(logs)
		LogDebug("普通日志:", time.Now().Format(TIMEFORMAT), string(s))
	}

	//状态更新
	url := "https://log.rts.dp.uc.cn:8083/api/v2_1"
	body := bytes.NewBuffer(logs)

	if self.Client == nil {
		self.Client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
				DisableCompression: false,
				Dial: func(netw, addr string) (net.Conn, error) {
					deadline := time.Now().Add(6 * time.Second)
					c, err := net.DialTimeout(netw, addr, time.Second*1)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
			},
		}
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("X-dp-app-id", GetCacheServer().Con.AppID)       // 1000024[王牌三国],1000008[血染]
	req.Header.Set("X-dp-api-token", GetCacheServer().Con.AppToken) // 对应的
	req.Header.Set("Connection", "Keep-Alive")

	if GetCacheServer().Con.LogFlag == 1 {
		//log.Println("AppId:", GetCacheServer().Con.AppID, "AppToken:", GetCacheServer().Con.AppToken)
	}
	resp, err := self.Client.Do(req)

	if err != nil {
		log.Println(err)
		return false
	}

	if resp.StatusCode != RESP_CODE_OK {
		if GetCacheServer().Con.LogFlag == 1 {
			log.Println("上报普通日志错误, resp status code:", resp.Status)
		}

		return true

	} else {
		//if GetCacheServer().Con.LogFlag == 1 {
		//	log.Println("上报普通日志成功!")
		//}
		return true
	}

	return false
}

//! 配置
type Config struct {
	Host         string       `json:"wshost"`       //! 服务器
	ServerId     int          `json:"serverid"`     //! 服
	ServerName   string       `json:"servername"`   //! 服务器名称
	InitGroupNum int          `json:"initgroupnum"` //! 初始消息组
	AddGroupNum  int          `json:"addgroupnum"`  //! 每次增加消息组
	MaxMsgChan   int          `json:"maxmsgchan"`   //! 最大消息队列
	AppID        string       `json:"appid"`
	AppToken     string       `json:"apptoken"`
	LogFlag      int          `json:"logflag"`   //! 是否打印调试信息
	MsgNum       int          `json:"msgnum"`    //! 消息个数
	AppConfig    []AppConfig  `json:"appconfig"` //! 应用配置
	LoggerConfig LoggerConfig `json:"log"`       //! 日志配置
}

type LoggerConfig struct {
	MaxFileSize int64 `json:"maxfilesize"`
	MaxFileNum  int   `json:"maxfilenum"`
}

type AppConfig struct {
	GameId string `json:"gameid"`
	Name   string `json:"name"`
}

/////////////////////////////////////////////////////////////////////////////////////
type CacheServer struct {
	Con            *Config         //! 配置
	Wait           *sync.WaitGroup //! 同步阻塞
	ShutDown       bool            //! 是否正在执行关闭
	MsgGroupCount  int64           //! 最大消息组
	MaxMsgLen      int             //! 消息队列长度
	WaitMsgCount   int64           //! 等待发送数量
	SendMsgCount   int64           //! 发送消息数量
	MinuteMsgCount int64           //! 每分钟新消息数量
	msgCountLock   *sync.RWMutex   //! 消息数读写锁
	MsgGroupMap *Map   			   //! 消息数量: map[string]*GroupInfo
	AppLock *sync.RWMutex 			//! app config lock
}

// 当等待消息数 = 最大消息数 * 70 开始增加group
type GroupInfo struct {
	MsgGroupArr []*MsgGroup   //! 消息组
	MsgIndex int32              //! 当前消息进行到的index
	Lock sync.RWMutex
}

// 增加group
func (self *GroupInfo) addGroup(num int)  {
	self.Lock.Lock()
	for i := 0; i < num; i++ {
		var msgGroup = new(MsgGroup)
		msgGroup.initGroup()
		if self.MsgGroupArr == nil {
			self.MsgGroupArr = make([]*MsgGroup, 0)
		}
		self.MsgGroupArr = append(self.MsgGroupArr, msgGroup)
	}
	self.Lock.Unlock()

	atomic.AddInt64(&GetCacheServer().MsgGroupCount, int64(num))
}

// 设置index
func (self *GroupInfo) addMsgIndex()  {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	self.MsgIndex += 1
	lenArr := len(self.MsgGroupArr)
	if lenArr == 0 {
		return
	}
	self.MsgIndex = self.MsgIndex % int32(lenArr)
}

func (self *GroupInfo) getMsgIndex() int32  {
	self.Lock.RLock()
	defer self.Lock.RUnlock()
	return self.MsgIndex
}

// 检查group
func (self *GroupInfo) checkGroup()  {
	waitMsg := 0
	for index := range self.MsgGroupArr {
		if self.MsgGroupArr[index] == nil {
			continue
		}
		waitMsg += self.MsgGroupArr[index].getLogNum()
	}

	// 可以容纳的最大消息量
	maxMsg := len(self.MsgGroupArr) * GetCacheServer().MaxMsgLen
	if maxMsg == 0 {
		LogDebug("maxMsg == 0")
		return
	}

	f := float32(waitMsg)/float32(maxMsg)
	if f > 0.7 {
		self.addGroup(GetCacheServer().Con.AddGroupNum)
	}
}

var cacheServ *CacheServer = nil

//! public
func GetCacheServer() *CacheServer {
	if cacheServ == nil {
		cacheServ = new(CacheServer)
		cacheServ.Wait = new(sync.WaitGroup)
		cacheServ.Con = new(Config)
		cacheServ.MaxMsgLen = 1000
		cacheServ.msgCountLock = new(sync.RWMutex)
		cacheServ.MsgGroupMap = new(Map)
		cacheServ.AppLock = new(sync.RWMutex)
	}

	return cacheServ
}

//! 载入配置文件
func (self *CacheServer) InitConfig() {
	configFile, err := ioutil.ReadFile("./config_cache.json") ///尝试打开配置文件
	if err != nil {
		log.Fatal("config err 1", err.Error())
	}
	err = json.Unmarshal(configFile, self.Con)
	if err != nil {
		log.Fatal("cache init config, err", err.Error())
	}

	if self.Con.MsgNum == 0 {
		self.Con.MsgNum = 10
	}

	LogDebug("上报AppID:", self.Con.AppID)
	LogDebug("上报AppToken:", self.Con.AppToken)
}

//! 定时器消息打印
func (self *CacheServer) Run() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		self.onTime()
	}

	ticker.Stop()
}

//! 时间响应
func (self *CacheServer) onTime() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println("CacheServer onTime:", x, string(debug.Stack()))
			LogDebug("CacheServer onTime:", x, string(debug.Stack()))
		}
	}()

	tNow := time.Now()
	second := tNow.Second()
	if second == 0 { //! 每分钟统计一次
		LogDebug("普通日志状态:", "线程数=", self.MsgGroupCount, "，等待消息数=", self.WaitMsgCount, ", 发送消息数=", self.SendMsgCount, "，分钟新消息：", self.MinuteMsgCount)
		atomic.StoreInt64(&self.SendMsgCount, 0)
		atomic.StoreInt64(&self.MinuteMsgCount, 0)
		if self.WaitMsgCount <= 0 {
			atomic.StoreInt64(&self.WaitMsgCount, 0)
		}
	}

	if second % 3 == 0 {
		self.checkMsgGroup()

	}
}

//! 日志session处理
func (self *CacheServer) GetConnectHandler() websocket.Handler {
	connectHandler := func(ws *websocket.Conn) {
		if self.ShutDown {
			return
		}
		session := GetSessionMgr().GetNewSession(ws)
		LogDebug("add session:", session.ID)
		session.Run()
	}
	return websocket.Handler(connectHandler)
}

// 一般日志处理
func (self *CacheServer) Log(msg *LogMsg) {
	atomic.AddInt64(&self.MinuteMsgCount, 1)
	atomic.AddInt64(&self.WaitMsgCount, 1)

	gameId := msg.GameId
	groupValueGet, ok := self.MsgGroupMap.Load(gameId)
	var groupInfo *GroupInfo
	if ok {
		//fmt.Printf("%#v\n", groupValueGet)
		groupInfo = groupValueGet.(*GroupInfo)
		msgGroup := groupInfo.MsgGroupArr
		if msgGroup == nil {
			return
		}

		pGroup := msgGroup[groupInfo.getMsgIndex()]
		if pGroup == nil {
			return
		}

		pGroup.Log(msg)
		if len(pGroup.MaxMsg) >= GetCacheServer().Con.MsgNum {
			groupInfo.addMsgIndex()
		}
	} else {
		self.AppLock.Lock()
		config := GetCacheServer().Con.AppConfig
		config = append(config, AppConfig{
			GameId: gameId,
			Name: "unkown",
		})
		self.AppLock.Unlock()

		groupInfo = new(GroupInfo)
		groupInfo.MsgGroupArr = make([]*MsgGroup, 0)
		groupInfo.MsgIndex = 0
		self.MsgGroupMap.Store(gameId, groupInfo)
		groupInfo.addGroup(GetCacheServer().Con.AddGroupNum)
		msgGroup := groupInfo.MsgGroupArr
		if msgGroup == nil {
			return
		}
		pGroup := msgGroup[groupInfo.getMsgIndex()]
		pGroup.Log(msg)
	}
}

// 初始化分配消息组
func (self *CacheServer) startMsgGroup(groupNum int) {
	if atomic.LoadInt64(&self.MsgGroupCount) >= 500 {
		return
	}

	appConfig := self.Con.AppConfig
	for _, conf := range appConfig {
		self.initGameGroup(conf.GameId, groupNum)
	}
}

func (self *CacheServer) initGameGroup(gameId string, groupNum int) {
	gropInfoGet, ok := self.MsgGroupMap.Load(gameId)
	var groupInfo *GroupInfo
	if !ok {
		groupInfo = new(GroupInfo)
		groupInfo.MsgGroupArr = make([]*MsgGroup, 0)
		groupInfo.MsgIndex = 0
		self.MsgGroupMap.Store(gameId, groupInfo)
	} else {
		groupInfo = gropInfoGet.(*GroupInfo)
	}

	if groupInfo != nil {
		groupInfo.addGroup(groupNum)
	}
}

func (self *CacheServer) updateMsgCount(lenCurrent int64) {
	atomic.AddInt64(&self.SendMsgCount, lenCurrent)
	atomic.AddInt64(&self.WaitMsgCount, -lenCurrent)
}

// 检查消息组
func (self *CacheServer) checkMsgGroup() {
	if atomic.LoadInt64(&self.MsgGroupCount) >= 500 {
		return
	}

	appConfig := self.Con.AppConfig
	for _, conf := range appConfig {
		gropInfoGet, ok := self.MsgGroupMap.Load(conf.GameId)
		if !ok  {
			continue
		}

		if groupInfo := gropInfoGet.(*GroupInfo); groupInfo != nil {
			groupInfo.checkGroup()
		}
	}
}


// "appid": "10000024",
//  "apptoken": "7c6bfc266f7a275358b618ecf14cd3ce",
//  "appid": "10000008",
//  "apptoken": "8c56d7c5708858670be880d37b94ebf6",

