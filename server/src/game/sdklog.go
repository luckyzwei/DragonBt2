package game

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	LOG_NAME_GOLD     = "money_gold_user"
	LOG_NAME_ITEM     = "knapsack_user"
	LOG_NAME_OFFLINE  = "offline_user"
	LOG_NAME_RECHARGE = "business_user"
	LOG_NAME_SURPLUS  = "stock_gold_user"
	LOG_NAME_SCORE    = "score_data_user"
)

const (
	LOG_TYPE_GOLD     = 0
	LOG_TYPE_ITEM     = 1
	LOG_TYPE_OFFLINE  = 2
	LOG_TYPE_RECHARGE = 3
	LOG_TYPE_SURPLUS  = 4
	LOG_TYPE_SCORE    = 5
	LOG_TYPE_MAX      = 6
)

var logname = [LOG_TYPE_MAX]string{"money_gold_user", "knapsack_user", "offline_user", "business_user", "stock_gold_user", "score_data_user"}

const (
	LOG_SDK_DAY  = "2006-01-02"
	LOG_SDK_TIME = "2006-01-02 15:04:05"
)

type SdkLoger struct {
	logger     [LOG_TYPE_MAX]*log.Logger //! 日志组件
	currentDay int                       //! 当前天
	logPath    string                    //! 当前日志文件路径
	showStd    bool                      //! 是否输出到标准输出
}

type SdkLoggerChecker struct {
	LoggerIndex [LOG_TYPE_MAX]int
}

var sdkloggerCheckSingleton *SdkLoggerChecker = nil

var sdklogersingleton *SdkLoger = nil

//! public
func GetSdkLogMgr() *SdkLoger {
	if sdklogersingleton == nil {
		sdklogersingleton = new(SdkLoger)
		sdklogersingleton.logPath = "../log/" + TimeServer().Format(LOG_SDK_DAY)
		sdklogersingleton.currentDay = -1
		sdklogersingleton.showStd = true
		os.Mkdir(sdklogersingleton.logPath, os.ModePerm)

		sdklogersingleton.changeDay()
	}

	return sdklogersingleton
}

//判断是否生成新日志文件
func (self *SdkLoger) changeDay() { ///跨天则生成新文件
	now := TimeServer()
	currentDay := now.Day()
	if self.currentDay == currentDay {
		return
	}

	self.logPath = "../log/" + TimeServer().Format(LOG_SDK_DAY)

	GetSdkLoggerCheckMgr().CheckPath(self.logPath)

	for i, v := range logname {
		fileName := self.logPath + "/" + GetServer().Con.GameID + "/" + HF_ItoA(GetServer().Con.ServerId) + "/" + v + "_" + GetServer().Con.GameName + ".txt"
		_, err1 := os.Stat(fileName)
		if nil == err1 && self.logger[i] != nil {
			continue
		}
		LoggerFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Println("open log err file fail!", err.Error())
		}
		self.logger[i] = log.New(LoggerFile, "", 0)
	}

	self.currentDay = currentDay

}

func (self *SdkLoger) doPrintf(logType int, content string) {
	if logType < LOG_TYPE_GOLD || logType >= LOG_TYPE_MAX {
		log.Println("type errer", content)
		return
	}
	if self.logger[logType] == nil {
		log.Println("logger closed", content)
		//panic("logger closed")
	}

	//self.changDay()

	self.logger[logType].Output(3, content)
}

func (self *SdkLoger) Close() {
	for i, _ := range self.logger {
		self.logger[i] = nil
	}
}

