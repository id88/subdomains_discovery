package utils

import (
	"fmt"
	"sync"
	"time"
)

// ProgressBar 进度条结构体
type ProgressBar struct {
	Total     int
	Current   int
	StartTime time.Time
	mu        sync.Mutex
}

// NewProgressBar 创建新的进度条
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		Total:     total,
		Current:   0,
		StartTime: time.Now(),
	}
}

// Increment 增加进度
func (p *ProgressBar) Increment() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Current++
}

// Display 显示进度
func (p *ProgressBar) Display() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Total == 0 {
		return
	}

	percentage := float64(p.Current) / float64(p.Total) * 100
	elapsed := time.Since(p.StartTime)

	// 计算预计剩余时间
	var eta time.Duration
	if p.Current > 0 {
		avgTime := elapsed / time.Duration(p.Current)
		remaining := p.Total - p.Current
		eta = avgTime * time.Duration(remaining)
	}

	// 创建进度条
	barWidth := 50
	filled := int(float64(barWidth) * percentage / 100)
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "="
		} else if i == filled {
			bar += ">"
		} else {
			bar += " "
		}
	}

	fmt.Printf("\r[%s] %.2f%% (%d/%d) 耗时: %v 预计剩余: %v",
		bar, percentage, p.Current, p.Total, elapsed.Round(time.Second), eta.Round(time.Second))
}

// Finish 完成进度条
func (p *ProgressBar) Finish() {
	p.mu.Lock()
	p.Current = p.Total
	p.mu.Unlock()
	p.Display()
	fmt.Println()
}
