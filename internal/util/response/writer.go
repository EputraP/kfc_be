package response

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func JSON(ctx *gin.Context, statusCode int, message string, data interface{}) {
	resp := NewResponse(statusCode, message, data)

	ctx.JSON(statusCode, resp)
}

func Error(ctx *gin.Context, statusCode int, message string) {
	JSON(ctx, statusCode, message, nil)
	ctx.Abort()
}

func ValidationError(ctx *gin.Context, errs validator.ValidationErrors) {
	err := errs[0]

	msg := fmt.Sprintf("validation error: field: [%v] rule: [%v] ", err.Field(), err.ActualTag())

	if err.Param() != "" {
		msg += fmt.Sprintf("rule value: [%v] ", err.Param())
	}

	msg += fmt.Sprintf("actual value: [%v]", err.Value())

	Error(ctx, http.StatusBadRequest, msg)
}

func UnknownError(ctx *gin.Context, err error) {
	log.Println(err)
	Error(ctx, http.StatusInternalServerError, "Something went wrong")
}

func FromRequest(ctx *gin.Context, response *http.Response, message string, data interface{}) {
	for key := range response.Header {
		ctx.Header(key, response.Header.Get(key))
		ctx.Status(http.StatusOK)
	}
}
