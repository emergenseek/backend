package main

import (
	"fmt"
)

// Request defines the expected body parameters of an ESServiceLocator invocation
type Request struct {
	// The user's location at the time of the request
	Location []float64 `json:"current_location"`
}

// Validate validates a request to the ESServiceLocator Lambda function
func (r *Request) Validate() error {
	if r.Location == nil {
		return fmt.Errorf("current_location field is required")
	}

	return nil
}
