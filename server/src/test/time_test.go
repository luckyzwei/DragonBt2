package test

import (
	"testing"
)
func TestDayTime(t *testing.T)  {
	//endOfDay  := game.EndOfDay().Add(5 * time.Hour)
	//t.Log(endOfDay.Format(game.DATEFORMAT))

	price := 8.0

	for i :=0; i <20; i++ {
		price = price * (1+0.1)
		if price >= 30.0 {
			println(int(i+1))
			break
		}
	}
}
