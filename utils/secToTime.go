package utils

import (
	"fmt"
)

func SecToTime(seconds int64) string {

	hours := seconds / 3600
	seconds = seconds % (60 * 60)
	minutes := seconds / 60
	seconds = seconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

}
