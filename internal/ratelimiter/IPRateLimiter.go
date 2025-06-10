package ratelimiter

import (
	"log"
	"sync"
	"time"
)

type IPRateLimiter struct {
	limiters map[string]*RateLimiter
	lastAccess map[string]time.Time
	config Config
	mutex    sync.RWMutex
	stopChan chan struct{}
	once sync.Once
}

func NewIPRateLimiter(config Config) *IPRateLimiter{
	ipr := &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
		lastAccess: make(map[string]time.Time),
		config: config,
		stopChan: make(chan struct{}),
	}

	ipr.once.Do(func(){
		go ipr.startCleanup()
	})

	return ipr
}


func (i *IPRateLimiter) GetLimiter(ip string) *RateLimiter {
	i.mutex.RLock()
	defer i.mutex.Unlock()

	limiter, exists := i.limiters[ip]
	if exists {
		i.lastAccess[ip] = time.Now()
		i.mutex.RUnlock()
		return limiter
	}
	i.mutex.RUnlock()

	i.mutex.Lock()
	defer i.mutex.Unlock()

	limiter, exists = i.limiters[ip]
	if exists {
		i.lastAccess[ip] = time.Now()
		return limiter
	}

	limiter = NewRateLimiter(
		float64(i.config.MaxRequests),
		i.config.RefillRate,
	)

	return limiter
}

func (i *IPRateLimiter) startCleanup() {
	ticker := time.NewTicker(i.config.CleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			i.cleanup()
		case <-i.stopChan:
			log.Println("Rate limiter cleanup stopped")
			return
		}
	}
}

func (i *IPRateLimiter) cleanup() {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	now := time.Now()
	removedCount := 0

	for ip, lastAccess := range i.lastAccess {
		if now.Sub(lastAccess) > i.config.IPTimeout {
			delete(i.limiters, ip)
			delete(i.lastAccess, ip)
			removedCount++
		}
	}

	if removedCount > 0{
		log.Printf("Rate limiter cleanup: removed %d inactive IPs, remaining: %d",
	removedCount, len(i.limiters))
	}
}

func (i *IPRateLimiter) Stop() {
	close(i.stopChan)
}

func (i *IPRateLimiter) Stats() (activeIPs int, totalRequest int64) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	return len(i.limiters), int64(len(i.limiters))
}