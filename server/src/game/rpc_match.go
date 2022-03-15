/*
@Time : 2020/5/10 11:10
@Author : 96121
@File : proto_player
@Software: GoLand
*/
package game

import (
	"net/rpc"
)

const (
	RPC_MATCH_GENERAL_GET_RANK_All = "RPC_Match.GetGeneralAllRankInfo" //! 获得初始排行
	RPC_MATCH_GENERAL_UPDATE       = "RPC_Match.UpdateGeneralInfo"     //! 更新数据

	RPC_MATCH_CROSSARENA_GET_RANK_All = "RPC_Match.GetCrossArenaAllRankInfo"  //! 获得初始排行
	RPC_MATCH_CROSSARENA_ADD          = "RPC_Match.CrossArenaAdd"             //! 参加竞技
	RPC_MATCH_CROSSARENA_GET_DEFENCE  = "RPC_Match.CrossArenaGetDefence"      //! 获得对手
	RPC_MATCH_CROSSARENA_GET_INFO     = "RPC_Match.CrossArenaGetInfo"         //! 查询信息
	RPC_MATCH_CROSSARENA_FIGHT_END    = "RPC_Match.CrossArenaFightEnd"        //! 查询信息
	RPC_MATCH_CROSSARENA_BATTLEINFO   = "RPC_Match.CrossArenaGetBattleInfo"   //!
	RPC_MATCH_CROSSARENA_BATTLERECORD = "RPC_Match.CrossArenaGetBattleRecord" //!

	RPC_MATCH_CONSUMERTOP_GET_RANK_All = "RPC_Match.GetConsumerTopAllRankInfo" //!
	RPC_MATCH_CONSUMERTOP_UPDATE       = "RPC_Match.UpdateConsumerTopInfo"     //!

	RPC_MATCH_CROSSARENA_GET_RANK_All_3V3 = "RPC_Match.GetCrossArena3V3AllRankInfo"  //! 获得初始排行
	RPC_MATCH_CROSSARENA_ADD_3V3          = "RPC_Match.CrossArena3V3Add"             //! 参加竞技
	RPC_MATCH_CROSSARENA_GET_DEFENCE_3V3  = "RPC_Match.CrossArena3V3GetDefence"      //! 获得对手
	RPC_MATCH_CROSSARENA_GET_INFO_3V3     = "RPC_Match.CrossArena3V3GetInfo"         //! 查询信息
	RPC_MATCH_CROSSARENA_FIGHT_END_3V3    = "RPC_Match.CrossArena3V3FightEnd"        //! 查询信息
	RPC_MATCH_CROSSARENA_BATTLEINFO_3V3   = "RPC_Match.CrossArena3V3GetBattleInfo"   //!
	RPC_MATCH_CROSSARENA_BATTLERECORD_3V3 = "RPC_Match.CrossArena3V3GetBattleRecord" //!

	CROSSARENA3V3_TEAM_MAX = 3
)

//! 角色消息主体
type RPC_Match struct {
	Client *rpc.Client
}

//! 事件请求
type RPC_GeneralActionReq struct {
	Uid           int64            //! 角色Id
	SelfInfo      *Js_GeneralUser  //!
	GeneralRecord []*GeneralRecord //! 公告
}

//! 事件请求
type RPC_GeneralActionReqAll struct {
	KeyId    int //! 活动Key
	ServerId int //! 服务器Id
}

//! 事件返回
type RPC_GeneralActionRes struct {
	RetCode       int               //! 返回结果
	RankInfo      []*Js_GeneralUser //! 排行信息
	SelfInfo      *Js_GeneralUser   //! 玩家信息
	GeneralRecord []*GeneralRecord  //! 公告
}

//! 事件请求
type RPC_CrossArenaActionReqAll struct {
	KeyId int //! 各种KEY  int
}

type RPC_CrossArenaAction64ReqAll struct {
	KeyId int64 //! 各种KEY  int64
}

