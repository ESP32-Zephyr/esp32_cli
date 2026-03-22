ARTIFACT_NAME := esp32_cli

build:
	@go build -o bin/${ARTIFACT_NAME} cmd/${ARTIFACT_NAME}/main.go

run:
	@go run  cmd/${ARTIFACT_NAME}/main.go
