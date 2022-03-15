package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const INT_MAXHORSE_NUM = 200
const INT_MAX_HORSE_PACKAGE = 4

const (
	SUMMON_COMMON   = 1 // 普通召唤
	SUMMON_COMPOUND = 2 // 合成碎片召唤
)

//! 魔宠数据结构
type San_Horse struct {
	Uid             int64
	Level           int    //! 马场等级
	Exp             int    //! 马场经验
	SummonTime      int64  //! 上次自动获得时间
	SummonNormal    int    //! 召唤次数
	SummonSenior    int    //! 高级召唤次数
	Discern         int    //! 鉴定次数
	Decompose       int    //! 分解等级
	Combine         int    //! 合成次数
	MaxHorseId      int    //! 最大Id
	Info            string //! 马匹信息1
	Soul            string //! 马魂信息
	Summon          string //! 召唤列表，1-野马，2-高级野马
	Info2           string //! 马匹信息2
	Info3           string //! 马匹信息3
	Info4           string //! 马匹信息4
	HorseTask       string //! 魔宠高级召唤任务
	HorseTotalFight int64  // !魔宠总战力

	summon    [2]int
	sinfo     [4]map[int]*JS_HorseInfo
	info      map[int]*JS_HorseInfo
	soul      map[int]*JS_HorseSoulInfo
	horseTask map[int]*HorseTask

	DataUpdate
}

// 魔宠召唤任务
type HorseTask struct {
	Id      int `json:"id"`
	Process int `json:"process"`
	Status  int `json:"status"`
}

//! 魔宠信息
type JS_HorseInfo struct {
	Id          int            `json:"id"`          //! 魔宠唯一Id
	Type        int            `json:"type"`        //! 魔宠配置Id
	Awaken      int            `json:"awaken"`      //! 觉醒
	Skill       int            `json:"skill"`       //! 技能数据
	Skilllv     int            `json:"skilllv"`     //! 技能等级
	Heroid      int            `json:"heroid"`      //! 英雄Id
	AttLst      []JS_HorseAttr `json:"attlst"`      //! 属性列表
	SoulLst     []int          `json:"soullst"`     //! 马魂列表-最多6个属性
	RandHorseId int            `json:"randhorseId"` //! 魔宠转换保留
	Fight       int64          `json:"fight"`       //! 战马战力

	chg bool //! 是否改变
}

//! 魔宠属性-动态
type JS_HorseAttr struct {
	AttrType  int   `json:"at"` //! 属性类型
	AttrValue int64 `json:"av"` //! 属性数值
}

//! 马魂信息
type JS_HorseSoulInfo struct {
	Id   int `json:"id"`   //! id，规则为type*100
	Rank int `json:"rank"` //! 阶级
	Num  int `json:"num"`  //! 数量
}

//! 抽奖
type ModHorse struct {
	player     *Player
	Sql_Horse  San_Horse //! 数据库结构
	MountNum   int       //! 使用数量
	DataLocker *sync.RWMutex
}

func (self *ModHorse) OnGetData(player *Player) {
	self.player = player
	self.DataLocker = new(sync.RWMutex)

	sql := fmt.Sprintf("select * from `san_userhorse` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Horse, "san_userhorse", self.player.ID)

	if self.Sql_Horse.Uid <= 0 {
		self.Sql_Horse.Uid = self.player.ID
		self.Sql_Horse.info = make(map[int]*JS_HorseInfo)
		for i := 0; i < INT_MAX_HORSE_PACKAGE; i++ {
			self.Sql_Horse.sinfo[i] = make(map[int]*JS_HorseInfo)
		}

		self.Sql_Horse.horseTask = make(map[int]*HorseTask)

		self.Sql_Horse.soul = make(map[int]*JS_HorseSoulInfo)
		self.Sql_Horse.Level = 1
		self.Sql_Horse.SummonNormal = 0
		self.Sql_Horse.Combine = 0
		self.Sql_Horse.SummonSenior = 0
		self.Sql_Horse.MaxHorseId = 1
		self.Encode()
		InsertTable("san_userhorse", &self.Sql_Horse, 0, true)
	} else {

		if self.Sql_Horse.soul == nil {
			self.Sql_Horse.soul = make(map[int]*JS_HorseSoulInfo)
		}

		self.Sql_Horse.info = make(map[int]*JS_HorseInfo)
		self.Decode()

		for _, horsedata := range self.Sql_Horse.info {
			if horsedata.Heroid == 0 && len(horsedata.SoulLst) > 0 {
				for j := 0; j < len(horsedata.SoulLst); j++ {
					if horsedata.SoulLst[j] > 0 {
						soulid := horsedata.SoulLst[j] / 100
						soulrank := horsedata.SoulLst[j] % 100

						self.AddHorseSoul(soulid, soulrank, 1)
					}
				}
				soulnum := len(horsedata.SoulLst)
				horsedata.SoulLst = make([]int, soulnum)
			}
		}
	}

	self.Sql_Horse.Init("san_userhorse", &self.Sql_Horse, true)
}

func (self *ModHorse) OnGetOtherData() {
	//! 检查是否由魔宠被销毁，数据残余
	self.player.GetModule("hero").(*ModHero).CheckHorse()

	self.MountNum = 0
	for _, horse := range self.Sql_Horse.info {
		if horse.Heroid > 0 {
			self.MountNum++
		}
	}

	//自动召唤
	self.AutoSummon()

	//if self.Sql_Horse.HorseTotalFight <= 0 {
	//	self.CountHorseFight()
	//}
}

func (self *ModHorse) OnSave(sql bool) {
	self.Encode()
	self.Sql_Horse.Update(sql)
}

//!刷新
func (self *ModHorse) OnRefresh() {
	self.Sql_Horse.SummonNormal = 0
	self.Sql_Horse.SummonSenior = 0
	self.freshHorseTask()
}

func (self *ModHorse) Decode() { //! 将数据库数据写入data
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	json.Unmarshal([]byte(self.Sql_Horse.Info), &self.Sql_Horse.sinfo[0])
	json.Unmarshal([]byte(self.Sql_Horse.Info2), &self.Sql_Horse.sinfo[1])
	json.Unmarshal([]byte(self.Sql_Horse.Info3), &self.Sql_Horse.sinfo[2])
	json.Unmarshal([]byte(self.Sql_Horse.Info4), &self.Sql_Horse.sinfo[3])

	for i := 0; i < INT_MAX_HORSE_PACKAGE; i++ {
		if self.Sql_Horse.sinfo[i] == nil {
			self.Sql_Horse.sinfo[i] = make(map[int]*JS_HorseInfo)
		} else {
			if len(self.Sql_Horse.sinfo[i]) > 0 {
				for hid, horse := range self.Sql_Horse.sinfo[i] {
					if horse.Type > 0 {
						self.Sql_Horse.info[hid] = horse
					} else {
						delete(self.Sql_Horse.sinfo[i], hid)
					}
				}
			}
		}
	}

	json.Unmarshal([]byte(self.Sql_Horse.Soul), &self.Sql_Horse.soul)
	json.Unmarshal([]byte(self.Sql_Horse.Summon), &self.Sql_Horse.summon)
	json.Unmarshal([]byte(self.Sql_Horse.HorseTask), &self.Sql_Horse.horseTask)

}

func (self *ModHorse) Encode() { //! 将data数据写入数据库
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	self.Sql_Horse.Info = HF_JtoA(self.Sql_Horse.sinfo[0])
	self.Sql_Horse.Info2 = HF_JtoA(self.Sql_Horse.sinfo[1])
	self.Sql_Horse.Info3 = HF_JtoA(self.Sql_Horse.sinfo[2])
	self.Sql_Horse.Info4 = HF_JtoA(self.Sql_Horse.sinfo[3])

	self.Sql_Horse.Soul = HF_JtoA(self.Sql_Horse.soul)
	self.Sql_Horse.Summon = HF_JtoA(self.Sql_Horse.summon)
	self.Sql_Horse.HorseTask = HF_JtoA(self.Sql_Horse.horseTask)

}

func (self *ModHorse) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "summonhorse": //! 购买水晶
		var msg C2S_SummonHorse
		json.Unmarshal(body, &msg)
		self.SummonHorse(msg.Index, msg.Num)
		return true
	case "discernhorse": // 召唤魔宠
		var msg C2S_IdentifyHorse
		json.Unmarshal(body, &msg)
		self.DiscernHorse(msg.Num)
		return true
	case "gethorsesoullst":
		self.SendHorseSoulInfo()
		return true
	case "gethorselst":
		self.SendHorseInfo()
		return true
	case "updatehorselst":
		self.UpdateHorse("updatehorselst")
		return true
	case "decomposehorse":
		var msg C2S_DecomposeHorse
		json.Unmarshal(body, &msg)
		self.DecomposeHorse(msg.Horselst)
		return true
	case "decomposesoul":
		var msg C2S_DecomposeSoul
		json.Unmarshal(body, &msg)
		self.DecomposeSoul(msg.Soullst)
		return true
	case "embedhorsesoul":
		var msg C2S_EmbedHorseSoul
		json.Unmarshal(body, &msg)
		self.EmbedHorseSoul(msg.HorseId, msg.SoulId, msg.Index)
		return true
	case "removehorsesoul":
		var msg C2S_EmbedHorseSoul
		json.Unmarshal(body, &msg)
		self.RemoveHorseSoul(msg.HorseId, msg.Index)
		return true
	case "uphorsesoul":
		var msg C2S_UpHorseSoul
		json.Unmarshal(body, &msg)
		self.UpHorseSoul(msg.HorseId, msg.SoulId, msg.Index)
		return true
	case "mounthorse":
		var msg C2S_MountHorse
		json.Unmarshal(body, &msg)
		self.MountHorse2(msg.Heroid, msg.Horseid, msg.Index, msg.TeamType)
		return true
	case "unmounthorse":
		var msg C2S_MountHorse
		json.Unmarshal(body, &msg)
		self.UnmountHorse(msg.Heroid, msg.TeamType)
		return true
	case "combinehorse": // 碎片合成魔宠 高级召唤
		var msg C2S_MountHorse
		json.Unmarshal(body, &msg)
		self.CombineHorse()
		return true
	case "uphorse":
		var msg C2S_UpHorse
		json.Unmarshal(body, &msg)
		self.UpHorse(msg.Star, msg.Material)
		return true
	case "awakehorse":
		var msg C2S_AwakeHorse
		json.Unmarshal(body, &msg)
		self.AwakeHorse(msg.HorseId, msg.Materail)
		return true
	case "exchangesoul":
		var msg C2S_ExchangeSoul
		json.Unmarshal(body, &msg)
		self.ExchangeSoul(msg.Index, msg.Soulid)
		return true
	case "awardhorsetask":
		self.awardHorseTask(ctrl, body)
		return true
	case "horseswitch":
		var msg C2S_SwitchHorse
		json.Unmarshal(body, &msg)
		self.switchHorse(ctrl, msg.KeyId)
		return true
	case "horsewash":
		var msg C2S_SwitchWash
		json.Unmarshal(body, &msg)
		self.washHorse(ctrl, msg.KeyId)
		return true
	case "saveswitch":
		var msg C2S_SaveSwitch
		json.Unmarshal(body, &msg)
		self.saveSwitch(ctrl, msg.KeyId)
		return true
	}

	return false
}

