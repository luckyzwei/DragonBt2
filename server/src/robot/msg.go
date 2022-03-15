package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"log"
	"math/rand"
)

func (self *Robot) SendLogin() {
	var msg C2S_Reg
	msg.Ctrl = "login_guest"
	msg.Uid = self.Uid
	msg.Account = self.Account
	msg.Password = self.Password
	msg.ServerId = self.ServerId

	self.Send("passport.php", &msg)
}

func (self *Robot) SendLoginOK() {
	var msg C2S_Uid
	msg.Ctrl = "loginok"
	msg.Uid = self.Uid

	self.Send("", &msg)
}

func (self *Robot) SendMission(missionid int, step int, warnum int, worknum int) {
	var msg C2S_SetMission
	msg.Ctrl = "setmission"
	msg.Uid = self.Uid
	msg.MissionId = missionid
	msg.Step = step
	msg.WarNum = warnum
	msg.WorkNum = worknum

	self.Send("", &msg)
}

func (self *Robot) SendZyId(zyid int) {
	var msg C2S_ZyInfo
	msg.Ctrl = "savezy"
	msg.Uid = self.Uid
	msg.Zyid = zyid

	self.Send("", &msg)
}

func (self *Robot) SendChapterevents(chapterId int) {
	var msg C2S_ChaptErevents
	msg.Ctrl = "chapterevents"
	msg.Uid = self.Uid
	msg.Chapter = chapterId
	self.Send("chapter.php", &msg)
}

func (self *Robot) SendGetCamp() {
	var msg C2S_Uid
	msg.Ctrl = "getcamp"
	msg.Uid = self.Uid

	self.Send("", &msg)
}

func (self *Robot) SendSetCamp(camp int) {
	var msg C2S_SetCamp
	msg.Ctrl = "setcamp"
	msg.Uid = self.Uid
	msg.Camp = camp

	self.Send("", &msg)
}

func (self *Robot) SendWorldEvent(ver int) {
	var msg C2S_WordEvent
	msg.Ctrl = "getworldevent"
	msg.Uid = self.Uid
	msg.Ver = 0

	self.Send("", &msg)
}

func (self *Robot) SendKingTop(camp int, topverking int) {
	var msg C2S_GetTopKing
	msg.Ctrl = "getkingtop"
	msg.Uid = self.Uid
	msg.Camp = camp
	msg.TopverKing = topverking

	self.Send("", &msg)
}

func (self *Robot) SendCreateRole() {
	var msg C2S_CreateRole
	msg.Ctrl = "createrole"
	msg.Uid = self.Uid
	msg.Name = fmt.Sprintf("%d", self.Uid+20190423)
	msg.Icon = 1002
	msg.Face = 1
	self.Send("find.php", &msg)
}

//! 进出大地图
func (self *Robot) EnterBigMap(_type int) {
	//log.Println("模拟进入大地图")
	var msg C2S_BigMap
	msg.Ctrl = "bigmap"
	msg.Uid = self.Uid
	msg.Type = _type
	self.Send("", &msg)
}

func (self *Robot) SetBossAction(id int,action int) {
	//log.Println("模拟进入大地图")
	var msg C2S_BossAction
	msg.Ctrl = "bossaction"
	msg.Uid = self.Uid
	msg.Id = id
	msg.Action = action
	self.Send("", &msg)
}

func (self *Robot) GetCityInfo(city int) {
	var msg C2S_GetCityInfo
	msg.Ctrl = "getcityinfo"
	msg.Uid = self.Uid
	msg.Cityid = city
	self.Send("", &msg)
}

func (self *Robot) GetGem() {
	var msg C2S_GetGem
	msg.Ctrl = "createitem"
	msg.Uid = self.Uid
	msg.ItemId = 91000002
	msg.Num = 1000000
	self.Send("", &msg)
}

func (self *Robot) GetMoney() {
	var msg C2S_GetGem
	msg.Ctrl = "createitem"
	msg.Uid = self.Uid
	msg.ItemId = 91000001
	msg.Num = 1000000
	self.Send("", &msg)
}

