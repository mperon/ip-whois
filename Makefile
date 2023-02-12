BINARY_NAME=ip-whois
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -maxdepth 1 -type f -name '*.go')

build:
	mkdir -p bin/
	go build -o bin/${BINARY_NAME} $(SOURCES)

linux:
	mkdir -p bin/
	GOOS=linux GOARCH=arm64 go build -o bin/${BINARY_NAME}-linux-arm64 $(SOURCES)

run:
	go run main.go ipdb.go routes.go envutil.go

dist:
	echo "Compiling for every OS and Platform"
	mkdir -p bin/
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux-amd64 $(SOURCES)
	GOOS=linux GOARCH=arm64 go build -o bin/${BINARY_NAME}-linux-arm64 $(SOURCES)
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}-darwin-amd64 $(SOURCES)
	GOOS=darwin GOARCH=arm64 go build -o bin/${BINARY_NAME}-darwin-arm64 $(SOURCES)
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-windows-amd64.exe $(SOURCES)

clean:
	go clean
	rm -rf bin/*

 test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all
