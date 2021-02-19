package utils

import (
	"strconv"
	"time"

	"github.com/utkarshsudhakar/PerfAPM/config"
)

func IntVal(str string) int64 {

	intval, _ := strconv.ParseInt(str, 10, 64)
	return intval
}

func EndToEndTime(jobData config.JobResponse) string {

	//fmt.Println(jobData)

	//fmt.Println(jobData[0].Type)
	//fmt.Println(jobData[0].EndTime)
	//sminTime := jobData[0].StartTime
	smaxTime := "1513561392896"
	sminTime := "1713561392896"

	maxTime, _ := strconv.ParseInt(smaxTime, 10, 64)
	minTime, _ := strconv.ParseInt(sminTime, 10, 64)
	//fmt.Println(len(jobData))
	for i := 0; i < len(jobData); i++ {

		if (jobData[i].Type != "Purge") && (jobData[i].Type != "Delete") && (jobData[i].Type != "Configuration purge") {
			//	fmt.Println(jobData[i].Type)
			//	fmt.Println(jobData[i].EndTime)
			//	fmt.Println(maxTime)
			//jobData[i].EndTime
			if IntVal(jobData[i].EndTime) > maxTime {
				maxTime = IntVal(jobData[i].EndTime)
				//	fmt.Println("hi")

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
