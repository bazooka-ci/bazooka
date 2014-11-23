default: docker

docker:
	docker build -t bazooka/parser-java .

devimage: gox
	docker build -t bazooka/parser-java dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"
