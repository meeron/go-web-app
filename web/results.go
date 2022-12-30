package web

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	ErrorCode string `json:"errorCode" example:"NotFound"`
	Message   string `json:"message" example:"A message"`
}

func NotFound() Error {
	return Error{
		ErrorCode: "NotFound",
	}
}

func Exists() Error {
	return Error{
		ErrorCode: "Exists",
	}
}

func BadRequest(err error) Error {
	var ve validator.ValidationErrors

	if !errors.As(err, &ve) {
		return Error{
			ErrorCode: "BadRequest",
			Message:   "invalid request",
		}
	}

	fieldName := ve[0].Field()
	message := func(tag string) string {
		switch tag {
		case "required":
			return fmt.Sprintf("'%s' is required", fieldName)
		case "gt":
			return fmt.Sprintf("'%s' must be greater than %s", fieldName, ve[0].Param())
		case "email":
			return fmt.Sprintf("'%s' is invalid", fieldName)
		}
		return "Unknown error"
	}(ve[0].Tag())

	return Error{
		ErrorCode: "BadRequest",
		Message:   message,
	}
}
