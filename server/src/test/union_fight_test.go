package test

import (
	. "game"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	openTime := time.Date(2019, 3, 1, 5, 0, 0, 0, time.Local)
	today := GetCsvMgr().GetTodyDay()
	hour := 20
	min := 0
	totalSec := int64((hour-5)*HOUR_SECS + min*MIN_SECS)
	var  cdDay int64 = 3
	var  openDay int64 = 5
	var conti int64 = 1
	diffTmp := today - openTime.Unix() - int64((openDay-1)*DAY_SECS)
	cdSecs := int64(cdDay * DAY_SECS)
	flag := int64(diffTmp % cdSecs)
	var total int64
	if flag == 0 {
		total = int64(diffTmp / cdSecs * cdDay)
	} else {
		total =int64( diffTmp/cdSecs*cdDay + cdDay)
	}

	// 2.28, 3,4(1,3),3.8(1,3),3.12(1,3)
	var n int64
	if time.Now().Sub(openTime) < 4*DAY_SECS  {
		n = 1
	} else {
		n = (time.Now().Unix() - openTime.Unix() - (openDay-1) * DAY_SECS) / ((cdDay+conti) * DAY_SECS) + 1
	}

	t.Log("n : ", n)
	s1 := openTime.Unix() + (openDay-1) * DAY_SECS + (n-1)*(cdDay + conti) * DAY_SECS
	t.Log("s1:", time.Unix(s1, 0).Format(DATEFORMAT))

	endTime := int64(0)
	startTime := int64(0)
	if endTime <= time.Now().Unix() {
		startTime = int64(openTime.Unix() + (openDay-1)*DAY_SECS + total*DAY_SECS + totalSec)
		endTime = startTime + 3600
	}

	t.Log("startTime:", time.Unix(startTime, 0).Format(DATEFORMAT))
	t.Log("endTime:", time.Unix(endTime, 0).Format(DATEFORMAT))
}
