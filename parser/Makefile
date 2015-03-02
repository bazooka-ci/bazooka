default: docker

docker:
	docker build -t bazooka/parser -f Dockerfile ..

devimage: gox
	docker build -t bazooka/parser -f Dockerfile.dev .

gox:
	gox -osarch="linux/amd64" -output="main"

push:
	docker push bazooka/parser
