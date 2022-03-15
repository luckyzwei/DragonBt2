package main

import (
	"encoding/json"

	//"fmt"
	"io/ioutil"
	"log"

	"code.google.com/p/go.net/websocket"
	//"math/rand"
	//"os"
	//"runtime"
	//"runtime/debug"
	//"reflect"
	"bytes"
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"
)

//!消息组
type MsgGroup struct {
	groupId int64
	Wait    *sync.WaitGroup
	logChan chan []byte
}

//! 初始化消息组
func (self *MsgGroup) initGroup() {
	self.Wait = new(sync.WaitGroup)
	self.logChan = make(chan []byte, GetChatServer().Con.MaxMsgChan)

	go self.RunLog()
}

//! 获取当前队列消息数量
func (self *MsgGroup) getLogNum() int {
	return len(self.logChan)
}

//! 加入新的日志
func (self *MsgGroup) Log(logs []byte) {
	//log.Println("加入日志:", self.groupId, string(logs), len(self.logChan))
	self.logChan <- logs
}

//! 数据上报
func (self *MsgGroup) RunLog() {
	self.Wait.Add(1)
	//log.Println("send log:", self.groupId, "step start...")
	for msg := range self.logChan {
		if string(msg) == "close" {
			break
		}
		for i := 0; i < 10; i++ {
			if self.SendLog(msg) {
				break
			}
		}

	}
	self.Wait.Done()
	//log.Println("send log:", self.groupId, "step ok...")
	close(self.logChan)
}

//! 发送日志
func (self *MsgGroup) SendLog(logs []byte) bool {
	//状态更新
	GetChatServer().lockCount.Lock()
	GetChatServer().SendMsgCount += 1
	GetChatServer().WaitMsgCount -= 1
	GetChatServer().lockCount.Unlock()

	//time.Sleep(time.Millisecond * 100)
	//log.Println("上报成功:", string(logs))
	//return true

	url := "https://log.rts.dp.uc.cn:8083/api/v2_1"
	body := bytes.NewBuffer(logs)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
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

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("X-dp-app-id", "496909086288")
	req.Header.Set("X-dp-api-token", "8c56d7c5708858670be880d37b94ebf6")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	if resp.StatusCode != 200 {
		log.Println("上报错误:", string(logs))
		//GetLogMgr().Output("上报错误:", string(logs))
		return false
	} else {
		//log.Println("上报成功:", string(logs))
	}

	return true
}

//! 配置
type Config struct {
	Host       string `json:"wshost"`     //! 服务器
	ServerId   int    `json:"serverid"`   //! 服
	ServerName string `json:"servername"` //! 服务器名称
	//OpenTime     string `json:"opentime"`     //! 开服时间
	InitGroupNum int `json:"initgroupnum"` //! 初始消息组
	AddGroupNum  int `json:"addgroupnum"`  //! 每次增加消息组
	MaxMsgChan   int `json:"maxmsgchan"`   //! 最大消息队列
}

type msgGroups map[int64]*MsgGroup ///定义客户列表类型
type userMsgMap map[int64]int64    ///角色绑定管道

/////////////////////////////////////////////////////////////////////////////////////
type ChatServer struct {
	Con            *Config         //! 配置
	Wait           *sync.WaitGroup //! 同步阻塞
	ShutDown       bool            //! 是否正在执行关闭
	MsgGroups      msgGroups       //! 消息组
	UserMsgMap     userMsgMap      //! 用户消息管道映射
	MsgGroupCount  int64           //! 最大消息组
	MsgGroupIndex  int64           //! 当前插入消息索引
	MaxMsgLen      int64           //! 消息队列长度
	WaitMsgCount   int64           //! 等待发送数量
	SendMsgCount   int64           //! 发送消息数量
	MinuteMsgCount int64           //! 每分钟新消息数量
	lockCount      *sync.RWMutex   //!统计消息数
}

var chatServerSingleton *ChatServer = nil

//! public
func GetChatServer() *ChatServer {
	if chatServerSingleton == nil {
		chatServerSingleton = new(ChatServer)
		chatServerSingleton.Wait = new(sync.WaitGroup)
		chatServerSingleton.Con = new(Config)

		chatServerSingleton.MsgGroups = make(msgGroups)
		chatServerSingleton.UserMsgMap = make(userMsgMap)
		chatServerSingleton.MsgGroupIndex = 0
		chatServerSingleton.MaxMsgLen = 1000
		chatServerSingleton.lockCount = new(sync.RWMutex)
		//cacheServerSingleton.msgGroupCount = 10
		//self.startMsgGroup(10)
	}

	return chatServerSingleton
}

