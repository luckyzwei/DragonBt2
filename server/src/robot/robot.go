package main

import (
	"log"
	"net"

	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"sync/atomic"
	"time"
)

const (
	dis_low  =  2 //短间隔
	dis_mid  =  3 //中间隔
	dis_high  =  5 //长间隔
	dis_highest  =  10 //超长间隔
)




type IRobotLogic interface {
	SetRobot(*Robot)
	Receive([]byte)
	OnTimer()
	OnLoginOK()
}

func (self *Robot) InitSocket() bool {
	origin := "http://127.0.0.1/"
	url := GetRobotMgr().RobotCon.Serurl
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println(err)
		return false
	}
	self.Ws = ws
	self.Dead = false
	self.IsLoginOut = false
	self.csvId = 1
	self.OperTime = time.Now().Unix()
	self.is_fightinfo = false
	//GetRobotMgr().AddRobot(self)

	if GetRobotMgr().RobotCon.RoboTtype == robot_type_1 {
		go self.Run()
	} else if GetRobotMgr().RobotCon.RoboTtype == robot_type_2 {
		go self.csvRun()
	}

	go self.Recive()

	return true
}


func (self *Robot) fightRun() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		if self.Dead {
			break
		}

		if self.OperTime == 0 {
			continue
		}

		if time.Now().Unix() >= self.OperTime {

			switch self.msgId {
			case 0:
				self.SendLogin()
				self.OperTime = time.Now().Unix() + 2

			case 4:
				log.Println("机器人:", self.Uname, "uid:", self.Uid, "关闭")
				self.OperTime = time.Now().Unix() + 1
				self.Dead = true
			}

			if self.msgId <= 2 {
				fmt.Println("fightRun.msgId =",self.msgId)
				self.msgId++
			}
		}

		if self.logic != nil {
			self.logic.OnTimer()
		}
	}
}



func (self *Robot) Recive() {
	for {

		if self.Dead { //! 关服
			break
		}

		if self.IsLoginOut == false {
			var msg []byte

			err := websocket.Message.Receive(self.Ws, &msg)
			if err != nil {
				neterr, ok := err.(net.Error)
				if ok && neterr.Timeout() {
					//LogError("receive timeout")
					continue
				}
				if err == io.EOF {
					log.Println("client disconnet")
				} else {
					log.Println("receive err:", err)
				}
				break
			}
			//self.RecvChan<-msg
			self.Receive(msg)
		}
	}

	if self.IsLoginOut == false {
		log.Println("server close")
		self.Ws.Close()
		self.Dead = true
	}

}

func (self *Robot) Send(head string, v interface{}) {
	self.msgId += 1
	str := HF_JtoA(v)
	if len(str) > 1 {
		str = fmt.Sprintf("{\"msgid\":%v,%s", self.MsgSN, str[1:])
	}

	atomic.AddInt64(&self.MsgSN, 1)
	websocket.Message.Send(self.Ws, HF_DecodeMsg(head, []byte(str)))
}

func (self *Robot) Start() bool {
	self.origin = GetRobotMgr().RobotCon.Origin
	self.ServerUrl = GetRobotMgr().RobotCon.Serurl
	ws, err := websocket.Dial(self.ServerUrl, "", self.origin)
	if err != nil {
		log.Println(fmt.Sprintf("机器人连接服务器：%s失败", self.ServerUrl))
		return false
	}
	self.Ws = ws
	self.Dead = false
	self.IsLoginOut = false
	self.csvId = 1
	self.OperTime = time.Now().Unix()
	self.is_fightinfo = false
	self.RobotType = GetRobotMgr().RobotCon.RoboTtype
	//fmt.Println("机器人Type:",self.RobotType)
	if self.RobotType == robot_type_1 {
		go self.Run()
	} else if self.RobotType == robot_type_2 {
		go self.csvRun()
	} else if self.RobotType == robot_type_3 {
		self.logic = new(RobotFight)
		self.logic.SetRobot(self)
		go self.fightRun()
	} else if self.RobotType == 4 {
		go self.Logic()
	} else if self.RobotType == robot_type_5 {
		go self.BrainRun()
	}
	go self.Recive()
	return true
}

func (self *Robot) Close()  {
	self.Ws.Close()
	fmt.Println("断开连接")
	time.Sleep(time.Second)
	fmt.Println("重新连接")
	self.Start()
}
