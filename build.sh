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


########################################################################################################################
#docker build -t member:v1  -f ./service/member/Dockerfile .
#
#docker build -t hello:v1  -f ./service/hello/Dockerfile .
#
#docker build -t business:v1 -f ./openapi/business/Dockerfile .



# 注意：如果打包成镜像，代码中所有访问另一个容器的地方，比如：localhost:8081 需要改成 ${容器别名}:${容器内部端口}
#
# docker network create my_default
#
# docker run -it --rm --name businesstest -p 8080:8080 --network my_default --network-alias business business:v1
#
# docker run -it --rm --name hellotest -p 8081:8081 --network my_default --network-alias hello hello:v1
#
# docker run -it --rm --name membertest -p 8082:8082 --network my_default --network-alias member member:v1
#
#
# docker run -itd --name mongotest -p 27012:27017 --network my_default --network-alias mongo mongo:latest --auth
# docker exec -it mongotest mongo admin
# db.createUser({ user:'admin',pwd:'123456',roles:[ { role:'userAdminAnyDatabase', db: 'admin'},"readWriteAnyDatabase"]});
# db.auth('admin', '123456')
########################################################################################################################