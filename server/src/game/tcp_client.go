package game

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

//! TCP 连接
type TCPClient struct {
	ID       int64        //! 连接ID
	Conn     *net.TCPConn //! 连接对象
	Dead     bool         //! 是否关闭
	OperTime int64        //! 执行下一次操作的时间
}

func (self *TCPClient) InitSocket(url string) bool {
	address := "api.chat.q-dazzle.com:5678"
	tcpaddr, err := net.ResolveTCPAddr("", address)
	if err != nil {
		fmt.Println("net ResolveTCPAddr error! ", err.Error())
		//os.Exit(-1)
		return false
	}

	fmt.Println("str tcpaddr = ", tcpaddr.String())
	fmt.Println("str Network = ", tcpaddr.Network())

	conn, err := net.DialTCP("tcp4", nil, tcpaddr)
	if err != nil {
		fmt.Println("net DialTcp Error!", err.Error())
		return false
	}
	//defer conn.Close()

	self.Conn = conn
	self.Dead = false

	self.OperTime = TimeServer().Unix()

	//go self.Run()
	//go self.Recive()

	return true
}

func (self *TCPClient) Run() {
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

		//if TimeServer().Unix() >= self.OperTime {
		//	//! 执行下一次操作
		//
		//	self.OperTime = 0
		//}
	}
}

func (self *TCPClient) Recive() {
	for {
		if self.Dead {
			//! 关服
			break
		}
		//var msg []byte
		//var tmp []byte

		result, err := ioutil.ReadAll(self.Conn)
		if err != nil {
			fmt.Println("ReadAll error: ", err.Error())
			self.Dead = true
		}

		if len(result) == 0 {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		LogDebug("result = %s", string(result))
		//err := websocket.Message.Receive(self.Conn, &tmp)
		if err != nil {
			if err == io.EOF {
				log.Println("server disconnet")
			} else {
				log.Println("receive err:", err)
			}
			break
		}

		//! 收到发送成功数据，则连接断开
		break

		//if len(tmp) > 1 || tmp[0] != '\n' {
		//	msg = append(msg, tmp...)
		//}
		//
		//self.Receive(result)
	}
	log.Println("server close", HF_GetGoroutineID())
	self.Conn.Close()

	//self.Dead = true
	self.InitSocket("")
}

func (self *TCPClient) Send(str string) {
	if self.Conn != nil {
		blen, err := self.Conn.Write([]byte(str))
		if err != nil {
			fmt.Println("err = ", err.Error())
			self.Dead = true
		}
		LogDebug("blen = %d", blen)
	}
}

func (self *TCPClient) Receive(msg []byte) {
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
