# Test the function locally using the SAM CLI tool
sam:
	# Building Go binary
	GOOS=linux go build -o main

	# Put the binary into a zip folder
	${GOBIN}/build-lambda-zip.exe -o main.zip main

	# Test
	sam local start-api

test:
	# curl -i -X PATCH http://127.0.0.1:3000/profile/a3cac301-f9a9-4df8-a224-d0bac718e4fa -d @patch_event.json --header "Content-Type: application/json"
	curl -i -X GET http://127.0.0.1:3000/profile/a3cac301-f9a9-4df8-a224-d0bac718e4fa