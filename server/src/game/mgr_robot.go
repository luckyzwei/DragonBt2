package game

import (
	"fmt"
	"sort"
	//"time"
)

// 战斗武将信息
type FightHero struct {
	Heroid     int   `json:"id"`         // 武将id
	Hp         int   `json:"hp"`         // hp
	Energy     int   `json:"rage"`       // 怒气
	Damage     int64 `json:"damage"`     //! 伤害
	TakeDamage int64 `json:"takedamage"` //! 承受伤害
	Healing    int64 `json:"healing"`    //! 治疗
}

const (
	SKILL_LEVEL = 1
)

// 兵法: 已废弃
type JS_ArmsSkill struct {
	Id    int `json:"id"`    // 兵法Id
	Level int `json:"level"` // 等级
}

type TeamCal struct {
	Camp int
	Num  int
}

type TeamCalList struct {
	list []TeamCal
}

// 天赋
type Js_TalentSkill struct {
	SkillId int `json:"skillid"` // 技能Id
	Level   int `json:"level"`   // 技能等级
}
type FriendHeroUserInfo struct {
	Uid       int64  `json:"uid"`    // id
	Uname     string `json:"uname"`  // 名字
	UnionName string `json:"union"`  //! 军团名字
	Iconid    int    `json:"iconid"` // icon
	Camp      int    `json:"camp"`   // 阵营
	Level     int    `json:"level"`  // 等级
	Vip       int    `json:"vip"`    // Vip 等级
	HeroKey   []int  `json:"herokey"`
}

// 战斗信息
type JS_FightInfo struct {
	Rankid       int              `json:"rankid"`       // 类型
	Uid          int64            `json:"uid"`          // id
	Uname        string           `json:"uname"`        // 名字
	UnionName    string           `json:"union"`        //! 军团名字
	Iconid       int              `json:"iconid"`       // icon
	Camp         int              `json:"camp"`         // 阵营
	Level        int              `json:"level"`        // 等级
	Vip          int              `json:"vip"`          // Vip 等级
	Defhero      []int            `json:"defhero"`      // 出战武将
	Heroinfo     []JS_HeroInfo    `json:"heroinfo"`     // 各武将信息
	HeroParam    []JS_HeroParam   `json:"heroparam"`    // 各武将属性
	Deffight     int64            `json:"deffight"`     // 玩家总战力
	FightTeam    int              `json:"fightteam"`    // 出战兵营
	FightTeamPos TeamPos          `json:"fightteampos"` // 出战兵营
	Portrait     int              `json:"portrait"`     // 边框  20190412 by zy
	LifeTreeInfo *JS_LifeTreeInfo `json:"lifetreeinfo"`
}

// 巨兽信息
type JsBoss struct {
	HeroId    int            `json:"heroid"`
	HeroParam *JS_HeroParam  `json:"heroinfo"`
	ArmsSkill []JS_ArmsSkill `json:"armsskill"`
}

// 鼓舞信息
type JsBuff struct {
	Lv        int           `json:"lv"`
	HeroParam *JS_HeroParam `json:"heroinfo"`
}

// 获取血量万分比
func (self *JS_FightInfo) GetHp() int {
	cur := float64(0)
	max := float64(0)
	for i := 0; i < len(self.HeroParam); i++ {
		cur += self.HeroParam[i].Hp
		max += self.HeroParam[i].Param[6]
	}

	return int(cur * 10000 / max)
}

// 加Buff
func (self *JS_FightInfo) AddBuffer(atttype int, value float64) {
	for i := 0; i < len(self.HeroParam); i++ {
		if atttype >= 100 {
			index := atttype - 100
			if index >= 0 && index < 33 {
				self.HeroParam[i].Param[index] = self.HeroParam[i].Param[index] * (1.0 + value/10000.0)
			}
		} else {
			if atttype >= 0 && atttype < 33 {
				self.HeroParam[i].Param[atttype] += value
			}
		}

	}
}

// 满血复活
func (self *JS_FightInfo) FullHp() {
	for i := 0; i < len(self.HeroParam); i++ {
		self.HeroParam[i].Hp = self.HeroParam[i].Param[6]
	}
}

// 设置战斗信息
func (self *JS_FightInfo) SetFightInfo(lst []FightHero) int {
	if len(lst) == 0 {
		return 0
	}

	killlevel := 0
	for i := 0; i < len(self.Defhero); {
		find := false
		for j := 0; j < len(lst); j++ {
			if self.Defhero[i] == lst[j].Heroid && lst[j].Hp != 0 {
				find = true
				self.HeroParam[i].Hp = self.HeroParam[i].Param[6] * float64(lst[j].Hp) / 10000
				self.HeroParam[i].Energy = lst[j].Energy
				//LogDebug(lst[j].HeroId, "最大血量:", self.HeroParam[i].Param[6])
				//LogDebug(lst[j].HeroId, "当前血量:", self.HeroParam[i].Hp)
				break
			}
		}

		if find {
			i++
			continue
		}

		if i < len(self.Heroinfo) {
			killlevel += self.Heroinfo[i].Levels + 2
			self.Defhero = append(self.Defhero[i:], self.Defhero[i+1:]...)
			self.Heroinfo = append(self.Heroinfo[i:], self.Heroinfo[i+1:]...)
			self.HeroParam = append(self.HeroParam[i:], self.HeroParam[i+1:]...)
		}
	}

	return killlevel
}

// 设置战斗信息不清除式
func (self *JS_FightInfo) SetFightInfoNoClear(lst []FightHero) int {
	if len(lst) == 0 {
		return 0
	}

	for j := 0; j < len(lst); j++ {
		LogDebug("1st-ID:", lst[j].Heroid, "|1st-Hp:", lst[j].Hp)
	}

	for j := 0; j < len(self.Defhero); j++ {
		LogDebug("self-ID:", self.Defhero[j])
	}

	bossDis := 0 //巨兽的位置不确定，会造成偏移
	for i := 0; i < len(self.Defhero); i++ {
		//巨兽后面处理 ZY 20190904
		if self.Defhero[i] > 5000 {
			//!巨兽
			/*
				bossDis = 1
				if self.BossInfo != nil {
					for j := 0; j < len(lst); j++ {
						if self.BossInfo.HeroId == lst[j].Heroid && lst[j].Hp != 0 {
							self.BossInfo.HeroParam.Hp = self.BossInfo.HeroParam.Param[6] * float32(lst[j].Hp) / 10000
							self.BossInfo.HeroParam.Energy = lst[j].Energy
							LogDebug(lst[j].Heroid, "巨兽血量：", self.BossInfo.HeroParam.Param[6], "当前血量：", self.BossInfo.HeroParam.Hp)
							break
						}
					}
				}
			*/
			continue
		}

		for j := 0; j < len(lst); j++ {
			if self.Defhero[i] == lst[j].Heroid {
				tempI := i - bossDis
				if tempI < len(self.HeroParam) && len(self.HeroParam[tempI].Param) > 6 {
					if lst[j].Hp != 0 {
						self.HeroParam[tempI].Hp = self.HeroParam[tempI].Param[6] * float64(lst[j].Hp) / 10000
					} else {
						self.HeroParam[tempI].Hp = 0
					}
					self.HeroParam[tempI].Energy = lst[j].Energy
					LogDebug(lst[j].Heroid, "最大血量:", self.HeroParam[tempI].Param[6], "当前血量:", self.HeroParam[tempI].Hp)
				}
				break
			}
		}
	}

	for j := 0; j < len(self.HeroParam); j++ {
		LogDebug("self.HeroParam-ID:", self.HeroParam[j].Heroid, "self.HeroParam-HP:", self.HeroParam[j].Hp)
	}

	LogDebug("战斗数据同步完毕")

	return 0
}

func (self *JS_FightInfo) GetBastHero() int { // 得到战力最高的
	var fight int64 = -1
	id := 0

	for i := 0; i < len(self.Heroinfo); i++ {
		if self.Heroinfo[i].Fight > fight {
			fight = self.Heroinfo[i].Fight
			id = self.Defhero[i]
		}
	}

	return id
}

// 得到简化版结构
func (self *JS_FightInfo) GetFightBase() *JS_FightBase {
	node := new(JS_FightBase)
	node.Rankid = self.Rankid
	node.Uid = self.Uid
	node.Uname = self.Uname
	node.Iconid = self.Iconid
	node.Portrait = self.Portrait
	node.Camp = self.Camp
	node.Level = self.Level
	node.Vip = self.Vip
	node.Deffight = self.Deffight
	//if self.BossInfo != nil {
	//	node.BossInfo = self.BossInfo
	//}

	for i := 0; i < len(self.Heroinfo); i++ {
		var herobase JS_HeroBase
		herobase.Heroid = self.Heroinfo[i].Heroid
		herobase.Stars = self.Heroinfo[i].Stars
		herobase.MainTalent = self.Heroinfo[i].MainTalent
		herobase.Color = self.Heroinfo[i].Color
		herobase.Skin = self.Heroinfo[i].Skin
		herobase.Levels = self.Heroinfo[i].Levels
		herobase.Fight = self.Heroinfo[i].Fight
		node.Heroinfo = append(node.Heroinfo, herobase)
	}

	return node
}

