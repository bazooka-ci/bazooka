default: docker

docker:
	$(PREFIX) docker build -t bazooka/orchestration .

devimage: gox
	$(PREFIX) docker build -t bazooka/orchestration dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"

push:
	$(PREFIX) docker tag bazooka/orchestration $(BZK_REGISTRY_HOST)/bazooka/orchestration
	$(PREFIX) docker push $(BZK_REGISTRY_HOST)/bazooka/orchestration
