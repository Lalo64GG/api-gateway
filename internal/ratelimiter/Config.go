package ratelimiter

import "time"

type Config struct {
	MaxRequests int
	RefillRate float64
	CleanupInterval time.Duration
	IPTimeout time.Duration
}

func DefaultConfig() Config{
	return Config{
		MaxRequests: 10,  // Maximum requests in bucket
		RefillRate: 0.1667, // 1 request every 6 seconds (10 request/minute)
		CleanupInterval: 5 * time.Minute, // Cleanup interval 5 minutes
		IPTimeout: 10 * time.Minute, // Delete IP after 10 minutes of inactivity
	}
}