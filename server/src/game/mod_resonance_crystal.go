package game

import (
	"encoding/json"
	"fmt"
	"sort"
	//"time"
)

const (
	RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX    = 5  // 祭司最大数量
	RESONANCE_CRYSTAL_RESONANCE_COUNT_BASE = 2  // 基础格子数量
	RESONANCE_CRYSTAL_RESONANCE_COUNT_MAX  = 60 // 最大格子数量
	RESONANCE_CRYSTAL_RESONANCE_CD         = 1  // 下阵cd时间 1天
)

const (
	RESONANCE_CRYSTAL_SET                 = "resonance_crystal_set"                 // 设置
	RESONANCE_CRYSTAL_CANCEL              = "resonance_crystal_cancel"              // 取消设置
	RESONANCE_CRYSTAL_INFO                = "resonance_crystal_info"                // 获取所有可用英雄信息
	RESONANCE_CRYSTAL_ADD_RESONANCE_COUNT = "resonance_crystal_add_resonance_count" // 增加格子
	RESONANCE_CRYSTAL_UPDATE_RESONANCE    = "resonance_crystal_update_resonance"    // 更新共鸣
	RESONANCE_CRYSTAL_UPDATE_PRIESTS      = "resonance_crystal_update_priests"      // 更新祭司
	RESONANCE_CRYSTAL_CLEAN_CD            = "resonance_crystal_clean_cd"            // 清理cd
	RESONANCE_CRYSTAL_FIGHT               = "resonance_crystal_fight"               // 同步战力
)

const (
	RESONANCE_CRYSTAL_ADD_GEM    = 1
	RESONANCE_CRYSTAL_ADD_NORMAL = 2
)

// 共鸣英雄结构体
type ResonanceHeros struct {
	HeroKey int   `json:"herokey"` // 英雄key值
	EndTime int64 `json:"endtime"` // cd结束时间
}

// 共鸣水晶
type San_ResonanceCrystal struct {
	Uid            int64  // 角色ID
	PriestsHeros   string // 祭司英雄
	ResonanceHeros string // 共鸣英雄
	ResonanceCount int    // 最大共鸣人数
	MaxFight       int64  // 历史最高战力
	MaxFightTime   int64  // 历史最高战力时间
	MaxFightAll    int64  // 共鸣水晶战力

	priestsHeros   []int             // 祭司英雄
	resonanceHeros []*ResonanceHeros // 共鸣英雄
	DataUpdate
}

// 共鸣水晶
type ModResonanceCrystal struct {
	player *Player // 玩家

	San_ResonanceCrystal San_ResonanceCrystal // 信息数据
}

