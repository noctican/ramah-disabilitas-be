package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func FormatResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	return Response{
		Meta: meta,
		Data: data,
	}
}

func JSONResponse(c *gin.Context, code int, message string, status string, data interface{}) {
	response := FormatResponse(message, code, status, data)
	c.JSON(code, response)
}
