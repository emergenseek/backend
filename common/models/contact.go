package models

import (
	"github.com/emergenseek/backend/common"
)

// Contact defines the information of a User's contacts
type Contact struct {
	PhoneNumber  string           `json:"phone_number"`
	Relationship string           `json:"relationship"`
	FirstName    string           `json:"first_name"`
	LastName     string           `json:"last_name"`
	EmailAddress string           `json:"email_address"`
	Tier         common.AlertTier `json:"tier"`
}
