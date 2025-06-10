package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	tokens  float64 // Current number of tokens
	maxTokens float64 // Maximum tokens allowed
	refillRate float64 // tokens added per second
	lastRefillTime time.Time // Last time tokens were refilled
	mutex sync.Mutex
}

func NewRateLimiter(maxTokens, refillRate float64) *RateLimiter {
	return &RateLimiter{
		tokens: maxTokens,
		maxTokens: maxTokens,
		refillRate: refillRate,
		lastRefillTime: time.Now(),
	}
}

func (r *RateLimiter) refillTokens() {
	now := time.Now()

	duration := now.Sub(r.lastRefillTime).Seconds()
	tokensToadd := duration * r.refillRate

	r.tokens += tokensToadd

	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}

	r.lastRefillTime = now
}

func (r *RateLimiter) Allow() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.refillTokens()

	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}