// 简化版
type JS_FightBase struct {
	Rankid    int           `json:"rankid"`
	Uid       int64         `json:"uid"`       // id
	Uname     string        `json:"uname"`     // 名字
	Iconid    int           `json:"iconid"`    // icon
	Portrait  int           `json:"portrait"`  // 边框  20190412 by zy
	Camp      int           `json:"camp"`      // 阵营
	Level     int           `json:"level"`     // 等级
	Vip       int           `json:"vip"`       // Vip
	Heroinfo  []JS_HeroBase `json:"heroinfo"`  // 各武将信息
	Deffight  int64         `json:"deffight"`  // 战力
	BossInfo  *JsBoss       `json:"bossinfo"`  // 巨兽信息
	Encourage int           `json:"encourage"` // 鼓舞
}

type JS_HeroBase struct {
	Heroid     int   `json:"heroid"` // 武将id
	Color      int   `json:"color"`  //
	Stars      int   `json:"stars"`
	Levels     int   `json:"levels"`
	Skin       int   `json:"skin"`
	Fight      int64 `json:"fight"`      // 战斗力
	MainTalent int   `json:"maintalent"` // 主天赋
}

type JS_HeroInfo struct {
	Heroid          int              `json:"heroid"`    // 武将id
	HeroKeyId       int              `json:"herokeyid"` // 武将keyid
	Color           int              `json:"color"`     //
	Stars           int              `json:"stars"`
	Levels          int              `json:"levels"`
	Soldiercolor    int              `json:"soldiercolor"` // 士兵星级
	Soldierid       int              `json:"soldierid"`    // 士兵id
	Skilllevel1     int              `json:"skilllevel1"`  // 技能等级
	Skilllevel2     int              `json:"skilllevel2"`
	Skilllevel3     int              `json:"skilllevel3"`
	Skilllevel4     int              `json:"skilllevel4"`
	Skilllevel5     int              `json:"skilllevel5"`
	Skilllevel6     int              `json:"skilllevel6"`
	Fervor1         int              `json:"fervor1"` // 战魂
	Fervor2         int              `json:"fervor2"`
	Fervor3         int              `json:"fervor3"`
	Fervor4         int              `json:"fervor4"`
	Fight           int64            `json:"fight"` // 战斗力
	ArmsSkill       []JS_ArmsSkill   `json:"armsskill"`
	TalentSkill     []Js_TalentSkill `json:"talentskill"` // 天赋技能
	ArmyId          int              `json:"army_id"`
	MainTalent      int              `json:"maintalent"`      // 英雄主天赋等级
	HeroQuality     int              `json:"heroquality"`     //英雄品质  (不变的属性)
	HeroArtifactId  int              `json:"heroartifactid"`  //英雄神器等级 (不变的属性)
	HeroArtifactLv  int              `json:"heroartifactlv"`  //英雄神器等级 (不变的属性)
	HeroExclusiveLv int              `json:"heroexclusivelv"` //英雄专属等级 (不变的属性)
	Skin            int              `json:"skin"`            //皮肤
	ExclusiveUnLock int              `json:"exclusiveunlock"` //专属等级是否解锁
	IsArmy          int              `json:"isarmy"`          //是否是佣兵
}

func (self *JS_HeroInfo) ProcAtt(valuetype int, value float64, att []float64, att_per map[int]float64, att_ext map[int]float64) int {

	if valuetype < 0 || value == 0 {
		return 0
	}

	if valuetype == 2500 {
		return int(value)
	} else if valuetype >= 100 {
		att_per[valuetype] += value
	} else if valuetype >= len(att) {
		att_ext[valuetype] += value
	} else {
		att[valuetype] += value
	}

	return 0
}

func (self *JS_HeroInfo) GetAtt(ewvaluetype []int, ewvalue []float64) ([]float64, []JS_HeroExtAttr, int) {
	energy := 0
	att := make([]float64, 32)
	att_per := make(map[int]float64)
	att_ext := make(map[int]float64)

	// 额外属性加成
	for i := 0; i < len(ewvaluetype); i++ {
		if ewvaluetype[i] < 0 {
			continue
		}
		energy += self.ProcAtt(ewvaluetype[i], ewvalue[i], att, att_per, att_ext)
	}

	for i := 0; i < AttrEnd; i++ {
		value, ok := att_per[i+AttrDisExt+1]
		if !ok {
			continue
		}
		att[i] = att[i] * (1.0 + value/10000.0)
	}

	att_ext_ret := []JS_HeroExtAttr{}
	for k := range att_ext {
		value, ok := att_per[k+AttrDisExt]
		if ok {
			att_ext[k] += att_ext[k] * (1.0 + value/10000.0)
		}

		att_ext_ret = append(att_ext_ret, JS_HeroExtAttr{k, att_ext[k]})
	}

	return att, att_ext_ret, energy
}

func (self *JS_HeroInfo) CountFight(ewvaluetype []int, ewvalue []float64) ([]float64, []JS_HeroExtAttr, int) {
	att, att_ext, energy := self.GetAtt(ewvaluetype, ewvalue)

	for i := 0; i < len(att_ext); i++ {
		if att_ext[i].Key == AttrFight {
			self.Fight = int64(att_ext[i].Value)
			break
		}
	}
	return att, att_ext, energy
}

type JS_HeroExtAttr struct {
	Key   int     `json:"k"` // 属性类型
	Value float64 `json:"v"` // 属性值
}

type JS_HeroParam struct {
	Heroid  int              `json:"heroid"` // 武将id
	Param   []float64        `json:"param"`  // 武将属性
	Hp      float64          `json:"hp"`     // 当前hp
	Energy  int              `json:"energy"` // 当前怒气
	Pos     int              `json:"pos"`
	ExtAttr []JS_HeroExtAttr `json:"ext"` //扩展属性
}

type RobotMgr struct {
}

var robotmgrsingleton *RobotMgr = nil

// public
func GetRobotMgr() *RobotMgr {
	if robotmgrsingleton == nil {
		robotmgrsingleton = new(RobotMgr)
	}

	return robotmgrsingleton
}

// 将玩家变成标准战斗结构, 这个是最新改装过的结构
func (self *RobotMgr) GetPlayerFightInfo(player *Player, inParam int, encourageType int) *JS_FightInfo {
	var pos = player.getTeamPos()

	data := new(JS_FightInfo)
	data.Rankid = inParam
	data.Uid = player.Sql_UserBase.Uid
	data.Uname = player.Sql_UserBase.UName
	data.Iconid = player.Sql_UserBase.IconId
	data.Portrait = player.Sql_UserBase.Portrait //边框  20190412 by zy
	data.Level = player.Sql_UserBase.Level
	data.Vip = player.Sql_UserBase.Vip
	data.UnionName = player.GetUnionName()
	data.Camp = player.Sql_UserBase.Camp
	if pos == nil {
		LogError("pos is nil..")
		return nil
	}

	data.FightTeamPos = *pos
	for i := 0; i < len(pos.FightPos); i++ {
		heroid := pos.FightPos[i]
		if heroid == 0 {
			continue
		}

		data.Defhero = append(data.Defhero, heroid)
	}
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)
	teamAttr := self.GetTeamAttr(data.Defhero)
	for _, key := range data.Defhero {
		heroData := player.getHero(key)
		if heroData == nil {
			continue
		}
		var heroinfo JS_HeroInfo
		heroinfo.Heroid = heroData.HeroId
		heroinfo.Color = 1
		heroinfo.Stars = 1
		if heroData.StarItem != nil {
			heroinfo.Stars = heroData.StarItem.UpStar
		}

		heroinfo.HeroKeyId = heroData.HeroKeyId
		heroinfo.Levels = heroData.HeroLv
		heroinfo.Skin = heroData.Skin
		heroinfo.Skilllevel1 = 0
		heroinfo.Skilllevel2 = 0
		heroinfo.Skilllevel3 = 0
		heroinfo.Skilllevel4 = SKILL_LEVEL
		heroinfo.Skilllevel5 = 0
		heroinfo.Skilllevel6 = 0
		heroinfo.Fervor1 = 0
		heroinfo.Fervor2 = 0
		heroinfo.Fervor3 = 0
		heroinfo.Fervor4 = 0
		heroinfo.ArmsSkill = make([]JS_ArmsSkill, 0)
		heroinfo.TalentSkill = []Js_TalentSkill{}

		//战斗结构中处理技能 20190522 by zy
		//升星技能
		for i := 0; i < len(heroData.StarItem.Skills)-1; i++ {
			if heroData.StarItem.Skills[i] == 0 {
				break
			}
			heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: heroData.StarItem.Skills[i], Level: SKILL_LEVEL})
		}

		if heroData.StageTalent != nil {
			// 天赋技能
			for i := 0; i < len(heroData.StageTalent.AllSkill); i++ {
				config := GetCsvMgr().GetStageTalent(heroData.StageTalent.AllSkill[i].ID)
				if nil == config {
					continue
				}
				if heroData.StageTalent.AllSkill[i].Pos <= 0 || heroData.StageTalent.AllSkill[i].Pos > len(config.Skill) {
					continue
				}
				heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: config.Skill[heroData.StageTalent.AllSkill[i].Pos-1], Level: 1})
			}
		}

		var param JS_HeroParam
		param.Heroid = heroinfo.Heroid
		pAttrWrapper := heroData.GetAttr2(player, encourageType, teamAttr)
		param.Param, param.ExtAttr, param.Energy = pAttrWrapper.Base, pAttrWrapper.ExtRet, pAttrWrapper.Energy
		param.Hp = param.Param[6]
		heroinfo.Fight = heroData.Fight
		data.Deffight += heroData.Fight
		data.Heroinfo = append(data.Heroinfo, heroinfo)
		data.HeroParam = append(data.HeroParam, param)
	}

	return data
}

