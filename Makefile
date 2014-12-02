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

test: errcheck devimages # Include errcheck in build phase
