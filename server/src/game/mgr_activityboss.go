package game

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

type ActivityBossTopNode struct {
	Uid           int64                `json:"uid"`           //玩家UID
	Name          string               `json:"name"`          //名字
	IconId        int                  `json:"icon"`          //头像
	Level         int                  `json:"level"`         //等级
	Portrait      int                  `json:"portrait"`      //头像框
	Score         int64                `json:"point"`         //分数
	TopRank       int                  `json:"toprank"`       //当前排名 纯数字意义上的排名
	Topsubsection int                  `json:"topsubsection"` //当前段位
	Topranking    int                  `json:"topranking"`    //当前分段
	UpdateTime    int64                `json:"updatetime"`    //更新时间
	FightRecord   []*ActivityBossFight `json:"fightrecordid"` //战报集
	BestRecord    *ActivityBossFight   `json:"bestrecord"`    //最佳战报
}

type ActivityBossOtherInfo struct {
	Subsection   int                    `json:"subsection"`
	Ranking      int                    `json:"ranking"`
	Num          int                    `json:"num"`
	Players      []*ActivityBossTopNode `json:"players"`      //用户资料
	NowBaseScore int64                  `json:"nowbasescore"` //基准分
}

type ActivityBossSelfInfo struct {
	Subsection    int   `json:"subsection"`
	Ranking       int   `json:"ranking"`
	Score         int64 `json:"score"`
	NextBaseScore int64 `json:"nowbasescore"` //下一档分数
}

type ActivityBossRank struct {
	Top   []*ActivityBossTopNode   `json:"top"`   //传奇
	Other []*ActivityBossOtherInfo `json:"other"` //其他段位
}

type ActivityBossFight struct {
	FightId  int64  `json:"fight_id"`      // 战斗Id
	Side     int    `json:"side"`          // 1 进攻方 0 防守方
	Result   int    `json:"attack_result"` // 0 进攻方成功 其他防守方胜利
	Score    int64  `json:"score"`         // 积分增减
	Uid      int64  `json:"uid"`           // Uid
	IconId   int    `json:"icon"`          // 头像Id
	Portrait int    `json:"portrait"`      // 头像框
	Name     string `json:"name"`          // 名字
	Level    int    `json:"level"`         // 等级
	Fight    int64  `json:"fight"`         // 战力
	Time     int64  `json:"time"`          // 发生的时间
}

type ActivityBossNodeArr []*ActivityBossTopNode

func (s ActivityBossNodeArr) Len() int      { return len(s) }
func (s ActivityBossNodeArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ActivityBossNodeArr) Less(i, j int) bool {
	if s[i].Score == s[j].Score {
		return s[i].UpdateTime < s[j].UpdateTime
	}
	return s[i].Score > s[j].Score
}

type ActivityBossInfo struct {
	Id              int    `json:"id"`
	Period          int    `json:"period"`
	ActivityBossTop string `json:"activitybosstop"`
	StartTime       int64  `json:"starttime"`  //开始时间
	EndTime         int64  `json:"endtime"`    //结束时间
	RewardTime      int64  `json:"rewardtime"` //发奖时间
	HasReward       int    `json:"hasreward"`  //0没发 1已发

	Mu                     *sync.RWMutex
	activityBossTop        map[int64]*ActivityBossTopNode
	activityBossTopNodeArr ActivityBossNodeArr
	activityBossRank       ActivityBossRank
	bossFightInfo          *JS_FightInfo
	DataUpdate
}

type ActivityBossMgr struct {
	ActivityBossInfo map[int]*ActivityBossInfo
}

var activityBossMgr *ActivityBossMgr = nil

func GetActivityBossMgr() *ActivityBossMgr {
	if activityBossMgr == nil {
		activityBossMgr = new(ActivityBossMgr)
		activityBossMgr.ActivityBossInfo = make(map[int]*ActivityBossInfo)
	}

	return activityBossMgr
}

func (self *ActivityBossInfo) Encode() {
	self.Mu.RLock()
	defer self.Mu.RUnlock()
	self.ActivityBossTop = HF_JtoA(self.activityBossTop)
}

