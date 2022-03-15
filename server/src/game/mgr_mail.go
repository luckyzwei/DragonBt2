package game

import (
	"encoding/json"
	"sync"
	"time"
)

type Sql_GMail struct {
	Id   int
	Info string
	info JS_Mail
	DataUpdate
	Locker *sync.RWMutex
}

type JS_Mail struct {
	Title    string         `json:"title"`   //! 标题
	Content  string         `json:"content"` //! 内容
	Item     []PassItem     `json:"item"`    //! 道具
	Users    map[string]int `json:"users"`   //! 哪些人领取过
	Time     int64          `json:"time"`    //! 时间
	Sender   string         `json:"sender"`  //! 发送者
	MinLevel int            `json:"minlevel"`
	MaxLevel int            `json:"maxlevel"`
}

type MailMgr struct {
	Sql_Mail []*Sql_GMail

	Locker *sync.RWMutex
}

func (self *Sql_GMail) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Info), &self.info)
}

func (self *Sql_GMail) Encode() { //! 将data数据写入数据库
	self.Info = HF_JtoA(&self.info)
}

var mailsingleton *MailMgr = nil

//! public
func GetMailMgr() *MailMgr {
	if mailsingleton == nil {
		mailsingleton = new(MailMgr)
		mailsingleton.Locker = new(sync.RWMutex)
		mailsingleton.Sql_Mail = make([]*Sql_GMail, 0)
	}

	return mailsingleton
}

func (self *MailMgr) Save() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	for i := 0; i < len(self.Sql_Mail); i++ {
		self.Sql_Mail[i].Encode()
		self.Sql_Mail[i].Update(true)
	}
}

func (self *MailMgr) GetData() {
	var dip Sql_GMail
	res := GetServer().DBUser.GetAllData("select * from `san_mail`", &dip)

	for i := 0; i < len(res); i++ {
		data := res[i].(*Sql_GMail)
		data.Decode()
		data.Init("san_mail", data, false)
		self.Sql_Mail = append(self.Sql_Mail, data)
	}
}

func (self *MailMgr) AddMail(title string, body string, items []PassItem, sender string, minlevel int, maxlevel int) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	mail := new(Sql_GMail)
	mail.info.Title = title
	mail.info.Content = body
	mail.info.Item = items
	mail.info.Time = TimeServer().Unix()
	mail.info.Sender = sender
	mail.info.MinLevel = minlevel
	mail.info.MaxLevel = maxlevel

	mail.Encode()
	mail.Id = int(InsertTable("san_mail", mail, 1, false))
	mail.Init("san_mail", mail, false)
	self.Sql_Mail = append(self.Sql_Mail, mail)

	//LogDebug("发送全局邮件：", mail.Id)
	GetPlayerMgr().AddMail(&mail.info, mail.Id)
}

func (self *MailMgr) GetMail(player *Player) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	mod := player.GetModule("mail").(*ModMail)
	playerLv := player.Sql_UserBase.Level
	for _, v := range self.Sql_Mail {
		rtime, _ := time.ParseInLocation(DATEFORMAT, player.Sql_UserBase.Regtime, time.Local)
		if rtime.Unix() > v.info.Time { //! 发邮件时还未注册
			continue
		}
		ok := mod.IsRecvSystemMail(v.Id)
		if ok {
			continue
		}
		mod.SetRecvSystemMail(v.Id)
		if playerLv < v.info.MinLevel {
			continue
		}
		if v.info.MaxLevel != 0 && playerLv > v.info.MaxLevel {
			continue
		}
		mod.AddMail(1, 2, 0, v.info.Title, v.info.Content, v.info.Sender, v.info.Item, false, v.info.Time)
	}
}
