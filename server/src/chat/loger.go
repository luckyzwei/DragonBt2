package main

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

type Loger struct {
	level      int         //! 输出日志等级
	infoLogger *log.Logger //! 日志组件
	errLogger  *log.Logger //! 异常日志组建
	currentDay int         //! 当前天
	logPath    string      //! 当前日志文件路径
	showStd    bool        //! 是否输出到标准输出
}

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
	fileName := self.logPath + "/" + now.Format("20060102.log")
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

	errFileName := self.logPath + "/" + now.Format("20060102-err.log")
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

	self.changDay()

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
