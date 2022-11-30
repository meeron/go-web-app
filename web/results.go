package web

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

func BadRequest(message string) Error {
	return Error{
		ErrorCode: "BadRequest",
		Message:   message,
	}
}
