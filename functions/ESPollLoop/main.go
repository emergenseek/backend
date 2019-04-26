package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	retry "github.com/avast/retry-go"
	"github.com/aws/aws-lambda-go/events"
	lambdaStart "github.com/aws/aws-lambda-go/lambda"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/emergenseek/backend/common"
	"github.com/emergenseek/backend/common/driver"
)

// Request defines the expected body parameters of an invocation
type Request struct {
	// The Tier of contacts this function works toward
	Tier common.AlertTier `json:"tier"`
}

func verifyRequest(request events.APIGatewayProxyRequest) (*Request, int, error) {
	// Create a new request object and unmarshal the request body into it
	req := new(Request)
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	// All checks passed, return req struct for use. http.StatusOK is ignored
	return req, http.StatusOK, nil
}

// Handler is the Lambda handler for ESPollLoop
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request
	req, status, err := verifyRequest(request)
	if err != nil {
		return driver.ErrorResponse(status, err), nil
	}

	// Initialize drivers
	db, twilio, _, mapsKey := driver.CreateAll()

	// Retrieve user from database
	userID := "a3cac301-f9a9-4df8-a224-d0bac718e4fa"
	user, err := db.GetUser(userID)
	if err != nil {
		return driver.ErrorResponse(http.StatusNotFound, err), nil
	}

	for _, contact := range user.Contacts {
		if contact.Tier == req.Tier {
			// Generate message based on the contact's tier
			message, err := driver.CreatePollMessage(user, mapsKey, user.LastKnownLocation, contact.Tier)
			if err != nil {
				return driver.ErrorResponse(http.StatusInternalServerError, err), nil
			}
			retry.Do(
				func() error { return twilio.SendSMS(contact.PhoneNumber, message) },
				retry.Attempts(3),
			)
		}
	}

	bodyContent := fmt.Sprintf("Successfully finised polling messages to users")
	return driver.SuccessfulResponse(bodyContent), nil
}

func main() {
	lambdaStart.Start(Handler)
}
