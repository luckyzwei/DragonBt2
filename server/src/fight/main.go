package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
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
		GetServer().Close()

	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	//! 读入配置
	log.Println("初始化服务器配置:")
	GetServer().InitConfig()

	log.Println("启动战斗代理服务器:")
	GetFightMgr().InitFightServers()

	log.Println("启动后台服务:")
	GetBackStageMgr().Init()

	go GetServer().Run()

	//! 注册信号量
	handleSignal(syscall.SIGPIPE, handlePIPE)
	handleSignal(syscall.SIGTRAP, handleTRAP)
	handleSignal(syscall.SIGINT, handleINT)

	http.HandleFunc("/fightserver", GetServer().FightServer)
	log.Println("绑定ip:", GetServer().Con.Host)
	log.Fatal(http.ListenAndServe(GetServer().Con.Host, nil))
}
