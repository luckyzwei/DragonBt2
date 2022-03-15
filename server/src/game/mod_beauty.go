package game

import (
	"encoding/json"
	"fmt"
	"time"
)

//! 圣物数据库
type San_Beauty struct {
	Uid         int64
	Beautyinfo  string
	Lastupdtime int64
	Fight       int64
	Count       int

	beautyinfo []BeautyInfo

	DataUpdate
}

type BeautyInfo struct {
	Beautyid       int   `json:"beautyid"`
	Stage_lv       int   `json:"stage_lv"`
	Treasurea_lv   []int `json:"treasurea_lv"`
	Legend_Chapter int   `json:"legend_chapter"`
	Legend_Index   int   `json:"legend_index"`
	Legend_Finish  bool  `json:"legend_finish"`
	Fight          int64 `json:"fight"` // 圣物战力
}

type JS_BeautyTop struct {
	Uid       int64  `json:"uid"`
	Uname     string `json:"uname"`
	Iconid    int    `json:"iconid"`
	Portrait  int    `json:"portrait"` // 边框  20190412 by zy
	Level     int    `json:"level"`
	Camp      int    `json:"camp"`
	Fight     int64  `json:"fight"`
	Vip       int    `json:"vip"`
	LastRank  int    `json:"-"`
	UnionName string `json:"union_name"`
}

//! 圣物
type ModBeauty struct {
	player     *Player
	Sql_Beauty San_Beauty         //! 数据库结构
	AttrMap    map[int]*Attribute //! 总属性, 登录以及操作成功后重新计算
}

func (self *ModBeauty) OnGetData(player *Player) {
	self.player = player
	if self.AttrMap == nil {
		self.AttrMap = make(map[int]*Attribute)
	}

	sql := fmt.Sprintf("select * from `san_userbeauty3` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Beauty, "san_userbeauty3", self.player.ID)

	if self.Sql_Beauty.Uid <= 0 {
		self.Sql_Beauty.Uid = self.player.ID
		self.Sql_Beauty.Lastupdtime = time.Now().Unix()
		self.Sql_Beauty.Fight = 0
		self.Sql_Beauty.beautyinfo = make([]BeautyInfo, 0)
		self.Encode()
		InsertTable("san_userbeauty3", &self.Sql_Beauty, 0, true)
	} else {
		self.Decode()
		self.CheckNew()
	}

	self.Sql_Beauty.Init("san_userbeauty3", &self.Sql_Beauty, true)
	self.countFight()
}

func (self *ModBeauty) OnGetOtherData() {

}

func (self *ModBeauty) CheckNew() {
	for i := 0; i < len(self.Sql_Beauty.beautyinfo); i++ {
		if self.Sql_Beauty.beautyinfo[i].Legend_Finish {
			config := GetCsvMgr().GetHolyLegend(self.Sql_Beauty.beautyinfo[i].Beautyid, self.Sql_Beauty.beautyinfo[i].Legend_Chapter+1)
			if config != nil {
				self.Sql_Beauty.beautyinfo[i].Legend_Chapter += 1
				self.Sql_Beauty.beautyinfo[i].Legend_Index = 1
				self.Sql_Beauty.beautyinfo[i].Legend_Finish = false
			}
		}
	}
}

func (self *ModBeauty) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "beauty_treasurea_gem_up": //! 宝物一键升级
		var c2s_msg C2S_Beauty_TreasureaUpLevel
		json.Unmarshal(body, &c2s_msg)
		self.TreasureaUpLevel_Gem(c2s_msg.Beautyid, c2s_msg.TreasureaIndex)
		return true
	case "beauty_gem_up": //! 圣物一键升级
		var c2s_msg C2S_Beauty_UpLevel
		json.Unmarshal(body, &c2s_msg)
		self.BeautyUpLevel_Gem(c2s_msg.Beautyid)
		return true
	case "beauty_treasurea_up": //! 宝物升级
		var c2s_msg C2S_Beauty_TreasureaUpLevel
		json.Unmarshal(body, &c2s_msg)
		self.TreasureaUpLevel(c2s_msg.Beautyid, c2s_msg.TreasureaIndex)
		return true
	case "beauty_up": //! 圣物升级
		LogDebug("beauty_up proc")
		var c2s_msg C2S_Beauty_UpLevel
		json.Unmarshal(body, &c2s_msg)
		self.BeautyUpLevel(c2s_msg.Beautyid)
		return true
	case "beautytop":
		var c2s_msg C2S_BeautyTop
		json.Unmarshal(body, &c2s_msg)
		self.GetBeautyTop(c2s_msg.Ver)
		return true
	case "legendover":
		var c2s_msg C2S_Beauty_LegendOver
		json.Unmarshal(body, &c2s_msg)
		self.LegendFightOver(c2s_msg.Beautyid, c2s_msg.Chapter, c2s_msg.Index)
		return true
	}

	return false
}

