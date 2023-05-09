## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage ./... && go tool cover -html=coverage

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build_cli: builds the command line tool hiSpeed and copies it to myapp
build_cli:
	@go build -o ../myappDemo/hiSpeed.exe ./cmd/cli