func (self *ActivityBossInfo) Decode() {
	self.Mu.Lock()
	defer self.Mu.Unlock()
	json.Unmarshal([]byte(self.ActivityBossTop), &self.activityBossTop)
}

// 存储数据库
func (self *ActivityBossMgr) Save() {
	for _, v := range self.ActivityBossInfo {
		v.Save()
	}
}

func (self *ActivityBossInfo) Save() {
	if self.Id == 0 || self.Period == 0 || self.EndTime < TimeServer().Unix() {
		return
	}

	self.Encode()
	self.Update(false)
}

func (self *ActivityBossMgr) NewActivityBossInfo(activityId int) *ActivityBossInfo {
	data := new(ActivityBossInfo)
	data.Mu = new(sync.RWMutex)
	data.Id = activityId
	data.activityBossTop = make(map[int64]*ActivityBossTopNode)
	return data
}

func (self *ActivityBossMgr) Run() {
	for i := ACT_BOSS_MIN; i < ACT_BOSS_MAX; i++ {
		self.ActivityBossInfo[i] = self.NewActivityBossInfo(i)
		self.ActivityBossInfo[i].GetData()
	}

	ticker := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-ticker.C:
			for _, v := range self.ActivityBossInfo {
				v.GetData()
				v.OnTimerState()
			}
		}
	}
}

func (self *ActivityBossInfo) GetData() {
	isOpen, _ := GetActivityMgr().JudgeOpenAllIndex(self.Id, self.Id)
	if !isOpen {
		return
	}

	activity := GetActivityMgr().GetActivity(self.Id)
	if activity == nil {
		return
	}
	self.StartTime = HF_CalTimeForConfig(activity.info.Start, "")
	self.EndTime = self.StartTime + int64(activity.info.Continued) + int64(activity.info.Show)
	self.RewardTime = self.StartTime + int64(activity.info.Continued)

	period := GetActivityMgr().getActN3(self.Id)
	if self.Period != period {
		self.Period = period
		queryStr := fmt.Sprintf("select * from `san_activityboss` where  `id` = %d and `period` = %d;", self.Id, self.Period)
		var msg ActivityBossInfo
		res := GetServer().DBUser.GetAllData(queryStr, &msg)

		if len(res) > 0 {
			self.ActivityBossTop = res[0].(*ActivityBossInfo).ActivityBossTop
			self.HasReward = res[0].(*ActivityBossInfo).HasReward
			self.Init("san_activityboss", self, false)
			self.Decode()
			self.MakeArr()
		} else {
			self.activityBossTop = make(map[int64]*ActivityBossTopNode, 0)
			self.HasReward = LOGIC_FALSE

			self.Encode()
			InsertTable("san_activityboss", self, 0, false)
			self.Init("san_activityboss", self, false)
		}

		nLen := len(GetCsvMgr().JJCRobotConfig)
		for i := 0; i < nLen; i++ {
			cfg := GetCsvMgr().JJCRobotConfig[i]
			if cfg.Type != self.Id {
				continue
			}
			self.bossFightInfo = self.GetRobot(cfg)
			break
		}
	}
	if self.activityBossTop == nil {
		self.activityBossTop = make(map[int64]*ActivityBossTopNode, 0)
	}
}

// 获取机器人信息
func (self *ActivityBossInfo) GetRobot(cfg *JJCRobotConfig) *JS_FightInfo {
	bossConfig := GetCsvMgr().GetActivityBossConfig(self.Id, self.Period)
	if bossConfig == nil {
		return nil
	}

	data := new(JS_FightInfo)
	data.Rankid = 0
	data.Uid = 0
	data.Uname = bossConfig.Name
	data.Iconid = 10000000 + bossConfig.BossId

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

	num := int(cfg.Fight[0] - cfg.Fight[1])
	if num <= 0 {
		num = 2
	}
	data.Deffight = cfg.Fight[1] + int64(HF_RandInt(1, num))

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
		hero.Color = cfg.NpcQuality
		hero.HeroKeyId = pos
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
	return data
}