func (self *ModBeauty) OnSave(sql bool) {
	self.Encode()
	self.Sql_Beauty.Update(sql)
}

func (self *ModBeauty) OnRefresh() {
}

func (self *ModBeauty) Decode() {
	//! 将数据库数据写入data
	json.Unmarshal([]byte(self.Sql_Beauty.Beautyinfo), &self.Sql_Beauty.beautyinfo)
}

func (self *ModBeauty) Encode() {
	//! 将data数据写入数据库
	s, _ := json.Marshal(&self.Sql_Beauty.beautyinfo)
	self.Sql_Beauty.Beautyinfo = string(s)
}

//！添加圣物数据
func (self *ModBeauty) AddBeauty(beautyid int) {
	for _, value := range self.Sql_Beauty.beautyinfo {
		if value.Beautyid == beautyid {
			return
		}
	}
	var info BeautyInfo
	info.Beautyid = beautyid
	info.Stage_lv = 0
	for i := 0; i < 6; i++ {
		info.Treasurea_lv = append(info.Treasurea_lv, 0)
	}
	info.Legend_Chapter = 1
	info.Legend_Index = 1
	info.Legend_Finish = false
	self.Sql_Beauty.beautyinfo = append(self.Sql_Beauty.beautyinfo, info)
}

//! 传奇关卡结束
func (self *ModBeauty) LegendFightOver(beautyid int, chapter int, index int) {
	//! 检测beautyid是否有效
	_, ok := GetCsvMgr().GetHolyAdvance(beautyid, 0)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_HOLY_ID_ERROR"))
		return
	}
	self.AddBeauty(beautyid)
	//! 检查章节和index是否与服务器一直
	chackok := true
	for _, value := range self.Sql_Beauty.beautyinfo {
		if value.Beautyid == beautyid {
			if value.Legend_Chapter != chapter || value.Legend_Index != index {
				chackok = false
			}
		}
	}
	if !chackok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_DATA_ERROR"))
		return
	}
	//! 检查章节等级需求
	csv_legend := GetCsvMgr().GetHolyLegend(beautyid, chapter)
	if csv_legend == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_CHAPTER_DATA_ANOMALIES"))
		return
	}
	if csv_legend.Chaptercondition > self.player.Sql_UserBase.Level {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PASS_INSUFFICIENT_GRADE"))
		return
	}
	csv_legendlevel := GetCsvMgr().GetHolyLengendLevel(csv_legend.Levelgroup, index)
	if csv_legendlevel == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_PVP_DATA_ABNORMITY"))
		return
	}
	//! 掉落处理
	var lstItem []PassItem
	for i := 0; i < len(csv_legendlevel.Items); i++ {
		itemid := csv_legendlevel.Items[i]
		if len(csv_legendlevel.Items) != len(csv_legendlevel.Nums) {
			continue
		}
		if itemid == 0 {
			continue
		}
		itemnum := csv_legendlevel.Nums[i]
		if itemnum == 0 {
			continue
		}
		lstItem = append(lstItem, PassItem{itemid, itemnum})
	}
	//! 添加物品
	for i := 0; i < len(lstItem); i++ {
		self.player.AddObject(lstItem[i].ItemID, lstItem[i].Num, beautyid, 0, 0, "圣物关卡通关")
	}
	//! 数据处理
	cur_index := 0
	for i := 0; i < len(self.Sql_Beauty.beautyinfo); i++ {
		if self.Sql_Beauty.beautyinfo[i].Beautyid == beautyid {
			cur_index = i
			config := GetCsvMgr().GetHolyLengendLevel(csv_legend.Levelgroup, index+1)
			if config == nil {
				pHoly := GetCsvMgr().GetHolyLegend(beautyid, chapter)
				if pHoly == nil {
					self.Sql_Beauty.beautyinfo[i].Legend_Finish = true
				} else {
					self.Sql_Beauty.beautyinfo[i].Legend_Chapter++
					self.Sql_Beauty.beautyinfo[i].Legend_Index = 1
				}
			} else {
				self.Sql_Beauty.beautyinfo[i].Legend_Index++
			}
		}
	}
	//curfight := self.player.GetModule("city").(*ModCity).GetCurFight()

	//GetServer().sendLog_pve_ex(self.player, int(curfight/100), csv_legendlevel.LevelName, "true", 0, "战斗", "圣物传记")

	self.player.HandleTask(LegendLevelTask, csv_legendlevel.Level, 0, 0)
	//! 发消息
	var msg S2C_Beauty_LegendOver
	msg.BeautyId = beautyid
	msg.Cid = "legendover"
	msg.Award = lstItem
	msg.Chapter = self.Sql_Beauty.beautyinfo[cur_index].Legend_Chapter
	msg.Index = self.Sql_Beauty.beautyinfo[cur_index].Legend_Index
	msg.Finish = self.Sql_Beauty.beautyinfo[cur_index].Legend_Finish

	self.player.SendMsg("legendover", HF_JtoB(&msg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEATUY_PASS_FINISH, beautyid, 0, 0, "圣物关卡通关", 0, 0, self.player)

	//self.player.GetModule("statistics").(*ModStatistics).CalValue(ST_1090, true, 0)
}

