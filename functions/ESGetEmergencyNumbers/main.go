package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/emergenseek/backend/common/driver"
)

// Handler is the Lambda handler for ESGetEmergencyInfo
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize drivers
	db, _, _, mapsKey := driver.CreateAll()

	// Extract latitude and longitude from path and attempt to convert to float64
	lat, err := strconv.ParseFloat(request.PathParameters["latitude"], 64)
	if err != nil {
		return driver.ErrorResponse(http.StatusUnprocessableEntity, err), nil
	}

	lng, err := strconv.ParseFloat(request.PathParameters["longitude"], 64)
	if err != nil {
		return driver.ErrorResponse(http.StatusUnprocessableEntity, err), nil
	}

	// Find the matching country code
	countryCode, cerr := driver.GetCountryCode([]float64{lat, lng}, mapsKey)
	if cerr != nil {
		return driver.ErrorResponse(http.StatusBadRequest, err), nil
	}

	// Retrieve data from database
	result, err := db.GetEmergencyNumbers(countryCode)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	// Report to client
	b, err := json.Marshal(result)
	if err != nil {
		return driver.ErrorResponse(http.StatusInternalServerError, err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(b),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
