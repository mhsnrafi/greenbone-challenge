package services

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"greenbone-task/logger"
	db "greenbone-task/models/db"
	"reflect"
	"time"
)

// CreateComputer function creates a new computer and assigns it to an employee.
func CreateComputer(computer db.Computer) (uint, error) {
	// Store computer details in database
	if err := DbConnection.Create(&computer).Error; err != nil {
		logger.Error("failed to save computer", zap.Error(err))
		return 0, err
	}

	// check if the employee exists
	// Assign computer to employee
	employee, err := FindByEmployeeAbbrev(computer.EmployeeAbbrev)
	if err != nil {
		logger.Error("failed to assign computer to employee", zap.Error(err))
		return 0, fmt.Errorf("error assigning computer to employee: %w", err)
	}

	if reflect.DeepEqual(employee, db.Employee{}) {
		logger.Error("failed to assign computer to employee", zap.Error(err))
		return 0, fmt.Errorf("error assigning computer to employee: employee not found")
	}

	// check if the employee already has 3 computers assigned
	count, err := CountComputersByEmployeeAbbreviation(computer.EmployeeAbbrev)
	if err != nil {
		return 0, fmt.Errorf("error assigning computer to employee: %w", err)
	}
	if count >= 3 {
		// notify system administrator about the assignment
		ns := NewNotificationService()
		message := fmt.Sprintf("Employee %s already has %d computers assigned.", computer.EmployeeAbbrev, count)
		err := ns.NotifySystemAdministrator(employee.Abbreviation, message)
		if err != nil {
			return 0, fmt.Errorf("error assigning computer to employee: %w", err)
		}
	}

	// assign the computer to the employee
	computer.EmployeeAbbrev = computer.EmployeeAbbrev
	err = DbConnection.Save(computer).Error
	if err != nil {
		return 0, fmt.Errorf("error assigning computer to employee: %w", err)
	}

	// create an entry in the employee_computer junction table
	err = DbConnection.Model(&employee).Association("Computers").Append(&computer).Error
	if err != nil {
		return 0, fmt.Errorf("error assigning computer to employee: %w", err)
	}

	return computer.ID, nil
}

func GetAllComputers() ([]db.Computer, error) {
	var computers []db.Computer
	err := DbConnection.Find(&computers).Error
	if err != nil {
		return nil, fmt.Errorf("error getting all computers: %w", err)
	}
	return computers, nil
}

// GetComputerByID function get computer information from id
func GetComputerByID(id int64) (*db.Computer, error) {
	// First, try to get the computer from the cache.
	computerKey := fmt.Sprintf("computer:%d", id)
	ctx := context.Background()
	cachedComputer, err := GetRedisDefaultClient().Get(ctx, computerKey).Result()
	if err == nil {
		var computer db.Computer
		if err := json.Unmarshal([]byte(cachedComputer), &computer); err != nil {
			return nil, fmt.Errorf("error unmarshaling cached computer data: %w", err)
		}
		return &computer, nil
	} else if err != redis.Nil {
		return nil, fmt.Errorf("error getting computer from Redis cache: %w", err)
	}

	// If the computer is not in the cache, query it from the database.
	var computer db.Computer
	err = DbConnection.Where("id = ?", id).First(&computer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no computer found with ID: %d", id)
		}
		return nil, fmt.Errorf("error getting computer by ID: %w", err)
	}

	// Store the computer in the cache.
	computerJSON, err := json.Marshal(computer)
	if err != nil {
		return nil, fmt.Errorf("error marshaling computer data: %w", err)
	}
	if err := GetRedisDefaultClient().Set(ctx, computerKey, computerJSON, time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("error setting computer in Redis cache: %w", err)
	}

	return &computer, nil
}

// DeleteComputer delete computer from the database from computer id
func DeleteComputer(id int64) error {
	err := DbConnection.Delete(&db.Computer{}, id).Error
	if err != nil {
		return fmt.Errorf("error deleting computer: %w", err)
	}
	return nil
}

// AssignComputerToEmployee assign employee computer to another employee
func AssignComputerToEmployee(computerID int64, newEmployeeAbbreviation string) error {
	// Get the computer record by ID
	computer, err := GetComputerByID(computerID)
	if err != nil {
		return fmt.Errorf("error getting computer by ID: %w", err)
	}
	if computer == nil {
		return fmt.Errorf("computer not found with ID %d", computerID)
	}

	// Get the new employee record by abbreviation
	newEmployee, err := FindByEmployeeAbbrev(newEmployeeAbbreviation)
	if err != nil {
		return fmt.Errorf("error getting employee by abbreviation: %w", err)
	}
	if reflect.DeepEqual(newEmployee, db.Employee{}) {
		return fmt.Errorf("employee not found with abbreviation %s", newEmployeeAbbreviation)
	}

	// Check if the record already exists in the employee_computers table
	var existingEmployeeComputer db.EmployeeComputer
	err = DbConnection.Where("computer_id = ?", computer.ID).First(&existingEmployeeComputer).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking employee_computers table: %w", err)
	}

	// If the existing record has the same employee ID as the new employee, there's nothing to do
	if existingEmployeeComputer.EmployeeID == newEmployee.ID {
		return nil
	}

	// If the existing record has a different employee ID, update it with the new employee ID
	if existingEmployeeComputer.EmployeeID != 0 {
		existingEmployeeComputer.EmployeeID = newEmployee.ID
		err = DbConnection.Save(&existingEmployeeComputer).Error
		if err != nil {
			return fmt.Errorf("error updating employee_computers record: %w", err)
		}
	} else {
		// If the record does not exist, insert a new record in the employee_computers junction table
		newEmployeeComputer := &db.EmployeeComputer{
			EmployeeID: newEmployee.ID,
			ComputerID: computer.ID,
		}
		err = DbConnection.Create(newEmployeeComputer).Error
		if err != nil {
			return fmt.Errorf("error inserting employee_computers record: %w", err)
		}
	}

	return nil
}

// FindByEmployeeAbbrev fetch data from employee table using abbreviation
func FindByEmployeeAbbrev(abbrev string) (db.Employee, error) {
	var employee db.Employee
	if err := DbConnection.Where("abbreviation = ?", abbrev).First(&employee).Error; err != nil {
		fmt.Println("Error finding card:", err)
		return db.Employee{}, errors.New("failed to find computers by employee abbreviation")
	}

	return employee, nil
}

// CountComputersByEmployeeAbbreviation count no of computer assign to employee
func CountComputersByEmployeeAbbreviation(abbreviation string) (int64, error) {
	var count int64
	result := DbConnection.Model(&db.Computer{}).Where("employee_abbrev = ?", abbreviation).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
