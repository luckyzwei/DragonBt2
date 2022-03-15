/*
@Time : 2020/5/10 11:10
@Author : 96121
@File : proto_player
@Software: GoLand
*/
package game

import (
	"master/utils"
	"net/rpc"
	"sync"
	//"time"
)

const (
	RPC_REGPLAYER      = "RPC_Player.RegPlayer"      //! 注册用户数据
	RPC_GETPLAYER      = "RPC_Player.GetPlayer"      //! 获取用户数据
	RPC_GETPLAYERARENA = "RPC_Player.GetPlayerArena" //! 获取JJC防守阵容
)

const (
	PLAYER_UPDATE_INTERVAL = 1800 //! 更新时间间隔 - 30分钟
)

//! 角色消息主体
type RPC_Player struct {
	Client       *rpc.Client
	PlayerLocker *sync.RWMutex                 //! 数据锁
	MapPlayer    map[int64]*RPC_PlayerData_Req //! 远程数据缓存
}

//---------------------------  Player Data  ---------------------
//! 角色信息结构
type JS_PlayerData struct {
	UId            int64         `json:"uid"`            //! 角色Id
	UName          string        `json:"uname"`          //! 昵称
	Level          int           `json:"level"`          //! 等级
	Fight          int64         `json:"fight"`          //! 战力
	PassId         int           `json:"passid"`         //! 最大关卡
	ServerId       int           `json:"serverid"`       //! 服务器ID
	Sex            int           `json:"sex"`            //! 性别
	IconId         int           `json:"iconid"`         //! 头像
	Portrait       int           `json:"portrait"`       //! 头像框
	Vip            int           `json:"vip"`            //! VIP
	RegTime        int64         `json:"reg"`            //! 注册时间
	LoginTime      int64         `json:"login"`          //! 登录时间
	LastUpdate     int64         `json:"update"`         //! 上次更新
	UnionId        int           `json:"unionid"`        //! 军团id
	ArenaFightInfo *JS_FightInfo `json:"arenafightinfo"` //! JJC防守阵容
	UnionName      string        `json:"unionname"`      //! 公会名
	Position       int           `json:"position"`       //! 职位
	BraveHand      int           `json:"bravehand"`      //! 无畏之手
	UserSignature  string        `json:"usersignature"`  //! 签名
	UserID         string        `json:"userid"`         //! sdk分配id
}

//! 英雄结构
type JS_PlayerHero struct {
	HeroId          int          `json:"heroid"`
	Level           int          `json:"level"`
	Star            int          `json:"star"`
	ArtifactId      int          `json:"artifactid"`      //! 神器ID
	ArtifactLv      int          `json:"artifactlv"`      //! 神器等级
	ExclusiveId     int          `json:"ExclusiveId"`     //! 专属装备Id
	ExclusiveLv     int          `json:"ExclusiveLv"`     //! 专属装备等级
	ExclusiveUnLock int          `json:"ExclusiveUnLock"` //! 专属装备解锁状态0未解锁  1解锁
	Skin            int          `json:"skin"`            //！皮肤
	Talent          *StageTalent `json:"talent"`          // 新天赋系统 和老天赋分开

	Attr map[int]int64 `json:"attr"` //！展示属性
}

//! 英雄装备
type JS_HeroEquip struct {
	ItemId int `json:"itemid"`
	Level  int `json:"level"`
}

//! RPC 角色数据
type RPC_PlayerData_Req struct {
	UId      int64             //! 角色ID
	Online   int               //! 在线状态状态
	Data     *JS_PlayerData    //! 基础数据
	Heros    []*JS_PlayerHero  //! 武将列表
	Equips   [][]*JS_HeroEquip //! 装备列表，最大5个，可为空
	LifeTree *JS_LifeTreeInfo  //! 生命树
}

type RPC_PlayerData_Res struct {
	RetCode int //! 操作结果
	Data    string
}

func (self *RPC_Player) Init() bool {
	self.PlayerLocker = new(sync.RWMutex)
	self.MapPlayer = make(map[int64]*RPC_PlayerData_Req)

	return true
}

