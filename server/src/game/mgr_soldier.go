package game

var soldierMgr *SoldierMgr

type SoldierMgr struct{}

func GetSoldierMgr() *SoldierMgr {
	if soldierMgr == nil {
		soldierMgr = new(SoldierMgr)
	}
	return soldierMgr
}

// 随机属性条数
func (self *SoldierMgr) RandSoldierAttrNum(id int, lv int) int {
	info, ok := GetCsvMgr().MercenaryLvMap[id]
	if !ok {
		return 0
	}

	v, ok := info[lv]
	if !ok {
		return 0
	}
	sum := 0
	for _, check := range v.Gets {
		sum += check
	}

	rand := HF_GetRandom(sum)
	total := 0
	for i, check := range v.Gets {
		total += check
		if rand < total {
			return i + 1
		}
	}
	return 0
}

// 通过Id，等级获得士兵配置
func (self *SoldierMgr) GetSoldierLvConfig(id int, lv int) *MercenaryLv {
	info, ok := GetCsvMgr().MercenaryLvMap[id]
	if !ok {
		return nil
	}

	config, ok := info[lv]
	if !ok {
		return nil
	}
	return config
}

// 通过Id，等级获得士兵配置
func (self *SoldierMgr) GetSoldierConfig(id int) *MercenaryConfig {
	config, ok := GetCsvMgr().MercenaryConfig[id]
	if !ok {
		return nil
	}

	return config
}

// 通过阶级class, 类型index来随机属性类型
func (self *SoldierMgr) RandAttrType(class int, index int) (randType int, randValue int) {
	info, ok := GetCsvMgr().MercenaryRandomGroup[class]
	if !ok {
		return
	}

	config, ok := info[index]
	if !ok {
		return
	}

	sum := 0
	for _, v := range config {
		sum += v.Weight
	}

	rand := HF_GetRandom(sum)
	total := 0

	var target *MercenaryRandom
	for _, v := range config {
		total += v.Weight
		if rand < total {
			randType = v.BaseType
			target = v
			break
		}
	}

	if target == nil {
		LogError("target == nil!")
		return
	}

	// 随机值
	randValue = HF_RandInt(1, target.InitValue)
	randValue = self.randValue(target, randValue)
	return
}

func (self *SoldierMgr) getClassAndIndex(id int) (int, int) {
	config, ok := GetCsvMgr().MercenaryConfig[id]
	if !ok {
		LogError("Mercenary config nil, id=", id)
		return 0, 0
	}
	class := config.Level
	index := config.Index

	return class, index
}

// 随机属性类型以及数值
func (self *SoldierMgr) RandAttr(id int, lv int) *Attribute {
	// 获得佣兵的配置
	config, ok := GetCsvMgr().MercenaryConfig[id]
	if !ok {
		LogError("Mercenary config nil, id=", id)
		return nil
	}

	class := config.Level
	index := config.Index
	randType, randValue := self.RandAttrType(class, index)
	//这里不好处理int64 暂时搁置 20190506 by zy
	return &Attribute{randType, int64(randValue)}
}

// 检查位置合法性
func (self *SoldierMgr) IsPosOK(config *MercenaryConfig, pos int) bool {
	switch config.Hole {
	case 1:
		if pos == 1 {
			return true
		}
	case 2:
		if pos >= 2 && pos <= 4 {
			return true
		}
	case 3:
		if pos >= 5 && pos <= 7 {
			return true
		}
	case 4:
		if pos >= 8 && pos <= 10 {
			return true
		}
	}
	return false
}

// 获取洗练消耗配置
func (self *SoldierMgr) getWashConfig(lockNum int) *TariffConfig {
	return GetCsvMgr().GetTariffConfig3(TariffSoldierWash, lockNum)
}

func (self *SoldierMgr) randValue(target *MercenaryRandom, randValue int) int {
	// 当前值与最大值比较
	threshold := float32(randValue) / float32(target.MaxValue)

	// 随机规则
	var changes []int
	if threshold <= float32(target.Threshold)/10000.0 { // 向下
		changes = append(changes, target.LowUp)
		changes = append(changes, target.LowDown)
		changes = append(changes, target.LowFlat)

	} else { // 向上
		changes = append(changes, target.HighUp)
		changes = append(changes, target.HighDown)
		changes = append(changes, target.HighFlat)
	}

	changesSum := 0
	for _, v := range changes {
		changesSum += v
	}

	mode := 2
	changesTotal := 0
	changesRand := HF_GetRandom(changesSum)
	for i, v := range changes {
		changesTotal += v
		if changesRand < changesTotal {
			mode = i
			break
		}
	}

	changeValue := 0
	switch mode {
	case 0: // 升
		changeValue = target.ChangeMax
	case 1: // 降
		changeValue = -target.ChangeMax
	case 2: // 不变
	}

	randValue += changeValue
	randValue = HF_MaxInt(randValue, 1)
	randValue = HF_MinInt(randValue, target.MaxValue)
	return randValue
}

// 通过阶级class, 类型index来随机属性类型
func (self *SoldierMgr) RandAttrTypeValues(class int, index int, attr *Attribute) *Attribute {
	info, ok := GetCsvMgr().MercenaryRandomGroup[class]
	if !ok {
		LogError("config error, class =", class)
		return nil
	}

	config, ok := info[index]
	if !ok {
		LogError("config error, index =", index)
		return nil
	}

	var target *MercenaryRandom
	for _, v := range config {
		if v.BaseType == attr.AttType {
			target = v
			break
		}
	}

	if target == nil {
		LogError("target == nil, attType=", attr.AttType, ", class :", class, ", index:", index)
		return nil
	}

	// 随机值
	randValue := self.randValue(target, int(attr.AttValue))
	return &Attribute{attr.AttType, int64(randValue)}
}