//! 分解魔宠
func (self *ModHorse) DecomposeHorse(horselst []int) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	soulitem := make([]JS_HorseSoulInfo, 0)
	outitem := make([]PassItem, 0)
	for i := 0; i < len(horselst); i++ {
		horsedata, ok := self.Sql_Horse.info[horselst[i]]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THERES_NO_SUCH_THING_AS"))
			return
		}

		//! 分解获得马魂
		if horsedata.Heroid > 0 {
			herodata := self.player.GetModule("hero").(*ModHero).GetHero(horsedata.Heroid)
			if herodata.Horse > 0 {
				herodata.Horse = 0
			}
		}

		for j := 0; j < len(horsedata.SoulLst); j++ {
			if horsedata.SoulLst[j] > 0 {
				soulid := horsedata.SoulLst[j] / 100
				soulrank := horsedata.SoulLst[j] % 100

				soulitem = append(soulitem, JS_HorseSoulInfo{soulid, soulrank, 1})
			}
		}

		csvhorse, ok := GetCsvMgr().Data["Horse_BattleSteed"][horsedata.Type]
		if !ok {
			continue
		}

		//! 获取基础的马魂
		baggroup := HF_DropForItemBagGroup(HF_Atoi(csvhorse["decompose_itembaggroup1"]))
		for j := 0; j < len(baggroup); j++ {
			find := false
			for k := 0; k < len(soulitem); k++ {
				if soulitem[k].Id == baggroup[j].ItemID {
					soulitem[k].Num += baggroup[j].Num
					find = true
					break
				}
			}
			if find == false {
				soulitem = append(soulitem, JS_HorseSoulInfo{baggroup[j].ItemID, 1, baggroup[j].Num})
			}
		}

		baggroup1 := HF_DropForItemBagGroup(HF_Atoi(csvhorse["decompose_itembaggroup2"]))
		for j := 0; j < len(baggroup1); j++ {
			find := false
			for k := 0; k < len(outitem); k++ {
				if outitem[k].ItemID == baggroup1[j].ItemID {
					outitem[k].Num += baggroup1[j].Num
					find = true
					break
				}
			}
			if find == false {
				outitem = append(outitem, PassItem{baggroup1[j].ItemID, baggroup1[j].Num})
			}
		}
		//消耗的魔宠
		GetServer().SqlLog(self.player.GetUid(), horsedata.Type, -1, 0, 0, "分解", 0, 0, self.player)
	}

	//! 增加道具
	for j := 0; j < len(soulitem); j++ {
		soulid := soulitem[j].Id*100 + soulitem[j].Rank
		usersoulitem, ok := self.Sql_Horse.soul[soulid]
		if ok {
			usersoulitem.Num += soulitem[j].Num
		} else {
			newitem := new(JS_HorseSoulInfo)
			newitem.Id = soulitem[j].Id
			newitem.Rank = 1
			newitem.Num = soulitem[j].Num
			self.Sql_Horse.soul[soulid] = newitem
		}

		//获得的马魂
		_, ok = GetCsvMgr().HorseSoulConfig[soulitem[j].Id]
		if ok {
			GetServer().SqlLog(self.player.GetUid(), soulitem[j].Id, soulitem[j].Num, 0, 0, "分解", 0, 0, self.player)
		}

	}
	//获得的道具
	for j := 0; j < len(outitem); j++ {
		self.player.AddObject(outitem[j].ItemID, outitem[j].Num, 0, 0, 0, "分解")
	}

	for i := 0; i < len(horselst); i++ {
		//! 删除魔宠
		//delete(self.Sql_Horse.info, horselst[i])
		self.DelHorse(horselst[i])
	}

	var msg S2C_DecomposeHorse
	msg.Cid = "decomposehorse"
	msg.Horselst = horselst
	msg.Soullst = soulitem
	msg.Item = outitem
	//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_Decomposehorse, 0, 0, 0, "分解魔宠", 0, 0, self.player)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_DECOMPOUND, 1, 0, 0, "分解", 0, 0, self.player)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HORSE_DECOMPOSE, 1, 0, 0, "魔宠分解", 0, 0, self.player)
	self.player.SendMsg("decomposehorse", HF_JtoB(&msg))

	self.player.HandleTask(HaveRunesCount, 0, 0, 0)
	self.player.HandleTask(HaveHorseCount, 0, 0, 0)

}

func (self *ModHorse) getDecomcsv(mode, level int) *HorseSoulUpgradeConfig {
	for _, v := range GetCsvMgr().HorseSoulUpgradeConfig {
		if v.Upgrademodelgroup == mode && v.Level == level {
			return v
		}
	}
	return nil
}

//! 分解马魂
func (self *ModHorse) DecomposeSoul(soullst []JS_HorseSoulInfo) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	outitem := make([]PassItem, 0)
	for i := 0; i < len(soullst); i++ {
		if soulinfo, ok := self.Sql_Horse.soul[soullst[i].Id*100+soullst[i].Rank]; !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_SOUL_ERR"))
			return
		} else {
			if soulinfo.Num < soullst[i].Num {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_SOUL_ERR"))
				return
			}
		}

		//LogDebug("soullst[i].Id:", soullst[i].Id)
		soulcsv, ok := GetCsvMgr().HorseSoulConfig[soullst[i].Id]
		if !ok {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_CONFIGURATION_DOES_NOT_EXIST"))
			return
		}

		upgrademode := soulcsv.Upgrademodel
		decomcsv := self.getDecomcsv(upgrademode, soullst[i].Rank)
		if decomcsv == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_UPGRADE_CONFIGURATION_DOES_NOT_EXIST"))
			return
		}

		find := false
		for m := 0; m < len(outitem); m++ {
			if outitem[m].ItemID == decomcsv.Decomposeitem {
				outitem[m].Num += soullst[i].Num * decomcsv.Decomposenum
				find = true
				break
			}
		}

		if find == false {
			outitem = append(outitem, PassItem{decomcsv.Decomposeitem, soullst[i].Num * decomcsv.Decomposenum})
		}
	}

	for i := 0; i < len(soullst); i++ {
		//! 删除马魂
		self.AddHorseSoul(soullst[i].Id, soullst[i].Rank, -soullst[i].Num)
		GetServer().SqlLog(self.player.GetUid(), soullst[i].Id, -soullst[i].Num, 0, 0, "分解", 0, 0, self.player)
		//delete(self.Sql_Horse.soul, soullst[i].Id*100+soullst[i].Rank)
	}

	for i := 0; i < len(outitem); i++ {
		self.player.AddObject(outitem[i].ItemID, outitem[i].Num, 0, 0, 0, "分解")
	}

	var msg S2C_DecomposeSoul
	msg.Cid = "decomposesoul"
	msg.CostSoul = soullst
	msg.Item = outitem
	//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_Decomposesoul, 0, 0, 0, "分解马魂", 0, 0, self.player)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_DECOMPOUND, 1, 0, 0, "分解", 0, 0, self.player)
	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_HORSE_SOUL_DECOMPOSE, 1, 0, 0, "魔宠符文分解", 0, 0, self.player)
	self.player.SendMsg("decomposesoul", HF_JtoB(&msg))

	self.player.HandleTask(HaveRunesCount, 0, 0, 0)

}

//! 镶嵌马魂
func (self *ModHorse) EmbedHorseSoul(horseid int, soulid int, index int) {
	self.DataLocker.Lock()

	if index <= 0 || index > 6 {
		self.DataLocker.Unlock()
		return
	}

	horsedata, ok := self.Sql_Horse.info[horseid]
	if !ok {
		self.DataLocker.Unlock()
		return
	}

	souldata, ok1 := self.Sql_Horse.soul[soulid]
	if !ok1 || souldata.Num == 0 {
		self.DataLocker.Unlock()
		return
	}

	if index > len(horsedata.SoulLst) {
		self.DataLocker.Unlock()
		return
	}

	outsoul := make([]JS_HorseSoulInfo, 0)

	if horsedata.SoulLst[index-1] != 0 {
		outsoul = append(outsoul, JS_HorseSoulInfo{horsedata.SoulLst[index-1] / 100, horsedata.SoulLst[index-1] % 100, 1})
		self.AddHorseSoul(horsedata.SoulLst[index-1]/100, horsedata.SoulLst[index-1]%100, 1)
	}

	souldata.Num -= 1
	if souldata.Num == 0 {
		delete(self.Sql_Horse.soul, soulid)
	}
	horsedata.SoulLst[index-1] = soulid
	self.CountFight(horsedata)
	var msg S2C_EmbedHorseSoul
	msg.Cid = "embedhorsesoul"
	msg.Horseid = horseid
	msg.Soulid = soulid
	msg.Soullst = outsoul
	msg.Horse = *horsedata

	self.player.SendMsg("embedhorsesoul", HF_JtoB(&msg))
	self.DataLocker.Unlock()
	self.Caculate(horsedata, ReasonEmbedHorseSoul)
}

func (self *ModHorse) Caculate(horseData *JS_HorseInfo, reason int) {
	if horseData == nil {
		return
	}
	hero := self.player.getHero(horseData.Heroid)
	if hero == nil {
		return
	}
	self.player.countHeroFight(hero, reason)

	//_, IsRefresh := self.CountHorseFight()

	//if IsRefresh {
	//GetTopHorseFightMgr().UpdateRank(count, self.player)
	//}
}

func (self *ModHorse) Caculate2(heroKeyId, reason int) {
	hero := self.player.getHero(heroKeyId)
	if hero == nil {
		return
	}
	self.player.countHeroFight(hero, reason)
}

//! 卸载马魂
func (self *ModHorse) RemoveHorseSoul(horseid int, index int) {
	self.DataLocker.Lock()

	horsedata, ok := self.Sql_Horse.info[horseid]
	if !ok {
		self.DataLocker.Unlock()
		return
	}

	if index > len(horsedata.SoulLst) {
		self.DataLocker.Unlock()
		return
	}

	soulid := horsedata.SoulLst[index-1]
	if soulid == 0 {
		self.DataLocker.Unlock()
		return
	}

	self.AddHorseSoul(soulid/100, soulid%100, 1)
	horsedata.SoulLst[index-1] = 0
	self.CountFight(horsedata)

	var msg S2C_EmbedHorseSoul
	msg.Cid = "removehorsesoul"
	msg.Horseid = horseid
	msg.Soulid = soulid
	msg.Horse = *horsedata

	self.player.SendMsg("removehorsesoul", HF_JtoB(&msg))
	self.DataLocker.Unlock()
	self.Caculate(horsedata, ReasonRemoveHorseSoul)

}

