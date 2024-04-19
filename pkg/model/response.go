package model

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func Pong() *Response {
	return &Response{
		Status:  true,
		Message: "pong",
	}
}

func Success() *Response {
	return &Response{
		Status:  true,
		Message: "success",
	}
}

func Error(err error) *Response {
	return &Response{
		Status:  false,
		Message: err.Error(),
	}
}