//! 宝物一键升级
func (self *ModBeauty) TreasureaUpLevel_Gem(beautyid int, index int) {
	if index <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_TREASURE_BIDDING_ERROR"))
		return
	}
	self.AddBeauty(beautyid)

	curlevel := 0
	for _, value := range self.Sql_Beauty.beautyinfo {
		if value.Beautyid == beautyid {
			curlevel = value.Treasurea_lv[index-1]
		}
	}
	csv_beauty, ok := GetCsvMgr().GetHolyAdvance(beautyid, 1)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	//	csv_cur, _ := GetCsvMgr().GetHarem_TreasureaAdv(HF_Atoi(csv_beauty[fmt.Sprintf("treasurea%d", index)]), curlevel)
	if index > len(csv_beauty.Treasureas) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BUILD_SUBSCRIPT_ERROR"))
		return
	}

	partId := csv_beauty.Treasureas[index-1]
	csv_finial, ok := GetCsvMgr().GetHolyPartsAdv(partId, curlevel+1)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_THE_TREASURE_HAS_REACHED_ITS"))
		return
	}

	csv, ok := GetCsvMgr().GetHolyPartsAdv(partId, curlevel)
	if !ok {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_BEAUTY_COMPONENT_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}

	var needitem []PassItem
	var chgitem []PassItem
	for i := 0; i < len(csv.Costitems); i++ {
		itemid := csv.Costitems[i]
		if itemid == 0 {
			continue
		}

		num := csv.Costnums[i]
		if num == 0 {
			continue
		}

		if self.player.GetObjectNum(itemid) < num {
			chgitem = append(chgitem, PassItem{itemid, num - self.player.GetObjectNum(itemid)})
			needitem = append(needitem, PassItem{itemid, -self.player.GetObjectNum(itemid)})
		} else {
			needitem = append(needitem, PassItem{itemid, -num})
		}
	}

	//! 计算购买物品需要多少元宝
	var needgem = 0
	for i := 0; i < len(chgitem); i++ {
		needgem += chgitem[i].Num * GetCsvMgr().GetItemGemPrice(chgitem[i].ItemID)
	}
	needgem = needgem * 9 / 10

	if self.player.Sql_UserBase.Gem < needgem {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_INSUFFICIENT_TREASURE"))
		return
	}
	needitem = append(needitem, PassItem{91000002, -needgem})

	for i := 0; i < len(needitem); i++ {
		if needitem[i].ItemID == 91000002 {
			self.player.AddObject(needitem[i].ItemID, needitem[i].Num, 6, 0, 0, "圣物宝物升级")
		} else {
			self.player.AddObject(needitem[i].ItemID, needitem[i].Num, 6, 0, 0, "圣物宝物升级")
		}
	}
	//! 修改数据
	ret_curlevel := 0
	for i := 0; i < len(self.Sql_Beauty.beautyinfo); i++ {
		if self.Sql_Beauty.beautyinfo[i].Beautyid == beautyid {
			self.Sql_Beauty.beautyinfo[i].Treasurea_lv[index-1] += 1
			ret_curlevel = self.Sql_Beauty.beautyinfo[i].Treasurea_lv[index-1]
			//!圣物宝物升级
			//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_UP_TREASURE, 1, beautyid, index, "圣物宝物一键升级", ret_curlevel, ret_curlevel, self.player)
			break
		}
	}

	self.countFight()
	self.player.countAllHero()
	//GetServer().sendLog_beautyequiplevelup(self.player, beautyid, index, ret_curlevel)
	//! 全服广播
	if csv_finial.Announcement > 0 {
		if csv_finial.Stagelvshow == 1 {
			GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_BEAUTY_ACTIVATE"),
				self.player.Sql_UserBase.UName, csv_beauty.Name,
				csv.Name), 0, 1)
		} else {
			GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_BEAUTY_LEVEL"),
				self.player.Sql_UserBase.UName, csv_beauty.Name,
				csv.Name, csv_finial.Name), 0, 1)
		}
	}

	//self.player.HandleTask(BeautyAdvance, 0, curlevel+1, 0)

	//self.player.HandleTask(TreasureaAdvance, curlevel+1, 0, 0)
	//! 返回消息
	var msg S2C_Beauty_TreasureaUpLevel
	msg.Cid = "beauty_treasurea_up"
	msg.Cost = needitem
	msg.CurLevel = ret_curlevel
	msg.Index = index
	msg.Fight = self.getFight(beautyid)
	msg.TotalFight = self.getTotalFight()

	self.player.SendMsg("beauty_treasurea_gem_up", HF_JtoB(&msg))

	if curlevel == 0 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_TREASURE_ACTIVE, beautyid, 0, 0, "圣物部件激活", 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_TREASURE_UP, curlevel+1, beautyid, 0, "圣物部件升级", 0, 0, self.player)
	}
}

