package game

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// 数据库字段检查
var mysqlMgr *MySqlMgr = nil

type MySqlMgr struct {
	tableCheck            []string
	logTableCheck         []string
	fieldCheck            [][2]string
	createTableStatements []string
	addField              [][3]string
	modifyField           [][4]string
	dropField             [][3]string
}

func GetSqlMgr() *MySqlMgr {
	if mysqlMgr == nil {
		mysqlMgr = &MySqlMgr{}
	}
	return mysqlMgr
}

func (self *MySqlMgr) initData() {
	self.tableCheck = []string{
		"san_act_rankaward",
		"san_cost",
		"san_invest",
		"san_investmsg",
		"san_smelt",
		"san_tiger",
	}

	self.fieldCheck = [][2]string{
		{"san_userbase", "blocktime"},
		{"san_userbase", "blockreason"},
		{"san_tiger", "buytimes"},
	}

	self.logTableCheck = []string{
		"san_mail",
	}

	self.addField = [][3]string{
		{"san_userluckegg", "total", "ALTER TABLE `san_userluckegg` ADD COLUMN `total`  int(11) NOT NULL DEFAULT 0 COMMENT '金额' AFTER `info`;"},
		{"san_fightresult", "cityid", "ALTER TABLE `san_fightresult` ADD COLUMN `cityid`  int(11) NOT NULL DEFAULT 0 COMMENT '城池' AFTER `id`;"},
		{"san_city", "gvgfight", "ALTER TABLE `san_city` ADD COLUMN `gvgfight`  text NOT NULL COMMENT 'gvg' AFTER `build`;"},
		{"san_expedition1", "process", "ALTER TABLE `san_expedition1` ADD COLUMN `process`  int(11) NOT NULL DEFAULT 0 COMMENT '总进度' AFTER `flowinfo`;"},
		{"san_expedition1", "lastprocess", "ALTER TABLE `san_expedition1` ADD COLUMN `lastprocess`  int(11) NOT NULL DEFAULT 0 COMMENT '上次进度' AFTER `process`;"},
		{"san_userbase", "soul", "ALTER TABLE `san_userbase` ADD COLUMN `soul`  int(11) NOT NULL DEFAULT 0 COMMENT '魂石' AFTER `blockreason`;"},
		{"san_userbase", "techpoint", "ALTER TABLE `san_userbase` ADD COLUMN `techpoint`  int(11) NOT NULL DEFAULT 0 COMMENT '科技点' AFTER `soul`;"},
		{"san_userbase", "bossmoney", "ALTER TABLE `san_userbase` ADD COLUMN `bossmoney`  int(11) NOT NULL DEFAULT 0 COMMENT '巨兽精魄' AFTER `techpoint`;"},
		{"san_armsarena1", "timenum", "ALTER TABLE `san_armsarena1` ADD COLUMN `timenum`  int(11) NOT NULL DEFAULT 0 COMMENT '时间重置次数' AFTER `topaward`;"},
		{"san_gemstone", "curchapter", "ALTER TABLE `san_gemstone` ADD COLUMN `curchapter`  int(11) NOT NULL DEFAULT 0 COMMENT '当前章节' AFTER `uid`;"},
		{"san_userhorse", "horsetask", "ALTER TABLE `san_userhorse` ADD COLUMN `horsetask`  text NOT NULL COMMENT '高级召唤战马任务' AFTER `info4`;"},
		{"san_userbase", "towerstone", "ALTER TABLE `san_userbase` ADD COLUMN `towerstone`  int(11) NOT NULL DEFAULT 0 COMMENT '镇魂石' AFTER `bossmoney`;"},
		{"san_userbase", "portrait", "ALTER TABLE `san_userbase` ADD COLUMN `portrait`  int(11) NOT NULL DEFAULT 0 COMMENT '头像挂件' AFTER `towerstone`;"},
		{"san_userbase", "campok", "ALTER TABLE `san_userbase` ADD COLUMN `campok`  int(11) NOT NULL DEFAULT 0 COMMENT '阵营ok' AFTER `portrait`;"},
		{"san_userbase", "nameok", "ALTER TABLE `san_userbase` ADD COLUMN `nameok`  int(11) NOT NULL DEFAULT 0 COMMENT '头像ok' AFTER `campok`;"},
		{"san_smelt", "smelttask", "ALTER TABLE `san_smelt` ADD COLUMN `smelttask`  text NOT NULL COMMENT '冶炼次数任务' AFTER `tigerchipnum`;"},
		{"san_smelt", "buytask", "ALTER TABLE `san_smelt` ADD COLUMN `buytask`  text NOT NULL COMMENT '购买商品次数' AFTER `smelttask`;"},
		{"san_smelt", "smelttasktime", "ALTER TABLE `san_smelt` ADD COLUMN `smelttasktime` bigint(20) NOT NULL DEFAULT 0 COMMENT '冶炼任务刷新时间' AFTER `buytask`;"},
		{"san_smelt", "buytasktime", "ALTER TABLE `san_smelt` ADD COLUMN `buytasktime`  bigint(20) NOT NULL DEFAULT 0 COMMENT '购买冶炼任务刷新时间' AFTER `smelttasktime`;"},
		{"san_smelt", "smelttimes", "ALTER TABLE `san_smelt` ADD COLUMN `smelttimes`  int(11) NOT NULL DEFAULT 0 COMMENT '冶炼次数' AFTER `buytasktime`;"},
		{"san_smelt", "buytimes", "ALTER TABLE `san_smelt` ADD COLUMN `buytimes`  int(11) NOT NULL DEFAULT 0 COMMENT '购买次数' AFTER `smelttimes`;"},
		{"san_account", "userid", "ALTER TABLE `san_account` ADD COLUMN `userid`  text NOT NULL COMMENT '用户id' AFTER `account`;"},
		{"san_military", "weektask", "ALTER TABLE `san_military` ADD COLUMN `weektask`  text NOT NULL COMMENT '周任务' AFTER `taskonce`;"},
		{"san_military", "weektime", "ALTER TABLE `san_military` ADD COLUMN `weektime`  int(11) NOT NULL DEFAULT 0 COMMENT '周刷新时间' AFTER `weektask`;"},
		{"san_weekplan", "boxnum", "ALTER TABLE `san_weekplan` ADD COLUMN `boxnum`  int(11) NOT NULL DEFAULT 0 COMMENT '宝箱数量' AFTER `taskinfo`;"},
		{"san_weekplan", "point", "ALTER TABLE `san_weekplan` ADD COLUMN `point`  int(11) NOT NULL DEFAULT 0 COMMENT '点数' AFTER `boxnum`;"},
		{"san_weekplan", "isget", "ALTER TABLE `san_weekplan` ADD COLUMN `isget`  int(11) NOT NULL DEFAULT 0 COMMENT '是否领取' AFTER `point`;"},
		{"san_weekplan", "isgetmark", "ALTER TABLE `san_weekplan` ADD COLUMN `isgetmark`   text NOT NULL COMMENT '七日状态' AFTER `isget`;"},
		{"san_weekplan", "stage", "ALTER TABLE `san_weekplan` ADD COLUMN `stage`  int(11) NOT NULL DEFAULT 0 COMMENT '活动阶段' AFTER `isgetmark`;"},
		{"san_halfmoon", "boxnum", "ALTER TABLE `san_halfmoon` ADD COLUMN `boxnum`  int(11) NOT NULL DEFAULT 0 COMMENT '宝箱数量' AFTER `taskinfo`;"},
		{"san_tower", "towerbox", "ALTER TABLE `san_tower` ADD COLUMN `towerbox`   text NOT NULL COMMENT '爬塔宝箱状态' AFTER `advancebuytimes`;"},
		{"san_userkingtask", "box", "ALTER TABLE `san_userkingtask` ADD COLUMN `box` text NOT NULL COMMENT '宝箱状态' AFTER `taskgroup`;"},
		{"san_userfind", "goldendtime", "ALTER TABLE `san_userfind` ADD COLUMN `goldendtime` bigint(20) NOT NULL DEFAULT 0 COMMENT '结束时间' AFTER `goldfreetimes`;"},
		{"san_userfind", "lootenergy", "ALTER TABLE `san_userfind` ADD COLUMN `lootenergy` int(11) NOT NULL DEFAULT 0 COMMENT '抽奖能量' AFTER `goldendtime`;"},
		{"san_userfind", "summontimes", "ALTER TABLE `san_userfind` ADD COLUMN `summontimes` int(11) NOT NULL DEFAULT 0 COMMENT '召唤打折次数' AFTER `lootenergy`;"},
		{"san_userhero2", "startime", "ALTER TABLE `san_userhero2` ADD COLUMN `startime`  bigint(20) NOT NULL DEFAULT 0 COMMENT '达到总天赋时间' AFTER `totalstars`;"},
		{"san_userhero2", "herototalstars", "ALTER TABLE `san_userhero2` ADD COLUMN `herototalstars`  int(11) NOT NULL DEFAULT 0 COMMENT '所有英雄总星级' AFTER `reborn`;"},
		{"san_userhero2", "herostartime", "ALTER TABLE `san_userhero2` ADD COLUMN `herostartime`  bigint(20) NOT NULL DEFAULT 0 COMMENT '达到总星级的时间' AFTER `herototalstars`;"},
		{"san_userhero2", "maxkey", "ALTER TABLE `san_userhero2` ADD COLUMN `maxkey`  bigint(20) NOT NULL DEFAULT 0 COMMENT '英雄最大KEY' AFTER `herostartime`;"},
		{"san_userhero2", "buyposnum", "ALTER TABLE `san_userhero2` ADD COLUMN `buyposnum`  bigint(20) NOT NULL DEFAULT 0 COMMENT '英雄栏扩展次数' AFTER `maxkey`;"},
		{"san_userhero2", "autofire", "ALTER TABLE `san_userhero2` ADD COLUMN `autofire`  bigint(20) NOT NULL DEFAULT 0 COMMENT '自动分解' AFTER `buyposnum`;"},
		{"san_userhero2", "backopen", "ALTER TABLE `san_userhero2` ADD COLUMN `backopen`  bigint(20) NOT NULL DEFAULT 0 COMMENT '自动分解' AFTER `autofire`;"},
		{"san_userhero2", "handbook", "ALTER TABLE `san_userhero2` ADD COLUMN `handbook` text NOT NULL COMMENT '图鉴记录' AFTER `backopen`;"},
		{"san_userhero2", "compoundsignnum", "ALTER TABLE `san_userhero2` ADD COLUMN `compoundsignnum` text NOT NULL COMMENT '合成标记' AFTER `handbook`;"},
		{"san_userhero2", "compoundsignscore", "ALTER TABLE `san_userhero2` ADD COLUMN `compoundsignscore` text NOT NULL COMMENT '合成标记' AFTER `compoundsignnum`;"},
		{"san_pass", "starttime", "ALTER TABLE `san_pass` ADD COLUMN `starttime`  bigint(20) NOT NULL DEFAULT 0 COMMENT '达到关卡星级的时间' AFTER `totalstars`;"},
		{"san_userhorse", "horsetotalfight", "ALTER TABLE `san_userhorse` ADD COLUMN `horsetotalfight`  bigint(20) NOT NULL DEFAULT 0 COMMENT '坐骑总战力' AFTER `horsetask`;"},
		{"san_tiger", "tigertotalfight", "ALTER TABLE `san_tiger` ADD COLUMN `tigertotalfight`  bigint(20) NOT NULL DEFAULT 0 COMMENT '纹章总战力' AFTER `gift`;"},
		{"san_userequip", "totalgemlevel", "ALTER TABLE `san_userequip` ADD COLUMN `totalgemlevel`  int(11) NOT NULL DEFAULT 0 COMMENT '总装备宝石等级' AFTER `info`;"},
		{"san_userequip", "starttime", "ALTER TABLE `san_userequip` ADD COLUMN `starttime`  bigint(20) NOT NULL DEFAULT 0 COMMENT '达到装备宝石的时间' AFTER `totalgemlevel`;"},
		{"san_armsarena2", "timenum", "ALTER TABLE `san_armsarena2` ADD COLUMN `timenum`  int(11) NOT NULL DEFAULT 0 COMMENT '时间重置次数' AFTER `topaward`;"},
		{"san_armsarena3", "timenum", "ALTER TABLE `san_armsarena3` ADD COLUMN `timenum`  int(11) NOT NULL DEFAULT 0 COMMENT '时间重置次数' AFTER `topaward`;"},
		{"san_userboss", "cost", "ALTER TABLE `san_userboss` ADD COLUMN `cost`  int(11) NOT NULL DEFAULT 0 COMMENT '消耗' AFTER `bossinfo`;"},
		{"san_mineglobal", "minevictorinfo", "ALTER TABLE `san_mineglobal` ADD COLUMN `minevictorinfo` text NOT NULL AFTER `minerelationinfo`;"},
		{"san_mineglobal", "minevictortime", "ALTER TABLE `san_mineglobal` ADD COLUMN `minevictortime` bigint(20) NOT NULL DEFAULT 0 COMMENT '展示时间' AFTER `minevictorinfo`;"},
		{"san_gveplayer", "needclear", "ALTER TABLE `san_gveplayer` ADD COLUMN `needclear` int(11) NOT NULL DEFAULT 0 COMMENT '是否清除鼓舞buff' AFTER `leveldata`;"},
		{"san_mineplayer", "needclear", "ALTER TABLE `san_mineplayer` ADD COLUMN `needclear` int(11) NOT NULL DEFAULT 0 COMMENT '是否清除鼓舞buff' AFTER `fightrecordinfo`;"},
		{"san_mineplayer", "movetime", "ALTER TABLE `san_mineplayer` ADD COLUMN `movetime` int(11) NOT NULL DEFAULT 0 COMMENT '移动动画时间' AFTER `needclear`;"},
		{"san_tiger", "gift", "ALTER TABLE `san_tiger` ADD COLUMN `gift` int(11) NOT NULL DEFAULT 0 COMMENT '初始道具状态' AFTER `buytimes`;"},
		{"san_armsarena1", "maxkey", "ALTER TABLE `san_armsarena1` ADD COLUMN `maxkey`  int(11) NOT NULL DEFAULT 0 COMMENT '个人最大战报Id' AFTER `timenum`;"},
		{"san_armsarena2", "maxkey", "ALTER TABLE `san_armsarena2` ADD COLUMN `maxkey`  int(11) NOT NULL DEFAULT 0 COMMENT '个人最大战报Id' AFTER `timenum`;"},
		{"san_armsarena3", "maxkey", "ALTER TABLE `san_armsarena3` ADD COLUMN `maxkey`  int(11) NOT NULL DEFAULT 0 COMMENT '个人最大战报Id' AFTER `timenum`;"},
		{"san_armsarena1", "pvpfight", "ALTER TABLE `san_armsarena1` ADD COLUMN `pvpfight`  text NOT NULL COMMENT '个人最大战报信息' AFTER `maxkey`;"},
		{"san_armsarena2", "pvpfight", "ALTER TABLE `san_armsarena2` ADD COLUMN `pvpfight`  text NOT NULL COMMENT '个人最大战报信息' AFTER `maxkey`;"},
		{"san_armsarena3", "pvpfight", "ALTER TABLE `san_armsarena3` ADD COLUMN `pvpfight`  text NOT NULL COMMENT '个人最大战报信息' AFTER `maxkey`;"},
		{"san_armsarena1", "randnum", "ALTER TABLE `san_armsarena1` ADD COLUMN `randnum`  bigint(20) NOT NULL DEFAULT 0 COMMENT '战斗随机值' AFTER `pvpfight`;"},
		{"san_armsarena2", "randnum", "ALTER TABLE `san_armsarena2` ADD COLUMN `randnum`  bigint(20) NOT NULL DEFAULT 0 COMMENT '战斗随机值' AFTER `pvpfight`;"},
		{"san_armsarena3", "randnum", "ALTER TABLE `san_armsarena3` ADD COLUMN `randnum`  bigint(20) NOT NULL DEFAULT 0 COMMENT '战斗随机值' AFTER `pvpfight`;"},
		{"san_userbase", "guildid", "ALTER TABLE `san_userbase` ADD COLUMN `guildid`  int(11) NOT NULL DEFAULT 0 COMMENT '指引Id' AFTER `nameok`;"},
		{"san_userevents", "info21", "ALTER TABLE `san_userevents` ADD COLUMN `info21`   text NOT NULL COMMENT '章节信息' AFTER `info20`;"},
		{"san_userevents", "info22", "ALTER TABLE `san_userevents` ADD COLUMN `info22`   text NOT NULL COMMENT '章节信息' AFTER `info21`;"},
		{"san_userevents", "info23", "ALTER TABLE `san_userevents` ADD COLUMN `info23`   text NOT NULL COMMENT '章节信息' AFTER `info22`;"},
		{"san_userevents", "info24", "ALTER TABLE `san_userevents` ADD COLUMN `info24`   text NOT NULL COMMENT '章节信息' AFTER `info23`;"},
		{"san_userevents", "info25", "ALTER TABLE `san_userevents` ADD COLUMN `info25`   text NOT NULL COMMENT '章节信息' AFTER `info24`;"},
		{"san_userevents", "info26", "ALTER TABLE `san_userevents` ADD COLUMN `info26`   text NOT NULL COMMENT '章节信息' AFTER `info25`;"},
		{"san_userevents", "info27", "ALTER TABLE `san_userevents` ADD COLUMN `info27`   text NOT NULL COMMENT '章节信息' AFTER `info26`;"},
		{"san_userevents", "info28", "ALTER TABLE `san_userevents` ADD COLUMN `info28`   text NOT NULL COMMENT '章节信息' AFTER `info27`;"},
		{"san_userevents", "info29", "ALTER TABLE `san_userevents` ADD COLUMN `info29`   text NOT NULL COMMENT '章节信息' AFTER `info28`;"},
		{"san_userevents", "info30", "ALTER TABLE `san_userevents` ADD COLUMN `info30`   text NOT NULL COMMENT '章节信息' AFTER `info29`;"},
		{"san_userarmy", "teamid", "ALTER TABLE `san_userarmy` ADD COLUMN `teamid`  bigint(20) NOT NULL DEFAULT 0 COMMENT '队伍Id' AFTER `buytimes`;"},
		{"san_userarmy", "playts", "ALTER TABLE `san_userarmy` ADD COLUMN `playts`  text NOT NULL COMMENT '时间' AFTER `teamid`;"},
		{"san_campdiplomacy", "union", "ALTER TABLE `san_campdiplomacy` ADD COLUMN `union`   text NOT NULL COMMENT '军机处' AFTER `info`;"},
		{"san_campdiplomacy", "campfight", "ALTER TABLE `san_campdiplomacy` ADD COLUMN `campfight`   text NOT NULL COMMENT '国战宣战' AFTER `union`;"},
		{"san_king_cityfight_award", "getts", "ALTER TABLE `san_king_cityfight_award` ADD COLUMN `getts`   int(11) NOT NULL COMMENT '领取时间' AFTER `lastts`;"},
		{"san_userunioninfo", "stateaward", "ALTER TABLE `san_userunioninfo` ADD COLUMN `stateaward`  int(11) NOT NULL COMMENT '州宝箱状态' AFTER `copyaward`;"},
		{"san_userarmy", "playts", "ALTER TABLE `san_userarmy` ADD COLUMN `playts`  text NOT NULL COMMENT '时间' AFTER `teamid`;"},
		{"san_userbase", "redicon", "ALTER TABLE `san_userbase` ADD COLUMN `redicon`  int(11) NOT NULL DEFAULT 0 COMMENT '红色图标' AFTER `guildid`;"},
		{"san_gveglobal", "state", "ALTER TABLE `san_gveglobal` ADD COLUMN `state`  int(11) NOT NULL DEFAULT 0 COMMENT '状态信息' AFTER `draginfo`;"},
		{"san_king_cityfight", "power", "ALTER TABLE `san_king_cityfight` ADD COLUMN `power`  int(11) NOT NULL COMMENT '行动力' AFTER `title`;"},
		{"san_king_cityfight", "buypower", "ALTER TABLE `san_king_cityfight` ADD COLUMN `buypower`  int(11) NOT NULL COMMENT '购买行动力' AFTER `power`;"},
		{"san_userhero2", "reborn", "ALTER TABLE `san_userhero2` ADD COLUMN `reborn`  int(11) NOT NULL COMMENT '今日重生次数' AFTER `startime`;"},
		{"san_king1", "face", "ALTER TABLE `san_king1` ADD COLUMN `face`  int(11) NOT NULL COMMENT '性别' AFTER `enemy`;"},
		{"san_king2", "face", "ALTER TABLE `san_king2` ADD COLUMN `face`  int(11) NOT NULL COMMENT '性别' AFTER `enemy`;"},
		{"san_king3", "face", "ALTER TABLE `san_king3` ADD COLUMN `face`  int(11) NOT NULL COMMENT '性别' AFTER `enemy`;"},
		{"san_activity", "becomestronger", "ALTER TABLE `san_activity` ADD COLUMN `becomestronger` text NOT NULL COMMENT '我要变强领取状态' AFTER `month`;"},
		{"san_userrecharge", "basecounts", "ALTER TABLE `san_userrecharge` ADD COLUMN `basecounts` text NOT NULL COMMENT '基础充值次数' AFTER `fundget`;"},
		{"san_userrecharge", "boxcounts", "ALTER TABLE `san_userrecharge` ADD COLUMN `boxcounts` text NOT NULL COMMENT '礼包充值次数' AFTER `basecounts`;"},
		{"san_activity", "fund", "ALTER TABLE `san_activity` ADD COLUMN `fund` text NOT NULL COMMENT '活动基金' AFTER `becomestronger`;"},
		{"san_dreamland", "guidetype", "ALTER TABLE `san_dreamland` ADD COLUMN `guidetype` int(11) NOT NULL DEFAULT 0 COMMENT '引导状态' AFTER `info`;"},
		{"san_userrecharge", "vipdailyreward", "ALTER TABLE `san_userrecharge` ADD COLUMN `vipdailyreward` bigint(20) NOT NULL COMMENT 'VIP每日奖励' AFTER `boxcounts`;"},
		{"san_userrecharge", "vipweekbuy", "ALTER TABLE `san_userrecharge` ADD COLUMN `vipweekbuy` bigint(20) NOT NULL COMMENT 'VIP每周购买' AFTER `vipdailyreward`;"},
		{"san_consumertopuser", "portrait", "ALTER TABLE `san_consumertopuser` ADD COLUMN `portrait` int(11) NOT NULL DEFAULT 0  COMMENT '头像框' AFTER `icon`;"},
		{"san_userpit", "maxkeyid", "ALTER TABLE `san_userpit` ADD COLUMN `maxkeyid`  int(11) NOT NULL DEFAULT 0 COMMENT '最大副本ID' AFTER `pitinfo`;"},
		{"san_userpit", "firstinfo", "ALTER TABLE `san_userpit` ADD COLUMN `firstinfo` text NOT NULL COMMENT '首次通关标记' AFTER `maxkeyid`;"},
		{"san_teampos", "preinfo", "ALTER TABLE `san_teampos` ADD COLUMN `preinfo`  text NOT NULL COMMENT '编组信息' AFTER `info`;"},
		{"san_userbase", "usersignature", "ALTER TABLE `san_userbase` ADD COLUMN `usersignature` text NOT NULL COMMENT '用户签名' AFTER `redicon`;"},
		{"san_userbase", "getallgem", "ALTER TABLE `san_userbase` ADD COLUMN `getallgem`  int(11) NOT NULL DEFAULT 0 COMMENT '历史总钻石数' AFTER `usersignature`;"},
		{"san_offlineinfo", "baseinfo", "ALTER TABLE `san_offlineinfo` ADD COLUMN `baseinfo`  text NOT NULL  COMMENT '英雄基础信息' AFTER `mailinfos`;"},
		{"san_friend", "applyhire", "ALTER TABLE `san_friend` ADD COLUMN `applyhire`  text NOT NULL  COMMENT '申请租用信息' AFTER `count`;"},
		{"san_friend", "hireHero", "ALTER TABLE `san_friend` ADD COLUMN `hireHero`  text NOT NULL  COMMENT '租用信息' AFTER `applyhire`;"},
		{"san_friend", "hiretime", "ALTER TABLE `san_friend` ADD COLUMN `hiretime` int(11) NOT NULL COMMENT '到期时间' AFTER `hireHero`;"},
		{"san_friend", "usesign", "ALTER TABLE `san_friend` ADD COLUMN `usesign`  text NOT NULL  COMMENT '使用标记' AFTER `hiretime`;"},
		{"san_friend", "giftsign", "ALTER TABLE `san_friend` ADD COLUMN `giftsign`  text NOT NULL  COMMENT '使用标记' AFTER `usesign`;"},
		{"san_onhook", "hangup", "ALTER TABLE `san_onhook` ADD COLUMN `hangup` int(11) NOT NULL DEFAULT 0  COMMENT '挂机库' AFTER `caltime`;"},
		{"san_onhook", "fasttimes", "ALTER TABLE `san_onhook` ADD COLUMN `fasttimes` int(11) NOT NULL DEFAULT 0  COMMENT '快速挂机次数' AFTER `hangup`;"},
		{"san_onhook", "onhookstagetime", "ALTER TABLE `san_onhook` ADD COLUMN `onhookstagetime` bigint(20) NOT NULL DEFAULT 0  COMMENT '历史最高关卡时间' AFTER `fasttimes`;"},
		{"san_userunioninfo", "stateaward", "ALTER TABLE `san_userunioninfo` ADD COLUMN `stateaward`  int(11) NOT NULL COMMENT '州宝箱' AFTER `copyaward`;"},
		{"san_userunioninfo", "huntlimit", "ALTER TABLE `san_userunioninfo` ADD COLUMN `huntlimit`  text NOT NULL COMMENT '狩猎限制' AFTER `stateaward`;"},
		{"san_unioninfo", "activitypoint", "ALTER TABLE `san_unioninfo` ADD COLUMN `activitypoint` int(11) NOT NULL DEFAULT 0 COMMENT '活跃点数' AFTER `copy`;"},
		{"san_unioninfo", "acitvitylimit", "ALTER TABLE `san_unioninfo` ADD COLUMN `acitvitylimit` int(11) NOT NULL DEFAULT 0 COMMENT '活跃点限制' AFTER `activitypoint`;"},
		{"san_unioninfo", "huntinfo", "ALTER TABLE `san_unioninfo` ADD COLUMN `huntinfo`  text NOT NULL COMMENT '狩猎信息' AFTER `acitvitylimit`;"},
		{"san_unioninfo", "board", "ALTER TABLE `san_unioninfo` ADD COLUMN `board`  text NOT NULL COMMENT '对外公告' AFTER `huntinfo`;"},
		{"san_unioninfo", "mailcd", "ALTER TABLE `san_unioninfo` ADD COLUMN `mailcd`  bigint(20) NOT NULL COMMENT '邮件cd' AFTER `board`;"},
		{"san_unioninfo", "bravehand", "ALTER TABLE `san_unioninfo` ADD COLUMN `bravehand`  text NOT NULL COMMENT '无畏之手' AFTER `mailcd`;"},
		{"san_usernewpit", "buffstore", "ALTER TABLE `san_usernewpit` ADD COLUMN `buffstore`  text NOT NULL COMMENT 'BUFF库存' AFTER `userpitinfo`;"},
		{"san_usernewpit", "winlevel", "ALTER TABLE `san_usernewpit` ADD COLUMN `winlevel` text NOT NULL COMMENT '胜利等级' AFTER `buffstore`;"},
		{"san_usernewpit", "loselevel", "ALTER TABLE `san_usernewpit` ADD COLUMN `loselevel` text NOT NULL COMMENT '失败等级' AFTER `winlevel`;"},
		{"san_resonancecrystal", "maxfight", "ALTER TABLE `san_resonancecrystal` ADD COLUMN `maxfight`  bigint(20) NOT NULL COMMENT '最高战力' AFTER `resonancecount`;"},
		{"san_resonancecrystal", "maxfighttime", "ALTER TABLE `san_resonancecrystal` ADD COLUMN `maxfighttime`  bigint(20) NOT NULL COMMENT '最高战力生成时间' AFTER `maxfight`;"},
		{"san_resonancecrystal", "maxfightall", "ALTER TABLE `san_resonancecrystal` ADD COLUMN `maxfightall`  bigint(20) NOT NULL COMMENT '共鸣水晶战力' AFTER `maxfighttime`;"},
		{"san_userrecharge", "warorder", "ALTER TABLE `san_userrecharge` ADD COLUMN `warorder` text NOT NULL COMMENT '赏金社' AFTER `vipweekbuy`;"},
		{"san_userrecharge", "warorderlimit", "ALTER TABLE `san_userrecharge` ADD COLUMN `warorderlimit` text NOT NULL COMMENT '限时战令' AFTER `warorder`;"},
		{"san_towerplayerrecord", "level", "ALTER TABLE `san_towerplayerrecord` ADD COLUMN `level` int(11) NOT NULL COMMENT '等级' AFTER `icon`;"},
		{"san_userfindpool", "astrology", "ALTER TABLE `san_userfindpool` ADD COLUMN `astrology` text NOT NULL COMMENT '占星' AFTER `wishinfo`;"},
		{"san_userfindpool", "selfselection", "ALTER TABLE `san_userfindpool` ADD COLUMN `selfselection` text NOT NULL COMMENT '自选召唤' AFTER `astrology`;"},
		{"san_userfindpool", "gemfindgift", "ALTER TABLE `san_userfindpool` ADD COLUMN `gemfindgift` text NOT NULL COMMENT '高级招募累计奖励' AFTER `selfselection`;"},
		{"san_userequip", "info2", "ALTER TABLE `san_userequip` ADD COLUMN `info2` mediumtext NOT NULL COMMENT '装备' AFTER `starttime`;"},
		{"san_userequip", "info3", "ALTER TABLE `san_userequip` ADD COLUMN `info3` mediumtext NOT NULL COMMENT '装备' AFTER `info2`;"},
		{"san_userequip", "info4", "ALTER TABLE `san_userequip` ADD COLUMN `info4` mediumtext NOT NULL COMMENT '装备' AFTER `info3`;"},
		{"san_userequip", "info5", "ALTER TABLE `san_userequip` ADD COLUMN `info5` mediumtext NOT NULL COMMENT '装备' AFTER `info4`;"},
		{"san_usertask", "weekrefreshtime", "ALTER TABLE `san_usertask` ADD COLUMN `weekrefreshtime` bigint(20) NOT NULL DEFAULT 0 COMMENT '周任务时间戳' AFTER `award`;"},
		{"san_specialpurchase", "recharge", "ALTER TABLE `san_specialpurchase` ADD COLUMN `recharge` int(11) NOT NULL DEFAULT 0 COMMENT '充值金额' AFTER `info`;"},
		{"san_userbase", "paygem", "ALTER TABLE `san_userbase` ADD COLUMN `paygem`  int(11) NOT NULL DEFAULT 0 COMMENT '付费钻石' AFTER `getallgem`;"},
		{"san_luckshop", "ver", "ALTER TABLE `san_luckshop` ADD COLUMN `ver`  int(11) NOT NULL DEFAULT 0 COMMENT '版本' AFTER `info`;"},
		{"san_luckshop", "starver", "ALTER TABLE `san_luckshop` ADD COLUMN `starver`  int(11) NOT NULL DEFAULT 0 COMMENT '星辰礼包版本' AFTER `ver`;"},
		{"san_luckshop", "discountver", "ALTER TABLE `san_luckshop` ADD COLUMN `discountver`  int(11) NOT NULL DEFAULT 0 COMMENT '特惠礼包版本' AFTER `starver`;"},
		{"san_luckshop", "vergroup", "ALTER TABLE `san_luckshop` ADD COLUMN `vergroup` text NOT NULL COMMENT '版本组' AFTER `discountver`;"},
		{"san_onhook", "passidrecord", "ALTER TABLE `san_onhook` ADD COLUMN `passidrecord` int(11) NOT NULL DEFAULT 0  COMMENT '换日关卡进度' AFTER `onhookstagetime`;"},
		{"san_onhook", "passidrecordtime", "ALTER TABLE `san_onhook` ADD COLUMN `passidrecordtime` bigint(20) NOT NULL DEFAULT 0  COMMENT '换日进度刷新时间' AFTER `passidrecord`;"},
		{"san_onhook", "caltimeprivilege", "ALTER TABLE `san_onhook` ADD COLUMN `caltimeprivilege` text NOT NULL COMMENT '特权掉落组' AFTER `passidrecordtime`;"},
		{"san_onhook", "caltimeextitems", "ALTER TABLE `san_onhook` ADD COLUMN `caltimeextitems` text NOT NULL COMMENT '特权掉落额外奖励' AFTER `caltimeprivilege`;"},
		{"san_usernobilitytask", "getreward", "ALTER TABLE `san_usernobilitytask` ADD COLUMN `getreward` text NOT NULL COMMENT '领取标记' AFTER `level`;"},
		{"san_usertargettask", "buylevel", "ALTER TABLE `san_usertargettask` ADD COLUMN `buylevel` text NOT NULL COMMENT '徽章等级' AFTER `nobilitylevel`;"},
		{"san_accesscard", "rank", "ALTER TABLE `san_accesscard` ADD COLUMN `rank`  int(11) NOT NULL DEFAULT 0 COMMENT '当前几人领取' AFTER `accesscardrecord`;"},
		{"san_accesscard", "accesscardtop", "ALTER TABLE `san_accesscard` ADD COLUMN `accesscardtop` text NOT NULL COMMENT '当前记录' AFTER `rank`;"},
		{"san_accesscard", "starttime", "ALTER TABLE `san_accesscard` ADD COLUMN `starttime`  int(11) NOT NULL DEFAULT 0 COMMENT '开始时间' AFTER `accesscardtop`;"},
		{"san_accesscard", "endtime", "ALTER TABLE `san_accesscard` ADD COLUMN `endtime`  int(11) NOT NULL DEFAULT 0 COMMENT '结束时间' AFTER `starttime`;"},
		{"san_accesscard", "rewardtime", "ALTER TABLE `san_accesscard` ADD COLUMN `rewardtime`  int(11) NOT NULL DEFAULT 0 COMMENT '发奖时间' AFTER `endtime`;"},
		{"san_accesscard", "hasreward", "ALTER TABLE `san_accesscard` ADD COLUMN `hasreward`  int(11) NOT NULL DEFAULT 0 COMMENT '是否发奖' AFTER `rewardtime`;"},
		{"san_activity", "activityresetsign", "ALTER TABLE `san_activity` ADD COLUMN `activityresetsign` text NOT NULL COMMENT '活动重置时间' AFTER `fund`;"},
		{"san_specialpurchase", "sign", "ALTER TABLE `san_specialpurchase` ADD COLUMN `sign` text NOT NULL COMMENT '弹出标记' AFTER `recharge`;"},
		{"san_crossarena", "showtime", "ALTER TABLE `san_crossarena` ADD COLUMN `showtime` int(11) NOT NULL DEFAULT 0 COMMENT '展示时间' AFTER `endtime`;"},
		{"san_usercrossarena", "showtime", "ALTER TABLE `san_usercrossarena` ADD COLUMN `showtime` int(11) NOT NULL DEFAULT 0 COMMENT '展示时间' AFTER `endtime`;"},
		{"san_usernewpit", "shop", "ALTER TABLE `san_usernewpit` ADD COLUMN `shop` text NOT NULL COMMENT '地牢商店' AFTER `loselevel`;"},
		{"san_userunioninfo", "huntstart", "ALTER TABLE `san_userunioninfo` ADD COLUMN `huntstart`  int(11) NOT NULL DEFAULT 0 COMMENT '开启的狩猎' AFTER `huntlimit`;"},
		{"san_userfund", "fundinfo", "ALTER TABLE `san_userfund` ADD COLUMN `fundinfo` text NOT NULL COMMENT '基金信息' AFTER `taskinfo`;"},
	}

	self.modifyField = [][4]string{
		{"san_usertask", "liveness", "int", "ALTER TABLE `san_usertask` MODIFY COLUMN `liveness` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '活跃度信息' AFTER `taskinfo`;"},
		{"san_userranktask", "getstate", "bigint", "ALTER TABLE `san_userranktask` MODIFY COLUMN `getstate` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '领取状态' AFTER `taskinfo`;"},
		{"san_userhero2", "totalstars", "int", "ALTER TABLE `san_userhero2` MODIFY COLUMN `totalstars` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '总星级信息' AFTER `info`;"},
		{"san_usertreasure", "info", "text", "ALTER TABLE `san_usertreasure` MODIFY COLUMN `info`  mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '宝物信息' AFTER `washfreetimes`;"},
		{"san_tiger", "info", "text", "ALTER TABLE `san_tiger` MODIFY COLUMN `info`  mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '虎符信息' AFTER `maxkey`;"},
		{"san_armsarena1", "pvpfight", "text", "ALTER TABLE `san_armsarena1` MODIFY COLUMN `pvpfight`  mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '战报信息' AFTER `maxkey`;"},
		{"san_armsarena2", "pvpfight", "text", "ALTER TABLE `san_armsarena2` MODIFY COLUMN `pvpfight`  mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '战报信息' AFTER `maxkey`;"},
		{"san_armsarena3", "pvpfight", "text", "ALTER TABLE `san_armsarena3` MODIFY COLUMN `pvpfight`  mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '战报信息' AFTER `maxkey`;"},
		{"san_activityboss", "activitybosstop", "text", "ALTER TABLE `san_activityboss` MODIFY COLUMN `activitybosstop`  mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '记录' AFTER `period`;"},
		{"san_accesscard", "accesscardtop", "text", "ALTER TABLE `san_accesscard` MODIFY COLUMN `accesscardtop`  mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '当前记录' AFTER `rank`;"},
		{"san_usercrossarena", "starttime", "int", "ALTER TABLE `san_usercrossarena` MODIFY COLUMN `starttime`  bigint(20) NOT NULL DEFAULT 0 COMMENT '开始时间' AFTER `buytimes`;"},
		{"san_usercrossarena", "endtime", "int", "ALTER TABLE `san_usercrossarena` MODIFY COLUMN `endtime`  bigint(20) NOT NULL DEFAULT 0 COMMENT '结束时间' AFTER `starttime`;"},
		{"san_usercrossarena", "showtime", "int", "ALTER TABLE `san_usercrossarena` MODIFY COLUMN `showtime`  bigint(20) NOT NULL DEFAULT 0 COMMENT '展示时间' AFTER `endtime`;"},
	}

	self.createTableStatements = []string{
		// 红包全局信息
		GetRedPacMgr().GlobalRedPacSql(),
		// 红包个人信息
		GetRedPacMgr().UserRedPacSql(),
		// 转盘
		GetDialMgr().UserDialSql(),
		// 翻牌
		self.UserDrawSql(),
		// 宝物
		self.UserTreasureSql(),
		// 砸金蛋
		self.UserLuckEggSql(),
		// 开工福利
		self.UserLuckStartSql(),
		// 连续充值
		self.UserDailyRechargeSql(),
		// 战报数据
		self.FightResultSql(),
		// GVG玩家数据
		self.GVGUserSql(),
		// GVG奖励数据
		self.GVGAwardSql(),
		// 布阵数据
		self.TeamPosSql(),
		// 地下城数据
		self.DungeonSql(),
		// 地下城玩家数据
		self.DungeonUserSql(),
		// 装备数据
		self.UserEquipSql(),
		// 科技数据
		self.UserTech(),
		// 镇魂塔数据
		self.UserTowerSql(),
		// 巨兽
		self.UserBoss(),
		// 宝石副本数据
		self.UserGemStoneSql(),
		// 镇魂塔商店
		self.TowerShopSql(),
		// 神格商店
		self.DinivityShopSql(),
		self.ShopSql1(),
		self.ShopSql2(),
		self.ShopSql3(),
		self.ShopSql4(),
		self.ShopSql5(),
		// 佣兵
		self.UserSoldier(),
		// 王国任务
		self.UserKingTask(),
		// 事件信息
		self.UserEvents1(),
		self.UserEvents2(),
		self.UserEvents3(),
		// 矿点相关
		self.MineGlobal(),
		self.MinePlayer(),
		// 孤山夺宝
		self.GvePlayer(),
		self.GveGlobal(),
		// 鼓舞
		self.Encourage(),
		// 头像
		self.HeadInfo(),
		// 赏金任务
		self.MoneyTask(),
		// 新版佣兵
		self.ArmyMySql(),
		// 军团战
		self.UnionFightSQL(),
		self.UnionFightGlobal(),
		self.Guides(),
		self.PassRecord(),
		self.DreamLandRecord(),
		//	系统总览
		self.Statistics(),

		//! 限时活动
		self.UserTimeGiftSql(),
		self.NobilityTask(),
		self.WholeShop(),
		self.AccessCard(),
		self.OnHook(),
		self.UserRuneSql(),
		//神兽
		self.UserHydra(),
		//地牢
		self.UserPit(),
		//新地牢
		self.UserNewPit(),
		//时光之巅
		self.UserInstance(),
		//客户端标记
		self.UserClienSign(),
		//神器
		self.UserArtifactEquipSql(),
		// 支援英雄
		self.SupportHero(),
		self.Entanglement(),
		self.RewardMod(),
		self.OfflineInfo(),
		self.HireHeroInfo(),
		self.UserRankTask(),
		self.MgrRankTask(),
		self.UserFindInfo(),
		self.MgrResonanceCrystal(),
		self.KingTower(),
		self.KingTowerMgr(),
		self.ArenaRank1(),
		self.ArenaRank2(),
		self.ArenaRank3(),
		self.ArenaPlayer1(),
		self.ArenaPlayer2(),
		self.ArenaPlayer3(),
		self.ActivityGift(),
		self.ArenaTime(),
		self.GrowthGift(),
		self.HeroSkin(),
		self.SpecialPurchase(),
		self.LifeTree(),
		self.UserTargetTask(),
		self.UserTurnTable(),
		self.UserVipRecharge(),
		self.UserAccessCard(),
		self.UserFund(),
		self.UserInterStellar(),
		self.ActivityBoss(),
		self.UserActivityBoss(),
		self.UserHeroGrow(),
		self.UserCrossArena(),
		self.UserCrossArena3V3(),
		self.CrossArena(),
		self.CrossArena3V3(),
		self.BattleRecord(),
		self.UserActivityBossFestival(),
		self.OfflineRecordInfo(),
		self.UserLotteryDraw(),
		self.HonourShop(),
		self.RankReward1732(),
		self.UserBase2(),
	}

	self.dropField = [][3]string{
		{"san_userhero2", "artifactstars", "ALTER TABLE `san_userhero2` DROP COLUMN `artifactstars`;"},
		{"san_userfind", "dicountinfo", "ALTER TABLE `san_userfind` DROP COLUMN `dicountinfo`;"},
		{"san_userbase", "goldfreetimes", "ALTER TABLE `san_userbase` DROP COLUMN `goldfreetimes`;"},
		{"san_userbase", "juqingid", "ALTER TABLE `san_userbase` DROP COLUMN `juqingid`;"},
		{"san_userbase", "juqingid2", "ALTER TABLE `san_userbase` DROP COLUMN `juqingid2`;"},
		{"san_userbase", "zhiyinid", "ALTER TABLE `san_userbase` DROP COLUMN `zhiyinid`;"},
		{"san_userbase", "zhiyinid1", "ALTER TABLE `san_userbase` DROP COLUMN `zhiyinid1`;"},
		{"san_mineplayer", "move_time", "ALTER TABLE `san_mineplayer` DROP COLUMN `move_time`;"},
	}
}

