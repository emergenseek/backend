package models

// Settings defines an active user's settings
type Settings struct {
	UserID            string `json:"user_id"`          // Foreign key association of the user's who's settings are being stored
	SOSSMS            bool   `json:"sos_sms"`          // Enable/Disable SOS SMS
	SOSCalls          bool   `json:"sos_calls"`        // Enable/Disable SOS Calls
	SOSLockscreenInfo bool   `json:"sos_lockscreen"`   // Enable/Disable SOS Lockscreen Information
	Updates           bool   `json:"updates"`          // Enable/Disable location updates via poller
	UpdateFrequency   int    `json:"update_frequency"` // Number of hours the poll notification is scheduled for
}
