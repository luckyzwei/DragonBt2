package core

import "game"

//! 程序入口接口
type IMasterApp interface {
	Wait()                            //! 增加一个等待计数
	Done()                            //! 完成一个等待计数
	GetConfig() *Config               //! 系统配置
	IsClosed() bool                   //! 系统是否关闭
	GetPlayerOnline(serverId int) int //! 在线人数

	StartService() //! 开启服务
	StopService()  //! 关闭
}

type ICenterApp interface {
	AddEvnet(int, interface{})                                                      //! 触发事件
	AddEvent(sid int, code int, uid int64, target int64, param1 int, param2 string) //! 增加事件
}

type IGateApp interface {
	AddEvnet(int, interface{}) //! 触发事件
}

//!玩家管理类
type IPlayerMgr interface {
	GetCorePlayer(int64, bool) IPlayer     //! 获取在线角色，不存在返回空
	GetOnline() int                        //! 获取在线人数
	BroadcastMsg(head string, body []byte) //! 发布消息
}

type IPlayer interface {
	GetUid() int64
	GetUname() string
	GetLevel() int
	GetIconId() int
	GetPortrait() int
	GetVip() int
	GetFight() int
	GetServerId() int
	GetUnionId() int
	GetLifeTree() *game.JS_LifeTreeInfo

	GetDataInt(attType int) int
	GetDataInt64(attType int) int64
	GetDataString(attType int) string

	GetSession() ISession
	OnClose()
}

type ISession interface {
}

//! 系统属性
const (
	PROPERTY_OPEN_SERVER = "OpenServerTime"
	PROPERTY_OPEN_DAY    = "OpenDay"
	PROPERTY_WORLD_LEVEL = "WorldLevel"
)

const (
	PLAYER_ATT_LEVEL = 1
	PLAYER_ATT_VIP   = 2
	PLAYER_ATT_UNAME = 3
	PLAYER_ATT_ICON  = 4
)
