#!/usr/bin/env bash

docker-compose -f docker-compose-infra.yml up -d
sleep 2
docker ps

export FN_API_URL="http://$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' twitterdaemon_fnserver_1):8080"
export POSTGRES_HOST=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' twitterdaemon_postgres_1)

source ${TWITTER_CONFIG:-"${GOPATH}/../omage2_twitter.rc"}

docker-compose -f docker-compose-twitter.yml up
