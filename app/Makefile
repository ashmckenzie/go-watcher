PROJECT_NAME ?= project_name
BINARY_RELEASE = release/$(PROJECT_NAME)
BINARY_VERSIONED_RELEASE = $(BINARY_RELEASE)_$(VERSION)

BUILD_IMAGE_NAME = $(PROJECT_NAME)-build

APP_DIR ?= $(PWD)

BIN_DIR = $(APP_DIR)/bin
RELEASE_DIR = $(APP_DIR)/release

BUILD_USER = builder
CURRENT_USER_UID = 1000

VERSION=$(shell cat VERSION)

.PHONY: create_build_image deps update_deps release_dir bin_dir run clean

default: bin_dir deps
	go build -o $(BIN_DIR)/${PROJECT_NAME}

static: release_dir deps
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY_VERSIONED_RELEASE}

release: create_build_image
	mkdir -p ${RELEASE_DIR}
	docker run --rm -v ${RELEASE_DIR}:/release -e BUILD_USER=${BUILD_USER} ${BUILD_IMAGE_NAME} make extract_releases

create_build_image:
	docker build --build-arg PROJECT_NAME=${PROJECT_NAME} --build-arg CURRENT_USER_UID=${CURRENT_USER_UID} --build-arg BUILD_USER=${BUILD_USER} -f Dockerfile -t ${BUILD_IMAGE_NAME} .

deps:
	go get -d ./...

update_deps:
	go get -u -d ./...

bin_dir:
	mkdir -p $(BIN_DIR)

release_dir:
	mkdir -p ${RELEASE_DIR}

run: deps
	go run main.go $(filter-out $@, $(MAKECMDGOALS))

clean:
	rm -rf $(BIN_DIR) $(RELEASE_DIR)
