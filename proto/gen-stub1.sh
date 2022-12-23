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


########===================================================########
#  protoc 安装(任何版本应该应可以)：
  ## 下载 https://github.com/protocolbuffers/protobuf/releases/tag/v3.0.2 把 protoc 可执行文件 cp 到 /usr/local/bin
#  protoc-gen-go 插件安装:
  ## 我这里用的是老版本 https://github.com/golang/protobuf/releases/tag/v1.3.3
  ## 下载下来之后 cd 到 protoc-gen-go 目录执行 go build, go install 在 GOBIN 目录下会有 protoc-gen-go 可执行文件, 我习惯 cp 到 /usr/local/bin 中
  ## 也可以直接 go install github.com/golang/protobuf/protoc-gen-go@v1.3.3 安装到 GOBIN 目录中
# protoc-gen-grpc-gateway 安装：
  ## 老版本: go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.10.0  不加版本号的时候需要在 go module 项目下执行
  ## 新版本 go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest 不加版本号的时候需要在 go module 项目下执行
  ## 或者下载 https://github.com/grpc-ecosystem/grpc-gateway/releases/tag/v1.10.0  进入到 protoc-gen-grpc-gateway 目录执行 go build 和 go install
  ## 直接下载可执行文件改名为 protoc-gen-grpc-gateway，但是需要 chmod +x protoc-gen-grpc-gateway 赋予可执行权限
  ## 会在 GOBIN 目录下生成可执行文件 protoc-gen-grpc-gateway
########===================================================########
## 下面命令使用的是老版本的 protoc-gen-go 插件
## 需要 cd 到 proto 执行才能把生成的 pb 文件跟 *.proto 文件在同一目录（主要是尾部写的是相对目录 ./common/response/response.proto）
##
#command_str="protoc -I. -I${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mhello/hello.proto=blogrpc/proto/hello,Mhello/service.proto=blogrpc/proto/hello,Mmember/member.proto=blogrpc/proto/member,Mmember/service.proto=blogrpc/proto/member,,plugins=grpc+tag:. ./common/response/response.proto"
#echo ${command_str}
#$command_str
#command_str="protoc -I. -I${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mhello/hello.proto=blogrpc/proto/hello,Mhello/service.proto=blogrpc/proto/hello,Mmember/member.proto=blogrpc/proto/member,Mmember/service.proto=blogrpc/proto/member,,plugins=grpc+tag:. ./hello/hello.proto ./hello/service.proto"
#echo ${command_str}
#$command_str

## 需要 cd 到 proto 执行才能把生成的 pb 文件跟 *.proto 文件在同一目录，好像不设置 M 也行
#command_str="protoc -I. -I${ABSOLUTE_PATH} --go_out=plugins=grpc+tag:. ./common/response/response.proto"
#echo ${command_str}
#$command_str
#command_str="protoc -I. -I${ABSOLUTE_PATH} --go_out=plugins=grpc+tag:. ./hello/hello.proto ./hello/service.proto"
#echo ${command_str}
#$command_str
#command_str="protoc -I. -I${ABSOLUTE_PATH} --go_out=plugins=grpc+tag:. ./member/member.proto ./member/service.proto"
#echo ${command_str}
#$command_str

## 在项目目录 blogrpc 下执行
## 把 pb 文件生成到 *.proto 同一目录下, 尾部的 proto 文件必须在一个包下
#protoc -I. -I${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=${ABSOLUTE_PATH}/common/response,,plugins=grpc+tag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/common/response/response.proto
#protoc -I. -I${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mhello/hello.proto=blogrpc/proto/hello,Mhello/service.proto=blogrpc/proto/hello,Mmember/member.proto=blogrpc/proto/member,Mmember/service.proto=blogrpc/proto/member,,plugins=grpc+tag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/hello/hello.proto ${ABSOLUTE_PATH}/hello/service.proto
#protoc -I. -I${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mhello/hello.proto=blogrpc/proto/hello,Mhello/service.proto=blogrpc/proto/hello,Mmember/member.proto=blogrpc/proto/member,Mmember/service.proto=blogrpc/proto/member,,plugins=grpc+tag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/member/member.proto ${ABSOLUTE_PATH}/member/service.proto

## 在项目目录 blogrpc 下执行
## 把 pb 文件生成到 *.proto 同一目录下， 尾部的 proto 文件必须在一个包下
#protoc -I=${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=${ABSOLUTE_PATH}/common/response,,plugins=grpc+tag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/common/response/response.proto
#protoc -I=${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mhello/hello.proto=blogrpc/proto/hello,Mhello/service.proto=blogrpc/proto/hello,Mmember/member.proto=blogrpc/proto/member,Mmember/service.proto=blogrpc/proto/member,,plugins=grpc+tag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/hello/hello.proto ${ABSOLUTE_PATH}/hello/service.proto
#protoc -I=${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mhello/hello.proto=blogrpc/proto/hello,Mhello/service.proto=blogrpc/proto/hello,Mmember/member.proto=blogrpc/proto/member,Mmember/service.proto=blogrpc/proto/member,,plugins=grpc+tag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/member/member.proto ${ABSOLUTE_PATH}/member/service.proto

