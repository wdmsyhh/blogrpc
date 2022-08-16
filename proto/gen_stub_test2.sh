#!/bin/bash

ABSOLUTE_PATH=$(cd `dirname $0`; pwd)

protoc -I=${ABSOLUTE_PATH}/hello \
    --go_out ${ABSOLUTE_PATH}/hello \
    --go_opt paths=source_relative \
    --go-grpc_out ${ABSOLUTE_PATH}/hello \
    --go-grpc_opt paths=source_relative \
    ${ABSOLUTE_PATH}/hello/*.proto

protoc -I=${ABSOLUTE_PATH}/hello \
    --grpc-gateway_out ${ABSOLUTE_PATH}/hello \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt grpc_api_configuration=${ABSOLUTE_PATH}/hello/business_service.yaml \
    ${ABSOLUTE_PATH}/hello/*.proto


protoc -I=${ABSOLUTE_PATH}/member \
  --go_out ${ABSOLUTE_PATH}/member \
  --go_opt paths=source_relative \
  --go-grpc_out ${ABSOLUTE_PATH}/member \
  --go-grpc_opt paths=source_relative \
  ${ABSOLUTE_PATH}/member/*.proto

protoc -I=${ABSOLUTE_PATH}/member \
  --grpc-gateway_out ${ABSOLUTE_PATH}/member \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt paths=source_relative \
  --grpc-gateway_opt grpc_api_configuration=${ABSOLUTE_PATH}/member/business_service.yaml \
  ${ABSOLUTE_PATH}/member/*.proto