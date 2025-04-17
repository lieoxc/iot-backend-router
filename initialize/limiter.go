package initialize

import (
	"sync"

	"golang.org/x/time/rate"
)

type AutomateLimiter struct {
	mu         sync.Mutex
	limiters   map[string]*rate.Limiter       // 存储不同键的限流器
	configFunc func(string) (rate.Limit, int) // 动态获取限流配置的函数
}

// NewAutomateLimiter 创建限流器实例
// configFunc: 接收键名，返回该键对应的速率限制和桶容量
func NewAutomateLimiter(configFunc func(string) (rate.Limit, int)) *AutomateLimiter {
	return &AutomateLimiter{
		limiters:   make(map[string]*rate.Limiter),
		configFunc: configFunc,
	}
}

// GetLimiter 获取或创建指定键的限流器
func (rl *AutomateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 如果已存在直接返回
	if limiter, ok := rl.limiters[key]; ok {
		return limiter
	}

	// 动态获取限流配置
	rateLimit, burst := rl.configFunc(key)
	limiter := rate.NewLimiter(rateLimit, burst)
	rl.limiters[key] = limiter

	return limiter
}

func (rl *AutomateLimiter) Allow(key string) bool {
	limiter := rl.GetLimiter(key)
	return limiter.Allow()
}

// 默认配置函数（可选）
func DefaultConfig(key string) (rate.Limit, int) {
	return rate.Limit(1.0 / 3.0), 10 // 默认值：3秒1次，突发10
}
func AutomateRateLimitConfig(key string) (rate.Limit, int) {
	return rate.Limit(1.0 / 60.0), 1 // 默认值：60秒1次，突发1
}
