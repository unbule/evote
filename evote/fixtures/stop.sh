#!/bin/bash

echo "关闭网络"
docker-compose -f docker-compose.yml down --volumes --remove-orphans
#docker rm -f $(docker ps)
docker volume prune
docker network prune

docker rm $(docker ps -aq)
docker rmi $(docker images dev-* -q)

rm -rf /tmp/evote-store
rm -rf /tmp/evote-msp


