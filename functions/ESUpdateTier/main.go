package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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

// Handler is the Lambda handler for ESUpdateTier
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request
	req, status, err := verifyRequest(request)
	if err != nil {
		return driver.ErrorResponse(status, err), nil
	}
	userID := request.PathParameters["user_id"]

	// Initialize drivers
	db, _, _, _ := driver.CreateAll()

	// Update the contact's tier
	uerr := db.UpdateTier(userID, req.PhoneNumber, req.NewTier)
	if uerr != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, uerr), nil
	}

	bodyContent := fmt.Sprintf("Successfully updated tier for contact of user %v", userID)
	return driver.SuccessfulResponse(bodyContent), nil

}

func main() {
	lambda.Start(Handler)
}
