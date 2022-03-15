package test

import (
	"fmt"
	. "game"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GameTestSuite struct {
	suite.Suite
	Player       *Player
	ModKingTask *ModKingTask
}


func (m *GameTestSuite) Dump(x ...interface{}) {
	fmt.Printf("%# v\n", pretty.Formatter(x))
}

// 初始化配置
func (m *GameTestSuite) SetupSuite() {
	//fmt.Println("所有次数前调用一次")
	GetServer().InitConfig()
	GetCsvMgr().InitData()
	//GetServer().ConnectDB()
	m.Player = NewPlayer(1)
	//m.Player.InitPlayerData()
	//m.Player.OtherPlayerData()

	//m.Dump(GetCsvMgr().PlayTimeMap)
}

func (m *GameTestSuite) SetupTest() {
	//fmt.Println("每次测试前调用一次")

}

func (m *GameTestSuite) TearDownSuite() {
	//GetServer().DBUser.Close()
	//GetServer().DBLog.Close()
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}