//! 升级马魂
func (self *ModHorse) UpHorseSoul(horseid int, upsoul int, index int) {
	horsedata, ok := self.Sql_Horse.info[horseid]
	if !ok {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_UNION_FIGHT_DATA_DOES_NOT_EXIST"))
		return
	}

	if index > len(horsedata.SoulLst) {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HORSE_INDEX_ERROR"))
		return
	}

	soulid := horsedata.SoulLst[index-1]
	if soulid == 0 || soulid != upsoul {
		return
	}

	soultype, rank := soulid/100, soulid%100
	soulcsv, ok := GetCsvMgr().HorseSoulConfig[soultype]
	if !ok {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_HORSE_THE_SOUL_TYPE_DOES_NOT"))
		return
	}

	upsoulid := soulcsv.Upgrademodel
	upcsv := self.getDecomcsv(upsoulid, rank)
	if upcsv == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_TECH_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	costitem, costnum := upcsv.Costitem, upcsv.Costnum
	if self.player.GetObjectNum(costitem) < costnum {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THE_MATERIAL_IS_INSUFFICIENT_AND"))
		return
	}

	self.player.AddObject(costitem, -costnum, upcsv.Id, upcsv.Level+1, 0, "魔宠符文强化")
	itemlst := make([]PassItem, 0)
	itemlst = append(itemlst, PassItem{costitem, -costnum})

	horsedata.SoulLst[index-1] = soultype*100 + rank + 1
	horsedata.chg = true
	self.CountFight(horsedata)

	self.player.HandleTask(CostEnergeTask, costnum, 0, costitem)

	var msg S2C_UpHorseSoul
	msg.Cid = "uphorsesoul"
	msg.Horseid = horseid
	msg.Index = index
	msg.Item = itemlst
	msg.Soulid = horsedata.SoulLst[index-1]
	msg.Horse = *horsedata
	//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_Uphorsesoul, 0, 0, 0, "马魂强化", 0, 0, self.player)
	self.player.SendMsg("uphorsesoul", HF_JtoB(&msg))

	self.Caculate(horsedata, ReasonUpHorseSoul)

	self.player.HandleTask(HaveRunesCount, 0, 0, 0)

	GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_SOUL_UP, soulid, rank, rank+1, "魔宠符文强化", 0, 0, self.player)

}

//! 召唤野马 1-普通召唤 2-高级召唤（元宝）3-高级召唤（道具）
func (self *ModHorse) SummonHorse(index int, num int) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	if num < 0 {
		return
	}

	//! 次数判断
	csvlevel, ok := GetCsvMgr().Data["Horse_Judge_level"][self.Sql_Horse.Level]
	if !ok {
		return
	}

	totalNum := GetCsvMgr().GetVipHorseCall(self.player.Sql_UserBase.Vip)
	if index == 1 {
		if self.Sql_Horse.SummonNormal+num > HF_Atoi(csvlevel["daily_normal_call"]) {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_NOT_ENOUGH_SUMMONS"))
			//return
		}
	} else if index == 2 {
		if self.Sql_Horse.SummonSenior+num > HF_Atoi(csvlevel["daily_higt_call"])+totalNum {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_NOT_ENOUGH_SUMMONS"))
			//return
		}
	}

	costitem := make([]PassItem, 0)
	getitem := make([]int, 0)
	lstcsv := GetCsvMgr().Data2["Horse_Judge_call"]
	addexp := 0
	addDec := "魔宠高级购买"
	if index == 1 {
		addDec = "魔宠普通购买"
		for j := 0; j < num; j++ {
			for k := 0; k < len(lstcsv); k++ {
				callnum := HF_Atoi(lstcsv[k]["call_num"])
				calltype := HF_Atoi(lstcsv[k]["type"])

				if calltype == index && self.Sql_Horse.SummonNormal+j+1 <= callnum {
					if len(costitem) > 0 {
						if costitem[0].ItemID == HF_Atoi(lstcsv[k]["cost_item1"]) {
							costitem[0].Num -= HF_Atoi(lstcsv[k]["cost_num1"])
						}
					} else {
						costitem = append(costitem, PassItem{HF_Atoi(lstcsv[k]["cost_item1"]), -HF_Atoi(lstcsv[k]["cost_num1"])})
					}

					totalweight := HF_Atoi(lstcsv[k]["weight1"]) + HF_Atoi(lstcsv[k]["weight2"])
					if HF_GetRandom(totalweight) < HF_Atoi(lstcsv[k]["weight1"]) {
						getitem = append(getitem, HF_Atoi(lstcsv[k]["dropsmall_horse_id1"]))
					}

					addexp += HF_Atoi(lstcsv[k]["get_judge_exp"])

					break
				}
			}
		}
	} else {
		csv, ok := GetCsvMgr().Data["Horse_Judge_call"][2001]
		if !ok {
			return
		}

		if index == 2 {
			for j := 0; j < len(lstcsv); j++ {
				callnum := HF_Atoi(lstcsv[j]["call_num"])
				calltype := HF_Atoi(lstcsv[j]["type"])
				if calltype == 2 && self.Sql_Horse.SummonSenior+num <= callnum {
					itemid1, itemnum1 := HF_Atoi(lstcsv[j]["cost_item1"]), HF_Atoi(lstcsv[j]["cost_num1"])
					if self.player.GetObjectNum(itemid1) >= itemnum1*num {
						costitem = append(costitem, PassItem{itemid1, 0 - itemnum1*num})
						getitem = append(getitem, HF_Atoi(lstcsv[j]["dropsmall_horse_id1"]))
					} else {
						self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_NOT_ENOUGH_DIAMONDS_TO_CALL"))
						return
					}

					addexp += HF_Atoi(lstcsv[j]["get_judge_exp"])

					break
				}
			}
		} else if index == 3 {
			//! 优先消耗道具2
			itemid2, itemnum2 := HF_Atoi(csv["cost_item2"]), HF_Atoi(csv["cost_num2"])
			if self.player.GetObjectNum(itemid2) >= itemnum2*num {
				costitem = append(costitem, PassItem{itemid2, 0 - itemnum2*num})
				getitem = append(getitem, HF_Atoi(csv["dropsmall_horse_id1"]))
			} else {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THE_PROPS_ARE_INSUFFICIENT_TO"))
				return
			}

			addexp += HF_Atoi(csv["get_judge_exp"])
		}
	}

	for i := 0; i < len(costitem); i++ {
		if self.player.GetObjectNum(costitem[i].ItemID) < 0-costitem[i].Num {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_SHORTAGE_OF_GOLD_COINS"))
			return
		}
	}

	for i := 0; i < len(getitem); i++ {
		if getitem[i] == 1001 {
			self.Sql_Horse.summon[0]++
		} else {
			self.Sql_Horse.summon[1]++
		}
	}

	if index == 1 {
		self.Sql_Horse.SummonNormal += num
	} else if index == 2 {
		self.Sql_Horse.SummonSenior += num
	}
	self.Sql_Horse.Exp += addexp
	if HF_Atoi(csvlevel["exp"]) > 0 && self.Sql_Horse.Exp >= HF_Atoi(csvlevel["exp"]) {
		self.Sql_Horse.Exp -= HF_Atoi(csvlevel["exp"])
		self.Sql_Horse.Level += 1
	}
	for i := 0; i < len(costitem); i++ {
		self.player.AddObject(costitem[i].ItemID, costitem[i].Num, num, self.Sql_Horse.Level, 0, addDec)
	}

	if index == 1 {
		//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_SummonHorse_Normal, num, 0, 0, addDec, self.Sql_Horse.Level, 0, self.player)
		GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_BUY_NOMAL, num, 0, self.Sql_Horse.Level, "魔宠普通购买", 0, 0, self.player)
	} else {
		//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_SummonHorse_Senior, num, 0, 0, addDec, self.Sql_Horse.Level, 0, self.player)
		GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_BUY_SPECIAL, num, 0, self.Sql_Horse.Level, "魔宠高级购买", 0, 0, self.player)
	}

	if index == 1 {
		self.player.HandleTask(SummonHorseTask, num, 1, 0)
	} else if index == 2 || index == 3 {
		self.player.HandleTask(SummonHorseTask, num, 2, 0)
		self.doHorseTask(1)
	}

	var msg S2C_SummonHorse
	msg.Cid = "summonhorse"
	msg.Level = self.Sql_Horse.Level
	msg.Exp = self.Sql_Horse.Exp
	msg.Gold = self.player.Sql_UserBase.Gold
	msg.Gem = self.player.Sql_UserBase.Gem
	msg.SummonNormal = self.Sql_Horse.SummonNormal
	msg.SummonSenior = self.Sql_Horse.SummonSenior
	msg.Summonlst = getitem
	msg.Summontime = self.Sql_Horse.SummonTime
	msg.Item = costitem
	for _, pTask := range self.Sql_Horse.horseTask {
		msg.HorseTask = append(msg.HorseTask, pTask)
	}

	self.player.SendMsg("summonhorse", HF_JtoB(&msg))
}

func (self *ModHorse) GetHorse(horseid int) *JS_HorseInfo {
	if horseinfo, ok := self.Sql_Horse.info[horseid]; ok == true {
		return horseinfo
	}

	return nil
}

func (self *ModHorse) GetHorseSafe(horseid int) *JS_HorseInfo {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	if horseinfo, ok := self.Sql_Horse.info[horseid]; ok == true {
		return horseinfo
	}

	return nil
}

