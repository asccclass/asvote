COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
APP?=app
PORT?=10028
RELEASE?=0.0.1
ImageName?=justgps/openhousevote
ContainerName?=vote
MKFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURDIR := $(dir $(MKFILE))

cleanDocker:
	sh clean.sh

clean:
	rm -f ${APP}

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 GO111MODULE=off go build -a -tags netgo \
	-ldflags "-s -w -X version.Release=${RELEASE} \
	-X version.Commit=${COMMIT} \
	-X version.BuildTime=${BUILD_TIME}" \
	-o ${APP}

docker: build
	docker build -t ${ImageName}:${RELEASE} .
	rm -f ${APP}
	docker images

run: docker cleanDocker
	docker run -d --rm --name ${ContainerName} \
	-v /etc/localtime:/etc/localtime:ro \
	-v /etc/ssl/certs:/etc/ssl/certs \
	-v /etc/pki/ca-trust/extracted/pem:/etc/pki/ca-trust/extracted/pem \
	-v /etc/pki/ca-trust/extracted/openssl:/etc/pki/ca-trust/extracted/openssl \
	-v ${CURDIR}www:/app/www  \
	-v ${CURDIR}data:/app/data  \
	-v ${CURDIR}envfile:/app/envfile  \
	-p ${PORT}:80 \
	--env-file ${CURDIR}envfile \
	${ImageName}:${RELEASE}
	make log

stop:
	docker stop ${ContainerName}

log:
	 docker logs -f -t --tail 20 ${ContainerName}
rm: stop
	docker rm ${ContainerName}

rmi:
	docker rmi ${ImageName}:${RELEASE}
	
re:stop run

rebuild:stop rm rmi

s:
	git push -u origin master
