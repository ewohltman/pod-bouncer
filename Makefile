.PHONY: lint test build docker push deploy

MAKEFILE_PATH=$(shell readlink -f "${0}")
MAKEFILE_DIR=$(shell dirname "${MAKEFILE_PATH}")

lint:
	golangci-lint run ./...

test:
	go test -v -race -coverprofile=coverage.out ./...

build:
	CGO_ENABLED=0 go build -o build/package/pod-bouncer/pod-bouncer cmd/pod-bouncer/pod-bouncer.go

image:
	docker image build -t ewohltman/pod-bouncer:latest build/package/pod-bouncer

push:
	docker login -u "${DOCKER_USER}" -p "${DOCKER_PASS}"
	docker push ewohltman/pod-bouncer:latest
	docker logout

deploy:
	${MAKEFILE_DIR}/scripts/deploy.sh
