package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
)

func ErrorMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		for _, err := range c.Errors {
			var appError *apperror.AppError

			switch e := err.Err.(type) {
			case *apperror.AppError:
				appError = e
			case validator.ValidationErrors:
				appError = apperror.HandleValidateErrors(e)
			default:
				appError = apperror.ErrInternalServerError
			}

			c.JSON(appError.StatusCode, gin.H{"error": appError.Msg})
			return
		}
	}
}
