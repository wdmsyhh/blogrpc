ABSOLUTE_PATH=$(cd `dirname $0`; pwd)

lm_traverse_dir(){
	#for file in $(ls $1)		#与下面一行等价
    for file in `ls $1`       	#注意两个反引号，获取命令执行的结果
    do
        if [ -d $1"/"$file ]  	#"-d" 判断是否为目录，注意此处之间一定要加上空格，否则会报错
        then
            if ls $1"/"$file/*.proto >/dev/null 2>&1;then # 判断当前文件夹中是否存在以 .proto 结尾的文件
              # "/" 双引号也可以不加
              protoc -I=$1"/"$file --go_out $1"/"$file --go_opt paths=source_relative --go-grpc_out $1"/"$file --go-grpc_opt paths=source_relative $1"/"$file/*.proto
            fi
            lm_traverse_dir $1"/"$file	#遍历子目录
        fi
    done
}
# 执行命令
#lm_traverse_dir ${ABSOLUTE_PATH}/common

genProto() {
    domain=$1
#    protoc -I=${ABSOLUTE_PATH}/${domain} --go_out ${ABSOLUTE_PATH}/${domain} --go_opt paths=source_relative --go-grpc_out ${ABSOLUTE_PATH}/${domain} --go-grpc_opt paths=source_relative ${ABSOLUTE_PATH}/${domain}/*.proto
#    protoc -I=${ABSOLUTE_PATH}/${domain} --grpc-gateway_out ${ABSOLUTE_PATH}/${domain} --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt grpc_api_configuration=${ABSOLUTE_PATH}/${domain}/business_service.yaml ${ABSOLUTE_PATH}/${domain}/*.proto

#    protoc -I=${ABSOLUTE_PATH} --go_out ${ABSOLUTE_PATH} --go_opt paths=source_relative --go-grpc_out ${ABSOLUTE_PATH} --go-grpc_opt paths=source_relative ${ABSOLUTE_PATH}/${domain}/*.proto
#    protoc -I=${ABSOLUTE_PATH} --grpc-gateway_out ${ABSOLUTE_PATH} --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt grpc_api_configuration=${ABSOLUTE_PATH}/${domain}/business_service.yaml ${ABSOLUTE_PATH}/${domain}/*.proto

command_str="protoc -I. -I${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/proto/common/response,plugins=grpc+tag:. ${ABSOLUTE_PATH}/common/response/response.proto"
echo ${command_str}
$command_str

#protoc -I. -I${ABSOLUTE_PATH} --go_out=Mcommon/response.proto=${ABSOLUTE_PATH}/common/response,plugins=grpc+tag:. ./common/response/response.proto
}

genProto hello
#genProto member