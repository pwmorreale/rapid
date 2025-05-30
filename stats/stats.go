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
	count     atomic.Uint64
	errors    atomic.Uint64
	totalTime atomic.Int64
	minTime   atomic.Int64
	maxTime   atomic.Int64
}

func (s *Statistics) setMin(d int64) {
	for {
		current := s.minTime.Load()
		if current > 0 && d >= current {
			return
		}
		if s.minTime.CompareAndSwap(current, d) {
			return
		}
	}
}

func (s *Statistics) setMax(d int64) {
	for {
		current := s.maxTime.Load()
		if d <= current {
			return
		}
		if s.maxTime.CompareAndSwap(current, d) {
			return
		}
	}
}

func (s *Statistics) updateTimes(start time.Time) {

	d := int64(time.Since(start))
	s.totalTime.Add(d)
	s.setMin(d)
	s.setMax(d)

}

// Success updates statistics and the recorded execution times.
func (s *Statistics) Success(start time.Time) {
	s.count.Add(1)
	s.updateTimes(start)
}

// Error updates statistics and the recorded execution times.
func (s *Statistics) Error(start time.Time) {
	s.errors.Add(1)
	s.updateTimes(start)
}

// GetCount returns the count.
func (s *Statistics) GetCount() int64 {
	return int64(s.count.Load())
}

// GetErrors returns the errors.
func (s *Statistics) GetErrors() int64 {
	return int64(s.errors.Load())
}

// GetMinDuration returns the minimum duration.
func (s *Statistics) GetMinDuration() time.Duration {
	return time.Duration(s.minTime.Load())
}

// GetDuration returns the total time duration.
func (s *Statistics) GetDuration() time.Duration {
	return time.Duration(s.totalTime.Load())
}

// GetMaxDuration returns the maximum duration.
func (s *Statistics) GetMaxDuration() time.Duration {
	return time.Duration(s.maxTime.Load())
}
