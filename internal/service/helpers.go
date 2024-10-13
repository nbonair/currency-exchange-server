package service

import "time"

func isStale(lastUpdated time.Time, threshold time.Duration) bool {
	return time.Since(lastUpdated) > threshold
}
