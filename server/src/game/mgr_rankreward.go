package game

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

type RankRewardPlayerInfo struct {
	Uid        int64  `json:"uid"`        //玩家UID
	KeyId      int    `json:"keyid"`      //活动key
	Name       string `json:"name"`       //名字
	IconId     int    `json:"icon"`       //头像
	Level      int    `json:"level"`      //等级
	Portrait   int    `json:"portrait"`   //头像框
	UpdateTime int64  `json:"updatetime"` //更新时间
	Score      int64  `json:"score"`      //当前成绩
	TopRank    int    `json:"toprank"`    //当前排名 纯数字意义上的排名
	Param      int64  `json:"param"`      //扩展参数
	HasReward  int    `json:"hasreward"`  //0没发 1已发
}

type RankRewardTopNode struct {
	Uid     int64  `json:"uid"`     //玩家UID
	KeyId   int    `json:"keyid"`   //活动key
	Score   int64  `json:"score"`   //当前成绩
	TopRank int    `json:"toprank"` //当前排名 纯数字意义上的排名
	Info    string `json:"info"`

	info *RankRewardPlayerInfo
	DataUpdate
}

type RankRewardNodeArr []*RankRewardTopNode

func (s RankRewardNodeArr) Len() int      { return len(s) }
func (s RankRewardNodeArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s RankRewardNodeArr) Less(i, j int) bool {
	if s[i].info.Score == s[j].info.Score {
		return s[i].info.UpdateTime < s[j].info.UpdateTime
	}
	return s[i].info.Score > s[j].info.Score
}

type RankRewardInfo struct {
	Id                   int
	KeyId                int
	StartTime            int64 //开始时间
	EndTime              int64 //结束时间
	RewardTime           int64 //发奖时间
	RankRewardTop        map[int64]*RankRewardTopNode
	RankRewardTopNodeArr RankRewardNodeArr
	TableName            string
	RewardCheck          bool
	Locker               *sync.RWMutex
}

//启动的时候初始化，以后只读不写
type RankRewardMgr struct {
	RankRewardInfo map[int]*RankRewardInfo
}

var rankRewardMgr *RankRewardMgr = nil

func GetRankRewardMgr() *RankRewardMgr {
	if rankRewardMgr == nil {
		rankRewardMgr = new(RankRewardMgr)
		rankRewardMgr.RankRewardInfo = make(map[int]*RankRewardInfo)
	}

	return rankRewardMgr
}

func (self *RankRewardTopNode) Encode() {
	self.Info = HF_JtoA(&self.info)
	self.KeyId = self.info.KeyId
	self.Score = self.info.Score
	self.TopRank = self.info.TopRank
}

func (self *RankRewardTopNode) Decode() {
	json.Unmarshal([]byte(self.Info), &self.info)
}

func (self *RankRewardMgr) Save() {
	for _, v := range self.RankRewardInfo {
		v.Save()
	}
}

func (self *RankRewardInfo) Save() {
	for _, v := range self.RankRewardTopNodeArr {
		v.Save()
	}
}

// 存储数据库
func (self *RankRewardTopNode) Save() {
	self.Encode()
	self.UpdateEx("keyid", self.KeyId)
}

func (self *RankRewardMgr) NewRankRewardInfoInfo(activityId int) *RankRewardInfo {
	data := new(RankRewardInfo)
	data.Locker = new(sync.RWMutex)
	data.Id = activityId
	data.RankRewardTop = make(map[int64]*RankRewardTopNode)
	data.TableName = fmt.Sprintf("rankreward%d", data.Id)
	return data
}

func (self *RankRewardMgr) Run() {
	self.RankRewardInfo[ACT_RANKREWARD_COST] = self.NewRankRewardInfoInfo(ACT_RANKREWARD_COST)
	//ticker := time.NewTicker(time.Minute * 1)
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			for _, v := range self.RankRewardInfo {
				v.GetData()
				v.OnTimerState()
			}
		}
	}
}

func (self *RankRewardInfo) GetData() {
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

	keyId := GetActivityMgr().getActN4(self.Id)*1000 + GetActivityMgr().getActN3(self.Id)
	if self.KeyId != keyId {
		self.KeyId = keyId
		self.RankRewardTop = make(map[int64]*RankRewardTopNode, 0)
		self.RankRewardTopNodeArr = make([]*RankRewardTopNode, 0)
		//初始化数据

		queryStr := fmt.Sprintf("select * from `%s` where `keyid` = %d;", self.TableName, self.KeyId)
		var msg RankRewardTopNode
		res := GetServer().DBUser.GetAllData(queryStr, &msg)

		for _, v := range res {
			data := v.(*RankRewardTopNode)
			data.Decode()
			data.Init(self.TableName, data, true)
			self.RankRewardTop[data.Uid] = data
		}
	}
	self.MakeArr()
}

func (self *RankRewardInfo) OnTimerState() {
	//主动发奖逻辑
	/*
		//如果结束了，则直接看有没新的一期开启
		now := TimeServer().Unix()
		if !self.RewardCheck && self.RewardTime < now {
			self.RewardCheck = true
			config := GetCsvMgr().GetRankRewardConfig(self.Id, self.KeyId)
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
		}
	*/
}

func (self *RankRewardInfo) MakeArr() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.RankRewardTopNodeArr = make([]*RankRewardTopNode, 0)
	for _, v := range self.RankRewardTop {
		self.RankRewardTopNodeArr = append(self.RankRewardTopNodeArr, v)
	}
	sort.Sort(self.RankRewardTopNodeArr)

	for i := 0; i < len(self.RankRewardTopNodeArr); i++ {
		self.RankRewardTopNodeArr[i].info.TopRank = i + 1
	}
}

