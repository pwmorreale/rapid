//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package stats implements RAPID statistics.
package stats

import (
	"sync/atomic"
	"time"
)

// Statistics defines measured statistics.
type Statistics struct {
	count     int64
	errors    int64
	totalTime int64
	minTime   int64
	maxTime   int64
}

func (s *Statistics) setMin(d int64) {
	for {
		current := atomic.LoadInt64(&s.minTime)
		if current > 0 && d >= current {
			return
		}
		if atomic.CompareAndSwapInt64(&s.minTime, current, d) {
			return
		}
	}
}

func (s *Statistics) setMax(d int64) {
	for {
		current := atomic.LoadInt64(&s.maxTime)
		if d <= current {
			return
		}
		if atomic.CompareAndSwapInt64(&s.maxTime, current, d) {
			return
		}
	}
}

func (s *Statistics) updateTimes(start time.Time) {

	d := int64(time.Since(start))
	atomic.AddInt64(&s.totalTime, d)
	s.setMin(d)
	s.setMax(d)

}

// Success updates statistics and the recorded execution times.
func (s *Statistics) Success(start time.Time) {
	atomic.AddInt64(&s.count, 1)
	s.updateTimes(start)
}

// Error updates statistics and the recorded execution times.
func (s *Statistics) Error(start time.Time) {
	atomic.AddInt64(&s.errors, 1)
	s.updateTimes(start)
}

// GetCount returns the count.
func (s *Statistics) GetCount() int64 {
	return atomic.LoadInt64(&s.count)
}

// GetErrors returns the errors.
func (s *Statistics) GetErrors() int64 {
	return atomic.LoadInt64(&s.errors)
}

// GetMinDuration returns the minimum duration.
func (s *Statistics) GetMinDuration() time.Duration {
	return time.Duration(atomic.LoadInt64(&s.minTime))
}

// GetDuration returns the total time duration.
func (s *Statistics) GetDuration() time.Duration {
	return time.Duration(atomic.LoadInt64(&s.totalTime))
}

// GetMaxDuration returns the maximum duration.
func (s *Statistics) GetMaxDuration() time.Duration {
	return time.Duration(atomic.LoadInt64(&s.maxTime))
}
