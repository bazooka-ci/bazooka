default: images

devimages:
	./scripts/build-devimages.sh

images:
	./scripts/build-images.sh

setup:
	./scripts/dev-setup.sh

run:
	./scripts/run.sh

errcheck:
	./scripts/errcheck.sh

updatedeps:
	go get -u -v ./...

test: errcheck devimages # Include errcheck in build phase
