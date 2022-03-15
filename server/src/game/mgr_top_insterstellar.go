//显示星座的排位值，没有排行榜，可以实时计算。如果策划以后说要做排行帮
package game

import (
	"sync"
)

const MAX_SIZE = 300

type JS_InterstellarTop struct {
	Uid          int64 `json:"uid"`
	Stellarcount int   `json:"stellarcount"`
}

type TopInterstellarMgr struct {
	// 排行数据
	TopInterstellar    []int // 排行榜
	TopInterstellarCur map[int64]int

	Locker *sync.RWMutex
}

var topInsterstellarMgr *TopInterstellarMgr = nil

// 初始排行
func (self *TopInterstellarMgr) GetData() {
	var top4 JS_InterstellarTop
	sql6 := "SELECT uid, stellarcount from  san_interstellar"
	res6 := GetServer().DBUser.GetAllDataEx(sql6, &top4)
	for i := 0; i < len(res6); i++ {
		data := res6[i].(*JS_InterstellarTop)
		self.TopInterstellar[data.Stellarcount]++
		self.TopInterstellarCur[data.Uid] = data.Stellarcount
	}
}

func GetTopInterstellarMgr() *TopInterstellarMgr {
	if topInsterstellarMgr == nil {
		topInsterstellarMgr = new(TopInterstellarMgr)
		topInsterstellarMgr.TopInterstellar = make([]int, MAX_SIZE)
		topInsterstellarMgr.TopInterstellarCur = make(map[int64]int, 0)
		topInsterstellarMgr.Locker = new(sync.RWMutex)
	}
	return topInsterstellarMgr
}

func (self *TopInterstellarMgr) GetCurPos(uid int64) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	_, ok := self.TopInterstellarCur[uid]
	if !ok {
		return 0
	}
	pos := 0
	all := 0
	for i := 0; i < len(self.TopInterstellar); i++ {
		if self.TopInterstellar[i] == 0 {
			continue
		}
		if i <= self.TopInterstellarCur[uid] {
			pos += self.TopInterstellar[i]
		}
		all += self.TopInterstellar[i]
	}
	all += 1
	return int((float64(pos) / float64(all)) * float64(10000))
}

func (self *TopInterstellarMgr) SetCurProgress(player *Player, value int) {
	self.Locker.Lock()
	_, ok := self.TopInterstellarCur[player.Sql_UserBase.Uid]
	if !ok {
		self.TopInterstellar[value]++
	} else {
		oldValue := self.TopInterstellarCur[player.Sql_UserBase.Uid]
		self.TopInterstellar[oldValue]--
		self.TopInterstellar[value]++
	}
	self.TopInterstellarCur[player.Sql_UserBase.Uid] = value
	self.Locker.Unlock()

	//加个消息通知客户端排位
	var msgRel S2C_UpdateInterStellarPos
	msgRel.Cid = "updateinterstellarpos"
	msgRel.Pos = self.GetCurPos(player.Sql_UserBase.Uid)
	msgRel.StellarCount = value
	player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
}
