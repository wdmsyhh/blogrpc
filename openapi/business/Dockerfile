# 启动编译环境
FROM golang:1.17-alpine AS builder

# 配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 拷贝源代码到镜像中
COPY . /go/src/blogrpc

# 编译
WORKDIR /go/src/blogrpc
RUN go install ./openapi/business/...

FROM alpine:3.13
COPY --from=builder /go/bin/business /bin/business

# 申明暴露的端口
EXPOSE 8080

# 设置服务入口
ENTRYPOINT [ "/bin/business" ]