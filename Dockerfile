# Create build stage based on buster image
FROM golang:1.20.6-bullseye

ARG BINARY_NAME=gotagv
# Create working directory under /app
WORKDIR /go/src/github.com/triabokon/${BINARY_NAME}/

# copy module cache
COPY .go* /go/

COPY go.mod go.sum ./
COPY Makefile ./main.go ./
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go mod download
RUN make build


# Make sure to expose the port the HTTP server is using
EXPOSE 8080

# Run the app binary when we run the container
ENTRYPOINT ["./bin/gotagv"]