//! 获取巨兽信息，根据等级生成
func (self *RobotMgr) GetBossInfoByLevel(bossId int, level int) *JsBoss {
	//LogDebug("gen boss info:", bossId, level)
	bossParam := &JS_HeroParam{}

	//heroId := self.player.GetModule("team").(*ModTeam).getBossId()
	pAttr := &AttrWrapper{
		Base:     make([]float64, 32),
		Ext:      make(map[int]float64),
		Per:      make(map[int]float64),
		Energy:   0,
		FightNum: 0,
	}

	if bossId == 0 {
		return nil
	}

	attMap := make(map[int]*Attribute)
	config := GetCsvMgr().GetBossAtt(bossId, level)
	if config == nil {
		return &JsBoss{HeroId: bossId, HeroParam: bossParam, ArmsSkill: nil}
	}
	AddAttrDirect(attMap, config.Basetypes, config.Basevalues)

	ProcAtt(attMap, pAttr)
	ProcLast(level, pAttr)
	pAttr.ExtRet = ProcExtAttr(pAttr)
	//return pAttr, heroId

	//pAttrWrapper, bossId := player.GetModule("boss").(*ModBoss).GetAttr()
	bossParam.Param, bossParam.ExtAttr, bossParam.Energy = pAttr.Base, pAttr.ExtRet, pAttr.Energy
	bossParam.Hp = bossParam.Param[6]

	//增加巨兽技能  20190701 by zy
	bossSkill := make([]JS_ArmsSkill, 0)
	configBoss := GetCsvMgr().GetBossConfig(bossId - 5000)
	if configBoss != nil {
		for _, v := range configBoss.Skills {
			if v > 0 {
				bossSkill = append(bossSkill, JS_ArmsSkill{Id: v, Level: 1})
			}
		}

		if configBoss.Openskill != 0 {
			bossSkill = append(bossSkill, JS_ArmsSkill{Id: configBoss.Openskill, Level: 1})
		}
	}

	return &JsBoss{HeroId: bossId, HeroParam: bossParam, ArmsSkill: bossSkill}
}

// 生成城防军(国战使用)
func (self *RobotMgr) GetDefNpc(name string, id int, level int, camp int) *JS_FightInfo {
	if level < 10 {
		level = GetServer().GetLevel(false)
	}
	LogDebug("Get Def Npc :", name, id, level, camp)
	csvlevel, _ := GetCsvMgr().WorldLevel[level-9]
	csvtype, _ := GetCsvMgr().Data["World_Lvtpye"][id]

	data := new(JS_FightInfo)
	data.Rankid = FIGHTTYPE_DEF
	data.Uid = int64(-id)
	data.Uname = name
	data.Iconid = HF_Atoi(csvtype["npcicon"])
	data.FightTeam = HF_Atoi(csvtype["npcteam"])
	data.Camp = camp
	data.Level = level
	data.Defhero = []int{HF_Atoi(csvtype["heronpc1"]), HF_Atoi(csvtype["heronpc2"]), HF_Atoi(csvtype["heronpc3"]),
		HF_Atoi(csvtype["heronpc4"]), HF_Atoi(csvtype["heronpc5"])}
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)
	data.Deffight = 0

	for i, key := range data.Defhero {
		//ishero := false
		herostar := 1
		if key >= 1000 && key <= 9999 {
			//	ishero = true
			itemheroid := 11000000 + (key * 100) + 1
			config := GetCsvMgr().GetItemConfig(itemheroid)
			if config != nil {
				herostar = config.Special
			}
		}
		var hero JS_HeroInfo
		//if ishero {
		//	hero.Heroid = key*10 + herostar
		//} else {
		hero.Heroid = key
		//}

		data.FightTeamPos.addFightPos(key)

		hero.Color = 1
		hero.Stars = herostar
		hero.Levels = data.Level + HF_Atoi(csvtype["npclvtable"])
		hero.Skin = 0
		hero.Soldiercolor = 6
		hero.Skilllevel4 = SKILL_LEVEL
		hero.Skilllevel5 = 0
		hero.Skilllevel6 = 0
		hero.Fervor1 = 0
		hero.Fervor2 = 0
		hero.Fervor3 = 0
		hero.Fervor4 = 0
		hero.ArmsSkill = make([]JS_ArmsSkill, 0)
		hero.TalentSkill = make([]Js_TalentSkill, 0)
		hero.MainTalent = 0

		hero.Soldierid = 0

		fightnum := float64(0)
		att := make([]float64, 32)
		//csv_attribute :=
		for i := 0; i < len(csvlevel.BaseTypes); i++ {
			if csvlevel.BaseTypes[i] == 99 {
				fightnum = float64(csvlevel.BaseValues[i])
				continue
			}

			value, ok := GetCsvMgr().HeroAttribute[csvlevel.BaseTypes[i]]
			//value, ok := csvlevel[csv_attribute[i]["name"]]
			if !ok {
				continue
			}

			valuecorrect, _ := csvtype[fmt.Sprintf("%scorrect", value.Name)]
			att[i] = float64(csvlevel.BaseValues[i]) * (10000 + HF_Atof64(valuecorrect)) / 10000
		}

		for i := 0; i < 32; i++ {
			value, ok := GetCsvMgr().HeroAttribute[i]
			//value, ok := csvlevel[csv_attribute[i]["name"]]
			if !ok {
				continue
			}
			valuecorrect, _ := csvtype[fmt.Sprintf("%scorrect", value.Name)]
			att[i] = 100 * (10000 + HF_Atof64(valuecorrect)) / 10000
		}

		//fightnum := float32(0)
		//for i := 6; i < len(csv_attribute); i++ {
		//	value, ok := csv_fightingnum[csv_attribute[i]["name"]]
		//	if !ok {
		//		continue
		//	}
		//	xs := HF_Atof(value)
		//	fightnum += xs * att[i]
		//}
		hero.Fight = int64(fightnum)
		data.Heroinfo = append(data.Heroinfo, hero)

		var param JS_HeroParam
		param.Heroid = hero.Heroid
		param.Param = att
		param.Hp = param.Param[6]
		param.Energy = 0
		param.Pos = i
		data.HeroParam = append(data.HeroParam, param)

		data.Deffight += hero.Fight
	}

	//data.BossInfo = self.GetBossInfoByLevel(HF_Atoi(csvtype["behemoth"]), level)
	//if data.BossInfo != nil {
	//	data.FightTeamPos.addFightPos(HF_Atoi(csvtype["behemoth"]))
	//	data.Defhero = append(data.Defhero, data.BossInfo.HeroId)
	//}

	return data
}

func (self *RobotMgr) GetMonsterConfig(monsterid int) *LevelMonsterConfig {
	info, ok := GetCsvMgr().LevelMonsterMap[monsterid]
	if !ok {
		LogError("monsterId:", monsterid, " not exist!")
		return nil
	}
	return info
}

