package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/beevik/etree"
	"github.com/emergenseek/backend/common"
	"github.com/emergenseek/backend/common/database"
	"github.com/emergenseek/backend/common/models"
	"github.com/emergenseek/backend/common/notification"
	"github.com/google/uuid"
	"github.com/jasonwinn/geocoder"
)

var headers = map[string]string{"Content-Type": "application/json"}

// ErrorResponse simplifies the sending of errors to the client from the API
func ErrorResponse(status int, errs ...error) events.APIGatewayProxyResponse {
	errorMessages := []string{}
	for _, err := range errs {
		errorMessages = append(errorMessages, err.Error())
	}

	// Create request body and send to handler
	body, _ := json.Marshal(map[string]string{"code": http.StatusText(status), "error": strings.Join(errorMessages, " | ")})
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers:    headers,
	}
}

// SuccessfulResponse prepares and sends a successful server response for this Lambda function
func SuccessfulResponse(bodyContent string, user *models.User) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"body": bodyContent})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers:    headers,
	}
}

// CreateAll initializes the necessary API providers for Lambda handlers
func CreateAll() (*database.DynamoConn, *notification.TwilioHandler, *session.Session, string) {
	// Create a shared session
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(common.Region)}))

	// Initialize database
	db := &database.DynamoConn{Region: common.Region}
	err := db.Create(sess)
	if err != nil {
		panic(err)
	}

	// Get MapQuest credentials
	mapsKey := db.MustGetMapsKey()

	// Get Twilio client credentials using database
	twilio := &notification.TwilioHandler{}
	err = twilio.GetCredentials(db)
	if err != nil {
		panic(err)
	}
	// Authenticate using credentials
	err = twilio.Authenticate()
	if err != nil {
		panic(err)
	}

	// Return for handler
	return db, twilio, sess, mapsKey
}

// CreateEmergencyMessage generates a message given a user's information and their severity
// Should not used with the CHECKIN emergency type
func CreateEmergencyMessage(emergency common.EmergencyType, user *models.User) string {
	name := user.FormattedName()
	message := fmt.Sprintf(`
		%v has just triggered a level %d emergency (%v). Please contact them at %v to ensure their safety -EmergenSeek
	`, name, emergency, emergency.String(), user.PhoneNumber)
	return message
}

// GetAddress is used to ReverseGeocode a latlng combination into a precise address
func GetAddress(latlng []float64, key string) (string, error) {
	geocoder.SetAPIKey(key)
	a, err := geocoder.ReverseGeocode(latlng[0], latlng[1])
	if err != nil {
		return "", err
	}
	address := fmt.Sprintf("%v, %v, %v, %v, %v", a.Street, a.City, a.State, a.PostalCode, a.CountryCode)
	return address, nil
}

// CreateTwilMLXML creates the XML necessary for the gotwilio.NewCallbackParameters invocation in notification.SendVoiceCall
func CreateTwilMLXML(user *models.User, lastLocation string) ([]byte, error) {
	// Split phone number so Twilio voice doesn't read it numerically
	splitPhoneNumber := fmt.Sprintf("%v", strings.Split(user.PhoneNumber, ""))

	// Create message using address and user's information
	name := user.FormattedName()
	message := fmt.Sprintf("This is an automated emergency call from EmergenSeek on behalf of %v. ", name)
	message = message + fmt.Sprintf("They are in need of emergency assistance. Please send help to %v. ", lastLocation)
	message = message + fmt.Sprintf("Please attempt to call %v at %v. Thank you.", name, splitPhoneNumber)

	// Create, format, and return the XML document
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	response := doc.CreateElement("Response")
	say := response.CreateElement("Say")
	say.CreateAttr("voice", common.TwilioVoice)
	say.CreateAttr("loop", "2")
	say.SetText(message)
	doc.Indent(2)
	twilML, err := doc.WriteToBytes()
	if err != nil {
		return nil, err
	}
	return twilML, nil
}

// UploadTwilMLXML uploads the XML generated in CreateTwilMLXML and returns the URL for the object
func UploadTwilMLXML(twilML []byte, sess *session.Session) (string, error) {
	// Create a unique id for the object
	id, _ := uuid.NewRandom()
	objectKey := string([]rune(id.String())[0:7]) + ".xml"

	// Upload to S3 bucket
	_, err := s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(common.S3BucketName),
		Key:           aws.String(objectKey),
		ACL:           aws.String("public-read"),
		Body:          bytes.NewReader(twilML),
		ContentLength: aws.Int64(int64(len(twilML))),
		ContentType:   aws.String("text/xml"),
	})
	if err != nil {
		return "", err
	}
	return common.S3BucketLocation + objectKey, nil
}