// 获得符合需求的魔宠数量 n2 品质 n3 星级 n4 觉醒等级  0为不限制
func (self *ModHorse) GetHorseCount(n2, n3, n4 int) int {
	nCount := 0 // 符合需求的总数量
	for _, info := range self.Sql_Horse.info { // 循环玩家魔宠
		nType := info.Type
		config, ok := GetCsvMgr().Data["Horse_BattleSteed"][nType]
		if !ok {
			self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_HORSE_NO_CONFIG"), nType))
			continue
		}

		// 品质 星级 觉醒等级
		var bStep, bStar, bLevel = false, false, false
		if n2 == 0 {
			bStep = true
		} else if n2 <= HF_Atoi(config["quality"]) {
			bStep = true
		}

		if n3 == 0 {
			bStar = true
		} else if n3 <= HF_Atoi(config["star"]) {
			bStar = true
		}

		if n4 == 0 {
			bLevel = true
		} else if n4 <= info.Awaken {
			bLevel = true
		}

		if bStep && bStar && bLevel {
			nCount++
		}
	}
	//
	//teamPos := self.player.getTeamPos()
	//
	//for _, info := range teamPos.TeamAttr{
	//	nID := info.HorseKeyId
	//	if nID == 0{
	//		continue
	//	}
	//
	//	horse, ok := self.Sql_Horse.info[nID]
	//	if !ok || horse == nil {
	//		continue
	//	}
	//
	//	nType := horse.Type
	//	config, ok := GetCsvMgr().Data["Horse_BattleSteed"][nType]
	//
	//	var bStep, bStar, bLevel  = false, false, false
	//	if n2 == 0{
	//		bStep = true
	//	}else if n2 <= HF_Atoi(config["quality"]){
	//		bStep = true
	//	}
	//
	//	if n3 == 0{
	//		bStar = true
	//	}else if n3 <= HF_Atoi(config["star"]) {
	//		bStar = true
	//	}
	//
	//	if n4 == 0{
	//		bLevel = true
	//	}else if n4 <= horse.Awaken{
	//		bLevel = true
	//	}
	//
	//	if bStep && bStar && bLevel{
	//		nCount++
	//	}
	//}

	return nCount
}

// 获得符合需求的符文数量 n2 符文等级 n3 符文品质  0为不限制
func (self *ModHorse) GetRunesCount(n2, n3 int) int {
	nCount := 0 // 符合需求的总数量
	for _, info := range self.Sql_Horse.soul { // 循环玩家符文表
		nRunesID := info.Id
		nRunesRank := info.Rank

		config, ok := GetCsvMgr().HorseSoulConfig[nRunesID]
		if !ok {
			continue
		}

		// 符文等级 符文品质
		var bLevel, bStep = false, false
		if n2 == 0 {
			bLevel = true
		} else if n2 <= nRunesRank {
			bLevel = true
		}

		if n3 == 0 {
			bStep = true
		} else if n3 <= config.Quality {
			bStep = true
		}

		if bStep && bLevel {
			nCount += info.Num
		}
	}

	for _, info := range self.Sql_Horse.info { // 循环玩家魔宠
		if info.Heroid > 0 {
			for _, nSoulID := range info.SoulLst {
				if nSoulID > 0 {

					nID := nSoulID / 100
					nRunesRank := nSoulID % 100

					config, ok := GetCsvMgr().HorseSoulConfig[nID]
					if !ok {
						continue
					}

					// 符文等级 符文品质
					var bLevel, bStep = false, false
					if n2 == 0 {
						bLevel = true
					} else if n2 <= nRunesRank {
						bLevel = true
					}

					if n3 == 0 {
						bStep = true
					} else if n3 <= config.Quality {
						bStep = true
					}

					if bStep && bLevel {
						nCount++
					}
				}
			}
		}
	}

	return nCount
}

//! 相马,通过召唤获得召唤次数,通过召唤次数获得抽马次数
func (self *ModHorse) DiscernHorse(num int) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	if num <= 0 {
		return
	}

	summoncount := self.Sql_Horse.summon[0] + self.Sql_Horse.summon[1]
	if summoncount < num {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_YOU_DONT_HAVE_ENOUGH_WILD"))
		return
	}

	if self.GetHorseNum() >= INT_MAXHORSE_NUM {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_YOU_CANT_HAVE_MORE_HORSES"))
		return
	}

	discernnumex := GetCsvMgr().GetHorseParamInt("Horse_Num_1")
	discernitemex := GetCsvMgr().GetHorseParamInt("Horse_Num_2")
	costitem := make([]PassItem, 0)
	getitem := make([]PassItem, 0)
	gethorse := make([]int, 0)
	addexp := 0
	for i := 0; i < num; i++ {
		horseid := 1002
		if self.Sql_Horse.summon[1] < i+1 {
			horseid = 1001
		}
		//horseid := self.Sql_Horse.summon[summoncount-i-1]
		if horseid == 0 {
			continue
		}
		csv, ok := GetCsvMgr().Data["Horse_Judge_Discern"][horseid]
		if !ok {
			continue
		}

		if len(costitem) > 0 {
			costitem[0].Num += HF_Atoi(csv["cost_num"])
		} else {
			costitem = append(costitem, PassItem{HF_Atoi(csv["cost_item"]), HF_Atoi(csv["cost_num"])})
		}

		dropitem := HF_DropForItemBagGroup(HF_Atoi(csv["itembaggroup"]))
		if len(dropitem) > 0 {
			for j := 0; j < len(dropitem); j++ {
				gethorse = append(gethorse, dropitem[j].ItemID)
			}
		}

		if HF_GetRandom(10000) < HF_Atoi(csv["odds"]) {
			getitem = append(getitem, PassItem{HF_Atoi(csv["dropitem_id"]), HF_Atoi(csv["dropitem_num"])})
		}

		if discernnumex > 0 && (self.Sql_Horse.Discern+i+1)%discernnumex == 0 {
			//!奖励道具
			getitem = append(getitem, PassItem{discernitemex, 1})
		}

		addexp += HF_Atoi(csv["get_judge_exp"])
	}

	if len(costitem) > 0 {
		if self.player.GetObjectNum(costitem[0].ItemID) < costitem[0].Num {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_SHORTAGE_OF_GOLD_COINS"))
			return
		}

		costitem[0].ItemID, costitem[0].Num = self.player.AddObject(costitem[0].ItemID, 0-costitem[0].Num, 0, 0, 0, "魔宠召唤")
	}

	for h := 0; h < len(getitem); h++ {
		getitem[h].ItemID, getitem[h].Num = self.player.AddObject(getitem[h].ItemID, getitem[h].Num, self.Sql_Horse.Level, 0, 0, "魔宠召唤")
	}

	csvlevel, ok := GetCsvMgr().Data["Horse_Judge_level"][self.Sql_Horse.Level]
	if ok {
		if self.player.Sql_UserBase.Level >= HF_Atoi(csvlevel["need_lv"]) {
			self.Sql_Horse.Exp += addexp
			if HF_Atoi(csvlevel["exp"]) > 0 && self.Sql_Horse.Exp >= HF_Atoi(csvlevel["exp"]) {
				self.Sql_Horse.Exp -= HF_Atoi(csvlevel["exp"])
				self.Sql_Horse.Level += 1
			}
		}
	}

	self.Sql_Horse.Discern += num
	//!删除现有的
	if num >= self.Sql_Horse.summon[1] {
		self.Sql_Horse.summon[0] -= (num - self.Sql_Horse.summon[1])
		self.Sql_Horse.summon[1] = 0
	} else {
		self.Sql_Horse.summon[1] -= num
	}

	csvlevel1, ok1 := GetCsvMgr().Data["Horse_Judge_level"][self.Sql_Horse.Level]
	if ok1 == true {
		if self.Sql_Horse.summon[0]+self.Sql_Horse.summon[1] < HF_Atoi(csvlevel1["horse_upper_limit"]) &&
			self.Sql_Horse.summon[0]+self.Sql_Horse.summon[1]+num >= HF_Atoi(csvlevel["horse_upper_limit"]) {
			self.Sql_Horse.SummonTime = time.Now().Unix()
		}
	}

	//self.Sql_Horse.summon = self.Sql_Horse.summon[0 : summoncount-num]

	var msg S2C_DiscernHorse
	msg.Cid = "discernhorse"
	msg.Level = self.Sql_Horse.Level
	msg.Exp = self.Sql_Horse.Exp
	msg.Gold = self.player.Sql_UserBase.Gold
	msg.Gem = self.player.Sql_UserBase.Gem
	msg.Discern = self.Sql_Horse.Discern % discernnumex
	msg.Summonlst = self.Sql_Horse.summon
	msg.SummonTime = self.Sql_Horse.SummonTime
	msg.Item = getitem
	msg.Horselst = make([]*JS_HorseInfo, 0)

	for i := 0; i < num; i++ {
		itemID := gethorse[num-i-1]
		csv_horse, ok := GetCsvMgr().Data["Horse_BattleSteed"][itemID]
		if ok {
			horseinfo := self.AddHorse(itemID, false, "魔宠召唤", HF_Atoi(csv_horse["quality"]), HF_Atoi(csv_horse["star"]), 0)
			if horseinfo != nil {
				msg.Horselst = append(msg.Horselst, horseinfo)
			}
		}
	}
	//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_Discernhorse, num, 0, 0, "相马", 0, 0, self.player)
	self.player.SendMsg("discernhorse", HF_JtoB(&msg))

	self.player.HandleTask(DiscernHorse, num, SUMMON_COMMON, 0)
	self.player.HandleTask(HaveHorseCount, 0, 0, 0)
	self.player.GetModule("task").(*ModTask).SendUpdate()

	GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_DISCERN, num, 0, 0, "魔宠召唤", 0, 0, self.player)
}

//! 获得魔宠属性
func (self *ModHorse) GetHorseAttr(group int, num int) []JS_HorseAttr {
	attrlst := make([]JS_HorseAttr, 0)
	attrbasic, ok := GetCsvMgr().HorseAttr_CSV[group]
	if !ok {
		return attrlst
	}

	if len(attrbasic) <= num {
		for i := 0; i < num; i++ {
			for attIndex := 1; attIndex <= 2; attIndex++ {
				nType := HF_Atoi(attrbasic[i][fmt.Sprintf("attribute_type%d", attIndex)])
				nValue := int64(HF_Atoi(attrbasic[i][fmt.Sprintf("attribute_valve%d", attIndex)]))
				attrlst = append(attrlst,
					JS_HorseAttr{nType, nValue})
			}
		}

		return attrlst
	}

	var node []CsvNode
	HF_DeepCopy(&node, &attrbasic)

	for k := 0; k < num; k++ {
		sum := 0
		for i := 0; i < len(node); i++ {
			sum += HF_Atoi(node[i]["weight"])
		}
		pro := HF_GetRandom(sum)
		cur := 0
		for i := 0; i < len(node); i++ {
			cur += HF_Atoi(node[i]["weight"])
			if pro < cur {
				for attIndex := 1; attIndex <= 2; attIndex++ {
					nType := HF_Atoi(node[i][fmt.Sprintf("attribute_type%d", attIndex)])
					nValue := int64(HF_Atoi(node[i][fmt.Sprintf("attribute_valve%d", attIndex)]))
					attrlst = append(attrlst, JS_HorseAttr{nType, nValue})
				}
				copy(node[i:], node[i+1:])
				node = node[:len(node)-1]
				break
			}
		}
	}

	return attrlst
}

