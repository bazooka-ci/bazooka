default: docker

docker:
	$(PREFIX) docker build -t bazooka/parser-python .

devimage: gox
	$(PREFIX) docker build -t bazooka/parser-python dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"

push:
	$(PREFIX) docker tag bazooka/parser-python $(BZK_REGISTRY_HOST)/bazooka/parser-python
	$(PREFIX) docker push $(BZK_REGISTRY_HOST)/bazooka/parser-python
