package game

import (
	"encoding/json"
	"fmt"
)

//! 用户消费结构
type San_UserConsumerTop struct {
	Uid       int64  //! Uid
	Rank      int    //! 排名-跨服
	Point     int    //! 积分
	Level     int    //! 进攻等级
	RankAward int    //! 排行奖励
	Step      int    //! 活动ID
	Record    string //! 伤害记录
	Award     string //! 领奖记录
	EndTime   int64  //! 过期时间，超过过期时间重置活动

	record []JS_DamageRecord //! 伤害记录
	award  map[int][]int     //! 领取积分奖励记录

	DataUpdate
}

//! 伤害记录
type JS_DamageRecord struct {
	Time int64 `json:"time"` //! 时间
	Dps  int   `json:"dps"`  //! 伤害
}

//! 消费者排行榜
type ModConsumerTop struct {
	player   *Player
	Sql_Data San_UserConsumerTop //! 数据库结构
}

func (self *ModConsumerTop) OnGetData(player *Player) {
	self.player = player
}

func (self *ModConsumerTop) OnGetOtherData() {
	sql := fmt.Sprintf("select * from `san_userconsumertop` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Data, "", self.player.ID)

	if self.Sql_Data.Uid <= 0 {
		self.Sql_Data.Uid = self.player.ID
		self.Sql_Data.record = make([]JS_DamageRecord, 0)
		self.Sql_Data.award = make(map[int][]int, 0)
		self.Sql_Data.Level = 1
		self.Sql_Data.EndTime = 0

		self.Encode()
		InsertTable("san_userconsumertop", &self.Sql_Data, 0, false)
	} else {
		self.Decode()
	}

	if self.Sql_Data.Level == 0 {
		self.Sql_Data.Level = 1
	}
	if self.Sql_Data.award == nil {
		self.Sql_Data.award = make(map[int][]int, 0)
		if self.Sql_Data.Award != "[]" {
			var award []int
			json.Unmarshal([]byte(self.Sql_Data.Award), &award)
			self.Sql_Data.award[self.Sql_Data.Step] = award
		} else {
			self.Sql_Data.award = make(map[int][]int, 0)
		}
	}

	if self.Sql_Data.award[self.Sql_Data.Step] == nil {
		self.Sql_Data.award[self.Sql_Data.Step] = make([]int, 0)
	}

	self.Sql_Data.Init("san_userconsumertop", &self.Sql_Data, false)

	topuser := GetConsumerTop().GetTopUser(self.player.Sql_UserBase.Uid)
	if topuser != nil {
		self.Sql_Data.Rank = topuser.Rank
	}

	if GetConsumerTop().IsValid() == true {
		mhero := GetConsumerTop().GetMHero()
		if mhero != nil {
			if mhero.Step != self.Sql_Data.Step {
				if self.Sql_Data.EndTime == mhero.Endtime {
					self.Sql_Data.Step = mhero.Step
				} else {
					self.Sql_Data.Level = 1
					self.Sql_Data.Point = 0
					self.Sql_Data.Rank = 0
					//self.Sql_Data.award = make(map[int][]int, 0)
					self.Sql_Data.record = make([]JS_DamageRecord, 0)
					self.Sql_Data.EndTime = mhero.Endtime
					self.Sql_Data.Step = mhero.Step
				}
			}
		}
	}
}

func (self *ModConsumerTop) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "consumertopinfo":
		self.SendInfo()
		return true
	case "consumertopuser":
		var msg C2S_Collection
		json.Unmarshal(body, &msg)
		self.GetTopList()
		return true
	case "consumertopserver":
		var msg C2S_Collection
		json.Unmarshal(body, &msg)
		self.GetServerTopList()
		return true
	case "consumertopattack":
		var msg C2S_Collection
		json.Unmarshal(body, &msg)
		self.AttackHero()
		return true
	case "consumertopdraw":
		var msg C2S_ConsumerTopDraw
		json.Unmarshal(body, &msg)
		self.Draw(msg.Id)
		return true
	case "consumertopenter":
		self.SendInfo()
		GetConsumerTop().AddPlayer(self.player.Sql_UserBase.Uid, false)
		return true
	case "consumertopleave":
		GetConsumerTop().AddPlayer(self.player.Sql_UserBase.Uid, true)
		return true
	}

	return false
}

