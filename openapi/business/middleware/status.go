package middleware

import (
	"fmt"
	"os"
	"sync"
	"time"

	"blogrpc/core/extension"

	"github.com/gin-gonic/gin"
)

// StatusMiddleware keeps track of various stats about the processed requests.
// It depends on context "STATUS_CODE" and context "ELAPSED_TIME",
type StatusMiddleware struct {
	lock              sync.RWMutex
	start             time.Time
	pid               int
	responseCounts    map[string]int
	totalResponseTime time.Time
}

// MiddlewareFunc makes StatusMiddleware implement the Middleware interface.
func (mw *StatusMiddleware) MiddlewareFunc() gin.HandlerFunc {
	mw.start = time.Now()
	mw.pid = os.Getpid()
	mw.responseCounts = map[string]int{}
	mw.totalResponseTime = time.Time{}

	return func(c *gin.Context) {
		c.Next()

		var responseTime *time.Duration
		if value, exists := c.Get("ELAPSED_TIME"); exists {
			responseTime = value.(*time.Duration)
		}

		mw.lock.Lock()
		mw.responseCounts[fmt.Sprintf("%d", c.Writer.Status())]++
		if responseTime != nil {
			mw.totalResponseTime = mw.totalResponseTime.Add(*responseTime)
		}
		mw.lock.Unlock()
	}
}

// Status contains stats and status information. It is returned by GetStatus.
// These information can be made available as an API endpoint, see the "status"
// example to install the following status route.
// GET /.status returns something like:
//
//	{
//	  "Pid": 21732,
//	  "UpTime": "1m15.926272s",
//	  "Time": "2013-03-04 08:00:27.152986 +0000 UTC",
//	  "StatusCodeCount": {
//	    "200": 53,
//	    "404": 11
//	  },
//	  "TotalCount": 64,
//	  "AverageResponseTime": "262.14us",
//	  "DBReachable": true,
//	}
type Status struct {
	Pid                 int
	UpTime              string
	Time                string
	StatusCodeCount     map[string]int
	TotalCount          int
	AverageResponseTime string
	DBReachable         bool
}

// GetStatus computes and returns a Status object based on the request informations accumulated
// since the start of the process.
func (mw *StatusMiddleware) GetStatus() *Status {

	mw.lock.RLock()

	now := time.Now()

	uptime := now.Sub(mw.start)

	totalCount := 0
	for _, count := range mw.responseCounts {
		totalCount += count
	}

	totalResponseTime := mw.totalResponseTime.Sub(time.Time{})

	averageResponseTime := time.Duration(0)

	if totalCount > 0 {
		avgNs := int64(totalResponseTime) / int64(totalCount)
		averageResponseTime = time.Duration(avgNs)
	}

	err := extension.PingMgo()

	status := &Status{
		Pid:                 mw.pid,
		UpTime:              uptime.String(),
		Time:                now.String(),
		StatusCodeCount:     mw.responseCounts,
		TotalCount:          totalCount,
		AverageResponseTime: averageResponseTime.String(),
		DBReachable:         err == nil,
	}

	mw.lock.RUnlock()

	return status
}
