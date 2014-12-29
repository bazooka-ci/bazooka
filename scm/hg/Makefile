default: docker

docker:
	$(PREFIX) docker build -t bazooka/scm-hg .

devimage: docker

push:
	$(PREFIX) docker tag bazooka/scm-hg $(BZK_REGISTRY_HOST)/bazooka/scm-hg
	$(PREFIX) docker push $(BZK_REGISTRY_HOST)/bazooka/scm-hg
