#!/usr/bin/env bash

docker-compose -f docker-compose-twitter.yml down
sleep 2
docker-compose -f docker-compose-infra.yml down