// 获取机器人信息
func (self *CsvMgr) GetRobot(cfg *JJCRobotConfig) *JS_FightInfo {
	data := new(JS_FightInfo)
	data.Rankid = 0
	data.Uid = 0
	if cfg.Name == "0" || cfg.Name == "" {
		data.Uname = GetCsvMgr().GetName()
	} else {
		data.Uname = cfg.Name
	}

	//data.Iconid = cfg.Head
	data.Portrait = 1000 //机器人边框  20190412 by zy
	data.Level = cfg.Level
	data.Defhero = make([]int, 0)
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)

	if cfg.Category == 0 {
		data.Camp = HF_GetRandom(3) + 1
	} else {
		data.Camp = cfg.Category
	}

	var heroes []int
	for i := 0; i < len(cfg.Hero); i++ {
		if cfg.Hero[i] == 0 {
			continue
		}

		heroes = append(heroes, cfg.Hero[i])
	}

	heroes = HF_GetRandomArr(heroes, 5)

	//20190422 by zy
	//HF_GetRandomArr参数是指针传递，函数内部有更新数组操作，导致数组元素丢失
	/*
		headArr := HF_GetRandomArr(heroes, 1)
		if len(headArr) > 0 {
			data.Iconid = headArr[0]
		} else {
			data.Iconid = 1001
		}
	*/

	rand := HF_GetRandom(10000)
	if rand >= 5000 {
		data.Iconid = 1002
	} else {
		data.Iconid = 1003
	}

	// 以前头像id是英雄id 现在改成了 10000000加英雄id
	/*
		rand.Seed(TimeServer().UnixNano())
		index := HF_RandInt(1, 6) - 1
		if heroes[index] != 0 {
			data.Iconid = 10000000 + heroes[index]
		} else {
			data.Iconid = 10000000 + heroes[0]
		}
	*/

	//LogDebug(heroes, data.Iconid, index)

	data.Deffight = cfg.Fight[1] + int64(HF_RandInt(1, int(cfg.Fight[0]-cfg.Fight[1])))

	for i := 0; i < len(heroes); i++ {
		err := data.FightTeamPos.addFightPos(i + 1)
		if err != nil {
			LogError("Get Robot err, err:", err.Error())
			continue
		}
		data.Defhero = append(data.Defhero, i)
		var hero JS_HeroInfo
		hero.Heroid = heroes[i]
		hero.Color = cfg.NpcQuality
		hero.HeroKeyId = i + 1
		hero.Stars = cfg.NpcStar[i]
		hero.HeroQuality = cfg.NpcStar[i]
		hero.Levels = cfg.NpcLv[i]
		hero.Skin = 0
		hero.Soldiercolor = 6
		hero.Skilllevel1 = 0
		hero.Skilllevel2 = 0
		hero.Skilllevel3 = 0
		hero.Skilllevel4 = SKILL_LEVEL
		hero.Skilllevel5 = 0
		hero.Skilllevel6 = 0
		hero.Fervor1 = 0
		hero.Fervor2 = 0
		hero.Fervor3 = 0
		hero.Fervor4 = 0
		hero.Fight = data.Deffight / 5

		hero.ArmsSkill = make([]JS_ArmsSkill, 0)
		hero.TalentSkill = []Js_TalentSkill{}
		hero.MainTalent = 0

		config := GetCsvMgr().HeroBreakConfigMap[hero.Heroid]
		if config == nil {
			continue
		}
		HeroBreakId := 0
		//计算突破等级
		for _, v := range config {
			if hero.Levels >= v.Break {
				HeroBreakId = v.Id
			}
		}

		skillBreakConfig := GetCsvMgr().HeroBreakConfigMap[hero.Heroid][HeroBreakId]
		if skillBreakConfig == nil {
			continue
		}

		for i := 0; i < len(skillBreakConfig.Skill); i++ {
			if skillBreakConfig.Skill[i] > 0 {
				hero.ArmsSkill = append(hero.ArmsSkill, JS_ArmsSkill{Id: skillBreakConfig.Skill[i] / 100, Level: skillBreakConfig.Skill[i] % 100})
			}
		}

		//monsterCfg := GetRobotMgr().GetMonsterConfig(cfg.MonsterId)
		//if monsterCfg == nil {
		//	continue
		//}
		att, att_ext, energy := hero.CountFight(cfg.BaseTypes, cfg.BaseValues)

		data.Heroinfo = append(data.Heroinfo, hero)
		var param JS_HeroParam
		param.Heroid = hero.Heroid
		param.Param = att
		param.ExtAttr = att_ext
		param.Hp = param.Param[AttrHp]
		param.Energy = energy
		param.Energy = 0

		data.HeroParam = append(data.HeroParam, param)

	}
	// 巨兽信息
	//data.BossInfo = GetRobotMgr().GetBossInfoByLevel(cfg.Hydra, cfg.NpcLv)

	//if data.BossInfo != nil {
	//	err := data.FightTeamPos.addFightPos(cfg.Hydra)
	//	if err != nil {
	//		LogError("Get Robot err, err:", err.Error())
	//	}
	//}

	//LogDebug("ROBOT FIGHT:", data.Deffight, cfg.Fight)

	return data
}

// 获取机器人信息，根据关卡表生成
func (self *RobotMgr) GetRobotByMonster(levelId int) *JS_FightInfo {

	data := new(JS_FightInfo)
	data.Uid = 0
	data.Defhero = make([]int, 0)
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)

	config := self.GetMonsterConfig(levelId)
	if config == nil {
		return data
	}

	data.Defhero = append(data.Defhero, config.HeroId)

	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)
	var heroinfo JS_HeroInfo
	heroinfo.Heroid = config.HeroId
	heroinfo.Color = 1
	heroinfo.Stars = config.MonsterIndex

	heroinfo.HeroKeyId = 1
	heroinfo.Levels = config.Level
	heroinfo.Skin = 0
	heroinfo.Skilllevel1 = 0
	heroinfo.Skilllevel2 = 0
	heroinfo.Skilllevel3 = 0
	heroinfo.Skilllevel4 = SKILL_LEVEL
	heroinfo.Skilllevel5 = 0
	heroinfo.Skilllevel6 = 0
	heroinfo.Fervor1 = 0
	heroinfo.Fervor2 = 0
	heroinfo.Fervor3 = 0
	heroinfo.Fervor4 = 0
	heroinfo.ArmsSkill = make([]JS_ArmsSkill, 0)
	heroinfo.TalentSkill = []Js_TalentSkill{}

	//新加的之前的不敢动
	heroinfo.HeroQuality = config.MonsterIndex
	/*
		artifact := player.GetModule("artifactequip").(*ModArtifactEquip).GetArtifactEquipItem(heroData.ArtifactEquipIds[0])
		if artifact != nil {
			heroinfo.HeroArtifactId = artifact.Id
			heroinfo.HeroArtifactLv = artifact.Lv
		}
	*/
	//专属
	/*
		if heroData.ExclusiveEquip != nil {
			if heroData.ExclusiveEquip.UnLock == LOGIC_TRUE {
				heroinfo.HeroExclusiveLv = heroData.ExclusiveEquip.Lv
			}
		}
	*/

	var param JS_HeroParam
	param.Heroid = config.HeroId

	pAttrWrapper := &AttrWrapper{
		Base:     make([]float64, AttrEnd),
		Ext:      make(map[int]float64),
		Per:      make(map[int]float64),
		Energy:   0,
		FightNum: 0,
	}
	attMap := make(map[int]*Attribute)
	for i := 0; i < len(config.BaseType); i++ {
		if config.BaseType[i] == 0 {
			continue
		}
		_, ok := attMap[config.BaseType[i]]
		if !ok {
			attMap[config.BaseType[i]] = &Attribute{config.BaseType[i], config.BaseValue[i]}
		} else {
			if attMap[config.BaseType[i]] != nil {
				attMap[config.BaseType[i]].AttValue += config.BaseValue[i]
			} else {
				LogError("Level_Monster:error")
			}
		}
		if config.BaseType[i] == AttrFight {
			heroinfo.Fight = config.BaseValue[i]
			data.Deffight += heroinfo.Fight
		}
	}
	ProcAtt(attMap, pAttrWrapper)
	ProcLast(1, pAttrWrapper)
	pAttrWrapper.ExtRet = ProcExtAttr(pAttrWrapper)
	param.Param, param.ExtAttr, param.Energy = pAttrWrapper.Base, pAttrWrapper.ExtRet, pAttrWrapper.Energy
	param.Hp = param.Param[AttrHp]
	if config.MainSkill != 0 {
		heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: config.MainSkill, Level: 1})
	}

	data.Heroinfo = append(data.Heroinfo, heroinfo)
	data.HeroParam = append(data.HeroParam, param)
	data.FightTeamPos.FightPos[3] = heroinfo.HeroKeyId
	return data
}

