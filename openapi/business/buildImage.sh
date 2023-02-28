#!/bin/bash -e
# Build and push image to docker registry

# Constant
#
# REGISTRY_HOST is docker registry host
REGISTRY_HOST=${REGISTRY_HOST:-"registry.ap-southeast-1.aliyuncs.com"}
# IMAGE_NAME is the image name for openapi
IMAGE_NAME="${REGISTRY_HOST}/yhhnamespace/blogrpc-openapi-business"
# ENV is the current environment
ENV=${ENV:-local}

echo 'Run go build'
./scripts/build bin

echo 'Build openapi image'
docker build --build-arg REGISTRY_HOST=${REGISTRY_HOST} --build-arg ENV=${ENV} -t "${IMAGE_NAME}:${ENV}" -f docker/Dockerfile .
docker push "${IMAGE_NAME}:${ENV}"
