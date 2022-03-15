package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

//! 注册信号量
func handleSignal(signalType os.Signal, handleFun func(*chan os.Signal)) {
	ch := make(chan os.Signal)
	signal.Notify(ch, signalType)
	go handleFun(&ch)
}

//! 管道破裂
func handlePIPE(ch *chan os.Signal) {
	for {
		<-*ch
		log.Println("get a SIGPIPE")
	}
}

//! ctrl+z
func handleTSTP(ch *chan os.Signal) {
	for {
		<-*ch
		log.Println("get a SIGTSTP")
	}
}

//! gdb trap
func handleTRAP(ch *chan os.Signal) {
	for {
		<-*ch
		log.Println("get a SIGTRAP")
	}
}

//! ctrl+c
func handleINT(ch *chan os.Signal) {
	for {
		<-*ch
		log.Println("get a SIGINT")
		log.Fatal("shotdown")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	//! 注册信号量
	handleSignal(syscall.SIGPIPE, handlePIPE)
	handleSignal(syscall.SIGTRAP, handleTRAP)
	handleSignal(syscall.SIGINT, handleINT)

	// 配置读取
	GetCacheServer().InitConfig()

	// 消息组初始化
	GetCacheServer().startMsgGroup(GetCacheServer().Con.InitGroupNum)

	// 定时器goroutine
	go GetCacheServer().Run()
	go GetLoggerCheckMgr().Run()

	//go addFailLog()

	// 启动http服务器
	http.Handle("/", GetCacheServer().GetConnectHandler())
	log.Println("绑定ip:", GetCacheServer().Con.Host)
	log.Fatal(http.ListenAndServe(GetCacheServer().Con.Host, nil))
}

//func addFailLog()  {
//	logFile, err := os.OpenFile("./log/fatal.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
//	if err != nil {
//		log.Println("服务启动出错",  "打开异常日志文件失败" , err)
//		return
//	}
//	// 将进程标准出错重定向至文件，进程崩溃时运行时将向该文件记录协程调用栈信息
//	syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
//}