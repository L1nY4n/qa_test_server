package web

import "github.com/gin-gonic/gin"

type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(200, apiResponse{
		Success: true,
		Data:    data,
	})
}

func fail(c *gin.Context, status int, message string) {
	c.JSON(status, apiResponse{
		Success: false,
		Error:   message,
	})
}
