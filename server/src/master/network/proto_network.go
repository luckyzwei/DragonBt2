package network

//! 登陆成功
type S2C_LoginRet struct {
	Cid        string `json:"cid"`
	Ret        int    `json:"ret"`
	CheckCode  string `json:"checkcode"`
	Servertime int64  `json:"servertime"`
}

//! 进入排队
type S2C_LineUp struct {
	Cid        string `json:"cid"`
	WaitCount  int    `json:"waitcount"`
	PassMinute int    `json:"passminute"`
	PassSecond int    `json:"passsecond"`
}

type SendRZ_EnvInfo_DevInfo struct {
	DeviceId string `json:"deviceId"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Os       string `json:"os"`
}

type SendRZ_EnvInfo_DevInfo_Ios struct {
	DeviceId string `json:"deviceId"`
	UUID     string `json:"uuid"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Os       string `json:"os"`
	Fr       string `json:"Fr"`
	Res      string `json:"res"`
	Net      string `json:"Net"`
	Mac      string `json:"mac"`
	Operator string `json:"operator"`
	Ip       string `json:"ip"`
}

//渠道信息
type SendRZ_EnvInfo_ChInfo_Ios struct {
	Ch    string `json:"ch"` //Sdk来源，如大圣，火速，大禹，天拓
	SubCh string `json:"subCh"`
}

//! 跑马灯
type S2C_Barrage struct {
	Cid   string `json:"cid"`
	Size  int    `json:"size"`
	Text  string `json:"text"`
	Red   int    `json:"red"`
	Green int    `json:"green"`
	Blue  int    `json:"blue"`
	Uid   int64  `json:"uid"`
}

//! 发送错误信息
type S2C_ErrInfo struct {
	Cid  string `json:"cid"`
	Info string `json:"info"`
}

//! 发送结果2
type S2C_Result2Msg struct {
	Cid string `json:"cid"`
}

type S2C_MsgId struct {
	Cid     string `json:"cid"`
	CurTime int64  `json:"curtime"`
}