//生成地牢秘宝守卫
func (self *RobotMgr) GetRobotByPitCart(level int, star int) *JS_FightInfo {

	data := new(JS_FightInfo)
	data.Uid = 0
	data.Defhero = make([]int, 0)
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)

	configHeroLst := GetCsvMgr().GetNewPitRobot(star, LOGIC_TRUE)
	if len(configHeroLst) == 0 {
		return data
	}

	rand := HF_GetRandom(len(configHeroLst))

	data.Defhero = append(data.Defhero, configHeroLst[rand].HeroId)

	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)
	var heroinfo JS_HeroInfo
	heroinfo.Heroid = configHeroLst[rand].HeroId
	heroinfo.Color = 1
	heroinfo.Stars = 1
	heroinfo.Stars = star

	heroinfo.HeroKeyId = 1
	heroinfo.Levels = level
	heroinfo.Skin = 0
	heroinfo.Skilllevel1 = 0
	heroinfo.Skilllevel2 = 0
	heroinfo.Skilllevel3 = 0
	heroinfo.Skilllevel4 = SKILL_LEVEL
	heroinfo.Skilllevel5 = 0
	heroinfo.Skilllevel6 = 0
	heroinfo.Fervor1 = 0
	heroinfo.Fervor2 = 0
	heroinfo.Fervor3 = 0
	heroinfo.Fervor4 = 0
	heroinfo.ArmsSkill = make([]JS_ArmsSkill, 0)
	heroinfo.TalentSkill = []Js_TalentSkill{}

	//新加的之前的不敢动
	heroinfo.HeroQuality = star

	var param JS_HeroParam
	param.Heroid = configHeroLst[rand].HeroId

	pAttrWrapper := &AttrWrapper{
		Base:     make([]float64, AttrEnd),
		Ext:      make(map[int]float64),
		Per:      make(map[int]float64),
		Energy:   0,
		FightNum: 0,
	}
	attMap := make(map[int]*Attribute)
	//先计算英雄属性  走英雄接口方便计算
	heroTemp := new(Hero)
	heroTemp.HeroId = configHeroLst[rand].HeroId
	heroTemp.checkStarItem(star)
	heroTemp.LvUp(level)
	skillBreakConfig := GetCsvMgr().HeroBreakConfigMap[heroTemp.HeroId][heroTemp.StarItem.HeroBreakId]
	if skillBreakConfig == nil {
		return data
	}

	for i := 0; i < len(skillBreakConfig.Skill); i++ {
		if skillBreakConfig.Skill[i] > 0 {
			heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: skillBreakConfig.Skill[i] / 100, Level: skillBreakConfig.Skill[i] % 100})
		}
	}

	//计算属性
	heroTemp.StarItem.AttrMap = make(map[int]*Attribute)
	//突破属性
	AddAttrHelperForTimes(heroTemp.StarItem.AttrMap, skillBreakConfig.BaseTypes, skillBreakConfig.BaseValues, 1)
	//星级属性
	configStar := GetCsvMgr().GetHeroMapConfig(heroTemp.HeroId, heroTemp.StarItem.UpStar)
	if configStar == nil {
		return data
	}

	AddAttrHelperForTimes(heroTemp.StarItem.AttrMap, configStar.BaseTypes, configStar.BaseValues, 1)
	//成长率
	configGrowth := GetCsvMgr().HeroGrowthConfigMap[heroTemp.HeroLv]
	if configGrowth != nil {
		AddAttrHelperForGrowth(heroTemp.StarItem.AttrMap, configStar.GrowthTypes, configStar.GrowthValues, configGrowth.GrowthType, configGrowth.GrowthValue, configStar.QuaType, configStar.QuaValue)
	}
	ProcAtt(heroTemp.StarItem.AttrMap, pAttrWrapper)
	//在计算战斗附加属性
	configAttr := GetCsvMgr().GetNewPitRobotAttr(configHeroLst[rand].GroupAttribute, level)
	for i := 0; i < len(configAttr.GrowthType); i++ {
		if configAttr.GrowthType[i] == 0 {
			continue
		}
		_, ok := attMap[configAttr.GrowthType[i]]
		if !ok {
			attMap[configAttr.GrowthType[i]] = &Attribute{configAttr.GrowthType[i], configAttr.GrowthValue[i]}
		} else {
			if attMap[configAttr.GrowthType[i]] != nil {
				attMap[configAttr.GrowthType[i]].AttValue += configAttr.GrowthValue[i]
			} else {
				LogError("Level_Monster:error")
			}
		}
	}
	ProcAtt(attMap, pAttrWrapper)
	ProcLast(1, pAttrWrapper)
	pAttrWrapper.ExtRet = ProcExtAttr(pAttrWrapper)
	param.Param, param.ExtAttr, param.Energy = pAttrWrapper.Base, pAttrWrapper.ExtRet, pAttrWrapper.Energy
	param.Hp = param.Param[AttrHp]
	heroinfo.Fight = pAttrWrapper.FightNum
	data.Deffight = pAttrWrapper.FightNum

	data.Heroinfo = append(data.Heroinfo, heroinfo)
	data.HeroParam = append(data.HeroParam, param)
	data.FightTeamPos.FightPos[3] = heroinfo.HeroKeyId
	return data
}

