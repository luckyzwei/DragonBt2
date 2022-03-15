package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
)

type RobotMgr struct {
	MapRobot  map[int64]*Robot
	Lock      *sync.RWMutex
	ID        int64
	RobotCon  *RobotConfig //! 配置
	PlayerIdx int
	RobotList []*Robot
}

var robotmgrsingleton *RobotMgr = nil

//! public
func GetRobotMgr() *RobotMgr {
	if robotmgrsingleton == nil {
		robotmgrsingleton = new(RobotMgr)
		robotmgrsingleton.MapRobot = make(map[int64]*Robot)
		robotmgrsingleton.Lock = new(sync.RWMutex)
		robotmgrsingleton.ID = 0
		robotmgrsingleton.RobotCon = new(RobotConfig)
		robotmgrsingleton.PlayerIdx = 0
	}

	return robotmgrsingleton
}

func (self *RobotMgr) Init() {
	self.InitType1()
	//self.InitType2()
	//self.InitType3()
	go self.Run()
}

func (self *RobotMgr) InitType1() {
	if self.RobotCon.RoboTtype != robot_type_1 {
		return
	}

	for j := 0; j < self.RobotCon.MaxNum; j++ {
		robot := new(Robot)
		robot.RobotType = self.RobotCon.RoboTtype
		if self.RobotCon.RoboTtype == robot_type_1 {
			robot.ServerUrl = self.RobotCon.Serurl
			self.RobotList = append(self.RobotList, robot)
		}
	}
}

func (self *RobotMgr) InitType2() {
	if self.RobotCon.RoboTtype != robot_type_2 {
		return
	}

	for j := 0; j < self.RobotCon.MaxNum; j++ {
		robot := new(Robot)
		robot.RobotType = self.RobotCon.RoboTtype
		if robot.RobotType == robot_type_2 {
			robot.ServerUrl = self.RobotCon.Serurl
			self.RobotList = append(self.RobotList, robot)
		}
	}
}

func (self *RobotMgr) InitType3() {
	if self.RobotCon.RoboTtype != robot_type_3 {
		return
	}

	for j := 0; j < self.RobotCon.MaxNum; j++ {
		robot := new(Robot)
		robot.RobotType = self.RobotCon.RoboTtype
		if robot.RobotType == robot_type_3 {
			robotfight := GetRobotCsvMgr().Data2["robotfight"]
			robot.ServerUrl = robotfight[j]["server"]
			robot.Account = robotfight[j]["account"]
			robot.Password = robotfight[j]["password"]
			robot.ServerId, _ = strconv.Atoi(robotfight[j]["serverid"])
		}
		self.RobotList = append(self.RobotList, robot)
	}
}

// 执行机器人逻辑
func (self *RobotMgr) Run() {
	fmt.Println("机器人个数：", GetRobotMgr().RobotCon.MaxNum)
	fmt.Println("在线机器人平均操作间隔：", GetRobotMgr().RobotCon.ActionDis)

	for i := 0; i < GetRobotMgr().RobotCon.MaxNum; i++ {
		// 先启动一个机器人
		robot := new(Robot)
		robot.ID = int64(i + 1)
		robot.Account = fmt.Sprintf("robot-%d", HF_GetRandom(100000000))
		robot.Password = "1"
		robot.Start()
	}

}

//! 载入配置文件
func (self *RobotMgr) InitConfig() {
	configFile, err := ioutil.ReadFile("robotconfig.json") ///尝试打开配置文件
	if err != nil {
		log.Fatal("robotconfig err 1")
	}
	err = json.Unmarshal(configFile, self.RobotCon)
	if err != nil {
		log.Fatal("robotconfig err 2", err.Error())
	}
	fmt.Println("初始化完表配置:", self.RobotCon)

	if self.RobotCon.ActionDis < 10 {
		self.RobotCon.ActionDis = 10
	}
}