//! 保存
func (self *ModConsumerTop) OnSave(sql bool) {
	self.Encode()
	self.Sql_Data.Update(sql)
}

//! 将数据库数据写入data
func (self *ModConsumerTop) Decode() {

	json.Unmarshal([]byte(self.Sql_Data.Record), &self.Sql_Data.record)
	json.Unmarshal([]byte(self.Sql_Data.Award), &self.Sql_Data.award)
}

//! 将data数据写入数据库
func (self *ModConsumerTop) Encode() {
	self.Sql_Data.Record = HF_JtoA(self.Sql_Data.record)
	self.Sql_Data.Award = HF_JtoA(self.Sql_Data.award)
}

func (self *ModConsumerTop) OnRefresh() {
	self.Sql_Data.Level = 1

	//! 活动结束时，清空数据
	if GetConsumerTop().IsValid() == false {
		self.CheckAward()
	}
}

//! 排行榜-个人
func (self *ModConsumerTop) GetTopList() {
	var msg S2C_ConsumerTopUser
	msg.Cid = "consumertopuser"
	msg.Rank = self.Sql_Data.Rank
	msg.TopUser = GetConsumerTop().GetTopUserList()

	for i := 0; i < len(msg.TopUser); i++ {
		if msg.TopUser[i].Uid == self.player.Sql_UserBase.Uid {
			self.Sql_Data.Rank = msg.TopUser[i].Rank
			msg.Rank = self.Sql_Data.Rank
			break
		}
	}

	self.player.SendMsg("consumertopuser", HF_JtoB(&msg))
}

//! 排行榜-服务器
func (self *ModConsumerTop) GetServerTopList() {
	var msg S2C_ConsumerTopServer
	msg.Cid = "consumertopserver"
	//msg.Rank = self.Sql_Data.Rank
	msg.SvrId = GetServer().Con.ServerId
	msg.TopSvr = GetConsumerTop().GetTopServerList()
	for i := 0; i < len(msg.TopSvr); i++ {
		if msg.TopSvr[i].SvrId == GetServer().Con.ServerId {
			msg.Rank = i + 1
			break
		}
	}

	self.player.SendMsg("consumertopserver", HF_JtoB(&msg))
}