func (self *ActivityBossInfo) GetBossFightInfo() *JS_FightInfo {
	if self.bossFightInfo == nil {
		nLen := len(GetCsvMgr().JJCRobotConfig)
		for i := 0; i < nLen; i++ {
			cfg := GetCsvMgr().JJCRobotConfig[i]
			if cfg.Type != self.Id {
				continue
			}
			self.bossFightInfo = self.GetRobot(cfg)
			break
		}
	}
	if self.bossFightInfo == nil {
		return nil
	} else {
		data := new(JS_FightInfo)
		HF_DeepCopy(data, self.bossFightInfo)
		return data
	}
}

func (self *ActivityBossMgr) GetBossFightInfo(id int) *JS_FightInfo {
	_, ok := self.ActivityBossInfo[id]
	if ok {
		return self.ActivityBossInfo[id].GetBossFightInfo()
	}
	return nil
}

func (self *ActivityBossInfo) OnTimerState() {
	//如果结束了，则直接看有没新的一期开启
	now := TimeServer().Unix()

	if self.HasReward == LOGIC_FALSE && self.RewardTime < now {

		strName := ""
		bossConfig := GetCsvMgr().GetActivityBossConfig(self.Id, self.Period)
		if bossConfig != nil {
			strName = bossConfig.Name
		}

		for i := 0; i < len(self.activityBossTopNodeArr); i++ {
			//配置
			config := GetCsvMgr().GetActivityBossRankConfig(self.Id, self.Period, self.activityBossTopNodeArr[i].Topsubsection, self.activityBossTopNodeArr[i].Topranking)
			if config == nil {
				LogError(fmt.Sprintf("配置不存在,期数%d,段位%d,分段%d", self.Period, self.activityBossTopNodeArr[i].Topsubsection, self.activityBossTopNodeArr[i].Topranking))
				continue
			}
			player := GetPlayerMgr().GetPlayer(self.activityBossTopNodeArr[i].Uid, true)
			if player == nil {
				continue
			}
			mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_ACTIVITYBOSS_RANK]
			if !ok {
				continue

			}

			pMail := player.GetModule("mail").(*ModMail)
			if pMail == nil {
				continue
			}
			var items []PassItem
			for i := 0; i < len(config.Item); i++ {
				if config.Item[i] == 0 {
					continue
				}

				if config.Num[i] == 0 {
					continue
				}
				items = append(items, PassItem{config.Item[i], config.Num[i]})
			}
			pMail.AddMail(1, 1, 0, mailConfig.Mailtitle, fmt.Sprintf(mailConfig.Mailtxt, strName, config.Name), GetCsvMgr().GetText("STR_SYS"), items, true, 0)
		}
		self.HasReward = LOGIC_TRUE
	}
}

func (self *ActivityBossInfo) MakeArr() {
	self.Mu.Lock()
	defer self.Mu.Unlock()

	self.activityBossTopNodeArr = make([]*ActivityBossTopNode, 0)
	for _, v := range self.activityBossTop {
		self.activityBossTopNodeArr = append(self.activityBossTopNodeArr, v)
	}
	sort.Sort(self.activityBossTopNodeArr)

	for i := 0; i < len(self.activityBossTopNodeArr); i++ {
		self.activityBossTopNodeArr[i].TopRank = i + 1
	}

	self.CalSubsect()
}

