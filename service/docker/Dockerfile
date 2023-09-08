ARG REGISTRY_HOST
ARG ENV
FROM registry.ap-southeast-1.aliyuncs.com/yhhnamespace/base-ubuntu:${ENV}

EXPOSE 1701
ARG rpcname
WORKDIR /app/${rpcname}

COPY scripts/start /etc/service/${rpcname}/run
COPY ${rpcname}/conf /app/${rpcname}/conf
COPY ${rpcname}/assets /app/${rpcname}/assets
COPY ${rpcname}/${rpcname} /app/${rpcname}/
