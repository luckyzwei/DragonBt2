package game

import "fmt"

// 关卡事件
type ModExample struct {
	player *Player
	Data   SanExample //! 数据库数据
}

type SanExample struct {
	Uid int64
	DataUpdate
}

func (m *ModExample) Decode() {

}

func (m *ModExample) Encode() {

}

func (m *ModExample) getTableName() string {
	return "san_userexample"
}

func (m *ModExample) init(uid int64) {
	m.Data.Uid = uid
	m.CheckInfo()
}

func (m *ModExample) CheckInfo() {

}

func (m *ModExample) OnGetData(player *Player) {
	m.player = player
}

func (m *ModExample) OnGetOtherData() {
	tableName := m.getTableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", tableName, m.player.ID)
	GetServer().DBUser.GetOneData(sql, &m.Data, tableName, m.player.ID)
	if m.Data.Uid <= 0 {
		m.init(m.player.ID)
		m.CheckInfo()
		m.Encode()
		InsertTable(tableName, &m.Data, 0, true)
	} else {
		m.Decode()
		m.CheckInfo()
	}

	m.Data.Init(tableName, &m.Data, true)
}

func (m *ModExample) OnSave(sql bool) {
	m.Encode()
	m.Data.Update(sql)
}

func (m *ModExample) OnRefresh() {

}

// 老的消息处理
func (m *ModExample) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (m *ModExample) onReg(handlers map[string]func(body []byte)) {
	handlers["exmaple_action"] = m.ExampleAction

}

func (m *ModExample) ExampleAction(body []byte) {

}

type S2C_ExampleInfo struct {
	Cid string `json:"cid"`
}

// 登录时同步消息
func (m *ModExample) sendInfo() {
	const cid = "exampleinfo"
	msg := &S2C_ExampleInfo{
		Cid: cid,
	}
	m.player.Send(cid, msg)
}
