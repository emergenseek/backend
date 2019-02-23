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

	// SEVERE defines a priority 0 emergency
	SEVERE EmergencyType = 1

	// MILD defines a priority 1 emergency
	MILD EmergencyType = 2

	// CHECKIN defines a priority 2 emergency (non-emergency)
	CHECKIN EmergencyType = 3
)
