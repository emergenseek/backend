package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	UsersTableName = "EmergenSeekUsers"
	Region         = "us-east-2"
)

type SMSNotification struct {
	PhoneNumber string
	Message     string
}

// NotifyPrimaryContacts will notify the primary contacts of the given user
// The user is matched using their user ID (uid)
func NotifyPrimaryContacts(uid string) error {
	// Authenticate with Twilio and send the SMS
	return nil
}

func Handler(request events.APIGatewayProxyRequest) (error, error) {
	return nil, nil
}

func main() {
	lambda.Start(Handler)
}
