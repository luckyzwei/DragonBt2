package game

type ModAll struct {
	player   *Player
	Module   map[string]ModBase
	Handlers map[string]func(body []byte)
}

func NewModAll(player *Player) *ModAll {
	p := new(ModAll)
	p.player = player
	// 这里注册
	p.Module = map[string]ModBase{
		"chat":          new(ModChat),
		"pass":          new(ModPass),
		"recharge":      new(ModRecharge),
		"weekplan":      new(ModWeekPlan),
		"task":          new(ModTask),
		"targettask":    new(ModTargetTask),
		"shop":          new(ModShop),
		"honourshop":    new(ModHonourShop),
		"viprecharge":   new(ModVipRecharge),
		"bag":           new(ModBag),
		"hero":          new(ModHero),
		"friend":        new(ModFriend),
		"find":          new(ModFind),
		"mail":          new(ModMail),
		"gm":            new(ModGm),
		"union":         new(ModUnion),
		"top":           new(ModTop),
		"activity":      new(ModActivity),
		"luckshop":      new(ModLuckShop),
		"fund":          new(ModFund),
		"actop":         new(ModActop),
		"redpac":        new(ModRedPac),
		"dailyrecharge": new(ModDailyRecharge),
		"equip":         new(ModEquip),
		"artifactequip": new(ModArtifactEquip),
		"tower":         new(ModTower),
		"head":          new(ModHead),
		"moneytask":     new(ModMoneyTask),
		"guide":         new(ModGuide),
		"timegift":      new(ModTimeGift),
		"nobilitytask":  new(ModNobilityTask),
		"turntable":     new(ModTurnTable),
		"accesscard":    new(ModAccessCard),
		"onhook":        new(ModOnHook),
		//"hydra":           new(ModHydra),
		"pit":                  new(ModPit),
		"clientsign":           new(ModClientSign),
		"newpit":               new(ModNewPit),
		"instance":             new(ModInstance),
		"support":              new(ModSupportHero),
		"entanglement":         new(ModEntanglement),
		"reward":               new(ModReward),
		"ranktask":             new(ModRankTask),
		"crystal":              new(ModResonanceCrystal),
		"arena":                new(ModArena),
		"battle":               new(ModBattle),
		"arenaspecial":         new(ModArenaSpecial),
		"activitygift":         new(ModActivityGift),
		"growthgift":           new(ModGrowthGift),
		"skin":                 new(ModSkin),
		"specialpurchase":      new(ModSpecialPurchase),
		"lifetree":             new(ModLifeTree),
		"interstellar":         new(ModInterStellar),
		"activityboss":         new(ModActivityBoss),
		"general":              new(ModGeneral),
		"herogrow":             new(ModHeroGrow),
		"crossarena":           new(ModCrossArena),
		"crossarena3v3":        new(ModCrossArena3V3),
		"activitybossfestival": new(ModActivityBossFestival),
		"lotterydraw":          new(ModLotteryDraw),
		"consumertop":          new(ModConsumerTop),
		"beauty":               new(ModBeauty),
		"horse":                new(ModHorse),
		"team":                 new(ModTeam),
	}
	p.RegHandle()

	return p
}

func (self *ModAll) RegHandle() {
	self.Handlers = make(map[string]func(body []byte))
	self.Module["bag"].(*ModBag).onReg(self.Handlers)
	self.Module["head"].(*ModHead).onReg(self.Handlers)
	self.Module["moneytask"].(*ModMoneyTask).onReg(self.Handlers)
	self.Module["onhook"].(*ModOnHook).onReg(self.Handlers)
	self.Module["accesscard"].(*ModAccessCard).onReg(self.Handlers)
	//self.Module["hydra"].(*ModHydra).onReg(self.Handlers)
	self.Module["pit"].(*ModPit).onReg(self.Handlers)
	self.Module["newpit"].(*ModNewPit).onReg(self.Handlers)
	self.Module["instance"].(*ModInstance).onReg(self.Handlers)
	self.Module["clientsign"].(*ModClientSign).onReg(self.Handlers)
	self.Module["team"].(*ModTeam).onReg(self.Handlers)
	self.Module["hero"].(*ModHero).onReg(self.Handlers)
	self.Module["friend"].(*ModFriend).onReg(self.Handlers)
	self.Module["find"].(*ModFind).onReg(self.Handlers)
	self.Module["interstellar"].(*ModInterStellar).onReg(self.Handlers)
	self.Module["viprecharge"].(*ModVipRecharge).onReg(self.Handlers)
	self.Module["nobilitytask"].(*ModNobilityTask).onReg(self.Handlers)
	self.Module["turntable"].(*ModTurnTable).onReg(self.Handlers)
	self.Module["equip"].(*ModEquip).onReg(self.Handlers)
	self.Module["artifactequip"].(*ModArtifactEquip).onReg(self.Handlers)
	self.Module["support"].(*ModSupportHero).onReg(self.Handlers)
	self.Module["entanglement"].(*ModEntanglement).onReg(self.Handlers)
	self.Module["reward"].(*ModReward).onReg(self.Handlers)
	self.Module["ranktask"].(*ModRankTask).onReg(self.Handlers)
	self.Module["crystal"].(*ModResonanceCrystal).onReg(self.Handlers)
	self.Module["battle"].(*ModBattle).onReg(self.Handlers)
	self.Module["arenaspecial"].(*ModArenaSpecial).onReg(self.Handlers)
	self.Module["activitygift"].(*ModActivityGift).onReg(self.Handlers)
	self.Module["growthgift"].(*ModGrowthGift).onReg(self.Handlers)
	self.Module["skin"].(*ModSkin).onReg(self.Handlers)
	self.Module["lifetree"].(*ModLifeTree).onReg(self.Handlers)
	self.Module["targettask"].(*ModTargetTask).onReg(self.Handlers)
	self.Module["activityboss"].(*ModActivityBoss).onReg(self.Handlers)
	self.Module["herogrow"].(*ModHeroGrow).onReg(self.Handlers)
	self.Module["crossarena"].(*ModCrossArena).onReg(self.Handlers)
	self.Module["crossarena3v3"].(*ModCrossArena3V3).onReg(self.Handlers)
	self.Module["activitybossfestival"].(*ModActivityBossFestival).onReg(self.Handlers)
	self.Module["lotterydraw"].(*ModLotteryDraw).onReg(self.Handlers)
	self.Module["honourshop"].(*ModHonourShop).onReg(self.Handlers)

}

func (self *ModAll) GetModule(name string) ModBase {
	return self.Module[name]
}

// 得到数据[拉取玩家时调用]
func (self *ModAll) GetData() {
	for _, value := range self.Module {
		value.OnGetData(self.player)
	}
}

// 登录时拉起玩家数据
func (self *ModAll) GetOtherData() {
	for _, value := range self.Module {
		value.OnGetOtherData()
	}
}

// 保存数据
func (self *ModAll) Save(sql bool) {
	for _, value := range self.Module {
		value.OnSave(sql)
	}
}

// 得到消息
func (self *ModAll) OnMsg(ctrl string, body []byte) {
	handle, ok := self.Handlers[ctrl]
	if !ok {
		found := false
		for _, value := range self.Module {
			//if name == "chat" {
			//	continue
			//}
			if found = value.OnMsg(ctrl, body); found {
				break
			}
		}

		if !found {
			LogError("消息:", ctrl, "没有找到")
		}
	} else {
		handle(body)
	}

}