//! 注册用户
func (self *RPC_Player) RegPlayer(player *Player) bool {
	//GetPlayerMgr().AddPlayer(req)
	//return nil
	if self.Client != nil {
		var req RPC_PlayerData_Req
		req.UId = player.Sql_UserBase.Uid
		regTime, _ := utils.Parse(player.Sql_UserBase.Regtime)
		logTime, _ := utils.Parse(player.Sql_UserBase.LastLoginTime)
		req.Data = &JS_PlayerData{
			UId:      player.Sql_UserBase.Uid,
			Level:    player.Sql_UserBase.Level,
			UName:    player.Sql_UserBase.UName,
			IconId:   player.Sql_UserBase.IconId,
			Portrait: player.Sql_UserBase.Portrait,
			Sex:      player.Sql_UserBase.Face,
			//ServerId:      player.Account.ServerId,
			ServerId:      GetServer().Con.ServerId,
			PassId:        player.Sql_UserBase.PassMax,
			RegTime:       regTime.Unix(),
			LoginTime:     logTime.Unix(),
			Fight:         player.GetModule("crystal").(*ModResonanceCrystal).San_ResonanceCrystal.MaxFightAll,
			LastUpdate:    TimeServer().Unix(),
			UserSignature: player.Sql_UserBase.UserSignature,
			UserID:        player.Account.UserId,
		}
		req.Data.UnionId = player.GetUnionId()
		req.Heros, req.Equips = player.GetModule("hero").(*ModHero).GetShowHeros()
		req.Data.ArenaFightInfo = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_1)

		//! 在线状态和时间更新
		req.Online = 1
		modLifeTree := player.GetModule("lifetree").(*ModLifeTree)
		if modLifeTree != nil {
			req.LifeTree = new(JS_LifeTreeInfo)
			req.LifeTree.MainLevel = modLifeTree.San_LifeTree.MainLevel
			req.LifeTree.Info = modLifeTree.San_LifeTree.info
		}

		var res RPC_PlayerData_Res
		err := GetMasterMgr().CallEx(self.Client,RPC_REGPLAYER, req, &res)
		if err == nil {
			player.NoticeBaseInfo = false
			//!调用成功
		} else {
			LogDebug("同步角色信息到中心服错误：", err.Error())
		}
	}

	return true
}

//! 载入新玩家
func (self *RPC_Player) GetPlayer(uid int64) *RPC_PlayerData_Req {
	/*  优化以后在做  zy
	self.PlayerLocker.RLock()
	playerData, ok := self.MapPlayer[uid]
	self.PlayerLocker.RUnlock()
	if ok {
		//! 离线状态直接返回
		if playerData.Online == 0 {
			return playerData
		} else if playerData.Data.LastUpdate >= TimeServer().Unix()-PLAYER_UPDATE_INTERVAL {
			//! 一小时内更新过，则直接返回
			return playerData
		}
	}
	*/

	//! 超时，或者未找到，则重新请求
	if self.Client != nil {
		var res RPC_PlayerData_Req
		err := GetMasterMgr().CallEx(self.Client,RPC_GETPLAYER, uid, &res)
		if err != nil {
			LogDebug("获取信息失败", err.Error())
			return nil
		} else {
			//! 获取信息成功
			//GetMasterMgr().AddPlayer(uid, &res)
			self.PlayerLocker.Lock()
			self.MapPlayer[uid] = &res
			self.PlayerLocker.Unlock()
			return &res
		}
	}

	return nil
}

func (self *RPC_Player) SetPlayerOffline(player *Player) {

	if self.Client != nil {
		var req RPC_PlayerData_Req
		req.UId = player.Sql_UserBase.Uid
		regTime, _ := utils.Parse(player.Sql_UserBase.Regtime)
		logTime, _ := utils.Parse(player.Sql_UserBase.LastLoginTime)
		req.Data = &JS_PlayerData{
			UId:      player.Sql_UserBase.Uid,
			Level:    player.Sql_UserBase.Level,
			UName:    player.Sql_UserBase.UName,
			IconId:   player.Sql_UserBase.IconId,
			Portrait: player.Sql_UserBase.Portrait,
			Sex:      player.Sql_UserBase.Face,
			//ServerId:      player.Account.ServerId,
			ServerId:      GetServer().Con.ServerId,
			PassId:        player.Sql_UserBase.PassMax,
			RegTime:       regTime.Unix(),
			LoginTime:     logTime.Unix(),
			Fight:         player.GetModule("crystal").(*ModResonanceCrystal).San_ResonanceCrystal.MaxFightAll,
			LastUpdate:    TimeServer().Unix(),
			UserSignature: player.Sql_UserBase.UserSignature,
			UserID:        player.Account.UserId,
		}
		req.Data.UnionId = player.GetUnionId()
		req.Heros, req.Equips = player.GetModule("hero").(*ModHero).GetShowHeros()
		req.Data.ArenaFightInfo = GetRobotMgr().GetPlayerFightInfoByPos(player, 0, 0, TEAMTYPE_ARENA_1)

		//! 在线状态和时间更新
		req.Online = 0

		modLifeTree := player.GetModule("lifetree").(*ModLifeTree)
		if modLifeTree != nil {
			req.LifeTree = new(JS_LifeTreeInfo)
			req.LifeTree.MainLevel = modLifeTree.San_LifeTree.MainLevel
			req.LifeTree.Info = modLifeTree.San_LifeTree.info
		}

		var res RPC_PlayerData_Res

		err := GetMasterMgr().CallEx(self.Client,RPC_REGPLAYER, req, &res)
		if err == nil {
			player.NoticeBaseInfo = false
			//!调用成功
		} else {
			LogDebug("同步角色信息到中心服错误：", err.Error())
		}
	}
}

//! 载入新玩家
func (self *RPC_Player) GetPlayerArena(uid int64) *RPC_PlayerData_Res {
	if self.Client != nil {
		var res RPC_PlayerData_Res
		err := GetMasterMgr().CallEx(self.Client,RPC_GETPLAYERARENA, uid, &res)
		if err != nil {
			LogDebug("获取信息失败", err.Error())
			return nil
		} else {
			return &res
		}
	}

	return nil
}
