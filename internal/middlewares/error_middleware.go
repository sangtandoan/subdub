package middlewares

import (
	"compress/gzip"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
)

func ErrorMiddleware(c *gin.Context) {
	c.Next()

	// close gzip writer if it ok
	gz, ok := c.Get("gz")
	if ok {
		defer gz.(*gzip.Writer).Close()
	}

	if len(c.Errors) > 0 {
		for _, err := range c.Errors {
			var appError *apperror.AppError

			switch e := err.Err.(type) {
			case *apperror.AppError:
				appError = e
			case validator.ValidationErrors:
				appError = apperror.HandleValidateErrors(e)
			default:
				appError = handleDefaultError(e)
			}

			fmt.Println(appError.Msg)

			c.JSON(appError.StatusCode, gin.H{"success": appError.Success, "errors": appError.Msg})
			return
		}
	}
}

func handleDefaultError(err error) *apperror.AppError {
	if strings.Contains(err.Error(), "token is expired") {
		return apperror.ErrTokenExpired
	}

	return apperror.ErrInternalServerError
}
