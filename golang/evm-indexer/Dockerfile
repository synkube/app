FROM --platform=linux/amd64 alpine:3.20

# Set hardcoded environment variables
ENV IMAGE_REPO="ghcr.io/synkube/app"
ENV OWNER="synkube"
ENV VCS_URL="https://github.com/synkube/app"
ENV DESCRIPTION="EVM indexer"

# Custom to this application
ARG APP_NAME=evm-indexer
ENV APP_NAME=${APP_NAME}

# Add labels
LABEL org.label-schema.description="${DESCRIPTION}"
LABEL org.label-schema.name="${IMAGE_REPO}/${APP_NAME}"
LABEL org.label-schema.schema-version="1.0.0"
LABEL org.label-schema.vcs-url="${VCS_URL}"
LABEL org.opencontainers.image.source="${VCS_URL}"
LABEL org.opencontainers.image.description="${DESCRIPTION}"

WORKDIR /app

COPY ${APP_NAME} /app/${APP_NAME}

# Run the Go application
CMD ["/app/evm-indexer"]
