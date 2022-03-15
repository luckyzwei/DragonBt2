package main

import (
	"code.google.com/p/go.net/websocket"
	"sync"
	"log"
	"runtime/debug"
	"time"
	"net"
	"io"
	"encoding/json"
)

type Session struct {
	ID        int64           //! 自增长id
	Ws        *websocket.Conn //! websocket
	SendChan  chan []byte     //! 发送消息管道
	ShutDown  bool            //! 是否关闭
	ShutTime  int64           //! 关闭时间
	IP        string          //! IP
	PlayerObj *Player         //! 角色对象
	TryNum    int             //! 重试次数
	//lockchan *sync.RWMutex
}

//! 消息run
func (self *Session) Run() {
	go self.sendMsgRun()
	self.receiveMsgRun()
}

//! 发送消息循环
func (self *Session) sendMsgRun() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	for msg := range self.SendChan {
		if GetServer().ShutDown { //! 关服
			break
		}

		if string(msg) == "" || self.ShutDown == true {
			break
		}

		exit := false
		for {
			self.Ws.SetWriteDeadline(time.Now().Add(socketTimeOut * time.Second))
			err := websocket.Message.Send(self.Ws, msg)
			if err != nil {
				neterr, ok := err.(net.Error)
				if ok && neterr.Timeout() {
					self.TryNum++
					if self.TryNum >= socketTryNum {
						//LogError("send break...", self.ID, len(self.SendChan), len(msg))
						exit = true
						break
					}
					//LogError("send timeout...", self.ID, len(self.SendChan), len(msg))
					continue
				}
				//LogError("send err", string(msg))
				exit = true
				break
			} else {
				break
			}
		}
		if exit {
			break
		}
	}

	LogInfo("client close send", self.ID)
	self.ShutDown = true
	self.ShutTime = time.Now().Unix()
	self.Ws.Close()
}

//! 接收消息循环
func (self *Session) receiveMsgRun() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	for {
		if GetServer().ShutDown { //! 关服
			break
		}

		if self.ShutDown == true {
			break
		}

		var msg []byte
		//var tmp []byte

		self.Ws.SetReadDeadline(time.Now().Add(socketTimeOut * time.Second))
		err := websocket.Message.Receive(self.Ws, &msg)
		if err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				//LogError("receive timeout")
				continue
			}
			if err == io.EOF {
				LogInfo("client disconnet")
			} else {
				LogInfo("receive err:", err)
			}
			break
		}

		self.onReceive(msg)
	}
	LogInfo("client close recv", self.ID)

	self.ShutDown = true
	self.ShutTime = time.Now().Unix()
	self.Ws.Close()
	self.onClose()
	GetSessionMgr().RemoveSession(self)
}

func (self *Session) onClose() {
	if self.PlayerObj != nil && self.PlayerObj.SessionObj == self {
		//LogInfo("玩家离线，同步离线数据", self.ID, self.PlayerObj.Sql_UserBase.Uid)
		self.PlayerObj = nil
	}
}


func (self *Session) onReceive(msg []byte) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	if GetServer().ShutDown {
		return
	}
}

func (self *Session) sendMsg(v interface{}) {

	if self.ShutDown == true || self.SendChan == nil {
		return
	}

	if data, err := json.Marshal(v); err == nil {
		self.SendChan <- data;
	}
}

const sendChanSize = 2000
const socketTimeOut = 5
const socketTryNum = 10
const broadcast_max = 10

type mapSession map[int64]*Session ///定义客户列表类型

var sessionindex int64 = 0

type SessionMgr struct {
	MapSession        mapSession    //
	Lock              *sync.RWMutex // 锁定
	RemoveLock        *sync.RWMutex
	RemoveSessionList []*Session // 保存删除的session做后的删除
}

var sessionmgrsingleton *SessionMgr = nil

func GetSessionMgr() *SessionMgr {
	if sessionmgrsingleton == nil {
		sessionmgrsingleton = new(SessionMgr)
		sessionmgrsingleton.MapSession = make(mapSession)
		sessionmgrsingleton.Lock = new(sync.RWMutex)
		sessionmgrsingleton.RemoveLock = new(sync.RWMutex)
		sessionmgrsingleton.RemoveSessionList = make([]*Session, 0)
	}

	return sessionmgrsingleton
}

func (self *SessionMgr) GetNewSession(ws *websocket.Conn) *Session {
	//if len(GetServer().AddSession) >= 10000 {
	//	return nil
	//}
	self.Lock.Lock()
	sessionindex += 1

	session := new(Session)
	session.ID = sessionindex
	session.Ws = ws
	session.SendChan = make(chan []byte, sendChanSize)
	session.ShutDown = false
	//session.lockchan = new(sync.RWMutex)
	session.IP = HF_GetHttpIP(session.Ws.Request())
	session.PlayerObj = nil


	self.MapSession[sessionindex] = session
	self.Lock.Unlock()
	//GetServer().AddSession <- session

	return session
}

//! 删除session
func (self *SessionMgr) RemoveSession(session *Session) {
	self.Lock.Lock()
	delete(self.MapSession, session.ID)
	self.Lock.Unlock()

	self.RemoveLock.Lock()
	self.RemoveSessionList = append(self.RemoveSessionList, session)
	self.RemoveLock.Unlock()
	//GetServer().DelSession <- session
}

//集中清理session
func (self *SessionMgr) ClearRemoveSession() {
	self.RemoveLock.Lock()
	defer self.RemoveLock.Unlock()
	waitCloseSession := 0
	tNow := time.Now().Unix()
	for i := 0; i < len(self.RemoveSessionList); i++ {
		if self.RemoveSessionList[i].ShutDown &&
			self.RemoveSessionList[i].SendChan != nil &&
			tNow > self.RemoveSessionList[i].ShutTime {
			close(self.RemoveSessionList[i].SendChan)
			self.RemoveSessionList[i].SendChan = nil
		} else {
			waitCloseSession++
		}
	}
	if waitCloseSession == 0 {
		self.RemoveSessionList = make([]*Session, 0)
	}
}