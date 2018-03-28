#!/usr/bin/env bash


docker rm -f pg | true
docker rm -f minio | true
docker rm -f fnserver | true

docker run --name pg -e "POSTGRES_DB=emokognition" -e "POSTGRES_PASSWORD=postgres" -e "POSTGRES_USER=postgres" -p 5432:5432 -d postgres:9.3-alpine
docker run -d -p 9000:9000 --name minio -e "MINIO_ACCESS_KEY=admin" -e "MINIO_SECRET_KEY=password" minio/minio server /data
sleep 10
docker run -v /var/run/docker.sock:/var/run/docker.sock -d -p 8080:8080 --name fnserver --link minio:minio -e "FN_LOG_STORE=s3://admin:password@minio:9000/us-east-1/fnlogs" -e "FN_LOG_LEVEL=DEBUG" fnproject/fnserver

sleep 2
docker inspect -f {{.State.Running}} fnserver | grep '^true$'
sleep 5

fn version

source ${GOPATH}/../omage2_twitter.rc

go build

export FN_API_URL="http://localhost:8080"
export INTERNAL_FN_API_URL="http://$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' fnserver):8080"

export POSTGRES_HOST=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' pg)
export POSTGRES_PORT=5432
export POSTGRES_DB=emokognition
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres

export TwitterBotType=${TwitterBotType:-emokognition}
export InitialTweetID=925300910162669567

./twitter-daemon