func (self *ModBeauty) getTotalFight() int64 {
	fight, ok := self.AttrMap[99]
	if !ok {
		return 0
	}
	return fight.AttValue
}

func (self *ModBeauty) getFight(id int) int64 {
	for _, value := range self.Sql_Beauty.beautyinfo {
		if value.Beautyid == id {
			return value.Fight
		}
	}
	return 0
}

//! 宝物升级
func (self *ModBeauty) TreasureaUpLevel(beautyid int, index int) {
	if index <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_TREASURE_BIDDING_ERROR"))
		return
	}
	self.AddBeauty(beautyid)

	curlevel := 0
	for _, value := range self.Sql_Beauty.beautyinfo {
		if value.Beautyid == beautyid {
			curlevel = value.Treasurea_lv[index-1]
		}
	}
	pHolyAdvance, ok := GetCsvMgr().GetHolyAdvance(beautyid, 1)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	if index > len(pHolyAdvance.Treasureas) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BUILD_SUBSCRIPT_ERROR"))
		return
	}

	partId := pHolyAdvance.Treasureas[index-1]
	csv_finial, ok := GetCsvMgr().GetHolyPartsAdv(pHolyAdvance.Treasureas[index-1], curlevel+1)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_THE_TREASURE_HAS_REACHED_ITS"))
		return
	}

	csv, ok := GetCsvMgr().GetHolyPartsAdv(partId, curlevel)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_THE_TREASURE_DOES_NOT_EXIST"))
		return
	}

	if err := self.player.HasObjectOk(csv.Costitems, csv.Costnums); err != nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_THIS_PART_HAS_NOT_BEEN"))
		return
	}

	needitem := self.player.RemoveObjectLst(csv.Costitems, csv.Costnums, "圣物部件升级", curlevel+1, beautyid, 0)

	//! 修改数据
	ret_curlevel := 0
	for i := 0; i < len(self.Sql_Beauty.beautyinfo); i++ {
		if self.Sql_Beauty.beautyinfo[i].Beautyid == beautyid {
			self.Sql_Beauty.beautyinfo[i].Treasurea_lv[index-1] += 1
			ret_curlevel = self.Sql_Beauty.beautyinfo[i].Treasurea_lv[index-1]

			//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_UP_TREASURE, index, beautyid, 0, "圣物部件升级", 0, ret_curlevel, self.player)
			break
		}
	}

	self.countFight()
	self.player.countAllHero()
	//GetServer().sendLog_beautyequiplevelup(self.player, beautyid, index, ret_curlevel)
	//! 全服广播
	if csv_finial.Announcement > 0 && self.player.Sql_UserBase.Camp != 0 {
		if csv_finial.Stagelvshow == 1 {
			GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_BEAUTY_ACTIVATE"),
				self.player.Sql_UserBase.UName, pHolyAdvance.Name, csv.Name), 0, 1)
		} else {
			strings:=GetCsvMgr().GetText("STR_BEAUTY_LEVEL")
			text := fmt.Sprintf(strings, self.player.Sql_UserBase.UName, pHolyAdvance.Name,
				csv.Name, csv_finial.Name)
			GetServer().Notice(text, 0, 1)
		}
	}

	//self.player.HandleTask(BeautyAdvance, 0, curlevel+1, 0)
	//self.player.HandleTask(TreasureaAdvance, curlevel+1, 0, 0)
	var msg S2C_Beauty_TreasureaUpLevel
	msg.Cid = "beauty_treasurea_up"
	msg.Cost = needitem
	msg.CurLevel = ret_curlevel
	msg.Index = index
	msg.Fight = self.getFight(beautyid)
	msg.TotalFight = self.getTotalFight()

	self.player.SendMsg("beauty_treasurea_up", HF_JtoB(&msg))

	if curlevel == 0 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_TREASURE_ACTIVE, beautyid, 0, 0, "圣物部件激活", 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_TREASURE_UP, curlevel+1, beautyid, 0, "圣物部件升级", 0, 0, self.player)
	}
}

