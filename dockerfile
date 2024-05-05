# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the Go app as a static binary.
RUN CGO_ENABLED=0 GOARCH=amd64 go build -o /goapp

# Command to run the executable.
CMD ["/goapp"]
