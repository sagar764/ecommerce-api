package utils

import (
	"ecommerce-api/errortools"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response structure.
// @Description Standard API response structure used in the application
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
type APIResponse struct {
	Status  string            `json:"status,omitempty"`
	Message string            `json:"message,omitempty"`
	Data    map[string]any    `json:"data,omitempty"`
	Error   *errortools.Error `json:"error,omitempty"`
}

func (resp *APIResponse) AddData(key string, data any) {
	if resp.Data == nil {
		resp.Data = make(map[string]any)
	}
	resp.Data[key] = data
}

func SuccessGenerator(data any, key, msg string) APIResponse {
	res := APIResponse{
		Status:  "success",
		Message: msg,
	}
	if data != nil {
		res.Data = map[string]any{
			key: data,
		}
	}

	return res
}

func ErrorGenerator(ctx *gin.Context, err error) {
	res := APIResponse{
		Status: "failure",
	}

	var srvErr *errortools.Error
	if errors.As(err, &srvErr) {
		res.Error = srvErr
		code := srvErr.Code.HTTPCode()
		if code == 0 {
			code = http.StatusBadRequest
		}
		ctx.AbortWithStatusJSON(code, res)
	} else {
		// Handle non-errortools errors
		res.Error = errortools.New("internal_error", errortools.WithDetail(err.Error()))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
	}
}

// Function to handle binding errors
func BindingError(ctx *gin.Context, err error) {
	var (
		unmarshalTypeErr *json.UnmarshalTypeError
		syntaxErr        *json.SyntaxError
	)
	resErr := errortools.Init()
	switch {
	case errors.As(err, &unmarshalTypeErr):
		field := unmarshalTypeErr.Field
		expectedType := unmarshalTypeErr.Type
		actualValue := unmarshalTypeErr.Value
		resErr.AddValidationError(field, errortools.InvalidFieldType, expectedType, actualValue)
	case errors.As(err, &syntaxErr):
		resErr = errortools.New(errortools.BindingError,
			errortools.WithDetail(fmt.Sprintf("%s : %d", syntaxErr, syntaxErr.Offset)))
	default:
		resErr = errortools.New(errortools.BindingError, errortools.WithDetail(err.Error()))
	}
	ErrorGenerator(ctx, resErr)
}