func AddSdkOfflineLog(player *Player) {
	if !GetServer().Con.LogCon.SDK {
		return
	}

	tll, _ := time.ParseInLocation(DATEFORMAT, player.Sql_UserBase.LastLoginTime, time.Local)
	onlinetime := TimeServer().Unix() - tll.Unix()

	ip := ""
	if player.GetSession() != nil {
		ip = player.GetSession().IP
	}

	var offlinecontent string
	offlinecontent = "[" + TimeServer().Format(LOG_SDK_TIME) + "]" + "[" +
		"op_type:OffLine " +
		"p_id:" + GetServer().Con.GameID + " " +
		"ditch_id:" + player.Account.Channelid + " " +
		"user:" + player.Account.UserId + " " +
		"role_id:" + HF_I64toA(player.GetUid()) + " " +
		"ip:(" + ip + ") " +
		"device:(" + player.Platform.Platform + ") " +
		"mac:(" + player.Platform.Mac + ") " +
		"device_type:(" + player.Platform.Brand + player.Platform.Model + ") " +
		"online_time:" + HF_I64toA(onlinetime) + " " +
		"accu_online_time:" + HF_I64toA(player.Sql_UserBase.LineTime) + " " +
		"login_time:" + HF_I64toA(tll.Unix()) + " " +
		"offline_time:" + HF_I64toA(TimeServer().Unix()) + " " +
		"role_level:" + HF_ItoA(player.GetLv()) + " " +
		"original_uid:" + "[" + GetServer().Con.GameID + " " + HF_ItoA(GetServer().Con.ServerId) + " " + HF_I64toA(player.GetUid()) + " " + player.Account.UserId + "]" +
		"]"
	GetSdkLogMgr().doPrintf(LOG_TYPE_OFFLINE, offlinecontent)

	var stockcontent string
	stockcontent = "[" + TimeServer().Format(LOG_SDK_TIME) + "]" + "[" +
		"op_type:RoleGold " +
		"p_id:" + GetServer().Con.GameID + " " +
		"ditch_id:" + player.Account.Channelid + " " +
		"user:" + player.Account.UserId + " " +
		"role_id:" + HF_I64toA(player.GetUid()) + " " +
		"original_uid:" + "[" + GetServer().Con.GameID + " " + HF_ItoA(GetServer().Con.ServerId) + " " + HF_I64toA(player.GetUid()) + " " + player.Account.UserId + "]" + " " +
		"charge_time:" + HF_I64toA(player.GetModule("recharge").(*ModRecharge).GetLastRechargeTime()) + " " +
		"remain_gold:" + HF_ItoA(player.Sql_UserBase.Gem) + " " +
		"history_gold:" + HF_ItoA(player.Sql_UserBase.GetAllGem) + " " +
		"role_level:" + HF_ItoA(player.GetLv()) + " " +
		"]"
	GetSdkLogMgr().doPrintf(LOG_TYPE_SURPLUS, stockcontent)
}

