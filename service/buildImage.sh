#!/bin/bash -e

log() {
  echo "====================>" $1
}

buildService() {
  # e.g. blogrpc-member => member
  SERVICE=${1:8}
  REGISTRY_HOST=${REGISTRY_HOST:-"registry.ap-southeast-1.aliyuncs.com"}
  # IMAGE_NAME is the image name for openapi
  IMAGE_NAME="${REGISTRY_HOST}/yhhnamespace/blogrpc-${SERVICE}"

  case ${SERVICE} in
    "cron")
      log "Build cron bin: todo"
      ;;
    *)
      log "Build ${SERVICE} bin"
      ./scripts/build bin ${SERVICE}
      log "Build ${SERVICE} image"
      docker build --build-arg REGISTRY_HOST=${REGISTRY_HOST} --build-arg ENV=${ENV} --build-arg rpcname=${SERVICE} -t $IMAGE_NAME:${ENV} -f ./docker/Dockerfile .
      ;;
  esac

  docker push $IMAGE_NAME:${ENV}
}

export ENV=${ENV:-local}
#export SERVICES=blogrpc-member,blogrpc-hello
export SERVICES=blogrpc-member
# 用空格替换','来作为分隔符
SERVICES=${SERVICES//,/ }

for service in ${SERVICES}; do
  buildService ${service}
done