// 添加一个带佣兵的数据结构
func (self *RobotMgr) GetPlayerFightInfoWithArmyByPos(player *Player, inParam int, encourageType int, teamPos int, armyInfo *ArmyInfo) *JS_FightInfo {
	hero := player.GetModule("friend").(*ModFriend).GetHireHero(armyInfo.SelfKey)
	if nil == hero {
		return nil
	}

	var pos = player.getTeamPosByType(teamPos)

	data := new(JS_FightInfo)
	data.Rankid = inParam
	data.Uid = player.Sql_UserBase.Uid
	data.Uname = player.Sql_UserBase.UName
	data.Iconid = player.Sql_UserBase.IconId
	data.Portrait = player.Sql_UserBase.Portrait //边框  20190412 by zy
	data.Level = player.Sql_UserBase.Level
	data.Vip = player.Sql_UserBase.Vip
	data.UnionName = player.GetUnionName()
	data.Camp = player.Sql_UserBase.Camp
	if pos == nil {
		return nil
	}
	heroIds := make([]int, 0)
	var army *NewHero = nil
	HF_DeepCopy(&data.FightTeamPos, pos)
	for i := 0; i < len(pos.FightPos); i++ {
		if i == armyInfo.Pos-1 {
			army = player.GetModule("friend").(*ModFriend).GetHireHero(armyInfo.SelfKey)
			if army == nil {
				continue
			}
			if army.HeroId == 0 {
				continue
			}
			data.Defhero = append(data.Defhero, armyInfo.SelfKey)
			heroIds = append(heroIds, army.HeroId)
			data.FightTeamPos.FightPos[i] = armyInfo.SelfKey
		} else {
			heroid := pos.FightPos[i]
			if heroid == 0 {
				continue
			}

			data.Defhero = append(data.Defhero, heroid)
			heroData := player.getHero(heroid)
			if heroData == nil {
				continue
			}
			heroIds = append(heroIds, heroData.HeroId)
		}
	}

	teamAttr := self.GetTeamAttr(heroIds)
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)
	for i, key := range data.FightTeamPos.FightPos {
		if key == 0 {
			continue
		}
		var heroData *Hero = nil
		if i != armyInfo.Pos-1 {
			heroData = player.getHero(key)
			if heroData == nil {
				continue
			}

			var heroinfo JS_HeroInfo
			heroinfo.Heroid = heroData.HeroId
			heroinfo.Color = 1
			heroinfo.Stars = 1
			if heroData.StarItem != nil {
				heroinfo.Stars = heroData.StarItem.UpStar
			}

			heroinfo.HeroKeyId = heroData.HeroKeyId
			heroinfo.Levels = heroData.HeroLv
			heroinfo.Skin = heroData.Skin
			heroinfo.Skilllevel1 = 0
			heroinfo.Skilllevel2 = 0
			heroinfo.Skilllevel3 = 0
			heroinfo.Skilllevel4 = SKILL_LEVEL
			heroinfo.Skilllevel5 = 0
			heroinfo.Skilllevel6 = 0
			heroinfo.Fervor1 = 0
			heroinfo.Fervor2 = 0
			heroinfo.Fervor3 = 0
			heroinfo.Fervor4 = 0
			heroinfo.ArmsSkill = make([]JS_ArmsSkill, 0)
			heroinfo.TalentSkill = []Js_TalentSkill{}

			//新加的之前的不敢动
			heroinfo.HeroQuality = heroinfo.Stars
			if len(heroData.ArtifactEquipIds) > 0 {
				artifact := player.GetModule("artifactequip").(*ModArtifactEquip).GetArtifactEquipItem(heroData.ArtifactEquipIds[0])
				if artifact != nil {
					heroinfo.HeroArtifactId = artifact.Id
					heroinfo.HeroArtifactLv = artifact.Lv
				}
			}
			//专属
			if heroData.ExclusiveEquip != nil {
				if heroData.ExclusiveEquip.UnLock == LOGIC_TRUE {
					heroinfo.ExclusiveUnLock = LOGIC_TRUE
					heroinfo.HeroExclusiveLv = heroData.ExclusiveEquip.Lv
				}
			}

			var param JS_HeroParam
			param.Heroid = heroinfo.Heroid
			pAttrWrapper := heroData.GetAttr2(player, encourageType, teamAttr)
			param.Param, param.ExtAttr, param.Energy = pAttrWrapper.Base, pAttrWrapper.ExtRet, pAttrWrapper.Energy
			param.Hp = param.Param[AttrHp]
			heroinfo.Fight = heroData.Fight

			//升星技能
			for i := 0; i < len(heroData.StarItem.Skills); i++ {
				if heroData.StarItem.Skills[i] == 0 {
					continue
				}
				heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: heroData.StarItem.Skills[i] / 100, Level: heroData.StarItem.Skills[i] % 100})
			}

			if heroData.StageTalent != nil {
				// 天赋技能
				for i := 0; i < len(heroData.StageTalent.AllSkill); i++ {
					config := GetCsvMgr().GetStageTalent(heroData.StageTalent.AllSkill[i].ID)
					if nil == config {
						continue
					}
					if heroData.StageTalent.AllSkill[i].Pos <= 0 || heroData.StageTalent.AllSkill[i].Pos > len(config.Skill) {
						continue
					}
					heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: config.Skill[heroData.StageTalent.AllSkill[i].Pos-1], Level: 1})
				}
			}

			// 生命树技能
			skill := player.GetModule("lifetree").(*ModLifeTree).GetAllSkill(heroinfo.Heroid, hero.StarItem.UpStar)
			for i := 0; i < len(skill); i++ {
				if skill[i] > 0 {
					heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: skill[i] / 100, Level: skill[i] % 100})
				}
			}

			//新加的之前的不敢动
			heroinfo.HeroQuality = heroinfo.Stars
			if len(heroData.ArtifactEquipIds) > 0 {
				artifact := player.GetModule("artifactequip").(*ModArtifactEquip).GetArtifactEquipItem(heroData.ArtifactEquipIds[0])
				if artifact != nil {
					heroinfo.HeroArtifactId = artifact.Id
					heroinfo.HeroArtifactLv = artifact.Lv
					//看看是否达到上限
					configArt := GetCsvMgr().GetArtifactStrengthenLvUpConfig(artifact.Id, artifact.Lv)
					if configArt != nil {
						for i := 0; i < len(configArt.Skill); i++ {
							heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: configArt.Skill[i] / 100, Level: configArt.Skill[i] % 100})
						}
					}
				}
			}

			if heroData.ExclusiveEquip != nil {
				if heroData.ExclusiveEquip.UnLock == LOGIC_TRUE {
					heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: heroData.ExclusiveEquip.Skill / 100, Level: heroData.ExclusiveEquip.Skill % 100})
				}
			}

			data.Deffight += heroData.Fight
			data.Heroinfo = append(data.Heroinfo, heroinfo)
			data.HeroParam = append(data.HeroParam, param)
		} else {
			if army == nil {
				continue
			}
			var heroinfo JS_HeroInfo
			heroinfo.Heroid = army.HeroId
			heroinfo.Color = 1
			if army.StarItem != nil {
				heroinfo.Stars = army.StarItem.UpStar
				heroinfo.HeroQuality = army.StarItem.UpStar
			}

			heroinfo.HeroKeyId = armyInfo.SelfKey
			heroinfo.Levels = armyInfo.Lv
			heroinfo.Skin = army.Skin
			heroinfo.IsArmy = LOGIC_TRUE
			heroinfo.Skilllevel1 = 0
			heroinfo.Skilllevel2 = 0
			heroinfo.Skilllevel3 = 0
			heroinfo.Skilllevel4 = SKILL_LEVEL
			heroinfo.Skilllevel5 = 0
			heroinfo.Skilllevel6 = 0
			heroinfo.Fervor1 = 0
			heroinfo.Fervor2 = 0
			heroinfo.Fervor3 = 0
			heroinfo.Fervor4 = 0
			heroinfo.ArmsSkill = make([]JS_ArmsSkill, 0)
			heroinfo.TalentSkill = []Js_TalentSkill{}

			//新加的之前的不敢动
			heroinfo.HeroQuality = heroinfo.Stars
			if len(army.ArtifactEquipIds) > 0 {
				heroinfo.HeroArtifactId = army.ArtifactEquipIds[0].Id
				heroinfo.HeroArtifactLv = army.ArtifactEquipIds[0].Lv
			}
			//专属
			if army.ExclusiveEquip != nil {
				if army.ExclusiveEquip.UnLock == LOGIC_TRUE {
					heroinfo.ExclusiveUnLock = LOGIC_TRUE
					heroinfo.HeroExclusiveLv = army.ExclusiveEquip.Lv
				}
			}

			var param JS_HeroParam
			param.Heroid = heroinfo.Heroid

			pAttr := &AttrWrapper{
				Base:     make([]float64, AttrEnd),
				Ext:      make(map[int]float64),
				Per:      make(map[int]float64),
				Energy:   0,
				FightNum: 0,
			}
			attMap := make(map[int]*Attribute, 0)
			for _, v := range armyInfo.Atts {
				data := new(Attribute)
				data.AttType = v.Type
				data.AttValue = v.Value
				attMap[v.Type] = data
			}
			ProcAtt(attMap, pAttr)
			param.Param = pAttr.Base
			param.ExtAttr = pAttr.ExtRet
			param.Hp = param.Param[AttrHp]
			heroinfo.Fight = pAttr.FightNum

			for _, v := range armyInfo.Skills {
				heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: v.Id, Level: v.Level})
			}

			//看看是否达到上限
			configArt := GetCsvMgr().GetArtifactStrengthenLvUpConfig(heroinfo.HeroArtifactId, heroinfo.HeroArtifactLv)
			if configArt != nil {
				for i := 0; i < len(configArt.Skill); i++ {
					heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: configArt.Skill[i] / 100, Level: configArt.Skill[i] % 100})
				}
			}

			if heroinfo.ExclusiveUnLock == LOGIC_TRUE {
				for _, v := range GetCsvMgr().ExclusiveEquipConfigMap {
					if v.HeroId == heroinfo.Heroid {
						//获得英雄专属和等级，加入技能
						config := GetCsvMgr().ExclusiveStrengthen[v.Id][heroinfo.HeroExclusiveLv]
						if config != nil {
							heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: config.Skill / 100, Level: config.Skill % 100})
						}
						break
					}
				}
			}

			data.Deffight += army.Fight
			data.Heroinfo = append(data.Heroinfo, heroinfo)
			data.HeroParam = append(data.HeroParam, param)
		}
	}

	return data
}

func (self TeamCalList) Len() int {
	return len(self.list)
}
func (self TeamCalList) Less(i, j int) bool {
	return self.list[i].Num > self.list[j].Num
}
func (self TeamCalList) Swap(i, j int) {
	self.list[i], self.list[j] = self.list[j], self.list[i]
}

func (self *RobotMgr) GetTeamAttr(hero []int) map[int]*Attribute {
	if len(hero) == 0 {
		return nil
	}
	attMap := make(map[int]*Attribute)
	calCampTeam := make(map[int]int)
	calCamp := new(TeamCalList)
	campRel_1 := 0
	campRel_2 := 0
	for _, v := range hero {
		config := GetCsvMgr().GetHeroMapConfig(v, 1)
		if config == nil {
			continue
		}
		calCampTeam[config.Attribute] += 1
	}
	for k, v := range calCampTeam {
		if (k >= HERO_ATTRIBUTE_WATER && k <= HERO_ATTRIBUTE_EARTH) || k == HERO_ATTRIBUTE_AIR {
			calCamp.list = append(calCamp.list, TeamCal{Camp: k, Num: v})
		}
	}
	sort.Sort(calCamp)
	if len(calCamp.list) > 0 {
		campRel_1 = calCamp.list[0].Num
	}
	if len(calCamp.list) > 1 && calCamp.list[1].Num == 2 {
		campRel_2 = calCamp.list[1].Num
	}
	campRel_1 += calCampTeam[HERO_ATTRIBUTE_LIGHT]
	//计算普通阵营加成
	config := GetCsvMgr().GetTeamAttrConfig(1, campRel_1, campRel_2)
	if config != nil {
		AddAttrHelperForTimes(attMap, config.Base_type, config.Base_value, 1)
	}
	//计算暗英雄
	for i := 1; i <= calCampTeam[HERO_ATTRIBUTE_DARK]; i++ {
		config = GetCsvMgr().GetTeamAttrConfig(2, i, 0)
		if config != nil {
			AddAttrHelperForTimes(attMap, config.Base_type, config.Base_value, 1)
		}
	}
	return attMap
}

