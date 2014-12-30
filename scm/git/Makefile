default: docker

docker:
	docker build -t bazooka/scm-git .

devimage: docker

push:
	docker push bazooka/scm-git
