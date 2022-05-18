package errors

import (
	"fmt"
	"runtime/debug"

	loggy "github.com/fliptable-io/subscription-service/src/utils/logging"

	"github.com/gin-gonic/gin"
)

type ServerError struct {
	error
	Code    int
	Message string
}

func (se ServerError) Error() string {
	return fmt.Sprint("[ServerError]", "|", se.Code, "|", se.Message, "|", se.error)
}

var UnknownError = ServerError{Code: 400, Message: "Unknown Error"}
var NotFoundError = ServerError{Code: 404, Message: "Not Found"}
var FormatError = ServerError{Code: 400, Message: "Format Error"}

func (se ServerError) Consume(err error) error {
	if err == nil {
		return nil
	}

	switch v := err.(type) {
	case *ServerError:
		return &se
	case error:
		se.error = err
		debug.PrintStack()
		loggy.Error(v.Error())
	}

	return &se
}

func (se ServerError) IfUnknown(err error) error {
	switch v := err.(type) {
	case *ServerError:
		return v
	default:
		return se.Consume(err)
	}
}

func (se *ServerError) Respond(ctx *gin.Context) {
	ctx.JSON(se.Code, se)
}

func Respond(ctx *gin.Context, err error, placeholder *ServerError) {
	if err == nil {
		placeholder.Respond(ctx)
		return
	}

	switch v := err.(type) {
	case *ServerError:
		ctx.JSON(v.Code, v)
		return
	}

	loggy.Error(err)
	placeholder.Respond(ctx)
}

func (se ServerError) WithMessage(message string) *ServerError {
	se.Message = message
	return &se
}
