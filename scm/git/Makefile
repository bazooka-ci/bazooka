default: docker

docker:
	$(PREFIX) docker build -t bazooka/scm-git .

devimage: docker

push:
	$(PREFIX) docker tag bazooka/scm-git $(BZK_REGISTRY_HOST)/bazooka/scm-git
	$(PREFIX) docker push $(BZK_REGISTRY_HOST)/bazooka/scm-git
