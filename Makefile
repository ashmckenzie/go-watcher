PROJECT_NAME ?= project_name
DOCKER_HUB_USERNAME ?= docker_hub_username
DOCKER_IMAGE_VERSION ?= latest

DOCKER_IMAGE_NAME = ${DOCKER_HUB_USERNAME}/${PROJECT_NAME}
DOCKER_FILE = Dockerfile
DOCKER_PORTS ?= -p 5000:5000

.PHONY: run shell attach clean

build:
	make -C app release
	docker build -f ${DOCKER_FILE} -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION} .

run: build
	docker run --rm -ti ${DOCKER_PORTS} -e PROJECT_NAME=${PROJECT_NAME} ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}

shell: build
	docker run --rm -ti ${DOCKER_PORTS} -e PROJECT_NAME=${PROJECT_NAME} ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION} bash

attach: build
	docker exec -ti `docker ps | grep '${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}' | awk '{ print $$1 }'` bash

clean:
	docker rmi ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}
