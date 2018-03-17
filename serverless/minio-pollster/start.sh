#!/usr/bin/env bash

docker rm -f pg | true
docker rm -f minio | true
docker rm -f fnserver | true

docker run --name pg -e "POSTGRES_DB=emokognition" -e "POSTGRES_PASSWORD=postgres" -e "POSTGRES_USER=postgres" -p 5432:5432 -d postgres:9.3-alpine
docker run -d -p 9000:9000 --name minio -e "MINIO_ACCESS_KEY=admin" -e "MINIO_SECRET_KEY=password" minio/minio server /data

echo -e "awaiting for Minio to become healthy..."
sleep 10

docker run -v /var/run/docker.sock:/var/run/docker.sock -d -p 8080:8080 --name fnserver --link minio:minio -e "FN_LOG_STORE=s3://admin:password@minio:9000/us-east-1/fnlogs" -e "FN_LOG_LEVEL=DEBUG" fnproject/fnserver
sleep 5

docker inspect -f {{.State.Running}} pg | grep '^true$'
docker inspect -f {{.State.Running}} minio | grep '^true$'
docker inspect -f {{.State.Running}} fnserver | grep '^true$'

fn version

export MINIO_URL="s3://admin:password@localhost:9000/us-east-1/emotions"
export FN_API_URL="http://localhost:8080"
export INTERNAL_FN_API_URL="http://$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' fnserver):8080"

export POSTGRES_HOST=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' pg)
export POSTGRES_PORT=5432
export POSTGRES_DB=emokognition
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres

go build

./minio-pollster
