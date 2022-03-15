package game

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

const (
	maxRank = 50
)

type DialInfo struct {
	Name string `json:"name"` //! 名字
	Id   int    `json:"id"`   //! 配置Id
}

type JsDialMsg struct {
	Id   int    `json:"id"`
	Info string `json:"info"`
}

type DialRank struct {
	Uid   int64  `json:"uid"`
	Times int    `json:"times"`
	Uname string `json:"uname"`
	Camp  int    `json:"camp"`
}

type Js_DialRank struct {
	Uid   int64  `json:"uid"`
	Times int    `json:"times"`
	Uname string `json:"uname"`
	Camp  int    `json:"camp"`
	Rank  int    `json:"rank"`
}

type lstDialRank []*Js_DialRank

func (s lstDialRank) Len() int      { return len(s) }
func (s lstDialRank) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstDialRank) Less(i, j int) bool {
	if s[i].Times > s[j].Times { // 由大到小
		return true
	}

	if s[i].Times < s[j].Times {
		return false
	}

	if s[i].Rank < s[j].Rank { // 由大到小
		return true
	}

	if s[i].Rank > s[j].Rank {
		return false
	}

	if s[i].Uid > s[j].Uid { // 由大到小
		return true
	}

	if s[i].Uid < s[j].Uid {
		return false
	}
	return false
}

type DialMgr struct {
	Id   int
	Info string

	Mu        *sync.RWMutex
	info      []*DialInfo
	dialRank  []*Js_DialRank
	dialMap   map[int64]*Js_DialRank
	dialMu    *sync.RWMutex
	dialMapMu *sync.RWMutex
}

var dialMgr *DialMgr = nil

func GetDialMgr() *DialMgr {
	if dialMgr == nil {
		dialMgr = new(DialMgr)
		dialMgr.Id = 2
		dialMgr.info = make([]*DialInfo, 0)
		dialMgr.Mu = new(sync.RWMutex)
		dialMgr.dialRank = make([]*Js_DialRank, 0)
		dialMgr.dialMu = new(sync.RWMutex)
		dialMgr.dialMap = make(map[int64]*Js_DialRank, 0)
		dialMgr.dialMapMu = new(sync.RWMutex)
	}

	return dialMgr
}

func (self *DialMgr) GetMsg() []*DialInfo {
	self.Mu.RLock()
	defer self.Mu.RUnlock()

	return self.info
}

func (self *DialMgr) AddMsg(name string, id int) *DialInfo {
	self.Mu.Lock()
	defer self.Mu.Unlock()

	if len(self.info) >= 50 {
		self.info = append(self.info[:0], self.info[1:]...)
	}

	pInfo := &DialInfo{Name: name, Id: id}
	self.info = append(self.info, pInfo)
	return pInfo
}

func (self *DialMgr) Encode() {
	self.Info = HF_JtoA(self.info)
}

// 存储数据库
func (self *DialMgr) Save() {
	self.Mu.RLock()
	defer self.Mu.RUnlock()

	// 先查找, 再更新
	queryStr := "select `id`, `info` from `san_investmsg`;"
	res := GetServer().DBUser.Query(queryStr)
	if res {
		// update
		self.Encode()
		updateStr := fmt.Sprintf("update `san_investmsg` set `id` = %d, `info` = '%s';", self.Id, self.Info)
		GetServer().DBUser.Exec(updateStr)
	} else {
		self.Encode()
		insertStr := fmt.Sprintf("insert into `san_investmsg` (`id`, `info`) values (%d, '%s');", self.Id, self.Info)
		GetServer().DBUser.Exec(insertStr)
	}
}

func (self *DialMgr) GetData() {
	queryStr := fmt.Sprintf("select `id`, `info` from `san_investmsg` where `id` = %d", self.Id)
	var msg JsDialMsg
	res := GetServer().DBUser.GetAllData(queryStr, &msg)
	if len(res) > 0 {
		self.info = make([]*DialInfo, 0)
		err := json.Unmarshal([]byte(msg.Info), &self.info)
		if err != nil {
			LogError(err.Error())
		}
	}

	self.GetDialRank()
}