// GetPlayerFightInfo 的升级版，根据预设阵容生成拷贝
func (self *RobotMgr) GetPlayerFightInfoByPos(player *Player, inParam int, encourageType int, teamPos int) *JS_FightInfo {
	var pos = player.getTeamPosByType(teamPos)
	if pos == nil {
		return nil
	}

	data := new(JS_FightInfo)
	data.Rankid = inParam
	data.Uid = player.Sql_UserBase.Uid
	data.Uname = player.Sql_UserBase.UName
	data.Iconid = player.Sql_UserBase.IconId
	data.Portrait = player.Sql_UserBase.Portrait //边框  20190412 by zy
	data.Level = player.Sql_UserBase.Level
	data.Vip = player.Sql_UserBase.Vip
	data.UnionName = player.GetUnionName()
	data.Camp = player.Sql_UserBase.Camp

	modLifeTree := player.GetModule("lifetree").(*ModLifeTree)
	if modLifeTree != nil {
		data.LifeTreeInfo = new(JS_LifeTreeInfo)
		data.LifeTreeInfo.MainLevel = modLifeTree.San_LifeTree.MainLevel
		data.LifeTreeInfo.Info = modLifeTree.San_LifeTree.info
	}

	data.FightTeamPos = *pos
	for i := 0; i < len(pos.FightPos); i++ {
		heroid := pos.FightPos[i]
		if heroid == 0 {
			continue
		}

		data.Defhero = append(data.Defhero, heroid)
	}
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)
	heroIds := make([]int, 0)
	for _, key := range data.Defhero {
		heroData := player.getHero(key)
		if heroData == nil {
			continue
		}
		heroIds = append(heroIds, heroData.HeroId)
	}
	teamAttr := self.GetTeamAttr(heroIds)
	for _, key := range data.Defhero {
		heroData := player.getHero(key)
		if heroData == nil {
			continue
		}
		var heroinfo JS_HeroInfo
		heroinfo.Heroid = heroData.HeroId
		heroinfo.Color = 1
		heroinfo.Stars = 1
		if heroData.StarItem != nil {
			heroinfo.Stars = heroData.StarItem.UpStar
		}

		heroinfo.HeroKeyId = heroData.HeroKeyId
		heroinfo.Levels = heroData.HeroLv
		heroinfo.Skin = heroData.Skin
		heroinfo.Skilllevel1 = 0
		heroinfo.Skilllevel2 = 0
		heroinfo.Skilllevel3 = 0
		heroinfo.Skilllevel4 = SKILL_LEVEL
		heroinfo.Skilllevel5 = 0
		heroinfo.Skilllevel6 = 0
		heroinfo.Fervor1 = 0
		heroinfo.Fervor2 = 0
		heroinfo.Fervor3 = 0
		heroinfo.Fervor4 = 0
		heroinfo.ArmsSkill = make([]JS_ArmsSkill, 0)
		heroinfo.TalentSkill = []Js_TalentSkill{}

		//新加的之前的不敢动
		heroinfo.HeroQuality = heroinfo.Stars
		if len(heroData.ArtifactEquipIds) > 0 {
			artifact := player.GetModule("artifactequip").(*ModArtifactEquip).GetArtifactEquipItem(heroData.ArtifactEquipIds[0])
			if artifact != nil {
				heroinfo.HeroArtifactId = artifact.Id
				heroinfo.HeroArtifactLv = artifact.Lv
				//看看是否达到上限
				configArt := GetCsvMgr().GetArtifactStrengthenLvUpConfig(artifact.Id, artifact.Lv)
				if configArt != nil {
					for i := 0; i < len(configArt.Skill); i++ {
						heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: configArt.Skill[i] / 100, Level: configArt.Skill[i] % 100})
					}
				}
			}
		}
		//专属
		if heroData.ExclusiveEquip != nil {
			if heroData.ExclusiveEquip.UnLock == LOGIC_TRUE {
				heroinfo.HeroExclusiveLv = heroData.ExclusiveEquip.Lv
			}
		}

		var param JS_HeroParam
		param.Heroid = heroinfo.Heroid
		pAttrWrapper := heroData.GetAttr2(player, encourageType, teamAttr)
		param.Param, param.ExtAttr, param.Energy = pAttrWrapper.Base, pAttrWrapper.ExtRet, pAttrWrapper.Energy
		param.Hp = param.Param[AttrHp]
		heroinfo.Fight = heroData.Fight

		//升星技能
		for i := 0; i < len(heroData.StarItem.Skills); i++ {
			if heroData.StarItem.Skills[i] == 0 {
				continue
			}
			heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: heroData.StarItem.Skills[i] / 100, Level: heroData.StarItem.Skills[i] % 100})
		}

		if heroData.StageTalent != nil {
			// 天赋技能
			for i := 0; i < len(heroData.StageTalent.AllSkill); i++ {
				config := GetCsvMgr().GetStageTalent(heroData.StageTalent.AllSkill[i].ID)
				if nil == config {
					continue
				}
				if heroData.StageTalent.AllSkill[i].Pos <= 0 || heroData.StageTalent.AllSkill[i].Pos > len(config.Skill) {
					continue
				}
				heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: config.Skill[heroData.StageTalent.AllSkill[i].Pos-1] / 100, Level: config.Skill[heroData.StageTalent.AllSkill[i].Pos-1] % 100})
			}
		}

		// 生命树技能
		skill := player.GetModule("lifetree").(*ModLifeTree).GetAllSkill(heroinfo.Heroid, heroinfo.Stars)
		for i := 0; i < len(skill); i++ {
			if skill[i] > 0 {
				heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: skill[i] / 100, Level: skill[i] % 100})
			}
		}

		//魔宠技能
		if heroData.Horse > 0 {
			horsedata := player.GetModule("horse").(*ModHorse).GetHorseSafe(heroData.Horse)
			if horsedata != nil && horsedata.Skill > 0 {
				heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: horsedata.Skill / 100, Level: horsedata.Skill % 100})
			}
		}

		if heroData.ExclusiveEquip != nil {
			if heroData.ExclusiveEquip.UnLock == LOGIC_TRUE {
				heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: heroData.ExclusiveEquip.Skill / 100, Level: heroData.ExclusiveEquip.Skill % 100})
			}
		}

		data.Deffight += heroData.Fight
		data.Heroinfo = append(data.Heroinfo, heroinfo)
		data.HeroParam = append(data.HeroParam, param)
	}

	return data
}

func (self *RobotMgr) GetPlayerFightInfoByPosShow(player *Player, data *JS_FightInfo) {
	if player == nil || data == nil {
		return
	}
	heroIds := make([]int, 0)
	for _, key := range data.Defhero {
		heroData := player.getHero(key)
		if heroData == nil {
			continue
		}
		heroIds = append(heroIds, heroData.HeroId)
	}
	teamAttr := self.GetTeamAttr(heroIds)
	for _, key := range data.Defhero {
		heroData := player.getHero(key)
		if heroData == nil {
			continue
		}
		var heroinfo JS_HeroInfo
		heroinfo.Heroid = heroData.HeroId
		heroinfo.Color = 1
		heroinfo.Stars = 1
		if heroData.StarItem != nil {
			heroinfo.Stars = heroData.StarItem.UpStar
		}

		heroinfo.HeroKeyId = heroData.HeroKeyId
		heroinfo.Levels = heroData.HeroLv
		heroinfo.Skin = heroData.Skin
		heroinfo.Skilllevel1 = 0
		heroinfo.Skilllevel2 = 0
		heroinfo.Skilllevel3 = 0
		heroinfo.Skilllevel4 = SKILL_LEVEL
		heroinfo.Skilllevel5 = 0
		heroinfo.Skilllevel6 = 0
		heroinfo.Fervor1 = 0
		heroinfo.Fervor2 = 0
		heroinfo.Fervor3 = 0
		heroinfo.Fervor4 = 0
		heroinfo.ArmsSkill = make([]JS_ArmsSkill, 0)
		heroinfo.TalentSkill = []Js_TalentSkill{}

		//新加的之前的不敢动
		heroinfo.HeroQuality = heroinfo.Stars
		if len(heroData.ArtifactEquipIds) > 0 {
			artifact := player.GetModule("artifactequip").(*ModArtifactEquip).GetArtifactEquipItem(heroData.ArtifactEquipIds[0])
			if artifact != nil {
				heroinfo.HeroArtifactId = artifact.Id
				heroinfo.HeroArtifactLv = artifact.Lv
			}
		}
		//专属
		if heroData.ExclusiveEquip != nil {
			if heroData.ExclusiveEquip.UnLock == LOGIC_TRUE {
				heroinfo.HeroExclusiveLv = heroData.ExclusiveEquip.Lv
			}
		}

		var param JS_HeroParam
		param.Heroid = heroinfo.Heroid
		pAttrWrapper := heroData.GetAttr2(player, 0, teamAttr)
		param.Param, param.ExtAttr, param.Energy = pAttrWrapper.Base, pAttrWrapper.ExtRet, pAttrWrapper.Energy
		param.Hp = param.Param[AttrHp]
		heroinfo.Fight = heroData.Fight

		//升星技能
		for i := 0; i < len(heroData.StarItem.Skills); i++ {
			if heroData.StarItem.Skills[i] == 0 {
				continue
			}
			heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: heroData.StarItem.Skills[i] / 100, Level: heroData.StarItem.Skills[i] % 100})
		}

		if heroData.StageTalent != nil {
			// 天赋技能
			for i := 0; i < len(heroData.StageTalent.AllSkill); i++ {
				config := GetCsvMgr().GetStageTalent(heroData.StageTalent.AllSkill[i].ID)
				if nil == config {
					continue
				}
				if heroData.StageTalent.AllSkill[i].Pos <= 0 || heroData.StageTalent.AllSkill[i].Pos > len(config.Skill) {
					continue
				}
				heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: config.Skill[heroData.StageTalent.AllSkill[i].Pos-1], Level: 1})
			}
		}

		// 生命树技能
		skill := player.GetModule("lifetree").(*ModLifeTree).GetAllSkill(heroinfo.Heroid, heroinfo.Stars)
		for i := 0; i < len(skill); i++ {
			if skill[i] > 0 {
				heroinfo.ArmsSkill = append(heroinfo.ArmsSkill, JS_ArmsSkill{Id: skill[i] / 100, Level: skill[i] % 100})
			}
		}

		data.Deffight += heroData.Fight
		data.Heroinfo = append(data.Heroinfo, heroinfo)
		data.HeroParam = append(data.HeroParam, param)
	}
}

