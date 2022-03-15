/*
@Time : 2020/4/22 10:18
@Author : 96121
@File : mod_player
@Software: GoLand
*/
package player

import (
	"encoding/json"
	"fmt"
	"game"
	"master/center/union"
	"master/core"
	"master/db"
	"master/utils"
	"time"
)

const (
	PLAYER_TABLE_SQL  = "select * from `%s` where uid = %d" //! 数据表地址
	PLAYER_TABLE_NAME = "tbl_user"                          //! 数据库表
)

//! 角色信息结构
type JS_PlayerData struct {
	UId            int64              `json:"uid"`            //! 角色Id
	UName          string             `json:"uname"`          //! 昵称
	Level          int                `json:"level"`          //! 等级
	Fight          int64              `json:"fight"`          //! 战力
	PassId         int                `json:"passid"`         //! 最大关卡
	ServerId       int                `json:"serverid"`       //! 服务器ID
	Sex            int                `json:"sex"`            //! 性别
	IconId         int                `json:"iconid"`         //! 头像
	Portrait       int                `json:"portrait"`       //! 头像框
	Vip            int                `json:"vip"`            //! VIP
	RegTime        int64              `json:"reg"`            //! 注册时间
	LoginTime      int64              `json:"login"`          //! 登录时间
	LastUpdate     int64              `json:"update"`         //! 上次更新
	UnionId        int                `json:"unionid"`        //! 军团id
	ArenaFightInfo *game.JS_FightInfo `json:"arenafightinfo"` //! JJC防守阵容
	UnionName      string             `json:"unionname"`      //! 公会名
	Position       int                `json:"position"`       //! 职位
	BraveHand      int                `json:"bravehand"`      //! 无畏之手
	UserSignature  string             `json:"usersignature"`  //! 签名
	UserID         string             `json:"userid"`         //! sdk分配id
}

//! 英雄结构
type JS_PlayerHero struct {
	HeroId          int           `json:"heroid"`
	Level           int           `json:"level"`
	Star            int           `json:"star"`
	ArtifactId      int           `json:"artifactid"`      //! 神器ID
	ArtifactLv      int           `json:"artifactlv"`      //! 神器等级
	ExclusiveId     int           `json:"ExclusiveId"`     //! 专属装备Id
	ExclusiveLv     int           `json:"ExclusiveLv"`     //! 专属装备等级
	ExclusiveUnLock int           `json:"ExclusiveUnLock"` //! 专属装备解锁状态0未解锁  1解锁
	Skin            int           `json:"skin"`            //！皮肤
	Attr            map[int]int64 `json:"attr"`            //！展示属性
}

//! 英雄装备
type JS_HeroEquip struct {
	ItemId int `json:"itemid"`
	Level  int `json:"level"`
}

//! 生命树
type JS_LifeTreeInfo struct {
	MainLevel int            `json:"mainlevel"`
	Info      []*JS_LifeTree `json:"info"`
}

type JS_LifeTree struct {
	Type  int `json:"type"`  // 类型
	Level int `json:"level"` // 等级
}

//! 数据库结构
type SQL_PlayerData struct {
	UId        int64  //! 角色ID
	UName      string //! 昵称
	Level      int    //! 等级
	Fight      int    //! 战力
	PassId     int    //! 最大关卡
	ServerId   int    //! 服务器ID那
	RegTime    int64  //! 注册时间
	LoginTime  int64  //! 登录时间
	LastUpdate int64  //! 上次更新时间
	Data       string //! 数据
	Heros      string //! 武将列表
	Equips     string //! 装备列表
	LifeTree   string //! 生命树

	data     *JS_PlayerData    //! 基础数据
	heros    []*JS_PlayerHero  //! 武将列表
	equips   [][]*JS_HeroEquip //! 装备列表
	lifetree *game.JS_LifeTreeInfo  //! 生命树
	db.DataUpdate              //! 数据库结构
}

func (self *SQL_PlayerData) Decode() {
	json.Unmarshal([]byte(self.Data), &self.data)
	json.Unmarshal([]byte(self.Heros), &self.heros)
	json.Unmarshal([]byte(self.Equips), &self.equips)
	json.Unmarshal([]byte(self.LifeTree), &self.lifetree)
}

func (self *SQL_PlayerData) Encode() {
	self.Data = utils.HF_JtoA(&self.data)
	self.Heros = utils.HF_JtoA(&self.heros)
	self.Equips = utils.HF_JtoA(&self.equips)
	self.LifeTree = utils.HF_JtoA(&self.lifetree)
}