type RPC_CrossArenaActionReq struct {
	Uid       int64              //! 角色Id
	KeyId     int                //!
	SelfInfo  *Js_CrossArenaUser //!
	FightInfo *JS_FightInfo      //!
}

type RPC_CrossArenaFightEndReq struct {
	Uid        int64         //! 角色Id
	KeyId      int           //!
	Attack     *JS_FightInfo //!
	Defend     *JS_FightInfo //!
	BattleInfo BattleInfo
}

//! 事件返回
type RPC_CrossArenaActionRes struct {
	RetCode    int                                  //! 返回结果
	RankInfo   map[int]map[int][]*Js_CrossArenaUser //! 排行信息
	SelfInfo   *Js_CrossArenaUser                   //! 玩家信息
	Result     int                                  //!
	NewFightId int64                                //! 中心服对战报ID的新修正值
}

//! 事件返回
type RPC_CrossArenaGetDefenceRes struct {
	RetCode   int                  //! 返回结果
	Info      []*Js_CrossArenaUser //! 排行信息
	FightInfo []*JS_FightInfo      //! 玩家信息
}

//! 事件返回
type RPC_CrossArenaGetInfoRes struct {
	RetCode      int                //! 返回结果
	Info         *Js_CrossArenaUser //! 信息
	FightInfo    *JS_FightInfo      //! 玩家信息
	LifeTreeInfo *JS_LifeTreeInfo   //!
}

//! 事件返回
type RPC_CrossArenaBattleInfoRes struct {
	RetCode    int //! 返回结果
	BattleInfo *BattleInfo
}

type RPC_CrossArenaBattleRecordRes struct {
	RetCode      int //! 返回结果
	BattleRecord *BattleRecord
}

//3V3
//! 事件请求
type RPC_CrossArena3V3ActionReqAll struct {
	KeyId int //! 各种KEY  int
}

type RPC_CrossArena3V3Action64ReqAll struct {
	KeyId int64 //! 各种KEY  int64
}

