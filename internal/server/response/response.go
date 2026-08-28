package response

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
