package driver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/emergenseek/backend/common"
	"github.com/emergenseek/backend/common/database"
	"github.com/emergenseek/backend/common/models"
	"github.com/emergenseek/backend/common/notification"
)

var headers = map[string]string{"Content-Type": "application/json"}

// ErrorResponse simplifies the sending of errors to the client from the API
func ErrorResponse(status int, errs ...error) events.APIGatewayProxyResponse {
	errorMessages := []string{}
	for _, err := range errs {
		errorMessages = append(errorMessages, err.Error())
	}

	// Create request body and send to handler
	body, _ := json.Marshal(map[string]string{"code": http.StatusText(status), "error": strings.Join(errorMessages, " | ")})
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers:    headers,
	}
}

// SuccessfulResponse prepares and sends a successful server response for this Lambda function
func SuccessfulResponse(bodyContent string, user *models.User) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"body": bodyContent})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers:    headers,
	}
}

// CreateAll initializes the necessary API providers for Lambda handlers
func CreateAll() (*database.DynamoConn, *notification.TwilioHandler) {
	// Initialize database
	db := &database.DynamoConn{Region: common.Region}
	err := db.Create()
	if err != nil {
		panic(err)
	}

	// Get Twilio client credentials using database
	twilio := &notification.TwilioHandler{}
	err = twilio.GetCredentials(db)
	if err != nil {
		panic(err)
	}
	// Authenticate using credentials
	err = twilio.Authenticate()
	if err != nil {
		panic(err)
	}

	// Return for handler
	return db, twilio
}

// CreateEmergencyMessage generates a message given a user's information and their severity
// Should not used with the CHECKIN emergency type
func CreateEmergencyMessage(emergency common.EmergencyType, user *models.User) string {
	name := user.FormattedName()
	message := fmt.Sprintf(`
		%v has just triggered a level %d emergency (%v). Please contact them at %v to ensure their safety -EmergenSeek
	`, name, emergency, emergency.String(), user.PhoneNumber)
	return message
}
