#!/bin/sh

docker rmi $(docker images | grep "none" | awk '{print $3}')
clear
docker ps -a
ls -al