func (self *RankRewardMgr) UpdateScore(id int, player *Player, score int64) {
	_, ok := self.RankRewardInfo[id]
	if !ok {
		player.SendErr(GetCsvMgr().GetText("STR_MOD_CONSUME_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	self.RankRewardInfo[id].UpdateScore(player, score)
}

func (self *RankRewardInfo) UpdateScore(player *Player, score int64) {
	//活动区间外,不更新排行榜
	if TimeServer().Unix() < self.StartTime || TimeServer().Unix() > self.RewardTime {
		return
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()

	info, ok := self.RankRewardTop[player.Sql_UserBase.Uid]
	if ok {
		info.info.Name = player.Sql_UserBase.UName
		info.info.IconId = player.Sql_UserBase.IconId
		info.info.Level = player.Sql_UserBase.Level
		info.info.Portrait = player.Sql_UserBase.Portrait
		info.info.Score += score
		info.info.UpdateTime = TimeServer().Unix()
	} else {
		data := new(RankRewardTopNode)
		data.Uid = player.Sql_UserBase.Uid
		data.info = new(RankRewardPlayerInfo)
		data.info.Score = score
		data.info.KeyId = self.KeyId
		data.info.Uid = player.Sql_UserBase.Uid
		data.info.Name = player.Sql_UserBase.UName
		data.info.IconId = player.Sql_UserBase.IconId
		data.info.Level = player.Sql_UserBase.Level
		data.info.Portrait = player.Sql_UserBase.Portrait
		data.info.Score += score
		data.info.TopRank = len(self.RankRewardTopNodeArr) + 1
		data.info.UpdateTime = TimeServer().Unix()

		self.RankRewardTop[player.Sql_UserBase.Uid] = data
		self.RankRewardTopNodeArr = append(self.RankRewardTopNodeArr, data)

		data.Encode()
		InsertTable(self.TableName, data, 0, true)
		data.Init(self.TableName, data, true)
	}

	infoNow, okNow := self.RankRewardTop[player.Sql_UserBase.Uid]
	if !okNow {
		return
	}

	for i := infoNow.info.TopRank - 2; i >= 0; i-- {
		if infoNow.info.Score > self.RankRewardTopNodeArr[i].info.Score {
			self.RankRewardTopNodeArr[i].info.TopRank++
			infoNow.info.TopRank--
			self.RankRewardTopNodeArr.Swap(infoNow.info.TopRank-1, self.RankRewardTopNodeArr[i].info.TopRank-1)
		} else {
			break
		}
	}
	return
}

func (self *RankRewardMgr) GetRank(player *Player, id int) {
	_, ok := self.RankRewardInfo[id]
	if ok {
		self.RankRewardInfo[id].GetRank(player)
	}
}

func (self *RankRewardInfo) GetRank(player *Player) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var msg S2C_GetRankRewardRank
	msg.Cid = "getrankrewardrank"
	msg.Id = self.Id
	if len(self.RankRewardTopNodeArr) > 20 {
		msg.Top = self.RankRewardTopNodeArr[:20]
	} else {
		msg.Top = self.RankRewardTopNodeArr
	}
	msg.SelfInfo = self.RankRewardTop[player.Sql_UserBase.Uid]
	msg.Config = GetCsvMgr().GetRankRewardConfig2C(self.Id, self.KeyId)
	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *RankRewardMgr) GetReward(player *Player, id int) {
	_, ok := self.RankRewardInfo[id]
	if ok {
		self.RankRewardInfo[id].GetReward(player)
	}
}

func (self *RankRewardInfo) GetReward(player *Player) {
	//活动区间外,不更新排行榜
	if TimeServer().Unix() < self.RewardTime || TimeServer().Unix() > self.EndTime {
		return
	}

	self.Locker.RLock()
	defer self.Locker.RUnlock()

	info, ok := self.RankRewardTop[player.Sql_UserBase.Uid]
	if !ok {
		player.SendErr(GetCsvMgr().GetText("STR_MOD_HORSE_DATA_DOES_NOT_EXIST"))
		return
	}
	if info.info.HasReward == LOGIC_TRUE {
		player.SendErr(GetCsvMgr().GetText("STR_MOD_ACTIVITY_AWARD_FOR_THE_EVENT_HAS"))
		return
	}

	config := GetCsvMgr().GetRankRewardConfig(self.Id, self.KeyId, info.info.TopRank)
	if config == nil {
		player.SendErr(GetCsvMgr().GetText("STR_MOD_CONSUME_CONFIGURATION_DOES_NOT_EXIST"))
		return
	}
	//发送奖励并同步
	getItems := player.AddObjectLst(config.Item, config.Num, "排行奖励", self.Id, info.info.TopRank, 0)
	info.info.HasReward = LOGIC_TRUE

	var msg S2C_GetRankRewardReward
	msg.Cid = "getrankrewardreward"
	msg.Id = self.Id
	msg.GetItems = getItems
	msg.SelfInfo = info.info
	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *RankRewardMgr) Rename(player *Player) {
	for _, v := range self.RankRewardInfo {
		v.Rename(player)
	}
}

func (self *RankRewardInfo) Rename(player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	info, ok := self.RankRewardTop[player.Sql_UserBase.Uid]
	if ok {
		info.info.Name = player.Sql_UserBase.UName
		info.info.IconId = player.Sql_UserBase.IconId
		info.info.Portrait = player.Sql_UserBase.Portrait
	}
}
