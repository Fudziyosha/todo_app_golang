main_package_path = ./cmd/app
binary_name = todo-app
APP_NAME=todo-app
IMAGE_NAME=$(APP_NAME)
TAG=latest
project_name = web_todos
DC = docker-compose

.PHONY: help tidy build run d-build dm-up dm-build
help:
	@awk 'BEGIN {FS=":"} \
	/^#/ {comment=substr($$0,3)} \
	/^[a-zA-Z0-9_-]+:/ {printf "\033[36m%-20s\033[0m %s\n", $$1, comment}' Makefile

tidy:
	go mod tidy -v
	go fmt ./...

build:
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

run: build
	/tmp/bin/${binary_name}

d-build:
	docker buildx build -t $(IMAGE_NAME):$(TAG) .

dm-up:
	$(DC) up -d

dm-build:
	$(DC) build