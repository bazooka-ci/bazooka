default: docker

docker:
	$(PREFIX) docker build -t bazooka/parser-nodejs .

devimage: gox
	$(PREFIX) docker build -t bazooka/parser-nodejs dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"

push:
	$(PREFIX) docker tag bazooka/parser-nodejs $(BZK_REGISTRY_HOST)/bazooka/parser-nodejs
	$(PREFIX) docker push $(BZK_REGISTRY_HOST)/bazooka/parser-nodejs
