package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// levels
const (
	debugLevel = 0
	infoLevel  = 1
	errorLevel = 2
	fatalLevel = 3
)

const (
	printDebugLevel = "[debug ] "
	printInfoLevel  = "[info  ] "
	printErrorLevel = "[error ] "
	printFatalLevel = "[fatal ] "
)

const (
	LOG_FORMAT_DAY          = "20060102"
	ERR_FORMAT_DAY          = "20060102-err"
	KB                int64 = 1024
	MB                int64 = 1024 * 1024
	INFO_LOGGER_TYPE        = 1
	ERROR_LOGGER_TYPE       = 2
)

type Loger struct {
	level      int         //! 输出日志等级
	infoLogger *log.Logger //! 日志组件
	errLogger  *log.Logger //! 异常日志组建
	currentDay int         //! 当前天
	logPath    string      //! 当前日志文件路径
	showStd    bool        //! 是否输出到标准输出
	logFormat  string      //! 普通日志分割[天或者小时]
	errFormat  string      //! 错误日志分割[天或者小时]
}

type LoggerChecker struct {
	LoggerIndex int
	ErrorIndex  int
}

var loggerCheckSingleton *LoggerChecker = nil

var logersingleton *Loger = nil

//! public
func GetLogMgr() *Loger {
	if logersingleton == nil {
		logersingleton = new(Loger)
		logersingleton.logPath = "./log"
		logersingleton.currentDay = -1
		logersingleton.level = debugLevel
		logersingleton.showStd = true
		os.Mkdir(logersingleton.logPath, os.ModePerm)

		logersingleton.changDay()
	}

	return logersingleton
}

//设置log输出等级，默认DebugLevel, true
func (self *Loger) SetLevel(level int, std bool) {
	if level >= debugLevel && level <= fatalLevel {
		self.level = level
	}
	self.showStd = std
}

//判断是否生成新日志文件
func (self *Loger) changDay() { ///跨天则生成新文件
	now := time.Now()
	currentDay := now.Day()
	if self.currentDay == currentDay {
		return
	}

	fileName := self.logPath + "/" + now.Format(LOG_FORMAT_DAY) + "[1].log"
	_, err := os.Stat(fileName)
	if nil == err && self.infoLogger != nil {
		return
	}

	self.currentDay = currentDay
	loggerFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println("open log file fail!", err.Error())
	}
	self.infoLogger = log.New(loggerFile, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)

	errFileName := self.logPath + "/" + now.Format(ERR_FORMAT_DAY) + "[1].log"
	_, err1 := os.Stat(fileName)
	if nil == err1 && self.errLogger != nil {
		return
	}
	errLoggerFile, err := os.OpenFile(errFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println("open log err file fail!", err.Error())
	}
	self.errLogger = log.New(errLoggerFile, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)

}

func (self *Loger) doPrintf(level int, printLevel string, a ...interface{}) {
	if level < self.level {
		return
	}
	if self.infoLogger == nil {
		log.Println("logger closed", fmt.Sprintln(a...))
		//panic("logger closed")
	}

	//self.changDay()

	format := printLevel + "%s"
	if level <= infoLevel {
		self.infoLogger.Output(3, fmt.Sprintf(format, fmt.Sprintln(a...)))
	} else {
		self.errLogger.Output(3, fmt.Sprintf(format, fmt.Sprintln(a...)))
	}

	if self.showStd {
		log.Println(a...)
	}

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (self *Loger) Close() {
	self.infoLogger = nil
	self.errLogger = nil
}

func (self *Loger) Debug(a ...interface{}) {
	self.doPrintf(debugLevel, printDebugLevel, a...)
}

func (self *Loger) Info(a ...interface{}) {
	self.doPrintf(infoLevel, printInfoLevel, a...)
}

func (self *Loger) Error(a ...interface{}) {
	self.doPrintf(errorLevel, printErrorLevel, a...)
}

func (self *Loger) Fatal(a ...interface{}) {
	self.doPrintf(fatalLevel, printFatalLevel, a...)
}

func LogDebug(a ...interface{}) {
	GetLogMgr().doPrintf(debugLevel, printDebugLevel, a...)
}

func LogInfo(a ...interface{}) {
	GetLogMgr().doPrintf(infoLevel, printInfoLevel, a...)
}

func LogError(a ...interface{}) {
	GetLogMgr().doPrintf(errorLevel, printErrorLevel, a...)
}

func LogFatal(a ...interface{}) {
	GetLogMgr().doPrintf(fatalLevel, printFatalLevel, a...)
}

func GetLoggerCheckMgr() *LoggerChecker {
	if loggerCheckSingleton == nil {
		loggerCheckSingleton = new(LoggerChecker)
		loggerCheckSingleton.LoggerIndex = 1
		loggerCheckSingleton.ErrorIndex = 1
	}
	return loggerCheckSingleton
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("logger check err", err)
	}
}