// 当获得物品时 param3 为-1 则不记录 由获得物品的最外层进行记录物品列表 大于等于0 则正常记录
func AddSdkItemLog(uid int64, _type int, value int, param1 int, param2 int, dec string, cur int, param3 int, player *Player) {
	if player == nil {
		return
	}

	if !GetServer().Con.LogCon.SDK {
		return
	}

	var content string
	// 得失钻石时
	if _type == DEFAULT_GEM {
		if value > 0 {
			content = "[" + TimeServer().Format(LOG_SDK_TIME) + "]" + "[" +
				"op_type:AddGold " +
				"p_id:" + GetServer().Con.GameID + " " +
				"ditch_id:" + player.Account.Channelid + " " +
				"user:" + player.Account.UserId + " " +
				"role_id:" + HF_I64toA(player.GetUid()) + " " +
				"type:" + dec + " " +
				"get_item:" + HF_ItoA(0) + " " +
				"add_gold:" + HF_ItoA(value) + " " +
				"use_gold:" + HF_ItoA(0) + " " +
				"remain_gold:" + HF_ItoA(cur) + " " +
				"role_level:" + HF_ItoA(player.GetLv()) + " " +
				"original_uid:" + "[" + GetServer().Con.GameID + " " + HF_ItoA(GetServer().Con.ServerId) + " " + HF_I64toA(player.GetUid()) + " " + player.Account.UserId + "]" +
				"]"
			GetSdkLogMgr().doPrintf(LOG_TYPE_GOLD, content)
		} else {
			content = "[" + TimeServer().Format(LOG_SDK_TIME) + "]" + "[" +
				"op_type:UseGold " +
				"p_id:" + GetServer().Con.GameID + " " +
				"ditch_id:" + player.Account.Channelid + " " +
				"user:" + player.Account.UserId + " " +
				"role_id:" + HF_I64toA(player.GetUid()) + " " +
				"type:" + dec + " " +
				"get_item:" + HF_ItoA(param3) + " " +
				"add_gold:" + HF_ItoA(0) + " " +
				"use_gold:" + HF_ItoA(-value) + " " +
				"remain_gold:" + HF_ItoA(cur) + " " +
				"role_level:" + HF_ItoA(player.GetLv()) + " " +
				"original_uid:" + "[" + GetServer().Con.GameID + " " + HF_ItoA(GetServer().Con.ServerId) + " " + HF_I64toA(player.GetUid()) + " " + player.Account.UserId + "]" +
				"]"
			GetSdkLogMgr().doPrintf(LOG_TYPE_GOLD, content)
		}
	} else { // 得失物品时
		if value > 0 {
			if param3 > 0 {
				content = "[" + TimeServer().Format(LOG_SDK_TIME) + "]" + "[" +
					"op_type:PutItem " +
					"p_id:" + GetServer().Con.GameID + " " +
					"ditch_id:" + player.Account.Channelid + " " +
					"user:" + player.Account.UserId + " " +
					"role_id:" + HF_I64toA(player.GetUid()) + " " +
					"role_level:" + HF_ItoA(player.GetLv()) + " " +
					"original_uid:" + "[" + GetServer().Con.GameID + " " + HF_ItoA(GetServer().Con.ServerId) + " " + HF_I64toA(player.GetUid()) + " " + player.Account.UserId + "]" + " " +
					"type:" + dec + " " +
					"use_gold:" + HF_ItoA(param3) + " " +
					"item_list:" + "[" + HF_ItoA(_type) + "," + HF_ItoA(value) + "," + HF_ItoA(0) + "]" +
					"]"
				GetSdkLogMgr().doPrintf(LOG_TYPE_ITEM, content)
			}
		} else {
			content = "[" + TimeServer().Format(LOG_SDK_TIME) + "]" + "[" +
				"op_type:UseItem " +
				"p_id:" + GetServer().Con.GameID + " " +
				"ditch_id:" + player.Account.Channelid + " " +
				"user:" + player.Account.UserId + " " +
				"role_id:" + HF_I64toA(player.GetUid()) + " " +
				"role_level:" + HF_ItoA(player.GetLv()) + " " +
				"original_uid:" + "[" + GetServer().Con.GameID + " " + HF_ItoA(GetServer().Con.ServerId) + " " + HF_I64toA(player.GetUid()) + " " + player.Account.UserId + "]" + " " +
				"type:" + dec + " " +
				"use_gold:" + HF_ItoA(0) + " " +
				"item_id:" + HF_ItoA(_type) + " " +
				"item_num:" + HF_ItoA(-value) +
				"]"
			GetSdkLogMgr().doPrintf(LOG_TYPE_ITEM, content)
		}

	}
}

func AddSpecialSdkItemListLog(player *Player, num int, item []PassItem, dec string) {
	if !GetServer().Con.LogCon.SDK {
		return
	}
	content := "[" + TimeServer().Format(LOG_SDK_TIME) + "]" + "[" +
		"op_type:PutItem " +
		"p_id:" + GetServer().Con.GameID + " " +
		"ditch_id:" + player.Account.Channelid + " " +
		"user:" + player.Account.UserId + " " +
		"role_id:" + HF_I64toA(player.GetUid()) + " " +
		"role_level:" + HF_ItoA(player.GetLv()) + " " +
		"original_uid:" + "[" + GetServer().Con.GameID + " " + HF_ItoA(GetServer().Con.ServerId) + " " + HF_I64toA(player.GetUid()) + " " + player.Account.UserId + "]" + " " +
		"type:" + dec + " " +
		"use_gold:" + HF_ItoA(num) + " " +
		"item_list:" + "[" + ItemToContent(item) + "]" +
		"]"
	GetSdkLogMgr().doPrintf(LOG_TYPE_ITEM, content)

}

