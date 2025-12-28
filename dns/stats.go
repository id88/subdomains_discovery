package dns

import (
	"sync"
	"time"
)

// Stat DNS性能统计
type Stat struct {
	Success int
	Failed  int
	AvgTime time.Duration
	mu      sync.RWMutex
}

// RecordSuccess 记录成功查询
func (s *Stat) RecordSuccess(duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Success++
	total := s.Success + s.Failed
	s.AvgTime = (s.AvgTime*time.Duration(total-1) + duration) / time.Duration(total)
}

// RecordFailure 记录失败查询
func (s *Stat) RecordFailure(duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Failed++
	total := s.Success + s.Failed
	s.AvgTime = (s.AvgTime*time.Duration(total-1) + duration) / time.Duration(total)
}

// GetStats 获取统计信息
func (s *Stat) GetStats() (success, failed int, avgTime time.Duration) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Success, s.Failed, s.AvgTime
}
