package middleware

import (
	"airbnb-clone/internal/cache"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiterConfig holds configuration for rate limiting
type RateLimiterConfig struct {
	RequestsPerMinute int           // Number of requests allowed per minute
	WindowSize        time.Duration // Time window for rate limiting
	KeyPrefix         string        // Redis key prefix
}

// DefaultRateLimiterConfig returns default rate limiter configuration
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerMinute: 100,
		WindowSize:        time.Minute,
		KeyPrefix:         "rate_limit",
	}
}

func RateLimiterMiddleware(redisClient *cache.RedisClient, config RateLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get client IP
		clientIP := c.ClientIP()

		// create Redis key for this IP
		key := fmt.Sprintf("%s:%s", config.KeyPrefix, clientIP)

		// get current count
		currentCount, err := redisClient.Get(key)
		if err != nil && err != redis.Nil {
			// if Redis is down it allows the request but logs the error
			c.Header("X-RateLimit-Error", "Redis unavailable")
			c.Next()
			return
		}

		var count int64
		if err == redis.Nil {
			// key doesn't exist, this is the first request
			count = 0
		} else {
			count, _ = strconv.ParseInt(currentCount, 10, 64)
		}

		// Check if limit exceeded
		if count >= int64(config.RequestsPerMinute) {
			// Get TTL for rate limit reset time
			ttl, err := redisClient.TTL(key)
			if err != nil {
				ttl = config.WindowSize
			}

			c.Header("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))
			c.Header("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     fmt.Sprintf("Too many requests. Limit: %d requests per minute", config.RequestsPerMinute),
				"retry_after": int64(ttl.Seconds()),
			})
			c.Abort()
			return
		}

		// Increment counter
		newCount, err := redisClient.Incr(key)
		if err != nil {
			// If Redis is down, allow the request but log the error
			c.Header("X-RateLimit-Error", "Redis unavailable")
			c.Next()
			return
		}

		// Set expiration only for the first request
		if newCount == 1 {
			err = redisClient.Expire(key, config.WindowSize)
			if err != nil {
				// Log error but continue
				c.Header("X-RateLimit-Error", "Failed to set expiration")
			}
		}

		// Add rate limit headers
		remaining := config.RequestsPerMinute - int(newCount)
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		// Get TTL for reset time
		ttl, err := redisClient.TTL(key)
		if err == nil {
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))
		}

		c.Next()
	}
}

// CreateRateLimiterForEndpoint creates a rate limiter for specific endpoints
func CreateRateLimiterForEndpoint(redisClient *cache.RedisClient, requestsPerMinute int, endpoint string) gin.HandlerFunc {
	config := RateLimiterConfig{
		RequestsPerMinute: requestsPerMinute,
		WindowSize:        time.Minute,
		KeyPrefix:         fmt.Sprintf("rate_limit:%s", endpoint),
	}

	return RateLimiterMiddleware(redisClient, config)
}

// StrictRateLimiterMiddleware creates a stricter rate limiter for sensitive endpoints
func StrictRateLimiterMiddleware(redisClient *cache.RedisClient) gin.HandlerFunc {
	config := RateLimiterConfig{
		RequestsPerMinute: 10,
		WindowSize:        time.Minute,
		KeyPrefix:         "strict_rate_limit",
	}

	return RateLimiterMiddleware(redisClient, config)
}
