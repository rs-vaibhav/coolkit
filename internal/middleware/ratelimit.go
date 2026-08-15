package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/coolkit-org/coolkit/pkg/response"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests allowed per window
	window   time.Duration // time window for rate limiting
}

// Visitor tracks request information for a single client
type Visitor struct {
	lastAccess time.Time
	count      int
}

// NewRateLimiter creates a new rate limiter with specified rate and window
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		window:   window,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// cleanup removes old visitors periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastAccess) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// isAllowed checks if a visitor is within rate limits
func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[ip]
	now := time.Now()

	if !exists || now.Sub(visitor.lastAccess) > rl.window {
		rl.visitors[ip] = &Visitor{
			lastAccess: now,
			count:      1,
		}
		return true
	}

	if visitor.count >= rl.rate {
		return false
	}

	visitor.count++
	visitor.lastAccess = now
	return true
}

// RateLimit creates a rate limiting middleware
func RateLimit(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.isAllowed(ip) {
			response.TooManyRequests(c, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}
