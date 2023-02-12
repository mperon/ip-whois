############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update \
    && apk add --no-cache 'git=~2' \
    && mkdir -p $GOPATH/src/packages/ip-whois/

# Install dependencies
ENV GO111MODULE=on
WORKDIR $GOPATH/src/packages/ip-whois/
COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/ip-whois *.go

############################
# STEP 2 build a small image
############################
FROM alpine:3

WORKDIR /

# Copy our static executable.
COPY --from=builder /go/ip-whois /go/ip-whois
#COPY public /go/public

ENV PORT 4444
ENV GIN_MODE release
ENV IPWHOIS_UPDATE_URL=http://ftp.registro.br/pub/numeracao/origin/nicbr-asn-blk-latest.txt
ENV IPWHOIS_UPDATE_INTERVAL=24h
ENV IPWHOIS_PORT=4444
EXPOSE 4444

WORKDIR /go

# Run the Go Gin binary.
ENTRYPOINT ["/go/ip-whois"]