func (self *Robot) CreateItem(itemid int, itemnum int) {
	var msg C2S_GetGem
	msg.Ctrl = "createitem"
	msg.Uid = self.Uid
	msg.ItemId = itemid
	msg.Num = itemnum
	self.Send("", &msg)
}

func (self *Robot) MoveTeam(index int, way []int) {
	var msg C2S_MoveTeam
	msg.Ctrl = "moveteam"
	msg.Uid = self.Uid
	msg.Index = index
	msg.Way = way
	self.Send("", &msg)
}

func (self *Robot) MoveTeamBegin(begin int, city int, index int) {
	var msg C2S_MoveTeamBegin
	msg.Ctrl = "moveteambegin"
	msg.Uid = self.Uid
	msg.Begin = begin
	msg.Cityid = city
	msg.Index = index
	self.Send("", &msg)
}

func (self *Robot) GetPassBegin(passid int) {
	var msg C2S_BeginPass
	msg.Ctrl = "passbegin"
	msg.Uid = self.Uid
	msg.Passid = passid
	self.Send("", &msg)
}

func (self *Robot) SendTakeNobilitytask(taskid int) {
	var msg C2S_TakeNobilitytask
	msg.Ctrl = "takenobilitytask"
	msg.Uid = self.Uid
	msg.TaskId = taskid
	self.Send("", &msg)
}


func (self *Robot) GetPassWin(missionid int) {
	var msg C2S_WinPass
	msg.Ctrl = "passwin"
	msg.Uid = self.Uid
	msg.Missionid = missionid
	self.Send("", &msg)
}

func (self *Robot) SendPassEnd(passid int, missionid int, step int, warnum int, worknum int) {
	var msg C2S_EndPass
	msg.Ctrl = "passresult"
	msg.Uid = self.Uid
	msg.Star = 3
	msg.Passid = passid
	msg.Index = 0
	msg.FightTime = 0
	msg.MissionId = missionid
	msg.Step = step
	msg.WarNum = warnum
	msg.WorkNum = worknum
	self.Send("", &msg)
}

func (self *Robot) SendBoxPass(passid int, nopass int, missionid int, step int, warnum int, worknum int) {
	var msg C2S_BoxPass
	msg.Ctrl = "boxpass"
	msg.Uid = self.Uid
	msg.Passid = passid
	msg.NoPass = nopass
	msg.MissionId = missionid
	msg.Step = step
	msg.WarNum = warnum
	msg.WorkNum = worknum
	self.Send("", &msg)
}

func (self *Robot) SendUpColor(heroid int) {
	var msg C2S_UpColor
	msg.Ctrl = "upcolor"
	msg.Uid = self.Uid
	msg.Heroid = heroid
	self.Send("", &msg)
}

func (self *Robot) SendPassWin(missionid int, step int, warnum int, worknum int) {
	var msg C2S_SetMission
	msg.Ctrl = "passwin"
	msg.Uid = self.Uid
	msg.MissionId = missionid
	msg.Step = step
	msg.WarNum = warnum
	msg.WorkNum = worknum
	self.Send("", &msg)
}

func (self *Robot) GetNewEvent() {
	var msg C2S_Uid
	msg.Ctrl = "getnewevent"
	msg.Uid = self.Uid
	self.Send("", &msg)
}

func (self *Robot) GetMailAllItem() {
	var msg C2S_Uid
	msg.Ctrl = "getmailallitem"
	msg.Uid = self.Uid
	self.Send("", &msg)
}

func (self *Robot) SendJJPass(passid int) {
	var msg C2S_JJPass
	msg.Ctrl = "jjpass"
	msg.Uid = self.Uid
	msg.Passid = passid
	self.Send("", &msg)
}

func (self *Robot) SendBoxSp(index int) {
	var msg C2S_LevelUpBoxSP
	msg.Ctrl = "levelupboxsp"
	msg.Uid = self.Uid
	msg.Index = index
	self.Send("", &msg)
}

