package database

// Reference: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-read-table-item.html

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/emergenseek/backend/common/models"
)

// DynamoConn encapsulates data and methods necessary for communicating with DynamoDB
type DynamoConn struct {
	Client *dynamodb.DynamoDB
	Region string
}

// Init initializes a DynamoDB session
func (d *DynamoConn) Init() error {
	// Assume client is may already be authorized
	if d.Client != nil {
		return nil
	}

	// Initialize session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(d.Region)},
	)

	if err != nil {
		return err
	}

	// Create DynamoDB client using session
	d.Client = dynamodb.New(sess)
	return nil
}

// GetUser will retrieve a user from Cognito
func GetUser(uid string) *models.User {
	return nil
}
