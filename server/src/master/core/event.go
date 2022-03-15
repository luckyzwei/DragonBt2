/*
@Time : 2020/5/10 23:02
@Author : 96121
@File : event
@Software: GoLand
*/
package core

const (
	PLAYER_EVENT_FIND_FRIEND             = 1  //! 寻找好友
	PLAYER_EVENT_ADD_FRIEND              = 2  //! 增加好友
	PLAYER_EVENT_DEL_FRIEND              = 3  //! 删除好友
	PLAYER_EVENT_AGREEE_FRIEND           = 4  //! 同意添加好友
	PLAYER_EVENT_REFUSE_FRIEND           = 5  //! 拒绝添加好友
	PLAYER_EVENT_BLACK_FRIEND            = 6  //! 拉黑好友
	PLAYER_EVENT_POWER_FRIEND            = 7  //! 友情点赠送
	PLAYER_EVENT_UPDATE_HIRE_HERO        = 11 //! 更新租借英雄信息
	PLAYER_EVENT_UPDATE_HIRE_HERO_SINGLE = 12 //! 更新租借英雄信息 单
	PLAYER_EVENT_AGREE_HIRE_HERO         = 13 //! 通知雇佣成功
	PLAYER_EVENT_REFUSE_HIRE_HERO        = 14 //! 通知雇佣拒绝
	PLAYER_EVENT_DELETE_HIRE             = 15 //! 移除雇佣

	UNION_EVENT_MASTER_FAIL          = 100 // 会长拒绝
	UNION_EVENT_MASTER_OK            = 101 // 会长同意
	UNION_EVENT_OUT_PLAYER           = 102 // 会长踢人
	UNION_EVENT_UNION_MODIFY         = 103 // 通知职位变更
	UNION_EVENT_UNION_SET_BRAVE_HAND = 104 // 设置无畏之手
	UNION_EVENT_UNION_HUNTER_AWARD   = 105 // 狩猎奖励
	UNION_EVENT_UNION_UPDATE         = 106 // 更新公会
	UNION_EVENT_UNION_SEND_MAIL      = 107 // 更新公会

	SYSTEM_EVENT_CHAT_UNION_OPEN = 200 //! 公会聊天频道开启
	SYSTEM_EVENT_CHAT_WORLD_OPEN = 201 //! 世界聊天频道开启

	CHAT_NEW_WORLD_MESSAGE   = 210 //! 聊天新消息世界
	CHAT_NEW_UNION_MESSAGE   = 211 //! 聊天新消息公会
	CHAT_NEW_PRIVATE_MESSAGE = 212 //! 聊天新消息私聊
	CHAT_GAP_PLAYER          = 213 //! 封禁聊天，通知客户端删除对应玩家的当前聊天记录
	//比赛消息
	MATCH_CROSSARENA_UPDATE = 300 //! 跨服竞技场的被动推送
)

const (
	BATTLE_TYPE_TOWER              = 1 // 爬塔
	BATTLE_TYPE_ARENA              = 2 // 普通竞技场
	BATTLE_TYPE_UNION_HUNT_NORMAl  = 3 // 公会狩猎1
	BATTLE_TYPE_UNION_HUNT_SPECIAL = 4 // 公会狩猎2
	BATTLE_TYPE_ARENA_SPECIAL      = 5 // 高阶竞技场
	BATTLE_TYPE_NORMAL             = 6 // 日常副本
	BATTLE_TYPE_RECORD_BOSS        = 7 // 暗域入侵
	BATTLE_TYPE_RECORD_CROSSARENA  = 8 // 跨服竞技战报
)
