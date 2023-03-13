package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	db "greenbone-task/models/db"
	"regexp"
)

type AuthRequest struct {
	Email string `json:"email"`
}

func (a AuthRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
	)
}

type RefreshRequest struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func (a RefreshRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Token,
			validation.Required,
			validation.Match(regexp.MustCompile("^\\S+$")).Error("cannot contain whitespaces"),
		),
	)
}

type ComputerRequest struct {
	MacAddress     string `json:"mac_address"`
	ComputerName   string `json:"computer_name"`
	IPAddress      string `json:"ip_address"`
	EmployeeAbbrev string `json:"employee_abbrev,omitempty"`
	Description    string `json:"description,omitempty"`
}

type EmployeeRequest struct {
	FirstName    string            `json:"first_name"`
	LastName     string            `json:"last_name"`
	Email        string            `json:"email"`
	Abbreviation string            `json:"abbreviation"`
	Computers    []ComputerRequest `json:"computers,omitempty"`
}

type CPaymentResponse struct {
	PaymentIdentifier string
	Status            string
}

func ValidateComputerRequest(computerReq db.Computer) error {
	// Check if required fields are present
	if computerReq.MacAddress == "" {
		return errors.New("mac_address is required")
	}
	if computerReq.ComputerName == "" {
		return errors.New("computer_name is required")
	}
	if computerReq.IPAddress == "" {
		return errors.New("ip_address is required")
	}
	return nil
}

func ValidateEmployeeRequest(req EmployeeRequest) error {
	if req.FirstName == "" {
		return fmt.Errorf("missing required field 'first_name'")
	}
	if req.LastName == "" {
		return fmt.Errorf("missing required field 'last_name'")
	}
	if req.Email == "" {
		return fmt.Errorf("missing required field 'email'")
	}
	if req.Abbreviation == "" {
		return fmt.Errorf("missing required field 'abbreviation'")
	}
	return nil
}