//! 攻击神将-获取道具
func (self *ModConsumerTop) AttackHero() {
	if GetConsumerTop().IsValid() == false {
		self.player.SendRet("consumertopattack", 1)
		//self.player.SendErrInfo("err", "活动已经结束")
		return
	}

	if GetConsumerTop().IsAttackable() == false {
		self.player.SendRet("consumertopattack", 2)
		//self.player.SendErrInfo("err", "当前无法攻击")
		return
	}

	if self.Sql_Data.Level == 0 {
		self.Sql_Data.Level = 1
	}

	csvattack, ok := GetCsvMgr().Data["Consumetop_Attack"][self.Sql_Data.Level]
	if !ok {
		self.player.SendRet("consumertopattack", 3)
		//self.player.SendErrInfo("err", "配置表读取错误")
		return
	}

	maxLevel := GetCsvMgr().GetVipConsumeTop(self.Sql_Data.Level)
	if maxLevel > 0 && self.Sql_Data.Level > maxLevel {
		self.player.SendRet("consumertopattack", 5)
		return
	}

	itemlst := make([]PassItem, 0)
	//!消耗
	if HF_Atoi(csvattack["cost"]) > 0 {
		if self.player.GetObjectNum(DEFAULT_GEM) >= HF_Atoi(csvattack["cost"]) {
			itemlst = append(itemlst, PassItem{DEFAULT_GEM, -HF_Atoi(csvattack["cost"])})
		} else {
			self.player.SendRet("consumertopattack", 4)
			//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_NOT_ENOUGH_GEM_FAIL"))
			return
		}
	}

	//! 奖励
	//for i := 0; i < 2; i++ {
	//	itemlst = append(itemlst, PassItem{
	//		HF_Atoi(csvattack[fmt.Sprintf("item%d", i+1)]),
	//		HF_Atoi(csvattack[fmt.Sprintf("num%d", i+1)])})
	//}
	//! 获得灭世币，item2+num2
	itemlst = append(itemlst, PassItem{HF_Atoi(csvattack["item2"]), HF_Atoi(csvattack["num2"])})
	itemlst = append(itemlst, PassItem{HF_Atoi(csvattack["item3"]), HF_Atoi(csvattack["num3"])})

	mhero := GetConsumerTop().GetMHero()
	//! 攻击等级每次+1
	self.Sql_Data.Level += 1
	// 获得积分
	point := HF_Atoi(csvattack["num1"])
	self.Sql_Data.Point += point //! 获得灭世币

	//! 添加积分到道具列表-用做显示
	itemlst = append(itemlst, PassItem{HF_Atoi(csvattack["item1"]), HF_Atoi(csvattack["num1"])})

	oldrank := self.Sql_Data.Rank
	//! 战斗数据保存
	kill, rank := GetConsumerTop().AddDamage(self.player, point, self.Sql_Data.Point, point)
	self.Sql_Data.Rank = rank

	nAttack := 0
	if kill == 1 {
		nAttack = 1
	} else {
		nAttack = 0
	}

	for i := 0; i < len(itemlst); i++ {
		self.player.AddObject(itemlst[i].ItemID, itemlst[i].Num, nAttack, 0, 0, "诸神黄昏战斗")
	}

	dropitem, rare := self.GetDropItem()
	for i := 0; i < len(dropitem); i++ {
		dropitem[i].ItemID, dropitem[i].Num = self.player.AddObject(dropitem[i].ItemID, dropitem[i].Num, nAttack, 0, 0, "诸神黄昏战斗")
		itemlst = append(itemlst, PassItem{dropitem[i].ItemID, dropitem[i].Num})

		if rare == true {
			GetServer().sendSysChat(fmt.Sprintf(GetCsvMgr().GetText("STR_CONSUME_AWARD"),
				HF_GetColorByCamp(self.player.Sql_UserBase.Camp), self.player.Sql_UserBase.UName,
				HF_GetColorByCamp(2), GetCsvMgr().GetItemName(dropitem[i].ItemID)))
		}
	}

	if oldrank != rank {
		GetServer().SqlLog(self.player.GetUid(), LOG_CONSUMER_CHANGE, rank, oldrank, 0, "诸神黄昏排名变化", 0, 0, self.player)
	}

	//! 击杀掉落-通过邮件发送奖励
	if kill == 1 {
		if herocsv, ok := GetCsvMgr().Data["Consumetop_Hp"][mhero.Level]; ok {
			lstItem := make([]PassItem, 0)
			for j := 0; j < 3; j++ {
				itemid := HF_Atoi(herocsv[fmt.Sprintf("item%d", j+1)])
				if itemid == 0 {
					continue
				}
				lstItem = append(lstItem, PassItem{itemid, HF_Atoi(herocsv[fmt.Sprintf("num%d", j+1)])})
			}
			info := fmt.Sprintf(GetCsvMgr().GetText("STR_CONSUME_SPECIAL_AWARD"), herocsv["level"], mhero.Name)
			self.player.GetModule("mail").(*ModMail).AddMail(1, 1, 0, GetCsvMgr().GetText("STR_CONSUME_SUCCESS_KILL"), info, GetCsvMgr().GetText("STR_SYS"), lstItem, true, 0)
		}

		GetServer().SqlLog(self.player.GetUid(), LOG_CONSUMER_FIGHT, 1, kill, 0, "诸神黄昏战斗", 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.GetUid(), LOG_CONSUMER_FIGHT, 1, 0, 0, "诸神黄昏战斗", 0, 0, self.player)
	}

	var msg S2C_ConsumerAttackHero
	msg.Cid = "consumertopattack"
	msg.Ret = 0
	msg.Kill = kill
	msg.Damage = point
	msg.Hp = mhero.HP
	msg.MaxHP = mhero.MaxHP
	msg.MHeroLevel = mhero.Level
	msg.Level = self.Sql_Data.Level
	msg.Point = self.Sql_Data.Point
	msg.Item = itemlst

	self.player.SendMsg("consumertopattack", HF_JtoB(&msg))
}

