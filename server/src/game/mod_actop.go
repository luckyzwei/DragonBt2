package game

import (
	"fmt"
)

// 活动消费排行榜, san_cost1
type ModActop struct {
	player *Player
	Data   DBActop //! 数据库结构
}

// 消费活动
type DBActop struct {
	Uid       int64
	Cost1     int
	Cost2     int
	Cost3     int
	Cost4     int
	Cost5     int
	TaskType1 int
	TaskType2 int
	TaskType3 int
	TaskType4 int
	TaskType5 int
	Step1     int
	Step2     int
	Step3     int
	Step4     int
	Step5     int

	DataUpdate
}

func (self *ModActop) Decode() {

}

func (self *ModActop) Encode() {

}

func (self *ModActop) OnGetData(player *Player) {
	self.player = player
}

func (self *ModActop) GetTableName() string {
	return "san_cost"
}

func (self *ModActop) OnGetOtherData() {
	name := self.GetTableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", name, self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Data, name, self.player.ID)
	if self.Data.Uid <= 0 {
		self.Data.Uid = self.player.ID
		self.Encode()
		InsertTable(name, &self.Data, 0, true)
	} else {
		self.Decode()
	}

	self.Data.Init(name, &self.Data, true)
}

func (self *ModActop) OnSave(sql bool) {
	self.Encode()
	self.Data.Update(sql)
}

func (self *ModActop) OnRefresh() {

}

func (self *ModActop) OnMsg(ctrl string, body []byte) bool {
	return false
}
