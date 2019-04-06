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

// Handler is the Lambda handler for ESAddContact
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Create a new request object and unmarshal the request body into it
	var contact models.Contact
	err := json.Unmarshal([]byte(request.Body), &contact)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}
	userID := request.PathParameters["user_id"]

	// Initialize drivers
	db, _, _, _ := driver.CreateAll()

	// Add the contact to the database
	cerr := db.AddContact(userID, &contact)
	if cerr != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, cerr), nil
	}

	bodyContent := fmt.Sprintf("Successfully associated contact to user %v", userID)
	return driver.SuccessfulResponse(bodyContent), nil

}

func main() {
	lambda.Start(Handler)
}
