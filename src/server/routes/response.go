package routes

type Response struct {
	Code int `json:"code"`
	Data any `json:"data"`
}
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func NewSuccess(data any) Response {
	return Response{
		Code: 200,
		Data: data,
	}
}

func NewError(code int, msg string) ErrorResponse {
	return ErrorResponse{
		Code:  code,
		Error: msg,
	}
}
