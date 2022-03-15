package game

import (
	"fmt"
	"sync"
	"time"
)

const (
	BEST_SUBSECTION = 1 //传奇段位
)

const (
	CROSSARENA_UPDATE_RANK_TIME = 300
)

type JS_CrossArenaBattleInfo struct {
	Id           int    `json:"id"`           //! 自增Id
	FightId      int64  `json:"fightid"`      //! 战斗Id
	RecordType   int    `json:"recordtype"`   //! 战报类型
	BattleInfo   string `json:"battleinfo"`   //! 简报
	BattleRecord string `json:"battlerecord"` //! 详细战报
	UpdateTime   int64  `json:"updatetime"`   //! 插入时间

	DataUpdate
}

type Js_CrossArenaUser struct {
	Uid         int64         `json:"uid"`
	SvrId       int           `json:"svrid"`
	SvrName     string        `json:"svrname"`
	Subsection  int           `json:"subsection"` //大段位
	Class       int           `json:"class"`      //小段位
	UName       string        `json:"uname"`
	Level       int           `json:"level"`
	Vip         int           `json:"vip"`
	Icon        int           `json:"icon"`
	Portrait    int           `json:"portrait"`
	Fight       int64         `json:"fight"`
	Robot       int           `json:"robot"`
	FightRecord []*ArenaFight `json:"fightrecord"` //战报集
}

// 神将活动单服务器配置
type CrossArenaMgr struct {
	KeyId      int   `json:"keyid"`      //! 活动期数
	StartTime  int64 `json:"starttime"`  //! 活动开始时间
	EndTime    int64 `json:"endtime"`    //! 活动结束时间
	ShowTime   int64 `json:"showtime"`   //! 活动展示时间
	IsReward   int   `json:"isreward"`   //！是否发奖
	UpdateTime int64 `json:"updatetime"` //! 请求时间

	Locker     *sync.RWMutex                        //! 数据锁
	topNew     map[int]map[int][]*Js_CrossArenaUser //! 排行
	rankNewMap map[int64]*Js_CrossArenaUser         //! 这个服务器里所有人的排名
	top        []*Js_CrossArenaUser                 //! 通知用

	DataUpdate
}

var crossarenaServ *CrossArenaMgr = nil

func GetCrossArenaMgr() *CrossArenaMgr {
	if crossarenaServ == nil {
		crossarenaServ = new(CrossArenaMgr)
		crossarenaServ.rankNewMap = make(map[int64]*Js_CrossArenaUser, 0)
		crossarenaServ.topNew = make(map[int]map[int][]*Js_CrossArenaUser, 0)
		crossarenaServ.Locker = new(sync.RWMutex)
		crossarenaServ.UpdateTime = time.Now().Unix() + 30
	}
	return crossarenaServ
}

func (self *CrossArenaMgr) OnTimer() {
	self.GetData()

	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		if self.UpdateTime < TimeServer().Unix() {
			self.ReqCrossArenaRank()
			self.UpdateTime = TimeServer().Unix() + CROSSARENA_UPDATE_RANK_TIME
			//LogDebug("开始向主服务器请求数据:", time.Now().Format(DATEFORMAT))
		}
	}
	ticker.Stop()
}

func (self *CrossArenaMgr) GetData() {
	activity := GetActivityMgr().GetActivity(ACT_AREAN_CROSS_SERVER)
	if activity == nil {
		return
	}

	self.KeyId = GetActivityMgr().getActN3(ACT_AREAN_CROSS_SERVER)
	self.StartTime = HF_CalTimeForConfig(activity.info.Start, "")
	self.EndTime = self.StartTime + int64(activity.info.Continued)
	self.ShowTime = self.StartTime + int64(activity.info.Continued) + int64(activity.info.Show)

	queryStr := fmt.Sprintf("select * from `san_crossarena` where  `keyid` = %d ;", self.KeyId)
	var msg CrossArenaMgr
	res := GetServer().DBUser.GetAllData(queryStr, &msg)

	if len(res) > 0 {
		self.IsReward = res[0].(*CrossArenaMgr).IsReward
		self.Init("san_crossarena", self, false)
	} else {
		self.rankNewMap = make(map[int64]*Js_CrossArenaUser, 0)
		self.topNew = make(map[int]map[int][]*Js_CrossArenaUser, 0)
		self.IsReward = LOGIC_FALSE
		InsertTable("san_crossarena", self, 0, false)
		self.Init("san_crossarena", self, false)
	}
}

