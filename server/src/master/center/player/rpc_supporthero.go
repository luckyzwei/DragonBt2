package player

import (
	"encoding/json"
	"master/utils"
)

//////////////////////////////////////////////////////////////

//! 操作请求
type RPC_SupportHeroAction struct {
	Data string
}

//! 操作响应
type RPC_SupportHeroActionRet struct {
	RetCode int //! 结果码
	Data    string
}

///////////添加支援英雄//////////////
type S2M_SupportHeroAdd struct {
	Uid       int64
	Index     int
	HeroKeyId int
	HeroID    int
	HeroStar  int
	HeroLv    int
	Skin      int
}

type M2S_SupportHeroAdd struct {
}

///////////取消支援英雄//////////////
type S2M_SupportHeroRemove struct {
	Uid       int64
	HeroKeyId int
}

type M2S_SupportHeroRemove struct {
}

///////////使用支援英雄//////////////
type S2M_SupportHeroUse struct {
	Uid       int64
	HeroKeyId int
	Useruid   int64
	Username  string
	Type      int
	Endtime   int64
}

type M2S_SupportHeroUse struct {
}

///////////取消使用支援英雄//////////////
type S2M_SupportHeroCancelUse struct {
	Uid       int64
	HeroKeyId int
	Useruid   int64
}

type M2S_SupportHeroCancelUse struct {
}

///////////获得玩家数据//////////////
type S2M_SupportHeroGetPlayerData struct {
	Uid int64
}

type M2S_SupportHeroGetPlayerData struct {
	Data []*SupportHero
}

///////////清理玩家数据//////////////
type S2M_SupportHeroCleanPlayerData struct {
	Uid int64
}

type M2S_SupportHeroCleanPlayerData struct {
}

///////////获得玩家数据//////////////
type S2M_SupportHeroGetCanUseHero struct {
	Uids   map[int64]int64
	HeroID int
}

type M2S_SupportHeroGetCanUseHero struct {
	Data []*SupportHero
}

///////////使用支援英雄//////////////
type S2M_SupportHeroUpdate struct {
	Uid       int64
	HeroKeyId int
	HeroStar  int
	HeroLv    int
	Skin      int
}

type M2S_SupportHeroUpdate struct {
}

///////////改名//////////////
type S2M_SupportHeroRename struct {
	Uid  int64
	Name string
}

type M2S_SupportHeroRename struct {
}

func (self *RPC_Friend) AddSupportHero(req RPC_SupportHeroAction, ret *RPC_SupportHeroActionRet) error {
	var msg S2M_SupportHeroAdd
	json.Unmarshal([]byte(req.Data), &msg)
	my := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if my != nil {
		if !my.DataFriend.AddSupportHero(msg.Index, &MsgSupportHero{msg.HeroKeyId, msg.HeroID, msg.HeroStar, msg.HeroLv, msg.Skin}) {
			ret.RetCode = UNION_NOT_MASTER
			return nil
		}
	}

	ret.RetCode = UNION_SUCCESS
	return nil
}

func (self *RPC_Friend) RemoveSupportHero(req RPC_SupportHeroAction, ret *RPC_SupportHeroActionRet) error {
	var msg S2M_SupportHeroRemove
	json.Unmarshal([]byte(req.Data), &msg)
	my := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if my != nil {
		if !my.DataFriend.RemoveSupportHero(msg.HeroKeyId) {
			ret.RetCode = UNION_NOT_MASTER
			return nil
		}
	}
	ret.RetCode = UNION_SUCCESS
	return nil
}

func (self *RPC_Friend) UseHero(req RPC_SupportHeroAction, ret *RPC_SupportHeroActionRet) error {
	var msg S2M_SupportHeroUse
	json.Unmarshal([]byte(req.Data), &msg)
	my := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if my != nil {
		if !my.DataFriend.UseHero(msg.HeroKeyId, msg.Useruid, msg.Username, msg.Type, msg.Endtime) {
			ret.RetCode = UNION_NOT_MASTER
			return nil
		}
	}
	ret.RetCode = UNION_SUCCESS
	return nil
}

