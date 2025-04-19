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

	// we have 2 options for not compressing img or video

	// 1.
	// Checks if we want to compress this content type in response header
	// because we want to comporess response data,
	// not checking in request header because that means user is sending img or video to us.
	//
	// There's actually a flaw in this implementation.
	// Since the middleware runs before the handler,
	// the Content-Type might not be set yet when this check runs
	//
	// contentType := c.Writer.Header().Get("Content-Type")
	// if strings.Contains(contentType, "image/") || strings.Contains(contentType, "video/") {
	// 	c.Next()
	// 	return
	// }

	// 2.
	// Check URL path for file extensions.
	// This is a more reliable way to check if we want to compress this content type
	path := c.Request.URL.Path
	if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".mp4") {
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

	// set appropriate headers
	c.Header("Content-Encoding", "gzip")
	c.Header("Vary", "Accept-Encoding")

	c.Writer = &gzipWriter{
		ResponseWriter: c.Writer,
		writer:         gz,
	}

	// stores the GZIP writer in the Gin context (`c`)
	// so it can be closed later (in ErrorMiddlware),
	// this is important to avoid memory leaks.
	c.Set("gz", gz)

	c.Next()
}
