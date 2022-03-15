package game

import (
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type LineUpMgr struct {
	Locker        *sync.RWMutex      //! 数据锁
	WaitTotal     int                //! 等待人数
	LoginTotal    int                //! 登录中人数
	WaitUser      map[int64]*Session //! 等待登录-sessionid
	LoginUser     map[int64]*Session //! 登录中
	WaitArr       []int64            //! 排队队列
	PassMinute    int                //! 每分登录人数
	PassMinuteNow int                //! 统计
	PassSecond    int                //! 每秒登录人数
	MaxLogin      int                //! 最大同时登录人数
}

var lineupmgrsingleton *LineUpMgr = nil

func GetLineUpMgr() *LineUpMgr {
	if lineupmgrsingleton == nil {
		lineupmgrsingleton = new(LineUpMgr)
		lineupmgrsingleton.Locker = new(sync.RWMutex)
		lineupmgrsingleton.WaitUser = make(map[int64]*Session)
		lineupmgrsingleton.LoginUser = make(map[int64]*Session)
		lineupmgrsingleton.WaitArr = make([]int64, 0)
		lineupmgrsingleton.PassMinute = 0
		lineupmgrsingleton.PassSecond = 0
		lineupmgrsingleton.MaxLogin = 100
	}

	return lineupmgrsingleton
}

//! 增加排队客户端
func (self *LineUpMgr) AddWaitClient(session *Session) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.WaitTotal++
	self.WaitUser[session.ID] = session
	self.WaitArr = append(self.WaitArr, session.ID)

	var waitLogin S2C_LineUp
	waitLogin.Cid = "lineup"
	waitLogin.WaitCount = 0
	waitLogin.PassMinute = self.PassMinute
	waitLogin.PassSecond = self.PassSecond
	session.SendMsg("lineup", HF_JtoB(&waitLogin))
}

//! 取消排队
func (self *LineUpMgr) CancelLogin(session *Session) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.WaitTotal--
	delete(self.WaitUser, session.ID)
}

//! 增加登录客户端, 暂时没有调用到.
func (self *LineUpMgr) AddLoginClient(session *Session) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	//! 如果超出登录人数上线，直接进入排队
	if len(self.LoginUser) > self.MaxLogin {
		self.WaitTotal++
		self.WaitUser[session.ID] = session

		var waitLogin S2C_LineUp
		waitLogin.Cid = "lineup"
		waitLogin.WaitCount = 0
		waitLogin.PassMinute = self.PassMinute
		waitLogin.PassSecond = self.PassSecond
		session.SendMsg("lineup", HF_JtoB(&waitLogin))
	}

	self.LoginUser[session.ID] = session
}

func (self *LineUpMgr) RemoveClient(session *Session) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	_, ok := self.WaitUser[session.ID]
	if ok {
		self.WaitTotal--
		delete(self.WaitUser, session.ID)
	}

	_, ok1 := self.LoginUser[session.ID]
	if ok1 {
		self.LoginTotal--
		delete(self.LoginUser, session.ID)
	}
}

func (self *LineUpMgr) RemoveClientSelf(session *Session) {
	//self.Locker.Lock()
	//defer self.Locker.Unlock()
	_, ok := self.WaitUser[session.ID]
	if ok {
		self.WaitTotal--
		delete(self.WaitUser, session.ID)
	}

	_, ok1 := self.LoginUser[session.ID]
	if ok1 {
		self.LoginTotal--
		delete(self.LoginUser, session.ID)
	}
}

//! 登录完成
func (self *LineUpMgr) CompleteLogin(sid int64) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	_, ok1 := self.LoginUser[sid]
	if ok1 {
		self.LoginTotal--
		self.PassMinuteNow++
		self.PassSecond++
		delete(self.LoginUser, sid)
	}
}

//! 登录完成
func (self *LineUpMgr) AutoCompleteLogin(sid int64) {
	//self.Locker.Lock()
	//defer self.Locker.Unlock()
	_, ok1 := self.LoginUser[sid]
	if ok1 {
		self.LoginTotal--
		self.PassMinuteNow++
		self.PassSecond++
		delete(self.LoginUser, sid)
	}
}

//!
func (self *LineUpMgr) Run() {
	ticker := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ticker.C:
			self.Fresh()
		}
	}
}

//! 同步客户端
func (self *LineUpMgr) Fresh() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	self.Locker.Lock()
	defer self.Locker.Unlock()

	if GetServer().ShutDown {
		return
	}

	if len(self.WaitArr) == 0 {
		return
	}

	//! 处理登录后的消息
	//online, _ := GetPlayerMgr().GetOnline()
	onlineplayer := GetPlayerMgr().GetPlayerOnline()
	//onlineplayer = len(online)
	LogInfo("处理排队信息", len(self.WaitArr), onlineplayer, GetServer().Con.NetworkCon.MaxPlayer)
	if onlineplayer < GetServer().Con.NetworkCon.MaxPlayer {
		for i := 0; i < GetServer().Con.NetworkCon.MaxPlayer-onlineplayer; i++ {
			if len(self.WaitArr) > 0 {
				sessionid := self.WaitArr[0]
				session, ok := self.WaitUser[sessionid]
				if !ok {
					self.WaitArr = self.WaitArr[1:]
					continue
				}

				var waitLogin S2C_LineUp
				waitLogin.Cid = "lineup"
				waitLogin.WaitCount = 0
				session.SendMsg("lineup", HF_JtoB(&waitLogin))

				//! 自动登录
				session.AutoLogin()

				//! 删除等待
				self.WaitArr = self.WaitArr[1:]
			} else {
				break
			}
		}
	}

	var waitLogin S2C_LineUp
	waitLogin.Cid = "lineup"
	for i := 0; i < len(self.WaitArr); i++ {
		session, ok := self.WaitUser[self.WaitArr[i]]
		if !ok {
			continue
		}
		waitLogin.WaitCount = i + 1
		waitLogin.PassMinute = self.PassMinute
		waitLogin.PassSecond = self.PassSecond

		session.SendMsg("lineup", HF_JtoB(&waitLogin))
	}

	self.PassSecond = 0
	if TimeServer().Second() == 0 {
		self.PassMinute = self.PassMinuteNow
		self.PassMinuteNow = 0
	}
}

//!
func (self *LineUpMgr) Add(player *Player, boxid int) {

}
