package test

import (
	. "game"
)

func (m *GameTestSuite) TestRandRankTest() {
	info := new(San_Pvp)
	info.Rankid = 1000
	fightid := [5]int{0, 0, 0, 0, 0}
	RandRank(info, &fightid)
	m.NotEqual(fightid, [5]int{0, 0, 0, 0, 0})
}
