package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/utkarshsudhakar/demo_app/config"
)

func IntVal(str string) int64 {

	intval, _ := strconv.ParseInt(str, 10, 64)
	return intval
}

func EndToEndTime(jobData config.JobResponse) string {

	fmt.Println(jobData)
	smaxTime := jobData[0].EndTime
	sminTime := jobData[0].StartTime
	maxTime, _ := strconv.ParseInt(smaxTime, 10, 64)
	minTime, _ := strconv.ParseInt(sminTime, 10, 64)
	for i := 1; i < len(jobData); i++ {

		if jobData[i].Type != "Purge" {
		if IntVal(jobData[i].EndTime) > maxTime {
			maxTime = IntVal(jobData[i].EndTime)
			//fmt.Println(maxTime)
		}

		if IntVal(jobData[i].StartTime) < minTime {
			minTime = IntVal(jobData[i].StartTime)
			//fmt.Println(minTime)
		}
	}
	}

	startTime := time.Unix(0, minTime*int64(time.Millisecond))
	endTime := time.Unix(0, maxTime*int64(time.Millisecond))

	//fmt.Println(endTime)
	//fmt.Println(startTime)
	diff := endTime.Sub(startTime)
	//fmt.Printf("%f", diff.Seconds())
	//p := fmt.Sprintf("%02d:%02d:%02d", int64(diff.Hours()), int64(diff.Minutes()), int64(diff.Seconds()))
	sdiff := SecToTime(int64(diff.Seconds()))

	return sdiff
}