// 魔宠上阵
func (self *ModHorse) MountHorse2(heroKeyId int, horseid int, index int, teamType int) {
	self.DataLocker.Lock()

	horsedata, ok := self.Sql_Horse.info[horseid]
	if horsedata == nil || !ok {
		self.DataLocker.Unlock()
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THERES_NO_SUCH_THING_AS"))
		return
	}

	herodata := self.player.GetModule("hero").(*ModHero).GetHero(heroKeyId)
	if herodata == nil {
		self.DataLocker.Unlock()
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THERE_IS_NO_HERO"))
		return
	}

	if horsedata.Heroid != 0 {
		self.DataLocker.Unlock()
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_FAVOURITE_HAS_BEEN_RIDDEN_AND"))
		return
	}

	for j := 0; j < len(horsedata.SoulLst); j++ {
		if horsedata.SoulLst[j] > 0 {
			soulid := horsedata.SoulLst[j] / 100
			soulrank := horsedata.SoulLst[j] % 100
			self.AddHorseSoul(soulid, soulrank, 1)
			horsedata.SoulLst[j] = 0
		}
	}

	soullst := make([]JS_HorseSoulInfo, 0)
	addSoullst := make([]JS_HorseSoulInfo, 0)
	oldHorseKey := herodata.Horse
	// 魔宠魂碎片
	if oldHorseKey != 0 {
		oldhorse, ok := self.Sql_Horse.info[oldHorseKey]
		if ok && oldhorse != nil {
			oldhorse.Heroid = 0
			oldhorse.chg = true

			oldCount := len(oldhorse.SoulLst)
			newCount := len(horsedata.SoulLst)
			if newCount >= oldCount {
				for j := 0; j < oldCount; j++ {
					horsedata.SoulLst[j] = oldhorse.SoulLst[j]
					soulid := oldhorse.SoulLst[j] / 100
					soulrank := oldhorse.SoulLst[j] % 100

					addSoullst = append(addSoullst, JS_HorseSoulInfo{soulid, soulrank, 1})
				}
			} else {
				for j := 0; j < newCount; j++ {
					horsedata.SoulLst[j] = oldhorse.SoulLst[j]
					soulid := oldhorse.SoulLst[j] / 100
					soulrank := oldhorse.SoulLst[j] % 100

					addSoullst = append(addSoullst, JS_HorseSoulInfo{soulid, soulrank, 1})
				}
				for k := newCount; k < oldCount; k++ {
					if oldhorse.SoulLst[k] > 0 {
						soulid := oldhorse.SoulLst[k] / 100
						soulrank := oldhorse.SoulLst[k] % 100

						soullst = append(soullst, JS_HorseSoulInfo{soulid, soulrank, 1})
						self.AddHorseSoul(soulid, soulrank, 1)
					}
				}
			}

			oldhorse.SoulLst = make([]int, oldCount)
		}
	}

	herodata.Horse = horseid

	horsedata.chg = true
	// this two line code will no use..
	horsedata.Heroid = herodata.HeroKeyId
	herodata.Horse = horseid
	self.CountFight(horsedata)

	self.MountNum++

	var msg S2C_MountHero
	msg.Cid = "mounthorse"
	msg.Heroid = heroKeyId
	msg.Horseid = horseid
	msg.HorseSoul = addSoullst
	msg.Soullst = soullst
	msg.Oldhorse = oldHorseKey
	//msg.TeamAttr = pTeamAttr
	msg.TeamType = teamType
	msg.Index = index
	self.player.SendMsg("mounthorse", HF_JtoB(&msg))

	self.DataLocker.Unlock()
	self.Caculate(horsedata, ReasonMountHorse)
}

// 魔宠下阵
func (self *ModHorse) UnmountHorse(heroKeyId int, teamType int) {
	self.DataLocker.Lock()

	herodata := self.player.GetModule("hero").(*ModHero).GetHero(heroKeyId)
	if herodata == nil {
		self.DataLocker.Unlock()
		return
	}

	horseId := herodata.Horse

	horsedata, ok := self.Sql_Horse.info[horseId]
	if horsedata == nil || !ok {
		self.DataLocker.Unlock()
		return
	}

	horsedata.chg = true
	horsedata.Heroid = 0
	herodata.Horse = 0

	//! 卸下马魂
	soullst := make([]JS_HorseSoulInfo, 0)
	for j := 0; j < len(horsedata.SoulLst); j++ {
		if horsedata.SoulLst[j] > 0 {
			soulid := horsedata.SoulLst[j] / 100
			soulrank := horsedata.SoulLst[j] % 100
			soullst = append(soullst, JS_HorseSoulInfo{soulid, soulrank, 1})
			self.AddHorseSoul(soulid, soulrank, 1)
		}
	}
	soulnum := len(horsedata.SoulLst)
	horsedata.SoulLst = make([]int, soulnum)
	self.CountFight(horsedata)
	self.MountNum--

	var msg S2C_MountHero
	msg.Cid = "unmounthorse"
	msg.Heroid = heroKeyId
	msg.Horseid = horseId
	msg.Soullst = soullst
	//msg.TeamAttr = pTeamAttr
	msg.TeamType = teamType
	//msg.Index = index
	self.player.SendMsg("unmounthorse", HF_JtoB(&msg))
	self.DataLocker.Unlock()
	self.Caculate2(heroKeyId, ReasonUnMountHorse)
}

/*
func (self *ModHorse) ChangeMountHorse(index int, teamType int, heroIdNew int, heroIdOld int) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	pTeamAttr, _, err := self.CheckTeam(index, teamType)
	if err != nil {
		self.player.SendErr(err.Error())
		return
	}

	horseId := pTeamAttr.HorseKeyId

	horsedata, ok := self.Sql_Horse.info[horseId]
	if horsedata == nil || !ok {
		return
	}

	herodatanew := self.player.GetModule("hero").(*ModHero).GetHero(heroIdNew)
	herodataold := self.player.GetModule("hero").(*ModHero).GetHero(heroIdOld)
	if herodatanew == nil || herodataold == nil {
		return
	}

	horsedata.chg = true
	horsedata.Heroid = heroIdNew
	herodatanew.Horse = horseId
	herodataold.Horse = 0
}

*/
//! 碎片合成魔宠
func (self *ModHorse) CombineHorse() {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	costitem := GetCsvMgr().GetHorseParamInt("Orange_WarHorse_1")
	costnum := GetCsvMgr().GetHorseParamInt("Orange_WarHorse_2")
	dropgroup := GetCsvMgr().GetHorseParamInt("Orange_WarHorse_3")

	neednum := GetCsvMgr().GetHorseParamInt("Orange_Surprised_1")
	needdrop := GetCsvMgr().GetHorseParamInt("Orange_Surprised_2")

	if self.player.GetObjectNum(costitem) < costnum {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_DIAMOND_SHORTAGE"))
		return
	}

	if self.GetHorseNum() >= INT_MAXHORSE_NUM {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_YOU_CANT_HAVE_MORE_HORSES"))

		return
	}

	costitems := make([]PassItem, 0)
	self.player.AddObject(costitem, 0-costnum, 0, 0, 0, "魔宠合成")
	costitems = append(costitems, PassItem{costitem, 0 - costnum})

	self.Sql_Horse.Combine++

	var msg S2C_CombineHorse
	msg.Cid = "combinehorse"
	msg.Cost = costitems
	msg.Combine = self.Sql_Horse.Combine % neednum

	if 0 == self.Sql_Horse.Combine%neednum {
		dropgroup = needdrop
	}
	dropitem := HF_DropForItemBagGroup(dropgroup)
	if len(dropitem) > 0 {
		for i := 0; i < len(dropitem); i++ {
			horseinfo := self.AddHorse(dropitem[i].ItemID, false, "魔宠合成", 0, 0, 0)
			msg.Horse = *horseinfo
			//GetServer().SqlLog(self.player.GetUid(), dropitem[i].ItemID, 1, 0, 0, "合成魔宠", 0, 0, self.player)
			GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_COMBINE, dropitem[i].ItemID, 0, 0, "魔宠合成", 0, 0, self.player)
		}
	}

	self.player.HandleTask(DiscernHorse, 1, SUMMON_COMPOUND, 0)
	self.player.HandleTask(HaveHorseCount, 0, 0, 0)
	//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_Combinehorse, 0, 0, 0, "合成魔宠", 0, 0, self.player)

	self.player.SendMsg("combinehorse", HF_JtoB(&msg))
}

//! 魔宠升星
func (self *ModHorse) UpHorse(star int, material []int) bool {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	needhorsestar := GetCsvMgr().GetHorseParamInt(fmt.Sprintf("%d_Orange_WarHorse_1", star))
	neednum := GetCsvMgr().GetHorseParamInt(fmt.Sprintf("%d_Orange_WarHorse_2", star))
	quality := GetCsvMgr().GetHorseParamInt(fmt.Sprintf("%d_Orange_WarHorse_6", star))

	if len(material) != neednum {
		return false
	}

	for i := 0; i < len(material); i++ {
		horsedata, ok := self.Sql_Horse.info[material[i]]
		if !ok {
			return false
		}

		horsecsv, ok1 := GetCsvMgr().Data["Horse_BattleSteed"][horsedata.Type]
		if !ok1 {
			return false
		}

		if HF_Atoi(horsecsv["star"]) != needhorsestar {
			return false
		}

		if HF_Atoi(horsecsv["quality"]) < quality {
			return false
		}
	}

	costmaterial := GetCsvMgr().GetHorseParamInt(fmt.Sprintf("%d_Orange_WarHorse_3", star))
	costnum := GetCsvMgr().GetHorseParamInt(fmt.Sprintf("%d_Orange_WarHorse_4", star))
	dropgroup := GetCsvMgr().GetHorseParamInt(fmt.Sprintf("%d_Orange_WarHorse_5", star))

	if self.player.GetObjectNum(costmaterial) < costnum {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TALENT_INSUFFICIENT_MATERIAL"))
		return false
	}

	self.player.AddObject(costmaterial, 0-costnum, 0, 0, 0, "魔宠升星")

	costitem := make([]PassItem, 0)
	soullst := make([]JS_HorseSoulInfo, 0)
	costitem = append(costitem, PassItem{costmaterial, -costnum})

	for i := 0; i < len(material); i++ {
		horsedata, _ := self.Sql_Horse.info[material[i]]
		if horsedata.Heroid != 0 {
			hero := self.player.GetModule("hero").(*ModHero).GetHero(horsedata.Heroid)
			if hero != nil {
				hero.Horse = 0
			}

			horsedata.Heroid = 0
		}

		//! 卸下马魂
		for j := 0; j < len(horsedata.SoulLst); j++ {
			if horsedata.SoulLst[j] > 0 {
				soulid := horsedata.SoulLst[j] / 100
				soulrank := horsedata.SoulLst[j] % 100

				soullst = append(soullst, JS_HorseSoulInfo{soulid, soulrank, 1})
				self.AddHorseSoul(soulid, soulrank, 1)
			}
		}

		//delete(self.Sql_Horse.info, material[i])
		//GetServer().SqlLog(self.player.GetUid(), horsedata.Type, -1, 0, 0, fmt.Sprintf("合成%d星橙马", star), 0, 0, self.player)
		self.DelHorse(material[i])
	}

	var msg S2C_UpHorse
	msg.Cid = "uphorse"
	msg.Cost = costitem
	msg.Soullst = soullst
	msg.Material = material

	dropitem := HF_DropForItemBagGroup(dropgroup)
	if len(dropitem) > 0 {
		for i := 0; i < len(dropitem); i++ {
			msg.Horse = *self.AddHorse(dropitem[i].ItemID, false, "魔宠升星", star, 0, 0)
			//GetServer().SqlLog(self.player.GetUid(), dropitem[i].ItemID, 1, star, 0, "魔宠升星", 0, 0, self.player)
		}
	}

	GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_UP, msg.Horse.Type, star, 0, "魔宠升星", 0, 0, self.player)

	//if star == 4 {
	//	//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_Uphorse_4, 0, 0, 0, fmt.Sprintf("合成%d星橙马", star), 0, 0, self.player)
	//	GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_UP, msg.Horse.Type, star, 0, "魔宠升星", 0, 0, self.player)
	//} else if star == 5 {
	//	//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_Uphorse_5, 0, 0, 0, fmt.Sprintf("合成%d星橙马", star), 0, 0, self.player)
	//	GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_UP, msg.Horse.Type, star, 0, "魔宠升星", 0, 0, self.player)
	//}

	self.player.HandleTask(HaveHorseCount, 0, 0, 0)

	self.player.SendMsg("uphorse", HF_JtoB(&msg))

	return true
}

