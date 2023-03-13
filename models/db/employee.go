package models

import (
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	FirstName    string     `json:"first_name" gorm:"not null"`
	LastName     string     `json:"last_name" gorm:"not null"`
	Email        string     `json:"email" gorm:"unique;not null"`
	Abbreviation string     `json:"abbreviation" gorm:"unique;not null"`
	Computers    []Computer `json:"computers" gorm:"many2many:employee_computers;"`
}

func (Employee) TableName() string {
	return "employees"
}
