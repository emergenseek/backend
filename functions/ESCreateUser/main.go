package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/emergenseek/backend/common/models"
)

// Request defines the expected body parameters of an ESSendEmergencyVoiceCall invocation
type CreateUser struct {
	// The verified user making the change
	CognitoID string `json:"cognito_id"`
}

function main() {

}
//func (c *DynamoDB) ListTables(input *ListTablesInput) (*ListTablesOutput, error)