func (self *ModConsumerTop) GetDropItem() ([]PassItem, bool) {
	outitem := make([]PassItem, 0)

	total := 0
	rare := false
	//lstdropgroup := GetCsvMgr().Data2["Consumetop_Luck"]
	for i := 0; i < len(GetCsvMgr().Data2["Consumetop_Luck"]); i++ {
		total += HF_Atoi(GetCsvMgr().Data2["Consumetop_Luck"][i]["change"])
	}
	if total == 0 {
		return outitem, rare
	}

	pro := HF_GetRandom(total)
	cur := 0
	var dropcsv CsvNode = nil
	for i := 0; i < len(GetCsvMgr().Data2["Consumetop_Luck"]); i++ {
		cur += HF_Atoi(GetCsvMgr().Data2["Consumetop_Luck"][i]["change"])
		if pro < cur {
			dropcsv = GetCsvMgr().Data2["Consumetop_Luck"][i]
			if HF_Atoi(GetCsvMgr().Data2["Consumetop_Luck"][i]["change"]) < 10 {
				rare = true
			}
			break
		}
	}

	if dropcsv != nil {
		total := 0
		for i := 0; i < 6; i++ {
			total += HF_Atoi(dropcsv[fmt.Sprintf("value%d", i+1)])
		}

		if total == 0 {
			return outitem, rare
		}
		pro = HF_GetRandom(total)
		cur = 0

		for i := 0; i < 6; i++ {
			cur += HF_Atoi(dropcsv[fmt.Sprintf("value%d", i+1)])
			if pro < cur {
				outitem = append(outitem, PassItem{HF_Atoi(dropcsv[fmt.Sprintf("item%d", i+1)]),
					HF_Atoi(dropcsv[fmt.Sprintf("num%d", i+1)])})

				break
			}
		}
	}

	return outitem, rare
}

//! 领取奖励-积分
func (self *ModConsumerTop) Draw(id int) {
	csv, ok := GetCsvMgr().Data["Consumetop_List"][id]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_CANT_FIND_SETUP"))
		return
	}

	getitem := make([]PassItem, 0)
	//! 积分活动
	if id > 200000 {
		if self.Sql_Data.Point < HF_Atoi(csv["n1"]) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_CONSUMERTOP_NOT_AWARD"))
			return
		}
	} else {
		activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
		if activity == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ACT_NOT_OPEN"))
			return
		} else {
			if activity.status.Status != ACTIVITY_STATUS_SHOW {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_CANT_NOT_GET"))
				return
			}
		}

		if GetConsumerTop().Sql_Data.hero.RankGroup != HF_Atoi(csv["group"]) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_CONSUMERTOP_NOT_AWARD"))
			return
		}

		if self.Sql_Data.Rank < HF_Atoi(csv["n1"]) || self.Sql_Data.Rank > HF_Atoi(csv["n2"]) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_CONSUMERTOP_NOT_AWARD"))
			return
		}

		for getid := range GetCsvMgr().Data["Consumetop_List"] {
			if getid < 200000 {
				for i := 0; i < len(self.Sql_Data.award[self.Sql_Data.Step]); i++ {
					if self.Sql_Data.award[self.Sql_Data.Step][i] == getid {
						//! 已经领过奖励，无法在领取
						self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_CONSUMERTOP_GOT_AWARD"))
						return
					}
				}
			}
		}
	}

	reason := ""

	if HF_Atoi(csv["type"]) == 1 {
		reason = "诸神黄昏排名奖励"
	} else {
		reason = "诸神黄昏积分奖励"
	}

	for i := 0; i < len(self.Sql_Data.award[self.Sql_Data.Step]); i++ {
		if self.Sql_Data.award[self.Sql_Data.Step][i] == id {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_CONSUMERTOP_NOT_AWARD"))
			return
		}
	}

	for i := 0; i < 4; i++ {
		itemid := HF_Atoi(csv[fmt.Sprintf("item%d", i+1)])
		num := HF_Atoi(csv[fmt.Sprintf("num%d", i+1)])

		if itemid > 0 {
			getitem = append(getitem, PassItem{itemid, num})

			self.player.AddObject(itemid, num, 0, 0, 0, reason)
		}
	}

	if id < 200000 {
		if self.Sql_Data.Point >= HF_Atoi(csv["n3"]) {
			for i := 4; i < 6; i++ {
				itemid := HF_Atoi(csv[fmt.Sprintf("item%d", i+1)])
				num := HF_Atoi(csv[fmt.Sprintf("num%d", i+1)])

				if itemid > 0 {
					getitem = append(getitem, PassItem{itemid, num})

					self.player.AddObject(itemid, num, 0, 0, 0, reason)
				}
			}
		}
	}

	self.Sql_Data.award[self.Sql_Data.Step] = append(self.Sql_Data.award[self.Sql_Data.Step], id)

	var msg S2C_ConsumerTopDraw
	msg.Cid = "consumertopdraw"
	msg.Id = id
	msg.Award = self.Sql_Data.award[self.Sql_Data.Step]
	msg.Item = getitem

	self.player.SendMsg("consumertopdraw", HF_JtoB(&msg))

}

