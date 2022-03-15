// 服务器跨服客户端上传或者获取跨服数据客户端

package game

import (
	"bytes"
	"compress/zlib"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"
)

const (
	serverAgentOrigin = "http://127.0.0.1/"
)

// 服务器上传客户端定义
type ServerAgent struct {
	ID   int64
	Ws   *websocket.Conn
	Dead bool
}

// 初始化ws客户端
func (self *ServerAgent) InitSocket(url string) bool {
	self.Dead = true
	ws, err := websocket.Dial(url, "", serverAgentOrigin)
	if err != nil {
		self.Ws = nil
		log.Println(err, "connect center server failed")
		return false
	}

	self.Ws = ws
	self.Dead = false
	go self.Recive()
	log.Println("连接center服务器成功...")

	return true
}

// 消息接收和解包
func (self *ServerAgent) Recive() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	for {
		if self.Dead {
			break
		}

		var msg []byte
		//var tmp []byte

		self.Ws.SetReadDeadline(TimeServer().Add(socketTimeOut * time.Second))
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

		if len(msg) == 0 {
			continue
		}

		self.Receive(msg)
	}

	log.Println("server agent close")
	self.Ws.Close()
	self.Dead = true
}

// 消息处理
func (self *ServerAgent) Receive(msg []byte) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	if msg == nil || len(msg) <= 0 {
		return
	}

	i := 0
	for i = 0; i < len(msg); i++ {
		if msg[i] == ' ' {
			break
		}
	}

	if len(msg) <= i+1 {
		return
	}

	msg = msg[i+1:]
	b := bytes.NewReader(msg)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)

	head, body, _ := HF_EncodeMsg(out.Bytes())
	LogDebug("receive msg s2center:", head, "....", string(body))
	if head == "mherotopinfo" {
		log.Println("mherotopinfo")
	} else if head == "servgeneralrank" { //! 收到排行榜rank
		log.Println("servgeneralrank")
	} else {
		log.Println("unkown msg :", head)
	}
}

// 向center服务器发送消息
func (self *ServerAgent) SendMsg(msg []byte) {
	if !self.Dead && self.Ws != nil {
		websocket.Message.Send(self.Ws, msg)
	}

}
