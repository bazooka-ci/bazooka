default: images

MODULES = cli orchestration parser server web

.PHONY: commons docs modules $(MODULES)

commons:
	cd commons && go install ./...

all: $(MODULES)

$(MODULES):
	$(MAKE) -C $@

cli: commons
orchestration: commons
parser: commons
server: commons

docs:
	mkdocs build

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

bintray: cli-gox
	./scripts/push-bintray.sh

cli-gox:
	gox -os="linux" github.com/bazooka-ci/bazooka/cli/bzk
	gox -os="darwin" github.com/bazooka-ci/bazooka/cli/bzk