//! 发送数据
func (self *ModConsumerTop) SendInfo() {
	///! 同步当前的排名
	rank := GetConsumerTop().GetUserRank(self.player.Sql_UserBase.Uid)
	if rank != 0 {
		self.Sql_Data.Rank = rank
	}

	//! 判断是否当前期数
	if GetConsumerTop().IsValid() == true {
		mhero := GetConsumerTop().GetMHero()
		if mhero != nil {
			if mhero.Step != self.Sql_Data.Step {
				if self.Sql_Data.EndTime == mhero.Endtime {
					self.Sql_Data.Step = mhero.Step
				} else {
					self.Sql_Data.Level = 1
					self.Sql_Data.Point = 0
					self.Sql_Data.Rank = 0
					//self.Sql_Data.award = make(map[int][]int, 0)
					self.Sql_Data.record = make([]JS_DamageRecord, 0)
					self.Sql_Data.EndTime = mhero.Endtime
					self.Sql_Data.Step = mhero.Step
				}
			}
		}
	} else {
		self.CheckAward()
	}

	var msg S2C_ConsumerTopInfo
	msg.Cid = "consumertopinfo"
	msg.Level = self.Sql_Data.Level
	msg.Rank = self.Sql_Data.Rank
	msg.ServerRank = GetConsumerTop().GetServerRank(GetServer().Con.ServerId)
	msg.ServerName = GetServer().Con.ServerName
	msg.Award = self.Sql_Data.award[self.Sql_Data.Step]
	msg.Point = self.Sql_Data.Point
	msg.Record = self.Sql_Data.record
	msg.Msg = make([]*JS_ConsumerMsg, 0)
	if GetConsumerTop().GetMHero() != nil {
		msg.MHero = append(msg.MHero, *(GetConsumerTop().GetMHero()))
		msg.Msg = GetConsumerTop().Sql_Data.msg
	} else {
		msg.MHero = make([]JS_MagicialHero, 0)
	}

	if msg.Msg == nil {
		msg.Msg = make([]*JS_ConsumerMsg, 0)
	}

	msg.BossConfig = GetCsvMgr().ConsumetopbossConfig
	msg.ListConfig1 = GetCsvMgr().GetConsumetoplistConfig1(GetConsumerTop().GetMHero().RankGroup)
	msg.ListConfig2 = GetCsvMgr().GetConsumetoplistConfig2()
	msg.ShopConfig = GetCsvMgr().GetConsumetopshopConfig(GetConsumerTop().GetMHero().ShopGroup)
	msg.BossFightInfo = GetConsumerTop().GetBossFightInfo()

	self.player.SendMsg("consumertopinfo", HF_JtoB(&msg))
}

