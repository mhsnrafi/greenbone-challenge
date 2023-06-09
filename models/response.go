package models

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response Base response
type Response struct {
	StatusCode int            `json:"-"`
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

func (response *Response) SendResponse(c *gin.Context) {
	c.AbortWithStatusJSON(response.StatusCode, response)
}

func SendResponseData(c *gin.Context, data gin.H) {
	response := &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       data,
	}
	response.SendResponse(c)
}

func SendErrorResponse(c *gin.Context, status int, message string) {
	response := &Response{
		StatusCode: status,
		Success:    false,
		Message:    message,
	}
	response.SendResponse(c)
}

type ComputerRequestResponse struct {
	EmployeeID           uint   `json:"employee_id"`
	EmployeeAbbreviation string `json:"employee_abbreviation"`
	ComputerName         string `json:"computer_name"`
}

type AssignComputerRequest struct {
	EmployeeAbbreviation string `json:"employee_abbreviation"`
	ComputerID           int64  `json:"computer_id"`
}
