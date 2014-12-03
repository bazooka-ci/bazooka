default: images

devimages:
	./build-devimages.sh

images:
	./build-images.sh

setup:
	./dev-setup.sh

run:
	./run.sh

godep-save:
	./godep-save.sh

godep-restore:
	./godep-restore.sh
