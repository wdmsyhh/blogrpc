ARG REGISTRY_HOST
ARG ENV
FROM ${REGISTRY_HOST}/yhhnamespace/base-ubuntu:${ENV}

COPY . /app/openapi/
RUN mkdir -p /etc/service/openapi
COPY docker/start.sh /etc/service/openapi/run
EXPOSE 9091
WORKDIR /app/openapi
