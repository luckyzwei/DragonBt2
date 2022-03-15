package gate

//! client2server
//! 版本验证
type C2S_CtrlHead struct {
	Ctrl  string `json:"ctrl"`
	Uid   int64  `json:"uid"`
	Os    string `json:"os"`
	Ver   int    `json:"ver"`
	MsgId int    `json:"msgid"`
}

//! 注册
type C2S_Reg struct {
	Account           string `json:"account"`
	Password          string `json:"password"`
	ServerId          int    `json:"serverid"`
	Platform_os       string `json:"os"`
	Platform_Brand    string `json:"brand"`
	Platform_DeviceId string `json:"deviceid"`
	Platform_Model    string `json:"model"`
}

//! sdk登陆
type C2S_SDKLogin struct {
	Password          string `json:"password"`
	ServerId          int    `json:"serverid"`
	Username          string `json:"username"`
	Third             string `json:"third"`
	Platform_os       string `json:"os"`
	Platform_Brand    string `json:"brand"`
	Platform_DeviceId string `json:"deviceid"`
	Platform_Model    string `json:"model"`
}

//! 重连
type C2S_AutoLogin struct {
	CheckCode string `json:"checkcode"`
}

//! server2client
//! msgid
type S2C_MsgId struct {
	Cid     string `json:"cid"`
	CurTime int64  `json:"curtime"`
}

//! 注册
type S2C_Reg struct {
	Cid      string `json:"cid"`
	Uid      int64  `json:"uid"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Creator  string `json:"creator"`
}

//! 获取平台信息
type S2C_PlatFromInfo struct {
	Cid string `json:"cid"`
}

//! 发送结果
type S2C_ResultMsg struct {
	Cid string `json:"cid"`
	Ret int    `json:"ret"`
}

//! 发送结果2
type S2C_Result2Msg struct {
	Cid string `json:"cid"`
}

//! 进入排队
type S2C_LineUp struct {
	Cid        string `json:"cid"`
	WaitCount  int    `json:"waitcount"`
	PassMinute int    `json:"passminute"`
	PassSecond int    `json:"passsecond"`
}

//! center服ctrl
type S2S_CenterCid struct {
	Cid string `json:"cid"`
}