//! 圣物一键升级
func (self *ModBeauty) BeautyUpLevel_Gem(beautyid int) {
	//! 检测beautyid是否有效
	_, ok := GetCsvMgr().GetHolyAdvance(beautyid, 0)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_SACRED_LOGISTICS_WATER_NUMBER_ERROR"))
		return
	}
	self.AddBeauty(beautyid)
	//! 获取圣物当前等级
	var reallevel = 0
	for _, value := range self.Sql_Beauty.beautyinfo {
		if value.Beautyid == beautyid {
			reallevel = value.Stage_lv
			break
		}
	}

	csv_finial, ok := GetCsvMgr().GetHolyAdvance(beautyid, reallevel+1)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_HAS_RISEN"))
		return
	}

	treasurea_lv := make([]int, 0)
	for _, value := range self.Sql_Beauty.beautyinfo {
		if value.Beautyid == beautyid {
			treasurea_lv = value.Treasurea_lv
			break
		}
	}

	csv, _ := GetCsvMgr().GetHolyAdvance(beautyid, reallevel)
	for _, value := range treasurea_lv {
		if value < csv_finial.Stagelv {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_INSUFFICIENT_GRADE_OF_TREASURES"))
			return
		}
	}
	var needitem []PassItem
	var chgitem []PassItem
	for i := 0; i < len(csv.Items); i++ {
		itemid := csv.Items[i]
		if itemid == 0 {
			continue
		}

		num := csv.Itemcosts[i]

		if self.player.GetObjectNum(itemid) < num {
			chgitem = append(chgitem, PassItem{itemid, num - self.player.GetObjectNum(itemid)})
			needitem = append(needitem, PassItem{itemid, -self.player.GetObjectNum(itemid)})
		} else {
			needitem = append(needitem, PassItem{itemid, -num})
		}
	}

	//! 计算购买物品需要多少元宝
	var needgem = 0
	for i := 0; i < len(chgitem); i++ {
		needgem += chgitem[i].Num * GetCsvMgr().GetItemGemPrice(chgitem[i].ItemID)
	}
	needgem = needgem * 9 / 10
	if self.player.Sql_UserBase.Gem < needgem {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_INSUFFICIENT_TREASURE"))
		return
	}
	needitem = append(needitem, PassItem{91000002, -needgem})

	for i := 0; i < len(needitem); i++ {
		self.player.AddObject(needitem[i].ItemID, needitem[i].Num, reallevel+1, beautyid, 0, "圣物升级")
	}

	//! 修改数据
	ret_curlevel := 0

	for i := 0; i < len(self.Sql_Beauty.beautyinfo); i++ {
		if self.Sql_Beauty.beautyinfo[i].Beautyid == beautyid {
			self.Sql_Beauty.beautyinfo[i].Stage_lv += 1
			ret_curlevel = self.Sql_Beauty.beautyinfo[i].Stage_lv

			//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_UP_LEVEL, 1, beautyid, i, "圣物元宝升级",
			//	ret_curlevel, ret_curlevel, self.player)
			break
		}
	}
	self.countFight()
	self.player.countAllHero()
	//! 全服广播
	if csv_finial.Announcement > 0 {
		if csv_finial.Stagelvshow == 1 {
			GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_BEAUTY_ACTIVATE_SUCCESS"),
				self.player.Sql_UserBase.UName, csv.Name), 0, 1)
		} else {
			GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_BEAUTY_LEVEL_SUCCESS"),
				self.player.Sql_UserBase.UName, csv.Name, csv_finial.Stagelvshow), 0, 1)
		}
	}

	self.player.HandleTask(BeautyAdvance, ret_curlevel, 0, 0)
	var msg S2C_Beauty_UpLevel
	msg.BeautyId = beautyid
	msg.Cid = "beauty_up"
	msg.Cost = needitem
	msg.CurLevel = ret_curlevel
	msg.Fight = self.getFight(beautyid)
	msg.TotalFight = self.getTotalFight()

	self.player.SendMsg("beauty_up", HF_JtoB(&msg))

	if reallevel == 0 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_ACTIVE, beautyid, 0, 0, "圣物激活", 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_UP, ret_curlevel, beautyid, 0, "圣物升级", 0, 0, self.player)
	}
}

