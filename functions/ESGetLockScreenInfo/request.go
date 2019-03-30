package main

import (
	"fmt"
)

// Request defines the expected body parameters of an ESGetLockScreenInfo invocation
type Request struct {

	// The ID of the user making the request
	UserID string `json:"user_id"`
}

// Validate validates a request to the ESGetLockScreenInfo Lambda function
func (r *Request) Validate() error {

	// Check if the UserID is present
	if r.UserID == "" {
		return fmt.Errorf("user_id field is required")
	}

	return nil
}
