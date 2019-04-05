package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	retry "github.com/avast/retry-go"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/emergenseek/backend/common"
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

// Handler is the Lambda handler for ESSendEmergencyVoiceCall
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request
	req, status, err := verifyRequest(request)
	if err != nil {
		return driver.ErrorResponse(status, err), nil
	}

	// Initialize drivers
	db, twilio, sess, mapKey := driver.CreateAll()

	// Retrieve user from database
	user, err := db.GetUser(req.UserID)
	if err != nil {
		return driver.ErrorResponse(http.StatusBadRequest, err), nil
	}

	// Update the user's last known location
	err = db.UpdateLocation(user.UserID, req.Location)

	// Convert the user's last known location to a human readable address
	address, err := driver.GetAddress(req.Location, mapKey, false, 0)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	// Create the Twilio Markdown Language necessary for the voice call
	twilML, err := driver.CreateTwilMLXML(user, address)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	// Upload the TwilML to S3 so Twilio can access it
	callbackURL, err := driver.UploadTwilMLXML(twilML, sess)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	// Call all tier 1 contacts
	for _, contact := range user.Contacts {
		if contact.Tier == common.FIRST {
			retry.Do(
				func() error { return twilio.SendVoiceCall(contact.PhoneNumber, callbackURL) },
				retry.Attempts(3),
			)
		}
	}

	bodyContent := fmt.Sprintf("Successfully sent emergency call to emergency services and contacts of user %v %v (%v)", user.FirstName, user.LastName, user.UserID)
	return driver.SuccessfulResponse(bodyContent), nil
}

func main() {
	lambda.Start(Handler)
}