func (self *MySqlMgr) createTables() {
	for _, stmt := range self.createTableStatements {
		_, _, res := GetServer().DBUser.Exec(stmt)
		if !res {
			LogError("创建table失败")
			os.Exit(1)
		}
	}
}

// 检查字段是否存在
func (self *MySqlMgr) CheckMysql() {
	self.initData()
	var checkErr error
	for _, filedName := range self.fieldCheck {
		checkErr = self.CheckFiled(filedName[0], filedName[1])
		if checkErr != nil {
			LogError(checkErr.Error())
			os.Exit(1)
		}
	}

	for index := range self.tableCheck {
		checkErr = self.CheckTable(self.tableCheck[index])
		if checkErr != nil {
			LogError(checkErr.Error())
			os.Exit(1)
		}
	}

	for index := range self.logTableCheck {
		checkErr = self.CheckLogTable(self.logTableCheck[index])
		if checkErr != nil {
			LogError(checkErr.Error())
			os.Exit(1)
		}
	}

	self.createTables()
	self.checkAddField()
	self.modifyColumnType()
	self.checkDropField()
}

func (self *MySqlMgr) CheckFiled(tableName string, filedName string) error {
	sql := "SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'  AND COLUMN_NAME = '%s'"
	dbName := self.getDbName()
	if dbName == "" {
		LogDebug("dbName is empty!")
		return nil
	}
	sqlStr := fmt.Sprintf(sql, dbName, tableName, filedName)
	res := GetServer().DBUser.Query(sqlStr)
	//LogDebug("sqlStr:", sqlStr, ", CheckFiledParam:", res)
	if res {
		return nil
	}
	return errors.New(tableName + " has no filed:" + filedName)
}