//! 圣物升级
func (self *ModBeauty) BeautyUpLevel(beautyid int) {
	//! 检测beautyid是否有效
	_, ok := GetCsvMgr().GetHolyAdvance(beautyid, 0)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_SACRED_LOGISTICS_WATER_NUMBER_ERROR"))
		return
	}
	self.AddBeauty(beautyid)
	//! 获取圣物当前等级
	var reallevel = 0
	for _, value := range self.Sql_Beauty.beautyinfo {
		if value.Beautyid == beautyid {
			reallevel = value.Stage_lv
			break
		}
	}
	csv_finial, ok := GetCsvMgr().GetHolyAdvance(beautyid, reallevel+1)
	if !ok {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_HAS_RISEN"))
		return
	}

	treasurea_lv := make([]int, 0)
	for _, value := range self.Sql_Beauty.beautyinfo {
		if value.Beautyid == beautyid {
			treasurea_lv = value.Treasurea_lv
			break
		}
	}
	LogDebug("reallevel = ", reallevel)
	csv, _ := GetCsvMgr().GetHolyAdvance(beautyid, reallevel)
	//LogDebug("csv = ", csv)
	for _, value := range treasurea_lv {
		if value < csv_finial.Stagelv {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_INSUFFICIENT_GRADE_OF_TREASURES"))
			return
		}
	}
	needitem := make([]PassItem, 0)
	for i := 0; i < len(csv.Items); i++ {
		itemid := csv.Items[i]
		if itemid == 0 {
			continue
		}

		num := csv.Itemcosts[i]

		if self.player.GetObjectNum(itemid) < num {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_BEAUTY_INSUFFICIENT_UPGRADE_ITEMS"))
			return
		}

		needitem = append(needitem, PassItem{itemid, -num})
	}

	for i := 0; i < len(needitem); i++ {
		self.player.AddObject(needitem[i].ItemID, needitem[i].Num, reallevel, beautyid, 0, "圣物升级")
	}
	//! 修改数据
	ret_curlevel := 0
	for i := 0; i < len(self.Sql_Beauty.beautyinfo); i++ {
		if self.Sql_Beauty.beautyinfo[i].Beautyid == beautyid {
			self.Sql_Beauty.beautyinfo[i].Stage_lv += 1
			ret_curlevel = self.Sql_Beauty.beautyinfo[i].Stage_lv

			//GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_UP_LEVEL, 1, beautyid, i, "圣物升级",
			//	ret_curlevel, ret_curlevel, self.player)
			break
		}
	}
	self.countFight()
	self.player.countAllHero()
	//GetServer().sendLog_beautyrank(self.player, beautyid, ret_curlevel)
	//! 全服广播
	if csv_finial.Announcement > 0 {
		if csv_finial.Stagelvshow == 1 {
			GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_BEAUTY_ACTIVATE_SUCCESS"),
				self.player.Sql_UserBase.UName, csv.Name), 0, 1)
		} else {
			GetServer().Notice(fmt.Sprintf(GetCsvMgr().GetText("STR_BEAUTY_LEVEL_SUCCESS"),
				self.player.Sql_UserBase.UName, csv.Name, csv_finial.Stagelvshow), 0, 1)
		}
	}
	self.player.HandleTask(BeautyAdvance, ret_curlevel, 0, 0)
	var msg S2C_Beauty_UpLevel
	msg.BeautyId = beautyid
	msg.Cid = "beauty_up"
	msg.Cost = needitem
	msg.CurLevel = ret_curlevel
	msg.Fight = self.getFight(beautyid)
	msg.TotalFight = self.getTotalFight()

	self.player.SendMsg("beauty_up", HF_JtoB(&msg))

	if reallevel == 0 {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_ACTIVE, beautyid, 0, 0, "圣物激活", 0, 0, self.player)
	} else {
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_BEAUTY_UP, ret_curlevel, beautyid, 0, "圣物升级", 0, 0, self.player)
	}
}

