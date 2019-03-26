package main

import (
	"fmt"
)

// Request defines the expected body parameters of an ESPollLocation invocation
type Request struct {
	// The ID of the user making the request
	UserID string `json:"user_id"`

	// The new location of the user
	Location []float64 `json:"last_known_location"`
}

// Validate validates a request to the ESPollLocation Lambda function
func (r *Request) Validate() error {
	// Check if the UserID is present
	if r.UserID == "" {
		return fmt.Errorf("user_id field is required")
	}

	// Check if the Location is present
	if r.Location == nil {
		return fmt.Errorf("last_known_location field is required")
	}

	return nil
}
