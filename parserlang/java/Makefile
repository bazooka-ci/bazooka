default: docker

docker:
	$(PREFIX) docker build -t bazooka/parser-java .

devimage: gox
	$(PREFIX) docker build -t bazooka/parser-java dev/

gox:
	gox -osarch="linux/amd64" -output="dev/main"
