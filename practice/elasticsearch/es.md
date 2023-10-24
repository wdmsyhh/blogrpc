# ES

## 启动 es

### 无密码启动

- 启动 ES 容器
```shell
docker run --name elasticsearch --rm -p 9200:9200 -p 9300:9300 \
    -e "discovery.type=single-node" \
    -e "xpack.security.enabled=false" \
    -e "ES_JAVA_OPTS=-Xmx512m -Xms512m" \
    --net my_default \
    docker.elastic.co/elasticsearch/elasticsearch:7.10.2
```
- 访问测试
```shell
curl localhost:9200
```

## 启动 kibana

默认可以直接连接容器名为 elasticsearch 的 es
```shell
docker run \
    --rm \
    --name kibana \
    --net my_default \
    -p 5601:5601 \
    docker.elastic.co/kibana/kibana:7.10.2
```

### 启动并设置密码

- es

```shell
docker run -d --name elasticsearch \
  -p 9200:9200 \
  -e "discovery.type=single-node" \
  -e "ELASTIC_PASSWORD=root123" \
  -e "ES_JAVA_OPTS=-Xmx512m -Xms512m" \
  --net my_default \
  docker.elastic.co/elasticsearch/elasticsearch:7.10.2
```

- kibana

```shell
docker run --rm --name kibana \
  -p 5601:5601 \
  -e "ELASTICSEARCH_HOSTS=http://elasticsearch:9200" \
  -e "ELASTICSEARCH_USERNAME=elastic" \
  -e "ELASTICSEARCH_PASSWORD=root123" \
  --net my_default \
  docker.elastic.co/kibana/kibana:7.10.2
```

```shell
# 这将在需要时递归创建 dockerdata、kibana 和 config 这些目录，以确保它们都存在。
sudo mkdir -p /home/dockerdata/kibana

ls dockerdata/kibana/config

#docker cp kibana:/usr/share/kibana/config/kibana.yml /home/dockerdata/kibana/config/kibana.yml
sudo docker cp kibana:/usr/share/kibana/config /home/dockerdata/kibana/config

sudo chmod 777 dockerdata/kibana/config/kibana.yml
# 设置中文 kibana.yml
i18n.locale: "zh-CN"
```

- 挂载

`/host/path:/container/path`

```shell
#   -v /home/dockerdata/kibana/config:/usr/share/kibana/config \
docker run --rm --name kibana \
  -p 5601:5601 \
  -e "ELASTICSEARCH_HOSTS=http://elasticsearch:9200" \
  -e "ELASTICSEARCH_USERNAME=elastic" \
  -e "ELASTICSEARCH_PASSWORD=root123" \
  --net my_default \
  docker.elastic.co/kibana/kibana:7.10.2
```

```shell
docker run --rm --name kibana \
  -p 5601:5601 \
  -e "ELASTICSEARCH_HOSTS=http://infras-elasticsearch:9200" \
  --net scrm_default \
  docker.elastic.co/kibana/kibana:7.10.1
```

```shell
docker run --rm --name es-head \
  -p 9100:9100 \
  mobz/elasticsearch-head:5
```