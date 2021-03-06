package database

// Reference: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-read-table-item.html

import (
	"errors"
	"fmt"

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

// GetSettings retrieves a user's settings from the Settings table
func (d *DynamoConn) GetSettings(uid string) (*models.Settings, error) {
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
		TableName: aws.String(common.SettingsTableName),
	}

	// Search for settings object matching uid in table
	result, err := d.Client.GetItem(input)
	if err != nil {
		return nil, err
	}

	// Unmarshal user into struct
	settings := &models.Settings{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &settings)
	if err != nil {
		return nil, err
	}

	// Cheap check for item not found
	if settings.UserID == "" {
		return nil, errors.New("settings for user not found")
	}
	fmt.Println(settings)
	return settings, nil
}

// UpdateSettings updates an existing user's settings
func (d *DynamoConn) UpdateSettings(settings *models.Settings) error {
	var SettingsUpdate struct {
		SOSSMS            bool `json:":a"`
		SOSCalls          bool `json:":b"`
		SOSLockscreenInfo bool `json:":c"`
		Updates           bool `json:":d"`
		UpdateFrequency   int  `json:":e"`
	}

	// Marshal the update expression struct for DynamoDB
	SettingsUpdate.SOSSMS = settings.SOSSMS
	SettingsUpdate.SOSCalls = settings.SOSCalls
	SettingsUpdate.SOSLockscreenInfo = settings.SOSLockscreenInfo
	SettingsUpdate.Updates = settings.Updates
	SettingsUpdate.UpdateFrequency = settings.UpdateFrequency

	expr, err := dynamodbattribute.MarshalMap(SettingsUpdate)
	if err != nil {
		return err

	}

	// Define table schema's key
	key := map[string]*dynamodb.AttributeValue{
		"user_id": {
			S: aws.String(settings.UserID),
		},
	}

	// Use marshalled map for UpdateItemInput
	item := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(common.SettingsTableName),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set sos_sms = :a, sos_calls = :b, sos_lockscreen = :c, updates = :d, update_frequency = :e"),
	}

	// Invoke the update
	_, err = d.Client.UpdateItem(item)
	if err != nil {
		return err
	}
	return nil
}

// AddContact associates an additional contact to the given user
func (d *DynamoConn) AddContact(userID string, contact *models.Contact) error {
	var ContactsUpdate struct {
		Contacts []*models.Contact `json:":c"`
	}

	// Retrieve user from database
	user, err := d.GetUser(userID)
	if err != nil {
		return err
	}
	// Add the new contact to the slice and marshal the update
	ContactsUpdate.Contacts = append(user.Contacts, contact)

	expr, err := dynamodbattribute.MarshalMap(ContactsUpdate)
	if err != nil {
		return err

	}

	// Define table schema's key
	key := map[string]*dynamodb.AttributeValue{
		"user_id": {
			S: aws.String(user.UserID),
		},
	}

	// Use marshalled map for UpdateItemInput
	item := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(common.UsersTableName),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set contacts = :c"),
	}

	// Invoke the update
	_, err = d.Client.UpdateItem(item)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTier updates an existing contact's tier
func (d *DynamoConn) UpdateTier(userID string, phoneNumber string, tier common.AlertTier) error {
	var ContactsUpdate struct {
		Contacts []*models.Contact `json:":c"`
	}

	// Retrieve user from database
	user, err := d.GetUser(userID)
	if err != nil {
		return err
	}

	// Search for the contact matching the provided phone number
	seen := false
	for _, contact := range user.Contacts {
		if contact.PhoneNumber == phoneNumber {
			contact.Tier = tier
			seen = true
			break
		}
	}

	// Hacky way to check if the contact was found
	if !seen {
		return fmt.Errorf("unable to find contact with phone number %v for user %v", phoneNumber, userID)
	}

	// Good enough until we figure out a way to update a single contact
	ContactsUpdate.Contacts = user.Contacts

	expr, err := dynamodbattribute.MarshalMap(ContactsUpdate)
	if err != nil {
		return err

	}

	// Define table schema's key
	key := map[string]*dynamodb.AttributeValue{
		"user_id": {
			S: aws.String(user.UserID),
		},
	}

	// Use marshalled map for UpdateItemInput
	// TODO: update contacts to be map[string]Contact where
	// the key is the phone number. This way we wont need to update the entire attribute
	item := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(common.UsersTableName),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set contacts = :c"),
	}

	// Invoke the update
	_, err = d.Client.UpdateItem(item)
	if err != nil {
		return err
	}
	return nil
}

// UpdateProfile updates an existing user's profile
func (d *DynamoConn) UpdateProfile(profile *models.User) error {
	var ProfileUpdate struct {
		FirstName        string          `json:":FirstName"`
		LastName         string          `json:":LastName"`
		BloodType        string          `json:":BloodType"`
		Age              uint32          `json:":Age"`
		PrimaryResidence *models.Address `json:":PrimaryResidence"`
		PhonePin         uint64          `json:":PhonePin"`
		EmailAddress     string          `json:":EmailAddress"`
		PhoneNumber      string          `json:":PhoneNumber"`
	}

	// Marshal the update expression struct for DynamoDB
	ProfileUpdate.FirstName = profile.FirstName
	ProfileUpdate.LastName = profile.LastName
	ProfileUpdate.BloodType = profile.BloodType
	ProfileUpdate.Age = profile.Age
	ProfileUpdate.PrimaryResidence = profile.PrimaryResidence
	ProfileUpdate.PhonePin = profile.PhonePin
	ProfileUpdate.EmailAddress = profile.EmailAddress
	ProfileUpdate.PhoneNumber = profile.PhoneNumber

	expr, err := dynamodbattribute.MarshalMap(ProfileUpdate)
	if err != nil {
		return err
	}

	// Define table schema's key
	key := map[string]*dynamodb.AttributeValue{
		"user_id": {
			S: aws.String(profile.UserID),
		},
	}

	// Use marshalled map for UpdateItemInput
	item := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(common.UsersTableName),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set first_name = :FirstName, last_name = :LastName, blood_type = :BloodType, age = :Age, primary_residence = :PrimaryResidence, phone_pin = :PhonePin, email_address = :EmailAddress, phone_number = :PhoneNumber"),
	}

	// Invoke the update
	_, err = d.Client.UpdateItem(item)
	if err != nil {
		return err
	}
	return nil
}

// GetEmergencyNumbers retrieves police, fire, and ambulance numbers, given a country code
func (d *DynamoConn) GetEmergencyNumbers(countryCode string) (*models.EmergencyInfo, error) {
	tableKey := struct {
		CountryCode string `json:"country_code"`
	}{
		countryCode,
	}

	// Create search key and Get request input
	key, err := dynamodbattribute.MarshalMap(tableKey)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(common.EmergencyNumsTableName),
	}

	// Search for item matching countyr code in table
	result, err := d.Client.GetItem(input)
	if err != nil {
		return nil, err
	}

	// Unmarshal information into struct
	data := &models.EmergencyInfo{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
