package models

import (
	"gorm.io/gorm"
)

type Computer struct {
	gorm.Model
	MacAddress     string `json:"mac_address" gorm:"unique;not null"`
	ComputerName   string `json:"computer_name" gorm:"not null"`
	IPAddress      string `json:"ip_address" gorm:"not null"`
	EmployeeAbbrev string `json:"employee_abbrev,omitempty"`
	Description    string `json:"description,omitempty"`
}

func (Computer) TableName() string {
	return "computers"
}

type EmployeeComputer struct {
	EmployeeID uint `json:"employee_id"`
	ComputerID uint `json:"computer_id"`
}
