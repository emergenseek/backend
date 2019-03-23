package notification

import (
	"errors"
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

	result, err := db.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(common.LambdaSecretsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				// IDs 1 and 2 are resevered for static Twilio Credentials
				// ID 2 is a non-trial account and should only be used in production
				N: aws.String(common.TwilioProduction),
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
func (t *TwilioHandler) SendVoiceCall(phoneNumber string) error {
	callbackParams := gotwilio.NewCallbackParameters("https://gist.githubusercontent.com/swoldemi/52535e13478c32af468eeab9d10f4813/raw/1d2a2542399af84209f8011953559665f17cb01f/voice.xml")
	callbackParams.Method = "GET"
	_, _, err := t.Client.CallWithUrlCallbacks(t.TwilioNumber, phoneNumber, callbackParams)
	if err != nil {
		return err
	}
	return nil
}
