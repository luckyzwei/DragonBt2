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
	//mapa := make(map[string]*game.JS_Hero)
	//hero := new(game.JS_Hero)
	//hero.Initaaa()
	//hero.Encode()
	//mapa["3001"] = hero
	//aa, _ := json.Marshal(&mapa)
	//log.Println(string(aa))

	//var mapb map[string]*game.JS_Hero
	//json.Unmarshal(aa, &mapb)
	//log.Println(mapb)

	//for i := 0; i < 100; i++ {
	//	game.GetServer().InsertBarrage(1000, "ddd", "adfdsa")
	//}

	//return
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	//! 读入配置
	log.Println("初始化服务器配置:")
	GetServer().InitConfig()

	//
	log.Println("初始化数据库连接:")
	GetServer().ConnectDB()

	log.Println("初始化模块:")
	GetModuleMgr().Init()

	go GetServer().Run()

	//! 注册信号量
	handleSignal(syscall.SIGPIPE, handlePIPE)
	handleSignal(syscall.SIGTRAP, handleTRAP)
	handleSignal(syscall.SIGINT, handleINT)

	http.HandleFunc("/api", GetServer().Handler)
	http.Handle("/", GetServer().GetConnectHandler())
	log.Println("绑定ip:", GetServer().Con.Host)
	log.Fatal(http.ListenAndServe(GetServer().Con.Host, nil))
}