//! 从中心服务器请求限时神将排行榜数据
func (self *CrossArenaMgr) ReqCrossArenaRank() {
	activity := GetActivityMgr().GetActivity(ACT_AREAN_CROSS_SERVER)
	if activity == nil {
		return
	}
	self.StartTime = HF_CalTimeForConfig(activity.info.Start, "")
	self.EndTime = self.StartTime + int64(activity.info.Continued)
	self.ShowTime = self.StartTime + int64(activity.info.Continued) + int64(activity.info.Show)

	isOpen, _ := GetActivityMgr().JudgeOpenAllIndex(ACT_AREAN_CROSS_SERVER, ACT_AREAN_CROSS_SERVER)
	if !isOpen {
		return
	}
	keyId := activity.getTaskN3()
	if keyId == 0 {
		return
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()
	if self.KeyId != keyId {
		self.KeyId = keyId
		queryStr := fmt.Sprintf("select * from `san_crossarena` where  `keyid` = %d ;", self.KeyId)
		var msg CrossArenaMgr
		res := GetServer().DBUser.GetAllData(queryStr, &msg)

		if len(res) > 0 {
			self.IsReward = res[0].(*CrossArenaMgr).IsReward
			self.Init("san_crossarena", self, false)
		} else {
			self.rankNewMap = make(map[int64]*Js_CrossArenaUser, 0)
			self.topNew = make(map[int]map[int][]*Js_CrossArenaUser, 0)
			self.IsReward = LOGIC_FALSE
			InsertTable("san_crossarena", self, 0, false)
			self.Init("san_crossarena", self, false)
		}
	}

	res := GetMasterMgr().MatchCrossArenaGetAllRank(keyId)
	if res != nil {
		self.topNew = res.RankInfo
		for _, subUser := range self.topNew {
			for _, classUser := range subUser {
				for _, user := range classUser {
					if user.SvrId == GetServer().Con.ServerId {
						self.rankNewMap[user.Uid] = user
					}
				}
			}
		}
		self.UpdateTime = TimeServer().Unix() + CROSSARENA_UPDATE_RANK_TIME
		self.MakeArr()

		now := TimeServer().Unix()
		if self.IsReward == LOGIC_FALSE && now > self.EndTime {
			//发送排行奖励
			self.IsReward = LOGIC_TRUE
			for _, v := range self.rankNewMap {
				player := GetPlayerMgr().GetPlayer(v.Uid, true)
				if player == nil {
					continue
				}
				pMail := player.GetModule("mail").(*ModMail)
				if pMail == nil {
					continue
				}
				itemMap := GetCsvMgr().GetCrossArenaRewardBySubsection(v.Subsection)
				if len(itemMap) > 0 {
					mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_CROSSARENA_SUB]
					if ok {
						itemLst := make([]PassItem, 0)
						for _, v := range itemMap {
							itemLst = append(itemLst, PassItem{ItemID: v.ItemId, Num: v.ItemNum})
						}
						TxtName := GetCsvMgr().GetCrossArenaSubsectionName(v.Subsection, v.Class)
						pMail.AddMail(1, 1, 0, mailConfig.Mailtitle, fmt.Sprintf(mailConfig.Mailtxt, TxtName), GetCsvMgr().GetText("STR_SYS"), itemLst, false, 0)
					}
				}
				configRank := GetCsvMgr().GetCrossArenaRewardByRank(v.Subsection, v.Class)
				if configRank != nil {
					mailConfig, ok := GetCsvMgr().MailConfig[MAIL_ID_CROSSARENA_RANK]
					if ok {
						itemLst := make([]PassItem, 0)
						for i := 0; i < len(configRank.Item); i++ {
							itemLst = append(itemLst, PassItem{ItemID: configRank.Item[i], Num: configRank.Num[i]})
						}
						pMail.AddMail(1, 1, 0, mailConfig.Mailtitle, fmt.Sprintf(mailConfig.Mailtxt, v.Class), GetCsvMgr().GetText("STR_SYS"), itemLst, false, 0)
					}
				}
			}
		}
	}
}

func (self *CrossArenaMgr) MakeArr() {

	self.top = make([]*Js_CrossArenaUser, 0)
	for _, config := range GetCsvMgr().CrossArenaSubsection {
		if config.Subsection != BEST_SUBSECTION {
			return
		}
		_, okSubsection := self.topNew[config.Subsection]
		if !okSubsection {
			continue
		}
		_, okClass := self.topNew[config.Subsection][config.Class]
		if !okClass {
			continue
		}
		self.top = append(self.top, self.topNew[config.Subsection][config.Class]...)
	}
}

func (self *CrossArenaMgr) GetSendInfo(player *Player) ([]*Js_CrossArenaUser, *Js_CrossArenaUser) {
	if player != nil {
		self.Locker.RLock()
		defer self.Locker.RUnlock()
		return self.top, self.rankNewMap[player.Sql_UserBase.Uid]
	}
	return nil, nil
}

