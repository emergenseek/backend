# Test the function locally using the SAM CLI tool
sam:
	# Build
	#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o main *.go

	# Test
	sam local start-api
