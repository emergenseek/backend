package main

import (
	"fmt"

	"github.com/emergenseek/backend/common"
)

// Request defines the expected body parameters of an ESSendSMSNotification invocation
type Request struct {
	// see "github.com/emergenseek/backend/common.const for EmergencyType
	Type common.EmergencyType `json:"type"`

	// The ID of the user making the request
	UserID string `json:"user_id"`

	// The message to send to primary contacts
	Message string `json:"message"`

	// The user's location at the time of the request
	Location []float64 `json:"last_known_location"`
}

// Validate validates a request to the ESSendSMSNotification Lambda function
func (r *Request) Validate() []error {
	errs := []error{}

	// Try to convert the EmergencyType to a string
	if r.Type.String() == "Unknown" {
		errs = append(errs, fmt.Errorf("%d is an invalid emergency type", r.Type))
	}

	// Check if the UserID is present
	if r.UserID == "" {
		errs = append(errs, fmt.Errorf("user_id field is required"))
	}

	// The message is only required if the emergency type is 3
	if r.Type.String() == "CHECKIN" && r.Message == "" {
		errs = append(errs, fmt.Errorf("message field is required because emergency type is %d", r.Type))
	}

	// If the errs slice is empty, return nil
	if len(errs) == 0 {
		return nil
	}

	// Check if the Location is present
	if r.Location == nil {
		errs = append(errs, fmt.Errorf("last_known_location field is required"))
	}

	return errs
}
