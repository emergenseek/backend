package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/emergenseek/backend/common/models"
)

//CreateUserDB to create user data, then READ
func CreateUserDB() {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-2")})

	svc := dynamodb.New(sess)

	user := models.User{
		FirstName: "Anni",
		LastName:  "Xcmneo",
		BloodType: "Rh+",
		Age:       1230,
		/*
			PrimaryContacts:
			SecondaryContacts:
			LastKnownLocation:
			PrimaryResidence:
		*/
		PhonePin:     12345,
		CognitoID:    "cognido_id",
		EmailAddress: "cognido_id@hotmail.com",
	}

	av, err := dynamodbattribute.MarshalMap(user)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("EmergenSeekUsers"),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	CreateUserDB()
}
