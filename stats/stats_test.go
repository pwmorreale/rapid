//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

package stats

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSuccessUpdatesCountAndTimes(t *testing.T) {

	var s Statistics

	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	s.Success(start)

	assert.Equal(t, int64(1), s.GetCount())
	assert.Equal(t, int64(0), s.GetErrors())
	assert.GreaterOrEqual(t, s.GetMinDuration(), 10*time.Millisecond)
	assert.GreaterOrEqual(t, s.GetMaxDuration(), 10*time.Millisecond)
	assert.GreaterOrEqual(t, s.GetDuration(), 10*time.Millisecond)
}

func TestErrorUpdatesErrorsAndTimes(t *testing.T) {

	var s Statistics

	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	s.Error(start)

	assert.Equal(t, int64(0), s.GetCount())
	assert.Equal(t, int64(1), s.GetErrors())
	assert.GreaterOrEqual(t, s.GetMinDuration(), 10*time.Millisecond)
}

func TestMinMaxTracking(t *testing.T) {

	var s Statistics

	// Record a short duration.
	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	s.Success(start)

	// Record a longer duration.
	start = time.Now()
	time.Sleep(50 * time.Millisecond)
	s.Success(start)

	assert.Equal(t, int64(2), s.GetCount())
	assert.LessOrEqual(t, s.GetMinDuration(), 30*time.Millisecond)
	assert.GreaterOrEqual(t, s.GetMaxDuration(), 50*time.Millisecond)
}

func TestStringZeroCount(t *testing.T) {

	var s Statistics

	str := s.String()
	assert.Contains(t, str, "count=0")
	assert.Contains(t, str, "errors=0")
}

func TestStringWithData(t *testing.T) {

	var s Statistics

	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	s.Success(start)
	s.Error(time.Now())

	str := s.String()
	assert.Contains(t, str, "count=1")
	assert.Contains(t, str, "errors=1")
	assert.Contains(t, str, "minTime=")
	assert.Contains(t, str, "maxTime=")
	assert.Contains(t, str, "avgTime=")
}

func TestConcurrentAccess(t *testing.T) {

	var s Statistics
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			s.Success(time.Now())
		}()
		go func() {
			defer wg.Done()
			s.Error(time.Now())
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(100), s.GetCount())
	assert.Equal(t, int64(100), s.GetErrors())
}
