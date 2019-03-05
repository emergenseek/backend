package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/emergenseek/backend/common/models"
)

//GetUserDB to read user data
func GetUserDB() {

	config := &aws.Config{Region: aws.String("us-east-2")}

	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)

	userKey := models.User{
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

	key, err := dynamodbattribute.MarshalMap(userKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String("EmergenSeekUsers"),
	}

	result, err := svc.GetItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	user := models.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	/*
		infoMap := movie.Info.(map[string]interface{})
		for k, v := range infoMap {
			switch vv := v.(type) {
			case string, float64:
				fmt.Println(k, ": ", vv)
			case []interface{}:
				for i, u := range vv {
					fmt.Println(i, u)
				}
			default:
				fmt.Println(k, "is of a type I don't know how to handle")
			}
		}
	*/
}

func main() {
	GetUserDB()
}
