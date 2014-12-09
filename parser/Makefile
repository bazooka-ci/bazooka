default: docker

docker:
	$(PREFIX) docker build -t bazooka/parser .

devimage: gox
	$(PREFIX) docker build -t bazooka/parser dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"

push:
	$(PREFIX) docker tag bazooka/parser $(BZK_REGISTRY_HOST)/bazooka/parser
	$(PREFIX) docker push $(BZK_REGISTRY_HOST)/bazooka/parser
