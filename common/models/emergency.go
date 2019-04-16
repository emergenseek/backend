package models

// EmergencyInfo defines the result from a lookup in the EmergencyNumbers table
type EmergencyInfo struct {
	CountryCode string `json:"country_code"`
	Police      string `json:"police"`
	Ambulance   string `json:"ambulance"`
	Fire        string `json:"fire"`
}
