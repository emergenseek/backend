package main

import (
	"fmt"
)

// Request defines the expected body parameters of an ESSendEmergencyVoiceCall invocation
type Request struct {
	// The ID of the user making the request
	UserID string `json:"user_id"`

	// When requests are made to this Lambda function,
	// it is assumed that the emergency type is 1 (SEVERE)
	// Type common.EmergencyType `json:"type"`

}

// Validate validates a request to the ESSendSMSNotification Lambda function
func (r *Request) Validate() error {
	// Check if the UserID is present
	if r.UserID == "" {
		return fmt.Errorf("user_id field is required")
	}
	return nil
}
