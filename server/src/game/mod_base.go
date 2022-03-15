package game

type ModBase interface {
	OnGetData(player *Player)            // 得到数据,缓存时会读取
	OnGetOtherData()                     // 得到数据,缓存时不会读取
	OnMsg(ctrl string, body []byte) bool // 模块收到消息
	OnSave(sql bool)                     // 保存数据
}
