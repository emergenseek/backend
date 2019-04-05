package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	successfulResponse := "{\"body\":\"Successfully sent emergency call to emergency services and contacts of user EmergenSeek User (4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9)\"}"
	tests := []struct {
		name              string
		UserID            string
		LastKnownLocation []float64
		ExpectedBody      string
	}{
		{
			"Simple Request",
			"4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9",
			[]float64{40.7648, -73.9808},
			successfulResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test: %v", tt.name)
			r := Request{
				UserID:   tt.UserID,
				Location: tt.LastKnownLocation,
			}

			b, _ := json.Marshal(r)
			request := events.APIGatewayProxyRequest{
				Path:       "/voice",
				HTTPMethod: "POST",
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       string(b),
			}

			expectedResponse := events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: tt.ExpectedBody,
			}

			response, err := Handler(request)
			assert.Equal(t, response.Headers, expectedResponse.Headers)
			assert.Contains(t, response.Body, expectedResponse.Body)
			assert.Equal(t, err, nil)
		})
	}

}
