package domain

import "time"

func readiness(from, to, now time.Time) float64 {
	if from.After(to) {
		return 0.0
	}
	if from.After(now) {
		return 0.0
	}
	if to.Before(now) {
		return 1.0
	}
	totalDiff := to.Sub(from)
	curDiff := to.Sub(now)
	if curDiff <= 0 {
		return 1.0
	}
	return 1.0 - float64(curDiff)/float64(totalDiff)
}
