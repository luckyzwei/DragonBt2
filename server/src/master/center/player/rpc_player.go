/*
@Time : 2020/5/10 11:10
@Author : 96121
@File : proto_player
@Software: GoLand
*/
package player

import (
	"errors"
	"fmt"
	"game"
	"master/center/union"
	"master/utils"
	"time"
)

//! 角色消息主体
type RPC_Player struct {
}

//! RPC 角色数据
type RPC_PlayerData_Req struct {
	UId      int64             //! 角色ID
	Online   int               //! 在线状态
	Data     *JS_PlayerData    //! 基础数据
	Heros    []*JS_PlayerHero  //! 武将列表
	Equips   [][]*JS_HeroEquip //! 装备列表，最大5个，可为空
	LifeTree *game.JS_LifeTreeInfo  //! 生命树
}

type RPC_PlayerData_Res struct {
	RetCode int //! 操作结果
	Data    string
}

//! 注册用户
func (self *RPC_Player) RegPlayer(req *RPC_PlayerData_Req, res *RPC_PlayerData_Res) error {
	utils.LogDebug("Reg Player :", req.UId, req.Data.UName)
	needSave := false
	player := GetPlayerMgr().GetPlayer(req.UId, true)
	if player != nil {
		if player.Data.UName == "" {
			needSave = true
		}
		player.Data.UId = req.Data.UId
		player.Data.Level = req.Data.Level
		player.Data.PassId = req.Data.PassId
		player.Data.ServerId = req.Data.ServerId
		player.Data.UName = req.Data.UName
		player.Data.Fight = int(req.Data.Fight)
		player.Data.LoginTime = req.Data.LoginTime
		player.Data.RegTime = req.Data.RegTime
		player.Data.LastUpdate = time.Now().Unix()
		player.Online = req.Online

		myunion := union.GetUnionMgr().GetUnion(player.Data.data.UnionId)
		if myunion != nil {
			member := myunion.GetMember(player.GetUid())
			if member != nil {
				if req.Online == game.LOGIC_TRUE {
					if member.Lastlogintime != 0 {
						member.Lastlogintime = 0
						member.Stage = req.Data.PassId
					}
				} else {
					if member.Lastlogintime == 0 {
						member.Lastlogintime = time.Now().Unix()
						member.Stage = req.Data.PassId
					}
				}
			}
		}

		player.Data.data = req.Data
		player.Data.heros = req.Heros
		player.Data.equips = req.Equips
		player.Data.lifetree = req.LifeTree

		if needSave == false {
			player.Save()
		} else {
			player.onSave(needSave)
		}

		player.DataFriend.UpdateData()
	}
	return nil
}

//! 载入新玩家
func (self *RPC_Player) GetPlayer(uid int64, res *RPC_PlayerData_Req) error {
	utils.LogDebug("Get Player :", uid)

	player := GetPlayerMgr().GetPlayer(uid, true)
	if player != nil {
		if player.Data.data.LastUpdate == 0 {
			player.Data.data.LastUpdate = time.Now().Unix()
		}
		res.UId = player.GetUid()
		res.Online = player.Online
		res.Data = player.Data.data
		res.Heros = player.Data.heros
		res.Equips = player.Data.equips
		res.LifeTree = player.Data.lifetree
		// 公会数据
		myunion := union.GetUnionMgr().GetUnion(player.Data.data.UnionId)
		if myunion != nil {
			member := myunion.GetMember(player.GetUid())
			if nil != member {
				res.Data.UnionName = myunion.Unionname
				res.Data.Position = member.Position
				res.Data.BraveHand = member.BraveHand
			}
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Player [%d] is nil .", uid))
}

//! 载入新玩家
func (self *RPC_Player) GetPlayerArena(uid int64, res *RPC_PlayerData_Res) error {
	utils.LogDebug("GetPlayerArena :", uid)

	GetPlayerMgr().GetPlayerArena(uid, res)
	if res.RetCode != RETCODE_OK {
		return errors.New(fmt.Sprintf("GetPlayerArena .", uid))
	}

	return nil
}
