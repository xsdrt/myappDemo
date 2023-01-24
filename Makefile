BINARY_NAME=hiSpeedApp.exe

## build: builds all binaries
build:
	@go mod vendor
	@go build -o tmp/${BINARY_NAME} .
	@echo HiSpeed built!

run: build
	@echo Starting HiSpeed...
	@.\tmp\${BINARY_NAME} &
	@echo HiSpeed started!

clean:
	@echo Cleaning...
	@go clean
	@del .\tmp\${BINARY_NAME}
	@echo Cleaned!

test:
	@echo Testing...
	@go test ./...
	@echo Done!

start: run
    
stop:
	@echo "Starting the front end..."
	@taskkill /IM ${BINARY_NAME} /F
	@echo Stopped HiSpeed

restart: stop start