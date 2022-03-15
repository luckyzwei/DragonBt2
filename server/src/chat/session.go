package main

import (
	//"crypto/md5"
	//"crypto/rand"
	//"encoding/base64"
	//"encoding/hex"

	"code.google.com/p/go.net/websocket"
	//"encoding/base64"
	"encoding/json"
	"log"

	//"github.com/garyburd/redigo/redis"
	//"net"
	//"time"
	//"github.com/aliyun/aliyun-oss-go-sdk/oss"
	//"github.com/golang/protobuf/proto"
	"io"
	//"os"
	//"bytes"
	//"fmt"
	//"io/ioutil"
	//"net/http"
	//"protobuf"
	//"net/http"
	//"net/url"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

//! client2server
//! 版本验证
type C2S_CtrlHead struct {
	Ctrl string `json:"ctrl"`
	Uid  int64  `json:"uid"`
	Os   string `json:"os"`
	Ver  int    `json:"ver"`
}

//! 发送错误信息
type S2C_ErrInfo struct {
	Cid  string `json:"cid"`
	Info string `json:"info"`
}

type San_Account struct {
	Uid      int64
	Account  string
	Password string
	Creator  string
	Time     int64
}

type Platform_Info struct {
	Platform string //! 平台 android/ios/wp/windows
}

type Session struct {
	ID       int64           //! 自增长id
	Ws       *websocket.Conn //! websocket
	SendChan [][]byte        //! 发送消息管道
	MsgTime  []int64         //! 时间切片
	IP       string          //! IP
	//PlayerObj *Player
	IsClose  bool  //! 是否关闭
	CloseNum int32 //!

	lockchan *sync.RWMutex
}

//! 发送消息
func (self *Session) SendMsg(head string, body []byte) {
	if self.IsClose {
		return
	}

	if head == "" || string(body) == "" {
		return
	}

	self.lockchan.Lock()
	defer self.lockchan.Unlock()

	if self.SendChan == nil {
		return
	}

	LogDebug("s2c:", head, "....", string(body))

	if head == "shutdown" {
		self.IsClose = true
		//self.SendChan = append(self.SendChan, []byte(""))
		return
	}

	self.SendChan = append(self.SendChan, HF_DecodeMsg(head, body))
	if len(self.SendChan) > sendChanSize {
		LogError("Send Chan Overlap...", self.ID)

		self.SendChan = append(self.SendChan, []byte(""))
	}
}

func (self *Session) CloseChan() {
	self.lockchan.Lock()
	defer self.lockchan.Unlock()

	self.SendChan = nil
}

//! 消息run
func (self *Session) Run() {
	go self.sendMsgRun()
	self.receiveMsgRun()
}

func (self *Session) Send() bool {
	self.lockchan.Lock()
	defer self.lockchan.Unlock()

	if self.IsClose {
		return false
	}

	if self.SendChan == nil {
		return false
	}

	if len(self.SendChan) == 0 {
		return true
	}

	msg := self.SendChan[0]

	if string(msg) == "" {
		return false
	}

	self.Ws.SetWriteDeadline(time.Now().Add(sockcetTimeOut * time.Second))
	err := websocket.Message.Send(self.Ws, msg)
	if err != nil {
		neterr, ok := err.(net.Error)
		if ok && neterr.Timeout() {
			LogError("send timeout")
			return true
		}
		LogError("send err", string(msg))
		return false
	}

	self.SendChan = self.SendChan[1:]
	return true
}

//! 发送消息循环
func (self *Session) sendMsgRun() {
	ticker := time.NewTicker(time.Millisecond * 50)
	for {
		<-ticker.C
		if !self.Send() {
			break
		}
	}
	ticker.Stop()
	LogInfo("client close", self.ID)
	self.IsClose = true
	atomic.AddInt32(&self.CloseNum, 1)
	if atomic.LoadInt32(&self.CloseNum) == 2 {
		self.Ws.Close()
		self.CloseChan()
		self.onClose()
		GetSessionMgr().RemoveSession(self)
	}
}

//! 接收消息循环
func (self *Session) receiveMsgRun() {
	for {
		if GetChatServer().ShutDown { //! 关服
			break
		}

		if self.IsClose {
			break
		}

		var msg []byte
		//var tmp []byte

		self.Ws.SetReadDeadline(time.Now().Add(sockcetTimeOut * time.Second))
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
				LogError("receive err:", err)
			}
			break
		}

		if len(self.MsgTime) >= 20 { //! 取20个消息的间隔如果小于1秒
			if self.MsgTime[19]-self.MsgTime[0] <= 1 {
				break
			}
			self.MsgTime = make([]int64, 0)
		} else {
			self.MsgTime = append(self.MsgTime, time.Now().Unix())
		}

		self.onReceive(msg)
	}
	LogInfo("client close", self.ID)
	self.IsClose = true
	atomic.AddInt32(&self.CloseNum, 1)
	if atomic.LoadInt32(&self.CloseNum) == 2 {
		self.Ws.Close()
		self.CloseChan()
		self.onClose()
		GetSessionMgr().RemoveSession(self)
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

	if GetChatServer().ShutDown {
		return
	}

	head, body, _ := HF_EncodeMsg(msg)
	LogDebug("c2s:", head, "....", string(body))

	if head == "chat" { //! 聊天单独处理
		//if self.PlayerObj == nil {
		//	return
		//}
		//self.PlayerObj.SetSessionID(self, false, self.IP, "")
		//self.PlayerObj.onReceive(head, "", body)

		GetSessionMgr().BroadCastMsg("chat", msg)
		return
	}

	var ctrlhead C2S_CtrlHead
	err := json.Unmarshal(body, &ctrlhead)
	if err != nil {
		LogError("CtrlType err:", err)
		return
	}
}

func (self *Session) onClose() {

}

func (self *Session) SendErrInfo(cid string, info string) {
	var msg S2C_ErrInfo
	msg.Cid = cid
	msg.Info = info
	smsg, _ := json.Marshal(&msg)
	self.SendMsg("1", smsg)
}

type JS_SDKData struct {
	Token string `json:"token"`
}

type JS_SDKGame struct {
	Id string `json:"id"`
}

type JS_SDKState struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type JS_SDKHData struct {
	UserId    string `json:"userId"`
	Creator   string `json:"creator"`
	ChannelId string `json:"channelid"`
}

type JS_SDKBody struct {
	Id    string      `json:"id"`
	State JS_SDKState `json:"state"`
	Data  JS_SDKHData `json:"data"`
}

type JS_SDKLogin struct {
	Id   int64      `json:"id"`
	Game JS_SDKGame `json:"game"`
	Data JS_SDKData `json:"data"`
	Sign string     `json:"sign"`
}

////////////////////////////////////////////////////////////////////
//! session 管理者

//! 发送消息管道缓冲
const sendChanSize = 1000
const sockcetTimeOut = 5

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
	session.SendChan = make([][]byte, 0)
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
	//_, ok := self.MapSession[session.ID]
	//if !ok {
	//	return
	//}
	self.Lock.Lock()
	defer self.Lock.Unlock()

	//LogInfo("remove session = ", session.ID)
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
