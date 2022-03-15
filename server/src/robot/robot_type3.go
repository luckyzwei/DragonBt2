package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"time"
)

func (self *RobotFight) SetRobot(robot *Robot) {
	self.Owner = robot
}

func (self *RobotFight) OnLoginOK() {
	self.Owner.SendTeamPowerReq()
}

func (self *RobotFight) IsCityOwner(cityid int) bool {
	for i := 0; i < len(self.CityInfo); i++ {
		if self.CityInfo[i].Id == cityid {
			return self.CityInfo[i].Camp == self.Owner.userbase.Camp
		}
	}

	return false
}

func (self *RobotFight) Receive(body []byte) {

	var cid S2C_MsgId
	json.Unmarshal(body, &cid)

	switch cid.Cid {
	case "cityinfo":
		var msg S2C_CityInfo
		json.Unmarshal(body, &msg)
		self.CityInfo = msg.Info
		return

	case "diplomacyinfo":

		log.Println("diplomacyinfo")

		//var msg S2C_DiplomacyInfo
		//json.Unmarshal(body, &msg)
		//findCity := false
		//for i := 0; i < 3; i++ {
		//	if msg.Status[i] == 1 && findCity == false {
		//		//尝试进入开战国战，只进入一次
		//		if msg.Align[0] > 0 { //! 结盟了任何国战都进入
		//			self.CampFightCity = msg.AttackCity[i][0]
		//			findCity = true
		//		} else {
		//			if self.Owner.userbase.Camp == i+1 {
		//				self.CampFightCity = msg.AttackCity[i][0]
		//				findCity = true
		//			}
		//		}
		//	}
		//}
		//
		//if findCity == false {
		//
		//	for i := 0; i < 3; i++ {
		//		if msg.Status[i] == 1 && findCity == false {
		//			//如果被打的城市是自己的势力
		//			if self.IsCityOwner(msg.AttackCity[i][0]) {
		//				self.CampFightCity = msg.AttackCity[i][0]
		//				findCity = true
		//			}
		//		}
		//	}
		//}
		//
		//if findCity == true {
		//	self.Owner.SendCampFightMoveEx(self.CampFightCity)
		//}
		//
		//self.State = 1

		return
	case "campfightsololist":

		log.Println("campfightsololist")

		var msg S2C_CampFightSoloList
		json.Unmarshal(body, &msg)

		self.fightinfo = msg.FightInfo

		self.Timer = 0
		self.Stage = 2
		return
	case "campfightresult":
		log.Println("campfightresult")
		self.Reset()
		return
	case "campfightover":
		log.Println("campfightover")
		self.Timer = 0
		self.Stage = 0
		return
	case "campfightsolo2begin":

		log.Println("campfightsolo2begin")
		self.Timer = 0
		self.Stage = 4
		return
	case "baglst":
		var msg S2C_BagInfo
		json.Unmarshal(body, &msg)
		self.Baglst = msg.Baglst
		return
	case "campfightplayinfo":
		log.Println("campfightplayinfo")
		var msg S2C_CampFightPlayInfo
		json.Unmarshal(body, &msg)

		for i := 0; i < len(msg.Playlist); i++ {
			if msg.Playlist[i].PlayId == 12001 || msg.Playlist[i].PlayId == 12002 { //! 开启国战
				self.Owner.SendCampFight56Req(msg.Playlist[i].PlayId, self.CampFightCity)
			}
		}
		return
	case "teampower":
		var msg S2C_TeamPower
		json.Unmarshal(body, &msg)
		self.Power = msg.Power
		return
	}
}

func (self *RobotFight) OnTimer() {

	if self.State == 0 {
		return
	}

	self.TotalTimer += 1
	if self.TotalTimer >= 20 {

		self.Owner.GetCampfightPlayList(self.CampFightCity, 0)

		if self.Power < GetRobotMgr().RobotCon.PowerCheckNum {
			self.Owner.SendBuyTeamPowerReq()
		}

		self.TotalTimer = 0
	}

	self.Timer += 1

	if self.Stage == 0 {
		//
		self.Owner.SendCampfightSoloReq(self.CampFightCity, 0)

		self.Stage = 1
	} else if self.Stage == 1 {

	} else if self.Stage == 2 {

		if len(self.fightinfo) > 0 {
			target := self.fightinfo[0]
			self.Owner.SendCampFightSolo2(target.Index, target.Uid, self.CampFightCity)
		}

		self.Stage = 3
	} else if self.Stage == 3 {

	} else if self.Stage == 4 {

		if self.Timer >= 30 {

			result := 1
			if self.Wins >= 9 {
				result = 0
			}
			self.Owner.SendCampFightSoloEnd(result, self.CampFightCity)
			self.Wins += 1
			self.Stage = 5
			self.Timer = 0
		}

	} else if self.Stage == 5 {
		if self.Timer >= 30 {
			self.Stage = 0
			self.Timer = 0
		}
	}

}

