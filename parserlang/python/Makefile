default: docker

docker:
	docker build -t bazooka/parser-python .

devimage: gox
	docker build -t bazooka/parser-python dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"

push:
	docker tag bazooka/parser-python bazooka/parser-python
	docker push bazooka/parser-python
