package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	"github.com/jasonwinn/geocoder" // For geocoding
	"googlemaps.github.io/maps"     // For Google Places API
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

// SuccessfulResponse prepares and sends a successful server response for the calling Lambda function
func SuccessfulResponse(bodyContent string) events.APIGatewayProxyResponse {
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
func CreateEmergencyMessage(emergency common.EmergencyType, user *models.User, mapsKey string, location []float64) (string, error) {
	name := user.FormattedName()
	address, err := GetAddress(location, mapsKey)
	if err != nil {
		return "", err
	}

	message := fmt.Sprintf("%v has just triggered a level %d emergency (%v). ", name, emergency, emergency.String())
	message = message + fmt.Sprintf("Their last known location is %v. ", address)
	message = message + fmt.Sprintf("Please contact them at %v to ensure their safety. -EmergenSeek", user.PhoneNumber)
	return message, nil
}

// GetAddress is used to ReverseGeocode a latlng combination into a precise address
func GetAddress(latlng []float64, key string) (string, error) {
	geocoder.SetAPIKey(key)
	a, err := geocoder.ReverseGeocode(latlng[0], latlng[1])
	if err != nil {
		return "", err
	}

	// Format address parts
	address := fmt.Sprintf("%v, ", a.Street)
	if a.City != "" {
		address = address + fmt.Sprintf("%v, ", a.City)
	}
	if a.State != "" {
		address = address + fmt.Sprintf("%v, ", a.State)
	}
	if a.PostalCode != "" {
		address = address + fmt.Sprintf("%v, ", a.PostalCode)
	}
	if a.CountryCode != "" {
		address = address + fmt.Sprintf("%v", a.CountryCode)
	}
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

// CreatePollMessage creates the
func CreatePollMessage(user *models.User, mapsKey string, location []float64) (string, error) {
	name := user.FormattedName()
	address, err := GetAddress(location, mapsKey)
	loc, _ := time.LoadLocation("UTC")
	if err != nil {
		return "", err
	}

	message := fmt.Sprintf("Location Update from %v! ", name)
	message = message + fmt.Sprintf("Location: %v. ", address)
	message = message + fmt.Sprintf("Date & Time: %v UTC.", time.Now().In(loc).Format("Mon 01-02-2006 15:04:05"))
	return message, nil
}

// GetEmergencyServices retrives hospitals and pharmacies within a 10 mile radius of the given location
func GetEmergencyServices(location []float64, db *database.DynamoConn) (string, error) {
	// Every element in the response body will have these attributes
	type LocationItem struct {
		Location maps.LatLng `json:"location"`
		Name     string      `json:"name"`
		Icon     string      `json:"icon"`
		Open     bool        `json:"open"`
	}

	// Retrive maps key
	gmapsKey := db.MustGetGMapsKey()

	// Authenticate a new Maps API client
	c, err := maps.NewClient(maps.WithAPIKey(gmapsKey))
	if err != nil {
		return "", err
	}

	// Create a search for nearby places using the current location
	// Reference: https://github.com/googlemaps/google-maps-services-go/blob/master/places.go#L134
	pharmacyRequest := &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: location[0],
			Lng: location[1],
		},
		Radius:  uint(20), // 10 mile radius
		Keyword: "pharmacy",
	}

	// Make request to Google Places API
	pharmacyResponse, err := c.NearbySearch(context.Background(), pharmacyRequest)
	if err != nil {
		return "", err
	}

	// Extract only necessary attributes from pharmacy query
	locationItems := []LocationItem{}
	open := false
	for _, i := range pharmacyResponse.Results {
		if i.OpeningHours != nil {
			if i.OpeningHours.OpenNow != nil {
				open = *i.OpeningHours.OpenNow
			}
		}
		item := LocationItem{
			Location: i.Geometry.Location,
			Name:     i.Name,
			Icon:     i.Icon,
			Open:     open,
		}
		locationItems = append(locationItems, item)
	}
	fmt.Printf("%+v", locationItems)

	// Do the same thing for hospital query
	hospitalRequest := &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: location[0],
			Lng: location[1],
		},
		Radius:  uint(20), // 10 mile radius
		Keyword: "hospital",
	}

	// Make request to Google Places API
	hospitalResponse, err := c.NearbySearch(context.Background(), hospitalRequest)
	if err != nil {
		return "", err
	}

	// continue appending to locationItems slice previously declared
	open = false
	for _, i := range hospitalResponse.Results {
		if i.OpeningHours != nil {
			if i.OpeningHours.OpenNow != nil {
				open = *i.OpeningHours.OpenNow
			}
		}
		item := LocationItem{
			Location: i.Geometry.Location,
			Name:     i.Name,
			Icon:     i.Icon,
			Open:     open,
		}
		locationItems = append(locationItems, item)
	}
	fmt.Printf("%+v", locationItems)

	j, err := json.Marshal(locationItems)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
