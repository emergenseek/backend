package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/emergenseek/backend/common/driver"
	"github.com/emergenseek/backend/common/models"
)

type lockScreenInfoResponse struct {
	FirstName        string          `json:"first_name,omitempty"`
	LastName         string          `json:"last_name,omitempty"`
	BloodType        string          `json:"blood_type,omitempty"`
	Age              uint32          `json:"age,omitempty"`
	PrimaryResidence *models.Address `json:"primary_residence,omitempty"`
	PhonePin         uint64          `json:"phone_pin,omitempty"`
	EmailAddress     string          `json:"email_address,omitempty"`
	PhoneNumber      string          `json:"phone_number,omitempty"`
}

var headers = map[string]string{"Content-Type": "application/json"}

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

	// Retrieve user from database
	user, err := db.GetUser(req.UserID)

	if err != nil {
		return driver.ErrorResponse(http.StatusBadRequest, err), nil
	}

	responseInfo := lockScreenInfoResponse{
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		BloodType:        user.BloodType,
		Age:              user.Age,
		PrimaryResidence: user.PrimaryResidence,
		PhonePin:         user.PhonePin,
		EmailAddress:     user.EmailAddress,
		PhoneNumber:      user.PhoneNumber,
	}

	b, err := json.Marshal(responseInfo)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(b),
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