//! 载入配置文件
func (self *ChatServer) InitConfig() {
	configFile, err := ioutil.ReadFile("./config_chat.json") ///尝试打开配置文件
	if err != nil {
		log.Fatal("config err 1")
	}
	err = json.Unmarshal(configFile, self.Con)
	if err != nil {
		log.Fatal("chat InitConfig err", err.Error())
	}

	//!
	//	if self.Con.UpRecord == 1 { //! 需要数据上报
	//		go self.RunLog()
	//	}

	//	if self.Con.NumRecord != "" { //! 上报人数
	//		go self.NumRecord()
	//	}
}

//! 消息循环
func (self *ChatServer) Run() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		self.onTime()
	}

	ticker.Stop()
}

//! 时间响应
func (self *ChatServer) onTime() {
	tNow := time.Now()

	if tNow.Second() == 0 { //! 每分钟统计一次
		log.Println("状态:", "线程数=", self.MsgGroupCount, "，等待消息数=", self.WaitMsgCount, ", 发送消息数=", self.SendMsgCount, "，分钟新消息：", self.MinuteMsgCount)
		self.lockCount.Lock()
		self.MinuteMsgCount = 0
		self.SendMsgCount = 0
		self.lockCount.Unlock()
	}

	if tNow.Second()%3 == 0 {
		self.lockCount.Lock()
		if self.WaitMsgCount/self.MsgGroupCount > self.MaxMsgLen/3 {
			self.startMsgGroup(self.Con.AddGroupNum)
		}
		self.lockCount.Unlock()
	}
}

//! 得到一个websocket处理句柄
func (self *ChatServer) GetConnectHandler() websocket.Handler {
	connectHandler := func(ws *websocket.Conn) {
		if self.ShutDown {
			return
		}
		session := GetSessionMgr().GetNewSession(ws)
		log.Println("add session:", session.ID)
		session.Run()
	}
	return websocket.Handler(connectHandler)
}

func (self *ChatServer) Log(uid int64, logs []byte) {
	self.lockCount.Lock()
	self.MinuteMsgCount += 1
	self.WaitMsgCount += 1
	self.lockCount.Unlock()

	if v, ok := self.UserMsgMap[uid]; ok {
		var msgGroup = self.MsgGroups[v]
		if msgGroup.getLogNum() < int(self.MaxMsgLen) {
			msgGroup.Log(logs)
		} else {
			//! 队列已满，本频道关联客户端过多
			var groupIndex int64
			for groupIndex = 0; groupIndex < self.MsgGroupCount; groupIndex++ {
				if self.MsgGroups[groupIndex].getLogNum() < int(self.MaxMsgLen)/3 {
					self.MsgGroups[groupIndex].Log(logs)
					self.MsgGroupIndex = groupIndex
					self.UserMsgMap[uid] = groupIndex
					break
				}
			}
		}
	} else {
		//! 新用户，按顺序分配一个序号
		if self.MsgGroupIndex < self.MsgGroupCount && self.MsgGroupIndex >= 0 {
			var msgGroup = self.MsgGroups[self.MsgGroupIndex]
			if msgGroup.getLogNum() < int(self.MaxMsgLen)/3 {
				msgGroup.Log(logs)
				self.UserMsgMap[uid] = self.MsgGroupIndex
			} else {
				var groupIndex int64
				for groupIndex = 0; groupIndex < self.MsgGroupCount; groupIndex++ {
					if self.MsgGroups[groupIndex].getLogNum() < int(self.MaxMsgLen)/3 {
						self.MsgGroups[groupIndex].Log(logs)
						self.MsgGroupIndex = groupIndex
						self.UserMsgMap[uid] = groupIndex
						break
					}
				}
			}

			//! 循环加入
			self.MsgGroupIndex += 1
			if self.MsgGroupIndex == self.MsgGroupCount {
				self.MsgGroupIndex = 0
			}
		}
	}

}

func (self *ChatServer) startMsgGroup(groupNum int) {
	for i := 0; i < groupNum; i++ {
		var msgGroup = new(MsgGroup)
		msgGroup.groupId = self.MsgGroupCount
		msgGroup.initGroup()
		self.MsgGroups[self.MsgGroupCount] = msgGroup
		self.MsgGroupCount += 1
	}
}
