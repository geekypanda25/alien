TARGET=alien

all: test build

build:
	go mod init alien
	@go build -o ./$(TARGET)

test:
	@go test -v ./...