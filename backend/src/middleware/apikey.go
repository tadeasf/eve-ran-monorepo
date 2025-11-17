package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth middleware validates the X-API-Key header
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := os.Getenv("API_KEY")

		// If API_KEY is not set, skip validation (for development)
		if apiKey == "" {
			c.Next()
			return
		}

		// Get the API key from the request header
		providedKey := c.GetHeader("X-API-Key")

		// Check if the API key matches
		if providedKey != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or missing API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAPIKeyAuth middleware that allows requests with or without API key
// but logs when API key is missing
func OptionalAPIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := os.Getenv("API_KEY")

		// If API_KEY is not set, skip validation
		if apiKey == "" {
			c.Next()
			return
		}

		providedKey := c.GetHeader("X-API-Key")

		// Set a flag in context to indicate if API key is valid
		if providedKey == apiKey {
			c.Set("api_key_valid", true)
		} else {
			c.Set("api_key_valid", false)
		}

		c.Next()
	}
}