func (self *RPC_Friend) CancelUseHero(req RPC_SupportHeroAction, ret *RPC_SupportHeroActionRet) error {
	var msg S2M_SupportHeroCancelUse
	json.Unmarshal([]byte(req.Data), &msg)
	my := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if my != nil {
		if !my.DataFriend.CancelUseHero(msg.HeroKeyId, msg.Useruid) {
			ret.RetCode = UNION_NOT_MASTER
			return nil
		}
	}
	ret.RetCode = UNION_SUCCESS
	return nil
}

// 获得玩家数据
func (self *RPC_Friend) GetPlayerData(req RPC_SupportHeroAction, ret *RPC_SupportHeroActionRet) error {
	var msg S2M_SupportHeroGetPlayerData
	json.Unmarshal([]byte(req.Data), &msg)
	my := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if my != nil {
		if my.DataFriend.Data.supportHero == nil {
			my.DataFriend.Data.supportHero = make([]*SupportHero, 0)
		}
		var backmsg M2S_SupportHeroGetPlayerData
		backmsg.Data = my.DataFriend.Data.supportHero
		ret.RetCode = UNION_SUCCESS
		ret.Data = utils.HF_JtoA(backmsg)
		return nil
	}
	ret.RetCode = UNION_NOT_MASTER
	return nil
}

// 获得玩家数据
func (self *RPC_Friend) CleanPlayerData(req RPC_SupportHeroAction, ret *RPC_SupportHeroActionRet) error {
	var msg S2M_SupportHeroCleanPlayerData
	json.Unmarshal([]byte(req.Data), &msg)
	my := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if my != nil {
		my.DataFriend.Data.supportHero = []*SupportHero{}
	}
	ret.RetCode = UNION_SUCCESS
	return nil
}

// 获得可使用的英雄
func (self *RPC_Friend) GetCanUseHero(req RPC_SupportHeroAction, ret *RPC_SupportHeroActionRet) error {
	var msg S2M_SupportHeroGetCanUseHero
	json.Unmarshal([]byte(req.Data), &msg)

	data := []*SupportHero{}
	for _, uid := range msg.Uids {
		my := GetPlayerMgr().GetPlayer(uid, true)
		if my == nil {
			continue
		}

		// 检测过期
		my.DataFriend.CheckEndTime()
		for _, v := range my.DataFriend.Data.supportHero {
			if msg.HeroID != 0 && msg.HeroID != v.HeroID {
				continue
			}
			data = append(data, v)
		}
	}
	var backmsg M2S_SupportHeroGetPlayerData
	backmsg.Data = data
	ret.RetCode = UNION_SUCCESS
	ret.Data = utils.HF_JtoA(backmsg)
	return nil
}
func (self *RPC_Friend) UpdateHero(req RPC_SupportHeroAction, ret *RPC_SupportHeroActionRet) error {
	var msg S2M_SupportHeroUpdate
	json.Unmarshal([]byte(req.Data), &msg)
	my := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if my != nil {
		my.DataFriend.UpdateHero(&MsgSupportHero{msg.HeroKeyId, 0, msg.HeroStar, msg.HeroLv, msg.Skin})
	}
	ret.RetCode = UNION_SUCCESS
	return nil
}

func (self *RPC_Friend) Rename(req RPC_SupportHeroAction, ret *RPC_SupportHeroActionRet) error {
	var msg S2M_SupportHeroRename
	json.Unmarshal([]byte(req.Data), &msg)
	my := GetPlayerMgr().GetPlayer(msg.Uid, true)
	if my != nil {
		my.DataFriend.Rename(msg.Name)
	}
	ret.RetCode = UNION_SUCCESS
	return nil
}