func (self *Robot) SendCheckin() {
	var msg C2S_Uid
	msg.Ctrl = "checkin"
	msg.Uid = self.Uid
	self.Send("", &msg)
}

func (self *Robot) SendCheckinAward(index int) {
	var msg C2S_CheckinAward
	msg.Ctrl = "checkinaward"
	msg.Uid = self.Uid
	msg.Index = index
	self.Send("", &msg)
}

func (self *Robot) SendFinishActivity(id int) {
	var msg C2S_FinishActivity
	msg.Ctrl = "finishactivity"
	msg.Uid = self.Uid
	msg.Id = id
	self.Send("", &msg)
}

func (self *Robot) SendSynthesis(heroid string) {
	var msg C2S_Synthesis
	msg.Ctrl = "synthesis"
	msg.Uid = self.Uid
	msg.Heroid = heroid
	self.Send("", &msg)
}

func (self *Robot) SendCampTeam(team [3]JS_CampTeam) {
	var msg C2S_CampTeam
	msg.Ctrl = "campteam"
	msg.Uid = self.Uid
	msg.Team = team
	self.Send("", &msg)
}

func (self *Robot) SendChat(channel int, content string, name string) {
	data := &Chat{
		Channel:   proto.Int32(int32(channel)),
		Content:   []byte(content),
		Name:      proto.String(name),
		Medianame: proto.String(""),
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	websocket.Message.Send(self.Ws, HF_DecodeMsg("chat", msg))
}

func (self *Robot) SendDanMu(setrDec string) {

	var msg C2S_Barrage
	msg.Ctrl = "barrage"
	msg.Uid = self.Uid
	msg.Cityid = 0
	msg.Size = 20
	msg.Text = setrDec
	msg.Red = 255
	msg.Green = 255
	msg.Blue = 255

	self.Send("", &msg)
}

//! 抽奖
func (self *Robot) Find(index int, free int) {
	var msg C2S_Find
	msg.Ctrl = "draw"
	msg.Uid = self.Uid
	msg.Findtype = index
	msg.Free = free
	self.Send("", &msg)
}

//!获取好友推荐
func (self *Robot) GetFriendCommend() {
	var msg C2S_FriendCommend
	msg.Ctrl = "friendcommend"
	msg.Uid = self.Uid
	msg.Refresh = 0
	self.Send("", &msg)
}

func (self *Robot) GetStatistcsinfo() {
	var msg C2S_StatisticsInfo
	msg.Ctrl = "statisticsinfo"
	msg.Uid = self.Uid
	self.Send("", &msg)
}

//!一键申请好友
func (self *Robot) SendFrienDapply() {
	var msg C2S_FrienDapply
	msg.Ctrl = "friendapply"
	msg.Uid = self.Uid
	msg.Pid = 0
	self.Send("", &msg)
}

//!一键通过好友
func (self *Robot) SendFrienDorder() {
	log.Println("一键通过好友")
	var msg C2S_FrienDorder
	msg.Ctrl = "friendorder"
	msg.Uid = self.Uid
	msg.Pid = 0
	msg.Agree = 1
	self.Send("", &msg)
}

//!一键通过章节
func (self *Robot) SendPassId(zhangjieid int) {
	var msg C2S_PassMission
	msg.Ctrl = "gm_pass_mission"
	msg.Uid = self.Uid
	msg.ChpaterId = zhangjieid
	self.Send("", &msg)
}

//!武将神器神级
//武将id ,升级类型（0 升级一次，1一键升级），装备ID
func (self *Robot) SendSenqiLvUp(heroid int, sjtype int, zhuangbeiid int) {
	var msg C2S_ShenqilvUp
	msg.Ctrl = "artifactuplevel"
	msg.Uid = self.Uid
	msg.Heroid = heroid
	msg.Type = sjtype
	msg.Id = zhuangbeiid
	self.Send("", &msg)
}

//!武将升星
func (self *Robot) SendHeroStarUp(heroid int) {
	var msg C2S_UpHeroStar
	msg.Ctrl = "heroupstar"
	msg.Uid = self.Uid
	msg.Heroid = heroid
	self.Send("", &msg)
}

//! 兵种升级
func (self *Robot) SendSoldierLvUp(heroid int) {
	var msg C2S_SoldierType
	msg.Ctrl = "armsuplevel"
	msg.Uid = self.Uid
	msg.Heroid = heroid
	msg.Type = 0
	self.Send("", &msg)
}

//!参加国家
func (self *Robot) SendCampFightMove() {

	var msg C2S_CampFightMove
	msg.Ctrl = "campfightmove"
	msg.Uid = self.Uid
	msg.Cityid = GetRobotMgr().RobotCon.GuoZhanCityId
	msg.Type = 0 //进国战
	self.Send("", &msg)
}

func (self *Robot) SendCampFightMoveEx(cityId int) {
	var msg C2S_CampFightMove
	msg.Ctrl = "campfightmove"
	msg.Uid = self.Uid
	msg.Cityid = cityId
	msg.Type = 0 //进国战
	self.Send("", &msg)
}

//!请求单挑队列
func (self *Robot) SendCampfightSoloReq(cityId int, cityPart int) {
	self.is_fightinfo = false
	var msg C2S_CampFightSoloReq
	msg.Ctrl = "campfightsoloreq"
	msg.Uid = self.Uid
	msg.CityId = cityId
	msg.WarPlayId = 1
	msg.Again = 0

	if cityPart == 99 {
		cityPart = rand.Intn(5)
	}

	msg.CityPart = cityPart
	self.Send("", &msg)
}

//!获取单挑对手
func (self *Robot) GetCampfightPlayList(cityId int, cityPart int) {
	self.is_fightinfo = false
	var msg C2S_CampFightSoloReq
	msg.Ctrl = "campfightplaylist"
	msg.Uid = self.Uid

	msg.CityId = cityId
	if cityId == 0 {
		msg.CityId = self.CampFightCity
	}

	msg.WarPlayId = 1
	msg.Again = 0

	if cityPart == 99 {
		cityPart = rand.Intn(5)
	}

	msg.CityPart = cityPart
	self.Send("", &msg)
}

//!选择单挑对手开战
func (self *Robot) SendCampFightSolo2(emindex int, uid int64, cityid int) {

	var msg C2S_CampFightSolo2
	msg.Ctrl = "campfightsolo2"
	msg.Uid = self.Uid

	msg.Cityid = GetRobotMgr().RobotCon.GuoZhanCityId
	if cityid != 0 {
		msg.Cityid = cityid
	}

	msg.MyIndex = 0
	msg.Pid = uid
	msg.Index = emindex
	self.Send("", &msg)
}

//!发送单挑结果 result 1胜利 2失败
func (self *Robot) SendCampFightSoloEnd(result int, cityid int) {
	var msg C2S_CampFightSoloEnd
	msg.Ctrl = "campfightsoloend"
	msg.Uid = self.Uid

	msg.Cityid = GetRobotMgr().RobotCon.GuoZhanCityId
	if cityid != 0 {
		msg.Cityid = cityid
	}

	msg.Soloid = self.SoloId
	msg.Result = result
	self.Send("", &msg)
}

//!报名国战决战
func (self *Robot) SendCampFight56Req(warplayid int, cityid int) {

	var msg C2S_CampFight55Req
	msg.Ctrl = "campfight56req"
	msg.Uid = self.Uid

	msg.Cityid = self.CampFightCity
	if cityid != 0 {
		msg.Cityid = cityid
	}

	msg.WarPlayid = warplayid
	self.Send("", &msg)
}

//!武将升级
func (self *Robot) SendHeroLvUp(heroid int, itemid string, num int) {
	var msg C2S_UseItem
	msg.Ctrl = "useitem"
	msg.Uid = self.Uid

	msg.Itemid = itemid
	msg.Num = num
	msg.Destid = heroid

	self.Send("", &msg)
}

//!请求军令数据
func (self *Robot) SendTeamPowerReq() {

	var msg C2S_TimePower
	msg.Index = 0
	msg.Ctrl = "teampower"
	msg.Uid = self.Uid

	self.Send("teampower", &msg)
}

func (self *Robot) SendGuildId(guild int) {
	//150以下是强制引导
	if guild<150{
		for ; self.guildId <= guild; self.guildId++ {
			var msg C2S_SetGuild
			msg.Ctrl = "set_guild_id"
			msg.GuildId = self.guildId
			msg.Uid = self.Uid
			self.Send("passport.php", &msg)
		}
	}else{
		var msg C2S_SetGuild
		msg.Ctrl = "set_guild_id"
		msg.GuildId = self.guildId
		msg.Uid = self.Uid
		self.Send("passport.php", &msg)
	}
}

func (self *Robot) SendGetMission(mission int) {
	var msg C2S_GetMission
	msg.Ctrl = "getmission"
	msg.MissionId = mission
	msg.Uid = self.Uid
	self.Send("", &msg)
}

func (self *Robot) SendSwapFightPos(index1 int, index2 int, teamtype int) {
	var msg C2S_SwapFightPos
	msg.Ctrl = "swapfightpos"
	msg.Index1 = index1
	msg.Index2 = index2
	msg.TeamType = teamtype
	msg.Uid = self.Uid

	self.Send("", &msg)
}

func (self *Robot) SendFinishTask(taskid string, tasktype int) {
	var msg C2S_TaskFinish
	msg.Ctrl = "finishtask"
	msg.Taskid = taskid
	msg.Tasktype = tasktype
	msg.Uid = self.Uid
	self.Send("task2.php", &msg)
}

func (self *Robot) SendAddGuide(guide_id int) {
	var msg C2S_AddGuide
	msg.Ctrl = "add_guide"
	msg.GuideId = guide_id
	msg.Uid = self.Uid
	self.Send("passport.php", &msg)
}

func (self *Robot) SendUpStarAuto(heroid int) {
	var msg C2S_UpStarAuto
	msg.Ctrl = "upstarauto"
	msg.HeroId = heroid
	msg.Uid = self.Uid

	self.Send("", &msg)
}

func (self *Robot) SendActivateStar(index int, heroid int) {
	var msg C2S_ActivateStar
	msg.Ctrl = "activatestar"
	msg.Index = index
	msg.HeroId = heroid
	msg.Uid = self.Uid

	self.Send("", &msg)
}

func (self *Robot) SendAddteampos(index int, heroid int, teamtype int) {
	var msg C2S_AddTeamUIPos
	msg.Ctrl = "addteampos"
	msg.Index = index
	msg.HeroId = heroid
	msg.TeamType = teamtype
	msg.Uid = self.Uid

	self.Send("", &msg)
}

func (self *Robot) SendSetRedicon(redicon int) {
	var msg C2S_SetRedIcon
	msg.Ctrl = "set_redicon"
	msg.Id = redicon

	self.Send("", &msg)
}

func (self *Robot) SendEquipAction(teamtype int, heroid int, keyid int, compundnum int, pos int, index int, action int, itemid int) {
	var msg C2S_EquipAction
	msg.Ctrl = "equipaction"
	msg.Uid = self.Uid
	msg.TeamType = teamtype
	msg.HeroId = heroid
	msg.KeyId = keyid
	msg.CompundNum = compundnum
	msg.Pos = pos
	msg.Index = index
	msg.Action = action
	msg.Itemid = itemid

	self.Send("", &msg)
}

func (self *Robot) SendFinishEvents(thindId int, eventId int) {
	var msg C2S_FinishEvents
	msg.Ctrl = "finishevents"
	msg.ThingId = thindId
	msg.EventId = eventId
	msg.Uid = self.Uid
	self.Send("chapter.php", &msg)
}

//!购买军令
func (self *Robot) SendBuyTeamPowerReq() {
	var msg C2S_Collection
	msg.Index = 3
	msg.Ctrl = "collection"
	msg.Uid = self.Uid

	self.Send("collection", &msg)
}
