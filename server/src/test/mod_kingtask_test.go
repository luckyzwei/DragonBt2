package test

import (
	"fmt"
	. "game"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/suite"
	"testing"
)

type KingTaskSuite struct {
	suite.Suite
	Player       *Player
	ModKingTask *ModKingTask
}


func (m *KingTaskSuite) Dump(x ...interface{}) {
	fmt.Printf("%# v\n", pretty.Formatter(x))
}

// 初始化配置
func (m *KingTaskSuite) SetupSuite() {
	GetServer().InitConfig()
	GetCsvMgr().InitData()
}

func TestKingTaskSuite(t *testing.T) {
	suite.Run(t, new(KingTaskSuite))
}

// white box test
// Growthtask_King.csv
func (m *KingTaskSuite) GetKingTask(hard int) *GrowthtaskKingConfig {
	for _, v := range GetCsvMgr().TaskKingConfigMap {
		if v.Type == hard {
			return v
		}
	}
	return nil
}

func (m *KingTaskSuite) NewTask() *KingTask {
	config := m.GetKingTask(1)
	m.NotEqual(config, nil, "should be equal")
	task := NewKingTask(config.Taskid)
	return task
}

func (m *KingTaskSuite) NewKingBox() *KingBox {
	task := m.NewTask()
	m.NotEqual(task, nil, "should be equal")
	box := NewKingBox(task.Taskid)
	return box
}


func (m *KingTaskSuite) TestFinishTask()  {
	m.Dump(GetCsvMgr().TariffConfig)
	info := GetCsvMgr().GetTariffConfig(9, 1)
	m.Dump(info)
}


