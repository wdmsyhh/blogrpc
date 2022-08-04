ABSOLUTE_PATH=$(cd `dirname $0`; pwd)

genProto() {
    domain=$1
    protoc -I=${ABSOLUTE_PATH}/${domain} --go_out ${ABSOLUTE_PATH}/${domain} --go_opt paths=source_relative --go-grpc_out ${ABSOLUTE_PATH}/${domain} --go-grpc_opt paths=source_relative ${ABSOLUTE_PATH}/${domain}/*.proto
    protoc -I=${ABSOLUTE_PATH}/${domain} --grpc-gateway_out ${ABSOLUTE_PATH}/${domain} --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt grpc_api_configuration=${ABSOLUTE_PATH}/${domain}/business_service.yaml ${ABSOLUTE_PATH}/${domain}/*.proto
}

genProto hello
genProto member