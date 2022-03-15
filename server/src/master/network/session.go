package network

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"master/core"
	"master/utils"
	"net"
	"runtime/debug"
	"time"
)

//! 连接类
type Session struct {
	ID        int64           //! 自增长id
	Ws        *websocket.Conn //! websocket
	SendChan  chan []byte     //! 发送消息管道
	RecvChan  chan []byte     //! 发送消息管道
	ShutDown  bool            //! 是否关闭
	ShutTime  int64           //! 关闭时间
	LogicTime int64           //! 逻辑处理时间
	MsgTime   []int64         //! 时间切片
	MsgMaxId  int             //! 当前最大消息Id
	IP        string          //! 公网IP
	PlayerObj core.IPlayer    //! 角色对象
	TryNum    int             //! 重试次数
	Os        string          //! 操作os
	LoginBuf  []byte          //! 缓存登录Buf，用做排队
	//Ctrl     string //! 操作
	//ServerId int    //! 服务器ID
	//Token    string //! ios专有token

	onMessage func(session *Session, body []byte)
	onClose   func(session *Session)

	DevInfo SendRZ_EnvInfo_DevInfo_Ios //! 保存设备信息
	ChInfo  SendRZ_EnvInfo_ChInfo_Ios  //! SDK相关信息
	//lockchan *sync.RWMutex
}

func (self *Session) SetOnMessage(hookFunc func(session *Session, body []byte)) {
	self.onMessage = hookFunc
}

func (self *Session) SetOnClose(hookFunc func(session *Session)) {
	self.onClose = hookFunc
}

func (self *Session) GetId() int64 {
	return self.ID
}

//! 主动关闭
func (self *Session) OnClose() {
	self.onClose(self)
}

//func (self *Session) onClose() {
//	if self.PlayerObj != nil && self.PlayerObj.SessionObj == self {
//		LogDebug("玩家离线，同步离线数据", self.ID, self.PlayerObj.GetUid())
//		GetPlayerMgr().SetPlayerOffline()
//		LogDebug("减少玩家离线数量")
//		self.PlayerObj.onClose()
//		self.PlayerObj = nil
//	}
//	GetLineUpMgr().RemoveClient(self)
//	LogDebug("删除玩家连接信息")
//}

//! 发送消息
func (self *Session) SendMsg(head string, body []byte) {
	if head == "" || string(body) == "" {
		return
	}

	//self.lockchan.Lock()
	//defer self.lockchan.Unlock()

	utils.LogDebug("s2c:", head, "....", string(body))

	if self.ShutDown == true || self.SendChan == nil {
		return
	}

	if head == "shutdown" {
		//self.ShutDown = true
		self.ShutTime = time.Now().Unix()
		self.SendChan <- []byte("")
		return
	}

	if len(self.SendChan) >= sendChanSize-100 {
		utils.LogError("Send Chan Overlap...", self.ID)
		self.SendChan <- []byte("")
		self.ShutDown = true
		self.ShutTime = time.Now().Unix()
		return
	}

	self.SendChan <- HF_DecodeMsg(head, body)
}

func (self *Session) SendMsgBatch(msg []byte) {
	if self.ShutDown == true {
		return
	}

	if self.SendChan == nil {
		return
	}

	self.SendChan <- msg
}

func (self *Session) CloseChan() {
	//self.lockchan.Lock()
	//defer self.lockchan.Unlock()

	if self.SendChan == nil {
		return
	}

	close(self.SendChan)
	self.SendChan = nil

	close(self.RecvChan)
	self.RecvChan = nil
}

//! 消息run
func (self *Session) Run() {
	go self.sendMsgRun()
	go self.logicRun()
	self.receiveMsgRun()
}

func (self *Session) logicRun() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			utils.LogError(x, string(debug.Stack()))
		}
	}()

	//! 0.1s 消息循环
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		if core.GetMasterApp().IsClosed() { //! 关服
			break
		}

		select {
		case packet := <-self.RecvChan:
			self.onMessage(self, packet)
		case <-ticker.C:
			self.onTimer()
		}
	}

	ticker.Stop()

	//! 关闭
	self.onClose(self)
	GetSessionMgr().RemoveSession(self)
}

//! 循环逻辑
func (self *Session) onTimer() {
	//if self.PlayerObj.IsOnline() == false {
	//	return
	//}
	//
	//tNowTime := time.Now().Unix()
	//if self.LogicTime == 0 || tNowTime > self.LogicTime {
	//	self.LogicTime = tNowTime
	//	self.PlayerObj.OnTimer()
	//}
}

//! 发送消息循环
func (self *Session) sendMsgRun() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			utils.LogError(x, string(debug.Stack()))
		}
	}()

	for msg := range self.SendChan {
		if GetSessionMgr().Shutdown { //! 关服
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
						//utils.LogError("send break...", self.ID, len(self.SendChan), len(msg))
						exit = true
						break
					}
					//utils.LogError("send timeout...", self.ID, len(self.SendChan), len(msg))
					continue
				}
				//utils.LogError("send err", string(msg))
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

	utils.LogInfo("client close send", self.ID)
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
			utils.LogError(x, string(debug.Stack()))
		}
	}()

	for {
		if core.GetMasterApp().IsClosed() { //! 关服
			break
		}

		if self.ShutDown == true {
			break
		}

		var msg []byte
		self.Ws.SetReadDeadline(time.Now().Add(socketTimeOut * time.Second))
		err := websocket.Message.Receive(self.Ws, &msg)
		if err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				//utils.LogError("receive timeout")
				continue
			}
			if err == io.EOF {
				utils.LogInfo("client disconnet")
			} else {
				utils.LogInfo("receive err:", err)
			}
			break
		}

		//! 增加筛选没意义，暂时屏蔽
		//if len(self.MsgTime) >= 30 { //! 取20个消息的间隔如果小于1秒
		//	if self.MsgTime[29]-self.MsgTime[0] <= 1 {
		//		break
		//	}
		//	self.MsgTime = make([]int64, 0)
		//} else {
		//	self.MsgTime = append(self.MsgTime, time.Now().Unix())
		//}

		if len(self.RecvChan) < recvChanSize {
			self.RecvChan <- msg
		} else {
			utils.LogError("receive err: overlap ", string(msg))
		}

	}
	utils.LogInfo("client close recv", self.ID)
	self.SendMsg("shutdown", []byte("shutdown"))
	self.ShutDown = true
	self.ShutTime = time.Now().Unix()
	self.Ws.Close()
	utils.LogInfo("关闭玩家websokcet连接")
}

func (self *Session) SendErrInfo(cid string, info string) {
	var msg S2C_ErrInfo
	msg.Cid = cid
	msg.Info = info
	smsg, _ := json.Marshal(&msg)
	self.SendMsg("1", smsg)
}

func (self *Session) SendReturn(cid string) {
	var msg S2C_Result2Msg
	msg.Cid = cid
	smsg, _ := json.Marshal(&msg)
	self.SendMsg("1", smsg)
}