func (self *ModBeauty) GetBeautyLevel(level int) int {
	amount := 0
	for i := 0; i < len(self.Sql_Beauty.beautyinfo); i++ {
		if self.Sql_Beauty.beautyinfo[i].Stage_lv >= level {
			amount++
		}
	}
	return amount
}

func (self *ModBeauty) GetTreasureLevel(level int) int {
	amount := 0
	for i := 0; i < len(self.Sql_Beauty.beautyinfo); i++ {

		for t := 0; t < len(self.Sql_Beauty.beautyinfo[i].Treasurea_lv); t++ {
			if self.Sql_Beauty.beautyinfo[i].Treasurea_lv[t] >= level {
				amount++
			}
		}
	}
	return amount
}

func (self *ModBeauty) CountBeautyLevel() int {
	amount := 0
	for i := 0; i < len(self.Sql_Beauty.beautyinfo); i++ {
		amount += self.Sql_Beauty.beautyinfo[i].Stage_lv
	}
	return amount
}

//! 获取排行数据
func (self *ModBeauty) GetBeautyTop(ver int) {
	/*
		var msg S2C_BeautyTop
		msg.Cid = "beautytop"
		if ver != GetTopBeautyMgr().Topver {
			msg.Top = GetTopBeautyMgr().TopBeauty
		} else {
			msg.Top = make([]*JS_BeautyTop, 0)
		}
		self.player.SendMsg("beautytop", HF_JtoB(&msg))
	*/
}

//////////////////////////////
func (self *ModBeauty) SendInfo() {
	var msg S2C_BeautyInfo
	msg.Cid = "beauty3lst"
	msg.Info.Uid = self.player.ID
	msg.Info.Fight = self.Sql_Beauty.Fight
	msg.Info.Beautyinfo = self.Sql_Beauty.beautyinfo
	msg.Info.Lastupdtime = self.Sql_Beauty.Lastupdtime
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("beauty3lst", smsg)
}

// 获取圣物最高等级
func (self *ModBeauty) GetMaxLv() int {
	maxLv := 0
	beautyinfo := self.Sql_Beauty.beautyinfo
	for _, value := range beautyinfo {
		for _, level := range value.Treasurea_lv {
			if maxLv < level {
				maxLv = level
			}
		}
	}
	return maxLv
}

