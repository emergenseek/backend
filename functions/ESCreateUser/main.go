package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/emergenseek/backend/common/driver"
	"github.com/emergenseek/backend/common/models"
)

// Handler is the Lambda handler for ESCreateUser
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Create a new request object and unmarshal the request body into it
	var profile models.User
	err := json.Unmarshal([]byte(request.Body), &profile)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	// Initialize drivers
	db, _, _, _ := driver.CreateAll()

	_, err = db.CreateUser(&profile)

	return driver.SuccessfulResponse(fmt.Sprintf("Successfully created user for user %v", profile.UserID)), nil

}

func main() {
	lambda.Start(Handler)
}
