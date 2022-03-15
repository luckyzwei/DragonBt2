package main

import (
	"net/http"
)

type BackStageMgr struct {
}

var backstagemgrsingleton *BackStageMgr = nil

//! public
func GetBackStageMgr() *BackStageMgr {
	if backstagemgrsingleton == nil {
		backstagemgrsingleton = new(BackStageMgr)
	}
	return backstagemgrsingleton
}

func (self *BackStageMgr) Init() {
	http.HandleFunc("/reloadconfig", self.ReloadConfig)
}

func (self *BackStageMgr) ReloadConfig(w http.ResponseWriter, r *http.Request) {

	GetServer().InitConfig()
}
