package main

//"fmt"

type San_Player struct {
	Uid     int64
	Uname   string
	Camp    int
	Unionid int
}

type PlayerMgr struct {
	PlayerMap  map[int64]*San_Player
	SessionObj *Session
}

func (self *PlayerMgr) GetPlayer(uid int64) *San_Player {
	return self.PlayerMap[uid]
}

func (self *PlayerMgr) AddPlayer(player *San_Player) {
	if self.GetPlayer(player.Uid) == nil {
		self.PlayerMap[player.Uid] = player
	}
}
