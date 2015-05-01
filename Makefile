default: images

.PHONY: cli devimages images docs web server orchestration parser

cli:
	cd cli && make

docs:
	mkdocs build

devimages: server parser orchestration

images: server parser orchestration web

setup:
	./scripts/dev-setup.sh

errcheck:
	./scripts/errcheck.sh

web:
	cd web && make

server:
	cd server && make

orchestration:
	cd orchestration && make

parser:
	cd parser && make


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
