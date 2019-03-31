table:
	aws dynamodb create-table --region us-east-2 --cli-input-json file://emergenseekusers_table.json

install:
	go get -u github.com/aws/aws-sdk-go/aws
	go get -u github.com/aws/aws-sdk-go/aws/session
	go get -u github.com/aws/aws-sdk-go/service/dynamodb
	go get -u github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute
	go get -u github.com/aws/aws-lambda-go/lambda
	go get -u github.com/aws/aws-lambda-go/events
	go get -u github.com/sfreiberg/gotwilio
	go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
	go get -u github.com/beevik/etree
	go get -u github.com/avast/retry-go
fmt:
	gofmt -s -w .
	goreportcard-cli -v -t 90
