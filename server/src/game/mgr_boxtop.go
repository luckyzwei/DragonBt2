package game

import (
	"fmt"
	"sort"
	"sync"
)

//! 宝箱排行
type JS_TopBox struct {
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Camp  int    `json:"camp"`
	BoxId int    `json:"boxid"`
	Num   int    `json:"num"`
	Icon  int    `json:"icon"`
	Vip   int    `json:"vip"`
}
type lstTopBox []*JS_TopBox

func (s lstTopBox) Len() int           { return len(s) }
func (s lstTopBox) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lstTopBox) Less(i, j int) bool { return s[i].Num > s[j].Num }

type TopBoxMgr struct {
	TopBox lstTopBox //! 功勋排行
	Lock   *sync.RWMutex
}

var topboxmgrsingleton *TopBoxMgr = nil

func GetTopBoxMgr() *TopBoxMgr {
	if topboxmgrsingleton == nil {
		topboxmgrsingleton = new(TopBoxMgr)
		topboxmgrsingleton.Lock = new(sync.RWMutex)
	}

	return topboxmgrsingleton
}

func (self *TopBoxMgr) Refresh() {
	for i := 0; i < len(self.TopBox); i++ {
		csv, ok := GetCsvMgr().Data["Gemsweeper_rank"][i+1]
		if !ok {
			break
		}
		player := GetPlayerMgr().GetPlayer(self.TopBox[i].Uid, true)
		if player == nil {
			continue
		}
		lstItem := make([]PassItem, 0)
		for j := 0; j < 4; j++ {
			itemid := HF_Atoi(csv[fmt.Sprintf("item_id%d", j+1)])
			if itemid == 0 {
				break
			}
			lstItem = append(lstItem, PassItem{itemid, HF_Atoi(csv[fmt.Sprintf("num%d", j+1)])})
		}
		player.GetModule("mail").(*ModMail).AddMail(1, 1, 0, csv["mail_title"], csv["mail_txt"], GetCsvMgr().GetText("STR_SYS"), lstItem, false, 0)
	}

	self.TopBox = make(lstTopBox, 0)
}

//!
func (self *TopBoxMgr) Add(player *Player, boxid int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	find := false
	for i := 0; i < len(self.TopBox); i++ {
		if self.TopBox[i].Uid == player.Sql_UserBase.Uid {
			self.TopBox[i].Num++
			self.TopBox[i].BoxId = boxid
			self.TopBox[i].Vip = player.Sql_UserBase.Vip
			if i > 0 && self.TopBox[i].Num < self.TopBox[i-1].Num {
				return
			}
			find = true
			break
		}
	}

	if !find {
		node := new(JS_TopBox)
		node.Uid = player.Sql_UserBase.Uid
		node.Name = player.Sql_UserBase.UName
		node.Camp = player.Sql_UserBase.Camp
		node.Num = 1
		node.Icon = player.Sql_UserBase.IconId
		node.BoxId = boxid
		node.Vip = player.Sql_UserBase.Vip
		self.TopBox = append(self.TopBox, node)
	}

	sort.Sort(lstTopBox(self.TopBox))
}

func (self *TopBoxMgr) Get() []*JS_TopBox {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	return self.TopBox[:HF_MinInt(20, len(self.TopBox))]
}
