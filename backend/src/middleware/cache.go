package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/cache"
	"github.com/tadeasf/eve-ran/src/utils"
)

// CacheMiddleware creates a caching middleware with specified TTL
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Generate cache key from path and query parameters
		cacheKey := generateCacheKey(c.Request.URL.Path, c.Request.URL.RawQuery)

		// Try to get cached response
		cachedResponse, err := cache.Get(cacheKey)
		if err != nil {
			utils.ErrorLogger.Printf("Error reading from cache (key: %s): %v", cacheKey, err)
		}

		// If cache hit, return cached response
		if cachedResponse != "" {
			utils.InfoLogger.Printf("Cache HIT: %s", cacheKey)
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json", []byte(cachedResponse))
			c.Abort()
			return
		}

		// Cache miss - capture the response
		utils.InfoLogger.Printf("Cache MISS: %s", cacheKey)
		c.Header("X-Cache", "MISS")

		// Create a custom response writer to capture the response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Execute the handler
		c.Next()

		// Only cache successful responses
		if c.Writer.Status() == http.StatusOK {
			responseBody := blw.body.String()

			// Store in cache with TTL
			err := cache.Set(cacheKey, responseBody, ttl)
			if err != nil {
				utils.ErrorLogger.Printf("Error writing to cache (key: %s): %v", cacheKey, err)
			} else {
				utils.InfoLogger.Printf("Response cached (key: %s, size: %d bytes, ttl: %s)", cacheKey, len(responseBody), ttl)
			}
		}
	}
}

// bodyLogWriter is a custom response writer that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString captures the response body as string
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// generateCacheKey creates a unique cache key from path and query
func generateCacheKey(path, query string) string {
	key := fmt.Sprintf("%s?%s", path, query)
	hash := sha256.Sum256([]byte(key))
	return fmt.Sprintf("cache:%s", hex.EncodeToString(hash[:]))
}

// InvalidateCache deletes cache entries matching a pattern
func InvalidateCache(pattern string) error {
	// Note: This is a simple implementation
	// For production, consider using Redis SCAN with pattern matching
	return cache.Delete(pattern)
}

// InvalidateAllCache clears all cache entries
func InvalidateAllCache() error {
	utils.InfoLogger.Println("Invalidating all cache entries")
	return cache.FlushAll()
}
