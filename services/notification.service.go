package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"greenbone-task/constants"
	"net/http"
)

type NotificationService interface {
	NotifySystemAdministrator(employeeAbbreviation string, message string) error
}

type notificationService struct{}

func NewNotificationService() NotificationService {
	return &notificationService{}
}

func (ns *notificationService) NotifySystemAdministrator(employeeAbbreviation string, message string) error {
	notification := map[string]string{
		"level":                constants.Warning,
		"employeeAbbreviation": employeeAbbreviation,
		"message":              message,
	}

	jsonBytes, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error marshalling notification JSON: %w", err)
	}

	// Send HTTP POST request to the notification service
	// create the request
	req, err := http.NewRequest("POST", constants.NOTIFICATION_URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		// handle error
	}

	// set request headers
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending notification request: unexpected status code %d", resp.StatusCode)

	}
	defer resp.Body.Close()

	// handle the response
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return fmt.Errorf("error in the json request")
	}

	return nil
}
