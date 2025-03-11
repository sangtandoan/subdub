package apperror

import "github.com/gin-gonic/gin"

func WriteError(c *gin.Context, err error) {
	if err := c.Error(err); err != nil {
		panic(err)
	}
}