//! 魔宠觉醒
func (self *ModHorse) AwakeHorse(horseid int, material []int) {
	self.DataLocker.Lock()

	horsedata, ok := self.Sql_Horse.info[horseid]
	if !ok {
		self.DataLocker.Unlock()
		return
	}

	LogDebug("AwakeHorse:", horsedata)

	if horsedata.Awaken >= MAX_HORSE_AWAKEN_LEVEL-1 {
		self.DataLocker.Unlock()
		return
	}

	awakecsv, ok := GetCsvMgr().Data["Horse_BattleSteed_Awaken"][horsedata.Type*100+horsedata.Awaken+1]
	if !ok {
		self.DataLocker.Unlock()
		return
	}

	for i := 0; i < len(material); i++ {
		horsedata, ok := self.Sql_Horse.info[material[i]]
		if !ok {
			self.DataLocker.Unlock()
			return
		}

		horsecsv, ok1 := GetCsvMgr().Data["Horse_BattleSteed"][horsedata.Type]
		if !ok1 {
			self.DataLocker.Unlock()
			return
		}

		if HF_Atoi(horsecsv["star"]) != HF_Atoi(awakecsv["cost_horse_star"]) {
			self.DataLocker.Unlock()
			return
		}
	}

	if self.player.GetObjectNum(HF_Atoi(awakecsv["cost_item"])) < HF_Atoi(awakecsv["cost_num"]) {
		self.DataLocker.Unlock()
		return
	}

	costitem := make([]PassItem, 0)

	self.player.AddObject(HF_Atoi(awakecsv["cost_item"]), 0-HF_Atoi(awakecsv["cost_num"]), horsedata.Awaken+1, 0, 0, "魔宠觉醒")
	costitem = append(costitem, PassItem{HF_Atoi(awakecsv["cost_item"]), 0 - HF_Atoi(awakecsv["cost_num"])})

	soullst := make([]JS_HorseSoulInfo, 0)

	for i := 0; i < len(material); i++ {
		horsedata, _ := self.Sql_Horse.info[material[i]]
		if horsedata.Heroid != 0 {
			hero := self.player.GetModule("hero").(*ModHero).GetHero(horsedata.Heroid)
			if hero != nil {
				hero.Horse = 0
			}

			horsedata.Heroid = 0
		}

		//! 卸下马魂
		for j := 0; j < len(horsedata.SoulLst); j++ {
			if horsedata.SoulLst[j] > 0 {
				soulid := horsedata.SoulLst[j] / 100
				soulrank := horsedata.SoulLst[j] % 100

				soullst = append(soullst, JS_HorseSoulInfo{soulid, soulrank, 1})
				self.AddHorseSoul(soulid, soulrank, 1)
			}
		}
		//消耗的魔宠
		GetServer().SqlLog(self.player.GetUid(), horsedata.Type, -1, horsedata.Awaken+1, 0, "魔宠觉醒", 0, 0, self.player)
		self.DelHorse(material[i])
		//delete(self.Sql_Horse.info, material[i])
		//for i := 0; i < INT_MAX_HORSE_PACKAGE; i++ {
		//	delete(self.Sql_Horse.sinfo[i], material[i])
		//}
	}

	horsedata.Awaken += 1
	self.CountFight(horsedata)
	horsedata.chg = true

	var msg S2C_UpHorse
	msg.Cid = "awakehorse"
	msg.Cost = costitem
	msg.Soullst = soullst
	msg.Material = material
	msg.Horse = *horsedata

	//().SqlLog(self.player.GetUid(), LOG_EVENT_Awakehorse, 0, 0, 0, "魔宠觉醒", 0, 0, self.player)

	GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_AWAKE, msg.Horse.Type, 0, 0, "魔宠分解", 0, 0, self.player)

	self.player.SendMsg("awakehorse", HF_JtoB(&msg))

	self.player.HandleTask(HaveHorseCount, 0, 0, 0)

	self.DataLocker.Unlock()
	self.Caculate(horsedata, ReasonAwakenHorse)
}

func (self *ModHorse) AddHorseSoulSafe(soulid int, rank int, num int) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	self.AddHorseSoul(soulid, rank, num)
}

func (self *ModHorse) AddHorseSoul(soulid int, rank int, num int) {
	_, ok := GetCsvMgr().HorseSoulConfig[soulid]
	if !ok {
		return
	}

	soulinfo, ok := self.Sql_Horse.soul[soulid*100+rank]
	if ok {
		if soulinfo.Num+num == 0 {
			delete(self.Sql_Horse.soul, soulid*100+rank)
		} else {
			soulinfo.Num += num
		}
	} else {
		soulinfo = new(JS_HorseSoulInfo)
		soulinfo.Id = soulid
		soulinfo.Rank = rank
		soulinfo.Num = num
		self.Sql_Horse.soul[soulid*100+rank] = soulinfo
	}
}

func (self *ModHorse) AddHorseSafe(horseid int, chg bool, dec string) *JS_HorseInfo {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	return self.AddHorse(horseid, chg, dec, 0, 0, 0)
}

// 增加魔宠
func (self *ModHorse) AddHorse(horseid int, chg bool, dec string, param1, param2, param3 int) *JS_HorseInfo {
	csv_horse, ok := GetCsvMgr().Data["Horse_BattleSteed"][horseid]
	if !ok {
		return nil
	}

	newhorse := new(JS_HorseInfo)
	newhorse.Id = self.Sql_Horse.MaxHorseId
	newhorse.Type = horseid
	newhorse.AttLst = make([]JS_HorseAttr, 0)
	newhorse.Awaken = 0
	//newhorse.Skill = 100306

	newhorse.SoulLst = make([]int, 0)

	newhorse.AttLst = self.GetHorseAttr(HF_Atoi(csv_horse["attribute_list_group"]), HF_Atoi(csv_horse["extract_attribute_num"]))
	newhorse.SoulLst = make([]int, HF_Atoi(csv_horse["hole_count"]))

	//newhorse.Skill = HF_Atoi(csv_horse["skill_id"]) + HF_Atoi(csv_horse["skill_lv"]) - 1
	newhorse.Skill = HF_Atoi(csv_horse["skill_id"])
	newhorse.chg = chg

	self.Sql_Horse.info[newhorse.Id] = newhorse
	addflag := false
	for i := 0; i < INT_MAX_HORSE_PACKAGE; i++ {
		if len(self.Sql_Horse.sinfo[i]) < (i+1)*100 {
			self.Sql_Horse.sinfo[i][newhorse.Id] = newhorse
			addflag = true
			break
		}
	}

	self.CountFight(newhorse)
	if addflag == false {
		self.Sql_Horse.sinfo[0][newhorse.Id] = newhorse
	}
	//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_GET_HORSE, horseid, newhorse.Id, 0, dec, 0, 0, self.player)

	GetServer().SqlLog(self.player.GetUid(), horseid, 1, param1, param2, dec, 0, param3, self.player)

	//! 获得公告
	if HF_Atoi(csv_horse["quality"]) >= 14 && (dec == "魔宠召唤" || dec == "魔宠合成") { //! 橙马公告
		content := fmt.Sprintf(GetCsvMgr().GetText("STR_HORSE_CONTENT"),
			HF_GetColorByCamp(3), CAMP_NAME[self.player.GetCamp()-1],
			HF_GetColorByCamp(3), self.player.Sql_UserBase.UName, dec,
			HF_GetColorByCamp(3), csv_horse["star"],
			"231#131#33#", csv_horse["name"])

		GetServer().sendSysChat(content)
	}

	self.Sql_Horse.MaxHorseId++

	return newhorse
}

func (self *ModHorse) GetHorseNum() int {
	return len(self.Sql_Horse.info) - self.MountNum
}

//! 删除魔宠
func (self *ModHorse) DelHorse(horseid int) {
	for i := 0; i < INT_MAX_HORSE_PACKAGE; i++ {
		delete(self.Sql_Horse.sinfo[i], horseid)
	}

	delete(self.Sql_Horse.info, horseid)
}

//！获取马魂
func (self *ModHorse) GetHorseSoul(soulid int, rank int) int {
	soulinfo, ok := self.Sql_Horse.soul[soulid*100+rank]
	if ok {
		return soulinfo.Num
	} else {
		return 0
	}
}

