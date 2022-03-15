//佣兵系统，管理所有玩家的佣兵状态  20200110 by zy
package game

import (
	"encoding/json"
	"fmt"
	"sync"
)

type HirePlayerBase struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"uname"`
	Face     int    `json:"face"`
	IconId   int    `json:"iconid"`
	Portrait int    `json:"portrait"`
}

type HireHero struct {
	HeroKeyId           int               `json:"herokeyid"`
	HeroId              int               `json:"heroid"`
	HeroQuality         int               `json:"heroquality"`         //英雄品质  (不变的属性)
	HeroArtifactId      int               `json:"heroartifactid"`      //英雄神器等级 (不变的属性)
	HeroArtifactLv      int               `json:"heroartifactlv"`      //英雄神器等级 (不变的属性)
	HeroExclusiveLv     int               `json:"heroexclusivelv"`     //英雄专属等级 (不变的属性)
	HeroExclusiveUnLock int               `json:"heroexclusiveunlock"` //英雄专属解锁 (不变的属性)
	Talent              *StageTalent      `json:"talent"`              //天赋
	OwnPlayer           *HirePlayerBase   `json:"ownplayer"`           //拥有者
	ApplyPlayer         []*HirePlayerBase `json:"applyplayer"`         //申请列表
	HirePlayer          *HirePlayerBase   `json:"hireplayer"`          //
	ReSetTime           int64             `json:"resettime"`           //
}

//数据库结构
type HireHeroInfo struct {
	Uid          int64
	HireHeroInfo string

	hireHeroInfo map[int]*HireHero //自己满足条件的英雄
	DataUpdate
}

// 将数据库数据写入dataf
func (self *HireHeroInfo) Decode() {
	err := json.Unmarshal([]byte(self.HireHeroInfo), &self.hireHeroInfo)
	if err != nil {
		LogError("HireHeroInfo Decode error:", err.Error())
	}
}

// 将data数据写入数据库
func (self *HireHeroInfo) Encode() {
	self.HireHeroInfo = HF_JtoA(&self.hireHeroInfo)
}

// 竞技场管理器
type HireHeroInfoMgr struct {
	HireHeroInfo map[int64]*HireHeroInfo //
	Lock         *sync.RWMutex           // 数据操作锁
}

var HireHeroInfomgr *HireHeroInfoMgr = nil

func GetHireHeroInfoMgr() *HireHeroInfoMgr {
	if HireHeroInfomgr == nil {
		HireHeroInfomgr = new(HireHeroInfoMgr)
		HireHeroInfomgr.HireHeroInfo = make(map[int64]*HireHeroInfo)
		HireHeroInfomgr.Lock = new(sync.RWMutex)
	}
	return HireHeroInfomgr
}

func (self *HireHeroInfoMgr) GetData() {
	var info HireHeroInfo
	tableName := "san_hireheroinfo"
	sql := fmt.Sprintf("select * from `%s`", tableName)
	res := GetServer().DBUser.GetAllData(sql, &info)
	for i := 0; i < len(res); i++ {
		data := res[i].(*HireHeroInfo)
		data.Init(tableName, data, false)
		data.Decode()
		self.HireHeroInfo[data.Uid] = data
	}
}

func (self *HireHeroInfoMgr) Save() {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	for _, value := range self.HireHeroInfo {
		value.Encode()
		value.Update(true)
	}
}

func (self *HireHeroInfoMgr) GetInfo(uid int64) *HireHeroInfo {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	_, ok := self.HireHeroInfo[uid]
	if !ok {
		info := new(HireHeroInfo)
		info.Uid = uid
		info.hireHeroInfo = make(map[int]*HireHero)
		info.Encode()
		tableName := "san_hireheroinfo"
		InsertTable(tableName, info, 0, false)
		info.Init(tableName, info, false)
		self.HireHeroInfo[info.Uid] = info
	}

	if self.HireHeroInfo[uid].hireHeroInfo == nil {
		self.HireHeroInfo[uid].hireHeroInfo = make(map[int]*HireHero)
	}
	return self.HireHeroInfo[uid]
}

func (self *HireHeroInfoMgr) Rename(player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.HireHeroInfo[player.Sql_UserBase.Uid]
	if ok {
		for _, v := range info.hireHeroInfo {
			if v.OwnPlayer != nil {
				v.OwnPlayer.Name = player.Sql_UserBase.UName
			}
		}
	}
}

func (self *HireHeroInfoMgr) GetHireInfo(uid int64, keyId int) *HireHero {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.HireHeroInfo[uid]
	if !ok {
		return nil
	} else {
		hireHero, okInfo := info.hireHeroInfo[keyId]
		if !okInfo {
			return nil
		} else {
			return hireHero
		}
	}
	return nil
}

func (self *HireHeroInfoMgr) AddHireApply(uid int64, keyId int, player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.HireHeroInfo[uid]
	if !ok {
		return
	} else {
		hireHero, okInfo := info.hireHeroInfo[keyId]
		if !okInfo {
			return
		} else {
			//已经在列表里就不加了
			for _, v := range hireHero.ApplyPlayer {
				if v.Uid == player.Sql_UserBase.Uid {
					return
				}
			}
			base := self.NewPlayerBase(player)
			hireHero.ApplyPlayer = append(hireHero.ApplyPlayer, base)
		}
	}
}

func (self *HireHeroInfoMgr) CancelHireApply(uid int64, keyId int, player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.HireHeroInfo[uid]
	if !ok {
		return
	} else {
		hireHero, okInfo := info.hireHeroInfo[keyId]
		if !okInfo {
			return
		} else {
			newLst := make([]*HirePlayerBase, 0)
			for _, v := range hireHero.ApplyPlayer {
				if v.Uid != player.Sql_UserBase.Uid {
					newLst = append(newLst, v)
				}
			}
			hireHero.ApplyPlayer = newLst
		}
	}
}

func (self *HireHeroInfoMgr) LostHireApply(uid int64, keyId int, player *Player) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.HireHeroInfo[uid]
	if !ok {
		return
	} else {
		hireHero, okInfo := info.hireHeroInfo[keyId]
		if !okInfo {
			return
		} else {
			if hireHero.HirePlayer == nil {
				return
			}
			if hireHero.HirePlayer.Uid != player.Sql_UserBase.Uid {
				return
			}
			hireHero.HirePlayer = nil
		}
	}
}

func (self *HireHeroInfoMgr) DeteleHire(player *Player, keyId int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.HireHeroInfo[player.Sql_UserBase.Uid]
	if !ok {
		return
	} else {
		_, okInfo := info.hireHeroInfo[keyId]
		if !okInfo {
			return
		} else {
			delete(info.hireHeroInfo, keyId)
		}
	}
}

/*
func (self *HireHeroInfoMgr) SetHireApply(uid int64, keyId int, state int, hireUid int64) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.HireHeroInfo[uid]
	if !ok {
		return
	}
	hireHero, okInfo := info.hireHeroInfo[keyId]
	if !okInfo {
		return
	}
	//已经在列表里就不加了
	for _, v := range hireHero.ApplyPlayer {
		if v.Uid == hireUid {
			if state == LOGIC_TRUE {
				hireHero.HirePlayer = new(HirePlayerBase)
				HF_DeepCopy(hireHero.HirePlayer, v)
				hireHero.ApplyPlayer = make([]*HirePlayerBase, 0)
				//在线情况下生成英雄，发送到列表
				player := GetPlayerMgr().GetPlayer(hireUid, false)
				if player != nil {
					friendHero := player.GetModule("friend").(*ModFriend).AddHireHero(hireHero)
					if friendHero != nil {
						var msgRel S2C_UpdateFriendHire
						msgRel.Cid = "updatefriendhire"
						msgRel.HireHero = friendHero
						player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
					}
				}
			} else {
				newList := make([]*HirePlayerBase, 0)
				for _, value := range hireHero.ApplyPlayer {
					if v.Uid != hireUid {
						newList = append(newList, value)
					}
				}
				hireHero.ApplyPlayer = newList
			}
			//看这个玩家是否在线，决定同步
			player := GetPlayerMgr().GetPlayer(hireUid, false)
			if player != nil {
				var msgRel S2C_HireStateUpdate
				msgRel.Cid = "hirestateupdate"
				msgRel.HireHero = hireHero
				player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
			}
			return
		}
	}
}
*/

/*
func (self *HireHeroInfoMgr) SetHireApplyAll(uid int64, keyId int) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	info, ok := self.HireHeroInfo[uid]
	if !ok {
		return
	}
	hireHero, okInfo := info.hireHeroInfo[keyId]
	if !okInfo {
		return
	}
	//已经在列表里就不加了
	needNotice := make([]int64, 0)
	state := LOGIC_TRUE
	for _, v := range hireHero.ApplyPlayer {
		if state == LOGIC_TRUE {
			hireHero.HirePlayer = new(HirePlayerBase)
			HF_DeepCopy(&hireHero.HirePlayer, &v)
			//在线情况下生成英雄，发送到列表
			player := GetPlayerMgr().GetPlayer(v.Uid, false)
			state = LOGIC_FALSE
			if player != nil {
				friendHero := player.GetModule("friend").(*ModFriend).AddHireHero(hireHero)
				if friendHero != nil {
					var msgRel S2C_UpdateFriendHire
					msgRel.Cid = "updatefriendhire"
					msgRel.HireHero = friendHero
					player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
				}
			}
		}
		needNotice = append(needNotice, v.Uid)
	}
	hireHero.ApplyPlayer = make([]*HirePlayerBase, 0)

	for i := 0; i < len(needNotice); i++ {
		//看这个玩家是否在线，决定同步
		player := GetPlayerMgr().GetPlayer(needNotice[i], false)
		if player != nil {
			var msgRel S2C_HireStateUpdate
			msgRel.Cid = "hirestateupdate"
			msgRel.HireHero = hireHero
			player.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
		}
	}
}
*/

func (self *HireHeroInfoMgr) NewPlayerBase(player *Player) *HirePlayerBase {

	relInfo := new(HirePlayerBase)
	relInfo.Uid = player.Sql_UserBase.Uid
	relInfo.IconId = player.Sql_UserBase.IconId
	relInfo.Name = player.Sql_UserBase.UName
	relInfo.Portrait = player.Sql_UserBase.Portrait
	relInfo.Face = player.Sql_UserBase.Face

	return relInfo
}