//根据原版HERO生成新版HERO
func (self *RobotMgr) HeroToNewHero(uid int64, hero *Hero) *NewHero {
	player := GetPlayerMgr().GetPlayer(uid, true)
	newHero := new(NewHero)
	if hero == nil {
		return newHero
	}
	heroConfig := GetCsvMgr().GetHeroMapConfig(hero.HeroId, hero.StarItem.UpStar)
	if heroConfig == nil {
		return newHero
	}
	newHero.Uid = hero.Uid
	newHero.HeroId = hero.HeroId
	newHero.HeroKeyId = hero.HeroKeyId
	newHero.Skin = hero.Skin
	//生成星级
	newHero.checkStarItem(hero.StarItem.UpStar)
	newHero.LvUp(hero.HeroLv)
	//生成专属
	newHero.ExclusiveEquip = new(ExclusiveEquip)
	HF_DeepCopy(&newHero.ExclusiveEquip, &hero.ExclusiveEquip)
	// 天赋
	if hero.StageTalent != nil {
		newHero.StageTalent = new(StageTalent)
		HF_DeepCopy(newHero.StageTalent, hero.StageTalent)
	}

	//生成装备
	equips := player.getEquips()
	for i := 0; i < len(equips); i++ {
		if equips[i] == nil {
			return newHero
		}
	}
	for i := 0; i < len(hero.EquipIds); i++ {
		isFind := false
		for j := 0; j < len(equips); j++ {
			if equips[j] == nil {
				continue
			}
			equipItem := equips[j][hero.EquipIds[i]]
			if equipItem == nil {
				continue
			}
			equip := &Equip{}
			equip.Id = equipItem.Id
			equip.Lv = equipItem.Lv
			newHero.EquipIds = append(newHero.EquipIds, equip)
			isFind = true
			break
		}
		if !isFind {
			equip := &Equip{}
			equip.Id = 0
			newHero.EquipIds = append(newHero.EquipIds, equip)
		}
	}
	//生成神器
	//ArtifactEquipIds
	//[]int
	//`json:"artifactequipids"`
	if len(hero.ArtifactEquipIds) > 0 {
		artifact := player.GetModule("artifactequip").(*ModArtifactEquip).GetArtifactEquipItem(hero.ArtifactEquipIds[0])
		if artifact != nil {
			newHero.ArtifactEquipIds = make([]*ArtifactEquip, 0)
			acti := &ArtifactEquip{}
			HF_DeepCopy(&acti, &artifact)
			newHero.ArtifactEquipIds = append(newHero.ArtifactEquipIds, acti)
		}
	}
	newHero.CalAttr()
	return newHero
}

// 获取机器人信息，根据共鸣水晶核心等级生成
func (self *RobotMgr) GetRobotByWorldLv(id int, level int) *JS_FightInfo {
	activity := GetActivityMgr().GetActivity(id)
	if activity == nil {
		return nil
	}

	period := GetActivityMgr().getActN3(id)

	bossConfig := GetCsvMgr().GetActivityBossConfig(id, period)
	if bossConfig == nil {
		return nil
	}
	//取动态配置
	configWorldLevelConfig := GetCsvMgr().GetWorldLevelConfig(level)
	if configWorldLevelConfig == nil {
		return nil
	}

	configWorldLevelTypeConfig := GetCsvMgr().GetWorldLevelTypeConfig(bossConfig.BossId, level)
	if configWorldLevelTypeConfig == nil {
		return nil
	}

	worldShowBossConfig := GetCsvMgr().GetWorldShowBossConfig(level)
	if worldShowBossConfig == nil {
		return nil
	}

	data := new(JS_FightInfo)
	data.Rankid = 0
	data.Uid = 0
	data.Uname = bossConfig.Name
	data.Iconid = 10000000 + bossConfig.BossId

	data.Portrait = 1000 //机器人边框  20190412 by zy
	data.Level = configWorldLevelConfig.NpcLv
	data.Defhero = make([]int, 0)
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)

	var heroes []int
	heroes = append(heroes, configWorldLevelTypeConfig.HeroNpc)

	heroes = HF_GetRandomArr(heroes, 5)

	for i := 0; i < len(heroes); i++ {
		pos := i + 1 //相当于KEY
		fightPos := 0
		for j := 0; j < len(bossConfig.Position); j++ {
			if bossConfig.Position[j] == heroes[i] {
				fightPos = j
			}
		}
		data.FightTeamPos.FightPos[fightPos] = pos
		data.Defhero = append(data.Defhero, pos)
		var hero JS_HeroInfo
		hero.Heroid = heroes[i]
		hero.Color = worldShowBossConfig.ShowHeroqua
		hero.HeroKeyId = pos
		hero.Stars = worldShowBossConfig.ShowHeroqua
		hero.HeroQuality = worldShowBossConfig.ShowHeroqua
		hero.Levels = configWorldLevelConfig.NpcLv
		hero.Skin = 0
		hero.Soldiercolor = 6
		hero.Skilllevel1 = 0
		hero.Skilllevel2 = 0
		hero.Skilllevel3 = 0
		hero.Skilllevel4 = SKILL_LEVEL
		hero.Skilllevel5 = 0
		hero.Skilllevel6 = 0
		hero.Fervor1 = 0
		hero.Fervor2 = 0
		hero.Fervor3 = 0
		hero.Fervor4 = 0
		hero.Fight = data.Deffight / 5

		hero.ArmsSkill = make([]JS_ArmsSkill, 0)
		hero.TalentSkill = []Js_TalentSkill{}
		hero.MainTalent = 0

		config := GetCsvMgr().HeroBreakConfigMap[hero.Heroid]
		if config == nil {
			continue
		}
		HeroBreakId := 0
		//计算突破等级
		for _, v := range config {
			if hero.Levels >= v.Break {
				HeroBreakId = v.Id
			}
		}

		skillBreakConfig := GetCsvMgr().HeroBreakConfigMap[hero.Heroid][HeroBreakId]
		if skillBreakConfig == nil {
			continue
		}

		for i := 0; i < len(skillBreakConfig.Skill); i++ {
			if skillBreakConfig.Skill[i] > 0 {
				hero.ArmsSkill = append(hero.ArmsSkill, JS_ArmsSkill{Id: skillBreakConfig.Skill[i] / 100, Level: skillBreakConfig.Skill[i] % 100})
			}
		}
		hero.ExclusiveUnLock = configWorldLevelTypeConfig.HeroEquip
		//计算专属
		if hero.ExclusiveUnLock == LOGIC_TRUE {

			pItem := &ExclusiveEquip{}
			pItem.Lv = worldShowBossConfig.ShowEquip
			for _, v := range GetCsvMgr().ExclusiveEquipConfigMap {
				if v.HeroId == hero.Heroid {
					pItem.Id = v.Id
					for i := 0; i < len(v.BaseType); i++ {
						if v.BaseType[i] > 0 {
							attr := new(AttrInfo)
							attr.AttrId = i + 1
							pItem.AttrInfo = append(pItem.AttrInfo, attr)
						}
					}
				}
			}
			pItem.CalAttr()

			hero.ArmsSkill = append(hero.ArmsSkill, JS_ArmsSkill{Id: pItem.Skill / 100, Level: pItem.Skill % 100})
		}

		newValue := make([]float64, 0)
		for i := 0; i < len(configWorldLevelTypeConfig.BaseValue); i++ {
			value := configWorldLevelConfig.BaseValue[i] * (PER_BIT + configWorldLevelTypeConfig.BaseValue[i]) / PER_BIT
			newValue = append(newValue, float64(value))
		}

		att, att_ext, energy := hero.CountFight(configWorldLevelTypeConfig.BaseType, newValue)

		data.Heroinfo = append(data.Heroinfo, hero)
		var param JS_HeroParam
		param.Heroid = hero.Heroid
		param.Param = att
		param.ExtAttr = att_ext
		param.Hp = param.Param[AttrHp]
		param.Energy = energy
		param.Energy = 0
		data.HeroParam = append(data.HeroParam, param)

		data.Deffight += hero.Fight
	}
	return data
}
