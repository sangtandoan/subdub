package middlewares

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

// Write implements the http.ResponseWriter interface
func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

// WriteString implements the http.ResponseWriter interface for strings
func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func GZipMiddleware(c *gin.Context) {
	// checks if client supports gzip
	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Next()
		return
	}

	// checks if we want to compress this content type
	contentType := c.Writer.Header().Get("Content-Type")
	if strings.Contains(contentType, "image/") || strings.Contains(contentType, "video/") {
		c.Next()
		return
	}

	// setup gzip writer
	gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestCompression)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	defer gz.Close()

	// set appropriate headers
	c.Header("Content-Encoding", "gzip")
	c.Header("Vary", "Accept-Encoding")

	c.Writer = &gzipWriter{
		ResponseWriter: c.Writer,
		writer:         gz,
	}

	c.Next()
}
