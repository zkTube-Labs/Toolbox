package helper

import "time"

func TimeDifferenceDays(t1, t2 time.Time) uint64 {
	return uint64(t1.Sub(t2).Hours()) / 24
}
