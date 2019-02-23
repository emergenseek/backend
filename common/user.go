package common

// User encapsulates the data necessary to represent the data of a single EmergenSeek User
type User struct {
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	BloodType         string     `json:"last_name"`
	Age               uint32     `json:"age"`
	PrimaryContacts   *[]Contact `json:"primary_contacts"`
	SecondaryContacts *[]Contact `json:"secondary_contacts"`
	LastKnownLocation []float32  `json:"location"`
	PrimaryResidence  *Address   `json:"primary_residence"`
	PhonePin          uint64     `json:"phone_pin"`
	CognitoID         string     `json:"cognito_id"`
	EmailAddress      string     `json:"email_address"`
}
