#!/bin/bash
serv_name=openapi
HOST=${HOST:-"0.0.0.0"}
PORT=${PORT:-9091}
ENV=${ENV:-local}

cd /app/openapi
echo "Service $serv_name is starting at $HOST:$PORT at env $ENV"
exec ./$serv_name --env=$ENV --port=$PORT --host=$HOST
