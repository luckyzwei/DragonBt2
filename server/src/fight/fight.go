package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"
)

type FightMgr struct {
	Locker              *sync.RWMutex
	AutoId              int64
	FightNodes          map[int64]*FightServerExtendNode
	FightResultNodeChan chan *FightResultExtendNode
	TempRecord          map[int]map[int64]int
}

var fightmgrsingleton *FightMgr = nil

//! public
func GetFightMgr() *FightMgr {
	if fightmgrsingleton == nil {
		fightmgrsingleton = new(FightMgr)
		fightmgrsingleton.Locker = new(sync.RWMutex)
		fightmgrsingleton.AutoId = 0
		fightmgrsingleton.FightNodes = make(map[int64]*FightServerExtendNode)
		fightmgrsingleton.TempRecord = make(map[int]map[int64]int)
	}

	return fightmgrsingleton
}

func (self *FightMgr) InitFightServers() {

	for i := 0; i < len(GetServer().Con.ServersCon); i++ {
		go self.GetServerFightsRoutine(GetServer().Con.ServersCon[i].Id)
	}

	self.FightResultNodeChan = make(chan *FightResultExtendNode, len(GetServer().Con.ServersCon)*16)

	routineCount := self.GetResultRoutineCount(len(GetServer().Con.ServersCon))
	for i := 0; i < routineCount; i++ {
		go self.SetServerFightResultRoutine()
	}
}

func (self *FightMgr) ReloadFightServers(newCon *Config) {

	//为新加的服务器启动goroutine
	for i := 0; i < len(newCon.ServersCon); i++ {
		if GetServer().GetServerCon(newCon.ServersCon[i].Id) == nil {
			go self.GetServerFightsRoutine(newCon.ServersCon[i].Id)
		}
	}

	//
	routineCount := self.GetResultRoutineCount(len(GetServer().Con.ServersCon))
	newRoutineCount := self.GetResultRoutineCount(len(newCon.ServersCon))
	for i := 0; i < (newRoutineCount - routineCount); i++ {
		go self.SetServerFightResultRoutine()
	}
}

func (self *FightMgr) GetResultRoutineCount(serverCount int) int {
	return serverCount * 2
}

func (self *FightMgr) GetServerFightsRoutine(serverid int) {

	GetServer().Wait.Add(1)
	ticker := time.NewTicker(time.Second * 1)
	for {

		if GetServer().ShutDown {
			break
		}

		select {
		case <-ticker.C:
			self.GetServerFightsImpl(serverid)
		}
	}

	ticker.Stop()
	GetServer().Wait.Done()
}

func (self *FightMgr) SetServerFightResultRoutine() {

	GetServer().Wait.Add(1)
	ticker := time.NewTicker(time.Second * 1)
	for {
		if GetServer().ShutDown {
			break
		}

		select {
		case node := <-self.FightResultNodeChan:
			self.SetServerFightResultImpl(node)
		case <-ticker.C:

		}
	}
	ticker.Stop()
	GetServer().Wait.Done()
}

func (self *FightMgr) SetServerFightResultImpl(node *FightResultExtendNode) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	serverConfig := GetServer().GetServerCon(node.node.serverid)
	if serverConfig == nil {
		return
	}

	node.resultNode.Fightid = node.node.FightServerNode.Id //发送给对应的游戏服务器之前需要重新decode成游戏服务器对应的id

	url := "http://" + serverConfig.Host + "/fightserver?msgtype=set&data="
	data, err := json.Marshal(node.resultNode)
	if err != nil {
		return
	}
	url += string(data[:])

	body := bytes.NewBuffer([]byte(""))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(6 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	_, err1 := ioutil.ReadAll(resp.Body)
	//log.Println("result:=", result)
	defer resp.Body.Close()
	if err1 != nil {
		return
	}

	self.Locker.Lock()
	delete(self.FightNodes, node.node.NewId)
	_, ok := self.TempRecord[node.node.serverid]
	if ok {
		_, okRecord := self.TempRecord[node.node.serverid][node.node.Id]
		if okRecord {
			delete(self.TempRecord[node.node.serverid], node.node.Id)
		}
	}
	self.Locker.Unlock()
}

func (self *FightMgr) IsFightNodeExists(serverid int, node *FightServerNode) bool {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, value := range self.FightNodes {
		if value.Id == node.Id && value.serverid == serverid {
			return true
		}
	}

	return false
}

func (self *FightMgr) AddFightNode(serverid int, node *FightServerNode) {

	self.Locker.Lock()

	_, ok := self.TempRecord[serverid]
	if !ok {
		self.TempRecord[serverid] = make(map[int64]int)
	}
	_, okRecord := self.TempRecord[serverid][node.Id]
	if !okRecord {
		self.AutoId += 1
		extendNode := new(FightServerExtendNode)
		extendNode.NewId = self.AutoId
		extendNode.serverid = serverid
		extendNode.FightServerNode = *node
		extendNode.status = 0
		extendNode.addTime = time.Now().Unix()
		self.FightNodes[self.AutoId] = extendNode
		self.TempRecord[serverid][node.Id] = 1
		for i := 0; i < len(GetServer().NeedCount); i++ {
			GetServer().NeedCount[i]++
		}
	}

	self.Locker.Unlock()
}

func (self *FightMgr) GetTopFight() *FightServerExtendNode {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	now := time.Now().Unix()
	for _, value := range self.FightNodes {
		if value.addTime+600 < now {
			continue
			//delete(self.FightNodes, value.NewId)
			//delete(self.TempRecord[value.serverid], value.Id)
		}
		if value.status == 0 {
			value.status = 1
			return value
		}
	}

	return nil
}

func (self *FightMgr) SetResult(node *FightResultNode) {

	self.Locker.Lock()
	result, ok := self.FightNodes[node.Fightid]
	self.Locker.Unlock()

	if !ok {
		return
	}

	resultNode := new(FightResultExtendNode)
	resultNode.resultNode = node
	resultNode.node = result

	self.FightResultNodeChan <- resultNode
	//go self.SetServerFightResultRoutine(result.serverid, node)
	//delete(self.FightNodes, node.Fightid)
}

func (self *FightMgr) GetServerFightsImpl(serverid int) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	serverConfig := GetServer().GetServerCon(serverid)
	if serverConfig == nil {
		return
	}

	url := "http://" + serverConfig.Host + "/fightserver?msgtype=get"
	body := bytes.NewBuffer([]byte(""))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(6 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	if string(result[:]) == "false" {
		return
	}

	var node FightServerNode
	err = json.Unmarshal([]byte(result), &node)
	if err != nil {
		log.Println("err=", err)
		return
	}

	if self.IsFightNodeExists(serverid, &node) == false {
		self.AddFightNode(serverid, &node)
	}
}

func (self *FightMgr) GetNeedNow() int {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	count := 0
	for _, value := range self.FightNodes {
		if value.status == 0 {
			count++
		}
	}

	return count
}
