default: docker

docker:
	$(PREFIX) docker build -t bazooka/server .

devimage: gox
	$(PREFIX) docker build -t bazooka/server dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"

godep-restore:
	godep restore ./...

godep-save:
	godep save ./...
