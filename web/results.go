package web

type Error struct {
	ErrorCode string `json:"errorCode" example:"NotFound"`
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