//! 角色信息结构
type Player struct {
	Data       SQL_PlayerData //! 角色数据
	Online     int            //! 是否在线
	NeedUpdate bool           //! 数据更新
	NeedSave   bool           //! 需要保存

	DataFriend *ModFriend //! 好友模块
}

//! 载入数据
func (self *Player) onGetData(uid int64) {
	sql := fmt.Sprintf(PLAYER_TABLE_SQL, PLAYER_TABLE_NAME, uid)
	ret := db.GetDBMgr().DBUser.GetOneData(sql, &self.Data, PLAYER_TABLE_NAME, uid)
	if ret == true {
		if self.Data.UId <= 0 {
			self.Data.UId = uid
			self.Data.data = new(JS_PlayerData)
			self.Data.heros = make([]*JS_PlayerHero, 0)
			self.Data.equips = make([][]*JS_HeroEquip, 0)
			self.Data.lifetree = new(game.JS_LifeTreeInfo)
			self.Online = game.LOGIC_TRUE

			myunion := union.GetUnionMgr().GetUnion(self.Data.data.UnionId)
			if myunion != nil {
				member := myunion.GetMember(self.GetUid())
				if member != nil {
					if member.Lastlogintime != 0 {
						member.Lastlogintime = 0
					}
				}
			}

			self.Data.Encode()
			db.InsertTable(PLAYER_TABLE_NAME, &self.Data, 0, true)
		} else {
			self.Data.Decode()
			self.Online = game.LOGIC_FALSE

			myunion := union.GetUnionMgr().GetUnion(self.Data.data.UnionId)
			if myunion != nil {
				member := myunion.GetMember(self.GetUid())
				if member != nil {
					if member.Lastlogintime == 0 {
						member.Lastlogintime = time.Now().Unix()
					}
				}
			}
		}
	}

	if self.Data.lifetree == nil {
		self.Data.lifetree = new(game.JS_LifeTreeInfo)
	}

	self.Data.Init(PLAYER_TABLE_NAME, &self.Data, true)

	if self.DataFriend == nil {
		self.DataFriend = new(ModFriend)
		self.DataFriend.onGetData(self)
	}
}

//! 保存数据
func (self *Player) onSave(force bool) {
	//! 强制保存，或者需要保存时
	if force || self.NeedSave {
		//! 基础数据保存
		self.Data.Encode()
		self.Data.Update(true, false)

		//! 好友数据保存
		self.DataFriend.onSave()
		self.NeedSave = false
	}
}

//! 保存数据
func (self *Player) Save() {
	self.NeedSave = true
}

func (self *Player) GetUId() int64 {
	return self.Data.UId
}

//! 同步线上数据
func (self *Player) UpdatePlayerData(req *RPC_PlayerData_Req) bool {
	self.Data.data = req.Data
	self.Data.heros = req.Heros
	self.Data.equips = req.Equips
	self.Data.lifetree = req.LifeTree

	self.Save()
	return true
}

func (self *Player) GetUid() int64 {
	return self.Data.UId
}

func (self *Player) GetUname() string {
	return self.Data.UName
}

func (self *Player) GetLevel() int {
	return self.Data.Level
}

func (self *Player) GetIconId() int {
	return self.Data.data.IconId
}

func (self *Player) GetPortrait() int {
	return self.Data.data.Portrait
}

func (self *Player) GetVip() int {
	return self.Data.Level
}

func (self *Player) GetFight() int {
	return self.Data.Fight
}

func (self *Player) GetServerId() int {
	return self.Data.ServerId
}

func (self *Player) GetLifeTree() *game.JS_LifeTreeInfo  {
	return self.Data.lifetree
}

//! 公会Id，暂时未加入
func (self *Player) GetUnionId() int {
	return self.Data.data.UnionId
}

func (self *Player) GetDataInt(attType int) int {
	switch attType {
	case core.PLAYER_ATT_LEVEL:
		return self.Data.data.Level
	case core.PLAYER_ATT_ICON:
		return self.Data.data.IconId
	}

	return 0
}

func (self *Player) GetDataInt64(attType int) int64 {
	return 0
}

func (self *Player) GetDataString(attType int) string {
	switch attType {
	case core.PLAYER_ATT_UNAME:
		return self.Data.data.UName
	}

	return ""
}

func (self *Player) OnClose() {

}

func (self *Player) GetSession() core.ISession {
	return nil
}