func ItemToContent(item []PassItem) string {
	content := ""
	len := len(item)
	for i, v := range item {
		content += "[" + HF_ItoA(v.ItemID) + "," + HF_ItoA(v.Num) + "," + HF_ItoA(0) + "]"
		if i != len-1 {
			content += ","
		}
	}
	return content
}

func GetGemNum(spenditem []PassItem) int {
	for _, v := range spenditem {
		if v.ItemID == DEFAULT_GEM && v.Num < 0 {
			return -v.Num
		}
	}
	return 0
}

//func CheckAddItemLog(player *Player, dec string, spenditem []PassItem, additem []PassItem) {
//	for _, v := range spenditem {
//		if v.ItemID == DEFAULT_GEM && v.Num < 0 && len(additem) > 0 {
//			content := "[" + TimeServer().Format(LOG_SDK_TIME) + "]" + "[" +
//				"op_type:UseGold " +
//				"p_id:" + HF_ItoA(10000) + " " +
//				"ditch_id:" + player.Account.Channelid + " " +
//				"user:" + player.Account.Account + " " +
//				"role_id:" + HF_I64toA(player.GetUid()) + " " +
//				"type:" + dec + " " +
//				"get_item:" + HF_ItoA(1) + " " +
//				"add_gold:" + HF_ItoA(0) + " " +
//				"use_gold:" + HF_ItoA(-v.Num) + " " +
//				"remain_gold:" + HF_ItoA(player.Sql_UserBase.Gem) + " " +
//				"role_level:" + HF_ItoA(player.GetLv()) + " " +
//				"original_uid:" + "[" + HF_ItoA(10000) + " " + HF_ItoA(GetServer().Con.ServerId) + " " + HF_I64toA(player.GetUid()) + " " + player.Account.Account + "]" +
//				"]"
//			GetSdkLogMgr().doPrintf(LOG_TYPE_GOLD, content)
//		}
//	}
//}

func GetSdkLoggerCheckMgr() *SdkLoggerChecker {
	if sdkloggerCheckSingleton == nil {
		sdkloggerCheckSingleton = new(SdkLoggerChecker)
		sdkloggerCheckSingleton.LoggerIndex = [LOG_TYPE_MAX]int{1, 1, 1, 1, 1, 1}
	}
	return sdkloggerCheckSingleton
}

//! 日志检查逻辑
func (self *SdkLoggerChecker) Run() {
	defer catchError()
	ticker := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-ticker.C:
			self.CheckFile()
		}
	}
}

func (self *SdkLoggerChecker) CheckFile() {
	// 传入文件格式,文件检查索引,设置*log.logger指针类型,文件大小
	GetSdkLogMgr().changeDay()
	//for i, v := range self.LoggerIndex {
	//	self.CheckLogger(LOG_SDK_DAY, &v, i)
	//}
}

func (self *SdkLoger) setInfoLog(p *log.Logger, loggerType int) {
	self.logger[loggerType] = p
}

