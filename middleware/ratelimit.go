package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/abrhamyalew/telebirr-payment-verifier/models"
)

type ipLimiter struct {
	tokens     float64
	lastUpdate time.Time
}

type RateLimiter struct {
	mu         sync.Mutex
	clients    map[string]*ipLimiter
	rate       float64 // tokens per second
	capacity   float64 // max tokens
}

func NewRateLimiter(rate float64, capacity float64) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*ipLimiter),
		rate:     rate,
		capacity: capacity,
	}

	// Periodic cleanup of stale clients (every 5 minutes)
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.mu.Lock()
			now := time.Now()
			for ip, lim := range rl.clients {
				if now.Sub(lim.lastUpdate) > 10*time.Minute {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	lim, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &ipLimiter{
			tokens:     rl.capacity - 1,
			lastUpdate: now,
		}
		return true
	}

	elapsed := now.Sub(lim.lastUpdate).Seconds()
	lim.tokens += elapsed * rl.rate
	if lim.tokens > rl.capacity {
		lim.tokens = rl.capacity
	}
	lim.lastUpdate = now

	if lim.tokens >= 1 {
		lim.tokens -= 1
		return true
	}

	return false
}

func RateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerSecond, float64(burst))
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Error: "Too many requests. Rate limit exceeded.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
