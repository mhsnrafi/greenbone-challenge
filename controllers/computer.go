package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"
	"greenbone-task/models"
	db "greenbone-task/models/db"
	"greenbone-task/services"
	"net/http"
)

// CreateComputer handles the request to create a new computer
// @Summary Create a new computer
// @Description Create a new computer with the given details
// @Tags Computers
// @Accept json
// @Produce json
// @Param computerReq body db.Computer true "Computer details"
// @Success 201 {object} models.Response
// @Failure 400 Bad Request models.Response
// @Router /computers [post]
func CreateComputer(c *gin.Context) {
	var computerReq db.Computer
	if err := c.ShouldBindBodyWith(&computerReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.ValidateComputerRequest(computerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// process the computer creation request
	computerID, err := services.CreateComputer(computerReq)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// Return success response
	response.Success = true
	response.StatusCode = http.StatusCreated
	response.Data = gin.H{
		"Computer ID": computerID,
		"Message":     "Computer created successfully",
	}
	response.SendResponse(c)
}

// GetComputerByID handles the request to get a computer by its ID
// @Summary Get a computer by ID
// @Description Get a computer with the given ID
// @Tags Computers
// @Accept json
// @Produce json
// @Param computer_id query int true "Computer ID"
// @Success 200 {object} models.Response
// @Failure 400 Bad Request models.Response
// @Router /computers{id} [get]
func GetComputerByID(c *gin.Context) {
	computerID := c.Param("computer_id")
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// Validate computer ID
	if len(computerID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "computer ID is required"})
		return
	}

	// process the computer creation request
	data, err := services.GetComputerByID(cast.ToInt64(computerID))
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// Return success response
	response.Success = true
	response.StatusCode = http.StatusOK
	response.Data = gin.H{
		"Message": "Computer Information Fetch Successfully",
		"Data":    data,
	}
	response.SendResponse(c)
}

// GetAllComputers handles the request to fetch all computers
// @Summary Fetch all computers
// @Description Fetch all computers and their details
// @Tags Computers
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 Bad Request models.Response
// @Router /computers [get]
func GetAllComputers(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// process the computer creation request
	data, err := services.GetAllComputers()
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// Return success response
	response.Success = true
	response.StatusCode = http.StatusOK
	response.Data = gin.H{
		"Message": "Fetch all computers successfully",
		"Data":    data,
	}
	response.SendResponse(c)
}

// UpdateComputer handles the request to update an existing computer  with the given ID and employee abbreviation. It updates the computer's employee association to the given employee.
// @Summary Update an existing computer
// @Description Update an existing computer with the given ID and employee abbreviation
// @Tags Computers
// @Accept json
// @Produce json
// @Param computer_id path int true "Computer ID to update"
// @Param employee_abbrev path string true "Employee abbreviation to assign computer to"
// @Success 201 {object} models.Response
// @Failure 400 Bad Request models.Response
// @Router /computers/{computer_id}/assign/{employee_abbrev} [put]
func UpdateComputer(c *gin.Context) {
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
	err := services.AssignComputerToEmployee(cast.ToInt64(computerID), employeeAbbrev)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// Return success response
	response.Success = true
	response.StatusCode = http.StatusCreated
	response.Data = gin.H{
		"Message": "Computer updated successfully",
	}
	response.SendResponse(c)
}

// DeleteComputer handles the request to delete a computer
// @Summary Delete a computer
// @Description Delete a computer with the given ID
// @Tags Computers
// @Accept json
// @Produce json
// @Param computer_id query int true "Computer ID"
// @Success 200 {object} models.Response
// @Failure 400 Bad Request
func DeleteComputer(c *gin.Context) {
	computerID := c.Request.URL.Query().Get("computer_id")
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// Validate computer ID
	if len(computerID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "computer ID is required"})
		return
	}

	// process the computer creation request
	err := services.DeleteComputer(cast.ToInt64(computerID))
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
