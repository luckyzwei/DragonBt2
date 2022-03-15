package utils

import (
	"errors"
	"regexp"
	"time"
)

// 加强版时间计算
// https://github.com/jinzhu/now
// WeekStartDay set week start day, default is sunday
var WeekStartDay = time.Sunday

// TimeFormats default time formats will be parsed as
var TimeFormats = []string{"1/2/2006", "1/2/2006 15:4:5", "2006", "2006-1", "2006-1-2", "2006-1-2 15", "2006-1-2 15:4", "2006-1-2 15:4:5", "1-2", "15:4:5", "15:4", "15", "15:4:5 Jan 2, 2006 MST", "2006-01-02 15:04:05.999999999 -0700 MST"}

// TimeUtil now struct
type TimeUtil struct {
	time.Time
}

// NewTimeUtil initialize TimeUtil with time
func NewTimeUtil(t time.Time) *TimeUtil {
	return &TimeUtil{t}
}

// BeginningOfMinute beginning of minute
func BeginningOfMinute() time.Time {
	return NewTimeUtil(time.Now()).BeginningOfMinute()
}

// BeginningOfHour beginning of hour
func BeginningOfHour() time.Time {
	return NewTimeUtil(time.Now()).BeginningOfHour()
}

// BeginningOfDay beginning of day
func BeginningOfDay() time.Time {
	return NewTimeUtil(time.Now()).BeginningOfDay()
}

// BeginningOfWeek beginning of week
func BeginningOfWeek() time.Time {
	return NewTimeUtil(time.Now()).BeginningOfWeek()
}

// BeginningOfMonth beginning of month
func BeginningOfMonth() time.Time {
	return NewTimeUtil(time.Now()).BeginningOfMonth()
}

// BeginningOfQuarter beginning of quarter
func BeginningOfQuarter() time.Time {
	return NewTimeUtil(time.Now()).BeginningOfQuarter()
}

// BeginningOfYear beginning of year
func BeginningOfYear() time.Time {
	return NewTimeUtil(time.Now()).BeginningOfYear()
}

// EndOfMinute end of minute
func EndOfMinute() time.Time {
	return NewTimeUtil(time.Now()).EndOfMinute()
}

// EndOfHour end of hour
func EndOfHour() time.Time {
	return NewTimeUtil(time.Now()).EndOfHour()
}

// EndOfDay end of day
func EndOfDay() time.Time {
	return NewTimeUtil(time.Now()).EndOfDay()
}

// EndOfWeek end of week
func EndOfWeek() time.Time {
	return NewTimeUtil(time.Now()).EndOfWeek()
}

// EndOfMonth end of month
func EndOfMonth() time.Time {
	return NewTimeUtil(time.Now()).EndOfMonth()
}

// EndOfQuarter end of quarter
func EndOfQuarter() time.Time {
	return NewTimeUtil(time.Now()).EndOfQuarter()
}

// EndOfYear end of year
func EndOfYear() time.Time {
	return NewTimeUtil(time.Now()).EndOfYear()
}

// Monday monday
func Monday() time.Time {
	return NewTimeUtil(time.Now()).Monday()
}

// Sunday sunday
func Sunday() time.Time {
	return NewTimeUtil(time.Now()).Sunday()
}

// EndOfSunday end of sunday
func EndOfSunday() time.Time {
	return NewTimeUtil(time.Now()).EndOfSunday()
}

//! today timestamp
func GetTodayTime(hour int, minute int, second int) int64 {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, hour, minute, second, 0, time.Now().Location()).Unix()
}

// Parse parse string to time
func Parse(strs ...string) (time.Time, error) {
	return NewTimeUtil(time.Now()).Parse(strs...)
}

// ParseInLocation parse string to time in location
func ParseInLocation(loc *time.Location, strs ...string) (time.Time, error) {
	return NewTimeUtil(time.Now().In(loc)).Parse(strs...)
}

// MustParse must parse string to time or will panic
func MustParse(strs ...string) time.Time {
	return NewTimeUtil(time.Now()).MustParse(strs...)
}

// MustParseInLocation must parse string to time in location or will panic
func MustParseInLocation(loc *time.Location, strs ...string) time.Time {
	return NewTimeUtil(time.Now().In(loc)).MustParse(strs...)
}

// Between check now between the begin, end time or not
func Between(time1, time2 string) bool {
	return NewTimeUtil(time.Now()).Between(time1, time2)
}

// BeginningOfMinute beginning of minute
func (now *TimeUtil) BeginningOfMinute() time.Time {
	return now.Truncate(time.Minute)
}

// BeginningOfHour beginning of hour
func (now *TimeUtil) BeginningOfHour() time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, now.Time.Hour(), 0, 0, 0, now.Time.Location())
}

// BeginningOfDay beginning of day
func (now *TimeUtil) BeginningOfDay() time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, now.Time.Location())
}

// BeginningOfWeek beginning of week
func (now *TimeUtil) BeginningOfWeek() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())

	if WeekStartDay != time.Sunday {
		weekStartDayInt := int(WeekStartDay)

		if weekday < weekStartDayInt {
			weekday = weekday + 7 - weekStartDayInt
		} else {
			weekday = weekday - weekStartDayInt
		}
	}
	return t.AddDate(0, 0, -weekday)
}