// 排行榜:select `uid`,`times`,`uname`, `camp` from san_userdial where `step` = 1 order by `times` desc, `starttime` asc, `uid` asc limit 50;
func (self *DialMgr) UserDialSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userdial (
		      uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
			  times int(11) NOT NULL DEFAULT '0' COMMENT '次数',
			  starttime bigint(20) NOT NULL DEFAULT '0' COMMENT '抽奖时间',
			  camp int(11) NOT NULL DEFAULT '1' COMMENT '阵营',
			  step int(11) NOT NULL DEFAULT '1' COMMENT '活动期数',
			  info text NOT NULL COMMENT '宝箱状态信息',
			  luck int(11) NOT NULL DEFAULT '0' COMMENT '幸运值',
			  freetime int(11) NOT NULL DEFAULT '0' COMMENT '免费次数',
			  uname varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT '用户名称',
			  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *DialMgr) GetJsRank(pRank *DialRank, rank int) *Js_DialRank {
	return &Js_DialRank{
		Uid:   pRank.Uid,
		Uname: pRank.Uname,
		Times: pRank.Times,
		Camp:  pRank.Camp,
		Rank:  rank,
	}
}

func (self *DialMgr) GetDialRank() {
	actN4 := GetActivityMgr().getActN4(ACT_DIAL)
	queryStr := fmt.Sprintf("select `uid`,`times`,`uname`, `camp` from san_userdial where `step` = %d and `times` > 0 order by `times` desc, `starttime` asc, `uid` asc limit 50", actN4)
	var msg DialRank
	res := GetServer().DBUser.GetAllData(queryStr, &msg)
	if len(res) <= 0 {
		return
	}

	self.dialMu.Lock()
	self.dialRank = make([]*Js_DialRank, 0)
	self.dialMu.Unlock()

	self.dialMu.Lock()
	self.dialMap = make(map[int64]*Js_DialRank)
	self.dialMu.Unlock()

	for index := range res {
		data := res[index].(*DialRank)
		pRank := self.GetJsRank(data, index+1)

		self.dialMu.Lock()
		self.dialRank = append(self.dialRank, pRank)
		self.dialMu.Unlock()

		self.dialMapMu.Lock()
		self.dialMap[data.Uid] = pRank
		self.dialMapMu.Unlock()
	}
}

func (self *DialMgr) GetRankData() []*Js_DialRank {
	self.dialMu.RLock()
	defer self.dialMu.RUnlock()
	return self.dialRank
}

func (self *DialMgr) updateRank(pRank *Js_DialRank) {
	// 检查排行榜有没有玩家
	self.dialMapMu.RLock()
	find, ok := self.dialMap[pRank.Uid]
	self.dialMapMu.RUnlock()
	if !ok {
		// 不在排行榜
		// 和最后一名比较
		self.dialMu.RLock()
		size := len(self.dialRank)
		self.dialMu.RUnlock()

		if size <= 0 {
			self.dialMu.Lock()
			pRank.Rank = 1
			self.dialRank = append(self.dialRank, pRank)
			self.dialMu.Unlock()

			// 设置排行榜
			self.dialMapMu.Lock()
			self.dialMap[pRank.Uid] = pRank
			self.dialMapMu.Unlock()

		} else if size < maxRank {
			self.dialMu.Lock()
			self.dialRank = append(self.dialRank, pRank)
			sort.Sort(lstDialRank(self.dialRank))
			for i := 0; i < len(self.dialRank); i++ {
				if self.dialRank[i] != nil {
					self.dialRank[i].Rank = i + 1
				}
			}
			self.dialMu.Unlock()

			// 设置排行榜
			self.dialMapMu.Lock()
			self.dialMap[pRank.Uid] = pRank
			self.dialMapMu.Unlock()

		} else if size >= maxRank {
			// 先和最后一名比较
			self.dialMu.RLock()
			lastRank := self.dialRank[size-1]
			self.dialMu.RUnlock()
			if lastRank.Times >= pRank.Times {
				// pass
			} else {
				self.dialMu.Lock()
				pRank.Rank = maxRank
				self.dialRank[size-1] = pRank
				self.dialMu.Unlock()

				// 设置排行榜
				self.dialMapMu.Lock()
				delete(self.dialMap, lastRank.Uid)
				self.dialMap[pRank.Uid] = pRank
				self.dialMapMu.Unlock()
			}
		}

	} else {
		// 在排行榜中,重新排序
		find.Times = pRank.Times
		find.Uname = pRank.Uname

		self.dialMu.Lock()
		sort.Sort(lstDialRank(self.dialRank))
		// 重新计算排行榜
		for i := 0; i < len(self.dialRank); i++ {
			if self.dialRank[i] != nil {
				self.dialRank[i].Rank = i + 1
			}
		}
		self.dialMu.Unlock()
	}

}
