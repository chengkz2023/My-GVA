// Package ratelimit 提供单实例内存限流器，用于登录防爆破等场景。
// 多实例部署时，客户项目应替换为基于共享存储（如 Redis）的实现。
package ratelimit

import (
	"sync"
	"time"
)

// Limiter 固定窗口计数限流器。
type Limiter struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	entries map[string]*windowEntry
}

type windowEntry struct {
	start time.Time
	count int
}

// NewLimiter 创建限流器；max<=0 时取 100，window<=0 时取 1 小时。
func NewLimiter(max int, window time.Duration) *Limiter {
	if max <= 0 {
		max = 100
	}
	if window <= 0 {
		window = time.Hour
	}
	return &Limiter{max: max, window: window, entries: make(map[string]*windowEntry)}
}

// Allow 报告 key 在当前窗口内是否尚未超限（不计数）。
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.entries[key]
	if !ok || time.Since(e.start) >= l.window {
		return true
	}
	return e.count < l.max
}

// Take 记录一次事件（如一次登录失败）。
func (l *Limiter) Take(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	e, ok := l.entries[key]
	if !ok || now.Sub(e.start) >= l.window {
		l.maybePurgeLocked(now)
		l.entries[key] = &windowEntry{start: now, count: 1}
		return
	}
	e.count++
}

// Reset 清除 key 的计数（如登录成功）。
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key)
}

// maybePurgeLocked 在条目过多时清理过期项（调用方需持锁）。
func (l *Limiter) maybePurgeLocked(now time.Time) {
	if len(l.entries) < 10000 {
		return
	}
	for k, e := range l.entries {
		if now.Sub(e.start) >= l.window {
			delete(l.entries, k)
		}
	}
}
