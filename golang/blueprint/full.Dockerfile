# Use the official Go image as the base image
FROM --platform=linux/amd64 golang:1.21.11-alpine3.20 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Set the application name
ARG APP_NAME=blueprint

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w -extldflags "-static"' -o ${APP_NAME}

###################################################################################################
# Use a minimal base image for the final container
FROM --platform=linux/amd64 alpine:3.20

# Set the working directory inside the container
WORKDIR /app

COPY version.txt .

# Set hardcoded environment variables
ENV GHCR_REPO="ghcr.io/synkube/app"
ENV OWNER="synkube"
ENV VCS_URL="https://github.com/synkube/app"

# Read the version from version.txt and set it as an environment variable
RUN VERSION=$(cat version.txt) && \
  echo $VERSION > /version.txt

# Add labels
LABEL maintainer="Borislav Grigorov <b.s.grigorov@gmail.com>"
LABEL org.label-schema.description="A simple Go web server that prints a message every 5 seconds."
LABEL org.label-schema.name="${GHCR_REPO}/${APP_NAME}"
LABEL org.label-schema.schema-version="${VERSION}"
LABEL org.label-schema.vcs-url="${VCS_URL}"
LABEL org.opencontainers.image.source="${VCS_URL}"
LABEL org.opencontainers.image.description="A simple Go web server that prints a message every 5 seconds."

# Copy the binary from the build stage to the final container
ARG APP_NAME=blueprint

ENV APP_NAME=${APP_NAME}
COPY --from=build /app/${APP_NAME} .

# Run the Go application
CMD ["./${APP_NAME}"]
