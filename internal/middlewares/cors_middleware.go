package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is in the allowed origins list
		allowedOrigin := "*"
		if len(allowedOrigins) > 0 && allowedOrigins[0] != "*" {
			for _, allowed := range allowedOrigins {
				if allowed == origin {
					allowedOrigin = allowed
					break
				}
			}
		}

		log.Info("CORS Middleware: Allowed Origin:", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			log.Info("OPTIONS request received, aborting with 204 No Content")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
