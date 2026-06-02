package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/tendo-mulira/tnotes-teams/internal/utils"
)

// RateLimiter implements a simple in-memory token bucket rate limiter.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // requests per window
	window   time.Duration // time window
}

type visitor struct {
	count    int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}
	// Clean up old entries periodically
	go rl.cleanup()
	return rl
}

// RateLimit creates rate limiting middleware.
func (m *Middleware) RateLimit(rate int, window time.Duration) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(rate, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			limiter.mu.Lock()
			v, exists := limiter.visitors[ip]
			if !exists || time.Since(v.lastSeen) > limiter.window {
				limiter.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
				limiter.mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if v.count >= limiter.rate {
				limiter.mu.Unlock()
				utils.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}

			v.count++
			v.lastSeen = time.Now()
			limiter.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.window * 2)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
