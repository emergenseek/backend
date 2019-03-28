# Test the function locally using the SAM CLI tool
sam:
	# Building Go binary
	GOOS=linux go build -o main

	# Put the binary into a zip folder
	${GOBIN}/build-lambda-zip.exe -o main.zip main

	# Test
	sam local start-api
