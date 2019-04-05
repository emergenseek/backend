// package main

// import (
// 	"fmt"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb"
// 	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
// 	"github.com/emergenseek/backend/common/models"
// )

// //CreateUserDB to create user data, then READ
// func CreateUserDB() {

// 	config := &aws.Config{Region: aws.String("us-east-2")}

// 	sess := session.Must(session.NewSession(config))

// 	svc := dynamodb.New(sess)

// 	user := models.User{
// 		FirstName: "Anni",
// 		LastName:  "Xcmneo",
// 		BloodType: "Rh+",
// 		Age:       1230,
// 		/*
// 			PrimaryContacts:
// 			SecondaryContacts:
// 			LastKnownLocation:
// 			PrimaryResidence:
// 		*/
// 		PhonePin:     12345,
// 		CognitoID:    "cognido_id",
// 		EmailAddress: "cognido_id@hotmail.com",
// 	}

// 	av, err := dynamodbattribute.MarshalMap(user)

// 	input := &dynamodb.PutItemInput{
// 		Item:      av,
// 		TableName: aws.String("EmergenSeekUsers"),
// 	}
// 	_, err = svc.PutItem(input)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// }

// func main() {
// 	CreateUserDB()
// }

/*
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/emergenseek/backend/common/driver"
)

func verifyRequest(request events.APIGatewayProxyRequest) (*Request, int, error) {
	// Create a new request object and unmarshal the request body into it
	req := new(Request)
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	// Make sure all of the necessary parameters are present
	err = req.Validate()
	if err != nil {
		return nil, http.StatusBadRequest, err

	}
	// All checks passed, return req struct for use. http.StatusOK is ignored
	return req, http.StatusOK, nil
}

// Handler is the Lambda handler for ESSendSMSNotification
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Verify the request
	req, status, err := verifyRequest(request)
	if err != nil {
		return driver.ErrorResponse(status, err), nil
	}

	// Initialize drivers
	db, _, _, _ := driver.CreateAll()

	// Retrieve emergency services near the provided latitude and longitude
	locations, err := driver.GetEmergencyServices(req.Location, db)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}
	fmt.Println(locations)
	// Return successful response containing locations within a 10 mile radius
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       locations,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
*/
