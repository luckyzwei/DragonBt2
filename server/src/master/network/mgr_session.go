package network

import (
	"bytes"
	"golang.org/x/net/websocket"
	"master/utils"
	"sync"
	"time"
)

////////////////////////////////////////////////////////////////////
//! session 管理者

//! 发送消息管道缓冲
const (
	sendChanSize  = 2000 //! 发送管道长度-消息累积长度
	recvChanSize  = 2000 //! 接受管道长度
	socketTimeOut = 5    //! 网络消息超时时间
	socketTryNum  = 10   //! 网络发送重试次数
	broadcast_max = 10
)

type mapSession map[int64]*Session ///定义客户列表类型

var sessionindex int64 = 0

type SessionMgr struct {
	MapSession        mapSession    //! 所有的队列
	AddSession        chan *Session //! 增加的队列
	DelSession        chan *Session //! 断开连接的队列
	Lock              *sync.RWMutex //! 删除锁定
	RemoveLock        *sync.RWMutex //! 删除连接锁
	RemoveSessionList []*Session    //! 保存删除的session做后的删除
	Shutdown          bool          //! 系统是否关闭，关闭后，不再处理消息
	BroadcastMsgArr   chan []byte   //! 广播消息队列
}

var s_sessionmgr *SessionMgr = nil

func GetSessionMgr() *SessionMgr {
	if s_sessionmgr == nil {
		s_sessionmgr = new(SessionMgr)
		s_sessionmgr.MapSession = make(mapSession)
		s_sessionmgr.Lock = new(sync.RWMutex)
		s_sessionmgr.RemoveLock = new(sync.RWMutex)
		s_sessionmgr.RemoveSessionList = make([]*Session, 0)
	}

	return s_sessionmgr
}

//消息广播
func (self *SessionMgr) Run() {
	for msg := range self.BroadcastMsgArr {
		if self.Shutdown { //! 关服
			break
		}

		utils.LogDebug("broadcast message:", string(msg))
		var buffer bytes.Buffer
		buffer.Write(HF_DecodeMsg("1", msg))
		self.Lock.RLock()
		for _, value := range self.MapSession {
			if value.PlayerObj != nil {
				value.SendMsgBatch(buffer.Bytes())
			}
		}
		self.Lock.RUnlock()
	}
	/*NEXT:
	for {
		if GetServer().ShutDown {
			break
		}

		msgCount := 0
		if len(GetServer().BroadCastMsg) > broadcast_max {
			msgCount = broadcast_max
		} else {
			msgCount = len(GetServer().BroadCastMsg)
		}

		if msgCount == 0{
			time.Sleep(time.Millisecond * 10)
			continue NEXT
		}
		var buffer bytes.Buffer
		head := "1"
		msgItr := 0
		for i := 0; i < msgCount; i++ {
			if len(GetServer().BroadCastMsg) > 0 {
				msgItr++
				msg := <-GetServer().BroadCastMsg
				buffer.Write(utils.HF_DecodeMsg(head, msg))
			} else {
				self.Lock.Lock()
				for _, value := range self.MapSession {
					value.SendMsgBatch(buffer.Bytes())
				}
				self.Lock.Unlock()

				continue NEXT
			}
		}

		if msgItr > 0 {
			self.Lock.Lock()
			for _, value := range self.MapSession {
				value.SendMsgBatch(buffer.Bytes())
			}
			self.Lock.Unlock()
		}

	}*/
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
	session.RecvChan = make(chan []byte, recvChanSize)
	session.ShutDown = false
	//session.lockchan = new(sync.RWMutex)
	session.IP = utils.HF_GetHttpIP(session.Ws.Request())
	session.PlayerObj = nil
	session.MsgMaxId = -1

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
			close(self.RemoveSessionList[i].RecvChan)
			self.RemoveSessionList[i].RecvChan = nil
		} else {
			waitCloseSession++
		}
	}
	if waitCloseSession == 0 {
		self.RemoveSessionList = make([]*Session, 0)
	}
}

func (self *SessionMgr) Barrage(size int, text string, red int, green int, blue int, uid int64) {
	var msg S2C_Barrage
	msg.Cid = "barrage"
	msg.Size = size
	msg.Text = text
	msg.Red = red
	msg.Green = green
	msg.Blue = blue
	msg.Uid = uid
	self.BroadCastMsg("1", utils.HF_JtoB(&msg))
}

//! 广播消息
func (self *SessionMgr) BroadCastMsg(head string, body []byte) {
	//return
	if len(self.BroadcastMsgArr) >= 5000 {
		return
	}
	self.BroadcastMsgArr <- body
}

//! 广播消息处理
func (self *SessionMgr) OrderBroadCastMsg(head string, body []byte) {
	//return
	//for _, value := range self.MapSession {
	//	value.SendMsg(head, body)
	//}
}

func (self *SessionMgr) GetSessionNum() int {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	return len(self.MapSession)
}
