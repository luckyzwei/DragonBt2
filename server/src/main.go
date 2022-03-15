package main

import (
	"game"
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

//! ctrl+z
func handleTSTP(ch *chan os.Signal) {
	for {
		<-*ch
		log.Println("get a SIGTSTP")
		game.GetCsvMgr().Reload()
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
		game.GetServer().Close()

	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	//! 读入系统配置
	game.GetServer().InitConfig()
	//! 读入csv
	game.GetCsvMgr().InitData()

	//! 连接数据库
	game.GetServer().ConnectDB()

	//! 注册信号量
	handleSignal(syscall.SIGPIPE, handlePIPE)
	//! 这个注释不删除,Linux下用
	//handleSignal(syscall.SIGTSTP, handleTSTP)
	handleSignal(syscall.SIGTRAP, handleTRAP)
	handleSignal(syscall.SIGINT, handleINT)


	game.GetServer().GoService()
	game.GetBackStageMgr().Init()

	http.HandleFunc("/fightserver", game.FightServer)
	http.Handle("/", game.GetServer().GetConnectHandler())
	log.Println("绑定ip:", game.GetServer().Con.Host)
	log.Println("服务器版本1001")
	log.Fatal(http.ListenAndServe(game.GetServer().Con.Host, nil))

}
