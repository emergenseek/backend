package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	retry "github.com/avast/retry-go"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/emergenseek/backend/common/driver"
)

func verifyRequest(request events.APIGatewayProxyRequest) (*Request, int, error) {
	// Create a new request object and unmarshal the request body into it
	req := new(Request)
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	// Make sure all of the necessary parameters are present
	err = req.Validate()
	if err != nil {
		return nil, http.StatusBadRequest, err

	}
	// All checks passed, return req struct for use. http.StatusOK is ignored
	return req, http.StatusOK, nil
}

// Handler is the Lambda handler for ESPollLocation
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request
	req, status, err := verifyRequest(request)
	if err != nil {
		return driver.ErrorResponse(status, err), nil
	}

	// Initialize drivers
	db, twilio, _, mapsKey := driver.CreateAll()

	// Retrieve user from database
	user, err := db.GetUser(req.UserID)
	if err != nil {
		return driver.ErrorResponse(http.StatusBadRequest, err), nil
	}

	// Update the user's last known location
	err = db.UpdateLocation(user.UserID, req.Location)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	// Generate the message to be sent to contacts
	message, err := driver.CreatePollMessage(user, mapsKey, req.Location)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	for _, contact := range user.PrimaryContacts {
		retry.Do(
			func() error { return twilio.SendSMS(contact.PhoneNumber, message) },
			retry.Attempts(3),
		)
	}
	for _, contact := range user.SecondaryContacts {
		retry.Do(
			func() error { return twilio.SendSMS(contact.PhoneNumber, message) },
			retry.Attempts(3),
		)
	}

	bodyContent := fmt.Sprintf("Successfully sent poll message to contents of user %v %v (%v)", user.FirstName, user.LastName, user.UserID)
	return driver.SuccessfulResponse(bodyContent, user), nil
}

func main() {
	lambda.Start(Handler)
}
