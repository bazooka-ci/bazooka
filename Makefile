default: images

.PHONY: scm runner docs web

devimages:
	./scripts/build-devimages.sh

images:
	./scripts/build-images.sh

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

push:
	./scripts/push-images.sh

deploy: setup devimages runner scm web push

updatedeps:
	go get -u -v ./...

test: errcheck devimages # Include errcheck in build phase
