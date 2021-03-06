package common

// EmergencyType defines integer constants for the types of emergencies that the system expects
type EmergencyType int

// AlertTier defines integer constants for the types of alert tiers available for contacts
type AlertTier int

const (
	// UsersTableName defines the DynamoDB table used to store EmergenSeek users
	UsersTableName = "EmergenSeekUsers"

	// LambdaSecretsTable defines the DynamoDB table used to store environment variables
	LambdaSecretsTable = "LambdaSecrets"

	// SettingsTableName defines the DynamoDB table used to store user settings
	SettingsTableName = "EmergenSeekSettings"

	// EmergencyNumsTableName defines the DynamoDB table used to global emergency phone numbers
	EmergencyNumsTableName = "EmergencyNumbers"

	// S3BucketName defines the name of the S3 bucket used by the application
	S3BucketName = "emergenseek.com"

	// S3BucketLocation defines the URL of the S3 bucket
	S3BucketLocation = "https://s3.us-east-2.amazonaws.com/emergenseek.com/"

	// Region defines the AWS VPC region used for development
	Region = "us-east-2"

	// SEVERE defines a priority 1 emergency
	SEVERE EmergencyType = 1

	// MILD defines a priority 2 emergency
	MILD EmergencyType = 2

	// CHECKIN defines a priority 3 emergency (non-emergency)
	CHECKIN EmergencyType = 3

	// FIRST defines a first priority alert tier
	FIRST AlertTier = 1

	// SECOND defines a second priority alert tier
	SECOND AlertTier = 2

	// THIRD defines a third priority alert tier
	THIRD AlertTier = 3

	// TwilioTrial defines the LambdaSecrets item ID of Twilio credentials for the trial account
	TwilioTrial = "1"

	// TwilioProduction defines the LambdaSecrets item ID of Twilio credentials for the paid account
	TwilioProduction = "2"

	// MapQuest defines the LambdaSecrets item ID containing MapQuest API credentials
	MapQuest = "3"

	// GoogleMaps defines the LambdaSecrets item ID containing Google Maps API credentials
	GoogleMaps = "4"

	// TwilioVoice defines the Amazon Polly voice used for voice calls
	TwilioVoice = "Polly.Kimberly"
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

// String converts an AlertTier to its string name
func (tier AlertTier) String() string {
	// Map the tier type string to an index
	types := [...]string{
		"FIRST",
		"SECOND",
		"THIRD",
	}

	// Check if the integer is between 1 and 3 inclusive
	if tier < FIRST || tier > THIRD {
		return "Unknown"
	}
	// Return at types index - 1 because the enum begins at 1 not 0
	return types[tier-1]
}