// BeginningOfMonth beginning of month
func (now *TimeUtil) BeginningOfMonth() time.Time {
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
}

// BeginningOfQuarter beginning of quarter
func (now *TimeUtil) BeginningOfQuarter() time.Time {
	month := now.BeginningOfMonth()
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

// BeginningOfYear BeginningOfYear beginning of year
func (now *TimeUtil) BeginningOfYear() time.Time {
	y, _, _ := now.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, now.Location())
}

// EndOfMinute end of minute
func (now *TimeUtil) EndOfMinute() time.Time {
	return now.BeginningOfMinute().Add(time.Minute - time.Nanosecond)
}

// EndOfHour end of hour
func (now *TimeUtil) EndOfHour() time.Time {
	return now.BeginningOfHour().Add(time.Hour - time.Nanosecond)
}

// EndOfDay end of day
func (now *TimeUtil) EndOfDay() time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())
}

// EndOfWeek end of week
func (now *TimeUtil) EndOfWeek() time.Time {
	return now.BeginningOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// EndOfMonth end of month
func (now *TimeUtil) EndOfMonth() time.Time {
	return now.BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// EndOfQuarter end of quarter
func (now *TimeUtil) EndOfQuarter() time.Time {
	return now.BeginningOfQuarter().AddDate(0, 3, 0).Add(-time.Nanosecond)
}

// EndOfYear end of year
func (now *TimeUtil) EndOfYear() time.Time {
	return now.BeginningOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// Monday monday
func (now *TimeUtil) Monday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return t.AddDate(0, 0, -weekday+1)
}

// Sunday sunday
func (now *TimeUtil) Sunday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	}
	return t.AddDate(0, 0, (7 - weekday))
}

// EndOfSunday end of sunday
func (now *TimeUtil) EndOfSunday() time.Time {
	return NewTimeUtil(now.Sunday()).EndOfDay()
}

func parseWithFormat(str string) (t time.Time, err error) {
	for _, format := range TimeFormats {
		t, err = time.Parse(format, str)
		if err == nil {
			return
		}
	}
	err = errors.New("Can't parse string as time: " + str)
	return
}

var hasTimeRegexp = regexp.MustCompile(`(\s+|^\s*)\d{1,2}((:\d{1,2})*|((:\d{1,2}){2}\.(\d{3}|\d{6}|\d{9})))\s*$`) // match 15:04:05, 15:04:05.000, 15:04:05.000000 15, 2017-01-01 15:04, etc
var onlyTimeRegexp = regexp.MustCompile(`^\s*\d{1,2}((:\d{1,2})*|((:\d{1,2}){2}\.(\d{3}|\d{6}|\d{9})))\s*$`)      // match 15:04:05, 15, 15:04:05.000, 15:04:05.000000, etc

// Parse parse string to time
func (now *TimeUtil) Parse(strs ...string) (t time.Time, err error) {
	var (
		setCurrentTime  bool
		parseTime       []int
		currentTime     = []int{now.Nanosecond(), now.Second(), now.Minute(), now.Hour(), now.Day(), int(now.Month()), now.Year()}
		currentLocation = now.Location()
		onlyTimeInStr   = true
	)

	for _, str := range strs {
		hasTimeInStr := hasTimeRegexp.MatchString(str) // match 15:04:05, 15
		onlyTimeInStr = hasTimeInStr && onlyTimeInStr && onlyTimeRegexp.MatchString(str)
		if t, err = parseWithFormat(str); err == nil {
			location := t.Location()
			if location.String() == "UTC" {
				location = currentLocation
			}

			parseTime = []int{t.Nanosecond(), t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}

			for i, v := range parseTime {
				// Don't reset hour, minute, second if current time str including time
				if hasTimeInStr && i <= 3 {
					continue
				}

				// If value is zero, replace it with current time
				if v == 0 {
					if setCurrentTime {
						parseTime[i] = currentTime[i]
					}
				} else {
					setCurrentTime = true
				}

				// if current time only includes time, should change day, month to current time
				if onlyTimeInStr {
					if i == 4 || i == 5 {
						parseTime[i] = currentTime[i]
						continue
					}
				}
			}

			t = time.Date(parseTime[6], time.Month(parseTime[5]), parseTime[4], parseTime[3], parseTime[2], parseTime[1], parseTime[0], location)
			currentTime = []int{t.Nanosecond(), t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
		}
	}
	return
}

// MustParse must parse string to time or it will panic
func (now *TimeUtil) MustParse(strs ...string) (t time.Time) {
	t, err := now.Parse(strs...)
	if err != nil {
		panic(err)
	}
	return t
}

// Between check time between the begin, end time or not
func (now *TimeUtil) Between(begin, end string) bool {
	beginTime := now.MustParse(begin)
	endTime := now.MustParse(end)
	return now.After(beginTime) && now.Before(endTime)
}

type TimeInfo struct {
	Hour   int
	Minute int
	Second int
}

func NewTimeInfo(timeStamp int) *TimeInfo {
	hour := timeStamp / 10000
	timeStamp = timeStamp % 10000
	minute := timeStamp / 100
	timeStamp = timeStamp % 100
	second := timeStamp
	return &TimeInfo{hour, minute, second}
}

func GetTodyDay() int64 {
	now := time.Now()
	timeSet := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local)
	return timeSet.Unix()
}
