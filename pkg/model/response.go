package model

type Response struct {
	Message string `json:"message"`
}

func Pong() *Response {
	return &Response{
		Message: "pong",
	}
}

func Success() *Response {
	return &Response{
		Message: "success",
	}
}

func Error(err error) *Response {
	return &Response{
		Message: err.Error(),
	}
}