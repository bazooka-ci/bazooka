default: images

devimages:
	./build-devimages.sh

images:
	./build-images.sh

setup:
	./dev-setup.sh

run:
	./run.sh

godepsave:
	./godep-save.sh

godeprestore:
	./godep-restore.sh
