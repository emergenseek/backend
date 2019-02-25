package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/emergenseek/backend/common"
)

// Request defines the expected body parameters of an ESSendEmergencyVoiceCall invocation
type Request struct {
	// see "github.com/emergenseek/backend/common.const for EmergencyType
	Type common.EmergencyType `json:"type"`

	// The ID of the user making the request
	UserID string `json:"user_id"`
}

// Handler is the Lambda handler for ESSendEmergencyVoiceCall
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Only allow JSON requests
	if request.Headers["Content-Type"] != "application/json" {
		return common.ClientError(http.StatusNotAcceptable, errors.New("invalid content type"))
	}

	// Create a new request object and unmarshal the request body into it
	req := new(Request)
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return common.ClientError(http.StatusUnprocessableEntity, err)
	}

	// Make sure all of the necessary parameters are present
	if req.Type == 0 || req.UserID == "" {
		return common.ClientError(http.StatusBadRequest, errors.New("invalid parameter"))
	}

	// Initialize drivers and send the voice call
	_, t, err := common.InitAll()
	if err != nil {
		return common.ClientError(http.StatusInternalServerError, err)
	}
	/* TODO
	1.  Look up user in database
	2. Extract phone number of all primary contacts
	3. Extract last known location
	4. Use last known location to get emergency service number (911 or 119, etc)
	5. Generate message for voice call
	6. Call all primary contacts and emergency service number
	*/

	// Update phoneNumbers to include all numbers necessary
	phoneNumbers := []string{t.TargetNumber}
	err = t.SendVoiceCall(phoneNumbers)
	if err != nil {
		return common.ClientError(http.StatusInternalServerError, err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("Successfully sent voice calls on behalf of user: %v", req.UserID),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
