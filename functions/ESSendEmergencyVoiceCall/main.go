package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	retry "github.com/avast/retry-go"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/emergenseek/backend/common/driver"
)

func verifyRequest(request events.APIGatewayProxyRequest) (*Request, int, error) {
	// Only allow JSON requests
	if request.Headers["Content-Type"] != "application/json" {
		return nil, http.StatusNotAcceptable, errors.New("content-type must be application/json")
	}

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

// Handler is the Lambda handler for ESSendEmergencyVoiceCall
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request
	req, status, err := verifyRequest(request)
	if err != nil {
		return driver.ErrorResponse(status, err), nil
	}

	// Initialize drivers
	db, twilio := driver.CreateAll()

	// Retrieve user from database
	user, err := db.GetUser(req.UserID)
	if err != nil {
		return driver.ErrorResponse(http.StatusBadRequest, err), nil
	}

	/* TODO
	3. Extract last known location
	4. Use last known location to get emergency service number (911 or 119, etc)
	5. Generate message for voice call
	6. Call all primary contacts and emergency service number
	*/

	// Update phoneNumbers to include all numbers necessary
	for _, contact := range user.PrimaryContacts {
		retry.Do(
			func() error { return twilio.SendVoiceCall(contact.PhoneNumber) },
			retry.Attempts(3),
		)
	}
	for _, contact := range user.SecondaryContacts {
		retry.Do(
			func() error { return twilio.SendVoiceCall(contact.PhoneNumber) },
			retry.Attempts(3),
		)
	}

	bodyContent := fmt.Sprintf("Successfully sent emergency call to emergency services and contacts of user %v %v (%v)", user.FirstName, user.LastName, user.CognitoID)
	return driver.SuccessfulResponse(bodyContent, user), nil

}

func main() {
	lambda.Start(Handler)
}
