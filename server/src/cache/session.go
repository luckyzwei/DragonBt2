package main

import (
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
	"runtime/debug"
	"strings"
	"sync"
)

type Session struct {
	ID       int64           //! 自增长id
	Ws       *websocket.Conn //! websocket
	SendChan chan []byte     //! 发送消息管道
	PID      int64           //! pid
	MsgTime  []int64         //! 时间切片

	lockchan *sync.RWMutex
}

//! 发送消息
func (self *Session) SendMsg(head string, body []byte) {
	if head == "" || string(body) == "" {
		return
	}

	self.lockchan.Lock()
	defer self.lockchan.Unlock()

	if self.SendChan == nil {
		return
	}

	self.SendChan <- HF_DecodeMsg(head, body)
}

func (self *Session) CloseChan() {
	self.lockchan.Lock()
	defer self.lockchan.Unlock()

	if self.SendChan != nil {
		close(self.SendChan)
		self.SendChan = nil
	}
}

//! 消息run
func (self *Session) Run() {
	go self.sendMsgRun()
	self.receiveMsgRun()
}

//! 发送消息循环
func (self *Session) sendMsgRun() {
	for msg := range self.SendChan {
		err := websocket.Message.Send(self.Ws, msg)
		if err != nil {
			//LogError("send err")
			break
		}
	}
	self.Ws.Close()
}

//! 接收消息循环
func (self *Session) receiveMsgRun() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println("receiveMsgRun onTime:", x, string(debug.Stack()))
			LogDebug("receiveMsgRun onTime:", x, string(debug.Stack()))
		}
	}()

	for {
		if GetCacheServer().ShutDown { //! 关服
			break
		}

		var msg []byte
		err := websocket.Message.Receive(self.Ws, &msg)
		if err != nil {
			if err == io.EOF {
				log.Println("receive eof", err.Error())
			} else {
				log.Println("error:", err.Error())
			}
			break
		}
		self.onReceive(msg)
	}
	//LogDebug("client close", self.ID)
	GetSessionMgr().RemoveSession(self)
	self.Ws.Close()
	self.CloseChan()
	self.onClose()
}

func (self *Session) onReceive(msg []byte) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println("onReceive:", x, string(debug.Stack()))
			LogDebug("onReceive:", string(debug.Stack()))
		}
	}()

	if GetCacheServer().ShutDown {
		return
	}

	content := string(msg)

	if GetCacheServer().Con.LogFlag == 1 {
		//log.Println("onReceive content...", content)
		//log.Println("onReceive content...", content)
	}

	msgarr := strings.Split(content, "##")
	if len(msgarr) != 3 {
		log.Println("len(msgarr) := ", len(msgarr), ", content:", content)
		return
	}
	gameId := msgarr[0]
	msgType := HF_Atoi(msgarr[1])
	body := msgarr[2]

	var log LogMsg
	log.MsgType = msgType
	log.GameId = gameId
	log.MsgBuf = []byte(body)
	delimter := []byte("\n")
	log.MsgBuf = append(log.MsgBuf, []byte(delimter)...)
	GetCacheServer().Log(&log)
}

func (self *Session) onClose() {

}

//! 发送消息管道缓冲
const sendChanSize = 1000

type mapSession map[int64]*Session ///定义客户列表类型

var sessionindex int64 = 0

type SessionMgr struct {
	MapSession mapSession
	Lock       *sync.RWMutex
}

var sessionmgrsingleton *SessionMgr = nil

func GetSessionMgr() *SessionMgr {
	if sessionmgrsingleton == nil {
		sessionmgrsingleton = new(SessionMgr)
		sessionmgrsingleton.MapSession = make(mapSession)
		sessionmgrsingleton.Lock = new(sync.RWMutex)
	}

	return sessionmgrsingleton
}

func (self *SessionMgr) GetNewSession(ws *websocket.Conn) *Session {
	sessionindex += 1
	session := new(Session)
	session.ID = sessionindex
	session.Ws = ws
	session.SendChan = make(chan []byte, sendChanSize)
	session.lockchan = new(sync.RWMutex)

	self.Lock.Lock()
	defer self.Lock.Unlock()

	self.MapSession[sessionindex] = session

	return session
}

func (self *SessionMgr) GetSession(id int64) *Session {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	session, ok := self.MapSession[id]

	if ok {
		return session
	}

	return nil
}

func (self *SessionMgr) RemoveSession(session *Session) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	delete(self.MapSession, session.ID)
}

//! 广播消息
func (self *SessionMgr) BroadCastMsg(head string, body []byte) {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	for _, value := range self.MapSession {
		value.SendMsg(head, body)
	}
}

func (self *SessionMgr) GetSessionNum() int {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	return len(self.MapSession)
}
