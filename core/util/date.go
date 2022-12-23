package util

import (
	"fmt"
	"time"
)

var (
	// (闰)年表
	leapYearMapper = map[int]int{1: 31, 2: 29, 3: 31, 4: 30, 5: 31, 6: 30, 7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31}
	yearMapper     = map[int]int{1: 31, 2: 28, 3: 31, 4: 30, 5: 31, 6: 30, 7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31}
	LAYOUT_BY_DAY  = "2006-1-2"
)

// 返回上月的时间范围，[)
func GetLastMonthDateRange(year int, month time.Month) (time.Time, time.Time) {
	switch month {
	case 1:
		return time.Date(year-1, time.December, 1, 0, 0, 0, 0, time.Local),
			time.Date(year-1, time.December, 31, 24, 0, 0, 0, time.Local)
	default:
		return time.Date(year, month-1, 1, 0, 0, 0, 0, time.Local),
			time.Date(year, month-1, GetMonthMaxDay(year, month-1), 24, 0, 0, 0, time.Local)
	}
}

// 获取前 N 天至今的时间范围
func GetLastNDayDateRange(day int) (time.Time, time.Time) {
	return time.Now().AddDate(0, 0, -1*day),
		time.Now()
}

// 获取指定年份月的天数
func GetMonthMaxDay(year int, month time.Month) int {
	if year%400 == 0 || (year%100 != 0 && year%4 == 0) {
		return leapYearMapper[int(month)]
	} else {
		return yearMapper[int(month)]
	}
}

func GetLastDayByYear(year int) time.Time {
	firstTimeByYear, _ := time.ParseInLocation(LAYOUT_BY_DAY, fmt.Sprintf("%d-%d-%d", year, 1, 1), time.Local)
	return firstTimeByYear.AddDate(1, 0, 0).Add(-1)
}

// 获取本月最后一天
func GetNowMonthLastDay() int {
	return GetMonthMaxDay(time.Now().Year(), time.Now().Month())
}

// 判断本月是否包含某一天
func NowMonthHasDay(day int) bool {
	return GetMonthMaxDay(time.Now().Year(), time.Now().Month()) >= day
}

// 判断指定月份是否包含某一天
func MonthHasDay(month, day int) bool {
	return GetMonthMaxDay(time.Now().Year(), time.Month(month)) >= day
}

func SetHourMinAndSecond(t time.Time, hour, min, sec int) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, hour, min, sec, 0, t.Location())
}
