package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/emergenseek/backend/common"
	"github.com/emergenseek/backend/common/database"
	"github.com/emergenseek/backend/common/models"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(common.Region)}))
	db := database.DynamoConn{Region: common.Region}
	err := db.Create(sess)
	if err != nil {
		fmt.Println(err.Error())
	}

	address := &models.Address{
		Line1:   "1 Main St",
		Line2:   "Apt 1",
		City:    "Houston",
		State:   "Texas",
		Country: "United States of America",
		ZipCode: "77071",
	}
	contact := &models.Contact{
		PhoneNumber:  "+1REDACTED",
		Relationship: "Brother",
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
		HomeAddress:  address,
	}
	user := &models.User{
		UserID:            "b4f2a0b9-5c63-4257-9655-a3ee2b0519a1",
		FirstName:         "John",
		LastName:          "Doe",
		BloodType:         "O+",
		Age:               50,
		PrimaryContacts:   []*models.Contact{contact},
		SecondaryContacts: []*models.Contact{contact},
		LastKnownLocation: []float64{40.7648, -73.9808},
		PrimaryResidence:  address,
		PhonePin:          12345,
		EmailAddress:      "john.doe@example.com",
		PhoneNumber:       "+1REDACTED",
	}

	insertedUser, err := db.CreateUser(user)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v", insertedUser)

	result, err := db.GetUser(user.UserID)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v", result)

	err = db.UpdateLocation(user.UserID, []float64{35.0116, 135.7681})
	if err != nil {
		fmt.Println(err.Error())
	}
}
