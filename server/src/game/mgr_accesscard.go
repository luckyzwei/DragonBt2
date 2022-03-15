package game

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

const (
	MAX_ACCESSCARD_RECORD_SIZE = 20
)

type AccessCardRecord struct {
	Uid      int64  `json:"uid"`      //玩家UID
	Name     string `json:"string"`   //玩家NAME
	GetTime  int64  `json:"gettime"`  //获取时间
	Type     int    `json:"type"`     //0领取最终奖励记录      1点数达到记录
	Rank     int    `json:"Rank"`     //第几个 type=0生效
	Point    int    `json:"point"`    //点数 type=1生效
	ItemName string `json:"itemname"` //物品名称 type=0生效
}

type AccessCardRecordTopNode struct {
	Uid        int64  `json:"uid"`        //玩家UID
	Name       string `json:"string"`     //玩家NAME
	TopRank    int    `json:"toprank"`    //当前排名
	Point      int    `json:"point"`      //点数
	UpdateTime int64  `json:"updatetime"` //更新时间
}

type AccessCardRecordTopNodeArr []*AccessCardRecordTopNode

func (s AccessCardRecordTopNodeArr) Len() int      { return len(s) }
func (s AccessCardRecordTopNodeArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s AccessCardRecordTopNodeArr) Less(i, j int) bool {
	if s[i].Point == s[j].Point {
		return s[i].UpdateTime < s[j].UpdateTime
	}
	return s[i].Point > s[j].Point
}

type AccessCardRecordMgr struct {
	Id               int    `json:"id"`
	NGroup           int    `json:"ngroup"`
	AccessCardRecord string `json:"accesscardrecord"`
	Rank             int    `json:"rank"`
	AccessCardTop    string `json:"accesscardtop"`
	StartTime        int64  `json:"starttime"`  //开始时间
	EndTime          int64  `json:"endtime"`    //结束时间
	RewardTime       int64  `json:"rewardtime"` //发奖时间
	HasReward        int    `json:"hasreward"`  //0没发 1已发

	Mu                         *sync.RWMutex
	accessCardRecord           []*AccessCardRecord
	accessCardTop              map[int64]*AccessCardRecordTopNode
	accessCardRecordTopNodeArr AccessCardRecordTopNodeArr
	DataUpdate
}

var accessCardRecordMgr *AccessCardRecordMgr = nil

func GetAccessCardRecordMgr() *AccessCardRecordMgr {
	if accessCardRecordMgr == nil {
		accessCardRecordMgr = new(AccessCardRecordMgr)
		accessCardRecordMgr.NGroup = 0
		accessCardRecordMgr.accessCardRecord = make([]*AccessCardRecord, 0)
		accessCardRecordMgr.Mu = new(sync.RWMutex)
	}

	return accessCardRecordMgr
}

func (self *AccessCardRecordMgr) CheckRefresh() {

}

func (self *AccessCardRecordMgr) Encode() {
	self.AccessCardRecord = HF_JtoA(self.accessCardRecord)
	self.AccessCardTop = HF_JtoA(self.accessCardTop)
}

func (self *AccessCardRecordMgr) Decode() {
	json.Unmarshal([]byte(self.AccessCardRecord), &self.accessCardRecord)
	json.Unmarshal([]byte(self.AccessCardTop), &self.accessCardTop)
}

// 存储数据库
func (self *AccessCardRecordMgr) Save() {

	if self.Id == 0 || self.NGroup == 0 || self.EndTime < TimeServer().Unix() {
		return
	}

	self.Mu.RLock()
	defer self.Mu.RUnlock()

	self.Encode()
	self.Update(true)
}

func (self *AccessCardRecordMgr) Run() {

	self.GetData()

	ticker := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-ticker.C:
			self.OnTimerState()
		}
	}
}

