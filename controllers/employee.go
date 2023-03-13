package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"
	"greenbone-task/models"
	"greenbone-task/services"
	"net/http"
)

// CreateEmployee handles the request to create a new employee
// @Summary Create a new employee
// @Description Create a new employee with the given details
// @Tags Employees
// @Accept json
// @Produce json
// @Param emp body models.EmployeeRequest true "Employee details"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /employees [post]
func CreateEmployee(c *gin.Context) {
	var emp models.EmployeeRequest
	if err := c.ShouldBindBodyWith(&emp, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.ValidateEmployeeRequest(emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// process the computer creation request
	err := services.CreateEmployee(emp)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// Return success response
	response.Success = true
	response.StatusCode = http.StatusCreated
	response.Data = gin.H{
		"Message": "Employee record created successfully",
	}
	response.SendResponse(c)
}

// GetEmployeeComputers handles the request to retrieve all computers assigned to a specific employee
// @Summary Retrieve all computers assigned to an employee
// @Description Retrieve all computers assigned to an employee with the given employee abbreviation
// @Tags Computers
// @Accept json
// @Produce json
// @Param employee_abbrev query string true "Employee abbreviation"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /employees/computers/{employee_abbrev} [get]
func GetEmployeeComputers(c *gin.Context) {
	employeeAbbrev := c.Param("employee_abbrev")

	// Validate employee abbreviation
	if len(employeeAbbrev) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee abbreviation is required"})
		return
	}

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// process the computer creation request
	empComputers, err := services.FindComputersByEmployeeAbbrev(employeeAbbrev)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// Return success response
	response.Success = true
	response.StatusCode = http.StatusOK
	response.Data = gin.H{
		"Message": "List of all the computers assigned to the employee",
		"Data":    empComputers,
	}
	response.SendResponse(c)
}

// DeleteEmployeeComputer handles the request to delete a computer assigned to an employee
// @Summary Delete a computer assigned to an employee
// @Description Delete a computer assigned to an employee with the given computer ID and employee abbreviation
// @Tags Computers
// @Accept json
// @Produce json
// @Param computer_id query int true "Computer ID to delete"
// @Param employee_abbrev query string true "Employee abbreviation"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /employees/{employee_abbrev}/computers/{computer_id} [delete]
func DeleteEmployeeComputer(c *gin.Context) {
	computerID := c.Param("computer_id")
	employeeAbbrev := c.Param("employee_abbrev")
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// Validate computer ID
	if len(computerID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "computer ID is required"})
		return
	}

	// Validate employee abbreviation
	if len(employeeAbbrev) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee abbreviation is required"})
		return
	}

	// process the computer creation request
	err := services.DeleteEmployeeComputer(cast.ToInt64(computerID), employeeAbbrev)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// Return success response
	response.Success = true
	response.StatusCode = http.StatusOK
	response.Data = gin.H{
		"Message": "Computer deleted successfully",
	}
	response.SendResponse(c)
}
