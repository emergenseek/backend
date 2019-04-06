# Test the function locally using the SAM CLI tool
sam:
	# Building Go binary
	GOOS=linux go build -o main

	# Put the binary into a zip folder
	${GOBIN}/build-lambda-zip.exe -o main.zip main

	# Test
	sam local start-api

test:
	curl -i -X POST http://127.0.0.1:3000/contact/4a35788a-e3fa-4eb4-b1c2-b9a3be8a58c9 -d @event.json -H "Content-Type: application/json"