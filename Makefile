default: images

.PHONY: scm runner docs web parser orchestration server

devimages:
	./scripts/build-devimages.sh

images: parser orchestration server parser-golang parser-java parser-nodejs parser-python

docs:
	mkdocs build

setup:
	./scripts/dev-setup.sh

errcheck:
	./scripts/errcheck.sh

runner:
	./scripts/build-runner.sh

scm:
	./scripts/build-scm.sh

web:
	cd web && make

parser:
	docker build -f Dockerfile-parser -t bazooka/parser .

orchestration:
	docker build -f Dockerfile-orchestration -t bazooka/orchestration .

server:
	docker build -f Dockerfile-server -t bazooka/server .

parser-golang:
	docker build -f Dockerfile-parser-golang -t bazooka/parser-golang .

parser-java:
	docker build -f Dockerfile-parser-java -t bazooka/parser-java .

parser-nodejs:
	docker build -f Dockerfile-parser-nodejs -t bazooka/parser-nodejs .

parser-python:
	docker build -f Dockerfile-parser-python -t bazooka/parser-python .

push:
	./scripts/push-images.sh

deploy: setup devimages runner scm web push

updatedeps:
	go get -u -v ./...

test: errcheck devimages # Include errcheck in build phase
