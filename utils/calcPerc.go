package utils

import (
	"time"
)

func CalcPerc(diff float64, relativeTime time.Time) float64 {
	/*	svOld := "05:00:00"
		svNew := "15:00:00"

		timeOld, _ := time.Parse(config.TimeFormat, svOld)
		timeNew, _ := time.Parse(config.TimeFormat, svNew)
		diff := timeNew.Sub(timeOld)
	*/
	relativeTimeInSecs := float64((relativeTime.Hour() * 3600) + (relativeTime.Minute() * 60) + (relativeTime.Second()))

	return float64((diff / relativeTimeInSecs) * 100)
}
