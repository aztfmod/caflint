build:
	go build -o bin/caflint

test:
	go test ./... -v -cover

cover:
	go test ./... -coverprofile=coverage.out	
	go tool cover -html=coverage.out
