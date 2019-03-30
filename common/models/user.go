package models

// User encapsulates the data necessary to represent the data of a single EmergenSeek User
type User struct {
	UserID            string     `json:"user_id,omitempty"`
	FirstName         string     `json:"first_name,omitempty"`
	LastName          string     `json:"last_name,omitempty"`
	BloodType         string     `json:"blood_type,omitempty"`
	Age               uint32     `json:"age,omitempty"`
	PrimaryContacts   []*Contact `json:"primary_contacts,omitempty"`
	SecondaryContacts []*Contact `json:"secondary_contacts,omitempty"`
	LastKnownLocation []float64  `json:"last_known_location,omitempty"`
	PrimaryResidence  *Address   `json:"primary_residence,omitempty"`
	PhonePin          uint64     `json:"phone_pin,omitempty"`
	EmailAddress      string     `json:"email_address,omitempty"`
	PhoneNumber       string     `json:"phone_number,omitempty"`
}

// FormattedName formats a users name for SMS messages
func (u *User) FormattedName() string {
	return u.FirstName + " " + u.LastName
}