func (self *MySqlMgr) CheckTable(tableName string) error {
	sql := "SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'"
	dbName := self.getDbName()
	if dbName == "" {
		LogDebug("dbName is empty!")
		return nil
	}
	sqlStr := fmt.Sprintf(sql, dbName, tableName)
	res := GetServer().DBUser.Query(sqlStr)
	//LogDebug("sqlStr:", sqlStr, ", CheckFiledParam:", res)
	if res {
		return nil
	}
	return errors.New(tableName + " not exists!")
}

func (self *MySqlMgr) getDbName() string {
	sqlSplit1 := strings.Split(GetServer().Con.DBCon.DBUser, "?")
	if len(sqlSplit1) < 1 {
		return ""
	}
	sqlSplit2 := strings.Split(sqlSplit1[0], "/")
	if len(sqlSplit2) < 2 {
		return ""
	}
	return sqlSplit2[1]
}

func (self *MySqlMgr) getDbLogName() string {
	sqlSplit1 := strings.Split(GetServer().Con.DBCon.DBLog, "?")
	if len(sqlSplit1) < 1 {
		return ""
	}
	sqlSplit2 := strings.Split(sqlSplit1[0], "/")
	if len(sqlSplit2) < 2 {
		return ""
	}
	return sqlSplit2[1]
}

