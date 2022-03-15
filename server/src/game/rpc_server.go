/*
@Time : 2020/5/11 0:56
@Author : 96121
@File : proto_server
@Software: GoLand
*/
package game

import (
	"net/rpc"
)

const (
	RPC_REGSERVER   = "RPC_Server.RegServer"   //! 注册服务器
	RPC_REQEVENT    = "RPC_Server.ReqEvent"    //! 请求事件
	RPC_REQEVENTARR = "RPC_Server.ReqEventArr" //! 事件列表
)

const (
	RPC_EVENT_MAX = 20000 //! 消息队列 = 20000个
	MAX_MSG_GET   = 10    //! 每次请求消息队列个数
)

type RPC_ServerEvent struct {
	EventCode int    `json:"code"`   //! 事件类型
	UId       int64  `json:"uid"`    //! 触发UID
	Target    int64  `json:"target"` //! 目标UID
	Param1    int    `json:"param1"` //! 参数1
	Param2    string `json:"param2"` //! 参数2
}

type RPC_Server struct {
	Client    *rpc.Client           //! RPC调用
	EventChan chan *RPC_ServerEvent //! 服务器消息队列
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

func (self *RPC_Server) Init() {
	self.EventChan = make(chan *RPC_ServerEvent)
}

func (self *RPC_Server) CheckErr(err error) {
	if self.Client == nil {
		GetMasterMgr().ResumeService(true)
	}

	//! 断开连接返回
	if err == nil {
		return
	}

	if err.Error() == "connection is shut down" {
		//! 连接丢失,重新初始化
		GetMasterMgr().ResumeService(true)
	}
}

//! 注册服务器
func (self *RPC_Server) RegServer(sid int, online int, sname string) bool {
	if self.Client != nil {
		var req RPC_RegServerReq
		req.ID = sid
		req.Online = online
		req.Name = sname

		var res RPC_RegServerRes

		req.ID = sid
		req.Name = sname
		req.Online = online

		err := GetMasterMgr().CallEx(self.Client,RPC_REGSERVER, req, &res)
		if err != nil {
			LogDebug("中心服->注册服务器失败：", err.Error())
			self.CheckErr(err)
		} else {
			retCode := res.RetCode
			switch retCode {
			case 0:
				//LogDebug("中心服->同步OK")
			}
		}
	} else {
		//! 连接丢失，则直接检查
		self.CheckErr(nil)
	}

	return false

	//req RPC_RegServerReq, res *RPC_RegServerRes
	//server := GetServerMgr().GetServer(req.ID, true)
	//if server != nil {
	//	server.Data.Online = req.Online
	//}
	//
	//return nil
}

func (self *RPC_Server) ReqEvent(sid int) bool {
	if self.Client != nil {
		var res RPC_ServerEvent

		err := GetMasterMgr().CallEx(self.Client,RPC_REQEVENT, sid, &res)
		if err == nil {
			LogDebug("req server event .. ", sid)
			self.EventChan <- &res
			return true
		} else {
			self.CheckErr(err)
		}
	} else {
		return GetMasterMgr().ResumeService(true)
	}

	return false

	//server := GetServerMgr().GetServer(sid, true)
	//if server != nil {
	//	evt := server.PopEvent()
	//	if evt != nil {
	//		res.EventCode = evt.EventCode
	//		res.Target = evt.Target
	//		res.Param1 = evt.Param1
	//		res.Param2 = evt.Param2
	//
	//		return nil
	//	}
	//}
	//
	//return errors.New("NULL")
	//
	//return nil
}

func (self *RPC_Server) ReqEventArr(sid int) bool {
	//LogDebug("req server event arr .. ", sid)
	if self.Client != nil {
		var res RPC_ReqEventArrRes
		err := GetMasterMgr().CallEx(self.Client,RPC_REQEVENTARR, sid, &res)
		if err == nil {
			eventArrLen := 0
			for i := 0; i < len(res.EventArr); i++ {
				if res.EventArr[i].EventCode > 0 {
					self.EventChan <- &res.EventArr[i]
					eventArrLen++
				}
			}

			if eventArrLen > 0 {
				return true
			}
		} else {
			self.CheckErr(err)
		}
	}

	return false

	//server := GetServerMgr().GetServer(sid, false)
	//if server != nil {
	//	evtNum := 0
	//	for i := 0; i < 10; i ++ {
	//		evt := server.PopEvent()
	//		if evt != nil {
	//			res[i].EventCode = evt.EventCode
	//			res[i].Target = evt.Target
	//			res[i].Param1 = evt.Param1
	//			res[i].Param2 = evt.Param2
	//			evtNum++
	//		}
	//	}
	//	if evtNum > 0 {
	//		return nil
	//	} else {
	//		return errors.New("NULL")
	//	}
	//}
	//
	//return errors.New("NULL")
	//
	//return nil
}
