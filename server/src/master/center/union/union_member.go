/*
@Time : 2020/5/8 23:40
@Author : 96121
@File : union_member
@Software: GoLand
*/
package union

type UserActivityRecord struct {
	Time     int64 `json:"time"`     // 时间
	AddCount int   `json:"addcount"` // 数量
}
type JS_UnionMember struct {
	Uid            int64                 `json:"uid"`
	Level          int                   `json:"level"`
	Uname          string                `json:"uname"`
	Iconid         int                   `json:"iconid"`
	Portrait       int                   `json:"portrait"`
	Vip            int                   `json:"vip"`
	Position       int                   `json:"position"`
	Fight          int64                 `json:"fight"`
	Lastlogintime  int64                 `json:"lastlogintime"`
	ActivityRecord []*UserActivityRecord `json:"activityrecord"` // 记录
	BraveHand      int                   `json:"bravehand"`      // 无畏之手
	Stage          int                   `json:"stage"`          // 关卡进度
	ServerID       int                   `json:"serverid"`
}

type JS_UnionApply struct {
	Uid       int64  `json:"uid"`
	Level     int    `json:"level"`
	Uname     string `json:"uname"`
	Iconid    int    `json:"iconid"`
	Portrait    int    `json:"portrait"`
	Vip       int    `json:"vip"`
	Fight     int64  `json:"fight"`
	Applytime int64  `json:"lastlogintime"`
	ServerID  int    `json:"serverid"`
}
