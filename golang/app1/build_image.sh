#!/usr/bin/env bash
if [[ -z "${DOCKER_REPO}" ]]; then
  export DOCKER_REPO="ghcr.io/synkube/gke/goproxy"
fi
if [[ -z "${IMAGE_TAG}" ]]; then
  export IMAGE_TAG=$(cat version.txt | tr -d '[:space:]')
fi

echo "Building Docker Image: $DOCKER_REPO:$TAG_ID based of $IMAGE_TAG"
docker build -t "$DOCKER_REPO:$IMAGE_TAG" ./src -f "./src/Dockerfile"