func (self *MySqlMgr) CheckLogTable(tableName string) error {
	sql := "SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'"
	dbName := self.getDbLogName()
	if dbName == "" {
		LogDebug("dbName is empty!")
		return nil
	}
	sqlStr := fmt.Sprintf(sql, dbName, tableName)
	res := GetServer().DBLog.Query(sqlStr)
	if res {
		return nil
	}
	return errors.New(tableName + " not exists!")
}

// 翻牌
func (self *MySqlMgr) UserDrawSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userdraw (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  times int(11) NOT NULL DEFAULT '0' COMMENT '翻牌次数',
		  step int(11) NOT NULL DEFAULT '1' COMMENT '活动期数',
		  freetime int(11) NOT NULL DEFAULT '0' COMMENT '免费次数',
		  lootgroup int(11) NOT NULL DEFAULT '1' COMMENT '抽奖组数',
		  resettimes int(11) NOT NULL DEFAULT '0' COMMENT '重置次数',
		  boxinfo text NOT NULL COMMENT '宝箱状态信息',
		  drawinfo text NOT NULL COMMENT '抽奖状态信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 宝物
func (self *MySqlMgr) UserTreasureSql() string {
	return `CREATE TABLE IF NOT EXISTS san_usertreasure (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  maxkey int(11) NOT NULL DEFAULT '0' COMMENT '宝物最大keyId',
		  washfreetimes int(11) NOT NULL DEFAULT '0' COMMENT '免费洗练次数',
		  info text NOT NULL COMMENT '宝物信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 砸金蛋信息
func (self *MySqlMgr) UserLuckEggSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userluckegg (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  rechargenum int(11) NOT NULL DEFAULT '0' COMMENT '充值获得锤子个数',
		  rechargeleft int(11) NOT NULL DEFAULT '0' COMMENT '充值剩余',
		  step int(11) NOT NULL DEFAULT '0' COMMENT '第几期',
		  lootgroup int(11) NOT NULL DEFAULT '0' COMMENT '当前是哪一组',
		  info text NOT NULL COMMENT '砸金蛋信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

//Uid  int64  //! 玩家uid
//Step int    //! 第几期
//Info string //! 任务状态
// 开工福利
func (self *MySqlMgr) UserLuckStartSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userluckstart (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  step int(11) NOT NULL DEFAULT '0' COMMENT '第几期',
		  info text NOT NULL COMMENT '活动信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 全服限时礼包
func (self *MySqlMgr) UserTimeGiftSql() string {
	return `CREATE TABLE IF NOT EXISTS san_timegift (
			uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
			step int(11) NOT NULL DEFAULT '0' COMMENT '第几期',
			info text NOT NULL COMMENT '活动信息',
			updatetime int(11) NOT NULL COMMENT '最后更新时间',
			PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

//Uid  int64  //! 玩家uid
//Step int    //! 第几期
//Info string //! 任务状态
// 连续充值
func (self *MySqlMgr) UserDailyRechargeSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userdailyrecharge (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  step int(11) NOT NULL DEFAULT '0' COMMENT '第几期',
		  info text NOT NULL COMMENT '活动信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) FightResultSql() string {
	return `CREATE TABLE IF NOT EXISTS san_fightresult (
		  id bigint(20) NOT NULL COMMENT '战斗id',
		  record text NOT NULL COMMENT '战斗内容',
		  ts bigint(20) DEFAULT NULL COMMENT '时间戳'
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) GVGUserSql() string {
	return `CREATE TABLE IF NOT EXISTS san_king_cityfight (
  uid bigint(20) NOT NULL COMMENT 'uid',
  info text NOT NULL COMMENT '玩家数据',
  cityid int(11) NOT NULL DEFAULT '0' COMMENT 'cityid',
  pos int(11) NOT NULL DEFAULT '0' COMMENT '位置',
  award text NOT NULL COMMENT '奖励',
  title int(11) NOT NULL DEFAULT '0' COMMENT '威名等级',
  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) GVGAwardSql() string {
	return `CREATE TABLE IF NOT EXISTS san_king_cityfight_award (
  id int(11) NOT NULL COMMENT 'id',
  city1 int(11) NOT NULL COMMENT '城池1id',
  award1 text NOT NULL COMMENT '城池1奖励',
  city2 int(11) NOT NULL COMMENT '城池2id',
  award2 text NOT NULL COMMENT '城池2奖励',
  city3 int(11) NOT NULL COMMENT '城池3id',
  award3 text NOT NULL COMMENT '城池3奖励',
  idtimes text NOT NULL COMMENT '奖励次数',
  lastts int(11) NOT NULL COMMENT '时间戳'
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) checkAddField() {
	var checkErr error
	for _, stmt := range self.addField {
		checkErr = self.CheckFiled(stmt[0], stmt[1])
		if checkErr != nil { // 没有才插入
			_, _, res := GetServer().DBUser.Exec(stmt[2])
			if !res {
				LogError("增加字段失败", stmt[2])
				os.Exit(1)
			}
		}
	}
}

func (self *MySqlMgr) checkDropField() {
	var checkErr error
	for _, stmt := range self.dropField {
		checkErr = self.CheckFiled(stmt[0], stmt[1])
		if checkErr == nil { // 没有才插入
			_, _, res := GetServer().DBUser.Exec(stmt[2])
			if !res {
				LogError("删除字段失败", stmt[2])
				os.Exit(1)
			}
		}
	}

}

// ALTER TABLE `san_usertreasure` MODIFY COLUMN `info`  mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '宝物信息' AFTER `washfreetimes`;
func (self *MySqlMgr) modifyColumnType() {
	for _, stmt := range self.modifyField {
		if len(stmt) != 4 {
			continue
		}

		isStmtOk := true
		for _, v := range stmt {
			if v == "" {
				isStmtOk = false
				break
			}
		}

		if !isStmtOk {
			break
		}

		sql := "SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'  AND COLUMN_NAME = '%s'"
		dbName := self.getDbName()
		if dbName == "" {
			LogError("dbName is empty!")
			os.Exit(1)
		}

		sqlStr := fmt.Sprintf(sql, dbName, stmt[0], stmt[1])
		res, fieldName := GetServer().DBUser.QueryColomn(sqlStr)
		if res && fieldName == stmt[2] {
			_, _, res := GetServer().DBUser.Exec(stmt[3])
			if !res {
				LogError("修改字段类型失败", stmt[3])
				os.Exit(1)
			} else {
				LogDebug("修改字段类型成功!")
			}
		}
	}

}

func (self *MySqlMgr) TeamPosSql() string {
	return `CREATE TABLE IF NOT EXISTS san_teampos (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  info text NOT NULL COMMENT '阵型信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 地下城
func (self *MySqlMgr) DungeonSql() string {
	return `CREATE TABLE IF NOT EXISTS san_dungeon (
instanceid int(11) NOT NULL COMMENT '组队副本编号',
lefttimes int(11) NOT NULL COMMENT '剩下次数',
fulltimes int(11) NOT NULL COMMENT '总共次数',
nextfreshts bigint(20) NOT NULL COMMENT '下次刷新时间戳'
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 地下城
func (self *MySqlMgr) DungeonUserSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userdungeon (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
playts text NOT NULL COMMENT '时间戳',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 装备
func (self *MySqlMgr) UserEquipSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userequip (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  maxkey int(11) NOT NULL DEFAULT '0' COMMENT '装备最大keyId',
		  info mediumtext NOT NULL COMMENT '装备信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 科技
func (self *MySqlMgr) UserTech() string {
	return `CREATE TABLE IF NOT EXISTS san_usertech (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  maxkey int(11) NOT NULL DEFAULT '0' COMMENT '科技最大keyId',
		  techinfo mediumtext NOT NULL COMMENT '科技信息',
		  workque text NOT NULL COMMENT '升级队列',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 镇魂塔
func (self *MySqlMgr) UserTowerSql() string {
	return `CREATE TABLE IF NOT EXISTS san_tower (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
maxlevel int(11) NOT NULL COMMENT '历史最大关卡',
maxlevelts bigint(20) NOT NULL COMMENT '历史最大关卡时间戳',
curlevel int(11) NOT NULL COMMENT '当前关卡',
curfailbox int(11) NOT NULL COMMENT '镇魂秘宝',
boxstate text NOT NULL  COMMENT '关卡宝箱状态',
levelstate text NOT NULL COMMENT '关卡通关状态',
resettimes int(11) NOT NULL COMMENT '重置次数',
resetbuytimes int(11) NOT NULL COMMENT '重置购买次数',
buff text NOT NULL COMMENT '祝福BUFF列表',
curbuff int(11) NOT NULL COMMENT '当前祝福',
bufftimes int(11) NOT NULL COMMENT '祝福重置次数',
buffbuytimes int(11) NOT NULL COMMENT '祝福重置购买次数',
advancelevelstate text NOT NULL COMMENT '精英关卡通关状态',
advancetimes int(11) NOT NULL COMMENT '精英关卡挑战次数',
advancebuytimes int(11) NOT NULL COMMENT '精英关卡购买次数',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 巨兽
func (self *MySqlMgr) UserBoss() string {
	return `CREATE TABLE IF NOT EXISTS san_userboss (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  bossinfo text NOT NULL COMMENT '巨兽信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 宝石副本
func (self *MySqlMgr) UserGemStoneSql() string {
	return `CREATE TABLE IF NOT EXISTS san_gemstone (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
curlevel int(11) NOT NULL COMMENT '当前关卡',
sweeptimes text NOT NULL COMMENT '扫荡次数',
buysweeptimes text NOT NULL COMMENT '购买次数',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 镇魂塔商店
func (self *MySqlMgr) TowerShopSql() string {
	return `CREATE TABLE IF NOT EXISTS san_shop10  (
  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
  shoptype int(11) NOT NULL,
  shopgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  refindex int(11) NOT NULL,
  shopnextgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  todayrefcount int(11) NOT NULL,
  sysreftime int(11) NOT NULL,
  lastupdtime int(11) NOT NULL,
  PRIMARY KEY (uid) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`
}

// 神格商店
func (self *MySqlMgr) DinivityShopSql() string {
	return `CREATE TABLE IF NOT EXISTS san_shop12  (
  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
  shoptype int(11) NOT NULL,
  shopgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  refindex int(11) NOT NULL,
  shopnextgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  todayrefcount int(11) NOT NULL,
  sysreftime int(11) NOT NULL,
  lastupdtime int(11) NOT NULL,
  PRIMARY KEY (uid) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`
}

// 镇魂塔商店
func (self *MySqlMgr) ShopSql1() string {
	return `CREATE TABLE IF NOT EXISTS san_shop1  (
  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
  shoptype int(11) NOT NULL,
  shopgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  refindex int(11) NOT NULL,
  shopnextgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  todayrefcount int(11) NOT NULL,
  sysreftime int(11) NOT NULL,
  lastupdtime int(11) NOT NULL,
  PRIMARY KEY (uid) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`
}

func (self *MySqlMgr) ShopSql2() string {
	return `CREATE TABLE IF NOT EXISTS san_shop2  (
  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
  shoptype int(11) NOT NULL,
  shopgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  refindex int(11) NOT NULL,
  shopnextgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  todayrefcount int(11) NOT NULL,
  sysreftime int(11) NOT NULL,
  lastupdtime int(11) NOT NULL,
  PRIMARY KEY (uid) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`
}

func (self *MySqlMgr) ShopSql3() string {
	return `CREATE TABLE IF NOT EXISTS san_shop3  (
  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
  shoptype int(11) NOT NULL,
  shopgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  refindex int(11) NOT NULL,
  shopnextgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  todayrefcount int(11) NOT NULL,
  sysreftime int(11) NOT NULL,
  lastupdtime int(11) NOT NULL,
  PRIMARY KEY (uid) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`
}

func (self *MySqlMgr) ShopSql4() string {
	return `CREATE TABLE IF NOT EXISTS san_shop4  (
  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
  shoptype int(11) NOT NULL,
  shopgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  refindex int(11) NOT NULL,
  shopnextgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  todayrefcount int(11) NOT NULL,
  sysreftime int(11) NOT NULL,
  lastupdtime int(11) NOT NULL,
  PRIMARY KEY (uid) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`
}

func (self *MySqlMgr) ShopSql5() string {
	return `CREATE TABLE IF NOT EXISTS san_shop5  (
  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
  shoptype int(11) NOT NULL,
  shopgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  refindex int(11) NOT NULL,
  shopnextgood text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  todayrefcount int(11) NOT NULL,
  sysreftime int(11) NOT NULL,
  lastupdtime int(11) NOT NULL,
  PRIMARY KEY (uid) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;`
}

// 佣兵
func (self *MySqlMgr) UserSoldier() string {
	return `CREATE TABLE IF NOT EXISTS san_usersoldier (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  maxkey int(11) NOT NULL DEFAULT '0' COMMENT '佣兵最大keyId',
		  info mediumtext NOT NULL COMMENT '佣兵信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 王国任务
func (self *MySqlMgr) UserKingTask() string {
	return `CREATE TABLE IF NOT EXISTS san_userkingtask (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  info text NOT NULL COMMENT '王国任务',
		  actiontimes int(11) NOT NULL,
		  flushtimes int(11) NOT NULL,
		  buytimes int(11) NOT NULL,
		  taskgroup int(11) NOT NULL,
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 关卡任务
func (self *MySqlMgr) UserEvents1() string {
	return `CREATE TABLE IF NOT EXISTS san_userevents1(
	uid bigint(20) NOT NULL,
	info1 text NOT NULL,
	info2 text NOT NULL,
	info3 text NOT NULL,
	info4 text NOT NULL,
	info5 text NOT NULL,
	info6 text NOT NULL,
	info7 text NOT NULL,
	info8 text NOT NULL,
	info9 text NOT NULL,
	info10 text NOT NULL,
	info11 text NOT NULL,
	info12 text NOT NULL,
	info13 text NOT NULL,
	info14 text NOT NULL,
	info15 text NOT NULL,
	info16 text NOT NULL,
	info17 text NOT NULL,
	info18 text NOT NULL,
	info19 text NOT NULL,
	info20 text NOT NULL,
	info21 text NOT NULL,
	info22 text NOT NULL,
	info23 text NOT NULL,
	info24 text NOT NULL,
	info25 text NOT NULL,
	info26 text NOT NULL,
	info27 text NOT NULL,
	info28 text NOT NULL,
	info29 text NOT NULL,
	info30 text NOT NULL,
	PRIMARY KEY (uid)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserEvents2() string {
	return `CREATE TABLE IF NOT EXISTS san_userevents2(
	uid bigint(20) NOT NULL,
	info1 text NOT NULL,
	info2 text NOT NULL,
	info3 text NOT NULL,
	info4 text NOT NULL,
	info5 text NOT NULL,
	info6 text NOT NULL,
	info7 text NOT NULL,
	info8 text NOT NULL,
	info9 text NOT NULL,
	info10 text NOT NULL,
	info11 text NOT NULL,
	info12 text NOT NULL,
	info13 text NOT NULL,
	info14 text NOT NULL,
	info15 text NOT NULL,
	info16 text NOT NULL,
	info17 text NOT NULL,
	info18 text NOT NULL,
	info19 text NOT NULL,
	info20 text NOT NULL,
	info21 text NOT NULL,
	info22 text NOT NULL,
	info23 text NOT NULL,
	info24 text NOT NULL,
	info25 text NOT NULL,
	info26 text NOT NULL,
	info27 text NOT NULL,
	info28 text NOT NULL,
	info29 text NOT NULL,
	info30 text NOT NULL,
	PRIMARY KEY (uid)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserEvents3() string {
	return `CREATE TABLE IF NOT EXISTS san_userevents3(
	uid bigint(20) NOT NULL,
	info1 text NOT NULL,
	info2 text NOT NULL,
	info3 text NOT NULL,
	info4 text NOT NULL,
	info5 text NOT NULL,
	info6 text NOT NULL,
	info7 text NOT NULL,
	info8 text NOT NULL,
	info9 text NOT NULL,
	info10 text NOT NULL,
	info11 text NOT NULL,
	info12 text NOT NULL,
	info13 text NOT NULL,
	info14 text NOT NULL,
	info15 text NOT NULL,
	info16 text NOT NULL,
	info17 text NOT NULL,
	info18 text NOT NULL,
	info19 text NOT NULL,
	info20 text NOT NULL,
	info21 text NOT NULL,
	info22 text NOT NULL,
	info23 text NOT NULL,
	info24 text NOT NULL,
	info25 text NOT NULL,
	info26 text NOT NULL,
	info27 text NOT NULL,
	info28 text NOT NULL,
	info29 text NOT NULL,
	info30 text NOT NULL,
	PRIMARY KEY (uid)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 全局矿点
func (self *MySqlMgr) MineGlobal() string {
	return `CREATE TABLE IF NOT EXISTS san_mineglobal(
	id bigint(20) NOT NULL,
	noticetime bigint(20) NOT NULL,
	minestarttime bigint(20) NOT NULL,
	mineendtime bigint(20) NOT NULL,
	buffedmineid int(11) NOT NULL,
	buffedtimecd bigint(20) NOT NULL,
	scorechanges text NOT NULL,
	minerankinfo text NOT NULL,
	buffcreatedinfo text NOT NULL,
	minerelationinfo text NOT NULL,
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 矿点玩家信息
func (self *MySqlMgr) MinePlayer() string {
	return `CREATE TABLE IF NOT EXISTS san_mineplayer(
	uid bigint(20) NOT NULL,
	camp int(11) NOT NULL,
	mineid int(11) NOT NULL,
	movecd bigint(20) NOT NULL,
	collectcd bigint(20) NOT NULL,
	spycd bigint(20) NOT NULL,
	fightcd bigint(20) NOT NULL,
	oldmineid int(11) NOT NULL,
	movereason int(11) NOT NULL,
	attacktimes int(11) NOT NULL,
	defencetimes int(11) NOT NULL,
	killnum int(11) NOT NULL,
	deadnum int(11) NOT NULL,
	collectnum int(11) NOT NULL,
	fightrecordinfo text NOT NULL,
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 孤点争夺玩家
func (self *MySqlMgr) GvePlayer() string {
	return `CREATE TABLE IF NOT EXISTS san_gveplayer(
	uid bigint(20) NOT NULL,
	camp int(11) NOT NULL,
	buildid int(11) NOT NULL,
	movecd bigint(20) NOT NULL,
	fightcd bigint(20) NOT NULL,
	gvematerial int(11) NOT NULL,
	result int(11) NOT NULL,
	lastlevelid int(11) NOT NULL,
	gveglory int(11) NOT NULL,
	skillcd bigint(20) NOT NULL,
	flushcd bigint(20) NOT NULL,
	leveldata text NOT NULL,
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) GveGlobal() string {
	return `CREATE TABLE IF NOT EXISTS san_gveglobal(
	id bigint(20) NOT NULL,
	noticetime bigint(20) NOT NULL,
	starttime bigint(20) NOT NULL,
	endtime bigint(20) NOT NULL,
	maxkey int(11) NOT NULL,
	rankinfo text NOT NULL,
	gvebuildinfo text NOT NULL,
	campmaterialinfo text NOT NULL,
	draginfo text NOT NULL,
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) Encourage() string {
	return `CREATE TABLE IF NOT EXISTS san_userencourage(
	uid bigint(20) NOT NULL,
	buffinfo text NOT NULL,
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) HeadInfo() string {
	return `CREATE TABLE IF NOT EXISTS san_userhead(
	uid bigint(20) NOT NULL,
	info text NOT NULL,
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) MoneyTask() string {
	return `CREATE TABLE IF NOT EXISTS san_moneytask(
	uid bigint(20) NOT NULL,
	flushtimes int(11) NOT NULL,
	star int(11) NOT NULL,
	info text NOT NULL,
	taskrecord text NOT NULL,
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) ArmyMySql() string {
	return `CREATE TABLE IF NOT EXISTS san_userarmy(
	uid bigint(20) NOT NULL,
	info text NOT NULL,
	flag text NOT NULL,
	freetimes int(11) NOT NULL,
	buytimes  int(11) NOT NULL,
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 军团战
func (self *MySqlMgr) UnionFightSQL() string {
	return `CREATE TABLE IF NOT EXISTS san_unionfight (
	stateid int(11) NOT NULL,
	attenders mediumtext NOT NULL,
	actionstate int(11) NOT NULL,
	records mediumtext NOT NULL,
	unionplans text NOT NULL,
	waittime bigint(20) NOT NULL,
	roundtime bigint(20) NOT NULL,
	PRIMARY KEY (stateid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UnionFightGlobal() string {
	return `CREATE TABLE IF NOT EXISTS san_unionfightglobalv5(
	id bigint(20) NOT NULL,
	noticetime bigint(20) NOT NULL,
	starttime bigint(20) NOT NULL,
	endtime bigint(20) NOT NULL,
	callstarttime bigint(20) NOT NULL,
	callendtime bigint(20) NOT NULL,
	enrollstarttime bigint(20) NOT NULL,
	enrollendtime bigint(20) NOT NULL,
	standardtime bigint(20) NOT NULL,
	state int(11) NOT NULL,
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) Guides() string {
	return `CREATE TABLE IF NOT EXISTS san_guides (
	uid bigint(20) NOT NULL,
	info text NOT NULL,
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) PassRecord() string {
	return `CREATE TABLE IF NOT EXISTS san_passrecords(
	passid int(11) NOT NULL,
	firstteam text NOT NULL,
	lowteam text NOT NULL,
	recentteam text NOT NULL,
	PRIMARY KEY (passid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 神格抽奖
func (self *MySqlMgr) DreamLandRecord() string {
	return `CREATE TABLE IF NOT EXISTS san_dreamland(
	uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
	info text NOT NULL COMMENT '抽奖信息',
	guidetype int(11) NOT NULL COMMENT '引导状态',
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) Statistics() string {
	return `CREATE TABLE IF NOT EXISTS san_statistics (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
statistics text NOT NULL COMMENT '统计info',
score int(11) NOT NULL COMMENT '积分',
rewardSign text NOT NULL COMMENT '奖励领取标记',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) NobilityTask() string {
	return `CREATE TABLE IF NOT EXISTS san_usernobilitytask (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
taskinfo text NOT NULL COMMENT '任务信息',
level int(11) NOT NULL COMMENT '爵位',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) WholeShop() string {
	return `CREATE TABLE IF NOT EXISTS san_wholeshop (
id  int(11) NOT NULL COMMENT '活动ID',
stage int(11) NOT NULL COMMENT '活动期数',
shopinfo text NOT NULL COMMENT '商品信息',
PRIMARY KEY (id,stage)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) AccessCard() string {
	return `CREATE TABLE IF NOT EXISTS san_accesscard (
id  int(11) NOT NULL COMMENT '活动ID期数',
ngroup int(11) NOT NULL COMMENT '活动组',
accesscardrecord text NOT NULL COMMENT '记录',
PRIMARY KEY (id,ngroup)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) OnHook() string {
	return `CREATE TABLE IF NOT EXISTS san_onhook (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
gettime int(11) NOT NULL COMMENT '领取的时间',
onhookstage int(11) NOT NULL COMMENT '挂机关卡',
caltime text NOT NULL COMMENT '挂机时间，辅助计算',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 装备
func (self *MySqlMgr) UserRuneSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userrune (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  maxkey int(11) NOT NULL DEFAULT '0' COMMENT '装备最大keyId',
		  info mediumtext NOT NULL COMMENT '符文信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 神兽
func (self *MySqlMgr) UserHydra() string {
	return `CREATE TABLE IF NOT EXISTS san_userhydra (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  hydrainfo text NOT NULL COMMENT '神兽信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 地牢
func (self *MySqlMgr) UserPit() string {
	return `CREATE TABLE IF NOT EXISTS san_userpit (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  pitinfo text NOT NULL COMMENT '地牢信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 新地牢
func (self *MySqlMgr) UserNewPit() string {
	return `CREATE TABLE IF NOT EXISTS san_usernewpit (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  newpitinfo text NOT NULL COMMENT '地牢信息',
		  userpitinfo text NOT NULL COMMENT '玩家信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 时光之巅
func (self *MySqlMgr) UserInstance() string {
	return `CREATE TABLE IF NOT EXISTS san_userinstance (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  instanceinfo text NOT NULL COMMENT '总体进度',
		  nowinstancestate text NOT NULL COMMENT '当前地图',
		  buffstore text NOT NULL COMMENT '库存',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 客户端标记
func (self *MySqlMgr) UserClienSign() string {
	return `CREATE TABLE IF NOT EXISTS san_clientsign (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  sign text NOT NULL COMMENT '标记',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 支援英雄
func (self *MySqlMgr) SupportHero() string {
	return `CREATE TABLE IF NOT EXISTS san_supporthero (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  supporthero text NOT NULL COMMENT '支援英雄状态',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 神器
func (self *MySqlMgr) UserArtifactEquipSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userartifactequip (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  maxkey int(11) NOT NULL DEFAULT '0' COMMENT '装备最大keyId',
		  info mediumtext NOT NULL COMMENT '装备信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 羁绊
func (self *MySqlMgr) Entanglement() string {
	return `CREATE TABLE IF NOT EXISTS san_entanglement (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  info text NOT NULL COMMENT '羁绊',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 悬赏
func (self *MySqlMgr) RewardMod() string {
	return `CREATE TABLE IF NOT EXISTS san_reward (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  info text NOT NULL COMMENT '悬赏',
		  level int(11) NOT NULL COMMENT '等级',
		  taskcount text NOT NULL COMMENT '任务进度',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 离线信息管理
func (self *MySqlMgr) OfflineInfo() string {
	return `CREATE TABLE IF NOT EXISTS san_offlineinfo (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  friendpowerinfos text NOT NULL COMMENT '好友赠送信息',
		  mailinfos text NOT NULL COMMENT '邮件信息（预留字段）',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 好友雇佣
func (self *MySqlMgr) HireHeroInfo() string {
	return `CREATE TABLE IF NOT EXISTS san_hireheroinfo (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  hireheroinfo text NOT NULL COMMENT '雇佣信息',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserFindInfo() string {
	return `CREATE TABLE IF NOT EXISTS san_userfindpool (
		  uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		  basefindinfo text NOT NULL COMMENT '基础池信息',
		  rewardinfo text NOT NULL COMMENT '奖励进度信息',
		  wishinfo text NOT NULL COMMENT '许愿阵容',
		  PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 排行任务
func (self *MySqlMgr) UserRankTask() string {
	return `CREATE TABLE IF NOT EXISTS san_userranktask (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		taskinfo text NOT NULL COMMENT '玩家任务列表',
		getstate text NOT NULL COMMENT '领取状态',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 排行任务
func (self *MySqlMgr) MgrRankTask() string {
	return `CREATE TABLE IF NOT EXISTS san_mgrranktask (
		id int(11) NOT NULL COMMENT '任务id',
		rankplayerinfo text NOT NULL COMMENT '玩家列表',
		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 共鸣水晶
func (self *MySqlMgr) MgrResonanceCrystal() string {
	return `CREATE TABLE IF NOT EXISTS san_resonancecrystal (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		priestsheros text NOT NULL COMMENT '祭司英雄',
		resonanceheros text NOT NULL COMMENT '共鸣英雄',
		resonancecount int(11) NOT NULL COMMENT '共鸣个数',
		maxfight bigint(20) NOT NULL COMMENT '最高战力',
		maxfighttime bigint(20) NOT NULL COMMENT '最高战力时间',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 王座之塔
func (self *MySqlMgr) KingTower() string {
	return `CREATE TABLE IF NOT EXISTS san_kingtower (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		info text NOT NULL COMMENT '祭司英雄',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 王座之塔
func (self *MySqlMgr) KingTowerMgr() string {
	return `CREATE TABLE IF NOT EXISTS san_towerplayerrecord (
		keyid bigint(20) NOT NULL COMMENT 'key值',		
		name text NOT NULL COMMENT '名字',
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		icon int(11) NOT NULL COMMENT '头像',
		level int(11) NOT NULL COMMENT '等级',
		playerfight bigint(20) NOT NULL COMMENT '玩家战力',
		battlefight bigint(20) NOT NULL COMMENT '战斗参与的战力',
		Time bigint(20) NOT NULL COMMENT '时间',
		PRIMARY KEY (keyid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) ArenaRank1() string {
	return `CREATE TABLE IF NOT EXISTS san_rankarena1 (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		point bigint(20) NOT NULL COMMENT '积分',
		rank int(11) NOT NULL COMMENT '排名',
		starttime bigint(20) NOT NULL COMMENT '到达时间',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}
func (self *MySqlMgr) ArenaRank2() string {
	return `CREATE TABLE IF NOT EXISTS san_rankarena2 (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		class int(11) NOT NULL COMMENT '层',
		dan int(11) NOT NULL COMMENT '段位',
		rank int(11) NOT NULL COMMENT '排名',
		point bigint(20) NOT NULL COMMENT '积分',
		starttime bigint(20) NOT NULL COMMENT '到达时间',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}
func (self *MySqlMgr) ArenaRank3() string {
	return `CREATE TABLE IF NOT EXISTS san_rankarena3 (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		point bigint(20) NOT NULL COMMENT '积分',
		rank int(11) NOT NULL COMMENT '排名',
		starttime bigint(20) NOT NULL COMMENT '到达时间',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) ArenaPlayer1() string {
	return `CREATE TABLE IF NOT EXISTS san_playerarena1 (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		rank int(11) NOT NULL COMMENT '排名',
		point bigint(20) NOT NULL COMMENT '积分',
		name text NOT NULL COMMENT '名字',
		format text NOT NULL COMMENT '阵容',
		count int(11) NOT NULL COMMENT '次数',
		randnum bigint(20) NOT NULL COMMENT '随机数',
		arenafight text NOT NULL COMMENT '战报',			
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}
func (self *MySqlMgr) ArenaPlayer2() string {
	return `CREATE TABLE IF NOT EXISTS san_playerarena2 (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		class int(11) NOT NULL COMMENT '阶',
		dan int(11) NOT NULL COMMENT '段位',
		coin int(11) NOT NULL COMMENT '硬币',
		point bigint(20) NOT NULL COMMENT '积分',
		name text NOT NULL COMMENT '名字',
		format text NOT NULL COMMENT '阵容',
		count int(11) NOT NULL COMMENT '次数',
		buycount int(11) NOT NULL COMMENT '购买次数',
		arenafight text NOT NULL COMMENT '战报',	
		state int(11) NOT NULL COMMENT '状态',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}
func (self *MySqlMgr) ArenaPlayer3() string {
	return `CREATE TABLE IF NOT EXISTS san_playerarena3 (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		rank int(11) NOT NULL COMMENT '排名',
		point bigint(20) NOT NULL COMMENT '积分',
		name text NOT NULL COMMENT '名字',
		format text NOT NULL COMMENT '阵容',
		count int(11) NOT NULL COMMENT '次数',
		randnum bigint(20) NOT NULL COMMENT '随机数',
		arenafight text NOT NULL COMMENT '战报',	
		arenabattleinfo text NOT NULL COMMENT '战报',	
		arenabattlerecord text NOT NULL COMMENT '战报',	
		enemy text NOT NULL COMMENT '敌人',	
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 活动礼包
func (self *MySqlMgr) ActivityGift() string {
	return `CREATE TABLE IF NOT EXISTS san_activitygift (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		info text NOT NULL COMMENT '礼包数据',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) GrowthGift() string {
	return `CREATE TABLE IF NOT EXISTS san_growthgift (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		info text NOT NULL COMMENT '礼包数据',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) ArenaTime() string {
	return `CREATE TABLE IF NOT EXISTS san_arenatime (
		type int(11) NOT NULL COMMENT '竞技场类型',
		starttime bigint(20) NOT NULL COMMENT '开始时间',
		refreshtime bigint(20) NOT NULL COMMENT '刷新时间戳',
		PRIMARY KEY (type)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) HeroSkin() string {
	return `CREATE TABLE IF NOT EXISTS san_skin (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		info text NOT NULL COMMENT '皮肤数据',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) SpecialPurchase() string {
	return `CREATE TABLE IF NOT EXISTS san_specialpurchase (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		info text NOT NULL COMMENT '领取信息',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}
func (self *MySqlMgr) LifeTree() string {
	return `CREATE TABLE IF NOT EXISTS san_lifetree (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		info text NOT NULL COMMENT '职业分支',
		award text NOT NULL COMMENT '奖励记录',
		heroget text NOT NULL COMMENT '获取记录',
		mainlevel int(11) NOT NULL COMMENT '主等级',
		isget int(11) NOT NULL COMMENT '是否获得初次结算',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 目标任务 试炼任务
func (self *MySqlMgr) UserTargetTask() string {
	return `CREATE TABLE IF NOT EXISTS san_usertargettask (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		taskinfo text NOT NULL COMMENT '玩家任务列表',
		systeminfo text NOT NULL COMMENT '任务组标记',
		nobilitylevel int(11) NOT NULL COMMENT '爵位等级',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// VIP每日充值
func (self *MySqlMgr) UserVipRecharge() string {
	return `CREATE TABLE IF NOT EXISTS san_userviprecharge (
		uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
		taskinfo text NOT NULL COMMENT '玩家任务列表',
		nexttime int(11) NOT NULL COMMENT '刷新时间',
		PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserTurnTable() string {
	return `CREATE TABLE IF NOT EXISTS san_userturntable (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
turntableinfo text NOT NULL COMMENT '转盘信息',
nowstage int(11) NOT NULL COMMENT '当前阶段',
nowcount int(11) NOT NULL COMMENT '当前次数',
nexttime int(11) NOT NULL COMMENT '下次时间',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserAccessCard() string {
	return `CREATE TABLE IF NOT EXISTS san_useraccesscard (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
taskinfo text NOT NULL COMMENT '任务信息',
awardinfo text NOT NULL COMMENT '积分奖励信息',
n3 int(11) NOT NULL COMMENT '期',
ngroup int(11) NOT NULL COMMENT '组',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserFund() string {
	return `CREATE TABLE IF NOT EXISTS san_userfund(
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
n3 int(11) NOT NULL COMMENT '期数',
ngroup int(11) NOT NULL COMMENT '组',
taskinfo text NOT NULL COMMENT '任务信息',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserInterStellar() string {
	return `CREATE TABLE IF NOT EXISTS san_interstellar(
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
galaxyinfo text NOT NULL COMMENT '星系',
privilegevalue text NOT NULL COMMENT '特权值',
stellarcount int(11) NOT NULL COMMENT '组',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserBase2() string {
	return `CREATE TABLE IF NOT EXISTS san_userbase2 (
  id int(11) NOT NULL AUTO_INCREMENT,
  uid bigint(20) NOT NULL COMMENT 'uid',
  uname varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT '用户名称',
  iconid int(5) NOT NULL DEFAULT '1002' COMMENT '图标',
  gem int(11) NOT NULL COMMENT '钻石',
  gold int(11) NOT NULL COMMENT '金币',
  exp int(11) NOT NULL COMMENT '经验',
  level int(11) NOT NULL COMMENT '等级',
  regtime timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '注册时间',
  lastlogintime timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '最后登录时间',
  lastupdtime timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '最后下线时间',
  lastlivetime timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '刷新次数',
  face int(10) NOT NULL DEFAULT '0',
  vip int(10) NOT NULL DEFAULT '0' COMMENT 'vip',
  vipexp int(10) NOT NULL DEFAULT '0',
  tili int(10) NOT NULL DEFAULT '0' COMMENT '体力',
  partyid int(10) NOT NULL DEFAULT '0' COMMENT '公会',
  skillpoint int(10) NOT NULL DEFAULT '10',
  tililastupdatatime int(10) NOT NULL DEFAULT '0' COMMENT '体力恢复时间',
  splastupdatatime int(10) NOT NULL DEFAULT '0',
  lastcheckintime timestamp NOT NULL DEFAULT '1971-01-02 00:00:00' COMMENT '最后签到时间',
  checkinnum int(10) NOT NULL DEFAULT '0' COMMENT '签到次数',
  checkinaward int(11) NOT NULL DEFAULT '0' COMMENT '签到奖励',
  position int(11) NOT NULL DEFAULT '10010',
  isrename int(2) NOT NULL DEFAULT '0',
  logindays int(11) NOT NULL DEFAULT '1' COMMENT '登录天数',
  loginaward int(11) NOT NULL DEFAULT '0' COMMENT '登录奖励',
  levelaward int(11) NOT NULL DEFAULT '0' COMMENT '等级奖励',
  worldaward int(10) NOT NULL DEFAULT '0',
  morale int(11) NOT NULL DEFAULT '0',
  citylevel int(10) NOT NULL DEFAULT '0' COMMENT '城市等级',
  camp int(10) NOT NULL DEFAULT '0',
  city int(11) NOT NULL DEFAULT '0' COMMENT '城市位置',
  fight int(11) NOT NULL DEFAULT '0' COMMENT '战力',
  isgag int(11) NOT NULL DEFAULT '0' COMMENT '是否禁言',
  isblock int(11) NOT NULL DEFAULT '0',
  blockday int(11) NOT NULL,
  ip varchar(64) NOT NULL DEFAULT '0' COMMENT '上次登录ip',
  promotebox int(11) NOT NULL DEFAULT '0',
  linetime bigint(20) NOT NULL DEFAULT '0',
  passmax int(11) NOT NULL DEFAULT '0',
  fitserver int(11) NOT NULL DEFAULT '0',
  blocktime int(11) NOT NULL DEFAULT '0',
  blockreason varchar(63) NOT NULL DEFAULT '',
  soul int(11) NOT NULL DEFAULT '0' COMMENT '魂石',
  techpoint int(11) NOT NULL DEFAULT '0' COMMENT '科技点',
  bossmoney int(11) NOT NULL DEFAULT '0' COMMENT '巨兽水晶',
  towerstone int(11) NOT NULL DEFAULT '0' COMMENT '镇魂石',
  portrait int(11) NOT NULL DEFAULT '0' COMMENT '头像挂件',
  campok int(11) NOT NULL DEFAULT '0' COMMENT '阵营ok',
  nameok int(11) NOT NULL DEFAULT '0' COMMENT '头像ok',
  guildid int(11) NOT NULL DEFAULT '0' COMMENT '指引Id',
  redicon int(11) NOT NULL DEFAULT '0' COMMENT '红色图标',
  usersignature text NOT NULL COMMENT '用户签名',
  getallgem int(11) NOT NULL DEFAULT '0' COMMENT '历史总钻石数',
  paygem int(11) NOT NULL DEFAULT '0' COMMENT '付费钻石',
  PRIMARY KEY (id)
)  ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) ActivityBoss() string {
	return `CREATE TABLE IF NOT EXISTS san_activityboss (
id  int(11) NOT NULL COMMENT '活动ID',
period int(11) NOT NULL COMMENT '活动期数',
activitybosstop text NOT NULL COMMENT '记录',
starttime int(11) NOT NULL COMMENT '活动开始时间',
endtime int(11) NOT NULL COMMENT '活动结束时间',
rewardtime int(11) NOT NULL COMMENT '发奖时间',
hasreward int(11) NOT NULL COMMENT '是否发奖',
PRIMARY KEY (id,period)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserActivityBoss() string {
	return `CREATE TABLE IF NOT EXISTS san_useractivityboss (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
activitybossinfo text NOT NULL COMMENT '任务信息',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserHeroGrow() string {
	return `CREATE TABLE IF NOT EXISTS san_userherogrow (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
herogrowinfos text NOT NULL COMMENT '礼包集合',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserCrossArena() string {
	return `CREATE TABLE IF NOT EXISTS san_usercrossarena (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
keyid int(11) NOT NULL COMMENT '活动key',
subsection int(11) NOT NULL COMMENT '当前大段位',
class int(11) NOT NULL COMMENT '当前小段位',
subsectionmax int(11)NOT NULL COMMENT '最高大段位',
classmax int(11) NOT NULL COMMENT '最高小段位',
times int(11) NOT NULL COMMENT '挑战次数',
buytimes int(11) NOT NULL COMMENT '购买挑战次数',
starttime int(11) NOT NULL COMMENT '活动开始时间',
endtime int(11) NOT NULL COMMENT '活动结束时间',
taskawardsign text NOT NULL COMMENT '任务领取标记',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) CrossArena() string {
	return `CREATE TABLE IF NOT EXISTS san_crossarena (
keyid  int(11) NOT NULL COMMENT '活动ID期数',
starttime  bigint(20) NOT NULL COMMENT '开始时间',
endtime  bigint(20) NOT NULL COMMENT '结束时间',
isreward  int(11) NOT NULL COMMENT '是否发奖',
updatetime  int(11) NOT NULL COMMENT '请求时间',
PRIMARY KEY (keyid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) BattleRecord() string {
	return `CREATE TABLE IF NOT EXISTS san_battlerecord(
		id bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '唯一Id',
  		fightid bigint(20) NULL DEFAULT NULL COMMENT '战斗ID',
  		recordtype int(11) NOT NULL DEFAULT 0 COMMENT '战报类型',
		battleinfo mediumtext NOT NULL COMMENT '简报',
		battlerecord mediumtext NOT NULL COMMENT '详细战报',
		updatetime int(11) NOT NULL DEFAULT 0 COMMENT '时间',
		PRIMARY KEY (id) USING BTREE,
  		INDEX FIGHT_ID(fightid) USING BTREE COMMENT '战斗ID'
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserActivityBossFestival() string {
	return `CREATE TABLE IF NOT EXISTS san_useractivitybossfestival (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
activitybossfestivalinfo text NOT NULL COMMENT '信息',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

// 离线全局战报管理
func (self *MySqlMgr) OfflineRecordInfo() string {
	return `CREATE TABLE IF NOT EXISTS san_offlinerecordinfo (
		  id int(11) NOT NULL COMMENT 'id',
		  period int(11) NOT NULL DEFAULT 0 COMMENT '期数',
		  records text NOT NULL COMMENT '战报集合',
		  PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserLotteryDraw() string {
	return `CREATE TABLE IF NOT EXISTS san_userlotterydraw (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
keyid int(11) NOT NULL DEFAULT 0 COMMENT '期数',
lotterydrawinfo text NOT NULL COMMENT '转盘信息',
alreadyget text NOT NULL COMMENT '大奖记录',
nowstage int(11) NOT NULL DEFAULT 0  COMMENT '当前阶段',
nowcount int(11) NOT NULL DEFAULT 0  COMMENT '当前次数',
lowchoose int(11) NOT NULL DEFAULT 0  COMMENT '普通大奖选择',
highchoose int(11) NOT NULL DEFAULT 0  COMMENT '终极大奖选择',
luckvalue int(11) NOT NULL DEFAULT 0  COMMENT '幸运值',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) HonourShop() string {
	return `CREATE TABLE IF NOT EXISTS san_honourshop (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
shopgood text NOT NULL COMMENT '商品信息',
nextdaytime int(11) NOT NULL DEFAULT 0  COMMENT '下次日刷新时间',
nextweektime int(11) NOT NULL DEFAULT 0  COMMENT '下次周刷新时间',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) UserCrossArena3V3() string {
	return `CREATE TABLE IF NOT EXISTS san_usercrossarena3v3 (
uid bigint(20) NOT NULL COMMENT '玩家唯一Id',
keyid int(11) NOT NULL COMMENT '活动key',
subsection int(11) NOT NULL COMMENT '当前大段位',
class int(11) NOT NULL COMMENT '当前小段位',
subsectionmax int(11)NOT NULL COMMENT '最高大段位',
classmax int(11) NOT NULL COMMENT '最高小段位',
times int(11) NOT NULL COMMENT '挑战次数',
buytimes int(11) NOT NULL COMMENT '购买挑战次数',
starttime bigint(20) NOT NULL COMMENT '活动开始时间',
endtime bigint(20) NOT NULL COMMENT '活动结束时间',
showtime bigint(20)NOT NULL COMMENT '活动展示时间',
taskawardsign text NOT NULL COMMENT '任务领取标记',
PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) CrossArena3V3() string {
	return `CREATE TABLE IF NOT EXISTS san_crossarena3v3 (
keyid  int(11) NOT NULL COMMENT '活动ID期数',
starttime  bigint(20) NOT NULL COMMENT '开始时间',
endtime  bigint(20) NOT NULL COMMENT '结束时间',
showtime bigint(20)NOT NULL COMMENT '活动展示时间',
isreward  int(11) NOT NULL COMMENT '是否发奖',
updatetime  int(11) NOT NULL COMMENT '请求时间',
PRIMARY KEY (keyid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *MySqlMgr) RankReward1732() string {
	return `CREATE TABLE IF NOT EXISTS rankreward1732 (
uid  int(11) NOT NULL COMMENT '玩家uid',
keyid int(11) NOT NULL COMMENT '活动keyid',
score bigint(20) NOT NULL COMMENT '分数',
toprank bigint(20) NOT NULL COMMENT '排名',
info text NOT NULL COMMENT '详细信息',
PRIMARY KEY (uid,keyid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}
