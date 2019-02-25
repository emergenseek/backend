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

// Request defines the expected body parameters of an ESSendSMSNotification invocation
type Request struct {
	// see "github.com/emergenseek/backend/common.const for EmergencyType
	Type common.EmergencyType `json:"type"`

	// The ID of the user making the request
	UserID string `json:"user_id"`

	// The message to send to primary contacts
	Message string `json:"message"`
}

// Handler is the Lambda handler for ESSendSMSNotification
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
	if req.Type == 0 || req.UserID == "" || req.Message == "" {
		return common.ClientError(http.StatusBadRequest, errors.New("invalid parameter"))
	}

	// Initialize drivers and send the SMS
	_, t, err := common.InitAll()
	if err != nil {
		return common.ClientError(http.StatusInternalServerError, err)
	}

	err = t.SendSMS(req.Message)
	if err != nil {
		return common.ClientError(http.StatusInternalServerError, err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("Successfully sent SMS to contacts of user: %v", req.UserID),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