func (self *RobotFight) Reset() {

	self.Stage = 0
	self.Timer = 0
	self.State = 0
	self.TotalTimer = 0
	self.Wins = 0
}

func (self *RobotFight) Send(head string, v interface{}) {

	self.Owner.Send(head, v)
}

func CheckErr(err error)  {
	if err != nil {
		log.Println(err.Error())
		return
	}
}


func (self *Robot) Receive(msg []byte) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			//LogError(x, string(debug.Stack()))
		}
	}()
	i := 0

	for i = 0; i < len(msg); i++ {
		if msg[i] == ' ' {
			break
		}
	}
	if i < len(msg)+3 {
		msg = msg[i+1:]
	} else {
		return
	}

	b := bytes.NewReader(msg)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)

	_, body, ok := HF_EncodeMsg(out.Bytes())
	if !ok {
		return
	}

	if self.logic != nil {
		self.logic.Receive(body)
	}

	//fmt.Println("receive body:", string(body))

	var cid S2C_MsgId
	json.Unmarshal(body, &cid)

	switch cid.Cid {
	case "reg": //! 注册
		var msg S2C_Reg
		err := json.Unmarshal(body, &msg)
		CheckErr(err)
		self.Account = msg.Account
		self.Password = msg.Password
		return
	case "userbaseinfo":
		//log.Println("收到基本信息")
		var msg S2C_UserBaseInfo
		err := json.Unmarshal(body, &msg)
		CheckErr(err)
		self.Uid = msg.Baseinfo.Uid
		self.Uname = msg.Baseinfo.Uname
		self.userbase.Gold = msg.Baseinfo.Gold
		self.userbase.Gem = msg.Baseinfo.Gem
		self.userbase.Camp = msg.Baseinfo.Camp
		self.userbase.Power = msg.Baseinfo.Tili
		self.userbase.Level = msg.Baseinfo.Level
		self.SendLoginOK()
		if msg.Baseinfo.CampOk == 0 {
			self.SetCamp()
		}

		if msg.Baseinfo.NameOk == 0 {
			self.SendCreateRole()
		}

		if self.logic != nil {
			self.logic.OnLoginOK()
		}
		return
	case "kickout":
		self.Ws.Close()
		return
	case "shutdown":
		self.Ws.Close()
		return

	case "loginokret":
		var msg S2C_LoginRet
		err := json.Unmarshal(body, &msg)
		CheckErr(err)
		if msg.Ret == 0 {
			fmt.Println("登录成功")
			//self.Close()
			self.SendMsg()
		}

		return
	case "createrole":
		fmt.Println("创建角色成功:", time.Now().Format(TIMEFORMAT))
		return
	case "onitem":
		var msg S2C_OnItem
		err := json.Unmarshal(body, &msg)
		CheckErr(err)

		for index:=0;index<len(msg.Itemlst);index++{
			self.AddItem(msg.Itemlst[index].ItemID,msg.Itemlst[index].Num)
		}
		return
	case "zdsj":
		var msg S2C_UpdateExp
		err := json.Unmarshal(body, &msg)
		CheckErr(err)
		self.userbase.Level=msg.New
		return
	case "updataenergy":
		var msg S2C_UpdateTiLi
		err := json.Unmarshal(body, &msg)
		CheckErr(err)
		self.userbase.Power=msg.Tili
		return
	case "drawok":
		var msg S2C_FindOK
		err := json.Unmarshal(body, &msg)
		CheckErr(err)
		for index:=0;index<len(msg.Cost);index++{
			self.AddItem(msg.Cost[index].ItemID,msg.Cost[index].Num)
		}
		return
	default:

	}
}

func (self *Robot)AddItem(itemId int,itemNum int)  {
	switch itemId{
	case 91000001:
		self.userbase.Gold+=itemNum
	case 91000002:
		self.userbase.Gem+=itemNum
	}
}
