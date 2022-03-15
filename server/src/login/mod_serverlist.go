package main

import "net/http"

type ModServerList struct {

}

func (self* ModServerList) GetName() string {
	return ModName_ServerList
}

func (self* ModServerList) Init() bool {
	GetServer().RegisterHandler(CMD_ServerList, self.GetList)
	return true;
}

func (self* ModServerList) Destory() {

}

func (self* ModServerList) GetList(r *http.Request) interface{}  {


	var ret S2C_Reg;

	return &ret;
}