func (self *ActivityBossInfo) MakeRank() {
	config := GetCsvMgr().ActivityBossRankConfig[self.Id]
	if config == nil {
		return
	}

	configRank := config[self.Period]
	if configRank == nil {
		return
	}

	self.activityBossRank.Top = make([]*ActivityBossTopNode, 0)
	self.activityBossRank.Other = make([]*ActivityBossOtherInfo, 0)

	nowIndex := 0
	nowCount := 0
	baseScore := self.GetBaseScore()
	nowBase := baseScore + 1
	players := make([]*ActivityBossTopNode, 0)
	for i := 0; i < len(configRank); i++ {
		if configRank[i].Subsection == 1 {
			for j := nowIndex; j < len(self.activityBossTopNodeArr); j++ {
				if self.activityBossTopNodeArr[j].Topsubsection == configRank[i].Subsection && self.activityBossTopNodeArr[j].Topranking == configRank[i].Ranking {
					self.activityBossRank.Top = append(self.activityBossRank.Top, self.activityBossTopNodeArr[j])
					nowIndex++
					break
				}
			}
			continue
		}

		for j := nowIndex; j < len(self.activityBossTopNodeArr); j++ {
			if self.activityBossTopNodeArr[j].Topsubsection == configRank[i].Subsection && self.activityBossTopNodeArr[j].Topranking == configRank[i].Ranking {
				nowIndex++
				nowCount++
				if len(players) < 5 {
					players = append(players, self.activityBossTopNodeArr[j])
				}
			} else {
				break
			}
		}

		if baseScore > 0 {
			data := new(ActivityBossOtherInfo)
			data.Subsection = configRank[i].Subsection
			data.Ranking = configRank[i].Ranking
			data.Num = nowCount
			data.Players = players
			data.NowBaseScore = nowBase
			nowBase = (baseScore * configRank[i].Section / 10000) + 1
			self.activityBossRank.Other = append(self.activityBossRank.Other, data)
			nowCount = 0
			players = make([]*ActivityBossTopNode, 0)
		} else {
			break
		}
	}
}

//计算段位，这个函数不要加锁,独立功能
func (self *ActivityBossInfo) CalSubsect() {
	config := GetCsvMgr().ActivityBossRankConfig[self.Id]
	if config == nil {
		return
	}

	configRank := config[self.Period]
	if configRank == nil {
		return
	}

	//当前计算位置
	nowSubsect := 1
	nowRank := 1
	index := 0
	baseScore := self.GetBaseScore()
	for i := 0; i < len(configRank); i++ {
		if configRank[i].Contain > 0 {
			if configRank[i].Subsection == nowSubsect && configRank[i].Ranking == nowRank {
				//依然是同段位分布，分段位增加
				for j := index; j < index+configRank[i].Contain; j++ {
					if j >= len(self.activityBossTopNodeArr) {
						self.MakeRank()
						return
					}
					self.activityBossTopNodeArr[j].Topsubsection = configRank[i].Subsection
					self.activityBossTopNodeArr[j].Topranking = configRank[i].Ranking
				}
				index += configRank[i].Contain
			}
			nowRank++
		} else {
			if configRank[i].Subsection == nowSubsect && configRank[i].Ranking == nowRank {
				//依然是同段位分布，分段位增加
				scoreNow := baseScore * configRank[i].Section / 10000
				for j := index; j < len(self.activityBossTopNodeArr); j++ {
					if self.activityBossTopNodeArr[j].Score >= scoreNow {
						self.activityBossTopNodeArr[j].Topsubsection = configRank[i].Subsection
						self.activityBossTopNodeArr[j].Topranking = configRank[i].Ranking
						index++
					}
				}
				nowRank++
			} else {
				//升段位了
				nowSubsect++
				nowRank = 1
				i--
			}
		}
	}

	//生成排行
	self.MakeRank()
}

//获得传奇基准分
func (self *ActivityBossInfo) GetBaseScore() int64 {
	//根据配置计算传奇人数
	//最后1个人的分数则为传奇基准分
	config := GetCsvMgr().ActivityBossRankConfig[self.Id]
	if config == nil {
		return 0
	}

	configRank := config[self.Period]
	if configRank == nil {
		return 0
	}
	count := 0
	for i := 0; i < len(configRank); i++ {
		if configRank[i].Subsection == 1 {
			count += configRank[i].Contain
		}
	}

	if count-1 >= len(self.activityBossTopNodeArr) {
		return 0
	} else {
		return self.activityBossTopNodeArr[count-1].Score
	}
}

