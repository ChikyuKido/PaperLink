package routes

type Response struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

func NewSuccess(data any) Response {
	return Response{
		Code: 200,
		Data: data,
	}
}

func NewError(code int, msg string) Response {
	return Response{
		Code: code,
		Data: msg,
	}
}
