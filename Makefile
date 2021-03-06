default: images

MODULES = cli orchestration parser server web

.PHONY: commons modules $(MODULES)

commons:
	cd commons && go install ./...

client: commons
	cd client && go install ./...

all: $(MODULES)

$(MODULES):
	$(MAKE) -C $@

cli: client commons
orchestration: commons
parser: commons
server: commons

devimages: server parser orchestration

images: server parser orchestration web

setup:
	./scripts/dev-setup.sh

errcheck:
	./scripts/errcheck.sh

push:
	./scripts/push-images.sh

deploy: setup devimages web push

updatedeps:
	go get -u -v ./...

test: errcheck devimages # Include errcheck in build phase

push-bintray:
	./scripts/push-bintray.sh

cli-multiplatform:
	GOOS=linux  GOARCH=amd64 go build -o=bzk_linux_amd64  github.com/bazooka-ci/bazooka/cli/bzk
	GOOS=darwin GOARCH=amd64 go build -o=bzk_darwin_amd64 github.com/bazooka-ci/bazooka/cli/bzk