//! 兑换物品
func (self *ModHorse) ExchangeSoul(index int, itemid int) {
	self.DataLocker.Lock()
	defer self.DataLocker.Unlock()

	csv, ok := GetCsvMgr().Data["Horse_Shop"][index]
	if !ok {
		return
	}

	if self.player.GetObjectNum(HF_Atoi(csv["cost_item"])) < HF_Atoi(csv["cost_num"]) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TIGER_LACK_OF_PROPS"))
		return
	}

	item := make([]PassItem, 0)
	soul := make([]JS_HorseSoulInfo, 0)

	//! 1-马魂，2-道具
	if HF_Atoi(csv["type"]) == 1 {
		self.AddHorseSoul(HF_Atoi(csv["item_id"]), HF_Atoi(csv["lv"]), HF_Atoi(csv["item_num"]))
		soul = append(soul, JS_HorseSoulInfo{HF_Atoi(csv["item_id"]), HF_Atoi(csv["lv"]), HF_Atoi(csv["item_num"])})

		csv_horse, ok := GetCsvMgr().HorseSoulConfig[HF_Atoi(csv["item_id"])]
		if ok {
			GetServer().SqlLog(self.player.GetUid(), HF_Atoi(csv["item_id"]), HF_Atoi(csv["item_num"]), csv_horse.Quality, 0, "魔宠符文兑换", 0, 1, self.player)
		}

	} else {
		self.player.AddObject(HF_Atoi(csv["item_id"]), HF_Atoi(csv["item_num"]), 0, 0, 0, "魔宠符文兑换")
		item = append(item, PassItem{HF_Atoi(csv["item_id"]), HF_Atoi(csv["item_num"])})
	}

	self.player.AddObject(HF_Atoi(csv["cost_item"]), 0-HF_Atoi(csv["cost_num"]), 0, 0, 0, "魔宠符文兑换")
	item = append(item, PassItem{HF_Atoi(csv["cost_item"]), 0 - HF_Atoi(csv["cost_num"])})

	var msg S2C_ExchangeSoul
	msg.Cid = "exchangesoul"
	msg.Index = index
	msg.Soul = soul
	msg.Item = item
	//GetServer().SqlLog(self.player.GetUid(), LOG_EVENT_Exchangesoul, 0, 0, 0, "马魂兑换", 0, 0, self.player)
	self.player.SendMsg("exchangesoul", HF_JtoB(&msg))

	self.player.HandleTask(HaveRunesCount, 0, 0, 0)

	self.player.HandleTask(CostEnergeTask, HF_Atoi(csv["cost_num"]), 0, HF_Atoi(csv["cost_item"]))

	GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_SOUL_EXCHANGE, HF_Atoi(csv["item_id"]), 0, 0, "魔宠符文兑换", 0, 0, self.player)

}

//! 自动召唤野马
func (self *ModHorse) AutoSummon() int64 {
	tNow := time.Now().Unix()
	if tNow >= self.Sql_Horse.SummonTime {
		if self.Sql_Horse.SummonTime == 0 {
			self.Sql_Horse.SummonTime = tNow
		}

		csv_level, ok := GetCsvMgr().Data["Horse_Judge_level"][self.Sql_Horse.Level]
		if !ok {
			return 0
		}

		if self.Sql_Horse.summon[0]+self.Sql_Horse.summon[1] >= HF_Atoi(csv_level["horse_upper_limit"]) {
			self.Sql_Horse.SummonTime = tNow
			return 0
		}

		passtime := int(tNow - self.Sql_Horse.SummonTime)
		recovernum := passtime / HF_Atoi(csv_level["recovery_horse_time"])

		if recovernum+len(self.Sql_Horse.summon) > HF_Atoi(csv_level["horse_upper_limit"]) {
			recovernum = HF_Atoi(csv_level["horse_upper_limit"]) - (self.Sql_Horse.summon[0] + self.Sql_Horse.summon[1])
		}

		for i := 0; i < recovernum; i++ {
			self.Sql_Horse.SummonTime += int64(HF_Atoi(csv_level["recovery_horse_time"]))
			if len(self.Sql_Horse.summon) < HF_Atoi(csv_level["horse_upper_limit"]) {
				self.Sql_Horse.summon[0]++
			}
		}
	}

	return self.Sql_Horse.SummonTime
}

//! 发送数据
func (self *ModHorse) SendHorseInfo() {
	self.player.GetModule("hero").(*ModHero).CheckHorse()

	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()
	self.checkHorseTask()
	discernnumex := GetCsvMgr().GetHorseParamInt("Horse_Num_1")
	self.AutoSummon()
	var msg S2C_HorseInfo
	msg.Cid = "horselst"
	msg.Level = self.Sql_Horse.Level
	msg.Exp = self.Sql_Horse.Exp
	msg.SummonNormal = self.Sql_Horse.SummonNormal
	msg.SummonSenior = self.Sql_Horse.SummonSenior
	msg.SummonTime = self.Sql_Horse.SummonTime
	msg.Summonlst = self.Sql_Horse.summon
	msg.Combine = self.Sql_Horse.Combine
	msg.Discern = self.Sql_Horse.Discern % discernnumex
	msg.Horselst = make([]*JS_HorseInfo, 0)
	for _, value := range self.Sql_Horse.info {
		msg.Horselst = append(msg.Horselst, value)
	}

	for _, pTask := range self.Sql_Horse.horseTask {
		msg.HorseTask = append(msg.HorseTask, pTask)
	}

	self.player.SendMsg("horselst", HF_JtoB(&msg))
}

//! 更新魔宠信息
func (self *ModHorse) UpdateHorse(cid string) {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	self.AutoSummon()
	discernnumex := GetCsvMgr().GetHorseParamInt("Horse_Num_1")

	var msg S2C_HorseInfo
	msg.Cid = cid
	msg.Level = self.Sql_Horse.Level
	msg.Exp = self.Sql_Horse.Exp
	msg.Combine = self.Sql_Horse.Combine
	msg.SummonNormal = self.Sql_Horse.SummonNormal
	msg.SummonSenior = self.Sql_Horse.SummonSenior
	msg.SummonTime = self.Sql_Horse.SummonTime
	msg.Summonlst = self.Sql_Horse.summon
	msg.Discern = self.Sql_Horse.Discern % discernnumex
	msg.Horselst = make([]*JS_HorseInfo, 0)
	for _, value := range self.Sql_Horse.info {
		if value.chg == true {
			msg.Horselst = append(msg.Horselst, value)
			value.chg = false
		}
	}
	self.player.SendMsg("horselst", HF_JtoB(&msg))
}

//! 发送马魂信息
func (self *ModHorse) SendHorseSoulInfo() {
	self.DataLocker.RLock()
	defer self.DataLocker.RUnlock()

	var msg S2C_HorseSoulInfo
	msg.Cid = "horsesoullst"
	msg.NewSoul = false
	msg.SummonNormal = self.Sql_Horse.SummonNormal
	msg.SummonSenior = self.Sql_Horse.SummonSenior
	msg.Summontime = self.Sql_Horse.SummonTime
	msg.Soullst = make([]*JS_HorseSoulInfo, 0)
	msg.Horselst = make([]*JS_HorseInfo, 0)

	for _, value := range self.Sql_Horse.soul {
		msg.Soullst = append(msg.Soullst, value)
	}

	for _, value := range self.Sql_Horse.info {
		if value.chg == true {
			msg.Horselst = append(msg.Horselst, value)
			value.chg = false
		}
	}
	self.player.SendMsg("horsesoullst", HF_JtoB(&msg))
}

func (self *ModHorse) addAttEx(inparam map[int]*Attribute, attMap map[int]*Attribute) {
	for _, att := range inparam {
		_, ok := attMap[att.AttType]
		if !ok {
			attMap[att.AttType] = &Attribute{
				AttType:  att.AttType,
				AttValue: att.AttValue,
			}
		} else {
			attMap[att.AttType].AttValue += att.AttValue
		}
	}
}

func (self *ModHorse) GetHorseAttrInfo(heroKeyId int) map[int]*Attribute {
	hero := self.player.getHero(heroKeyId)
	return self.getAttr(hero.Horse)
}

// 获取魔宠总属性: 魔宠属性, 觉醒属性, 马魂属性, 魔宠技能属性
func (self *ModHorse) getAttr(horseId int) map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	horsedata := self.GetHorse(horseId)
	if horsedata == nil {
		return attMap
	}

	// 基础属性 + 觉醒 + 马魂 + 马魂技能
	base := self.getBaseAttr(horsedata)
	self.addAttEx(base, attMap)
	// 觉醒
	awaken := self.getAwakenAttr(horsedata)
	self.addAttEx(awaken, attMap)
	// 马魂
	soul := self.getSoulAttr(horsedata)
	self.addAttEx(soul, attMap)
	// 马魂技能
	skill := self.getSkillAttr(horsedata)
	self.addAttEx(skill, attMap)

	return attMap
}

// 魔宠属性
func (self *ModHorse) getBaseAttr(horsedata *JS_HorseInfo) map[int]*Attribute {
	res := make(map[int]*Attribute)
	if horsedata == nil {
		return res
	}

	for k := 0; k < len(horsedata.AttLst); k++ {
		attrType := horsedata.AttLst[k].AttrType
		attrValue := horsedata.AttLst[k].AttrValue
		if attrValue == 0 {
			continue
		}

		v, ok := res[attrType]
		if !ok {
			res[attrType] = &Attribute{AttType: attrType, AttValue: attrValue}
		} else {
			v.AttValue += attrValue
		}
	}
	return res
}

// 魔宠觉醒属性
func (self *ModHorse) getAwakenAttr(horsedata *JS_HorseInfo) map[int]*Attribute {
	res := make(map[int]*Attribute)
	if horsedata == nil {
		return res
	}

	if horsedata.Awaken <= 0 {
		return res
	}

	awake_csv, ok := GetCsvMgr().Data["Horse_BattleSteed_Awaken"][horsedata.Type*100+horsedata.Awaken+1]
	if !ok {
		return res
	}

	for k := 0; k < 6; k++ {
		attrType := HF_Atoi(awake_csv[fmt.Sprintf("attribute_type%d", k+1)])
		attrValue := int64(HF_Atoi(awake_csv[fmt.Sprintf("attribute_valve%d", k+1)]))
		if attrValue == 0 {
			continue
		}

		v, ok := res[attrType]
		if !ok {
			res[attrType] = &Attribute{AttType: attrType, AttValue: attrValue}
		} else {
			v.AttValue += attrValue
		}
	}

	return res
}

