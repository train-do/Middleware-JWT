package helper

import (
	"github.com/gin-gonic/gin"
)

type HTTPResponse struct {
	ErrorMsg string      `json:"error_msg,omitempty"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func ResponseOK(c *gin.Context, data interface{}, message string, httpStatusCode int) {
	c.JSON(httpStatusCode, HTTPResponse{
		Message: message,
		Data:    data,
	})
}

func ResponseError(c *gin.Context, errorMsg string, message string, httpStatusCode int) {
	c.JSON(httpStatusCode, HTTPResponse{
		ErrorMsg: errorMsg,
		Message:  message,
	})
}