func (self *AccessCardRecordMgr) OnTimerState() {
	//如果结束了，则直接看有没新的一期开启
	now := TimeServer().Unix()
	if self.EndTime < now {
		self.GetData()
		return
	}

	if self.HasReward == LOGIC_FALSE && self.RewardTime < now {
		for i := 0; i < len(self.accessCardRecordTopNodeArr); i++ {
			config := GetCsvMgr().GetAccessTaskConfig(self.NGroup, i+1)
			if config == nil {
				//发放完毕，后面的没了
				break
			}
			player := GetPlayerMgr().GetPlayer(self.accessCardRecordTopNodeArr[i].Uid, true)
			if player == nil {
				continue
			}
			pMail := player.GetModule("mail").(*ModMail)
			if pMail == nil {
				continue
			}
			itemMap := make(map[int]*Item)
			content := config.MailTxt[2] //默认没有额外奖励
			text := fmt.Sprintf(content, i+1, self.accessCardRecordTopNodeArr[i].Point)
			AddItemMapHelper(itemMap, config.NormalAward, config.NormalNum)
			if config.ExtraAward[0] > 0 {
				content = config.MailTxt[1]
				text = fmt.Sprintf(content, i+1, self.accessCardRecordTopNodeArr[i].Point)
				if self.accessCardRecordTopNodeArr[i].Point < config.NeedPoint {
					content = config.MailTxt[0]
					text = fmt.Sprintf(content, i+1, self.accessCardRecordTopNodeArr[i].Point, config.NeedPoint)
				} else {
					AddItemMapHelper(itemMap, config.ExtraAward, config.ExtraNum)
				}
			}
			var items []PassItem
			for _, value := range itemMap {
				if value.ItemId == 0 {
					continue
				}

				if value.ItemNum == 0 {
					continue
				}
				items = append(items, PassItem{value.ItemId, value.ItemNum})
			}
			pMail.AddMail(1, 1, 0, fmt.Sprintf(config.MailTitle, i+1), text, GetCsvMgr().GetText("STR_SYS"), items, true, 0)
		}
		self.HasReward = LOGIC_TRUE
	}
}

func (self *AccessCardRecordMgr) GetData() {

	isOpen, index := GetActivityMgr().JudgeOpenAllIndex(ACT_ACCESSCARD_MIN, ACT_ACCESSCARD_MAX)
	if !isOpen {
		return
	}

	activity := GetActivityMgr().GetActivity(index)
	if activity == nil {
		return
	}

	self.Id = GetActivityMgr().getActN3(index)
	self.NGroup = GetActivityMgr().getActN4(index)

	queryStr := fmt.Sprintf("select * from `san_accesscard` where  `id` = %d and `ngroup` = %d;", self.Id, self.NGroup)
	var msg AccessCardRecordMgr
	res := GetServer().DBUser.GetAllData(queryStr, &msg)

	if len(res) > 0 {
		self.AccessCardRecord = res[0].(*AccessCardRecordMgr).AccessCardRecord
		self.AccessCardTop = res[0].(*AccessCardRecordMgr).AccessCardTop
		self.HasReward = res[0].(*AccessCardRecordMgr).HasReward
		self.Init("san_accesscard", self, false)
		self.Decode()
		self.MakeArr()
		startday := HF_Atoi(activity.info.Start)
		accessCardRecordMgr.StartTime = GetServer().GetOpenServer() + int64(startday-1)*86400
		accessCardRecordMgr.EndTime = accessCardRecordMgr.StartTime + int64(activity.info.Continued) + int64(activity.info.Show)
		for _, v := range GetCsvMgr().AccessRankConfig {
			accessCardRecordMgr.RewardTime = accessCardRecordMgr.StartTime + int64(activity.info.Continued) + v.RankOverTime
			break
		}
	} else {
		accessCardRecordMgr.accessCardRecord = make([]*AccessCardRecord, 0)
		accessCardRecordMgr.accessCardTop = make(map[int64]*AccessCardRecordTopNode, 0)
		accessCardRecordMgr.accessCardRecordTopNodeArr = make([]*AccessCardRecordTopNode, 0)
		accessCardRecordMgr.HasReward = LOGIC_FALSE

		startday := HF_Atoi(activity.info.Start)
		accessCardRecordMgr.StartTime = GetServer().GetOpenServer() + int64(startday-1)*86400
		accessCardRecordMgr.EndTime = accessCardRecordMgr.StartTime + int64(activity.info.Continued) + int64(activity.info.Show)
		for _, v := range GetCsvMgr().AccessRankConfig {
			accessCardRecordMgr.RewardTime = accessCardRecordMgr.StartTime + int64(activity.info.Continued) + v.RankOverTime
			break
		}
		self.Encode()
		InsertTable("san_accesscard", self, 0, false)
		self.Init("san_accesscard", self, false)
	}

	if accessCardRecordMgr.accessCardTop == nil {
		accessCardRecordMgr.accessCardTop = make(map[int64]*AccessCardRecordTopNode, 0)
	}
}