//! 事件返回
type RPC_CrossArena3V3ActionReq struct {
	Uid       int64                                 //! 角色Id
	KeyId     int                                   //!
	SelfInfo  *Js_CrossArena3V3User                 //!
	FightInfo [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo //!
}

type RPC_CrossArena3V3FightEndReq struct {
	Uid        int64                                 //! 角色Id
	KeyId      int                                   //!
	Attack     [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo //!
	Defend     [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo //!
	BattleInfo [CROSSARENA3V3_TEAM_MAX]BattleInfo
}

type RPC_CrossArena3V3ActionRes struct {
	RetCode    int                                     //! 返回结果
	RankInfo   map[int]map[int][]*Js_CrossArena3V3User //! 排行信息
	SelfInfo   *Js_CrossArena3V3User                   //! 玩家信息
	Result     int                                     //!
	NewFightId [CROSSARENA3V3_TEAM_MAX]int64           //! 中心服对战报ID的新修正值
}

type RPC_CrossArena3V3GetDefenceRes struct {
	RetCode   int                                     //! 返回结果
	Info      []*Js_CrossArena3V3User                 //! 排行信息
	FightInfo [][CROSSARENA3V3_TEAM_MAX]*JS_FightInfo //! 玩家信息
}

//! 事件返回
type RPC_CrossArena3V3GetInfoRes struct {
	RetCode      int                                   //! 返回结果
	Info         *Js_CrossArena3V3User                 //! 信息
	FightInfo    [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo //! 玩家信息
	LifeTreeInfo *JS_LifeTreeInfo                      //!
}

//! 事件返回
type RPC_CrossArena3V3BattleInfoRes struct {
	RetCode    int //! 返回结果
	BattleInfo []*BattleInfo
}

type RPC_CrossArena3V3BattleRecordRes struct {
	RetCode      int //! 返回结果
	BattleRecord *BattleRecord
}

//ConsumerTop
type RPC_ConsumerTopActionReqAll struct {
	KeyId    int //! 活动Key
	ServerId int //! 服务器Id
}

type RPC_ConsumerTopActionRes struct {
	RetCode  int                     //! 返回结果
	SelfInfo *JS_ConsumerTopUser     //! 玩家信息
	TopUser  []*JS_ConsumerTopUser   //! 跨服个人排行数据
	TopSvr   []*JS_ConsumerTopServer //! 跨服服务器排行数据
}

type RPC_ConsumerTopActionReq struct {
	Uid      int64               //! 角色Id
	SelfInfo *JS_ConsumerTopUser //!
}

func (self *RPC_Match) Init() bool {
	return true
}

func (self *RPC_Match) MatchGeneralUpdate(top *Js_GeneralUser, records []*GeneralRecord) *RPC_GeneralActionRes {
	if self.Client != nil {
		var req RPC_GeneralActionReq
		req.Uid = top.Uid
		req.GeneralRecord = records
		req.SelfInfo = top

		var res RPC_GeneralActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_GENERAL_UPDATE, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchGeneralGetAllRank(keyId int, serverId int) *RPC_GeneralActionRes {
	if self.Client != nil {
		var req RPC_GeneralActionReqAll
		req.KeyId = keyId
		req.ServerId = serverId
		var res RPC_GeneralActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_GENERAL_GET_RANK_All, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArenaGetAllRank(keyId int) *RPC_CrossArenaActionRes {
	if self.Client != nil {
		var req RPC_CrossArenaActionReqAll
		req.KeyId = keyId
		var res RPC_CrossArenaActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_GET_RANK_All, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArenaAdd(keyId int, data *Js_CrossArenaUser, fightInfo *JS_FightInfo) *RPC_CrossArenaActionRes {
	if self.Client != nil {
		var req RPC_CrossArenaActionReq
		req.Uid = data.Uid
		req.KeyId = keyId
		req.SelfInfo = data
		req.FightInfo = fightInfo
		var res RPC_CrossArenaActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_ADD, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArenaGetDefence(keyId int, player *Player) *RPC_CrossArenaGetDefenceRes {
	if self.Client != nil {
		var req RPC_CrossArenaActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.KeyId = keyId
		var res RPC_CrossArenaGetDefenceRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_GET_DEFENCE, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArenaGetInfo(keyId int, uid int64) *RPC_CrossArenaGetInfoRes {
	if self.Client != nil {
		var req RPC_CrossArenaActionReq
		req.Uid = uid
		req.KeyId = keyId
		var res RPC_CrossArenaGetInfoRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_GET_INFO, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArenaFightEnd(keyId int, attack *JS_FightInfo, defend *JS_FightInfo, battleInfo BattleInfo) *RPC_CrossArenaActionRes {
	if self.Client != nil {
		var req RPC_CrossArenaFightEndReq
		req.KeyId = keyId
		req.Attack = attack
		req.Defend = defend
		req.BattleInfo = battleInfo
		var res RPC_CrossArenaActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_FIGHT_END, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArenaGetBattleInfo(keyId int64) *RPC_CrossArenaBattleInfoRes {
	if self.Client != nil {
		var req RPC_CrossArenaAction64ReqAll
		req.KeyId = keyId
		var res RPC_CrossArenaBattleInfoRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_BATTLEINFO, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}
func (self *RPC_Match) MatchCrossArenaGetBattleRecord(keyId int64) *RPC_CrossArenaBattleRecordRes {
	if self.Client != nil {
		var req RPC_CrossArenaAction64ReqAll
		req.KeyId = keyId
		var res RPC_CrossArenaBattleRecordRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_BATTLERECORD, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchConsumerTopGetAllRank(keyId int, serverId int) *RPC_ConsumerTopActionRes {
	if self.Client != nil {
		var req RPC_ConsumerTopActionReqAll
		req.KeyId = keyId
		req.ServerId = serverId
		var res RPC_ConsumerTopActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CONSUMERTOP_GET_RANK_All, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchConsumerTopUpdate(top *JS_ConsumerTopUser) *RPC_ConsumerTopActionRes {
	if self.Client != nil {
		var req RPC_ConsumerTopActionReq
		req.Uid = top.Uid
		req.SelfInfo = top

		var res RPC_ConsumerTopActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CONSUMERTOP_UPDATE, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArena3V3GetAllRank(keyId int) *RPC_CrossArena3V3ActionRes {
	if self.Client != nil {
		var req RPC_CrossArena3V3ActionReqAll
		req.KeyId = keyId
		var res RPC_CrossArena3V3ActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_GET_RANK_All_3V3, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArena3V3Add(keyId int, data *Js_CrossArena3V3User, fightInfo [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo) *RPC_CrossArena3V3ActionRes {
	if self.Client != nil {
		var req RPC_CrossArena3V3ActionReq
		req.Uid = data.Uid
		req.KeyId = keyId
		req.SelfInfo = data
		for k, v := range fightInfo {
			req.FightInfo[k] = v
		}
		var res RPC_CrossArena3V3ActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_ADD_3V3, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArena3V3GetDefence(keyId int, player *Player) *RPC_CrossArena3V3GetDefenceRes {
	if self.Client != nil {
		var req RPC_CrossArena3V3ActionReq
		req.Uid = player.Sql_UserBase.Uid
		req.KeyId = keyId

		data := new(JS_FightInfo)
		for i := 0; i < CROSSARENA3V3_TEAM_MAX; i++ {
			if req.FightInfo[i] == nil {
				req.FightInfo[i] = data
			}
		}
		var res RPC_CrossArena3V3GetDefenceRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_GET_DEFENCE_3V3, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArena3V3GetInfo(keyId int, uid int64) *RPC_CrossArena3V3GetInfoRes {
	if self.Client != nil {
		var req RPC_CrossArena3V3ActionReq
		req.Uid = uid
		req.KeyId = keyId
		temp:=new(JS_FightInfo)
		req.FightInfo[0]=temp
		req.FightInfo[1]=temp
		req.FightInfo[2]=temp
		var res RPC_CrossArena3V3GetInfoRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_GET_INFO_3V3, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArena3V3FightEnd(keyId int, attack [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo, defend [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo, battleInfo [CROSSARENA3V3_TEAM_MAX]BattleInfo) *RPC_CrossArena3V3ActionRes {
	if self.Client != nil {
		var req RPC_CrossArena3V3FightEndReq
		req.KeyId = keyId
		req.Attack = attack
		req.Defend = defend
		req.BattleInfo = battleInfo
		var res RPC_CrossArena3V3ActionRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_FIGHT_END_3V3, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}

func (self *RPC_Match) MatchCrossArena3V3GetBattleInfo(keyId int64) *RPC_CrossArena3V3BattleInfoRes {
	if self.Client != nil {
		var req RPC_CrossArenaAction64ReqAll
		req.KeyId = keyId
		var res RPC_CrossArena3V3BattleInfoRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_BATTLEINFO_3V3, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}
func (self *RPC_Match) MatchCrossArena3V3GetBattleRecord(keyId int64) *RPC_CrossArena3V3BattleRecordRes {
	if self.Client != nil {
		var req RPC_CrossArena3V3Action64ReqAll
		req.KeyId = keyId
		var res RPC_CrossArena3V3BattleRecordRes
		err := GetMasterMgr().CallEx(self.Client, RPC_MATCH_CROSSARENA_BATTLERECORD_3V3, req, &res)

		if err == nil {
			//! 添加成功
			return &res
		} else {
			print(err.Error())
			//! 添加失败
			return nil
		}
	}
	return nil
}
