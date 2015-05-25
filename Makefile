default: images

MODULES = cli orchestration parser server worker web

.PHONY: commons modules $(MODULES)

commons:
	cd commons && go install ./...

client: commons
	cd client && go install ./...

all: $(MODULES)

$(MODULES):
	$(MAKE) -C $@

cli: client commons
orchestration: client commons
parser: commons
server: commons
worker: client commons

devimages: server parser orchestration worker

images: server parser orchestration worker web

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

cli-gox:
	gox -os="linux" github.com/bazooka-ci/bazooka/cli/bzk
	gox -os="darwin" github.com/bazooka-ci/bazooka/cli/bzk

git-server:
	-docker rm -vf bzk_git_server
	docker run -d --name bzk_git_server -p 9418:9418 -v /Users/jawher/temp/bzk-example:/repo bazooka/e2e-git