package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/emergenseek/backend/common"
	"github.com/emergenseek/backend/common/models"
)

const (
	u = "Unknown"
)

func main() {
	// numbers.json is a alightly modified version of
	// https://github.com/Alex0x47/EmergencyNumbers/blob/master/emergency_numbers.js
	jsonFile, err := os.Open("numbers.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(common.Region)}))
	client := dynamodb.New(sess)

	for k, v := range result {
		police := v["Police"]
		if police == "" {
			police = u
		}

		ambulance := v["Ambulance"]
		if ambulance == "" {
			ambulance = u
		}

		fire := v["Fire"]
		if fire == "" {
			fire = u
		}

		data := &models.EmergencyInfo{
			CountryCode: k,
			Police:      police,
			Ambulance:   ambulance,
			Fire:        fire,
		}

		// Marshal user struct into map for DynamoDB
		item, err := dynamodbattribute.MarshalMap(data)
		if err != nil {
			panic(err)
		}

		// Use marshalled map for PutItemInput
		input := &dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String(common.EmergencyNumsTableName),
		}

		// Insert into database
		_, perr := client.PutItem(input)
		if perr != nil {
			panic(perr)
		}

		fmt.Println(data)

	}
}
