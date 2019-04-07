package notification

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/emergenseek/backend/common"
	"github.com/emergenseek/backend/common/database"
	"github.com/sfreiberg/gotwilio"
)

// TwilioHandler encapsulates credentials and methods for using Twilio's API
type TwilioHandler struct {
	AccountSID   string `json:"TWILIO_ACCOUNT_SID"`
	AuthToken    string `json:"TWILIO_AUTH_TOKEN"`
	TwilioNumber string `json:"TWILIO_NUMBER"`
	TargetNumber string `json:"TWILIO_TARGET"`
	Client       *gotwilio.Twilio
}

// GetCredentials instantiates the struct with credentials from DynamoDB
func (t *TwilioHandler) GetCredentials(db *database.DynamoConn) error {
	if db.Client == nil {
		return errors.New("db: please run DynamoConn.Init() first")
	}

	// Check if in CI to prevent waste of production credits
	// IDs 1 and 2 are resevered for static Twilio credentials
	// ID 2 is a non-trial account and should only be used in production
	var index string
	if os.Getenv("PARTITION") == "aws" {
		index = common.TwilioTrial
	} else {
		index = common.TwilioProduction
	}

	result, err := db.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(common.LambdaSecretsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(index),
			},
		},
	})
	if err != nil {
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &t)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates the struct with a Twilio client
func (t *TwilioHandler) Authenticate() error {
	if t.Client != nil {
		return nil
	}

	t.Client = gotwilio.NewTwilioClient(t.AccountSID, t.AuthToken)
	return nil
}

// SendSMS provides EmergenSeek with SMS notification functionality
func (t *TwilioHandler) SendSMS(phoneNumber string, message string) error {
	_, _, err := t.Client.SendSMS(t.TwilioNumber, phoneNumber, message, "", "")
	if err != nil {
		return err
	}
	return nil
}

// SendVoiceCall provides EmergenSeek with voice call functionality
func (t *TwilioHandler) SendVoiceCall(phoneNumber string, callbackURL string) error {
	callbackParams := gotwilio.NewCallbackParameters(callbackURL)
	callbackParams.Method = "GET"
	_, _, err := t.Client.CallWithUrlCallbacks(t.TwilioNumber, phoneNumber, callbackParams)
	if err != nil {
		return err
	}
	return nil
}
