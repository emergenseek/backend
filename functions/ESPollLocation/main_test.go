package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	successfulResponse := "{\"body\":\"Successfully sent poll message to contents of user John Doe (e78e0c86-f9ba-4375-9554-6dc1426f5605)\"}"
	tests := []struct {
		name              string
		UserID            string
		LastKnownLocation []float64
		ExpectedBody      string
	}{
		{
			"Simple Request",
			"e78e0c86-f9ba-4375-9554-6dc1426f5605",
			[]float64{-31.9517231, 115.8603252},
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
				Path:       "/poll",
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
