package match

import "game"

const (
	RETCODE_OK         = 0 //! 没有错误
	RETCODE_DATA_ERROR = 1 //! 数据异常
)

const (
	RECORD_MAX = 50
)

//! 消息主体
type RPC_Match struct {
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

//! 事件返回
type RPC_CrossArenaActionReq struct {
	Uid       int64              //! 角色Id
	KeyId     int                //!
	SelfInfo  *Js_CrossArenaUser //!
	FightInfo *game.JS_FightInfo //!
}

type RPC_CrossArenaFightEndReq struct {
	Uid        int64              //! 角色Id
	KeyId      int                //!
	Attack     *game.JS_FightInfo //!
	Defend     *game.JS_FightInfo //!
	BattleInfo game.BattleInfo
}

type RPC_CrossArenaActionRes struct {
	RetCode    int                                  //! 返回结果
	RankInfo   map[int]map[int][]*Js_CrossArenaUser //! 排行信息
	SelfInfo   *Js_CrossArenaUser                   //! 玩家信息
	Result     int                                  //!
	NewFightId int64                                //! 中心服对战报ID的新修正值
}

type RPC_CrossArenaGetDefenceRes struct {
	RetCode   int                  //! 返回结果
	Info      []*Js_CrossArenaUser //! 排行信息
	FightInfo []*game.JS_FightInfo //! 玩家信息
}

//! 事件返回
type RPC_CrossArenaGetInfoRes struct {
	RetCode      int                   //! 返回结果
	Info         *Js_CrossArenaUser    //! 信息
	FightInfo    *game.JS_FightInfo    //! 玩家信息
	LifeTreeInfo *game.JS_LifeTreeInfo //!
}

//! 事件返回
type RPC_CrossArenaBattleInfoRes struct {
	RetCode    int //! 返回结果
	BattleInfo *game.BattleInfo
}

type RPC_CrossArenaBattleRecordRes struct {
	RetCode      int //! 返回结果
	BattleRecord *game.BattleRecord
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
	Uid      int64           //! 角色Id
	SelfInfo *JS_ConsumerTopUser //!
}

//获取排行信息
func (self *RPC_Match) GetGeneralAllRankInfo(req *RPC_GeneralActionReqAll, res *RPC_GeneralActionRes) error {
	GetGeneralMgr().GetAllRank(req, res)
	return nil
}

//上传玩家信息
func (self *RPC_Match) UpdateGeneralInfo(req *RPC_GeneralActionReq, res *RPC_GeneralActionRes) error {
	GetGeneralMgr().UpdatePoint(req, res)
	return nil
}

//获取排行信息
func (self *RPC_Match) GetCrossArenaAllRankInfo(req *RPC_CrossArenaActionReqAll, res *RPC_CrossArenaActionRes) error {
	GetCrossArenaMgr().GetAllRank(req, res)
	return nil
}

func (self *RPC_Match) CrossArenaAdd(req *RPC_CrossArenaActionReq, res *RPC_CrossArenaActionRes) error {
	GetCrossArenaMgr().AddInfo(req, res)
	return nil
}

func (self *RPC_Match) CrossArenaGetDefence(req *RPC_CrossArenaActionReq, res *RPC_CrossArenaGetDefenceRes) error {
	GetCrossArenaMgr().GetDefence(req, res)
	return nil
}

func (self *RPC_Match) CrossArenaGetInfo(req *RPC_CrossArenaActionReq, res *RPC_CrossArenaGetInfoRes) error {
	GetCrossArenaMgr().GetPlayerInfo(req, res)
	return nil
}

func (self *RPC_Match) CrossArenaFightEnd(req *RPC_CrossArenaFightEndReq, res *RPC_CrossArenaActionRes) error {
	GetCrossArenaMgr().FightEnd(req, res)
	return nil
}

func (self *RPC_Match) CrossArenaGetBattleInfo(req *RPC_CrossArenaAction64ReqAll, res *RPC_CrossArenaBattleInfoRes) error {
	GetCrossArenaMgr().GetBattleInfo(req, res)
	return nil
}

func (self *RPC_Match) CrossArenaGetBattleRecord(req *RPC_CrossArenaAction64ReqAll, res *RPC_CrossArenaBattleRecordRes) error {
	GetCrossArenaMgr().GetBattleRecord(req, res)
	return nil
}

//获取排行信息
func (self *RPC_Match) GetConsumerTopAllRankInfo(req *RPC_ConsumerTopActionReqAll, res *RPC_ConsumerTopActionRes) error {
	GetConsumerTopMgr().GetAllRank(req, res)
	return nil
}

//上传玩家信息
func (self *RPC_Match) UpdateConsumerTopInfo(req *RPC_ConsumerTopActionReq, res *RPC_ConsumerTopActionRes) error {
	GetConsumerTopMgr().UpdatePoint(req, res)
	return nil
}