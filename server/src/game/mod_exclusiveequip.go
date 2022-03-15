package game

type ExclusiveEquip struct {
	Id       int         `json:"id"`       //! 装备配置Id
	AttrInfo []*AttrInfo `json:"attrinfo"` //! 属性信息
	Lv       int         `json:"lv"`       //!
	Skill    int         `json:"skill"`    //!
	UnLock   int         `json:"UnLock"`   //!  0未解锁  1解锁
}

// 装备计算属性
func (self *ExclusiveEquip) CalAttr() {
	configBase := GetCsvMgr().ExclusiveEquipConfigMap[self.Id]
	if configBase == nil {
		return
	}

	config := GetCsvMgr().ExclusiveStrengthen[self.Id][self.Lv]

	for _, v := range self.AttrInfo {
		if v.AttrId <= 0 {
			continue
		}
		index := v.AttrId - 1
		v.AttrType = configBase.BaseType[index]
		v.AttrValue = configBase.BaseValue[index]

		if config != nil {
			v.AttrValue += config.Value[index]
		}
	}

	if config != nil {
		self.Skill = config.Skill
	} else {
		self.Skill = configBase.Skill
	}
}
