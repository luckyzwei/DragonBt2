package test

import (
	. "game"
)

func (m *GameTestSuite) TestShuffle() {
	vals := []int{1, 3, 4, 6, 10, 11, 12, 13}
	for i := 0; i < 10000; i++ {
		res := make(map[int]int)
		GetMineMgr().Shuffle(vals)
		for _, v := range vals {
			res[v] = 1
		}

		if len(res) != 8 {
			m.T().Error("len(res) != 8")
		}
	}

}

