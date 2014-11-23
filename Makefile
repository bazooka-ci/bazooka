default: images

devimages:
	./build-devimages.sh

images:
	./build-images.sh

setup:
	./dev-setup.sh

run:
	./run.sh