func (self *CrossArenaMgr) PlayerToCrossArena(player *Player) *Js_CrossArenaUser {
	data := new(Js_CrossArenaUser)
	data.Uid = player.Sql_UserBase.Uid
	data.SvrId = GetServer().Con.ServerId
	data.SvrName = GetServer().Con.ServerName
	data.UName = player.Sql_UserBase.UName
	data.Level = player.Sql_UserBase.Level
	data.Vip = player.Sql_UserBase.Vip
	data.Icon = player.Sql_UserBase.IconId
	data.Portrait = player.Sql_UserBase.Portrait

	return data
}

func (self *CrossArenaMgr) GetStartLv() (int, int) {
	size := len(GetCsvMgr().CrossArenaSubsection)
	if size > 0 {
		return GetCsvMgr().CrossArenaSubsection[size-1].Subsection, GetCsvMgr().CrossArenaSubsection[size-1].Class
	}
	return 0, 0
}

func (self *CrossArenaMgr) AddInfo(player *Player, fightInfo *JS_FightInfo) *Js_CrossArenaUser {
	data := self.PlayerToCrossArena(player)
	data.Subsection, data.Class = self.GetStartLv()
	data.Fight = fightInfo.Deffight
	data.SvrId = GetServer().Con.ServerId
	res := GetMasterMgr().MatchCrossArenaAdd(self.KeyId, data, fightInfo)
	if res != nil {
		self.topNew = res.RankInfo
		for _, subUser := range self.topNew {
			for _, classUser := range subUser {
				for _, user := range classUser {
					if user.SvrId == GetServer().Con.ServerId {
						self.rankNewMap[user.Uid] = user
					}
				}
			}
		}
		if res.SelfInfo != nil {
			self.rankNewMap[res.SelfInfo.Uid] = res.SelfInfo
		}
		self.MakeArr()
		return res.SelfInfo
	}
	return nil
}

func (self *CrossArenaMgr) GetDefenceList(player *Player) ([]*Js_CrossArenaUser, []*JS_FightInfo, int) {

	res := GetMasterMgr().MatchCrossArenaGetDefence(self.KeyId, player)
	if res != nil {
		return res.Info, res.FightInfo, res.RetCode
	}
	return nil, nil, 0
}

func (self *CrossArenaMgr) GetInfo(uid int64) (*Js_CrossArenaUser, *JS_FightInfo, *JS_LifeTreeInfo) {
	res := GetMasterMgr().MatchCrossArenaGetInfo(self.KeyId, uid)
	if res != nil {
		return res.Info, res.FightInfo, res.LifeTreeInfo
	}
	return nil, nil, nil
}

func (self *CrossArenaMgr) FightEnd(player *Player, attack *JS_FightInfo, defend *JS_FightInfo, battleInfo BattleInfo) *RPC_CrossArenaActionRes {
	//通知中心服结果，并返回相应信息
	res := GetMasterMgr().MatchCrossArenaFightEnd(self.KeyId, attack, defend, battleInfo)
	if res != nil {
		self.topNew = res.RankInfo
		for _, subUser := range self.topNew {
			for _, classUser := range subUser {
				for _, user := range classUser {
					if user.SvrId == GetServer().Con.ServerId {
						self.rankNewMap[user.Uid] = user
					}
				}
			}
		}
		if res.SelfInfo != nil {
			self.rankNewMap[res.SelfInfo.Uid] = res.SelfInfo
		}
		self.UpdateTime = TimeServer().Unix() + CROSSARENA_UPDATE_RANK_TIME
		self.MakeArr()
	}
	return res
}

func (self *CrossArenaMgr) GetBattleInfo(key int64) *BattleInfo {
	res := GetMasterMgr().MatchCrossArenaGetBattleInfo(key)
	if res != nil {
		return res.BattleInfo
	}
	return nil
}

func (self *CrossArenaMgr) GetBattleRecord(key int64) *BattleRecord {
	res := GetMasterMgr().MatchCrossArenaGetBattleRecord(key)
	if res != nil {
		return res.BattleRecord
	}
	return nil
}

func (self *CrossArenaMgr) UpdateInfo(data *Js_CrossArenaUser) {
	if data == nil {
		return
	}
	if data.SvrId == GetServer().Con.ServerId {
		self.Locker.Lock()
		self.rankNewMap[data.Uid] = data
		self.Locker.Unlock()
		player := GetPlayerMgr().GetPlayer(data.Uid, false)
		if player != nil {
			var msg S2C_CrossArenaUpdate
			msg.Cid = "crossarenaupdateinfo"
			msg.SelfInfo = data
			player.SendMsg(msg.Cid, HF_JtoB(&msg))
		}
	}
}

// 存储数据库
func (self *CrossArenaMgr) Save() {
	self.Update(true)
}
