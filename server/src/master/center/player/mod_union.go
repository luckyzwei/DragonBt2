/*
@Time : 2020/5/8 22:54 
@Author : 96121
@File : mod_union
@Software: GoLand
*/
package player

import (
	"master/db"
)

//! 好友数据
type JS_UnionInfo struct {
	Uid       int64  //! 好友Uid
	FId       int64  //! 好友Id
	FName     string //! 好友昵称
	Level     int    //! 好友等级
	Fight     int    //! 战力
	ServerId  int    //! 区服Id
	Sex       int    //! 性别
	Vip       int    //! VIP
	Icon      int    //! 头像
	Portrait  int    //! 头像边框
	LastLogin int64  //! 上次登录时间
}

//! 好友数据库结构
type SQL_UnionInfo struct {
	Uid     int64  //! 好友Uid
	Friends string //! 好友列表
	Applys  string //! 申请列表

	friends map[int64]*JS_Friend //! 好友数据
	applys  map[int64]*JS_Friend //! 申请数据
	db.DataUpdate                //! 数据库操作结构
}

//! 加密数据
func (self *SQL_UnionInfo) Decode() {

}

//! 加密数据
func (self *SQL_UnionInfo) Encode() {

}

//! 好友模块
type ModUnion struct {
	Data SQL_UnionInfo //! 公会数据
}

func (self *ModUnion) onGetData(uid int64) {
	self.Data.Uid = uid
	///sql := fmt.Sprintf("select * from user_unioninfo where uid = %d", self.Data.Uid)
}

func (self *ModUnion) onSave() {

}
