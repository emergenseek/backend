package common

// EmergencyType defines integer constants for the types of emergencies that the system expects
type EmergencyType int

const (
	// UsersTableName defines the DynamoDB table used to store EmergenSeek users
	UsersTableName = "EmergenSeekUsers"

	// LambdaSecretsTable defines the DyamoDB table used to store environment variables
	LambdaSecretsTable = "LambdaSecrets"

	// Region defines the AWS VPC region used for development
	Region = "us-east-2"

	// SEVERE defines a priority 1 emergency
	SEVERE EmergencyType = 1

	// MILD defines a priority 2 emergency
	MILD EmergencyType = 2

	// CHECKIN defines a priority 3 emergency (non-emergency)
	CHECKIN EmergencyType = 3

	// TwilioTrial defines the LambdaSecrets item ID of Twilio credentials for the trial account
	TwilioTrial = "1"

	// TwilioProduction defines the LambdaSecrets item ID of Twilio credentials for the paid account
	TwilioProduction = "2"
)

// String converts an EmergencyType to its string name
func (emergency EmergencyType) String() string {
	// Map the emergency type string to an index
	types := [...]string{
		"SEVERE",
		"MILD",
		"CHECKIN",
	}

	// Check if the integer is between 1 and 3 inclusive
	if emergency < SEVERE || emergency > CHECKIN {
		return "Unknown"
	}
	// Return at types index - 1 because the enum begins at 1 not 0
	return types[emergency-1]
}
