package core

//! 全局接口变量
var (
	MasterApp IMasterApp = nil
	CenterApp ICenterApp = nil
	PlayerMgr IPlayerMgr = nil
	GateApp   IGateApp   = nil
)

//! 对外接口，需保证获取全局句柄
func GetMasterApp() IMasterApp {
	return MasterApp
}

func GetPlayerMgr() IPlayerMgr {
	return PlayerMgr
}

func GetCenterApp() ICenterApp {
	return CenterApp
}

func GetGateApp() IGateApp {
	return GateApp
}
