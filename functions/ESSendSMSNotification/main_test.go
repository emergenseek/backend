package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/emergenseek/backend/common"
)

func TestHandler(t *testing.T) {
	message := "Hello from Lambda"
	successfulResponse := "{\"body\":\"Successfully sent SMS to contacts of user John Doe (e78e0c86-f9ba-4375-9554-6dc1426f5605)\"}"
	tests := []struct {
		name              string
		UserID            string
		Type              common.EmergencyType
		Message           string
		LastKnownLocation []float64
		ExpectedBody      string
	}{
		{
			"SEVERE Emergency Request",
			"e78e0c86-f9ba-4375-9554-6dc1426f5605",
			1,
			"",
			[]float64{40.7648, -73.9808},
			successfulResponse,
		},
		{
			"MILD Emergency Request",
			"e78e0c86-f9ba-4375-9554-6dc1426f5605",
			2,
			"",
			[]float64{40.7648, -73.9808},
			successfulResponse,
		},
		{
			"CHECKIN Emergency Request",
			"e78e0c86-f9ba-4375-9554-6dc1426f5605",
			3,
			message,
			[]float64{40.7648, -73.9808},
			successfulResponse,
		},
		{
			"Invalid UserID",
			"e78e0c86-f9ba-4375-9554-6dc1426f5600",
			1,
			"",
			[]float64{40.7648, -73.9808},
			"{\"code\":\"Bad Request\",\"error\":\"user not found\"}",
		},
		{
			"Missing UserID",
			"",
			1,
			"",
			[]float64{40.7648, -73.9808},
		},
			"Invalid EmergencyType",
			"e78e0c86-f9ba-4375-9554-6dc1426f5605",
			50,
			"",
			[]float64{40.7648, -73.9808},
			"{\"code\":\"Bad Request\",\"error\":\"50 is an invalid emergency type\"}",
		},
		{
			"Missing Message w/ CHECKIN EmergencyType",
			"e78e0c86-f9ba-4375-9554-6dc1426f5605",
			3,
			"",
			[]float64{40.7648, -73.9808},
			"{\"code\":\"Bad Request\",\"error\":\"message field is required because emergency type is 3\"}",
		},
		{
			"Successful Request (MILD)",
			"e78e0c86-f9ba-4375-9554-6dc1426f5605",
			2,
			"",
			[]float64{40.7648, -73.9808},
			successfulResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test: %v", tt.name)
			r := Request{
				UserID:   tt.UserID,
				Type:     tt.Type,
				Message:  tt.Message,
				Location: tt.LastKnownLocation,
			}

			b, _ := json.Marshal(r)
			request := events.APIGatewayProxyRequest{
				Path:       "/sms",
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
