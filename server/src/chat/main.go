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
	//handleSignal(syscall.SIGTSTP, handleTSTP)
	handleSignal(syscall.SIGTRAP, handleTRAP)
	handleSignal(syscall.SIGINT, handleINT)

	GetChatServer().InitConfig()

	GetChatServer().startMsgGroup(GetChatServer().Con.InitGroupNum)

	go GetChatServer().Run()

	http.Handle("/", GetChatServer().GetConnectHandler())
	log.Println("绑定ip:", GetChatServer().Con.Host)
	log.Fatal(http.ListenAndServe(GetChatServer().Con.Host, nil))
}
