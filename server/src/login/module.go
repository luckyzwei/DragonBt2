package main


type IModule interface {
	GetName() string
	Init() bool
	Destory()
}

type ModuleMgr struct {
	modules map[string]IModule
}

func (self *ModuleMgr) RegisterModule(module IModule)  {
	self.modules[module.GetName()] = module;
}

func (self *ModuleMgr) Init()  {
	for _, v := range self.modules {
		v.Init();
	}
}

func (self *ModuleMgr) Destory() {
	for _, v := range self.modules {
		v.Destory();
	}
}

var gModuleMgr *ModuleMgr = nil;

func GetModuleMgr()*ModuleMgr  {

	if (gModuleMgr == nil) {
		gModuleMgr = new(ModuleMgr);
		gModuleMgr.modules = make(map[string]IModule)
		gModuleMgr.RegisterModule(new(ModLogin))
	}
	return gModuleMgr;
}




