#!/bin/bash -e

PROJECT_PATH=$(cd `dirname ${0}`; pwd)
BUSINESS_PATH="${PROJECT_PATH}/openapi/business"
BLOGRPC_PATH="${PROJECT_PATH}/service"

main() {
  echo "=================================================================="
  echo "开始启动服务。"
  echo "=================================================================="
  case $1 in
  service)
    ${BLOGRPC_PATH}/scripts/build ${@:2}
    ;;
  business)
    ${BUSINESS_PATH}/scripts/build ${@:2}
    ;;
  *)
    usage
    ;;
  esac
}

usage() {
  echo "USAGE: $0 option"
  echo -e "\nOptions:"
  echo "    service"
  exit 1
}

main $@