func (self *SdkLoggerChecker) CheckPath(logPath string) {
	exist, err := Exists(logPath)
	if err != nil {
		fmt.Printf("get logPath error![%v]\n", err)
		return
	}

	if exist {
		fmt.Printf("has logPath![%v]\n", logPath)
	} else {
		fmt.Printf("no logPath![%v]\n", logPath)
		// 创建文件夹
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			fmt.Printf("logPath failed![%v]\n", err)
		} else {
			fmt.Printf("logPath success!\n")
		}
	}

	second := logPath + "/" + GetServer().Con.GameID
	exist, err = Exists(second)
	if err != nil {
		fmt.Printf("get second error![%v]\n", err)
		return
	}

	if exist {
		fmt.Printf("has second![%v]\n", second)
	} else {
		fmt.Printf("no second![%v]\n", second)
		// 创建文件夹
		err := os.Mkdir(second, os.ModePerm)
		if err != nil {
			fmt.Printf("second failed![%v]\n", err)
		} else {
			fmt.Printf("second success!\n")
		}
	}

	third := second + "/" + HF_ItoA(GetServer().Con.ServerId)
	exist, err = Exists(third)
	if err != nil {
		fmt.Printf("get third error![%v]\n", err)
		return
	}

	if exist {
		fmt.Printf("has third![%v]\n", third)
	} else {
		fmt.Printf("no third![%v]\n", third)
		// 创建文件夹
		err := os.Mkdir(third, os.ModePerm)
		if err != nil {
			fmt.Printf("third failed![%v]\n", err)
		} else {
			fmt.Printf("third success!\n")
		}
	}
}

// 1.获取当前日期
// 2.遍历检查超过大小的日志
// 3.创建新的日志文件,索引++
func (self *SdkLoggerChecker) CheckLogger(format string, pIndex *int, logType int) {
	//now := TimeServer()
	//logPath := "./log/" + now.Format(format)
	//fileIndex := 1
	//str := logname[logType]

	//
	//// 检查最后一个大小不满的文件
	//maxFileNum := GetServer().Con.LogCon.MaxFileNum
	//for i := 1; i <= maxFileNum; i++ {
	//	if i < *pIndex {
	//		continue
	//	}
	//	// 先检查文件是否存在
	//	loggerFn := logPath + "/" + "s001" + "/" + HF_ItoA(GetServer().Con.ServerId) + "/" + str + ".txt"
	//
	//	//返回最后一个文件
	//	if exists, _ := Exists(loggerFn); exists {
	//		fileIndex += 1
	//	} else {
	//		fileIndex -= 1
	//		break
	//	}
	//}
	//
	//if fileIndex <= 0 {
	//	fileIndex = 1
	//}
	//
	//lastFn := logPath + "/" + "001" + "/" + HF_ItoA(GetServer().Con.ServerId) + "/" + str + fmt.Sprintf("[%d].txt", fileIndex)
	//// 只读方式打开这个文件
	//loggerFile, err := os.OpenFile(lastFn, os.O_RDONLY, os.ModePerm)
	//if err != nil {
	//	// 重新生成!
	//	loggerFile, err = os.OpenFile(lastFn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	//	if err != nil {
	//		log.Println("open log file fail!", err.Error())
	//	} else {
	//		pLogger := log.New(loggerFile, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	//		GetSdkLogMgr().setInfoLog(pLogger, logType)
	//	}
	//
	//}
	//
	//loggerStat, err := loggerFile.Stat()
	//if err != nil {
	//	log.Println("logger file stat error", err.Error())
	//	return
	//}
	//
	//// 检查这个文件的大小
	//size := loggerStat.Size()
	//// 超过指定大小,则创建一个新的文件,并重新设置对应的logger指针
	//maxFileSize := GetServer().Con.LogCon.MaxFileSize * MB
	//if size >= maxFileSize {
	//	fileIndex++
	//	nextFileName := logPath + "/" + "001" + "/" + string(GetServer().Con.ServerId) + "/" + str + fmt.Sprintf("[%d].txt", fileIndex)
	//	*pIndex = fileIndex
	//	nextFile, err := os.OpenFile(nextFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	//	if err != nil {
	//		log.Println("open log file fail!", err.Error())
	//		return
	//	}
	//	pLogger := log.New(nextFile, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	//	GetSdkLogMgr().setInfoLog(pLogger, logType)
	//}
}
