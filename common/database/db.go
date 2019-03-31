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
func (d *DynamoConn) Create(sess *session.Session) error {
	// Assume client is may already be authorized
	if d.Client != nil {
		return errors.New("db: dynamodb client already exists")
	}

	// Create DynamoDB client using session
	d.Client = dynamodb.New(sess)
	return nil
}

// GetUser will retrieve a user from the database
func (d *DynamoConn) GetUser(uid string) (*models.User, error) {
	// Create user struct to be searched for using provided uid
	userKey := &models.User{
		UserID: uid,
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

// MustGetMapsKey retrives the MapQuest API key from the database
func (d *DynamoConn) MustGetMapsKey() string {
	result, err := d.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(common.LambdaSecretsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(common.MapQuest),
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return *result.Item["MAPQUEST_CONSUMER_KEY"].S

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

// UpdateLocation updates the location of a user when a location poll is invoked
func (d *DynamoConn) UpdateLocation(userID string, location []float64) error {
	var LocationUpdate struct {
		LastKnownLocation []float64 `json:":l"`
	}

	// Marshal the update expression struct for DynamoDB
	LocationUpdate.LastKnownLocation = location
	expr, err := dynamodbattribute.MarshalMap(LocationUpdate)
	if err != nil {
		return err

	}

	// Define table schema's key
	key := map[string]*dynamodb.AttributeValue{
		"user_id": {
			S: aws.String(userID),
		},
	}

	// Use marshalled map for UpdateItemInput
	item := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(common.UsersTableName),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set last_known_location = :l"),
	}

	_, err = d.Client.UpdateItem(item)
	if err != nil {
		return err
	}
	return nil
}

// MustGetGMapsKey retrives the Google Maps API
func (d *DynamoConn) MustGetGMapsKey() string {
	result, err := d.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(common.LambdaSecretsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(common.GoogleMaps),
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return *result.Item["MAPS_API_KEY"].S
}
