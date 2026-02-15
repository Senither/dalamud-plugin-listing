FROM node:22-alpine AS node-build

# Set the working directory
WORKDIR /app

# Copy everything into the working directory
COPY . ./

# Install the dependencies
RUN npm install

# Build the assets
RUN npm run build

# Setup the Go build stage
FROM golang:1.22 AS go-build

# Enable Go modules
ARG CGO_ENABLED=0

# Set the working directory
WORKDIR /app

# Copy the Go modules and source files
COPY . .

# Install the dependencies
RUN go mod download && go mod verify

# Build the application
RUN go build -o /app/dalamud-plugin-listing

# Setup a lean image to run the application
FROM gcr.io/distroless/base-debian11 AS build-release-stage

# Set the working directory
WORKDIR /app

# Copy the built application and the views
COPY --from=go-build /app/dalamud-plugin-listing /app/dalamud-plugin-listing
COPY --from=go-build /app/repositories.txt /app/repositories.txt
COPY --from=go-build /app/plugins.txt /app/plugins.txt
COPY --from=node-build /app/assets /app/assets
COPY --from=node-build /app/views /app/views

# Expose the port the application runs on
EXPOSE 8080

# Run the application
CMD ["/app/dalamud-plugin-listing"]
