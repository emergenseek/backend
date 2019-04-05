package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/events"
	"github.com/emergenseek/backend/common"
)

func TestHandler(t *testing.T) {
	message := "Hello from Lambda"
	successfulResponse := "{\"body\":\"Successfully sent SMS to contacts of user EmergenSeek User (4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9)\"}"
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
			"4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9",
			1,
			"",
			[]float64{40.7648, -73.9808},
			successfulResponse,
		},
		{
			"MILD Emergency Request",
			"4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9",
			2,
			"",
			[]float64{40.7648, -73.9808},
			successfulResponse,
		},
		{
			"CHECKIN Emergency Request",
			"4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9",
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
			"{\"code\":\"Bad Request\",\"error\":\"user_id field is required\"}",
		},
		{
			"Invalid EmergencyType",
			"4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9",
			50,
			"",
			[]float64{40.7648, -73.9808},
			"{\"code\":\"Bad Request\",\"error\":\"50 is an invalid emergency type\"}",
		},
		{
			"Missing Message w/ CHECKIN EmergencyType",
			"4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9",
			3,
			"",
			[]float64{40.7648, -73.9808},
			"{\"code\":\"Bad Request\",\"error\":\"message field is required because emergency type is 3\"}",
		},
		{
			"Successful Request (MILD)",
			"4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9",
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
