package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/emergenseek/backend/common/database"
	"github.com/emergenseek/backend/common/driver"
	"github.com/emergenseek/backend/common/models"
)

var headers = map[string]string{"Content-Type": "application/json"}

// Driver method for PATCH requests made to ESManageProfile
func processPatch(db *database.DynamoConn, profile *models.User) (events.APIGatewayProxyResponse, error) {
	// Update the user's profile from the database
	err := db.UpdateProfile(profile)
	if err != nil {
		return driver.ErrorResponse(http.StatusNotFound, err), nil
	}

	return driver.SuccessfulResponse(fmt.Sprintf("Successfully updated profile for user %v", profile.UserID)), nil
}

// Driver method for GET requests made to ESManageProfile
func processGet(db *database.DynamoConn, userID string) (events.APIGatewayProxyResponse, error) {
	// Retrieve the user's profile from the database
	result, err := db.GetUser(userID)
	if err != nil {
		return driver.ErrorResponse(http.StatusNotFound, err), nil
	}

	// Convert profile struct into JSON and return to client
	b, err := json.Marshal(result)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(b),
		Headers:    headers,
	}, nil

}

// Handler is the Lambda handler for ESManageProfile
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize drivers
	db, _, _, _ := driver.CreateAll()

	// Switch on request method
	switch request.HTTPMethod {
	case "PATCH":
		// Create a new request object and unmarshal the request body into it
		var profile models.User
		err := json.Unmarshal([]byte(request.Body), &profile)
		if err != nil {
			return driver.ErrorResponse(http.StatusInternalServerError, err), nil
		}
		profile.UserID = request.PathParameters["user_id"]
		return processPatch(db, &profile)
	case "GET":
		return processGet(db, request.PathParameters["user_id"])
	}

	return driver.ErrorResponse(http.StatusMethodNotAllowed, errors.New("method not allowed")), nil

}

func main() {
	lambda.Start(Handler)
}
