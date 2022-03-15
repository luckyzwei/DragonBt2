package game

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"time"
)

type CacheClient struct {
	ID       int64
	Ws       *websocket.Conn
	Dead     bool
	OperTime int64 //! 执行下一次操作的时间
}

func (self *CacheClient) InitSocket(url string) bool {
	origin := "http://127.0.0.1/"
	self.Dead = true
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		self.Ws = nil
		log.Println("InitSocket:", err)
		return false
	}
	self.Ws = ws
	self.Dead = false

	self.OperTime = TimeServer().Unix()

	go self.Run()
	go self.Recive()

	return true
}

func (self *CacheClient) Run() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		if self.Dead {
			break
		}

		if GetServer().ShutDown {
			break
		}

		if self.OperTime == 0 {
			continue
		}

		if TimeServer().Unix() >= self.OperTime {
			//! 执行下一次操作

			self.OperTime = 0
		}
	}
}

func (self *CacheClient) Recive() {
	for {
		if self.Dead {
			//! 关服
			break
		}

		var msg []byte
		var tmp []byte

		err := websocket.Message.Receive(self.Ws, &tmp)
		if err != nil {
			if err == io.EOF {
				log.Println("server disconnet")
			} else {
				log.Println("receive err:", err)
			}
			break
		}
		for len(tmp) >= 2048 {
			msg = append(msg, tmp...)
			err = websocket.Message.Receive(self.Ws, &tmp)
			if err != nil {
				if err == io.EOF {
					log.Println("server disconnet")
				} else {
					log.Println("receive err:", err)
				}
				break
			}
		}
		if len(tmp) > 1 || tmp[0] != '\n' {
			msg = append(msg, tmp...)
		}

		self.Receive(msg)
	}
	log.Println("server close")
	self.Ws.Close()
	self.Dead = true
}

func (self *CacheClient) Send(head string, v interface{}) {
	if self.Ws != nil {
		websocket.Message.Send(self.Ws, HF_DecodeMsg(head, HF_JtoB(v)))
	}

}

func (self *CacheClient) Sendstring(gameId string, msgType string, content string) {
	if self.Ws != nil {
		var buf bytes.Buffer
		buf.WriteString(gameId)
		buf.WriteString("##")
		buf.WriteString(msgType)
		buf.WriteString("##")
		buf.WriteString(content)
		websocket.Message.Send(self.Ws, buf.String())
	}

}

func (self *CacheClient) Receive(msg []byte) {
	i := 0
	for i = 0; i < len(msg); i++ {
		if msg[i] == ' ' {
			break
		}
	}

	msg = msg[i+1:]

	b := bytes.NewReader(msg)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)

	_, body, ok := HF_EncodeMsg(out.Bytes())
	if !ok {
		return
	}

	var cid S2C_MsgId
	json.Unmarshal(body, &cid)

	switch cid.Cid {
	default:
		log.Println("收到其他消息:", cid.Cid)
	}
}