//! 日志检查逻辑
func (self *LoggerChecker) Run() {
	defer catchError()
	ticker := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-ticker.C:
			self.checkFile()
		}
	}
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (self *LoggerChecker) checkFile() {
	// 传入文件格式,文件检查索引,设置*log.logger指针类型,文件大小
	self.checkLogger(LOG_FORMAT_DAY, &self.LoggerIndex, INFO_LOGGER_TYPE)
	self.checkLogger(ERR_FORMAT_DAY, &self.ErrorIndex, ERROR_LOGGER_TYPE)
}

func (self *Loger) setInfoLog(p *log.Logger, loggerType int) {
	if loggerType == 1 {
		self.infoLogger = p
	} else if loggerType == 2 {
		self.errLogger = p
	}
}

// 1.获取当前日期
// 2.遍历检查超过大小的日志
// 3.创建新的日志文件,索引++
func (self *LoggerChecker) checkLogger(format string, pIndex *int, logType int) {
	now := time.Now()
	logPath := "./log"
	fileIndex := 1
	str := now.Format(format)
	// 检查最后一个大小不满的文件
	maxFileNum := 500
	for i := 1; i <= maxFileNum; i++ {
		if i < *pIndex {
			continue
		}
		// 先检查文件是否存在
		loggerFn := logPath + "/" + str + fmt.Sprintf("[%d].log", fileIndex)

		//返回最后一个文件
		if exists, _ := Exists(loggerFn); exists {
			fileIndex += 1
		} else {
			fileIndex -= 1
			break
		}
	}

	if fileIndex <= 0 {
		fileIndex = 1
	}

	lastFn := logPath + "/" + str + fmt.Sprintf("[%d].log", fileIndex)
	// 只读方式打开这个文件
	loggerFile, err := os.OpenFile(lastFn, os.O_RDONLY, os.ModePerm)
	if err != nil {
		// 重新生成!
		loggerFile, err = os.OpenFile(lastFn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Println("open log file fail!", err.Error())
		} else {
			pLogger := log.New(loggerFile, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
			GetLogMgr().setInfoLog(pLogger, logType)
		}

	}

	loggerStat, err := loggerFile.Stat()
	if err != nil {
		log.Println("logger file stat error", err.Error())
		return
	}

	// 检查这个文件的大小
	size := loggerStat.Size()
	// 超过指定大小,则创建一个新的文件,并重新设置对应的logger指针
	maxFileSize := 500 * MB
	if size >= maxFileSize {
		fileIndex++
		nextFileName := logPath + "/" + str + fmt.Sprintf("[%d].log", fileIndex)
		*pIndex = fileIndex
		nextFile, err := os.OpenFile(nextFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Println("open log file fail!", err.Error())
			return
		}
		pLogger := log.New(nextFile, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
		GetLogMgr().setInfoLog(pLogger, logType)
	}
}
