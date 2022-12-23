package util

import (
	"blogrpc/proto/common/types"
	"strconv"
	"strings"
	"time"
)

func MonthRange(t time.Time) (time.Time, time.Time) {
	year, month, _ := t.Date()
	thisMonth := time.Date(year, month, 0, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, 0, 1)
	end := thisMonth.AddDate(0, 1, -1).Add(-time.Second * 1)
	return start, end
}

func GetTodayDateRange(rangeType types.RangeType) *types.StringDateRange {
	t := time.Now()
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 0, 1)
	return &types.StringDateRange{
		Start: start.Format(RFC3339Mili),
		End:   end.Format(RFC3339Mili),
		Type:  rangeType,
	}
}

func GetYesterdayZeroTime() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, -1)
}

func GetTomorrowEndTime() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).AddDate(0, 0, 1)
}

func GetLatestTime(times []time.Time) time.Time {
	if len(times) == 0 {
		panic("times should not be empty")
	}
	t := times[0]
	for i := range times {
		if times[i].After(t) {
			t = times[i]
		}
	}
	return t
}

// 获取至今的时间范围
func GetToTodayDateRange(rangeType types.RangeType) *types.StringDateRange {
	t := time.Now()
	end := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	start := time.Time{}
	return &types.StringDateRange{
		Start: start.Format(RFC3339Mili),
		End:   end.Format(RFC3339Mili),
		Type:  rangeType,
	}
}

func SplitTimeRangeByDuration(startTime time.Time, endTime time.Time, duration time.Duration) []*types.DateRange {
	dateRangeList := []*types.DateRange{}
	startTimestamp := startTime.Unix()
	endTimestamp := endTime.Unix()
	if endTimestamp > startTimestamp {
		durationSeconds := int64(duration.Seconds())
		if endTimestamp-startTimestamp <= durationSeconds {
			dateRangeList = append(dateRangeList, &types.DateRange{
				Start: startTimestamp,
				End:   endTimestamp,
				Type:  types.RangeType_CLOSE_OPEN,
			})
		} else {
			currentRangeStart := startTimestamp
			for {
				currentRangeEnd := currentRangeStart + durationSeconds
				if currentRangeEnd < endTimestamp {
					dateRangeList = append(dateRangeList, &types.DateRange{
						Start: currentRangeStart,
						End:   currentRangeEnd,
						Type:  types.RangeType_CLOSE_OPEN,
					})
				} else {
					dateRangeList = append(dateRangeList, &types.DateRange{
						Start: currentRangeStart,
						End:   endTimestamp,
						Type:  types.RangeType_CLOSE_CLOSE,
					})
					break
				}
				currentRangeStart = currentRangeEnd
			}
		}
	}
	return dateRangeList
}

func FormatClockTimeStr(clockTime string, delimiter string) (int, int) {
	clockTimes := strings.Split(clockTime, delimiter)
	hours, _ := strconv.Atoi(clockTimes[0])
	minutes, _ := strconv.Atoi(clockTimes[1])
	return hours, minutes
}

func FormatTimeToRFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
