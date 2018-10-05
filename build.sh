#!/usr/bin/env bash

NAME="limiter"
TAG=latest
USER="kgantsov"
DOCKER_ID_USER="kgantsov"

docker build -t $USER/$NAME:$TAG --no-cache .

docker tag $USER/$NAME:$TAG $USER/$NAME:$TAG
docker push $USER/$NAME:$TAG


docker rmi $USER/$NAME:$TAG
