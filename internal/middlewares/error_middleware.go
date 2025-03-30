package middlewares

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
)

func ErrorMiddleware(c *gin.Context) {
	c.Next()

	gz, exists := c.Get("gz")
	if exists {
		defer gz.(*gzip.Writer).Close()
	} else {
		_ = c.Error(apperror.ErrInternalServerError)
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
