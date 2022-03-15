package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func addSig(s chan os.Signal) chan os.Signal {
	signal.Notify(s, syscall.SIGTERM)
	signal.Notify(s, syscall.SIGINT)
	return s
}

func exit(s chan os.Signal) {
	sig := <-s
	fmt.Printf("caught sig: %+v", sig)
	os.Exit(0)
}

func Stop()  {
	exit(addSig(make(chan os.Signal)))
}

func main() {
	fmt.Println("启动机器人")
	GetRobotCsvMgr().Reload()
	GetRobotMgr().InitConfig()
	GetRobotMgr().Init()

	Stop()
}