func (self *ModBeauty) addAttEx(inparam map[int]*Attribute, attMap map[int]*Attribute) {
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

// 获取圣物属性: 圣物属性, 宝物属性, 技能属性
func (self *ModBeauty) cacAttr() map[int]*Attribute {
	attMap := make(map[int]*Attribute)
	base := self.getBaseAttr()
	self.addAttEx(base, attMap)

	return attMap
}

// 计算单个圣物战斗力
func (self *ModBeauty) getBaseBeauty(value *BeautyInfo) map[int]*Attribute {
	res := make(map[int]*Attribute)
	for i := 0; i <= value.Stage_lv; i++ {
		csv_beautyadv, ok := GetCsvMgr().GetHolyAdvance(value.Beautyid, i)
		if !ok {
			continue
		}
		AddAttrDirect(res, csv_beautyadv.Extraattributestypes, csv_beautyadv.Extraattributeaddvalues)
	}
	return res
}

// 圣物属性
func (self *ModBeauty) getBaseAttr() map[int]*Attribute {
	res := make(map[int]*Attribute)
	for i, value := range self.Sql_Beauty.beautyinfo {
		att := make(map[int]*Attribute)
		base := self.getBaseBeauty(&value)
		AddAttrMapHelper(att, base)
		skill := self.getBaseSkill(&value)
		AddAttrMapHelper(att, skill)
		treasure := self.getBaseTreasure(&value)
		AddAttrMapHelper(att, treasure)
		fight, ok := att[99]
		if ok {
			self.Sql_Beauty.beautyinfo[i].Fight = fight.AttValue
		}
		AddAttrMapHelper(res, att)
	}
	return res
}

func (self *ModBeauty) getBaseSkill(value *BeautyInfo) map[int]*Attribute {
	res := make(map[int]*Attribute)
	csv_beautyadv, ok := GetCsvMgr().GetHolyAdvance(value.Beautyid, value.Stage_lv)
	if !ok {
		return res
	}

	//config, ok2 := GetCsvMgr().SkillConfig[csv_beautyadv.Skillid]
	config, ok2 := GetCsvMgr().SkillConfigMap[csv_beautyadv.Skillid]
	if !ok2 {
		return res
	}

	var skillCounts []int64
	for i := 0; i < len(config.SkillValueType); i++ {
		skillCounts = append(skillCounts, int64(config.SkillCount[i]))
	}

	AddAttrDirect(res, config.SkillValueType, skillCounts)
	return res
}

func (self *ModBeauty) getBaseTreasure(value *BeautyInfo) map[int]*Attribute {
	res := make(map[int]*Attribute)
	//! 圣物宝物属性
	csv_beautyadv, ok := GetCsvMgr().GetHolyAdvance(value.Beautyid, value.Stage_lv)
	if !ok {
		return res
	}
	for k := 0; k < len(csv_beautyadv.Treasureas); k++ {
		treasurea := csv_beautyadv.Treasureas[k]
		pHolyParts, ok2 := GetCsvMgr().GetHolyPartsAdv(treasurea, value.Treasurea_lv[k])
		if !ok2 {
			continue
		}
		AddAttrDirect(res, pHolyParts.Attributetypes, pHolyParts.Attributes)
		for i := 1; i <= value.Treasurea_lv[k]; i++ {
			config, has := GetCsvMgr().GetHolyPartsAdv(treasurea, i)
			if !has {
				continue
			}
			AddAttrDirect(res, config.Extraattributestypes, config.Extraattributeaddvalues)
		}
	}
	return res
}

func (self *ModBeauty) procAtt(attMap map[int]*Attribute) int64 {
	var fightNum int64 = 0
	for _, pAttribute := range attMap {
		valuetype := pAttribute.AttType
		if valuetype != 99 {
			continue
		}

		value := pAttribute.AttValue
		if value == 0 {
			continue
		}
		fightNum += value
	}
	return fightNum
}

// 计算战斗力
func (self *ModBeauty) countFight() int64 {
	attMap := self.cacAttr()
	fightNum := self.procAtt(attMap)
	self.Sql_Beauty.Fight = fightNum
	//GetTopBeautyMgr().updateRank(self.Sql_Beauty.Fight, self.player)
	//GetTopBeautyMgr().updateCampRank()
	self.AttrMap = attMap
	return self.Sql_Beauty.Fight
}

// 获取圣物缓存属性
func (self *ModBeauty) getAttr() map[int]*Attribute {
	return self.AttrMap
}

//! 传奇关卡结束
func (self *ModBeauty) GetStatisticsValue1090() int {

	//没有数据表适用，为了提高性能只能写死了，否则遍历太多
	for i := 5001; i <= 5007; i++ {
		isHas := false
		//先看看有没打过的记录
		for _, value := range self.Sql_Beauty.beautyinfo {
			if value.Beautyid == i {
				isHas = true
				csv_legend := GetCsvMgr().GetHolyLegend(value.Beautyid, value.Legend_Chapter)
				if csv_legend != nil {
					return value.Beautyid
				}
			}
		}
		//如果没记录就看看是否可以开
		if !isHas {
			csv_legend := GetCsvMgr().GetHolyLegend(i, 1)
			if self.player.Sql_UserBase.Level >= csv_legend.Chaptercondition {
				csv_legendlevel := GetCsvMgr().GetHolyLengendLevel(csv_legend.Levelgroup, 1)
				if csv_legendlevel != nil {
					return i
				}
			}
		}
	}

	return 0
}
