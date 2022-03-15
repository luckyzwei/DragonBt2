/*
@Time : 2020/5/11 0:56
@Author : 96121
@File : proto_server
@Software: GoLand
*/
package server

import (
	"errors"
	"master/utils"
)

type RPC_ServerEvent struct {
	EventCode int    `json:"code"`    //! 事件类型
	UId       int64  `json:"uid"`     //! 触发UID
	Target    int64  `json:"target"`  //! 目标UID
	Param1    int    `json:"param1""` //! 参数1
	Param2    string `json:"param2"`  //! 参数2
}

const (
	MAX_MSG_GET = 10
)

type RPC_Server struct {
}

type RPC_RegServerReq struct {
	ID     int    //! 服务器Id
	Name   string //! 名字
	Online int    //! 在线人数
}

type RPC_RegServerRes struct {
	RetCode int //! 返回错误码
}

type RPC_ReqEventArrRes struct {
	RetCode  int                          //!
	EventArr [MAX_MSG_GET]RPC_ServerEvent //! 消息队列
}

//! 注册服务器
func (self *RPC_Server) RegServer(req RPC_RegServerReq, res *RPC_RegServerRes) error {
	server := GetServerMgr().GetServer(req.ID, true)
	if server != nil {
		server.Data.Online = req.Online
		if server.Data.Name == "" {
			server.Data.Name = req.Name
		}
	}

	//utils.LogDebug("Server Reg Info:", req.ID, req.Name, req.Online)
	return nil
}

func (self *RPC_Server) ReqEvent(sid int, res *RPC_ServerEvent) error {
	server := GetServerMgr().GetServer(sid, true)
	if server != nil {
		evt := server.PopEvent()
		if evt != nil {
			res.EventCode = evt.EventCode
			res.UId = evt.UId
			res.Target = evt.Target
			res.Param1 = evt.Param1
			res.Param2 = evt.Param2

			return nil
		}
	}

	//utils.LogDebug("Server Req Event :", sid)

	return errors.New("NULL")
}

func (self *RPC_Server) ReqEventArr(sid int, res *RPC_ReqEventArrRes) error {
	utils.LogDebug("Server Req Event Arr:", sid)

	server := GetServerMgr().GetServer(sid, false)
	if server != nil {
		evtNum := 0
		for i := 0; i < MAX_MSG_GET; i++ {
			evt := server.PopEvent()
			if evt != nil {
				res.EventArr[i].EventCode = evt.EventCode
				res.EventArr[i].UId = evt.UId
				res.EventArr[i].Target = evt.Target
				res.EventArr[i].Param1 = evt.Param1
				res.EventArr[i].Param2 = evt.Param2
				evtNum++
			}
		}
		if evtNum > 0 {
			return nil
		} else {
			return errors.New("NULL")
		}
	}

	return errors.New("NULL")
}