func (self *ModConsumerTop) CheckAward() {
	//活动关了之后，只检查当然身上的任务是否有领取
	itemMap := make(map[int]*Item)
	for getid, v := range GetCsvMgr().Data["Consumetop_List"] {
		if getid < 200000 {
			continue
		}
		if self.Sql_Data.Point < HF_Atoi(v["n1"]) {
			continue
		}
		isAlready := false
		for i := 0; i < len(self.Sql_Data.award[self.Sql_Data.Step]); i++ {
			if self.Sql_Data.award[self.Sql_Data.Step][i] == getid {
				isAlready = true
				break
			}
		}
		if isAlready {
			continue
		}
		for i := 0; i < 4; i++ {
			itemid := HF_Atoi(v[fmt.Sprintf("item%d", i+1)])
			num := HF_Atoi(v[fmt.Sprintf("num%d", i+1)])

			if itemid > 0 {
				AddItemMapHelper3(itemMap, itemid, num)
			}
		}
		self.Sql_Data.award[self.Sql_Data.Step] = append(self.Sql_Data.award[self.Sql_Data.Step], getid)
	}

	//发送邮件
	if len(itemMap) > 0 {
		mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_CONSUMERTOP_SCORE]
		if ok {
			itemLst := make([]PassItem, 0)
			for _, v := range itemMap {
				itemLst = append(itemLst, PassItem{ItemID: v.ItemId, Num: v.ItemNum})
			}
			self.player.GetModule("mail").(*ModMail).AddMail(1,
				1, 0, mailConfig.Mailtitle, mailConfig.Mailtxt, GetCsvMgr().GetText("STR_SYS"), itemLst, true, 0)
		}
	}

	//排行奖励
	itemMapRank := make(map[int]*Item)
	csvBoss, ok := GetCsvMgr().Data["Consumetop_Boss"][self.Sql_Data.Step/1000]
	if !ok {
		return
	}
	group := HF_Atoi(csvBoss["list"])
	for getid, v := range GetCsvMgr().Data["Consumetop_List"] {
		if getid > 200000 {
			continue
		}

		if group != HF_Atoi(v["group"]) {
			continue
		}

		if self.Sql_Data.Rank < HF_Atoi(v["n1"]) || self.Sql_Data.Rank > HF_Atoi(v["n2"]) {
			continue
		}
		for i := 0; i < len(self.Sql_Data.award[self.Sql_Data.Step]); i++ {
			if self.Sql_Data.award[self.Sql_Data.Step][i] == getid {
				return
			}
		}
		for i := 0; i < 4; i++ {
			itemid := HF_Atoi(v[fmt.Sprintf("item%d", i+1)])
			num := HF_Atoi(v[fmt.Sprintf("num%d", i+1)])

			if itemid > 0 {
				AddItemMapHelper3(itemMapRank, itemid, num)
			}
		}
		if self.Sql_Data.Point >= HF_Atoi(v["n3"]) {
			for i := 4; i < 6; i++ {
				itemid := HF_Atoi(v[fmt.Sprintf("item%d", i+1)])
				num := HF_Atoi(v[fmt.Sprintf("num%d", i+1)])

				if itemid > 0 {
					AddItemMapHelper3(itemMapRank, itemid, num)
				}
			}
		}

		if len(itemMapRank) > 0 {
			mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_CONSUMERTOP_RANK]
			if ok {
				itemLst := make([]PassItem, 0)
				for _, v := range itemMapRank {
					itemLst = append(itemLst, PassItem{ItemID: v.ItemId, Num: v.ItemNum})
				}
				self.player.GetModule("mail").(*ModMail).AddMail(1,
					1, 0, mailConfig.Mailtitle, fmt.Sprintf(mailConfig.Mailtxt, self.Sql_Data.Rank), GetCsvMgr().GetText("STR_SYS"), itemLst, true, 0)
			}
			self.Sql_Data.award[self.Sql_Data.Step] = append(self.Sql_Data.award[self.Sql_Data.Step], getid)
		}
	}
}
