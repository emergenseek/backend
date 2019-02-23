table:
	aws dynamodb create-table --region us-east-2 --cli-input-json file://emergenseekusers_table.json

install:
	go get -u github.com/aws/aws-sdk-go/aws
	go get -u github.com/aws/aws-sdk-go/aws/session
	go get -u github.com/aws/aws-sdk-go/service/dynamodb
	go get -u github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute

fmt:
	gofmt -s -w .