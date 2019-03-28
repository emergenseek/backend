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

func verifyRequest(request events.APIGatewayProxyRequest) (*Request, int, []error) {
	errs := []error{}

	// Create a new request object and unmarshal the request body into it
	req := new(Request)
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, append(errs, err)
	}

	// Make sure all of the necessary parameters are present
	errs = req.Validate()
	if errs != nil {
		return nil, http.StatusBadRequest, errs

	}
	// All checks passed, return req struct for use. http.StatusOK is ignored
	return req, http.StatusOK, nil
}

// Handler is the Lambda handler for ESSendSMSNotification
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request
	req, status, errs := verifyRequest(request)
	if errs != nil {
		return driver.ErrorResponse(status, errs...), nil
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

	// Send an SMS to all of their contacts depending on the servity of the emergency
	switch req.Type {
	case common.SEVERE:
		// If the request type is SEVERE
		// send an SMS message to all primary and secondary contacts with retry
		message, err := driver.CreateEmergencyMessage(common.SEVERE, user, mapsKey, req.Location)
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

	case common.MILD:
		// If the request type is MILD
		// send an SMS message only to primary contacts with retry
		message, err := driver.CreateEmergencyMessage(common.MILD, user, mapsKey, req.Location)
		if err != nil {
			return driver.ErrorResponse(http.StatusInternalServerError, err), nil
		}

		for _, contact := range user.PrimaryContacts {
			retry.Do(
				func() error { return twilio.SendSMS(contact.PhoneNumber, message) },
				retry.Attempts(3),
			)
		}

	case common.CHECKIN:
		// If the request type is CHECKIN
		// send an SMS message to all primary contacts using the message in the body
		for _, contact := range user.PrimaryContacts {
			retry.Do(
				func() error { return twilio.SendSMS(contact.PhoneNumber, req.Message) },
				retry.Attempts(3),
			)
		}

	}

	// Return successful response
	bodyContent := fmt.Sprintf("Successfully sent SMS to contacts of user %v %v (%v)", user.FirstName, user.LastName, user.UserID)
	return driver.SuccessfulResponse(bodyContent), nil
}

func main() {
	lambda.Start(Handler)
}
