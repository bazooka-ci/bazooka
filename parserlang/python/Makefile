default: docker

docker:
	docker build -t bazooka/parser-python -f Dockerfile ../..

devimage: gox
	docker build -t bazooka/parser-python -f Dockerfile.dev .

gox:
	gox -osarch="linux/amd64" -output="main"

push:
	docker tag bazooka/parser-python bazooka/parser-python
	docker push bazooka/parser-python
