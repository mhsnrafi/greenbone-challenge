package main

import (
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"greenbone-task/models"
	db "greenbone-task/models/db"
	"greenbone-task/services"
	"testing"
)

func TestCreateComputer(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	employee := models.EmployeeRequest{
		FirstName:    "Test",
		LastName:     "Dummy",
		Abbreviation: "DTT",
		Email:        "test.dummy3@tes2t.com",
	}
	err := services.CreateEmployee(employee)
	require.NoError(t, err)

	computer := db.Computer{
		MacAddress:     "aa:ac:cc:ed:ee:ff",
		ComputerName:   "Dummy'ss Desktop",
		EmployeeAbbrev: employee.Abbreviation,
		IPAddress:      "192.168.4.105",
		Description:    "Custom-built PC1",
	}

	id, err := services.CreateComputer(computer)
	require.NoError(t, err)
	require.NotEqual(t, uint(0), id)

}

func TestGetAllComputers(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	// Test
	computers, err := services.GetAllComputers()
	require.NoError(t, err)
	require.NotNil(t, computers)
	assert.True(t, len(computers) > 0)
}

func TestAssignComputerToEmployee(t *testing.T) {
	services.LoadConfig()
	services.ConnectDB()

	// Create test data
	testEmployeeAbbrev := "JAD"
	testComputer := &db.Computer{
		MacAddress:     "az:bx:cd:ed:ee:ff",
		ComputerName:   "Dummy'ss Desktop",
		EmployeeAbbrev: "JAD",
		IPAddress:      "192.168.7.121",
		Description:    "Custom-built PC1",
	}
	err := services.DbConnection.Create(testComputer).Error
	require.NoError(t, err)

	// Test case 1: Assign computer to employee for the first time
	err = services.AssignComputerToEmployee(cast.ToInt64(testComputer.ID), testEmployeeAbbrev)
	require.NoError(t, err)

	employee, err := services.FindByEmployeeAbbrev(testEmployeeAbbrev)
	require.NoError(t, err)

	var employeeComputer db.EmployeeComputer
	err = services.DbConnection.Where("employee_id = ? AND computer_id = ?", employee.ID, testComputer.ID).First(&employeeComputer).Error
	require.NoError(t, err)
	assert.Equal(t, employee.ID, employeeComputer.EmployeeID)

	// Test case 2: Assign computer to a different employee
	otherTestEmployeeAbbrev := "JDE"
	otherEmployee, err := services.FindByEmployeeAbbrev(otherTestEmployeeAbbrev)
	require.NoError(t, err)

	// Add existing relationship to the employee_computers table
	existingEmployeeComputer := &db.EmployeeComputer{
		EmployeeID: otherEmployee.ID,
		ComputerID: testComputer.ID,
	}
	err = services.DbConnection.Create(existingEmployeeComputer).Error
	require.NoError(t, err)

	err = services.AssignComputerToEmployee(cast.ToInt64(testComputer.ID), testEmployeeAbbrev)
	require.NoError(t, err)

	var newEmployeeComputer db.EmployeeComputer
	err = services.DbConnection.Where("employee_id = ? AND computer_id = ?", employee.ID, testComputer.ID).First(&newEmployeeComputer).Error
	require.NoError(t, err)
	assert.Equal(t, employee.ID, newEmployeeComputer.EmployeeID)

	// Test case 3: Assign computer to the same employee
	err = services.AssignComputerToEmployee(cast.ToInt64(testComputer.ID), testEmployeeAbbrev)
	require.NoError(t, err)

	err = services.DbConnection.Where("employee_id = ? AND computer_id = ?", employee.ID, testComputer.ID).First(&newEmployeeComputer).Error
	require.NoError(t, err)
	assert.Equal(t, employee.ID, newEmployeeComputer.EmployeeID)

	// Cleanup
	err = services.DeleteComputer(cast.ToInt64(testComputer.ID))
	require.NoError(t, err)

}
