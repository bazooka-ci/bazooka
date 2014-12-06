default: docker

docker:
	$(PREFIX) docker build -t bazooka/parser-golang .

devimage: gox
	$(PREFIX) docker build -t bazooka/parser-golang dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"

push:
	$(PREFIX) docker tag bazooka/parser-golang $(BZK_REGISTRY_HOST)/bazooka/parser-golang
	$(PREFIX) docker push $(BZK_REGISTRY_HOST)/bazooka/parser-golang
