package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"greenbone-task/logger"
	"greenbone-task/models"
	db "greenbone-task/models/db"
	"time"
)

var CacheExpiration = 30 * time.Minute

// CreateEmployee function creates a new emplpyee records
func CreateEmployee(employee models.EmployeeRequest) error {
	// Store employee details in database

	emp := db.Employee{
		FirstName:    employee.FirstName,
		LastName:     employee.LastName,
		Abbreviation: employee.Abbreviation,
		Email:        employee.Email,
		Computers:    []db.Computer{},
	}
	if err := DbConnection.Create(&emp).Error; err != nil {
		logger.Error("failed to save computer", zap.Error(err))
		return err
	}
	return nil
}

// DeleteEmployeeComputer delete specific employee computer
func DeleteEmployeeComputer(computerID int64, abbrev string) error {
	err := DbConnection.Delete(&db.Computer{}, computerID, abbrev).Error
	if err != nil {
		return fmt.Errorf("error deleting computer: %w", err)
	}
	return nil
}

// FindComputersByEmployeeAbbrev fidn the computers from the database using abbrev
func FindComputersByEmployeeAbbrev(abbrev string) ([]db.Computer, error) {
	// check cache first
	cacheKey := fmt.Sprintf("computers_by_employee_%s", abbrev)
	if cachedResult, err := GetRedisDefaultClient().Get(context.Background(), cacheKey).Result(); err == nil {
		var cachedComputers []db.Computer
		if err := json.Unmarshal([]byte(cachedResult), &cachedComputers); err == nil {
			return cachedComputers, nil
		}
		// cache hit but unmarshal error, fallback to DB query
	}

	var computers []db.Computer

	// find the employee
	employee, err := FindByEmployeeAbbrev(abbrev)
	if err != nil {
		return nil, fmt.Errorf("error finding employee: %w", err)
	}

	// fetch the associated computers using explicit join condition
	if err := DbConnection.Joins("JOIN employee_computers ON computers.id = employee_computers.computer_id").
		Where("employee_computers.employee_id = ?", employee.ID).
		Find(&computers).Error; err != nil {
		return nil, fmt.Errorf("error finding computers: %w", err)
	}

	// cache the result for future use
	cachedResult, err := json.Marshal(computers)
	if err != nil {
		return computers, nil
	}
	GetRedisDefaultClient().Set(context.Background(), cacheKey, string(cachedResult), CacheExpiration)

	return computers, nil
}
