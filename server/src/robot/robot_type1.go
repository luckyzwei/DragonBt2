package main

import (
	"fmt"
	"time"
)

// 最老版的机器人 用于密集登录压测
func (self *Robot) Run() {
	ticker := time.NewTicker(time.Second)
	self.csvId = 0
	self.guildId = 0
	for {
		<-ticker.C
		if self.Dead {
			break
		}
		if self.OperTime == 0 {
			continue
		}
		//fmt.Println("csvId =",self.csvId)
		if time.Now().Unix() >= self.OperTime {
			//! 执行下一次操作
			//fmt.Println("csvId =",self.csvId)
			switch self.csvId {
			case 0:
				self.SendLogin()
				self.OperTime = time.Now().Unix() + dis_highest
			case 1:
				self.SendCampTeam([3]JS_CampTeam{{[5]int{0, 0, 0, 0, 0}, 0, 0, 1104, "", 0, 0},
					{[5]int{0, 0, 0, 0, 0}, 0, 0, 1104, "", 0, 0},
					{[5]int{0, 0, 0, 0, 0}, 0, 0, 1104, "", 0, 0}})
				self.OperTime = time.Now().Unix() + dis_mid
			case 2:
				self.SendChapterevents(1)
				self.SendGuildId(1)
				self.OperTime = time.Now().Unix() + dis_mid
			case 3:
				self.SendTeamPowerReq()
				self.OperTime = time.Now().Unix() + dis_mid
			case 4:
				self.SendCreateRole()
				self.OperTime = time.Now().Unix() + dis_mid
			case 5:
				self.SendSetCamp(GetRobotMgr().RobotCon.Camp)
				self.OperTime = time.Now().Unix() + dis_mid
			case 6:
				self.SendGuildId(4)
				self.OperTime = time.Now().Unix() + dis_mid
			case 7:
				self.Send("", C2S_Uid{Ctrl: "armsarenainfo", Uid: self.Uid})
				self.OperTime = time.Now().Unix() + dis_mid
			case 8:
				self.Send("", C2S_Uid{Ctrl: "getarmsarenafight", Uid: self.Uid})
				self.OperTime = time.Now().Unix() + dis_mid
			case 9:
				self.SendChapterevents(1)
				self.OperTime = time.Now().Unix() + dis_mid
			case 10:
				self.EnterBigMap(1)
				self.OperTime = time.Now().Unix() + dis_high
			case 11:
				self.SendGuildId(6)
				self.OperTime = time.Now().Unix() + dis_high
			case 12:
				self.GetPassBegin(110101)
				self.OperTime = time.Now().Unix() + dis_mid
			case 13:
				self.EnterBigMap(1)
				self.OperTime = time.Now().Unix() + dis_high
			case 14:
				self.SendGuildId(7)
				self.OperTime = time.Now().Unix() + dis_high
			case 15:
				self.SendFinishEvents(105, 175)
				self.OperTime = time.Now().Unix() + dis_mid
			case 16:
				self.SendPassEnd(110101, 0, 0, 0, 0)
				self.OperTime = time.Now().Unix() + dis_mid
			case 17:
				self.SendGuildId(8)
				self.OperTime = time.Now().Unix() + dis_high
			case 18:
				self.EnterBigMap(1)
				self.OperTime = time.Now().Unix() + dis_high
			case 19:
				self.SendGuildId(10)
				self.OperTime = time.Now().Unix() + dis_high
			case 20:
				fmt.Println("csvId =",self.csvId)
				self.SendSetRedicon(1)
				self.OperTime = time.Now().Unix() + dis_low
			case 21:
				self.SendGuildId(14)
				self.OperTime = time.Now().Unix() + dis_high
			case 22:
				self.SendActivateStar(1, 3004)
				self.OperTime = time.Now().Unix() + dis_mid
			case 23:
				self.SendGuildId(17)
				self.OperTime = time.Now().Unix() + dis_high
			case 24:
				self.SendGetMission(110103)
				self.OperTime = time.Now().Unix() + dis_low
			case 25:
				self.SendGuildId(20)
				self.OperTime = time.Now().Unix() + dis_high
			case 26:
				self.SendUpStarAuto(3004)
				self.OperTime = time.Now().Unix() + dis_mid
			case 27:
				self.SendGuildId(25)
				self.OperTime = time.Now().Unix() + dis_high
			case 28:
				self.SendFinishEvents(110, 180)
				self.OperTime = time.Now().Unix() + dis_mid
			case 29:
				self.GetPassWin(110103)
				self.OperTime = time.Now().Unix() + dis_mid
			case 30:
				self.SendSetRedicon(3)
				self.SendGuildId(29)
				self.OperTime = time.Now().Unix() + dis_high
			case 31:
				self.GetPassBegin(110104)
				self.OperTime = time.Now().Unix() + dis_mid
			case 32:
				self.EnterBigMap(1)
				self.OperTime = time.Now().Unix() + dis_mid
			case 33:
				self.SendFinishEvents(120, 185)
				self.OperTime = time.Now().Unix() + dis_mid
			case 34:
				self.SendPassEnd(110104, 0, 0, 0, 0)
				self.OperTime = time.Now().Unix() + dis_mid
			case 35:
				self.SendGuildId(32)
				self.OperTime = time.Now().Unix() + dis_high
			case 36:
				self.SendFinishEvents(106, 83)
				self.OperTime = time.Now().Unix() + dis_mid
			case 37:
				self.SendGuildId(33)
				self.SendSetRedicon(7)
				self.OperTime = time.Now().Unix() + dis_high
			case 38:
				self.SendGuildId(38)
				self.OperTime = time.Now().Unix() + dis_high
			case 39:
				self.SendEquipAction(1, 3004, 1, 0, 3, 1, 3, 20000013)
				self.SendGuildId(45)
				self.OperTime = time.Now().Unix() + dis_high
			case 40:
				fmt.Println("csvId =",self.csvId)
				self.SendEquipAction(1, 3, 1, 0, 0, 0, 5, 20000013)
				self.OperTime = time.Now().Unix() + dis_low
			case 41:
				self.SendSetRedicon(131079)
				self.OperTime = time.Now().Unix() + dis_low
			case 42:
				self.SendGuildId(46)
				self.SendEquipAction(1, 3, 1, 0, 0, 0, 6, 20000013)
				self.OperTime = time.Now().Unix() + dis_mid
			case 43:
				self.SendSetRedicon(131075)
				self.SendGuildId(49)
				self.OperTime = time.Now().Unix() + dis_high
			case 44:
				self.GetPassBegin(110105)
				self.OperTime = time.Now().Unix() + dis_mid
			case 45:
				self.EnterBigMap(1)
				self.OperTime = time.Now().Unix() + dis_mid
			case 46:
				self.SendFinishEvents(140, 195)
				self.OperTime = time.Now().Unix() + dis_mid
			case 47:
				self.SendPassEnd(110105, 0, 0, 0, 0)
				self.SendGuildId(49)
				self.OperTime = time.Now().Unix() + dis_high
			case 48:
				self.EnterBigMap(1)
				self.SendGuildId(54)
				self.OperTime = time.Now().Unix() + dis_mid
			case 49:
				self.SetBossAction(1,1)
				self.SendGuildId(55)
				self.OperTime = time.Now().Unix() + dis_mid
			case 50:
				self.SendAddteampos(6, 5001, 1)
				self.OperTime = time.Now().Unix() + dis_mid
			case 51:
				self.EnterBigMap(1)
				self.Send("", C2S_Uid{Ctrl: "armsarenainfo", Uid: self.Uid})
				self.OperTime = time.Now().Unix() + dis_mid
			case 52:
				self.Send("", C2S_Uid{Ctrl: "getarmsarenafight", Uid: self.Uid})
				self.OperTime = time.Now().Unix() + dis_mid
			case 53:
				self.GetFriendCommend()
				self.OperTime = time.Now().Unix() + dis_low
			case 54:
				self.Send("", C2S_Uid{Ctrl: "campfightinfolist", Uid: self.Uid})
				self.OperTime = time.Now().Unix() + dis_low
			case 55:
				self.GetStatistcsinfo()
				self.OperTime = time.Now().Unix() + dis_low
			case 56:
				self.Send("", C2S_Uid{Ctrl: "campunioninfo", Uid: self.Uid})
				self.OperTime = time.Now().Unix() + dis_low
			case 57:
				self.Send("", C2S_Uid{Ctrl: "get_dreamland_base_info", Uid: self.Uid})
				self.OperTime = time.Now().Unix() + dis_low
			case 58:
				self.Send("", C2S_Uid{Ctrl: "getactivitymask", Uid: self.Uid})
				self.OperTime = time.Now().Unix() + dis_low
			case 59:
				self.SendGuildId(60)
				self.SendChapterevents(2)
				self.OperTime = time.Now().Unix() + dis_low
			case 60:
				fmt.Println("csvId =",self.csvId)
				self.EnterBigMap(1)
				self.SendGuildId(62)
				self.OperTime = time.Now().Unix() + dis_high
			case 61:
				self.GetPassBegin(110201)
				self.EnterBigMap(1)
				self.OperTime = time.Now().Unix() + dis_mid
			case 62:
				self.SendFinishEvents(201, 205)
				self.OperTime = time.Now().Unix() + dis_low
			case 63:
				self.SendPassEnd(110201, 0, 0, 0, 0)
				self.OperTime = time.Now().Unix() + dis_high
			case 64:
				self.SendGuildId(68)
				self.EnterBigMap(1)
				self.OperTime = time.Now().Unix() + dis_high
			case 65:
				self.GetPassBegin(110202)
				self.SendFinishEvents(203, 210)
				self.SendPassEnd(110202, 0, 0, 0, 0)
				self.SendSetRedicon(131079)
				self.SendGuildId(69)
				self.OperTime = time.Now().Unix() + dis_high
			case 66:
				self.EnterBigMap(1)
				self.SendFinishEvents(203, 172)
				self.OperTime = time.Now().Unix() + dis_low
			case 67:
				self.Send("", C2S_Uid{Ctrl: "armsarenainfo", Uid: self.Uid})
				self.Send("", C2S_Uid{Ctrl: "getarmsarenafight", Uid: self.Uid})
				self.GetFriendCommend()
				self.Send("", C2S_Uid{Ctrl: "campfightinfolist", Uid: self.Uid})
				self.GetStatistcsinfo()
				self.Send("", C2S_Uid{Ctrl: "campunioninfo", Uid: self.Uid})
				self.Send("", C2S_Uid{Ctrl: "get_dreamland_base_info", Uid: self.Uid})
				self.Send("", C2S_Uid{Ctrl: "getactivitymask", Uid: self.Uid})
				self.SendGuildId(71)
				self.OperTime = time.Now().Unix() + dis_high
			case 68:
				self.Find(1, 1)
				self.SendGuildId(74)
				self.OperTime = time.Now().Unix() + dis_low
			case 69:
				self.Find(3, 1)
				self.SendGuildId(76)
				self.OperTime = time.Now().Unix() + dis_low
			case 70:
				self.SendSetRedicon(131074)
				self.SendGuildId(78)
				self.OperTime = time.Now().Unix() + dis_high
			case 71:
				self.SendAddteampos(4, 1002, 1)
				self.SendGuildId(81)
				self.OperTime = time.Now().Unix() + dis_mid
			case 72:
				self.SendAddteampos(5, 4010, 1)
				self.SendGuildId(84)
				self.OperTime = time.Now().Unix() + dis_mid
			case 73:
				self.SendSwapFightPos(4, 2, 1)
				self.SendGuildId(87)
				self.OperTime = time.Now().Unix() + dis_mid
			case 74:
				self.SendChapterevents(2)
				self.EnterBigMap(1)
				self.SendGuildId(89)
				self.OperTime = time.Now().Unix() + dis_high
			case 75:
				self.GetPassBegin(110203)
				self.EnterBigMap(1)
				self.SendFinishEvents(202, 229)
				self.SendPassEnd(110203, 0, 0, 0, 0)
				self.SendGuildId(90)
				self.OperTime = time.Now().Unix() + dis_high
			case 76:
				self.EnterBigMap(1)
				self.SendGuildId(92)
				self.OperTime = time.Now().Unix() + dis_mid
			case 77:
				self.Send("", C2S_Uid{Ctrl: "nobilitytaskinfo", Uid: self.Uid})
				self.SendGuildId(93)
				self.OperTime = time.Now().Unix() + dis_mid
			case 78:
				self.SendTakeNobilitytask(1001)
				self.SendGuildId(97)
				self.OperTime = time.Now().Unix() + dis_mid
			case 79:
				self.SendTakeNobilitytask(1002)
				self.SendGuildId(99)
				self.SendGuildId(100003)
				self.OperTime = time.Now().Unix() + dis_mid
			case 80:
				fmt.Println("over！！！")
				self.Dead = true
			}
			self.csvId++
			//log.Println("下一步操作:", self.msgId)
		}
	}
}