// 获得数据
func (self *ModResonanceCrystal) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_resonancecrystal` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.San_ResonanceCrystal, "san_resonancecrystal", self.player.ID)

	if self.San_ResonanceCrystal.Uid <= 0 {
		self.San_ResonanceCrystal.Uid = self.player.ID
		self.San_ResonanceCrystal.ResonanceCount = RESONANCE_CRYSTAL_RESONANCE_COUNT_BASE
		self.Encode()
		InsertTable("san_resonancecrystal", &self.San_ResonanceCrystal, 0, true)
		self.San_ResonanceCrystal.Init("san_resonancecrystal", &self.San_ResonanceCrystal, true)
	} else {
		self.Decode()
		self.San_ResonanceCrystal.Init("san_resonancecrystal", &self.San_ResonanceCrystal, true)
	}

	nLen := len(self.San_ResonanceCrystal.resonanceHeros)
	if self.San_ResonanceCrystal.ResonanceCount > nLen {
		for i := 0; i < self.San_ResonanceCrystal.ResonanceCount-nLen; i++ {
			self.San_ResonanceCrystal.resonanceHeros = append(self.San_ResonanceCrystal.resonanceHeros, &ResonanceHeros{})
		}
	}
}

type topHeros []*Hero

func (s topHeros) Len() int      { return len(s) }
func (s topHeros) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s topHeros) Less(i, j int) bool {
	// 如果是已经上阵共鸣水晶的英雄则判断原等级
	leveli := s[i].HeroLv
	if s[i].UseType[HERO_USE_TYPE_CRYSTAL_RESONANCE] == 1 {
		leveli = s[i].OriginalLevel
	}
	levelj := s[j].HeroLv
	if s[j].UseType[HERO_USE_TYPE_CRYSTAL_RESONANCE] == 1 {
		levelj = s[j].OriginalLevel
	}

	if leveli > levelj {
		return true
	}

	if leveli < levelj {
		return false
	}

	if s[i].Fight > s[j].Fight {
		return true
	}

	if s[i].Fight < s[j].Fight {
		return false
	}

	if s[i].HeroId < s[j].HeroId {
		return true
	}

	return false
}

// 获得数据
func (self *ModResonanceCrystal) OnGetOtherData() {

}

// save
func (self *ModResonanceCrystal) Decode() {
	json.Unmarshal([]byte(self.San_ResonanceCrystal.PriestsHeros), &self.San_ResonanceCrystal.priestsHeros)
	json.Unmarshal([]byte(self.San_ResonanceCrystal.ResonanceHeros), &self.San_ResonanceCrystal.resonanceHeros)
}

// read
func (self *ModResonanceCrystal) Encode() {
	self.San_ResonanceCrystal.PriestsHeros = HF_JtoA(self.San_ResonanceCrystal.priestsHeros)
	self.San_ResonanceCrystal.ResonanceHeros = HF_JtoA(self.San_ResonanceCrystal.resonanceHeros)
}

// 存储
func (self *ModResonanceCrystal) OnSave(sql bool) {
	self.Encode()
	self.San_ResonanceCrystal.Update(sql)
}

// 消息
func (self *ModResonanceCrystal) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (self *ModResonanceCrystal) onReg(handlers map[string]func(body []byte)) {
	handlers[RESONANCE_CRYSTAL_SET] = self.SetResonanceHeros                 // 设置
	handlers[RESONANCE_CRYSTAL_CANCEL] = self.CancelResonanceHerosMsg        // 取消设置
	handlers[RESONANCE_CRYSTAL_INFO] = self.SendInfo                         // 获取信息
	handlers[RESONANCE_CRYSTAL_ADD_RESONANCE_COUNT] = self.AddResonanceCount // 增加格子
	handlers[RESONANCE_CRYSTAL_CLEAN_CD] = self.CleanCD                      // 清理cd
}

// 设置祭司英雄
func (self *ModResonanceCrystal) SetPriestsHeros() {
	// 取消之前的全部祭司英雄
	self.CancelPriestsHeros()
	if len(self.San_ResonanceCrystal.priestsHeros) <= 0 {
		// 英雄排序
		data := topHeros{}
		heros := self.player.GetModule("hero").(*ModHero).Sql_Hero.info
		for _, v := range heros {
			data = append(data, v)
		}
		sort.Sort(topHeros(data))
		// 设置祭司英雄
		nCount := 0
		for _, v := range data {
			// 虚空英雄 且有共鸣对象则跳过
			if v.VoidHero != 0 && v.Resonance != 0 {
				continue
			}
			// 如果是共鸣英雄
			if v.UseType[HERO_USE_TYPE_CRYSTAL_RESONANCE] == 1 {
				index := self.GetResonanceIndex(v.HeroKeyId)
				if index >= 0 {
					// 取消共鸣英雄
					self.CancelResonanceHeros(index, v.HeroKeyId, false)
				}
			}
			// 设置祭司
			v.UseType[HERO_USE_TYPE_CRYSTAL_PRIESTS] = 1
			self.San_ResonanceCrystal.priestsHeros = append(self.San_ResonanceCrystal.priestsHeros, v.HeroKeyId)
			nCount++
			// 个数达到
			if nCount >= RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX {
				break
			}
		}
		if nCount > 0 {
			var backmsg S2C_ResonanceCrystaUpdatePriests
			backmsg.Cid = RESONANCE_CRYSTAL_UPDATE_PRIESTS
			backmsg.Heros = self.San_ResonanceCrystal.priestsHeros
			self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
		}
	} else {
		return
	}
	// 计算设置等级
	self.CheckLevel()
	self.CheckMaxFight()
}

// 取消祭司英雄
func (self *ModResonanceCrystal) CancelPriestsHeros() {
	if len(self.San_ResonanceCrystal.priestsHeros) > 0 {
		for _, v := range self.San_ResonanceCrystal.priestsHeros {
			hero := self.player.getHero(v)
			if nil != hero {
				// 恢复使用类型
				hero.UseType[HERO_USE_TYPE_CRYSTAL_PRIESTS] = 0
			}
		}
	}
	// 清理数据
	self.San_ResonanceCrystal.priestsHeros = []int{}
}

func (self *ModResonanceCrystal) GetPriestsEquipFight() int64 {
	fight := int64(0)
	if len(self.San_ResonanceCrystal.priestsHeros) > 0 {
		for _, v := range self.San_ResonanceCrystal.priestsHeros {
			equip := self.player.GetModule("equip").(*ModEquip).getAttr(v)
			_, ok := equip[99]
			if ok {
				fight += equip[99].AttValue
			}
		}
	}
	return fight
}

// 设置共鸣英雄
func (self *ModResonanceCrystal) SetResonanceHeros(body []byte) {
	var msg C2S_ResonanceCrystalSet
	json.Unmarshal(body, &msg)
	index, heroKey := msg.Index, msg.HeroKey
	// 判断格子数和数据结构
	if len(self.San_ResonanceCrystal.resonanceHeros) != self.San_ResonanceCrystal.ResonanceCount {
		self.player.SendErrInfo("err", "判断格子数和数据结构"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 设置的是否已经是共鸣英雄
	if self.IsResonanceHero(heroKey) {
		self.player.SendErrInfo("err", "设置的是否已经是共鸣英雄"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 判断index是否合法
	if index > self.San_ResonanceCrystal.ResonanceCount || index <= 0 {
		self.player.SendErrInfo("err", "判断index是否合法"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 判断格子是否已经被占用了
	if self.San_ResonanceCrystal.resonanceHeros[index-1].HeroKey != 0 {
		self.player.SendErrInfo("err", "判断格子是否已经被占用了"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 判断是否在cd中
	if self.San_ResonanceCrystal.resonanceHeros[index-1].EndTime > TimeServer().Unix() {
		self.player.SendErrInfo("err", "判断是否在cd中"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 获得英雄
	hero := self.player.getHero(heroKey)
	if nil == hero {
		self.player.SendErrInfo("err", "设置共鸣英雄获得英雄"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 虚空英雄参加了共鸣则返回
	if hero.VoidHero != 0 && hero.Resonance != 0 {
		self.player.SendErrInfo("err", "共鸣中的虚空英雄不能上阵"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	config := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if config == nil {
		self.player.SendErrInfo("err", "设置共鸣英雄配置"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	//if hero.HeroLv != 1 {
	//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
	//	return
	//}
	// 判断是不是已经是共鸣英雄
	if hero.UseType[HERO_USE_TYPE_CRYSTAL_RESONANCE] != 0 {
		self.player.SendErrInfo("err", "判断是不是已经是共鸣英雄"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 获得等级
	level := 0
	nLen := len(self.San_ResonanceCrystal.priestsHeros)
	if nLen == RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX {
		lastHero := self.player.getHero(self.San_ResonanceCrystal.priestsHeros[nLen-1])
		if lastHero != nil {
			level = lastHero.HeroLv
		}
	}
	// 判断等级
	if level <= 0 {
		self.player.SendErrInfo("err", "判断等级"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 设置状态
	hero.UseType[HERO_USE_TYPE_CRYSTAL_RESONANCE] = 1
	hero.OriginalLevel = hero.HeroLv

	nCount := 0
	if config.FinalLevel >= level {
		nCount = level - hero.HeroLv
	} else {
		nCount = config.FinalLevel - hero.HeroLv
	}
	hero.LvUp(nCount)

	self.player.countHeroFight(hero, ReasonHeroLvUp)

	self.San_ResonanceCrystal.resonanceHeros[index-1].HeroKey = heroKey
	self.San_ResonanceCrystal.resonanceHeros[index-1].EndTime = 0

	self.player.HandleTask(TASK_TYPE_RESONANCE_CRYSTAL_SET, 0, 0, 0)

	self.player.GetModule("hero").(*ModHero).ChangeVoidHeroResonance(hero.HeroKeyId, 0)

	// 发送信息同步
	var backmsg2 S2C_ResonanceCrystaUpdateResonance
	backmsg2.Cid = RESONANCE_CRYSTAL_UPDATE_RESONANCE
	backmsg2.Heros = append(backmsg2.Heros, hero)
	self.player.SendMsg(backmsg2.Cid, HF_JtoB(&backmsg2))
	// 返回设置消息
	var backmsg S2C_ResonanceCrystalSet
	backmsg.Cid = RESONANCE_CRYSTAL_SET
	backmsg.Index = index
	backmsg.HeroKey = heroKey
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_RESONANCE_CRYSTAL_SET, hero.HeroId, index, hero.OriginalLevel, "设置法阵列表英雄", 0, hero.HeroLv, self.player)

	self.UpdateMaxFightAll()
}

// 取消共鸣英雄
func (self *ModResonanceCrystal) CancelResonanceHerosMsg(body []byte) {
	var msg C2S_ResonanceCrystalCancel
	json.Unmarshal(body, &msg)
	heroKey, index := msg.HeroKey, self.GetResonanceIndex(msg.HeroKey)
	// 主动取消设置cd
	self.CancelResonanceHeros(index, heroKey, true)
}

// 取消共鸣英雄
func (self *ModResonanceCrystal) CancelResonanceHeros(index int, heroKey int, isCD bool) {
	// 是否是共鸣英雄
	if !self.IsResonanceHero(heroKey) {
		self.player.SendErrInfo("err", "取消共鸣英雄是否是共鸣英雄"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// index是否合法
	if index >= self.San_ResonanceCrystal.ResonanceCount || index < 0 {
		self.player.SendErrInfo("err", "取消共鸣英雄index是否合法"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 是否没被占用了
	if self.San_ResonanceCrystal.resonanceHeros[index].HeroKey == 0 {
		self.player.SendErrInfo("err", "取消共鸣英雄是否没被占用了"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	oldlevel := 0
	// 获得英雄
	hero := self.player.getHero(heroKey)
	if nil != hero {
		// 设置属性
		if hero.UseType[HERO_USE_TYPE_CRYSTAL_RESONANCE] != 1 {
			self.player.SendErrInfo("err", "取消共鸣英雄设置属性"+GetCsvMgr().GetText("STR_ERROR"))
			return
		}
		oldlevel = hero.HeroLv
		hero.UseType[HERO_USE_TYPE_CRYSTAL_RESONANCE] = 0
		nCount := hero.OriginalLevel - hero.HeroLv
		hero.LvUp(nCount)
	}
	// 设置状态
	self.San_ResonanceCrystal.resonanceHeros[index].HeroKey = 0
	// 是否添加cd
	if isCD {
		self.San_ResonanceCrystal.resonanceHeros[index].EndTime = TimeServer().Unix() + RESONANCE_CRYSTAL_RESONANCE_CD*DAY_SECS
		//self.San_ResonanceCrystal.resonanceHeros[index].EndTime = TimeServer().Unix() + MIN_SECS
	}
	heroID := 0
	herolv := 0
	if hero != nil {
		// 同步属性
		var backmsg2 S2C_ResonanceCrystaUpdateResonance
		backmsg2.Cid = RESONANCE_CRYSTAL_UPDATE_RESONANCE
		backmsg2.Heros = append(backmsg2.Heros, hero)
		self.player.SendMsg(backmsg2.Cid, HF_JtoB(&backmsg2))

		heroID = hero.HeroId
		herolv = hero.HeroLv

		self.player.GetModule("hero").(*ModHero).ChangeVoidHeroResonance(hero.HeroKeyId, 0)
	}

	// 取消回复
	var backmsg S2C_ResonanceCrystaCancel
	backmsg.Cid = RESONANCE_CRYSTAL_CANCEL
	backmsg.Index = index + 1
	backmsg.EndTime = self.San_ResonanceCrystal.resonanceHeros[index].EndTime
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_RESONANCE_CRYSTAL_CANCEL, heroID, index, oldlevel, "解除法阵列表英雄", 0, herolv, self.player)

	self.UpdateMaxFightAll()
}

// 发送信息
func (self *ModResonanceCrystal) SendInfo(body []byte) {
	self.CheckEndTime()
	self.CheckHeroExist()
	self.CheckLevel()
	var backmsg S2C_ResonanceCrystaInfo
	backmsg.Cid = RESONANCE_CRYSTAL_INFO
	backmsg.ResonanceCount = self.San_ResonanceCrystal.ResonanceCount
	backmsg.PriestsHeros = self.San_ResonanceCrystal.priestsHeros
	backmsg.ResonanceHeros = self.San_ResonanceCrystal.resonanceHeros
	backmsg.FightAll = self.San_ResonanceCrystal.MaxFightAll
	backmsg.Level = GetOfflineInfoMgr().GetMaxLevel(self.player.Sql_UserBase.Uid)
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 发送信息
func (self *ModResonanceCrystal) CleanCD(body []byte) {
	// 清理完成cd时间
	self.CheckEndTime()
	var msg C2S_ResonanceCrystalCleanCD
	json.Unmarshal(body, &msg)
	index := msg.Index
	// 判断index合法
	if index > self.San_ResonanceCrystal.ResonanceCount || index <= 0 {
		self.player.SendErrInfo("err", "清理完成cd时间判断index合法"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 如果没有cd则返回
	if self.San_ResonanceCrystal.resonanceHeros[index-1].EndTime == 0 {
		self.player.SendErrInfo("err", "清理完成cd时间没有cd则返回"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 获得消耗配置
	config := GetCsvMgr().GetTariffConfig2(TARIFF_RESONANCE_CRYSTAL_CLEAN_CD)
	if config == nil {
		self.player.SendErrInfo("err", "清理完成cd时间获得消耗配置"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 是否足够
	if err := self.player.HasObjectOk(config.ItemIds, config.ItemNums); err != nil {
		self.player.SendErrInfo("err", err.Error())
		return
	}
	// 回复消息
	var backmsg S2C_ResonanceCrystalCleanCD
	backmsg.Cid = RESONANCE_CRYSTAL_CLEAN_CD
	backmsg.Index = msg.Index
	backmsg.Items = self.player.RemoveObjectLst(config.ItemIds, config.ItemNums, "清理cd", 0, 0, 0)
	self.San_ResonanceCrystal.resonanceHeros[index-1].EndTime = 0
	backmsg.Data = self.San_ResonanceCrystal.resonanceHeros[index-1]
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

// 增加共鸣格子
func (self *ModResonanceCrystal) AddResonanceCount(body []byte) {
	var msg C2S_ResonanceCrystalAddResonanceCount
	json.Unmarshal(body, &msg)
	// 判断最大格子数
	if self.San_ResonanceCrystal.ResonanceCount >= RESONANCE_CRYSTAL_RESONANCE_COUNT_MAX {
		self.player.SendErrInfo("err", "增加共鸣格子判断最大格子数"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 获得消耗配置
	costID := 0
	/// 魔龙新
	if msg.Type == RESONANCE_CRYSTAL_ADD_GEM {
		costID = TARIFF_RESONANCE_CRYSTAL_RESONANCE_GEM
	} else {
		costID = TARIFF_RESONANCE_CRYSTAL_RESONANCE
	}
	config := GetCsvMgr().GetTariffConfig(costID, self.San_ResonanceCrystal.ResonanceCount+1)
	if config == nil {
		self.player.SendErrInfo("err", "增加共鸣格子配置"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}
	// 是否足够
	if err := self.player.HasObjectOk(config.ItemIds, config.ItemNums); err != nil {
		if msg.Type == RESONANCE_CRYSTAL_ADD_GEM {
			self.player.SendErrInfo("err", "钻石不足")
		} else {
			self.player.SendErrInfo("err", "命运精华不足")
		}
		return
	}
	// 扣除并返回消息
	var backmsg S2C_ResonanceCrystaAddResonanceCount
	backmsg.Cid = RESONANCE_CRYSTAL_ADD_RESONANCE_COUNT
	backmsg.Items = self.player.RemoveObjectLst(config.ItemIds, config.ItemNums, "解锁法阵格子", 0, 0, 0)
	self.San_ResonanceCrystal.ResonanceCount++
	self.San_ResonanceCrystal.resonanceHeros = append(self.San_ResonanceCrystal.resonanceHeros, &ResonanceHeros{})

	self.player.HandleTask(TASK_TYPE_RESONANCE_CRYSTAL_COUNT, self.San_ResonanceCrystal.ResonanceCount, 0, 0)
	backmsg.Type = msg.Type
	backmsg.ResonanceCount = self.San_ResonanceCrystal.ResonanceCount
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_RESONANCE_CRYSTAL_ADD, msg.Type, backmsg.ResonanceCount, backmsg.ResonanceCount-1, "解锁法阵格子", 0, 0, self.player)

}

// 能否升级
func (self *ModResonanceCrystal) CanUpLevel(heroKey int) bool {
	//// 如果是共鸣英雄则直接返回false
	//if self.IsResonanceHero(heroKey) {
	//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RESONANCE_CRYSTAL_IS_RESONANCE_HERO"))
	//	return false
	//}
	//// 如果是祭祀
	//if self.IsPriestsHero(heroKey) {
	//	hero := self.player.getHero(heroKey)
	//	if hero == nil {
	//		return false
	//	}
	//	config := GetCsvMgr().GetResonanceCrystalconfig(hero.HeroLv)
	//	if config == nil {
	//		return false
	//	}
	//	// 判断是否等级差超过配置
	//	needLevel := config.CrystalLevel
	//	for _, v := range self.San_ResonanceCrystal.priestsHeros {
	//		priestsHero := self.player.getHero(v)
	//		if nil == priestsHero {
	//			continue
	//		}
	//		if priestsHero.HeroLv < needLevel {
	//			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RESONANCE_CRYSTAL_PRIESTS_HERO_LEVEL_LIMIT"))
	//			return false
	//		}
	//	}
	//}
	return false
}

func (self *ModResonanceCrystal) CanUpLevelByAim(heroKey int, aimLv int) bool {
	// 如果是共鸣英雄则直接返回false
	if self.IsResonanceHero(heroKey) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RESONANCE_CRYSTAL_IS_RESONANCE_HERO"))
		return false
	}

	config := GetCsvMgr().GetResonanceCrystalconfig(aimLv)
	if config == nil {
		return false
	}

	// 升级到有需求等级的等级 祭司没有满 则直接返回
	needLevel := config.CrystalLevel
	len := len(self.San_ResonanceCrystal.priestsHeros)
	if needLevel > 1 && len != RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX {
		return false
	}

	if needLevel <= 1 {
		return true
	}

	lastheroKey := self.San_ResonanceCrystal.priestsHeros[len-1]
	lastHero := self.player.getHero(lastheroKey)
	if lastHero == nil {
		return false
	}

	// 如果是祭祀
	if self.IsPriestsHero(heroKey) {
		// 是最后一个祭司
		if heroKey == lastheroKey {
			// 如果目标等级从根本上改变了祭司英雄的最低等级 则判断升级后等级 及前四个祭司等级
			if aimLv < needLevel {
				return false
			}
			// 判断前四个
			for i, v := range self.San_ResonanceCrystal.priestsHeros {
				if i == len-1 {
					break
				}
				priestsHero := self.player.getHero(v)
				if nil == priestsHero {
					continue
				}
				if priestsHero.HeroLv < needLevel {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RESONANCE_CRYSTAL_PRIESTS_HERO_LEVEL_LIMIT"))
					return false
				}
			}
		} else {
			// 判断是否等级差超过配置
			for _, v := range self.San_ResonanceCrystal.priestsHeros {
				priestsHero := self.player.getHero(v)
				if nil == priestsHero {
					continue
				}
				if priestsHero.HeroLv < needLevel {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RESONANCE_CRYSTAL_PRIESTS_HERO_LEVEL_LIMIT"))
					return false
				}
			}
		}
	} else {
		// 如果目标等级超过了最低祭司的等级则只判断前四个和自身
		if aimLv > lastHero.HeroLv {
			if aimLv < needLevel {
				return false
			}
			// 判断前四个
			for i, v := range self.San_ResonanceCrystal.priestsHeros {
				if i == len-1 {
					break
				}
				priestsHero := self.player.getHero(v)
				if nil == priestsHero {
					continue
				}
				if priestsHero.HeroLv < needLevel {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RESONANCE_CRYSTAL_PRIESTS_HERO_LEVEL_LIMIT"))
					return false
				}
			}
		} else {
			// 判断是否等级差超过配置
			for _, v := range self.San_ResonanceCrystal.priestsHeros {
				priestsHero := self.player.getHero(v)
				if nil == priestsHero {
					continue
				}
				if priestsHero.HeroLv < needLevel {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RESONANCE_CRYSTAL_PRIESTS_HERO_LEVEL_LIMIT"))
					return false
				}
			}
		}
	}

	return true
}

// 是否是祭司英雄
func (self *ModResonanceCrystal) IsPriestsHero(heroKey int) bool {
	for _, v := range self.San_ResonanceCrystal.priestsHeros {
		if heroKey == v {
			return true
		}
	}
	return false
}

// 是否是共鸣英雄
func (self *ModResonanceCrystal) IsResonanceHero(heroKey int) bool {
	for _, v := range self.San_ResonanceCrystal.resonanceHeros {
		if v.HeroKey == 0 {
			continue
		}

		if heroKey == v.HeroKey {
			return true
		}
	}
	return false
}

func (self *ModResonanceCrystal) CheckMaxFight() {
	maxFight := int64(0)
	maxLevel := 0
	for i, v := range self.San_ResonanceCrystal.priestsHeros {
		hero := self.player.getHero(v)
		if nil != hero {
			maxFight += hero.Fight

			if i < 3 {
				maxLevel += hero.HeroLv
			}
		}
	}

	if maxFight > self.San_ResonanceCrystal.MaxFight {
		self.San_ResonanceCrystal.MaxFight = maxFight
		self.San_ResonanceCrystal.MaxFightTime = TimeServer().Unix()
		GetOfflineInfoMgr().SetMaxFight(self.player.Sql_UserBase.Uid, maxFight)
	}

	GetOfflineInfoMgr().SetNewHeroLv(self.player.GetUid(), maxLevel/3)
}

//返回祭司的专属总等级，和开启专属人数
func (self *ModResonanceCrystal) CalExclusiveLvMax() (int, int) {
	maxLv := 0
	maxNum := 0
	for _, v := range self.San_ResonanceCrystal.priestsHeros {
		hero := self.player.getHero(v)
		if nil != hero && hero.ExclusiveEquip != nil {
			if hero.ExclusiveEquip.UnLock == LOGIC_TRUE {
				maxLv += hero.ExclusiveEquip.Lv
				maxNum += 1
			}
		}
	}

	return maxLv, maxNum
}

// 检测设置等级
func (self *ModResonanceCrystal) CheckLevel() {
	// 祭司个数不够
	if len(self.San_ResonanceCrystal.priestsHeros) != RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX {
		return
	}
	// 获得最低等级
	level := 0
	lastHero := self.player.getHero(self.San_ResonanceCrystal.priestsHeros[len(self.San_ResonanceCrystal.priestsHeros)-1])
	if lastHero != nil {
		level = lastHero.HeroLv
	}
	if level < 1 {
		return
	}
	// 设置英雄等级
	update := []*Hero{}
	for _, v := range self.San_ResonanceCrystal.resonanceHeros {
		hero := self.player.getHero(v.HeroKey)
		if nil != hero {
			config := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
			if hero.HeroLv != level && config != nil {
				nCount := 0
				if config.FinalLevel >= level {
					nCount = level - hero.HeroLv
				} else {
					nCount = config.FinalLevel - hero.HeroLv
				}
				hero.LvUp(nCount)
				self.player.countHeroFight(hero, ReasonHeroLvUp)
				update = append(update, hero)

				// 非虚空英雄 有共鸣的虚空英雄
				if hero.VoidHero == 0 && hero.Resonance != 0 {
					self.player.GetModule("hero").(*ModHero).ChangeVoidHeroResonance(hero.HeroKeyId, 0)
				}
			}
		}
	}

	self.player.HandleTask(TASK_TYPE_RESONANCE_CRYSTAL_LEVEL, level, 0, 0)
	if len(update) > 0 {
		var backmsg S2C_ResonanceCrystaUpdateResonance
		backmsg.Cid = RESONANCE_CRYSTAL_UPDATE_RESONANCE
		backmsg.Heros = update
		self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
	}

	GetOfflineInfoMgr().UpdateMaxLevel(self.player, level)
	return
}

//得到祭祀平均等级，地牢用  返回 平均等级，最大星级，平均星级
func (self *ModResonanceCrystal) GetPriestsPitLv() (int, int, int) {
	level := 1
	starMax := 2
	starAll := 0
	count := 0
	if len(self.San_ResonanceCrystal.priestsHeros) != RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX {
		return level, starMax, 2
	}
	for i := 0; i < len(self.San_ResonanceCrystal.priestsHeros); i++ {
		hero := self.player.getHero(self.San_ResonanceCrystal.priestsHeros[i])
		if nil != hero {
			level += hero.HeroLv
			if hero.StarItem == nil {
				continue
			}
			starAll += hero.StarItem.UpStar
			count++
			if hero.StarItem.UpStar > starMax {
				starMax = hero.StarItem.UpStar
			}
		}
	}

	if level > RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX {
		level = level / RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX
	}
	if count > 0 {
		starAll = starAll / count
	}
	return level, starMax, starAll
}

// 升级 分解 重置英雄调用
func (self *ModResonanceCrystal) UpdatePriestsHeros(heroKey int) {
	// 是否是祭司英雄
	if self.IsPriestsHero(heroKey) {
		// 重新判断
		self.SetPriestsHeros()
	} else {
		// 不是祭司则判断
		hero := self.player.getHero(heroKey)
		// 英雄还存在
		if hero != nil {
			// 祭司满了
			if len(self.San_ResonanceCrystal.priestsHeros) >= RESONANCE_CRYSTAL_PRIESTS_COUNT_MAX {
				// 判断是否比最低祭司等级高 是就重新设置
				lastHero := self.player.getHero(self.San_ResonanceCrystal.priestsHeros[len(self.San_ResonanceCrystal.priestsHeros)-1])
				if nil != lastHero {
					if hero.HeroLv >= lastHero.HeroLv {
						self.SetPriestsHeros()
					}
				}
			} else {
				// 祭司没满
				self.SetPriestsHeros()
			}
		} else {
			// 是共鸣英雄
			if self.IsResonanceHero(heroKey) {
				index := self.GetResonanceIndex(heroKey)
				if index >= 0 {
					// 取消共鸣英雄
					self.CancelResonanceHeros(index, heroKey, false)
				}
			}
		}
	}

	self.UpdateMaxFightAll()
}

func (self *ModResonanceCrystal) UpdateMaxFightAll() {
	fight := int64(0)
	for _, v := range self.San_ResonanceCrystal.priestsHeros {
		hero := self.player.getHero(v)
		if nil != hero {
			// 恢复使用类型
			fight += hero.Fight
		}
	}

	for _, v := range self.San_ResonanceCrystal.resonanceHeros {
		hero := self.player.getHero(v.HeroKey)
		if nil != hero {
			fight += hero.Fight
		}
	}

	if self.San_ResonanceCrystal.MaxFightAll != fight {
		self.San_ResonanceCrystal.MaxFightAll = fight

		self.player.HandleTask(TASK_TYPE_CRYSTAL_FIGHT, int(self.San_ResonanceCrystal.MaxFightAll)/100, 0, 0)

		var backmsg S2C_ResonanceCrystalFight
		backmsg.Cid = RESONANCE_CRYSTAL_FIGHT
		backmsg.FightAll = self.San_ResonanceCrystal.MaxFightAll
		self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
	}

}

// 获得共鸣index
func (self *ModResonanceCrystal) GetResonanceIndex(heroKey int) int {
	for i, v := range self.San_ResonanceCrystal.resonanceHeros {
		if v.HeroKey == 0 {
			continue
		}

		if heroKey == v.HeroKey {
			return i
		}
	}
	return -1
}

//func (self *ModResonanceCrystal) GetResonanceHero(index int) *ResonanceHeros {
//	for _, v := range self.San_ResonanceCrystal.resonanceHeros {
//		if index == v.Index {
//			return v
//		}
//	}
//	return nil
//}

// 检查完成的cd时间
func (self *ModResonanceCrystal) CheckEndTime() {
	now := TimeServer().Unix()
	for _, v := range self.San_ResonanceCrystal.resonanceHeros {
		if now >= v.EndTime {
			v.EndTime = 0
		}
	}
	return
}

// 检查英雄存在
func (self *ModResonanceCrystal) CheckHeroExist() {
	for _, v := range self.San_ResonanceCrystal.priestsHeros {
		if v == 0 {
			continue
		}
		hero := self.player.getHero(v)
		if hero == nil {
			self.SetPriestsHeros()
			break
		}
	}

	for i, v := range self.San_ResonanceCrystal.resonanceHeros {
		if v.HeroKey == 0 {
			continue
		}

		hero := self.player.getHero(v.HeroKey)
		if hero == nil {
			self.CancelResonanceHeros(i, v.HeroKey, false)
		}
	}
	return
}

// 获得所有共鸣水晶系统中的英雄KEYID
func (self *ModResonanceCrystal) GetHeroInThis() []int {
	keyids := make([]int, 0)
	for _, v := range self.San_ResonanceCrystal.priestsHeros {
		if v == 0 {
			continue
		}
		keyids = append(keyids, v)
	}

	for _, v := range self.San_ResonanceCrystal.resonanceHeros {
		if v.HeroKey == 0 {
			continue
		}
		keyids = append(keyids, v.HeroKey)
	}
	return keyids
}

func (self *ModResonanceCrystal) CheckTask() {
	self.player.HandleTask(TASK_TYPE_CRYSTAL_FIGHT, int(self.San_ResonanceCrystal.MaxFightAll)/100, 0, 0)

	level := 0
	if len(self.San_ResonanceCrystal.priestsHeros) > 0 {
		lastHero := self.player.getHero(self.San_ResonanceCrystal.priestsHeros[len(self.San_ResonanceCrystal.priestsHeros)-1])
		if lastHero != nil {
			level = lastHero.HeroLv
		}
	}
	self.player.HandleTask(TASK_TYPE_RESONANCE_CRYSTAL_LEVEL, level, 0, 0)

	self.player.HandleTask(TASK_TYPE_RESONANCE_CRYSTAL_SET, 0, 0, 0)
}
