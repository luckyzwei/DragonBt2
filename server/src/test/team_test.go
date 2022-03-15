package test

import (
	"game"
	"testing"
)

func TestWin(t *testing.T)  {
	t.Log()
	winTimes := 0
	failTimes := 0
	for i := 0; i < 10000; i++ {
		if game.GetCsvMgr().IsTeamWin(4001, 5000, 100, 80, 200, 90, 60) {
			winTimes += 1
		} else {
			failTimes += 1
		}
	}
	//if failTimes <= 0 {
	//	t.Error("failTimes <= 0")
	//}

	if failTimes <= 0 {
		t.Error("failTimes > 0")
	}
	t.Log("winTimes:", winTimes, ", failTimes:", failTimes)

}
