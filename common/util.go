package common

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/emergenseek/backend/common/database"
	"github.com/emergenseek/backend/common/notification"
)

// ClientError simplifies the sending of errors to the client from the API
func ClientError(status int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, err
}

// InitAll initializes the necessary API providers for Lambda handlers
func InitAll() (*database.DynamoConn, *notification.TwilioHandler, error) {
	// Initalize database
	db := &database.DynamoConn{Region: Region}
	err := db.Init()
	if err != nil {
		return nil, nil, err
	}

	// Get Twilio client credentials using database
	t := &notification.TwilioHandler{}
	err = t.GetCredentials(db)
	if err != nil {
		return nil, nil, err
	}
	// Authenticate using credentials
	err = t.Authenticate()
	if err != nil {
		return nil, nil, err
	}

	// Return for handler
	return db, t, nil
}
