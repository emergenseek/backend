package database

// Reference: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-read-table-item.html

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/emergenseek/backend/common"
	"github.com/emergenseek/backend/common/models"
)

// DynamoConn encapsulates data and methods necessary for communicating with DynamoDB
type DynamoConn struct {
	Client *dynamodb.DynamoDB
	Region string
}

// Create creates a new, private DynamoDB session
func (d *DynamoConn) Create() error {
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

// GetUser will retrieve a user from the database
func (d *DynamoConn) GetUser(uid string) (*models.User, error) {
	// Create user struct to be searched for using provided uid
	userKey := &models.User{
		CognitoID: uid,
	}
	key, err := dynamodbattribute.MarshalMap(userKey)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(common.UsersTableName),
	}

	// Search for user matching uid in table
	result, err := d.Client.GetItem(input)
	if err != nil {
		return nil, err
	}

	// Unmarshal user into struct
	user := &models.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	// Cheap check for item not found
	if user.FirstName == "" {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// CreateUser will create a user for the application
func (d *DynamoConn) CreateUser(user *models.User) (*dynamodb.PutItemOutput, error) {
	// Marshal user struct into map for DynamoDB
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	// Use marshalled map for PutItemInput
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(common.UsersTableName),
	}

	// Insert into database
	output, err := d.Client.PutItem(input)
	if err != nil {
		return nil, err
	}

	return output, nil
}
