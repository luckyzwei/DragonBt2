package main

import (
	"log"
	"time"
)

// 按照csv流程跑的机器人
func (self *Robot) csvRun() {
	ticker := time.NewTicker(time.Second)

	self.csvId=1
	robotnromal := GetRobotCsvMgr().Data["robotnormal"]
	csv_max := len(robotnromal)
	for {
		<-ticker.C
		if self.Dead {
			break
		}

		if self.OperTime == 0 {
			continue
		}


		if time.Now().Unix() >= self.OperTime {
			//不需要下线在上，所以过滤掉上线环节  如果要包括上下线环节则（self.csvId>csv_max,self.csvId=1）
			if self.csvId >= csv_max {
				if GetRobotMgr().RobotCon.IsLoop {
					self.csvId = 2
					log.Println("机器人:", self.Uname, "uid:", self.Uid, "开始循环执行")
				} else {
					self.Dead=true
					continue
				}
			}

			csv_data, ok := robotnromal[self.csvId]

			if !ok {
				log.Println("robotnormal err", self.csvId)
				return
			}

			if self.csvGoAct(csv_data) {
				self.csvId++
			}
		}
	}
}

// deprecated
func (self *Robot) csvGoAct(csvdata CsvNode) bool {
	acttype := HF_Atoi(csvdata["acttype"])
	log.Println("机器人:", self.Uname, "uid:", self.Uid, "执行命令,类型", acttype, "csvid", self.csvId)
	ret_ok := true
	switch acttype {
	case Sendlogin:
		//ws, err := websocket.Dial(self.ServerUrl, "", self.origin)
		//if err != nil {
		//	log.Println(fmt.Sprintf("机器人连接服务器：%s失败", self.ServerUrl))
		//	return false
		//}
		//self.Ws = ws

		self.IsLoginOut = false
		self.SendLogin()
		//go self.Recive()
	case Passmission:
		missionid := HF_Atoi(csvdata["var1"])
		step := HF_Atoi(csvdata["var2"])
		self.SendMission(missionid, step, 0, 0)
	case Passlevelbeggin:
		passid := HF_Atoi(csvdata["var1"])
		self.GetPassBegin(passid)
	case Passlevelend:
		passid := HF_Atoi(csvdata["var1"])
		missionid := HF_Atoi(csvdata["var2"])
		step := HF_Atoi(csvdata["var3"])
		self.SendPassEnd(passid, missionid, step, 0, 0)
	case Getfriendcommend:
		self.GetFriendCommend()
	case Sendfriendapply:
		self.SendFrienDapply()
	case Sendfriendorder:
		self.SendFrienDorder()
	case Chatsystem:
		channel := HF_Atoi(csvdata["var1"])
		strDec := csvdata["var2"]
		strname := csvdata["var3"]

		if strname == "0" {
			strname = ""
		}
		self.SendChat(channel, strDec, strname)
	case DanMusystem:
		strDec := csvdata["var1"]
		self.SendDanMu(strDec)
	case Createitem:
		itemid := HF_Atoi(csvdata["var1"])
		itemnum := HF_Atoi(csvdata["var2"])
		self.CreateItem(itemid, itemnum)
	case GmPassid:
		passid := HF_Atoi(csvdata["var1"])
		self.SendPassId(passid)
		/*
	case ShenqiLvUp:
		heroid := HF_Atoi(csvdata["var1"])
		sjtype := HF_Atoi(csvdata["var2"])
		zhuangbei_id := HF_Atoi(csvdata["var3"])
		self.SendSenqiLvUp(heroid, sjtype, zhuangbei_id)
		*/
		/*
	case HeroColorUp:
		heroid := HF_Atoi(csvdata["var1"])
		self.SendUpColor(heroid)
		 */
		/*
	case HeroSatrUp:
		heroid := HF_Atoi(csvdata["var1"])
		self.SendHeroStarUp(heroid)
		 */
		/*
	case SoldierLvUp:
		heroid := HF_Atoi(csvdata["var1"])
		self.SendSoldierLvUp(heroid)
		 */
	case LuckDraw:
		drawindx := HF_Atoi(csvdata["var1"])
		drawtype := HF_Atoi(csvdata["var2"])
		self.Find(drawindx,drawtype)
		/*
	case CampFightMove:
		self.SendCampFightMove()
	case CampFightPlayList:
		cityPart := HF_Atoi(csvdata["var1"])
		self.GetCampfightPlayList(0, cityPart)
	case CampFightSolo2:

		fightidx := HF_Atoi(csvdata["var1"])
		if self.is_fightinfo {

			self.SendCampFightSolo2(self.fightinfo[fightidx].Index, self.fightinfo[fightidx].Uid, 0)

		} else {
			ret_ok = false
		}

	case CampFightSoloEnd:
		jieguo := HF_Atoi(csvdata["var1"])
		self.SendCampFightSoloEnd(jieguo, 0)
	case CampFight56Req:
		self.SendCampFight56Req(1200, 0)
		self.SendCampFight56Req(12002, 0)
		 */
	/*
	case HeroLvUp:
		heroid := HF_Atoi(csvdata["var1"])
		itemid := csvdata["var2"]
		itemnum := HF_Atoi(csvdata["var3"])
		self.SendHeroLvUp(heroid, itemid, itemnum)
	*/
	case SendLoginOut:
		log.Println("机器人:", self.Uname, "uid:", self.Uid, "登出游戏")
		self.IsLoginOut = true
		self.Ws.Close()
	case ConsumerTopAttack:
		var msg C2S_Uid
		msg.Ctrl = "consumertopattack"
		msg.Uid = self.Uid
		self.Send("", &msg)
	}

	if ret_ok {
		starttime := HF_Atoi64(csvdata["starttime"])
		endtime := HF_Atoi64(csvdata["endtime"])
		timeinterval := HF_Atoi64(csvdata["timeinterval"])

		if starttime != 0 && endtime != 0 {
			timeinterval = RandInt64(starttime, endtime)
			log.Println("随机时间:", timeinterval)
		}
		self.OperTime = time.Now().Unix() + timeinterval

	} else {
		self.OperTime = time.Now().Unix() + 3
	}
	return ret_ok

}