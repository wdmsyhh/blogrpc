ABSOLUTE_PATH=$(cd `dirname $0`; pwd)

lm_traverse_dir(){
	#for file in $(ls $1)		#与下面一行等价
    for file in `ls $1`       	#注意两个反引号，获取命令执行的结果
    do
        if [ -d $1"/"$file ]  	#"-d" 判断是否为目录，注意此处之间一定要加上空格，否则会报错
        then
            # "/" 双引号也可以不加
            if ls $1"/"$file/*.proto >/dev/null 2>&1;then # 判断当前文件夹中是否存在以 .proto 结尾的文件
              protoc -I=$1"/"$file --go_out $1"/"$file --go_opt paths=source_relative --go-grpc_out $1"/"$file --go-grpc_opt paths=source_relative $1"/"$file/*.proto
            fi

            lm_traverse_dir $1"/"$file	#遍历子目录
        #else
            # 可以在这里处理文件，比如改名、删除等
            #effect_name=$1"/"$file		#注意"="前后不要留空格
            #echo $effect_name			#输出文件名
            #rm -rf $effect_name
            #mv $effect_name "new_name"
        fi
    done
}

# 执行命令
lm_traverse_dir ${ABSOLUTE_PATH}/common