docker build -t member:v1  -f ./service/member/Dockerfile .

docker build -t hello:v1  -f ./service/hello/Dockerfile .

docker build -t business:v1 -f ./openapi/business/Dockerfile .

