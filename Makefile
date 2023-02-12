BINARY_NAME=ip-whois
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -maxdepth 1 -type f -name '*.go')

build:
	mkdir -p bin/
	go build -o bin/${BINARY_NAME} $(SOURCES)

run:
	go run main.go ipdb.go routes.go envutil.go

compile:
	echo "Compiling for every OS and Platform"
	mkdir -p bin/
	GOOS=linux GOARCH=arm go build -o bin/${BINARY_NAME}-linux-arm $(SOURCES)
	GOOS=linux GOARCH=arm64 go build -o bin/${BINARY_NAME}-linux-arm64 $(SOURCES)

clean:
	go clean
	rm bin/${BINARY_NAME}-linux-arm
	rm bin/${BINARY_NAME}-linux-arm64

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
