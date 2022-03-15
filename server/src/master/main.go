package main

import (
	"game"
	"log"
	"master/app"
	"master/gate"
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
		//CsvUtilMgr{}.Reload()
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
		app.GetMasterApp().Close()
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	//! 读入系统配置
	//app.GetMainApp().InitConfig()
	//game.GetServer().InitConfig()
	//! 读入csv
	//game.InitData()

	//! 连接数据库
	//app.GetMasterApp().ConnectDB()

	//! 注册信号量
	handleSignal(syscall.SIGPIPE, handlePIPE)
	//! 这个注释不删除,Linux下用
	//handleSignal(syscall.SIGTSTP, handleTSTP)
	handleSignal(syscall.SIGTRAP, handleTRAP)
	handleSignal(syscall.SIGINT, handleINT)

	//game.GetServer().GoService()
	//game.GetBackStageMgr().Init()

	conf := app.GetMasterApp().GetConfig()

	http.HandleFunc("/fightserver", game.FightServer)
	http.Handle("/", gate.GetGateApp().GetConnectHandler())
	log.Println("绑定ip:", conf.Host)
	log.Println("服务器版本:", conf.ServerVer)

	app.GetMasterApp().StartService()

	//log.Fatal(http.ListenAndServe(conf.Host, nil))

}
