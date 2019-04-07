package main

import (
	"fmt"

	"github.com/emergenseek/backend/common"
)

// Request defines the expected body parameters of an ESUpdateTier invocation
type Request struct {
	// The phone number of the user who's tier is being updated
	PhoneNumber string `json:"phone_number"`

	// The new tier to be applied to the user's contact
	NewTier common.AlertTier `json:"new_tier"`
}

// Validate validates a request to the ESUpdateTier Lambda function
func (r *Request) Validate() error {
	// Check if the PhoneNumber is present
	if r.PhoneNumber == "" {
		return fmt.Errorf("phone_number field is required")
	}

	// Try to convert the EmergencyType to a string
	if r.NewTier.String() == "Unknown" {
		return fmt.Errorf("%d is an invalid tier", r.NewTier)
	}

	return nil
}
