default: images

.PHONY: cli devimages images docs web

cli:
	cd cli && make

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

web:
	cd web && make

push:
	./scripts/push-images.sh

deploy: setup devimages runner web push

updatedeps:
	go get -u -v ./...

test: errcheck devimages # Include errcheck in build phase

bintray: gox
	./scripts/push-bintray.sh

gox:
	gox -os="linux" github.com/bazooka-ci/bazooka/cli/bzk
	gox -os="darwin" github.com/bazooka-ci/bazooka/cli/bzk
