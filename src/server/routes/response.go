package routes

import "github.com/gin-gonic/gin"

type Response struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func JSONSuccess(c *gin.Context, httpCode int, data any) {
	c.JSON(httpCode, Response{
		Code: httpCode,
		Data: data,
	})
}

func JSONSuccessOK(c *gin.Context, data any) {
	JSONSuccess(c, 200, data)
}

func JSONError(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, ErrorResponse{
		Code:  httpCode,
		Error: msg,
	})
}
