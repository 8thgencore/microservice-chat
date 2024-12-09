# Use the official Golang image as the base for building the application
FROM golang:1.23.4-alpine3.20 AS builder

# Argument for the environment to be passed during the build
ARG ENV

# Update and upgrade the Alpine packages, then install 'make'
RUN apk update && apk upgrade --available && \
    apk add make && \
    # Create a new user 'chat' with specific parameters
    adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "10001" \
    "chat"

WORKDIR /opt/app/

# Copy the go.mod and go.sum first to install dependencies
COPY go.mod go.sum ./

# Download the Go module dependencies and verify them
RUN go mod download && go mod verify

# Copy the entire application code into the working directory
COPY . .

# Build the application using the 'make' command, passing the environment as a variable
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/main cmd/chat/main.go

###########
# 2 stage #
###########
# Use a minimal base image to run the application
FROM scratch

# Set the working directory in the new image
WORKDIR /opt/app/

# Copy the passwd and group files from the builder stage for the user 'chat'
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the compiled binary and configuration file from the builder stage
# Ensure the ownership is set to the 'chat' user and group
COPY --from=builder --chown=chat:chat /opt/app/bin/main .

# Set the user and group for running the application
USER chat:chat

# Command to run the application with the specified configuration file
ENTRYPOINT ["./main"]
