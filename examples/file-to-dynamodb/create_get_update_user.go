package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/emergenseek/backend/common"
	"github.com/emergenseek/backend/common/database"
	"github.com/emergenseek/backend/common/models"
	"github.com/google/uuid"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(common.Region)}))
	db := database.DynamoConn{Region: common.Region}
	err := db.Create(sess)
	if err != nil {
		fmt.Println(err.Error())
	}

	// simon, derek, kevon, annie
	hostNumber := []string{"+18326720272"}
	for _, phone := range hostNumber {
		address := &models.Address{
			Line1:   "1 Main St",
			Line2:   "Apt 1",
			City:    "Houston",
			State:   "Texas",
			Country: "United States of America",
			ZipCode: "77071",
		}
		tierOneContact := &models.Contact{
			PhoneNumber:  phone,
			Relationship: "Brother",
			FirstName:    "Primary",
			LastName:     "Simon",
			EmailAddress: "simon.woldemichael@ttu.edu",
			Tier:         common.FIRST,
		}

		tierTwoContact := &models.Contact{
			PhoneNumber:  phone,
			Relationship: "Brother",
			FirstName:    "Secondary",
			LastName:     "Simon",
			EmailAddress: "simon.woldemichael@ttu.edu",
			Tier:         common.SECOND,
		}

		tierThreeContact := &models.Contact{
			PhoneNumber:  phone,
			Relationship: "Brother",
			FirstName:    "Tertiary",
			LastName:     "Simon",
			EmailAddress: "simon.woldemichael@ttu.edu",
			Tier:         common.THIRD,
		}

		user := &models.User{
			UserID:            uuid.New().String(),
			FirstName:         "EmergenSeek",
			LastName:          "User",
			BloodType:         "O+",
			Age:               20,
			Contacts:          []*models.Contact{tierOneContact, tierTwoContact, tierThreeContact},
			LastKnownLocation: []float64{35.6795351, 139.7204474},
			PrimaryResidence:  address,
			PhonePin:          12345,
			EmailAddress:      fmt.Sprintf("EmergenSeek.User.%v@example.com", phone),
			PhoneNumber:       phone,
		}
		fmt.Println(user.UserID)

		// Create the user
		insertedUser, err := db.CreateUser(user)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("%+v", insertedUser)

		// Retrieve the user
		result, err := db.GetUser(user.UserID)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("%+v", result)

		// // Update the user's location
		// err = db.UpdateLocation(user.UserID, []float64{35.0116, 135.7681})
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
	}

}
