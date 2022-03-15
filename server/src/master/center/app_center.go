package center

import (
	"master/center/chat"
	"master/center/match"
	"master/center/player"
	"master/center/server"
	"master/center/tower"
	"master/center/union"
	"master/core"
	"master/utils"
	"net/rpc"
	"time"
)

const (
	PER_SAVE_TIME      = 60 //! 60s保存一次
	INIT_WORLD_CHANNEL = 10 //! 默认初始化10个聊天频道

)

//! 登录模块
type CenterApp struct {
	StartTime int64 //! 服务器开始时间
	LastSave  int64 //! 上次保存时间

	EventArr chan *server.ServerEvent //! 事件推送
}

var s_centerapp *CenterApp

func GetCenterApp() *CenterApp {
	if s_centerapp == nil {
		s_centerapp = new(CenterApp)
		s_centerapp.StartTime = time.Now().Unix()
		s_centerapp.LastSave = time.Now().Unix()
		s_centerapp.EventArr = make(chan *server.ServerEvent, 10000)
		core.CenterApp = s_centerapp
	}

	return s_centerapp
}

//! 注册RPC服务,该消息是阻塞模式
func (self *CenterApp) RegisterService() {
	//! 采用HTTP作为调用载体
	rpc.HandleHTTP()

	//! 注册远程服务对象
	rpc.Register(new(chat.RPC_Chat))     //! 聊天
	rpc.Register(new(player.RPC_Player)) //! 角色
	rpc.Register(new(player.RPC_Friend)) //! 好友
	rpc.Register(new(player.RPC_Union))  //! 工会
	rpc.Register(new(player.RPC_Tower))  //! 爬塔
	rpc.Register(new(server.RPC_Server)) //! 服务器，事件
	rpc.Register(new(match.RPC_Match))   //! 跨服竞技
}

func (self *CenterApp) StartService() {
	//! 初始化10个世界聊天频道，后续会根据需求进行动态添加
	for i := 0; i < INIT_WORLD_CHANNEL; i++ {
		self.AddEvent(0, core.SYSTEM_EVENT_CHAT_WORLD_OPEN, 0, 0, 0, "")
	}

	player.GetPlayerMgr().Init()

	//! 初始化所有的公会
	union.GetUnionMgr().GetAllData()
	tower.GetTowerMgr().GetAllData()
	match.GetGeneralMgr().GetAllData()
	match.GetCrossArenaMgr().GetAllData()
	match.GetConsumerTopMgr().GetAllData()

	////! 注册TCP，端口
	//handle, err := net.Listen("tcp", "127.0.0.1:9000")
	//if err != nil {
	//	log.Fatalln("listen rpc fatal error: ", err)
	//}
	//
	//log.Println("Start RPC Service...", 9000)
	//http.Serve(handle, nil)
}

//! 保存所有的数据
func (self *CenterApp) StopService() {
	//! 公会 保存数据
	union.GetUnionMgr().OnSave()

	//! 爬塔 保存数据
	tower.GetTowerMgr().OnSave()
}

func (self *CenterApp) OnTimer() {
	self.LastSave = time.Now().Unix()
	//! 定时逻辑
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		self.OnLogic()
	}

	ticker.Stop()
}

//! 逻辑处理
func (self *CenterApp) OnLogic() {
	tNow := time.Now().Unix()
	if tNow-self.LastSave > PER_SAVE_TIME {
		self.LastSave = tNow

		//! 服务器数据保存
		server.GetServerMgr().OnSave()

		//! 角色数据保存
		player.GetPlayerMgr().OnSave()

		//! 公会数据保存
		union.GetUnionMgr().OnSave()

		//跨服限时神将信息保存
		match.GetGeneralMgr().OnSave()
		//跨服竞技场信息保存
		match.GetCrossArenaMgr().OnSave()
		//无双神将
		match.GetConsumerTopMgr().OnSave()

		//! 事件频道检查
		self.CheckWorldChannel()
	}
}

//! 检查聊天频道负载
func (self *CenterApp) CheckWorldChannel() {
	worldChannelNum := chat.GetChatMgr().GetWorldCount()
	worldPlayerNum := chat.GetChatMgr().GetPlayerCount()

	if worldChannelNum*chat.MAX_PLAYER_NUM*80/100 < worldPlayerNum {
		//! 超过总负载人数的80%，则增加2个频道
	}
	for i := 0; i < worldChannelNum; i++ {
		channel := chat.GetChatMgr().GetWorldChannel(i + 1)
		if channel != nil {
			if channel.GetPlayerCount() > chat.MAX_PLAYER_NUM*90/100 {
				//! 超过90%则补充频道
			}
		}
	}
}

func (self *CenterApp) AddEvent(sid int, code int, uid int64, target int64, param1 int, param2 string) {
	if sid > 0 {
		server := server.GetServerMgr().GetServer(sid, true)
		if server != nil {
			server.PushEvent(code, uid, target, param1, param2)
		}
	} else {
		//! 中心服事件
		evt := &server.ServerEvent{
			EventCode: code,
			UId:       uid,
			Target:    target,
			Param1:    param1,
			Param2:    param2,
		}

		self.EventArr <- evt
	}

}

func (self *CenterApp) AddEvnet(int, interface{}) {

}

//! 100ms处理一次事件
func (self *CenterApp) OnLogicEvent() {
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ticker.C:
		case event := <-self.EventArr:
			self.ProcessEvent(event)
		}
	}

	ticker.Stop()
}

func (self *CenterApp) ProcessEvent(evt *server.ServerEvent) {
	utils.LogDebug("Proc Event...", evt.EventCode, evt.Target, evt.Param1, evt.Param2)
	switch evt.EventCode {
	case core.SYSTEM_EVENT_CHAT_UNION_OPEN:
		chat.GetChatMgr().AddUnionChannel(evt.Param1)
	case core.SYSTEM_EVENT_CHAT_WORLD_OPEN:
		chat.GetChatMgr().AddWorldChannel()
	}
}