func (self *ActivityBossMgr) UpdatePoint(player *Player, score int64, bossId int, attack *JS_FightInfo, defend *JS_FightInfo, battleInfo BattleInfo) {

	fight := self.NewPvpFight(battleInfo.Id, attack, score)

	_, ok := self.ActivityBossInfo[bossId]
	if !ok {
		return
	}

	self.ActivityBossInfo[bossId].UpdatePoint(player, fight, score)

	data2 := BattleRecord{}
	data2.Level = 0
	data2.Side = 1
	data2.Time = TimeServer().Unix()
	data2.Id = battleInfo.Id
	data2.LevelID = battleInfo.LevelID
	//data2.Result = result
	data2.Type = BATTLE_TYPE_PVP
	data2.RandNum = battleInfo.Random
	data2.FightInfo[0] = attack
	data2.FightInfo[1] = defend

	HMSetRedisEx("san_activitybossbattleinfo", battleInfo.Id, &battleInfo, DAY_SECS*10)
	HMSetRedisEx("san_activitybossbattlerecord", data2.Id, &data2, DAY_SECS*10)
	return
}

func (self *ActivityBossMgr) NewPvpFight(FightID int64, enemy *JS_FightInfo, Score int64) *ActivityBossFight {

	p := &ActivityBossFight{}
	p.FightId = FightID
	p.Side = 1
	p.Score = Score
	if enemy != nil {
		p.Uid = enemy.Uid
		p.IconId = enemy.Iconid
		p.Portrait = enemy.Portrait
		p.Name = enemy.Uname
		p.Level = enemy.Level
		p.Fight = enemy.Deffight
	}
	p.Time = TimeServer().Unix()

	return p
}

func (self *ActivityBossInfo) UpdatePoint(player *Player, activityboss *ActivityBossFight, score int64) {
	//如果时间过了发奖时间,则不更新排行榜
	if TimeServer().Unix() > self.RewardTime {
		return
	}

	self.Mu.Lock()
	defer self.Mu.Unlock()

	baseScore := self.GetBaseScore()

	isNew := false

	info, ok := self.activityBossTop[player.Sql_UserBase.Uid]
	if ok {

		info.Name = player.Sql_UserBase.UName
		info.IconId = player.Sql_UserBase.IconId
		info.Portrait = player.Sql_UserBase.Portrait
		info.Level = player.Sql_UserBase.Level

		if len(info.FightRecord) >= 10 {
			info.FightRecord = info.FightRecord[1:]
		}
		info.FightRecord = append(info.FightRecord, activityboss)
		if info.Score > score {
			return
		}
		info.Score = score
		isNew = true
		info.BestRecord = activityboss
		info.UpdateTime = TimeServer().Unix()
	} else {
		data := new(ActivityBossTopNode)
		data.Score = score
		data.UpdateTime = TimeServer().Unix()
		data.Uid = player.Sql_UserBase.Uid
		data.Name = player.Sql_UserBase.UName
		data.IconId = player.Sql_UserBase.IconId
		data.Portrait = player.Sql_UserBase.Portrait
		data.TopRank = len(self.activityBossTopNodeArr) + 1
		data.Level = player.Sql_UserBase.Level
		data.FightRecord = append(data.FightRecord, activityboss)
		data.BestRecord = activityboss

		self.activityBossTop[player.Sql_UserBase.Uid] = data
		self.activityBossTopNodeArr = append(self.activityBossTopNodeArr, data)

		isNew = true
	}

	infoNow, okNow := self.activityBossTop[player.Sql_UserBase.Uid]
	if !okNow {
		return
	}

	for i := infoNow.TopRank - 2; i >= 0; i-- {
		if infoNow.Score > self.activityBossTopNodeArr[i].Score {
			self.activityBossTopNodeArr[i].TopRank++
			infoNow.TopRank--
			self.activityBossTopNodeArr.Swap(infoNow.TopRank-1, self.activityBossTopNodeArr[i].TopRank-1)
		} else {
			break
		}
	}

	if isNew || score >= baseScore {
		self.CalSubsect()
	}
	return
}
func (self *ActivityBossMgr) GetRank(player *Player, id int) {
	_, ok := self.ActivityBossInfo[id]
	if ok {
		self.ActivityBossInfo[id].GetRank(player)
	}
}

func (self *ActivityBossMgr) GetRecord(player *Player, id int, uid int64) {
	_, ok := self.ActivityBossInfo[id]
	if ok {
		self.ActivityBossInfo[id].GetRecord(player, uid)
	}
}

