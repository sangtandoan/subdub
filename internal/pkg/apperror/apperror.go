package apperror

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidJSON         = NewAppError(http.StatusBadRequest, "invalid json format")
	ErrInternalServerError = NewAppError(
		http.StatusInternalServerError,
		"an unexpected error occured",
	)
	ErrExisted = NewAppError(http.StatusBadRequest, "resource has already existed")
)

type AppError struct {
	Msg        any
	StatusCode int
	Success    bool
}

type ValidateError struct {
	Field string `json:"field,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("app error code: %d", e.StatusCode)
}

func NewAppError(statusCode int, msg any) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Msg:        msg,
		Success:    false,
	}
}

func HandleValidateErrors(errors validator.ValidationErrors) *AppError {
	var errorArr []ValidateError
	for _, err := range errors {
		var validateError ValidateError
		validateError.Field = strings.ToLower(err.Field())
		validateError.Msg = getErrMsg(err)

		errorArr = append(errorArr, validateError)
	}

	return &AppError{
		StatusCode: http.StatusBadRequest,
		Msg:        errorArr,
	}
}

func getErrMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "min":
		return "should be at least " + err.Param() + " length"
	case "max":
		return "should be at most " + err.Param() + " length"
	case "gte":
		return "should be greater than or equal to " + err.Param()
	default:
		return "invalid value"
	}
}
