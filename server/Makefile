default: docker

docker:
	$(PREFIX) docker build -t bazooka/server .


devimage: gox
	$(PREFIX) docker build -t bazooka/server dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"

push:
	$(PREFIX) docker tag bazooka/server $(BZK_REGISTRY_HOST)/bazooka/server
	$(PREFIX) docker push $(BZK_REGISTRY_HOST)/bazooka/server