func (self *AccessCardRecordMgr) AddRecord(player *Player, num int, itemName string, nType int) int {
	if player == nil {
		return 0
	}

	self.Mu.Lock()
	defer self.Mu.Unlock()

	data := new(AccessCardRecord)
	data.Uid = player.GetUid()
	data.Name = player.GetName()
	data.GetTime = TimeServer().Unix()
	if nType == 0 {
		self.Rank = self.Rank + 1
		data.Rank = self.Rank
	}
	data.Point = num
	data.Type = nType
	data.ItemName = itemName

	if len(self.accessCardRecord) >= MAX_ACCESSCARD_RECORD_SIZE {
		self.accessCardRecord = append(self.accessCardRecord[1:], data)
	} else {
		self.accessCardRecord = append(self.accessCardRecord, data)
	}

	return self.Rank
}

func (self *AccessCardRecordMgr) GetRecord() []*AccessCardRecord {
	self.Mu.RLock()
	defer self.Mu.RUnlock()

	return self.accessCardRecord
}

func (self *AccessCardRecordMgr) IsCanGet() bool {
	if TimeServer().Unix() > self.RewardTime {
		return false
	}

	return true
}

func (self *AccessCardRecordMgr) MakeArr() {
	self.Mu.Lock()
	defer self.Mu.Unlock()

	self.accessCardRecordTopNodeArr = make([]*AccessCardRecordTopNode, 0)
	for _, v := range self.accessCardTop {
		self.accessCardRecordTopNodeArr = append(self.accessCardRecordTopNodeArr, v)
	}
	sort.Sort(self.accessCardRecordTopNodeArr)

	for i := 0; i < len(self.accessCardRecordTopNodeArr); i++ {
		self.accessCardRecordTopNodeArr[i].TopRank = i + 1
	}
}

func (self *AccessCardRecordMgr) UpdatePoint(player *Player, point int) {
	//如果时间过了发奖时间,则不更新排行榜
	if TimeServer().Unix() > self.RewardTime {
		return
	}

	self.Mu.Lock()
	defer self.Mu.Unlock()

	info, ok := self.accessCardTop[player.Sql_UserBase.Uid]
	if ok {
		info.Point = point
		info.UpdateTime = TimeServer().Unix()
		for i := info.TopRank - 2; i >= 0; i-- {
			if info.Point > self.accessCardRecordTopNodeArr[i].Point {
				self.accessCardRecordTopNodeArr[i].TopRank++
				info.TopRank--
				self.accessCardRecordTopNodeArr.Swap(info.TopRank-1, self.accessCardRecordTopNodeArr[i].TopRank-1)
			} else {
				break
			}
		}
	} else {
		data := new(AccessCardRecordTopNode)
		data.Point = point
		data.UpdateTime = TimeServer().Unix()
		data.Uid = player.Sql_UserBase.Uid
		data.Name = player.Sql_UserBase.UName
		data.TopRank = len(self.accessCardRecordTopNodeArr) + 1
		self.accessCardTop[player.Sql_UserBase.Uid] = data
		self.accessCardRecordTopNodeArr = append(self.accessCardRecordTopNodeArr, data)

		for i := data.TopRank - 2; i >= 0; i-- {
			if data.Point > self.accessCardRecordTopNodeArr[i].Point {
				self.accessCardRecordTopNodeArr[i].TopRank++
				data.TopRank--
				self.accessCardRecordTopNodeArr.Swap(data.TopRank-1, self.accessCardRecordTopNodeArr[i].TopRank-1)
			} else {
				break
			}
		}
	}
}

func (self *AccessCardRecordMgr) GetRank(player *Player) {
	self.Mu.RLock()
	defer self.Mu.RUnlock()

	var msg S2C_AccessGetRank
	msg.Cid = "accessgetrank"
	if len(self.accessCardRecordTopNodeArr) > 50 {
		msg.Rank = self.accessCardRecordTopNodeArr[0:50]
	} else {
		msg.Rank = self.accessCardRecordTopNodeArr
	}
	msg.Self = self.accessCardTop[player.GetUid()]
	msg.StartTime = self.StartTime
	msg.EndTime = self.EndTime
	msg.RewardTime = self.RewardTime
	msg.HasReward = self.HasReward
	player.SendMsg(msg.Cid, HF_JtoB(&msg))
}
