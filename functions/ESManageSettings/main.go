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

// Driver method for PATCH requests made to ESManageSettings
func processPatch(db *database.DynamoConn, settings *models.Settings) (events.APIGatewayProxyResponse, error) {
	// Update the user's settings from the database
	err := db.UpdateSettings(settings)
	if err != nil {
		return driver.ErrorResponse(http.StatusNotFound, err), nil
	}

	return driver.SuccessfulResponse(fmt.Sprintf("Successfully updated settings for user %v", settings.UserID)), nil
}

// Driver method for GET requests made to ESManageSettings
func processGet(db *database.DynamoConn, userID string) (events.APIGatewayProxyResponse, error) {
	// Retrieve the user's settings from the database
	result, err := db.GetSettings(userID)
	if err != nil {
		return driver.ErrorResponse(http.StatusNotFound, err), nil
	}

	// Convert settings struct into JSON and return to client
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

// Handler is the Lambda handler for ESManageSettings
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize drivers
	db, _, _, _ := driver.CreateAll()

	// Switch on request method
	switch request.HTTPMethod {
	case "PATCH":
		// Create a new request object and unmarshal the request body into it
		var settings *models.Settings
		err := json.Unmarshal([]byte(request.Body), settings)
		if err != nil {
			return driver.ErrorResponse(http.StatusInternalServerError, err), nil
		}
		settings.UserID = request.PathParameters["user_id"]
		return processPatch(db, settings)
	case "GET":
		return processGet(db, request.PathParameters["user_id"])
	}

	return driver.ErrorResponse(http.StatusMethodNotAllowed, errors.New("method not allowed")), nil

}

func main() {
	lambda.Start(Handler)
}
