/*
@Time : 2020/4/22 10:28 
@Author : 96121
@File : mod_server
@Software: GoLand
*/
package server

import (
	"master/db"
)

const (
	Max_Server_Event = 20000 //! 事件列表最长
)

//! 游戏服节点
type GameSvrNode struct {
	ID         int    //! 服务器Id
	Name       string //! 名字
	Online     int    //! 在线人数
	LastUpdate int64  //! 上次更新时间
}

type ServerEvent struct {
	EventCode int    `json:"code"`   //! 事件类型
	UId       int64  `json:"uid"`    //! 触发UID
	Target    int64  `json:"target"` //! 目标UID
	Param1    int    `json:"param1"` //! 参数1
	Param2    string `json:"param2"` //! 参数2
}

type SQL_GameServer struct {
	Id         int    //! 当前Id
	ServerId   int    //! 服务器ID
	Name       string //! 名字
	Online     int    //! 在线
	LastUpdate int64  //! 上次更新

	db.DataUpdate //! 数据库操作接口
}

func (self *SQL_GameServer) Encode() {

}

func (self *SQL_GameServer) Decode() {

}

type GameServer struct {
	Data     SQL_GameServer    //! 游戏服数据
	EventArr chan *ServerEvent //! 事件推送
}

func (self *GameServer) OnSave() {
	self.Data.Update(true, false)
}

func (self *GameServer) PushEvent(code int, uid, target int64, param1 int, param2 string) {
	evt := &ServerEvent{
		EventCode: code,
		UId:       uid,
		Target:    target,
		Param1:    param1,
		Param2:    param2,
	}

	self.EventArr <- evt
}

func (self *GameServer) PopEvent() *ServerEvent {
	if len(self.EventArr) == 0 {
		return nil
	}

	evt, ok := <-self.EventArr
	if ok {
		return evt
	}

	return nil
}
