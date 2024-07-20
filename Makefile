APP_NAME=server

vet:
	go vet -v ./...
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux/$(APP_NAME) main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows/$(APP_NAME).exe main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/mac/$(APP_NAME) main.go
