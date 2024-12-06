package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

// Client stores the limiter and last access for each IP
type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter interface {
	LimitRate() gin.HandlerFunc
}

// RateLimiterImpl handles the rate limiting for clients
type RateLimiterImpl struct {
	clients    map[string]*Client
	mu         sync.RWMutex
	rate       rate.Limit
	burst      int
	cleanupInt time.Duration
}

// NewRateLimiterImpl creates a new rate limiter
func NewRateLimiterImpl(r rate.Limit, burst int) *RateLimiterImpl {
	rl := &RateLimiterImpl{
		clients:    make(map[string]*Client),
		rate:       r,
		burst:      burst,
		cleanupInt: time.Hour,
	}

	// Cleanup old clients periodically
	go rl.cleanup()

	return rl
}

// getClient fetches or creates a client
func (rl *RateLimiterImpl) getClient(ip string) *Client {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if client, exists := rl.clients[ip]; exists {
		client.lastSeen = time.Now()
		return client
	}

	client := &Client{
		limiter:  rate.NewLimiter(rl.rate, rl.burst),
		lastSeen: time.Now(),
	}
	rl.clients[ip] = client

	return client
}

// cleanup deletes inactive clients
func (rl *RateLimiterImpl) cleanup() {
	for {
		time.Sleep(rl.cleanupInt)

		rl.mu.Lock()
		for ip, client := range rl.clients {
			if time.Since(client.lastSeen) > rl.cleanupInt {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware returns the Gin middleware for the rate limiting
func (rl *RateLimiterImpl) LimitRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client's IP
		ip := c.ClientIP()

		// Get or create the limiter for this IP
		client := rl.getClient(ip)

		// Check if request is allowed
		if !client.limiter.Allow() {
			c.JSON(429, gin.H{
				"error":       "Rate limit exceeded",
				"retry-after": time.Until(time.Now().Add(time.Second)),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
