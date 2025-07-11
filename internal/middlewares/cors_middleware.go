package middlewares

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Determine if we should set CORS headers
		var allowedOrigin string
		isAllowed := false

		// If no origins specified or "*" is explicitly allowed, allow all
		if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "*") {
			allowedOrigin = "*"
			isAllowed = true
		} else {
			// Check if the origin is in the allowed origins list
			if slices.Contains(allowedOrigins, origin) {
				allowedOrigin = origin
				isAllowed = true
			}
		}

		// Only set CORS headers if origin is allowed
		if isAllowed {
			log.Info("CORS Middleware: Allowing origin:", allowedOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().
				Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			log.Info("CORS Middleware: Origin not allowed:", origin)
			// Don't set CORS headers - let browser handle the rejection
		}

		if c.Request.Method == "OPTIONS" {
			if isAllowed {
				log.Info("OPTIONS request received for allowed origin, returning 204")
				c.AbortWithStatus(204)
			} else {
				log.Info("OPTIONS request received for disallowed origin, returning 403")
				c.AbortWithStatus(403)
			}
			return
		}

		c.Next()
	}
}