func (self *ActivityBossInfo) GetRank(player *Player) {
	self.Mu.RLock()
	defer self.Mu.RUnlock()

	var msg S2C_ActivityBossGetRank
	msg.Cid = "activitybossgetrank"
	msg.Id = self.Id
	msg.ActivityBossRank = self.activityBossRank
	msg.SelfInfo = self.GetPlayerInfo(player)
	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ActivityBossMgr) GetPlayerInfo(player *Player, id int) *ActivityBossSelfInfo {
	_, ok := self.ActivityBossInfo[id]
	if ok {
		return self.ActivityBossInfo[id].GetPlayerInfo(player)
	}
	return nil
}

func (self *ActivityBossInfo) GetPlayerInfo(player *Player) *ActivityBossSelfInfo {
	data := new(ActivityBossSelfInfo)

	if player == nil {
		return data
	}

	info, ok := self.activityBossTop[player.Sql_UserBase.Uid]
	if ok {
		data.Score = info.Score
		data.Ranking = info.Topranking
		data.Subsection = info.Topsubsection
		data.NextBaseScore = self.GetNextBaseScore(data.Subsection, data.Ranking)
	}

	return data
}

func (self *ActivityBossInfo) GetNextBaseScore(subsection int, ranking int) int64 {
	_, ok := GetCsvMgr().ActivityBossRankConfig[self.Id]
	if !ok {
		return 0
	}
	nextSubsection := 0
	nextRanking := 0
	for _, v := range GetCsvMgr().ActivityBossRankConfig[self.Id][self.Period] {
		if v.Subsection == subsection && v.Ranking == ranking {
			if nextSubsection == 0 {
				return 0
			} else {
				//非传奇分段
				if nextSubsection > 1 {
					config := GetCsvMgr().GetActivityBossRankConfig(self.Id, self.Period, nextSubsection, nextRanking)
					if config != nil {
						baseScore := self.GetBaseScore()
						return baseScore * config.Section / 10000
					} else {
						return 0
					}
				} else {
					for _, v := range self.activityBossTopNodeArr {
						if v.Topsubsection == nextSubsection && v.Topranking == nextRanking {
							return v.Score
						}
					}
				}
			}
		} else {
			nextSubsection = v.Subsection
			nextRanking = v.Ranking
		}
	}
	return 0
}

func (self *ActivityBossInfo) GetRecord(player *Player, uid int64) {
	self.Mu.RLock()
	defer self.Mu.RUnlock()

	var msg S2C_ActivityBossGetRecord
	msg.Cid = "activitybossgetrecord"
	msg.Target = uid
	_, ok := self.activityBossTop[uid]
	if ok {
		msg.FightRecord = self.activityBossTop[uid].FightRecord
	} else {
		msg.FightRecord = make([]*ActivityBossFight, 0)
	}
	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ActivityBossMgr) DeleteUserRecord(id int, uid int64) bool {

	_, ok := self.ActivityBossInfo[id]
	if !ok {
		return false
	}

	return self.ActivityBossInfo[id].DeleteUserRecord(uid)
}

func (self *ActivityBossInfo) DeleteUserRecord(uid int64) bool {

	_, ok := self.activityBossTop[uid]
	if ok {
		delete(self.activityBossTop, uid)
	}

	self.MakeArr()
	return true
}

func (self *ActivityBossMgr) Rename(player *Player) {
	for _, v := range self.ActivityBossInfo {
		v.Rename(player)
	}
}

func (self *ActivityBossInfo) Rename(player *Player) {
	self.Mu.Lock()
	defer self.Mu.Unlock()

	info, ok := self.activityBossTop[player.Sql_UserBase.Uid]
	if ok {
		info.Name = player.Sql_UserBase.UName
		info.IconId = player.Sql_UserBase.IconId
		info.Portrait = player.Sql_UserBase.Portrait

		for _, v := range info.FightRecord {
			v.Name = info.Name
			v.IconId = info.IconId
			v.Portrait = info.Portrait
		}

		if info.BestRecord != nil {
			info.BestRecord.Name = info.Name
			info.BestRecord.IconId = info.IconId
			info.BestRecord.Portrait = info.Portrait
		}
	}
}
