FROM gcr.io/distroless/static:latest

# Set hardcoded environment variables
ENV DOCKER_REPO="ghcr.io/synkube/app"
ENV OWNER="synkube"
ENV VCS_URL="https://github.com/synkube/app"
ENV DESCRIPTION="A simple Go web server that prints a message every 5 seconds."

# Add labels
LABEL maintainer="Borislav Grigorov <b.s.grigorov@gmail.com>"
LABEL org.label-schema.description="${DESCRIPTION}"
LABEL org.label-schema.name="${DOCKER_REPO}/${APP_NAME}"
LABEL org.label-schema.schema-version="1.0.0"
LABEL org.label-schema.vcs-url="${VCS_URL}"
LABEL org.opencontainers.image.source="${VCS_URL}"
LABEL org.opencontainers.image.description="${DESCRIPTION}"

USER nonroot:nonroot

WORKDIR /app

# Custom to this application
ARG APP_NAME=app1

COPY --chown=nonroot:nonroot ${APP_NAME} /app/${APP_NAME}

ENTRYPOINT ["/app/${APP_NAME}"]
