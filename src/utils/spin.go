package utils

import "time"

func Spin(duration time.Duration) {
	startTime := time.Now()
	for {
		if time.Since(startTime) > duration {
			return
		}
	}
}