## 既可在项目目录 blogrpc 下执行，也可在 proto 目录下执行
## 要想把 pb 文件生成到 *.proto 同一目录下， 尾部的 proto 文件必须在一个包下
protoc -I=${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/common/response,,plugins=grpc+retag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/common/response/response.proto
protoc -I=${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mhello/hello.proto=blogrpc/proto/hello,Mhello/service.proto=blogrpc/proto/hello,Mmember/member.proto=blogrpc/proto/member,Mmember/service.proto=blogrpc/proto/member,,plugins=grpc+retag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/hello/hello.proto ${ABSOLUTE_PATH}/hello/service.proto
protoc -I=${ABSOLUTE_PATH} --go_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mhello/hello.proto=blogrpc/proto/hello,Mhello/service.proto=blogrpc/proto/hello,Mmember/member.proto=blogrpc/proto/member,Mmember/service.proto=blogrpc/proto/member,,plugins=grpc+retag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/member/member.proto ${ABSOLUTE_PATH}/member/service.proto

## 或者 -I 后面可以不写 = ，也可以没有空格  --go_out 后面也可以不写 =
#protoc -I ${ABSOLUTE_PATH} --go_out=plugins=grpc+tag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/common/response/response.proto
#protoc -I ${ABSOLUTE_PATH} --go_out=plugins=grpc+tag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/hello/hello.proto ${ABSOLUTE_PATH}/hello/service.proto
#protoc -I ${ABSOLUTE_PATH} --go_out=plugins=grpc+tag:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/member/member.proto ${ABSOLUTE_PATH}/member/service.proto

#command_str="${COMMAND_START} --grpc-gateway_out=${package_map}grpc_api_configuration=$folder/${target_config},allow_delete_body=true:. $folder/*.proto"
#echo ${command_str}

## 需要 cd 到 proto 执行才能把生成的 pb 文件跟 *.proto 文件在同一目录，好像不设置 M 也行
## 把 pb.gw 文件生成到 *.proto 同一目录下， 尾部的 proto 文件必须在一个包下
#protoc -I. -I/${ABSOLUTE_PATH} --grpc-gateway_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mcommon/types/types.proto=blogrpc/proto/common/types,grpc_api_configuration=./hello/business_service.yaml,allow_delete_body=true:. ./hello/service.proto ./hello/hello.proto
#protoc -I. -I/${ABSOLUTE_PATH} --grpc-gateway_out=Mcommon/benchmark/benchmark.proto=blogrpc/proto/common/benchmark,Mcommon/ec/ec.proto=blogrpc/proto/common/ec,Mcommon/origin/origin.proto=blogrpc/proto/common/origin,Mcommon/request/request.proto=blogrpc/proto/common/request,Mcommon/response/response.proto=blogrpc/proto/common/response,Mcommon/types/types.proto=blogrpc/proto/common/types,grpc_api_configuration=./member/business_service.yaml,allow_delete_body=true:. ./member/service.proto ./member/member.proto

## 在 blogrpc 目录下执行 ./proto/gen-stub1.sh 可以，但是 cd 到 proto 中执行 ./gen-stub1.sh 会报错：
### /home/user/GolandProjects/blogrpc/proto/hello/service.proto: Input is shadowed in the --proto_path by "hello/service.proto".  Either use the latter file as your input or reorder the --proto_path so that the former file's location comes first.
#protoc -I. -I/${ABSOLUTE_PATH} --grpc-gateway_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mcommon/types/types.proto=blogrpc/proto/common/types,grpc_api_configuration=${ABSOLUTE_PATH}/hello/business_service.yaml,allow_delete_body=true:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/hello/service.proto ${ABSOLUTE_PATH}/hello/hello.proto
#protoc -I. -I/${ABSOLUTE_PATH} --grpc-gateway_out=Mcommon/request/request.proto=blogrpc/proto/common/request,Mcommon/response/response.proto=blogrpc/proto/common/response,Mcommon/types/types.proto=blogrpc/proto/common/types,grpc_api_configuration=${ABSOLUTE_PATH}/member/business_service.yaml,allow_delete_body=true:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/member/service.proto ${ABSOLUTE_PATH}/member/member.proto

## 既可在项目目录 blogrpc 下执行，也可在 proto 目录下执行
## 要想把 pb 文件生成到 *.proto 同一目录下， 尾部的 proto 文件必须在一个包下
protoc -I=${ABSOLUTE_PATH} --grpc-gateway_out=Mcommon/response/response.proto=blogrpc/proto/common/response,Mcommon/types/types.proto=blogrpc/proto/common/types,grpc_api_configuration=${ABSOLUTE_PATH}/hello/business_service.yaml,allow_delete_body=true:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/hello/service.proto ${ABSOLUTE_PATH}/hello/hello.proto
protoc -I=${ABSOLUTE_PATH} --grpc-gateway_out=Mcommon/request/request.proto=blogrpc/proto/common/request,Mcommon/response/response.proto=blogrpc/proto/common/response,Mcommon/types/types.proto=blogrpc/proto/common/types,grpc_api_configuration=${ABSOLUTE_PATH}/member/business_service.yaml,allow_delete_body=true:${ABSOLUTE_PATH} ${ABSOLUTE_PATH}/member/service.proto ${ABSOLUTE_PATH}/member/member.proto

}

genProto hello
#genProto member