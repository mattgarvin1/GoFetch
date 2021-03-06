FROM golang:1.15-alpine as build

# Commands to build and run the points server as a Docker container
### MacBook-Air:GoFetch matthewgarvin$ docker build -t points .
### MacBook-Air:GoFetch matthewgarvin$ docker run -d -p 8080:8080 points

# Install SSL certificates
RUN apk update && apk add --no-cache git ca-certificates gcc musl-dev

# Build static points binary
RUN mkdir -p /go/src/github.com/mattgarvin1/GoFetch
WORKDIR /go/src/github.com/mattgarvin1/GoFetch
COPY . .
WORKDIR /go/src/github.com/mattgarvin1/GoFetch/points
RUN go build -ldflags "-linkmode external -extldflags -static" -o /points

# Small image
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /points /

# Q. Does this server need any args?
# ENTRYPOINT ["/points", "listen"]
ENTRYPOINT ["/points"]
