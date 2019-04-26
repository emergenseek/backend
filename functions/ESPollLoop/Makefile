# Test the function locally using the SAM CLI tool
sam:
	# Building Go binary
	GOOS=linux go build -o main

	# Put the binary into a zip folder
	${GOBIN}/build-lambda-zip.exe -o main.zip main

	# Test
	# sam local start-api 
	sam local invoke

test:
	curl -i -X POST http://127.0.0.1:3000/loop -d @event.json --header "Content-Type: application/json"