// 马魂属性
func (self *ModHorse) getSoulAttr(horsedata *JS_HorseInfo) map[int]*Attribute {
	res := make(map[int]*Attribute)
	if horsedata == nil {
		return res
	}

	for k := 0; k < len(horsedata.SoulLst); k++ {
		if horsedata.SoulLst[k] == 0 {
			continue
		}

		soulid := horsedata.SoulLst[k] / 100
		soullevel := horsedata.SoulLst[k] % 100
		if soulcsv, ok := GetCsvMgr().HorseSoulConfig[soulid]; ok {
			var attrValue = int64(0)
			upgradecsv := self.getDecomcsv(soulcsv.Upgrademodel, soullevel)
			if upgradecsv == nil {
				continue
			}

			// 孔的属性
			attrValue = upgradecsv.Value
			v, ok := res[soulcsv.AttType]
			if !ok {
				res[soulcsv.AttType] = &Attribute{AttType: soulcsv.AttType, AttValue: upgradecsv.Value}
			} else {
				v.AttValue += attrValue
			}

			vOne, okOne := res[soulcsv.AttTypeOne]
			if !okOne {
				res[soulcsv.AttTypeOne] = &Attribute{AttType: soulcsv.AttTypeOne, AttValue: upgradecsv.ValueOne}
			} else {
				vOne.AttValue += upgradecsv.ValueOne
			}

			// 孔的战力
			_, ok1 := res[upgradecsv.AttType1]
			if !ok1 {
				res[upgradecsv.AttType1] = &Attribute{AttType: upgradecsv.AttType1, AttValue: upgradecsv.AttValue1}
			} else {
				res[upgradecsv.AttType1].AttValue += upgradecsv.AttValue1
			}

			// 觉醒附带属性
			if horsedata.Awaken > k {
				res[soulcsv.AttType].AttValue += soulcsv.AttValue
				res[soulcsv.AttTypeOne].AttValue += soulcsv.AttValueOne
				AddFightAtt(res, soulcsv.DragonAtt, soulcsv.DragonValue)
			}
		}
	}

	return res
}

// 马魂技能属性
func (self *ModHorse) getSkillAttr(horsedata *JS_HorseInfo) map[int]*Attribute {
	res := make(map[int]*Attribute)
	if horsedata == nil {
		return res
	}

	skill := horsedata.Skill
	skillConfig, ok := GetCsvMgr().SkillConfig[skill]
	if !ok {
		return res
	}

	if skillConfig.Skilltype != 1 {
		return res
	}

	var skillCounts []int64
	if len(skillConfig.Skillcounts) == len(skillConfig.Skilladdvalues) {
		for i := 0; i < len(skillConfig.Skillcounts); i++ {
			skillCounts = append(skillCounts, int64(skillConfig.Skillcounts[i]+skillConfig.Skilladdvalues[i]))
		}
	}

	AddAttrDirect(res, skillConfig.Skillvaluetypes, skillCounts)

	return res
}

func NewHorseTask(id int) *HorseTask {
	return &HorseTask{Id: id}
}

// 检查魔宠召唤次数
func (self *ModHorse) checkHorseTask() {
	configs := GetCsvMgr().HorseAward
	if self.Sql_Horse.horseTask == nil {
		self.Sql_Horse.horseTask = make(map[int]*HorseTask)
	}

	for _, config := range configs {
		_, ok := self.Sql_Horse.horseTask[config.Id]
		if !ok {
			self.Sql_Horse.horseTask[config.Id] = NewHorseTask(config.Id)
		}
	}
}

// 刷购召唤次数
func (self *ModHorse) freshHorseTask() {
	for _, pTask := range self.Sql_Horse.horseTask {
		if pTask == nil {
			continue
		}
		pTask.Status = 0
		pTask.Process = 0
	}
}

//! 完成高级召唤
func (self *ModHorse) doHorseTask(times int) {
	for _, pTask := range self.Sql_Horse.horseTask {
		if pTask.Status != 0 {
			continue
		}

		config, ok := GetCsvMgr().HorseAwardMap[pTask.Id]
		if !ok {
			LogError("doHorseTask config error, id:", pTask.Id)
			continue
		}
		pTask.Process += times
		if pTask.Process >= config.Times {
			pTask.Process = config.Times
			pTask.Status = 1
		}
	}
}

//! 领取高级召唤奖励
func (self *ModHorse) awardHorseTask(ctrl string, body []byte) {
	var msg C2S_AwardHorseTask
	json.Unmarshal(body, &msg)

	taskId := msg.Id
	pTask, ok := self.Sql_Horse.horseTask[taskId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_MISSION_DOES_NOT_EXIST"))
		return
	}

	if pTask.Status == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_MISSION_NOT_COMPLETED"))
		return
	}

	if pTask.Status == 2 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_TASKS_HAVE_BEEN_RECEIVED"))
		return
	}

	config, ok := GetCsvMgr().HorseAwardMap[taskId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	outItem := self.player.AddObjectLst(config.Rewards, config.Nums, "领取高级召唤任务奖励", 0, 0, 0)

	pTask.Status = 2
	data := &S2C_AwardHorseTask{
		Cid:   ctrl,
		Info:  pTask,
		Items: outItem,
	}
	self.player.SendMsg(data.Cid, HF_JtoB(data))
}

// 魔宠转换
func (self *ModHorse) switchHorse(ctrl string, keyId int) {
	horseData, ok := self.Sql_Horse.info[keyId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_FAVOURITE_DOES_NOT_EXIST"))
		return
	}

	//if horseData.Heroid != 0 {
	//	self.player.SendErrInfo("err", "魔宠必须下阵才能转换")
	//	return
	//}

	// 检查消耗
	itemId := GetCsvMgr().HorseSwitchId
	itemNum := GetCsvMgr().HorseSwitchNum
	if itemId == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	if itemNum == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_SMELT_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	if self.player.GetObjectNum(itemId) < itemNum {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TIGER_LACK_OF_PROPS"))
		return
	}

	// 检查配置是否存在
	horseId := horseData.Type
	_, ok = GetCsvMgr().HorseSwitchMap[horseId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_CURRENT_SPOILERS_CANNOT_BE_CONVERTED"))
		return
	}

	randHorseId, err := GetCsvMgr().randHorseId(horseId)
	if err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}

	horseData.RandHorseId = randHorseId

	//horseData.Type = randHorseId
	items := self.player.RemoveObjectSimple(itemId, itemNum, "魔宠转换", 0, 0, 0)

	var msg S2C_SwitchHero
	msg.Cid = ctrl
	msg.Items = items
	msg.Horse = horseData
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))

	GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_SWITCH, horseId, randHorseId, 0, "魔宠转换", 0, 0, self.player)
}

// 洗练
func (self *ModHorse) washHorse(ctrl string, keyId int) {
	horseData, ok := self.Sql_Horse.info[keyId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_FAVOURITE_DOES_NOT_EXIST"))
		return
	}

	// 获取配置
	config, ok := GetCsvMgr().HorseBattleSteedMap[horseData.Type]
	if !ok {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_HORSE_NO_CONFIG"), horseData.Type))
		return
	}

	// 判断品质
	if config.Quality <= 4 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THE_QUALITY_IS_NOT_ENOUGH"))
		return
	}

	// 检查消耗
	itemId := GetCsvMgr().ClearWarHorseId
	itemNum := GetCsvMgr().ClearWarHorseNum
	// 检查配置
	if itemId == 0 {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_HORSE_NO_CONFIG2"), GetCsvMgr().ClearWarHorseId))
		return
	}

	if itemNum == 0 {
		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_HORSE_NO_CONFIG3"), GetCsvMgr().ClearWarHorseNum))
		return
	}

	// 检查消耗
	if self.player.GetObjectNum(itemId) < itemNum {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_TIGER_LACK_OF_PROPS"))
		return
	}

	// 扣除物品
	items := self.player.RemoveObjectSimple(itemId, itemNum, "魔宠洗练", horseData.Type, 0, 0)
	horseData.AttLst = self.GetHorseAttr(config.ExtractListGroup, config.ExtractAttNum)

	GetServer().SqlLog(self.player.GetUid(), LOG_HORSE_WASH, horseData.Type, 0, 0, "魔宠洗练", 0, 0, self.player)

	var msg S2C_WashHero
	msg.Cid = ctrl
	msg.Horse = horseData
	msg.Items = items
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

// 保存洗练
func (self *ModHorse) saveSwitch(ctrl string, keyId int) {
	horseData, ok := self.Sql_Horse.info[keyId]
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_FAVOURITE_DOES_NOT_EXIST"))
		return
	}

	if horseData.RandHorseId == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_HORSE_THERE_ARE_NO_TRANSFORMATION_ATTRIBUTES"))
		return
	}

	horseData.Type = horseData.RandHorseId
	csv_horse, ok := GetCsvMgr().Data["Horse_BattleSteed"][horseData.Type]
	if ok {
		//horseData.Skill = HF_Atoi(csv_horse["skill_id"]) + HF_Atoi(csv_horse["skill_lv"]) - 1
		horseData.Skill = HF_Atoi(csv_horse["skill_id"])
	}
	horseData.RandHorseId = 0
	self.CountFight(horseData)

	var msg S2C_SaveSwitch
	msg.Cid = ctrl
	msg.Horse = horseData
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
	self.Caculate(horseData, ReasonSaveSwicthHorse)
}

// 计算战力值
func (self *ModHorse) CountFight(info *JS_HorseInfo) {
	horseId := info.Id
	attr := self.getAttr(horseId)
	fightNum, ok := attr[99]
	if ok {
		info.Fight = fightNum.AttValue
	} else {
		info.Fight = 0
	}
}

func (self *ModHorse) GetStatisticsValue2021() int {
	//! 次数判断
	csvlevel, ok := GetCsvMgr().Data["Horse_Judge_level"][self.Sql_Horse.Level]
	if !ok {
		return 0
	}
	return HF_Atoi(csvlevel["daily_normal_call"]) - self.Sql_Horse.SummonNormal
}

func (self *ModHorse) GetStatisticsValue2030() (itemNum int, countNum int) {
	relItemNum := 0
	relCountNum := 0

	//! 次数判断
	csvlevel, ok := GetCsvMgr().Data["Horse_Judge_level"][self.Sql_Horse.Level]
	if !ok {
		return relItemNum, relCountNum
	}

	totalNum := GetCsvMgr().GetVipHorseCall(self.player.Sql_UserBase.Vip)
	relCountNum = HF_Atoi(csvlevel["daily_higt_call"]) + totalNum - self.Sql_Horse.SummonSenior

	csv, ok := GetCsvMgr().Data["Horse_Judge_call"][2001]
	itemid2 := HF_Atoi(csv["cost_item2"])
	relItemNum = self.player.GetObjectNum(itemid2)

	return relItemNum, relCountNum
}

/*
func (self *ModHorse) CountHorseFight() (int64, bool) {
	count := (int64)(0)

	teamPos := self.player.getTeamPos()
	if teamPos == nil {
		return count, false
	}

	if len(teamPos.TeamAttr) <= 0 {
		return count, false
	}

	for _, v := range teamPos.TeamAttr {
		horseId := v.HorseKeyId

		pAttr := &AttrWrapper{
			Base:     make([]float32, 32),
			Ext:      make(map[int]float32),
			Per:      make(map[int]float32),
			Energy:   0,
			FightNum: 0,
		}
		attr := self.getAttr(horseId)

		ProcAtt(attr, pAttr)
		count += pAttr.FightNum
	}

	if self.Sql_Horse.HorseTotalFight < count {
		self.Sql_Horse.HorseTotalFight = count

		self.player.HandleTask(HorseTotalFight, 0, 0, 0)

		return count, true
	}

	return count, false
}
*/
